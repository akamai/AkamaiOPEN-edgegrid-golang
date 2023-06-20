// Package clientlists provides access to Akamai Client Lists APIs
//
// See: https://techdocs.akamai.com/client-lists/reference/api
package clientlists

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// ClientLists is the clientlists api interface
	ClientLists interface {
		Lists
	}

	clientlists struct {
		session.Session
	}

	// Option defines a clientlists option
	Option func(*clientlists)

	// ClientFunc is a clientlists client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) ClientLists
)

// Client returns a new clientlists Client instance with the specified controller
func Client(sess session.Session, opts ...Option) ClientLists {
	p := &clientlists{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
