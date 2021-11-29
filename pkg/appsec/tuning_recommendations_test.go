package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_GetTuningRecommendations(t *testing.T) {

	result := GetTuningRecommendationsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestTuningRecommendations/Recommendations.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetTuningRecommendationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetTuningRecommendationsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetTuningRecommendationsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/recommendations?standardException=true",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetTuningRecommendationsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching propertys",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/recommendations?standardException=true",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching propertys",
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
			result, err := client.GetTuningRecommendations(
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

func TestAppSec_GetAttackGroupRecommendations(t *testing.T) {

	result := GetAttackGroupRecommendationsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestTuningRecommendations/AttackGroupRecommendations.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAttackGroupRecommendationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAttackGroupRecommendationsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAttackGroupRecommendationsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Group:    "XSS",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/recommendations/attack-groups/XSS?standardException=true",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAttackGroupRecommendationsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Group:    "XSS",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching propertys",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/recommendations/attack-groups/XSS?standardException=true",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching propertys",
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
			result, err := client.GetAttackGroupRecommendations(
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
