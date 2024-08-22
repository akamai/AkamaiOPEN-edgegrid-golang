package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiGetClientSettings(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ClientSettingsBody
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "ruleFormat": "v2015-08-08",
    "usePrefixes": true
}
`,
			expectedPath: "/papi/v1/client-settings",
			expectedResponse: &ClientSettingsBody{
				RuleFormat:  "v2015-08-08",
				UsePrefixes: true,
			},
		},
		"500 server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching client settings",
    "status": 500
}
`,
			expectedPath: "/papi/v1/client-settings",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching client settings",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
			result, err := client.GetClientSettings(context.Background())
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapiUpdateClientSettings(t *testing.T) {
	tests := map[string]struct {
		params           ClientSettingsBody
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ClientSettingsBody
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ClientSettingsBody{
				RuleFormat:  "v2015-08-08",
				UsePrefixes: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "ruleFormat": "v2015-08-08",
    "usePrefixes": true
}
`,
			expectedPath: "/papi/v1/client-settings",
			expectedResponse: &ClientSettingsBody{
				RuleFormat:  "v2015-08-08",
				UsePrefixes: true,
			},
		},
		"500 OK": {
			params: ClientSettingsBody{
				RuleFormat:  "v2015-08-08",
				UsePrefixes: true,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching client settings",
    "status": 500
}
`,
			expectedPath: "/papi/v1/client-settings",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching client settings",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateClientSettings(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
