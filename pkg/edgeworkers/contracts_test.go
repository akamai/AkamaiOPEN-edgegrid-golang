package edgeworkers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListContracts(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListContractsResponse
		withError        error
	}{
		"200 OK": {
			responseStatus: 200,
			responseBody: `
{
	"contractIds": [
		"1-599K",
		"B-M-28QYF3M"
	]
}`,
			expectedPath: "/edgeworkers/v1/contracts",
			expectedResponse: &ListContractsResponse{
				[]string{
					"1-599K",
					"B-M-28QYF3M",
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/edgeworkers/error-instances/abc",
    "status": 500,
    "errorCode": "EW4303"
}`,
			expectedPath: "/edgeworkers/v1/contracts",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
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
			result, err := client.ListContracts(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
