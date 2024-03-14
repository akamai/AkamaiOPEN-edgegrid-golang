package dns

import (
	"context"
	"fmt"
	"net/http"

	"sync"
)

// Records contains operations available on a Record resource.
type Records interface {
	// GetRecordList retrieves recordset list based on type.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-recordsets
	GetRecordList(context.Context, string, string, string) (*RecordSetResponse, error)
	// GetRdata retrieves record rdata, e.g. target.
	GetRdata(context.Context, string, string, string) ([]string, error)
	// ProcessRdata process rdata.
	ProcessRdata(context.Context, []string, string) []string
	// ParseRData parses rdata. returning map.
	ParseRData(context.Context, string, []string) map[string]interface{}
	// GetRecord retrieves a recordset and returns as RecordBody.
	//
	// See:  https://techdocs.akamai.com/edge-dns/reference/get-zone-name-type
	GetRecord(context.Context, string, string, string) (*RecordBody, error)
	// CreateRecord creates recordset.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-zone-names-name-types-type
	CreateRecord(context.Context, *RecordBody, string, ...bool) error
	// DeleteRecord removes recordset.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/delete-zone-name-type
	DeleteRecord(context.Context, *RecordBody, string, ...bool) error
	// UpdateRecord replaces the recordset.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/put-zones-zone-names-name-types-type
	UpdateRecord(context.Context, *RecordBody, string, ...bool) error
}

// RecordBody contains request body for dns record
type RecordBody struct {
	Name       string   `json:"name,omitempty"`
	RecordType string   `json:"type,omitempty"`
	TTL        int      `json:"ttl,omitempty"`
	Active     bool     `json:"active,omitempty"`
	Target     []string `json:"rdata,omitempty"`
}

var (
	zoneRecordWriteLock sync.Mutex
)

// Validate validates RecordBody
func (rec *RecordBody) Validate() error {
	if len(rec.Name) < 1 {
		return fmt.Errorf("RecordBody is missing Name")
	}
	if len(rec.RecordType) < 1 {
		return fmt.Errorf("RecordBody is missing RecordType")
	}
	if rec.TTL == 0 {
		return fmt.Errorf("RecordBody is missing TTL")
	}
	if rec.Target == nil || len(rec.Target) < 1 {
		return fmt.Errorf("RecordBody is missing Target")
	}

	return nil
}

// Eval option lock arg passed into writable endpoints. Default is true, e.g. lock
func localLock(lockArg []bool) bool {
	for _, lock := range lockArg {
		// should only be one entry
		return lock
	}

	return true
}

func (d *dns) CreateRecord(ctx context.Context, record *RecordBody, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone,
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("CreateRecord")
	logger.Debugf("DNS Lib Create Record: [%v]", record)
	if err := record.Validate(); err != nil {
		logger.Errorf("Record content not valid: %w", err)
		return fmt.Errorf("CreateRecord content not valid. [%w]", err)
	}

	reqBody, err := convertStructToReqBody(record)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", zone, record.Name, record.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create CreateRecord request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("CreateRecord request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) UpdateRecord(ctx context.Context, record *RecordBody, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("UpdateRecord")
	logger.Debugf("DNS Lib Update Record: [%v]", record)
	if err := record.Validate(); err != nil {
		logger.Errorf("Record content not valid: %s", err.Error())
		return fmt.Errorf("UpdateRecord content not valid. [%w]", err)
	}

	reqBody, err := convertStructToReqBody(record)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", zone, record.Name, record.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create UpdateRecord request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("UpdateRecord request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) DeleteRecord(ctx context.Context, record *RecordBody, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("DeleteRecord")

	if err := record.Validate(); err != nil {
		logger.Errorf("Record content not valid: %w", err)
		return fmt.Errorf("DeleteRecord content not valid. [%w]", err)
	}

	deleteURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", zone, record.Name, record.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DeleteRecord request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("DeleteRecord request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}
