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

func TestAppSec_ListApiHostnameCoverageMatchTargets(t *testing.T) {

	result := GetApiHostnameCoverageMatchTargetsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiHostnameCoverageMatchTargets/ApiHostnameCoverageMatchTargets.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetApiHostnameCoverageMatchTargetsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetApiHostnameCoverageMatchTargetsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetApiHostnameCoverageMatchTargetsRequest{
				ConfigID: 43253,
				Version:  15,
				Hostname: "www.example.com",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/hostname-coverage/match-targets?hostname=www.example.com",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetApiHostnameCoverageMatchTargetsRequest{
				ConfigID: 43253,
				Version:  15,
				Hostname: "www.example.com",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching ApiHostnameCoverageMatchTargets",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/hostname-coverage/match-targets?hostname=www.example.com",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ApiHostnameCoverageMatchTargets",
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
			result, err := client.GetApiHostnameCoverageMatchTargets(
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

// Test ApiHostnameCoverageMatchTargets
func TestAppSec_GetApiHostnameCoverageMatchTargets(t *testing.T) {

	result := GetApiHostnameCoverageMatchTargetsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiHostnameCoverageMatchTargets/ApiHostnameCoverageMatchTargets.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetApiHostnameCoverageMatchTargetsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetApiHostnameCoverageMatchTargetsResponse
		withError        error
	}{
		"200 OK": {
			params: GetApiHostnameCoverageMatchTargetsRequest{
				ConfigID: 43253,
				Version:  15,
				Hostname: "www.example.com",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/hostname-coverage/match-targets?hostname=www.example.com",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetApiHostnameCoverageMatchTargetsRequest{
				ConfigID: 43253,
				Version:  15,
				Hostname: "www.example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching ApiHostnameCoverageMatchTargets"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/hostname-coverage/match-targets?hostname=www.example.com",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ApiHostnameCoverageMatchTargets",
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
			result, err := client.GetApiHostnameCoverageMatchTargets(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
