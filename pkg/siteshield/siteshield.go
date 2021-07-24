// Package siteshield provides access to the Akamai Site Shield APIs
package siteshield

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// SSMAPS is the siteshieldmap api interface
	SSMAPS interface {
		SiteShieldMap
	}

	siteshieldmap struct {
		session.Session
		usePrefixes bool
	}

	// Option defines a siteshieldmap option
	Option func(*siteshieldmap)

	// ClientFunc is a siteshieldmap client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) SSMAPS
)

// Client returns a new siteshieldmap Client instance with the specified controller
func Client(sess session.Session, opts ...Option) SSMAPS {
	s := &siteshieldmap{
		Session: sess,
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}
