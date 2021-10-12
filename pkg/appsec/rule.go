package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The Rule interface supports retrieving and modifying the rules in a policy together with their
	// actions, conditions and exceptions, or the action, condition and exceptions of a specific rule.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#rule
	Rule interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getrules
		GetRules(ctx context.Context, params GetRulesRequest) (*GetRulesResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getruleaction
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getruleconditionexception
		GetRule(ctx context.Context, params GetRuleRequest) (*GetRuleResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleaction
		UpdateRule(ctx context.Context, params UpdateRuleRequest) (*UpdateRuleResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleconditionexception
		UpdateRuleConditionException(ctx context.Context, params UpdateConditionExceptionRequest) (*UpdateConditionExceptionResponse, error)
	}

	// GetRulesRequest is used to retrieve the rules for a configuration and policy, together with their actions and condition and exception information.
	GetRulesRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	// GetRulesResponse is returned from a call to GetRules.
	GetRulesResponse struct {
		Rules []struct {
			ID                 int                     `json:"id,omitempty"`
			Action             string                  `json:"action,omitempty"`
			ConditionException *RuleConditionException `json:"conditionException,omitempty"`
		} `json:"ruleActions,omitempty"`
	}

	// GetRuleRequest is used to retrieve a rule together with its action and its condition and exception information.
	GetRuleRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	// GetRuleResponse is returned from a call to GetRule.
	GetRuleResponse struct {
		Action             string                  `json:"action,omitempty"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}

	// UpdateRuleRequest is used to modify the settings for a rule.
	UpdateRuleRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		RuleID         int             `json:"-"`
		Action         string          `json:"action"`
		JsonPayloadRaw json.RawMessage `json:"conditionException,omitempty"`
	}

	// UpdateRuleResponse is returned from a call to UpdateRule.
	UpdateRuleResponse struct {
		Action             string                  `json:"action,omitempty"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}

	// AdvancedExceptions is used to describe advanced exceptions used in Adaptive Security Engine(ASE) rules.
	AdvancedExceptions AttackGroupAdvancedExceptions

	// RuleConditionException is used to describe the conditions and exceptions for a rule.
	RuleConditionException struct {
		Conditions *RuleConditions `json:"conditions,omitempty"`
		Exception  *RuleException  `json:"exception,omitempty"`
		AdvancedExceptionsList *AdvancedExceptions `json:"advancedExceptions,omitempty"`
	}

	// RuleConditions is used to describe the conditions for a rule.
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


	// RuleException is used to describe the exceptions for a rule.
	RuleException struct {
		AnyHeaderCookieOrParam               []string                                 `json:"anyHeaderCookieOrParam,omitempty"`
		HeaderCookieOrParamValues            []string                                 `json:"headerCookieOrParamValues,omitempty"`
		SpecificHeaderCookieOrParamNameValue *SpecificHeaderCookieOrParamNameValuePtr `json:"specificHeaderCookieOrParamNameValue,omitempty"`
		SpecificHeaderCookieOrParamNames     *SpecificHeaderCookieOrParamNamesPtr     `json:"specificHeaderCookieOrParamNames,omitempty"`
		SpecificHeaderCookieOrParamPrefix    *SpecificHeaderCookieOrParamPrefixPtr    `json:"specificHeaderCookieOrParamPrefix,omitempty"`
		SpecificHeaderCookieParamXMLOrJSONNames *SpecificHeaderCookieParamXMLOrJSONNames `json:"specificHeaderCookieParamXmlOrJsonNames,omitempty"`
	}

	// SpecificHeaderCookieOrParamNamesPtr is used as part of condition and exception information for a rule.
	SpecificHeaderCookieOrParamNamesPtr []struct {
		Names    []string `json:"names,omitempty"`
		Selector string   `json:"selector,omitempty"`
	}

	// SpecificHeaderCookieOrParamPrefixPtr is used as part of condition and exception information for a rule.
	SpecificHeaderCookieOrParamPrefixPtr struct {
		Prefix   string `json:"prefix,omitempty"`
		Selector string `json:"selector,omitempty"`
	}

	// SpecificHeaderCookieOrParamNameValuePtr is used as part of condition and exception information for a rule.
	SpecificHeaderCookieOrParamNameValuePtr struct {
		Name     string `json:"name,omitempty"`
		Selector string `json:"selector,omitempty"`
		Value    string `json:"value,omitempty"`
	}

	// SpecificHeaderCookieParamXMLOrJSONNames is used as part of condition and exception information for an ASE rule.
	SpecificHeaderCookieParamXMLOrJSONNames AttackGroupSpecificHeaderCookieParamXMLOrJSONNames


	// UpdateConditionExceptionRequest is used to update the condition and exception information for a rule.
	UpdateConditionExceptionRequest struct {
		ConfigID   int             `json:"-"`
		Version    int             `json:"-"`
		PolicyID   string          `json:"-"`
		RuleID     int             `json:"-"`
		Conditions *RuleConditions `json:"conditions,omitempty"`
		Exception  *RuleException  `json:"exception,omitempty"`
	}

	// UpdateConditionExceptionResponse is returned from a call to UpdateConditionException.
	UpdateConditionExceptionResponse struct {
		Conditions *RuleConditions `json:"conditions,omitempty"`
		Exception  *RuleException  `json:"exception,omitempty"`
	}
)

// IsEmptyConditionException checks whether a rule's condition and exception information is empty.
func (r *GetRuleResponse) IsEmptyConditionException() bool {
	return r.ConditionException == nil
}

// Validate validates a GetRuleRequest.
func (v GetRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates a GetRulesRequest.
func (v GetRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateRuleRequest.
func (v UpdateRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateConditionExceptionRequest.
func (v UpdateConditionExceptionRequest) Validate() error {
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
		return nil, fmt.Errorf("failed to create GetRule request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetRule request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create GetRules request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetRules request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create UpdateRule request: %w", err)
	}

	var rval UpdateRuleResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *appsec) UpdateRuleConditionException(ctx context.Context, params UpdateConditionExceptionRequest) (*UpdateConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRuleConditionException")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRuleConditionException request: %w", err)
	}

	var rval UpdateConditionExceptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateRuleConditionException request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
