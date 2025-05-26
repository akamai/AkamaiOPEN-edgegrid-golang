package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// RecordBody contains request body for dns record
	RecordBody struct {
		Name       string   `json:"name,omitempty"`
		RecordType string   `json:"type,omitempty"`
		TTL        int      `json:"ttl,omitempty"`
		Active     bool     `json:"active,omitempty"`
		Target     []string `json:"rdata,omitempty"`
	}

	// RecordRequest contains request parameters
	RecordRequest struct {
		Record  *RecordBody
		Zone    string
		RecLock []bool
	}
	// CreateRecordRequest contains request parameters for CreateRecord
	CreateRecordRequest RecordRequest

	// UpdateRecordRequest contains request parameters for UpdateRecord
	UpdateRecordRequest RecordRequest

	// DeleteRecordRequest contains request parameters for DeleteRecord
	DeleteRecordRequest struct {
		Zone       string
		Name       string
		RecordType string
		RecLock    []bool
	}
)

var (
	zoneRecordWriteLock sync.Mutex
)

var (
	// ErrCreateRecord is returned when CreateRecord fails
	ErrCreateRecord = errors.New("create record")
	// ErrUpdateRecord is returned when UpdateRecord fails
	ErrUpdateRecord = errors.New("update record")
	// ErrDeleteRecord is returned when UpdateRecord fails
	ErrDeleteRecord = errors.New("delete record")
)

// Validate validates CreateRecordRequest
func (r CreateRecordRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone":   validation.Validate(r.Zone, validation.Required),
		"Record": validation.Validate(r.Record, validation.Required),
	})
}

// Validate validates UpdateRecordRequest
func (r UpdateRecordRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone":   validation.Validate(r.Zone, validation.Required),
		"Record": validation.Validate(r.Record, validation.Required),
	})
}

// Validate validates DeleteRecordRequest
func (r DeleteRecordRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Zone":       validation.Validate(r.Zone, validation.Required),
		"Name":       validation.Validate(r.Name, validation.Required),
		"RecordType": validation.Validate(r.RecordType, validation.Required),
	})
}

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

func (d *dns) CreateRecord(ctx context.Context, params CreateRecordRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone,
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(params.RecLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("CreateRecord")
	logger.Debugf("DNS Lib Create Record: [%v]", params.Record)
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrCreateRecord, ErrStructValidation, err)
	}

	reqBody, err := convertStructToReqBody(params.Record)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", params.Zone,
		params.Record.Name, params.Record.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create CreateRecord request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("CreateRecord request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) UpdateRecord(ctx context.Context, params UpdateRecordRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(params.RecLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("UpdateRecord")
	logger.Debugf("DNS Lib Update Record: [%v]", params.Record)

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrUpdateRecord, ErrStructValidation, err)
	}

	reqBody, err := convertStructToReqBody(params.Record)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", params.Zone,
		params.Record.Name, params.Record.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create UpdateRecord request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("UpdateRecord request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return d.Error(resp)
	}

	return nil
}

func (d *dns) DeleteRecord(ctx context.Context, params DeleteRecordRequest) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(params.RecLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := d.Log(ctx)
	logger.Debug("DeleteRecord")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteRecord, ErrStructValidation, err)
	}

	deleteURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", params.Zone,
		params.Name, params.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DeleteRecord request: %w", err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("DeleteRecord request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return d.Error(resp)
	}

	return nil
}
