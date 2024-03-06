package gtm

import (
	"context"
	"fmt"
	"net/http"
)

// Resources contains operations available on a Resource resource.
type Resources interface {
	// NewResourceInstance instantiates a new ResourceInstance.
	NewResourceInstance(context.Context, *Resource, int) *ResourceInstance
	// NewResource creates a new Resource object.
	NewResource(context.Context, string) *Resource
	// ListResources retrieves all Resources
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-resources
	ListResources(context.Context, string) ([]*Resource, error)
	// GetResource retrieves a Resource with the given name.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-resource
	GetResource(context.Context, string, string) (*Resource, error)
	// CreateResource creates the datacenter identified by the receiver argument in the specified domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-resource
	CreateResource(context.Context, *Resource, string) (*ResourceResponse, error)
	// DeleteResource deletes the datacenter identified by the receiver argument from the domain specified.
	//
	// See: https://techdocs.akamai.com/gtm/reference/delete-resource
	DeleteResource(context.Context, *Resource, string) (*ResponseStatus, error)
	// UpdateResource updates the datacenter identified in the receiver argument in the provided domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-resource
	UpdateResource(context.Context, *Resource, string) (*ResponseStatus, error)
}

// ResourceInstance contains information about the resources that constrain the properties within the data center
type ResourceInstance struct {
	DatacenterID         int  `json:"datacenterId"`
	UseDefaultLoadObject bool `json:"useDefaultLoadObject"`
	LoadObject
}

// Resource represents a GTM resource
type Resource struct {
	Type                        string              `json:"type"`
	HostHeader                  string              `json:"hostHeader,omitempty"`
	LeastSquaresDecay           float64             `json:"leastSquaresDecay,omitempty"`
	Description                 string              `json:"description,omitempty"`
	LeaderString                string              `json:"leaderString,omitempty"`
	ConstrainedProperty         string              `json:"constrainedProperty,omitempty"`
	ResourceInstances           []*ResourceInstance `json:"resourceInstances,omitempty"`
	AggregationType             string              `json:"aggregationType,omitempty"`
	Links                       []*Link             `json:"links,omitempty"`
	LoadImbalancePercentage     float64             `json:"loadImbalancePercentage,omitempty"`
	UpperBound                  int                 `json:"upperBound,omitempty"`
	Name                        string              `json:"name"`
	MaxUMultiplicativeIncrement float64             `json:"maxUMultiplicativeIncrement,omitempty"`
	DecayRate                   float64             `json:"decayRate,omitempty"`
}

// ResourceList is the structure returned by List Resources
type ResourceList struct {
	ResourceItems []*Resource `json:"items"`
}

// Validate validates Resource
func (r *Resource) Validate() error {
	if len(r.Name) < 1 {
		return fmt.Errorf("resource is missing Name")
	}
	if len(r.Type) < 1 {
		return fmt.Errorf("resource is missing Type")
	}

	return nil
}

func (g *gtm) NewResourceInstance(ctx context.Context, _ *Resource, dcID int) *ResourceInstance {
	logger := g.Log(ctx)
	logger.Debug("NewResourceInstance")

	return &ResourceInstance{DatacenterID: dcID}
}

func (g *gtm) NewResource(ctx context.Context, name string) *Resource {
	logger := g.Log(ctx)
	logger.Debug("NewResource")

	resource := &Resource{Name: name}
	return resource
}

func (g *gtm) ListResources(ctx context.Context, domainName string) ([]*Resource, error) {
	logger := g.Log(ctx)
	logger.Debug("ListResources")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources", domainName)
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

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.ResourceItems, nil
}

func (g *gtm) GetResource(ctx context.Context, resourceName, domainName string) (*Resource, error) {
	logger := g.Log(ctx)
	logger.Debug("GetResource")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, resourceName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetResource request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result Resource
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetResource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateResource(ctx context.Context, resource *Resource, domainName string) (*ResourceResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateResource")

	return resource.save(ctx, g, domainName)
}

func (g *gtm) UpdateResource(ctx context.Context, resource *Resource, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateResource")

	stat, err := resource.save(ctx, g, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// save is a function that saves Resource in given domain. Common path for Create and Update.
func (r *Resource) save(ctx context.Context, g *gtm, domainName string) (*ResourceResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("resource validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, r.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Resource request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ResourceResponse
	resp, err := g.Exec(req, &result, r)
	if err != nil {
		return nil, fmt.Errorf("resource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteResource(ctx context.Context, resource *Resource, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteResource")

	if err := resource.Validate(); err != nil {
		logger.Errorf("Resource validation failed. %w", err)
		return nil, fmt.Errorf("DeleteResource validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, resource.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ResponseBody
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("DeleteResource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.Status, nil
}
