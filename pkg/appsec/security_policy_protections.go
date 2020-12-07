package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// PolicyProtections represents a collection of PolicyProtections
//
// See: PolicyProtections.GetPolicyProtections()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// PolicyProtections  contains operations available on PolicyProtections  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpolicyprotections
	PolicyProtections interface {
		GetPolicyProtections(ctx context.Context, params GetPolicyProtectionsRequest) (*GetPolicyProtectionsResponse, error)
		UpdatePolicyProtections(ctx context.Context, params UpdatePolicyProtectionsRequest) (*UpdatePolicyProtectionsResponse, error)
		RemovePolicyProtections(ctx context.Context, params RemovePolicyProtectionsRequest) (*RemovePolicyProtectionsResponse, error)
	}

	GetPolicyProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	GetPolicyProtectionsRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyAPIConstraints           bool   `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool   `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool   `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool   `json:"applyRateControls"`
		ApplyReputationControls       bool   `json:"applyReputationControls"`
		ApplySlowPostControls         bool   `json:"applySlowPostControls"`
	}

	UpdatePolicyProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	UpdatePolicyProtectionsRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyAPIConstraints           bool   `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool   `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool   `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool   `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool   `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool   `json:"applySlowPostControls,omitempty"`
	}

	RemovePolicyProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	RemovePolicyProtectionsRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyAPIConstraints           bool   `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool   `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool   `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool   `json:"applyRateControls"`
		ApplyReputationControls       bool   `json:"applyReputationControls"`
		ApplySlowPostControls         bool   `json:"applySlowPostControls"`
	}
)

// Validate validates GetPolicyProtectionsRequest
func (v GetPolicyProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdatePolicyProtectionsRequest
func (v UpdatePolicyProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates RemovePolicyProtectionsRequest
func (v RemovePolicyProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetPolicyProtections(ctx context.Context, params GetPolicyProtectionsRequest) (*GetPolicyProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetPolicyProtections")

	var rval GetPolicyProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getpolicyprotections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getpolicyprotections  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a PolicyProtections.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpolicyprotections

func (p *appsec) UpdatePolicyProtections(ctx context.Context, params UpdatePolicyProtectionsRequest) (*UpdatePolicyProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdatePolicyProtections")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create PolicyProtectionsrequest: %w", err)
	}

	var rval UpdatePolicyProtectionsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create PolicyProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Update will update a PolicyProtections.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpolicyprotections

func (p *appsec) RemovePolicyProtections(ctx context.Context, params RemovePolicyProtectionsRequest) (*RemovePolicyProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("RemovePolicyProtections")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create PolicyProtectionsrequest: %w", err)
	}

	var rval RemovePolicyProtectionsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create PolicyProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
