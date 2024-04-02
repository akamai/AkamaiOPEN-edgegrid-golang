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

func (m *Mock) CreateAccessKey(ctx context.Context, r CreateAccessKeyRequest) (*CreateAccessKeyResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateAccessKeyResponse), args.Error(1)
}

func (m *Mock) GetAccessKey(ctx context.Context, r AccessKeyRequest) (*GetAccessKeyResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetAccessKeyResponse), args.Error(1)
}

func (m *Mock) ListAccessKeys(ctx context.Context, r ListAccessKeysRequest) (*ListAccessKeysResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListAccessKeysResponse), args.Error(1)
}

func (m *Mock) UpdateAccessKey(ctx context.Context, request UpdateAccessKeyRequest, param AccessKeyRequest) (*UpdateAccessKeyResponse, error) {
	args := m.Called(ctx, request, param)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateAccessKeyResponse), args.Error(1)
}

func (m *Mock) DeleteAccessKey(ctx context.Context, r AccessKeyRequest) error {
	args := m.Called(ctx, r)

	return args.Error(0)
}
