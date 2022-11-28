// Package appsec provides access to the Akamai Application Security APIs
package appsec

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// APPSEC is the appsec api interface
	APPSEC interface {
		Activations
		AdvancedSettingsEvasivePathMatch
		AdvancedSettingsLogging
		AdvancedSettingsPragma
		AdvancedSettingsPrefetch
		ApiConstraintsProtection
		ApiEndpoints
		ApiHostnameCoverage
		ApiHostnameCoverageMatchTargets
		ApiHostnameCoverageOverlapping
		ApiRequestConstraints
		AttackGroup
		BypassNetworkLists
		Configuration
		ConfigurationClone
		ConfigurationVersion
		ConfigurationVersionClone
		ContractsGroups
		CustomDeny
		CustomRule
		CustomRuleAction
		Eval
		EvalGroup
		EvalHost
		EvalPenaltyBox
		EvalProtectHost
		EvalRule
		ExportConfiguration
		FailoverHostnames
		IPGeo
		IPGeoProtection
		MalwareContentTypes
		MalwarePolicy
		MalwarePolicyAction
		MalwareProtection
		MatchTarget
		MatchTargetSequence
		NetworkLayerProtection
		PenaltyBox
		PolicyProtections
		RatePolicy
		RatePolicyAction
		RateProtection
		ReputationAnalysis
		ReputationProfile
		ReputationProfileAction
		ReputationProtection
		Rule
		RuleUpgrade
		SecurityPolicy
		SecurityPolicyClone
		SelectableHostnames
		SelectedHostname
		SiemDefinitions
		SiemSettings
		SlowPostProtection
		SlowPostProtectionSetting
		ThreatIntel
		TuningRecommendations
		VersionNotes
		WAFMode
		WAFProtection
		WAPBypassNetworkLists
		WAPSelectedHostnames
	}

	appsec struct {
		session.Session
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
