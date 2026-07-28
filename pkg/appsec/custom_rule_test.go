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

func TestAppSec_ListCustomRules(t *testing.T) {

	result := GetCustomRulesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRules.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetCustomRulesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCustomRulesResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetCustomRulesRequest{
				ConfigID: 43253,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/custom-rules",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetCustomRulesRequest{
				ConfigID: 43253,
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
			expectedPath: "/appsec/v1/configs/43253/custom-rules",
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
			result, err := client.GetCustomRules(
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

// Test CustomRule
func TestAppSec_GetCustomRule(t *testing.T) {

	result := GetCustomRuleResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRule.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetCustomRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCustomRuleResponse
		withError        error
	}{
		"200 OK": {
			params: GetCustomRuleRequest{
				ConfigID: 43253,
				ID:       60039625,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/custom-rules/60039625",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetCustomRuleRequest{
				ConfigID: 43253,
				ID:       60039625,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/custom-rules/60039625",
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
			result, err := client.GetCustomRule(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create CustomRule
func TestAppSec_CreateCustomRule(t *testing.T) {

	result := CreateCustomRuleResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRule.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRule.json"))
	req := CreateCustomRuleRequest{
		ConfigID:       43253,
		JsonPayloadRaw: json.RawMessage(reqData),
	}

	tests := map[string]struct {
		params              CreateCustomRuleRequest
		prop                *CreateCustomRuleRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateCustomRuleResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{
		"201 Created": {
			params:              req,
			expectedRequestBody: reqData,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/custom-rules",
		},
		"500 internal server error": {
			params:         req,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating domain"
			}`,
			expectedPath: "/appsec/v1/configs/43253/custom-rules",
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
			result, err := client.CreateCustomRule(
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

// Test Update CustomRule
func TestAppSec_UpdateCustomRule(t *testing.T) {
	result := UpdateCustomRuleResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRule.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRule.json"))
	req := UpdateCustomRuleRequest{
		ConfigID:       43253,
		ID:             60039625,
		JsonPayloadRaw: json.RawMessage(reqData),
	}

	tests := map[string]struct {
		params              UpdateCustomRuleRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *UpdateCustomRuleResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{
		"200 Success": {
			params:              req,
			expectedRequestBody: reqData,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/custom-rules/%d",
		},
		"500 internal server error": {
			params:         req,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/custom-rules/%d",
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
			result, err := client.UpdateCustomRule(
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

// Test Remove CustomRule
func TestAppSec_RemoveCustomRule(t *testing.T) {

	result := RemoveCustomRuleResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRulesEmpty.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := RemoveCustomRuleRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRulesEmpty.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           RemoveCustomRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveCustomRuleResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveCustomRuleRequest{
				ConfigID: 43253,
				ID:       60039625,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/custom-rules/%d",
		},
		"500 internal server error": {
			params: RemoveCustomRuleRequest{
				ConfigID: 43253,
				ID:       60039625,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/custom-rules/%d",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting match target",
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
			result, err := client.RemoveCustomRule(
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

func TestAppSec_GetCustomRuleUsage(t *testing.T) {
	result := GetCustomRulesUsageResponse{}
	respData := `{
		"rules": [
			{
				"ruleId": 12345,
				"policies": [
					{
						"policyId": "POLICY_ID1",
						"policyName": "Policy One"
					},
					{
						"policyId": "POLICY_ID2",
						"policyName": "Policy Two"
					}
				]
			},
			{
				"ruleId": 67890,
				"policies": [
					{
						"policyId": "POLICY_ID3",
						"policyName": "Policy Three"
					}
				]
			}
		]
	}`
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetCustomRulesUsageRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCustomRulesUsageResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  7,
				RequestBody: RuleIDs{
					IDs: []int64{12345, 67890},
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/7/custom-rules/usage",
			expectedResponse: &result,
		},
		"200 OK - empty rules": {
			params: GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  7,
				RequestBody: RuleIDs{
					IDs: []int64{12345, 67890},
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"rules": []}`,
			expectedPath:   "/appsec/v1/configs/43253/versions/7/custom-rules/usage",
			expectedResponse: &GetCustomRulesUsageResponse{
				Rules: []CustomRuleUsage{},
			},
		},
		"403 forbidden": {
			params: GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  7,
				RequestBody: RuleIDs{
					IDs: []int64{12345, 67890},
				},
			},
			headers:        http.Header{},
			responseStatus: http.StatusForbidden,
			responseBody: `{
				"type": "https://problems.luna.akamaiapis.net/appsec/error-types/UNAUTHORIZED",
				"status": 403,
				"title": "Unauthorized Access/Action",
				"detail": "You do not have the necessary access to perform this operation.",
				"instance": "https://problems.luna.akamaiapis.net/appsec/error-instances/65ecbe9f46a3eb5c"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/7/custom-rules/usage",
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/UNAUTHORIZED",
				Title:      "Unauthorized Access/Action",
				Detail:     "You do not have the necessary access to perform this operation.",
				StatusCode: http.StatusForbidden,
			},
		},
		"500 internal server error": {
			params: GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  7,
				RequestBody: RuleIDs{
					IDs: []int64{12345, 67890},
				},
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching custom rule usage",
				"status": 500
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/7/custom-rules/usage",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching custom rule usage",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)

				// Check request body
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)

				var requestBody RuleIDs
				err = json.Unmarshal(body, &requestBody)
				assert.NoError(t, err)
				assert.Equal(t, test.params.RequestBody.IDs, requestBody.IDs)

				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCustomRulesUsage(
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

func TestGetCustomRuleUsageRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		params    GetCustomRulesUsageRequest
		withError func(*testing.T, error)
	}{
		"valid request": {
			params: GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  7,
				RequestBody: RuleIDs{
					IDs: []int64{12345},
				},
			},
		},
		"missing ConfigID": {
			params: GetCustomRulesUsageRequest{
				Version: 7,
				RequestBody: RuleIDs{
					IDs: []int64{12345},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ConfigID: cannot be blank", err.Error())
			},
		},
		"missing Version": {
			params: GetCustomRulesUsageRequest{
				ConfigID: 43253,
				RequestBody: RuleIDs{
					IDs: []int64{12345},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "Version: cannot be blank", err.Error())
			},
		},
		"nil rule IDs": {
			params: GetCustomRulesUsageRequest{
				ConfigID:    43253,
				Version:     7,
				RequestBody: RuleIDs{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "RuleIDs: {\n\tIDs: cannot be blank\n}", err.Error())
			},
		},
		"empty rule IDs": {
			params: GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  7,
				RequestBody: RuleIDs{
					IDs: []int64{},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "RuleIDs: {\n\tIDs: cannot be blank\n}", err.Error())
			},
		},
		"multiple rule IDs": {
			params: GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  7,
				RequestBody: RuleIDs{
					IDs: []int64{12345, 67890, 98765},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.params.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
