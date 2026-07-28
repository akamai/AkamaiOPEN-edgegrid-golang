package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var badRequestResponse = `
{
    "type": "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
    "title": "Bad Request",
    "detail": "The request could not be understood by the server due to malformed syntax.",
	"status": 400
}`

var forbiddenResponse = `
{
  "detail": "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
  "status": 403,
  "title": "Forbidden",
  "type": "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED"
}`

var notFoundResponse = `
{
  "detail": "The requested resource is not found",
  "status": 404,
  "title": "Not Found",
  "type": "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND"
}`

var internalServerErrorResponse = `
{
  "detail": "Error fetching WAF composite ruleset",
  "status": 500,
  "title": "Internal Server Error",
  "type": "https://problems.luna.akamaiapis.net/appsec/error-types/INTERNAL-SERVER-ERROR"
}`

func TestAppSec_GetWAFCompositeRuleset(t *testing.T) {

	result := CompositeRulesetResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestWAFCompositeRuleset/WAFCompositeRuleset.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetWAFCompositeRulesetRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CompositeRulesetResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: GetWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResponse,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: GetWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResponse,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: GetWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResponse,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: GetWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerErrorResponse,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/INTERNAL-SERVER-ERROR",
				Title:      "Internal Server Error",
				Detail:     "Error fetching WAF composite ruleset",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetWAFCompositeRuleset(
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

func TestAppSec_UpdateWAFCompositeRuleset(t *testing.T) {
	result := CompositeRulesetResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestWAFCompositeRuleset/WAFCompositeRuleset.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params              UpdateWAFCompositeRulesetRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CompositeRulesetResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{
		"200 OK": {
			params: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			expectedRequestBody: `{"ConfigID":43253,"Version":15,"PolicyID":"AAAA_81230"}`,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResponse,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResponse,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResponse,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerErrorResponse,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/INTERNAL-SERVER-ERROR",
				Title:      "Internal Server Error",
				Detail:     "Error fetching WAF composite ruleset",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"400 invalid input - invalid condition type": {
			params: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody: `
				{
				  "type": "https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR",
				  "status": 400,
				  "title": "Invalid Input Error",
				  "detail": "Invalid input specified for attackGroups.conditionException.advancedExceptions.conditions",
				  "instance": "https://problems.luna.akamaiapis.net/appsec/error-instances/5cb90e7bf784d05a"
				}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR",
				Title:      "Invalid Input Error",
				Detail:     "Invalid input specified for attackGroups.conditionException.advancedExceptions.conditions",
				StatusCode: http.StatusBadRequest,
			},
		},
		"400 validation error - invalid values": {
			params: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody: `
				{
				  "type": "https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR",
				  "status": 400,
				  "title": "Invalid Input Error",
				  "detail": "Validation failed: attackGroups[0].conditionException.advancedExceptions.conditions[0].hosts: com is not valid.; attackGroups[0].conditionException.advancedExceptions.conditions[1].ips: Invalid IP(s): 1; attackGroups[0].conditionException.advancedExceptions.conditions[2].paths: Value [test] contains invalid relative paths.",
				  "instance": "https://problems.luna.akamaiapis.net/appsec/error-instances/daab3e01537b04f2",
				  "fieldErrors": {
				    "attackGroups[0].conditionException.advancedExceptions.conditions[0].hosts": [
				      "com is not valid."
				    ],
				    "attackGroups[0].conditionException.advancedExceptions.conditions[1].ips": [
				      "Invalid IP(s): 1"
				    ],
				    "attackGroups[0].conditionException.advancedExceptions.conditions[2].paths": [
				      "Value [test] contains invalid relative paths."
				    ]
				  }
				}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/web-application-firewall/ruleset",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR",
				Title:      "Invalid Input Error",
				Detail:     "Validation failed: attackGroups[0].conditionException.advancedExceptions.conditions[0].hosts: com is not valid.; attackGroups[0].conditionException.advancedExceptions.conditions[1].ips: Invalid IP(s): 1; attackGroups[0].conditionException.advancedExceptions.conditions[2].paths: Value [test] contains invalid relative paths.",
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPatch, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateWAFCompositeRuleset(
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

func TestGetWAFCompositeRulesetRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           GetWAFCompositeRulesetRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: GetWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: GetWAFCompositeRulesetRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing Version": {
			req: GetWAFCompositeRulesetRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing PolicyID": {
			req: GetWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
			},
			errorExpected: true,
		},
		"all fields missing": {
			req:           GetWAFCompositeRulesetRequest{},
			errorExpected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateWAFCompositeRulesetRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           UpdateWAFCompositeRulesetRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: UpdateWAFCompositeRulesetRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing Version": {
			req: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			errorExpected: true,
		},
		"missing PolicyID": {
			req: UpdateWAFCompositeRulesetRequest{
				ConfigID: 43253,
				Version:  15,
			},
			errorExpected: true,
		},
		"all fields missing": {
			req:           UpdateWAFCompositeRulesetRequest{},
			errorExpected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
