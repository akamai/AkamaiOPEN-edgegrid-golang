package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Eval represents a collection of Eval
//
// See: Eval.GetEval()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// Eval  contains operations available on Eval  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#geteval
	Eval interface {
		GetEvals(ctx context.Context, params GetEvalsRequest) (*GetEvalsResponse, error)
		GetEval(ctx context.Context, params GetEvalRequest) (*GetEvalResponse, error)
		UpdateEval(ctx context.Context, params UpdateEvalRequest) (*UpdateEvalResponse, error)
		RemoveEval(ctx context.Context, params RemoveEvalRequest) (*RemoveEvalResponse, error)
	}

	GetEvalsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	GetEvalRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	GetEvalsResponse struct {
		Current    string `json:"current,omitempty"`
		Mode       string `json:"mode,omitempty"`
		Eval       string `json:"eval,omitempty"`
		Evaluating string `json:"evaluating,omitempty"`
		Expires    string `json:"expires,omitempty"`
	}

	GetEvalResponse struct {
		Current    string `json:"current,omitempty"`
		Mode       string `json:"mode,omitempty"`
		Eval       string `json:"eval,omitempty"`
		Evaluating string `json:"evaluating,omitempty"`
		Expires    string `json:"expires,omitempty"`
	}

	RemoveEvalRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"-"`
		Mode     string `json:"-"`
		Eval     string `json:"eval"`
	}

	RemoveEvalResponse struct {
		Current string `json:"current"`
		Eval    string `json:"eval"`
		Mode    string `json:"mode"`
	}

	UpdateEvalRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"-"`
		Mode     string `json:"-"`
		Eval     string `json:"eval"`
	}

	UpdateEvalResponse struct {
		Current string `json:"current"`
		Eval    string `json:"eval"`
		Mode    string `json:"mode"`
	}
)

// Validate validates GetEvalRequest
func (v GetEvalRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetEvalsRequest
func (v GetEvalsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateEvalRequest
func (v UpdateEvalRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateEvalRequest
func (v RemoveEvalRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetEval(ctx context.Context, params GetEvalRequest) (*GetEvalResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEval")

	var rval GetEvalResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/mode",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create geteval request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("geteval  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetEvals(ctx context.Context, params GetEvalsRequest) (*GetEvalsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvals")

	var rval GetEvalsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/mode",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevals request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevals request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a Eval.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#puteval

func (p *appsec) UpdateEval(ctx context.Context, params UpdateEvalRequest) (*UpdateEvalResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEval")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create Evalrequest: %w", err)
	}

	var rval UpdateEvalResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create Eval request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Remove will update a Eval.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#puteval

func (p *appsec) RemoveEval(ctx context.Context, params RemoveEvalRequest) (*RemoveEvalResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEval")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create Evalrequest: %w", err)
	}

	var rval RemoveEvalResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create Eval request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
