package iam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAM_LockUser(t *testing.T) {
	tests := map[string]struct {
		params         LockUserRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			params: LockUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusOK,
			responseBody:   "",
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE/lock",
		},
		"204 No Content": {
			params: LockUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE/lock",
		},
		"404 Not Found": {
			params: LockUserRequest{
				IdentityID: "X1-ABCDE",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"instance": "",
				"httpStatus": 404,
				"detail": "",
				"title": "User not found",
				"type": "/useradmin-api/error-types/1100"
			}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/X1-ABCDE/lock",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusNotFound,
					Detail:     "",
					Title:      "User not found",
					Type:       "/useradmin-api/error-types/1100",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: LockUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
			}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/lock",
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
			err := client.LockUser(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_UnlockUser(t *testing.T) {
	tests := map[string]struct {
		params         UnlockUserRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			params: UnlockUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusOK,
			responseBody:   "",
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE/unlock",
		},
		"204 No Content": {
			params: UnlockUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE/unlock",
		},
		"404 Not Found": {
			params: UnlockUserRequest{
				IdentityID: "X1-ABCDE",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"instance": "",
				"httpStatus": 404,
				"detail": "",
				"title": "User not found",
				"type": "/useradmin-api/error-types/1100"
			}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/X1-ABCDE/unlock",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusNotFound,
					Detail:     "",
					Title:      "User not found",
					Type:       "/useradmin-api/error-types/1100",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UnlockUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
			}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/unlock",
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
			err := client.UnlockUser(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
