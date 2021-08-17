package appsec

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// The ContractsGroups interface supports listing the contracts and groups for the current
	// account. Each object contains the contract, groups associated with the contract, and whether
	// Kona Site Defender or Web Application Protector is the product for that contract.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#contractgroup
	ContractsGroups interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcontractsandgroupswithksdorwaf
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

	var rval GetContractsGroupsResponse
	var rvalfiltered GetContractsGroupsResponse

	uri :=
		"/appsec/v1/contracts-groups"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetContractsGroups request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetContractsGroups request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.GroupID != 0 {
		for _, val := range rval.ContractGroups {
			if val.ContractID == params.ContractID && val.GroupID == params.GroupID {
				rvalfiltered.ContractGroups = append(rvalfiltered.ContractGroups, val)
			}
		}
	} else {
		rvalfiltered = rval
	}
	return &rvalfiltered, nil

}
