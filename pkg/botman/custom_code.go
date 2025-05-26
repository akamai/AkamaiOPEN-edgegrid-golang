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
	// The CustomCode interface supports retrieving and updating custom code
	CustomCode interface {
		GetCustomCode(ctx context.Context, params GetCustomCodeRequest) (map[string]interface{}, error)

		UpdateCustomCode(ctx context.Context, params UpdateCustomCodeRequest) (map[string]interface{}, error)
	}

	// GetCustomCodeRequest is used to retrieve custom code
	GetCustomCodeRequest struct {
		ConfigID int64
		Version  int64
	}

	// UpdateCustomCodeRequest is used to modify custom code
	UpdateCustomCodeRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}
)

// Validate validates a GetCustomCodeRequest.
func (v GetCustomCodeRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomCodeRequest.
func (v UpdateCustomCodeRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetCustomCode(ctx context.Context, params GetCustomCodeRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/transactional-endpoint-protection/custom-code",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomCode request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomCode request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateCustomCode(ctx context.Context, params UpdateCustomCodeRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/transactional-endpoint-protection/custom-code",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomCode request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomCode request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, b.Error(resp)
	}

	return result, nil
}
