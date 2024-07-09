package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	zoneWriteLock sync.Mutex
)

type (
	// Zones contains operations available on Zone resources.
	Zones interface {
		// ListZones retrieves a list of all zones user can access.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones
		ListZones(context.Context, ...ZoneListQueryArgs) (*ZoneListResponse, error)
		// GetZone retrieves Zone metadata.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone
		GetZone(context.Context, string) (*ZoneResponse, error)
		//GetChangeList retrieves Zone changelist.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-changelists-zone
		GetChangeList(context.Context, string) (*ChangeListResponse, error)
		// GetMasterZoneFile retrieves master zone file.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-zone-file
		GetMasterZoneFile(context.Context, string) (string, error)
		// PostMasterZoneFile updates master zone file.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-zone-zone-file
		PostMasterZoneFile(context.Context, string, string) error
		// CreateZone creates new zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zone
		CreateZone(context.Context, *ZoneCreate, ZoneQueryString, ...bool) error
		// SaveChangelist creates a new Change List based on the most recent version of a zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-changelists
		SaveChangelist(context.Context, *ZoneCreate) error
		// SubmitChangelist submits changelist for the Zone to create default NS SOA records.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-changelists-zone-submit
		SubmitChangelist(context.Context, *ZoneCreate) error
		// UpdateZone updates zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/put-zone
		UpdateZone(context.Context, *ZoneCreate, ZoneQueryString) error
		// GetZoneNames retrieves a list of a zone's record names.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone-names
		GetZoneNames(context.Context, string) (*ZoneNamesResponse, error)
		// GetZoneNameTypes retrieves a zone name's record types.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone-name-types
		GetZoneNameTypes(context.Context, string, string) (*ZoneNameTypesResponse, error)
		// CreateBulkZones submits create bulk zone request.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-create-requests
		CreateBulkZones(context.Context, *BulkZonesCreate, ZoneQueryString) (*BulkZonesResponse, error)
		// DeleteBulkZones submits delete bulk zone request.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-delete-requests
		DeleteBulkZones(context.Context, *ZoneNameListResponse, ...bool) (*BulkZonesResponse, error)
		// GetBulkZoneCreateStatus retrieves submit request status.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-create-requests-requestid
		GetBulkZoneCreateStatus(context.Context, string) (*BulkStatusResponse, error)
		//GetBulkZoneDeleteStatus retrieves submit request status.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-delete-requests-requestid
		GetBulkZoneDeleteStatus(context.Context, string) (*BulkStatusResponse, error)
		// GetBulkZoneCreateResult retrieves create request result.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-create-requests-requestid-result
		GetBulkZoneCreateResult(ctx context.Context, requestid string) (*BulkCreateResultResponse, error)
		// GetBulkZoneDeleteResult retrieves delete request result.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-delete-requests-requestid-result
		GetBulkZoneDeleteResult(context.Context, string) (*BulkDeleteResultResponse, error)
		// GetZonesDNSSecStatus returns the current DNSSEC status for one or more zones.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-dns-sec-status
		GetZonesDNSSecStatus(context.Context, GetZonesDNSSecStatusRequest) (*GetZonesDNSSecStatusResponse, error)
	}

	// ZoneQueryString contains zone query parameters
	ZoneQueryString struct {
		Contract string
		Group    string
	}

	// ZoneCreate contains zone create request
	ZoneCreate struct {
		Zone                  string   `json:"zone"`
		Type                  string   `json:"type"`
		Masters               []string `json:"masters,omitempty"`
		Comment               string   `json:"comment,omitempty"`
		SignAndServe          bool     `json:"signAndServe"`
		SignAndServeAlgorithm string   `json:"signAndServeAlgorithm,omitempty"`
		TSIGKey               *TSIGKey `json:"tsigKey,omitempty"`
		Target                string   `json:"target,omitempty"`
		EndCustomerID         string   `json:"endCustomerId,omitempty"`
		ContractID            string   `json:"contractId,omitempty"`
	}

	// ZoneResponse contains zone create response
	ZoneResponse struct {
		Zone                  string   `json:"zone,omitempty"`
		Type                  string   `json:"type,omitempty"`
		Masters               []string `json:"masters,omitempty"`
		Comment               string   `json:"comment,omitempty"`
		SignAndServe          bool     `json:"signAndServe"`
		SignAndServeAlgorithm string   `json:"signAndServeAlgorithm,omitempty"`
		TSIGKey               *TSIGKey `json:"tsigKey,omitempty"`
		Target                string   `json:"target,omitempty"`
		EndCustomerID         string   `json:"endCustomerId,omitempty"`
		ContractID            string   `json:"contractId,omitempty"`
		AliasCount            int64    `json:"aliasCount,omitempty"`
		ActivationState       string   `json:"activationState,omitempty"`
		LastActivationDate    string   `json:"lastActivationDate,omitempty"`
		LastModifiedBy        string   `json:"lastModifiedBy,omitempty"`
		LastModifiedDate      string   `json:"lastModifiedDate,omitempty"`
		VersionID             string   `json:"versionId,omitempty"`
	}

	// ZoneListQueryArgs contains parameters for List Zones query
	ZoneListQueryArgs struct {
		ContractIDs string
		Page        int
		PageSize    int
		Search      string
		ShowAll     bool
		SortBy      string
		Types       string
	}

	// ListMetadata contains metadata for List Zones request
	ListMetadata struct {
		ContractIDs   []string `json:"contractIds"`
		Page          int      `json:"page"`
		PageSize      int      `json:"pageSize"`
		ShowAll       bool     `json:"showAll"`
		TotalElements int      `json:"totalElements"`
	}

	// ZoneListResponse contains response for List Zones request
	ZoneListResponse struct {
		Metadata *ListMetadata   `json:"metadata,omitempty"`
		Zones    []*ZoneResponse `json:"zones,omitempty"`
	}

	// ChangeListResponse contains metadata about a change list
	ChangeListResponse struct {
		Zone             string `json:"zone,omitempty"`
		ChangeTag        string `json:"changeTag,omitempty"`
		ZoneVersionID    string `json:"zoneVersionId,omitempty"`
		LastModifiedDate string `json:"lastModifiedDate,omitempty"`
		Stale            bool   `json:"stale,omitempty"`
	}

	// ZoneNameListResponse contains response with a list of zone's names and aliases
	ZoneNameListResponse struct {
		Zones   []string `json:"zones"`
		Aliases []string `json:"aliases,omitempty"`
	}

	// ZoneNamesResponse contains record set names for zone
	ZoneNamesResponse struct {
		Names []string `json:"names"`
	}

	// ZoneNameTypesResponse contains record set types for zone
	ZoneNameTypesResponse struct {
		Types []string `json:"types"`
	}

	// GetZonesDNSSecStatusRequest is used to get the DNSSEC status for one or more zones
	GetZonesDNSSecStatusRequest struct {
		Zones []string `json:"zones"`
	}

	// GetZonesDNSSecStatusResponse represents a list of DNSSEC statuses for DNS zones specified
	// in the GetZonesDNSSecStatus request
	GetZonesDNSSecStatusResponse struct {
		DNSSecStatuses []SecStatus `json:"dnsSecStatuses"`
	}

	// SecStatus represents the DNSSEC status for a DNS zone
	SecStatus struct {
		Zone           string      `json:"zone"`
		Alerts         []string    `json:"alerts"`
		CurrentRecords SecRecords  `json:"currentRecords"`
		NewRecords     *SecRecords `json:"newRecords"`
	}

	// SecRecords represents a set of DNSSEC records for a DNS zone
	SecRecords struct {
		DNSKeyRecord     string    `json:"dnskeyRecord"`
		DSRecord         string    `json:"dsRecord"`
		ExpectedTTL      int64     `json:"expectedTtl"`
		LastModifiedDate time.Time `json:"lastModifiedDate"`
	}
)

