package gtm

import (
	"context"
	"fmt"
	"net/http"
)

//
// Handle Operations on gtm asmaps
// Based on 1.4 schema
//

// ASMaps contains operations available on a ASmap resource.
type ASMaps interface {
	// NewAsMap creates a new AsMap object.
	NewAsMap(context.Context, string) *AsMap
	// NewASAssignment instantiates new Assignment struct.
	NewASAssignment(context.Context, *AsMap, int, string) *AsAssignment
	// ListAsMaps retrieves all AsMaps.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-as-maps
	ListAsMaps(context.Context, string) ([]*AsMap, error)
	// GetAsMap retrieves a AsMap with the given name.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-as-map
	GetAsMap(context.Context, string, string) (*AsMap, error)
	// CreateAsMap creates the datacenter identified by the receiver argument in the specified domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-as-map
	CreateAsMap(context.Context, *AsMap, string) (*AsMapResponse, error)
	// DeleteAsMap deletes the datacenter identified by the receiver argument from the domain specified.
	//
	// See: https://techdocs.akamai.com/gtm/reference/delete-as-map
	DeleteAsMap(context.Context, *AsMap, string) (*ResponseStatus, error)
	// UpdateAsMap updates the datacenter identified in the receiver argument in the provided domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-as-map
	UpdateAsMap(context.Context, *AsMap, string) (*ResponseStatus, error)
}

// AsAssignment represents a GTM asmap assignment structure
type AsAssignment struct {
	DatacenterBase
	AsNumbers []int64 `json:"asNumbers"`
}

// AsMap  represents a GTM AsMap
type AsMap struct {
	DefaultDatacenter *DatacenterBase `json:"defaultDatacenter"`
	Assignments       []*AsAssignment `json:"assignments,omitempty"`
	Name              string          `json:"name"`
	Links             []*Link         `json:"links,omitempty"`
}

// AsMapList represents the returned GTM AsMap List body
type AsMapList struct {
	AsMapItems []*AsMap `json:"items"`
}

// Validate validates AsMap
func (asm *AsMap) Validate() error {

	if len(asm.Name) < 1 {
		return fmt.Errorf("AsMap is missing Name")
	}
	if asm.DefaultDatacenter == nil {
		return fmt.Errorf("AsMap is missing DefaultDatacenter")
	}

	return nil
}

func (p *gtm) NewAsMap(ctx context.Context, name string) *AsMap {

	logger := p.Log(ctx)
	logger.Debug("NewAsMap")

	asmap := &AsMap{Name: name}
	return asmap
}

func (p *gtm) ListAsMaps(ctx context.Context, domainName string) ([]*AsMap, error) {

	logger := p.Log(ctx)
	logger.Debug("ListAsMaps")

	var aslist AsMapList
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListAsMaps request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &aslist)
	if err != nil {
		return nil, fmt.Errorf("ListAsMaps request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return aslist.AsMapItems, nil
}

func (p *gtm) GetAsMap(ctx context.Context, name, domainName string) (*AsMap, error) {

	logger := p.Log(ctx)
	logger.Debug("GetAsMap")

	var as AsMap
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAsMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &as)
	if err != nil {
		return nil, fmt.Errorf("GetAsMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &as, nil
}

func (p *gtm) NewASAssignment(ctx context.Context, _ *AsMap, dcID int, nickname string) *AsAssignment {

	logger := p.Log(ctx)
	logger.Debug("NewAssignment")

	asAssign := &AsAssignment{}
	asAssign.DatacenterId = dcID
	asAssign.Nickname = nickname

	return asAssign
}

func (p *gtm) CreateAsMap(ctx context.Context, as *AsMap, domainName string) (*AsMapResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateAsMap")

	// Use common code. Any specific validation needed?
	return as.save(ctx, p, domainName)
}

func (p *gtm) UpdateAsMap(ctx context.Context, as *AsMap, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateAsMap")

	// common code
	stat, err := as.save(ctx, p, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// save AsMap in given domain. Common path for Create and Update.
func (asm *AsMap) save(ctx context.Context, p *gtm, domainName string) (*AsMapResponse, error) {

	if err := asm.Validate(); err != nil {
		return nil, fmt.Errorf("AsMap validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, asm.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create AsMap request: %w", err)
	}

	var mapresp AsMapResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &mapresp, asm)
	if err != nil {
		return nil, fmt.Errorf("AsMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &mapresp, nil
}

func (p *gtm) DeleteAsMap(ctx context.Context, as *AsMap, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteAsMap")

	if err := as.Validate(); err != nil {
		return nil, fmt.Errorf("Resource validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, as.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}

	var mapresp ResponseBody
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &mapresp)
	if err != nil {
		return nil, fmt.Errorf("AsMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return mapresp.Status, nil
}
