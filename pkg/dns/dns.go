// Package dns provides access to the Akamai DNS V2 APIs
//
// See: https://techdocs.akamai.com/edge-dns/reference/edge-dns-api
package dns

import (
	"context"
	"errors"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")
)

type (
	// DNS is the dns api interface
	DNS interface {
		// Authorities

		// GetAuthorities provides a list of structured read-only list of name servers.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-data-authorities
		GetAuthorities(context.Context, GetAuthoritiesRequest) (*GetAuthoritiesResponse, error)

		// GetNameServerRecordList provides a list of name server records.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-data-authorities
		GetNameServerRecordList(context.Context, GetNameServerRecordListRequest) ([]string, error)

		// ListGroups returns group list associated with particular user
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-data-groups
		ListGroups(context.Context, ListGroupRequest) (*ListGroupResponse, error)

		// Data

		// GetRdata retrieves record rdata, e.g. target.
		GetRdata(ctx context.Context, params GetRdataRequest) ([]string, error)
		// ProcessRdata process rdata.
		ProcessRdata(context.Context, []string, string) []string
		// ParseRData parses rdata. returning map.
		ParseRData(context.Context, string, []string) map[string]interface{}

		// Recordsets

		// GetRecordSets retrieves record sets with Query Args. No formatting of arg values.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-recordsets
		GetRecordSets(context.Context, GetRecordSetsRequest) (*GetRecordSetsResponse, error)
		// CreateRecordSets creates multiple record sets.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-zone-recordsets
		CreateRecordSets(context.Context, CreateRecordSetsRequest) error
		// UpdateRecordSets replaces list of record sets.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/put-zones-zone-recordsets
		UpdateRecordSets(context.Context, UpdateRecordSetsRequest) error

		// GetRecordList retrieves recordset list based on type.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-recordsets
		GetRecordList(context.Context, GetRecordListRequest) (*GetRecordListResponse, error)

		// Records

		// GetRecord retrieves a recordset and returns as RecordBody.
		//
		// See:  https://techdocs.akamai.com/edge-dns/reference/get-zone-name-type
		GetRecord(context.Context, GetRecordRequest) (*GetRecordResponse, error)
		// CreateRecord creates recordset.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-zone-names-name-types-type
		CreateRecord(context.Context, CreateRecordRequest) error
		// DeleteRecord removes recordset.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/delete-zone-name-type
		DeleteRecord(context.Context, DeleteRecordRequest) error
		// UpdateRecord replaces the recordset.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/put-zones-zone-names-name-types-type
		UpdateRecord(context.Context, UpdateRecordRequest) error

		// TSIGKeys

		// ListTSIGKeys lists the TSIG keys used by zones that you are allowed to manage.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-keys
		ListTSIGKeys(context.Context, ListTSIGKeysRequest) (*ListTSIGKeysResponse, error)
		// GetTSIGKeyZones retrieves DNS Zones using TSIG key.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-keys-used-by
		GetTSIGKeyZones(context.Context, GetTSIGKeyZonesRequest) (*GetTSIGKeyZonesResponse, error)
		// GetTSIGKeyAliases retrieves a DNS Zone's aliases.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-key-used-by
		GetTSIGKeyAliases(context.Context, GetTSIGKeyAliasesRequest) (*GetTSIGKeyAliasesResponse, error)
		// UpdateTSIGKeyBulk updates Bulk Zones TSIG key.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-keys-bulk-update
		UpdateTSIGKeyBulk(context.Context, UpdateTSIGKeyBulkRequest) error
		// GetTSIGKey retrieves a TSIG key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-key
		GetTSIGKey(context.Context, GetTSIGKeyRequest) (*GetTSIGKeyResponse, error)
		// DeleteTSIGKey deletes TSIG key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/delete-zones-zone-key
		DeleteTSIGKey(context.Context, DeleteTSIGKeyRequest) error
		// UpdateTSIGKey updates TSIG key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/put-zones-zone-key
		UpdateTSIGKey(context.Context, UpdateTSIGKeyRequest) error

		// Zones

		// ListZones retrieves a list of all zones user can access.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones
		ListZones(context.Context, ListZonesRequest) (*ZoneListResponse, error)

		// GetZone retrieves Zone metadata.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone
		GetZone(context.Context, GetZoneRequest) (*GetZoneResponse, error)
		//GetChangeList retrieves Zone changelist.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-changelists-zone
		GetChangeList(context.Context, GetChangeListRequest) (*GetChangeListResponse, error)
		// GetMasterZoneFile retrieves master zone file.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-zone-file
		GetMasterZoneFile(context.Context, GetMasterZoneFileRequest) (string, error)
		// PostMasterZoneFile updates master zone file.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-zone-zone-file
		PostMasterZoneFile(context.Context, PostMasterZoneFileRequest) error
		// CreateZone creates new zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zone
		CreateZone(context.Context, CreateZoneRequest) error
		// SaveChangeList creates a new Change List based on the most recent version of a zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-changelists
		SaveChangeList(context.Context, SaveChangeListRequest) error
		// SubmitChangeList submits changelist for the Zone to create default NS SOA records.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-changelists-zone-submit
		SubmitChangeList(context.Context, SubmitChangeListRequest) error
		// UpdateZone updates zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/put-zone
		UpdateZone(context.Context, UpdateZoneRequest) error

		// GetZoneNames retrieves a list of a zone's record names.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone-names
		GetZoneNames(context.Context, GetZoneNamesRequest) (*GetZoneNamesResponse, error)
		// GetZoneNameTypes retrieves a zone name's record types.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone-name-types
		GetZoneNameTypes(context.Context, GetZoneNameTypesRequest) (*GetZoneNameTypesResponse, error)
		// CreateBulkZones submits create bulk zone request.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-create-requests
		CreateBulkZones(context.Context, CreateBulkZonesRequest) (*CreateBulkZonesResponse, error)
		// DeleteBulkZones submits delete bulk zone request.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-delete-requests
		DeleteBulkZones(context.Context, DeleteBulkZonesRequest) (*DeleteBulkZonesResponse, error)
		// GetBulkZoneCreateStatus retrieves submit request status.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-create-requests-requestid
		GetBulkZoneCreateStatus(context.Context, GetBulkZoneCreateStatusRequest) (*GetBulkZoneCreateStatusResponse, error)
		//GetBulkZoneDeleteStatus retrieves submit request status.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-delete-requests-requestid
		GetBulkZoneDeleteStatus(context.Context, GetBulkZoneDeleteStatusRequest) (*GetBulkZoneDeleteStatusResponse, error)
		// GetBulkZoneCreateResult retrieves create request result.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-create-requests-requestid-result
		GetBulkZoneCreateResult(ctx context.Context, request GetBulkZoneCreateResultRequest) (*GetBulkZoneCreateResultResponse, error)
		// GetBulkZoneDeleteResult retrieves delete request result.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-delete-requests-requestid-result
		GetBulkZoneDeleteResult(context.Context, GetBulkZoneDeleteResultRequest) (*GetBulkZoneDeleteResultResponse, error)
		// GetZonesDNSSecStatus returns the current DNSSEC status for one or more zones.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-dns-sec-status
		GetZonesDNSSecStatus(context.Context, GetZonesDNSSecStatusRequest) (*GetZonesDNSSecStatusResponse, error)
	}

	dns struct {
		session.Session
	}

	// Option defines a DNS option
	Option func(*dns)

	// ClientFunc is a dns client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) DNS
)

// Client returns a new dns Client instance with the specified controller
func Client(sess session.Session, opts ...Option) DNS {
	d := &dns{
		Session: sess,
	}

	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Exec overrides the session.Exec to add dns options
func (d *dns) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	return d.Session.Exec(r, out, in...)
}
