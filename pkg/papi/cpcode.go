package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
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

	// AccessGroup contains the access group information for a CP code details.
	AccessGroup struct {
		// ContractID identifies the contract assigned to the access control group.
		ContractID string `json:"contractId"`

		// GroupID identifies the access control group.
		GroupID *int64 `json:"groupId"`
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
		AccessGroup      AccessGroup      `json:"accessGroup"`
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
	logger := p.Log(ctx)
	logger.Debug("GetCPCodes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCPCodes, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/papi/v1/cpcodes").
		AddQueryParam("groupId", params.GroupID).
		AddQueryParam("contractId", params.ContractID).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCPCodes, err)
	}

	var cpCodes GetCPCodesResponse
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCPCodes, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCPCodes, p.Error(resp))
	}

	return &cpCodes, nil
}

// GetCPCode is used to fetch a CP code with provided ID
func (p *papi) GetCPCode(ctx context.Context, params GetCPCodeRequest) (*GetCPCodesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetCPCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCPCode, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/papi/v1/cpcodes/%s", params.CPCodeID).
		AddQueryParam("groupId", params.GroupID).
		AddQueryParam("contractId", params.ContractID).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCPCode, err)
	}

	var cpCodes GetCPCodesResponse
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCPCode, err)
	}
	defer session.CloseResponseBody(resp)

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
//
// Deprecated: Use reportinggroups.GetCPCode instead.
func (p *papi) GetCPCodeDetail(ctx context.Context, ID int) (*CPCodeDetailResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetCPCodeDetail")

	req, err := request.NewGet(ctx, "/cprg/v1/cpcodes/%d", ID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCPCodeDetail, err)
	}

	var cpCodeDetail CPCodeDetailResponse
	resp, err := p.Exec(req, &cpCodeDetail)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCPCodeDetail, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCPCodeDetail, p.Error(resp))
	}

	return &cpCodeDetail, nil
}

// CreateCPCode creates a new CP code with provided CreateCPCodeRequest data
func (p *papi) CreateCPCode(ctx context.Context, params CreateCPCodeRequest) (*CreateCPCodeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateCPCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", ErrCreateCPCode, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/papi/v1/cpcodes").
		AddQueryParam("groupId", params.GroupID).
		AddQueryParam("contractId", params.ContractID).
		WithBody(params.CPCode).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateCPCode, err)
	}

	var createResponse CreateCPCodeResponse
	resp, err := p.Exec(req, &createResponse)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateCPCode, err)
	}
	defer session.CloseResponseBody(resp)

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
//
// Deprecated: Use reportinggroups.UpdateCPCode instead.
func (p *papi) UpdateCPCode(ctx context.Context, r UpdateCPCodeRequest) (*CPCodeDetailResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateCPCode")

	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", ErrUpdateCPCode, ErrStructValidation, err)
	}

	req, err := request.NewPut(ctx, "/cprg/v1/cpcodes/%d", r.ID).
		WithBody(r).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateCPCode, err)
	}

	var cpCodeDetail CPCodeDetailResponse
	resp, err := p.Exec(req, &cpCodeDetail)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateCPCode, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateCPCode, p.Error(resp))
	}

	return &cpCodeDetail, nil
}
