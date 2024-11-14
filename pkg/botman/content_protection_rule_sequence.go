package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ContentProtectionRuleSequence interface supports retrieving and updating content protection rule sequence
	ContentProtectionRuleSequence interface {
		// GetContentProtectionRuleSequence https://techdocs.akamai.com/content-protector/reference/get-content-protection-rule-sequence
		GetContentProtectionRuleSequence(ctx context.Context, params GetContentProtectionRuleSequenceRequest) (*GetContentProtectionRuleSequenceResponse, error)

		// UpdateContentProtectionRuleSequence https://techdocs.akamai.com/content-protector/reference/put-content-protection-rule-sequence
		UpdateContentProtectionRuleSequence(ctx context.Context, params UpdateContentProtectionRuleSequenceRequest) (*UpdateContentProtectionRuleSequenceResponse, error)
	}

	// GetContentProtectionRuleSequenceRequest is used to retrieve content protection rule sequence
	GetContentProtectionRuleSequenceRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
	}

	// GetContentProtectionRuleSequenceResponse contains the sequence of content protection rule
	GetContentProtectionRuleSequenceResponse ContentProtectionRuleUUIDSequence

	// UpdateContentProtectionRuleSequenceResponse contains the sequence of content protection rule
	UpdateContentProtectionRuleSequenceResponse ContentProtectionRuleUUIDSequence

	// UpdateContentProtectionRuleSequenceRequest is used to update content protection rule sequence
	UpdateContentProtectionRuleSequenceRequest struct {
		ConfigID                      int64
		Version                       int64
		SecurityPolicyID              string
		ContentProtectionRuleSequence ContentProtectionRuleUUIDSequence
	}

	// ContentProtectionRuleUUIDSequence is a sequence of UUIDs.
	ContentProtectionRuleUUIDSequence struct {
		ContentProtectionRuleSequence []string `json:"contentProtectionRuleSequence"`
	}
)

// Validate validates a GetContentProtectionRuleSequenceRequest.
func (v GetContentProtectionRuleSequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateContentProtectionRuleSequenceRequest.
func (v UpdateContentProtectionRuleSequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID":                      validation.Validate(v.ConfigID, validation.Required),
		"Version":                       validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID":              validation.Validate(v.SecurityPolicyID, validation.Required),
		"ContentProtectionRuleSequence": validation.Validate(v.ContentProtectionRuleSequence.ContentProtectionRuleSequence, validation.Required),
	}.Filter()
}

func (b *botman) GetContentProtectionRuleSequence(ctx context.Context, params GetContentProtectionRuleSequenceRequest) (*GetContentProtectionRuleSequenceResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetContentProtectionRuleSequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-rule-sequence",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetContentProtectionRuleSequence request: %w", err)
	}

	var result GetContentProtectionRuleSequenceResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetContentProtectionRuleSequence request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}

func (b *botman) UpdateContentProtectionRuleSequence(ctx context.Context, params UpdateContentProtectionRuleSequenceRequest) (*UpdateContentProtectionRuleSequenceResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateContentProtectionRuleSequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/content-protection-rule-sequence",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateContentProtectionRuleSequence request: %w", err)
	}

	var result UpdateContentProtectionRuleSequenceResponse
	resp, err := b.Exec(req, &result, params.ContentProtectionRuleSequence)
	if err != nil {
		return nil, fmt.Errorf("UpdateContentProtectionRuleSequence request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}
