package papi

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// Group  represents a property group resource
	Group struct {
		GroupID       string   `json:"groupId"`
		GroupName     string   `json:"groupName"`
		ParentGroupID string   `json:"parentGroupId,omitempty"`
		ContractIDs   []string `json:"contractIds"`
	}

	// GroupItems represents sub-compent of the group response
	GroupItems struct {
		Items []*Group `json:"items"`
	}

	// GetGroupsResponse represents a collection of groups
	// This is the reponse to the /papi/v1/groups request
	GetGroupsResponse struct {
		AccountID   string     `json:"accountId"`
		AccountName string     `json:"accountName"`
		Groups      GroupItems `json:"groups"`
	}
)

func (p *papi) GetGroups(ctx context.Context) (GetGroupsResponse, error) {
	var groups GetGroupsResponse

	p.Log(ctx).Debug("GetGroups")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/papi/v1/groups", nil)
	if err != nil {
		return groups, fmt.Errorf("failed to create getgroups request: %w", err)
	}

	resp, err := p.Exec(req, &groups)
	if err != nil {
		return groups, fmt.Errorf("getgroups request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return groups, fmt.Errorf("getgroups request failed with status code: %d", resp.StatusCode)
	}

	return groups, nil
}
