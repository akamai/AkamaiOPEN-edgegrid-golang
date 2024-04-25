package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	zoneRecordSetsWriteLock sync.Mutex
)

type (
	// RecordSetQueryArgs contains query parameters for recordset request
	RecordSetQueryArgs struct {
		Page     int
		PageSize int
		Search   string
		ShowAll  bool
		SortBy   string
		Types    string
	}

	// RecordSets Struct. Used for Create and Update record sets. Contains a list of RecordSet objects
	RecordSets struct {
		RecordSets []RecordSet `json:"recordsets"`
	}

	// RecordSet contains record set metadata
	RecordSet struct {
		Name  string   `json:"name"`
		Type  string   `json:"type"`
		TTL   int      `json:"ttl"`
		Rdata []string `json:"rdata"`
	}

	// Metadata contains metadata of RecordSet response
	Metadata struct {
		LastPage      int  `json:"lastPage"`
		Page          int  `json:"page"`
		PageSize      int  `json:"pageSize"`
		ShowAll       bool `json:"showAll"`
		TotalElements int  `json:"totalElements"`
	}

	// GetRecordSetsRequest contains request parameters for GetRecordSets
	GetRecordSetsRequest struct {
		Zone      string
		QueryArgs *RecordSetQueryArgs
	}

	// GetRecordSetsResponse contains the response data from GetRecordSets operation
	GetRecordSetsResponse struct {
		Metadata   Metadata    `json:"metadata"`
		RecordSets []RecordSet `json:"recordsets"`
	}

	// RecordSetsRequest contains request parameters
	RecordSetsRequest struct {
		RecordSets *RecordSets
		Zone       string
		RecLock    []bool
	}

	// CreateRecordSetsRequest contains request parameters for CreateRecordSets
	CreateRecordSetsRequest RecordSetsRequest

	// UpdateRecordSetsRequest contains request parameters for UpdateRecordSets
	UpdateRecordSetsRequest RecordSetsRequest
)

var (
	// ErrCreateRecordSets is returned when CreateRecordSets fails
	ErrCreateRecordSets = errors.New("create record sets")
	// ErrGetRecordSets is returned when GetRecordSets fails
	ErrGetRecordSets = errors.New("get record sets")
	// ErrUpdateRecordSets is returned when UpdateRecordSets fails
	ErrUpdateRecordSets = errors.New("update record sets")
)

// Validate validates GetRecordSetsRequest
func (r GetRecordSetsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone": validation.Validate(r.Zone, validation.Required),
	})
}

// Validate validates CreateRecordSetsRequest
func (r CreateRecordSetsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone":       validation.Validate(r.Zone, validation.Required),
		"RecordSets": validation.Validate(r.RecordSets, validation.Required),
	})
}

// Validate validates UpdateRecordSetsRequest
func (r UpdateRecordSetsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone":       validation.Validate(r.Zone, validation.Required),
		"RecordSets": validation.Validate(r.RecordSets, validation.Required),
	})
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

func (d *dns) GetRecordSets(ctx context.Context, params GetRecordSetsRequest) (*GetRecordSetsResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetRecordSets")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetRecordSets, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", params.Zone)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRecordsets request: %w", err)
	}

	if params.QueryArgs != nil {
		q := req.URL.Query()
		if params.QueryArgs.Page > 0 {
			q.Add("page", strconv.Itoa(params.QueryArgs.Page))
		}
		if params.QueryArgs.PageSize > 0 {
			q.Add("pageSize", strconv.Itoa(params.QueryArgs.PageSize))
		}
		if params.QueryArgs.Search != "" {
			q.Add("search", params.QueryArgs.Search)
		}
		q.Add("showAll", strconv.FormatBool(params.QueryArgs.ShowAll))
		if params.QueryArgs.SortBy != "" {
			q.Add("sortBy", params.QueryArgs.SortBy)
		}
		if params.QueryArgs.Types != "" {
			q.Add("types", params.QueryArgs.Types)
		}
		req.URL.RawQuery = q.Encode()
	}

	var result GetRecordSetsResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRecordsets request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) CreateRecordSets(ctx context.Context, params CreateRecordSetsRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(params.RecLock) {
		zoneRecordSetsWriteLock.Lock()
		defer zoneRecordSetsWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("CreateRecordSets")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrCreateRecordSets, ErrStructValidation, err)
	}

	if err := params.RecordSets.Validate(); err != nil {
		return err
	}

	reqBody, err := convertStructToReqBody(params.RecordSets)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", params.Zone)
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

func (d *dns) UpdateRecordSets(ctx context.Context, params UpdateRecordSetsRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(params.RecLock) {
		zoneRecordSetsWriteLock.Lock()
		defer zoneRecordSetsWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("UpdateRecordsets")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrUpdateRecordSets, ErrStructValidation, err)
	}

	if err := params.RecordSets.Validate(); err != nil {
		return err
	}

	reqBody, err := convertStructToReqBody(params.RecordSets)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", params.Zone)
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
