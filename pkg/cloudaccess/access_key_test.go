package cloudaccess

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccessKeyStatus(t *testing.T) {

	var result GetAccessKeyStatusResponse
	var resultMinimal GetAccessKeyStatusResponse

	respData, err := loadTestData("AccessKeyStatus/GetAccessKeyStatus.resp.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	respDataMinimal, err := loadTestData("AccessKeyStatus/GetAccessKeyStatus.resp.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(bytes.NewBuffer(respDataMinimal)).Decode(&resultMinimal); err != nil {
		t.Fatal(err)
	}

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
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/cam/v1/access-key-create-requests/1",
			expectedResponse: &result,
		},
		"200 OK - minimal": {
			params: GetAccessKeyStatusRequest{
				RequestID: 1,
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respDataMinimal),
			expectedPath:     "/cam/v1/access-key-create-requests/1",
			expectedResponse: &resultMinimal,
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

func TestCreateAccessKey(t *testing.T) {

	var req CreateAccessKeyRequest
	reqData, err := loadTestData("AccessKey/CreateAccessKey.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		accessKey        CreateAccessKeyRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *CreateAccessKeyResponse
		responseHeaders  map[string]string
		withError        func(*testing.T, error)
	}{
		"202 Accepted": {
			accessKey:      req,
			expectedPath:   "/cam/v1/access-keys",
			responseStatus: http.StatusAccepted,
			responseBody: `
			{
  				"requestId": 195,
  				"retryAfter": 4

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
		"missing required request body - validation error": {
			accessKey: CreateAccessKeyRequest{
				Credentials:          Credentials{},
				NetworkConfiguration: SecureNetwork{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create an access key: struct validation: AccessKeyName: cannot be blank\nAuthenticationMethod: cannot be blank\nCloudAccessKeyID: cannot be blank\nCloudSecretAccessKey: cannot be blank\nContractID: cannot be blank\nGroupID: cannot be blank\nSecurityNetwork: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			accessKey:      req,
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
				if len(test.responseHeaders) > 0 {
					for header, value := range test.responseHeaders {
						w.Header().Set(header, value)
					}
				}

				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
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

	var result GetAccessKeyResponse

	respData, err := loadTestData("AccessKey/GetAccessKey.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

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
			expectedPath:     "/cam/v1/access-keys/1",
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedResponse: &result,
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
			responseBody: `{
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

	var result ListAccessKeysResponse

	respData, err := loadTestData("AccessKey/ListAccessKey.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

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
			expectedPath:     "/cam/v1/access-keys?versionGuid=1",
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedResponse: &result,
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
		accessKey        UpdateAccessKeyRequest
		params           AccessKeyRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *UpdateAccessKeyResponse
		withError        func(*testing.T, error)
	}{
		"201 OK": {
			accessKey: UpdateAccessKeyRequest{
				AccessKeyName: "key2",
			},
			params: AccessKeyRequest{
				AccessKeyUID: 1,
			},
			expectedPath:   "/cam/v1/access-keys/1",
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
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
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

func loadTestData(name string) ([]byte, error) {
	data, err := os.ReadFile(fmt.Sprintf("./testdata/%s", name))
	if err != nil {
		return nil, err
	}

	return data, nil
}
