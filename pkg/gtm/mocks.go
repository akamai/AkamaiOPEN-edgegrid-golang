//revive:disable:exported

package gtm

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ GTM = &Mock{}

func (p *Mock) NullFieldMap(ctx context.Context, domain *Domain) (*NullFieldMapStruct, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*NullFieldMapStruct), args.Error(1)
}

func (p *Mock) GetDomainStatus(ctx context.Context, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) ListDomains(ctx context.Context) ([]*DomainItem, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*DomainItem), args.Error(1)
}

func (p *Mock) GetDomain(ctx context.Context, domain string) (*Domain, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Domain), args.Error(1)
}

func (p *Mock) CreateDomain(ctx context.Context, domain *Domain, queryArgs map[string]string) (*DomainResponse, error) {
	args := p.Called(ctx, domain, queryArgs)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DomainResponse), args.Error(1)
}

func (p *Mock) DeleteDomain(ctx context.Context, domain *Domain) (*ResponseStatus, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateDomain(ctx context.Context, domain *Domain, queryArgs map[string]string) (*ResponseStatus, error) {
	args := p.Called(ctx, domain, queryArgs)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) GetProperty(ctx context.Context, prop string, domain string) (*Property, error) {
	args := p.Called(ctx, prop, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Property), args.Error(1)
}

func (p *Mock) DeleteProperty(ctx context.Context, prop *Property, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, prop, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) CreateProperty(ctx context.Context, prop *Property, domain string) (*PropertyResponse, error) {
	args := p.Called(ctx, prop, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*PropertyResponse), args.Error(1)
}

func (p *Mock) UpdateProperty(ctx context.Context, prop *Property, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, prop, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) ListProperties(ctx context.Context, domain string) ([]*Property, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Property), args.Error(1)
}

func (p *Mock) GetDatacenter(ctx context.Context, dcID int, domain string) (*Datacenter, error) {
	args := p.Called(ctx, dcID, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Datacenter), args.Error(1)
}

func (p *Mock) CreateDatacenter(ctx context.Context, dc *Datacenter, domain string) (*DatacenterResponse, error) {
	args := p.Called(ctx, dc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DatacenterResponse), args.Error(1)
}

func (p *Mock) DeleteDatacenter(ctx context.Context, dc *Datacenter, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, dc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateDatacenter(ctx context.Context, dc *Datacenter, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, dc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) ListDatacenters(ctx context.Context, domain string) ([]*Datacenter, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Datacenter), args.Error(1)
}

func (p *Mock) CreateIPv4DefaultDatacenter(ctx context.Context, domain string) (*Datacenter, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Datacenter), args.Error(1)
}

func (p *Mock) CreateIPv6DefaultDatacenter(ctx context.Context, domain string) (*Datacenter, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Datacenter), args.Error(1)
}

func (p *Mock) CreateMapsDefaultDatacenter(ctx context.Context, domainName string) (*Datacenter, error) {
	args := p.Called(ctx, domainName)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Datacenter), args.Error(1)
}

func (p *Mock) GetResource(ctx context.Context, resource string, domain string) (*Resource, error) {
	args := p.Called(ctx, resource, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Resource), args.Error(1)
}

func (p *Mock) CreateResource(ctx context.Context, resource *Resource, domain string) (*ResourceResponse, error) {
	args := p.Called(ctx, resource, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResourceResponse), args.Error(1)
}

func (p *Mock) DeleteResource(ctx context.Context, resource *Resource, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, resource, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateResource(ctx context.Context, resource *Resource, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, resource, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) ListResources(ctx context.Context, domain string) ([]*Resource, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Resource), args.Error(1)
}

func (p *Mock) GetASMap(ctx context.Context, asMap string, domain string) (*ASMap, error) {
	args := p.Called(ctx, asMap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ASMap), args.Error(1)
}

func (p *Mock) CreateASMap(ctx context.Context, asMap *ASMap, domain string) (*ASMapResponse, error) {
	args := p.Called(ctx, asMap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ASMapResponse), args.Error(1)
}

func (p *Mock) DeleteASMap(ctx context.Context, asMap *ASMap, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, asMap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateASMap(ctx context.Context, asMap *ASMap, domain string) (*ResponseStatus, error) {

	args := p.Called(ctx, asMap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) ListASMaps(ctx context.Context, domain string) ([]*ASMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*ASMap), args.Error(1)
}

func (p *Mock) GetGeoMap(ctx context.Context, geo string, domain string) (*GeoMap, error) {
	args := p.Called(ctx, geo, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GeoMap), args.Error(1)
}

func (p *Mock) CreateGeoMap(ctx context.Context, geo *GeoMap, domain string) (*GeoMapResponse, error) {
	args := p.Called(ctx, geo, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GeoMapResponse), args.Error(1)
}

func (p *Mock) DeleteGeoMap(ctx context.Context, geo *GeoMap, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, geo, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateGeoMap(ctx context.Context, geo *GeoMap, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, geo, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) ListGeoMaps(ctx context.Context, domain string) ([]*GeoMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*GeoMap), args.Error(1)
}

func (p *Mock) GetCIDRMap(ctx context.Context, cidr string, domain string) (*CIDRMap, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CIDRMap), args.Error(1)
}

func (p *Mock) CreateCIDRMap(ctx context.Context, cidr *CIDRMap, domain string) (*CIDRMapResponse, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CIDRMapResponse), args.Error(1)
}

func (p *Mock) DeleteCIDRMap(ctx context.Context, cidr *CIDRMap, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateCIDRMap(ctx context.Context, cidr *CIDRMap, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) ListCIDRMaps(ctx context.Context, domain string) ([]*CIDRMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*CIDRMap), args.Error(1)
}