// Validate validates GetZonesDNSSecStatusRequest
func (r GetZonesDNSSecStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zones": validation.Validate(r.Zones, validation.Required),
	})
}

var zoneStructMap = map[string]string{
	"Zone":                  "zone",
	"Type":                  "type",
	"Masters":               "masters",
	"Comment":               "comment",
	"SignAndServe":          "signAndServe",
	"SignAndServeAlgorithm": "signAndServeAlgorithm",
	"TSIGKey":               "tsigKey",
	"Target":                "target",
	"EndCustomerID":         "endCustomerId",
	"ContractId":            "contractId"}

// Util to convert struct to http request body, eg. io.reader
func convertStructToReqBody(srcStruct interface{}) (io.Reader, error) {
	reqBody, err := json.Marshal(srcStruct)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(reqBody), nil
}

func (d *dns) ListZones(ctx context.Context, queryArgs ...ZoneListQueryArgs) (*ZoneListResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("ListZones")

	getURL := fmt.Sprintf("/config-dns/v2/zones")
	if len(queryArgs) > 1 {
		return nil, fmt.Errorf("ListZones QueryArgs invalid")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create listzones request: %w", err)
	}

	q := req.URL.Query()
	if len(queryArgs) > 0 {
		if queryArgs[0].Page > 0 {
			q.Add("page", strconv.Itoa(queryArgs[0].Page))
		}
		if queryArgs[0].PageSize > 0 {
			q.Add("pageSize", strconv.Itoa(queryArgs[0].PageSize))
		}
		if queryArgs[0].Search != "" {
			q.Add("search", queryArgs[0].Search)
		}
		q.Add("showAll", strconv.FormatBool(queryArgs[0].ShowAll))
		if queryArgs[0].SortBy != "" {
			q.Add("sortBy", queryArgs[0].SortBy)
		}
		if queryArgs[0].Types != "" {
			q.Add("types", queryArgs[0].Types)
		}
		if queryArgs[0].ContractIDs != "" {
			q.Add("contractIds", queryArgs[0].ContractIDs)
		}
		req.URL.RawQuery = q.Encode()
	}

	var result ZoneListResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("listzones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetZone(ctx context.Context, zoneName string) (*ZoneResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetZone")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s", zoneName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZone request: %w", err)
	}

	var result ZoneResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetZone request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetChangeList(ctx context.Context, zone string) (*ChangeListResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetChangeList")

	getURL := fmt.Sprintf("/config-dns/v2/changelists/%s", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetChangeList request: %w", err)
	}

	var result ChangeListResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetChangeList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetMasterZoneFile(ctx context.Context, zone string) (string, error) {
	logger := d.Log(ctx)
	logger.Debug("GetMasterZoneFile")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/zone-file", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GetMasterZoneFile request: %w", err)
	}
	req.Header.Add("Accept", "text/dns")

	resp, err := d.Exec(req, nil)
	if err != nil {
		return "", fmt.Errorf("GetMasterZoneFile request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", d.Error(resp)
	}

	masterFile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("GetMasterZoneFile request failed: %w", err)
	}

	return string(masterFile), nil
}

