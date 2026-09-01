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
	//	GetUserRiskResponseStrategyRequest represents the request to get user risk response strategy.
	GetUserRiskResponseStrategyRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64
	}

	//	UpsertUserRiskResponseStrategyRequest represents the request to upsert user risk response strategy.
	UpsertUserRiskResponseStrategyRequest struct {
		// ConfigID is the ID of the security configuration.
		ConfigID int64

		// Version is the version of the security configuration.
		Version int64

		// JsonPayload contains the request payload for user risk response strategy
		JsonPayload json.RawMessage
	}
)

// Validate a GetUserRiskResponseStrategyRequest.
func (v GetUserRiskResponseStrategyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate an UpdateUserRiskResponseStrategyRequest.
func (v UpsertUserRiskResponseStrategyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (ap *accountProtection) GetUserRiskResponseStrategy(ctx context.Context, params GetUserRiskResponseStrategyRequest) (map[string]any, error) {
	logger := ap.Log(ctx)
	logger.Debug("GetUserRiskResponseStrategy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/account-protection/user-risk-response-strategy",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetUserRiskResponseStrategy request: %w", err)
	}

	var result map[string]interface{}
	resp, err := ap.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetUserRiskResponseStrategy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}

	return result, nil
}

func (ap *accountProtection) UpsertUserRiskResponseStrategy(ctx context.Context, params UpsertUserRiskResponseStrategyRequest) (map[string]any, error) {
	logger := ap.Log(ctx)
	logger.Debug("UpsertUserRiskResponseStrategy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/account-protection/user-risk-response-strategy",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpsertUserRiskResponseStrategyRequest request: %w", err)
	}

	var result map[string]interface{}
	resp, err := ap.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpsertUserRiskResponseStrategyRequest failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, ap.Error(resp)
	}

	return result, nil
}
