package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The BotAnalyticsSettings interface supports retrieving and updating bot analytics settings.
	BotAnalyticsSettings interface {

		// GetBotAnalyticsSettings https://techdocs.akamai.com/bot-manager/reference/get-bot-analytics-settings
		GetBotAnalyticsSettings(ctx context.Context, params GetBotAnalyticsSettingsRequest) (map[string]any, error)

		// UpdateBotAnalyticsSettings https://techdocs.akamai.com/bot-manager/reference/put-bot-analytics-settings
		UpdateBotAnalyticsSettings(ctx context.Context, params UpdateBotAnalyticsSettingsRequest) (map[string]any, error)
	}

	// GetBotAnalyticsSettingsRequest is used to retrieve bot analytics settings.
	GetBotAnalyticsSettingsRequest struct {
		// ConfigID is the unique identifier of the security configuration.
		ConfigID int64

		// Version is the version number of the security configuration.
		Version int64
	}

	// UpdateBotAnalyticsSettingsRequest is used to modify bot analytics settings.
	UpdateBotAnalyticsSettingsRequest struct {
		// ConfigID is the unique identifier of the security configuration.
		ConfigID int64

		// Version is the version number of the security configuration.
		Version int64

		// JSONPayload contains the bot analytics settings to be updated.
		JSONPayload json.RawMessage
	}
)

// Validate validates a GetBotAnalyticsSettingsRequest.
func (v GetBotAnalyticsSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateBotAnalyticsSettingsRequest.
func (v UpdateBotAnalyticsSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JSONPayload": validation.Validate(v.JSONPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetBotAnalyticsSettings(ctx context.Context, params GetBotAnalyticsSettingsRequest) (map[string]any, error) {
	logger := b.Log(ctx)
	logger.Debug("GetBotAnalyticsSettings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewGet(ctx, "/appsec/v1/configs/%d/versions/%d/advanced-settings/bot-analytics-settings", params.ConfigID, params.Version).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotAnalyticsSettings request: %w", err)
	}

	var result map[string]any
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotAnalyticsSettings request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateBotAnalyticsSettings(ctx context.Context, params UpdateBotAnalyticsSettingsRequest) (map[string]any, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateBotAnalyticsSettings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewPut(ctx, "/appsec/v1/configs/%d/versions/%d/advanced-settings/bot-analytics-settings", params.ConfigID, params.Version).
		WithBody(params.JSONPayload).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateBotAnalyticsSettings request: %w", err)
	}

	var result map[string]any
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("UpdateBotAnalyticsSettings request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
