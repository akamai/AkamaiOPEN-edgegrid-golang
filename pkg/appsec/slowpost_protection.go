package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SlowPostProtection interface supports retrieving and updating slow post protection for a configuration and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#slowpostprotection
	SlowPostProtection interface {
		// GetSlowPostProtections retrieves the current SLOW post protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		// Deprecated: this method will be removed in a future release. Use GetSlowPostProtection instead.
		GetSlowPostProtections(ctx context.Context, params GetSlowPostProtectionsRequest) (*GetSlowPostProtectionsResponse, error)

		// GetSlowPostProtection retrieves the current SLOW post protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetSlowPostProtection(ctx context.Context, params GetSlowPostProtectionRequest) (*GetSlowPostProtectionResponse, error)

		// UpdateSlowPostProtection updates the SLOW post protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateSlowPostProtection(ctx context.Context, params UpdateSlowPostProtectionRequest) (*UpdateSlowPostProtectionResponse, error)
	}

	// GetSlowPostProtectionRequest is used to retrieve the slow post protecton setting for a policy.
	GetSlowPostProtectionRequest struct {
		ConfigID              int    `json:"-"`
		Version               int    `json:"-"`
		PolicyID              string `json:"-"`
		ApplySlowPostControls bool   `json:"applySlowPostControls"`
	}

	// GetSlowPostProtectionResponse is returned from a call to GetSlowPostProtection.
	GetSlowPostProtectionResponse ProtectionsResponse

	// GetSlowPostProtectionsRequest is used to retrieve the slow post protecton setting for a policy.
	// Deprecated: this struct will be removed in a future release.
	GetSlowPostProtectionsRequest struct {
		ConfigID              int    `json:"-"`
		Version               int    `json:"-"`
		PolicyID              string `json:"-"`
		ApplySlowPostControls bool   `json:"applySlowPostControls"`
	}

	// GetSlowPostProtectionsResponse is returned from a call to GetSlowPostProtections.
	// Deprecated: this struct will be removed in a future release.
	GetSlowPostProtectionsResponse ProtectionsResponse

	// UpdateSlowPostProtectionRequest is used to modify the slow post protection setting.
	UpdateSlowPostProtectionRequest struct {
		ConfigID              int    `json:"-"`
		Version               int    `json:"-"`
		PolicyID              string `json:"-"`
		ApplySlowPostControls bool   `json:"applySlowPostControls"`
	}

	// UpdateSlowPostProtectionResponse is returned from a call to UpdateSlowPostProtection.
	UpdateSlowPostProtectionResponse ProtectionsResponse
)

// Validate validates a GetSlowPostProtectionRequest.
func (v GetSlowPostProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetSlowPostProtectionsRequest.
func (v GetSlowPostProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateSlowPostProtectionRequest.
func (v UpdateSlowPostProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetSlowPostProtection(ctx context.Context, params GetSlowPostProtectionRequest) (*GetSlowPostProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSlowPostProtection")

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
		return nil, fmt.Errorf("failed to create GetSlowPostProtection request: %w", err)
	}

	var result GetSlowPostProtectionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get slow post protection request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetSlowPostProtections(ctx context.Context, params GetSlowPostProtectionsRequest) (*GetSlowPostProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSlowPostProtections")

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
		return nil, fmt.Errorf("failed to create GetSlowPostProtections request: %w", err)
	}

	var result GetSlowPostProtectionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get slow post protections request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateSlowPostProtection(ctx context.Context, params UpdateSlowPostProtectionRequest) (*UpdateSlowPostProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateSlowPostProtection")

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
		return nil, fmt.Errorf("failed to create UpdateSlowPostProtection request: %w", err)
	}

	var result UpdateSlowPostProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update slow post protection request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
