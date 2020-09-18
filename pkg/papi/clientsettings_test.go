package papi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestPapi_GetClientSettings(t *testing.T) {
	tests := map[string]struct {
		//params ClientSettingsBody
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
		"404 not found": {
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "not_found",
    "title": "Not Found",
    "detail": "Could not find client settings",
    "status": 404
}
`,
			expectedPath: "/papi/v1/client-settings",
			/*expectedResponse: &ClientSettingsBody{
				RuleFormat:  "v2015-08-08",
				UsePrefixes: true,
			},*/
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

/*
func TestPapi_UpdateClientSettings(t *testing.T) {

}
*/
