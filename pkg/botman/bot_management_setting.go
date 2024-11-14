package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The BotManagementSetting interface supports retrieving and updating bot management settings
	BotManagementSetting interface {

		// GetBotManagementSetting todo: add link
		GetBotManagementSetting(ctx context.Context, params GetBotManagementSettingRequest) (map[string]interface{}, error)

		// UpdateBotManagementSetting todo: add link
		UpdateBotManagementSetting(ctx context.Context, params UpdateBotManagementSettingRequest) (map[string]interface{}, error)
	}

	// GetBotManagementSettingRequest is used to retrieve the bot management settings
	GetBotManagementSettingRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
	}

	// UpdateBotManagementSettingRequest is used to modify bot management settings
	UpdateBotManagementSettingRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		JsonPayload      json.RawMessage
	}
)

// Validate validates a GetBotManagementSettingRequest.
func (v GetBotManagementSettingRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateBotManagementSettingRequest.
func (v UpdateBotManagementSettingRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetBotManagementSetting(ctx context.Context, params GetBotManagementSettingRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetBotManagementSetting")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bot-management-settings",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotManagementSetting request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotManagementSetting request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateBotManagementSetting(ctx context.Context, params UpdateBotManagementSettingRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateBotManagementSetting")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bot-management-settings",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateBotManagementSetting request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateBotManagementSetting request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
