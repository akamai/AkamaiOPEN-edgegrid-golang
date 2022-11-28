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

func TestAppSec_ListEvalRule(t *testing.T) {

	result := GetEvalRulesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestEvalRule/EvalRules.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetEvalRulesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetEvalRulesResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetEvalRulesRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-rules?includeConditionException=true",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetEvalRulesRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-rules?includeConditionException=true",
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
			result, err := client.GetEvalRules(
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

// Test EvalRule
func TestAppSec_GetEvalRule(t *testing.T) {

	result := GetEvalRuleResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestEvalRule/EvalRule.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetEvalRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetEvalRuleResponse
		withError        error
	}{
		"200 OK": {
			params: GetEvalRuleRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   12345,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-rules/12345?includeConditionException=true",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetEvalRuleRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   12345,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-rules/12345?includeConditionException=true",
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
			result, err := client.GetEvalRule(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update EvalRule.
func TestAppSec_UpdateEvalRule(t *testing.T) {
	result := UpdateEvalRuleResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestEvalRule/EvalRule.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateEvalRuleRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestEvalRule/EvalRule.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateEvalRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateEvalRuleResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateEvalRuleRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   12345,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-rules/12345/action-condition-exceptions",
		},
		"500 internal server error": {
			params: UpdateEvalRuleRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   12345,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-rules/12345/action-condition-exceptions",
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
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateEvalRule(
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
