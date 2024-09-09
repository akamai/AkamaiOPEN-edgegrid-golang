package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Helper is a list of IAM helper API interfaces
	Helper interface {
		// ListAllowedCPCodes lists available CP codes for a user
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-api-clients-users-allowed-cpcodes
		ListAllowedCPCodes(context.Context, ListAllowedCPCodesRequest) (ListAllowedCPCodesResponse, error)

		// ListAuthorizedUsers lists authorized API client users
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-clients-users
		ListAuthorizedUsers(context.Context) (ListAuthorizedUsersResponse, error)

		// ListAllowedAPIs lists available APIs for a user
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-clients-users-allowed-apis
		ListAllowedAPIs(context.Context, ListAllowedAPIsRequest) (ListAllowedAPIsResponse, error)

		// ListAccessibleGroups lists groups available to a user
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-clients-users-group-access
		ListAccessibleGroups(context.Context, ListAccessibleGroupsRequest) (ListAccessibleGroupsResponse, error)
	}

	// ListAllowedCPCodesRequest contains the request parameter for the list of allowed CP codes endpoint
	ListAllowedCPCodesRequest struct {
		UserName string
		Body     ListAllowedCPCodesRequestBody
	}

	// ListAllowedAPIsRequest contains the request parameters for the list of allowed APIs endpoint
	ListAllowedAPIsRequest struct {
		UserName           string
		ClientType         ClientType
		AllowAccountSwitch bool
	}

	// ListAccessibleGroupsRequest contains the request parameter for the list of accessible groups endpoint
	ListAccessibleGroupsRequest struct {
		UserName string
	}

	// ListAllowedCPCodesRequestBody contains the filtering parameters for the list of allowed CP codes endpoint
	ListAllowedCPCodesRequestBody struct {
		ClientType ClientType            `json:"clientType"`
		Groups     []AllowedCPCodesGroup `json:"groups"`
	}

	// AllowedCPCodesGroup contains the group parameters for the list of allowed CP codes endpoint
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

	// ListAllowedCPCodesResponse contains response for the list of allowed CP codes endpoint
	ListAllowedCPCodesResponse []ListAllowedCPCodesResponseItem

	// ListAllowedCPCodesResponseItem contains single item of the response for allowed CP codes endpoint
	ListAllowedCPCodesResponseItem struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	// ListAuthorizedUsersResponse contains the response for the list of authorized users endpoint
	ListAuthorizedUsersResponse []AuthorizedUser

	// AuthorizedUser contains the details about the authorized user
	AuthorizedUser struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		UIIdentityID string `json:"uiIdentityId"`
	}

	// ListAccessibleGroupsResponse contains the response for the list of accessible groups endpoint
	ListAccessibleGroupsResponse []AccessibleGroup

	// AccessibleGroup contains the details about accessible group
	AccessibleGroup struct {
		GroupID         int64                `json:"groupId"`
		RoleID          int64                `json:"roleId"`
		GroupName       string               `json:"groupName"`
		RoleName        string               `json:"roleName"`
		IsBlocked       bool                 `json:"isBlocked"`
		RoleDescription string               `json:"roleDescription"`
		SubGroups       []AccessibleSubGroup `json:"subGroups"`
	}

	// AccessibleSubGroup contains the details about subgroup
	AccessibleSubGroup struct {
		GroupID       int64                `json:"groupId"`
		GroupName     string               `json:"groupName"`
		ParentGroupID int64                `json:"parentGroupId"`
		SubGroups     []AccessibleSubGroup `json:"subGroups"`
	}

	// ListAllowedAPIsResponse contains the response for the list of allowed APIs endpoint
	ListAllowedAPIsResponse []AllowedAPI

	// AllowedAPI contains the details about the API
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

	// ClientType represents the type of the client
	ClientType string
)

const (
	// UserClientType is the `USER_CLIENT` client type
	UserClientType ClientType = "USER_CLIENT"
	// ServiceAccountClientType is the `SERVICE_ACCOUNT` client type
	ServiceAccountClientType ClientType = "SERVICE_ACCOUNT"
	// ClientClientType is the `CLIENT` client type
	ClientClientType ClientType = "CLIENT"
)

var (
	// ErrListAllowedCPCodes is returned when ListAllowedCPCodes fails
	ErrListAllowedCPCodes = errors.New("list allowed CP codes")
	// ErrListAuthorizedUsers is returned when ListAuthorizedUsers fails
	ErrListAuthorizedUsers = errors.New("list authorized users")
	// ErrListAllowedAPIs is returned when ListAllowedAPIs fails
	ErrListAllowedAPIs = errors.New("list allowed APIs")
	// ErrAccessibleGroups is returned when ListAccessibleGroups fails
	ErrAccessibleGroups = errors.New("list accessible groups")
)

// Validate validates ListAllowedCPCodesRequest
func (r ListAllowedCPCodesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"UserName": validation.Validate(r.UserName, validation.Required),
		"Body":     validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates ListAllowedCPCodesRequestBody
func (r ListAllowedCPCodesRequestBody) Validate() error {
	return validation.Errors{
		"ClientType": validation.Validate(r.ClientType, validation.Required, validation.In(ClientClientType, UserClientType, ServiceAccountClientType).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'CLIENT' or 'USER_CLIENT' or 'SERVICE_ACCOUNT'", r.ClientType))),
		"Groups":     validation.Validate(r.Groups, validation.Required.When(r.ClientType == ServiceAccountClientType)),
	}.Filter()
}

// Validate validates ListAllowedAPIsRequest
func (r ListAllowedAPIsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"UserName":   validation.Validate(r.UserName, validation.Required),
		"ClientType": validation.Validate(r.ClientType, validation.In(ClientClientType, UserClientType, ServiceAccountClientType).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'CLIENT' or 'USER_CLIENT' or 'SERVICE_ACCOUNT'", r.ClientType))),
	})
}

// Validate validates ListAccessibleGroupsRequest
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

	u := fmt.Sprintf("/identity-management/v3/users/%s/allowed-cpcodes", params.UserName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAllowedCPCodes, err)
	}

	var result ListAllowedCPCodesResponse
	resp, err := i.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAllowedCPCodes, err)
	}

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

	u, err := url.Parse(fmt.Sprintf("/identity-management/v3/users/%s/allowed-apis", params.UserName))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAllowedAPIs, err)
	}

	q := u.Query()
	if params.ClientType != "" {
		q.Add("clientType", string(params.ClientType))

	}
	q.Add("allowAccountSwitch", strconv.FormatBool(params.AllowAccountSwitch))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAllowedAPIs, err)
	}

	var result ListAllowedAPIsResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAllowedAPIs, err)
	}

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

	u := fmt.Sprintf("/identity-management/v3/users/%s/group-access", params.UserName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrAccessibleGroups, err)
	}

	var result ListAccessibleGroupsResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrAccessibleGroups, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrAccessibleGroups, i.Error(resp))
	}

	return result, nil
}
