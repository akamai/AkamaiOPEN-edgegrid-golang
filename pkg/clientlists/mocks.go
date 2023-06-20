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
