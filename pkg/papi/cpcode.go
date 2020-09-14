package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/spf13/cast"
)

type (
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

var (
	// ErrGroupEmpty is returned when a required 'groupId' param is missing from the request
	ErrGroupEmpty = errors.New("provided group ID cannot be empty")
	// ErrContractEmpty is returned when a required 'contractId' param is missing from the request
	ErrContractEmpty = errors.New("provided contract ID cannot be empty")
	// ErrIDEmpty is returned when a required resource ID param is missing from the request
	ErrIDEmpty = errors.New("provided CP code ID cannot be empty")
	// ErrInvalidLocation is returned when there was an error while fetching ID from location response object
	ErrInvalidLocation = errors.New("response location URL is invalid")
)

// GetCPCodes is used to list all available CP codes for given group and contract
func (p *papi) GetCPCodes(ctx context.Context, params GetCPCodesRequest) (*GetCPCodesResponse, error) {
	if params.ContractID == "" {
		return nil, ErrContractEmpty
	}
	if params.GroupID == "" {
		return nil, ErrGroupEmpty
	}

	logger := p.Log(ctx)
	logger.Debug("GetCPCodes")

	getURL := fmt.Sprintf(
		"/papi/v1/cpcodes?contractId=%s&groupId=%s",
		params.ContractID,
		params.GroupID,
	)
	if len(params.Options) > 0 {
		getURL = fmt.Sprintf("%s&options=%s", getURL, strings.Join(params.Options, ","))
	}
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
	if params.ContractID == "" {
		return nil, ErrContractEmpty
	}
	if params.GroupID == "" {
		return nil, ErrGroupEmpty
	}
	if params.CPCodeID == "" {
		return nil, ErrIDEmpty
	}

	logger := p.Log(ctx)
	logger.Debug("GetCPCode")

	createURL := fmt.Sprintf("/papi/v1/cpcodes/%s?contractId=%s&groupId=%s", params.CPCodeID, params.ContractID, params.GroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, createURL, nil)
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
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, createURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &cpCodes, nil
}

// CreateCPCode creates a new CP code with provided CreateCPCodeRequest data
func (p *papi) CreateCPCode(ctx context.Context, r CreateCPCodeRequest) (*CreateCPCodeResponse, error) {
	if r.ContractID == "" {
		return nil, ErrContractEmpty
	}
	if r.GroupID == "" {
		return nil, ErrGroupEmpty
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
	id, err := fetchIDFromLocation(createResponse.CPCodeLink)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidLocation, err.Error())
	}
	createResponse.CPCodeID = id
	return &createResponse, nil
}

func fetchIDFromLocation(loc string) (string, error) {
	locURL, err := url.Parse(loc)
	if err != nil {
		return "", err
	}
	pathSplit := strings.Split(locURL.Path, "/")
	return pathSplit[len(pathSplit)-1], nil
}
