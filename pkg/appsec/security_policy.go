package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SecurityPolicy interface supports creating, retrieving, modifying and removing security policies.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#securitypolicy
	SecurityPolicy interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsecuritypolicies
		GetSecurityPolicies(ctx context.Context, params GetSecurityPoliciesRequest) (*GetSecurityPoliciesResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsecuritypolicy
		GetSecurityPolicy(ctx context.Context, params GetSecurityPolicyRequest) (*GetSecurityPolicyResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postsecuritypolicies
		CreateSecurityPolicy(ctx context.Context, params CreateSecurityPolicyRequest) (*CreateSecurityPolicyResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putsecuritypolicy
		UpdateSecurityPolicy(ctx context.Context, params UpdateSecurityPolicyRequest) (*UpdateSecurityPolicyResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deletesecuritypolicy
		RemoveSecurityPolicy(ctx context.Context, params RemoveSecurityPolicyRequest) (*RemoveSecurityPolicyResponse, error)
	}

	// GetSecurityPoliciesRequest is used to retrieve the security policies for a configuration.
	GetSecurityPoliciesRequest struct {
		ConfigID   int    `json:"configId"`
		Version    int    `json:"version"`
		PolicyName string `json:"-"`
	}

	// GetSecurityPoliciesResponse is returned from a call to GetSecurityPolicies.
	GetSecurityPoliciesResponse struct {
		ConfigID int `json:"configId,omitempty"`
		Version  int `json:"version,omitempty"`
		Policies []struct {
			PolicyID                string            `json:"policyId,omitempty"`
			PolicyName              string            `json:"policyName,omitempty"`
			HasRatePolicyWithAPIKey bool              `json:"hasRatePolicyWithApiKey,omitempty"`
			PolicySecurityControls  *SecurityControls `json:"policySecurityControls,omitempty"`
		} `json:"policies,omitempty"`
	}

	// GetSecurityPolicyRequest is used to retrieve information about a security policy.
	GetSecurityPolicyRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
	}

	// GetSecurityPolicyResponse is returned from a call to GetSecurityPolicy.
	GetSecurityPolicyResponse struct {
		ConfigID               int               `json:"configId,omitempty"`
		PolicyID               string            `json:"policyId,omitempty"`
		PolicyName             string            `json:"policyName,omitempty"`
		DefaultSettings        bool              `json:"defaultSettings,omitempty"`
		PolicySecurityControls *SecurityControls `json:"policySecurityControls,omitempty"`
		Version                int               `json:"version,omitempty"`
	}

	// CreateSecurityPolicyRequest is used to create a ecurity policy.
	CreateSecurityPolicyRequest struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		PolicyName      string `json:"policyName"`
		PolicyPrefix    string `json:"policyPrefix"`
		DefaultSettings bool   `json:"defaultSettings"`
	}

	// CreateSecurityPolicyResponse is returned from a call to CreateSecurityPolicy.
	CreateSecurityPolicyResponse struct {
		ConfigID               int               `json:"configId"`
		PolicyID               string            `json:"policyId"`
		PolicyName             string            `json:"policyName"`
		DefaultSettings        bool              `json:"defaultSettings,omitempty"`
		PolicySecurityControls *SecurityControls `json:"policySecurityControls,omitempty"`
		Version                int               `json:"version"`
	}

	// UpdateSecurityPolicyRequest is used to modify a security policy.
	UpdateSecurityPolicyRequest struct {
		ConfigID   int    `json:"-"`
		Version    int    `json:"-"`
		PolicyID   string `json:"-"`
		PolicyName string `json:"policyName"`
	}

	// UpdateSecurityPolicyResponse is returned from a call to UpdateSecurityPolicy.
	UpdateSecurityPolicyResponse struct {
		ConfigID               int               `json:"configId"`
		PolicyID               string            `json:"policyId"`
		PolicyName             string            `json:"policyName"`
		DefaultSettings        bool              `json:"defaultSettings,omitempty"`
		PolicySecurityControls *SecurityControls `json:"policySecurityControls,omitempty"`
		Version                int               `json:"version"`
	}

	// RemoveSecurityPolicyRequest is used to remove a security policy.
	RemoveSecurityPolicyRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
	}

	// RemoveSecurityPolicyResponse is returned from a call to RemoveSecurityPolicy.
	RemoveSecurityPolicyResponse struct {
		ConfigID               int               `json:"configId"`
		PolicyID               string            `json:"policyId"`
		PolicyName             string            `json:"policyName"`
		PolicySecurityControls *SecurityControls `json:"policySecurityControls,omitempty"`
		Version                int               `json:"version"`
	}

	// SecurityControls is returned as part of GetSecurityPoliciesResponse and similar responses.
	SecurityControls struct {
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}
)

// Validate validates a GetSecurityPolicyRequest.
func (v GetSecurityPolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a GetSecurityPolicysRequest.
func (v GetSecurityPoliciesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateSecurityPolicyRequest.
func (v CreateSecurityPolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateSecurityPolicyRequest.
func (v UpdateSecurityPolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveSecurityPolicyRequest.
func (v RemoveSecurityPolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetSecurityPolicies(ctx context.Context, params GetSecurityPoliciesRequest) (*GetSecurityPoliciesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicies")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSecurityPolicies request: %w", err)
	}

	var result GetSecurityPoliciesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get security policies request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.PolicyName != "" {
		var filteredResult GetSecurityPoliciesResponse
		for _, val := range result.Policies {
			if val.PolicyName == params.PolicyName {
				filteredResult.Policies = append(filteredResult.Policies, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) GetSecurityPolicy(ctx context.Context, params GetSecurityPolicyRequest) (*GetSecurityPolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicy")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSecurityPolicy request: %w", err)
	}

	var result GetSecurityPolicyResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get security policy request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateSecurityPolicy(ctx context.Context, params UpdateSecurityPolicyRequest) (*UpdateSecurityPolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateSecurityPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateSecurityPolicy request: %w", err)
	}

	var result UpdateSecurityPolicyResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update security policy request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateSecurityPolicy(ctx context.Context, params CreateSecurityPolicyRequest) (*CreateSecurityPolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateSecurityPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateSecurityPolicy request: %w", err)
	}

	var result CreateSecurityPolicyResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("create security policy request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveSecurityPolicy(ctx context.Context, params RemoveSecurityPolicyRequest) (*RemoveSecurityPolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveSecurityPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/security-policies/%s", params.ConfigID, params.Version, params.PolicyID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveSecurityPolicy request: %w", err)
	}

	var result RemoveSecurityPolicyResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("remove security policy request failed: %w", err)
	}
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
