package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The NetworkLayerProtection interface supports retrieving, modifying and removing network layer protection.
	//
	// Deprecated: this interface will be removed in a future release. Use the PolicyProtections interface instead.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	NetworkLayerProtection interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetNetworkLayerProtections(ctx context.Context, params GetNetworkLayerProtectionsRequest) (*GetNetworkLayerProtectionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetNetworkLayerProtection(ctx context.Context, params GetNetworkLayerProtectionRequest) (*GetNetworkLayerProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateNetworkLayerProtection(ctx context.Context, params UpdateNetworkLayerProtectionRequest) (*UpdateNetworkLayerProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		RemoveNetworkLayerProtection(ctx context.Context, params RemoveNetworkLayerProtectionRequest) (*RemoveNetworkLayerProtectionResponse, error)
	}

	// GetNetworkLayerProtectionRequest is used to retrieve the network layer protection setting.
	GetNetworkLayerProtectionRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// GetNetworkLayerProtectionResponse is returned from a call to GetNetworkLayerProtection.
	GetNetworkLayerProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// GetNetworkLayerProtectionsRequest is used to retrieve the network layer protection setting.
	GetNetworkLayerProtectionsRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// GetNetworkLayerProtectionsResponse is returned from a call to GetNetworkLayerProtection.
	GetNetworkLayerProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	// UpdateNetworkLayerProtectionRequest is used to modify the network layer protection setting.
	UpdateNetworkLayerProtectionRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// UpdateNetworkLayerProtectionResponse is returned from a call to UpdateNetworkLayerProtection
	UpdateNetworkLayerProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	// RemoveNetworkLayerProtectionRequest is used to remove the network layer protection setting.
	RemoveNetworkLayerProtectionRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// RemoveNetworkLayerProtectionResponse is returned from a call to RemoveNetworkLayerProtection.
	RemoveNetworkLayerProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates a GetNetworkLayerProtectionRequest.
func (v GetNetworkLayerProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetNetworkLayerProtectionsRequest.
func (v GetNetworkLayerProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateNetworkLayerProtectionRequest.
func (v UpdateNetworkLayerProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveNetworkLayerProtectionRequest.
func (v RemoveNetworkLayerProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetNetworkLayerProtection(ctx context.Context, params GetNetworkLayerProtectionRequest) (*GetNetworkLayerProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetNetworkLayerProtection")

	var rval GetNetworkLayerProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetNetworkLayerProtection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetNetworkLayerProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetNetworkLayerProtections(ctx context.Context, params GetNetworkLayerProtectionsRequest) (*GetNetworkLayerProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetNetworkLayerProtections")

	var rval GetNetworkLayerProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetNetworkLayerProtections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetNetworkLayerProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateNetworkLayerProtection(ctx context.Context, params UpdateNetworkLayerProtectionRequest) (*UpdateNetworkLayerProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateNetworkLayerProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateNetworkLayerProtection request: %w", err)
	}

	var rval UpdateNetworkLayerProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateNetworkLayerProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *appsec) RemoveNetworkLayerProtection(ctx context.Context, params RemoveNetworkLayerProtectionRequest) (*RemoveNetworkLayerProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("RemoveNetworkLayerProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveNetworkLayerProtection request: %w", err)
	}

	var rval RemoveNetworkLayerProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveNetworkLayerProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
