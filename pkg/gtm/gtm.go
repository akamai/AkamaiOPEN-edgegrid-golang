// Package gtm provides access to the Akamai GTM V1_4 APIs
//
// See: https://techdocs.akamai.com/gtm/reference/api
package gtm

import (
	"context"
	"errors"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")
)

type (
	// GTM is the gtm api interface
	GTM interface {
		// Domains

		// NullFieldMap retrieves map of null fields.
		NullFieldMap(context.Context, *Domain) (*NullFieldMapStruct, error)

		// GetDomainStatus retrieves current status for the given domain name.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-status-current
		GetDomainStatus(context.Context, GetDomainStatusRequest) (*GetDomainStatusResponse, error)

		// ListDomains retrieves all Domains.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-domains
		ListDomains(context.Context) ([]DomainItem, error)

		// GetDomain retrieves a Domain with the given domain name.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-domain
		GetDomain(context.Context, GetDomainRequest) (*GetDomainResponse, error)

		// CreateDomain creates domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/post-domain
		CreateDomain(context.Context, CreateDomainRequest) (*CreateDomainResponse, error)

		// Deprecated: DeleteDomain is deprecated and may be removed in future versions.
		// DeleteDomain is a method applied to a domain object resulting in removal.
		//
		// See: ** Not Supported by API **
		DeleteDomain(context.Context, DeleteDomainRequest) (*DeleteDomainResponse, error)

		// DeleteDomains submits a request to delete one or more new domains.
		//
		// See: ** API documentation (to be published in Akamai TechDocs). **
		DeleteDomains(context.Context, DeleteDomainsRequest) (*DeleteDomainsResponse, error)

		// GetDeleteDomainsStatus retrieves the current status of delete domains request.
		//
		// See: ** API documentation (to be published in Akamai TechDocs). **
		GetDeleteDomainsStatus(ctx context.Context, params DeleteDomainsStatusRequest) (*DeleteDomainsStatusResponse, error)

		// UpdateDomain is a method applied to a domain object resulting in an update.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-domain
		UpdateDomain(context.Context, UpdateDomainRequest) (*UpdateDomainResponse, error)

		// Properties

		// ListProperties retrieves all Properties for the provided domainName.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-properties
		ListProperties(context.Context, ListPropertiesRequest) ([]Property, error)

		// GetProperty retrieves a Property with the given domain and property names.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-property
		GetProperty(context.Context, GetPropertyRequest) (*GetPropertyResponse, error)

		// CreateProperty creates property.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-property
		CreateProperty(context.Context, CreatePropertyRequest) (*CreatePropertyResponse, error)

		// DeleteProperty is a method applied to a property object resulting in removal.
		//
		// See: https://techdocs.akamai.com/gtm/reference/delete-property
		DeleteProperty(context.Context, DeletePropertyRequest) (*DeletePropertyResponse, error)

		// UpdateProperty is a method applied to a property object resulting in an update.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-property
		UpdateProperty(context.Context, UpdatePropertyRequest) (*UpdatePropertyResponse, error)

		// Datacenters

		// ListDatacenters retrieves all Datacenters.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-datacenters
		ListDatacenters(context.Context, ListDatacentersRequest) ([]Datacenter, error)

		// GetDatacenter retrieves a Datacenter with the given name. NOTE: Id arg is int!
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-datacenter
		GetDatacenter(context.Context, GetDatacenterRequest) (*Datacenter, error)

		// CreateDatacenter creates the datacenter identified by the receiver argument in the specified domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/post-datacenter
		CreateDatacenter(context.Context, CreateDatacenterRequest) (*CreateDatacenterResponse, error)

		// DeleteDatacenter deletes the datacenter identified by the receiver argument from the domain specified.
		//
		// See: https://techdocs.akamai.com/gtm/reference/delete-datacenter
		DeleteDatacenter(context.Context, DeleteDatacenterRequest) (*DeleteDatacenterResponse, error)

		// UpdateDatacenter updates the datacenter identified in the receiver argument in the provided domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-datacenter
		UpdateDatacenter(context.Context, UpdateDatacenterRequest) (*UpdateDatacenterResponse, error)

		// CreateMapsDefaultDatacenter creates Default Datacenter for Maps.
		CreateMapsDefaultDatacenter(context.Context, string) (*Datacenter, error)

		// CreateIPv4DefaultDatacenter creates Default Datacenter for IPv4 Selector.
		CreateIPv4DefaultDatacenter(context.Context, string) (*Datacenter, error)

		// CreateIPv6DefaultDatacenter creates Default Datacenter for IPv6 Selector.
		CreateIPv6DefaultDatacenter(context.Context, string) (*Datacenter, error)

		// Resources

		// ListResources retrieves all Resources
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-resources
		ListResources(context.Context, ListResourcesRequest) ([]Resource, error)

		// GetResource retrieves a Resource with the given name.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-resource
		GetResource(context.Context, GetResourceRequest) (*GetResourceResponse, error)

		// CreateResource creates the datacenter identified by the receiver argument in the specified domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-resource
		CreateResource(context.Context, CreateResourceRequest) (*CreateResourceResponse, error)

		// DeleteResource deletes the datacenter identified by the receiver argument from the domain specified.
		//
		// See: https://techdocs.akamai.com/gtm/reference/delete-resource
		DeleteResource(context.Context, DeleteResourceRequest) (*DeleteResourceResponse, error)

		// UpdateResource updates the datacenter identified in the receiver argument in the provided domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-resource
		UpdateResource(context.Context, UpdateResourceRequest) (*UpdateResourceResponse, error)

		// ASMaps

		// ListASMaps retrieves all AsMaps.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-as-maps
		ListASMaps(context.Context, ListASMapsRequest) ([]ASMap, error)

		// GetASMap retrieves a AsMap with the given name.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-as-map
		GetASMap(context.Context, GetASMapRequest) (*GetASMapResponse, error)

		// CreateASMap creates the datacenter identified by the receiver argument in the specified domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-as-map
		CreateASMap(context.Context, CreateASMapRequest) (*CreateASMapResponse, error)

		// DeleteASMap deletes the datacenter identified by the receiver argument from the domain specified.
		//
		// See: https://techdocs.akamai.com/gtm/reference/delete-as-map
		DeleteASMap(context.Context, DeleteASMapRequest) (*DeleteASMapResponse, error)

		// UpdateASMap updates the datacenter identified in the receiver argument in the provided domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-as-map
		UpdateASMap(context.Context, UpdateASMapRequest) (*UpdateASMapResponse, error)

		// GeoMaps

		// ListGeoMaps retrieves all GeoMaps.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-geographic-maps
		ListGeoMaps(context.Context, ListGeoMapsRequest) ([]GeoMap, error)

		// GetGeoMap retrieves a GeoMap with the given name.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-geographic-map
		GetGeoMap(context.Context, GetGeoMapRequest) (*GetGeoMapResponse, error)

		// CreateGeoMap creates the datacenter identified by the receiver argument in the specified domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-geographic-map
		CreateGeoMap(context.Context, CreateGeoMapRequest) (*CreateGeoMapResponse, error)

		// DeleteGeoMap deletes the datacenter identified by the receiver argument from the domain specified.
		//
		// See: https://techdocs.akamai.com/gtm/reference/delete-geographic-map
		DeleteGeoMap(context.Context, DeleteGeoMapRequest) (*DeleteGeoMapResponse, error)

		// UpdateGeoMap updates the datacenter identified in the receiver argument in the provided domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-geographic-map
		UpdateGeoMap(context.Context, UpdateGeoMapRequest) (*UpdateGeoMapResponse, error)

		// CIDRMaps

		// ListCIDRMaps retrieves all CIDRMaps.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-cidr-maps
		ListCIDRMaps(context.Context, ListCIDRMapsRequest) ([]CIDRMap, error)

		// GetCIDRMap retrieves a CIDRMap with the given name.
		//
		// See: https://techdocs.akamai.com/gtm/reference/get-cidr-map
		GetCIDRMap(context.Context, GetCIDRMapRequest) (*GetCIDRMapResponse, error)

		// CreateCIDRMap creates the datacenter identified by the receiver argument in the specified domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-cidr-map
		CreateCIDRMap(context.Context, CreateCIDRMapRequest) (*CreateCIDRMapResponse, error)

		// DeleteCIDRMap deletes the datacenter identified by the receiver argument from the domain specified.
		//
		// See: https://techdocs.akamai.com/gtm/reference/delete-cidr-maps
		DeleteCIDRMap(context.Context, DeleteCIDRMapRequest) (*DeleteCIDRMapResponse, error)

		// UpdateCIDRMap updates the datacenter identified in the receiver argument in the provided domain.
		//
		// See: https://techdocs.akamai.com/gtm/reference/put-cidr-map
		UpdateCIDRMap(context.Context, UpdateCIDRMapRequest) (*UpdateCIDRMapResponse, error)
	}

	gtm struct {
		session.Session
	}

	// Option defines a GTM option
	Option func(*gtm)

	// ClientFunc is a gtm client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) GTM
)

// Client returns a new dns Client instance with the specified controller
func Client(sess session.Session, opts ...Option) GTM {
	p := &gtm{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Exec overrides the session.Exec to add dns options
func (g *gtm) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	return g.Session.Exec(r, out, in...)
}