func (d *dns) PostMasterZoneFile(ctx context.Context, zone string, fileData string) error {
	logger := d.Log(ctx)
	logger.Debug("PostMasterZoneFile")

	mtResp := ""
	pmzfURL := fmt.Sprintf("/config-dns/v2/zones/%s/zone-file", zone)
	buf := bytes.NewReader([]byte(fileData))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pmzfURL, buf)
	if err != nil {
		return fmt.Errorf("failed to create PostMasterZoneFile request: %w", err)
	}

	req.Header.Set("Content-Type", "text/dns")

	resp, err := d.Exec(req, &mtResp)
	if err != nil {
		return fmt.Errorf("Create PostMasterZoneFile failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) CreateZone(ctx context.Context, zone *ZoneCreate, zoneQueryString ZoneQueryString, clearConn ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone,
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := d.Log(ctx)
	logger.Debug("Zone Create")

	if err := ValidateZone(zone); err != nil {
		return err
	}

	zoneMap := filterZoneCreate(zone)

	var zoneResponse ZoneResponse
	zoneURL := "/config-dns/v2/zones/?contractId=" + zoneQueryString.Contract
	if len(zoneQueryString.Group) > 0 {
		zoneURL += "&gid=" + zoneQueryString.Group
	}

	reqBody, err := convertStructToReqBody(zoneMap)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, zoneURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create Zone Create request: %w", err)
	}

	resp, err := d.Exec(req, &zoneResponse)
	if err != nil {
		return fmt.Errorf("Create Zone request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return d.Error(resp)
	}

	if strings.ToUpper(zone.Type) == "PRIMARY" {
		// Timing issue with Create immediately followed by SaveChangelist
		for _, clear := range clearConn {
			// should only be one entry
			if clear {
				logger.Info("Clearing Idle Connections")
				d.Client().CloseIdleConnections()
			}
		}
	}

	return nil
}

func (d *dns) SaveChangelist(ctx context.Context, zone *ZoneCreate) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := d.Log(ctx)
	logger.Debug("SaveChangeList")

	reqBody, err := convertStructToReqBody("")
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/changelists/?zone=%s", zone.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create SaveChangeList request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("SaveChangeList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) SubmitChangelist(ctx context.Context, zone *ZoneCreate) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := d.Log(ctx)
	logger.Debug("SubmitChangeList")

	reqBody, err := convertStructToReqBody("")
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/changelists/%s/submit", zone.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create SubmitChangeList request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("SubmitChangeList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) UpdateZone(ctx context.Context, zone *ZoneCreate, _ ZoneQueryString) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := d.Log(ctx)
	logger.Debug("Zone Update")

	if err := ValidateZone(zone); err != nil {
		return err
	}

	zoneMap := filterZoneCreate(zone)
	reqBody, err := convertStructToReqBody(zoneMap)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s", zone.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create Get Update request: %w", err)
	}

	var result ZoneResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("Zone Update request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return d.Error(resp)
	}

	return nil
}

