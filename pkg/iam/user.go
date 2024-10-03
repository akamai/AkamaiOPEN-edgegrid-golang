package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	// CreateUserRequest contains the request parameters for the CreateUser endpoint.
	CreateUserRequest struct {
		UserBasicInfo
		AuthGrants    []AuthGrantRequest `json:"authGrants,omitempty"`
		Notifications *UserNotifications `json:"notifications,omitempty"`
		SendEmail     bool               `json:"-"`
	}

	// ListUsersRequest contains the request parameters for the ListUsers endpoint.
	ListUsersRequest struct {
		GroupID    *int64
		AuthGrants bool
		Actions    bool
	}

	// GetUserRequest contains the request parameters for the GetUser endpoint.
	GetUserRequest struct {
		IdentityID    string
		Actions       bool
		AuthGrants    bool
		Notifications bool
	}

	// UpdateUserInfoRequest contains the request parameters for the UpdateUser endpoint.
	UpdateUserInfoRequest struct {
		IdentityID string
		User       UserBasicInfo
	}

	// UpdateUserNotificationsRequest contains the request parameters for the UpdateUserNotifications endpoint.
	UpdateUserNotificationsRequest struct {
		IdentityID    string
		Notifications *UserNotifications
	}

	// UpdateUserAuthGrantsRequest contains the request parameters for the UpdateUserAuthGrants endpoint.
	UpdateUserAuthGrantsRequest struct {
		IdentityID string
		AuthGrants []AuthGrantRequest
	}

	// RemoveUserRequest contains the request parameters of the RemoveUser endpoint.
	RemoveUserRequest struct {
		IdentityID string
	}

	// User describes the response from the GetUser and CreateUser endpoints.
	User struct {
		UserBasicInfo
		IdentityID                         string            `json:"uiIdentityId"`
		IsLocked                           bool              `json:"isLocked"`
		LastLoginDate                      time.Time         `json:"lastLoginDate,omitempty"`
		PasswordExpiryDate                 time.Time         `json:"passwordExpiryDate,omitempty"`
		TFAConfigured                      bool              `json:"tfaConfigured"`
		EmailUpdatePending                 bool              `json:"emailUpdatePending"`
		AuthGrants                         []AuthGrant       `json:"authGrants,omitempty"`
		Notifications                      UserNotifications `json:"notifications,omitempty"`
		Actions                            *UserActions      `json:"actions,omitempty"`
		UserStatus                         string            `json:"userStatus"`
		AccountID                          string            `json:"accountId"`
		AdditionalAuthenticationConfigured bool              `json:"additionalAuthenticationConfigured"`
	}

	// UserListItem describes the response from the ListUsers endpoint.
	UserListItem struct {
		FirstName                          string         `json:"firstName"`
		LastName                           string         `json:"lastName"`
		UserName                           string         `json:"uiUserName,omitempty"`
		Email                              string         `json:"email"`
		TFAEnabled                         bool           `json:"tfaEnabled"`
		IdentityID                         string         `json:"uiIdentityId"`
		IsLocked                           bool           `json:"isLocked"`
		LastLoginDate                      time.Time      `json:"lastLoginDate,omitempty"`
		TFAConfigured                      bool           `json:"tfaConfigured"`
		AccountID                          string         `json:"accountId"`
		Actions                            *UserActions   `json:"actions,omitempty"`
		AuthGrants                         []AuthGrant    `json:"authGrants,omitempty"`
		AdditionalAuthentication           Authentication `json:"additionalAuthentication"`
		AdditionalAuthenticationConfigured bool           `json:"additionalAuthenticationConfigured"`
	}

	// UserBasicInfo is the user basic info structure.
	UserBasicInfo struct {
		FirstName                string         `json:"firstName"`
		LastName                 string         `json:"lastName"`
		UserName                 string         `json:"uiUserName,omitempty"`
		Email                    string         `json:"email"`
		Phone                    string         `json:"phone,omitempty"`
		TimeZone                 string         `json:"timeZone,omitempty"`
		JobTitle                 string         `json:"jobTitle"`
		TFAEnabled               bool           `json:"tfaEnabled"`
		SecondaryEmail           string         `json:"secondaryEmail,omitempty"`
		MobilePhone              string         `json:"mobilePhone,omitempty"`
		Address                  string         `json:"address,omitempty"`
		City                     string         `json:"city,omitempty"`
		State                    string         `json:"state,omitempty"`
		ZipCode                  string         `json:"zipCode,omitempty"`
		Country                  string         `json:"country"`
		ContactType              string         `json:"contactType,omitempty"`
		PreferredLanguage        string         `json:"preferredLanguage,omitempty"`
		SessionTimeOut           *int           `json:"sessionTimeOut,omitempty"`
		AdditionalAuthentication Authentication `json:"additionalAuthentication"`
	}

	// UserActions encapsulates permissions available to the user for this group.
	UserActions struct {
		APIClient             bool `json:"apiClient"`
		Delete                bool `json:"delete"`
		Edit                  bool `json:"edit"`
		IsCloneable           bool `json:"isCloneable"`
		ResetPassword         bool `json:"resetPassword"`
		ThirdPartyAccess      bool `json:"thirdPartyAccess"`
		CanEditTFA            bool `json:"canEditTFA"`
		CanEditMFA            bool `json:"canEditMFA"`
		CanEditNone           bool `json:"canEditNone"`
		EditProfile           bool `json:"editProfile"`
		EditRole              bool `json:"editRole"`
		CanGenerateBypassCode bool `json:"canGenerateBypassCode"`
	}

	// AuthGrant is user’s role assignments, per group.
	AuthGrant struct {
		GroupID         int64       `json:"groupId"`
		GroupName       string      `json:"groupName"`
		IsBlocked       bool        `json:"isBlocked"`
		RoleDescription string      `json:"roleDescription"`
		RoleID          *int        `json:"roleId,omitempty"`
		RoleName        string      `json:"roleName"`
		Subgroups       []AuthGrant `json:"subGroups,omitempty"`
	}

	// AuthGrantRequest is user’s role assignments, per group for the create/update operation.
	AuthGrantRequest struct {
		GroupID   int64              `json:"groupId"`
		IsBlocked bool               `json:"isBlocked"`
		RoleID    *int               `json:"roleId,omitempty"`
		Subgroups []AuthGrantRequest `json:"subGroups,omitempty"`
	}

	// UserNotifications types of notification emails the user receives.
	UserNotifications struct {
		EnableEmail bool                    `json:"enableEmailNotifications"`
		Options     UserNotificationOptions `json:"options"`
	}

	// UserNotificationOptions types of notification emails the user receives.
	UserNotificationOptions struct {
		NewUser                   bool     `json:"newUserNotification"`
		PasswordExpiry            bool     `json:"passwordExpiry"`
		Proactive                 []string `json:"proactive"`
		Upgrade                   []string `json:"upgrade"`
		APIClientCredentialExpiry bool     `json:"apiClientCredentialExpiryNotification"`
	}

	// TFAActionType is a type for tfa action constants.
	TFAActionType string

	// UpdateTFARequest contains the request parameters for the UpdateTFA endpoint.
	UpdateTFARequest struct {
		IdentityID string
		Action     TFAActionType
	}

	// Authentication is a type of additional authentication.
	Authentication string

	// UpdateMFARequest contains the request body for the UpdateMFA endpoint.
	UpdateMFARequest struct {
		IdentityID string
		Value      Authentication
	}

	// ResetMFARequest contains the request parameters for the ResetMFA endpoint.
	ResetMFARequest struct {
		IdentityID string
	}
)

