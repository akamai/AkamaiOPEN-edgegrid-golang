package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomBotCategoryAction interface supports retrieving and updating the actions for the custom bot categories of
	// a configuration
	CustomBotCategoryAction interface {
		// GetCustomBotCategoryActionList https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-actions
		GetCustomBotCategoryActionList(ctx context.Context, params GetCustomBotCategoryActionListRequest) (*GetCustomBotCategoryActionListResponse, error)

		// GetCustomBotCategoryAction https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-action
		GetCustomBotCategoryAction(ctx context.Context, params GetCustomBotCategoryActionRequest) (map[string]interface{}, error)

		// UpdateCustomBotCategoryAction https://techdocs.akamai.com/bot-manager/reference/put-custom-bot-category-action
		UpdateCustomBotCategoryAction(ctx context.Context, params UpdateCustomBotCategoryActionRequest) (map[string]interface{}, error)
	}

	// GetCustomBotCategoryActionListRequest is used to retrieve the custom bot category actions for a policy.
	GetCustomBotCategoryActionListRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		CategoryID       string
	}

	// GetCustomBotCategoryActionListResponse is returned from a call to GetCustomBotCategoryActionList.
	GetCustomBotCategoryActionListResponse struct {
		Actions []map[string]interface{} `json:"actions"`
	}

	// GetCustomBotCategoryActionRequest is used to retrieve the action for a custom bot category
	GetCustomBotCategoryActionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		CategoryID       string
	}

	// UpdateCustomBotCategoryActionRequest is used to modify an existing custom bot category action
	UpdateCustomBotCategoryActionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		CategoryID       string
		JsonPayload      json.RawMessage
	}
)

// Validate validates a GetCustomBotCategoryActionRequest.
func (v GetCustomBotCategoryActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"CategoryID":       validation.Validate(v.CategoryID, validation.Required),
	}.Filter()
}

// Validate validates a GetCustomBotCategoryActionListRequest.
func (v GetCustomBotCategoryActionListRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomBotCategoryActionRequest.
func (v UpdateCustomBotCategoryActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"CategoryID":       validation.Validate(v.CategoryID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetCustomBotCategoryAction(ctx context.Context, params GetCustomBotCategoryActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomBotCategoryAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-bot-category-actions/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.CategoryID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomBotCategoryAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomBotCategoryAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetCustomBotCategoryActionList(ctx context.Context, params GetCustomBotCategoryActionListRequest) (*GetCustomBotCategoryActionListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomBotCategoryActionList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-bot-category-actions",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomBotCategoryActionList request: %w", err)
	}

	var result GetCustomBotCategoryActionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomBotCategoryActionList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetCustomBotCategoryActionListResponse
	if params.CategoryID != "" {
		for _, val := range result.Actions {
			if val["categoryId"].(string) == params.CategoryID {
				filteredResult.Actions = append(filteredResult.Actions, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateCustomBotCategoryAction(ctx context.Context, params UpdateCustomBotCategoryActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomBotCategoryAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-bot-category-actions/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.CategoryID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomBotCategoryAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomBotCategoryAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
