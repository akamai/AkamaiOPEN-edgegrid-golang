package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The PolicyProtections interface supports retrieving and updating protections for a configuration and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	PolicyProtections interface {
		// GetPolicyProtections retrieves the current protection settings for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetPolicyProtections(ctx context.Context, params GetPolicyProtectionsRequest) (*PolicyProtectionsResponse, error)

		// UpdatePolicyProtections updates the protection settings for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdatePolicyProtections(ctx context.Context, params UpdatePolicyProtectionsRequest) (*PolicyProtectionsResponse, error)

		// RemovePolicyProtections removes protection settings for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		// Deprecated: this method will be removed in a future release. Use UpdatePolicyProtections instead.
		RemovePolicyProtections(ctx context.Context, params UpdatePolicyProtectionsRequest) (*PolicyProtectionsResponse, error)
	}

	// GetPolicyProtectionsRequest is used to retrieve the policy protection setting.
	GetPolicyProtectionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
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
		ApplyMalwareControls          bool   `json:"applyMalwareControls"`
	}

	// RemovePolicyProtectionsRequest is used to remove the policy protection setting.
	// Deprecated: this struct will be removed in a future release.
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
		ApplyMalwareControls          bool   `json:"applyMalwareControls"`
	}

	// PolicyProtectionsResponse is returned from GetPolicyProtections, UpdatePolicyProtections, and RemovePolicyProtections.
	PolicyProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
		ApplyMalwareControls          bool `json:"applyMalwareControls"`
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

func (p *appsec) GetPolicyProtections(ctx context.Context, params GetPolicyProtectionsRequest) (*PolicyProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetPolicyProtections")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetPolicyProtections request: %w", err)
	}

	var result PolicyProtectionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get policy protections request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdatePolicyProtections(ctx context.Context, params UpdatePolicyProtectionsRequest) (*PolicyProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdatePolicyProtections")

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
		return nil, fmt.Errorf("failed to create UpdatePolicyProtections request: %w", err)
	}

	var result PolicyProtectionsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update policy protections request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemovePolicyProtections(ctx context.Context, params UpdatePolicyProtectionsRequest) (*PolicyProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemovePolicyProtections")

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
		return nil, fmt.Errorf("failed to create RemovePolicyProtections request: %w", err)
	}

	var result PolicyProtectionsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove policy protections request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
