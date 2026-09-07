package botman

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Get BotAnalyticsSettingsValues
func TestBotman_GetBotAnalyticsSettingsValues(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]any
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
                {
                    "values": [
                        {"testKey":"testValue1"},
                        {"testKey":"testValue2"},
                        {"testKey":"testValue3"}
                    ]
                }`,
			expectedPath: "/appsec/v1/bot-analytics-settings/values",
			expectedResponse: map[string]any{
				"values": []any{
					map[string]any{"testKey": "testValue1"},
					map[string]any{"testKey": "testValue2"},
					map[string]any{"testKey": "testValue3"},
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
                {
                    "type": "internal_error",
                    "title": "Internal Server Error",
                    "detail": "Error fetching data",
                    "status": 500
                }`,
			expectedPath: "/appsec/v1/bot-analytics-settings/values",
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
			result, err := client.GetBotAnalyticsSettingsValues(
				session.ContextWithOptions(
					context.Background(),
				))
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