const (
	// TFAActionEnable is an action value to use to enable tfa.
	TFAActionEnable TFAActionType = "enable"
	// TFAActionDisable is an action value to use to disable tfa.
	TFAActionDisable TFAActionType = "disable"
	// TFAActionReset is an action value to use to reset tfa.
	TFAActionReset TFAActionType = "reset"
	// MFAAuthentication is authentication of type MFA.
	MFAAuthentication Authentication = "MFA"
	// TFAAuthentication is authentication of type TFA.
	TFAAuthentication Authentication = "TFA"
	// NoneAuthentication represents a state where no authentication method is configured.
	NoneAuthentication Authentication = "NONE"
)

var (
	// ErrCreateUser is returned when CreateUser fails.
	ErrCreateUser = errors.New("create user")

	// ErrGetUser is returned when GetUser fails.
	ErrGetUser = errors.New("get user")

	// ErrListUsers is returned when ListUsers fails.
	ErrListUsers = errors.New("list users")

	// ErrRemoveUser is returned when RemoveUser fails.
	ErrRemoveUser = errors.New("remove user")

	// ErrUpdateUserAuthGrants is returned when UpdateUserAuthGrants fails.
	ErrUpdateUserAuthGrants = errors.New("update user auth grants")

	// ErrUpdateUserInfo is returned when UpdateUserInfo fails.
	ErrUpdateUserInfo = errors.New("update user info")

	// ErrUpdateUserNotifications is returned when UpdateUserNotifications fails.
	ErrUpdateUserNotifications = errors.New("update user notifications")

	// ErrUpdateTFA is returned when UpdateTFA fails.
	ErrUpdateTFA = errors.New("update user's two-factor authentication")

	// ErrUpdateMFA is returned when UpdateMFA fails.
	ErrUpdateMFA = errors.New("update user's authentication method")

	// ErrResetMFA is returned when ResetMFA fails.
	ErrResetMFA = errors.New("reset user's authentication method")
)

