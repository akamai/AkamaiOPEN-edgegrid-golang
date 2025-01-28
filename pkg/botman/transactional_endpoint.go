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
	// The TransactionalEndpoint interface supports creating, retrieving, modifying and removing transactional endpoints
	// for a configuration.
	TransactionalEndpoint interface {
		// GetTransactionalEndpointList https://techdocs.akamai.com/bot-manager/reference/get-transactional-endpoints
		GetTransactionalEndpointList(ctx context.Context, params GetTransactionalEndpointListRequest) (*GetTransactionalEndpointListResponse, error)

		// GetTransactionalEndpoint https://techdocs.akamai.com/bot-manager/reference/get-transactional-endpoint
		GetTransactionalEndpoint(ctx context.Context, params GetTransactionalEndpointRequest) (map[string]interface{}, error)

		// CreateTransactionalEndpoint https://techdocs.akamai.com/bot-manager/reference/post-transactional-endpoint
		CreateTransactionalEndpoint(ctx context.Context, params CreateTransactionalEndpointRequest) (map[string]interface{}, error)

		// UpdateTransactionalEndpoint https://techdocs.akamai.com/bot-manager/reference/put-transactional-endpoint
		UpdateTransactionalEndpoint(ctx context.Context, params UpdateTransactionalEndpointRequest) (map[string]interface{}, error)

		// RemoveTransactionalEndpoint https://techdocs.akamai.com/bot-manager/reference/delete-transactional-endpoint
		RemoveTransactionalEndpoint(ctx context.Context, params RemoveTransactionalEndpointRequest) error
	}

	// GetTransactionalEndpointListRequest is used to retrieve the transactional endpoint for a configuration.
	GetTransactionalEndpointListRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		OperationID      string
	}

	// GetTransactionalEndpointListResponse is used to retrieve the transactional endpoints for a configuration.
	GetTransactionalEndpointListResponse struct {
		Operations []map[string]interface{} `json:"operations"`
	}

	// GetTransactionalEndpointRequest is used to retrieve a specific transactional endpoint.
	GetTransactionalEndpointRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		OperationID      string
	}

	// CreateTransactionalEndpointRequest is used to create a new transactional endpoint for a specific configuration.
	CreateTransactionalEndpointRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		JsonPayload      json.RawMessage
	}

	// UpdateTransactionalEndpointRequest is used to update details for a specific transactional endpoint
	UpdateTransactionalEndpointRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		OperationID      string
		JsonPayload      json.RawMessage
	}

	// RemoveTransactionalEndpointRequest is used to remove an existing transactional endpoint
	RemoveTransactionalEndpointRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		OperationID      string
	}
)

// Validate validates a GetTransactionalEndpointRequest.
func (v GetTransactionalEndpointRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"OperationID":      validation.Validate(v.OperationID, validation.Required),
	}.Filter()
}

// Validate validates a GetTransactionalEndpointListRequest.
func (v GetTransactionalEndpointListRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateTransactionalEndpointRequest.
func (v CreateTransactionalEndpointRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateTransactionalEndpointRequest.
func (v UpdateTransactionalEndpointRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"OperationID":      validation.Validate(v.OperationID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveTransactionalEndpointRequest.
func (v RemoveTransactionalEndpointRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"OperationID":      validation.Validate(v.OperationID, validation.Required),
	}.Filter()
}

func (b *botman) GetTransactionalEndpoint(ctx context.Context, params GetTransactionalEndpointRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetTransactionalEndpoint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/bot-protection/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.OperationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTransactionalEndpoint request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTransactionalEndpoint request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetTransactionalEndpointList(ctx context.Context, params GetTransactionalEndpointListRequest) (*GetTransactionalEndpointListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetTransactionalEndpointList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/bot-protection",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetlustomDenyList request: %w", err)
	}

	var result GetTransactionalEndpointListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTransactionalEndpointList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetTransactionalEndpointListResponse
	if params.OperationID != "" {
		for _, val := range result.Operations {
			if val["operationId"].(string) == params.OperationID {
				filteredResult.Operations = append(filteredResult.Operations, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateTransactionalEndpoint(ctx context.Context, params UpdateTransactionalEndpointRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateTransactionalEndpoint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/bot-protection/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.OperationID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateTransactionalEndpoint request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateTransactionalEndpoint request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateTransactionalEndpoint(ctx context.Context, params CreateTransactionalEndpointRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateTransactionalEndpoint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/bot-protection",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateTransactionalEndpoint request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateTransactionalEndpoint request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveTransactionalEndpoint(ctx context.Context, params RemoveTransactionalEndpointRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveTransactionalEndpoint")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/security-policies/%s/transactional-endpoints/bot-protection/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.OperationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveTransactionalEndpoint request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveTransactionalEndpoint request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
