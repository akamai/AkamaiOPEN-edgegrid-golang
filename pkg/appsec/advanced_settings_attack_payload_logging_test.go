package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListAdvancedSettingsAttackPayloadLogging(t *testing.T) {

	result := GetAdvancedSettingsAttackPayloadLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsAttackPayloadLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsAttackPayloadLoggingResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/logging/attack-payload",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error fetching AdvancedSettingsAttackPayloadLogging",
					"status": 500
				}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/logging/attack-payload",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsAttackPayloadLogging",
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
			result, err := client.GetAdvancedSettingsAttackPayloadLogging(
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
func TestAppSec_GetAdvancedSettingsAttackPayloadLoggingPolicy(t *testing.T) {

	result := GetAdvancedSettingsAttackPayloadLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsAttackPayloadLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsAttackPayloadLoggingResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "test_policy",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/logging/attack-payload",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "test_policy",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
							{
								"type": "internal_error",
								"title": "Internal Server Error",
								"detail": "Error fetching AdvancedSettingsLogging"
							}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/logging/attack-payload",
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
			result, err := client.GetAdvancedSettingsAttackPayloadLogging(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update AdvancedSettingsAttackPayloadLogging config level.
func TestAppSec_UpdateAdvancedSettingsAttackPayloadLogging(t *testing.T) {
	result := UpdateAdvancedSettingsAttackPayloadLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsAttackPayloadLoggingRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsAttackPayloadLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsAttackPayloadLoggingResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/logging/attack-payload",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsAttackPayloadLogging"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/logging/attack-payload",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsAttackPayloadLogging",
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
			result, err := client.UpdateAdvancedSettingsAttackPayloadLogging(
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

// Test Update AdvancedSettingsAttackPayloadLogging policy level.
func TestAppSec_UpdateAdvancedSettingsAttackPayloadLoggingPolicy(t *testing.T) {
	result := UpdateAdvancedSettingsAttackPayloadLoggingResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsAttackPayloadLoggingRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsAttackPayloadLoggingRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsAttackPayloadLoggingResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "test_policy",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/logging/attack-payload",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "test_policy",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsAttackPayloadLogging"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/logging/attack-payload",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsAttackPayloadLogging",
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
			result, err := client.UpdateAdvancedSettingsAttackPayloadLogging(
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
