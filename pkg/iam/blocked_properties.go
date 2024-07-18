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
	// BlockedProperties is the IAM user blocked properties API interface
	BlockedProperties interface {
		// ListBlockedProperties returns all properties a user doesn't have access to in a group
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-blocked-properties
		ListBlockedProperties(context.Context, ListBlockedPropertiesRequest) ([]int64, error)

		// UpdateBlockedProperties removes or grants user access to properties
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-blocked-properties
		UpdateBlockedProperties(context.Context, UpdateBlockedPropertiesRequest) ([]int64, error)
	}

	// ListBlockedPropertiesRequest contains the request parameters for the list blocked properties endpoint
	ListBlockedPropertiesRequest struct {
		IdentityID string
		GroupID    int64
	}

	// UpdateBlockedPropertiesRequest contains the request parameters for the update blocked properties endpoint
	UpdateBlockedPropertiesRequest struct {
		IdentityID string
		GroupID    int64
		Properties []int64
	}
)

var (
	// ErrListBlockedProperties is returned when ListBlockedPropertiesRequest fails
	ErrListBlockedProperties = errors.New("list blocked properties")

	// ErrUpdateBlockedProperties is returned when UpdateBlockedPropertiesRequest fails
	ErrUpdateBlockedProperties = errors.New("update blocked properties")
)

// Validate validates ListBlockedPropertiesRequest
func (r ListBlockedPropertiesRequest) Validate() error {
	return validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	}.Filter()
}

// Validate validates UpdateBlockedPropertiesRequest
func (r UpdateBlockedPropertiesRequest) Validate() error {
	return validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	}.Filter()
}

func (i *iam) ListBlockedProperties(ctx context.Context, params ListBlockedPropertiesRequest) ([]int64, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListBlockedProperties, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/groups/%d/blocked-properties", params.IdentityID, params.GroupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListBlockedProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListBlockedProperties, err)
	}

	var blockedProperties []int64
	resp, err := i.Exec(req, &blockedProperties)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListBlockedProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListBlockedProperties, i.Error(resp))
	}

	return blockedProperties, nil
}

func (i *iam) UpdateBlockedProperties(ctx context.Context, params UpdateBlockedPropertiesRequest) ([]int64, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateBlockedProperties, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/groups/%d/blocked-properties", params.IdentityID, params.GroupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateBlockedProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateBlockedProperties, err)
	}

	var blockedProperties []int64
	resp, err := i.Exec(req, &blockedProperties, params.Properties)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateBlockedProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateBlockedProperties, i.Error(resp))
	}

	return blockedProperties, nil
}
