package appsec

type (
	// RulesetType is a ruleset type value.
	RulesetType string
	ActionType  string
)

const (
	// RulesetTypeActive for active rulesets.
	RulesetTypeActive RulesetType = "active"

	// RulesetTypeEvaluation for evaluation rulesets.
	RulesetTypeEvaluation RulesetType = "evaluation"

	// ActionTypeDeny firewall deny action.
	ActionTypeDeny ActionType = "deny"
	// ActionTypeAlert firewall alert action.
	ActionTypeAlert ActionType = "alert"
	// ActionTypeNone firewall no action.
	ActionTypeNone ActionType = "none"
)
