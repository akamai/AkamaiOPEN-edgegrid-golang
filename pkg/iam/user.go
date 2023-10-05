package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	// Users is the IAM user identity API interface
	Users interface {
		// CreateUser creates a user in the account specified in your own API client credentials or clone an existing user's role assignments
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/post-ui-identity
		CreateUser(context.Context, CreateUserRequest) (*User, error)

		// GetUser gets  a specific user's profile
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-ui-identity
		GetUser(context.Context, GetUserRequest) (*User, error)

		// ListUsers returns a list of users who have access on this account
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-ui-identities
		ListUsers(context.Context, ListUsersRequest) ([]UserListItem, error)

		// RemoveUser removes a user identity
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/delete-ui-identity
		RemoveUser(context.Context, RemoveUserRequest) error

		// UpdateUserAuthGrants edits what groups a user has access to, and how the user can interact with the objects in those groups
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/put-ui-uiidentity-auth-grants
		UpdateUserAuthGrants(context.Context, UpdateUserAuthGrantsRequest) ([]AuthGrant, error)

		// UpdateUserInfo updates a user's information
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/put-ui-identity-basic-info
		UpdateUserInfo(context.Context, UpdateUserInfoRequest) (*UserBasicInfo, error)

		// UpdateUserNotifications subscribes or un-subscribe user to product notification emails
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/put-notifications
		UpdateUserNotifications(context.Context, UpdateUserNotificationsRequest) (*UserNotifications, error)

		// UpdateTFA updates a user's two-factor authentication setting and can reset tfa
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/put-ui-identity-tfa
		UpdateTFA(context.Context, UpdateTFARequest) error
	}

	// CreateUserRequest contains the request parameters for the create user endpoint
	CreateUserRequest struct {
		UserBasicInfo
		AuthGrants    []AuthGrantRequest `json:"authGrants,omitempty"`
		Notifications UserNotifications  `json:"notifications,omitempty"`
		SendEmail     bool               `json:"-"`
	}

	// ListUsersRequest contains the request parameters for the list users endpoint
	ListUsersRequest struct {
		GroupID    *int64
		AuthGrants bool
		Actions    bool
	}

	// GetUserRequest contains the request parameters of the get user endpoint
	GetUserRequest struct {
		IdentityID    string
		Actions       bool
		AuthGrants    bool
		Notifications bool
	}

	// UpdateUserInfoRequest contains the request parameters of the update user endpoint
	UpdateUserInfoRequest struct {
		IdentityID string
		User       UserBasicInfo
	}

	// UpdateUserNotificationsRequest contains the request parameters of the update user notifications endpoint
	UpdateUserNotificationsRequest struct {
		IdentityID    string
		Notifications UserNotifications
	}

	// UpdateUserAuthGrantsRequest contains the request parameters of the update user auth grants endpoint
	UpdateUserAuthGrantsRequest struct {
		IdentityID string
		AuthGrants []AuthGrantRequest
	}

	// RemoveUserRequest contains the request parameters of the remove user endpoint
	RemoveUserRequest struct {
		IdentityID string
	}

	// User describes the response of the get and create user endpoints
	User struct {
		UserBasicInfo
		IdentityID         string            `json:"uiIdentityId"`
		IsLocked           bool              `json:"isLocked"`
		LastLoginDate      string            `json:"lastLoginDate,omitempty"`
		PasswordExpiryDate string            `json:"passwordExpiryDate,omitempty"`
		TFAConfigured      bool              `json:"tfaConfigured"`
		EmailUpdatePending bool              `json:"emailUpdatePending"`
		AuthGrants         []AuthGrant       `json:"authGrants,omitempty"`
		Notifications      UserNotifications `json:"notifications,omitempty"`
	}

	// UserListItem describes the response of the list endpoint
	UserListItem struct {
		FirstName     string       `json:"firstName"`
		LastName      string       `json:"lastName"`
		UserName      string       `json:"uiUserName,omitempty"`
		Email         string       `json:"email"`
		TFAEnabled    bool         `json:"tfaEnabled"`
		IdentityID    string       `json:"uiIdentityId"`
		IsLocked      bool         `json:"isLocked"`
		LastLoginDate string       `json:"lastLoginDate,omitempty"`
		TFAConfigured bool         `json:"tfaConfigured"`
		AccountID     string       `json:"accountId"`
		Actions       *UserActions `json:"actions,omitempty"`
		AuthGrants    []AuthGrant  `json:"authGrants,omitempty"`
	}

	// UserBasicInfo is the user basic info structure
	UserBasicInfo struct {
		FirstName         string `json:"firstName"`
		LastName          string `json:"lastName"`
		UserName          string `json:"uiUserName,omitempty"`
		Email             string `json:"email"`
		Phone             string `json:"phone,omitempty"`
		TimeZone          string `json:"timeZone,omitempty"`
		JobTitle          string `json:"jobTitle"`
		TFAEnabled        bool   `json:"tfaEnabled"`
		SecondaryEmail    string `json:"secondaryEmail,omitempty"`
		MobilePhone       string `json:"mobilePhone,omitempty"`
		Address           string `json:"address,omitempty"`
		City              string `json:"city,omitempty"`
		State             string `json:"state,omitempty"`
		ZipCode           string `json:"zipCode,omitempty"`
		Country           string `json:"country"`
		ContactType       string `json:"contactType,omitempty"`
		PreferredLanguage string `json:"preferredLanguage,omitempty"`
		SessionTimeOut    *int   `json:"sessionTimeOut,omitempty"`
	}

	// UserActions encapsulates permissions available to the user for this group
	UserActions struct {
		APIClient        bool `json:"apiClient"`
		Delete           bool `json:"delete"`
		Edit             bool `json:"edit"`
		IsCloneable      bool `json:"isCloneable"`
		ResetPassword    bool `json:"resetPassword"`
		ThirdPartyAccess bool `json:"thirdPartyAccess"`
		CanEditTFA       bool `json:"canEditTFA"`
		EditProfile      bool `json:"editProfile"`
	}

	// AuthGrant is user’s role assignments, per group
	AuthGrant struct {
		GroupID         int64       `json:"groupId"`
		GroupName       string      `json:"groupName"`
		IsBlocked       bool        `json:"isBlocked"`
		RoleDescription string      `json:"roleDescription"`
		RoleID          *int        `json:"roleId,omitempty"`
		RoleName        string      `json:"roleName"`
		Subgroups       []AuthGrant `json:"subGroups,omitempty"`
	}

	// AuthGrantRequest is user’s role assignments, per group for the create/update operation
	AuthGrantRequest struct {
		GroupID   int64              `json:"groupId"`
		IsBlocked bool               `json:"isBlocked"`
		RoleID    *int               `json:"roleId,omitempty"`
		Subgroups []AuthGrantRequest `json:"subGroups,omitempty"`
	}

	// UserNotifications types of notification emails the user receives
	UserNotifications struct {
		EnableEmail bool                    `json:"enableEmailNotifications"`
		Options     UserNotificationOptions `json:"options"`
	}

	// UserNotificationOptions types of notification emails the user receives
	UserNotificationOptions struct {
		NewUser        bool     `json:"newUserNotification"`
		PasswordExpiry bool     `json:"passwordExpiry"`
		Proactive      []string `json:"proactive"`
		Upgrade        []string `json:"upgrade"`
	}

	// TFAActionType is a type for tfa action constants
	TFAActionType string

	// UpdateTFARequest contains the request parameters of the tfa user endpoint
	UpdateTFARequest struct {
		IdentityID string
		Action     TFAActionType
	}
)

