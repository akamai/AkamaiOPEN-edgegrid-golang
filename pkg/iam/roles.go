package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Roles is the IAM role API interface
	Roles interface {
		// CreateRole creates a custom role
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/post-role
		CreateRole(context.Context, CreateRoleRequest) (*Role, error)

		// GetRole gets details for a specific role
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-role
		GetRole(context.Context, GetRoleRequest) (*Role, error)

		// UpdateRole adds or removes permissions from a role and updates other parameters
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/put-role
		UpdateRole(context.Context, UpdateRoleRequest) (*Role, error)

		// DeleteRole deletes a role. This operation is only allowed if the role isn't assigned to any users.
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/delete-role
		DeleteRole(context.Context, DeleteRoleRequest) error

		// ListRoles lists roles for the current account and contract type
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-roles
		ListRoles(context.Context, ListRolesRequest) ([]Role, error)

		// ListGrantableRoles lists which grantable roles can be included in a new custom role or added to an existing custom role
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-grantable-roles
		ListGrantableRoles(context.Context) ([]RoleGrantedRole, error)
	}

	// RoleRequest describes request parameters of the create and update role endpoint
	RoleRequest struct {
		Name         string          `json:"roleName,omitempty"`
		Description  string          `json:"roleDescription,omitempty"`
		GrantedRoles []GrantedRoleID `json:"grantedRoles,omitempty"`
	}

	// CreateRoleRequest describes the request parameters of the create role endpoint
	CreateRoleRequest RoleRequest

	// GrantedRoleID describes a unique identifier for a granted role
	GrantedRoleID struct {
		ID int64 `json:"grantedRoleId"`
	}

	// GetRoleRequest describes the request parameters of the get role endpoint
	GetRoleRequest struct {
		ID           int64
		Actions      bool
		GrantedRoles bool
		Users        bool
	}

	// UpdateRoleRequest describes the request parameters of the update role endpoint.
	// It works as patch request. You need to provide only fields which you want to update.
	UpdateRoleRequest struct {
		ID int64
		RoleRequest
	}

	// DeleteRoleRequest describes the request parameters of the delete role endpoint
	DeleteRoleRequest struct {
		ID int64
	}

	// ListRolesRequest describes the request parameters of the list roles endpoint
	ListRolesRequest struct {
		GroupID       *int64
		Actions       bool
		IgnoreContext bool
		Users         bool
	}

	// RoleAction encapsulates permissions available to the user for this role
	RoleAction struct {
		Delete bool `json:"delete"`
		Edit   bool `json:"edit"`
	}

	// RoleGrantedRole is a list of granted roles, giving the user access to objects in a group
	RoleGrantedRole struct {
		Description string `json:"grantedRoleDescription,omitempty"`
		RoleID      int64  `json:"grantedRoleId"`
		RoleName    string `json:"grantedRoleName"`
	}

	// RoleUser user who shares the same role
	RoleUser struct {
		AccountID     string `json:"accountId"`
		Email         string `json:"email"`
		FirstName     string `json:"firstName"`
		LastLoginDate string `json:"lastLoginDate"`
		LastName      string `json:"lastName"`
		UIIdentityID  string `json:"uiIdentityId"`
	}

	// Role encapsulates the response of the list roles endpoint
	Role struct {
		Actions         *RoleAction       `json:"actions,omitempty"`
		CreatedBy       string            `json:"createdBy"`
		CreatedDate     string            `json:"createdDate"`
		GrantedRoles    []RoleGrantedRole `json:"grantedRoles,omitempty"`
		ModifiedBy      string            `json:"modifiedBy"`
		ModifiedDate    string            `json:"modifiedDate"`
		RoleDescription string            `json:"roleDescription"`
		RoleID          int64             `json:"roleId"`
		RoleName        string            `json:"roleName"`
		Users           []RoleUser        `json:"users,omitempty"`
		RoleType        RoleType          `json:"type"`
	}

	// RoleType is an enum of role types
	RoleType string
)

// Validate validates CreateRoleRequest
func (r CreateRoleRequest) Validate() error {
	return validation.Errors{
		"Name":         validation.Validate(r.Name, validation.Required),
		"Description":  validation.Validate(r.Description, validation.Required),
		"GrantedRoles": validation.Validate(r.GrantedRoles, validation.Required),
	}.Filter()
}

