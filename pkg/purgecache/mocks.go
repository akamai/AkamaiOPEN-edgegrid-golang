//revive:disable:exported

package purgecache

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ PurgeCache = &Mock{}

func (m *Mock) RateLimitStatus(ctx context.Context, params RateLimitStatusRequest) (*RateLimitStatusResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RateLimitStatusResponse), args.Error(1)
}

func (m *Mock) DeleteByURL(ctx context.Context, params DeleteByURLRequest) (*DeleteResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteResponse), args.Error(1)
}

func (m *Mock) DeleteByTag(ctx context.Context, params DeleteByTagRequest) (*DeleteResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteResponse), args.Error(1)
}

func (m *Mock) DeleteByCPCode(ctx context.Context, params DeleteByCPCodeRequest) (*DeleteResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteResponse), args.Error(1)
}

func (m *Mock) InvalidateByURL(ctx context.Context, params InvalidateByURLRequest) (*InvalidateResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*InvalidateResponse), args.Error(1)
}

func (m *Mock) InvalidateByTag(ctx context.Context, params InvalidateByTagRequest) (*InvalidateResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*InvalidateResponse), args.Error(1)
}

func (m *Mock) InvalidateByCPCode(ctx context.Context, params InvalidateByCPCodeRequest) (*InvalidateResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*InvalidateResponse), args.Error(1)
}
