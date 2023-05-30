// Package cloudwrapper provides access to the Akamai Cloud Wrapper API
package cloudwrapper

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// CloudWrapper is the api interface for Cloud Wrapper
	CloudWrapper interface {
		Capacities
		Configurations
		Locations
		Properties
	}

	cloudwrapper struct {
		session.Session
	}

	// Option defines an CloudWrapper option
	Option func(*cloudwrapper)
)

// Client returns a new cloudwrapper Client instance with the specified controller
func Client(sess session.Session, opts ...Option) CloudWrapper {
	c := &cloudwrapper{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
