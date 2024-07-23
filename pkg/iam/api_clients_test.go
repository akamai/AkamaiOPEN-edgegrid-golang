package iam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/internal/test"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestIAM_LockAPIClient(t *testing.T) {
	tests := map[string]struct {
		params           LockAPIClientRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *LockAPIClientResponse
		withError        func(*testing.T, error)
	}{
		"200 OK with specified client": {
			params: LockAPIClientRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/lock",
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessToken": "test_token1234",
  "activeCredentialCount": 1,
  "allowAccountSwitch": false,
  "authorizedUsers": [
    "jdoe"
  ],
  "clientDescription": "Test",
  "clientId": "abcd1234",
  "clientName": "test",
  "clientType": "CLIENT",
  "createdBy": "jdoe",
  "createdDate": "2022-05-13T20:04:35.000Z",
  "isLocked": true,
  "notificationEmails": [
    "jdoe@example.com"
  ],
  "serviceConsumerToken": "test_token12345"
}`,
			expectedResponse: &LockAPIClientResponse{
				AccessToken:             "test_token1234",
				ActiveCredentialCount:   1,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"jdoe"},
				CanAutoCreateCredential: false,
				ClientDescription:       "Test",
				ClientID:                "abcd1234",
				ClientName:              "test",
				ClientType:              "CLIENT",
				CreatedBy:               "jdoe",
				CreatedDate:             test.NewTimeFromString(t, "2022-05-13T20:04:35.000Z"),
				IsLocked:                true,
				NotificationEmails:      []string{"jdoe@example.com"},
				ServiceConsumerToken:    "test_token12345",
			},
		},
		"200 OK - self": {
			params:         LockAPIClientRequest{},
			expectedPath:   "/identity-management/v3/api-clients/self/lock",
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessToken": "test_token1234",
  "activeCredentialCount": 1,
  "allowAccountSwitch": false,
  "authorizedUsers": [
    "jdoe"
  ],
  "clientDescription": "Test",
  "clientId": "abcd1234",
  "clientName": "test",
  "clientType": "CLIENT",
  "createdBy": "jdoe",
  "createdDate": "2022-05-13T20:04:35.000Z",
  "isLocked": true,
  "notificationEmails": [
    "jdoe@example.com"
  ],
  "serviceConsumerToken": "test_token12345"
}`,
			expectedResponse: &LockAPIClientResponse{
				AccessToken:             "test_token1234",
				ActiveCredentialCount:   1,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"jdoe"},
				CanAutoCreateCredential: false,
				ClientDescription:       "Test",
				ClientID:                "abcd1234",
				ClientName:              "test",
				ClientType:              "CLIENT",
				CreatedBy:               "jdoe",
				CreatedDate:             test.NewTimeFromString(t, "2022-05-13T20:04:35.000Z"),
				IsLocked:                true,
				NotificationEmails:      []string{"jdoe@example.com"},
				ServiceConsumerToken:    "test_token12345",
			},
		},
		"404 Not Found": {
			params: LockAPIClientRequest{
				ClientID: "test12344",
			},
			expectedPath:   "/identity-management/v3/api-clients/test12344/lock",
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"instance": "",
				"httpStatus": 404,
				"detail": "",
				"title": "invalid open identity",
				"type": "/identity-management/error-types/2"
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusNotFound,
					Detail:     "",
					Title:      "invalid open identity",
					Type:       "/identity-management/error-types/2",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: LockAPIClientRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/lock",
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
			response, err := client.LockAPIClient(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, response)
		})
	}
}

func TestIAM_UnlockAPIClient(t *testing.T) {
	tests := map[string]struct {
		params           UnlockAPIClientRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *UnlockAPIClientResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UnlockAPIClientRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/unlock",
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessToken": "test_token1234",
  "activeCredentialCount": 1,
  "allowAccountSwitch": false,
  "authorizedUsers": [
    "jdoe"
  ],
  "clientDescription": "Test",
  "clientId": "abcd1234",
  "clientName": "test",
  "clientType": "CLIENT",
  "createdBy": "jdoe",
  "createdDate": "2022-05-13T20:04:35.000Z",
  "isLocked": true,
  "notificationEmails": [
    "jdoe@example.com"
  ],
  "serviceConsumerToken": "test_token12345"
}`,
			expectedResponse: &UnlockAPIClientResponse{
				AccessToken:             "test_token1234",
				ActiveCredentialCount:   1,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"jdoe"},
				CanAutoCreateCredential: false,
				ClientDescription:       "Test",
				ClientID:                "abcd1234",
				ClientName:              "test",
				ClientType:              "CLIENT",
				CreatedBy:               "jdoe",
				CreatedDate:             test.NewTimeFromString(t, "2022-05-13T20:04:35.000Z"),
				IsLocked:                true,
				NotificationEmails:      []string{"jdoe@example.com"},
				ServiceConsumerToken:    "test_token12345",
			},
		},
		"validation errors": {
			params: UnlockAPIClientRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "unlock api client: struct validation:\nClientID: cannot be blank", err.Error())
			},
		},
		"404 Not Found": {
			params: UnlockAPIClientRequest{
				ClientID: "test12344",
			},
			expectedPath:   "/identity-management/v3/api-clients/test12344/unlock",
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"instance": "",
				"httpStatus": 404,
				"detail": "",
				"title": "invalid open identity",
				"type": "/identity-management/error-types/2"
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusNotFound,
					Detail:     "",
					Title:      "invalid open identity",
					Type:       "/identity-management/error-types/2",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UnlockAPIClientRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/unlock",
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
			response, err := client.UnlockAPIClient(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, response)
		})
	}
}
