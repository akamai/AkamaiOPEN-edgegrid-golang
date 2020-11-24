package iam

import "context"

type (
	// Roles is the iam role management interface
	Roles interface {
		ListRoles(context.Context, ListRolesRequest) ([]Role, error)
	}

	// ListRolesRequest is option query parameters for the list roles endpoint
	ListRolesRequest struct {
		GroupID       *int64 `json:"groupId,omitempty"`
		Actions       bool   `json:"actions"`
		IgnoreContext bool   `json:"ignoreContext"`
		Users         bool   `json:"users"`
	}

	// RoleAction encapsulates permissions available to the user for this role
	RoleAction struct {
		Delete bool `json:"delete"`
		Edit   bool `json:"edit"`
	}

	// RoleGrantedRole is a list of granted roles, giving the user access to objects in a group.
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

	// Role is a role that includes granted roles.
	Role struct {
		Actions         *RoleAction       `json:"actions,omitempty"`
		CreatedBy       string            `json:"createdBy"`
		CreatedDate     striong           `json:"createdDate"`
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

var (
	// RoleTypeStandard is a standard type provided by Akamai
	RoleTypeStandard RoleType = "standard"

	// RoleTypeCustom is a custom role provided by the account
	RoleTypeCustom RoleType = "custom"
)
