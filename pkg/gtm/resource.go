package gtm

import (
	"context"
	"fmt"
	"net/http"
)

//
// Handle Operations on gtm resources
// Based on 1.4 schema
//

// Resources contains operations available on a Resource resource.
type Resources interface {
	// NewResourceInstance instantiates a new ResourceInstance.
	NewResourceInstance(context.Context, *Resource, int) *ResourceInstance
	// NewResource creates a new Resource object.
	NewResource(context.Context, string) *Resource
	// ListResources retreieves all Resources
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
	DatacenterId         int  `json:"datacenterId"`
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
func (rsrc *Resource) Validate() error {

	if len(rsrc.Name) < 1 {
		return fmt.Errorf("Resource is missing Name")
	}
	if len(rsrc.Type) < 1 {
		return fmt.Errorf("Resource is missing Type")
	}

	return nil
}

func (p *gtm) NewResourceInstance(ctx context.Context, _ *Resource, dcID int) *ResourceInstance {

	logger := p.Log(ctx)
	logger.Debug("NewResourceInstance")

	return &ResourceInstance{DatacenterId: dcID}

}

func (p *gtm) NewResource(ctx context.Context, name string) *Resource {

	logger := p.Log(ctx)
	logger.Debug("NewResource")

	resource := &Resource{Name: name}
	return resource
}

func (p *gtm) ListResources(ctx context.Context, domainName string) ([]*Resource, error) {

	logger := p.Log(ctx)
	logger.Debug("ListResources")

	var rsrcs ResourceList
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListResources request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &rsrcs)
	if err != nil {
		return nil, fmt.Errorf("ListResources request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return rsrcs.ResourceItems, nil
}

func (p *gtm) GetResource(ctx context.Context, name, domainName string) (*Resource, error) {

	logger := p.Log(ctx)
	logger.Debug("GetResource")

	var rsc Resource
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetResource request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &rsc)
	if err != nil {
		return nil, fmt.Errorf("GetResource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rsc, nil
}

func (p *gtm) CreateResource(ctx context.Context, rsrc *Resource, domainName string) (*ResourceResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateResource")

	// Use common code. Any specific validation needed?
	return rsrc.save(ctx, p, domainName)

}

func (p *gtm) UpdateResource(ctx context.Context, rsrc *Resource, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateResource")

	// common code
	stat, err := rsrc.save(ctx, p, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err

}

// save is a function that saves Resource in given domain. Common path for Create and Update.
func (rsrc *Resource) save(ctx context.Context, p *gtm, domainName string) (*ResourceResponse, error) {

	if err := rsrc.Validate(); err != nil {
		return nil, fmt.Errorf("Resource validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, rsrc.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Resource request: %w", err)
	}

	var rscresp ResourceResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &rscresp, rsrc)
	if err != nil {
		return nil, fmt.Errorf("Resource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rscresp, nil

}

func (p *gtm) DeleteResource(ctx context.Context, rsrc *Resource, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteResource")

	if err := rsrc.Validate(); err != nil {
		logger.Errorf("Resource validation failed. %w", err)
		return nil, fmt.Errorf("Resource validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, rsrc.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}

	var rscresp ResponseBody
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &rscresp)
	if err != nil {
		return nil, fmt.Errorf("Resource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return rscresp.Status, nil

}
