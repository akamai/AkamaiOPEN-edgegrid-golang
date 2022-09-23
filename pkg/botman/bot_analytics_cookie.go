package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The BotAnalyticsCookie interface supports retrieving and updating bot analytics cookie settings
	BotAnalyticsCookie interface {

		// GetBotAnalyticsCookie https://techdocs.akamai.com/bot-manager/reference/get-bot-analytics-cookie
		GetBotAnalyticsCookie(ctx context.Context, params GetBotAnalyticsCookieRequest) (map[string]interface{}, error)

		// UpdateBotAnalyticsCookie https://techdocs.akamai.com/bot-manager/reference/put-bot-analytics-cookie
		UpdateBotAnalyticsCookie(ctx context.Context, params UpdateBotAnalyticsCookieRequest) (map[string]interface{}, error)
	}

	// GetBotAnalyticsCookieRequest is used to retrieve the bot analytics cookie settings
	GetBotAnalyticsCookieRequest struct {
		ConfigID int64
		Version  int64
	}

	// UpdateBotAnalyticsCookieRequest is used to modify bot analytics cookie settings
	UpdateBotAnalyticsCookieRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}
)

// Validate validates a GetBotAnalyticsCookieRequest.
func (v GetBotAnalyticsCookieRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateBotAnalyticsCookieRequest.
func (v UpdateBotAnalyticsCookieRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetBotAnalyticsCookie(ctx context.Context, params GetBotAnalyticsCookieRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetBotAnalyticsCookie")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/bot-analytics-cookie",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotAnalyticsCookie request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotAnalyticsCookie request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateBotAnalyticsCookie(ctx context.Context, params UpdateBotAnalyticsCookieRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateBotAnalyticsCookie")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/bot-analytics-cookie",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateBotAnalyticsCookie request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateBotAnalyticsCookie request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
