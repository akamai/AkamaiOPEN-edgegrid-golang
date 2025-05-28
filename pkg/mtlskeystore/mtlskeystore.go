package mtlskeystore

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
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

		// RotateClientCertificateVersion creates a new version in the client certificate.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/post-client-cert-version
		RotateClientCertificateVersion(ctx context.Context, params RotateClientCertificateVersionRequest) (*RotateClientCertificateVersionResponse, error)

		// GetClientCertificateVersions lists versions of the client certificate specified by certificateID.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/get-client-cert-versions
		GetClientCertificateVersions(ctx context.Context, params GetClientCertificateVersionsRequest) (*GetClientCertificateVersionsResponse, error)

		// DeleteClientCertificateVersion deletes a client certificate version with the provided certificateID and version.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/delete-client-certificate
		DeleteClientCertificateVersion(ctx context.Context, params DeleteClientCertificateVersionRequest) (*DeleteClientCertificateVersionResponse, error)

		// UploadSignedClientCertificate uploads a signed THIRD_PARTY client certificate.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/post-cert-block
		UploadSignedClientCertificate(ctx context.Context, params UploadSignedClientCertificateRequest) error

		// ListAccountCACertificates lists CA certificates under the account.
		//
		// See: https://techdocs.akamai.com/mtls-origin-keystore/reference/get-ca-certs
		ListAccountCACertificates(ctx context.Context, params ListAccountCACertificatesRequest) (*ListAccountCACertificatesResponse, error)
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
