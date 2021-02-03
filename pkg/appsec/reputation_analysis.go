package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ReputationAnalysis represents a collection of ReputationAnalysis
//
// See: ReputationAnalysis.GetReputationAnalysis()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ReputationAnalysis  contains operations available on ReputationAnalysis  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getreputationanalysis
	ReputationAnalysis interface {
		GetReputationAnalysis(ctx context.Context, params GetReputationAnalysisRequest) (*GetReputationAnalysisResponse, error)
		UpdateReputationAnalysis(ctx context.Context, params UpdateReputationAnalysisRequest) (*UpdateReputationAnalysisResponse, error)
		RemoveReputationAnalysis(ctx context.Context, params RemoveReputationAnalysisRequest) (*RemoveReputationAnalysisResponse, error)
	}

	GetReputationAnalysisResponse struct {
		ConfigID                           int    `json:"-"`
		Version                            int    `json:"-"`
		PolicyID                           string `json:"-"`
		ForwardToHTTPHeader                bool   `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool   `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}

	UpdateReputationAnalysisResponse struct {
		ForwardToHTTPHeader                bool `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}

	RemoveReputationAnalysisResponse struct {
		ForwardToHTTPHeader                bool `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}

	GetReputationAnalysisRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
	}

	UpdateReputationAnalysisRequest struct {
		ConfigID                           int    `json:"-"`
		Version                            int    `json:"-"`
		PolicyID                           string `json:"-"`
		ForwardToHTTPHeader                bool   `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool   `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}
	RemoveReputationAnalysisRequest struct {
		ConfigID                           int    `json:"-"`
		Version                            int    `json:"-"`
		PolicyID                           string `json:"-"`
		ForwardToHTTPHeader                bool   `json:"forwardToHTTPHeader"`
		ForwardSharedIPToHTTPHeaderAndSIEM bool   `json:"forwardSharedIPToHTTPHeaderAndSIEM"`
	}
)

// Validate validates GetReputationAnalysisRequest
func (v GetReputationAnalysisRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateReputationAnalysisRequest
func (v UpdateReputationAnalysisRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates RemoveReputationAnalysisRequest
func (v RemoveReputationAnalysisRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetReputationAnalysis(ctx context.Context, params GetReputationAnalysisRequest) (*GetReputationAnalysisResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetReputationAnalysis")

	var rval GetReputationAnalysisResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-analysis",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getreputationanalysis request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a ReputationAnalysis.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putreputationanalysis

func (p *appsec) UpdateReputationAnalysis(ctx context.Context, params UpdateReputationAnalysisRequest) (*UpdateReputationAnalysisResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateReputationAnalysis")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-analysis",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create ReputationAnalysisrequest: %w", err)
	}

	var rval UpdateReputationAnalysisResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create ReputationAnalysis request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Delete will delete a ReputationAnalysis
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deletereputationanalysis

func (p *appsec) RemoveReputationAnalysis(ctx context.Context, params RemoveReputationAnalysisRequest) (*RemoveReputationAnalysisResponse, error) {

	var rval RemoveReputationAnalysisResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveReputationAnalysis")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-analysis",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create ReputationAnalysisrequest: %w", err)
	}

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("delreputationanalysis request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
