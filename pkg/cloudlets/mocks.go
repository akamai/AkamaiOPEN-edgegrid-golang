//revive:disable:exported

package cloudlets

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ Cloudlets = &Mock{}

func (m *Mock) DeletePolicyProperty(ctx context.Context, req DeletePolicyPropertyRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *Mock) CreateLoadBalancerVersion(ctx context.Context, req CreateLoadBalancerVersionRequest) (*LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoadBalancerVersion), args.Error(1)
}

func (m *Mock) GetLoadBalancerVersion(ctx context.Context, req GetLoadBalancerVersionRequest) (*LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoadBalancerVersion), args.Error(1)
}

func (m *Mock) UpdateLoadBalancerVersion(ctx context.Context, req UpdateLoadBalancerVersionRequest) (*LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoadBalancerVersion), args.Error(1)
}

func (m *Mock) ListLoadBalancerActivations(ctx context.Context, req ListLoadBalancerActivationsRequest) ([]LoadBalancerActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LoadBalancerActivation), args.Error(1)
}

func (m *Mock) ActivateLoadBalancerVersion(ctx context.Context, req ActivateLoadBalancerVersionRequest) (*LoadBalancerActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoadBalancerActivation), args.Error(1)
}

func (m *Mock) ListPolicyActivations(ctx context.Context, req ListPolicyActivationsRequest) ([]PolicyActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]PolicyActivation), args.Error(1)
}

func (m *Mock) ActivatePolicyVersion(ctx context.Context, req ActivatePolicyVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *Mock) ListOrigins(ctx context.Context, req ListOriginsRequest) ([]OriginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]OriginResponse), args.Error(1)
}

func (m *Mock) GetOrigin(ctx context.Context, req GetOriginRequest) (*Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Origin), args.Error(1)
}

func (m *Mock) CreateOrigin(ctx context.Context, req CreateOriginRequest) (*Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Origin), args.Error(1)
}

func (m *Mock) UpdateOrigin(ctx context.Context, req UpdateOriginRequest) (*Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Origin), args.Error(1)
}

func (m *Mock) ListPolicies(ctx context.Context, request ListPoliciesRequest) ([]Policy, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Policy), args.Error(1)
}

func (m *Mock) GetPolicy(ctx context.Context, policyID GetPolicyRequest) (*Policy, error) {
	args := m.Called(ctx, policyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) CreatePolicy(ctx context.Context, req CreatePolicyRequest) (*Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) RemovePolicy(ctx context.Context, policyID RemovePolicyRequest) error {
	args := m.Called(ctx, policyID)
	return args.Error(0)
}

func (m *Mock) UpdatePolicy(ctx context.Context, req UpdatePolicyRequest) (*Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Policy), args.Error(1)
}

func (m *Mock) ListPolicyVersions(ctx context.Context, request ListPolicyVersionsRequest) ([]PolicyVersion, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]PolicyVersion), args.Error(1)
}

func (m *Mock) GetPolicyVersion(ctx context.Context, req GetPolicyVersionRequest) (*PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyVersion), args.Error(1)
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

func (m *Mock) UpdatePolicyVersion(ctx context.Context, req UpdatePolicyVersionRequest) (*PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyVersion), args.Error(1)
}

func (m *Mock) GetPolicyProperties(ctx context.Context, req GetPolicyPropertiesRequest) (map[string]PolicyProperty, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]PolicyProperty), args.Error(1)
}

func (m *Mock) ListLoadBalancerVersions(ctx context.Context, req ListLoadBalancerVersionsRequest) ([]LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LoadBalancerVersion), args.Error(1)
}
