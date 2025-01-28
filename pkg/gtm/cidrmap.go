package gtm

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// CIDRAssignment represents a GTM CIDR assignment element
	CIDRAssignment struct {
		DatacenterBase
		Blocks []string `json:"blocks"`
	}

	// CIDRMap represents a GTM CIDRMap element
	CIDRMap struct {
		DefaultDatacenter *DatacenterBase  `json:"defaultDatacenter"`
		Assignments       []CIDRAssignment `json:"assignments,omitempty"`
		Name              string           `json:"name"`
		Links             []Link           `json:"links,omitempty"`
	}

	// CIDRMapList represents a GTM returned CIDRMap list body
	CIDRMapList struct {
		CIDRMapItems []CIDRMap `json:"items"`
	}
	// ListCIDRMapsRequest contains request parameters for ListCIDRMaps
	ListCIDRMapsRequest struct {
		DomainName string
	}

	// GetCIDRMapRequest contains request parameters for GetCIDRMap
	GetCIDRMapRequest struct {
		MapName    string
		DomainName string
	}

	// GetCIDRMapResponse contains the response data from GetCIDRMap operation
	GetCIDRMapResponse CIDRMap

	// CIDRMapRequest contains request parameters
	CIDRMapRequest struct {
		CIDR       *CIDRMap
		DomainName string
	}

	// CreateCIDRMapRequest contains request parameters for CreateCIDRMap
	CreateCIDRMapRequest CIDRMapRequest

	// CreateCIDRMapResponse contains the response data from CreateCIDRMap operation
	CreateCIDRMapResponse struct {
		Resource *CIDRMap        `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// UpdateCIDRMapRequest contains request parameters for UpdateCIDRMap
	UpdateCIDRMapRequest CIDRMapRequest

	// UpdateCIDRMapResponse contains the response data from UpdateCIDRMap operation
	UpdateCIDRMapResponse struct {
		Resource *CIDRMap        `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// DeleteCIDRMapRequest contains request parameters for DeleteCIDRMap
	DeleteCIDRMapRequest struct {
		MapName    string
		DomainName string
	}

	// DeleteCIDRMapResponse contains the response data from DeleteCIDRMap operation
	DeleteCIDRMapResponse struct {
		Resource *CIDRMap        `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}
)

var (
	// ErrListCIDRMaps is returned when ListCIDRMaps fails
	ErrListCIDRMaps = errors.New("list cidrmaps")
	// ErrGetCIDRMap is returned when GetCIDRMap fails
	ErrGetCIDRMap = errors.New("get cidrmap")
	// ErrCreateCIDRMap is returned when CreateCIDRMap fails
	ErrCreateCIDRMap = errors.New("create cidrmap")
	// ErrUpdateCIDRMap is returned when UpdateCIDRMap fails
	ErrUpdateCIDRMap = errors.New("update cidrmap")
	// ErrDeleteCIDRMap is returned when DeleteCIDRMap fails
	ErrDeleteCIDRMap = errors.New("delete cidrmap")
)

// Validate validates ListCIDRMapsRequest
func (r ListCIDRMapsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates GetCIDRMapRequest
func (r GetCIDRMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"MapName":    validation.Validate(r.MapName, validation.Required),
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates CreateCIDRMapRequest
func (r CreateCIDRMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"CIDRMap":    validation.Validate(r.CIDR, validation.Required),
	})
}

// Validate validates UpdateCIDRMapRequest
func (r UpdateCIDRMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"CIDRMap":    validation.Validate(r.CIDR, validation.Required),
	})
}

// Validate validates DeleteCIDRMapRequest
func (r DeleteCIDRMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"MapName":    validation.Validate(r.MapName, validation.Required),
	})
}

// Validate validates CIDRMap
func (c CIDRMap) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Name": validation.Validate(c.Name, validation.Required),
	})
}

func (g *gtm) ListCIDRMaps(ctx context.Context, params ListCIDRMapsRequest) ([]CIDRMap, error) {
	logger := g.Log(ctx)
	logger.Debug("ListCIDRMaps")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListCIDRMaps, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps", params.DomainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListCIDRMaps request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CIDRMapList
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListCIDRMaps request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.CIDRMapItems, nil
}

func (g *gtm) GetCIDRMap(ctx context.Context, params GetCIDRMapRequest) (*GetCIDRMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("GetCIDRMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCIDRMap, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", params.DomainName, params.MapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCIDRMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result GetCIDRMapResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCIDRMap request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateCIDRMap(ctx context.Context, params CreateCIDRMapRequest) (*CreateCIDRMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateCIDRMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateCIDRMap, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", params.DomainName, params.CIDR.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CIDRMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CreateCIDRMapResponse
	resp, err := g.Exec(req, &result, params.CIDR)
	if err != nil {
		return nil, fmt.Errorf("CIDRMap request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) UpdateCIDRMap(ctx context.Context, params UpdateCIDRMapRequest) (*UpdateCIDRMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateCIDRMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateCIDRMap, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", params.DomainName, params.CIDR.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CIDRMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result UpdateCIDRMapResponse
	resp, err := g.Exec(req, &result, params.CIDR)
	if err != nil {
		return nil, fmt.Errorf("CIDRMap request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteCIDRMap(ctx context.Context, params DeleteCIDRMapRequest) (*DeleteCIDRMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteCIDRMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteCIDRMap, ErrStructValidation, err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", params.DomainName, params.MapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DeleteCIDRMapResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("CIDRMap request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}
