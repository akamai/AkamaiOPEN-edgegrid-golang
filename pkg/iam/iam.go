// Package iam provides access to the Akamai Property APIs
package iam

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// IAM is the IAM api interface
	IAM interface {

		// API Clients

		// LockAPIClient locks an API client based on `ClientID` parameter. If `ClientID` is not provided, it locks your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-lock-api-client, https://techdocs.akamai.com/iam-api/reference/put-lock-api-client-self
		LockAPIClient(ctx context.Context, params LockAPIClientRequest) (*LockAPIClientResponse, error)

		// UnlockAPIClient unlocks an API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-unlock-api-client
		UnlockAPIClient(ctx context.Context, params UnlockAPIClientRequest) (*UnlockAPIClientResponse, error)

		// ListAPIClients lists API clients an administrator can manage.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-clients
		ListAPIClients(ctx context.Context, params ListAPIClientsRequest) (ListAPIClientsResponse, error)

		// GetAPIClient provides details about an API client. If `ClientID` is not provided, it returns details about your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-client and https://techdocs.akamai.com/iam-api/reference/get-api-client-self
		GetAPIClient(ctx context.Context, params GetAPIClientRequest) (*GetAPIClientResponse, error)

		// CreateAPIClient creates a new API client. Optionally, it can automatically assign a credential for the client when creating it.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-api-clients
		CreateAPIClient(ctx context.Context, params CreateAPIClientRequest) (*CreateAPIClientResponse, error)

		// UpdateAPIClient updates an API client. If `ClientID` is not provided, it updates your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-api-clients and https://techdocs.akamai.com/iam-api/reference/put-api-clients-self
		UpdateAPIClient(ctx context.Context, params UpdateAPIClientRequest) (*UpdateAPIClientResponse, error)

		// DeleteAPIClient permanently deletes the API client, breaking any API connections with the client.
		// If `ClientID` is not provided, it deletes your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/delete-api-client and https://techdocs.akamai.com/iam-api/reference/delete-api-client-self
		DeleteAPIClient(ctx context.Context, params DeleteAPIClientRequest) error

		// API Client Credentials

		// CreateCredential creates a new credential for the API client.  If `ClientID` is not provided, it creates credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-self-credentials, https://techdocs.akamai.com/iam-api/reference/post-client-credentials
		CreateCredential(context.Context, CreateCredentialRequest) (*CreateCredentialResponse, error)

		// ListCredentials lists credentials for an API client. If `ClientID` is not provided, it lists credentials for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-self-credentials, https://techdocs.akamai.com/iam-api/reference/get-client-credentials
		ListCredentials(context.Context, ListCredentialsRequest) (ListCredentialsResponse, error)

		// GetCredential returns details about a specific credential for an API client. If `ClientID` is not provided, it gets credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-self-credential, https://techdocs.akamai.com/iam-api/reference/get-client-credential
		GetCredential(context.Context, GetCredentialRequest) (*GetCredentialResponse, error)

		// UpdateCredential updates a specific credential for an API client. If `ClientID` is not provided, it updates credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-self-credential, https://techdocs.akamai.com/iam-api/reference/put-client-credential
		UpdateCredential(context.Context, UpdateCredentialRequest) (*UpdateCredentialResponse, error)

		// DeleteCredential deletes a specific credential from an API client. If `ClientID` is not provided, it deletes credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/delete-self-credential, https://techdocs.akamai.com/iam-api/reference/delete-client-credential
		DeleteCredential(context.Context, DeleteCredentialRequest) error

		// DeactivateCredential deactivates a specific credential for an API client. If `ClientID` is not provided, it deactivates credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-self-credential-deactivate, https://techdocs.akamai.com/iam-api/reference/post-client-credential-deactivate
		DeactivateCredential(context.Context, DeactivateCredentialRequest) error

		// DeactivateCredentials deactivates all credentials for a specific API client. If `ClientID` is not provided, it deactivates all credentials for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-self-credentials-deactivate, https://techdocs.akamai.com/iam-api/reference/post-client-credentials-deactivate
		DeactivateCredentials(context.Context, DeactivateCredentialsRequest) error

		// Blocked Properties

		// ListBlockedProperties returns all properties a user doesn't have access to in a group.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-blocked-properties
		ListBlockedProperties(context.Context, ListBlockedPropertiesRequest) ([]int64, error)

		// UpdateBlockedProperties removes or grants user access to properties.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-blocked-properties
		UpdateBlockedProperties(context.Context, UpdateBlockedPropertiesRequest) ([]int64, error)

		// CIDR Blocks

		// ListCIDRBlocks lists all CIDR blocks on selected account's allowlist.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-allowlist
		ListCIDRBlocks(context.Context, ListCIDRBlocksRequest) (ListCIDRBlocksResponse, error)

		// CreateCIDRBlock adds CIDR blocks to your account's allowlist.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-allowlist
		CreateCIDRBlock(context.Context, CreateCIDRBlockRequest) (*CreateCIDRBlockResponse, error)

		// GetCIDRBlock retrieves a CIDR block's details.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-allowlist-cidrblockid
		GetCIDRBlock(context.Context, GetCIDRBlockRequest) (*GetCIDRBlockResponse, error)

		// UpdateCIDRBlock modifies an existing CIDR block.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-allowlist-cidrblockid
		UpdateCIDRBlock(context.Context, UpdateCIDRBlockRequest) (*UpdateCIDRBlockResponse, error)

		// DeleteCIDRBlock deletes an existing CIDR block from the IP allowlist.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/delete-allowlist-cidrblockid
		DeleteCIDRBlock(context.Context, DeleteCIDRBlockRequest) error

		// ValidateCIDRBlock checks the format of CIDR block.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-allowlist-validate
		ValidateCIDRBlock(context.Context, ValidateCIDRBlockRequest) error

		// Groups

		// CreateGroup creates a new group within a parent group_id specified in the request.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-group
		CreateGroup(context.Context, GroupRequest) (*Group, error)

		// GetGroup returns a group's details.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-group
		GetGroup(context.Context, GetGroupRequest) (*Group, error)

		// ListAffectedUsers lists users who are affected when a group is moved.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-move-affected-users
		ListAffectedUsers(context.Context, ListAffectedUsersRequest) ([]GroupUser, error)

		// ListGroups lists all groups in which you have a scope of admin for the current account and contract type.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-groups
		ListGroups(context.Context, ListGroupsRequest) ([]Group, error)

		// RemoveGroup removes a group based on group_id. We can only delete a sub-group, and only if that sub-group doesn't include any users.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/delete-group
		RemoveGroup(context.Context, RemoveGroupRequest) error

		// UpdateGroupName changes the name of the group.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-group
		UpdateGroupName(context.Context, GroupRequest) (*Group, error)

		// MoveGroup moves a nested group under another group within the same parent hierarchy.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-groups-move
		MoveGroup(context.Context, MoveGroupRequest) error

		// Helpers

		// ListAllowedCPCodes lists available CP codes for a user.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-api-clients-users-allowed-cpcodes
		ListAllowedCPCodes(context.Context, ListAllowedCPCodesRequest) (ListAllowedCPCodesResponse, error)

		// ListAuthorizedUsers lists authorized API client users.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-clients-users
		ListAuthorizedUsers(context.Context) (ListAuthorizedUsersResponse, error)

		// ListAllowedAPIs lists available APIs for a user.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-clients-users-allowed-apis
		ListAllowedAPIs(context.Context, ListAllowedAPIsRequest) (ListAllowedAPIsResponse, error)

		// ListAccessibleGroups lists groups available to a user.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-clients-users-group-access
		ListAccessibleGroups(context.Context, ListAccessibleGroupsRequest) (ListAccessibleGroupsResponse, error)

		// IP Allowlist

		// DisableIPAllowlist disables IP allowlist on your account. After you disable IP allowlist,
		// users can access Control Center regardless of their IP address or who assigns it.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-allowlist-disable
		DisableIPAllowlist(context.Context) error

		// EnableIPAllowlist enables IP allowlist on your account. Before you enable IP allowlist,
		// add at least one IP address to allow access to Control Center.
		// The allowlist can't be empty with IP allowlist enabled.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-allowlist-enable
		EnableIPAllowlist(context.Context) error

		// GetIPAllowlistStatus indicates whether IP allowlist is enabled or disabled on your account.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-allowlist-status
		GetIPAllowlistStatus(context.Context) (*GetIPAllowlistStatusResponse, error)

		// Properties

		// ListProperties lists the properties for the current account or other managed accounts using the accountSwitchKey parameter.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-properties
		ListProperties(context.Context, ListPropertiesRequest) (ListPropertiesResponse, error)

		// ListUsersForProperty lists users who can access a property.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-property-users
		ListUsersForProperty(context.Context, ListUsersForPropertyRequest) (ListUsersForPropertyResponse, error)

		// GetProperty lists a property's details.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-property
		GetProperty(context.Context, GetPropertyRequest) (*GetPropertyResponse, error)

		// MoveProperty moves a property from one group to another group.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-property
		MoveProperty(context.Context, MovePropertyRequest) error

		// MapPropertyIDToName returns property name for given (IAM) property ID
		// Mainly to be used to map (IAM) Property ID to (PAPI) Property ID
		// To finish the mapping, please use papi.MapPropertyNameToID
		MapPropertyIDToName(context.Context, MapPropertyIDToNameRequest) (*string, error)

		// BlockUsers blocks the users on a property.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-property-users-block
		BlockUsers(context.Context, BlockUsersRequest) (*BlockUsersResponse, error)

		// Roles

		// CreateRole creates a custom role.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-role
		CreateRole(context.Context, CreateRoleRequest) (*Role, error)

		// GetRole gets details for a specific role.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-role
		GetRole(context.Context, GetRoleRequest) (*Role, error)

		// UpdateRole adds or removes permissions from a role and updates other parameters.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-role
		UpdateRole(context.Context, UpdateRoleRequest) (*Role, error)

		// DeleteRole deletes a role. This operation is only allowed if the role isn't assigned to any users.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/delete-role
		DeleteRole(context.Context, DeleteRoleRequest) error

		// ListRoles lists roles for the current account and contract type.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-roles
		ListRoles(context.Context, ListRolesRequest) ([]Role, error)

		// ListGrantableRoles lists which grantable roles can be included in a new custom role or added to an existing custom role.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-grantable-roles
		ListGrantableRoles(context.Context) ([]RoleGrantedRole, error)

		// Support

		// GetPasswordPolicy gets the password policy for the account.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-password-policy
		GetPasswordPolicy(ctx context.Context) (*GetPasswordPolicyResponse, error)

		// ListProducts lists products a user can subscribe to and receive notifications for on the account.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-notification-products
		ListProducts(context.Context) ([]string, error)

		// ListStates lists U.S. states or Canadian provinces.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-states
		ListStates(context.Context, ListStatesRequest) ([]string, error)

		// ListTimeoutPolicies lists all the possible session timeout policies.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-timeout-policies
		ListTimeoutPolicies(context.Context) ([]TimeoutPolicy, error)

		// ListAccountSwitchKeys lists account switch keys available for a specific API client. If `ClientID` is not provided, it lists account switch keys available for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-client-account-switch-keys, https://techdocs.akamai.com/iam-api/reference/get-self-account-switch-keys
		ListAccountSwitchKeys(context.Context, ListAccountSwitchKeysRequest) (ListAccountSwitchKeysResponse, error)

		// SupportedContactTypes lists supported contact types.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-contact-types
		SupportedContactTypes(context.Context) ([]string, error)

		// SupportedCountries lists supported countries.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-countries
		SupportedCountries(context.Context) ([]string, error)

		// SupportedLanguages lists supported languages.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-languages
		SupportedLanguages(context.Context) ([]string, error)

		// SupportedTimezones lists supported timezones.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-timezones
		SupportedTimezones(context.Context) ([]Timezone, error)

		// Users

		// LockUser locks the user.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-ui-identity-lock
		LockUser(context.Context, LockUserRequest) error

		// UnlockUser releases the lock on a user's account.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-ui-identity-unlock
		UnlockUser(context.Context, UnlockUserRequest) error

		// ResetUserPassword optionally sends a one-time use password to the user.
		// If you send the email with the password directly to the user, the response for this operation doesn't include that password.
		// If you don't send the password to the user through email, the password is included in the response.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-reset-password
		ResetUserPassword(context.Context, ResetUserPasswordRequest) (*ResetUserPasswordResponse, error)

		// SetUserPassword sets a specific password for a user.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-set-password
		SetUserPassword(context.Context, SetUserPasswordRequest) error

		// CreateUser creates a user in the account specified in your own API client credentials or clone an existing user's role assignments.
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/post-ui-identity
		CreateUser(context.Context, CreateUserRequest) (*User, error)

		// GetUser gets  a specific user's profile.
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-ui-identity
		GetUser(context.Context, GetUserRequest) (*User, error)

		// ListUsers returns a list of users who have access on this account.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-ui-identities
		ListUsers(context.Context, ListUsersRequest) ([]UserListItem, error)

		// RemoveUser removes a user identity.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/delete-ui-identity
		RemoveUser(context.Context, RemoveUserRequest) error

		// UpdateUserAuthGrants edits what groups a user has access to, and how the user can interact with the objects in those groups.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-ui-uiidentity-auth-grants
		UpdateUserAuthGrants(context.Context, UpdateUserAuthGrantsRequest) ([]AuthGrant, error)

		// UpdateUserInfo updates a user's information.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-ui-identity-basic-info
		UpdateUserInfo(context.Context, UpdateUserInfoRequest) (*UserBasicInfo, error)

		// UpdateUserNotifications subscribes or un-subscribes user to product notification emails.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-notifications
		UpdateUserNotifications(context.Context, UpdateUserNotificationsRequest) (*UserNotifications, error)

		// UpdateMFA updates a user's profile authentication method.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-user-profile-additional-authentication
		UpdateMFA(context.Context, UpdateMFARequest) error

		// ResetMFA resets a user's profile authentication method.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-ui-identity-reset-additional-authentication
		ResetMFA(context.Context, ResetMFARequest) error
	}

	iam struct {
		session.Session
	}

	// Option defines a IAM option.
	Option func(*iam)

	// ClientFunc is an IAM client new method, this can be used for mocking.
	ClientFunc func(sess session.Session, opts ...Option) IAM
)

// Client returns a new IAM Client instance with the specified controller.
func Client(sess session.Session, opts ...Option) IAM {
	p := &iam{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
