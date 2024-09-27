package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListPropertiesRequest contains the request parameters for the ListProperties endpoint.
	ListPropertiesRequest struct {
		GroupID int64
		Actions bool
	}

	// ListUsersForPropertyRequest contains the request parameters for the ListUsersForProperty endpoint.
	ListUsersForPropertyRequest struct {
		PropertyID int64
		UserType   PropertyUserType
	}

	// GetPropertyRequest contains the request parameters for the GetProperty endpoint.
	GetPropertyRequest struct {
		PropertyID int64
		GroupID    int64
	}

	// MapPropertyNameToIDRequest is the argument for MapPropertyNameToID.
	MapPropertyNameToIDRequest string

	// BlockUsersRequest contains the request parameters for the BlockUsers endpoint.
	BlockUsersRequest struct {
		PropertyID int64
		Body       BlockUsersRequestBody
	}

	// ListPropertiesResponse holds the response data from the ListProperties endpoint.
	ListPropertiesResponse []Property

	// ListUsersForPropertyResponse holds the response data from the ListUsersForProperty endpoint.
	ListUsersForPropertyResponse []UsersForProperty

	// GetPropertyResponse holds the response data from the GetProperty endpoint.
	GetPropertyResponse struct {
		ARLConfigFile string    `json:"arlConfigFile"`
		CreatedBy     string    `json:"createdBy"`
		CreatedDate   time.Time `json:"createdDate"`
		GroupID       int64     `json:"groupId"`
		GroupName     string    `json:"groupName"`
		ModifiedBy    string    `json:"modifiedBy"`
		ModifiedDate  time.Time `json:"modifiedDate"`
		PropertyID    int64     `json:"propertyId"`
		PropertyName  string    `json:"propertyName"`
	}

	// BlockUsersResponse holds the response data from the BlockUsers endpoint.
	BlockUsersResponse []UsersForProperty

	// MovePropertyRequest contains the request parameters for the MoveProperty endpoint.
	MovePropertyRequest struct {
		PropertyID int64
		Body       MovePropertyRequestBody
	}

	// MovePropertyRequestBody contains body parameters for the MoveProperty endpoint.
	MovePropertyRequestBody struct {
		DestinationGroupID int64 `json:"destinationGroupId"`
		SourceGroupID      int64 `json:"sourceGroupId"`
	}

	// BlockUsersRequestBody hold the request body parameters for the BlockUsers endpoint.
	BlockUsersRequestBody []BlockUserItem

	// BlockUserItem contains body parameters for the BlockUsers endpoint.
	BlockUserItem struct {
		UIIdentityID string `json:"uiIdentityId"`
	}

	// Property holds the property details.
	Property struct {
		PropertyID              int64           `json:"propertyId"`
		PropertyName            string          `json:"propertyName"`
		PropertyTypeDescription string          `json:"propertyTypeDescription"`
		GroupID                 int64           `json:"groupId"`
		GroupName               string          `json:"groupName"`
		Actions                 PropertyActions `json:"actions"`
	}

	// PropertyActions specifies activities available for the property.
	PropertyActions struct {
		Move bool `json:"move"`
	}

	// MapPropertyIDToNameRequest is the argument for MapPropertyIDToName.
	MapPropertyIDToNameRequest struct {
		PropertyID int64
		GroupID    int64
	}

	// UsersForProperty holds details about the users accessing the property.
	UsersForProperty struct {
		FirstName    string `json:"firstName"`
		IsBlocked    bool   `json:"isBlocked"`
		LastName     string `json:"lastName"`
		UIIdentityID string `json:"uiIdentityId"`
		UIUserName   string `json:"uiUserName"`
	}

	// PropertyUserType filters property users based on their access to the property.
	PropertyUserType string
)

const (
	// PropertyUserTypeAll selects all property users.
	PropertyUserTypeAll PropertyUserType = "all"
	// PropertyUserTypeAssigned selects users that have access to the property.
	PropertyUserTypeAssigned PropertyUserType = "assigned"
	// PropertyUserTypeBlocked selects users whose access to the property is blocked.
	PropertyUserTypeBlocked PropertyUserType = "blocked"
)

// Validate validates PropertyUserType.
func (p PropertyUserType) Validate() error {
	return validation.In(PropertyUserTypeAll, PropertyUserTypeAssigned, PropertyUserTypeBlocked).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' or '%s'",
			p, PropertyUserTypeAll, PropertyUserTypeAssigned, PropertyUserTypeBlocked)).
		Validate(p)
}

// Validate validates ListUsersForPropertyRequest.
func (r ListUsersForPropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"UserType":   r.UserType.Validate(),
	})
}

// Validate validates GetPropertyRequest.
func (r GetPropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	})
}

// Validate validates MapPropertyIDToNameRequest.
func (r MapPropertyIDToNameRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	})
}

