package cloudcertificates

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

func (c *cloudcertificates) ListCertificateBindings(ctx context.Context, params ListCertificateBindingsRequest) (*ListCertificateBindingsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListCertificateBindings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListCertificateBindings, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v1/certificates/%s/certificate-bindings", params.CertificateID).
		AddQueryParamIf("pageSize", strconv.FormatInt(params.PageSize, 10), params.PageSize > 0).
		AddQueryParamIf("page", strconv.FormatInt(params.Page, 10), params.Page > 0).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListCertificateBindings, err)
	}

	var result ListCertificateBindingsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListCertificateBindings, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListCertificateBindings, c.Error(resp))
	}

	result.RateLimits = extractRateLimitHeaders(resp, logger)

	return &result, nil
}

func (c *cloudcertificates) ListBindings(ctx context.Context, params ListBindingsRequest) (*ListBindingsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListBindings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListBindings, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v1/certificate-bindings").
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParamIf("groupId", params.GroupID, params.GroupID != "").
		AddQueryParamFunc("expiringInDays", func() string {
			return strconv.FormatInt(*params.ExpiringInDays, 10)
		}, params.ExpiringInDays != nil).
		AddQueryParamIf("domain", params.Domain, params.Domain != "").
		AddQueryParamIf("network", string(params.Network), params.Network != "").
		AddQueryParamIf("pageSize", strconv.FormatInt(params.PageSize, 10), params.PageSize > 0).
		AddQueryParamIf("page", strconv.FormatInt(params.Page, 10), params.Page > 0).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListBindings, err)
	}

	var result ListBindingsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListBindings, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListBindings, c.Error(resp))
	}

	result.RateLimits = extractRateLimitHeaders(resp, logger)

	return &result, nil
}
