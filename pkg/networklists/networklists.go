// Package networklists provides access to the Akamai Networklist APIs
//
// See: https://techdocs.akamai.com/network-lists/reference/api#networklist
package networklists

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// NetworkList is the networklist api interface
	NetworkList interface {
		// Activations

		// GetActivations retrieves list of network list activations.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/get-network-list-status
		GetActivations(ctx context.Context, params GetActivationsRequest) (*GetActivationsResponse, error)

		// GetActivation retrieves network list activation.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/get-activation
		GetActivation(ctx context.Context, params GetActivationRequest) (*GetActivationResponse, error)

		// CreateActivations activates network list.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/post-network-list-activate
		CreateActivations(ctx context.Context, params CreateActivationsRequest) (*CreateActivationsResponse, error)

		// RemoveActivations deactivates network list.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/post-network-list-activate
		RemoveActivations(ctx context.Context, params RemoveActivationsRequest) (*RemoveActivationsResponse, error)

		// NetworkList

		// GetNetworkLists lists all network lists available for an authenticated user.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/get-network-lists
		GetNetworkLists(ctx context.Context, params GetNetworkListsRequest) (*GetNetworkListsResponse, error)

		// GetNetworkList retrieves network list with specific network list id.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/get-network-list
		GetNetworkList(ctx context.Context, params GetNetworkListRequest) (*GetNetworkListResponse, error)

		// CreateNetworkList creates a new network list.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/post-network-lists
		CreateNetworkList(ctx context.Context, params CreateNetworkListRequest) (*CreateNetworkListResponse, error)

		// UpdateNetworkList modifies the network list.
		//
		//See: https://techdocs.akamai.com/network-lists/reference/put-network-list
		UpdateNetworkList(ctx context.Context, params UpdateNetworkListRequest) (*UpdateNetworkListResponse, error)

		// RemoveNetworkList removes a network list.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/delete-network-list
		RemoveNetworkList(ctx context.Context, params RemoveNetworkListRequest) (*RemoveNetworkListResponse, error)

		// NetworkListDescription

		// GetNetworkListDescription retrieves network list with description.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/get-network-list
		GetNetworkListDescription(ctx context.Context, params GetNetworkListDescriptionRequest) (*GetNetworkListDescriptionResponse, error)

		// UpdateNetworkListDescription modifies network list description.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/put-network-list-details
		UpdateNetworkListDescription(ctx context.Context, params UpdateNetworkListDescriptionRequest) (*UpdateNetworkListDescriptionResponse, error)

		// NetworkListSubscription

		// GetNetworkListSubscription retrieves networklist subscription.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/post-notifications-subscribe
		GetNetworkListSubscription(ctx context.Context, params GetNetworkListSubscriptionRequest) (*GetNetworkListSubscriptionResponse, error)

		// UpdateNetworkListSubscription updates networklist subscription.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/post-notifications-subscribe
		UpdateNetworkListSubscription(ctx context.Context, params UpdateNetworkListSubscriptionRequest) (*UpdateNetworkListSubscriptionResponse, error)

		// RemoveNetworkListSubscription unsubscribes networklist.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/post-notifications-unsubscribe
		RemoveNetworkListSubscription(ctx context.Context, params RemoveNetworkListSubscriptionRequest) (*RemoveNetworkListSubscriptionResponse, error)
	}

	networklists struct {
		session.Session
		usePrefixes bool
	}

	// Option defines a networklist option
	Option func(*networklists)

	// ClientFunc is a networklist client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) NetworkList
)

// Client returns a new networklist Client instance with the specified controller
func Client(sess session.Session, opts ...Option) NetworkList {
	p := &networklists{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
