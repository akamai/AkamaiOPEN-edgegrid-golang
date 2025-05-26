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
	// The CustomBotCategory interface supports creating, retrieving, modifying and removing custom bot categories for a
	// configuration.
	CustomBotCategory interface {
		// GetCustomBotCategoryList https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-categories
		GetCustomBotCategoryList(ctx context.Context, params GetCustomBotCategoryListRequest) (*GetCustomBotCategoryListResponse, error)

		// GetCustomBotCategory https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category
		GetCustomBotCategory(ctx context.Context, params GetCustomBotCategoryRequest) (map[string]interface{}, error)

		// CreateCustomBotCategory https://techdocs.akamai.com/bot-manager/reference/post-custom-bot-category
		CreateCustomBotCategory(ctx context.Context, params CreateCustomBotCategoryRequest) (map[string]interface{}, error)

		// UpdateCustomBotCategory https://techdocs.akamai.com/bot-manager/reference/put-custom-bot-category
		UpdateCustomBotCategory(ctx context.Context, params UpdateCustomBotCategoryRequest) (map[string]interface{}, error)

		// RemoveCustomBotCategory https://techdocs.akamai.com/bot-manager/reference/delete-custom-bot-category
		RemoveCustomBotCategory(ctx context.Context, params RemoveCustomBotCategoryRequest) error
	}

	// GetCustomBotCategoryListRequest is used to retrieve custom bot categories for a configuration.
	GetCustomBotCategoryListRequest struct {
		ConfigID   int64
		Version    int64
		CategoryID string
	}

	// GetCustomBotCategoryListResponse is used to retrieve custom bot categories for a configuration.
	GetCustomBotCategoryListResponse struct {
		Categories []map[string]interface{} `json:"categories"`
	}

	// GetCustomBotCategoryRequest is used to retrieve a specific custom bot category
	GetCustomBotCategoryRequest struct {
		ConfigID   int64
		Version    int64
		CategoryID string
	}

	// CreateCustomBotCategoryRequest is used to create a new custom bot category for a specific configuration.
	CreateCustomBotCategoryRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}

	// UpdateCustomBotCategoryRequest is used to update an existing custom bot category
	UpdateCustomBotCategoryRequest struct {
		ConfigID    int64
		Version     int64
		CategoryID  string
		JsonPayload json.RawMessage
	}

	// RemoveCustomBotCategoryRequest is used to remove an existing custom bot category
	RemoveCustomBotCategoryRequest struct {
		ConfigID   int64
		Version    int64
		CategoryID string
	}
)

// Validate validates a GetCustomBotCategoryRequest.
func (v GetCustomBotCategoryRequest) Validate() error {
	return validation.Errors{
		"ConfigID":   validation.Validate(v.ConfigID, validation.Required),
		"Version":    validation.Validate(v.Version, validation.Required),
		"CategoryID": validation.Validate(v.CategoryID, validation.Required),
	}.Filter()
}

// Validate validates a GetCustomBotCategoryListRequest.
func (v GetCustomBotCategoryListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateCustomBotCategoryRequest.
func (v CreateCustomBotCategoryRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomBotCategoryRequest.
func (v UpdateCustomBotCategoryRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"CategoryID":  validation.Validate(v.CategoryID, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveCustomBotCategoryRequest.
func (v RemoveCustomBotCategoryRequest) Validate() error {
	return validation.Errors{
		"ConfigID":   validation.Validate(v.ConfigID, validation.Required),
		"Version":    validation.Validate(v.Version, validation.Required),
		"CategoryID": validation.Validate(v.CategoryID, validation.Required),
	}.Filter()
}

func (b *botman) GetCustomBotCategory(ctx context.Context, params GetCustomBotCategoryRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomBotCategory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-bot-categories/%s",
		params.ConfigID,
		params.Version,
		params.CategoryID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomBotCategory request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomBotCategory request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetCustomBotCategoryList(ctx context.Context, params GetCustomBotCategoryListRequest) (*GetCustomBotCategoryListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomBotCategoryList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-bot-categories",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetlustomDenyList request: %w", err)
	}

	var result GetCustomBotCategoryListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomBotCategoryList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetCustomBotCategoryListResponse
	if params.CategoryID != "" {
		for _, val := range result.Categories {
			if val["categoryId"].(string) == params.CategoryID {
				filteredResult.Categories = append(filteredResult.Categories, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateCustomBotCategory(ctx context.Context, params UpdateCustomBotCategoryRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomBotCategory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-bot-categories/%s",
		params.ConfigID,
		params.Version,
		params.CategoryID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomBotCategory request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomBotCategory request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateCustomBotCategory(ctx context.Context, params CreateCustomBotCategoryRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateCustomBotCategory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-bot-categories",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateCustomBotCategory request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateCustomBotCategory request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveCustomBotCategory(ctx context.Context, params RemoveCustomBotCategoryRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveCustomBotCategory")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/custom-bot-categories/%s",
		params.ConfigID,
		params.Version,
		params.CategoryID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveCustomBotCategory request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveCustomBotCategory request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
