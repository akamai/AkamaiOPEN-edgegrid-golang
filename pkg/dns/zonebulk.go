package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// BulkZonesCreate contains a list of one or more new Zones to create
	BulkZonesCreate struct {
		Zones []ZoneCreate `json:"zones"`
	}

	// BulkZonesResponse contains response from bulk-create request
	BulkZonesResponse struct {
		RequestID      string `json:"requestId"`
		ExpirationDate string `json:"expirationDate"`
	}

	// BulkRequest contains request parameter
	BulkRequest struct {
		RequestID string
	}

	// BulkStatusResponse contains current status of a running or completed bulk-create request
	BulkStatusResponse struct {
		RequestID      string `json:"requestId"`
		ZonesSubmitted int    `json:"zonesSubmitted"`
		SuccessCount   int    `json:"successCount"`
		FailureCount   int    `json:"failureCount"`
		IsComplete     bool   `json:"isComplete"`
		ExpirationDate string `json:"expirationDate"`
	}

	// BulkFailedZone contains information about failed zone
	BulkFailedZone struct {
		Zone          string `json:"zone"`
		FailureReason string `json:"failureReason"`
	}

	// BulkCreateResultResponse contains the response from a completed bulk-create request
	BulkCreateResultResponse struct {
		RequestID                string           `json:"requestId"`
		SuccessfullyCreatedZones []string         `json:"successfullyCreatedZones"`
		FailedZones              []BulkFailedZone `json:"failedZones"`
	}

	// BulkDeleteResultResponse contains the response from a completed bulk-delete request
	BulkDeleteResultResponse struct {
		RequestID                string           `json:"requestId"`
		SuccessfullyDeletedZones []string         `json:"successfullyDeletedZones"`
		FailedZones              []BulkFailedZone `json:"failedZones"`
	}

	// GetBulkZoneCreateStatusRequest contains request parameters for GetBulkZoneCreateStatus
	GetBulkZoneCreateStatusRequest BulkRequest

	// GetBulkZoneCreateStatusResponse contains the response data from GetBulkZoneCreateStatus operation
	GetBulkZoneCreateStatusResponse struct {
		RequestID      string `json:"requestId"`
		ZonesSubmitted int    `json:"zonesSubmitted"`
		SuccessCount   int    `json:"successCount"`
		FailureCount   int    `json:"failureCount"`
		IsComplete     bool   `json:"isComplete"`
		ExpirationDate string `json:"expirationDate"`
	}

	// GetBulkZoneDeleteStatusRequest contains request parameters for GetBulkZoneDeleteStatus
	GetBulkZoneDeleteStatusRequest BulkRequest

	// GetBulkZoneDeleteStatusResponse contains the response data from GetBulkZoneDeleteStatus operation
	GetBulkZoneDeleteStatusResponse struct {
		RequestID      string `json:"requestId"`
		ZonesSubmitted int    `json:"zonesSubmitted"`
		SuccessCount   int    `json:"successCount"`
		FailureCount   int    `json:"failureCount"`
		IsComplete     bool   `json:"isComplete"`
		ExpirationDate string `json:"expirationDate"`
	}

	// GetBulkZoneCreateResultRequest contains request parameters for GetBulkZoneCreateResult
	GetBulkZoneCreateResultRequest BulkRequest

	// GetBulkZoneCreateResultResponse contains the response data from GetBulkZoneCreateResult operation
	GetBulkZoneCreateResultResponse struct {
		RequestID                string           `json:"requestId"`
		SuccessfullyCreatedZones []string         `json:"successfullyCreatedZones"`
		FailedZones              []BulkFailedZone `json:"failedZones"`
	}

	// GetBulkZoneDeleteResultRequest contains request parameters for GetBulkZoneDeleteResult
	GetBulkZoneDeleteResultRequest BulkRequest

	// GetBulkZoneDeleteResultResponse contains the response data from GetBulkZoneDeleteResult operation
	GetBulkZoneDeleteResultResponse struct {
		RequestID                string           `json:"requestId"`
		SuccessfullyDeletedZones []string         `json:"successfullyDeletedZones"`
		FailedZones              []BulkFailedZone `json:"failedZones"`
	}

	// CreateBulkZonesRequest contains request parameters for CreateBulkZones
	CreateBulkZonesRequest struct {
		BulkZones       *BulkZonesCreate
		ZoneQueryString ZoneQueryString
	}

	// CreateBulkZonesResponse contains the response data from CreateBulkZones operation
	CreateBulkZonesResponse struct {
		RequestID      string `json:"requestId"`
		ExpirationDate string `json:"expirationDate"`
	}

	// DeleteBulkZonesRequest contains request parameters for DeleteBulkZones
	DeleteBulkZonesRequest struct {
		ZonesList          *ZoneNameListResponse
		BypassSafetyChecks *bool
	}

	// DeleteBulkZonesResponse contains the response data from DeleteBulkZones operation
	DeleteBulkZonesResponse struct {
		RequestID      string `json:"requestId"`
		ExpirationDate string `json:"expirationDate"`
	}
)

