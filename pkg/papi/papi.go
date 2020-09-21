// Package papi provides access to the Akamai Property APIs
package papi

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// PAPI is the papi api interface
	PAPI interface {
		Groups
		Contracts
		Activations
		CPCodes
		Properties
		PropertyVersions
		EdgeHostnames
		Products
		Search
		PropertyVersionHostnames
		ClientSettings
		PropertyRules
	}

	papi struct {
		session.Session
		usePrefixes bool
	}

	// Option defines a PAPI option
	Option func(*papi)

	// ClientFunc is a papi client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) PAPI
)

// Client returns a new papi Client instance with the specified controller
func Client(sess session.Session, opts ...Option) PAPI {
	p := &papi{
		Session:     sess,
		usePrefixes: true,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// WithUsePrefixes sets the `PAPI-Use-Prefixes` header on requests
// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#prefixes
func WithUsePrefixes(usePrefixes bool) Option {
	return func(p *papi) {
		p.usePrefixes = usePrefixes
	}
}
