package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/log"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRequest(t *testing.T, method, url string) *http.Request {
	r, err := http.NewRequest(method, url, nil)
	assert.NoError(t, err)
	return r
}

func TestOverrideRetryPolicy(t *testing.T) {
	basePolicy := func(_ context.Context, _ *http.Response, _ error) (bool, error) {
		return false, errors.New("base policy: dummy, not implemented")
	}
	policy := overrideRetryPolicy(basePolicy, []string{"/excluded"})

	tests := map[string]struct {
		ctx            context.Context
		resp           *http.Response
		err            error
		expectedResult bool
		expectedError  string
	}{
		"should retry for PAPI GET with status 429": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodGet, "/papi/v1/sth"),
				StatusCode: http.StatusTooManyRequests,
			},
			expectedResult: true,
		},
		"should retry for GET with status 409 conflict": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    &http.Request{Method: http.MethodGet},
				StatusCode: http.StatusConflict,
			},
			expectedResult: true,
		},
		"should call base policy for other GETs": {
			ctx:           context.Background(),
			resp:          &http.Response{Request: &http.Request{Method: http.MethodGet}},
			expectedError: "base policy: dummy, not implemented",
		},
		"should forward context error when present": {
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			resp:          &http.Response{Request: &http.Request{Method: http.MethodGet}},
			expectedError: "context canceled",
		},
		"should not retry for POST": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: &http.Request{Method: http.MethodPost}},
			expectedResult: false,
		},
		"should not retry for PUT": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: &http.Request{Method: http.MethodPut}},
			expectedResult: false,
		},
		"should not retry for PATCH": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: &http.Request{Method: http.MethodPatch}},
			expectedResult: false,
		},
		"should not retry for HEAD": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: &http.Request{Method: http.MethodHead}},
			expectedResult: false,
		},
		"should not retry for DELETE": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: &http.Request{Method: http.MethodDelete}},
			expectedResult: false,
		},
		"should not retry excluded endpoints": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/excluded"}}},
			expectedResult: false,
		},
		"nil request with error": {
			ctx:            context.Background(),
			resp:           nil,
			err:            errors.New("test error"),
			expectedResult: false,
			expectedError:  "test error",
		},
		"should call base policy for GET url.Error": {
			ctx:  context.Background(),
			resp: nil,
			err: &url.Error{
				Op:  http.MethodGet,
				URL: "",
				Err: nil,
			},
			expectedError: "base policy: dummy, not implemented",
		},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			shouldRetry, err := policy(tst.ctx, tst.resp, tst.err)
			if len(tst.expectedError) > 0 {
				assert.ErrorContains(t, err, tst.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tst.expectedResult, shouldRetry)
			}
		})
	}
}

func stat429ResponseWaiting(wait time.Duration) *http.Response {
	res := http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{},
	}

	now := time.Now().UTC().Round(time.Second)
	date := strings.Replace(now.Format(time.RFC1123), "UTC", "GMT", 1)
	res.Header.Add("Date", date)
	if wait != 0 {
		// Add: allow to canonicalize to X-Ratelimit-Next or the header won't be recognized
		res.Header.Add("X-RateLimit-Next", now.Add(wait).Format(time.RFC3339Nano))
	}
	return &res
}

func Test_overrideBackoff(t *testing.T) {
	baseWait := time.Duration(24) * time.Hour
	baseBackoff := func(_, _ time.Duration, _ int, _ *http.Response) time.Duration {
		return baseWait
	}
	backoff := overrideBackoff(baseBackoff, nil)

	tests := map[string]struct {
		resp           *http.Response
		expectedResult time.Duration
	}{
		"correctly calculates backoff from X-RateLimit-Next": {
			resp:           stat429ResponseWaiting(time.Duration(5729) * time.Millisecond),
			expectedResult: time.Duration(5729) * time.Millisecond,
		},
		"falls back for next in the past": {
			resp:           stat429ResponseWaiting(-time.Duration(5729) * time.Millisecond),
			expectedResult: baseWait,
		},
		"falls back for no X-RateLimit-Next header": {
			resp:           stat429ResponseWaiting(0),
			expectedResult: baseWait,
		},
		"falls back for invalid X-RateLimit-Next header": {
			resp: func() *http.Response {
				r := stat429ResponseWaiting(time.Duration(5729) * time.Millisecond)
				r.Header.Set("X-RateLimit-Next", "2024-07-01T14:32:28.645???")
				return r
			}(),
			expectedResult: baseWait,
		},
		"falls back for no Date header": {
			resp: func() *http.Response {
				r := stat429ResponseWaiting(time.Duration(5729) * time.Millisecond)
				r.Header.Del("Date")
				return r
			}(),
			expectedResult: baseWait,
		},
		"falls back for invalid Date header": {
			resp: func() *http.Response {
				r := stat429ResponseWaiting(time.Duration(5729) * time.Millisecond)
				r.Header.Set("Date", "Mon, 01 Jul 2024 99:99:99 GMT")
				return r
			}(),
			expectedResult: baseWait,
		},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			wait := backoff(1, 30, 1, tst.resp)
			assert.Equal(t, tst.expectedResult, wait)
		})
	}
}

