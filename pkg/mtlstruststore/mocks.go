//revive:disable:exported

package mtlstruststore

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ MTLSTruststore = &Mock{}

// Mock is MTLS Truststore API Mock
type Mock struct {
	mock.Mock
}

func (m *Mock) CreateCASet(ctx context.Context, params CreateCASetRequest) (*CreateCASetResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateCASetResponse), args.Error(1)
}

func (m *Mock) GetCASet(ctx context.Context, params GetCASetRequest) (*GetCASetResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCASetResponse), args.Error(1)
}

func (m *Mock) ListCASets(ctx context.Context, params ListCASetsRequest) (*ListCASetsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListCASetsResponse), args.Error(1)
}

func (m *Mock) DeleteCASet(ctx context.Context, params DeleteCASetRequest) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *Mock) CreateCASetVersion(ctx context.Context, params CreateCASetVersionRequest) (*CreateCASetVersionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateCASetVersionResponse), args.Error(1)
}

func (m *Mock) CloneCASetVersion(ctx context.Context, params CloneCASetVersionRequest) (*CloneCASetVersionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CloneCASetVersionResponse), args.Error(1)
}

func (m *Mock) UpdateCASetVersion(ctx context.Context, params UpdateCASetVersionRequest) (*UpdateCASetVersionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateCASetVersionResponse), args.Error(1)
}

func (m *Mock) GetCASetVersion(ctx context.Context, params GetCASetVersionRequest) (*GetCASetVersionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCASetVersionResponse), args.Error(1)
}

func (m *Mock) ListCASetVersions(ctx context.Context, params ListCASetVersionsRequest) (*ListCASetVersionsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListCASetVersionsResponse), args.Error(1)
}

func (m *Mock) GetCASetVersionCertificates(ctx context.Context, params GetCASetVersionCertificatesRequest) (*GetCASetVersionCertificatesResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCASetVersionCertificatesResponse), args.Error(1)
}

func (m *Mock) ActivateCASetVersion(ctx context.Context, params ActivateCASetVersionRequest) (*ActivateCASetVersionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ActivateCASetVersionResponse), args.Error(1)
}

func (m *Mock) DeactivateCASetVersion(ctx context.Context, params DeactivateCASetVersionRequest) (*DeactivateCASetVersionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeactivateCASetVersionResponse), args.Error(1)
}

func (m *Mock) GetCASetVersionActivation(ctx context.Context, params GetCASetVersionActivationRequest) (*GetCASetVersionActivationResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCASetVersionActivationResponse), args.Error(1)
}

func (m *Mock) ListCASetActivations(ctx context.Context, params ListCASetActivationsRequest) (*ListCASetActivationsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListCASetActivationsResponse), args.Error(1)
}

func (m *Mock) ListCASetVersionActivations(ctx context.Context, params ListCASetVersionActivationsRequest) (*ListCASetVersionActivationsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListCASetVersionActivationsResponse), args.Error(1)
}

func (m *Mock) ListCASetAssociations(ctx context.Context, params ListCASetAssociationsRequest) (*ListCASetAssociationsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListCASetAssociationsResponse), args.Error(1)
}

func (m *Mock) CloneCASet(ctx context.Context, params CloneCASetRequest) (*CloneCASetResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CloneCASetResponse), args.Error(1)
}

func (m *Mock) GetCASetDeletionStatus(ctx context.Context, params GetCASetDeletionStatusRequest) (*GetCASetDeletionStatusResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCASetDeletionStatusResponse), args.Error(1)
}

func (m *Mock) ListCASetActivities(ctx context.Context, params ListCASetActivitiesRequest) (*ListCASetActivitiesResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListCASetActivitiesResponse), args.Error(1)
}

func (m *Mock) ValidateCertificates(ctx context.Context, params ValidateCertificatesRequest) (*ValidateCertificatesResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ValidateCertificatesResponse), args.Error(1)
}

func (m *Mock) DeleteCASetVersion(ctx context.Context, params DeleteCASetVersionRequest) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}
