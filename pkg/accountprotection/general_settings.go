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
	// GetGeneralSettingsRequest represent the request to get account protection general settings.
	GetGeneralSettingsRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// SecurityPolicyID is the ID of the security policy
		SecurityPolicyID string
	}

	// UpsertGeneralSettingsRequest represents the request to upsert account protection general settings.
	UpsertGeneralSettingsRequest struct {

		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// SecurityPolicyID is the ID of the security policy
		SecurityPolicyID string

		// JsonPayload contains the values of the general settings
		JsonPayload json.RawMessage
	}
)

// Validate a GetAccountProtectionGeneralSettingsRequest.
func (v GetGeneralSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate an UpdateAccountProtectionGeneralSettingsRequest.
func (v UpsertGeneralSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (ap *accountProtection) GetGeneralSettings(ctx context.Context, params GetGeneralSettingsRequest) (map[string]any, error) {
	logger := ap.Log(ctx)
	logger.Debug("GetGeneralSettings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/account-protection-settings",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetGeneralSettingsRequest request: %w", err)
	}

	var result map[string]interface{}
	resp, err := ap.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetGeneralSettingsRequest failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}

	return result, nil
}

func (ap *accountProtection) UpsertGeneralSettings(ctx context.Context, params UpsertGeneralSettingsRequest) (map[string]any, error) {
	logger := ap.Log(ctx)
	logger.Debug("UpsertGeneralSettings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/account-protection-settings",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpsertGeneralSettingsRequest: %w", err)
	}

	var result map[string]interface{}
	resp, err := ap.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpsertGeneralSettingsRequest failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}

	return result, nil
}
