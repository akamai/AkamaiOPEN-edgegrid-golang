package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// CPCodes contains operations available on CPCode resource
	CPCodes interface {
		// GetCPCodes lists all available CP codes
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-cpcodes
		GetCPCodes(context.Context, GetCPCodesRequest) (*GetCPCodesResponse, error)

		// GetCPCode gets the CP code with provided ID
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-cpcode
		GetCPCode(context.Context, GetCPCodeRequest) (*GetCPCodesResponse, error)

		// GetCPCodeDetail lists detailed information about a specific CP code
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/get-cpcode
		GetCPCodeDetail(context.Context, int) (*CPCodeDetailResponse, error)

		// CreateCPCode creates a new CP code
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-cpcodes
		CreateCPCode(context.Context, CreateCPCodeRequest) (*CreateCPCodeResponse, error)

		// UpdateCPCode modifies a specific CP code. You should only modify a CP code's name, time zone, and purgeable member
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/put-cpcode
		UpdateCPCode(context.Context, UpdateCPCodeRequest) (*CPCodeDetailResponse, error)
	}

	// CPCode contains CP code resource data
	CPCode struct {
		ID          string   `json:"cpcodeId"`
		Name        string   `json:"cpcodeName"`
		CreatedDate string   `json:"createdDate"`
		ProductIDs  []string `json:"productIds"`
	}

	// CPCodeContract contains contract data used in CPRG API calls
	CPCodeContract struct {
		ContractID string `json:"contractId"`
		Status     string `json:"status,omitempty"`
	}

	// CPCodeDetailResponse is a response returned while fetching CP code details using CPRG API call
	CPCodeDetailResponse struct {
		ID               int              `json:"cpcodeId"`
		Name             string           `json:"cpcodeName"`
		Purgeable        bool             `json:"purgeable"`
		AccountID        string           `json:"accountId"`
		DefaultTimeZone  string           `json:"defaultTimezone"`
		OverrideTimeZone CPCodeTimeZone   `json:"overrideTimezone"`
		Type             string           `json:"type"`
		Contracts        []CPCodeContract `json:"contracts"`
		Products         []CPCodeProduct  `json:"products"`
	}

	// CPCodeItems contains a list of CPCode items
	CPCodeItems struct {
		Items []CPCode `json:"items"`
	}

	// CPCodeProduct contains product data used in CPRG API calls
	CPCodeProduct struct {
		ProductID   string `json:"productId"`
		ProductName string `json:"productName,omitempty"`
	}

	// GetCPCodesResponse is a response returned while fetching CP codes
	GetCPCodesResponse struct {
		AccountID  string      `json:"accountId"`
		ContractID string      `json:"contractId"`
		GroupID    string      `json:"groupId"`
		CPCodes    CPCodeItems `json:"cpcodes"`
		CPCode     CPCode
	}

	// CPCodeTimeZone contains time zone data used in CPRG API calls
	CPCodeTimeZone struct {
		TimeZoneID    string `json:"timezoneId"`
		TimeZoneValue string `json:"timezoneValue,omitempty"`
	}

	// CreateCPCodeRequest contains data required to create CP code (both request body and group/contract information
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

	// GetCPCodesRequest contains parameters required to list/create CP codes
	// GroupID and ContractID are required as part of every CP code operation, ID is required only for operating on specific CP code
	GetCPCodesRequest struct {
		ContractID string
		GroupID    string
	}

	// UpdateCPCodeRequest contains parameters required to update CP code, using CPRG API call
	UpdateCPCodeRequest struct {
		ID               int              `json:"cpcodeId"`
		Name             string           `json:"cpcodeName"`
		Purgeable        *bool            `json:"purgeable,omitempty"`
		OverrideTimeZone *CPCodeTimeZone  `json:"overrideTimezone,omitempty"`
		Contracts        []CPCodeContract `json:"contracts"`
		Products         []CPCodeProduct  `json:"products"`
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

// Validate validates CPCodeContract
func (contract CPCodeContract) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(contract.ContractID, validation.Required),
	}.Filter()
}

// Validate validates CPCodeProduct
func (product CPCodeProduct) Validate() error {
	return validation.Errors{
		"ProductID": validation.Validate(product.ProductID, validation.Required),
	}.Filter()
}

