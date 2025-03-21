package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// SortOrder represents SortOrder enum.
	SortOrder string

	// CertType represents CertType enum.
	CertType string

	// ListActivePropertyHostnamesRequest contains parameters required to list active property hostnames.
	ListActivePropertyHostnamesRequest struct {
		PropertyID        string
		Offset            int
		Limit             int
		Sort              SortOrder
		Hostname          string
		CnameTo           string
		Network           ActivationNetwork
		ContractID        string
		GroupID           string
		IncludeCertStatus bool
	}

	// GetActivePropertyHostnamesDiffRequest contains parameters required to list active property hostnames diff.
	GetActivePropertyHostnamesDiffRequest struct {
		PropertyID string
		Offset     int
		Limit      int
		ContractID string
		GroupID    string
	}

	// ListActivePropertyHostnamesResponse contains information about each of the active property hostnames.
	ListActivePropertyHostnamesResponse struct {
		AccountID     string                 `json:"accountId"`
		AvailableSort []SortOrder            `json:"availableSort"`
		ContractID    string                 `json:"contractId"`
		CurrentSort   SortOrder              `json:"currentSort"`
		DefaultSort   SortOrder              `json:"defaultSort"`
		GroupID       string                 `json:"groupId"`
		PropertyID    string                 `json:"propertyId"`
		PropertyName  string                 `json:"propertyName"`
		Hostnames     HostnamesResponseItems `json:"hostnames"`
	}

	// GetActivePropertyHostnamesDiffResponse contains information about each of the active property hostnames diff.
	GetActivePropertyHostnamesDiffResponse struct {
		AccountID  string                     `json:"accountId"`
		ContractID string                     `json:"contractId"`
		GroupID    string                     `json:"groupId"`
		PropertyID string                     `json:"propertyId"`
		Hostnames  HostnamesDiffResponseItems `json:"hostnames"`
	}

	// HostnamesResponseItems contains the response body for ListActivePropertyHostnamesResponse.
	HostnamesResponseItems struct {
		Items            []HostnameItem `json:"items"`
		CurrentItemCount int            `json:"currentItemCount"`
		NextLink         *string        `json:"nextLink"`
		PreviousLink     *string        `json:"previousLink"`
		TotalItems       int            `json:"totalItems"`
	}

	// HostnamesDiffResponseItems contains the response body for GetActivePropertyHostnamesDiffResponse.
	HostnamesDiffResponseItems struct {
		Items            []HostnameDiffItem `json:"items"`
		CurrentItemCount int                `json:"currentItemCount"`
		NextLink         *string            `json:"nextLink"`
		PreviousLink     *string            `json:"previousLink"`
		TotalItems       int                `json:"totalItems"`
	}

	// HostnameItem contains information about each of the HostnamesResponseItems.
	HostnameItem struct {
		CertStatus               *CertStatusItem   `json:"certStatus"`
		CnameFrom                string            `json:"cnameFrom"`
		CnameType                HostnameCnameType `json:"cnameType"`
		ProductionCertType       CertType          `json:"productionCertType"`
		ProductionCnameTo        string            `json:"productionCnameTo"`
		ProductionEdgeHostnameID string            `json:"productionEdgeHostnameId"`
		StagingCertType          CertType          `json:"stagingCertType"`
		StagingCnameTo           string            `json:"StagingCnameTo"`
		StagingEdgeHostnameID    string            `json:"stagingEdgeHostnameId"`
	}

	// HostnameDiffItem contains information about each of the HostnamesDiffResponseItems.
	HostnameDiffItem struct {
		CnameFrom                      string            `json:"cnameFrom"`
		ProductionCertProvisioningType CertType          `json:"productionCertProvisioningType"`
		ProductionCnameTo              string            `json:"productionCnameTo"`
		ProductionCnameType            HostnameCnameType `json:"productionCnameType"`
		ProductionEdgeHostnameID       string            `json:"productionEdgeHostnameId"`
		StagingCertProvisioningType    CertType          `json:"stagingCertProvisioningType"`
		StagingCnameTo                 string            `json:"stagingCnameTo"`
		StagingCnameType               HostnameCnameType `json:"stagingCnameType"`
		StagingEdgeHostnameID          string            `json:"stagingEdgeHostnameId"`
	}
)

