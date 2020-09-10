package papi

import (
	"context"
	"errors"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/spf13/cast"
	"net/http"
)

type (
	CPCode struct {
		ID          string   `json:"cpcodeId"`
		Name        string   `json:"cpcodeName"`
		CreatedDate string   `json:"createdDate"`
		ProductId   string   `json:"productId"`
		ProductIds  []string `json:"productIds"`
	}

	CPCodeItems struct {
		Items []CPCode `json:"items"`
	}

	GetCPCodesResponse struct {
		AccountID  string      `json:"accountId"`
		ContractID string      `json:"contractId"`
		GroupId    string      `json:"groupId"`
		CPCodes    CPCodeItems `json:"cpcodes"`
	}
)

var (
	ErrGroupEmpty    = errors.New("provided group ID cannot be empty")
	ErrContractEmpty = errors.New("provided contract ID cannot be empty")
	ErrIDEmpty       = errors.New("provided CP code ID cannot be empty")
)

func (p *papi) GetCPCodes(ctx context.Context, contractID, groupID string) (*GetCPCodesResponse, error) {
	if contractID == "" {
		return nil, ErrContractEmpty
	}
	if groupID == "" {
		return nil, ErrGroupEmpty
	}
	var cpCodes GetCPCodesResponse

	logger := p.Log(ctx)
	logger.Debug("GetCPCodes")

	url := fmt.Sprintf("/papi/v1/cpcodes?contractId=%s&groupId=%s", contractID, groupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcpcodes request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(UsePrefixes))
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("getcpcodes request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &cpCodes, nil
}

func (p *papi) GetCPCode(ctx context.Context, id, contractID, groupID string) (*GetCPCodesResponse, error) {
	if contractID == "" {
		return nil, ErrContractEmpty
	}
	if groupID == "" {
		return nil, ErrGroupEmpty
	}
	if id == "" {
		return nil, ErrIDEmpty
	}
	var cpCodes GetCPCodesResponse

	logger := p.Log(ctx)
	logger.Debug("GetCPCode")

	url := fmt.Sprintf("/papi/v1/cpcodes/%s?contractId=%s&groupId=%s", id, contractID, groupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcpcode request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(UsePrefixes))
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("getcpcode request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &cpCodes, nil
}
