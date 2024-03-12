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

func TestGetAccessKeyStatus(t *testing.T) {
	tests := map[string]struct {
		params           GetAccessKeyStatusRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAccessKeyStatusResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetAccessKeyStatusRequest{
				RequestID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"accessKey": {
		"accessKeyUid": 123,
		"link": "/cam/v1/access-keys/123"
	  },
	  "accessKeyVersion": {
		"accessKeyUid": 123,
		"link": "/cam/v1/access-keys/123/versions/1",
		"version": 1
	  },
	  "processingStatus": "IN_PROGRESS",
	  "request": {
		"accessKeyName": "TestAccessKeyName",
		"authenticationMethod": "AWS4_HMAC_SHA256",
		"contractId": "TestContractID",
		"groupId": 123,
		"networkConfiguration": {
		  "additionalCdn": "CHINA_CDN",
		  "securityNetwork": "ENHANCED_TLS"
		}
	  },
	  "requestDate": "2021-02-26T13:34:36.715643Z",
	  "requestId": 1,
	  "requestedBy": "user"
}`,
			expectedPath: "/cam/v1/access-key-create-requests/1",
			expectedResponse: &GetAccessKeyStatusResponse{
				ProcessingStatus: ProcessingInProgress,
				RequestDate:      "2021-02-26T13:34:36.715643Z",
				RequestID:        1,
				RequestedBy:      "user",
				AccessKey: &KeyLink{
					AccessKeyUID: 123,
					Link:         "/cam/v1/access-keys/123",
				},
				AccessKeyVersion: &KeyVersion{
					AccessKeyUID: 123,
					Link:         "/cam/v1/access-keys/123/versions/1",
					Version:      1,
				},
				Request: &RequestInformation{
					AccessKeyName:        "TestAccessKeyName",
					AuthenticationMethod: AuthAWS,
					ContractID:           "TestContractID",
					GroupID:              123,
					NetworkConfiguration: &SecureNetwork{
						AdditionalCDN:   ChinaCDN,
						SecurityNetwork: NetworkEnhanced,
					},
				},
			},
		},
		"200 OK - minimal": {
			params: GetAccessKeyStatusRequest{
				RequestID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"accessKey": null,
    "accessKeyVersion": null,
    "processingStatus": "IN_PROGRESS",
    "request": null,
    "requestDate": "2021-02-26T13:34:36.715643Z",
    "requestId": 1,
    "requestedBy": "user"
}`,
			expectedPath: "/cam/v1/access-key-create-requests/1",
			expectedResponse: &GetAccessKeyStatusResponse{
				ProcessingStatus: ProcessingInProgress,
				RequestDate:      "2021-02-26T13:34:36.715643Z",
				RequestID:        1,
				RequestedBy:      "user",
			},
		},
		"missing required params - validation error": {
			params: GetAccessKeyStatusRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get the status of an access key: struct validation: RequestID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: GetAccessKeyStatusRequest{
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
			expectedPath: "/cam/v1/access-key-create-requests/123",
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
			result, err := client.GetAccessKeyStatus(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
