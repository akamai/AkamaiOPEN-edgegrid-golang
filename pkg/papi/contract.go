package papi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spf13/cast"
)

type (
	// Contract represents a property contract resource
	Contract struct {
		ContractID       string `json:"contractId"`
		ContractTypeName string `json:"contractTypeName"`
	}

	// GetContractResponse represents a collection of property manager contracts
	// This is the reponse to the /papi/v1/contracts request
	GetContractResponse struct {
		AccountID string `json:"accountId"`

		Contracts struct {
			Items []*Contract `json:"items"`
		} `json:"contracts"`
	}
)

func (p *papi) GetContracts(ctx context.Context) (*GetContractResponse, error) {
	var contracts GetContractResponse

	p.Log(ctx).Debug("GetContracts")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/papi/v1/contracts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcontracts request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(UsePrefixes))

	resp, err := p.Exec(req, &contracts)
	if err != nil {
		return nil, fmt.Errorf("getcontracts request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getcontracts request failed with status code: %d", resp.StatusCode)
	}

	return &contracts, nil
}
