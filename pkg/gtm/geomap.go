package gtm

import (
	"context"
	"fmt"
	"net/http"
)

//
// Handle Operations on gtm geomaps
// Based on 1.4 schema
//

// GeoMaps contains operations available on a GeoMap resource.
type GeoMaps interface {
	// NewGeoMap creates a new GeoMap object.
	NewGeoMap(context.Context, string) *GeoMap
	// NewGeoAssignment instantiates new Assignment struct.
	NewGeoAssignment(context.Context, *GeoMap, int, string) *GeoAssignment
	// ListGeoMaps retrieves all GeoMaps.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-geographic-maps
	ListGeoMaps(context.Context, string) ([]*GeoMap, error)
	// GetGeoMap retrieves a GeoMap with the given name.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-geographic-map
	GetGeoMap(context.Context, string, string) (*GeoMap, error)
	// CreateGeoMap creates the datacenter identified by the receiver argument in the specified domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-geographic-map
	CreateGeoMap(context.Context, *GeoMap, string) (*GeoMapResponse, error)
	// DeleteGeoMap deletes the datacenter identified by the receiver argument from the domain specified.
	//
	// See: https://techdocs.akamai.com/gtm/reference/delete-geographic-map
	DeleteGeoMap(context.Context, *GeoMap, string) (*ResponseStatus, error)
	// UpdateGeoMap updates the datacenter identified in the receiver argument in the provided domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-geographic-map
	UpdateGeoMap(context.Context, *GeoMap, string) (*ResponseStatus, error)
}

// GeoAssignment represents a GTM geo assignment element
type GeoAssignment struct {
	DatacenterBase
	Countries []string `json:"countries"`
}

// GeoMap  represents a GTM GeoMap
type GeoMap struct {
	DefaultDatacenter *DatacenterBase  `json:"defaultDatacenter"`
	Assignments       []*GeoAssignment `json:"assignments,omitempty"`
	Name              string           `json:"name"`
	Links             []*Link          `json:"links,omitempty"`
}

// GeoMapList represents the returned GTM GeoMap List body
type GeoMapList struct {
	GeoMapItems []*GeoMap `json:"items"`
}

// Validate validates GeoMap
func (geo *GeoMap) Validate() error {

	if len(geo.Name) < 1 {
		return fmt.Errorf("GeoMap is missing Name")
	}
	if geo.DefaultDatacenter == nil {
		return fmt.Errorf("GeoMap is missing DefaultDatacenter")
	}

	return nil
}

func (p *gtm) NewGeoMap(ctx context.Context, name string) *GeoMap {

	logger := p.Log(ctx)
	logger.Debug("NewGeoMap")

	geomap := &GeoMap{Name: name}
	return geomap
}

func (p *gtm) ListGeoMaps(ctx context.Context, domainName string) ([]*GeoMap, error) {

	logger := p.Log(ctx)
	logger.Debug("ListGeoMaps")

	var geos GeoMapList
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListGeoMaps request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &geos)
	if err != nil {
		return nil, fmt.Errorf("ListGeoMaps request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return geos.GeoMapItems, nil
}

func (p *gtm) GetGeoMap(ctx context.Context, name, domainName string) (*GeoMap, error) {

	logger := p.Log(ctx)
	logger.Debug("GetGeoMap")

	var geo GeoMap
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", domainName, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetGeoMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &geo)
	if err != nil {
		return nil, fmt.Errorf("GetGeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &geo, nil
}

func (p *gtm) NewGeoAssignment(ctx context.Context, _ *GeoMap, dcID int, nickname string) *GeoAssignment {

	logger := p.Log(ctx)
	logger.Debug("NewGeoAssignment")

	geoAssign := &GeoAssignment{}
	geoAssign.DatacenterId = dcID
	geoAssign.Nickname = nickname

	return geoAssign
}

func (p *gtm) CreateGeoMap(ctx context.Context, geo *GeoMap, domainName string) (*GeoMapResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateGeoMap")

	// Use common code. Any specific validation needed?
	return geo.save(ctx, p, domainName)
}

func (p *gtm) UpdateGeoMap(ctx context.Context, geo *GeoMap, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateGeoMap")

	// common code
	stat, err := geo.save(ctx, p, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Save GeoMap in given domain. Common path for Create and Update.
func (geo *GeoMap) save(ctx context.Context, p *gtm, domainName string) (*GeoMapResponse, error) {

	if err := geo.Validate(); err != nil {
		return nil, fmt.Errorf("GeoMap validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", domainName, geo.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GeoMap request: %w", err)
	}

	var mapresp GeoMapResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &mapresp, geo)
	if err != nil {
		return nil, fmt.Errorf("GeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &mapresp, nil
}

func (p *gtm) DeleteGeoMap(ctx context.Context, geo *GeoMap, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteGeoMap")

	if err := geo.Validate(); err != nil {
		logger.Errorf("Resource validation failed. %w", err)
		return nil, fmt.Errorf("GeoMap validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", domainName, geo.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}

	var mapresp ResponseBody
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &mapresp)
	if err != nil {
		return nil, fmt.Errorf("GeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return mapresp.Status, nil
}
