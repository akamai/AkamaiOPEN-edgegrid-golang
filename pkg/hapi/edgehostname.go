package hapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

func (h *hapi) DeleteEdgeHostname(ctx context.Context, params DeleteEdgeHostnameRequest) (*DeleteEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("DeleteEdgeHostname")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteEdgeHostname, ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/hapi/v1/dns-zones/%s/edge-hostnames/%s",
		params.DNSZone,
		params.RecordName,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeleteEdgeHostname, err)
	}

	q := req.URL.Query()
	if len(params.StatusUpdateEmail) > 0 {
		emails := strings.Join(params.StatusUpdateEmail, ",")
		q.Add("statusUpdateEmail", emails)
	}
	if params.Comments != "" {
		q.Add("comments", params.Comments)
	}
	req.URL.RawQuery = q.Encode()

	var rval DeleteEdgeHostnameResponse

	resp, err := h.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeleteEdgeHostname, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrDeleteEdgeHostname, h.Error(resp))
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

func (h *hapi) UpdateEdgeHostname(ctx context.Context, request UpdateEdgeHostnameRequest) (*UpdateEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("UpdateEdgeHostname")

	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s: %s", ErrUpdateEdgeHostname, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s", request.DNSZone, request.RecordName)

	body, err := buildBody(request.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request body", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, uri, body)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if len(request.StatusUpdateEmail) > 0 {
		emails := strings.Join(request.StatusUpdateEmail, ",")
		q.Add("statusUpdateEmail", emails)
	}
	if request.Comments != "" {
		q.Add("comments", request.Comments)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json-patch+json")

	var result UpdateEdgeHostnameResponse

	resp, err := h.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateEdgeHostname, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%w: %s", ErrUpdateEdgeHostname, h.Error(resp))
	}

	return &result, nil
}

func (h *hapi) GetCertificate(ctx context.Context, params GetCertificateRequest) (*GetCertificateResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("GetCertificate")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCertificate, ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/hapi/v1/dns-zones/%s/edge-hostnames/%s/certificate",
		params.DNSZone,
		params.RecordName,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCertificate, err)
	}

	var result GetCertificateResponse

	resp, err := h.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCertificate, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s: %s: %w", ErrGetCertificate, ErrNotFound, h.Error(resp))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCertificate, h.Error(resp))
	}

	return &result, nil
}

func buildBody(body []UpdateEdgeHostnameRequestBody) (io.Reader, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(reqBody), nil
}
