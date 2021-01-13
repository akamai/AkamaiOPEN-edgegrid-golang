package iam

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/spf13/cast"
)

type (
	// Groups is the IAM group management interface
	Groups interface {
		ListGroups(context.Context, ListGroupsRequest) ([]Group, error)
	}

	// ListGroupsRequest is the request for listing groups
	ListGroupsRequest struct {
		Actions bool `json:"actions"`
	}

	// GroupActions encapsulates permissions available to the user for this group.
	GroupActions struct {
		Delete bool `json:"bool"`
		Edit   bool `json:"edit"`
	}

	// Group encapsulates information about a group.
	Group struct {
		Actions       *GroupActions `json:"actions,omitempty"`
		CreatedBy     string        `json:"createdBy"`
		CreatedDate   string        `json:"createdDate"`
		GroupID       int64         `json:"groupId"`
		GroupName     string        `json:"groupName"`
		ModifiedBy    string        `json:"modifiedBy"`
		ModifiedDate  string        `json:"modifiedDate"`
		ParentGroupID int64         `json:"parentGroupId"`
		SubGroups     []Group       `json:"subGroups,omitempty"`
	}
)

func (i *iam) ListGroups(ctx context.Context, params ListGroupsRequest) ([]Group, error) {
	logger := i.Log(ctx)
	logger.Debug("ListGroups")

	u, err := url.Parse(path.Join(UserAdminEP, "groups"))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "ListGroups", err)
	}
	q := u.Query()
	q.Add("actions", cast.ToString(params.Actions))

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "ListGroups", err)
	}

	var rval []Group
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "ListGroups", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "ListGroups", i.Error(resp))
	}

	return rval, nil
}