// Validate validates Authentication.
func (a Authentication) Validate() error {
	return validation.In(MFAAuthentication, TFAAuthentication, NoneAuthentication).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' or '%s'",
			a, MFAAuthentication, TFAAuthentication, NoneAuthentication)).
		Validate(a)
}

// Validate validates validation on AuthGrant.
func (r AuthGrant) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"GroupID": validation.Validate(r.GroupID, validation.Required),
		"RoleID":  validation.Validate(r.RoleID, validation.Required),
	})
}

// Validate validates CreateUserRequest.
func (r CreateUserRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Country":                  validation.Validate(r.Country, validation.Required),
		"Email":                    validation.Validate(r.Email, validation.Required, is.EmailFormat),
		"FirstName":                validation.Validate(r.FirstName, validation.Required),
		"LastName":                 validation.Validate(r.LastName, validation.Required),
		"AuthGrants":               validation.Validate(r.AuthGrants, validation.Required),
		"Notifications":            validation.Validate(r.Notifications),
		"AdditionalAuthentication": validation.Validate(r.AdditionalAuthentication, validation.Required),
	})
}

// Validate validates GetUserRequest.
func (r GetUserRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
	})
}

// Validate validates UpdateUserInfoRequest.
func (r UpdateUserInfoRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID":        validation.Validate(r.IdentityID, validation.Required),
		"FirstName":         validation.Validate(r.User.FirstName, validation.Required),
		"LastName":          validation.Validate(r.User.LastName, validation.Required),
		"Country":           validation.Validate(r.User.Country, validation.Required),
		"TimeZone":          validation.Validate(r.User.TimeZone, validation.Required),
		"PreferredLanguage": validation.Validate(r.User.PreferredLanguage, validation.Required),
		"SessionTimeOut":    validation.Validate(r.User.SessionTimeOut, validation.Required),
	})
}

// Validate validates UpdateUserNotificationsRequest.
func (r UpdateUserNotificationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID":    validation.Validate(r.IdentityID, validation.Required),
		"Notifications": validation.Validate(r.Notifications, validation.Required),
	})
}

// Validate validates UpdateUserAuthGrantsRequest.
func (r UpdateUserAuthGrantsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
		"AuthGrants": validation.Validate(r.AuthGrants, validation.Required),
	})
}

// Validate validates RemoveUserRequest.
func (r RemoveUserRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"uiIdentity": validation.Validate(r.IdentityID, validation.Required),
	})
}

// Validate validates UpdateTFARequest.
func (r UpdateTFARequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
		"Action": validation.Validate(r.Action, validation.Required, validation.In(TFAActionEnable, TFAActionDisable, TFAActionReset).
			Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'enable', 'disable' or 'reset'", r.Action))),
	})
}

