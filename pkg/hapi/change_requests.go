package hapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

func (h *hapi) GetChangeRequest(ctx context.Context, params GetChangeRequestRequest) (*GetChangeRequestResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("GetChangeRequest")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetChangeRequest, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/hapi/v1/change-requests/%d", params.ChangeID).
		AddQueryParamFunc("includeDeliveryPolicy", func() string {
			return strconv.FormatBool(*params.IncludeDeliveryPolicy)
		}, params.IncludeDeliveryPolicy != nil).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetChangeRequest, err)
	}

	var response GetChangeRequestResponse

	resp, err := h.Exec(req, &response)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetChangeRequest, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetChangeRequest, h.Error(resp))
	}

	return &response, nil
}

func (h *hapi) ListChangeRequests(ctx context.Context, params ListChangeRequestsRequest) (*ListChangeRequestsResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("ListChangeRequests")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListChangeRequests, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/hapi/v1/change-requests").
		AddQueryParamIf("status", string(params.Status), params.Status != "").
		AddQueryParamFunc("includeDeliveryPolicy", func() string {
			return strconv.FormatBool(*params.IncludeDeliveryPolicy)
		}, params.IncludeDeliveryPolicy != nil).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListChangeRequests, err)
	}

	var result ListChangeRequestsResponse
	resp, err := h.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListChangeRequests, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListChangeRequests, h.Error(resp))
	}

	return &result, nil
}

func (h *hapi) ListChangeRequestsPerEdgeHostname(ctx context.Context, params ListChangeRequestsPerEdgeHostnameRequest) (*ListChangeRequestsPerEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("ListChangeRequestsPerEdgeHostname")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListChangeRequestsPerEdgeHostname, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/hapi/v1/dns-zones/%s/edge-hostnames/%s/change-requests", url.PathEscape(params.DNSZone), url.PathEscape(params.RecordName)).
		AddQueryParamFunc("includeDeliveryPolicy", func() string {
			return strconv.FormatBool(*params.IncludeDeliveryPolicy)
		}, params.IncludeDeliveryPolicy != nil).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListChangeRequestsPerEdgeHostname, err)
	}

	var result ListChangeRequestsPerEdgeHostnameResponse
	resp, err := h.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListChangeRequestsPerEdgeHostname, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListChangeRequestsPerEdgeHostname, h.Error(resp))
	}

	return &result, nil
}
