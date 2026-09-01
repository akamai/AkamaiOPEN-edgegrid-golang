package appsec

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// AIRule interface supports retrieving and modifying AI rules for a security policy.
	AIRule interface {
		// ListAIRules returns all AI rules for a configuration version and security policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-ai-rules (pending publication)
		ListAIRules(ctx context.Context, params ListAIRulesRequest) (*ListAIRulesResponse, error)

		// GetAIRulesStatus returns whether AI rules are enabled or disabled for a policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-ai-rules-status (pending publication)
		GetAIRulesStatus(ctx context.Context, params GetAIRulesStatusRequest) (*GetAIRulesStatusResponse, error)

		// UpdateAIRulesStatus enables or disables AI rules for a policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-ai-rules-status (pending publication)
		UpdateAIRulesStatus(ctx context.Context, params UpdateAIRulesStatusRequest) (*UpdateAIRulesStatusResponse, error)

		// GetAIRuleAction returns the current action for a specific AI rule version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-ai-rule-action (pending publication)
		GetAIRuleAction(ctx context.Context, params GetAIRuleActionRequest) (*GetAIRuleActionResponse, error)

		// UpdateAIRuleAction sets the action for a specific AI rule.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-ai-rule-action (pending publication)
		UpdateAIRuleAction(ctx context.Context, params UpdateAIRuleActionRequest) (*UpdateAIRuleActionResponse, error)
	}

	// ListAIRulesRequest is used to retrieve AI rules for a configuration and policy.
	ListAIRulesRequest struct {
		ConfigID int64
		Version  int
		PolicyID string
	}

	// ListAIRulesResponse is returned from a call to ListAIRules.
	ListAIRulesResponse struct {
		AIRuleStatus string         `json:"aiRuleStatus"`
		AIRules      []PolicyAIRule `json:"aiRules"`
	}

	// PolicyAIRule represents an AI rule returned as part of ListAIRulesResponse.
	PolicyAIRule struct {
		RuleID              int64                      `json:"ruleId"`
		RuleVersion         int64                      `json:"ruleVersion"`
		Title               string                     `json:"title"`
		RiskScoreGroup      string                     `json:"riskScoreGroup"`
		RuleDescription     string                     `json:"ruleDescription"`
		Action              string                     `json:"action"`
		ConditionExceptions []AIRuleConditionException `json:"conditionExceptions"`
	}

	// AIRuleConditionException is a condition exception for an AI rule.
	AIRuleConditionException struct {
		Exception *AIRuleExceptionBody `json:"exception"`
	}

	// AIRuleExceptionBody holds the exception selectors for an AI rule.
	AIRuleExceptionBody struct {
		Selectors []AIRuleExceptionSelector `json:"selectors"`
	}

	// AIRuleExceptionSelector is a single selector entry within an AI rule exception.
	AIRuleExceptionSelector struct {
		Names    []string `json:"names"`
		Selector string   `json:"selector"`
		Wildcard bool     `json:"wildcard"`
		Type     string   `json:"type"`
	}

	// GetAIRulesStatusRequest is used to retrieve the AI rules status (enabled/disabled).
	GetAIRulesStatusRequest ListAIRulesRequest

	// GetAIRulesStatusResponse is returned from a call to GetAIRulesStatus.
	GetAIRulesStatusResponse struct {
		AIRuleStatus string `json:"aiRuleStatus"`
	}

	// UpdateAIRulesStatusRequestBody is the body for an AI rules status update request.
	UpdateAIRulesStatusRequestBody struct {
		AIRuleStatus string `json:"aiRuleStatus"`
	}

	// UpdateAIRulesStatusRequest is used to enable or disable AI rules for a policy.
	UpdateAIRulesStatusRequest struct {
		ConfigID int64
		Version  int
		PolicyID string
		Body     UpdateAIRulesStatusRequestBody
	}

	// UpdateAIRulesStatusResponse is returned from a call to UpdateAIRulesStatus.
	UpdateAIRulesStatusResponse GetAIRulesStatusResponse

	// UpdateAIRuleActionRequestBody is the body for an AI rule action update request.
	UpdateAIRuleActionRequestBody struct {
		Action string `json:"action"`
	}

	// GetAIRuleActionRequest is used to retrieve the action for a specific AI rule version.
	GetAIRuleActionRequest struct {
		ConfigID    int64
		Version     int
		PolicyID    string
		RuleID      int64
		RuleVersion int64
	}

	// UpdateAIRuleActionRequest is used to set the action for a specific AI rule.
	UpdateAIRuleActionRequest struct {
		ConfigID    int64
		Version     int
		PolicyID    string
		RuleID      int64
		RuleVersion int64
		Body        UpdateAIRuleActionRequestBody
	}

	// GetAIRuleActionResponse is returned from a call to GetAIRuleAction.
	GetAIRuleActionResponse struct {
		Action string `json:"action"`
	}

	// UpdateAIRuleActionResponse is returned from a call to UpdateAIRuleAction.
	UpdateAIRuleActionResponse GetAIRuleActionResponse
)

