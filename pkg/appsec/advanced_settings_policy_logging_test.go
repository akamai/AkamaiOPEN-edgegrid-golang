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

func TestApsec_ListAdvancedSettingsPolicyLogging(t *testing.T) {

	result := GetAdvancedSettingsPolicyLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPolicyLogging/AdvancedSettingsPolicyLogging.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAdvancedSettingsPolicyLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsPolicyLoggingResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAdvancedSettingsPolicyLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/advanced-settings/logging",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsPolicyLoggingRequest{
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
    "detail": "Error fetching AdvancedSettingsPolicyLogging",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/advanced-settings/logging",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsPolicyLogging",
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
			result, err := client.GetAdvancedSettingsPolicyLogging(
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

// Test AdvancedSettingsPolicyLogging
func TestAppSec_GetAdvancedSettingsPolicyLogging(t *testing.T) {

	result := GetAdvancedSettingsPolicyLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPolicyLogging/AdvancedSettingsPolicyLogging.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAdvancedSettingsPolicyLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsPolicyLoggingResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsPolicyLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/advanced-settings/logging",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsPolicyLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching AdvancedSettingsPolicyLogging"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/advanced-settings/logging",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsPolicyLogging",
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
			result, err := client.GetAdvancedSettingsPolicyLogging(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update AdvancedSettingsPolicyLogging.
func TestAppSec_UpdateAdvancedSettingsPolicyLogging(t *testing.T) {
	result := UpdateAdvancedSettingsPolicyLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPolicyLogging/AdvancedSettingsPolicyLogging.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateAdvancedSettingsPolicyLoggingRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPolicyLogging/AdvancedSettingsPolicyLogging.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsPolicyLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsPolicyLoggingResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsPolicyLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/advanced-settings/logging",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsPolicyLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating AdvancedSettingsPolicyLogging"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/advanced-settings/logging",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsPolicyLogging",
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
			result, err := client.UpdateAdvancedSettingsPolicyLogging(
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
