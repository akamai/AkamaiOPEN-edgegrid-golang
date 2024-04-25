package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
		Metadata *ListMetadata  `json:"metadata,omitempty"`
		Zones    []ZoneResponse `json:"zones,omitempty"`
	}

	// ZoneRequest contains request parameters
	ZoneRequest struct {
		Zone string
	}

	// GetZoneResponse contains the response data from GetZone operation
	GetZoneResponse ZoneResponse

	// GetChangeListResponse contains metadata about a change list
	GetChangeListResponse struct {
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

	// GetZoneNamesResponse contains record set names for zone
	GetZoneNamesResponse struct {
		Names []string `json:"names"`
	}

	// GetZoneNameTypesResponse contains record set types for zone
	GetZoneNameTypesResponse struct {
		Types []string `json:"types"`
	}
	// GetZoneRequest contains request parameters for GetZone
	GetZoneRequest ZoneRequest

	// GetChangeListRequest contains request parameters for GetChangeList
	GetChangeListRequest ZoneRequest

	// ListZonesRequest contains request parameters for ListZones
	ListZonesRequest struct {
		ContractIDs string
		Page        int
		PageSize    int
		Search      string
		ShowAll     bool
		SortBy      string
		Types       string
	}
	// GetMasterZoneFileRequest contains request parameters for GetMasterZoneFile
	GetMasterZoneFileRequest ZoneRequest

	// PostMasterZoneFileRequest contains request parameters for PostMasterZoneFile
	PostMasterZoneFileRequest struct {
		Zone     string
		FileData string
	}
	// CreateZoneRequest contains request parameters for CreateZone
	CreateZoneRequest struct {
		CreateZone      *ZoneCreate
		ZoneQueryString ZoneQueryString
		ClearConn       []bool
	}
	// SaveChangeListRequest contains request parameters for SaveChangelist
	SaveChangeListRequest ZoneCreate

	// SubmitChangeListRequest contains request parameters for SubmitChangeList
	SubmitChangeListRequest ZoneCreate

	// UpdateZoneRequest contains request parameters for UpdateZone
	UpdateZoneRequest struct {
		CreateZone *ZoneCreate
	}
	// GetZoneNamesRequest contains request parameters for GetZoneNames
	GetZoneNamesRequest ZoneRequest

	// GetZoneNameTypesRequest contains request parameters for GetZoneNameTypes
	GetZoneNameTypesRequest struct {
		Zone     string
		ZoneName string
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

var (
	// ErrGetZone is returned when GetZone fails
	ErrGetZone = errors.New("get zone")
	// ErrGetChangeList is returned when GetChangeList fails
	ErrGetChangeList = errors.New("get change list")
	// ErrGetMasterZoneFile is returned when GetMasterZoneFile fails
	ErrGetMasterZoneFile = errors.New("get master zone file")
	// ErrPostMasterZoneFile is returned when PostMasterZoneFile fails
	ErrPostMasterZoneFile = errors.New("post master zone file")
	// ErrCreateZone is returned when CreateZone fails
	ErrCreateZone = errors.New("create zone")
	// ErrSaveChangeList is returned when SaveChangeList fails
	ErrSaveChangeList = errors.New("save change list")
	// ErrSubmitChangeList is returned when SubmitChangeList fails
	ErrSubmitChangeList = errors.New("submit change list")
	// ErrGetZoneNames is returned when GetZoneNames fails
	ErrGetZoneNames = errors.New("get zone names")
	// ErrGetZoneNameTypes is returned when GetZoneNameTypes fails
	ErrGetZoneNameTypes = errors.New("get zone name types")
)

// Validate validates GetZoneNameTypesRequest
func (r GetZoneNameTypesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone":     validation.Validate(r.Zone, validation.Required),
		"ZoneName": validation.Validate(r.ZoneName, validation.Required),
	})
}

