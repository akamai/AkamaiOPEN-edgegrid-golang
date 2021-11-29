package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The IPGeoProtection interface supports retrieving and modifying the protections for a policy,
	// and whether each is enabled or disabled.
	//
	// Deprecated: this interface will be removed in a future release. Use the PolicyProtections interface instead.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	IPGeoProtection interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetIPGeoProtections(ctx context.Context, params GetIPGeoProtectionsRequest) (*GetIPGeoProtectionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetIPGeoProtection(ctx context.Context, params GetIPGeoProtectionRequest) (*GetIPGeoProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateIPGeoProtection(ctx context.Context, params UpdateIPGeoProtectionRequest) (*UpdateIPGeoProtectionResponse, error)
	}

	// GetIPGeoProtectionRequest is used to retrieve the IPGeo protection settings.
	GetIPGeoProtectionRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// GetIPGeoProtectionResponse is returned from a call to GetIPGeoProtection.
	GetIPGeoProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// GetIPGeoProtectionsRequest is used to retrieve the IPGeo protection settings.
	GetIPGeoProtectionsRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// GetIPGeoProtectionsResponse is returned from a call to GetIPGeoProtections.
	GetIPGeoProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// UpdateIPGeoProtectionResponse is used to modify the IPGeo protection settings.
	UpdateIPGeoProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}

	// UpdateIPGeoProtectionRequest is returned from a call to UpdateIPGeoProtection.
	UpdateIPGeoProtectionRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}
)

// Validate validates a GetIPGeoProtectionRequest.
func (v GetIPGeoProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetIPGeoProtectionsRequest.
func (v GetIPGeoProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateIPGeoProtectionRequest.
func (v UpdateIPGeoProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetIPGeoProtection(ctx context.Context, params GetIPGeoProtectionRequest) (*GetIPGeoProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetIPGeoProtection")

	var rval GetIPGeoProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetIPGeoProtection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetIPGeoProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetIPGeoProtections(ctx context.Context, params GetIPGeoProtectionsRequest) (*GetIPGeoProtectionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetIPGeoProtections")

	var rval GetIPGeoProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetIPGeoProtections request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetIPGeoProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateIPGeoProtection(ctx context.Context, params UpdateIPGeoProtectionRequest) (*UpdateIPGeoProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateIPGeoProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateIPGeoProtection request: %w", err)
	}

	var rval UpdateIPGeoProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateIPGeoProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