// Validate validates a ListAIRulesRequest.
func (v ListAIRulesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	})
}

// Validate validates a GetAIRulesStatusRequest.
func (v GetAIRulesStatusRequest) Validate() error {
	return ListAIRulesRequest(v).Validate()
}

// Validate validates an UpdateAIRulesStatusRequest.
func (v UpdateAIRulesStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Body":     validation.Validate(v.Body, validation.Required),
	})
}

// Validate validates an UpdateAIRulesStatusRequestBody.
func (v UpdateAIRulesStatusRequestBody) Validate() error {
	return validation.Errors{
		"AIRuleStatus": validation.Validate(v.AIRuleStatus, validation.Required, validation.In("ENABLED", "DISABLED")),
	}.Filter()
}

const (
	// AIRuleIDSQL is the rule ID for the AI-Detected SQL Injection Attack rule.
	AIRuleIDSQL int64 = 3001000
	// AIRuleIDXSS is the rule ID for the AI-Detected XSS Attack rule.
	AIRuleIDXSS int64 = 3001001
)

var validAIRuleIDs = validation.In(AIRuleIDSQL, AIRuleIDXSS)

// Validate validates a GetAIRuleActionRequest.
func (v GetAIRuleActionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"PolicyID":    validation.Validate(v.PolicyID, validation.Required),
		"RuleID":      validation.Validate(v.RuleID, validation.Required, validAIRuleIDs),
		"RuleVersion": validation.Validate(v.RuleVersion, validation.Required),
	})
}

// Validate validates an UpdateAIRuleActionRequest.
func (v UpdateAIRuleActionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"PolicyID":    validation.Validate(v.PolicyID, validation.Required),
		"RuleID":      validation.Validate(v.RuleID, validation.Required, validAIRuleIDs),
		"RuleVersion": validation.Validate(v.RuleVersion, validation.Required),
		"Body":        validation.Validate(v.Body, validation.Required),
	})
}

var validAIRuleAction = regexp.MustCompile(`^(alert|deny|deny_custom_.+|none)$`)

// Validate validates an UpdateAIRuleActionRequestBody.
func (v UpdateAIRuleActionRequestBody) Validate() error {
	return validation.Errors{
		"Action": validation.Validate(v.Action, validation.Required, validation.Match(validAIRuleAction)),
	}.Filter()
}

func (p *appsec) ListAIRules(ctx context.Context, params ListAIRulesRequest) (*ListAIRulesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListAIRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewGet(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ai-rules",
		params.ConfigID, params.Version, params.PolicyID).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create ListAIRules request: %w", err)
	}

	var result ListAIRulesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("list AI rules request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetAIRulesStatus(ctx context.Context, params GetAIRulesStatusRequest) (*GetAIRulesStatusResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAIRulesStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewGet(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ai-rules/status",
		params.ConfigID, params.Version, params.PolicyID).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAIRulesStatus request: %w", err)
	}

	var result GetAIRulesStatusResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get AI rules status request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateAIRulesStatus(ctx context.Context, params UpdateAIRulesStatusRequest) (*UpdateAIRulesStatusResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAIRulesStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewPut(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ai-rules/status",
		params.ConfigID, params.Version, params.PolicyID).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAIRulesStatus request: %w", err)
	}

	var result UpdateAIRulesStatusResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("update AI rules status request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateAIRuleAction(ctx context.Context, params UpdateAIRuleActionRequest) (*UpdateAIRuleActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAIRuleAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewPut(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ai-rules/%d/versions/%d/action",
		params.ConfigID, params.Version, params.PolicyID, params.RuleID, params.RuleVersion).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAIRuleAction request: %w", err)
	}

	var result UpdateAIRuleActionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("update AI rule action request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetAIRuleAction(ctx context.Context, params GetAIRuleActionRequest) (*GetAIRuleActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAIRuleAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewGet(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ai-rules/%d/versions/%d/action",
		params.ConfigID, params.Version, params.PolicyID, params.RuleID, params.RuleVersion).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAIRuleAction request: %w", err)
	}

	var result GetAIRuleActionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get AI rule action request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
