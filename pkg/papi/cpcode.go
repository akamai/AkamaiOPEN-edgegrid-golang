package papi

import (
	"context"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi/tools"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
	"net/http"
)

type (
	// CPCodes contains operations available on CPCode resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#cpcodesgroup
	CPCodes interface {
		// GetCPCodes lists all available CP codes
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getcpcodes
		GetCPCodes(context.Context, GetCPCodesRequest) (*GetCPCodesResponse, error)

		// GetCPCode gets CP code with provided ID
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getcpcode
		GetCPCode(context.Context, GetCPCodeRequest) (*GetCPCodesResponse, error)

		// CreateCPCode creates a new CP code
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#postcpcodes
		CreateCPCode(context.Context, CreateCPCodeRequest) (*CreateCPCodeResponse, error)
	}

	// CPCode contains CP code resource data
	CPCode struct {
		ID          string   `json:"cpcodeId"`
		Name        string   `json:"cpcodeName"`
		CreatedDate string   `json:"createdDate"`
		ProductID   string   `json:"productId"`
		ProductIDs  []string `json:"productIds"`
	}

	// CPCodeItems contains a list of CPCode items
	CPCodeItems struct {
		Items []CPCode `json:"items"`
	}

	// GetCPCodesResponse is a response returned while fetching CP codes
	GetCPCodesResponse struct {
		AccountID  string      `json:"accountId"`
		ContractID string      `json:"contractId"`
		GroupID    string      `json:"groupId"`
		CPCodes    CPCodeItems `json:"cpcodes"`
	}

	// CreateCPCodeRequest contains data required to create CP code (both request body and group/contract infromation
	CreateCPCodeRequest struct {
		ContractID string
		GroupID    string
		CPCode     CreateCPCode
	}

	// CreateCPCode contains the request body for CP code creation
	CreateCPCode struct {
		ProductID  string `json:"productId"`
		CPCodeName string `json:"cpcodeName"`
	}

	// CreateCPCodeResponse contains the response from CP code creation as well as the ID of created resource
	CreateCPCodeResponse struct {
		CPCodeLink string `json:"cpcodeLink"`
		CPCodeID   string `json:"-"`
	}

	// GetCPCodeRequest gets details about a CP code.
	GetCPCodeRequest struct {
		CPCodeID   string
		ContractID string
		GroupID    string
	}

	// GetCPCodesRequest contains parameters require to list/create CP codes
	// GroupID and ContractID are required as part of every CP code operation, ID is required only for operating on specific CP code
	GetCPCodesRequest struct {
		ContractID string
		GroupID    string
	}
)

// Validate validates GetCPCodesRequest
func (cp GetCPCodesRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(cp.ContractID, validation.Required),
		"GroupID":    validation.Validate(cp.GroupID, validation.Required),
	}.Filter()
}

// Validate validates GetCPCodeRequest
func (cp GetCPCodeRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(cp.ContractID, validation.Required),
		"GroupID":    validation.Validate(cp.GroupID, validation.Required),
		"CPCodeID":   validation.Validate(cp.CPCodeID, validation.Required),
	}.Filter()
}

// Validate validates CreateCPCodeRequest
func (cp CreateCPCodeRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(cp.ContractID, validation.Required),
		"GroupID":    validation.Validate(cp.GroupID, validation.Required),
		"CPCode":     validation.Validate(cp.CPCode, validation.Required),
	}.Filter()
}

// Validate validates CreateCPCode
func (cp CreateCPCode) Validate() error {
	return validation.Errors{
		"ProductID":  validation.Validate(cp.ProductID, validation.Required),
		"CPCodeName": validation.Validate(cp.CPCodeName, validation.Required),
	}.Filter()
}

// GetCPCodes is used to list all available CP codes for given group and contract
func (p *papi) GetCPCodes(ctx context.Context, params GetCPCodesRequest) (*GetCPCodesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetCPCodes")

	getURL := fmt.Sprintf(
		"/papi/v1/cpcodes?contractId=%s&groupId=%s",
		params.ContractID,
		params.GroupID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcpcodes request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))
	var cpCodes GetCPCodesResponse
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("getcpcodes request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &cpCodes, nil
}

// GetCPCodes is used to fetch a CP code with provided ID
func (p *papi) GetCPCode(ctx context.Context, params GetCPCodeRequest) (*GetCPCodesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetCPCode")

	getURL := fmt.Sprintf("/papi/v1/cpcodes/%s?contractId=%s&groupId=%s", params.CPCodeID, params.ContractID, params.GroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcpcode request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))
	var cpCodes GetCPCodesResponse
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("getcpcode request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &cpCodes, nil
}

// CreateCPCode creates a new CP code with provided CreateCPCodeRequest data
func (p *papi) CreateCPCode(ctx context.Context, r CreateCPCodeRequest) (*CreateCPCodeResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("CreateCPCode")

	createURL := fmt.Sprintf("/papi/v1/cpcodes?contractId=%s&groupId=%s", r.ContractID, r.GroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create createcpcode request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))
	var createResponse CreateCPCodeResponse
	resp, err := p.Exec(req, &createResponse, r.CPCode)
	if err != nil {
		return nil, fmt.Errorf("getcpcode request failed: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, session.NewAPIError(resp, logger)
	}
	id, err := tools.FetchIDFromLocation(createResponse.CPCodeLink)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrInvalidLocation, err.Error())
	}
	createResponse.CPCodeID = id
	return &createResponse, nil
}
