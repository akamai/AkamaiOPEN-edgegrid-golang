package iam

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/spf13/cast"
)

type (
	// Users is the IAM user identity management interface
	Users interface {
		CreateUser(context.Context, CreateUserRequest) (*User, error)
		GetUser(context.Context, GetUserRequest) (*User, error)
		UpdateUserInfo(context.Context, UpdateUserInfoRequest) (*UserBasicInfo, error)
		UpdateUserNotifications(context.Context, UpdateUserNotificationsRequest) (*UserNotifications, error)
		UpdateUserAuthGrants(context.Context, UpdateUserAuthGrantsRequest) ([]AuthGrant, error)
		RemoveUser(context.Context, RemoveUserRequest) error
	}

	// CreateUserRequest is the input to CreateUser
	CreateUserRequest struct {
		User          UserBasicInfo      `json:"user"`
		Notifications *UserNotifications `json:"notifications,omitempty"`
		AuthGrants    []AuthGrant        `json:"authGrants,omitempty"`
		SendEmail     bool               `json:"sendEmail"`
	}

	// GetUserRequest is the input for GetUser
	GetUserRequest struct {
		IdentityID    string `json:"uiIdentityId"`
		Actions       bool   `json:"actions"`
		AuthGrants    bool   `json:"authGrants"`
		Notifications bool   `json:"notificiations"`
	}

	// UpdateUserInfoRequest is the input to UpdateUserInfo
	UpdateUserInfoRequest struct {
		IdentityID string        `json:"uiIdentityId"`
		User       UserBasicInfo `json:"user"`
	}

	// UpdateUserNotificationsRequest is the input to update user notifications
	UpdateUserNotificationsRequest struct {
		IdentityID    string            `json:"uiIdentityId"`
		Notifications UserNotifications `json:"notifications,omitempty"`
	}

	// UpdateUserAuthGrantsRequest is the input to update user auth grants
	UpdateUserAuthGrantsRequest struct {
		IdentityID string      `json:"uiIdentityId"`
		AuthGrants []AuthGrant `json:"authGrants,omitempty"`
	}

	// RemoveUserRequest is the input for RemoveUser
	RemoveUserRequest struct {
		IdentityID string `json:"uiIdentityId"`
	}

	// User encapsulates information about each user.
	User struct {
		UserBasicInfo
		IdentityID         string             `json:"uiIdentityId"`
		IsLocked           bool               `json:"isLocked"`
		LastLoginDate      string             `json:"lastLoginDate,omitempty"`
		PasswordExpiryDate string             `json:"passwordExpiryDate,omitempty"`
		TFAConfigured      bool               `json:"tfaConfigured"`
		EmailUpdatePending bool               `json:"emailUpdatePending"`
		AuthGrants         []AuthGrant        `json:"authGrants,omitempty"`
		Notifications      *UserNotifications `json:"notifications,omitempty"`
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
		ContactType       string `json:"contactType"`
		PreferredLanguage string `json:"preferredLanguage,omitempty"`
		SessionTimeOut    *int   `json:"sessionTimeOut,omitempty"`
	}

	// UserActions encapsulates permissions available to the user for this group.
	UserActions struct {
		APIClient        bool `json:"apiClient"`
		Delete           bool `json:"delete"`
		Edit             bool `json:"edit"`
		IsCloneable      bool `json:"isCloneable"`
		ResetPassword    bool `json:"resetPassword"`
		ThirdPartyAccess bool `json:"thirdPartyAccess"`
	}

	// AuthGrant is user’s role assignments, per group.
	AuthGrant struct {
		GroupID         int         `json:"groupId"`
		GroupName       string      `json:"groupName"`
		IsBlocked       bool        `json:"isBlocked"`
		RoleDescription string      `json:"roleDescription"`
		RoleID          *int        `json:"roleId,omitempty"`
		RoleName        string      `json:"roleName"`
		Subgroups       []AuthGrant `json:"subGroups,omitempty"`
	}

	// UserNotifications types of notification emails the user receives.
	UserNotifications struct {
		EnableEmail bool                     `json:"enableEmailNotifications"`
		Options     *UserNotificationOptions `json:"options,omitempty"`
	}

	// UserNotificationOptions types of notification emails the user receives.
	UserNotificationOptions struct {
		NewUser        bool     `json:"newUserNotification"`
		PasswordExpiry bool     `json:"passwordExpiry"`
		Proactive      []string `json:"proactive,omitempty"`
		Upgrade        []string `json:"upgrade,omitempty"`
	}
)

// Validate performs the input validation for CreateUserRequest
func (r CreateUserRequest) Validate() error {
	return validation.Errors{
		"country":    validation.Validate(r.User.Country, validation.Required),
		"email":      validation.Validate(r.User.Email, validation.Required, is.EmailFormat),
		"firstName":  validation.Validate(r.User.FirstName, validation.Required),
		"lastName":   validation.Validate(r.User.LastName, validation.Required),
		"phone":      validation.Validate(r.User.Phone, validation.Required),
		"authGrants": validation.Validate(r.AuthGrants, validation.NilOrNotEmpty),
	}.Filter()
}

