//revive:disable:exported

package hapi

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ HAPI = &Mock{}

func (m *Mock) DeleteEdgeHostname(ctx context.Context, request DeleteEdgeHostnameRequest) (*DeleteEdgeHostnameResponse, error) {
	args := m.Called(ctx, request)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteEdgeHostnameResponse), nil
}

func (m *Mock) GetEdgeHostname(ctx context.Context, request GetEdgeHostnameRequest) (*GetEdgeHostnameResponse, error) {
	args := m.Called(ctx, request)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetEdgeHostnameResponse), nil
}

func (m *Mock) UpdateEdgeHostname(ctx context.Context, request UpdateEdgeHostnameRequest) (*UpdateEdgeHostnameResponse, error) {
	args := m.Called(ctx, request)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateEdgeHostnameResponse), nil
}

func (m *Mock) GetChangeRequest(ctx context.Context, req GetChangeRequestRequest) (*GetChangeRequestResponse, error) {
	args := m.Called(ctx, req)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetChangeRequestResponse), nil
}

func (m *Mock) ListChangeRequests(ctx context.Context, req ListChangeRequestsRequest) (*ListChangeRequestsResponse, error) {
	args := m.Called(ctx, req)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListChangeRequestsResponse), nil
}

func (m *Mock) ListChangeRequestsPerEdgeHostname(ctx context.Context, req ListChangeRequestsPerEdgeHostnameRequest) (*ListChangeRequestsPerEdgeHostnameResponse, error) {
	args := m.Called(ctx, req)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListChangeRequestsPerEdgeHostnameResponse), nil
}

func (m *Mock) GetCertificate(ctx context.Context, req GetCertificateRequest) (*GetCertificateResponse, error) {
	args := m.Called(ctx, req)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCertificateResponse), nil
}
