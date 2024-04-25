package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// TSIGQueryString contains TSIG query parameters
	TSIGQueryString struct {
		ContractIDs []string `json:"contractIds,omitempty"`
		Search      string   `json:"search,omitempty"`
		SortBy      []string `json:"sortBy,omitempty"`
		GID         int64    `json:"gid,omitempty"`
	}

	// TSIGKey contains TSIG key POST response
	TSIGKey struct {
		Name      string `json:"name"`
		Algorithm string `json:"algorithm,omitempty"`
		Secret    string `json:"secret,omitempty"`
	}

	// TSIGKeyRequest contains request parameter
	TSIGKeyRequest struct {
		Zone string
	}

	// GetTSIGKeyRequest contains request parameters for GetTSIGKey
	GetTSIGKeyRequest TSIGKeyRequest

	// GetTSIGKeyResponse contains the response data from GetTSIGKey operation
	GetTSIGKeyResponse struct {
		TSIGKey
		ZoneCount int64 `json:"zonesCount,omitempty"`
	}

	// DeleteTSIGKeyRequest contains request parameters for DeleteTSIGKey
	DeleteTSIGKeyRequest TSIGKeyRequest

	// GetTSIGKeyAliasesRequest contains request parameters for GetTSIGKeyAliases
	GetTSIGKeyAliasesRequest TSIGKeyRequest

	// GetTSIGKeyAliasesResponse contains the response data from GetTSIGKeyAliases operation
	GetTSIGKeyAliasesResponse struct {
		Zones   []string `json:"zones"`
		Aliases []string `json:"aliases"`
	}

	// TSIGKeyResponse contains TSIG key GET response
	TSIGKeyResponse struct {
		TSIGKey
		ZoneCount int64 `json:"zonesCount,omitempty"`
	}

	// TSIGKeyBulkPost contains TSIG key and a list of names of zones that should use the key. Used with update function.
	TSIGKeyBulkPost struct {
		Key   *TSIGKey `json:"key"`
		Zones []string `json:"zones"`
	}

	// TSIGZoneAliases contains list of zone aliases
	TSIGZoneAliases struct {
		Aliases []string `json:"aliases"`
	}

	// TSIGReportMeta contains metadata for TSIGReport response
	TSIGReportMeta struct {
		TotalElements int64    `json:"totalElements"`
		Search        string   `json:"search,omitempty"`
		Contracts     []string `json:"contracts,omitempty"`
		GID           int64    `json:"gid,omitempty"`
		SortBy        []string `json:"sortBy,omitempty"`
	}

	// TSIGReportResponse contains response with a list of the TSIG keys used by zones.
	TSIGReportResponse struct {
		Metadata *TSIGReportMeta   `json:"metadata"`
		Keys     []TSIGKeyResponse `json:"keys,omitempty"`
	}

	// UpdateTSIGKeyRequest contains request parameters for UpdateTSIGKey
	UpdateTSIGKeyRequest struct {
		TsigKey *TSIGKey
		Zone    string
	}

	// UpdateTSIGKeyBulkRequest contains request parameters for UpdateTSIGKeyBulk
	UpdateTSIGKeyBulkRequest struct {
		TSIGKeyBulk *TSIGKeyBulkPost
	}

	// GetTSIGKeyZonesRequest contains request parameters for GetTSIGKeyZones
	GetTSIGKeyZonesRequest struct {
		TsigKey *TSIGKey
	}

	// GetTSIGKeyZonesResponse contains the response data from GetTSIGKeyZones operation
	GetTSIGKeyZonesResponse struct {
		Zones   []string `json:"zones"`
		Aliases []string `json:"aliases"`
	}

	// ListTSIGKeysRequest contains request parameters for ListTSIGKeys
	ListTSIGKeysRequest struct {
		TsigQuery *TSIGQueryString
	}

	// ListTSIGKeysResponse contains the response data from ListTSIGKeys operation
	ListTSIGKeysResponse struct {
		Metadata *TSIGReportMeta   `json:"metadata"`
		Keys     []TSIGKeyResponse `json:"keys,omitempty"`
	}
)

var (
	// ErrGetTSIGKey is returned when GetTSIGKey fails
	ErrGetTSIGKey = errors.New("get tsig key")
	// ErrDeleteTSIGKey is returned when DeleteTSIGKey fails
	ErrDeleteTSIGKey = errors.New("delete tsig key")
	// ErrGetTSIGKeyAliases is returned when GetTSIGKeyAliases fails
	ErrGetTSIGKeyAliases = errors.New("get tsig key aliases")
	// ErrUpdateTSIGKey is returned when UpdateTSIGKey fails
	ErrUpdateTSIGKey = errors.New("updated tsig key")
	// ErrUpdateTSIGKeyBulk is returned when UpdateTSIGKeyBulk fails
	ErrUpdateTSIGKeyBulk = errors.New("update tsig key for multiple zones")
	// ErrGetTSIGKeyZones is returned when GetTSIGKeyZones fails
	ErrGetTSIGKeyZones = errors.New("list zones using tsig key")
	// ErrListTSIGKeys is returned when ListTSIGKeys fails
	ErrListTSIGKeys = errors.New("get a list of the tsig keys")
)

