package gtm

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ResourceInstance contains information about the resources that constrain the properties within the data center
	ResourceInstance struct {
		DatacenterID         int  `json:"datacenterId"`
		UseDefaultLoadObject bool `json:"useDefaultLoadObject"`
		LoadObject
	}

	// Resource represents a GTM resource
	Resource struct {
		Type                        string             `json:"type"`
		HostHeader                  string             `json:"hostHeader,omitempty"`
		LeastSquaresDecay           float64            `json:"leastSquaresDecay,omitempty"`
		Description                 string             `json:"description,omitempty"`
		LeaderString                string             `json:"leaderString,omitempty"`
		ConstrainedProperty         string             `json:"constrainedProperty,omitempty"`
		ResourceInstances           []ResourceInstance `json:"resourceInstances,omitempty"`
		AggregationType             string             `json:"aggregationType,omitempty"`
		Links                       []Link             `json:"links,omitempty"`
		LoadImbalancePercentage     float64            `json:"loadImbalancePercentage,omitempty"`
		UpperBound                  int                `json:"upperBound,omitempty"`
		Name                        string             `json:"name"`
		MaxUMultiplicativeIncrement float64            `json:"maxUMultiplicativeIncrement,omitempty"`
		DecayRate                   float64            `json:"decayRate,omitempty"`
	}

	// ResourceList is the structure returned by List Resources
	ResourceList struct {
		ResourceItems []Resource `json:"items"`
	}

	// ListResourcesRequest contains request parameters for ListResources
	ListResourcesRequest struct {
		DomainName string
	}

	// GetResourceRequest contains request parameters for GetResource
	GetResourceRequest struct {
		DomainName   string
		ResourceName string
	}

	// GetResourceResponse contains the response data from GetResource operation
	GetResourceResponse Resource

	// ResourceRequest contains request parameters
	ResourceRequest struct {
		Resource   *Resource
		DomainName string
	}

	// CreateResourceRequest contains request parameters for CreateResource
	CreateResourceRequest ResourceRequest

	// CreateResourceResponse contains the response data from CreateResource operation
	CreateResourceResponse struct {
		Resource *Resource       `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// UpdateResourceRequest contains request parameters for UpdateResource
	UpdateResourceRequest ResourceRequest

	// UpdateResourceResponse contains the response data from UpdateResource operation
	UpdateResourceResponse struct {
		Resource *Resource       `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// DeleteResourceRequest contains request parameters for DeleteResource
	DeleteResourceRequest struct {
		DomainName   string
		ResourceName string
	}

	// DeleteResourceResponse contains the response data from DeleteResource operation
	DeleteResourceResponse struct {
		Resource *Resource       `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}
)

var (
	// ErrListResources is returned when ListResources fails
	ErrListResources = errors.New("list resources")
	// ErrGetResource is returned when GetResource fails
	ErrGetResource = errors.New("get resource")
	// ErrCreateResource is returned when CreateResource fails
	ErrCreateResource = errors.New("create resource")
	// ErrUpdateResource is returned when UpdateResource fails
	ErrUpdateResource = errors.New("update resource")
	// ErrDeleteResource is returned when DeleteResource fails
	ErrDeleteResource = errors.New("delete resource")
)

// Validate validates ListResourcesRequest
func (r ListResourcesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates GetResourceRequest
func (r GetResourceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName":   validation.Validate(r.DomainName, validation.Required),
		"ResourceName": validation.Validate(r.ResourceName, validation.Required),
	})
}

// Validate validates CreateResourceRequest
func (r CreateResourceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"Resource":   validation.Validate(r.Resource, validation.Required),
	})
}

// Validate validates UpdateResourceRequest
func (r UpdateResourceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"Resource":   validation.Validate(r.Resource, validation.Required),
	})
}

// Validate validates DeleteResourceRequest
func (r DeleteResourceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName":   validation.Validate(r.DomainName, validation.Required),
		"ResourceName": validation.Validate(r.ResourceName, validation.Required),
	})
}

// Validate validates Resource
func (r *Resource) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Name":            validation.Validate(r.Name, validation.Required),
		"Type":            validation.Validate(r.Type, validation.Required),
		"AggregationType": validation.Validate(r.AggregationType, validation.Required),
	})
}

func (g *gtm) ListResources(ctx context.Context, params ListResourcesRequest) ([]Resource, error) {
	logger := g.Log(ctx)
	logger.Debug("ListResources")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListResources, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources", params.DomainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListResources request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ResourceList
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListResources request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.ResourceItems, nil
}

func (g *gtm) GetResource(ctx context.Context, params GetResourceRequest) (*GetResourceResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("GetResource")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetResource, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", params.DomainName, params.ResourceName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetResource request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result GetResourceResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetResource request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateResource(ctx context.Context, params CreateResourceRequest) (*CreateResourceResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateResource")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateResource, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", params.DomainName, params.Resource.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Resource request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CreateResourceResponse
	resp, err := g.Exec(req, &result, params.Resource)
	if err != nil {
		return nil, fmt.Errorf("resource request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) UpdateResource(ctx context.Context, params UpdateResourceRequest) (*UpdateResourceResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateResource")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateResource, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", params.DomainName, params.Resource.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Resource request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result UpdateResourceResponse
	resp, err := g.Exec(req, &result, params.Resource)
	if err != nil {
		return nil, fmt.Errorf("resource request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteResource(ctx context.Context, params DeleteResourceRequest) (*DeleteResourceResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteResource")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteResource, ErrStructValidation, err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", params.DomainName, params.ResourceName)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DeleteResourceResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("DeleteResource request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}
