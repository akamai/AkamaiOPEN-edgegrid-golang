package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// LockUserRequest contains the request parameters for the LockUser endpoint.
	LockUserRequest struct {
		IdentityID string
	}

	// UnlockUserRequest contains the request parameters for the UnlockUser endpoint.
	UnlockUserRequest struct {
		IdentityID string
	}
)

var (
	// ErrLockUser is returned when LockUser fails.
	ErrLockUser = errors.New("lock user")

	// ErrUnlockUser is returned when UnlockUser fails.
	ErrUnlockUser = errors.New("unlock user")
)

// Validate validates LockUserRequest.
func (r LockUserRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
	})
}

// Validate validates UnlockUserRequest.
func (r UnlockUserRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
	})
}

func (i *iam) LockUser(ctx context.Context, params LockUserRequest) error {
	logger := i.Log(ctx)
	logger.Debug("LockUser")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrLockUser, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/lock", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrLockUser, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrLockUser, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrLockUser, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrLockUser, i.Error(resp))
	}

	return nil
}

func (i *iam) UnlockUser(ctx context.Context, params UnlockUserRequest) error {
	logger := i.Log(ctx)
	logger.Debug("UnlockUser")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrUnlockUser, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/unlock", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUnlockUser, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUnlockUser, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrUnlockUser, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrUnlockUser, i.Error(resp))
	}

	return nil
}