const (
	// SortAscending represents ascending sorting by hostname.
	SortAscending SortOrder = "hostname:a"
	// SortDescending represents descending sorting by hostname.
	SortDescending SortOrder = "hostname:d"
	// CertTypeCPSManaged indicates that the certificate is provisioned using the Certificate Provisioning System (CPS).
	CertTypeCPSManaged CertType = "CPS_MANAGED"
	// CertTypeDefault indicates that the certificate is a Default Domain Validation (DV) certificate.
	CertTypeDefault CertType = "DEFAULT"
	// maxHostnamesPerPage indicates the maximum possible value for 'limit' parameter for Get and List active property hostnames.
	maxHostnamesPerPage int = 999
)

var (
	// ErrListActivePropertyHostnames represents error when fetching active property hostnames fails.
	ErrListActivePropertyHostnames = errors.New("fetching active property hostnames")

	// ErrGetActivePropertyHostnamesDiff represents error when fetching active property hostnames diff fails.
	ErrGetActivePropertyHostnamesDiff = errors.New("fetching active property hostnames diff")
)

// Validate validates SortOrder.
func (o SortOrder) Validate() validation.InRule {
	return validation.In(SortAscending, SortDescending).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'",
			o, SortAscending, SortDescending))
}

// Validate validates CertType.
func (t CertType) Validate() validation.InRule {
	return validation.In(CertTypeCPSManaged, CertTypeDefault).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'",
			t, CertTypeCPSManaged, CertTypeDefault))
}

// Validate validates ListActivePropertyHostnamesRequest.
func (r ListActivePropertyHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"Network":    validation.Validate(r.Network, r.Network.Validate()),
		"Sort":       validation.Validate(r.Sort, r.Sort.Validate()),
		"Offset":     validation.Validate(r.Offset, validation.Min(0)),
		"Limit":      validation.Validate(r.Limit, validation.Min(1), validation.Max(maxHostnamesPerPage)),
	}.Filter()
}

// Validate validates GetActivePropertyHostnamesDiffRequest.
func (r GetActivePropertyHostnamesDiffRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"Offset":     validation.Validate(r.Offset, validation.Min(0)),
		"Limit":      validation.Validate(r.Limit, validation.Min(1), validation.Max(maxHostnamesPerPage)),
	}.Filter()
}

func (p *papi) ListActivePropertyHostnames(ctx context.Context, params ListActivePropertyHostnamesRequest) (*ListActivePropertyHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListActivePropertyHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListActivePropertyHostnames, ErrStructValidation, err)
	}

	baseURL := fmt.Sprintf("/papi/v1/properties/%s/hostnames", params.PropertyID)

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse base URL: %s", ErrListActivePropertyHostnames, err)
	}

	query := parsedURL.Query()
	if params.ContractID != "" {
		query.Set("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		query.Set("groupId", params.GroupID)
	}
	if params.Sort != "" {
		query.Set("sort", string(params.Sort))
	}
	if params.Hostname != "" {
		query.Set("hostname", params.Hostname)
	}
	if params.CnameTo != "" {
		query.Set("cnameTo", params.CnameTo)
	}
	if params.Network != "" {
		query.Set("network", string(params.Network))
	}
	if params.IncludeCertStatus {
		query.Set("includeCertStatus", fmt.Sprintf("%t", params.IncludeCertStatus))
	}
	if params.Limit != 0 {
		query.Set("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", fmt.Sprintf("%d", params.Offset))
	}

	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListActivePropertyHostnames, err)
	}

	var hostnames ListActivePropertyHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListActivePropertyHostnames, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListActivePropertyHostnames, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) GetActivePropertyHostnamesDiff(ctx context.Context, params GetActivePropertyHostnamesDiffRequest) (*GetActivePropertyHostnamesDiffResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetActivePropertyHostnamesDiff")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetActivePropertyHostnamesDiff, ErrStructValidation, err)
	}

	baseURL := fmt.Sprintf("/papi/v1/properties/%s/hostnames/diff", params.PropertyID)

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse base URL: %s", ErrGetActivePropertyHostnamesDiff, err)
	}

	query := parsedURL.Query()
	if params.ContractID != "" {
		query.Set("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		query.Set("groupId", params.GroupID)
	}
	if params.Limit != 0 {
		query.Set("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", fmt.Sprintf("%d", params.Offset))
	}

	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetActivePropertyHostnamesDiff, err)
	}

	var hostnamesDiff GetActivePropertyHostnamesDiffResponse
	resp, err := p.Exec(req, &hostnamesDiff)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetActivePropertyHostnamesDiff, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetActivePropertyHostnamesDiff, p.Error(resp))
	}

	return &hostnamesDiff, nil
}
