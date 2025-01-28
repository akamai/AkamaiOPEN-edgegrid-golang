package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// Contract represents a property contract resource
	Contract struct {
		ContractID       string `json:"contractId"`
		ContractTypeName string `json:"contractTypeName"`
	}

	// ContractsItems is the response items array
	ContractsItems struct {
		Items []*Contract `json:"items"`
	}

	// GetContractsResponse is the response to the /papi/v1/contracts request
	GetContractsResponse struct {
		AccountID string         `json:"accountId"`
		Contracts ContractsItems `json:"contracts"`
	}
)

var (
	// ErrGetContracts represents error when fetching contracts fails
	ErrGetContracts = errors.New("fetching contracts")
)

func (p *papi) GetContracts(ctx context.Context) (*GetContractsResponse, error) {
	var contracts GetContractsResponse

	logger := p.Log(ctx)
	logger.Debug("GetContracts")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/papi/v1/contracts", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetContracts, err)
	}

	resp, err := p.Exec(req, &contracts)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetContracts, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetContracts, p.Error(resp))
	}

	return &contracts, nil
}
