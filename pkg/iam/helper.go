package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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
	}

	// ListAllowedCPCodesRequest contains the request parameter for the list of allowed CP codes endpoint
	ListAllowedCPCodesRequest struct {
		UserName string
		ListAllowedCPCodesRequestBody
	}

	// ListAllowedCPCodesRequestBody contains the filtering parameters for the list of allowed CP codes endpoint
	ListAllowedCPCodesRequestBody struct {
		ClientType string                `json:"clientType"`
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
)

var (
	// ErrListAllowedCPCodes is returned when ListAllowedCPCodes fails
	ErrListAllowedCPCodes = errors.New("list allowed CP codes")
)

// Validate validates ListAllowedCPCodesRequest
func (r ListAllowedCPCodesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"UserName": validation.Validate(r.UserName, validation.Required),
		"Body":     validation.Validate(r.ListAllowedCPCodesRequestBody, validation.Required),
	})
}

// Validate validates ListAllowedCPCodesRequestBody
func (r ListAllowedCPCodesRequestBody) Validate() error {
	return validation.Errors{
		"ClientType": validation.Validate(r.ClientType, validation.Required, validation.In("CLIENT", "USER_CLIENT", "SERVICE_ACCOUNT").Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'CLIENT' or 'USER_CLIENT' or 'SERVICE_ACCOUNT'", r.ClientType))),
		"Groups":     validation.Validate(r.Groups, validation.Required.When(r.ClientType == "SERVICE_ACCOUNT")),
	}.Filter()
}

func (i *iam) ListAllowedCPCodes(ctx context.Context, params ListAllowedCPCodesRequest) (ListAllowedCPCodesResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListAllowedCPCodes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListAllowedCPCodes, ErrStructValidation, err)
	}

	url := fmt.Sprintf("/identity-management/v3/users/%s/allowed-cpcodes", params.UserName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAllowedCPCodes, err)
	}

	var result ListAllowedCPCodesResponse
	resp, err := i.Exec(req, &result, params.ListAllowedCPCodesRequestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAllowedCPCodes, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAllowedCPCodes, i.Error(resp))
	}

	return result, nil
}