const (
	// TFAActionEnable ia an action value to use to enable tfa
	TFAActionEnable TFAActionType = "enable"
	// TFAActionDisable ia an action value to use to disable tfa
	TFAActionDisable TFAActionType = "disable"
	// TFAActionReset ia an action value to use to reset tfa
	TFAActionReset TFAActionType = "reset"
)

var (
	// ErrCreateUser is returned when CreateUser fails
	ErrCreateUser = errors.New("create user")

	// ErrGetUser is returned when GetUser fails
	ErrGetUser = errors.New("get user")

	// ErrListUsers is returned when GetUser fails
	ErrListUsers = errors.New("list users")

	// ErrRemoveUser is returned when RemoveUser fails
	ErrRemoveUser = errors.New("remove user")

	// ErrUpdateUserAuthGrants is returned when UpdateUserAuthGrants fails
	ErrUpdateUserAuthGrants = errors.New("update user auth grants")

	// ErrUpdateUserInfo is returned when UpdateUserInfo fails
	ErrUpdateUserInfo = errors.New("update user info")

	// ErrUpdateUserNotifications is returned when UpdateUserNotifications fails
	ErrUpdateUserNotifications = errors.New("update user notifications")

	// ErrUpdateTFA is returned when UpdateTFA fails
	ErrUpdateTFA = errors.New("update user's two-factor authentication")
)

// Validate performs validation on AuthGrant
func (r AuthGrant) Validate() error {
	return validation.Errors{
		"group_id": validation.Validate(r.GroupID, validation.Required),
		"role_id":  validation.Validate(r.RoleID, validation.Required),
	}.Filter()
}

// Validate validates CreateUserRequest
func (r CreateUserRequest) Validate() error {
	return validation.Errors{
		"country":       validation.Validate(r.Country, validation.Required),
		"email":         validation.Validate(r.Email, validation.Required, is.EmailFormat),
		"firstName":     validation.Validate(r.FirstName, validation.Required),
		"lastName":      validation.Validate(r.LastName, validation.Required),
		"authGrants":    validation.Validate(r.AuthGrants, validation.Required),
		"notifications": validation.Validate(r.Notifications, validation.Required),
	}.Filter()
}

