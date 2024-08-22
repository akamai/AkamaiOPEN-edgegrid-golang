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

func (m *Mock) GetPasswordPolicy(ctx context.Context) (*GetPasswordPolicyResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPasswordPolicyResponse), args.Error(1)
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

func (m *Mock) ListAccountSwitchKeys(ctx context.Context, request ListAccountSwitchKeysRequest) (ListAccountSwitchKeysResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(ListAccountSwitchKeysResponse), args.Error(1)
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

func (m *Mock) UpdateMFA(ctx context.Context, request UpdateMFARequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) ResetMFA(ctx context.Context, request ResetMFARequest) error {
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

func (m *Mock) ListProperties(ctx context.Context, request ListPropertiesRequest) (*ListPropertiesResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListPropertiesResponse), args.Error(1)
}

func (m *Mock) ListUsersForProperty(ctx context.Context, request ListUsersForPropertyRequest) (*ListUsersForPropertyResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListUsersForPropertyResponse), args.Error(1)
}

func (m *Mock) GetProperty(ctx context.Context, request GetPropertyRequest) (*GetPropertyResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPropertyResponse), args.Error(1)
}

func (m *Mock) MoveProperty(ctx context.Context, request MovePropertyRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) MapPropertyIDToName(ctx context.Context, request MapPropertyIDToNameRequest) (*string, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*string), args.Error(1)
}

func (m *Mock) MapPropertyNameToID(ctx context.Context, request MapPropertyNameToIDRequest) (*int64, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*int64), args.Error(1)
}

func (m *Mock) LockAPIClient(ctx context.Context, request LockAPIClientRequest) (*LockAPIClientResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*LockAPIClientResponse), args.Error(1)
}

func (m *Mock) UnlockAPIClient(ctx context.Context, request UnlockAPIClientRequest) (*UnlockAPIClientResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UnlockAPIClientResponse), args.Error(1)
}

func (m *Mock) DisableIPAllowlist(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *Mock) EnableIPAllowlist(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *Mock) GetIPAllowlistStatus(ctx context.Context) (*GetIPAllowlistStatusResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetIPAllowlistStatusResponse), args.Error(1)
}

func (m *Mock) ListAllowedCPCodes(ctx context.Context, params ListAllowedCPCodesRequest) (ListAllowedCPCodesResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(ListAllowedCPCodesResponse), args.Error(1)
}

func (m *Mock) ListCIDRBlocks(ctx context.Context, request ListCIDRBlocksRequest) (*ListCIDRBlocksResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListCIDRBlocksResponse), args.Error(1)
}

func (m *Mock) CreateCIDRBlock(ctx context.Context, request CreateCIDRBlockRequest) (*CreateCIDRBlockResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateCIDRBlockResponse), args.Error(1)
}

func (m *Mock) GetCIDRBlock(ctx context.Context, request GetCIDRBlockRequest) (*GetCIDRBlockResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCIDRBlockResponse), args.Error(1)
}

func (m *Mock) UpdateCIDRBlock(ctx context.Context, request UpdateCIDRBlockRequest) (*UpdateCIDRBlockResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateCIDRBlockResponse), args.Error(1)
}

func (m *Mock) DeleteCIDRBlock(ctx context.Context, request DeleteCIDRBlockRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) ValidateCIDRBlock(ctx context.Context, request ValidateCIDRBlockRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) BlockUsers(ctx context.Context, request BlockUsersRequest) (*BlockUsersResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*BlockUsersResponse), args.Error(1)
}

func (m *Mock) ListAuthorizedUsers(ctx context.Context) (ListAuthorizedUsersResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(ListAuthorizedUsersResponse), args.Error(1)
}

func (m *Mock) ListAllowedAPIs(ctx context.Context, request ListAllowedAPIsRequest) (ListAllowedAPIsResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(ListAllowedAPIsResponse), args.Error(1)
}

func (m *Mock) ListAccessibleGroups(ctx context.Context, request ListAccessibleGroupsRequest) (ListAccessibleGroupsResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(ListAccessibleGroupsResponse), args.Error(1)
}

func (m *Mock) CreateCredential(ctx context.Context, request CreateCredentialRequest) (*CreateCredentialResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateCredentialResponse), args.Error(1)
}

func (m *Mock) ListCredentials(ctx context.Context, request ListCredentialsRequest) (ListCredentialsResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(ListCredentialsResponse), args.Error(1)
}

func (m *Mock) GetCredential(ctx context.Context, request GetCredentialRequest) (*GetCredentialResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCredentialResponse), args.Error(1)
}

func (m *Mock) UpdateCredential(ctx context.Context, request UpdateCredentialRequest) (*UpdateCredentialResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateCredentialResponse), args.Error(1)
}

func (m *Mock) DeleteCredential(ctx context.Context, request DeleteCredentialRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) DeactivateCredential(ctx context.Context, request DeactivateCredentialRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) DeactivateCredentials(ctx context.Context, request DeactivateCredentialsRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *Mock) ListAPIClients(ctx context.Context, request ListAPIClientsRequest) (ListAPIClientsResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(ListAPIClientsResponse), args.Error(1)
}

func (m *Mock) GetAPIClient(ctx context.Context, request GetAPIClientRequest) (*GetAPIClientResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetAPIClientResponse), args.Error(1)
}

func (m *Mock) CreateAPIClient(ctx context.Context, request CreateAPIClientRequest) (*CreateAPIClientResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateAPIClientResponse), args.Error(1)
}

func (m *Mock) UpdateAPIClient(ctx context.Context, request UpdateAPIClientRequest) (*UpdateAPIClientResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateAPIClientResponse), args.Error(1)
}

func (m *Mock) DeleteAPIClient(ctx context.Context, request DeleteAPIClientRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}
