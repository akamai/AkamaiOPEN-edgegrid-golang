//revive:disable:exported

package cloudwrapper

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ CloudWrapper = &Mock{}

// ListCapacities implements CloudWrapper
func (m *Mock) ListCapacities(ctx context.Context, req ListCapacitiesRequest) (*ListCapacitiesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListCapacitiesResponse), args.Error(1)
}

// ListLocations implements CloudWrapper
func (m *Mock) ListLocations(ctx context.Context) (*ListLocationResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListLocationResponse), args.Error(1)
}
