package gtm

import (
	"context"
	"fmt"
	"net/http"
)

// CIDRMaps contains operations available on a CIDR map resource.
type CIDRMaps interface {
	// NewCIDRMap creates a new CIDRMap object.
	NewCIDRMap(context.Context, string) *CIDRMap
	// NewCIDRAssignment instantiates new Assignment struct.
	NewCIDRAssignment(context.Context, *CIDRMap, int, string) *CIDRAssignment
	// ListCIDRMaps retrieves all CIDRMaps.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-cidr-maps
	ListCIDRMaps(context.Context, string) ([]*CIDRMap, error)
	// GetCIDRMap retrieves a CIDRMap with the given name.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-cidr-map
	GetCIDRMap(context.Context, string, string) (*CIDRMap, error)
	// CreateCIDRMap creates the datacenter identified by the receiver argument in the specified domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-cidr-map
	CreateCIDRMap(context.Context, *CIDRMap, string) (*CIDRMapResponse, error)
	// DeleteCIDRMap deletes the datacenter identified by the receiver argument from the domain specified.
	//
	// See: https://techdocs.akamai.com/gtm/reference/delete-cidr-maps
	DeleteCIDRMap(context.Context, *CIDRMap, string) (*ResponseStatus, error)
	// UpdateCIDRMap updates the datacenter identified in the receiver argument in the provided domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-cidr-map
	UpdateCIDRMap(context.Context, *CIDRMap, string) (*ResponseStatus, error)
}

// CIDRAssignment represents a GTM CIDR assignment element
type CIDRAssignment struct {
	DatacenterBase
	Blocks []string `json:"blocks"`
}

// CIDRMap represents a GTM CIDRMap element
type CIDRMap struct {
	DefaultDatacenter *DatacenterBase   `json:"defaultDatacenter"`
	Assignments       []*CIDRAssignment `json:"assignments,omitempty"`
	Name              string            `json:"name"`
	Links             []*Link           `json:"links,omitempty"`
}

// CIDRMapList represents a GTM returned CIDRMap list body
type CIDRMapList struct {
	CIDRMapItems []*CIDRMap `json:"items"`
}

// Validate validates CIDRMap
func (c *CIDRMap) Validate() error {
	if len(c.Name) < 1 {
		return fmt.Errorf("CIDRMap is missing Name")
	}
	if c.DefaultDatacenter == nil {
		return fmt.Errorf("CIDRMap is missing DefaultDatacenter")
	}

	return nil
}

func (g *gtm) NewCIDRMap(ctx context.Context, name string) *CIDRMap {
	logger := g.Log(ctx)
	logger.Debug("NewCIDRMap")

	cidrMap := &CIDRMap{Name: name}
	return cidrMap
}

func (g *gtm) ListCIDRMaps(ctx context.Context, domainName string) ([]*CIDRMap, error) {
	logger := g.Log(ctx)
	logger.Debug("ListCIDRMaps")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps", domainName)
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

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.CIDRMapItems, nil
}

func (g *gtm) GetCIDRMap(ctx context.Context, mapName, domainName string) (*CIDRMap, error) {
	logger := g.Log(ctx)
	logger.Debug("GetCIDRMap")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", domainName, mapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCIDRMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CIDRMap
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCIDRMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) NewCIDRAssignment(ctx context.Context, _ *CIDRMap, dcID int, nickname string) *CIDRAssignment {
	logger := g.Log(ctx)
	logger.Debug("NewCIDRAssignment")

	cidrAssign := &CIDRAssignment{}
	cidrAssign.DatacenterID = dcID
	cidrAssign.Nickname = nickname

	return cidrAssign
}

func (g *gtm) CreateCIDRMap(ctx context.Context, cidr *CIDRMap, domainName string) (*CIDRMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateCIDRMap")

	return cidr.save(ctx, g, domainName)
}

func (g *gtm) UpdateCIDRMap(ctx context.Context, cidr *CIDRMap, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateCIDRMap")

	stat, err := cidr.save(ctx, g, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Save CIDRMap in given domain. Common path for Create and Update.
func (c *CIDRMap) save(ctx context.Context, g *gtm, domainName string) (*CIDRMapResponse, error) {

	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("CIDRMap validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", domainName, c.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CIDRMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CIDRMapResponse
	resp, err := g.Exec(req, &result, c)
	if err != nil {
		return nil, fmt.Errorf("CIDRMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteCIDRMap(ctx context.Context, cidr *CIDRMap, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteCIDRMap")

	if err := cidr.Validate(); err != nil {
		return nil, fmt.Errorf("CIDRMap validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", domainName, cidr.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ResponseBody
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("CIDRMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.Status, nil
}
