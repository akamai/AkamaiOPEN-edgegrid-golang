//revive:disable:exported

package datastream

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ Stream = &Mock{}

func (m *Mock) CreateStream(ctx context.Context, r CreateStreamRequest) (*StreamUpdate, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*StreamUpdate), args.Error(1)
}

func (m *Mock) GetStream(ctx context.Context, r GetStreamRequest) (*DetailedStreamVersion, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DetailedStreamVersion), args.Error(1)
}

func (m *Mock) UpdateStream(ctx context.Context, r UpdateStreamRequest) (*StreamUpdate, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*StreamUpdate), args.Error(1)
}

func (m *Mock) DeleteStream(ctx context.Context, r DeleteStreamRequest) (*DeleteStreamResponse, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteStreamResponse), args.Error(1)
}

func (m *Mock) ListStreams(ctx context.Context, r ListStreamsRequest) ([]StreamDetails, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]StreamDetails), args.Error(1)
}

func (m *Mock) ActivateStream(ctx context.Context, r ActivateStreamRequest) (*ActivateStreamResponse, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ActivateStreamResponse), args.Error(1)
}

func (m *Mock) DeactivateStream(ctx context.Context, r DeactivateStreamRequest) (*DeactivateStreamResponse, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeactivateStreamResponse), args.Error(1)
}

func (m *Mock) GetActivationHistory(ctx context.Context, r GetActivationHistoryRequest) ([]ActivationHistoryEntry, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]ActivationHistoryEntry), args.Error(1)
}

func (m *Mock) GetProperties(ctx context.Context, r GetPropertiesRequest) ([]Property, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Property), args.Error(1)
}

func (m *Mock) GetPropertiesByGroup(ctx context.Context, r GetPropertiesByGroupRequest) ([]Property, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Property), args.Error(1)
}

func (m *Mock) GetDatasetFields(ctx context.Context, r GetDatasetFieldsRequest) ([]DataSets, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]DataSets), args.Error(1)
}
