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

func (m *Mock) ListActivePolicyProperties(ctx context.Context, req ListActivePolicyPropertiesRequest) (*ListActivePolicyPropertiesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListActivePolicyPropertiesResponse), args.Error(1)
}

func (m *Mock) ListPolicyActivations(ctx context.Context, req ListPolicyActivationsRequest) (*PolicyActivations, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyActivations), args.Error(1)
}

func (m *Mock) ActivatePolicy(ctx context.Context, req ActivatePolicyRequest) (*PolicyActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyActivation), args.Error(1)
}

func (m *Mock) DeactivatePolicy(ctx context.Context, req DeactivatePolicyRequest) (*PolicyActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyActivation), args.Error(1)
}

func (m *Mock) GetPolicyActivation(ctx context.Context, req GetPolicyActivationRequest) (*PolicyActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyActivation), args.Error(1)
}

func (m *Mock) ListPolicies(ctx context.Context, req ListPoliciesRequest) (*ListPoliciesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListPoliciesResponse), args.Error(1)
}

func (m *Mock) CreatePolicy(ctx context.Context, req CreatePolicyRequest) (*Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) DeletePolicy(ctx context.Context, req DeletePolicyRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *Mock) GetPolicy(ctx context.Context, req GetPolicyRequest) (*Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) UpdatePolicy(ctx context.Context, req UpdatePolicyRequest) (*Policy, error) {
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

func (m *Mock) ListCloudlets(ctx context.Context) ([]ListCloudletsItem, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ListCloudletsItem), args.Error(1)
}
