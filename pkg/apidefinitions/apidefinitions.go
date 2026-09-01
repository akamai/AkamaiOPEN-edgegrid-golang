// Package apidefinitions provides access to the Akamai APIDefinitions API
package apidefinitions

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// APIDefinitions is the api definitions API interface
	APIDefinitions interface {

		// GetEndpoint returns information about API endpoint
		GetEndpoint(context.Context, GetEndpointRequest) (*GetEndpointResponse, error)

		// RegisterEndpoint creates the first version of an API endpoint configuration
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/post-endpoints
		RegisterEndpoint(context.Context, RegisterEndpointRequest) (*RegisterEndpointResponse, error)

		// RegisterEndpointFromFile imports an API definition file and creates a new endpoint based on the file contents
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/post-endpoints-file
		RegisterEndpointFromFile(context.Context, RegisterEndpointFromFileRequest) (*RegisterEndpointFromFileResponse, error)

		// ShowEndpoint reveals a hidden endpoint and all its versions
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/post-endpoint-show
		ShowEndpoint(context.Context, ShowEndpointRequest) (*ShowEndpointResponse, error)

		// HideEndpoint hides an endpoint and all its versions
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/post-endpoint-hide
		HideEndpoint(context.Context, HideEndpointRequest) (*HideEndpointResponse, error)

		// DeleteEndpoint removes an endpoint configuration from API Gateway
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/delete-endpoint
		DeleteEndpoint(context.Context, DeleteEndpointRequest) error

		// ListEndpoints lists the available API endpoints, with results optionally paginated, sorted, and filtered
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/get-endpoints
		ListEndpoints(context.Context, ListEndpointsRequest) (*ListEndpointsResponse, error)

		// ListUserEntitlements lists user entitlements based on your assigned permissions
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/get-user-entitlements
		ListUserEntitlements(context.Context) (ListUserEntitlementsResponse, error)

		// VerifyVersion verify an endpoint version on the staging or production network
		VerifyVersion(context.Context, VerifyVersionRequest) (VerifyVersionResponse, error)

		// ActivateVersion activates an endpoint version on the staging or production network
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/post-endpoint-version-activate
		ActivateVersion(context.Context, ActivateVersionRequest) (*ActivateVersionResponse, error)

		// DeactivateVersion deactivates an endpoint version on the staging or production network
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/post-endpoint-version-deactivate
		DeactivateVersion(context.Context, DeactivateVersionRequest) (*DeactivateVersionResponse, error)

		// ListEndpointVersions returns all versions of an endpoint
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/get-endpoint-versions
		ListEndpointVersions(context.Context, ListEndpointVersionsRequest) (*ListEndpointVersionsResponse, error)

		// GetEndpointVersion returns an endpoint version. Use this operation's response object when modifying an endpoint version through Edit a version.
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/get-version-details
		GetEndpointVersion(context.Context, GetEndpointVersionRequest) (*GetEndpointVersionResponse, error)

		// UpdateEndpointVersion updates details about an endpoint version that has never been activated on the staging or production network.
		// You can configure the endpoint's security settings and other top-level metadata, or modify the endpoint's entire set of resources as an alternative to separate calls to Edit a resource.
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/put-endpoint-version
		UpdateEndpointVersion(context.Context, UpdateEndpointVersionRequest) (*UpdateEndpointVersionResponse, error)

		// CloneEndpointVersion creates a new endpoint version as a clone of an existing version.
		// The system assigns a new number to the version that you clone
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/post-endpoint-version-clone
		CloneEndpointVersion(context.Context, CloneEndpointVersionRequest) (*CloneEndpointVersionResponse, error)

		// DeleteEndpointVersion removes an endpoint version from API Gateway if the version has never been activated on the staging or production network
		//
		// See: https://techdocs.akamai.com/api-definitions/reference/delete-endpoint-version
		DeleteEndpointVersion(context.Context, DeleteEndpointVersionRequest) error

		// SearchResourceOperations searches for resources and operations within an API endpoint
		SearchResourceOperations(ctx context.Context) (*SearchResourceOperationsResponse, error)
	}

	apidefinitions struct {
		session.Session
	}

	// Option defines a api definition option
	Option func(*apidefinitions)

	// ClientFunc is a apidefinitions client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) APIDefinitions
)

// Client returns a new apidefinitions Client instance with the specified controller
func Client(sess session.Session, opts ...Option) APIDefinitions {
	a := &apidefinitions{
		Session: sess,
	}

	for _, opt := range opts {
		opt(a)
	}
	return a
}
