// Package appsec provides access to the Akamai Application Security APIs
package appsec

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")

	// ErrRequestCreation is returned when creating an HTTP request failed
	ErrRequestCreation = errors.New("HTTP request failure")

	// ErrAPICallFailure is returned when an Appsec OpenAPI call failed
	ErrAPICallFailure = errors.New("API call failure")
)

type (
	// APPSEC is the appsec api interface
	APPSEC interface {
		Activations
		AdvancedSettingsAsePenaltyBox
		AdvancedSettingsAttackPayloadLogging
		AdvancedSettingsEvasivePathMatch
		AdvancedSettingsLogging
		AdvancedSettingsPIILearning
		AdvancedSettingsPragma
		AdvancedSettingsPrefetch
		AdvancedSettingsRequestBody
		AdvancedSettingsJA4Fingerprint
		AdvancedSettingsURLEvasionDefense
		ApiConstraintsProtection
		ApiEndpoints
		ApiHostnameCoverage
		ApiHostnameCoverageMatchTargets
		ApiHostnameCoverageOverlapping
		ApiRequestConstraints
		AttackGroup
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
		EvalPenaltyBox
		EvalPenaltyBoxConditions
		EvalRule
		ExportConfiguration
		FailoverHostnames
		HostMoveActivations
		IPGeo
		IPGeoProtection
		MalwareContentTypes
		MalwarePolicy
		MalwarePolicyAction
		MalwareProtection
		MatchTarget
		MatchTargetSequence
		PenaltyBox
		PenaltyBoxConditions
		PolicyProtections
		RapidRule
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
		URLProtection
		URLProtectionAction
		VersionNotes
		WAFCompositeRuleset
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
