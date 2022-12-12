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

func (m *Mock) GetEdgeHostname(ctx context.Context, id int) (*GetEdgeHostnameResponse, error) {
	args := m.Called(ctx, id)

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
