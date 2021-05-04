package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// EvalRule represents a collection of EvalRule
//
// See: EvalRule.GetEvalRule()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// EvalRule  contains operations available on EvalRule  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevalruleconditionexception
	EvalRule interface {
		GetEvalRules(ctx context.Context, params GetEvalRulesRequest) (*GetEvalRulesResponse, error)
		GetEvalRule(ctx context.Context, params GetEvalRuleRequest) (*GetEvalRuleResponse, error)
		UpdateEvalRule(ctx context.Context, params UpdateEvalRuleRequest) (*UpdateEvalRuleResponse, error)
	}

	GetEvalRulesRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetEvalRulesResponse struct {
		Rules []struct {
			ID                 int                     `json:"id,omitempty"`
			Action             string                  `json:"action,omitempty"`
			ConditionException *RuleConditionException `json:"conditionException,omitempty"`
		} `json:"evalRuleActions,omitempty"`
	}

	GetEvalRuleRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"ruleId"`
	}

	GetEvalRuleResponse struct {
		Action             string                  `json:"action,omitempty"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}

	UpdateEvalRuleRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		RuleID         int             `json:"-"`
		Action         string          `json:"action"`
		JsonPayloadRaw json.RawMessage `json:"conditionException,omitempty"`
	}

	UpdateEvalRuleResponse struct {
		Action             string                  `json:"action,omitempty"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}
)

// Check Condition Exception is Empty
func (r *GetEvalRuleResponse) IsEmptyCodnitionException() bool {

	if r.ConditionException == nil {
		return true
	}
	return false
}

// Validate validates GetEvalRuleRequest
func (v GetEvalRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates GetEvalRulesRequest
func (v GetEvalRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateEvalRuleRequest
func (v UpdateEvalRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetEvalRule(ctx context.Context, params GetEvalRuleRequest) (*GetEvalRuleResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalRule")

	var rval GetEvalRuleResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalruleconditionexception request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalruleconditionexception request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetEvalRules(ctx context.Context, params GetEvalRulesRequest) (*GetEvalRulesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalRules")

	var rval GetEvalRulesResponse
	var rvalfiltered GetEvalRulesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalruleconditionexceptions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalruleconditionexceptions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RuleID != 0 {
		for _, val := range rval.Rules {
			if val.ID == params.RuleID {
				rvalfiltered.Rules = append(rvalfiltered.Rules, val)
			}
		}
	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

// Update will update a EvalRule.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalruleconditionexception

func (p *appsec) UpdateEvalRule(ctx context.Context, params UpdateEvalRuleRequest) (*UpdateEvalRuleResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEvalRule")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d/action-condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create EvalRulerequest: %w", err)
	}

	var rval UpdateEvalRuleResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create EvalRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
