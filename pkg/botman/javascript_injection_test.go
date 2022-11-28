package botman

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

// Test Get JavascriptInjection
func TestBotman_GetJavascriptInjection(t *testing.T) {
	tests := map[string]struct {
		params           GetJavascriptInjectionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"testKey":"testValue3"}`,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/javascript-injection",
			expectedResponse: map[string]interface{}{"testKey": "testValue3"},
		},
		"500 internal server error": {
			params: GetJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/javascript-injection",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching match target",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: GetJavascriptInjectionRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetJavascriptInjectionRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: GetJavascriptInjectionRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
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
			result, err := client.GetJavascriptInjection(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update JavascriptInjection.
func TestBotman_UpdateJavascriptInjection(t *testing.T) {
	tests := map[string]struct {
		params           UpdateJavascriptInjectionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/javascript-injection",
		},
		"500 internal server error": {
			params: UpdateJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/javascript-injection",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating zone",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: UpdateJavascriptInjectionRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateJavascriptInjectionRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: UpdateJavascriptInjectionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: json.RawMessage(`{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing JsonPayload": {
			params: UpdateJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
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
			result, err := client.UpdateJavascriptInjection(
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
