package ivm

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// IVM is the api interface for Image and Video Manager
	IVM interface {
		Policies
		PolicySets
	}

	ivm struct {
		session.Session
	}

	// Option defines an Image and Video Manager option
	Option func(*ivm)

	// ClientFunc is a Image and Video Manager client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) IVM
)

// Client returns a new ivmanager Client instance with the specified controller
func Client(sess session.Session, opts ...Option) IVM {
	c := &ivm{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
