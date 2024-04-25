package gtm

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GeoAssignment represents a GTM Geo assignment element
	GeoAssignment struct {
		DatacenterBase
		Countries []string `json:"countries"`
	}

	// GeoMap represents a GTM GeoMap
	GeoMap struct {
		DefaultDatacenter *DatacenterBase `json:"defaultDatacenter"`
		Assignments       []GeoAssignment `json:"assignments,omitempty"`
		Name              string          `json:"name"`
		Links             []Link          `json:"links,omitempty"`
	}

	// GeoMapList represents the returned GTM GeoMap List body
	GeoMapList struct {
		GeoMapItems []GeoMap `json:"items"`
	}

	// ListGeoMapsRequest contains request parameters for ListGeoMaps
	ListGeoMapsRequest struct {
		DomainName string
	}

	// GetGeoMapRequest contains request parameters for GetGeoMap
	GetGeoMapRequest struct {
		MapName    string
		DomainName string
	}

	// GetGeoMapResponse contains the response data from GetGeoMap operation
	GetGeoMapResponse GeoMap

	// GeoMapRequest contains request parameters
	GeoMapRequest struct {
		GeoMap     *GeoMap
		DomainName string
	}

	// CreateGeoMapRequest contains request parameters for CreateGeoMap
	CreateGeoMapRequest GeoMapRequest

	// CreateGeoMapResponse contains the response data from CreateGeoMap operation
	CreateGeoMapResponse struct {
		Resource *GeoMap         `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// UpdateGeoMapRequest contains request parameters for UpdateGeoMap
	UpdateGeoMapRequest GeoMapRequest

	// UpdateGeoMapResponse contains the response data from UpdateGeoMap operation
	UpdateGeoMapResponse struct {
		Resource *GeoMap         `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// DeleteGeoMapRequest contains request parameters for DeleteGeoMap
	DeleteGeoMapRequest struct {
		MapName    string
		DomainName string
	}

	// DeleteGeoMapResponse contains the response data from DeleteGeoMap operation
	DeleteGeoMapResponse struct {
		Resource *GeoMap         `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}
)

var (
	// ErrListGeoMaps is returned when ListGeoMaps fails
	ErrListGeoMaps = errors.New("list geomaps")
	// ErrGetGeoMap is returned when GetGeoMap fails
	ErrGetGeoMap = errors.New("get geomap")
	// ErrCreateGeoMap is returned when CreateGeoMap fails
	ErrCreateGeoMap = errors.New("create geomap")
	// ErrUpdateGeoMap is returned when UpdateGeoMap fails
	ErrUpdateGeoMap = errors.New("update geomap")
	// ErrDeleteGeoMap is returned when DeleteGeoMap fails
	ErrDeleteGeoMap = errors.New("delete geomap")
)

// Validate validates ListGeoMapsRequest
func (r ListGeoMapsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates GetGeoMapRequest
func (r GetGeoMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"MapName":    validation.Validate(r.MapName, validation.Required),
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates CreateGeoMapRequest
func (r CreateGeoMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"GeoMap":     validation.Validate(r.GeoMap, validation.Required),
	})
}

// Validate validates UpdateGeoMapRequest
func (r UpdateGeoMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"GeoMap":     validation.Validate(r.GeoMap, validation.Required),
	})
}

// Validate validates DeleteGeoMapRequest
func (r DeleteGeoMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"MapName":    validation.Validate(r.MapName, validation.Required),
	})
}

// Validate validates GeoMap
func (g *GeoMap) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Name":              validation.Validate(g.Name, validation.Required),
		"DefaultDatacenter": validation.Validate(g.DefaultDatacenter, validation.Required),
	})
}

func (g *gtm) ListGeoMaps(ctx context.Context, params ListGeoMapsRequest) ([]GeoMap, error) {
	logger := g.Log(ctx)
	logger.Debug("ListGeoMaps")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListGeoMaps, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps", params.DomainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListGeoMaps request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result GeoMapList
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListGeoMaps request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.GeoMapItems, nil
}

func (g *gtm) GetGeoMap(ctx context.Context, params GetGeoMapRequest) (*GetGeoMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("GetGeoMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetGeoMap, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", params.DomainName, params.MapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetGeoMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result GetGeoMapResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetGeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateGeoMap(ctx context.Context, params CreateGeoMapRequest) (*CreateGeoMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateGeoMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateGeoMap, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", params.DomainName, params.GeoMap.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GeoMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CreateGeoMapResponse
	resp, err := g.Exec(req, &result, params.GeoMap)
	if err != nil {
		return nil, fmt.Errorf("GeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) UpdateGeoMap(ctx context.Context, params UpdateGeoMapRequest) (*UpdateGeoMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateGeoMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateGeoMap, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", params.DomainName, params.GeoMap.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GeoMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result UpdateGeoMapResponse
	resp, err := g.Exec(req, &result, params.GeoMap)
	if err != nil {
		return nil, fmt.Errorf("GeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteGeoMap(ctx context.Context, params DeleteGeoMapRequest) (*DeleteGeoMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteGeoMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteGeoMap, ErrStructValidation, err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", params.DomainName, params.MapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DeleteGeoMapResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}
