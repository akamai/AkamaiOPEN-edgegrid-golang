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
	// The ContentProtectionRule interface supports creating, retrieving, modifying and removing content protection rule
	// for a policy.
	ContentProtectionRule interface {
		// GetContentProtectionRuleList https://techdocs.akamai.com/content-protector/reference/get-content-protection-rules
		GetContentProtectionRuleList(ctx context.Context, params GetContentProtectionRuleListRequest) (*GetContentProtectionRuleListResponse, error)

		// GetContentProtectionRule https://techdocs.akamai.com/content-protector/reference/get-content-protection-rule
		GetContentProtectionRule(ctx context.Context, params GetContentProtectionRuleRequest) (map[string]interface{}, error)

		// CreateContentProtectionRule https://techdocs.akamai.com/content-protector/reference/post-content-protection-rule
		CreateContentProtectionRule(ctx context.Context, params CreateContentProtectionRuleRequest) (map[string]interface{}, error)

		// UpdateContentProtectionRule https://techdocs.akamai.com/content-protector/reference/put-content-protection-rule
		UpdateContentProtectionRule(ctx context.Context, params UpdateContentProtectionRuleRequest) (map[string]interface{}, error)

		// RemoveContentProtectionRule https://techdocs.akamai.com/content-protector/reference/delete-content-protection-rule
		RemoveContentProtectionRule(ctx context.Context, params RemoveContentProtectionRuleRequest) error
	}

	// GetContentProtectionRuleListRequest is used to retrieve the content protection rules for a policy.
	GetContentProtectionRuleListRequest struct {
		ConfigID                int64
		Version                 int64
		SecurityPolicyID        string
		ContentProtectionRuleID string
	}

	// GetContentProtectionRuleListResponse is used to retrieve the content protection rules for a policy.
	GetContentProtectionRuleListResponse struct {
		ContentProtectionRules []map[string]interface{} `json:"contentProtectionRules"`
	}

	// GetContentProtectionRuleRequest is used to retrieve a specific content protection rule.
	GetContentProtectionRuleRequest struct {
		ConfigID                int64
		Version                 int64
		SecurityPolicyID        string
		ContentProtectionRuleID string
	}

	// CreateContentProtectionRuleRequest is used to create a new content protection rule for a specific policy.
	CreateContentProtectionRuleRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		JsonPayload      json.RawMessage
	}

	// UpdateContentProtectionRuleRequest is used to update details for a content protection rule.
	UpdateContentProtectionRuleRequest struct {
		ConfigID                int64
		Version                 int64
		SecurityPolicyID        string
		ContentProtectionRuleID string
		JsonPayload             json.RawMessage
	}

	// RemoveContentProtectionRuleRequest is used to remove an existing content protection rule
	RemoveContentProtectionRuleRequest struct {
		ConfigID                int64
		Version                 int64
		SecurityPolicyID        string
		ContentProtectionRuleID string
	}
)

// Validate validates a GetContentProtectionRuleRequest.
func (v GetContentProtectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":                validation.Validate(v.ConfigID, validation.Required),
		"Version":                 validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID":        validation.Validate(v.SecurityPolicyID, validation.Required),
		"ContentProtectionRuleID": validation.Validate(v.ContentProtectionRuleID, validation.Required),
	}.Filter()
}

// Validate validates a GetContentProtectionRuleListRequest.
func (v GetContentProtectionRuleListRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateContentProtectionRuleRequest.
func (v CreateContentProtectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateContentProtectionRuleRequest.
func (v UpdateContentProtectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":                validation.Validate(v.ConfigID, validation.Required),
		"Version":                 validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID":        validation.Validate(v.SecurityPolicyID, validation.Required),
		"ContentProtectionRuleID": validation.Validate(v.ContentProtectionRuleID, validation.Required),
		"JsonPayload":             validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveContentProtectionRuleRequest.
func (v RemoveContentProtectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":                validation.Validate(v.ConfigID, validation.Required),
		"Version":                 validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID":        validation.Validate(v.SecurityPolicyID, validation.Required),
		"ContentProtectionRuleID": validation.Validate(v.ContentProtectionRuleID, validation.Required),
	}.Filter()
}

func (b *botman) GetContentProtectionRule(ctx context.Context, params GetContentProtectionRuleRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetContentProtectionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-rules/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.ContentProtectionRuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetContentProtectionRule request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetContentProtectionRule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetContentProtectionRuleList(ctx context.Context, params GetContentProtectionRuleListRequest) (*GetContentProtectionRuleListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetContentProtectionRuleList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-rules",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetContentProtectionRuleList request: %w", err)
	}

	var result GetContentProtectionRuleListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetContentProtectionRuleList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetContentProtectionRuleListResponse
	if params.ContentProtectionRuleID != "" {
		for _, val := range result.ContentProtectionRules {
			if val["contentProtectionRuleId"].(string) == params.ContentProtectionRuleID {
				filteredResult.ContentProtectionRules = append(filteredResult.ContentProtectionRules, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateContentProtectionRule(ctx context.Context, params UpdateContentProtectionRuleRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateContentProtectionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-rules/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.ContentProtectionRuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateContentProtectionRule request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateContentProtectionRule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateContentProtectionRule(ctx context.Context, params CreateContentProtectionRuleRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateContentProtectionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-rules",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateContentProtectionRule request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateContentProtectionRule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveContentProtectionRule(ctx context.Context, params RemoveContentProtectionRuleRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveContentProtectionRule")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-rules/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.ContentProtectionRuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveContentProtectionRule request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveContentProtectionRule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
