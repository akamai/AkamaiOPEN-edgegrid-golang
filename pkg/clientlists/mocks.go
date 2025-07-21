//revive:disable:exported

package clientlists

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Mock is ClientList API Mock
type Mock struct {
	mock.Mock
}

var _ ClientLists = &Mock{}

// GetClientLists return list of client lists
func (p *Mock) GetClientLists(ctx context.Context, params GetClientListsRequest) (*GetClientListsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetClientListsResponse), args.Error(1)
}

func (p *Mock) GetClientList(ctx context.Context, params GetClientListRequest) (*GetClientListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetClientListResponse), args.Error(1)
}

func (p *Mock) CreateClientList(ctx context.Context, params CreateClientListRequest) (*CreateClientListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateClientListResponse), args.Error(1)
}

func (p *Mock) UpdateClientList(ctx context.Context, params UpdateClientListRequest) (*UpdateClientListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateClientListResponse), args.Error(1)
}

func (p *Mock) UpdateClientListItems(ctx context.Context, params UpdateClientListItemsRequest) (*UpdateClientListItemsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateClientListItemsResponse), args.Error(1)
}

func (p *Mock) DeleteClientList(ctx context.Context, params DeleteClientListRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}

func (p *Mock) GetActivation(ctx context.Context, params GetActivationRequest) (*GetActivationResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetActivationResponse), args.Error(1)
}

func (p *Mock) GetActivationStatus(ctx context.Context, params GetActivationStatusRequest) (*GetActivationStatusResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetActivationStatusResponse), args.Error(1)
}

func (p *Mock) CreateActivation(ctx context.Context, params CreateActivationRequest) (*CreateActivationResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateActivationResponse), args.Error(1)
}

func (p *Mock) CreateDeactivation(ctx context.Context, params CreateDeactivationRequest) (*CreateDeactivationResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateDeactivationResponse), args.Error(1)
}
