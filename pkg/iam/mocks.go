//revive:disable:exported

package iam

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ IAM = &Mock{}

func (m *Mock) ListGroups(ctx context.Context, request ListGroupsRequest) ([]Group, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Group), args.Error(1)
}

func (m *Mock) ListRoles(ctx context.Context, request ListRolesRequest) ([]Role, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Role), args.Error(1)
}

func (m *Mock) SupportedCountries(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *Mock) SupportedContactTypes(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *Mock) SupportedLanguages(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *Mock) SupportedTimezones(ctx context.Context) ([]Timezone, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Timezone), args.Error(1)
}

func (m *Mock) ListProducts(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *Mock) ListTimeoutPolicies(ctx context.Context) ([]TimeoutPolicy, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]TimeoutPolicy), args.Error(1)
}

func (m *Mock) ListStates(ctx context.Context, request ListStatesRequest) ([]string, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *Mock) CreateUser(ctx context.Context, request CreateUserRequest) (*User, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (m *Mock) GetUser(ctx context.Context, request GetUserRequest) (*User, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (m *Mock) UpdateUserInfo(ctx context.Context, request UpdateUserInfoRequest) (*UserBasicInfo, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UserBasicInfo), args.Error(1)
}

func (m *Mock) UpdateUserNotifications(ctx context.Context, request UpdateUserNotificationsRequest) (*UserNotifications, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UserNotifications), args.Error(1)
}

func (m *Mock) UpdateUserAuthGrants(ctx context.Context, request UpdateUserAuthGrantsRequest) ([]AuthGrant, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]AuthGrant), args.Error(1)
}

func (m *Mock) RemoveUser(ctx context.Context, request RemoveUserRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) ListBlockedProperties(ctx context.Context, request ListBlockedPropertiesRequest) ([]int64, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]int64), args.Error(1)
}

func (m *Mock) UpdateBlockedProperties(ctx context.Context, request UpdateBlockedPropertiesRequest) ([]int64, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]int64), args.Error(1)
}

func (m *Mock) CreateGroup(ctx context.Context, request GroupRequest) (*Group, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Group), args.Error(1)
}

func (m *Mock) GetGroup(ctx context.Context, request GetGroupRequest) (*Group, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Group), args.Error(1)
}

func (m *Mock) ListAffectedUsers(ctx context.Context, request ListAffectedUsersRequest) ([]GroupUser, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]GroupUser), args.Error(1)
}

func (m *Mock) RemoveGroup(ctx context.Context, request RemoveGroupRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) UpdateGroupName(ctx context.Context, request GroupRequest) (*Group, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Group), args.Error(1)
}

func (m *Mock) MoveGroup(ctx context.Context, request MoveGroupRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) CreateRole(ctx context.Context, request CreateRoleRequest) (*Role, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Role), args.Error(1)
}

func (m *Mock) GetRole(ctx context.Context, request GetRoleRequest) (*Role, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Role), args.Error(1)
}

func (m *Mock) UpdateRole(ctx context.Context, request UpdateRoleRequest) (*Role, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Role), args.Error(1)
}

func (m *Mock) DeleteRole(ctx context.Context, request DeleteRoleRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) ListGrantableRoles(ctx context.Context) ([]RoleGrantedRole, error) {
	args := m.Called(ctx)
	return args.Get(0).([]RoleGrantedRole), args.Error(1)
}

func (m *Mock) LockUser(ctx context.Context, request LockUserRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) UnlockUser(ctx context.Context, request UnlockUserRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) ListUsers(ctx context.Context, request ListUsersRequest) ([]UserListItem, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]UserListItem), args.Error(1)
}

func (m *Mock) UpdateTFA(ctx context.Context, request UpdateTFARequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) ResetUserPassword(ctx context.Context, request ResetUserPasswordRequest) (*ResetUserPasswordResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResetUserPasswordResponse), args.Error(1)
}

func (m *Mock) SetUserPassword(ctx context.Context, request SetUserPasswordRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}
