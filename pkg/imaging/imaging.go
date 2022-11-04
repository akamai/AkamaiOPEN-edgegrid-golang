package imaging

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Imaging is the api interface for Image and Video Manager
	Imaging interface {
		Policies
		PolicySets
	}

	imaging struct {
		session.Session
	}

	// Option defines an Image and Video Manager option
	Option func(*imaging)

	// ClientFunc is a Image and Video Manager client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) Imaging
)

// Client returns a new Image and Video Manager Client instance with the specified controller
func Client(sess session.Session, opts ...Option) Imaging {
	c := &imaging{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
