package gtm

import (
	"context"
	"fmt"
	"net/http"

	"strconv"
)

// Datacenters contains operations available on a Datacenter resource.
type Datacenters interface {
	// NewDatacenterResponse instantiates a new DatacenterResponse structure.
	NewDatacenterResponse(context.Context) *DatacenterResponse
	// NewDatacenter creates a new Datacenter object.
	NewDatacenter(context.Context) *Datacenter
	// ListDatacenters retrieves all Datacenters.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-datacenters
	ListDatacenters(context.Context, string) ([]*Datacenter, error)
	// GetDatacenter retrieves a Datacenter with the given name. NOTE: Id arg is int!
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-datacenter
	GetDatacenter(context.Context, int, string) (*Datacenter, error)
	// CreateDatacenter creates the datacenter identified by the receiver argument in the specified domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/post-datacenter
	CreateDatacenter(context.Context, *Datacenter, string) (*DatacenterResponse, error)
	// DeleteDatacenter deletes the datacenter identified by the receiver argument from the domain specified.
	//
	// See: https://techdocs.akamai.com/gtm/reference/delete-datacenter
	DeleteDatacenter(context.Context, *Datacenter, string) (*ResponseStatus, error)
	// UpdateDatacenter updates the datacenter identified in the receiver argument in the provided domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-datacenter
	UpdateDatacenter(context.Context, *Datacenter, string) (*ResponseStatus, error)
	// CreateMapsDefaultDatacenter creates Default Datacenter for Maps.
	CreateMapsDefaultDatacenter(context.Context, string) (*Datacenter, error)
	// CreateIPv4DefaultDatacenter creates Default Datacenter for IPv4 Selector.
	CreateIPv4DefaultDatacenter(context.Context, string) (*Datacenter, error)
	// CreateIPv6DefaultDatacenter creates Default Datacenter for IPv6 Selector.
	CreateIPv6DefaultDatacenter(context.Context, string) (*Datacenter, error)
}

// Datacenter represents a GTM datacenter
type Datacenter struct {
	City                          string      `json:"city,omitempty"`
	CloneOf                       int         `json:"cloneOf,omitempty"`
	CloudServerHostHeaderOverride bool        `json:"cloudServerHostHeaderOverride"`
	CloudServerTargeting          bool        `json:"cloudServerTargeting"`
	Continent                     string      `json:"continent,omitempty"`
	Country                       string      `json:"country,omitempty"`
	DefaultLoadObject             *LoadObject `json:"defaultLoadObject,omitempty"`
	Latitude                      float64     `json:"latitude,omitempty"`
	Links                         []*Link     `json:"links,omitempty"`
	Longitude                     float64     `json:"longitude,omitempty"`
	Nickname                      string      `json:"nickname,omitempty"`
	PingInterval                  int         `json:"pingInterval,omitempty"`
	PingPacketSize                int         `json:"pingPacketSize,omitempty"`
	DatacenterID                  int         `json:"datacenterId,omitempty"`
	ScorePenalty                  int         `json:"scorePenalty,omitempty"`
	ServermonitorLivenessCount    int         `json:"servermonitorLivenessCount,omitempty"`
	ServermonitorLoadCount        int         `json:"servermonitorLoadCount,omitempty"`
	ServermonitorPool             string      `json:"servermonitorPool,omitempty"`
	StateOrProvince               string      `json:"stateOrProvince,omitempty"`
	Virtual                       bool        `json:"virtual"`
}

// DatacenterList contains a list of Datacenters
type DatacenterList struct {
	DatacenterItems []*Datacenter `json:"items"`
}

func (g *gtm) NewDatacenterResponse(ctx context.Context) *DatacenterResponse {
	logger := g.Log(ctx)
	logger.Debug("NewDatacenterResponse")

	dcResp := &DatacenterResponse{}
	return dcResp
}

func (g *gtm) NewDatacenter(ctx context.Context) *Datacenter {
	logger := g.Log(ctx)
	logger.Debug("NewDatacenter")

	dc := &Datacenter{}
	return dc
}

