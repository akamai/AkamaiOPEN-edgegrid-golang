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
	// The ConditionalAction interface supports creating, retrieving, modifying and removing conditional action for a
	// configuration.
	ConditionalAction interface {
		// GetConditionalActionList https://techdocs.akamai.com/bot-manager/reference/get-conditional-actions
		GetConditionalActionList(ctx context.Context, params GetConditionalActionListRequest) (*GetConditionalActionListResponse, error)

		// GetConditionalAction https://techdocs.akamai.com/bot-manager/reference/get-conditional-action
		GetConditionalAction(ctx context.Context, params GetConditionalActionRequest) (map[string]interface{}, error)

		// CreateConditionalAction https://techdocs.akamai.com/bot-manager/reference/post-conditional-action
		CreateConditionalAction(ctx context.Context, params CreateConditionalActionRequest) (map[string]interface{}, error)

		// UpdateConditionalAction https://techdocs.akamai.com/bot-manager/reference/put-conditional-action
		UpdateConditionalAction(ctx context.Context, params UpdateConditionalActionRequest) (map[string]interface{}, error)

		// RemoveConditionalAction https://techdocs.akamai.com/bot-manager/reference/delete-conditional-action
		RemoveConditionalAction(ctx context.Context, params RemoveConditionalActionRequest) error
	}

	// GetConditionalActionListRequest is used to retrieve conditional actions for a configuration.
	GetConditionalActionListRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// GetConditionalActionListResponse is used to retrieve conditional actions for a configuration.
	GetConditionalActionListResponse struct {
		ConditionalActions []map[string]interface{} `json:"conditionalActions"`
	}

	// GetConditionalActionRequest is used to retrieve a specific conditional action
	GetConditionalActionRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// CreateConditionalActionRequest is used to create a new conditional action for a specific configuration.
	CreateConditionalActionRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}

	// UpdateConditionalActionRequest is used to update an existing conditional action
	UpdateConditionalActionRequest struct {
		ConfigID    int64
		Version     int64
		ActionID    string
		JsonPayload json.RawMessage
	}

	// RemoveConditionalActionRequest is used to remove an existing conditional action
	RemoveConditionalActionRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}
)

// Validate validates a GetConditionalActionRequest.
func (v GetConditionalActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ActionID": validation.Validate(v.ActionID, validation.Required),
	}.Filter()
}

// Validate validates a GetConditionalActionListRequest.
func (v GetConditionalActionListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateConditionalActionRequest.
func (v CreateConditionalActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateConditionalActionRequest.
func (v UpdateConditionalActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"ActionID":    validation.Validate(v.ActionID, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveConditionalActionRequest.
func (v RemoveConditionalActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ActionID": validation.Validate(v.ActionID, validation.Required),
	}.Filter()
}

func (b *botman) GetConditionalAction(ctx context.Context, params GetConditionalActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetConditionalAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/conditional-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetConditionalAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetConditionalAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetConditionalActionList(ctx context.Context, params GetConditionalActionListRequest) (*GetConditionalActionListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetConditionalActionList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/conditional-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetlustomDenyList request: %w", err)
	}

	var result GetConditionalActionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetConditionalActionList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetConditionalActionListResponse
	if params.ActionID != "" {
		for _, val := range result.ConditionalActions {
			if val["actionId"].(string) == params.ActionID {
				filteredResult.ConditionalActions = append(filteredResult.ConditionalActions, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateConditionalAction(ctx context.Context, params UpdateConditionalActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateConditionalAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/conditional-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateConditionalAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateConditionalAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateConditionalAction(ctx context.Context, params CreateConditionalActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateConditionalAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/conditional-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateConditionalAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateConditionalAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveConditionalAction(ctx context.Context, params RemoveConditionalActionRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveConditionalAction")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/response-actions/conditional-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveConditionalAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveConditionalAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
