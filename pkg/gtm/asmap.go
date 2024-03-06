package gtm

import (
	"context"
	"fmt"
	"net/http"
)

// ASMaps contains operations available on a ASmap resource.
type (
	ASMaps interface {
		// NewASMap creates a new AsMap object.
		NewASMap(context.Context, string) *ASMap
		// NewASAssignment instantiates new Assignment struct.
		NewASAssignment(context.Context, *ASMap, int, string) *ASAssignment
		// ListASMaps retrieves all AsMaps.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-as-maps
		ListASMaps(context.Context, string) ([]*ASMap, error)
		// GetASMap retrieves a AsMap with the given name.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-as-map
		GetASMap(context.Context, string, string) (*ASMap, error)
		// CreateASMap creates the datacenter identified by the receiver argument in the specified domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-as-map
		CreateASMap(context.Context, *ASMap, string) (*ASMapResponse, error)
		// DeleteASMap deletes the datacenter identified by the receiver argument from the domain specified.
		//
		// See: https://techdocs.akamai.com/gtm/reference/delete-as-map
		DeleteASMap(context.Context, *ASMap, string) (*ResponseStatus, error)
		// UpdateASMap updates the datacenter identified in the receiver argument in the provided domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-as-map
		UpdateASMap(context.Context, *ASMap, string) (*ResponseStatus, error)
	}
	// ASAssignment represents a GTM as map assignment structure
	ASAssignment struct {
		DatacenterBase
		ASNumbers []int64 `json:"asNumbers"`
	}

	// ASMap  represents a GTM ASMap
	ASMap struct {
		DefaultDatacenter *DatacenterBase `json:"defaultDatacenter"`
		Assignments       []*ASAssignment `json:"assignments,omitempty"`
		Name              string          `json:"name"`
		Links             []*Link         `json:"links,omitempty"`
	}

	// ASMapList represents the returned GTM ASMap List body
	ASMapList struct {
		ASMapItems []*ASMap `json:"items"`
	}
)

// Validate validates ASMap
func (a *ASMap) Validate() error {
	if len(a.Name) < 1 {
		return fmt.Errorf("ASMap is missing Name")
	}
	if a.DefaultDatacenter == nil {
		return fmt.Errorf("ASMap is missing DefaultDatacenter")
	}

	return nil
}

func (g *gtm) NewASMap(ctx context.Context, name string) *ASMap {
	logger := g.Log(ctx)
	logger.Debug("NewASMap")

	asMap := &ASMap{Name: name}
	return asMap
}

func (g *gtm) ListASMaps(ctx context.Context, domainName string) ([]*ASMap, error) {
	logger := g.Log(ctx)
	logger.Debug("ListASMaps")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListASMaps request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ASMapList
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListASMaps request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.ASMapItems, nil
}

func (g *gtm) GetASMap(ctx context.Context, asMapName, domainName string) (*ASMap, error) {
	logger := g.Log(ctx)
	logger.Debug("GetASMap")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, asMapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetASMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ASMap
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetASMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) NewASAssignment(ctx context.Context, _ *ASMap, dcID int, nickname string) *ASAssignment {
	logger := g.Log(ctx)
	logger.Debug("NewASAssignment")

	asAssign := &ASAssignment{}
	asAssign.DatacenterID = dcID
	asAssign.Nickname = nickname

	return asAssign
}

func (g *gtm) CreateASMap(ctx context.Context, asMap *ASMap, domainName string) (*ASMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateASMap")

	return asMap.save(ctx, g, domainName)
}

func (g *gtm) UpdateASMap(ctx context.Context, asMap *ASMap, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateASMap")

	stat, err := asMap.save(ctx, g, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// save AsMap in given domain. Common path for Create and Update.
func (a *ASMap) save(ctx context.Context, g *gtm, domainName string) (*ASMapResponse, error) {
	if err := a.Validate(); err != nil {
		return nil, fmt.Errorf("ASMap validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, a.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ASMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ASMapResponse
	resp, err := g.Exec(req, &result, a)
	if err != nil {
		return nil, fmt.Errorf("ASMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteASMap(ctx context.Context, asMap *ASMap, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteASMap")

	if err := asMap.Validate(); err != nil {
		return nil, fmt.Errorf("resource validation failed: %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, asMap.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ResponseBody
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ASMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.Status, nil
}
