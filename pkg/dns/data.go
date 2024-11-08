package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

type (
	// ListGroupResponse lists the groups accessible to the current user
	ListGroupResponse struct {
		Groups []Group `json:"groups"`
	}

	// ListGroupRequest is a request struct
	ListGroupRequest struct {
		GroupID string
	}

	// Group contain the information of the particular group
	Group struct {
		GroupID     int      `json:"groupId"`
		GroupName   string   `json:"groupName"`
		ContractIDs []string `json:"contractIds"`
		Permissions []string `json:"permissions"`
	}
)

var (
	// ErrListGroups is returned in case an error occurs on ListGroups operation
	ErrListGroups = errors.New("list groups")
)

func (d *dns) ListGroups(ctx context.Context, params ListGroupRequest) (*ListGroupResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("ListGroups")

	uri, err := url.Parse("/config-dns/v2/data/groups/")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListGroups, err)
	}

	q := uri.Query()
	if params.GroupID != "" {
		q.Add("gid", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list listZoneGroups request: %w", err)
	}

	var result ListGroupResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListZoneGroups request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}
