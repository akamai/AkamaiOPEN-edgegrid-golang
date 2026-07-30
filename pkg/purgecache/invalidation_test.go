package purgecache

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidateByURL(t *testing.T) {
	t.Parallel()

	result := InvalidateResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestInvalidate/Invalidate.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	tests := map[string]struct {
		params          InvalidateByURLRequest
		responseStatus  int
		responseBody    string
		responseHeaders map[string]string
		expectedPath    string
		expectedBody    string
		expectedResp    *InvalidateResponse
		withError       func(*testing.T, error)
	}{
		"201 - with rate limit headers": {
			params: InvalidateByURLRequest{
				Objects: []string{"http://example.com"},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			responseHeaders: map[string]string{
				"X-Ratelimit-Limit":                            "200",
				"X-Ratelimit-Limit-Objects":                    "2000",
				"X-Ratelimit-Limit-Per-Second":                 "10.5",
				"X-Ratelimit-Limit-Per-Second-Objects":         "100.0",
				"X-Ratelimit-Remaining":                        "195",
				"X-Ratelimit-Remaining-Objects":                "1950",
				"X-Ratelimit-Seconds-To-Refresh-Limit":         "0.5",
				"X-Ratelimit-Seconds-To-Refresh-Limit-Objects": "0.1",
			},
			expectedPath: "/ccu/v3/invalidate/url",
			expectedBody: `{"objects":["http://example.com"]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{
					Limit:                        ptr.To(int64(200)),
					LimitObjects:                 ptr.To(int64(2000)),
					LimitPerSecond:               ptr.To(10.5),
					LimitPerSecondObjects:        ptr.To(100.0),
					Remaining:                    ptr.To(int64(195)),
					RemainingObjects:             ptr.To(int64(1950)),
					SecondsToRefreshLimit:        ptr.To(0.5),
					SecondsToRefreshLimitObjects: ptr.To(0.1),
				},
			},
		},
		"201 - no rate limit headers": {
			params: InvalidateByURLRequest{
				Objects: []string{"http://example.com"},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/url",
			expectedBody:   `{"objects":["http://example.com"]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"201 - staging network": {
			params: InvalidateByURLRequest{
				Network: PurgeNetworkStaging,
				Objects: []string{"http://example.com"},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/url/staging",
			expectedBody:   `{"objects":["http://example.com"]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"201 - production network": {
			params: InvalidateByURLRequest{
				Network: PurgeNetworkProduction,
				Objects: []string{"http://example.com"},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/url/production",
			expectedBody:   `{"objects":["http://example.com"]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"400 bad request": {
			params: InvalidateByURLRequest{
				Objects: []string{"not-a-valid-url"},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"httpStatus": 400,
				"detail": "invalid request",
				"supportId": "aaaa-1234",
				"title": "Bad Request",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/400"
			}`,
			expectedPath: "/ccu/v3/invalidate/url",
			expectedBody: `{"objects":["not-a-valid-url"]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:  http.StatusBadRequest,
					Detail:      "invalid request",
					SupportID:   "aaaa-1234",
					Title:       "Bad Request",
					DescribedBy: "https://techdocs.akamai.com/purge-cache/reference/400",
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"429 too many requests": {
			params: InvalidateByURLRequest{
				Objects: []string{"http://example.com"},
			},
			responseStatus: http.StatusTooManyRequests,
			responseBody: `{
				"httpStatus": 429,
				"detail": "Rate Limit exceeded",
				"supportId": "bbbb-5678",
				"title": "Rate Limit exceeded",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/429",
				"rateLimit": 50,
				"rateLimitCurrentRequestSize": 1,
				"rateLimitRemaining": 0
			}`,
			expectedPath: "/ccu/v3/invalidate/url",
			expectedBody: `{"objects":["http://example.com"]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:                  http.StatusTooManyRequests,
					Detail:                      "Rate Limit exceeded",
					SupportID:                   "bbbb-5678",
					Title:                       "Rate Limit exceeded",
					DescribedBy:                 "https://techdocs.akamai.com/purge-cache/reference/429",
					RateLimit:                   ptr.To(int64(50)),
					RateLimitCurrentRequestSize: ptr.To(int64(1)),
					RateLimitRemaining:          ptr.To(int64(0)),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: InvalidateByURLRequest{
				Objects: []string{"http://example.com"},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"httpStatus": 500,
				"detail": "Internal Server Error",
				"supportId": "cccc-9012",
				"title": "Internal Server Error",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/500"
			}`,
			expectedPath: "/ccu/v3/invalidate/url",
			expectedBody: `{"objects":["http://example.com"]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:  http.StatusInternalServerError,
					Detail:      "Internal Server Error",
					SupportID:   "cccc-9012",
					Title:       "Internal Server Error",
					DescribedBy: "https://techdocs.akamai.com/purge-cache/reference/500",
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error - missing objects": {
			params: InvalidateByURLRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t,
					"invalidate cache: struct validation: Objects: must contain at least one URL or ARL to invalidate",
					err.Error())
				assert.ErrorIs(t, err, ErrInvalidate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid network": {
			params: InvalidateByURLRequest{
				Network: "invalid",
				Objects: []string{"http://example.com"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t,
					"invalidate cache: struct validation: Network: value 'invalid' is invalid. Must be one of: 'staging', 'production'",
					err.Error())
				assert.ErrorIs(t, err, ErrInvalidate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if test.responseStatus == 0 {
				sess, err := session.New()
				require.NoError(t, err)
				client := Client(sess)
				_, err = client.InvalidateByURL(context.Background(), test.params)
				require.Error(t, err)
				test.withError(t, err)
				return
			}

			mockServer := getMockTestServer(t, http.MethodPost, test.expectedPath, test.responseStatus, test.responseBody, test.expectedBody, test.responseHeaders)

			client := mockAPIClient(t, mockServer)
			resp, err := client.InvalidateByURL(context.Background(), test.params)
			if test.withError != nil {
				require.Error(t, err)
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResp, resp)
		})
	}
}

func TestInvalidateByTag(t *testing.T) {
	t.Parallel()

	result := InvalidateResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestInvalidate/Invalidate.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	tests := map[string]struct {
		params          InvalidateByTagRequest
		responseStatus  int
		responseBody    string
		responseHeaders map[string]string
		expectedPath    string
		expectedBody    string
		expectedResp    *InvalidateResponse
		withError       func(*testing.T, error)
	}{
		"201 - with rate limit headers": {
			params: InvalidateByTagRequest{
				Objects: []string{"tag-1", "tag-2"},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			responseHeaders: map[string]string{
				"X-Ratelimit-Limit":                            "200",
				"X-Ratelimit-Limit-Objects":                    "2000",
				"X-Ratelimit-Limit-Per-Second":                 "10.5",
				"X-Ratelimit-Limit-Per-Second-Objects":         "100.0",
				"X-Ratelimit-Remaining":                        "195",
				"X-Ratelimit-Remaining-Objects":                "1950",
				"X-Ratelimit-Seconds-To-Refresh-Limit":         "0.5",
				"X-Ratelimit-Seconds-To-Refresh-Limit-Objects": "0.1",
			},
			expectedPath: "/ccu/v3/invalidate/tag",
			expectedBody: `{"objects":["tag-1","tag-2"]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{
					Limit:                        ptr.To(int64(200)),
					LimitObjects:                 ptr.To(int64(2000)),
					LimitPerSecond:               ptr.To(10.5),
					LimitPerSecondObjects:        ptr.To(100.0),
					Remaining:                    ptr.To(int64(195)),
					RemainingObjects:             ptr.To(int64(1950)),
					SecondsToRefreshLimit:        ptr.To(0.5),
					SecondsToRefreshLimitObjects: ptr.To(0.1),
				},
			},
		},
		"201 - no rate limit headers": {
			params: InvalidateByTagRequest{
				Objects: []string{"tag-1"},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/tag",
			expectedBody:   `{"objects":["tag-1"]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"201 - staging network": {
			params: InvalidateByTagRequest{
				Network: PurgeNetworkStaging,
				Objects: []string{"tag-1"},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/tag/staging",
			expectedBody:   `{"objects":["tag-1"]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"201 - production network": {
			params: InvalidateByTagRequest{
				Network: PurgeNetworkProduction,
				Objects: []string{"tag-1"},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/tag/production",
			expectedBody:   `{"objects":["tag-1"]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"400 bad request": {
			params: InvalidateByTagRequest{
				Objects: []string{"tag-1"},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"httpStatus": 400,
				"detail": "invalid request",
				"supportId": "aaaa-1234",
				"title": "Bad Request",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/400"
			}`,
			expectedPath: "/ccu/v3/invalidate/tag",
			expectedBody: `{"objects":["tag-1"]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:  http.StatusBadRequest,
					Detail:      "invalid request",
					SupportID:   "aaaa-1234",
					Title:       "Bad Request",
					DescribedBy: "https://techdocs.akamai.com/purge-cache/reference/400",
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"429 too many requests": {
			params: InvalidateByTagRequest{
				Objects: []string{"tag-1"},
			},
			responseStatus: http.StatusTooManyRequests,
			responseBody: `{
				"httpStatus": 429,
				"detail": "Rate Limit exceeded",
				"supportId": "bbbb-5678",
				"title": "Rate Limit exceeded",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/429",
				"rateLimit": 50,
				"rateLimitCurrentRequestSize": 1,
				"rateLimitRemaining": 0
			}`,
			expectedPath: "/ccu/v3/invalidate/tag",
			expectedBody: `{"objects":["tag-1"]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:                  http.StatusTooManyRequests,
					Detail:                      "Rate Limit exceeded",
					SupportID:                   "bbbb-5678",
					Title:                       "Rate Limit exceeded",
					DescribedBy:                 "https://techdocs.akamai.com/purge-cache/reference/429",
					RateLimit:                   ptr.To(int64(50)),
					RateLimitCurrentRequestSize: ptr.To(int64(1)),
					RateLimitRemaining:          ptr.To(int64(0)),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: InvalidateByTagRequest{
				Objects: []string{"tag-1"},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"httpStatus": 500,
				"detail": "Internal Server Error",
				"supportId": "cccc-9012",
				"title": "Internal Server Error",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/500"
			}`,
			expectedPath: "/ccu/v3/invalidate/tag",
			expectedBody: `{"objects":["tag-1"]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:  http.StatusInternalServerError,
					Detail:      "Internal Server Error",
					SupportID:   "cccc-9012",
					Title:       "Internal Server Error",
					DescribedBy: "https://techdocs.akamai.com/purge-cache/reference/500",
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error - missing objects": {
			params: InvalidateByTagRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t,
					"invalidate cache: struct validation: Objects: must contain at least one tag to invalidate",
					err.Error())
				assert.ErrorIs(t, err, ErrInvalidate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid network": {
			params: InvalidateByTagRequest{
				Network: "invalid",
				Objects: []string{"tag-1"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t,
					"invalidate cache: struct validation: Network: value 'invalid' is invalid. Must be one of: 'staging', 'production'",
					err.Error())
				assert.ErrorIs(t, err, ErrInvalidate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if test.responseStatus == 0 {
				sess, err := session.New()
				require.NoError(t, err)
				client := Client(sess)
				_, err = client.InvalidateByTag(context.Background(), test.params)
				require.Error(t, err)
				test.withError(t, err)
				return
			}

			mockServer := getMockTestServer(t, http.MethodPost, test.expectedPath, test.responseStatus, test.responseBody, test.expectedBody, test.responseHeaders)

			client := mockAPIClient(t, mockServer)
			resp, err := client.InvalidateByTag(context.Background(), test.params)
			if test.withError != nil {
				require.Error(t, err)
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResp, resp)
		})
	}
}

func TestInvalidateByCPCode(t *testing.T) {
	t.Parallel()

	result := InvalidateResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestInvalidate/Invalidate.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	tests := map[string]struct {
		params          InvalidateByCPCodeRequest
		responseStatus  int
		responseBody    string
		responseHeaders map[string]string
		expectedPath    string
		expectedBody    string
		expectedResp    *InvalidateResponse
		withError       func(*testing.T, error)
	}{
		"201 - with rate limit headers": {
			params: InvalidateByCPCodeRequest{
				Objects: []int64{12345, 67890},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			responseHeaders: map[string]string{
				"X-Ratelimit-Limit":                            "200",
				"X-Ratelimit-Limit-Objects":                    "2000",
				"X-Ratelimit-Limit-Per-Second":                 "10.5",
				"X-Ratelimit-Limit-Per-Second-Objects":         "100.0",
				"X-Ratelimit-Remaining":                        "195",
				"X-Ratelimit-Remaining-Objects":                "1950",
				"X-Ratelimit-Seconds-To-Refresh-Limit":         "0.5",
				"X-Ratelimit-Seconds-To-Refresh-Limit-Objects": "0.1",
			},
			expectedPath: "/ccu/v3/invalidate/cpcode",
			expectedBody: `{"objects":[12345,67890]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{
					Limit:                        ptr.To(int64(200)),
					LimitObjects:                 ptr.To(int64(2000)),
					LimitPerSecond:               ptr.To(10.5),
					LimitPerSecondObjects:        ptr.To(100.0),
					Remaining:                    ptr.To(int64(195)),
					RemainingObjects:             ptr.To(int64(1950)),
					SecondsToRefreshLimit:        ptr.To(0.5),
					SecondsToRefreshLimitObjects: ptr.To(0.1),
				},
			},
		},
		"201 - no rate limit headers": {
			params: InvalidateByCPCodeRequest{
				Objects: []int64{12345},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/cpcode",
			expectedBody:   `{"objects":[12345]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"201 - staging network": {
			params: InvalidateByCPCodeRequest{
				Network: PurgeNetworkStaging,
				Objects: []int64{12345},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/cpcode/staging",
			expectedBody:   `{"objects":[12345]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"201 - production network": {
			params: InvalidateByCPCodeRequest{
				Network: PurgeNetworkProduction,
				Objects: []int64{12345},
			},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/invalidate/cpcode/production",
			expectedBody:   `{"objects":[12345]}`,
			expectedResp: &InvalidateResponse{
				Detail:           result.Detail,
				EstimatedSeconds: result.EstimatedSeconds,
				HTTPStatus:       result.HTTPStatus,
				PurgeID:          result.PurgeID,
				SupportID:        result.SupportID,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"400 bad request": {
			params: InvalidateByCPCodeRequest{
				Objects: []int64{99999},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"httpStatus": 400,
				"detail": "invalid request",
				"supportId": "aaaa-1234",
				"title": "Bad Request",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/400"
			}`,
			expectedPath: "/ccu/v3/invalidate/cpcode",
			expectedBody: `{"objects":[99999]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:  http.StatusBadRequest,
					Detail:      "invalid request",
					SupportID:   "aaaa-1234",
					Title:       "Bad Request",
					DescribedBy: "https://techdocs.akamai.com/purge-cache/reference/400",
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"429 too many requests": {
			params: InvalidateByCPCodeRequest{
				Objects: []int64{12345},
			},
			responseStatus: http.StatusTooManyRequests,
			responseBody: `{
				"httpStatus": 429,
				"detail": "Rate Limit exceeded",
				"supportId": "bbbb-5678",
				"title": "Rate Limit exceeded",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/429",
				"rateLimit": 50,
				"rateLimitCurrentRequestSize": 1,
				"rateLimitRemaining": 0
			}`,
			expectedPath: "/ccu/v3/invalidate/cpcode",
			expectedBody: `{"objects":[12345]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:                  http.StatusTooManyRequests,
					Detail:                      "Rate Limit exceeded",
					SupportID:                   "bbbb-5678",
					Title:                       "Rate Limit exceeded",
					DescribedBy:                 "https://techdocs.akamai.com/purge-cache/reference/429",
					RateLimit:                   ptr.To(int64(50)),
					RateLimitCurrentRequestSize: ptr.To(int64(1)),
					RateLimitRemaining:          ptr.To(int64(0)),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: InvalidateByCPCodeRequest{
				Objects: []int64{12345},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"httpStatus": 500,
				"detail": "Internal Server Error",
				"supportId": "cccc-9012",
				"title": "Internal Server Error",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/500"
			}`,
			expectedPath: "/ccu/v3/invalidate/cpcode",
			expectedBody: `{"objects":[12345]}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					HTTPStatus:  http.StatusInternalServerError,
					Detail:      "Internal Server Error",
					SupportID:   "cccc-9012",
					Title:       "Internal Server Error",
					DescribedBy: "https://techdocs.akamai.com/purge-cache/reference/500",
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error - missing objects": {
			params: InvalidateByCPCodeRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t,
					"invalidate cache: struct validation: Objects: must contain at least one CP Code to invalidate",
					err.Error())
				assert.ErrorIs(t, err, ErrInvalidate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid network": {
			params: InvalidateByCPCodeRequest{
				Network: "invalid",
				Objects: []int64{12345},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t,
					"invalidate cache: struct validation: Network: value 'invalid' is invalid. Must be one of: 'staging', 'production'",
					err.Error())
				assert.ErrorIs(t, err, ErrInvalidate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if test.responseStatus == 0 {
				sess, err := session.New()
				require.NoError(t, err)
				client := Client(sess)
				_, err = client.InvalidateByCPCode(context.Background(), test.params)
				require.Error(t, err)
				test.withError(t, err)
				return
			}

			mockServer := getMockTestServer(t, http.MethodPost, test.expectedPath, test.responseStatus, test.responseBody, test.expectedBody, test.responseHeaders)

			client := mockAPIClient(t, mockServer)
			resp, err := client.InvalidateByCPCode(context.Background(), test.params)
			if test.withError != nil {
				require.Error(t, err)
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResp, resp)
		})
	}
}
