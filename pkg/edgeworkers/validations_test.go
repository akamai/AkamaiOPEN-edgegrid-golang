package edgeworkers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateBundle(t *testing.T) {
	tests := map[string]struct {
		params           ValidateBundleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ValidateBundleResponse
		withError        error
	}{
		"200 OK": {
			params:         ValidateBundleRequest{Bundle{strings.NewReader("a valid bundle")}},
			responseStatus: 200,
			responseBody: `
{
    "errors": [],
    "warnings": []
}`,
			expectedPath: "/edgeworkers/v1/validations",
			expectedResponse: &ValidateBundleResponse{
				Errors:   []ValidationIssue{},
				Warnings: []ValidationIssue{},
			},
		},
		"200 OK with invalid gzip format error": {
			params:         ValidateBundleRequest{Bundle{strings.NewReader("invalid bundle format")}},
			responseStatus: 200,
			responseBody: `
{
    "errors": [
        {
            "type": "INVALID_GZIP_FORMAT",
            "message": "invalid GZIP file format"
        }
    ],
    "warnings": []
}`,
			expectedPath: "/edgeworkers/v1/validations",
			expectedResponse: &ValidateBundleResponse{
				Errors: []ValidationIssue{
					{
						Type:    "INVALID_GZIP_FORMAT",
						Message: "invalid GZIP file format",
					},
				},
				Warnings: []ValidationIssue{},
			},
		},
		"200 OK with expiring token warning": {
			params:         ValidateBundleRequest{Bundle{bytes.NewReader([]byte("bundle with expiring token"))}},
			responseStatus: 200,
			responseBody: `
{
    "errors": [],
    "warnings": [
        {
            "type": "ACCESS_TOKEN_EXPIRING_SOON",
            "message": "token expiring soon"
        }
	]
}`,
			expectedPath: "/edgeworkers/v1/validations",
			expectedResponse: &ValidateBundleResponse{
				Errors: []ValidationIssue{},
				Warnings: []ValidationIssue{
					{
						Type:    "ACCESS_TOKEN_EXPIRING_SOON",
						Message: "token expiring soon",
					},
				},
			},
		}, "500 internal server error": {
			params:         ValidateBundleRequest{Bundle{bytes.NewReader([]byte("a valid bundle"))}},
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
			expectedPath: "/edgeworkers/v1/validations",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
			},
		},
		"missing bundle reader": {
			params:    ValidateBundleRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ValidateBundle(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
