//revive:disable:exported

package cloudcertificates

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ CloudCertificates = &Mock{}

func (m *Mock) CreateLineage(ctx context.Context, req CreateLineageRequest) (*CreateLineageResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateLineageResponse), args.Error(1)
}

func (m *Mock) GetLineage(ctx context.Context, req GetLineageRequest) (*GetLineageResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetLineageResponse), args.Error(1)
}

func (m *Mock) ListLineages(ctx context.Context, req ListLineagesRequest) (*ListLineagesResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListLineagesResponse), args.Error(1)
}

func (m *Mock) RenameLineage(ctx context.Context, req RenameLineageRequest) (*RenameLineageResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RenameLineageResponse), args.Error(1)
}

func (m *Mock) DeleteLineage(ctx context.Context, req DeleteLineageRequest) error {
	args := m.Called(ctx, req)

	return args.Error(0)
}

func (m *Mock) RenewLineage(ctx context.Context, req RenewLineageRequest) (*RenewLineageResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RenewLineageResponse), args.Error(1)
}

func (m *Mock) CompleteLineage(ctx context.Context, req CompleteLineageRequest) (*CompleteLineageResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CompleteLineageResponse), args.Error(1)
}

func (m *Mock) ListLineageActivity(ctx context.Context, req ListLineageActivityRequest) (*ListLineageActivityResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListLineageActivityResponse), args.Error(1)
}

func (m *Mock) GetGeneration(ctx context.Context, req GetGenerationRequest) (*GetGenerationResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetGenerationResponse), args.Error(1)
}

func (m *Mock) DeleteGeneration(ctx context.Context, req DeleteGenerationRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *Mock) ListArchivedGenerations(ctx context.Context, req ListArchivedGenerationsRequest) (*ListArchivedGenerationsResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListArchivedGenerationsResponse), args.Error(1)
}

func (m *Mock) UploadSignedCertificate(ctx context.Context, req UploadSignedCertificateRequest) (*UploadSignedCertificateResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UploadSignedCertificateResponse), args.Error(1)
}

func (m *Mock) PromoteLineage(ctx context.Context, req PromoteLineageRequest) (*PromoteLineageResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*PromoteLineageResponse), args.Error(1)
}

func (m *Mock) RollbackLineage(ctx context.Context, req RollbackLineageRequest) (*RollbackLineageResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RollbackLineageResponse), args.Error(1)
}

func (m *Mock) ReplaceStagingLineage(ctx context.Context, req ReplaceStagingLineageRequest) (*ReplaceStagingLineageResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ReplaceStagingLineageResponse), args.Error(1)
}

func (m *Mock) GetActivationStatus(ctx context.Context, req GetActivationStatusRequest) (*GetActivationStatusResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetActivationStatusResponse), args.Error(1)
}

func (m *Mock) ListActivations(ctx context.Context, req ListActivationsRequest) (*ListActivationsResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListActivationsResponse), args.Error(1)
}

func (m *Mock) ListLineageBindings(ctx context.Context, req ListLineageBindingsRequest) (*ListLineageBindingsResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListLineageBindingsResponse), args.Error(1)
}
