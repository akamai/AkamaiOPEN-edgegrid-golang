package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AkamaiBotCategoryAction interface supports retrieving and updating the actions for the akamai bot categories of
	// a configuration
	AkamaiBotCategoryAction interface {
		// GetAkamaiBotCategoryActionList https://techdocs.akamai.com/bot-manager/reference/get-akamai-bot-category-actions
		GetAkamaiBotCategoryActionList(ctx context.Context, params GetAkamaiBotCategoryActionListRequest) (*GetAkamaiBotCategoryActionListResponse, error)

		// GetAkamaiBotCategoryAction https://techdocs.akamai.com/bot-manager/reference/get-akamai-bot-category-action
		GetAkamaiBotCategoryAction(ctx context.Context, params GetAkamaiBotCategoryActionRequest) (map[string]interface{}, error)

		// UpdateAkamaiBotCategoryAction https://techdocs.akamai.com/bot-manager/reference/put-akamai-bot-category-action
		UpdateAkamaiBotCategoryAction(ctx context.Context, params UpdateAkamaiBotCategoryActionRequest) (map[string]interface{}, error)
	}

	// GetAkamaiBotCategoryActionListRequest is used to retrieve the akamai bot category actions for a policy.
	GetAkamaiBotCategoryActionListRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		CategoryID       string
	}

	// GetAkamaiBotCategoryActionListResponse is returned from a call to GetAkamaiBotCategoryActionList.
	GetAkamaiBotCategoryActionListResponse struct {
		Actions []map[string]interface{} `json:"actions"`
	}

	// GetAkamaiBotCategoryActionRequest is used to retrieve the action for an akamai bot category.
	GetAkamaiBotCategoryActionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		CategoryID       string
	}

	// UpdateAkamaiBotCategoryActionRequest is used to modify an akamai bot category action.
	UpdateAkamaiBotCategoryActionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		CategoryID       string
		JsonPayload      json.RawMessage
	}
)

// Validate validates a GetAkamaiBotCategoryActionRequest.
func (v GetAkamaiBotCategoryActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"CategoryID":       validation.Validate(v.CategoryID, validation.Required),
	}.Filter()
}

// Validate validates a GetAkamaiBotCategoryActionListRequest.
func (v GetAkamaiBotCategoryActionListRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateAkamaiBotCategoryActionRequest.
func (v UpdateAkamaiBotCategoryActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"CategoryID":       validation.Validate(v.CategoryID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetAkamaiBotCategoryAction(ctx context.Context, params GetAkamaiBotCategoryActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetAkamaiBotCategoryAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/akamai-bot-category-actions/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.CategoryID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAkamaiBotCategoryAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAkamaiBotCategoryAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetAkamaiBotCategoryActionList(ctx context.Context, params GetAkamaiBotCategoryActionListRequest) (*GetAkamaiBotCategoryActionListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetAkamaiBotCategoryActionList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/akamai-bot-category-actions",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAkamaiBotCategoryActionList request: %w", err)
	}

	var result GetAkamaiBotCategoryActionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAkamaiBotCategoryActionList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetAkamaiBotCategoryActionListResponse
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

func (b *botman) UpdateAkamaiBotCategoryAction(ctx context.Context, params UpdateAkamaiBotCategoryActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateAkamaiBotCategoryAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/akamai-bot-category-actions/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.CategoryID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAkamaiBotCategoryAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateAkamaiBotCategoryAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
