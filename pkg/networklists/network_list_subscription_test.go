package networklists

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApsec_ListNetworkListSubscription(t *testing.T) {

	result := GetNetworkListSubscriptionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkListSubscription/NetworkListSubscription.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetNetworkListSubscriptionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetNetworkListSubscriptionResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetNetworkListSubscriptionRequest{},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/notifications/subscriptions",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetNetworkListSubscriptionRequest{},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching subscriptions",
    "status": 500
}`,
			expectedPath: "/network-list/v2/notifications/subscriptions",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching subscriptions",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetNetworkListSubscription(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test NetworkListSubscription
func TestAppSec_GetNetworkListSubscription(t *testing.T) {

	result := GetNetworkListSubscriptionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkListSubscription/NetworkListSubscription.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetNetworkListSubscriptionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetNetworkListSubscriptionResponse
		withError        error
	}{
		"200 OK": {
			params:           GetNetworkListSubscriptionRequest{},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/notifications/subscriptions",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetNetworkListSubscriptionRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching subscriptions"
}`,
			expectedPath: "/network-list/v2/notifications/subscriptions",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching subscriptions",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetNetworkListSubscription(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update NetworkListSubscription.
func TestAppSec_UpdateNetworkListSubscription(t *testing.T) {
	result := UpdateNetworkListSubscriptionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkListSubscription/NetworkListSubscription.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateNetworkListSubscriptionRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestNetworkListSubscription/NetworkListSubscription.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateNetworkListSubscriptionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateNetworkListSubscriptionResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateNetworkListSubscriptionRequest{},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/network-list/v2/notifications/subscribe",
		},
		"500 internal server error": {
			params:         UpdateNetworkListSubscriptionRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating subscription"
}`,
			expectedPath: "/network-list/v2/notifications/subscribe",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating subscription",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, test.expectedPath, r.URL.String())
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateNetworkListSubscription(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
