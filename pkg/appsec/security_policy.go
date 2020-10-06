package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
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
	}

	GetSecurityPoliciesRequest struct {
		ConfigID int `json:"configId"`
		Version  int `json:"version"`
	}

	GetSecurityPoliciesResponse struct {
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
)

func (p *appsec) GetSecurityPolicies(ctx context.Context, params GetSecurityPoliciesRequest) (*GetSecurityPoliciesResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetSecurityPolicys")

	var rval GetSecurityPoliciesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies?notMatched=false&detail=true",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getsecuritypolicys request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getsecuritypolicys request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &rval, nil

}
