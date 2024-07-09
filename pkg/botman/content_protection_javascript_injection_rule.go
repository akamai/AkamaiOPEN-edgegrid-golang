package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ContentProtectionJavaScriptInjectionRule interface supports creating, retrieving, modifying and removing content protection JavaScript injection rule
	// for a policy.
	ContentProtectionJavaScriptInjectionRule interface {
		// GetContentProtectionJavaScriptInjectionRuleList https://techdocs.akamai.com/content-protector/reference/get-content-protection-javascript-injection-rules
		GetContentProtectionJavaScriptInjectionRuleList(ctx context.Context, params GetContentProtectionJavaScriptInjectionRuleListRequest) (*GetContentProtectionJavaScriptInjectionRuleListResponse, error)

		// GetContentProtectionJavaScriptInjectionRule https://techdocs.akamai.com/content-protector/reference/get-content-protection-javascript-injection-rule
		GetContentProtectionJavaScriptInjectionRule(ctx context.Context, params GetContentProtectionJavaScriptInjectionRuleRequest) (map[string]interface{}, error)

		// CreateContentProtectionJavaScriptInjectionRule https://techdocs.akamai.com/content-protector/reference/post-content-protection-javascript-injection-rule
		CreateContentProtectionJavaScriptInjectionRule(ctx context.Context, params CreateContentProtectionJavaScriptInjectionRuleRequest) (map[string]interface{}, error)

		// UpdateContentProtectionJavaScriptInjectionRule https://techdocs.akamai.com/content-protector/reference/put-content-protection-javascript-injection-rule
		UpdateContentProtectionJavaScriptInjectionRule(ctx context.Context, params UpdateContentProtectionJavaScriptInjectionRuleRequest) (map[string]interface{}, error)

		// RemoveContentProtectionJavaScriptInjectionRule https://techdocs.akamai.com/content-protector/reference/delete-content-protection-javascript-injection-rule
		RemoveContentProtectionJavaScriptInjectionRule(ctx context.Context, params RemoveContentProtectionJavaScriptInjectionRuleRequest) error
	}

	// GetContentProtectionJavaScriptInjectionRuleListRequest is used to retrieve the content protection JavaScript injection rules for a policy.
	GetContentProtectionJavaScriptInjectionRuleListRequest struct {
		ConfigID                                   int64
		Version                                    int64
		SecurityPolicyID                           string
		ContentProtectionJavaScriptInjectionRuleID string
	}

	// GetContentProtectionJavaScriptInjectionRuleListResponse is used to retrieve the content protection JavaScript injection rules for a policy.
	GetContentProtectionJavaScriptInjectionRuleListResponse struct {
		ContentProtectionJavaScriptInjectionRules []map[string]interface{} `json:"contentProtectionJavaScriptInjectionRules"`
	}

	// GetContentProtectionJavaScriptInjectionRuleRequest is used to retrieve a specific content protection JavaScript injection rule.
	GetContentProtectionJavaScriptInjectionRuleRequest struct {
		ConfigID                                   int64
		Version                                    int64
		SecurityPolicyID                           string
		ContentProtectionJavaScriptInjectionRuleID string
	}

	// CreateContentProtectionJavaScriptInjectionRuleRequest is used to create a new content protection JavaScript injection rule for a specific policy.
	CreateContentProtectionJavaScriptInjectionRuleRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		JsonPayload      json.RawMessage
	}

	// UpdateContentProtectionJavaScriptInjectionRuleRequest is used to update details for a content protection JavaScript injection rule.
	UpdateContentProtectionJavaScriptInjectionRuleRequest struct {
		ConfigID                                   int64
		Version                                    int64
		SecurityPolicyID                           string
		ContentProtectionJavaScriptInjectionRuleID string
		JsonPayload                                json.RawMessage
	}

	// RemoveContentProtectionJavaScriptInjectionRuleRequest is used to remove an existing content protection JavaScript injection rule
	RemoveContentProtectionJavaScriptInjectionRuleRequest struct {
		ConfigID                                   int64
		Version                                    int64
		SecurityPolicyID                           string
		ContentProtectionJavaScriptInjectionRuleID string
	}
)

// Validate validates a GetContentProtectionJavaScriptInjectionRuleRequest.
func (v GetContentProtectionJavaScriptInjectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"ContentProtectionJavaScriptInjectionRuleID": validation.Validate(v.ContentProtectionJavaScriptInjectionRuleID, validation.Required),
	}.Filter()
}

// Validate validates a GetContentProtectionJavaScriptInjectionRuleListRequest.
func (v GetContentProtectionJavaScriptInjectionRuleListRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateContentProtectionJavaScriptInjectionRuleRequest.
func (v CreateContentProtectionJavaScriptInjectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateContentProtectionJavaScriptInjectionRuleRequest.
func (v UpdateContentProtectionJavaScriptInjectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"ContentProtectionJavaScriptInjectionRuleID": validation.Validate(v.ContentProtectionJavaScriptInjectionRuleID, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveContentProtectionJavaScriptInjectionRuleRequest.
func (v RemoveContentProtectionJavaScriptInjectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"ContentProtectionJavaScriptInjectionRuleID": validation.Validate(v.ContentProtectionJavaScriptInjectionRuleID, validation.Required),
	}.Filter()
}

func (b *botman) GetContentProtectionJavaScriptInjectionRule(ctx context.Context, params GetContentProtectionJavaScriptInjectionRuleRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetContentProtectionJavaScriptInjectionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-javascript-injection-rules/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.ContentProtectionJavaScriptInjectionRuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetContentProtectionJavaScriptInjectionRule request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetContentProtectionJavaScriptInjectionRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetContentProtectionJavaScriptInjectionRuleList(ctx context.Context, params GetContentProtectionJavaScriptInjectionRuleListRequest) (*GetContentProtectionJavaScriptInjectionRuleListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetContentProtectionJavaScriptInjectionRuleList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-javascript-injection-rules",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetContentProtectionJavaScriptInjectionRuleList request: %w", err)
	}

	var result GetContentProtectionJavaScriptInjectionRuleListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetContentProtectionJavaScriptInjectionRuleList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetContentProtectionJavaScriptInjectionRuleListResponse
	if params.ContentProtectionJavaScriptInjectionRuleID != "" {
		for _, val := range result.ContentProtectionJavaScriptInjectionRules {
			if val["contentProtectionJavaScriptInjectionRuleId"].(string) == params.ContentProtectionJavaScriptInjectionRuleID {
				filteredResult.ContentProtectionJavaScriptInjectionRules = append(filteredResult.ContentProtectionJavaScriptInjectionRules, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateContentProtectionJavaScriptInjectionRule(ctx context.Context, params UpdateContentProtectionJavaScriptInjectionRuleRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateContentProtectionJavaScriptInjectionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-javascript-injection-rules/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.ContentProtectionJavaScriptInjectionRuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateContentProtectionJavaScriptInjectionRule request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateContentProtectionJavaScriptInjectionRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateContentProtectionJavaScriptInjectionRule(ctx context.Context, params CreateContentProtectionJavaScriptInjectionRuleRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateContentProtectionJavaScriptInjectionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-javascript-injection-rules",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateContentProtectionJavaScriptInjectionRule request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateContentProtectionJavaScriptInjectionRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveContentProtectionJavaScriptInjectionRule(ctx context.Context, params RemoveContentProtectionJavaScriptInjectionRuleRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveContentProtectionJavaScriptInjectionRule")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-javascript-injection-rules/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.ContentProtectionJavaScriptInjectionRuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveContentProtectionJavaScriptInjectionRule request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveContentProtectionJavaScriptInjectionRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
