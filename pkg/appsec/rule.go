package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Rule represents a collection of Rule, with Action and optional Condition/Exception information
//
// See: Rule.GetRule()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// Rule  contains operations available on Rule  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getRule
	Rule interface {
		GetRules(ctx context.Context, params GetRulesRequest) (*GetRulesResponse, error)
		GetRule(ctx context.Context, params GetRuleRequest) (*GetRuleResponse, error)
		UpdateRule(ctx context.Context, params UpdateRuleRequest) (*UpdateRuleResponse, error)
	}

	RuleConditionException struct {
		Conditions *RuleConditions `json:"conditions,omitempty"`
		Exception  *RuleException  `json:"exception,omitempty"`
	}

	RuleConditions []struct {
		Type          string   `json:"type,omitempty"`
		Extensions    []string `json:"extensions,omitempty"`
		Filenames     []string `json:"filenames,omitempty"`
		Hosts         []string `json:"hosts,omitempty"`
		Ips           []string `json:"ips,omitempty"`
		Methods       []string `json:"methods,omitempty"`
		Paths         []string `json:"paths,omitempty"`
		Header        string   `json:"header,omitempty"`
		CaseSensitive bool     `json:"caseSensitive,omitempty"`
		Name          string   `json:"name,omitempty"`
		NameCase      bool     `json:"nameCase,omitempty"`
		PositiveMatch bool     `json:"positiveMatch"`
		Value         string   `json:"value,omitempty"`
		Wildcard      bool     `json:"wildcard,omitempty"`
		ValueCase     bool     `json:"valueCase,omitempty"`
		ValueWildcard bool     `json:"valueWildcard,omitempty"`
		UseHeaders    bool     `json:"useHeaders,omitempty"`
	}

	RuleException struct {
		AnyHeaderCookieOrParam               []string                                 `json:"anyHeaderCookieOrParam,omitempty"`
		HeaderCookieOrParamValues            []string                                 `json:"headerCookieOrParamValues,omitempty"`
		SpecificHeaderCookieOrParamNameValue *SpecificHeaderCookieOrParamNameValuePtr `json:"specificHeaderCookieOrParamNameValue,omitempty"`
		SpecificHeaderCookieOrParamNames     *SpecificHeaderCookieOrParamNamesPtr     `json:"specificHeaderCookieOrParamNames,omitempty"`
		SpecificHeaderCookieOrParamPrefix    *SpecificHeaderCookieOrParamPrefixPtr    `json:"specificHeaderCookieOrParamPrefix,omitempty"`
	}

	SpecificHeaderCookieOrParamNamesPtr []struct {
		Names    []string `json:"names,omitempty"`
		Selector string   `json:"selector,omitempty"`
	}

	SpecificHeaderCookieOrParamPrefixPtr struct {
		Prefix   string `json:"prefix,omitempty"`
		Selector string `json:"selector,omitempty"`
	}

	SpecificHeaderCookieOrParamNameValuePtr struct {
		Name     string `json:"name,omitempty"`
		Selector string `json:"selector,omitempty"`
		Value    string `json:"value,omitempty"`
	}

	GetRuleRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetRuleResponse struct {
		Action             string                  `json:"action,omitempty"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}

	GetRulesRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetRulesResponse struct {
		Rules []struct {
			ID                 int                     `json:"id,omitempty"`
			Action             string                  `json:"action,omitempty"`
			ConditionException *RuleConditionException `json:"conditionException,omitempty"`
		} `json:"ruleActions,omitempty"`
	}

	UpdateRuleRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		RuleID         int             `json:"-"`
		Action         string          `json:"action"`
		JsonPayloadRaw json.RawMessage `json:"conditionException,omitempty"`
	}

	UpdateRuleResponse struct {
		Action             string                  `json:"action,omitempty"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}
)

// Check Condition Exception is Empty
func (r *GetRuleResponse) IsEmptyConditionException() bool {

	if r.ConditionException == nil {
		return true
	}
	return false
}

// Validate validates GetRuleRequest
func (v GetRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates GetRulesRequest
func (v GetRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRuleRequest
func (v UpdateRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetRule(ctx context.Context, params GetRuleRequest) (*GetRuleResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRule")

	var rval GetRuleResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getRule request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetRules(ctx context.Context, params GetRulesRequest) (*GetRulesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRules")

	var rval GetRulesResponse
	var rvalfiltered GetRulesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getRules request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getRules request failed: %w", err)
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

// Update will update a Rule.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putRule

func (p *appsec) UpdateRule(ctx context.Context, params UpdateRuleRequest) (*UpdateRuleResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRule")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d/action-condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create Rule request: %w", err)
	}

	var rval UpdateRuleResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create Rule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
