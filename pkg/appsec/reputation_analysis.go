package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ReputationAnalysis interface supports retrieving and modifying the reputation analysis
	// settings for a configuration and policy.
	ReputationAnalysis interface {
		// GetReputationAnalysis returns the current reputation analysis settings.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-reputation-analysis
		GetReputationAnalysis(ctx context.Context, params GetReputationAnalysisRequest) (*GetReputationAnalysisResponse, error)

		// UpdateReputationAnalysis updates the reputation analysis settings.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-reputation-analysis
		UpdateReputationAnalysis(ctx context.Context, params UpdateReputationAnalysisRequest) (*UpdateReputationAnalysisResponse, error)

		// RemoveReputationAnalysis removes the reputation analysis settings.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-reputation-analysis
		RemoveReputationAnalysis(ctx context.Context, params RemoveReputationAnalysisRequest) (*RemoveReputationAnalysisResponse, error)
	}

	// GetReputationAnalysisRequest is used to retrieve the reputation analysis settings for a security policy.
	GetReputationAnalysisRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
	}

	// GetReputationAnalysisResponse is returned from a call to GetReputationAnalysis.
	GetReputationAnalysisResponse struct {
		ConfigID                           int    `json:"-"`
		Version                            int    `json:"-"`
		PolicyID                           string `json:"-"`
		ForwardToHTTPHeader                bool   `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool   `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}

	// UpdateReputationAnalysisRequest is used to modify the reputation analysis settings for a security poliyc.
	UpdateReputationAnalysisRequest struct {
		ConfigID                           int    `json:"-"`
		Version                            int    `json:"-"`
		PolicyID                           string `json:"-"`
		ForwardToHTTPHeader                bool   `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool   `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}

	// UpdateReputationAnalysisResponse is returned from a call to UpdateReputationAnalysis.
	UpdateReputationAnalysisResponse struct {
		ForwardToHTTPHeader                bool `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}

	// RemoveReputationAnalysisRequest is used to remove the reputation analysis settings for a security policy.
	RemoveReputationAnalysisRequest struct {
		ConfigID                           int    `json:"-"`
		Version                            int    `json:"-"`
		PolicyID                           string `json:"-"`
		ForwardToHTTPHeader                bool   `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool   `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}

	// RemoveReputationAnalysisResponse is returned from a call to RemoveReputationAnalysis.
	RemoveReputationAnalysisResponse struct {
		ForwardToHTTPHeader                bool `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}
)

// Validate validates a GetReputationAnalysisRequest.
func (v GetReputationAnalysisRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateReputationAnalysisRequest.
func (v UpdateReputationAnalysisRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveReputationAnalysisRequest.
func (v RemoveReputationAnalysisRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetReputationAnalysis(ctx context.Context, params GetReputationAnalysisRequest) (*GetReputationAnalysisResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetReputationAnalysis")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-analysis",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetReputationAnalysis request: %w", err)
	}

	var result GetReputationAnalysisResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get reputation analysis request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateReputationAnalysis(ctx context.Context, params UpdateReputationAnalysisRequest) (*UpdateReputationAnalysisResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateReputationAnalysis")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-analysis",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateReputationAnalysis request: %w", err)
	}

	var result UpdateReputationAnalysisResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update reputation analysis request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveReputationAnalysis(ctx context.Context, params RemoveReputationAnalysisRequest) (*RemoveReputationAnalysisResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveReputationAnalysis")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-analysis",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveReputationAnalysis request: %w", err)
	}

	var result RemoveReputationAnalysisResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove reputation analysis request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
