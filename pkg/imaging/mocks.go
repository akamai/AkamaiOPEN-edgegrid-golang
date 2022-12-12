//revive:disable:exported

package imaging

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ Imaging = &Mock{}

func (m *Mock) GetPolicy(ctx context.Context, req GetPolicyRequest) (PolicyOutput, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(PolicyOutput), args.Error(1)
}

func (m *Mock) UpsertPolicy(ctx context.Context, req UpsertPolicyRequest) (*PolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyResponse), args.Error(1)
}

func (m *Mock) DeletePolicy(ctx context.Context, req DeletePolicyRequest) (*PolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyResponse), args.Error(1)
}

func (m *Mock) RollbackPolicy(ctx context.Context, req RollbackPolicyRequest) (*PolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyResponse), args.Error(1)
}

func (m *Mock) GetPolicyHistory(ctx context.Context, req GetPolicyHistoryRequest) (*GetPolicyHistoryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPolicyHistoryResponse), args.Error(1)
}

func (m *Mock) ListPolicies(ctx context.Context, req ListPoliciesRequest) (*ListPoliciesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListPoliciesResponse), args.Error(1)
}

func (m *Mock) ListPolicySets(ctx context.Context, req ListPolicySetsRequest) ([]PolicySet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]PolicySet), args.Error(1)
}

func (m *Mock) GetPolicySet(ctx context.Context, req GetPolicySetRequest) (*PolicySet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicySet), args.Error(1)
}

func (m *Mock) CreatePolicySet(ctx context.Context, req CreatePolicySetRequest) (*PolicySet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicySet), args.Error(1)
}

func (m *Mock) UpdatePolicySet(ctx context.Context, req UpdatePolicySetRequest) (*PolicySet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicySet), args.Error(1)
}

func (m *Mock) DeletePolicySet(ctx context.Context, req DeletePolicySetRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
