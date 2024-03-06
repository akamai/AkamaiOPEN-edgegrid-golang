// Package dns provides access to the Akamai DNS V2 APIs
//
// See: https://techdocs.akamai.com/edge-dns/reference/edge-dns-api
package dns

import (
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
)

type (
	// DNS is the dns api interface
	DNS interface {
		Authorities
		Data
		Recordsets
		Records
		TSIGKeys
		Zones
	}

	dns struct {
		session.Session
	}

	// Option defines a DNS option
	Option func(*dns)

	// ClientFunc is a dns client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) DNS
)

// Client returns a new dns Client instance with the specified controller
func Client(sess session.Session, opts ...Option) DNS {
	d := &dns{
		Session: sess,
	}

	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Exec overrides the session.Exec to add dns options
func (d *dns) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	return d.Session.Exec(r, out, in...)
}
