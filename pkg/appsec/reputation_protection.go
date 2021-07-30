package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ReputationProtection interface supports retrieving, modifying and removing reputation
	// protection.
	//
	// Note: this interface is DEPRECATED and will be removed in a future release. Use the PolicyProtections interface instead.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	ReputationProtection interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetReputationProtections(ctx context.Context, params GetReputationProtectionsRequest) (*GetReputationProtectionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetReputationProtection(ctx context.Context, params GetReputationProtectionRequest) (*GetReputationProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateReputationProtection(ctx context.Context, params UpdateReputationProtectionRequest) (*UpdateReputationProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		RemoveReputationProtection(ctx context.Context, params RemoveReputationProtectionRequest) (*RemoveReputationProtectionResponse, error)
	}

	// GetReputationProtectionRequest is used to retrieve the reputation protection setting.
	GetReputationProtectionRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	// GetReputationProtectionResponse is returned from a call to GetReputationProtection.
	GetReputationProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// GetReputationProtectionsRequest is used to retrieve the reputation protection setting.
	GetReputationProtectionsRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	// GetReputationProtectionsResponse is returned from a call to GetReputationProtection.
	GetReputationProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	// UpdateReputationProtectionRequest is used to modify the reputation protection setting.
	UpdateReputationProtectionRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	// UpdateReputationProtectionResponse is returned from a call to UpdateReputationProtection.
	UpdateReputationProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	// RemoveReputationProtectionRequest is used to remove the reputation protection settings.
	RemoveReputationProtectionRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	// RemoveReputationProtectionResponse is returned from a call to RemoveReputationProtection.
	RemoveReputationProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates a GetReputationProtectionRequest.
func (v GetReputationProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetReputationProtectionsRequest.
func (v GetReputationProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateReputationProtectionRequest.
func (v UpdateReputationProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveReputationProtectionRequest.
func (v RemoveReputationProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
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
		return nil, fmt.Errorf("failed to create GetReputationProtection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetReputationProtection request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create GetReputationProtections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetReputationProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create UpdateReputationProtection request: %w", err)
	}

	var rval UpdateReputationProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateReputationProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *appsec) RemoveReputationProtection(ctx context.Context, params RemoveReputationProtectionRequest) (*RemoveReputationProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("RemoveReputationProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveReputationProtection request: %w", err)
	}

	var rval RemoveReputationProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveReputationProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
