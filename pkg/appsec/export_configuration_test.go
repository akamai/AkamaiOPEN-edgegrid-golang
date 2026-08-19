package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListExportConfiguration(t *testing.T) {

	result := GetExportConfigurationResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestExportConfiguration/ExportConfiguration.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	aiRulesResult := GetExportConfigurationResponse{}
	aiRulesRespData := compactJSON(loadFixtureBytes("testdata/TestAIRule/ExportAIRules.json"))
	err = json.Unmarshal([]byte(aiRulesRespData), &aiRulesResult)
	require.NoError(t, err)

	aiRulesDisabledResult := GetExportConfigurationResponse{}
	aiRulesDisabledRespData := compactJSON(loadFixtureBytes("testdata/TestExportConfiguration/ExportAIRulesDisabled.json"))
	err = json.Unmarshal([]byte(aiRulesDisabledRespData), &aiRulesDisabledResult)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetExportConfigurationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetExportConfigurationResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetExportConfigurationRequest{
				ConfigID: 43253,
				Version:  15,
				Source:   "TF",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/export/configs/43253/versions/15?source=TF",
			expectedResponse: &result,
		},
		"200 OK - with AI rules data": {
			params: GetExportConfigurationRequest{
				ConfigID: 77653,
				Version:  25,
				Source:   "TF",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(aiRulesRespData),
			expectedPath:     "/appsec/v1/export/configs/77653/versions/25?source=TF",
			expectedResponse: &aiRulesResult,
		},
		"200 OK - DISABLED AI rules with empty list": {
			params: GetExportConfigurationRequest{
				ConfigID: 77653,
				Version:  25,
				Source:   "TF",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(aiRulesDisabledRespData),
			expectedPath:     "/appsec/v1/export/configs/77653/versions/25?source=TF",
			expectedResponse: &aiRulesDisabledResult,
		},
		"500 internal server error": {
			params: GetExportConfigurationRequest{
				ConfigID: 43253,
				Version:  15,
				Source:   "TF",
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
			expectedPath: "/appsec/v1/export/configs/43253/versions/15?source=TF",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetExportConfiguration(
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
