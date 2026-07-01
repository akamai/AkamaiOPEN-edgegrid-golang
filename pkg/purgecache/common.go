package purgecache

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PurgeType represents the type of content to purge.
	PurgeType string

	// PurgeNetwork represents the network environment to target.
	PurgeNetwork string

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

func extractRateLimitStatusHeaders(resp *http.Response, logger log.Interface) RateLimitHeaders {
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
