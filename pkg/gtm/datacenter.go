package gtm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Datacenter represents a GTM datacenter
	Datacenter struct {
		City                          string      `json:"city,omitempty"`
		CloneOf                       int         `json:"cloneOf,omitempty"`
		CloudServerHostHeaderOverride bool        `json:"cloudServerHostHeaderOverride"`
		CloudServerTargeting          bool        `json:"cloudServerTargeting"`
		Continent                     string      `json:"continent,omitempty"`
		Country                       string      `json:"country,omitempty"`
		DefaultLoadObject             *LoadObject `json:"defaultLoadObject,omitempty"`
		Latitude                      float64     `json:"latitude,omitempty"`
		Links                         []Link      `json:"links,omitempty"`
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
	DatacenterList struct {
		DatacenterItems []Datacenter `json:"items"`
	}

	// ListDatacentersRequest contains request parameters for ListDatacenters
	ListDatacentersRequest struct {
		DomainName string
	}

	// GetDatacenterRequest contains request parameters for GetDatacenter
	GetDatacenterRequest struct {
		DatacenterID int
		DomainName   string
	}

	// DatacenterRequest contains request parameters
	DatacenterRequest struct {
		Datacenter *Datacenter
		DomainName string
	}

	// CreateDatacenterRequest contains request parameters for CreateDatacenter
	CreateDatacenterRequest DatacenterRequest

	// CreateDatacenterResponse contains the response data from CreateDatacenter operation
	CreateDatacenterResponse struct {
		Status   *ResponseStatus `json:"status"`
		Resource *Datacenter     `json:"resource"`
	}

	// UpdateDatacenterRequest contains request parameters for UpdateDatacenter
	UpdateDatacenterRequest DatacenterRequest

	// UpdateDatacenterResponse contains the response data from UpdateDatacenter operation
	UpdateDatacenterResponse struct {
		Status   *ResponseStatus `json:"status"`
		Resource *Datacenter     `json:"resource"`
	}

	// DeleteDatacenterRequest contains request parameters for DeleteDatacenter
	DeleteDatacenterRequest struct {
		DatacenterID int
		DomainName   string
	}

	// DeleteDatacenterResponse contains the response data from DeleteDatacenter operation
	DeleteDatacenterResponse struct {
		Status   *ResponseStatus `json:"status"`
		Resource *Datacenter     `json:"resource"`
	}
)

var (
	// ErrListDatacenters is returned when ListDatacenters fails
	ErrListDatacenters = errors.New("list datacenters")
	// ErrGetDatacenter is returned when GetDatacenter fails
	ErrGetDatacenter = errors.New("get datacenter")
	// ErrCreateDatacenter is returned when CreateDatacenter fails
	ErrCreateDatacenter = errors.New("create datacenter")
	// ErrUpdateDatacenter is returned when UpdateDatacenter fails
	ErrUpdateDatacenter = errors.New("update datacenter")
	// ErrDeleteDatacenter is returned when DeleteDatacenter fails
	ErrDeleteDatacenter = errors.New("delete datacenter")
)

// Validate validates ListDatacentersRequest
func (r ListDatacentersRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates GetDatacenterRequest
func (r GetDatacenterRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DatacenterID": validation.Validate(r.DatacenterID, validation.Required),
		"DomainName":   validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates CreateDatacenterRequest
func (r CreateDatacenterRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"Datacenter": validation.Validate(r.Datacenter, validation.Required),
	})
}

// Validate validates UpdateDatacenterRequest
func (r UpdateDatacenterRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"Datacenter": validation.Validate(r.Datacenter, validation.Required),
	})
}

// Validate validates DeleteDatacenterRequest
func (r DeleteDatacenterRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName":   validation.Validate(r.DomainName, validation.Required),
		"DatacenterID": validation.Validate(r.DatacenterID, validation.Required),
	})
}

// Validate validates Datacenter
func (d Datacenter) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DatacenterID": validation.Validate(d.DatacenterID),
	})
}

func (g *gtm) ListDatacenters(ctx context.Context, params ListDatacentersRequest) ([]Datacenter, error) {
	logger := g.Log(ctx)
	logger.Debug("ListDatacenters")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListDatacenters, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters", params.DomainName)
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.DatacenterItems, nil
}

func (g *gtm) GetDatacenter(ctx context.Context, params GetDatacenterRequest) (*Datacenter, error) {
	logger := g.Log(ctx)
	logger.Debug("GetDatacenter")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetDatacenter, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", params.DomainName, strconv.Itoa(params.DatacenterID))
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateDatacenter(ctx context.Context, params CreateDatacenterRequest) (*CreateDatacenterResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateDatacenter")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateDatacenter, ErrStructValidation, err)
	}

	postURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters", params.DomainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Datacenter request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CreateDatacenterResponse
	resp, err := g.Exec(req, &result, params.Datacenter)
	if err != nil {
		return nil, fmt.Errorf("CreateDatacenter request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

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
	dc, err := g.GetDatacenter(ctx, GetDatacenterRequest{
		DatacenterID: defaultID,
		DomainName:   domainName,
	})
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return result.Resource, nil
}

func (g *gtm) UpdateDatacenter(ctx context.Context, params UpdateDatacenterRequest) (*UpdateDatacenterResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateDatacenter")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateDatacenter, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", params.DomainName, strconv.Itoa(params.Datacenter.DatacenterID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Update Datacenter request: %w", err)
	}

	setVersionHeader(req, schemaVersion)
	var result UpdateDatacenterResponse
	resp, err := g.Exec(req, &result, params.Datacenter)
	if err != nil {
		return nil, fmt.Errorf("UpdateDatacenter request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteDatacenter(ctx context.Context, params DeleteDatacenterRequest) (*DeleteDatacenterResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteDatacenter")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteDatacenter, ErrStructValidation, err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", params.DomainName, strconv.Itoa(params.DatacenterID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete Datacenter request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DeleteDatacenterResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("DeleteDatacenter request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}
