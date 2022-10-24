package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ReputationProtection interface supports retrieving and updating reputation protection for a configuration and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	ReputationProtection interface {
		// GetReputationProtections retrieves the current reputation protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		// Deprecated: this method will be removed in a future release. Use GetReputationProtection instead.
		GetReputationProtections(ctx context.Context, params GetReputationProtectionsRequest) (*GetReputationProtectionsResponse, error)

		// GetReputationProtection retrieves the current reputation protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetReputationProtection(ctx context.Context, params GetReputationProtectionRequest) (*GetReputationProtectionResponse, error)

		// UpdateReputationProtection updates the reputation protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateReputationProtection(ctx context.Context, params UpdateReputationProtectionRequest) (*UpdateReputationProtectionResponse, error)

		// RemoveReputationProtection removes reputation protection for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		// Deprecated: this method will be removed in a future release. Use UpdateReputationProtection instead.
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
	GetReputationProtectionResponse ProtectionsResponse

	// GetReputationProtectionsRequest is used to retrieve the reputation protection setting.
	// Deprecated: this struct will be removed in a future release.
	GetReputationProtectionsRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	// GetReputationProtectionsResponse is returned from a call to GetReputationProtection.
	// Deprecated: this struct will be removed in a future release.
	GetReputationProtectionsResponse ProtectionsResponse

	// UpdateReputationProtectionRequest is used to modify the reputation protection setting.
	UpdateReputationProtectionRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	// UpdateReputationProtectionResponse is returned from a call to UpdateReputationProtection.
	UpdateReputationProtectionResponse ProtectionsResponse

	// RemoveReputationProtectionRequest is used to remove the reputation protection settings.
	// Deprecated: this struct will be removed in a future release.
	RemoveReputationProtectionRequest struct {
		ConfigID                int    `json:"-"`
		Version                 int    `json:"-"`
		PolicyID                string `json:"-"`
		ApplyReputationControls bool   `json:"applyReputationControls"`
	}

	// RemoveReputationProtectionResponse is returned from a call to RemoveReputationProtection.
	// Deprecated: this struct will be removed in a future release.
	RemoveReputationProtectionResponse ProtectionsResponse
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
	logger := p.Log(ctx)
	logger.Debug("GetReputationProtection")

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
		return nil, fmt.Errorf("failed to create GetReputationProtection request: %w", err)
	}

	var result GetReputationProtectionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get reputation protection request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetReputationProtections(ctx context.Context, params GetReputationProtectionsRequest) (*GetReputationProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetReputationProtections")

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
		return nil, fmt.Errorf("failed to create GetReputationProtections request: %w", err)
	}

	var result GetReputationProtectionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get reputation protections request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateReputationProtection(ctx context.Context, params UpdateReputationProtectionRequest) (*UpdateReputationProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateReputationProtection")

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
		return nil, fmt.Errorf("failed to create UpdateReputationProtection request: %w", err)
	}

	var result UpdateReputationProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update reputation protection request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveReputationProtection(ctx context.Context, params RemoveReputationProtectionRequest) (*RemoveReputationProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveReputationProtection")

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
		return nil, fmt.Errorf("failed to create RemoveReputationProtection request: %w", err)
	}

	var result RemoveReputationProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove reputation protection request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
