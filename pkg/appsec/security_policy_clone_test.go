package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListSecurityPolicyClone(t *testing.T) {

	result := GetSecurityPolicyCloneResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSecurityPolicyClone/SecurityPolicyClone.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetSecurityPolicyCloneRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSecurityPolicyCloneResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetSecurityPolicyCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetSecurityPolicyCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching propertys",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching propertys",
				StatusCode: http.StatusInternalServerError,
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
			result, err := client.GetSecurityPolicyClone(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test SecurityPolicyClone
func TestAppSec_GetSecurityPolicyClone(t *testing.T) {

	result := GetSecurityPolicyCloneResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSecurityPolicyClone/SecurityPolicyClone.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetSecurityPolicyCloneRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSecurityPolicyCloneResponse
		withError        error
	}{
		"200 OK": {
			params: GetSecurityPolicyCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetSecurityPolicyCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching match target",
				StatusCode: http.StatusInternalServerError,
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
			result, err := client.GetSecurityPolicyClone(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create SecurityPolicyClone
// Test Create SecurityPolicyClone
func TestAppSec_CreateSecurityPolicyClone(t *testing.T) {

	result := CreateSecurityPolicyCloneResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSecurityPolicyClone/SecurityPolicyClone.json"))
	json.Unmarshal([]byte(respData), &result)

	req := CreateSecurityPolicyCloneRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestSecurityPolicyClone/SecurityPolicyClone.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           CreateSecurityPolicyCloneRequest
		prop             *CreateSecurityPolicyCloneRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateSecurityPolicyCloneResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateSecurityPolicyCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/",
		},
		"500 internal server error": {
			params: CreateSecurityPolicyCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating domain"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating domain",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateSecurityPolicyClone(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
