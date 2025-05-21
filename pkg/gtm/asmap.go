package gtm

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ASAssignment represents a GTM as map assignment structure
	ASAssignment struct {
		DatacenterBase
		ASNumbers []int64 `json:"asNumbers"`
	}

	// ASMap  represents a GTM ASMap
	ASMap struct {
		DefaultDatacenter *DatacenterBase `json:"defaultDatacenter"`
		Assignments       []ASAssignment  `json:"assignments,omitempty"`
		Name              string          `json:"name"`
		Links             []Link          `json:"links,omitempty"`
	}

	// ASMapList represents the returned GTM ASMap List body
	ASMapList struct {
		ASMapItems []ASMap `json:"items"`
	}

	// ListASMapsRequest contains request parameters for ListASMaps
	ListASMapsRequest struct {
		DomainName string
	}

	// GetASMapRequest contains request parameters for GetASMap
	GetASMapRequest struct {
		ASMapName  string
		DomainName string
	}

	// GetASMapResponse contains the response data from GetASMap operation
	GetASMapResponse ASMap

	// ASMapRequest contains request parameters
	ASMapRequest struct {
		ASMap      *ASMap
		DomainName string
	}

	// CreateASMapRequest contains request parameters for CreateASMap
	CreateASMapRequest ASMapRequest

	// CreateASMapResponse contains the response data from CreateASMap operation
	CreateASMapResponse struct {
		Resource *ASMap          `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// UpdateASMapRequest contains request parameters for UpdateASMap
	UpdateASMapRequest ASMapRequest

	// UpdateASMapResponse contains the response data from UpdateASMap operation
	UpdateASMapResponse struct {
		Resource *ASMap          `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// DeleteASMapRequest contains request parameters for DeleteASMap
	DeleteASMapRequest struct {
		ASMapName  string
		DomainName string
	}

	// DeleteASMapResponse contains the response data from DeleteASMap operation
	DeleteASMapResponse struct {
		Resource *ASMap          `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}
)

var (
	// ErrListASMaps is returned when ListASMaps fails
	ErrListASMaps = errors.New("list asmaps")
	// ErrGetASMap is returned when GetASMap fails
	ErrGetASMap = errors.New("get asmap")
	// ErrCreateASMap is returned when CreateASMap fails
	ErrCreateASMap = errors.New("create asmap")
	// ErrUpdateASMap is returned when UpdateASMap fails
	ErrUpdateASMap = errors.New("update asmap")
	// ErrDeleteASMap is returned when DeleteASMap fails
	ErrDeleteASMap = errors.New("delete asmap")
)

// Validate validates ListASMapsRequest
func (r ListASMapsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates GetASMapRequest
func (r GetASMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ASMapName":  validation.Validate(r.ASMapName, validation.Required),
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates CreateASMapRequest
func (r CreateASMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"ASMap":      validation.Validate(r.ASMap, validation.Required),
	})
}

// Validate validates ASMap
func (a ASMap) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Name":              validation.Validate(a.Name, validation.Required),
		"DefaultDatacenter": validation.Validate(a.DefaultDatacenter, validation.Required),
		"Assignments":       validation.Validate(a.Assignments, validation.Required),
	})
}

// Validate validates UpdateASMapRequest
func (r UpdateASMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"ASMap":      validation.Validate(r.ASMap, validation.Required),
	})
}

// Validate validates DeleteASMapRequest
func (r DeleteASMapRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"ASMapName":  validation.Validate(r.ASMapName, validation.Required),
	})
}

func (g *gtm) ListASMaps(ctx context.Context, params ListASMapsRequest) ([]ASMap, error) {
	logger := g.Log(ctx)
	logger.Debug("ListASMaps")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListASMaps, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps", params.DomainName)
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.ASMapItems, nil
}

func (g *gtm) GetASMap(ctx context.Context, params GetASMapRequest) (*GetASMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("GetASMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetASMap, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", params.DomainName, params.ASMapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetASMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result GetASMapResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetASMap request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateASMap(ctx context.Context, params CreateASMapRequest) (*CreateASMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateASMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateASMap, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", params.DomainName, params.ASMap.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ASMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CreateASMapResponse
	resp, err := g.Exec(req, &result, params.ASMap)
	if err != nil {
		return nil, fmt.Errorf("ASMap request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) UpdateASMap(ctx context.Context, params UpdateASMapRequest) (*UpdateASMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateASMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateASMap, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", params.DomainName, params.ASMap.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ASMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result UpdateASMapResponse
	resp, err := g.Exec(req, &result, params.ASMap)
	if err != nil {
		return nil, fmt.Errorf("ASMap request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteASMap(ctx context.Context, params DeleteASMapRequest) (*DeleteASMapResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteASMap")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteASMap, ErrStructValidation, err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", params.DomainName, params.ASMapName)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DeleteASMapResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ASMap request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}
