package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// Groups contains operations available on Group resource
	Groups interface {
		// GetGroups provides a read-only list of groups, which may contain properties.
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-groups
		GetGroups(context.Context) (*GetGroupsResponse, error)
	}

	// Group represents a property group resource
	Group struct {
		GroupID       string   `json:"groupId"`
		GroupName     string   `json:"groupName"`
		ParentGroupID string   `json:"parentGroupId,omitempty"`
		ContractIDs   []string `json:"contractIds"`
	}

	// GroupItems represents sub-component of the group response
	GroupItems struct {
		Items []*Group `json:"items"`
	}

	// GetGroupsResponse represents a collection of groups
	// This is the response to the /papi/v1/groups request
	GetGroupsResponse struct {
		AccountID   string     `json:"accountId"`
		AccountName string     `json:"accountName"`
		Groups      GroupItems `json:"groups"`
	}
)

var (
	// ErrGetGroups represents error when fetching groups fails
	ErrGetGroups = errors.New("fetching groups")
)

func (p *papi) GetGroups(ctx context.Context) (*GetGroupsResponse, error) {
	var groups GetGroupsResponse

	logger := p.Log(ctx)
	logger.Debug("GetGroups")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/papi/v1/groups", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetGroups, err)
	}

	resp, err := p.Exec(req, &groups)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetGroups, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetGroups, p.Error(resp))
	}

	return &groups, nil
}
