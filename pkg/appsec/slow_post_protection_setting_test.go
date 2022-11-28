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

func TestAppSec_ListSlowPostProtectionSettings(t *testing.T) {

	result := GetSlowPostProtectionSettingsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSlowPostProtectionSettings/SlowPostProtectionSettings.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetSlowPostProtectionSettingsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSlowPostProtectionSettingsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetSlowPostProtectionSettingsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/slow-post",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetSlowPostProtectionSettingsRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/slow-post",
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
			result, err := client.GetSlowPostProtectionSettings(
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

// Test SlowPostProtectionSetting
func TestAppSec_GetSlowPostProtectionSetting(t *testing.T) {

	result := GetSlowPostProtectionSettingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSlowPostProtectionSettings/SlowPostProtectionSettings.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetSlowPostProtectionSettingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSlowPostProtectionSettingResponse
		withError        error
	}{
		"200 OK": {
			params: GetSlowPostProtectionSettingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/slow-post",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetSlowPostProtectionSettingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/slow-post",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching match target",
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
			result, err := client.GetSlowPostProtectionSetting(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update SlowPostProtectionSetting.
func TestAppSec_UpdateSlowPostProtectionSetting(t *testing.T) {
	result := UpdateSlowPostProtectionSettingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSlowPostProtectionSettings/SlowPostProtectionSettings.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateSlowPostProtectionSettingRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestSlowPostProtectionSettings/SlowPostProtectionSettings.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateSlowPostProtectionSettingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateSlowPostProtectionSettingResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateSlowPostProtectionSettingRequest{
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
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/slow-post",
		},
		"500 internal server error": {
			params: UpdateSlowPostProtectionSettingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/slow-post",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
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
			result, err := client.UpdateSlowPostProtectionSetting(
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
