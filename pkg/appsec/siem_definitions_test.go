package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test SiemDefinitions
func TestAppSec_GetSiemDefinitions(t *testing.T) {

	result := GetSiemDefinitionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSiemDefinitions/SiemDefinitions.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetSiemDefinitionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSiemDefinitionsResponse
		withError        error
	}{
		"200 OK": {
			params:           GetSiemDefinitionsRequest{},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/siem-definitions",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetSiemDefinitionsRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching SiemDefinitions"
}`),
			expectedPath: "/appsec/v1/siem-definitions",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching SiemDefinitions",
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
			result, err := client.GetSiemDefinitions(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
