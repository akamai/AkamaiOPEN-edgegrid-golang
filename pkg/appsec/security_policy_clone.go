package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SecurityPolicyClone represents a collection of SecurityPolicyClone
//
// See: SecurityPolicyClone.GetSecurityPolicyClone()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// SecurityPolicyClone  contains operations available on SecurityPolicyClone  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsecuritypolicyclone
	SecurityPolicyClone interface {
		GetSecurityPolicyClones(ctx context.Context, params GetSecurityPolicyClonesRequest) (*GetSecurityPolicyClonesResponse, error)
		GetSecurityPolicyClone(ctx context.Context, params GetSecurityPolicyCloneRequest) (*GetSecurityPolicyCloneResponse, error)
		CreateSecurityPolicyClone(ctx context.Context, params CreateSecurityPolicyCloneRequest) (*CreateSecurityPolicyCloneResponse, error)
	}

	GetSecurityPolicyClonesRequest struct {
		ConfigID int `json:"configId"`
		Version  int `json:"version"`
	}

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

	GetSecurityPolicyCloneRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
	}

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

	CreateSecurityPolicyCloneRequest struct {
		ConfigID                 int    `json:"configId"`
		Version                  int    `json:"version"`
		CreateFromSecurityPolicy string `json:"createFromSecurityPolicy"`
		PolicyName               string `json:"policyName"`
		PolicyPrefix             string `json:"policyPrefix"`
	}

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

	SecurityPolicyCloneResponse struct {
		ConfigID int        `json:"configId"`
		Policies []Policies `json:"policies"`
		Version  int        `json:"version"`
	}

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

	CreateSecurityPolicyClonePost struct {
		CreateFromSecurityPolicy string `json:"createFromSecurityPolicy"`
		PolicyName               string `json:"policyName"`
		PolicyPrefix             string `json:"policyPrefix"`
	}

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

// Validate validates GetSecurityPolicyCloneRequest
func (v GetSecurityPolicyCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates GetSecurityPolicyClonesRequest
func (v GetSecurityPolicyClonesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates CreateSecurityPolicyCloneRequest
func (v CreateSecurityPolicyCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetSecurityPolicyClone(ctx context.Context, params GetSecurityPolicyCloneRequest) (*GetSecurityPolicyCloneResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicyClone")

	var rvals GetSecurityPolicyCloneResponse

	uri := fmt.Sprintf(
		//	"/appsec/v1/configs/%d/versions/%d/security-policies?notMatched=false&detail=true",
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getsecuritypolicyclone request: %w", err)
	}

	resp, err := p.Exec(req, &rvals)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rvals, nil

}

func (p *appsec) GetSecurityPolicyClones(ctx context.Context, params GetSecurityPolicyClonesRequest) (*GetSecurityPolicyClonesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicyClone")

	var rval GetSecurityPolicyClonesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies?detail=true&notMatched=false",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getsecuritypolicyclone request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("gGetSecurityPolicyClone request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

/// Create will create a new securitypolicyclone.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postsecuritypolicyclone
func (p *appsec) CreateSecurityPolicyClone(ctx context.Context, params CreateSecurityPolicyCloneRequest) (*CreateSecurityPolicyCloneResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateSecurityPolicyClone")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create securitypolicyclone request: %w", err)
	}

	var rval CreateSecurityPolicyCloneResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create securitypolicyclonerequest failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
