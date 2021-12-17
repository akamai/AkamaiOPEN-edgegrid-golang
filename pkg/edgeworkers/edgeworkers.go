package edgeworkers

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Edgeworkers is the api interface for edgeworkers
	Edgeworkers interface {
		Activations
		EdgeWorkerIDs
		PermissionGroups
		Properties
		ResourceTiers
	}

	edgeworkers struct {
		session.Session
	}

	// Option defines an Edgeworkers option
	Option func(*edgeworkers)

	// ClientFunc is a Edgeworkers client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) Edgeworkers
)

// Client returns a new edgeworkers Client instance with the specified controller
func Client(sess session.Session, opts ...Option) Edgeworkers {
	c := &edgeworkers{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
