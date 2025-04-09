package mtlskeystore

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")
)

type (
	// MTLSKeystore is the interface for the mTLS Keystore API.
	MTLSKeystore interface{}

	mtlskeystore struct {
		session.Session
	}

	// Option is a function that configures the mTLS Keystore.
	Option func(*mtlskeystore)
)

// Client creates a new MTLSKeystore client.
func Client(sess session.Session, opts ...Option) MTLSKeystore {
	c := &mtlskeystore{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
