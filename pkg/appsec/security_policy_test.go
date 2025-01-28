package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListSecurityPolicies(t *testing.T) {

	result := GetSecurityPoliciesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSecurityPolicy/SecurityPolicy.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetSecurityPoliciesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSecurityPoliciesResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetSecurityPoliciesRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetSecurityPoliciesRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies",
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
			result, err := client.GetSecurityPolicies(
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

func TestAppSec_CreateSecurityPolicyWithDefaultProtectionsSuccess(t *testing.T) {

	response := CreateSecurityPolicyResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestSecurityPolicy/CreateSecurityPolicyWithDefaultProtectionsResponse.json"))
	err := json.Unmarshal([]byte(respData), &response)
	require.NoError(t, err)

	requestParams := CreateSecurityPolicyWithDefaultProtectionsRequest{
		ConfigVersion: ConfigVersion{
			ConfigID: 43253,
			Version:  15,
		},
		PolicyName:   "akamaitools",
		PolicyPrefix: "AK01",
	}

	tests := map[string]struct {
		params           CreateSecurityPolicyWithDefaultProtectionsRequest
		responseStatus   int
		responseBody     string
		expectedMethod   string
		expectedPath     string
		expectedResponse *CreateSecurityPolicyResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: requestParams,
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedMethod:   http.MethodPost,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/protections",
			expectedResponse: &response,
		},
		"500 internal server error": {
			params:         requestParams,
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating security policy with default protections",
    "status": 500
}`,
			expectedMethod: http.MethodPost,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/protections",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating security policy with default protections",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, test.expectedMethod, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			_, err := client.CreateSecurityPolicyWithDefaultProtections(
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
		})
	}
}
