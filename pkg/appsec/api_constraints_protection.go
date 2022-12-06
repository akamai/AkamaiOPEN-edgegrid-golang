package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ApiConstraintsProtection interface supports retrieving and updating API request constraints protection for a configuration and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	ApiConstraintsProtection interface {
		// GetAPIConstraintsProtection retrieves the current API constraints protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapirequestconstraints
		GetAPIConstraintsProtection(ctx context.Context, params GetAPIConstraintsProtectionRequest) (*GetAPIConstraintsProtectionResponse, error)

		// UpdateAPIConstraintsProtection updates the API constraints protection setting for a configuration and policy.
		//
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

	// ProtectionsResponse is returned from a call to GetAPIConstraintsProtection and similar security policy protection requests.
	ProtectionsResponse struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyMalwareControls          bool `json:"applyMalwareControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// GetAPIConstraintsProtectionResponse contains the status of various protection flags assigned to a security policy.
	GetAPIConstraintsProtectionResponse ProtectionsResponse

	// UpdateAPIConstraintsProtectionRequest is used to modify the API constraints protection setting.
	UpdateAPIConstraintsProtectionRequest struct {
		ConfigID            int    `json:"-"`
		Version             int    `json:"-"`
		PolicyID            string `json:"-"`
		ApplyAPIConstraints bool   `json:"applyApiConstraints"`
	}

	// UpdateAPIConstraintsProtectionResponse is returned from a call to UpdateAPIConstraintsProtection.
	UpdateAPIConstraintsProtectionResponse ProtectionsResponse
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
	logger := p.Log(ctx)
	logger.Debug("GetAPIConstraintsProtection")

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
		return nil, fmt.Errorf("failed to create GetAPIConstraintsProtection request: %w", err)
	}

	var result GetAPIConstraintsProtectionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get API constraints protection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateAPIConstraintsProtection(ctx context.Context, params UpdateAPIConstraintsProtectionRequest) (*UpdateAPIConstraintsProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAPIConstraintsProtection")

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
		return nil, fmt.Errorf("failed to create UpdateAPIConstraintsProtection request: %w", err)
	}

	var result UpdateAPIConstraintsProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update API constraints protection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
