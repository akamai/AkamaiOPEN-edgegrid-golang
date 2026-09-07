package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsJA4Fingerprint interface supports retrieving, updating or removing settings
	// related to JA4 TLS Fingerprint.
	AdvancedSettingsJA4Fingerprint interface {
		// GetAdvancedSettingsJA4Fingerprint lists the JA4 TLS Fingerprint settings for a configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-ja4-fingerprint-settings
		GetAdvancedSettingsJA4Fingerprint(ctx context.Context, params GetAdvancedSettingsJA4FingerprintRequest) (*GetAdvancedSettingsJA4FingerprintResponse, error)

		// UpdateAdvancedSettingsJA4Fingerprint enables, disables, or updates the JA4 TLS Fingerprint settings for a
		// configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-ja4-fingerprint-settings
		UpdateAdvancedSettingsJA4Fingerprint(ctx context.Context, params UpdateAdvancedSettingsJA4FingerprintRequest) (*UpdateAdvancedSettingsJA4FingerprintResponse, error)

		// RemoveAdvancedSettingsJA4Fingerprint disables JA4 TLS Fingerprint for a configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-ja4-fingerprint-settings
		RemoveAdvancedSettingsJA4Fingerprint(ctx context.Context, params RemoveAdvancedSettingsJA4FingerprintRequest) (*RemoveAdvancedSettingsJA4FingerprintResponse, error)
	}

	// GetAdvancedSettingsJA4FingerprintRequest is used to retrieve the JA4 TLS Fingerprint settings for a configuration.
	GetAdvancedSettingsJA4FingerprintRequest struct {
		ConfigID int
		Version  int
	}

	// GetAdvancedSettingsJA4FingerprintResponse is returned from a call to GetAdvancedSettingsJA4Fingerprint.
	GetAdvancedSettingsJA4FingerprintResponse struct {
		HeaderNames []string `json:"headerNames"`
	}

	// UpdateAdvancedSettingsJA4FingerprintRequest is used to update the JA4 TLS Fingerprint settings for a configuration.
	UpdateAdvancedSettingsJA4FingerprintRequest struct {
		ConfigID    int
		Version     int
		HeaderNames []string `json:"headerNames,omitempty"`
	}

	// UpdateAdvancedSettingsJA4FingerprintResponse is returned from a call to UpdateAdvancedSettingsJA4Fingerprint.
	UpdateAdvancedSettingsJA4FingerprintResponse struct {
		HeaderNames []string `json:"headerNames"`
	}

	// RemoveAdvancedSettingsJA4FingerprintRequest is used to disable JA4 TLS Fingerprint for a configuration.
	RemoveAdvancedSettingsJA4FingerprintRequest struct {
		ConfigID int
		Version  int
	}

	// RemoveAdvancedSettingsJA4FingerprintResponse is returned from a call to RemoveAdvancedSettingsJA4Fingerprint.
	RemoveAdvancedSettingsJA4FingerprintResponse struct {
		HeaderNames []string `json:"headerNames"`
	}
)

// Validate validates a GetAdvancedSettingsJA4FingerprintRequest.
func (v GetAdvancedSettingsJA4FingerprintRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates an UpdateAdvancedSettingsJA4FingerprintRequest.
func (v UpdateAdvancedSettingsJA4FingerprintRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates a RemoveAdvancedSettingsJA4FingerprintRequest.
func (v RemoveAdvancedSettingsJA4FingerprintRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

func (p *appsec) GetAdvancedSettingsJA4Fingerprint(ctx context.Context, params GetAdvancedSettingsJA4FingerprintRequest) (*GetAdvancedSettingsJA4FingerprintResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsJA4Fingerprint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/ja4-fingerprint",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAdvancedSettingsJA4Fingerprint request: %w", err)
	}

	var result GetAdvancedSettingsJA4FingerprintResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get advanced settings JA4 Fingerprint failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateAdvancedSettingsJA4Fingerprint(ctx context.Context, params UpdateAdvancedSettingsJA4FingerprintRequest) (*UpdateAdvancedSettingsJA4FingerprintResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsJA4Fingerprint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/ja4-fingerprint",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsJA4Fingerprint request: %w", err)
	}

	var result UpdateAdvancedSettingsJA4FingerprintResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update advanced settings JA4 Fingerprint failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveAdvancedSettingsJA4Fingerprint(ctx context.Context, params RemoveAdvancedSettingsJA4FingerprintRequest) (*RemoveAdvancedSettingsJA4FingerprintResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveAdvancedSettingsJA4Fingerprint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/ja4-fingerprint",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveAdvancedSettingsJA4Fingerprint request: %w", err)
	}

	request := UpdateAdvancedSettingsJA4FingerprintRequest{
		ConfigID:    params.ConfigID,
		Version:     params.Version,
		HeaderNames: nil,
	}

	var result RemoveAdvancedSettingsJA4FingerprintResponse
	resp, err := p.Exec(req, &result, request)
	if err != nil {
		return nil, fmt.Errorf("remove advanced settings JA4 Fingerprint failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
