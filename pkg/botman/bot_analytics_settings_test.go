package botman

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Get BotAnalyticsSettings
func TestBotman_GetBotAnalyticsSettings(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params           GetBotAnalyticsSettingsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]any
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetBotAnalyticsSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"testKey":"testValue3"}`,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/bot-analytics-settings",
			expectedResponse: map[string]any{"testKey": "testValue3"},
		},
		"500 internal server error": {
			params: GetBotAnalyticsSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/bot-analytics-settings",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: GetBotAnalyticsSettingsRequest{
				Version: 15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetBotAnalyticsSettingsRequest{
				ConfigID: 43253,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetBotAnalyticsSettings(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update BotAnalyticsSettings.
func TestBotman_UpdateBotAnalyticsSettings(t *testing.T) {
	tests := map[string]struct {
		params              UpdateBotAnalyticsSettingsRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    map[string]any
		withError           func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateBotAnalyticsSettingsRequest{
				ConfigID:    43253,
				Version:     15,
				JSONPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			expectedRequestBody: `{"testKey":"testValue3"}`,
			responseStatus:      http.StatusOK,
			responseBody:        `{"testKey":"testValue3"}`,
			expectedResponse:    map[string]any{"testKey": "testValue3"},
			expectedPath:        "/appsec/v1/configs/43253/versions/15/advanced-settings/bot-analytics-settings",
		},
		"500 internal server error": {
			params: UpdateBotAnalyticsSettingsRequest{
				ConfigID:    43253,
				Version:     15,
				JSONPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error updating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/bot-analytics-settings",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error updating data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: UpdateBotAnalyticsSettingsRequest{
				Version:     15,
				JSONPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateBotAnalyticsSettingsRequest{
				ConfigID:    43253,
				JSONPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing JSONPayload": {
			params: UpdateBotAnalyticsSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "JSONPayload")
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateBotAnalyticsSettings(
				session.ContextWithOptions(
					context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
