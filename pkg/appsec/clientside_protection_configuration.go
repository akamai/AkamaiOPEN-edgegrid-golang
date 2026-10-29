package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ClientsideProtectionConfiguration interface supports retrieving and modifying clientside protection configuration.
	ClientsideProtectionConfiguration interface {
		// GetClientsideProtectionConfiguration retrieves the clientside protection configuration for a security policy.
		GetClientsideProtectionConfiguration(ctx context.Context, params GetClientsideProtectionConfigurationRequest) (*GetClientsideProtectionConfigurationResponse, error)

		// UpdateClientsideProtectionConfiguration creates or updates the clientside protection configuration for a security policy (upsert).
		UpdateClientsideProtectionConfiguration(ctx context.Context, params UpdateClientsideProtectionConfigurationRequest) (*UpdateClientsideProtectionConfigurationResponse, error)
	}

	// GetClientsideProtectionConfigurationRequest is used to retrieve the clientside protection configuration
	GetClientsideProtectionConfigurationRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
	}

	// ClientsideProtectionEdgeInjection contains edge injection settings
	ClientsideProtectionEdgeInjection struct {
		ClientsideProtectionConfigID int  `json:"clientSideProtectionConfigId,omitempty"`
		LoadScriptAsync              bool `json:"loadScriptAsync"`
	}

	// ClientsideProtectionEdgeInjectionRequest contains edge injection settings for the request (excludes read-only fields)
	ClientsideProtectionEdgeInjectionRequest struct {
		LoadScriptAsync bool `json:"loadScriptAsync"`
	}

	// ClientsideProtectionEdgeTestParameters contains test parameters for enabling or disabling injection
	ClientsideProtectionEdgeTestParameters struct {
		DisableInjectionKey   string `json:"disableInjectionKey"`
		DisableInjectionValue string `json:"disableInjectionValue"`
		ForceInjectionKey     string `json:"forceInjectionKey"`
		ForceInjectionValue   string `json:"forceInjectionValue"`
	}

	// ClientsideProtectionInjectionCriteria contains criteria controlling where and how scripts are injected
	ClientsideProtectionInjectionCriteria struct {
		InjectionPathIncludes *string  `json:"injectionPathIncludes,omitempty"`
		InjectionPathExcludes *string  `json:"injectionPathExcludes,omitempty"`
		InjectionPercent      int      `json:"injectionPercent"`
		InjectionPolicy       *string  `json:"injectionPolicy,omitempty"`
		ScriptTagAttributes   []string `json:"scriptTagAttributes,omitempty"`
	}

	// ClientsideProtectionConfigurationResponse is the common response shape for get and update operations
	ClientsideProtectionConfigurationResponse struct {
		EdgeInjection      ClientsideProtectionEdgeInjection      `json:"edgeInjection"`
		EdgeTestParameters ClientsideProtectionEdgeTestParameters `json:"edgeTestParameters"`
		InjectionCriteria  ClientsideProtectionInjectionCriteria  `json:"injectionCriteria"`
	}

	// GetClientsideProtectionConfigurationResponse is an alias for ClientsideProtectionConfigurationResponse
	GetClientsideProtectionConfigurationResponse = ClientsideProtectionConfigurationResponse

	// UpdateClientsideProtectionConfigurationResponse is an alias for ClientsideProtectionConfigurationResponse
	UpdateClientsideProtectionConfigurationResponse = ClientsideProtectionConfigurationResponse

	// UpdateClientsideProtectionConfigurationRequest is used to create or update the clientside protection configuration
	UpdateClientsideProtectionConfigurationRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		Body             UpdateClientsideProtectionConfigurationRequestBody
	}

	// UpdateClientsideProtectionConfigurationRequestBody contains the request body for updating the clientside protection configuration
	UpdateClientsideProtectionConfigurationRequestBody struct {
		EdgeInjection      ClientsideProtectionEdgeInjectionRequest `json:"edgeInjection"`
		EdgeTestParameters ClientsideProtectionEdgeTestParameters   `json:"edgeTestParameters"`
		InjectionCriteria  ClientsideProtectionInjectionCriteria    `json:"injectionCriteria"`
	}
)

// Validate validates GetClientsideProtectionConfigurationRequest
func (v GetClientsideProtectionConfigurationRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required, validation.Min(int64(1))),
		"Version":          validation.Validate(v.Version, validation.Required, validation.Min(int64(1))),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateClientsideProtectionConfigurationRequest
func (v UpdateClientsideProtectionConfigurationRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required, validation.Min(int64(1))),
		"Version":          validation.Validate(v.Version, validation.Required, validation.Min(int64(1))),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"Body.InjectionCriteria.InjectionPercent": validation.Validate(v.Body.InjectionCriteria.InjectionPercent, validation.Min(0), validation.Max(100)),
	}.Filter()
}

func (p *appsec) GetClientsideProtectionConfiguration(ctx context.Context, params GetClientsideProtectionConfigurationRequest) (*GetClientsideProtectionConfigurationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetClientsideProtectionConfiguration")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewGet(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/client-side-protection",
		params.ConfigID, params.Version, params.SecurityPolicyID).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create GetClientsideProtectionConfiguration request: %w", err)
	}

	var result GetClientsideProtectionConfigurationResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get clientside protection configuration request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.cpcError(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateClientsideProtectionConfiguration(ctx context.Context, params UpdateClientsideProtectionConfigurationRequest) (*UpdateClientsideProtectionConfigurationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateClientsideProtectionConfiguration")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewPut(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/client-side-protection",
		params.ConfigID, params.Version, params.SecurityPolicyID).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateClientsideProtectionConfiguration request: %w", err)
	}

	var result UpdateClientsideProtectionConfigurationResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("update clientside protection configuration request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.cpcError(resp)
	}

	return &result, nil
}
