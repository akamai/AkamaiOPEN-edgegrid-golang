package dns

import (
	"context"
	"fmt"
	"net/http"
)

// BulkZonesCreate contains a list of one or more new Zones to create
type BulkZonesCreate struct {
	Zones []*ZoneCreate `json:"zones"`
}

// BulkZonesResponse contains response from bulk-create request
type BulkZonesResponse struct {
	RequestID      string `json:"requestId"`
	ExpirationDate string `json:"expirationDate"`
}

// BulkStatusResponse contains current status of a running or completed bulk-create request
type BulkStatusResponse struct {
	RequestID      string `json:"requestId"`
	ZonesSubmitted int    `json:"zonesSubmitted"`
	SuccessCount   int    `json:"successCount"`
	FailureCount   int    `json:"failureCount"`
	IsComplete     bool   `json:"isComplete"`
	ExpirationDate string `json:"expirationDate"`
}

// BulkFailedZone contains information about failed zone
type BulkFailedZone struct {
	Zone          string `json:"zone"`
	FailureReason string `json:"failureReason"`
}

// BulkCreateResultResponse contains the response from a completed bulk-create request
type BulkCreateResultResponse struct {
	RequestID                string            `json:"requestId"`
	SuccessfullyCreatedZones []string          `json:"successfullyCreatedZones"`
	FailedZones              []*BulkFailedZone `json:"failedZones"`
}

// BulkDeleteResultResponse contains the response from a completed bulk-delete request
type BulkDeleteResultResponse struct {
	RequestID                string            `json:"requestId"`
	SuccessfullyDeletedZones []string          `json:"successfullyDeletedZones"`
	FailedZones              []*BulkFailedZone `json:"failedZones"`
}

func (d *dns) GetBulkZoneCreateStatus(ctx context.Context, requestID string) (*BulkStatusResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetBulkZoneCreateStatus")

	bulkZonesURL := fmt.Sprintf("/config-dns/v2/zones/create-requests/%s", requestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneCreateStatus request: %w", err)
	}

	var result BulkStatusResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneCreateStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetBulkZoneDeleteStatus(ctx context.Context, requestID string) (*BulkStatusResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetBulkZoneDeleteStatus")

	bulkZonesURL := fmt.Sprintf("/config-dns/v2/zones/delete-requests/%s", requestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneDeleteStatus request: %w", err)
	}

	var result BulkStatusResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneDeleteStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetBulkZoneCreateResult(ctx context.Context, requestID string) (*BulkCreateResultResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetBulkZoneCreateResult")

	bulkZonesURL := fmt.Sprintf("/config-dns/v2/zones/create-requests/%s/result", requestID)
	var status BulkCreateResultResponse

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

func (d *dns) GetBulkZoneDeleteResult(ctx context.Context, requestID string) (*BulkDeleteResultResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetBulkZoneDeleteResult")

	bulkZonesURL := fmt.Sprintf("/config-dns/v2/zones/delete-requests/%s/result", requestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneDeleteResult request: %w", err)
	}

	var result BulkDeleteResultResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneDeleteResult request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) CreateBulkZones(ctx context.Context, bulkZones *BulkZonesCreate, zoneQueryString ZoneQueryString) (*BulkZonesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("CreateBulkZones")

	bulkZonesURL := "/config-dns/v2/zones/create-requests?contractId=" + zoneQueryString.Contract
	if len(zoneQueryString.Group) > 0 {
		bulkZonesURL += "&gid=" + zoneQueryString.Group
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateBulkZones request: %w", err)
	}

	var result BulkZonesResponse
	resp, err := d.Exec(req, &result, bulkZones)
	if err != nil {
		return nil, fmt.Errorf("CreateBulkZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) DeleteBulkZones(ctx context.Context, zonesList *ZoneNameListResponse, bypassSafetyChecks ...bool) (*BulkZonesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("DeleteBulkZones")

	bulkZonesURL := "/config-dns/v2/zones/delete-requests"
	if len(bypassSafetyChecks) > 0 {
		bulkZonesURL += fmt.Sprintf("?bypassSafetyChecks=%t", bypassSafetyChecks[0])
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bulkZonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeleteBulkZones request: %w", err)
	}

	var result BulkZonesResponse
	resp, err := d.Exec(req, &result, zonesList)
	if err != nil {
		return nil, fmt.Errorf("DeleteBulkZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, d.Error(resp)
	}

	return &result, nil
}
