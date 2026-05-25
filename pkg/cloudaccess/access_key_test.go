package cloudaccess

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
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
				"accessKey": 
					{
						"accessKeyUid": 123,
						"link": "/cam/v1/access-keys/123"
					},
				"accessKeyVersion": 
					{
						"accessKeyUid": 123,
						"link": "/cam/v1/access-keys/123/versions/1",
						"version": 1
					},
				"processingStatus": "IN_PROGRESS",
				"request": 
					{
						"accessKeyName": "TestAccessKeyName",
						"authenticationMethod": "AWS4_HMAC_SHA256",
						"contractId": "TestContractID",
						"groupId": 123,
						"networkConfiguration": 
							{
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
				AccessKey:        &KeyLink{AccessKeyUID: 123, Link: "/cam/v1/access-keys/123"},
				AccessKeyVersion: &KeyVersion{AccessKeyUID: 123, Link: "/cam/v1/access-keys/123/versions/1", Version: 1},
				ProcessingStatus: ProcessingInProgress,
				Request: &RequestInformation{
					AccessKeyName:        "TestAccessKeyName",
					AuthenticationMethod: AuthAWS,
					ContractID:           "TestContractID",
					GroupID:              123,
					NetworkConfiguration: &SecureNetwork{
						AdditionalCDN:   ptr.To(ChinaCDN),
						SecurityNetwork: NetworkEnhanced,
					},
				},
				RequestDate: time.Date(2021, 2, 26, 13, 34, 36, 715643000, time.UTC),
				RequestID:   1,
				RequestedBy: "user",
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
				RequestDate:      time.Date(2021, 2, 26, 13, 34, 36, 715643000, time.UTC),
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
			defer mockServer.Close()
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

func TestCreateAccessKey(t *testing.T) {
	tests := map[string]struct {
		accessKey           CreateAccessKeyRequest
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedResponse    *CreateAccessKeyResponse
		responseHeaders     map[string]string
		withError           func(*testing.T, error)
	}{
		"202 Accepted": {
			accessKey: CreateAccessKeyRequest{
				AccessKeyName:        "key1",
				AuthenticationMethod: string(AuthAWS),
				ContractID:           "TestContractID",
				Credentials: Credentials{
					CloudAccessKeyID:     "456",
					CloudSecretAccessKey: "testKey",
				},
				GroupID: 123,
				NetworkConfiguration: SecureNetwork{
					AdditionalCDN:   ptr.To(ChinaCDN),
					SecurityNetwork: NetworkEnhanced,
				},
			},
			expectedPath:   "/cam/v1/access-keys",
			responseStatus: http.StatusAccepted,
			responseBody: `
			{
				"requestId": 195,
				"retryAfter": 4
			}`,
			expectedRequestBody: `
			{
				"accessKeyName": "key1",
				"authenticationMethod": "AWS4_HMAC_SHA256",
				"contractId": "TestContractID",
				"credentials": 
					{
						"cloudAccessKeyId": "456",
						"cloudSecretAccessKey": "testKey"
					},
				"groupId": 123,
				"networkConfiguration": 
					{
						"additionalCdn": "CHINA_CDN",
						"securityNetwork": "ENHANCED_TLS"
					}
			}`,
			expectedResponse: &CreateAccessKeyResponse{
				RequestID:  195,
				RetryAfter: 4,
				Location:   "https://abc.com",
			},
			responseHeaders: map[string]string{
				"Location": "https://abc.com",
			},
		},
		"cloudAccessKeyID not required for VP_QUEUE_IT": {
			accessKey: CreateAccessKeyRequest{
				AccessKeyName:        "vp-queue-it-key",
				AuthenticationMethod: string(AuthVPQueueIt),
				ContractID:           "TestContractID",
				Credentials: Credentials{
					CloudSecretAccessKey: "testKey",
				},
				GroupID: 123,
				NetworkConfiguration: SecureNetwork{
					SecurityNetwork: NetworkEnhanced,
				},
			},
			expectedPath:   "/cam/v1/access-keys",
			responseStatus: http.StatusAccepted,
			responseBody: `
			{
				"requestId": 196,
				"retryAfter": 3
			}`,
			expectedRequestBody: `
			{
				"accessKeyName": "vp-queue-it-key",
				"authenticationMethod": "VP_QUEUE_IT",
				"contractId": "TestContractID",
				"credentials": 
					{
						"cloudSecretAccessKey": "testKey"
					},
				"groupId": 123,
				"networkConfiguration": 
					{
						"securityNetwork": "ENHANCED_TLS"
					}
			}`,
			expectedResponse: &CreateAccessKeyResponse{
				RequestID:  196,
				RetryAfter: 3,
				Location:   "https://example.com/cam/v1/access-key-create-requests/196",
			},
			responseHeaders: map[string]string{
				"Location": "https://example.com/cam/v1/access-key-create-requests/196",
			},
		},
		"cloudAccessKeyID not required for AVM_CLOUDINARY": {
			accessKey: CreateAccessKeyRequest{
				AccessKeyName:        "avm-cloudinary-key",
				AuthenticationMethod: string(AuthAVMCloudinary),
				ContractID:           "TestContractID",
				Credentials: Credentials{
					CloudSecretAccessKey: "testKey",
				},
				GroupID: 123,
				NetworkConfiguration: SecureNetwork{
					SecurityNetwork: NetworkEnhanced,
				},
			},
			expectedPath:   "/cam/v1/access-keys",
			responseStatus: http.StatusAccepted,
			responseBody: `
			{
				"requestId": 197,
				"retryAfter": 3
			}`,
			expectedRequestBody: `
			{
				"accessKeyName": "avm-cloudinary-key",
				"authenticationMethod": "AVM_CLOUDINARY",
				"contractId": "TestContractID",
				"credentials": 
					{
						"cloudSecretAccessKey": "testKey"
					},
				"groupId": 123,
				"networkConfiguration": 
					{
						"securityNetwork": "ENHANCED_TLS"
					}
			}`,
			expectedResponse: &CreateAccessKeyResponse{
				RequestID:  197,
				RetryAfter: 3,
				Location:   "https://example.com/cam/v1/access-key-create-requests/197",
			},
			responseHeaders: map[string]string{
				"Location": "https://example.com/cam/v1/access-key-create-requests/197",
			},
		},
		"missing required request body - validation error": {
			accessKey: CreateAccessKeyRequest{
				Credentials:          Credentials{},
				NetworkConfiguration: SecureNetwork{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create an access key: struct validation: AccessKeyName: cannot be blank\nAuthenticationMethod: cannot be blank\nCloudAccessKeyID: cannot be blank\nCloudSecretAccessKey: cannot be blank\nContractID: cannot be blank\nGroupID: cannot be blank\nSecurityNetwork: cannot be blank", err.Error())
			},
		},
		"invalid Authentication Method - validation error": {
			accessKey: CreateAccessKeyRequest{
				AccessKeyName:        "key1",
				AuthenticationMethod: "AVM_CLOUDINARY_INVALID",
				ContractID:           "TestContractID",
				Credentials: Credentials{
					CloudAccessKeyID:     "456",
					CloudSecretAccessKey: "testKey",
				},
				GroupID: 123,
				NetworkConfiguration: SecureNetwork{
					AdditionalCDN:   ptr.To(ChinaCDN),
					SecurityNetwork: NetworkEnhanced,
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create an access key: struct validation: AuthenticationMethod: must be a valid value", err.Error())
			},
		},
		"cloudAccessKeyID missing when required - validation error": {
			accessKey: CreateAccessKeyRequest{
				AccessKeyName:        "key-without-cloud-id",
				AuthenticationMethod: string(AuthAWS),
				ContractID:           "TestContractID",
				Credentials: Credentials{
					CloudSecretAccessKey: "testSecret",
				},
				GroupID: 123,
				NetworkConfiguration: SecureNetwork{
					SecurityNetwork: NetworkEnhanced,
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create an access key: struct validation: CloudAccessKeyID: cannot be blank", err.Error())
			},
		},
		"additionalCDN present when not allowed for VP_QUEUE_IT - validation error": {
			accessKey: CreateAccessKeyRequest{
				AccessKeyName:        "vp-queue-it-key-with-cdn",
				AuthenticationMethod: string(AuthVPQueueIt),
				ContractID:           "TestContractID",
				Credentials: Credentials{
					CloudSecretAccessKey: "testKey",
				},
				GroupID: 123,
				NetworkConfiguration: SecureNetwork{
					AdditionalCDN:   ptr.To(ChinaCDN),
					SecurityNetwork: NetworkEnhanced,
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create an access key: struct validation: AdditionalCDN: must be blank", err.Error())
			},
		},
		"additionalCDN present when not allowed for AVM_CLOUDINARY - validation error": {
			accessKey: CreateAccessKeyRequest{
				AccessKeyName:        "avm-cloudinary-key-with-cdn",
				AuthenticationMethod: string(AuthAVMCloudinary),
				ContractID:           "TestContractID",
				Credentials: Credentials{
					CloudSecretAccessKey: "testKey",
				},
				GroupID: 123,
				NetworkConfiguration: SecureNetwork{
					AdditionalCDN:   ptr.To(ChinaCDN),
					SecurityNetwork: NetworkEnhanced,
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create an access key: struct validation: AdditionalCDN: must be blank", err.Error())
			},
		},
		"500 internal server error": {
			accessKey: CreateAccessKeyRequest{
				AccessKeyName:        "key1",
				AuthenticationMethod: string(AuthAWS),
				ContractID:           "TestContractID",
				Credentials: Credentials{
					CloudAccessKeyID:     "456",
					CloudSecretAccessKey: "testKey",
				},
				GroupID: 123,
				NetworkConfiguration: SecureNetwork{
					AdditionalCDN:   ptr.To(ChinaCDN),
					SecurityNetwork: NetworkEnhanced,
				},
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
			expectedPath: "/cam/v1/access-keys",
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
				assert.Equal(t, http.MethodPost, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				if len(test.responseHeaders) > 0 {
					for header, value := range test.responseHeaders {
						w.Header().Set(header, value)
					}
				}

				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateAccessKey(context.Background(), test.accessKey)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetAccessKey(t *testing.T) {

	tests := map[string]struct {
		params           AccessKeyRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *GetAccessKeyResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: AccessKeyRequest{
				AccessKeyUID: 1,
			},
			expectedPath:   "/cam/v1/access-keys/1",
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accessKeyUid": 1,
				"accessKeyName": "key1",
				"authenticationMethod": "AWS4_HMAC_SHA256",
				"createdBy": "user1",
				"groups": [
					{
						"contractIds": ["TestContractID"],
						"groupId": 123
					}
				],
				"latestVersion": 1,
				"networkConfiguration": 
					{
						"additionalCdn": "RUSSIA_CDN",
						"securityNetwork": "ENHANCED_TLS"
					},
				"note": "some note"
			}`,
			expectedResponse: &GetAccessKeyResponse{
				AccessKeyUID:         1,
				AccessKeyName:        "key1",
				AuthenticationMethod: "AWS4_HMAC_SHA256",
				NetworkConfiguration: &SecureNetwork{
					AdditionalCDN:   ptr.To(RussiaCDN),
					SecurityNetwork: NetworkEnhanced,
				},
				LatestVersion: 1,
				Groups: []Group{
					{
						ContractIDs: []string{"TestContractID"},
						GroupID:     123,
					},
				},
				CreatedBy: "user1",
				Note:      ptr.To("some note"),
			},
		},
		"missing required params - validation error": {
			params: AccessKeyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get an access key: struct validation: AccessKeyUID: cannot be blank", err.Error())
			},
		},
		"404 access key not found - custom error check": {
			params: AccessKeyRequest{
				AccessKeyUID: 2,
			},
			expectedPath:   "/cam/v1/access-keys/2",
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"type": "/cam/error-types/access-key-does-not-exist",
				"title": "Domain Error",
				"detail": "Access key with accessKeyUID '2' does not exist.",
				"instance": "test-instance-123",
				"status": 404,
				"accessKeyUid": 2
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrAccessKeyNotFound))
			},
		},
		"500 internal server error": {
			params: AccessKeyRequest{
				AccessKeyUID: 1,
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
			expectedPath: "/cam/v1/access-keys/1",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAccessKey(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListAccessKey(t *testing.T) {

	tests := map[string]struct {
		params           ListAccessKeysRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *ListAccessKeysResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListAccessKeysRequest{
				VersionGUID: "1",
			},
			expectedPath:   "/cam/v1/access-keys?versionGuid=1",
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accessKeys": [
					{
						"accessKeyUid": 1,
						"accessKeyName": "key1",
						"authenticationMethod": "AWS4_HMAC_SHA256",
						"createdBy": "user1",
						"groups": [
							{
								"contractIds": ["TestContractID"],
								"groupId": 123
							}
						],
						"latestVersion": 1,
						"networkConfiguration": 
							{
								"additionalCdn": "RUSSIA_CDN",
								"securityNetwork": "ENHANCED_TLS"
							}
					}
				]
			}`,
			expectedResponse: &ListAccessKeysResponse{
				AccessKeys: []AccessKeyResponse{
					{
						AccessKeyUID:         1,
						AccessKeyName:        "key1",
						AuthenticationMethod: "AWS4_HMAC_SHA256",
						NetworkConfiguration: &SecureNetwork{
							AdditionalCDN:   ptr.To(RussiaCDN),
							SecurityNetwork: NetworkEnhanced,
						},
						LatestVersion: 1,
						Groups: []Group{
							{
								ContractIDs: []string{"TestContractID"},
								GroupID:     123,
							},
						},
						CreatedBy: "user1",
					},
				},
			},
		},
		"500 internal server error": {
			params: ListAccessKeysRequest{
				VersionGUID: "1",
			},
			expectedPath:   "/cam/v1/access-keys?versionGuid=1",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal-server-error",
				"title": "Internal Server Error",
				"detail": "Error processing request",
				"instance": "TestInstances",
				"status": 500
			}`,
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListAccessKeys(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteAccessKey(t *testing.T) {

	tests := map[string]struct {
		params         AccessKeyRequest
		expectedPath   string
		responseStatus int
		responseBody   string
		withError      func(*testing.T, error)
	}{
		"204 No Content": {
			params: AccessKeyRequest{
				AccessKeyUID: 1,
			},
			expectedPath:   "/cam/v1/access-keys/1",
			responseStatus: http.StatusNoContent,
		},
		"missing required params - validation error": {
			params: AccessKeyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete an access key: struct validation: AccessKeyUID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: AccessKeyRequest{
				AccessKeyUID: 1,
			},
			expectedPath:   "/cam/v1/access-keys/1",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal-server-error",
				"title": "Internal Server Error",
				"detail": "Error processing request",
				"instance": "TestInstances",
				"status": 500
			}`,
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
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			err := client.DeleteAccessKey(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdateAccessKey(t *testing.T) {
	tests := map[string]struct {
		accessKey           UpdateAccessKeyRequest
		params              AccessKeyRequest
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedResponse    *UpdateAccessKeyResponse
		withError           func(*testing.T, error)
	}{
		"201 OK": {
			accessKey: UpdateAccessKeyRequest{
				AccessKeyName: "key2",
			},
			params: AccessKeyRequest{
				AccessKeyUID: 1,
			},
			expectedPath: "/cam/v1/access-keys/1",
			expectedRequestBody: `
			{
				"accessKeyName": "key2"
			}`,
			responseStatus: http.StatusOK,
			responseBody: `
			{
  				 "accessKeyName": "key2",
                 "AccessKeyUID": 1

			}`,
			expectedResponse: &UpdateAccessKeyResponse{
				AccessKeyUID:  1,
				AccessKeyName: "key2",
			},
		},
		"missing required params - validation error": {
			params: AccessKeyRequest{},
			accessKey: UpdateAccessKeyRequest{
				AccessKeyName: "key2",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update an access key: struct validation: AccessKeyUID: cannot be blank", err.Error())
			},
		},
		"missing required request body - validation error": {
			params: AccessKeyRequest{
				AccessKeyUID: 1,
			},
			accessKey: UpdateAccessKeyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update an access key: struct validation: AccessKeyName: cannot be blank", err.Error())
			},
		},
		"max length - validation error": {
			params: AccessKeyRequest{
				AccessKeyUID: 1,
			},
			accessKey: UpdateAccessKeyRequest{
				AccessKeyName: "asdfghjkloasdfghjkloasdfghjkloasdfghjkloasdfghjkloasdfghjkloasdfghjkloasdfghjklo",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update an access key: struct validation: AccessKeyName: the length must be between 1 and 50", err.Error())
			},
		},
		"500 internal server error": {
			accessKey: UpdateAccessKeyRequest{
				AccessKeyName: "key2",
			},
			params: AccessKeyRequest{
				AccessKeyUID: 1,
			},
			expectedPath: "/cam/v1/access-keys/1",
			expectedRequestBody: `
			{
				"accessKeyName": "key2"
			}`,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal-server-error",
				"title": "Internal Server Error",
				"detail": "Error processing request",
				"instance": "TestInstances",
				"status": 500
			}`,
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
				assert.Equal(t, http.MethodPut, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAccessKey(context.Background(), test.accessKey, test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
