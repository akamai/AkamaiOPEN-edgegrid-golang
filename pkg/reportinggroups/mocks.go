//revive:disable:exported

package reportinggroups

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ ReportingGroups = &Mock{}

func (m *Mock) CreateReportingGroup(ctx context.Context, req CreateReportingGroupRequest) (*CreateReportingGroupResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateReportingGroupResponse), args.Error(1)
}

func (m *Mock) GetReportingGroup(ctx context.Context, req GetReportingGroupsRequest) (*GetReportingGroupResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetReportingGroupResponse), args.Error(1)
}

func (m *Mock) UpdateReportingGroup(ctx context.Context, req UpdateReportingGroupRequest) (*UpdateReportingGroupResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateReportingGroupResponse), args.Error(1)
}

func (m *Mock) DeleteReportingGroup(ctx context.Context, req DeleteReportingGroupRequest) error {
	args := m.Called(ctx, req)

	return args.Error(0)
}
