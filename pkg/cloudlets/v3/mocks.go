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

func (m *Mock) ListActivePolicyProperties(ctx context.Context, req ListActivePolicyPropertiesRequest) (*PolicyProperty, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyProperty), args.Error(1)
}

func (m *Mock) ListSharedPolicies(ctx context.Context, req ListSharedPoliciesRequest) (*ListSharedPoliciesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListSharedPoliciesResponse), args.Error(1)
}

func (m *Mock) CreateSharedPolicy(ctx context.Context, req CreateSharedPolicyRequest) (*Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) DeleteSharedPolicy(ctx context.Context, req DeleteSharedPolicyRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *Mock) GetSharedPolicy(ctx context.Context, req GetSharedPolicyRequest) (*Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) UpdateSharedPolicy(ctx context.Context, req UpdateSharedPolicyRequest) (*Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) ClonePolicy(ctx context.Context, req ClonePolicyRequest) (*Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) CreatePolicyVersion(ctx context.Context, req CreatePolicyVersionRequest) (*PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyVersion), args.Error(1)
}

func (m *Mock) DeletePolicyVersion(ctx context.Context, req DeletePolicyVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *Mock) GetPolicyVersion(ctx context.Context, req GetPolicyVersionRequest) (*PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyVersion), args.Error(1)
}

func (m *Mock) ListPolicyVersions(ctx context.Context, req ListPolicyVersionsRequest) (*ListPolicyVersions, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListPolicyVersions), args.Error(1)
}

func (m *Mock) UpdatePolicyVersion(ctx context.Context, req UpdatePolicyVersionRequest) (*PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyVersion), args.Error(1)
}
