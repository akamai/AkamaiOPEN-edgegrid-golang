package purgecache

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PurgeType represents the type of content to purge.
	PurgeType string

	// PurgeNetwork represents the network environment to target.
	PurgeNetwork string

	// PurgeByURLRequest is the common request struct for purge operations by URL or ARL.
	PurgeByURLRequest struct {
		// Network specifies the network environment: `staging` or `production`.
		// When omitted, the API defaults to production.
		Network PurgeNetwork `json:"-"`

		// Objects is the list of URLs or ARLs to purge.
		Objects []string `json:"objects"`
	}

	// PurgeByTagRequest is the common request struct for purge operations by cache tag.
	PurgeByTagRequest struct {
		// Network specifies the network environment: `staging` or `production`.
		// When omitted, the API defaults to production.
		Network PurgeNetwork `json:"-"`

		// Objects is the list of cache tags to purge.
		Objects []string `json:"objects"`
	}

	// PurgeByCPCodeRequest is the common request struct for purge operations by CP Code.
	PurgeByCPCodeRequest struct {
		// Network specifies the network environment: `staging` or `production`.
		// When omitted, the API defaults to production.
		Network PurgeNetwork `json:"-"`

		// Objects is the list of CP Codes to purge.
		Objects []int64 `json:"objects"`
	}

	// PurgeResponse is the common response struct for cache purge operations.
	PurgeResponse struct {
		// DescribedBy is a URL that describes the API's error response.
		DescribedBy string `json:"describedBy"`

		// Detail contains detailed information about the HTTP status code returned.
		Detail string `json:"detail"`

		// EstimatedSeconds is the estimated number of seconds before the purge is complete.
		EstimatedSeconds int64 `json:"estimatedSeconds"`

		// HTTPStatus is the HTTP code that indicates the status of the purge request.
		HTTPStatus int `json:"httpStatus"`

		// PurgeID is the unique identifier for the purge request.
		PurgeID string `json:"purgeId"`

		// SupportID is an identifier to provide Akamai Technical Support if issues arise.
		SupportID string `json:"supportId"`

		// Title describes the response type.
		Title string `json:"title"`

		// RateLimitHeaders contains rate limit information extracted from response headers.
		RateLimitHeaders RateLimitHeaders `json:"-"`
	}

	// RateLimitHeaders contains rate limit information extracted from response headers.
	RateLimitHeaders struct {
		// Limit is the burst limit specifying how many requests you can make in excess of the rate limit.
		// Nil if the header is absent.
		Limit *int64

		// LimitObjects is the burst limit specifying how many objects you can add in excess of the rate limit.
		// Nil if the header is absent.
		LimitObjects *int64

		// LimitPerSecond is the maximum sustainable number of requests allowed each second per account.
		// Nil if the header is absent.
		LimitPerSecond *float64

		// LimitPerSecondObjects is the maximum sustainable number of objects allowed each second per request.
		// Nil if the header is absent.
		LimitPerSecondObjects *float64

		// Remaining is the number of requests left before exceeding the limit.
		// Nil if the header is absent.
		Remaining *int64

		// RemainingObjects is the number of objects left before exceeding the limit.
		// Nil if the header is absent.
		RemainingObjects *int64

		// SecondsToRefreshLimit is the minimum number of seconds to wait before sending a request burst.
		// Nil if the header is absent.
		SecondsToRefreshLimit *float64

		// SecondsToRefreshLimitObjects is the minimum number of seconds to wait before sending an object burst.
		// Nil if the header is absent.
		SecondsToRefreshLimitObjects *float64
	}
)

const (
	// PurgeTypeCPCode purges by CP code.
	PurgeTypeCPCode PurgeType = "cpcode"

	// PurgeTypeTag purges by cache tag.
	PurgeTypeTag PurgeType = "tag"

	// PurgeTypeURL purges by URL or ARL.
	PurgeTypeURL PurgeType = "url"

	// PurgeNetworkStaging targets the staging network.
	PurgeNetworkStaging PurgeNetwork = "staging"

	// PurgeNetworkProduction targets the production network.
	PurgeNetworkProduction PurgeNetwork = "production"
)

// Validate implements the validation.Validatable interface.
// It checks that the PurgeType is one of the valid values: `cpcode`, `tag`, `url`.
func (p PurgeType) Validate() error {
	return validation.In(PurgeTypeCPCode, PurgeTypeTag, PurgeTypeURL).Error(
		fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s'",
			p, PurgeTypeCPCode, PurgeTypeTag, PurgeTypeURL)).Validate(p)
}

// Validate implements the validation.Validatable interface.
// An empty value is valid (network is optional), but equal to `production`; non-empty values must be `staging` or `production`.
func (n PurgeNetwork) Validate() error {
	return validation.In(PurgeNetworkStaging, PurgeNetworkProduction).Error(
		fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'",
			n, PurgeNetworkStaging, PurgeNetworkProduction)).Validate(n)
}

func extractRateLimitHeaders(resp *http.Response, logger log.Interface) RateLimitHeaders {
	var h RateLimitHeaders

	parseIntHeader := func(name string) *int64 {
		val := resp.Header.Get(name)
		if val == "" {
			return nil
		}
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			logger.Warnf("failed to parse %s header: %v", name, err)
			return nil
		}
		return &n
	}

	parseFloatHeader := func(name string) *float64 {
		val := resp.Header.Get(name)
		if val == "" {
			return nil
		}
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			logger.Warnf("failed to parse %s header: %v", name, err)
			return nil
		}
		return &f
	}

	h.Limit = parseIntHeader("X-Ratelimit-Limit")
	h.LimitObjects = parseIntHeader("X-Ratelimit-Limit-Objects")
	h.LimitPerSecond = parseFloatHeader("X-Ratelimit-Limit-Per-Second")
	h.LimitPerSecondObjects = parseFloatHeader("X-Ratelimit-Limit-Per-Second-Objects")
	h.Remaining = parseIntHeader("X-Ratelimit-Remaining")
	h.RemainingObjects = parseIntHeader("X-Ratelimit-Remaining-Objects")
	h.SecondsToRefreshLimit = parseFloatHeader("X-Ratelimit-Seconds-To-Refresh-Limit")
	h.SecondsToRefreshLimitObjects = parseFloatHeader("X-Ratelimit-Seconds-To-Refresh-Limit-Objects")

	return h
}

// performPurgeRequest builds and executes a purge request for the given action and purgeType.
func (p *purgecache) performPurgeRequest(ctx context.Context, wrapErr error, action string, purgeType PurgeType, network PurgeNetwork, body any) (*PurgeResponse, error) {
	var err error
	var req *http.Request
	if network == "" {
		req, err = request.NewPost(ctx, "/ccu/v3/%s/%s", action, purgeType).WithBody(body).Build()
	} else {
		req, err = request.NewPost(ctx, "/ccu/v3/%s/%s/%s", action, purgeType, network).WithBody(body).Build()
	}
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", wrapErr, err)
	}

	var result PurgeResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", wrapErr, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	result.RateLimitHeaders = extractRateLimitHeaders(resp, p.Log(ctx))
	return &result, nil
}
