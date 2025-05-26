package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Contract contains contractID and a list of currently assigned Akamai authoritative nameservers
	Contract struct {
		ContractID  string   `json:"contractId"`
		Authorities []string `json:"authorities"`
	}
	// AuthorityResponse contains response with a list of one or more Contracts
	AuthorityResponse struct {
		Contracts []Contract `json:"contracts"`
	}
	// GetAuthoritiesRequest contains request parameters for GetAuthorities
	GetAuthoritiesRequest struct {
		ContractIDs string
	}

	// GetAuthoritiesResponse contains the response data from GetAuthorities operation
	GetAuthoritiesResponse struct {
		Contracts []Contract `json:"contracts"`
	}

	// GetNameServerRecordListRequest contains request parameters for GetNameServerRecordList
	GetNameServerRecordListRequest struct {
		ContractIDs string
	}
)

var (
	// ErrGetAuthorities is returned when GetAuthorities fails
	ErrGetAuthorities = errors.New("get authorities")
	// ErrGetNameServerRecordList is returned when GetNameServerRecordList fails
	ErrGetNameServerRecordList = errors.New("get name server record list")
)

// Validate validates GetAuthoritiesRequest
func (r GetAuthoritiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ContractIDs": validation.Validate(r.ContractIDs, validation.Required),
	})
}

// Validate validates GetNameServerRecordListRequest
func (r GetNameServerRecordListRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ContractIDs": validation.Validate(r.ContractIDs, validation.Required),
	})
}

func (d *dns) GetAuthorities(ctx context.Context, params GetAuthoritiesRequest) (*GetAuthoritiesResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("GetAuthorities")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAuthorities, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-dns/v2/data/authorities?contractIds=%s", params.ContractIDs)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getauthorities request: %w", err)
	}

	var result GetAuthoritiesResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAuthorities request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, d.Error(resp)
	}

	return &result, nil
}

func (d *dns) GetNameServerRecordList(ctx context.Context, params GetNameServerRecordListRequest) ([]string, error) {
	logger := d.Log(ctx)
	logger.Debug("GetNameServerRecordList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetNameServerRecordList, ErrStructValidation, err)
	}

	NSrecords, err := d.GetAuthorities(ctx, GetAuthoritiesRequest(params))
	if err != nil {
		return nil, err
	}

	var result []string
	for _, r := range NSrecords.Contracts {
		result = append(result, r.Authorities...)
	}

	return result, nil
}
