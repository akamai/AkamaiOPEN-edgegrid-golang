package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RapidRule interface supports retrieving and modifying the rapid rules in a policy together with their
	// actions, actions lock, exceptions to a specific rapid rule.
	RapidRule interface {

		// GetRapidRules returns the action taken for each rule in a policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rapid-rules
		GetRapidRules(ctx context.Context, params GetRapidRulesRequest) (*GetRapidRulesResponse, error)

		// GetRapidRulesDefaultAction returns the rapid ruleset default action.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rapid-rules-action
		GetRapidRulesDefaultAction(ctx context.Context, params GetRapidRulesDefaultActionRequest) (*GetRapidRulesDefaultActionResponse, error)

		// GetRapidRulesStatus returns the information whether rapid rules is enabled or disabled.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rapid-rules-status
		GetRapidRulesStatus(ctx context.Context, params GetRapidRulesStatusRequest) (*GetRapidRulesStatusResponse, error)

		// UpdateRapidRulesStatus enables or disables the rapid rules.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rapid-rules-status
		UpdateRapidRulesStatus(ctx context.Context, params UpdateRapidRulesStatusRequest) (*UpdateRapidRulesStatusResponse, error)

		// UpdateRapidRulesDefaultAction updates a rapid rules default action.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rapid-rules-action
		UpdateRapidRulesDefaultAction(ctx context.Context, params UpdateRapidRulesDefaultActionRequest) (*UpdateRapidRulesDefaultActionResponse, error)

		// UpdateRapidRuleActionLock updates a rapid rule action lock.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rapid-rule-lock
		UpdateRapidRuleActionLock(ctx context.Context, params UpdateRapidRuleActionLockRequest) (*UpdateRapidRuleActionLockResponse, error)

		// UpdateRapidRuleAction updates what action a rapid rule takes when it's triggered.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rapid-rule-action
		UpdateRapidRuleAction(ctx context.Context, params UpdateRapidRuleActionRequest) (*UpdateRapidRuleActionResponse, error)

		// UpdateRapidRuleException updates a rapid rule exception.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rapid-rule-condition-exception
		UpdateRapidRuleException(ctx context.Context, params UpdateRapidRuleExceptionRequest) (*UpdateRapidRuleExceptionResponse, error)
	}

	// GetRapidRulesRequest is used to retrieve the rapid rules for a configuration and policy,
	// together with their actions, actions lock and condition and exception information.
	GetRapidRulesRequest struct {
		ConfigID int64
		Version  int
		PolicyID string
		RuleID   *int64
	}

	// GetRapidRulesResponse is returned from a call to GetRapidRules.
	GetRapidRulesResponse struct {
		Rules []PolicyRapidRule `json:"policyRules"`
	}

	// PolicyRapidRule represents a rapid rule returned as part of GetRapidRulesResponse.
	PolicyRapidRule struct {
		ID                 int64                   `json:"id"`
		Action             string                  `json:"action"`
		Lock               bool                    `json:"lock"`
		Name               string                  `json:"title"`
		Version            int                     `json:"version"`
		RiskScoreGroups    []string                `json:"riskScoreGroups"`
		ConditionException *RuleConditionException `json:"conditionException"`
	}

	// GetRapidRulesDefaultActionRequest is used to retrieve the rapid rules default action.
	GetRapidRulesDefaultActionRequest GetRapidRulesRequest

	// GetRapidRulesDefaultActionResponse is returned from a call to GetRapidRulesDefaultAction.
	GetRapidRulesDefaultActionResponse struct {
		Action string `json:"action"`
	}

	// GetRapidRulesStatusRequest is used to retrieve the rapid rules status (enabled/disabled).
	GetRapidRulesStatusRequest GetRapidRulesRequest

	// GetRapidRulesStatusResponse is returned from a call to GetRapidRulesStatus.
	GetRapidRulesStatusResponse struct {
		Enabled bool `json:"enabled"`
	}

	// UpdateRapidRulesStatusRequestBody body for rapid ruleset status update request
	UpdateRapidRulesStatusRequestBody struct {
		Enabled *bool `json:"enabled"`
	}

	// UpdateRapidRulesStatusRequest is used to modify the status (enabled/disabled) for a rapid ruleset.
	UpdateRapidRulesStatusRequest struct {
		ConfigID int64
		Version  int
		PolicyID string
		Body     UpdateRapidRulesStatusRequestBody
	}

	// UpdateRapidRulesStatusResponse is returned from a call to UpdateRapidRulesStatus.
	UpdateRapidRulesStatusResponse GetRapidRulesStatusResponse

	// UpdateRapidRulesDefaultActionRequestBody body for rapid rules default action update request
	UpdateRapidRulesDefaultActionRequestBody struct {
		Action string `json:"action"`
	}

	// UpdateRapidRulesDefaultActionRequest is used to modify a new rapid rules default action.
	UpdateRapidRulesDefaultActionRequest struct {
		ConfigID int64
		Version  int
		PolicyID string
		Body     UpdateRapidRulesDefaultActionRequestBody
	}

	// UpdateRapidRulesDefaultActionResponse is returned from a call to UpdateRapidRulesDefaultAction.
	UpdateRapidRulesDefaultActionResponse GetRapidRulesDefaultActionResponse

	// UpdateRapidRuleActionLockRequestBody body for rapid rule action lock update request
	UpdateRapidRuleActionLockRequestBody struct {
		Enabled *bool `json:"enabled"`
	}

	// UpdateRapidRuleActionLockRequest is used to modify a rapid rule action lock.
	UpdateRapidRuleActionLockRequest struct {
		ConfigID int64
		Version  int
		PolicyID string
		RuleID   int64
		Body     UpdateRapidRuleActionLockRequestBody
	}

	// UpdateRapidRuleActionLockResponse is returned from a call to UpdateRapidRuleActionLock.
	UpdateRapidRuleActionLockResponse GetRapidRulesStatusResponse

	// UpdateRapidRuleActionRequestBody body for rapid rule's action update request
	UpdateRapidRuleActionRequestBody struct {
		Action string `json:"action"`
	}

	// UpdateRapidRuleActionRequest is used to modify a rapid rule action.
	UpdateRapidRuleActionRequest struct {
		ConfigID    int64
		Version     int
		PolicyID    string
		RuleID      int64
		RuleVersion int
		Body        UpdateRapidRuleActionRequestBody
	}

	// UpdateRapidRuleActionResponse is returned from a call to UpdateRapidRuleAction.
	UpdateRapidRuleActionResponse struct {
		Action string `json:"action"`
		Lock   bool   `json:"lock"`
	}

	// UpdateRapidRuleExceptionRequest is used to modify a rapid rule exception.
	UpdateRapidRuleExceptionRequest struct {
		ConfigID int64
		Version  int
		PolicyID string
		RuleID   int64
		Body     RuleConditionException
	}

	// UpdateRapidRuleExceptionResponse is returned from a call to UpdateRapidRuleException.
	UpdateRapidRuleExceptionResponse RuleConditionException

	// RapidRuleDetails represents a rule details with an ID, action, action lock, attack group, optional condition exception and attack group exception.
	RapidRuleDetails struct {
		ID                   int64                          `json:"id"`
		Action               string                         `json:"action"`
		Lock                 bool                           `json:"lock"`
		Name                 string                         `json:"name"`
		AttackGroup          string                         `json:"attack_group,omitempty"`
		AttackGroupException *AttackGroupConditionException `json:"attack_group_exception,omitempty"`
		ConditionException   *RuleConditionException        `json:"condition_exception,omitempty"`
	}

	// RuleDefinition represents a rule configuration with an ID, action, action lock and optional condition exception.
	RuleDefinition struct {
		ID                 *int64                  `json:"id"`
		Action             *string                 `json:"action"`
		Lock               *bool                   `json:"lock"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}
)

// Validate validates a GetRapidRulesRequest.
func (v GetRapidRulesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	})
}

// Validate validates a GetRapidRulesDefaultActionRequest.
func (v GetRapidRulesDefaultActionRequest) Validate() error {
	return GetRapidRulesRequest(v).Validate()
}

// Validate validates a GetRapidRulesStatusRequest.
func (v GetRapidRulesStatusRequest) Validate() error {
	return GetRapidRulesRequest(v).Validate()
}

// Validate validates a UpdateRapidRulesStatusRequest.
func (v UpdateRapidRulesStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Body":     validation.Validate(v.Body, validation.Required),
	})
}

// Validate validates a UpdateRapidRulesStatusRequestBody.
func (v UpdateRapidRulesStatusRequestBody) Validate() error {
	return validation.Errors{
		"Enabled": validation.Validate(v.Enabled, validation.NotNil),
	}.Filter()
}

// Validate validates a UpdateRapidRulesDefaultActionRequest.
func (v UpdateRapidRulesDefaultActionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Body":     validation.Validate(v.Body, validation.Required),
	})
}

// Validate validates a UpdateRapidRulesDefaultActionRequestBody.
func (v UpdateRapidRulesDefaultActionRequestBody) Validate() error {
	return validation.Errors{
		"Action": validation.Validate(v.Action, validation.Required),
	}.Filter()
}

// Validate validates a UpdateRapidRuleActionLockRequest.
func (v UpdateRapidRuleActionLockRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
		"Body":     validation.Validate(v.Body, validation.Required),
	})
}

// Validate validates a UpdateRapidRuleActionLockRequestBody.
func (v UpdateRapidRuleActionLockRequestBody) Validate() error {
	return validation.Errors{
		"Enabled": validation.Validate(v.Enabled, validation.NotNil),
	}.Filter()
}

// Validate validates a UpdateRapidRuleActionRequest.
func (v UpdateRapidRuleActionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"PolicyID":    validation.Validate(v.PolicyID, validation.Required),
		"RuleID":      validation.Validate(v.RuleID, validation.Required),
		"RuleVersion": validation.Validate(v.RuleVersion, validation.Required),
		"Body":        validation.Validate(v.Body, validation.Required),
	})
}

// Validate validates a UpdateRapidRuleActionRequestBody.
func (v UpdateRapidRuleActionRequestBody) Validate() error {
	return validation.Errors{
		"Action": validation.Validate(v.Action, validation.Required),
	}.Filter()
}

// Validate validates a UpdateRapidRuleExceptionRequest.
func (v UpdateRapidRuleExceptionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
		"Body":     validation.Validate(v.Body, validation.Required),
	})
}

func (p *appsec) GetRapidRules(ctx context.Context, params GetRapidRulesRequest) (*GetRapidRulesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRapidRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rapid-rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRapidRules request: %w", err)
	}

	var result GetRapidRulesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rapid rules request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RuleID != nil {
		for _, val := range result.Rules {
			if val.ID == *params.RuleID {
				return &GetRapidRulesResponse{Rules: []PolicyRapidRule{val}}, nil
			}
		}
		return nil, fmt.Errorf("get rapid rule failure. rapid rule with ID: %d not found", *params.RuleID)
	}

	return &result, nil
}

func (p *appsec) GetRapidRulesDefaultAction(ctx context.Context, params GetRapidRulesDefaultActionRequest) (*GetRapidRulesDefaultActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRapidRulesDefaultAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getRapidRulesDefaultActionURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRapidRulesDefaultAction request: %w", err)
	}

	var result GetRapidRulesDefaultActionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rapid rules default action request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetRapidRulesStatus(ctx context.Context, params GetRapidRulesStatusRequest) (*GetRapidRulesStatusResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRapidRulesStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getRapidRulesStatusURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRapidRulesStatus request: %w", err)
	}

	var result GetRapidRulesStatusResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rapid rules status request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateRapidRulesStatus(ctx context.Context, params UpdateRapidRulesStatusRequest) (*UpdateRapidRulesStatusResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRapidRulesStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getRapidRulesStatusURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRapidRulesStatus request: %w", err)
	}

	var result UpdateRapidRulesStatusResponse
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("update rapid rules status request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateRapidRulesDefaultAction(ctx context.Context, params UpdateRapidRulesDefaultActionRequest) (*UpdateRapidRulesDefaultActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRapidRulesDefaultAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getRapidRulesDefaultActionURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRapidRulesDefaultActionRequest request: %w", err)
	}

	var result UpdateRapidRulesDefaultActionResponse
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("update rapid rules default action request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateRapidRuleActionLock(ctx context.Context, params UpdateRapidRuleActionLockRequest) (*UpdateRapidRuleActionLockResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRapidRuleActionLock")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rapid-rules/%d/lock",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRapidRuleActionLockRequest request: %w", err)
	}

	var result UpdateRapidRuleActionLockResponse
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("update rapid rule action lock request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateRapidRuleAction(ctx context.Context, params UpdateRapidRuleActionRequest) (*UpdateRapidRuleActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRapidRuleAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rapid-rules/%d/versions/%d/action",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
		params.RuleVersion)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRapidRuleActionRequest request: %w", err)
	}

	var result UpdateRapidRuleActionResponse
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("update rapid rule action request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateRapidRuleException(ctx context.Context, params UpdateRapidRuleExceptionRequest) (*UpdateRapidRuleExceptionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRapidRuleException")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rapid-rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRapidRuleExceptionRequest request: %w", err)
	}

	var result UpdateRapidRuleExceptionResponse
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("update rapid rule exception request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func getRapidRulesStatusURI(configID int64, version int, policyID string) string {
	return fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rapid-rules/status",
		configID,
		version,
		policyID)
}

func getRapidRulesDefaultActionURI(configID int64, version int, policyID string) string {
	return fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rapid-rules/action",
		configID,
		version,
		policyID)
}
