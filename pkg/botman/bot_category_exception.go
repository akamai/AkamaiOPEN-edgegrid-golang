package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The BotCategoryException interface supports retrieving bot category exceptions
	BotCategoryException interface {
		// GetBotCategoryException https://techdocs.akamai.com/bot-manager/reference/get-bot-category-exception
		GetBotCategoryException(ctx context.Context, params GetBotCategoryExceptionRequest) (map[string]interface{}, error)

		// UpdateBotCategoryException https://techdocs.akamai.com/bot-manager/reference/put-bot-category-exception
		UpdateBotCategoryException(ctx context.Context, params UpdateBotCategoryExceptionRequest) (map[string]interface{}, error)
	}

	// GetBotCategoryExceptionRequest is used to retrieve bot category exceptions
	GetBotCategoryExceptionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
	}

	// UpdateBotCategoryExceptionRequest is used to update bot category exceptions
	UpdateBotCategoryExceptionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		JsonPayload      json.RawMessage
	}
)

// Validate validates a GetBotCategoryExceptionRequest.
func (v GetBotCategoryExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateBotCategoryExceptionRequest.
func (v UpdateBotCategoryExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetBotCategoryException(ctx context.Context, params GetBotCategoryExceptionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetBotCategoryException")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/bot-protection-exceptions",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotCategoryException request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotCategoryException request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateBotCategoryException(ctx context.Context, params UpdateBotCategoryExceptionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateBotCategoryException")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/bot-protection-exceptions",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateBotCategoryException request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateBotCategoryException request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
