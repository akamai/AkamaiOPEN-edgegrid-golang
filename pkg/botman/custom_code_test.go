package botman

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Get CustomCode
func TestBotman_GetCustomCode(t *testing.T) {
	tests := map[string]struct {
		params           GetCustomCodeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetCustomCodeRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"testKey":"testValue3"}`,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/transactional-endpoint-protection/custom-code",
			expectedResponse: map[string]interface{}{"testKey": "testValue3"},
		},
		"500 internal server error": {
			params: GetCustomCodeRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/transactional-endpoint-protection/custom-code",
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
			params: GetCustomCodeRequest{
				Version: 15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetCustomCodeRequest{
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
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCustomCode(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update CustomCode.
func TestBotman_UpdateCustomCode(t *testing.T) {
	tests := map[string]struct {
		params           UpdateCustomCodeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateCustomCodeRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/transactional-endpoint-protection/custom-code",
		},
		"500 internal server error": {
			params: UpdateCustomCodeRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error updating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/transactional-endpoint-protection/custom-code",
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
			params: UpdateCustomCodeRequest{
				Version:     15,
				JsonPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateCustomCodeRequest{
				ConfigID:    43253,
				JsonPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing JsonPayload": {
			params: UpdateCustomCodeRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "JsonPayload")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateCustomCode(
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