var (
	// ErrGetBulkZoneCreateStatus is returned when GetBulkZoneCreateStatus fails
	ErrGetBulkZoneCreateStatus = errors.New("get bulk zone create status")
	// ErrGetBulkZoneDeleteStatus is returned when GetBulkZoneDeleteStatus fails
	ErrGetBulkZoneDeleteStatus = errors.New("get bulk zone delete status")
	// ErrGetBulkZoneCreateResult is returned when GetBulkZoneCreateResult fails
	ErrGetBulkZoneCreateResult = errors.New("get bulk zone create result")
	// ErrGetBulkZoneDeleteResult is returned when GetBulkZoneDeleteResult fails
	ErrGetBulkZoneDeleteResult = errors.New("get bulk zone delete result")
	// ErrCreateBulkZones is returned when CreateBulkZones fails
	ErrCreateBulkZones = errors.New("create bulk zones")
	// ErrDeleteBulkZones is returned when DeleteBulkZones fails
	ErrDeleteBulkZones = errors.New("delete bulk zones")
)

// Validate validates GetBulkZoneCreateStatusRequest
func (r GetBulkZoneCreateStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"RequestID": validation.Validate(r.RequestID, validation.Required),
	})
}

// Validate validates GetBulkZoneDeleteStatusRequest
func (r GetBulkZoneDeleteStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"RequestID": validation.Validate(r.RequestID, validation.Required),
	})
}

// Validate validates GetBulkZoneCreateResultRequest
func (r GetBulkZoneCreateResultRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"RequestID": validation.Validate(r.RequestID, validation.Required),
	})
}

// Validate validates GetBulkZoneDeleteResultRequest
func (r GetBulkZoneDeleteResultRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"RequestID": validation.Validate(r.RequestID, validation.Required),
	})
}

// Validate validates CreateBulkZonesRequest
func (r CreateBulkZonesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"BulkZones": validation.Validate(r.BulkZones, validation.Required),
	})
}

// Validate validates DeleteBulkZonesRequest
func (r DeleteBulkZonesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ZonesList": validation.Validate(r.ZonesList, validation.Required),
	})
}

func (d *dns) GetBulkZoneCreateStatus(ctx context.Context, params GetBulkZoneCreateStatusRequest) (*GetBulkZoneCreateStatusResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetBulkZoneCreateStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetBulkZoneCreateStatus, ErrStructValidation, err)
	}

	bulkZonesURL := fmt.Sprintf("/config-dns/v2/zones/create-requests/%s", params.RequestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneCreateStatus request: %w", err)
	}

	var result GetBulkZoneCreateStatusResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneCreateStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetBulkZoneDeleteStatus(ctx context.Context, params GetBulkZoneDeleteStatusRequest) (*GetBulkZoneDeleteStatusResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetBulkZoneDeleteStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetBulkZoneDeleteStatus, ErrStructValidation, err)
	}

	bulkZonesURL := fmt.Sprintf("/config-dns/v2/zones/delete-requests/%s", params.RequestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneDeleteStatus request: %w", err)
	}

	var result GetBulkZoneDeleteStatusResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneDeleteStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetBulkZoneCreateResult(ctx context.Context, params GetBulkZoneCreateResultRequest) (*GetBulkZoneCreateResultResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetBulkZoneCreateResult")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetBulkZoneCreateResult, ErrStructValidation, err)
	}

	bulkZonesURL := fmt.Sprintf("/config-dns/v2/zones/create-requests/%s/result", params.RequestID)
	var status GetBulkZoneCreateResultResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneCreateResult request: %w", err)
	}

	resp, err := d.Exec(req, &status)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneCreateResult request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &status, nil
}

func (d *dns) GetBulkZoneDeleteResult(ctx context.Context, params GetBulkZoneDeleteResultRequest) (*GetBulkZoneDeleteResultResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetBulkZoneDeleteResult")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetBulkZoneDeleteResult, ErrStructValidation, err)
	}

	bulkZonesURL := fmt.Sprintf("/config-dns/v2/zones/delete-requests/%s/result", params.RequestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneDeleteResult request: %w", err)
	}

	var result GetBulkZoneDeleteResultResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneDeleteResult request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) CreateBulkZones(ctx context.Context, params CreateBulkZonesRequest) (*CreateBulkZonesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("CreateBulkZones")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateBulkZones, ErrStructValidation, err)
	}

	bulkZonesURL := "/config-dns/v2/zones/create-requests?contractId=" + params.ZoneQueryString.Contract
	if len(params.ZoneQueryString.Group) > 0 {
		bulkZonesURL += "&gid=" + params.ZoneQueryString.Group
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateBulkZones request: %w", err)
	}

	var result CreateBulkZonesResponse
	resp, err := d.Exec(req, &result, params.BulkZones)
	if err != nil {
		return nil, fmt.Errorf("CreateBulkZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) DeleteBulkZones(ctx context.Context, params DeleteBulkZonesRequest) (*DeleteBulkZonesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("DeleteBulkZones")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteBulkZones, ErrStructValidation, err)
	}

	bulkZonesURL := "/config-dns/v2/zones/delete-requests"
	if params.BypassSafetyChecks != nil {
		bulkZonesURL += fmt.Sprintf("?bypassSafetyChecks=%t", *params.BypassSafetyChecks)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeleteBulkZones request: %w", err)
	}

	var result DeleteBulkZonesResponse
	resp, err := d.Exec(req, &result, params.ZonesList)
	if err != nil {
		return nil, fmt.Errorf("DeleteBulkZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, d.Error(resp)
	}

	return &result, nil
}