func TestXRateLimitGet(t *testing.T) {
	xrlHandler := test.XRateLimitHTTPHandler{
		T:           t,
		SuccessCode: http.StatusOK,
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/papi/test", r.URL.String())
		assert.Equal(t, http.MethodGet, r.Method)
		xrlHandler.ServeHTTP(w, r)
	}))
	defer mockServer.Close()

	mockSession := mockRetrySession(t, mockServer)
	client := mockSession.Client()
	_, err := client.Get(fmt.Sprintf("%s/papi/test", mockServer.URL))
	assert.NoError(t, err)

	// We expect exactly two requests to the server:
	// - the first resulting in code 429
	// - the second after a proper backoff, resulting in status 200
	assert.Equal(t, []int{http.StatusTooManyRequests, http.StatusOK}, xrlHandler.ReturnedCodes())
	assert.Less(t,
		xrlHandler.ReturnTimes()[1],
		xrlHandler.AvailableAt().Add(time.Duration(time.Millisecond)*1100))
}

func mockRetrySession(t *testing.T, mockServer *httptest.Server) Session {

	serverURL, err := url.Parse(mockServer.URL)
	require.NoError(t, err)
	config := edgegrid.Config{Host: serverURL.Host}
	retryConf := NewRetryConfig()

	retrySession, err := New(WithRetries(retryConf), WithSigner(&config))
	assert.NoError(t, err)

	return retrySession
}

func Test_configureRetryClient(t *testing.T) {
	tests := map[string]struct {
		retryMax          int
		retryWaitMin      time.Duration
		retryWaitMax      time.Duration
		excludedEndpoints []string
		expected          *retryablehttp.Client
		expectedError     string
	}{
		"happy path": {
			retryMax:          5,
			retryWaitMin:      5 * time.Second,
			retryWaitMax:      31 * time.Second,
			excludedEndpoints: []string{"aaa/bbb/ccc", "aaa/*/ccc"},
			expected: &retryablehttp.Client{
				RetryWaitMin: 5 * time.Second,
				RetryWaitMax: 31 * time.Second,
				RetryMax:     5,
			},
		},
		"negative number of retries": {
			retryMax:      -5,
			retryWaitMin:  5 * time.Second,
			retryWaitMax:  30 * time.Second,
			expected:      nil,
			expectedError: `maximum number of retries cannot be negative`,
		},
		"negative wait times": {
			retryMax:     5,
			retryWaitMin: -5 * time.Second,
			retryWaitMax: -5 * time.Second,
			expected:     nil,
			expectedError: `minimum retry wait time cannot be negative
maximum retry wait time cannot be negative`,
		},
		"minimum wait time higher that maximum wait time": {
			retryMax:      5,
			retryWaitMin:  30 * time.Second,
			retryWaitMax:  5 * time.Second,
			expected:      nil,
			expectedError: `maximum retry wait time cannot be shorter than minimum retry wait time`,
		},
		"malformed excluded endpoint pattern": {
			excludedEndpoints: []string{"[-]"},
			expected:          nil,
			expectedError:     "malformed exclude endpoint pattern: syntax error in pattern: [-]",
		},
		"test error formation": {
			retryMax:          -5,
			retryWaitMin:      -3 * time.Second,
			retryWaitMax:      -7 * time.Second,
			excludedEndpoints: []string{"[-]"},
			expected:          nil,
			expectedError: `maximum number of retries cannot be negative
minimum retry wait time cannot be negative
maximum retry wait time cannot be negative
maximum retry wait time cannot be shorter than minimum retry wait time
malformed exclude endpoint pattern: syntax error in pattern: [-]`,
		},
	}

	sessionLogger := log.NOPLogger()

	testSession := &session{log: sessionLogger}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			conf := RetryConfig{
				RetryMax:          test.retryMax,
				RetryWaitMin:      test.retryWaitMin,
				RetryWaitMax:      test.retryWaitMax,
				ExcludedEndpoints: test.excludedEndpoints,
			}
			got, err := configureRetryClient(conf, testSession.Sign, testSession.log)

			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}
			if test.expected != nil {
				assert.Equal(t, test.expected.RetryMax, got.RetryMax)
				assert.Equal(t, test.expected.RetryWaitMax, got.RetryWaitMax)
				assert.Equal(t, test.expected.RetryWaitMin, got.RetryWaitMin)
			} else {
				assert.Nil(t, got)
			}

		})
	}
}

func Test_isBlocked(t *testing.T) {
	tests := map[string]struct {
		url      string
		blocked  []string
		expected bool
	}{
		"blocked - no wildcards": {
			url:      "/api/blocked",
			blocked:  []string{"/api/blocked"},
			expected: true,
		},
		"blocked - with wildcard": {
			url:      "/api/blocked/123",
			blocked:  []string{"/api/blocked/*"},
			expected: true,
		},
		"allowed - all": {
			url:      "/api/not-blocked/123/data",
			blocked:  []string{},
			expected: false,
		},
		"allowed - with wildcard": {
			url:      "/api/not-blocked/123/data",
			blocked:  []string{"/api/not-blocked/*"},
			expected: false,
		},
		"allowed - wildcard does not match beginning of url": {
			url:      "/api/not-blocked/123/data",
			blocked:  []string{"/not-blocked/*/data"},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := isBlocked(test.url, test.blocked)
			assert.Equal(t, test.expected, got)

		})
	}
}
