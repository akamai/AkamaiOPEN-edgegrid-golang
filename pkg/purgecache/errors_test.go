package purgecache

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	t.Parallel()
	sess, err := session.New()
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodHead,
		"/",
		nil)
	require.NoError(t, err)

	tests := map[string]struct {
		response *http.Response
		expected *Error
	}{
		"Bad request 400": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`
				{
					"httpStatus": 400,
					"detail": "'objects' attribute cannot be missing or empty",
					"supportId": "aaaa-1234",
					"title": "invalid attribute value",
					"describedBy": "https://www.example.com/purge-cache/error-codes"
				}`)),
				Request: req,
			},
			expected: &Error{
				HTTPStatus:  http.StatusBadRequest,
				Detail:      "'objects' attribute cannot be missing or empty",
				SupportID:   "aaaa-1234",
				Title:       "invalid attribute value",
				DescribedBy: "https://www.example.com/purge-cache/error-codes",
			},
		},
		"Forbidden 403 unauthorized arl": {
			response: &http.Response{
				StatusCode: http.StatusForbidden,
				Body: io.NopCloser(strings.NewReader(`
				{
					"httpStatus": 403,
					"detail": "https://www.example.com/some/path",
					"supportId": "aaaa-1234",
					"title": "unauthorized arl",
					"describedBy": "https://www.example.com/purge-cache/error-codes"
				}`)),
				Request: req,
			},
			expected: &Error{
				HTTPStatus:  http.StatusForbidden,
				Detail:      "https://www.example.com/some/path",
				SupportID:   "aaaa-1234",
				Title:       "unauthorized arl",
				DescribedBy: "https://www.example.com/purge-cache/error-codes",
			},
		},
		"Not found 404": {
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(strings.NewReader(`
				{
					"httpStatus": 404,
					"detail": "Requested API endpoint does not exist. Check URI path for typos.",
					"supportId": "aaaa-1234",
					"title": "not found",
					"describedBy": "https://www.example.com/purge-cache/error-codes"
				}`)),
				Request: req,
			},
			expected: &Error{
				HTTPStatus:  http.StatusNotFound,
				Detail:      "Requested API endpoint does not exist. Check URI path for typos.",
				SupportID:   "aaaa-1234",
				Title:       "not found",
				DescribedBy: "https://www.example.com/purge-cache/error-codes",
			},
		},
		"Method not allowed 405": {
			response: &http.Response{
				StatusCode: http.StatusMethodNotAllowed,
				Body: io.NopCloser(strings.NewReader(`
				{
					"httpStatus": 405,
					"detail": "Specified HTTP method is invalid for this resource.",
					"supportId": "aaaa-1234",
					"title": "method not allowed",
					"describedBy": "https://www.example.com/purge-cache/error-codes"
				}`)),
				Request: req,
			},
			expected: &Error{
				HTTPStatus:  http.StatusMethodNotAllowed,
				Detail:      "Specified HTTP method is invalid for this resource.",
				SupportID:   "aaaa-1234",
				Title:       "method not allowed",
				DescribedBy: "https://www.example.com/purge-cache/error-codes",
			},
		},
		"Too many requests 429": {
			response: &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body: io.NopCloser(strings.NewReader(`
				{
					"describedBy": "https://www.example.com/purge-cache/rate-limit-errors",
					"detail": "",
					"httpStatus": 429,
					"rateLimit": 10000,
					"rateLimitCurrentRequestSize": 110,
					"rateLimitRemaining": 94,
					"supportId": "aaaa-1234",
					"title": "URL Rate Limit exceeded"
				}`)),
				Request: req,
			},
			expected: &Error{
				HTTPStatus:                  http.StatusTooManyRequests,
				Detail:                      "",
				SupportID:                   "aaaa-1234",
				Title:                       "URL Rate Limit exceeded",
				DescribedBy:                 "https://www.example.com/purge-cache/rate-limit-errors",
				RateLimit:                   10000,
				RateLimitCurrentRequestSize: 110,
				RateLimitRemaining:          94,
			},
		},
		"Gateway timeout 504": {
			response: &http.Response{
				StatusCode: http.StatusGatewayTimeout,
				Body: io.NopCloser(strings.NewReader(`
				{
					"httpStatus": 504,
					"detail": "Timed out while completing request. Please try again in a few seconds.",
					"supportId": "aaaa-1234",
					"title": "gateway timeout",
					"describedBy": "https://www.example.com/purge-cache/error-codes"
				}`)),
				Request: req,
			},
			expected: &Error{
				HTTPStatus:  http.StatusGatewayTimeout,
				Detail:      "Timed out while completing request. Please try again in a few seconds.",
				SupportID:   "aaaa-1234",
				Title:       "gateway timeout",
				DescribedBy: "https://www.example.com/purge-cache/error-codes",
			},
		},
		"Invalid response body": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				HTTPStatus: http.StatusInternalServerError,
				Title:      "Failed to unmarshal error body. Purge Cache API failed. Check details for more information.",
				Detail:     "test",
			},
		},
		"Empty response body": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			},
			expected: &Error{
				HTTPStatus: http.StatusInternalServerError,
				Title:      "Failed to unmarshal error body. Purge Cache API failed. Check details for more information.",
				Detail:     "",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*purgecache).Error(tc.response)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestIs(t *testing.T) {
	t.Parallel()
	someError := Error{
		HTTPStatus:  http.StatusNotFound,
		Title:       "not found",
		Detail:      "Requested API endpoint does not exist. Check URI path for typos.",
		DescribedBy: "https://www.example.com/purge-cache/error-codes",
		SupportID:   "aaaa-1234",
	}

	tests := map[string]struct {
		source   *Error
		target   Error
		expected bool
	}{
		"same HTTP status": {
			target:   Error{HTTPStatus: 404},
			expected: true,
		},
		"different HTTP status": {
			target:   Error{HTTPStatus: 400},
			expected: false,
		},
		"same title": {
			target:   Error{Title: "not found"},
			expected: true,
		},
		"different title": {
			target:   Error{Title: "gateway timeout"},
			expected: false,
		},
		"same detail": {
			target:   Error{Detail: "Requested API endpoint does not exist. Check URI path for typos."},
			expected: true,
		},
		"different detail": {
			target:   Error{Detail: "some other detail"},
			expected: false,
		},
		"same describedBy": {
			target:   Error{DescribedBy: "https://www.example.com/purge-cache/error-codes"},
			expected: true,
		},
		"different describedBy": {
			target:   Error{DescribedBy: "https://www.example.com/purge-cache/rate-limit-errors"},
			expected: false,
		},
		"same HTTP status and title": {
			target:   Error{HTTPStatus: 404, Title: "not found"},
			expected: true,
		},
		"same HTTP status but different title": {
			target:   Error{HTTPStatus: 404, Title: "gateway timeout"},
			expected: false,
		},
		"same title but different HTTP status": {
			target:   Error{Title: "not found", HTTPStatus: 400},
			expected: false,
		},
		"all fields zero values match any": {
			target:   Error{},
			expected: true,
		},
		"all matching fields set": {
			target: Error{
				HTTPStatus:  404,
				Title:       "not found",
				Detail:      "Requested API endpoint does not exist. Check URI path for typos.",
				DescribedBy: "https://www.example.com/purge-cache/error-codes",
			},
			expected: true,
		},
		"supportId not used in comparison - different supportId still matches": {
			target: Error{
				HTTPStatus: 404,
				Title:      "not found",
				SupportID:  "bbbb-5678",
			},
			expected: true,
		},
		"empty source matches empty target": {
			source:   &Error{},
			target:   Error{},
			expected: true,
		},
		"rate limit error matches by status and title": {
			source: &Error{
				HTTPStatus:                  429,
				Title:                       "URL Rate Limit exceeded",
				DescribedBy:                 "https://www.example.com/purge-cache/rate-limit-errors",
				RateLimit:                   10000,
				RateLimitCurrentRequestSize: 110,
				RateLimitRemaining:          94,
			},
			target:   Error{HTTPStatus: 429, Title: "URL Rate Limit exceeded"},
			expected: true,
		},
		"rate limit error does not match different title": {
			source: &Error{
				HTTPStatus: 429,
				Title:      "URL Rate Limit exceeded",
			},
			target:   Error{HTTPStatus: 429, Title: "not found"},
			expected: false,
		},
		"zero value HTTP status ignored": {
			target:   Error{Title: "not found", HTTPStatus: 0},
			expected: true,
		},
		"empty detail ignored": {
			target:   Error{HTTPStatus: 404, Detail: ""},
			expected: true,
		},
		"empty describedBy ignored": {
			target:   Error{HTTPStatus: 404, DescribedBy: ""},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			source := test.source
			if source == nil {
				source = &someError
			}
			assert.Equal(t, test.expected, source.Is(&test.target))
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	e := &Error{
		HTTPStatus:  http.StatusBadRequest,
		Title:       "invalid attribute value",
		Detail:      "'objects' attribute cannot be missing or empty",
		DescribedBy: "https://www.example.com/purge-cache/error-codes",
		SupportID:   "aaaa-1234",
	}
	expected := `API error: 
{
	"httpStatus": 400,
	"title": "invalid attribute value",
	"detail": "'objects' attribute cannot be missing or empty",
	"describedBy": "https://www.example.com/purge-cache/error-codes",
	"supportId": "aaaa-1234"
}`
	assert.EqualError(t, e, expected)
}

func TestError_RateLimit(t *testing.T) {
	t.Parallel()
	e := &Error{
		HTTPStatus:                  http.StatusTooManyRequests,
		Title:                       "URL Rate Limit exceeded",
		Detail:                      "",
		DescribedBy:                 "https://www.example.com/purge-cache/rate-limit-errors",
		SupportID:                   "aaaa-1234",
		RateLimit:                   10000,
		RateLimitCurrentRequestSize: 110,
		RateLimitRemaining:          94,
	}
	expected := `API error: 
{
	"httpStatus": 429,
	"title": "URL Rate Limit exceeded",
	"detail": "",
	"describedBy": "https://www.example.com/purge-cache/rate-limit-errors",
	"supportId": "aaaa-1234",
	"rateLimit": 10000,
	"rateLimitCurrentRequestSize": 110,
	"rateLimitRemaining": 94
}`
	assert.EqualError(t, e, expected)
}
