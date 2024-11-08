package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListBlockedPropertiesRequest contains the request parameters for the ListBlockedProperties endpoint.
	ListBlockedPropertiesRequest struct {
		IdentityID string
		GroupID    int64
	}

	// UpdateBlockedPropertiesRequest contains the request parameters for the UpdateBlockedProperties endpoint.
	UpdateBlockedPropertiesRequest struct {
		IdentityID string
		GroupID    int64
		Properties []int64
	}
)

var (
	// ErrListBlockedProperties is returned when ListBlockedPropertiesRequest fails.
	ErrListBlockedProperties = errors.New("list blocked properties")

	// ErrUpdateBlockedProperties is returned when UpdateBlockedPropertiesRequest fails.
	ErrUpdateBlockedProperties = errors.New("update blocked properties")
)

// Validate validates ListBlockedPropertiesRequest.
func (r ListBlockedPropertiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	})
}

// Validate validates UpdateBlockedPropertiesRequest.
func (r UpdateBlockedPropertiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IdentityID": validation.Validate(r.IdentityID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	})
}

func (i *iam) ListBlockedProperties(ctx context.Context, params ListBlockedPropertiesRequest) ([]int64, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListBlockedProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/groups/%d/blocked-properties", params.IdentityID, params.GroupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListBlockedProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListBlockedProperties, err)
	}

	var blockedProperties []int64
	resp, err := i.Exec(req, &blockedProperties)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListBlockedProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListBlockedProperties, i.Error(resp))
	}

	return blockedProperties, nil
}

func (i *iam) UpdateBlockedProperties(ctx context.Context, params UpdateBlockedPropertiesRequest) ([]int64, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateBlockedProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ui-identities/%s/groups/%d/blocked-properties", params.IdentityID, params.GroupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateBlockedProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateBlockedProperties, err)
	}

	var blockedProperties []int64
	resp, err := i.Exec(req, &blockedProperties, params.Properties)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateBlockedProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateBlockedProperties, i.Error(resp))
	}

	return blockedProperties, nil
}
