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

	if args.Get(0) == nil {
		return args.Error(1)
	}

	return args.Error(1)
}
