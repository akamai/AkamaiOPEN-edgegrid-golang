// Package iam provides access to the Akamai Property APIs
package iam

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// IAM is the IAM api interface
	IAM interface {
		APIClients
		APIClientsCredentials
		BlockedProperties
		CIDR
		Groups
		Helper
		IPAllowlist
		Properties
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

func convertStructToReqBody(srcStruct interface{}) (io.Reader, error) {
	reqBody, err := json.Marshal(srcStruct)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(reqBody), nil
}
