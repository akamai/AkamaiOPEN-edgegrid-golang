package dns

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"strconv"
	"sync"
)

var (
	zoneRecordSetsWriteLock sync.Mutex
)

// Recordsets contains operations available on a record sets.
type Recordsets interface {
	// NewRecordSetResponse returns new response object.
	NewRecordSetResponse(context.Context, string) *RecordSetResponse
	// GetRecordSets retrieves record sets with Query Args. No formatting of arg values.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-recordsets
	GetRecordSets(context.Context, string, ...RecordSetQueryArgs) (*RecordSetResponse, error)
	// CreateRecordSets creates multiple record sets.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-zone-recordsets
	CreateRecordSets(context.Context, *RecordSets, string, ...bool) error
	// UpdateRecordSets replaces list of record sets.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/put-zones-zone-recordsets
	UpdateRecordSets(context.Context, *RecordSets, string, ...bool) error
}

// RecordSetQueryArgs contains query parameters for recordset request
type RecordSetQueryArgs struct {
	Page     int
	PageSize int
	Search   string
	ShowAll  bool
	SortBy   string
	Types    string
}

// RecordSets Struct. Used for Create and Update record sets. Contains a list of RecordSet objects
type RecordSets struct {
	RecordSets []RecordSet `json:"recordsets"`
}

// RecordSet contains record set metadata
type RecordSet struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	TTL   int      `json:"ttl"`
	Rdata []string `json:"rdata"`
}

// Metadata contains metadata of RecordSet response
type Metadata struct {
	LastPage      int  `json:"lastPage"`
	Page          int  `json:"page"`
	PageSize      int  `json:"pageSize"`
	ShowAll       bool `json:"showAll"`
	TotalElements int  `json:"totalElements"`
}

// RecordSetResponse contains a response with a list of record sets
type RecordSetResponse struct {
	Metadata   Metadata    `json:"metadata"`
	RecordSets []RecordSet `json:"recordsets"`
}

// Validate validates RecordSets
func (rs *RecordSets) Validate() error {
	if len(rs.RecordSets) < 1 {
		return fmt.Errorf("request initiated with empty recordsets list")
	}
	for _, rec := range rs.RecordSets {
		err := validation.Errors{
			"Name":  validation.Validate(rec.Name, validation.Required),
			"Type":  validation.Validate(rec.Type, validation.Required),
			"TTL":   validation.Validate(rec.TTL, validation.Required),
			"Rdata": validation.Validate(rec.Rdata, validation.Required),
		}.Filter()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dns) NewRecordSetResponse(_ context.Context, _ string) *RecordSetResponse {
	recordset := &RecordSetResponse{}
	return recordset
}

func (d *dns) GetRecordSets(ctx context.Context, zone string, queryArgs ...RecordSetQueryArgs) (*RecordSetResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetRecordSets")

	if len(queryArgs) > 1 {
		return nil, fmt.Errorf("invalid arguments GetRecordSets QueryArgs")
	}

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", zone)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRecordsets request: %w", err)
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
		req.URL.RawQuery = q.Encode()
	}

	var result RecordSetResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRecordsets request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) CreateRecordSets(ctx context.Context, recordSets *RecordSets, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordSetsWriteLock.Lock()
		defer zoneRecordSetsWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("CreateRecordSets")

	if err := recordSets.Validate(); err != nil {
		return err
	}

	reqBody, err := convertStructToReqBody(recordSets)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create CreateRecordsets request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("CreateRecordsets request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) UpdateRecordSets(ctx context.Context, recordSets *RecordSets, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordSetsWriteLock.Lock()
		defer zoneRecordSetsWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("UpdateRecordsets")

	if err := recordSets.Validate(); err != nil {
		return err
	}

	reqBody, err := convertStructToReqBody(recordSets)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create UpdateRecordsets request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("UpdateRecordsets request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}
