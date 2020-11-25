// Package iam provides access to the Akamai Property APIs
package iam

import (
	"errors"
	"path"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")

	// ErrNotFound is returned when requested resource was not found
	ErrNotFound = errors.New("resource not found")
)

type (
	// IAM is the iam api interface
	IAM interface {
		Groups
		Roles
		Support
	}

	iam struct {
		session.Session
	}

	// Option defines a IAM option
	Option func(*iam)

	// ClientFunc is a iam client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) IAM

	// Response is a base IAM Response type
	Response struct {
		AccountID  string   `json:"accountId,omitempty"`
		ContractID string   `json:"contractId,omitempty"`
		GroupID    string   `json:"groupId,omitempty"`
		Etag       string   `json:"etag,omitempty"`
		Errors     []*Error `json:"errors,omitempty"`
		Warnings   []*Error `json:"warnings,omitempty"`
	}
)

var (
	// BaseEndPoint is the IAM basepath
	BaseEndPoint = "/identity-management/v2"

	// UserAdminEP is the IAM user-admin endpoint
	UserAdminEP = path.Join(BaseEndPoint, "user-admin")
)

// Client returns a new iam Client instance with the specified controller
func Client(sess session.Session, opts ...Option) IAM {
	p := &iam{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
