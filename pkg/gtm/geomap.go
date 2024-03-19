package gtm

import (
	"context"
	"fmt"
	"net/http"
)

// GeoMaps contains operations available on a GeoMap resource.
type GeoMaps interface {
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

// GeoAssignment represents a GTM Geo assignment element
type GeoAssignment struct {
	DatacenterBase
	Countries []string `json:"countries"`
}

// GeoMap represents a GTM GeoMap
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
func (m *GeoMap) Validate() error {
	if len(m.Name) < 1 {
		return fmt.Errorf("GeoMap is missing Name")
	}
	if m.DefaultDatacenter == nil {
		return fmt.Errorf("GeoMap is missing DefaultDatacenter")
	}

	return nil
}

func (g *gtm) ListGeoMaps(ctx context.Context, domainName string) ([]*GeoMap, error) {
	logger := g.Log(ctx)
	logger.Debug("ListGeoMaps")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps", domainName)
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

func (g *gtm) GetGeoMap(ctx context.Context, mapName, domainName string) (*GeoMap, error) {
	logger := g.Log(ctx)
	logger.Debug("GetGeoMap")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", domainName, mapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetGeoMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result GeoMap
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetGeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateGeoMap(ctx context.Context, geo *GeoMap, domainName string) (*GeoMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateGeoMap")

	return geo.save(ctx, g, domainName)
}

func (g *gtm) UpdateGeoMap(ctx context.Context, geo *GeoMap, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateGeoMap")

	stat, err := geo.save(ctx, g, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Save GeoMap in given domain. Common path for Create and Update.
func (m *GeoMap) save(ctx context.Context, g *gtm, domainName string) (*GeoMapResponse, error) {
	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("GeoMap validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", domainName, m.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GeoMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result GeoMapResponse
	resp, err := g.Exec(req, &result, m)
	if err != nil {
		return nil, fmt.Errorf("GeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteGeoMap(ctx context.Context, geo *GeoMap, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
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
	setVersionHeader(req, schemaVersion)

	var result ResponseBody
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GeoMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.Status, nil
}