// Validate validates MovePropertyRequest.
func (r MovePropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"Body":       validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates MovePropertyRequestBody.
func (r MovePropertyRequestBody) Validate() error {
	return validation.Errors{
		"DestinationGroupID": validation.Validate(r.DestinationGroupID, validation.Required),
		"SourceGroupID":      validation.Validate(r.SourceGroupID, validation.Required),
	}.Filter()
}

// Validate validates BlockUsersRequest.
func (r BlockUsersRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"Body":       validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates BlockUserItem.
func (r BlockUserItem) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"UIIdentityID": validation.Validate(r.UIIdentityID, validation.Required),
	})
}

var (
	// ErrListProperties is returned when ListProperties fails.
	ErrListProperties = errors.New("list properties")
	// ErrListUsersForProperty is returned when ListUsersForProperty fails.
	ErrListUsersForProperty = errors.New("list users for property")
	// ErrGetProperty is returned when GetProperty fails.
	ErrGetProperty = errors.New("get property")
	// ErrMoveProperty is returned when MoveProperty fails.
	ErrMoveProperty = errors.New("move property")
	// ErrMapPropertyIDToName is returned when MapPropertyIDToName fails.
	ErrMapPropertyIDToName = errors.New("map property by id")
	// ErrMapPropertyNameToID is returned when MapPropertyNameToID fails.
	ErrMapPropertyNameToID = errors.New("map property by name")
	// ErrNoProperty is returned when MapPropertyNameToID did not find given property.
	ErrNoProperty = errors.New("no such property")
	// ErrBlockUsers is returned when BlockUsers fails.
	ErrBlockUsers = errors.New("block users")
)

func (i *iam) ListProperties(ctx context.Context, params ListPropertiesRequest) (ListPropertiesResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListProperties")

	uri, err := url.Parse("/identity-management/v3/user-admin/properties")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListProperties, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	if params.GroupID != 0 {
		q.Add("groupId", strconv.FormatInt(params.GroupID, 10))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListProperties, err)
	}

	var result ListPropertiesResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListProperties, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListUsersForProperty(ctx context.Context, params ListUsersForPropertyRequest) (ListUsersForPropertyResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListUsersForProperty")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListUsersForProperty, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/properties/%d/users", params.PropertyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListUsersForProperty, err)
	}

	if params.UserType != "" {
		q := uri.Query()
		q.Add("userType", string(params.UserType))
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListUsersForProperty, err)
	}

	var result ListUsersForPropertyResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListUsersForProperty, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListUsersForProperty, i.Error(resp))
	}

	return result, nil
}

func (i *iam) GetProperty(ctx context.Context, params GetPropertyRequest) (*GetPropertyResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("GetProperty")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetProperty, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/properties/%d", params.PropertyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetProperty, err)
	}

	q := uri.Query()
	q.Add("groupId", strconv.FormatInt(params.GroupID, 10))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetProperty, err)
	}

	var result GetPropertyResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetProperty, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetProperty, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) MoveProperty(ctx context.Context, params MovePropertyRequest) error {
	logger := i.Log(ctx)
	logger.Debug("MoveProperty")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrMoveProperty, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/identity-management/v3/user-admin/properties/%d", params.PropertyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrMoveProperty, err)
	}

	resp, err := i.Exec(req, nil, params.Body)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrMoveProperty, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrMoveProperty, i.Error(resp))
	}

	return nil
}

func (i *iam) MapPropertyIDToName(ctx context.Context, params MapPropertyIDToNameRequest) (*string, error) {
	logger := i.Log(ctx)
	logger.Debug("MapPropertyIDToName")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrMapPropertyIDToName, ErrStructValidation, err)
	}

	req := GetPropertyRequest{
		PropertyID: params.PropertyID,
		GroupID:    params.GroupID,
	}

	property, err := i.GetProperty(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrMapPropertyIDToName, err)
	}

	return &property.PropertyName, nil
}

func (i *iam) MapPropertyNameToID(ctx context.Context, name MapPropertyNameToIDRequest) (*int64, error) {
	logger := i.Log(ctx)
	logger.Debug("MapPropertyNameToID")

	if name == "" {
		return nil, fmt.Errorf("%s: %w:\n name cannot be blank", ErrMapPropertyNameToID, ErrStructValidation)
	}

	properties, err := i.ListProperties(ctx, ListPropertiesRequest{})
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrMapPropertyNameToID, err)
	}

	for _, property := range properties {
		if property.PropertyName == string(name) {
			return &property.PropertyID, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrNoProperty, name)
}

func (i *iam) BlockUsers(ctx context.Context, params BlockUsersRequest) (*BlockUsersResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("BlockUsers")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrBlockUsers, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/identity-management/v3/user-admin/properties/%d/users/block", params.PropertyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrBlockUsers, err)
	}

	var result BlockUsersResponse
	resp, err := i.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrBlockUsers, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrBlockUsers, i.Error(resp))
	}

	return &result, nil
}
