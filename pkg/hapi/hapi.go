// Package hapi provides access to the Akamai Edge Hostnames APIs
//
// See: https://techdocs.akamai.com/edge-hostnames/reference/api
package hapi

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// HAPI is the hapi api interface
	HAPI interface {
		EdgeHostnames
	}

	hapi struct {
		session.Session
	}

	// Option defines a HAPI option
	Option func(*hapi)

	// ClientFunc is a hapi client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) HAPI
)

// Client returns a new hapi Client instance with the specified controller
func Client(sess session.Session, opts ...Option) HAPI {
	h := &hapi{
		Session: sess,
	}

	for _, opt := range opts {
		opt(h)
	}
	return h
}
