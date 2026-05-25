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

func TestAppSec_GetAdvancedSettingsJA4Fingerprint(t *testing.T) {

	result := GetAdvancedSettingsJA4FingerprintResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsJA4Fingerprint/AdvancedSettingsJA4Fingerprint.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsJA4FingerprintRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsJA4FingerprintResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsJA4FingerprintRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/ja4-fingerprint",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsJA4FingerprintRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching AdvancedSettingsJA4Fingerprint"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/ja4-fingerprint",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsJA4Fingerprint",
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
			result, err := client.GetAdvancedSettingsJA4Fingerprint(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateAdvancedSettingsJA4Fingerprint(t *testing.T) {
	result := UpdateAdvancedSettingsJA4FingerprintResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsJA4Fingerprint/AdvancedSettingsJA4Fingerprint.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsJA4FingerprintRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsJA4Fingerprint/AdvancedSettingsJA4Fingerprint.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params              UpdateAdvancedSettingsJA4FingerprintRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *UpdateAdvancedSettingsJA4FingerprintResponse
		withError           error
		headers             http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsJA4FingerprintRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:      http.StatusCreated,
			responseBody:        respData,
			expectedResponse:    &result,
			expectedRequestBody: `{"headerNames":["ja4-fingerprint"]}`,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/advanced-settings/ja4-fingerprint",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsJA4FingerprintRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsJA4Fingerprint"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/ja4-fingerprint",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsJA4Fingerprint",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAdvancedSettingsJA4Fingerprint(
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

func TestAppSec_RemoveAdvancedSettingsJA4Fingerprint(t *testing.T) {
	result := RemoveAdvancedSettingsJA4FingerprintResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsJA4Fingerprint/AdvancedSettingsJA4Fingerprint.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := RemoveAdvancedSettingsJA4FingerprintRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsJA4Fingerprint/AdvancedSettingsJA4Fingerprint.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           RemoveAdvancedSettingsJA4FingerprintRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveAdvancedSettingsJA4FingerprintResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveAdvancedSettingsJA4FingerprintRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/ja4-fingerprint",
		},
		"500 internal server error": {
			params: RemoveAdvancedSettingsJA4FingerprintRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsJA4Fingerprint"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/ja4-fingerprint",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsJA4Fingerprint",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.RemoveAdvancedSettingsJA4Fingerprint(
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
