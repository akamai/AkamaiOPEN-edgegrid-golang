package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// WAFProtection represents a collection of WAFProtection
//
// See: WAFProtection.GetWAFProtection()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// WAFProtection  contains operations available on WAFProtection  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getwafprotection
	WAFProtection interface {
		GetWAFProtections(ctx context.Context, params GetWAFProtectionsRequest) (*GetWAFProtectionsResponse, error)
		GetWAFProtection(ctx context.Context, params GetWAFProtectionRequest) (*GetWAFProtectionResponse, error)
		UpdateWAFProtection(ctx context.Context, params UpdateWAFProtectionRequest) (*UpdateWAFProtectionResponse, error)
	}

	GetWAFProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	GetWAFProtectionRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}

	GetWAFProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	GetWAFProtectionsRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}

	UpdateWAFProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	UpdateWAFProtectionRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}
)

// Validate validates GetWAFProtectionRequest
func (v GetWAFProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetWAFProtectionsRequest
func (v GetWAFProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateWAFProtectionRequest
func (v UpdateWAFProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetWAFProtection(ctx context.Context, params GetWAFProtectionRequest) (*GetWAFProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetWAFProtection")

	var rval GetWAFProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getwafprotection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getwafprotection  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetWAFProtections(ctx context.Context, params GetWAFProtectionsRequest) (*GetWAFProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetWAFProtections")

	var rval GetWAFProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getwafprotections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getwafprotections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a WAFProtection.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putwafprotection

func (p *appsec) UpdateWAFProtection(ctx context.Context, params UpdateWAFProtectionRequest) (*UpdateWAFProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateWAFProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create WAFProtectionrequest: %w", err)
	}

	var rval UpdateWAFProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create WAFProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
