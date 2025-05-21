package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListApiHostnameCoverageOverlapping(t *testing.T) {

	result := GetApiHostnameCoverageOverlappingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiHostnameCoverageOverlapping/ApiHostnameCoverageOverlapping.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetApiHostnameCoverageOverlappingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetApiHostnameCoverageOverlappingResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetApiHostnameCoverageOverlappingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/hostname-coverage/overlapping?hostname=",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetApiHostnameCoverageOverlappingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching ApiHostnameCoverageOverlapping",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/hostname-coverage/overlapping?hostname=",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ApiHostnameCoverageOverlapping",
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
			result, err := client.GetApiHostnameCoverageOverlapping(
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

// Test ApiHostnameCoverageOverlapping
func TestAppSec_GetApiHostnameCoverageOverlapping(t *testing.T) {

	result := GetApiHostnameCoverageOverlappingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiHostnameCoverageOverlapping/ApiHostnameCoverageOverlapping.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetApiHostnameCoverageOverlappingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetApiHostnameCoverageOverlappingResponse
		withError        error
	}{
		"200 OK": {
			params: GetApiHostnameCoverageOverlappingRequest{
				ConfigID: 43253,
				Version:  15,
				Hostname: "www.example.com",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/hostname-coverage/overlapping?hostname=www.example.com",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetApiHostnameCoverageOverlappingRequest{
				ConfigID: 43253,
				Version:  15,
				Hostname: "www.example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching ApiHostnameCoverageOverlapping"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/hostname-coverage/overlapping?hostname=www.example.com",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ApiHostnameCoverageOverlapping",
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
			result, err := client.GetApiHostnameCoverageOverlapping(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
