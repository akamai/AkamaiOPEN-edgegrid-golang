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

var badRequestResp = `
{
    "type": "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
    "title": "Bad Request",
    "detail": "The request could not be understood by the server due to malformed syntax.",
	"status": 400
}`

var forbiddenResp = `
{
  "detail": "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
  "status": 403,
  "title": "Forbidden",
  "type": "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED"
}`

var notFoundResp = `
{
  "detail": "The requested resource is not found",
  "status": 404,
  "title": "Not Found",
  "type": "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND"
}`

// Test GET URLProtectionPolicies
func TestAppSec_ListURLProtectionPolicies(t *testing.T) {

	result := ListURLProtectionPoliciesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestURLProtection/URLProtectionPolicies.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           ListURLProtectionPoliciesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListURLProtectionPoliciesResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: ListURLProtectionPoliciesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/url-protections",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: ListURLProtectionPoliciesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: ListURLProtectionPoliciesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: ListURLProtectionPoliciesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: ListURLProtectionPoliciesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/url-protections",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListURLProtectionPolicies(
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

// Test GET URLProtectionPolicy
func TestAppSec_GetURLProtectionPolicy(t *testing.T) {

	result := GetURLProtectionPolicyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestURLProtection/URLProtectionPolicy.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetURLProtectionPolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetURLProtectionPolicyResponse
		withError        error
	}{
		"200 OK": {
			params: GetURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: GetURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: GetURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: GetURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: GetURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/url-protections/134644",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetURLProtectionPolicy(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create URLProtectionPolicy with Hostname Paths and APIDefinitions Separately
func TestAppSec_CreateURLProtectionPolicyHostnamePaths(t *testing.T) {

	resultHostnamePath := CreateURLProtectionPolicyResponse{}
	resultAPIDefinition := CreateURLProtectionPolicyResponse{}

	respDataForHostnamePath := compactJSON(loadFixtureBytes("testdata/TestURLProtection/URLProtectionPolicyHostnamePaths.json"))
	err := json.Unmarshal([]byte(respDataForHostnamePath), &resultHostnamePath)
	require.NoError(t, err)

	respDataForAPIDefinition := compactJSON(loadFixtureBytes("testdata/TestURLProtection/URLProtectionPolicyApiDefinitions.json"))
	err = json.Unmarshal([]byte(respDataForAPIDefinition), &resultAPIDefinition)
	require.NoError(t, err)

	tests := map[string]struct {
		params              CreateURLProtectionPolicyRequest
		prop                *CreateURLProtectionPolicyRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateURLProtectionPolicyResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{

		"201 Created with HostnamePaths": {
			params: CreateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
					HostnamePaths: []HostnamePath{
						{
							Hostname: "custom.com",
							Paths:    []string{"/asd", "/my-test-path"},
						},
					},
					Categories: []Category{
						{
							Type: "BOTS",
						},
						{
							Type:          "CLIENT_LIST",
							PositiveMatch: func() *bool { b := true; return &b }(),
							ListIDs: []string{"12345_10CLIENTLIST",
								"54321_123"},
						},
					},
					BypassCondition: &BypassCondition{
						AtomicConditions: []AtomicCondition{
							{
								Type: "NetworkListCondition",
								Values: []string{
									"12345_10CLIENTLIST",
									"54321_123",
								},
							},
							{
								Type: "RequestHeaderCondition",
								Names: []string{
									"my-custom-header",
								},
								NameWildcard: func() *bool { b := true; return &b }(),
								Values: []string{
									"my-custom-value",
								},
								ValueCase:     func() *bool { b := false; return &b }(),
								ValueWildcard: func() *bool { b := false; return &b }(),
							},
						},
					},
				},
			},
			expectedRequestBody: `{"name":"Test URL Protection Policy with Hostname Paths","bypassCondition":{"atomicConditions":[{"className":"NetworkListCondition","value":["12345_10CLIENTLIST","54321_123"]},{"className":"RequestHeaderCondition","name":["my-custom-header"],"nameWildcard":true,"value":["my-custom-value"],"valueCase":false,"valueWildcard":false}]},"rateThreshold":400,"hostnamePaths":[{"hostname":"custom.com","paths":["/asd","/my-test-path"]}],"intelligentLoadShedding":false,"categories":[{"type":"BOTS","positiveMatch":null},{"type":"CLIENT_LIST","positiveMatch":true,"listIds":["12345_10CLIENTLIST","54321_123"]}]}`,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respDataForHostnamePath,
			expectedResponse: &resultHostnamePath,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/url-protections",
		},
		"201 Created with APIDefinitions": {
			params: CreateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 195,
					APIDefinitions: []APIDefinition{
						{
							APIDefinitionID:    3216157,
							DefinedResources:   true,
							ResourceIDs:        []int64{},
							UndefinedResources: true,
						},
					},
					Categories: []Category{
						{
							Type: "BOTS",
						},
						{
							Type:          "CLIENT_LIST",
							PositiveMatch: func() *bool { b := true; return &b }(),
							ListIDs: []string{"12345_10CLIENTLIST",
								"54321_123"},
						},
					},
					BypassCondition: &BypassCondition{
						AtomicConditions: []AtomicCondition{
							{
								Type: "NetworkListCondition",
								Values: []string{
									"12345_10CLIENTLIST",
									"54321_123",
								},
							},
							{
								Type: "RequestHeaderCondition",
								Names: []string{
									"my-custom-header",
								},
								NameWildcard: func() *bool { b := true; return &b }(),
								Values: []string{
									"my-custom-value",
								},
								ValueCase:     func() *bool { b := false; return &b }(),
								ValueWildcard: func() *bool { b := false; return &b }(),
							},
						},
					},
				},
			},
			expectedRequestBody: `{"name":"Test URL Protection Policy with Hostname Paths","bypassCondition":{"atomicConditions":[{"className":"NetworkListCondition","value":["12345_10CLIENTLIST","54321_123"]},{"className":"RequestHeaderCondition","name":["my-custom-header"],"nameWildcard":true,"value":["my-custom-value"],"valueCase":false,"valueWildcard":false}]},"rateThreshold":195,"apiDefinitions":[{"apiDefinitionId":3216157,"definedResources":true,"resourceIds":[],"undefinedResources":true}],"intelligentLoadShedding":false,"categories":[{"type":"BOTS","positiveMatch":null},{"type":"CLIENT_LIST","positiveMatch":true,"listIds":["12345_10CLIENTLIST","54321_123"]}]}`,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respDataForAPIDefinition,
			expectedResponse: &resultAPIDefinition,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/url-protections",
		},
		"400 bad request": {
			params: CreateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 200,
					HostnamePaths: []HostnamePath{
						{
							Hostname: "custom.com",
							Paths:    []string{"/asd", "/my-test-path"},
						},
					},
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: CreateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 200,
					HostnamePaths: []HostnamePath{
						{
							Hostname: "custom.com",
							Paths:    []string{"/asd", "/my-test-path"},
						},
					},
				},
			},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: CreateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 200,
					HostnamePaths: []HostnamePath{
						{
							Hostname: "custom.com",
							Paths:    []string{"/asd", "/my-test-path"},
						},
					},
				},
			},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: CreateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 200,
					HostnamePaths: []HostnamePath{
						{
							Hostname: "custom.com",
							Paths:    []string{"/asd", "/my-test-path"},
						},
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating domain"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/url-protections",
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
			result, err := client.CreateURLProtectionPolicy(
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

// Test Update URLProtectionPolicy
func TestAppSec_UpdateURLProtectionPolicy(t *testing.T) {
	result := UpdateURLProtectionPolicyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestURLProtection/URLProtectionPolicy.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateURLProtectionPolicyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestURLProtection/URLProtectionPolicy.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params              UpdateURLProtectionPolicyRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *UpdateURLProtectionPolicyResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{
		"200 Success": {
			params: UpdateURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			expectedRequestBody: `{"name":"Test URL Protection Policy with Hostname Paths","rateThreshold":400,"intelligentLoadShedding":false}`,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/url-protections/134644",
		},
		"400 bad request": {
			params: UpdateURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: UpdateURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: UpdateURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: UpdateURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/url-protections/134644",
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
			result, err := client.UpdateURLProtectionPolicy(
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

// Test Remove URLProtectionPolicy
func TestAppSec_RemoveURLProtectionPolicy(t *testing.T) {

	tests := map[string]struct {
		params         RemoveURLProtectionPolicyRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
		headers        http.Header
	}{
		"204 No Content": {
			params: RemoveURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
		},
		"400 bad request": {
			params: RemoveURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequestResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/BAD-REQUEST",
				Title:      "Bad Request",
				Detail:     "The request could not be understood by the server due to malformed syntax.",
				StatusCode: http.StatusBadRequest,
			},
		},
		"403 Forbidden": {
			params: RemoveURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus: http.StatusForbidden,
			responseBody:   forbiddenResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec-resource/error-types/ACCESS-DENIED",
				Title:      "Forbidden",
				Detail:     "You do not have the necessary access to perform this operation or the requested resource cannot be modified",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 Not Found": {
			params: RemoveURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus: http.StatusNotFound,
			responseBody:   notFoundResp,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "The requested resource is not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params: RemoveURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting url protection policy"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/url-protections/134644",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting url protection policy",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			err := client.RemoveURLProtectionPolicy(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// Test RemoveURLProtectionPolicyRequest Validate
func TestRemoveURLProtectionPolicyRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           RemoveURLProtectionPolicyRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: RemoveURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: RemoveURLProtectionPolicyRequest{
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigVersion": {
			req: RemoveURLProtectionPolicyRequest{
				ConfigID:              43253,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing URLProtectionPolicyID": {
			req: RemoveURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing all required fields": {
			req:           RemoveURLProtectionPolicyRequest{},
			errorExpected: true,
		},
		"missing ConfigID and ConfigVersion": {
			req: RemoveURLProtectionPolicyRequest{
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigID and URLProtectionPolicyID": {
			req: RemoveURLProtectionPolicyRequest{
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing ConfigVersion and URLProtectionPolicyID": {
			req: RemoveURLProtectionPolicyRequest{
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

// Test GetURLProtectionPolicyRequest Validate
func TestGetURLProtectionPolicyRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           GetURLProtectionPolicyRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: GetURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: GetURLProtectionPolicyRequest{
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigVersion": {
			req: GetURLProtectionPolicyRequest{
				ConfigID:              43253,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing URLProtectionPolicyID": {
			req: GetURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing all required fields": {
			req:           GetURLProtectionPolicyRequest{},
			errorExpected: true,
		},
		"missing ConfigID and ConfigVersion": {
			req: GetURLProtectionPolicyRequest{
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigID and URLProtectionPolicyID": {
			req: GetURLProtectionPolicyRequest{
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing ConfigVersion and URLProtectionPolicyID": {
			req: GetURLProtectionPolicyRequest{
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

// Test ListURLProtectionPoliciesRequest Validate
func TestListURLProtectionPoliciesRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           ListURLProtectionPoliciesRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: ListURLProtectionPoliciesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: ListURLProtectionPoliciesRequest{
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing ConfigVersion": {
			req: ListURLProtectionPoliciesRequest{
				ConfigID: 43253,
			},
			errorExpected: true,
		},
		"missing all required fields": {
			req:           ListURLProtectionPoliciesRequest{},
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

// Test CreateURLProtectionPolicyRequest Validate
func TestCreateURLProtectionPolicyRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           CreateURLProtectionPolicyRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: CreateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: CreateURLProtectionPolicyRequest{
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing ConfigVersion": {
			req: CreateURLProtectionPolicyRequest{
				ConfigID: 43253,
			},
			errorExpected: true,
		},
		"missing Body": {
			req: CreateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing all required fields": {
			req:           CreateURLProtectionPolicyRequest{},
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

// Test UpdateURLProtectionPolicyRequest Validate
func TestUpdateURLProtectionPolicyRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           UpdateURLProtectionPolicyRequest
		errorExpected bool
	}{
		"valid request - all fields populated": {
			req: UpdateURLProtectionPolicyRequest{
				ConfigID:              43253,
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			errorExpected: false,
		},
		"missing ConfigID": {
			req: UpdateURLProtectionPolicyRequest{
				ConfigVersion:         15,
				URLProtectionPolicyID: 134644,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			errorExpected: true,
		},
		"missing ConfigVersion": {
			req: UpdateURLProtectionPolicyRequest{
				ConfigID:              43253,
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing URLProtectionPolicyID": {
			req: UpdateURLProtectionPolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Body: URLProtectionPolicyRequestBody{
					Name:             "Test URL Protection Policy with Hostname Paths",
					MaxRateThreshold: 400,
				},
			},
			errorExpected: true,
		},
		"missing all required fields": {
			req:           UpdateURLProtectionPolicyRequest{},
			errorExpected: true,
		},
		"missing ConfigID and ConfigVersion": {
			req: UpdateURLProtectionPolicyRequest{
				URLProtectionPolicyID: 134644,
			},
			errorExpected: true,
		},
		"missing ConfigID and URLProtectionPolicyID": {
			req: UpdateURLProtectionPolicyRequest{
				ConfigVersion: 15,
			},
			errorExpected: true,
		},
		"missing ConfigVersion and URLProtectionPolicyID": {
			req: UpdateURLProtectionPolicyRequest{
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
