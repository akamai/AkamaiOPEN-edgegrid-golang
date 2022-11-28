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

func TestApsec_ListAdvancedSettingsEvasivePathMatch(t *testing.T) {

	result := GetAdvancedSettingsEvasivePathMatchResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsEvasivePathMatch/AdvancedSettingsEvasivePathMatch.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAdvancedSettingsEvasivePathMatchRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsEvasivePathMatchResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAdvancedSettingsEvasivePathMatchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/evasive-path-match",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsEvasivePathMatchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching AdvancedSettingsEvasivePathMatch",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/evasive-path-match",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsEvasivePathMatch",
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
			result, err := client.GetAdvancedSettingsEvasivePathMatch(
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

// Test AdvancedSettingsEvasivePathMatch
func TestAppSec_GetAdvancedSettingsEvasivePathmatch(t *testing.T) {

	result := GetAdvancedSettingsEvasivePathMatchResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsEvasivePathMatch/AdvancedSettingsEvasivePathMatch.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAdvancedSettingsEvasivePathMatchRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsEvasivePathMatchResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsEvasivePathMatchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/evasive-path-match",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsEvasivePathMatchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching AdvancedSettingsEvasivePathMatch"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/evasive-path-match",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsEvasivePathMatch",
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
			result, err := client.GetAdvancedSettingsEvasivePathMatch(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update AdvancedSettingsEvasivePathMatch.
func TestAppSec_UpdateAdvancedSettingsEvasivePathMatch(t *testing.T) {
	result := UpdateAdvancedSettingsEvasivePathMatchResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsEvasivePathMatch/AdvancedSettingsEvasivePathMatch.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateAdvancedSettingsEvasivePathMatchRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsEvasivePathMatch/AdvancedSettingsEvasivePathMatch.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsEvasivePathMatchRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsEvasivePathMatchResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsEvasivePathMatchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/evasive-path-match",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsEvasivePathMatchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating AdvancedSettingsEvasivePathMatch"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/evasive-path-match",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsEvasivePathMatch",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAdvancedSettingsEvasivePathMatch(
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
