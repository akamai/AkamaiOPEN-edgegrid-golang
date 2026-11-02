package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
)

type (
	// ListContractsResponse represents a response object returned by ListContracts
	ListContractsResponse struct {
		ContractIDs []string `json:"contractIds"`
	}
)

var (
	// ErrListContracts is returned in case an error occurs on ListContracts operation
	ErrListContracts = errors.New("list contracts")
)

func (e *edgeworkers) ListContracts(ctx context.Context) (*ListContractsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListContracts")

	req, err := request.NewGet(ctx, "/edgeworkers/v1/contracts").Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListContracts, err)
	}

	var result ListContractsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListContracts, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListContracts, e.Error(resp))
	}

	return &result, nil
}
