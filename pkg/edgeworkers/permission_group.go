package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetPermissionGroupRequest contains parameters used to get a permission group
	GetPermissionGroupRequest struct {
		GroupID string
	}

	// PermissionGroup represents a single permission group object
	PermissionGroup struct {
		ID           int64    `json:"groupId"`
		Name         string   `json:"groupName"`
		Capabilities []string `json:"capabilities"`
	}

	// ListPermissionGroupsResponse represents a response object returned by ListPermissionGroups
	ListPermissionGroupsResponse struct {
		PermissionGroups []PermissionGroup `json:"groups"`
	}
)

// Validate validates GetPermissionGroupRequest
func (g GetPermissionGroupRequest) Validate() error {
	return validation.Errors{
		"GroupID": validation.Validate(g.GroupID, validation.Required),
	}.Filter()
}

var (
	// ErrGetPermissionGroup is returned in case an error occurs on GetPermissionGroup operation
	ErrGetPermissionGroup = errors.New("get a permission group")
	// ErrListPermissionGroups is returned in case an error occurs on ListPermissionGroups operation
	ErrListPermissionGroups = errors.New("list permission groups")
)

func (e *edgeworkers) GetPermissionGroup(ctx context.Context, params GetPermissionGroupRequest) (*PermissionGroup, error) {
	logger := e.Log(ctx)
	logger.Debug("GetPermissionGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPermissionGroup, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/groups/%s", params.GroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPermissionGroup, err)
	}

	var result PermissionGroup
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPermissionGroup, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPermissionGroup, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) ListPermissionGroups(ctx context.Context) (*ListPermissionGroupsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListPermissionGroups")

	uri := "/edgeworkers/v1/groups"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPermissionGroups, err)
	}

	var result ListPermissionGroupsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPermissionGroups, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPermissionGroups, e.Error(resp))
	}

	return &result, nil
}
