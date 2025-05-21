package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SecurityPolicy interface supports creating, retrieving, modifying and removing security policies.
	SecurityPolicy interface {
		// GetSecurityPolicies returns a list of security policies available for the specified security configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policies
		GetSecurityPolicies(ctx context.Context, params GetSecurityPoliciesRequest) (*GetSecurityPoliciesResponse, error)

		// GetSecurityPolicy returns the specified security policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy
		GetSecurityPolicy(ctx context.Context, params GetSecurityPolicyRequest) (*GetSecurityPolicyResponse, error)

		// CreateSecurityPolicy creates a new copy of an existing security policy or creates a new security policy from scratch
		// when you don't specify a policy to clone in the request.
		// Deprecated: this method will be removed in a future release. Use the CreateSecurityPolicyWithDefaultProtections method instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-policy
		CreateSecurityPolicy(ctx context.Context, params CreateSecurityPolicyRequest) (*CreateSecurityPolicyResponse, error)

		// CreateSecurityPolicyWithDefaultProtections creates a new security policy with a specified set of security protections.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-policy, https://techdocs.akamai.com/application-security/reference/put-policy-protections
		CreateSecurityPolicyWithDefaultProtections(ctx context.Context, params CreateSecurityPolicyWithDefaultProtectionsRequest) (*CreateSecurityPolicyResponse, error)

		// UpdateSecurityPolicy updates the name of a specific security policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-policy
		UpdateSecurityPolicy(ctx context.Context, params UpdateSecurityPolicyRequest) (*UpdateSecurityPolicyResponse, error)

		// RemoveSecurityPolicy deletes the specified security policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/delete-policy
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

	// CreateSecurityPolicyRequest is used to create a security policy.
	CreateSecurityPolicyRequest struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		PolicyName      string `json:"policyName"`
		PolicyPrefix    string `json:"policyPrefix"`
		DefaultSettings bool   `json:"defaultSettings"`
	}

	// CreateSecurityPolicyWithDefaultProtectionsRequest is used to create a security policy with a specified set of protections.
	CreateSecurityPolicyWithDefaultProtectionsRequest struct {
		ConfigVersion
		PolicyName   string `json:"policyName"`
		PolicyPrefix string `json:"policyPrefix"`
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
		ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
		ApplyMalwareControls          bool `json:"applyMalwareControls,omitempty"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}
)

// Validate validates a GetSecurityPolicyRequest.
func (v GetSecurityPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates a GetSecurityPolicysRequest.
func (v GetSecurityPoliciesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates a CreateSecurityPolicyRequest.
func (v CreateSecurityPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":     validation.Validate(v.ConfigID, validation.Required),
		"Version":      validation.Validate(v.Version, validation.Required),
		"PolicyName":   validation.Validate(v.PolicyName, validation.Required),
		"PolicyPrefix": validation.Validate(v.PolicyPrefix, validation.Required),
	})
}

// Validate validates a CreateSecurityPolicyWithDefaultProtectionsRequest.
func (v CreateSecurityPolicyWithDefaultProtectionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":     validation.Validate(v.ConfigID, validation.Required),
		"Version":      validation.Validate(v.Version, validation.Required),
		"PolicyName":   validation.Validate(v.PolicyName, validation.Required),
		"PolicyPrefix": validation.Validate(v.PolicyPrefix, validation.Required),
	})
}

// Validate validates an UpdateSecurityPolicyRequest.
func (v UpdateSecurityPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	})
}

// Validate validates a RemoveSecurityPolicyRequest.
func (v RemoveSecurityPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	})
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
	defer session.CloseResponseBody(resp)

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
	defer session.CloseResponseBody(resp)

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
	defer session.CloseResponseBody(resp)

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateSecurityPolicyWithDefaultProtections(ctx context.Context, params CreateSecurityPolicyWithDefaultProtectionsRequest) (*CreateSecurityPolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateSecurityPolicyWithDefaultProtections")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/protections",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateSecurityPolicyWithDefaultProtections request: %w", err)
	}

	var result CreateSecurityPolicyResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("create security policy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
