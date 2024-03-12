package cloudaccess

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccessKeyVersionStatus(t *testing.T) {
	tests := map[string]struct {
		params           GetAccessKeyVersionStatusRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAccessKeyVersionStatusResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetAccessKeyVersionStatusRequest{
				RequestID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessKeyVersion": {
    "accessKeyUid": 123,
    "link": "/cam/v1/access-keys/123/versions/2",
    "version": 2
  },
  "processingStatus": "IN_PROGRESS",
  "requestDate": "2021-02-26T14:54:38.622074Z",
  "requestedBy": "user"
}`,
			expectedPath: "/cam/v1/access-key-version-create-requests/1",
			expectedResponse: &GetAccessKeyVersionStatusResponse{
				ProcessingStatus: ProcessingInProgress,
				RequestDate:      "2021-02-26T14:54:38.622074Z",
				RequestedBy:      "user",
				AccessKeyVersion: &KeyVersion{
					AccessKeyUID: 123,
					Link:         "/cam/v1/access-keys/123/versions/2",
					Version:      2,
				},
			},
		},
		"200 OK - minimal": {
			params: GetAccessKeyVersionStatusRequest{
				RequestID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessKeyVersion": null,
  "processingStatus": "IN_PROGRESS",
  "requestDate": "2021-02-26T14:54:38.622074Z",
  "requestedBy": "user"
}`,
			expectedPath: "/cam/v1/access-key-version-create-requests/1",
			expectedResponse: &GetAccessKeyVersionStatusResponse{
				ProcessingStatus: ProcessingInProgress,
				RequestDate:      "2021-02-26T14:54:38.622074Z",
				RequestedBy:      "user",
			},
		},
		"missing required params - validation error": {
			params: GetAccessKeyVersionStatusRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get the status of an access key version: struct validation: RequestID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: GetAccessKeyVersionStatusRequest{
				RequestID: 123,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal-server-error",
    "title": "Internal Server Error",
    "detail": "Error processing request",
    "instance": "TestInstances",
    "status": 500
}`,
			expectedPath: "/cam/v1/access-key-version-create-requests/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
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
			result, err := client.GetAccessKeyVersionStatus(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
