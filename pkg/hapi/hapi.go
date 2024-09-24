// Package hapi provides access to the Akamai Edge Hostnames APIs
//
// See: https://techdocs.akamai.com/edge-hostnames/reference/api
package hapi

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// HAPI is the hapi api interface
	HAPI interface {
		// ChangeRequests

		// GetChangeRequest request status and details specified by the change ID
		// that is provided when you make a change request.
		//
		// See: https://techdocs.akamai.com/edge-hostnames/reference/get-changeid
		GetChangeRequest(context.Context, GetChangeRequest) (*ChangeRequest, error)

		// EdgeHostnames

		// DeleteEdgeHostname allows deleting a specific edge hostname.
		// You must have an Admin or Technical role in order to delete an edge hostname.
		// You can delete any hostname that’s not currently part of an active Property Manager configuration.
		//
		// See: https://techdocs.akamai.com/edge-hostnames/reference/delete-edgehostname
		DeleteEdgeHostname(context.Context, DeleteEdgeHostnameRequest) (*DeleteEdgeHostnameResponse, error)

		// GetEdgeHostname gets a specific edge hostname's details including its product ID, IP version behavior,
		// and China CDN or Edge IP Binding status.
		//
		// See: https://techdocs.akamai.com/edge-hostnames/reference/get-edgehostnameid
		GetEdgeHostname(context.Context, int) (*GetEdgeHostnameResponse, error)

		// UpdateEdgeHostname allows update ttl (path = "/ttl") or IpVersionBehaviour (path = "/ipVersionBehavior")
		//
		// See: https://techdocs.akamai.com/edge-hostnames/reference/patch-edgehostnames
		UpdateEdgeHostname(context.Context, UpdateEdgeHostnameRequest) (*UpdateEdgeHostnameResponse, error)

		// GetCertificate gets the certificate associated with an enhanced TLS edge hostname
		//
		// See: https://techdocs.akamai.com/edge-hostnames/reference/get-edge-hostname-certificate
		GetCertificate(context.Context, GetCertificateRequest) (*GetCertificateResponse, error)
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
