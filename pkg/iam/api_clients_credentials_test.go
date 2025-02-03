package iam

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAM_CreateCredential(t *testing.T) {
	tests := map[string]struct {
		params           CreateCredentialRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *CreateCredentialResponse
		withError        func(*testing.T, error)
	}{
		"201 Created with specified client": {
			params: CreateCredentialRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials",
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "credentialId": 123,
    "clientToken": "test-token",
    "clientSecret": "test-secret",
    "createdOn": "2024-07-25T11:02:28.000Z",
    "expiresOn": "2026-07-25T11:02:28.000Z",
    "status": "ACTIVE",
    "description": ""
}
`,
			expectedResponse: &CreateCredentialResponse{
				ClientSecret: "test-secret",
				ClientToken:  "test-token",
				CreatedOn:    test.NewTimeFromString(t, "2024-07-25T11:02:28.000Z"),
				CredentialID: 123,
				Description:  "",
				ExpiresOn:    test.NewTimeFromString(t, "2026-07-25T11:02:28.000Z"),
				Status:       CredentialActive,
			},
		},
		"200 OK - self": {
			params:         CreateCredentialRequest{},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials",
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "credentialId": 123,
    "clientToken": "test-token",
    "clientSecret": "test-secret",
    "createdOn": "2024-07-25T11:02:28.000Z",
    "expiresOn": "2026-07-25T11:02:28.000Z",
    "status": "ACTIVE",
    "description": ""
}
`,
			expectedResponse: &CreateCredentialResponse{
				ClientSecret: "test-secret",
				ClientToken:  "test-token",
				CreatedOn:    test.NewTimeFromString(t, "2024-07-25T11:02:28.000Z"),
				CredentialID: 123,
				Description:  "",
				ExpiresOn:    test.NewTimeFromString(t, "2026-07-25T11:02:28.000Z"),
				Status:       CredentialActive,
			},
		},
		"404 Not Found": {
			params: CreateCredentialRequest{
				ClientID: "test12344",
			},
			expectedPath:   "/identity-management/v3/api-clients/test12344/credentials",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/identity-management/error-types/2",
    "status": 404,
    "title": "invalid open identity",
    "detail": "",
    "instance": "",
    "errors": []
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:      "invalid open identity",
					Type:       "/identity-management/error-types/2",
					StatusCode: http.StatusNotFound,
					Errors:     json.RawMessage("[]"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params:         CreateCredentialRequest{},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			response, err := client.CreateCredential(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestIAM_ListCredentials(t *testing.T) {
	tests := map[string]struct {
		params           ListCredentialsRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse ListCredentialsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK with specified client": {
			params: ListCredentialsRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials?actions=false",
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "credentialId": 1,
        "clientToken": "test-token1",
        "status": "ACTIVE",
        "createdOn": "2024-05-14T11:10:25.000Z",
        "description": "",
        "expiresOn": "2026-05-14T11:10:25.000Z",
        "maxAllowedExpiry": "2026-07-25T11:09:30.658Z"
    },
    {
        "credentialId": 2,
        "clientToken": "test-token2",
        "status": "DELETED",
        "createdOn": "2024-05-28T06:53:36.000Z",
        "description": "deactivate for deletion",
        "expiresOn": "2025-10-11T23:06:59.000Z",
        "maxAllowedExpiry": "2026-07-25T11:09:30.658Z"
    },
    {
        "credentialId": 3,
        "clientToken": "test-token3",
        "status": "ACTIVE",
        "createdOn": "2024-07-25T11:02:28.000Z",
        "description": "",
        "expiresOn": "2026-07-25T11:02:28.000Z",
        "maxAllowedExpiry": "2026-07-25T11:09:30.658Z"
    }
]
`,
			expectedResponse: ListCredentialsResponse{
				{
					ClientToken:      "test-token1",
					CreatedOn:        test.NewTimeFromString(t, "2024-05-14T11:10:25.000Z"),
					CredentialID:     1,
					Description:      "",
					ExpiresOn:        test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
					Status:           CredentialActive,
					MaxAllowedExpiry: test.NewTimeFromString(t, "2026-07-25T11:09:30.658Z"),
				},
				{
					ClientToken:      "test-token2",
					CreatedOn:        test.NewTimeFromString(t, "2024-05-28T06:53:36.000Z"),
					CredentialID:     2,
					Description:      "deactivate for deletion",
					ExpiresOn:        test.NewTimeFromString(t, "2025-10-11T23:06:59.000Z"),
					Status:           CredentialDeleted,
					MaxAllowedExpiry: test.NewTimeFromString(t, "2026-07-25T11:09:30.658Z"),
				},
				{
					ClientToken:      "test-token3",
					CreatedOn:        test.NewTimeFromString(t, "2024-07-25T11:02:28.000Z"),
					CredentialID:     3,
					Description:      "",
					ExpiresOn:        test.NewTimeFromString(t, "2026-07-25T11:02:28.000Z"),
					Status:           CredentialActive,
					MaxAllowedExpiry: test.NewTimeFromString(t, "2026-07-25T11:09:30.658Z"),
				},
			},
		},
		"200 OK - self and actions query param": {
			params: ListCredentialsRequest{
				Actions: true,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials?actions=true",
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "credentialId": 1,
        "clientToken": "test-token1",
        "status": "ACTIVE",
        "createdOn": "2024-05-14T11:10:25.000Z",
        "description": "",
        "expiresOn": "2026-05-14T11:10:25.000Z",
        "maxAllowedExpiry": "2026-07-25T11:09:30.658Z",
		"actions": {
            "deactivate": true,
            "delete": true,
            "activate": true,
            "editDescription": true,
            "editExpiration": true
        }
    }
]
`,
			expectedResponse: ListCredentialsResponse{
				{
					ClientToken:      "test-token1",
					CreatedOn:        test.NewTimeFromString(t, "2024-05-14T11:10:25.000Z"),
					CredentialID:     1,
					Description:      "",
					ExpiresOn:        test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
					Status:           CredentialActive,
					MaxAllowedExpiry: test.NewTimeFromString(t, "2026-07-25T11:09:30.658Z"),
					Actions: &CredentialActions{
						Deactivate:      true,
						Delete:          true,
						Activate:        true,
						EditDescription: true,
						EditExpiration:  true,
					},
				},
			},
		},
		"404 Not Found": {
			params: ListCredentialsRequest{
				ClientID: "test12344",
			},
			expectedPath:   "/identity-management/v3/api-clients/test12344/credentials?actions=false",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/identity-management/error-types/2",
    "status": 404,
    "title": "invalid open identity",
    "detail": "",
    "instance": "",
    "errors": []
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:      "invalid open identity",
					Type:       "/identity-management/error-types/2",
					StatusCode: http.StatusNotFound,
					Errors:     json.RawMessage("[]"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params:         ListCredentialsRequest{},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials?actions=false",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			response, err := client.ListCredentials(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestIAM_GetCredential(t *testing.T) {
	tests := map[string]struct {
		params           GetCredentialRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *GetCredentialResponse
		withError        func(*testing.T, error)
	}{
		"200 OK with specified client": {
			params: GetCredentialRequest{
				ClientID:     "test1234",
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials/123?actions=false",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"credentialId": 1,
	"clientToken": "test-token1",
	"status": "ACTIVE",
	"createdOn": "2024-05-14T11:10:25.000Z",
	"description": "",
	"expiresOn": "2026-05-14T11:10:25.000Z",
	"maxAllowedExpiry": "2026-07-25T11:09:30.658Z"
}
`,
			expectedResponse: &GetCredentialResponse{
				ClientToken:      "test-token1",
				CreatedOn:        test.NewTimeFromString(t, "2024-05-14T11:10:25.000Z"),
				CredentialID:     1,
				Description:      "",
				ExpiresOn:        test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
				Status:           CredentialActive,
				MaxAllowedExpiry: test.NewTimeFromString(t, "2026-07-25T11:09:30.658Z"),
			},
		},
		"200 OK - self with actions query param": {
			params: GetCredentialRequest{
				CredentialID: 123,
				Actions:      true,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123?actions=true",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"credentialId": 1,
	"clientToken": "test-token1",
	"status": "ACTIVE",
	"createdOn": "2024-05-14T11:10:25.000Z",
	"description": "",
	"expiresOn": "2026-05-14T11:10:25.000Z",
	"maxAllowedExpiry": "2026-07-25T11:09:30.658Z",
	"actions": {
		"deactivate": true,
		"delete": true,
		"activate": true,	
		"editDescription": true,
		"editExpiration": false
	}
}
`,
			expectedResponse: &GetCredentialResponse{
				ClientToken:      "test-token1",
				CreatedOn:        test.NewTimeFromString(t, "2024-05-14T11:10:25.000Z"),
				CredentialID:     1,
				Description:      "",
				ExpiresOn:        test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
				Status:           CredentialActive,
				MaxAllowedExpiry: test.NewTimeFromString(t, "2026-07-25T11:09:30.658Z"),
				Actions: &CredentialActions{
					Deactivate:      true,
					Delete:          true,
					Activate:        true,
					EditDescription: true,
					EditExpiration:  false,
				},
			},
		},
		"validation errors": {
			params: GetCredentialRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get credential: struct validation: CredentialID: cannot be blank", err.Error())
			},
		},
		"404 Not Found": {
			params: GetCredentialRequest{
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123?actions=false",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/identity-management/error-types/25",
    "status": 404,
    "title": "ERROR_NO_CREDENTIAL",
    "detail": "",
    "instance": "",
    "errors": []
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:      "ERROR_NO_CREDENTIAL",
					Type:       "/identity-management/error-types/25",
					StatusCode: http.StatusNotFound,
					Errors:     json.RawMessage("[]"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: GetCredentialRequest{
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123?actions=false",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			response, err := client.GetCredential(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestIAM_UpdateCredential(t *testing.T) {
	tests := map[string]struct {
		params              UpdateCredentialRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *UpdateCredentialResponse
		withError           func(*testing.T, error)
	}{
		"200 OK with zeros as milliseconds - add nanosecond to the request": {
			params: UpdateCredentialRequest{
				ClientID:     "test1234",
				CredentialID: 123,
				Body: UpdateCredentialRequestBody{
					ExpiresOn: test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
					Status:    CredentialActive,
				},
			},
			expectedRequestBody: `
{
	"expiresOn": "2026-05-14T11:10:25.000000001Z",	
	"status": "ACTIVE"
}
`,
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials/123",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"status": "ACTIVE",
	"expiresOn": "2026-05-14T11:10:25.000Z"
}
`,
			expectedResponse: &UpdateCredentialResponse{
				ExpiresOn: test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
				Status:    CredentialActive,
			},
		},
		"200 OK with no milliseconds provided - add nanosecond to the request": {
			params: UpdateCredentialRequest{
				ClientID:     "test1234",
				CredentialID: 123,
				Body: UpdateCredentialRequestBody{
					ExpiresOn: test.NewTimeFromString(t, "2026-05-14T11:10:25Z"),
					Status:    CredentialActive,
				},
			},
			expectedRequestBody: `
{
	"expiresOn": "2026-05-14T11:10:25.000000001Z",	
	"status": "ACTIVE"
}
`,
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials/123",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"status": "ACTIVE",
	"expiresOn": "2026-05-14T11:10:25.000Z"
}
`,
			expectedResponse: &UpdateCredentialResponse{
				ExpiresOn: test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
				Status:    CredentialActive,
			},
		},
		"200 OK with specified client, without description": {
			params: UpdateCredentialRequest{
				ClientID:     "test1234",
				CredentialID: 123,
				Body: UpdateCredentialRequestBody{
					ExpiresOn: test.NewTimeFromString(t, "2026-05-14T11:10:25.123Z"),
					Status:    CredentialActive,
				},
			},
			expectedRequestBody: `
{
	"expiresOn": "2026-05-14T11:10:25.123Z",	
	"status": "ACTIVE"
}
`,
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials/123",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"status": "ACTIVE",
	"expiresOn": "2026-05-14T11:10:25.000Z"
}
`,
			expectedResponse: &UpdateCredentialResponse{
				ExpiresOn: test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
				Status:    CredentialActive,
			},
		},
		"200 OK without specified client, with description": {
			params: UpdateCredentialRequest{
				CredentialID: 123,
				Body: UpdateCredentialRequestBody{
					ExpiresOn:   test.NewTimeFromString(t, "2026-05-14T11:10:25.123Z"),
					Status:      CredentialInactive,
					Description: "test description",
				},
			},
			expectedRequestBody: `
{
	"description": "test description",
	"expiresOn": "2026-05-14T11:10:25.123Z",	
	"status": "INACTIVE"
}
`,
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"status": "INACTIVE",
	"expiresOn": "2026-05-14T11:10:25.000Z",
	"description": "test description"
}
`,
			expectedResponse: &UpdateCredentialResponse{
				ExpiresOn:   test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
				Status:      CredentialInactive,
				Description: ptr.To("test description"),
			},
		},
		"validation errors": {
			params: UpdateCredentialRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update credential: struct validation: Body: {\n\tExpiresOn: cannot be blank\n\tStatus: cannot be blank\n}\nCredentialID: cannot be blank", err.Error())
			},
		},
		"404 Not Found": {
			params: UpdateCredentialRequest{
				CredentialID: 123,
				Body: UpdateCredentialRequestBody{
					ExpiresOn: test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
					Status:    "ACTIVE",
				},
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123",
			responseStatus: http.StatusNotFound,
			responseBody: `
		{
		   "type": "/identity-management/error-types/25",
		   "status": 404,
		   "title": "ERROR_NO_CREDENTIAL",
		   "detail": "",
		   "instance": "",
		   "errors": []
		}
		`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:      "ERROR_NO_CREDENTIAL",
					Type:       "/identity-management/error-types/25",
					StatusCode: http.StatusNotFound,
					Errors:     json.RawMessage("[]"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UpdateCredentialRequest{
				CredentialID: 123,
				Body: UpdateCredentialRequestBody{
					ExpiresOn: test.NewTimeFromString(t, "2026-05-14T11:10:25.000Z"),
					Status:    CredentialActive,
				},
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
					{
						"type": "internal_error",
						"title": "Internal Server Error",
						"detail": "Error making request",
						"status": 500
					}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
				if tc.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			response, err := client.UpdateCredential(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestIAM_DeleteCredential(t *testing.T) {
	tests := map[string]struct {
		params         DeleteCredentialRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 with specified client": {
			params: DeleteCredentialRequest{
				ClientID:     "test1234",
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials/123",
			responseStatus: http.StatusNoContent,
		},
		"204 without specified client": {
			params: DeleteCredentialRequest{
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123",
			responseStatus: http.StatusNoContent,
		},
		"validation errors": {
			params: DeleteCredentialRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete credential: struct validation: CredentialID: cannot be blank", err.Error())
			},
		},
		"404 Not Found": {
			params: DeleteCredentialRequest{
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123",
			responseStatus: http.StatusNotFound,
			responseBody: `
		{
		   "type": "/identity-management/error-types/25",
		   "status": 404,
		   "title": "ERROR_NO_CREDENTIAL",
		   "detail": "",
		   "instance": "",
		   "errors": []
		}
		`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:      "ERROR_NO_CREDENTIAL",
					Type:       "/identity-management/error-types/25",
					StatusCode: http.StatusNotFound,
					Errors:     json.RawMessage("[]"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: DeleteCredentialRequest{
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
					{
						"type": "internal_error",
						"title": "Internal Server Error",
						"detail": "Error making request",
						"status": 500
					}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeleteCredential(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_DeactivateCredential(t *testing.T) {
	tests := map[string]struct {
		params         DeactivateCredentialRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 with specified client": {
			params: DeactivateCredentialRequest{
				ClientID:     "test1234",
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials/123/deactivate",
			responseStatus: http.StatusNoContent,
		},
		"204 without specified client": {
			params: DeactivateCredentialRequest{
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123/deactivate",
			responseStatus: http.StatusNoContent,
		},
		"validation errors": {
			params: DeactivateCredentialRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "deactivate credential: struct validation: CredentialID: cannot be blank", err.Error())
			},
		},
		"404 Not Found": {
			params: DeactivateCredentialRequest{
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123/deactivate",
			responseStatus: http.StatusNotFound,
			responseBody: `
		{
		   "type": "/identity-management/error-types/25",
		   "status": 404,
		   "title": "ERROR_NO_CREDENTIAL",
		   "detail": "",
		   "instance": "",
		   "errors": []
		}
		`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:      "ERROR_NO_CREDENTIAL",
					Type:       "/identity-management/error-types/25",
					StatusCode: http.StatusNotFound,
					Errors:     json.RawMessage("[]"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: DeactivateCredentialRequest{
				CredentialID: 123,
			},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/123/deactivate",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
					{
						"type": "internal_error",
						"title": "Internal Server Error",
						"detail": "Error making request",
						"status": 500
					}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeactivateCredential(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_DeactivateCredentials(t *testing.T) {
	tests := map[string]struct {
		params         DeactivateCredentialsRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 with specified client": {
			params: DeactivateCredentialsRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/credentials/deactivate",
			responseStatus: http.StatusNoContent,
		},
		"204 without specified client": {
			params:         DeactivateCredentialsRequest{},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/deactivate",
			responseStatus: http.StatusNoContent,
		},
		"404 Not Found": {
			params:         DeactivateCredentialsRequest{},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/deactivate",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/identity-management/error-types/2",
    "status": 404,
    "title": "invalid open identity",
    "detail": "",
    "instance": "",
    "errors": []
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Title:      "invalid open identity",
					Type:       "/identity-management/error-types/2",
					StatusCode: http.StatusNotFound,
					Errors:     json.RawMessage("[]"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params:         DeactivateCredentialsRequest{},
			expectedPath:   "/identity-management/v3/api-clients/self/credentials/deactivate",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
					{
						"type": "internal_error",
						"title": "Internal Server Error",
						"detail": "Error making request",
						"status": 500
					}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeactivateCredentials(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
