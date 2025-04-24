package mtlskeystore

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")
)

type (
	// MTLSKeystore is the interface for the mTLS Keystore API.
	MTLSKeystore interface {

		// Client certificates

		// ListClientCertificates lists client certificates under the account.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/get-client-certs
		ListClientCertificates(ctx context.Context) (*ListClientCertificatesResponse, error)

		// GetClientCertificate gets details of a client certificate.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/get-client-cert
		GetClientCertificate(ctx context.Context, params GetClientCertificateRequest) (*GetClientCertificateResponse, error)

		// CreateClientCertificate creates a client certificate with the provided name.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/post-client-cert
		CreateClientCertificate(ctx context.Context, params CreateClientCertificateRequest) (*CreateClientCertificateResponse, error)

		// PatchClientCertificate updates the client certificate's name or notification emails.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/patch-client-cert
		PatchClientCertificate(ctx context.Context, params PatchClientCertificateRequest) error
	}

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