func (g *gtm) ListDatacenters(ctx context.Context, domainName string) ([]*Datacenter, error) {
	logger := g.Log(ctx)
	logger.Debug("ListDatacenters")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListDatacenters request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DatacenterList
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListDatacenters request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.DatacenterItems, nil
}

func (g *gtm) GetDatacenter(ctx context.Context, dcID int, domainName string) (*Datacenter, error) {
	logger := g.Log(ctx)
	logger.Debug("GetDatacenter")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dcID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDatacenter request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result Datacenter
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetDatacenter request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateDatacenter(ctx context.Context, dc *Datacenter, domainName string) (*DatacenterResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateDatacenter")

	postURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Datacenter request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DatacenterResponse
	resp, err := g.Exec(req, &result, dc)
	if err != nil {
		return nil, fmt.Errorf("CreateDatacenter request failed: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

var (
	// MapDefaultDC is a default Datacenter ID for Maps
	MapDefaultDC = 5400
	// Ipv4DefaultDC is a default Datacenter ID for IPv4
	Ipv4DefaultDC = 5401
	// Ipv6DefaultDC is a default Datacenter ID for IPv6
	Ipv6DefaultDC = 5402
)

func (g *gtm) CreateMapsDefaultDatacenter(ctx context.Context, domainName string) (*Datacenter, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateMapsDefaultDatacenter")

	return createDefaultDC(ctx, g, MapDefaultDC, domainName)
}

func (g *gtm) CreateIPv4DefaultDatacenter(ctx context.Context, domainName string) (*Datacenter, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateIPv4DefaultDatacenter")

	return createDefaultDC(ctx, g, Ipv4DefaultDC, domainName)
}

func (g *gtm) CreateIPv6DefaultDatacenter(ctx context.Context, domainName string) (*Datacenter, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateIPv6DefaultDatacenter")

	return createDefaultDC(ctx, g, Ipv6DefaultDC, domainName)
}

// createDefaultDC is worker function used to create Default Datacenter identified id in the specified domain.
func createDefaultDC(ctx context.Context, g *gtm, defaultID int, domainName string) (*Datacenter, error) {
	if defaultID != MapDefaultDC && defaultID != Ipv4DefaultDC && defaultID != Ipv6DefaultDC {
		return nil, fmt.Errorf("invalid default datacenter id provided for creation")
	}
	// check if already exists
	dc, err := g.GetDatacenter(ctx, defaultID, domainName)
	if err == nil {
		return dc, err
	}
	apiError, ok := err.(*Error)
	if !ok || apiError.StatusCode != http.StatusNotFound {
		return nil, err
	}

	defaultURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/", domainName)
	switch defaultID {
	case MapDefaultDC:
		defaultURL += "default-datacenter-for-maps"
	case Ipv4DefaultDC:
		defaultURL += "datacenter-for-ip-version-selector-ipv4"
	case Ipv6DefaultDC:
		defaultURL += "datacenter-for-ip-version-selector-ipv6"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, defaultURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Default Datacenter request: %w", err)
	}

	setVersionHeader(req, schemaVersion)
	var result DatacenterResponse
	resp, err := g.Exec(req, &result, "")
	if err != nil {
		return nil, fmt.Errorf("DefaultDatacenter request failed: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return result.Resource, nil
}

func (g *gtm) UpdateDatacenter(ctx context.Context, dc *Datacenter, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateDatacenter")

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dc.DatacenterID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Update Datacenter request: %w", err)
	}

	setVersionHeader(req, schemaVersion)
	var result DatacenterResponse
	resp, err := g.Exec(req, &result, dc)
	if err != nil {
		return nil, fmt.Errorf("UpdateDatacenter request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.Status, nil
}

func (g *gtm) DeleteDatacenter(ctx context.Context, dc *Datacenter, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteDatacenter")

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dc.DatacenterID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete Datacenter request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DatacenterResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("DeleteDatacenter request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.Status, nil
}
