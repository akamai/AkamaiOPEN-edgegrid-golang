package dns

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"reflect"
	"strings"
)

type (
	// TSIGKeys contains operations available on TSIKeyG resource.
	TSIGKeys interface {
		// ListTSIGKeys lists the TSIG keys used by zones that you are allowed to manage.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-keys
		ListTSIGKeys(context.Context, *TSIGQueryString) (*TSIGReportResponse, error)
		// GetTSIGKeyZones retrieves DNS Zones using TSIG key.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-keys-used-by
		GetTSIGKeyZones(context.Context, *TSIGKey) (*ZoneNameListResponse, error)
		// GetTSIGKeyAliases retrieves a DNS Zone's aliases.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-key-used-by
		GetTSIGKeyAliases(context.Context, string) (*ZoneNameListResponse, error)
		// TSIGKeyBulkUpdate updates Bulk Zones TSIG key.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-keys-bulk-update
		TSIGKeyBulkUpdate(context.Context, *TSIGKeyBulkPost) error
		// GetTSIGKey retrieves a TSIG key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-key
		GetTSIGKey(context.Context, string) (*TSIGKeyResponse, error)
		// DeleteTSIGKey deletes TSIG key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/delete-zones-zone-key
		DeleteTSIGKey(context.Context, string) error
		// UpdateTSIGKey updates TSIG key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/put-zones-zone-key
		UpdateTSIGKey(context.Context, *TSIGKey, string) error
	}

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
		Metadata *TSIGReportMeta    `json:"metadata"`
		Keys     []*TSIGKeyResponse `json:"keys,omitempty"`
	}
)

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

func (d *dns) ListTSIGKeys(ctx context.Context, tsigQueryString *TSIGQueryString) (*TSIGReportResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("ListTSIGKeys")

	getURL := fmt.Sprintf("/config-dns/v2/keys%s", constructTSIGQueryString(tsigQueryString))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListTsigKeyss request: %w", err)
	}

	var result TSIGReportResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf(" ListTsigKeys request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetTSIGKeyZones(ctx context.Context, tsigKey *TSIGKey) (*ZoneNameListResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetTSIGKeyZones")

	if err := tsigKey.Validate(); err != nil {
		return nil, err
	}

	reqBody, err := convertStructToReqBody(tsigKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := "/config-dns/v2/keys/used-by"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKeyZones request: %w", err)
	}

	var result ZoneNameListResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKeyZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetTSIGKeyAliases(ctx context.Context, zone string) (*ZoneNameListResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetTSIGKeyAliases")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/key/used-by", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKeyAliases request: %w", err)
	}

	var result ZoneNameListResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKeyAliases request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) TSIGKeyBulkUpdate(ctx context.Context, tsigBulk *TSIGKeyBulkPost) error {
	logger := d.Log(ctx)
	logger.Debug("TSIGKeyBulkUpdate")

	if err := tsigBulk.Validate(); err != nil {
		return err
	}

	reqBody, err := convertStructToReqBody(tsigBulk)
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

func (d *dns) GetTSIGKey(ctx context.Context, zone string) (*TSIGKeyResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetTSIGKey")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKey request: %w", err)
	}

	var result TSIGKeyResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) DeleteTSIGKey(ctx context.Context, zone string) error {
	logger := d.Log(ctx)
	logger.Debug("DeleteTSIGKey")

	delURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", zone)
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

func (d *dns) UpdateTSIGKey(ctx context.Context, tsigKey *TSIGKey, zone string) error {
	logger := d.Log(ctx)
	logger.Debug("UpdateTSIGKey")

	if err := tsigKey.Validate(); err != nil {
		return err
	}

	reqBody, err := convertStructToReqBody(tsigKey)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", zone)
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
