package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
)

type (
	// Groups is the IAM group API interface
	Groups interface {
		// CreateGroup creates a new group within a parent group_id specified in the request
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/post-group
		CreateGroup(context.Context, GroupRequest) (*Group, error)

		// GetGroup returns a group's details
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-group
		GetGroup(context.Context, GetGroupRequest) (*Group, error)

		// ListAffectedUsers lists users who are affected when a group is moved
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-move-affected-users
		ListAffectedUsers(context.Context, ListAffectedUsersRequest) ([]GroupUser, error)

		// ListGroups lists all groups in which you have a scope of admin for the current account and contract type
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-groups
		ListGroups(context.Context, ListGroupsRequest) ([]Group, error)

		// RemoveGroup removes a group based on group_id. We can only delete a sub-group, and only if that sub-group doesn't include any users
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/delete-group
		RemoveGroup(context.Context, RemoveGroupRequest) error

		// UpdateGroupName changes the name of the group
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/put-group
		UpdateGroupName(context.Context, GroupRequest) (*Group, error)
	}

	// GetGroupRequest describes the request parameters of the get group endpoint
	GetGroupRequest struct {
		GroupID int64
		Actions bool
	}

	// Group describes the response of the list groups endpoint
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

	// GroupActions encapsulates permissions available to the user for this group
	GroupActions struct {
		Delete bool `json:"delete"`
		Edit   bool `json:"edit"`
	}

	// GroupUser describes the response of the list affected users endpoint
	GroupUser struct {
		AccountID     string `json:"accountId"`
		Email         string `json:"email"`
		FirstName     string `json:"firstName"`
		IdentityID    string `json:"uiIdentityId"`
		LastLoginDate string `json:"lastLoginDate"`
		LastName      string `json:"lastName"`
		UserName      string `json:"uiUserName"`
	}

	// GroupRequest describes the request and body parameters for creating new group or updating a group name endpoint
	GroupRequest struct {
		GroupID   int64  `json:"-"`
		GroupName string `json:"groupName"`
	}

	// ListAffectedUsersRequest describes the request and body parameters of the list affected users endpoint
	ListAffectedUsersRequest struct {
		DestinationGroupID int64
		SourceGroupID      int64
		UserType           string
	}

	// ListGroupsRequest describes the request parameters of the list groups endpoint
	ListGroupsRequest struct {
		Actions bool
	}

	// RemoveGroupRequest describes the request parameter for removing a group
	RemoveGroupRequest struct {
		GroupID int64
	}
)

const (
	// LostAccessUsers with a userType of lostAccess lose their access to the source group
	LostAccessUsers = "lostAccess"
	// GainAccessUsers with a userType of gainAccess gain their access to the source group
	GainAccessUsers = "gainAccess"
)

var (
	// ErrCreateGroup is returned when CreateGroup fails
	ErrCreateGroup = errors.New("create group")
	// ErrGetGroup is returned when GetGroup fails
	ErrGetGroup = errors.New("get group")
	// ErrListAffectedUsers is returned when ListAffectedUsers fails
	ErrListAffectedUsers = errors.New("list affected users")
	// ErrListGroups is returned when ListGroups fails
	ErrListGroups = errors.New("list groups")
	// ErrUpdateGroupName is returned when UpdateGroupName fails
	ErrUpdateGroupName = errors.New("update group name")
	// ErrRemoveGroup is returned when RemoveGroup fails
	ErrRemoveGroup = errors.New("remove group")
)

// Validate validates GetGroupRequest
func (r GetGroupRequest) Validate() error {
	return validation.Errors{
		"groupID": validation.Validate(r.GroupID, validation.Required),
	}.Filter()
}

// Validate validates GroupRequest
func (r GroupRequest) Validate() error {
	return validation.Errors{
		"groupID":   validation.Validate(r.GroupID, validation.Required),
		"groupName": validation.Validate(r.GroupName, validation.Required),
	}.Filter()
}

// Validate validates ListAffectedUsersRequest
func (r ListAffectedUsersRequest) Validate() error {
	return validation.Errors{
		"destinationGroupID": validation.Validate(r.DestinationGroupID, validation.Required),
		"sourceGroupID":      validation.Validate(r.SourceGroupID, validation.Required),
		"userType":           validation.Validate(r.UserType, validation.In(LostAccessUsers, GainAccessUsers)),
	}.Filter()
}

// Validate validates RemoveGroupRequest
func (r RemoveGroupRequest) Validate() error {
	return validation.Errors{
		"groupID": validation.Validate(r.GroupID, validation.Required),
	}.Filter()
}

func (i *iam) CreateGroup(ctx context.Context, params GroupRequest) (*Group, error) {
	logger := i.Log(ctx)
	logger.Debug("CreateGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateGroup, ErrStructValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "groups", strconv.FormatInt(params.GroupID, 10)))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateGroup, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateGroup, err)
	}

	var result Group
	resp, err := i.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateGroup, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateGroup, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) GetGroup(ctx context.Context, params GetGroupRequest) (*Group, error) {
	logger := i.Log(ctx)
	logger.Debug("GetGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetGroup, ErrStructValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "groups", strconv.FormatInt(params.GroupID, 10)))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetGroup, err)
	}
	q := u.Query()
	q.Add("actions", cast.ToString(params.Actions))

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetGroup, err)
	}

	var result Group
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetGroup, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetGroup, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) ListAffectedUsers(ctx context.Context, params ListAffectedUsersRequest) ([]GroupUser, error) {
	logger := i.Log(ctx)
	logger.Debug("ListAffectedUsers")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListAffectedUsers, ErrStructValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "groups", "move", strconv.FormatInt(params.SourceGroupID, 10),
		strconv.FormatInt(params.DestinationGroupID, 10), "affected-users"))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAffectedUsers, err)
	}

	if params.UserType != "" {
		q := u.Query()
		q.Add("userType", cast.ToString(params.UserType))
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAffectedUsers, err)
	}

	var result []GroupUser
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAffectedUsers, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAffectedUsers, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListGroups(ctx context.Context, params ListGroupsRequest) ([]Group, error) {
	logger := i.Log(ctx)
	logger.Debug("ListGroups")

	u, err := url.Parse(path.Join(UserAdminEP, "groups"))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListGroups, err)
	}
	q := u.Query()
	q.Add("actions", cast.ToString(params.Actions))

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListGroups, err)
	}

	var result []Group
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListGroups, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListGroups, i.Error(resp))
	}

	return result, nil
}

func (i *iam) RemoveGroup(ctx context.Context, params RemoveGroupRequest) error {
	logger := i.Log(ctx)
	logger.Debug("RemoveGroup")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrRemoveGroup, ErrStructValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "groups", strconv.FormatInt(params.GroupID, 10)))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrRemoveGroup, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrRemoveGroup, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrRemoveGroup, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrRemoveGroup, i.Error(resp))
	}

	return nil
}

func (i *iam) UpdateGroupName(ctx context.Context, params GroupRequest) (*Group, error) {
	logger := i.Log(ctx)
	logger.Debug("UpdateGroupName")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateGroupName, ErrStructValidation, err)
	}

	u, err := url.Parse(path.Join(UserAdminEP, "groups", strconv.FormatInt(params.GroupID, 10)))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateGroupName, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateGroupName, err)
	}

	var result Group
	resp, err := i.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateGroupName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateGroupName, i.Error(resp))
	}

	return &result, nil
}
