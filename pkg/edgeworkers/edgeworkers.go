// Package edgeworkers provides access to the Akamai EdgeWorkers and EdgeKV APIs
package edgeworkers

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Edgeworkers is the api interface for EdgeWorkers and EdgeKV
	Edgeworkers interface {
		// Activations

		// ListActivations lists all activations for an EdgeWorker
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-activations-1
		ListActivations(context.Context, ListActivationsRequest) (*ListActivationsResponse, error)

		// GetActivation fetches an EdgeWorker activation by id
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-activation-1
		GetActivation(context.Context, GetActivationRequest) (*Activation, error)

		// ActivateVersion activates an EdgeWorker version on a given network
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/post-activations-1
		ActivateVersion(context.Context, ActivateVersionRequest) (*Activation, error)

		// CancelPendingActivation cancels pending activation with a given id
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/cancel-activation
		CancelPendingActivation(context.Context, CancelActivationRequest) (*Activation, error)

		// Contracts

		// ListContracts lists contract IDs that can be used to list resource tiers
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-contracts-1
		ListContracts(context.Context) (*ListContractsResponse, error)

		// Deactivations

		// ListDeactivations lists all deactivations for a given EdgeWorker ID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-deactivations-1
		ListDeactivations(context.Context, ListDeactivationsRequest) (*ListDeactivationsResponse, error)

		// GetDeactivation gets details for a specific deactivation
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-deactivation-1
		GetDeactivation(context.Context, GetDeactivationRequest) (*Deactivation, error)

		// DeactivateVersion deactivates an existing EdgeWorker version on the akamai network
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/post-deactivations-1
		DeactivateVersion(context.Context, DeactivateVersionRequest) (*Deactivation, error)

		// EdgeKVAccessTokens

		// CreateEdgeKVAccessToken generates EdgeKV specific access token
		//
		// See: https://techdocs.akamai.com/edgekv/reference/post-tokens
		CreateEdgeKVAccessToken(context.Context, CreateEdgeKVAccessTokenRequest) (*CreateEdgeKVAccessTokenResponse, error)

		// GetEdgeKVAccessToken retrieves an EdgeKV access token
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-token
		GetEdgeKVAccessToken(context.Context, GetEdgeKVAccessTokenRequest) (*GetEdgeKVAccessTokenResponse, error)

		// ListEdgeKVAccessTokens lists EdgeKV access tokens
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-tokens
		ListEdgeKVAccessTokens(context.Context, ListEdgeKVAccessTokensRequest) (*ListEdgeKVAccessTokensResponse, error)

		// DeleteEdgeKVAccessToken revokes an EdgeKV access token
		//
		// See: https://techdocs.akamai.com/edgekv/reference/delete-token
		DeleteEdgeKVAccessToken(context.Context, DeleteEdgeKVAccessTokenRequest) (*DeleteEdgeKVAccessTokenResponse, error)

		// EdgeKVInitialize

		// InitializeEdgeKV Initialize the EdgeKV database
		//
		// See: https://techdocs.akamai.com/edgekv/reference/put-initialize
		InitializeEdgeKV(ctx context.Context) (*EdgeKVInitializationStatus, error)

		// GetEdgeKVInitializationStatus is used to check on the current initialization status
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-initialize
		GetEdgeKVInitializationStatus(ctx context.Context) (*EdgeKVInitializationStatus, error)

		// EdgeKVItems

		// ListItems lists items in EdgeKV group
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-group-1
		ListItems(context.Context, ListItemsRequest) (*ListItemsResponse, error)

		// GetItem reads an item from EdgeKV group
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-item
		GetItem(context.Context, GetItemRequest) (*Item, error)

		// UpsertItem creates or updates an item in EdgeKV group
		//
		// See: https://techdocs.akamai.com/edgekv/reference/put-item
		UpsertItem(context.Context, UpsertItemRequest) (*string, error)

		// DeleteItem deletes an item from EdgeKV group
		//
		// See: https://techdocs.akamai.com/edgekv/reference/delete-item
		DeleteItem(context.Context, DeleteItemRequest) (*string, error)

		// EdgeKVNamespaces

		// ListEdgeKVNamespaces lists all namespaces in the given network
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-namespaces
		ListEdgeKVNamespaces(context.Context, ListEdgeKVNamespacesRequest) (*ListEdgeKVNamespacesResponse, error)

		// GetEdgeKVNamespace fetches a namespace by name
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-namespace
		GetEdgeKVNamespace(context.Context, GetEdgeKVNamespaceRequest) (*Namespace, error)

		// CreateEdgeKVNamespace creates a namespace on the given network
		//
		// See: https://techdocs.akamai.com/edgekv/reference/post-namespace
		CreateEdgeKVNamespace(context.Context, CreateEdgeKVNamespaceRequest) (*Namespace, error)

		// UpdateEdgeKVNamespace updates a namespace
		//
		// See: https://techdocs.akamai.com/edgekv/reference/put-namespace
		UpdateEdgeKVNamespace(context.Context, UpdateEdgeKVNamespaceRequest) (*Namespace, error)

		// DeleteEdgeKVNamespace deletes a namespace and all of its contents.
		//
		// See: https://techdocs.akamai.com/edgekv/reference/delete-namespace
		DeleteEdgeKVNamespace(context.Context, DeleteEdgeKVNamespaceRequest) (*DeleteEdgeKVNamespacesResponse, error)

		// GetNamespaceScheduledDeleteTime gets the scheduled time for a namespace delete.
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-scheduled-delete
		GetNamespaceScheduledDeleteTime(context.Context, GetScheduledDeleteTimeRequest) (*ScheduledDeleteTimeResponse, error)

		// RescheduleNamespaceDelete changes the scheduled time of a namespace delete.
		//
		// See: https://techdocs.akamai.com/edgekv/reference/put-scheduled-delete
		RescheduleNamespaceDelete(context.Context, RescheduleNamespaceDeleteRequest) (*RescheduleNamespaceDeleteResponse, error)

		// CancelScheduledNamespaceDelete deletes the scheduled time for a namespace delete, effectively canceling the deletion.
		//
		// See: https://techdocs.akamai.com/edgekv/reference/delete-scheduled-delete
		CancelScheduledNamespaceDelete(context.Context, CancelScheduledNamespaceDeleteRequest) error

		// EdgeWorkerIDs

		// GetEdgeWorkerID gets details for a specific EdgeWorkerID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-id
		GetEdgeWorkerID(context.Context, GetEdgeWorkerIDRequest) (*EdgeWorkerID, error)

		// ListEdgeWorkersID lists EdgeWorkerIDs in the identified group
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-ids
		ListEdgeWorkersID(context.Context, ListEdgeWorkersIDRequest) (*ListEdgeWorkersIDResponse, error)

		// CreateEdgeWorkerID creates a new EdgeWorkerID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/post-ids
		CreateEdgeWorkerID(context.Context, CreateEdgeWorkerIDRequest) (*EdgeWorkerID, error)

		// UpdateEdgeWorkerID updates an EdgeWorkerID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/put-id
		UpdateEdgeWorkerID(context.Context, UpdateEdgeWorkerIDRequest) (*EdgeWorkerID, error)

		// CloneEdgeWorkerID clones an EdgeWorkerID to change the resource tier
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/post-id-clone
		CloneEdgeWorkerID(context.Context, CloneEdgeWorkerIDRequest) (*EdgeWorkerID, error)

		// DeleteEdgeWorkerID deletes an EdgeWorkerID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/delete-id
		DeleteEdgeWorkerID(context.Context, DeleteEdgeWorkerIDRequest) error

		// EdgeWorkerVersions

		// GetEdgeWorkerVersion gets details for a specific EdgeWorkerVersion
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-version
		GetEdgeWorkerVersion(context.Context, GetEdgeWorkerVersionRequest) (*EdgeWorkerVersion, error)

		// ListEdgeWorkerVersions lists EdgeWorkerVersions in the identified group
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-versions
		ListEdgeWorkerVersions(context.Context, ListEdgeWorkerVersionsRequest) (*ListEdgeWorkerVersionsResponse, error)

		// GetEdgeWorkerVersionContent gets content bundle for a specific EdgeWorkerVersion
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-version-content
		GetEdgeWorkerVersionContent(context.Context, GetEdgeWorkerVersionContentRequest) (*Bundle, error)

		// CreateEdgeWorkerVersion creates a new EdgeWorkerVersion
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/post-versions
		CreateEdgeWorkerVersion(context.Context, CreateEdgeWorkerVersionRequest) (*EdgeWorkerVersion, error)

		// DeleteEdgeWorkerVersion deletes an EdgeWorkerVersion
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/delete-version
		DeleteEdgeWorkerVersion(context.Context, DeleteEdgeWorkerVersionRequest) error

		// Groups

		// ListGroupsWithinNamespace lists group identifiers created when writing items to a namespace
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get-namespace-groups
		ListGroupsWithinNamespace(context.Context, ListGroupsWithinNamespaceRequest) ([]string, error)

		// PermissionGroups

		// GetPermissionGroup gets details on the capabilities enabled within a specified group
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-group
		GetPermissionGroup(context.Context, GetPermissionGroupRequest) (*PermissionGroup, error)

		// ListPermissionGroups lists groups and the associated permission capabilities
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-groups
		ListPermissionGroups(context.Context) (*ListPermissionGroupsResponse, error)

		// Properties

		// ListProperties lists all properties for a given edgeworker ID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-properties
		ListProperties(context.Context, ListPropertiesRequest) (*ListPropertiesResponse, error)

		// Reports

		// GetSummaryReport gets summary overview for EdgeWorker reports. Report id is  1
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-report
		GetSummaryReport(context.Context, GetSummaryReportRequest) (*GetSummaryReportResponse, error)

		// GetReport gets details for an EdgeWorker
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-report
		GetReport(context.Context, GetReportRequest) (*GetReportResponse, error)

		// ListReports lists EdgeWorker reports
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-reports
		ListReports(context.Context) (*ListReportsResponse, error)

		// ResourceTiers

		// ListResourceTiers lists all resource tiers for a given contract
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-resource-tiers
		ListResourceTiers(context.Context, ListResourceTiersRequest) (*ListResourceTiersResponse, error)

		// GetResourceTier returns resource tier for a given edgeworker ID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-id-resource-tier
		GetResourceTier(context.Context, GetResourceTierRequest) (*ResourceTier, error)

		// SecureTokens

		// CreateSecureToken creates a new secure token
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/post-secure-token
		CreateSecureToken(context.Context, CreateSecureTokenRequest) (*CreateSecureTokenResponse, error)

		// Validations

		// ValidateBundle given bundle validates it and returns a list of errors and/or warnings
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/post-validations
		ValidateBundle(context.Context, ValidateBundleRequest) (*ValidateBundleResponse, error)
	}

	edgeworkers struct {
		session.Session
	}

	// Option defines an Edgeworkers option
	Option func(*edgeworkers)

	// ClientFunc is a Edgeworkers client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) Edgeworkers
)

// Client returns a new edgeworkers Client instance with the specified controller
func Client(sess session.Session, opts ...Option) Edgeworkers {
	e := &edgeworkers{
		Session: sess,
	}

	for _, opt := range opts {
		opt(e)
	}
	return e
}
