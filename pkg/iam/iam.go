// Package iam provides access to the Akamai Property APIs
package iam

import (
	"errors"
	"path"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
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

var (
	// BaseEndPoint is the IAM basepath
	BaseEndPoint = "/identity-management/v2"

	// UserAdminEP is the IAM user-admin endpoint
	UserAdminEP = path.Join(BaseEndPoint, "user-admin")
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
