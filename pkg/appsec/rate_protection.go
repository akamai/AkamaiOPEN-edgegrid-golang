package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RateProtection interface supports retrieving and updating rate protection for a configuration and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	RateProtection interface {
		// GetRateProtections retrieves the current rate protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		// Deprecated: this method will be removed in a future release. Use GetRateProtection instead.
		GetRateProtections(ctx context.Context, params GetRateProtectionsRequest) (*GetRateProtectionsResponse, error)

		// GetRateProtection retrieves the current rate protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetRateProtection(ctx context.Context, params GetRateProtectionRequest) (*GetRateProtectionResponse, error)

		// UpdateRateProtection updates the rate protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateRateProtection(ctx context.Context, params UpdateRateProtectionRequest) (*UpdateRateProtectionResponse, error)
	}

	// GetRateProtectionRequest is used to retrieve the rate protection setting.
	GetRateProtectionRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}

	// GetRateProtectionResponse is returned from a call to GetRateProtection.
	GetRateProtectionResponse ProtectionsResponse

	// GetRateProtectionsRequest is used to retrieve the rate protection setting.
	// Deprecated: this struct will be removed in a future release.
	GetRateProtectionsRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}

	// GetRateProtectionsResponse is returned from a call to GetRateProtection.
	// Deprecated: this struct will be removed in a future release.
	GetRateProtectionsResponse ProtectionsResponse

	// UpdateRateProtectionRequest is used to modify the rate protection setting.
	UpdateRateProtectionRequest struct {
		ConfigID          int    `json:"-"`
		Version           int    `json:"-"`
		PolicyID          string `json:"-"`
		ApplyRateControls bool   `json:"applyRateControls"`
	}

	// UpdateRateProtectionResponse is returned from a call to UpdateRateProtection.
	UpdateRateProtectionResponse ProtectionsResponse
)

// Validate validates a GetRateProtectionRequest.
func (v GetRateProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetRateProtectionsRequest.
func (v GetRateProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateRateProtectionRequest.
func (v UpdateRateProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetRateProtection(ctx context.Context, params GetRateProtectionRequest) (*GetRateProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRateProtection")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetRateProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRateProtection request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRateProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) GetRateProtections(ctx context.Context, params GetRateProtectionsRequest) (*GetRateProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRateProtections")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetRateProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRateProtections request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRateProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) UpdateRateProtection(ctx context.Context, params UpdateRateProtectionRequest) (*UpdateRateProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRateProtection")

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
		return nil, fmt.Errorf("failed to create UpdateRateProtection request: %w", err)
	}

	var result UpdateRateProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateRateProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
