package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The JavascriptInjection interface supports retrieving and updating the javascript injection settings for a
	// configuration
	JavascriptInjection interface {
		// GetJavascriptInjection https://techdocs.akamai.com/bot-manager/reference/get-javascript-injection-rules
		GetJavascriptInjection(ctx context.Context, params GetJavascriptInjectionRequest) (map[string]interface{}, error)
		// UpdateJavascriptInjection https://techdocs.akamai.com/bot-manager/reference/put-javascript-injection-rules
		UpdateJavascriptInjection(ctx context.Context, params UpdateJavascriptInjectionRequest) (map[string]interface{}, error)
	}

	// GetJavascriptInjectionRequest is used to retrieve javascript injection settings
	GetJavascriptInjectionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
	}

	// UpdateJavascriptInjectionRequest is used to modify javascript injection settings
	UpdateJavascriptInjectionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		JsonPayload      json.RawMessage
	}
)

// Validate validates a GetJavascriptInjectionRequest.
func (v GetJavascriptInjectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateJavascriptInjectionRequest.
func (v UpdateJavascriptInjectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetJavascriptInjection(ctx context.Context, params GetJavascriptInjectionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetJavascriptInjection")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/javascript-injection",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetJavascriptInjection request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetJavascriptInjection request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateJavascriptInjection(ctx context.Context, params UpdateJavascriptInjectionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateJavascriptInjection")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/javascript-injection",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateJavascriptInjection request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateJavascriptInjection request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
