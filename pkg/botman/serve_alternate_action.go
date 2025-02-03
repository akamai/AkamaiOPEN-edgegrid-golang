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
	// The ServeAlternateAction interface supports creating, retrieving, modifying and removing serve alternate action for a
	// configuration.
	ServeAlternateAction interface {
		// GetServeAlternateActionList https://techdocs.akamai.com/bot-manager/reference/get-serve-alternate-actions
		GetServeAlternateActionList(ctx context.Context, params GetServeAlternateActionListRequest) (*GetServeAlternateActionListResponse, error)

		// GetServeAlternateAction https://techdocs.akamai.com/bot-manager/reference/get-serve-alternate-action
		GetServeAlternateAction(ctx context.Context, params GetServeAlternateActionRequest) (map[string]interface{}, error)

		// CreateServeAlternateAction https://techdocs.akamai.com/bot-manager/reference/post-serve-alternate-action
		CreateServeAlternateAction(ctx context.Context, params CreateServeAlternateActionRequest) (map[string]interface{}, error)

		// UpdateServeAlternateAction https://techdocs.akamai.com/bot-manager/reference/put-serve-alternate-action
		UpdateServeAlternateAction(ctx context.Context, params UpdateServeAlternateActionRequest) (map[string]interface{}, error)

		// RemoveServeAlternateAction https://techdocs.akamai.com/bot-manager/reference/delete-serve-alternate-action
		RemoveServeAlternateAction(ctx context.Context, params RemoveServeAlternateActionRequest) error
	}

	// GetServeAlternateActionListRequest is used to retrieve serve alternate actions for a configuration.
	GetServeAlternateActionListRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// GetServeAlternateActionListResponse is used to retrieve serve alternate actions for a configuration.
	GetServeAlternateActionListResponse struct {
		ServeAlternateActions []map[string]interface{} `json:"serveAlternateActions"`
	}

	// GetServeAlternateActionRequest is used to retrieve a specific serve alternate action
	GetServeAlternateActionRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// CreateServeAlternateActionRequest is used to create a new serve alternate action for a specific configuration.
	CreateServeAlternateActionRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}

	// UpdateServeAlternateActionRequest is used to update an existing serve alternate action
	UpdateServeAlternateActionRequest struct {
		ConfigID    int64
		Version     int64
		ActionID    string
		JsonPayload json.RawMessage
	}

	// RemoveServeAlternateActionRequest is used to remove an existing serve alternate action
	RemoveServeAlternateActionRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}
)

// Validate validates a GetServeAlternateActionRequest.
func (v GetServeAlternateActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ActionID": validation.Validate(v.ActionID, validation.Required),
	}.Filter()
}

// Validate validates a GetServeAlternateActionListRequest.
func (v GetServeAlternateActionListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateServeAlternateActionRequest.
func (v CreateServeAlternateActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateServeAlternateActionRequest.
func (v UpdateServeAlternateActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"ActionID":    validation.Validate(v.ActionID, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveServeAlternateActionRequest.
func (v RemoveServeAlternateActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ActionID": validation.Validate(v.ActionID, validation.Required),
	}.Filter()
}

func (b *botman) GetServeAlternateAction(ctx context.Context, params GetServeAlternateActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetServeAlternateAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/serve-alternate-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetServeAlternateAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetServeAlternateAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetServeAlternateActionList(ctx context.Context, params GetServeAlternateActionListRequest) (*GetServeAlternateActionListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetServeAlternateActionList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/serve-alternate-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetlustomDenyList request: %w", err)
	}

	var result GetServeAlternateActionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetServeAlternateActionList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetServeAlternateActionListResponse
	if params.ActionID != "" {
		for _, val := range result.ServeAlternateActions {
			if val["actionId"].(string) == params.ActionID {
				filteredResult.ServeAlternateActions = append(filteredResult.ServeAlternateActions, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateServeAlternateAction(ctx context.Context, params UpdateServeAlternateActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateServeAlternateAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/serve-alternate-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateServeAlternateAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateServeAlternateAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateServeAlternateAction(ctx context.Context, params CreateServeAlternateActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateServeAlternateAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/serve-alternate-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateServeAlternateAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateServeAlternateAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveServeAlternateAction(ctx context.Context, params RemoveServeAlternateActionRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveServeAlternateAction")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/response-actions/serve-alternate-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveServeAlternateAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveServeAlternateAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
