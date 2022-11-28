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

func TestAppSec_ListAdvancedSettingsLogging(t *testing.T) {

	result := GetAdvancedSettingsLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAdvancedSettingsLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsLoggingResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAdvancedSettingsLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/logging",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching AdvancedSettingsLogging",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/logging",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsLogging",
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
			result, err := client.GetAdvancedSettingsLogging(
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

// Test AdvancedSettingsLogging
func TestAppSec_GetAdvancedSettingsLogging(t *testing.T) {

	result := GetAdvancedSettingsLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAdvancedSettingsLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsLoggingResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/logging",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching AdvancedSettingsLogging"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/logging",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsLogging",
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
			result, err := client.GetAdvancedSettingsLogging(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update AdvancedSettingsLogging.
func TestAppSec_UpdateAdvancedSettingsLogging(t *testing.T) {
	result := UpdateAdvancedSettingsLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateAdvancedSettingsLoggingRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsLoggingResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/logging",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating AdvancedSettingsLogging"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/logging",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsLogging",
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
			result, err := client.UpdateAdvancedSettingsLogging(
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
