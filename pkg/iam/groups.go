package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetGroupRequest describes the request parameters for the GetGroup endpoint.
	GetGroupRequest struct {
		GroupID int64
		Actions bool
	}

	// Group describes the response from the ListGroups endpoint.
	Group struct {
		Actions       *GroupActions `json:"actions,omitempty"`
		CreatedBy     string        `json:"createdBy"`
		CreatedDate   time.Time     `json:"createdDate"`
		GroupID       int64         `json:"groupId"`
		GroupName     string        `json:"groupName"`
		ModifiedBy    string        `json:"modifiedBy"`
		ModifiedDate  time.Time     `json:"modifiedDate"`
		ParentGroupID int64         `json:"parentGroupId"`
		SubGroups     []Group       `json:"subGroups,omitempty"`
	}

	// GroupActions encapsulates permissions available to the user for this group.
	GroupActions struct {
		Delete bool `json:"delete"`
		Edit   bool `json:"edit"`
	}

	// GroupUser describes the response from the ListAffectedUsers endpoint.
	GroupUser struct {
		AccountID     string    `json:"accountId"`
		Email         string    `json:"email"`
		FirstName     string    `json:"firstName"`
		IdentityID    string    `json:"uiIdentityId"`
		LastLoginDate time.Time `json:"lastLoginDate"`
		LastName      string    `json:"lastName"`
		UserName      string    `json:"uiUserName"`
	}

	// GroupRequest describes the request and body parameters for creating new group or updating a group name endpoint.
	GroupRequest struct {
		GroupID   int64  `json:"-"`
		GroupName string `json:"groupName"`
	}

	// MoveGroupRequest describes the request body for the MoveGroup endpoint.
	MoveGroupRequest struct {
		SourceGroupID      int64 `json:"sourceGroupId"`
		DestinationGroupID int64 `json:"destinationGroupId"`
	}

	// ListAffectedUsersRequest describes the request and body parameters of the ListAffectedUsers endpoint.
	ListAffectedUsersRequest struct {
		DestinationGroupID int64
		SourceGroupID      int64
		UserType           string
	}

	// ListGroupsRequest describes the request parameters of the ListGroups endpoint.
	ListGroupsRequest struct {
		Actions bool
	}

	// RemoveGroupRequest describes the request parameter for the RemoveGroup endpoint.
	RemoveGroupRequest struct {
		GroupID int64
	}
)

const (
	// LostAccessUsers with a userType of lostAccess lose their access to the source group.
	LostAccessUsers = "lostAccess"
	// GainAccessUsers with a userType of gainAccess gain their access to the source group.
	GainAccessUsers = "gainAccess"
)

var (
	// ErrCreateGroup is returned when CreateGroup fails.
	ErrCreateGroup = errors.New("create group")
	// ErrGetGroup is returned when GetGroup fails.
	ErrGetGroup = errors.New("get group")
	// ErrListAffectedUsers is returned when ListAffectedUsers fails.
	ErrListAffectedUsers = errors.New("list affected users")
	// ErrListGroups is returned when ListGroups fails.
	ErrListGroups = errors.New("list groups")
	// ErrUpdateGroupName is returned when UpdateGroupName fails.
	ErrUpdateGroupName = errors.New("update group name")
	// ErrRemoveGroup is returned when RemoveGroup fails.
	ErrRemoveGroup = errors.New("remove group")
	// ErrMoveGroup is returned when MoveGroup fails.
	ErrMoveGroup = errors.New("move group")
)

// Validate validates GetGroupRequest.
func (r GetGroupRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"GroupID": validation.Validate(r.GroupID, validation.Required),
	})
}

// Validate validates GroupRequest.
func (r GroupRequest) Validate() error {
	return validation.Errors{
		"GroupID":   validation.Validate(r.GroupID, validation.Required),
		"GroupName": validation.Validate(r.GroupName, validation.Required),
	}.Filter()
}

// Validate validates MoveGroupRequest.
func (r MoveGroupRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DestinationGroupID": validation.Validate(r.DestinationGroupID, validation.Required),
		"SourceGroupID":      validation.Validate(r.SourceGroupID, validation.Required),
	})
}

// Validate validates ListAffectedUsersRequest.
func (r ListAffectedUsersRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DestinationGroupID": validation.Validate(r.DestinationGroupID, validation.Required),
		"SourceGroupID":      validation.Validate(r.SourceGroupID, validation.Required),
		"UserType":           validation.Validate(r.UserType, validation.In(LostAccessUsers, GainAccessUsers)),
	})
}

// Validate validates RemoveGroupRequest.
func (r RemoveGroupRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"GroupID": validation.Validate(r.GroupID, validation.Required),
	})
}

func (i *iam) CreateGroup(ctx context.Context, params GroupRequest) (*Group, error) {
	logger := i.Log(ctx)
	logger.Debug("CreateGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateGroup, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/groups/%d", params.GroupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateGroup, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateGroup, err)
	}

	var result Group
	resp, err := i.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateGroup, err)
	}
	defer session.CloseResponseBody(resp)

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

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/groups/%d", params.GroupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetGroup, err)
	}
	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetGroup, err)
	}

	var result Group
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetGroup, err)
	}
	defer session.CloseResponseBody(resp)

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

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/groups/move/%d/%d/affected-users", params.SourceGroupID, params.DestinationGroupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAffectedUsers, err)
	}

	if params.UserType != "" {
		q := uri.Query()
		q.Add("userType", params.UserType)
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAffectedUsers, err)
	}

	var result []GroupUser
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAffectedUsers, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAffectedUsers, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListGroups(ctx context.Context, params ListGroupsRequest) ([]Group, error) {
	logger := i.Log(ctx)
	logger.Debug("ListGroups")

	uri, err := url.Parse("/identity-management/v3/user-admin/groups")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListGroups, err)
	}
	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListGroups, err)
	}

	var result []Group
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListGroups, err)
	}
	defer session.CloseResponseBody(resp)

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

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/groups/%d", params.GroupID))
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrRemoveGroup, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrRemoveGroup, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrRemoveGroup, err)
	}
	defer session.CloseResponseBody(resp)

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

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/groups/%d", params.GroupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateGroupName, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateGroupName, err)
	}

	var result Group
	resp, err := i.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateGroupName, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateGroupName, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) MoveGroup(ctx context.Context, params MoveGroupRequest) error {
	logger := i.Log(ctx)
	logger.Debug("MoveGroup")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrMoveGroup, ErrStructValidation, err)
	}

	uri, err := url.Parse("/identity-management/v3/user-admin/groups/move")
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrMoveGroup, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrMoveGroup, err)
	}

	resp, err := i.Exec(req, nil, params)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrMoveGroup, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w: %s", ErrMoveGroup, i.Error(resp))
	}

	return nil
}
