//revive:disable:exported

package dns

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ DNS = &Mock{}

func (d *Mock) ListZones(ctx context.Context, req ListZonesRequest) (*ZoneListResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ZoneListResponse), args.Error(1)
}

func (d *Mock) GetZonesDNSSecStatus(ctx context.Context, params GetZonesDNSSecStatusRequest) (*GetZonesDNSSecStatusResponse, error) {
	args := d.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetZonesDNSSecStatusResponse), args.Error(1)
}

func (d *Mock) GetZone(ctx context.Context, req GetZoneRequest) (*GetZoneResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetZoneResponse), args.Error(1)
}

func (d *Mock) GetChangeList(ctx context.Context, req GetChangeListRequest) (*GetChangeListResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetChangeListResponse), args.Error(1)
}

func (d *Mock) GetMasterZoneFile(ctx context.Context, req GetMasterZoneFileRequest) (string, error) {
	args := d.Called(ctx, req)

	return args.String(0), args.Error(1)
}

func (d *Mock) CreateZone(ctx context.Context, req CreateZoneRequest) error {
	args := d.Called(ctx, req)

	return args.Error(0)
}

func (d *Mock) SaveChangeList(ctx context.Context, req SaveChangeListRequest) error {
	args := d.Called(ctx, req)

	return args.Error(0)
}

func (d *Mock) SubmitChangeList(ctx context.Context, req SubmitChangeListRequest) error {
	args := d.Called(ctx, req)

	return args.Error(0)
}

func (d *Mock) UpdateZone(ctx context.Context, req UpdateZoneRequest) error {
	args := d.Called(ctx, req)

	return args.Error(0)
}

func (d *Mock) GetZoneNames(ctx context.Context, req GetZoneNamesRequest) (*GetZoneNamesResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetZoneNamesResponse), args.Error(1)
}

func (d *Mock) GetZoneNameTypes(ctx context.Context, req GetZoneNameTypesRequest) (*GetZoneNameTypesResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetZoneNameTypesResponse), args.Error(1)
}

func (d *Mock) ListTSIGKeys(ctx context.Context, req ListTSIGKeysRequest) (*ListTSIGKeysResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListTSIGKeysResponse), args.Error(1)
}

func (d *Mock) GetTSIGKeyZones(ctx context.Context, req GetTSIGKeyZonesRequest) (*GetTSIGKeyZonesResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetTSIGKeyZonesResponse), args.Error(1)
}

func (d *Mock) GetTSIGKeyAliases(ctx context.Context, req GetTSIGKeyAliasesRequest) (*GetTSIGKeyAliasesResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetTSIGKeyAliasesResponse), args.Error(1)
}

func (d *Mock) UpdateTSIGKeyBulk(ctx context.Context, req UpdateTSIGKeyBulkRequest) error {
	args := d.Called(ctx, req)

	return args.Error(0)
}

func (d *Mock) GetTSIGKey(ctx context.Context, req GetTSIGKeyRequest) (*GetTSIGKeyResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetTSIGKeyResponse), args.Error(1)
}

func (d *Mock) DeleteTSIGKey(ctx context.Context, req DeleteTSIGKeyRequest) error {
	args := d.Called(ctx, req)

	return args.Error(0)
}

func (d *Mock) UpdateTSIGKey(ctx context.Context, req UpdateTSIGKeyRequest) error {
	args := d.Called(ctx, req)

	return args.Error(0)
}

func (d *Mock) GetAuthorities(ctx context.Context, req GetAuthoritiesRequest) (*GetAuthoritiesResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetAuthoritiesResponse), args.Error(1)
}

func (d *Mock) GetNameServerRecordList(ctx context.Context, req GetNameServerRecordListRequest) ([]string, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (d *Mock) GetRecordList(ctx context.Context, req GetRecordListRequest) (*GetRecordListResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetRecordListResponse), args.Error(1)
}

func (d *Mock) GetRdata(ctx context.Context, req GetRdataRequest) ([]string, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (d *Mock) ProcessRdata(ctx context.Context, param []string, param2 string) []string {
	args := d.Called(ctx, param, param2)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

func (d *Mock) ParseRData(ctx context.Context, param string, param2 []string) map[string]interface{} {
	args := d.Called(ctx, param, param2)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]interface{})
}

func (d *Mock) GetRecord(ctx context.Context, req GetRecordRequest) (*GetRecordResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetRecordResponse), args.Error(1)
}

func (d *Mock) CreateRecord(ctx context.Context, req CreateRecordRequest) error {
	args := d.Called(ctx, req)
	return args.Error(0)
}

func (d *Mock) DeleteRecord(ctx context.Context, req DeleteRecordRequest) error {
	args := d.Called(ctx, req)
	return args.Error(0)
}

func (d *Mock) UpdateRecord(ctx context.Context, req UpdateRecordRequest) error {
	args := d.Called(ctx, req)
	return args.Error(0)
}

func (d *Mock) GetRecordSets(ctx context.Context, req GetRecordSetsRequest) (*GetRecordSetsResponse, error) {
	args := d.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetRecordSetsResponse), args.Error(1)
}

func (d *Mock) CreateRecordSets(ctx context.Context, req CreateRecordSetsRequest) error {
	args := d.Called(ctx, req)
	return args.Error(0)
}

func (d *Mock) UpdateRecordSets(ctx context.Context, req UpdateRecordSetsRequest) error {
	args := d.Called(ctx, req)
	return args.Error(0)
}

func (d *Mock) PostMasterZoneFile(ctx context.Context, req PostMasterZoneFileRequest) error {
	args := d.Called(ctx, req)
	return args.Error(0)
}
func (d *Mock) CreateBulkZones(ctx context.Context, req CreateBulkZonesRequest) (*CreateBulkZonesResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateBulkZonesResponse), args.Error(1)
}
func (d *Mock) DeleteBulkZones(ctx context.Context, req DeleteBulkZonesRequest) (*DeleteBulkZonesResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DeleteBulkZonesResponse), args.Error(1)
}
func (d *Mock) GetBulkZoneCreateStatus(ctx context.Context, req GetBulkZoneCreateStatusRequest) (*GetBulkZoneCreateStatusResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetBulkZoneCreateStatusResponse), args.Error(1)
}
func (d *Mock) GetBulkZoneDeleteStatus(ctx context.Context, req GetBulkZoneDeleteStatusRequest) (*GetBulkZoneDeleteStatusResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetBulkZoneDeleteStatusResponse), args.Error(1)
}
func (d *Mock) GetBulkZoneCreateResult(ctx context.Context, req GetBulkZoneCreateResultRequest) (*GetBulkZoneCreateResultResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetBulkZoneCreateResultResponse), args.Error(1)
}
func (d *Mock) GetBulkZoneDeleteResult(ctx context.Context, req GetBulkZoneDeleteResultRequest) (*GetBulkZoneDeleteResultResponse, error) {
	args := d.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetBulkZoneDeleteResultResponse), args.Error(1)
}

func (d *Mock) ListGroups(ctx context.Context, request ListGroupRequest) (*ListGroupResponse, error) {
	args := d.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListGroupResponse), args.Error(1)
}