// Validate validates CPCodeTimeZone
func (timeZone CPCodeTimeZone) Validate() error {
	return validation.Errors{
		"TimeZoneID": validation.Validate(timeZone.TimeZoneID, validation.Required),
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

// Validate validates UpdateCPCodeRequest
func (cp UpdateCPCodeRequest) Validate() error {
	return validation.Errors{
		"ID":               validation.Validate(cp.ID, validation.Required),
		"Name":             validation.Validate(cp.Name, validation.Required),
		"Contracts":        validation.Validate(cp.Contracts, validation.Required),
		"Products":         validation.Validate(cp.Products, validation.Required),
		"OverrideTimeZone": validation.Validate(cp.OverrideTimeZone),
	}.Filter()
}

var (
	// ErrGetCPCodes represents error when fetching CP Codes fails
	ErrGetCPCodes = errors.New("fetching CP Codes")
	// ErrGetCPCode represents error when fetching CP Code fails
	ErrGetCPCode = errors.New("fetching CP Code")
	// ErrGetCPCodeDetail represents error when fetching CP Code Details fails
	ErrGetCPCodeDetail = errors.New("fetching CP Code Detail")
	// ErrCreateCPCode represents error when creating CP Code fails
	ErrCreateCPCode = errors.New("creating CP Code")
	// ErrUpdateCPCode represents error when updating CP Code
	ErrUpdateCPCode = errors.New("updating CP Code")
)

// GetCPCodes is used to list all available CP codes for given group and contract
func (p *papi) GetCPCodes(ctx context.Context, params GetCPCodesRequest) (*GetCPCodesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCPCodes, ErrStructValidation, err)
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
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCPCodes, err)
	}

	var cpCodes GetCPCodesResponse
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCPCodes, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCPCodes, p.Error(resp))
	}

	return &cpCodes, nil
}

// GetCPCode is used to fetch a CP code with provided ID
func (p *papi) GetCPCode(ctx context.Context, params GetCPCodeRequest) (*GetCPCodesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCPCode, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetCPCode")

	getURL := fmt.Sprintf("/papi/v1/cpcodes/%s?contractId=%s&groupId=%s", params.CPCodeID, params.ContractID, params.GroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCPCode, err)
	}

	var cpCodes GetCPCodesResponse
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCPCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCPCode, p.Error(resp))
	}
	if len(cpCodes.CPCodes.Items) == 0 {
		return nil, fmt.Errorf("%s: %w: CPCodeID: %s", ErrGetCPCode, ErrNotFound, params.CPCodeID)
	}
	cpCodes.CPCode = cpCodes.CPCodes.Items[0]

	return &cpCodes, nil
}

// GetCPCodeDetail is used to fetch CP code detail with provided ID using CPRG API
func (p *papi) GetCPCodeDetail(ctx context.Context, ID int) (*CPCodeDetailResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetCPCodeDetail")

	getURL := fmt.Sprintf("/cprg/v1/cpcodes/%d", ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCPCodeDetail, err)
	}

	var cpCodeDetail CPCodeDetailResponse
	resp, err := p.Exec(req, &cpCodeDetail)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCPCodeDetail, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCPCodeDetail, p.Error(resp))
	}

	return &cpCodeDetail, nil
}

// CreateCPCode creates a new CP code with provided CreateCPCodeRequest data
func (p *papi) CreateCPCode(ctx context.Context, r CreateCPCodeRequest) (*CreateCPCodeResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", ErrCreateCPCode, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("CreateCPCode")

	createURL := fmt.Sprintf("/papi/v1/cpcodes?contractId=%s&groupId=%s", r.ContractID, r.GroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateCPCode, err)
	}

	var createResponse CreateCPCodeResponse
	resp, err := p.Exec(req, &createResponse, r.CPCode)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateCPCode, err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateCPCode, p.Error(resp))
	}
	id, err := ResponseLinkParse(createResponse.CPCodeLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateCPCode, ErrInvalidResponseLink, err)
	}
	createResponse.CPCodeID = id
	return &createResponse, nil
}

// UpdateCPCode is used to update CP code using CPRG API
func (p *papi) UpdateCPCode(ctx context.Context, r UpdateCPCodeRequest) (*CPCodeDetailResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", ErrUpdateCPCode, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateCPCode")

	updateURL := fmt.Sprintf("/cprg/v1/cpcodes/%d", r.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, updateURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateCPCode, err)
	}

	var cpCodeDetail CPCodeDetailResponse
	resp, err := p.Exec(req, &cpCodeDetail, r)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateCPCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateCPCode, p.Error(resp))
	}

	return &cpCodeDetail, nil
}