func filterZoneCreate(zone *ZoneCreate) map[string]interface{} {
	zoneType := strings.ToUpper(zone.Type)
	filteredZone := make(map[string]interface{})
	zoneElems := reflect.ValueOf(zone).Elem()
	for i := 0; i < zoneElems.NumField(); i++ {
		varName := zoneElems.Type().Field(i).Name
		varLower := zoneStructMap[varName]
		varValue := zoneElems.Field(i).Interface()
		switch varName {
		case "Target":
			if zoneType == "ALIAS" {
				filteredZone[varLower] = varValue
			}
		case "TsigKey":
			if zoneType == "SECONDARY" {
				filteredZone[varLower] = varValue
			}
		case "Masters":
			if zoneType == "SECONDARY" {
				filteredZone[varLower] = varValue
			}
		case "SignAndServe":
			if zoneType != "ALIAS" {
				filteredZone[varLower] = varValue
			}
		case "SignAndServeAlgorithm":
			if zoneType != "ALIAS" {
				filteredZone[varLower] = varValue
			}
		default:
			filteredZone[varLower] = varValue
		}
	}

	return filteredZone
}

// ValidateZone validates ZoneCreate Object
func ValidateZone(zone *ZoneCreate) error {
	if len(zone.Zone) == 0 {
		return fmt.Errorf("Zone name is required")
	}
	zType := strings.ToUpper(zone.Type)
	if zType != "PRIMARY" && zType != "SECONDARY" && zType != "ALIAS" {
		return fmt.Errorf("Invalid zone type")
	}
	if zType != "SECONDARY" && zone.TSIGKey != nil {
		return fmt.Errorf("TsigKey is invalid for %s zone type", zType)
	}
	if zType == "ALIAS" {
		if len(zone.Target) == 0 {
			return fmt.Errorf("Target is required for Alias zone type")
		}
		if zone.Masters != nil && len(zone.Masters) > 0 {
			return fmt.Errorf("Masters is invalid for Alias zone type")
		}
		if zone.SignAndServe {
			return fmt.Errorf("SignAndServe is invalid for Alias zone type")
		}
		if len(zone.SignAndServeAlgorithm) > 0 {
			return fmt.Errorf("SignAndServeAlgorithm is invalid for Alias zone type")
		}
		return nil
	}
	// Primary or Secondary
	if len(zone.Target) > 0 {
		return fmt.Errorf("Target is invalid for %s zone type", zType)
	}
	if zone.Masters != nil && len(zone.Masters) > 0 && zType == "PRIMARY" {
		return fmt.Errorf("Masters is invalid for Primary zone type")
	}

	return nil
}

func (d *dns) GetZoneNames(ctx context.Context, zone string) (*ZoneNamesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetZoneNames")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/names", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZoneNames request: %w", err)
	}

	var result ZoneNamesResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetZoneNames request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetZoneNameTypes(ctx context.Context, zName, zone string) (*ZoneNameTypesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug(" GetZoneNameTypes")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types", zone, zName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZoneNameTypes request: %w", err)
	}

	var result ZoneNameTypesResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetZoneNameTypes request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetZonesDNSSecStatus(ctx context.Context, params GetZonesDNSSecStatusRequest) (*GetZonesDNSSecStatusResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetZonesDNSSecStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/config-dns/v2/zones/dns-sec-status", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZonesDNSSecStatus request: %w", err)
	}

	var result GetZonesDNSSecStatusResponse
	resp, err := d.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("GetZonesDNSSecStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}
