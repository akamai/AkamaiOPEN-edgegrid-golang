package cloudlets

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Cloudlets is the api interface for cloudlets
	Cloudlets interface {
		LoadBalancers
		LoadBalancerVersions
		LoadBalancerActivations
		Policies
		PolicyProperties
		PolicyVersions
		PolicyVersionActivations
	}

	cloudlets struct {
		session.Session
	}

	// Option defines a Cloudlets option
	Option func(*cloudlets)

	// ClientFunc is a Cloudlets client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) Cloudlets
)

// Client returns a new cloudlets Client instance with the specified controller
func Client(sess session.Session, opts ...Option) Cloudlets {
	c := &cloudlets{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
