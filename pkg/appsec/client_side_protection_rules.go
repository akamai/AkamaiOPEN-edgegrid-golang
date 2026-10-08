package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ClientSideProtectionRule interface supports managing script behavior rules for a security policy.
	ClientSideProtectionRule interface {
		// ListClientSideProtectionRules returns all script behavior rules for a policy.
		ListClientSideProtectionRules(ctx context.Context, params ListClientSideProtectionRulesRequest) (*ListClientSideProtectionRulesResponse, error)

		// CreateClientSideProtectionRule creates a new script behavior rule.
		CreateClientSideProtectionRule(ctx context.Context, params CreateClientSideProtectionRuleRequest) (*CreateClientSideProtectionRuleResponse, error)

		// UpdateClientSideProtectionRule updates an existing script behavior rule.
		UpdateClientSideProtectionRule(ctx context.Context, params UpdateClientSideProtectionRuleRequest) (*UpdateClientSideProtectionRuleResponse, error)

		// DeleteClientSideProtectionRule deletes a script behavior rule.
		DeleteClientSideProtectionRule(ctx context.Context, params DeleteClientSideProtectionRuleRequest) error
	}

	// ListClientSideProtectionRulesResponse is returned from a call to ListClientSideProtectionRules.
	ListClientSideProtectionRulesResponse []ClientSideProtectionRuleBody

	// CreateClientSideProtectionRuleResponse is returned from a call to CreateClientSideProtectionRule.
	CreateClientSideProtectionRuleResponse = ClientSideProtectionRuleBody

	// UpdateClientSideProtectionRuleResponse is returned from a call to UpdateClientSideProtectionRule.
	UpdateClientSideProtectionRuleResponse = ClientSideProtectionRuleBody

	// ListClientSideProtectionRulesRequest is used to retrieve all script behavior rules for a policy.
	ListClientSideProtectionRulesRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	// CreateClientSideProtectionRuleRequest is used to create a new script behavior rule.
	CreateClientSideProtectionRuleRequest struct {
		ConfigID int                          `json:"-"`
		Version  int                          `json:"-"`
		PolicyID string                       `json:"-"`
		Body     ClientSideProtectionRuleBody `json:"-"`
	}

	// UpdateClientSideProtectionRuleRequest is used to update an existing script behavior rule.
	UpdateClientSideProtectionRuleRequest struct {
		ConfigID int                          `json:"-"`
		Version  int                          `json:"-"`
		PolicyID string                       `json:"-"`
		RuleID   int64                        `json:"-"`
		Body     ClientSideProtectionRuleBody `json:"-"`
	}

	// DeleteClientSideProtectionRuleRequest is used to delete a script behavior rule.
	DeleteClientSideProtectionRuleRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int64  `json:"-"`
	}

	// ClientSideProtectionRuleBody represents a script behavior rule.
	ClientSideProtectionRuleBody struct {
		RuleID       *int64                         `json:"ruleId,omitempty"`
		RuleName     string                         `json:"ruleName"`
		Description  string                         `json:"description,omitempty"`
		Conditions   ClientSideProtectionConditions `json:"conditions"`
		ActivityType string                         `json:"activityType"`
		Action       string                         `json:"action"`
	}

	// ClientSideProtectionConditions holds matching conditions for a script behavior rule.
	ClientSideProtectionConditions struct {
		Page               ClientSideProtectionConditionCriteria       `json:"page"`
		Script             ClientSideProtectionScriptConditionCriteria `json:"script"`
		ActivityTypeValues ClientSideProtectionConditionCriteria       `json:"activityTypeValues"`
	}

	// ClientSideProtectionConditionCriteria holds condition matching criteria.
	ClientSideProtectionConditionCriteria struct {
		Condition string   `json:"condition"`
		Values    []string `json:"values"`
	}

	// ClientSideProtectionScriptConditionCriteria holds script-specific condition criteria.
	ClientSideProtectionScriptConditionCriteria struct {
		ScriptType string   `json:"scriptType,omitempty"`
		Condition  string   `json:"condition"`
		Values     []string `json:"values"`
	}
)

// Validate validates a ListClientSideProtectionRulesRequest.
func (v ListClientSideProtectionRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a CreateClientSideProtectionRuleRequest.
func (v CreateClientSideProtectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":          validation.Validate(v.ConfigID, validation.Required),
		"Version":           validation.Validate(v.Version, validation.Required),
		"PolicyID":          validation.Validate(v.PolicyID, validation.Required),
		"Body.RuleName":     validation.Validate(v.Body.RuleName, validation.Required),
		"Body.ActivityType": validation.Validate(v.Body.ActivityType, validation.Required),
		"Body.Action":       validation.Validate(v.Body.Action, validation.Required),
	}.Filter()
}

// Validate validates an UpdateClientSideProtectionRuleRequest.
func (v UpdateClientSideProtectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID":          validation.Validate(v.ConfigID, validation.Required),
		"Version":           validation.Validate(v.Version, validation.Required),
		"PolicyID":          validation.Validate(v.PolicyID, validation.Required),
		"RuleID":            validation.Validate(v.RuleID, validation.Required),
		"Body.RuleName":     validation.Validate(v.Body.RuleName, validation.Required),
		"Body.ActivityType": validation.Validate(v.Body.ActivityType, validation.Required),
		"Body.Action":       validation.Validate(v.Body.Action, validation.Required),
	}.Filter()
}

// Validate validates a DeleteClientSideProtectionRuleRequest.
func (v DeleteClientSideProtectionRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) ListClientSideProtectionRules(ctx context.Context, params ListClientSideProtectionRulesRequest) (*ListClientSideProtectionRulesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListClientSideProtectionRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewGet(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/client-side-protection-rules",
		params.ConfigID, params.Version, params.PolicyID).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create ListClientSideProtectionRules request: %w", err)
	}

	var result ListClientSideProtectionRulesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("list client-side protection rules request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateClientSideProtectionRule(ctx context.Context, params CreateClientSideProtectionRuleRequest) (*CreateClientSideProtectionRuleResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateClientSideProtectionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewPost(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/client-side-protection-rules",
		params.ConfigID, params.Version, params.PolicyID).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateClientSideProtectionRule request: %w", err)
	}

	var result CreateClientSideProtectionRuleResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("create client-side protection rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateClientSideProtectionRule(ctx context.Context, params UpdateClientSideProtectionRuleRequest) (*UpdateClientSideProtectionRuleResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateClientSideProtectionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewPut(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/client-side-protection-rules/%d",
		params.ConfigID, params.Version, params.PolicyID, params.RuleID).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateClientSideProtectionRule request: %w", err)
	}

	var result UpdateClientSideProtectionRuleResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("update client-side protection rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) DeleteClientSideProtectionRule(ctx context.Context, params DeleteClientSideProtectionRuleRequest) error {
	logger := p.Log(ctx)
	logger.Debug("DeleteClientSideProtectionRule")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewDelete(ctx,
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/client-side-protection-rules/%d",
		params.ConfigID, params.Version, params.PolicyID, params.RuleID).Build()
	if err != nil {
		return fmt.Errorf("failed to create DeleteClientSideProtectionRule request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("delete client-side protection rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}
