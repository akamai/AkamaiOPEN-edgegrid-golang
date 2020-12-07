package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RateProtection represents a collection of RateProtection
//
// See: RateProtection.GetRateProtection()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// RateProtection  contains operations available on RateProtection  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getrateprotection
	RateProtection interface {
		GetRateProtections(ctx context.Context, params GetRateProtectionsRequest) (*GetRateProtectionsResponse, error)
		GetRateProtection(ctx context.Context, params GetRateProtectionRequest) (*GetRateProtectionResponse, error)
		UpdateRateProtection(ctx context.Context, params UpdateRateProtectionRequest) (*UpdateRateProtectionResponse, error)
	}

	GetRateProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	GetRateProtectionRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}

	GetRateProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	GetRateProtectionsRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}

	UpdateRateProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	UpdateRateProtectionRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}
)

// Validate validates GetRateProtectionRequest
func (v GetRateProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetRateProtectionsRequest
func (v GetRateProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRateProtectionRequest
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
		return nil, fmt.Errorf("failed to create getrateprotection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getrateprotection  request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create getrateprotections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getrateprotections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a RateProtection.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putrateprotection

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
		return nil, fmt.Errorf("failed to create create RateProtectionrequest: %w", err)
	}

	var rval UpdateRateProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create RateProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
