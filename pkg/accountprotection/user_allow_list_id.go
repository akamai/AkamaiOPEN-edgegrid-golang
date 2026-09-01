package accountprotection

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetUserAllowListIDRequest represents the request to get user allow list ID.
	GetUserAllowListIDRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64
	}

	//	UpsertUserAllowListIDRequest represents the request to upsert user allow list ID.
	UpsertUserAllowListIDRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// JsonPayload contains the request payload for user allow list
		JsonPayload json.RawMessage
	}

	// 	DeleteUserAllowListIDRequest represents the request to delete user allow list ID.
	DeleteUserAllowListIDRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64
	}
)

// Validate a GetUserAllowListIDRequest.
func (v GetUserAllowListIDRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate an UpsertUserAllowListIDRequest.
func (v UpsertUserAllowListIDRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate a DeleteUserAllowListIDRequest.
func (v DeleteUserAllowListIDRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (ap *accountProtection) GetUserAllowListID(ctx context.Context, params GetUserAllowListIDRequest) (map[string]any, error) {
	logger := ap.Log(ctx)
	logger.Debug("GetUserAllowListID")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/account-protection/user-allow-list-id",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetUserAllowListID request: %w", err)
	}

	var result map[string]interface{}
	resp, err := ap.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetUserAllowListID request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}

	return result, nil
}

func (ap *accountProtection) UpsertUserAllowListID(ctx context.Context, params UpsertUserAllowListIDRequest) (map[string]any, error) {
	logger := ap.Log(ctx)
	logger.Debug("UpsertUserAllowListID")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/account-protection/user-allow-list-id",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpsertUserAllowListIDRequest request: %w", err)
	}

	var result map[string]interface{}
	resp, err := ap.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpsertUserAllowListIDRequest request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}

	return result, nil
}

func (ap *accountProtection) DeleteUserAllowListID(ctx context.Context, params DeleteUserAllowListIDRequest) error {
	logger := ap.Log(ctx)
	logger.Debug("DeleteUserAllowListID")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/advanced-settings/account-protection/user-allow-list-id",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create DeleteUserAllowListID request: %w", err)
	}

	resp, err := ap.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("DeleteUserAllowListID request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return ap.Error(resp)
	}

	return nil
}
