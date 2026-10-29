package hapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
)

func (h *hapi) DeleteEdgeHostname(ctx context.Context, params DeleteEdgeHostnameRequest) (*DeleteEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("DeleteEdgeHostname")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrDeleteEdgeHostname, ErrStructValidation, err)
	}

	req, err := request.NewDelete(ctx, "/hapi/v1/dns-zones/%s/edge-hostnames/%s", params.DNSZone, params.RecordName).
		AddQueryParamIf("statusUpdateEmail", strings.Join(params.StatusUpdateEmail, ","), len(params.StatusUpdateEmail) > 0).
		AddQueryParamIf("comments", params.Comments, params.Comments != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrDeleteEdgeHostname, err)
	}

	var rval DeleteEdgeHostnameResponse

	resp, err := h.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrDeleteEdgeHostname, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%w: %w", ErrDeleteEdgeHostname, h.Error(resp))
	}

	return &rval, nil
}

func (h *hapi) GetEdgeHostname(ctx context.Context, params GetEdgeHostnameRequest) (*GetEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("GetEdgeHostname")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetEdgeHostname, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/hapi/v1/edge-hostnames/%d", params.EdgeHostnameID).
		AddQueryParamFunc("includeDeliveryPolicy", func() string {
			return fmt.Sprintf("%t", *params.IncludeDeliveryPolicy)
		}, params.IncludeDeliveryPolicy != nil).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetEdgeHostname, err)
	}

	var rval GetEdgeHostnameResponse

	resp, err := h.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetEdgeHostname, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetEdgeHostname, h.Error(resp))
	}

	return &rval, nil
}

func (h *hapi) UpdateEdgeHostname(ctx context.Context, params UpdateEdgeHostnameRequest) (*UpdateEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("UpdateEdgeHostname")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrUpdateEdgeHostname, ErrStructValidation, err)
	}

	req, err := request.NewPatch(ctx, "/hapi/v1/dns-zones/%s/edge-hostnames/%s", params.DNSZone, params.RecordName).
		AddQueryParamIf("statusUpdateEmail", strings.Join(params.StatusUpdateEmail, ","), len(params.StatusUpdateEmail) > 0).
		AddQueryParamIf("comments", params.Comments, params.Comments != "").
		AddHeader("Content-Type", "application/json-patch+json").
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrUpdateEdgeHostname, err)
	}

	var result UpdateEdgeHostnameResponse

	resp, err := h.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrUpdateEdgeHostname, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%w: %w", ErrUpdateEdgeHostname, h.Error(resp))
	}

	return &result, nil
}

func (h *hapi) GetCertificate(ctx context.Context, params GetCertificateRequest) (*GetCertificateResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("GetCertificate")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetCertificate, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/hapi/v1/dns-zones/%s/edge-hostnames/%s/certificate", params.DNSZone, params.RecordName).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetCertificate, err)
	}

	var result GetCertificateResponse

	resp, err := h.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetCertificate, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetCertificate, ErrNotFound, h.Error(resp))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetCertificate, h.Error(resp))
	}

	return &result, nil
}
