package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RateProtection interface supports retrieving and updating rate controls for a configuration
	// and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	RateProtection interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetRateProtections(ctx context.Context, params GetRateProtectionsRequest) (*GetRateProtectionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetRateProtection(ctx context.Context, params GetRateProtectionRequest) (*GetRateProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateRateProtection(ctx context.Context, params UpdateRateProtectionRequest) (*UpdateRateProtectionResponse, error)
	}

	// GetRateProtectionRequest is used to retrieve the rate protection setting.
	GetRateProtectionRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}

	// GetRateProtectionResponse is returned from a call to GetRateProtection.
	GetRateProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// GetRateProtectionsRequest is used to retrieve the rate protection setting.
	GetRateProtectionsRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}

	// GetRateProtectionsResponse is returned from a call to GetRateProtection.
	GetRateProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// UpdateRateProtectionRequest is used to modify the rate protection setting.
	UpdateRateProtectionRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}

	// UpdateRateProtectionResponse is returned from a call to UpdateRateProtection.
	UpdateRateProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates a GetRateProtectionRequest.
func (v GetRateProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetRateProtectionsRequest.
func (v GetRateProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateRateProtectionRequest.
func (v UpdateRateProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetRateProtection(ctx context.Context, params GetRateProtectionRequest) (*GetRateProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRateProtection")

	var rval GetRateProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRateProtection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetRateProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetRateProtections(ctx context.Context, params GetRateProtectionsRequest) (*GetRateProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRateProtections")

	var rval GetRateProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRateProtections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetRateProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateRateProtection(ctx context.Context, params UpdateRateProtectionRequest) (*UpdateRateProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRateProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRateProtection request: %w", err)
	}

	var rval UpdateRateProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateRateProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
