// Package cloudaccess provides access to the Akamai Cloud Access Manager API
package cloudaccess

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// CloudAccess is the api interface
	CloudAccess interface {
	}

	cloudaccess struct {
		session.Session
	}

	// Option defines an CloudAccess option
	Option func(*cloudaccess)
)

// Client returns a new cloudaccess Client instance with the specified controller
func Client(sess session.Session, opts ...Option) CloudAccess {
	c := &cloudaccess{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
