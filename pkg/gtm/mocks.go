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

func (p *Mock) NewDomain(ctx context.Context, _ string, _ string) *Domain {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*Domain)
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

func (p *Mock) NewTrafficTarget(ctx context.Context) *TrafficTarget {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*TrafficTarget)
}

func (p *Mock) NewStaticRRSet(ctx context.Context) *StaticRRSet {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*StaticRRSet)
}

func (p *Mock) NewLivenessTest(ctx context.Context, a string, b string, c int, d float32) *LivenessTest {
	args := p.Called(ctx, a, b, c, d)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*LivenessTest)
}

func (p *Mock) NewProperty(ctx context.Context, prop string) *Property {
	args := p.Called(ctx, prop)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*Property)
}

func (p *Mock) ListProperties(ctx context.Context, domain string) ([]*Property, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Property), args.Error(1)
}

func (p *Mock) GetDatacenter(ctx context.Context, dcid int, domain string) (*Datacenter, error) {
	args := p.Called(ctx, dcid, domain)

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

func (p *Mock) NewDatacenterResponse(ctx context.Context) *DatacenterResponse {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DatacenterResponse)
}

func (p *Mock) NewDatacenter(ctx context.Context) *Datacenter {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*Datacenter)
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

func (p *Mock) GetResource(ctx context.Context, rsrc string, domain string) (*Resource, error) {
	args := p.Called(ctx, rsrc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Resource), args.Error(1)
}

func (p *Mock) CreateResource(ctx context.Context, rsrc *Resource, domain string) (*ResourceResponse, error) {
	args := p.Called(ctx, rsrc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResourceResponse), args.Error(1)
}

func (p *Mock) DeleteResource(ctx context.Context, rsrc *Resource, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, rsrc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateResource(ctx context.Context, rsrc *Resource, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, rsrc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) NewResourceInstance(ctx context.Context, ri *Resource, a int) *ResourceInstance {
	args := p.Called(ctx, ri, a)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*ResourceInstance)
}

func (p *Mock) NewResource(ctx context.Context, rname string) *Resource {
	args := p.Called(ctx, rname)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*Resource)
}

func (p *Mock) ListResources(ctx context.Context, domain string) ([]*Resource, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Resource), args.Error(1)
}

func (p *Mock) GetAsMap(ctx context.Context, asmap string, domain string) (*AsMap, error) {
	args := p.Called(ctx, asmap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*AsMap), args.Error(1)
}

func (p *Mock) CreateAsMap(ctx context.Context, asmap *AsMap, domain string) (*AsMapResponse, error) {
	args := p.Called(ctx, asmap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*AsMapResponse), args.Error(1)
}

func (p *Mock) DeleteAsMap(ctx context.Context, asmap *AsMap, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, asmap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateAsMap(ctx context.Context, asmap *AsMap, domain string) (*ResponseStatus, error) {

	args := p.Called(ctx, asmap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) NewAsMap(ctx context.Context, mname string) *AsMap {
	args := p.Called(ctx, mname)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*AsMap)
}

func (p *Mock) NewASAssignment(ctx context.Context, as *AsMap, a int, b string) *AsAssignment {
	args := p.Called(ctx, as, a, b)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*AsAssignment)
}

func (p *Mock) ListAsMaps(ctx context.Context, domain string) ([]*AsMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*AsMap), args.Error(1)
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

func (p *Mock) NewGeoMap(ctx context.Context, mname string) *GeoMap {
	args := p.Called(ctx, mname)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*GeoMap)
}

func (p *Mock) NewGeoAssignment(ctx context.Context, as *GeoMap, a int, b string) *GeoAssignment {
	args := p.Called(ctx, as, a, b)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*GeoAssignment)
}

func (p *Mock) ListGeoMaps(ctx context.Context, domain string) ([]*GeoMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*GeoMap), args.Error(1)
}

func (p *Mock) GetCidrMap(ctx context.Context, cidr string, domain string) (*CidrMap, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CidrMap), args.Error(1)
}

func (p *Mock) CreateCidrMap(ctx context.Context, cidr *CidrMap, domain string) (*CidrMapResponse, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CidrMapResponse), args.Error(1)
}

func (p *Mock) DeleteCidrMap(ctx context.Context, cidr *CidrMap, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) UpdateCidrMap(ctx context.Context, cidr *CidrMap, domain string) (*ResponseStatus, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ResponseStatus), args.Error(1)
}

func (p *Mock) NewCidrMap(ctx context.Context, mname string) *CidrMap {
	args := p.Called(ctx, mname)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*CidrMap)
}

func (p *Mock) NewCidrAssignment(ctx context.Context, as *CidrMap, a int, b string) *CidrAssignment {
	args := p.Called(ctx, as, a, b)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*CidrAssignment)
}

func (p *Mock) ListCidrMaps(ctx context.Context, domain string) ([]*CidrMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*CidrMap), args.Error(1)
}
