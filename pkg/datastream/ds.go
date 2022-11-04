package datastream

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// DS is the ds api interface
	DS interface {
		Activation
		Properties
		Stream
	}

	ds struct {
		session.Session
	}

	// Option defines a DS option
	Option func(*ds)

	// ClientFunc is a ds client new method, this can be used for mocking
	ClientFunc func(sess session.Session, ops ...Option) DS
)

// Client returns a new ds Client instance with the specified controller
func Client(sess session.Session, opts ...Option) DS {
	c := &ds{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// DelimiterTypePtr returns the address of the DelimiterType
func DelimiterTypePtr(d DelimiterType) *DelimiterType {
	return &d
}
