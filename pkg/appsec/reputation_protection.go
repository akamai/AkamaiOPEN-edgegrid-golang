package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ReputationProtection represents a collection of ReputationProtection
//
// See: ReputationProtection.GetReputationProtection()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ReputationProtection  contains operations available on ReputationProtection  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getreputationprotection
	ReputationProtection interface {
		GetReputationProtections(ctx context.Context, params GetReputationProtectionsRequest) (*GetReputationProtectionsResponse, error)
		GetReputationProtection(ctx context.Context, params GetReputationProtectionRequest) (*GetReputationProtectionResponse, error)
		UpdateReputationProtection(ctx context.Context, params UpdateReputationProtectionRequest) (*UpdateReputationProtectionResponse, error)
	}

	GetReputationProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	GetReputationProtectionRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	GetReputationProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	GetReputationProtectionsRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	UpdateReputationProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	UpdateReputationProtectionRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}
)

// Validate validates GetReputationProtectionRequest
func (v GetReputationProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetReputationProtectionsRequest
func (v GetReputationProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateReputationProtectionRequest
func (v UpdateReputationProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		//"ApplyReputationControls": validation.Validate(v.ApplyReputationControls, validation.Required),
	}.Filter()
}

func (p *appsec) GetReputationProtection(ctx context.Context, params GetReputationProtectionRequest) (*GetReputationProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetReputationProtection")

	var rval GetReputationProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getreputationprotection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getreputationprotection  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetReputationProtections(ctx context.Context, params GetReputationProtectionsRequest) (*GetReputationProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetReputationProtections")

	var rval GetReputationProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getreputationprotections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getreputationprotections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a ReputationProtection.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putreputationprotection

func (p *appsec) UpdateReputationProtection(ctx context.Context, params UpdateReputationProtectionRequest) (*UpdateReputationProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateReputationProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create ReputationProtectionrequest: %w", err)
	}

	var rval UpdateReputationProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create ReputationProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
