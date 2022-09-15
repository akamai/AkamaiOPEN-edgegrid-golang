package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The NetworkLayerProtection interface supports retrieving and updating network layer protection for a configuration and policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#protections
	NetworkLayerProtection interface {
		// GetNetworkLayerProtections retrieves the current network layer protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		// Deprecated: this method will be removed in a future release. Use GetNetworkLayerProtection instead.
		GetNetworkLayerProtections(ctx context.Context, params GetNetworkLayerProtectionsRequest) (*GetNetworkLayerProtectionsResponse, error)

		// GetNetworkLayerProtection retrieves the current network layer protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprotections
		GetNetworkLayerProtection(ctx context.Context, params GetNetworkLayerProtectionRequest) (*GetNetworkLayerProtectionResponse, error)

		// UpdateNetworkLayerProtection updates the network layer protection setting for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		UpdateNetworkLayerProtection(ctx context.Context, params UpdateNetworkLayerProtectionRequest) (*UpdateNetworkLayerProtectionResponse, error)

		// UpdateNetworkLayerProtection removes network layer protection for a configuration and policy.
		//
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprotections
		// Deprecated: this method will be removed in a future release. Use UpdateNetworkLayerProtection instead.
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
	GetNetworkLayerProtectionResponse ProtectionsResponse

	// GetNetworkLayerProtectionsRequest is used to retrieve the network layer protection setting.
	// Deprecated: this struct will be removed in a future release.
	GetNetworkLayerProtectionsRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// GetNetworkLayerProtectionsResponse is returned from a call to GetNetworkLayerProtection.
	// Deprecated: this struct will be removed in a future release.
	GetNetworkLayerProtectionsResponse ProtectionsResponse

	// UpdateNetworkLayerProtectionRequest is used to modify the network layer protection setting.
	UpdateNetworkLayerProtectionRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// UpdateNetworkLayerProtectionResponse is returned from a call to UpdateNetworkLayerProtection
	UpdateNetworkLayerProtectionResponse ProtectionsResponse

	// RemoveNetworkLayerProtectionRequest is used to remove the network layer protection setting.
	// Deprecated: this struct will be removed in a future release.
	RemoveNetworkLayerProtectionRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// RemoveNetworkLayerProtectionResponse is returned from a call to RemoveNetworkLayerProtection.
	// Deprecated: this struct will be removed in a future release.
	RemoveNetworkLayerProtectionResponse ProtectionsResponse
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
	logger := p.Log(ctx)
	logger.Debug("GetNetworkLayerProtection")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetNetworkLayerProtectionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetNetworkLayerProtection request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetNetworkLayerProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) GetNetworkLayerProtections(ctx context.Context, params GetNetworkLayerProtectionsRequest) (*GetNetworkLayerProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetNetworkLayerProtections")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetNetworkLayerProtectionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/protections",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetNetworkLayerProtections request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetNetworkLayerProtections request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) UpdateNetworkLayerProtection(ctx context.Context, params UpdateNetworkLayerProtectionRequest) (*UpdateNetworkLayerProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateNetworkLayerProtection")

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
		return nil, fmt.Errorf("failed to create UpdateNetworkLayerProtection request: %w", err)
	}

	var result UpdateNetworkLayerProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateNetworkLayerProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveNetworkLayerProtection(ctx context.Context, params RemoveNetworkLayerProtectionRequest) (*RemoveNetworkLayerProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveNetworkLayerProtection")

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
		return nil, fmt.Errorf("failed to create RemoveNetworkLayerProtection request: %w", err)
	}

	var result RemoveNetworkLayerProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveNetworkLayerProtection request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
