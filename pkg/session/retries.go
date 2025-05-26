package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/log"
	"github.com/hashicorp/go-retryablehttp"
)

// RetryConfig struct contains http retry configuration.
//
// ExcludedEndpoints field is a list of  shell patterns.
// The pattern syntax is:
//
//	pattern:
//		{ term }
//	term:
//		'*'         matches any sequence of non-/ characters
//		'?'         matches any single non-/ character
//		'[' [ '^' ] { character-range } ']'
//		            character class (must be non-empty)
//		c           matches character c (c != '*', '?', '\\', '[')
//		'\\' c      matches character c
//
//	character-range:
//		c           matches character c (c != '\\', '-', ']')
//		'\\' c      matches character c
//		lo '-' hi   matches character c for lo <= c <= hi
type RetryConfig struct {
	RetryMax          int
	RetryWaitMin      time.Duration
	RetryWaitMax      time.Duration
	ExcludedEndpoints []string
}

// NewRetryConfig creates a new retry config with default settings.
func NewRetryConfig() RetryConfig {
	return RetryConfig{
		RetryMax:          10,
		RetryWaitMin:      1 * time.Second,
		RetryWaitMax:      30 * time.Second,
		ExcludedEndpoints: []string{},
	}
}

func configureRetryClient(conf RetryConfig, signFunc func(r *http.Request) error, log log.Interface) (*retryablehttp.Client, error) {
	retryClient := retryablehttp.NewClient()

	err := validateRetryConf(conf)
	if err != nil {
		return nil, err
	}
	retryClient.RetryMax = conf.RetryMax
	retryClient.RetryWaitMin = conf.RetryWaitMin
	retryClient.RetryWaitMax = conf.RetryWaitMax

	retryClient.PrepareRetry = signFunc
	retryClient.HTTPClient.CheckRedirect = func(r *http.Request, _ []*http.Request) error {
		return signFunc(r)
	}
	retryClient.CheckRetry = overrideRetryPolicy(retryablehttp.DefaultRetryPolicy, conf.ExcludedEndpoints)
	retryClient.Backoff = overrideBackoff(retryablehttp.DefaultBackoff, log)
	retryClient.Logger = GetRetryableLogger(log)

	return retryClient, err
}

func validateRetryConf(conf RetryConfig) error {
	errs := []error{}

	if conf.RetryMax < 0 {
		errs = append(errs, errors.New("maximum number of retries cannot be negative"))
	}
	if conf.RetryWaitMin < 0 {
		errs = append(errs, errors.New("minimum retry wait time cannot be negative"))
	}
	if conf.RetryWaitMax < 0 {
		errs = append(errs, errors.New("maximum retry wait time cannot be negative"))
	}
	if conf.RetryWaitMax < conf.RetryWaitMin {
		errs = append(errs, errors.New("maximum retry wait time cannot be shorter than minimum retry wait time"))
	}
	for _, pattern := range conf.ExcludedEndpoints {
		if _, err := path.Match(pattern, ""); err != nil {
			errs = append(errs, fmt.Errorf("malformed exclude endpoint pattern: %v: %s", err, pattern))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func overrideRetryPolicy(basePolicy retryablehttp.CheckRetry, excludedEndpoints []string) retryablehttp.CheckRetry {
	return func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// do not retry on context.Canceled or context.DeadlineExceeded
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		if resp == nil || resp.Request.Method != http.MethodGet ||
			(resp.Request.URL != nil && isBlocked(resp.Request.URL.Path, excludedEndpoints)) {
			var urlErr *url.Error
			if resp == nil && errors.As(err, &urlErr) && strings.ToUpper(urlErr.Op) == http.MethodGet {
				return basePolicy(ctx, resp, err)
			}
			return false, err
		}

		// Retry all PAPI GET requests resulting status code 429
		// The backoff time is calculated in getXRateLimitBackoff
		is429 := resp.StatusCode == http.StatusTooManyRequests
		if is429 && (resp.Request.URL != nil && strings.HasPrefix(resp.Request.URL.Path, "/papi/")) {
			return true, nil
		}
		if resp.StatusCode == http.StatusConflict {
			return true, nil
		}
		return basePolicy(ctx, resp, err)
	}
}

func overrideBackoff(baseBackoff retryablehttp.Backoff, logger log.Interface) retryablehttp.Backoff {
	return func(minT, maxT time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests {
				if wait, ok := getXRateLimitBackoff(resp, logger); ok {
					return wait
				}
			}
		}
		return baseBackoff(minT, maxT, attemptNum, resp)
	}
}

// Note that Date's resolution is seconds (e.g. Mon, 01 Jul 2024 14:32:14 GMT),
// while X-RateLimit-Next's resolution is milliseconds (2024-07-01T14:32:28.645Z).
// This may cause the wait time to be inflated by at most one second, like for the
// actual server response time around 2024-07-01T14:32:14.999Z. This is acceptable behavior
// as retry does not occur earlier than expected.
func getXRateLimitBackoff(resp *http.Response, logger log.Interface) (time.Duration, bool) {
	nextHeader := resp.Header.Get("X-RateLimit-Next")
	if nextHeader == "" {
		return 0, false
	}
	next, err := time.Parse(time.RFC3339Nano, nextHeader)
	if err != nil {
		if logger != nil {
			logger.Error("Could not parse X-RateLimit-Next header", "error", err)
		}
		return 0, false
	}

	dateHeader := resp.Header.Get("Date")
	if dateHeader == "" {
		if logger != nil {
			logger.Warnf("No Date header for X-RateLimit-Next: %s", nextHeader)
		}
		return 0, false
	}
	date, err := time.Parse(time.RFC1123, dateHeader)
	if err != nil {
		if logger != nil {
			logger.Error("Could not parse Date header", "error", err)
		}
		return 0, false
	}

	// Next in the past does not make sense
	if next.Before(date) {
		if logger != nil {
			logger.Warnf("X-RateLimit-Next: %s before Date: %s", nextHeader, dateHeader)
		}
		return 0, false
	}
	return next.Sub(date), true
}

func isBlocked(url string, disabledPatterns []string) bool {
	for _, pattern := range disabledPatterns {
		match, err := path.Match(pattern, url)
		if err != nil {
			return false
		}
		if match {
			return true
		}
	}
	return false
}
