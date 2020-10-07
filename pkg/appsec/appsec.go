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
		Configuration
		ConfigurationClone
		ConfigurationVersion
		CustomRule
		CustomRuleAction
		ExportConfiguration
		MatchTarget
		RatePolicy
		RatePolicyAction
		SecurityPolicy
		SecurityPolicyClone
		SelectedHostname
		SelectableHostnames
		SlowPostProtectionSetting
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
		Session:     sess,
		usePrefixes: true,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

/*
// WithUsePrefixes sets the `PAPI-Use-Prefixes` header on requests
// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#prefixes
func WithUsePrefixes(usePrefixes bool) Option {
	return func(p *papi) {
		p.usePrefixes = usePrefixes
	}
}

// Exec overrides the session.Exec to add papi options
func (p *papi) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	// explicitly add the PAPI-Use-Prefixes header
	r.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	return p.Session.Exec(r, out, in...)
}
*/
