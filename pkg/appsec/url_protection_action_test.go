package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test ListURLProtectionPoliciesActionsRequest Validate
func TestListURLProtectionPoliciesActionsRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           ListURLProtectionPoliciesActionsRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: ListURLProtectionPoliciesActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: ListURLProtectionPoliciesActionsRequest{
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing Version": {
			req: ListURLProtectionPoliciesActionsRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing PolicyID": {
			req: ListURLProtectionPoliciesActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing all required fields": {
			req:           ListURLProtectionPoliciesActionsRequest{},
			errorExpected: true,
		},
		"missing ConfigID and Version": {
			req: ListURLProtectionPoliciesActionsRequest{
				PolicyID: "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing ConfigID and PolicyID": {
			req: ListURLProtectionPoliciesActionsRequest{
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing Version and PolicyID": {
			req: ListURLProtectionPoliciesActionsRequest{
				ConfigID: 43253,
			},
			errorExpected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.errorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// Test GetURLProtectionPolicyActionsRequest Validate
func TestGetURLProtectionPolicyActionsRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           GetURLProtectionPolicyActionsRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing Version": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing PolicyID": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing URLProtectionPolicyID": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing all required fields": {
			req:           GetURLProtectionPolicyActionsRequest{},
			errorExpected: true,
		},
		"missing ConfigID and Version": {
			req: GetURLProtectionPolicyActionsRequest{
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigID and PolicyID": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigID and URLProtectionPolicyID": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing Version and PolicyID": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing Version and URLProtectionPolicyID": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing PolicyID and URLProtectionPolicyID": {
			req: GetURLProtectionPolicyActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.errorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// Test UpdateURLProtectionPolicyActionsRequest Validate
func TestUpdateURLProtectionPolicyActionsRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           UpdateURLProtectionPolicyActionsRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyActions{
					MaxRateThresholdAction: "alert",
					LoadSheddingAction:     "deny",
				},
			},
			errorExpected: false,
		},
		"valid request - missing MaxRateThresholdAction": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyActions{
					LoadSheddingAction: "deny",
				},
			},
			errorExpected: true,
		},
		"missing ConfigID": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing Version": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing PolicyID": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing URLProtectionPolicyID": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing all required fields": {
			req:           UpdateURLProtectionPolicyActionsRequest{},
			errorExpected: true,
		},
		"missing ConfigID and Version": {
			req: UpdateURLProtectionPolicyActionsRequest{
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigID and PolicyID": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigID and URLProtectionPolicyID": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing Version and PolicyID": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing Version and URLProtectionPolicyID": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing PolicyID and URLProtectionPolicyID": {
			req: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.errorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// Test ListURLProtectionPoliciesActions.
func TestAppSec_ListURLProtectionPoliciesActions(t *testing.T) {
	result := ListURLProtectionPoliciesActionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestURLProtectionActions/URLProtectionPoliciesActions.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           ListURLProtectionPoliciesActionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListURLProtectionPoliciesActionsResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: ListURLProtectionPoliciesActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
		},
		"400 bad request": {
			params: ListURLProtectionPoliciesActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: ListURLProtectionPoliciesActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: ListURLProtectionPoliciesActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: ListURLProtectionPoliciesActionsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				PolicyID:      "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListURLProtectionPoliciesActions(
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

// Test GetURLProtectionPolicyActions.
func TestAppSec_GetURLProtectionPolicyActions(t *testing.T) {
	result := GetURLProtectionPolicyActionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestURLProtectionActions/URLProtectionPolicyActions.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetURLProtectionPolicyActionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetURLProtectionPolicyActionsResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     compactJSON(loadFixtureBytes("testdata/TestURLProtectionActions/URLProtectionPoliciesActions.json")),
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
		},
		"400 bad request": {
			params: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 123,
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 234,
			},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 234,
			},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: GetURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 234,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetURLProtectionPolicyActions(
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

// Test UpdateURLProtectionPoliciesAction.
func TestAppSec_UpdateURLProtectionPoliciesActions(t *testing.T) {
	result := UpdateURLProtectionPolicyActionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestURLProtectionActions/URLProtectionPolicyActions.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateURLProtectionPolicyActionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateURLProtectionPolicyActionsResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyActions{
					MaxRateThresholdAction: "alert",
					LoadSheddingAction:     "deny",
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections/134644",
		},
		"400 bad request": {
			params: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyActions{
					MaxRateThresholdAction: "alert",
					LoadSheddingAction:     "deny",
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyActions{
					MaxRateThresholdAction: "alert",
					LoadSheddingAction:     "deny",
				},
			},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyActions{
					MaxRateThresholdAction: "alert",
					LoadSheddingAction:     "deny",
				},
			},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: UpdateURLProtectionPolicyActionsRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				PolicyID:              "AAAA_81230",
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyActions{
					MaxRateThresholdAction: "alert",
					LoadSheddingAction:     "deny",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/url-protections/134644",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateURLProtectionPolicyActions(
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
