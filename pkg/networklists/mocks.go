//revive:disable:exported

package networklists

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ NetworkList = &Mock{}

func (p *Mock) CreateActivations(ctx context.Context, params CreateActivationsRequest) (*CreateActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateActivationsResponse), args.Error(1)
}

func (p *Mock) GetActivations(ctx context.Context, params GetActivationsRequest) (*GetActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetActivationsResponse), args.Error(1)
}

func (p *Mock) GetActivation(ctx context.Context, params GetActivationRequest) (*GetActivationResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetActivationResponse), args.Error(1)
}

func (p *Mock) RemoveActivations(ctx context.Context, params RemoveActivationsRequest) (*RemoveActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RemoveActivationsResponse), args.Error(1)
}

func (p *Mock) CreateNetworkList(ctx context.Context, params CreateNetworkListRequest) (*CreateNetworkListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateNetworkListResponse), args.Error(1)
}

func (p *Mock) RemoveNetworkList(ctx context.Context, params RemoveNetworkListRequest) (*RemoveNetworkListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RemoveNetworkListResponse), args.Error(1)
}

func (p *Mock) UpdateNetworkList(ctx context.Context, params UpdateNetworkListRequest) (*UpdateNetworkListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateNetworkListResponse), args.Error(1)
}

func (p *Mock) GetNetworkList(ctx context.Context, params GetNetworkListRequest) (*GetNetworkListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetNetworkListResponse), args.Error(1)
}

func (p *Mock) GetNetworkLists(ctx context.Context, params GetNetworkListsRequest) (*GetNetworkListsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetNetworkListsResponse), args.Error(1)
}

func (p *Mock) GetNetworkListDescription(ctx context.Context, params GetNetworkListDescriptionRequest) (*GetNetworkListDescriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetNetworkListDescriptionResponse), args.Error(1)
}

func (p *Mock) UpdateNetworkListDescription(ctx context.Context, params UpdateNetworkListDescriptionRequest) (*UpdateNetworkListDescriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateNetworkListDescriptionResponse), args.Error(1)
}

func (p *Mock) GetNetworkListSubscription(ctx context.Context, params GetNetworkListSubscriptionRequest) (*GetNetworkListSubscriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetNetworkListSubscriptionResponse), args.Error(1)
}

func (p *Mock) RemoveNetworkListSubscription(ctx context.Context, params RemoveNetworkListSubscriptionRequest) (*RemoveNetworkListSubscriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RemoveNetworkListSubscriptionResponse), args.Error(1)
}

func (p *Mock) UpdateNetworkListSubscription(ctx context.Context, params UpdateNetworkListSubscriptionRequest) (*UpdateNetworkListSubscriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateNetworkListSubscriptionResponse), args.Error(1)
}
