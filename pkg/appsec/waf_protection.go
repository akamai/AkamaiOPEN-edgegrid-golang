package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The WAFProtection interface supports retrieving and updating application layer protection for a configuration and policy.
	// Deprecated: this interface will be removed in a future release. Use the SecurityPolicy interface instead.
	WAFProtection interface {
		// GetWAFProtections retrieves the current WAF protection setting for a configuration and policy.
		// Deprecated: this method will be removed in a future release. Use the GetPolicyProtections method of the PolicyProtections interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-protections
		GetWAFProtections(ctx context.Context, params GetWAFProtectionsRequest) (*GetWAFProtectionsResponse, error)

		// GetWAFProtection retrieves the current WAF protection setting for a configuration and policy.
		// Deprecated: this method will be removed in a future release. Use the GetPolicyProtections method of the PolicyProtections interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-protections
		GetWAFProtection(ctx context.Context, params GetWAFProtectionRequest) (*GetWAFProtectionResponse, error)

		// UpdateWAFProtection updates the WAF protection setting for a configuration and policy.
		// Deprecated: this method will be removed in a future release. Use the CreateSecurityPolicyWithDefaultProtections method of the SecurityPolicy interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-policy-protections
		UpdateWAFProtection(ctx context.Context, params UpdateWAFProtectionRequest) (*UpdateWAFProtectionResponse, error)
	}

	// GetWAFProtectionRequest is used to retrieve the WAF protection setting.
	GetWAFProtectionRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}

	// GetWAFProtectionResponse is returned from a call to GetWAFProtection.
	GetWAFProtectionResponse ProtectionsResponse

	// GetWAFProtectionsRequest is used to retrieve the WAF protection setting.
	// Deprecated: this struct will be removed in a future release.
	GetWAFProtectionsRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}

	// GetWAFProtectionsResponse is returned from a call to GetWAFProtections.
	// Deprecated: this struct will be removed in a future release.
	GetWAFProtectionsResponse ProtectionsResponse

	// UpdateWAFProtectionRequest is used to modify the WAF protection setting.
	UpdateWAFProtectionRequest struct {
		ConfigID                      int    `json:"-"`
		Version                       int    `json:"-"`
		PolicyID                      string `json:"-"`
		ApplyApplicationLayerControls bool   `json:"applyApplicationLayerControls"`
	}

	// UpdateWAFProtectionResponse is returned from a call to UpdateWAFProtection.
	UpdateWAFProtectionResponse ProtectionsResponse
)

// Validate validates a GetWAFProtectionRequest.
func (v GetWAFProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetWAFProtectionsRequest.
func (v GetWAFProtectionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateWAFProtectionRequest.
func (v UpdateWAFProtectionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetWAFProtection(ctx context.Context, params GetWAFProtectionRequest) (*GetWAFProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetWAFProtection")

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
		return nil, fmt.Errorf("failed to create GetWAFProtection request: %w", err)
	}

	var result GetWAFProtectionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get WAF protection request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetWAFProtections(ctx context.Context, params GetWAFProtectionsRequest) (*GetWAFProtectionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetWAFProtections")

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
		return nil, fmt.Errorf("failed to create GetWAFProtections request: %w", err)
	}

	var result GetWAFProtectionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get WAF protections request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateWAFProtection(ctx context.Context, params UpdateWAFProtectionRequest) (*UpdateWAFProtectionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateWAFProtection")

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
		return nil, fmt.Errorf("failed to create UpdateWAFProtection request: %w", err)
	}

	var result UpdateWAFProtectionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update WAF protection request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
