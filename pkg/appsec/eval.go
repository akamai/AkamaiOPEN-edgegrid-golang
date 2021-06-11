package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The Eval interface supports retrieving and updating the way evaluation rules would respond if
	// they were applied to live traffic.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#evalmode
	Eval interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getmode
		GetEvals(ctx context.Context, params GetEvalsRequest) (*GetEvalsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getmode
		GetEval(ctx context.Context, params GetEvalRequest) (*GetEvalResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postevaluationmode
		UpdateEval(ctx context.Context, params UpdateEvalRequest) (*UpdateEvalResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postevaluationmode
		RemoveEval(ctx context.Context, params RemoveEvalRequest) (*RemoveEvalResponse, error)
	}

	// GetEvalsRequest is used to retrieve the mode setting that conveys how rules will be kept up to date.
	GetEvalsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	// GetEvalsResponse is returned from a call to GetEvalsResponse.
	GetEvalsResponse struct {
		Current    string `json:"current,omitempty"`
		Mode       string `json:"mode,omitempty"`
		Eval       string `json:"eval,omitempty"`
		Evaluating string `json:"evaluating,omitempty"`
		Expires    string `json:"expires,omitempty"`
	}

	// GetEvalRequest is used to retrieve the mode setting that conveys how rules will be kept up to date.
	GetEvalRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	// GetEvalResponse is returned from a call to GetEvalResponse.
	GetEvalResponse struct {
		Current    string `json:"current,omitempty"`
		Mode       string `json:"mode,omitempty"`
		Eval       string `json:"eval,omitempty"`
		Evaluating string `json:"evaluating,omitempty"`
		Expires    string `json:"expires,omitempty"`
	}

	// RemoveEvalRequest is used to remove an evaluation mode setting.
	RemoveEvalRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"-"`
		Mode     string `json:"-"`
		Eval     string `json:"eval"`
	}

	// RemoveEvalResponse is returned from a call to RemoveEval.
	RemoveEvalResponse struct {
		Current string `json:"current"`
		Eval    string `json:"eval"`
		Mode    string `json:"mode"`
	}

	// UpdateEvalRequest is used to modify an evaluation mode setting.
	UpdateEvalRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"-"`
		Mode     string `json:"-"`
		Eval     string `json:"eval"`
	}

	// UpdateEvalResponse is returned from a call to UpdateEval.
	UpdateEvalResponse struct {
		Current string `json:"current"`
		Eval    string `json:"eval"`
		Mode    string `json:"mode"`
	}
)

// Validate validates a GetEvalRequest.
func (v GetEvalRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetEvalsRequest.
func (v GetEvalsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateEvalRequest.
func (v UpdateEvalRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveEvalRequest.
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
		return nil, fmt.Errorf("failed to create GetEval request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetEval request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create GetEvals request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetEvals request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create UpdateEval request: %w", err)
	}

	var rval UpdateEvalResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateEval request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

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
		return nil, fmt.Errorf("failed to create RemoveEval request: %w", err)
	}

	var rval RemoveEvalResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveEval request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
