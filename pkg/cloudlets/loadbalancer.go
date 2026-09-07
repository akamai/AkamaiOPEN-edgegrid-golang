package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// OriginResponse is an Origin returned in ListOrigins
	OriginResponse struct {
		Hostname string `json:"hostname"`
		Origin
	}

	// OriginType is a type for Origin Type
	OriginType string

	// ListOriginsRequest describes the parameters of the ListOrigins request
	ListOriginsRequest struct {
		Type OriginType
	}

	// GetOriginRequest describes the parameters of the get origins request
	GetOriginRequest struct {
		OriginID string
	}

	// CreateOriginRequest describes the parameters of the create origin request
	CreateOriginRequest struct {
		OriginID string `json:"originId"`
		Description
	}

	// UpdateOriginRequest describes the parameters of the update origin request
	UpdateOriginRequest struct {
		OriginID string
		Description
	}

	// Description describes description for the Origin
	Description struct {
		Description string `json:"description,omitempty"`
	}

	// Origin is a response returned by CreateOrigin
	Origin struct {
		OriginID    string     `json:"originId"`
		Description string     `json:"description"`
		Akamaized   bool       `json:"akamaized"`
		Checksum    string     `json:"checksum"`
		Type        OriginType `json:"type"`
	}
)

const (
	// OriginTypeAll is a value to use when you want ListOrigins to return origins of all types
	OriginTypeAll OriginType = ""
	// OriginTypeCustomer is a value to use when you want ListOrigins to return only origins of CUSTOMER type
	OriginTypeCustomer OriginType = "CUSTOMER"
	// OriginTypeApplicationLoadBalancer is a value to use when you want ListOrigins to return only origins of APPLICATION_LOAD_BALANCER type
	OriginTypeApplicationLoadBalancer OriginType = "APPLICATION_LOAD_BALANCER"
	// OriginTypeNetStorage is a value to use when you want ListOrigins to return only origins of NETSTORAGE type
	OriginTypeNetStorage OriginType = "NETSTORAGE"
)

var (
	// ErrListOrigins is returned when ListOrigins fails
	ErrListOrigins = errors.New("list origins")
	// ErrGetOrigin is returned when GetOrigin fails
	ErrGetOrigin = errors.New("get origin")
	// ErrCreateOrigin is returned when CreateOrigin fails
	ErrCreateOrigin = errors.New("create origin")
	// ErrUpdateOrigin is returned when UpdateOrigin fails
	ErrUpdateOrigin = errors.New("update origin")
)

// Validate validates ListOriginsRequest
func (v ListOriginsRequest) Validate() error {
	errs := validation.Errors{
		"Type": validation.Validate(v.Type, validation.In(OriginTypeCustomer, OriginTypeApplicationLoadBalancer, OriginTypeNetStorage, OriginTypeAll).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CUSTOMER', 'APPLICATION_LOAD_BALANCER', 'NETSTORAGE' or '' (empty)", (&v).Type))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates CreateOriginRequest
func (v CreateOriginRequest) Validate() error {
	errs := validation.Errors{
		"OriginID":    validation.Validate(v.OriginID, validation.Required, validation.Length(2, 63)),
		"Description": validation.Validate(v.Description.Description, validation.Length(0, 255)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates UpdateOriginRequest
func (v UpdateOriginRequest) Validate() error {
	errs := validation.Errors{
		"OriginID":    validation.Validate(v.OriginID, validation.Required, validation.Length(2, 63)),
		"Description": validation.Validate(v.Description.Description, validation.Length(0, 255)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

func (c *cloudlets) ListOrigins(ctx context.Context, params ListOriginsRequest) ([]OriginResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListOrigins")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListOrigins, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cloudlets/api/v2/origins").
		AddQueryParamIf("type", string(params.Type), params.Type != OriginTypeAll).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListOrigins, err)
	}

	var result []OriginResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListOrigins, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListOrigins, c.Error(resp))
	}

	return result, nil
}

func (c *cloudlets) GetOrigin(ctx context.Context, params GetOriginRequest) (*Origin, error) {
	logger := c.Log(ctx)
	logger.Debug("GetOrigin")

	req, err := request.NewGet(ctx, "/cloudlets/api/v2/origins/%s", params.OriginID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetOrigin, err)
	}

	var result Origin
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetOrigin, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetOrigin, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) CreateOrigin(ctx context.Context, params CreateOriginRequest) (*Origin, error) {
	logger := c.Log(ctx)
	logger.Debug("CreateOrigin")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateOrigin, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/cloudlets/api/v2/origins").
		WithBody(params).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateOrigin, err)
	}

	var result Origin

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateOrigin, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateOrigin, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) UpdateOrigin(ctx context.Context, params UpdateOriginRequest) (*Origin, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdateOrigin")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateOrigin, ErrStructValidation, err)
	}

	req, err := request.NewPut(ctx, "/cloudlets/api/v2/origins/%s", params.OriginID).
		WithBody(params.Description).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateOrigin, err)
	}

	var result Origin

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateOrigin, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateOrigin, c.Error(resp))
	}

	return &result, nil
}
