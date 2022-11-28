package appsec

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

func TestAppSec_ListApiHostnameCoverage(t *testing.T) {

	result := GetApiHostnameCoverageResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiHostnameCoverage/ApiHostnameCoverage.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetApiHostnameCoverageRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetApiHostnameCoverageResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetApiHostnameCoverageRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/hostname-coverage",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetApiHostnameCoverageRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching ApiHostnameCoverage",
    "status": 500
}`,
			expectedPath: "/appsec/v1/hostname-coverage",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ApiHostnameCoverage",
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
			result, err := client.GetApiHostnameCoverage(
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

// Test ApiHostnameCoverage
func TestAppSec_GetApiHostnameCoverage(t *testing.T) {

	result := GetApiHostnameCoverageResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiHostnameCoverage/ApiHostnameCoverage.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetApiHostnameCoverageRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetApiHostnameCoverageResponse
		withError        error
	}{
		"200 OK": {
			params: GetApiHostnameCoverageRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/hostname-coverage",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetApiHostnameCoverageRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching ApiHostnameCoverage"
}`),
			expectedPath: "/appsec/v1/hostname-coverage",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ApiHostnameCoverage",
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
			result, err := client.GetApiHostnameCoverage(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
