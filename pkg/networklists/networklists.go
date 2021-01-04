// Package networklist provides access to the Akamai Networklist APIs
package networklists

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// NETWORKLIST is the networklist api interface
	NETWORKLISTS interface {
		Activations
		NetworkList
		NetworkListDescription
		NetworkListSubscription
	}

	networklists struct {
		session.Session
		usePrefixes bool
	}

	// Option defines a networklist option
	Option func(*networklists)

	// ClientFunc is a networklist client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) NETWORKLISTS
)

// Client returns a new networklist Client instance with the specified controller
func Client(sess session.Session, opts ...Option) NETWORKLISTS {
	p := &networklists{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
