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
		AttackGroupAction
		AttackGroupConditionException
		Configuration
		ConfigurationClone
		ConfigurationVersion
		CustomRule
		CustomRuleAction
		Eval
		EvalRuleAction
		EvalRuleConditionException
		ExportConfiguration
		IPGeo
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
		SlowPostProtectionSetting
		SlowPostProtection
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