// Validate performs the input validation for GetUserRequest
func (r GetUserRequest) Validate() error {
	return validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

// Validate performs the input validation for UpdateUserRequest
func (r UpdateUserInfoRequest) Validate() error {
	return validation.Errors{
		"uiIdentity":        validation.Validate(r.IdentityID, validation.Required),
		"firstName":         validation.Validate(r.User.FirstName, validation.Required),
		"lastName":          validation.Validate(r.User.LastName, validation.Required),
		"country":           validation.Validate(r.User.Country, validation.Required),
		"phone":             validation.Validate(r.User.Phone, validation.Required),
		"contactType":       validation.Validate(r.User.ContactType, validation.Required),
		"timeZone":          validation.Validate(r.User.TimeZone, validation.Required),
		"preferredLanguage": validation.Validate(r.User.PreferredLanguage, validation.Required),
		"sessionTimeOut":    validation.Validate(r.User.SessionTimeOut, validation.Required),
	}.Filter()
}

// Validate performs the input validation for UpdateUserNotificationsRequest
func (r UpdateUserNotificationsRequest) Validate() error {
	return validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

// Validate performs the input validation for UpdateUserAuthGrantsRequest
func (r UpdateUserAuthGrantsRequest) Validate() error {
	return validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
		"authGrants": validation.Validate(r.AuthGrants, validation.Required),
	}.Filter()
}

// Validate performs the input validation for RemoveUserRequest
func (r RemoveUserRequest) Validate() error {
	return validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

// CreateUser creates a new iam user
func (i *iam) CreateUser(ctx context.Context, params CreateUserRequest) (*User, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInputValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "ui-identities"))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "CreateUser", err)
	}

	q := u.Query()
	q.Add("sendEmail", cast.ToString(params.SendEmail))

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "CreateUser", err)
	}

	user := User{
		UserBasicInfo: params.User,
		AuthGrants:    params.AuthGrants,
		Notifications: params.Notifications,
	}

	var rval User
	resp, err := i.Exec(req, &rval, user)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "CreateUser", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", "CreateUser", i.Error(resp))
	}

	return &rval, nil
}

// GetUser gets a user by id
func (i *iam) GetUser(ctx context.Context, params GetUserRequest) (*User, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInputValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "ui-identities", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "GetUser", err)
	}

	q := u.Query()
	q.Add("actions", cast.ToString(params.Actions))
	q.Add("authGrants", cast.ToString(params.AuthGrants))
	q.Add("notifications", cast.ToString(params.Notifications))

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "GetUser", err)
	}

	var rval User
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "GetUser", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "GetUser", i.Error(resp))
	}

	return &rval, nil
}

// UpdateUserInfo updates a user's basic info
func (i *iam) UpdateUserInfo(ctx context.Context, params UpdateUserInfoRequest) (*UserBasicInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInputValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "ui-identities", params.IdentityID, "basic-info"))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "UpdateUser", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "UpdateUser", err)
	}

	var rval UserBasicInfo
	resp, err := i.Exec(req, &rval, params.User)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "UpdateUser", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "UpdateUser", i.Error(resp))
	}

	return &rval, nil
}

// UpdateUserNotifications updates a user's notifications
func (i *iam) UpdateUserNotifications(ctx context.Context, params UpdateUserNotificationsRequest) (*UserNotifications, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInputValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "ui-identities", params.IdentityID, "notifications"))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "UpdateUserNotifications", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "UpdateUserNotifications", err)
	}

	var rval UserNotifications
	resp, err := i.Exec(req, &rval, params.Notifications)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "UpdateUserNotifications", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "UpdateUserNotifications", i.Error(resp))
	}

	return &rval, nil
}

// UpdateUserAuthGrants updates a user's notifications
func (i *iam) UpdateUserAuthGrants(ctx context.Context, params UpdateUserAuthGrantsRequest) ([]AuthGrant, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInputValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "ui-identities", params.IdentityID, "auth-grants"))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "UpdateUserAuthGrants", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "UpdateUserAuthGrants", err)
	}

	rval := make([]AuthGrant, 0)

	resp, err := i.Exec(req, &rval, params.AuthGrants)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "UpdateUserAuthGrants", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "UpdateUserAuthGrants", i.Error(resp))
	}

	return rval, nil
}

// RemoveUser removes a user identity
func (i *iam) RemoveUser(ctx context.Context, params RemoveUserRequest) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrInputValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "ui-identities", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%s: failed to create request: %s", "RemoveUser", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf("%s: failed to create request: %s", "RemoveUser", err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%s: request failed: %s", "RemoveUser", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", "RemoveUser", i.Error(resp))
	}

	return nil
}
