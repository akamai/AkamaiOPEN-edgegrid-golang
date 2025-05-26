package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListPropertiesRequest holds parameters for ListProperties
	ListPropertiesRequest struct {
		Unused      bool
		ContractIDs []string
	}

	// ListOriginsRequest holds parameters for ListOrigins
	ListOriginsRequest struct {
		PropertyID int64
		ContractID string
		GroupID    int64
	}

	// ListPropertiesResponse contains response from ListProperties
	ListPropertiesResponse struct {
		Properties []Property `json:"properties"`
	}

	// ListOriginsResponse contains response from ListOrigins
	ListOriginsResponse struct {
		Children []Child    `json:"children"`
		Default  []Behavior `json:"default"`
	}

	// Child represents children rules in a property
	Child struct {
		Name      string     `json:"name"`
		Behaviors []Behavior `json:"behaviors"`
	}

	// Behavior contains behavior information
	Behavior struct {
		Hostname   string     `json:"hostname"`
		OriginType OriginType `json:"originType"`
	}

	// Property represents property object
	Property struct {
		GroupID      int64        `json:"groupId"`
		ContractID   string       `json:"contractId"`
		PropertyID   int64        `json:"propertyId"`
		PropertyName string       `json:"propertyName"`
		Type         PropertyType `json:"type"`
	}

	// OriginType represents the type of the origin
	OriginType string

	// PropertyType represents the type of the property
	PropertyType string
)

const (
	// PropertyTypeWeb is the web type of the property
	PropertyTypeWeb PropertyType = "WEB"
	// PropertyTypeMedia is the media type of the property
	PropertyTypeMedia PropertyType = "MEDIA"
	// OriginTypeCustomer is the customer type of the origin
	OriginTypeCustomer OriginType = "CUSTOMER"
	// OriginTypeNetStorage is the net storage type of the origin
	OriginTypeNetStorage OriginType = "NET_STORAGE"
)

// Validate validates ListOriginsRequest
func (r ListOriginsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"ContractID": validation.Validate(r.ContractID, validation.Required),
		"GroupID":    validation.Validate(r.GroupID, validation.Required),
	})
}

var (
	// ErrListProperties is returned when ListProperties fails
	ErrListProperties = errors.New("list properties")
	// ErrListOrigins is returned when ListOrigins fails
	ErrListOrigins = errors.New("list origins")
)

func (c *cloudwrapper) ListProperties(ctx context.Context, params ListPropertiesRequest) (*ListPropertiesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListProperties")

	uri, err := url.Parse("/cloud-wrapper/v1/properties")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListProperties, err)
	}

	q := uri.Query()
	q.Add("unused", strconv.FormatBool(params.Unused))
	for _, ctr := range params.ContractIDs {
		q.Add("contractIds", ctr)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListProperties, err)
	}

	var result ListPropertiesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListProperties, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudwrapper) ListOrigins(ctx context.Context, params ListOriginsRequest) (*ListOriginsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListOrigins")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListOrigins, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloud-wrapper/v1/properties/%d/origins", params.PropertyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListOrigins, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", strconv.FormatInt(params.GroupID, 10))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListOrigins, err)
	}

	var result ListOriginsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListOrigins, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListOrigins, c.Error(resp))
	}

	return &result, nil
}
