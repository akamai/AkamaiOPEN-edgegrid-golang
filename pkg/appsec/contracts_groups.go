package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// The ContractsGroups interface supports listing the contracts and groups for the current
	// account. Each object contains the contract, groups associated with the contract, and whether
	// Kona Site Defender or Web Application Protector is the product for that contract.
	ContractsGroups interface {
		// GetContractsGroups lists the contracts and groups for your account.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-contracts-groups
		GetContractsGroups(ctx context.Context, params GetContractsGroupsRequest) (*GetContractsGroupsResponse, error)
	}

	// GetContractsGroupsRequest is used to retrieve the list of contracts and groups for your account.
	GetContractsGroupsRequest struct {
		ConfigID   int    `json:"-"`
		Version    int    `json:"-"`
		PolicyID   string `json:"-"`
		ContractID string `json:"-"`
		GroupID    int    `json:"-"`
	}

	// GetContractsGroupsResponse is returned from a call to GetContractsGroups.
	GetContractsGroupsResponse struct {
		ContractGroups []struct {
			ContractID  string `json:"contractId"`
			DisplayName string `json:"displayName"`
			GroupID     int    `json:"groupId"`
		} `json:"contract_groups"`
	}
)

func (p *appsec) GetContractsGroups(ctx context.Context, params GetContractsGroupsRequest) (*GetContractsGroupsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetContractsGroups")

	uri :=
		"/appsec/v1/contracts-groups"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetContractsGroups request: %w", err)
	}

	var result GetContractsGroupsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get contracts groups request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.GroupID != 0 {
		var filteredResult GetContractsGroupsResponse
		for _, val := range result.ContractGroups {
			if val.ContractID == params.ContractID && val.GroupID == params.GroupID {
				filteredResult.ContractGroups = append(filteredResult.ContractGroups, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}
