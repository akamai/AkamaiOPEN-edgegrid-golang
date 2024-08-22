package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Properties is the IAM properties API interface
	Properties interface {
		// ListProperties lists the properties for the current account or other managed accounts using the accountSwitchKey parameter.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-properties
		ListProperties(context.Context, ListPropertiesRequest) (*ListPropertiesResponse, error)

		// GetProperty lists a property's details.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-property
		GetProperty(context.Context, GetPropertyRequest) (*GetPropertyResponse, error)

		// MoveProperty moves a property from one group to another group.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-property
		MoveProperty(context.Context, MovePropertyRequest) error

		// MapPropertyIDToName returns property name for given (IAM) property ID
		// Mainly to be used to map (IAM) Property ID to (PAPI) Property ID
		// To finish the mapping, please use papi.MapPropertyNameToID
		MapPropertyIDToName(context.Context, MapPropertyIDToNameRequest) (*string, error)
	}

	// ListPropertiesRequest contains the request parameters for the list properties operation.
	ListPropertiesRequest struct {
		GroupID int64
		Actions bool
	}

	// GetPropertyRequest contains the request parameters for the get property operation.
	GetPropertyRequest struct {
		PropertyID int64
		GroupID    int64
	}

	// MapPropertyNameToIDRequest is the argument for MapPropertyNameToID
	MapPropertyNameToIDRequest string

	// ListPropertiesResponse holds the response data from ListProperties.
	ListPropertiesResponse []Property

	// GetPropertyResponse holds the response data from GetProperty.
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

	// MovePropertyRequest contains the request parameters for the MoveProperty operation.
	MovePropertyRequest struct {
		PropertyID int64
		BodyParams MovePropertyReqBody
	}

	// MovePropertyReqBody contains body parameters for the MoveProperty operation.
	MovePropertyReqBody struct {
		DestinationGroupID int64 `json:"destinationGroupId"`
		SourceGroupID      int64 `json:"sourceGroupId"`
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

	// MapPropertyIDToNameRequest is the argument for MapPropertyIDToName
	MapPropertyIDToNameRequest struct {
		PropertyID int64
		GroupID    int64
	}
)

// Validate validates GetPropertyRequest
func (r GetPropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	})
}

// Validate validates MapPropertyIDToNameRequest
func (r MapPropertyIDToNameRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	})
}

// Validate validates MovePropertyRequest
func (r MovePropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"BodyParams": validation.Validate(r.BodyParams, validation.Required),
	})
}

// Validate validates MovePropertyReqBody
func (r MovePropertyReqBody) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DestinationGroupID": validation.Validate(r.DestinationGroupID, validation.Required),
		"SourceGroupID":      validation.Validate(r.SourceGroupID, validation.Required),
	})
}

var (
	// ErrListProperties is returned when ListProperties fails
	ErrListProperties = errors.New("list properties")
	// ErrGetProperty is returned when GetProperty fails
	ErrGetProperty = errors.New("get property")
	// ErrMoveProperty is returned when MoveProperty fails
	ErrMoveProperty = errors.New("move property")
	// ErrMapPropertyIDToName is returned when MapPropertyIDToName fails
	ErrMapPropertyIDToName = errors.New("map property by id")
	// ErrMapPropertyNameToID is returned when MapPropertyNameToID fails
	ErrMapPropertyNameToID = errors.New("map property by name")
	// ErrNoProperty is returned when MapPropertyNameToID did not find given property
	ErrNoProperty = errors.New("no such property")
)

func (i *iam) ListProperties(ctx context.Context, params ListPropertiesRequest) (*ListPropertiesResponse, error) {
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

	return &result, nil
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

	resp, err := i.Exec(req, nil, params.BodyParams)
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

	for _, property := range *properties {
		if property.PropertyName == string(name) {
			return &property.PropertyID, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrNoProperty, name)
}
