// Package papi provides access to the Akamai Property APIs
package papi

import (
	"errors"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/spf13/cast"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")

	// ErrNotFound is returned when requested resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrSBDNotEnabled indicates that secure-by-default is not enabled on the given account
	ErrSBDNotEnabled = errors.New("secure-by-default is not enabled")

	// ErrDefaultCertLimitReached indicates that the limit for DEFAULT certificates has been reached
	ErrDefaultCertLimitReached = errors.New("the limit for DEFAULT certificates has been reached")
)

type (
	// PAPI is the papi api interface
	PAPI interface {
		Activations
		ClientSettings
		Contracts
		CPCodes
		EdgeHostnames
		Groups
		Includes
		IncludeRules
		IncludeActivations
		IncludeVersions
		Products
		Properties
		PropertyRules
		PropertyVersionHostnames
		PropertyVersions
		RuleFormats
		Search
	}

	papi struct {
		session.Session
		usePrefixes bool
	}

	// Option defines a PAPI option
	Option func(*papi)

	// ClientFunc is a papi client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) PAPI

	// Response is a base PAPI Response type
	Response struct {
		AccountID  string   `json:"accountId,omitempty"`
		ContractID string   `json:"contractId,omitempty"`
		GroupID    string   `json:"groupId,omitempty"`
		Etag       string   `json:"etag,omitempty"`
		Errors     []*Error `json:"errors,omitempty"`
		Warnings   []*Error `json:"warnings,omitempty"`
	}
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
// See: https://techdocs.akamai.com/property-mgr/reference/id-prefixes
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
