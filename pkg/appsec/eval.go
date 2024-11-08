package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The Eval interface supports retrieving and updating the way evaluation rules would respond if
	// they were applied to live traffic.
	Eval interface {
		// GetEvals returns which modes your rules are currently set to.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-mode-1
		GetEvals(ctx context.Context, params GetEvalsRequest) (*GetEvalsResponse, error)

		// GetEval returns which mode your rules are currently set to.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-mode-1
		GetEval(ctx context.Context, params GetEvalRequest) (*GetEvalResponse, error)

		// UpdateEval updated the rule evaluation mode.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-policy-eval
		UpdateEval(ctx context.Context, params UpdateEvalRequest) (*UpdateEvalResponse, error)

		// RemoveEval removes the rule evaluation mode.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-policy-eval
		RemoveEval(ctx context.Context, params RemoveEvalRequest) (*RemoveEvalResponse, error)
	}

	// GetEvalsRequest is used to retrieve the mode setting that conveys how rules will be kept up to date.
	// Deprecated: this struct will be removed in a future release.
	GetEvalsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	// GetEvalsResponse is returned from a call to GetEvalsResponse.
	// Deprecated: this struct will be removed in a future release.
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
// Deprecated: this method will be removed in a future release.
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
	logger := p.Log(ctx)
	logger.Debug("GetEval")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/mode",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEval request: %w", err)
	}

	var result GetEvalResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get eval request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

// Deprecated: this method will be removed in a future release.
func (p *appsec) GetEvals(ctx context.Context, params GetEvalsRequest) (*GetEvalsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetEvals")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/mode",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvals request: %w", err)
	}

	var result GetEvalsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get evals request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateEval(ctx context.Context, params UpdateEvalRequest) (*UpdateEvalResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateEval")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateEval request: %w", err)
	}

	var result UpdateEvalResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update eval request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveEval(ctx context.Context, params RemoveEvalRequest) (*RemoveEvalResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveEval")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveEval request: %w", err)
	}

	var result RemoveEvalResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove eval request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
