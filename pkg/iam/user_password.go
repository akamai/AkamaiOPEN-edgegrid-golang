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
	// UserPassword is the IAM managing user's password API interface
	UserPassword interface {
		// ResetUserPassword optionally sends a one-time use password to the user.
		// If you send the email with the password directly to the user, the response for this operation doesn't include that password.
		// If you don't send the password to the user through email, the password is included in the response.
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/post-reset-password
		ResetUserPassword(context.Context, ResetUserPasswordRequest) (*ResetUserPasswordResponse, error)

		// SetUserPassword sets a specific password for a user
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/post-set-password
		SetUserPassword(context.Context, SetUserPasswordRequest) error
	}

	// ResetUserPasswordRequest contains the request parameters of the reset user password endpoint
	ResetUserPasswordRequest struct {
		IdentityID string
		SendEmail  bool
	}

	// ResetUserPasswordResponse contains the response from the reset user password endpoint
	ResetUserPasswordResponse struct {
		NewPassword string `json:"newPassword"`
	}

	// SetUserPasswordRequest contains the request parameters of the set user password endpoint
	SetUserPasswordRequest struct {
		IdentityID  string
		NewPassword string
	}
)

var (
	// ErrResetUserPassword is returned when ResetUserPassword fails
	ErrResetUserPassword = errors.New("reset user password")

	// ErrSetUserPassword is returned when SetUserPassword fails
	ErrSetUserPassword = errors.New("set user password")
)

// Validate validates ResetUserPasswordRequest
func (r ResetUserPasswordRequest) Validate() error {
	return validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

// Validate validates SetUserPasswordRequest
func (r SetUserPasswordRequest) Validate() error {
	return validation.Errors{
		"IdentityID":  validation.Validate(r.IdentityID, validation.Required),
		"NewPassword": validation.Validate(r.NewPassword, validation.Required),
	}.Filter()
}

func (i *iam) ResetUserPassword(ctx context.Context, params ResetUserPasswordRequest) (*ResetUserPasswordResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrResetUserPassword, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/reset-password", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrResetUserPassword, err)
	}

	q := u.Query()
	q.Add("actions", strconv.FormatBool(params.SendEmail))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrResetUserPassword, err)
	}

	var result ResetUserPasswordResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrResetUserPassword, err)
	}

	if !((!params.SendEmail && resp.StatusCode == http.StatusOK) || (params.SendEmail && resp.StatusCode == http.StatusNoContent)) {
		return nil, fmt.Errorf("%s: %w", ErrResetUserPassword, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) SetUserPassword(ctx context.Context, params SetUserPasswordRequest) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrSetUserPassword, ErrStructValidation, err)
	}

	u := fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/set-password", params.IdentityID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrSetUserPassword, err)
	}

	resp, err := i.Exec(req, nil, params.NewPassword)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrSetUserPassword, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrSetUserPassword, i.Error(resp))
	}

	return nil
}
