package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The APIConstraintsProtection interface supports retrieving and updating API request constraints
	// for a configuration and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	APIConstraintsProtection interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapirequestconstraints
		GetAPIConstraintsProtection(ctx context.Context, params GetAPIConstraintsProtectionRequest) (*GetAPIConstraintsProtectionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putapirequestconstraints
		UpdateAPIConstraintsProtection(ctx context.Context, params UpdateAPIConstraintsProtectionRequest) (*UpdateAPIConstraintsProtectionResponse, error)
	}

	// GetAPIConstraintsProtectionRequest is used to retrieve the API constraints protection setting.
	GetAPIConstraintsProtectionRequest struct {
		ConfigID            int    `json:"-"`
		Version             int    `json:"-"`
		PolicyID            string `json:"-"`
		ApplyAPIConstraints bool   `json:"applyApiConstraints"`
	}

	// GetAPIConstraintsProtectionResponse is returned from a call to GetAPIConstraintsProtection.
	GetAPIConstraintsProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// UpdateAPIConstraintsProtectionRequest is used to modify the API constraints protection setting.
	UpdateAPIConstraintsProtectionRequest struct {
		ConfigID            int    `json:"-"`
		Version             int    `json:"-"`
		PolicyID            string `json:"-"`
		ApplyAPIConstraints bool   `json:"applyApiConstraints"`
	}

	// UpdateAPIConstraintsProtectionResponse is returned from a call to UpdateAPIConstraintsProtection.
	UpdateAPIConstraintsProtectionResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates a GetAPIConstraintsProtectionRequest.
func (v GetAPIConstraintsProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateAPIConstraintsProtectionRequest.
func (v UpdateAPIConstraintsProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetAPIConstraintsProtection(ctx context.Context, params GetAPIConstraintsProtectionRequest) (*GetAPIConstraintsProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAPIConstraintsProtection")

	var rval GetAPIConstraintsProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAPIConstraintsProtection request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetAPIConstraintsProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateAPIConstraintsProtection(ctx context.Context, params UpdateAPIConstraintsProtectionRequest) (*UpdateAPIConstraintsProtectionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAPIConstraintsProtection")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAPIConstraintsProtection request: %w", err)
	}

	var rval UpdateAPIConstraintsProtectionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateAPIConstraintsProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
