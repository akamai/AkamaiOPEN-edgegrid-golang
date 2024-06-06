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

func (d *Mock) ListZones(ctx context.Context, query ...ZoneListQueryArgs) (*ZoneListResponse, error) {
	var args mock.Arguments

	if len(query) > 0 {
		args = d.Called(ctx, query[0])
	} else {
		args = d.Called(ctx)
	}

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

func (d *Mock) GetZone(ctx context.Context, name string) (*ZoneResponse, error) {
	args := d.Called(ctx, name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ZoneResponse), args.Error(1)
}

func (d *Mock) GetChangeList(ctx context.Context, param string) (*ChangeListResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ChangeListResponse), args.Error(1)
}

func (d *Mock) GetMasterZoneFile(ctx context.Context, param string) (string, error) {
	args := d.Called(ctx, param)

	return args.String(0), args.Error(1)
}

func (d *Mock) CreateZone(ctx context.Context, param1 *ZoneCreate, param2 ZoneQueryString, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param1, param2, param3[0])
	} else {
		args = d.Called(ctx, param1, param2)
	}

	return args.Error(0)
}

func (d *Mock) SaveChangelist(ctx context.Context, param *ZoneCreate) error {
	args := d.Called(ctx, param)

	return args.Error(0)
}

func (d *Mock) SubmitChangelist(ctx context.Context, param *ZoneCreate) error {
	args := d.Called(ctx, param)

	return args.Error(0)
}

func (d *Mock) UpdateZone(ctx context.Context, param1 *ZoneCreate, param2 ZoneQueryString) error {
	args := d.Called(ctx, param1, param2)

	return args.Error(0)
}

func (d *Mock) DeleteZone(ctx context.Context, param1 *ZoneCreate, param2 ZoneQueryString) error {
	args := d.Called(ctx, param1, param2)

	return args.Error(0)
}

func (d *Mock) GetZoneNames(ctx context.Context, param string) (*ZoneNamesResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ZoneNamesResponse), args.Error(1)
}

func (d *Mock) GetZoneNameTypes(ctx context.Context, param1 string, param2 string) (*ZoneNameTypesResponse, error) {
	args := d.Called(ctx, param1, param2)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ZoneNameTypesResponse), args.Error(1)
}

func (d *Mock) ListTSIGKeys(ctx context.Context, param *TSIGQueryString) (*TSIGReportResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*TSIGReportResponse), args.Error(1)
}

func (d *Mock) GetTSIGKeyZones(ctx context.Context, param *TSIGKey) (*ZoneNameListResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ZoneNameListResponse), args.Error(1)
}

func (d *Mock) GetTSIGKeyAliases(ctx context.Context, param string) (*ZoneNameListResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ZoneNameListResponse), args.Error(1)
}

func (d *Mock) TSIGKeyBulkUpdate(ctx context.Context, param1 *TSIGKeyBulkPost) error {
	args := d.Called(ctx, param1)

	return args.Error(0)
}

func (d *Mock) GetTSIGKey(ctx context.Context, param string) (*TSIGKeyResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*TSIGKeyResponse), args.Error(1)
}

func (d *Mock) DeleteTSIGKey(ctx context.Context, param1 string) error {
	args := d.Called(ctx, param1)

	return args.Error(0)
}

func (d *Mock) UpdateTSIGKey(ctx context.Context, param1 *TSIGKey, param2 string) error {
	args := d.Called(ctx, param1, param2)

	return args.Error(0)
}

func (d *Mock) GetAuthorities(ctx context.Context, param string) (*AuthorityResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*AuthorityResponse), args.Error(1)
}

func (d *Mock) GetNameServerRecordList(ctx context.Context, param string) ([]string, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (d *Mock) GetRecordList(ctx context.Context, param string, param2 string, param3 string) (*RecordSetResponse, error) {
	args := d.Called(ctx, param, param2, param3)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RecordSetResponse), args.Error(1)
}

func (d *Mock) GetRdata(ctx context.Context, param string, param2 string, param3 string) ([]string, error) {
	args := d.Called(ctx, param, param2, param3)

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

func (d *Mock) GetRecord(ctx context.Context, param string, param2 string, param3 string) (*RecordBody, error) {
	args := d.Called(ctx, param, param2, param3)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RecordBody), args.Error(1)
}

func (d *Mock) CreateRecord(ctx context.Context, param *RecordBody, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *Mock) DeleteRecord(ctx context.Context, param *RecordBody, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *Mock) UpdateRecord(ctx context.Context, param *RecordBody, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *Mock) GetRecordSets(ctx context.Context, param string, param2 ...RecordSetQueryArgs) (*RecordSetResponse, error) {
	var args mock.Arguments

	if len(param2) > 0 {
		args = d.Called(ctx, param, param2)
	} else {
		args = d.Called(ctx, param)
	}

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RecordSetResponse), args.Error(1)
}

func (d *Mock) CreateRecordSets(ctx context.Context, param *RecordSets, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *Mock) UpdateRecordSets(ctx context.Context, param *RecordSets, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *Mock) PostMasterZoneFile(ctx context.Context, param string, param2 string) error {
	args := d.Called(ctx, param, param2)

	return args.Error(0)
}
func (d *Mock) CreateBulkZones(ctx context.Context, param *BulkZonesCreate, param2 ZoneQueryString) (*BulkZonesResponse, error) {
	args := d.Called(ctx, param, param2)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BulkZonesResponse), args.Error(1)
}
func (d *Mock) DeleteBulkZones(ctx context.Context, param *ZoneNameListResponse, param2 ...bool) (*BulkZonesResponse, error) {
	var args mock.Arguments

	if len(param2) > 0 {
		args = d.Called(ctx, param, param2[0])
	} else {
		args = d.Called(ctx, param)
	}

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BulkZonesResponse), args.Error(1)
}
func (d *Mock) GetBulkZoneCreateStatus(ctx context.Context, param string) (*BulkStatusResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*BulkStatusResponse), args.Error(1)
}
func (d *Mock) GetBulkZoneDeleteStatus(ctx context.Context, param string) (*BulkStatusResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*BulkStatusResponse), args.Error(1)
}
func (d *Mock) GetBulkZoneCreateResult(ctx context.Context, param string) (*BulkCreateResultResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*BulkCreateResultResponse), args.Error(1)
}
func (d *Mock) GetBulkZoneDeleteResult(ctx context.Context, param string) (*BulkDeleteResultResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*BulkDeleteResultResponse), args.Error(1)
}

func (d *Mock) ListGroups(ctx context.Context, request ListGroupRequest) (*ListGroupResponse, error) {
	args := d.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListGroupResponse), args.Error(1)
}
