package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The PolicyProtections interface supports retrieving, modifying and removing protections for a
	// security policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	PolicyProtections interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetPolicyProtections(ctx context.Context, params GetPolicyProtectionsRequest) (*GetPolicyProtectionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdatePolicyProtections(ctx context.Context, params UpdatePolicyProtectionsRequest) (*UpdatePolicyProtectionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		RemovePolicyProtections(ctx context.Context, params RemovePolicyProtectionsRequest) (*RemovePolicyProtectionsResponse, error)
	}

	// GetPolicyProtectionsRequest is used to retrieve the policy protection setting.
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

	// GetPolicyProtectionsResponse is returned from a call to GetPolicyProtections.
	GetPolicyProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	// UpdatePolicyProtectionsRequest is used to modify the policy protection setting.
	UpdatePolicyProtectionsRequest struct {
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

	// UpdatePolicyProtectionsResponse is returned from a call to UpdatePolicyProtections.
	UpdatePolicyProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	// RemovePolicyProtectionsRequest is used to remove the policy protection setting.
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

	// RemovePolicyProtectionsResponse is returned from a call to RemovePolicyProtections.
	RemovePolicyProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates a GetPolicyProtectionsRequest.
func (v GetPolicyProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdatePolicyProtectionsRequest.
func (v UpdatePolicyProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemovePolicyProtectionsRequest.
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
		return nil, fmt.Errorf("failed to create GetPolicyProtections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetPolicyProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create UpdatePolicyProtections request: %w", err)
	}

	var rval UpdatePolicyProtectionsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdatePolicyProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

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
		return nil, fmt.Errorf("failed to create RemovePolicyProtections request: %w", err)
	}

	var rval RemovePolicyProtectionsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("RemovePolicyProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
