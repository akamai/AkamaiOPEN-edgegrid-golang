//revive:disable:exported

package v3

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ Cloudlets = &Mock{}

func (m *Mock) GetPolicyProperties(ctx context.Context, req GetPolicyPropertiesRequest) (*PolicyProperty, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyProperty), args.Error(1)
}
