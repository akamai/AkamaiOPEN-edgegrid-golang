//revive:disable:exported

package edgeworkers

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ Edgeworkers = &Mock{}

// Activations

func (m *Mock) ListActivations(ctx context.Context, req ListActivationsRequest) (*ListActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListActivationsResponse), args.Error(1)
}

func (m *Mock) GetActivation(ctx context.Context, req GetActivationRequest) (*Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Activation), args.Error(1)
}

func (m *Mock) ActivateVersion(ctx context.Context, req ActivateVersionRequest) (*Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Activation), args.Error(1)
}

func (m *Mock) CancelPendingActivation(ctx context.Context, req CancelActivationRequest) (*Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Activation), args.Error(1)
}

// Contracts

func (m *Mock) ListContracts(ctx context.Context) (*ListContractsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListContractsResponse), args.Error(1)
}

// Deactivations

func (m *Mock) ListDeactivations(ctx context.Context, req ListDeactivationsRequest) (*ListDeactivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListDeactivationsResponse), args.Error(1)
}

func (m *Mock) GetDeactivation(ctx context.Context, req GetDeactivationRequest) (*Deactivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Deactivation), args.Error(1)
}

func (m *Mock) DeactivateVersion(ctx context.Context, req DeactivateVersionRequest) (*Deactivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Deactivation), args.Error(1)
}

// EdgeKVAccessTokens

func (m *Mock) CreateEdgeKVAccessToken(ctx context.Context, req CreateEdgeKVAccessTokenRequest) (*CreateEdgeKVAccessTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateEdgeKVAccessTokenResponse), args.Error(1)
}

func (m *Mock) GetEdgeKVAccessToken(ctx context.Context, req GetEdgeKVAccessTokenRequest) (*GetEdgeKVAccessTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEdgeKVAccessTokenResponse), args.Error(1)
}

func (m *Mock) ListEdgeKVAccessTokens(ctx context.Context, req ListEdgeKVAccessTokensRequest) (*ListEdgeKVAccessTokensResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListEdgeKVAccessTokensResponse), args.Error(1)
}

func (m *Mock) DeleteEdgeKVAccessToken(ctx context.Context, req DeleteEdgeKVAccessTokenRequest) (*DeleteEdgeKVAccessTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DeleteEdgeKVAccessTokenResponse), args.Error(1)
}

// EdgeKVInitialize

func (m *Mock) InitializeEdgeKV(ctx context.Context) (*EdgeKVInitializationStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EdgeKVInitializationStatus), args.Error(1)
}

func (m *Mock) GetEdgeKVInitializationStatus(ctx context.Context) (*EdgeKVInitializationStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EdgeKVInitializationStatus), args.Error(1)
}

// EdgeKVItems

func (m *Mock) ListItems(ctx context.Context, req ListItemsRequest) (*ListItemsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListItemsResponse), args.Error(1)
}

func (m *Mock) GetItem(ctx context.Context, req GetItemRequest) (*Item, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Item), args.Error(1)
}

func (m *Mock) UpsertItem(ctx context.Context, req UpsertItemRequest) (*string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *Mock) DeleteItem(ctx context.Context, req DeleteItemRequest) (*string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

// EdgeKVNamespaces

func (m *Mock) ListEdgeKVNamespaces(ctx context.Context, req ListEdgeKVNamespacesRequest) (*ListEdgeKVNamespacesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListEdgeKVNamespacesResponse), args.Error(1)
}

func (m *Mock) GetEdgeKVNamespace(ctx context.Context, req GetEdgeKVNamespaceRequest) (*Namespace, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Namespace), args.Error(1)
}

func (m *Mock) CreateEdgeKVNamespace(ctx context.Context, req CreateEdgeKVNamespaceRequest) (*Namespace, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Namespace), args.Error(1)
}

