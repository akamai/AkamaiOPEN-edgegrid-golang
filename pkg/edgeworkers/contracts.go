package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// Contracts is an edgeworkers contracts API interface
	Contracts interface {
		// ListContracts lists contract IDs that can be used to list resource tiers
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-contracts-1
		ListContracts(context.Context) (*ListContractsResponse, error)
	}

	// ListContractsResponse represents a response object returned by ListContracts
	ListContractsResponse struct {
		ContractIDs []string `json:"contractIds"`
	}
)

var (
	// ErrListContracts is returned in case an error occurs on ListContracts operation
	ErrListContracts = errors.New("list contracts")
)

func (e *edgeworkers) ListContracts(ctx context.Context) (*ListContractsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListContracts")

	uri := "/edgeworkers/v1/contracts"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListContracts, err)
	}

	var result ListContractsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListContracts, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListContracts, e.Error(resp))
	}

	return &result, nil
}
