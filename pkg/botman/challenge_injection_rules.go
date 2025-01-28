package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ChallengeInjectionRules interface supports retrieving and updating the challenge injection rules for a configuration
	ChallengeInjectionRules interface {
		// GetChallengeInjectionRules https://techdocs.akamai.com/bot-manager/reference/get-challenge-injection-rules
		GetChallengeInjectionRules(ctx context.Context, params GetChallengeInjectionRulesRequest) (map[string]interface{}, error)
		// UpdateChallengeInjectionRules https://techdocs.akamai.com/bot-manager/reference/put-challenge-injection-rules
		UpdateChallengeInjectionRules(ctx context.Context, params UpdateChallengeInjectionRulesRequest) (map[string]interface{}, error)
	}

	// GetChallengeInjectionRulesRequest is used to retrieve challenge injection rules
	GetChallengeInjectionRulesRequest struct {
		ConfigID int64
		Version  int64
	}

	// UpdateChallengeInjectionRulesRequest is used to modify challenge injection rules
	UpdateChallengeInjectionRulesRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}
)

// Validate validates a GetChallengeInjectionRulesRequest.
func (v GetChallengeInjectionRulesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates an UpdateChallengeInjectionRulesRequest.
func (v UpdateChallengeInjectionRulesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	})
}

func (b *botman) GetChallengeInjectionRules(ctx context.Context, params GetChallengeInjectionRulesRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetChallengeInjectionRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-injection-rules",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetChallengeInjectionRules request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetChallengeInjectionRules request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateChallengeInjectionRules(ctx context.Context, params UpdateChallengeInjectionRulesRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateChallengeInjectionRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-injection-rules",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateChallengeInjectionRules request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateChallengeInjectionRules request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
