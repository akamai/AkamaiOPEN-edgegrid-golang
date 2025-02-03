package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The IPGeoProtection interface supports retrieving and updating IPGeo protection for a configuration and policy.
	// Deprecated: this interface will be removed in a future release. Use the SecurityPolicy interface instead.
	IPGeoProtection interface {
		// GetIPGeoProtections retrieves the current IPGeo protection protection setting for a configuration and policy.
		// Deprecated: this method will be removed in a future release. Use the GetPolicyProtections method of the PolicyProtections interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-protections
		GetIPGeoProtections(ctx context.Context, params GetIPGeoProtectionsRequest) (*GetIPGeoProtectionsResponse, error)

		// GetIPGeoProtection retrieves the current IPGeo protection protection setting for a configuration and policy.
		// Deprecated: this method will be removed in a future release. Use the GetPolicyProtections method of the PolicyProtections interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-protections
		GetIPGeoProtection(ctx context.Context, params GetIPGeoProtectionRequest) (*GetIPGeoProtectionResponse, error)

		// UpdateIPGeoProtection updates the IPGeo protection protection setting for a configuration and policy.
		// Deprecated: this method will be removed in a future release. Use the CreateSecurityPolicyWithDefaultProtections method of the SecurityPolicy interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-policy-protections
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
	GetIPGeoProtectionResponse ProtectionsResponse

	// GetIPGeoProtectionsRequest is used to retrieve the IPGeo protection settings.
	// Deprecated: this struct will be removed in a future release.
	GetIPGeoProtectionsRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// GetIPGeoProtectionsResponse is returned from a call to GetIPGeoProtections.
	// Deprecated: this struct will be removed in a future release.
	GetIPGeoProtectionsResponse ProtectionsResponse

	// UpdateIPGeoProtectionRequest is used to modify the IPGeo protection settings.
	UpdateIPGeoProtectionRequest struct {
		ConfigID                  int    `json:"-"`
		Version                   int    `json:"-"`
		PolicyID                  string `json:"-"`
		ApplyNetworkLayerControls bool   `json:"applyNetworkLayerControls"`
	}

	// UpdateIPGeoProtectionResponse is returned from a call to UpdateIPGeoProtection.
	UpdateIPGeoProtectionResponse ProtectionsResponse
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
	logger := p.Log(ctx)
	logger.Debug("GetIPGeoProtection")

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
		return nil, fmt.Errorf("failed to create GetIPGeoProtection request: %w", err)
	}

	var result GetIPGeoProtectionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get IPGeo protection request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetIPGeoProtections(ctx context.Context, params GetIPGeoProtectionsRequest) (*GetIPGeoProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetIPGeoProtections")

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
		return nil, fmt.Errorf("failed to create GetIPGeoProtections request: %w", err)
	}

	var result GetIPGeoProtectionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get IPGeo protections request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateIPGeoProtection(ctx context.Context, params UpdateIPGeoProtectionRequest) (*UpdateIPGeoProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateIPGeoProtection")

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
		return nil, fmt.Errorf("failed to create UpdateIPGeoProtection request: %w", err)
	}

	var result UpdateIPGeoProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update IPGeo protection request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
