package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SecurityPolicyClone interface supports cloning an existing security policy and retrieving
	// existing security policies.
	SecurityPolicyClone interface {
		// GetSecurityPolicyClones returns a list of security policies available for the specified security configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policies
		GetSecurityPolicyClones(ctx context.Context, params GetSecurityPolicyClonesRequest) (*GetSecurityPolicyClonesResponse, error)

		// GetSecurityPolicyClone returns the specified security policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy
		GetSecurityPolicyClone(ctx context.Context, params GetSecurityPolicyCloneRequest) (*GetSecurityPolicyCloneResponse, error)

		// CreateSecurityPolicyClone Creates a new copy of an existing security policy or creates a new security policy from scratch
		// when you don't specify a policy to clone in the request.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-policy
		CreateSecurityPolicyClone(ctx context.Context, params CreateSecurityPolicyCloneRequest) (*CreateSecurityPolicyCloneResponse, error)
	}

	// GetSecurityPolicyClonesRequest is used to retrieve the available security policies.
	GetSecurityPolicyClonesRequest struct {
		ConfigID int64 `json:"configId"`
		Version  int64 `json:"version"`
	}

	// GetSecurityPolicyClonesResponse is returned from a call to GetSecurityPolicyClones.
	GetSecurityPolicyClonesResponse struct {
		ConfigID int64                                         `json:"configId"`
		Version  int64                                         `json:"version"`
		Policies []GetSecurityPolicyClonesResponsePoliciesItem `json:"policies"`
	}

	// GetSecurityPolicyClonesResponsePoliciesItem represents an individual policy item in GetSecurityPolicyClonesResponse.
	GetSecurityPolicyClonesResponsePoliciesItem struct {
		PolicyID                string                 `json:"policyId"`
		PolicyName              string                 `json:"policyName"`
		HasRatePolicyWithAPIKey bool                   `json:"hasRatePolicyWithApiKey"`
		PolicySecurityControls  PolicySecurityControls `json:"policySecurityControls"`
	}

	// GetSecurityPolicyCloneRequest is used to retrieve a security policy.
	GetSecurityPolicyCloneRequest struct {
		ConfigID int64  `json:"configId"`
		Version  int64  `json:"version"`
		PolicyID string `json:"policyId"`
	}

	// GetSecurityPolicyCloneResponse is returned from a call to GetSecurityPolicyClone.
	GetSecurityPolicyCloneResponse struct {
		ConfigID               int64                  `json:"configId,omitempty"`
		PolicyID               string                 `json:"policyId,omitempty"`
		PolicyName             string                 `json:"policyName,omitempty"`
		PolicySecurityControls PolicySecurityControls `json:"policySecurityControls,omitempty"`
		Version                int64                  `json:"version,omitempty"`
	}

	// CreateSecurityPolicyCloneRequest is used to clone a security policy.
	CreateSecurityPolicyCloneRequest struct {
		ConfigID                 int64  `json:"configId"`
		Version                  int64  `json:"version"`
		CreateFromSecurityPolicy string `json:"createFromSecurityPolicy"`
		PolicyName               string `json:"policyName"`
		PolicyPrefix             string `json:"policyPrefix"`
	}

	// CreateSecurityPolicyCloneResponse is returned from a call to CreateSecurityPolicyClone.
	CreateSecurityPolicyCloneResponse struct {
		HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey"`
		PolicyID                string `json:"policyId"`
		PolicyName              string `json:"policyName"`
		PolicySecurityControls  PolicySecurityControls
	}

	// SecurityPolicyCloneResponse is currently unused.
	SecurityPolicyCloneResponse struct {
		ConfigID int64      `json:"configId"`
		Policies []Policies `json:"policies"`
		Version  int64      `json:"version"`
	}

	// Policies is used as part of a description of available security policies.
	Policies CreateSecurityPolicyCloneResponse

	// CreateSecurityPolicyClonePost is currently unused.
	CreateSecurityPolicyClonePost struct {
		CreateFromSecurityPolicy string `json:"createFromSecurityPolicy"`
		PolicyName               string `json:"policyName"`
		PolicyPrefix             string `json:"policyPrefix"`
	}

	// CreateSecurityPolicyClonePostResponse is currently unused.
	CreateSecurityPolicyClonePostResponse struct {
		ConfigID               int64                  `json:"configId"`
		PolicyID               string                 `json:"policyId"`
		PolicyName             string                 `json:"policyName"`
		PolicySecurityControls PolicySecurityControls `json:"policySecurityControls"`
		Version                int64                  `json:"version"`
	}

	// PolicySecurityControls represents the security controls for a policy clone.
	PolicySecurityControls struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates a GetSecurityPolicyCloneRequest.
func (v GetSecurityPolicyCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a GetSecurityPolicyClonesRequest.
func (v GetSecurityPolicyClonesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateSecurityPolicyCloneRequest.
func (v CreateSecurityPolicyCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetSecurityPolicyClone(ctx context.Context, params GetSecurityPolicyCloneRequest) (*GetSecurityPolicyCloneResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicyClone")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSecurityPolicyClone request: %w", err)
	}

	var results GetSecurityPolicyCloneResponse
	resp, err := p.Exec(req, &results)
	if err != nil {
		return nil, fmt.Errorf("get security policy clone request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &results, nil
}

func (p *appsec) GetSecurityPolicyClones(ctx context.Context, params GetSecurityPolicyClonesRequest) (*GetSecurityPolicyClonesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicyClone")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies?detail=true&notMatched=false",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSecurityPolicyClones request: %w", err)
	}

	var result GetSecurityPolicyClonesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get security policy clones request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateSecurityPolicyClone(ctx context.Context, params CreateSecurityPolicyCloneRequest) (*CreateSecurityPolicyCloneResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateSecurityPolicyClone")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateSecurityPolicyClone request: %w", err)
	}

	var result CreateSecurityPolicyCloneResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("create security policy clone request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
