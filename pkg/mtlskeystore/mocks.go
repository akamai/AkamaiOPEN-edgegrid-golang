//revive:disable:exported

package mtlskeystore

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Mock is MTLS Keystore API Mock.
type Mock struct {
	mock.Mock
}

var _ MTLSKeystore = &Mock{}

func (m *Mock) ListClientCertificates(ctx context.Context) (*ListClientCertificatesResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListClientCertificatesResponse), args.Error(1)
}

func (m *Mock) GetClientCertificate(ctx context.Context, params GetClientCertificateRequest) (*GetClientCertificateResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetClientCertificateResponse), args.Error(1)
}

func (m *Mock) CreateClientCertificate(ctx context.Context, params CreateClientCertificateRequest) (*CreateClientCertificateResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateClientCertificateResponse), args.Error(1)
}

func (m *Mock) PatchClientCertificate(ctx context.Context, params PatchClientCertificateRequest) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *Mock) ListClientCertificateVersions(ctx context.Context, params ListClientCertificateVersionsRequest) (*ListClientCertificateVersionsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListClientCertificateVersionsResponse), args.Error(1)
}

func (m *Mock) RotateClientCertificateVersion(ctx context.Context, params RotateClientCertificateVersionRequest) (*RotateClientCertificateVersionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RotateClientCertificateVersionResponse), args.Error(1)
}

func (m *Mock) DeleteClientCertificateVersion(ctx context.Context, params DeleteClientCertificateVersionRequest) (*DeleteClientCertificateVersionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteClientCertificateVersionResponse), args.Error(1)
}

func (m *Mock) UploadSignedClientCertificate(ctx context.Context, params UploadSignedClientCertificateRequest) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *Mock) ListAccountCACertificates(ctx context.Context, params ListAccountCACertificatesRequest) (*ListAccountCACertificatesResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListAccountCACertificatesResponse), args.Error(1)
}
