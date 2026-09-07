package purgecache

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimitStatus(t *testing.T) {
	t.Parallel()

	result := RateLimitStatusResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestRateLimitStatus/RateLimitStatus.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	tests := map[string]struct {
		params          RateLimitStatusRequest
		responseStatus  int
		responseBody    string
		responseHeaders map[string]string
		expectedPath    string
		expectedResp    *RateLimitStatusResponse
		withError       func(*testing.T, error)
	}{
		"201 cpcode - with rate limit headers": {
			params:         RateLimitStatusRequest{PurgeType: PurgeTypeCPCode},
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
			expectedPath: "/ccu/v3/rate-limit-status/cpcode",
			expectedResp: &RateLimitStatusResponse{
				DescribedBy: result.DescribedBy,
				Detail:      result.Detail,
				HTTPStatus:  result.HTTPStatus,
				SupportID:   result.SupportID,
				Title:       result.Title,
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
		"201 tag - no rate limit headers": {
			params:         RateLimitStatusRequest{PurgeType: PurgeTypeTag},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/rate-limit-status/tag",
			expectedResp: &RateLimitStatusResponse{
				DescribedBy:      result.DescribedBy,
				Detail:           result.Detail,
				HTTPStatus:       result.HTTPStatus,
				SupportID:        result.SupportID,
				Title:            result.Title,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"201 url": {
			params:         RateLimitStatusRequest{PurgeType: PurgeTypeURL},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/rate-limit-status/url",
			expectedResp: &RateLimitStatusResponse{
				DescribedBy:      result.DescribedBy,
				Detail:           result.Detail,
				HTTPStatus:       result.HTTPStatus,
				SupportID:        result.SupportID,
				Title:            result.Title,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"400 bad request": {
			params:         RateLimitStatusRequest{PurgeType: PurgeTypeCPCode},
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"httpStatus": 400,
				"detail": "invalid request",
				"supportId": "aaaa-1234",
				"title": "Bad Request",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/400"
			}`,
			expectedPath: "/ccu/v3/rate-limit-status/cpcode",
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
			params:         RateLimitStatusRequest{PurgeType: PurgeTypeCPCode},
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
			expectedPath: "/ccu/v3/rate-limit-status/cpcode",
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
			params:         RateLimitStatusRequest{PurgeType: PurgeTypeCPCode},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"httpStatus": 500,
				"detail": "Internal Server Error",
				"supportId": "cccc-9012",
				"title": "Internal Server Error",
				"describedBy": "https://techdocs.akamai.com/purge-cache/reference/500"
			}`,
			expectedPath: "/ccu/v3/rate-limit-status/cpcode",
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
		"validation error - missing purge type": {
			params: RateLimitStatusRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "checking rate limit status: struct validation: PurgeType: cannot be blank",
					err.Error())
				assert.ErrorIs(t, err, ErrRateLimitStatus)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid purge type": {
			params: RateLimitStatusRequest{PurgeType: "invalid"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t,
					"checking rate limit status: struct validation: PurgeType: value 'invalid' is invalid. Must be one of: 'cpcode', 'tag', 'url'",
					err.Error())
				assert.ErrorIs(t, err, ErrRateLimitStatus)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"201 cpcode - staging network": {
			params:         RateLimitStatusRequest{PurgeType: PurgeTypeCPCode, Network: PurgeNetworkStaging},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/rate-limit-status/cpcode/staging",
			expectedResp: &RateLimitStatusResponse{
				DescribedBy:      result.DescribedBy,
				Detail:           result.Detail,
				HTTPStatus:       result.HTTPStatus,
				SupportID:        result.SupportID,
				Title:            result.Title,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"201 cpcode - production network": {
			params:         RateLimitStatusRequest{PurgeType: PurgeTypeCPCode, Network: PurgeNetworkProduction},
			responseStatus: http.StatusCreated,
			responseBody:   respData,
			expectedPath:   "/ccu/v3/rate-limit-status/cpcode/production",
			expectedResp: &RateLimitStatusResponse{
				DescribedBy:      result.DescribedBy,
				Detail:           result.Detail,
				HTTPStatus:       result.HTTPStatus,
				SupportID:        result.SupportID,
				Title:            result.Title,
				RateLimitHeaders: RateLimitHeaders{},
			},
		},
		"validation error - invalid network": {
			params: RateLimitStatusRequest{PurgeType: PurgeTypeCPCode, Network: "invalid"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t,
					"checking rate limit status: struct validation: Network: value 'invalid' is invalid. Must be one of: 'staging', 'production'",
					err.Error())
				assert.ErrorIs(t, err, ErrRateLimitStatus)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if test.responseStatus == 0 {
				// Validation error — no HTTP server needed.
				sess, err := session.New()
				require.NoError(t, err)
				client := Client(sess)
				_, err = client.RateLimitStatus(context.Background(), test.params)
				require.Error(t, err)
				test.withError(t, err)
				return
			}

			mockServer := getMockTestServer(t, http.MethodPost, test.expectedPath, test.responseStatus, test.responseBody, "", test.responseHeaders)

			client := mockAPIClient(t, mockServer)
			resp, err := client.RateLimitStatus(context.Background(), test.params)
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
