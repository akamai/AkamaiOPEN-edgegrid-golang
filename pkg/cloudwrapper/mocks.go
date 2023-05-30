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

func (m *Mock) ListCapacities(ctx context.Context, req ListCapacitiesRequest) (*ListCapacitiesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListCapacitiesResponse), args.Error(1)
}

func (m *Mock) ListLocations(ctx context.Context) (*ListLocationResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListLocationResponse), args.Error(1)
}

func (m *Mock) ListProperties(ctx context.Context, r ListPropertiesRequest) (*ListPropertiesResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListPropertiesResponse), args.Error(1)
}

func (m *Mock) ListOrigins(ctx context.Context, r ListOriginsRequest) (*ListOriginsResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListOriginsResponse), args.Error(1)
}

func (m *Mock) GetConfiguration(ctx context.Context, r GetConfigurationRequest) (*Configuration, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Configuration), args.Error(1)
}

func (m *Mock) ListConfigurations(ctx context.Context) (*ListConfigurationsResponse, error) {
	args := m.Called(ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListConfigurationsResponse), args.Error(1)
}

func (m *Mock) CreateConfiguration(ctx context.Context, r CreateConfigurationRequest) (*Configuration, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Configuration), args.Error(1)
}

func (m *Mock) UpdateConfiguration(ctx context.Context, r UpdateConfigurationRequest) (*Configuration, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Configuration), args.Error(1)
}

func (m *Mock) ActivateConfiguration(ctx context.Context, r ActivateConfigurationRequest) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
