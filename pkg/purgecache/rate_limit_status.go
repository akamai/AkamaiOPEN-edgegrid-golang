package purgecache

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// RateLimitStatusRequest contains request parameters for checking rate limit status.
	RateLimitStatusRequest struct {
		// PurgeType specifies the type of purge objects. Must be one of: cpcode, tag, url.
		PurgeType PurgeType

		// Network specifies the network environment: staging or production.
		// When omitted, the API defaults to production.
		Network PurgeNetwork
	}

	// RateLimitStatusResponse contains the response for a rate limit status check.
	RateLimitStatusResponse struct {
		// DescribedBy is a URL that describes the API's error response.
		DescribedBy string `json:"describedBy"`

		// Detail contains detailed information about the HTTP status code returned.
		Detail string `json:"detail"`

		// HTTPStatus is the HTTP code that indicates the status of the rate limit status request.
		HTTPStatus int `json:"httpStatus"`

		// SupportID is an identifier to provide Akamai Technical Support if issues arise.
		SupportID string `json:"supportId"`

		// Title describes the response type.
		Title string `json:"title"`

		// RateLimitHeaders contains rate limit information extracted from response headers.
		RateLimitHeaders RateLimitHeaders `json:"-"`
	}
)

// Validate validates RateLimitStatusRequest.
func (r RateLimitStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PurgeType": validation.Validate(r.PurgeType, validation.Required),
		"Network":   validation.Validate(r.Network),
	})
}

// RateLimitStatus checks the rate and object limit statuses for a specific purge type.
//
// See: https://techdocs.akamai.com/purge-cache/reference/post-rate-limit-status
func (p *purgecache) RateLimitStatus(ctx context.Context, params RateLimitStatusRequest) (*RateLimitStatusResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RateLimitStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrRateLimitStatus, ErrStructValidation, err)
	}

	var err error
	var req *http.Request
	if params.Network == "" {
		req, err = request.NewPost(ctx, "/ccu/v3/rate-limit-status/%s", params.PurgeType).Build()
	} else {
		req, err = request.NewPost(ctx, "/ccu/v3/rate-limit-status/%s/%s", params.PurgeType, params.Network).Build()
	}
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrRateLimitStatus, err)
	}

	var result RateLimitStatusResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrRateLimitStatus, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	result.RateLimitHeaders = extractRateLimitStatusHeaders(resp, logger)

	return &result, nil
}
