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

func (p *Mock) GetDomainStatus(ctx context.Context, req GetDomainStatusRequest) (*GetDomainStatusResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetDomainStatusResponse), args.Error(1)
}

func (p *Mock) ListDomains(ctx context.Context) ([]DomainItem, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]DomainItem), args.Error(1)
}

func (p *Mock) GetDomain(ctx context.Context, req GetDomainRequest) (*GetDomainResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetDomainResponse), args.Error(1)
}

func (p *Mock) CreateDomain(ctx context.Context, req CreateDomainRequest) (*CreateDomainResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateDomainResponse), args.Error(1)
}

func (p *Mock) DeleteDomain(ctx context.Context, req DeleteDomainRequest) (*DeleteDomainResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteDomainResponse), args.Error(1)
}

func (p *Mock) UpdateDomain(ctx context.Context, req UpdateDomainRequest) (*UpdateDomainResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateDomainResponse), args.Error(1)
}

func (p *Mock) GetProperty(ctx context.Context, req GetPropertyRequest) (*GetPropertyResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPropertyResponse), args.Error(1)
}

func (p *Mock) DeleteProperty(ctx context.Context, req DeletePropertyRequest) (*DeletePropertyResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeletePropertyResponse), args.Error(1)
}

func (p *Mock) CreateProperty(ctx context.Context, req CreatePropertyRequest) (*CreatePropertyResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreatePropertyResponse), args.Error(1)
}

func (p *Mock) UpdateProperty(ctx context.Context, req UpdatePropertyRequest) (*UpdatePropertyResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdatePropertyResponse), args.Error(1)
}

func (p *Mock) ListProperties(ctx context.Context, req ListPropertiesRequest) ([]Property, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Property), args.Error(1)
}

func (p *Mock) GetDatacenter(ctx context.Context, req GetDatacenterRequest) (*Datacenter, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Datacenter), args.Error(1)
}

func (p *Mock) CreateDatacenter(ctx context.Context, req CreateDatacenterRequest) (*CreateDatacenterResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateDatacenterResponse), args.Error(1)
}

func (p *Mock) DeleteDatacenter(ctx context.Context, req DeleteDatacenterRequest) (*DeleteDatacenterResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteDatacenterResponse), args.Error(1)
}

func (p *Mock) UpdateDatacenter(ctx context.Context, req UpdateDatacenterRequest) (*UpdateDatacenterResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateDatacenterResponse), args.Error(1)
}

func (p *Mock) ListDatacenters(ctx context.Context, req ListDatacentersRequest) ([]Datacenter, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Datacenter), args.Error(1)
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

func (p *Mock) GetResource(ctx context.Context, req GetResourceRequest) (*GetResourceResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetResourceResponse), args.Error(1)
}

func (p *Mock) CreateResource(ctx context.Context, req CreateResourceRequest) (*CreateResourceResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateResourceResponse), args.Error(1)
}

func (p *Mock) DeleteResource(ctx context.Context, req DeleteResourceRequest) (*DeleteResourceResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteResourceResponse), args.Error(1)
}

func (p *Mock) UpdateResource(ctx context.Context, req UpdateResourceRequest) (*UpdateResourceResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateResourceResponse), args.Error(1)
}

func (p *Mock) ListResources(ctx context.Context, req ListResourcesRequest) ([]Resource, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Resource), args.Error(1)
}

func (p *Mock) GetASMap(ctx context.Context, req GetASMapRequest) (*GetASMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetASMapResponse), args.Error(1)
}

func (p *Mock) CreateASMap(ctx context.Context, req CreateASMapRequest) (*CreateASMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateASMapResponse), args.Error(1)
}

func (p *Mock) DeleteASMap(ctx context.Context, req DeleteASMapRequest) (*DeleteASMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteASMapResponse), args.Error(1)
}

func (p *Mock) UpdateASMap(ctx context.Context, req UpdateASMapRequest) (*UpdateASMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateASMapResponse), args.Error(1)
}

func (p *Mock) ListASMaps(ctx context.Context, req ListASMapsRequest) ([]ASMap, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]ASMap), args.Error(1)
}

func (p *Mock) GetGeoMap(ctx context.Context, req GetGeoMapRequest) (*GetGeoMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetGeoMapResponse), args.Error(1)
}

func (p *Mock) CreateGeoMap(ctx context.Context, req CreateGeoMapRequest) (*CreateGeoMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateGeoMapResponse), args.Error(1)
}

func (p *Mock) DeleteGeoMap(ctx context.Context, req DeleteGeoMapRequest) (*DeleteGeoMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteGeoMapResponse), args.Error(1)
}

func (p *Mock) UpdateGeoMap(ctx context.Context, req UpdateGeoMapRequest) (*UpdateGeoMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateGeoMapResponse), args.Error(1)
}

func (p *Mock) ListGeoMaps(ctx context.Context, req ListGeoMapsRequest) ([]GeoMap, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]GeoMap), args.Error(1)
}

func (p *Mock) GetCIDRMap(ctx context.Context, req GetCIDRMapRequest) (*GetCIDRMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCIDRMapResponse), args.Error(1)
}

func (p *Mock) CreateCIDRMap(ctx context.Context, req CreateCIDRMapRequest) (*CreateCIDRMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateCIDRMapResponse), args.Error(1)
}

func (p *Mock) DeleteCIDRMap(ctx context.Context, req DeleteCIDRMapRequest) (*DeleteCIDRMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteCIDRMapResponse), args.Error(1)
}

func (p *Mock) UpdateCIDRMap(ctx context.Context, req UpdateCIDRMapRequest) (*UpdateCIDRMapResponse, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateCIDRMapResponse), args.Error(1)
}

func (p *Mock) ListCIDRMaps(ctx context.Context, req ListCIDRMapsRequest) ([]CIDRMap, error) {
	args := p.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]CIDRMap), args.Error(1)
}
