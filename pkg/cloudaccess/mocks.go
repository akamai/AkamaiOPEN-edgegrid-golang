//revive:disable:exported

package cloudaccess

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ CloudAccess = &Mock{}

func (m *Mock) GetAccessKeyStatus(ctx context.Context, r GetAccessKeyStatusRequest) (*GetAccessKeyStatusResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetAccessKeyStatusResponse), args.Error(1)
}

func (m *Mock) GetAccessKeyVersionStatus(ctx context.Context, r GetAccessKeyVersionStatusRequest) (*GetAccessKeyVersionStatusResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetAccessKeyVersionStatusResponse), args.Error(1)
}
