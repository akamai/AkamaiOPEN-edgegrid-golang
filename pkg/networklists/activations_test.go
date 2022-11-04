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

func TestApsec_ListActivations(t *testing.T) {
	result := GetActivationsResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestActivations/Activations.json"))
	json.Unmarshal([]byte(respData), &result)
	tests := map[string]struct {
		params           GetActivationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetActivationsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetActivationsRequest{UniqueID: "38069_INTERNALWHITELIST", Network: "STAGING"},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/network-lists/38069_INTERNALWHITELIST/environments/STAGING/status",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetActivationsRequest{UniqueID: "38069_INTERNALWHITELIST", Network: "STAGING"},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching Activations",
    "status": 500
}`,
			expectedPath: "/network-list/v2/network-lists/38069_INTERNALWHITELIST/environments/STAGING/status",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching Activations",
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
			result, err := client.GetActivations(
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

// Test Activations
func TestAppSec_GetActivations(t *testing.T) {
	result := GetActivationsResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestActivations/Activations.json"))
	json.Unmarshal([]byte(respData), &result)
	tests := map[string]struct {
		params           GetActivationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetActivationsResponse
		withError        error
	}{
		"200 OK": {
			params:           GetActivationsRequest{UniqueID: "38069_INTERNALWHITELIST", Network: "STAGING"},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/network-lists/38069_INTERNALWHITELIST/environments/STAGING/status",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetActivationsRequest{UniqueID: "38069_INTERNALWHITELIST", Network: "STAGING"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching Activations"
}`,
			expectedPath: "/network-list/v2/network-lists/38069_INTERNALWHITELIST/environments/STAGING/status",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching Activations",
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
			result, err := client.GetActivations(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
