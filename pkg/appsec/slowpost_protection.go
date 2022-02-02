package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SlowPostProtection interface supports retrieving and updating slow post protection
	// for a configuration and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#slowpostprotection
	SlowPostProtection interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetSlowPostProtections(ctx context.Context, params GetSlowPostProtectionsRequest) (*GetSlowPostProtectionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetSlowPostProtection(ctx context.Context, params GetSlowPostProtectionRequest) (*GetSlowPostProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateSlowPostProtection(ctx context.Context, params UpdateSlowPostProtectionRequest) (*UpdateSlowPostProtectionResponse, error)
	}

	// GetSlowPostProtectionRequest is used to retrieve the slow post protecton setting for a policy.
	GetSlowPostProtectionRequest struct {
		ConfigID              int    `json:"-"`
		Version               int    `json:"-"`
		PolicyID              string `json:"-"`
		ApplySlowPostControls bool   `json:"applySlowPostControls"`
	}

	// GetSlowPostProtectionResponse is returned from a call to GetSlowPostProtection.
	GetSlowPostProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// GetSlowPostProtectionsRequest is used to retrieve the slow post protecton setting for a policy.
	GetSlowPostProtectionsRequest struct {
		ConfigID              int    `json:"-"`
		Version               int    `json:"-"`
		PolicyID              string `json:"-"`
		ApplySlowPostControls bool   `json:"applySlowPostControls"`
	}

	// GetSlowPostProtectionsResponse is returned from a call to GetSlowPostProtections.
	GetSlowPostProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// UpdateSlowPostProtectionRequest is used to modify the slow post protection setting.
	UpdateSlowPostProtectionRequest struct {
		ConfigID              int    `json:"-"`
		Version               int    `json:"-"`
		PolicyID              string `json:"-"`
		ApplySlowPostControls bool   `json:"applySlowPostControls"`
	}

	// UpdateSlowPostProtectionResponse is returned from a call to UpdateSlowPostProtection.
	UpdateSlowPostProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates a GetSlowPostProtectionRequest.
func (v GetSlowPostProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetSlowPostProtectionsRequest.
func (v GetSlowPostProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateSlowPostProtectionRequest.
func (v UpdateSlowPostProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetSlowPostProtection(ctx context.Context, params GetSlowPostProtectionRequest) (*GetSlowPostProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetSlowPostProtection")

	var rval GetSlowPostProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSlowPostProtection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetSlowPostProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetSlowPostProtections(ctx context.Context, params GetSlowPostProtectionsRequest) (*GetSlowPostProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetSlowPostProtections")

	var rval GetSlowPostProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSlowPostProtections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetSlowPostProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateSlowPostProtection(ctx context.Context, params UpdateSlowPostProtectionRequest) (*UpdateSlowPostProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateSlowPostProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateSlowPostProtection request: %w", err)
	}

	var rval UpdateSlowPostProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSlowPostProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
