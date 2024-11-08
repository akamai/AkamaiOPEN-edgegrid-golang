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
	// The CustomClient interface supports creating, retrieving, modifying and removing custom clients for a configuration.
	CustomClient interface {
		// GetCustomClientList https://techdocs.akamai.com/bot-manager/reference/get-custom-clients
		GetCustomClientList(ctx context.Context, params GetCustomClientListRequest) (*GetCustomClientListResponse, error)

		// GetCustomClient https://techdocs.akamai.com/bot-manager/reference/get-custom-client
		GetCustomClient(ctx context.Context, params GetCustomClientRequest) (map[string]interface{}, error)

		// CreateCustomClient https://techdocs.akamai.com/bot-manager/reference/post-custom-clients
		CreateCustomClient(ctx context.Context, params CreateCustomClientRequest) (map[string]interface{}, error)

		// UpdateCustomClient https://techdocs.akamai.com/bot-manager/reference/put-custom-clients
		UpdateCustomClient(ctx context.Context, params UpdateCustomClientRequest) (map[string]interface{}, error)

		// RemoveCustomClient https://techdocs.akamai.com/bot-manager/reference/delete-custom-clients
		RemoveCustomClient(ctx context.Context, params RemoveCustomClientRequest) error
	}

	// GetCustomClientListRequest is used to retrieve the custom clients for a configuration.
	GetCustomClientListRequest struct {
		ConfigID       int64
		Version        int64
		CustomClientID string
	}

	// GetCustomClientListResponse is used to retrieve the custom clients for a configuration.
	GetCustomClientListResponse struct {
		CustomClients []map[string]interface{} `json:"customClients"`
	}

	// GetCustomClientRequest is used to retrieve a specific custom client.
	GetCustomClientRequest struct {
		ConfigID       int64
		Version        int64
		CustomClientID string
	}

	// CreateCustomClientRequest is used to create a new custom client for a specific configuration.
	CreateCustomClientRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}

	// UpdateCustomClientRequest is used to update details for a specific custom client.
	UpdateCustomClientRequest struct {
		ConfigID       int64
		Version        int64
		CustomClientID string
		JsonPayload    json.RawMessage
	}

	// RemoveCustomClientRequest is used to remove an existing custom client.
	RemoveCustomClientRequest struct {
		ConfigID       int64
		Version        int64
		CustomClientID string
	}
)

// Validate validates a GetCustomClientRequest.
func (v GetCustomClientRequest) Validate() error {
	return validation.Errors{
		"ConfigID":       validation.Validate(v.ConfigID, validation.Required),
		"Version":        validation.Validate(v.Version, validation.Required),
		"CustomClientID": validation.Validate(v.CustomClientID, validation.Required),
	}.Filter()
}

// Validate validates a GetCustomClientsRequest.
func (v GetCustomClientListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateCustomClientRequest.
func (v CreateCustomClientRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomClientRequest.
func (v UpdateCustomClientRequest) Validate() error {
	return validation.Errors{
		"ConfigID":       validation.Validate(v.ConfigID, validation.Required),
		"Version":        validation.Validate(v.Version, validation.Required),
		"CustomClientID": validation.Validate(v.CustomClientID, validation.Required),
		"JsonPayload":    validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveCustomClientRequest.
func (v RemoveCustomClientRequest) Validate() error {
	return validation.Errors{
		"ConfigID":       validation.Validate(v.ConfigID, validation.Required),
		"Version":        validation.Validate(v.Version, validation.Required),
		"CustomClientID": validation.Validate(v.CustomClientID, validation.Required),
	}.Filter()
}

func (b *botman) GetCustomClient(ctx context.Context, params GetCustomClientRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomClient")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-clients/%s",
		params.ConfigID,
		params.Version,
		params.CustomClientID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomClient request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomClient request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetCustomClientList(ctx context.Context, params GetCustomClientListRequest) (*GetCustomClientListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomClientList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-clients",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomClientList request: %w", err)
	}

	var result GetCustomClientListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomClientList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetCustomClientListResponse
	if params.CustomClientID != "" {
		for _, val := range result.CustomClients {
			if val["customClientId"].(string) == params.CustomClientID {
				filteredResult.CustomClients = append(filteredResult.CustomClients, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateCustomClient(ctx context.Context, params UpdateCustomClientRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomClient")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-clients/%s",
		params.ConfigID,
		params.Version,
		params.CustomClientID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomClient request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomClient request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateCustomClient(ctx context.Context, params CreateCustomClientRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateCustomClient")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-clients",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateCustomClient request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateCustomClient request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveCustomClient(ctx context.Context, params RemoveCustomClientRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveCustomClient")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/custom-clients/%s",
		params.ConfigID,
		params.Version,
		params.CustomClientID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveCustomClient request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveCustomClient request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
