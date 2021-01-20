package appsec

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SecurityPolicy represents a collection of SecurityPolicy
//
// See: SecurityPolicy.GetSecurityPolicy()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// SecurityPolicy  contains operations available on SecurityPolicy  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsecuritypolicy
	SecurityPolicy interface {
		GetSecurityPolicies(ctx context.Context, params GetSecurityPoliciesRequest) (*GetSecurityPoliciesResponse, error)
		GetSecurityPolicy(ctx context.Context, params GetSecurityPolicyRequest) (*GetSecurityPolicyResponse, error)
		CreateSecurityPolicy(ctx context.Context, params CreateSecurityPolicyRequest) (*CreateSecurityPolicyResponse, error)
		UpdateSecurityPolicy(ctx context.Context, params UpdateSecurityPolicyRequest) (*UpdateSecurityPolicyResponse, error)
		RemoveSecurityPolicy(ctx context.Context, params RemoveSecurityPolicyRequest) (*RemoveSecurityPolicyResponse, error)
	}

	GetSecurityPoliciesRequest struct {
		ConfigID   int    `json:"configId"`
		Version    int    `json:"version"`
		PolicyName string `json:"-"`
	}

	GetSecurityPoliciesResponse struct {
		ConfigID int `json:"configId,omitempty"`
		Version  int `json:"version,omitempty"`
		Policies []struct {
			PolicyID                string `json:"policyId,omitempty"`
			PolicyName              string `json:"policyName,omitempty"`
			HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey,omitempty"`
			PolicySecurityControls  struct {
				ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
				ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls,omitempty"`
				ApplyRateControls             bool `json:"applyRateControls,omitempty"`
				ApplyReputationControls       bool `json:"applyReputationControls,omitempty"`
				ApplyBotmanControls           bool `json:"applyBotmanControls,omitempty"`
				ApplyAPIConstraints           bool `json:"applyApiConstraints,omitempty"`
				ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
			} `json:"policySecurityControls,omitempty"`
		} `json:"policies,omitempty"`
	}

	GetSecurityPolicyRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
	}

	GetSecurityPolicyResponse struct {
		ConfigID               int    `json:"configId,omitempty"`
		PolicyID               string `json:"policyId,omitempty"`
		PolicyName             string `json:"policyName,omitempty"`
		DefaultSettings        bool   `json:"defaultSettings,omitempty"`
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

	CreateSecurityPolicyRequest struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		PolicyName      string `json:"policyName"`
		PolicyPrefix    string `json:"policyPrefix"`
		DefaultSettings bool   `json:"defaultSettings"`
	}

	CreateSecurityPolicyResponse struct {
		ConfigID               int    `json:"configId"`
		PolicyID               string `json:"policyId"`
		PolicyName             string `json:"policyName"`
		DefaultSettings        bool   `json:"defaultSettings,omitempty"`
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

	UpdateSecurityPolicyRequest struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		PolicyName      string `json:"policyName"`
		DefaultSettings bool   `json:"defaultSettings,omitempty"`
		PolicyPrefix    string `json:"policyPrefix"`
	}

	UpdateSecurityPolicyResponse struct {
		ConfigID               int    `json:"configId"`
		PolicyID               string `json:"policyId"`
		PolicyName             string `json:"policyName"`
		DefaultSettings        bool   `json:"defaultSettings,omitempty"`
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

	RemoveSecurityPolicyRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
	}

	RemoveSecurityPolicyResponse struct {
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

// Validate validates GetSecurityPolicyRequest
func (v GetSecurityPolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates GetSecurityPolicysRequest
func (v GetSecurityPoliciesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates CreateSecurityPolicyRequest
func (v CreateSecurityPolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateSecurityPolicyRequest
func (v UpdateSecurityPolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates RemoveSecurityPolicyRequest
func (v RemoveSecurityPolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetSecurityPolicies(ctx context.Context, params GetSecurityPoliciesRequest) (*GetSecurityPoliciesResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicys")

	var rval GetSecurityPoliciesResponse
	var rvalfiltered GetSecurityPoliciesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getsecuritypolicies request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getsecuritypolicies request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.PolicyName != "" {
		for _, val := range rval.Policies {
			if val.PolicyName == params.PolicyName {
				rvalfiltered.Policies = append(rvalfiltered.Policies, val)
			}
		}

	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

func (p *appsec) GetSecurityPolicy(ctx context.Context, params GetSecurityPolicyRequest) (*GetSecurityPolicyResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicys")

	var rval GetSecurityPolicyResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getsecuritypolicies request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getsecuritypolicies request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a SecurityPolicy.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putsecuritypolicy

func (p *appsec) UpdateSecurityPolicy(ctx context.Context, params UpdateSecurityPolicyRequest) (*UpdateSecurityPolicyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateSecurityPolicy")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create SecurityPolicyrequest: %w", err)
	}

	var rval UpdateSecurityPolicyResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create SecurityPolicy request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Create will create a new securitypolicy.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postsecuritypolicy
func (p *appsec) CreateSecurityPolicy(ctx context.Context, params CreateSecurityPolicyRequest) (*CreateSecurityPolicyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateSecurityPolicy")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create securitypolicy request: %w", err)
	}

	var rval CreateSecurityPolicyResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create securitypolicyrequest failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Delete will delete a SecurityPolicy
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deletesecuritypolicy

func (p *appsec) RemoveSecurityPolicy(ctx context.Context, params RemoveSecurityPolicyRequest) (*RemoveSecurityPolicyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval RemoveSecurityPolicyResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveSecurityPolicy")

	uri, err := url.Parse(fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create delsecuritypolicy request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("delsecuritypolicy request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
