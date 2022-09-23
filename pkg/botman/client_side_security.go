package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ClientSideSecurity interface supports retrieving and updating client side security settings
	ClientSideSecurity interface {

		// GetClientSideSecurity https://techdocs.akamai.com/bot-manager/reference/get-client-side-security
		GetClientSideSecurity(ctx context.Context, params GetClientSideSecurityRequest) (map[string]interface{}, error)

		// UpdateClientSideSecurity https://techdocs.akamai.com/bot-manager/reference/put-client-side-security
		UpdateClientSideSecurity(ctx context.Context, params UpdateClientSideSecurityRequest) (map[string]interface{}, error)
	}

	// GetClientSideSecurityRequest is used to retrieve client side security settings
	GetClientSideSecurityRequest struct {
		ConfigID int64 `json:"configId"`
		Version  int64 `json:"version"`
	}

	// UpdateClientSideSecurityRequest is used to modify client side security settings
	UpdateClientSideSecurityRequest struct {
		ConfigID    int64           `json:"-"`
		Version     int64           `json:"-"`
		JsonPayload json.RawMessage `json:"-"`
	}
)

// Validate validates a GetClientSideSecurityRequest.
func (v GetClientSideSecurityRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateClientSideSecurityRequest.
func (v UpdateClientSideSecurityRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetClientSideSecurity(ctx context.Context, params GetClientSideSecurityRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetClientSideSecurity")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/client-side-security",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetClientSideSecurity request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetClientSideSecurity request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) UpdateClientSideSecurity(ctx context.Context, params UpdateClientSideSecurityRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateClientSideSecurity")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/client-side-security",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateClientSideSecurity request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateClientSideSecurity request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, b.Error(resp)
	}

	return result, nil
}
