package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// GetIPAllowlistStatusResponse contains response from the GetIPAllowlistStatus endpoint.
	GetIPAllowlistStatusResponse struct {
		Enabled bool `json:"enabled"`
	}
)

var (
	// ErrDisableIPAllowlist is returned when DisableIPAllowlist fails.
	ErrDisableIPAllowlist = errors.New("disable ip allowlist")
	// ErrEnableIPAllowlist is returned when EnableIPAllowlist fails.
	ErrEnableIPAllowlist = errors.New("enable ip allowlist")
	// ErrGetIPAllowlistStatus is returned when GetIPAllowlistStatus fails.
	ErrGetIPAllowlistStatus = errors.New("get ip allowlist status")
)

func (i *iam) DisableIPAllowlist(ctx context.Context) error {
	logger := i.Log(ctx)
	logger.Debug("DisableIPAllowlist")

	uri := "/identity-management/v3/user-admin/ip-acl/allowlist/disable"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDisableIPAllowlist, err)
	}

	resp, err := i.Exec(req, nil, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDisableIPAllowlist, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDisableIPAllowlist, i.Error(resp))
	}

	return nil
}

func (i *iam) EnableIPAllowlist(ctx context.Context) error {
	logger := i.Log(ctx)
	logger.Debug("EnableIPAllowlist")

	uri := "/identity-management/v3/user-admin/ip-acl/allowlist/enable"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrEnableIPAllowlist, err)
	}

	resp, err := i.Exec(req, nil, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrEnableIPAllowlist, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrEnableIPAllowlist, i.Error(resp))
	}

	return nil
}

func (i *iam) GetIPAllowlistStatus(ctx context.Context) (*GetIPAllowlistStatusResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("GetIPAllowlistStatus")

	uri := "/identity-management/v3/user-admin/ip-acl/allowlist/status"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetIPAllowlistStatus, err)
	}

	var result GetIPAllowlistStatusResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetIPAllowlistStatus, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetIPAllowlistStatus, i.Error(resp))
	}

	return &result, nil
}
