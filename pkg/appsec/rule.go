package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The Rule interface supports retrieving and modifying the rules in a policy together with their
	// actions, conditions and exceptions, or the action, condition and exceptions of a specific rule.
	Rule interface {
		// GetRules returns the action taken for each rule in a policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-rules
		GetRules(ctx context.Context, params GetRulesRequest) (*GetRulesResponse, error)

		// GetRule returns the action a rule takes when triggered with conditions and exceptions.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rule-condition-exception-1
		// See: https://techdocs.akamai.com/application-security/reference/get-rule-1
		GetRule(ctx context.Context, params GetRuleRequest) (*GetRuleResponse, error)

		// UpdateRule updates what action a rule takes when it's triggered.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rule-1
		UpdateRule(ctx context.Context, params UpdateRuleRequest) (*UpdateRuleResponse, error)

		// UpdateRuleConditionException updates a rule's conditions and exceptions.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rule-condition-exception-1
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
		Conditions             *RuleConditions     `json:"conditions,omitempty"`
		Exception              *RuleException      `json:"exception,omitempty"`
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
		ClientLists   []string `json:"clientLists,omitempty"`
	}

	// RuleException is used to describe the exceptions for a rule.
	RuleException struct {
		AnyHeaderCookieOrParam                  []string                                 `json:"anyHeaderCookieOrParam,omitempty"`
		HeaderCookieOrParamValues               []string                                 `json:"headerCookieOrParamValues,omitempty"`
		SpecificHeaderCookieOrParamNameValue    *SpecificHeaderCookieOrParamNameValuePtr `json:"specificHeaderCookieOrParamNameValue,omitempty"`
		SpecificHeaderCookieOrParamNames        *SpecificHeaderCookieOrParamNamesPtr     `json:"specificHeaderCookieOrParamNames,omitempty"`
		SpecificHeaderCookieOrParamPrefix       *SpecificHeaderCookieOrParamPrefixPtr    `json:"specificHeaderCookieOrParamPrefix,omitempty"`
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
		ConfigID               int                 `json:"-"`
		Version                int                 `json:"-"`
		PolicyID               string              `json:"-"`
		RuleID                 int                 `json:"-"`
		Conditions             *RuleConditions     `json:"conditions,omitempty"`
		Exception              *RuleException      `json:"exception,omitempty"`
		AdvancedExceptionsList *AdvancedExceptions `json:"advancedExceptions,omitempty"`
	}

	// UpdateConditionExceptionResponse is returned from a call to UpdateConditionException.
	UpdateConditionExceptionResponse RuleConditionException
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
	logger := p.Log(ctx)
	logger.Debug("GetRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	var result GetRuleResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetRules(ctx context.Context, params GetRulesRequest) (*GetRulesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRules request: %w", err)
	}

	var result GetRulesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rules request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RuleID != 0 {
		var filteredResult GetRulesResponse
		for _, val := range result.Rules {
			if val.ID == params.RuleID {
				filteredResult.Rules = append(filteredResult.Rules, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateRule(ctx context.Context, params UpdateRuleRequest) (*UpdateRuleResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d/action-condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRule request: %w", err)
	}

	var result UpdateRuleResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateRuleConditionException(ctx context.Context, params UpdateConditionExceptionRequest) (*UpdateConditionExceptionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRuleConditionException")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRuleConditionException request: %w", err)
	}

	var result UpdateConditionExceptionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update rule condition exception request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
