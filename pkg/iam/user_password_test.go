package iam

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestIAM_ResetUserPassword(t *testing.T) {
	tests := map[string]struct {
		params           ResetUserPasswordRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse ResetUserPasswordResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ResetUserPasswordRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"newPassword": "K8QVa7Q2"
			}`,
			expectedResponse: ResetUserPasswordResponse{
				NewPassword: "K8QVa7Q2",
			},
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/reset-password?sendEmail=false",
		},
		"204 No Content": {
			params: ResetUserPasswordRequest{
				IdentityID: "1-ABCDE",
				SendEmail:  true,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE/reset-password?sendEmail=true",
		},
		"404 Not Found": {
			params: ResetUserPasswordRequest{
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/X1-ABCDE/reset-password?sendEmail=false",
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
			params: ResetUserPasswordRequest{
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/reset-password?sendEmail=false",
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
			response, err := client.ResetUserPassword(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, *response)
		})
	}
}

func TestIAM_SetUserPassword(t *testing.T) {
	tests := map[string]struct {
		params              SetUserPasswordRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"204 No Content": {
			params: SetUserPasswordRequest{
				IdentityID:  "1-ABCDE",
				NewPassword: "newpwd",
			},
			responseStatus:      http.StatusNoContent,
			responseBody:        "",
			expectedRequestBody: `{"newPassword":"newpwd"}`,
			expectedPath:        "/identity-management/v3/user-admin/ui-identities/1-ABCDE/set-password",
		},
		"400 Bad Request - same password": {
			params: SetUserPasswordRequest{
				IdentityID:  "X1-ABCDE",
				NewPassword: "newpwd",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"instance": "",
				"httpStatus": 400,
				"detail": "Must not match a previously used password.",
				"title": "Validation Exception",
				"type": "/useradmin-api/error-types/1508"
			}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/X1-ABCDE/set-password",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusBadRequest,
					Detail:     "Must not match a previously used password.",
					Title:      "Validation Exception",
					Type:       "/useradmin-api/error-types/1508",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"404 Not Found": {
			params: SetUserPasswordRequest{
				IdentityID:  "X1-ABCDE",
				NewPassword: "newpwd",
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/X1-ABCDE/set-password",
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
			params: SetUserPasswordRequest{
				IdentityID:  "1-ABCDE",
				NewPassword: "newpwd",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
			}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/set-password",
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

				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.SetUserPassword(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
