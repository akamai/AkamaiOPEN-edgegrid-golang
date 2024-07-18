package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// IPAllowlist is the IAM IP allowlist API interface
	IPAllowlist interface {
		// DisableIPAllowlist disables IP allowlist on your account. After you disable IP allowlist,
		// users can access Control Center regardless of their IP address or who assigns it.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-allowlist-disable
		DisableIPAllowlist(context.Context) error

		// EnableIPAllowlist enables IP allowlist on your account. Before you enable IP allowlist,
		// add at least one IP address to allow access to Control Center.
		// The allowlist can't be empty with IP allowlist enabled.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-allowlist-enable
		EnableIPAllowlist(context.Context) error

		// GetIPAllowlistStatus indicates whether IP allowlist is enabled or disabled on your account.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-allowlist-status
		GetIPAllowlistStatus(context.Context) (*GetIPAllowlistStatusResponse, error)
	}

	// GetIPAllowlistStatusResponse contains response from GetIPAllowlistStatus operation
	GetIPAllowlistStatusResponse struct {
		Enabled bool `json:"enabled"`
	}
)

var (
	// ErrDisableIPAllowlist is returned when DisableIPAllowlist fails
	ErrDisableIPAllowlist = errors.New("disable ip allowlist")
	// ErrEnableIPAllowlist is returned when EnableIPAllowlist fails
	ErrEnableIPAllowlist = errors.New("enable ip allowlist")
	// ErrGetIPAllowlistStatus is returned when GetIPAllowlistStatus fails
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