// Validate validates UpdateMFARequest.
func (r UpdateMFARequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Value": validation.Validate(r.Value, validation.Required),
	})
}

func (i *iam) CreateUser(ctx context.Context, params CreateUserRequest) (*User, error) {
	logger := i.Log(ctx)
	logger.Debug("CreateUser")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateUser, ErrStructValidation, err)
	}

	uri, err := url.Parse("/identity-management/v3/user-admin/ui-identities")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateUser, err)
	}

	q := uri.Query()
	q.Add("sendEmail", strconv.FormatBool(params.SendEmail))

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
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
	logger := i.Log(ctx)
	logger.Debug("GetUser")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetUser, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetUser, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	q.Add("authGrants", strconv.FormatBool(params.AuthGrants))
	q.Add("notifications", strconv.FormatBool(params.Notifications))

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetUser, err)
	}

	var result User
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetUser, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetUser, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) ListUsers(ctx context.Context, params ListUsersRequest) ([]UserListItem, error) {
	logger := i.Log(ctx)
	logger.Debug("ListUsers")

	uri, err := url.Parse("/identity-management/v3/user-admin/ui-identities")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse the URL:\n%s", ErrListUsers, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	q.Add("authGrants", strconv.FormatBool(params.AuthGrants))
	if params.GroupID != nil {
		q.Add("groupId", strconv.FormatInt(*params.GroupID, 10))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
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
	logger := i.Log(ctx)
	logger.Debug("RemoveUser")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrRemoveUser, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrRemoveUser, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
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
	logger := i.Log(ctx)
	logger.Debug("UpdateUserAuthGrants")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateUserAuthGrants, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/auth-grants", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserAuthGrants, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserAuthGrants, err)
	}

	var result []AuthGrant

	resp, err := i.Exec(req, &result, params.AuthGrants)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateUserAuthGrants, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateUserAuthGrants, i.Error(resp))
	}

	return result, nil
}

func (i *iam) UpdateUserInfo(ctx context.Context, params UpdateUserInfoRequest) (*UserBasicInfo, error) {
	logger := i.Log(ctx)
	logger.Debug("UpdateUserInfo")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateUserInfo, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/basic-info", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserInfo, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserInfo, err)
	}

	var result UserBasicInfo
	resp, err := i.Exec(req, &result, params.User)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateUserInfo, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateUserInfo, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) UpdateUserNotifications(ctx context.Context, params UpdateUserNotificationsRequest) (*UserNotifications, error) {
	logger := i.Log(ctx)
	logger.Debug("UpdateUserNotifications")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateUserNotifications, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/notifications", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserNotifications, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateUserNotifications, err)
	}

	var result UserNotifications
	resp, err := i.Exec(req, &result, params.Notifications)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateUserNotifications, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateUserNotifications, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) UpdateTFA(ctx context.Context, params UpdateTFARequest) error {
	logger := i.Log(ctx)
	logger.Debug("UpdateTFA")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrUpdateTFA, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/tfa", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUpdateTFA, err)
	}

	q := uri.Query()
	q.Add("action", string(params.Action))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
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

func (i *iam) UpdateMFA(ctx context.Context, params UpdateMFARequest) error {
	logger := i.Log(ctx)
	logger.Debug("UpdateMFA")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrUpdateMFA, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/additionalAuthentication", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUpdateMFA, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUpdateMFA, err)
	}

	resp, err := i.Exec(req, nil, params.Value)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrUpdateMFA, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrUpdateMFA, i.Error(resp))
	}

	return nil
}

func (i *iam) ResetMFA(ctx context.Context, params ResetMFARequest) error {
	logger := i.Log(ctx)
	logger.Debug("ResetMFA")

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/additionalAuthentication/reset", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrResetMFA, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrResetMFA, err)
	}

	resp, err := i.Exec(req, nil, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrResetMFA, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrResetMFA, i.Error(resp))
	}

	return nil
}
