package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// History is a CPS interface for History management
	History interface {
		// GetDVHistory is a domain name validation history for the enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/get-dv-history
		GetDVHistory(context.Context, GetDVHistoryRequest) (*GetDVHistoryResponse, error)

		// GetCertificateHistory views the certificate history.
		//
		// See: https://techdocs.akamai.com/cps/reference/get-history-certificates
		GetCertificateHistory(context.Context, GetCertificateHistoryRequest) (*GetCertificateHistoryResponse, error)

		// GetChangeHistory views the change history for enrollment.
		//
		// See: https://techdocs.akamai.com/cps/reference/get-history-changes
		GetChangeHistory(context.Context, GetChangeHistoryRequest) (*GetChangeHistoryResponse, error)
	}

	// GetDVHistoryRequest represents request for GetDVHistory operation
	GetDVHistoryRequest struct {
		EnrollmentID int
	}

	// GetDVHistoryResponse represents response for GetDVHistory operation
	GetDVHistoryResponse struct {
		Results []HistoryResult `json:"results"`
	}

	// HistoryResult represents a piece of history for GetDVHistory operation
	HistoryResult struct {
		Domain        string          `json:"domain"`
		DomainHistory []DomainHistory `json:"domainHistory"`
	}

	// DomainHistory represents a history for single domain for GetDVHistory operation
	DomainHistory struct {
		Domain             string             `json:"domain"`
		Challenges         []Challenge        `json:"challenges"`
		Error              string             `json:"error"`
		Expires            string             `json:"expires"`
		FullPath           string             `json:"fullPath"`
		RedirectFullPath   string             `json:"redirectFullPath"`
		RequestTimestamp   string             `json:"requestTimestamp"`
		ResponseBody       string             `json:"responseBody"`
		Status             string             `json:"status"`
		Token              string             `json:"token"`
		ValidatedTimestamp string             `json:"validatedTimestamp"`
		ValidationRecords  []ValidationRecord `json:"validationRecords"`
		ValidationStatus   string             `json:"validationStatus"`
	}

	// GetCertificateHistoryRequest represents request for GetCertificateHistory operation
	GetCertificateHistoryRequest struct {
		EnrollmentID int
	}

	// GetCertificateHistoryResponse represents response for GetCertificateHistory operation
	GetCertificateHistoryResponse struct {
		Certificates []HistoryCertificate `json:"certificates"`
	}

	// HistoryCertificate represents a piece of enrollment's certificate history for GetCertificateHistory operation
	HistoryCertificate struct {
		DeploymentStatus         string              `json:"deploymentStatus"`
		Geography                string              `json:"geography"`
		MultiStackedCertificates []CertificateObject `json:"multiStackedCertificates"`
		PrimaryCertificate       CertificateObject   `json:"primaryCertificate"`
		RA                       string              `json:"ra"`
		Slots                    []int               `json:"slots"`
		StagingStatus            string              `json:"stagingStatus"`
		Type                     string              `json:"type"`
	}

	//CertificateObject represent certificate for enrollment
	CertificateObject struct {
		Certificate  string `json:"certificate"`
		Expiry       string `json:"expiry"`
		KeyAlgorithm string `json:"keyAlgorithm"`
		TrustChain   string `json:"trustChain"`
	}

	// GetChangeHistoryRequest represents request for GetChangeHistory operation
	GetChangeHistoryRequest struct {
		EnrollmentID int
	}

	// GetChangeHistoryResponse represents response for GetChangeHistory operation
	GetChangeHistoryResponse struct {
		Changes []ChangeHistory `json:"changes"`
	}

	// ChangeHistory represents a piece of enrollment's history for a single change for GetChangeHistory operation
	ChangeHistory struct {
		Action                         string                     `json:"action"`
		ActionDescription              string                     `json:"actionDescription"`
		BusinessCaseID                 string                     `json:"businessCaseId"`
		CreatedBy                      string                     `json:"createdBy"`
		CreatedOn                      string                     `json:"createdOn"`
		LastUpdated                    string                     `json:"lastUpdated"`
		MultiStackedCertificates       []CertificateChangeHistory `json:"multiStackedCertificates"`
		PrimaryCertificate             CertificateChangeHistory   `json:"primaryCertificate"`
		PrimaryCertificateOrderDetails CertificateOrderDetails    `json:"primaryCertificateOrderDetails"`
		RA                             string                     `json:"ra"`
		Status                         string                     `json:"status"`
	}

	// CertificateChangeHistory represents certificate returned in GetChangeHistory operation
	CertificateChangeHistory struct {
		Certificate  string `json:"certificate"`
		TrustChain   string `json:"trustChain"`
		CSR          string `json:"csr"`
		KeyAlgorithm string `json:"keyAlgorithm"`
	}

	// CertificateOrderDetails represents CA order details for a Change
	CertificateOrderDetails struct {
		OrderID string `json:"orderId"`
	}
)

var (
	// ErrGetDVHistory is returned when GetDVHistory fails
	ErrGetDVHistory = errors.New("get dv history")
	// ErrGetCertificateHistory is returned when GetDVHistory fails
	ErrGetCertificateHistory = errors.New("get certificate history")
	// ErrGetChangeHistory is returned when GetDVHistory fails
	ErrGetChangeHistory = errors.New("get change history")
)

// Validate validates GetDVHistoryRequest
func (r GetDVHistoryRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(r.EnrollmentID, validation.Required),
	}.Filter()
}

// Validate validates GetCertificateHistoryRequest
func (r GetCertificateHistoryRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(r.EnrollmentID, validation.Required),
	}.Filter()
}

// Validate validates GetChangeHistoryRequest
func (r GetChangeHistoryRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(r.EnrollmentID, validation.Required),
	}.Filter()
}

func (c *cps) GetDVHistory(ctx context.Context, params GetDVHistoryRequest) (*GetDVHistoryResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetDVHistory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetDVHistory, ErrStructValidation, err)
	}

	url := fmt.Sprintf("/cps/v2/enrollments/%d/dv-history", params.EnrollmentID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetDVHistory, err)
	}
	var result GetDVHistoryResponse
	req.Header.Set("Accept", "application/vnd.akamai.cps.dv-history.v1+json")
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetDVHistory, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetDVHistory, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) GetCertificateHistory(ctx context.Context, params GetCertificateHistoryRequest) (*GetCertificateHistoryResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetCertificateHistory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCertificateHistory, ErrStructValidation, err)
	}

	url := fmt.Sprintf("/cps/v2/enrollments/%d/history/certificates", params.EnrollmentID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetCertificateHistory, err)
	}
	var result GetCertificateHistoryResponse
	req.Header.Set("Accept", "application/vnd.akamai.cps.certificate-history.v2+json")
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCertificateHistory, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCertificateHistory, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) GetChangeHistory(ctx context.Context, params GetChangeHistoryRequest) (*GetChangeHistoryResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetChangeHistory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangeHistory, ErrStructValidation, err)
	}

	url := fmt.Sprintf("/cps/v2/enrollments/%d/history/changes", params.EnrollmentID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetChangeHistory, err)
	}
	var result GetChangeHistoryResponse
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-history.v5+json")
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangeHistory, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangeHistory, c.Error(resp))
	}

	return &result, nil
}