func (m *Mock) UpdateEdgeKVNamespace(ctx context.Context, req UpdateEdgeKVNamespaceRequest) (*Namespace, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Namespace), args.Error(1)
}

// EdgeWorkerID

func (m *Mock) GetEdgeWorkerID(ctx context.Context, req GetEdgeWorkerIDRequest) (*EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EdgeWorkerID), args.Error(1)
}

func (m *Mock) ListEdgeWorkersID(ctx context.Context, req ListEdgeWorkersIDRequest) (*ListEdgeWorkersIDResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListEdgeWorkersIDResponse), args.Error(1)
}

func (m *Mock) CreateEdgeWorkerID(ctx context.Context, req CreateEdgeWorkerIDRequest) (*EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EdgeWorkerID), args.Error(1)
}

func (m *Mock) UpdateEdgeWorkerID(ctx context.Context, req UpdateEdgeWorkerIDRequest) (*EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EdgeWorkerID), args.Error(1)
}

func (m *Mock) CloneEdgeWorkerID(ctx context.Context, req CloneEdgeWorkerIDRequest) (*EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EdgeWorkerID), args.Error(1)
}

func (m *Mock) DeleteEdgeWorkerID(ctx context.Context, req DeleteEdgeWorkerIDRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// EdgeWorkerVersion

func (m *Mock) GetEdgeWorkerVersion(ctx context.Context, req GetEdgeWorkerVersionRequest) (*EdgeWorkerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EdgeWorkerVersion), args.Error(1)
}

func (m *Mock) ListEdgeWorkerVersions(ctx context.Context, req ListEdgeWorkerVersionsRequest) (*ListEdgeWorkerVersionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListEdgeWorkerVersionsResponse), args.Error(1)
}

func (m *Mock) GetEdgeWorkerVersionContent(ctx context.Context, req GetEdgeWorkerVersionContentRequest) (*Bundle, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Bundle), args.Error(1)
}

func (m *Mock) CreateEdgeWorkerVersion(ctx context.Context, req CreateEdgeWorkerVersionRequest) (*EdgeWorkerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EdgeWorkerVersion), args.Error(1)
}

func (m *Mock) DeleteEdgeWorkerVersion(ctx context.Context, req DeleteEdgeWorkerVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// PermissionGroups

func (m *Mock) GetPermissionGroup(ctx context.Context, req GetPermissionGroupRequest) (*PermissionGroup, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PermissionGroup), args.Error(1)
}

func (m *Mock) ListPermissionGroups(ctx context.Context) (*ListPermissionGroupsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListPermissionGroupsResponse), args.Error(1)
}

// Properties

func (m *Mock) ListProperties(ctx context.Context, req ListPropertiesRequest) (*ListPropertiesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListPropertiesResponse), args.Error(1)
}

// Reports

func (m *Mock) GetReport(ctx context.Context, req GetReportRequest) (*GetReportResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReportResponse), args.Error(1)
}

func (m *Mock) GetSummaryReport(ctx context.Context, req GetSummaryReportRequest) (*GetSummaryReportResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSummaryReportResponse), args.Error(1)
}

func (m *Mock) ListReports(ctx context.Context) (*ListReportsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListReportsResponse), args.Error(1)
}

// ResourceTiers

func (m *Mock) ListResourceTiers(ctx context.Context, req ListResourceTiersRequest) (*ListResourceTiersResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListResourceTiersResponse), args.Error(1)
}

func (m *Mock) GetResourceTier(ctx context.Context, req GetResourceTierRequest) (*ResourceTier, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ResourceTier), args.Error(1)
}

// SecureTokens

func (m *Mock) CreateSecureToken(ctx context.Context, req CreateSecureTokenRequest) (*CreateSecureTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateSecureTokenResponse), args.Error(1)
}

// Validations

func (m *Mock) ValidateBundle(ctx context.Context, req ValidateBundleRequest) (*ValidateBundleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ValidateBundleResponse), args.Error(1)
}
