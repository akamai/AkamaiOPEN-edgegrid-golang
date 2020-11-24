package iam

import "context"

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
