package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// EvalRuleAction represents a collection of EvalRuleAction
//
// See: EvalRuleAction.GetEvalRuleAction()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// EvalRuleAction  contains operations available on EvalRuleAction  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevalruleaction
	EvalRuleAction interface {
		GetEvalRuleActions(ctx context.Context, params GetEvalRuleActionsRequest) (*GetEvalRuleActionsResponse, error)
		GetEvalRuleAction(ctx context.Context, params GetEvalRuleActionRequest) (*GetEvalRuleActionResponse, error)
		UpdateEvalRuleAction(ctx context.Context, params UpdateEvalRuleActionRequest) (*UpdateEvalRuleActionResponse, error)
	}

	GetEvalRuleActionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetEvalRuleActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetEvalRuleActionsResponse struct {
		RuleActions []struct {
			Action string `json:"action"`
			ID     int    `json:"id"`
		} `json:"evalRuleActions"`
	}

	GetEvalRuleActionResponse struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}

	UpdateEvalRuleActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"ruleId"`
		Action   string `json:"action"`
	}

	UpdateEvalRuleActionResponse struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}
)

// Validate validates GetEvalRuleActionRequest
func (v GetEvalRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates GetEvalRuleActionsRequest
func (v GetEvalRuleActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateEvalRuleActionRequest
func (v UpdateEvalRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetEvalRuleAction(ctx context.Context, params GetEvalRuleActionRequest) (*GetEvalRuleActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalRuleAction")

	var rval GetEvalRuleActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalruleaction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalruleaction  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetEvalRuleActions(ctx context.Context, params GetEvalRuleActionsRequest) (*GetEvalRuleActionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalRuleActions")

	var rval GetEvalRuleActionsResponse
	var rvalfiltered GetEvalRuleActionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalruleactions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalruleactions request failed: %w", err)
	}
	logger.Debugf("GetEvalRuleActions %v %v", params.RuleID, resp)
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RuleID != 0 {
		for _, val := range rval.RuleActions {
			if val.ID == params.RuleID {
				rvalfiltered.RuleActions = append(rvalfiltered.RuleActions, val)
			}
		}
	} else {
		rvalfiltered = rval
	}
	return &rvalfiltered, nil

}

// Update will update a EvalRuleAction.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalruleaction

func (p *appsec) UpdateEvalRuleAction(ctx context.Context, params UpdateEvalRuleActionRequest) (*UpdateEvalRuleActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEvalRuleAction")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create EvalRuleActionrequest: %w", err)
	}

	var rval UpdateEvalRuleActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create EvalRuleAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
