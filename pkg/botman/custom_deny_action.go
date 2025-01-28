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
	// The CustomDenyAction interface supports creating, retrieving, modifying and removing custom deny action for a
	// configuration.
	CustomDenyAction interface {
		// GetCustomDenyActionList https://techdocs.akamai.com/bot-manager/reference/get-custom-deny-actions
		GetCustomDenyActionList(ctx context.Context, params GetCustomDenyActionListRequest) (*GetCustomDenyActionListResponse, error)

		// GetCustomDenyAction https://techdocs.akamai.com/bot-manager/reference/get-custom-deny-action
		GetCustomDenyAction(ctx context.Context, params GetCustomDenyActionRequest) (map[string]interface{}, error)

		// CreateCustomDenyAction https://techdocs.akamai.com/bot-manager/reference/post-custom-deny-action
		CreateCustomDenyAction(ctx context.Context, params CreateCustomDenyActionRequest) (map[string]interface{}, error)

		// UpdateCustomDenyAction https://techdocs.akamai.com/bot-manager/reference/put-custom-deny-action
		UpdateCustomDenyAction(ctx context.Context, params UpdateCustomDenyActionRequest) (map[string]interface{}, error)

		// RemoveCustomDenyAction https://techdocs.akamai.com/bot-manager/reference/delete-custom-deny-action
		RemoveCustomDenyAction(ctx context.Context, params RemoveCustomDenyActionRequest) error
	}

	// GetCustomDenyActionListRequest is used to retrieve custom deny actions for a configuration.
	GetCustomDenyActionListRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// GetCustomDenyActionListResponse is used to retrieve custom deny actions for a configuration.
	GetCustomDenyActionListResponse struct {
		CustomDenyActions []map[string]interface{} `json:"customDenyActions"`
	}

	// GetCustomDenyActionRequest is used to retrieve a specific custom deny action
	GetCustomDenyActionRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// CreateCustomDenyActionRequest is used to create a new custom deny action for a specific configuration.
	CreateCustomDenyActionRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}

	// UpdateCustomDenyActionRequest is used to update an existing custom deny action
	UpdateCustomDenyActionRequest struct {
		ConfigID    int64
		Version     int64
		ActionID    string
		JsonPayload json.RawMessage
	}

	// RemoveCustomDenyActionRequest is used to remove an existing custom deny action
	RemoveCustomDenyActionRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}
)

// Validate validates a GetCustomDenyActionRequest.
func (v GetCustomDenyActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ActionID": validation.Validate(v.ActionID, validation.Required),
	}.Filter()
}

// Validate validates a GetCustomDenyActionListRequest.
func (v GetCustomDenyActionListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateCustomDenyActionRequest.
func (v CreateCustomDenyActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomDenyActionRequest.
func (v UpdateCustomDenyActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"ActionID":    validation.Validate(v.ActionID, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveCustomDenyActionRequest.
func (v RemoveCustomDenyActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ActionID": validation.Validate(v.ActionID, validation.Required),
	}.Filter()
}

func (b *botman) GetCustomDenyAction(ctx context.Context, params GetCustomDenyActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomDenyAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/custom-deny-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomDenyAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomDenyAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetCustomDenyActionList(ctx context.Context, params GetCustomDenyActionListRequest) (*GetCustomDenyActionListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomDenyActionList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/custom-deny-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetlustomDenyList request: %w", err)
	}

	var result GetCustomDenyActionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomDenyActionList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetCustomDenyActionListResponse
	if params.ActionID != "" {
		for _, val := range result.CustomDenyActions {
			if val["actionId"].(string) == params.ActionID {
				filteredResult.CustomDenyActions = append(filteredResult.CustomDenyActions, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateCustomDenyAction(ctx context.Context, params UpdateCustomDenyActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomDenyAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/custom-deny-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomDenyAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomDenyAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateCustomDenyAction(ctx context.Context, params CreateCustomDenyActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateCustomDenyAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/custom-deny-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateCustomDenyAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateCustomDenyAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveCustomDenyAction(ctx context.Context, params RemoveCustomDenyActionRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveCustomDenyAction")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/custom-deny-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveCustomDenyAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveCustomDenyAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
