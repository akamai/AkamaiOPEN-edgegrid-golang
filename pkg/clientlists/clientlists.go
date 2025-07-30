// Package clientlists provides access to Akamai Client Lists APIs
//
// See: https://techdocs.akamai.com/client-lists/reference/api
package clientlists

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
	// ClientLists is the API interface for Client Lists
	ClientLists interface {
		// Activations

		// GetActivation retrieves details of a specified activation ID.
		//
		// See: https://techdocs.akamai.com/client-lists/reference/get-retrieve-activation-status
		GetActivation(ctx context.Context, params GetActivationRequest) (*GetActivationResponse, error)

		// GetActivationStatus retrieves activation status for a client list in a network environment.
		//
		// See: https://techdocs.akamai.com/client-lists/reference/get-activation-status
		GetActivationStatus(ctx context.Context, params GetActivationStatusRequest) (*GetActivationStatusResponse, error)

		// CreateActivation activates a client list
		//
		// See: https://techdocs.akamai.com/client-lists/reference/post-activate-list
		CreateActivation(ctx context.Context, params CreateActivationRequest) (*CreateActivationResponse, error)

		// CreateDeactivation deactivates a client list
		//
		// See: https://techdocs.akamai.com/client-lists/reference/post-activate-list
		CreateDeactivation(ctx context.Context, params CreateDeactivationRequest) (*CreateDeactivationResponse, error)

		// Lists

		// GetClientLists lists all client lists accessible for an authenticated user
		//
		// See: https://techdocs.akamai.com/client-lists/reference/get-lists
		GetClientLists(ctx context.Context, params GetClientListsRequest) (*GetClientListsResponse, error)

		// GetClientList retrieves client list with specific list id
		//
		// See: https://techdocs.akamai.com/client-lists/reference/get-list
		GetClientList(ctx context.Context, params GetClientListRequest) (*GetClientListResponse, error)

		// CreateClientList creates a new client list
		//
		// See: https://techdocs.akamai.com/client-lists/reference/post-create-list
		CreateClientList(ctx context.Context, params CreateClientListRequest) (*CreateClientListResponse, error)

		// UpdateClientList updates existing client list
		//
		// See: https://techdocs.akamai.com/client-lists/reference/put-update-list
		UpdateClientList(ctx context.Context, params UpdateClientListRequest) (*UpdateClientListResponse, error)

		// UpdateClientListItems updates items/entries of an existing client lists
		//
		// See: https://techdocs.akamai.com/client-lists/reference/post-update-items
		UpdateClientListItems(ctx context.Context, params UpdateClientListItemsRequest) (*UpdateClientListItemsResponse, error)

		// DeleteClientList removes a client list
		//
		// See: https://techdocs.akamai.com/client-lists/reference/delete-list
		DeleteClientList(ctx context.Context, params DeleteClientListRequest) error
	}

	clientlists struct {
		session.Session
	}

	// Option defines a clientlists option
	Option func(*clientlists)

	// ClientFunc is a clientlists client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) ClientLists
)

// Client returns a new clientlists Client instance with the specified controller
func Client(sess session.Session, opts ...Option) ClientLists {
	p := &clientlists{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
