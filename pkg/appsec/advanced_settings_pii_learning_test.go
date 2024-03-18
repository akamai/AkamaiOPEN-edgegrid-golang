package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppsec_GetAdvancedSettingsPIILearning(t *testing.T) {

	result := AdvancedSettingsPIILearningResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPIILearning/AdvancedSettingsPIILearning.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsPIILearningRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *AdvancedSettingsPIILearningResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsPIILearningRequest{
				ConfigVersion: ConfigVersion{
					ConfigID: 43253,
					Version:  15,
				},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/pii-learning",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsPIILearningRequest{
				ConfigVersion: ConfigVersion{
					ConfigID: 43253,
					Version:  15,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching AdvancedSettingsPIILearning",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/pii-learning",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsPIILearning",
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
			result, err := client.GetAdvancedSettingsPIILearning(
				session.ContextWithOptions(context.Background()),
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

func TestAppSec_UpdateAdvancedSettingsPIILearning(t *testing.T) {
	result := AdvancedSettingsPIILearningResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPIILearning/AdvancedSettingsPIILearning.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsPIILearningRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPIILearning/AdvancedSettingsPIILearning.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params                UpdateAdvancedSettingsPIILearningRequest
		responseStatus        int
		responseBody          string
		expectedPath          string
		expectedResponse      *AdvancedSettingsPIILearningResponse
		withError             error
		expectedBody          string
		expectValidationError bool
	}{
		"validation error": {
			params: UpdateAdvancedSettingsPIILearningRequest{
				ConfigVersion: ConfigVersion{
					ConfigID: 43253,
				},
				EnablePIILearning: false,
			},
			expectValidationError: true,
		},
		"200 Success": {
			params: UpdateAdvancedSettingsPIILearningRequest{
				ConfigVersion: ConfigVersion{
					ConfigID: 43253,
					Version:  15,
				},
				EnablePIILearning: true,
			},
			expectedBody: `
		{
			"enablePiiLearning": true
		}`,
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/pii-learning",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsPIILearningRequest{
				ConfigVersion: ConfigVersion{
					ConfigID: 43253,
					Version:  15,
				},
				EnablePIILearning: true,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
					{
						"type": "internal_error",
						"title": "Internal Server Error",
						"detail": "Error creating AdvancedSettingsPIILearning"
					}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/pii-learning",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsPIILearning",
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
				if test.withError == nil {
					var reqBody interface{}
					err = json.NewDecoder(r.Body).Decode(&reqBody)
					require.NoError(t, err, "Error while decoding request body")

					var expectedBody interface{}
					err = json.Unmarshal([]byte(test.expectedBody), &expectedBody)
					require.NoError(t, err, "Error while parsing expected body to JSON")

					assert.Equal(t, expectedBody, reqBody)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAdvancedSettingsPIILearning(
				session.ContextWithOptions(context.Background()),
				test.params)
			if test.expectValidationError {
				assert.True(t, strings.Contains(err.Error(), "struct validation"))
				return
			}
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
