package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListAllowedCPCodesRequest contains the request parameters for the ListAllowedCPCodes endpoint.
	ListAllowedCPCodesRequest struct {
		UserName string
		Body     ListAllowedCPCodesRequestBody
	}

	// ListAllowedAPIsRequest contains the request parameters for the ListAllowedAPIs endpoint.
	ListAllowedAPIsRequest struct {
		UserName           string
		ClientType         ClientType
		AllowAccountSwitch bool
	}

	// ListAccessibleGroupsRequest contains the request parameter for the ListAccessibleGroups endpoint.
	ListAccessibleGroupsRequest struct {
		UserName string
	}

	// ListAllowedCPCodesRequestBody contains the filtering parameters for the ListAllowedCPCodes endpoint.
	ListAllowedCPCodesRequestBody struct {
		ClientType ClientType               `json:"clientType"`
		Groups     []ClientGroupRequestItem `json:"groups"`
	}

	// AllowedCPCodesGroup contains the group parameters for the ListAllowedCPCodes endpoint.
	AllowedCPCodesGroup struct {
		GroupID         int64                 `json:"groupId,omitempty"`
		RoleID          int64                 `json:"roleId,omitempty"`
		GroupName       string                `json:"groupName,omitempty"`
		IsBlocked       bool                  `json:"isBlocked,omitempty"`
		ParentGroupID   int64                 `json:"parentGroupId,omitempty"`
		RoleDescription string                `json:"roleDescription,omitempty"`
		RoleName        string                `json:"roleName,omitempty"`
		SubGroups       []AllowedCPCodesGroup `json:"subGroups,omitempty"`
	}

	// ListAllowedCPCodesResponse contains response for the ListAllowedCPCodes endpoint.
	ListAllowedCPCodesResponse []ListAllowedCPCodesResponseItem

	// ListAllowedCPCodesResponseItem contains single item of the response from the ListAllowedCPCodes endpoint.
	ListAllowedCPCodesResponseItem struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	// ListAuthorizedUsersResponse contains the response from the ListAuthorizedUsers endpoint.
	ListAuthorizedUsersResponse []AuthorizedUser

	// AuthorizedUser contains the details about the authorized user.
	AuthorizedUser struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		UIIdentityID string `json:"uiIdentityId"`
	}

	// ListAccessibleGroupsResponse contains the response from the ListAccessibleGroups endpoint.
	ListAccessibleGroupsResponse []AccessibleGroup

	// AccessibleGroup contains the details about accessible group.
	AccessibleGroup struct {
		GroupID         int64                `json:"groupId"`
		RoleID          int64                `json:"roleId"`
		GroupName       string               `json:"groupName"`
		RoleName        string               `json:"roleName"`
		IsBlocked       bool                 `json:"isBlocked"`
		RoleDescription string               `json:"roleDescription"`
		SubGroups       []AccessibleSubGroup `json:"subGroups"`
	}

	// AccessibleSubGroup contains the details about subgroup.
	AccessibleSubGroup struct {
		GroupID       int64                `json:"groupId"`
		GroupName     string               `json:"groupName"`
		ParentGroupID int64                `json:"parentGroupId"`
		SubGroups     []AccessibleSubGroup `json:"subGroups"`
	}

	// ListAllowedAPIsResponse contains the response from the ListAllowedAPIs endpoint.
	ListAllowedAPIsResponse []AllowedAPI

	// AllowedAPI contains the details about the API.
	AllowedAPI struct {
		AccessLevels      []AccessLevel `json:"accessLevels"`
		APIID             int64         `json:"apiId"`
		APIName           string        `json:"apiName"`
		Description       string        `json:"description"`
		DocumentationURL  string        `json:"documentationUrl"`
		Endpoint          string        `json:"endpoint"`
		HasAccess         bool          `json:"hasAccess"`
		ServiceProviderID int64         `json:"serviceProviderId"`
	}

	// ClientType represents the type of the client.
	ClientType string
)

const (
	// UserClientType is the `USER_CLIENT` client type.
	UserClientType ClientType = "USER_CLIENT"
	// ServiceAccountClientType is the `SERVICE_ACCOUNT` client type.
	ServiceAccountClientType ClientType = "SERVICE_ACCOUNT"
	// ClientClientType is the `CLIENT` client type.
	ClientClientType ClientType = "CLIENT"
)