// Validate validates GetTSIGKeyRequest
func (r GetTSIGKeyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates DeleteTSIGKeyRequest
func (r DeleteTSIGKeyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates GetTSIGKeyAliasesRequest
func (r GetTSIGKeyAliasesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates UpdateTSIGKeyRequest
func (r UpdateTSIGKeyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone":    validation.Validate(r.Zone, validation.Required),
		"TsigKey": validation.Validate(r.TsigKey),
	})
}

// Validate validates UpdateTSIGKeyBulkRequest
func (r UpdateTSIGKeyBulkRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"TSIGKeyBulk": validation.Validate(r.TSIGKeyBulk, validation.Required),
	})
}

// Validate validates GetTSIGKeyZonesRequest
func (r GetTSIGKeyZonesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"TsigKey": validation.Validate(r.TsigKey, validation.Required),
	})
}

// Validate validates TSIGKey
func (key *TSIGKey) Validate() error {
	return validation.Errors{
		"Name":      validation.Validate(key.Name, validation.Required),
		"Algorithm": validation.Validate(key.Algorithm, validation.Required),
		"Secret":    validation.Validate(key.Secret, validation.Required),
	}.Filter()
}

// Validate validates TSIGKeyBulkPost
func (bulk *TSIGKeyBulkPost) Validate() error {
	return validation.Errors{
		"Key":   validation.Validate(bulk.Key, validation.Required),
		"Zones": validation.Validate(bulk.Zones, validation.Required),
	}.Filter()
}

func constructTSIGQueryString(tsigQueryString *TSIGQueryString) string {
	queryString := ""
	qsElems := reflect.ValueOf(tsigQueryString).Elem()
	for i := 0; i < qsElems.NumField(); i++ {
		varName := qsElems.Type().Field(i).Name
		varValue := qsElems.Field(i).Interface()
		keyVal := fmt.Sprint(varValue)
		switch varName {
		case "ContractIDs":
			contractList := ""
			for j, id := range varValue.([]string) {
				contractList += id
				if j < len(varValue.([]string))-1 {
					contractList += "%2C"
				}
			}
			if len(varValue.([]string)) > 0 {
				queryString += "contractIds=" + contractList
			}
		case "SortBy":
			sortByList := ""
			for j, sb := range varValue.([]string) {
				sortByList += sb
				if j < len(varValue.([]string))-1 {
					sortByList += "%2C"
				}
			}
			if len(varValue.([]string)) > 0 {
				queryString += "sortBy=" + sortByList
			}
		case "Search":
			if keyVal != "" {
				queryString += "search=" + keyVal
			}
		case "GID":
			if varValue.(int64) != 0 {
				queryString += "gid=" + keyVal
			}
		}
		if i < qsElems.NumField()-1 {
			queryString += "&"
		}
	}
	queryString = strings.TrimRight(queryString, "&")
	if len(queryString) > 0 {
		return "?" + queryString
	}
	return ""
}

func (d *dns) ListTSIGKeys(ctx context.Context, params ListTSIGKeysRequest) (*ListTSIGKeysResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("ListTSIGKeys")

	getURL := fmt.Sprintf("/config-dns/v2/keys%s", constructTSIGQueryString(params.TsigQuery))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListTsigKeyss request: %w", err)
	}

	var result ListTSIGKeysResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf(" ListTsigKeys request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetTSIGKeyZones(ctx context.Context, params GetTSIGKeyZonesRequest) (*GetTSIGKeyZonesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetTSIGKeyZones")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetTSIGKeyZones, ErrStructValidation, err)
	}

	reqBody, err := convertStructToReqBody(params.TsigKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := "/config-dns/v2/keys/used-by"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKeyZones request: %w", err)
	}

	var result GetTSIGKeyZonesResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKeyZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetTSIGKeyAliases(ctx context.Context, params GetTSIGKeyAliasesRequest) (*GetTSIGKeyAliasesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetTSIGKeyAliases")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetTSIGKeyAliases, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/key/used-by", params.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKeyAliases request: %w", err)
	}

	var result GetTSIGKeyAliasesResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKeyAliases request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) UpdateTSIGKeyBulk(ctx context.Context, params UpdateTSIGKeyBulkRequest) error {
	logger := d.Log(ctx)
	logger.Debug("TSIGKeyBulkUpdate")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrUpdateTSIGKeyBulk, ErrStructValidation, err)
	}

	reqBody, err := convertStructToReqBody(params.TSIGKeyBulk)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := "/config-dns/v2/keys/bulk-update"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create TsigKeyBulkUpdate request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("TsigKeyBulkUpdate request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) GetTSIGKey(ctx context.Context, params GetTSIGKeyRequest) (*GetTSIGKeyResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetTSIGKey")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetTSIGKey, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", params.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKey request: %w", err)
	}

	var result GetTSIGKeyResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) DeleteTSIGKey(ctx context.Context, params DeleteTSIGKeyRequest) error {
	logger := d.Log(ctx)
	logger.Debug("DeleteTSIGKey")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteTSIGKey, ErrStructValidation, err)
	}

	delURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", params.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DeleteTsigKey request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("DeleteTsigKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) UpdateTSIGKey(ctx context.Context, params UpdateTSIGKeyRequest) error {
	logger := d.Log(ctx)
	logger.Debug("UpdateTSIGKey")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrUpdateTSIGKey, ErrStructValidation, err)
	}

	reqBody, err := convertStructToReqBody(params.TsigKey)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", params.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create UpdateTsigKey request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("UpdateTsigKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}
