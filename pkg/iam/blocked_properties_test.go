package iam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIam_ListBlockedProperties(t *testing.T) {
	tests := map[string]struct {
		params           ListBlockedPropertiesRequest
		responseStatus   int
		expectedPath     string
		responseBody     string
		expectedResponse []int64
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListBlockedPropertiesRequest{
				GroupID:    12345,
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/groups/12345/blocked-properties",
			responseBody: `[
								10977166
							]`,
			expectedResponse: []int64{
				10977166,
			},
		},
		"200 OK, no blocked property": {
			params: ListBlockedPropertiesRequest{
				GroupID:    12345,
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/groups/12345/blocked-properties",
			responseBody: `[
							
							]`,
			expectedResponse: []int64{},
		},

		"404 not found error": {
			params: ListBlockedPropertiesRequest{
				GroupID:    123450000,
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/groups/123450000/blocked-properties",
			responseBody: `
			{
    "instance": "",
    "httpStatus": 404,
    "detail": "group not found",
    "title": "Not found",
    "type": "/useradmin-api/error-types/1700"
}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusNotFound,
					Detail:     "group not found",
					Title:      "Not found",
					Type:       "/useradmin-api/error-types/1700",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},

		"500 internal server error": {
			params: ListBlockedPropertiesRequest{
				GroupID:    12345,
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/groups/12345/blocked-properties",
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error processing request",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error processing request",
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
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			users, err := client.ListBlockedProperties(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}

func TestIam_UpdateBlockedProperties(t *testing.T) {
	tests := map[string]struct {
		params           UpdateBlockedPropertiesRequest
		responseStatus   int
		expectedPath     string
		responseBody     string
		expectedResponse []int64
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateBlockedPropertiesRequest{
				GroupID:    12345,
				IdentityID: "1-ABCDE",
				Properties: []int64{10977166, 10977167},
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/groups/12345/blocked-properties",
			responseBody: `[
								10977166,10977167
							]`,
			expectedResponse: []int64{
				10977166, 10977167,
			},
		},
		"400 bad request": {
			params: UpdateBlockedPropertiesRequest{
				GroupID:    12345,
				IdentityID: "1-ABCDE",
				Properties: []int64{0, 1},
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/groups/12345/blocked-properties",
			responseBody: `
			{
    "instance": "",
    "httpStatus": 400,
    "detail": "",
    "title": "Validation Exception",
    "type": "/useradmin-api/error-types/1003"
}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusBadRequest,
					Detail:     "",
					Title:      "Validation Exception",
					Type:       "/useradmin-api/error-types/1003",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},

		"404 not found error": {
			params: UpdateBlockedPropertiesRequest{
				GroupID:    123450000,
				IdentityID: "1-ABCDE",
				Properties: []int64{10977166, 10977167},
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/groups/123450000/blocked-properties",
			responseBody: `
			{
    "instance": "",
    "httpStatus": 404,
    "detail": "group not found",
    "title": "Not found",
    "type": "/useradmin-api/error-types/1700"
}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusNotFound,
					Detail:     "group not found",
					Title:      "Not found",
					Type:       "/useradmin-api/error-types/1700",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},

		"500 internal server error": {
			params: UpdateBlockedPropertiesRequest{
				GroupID:    12345,
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/groups/12345/blocked-properties",
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error processing request",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error processing request",
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
			users, err := client.UpdateBlockedProperties(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}