var (
	// ErrListAllowedCPCodes is returned when ListAllowedCPCodes fails.
	ErrListAllowedCPCodes = errors.New("list allowed CP codes")
	// ErrListAuthorizedUsers is returned when ListAuthorizedUsers fails.
	ErrListAuthorizedUsers = errors.New("list authorized users")
	// ErrListAllowedAPIs is returned when ListAllowedAPIs fails.
	ErrListAllowedAPIs = errors.New("list allowed APIs")
	// ErrAccessibleGroups is returned when ListAccessibleGroups fails.
	ErrAccessibleGroups = errors.New("list accessible groups")
)

// Validate validates ClientType.
func (c ClientType) Validate() error {
	return validation.In(ClientClientType, ServiceAccountClientType, UserClientType).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' or '%s'",
			c, ClientClientType, ServiceAccountClientType, UserClientType)).
		Validate(c)
}

// Validate validates ListAllowedCPCodesRequest.
func (r ListAllowedCPCodesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"UserName": validation.Validate(r.UserName, validation.Required),
		"Body":     validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates ListAllowedCPCodesRequestBody.
func (r ListAllowedCPCodesRequestBody) Validate() error {
	return validation.Errors{
		"ClientType": validation.Validate(r.ClientType, validation.Required, validation.In(ClientClientType, UserClientType, ServiceAccountClientType).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'CLIENT' or 'USER_CLIENT' or 'SERVICE_ACCOUNT'", r.ClientType))),
		"Groups":     validation.Validate(r.Groups, validation.Required.When(r.ClientType == ServiceAccountClientType)),
	}.Filter()
}

// Validate validates ListAllowedAPIsRequest.
func (r ListAllowedAPIsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"UserName":   validation.Validate(r.UserName, validation.Required),
		"ClientType": validation.Validate(r.ClientType, validation.In(ClientClientType, UserClientType, ServiceAccountClientType).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'CLIENT' or 'USER_CLIENT' or 'SERVICE_ACCOUNT'", r.ClientType))),
	})
}

// Validate validates ListAccessibleGroupsRequest.
func (r ListAccessibleGroupsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"UserName": validation.Validate(r.UserName, validation.Required),
	})
}

func (i *iam) ListAllowedCPCodes(ctx context.Context, params ListAllowedCPCodesRequest) (ListAllowedCPCodesResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListAllowedCPCodes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListAllowedCPCodes, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/identity-management/v3/users/%s/allowed-cpcodes", params.UserName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAllowedCPCodes, err)
	}

	var result ListAllowedCPCodesResponse
	resp, err := i.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAllowedCPCodes, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAllowedCPCodes, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListAuthorizedUsers(ctx context.Context) (ListAuthorizedUsersResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListAuthorizedUsers")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/identity-management/v3/users", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAuthorizedUsers, err)
	}

	var result ListAuthorizedUsersResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAuthorizedUsers, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAuthorizedUsers, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListAllowedAPIs(ctx context.Context, params ListAllowedAPIsRequest) (ListAllowedAPIsResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListAllowedAPIs")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListAllowedAPIs, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/users/%s/allowed-apis", params.UserName))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAllowedAPIs, err)
	}

	q := uri.Query()
	if params.ClientType != "" {
		q.Add("clientType", string(params.ClientType))

	}
	q.Add("allowAccountSwitch", strconv.FormatBool(params.AllowAccountSwitch))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAllowedAPIs, err)
	}

	var result ListAllowedAPIsResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAllowedAPIs, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAllowedAPIs, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListAccessibleGroups(ctx context.Context, params ListAccessibleGroupsRequest) (ListAccessibleGroupsResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListAccessibleGroups")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrAccessibleGroups, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/identity-management/v3/users/%s/group-access", params.UserName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrAccessibleGroups, err)
	}

	var result ListAccessibleGroupsResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrAccessibleGroups, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrAccessibleGroups, i.Error(resp))
	}

	return result, nil
}
