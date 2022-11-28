// Package networklists provides access to the Akamai Networklist APIs
package networklists

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// NTWRKLISTS is the networklist api interface
	NTWRKLISTS interface {
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
	ClientFunc func(sess session.Session, opts ...Option) NTWRKLISTS
)

// Client returns a new networklist Client instance with the specified controller
func Client(sess session.Session, opts ...Option) NTWRKLISTS {
	p := &networklists{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
