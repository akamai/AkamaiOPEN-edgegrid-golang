package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The WAFProtection interface supports retrieving, modifying and removing protections for a
	// security policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	WAFProtection interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetWAFProtections(ctx context.Context, params GetWAFProtectionsRequest) (*GetWAFProtectionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetWAFProtection(ctx context.Context, params GetWAFProtectionRequest) (*GetWAFProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateWAFProtection(ctx context.Context, params UpdateWAFProtectionRequest) (*UpdateWAFProtectionResponse, error)
	}

	// GetWAFProtectionRequest is used to retrieve the WAF protection setting.
	GetWAFProtectionRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}

	// GetWAFProtectionResponse is returned from a call to GetWAFProtection.
	GetWAFProtectionResponse ProtectionsResponse

	// GetWAFProtectionsRequest is used to retrieve the WAF protection setting.
	GetWAFProtectionsRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}

	// GetWAFProtectionsResponse is returned from a call to GetWAFProtections.
	GetWAFProtectionsResponse ProtectionsResponse

	// UpdateWAFProtectionRequest is used to modify the WAF protection setting.
	UpdateWAFProtectionRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}

	// UpdateWAFProtectionResponse is returned from a call to UpdateWAFProtection.
	UpdateWAFProtectionResponse ProtectionsResponse
)

// Validate validates a GetWAFProtectionRequest.
func (v GetWAFProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetWAFProtectionsRequest.
func (v GetWAFProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateWAFProtectionRequest.
func (v UpdateWAFProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetWAFProtection(ctx context.Context, params GetWAFProtectionRequest) (*GetWAFProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetWAFProtection")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetWAFProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetWAFProtection request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetWAFProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) GetWAFProtections(ctx context.Context, params GetWAFProtectionsRequest) (*GetWAFProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetWAFProtections")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetWAFProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetWAFProtections request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetWAFProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) UpdateWAFProtection(ctx context.Context, params UpdateWAFProtectionRequest) (*UpdateWAFProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateWAFProtection")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateWAFProtection request: %w", err)
	}

	var result UpdateWAFProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateWAFProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
