package edgeworkers

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Edgeworkers is the api interface for EdgeWorkers and EdgeKV
	Edgeworkers interface {
		Activations
		Contracts
		Deactivations
		EdgeKVAccessTokens
		EdgeKVInitialize
		EdgeKVItems
		EdgeKVNamespaces
		EdgeWorkerIDs
		EdgeWorkerVersions
		PermissionGroups
		Properties
		Reports
		ResourceTiers
		SecureTokens
		Validations
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
	e := &edgeworkers{
		Session: sess,
	}

	for _, opt := range opts {
		opt(e)
	}
	return e
}
