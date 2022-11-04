// Package iam provides access to the Akamai Property APIs
package iam

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// IAM is the IAM api interface
	IAM interface {
		BlockedProperties
		Groups
		Roles
		Support
		UserLock
		UserPassword
		Users
	}

	iam struct {
		session.Session
	}

	// Option defines a IAM option
	Option func(*iam)

	// ClientFunc is an IAM client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) IAM
)

// Client returns a new IAM Client instance with the specified controller
func Client(sess session.Session, opts ...Option) IAM {
	p := &iam{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
