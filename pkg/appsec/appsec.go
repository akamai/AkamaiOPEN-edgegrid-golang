// Package appsec provides access to the Akamai Application Security APIs
package appsec

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// APPSEC is the appsec api interface
	APPSEC interface {
		Activations
		AdvancedSettingsLogging
		AdvancedSettingsPrefetch
		ApiEndpoints
		ApiHostnameCoverage
		ApiHostnameCoverageOverlapping
		ApiHostnameCoverageMatchTargets
		ApiRequestConstraints
		AttackGroupAction
		AttackGroupConditionException
		BypassNetworkLists
		Configuration
		ConfigurationClone
		ConfigurationVersionClone
		ConfigurationVersion
		ContractsGroups
		CustomDeny
		CustomRule
		CustomRuleAction
		Eval
		EvalHost
		EvalProtectHost
		EvalRuleAction
		EvalRuleConditionException
		ExportConfiguration
		FailoverHostnames
		IPGeo
		ReputationAnalysis
		RuleAction
		RuleConditionException
		MatchTarget
		MatchTargetSequence
		NetworkLayerProtection
		PenaltyBox
		PolicyProtections
		RatePolicy
		RatePolicyAction
		RateProtection
		ReputationProfile
		ReputationProfileAction
		ReputationProtection
		RuleUpgrade
		SecurityPolicy
		SecurityPolicyClone
		SelectedHostname
		SelectableHostnames
		SiemDefinitions
		SiemSettings
		SlowPostProtectionSetting
		SlowPostProtection
		VersionNotes
		WAFMode
		WAFProtection
	}

	appsec struct {
		session.Session
		usePrefixes bool
	}

	// Option defines a PAPI option
	Option func(*appsec)

	// ClientFunc is a appsec client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) APPSEC
)

// Client returns a new appsec Client instance with the specified controller
func Client(sess session.Session, opts ...Option) APPSEC {
	p := &appsec{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
