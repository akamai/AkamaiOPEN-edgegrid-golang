package cloudaccess

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestCreateAccessKeyVersion(t *testing.T) {
	tests := map[string]struct {
		params              CreateAccessKeyVersionRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *CreateAccessKeyVersionResponse
		withError           func(*testing.T, error)
	}{
		"202 ACCEPTED": {
			params: CreateAccessKeyVersionRequest{
				AccessKeyUID: 1,
				BodyParams: CreateAccessKeyVersionBodyParams{
					CloudAccessKeyID:     "key-1",
					CloudSecretAccessKey: "secret-1",
				},
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
  "requestId": 111,
  "retryAfter": 6
}`,
			expectedPath: "/cam/v1/access-keys/1/versions",
			expectedRequestBody: `
{
	"cloudAccessKeyId": "key-1", 
	"cloudSecretAccessKey": "secret-1"
}
`,
			expectedResponse: &CreateAccessKeyVersionResponse{
				RequestID:  111,
				RetryAfter: 6,
			},
		},
		"missing required params - validation error": {
			params: CreateAccessKeyVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create access key version: struct validation: AccessKeyUID: cannot be blank\nBodyParams: CloudAccessKeyID: cannot be blank\nCloudSecretAccessKey: cannot be blank", err.Error())
			},
		},
		"404 error": {
			params: CreateAccessKeyVersionRequest{
				AccessKeyUID: 1,
				BodyParams: CreateAccessKeyVersionBodyParams{
					CloudAccessKeyID:     "key-1",
					CloudSecretAccessKey: "secret-1",
				},
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
  "accessKeyUid": 1,
  "detail": "Access key with accessKeyUid '1' does not exist.",
  "instance": "c111eff1-22ec-4d4e-99c9-55efb5d55f55",
  "status": 404,
  "title": "Domain Error",
  "type": "/cam/error-types/access-key-does-not-exist"
}`,
			expectedPath: "/cam/v1/access-keys/1/versions",
			expectedRequestBody: `
{
	"cloudAccessKeyId": "key-1", 
	"cloudSecretAccessKey": "secret-1"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					AccessKeyUID: 1,
					Type:         "/cam/error-types/access-key-does-not-exist",
					Title:        "Domain Error",
					Detail:       "Access key with accessKeyUid '1' does not exist.",
					Instance:     "c111eff1-22ec-4d4e-99c9-55efb5d55f55",
					Status:       http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"409 error": {
			params: CreateAccessKeyVersionRequest{
				AccessKeyUID: 1,
				BodyParams: CreateAccessKeyVersionBodyParams{
					CloudAccessKeyID:     "key-1",
					CloudSecretAccessKey: "secret-1",
				},
			},
			responseStatus: http.StatusConflict,
			responseBody: `
{
  "accessKeyName": "Sales-s3",
  "detail": "Access key with name 'Sales-s3' already exists.",
  "instance": "109443e6-f347-43f1-922c-fa0fd480973f",
  "status": 409,
  "title": "Domain Error",
  "type": "/cam/error-types/access-key-already-exists"
}`,
			expectedPath: "/cam/v1/access-keys/1/versions",
			expectedRequestBody: `
{
	"cloudAccessKeyId": "key-1", 
	"cloudSecretAccessKey": "secret-1"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					AccessKeyName: "Sales-s3",
					Type:          "/cam/error-types/access-key-already-exists",
					Title:         "Domain Error",
					Detail:        "Access key with name 'Sales-s3' already exists.",
					Instance:      "109443e6-f347-43f1-922c-fa0fd480973f",
					Status:        http.StatusConflict,
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
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateAccessKeyVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetAccessKeyVersion(t *testing.T) {
	tests := map[string]struct {
		params           GetAccessKeyVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAccessKeyVersionResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetAccessKeyVersionRequest{
				AccessKeyUID: 12345,
				Version:      1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessKeyUid": 12345,
  "cloudAccessKeyId": null,
  "createdBy": "testUser",
  "createdTime": "2021-02-26T13:34:37.916873Z",
  "deploymentStatus": "ACTIVE",
  "version": 1,
  "versionGuid": "aaaa-bbbb-1111"
}`,
			expectedPath: "/cam/v1/access-keys/12345/versions/1",
			expectedResponse: &GetAccessKeyVersionResponse{
				AccessKeyUID:     12345,
				CreatedBy:        "testUser",
				CreatedTime:      *newTimeFromString(t, "2021-02-26T13:34:37.916873Z"),
				DeploymentStatus: Active,
				Version:          1,
				VersionGUID:      "aaaa-bbbb-1111",
			},
		},
		"missing required params - validation error": {
			params: GetAccessKeyVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get access key version: struct validation: AccessKeyUID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"404 error": {
			params: GetAccessKeyVersionRequest{
				AccessKeyUID: 1,
				Version:      1,
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
		{
		 "accessKeyUid": 1,
		 "detail": "Access key with accessKeyUid '1' does not exist.",
		 "instance": "c111eff1-22ec-4d4e-99c9-55efb5d55f55",
		 "status": 404,
		 "title": "Domain Error",
		 "type": "/cam/error-types/access-key-does-not-exist"
		}`,
			expectedPath: "/cam/v1/access-keys/1/versions/1",
			withError: func(t *testing.T, err error) {
				want := &Error{
					AccessKeyUID: 1,
					Type:         "/cam/error-types/access-key-does-not-exist",
					Title:        "Domain Error",
					Detail:       "Access key with accessKeyUid '1' does not exist.",
					Instance:     "c111eff1-22ec-4d4e-99c9-55efb5d55f55",
					Status:       http.StatusNotFound,
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
			result, err := client.GetAccessKeyVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListAccessKeyVersions(t *testing.T) {
	tests := map[string]struct {
		params           ListAccessKeyVersionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListAccessKeyVersionsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListAccessKeyVersionsRequest{
				AccessKeyUID: 2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessKeyVersions": [
    {
      "accessKeyUid": 2,
      "cloudAccessKeyId": null,
      "createdBy": "testUser2",
      "createdTime": "2021-02-26T14:48:27.355346Z",
      "deploymentStatus": "PENDING_ACTIVATION",
      "version": 2,
      "versionGuid": "bbbb-2222"
    },
    {
      "accessKeyUid": 2,
      "cloudAccessKeyId": null,
      "createdBy": "testUser1",
      "createdTime": "2021-02-26T13:34:37.916873Z",
      "deploymentStatus": "ACTIVE",
      "version": 1,
      "versionGuid": "aaaa-1111"
    }
  ]
}`,
			expectedPath: "/cam/v1/access-keys/2/versions",
			expectedResponse: &ListAccessKeyVersionsResponse{
				AccessKeyVersions: []AccessKeyVersion{
					{
						AccessKeyUID:     2,
						CreatedBy:        "testUser2",
						CreatedTime:      *newTimeFromString(t, "2021-02-26T14:48:27.355346Z"),
						DeploymentStatus: PendingActivation,
						Version:          2,
						VersionGUID:      "bbbb-2222",
					},
					{
						AccessKeyUID:     2,
						CreatedBy:        "testUser1",
						CreatedTime:      *newTimeFromString(t, "2021-02-26T13:34:37.916873Z"),
						DeploymentStatus: Active,
						Version:          1,
						VersionGUID:      "aaaa-1111",
					},
				},
			},
		},
		"200 OK - single version": {
			params: ListAccessKeyVersionsRequest{
				AccessKeyUID: 2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessKeyVersions": [
    {
      "accessKeyUid": 2,
      "cloudAccessKeyId": null,
      "createdBy": "testUser2",
      "createdTime": "2021-02-26T14:48:27.355346Z",
      "deploymentStatus": "PENDING_ACTIVATION",
      "version": 2,
      "versionGuid": "bbbb-2222"
    }
  ]
}`,
			expectedPath: "/cam/v1/access-keys/2/versions",
			expectedResponse: &ListAccessKeyVersionsResponse{
				AccessKeyVersions: []AccessKeyVersion{
					{
						AccessKeyUID:     2,
						CreatedBy:        "testUser2",
						CreatedTime:      *newTimeFromString(t, "2021-02-26T14:48:27.355346Z"),
						DeploymentStatus: PendingActivation,
						Version:          2,
						VersionGUID:      "bbbb-2222",
					},
				},
			},
		},
		"200 OK - no versions": {
			params: ListAccessKeyVersionsRequest{
				AccessKeyUID: 2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessKeyVersions": []
}`,
			expectedPath: "/cam/v1/access-keys/2/versions",
			expectedResponse: &ListAccessKeyVersionsResponse{
				AccessKeyVersions: []AccessKeyVersion{},
			},
		},
		"missing required params - validation error": {
			params: ListAccessKeyVersionsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list access key versions: struct validation: AccessKeyUID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: ListAccessKeyVersionsRequest{
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
			expectedPath: "/cam/v1/access-keys/1/versions",
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
			result, err := client.ListAccessKeyVersions(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteAccessKeyVersion(t *testing.T) {
	tests := map[string]struct {
		params           DeleteAccessKeyVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DeleteAccessKeyVersionResponse
		withError        func(*testing.T, error)
	}{
		"202 ACCEPTED": {
			params: DeleteAccessKeyVersionRequest{
				AccessKeyUID: 12345,
				Version:      1,
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
  "accessKeyUid": 12345,
  "cloudAccessKeyId": null,
  "createdBy": "testUser",
  "createdTime": "2021-02-26T09:09:53.762230Z",
  "deploymentStatus": "PENDING_DELETION",
  "version": 1,
  "versionGuid": "aaaa-bbbb-1111"
}`,
			expectedPath: "/cam/v1/access-keys/12345/versions/1",
			expectedResponse: &DeleteAccessKeyVersionResponse{
				AccessKeyUID:     12345,
				CreatedBy:        "testUser",
				CreatedTime:      *newTimeFromString(t, "2021-02-26T09:09:53.762230Z"),
				DeploymentStatus: PendingDeletion,
				Version:          1,
				VersionGUID:      "aaaa-bbbb-1111",
			},
		},
		"missing required params - validation error": {
			params: DeleteAccessKeyVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete access key version: struct validation: AccessKeyUID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"404 error": {
			params: DeleteAccessKeyVersionRequest{
				AccessKeyUID: 1,
				Version:      1,
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
		{
		 "accessKeyUid": 1,
		 "detail": "Access key with accessKeyUid '1' does not exist.",
		 "instance": "c111eff1-22ec-4d4e-99c9-55efb5d55f55",
		 "status": 404,
		 "title": "Domain Error",
		 "type": "/cam/error-types/access-key-does-not-exist"
		}`,
			expectedPath: "/cam/v1/access-keys/1/versions/1",
			withError: func(t *testing.T, err error) {
				want := &Error{
					AccessKeyUID: 1,
					Type:         "/cam/error-types/access-key-does-not-exist",
					Title:        "Domain Error",
					Detail:       "Access key with accessKeyUid '1' does not exist.",
					Instance:     "c111eff1-22ec-4d4e-99c9-55efb5d55f55",
					Status:       http.StatusNotFound,
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
			result, err := client.DeleteAccessKeyVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func newTimeFromString(t *testing.T, s string) *time.Time {
	parsedTime, err := time.Parse(time.RFC3339, s)
	require.NoError(t, err)
	return &parsedTime
}
