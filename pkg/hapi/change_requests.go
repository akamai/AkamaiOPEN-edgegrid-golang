package hapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// GetChangeRequest is a request struct
	GetChangeRequest struct {
		ChangeID int
	}

	// ChangeRequest represents change response from api
	ChangeRequest struct {
		Action            string         `json:"action"`
		ChangeID          int64          `json:"changeId"`
		Comments          string         `json:"comments"`
		EdgeHostnames     []EdgeHostname `json:"edgeHostnames"`
		Status            string         `json:"status"`
		StatusMessage     string         `json:"statusMessage"`
		StatusUpdateEmail string         `json:"statusUpdateEmail"`
		StatusUpdateDate  string         `json:"statusUpdateDate"`
		SubmitDate        string         `json:"submitDate"`
		Submitter         string         `json:"submitter"`
		SubmitterEmail    string         `json:"submitterEmail"`
	}
)

// ErrGetChangeRequest returned when get change request fails
var ErrGetChangeRequest = errors.New("get change request")

func (h *hapi) GetChangeRequest(ctx context.Context, prop GetChangeRequest) (*ChangeRequest, error) {
	logger := h.Log(ctx)
	logger.Debug("GetChangeRequest")

	uri := fmt.Sprintf("/hapi/v1/change-requests/%d", prop.ChangeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangeRequest, err)
	}

	var rval ChangeRequest

	resp, err := h.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangeRequest, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangeRequest, h.Error(resp))
	}

	return &rval, nil
}
