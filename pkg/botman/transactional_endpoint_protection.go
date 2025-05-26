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
	// The TransactionalEndpointProtection interface supports retrieving and updating transactional endpoint protection settings for a configuration.
	TransactionalEndpointProtection interface {

		// GetTransactionalEndpointProtection https://techdocs.akamai.com/bot-manager/reference/get-transactional-endpoint-protection
		GetTransactionalEndpointProtection(ctx context.Context, params GetTransactionalEndpointProtectionRequest) (map[string]interface{}, error)

		// UpdateTransactionalEndpointProtection https://techdocs.akamai.com/bot-manager/reference/put-transactional-endpoint-protection
		UpdateTransactionalEndpointProtection(ctx context.Context, params UpdateTransactionalEndpointProtectionRequest) (map[string]interface{}, error)
	}

	// GetTransactionalEndpointProtectionRequest is used to retrieve transactional endpoint protection settings
	GetTransactionalEndpointProtectionRequest struct {
		ConfigID int64
		Version  int64
	}

	// UpdateTransactionalEndpointProtectionRequest is used to modify transactional endpoint protection settings
	UpdateTransactionalEndpointProtectionRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}
)

// Validate validates a GetTransactionalEndpointProtectionRequest.
func (v GetTransactionalEndpointProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateTransactionalEndpointProtectionRequest.
func (v UpdateTransactionalEndpointProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetTransactionalEndpointProtection(ctx context.Context, params GetTransactionalEndpointProtectionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetTransactionalEndpointProtection")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/transactional-endpoint-protection",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTransactionalEndpointProtection request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTransactionalEndpointProtection request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateTransactionalEndpointProtection(ctx context.Context, params UpdateTransactionalEndpointProtectionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateTransactionalEndpointProtection")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/transactional-endpoint-protection",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateTransactionalEndpointProtection request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateTransactionalEndpointProtection request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
