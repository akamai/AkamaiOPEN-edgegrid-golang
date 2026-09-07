// Package cloudcertificates provides access to the Akamai Cloud Certificate Manager API.
package cloudcertificates

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

// CloudCertificates defines the interface for Akamai Cloud Certificate Manager operations.
type (
	CloudCertificates interface {
		// CreateCertificate creates a third party certificate.
		//
		// See: https://techdocs.akamai.com/ccm/reference/post-certificates
		CreateCertificate(ctx context.Context, params CreateCertificateRequest) (*CreateCertificateResponse, error)

		// GetCertificate retrieves a single certificate by its certificateID.
		// This includes details such as the certificate name, SANs, subject, certificate signing request, and, if uploaded, the signed certificate and trust chain.
		//
		// See: https://techdocs.akamai.com/ccm/reference/get-cert
		GetCertificate(ctx context.Context, params GetCertificateRequest) (*GetCertificateResponse, error)

		// PatchCertificate patches some fields of the certificate: name, signed certificate, and trust chain.
		// It allows to rename or reset the name of certificate or upload of signed certificate and optional trust chain or both in one request.
		//
		// See: https://techdocs.akamai.com/ccm/reference/patch-certificate
		PatchCertificate(ctx context.Context, params PatchCertificateRequest) (*PatchCertificateResponse, error)

		// UpdateCertificate updates a certificate. This includes renaming the certificate, uploading a signed certificate, and an optional trust chain.
		//
		// See: https://techdocs.akamai.com/ccm/reference/put-cert
		UpdateCertificate(ctx context.Context, params UpdateCertificateRequest) (*UpdateCertificateResponse, error)

		// DeleteCertificate deletes a certificate by its certificateID. Note that only certificates that are not ACTIVE can be deleted.
		//
		// See: https://techdocs.akamai.com/ccm/reference/delete-certificate
		DeleteCertificate(ctx context.Context, params DeleteCertificateRequest) (*DeleteCertificateResponse, error)

		// ListCertificateBindings provides hostname bindings for the given certificate.
		//
		// See: https://techdocs.akamai.com/ccm/reference/get-single-cert-bindings
		ListCertificateBindings(ctx context.Context, params ListCertificateBindingsRequest) (*ListCertificateBindingsResponse, error)

		// ListCertificates lists all certificates that are accessible for the requesting end user.
		//
		// See: https://techdocs.akamai.com/ccm/reference/get-certificates
		ListCertificates(ctx context.Context, params ListCertificatesRequest) (*ListCertificatesResponse, error)

		// ListBindings provides hostname bindings for user accessible certificates, optionally filtered by contract, group, domain, or expiration days.
		//
		// See: https://techdocs.akamai.com/ccm/reference/get-all-cert-bindings
		ListBindings(ctx context.Context, params ListBindingsRequest) (*ListBindingsResponse, error)
	}

	cloudcertificates struct {
		session.Session
	}

	// Option is a function that configures the CloudCertificates.
	Option func(*cloudcertificates)
)

// Client creates a new CloudCertificates client.
func Client(sess session.Session, opts ...Option) CloudCertificates {
	c := &cloudcertificates{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