// Validate validates GetZoneNamesRequest
func (r GetZoneNamesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates SubmitChangeListRequest
func (r SubmitChangeListRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates SaveChangelistRequest
func (r SaveChangeListRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates PostMasterZoneFileRequest
func (r PostMasterZoneFileRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates CreateZoneRequest
func (r CreateZoneRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ZoneQueryString": validation.Validate(r.ZoneQueryString, validation.Required),
	})
}

// Validate validates GetZoneRequest
func (r GetZoneRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates GetMasterZoneFileRequest
func (r GetMasterZoneFileRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates GetChangeListRequest
func (r GetChangeListRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

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

func (d *dns) ListZones(ctx context.Context, params ListZonesRequest) (*ZoneListResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("ListZones")

	getURL := fmt.Sprintf("/config-dns/v2/zones")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create listzones request: %w", err)
	}

	q := req.URL.Query()
	if params.Page > 0 {
		q.Add("page", strconv.Itoa(params.Page))
	}
	if params.PageSize > 0 {
		q.Add("pageSize", strconv.Itoa(params.PageSize))
	}
	if params.Search != "" {
		q.Add("search", params.Search)
	}
	q.Add("showAll", strconv.FormatBool(params.ShowAll))
	if params.SortBy != "" {
		q.Add("sortBy", params.SortBy)
	}
	if params.Types != "" {
		q.Add("types", params.Types)
	}
	if params.ContractIDs != "" {
		q.Add("contractIds", params.ContractIDs)
	}
	req.URL.RawQuery = q.Encode()

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

func (d *dns) GetZone(ctx context.Context, params GetZoneRequest) (*GetZoneResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetZone")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetZone, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s", params.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZone request: %w", err)
	}

	var result GetZoneResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetZone request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetChangeList(ctx context.Context, params GetChangeListRequest) (*GetChangeListResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetChangeList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangeList, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/changelists/%s", params.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetChangeList request: %w", err)
	}

	var result GetChangeListResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetChangeList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetMasterZoneFile(ctx context.Context, params GetMasterZoneFileRequest) (string, error) {
	logger := d.Log(ctx)
	logger.Debug("GetMasterZoneFile")

	if err := params.Validate(); err != nil {
		return "", fmt.Errorf("%s: %w: %s", ErrGetMasterZoneFile, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/zone-file", params.Zone)
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

func (d *dns) PostMasterZoneFile(ctx context.Context, params PostMasterZoneFileRequest) error {
	logger := d.Log(ctx)
	logger.Debug("PostMasterZoneFile")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrPostMasterZoneFile, ErrStructValidation, err)
	}

	mtResp := ""
	pmzfURL := fmt.Sprintf("/config-dns/v2/zones/%s/zone-file", params.Zone)
	buf := bytes.NewReader([]byte(params.FileData))
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

func (d *dns) CreateZone(ctx context.Context, params CreateZoneRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone,
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := d.Log(ctx)
	logger.Debug("Zone Create")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrCreateZone, ErrStructValidation, err)
	}

	if err := ValidateZone(params.CreateZone); err != nil {
		return err
	}

	zoneMap := filterZoneCreate(params.CreateZone)

	var zoneResponse ZoneResponse
	zoneURL := "/config-dns/v2/zones/?contractId=" + params.ZoneQueryString.Contract
	if len(params.ZoneQueryString.Group) > 0 {
		zoneURL += "&gid=" + params.ZoneQueryString.Group
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

	if strings.ToUpper(params.CreateZone.Type) == "PRIMARY" {
		// Timing issue with Create immediately followed by SaveChangelist
		for _, clear := range params.ClearConn {
			// should only be one entry
			if clear {
				logger.Info("Clearing Idle Connections")
				d.Client().CloseIdleConnections()
			}
		}
	}

	return nil
}

func (d *dns) SaveChangeList(ctx context.Context, params SaveChangeListRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := d.Log(ctx)
	logger.Debug("SaveChangeList")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrSaveChangeList, ErrStructValidation, err)
	}

	reqBody, err := convertStructToReqBody("")
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/changelists/?zone=%s", params.Zone)
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

func (d *dns) SubmitChangeList(ctx context.Context, params SubmitChangeListRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := d.Log(ctx)
	logger.Debug("SubmitChangeList")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrSubmitChangeList, ErrStructValidation, err)
	}

	reqBody, err := convertStructToReqBody("")
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/changelists/%s/submit", params.Zone)
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

func (d *dns) UpdateZone(ctx context.Context, params UpdateZoneRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := d.Log(ctx)
	logger.Debug("Zone Update")

	if err := ValidateZone(params.CreateZone); err != nil {
		return err
	}

	zoneMap := filterZoneCreate(params.CreateZone)
	reqBody, err := convertStructToReqBody(zoneMap)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s", params.CreateZone.Zone)
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

func (d *dns) GetZoneNames(ctx context.Context, params GetZoneNamesRequest) (*GetZoneNamesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetZoneNames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetZoneNames, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/names", params.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZoneNames request: %w", err)
	}

	var result GetZoneNamesResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetZoneNames request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetZoneNameTypes(ctx context.Context, params GetZoneNameTypesRequest) (*GetZoneNameTypesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug(" GetZoneNameTypes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetZoneNameTypes, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types", params.Zone, params.ZoneName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZoneNameTypes request: %w", err)
	}

	var result GetZoneNameTypesResponse
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
