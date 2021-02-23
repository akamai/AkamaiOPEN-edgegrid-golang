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

// Test ContractsGroups
func TestAppSec_GetContractsGroups(t *testing.T) {

	result := GetContractsGroupsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestContractsGroups/ContractsGroups.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetContractsGroupsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetContractsGroupsResponse
		withError        error
	}{
		"200 OK": {
			params: GetContractsGroupsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/contracts-groups",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetContractsGroupsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching ContractsGroups"
}`),
			expectedPath: "/appsec/v1/contracts-groups",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ContractsGroups",
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
			result, err := client.GetContractsGroups(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