// Validate validates GetRoleRequest
func (r GetRoleRequest) Validate() error {
	return validation.Errors{
		"ID": validation.Validate(r.ID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRoleRequest
func (r UpdateRoleRequest) Validate() error {
	return validation.Errors{
		"ID": validation.Validate(r.ID, validation.Required),
	}.Filter()
}

// Validate validates DeleteRoleRequest
func (r DeleteRoleRequest) Validate() error {
	return validation.Errors{
		"ID": validation.Validate(r.ID, validation.Required),
	}.Filter()
}

var (
	// RoleTypeStandard is a standard type provided by Akamai
	RoleTypeStandard RoleType = "standard"

	// RoleTypeCustom is a custom role provided by the account
	RoleTypeCustom RoleType = "custom"
)

var (
	// ErrCreateRole is returned when CreateRole fails
	ErrCreateRole = errors.New("create a role")
	// ErrGetRole is returned when GetRole fails
	ErrGetRole = errors.New("get a role")
	// ErrUpdateRole is returned when UpdateRole fails
	ErrUpdateRole = errors.New("update a role")
	// ErrDeleteRole is returned when DeleteRole fails
	ErrDeleteRole = errors.New("delete a role")
	// ErrListRoles is returned when ListRoles fails
	ErrListRoles = errors.New("list roles")
	// ErrListGrantableRoles is returned when ListGrantableRoles fails
	ErrListGrantableRoles = errors.New("list grantable roles")
)

func (i *iam) CreateRole(ctx context.Context, params CreateRoleRequest) (*Role, error) {
	logger := i.Log(ctx)
	logger.Debug("CreateRole")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateRole, ErrStructValidation, err)
	}

	uri, err := url.Parse("/identity-management/v2/user-admin/roles")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreateRole, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateRole, err)
	}

	var result Role
	resp, err := i.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateRole, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateRole, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) GetRole(ctx context.Context, params GetRoleRequest) (*Role, error) {
	logger := i.Log(ctx)
	logger.Debug("GetRole")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetRole, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/roles/%d", params.ID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetRole, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	q.Add("grantedRoles", strconv.FormatBool(params.GrantedRoles))
	q.Add("users", strconv.FormatBool(params.Users))

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetRole, err)
	}

	var result Role
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetRole, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetRole, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) UpdateRole(ctx context.Context, params UpdateRoleRequest) (*Role, error) {
	logger := i.Log(ctx)
	logger.Debug("UpdateRole")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateRole, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/roles/%d", params.ID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateRole, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateRole, err)
	}

	var result Role
	resp, err := i.Exec(req, &result, params.RoleRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateRole, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateRole, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) DeleteRole(ctx context.Context, params DeleteRoleRequest) error {
	logger := i.Log(ctx)
	logger.Debug("DeleteRole")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrDeleteRole, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/roles/%d", params.ID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeleteRole, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteRole, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteRole, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteRole, i.Error(resp))
	}

	return nil
}

func (i *iam) ListRoles(ctx context.Context, params ListRolesRequest) ([]Role, error) {
	logger := i.Log(ctx)
	logger.Debug("ListRoles")

	u, err := url.Parse("/identity-management/v2/user-admin/roles")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListRoles, err)
	}
	q := u.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	q.Add("ignoreContext", strconv.FormatBool(params.IgnoreContext))
	q.Add("users", strconv.FormatBool(params.Users))

	if params.GroupID != nil {
		q.Add("groupId", strconv.FormatInt(*params.GroupID, 10))
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListRoles, err)
	}

	var result []Role
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListRoles, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListRoles, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListGrantableRoles(ctx context.Context) ([]RoleGrantedRole, error) {
	logger := i.Log(ctx)
	logger.Debug("ListGrantableRoles")

	uri, err := url.Parse("/identity-management/v2/user-admin/roles/grantable-roles")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListGrantableRoles, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListGrantableRoles, err)
	}

	var result []RoleGrantedRole
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListGrantableRoles, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListGrantableRoles, i.Error(resp))
	}

	return result, nil
}
