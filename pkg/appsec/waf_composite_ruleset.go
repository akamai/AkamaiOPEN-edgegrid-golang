package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The WAFCompositeRuleset interface supports retrieving and updating the complete WAF ruleset
	WAFCompositeRuleset interface {
		// GetWAFCompositeRuleset returns the complete WAF ruleset configuration including adaptive
		// intelligence, attack groups with their actions and condition exceptions, and rules with their
		// actions and condition exceptions.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-waf-policy-ruleset-composite
		GetWAFCompositeRuleset(ctx context.Context, params GetWAFCompositeRulesetRequest) (*CompositeRulesetResponse, error)

		// UpdateWAFCompositeRuleset updates the WAF ruleset configuration including attack groups
		// and rules together with their actions and conditions.
		//
		// See: https://techdocs.akamai.com/application-security/reference/patch-waf-policy-ruleset-composite
		UpdateWAFCompositeRuleset(ctx context.Context, params UpdateWAFCompositeRulesetRequest) (*CompositeRulesetResponse, error)
	}

	// GetWAFCompositeRulesetRequest is used to retrieve the complete WAF ruleset configuration.
	GetWAFCompositeRulesetRequest struct {
		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// Version is the version number of the security configuration
		Version int64

		// PolicyID is the unique identifier of the security policy
		PolicyID string
	}

	// WAFRulesetInfo contains information about the ruleset version.
	WAFRulesetInfo struct {
		// RulesetVersion is the version number of the WAF ruleset
		RulesetVersion int64 `json:"rulesetVersion,omitempty"`
	}

	// AdaptiveIntelligence contains adaptive intelligence settings for the WAF.
	AdaptiveIntelligence struct {
		// Indicates whether threat intelligence is enabled or not
		ThreatIntel bool `json:"threatIntel,omitempty"`
	}

	// CompositeRulesetResponse is returned from a call to GetWAFCompositeRuleset.
	CompositeRulesetResponse struct {
		// RulesetInfo contains information about the ruleset version
		RulesetInfo *WAFRulesetInfo `json:"rulesetInfo,omitempty"`

		// AdaptiveIntelligence contains adaptive intelligence settings for the WAF
		AdaptiveIntelligence *AdaptiveIntelligence `json:"adaptiveIntelligence,omitempty"`

		// AttackGroups is the list of attack group configurations
		AttackGroups []WAFCompositeAttackGroup `json:"attackGroups,omitempty"`

		// Rules is the list of rule objects including action and condition exceptions
		Rules []WAFCompositeRule `json:"rules,omitempty"`
	}

	// WAFCompositeAttackGroup represents an attack group with its action and condition exception.
	WAFCompositeAttackGroup struct {
		// Group is the unique name of the attack group
		Group string `json:"group,omitempty"`

		// Action taken anytime the attack group is triggered. Possible values: alert, deny, deny_custom_{custom_deny_id}, none
		Action string `json:"action,omitempty"`

		// ConditionException contains conditions and exceptions associated with the attack group
		ConditionException *AttackGroupConditionException `json:"conditionException,omitempty"`
	}

	// WAFCompositeRule represents a rule with its action and condition exception.
	WAFCompositeRule struct {
		// RuleID is the unique identifier for the rule
		RuleID int64 `json:"ruleId,omitempty"`

		// RuleName is the descriptive name of the rule
		RuleName string `json:"ruleName,omitempty"`

		// Action taken anytime the rule is triggered. Possible values: alert, deny, deny_custom_{custom_deny_id}, none
		Action string `json:"action,omitempty"`

		// ConditionException contains conditions and exceptions associated with the rule
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}

	// UpdateWAFCompositeRulesetRequest is used to update the WAF ruleset configuration.
	UpdateWAFCompositeRulesetRequest struct {
		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// Version is the version number of the security configuration
		Version int64

		// PolicyID is the unique identifier of the security policy
		PolicyID string

		// AttackGroups is the list of attack group configurations to update
		AttackGroups []WAFCompositeAttackGroupUpdate `json:"attackGroups,omitempty"`

		// Rules is the list of rule objects including action and condition exceptions to update
		Rules []WAFCompositeRuleUpdate `json:"rules,omitempty"`
	}

	// WAFCompositeAttackGroupUpdate represents an attack group update with its action and condition exception.
	WAFCompositeAttackGroupUpdate struct {
		// Group is the unique name of the attack group to update
		Group string `json:"group,omitempty"`

		// Action taken anytime the attack group is triggered. Possible values: alert, deny, deny_custom_{custom_deny_id}, none
		Action string `json:"action,omitempty"`

		// ConditionException contains conditions and exceptions associated with the attack group
		ConditionException *AttackGroupConditionException `json:"conditionException,omitempty"`
	}

	// WAFCompositeRuleUpdate represents a rule update with its action and condition exception.
	WAFCompositeRuleUpdate struct {
		// RuleID is the unique identifier for the rule to update
		RuleID int64 `json:"ruleId,omitempty"`

		// Action taken anytime the rule is triggered. Possible values: alert, deny, deny_custom_{custom_deny_id}, none
		Action string `json:"action,omitempty"`

		// ConditionException contains conditions and exceptions associated with the rule
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}
)

// Validate validates a GetWAFCompositeRulesetRequest.
func (v GetWAFCompositeRulesetRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	})
}

// Validate validates an UpdateWAFCompositeRulesetRequest.
func (v UpdateWAFCompositeRulesetRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	})
}

func (p *appsec) GetWAFCompositeRuleset(ctx context.Context, params GetWAFCompositeRulesetRequest) (*CompositeRulesetResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetWAFCompositeRuleset")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/web-application-firewall/ruleset",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetWAFCompositeRuleset request: %w", err)
	}

	var result CompositeRulesetResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get WAF composite ruleset request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateWAFCompositeRuleset(ctx context.Context, params UpdateWAFCompositeRulesetRequest) (*CompositeRulesetResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateWAFCompositeRuleset")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/web-application-firewall/ruleset",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateWAFCompositeRuleset request: %w", err)
	}

	var result CompositeRulesetResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update WAF composite ruleset request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
