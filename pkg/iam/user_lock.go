package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// UserLock is the IAM user lock/unlock API interface
	UserLock interface {
		// LockUser lock the user
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/post-ui-identity-lock
		LockUser(context.Context, LockUserRequest) error

		// UnlockUser release the lock on a user's account
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/post-ui-identity-unlock
		UnlockUser(context.Context, UnlockUserRequest) error
	}

	// LockUserRequest contains the request parameters of the lock user endpoint
	LockUserRequest struct {
		IdentityID string
	}

	// UnlockUserRequest contains the request parameters of the unlock user endpoint
	UnlockUserRequest struct {
		IdentityID string
	}
)

var (
	// ErrLockUser is returned when LockUser fails
	ErrLockUser = errors.New("lock user")

	// ErrUnlockUser is returned when UnlockUser fails
	ErrUnlockUser = errors.New("unlock user")
)

// Validate validates LockUserRequest
func (r LockUserRequest) Validate() error {
	return validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

// Validate validates UnlockUserRequest
func (r UnlockUserRequest) Validate() error {
	return validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
	}.Filter()
}

func (i *iam) LockUser(ctx context.Context, params LockUserRequest) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrLockUser, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/lock", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrLockUser, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrLockUser, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrLockUser, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrLockUser, i.Error(resp))
	}

	return nil
}

func (i *iam) UnlockUser(ctx context.Context, params UnlockUserRequest) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrUnlockUser, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v2/user-admin/ui-identities/%s/unlock", params.IdentityID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUnlockUser, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUnlockUser, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrUnlockUser, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrUnlockUser, i.Error(resp))
	}

	return nil
}
