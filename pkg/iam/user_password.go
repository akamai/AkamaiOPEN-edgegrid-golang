package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ResetUserPasswordRequest contains the request parameters for the ResetUserPassword endpoint.
	ResetUserPasswordRequest struct {
		IdentityID string
		SendEmail  bool
	}

	// ResetUserPasswordResponse contains the response from the ResetUserPassword endpoint.
	ResetUserPasswordResponse struct {
		NewPassword string `json:"newPassword"`
	}

	// SetUserPasswordRequest contains the request parameters for the SetUserPassword endpoint.
	SetUserPasswordRequest struct {
		IdentityID  string `json:"-"`
		NewPassword string `json:"newPassword"`
	}
)

var (
	// ErrResetUserPassword is returned when ResetUserPassword fails.
	ErrResetUserPassword = errors.New("reset user password")

	// ErrSetUserPassword is returned when SetUserPassword fails.
	ErrSetUserPassword = errors.New("set user password")
)

// Validate validates ResetUserPasswordRequest.
func (r ResetUserPasswordRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
	})
}

// Validate validates SetUserPasswordRequest.
func (r SetUserPasswordRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID":  validation.Validate(r.IdentityID, validation.Required),
		"NewPassword": validation.Validate(r.NewPassword, validation.Required),
	})
}

func (i *iam) ResetUserPassword(ctx context.Context, params ResetUserPasswordRequest) (*ResetUserPasswordResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ResetUserPassword")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrResetUserPassword, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/reset-password", params.IdentityID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrResetUserPassword, err)
	}

	q := uri.Query()
	q.Add("sendEmail", strconv.FormatBool(params.SendEmail))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrResetUserPassword, err)
	}

	var result ResetUserPasswordResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrResetUserPassword, err)
	}
	defer session.CloseResponseBody(resp)

	if !((!params.SendEmail && resp.StatusCode == http.StatusOK) || (params.SendEmail && resp.StatusCode == http.StatusNoContent)) {
		return nil, fmt.Errorf("%s: %w", ErrResetUserPassword, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) SetUserPassword(ctx context.Context, params SetUserPasswordRequest) error {
	logger := i.Log(ctx)
	logger.Debug("SetUserPassword")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrSetUserPassword, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/set-password", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrSetUserPassword, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrSetUserPassword, err)
	}

	resp, err := i.Exec(req, nil, params)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrSetUserPassword, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrSetUserPassword, i.Error(resp))
	}

	return nil
}