// Validate validates GetUserRequest
func (r GetUserRequest) Validate() error {
	return validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

// Validate validates UpdateUserInfoRequest
func (r UpdateUserInfoRequest) Validate() error {
	return validation.Errors{
		"uiIdentity":        validation.Validate(r.IdentityID, validation.Required),
		"firstName":         validation.Validate(r.User.FirstName, validation.Required),
		"lastName":          validation.Validate(r.User.LastName, validation.Required),
		"country":           validation.Validate(r.User.Country, validation.Required),
		"timeZone":          validation.Validate(r.User.TimeZone, validation.Required),
		"preferredLanguage": validation.Validate(r.User.PreferredLanguage, validation.Required),
		"sessionTimeOut":    validation.Validate(r.User.SessionTimeOut, validation.Required),
	}.Filter()
}

// Validate validates UpdateUserNotificationsRequest
func (r UpdateUserNotificationsRequest) Validate() error {
	return validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

// Validate validates UpdateUserAuthGrantsRequest
func (r UpdateUserAuthGrantsRequest) Validate() error {
	return validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
		"authGrants": validation.Validate(r.AuthGrants, validation.Required),
	}.Filter()
}

// Validate validates RemoveUserRequest
func (r RemoveUserRequest) Validate() error {
	return validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

// Validate validates UpdateTFARequest
func (r UpdateTFARequest) Validate() error {
	return validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
		"Action": validation.Validate(r.Action, validation.Required, validation.In(TFAActionEnable, TFAActionDisable, TFAActionReset).
			Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'enable', 'disable' or 'reset'", r.Action))),
	}.Filter()
}

func (i *iam) CreateUser(ctx context.Context, params CreateUserRequest) (*User, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateUser, ErrStructValidation, err)
	}

	u, err := url.Parse(path.Join("/identity-management/v2/user-admin", "ui-identities"))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateUser, err)
	}

	q := u.Query()
	q.Add("sendEmail", strconv.FormatBool(params.SendEmail))

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateUser, err)
	}

	var result User
	resp, err := i.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateUser, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateUser, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) GetUser(ctx context.Context, params GetUserRequest) (*User, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetUser, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetUser, err)
	}

	q := u.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	q.Add("authGrants", strconv.FormatBool(params.AuthGrants))
	q.Add("notifications", strconv.FormatBool(params.Notifications))

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetUser, err)
	}

	var rval User
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetUser, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetUser, i.Error(resp))
	}

	return &rval, nil
}

func (i *iam) ListUsers(ctx context.Context, params ListUsersRequest) ([]UserListItem, error) {
	u, err := url.Parse("/identity-management/v2/user-admin/ui-identities")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse the URL:\n%s", ErrListUsers, err)
	}

	q := u.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	q.Add("authGrants", strconv.FormatBool(params.AuthGrants))
	if params.GroupID != nil {
		q.Add("groupId", strconv.FormatInt(int64(*params.GroupID), 10))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request:\n%s", ErrListUsers, err)
	}

	var users []UserListItem
	resp, err := i.Exec(req, &users)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed:\n%s", ErrListUsers, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListUsers, i.Error(resp))
	}

	return users, nil
}

func (i *iam) RemoveUser(ctx context.Context, params RemoveUserRequest) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrRemoveUser, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrRemoveUser, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrRemoveUser, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrRemoveUser, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrRemoveUser, i.Error(resp))
	}

	return nil
}

func (i *iam) UpdateUserAuthGrants(ctx context.Context, params UpdateUserAuthGrantsRequest) ([]AuthGrant, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateUserAuthGrants, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/auth-grants", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserAuthGrants, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserAuthGrants, err)
	}

	rval := make([]AuthGrant, 0)

	resp, err := i.Exec(req, &rval, params.AuthGrants)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateUserAuthGrants, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateUserAuthGrants, i.Error(resp))
	}

	return rval, nil
}

func (i *iam) UpdateUserInfo(ctx context.Context, params UpdateUserInfoRequest) (*UserBasicInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateUserInfo, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/basic-info", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserInfo, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserInfo, err)
	}

	var rval UserBasicInfo
	resp, err := i.Exec(req, &rval, params.User)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateUserInfo, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateUserInfo, i.Error(resp))
	}

	return &rval, nil
}

func (i *iam) UpdateUserNotifications(ctx context.Context, params UpdateUserNotificationsRequest) (*UserNotifications, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateUserNotifications, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/notifications", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserNotifications, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserNotifications, err)
	}

	var rval UserNotifications
	resp, err := i.Exec(req, &rval, params.Notifications)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateUserNotifications, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateUserNotifications, i.Error(resp))
	}

	return &rval, nil
}

func (i *iam) UpdateTFA(ctx context.Context, params UpdateTFARequest) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrUpdateTFA, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/tfa", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUpdateTFA, err)
	}

	q := u.Query()
	q.Add("action", string(params.Action))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUpdateTFA, err)
	}

	resp, err := i.Exec(req, nil, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrUpdateTFA, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrUpdateTFA, i.Error(resp))
	}

	return nil
}
