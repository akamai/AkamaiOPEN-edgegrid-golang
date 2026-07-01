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
