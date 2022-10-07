package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SecurityPolicyClone interface supports cloning an existing security policy and retrieving
	// existing security policies.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#securitypolicyclone
	SecurityPolicyClone interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsecuritypolicies
		GetSecurityPolicyClones(ctx context.Context, params GetSecurityPolicyClonesRequest) (*GetSecurityPolicyClonesResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsecuritypolicies
		GetSecurityPolicyClone(ctx context.Context, params GetSecurityPolicyCloneRequest) (*GetSecurityPolicyCloneResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postsecuritypolicies
		CreateSecurityPolicyClone(ctx context.Context, params CreateSecurityPolicyCloneRequest) (*CreateSecurityPolicyCloneResponse, error)
	}

	// GetSecurityPolicyClonesRequest is used to retrieve the available security policies.
	GetSecurityPolicyClonesRequest struct {
		ConfigID int `json:"configId"`
		Version  int `json:"version"`
	}

	// GetSecurityPolicyClonesResponse is returned from a call to GetSecurityPolicyClones.
	GetSecurityPolicyClonesResponse struct {
		ConfigID int `json:"configId"`
		Version  int `json:"version"`
		Policies []struct {
			PolicyID                string `json:"policyId"`
			PolicyName              string `json:"policyName"`
			HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey"`
			PolicySecurityControls  struct {
				ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
				ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
				ApplyRateControls             bool `json:"applyRateControls"`
				ApplyReputationControls       bool `json:"applyReputationControls"`
				ApplyBotmanControls           bool `json:"applyBotmanControls"`
				ApplyAPIConstraints           bool `json:"applyApiConstraints"`
				ApplySlowPostControls         bool `json:"applySlowPostControls"`
			} `json:"policySecurityControls"`
		} `json:"policies"`
	}

	// GetSecurityPolicyCloneRequest is used to retrieve a security policy.
	GetSecurityPolicyCloneRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
	}

	// GetSecurityPolicyCloneResponse is returned from a call to GetSecurityPolicyClone.
	GetSecurityPolicyCloneResponse struct {
		ConfigID               int    `json:"configId,omitempty"`
		PolicyID               string `json:"policyId,omitempty"`
		PolicyName             string `json:"policyName,omitempty"`
		PolicySecurityControls struct {
			ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
			ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
			ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
			ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
			ApplyRateControls             bool `json:"applyRateControls,omitempty"`
			ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
			ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
		} `json:"policySecurityControls,omitempty"`
		Version int `json:"version,omitempty"`
	}

	// CreateSecurityPolicyCloneRequest is used to clone a security policy.
	CreateSecurityPolicyCloneRequest struct {
		ConfigID                 int    `json:"configId"`
		Version                  int    `json:"version"`
		CreateFromSecurityPolicy string `json:"createFromSecurityPolicy"`
		PolicyName               string `json:"policyName"`
		PolicyPrefix             string `json:"policyPrefix"`
	}

	// CreateSecurityPolicyCloneResponse is returned from a call to CreateSecurityPolicyClone.
	CreateSecurityPolicyCloneResponse struct {
		HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey"`
		PolicyID                string `json:"policyId"`
		PolicyName              string `json:"policyName"`
		PolicySecurityControls  struct {
			ApplyAPIConstraints           bool `json:"applyApiConstraints"`
			ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
			ApplyBotmanControls           bool `json:"applyBotmanControls"`
			ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
			ApplyRateControls             bool `json:"applyRateControls"`
			ApplyReputationControls       bool `json:"applyReputationControls"`
			ApplySlowPostControls         bool `json:"applySlowPostControls"`
		}
	}

	// SecurityPolicyCloneResponse is currently unused.
	SecurityPolicyCloneResponse struct {
		ConfigID int        `json:"configId"`
		Policies []Policies `json:"policies"`
		Version  int        `json:"version"`
	}

	// Policies is used as part of a description of available security policies.
	Policies struct {
		HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey"`
		PolicyID                string `json:"policyId"`
		PolicyName              string `json:"policyName"`
		PolicySecurityControls  struct {
			ApplyAPIConstraints           bool `json:"applyApiConstraints"`
			ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
			ApplyBotmanControls           bool `json:"applyBotmanControls"`
			ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
			ApplyRateControls             bool `json:"applyRateControls"`
			ApplyReputationControls       bool `json:"applyReputationControls"`
			ApplySlowPostControls         bool `json:"applySlowPostControls"`
		}
	}

	// CreateSecurityPolicyClonePost is currently unused.
	CreateSecurityPolicyClonePost struct {
		CreateFromSecurityPolicy string `json:"createFromSecurityPolicy"`
		PolicyName               string `json:"policyName"`
		PolicyPrefix             string `json:"policyPrefix"`
	}

	// CreateSecurityPolicyClonePostResponse is currently unused.
	CreateSecurityPolicyClonePostResponse struct {
		ConfigID               int    `json:"configId"`
		PolicyID               string `json:"policyId"`
		PolicyName             string `json:"policyName"`
		PolicySecurityControls struct {
			ApplyAPIConstraints           bool `json:"applyApiConstraints"`
			ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
			ApplyBotmanControls           bool `json:"applyBotmanControls"`
			ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
			ApplyRateControls             bool `json:"applyRateControls"`
			ApplyReputationControls       bool `json:"applyReputationControls"`
			ApplySlowPostControls         bool `json:"applySlowPostControls"`
		} `json:"policySecurityControls"`
		Version int `json:"version"`
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
