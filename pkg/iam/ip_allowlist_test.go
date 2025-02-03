package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAM_DisableIPAllowlist(t *testing.T) {
	tests := map[string]struct {
		responseStatus int
		expectedPath   string
		responseBody   string
		withError      func(*testing.T, error)
	}{
		"204 no content": {
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v3/user-admin/ip-acl/allowlist/disable",
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/disable",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				if tc.responseBody != "" {
					_, err := w.Write([]byte(tc.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DisableIPAllowlist(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_EnableIPAllowlist(t *testing.T) {
	tests := map[string]struct {
		responseStatus int
		expectedPath   string
		responseBody   string
		withError      func(*testing.T, error)
	}{
		"204 no content": {
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v3/user-admin/ip-acl/allowlist/enable",
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/enable",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				if tc.responseBody != "" {
					_, err := w.Write([]byte(tc.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.EnableIPAllowlist(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_GetIPAllowlistStatus(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetIPAllowlistStatusResponse
		withError        func(*testing.T, error)
	}{
		"200 OK enabled true": {
			responseStatus: 200,
			expectedPath:   "/identity-management/v3/user-admin/ip-acl/allowlist/status",
			responseBody: `
			{
				"enabled": true
			}`,
			expectedResponse: &GetIPAllowlistStatusResponse{
				Enabled: true,
			},
		},
		"200 OK enabled false": {
			responseStatus: 200,
			expectedPath:   "/identity-management/v3/user-admin/ip-acl/allowlist/status",
			responseBody: `
			{
				"enabled": false
			}`,
			expectedResponse: &GetIPAllowlistStatusResponse{
				Enabled: false,
			},
		},
		"500 internal server error": {
			responseStatus: 500,
			expectedPath:   "/identity-management/v3/user-admin/ip-acl/allowlist/status",
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"detail": "Error making request",
	"status": 500
}
`,
			withError: func(t *testing.T, e error) {
				err := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					StatusCode: 500,
					Detail:     "Error making request",
				}
				assert.Equal(t, true, err.Is(e))
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
			result, err := client.GetIPAllowlistStatus(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
