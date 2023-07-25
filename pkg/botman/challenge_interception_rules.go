package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ChallengeInterceptionRules interface supports retrieving and updating the challenge interception rules for a
	// configuration
	// Deprecated: this interface will be removed in a future release. Use ChallengeInjectionRules instead.
	ChallengeInterceptionRules interface {
		// GetChallengeInterceptionRules https://techdocs.akamai.com/bot-manager/reference/get-challenge-interception-rules
		// Deprecated: this method will be removed in a future release. Use GetChallengeInjectionRules instead.
		GetChallengeInterceptionRules(ctx context.Context, params GetChallengeInterceptionRulesRequest) (map[string]interface{}, error)
		// UpdateChallengeInterceptionRules https://techdocs.akamai.com/bot-manager/reference/put-challenge-interception-rules
		// Deprecated: this method will be removed in a future release. Use UpdateChallengeInjectionRules instead.
		UpdateChallengeInterceptionRules(ctx context.Context, params UpdateChallengeInterceptionRulesRequest) (map[string]interface{}, error)
	}

	// GetChallengeInterceptionRulesRequest is used to retrieve challenge interception rules
	// Deprecated: this struct will be removed in a future release. Use GetChallengeInjectionRulesRequest instead.
	GetChallengeInterceptionRulesRequest struct {
		ConfigID int64
		Version  int64
	}

	// UpdateChallengeInterceptionRulesRequest is used to modify challenge interception rules
	// Deprecated: this struct will be removed in a future release. Use UpdateChallengeInjectionRulesRequest instead.
	UpdateChallengeInterceptionRulesRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}
)

// Validate validates a GetChallengeInterceptionRulesRequest.
func (v GetChallengeInterceptionRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateChallengeInterceptionRulesRequest.
func (v UpdateChallengeInterceptionRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetChallengeInterceptionRules(ctx context.Context, params GetChallengeInterceptionRulesRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetChallengeInterceptionRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-interception-rules",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetChallengeInterceptionRules request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetChallengeInterceptionRules request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateChallengeInterceptionRules(ctx context.Context, params UpdateChallengeInterceptionRulesRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateChallengeInterceptionRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-interception-rules",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateChallengeInterceptionRules request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateChallengeInterceptionRules request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
