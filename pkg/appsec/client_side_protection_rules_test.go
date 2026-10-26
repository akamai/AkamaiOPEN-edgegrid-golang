package appsec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListClientSideProtectionRules(t *testing.T) {
	var result ListClientSideProtectionRulesResponse
	respData := compactJSON(loadFixtureBytes("testdata/TestClientSideProtectionRules/Rules.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	tests := map[string]struct {
		params           ListClientSideProtectionRulesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListClientSideProtectionRulesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:           ListClientSideProtectionRulesRequest{ConfigID: 40576, Version: 6, PolicyID: "ab_86352"},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules",
			expectedResponse: &result,
		},
		"400 bad request": {
			params:         ListClientSideProtectionRulesRequest{ConfigID: 40576, Version: 6, PolicyID: "ab_86352"},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"500 internal server error": {
			params:         ListClientSideProtectionRulesRequest{ConfigID: 40576, Version: 6, PolicyID: "ab_86352"},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error500), "want: %s; got: %s", error500, err)
			},
		},
		"validate - missing ConfigID": {
			params: ListClientSideProtectionRulesRequest{Version: 6, PolicyID: "ab_86352"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank.", err.Error())
			},
		},
		"validate - missing Version": {
			params: ListClientSideProtectionRulesRequest{ConfigID: 40576, PolicyID: "ab_86352"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank.", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: ListClientSideProtectionRulesRequest{ConfigID: 40576, Version: 6},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank.", err.Error())
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
			result, err := client.ListClientSideProtectionRules(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_CreateClientSideProtectionRule(t *testing.T) {
	var result ClientSideProtectionRuleBody
	respData := compactJSON(loadFixtureBytes("testdata/TestClientSideProtectionRules/Rule.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	body := ClientSideProtectionRuleBody{
		RuleName:     "Block PCI Data Access",
		Description:  "Prevents third-party scripts from accessing payment card data",
		ActivityType: "DATA_ACCESS",
		Action:       "deny",
		Conditions: ClientSideProtectionConditions{
			Page: ClientSideProtectionConditionCriteria{
				Condition: "MATCHES",
				Values:    []string{"/checkout", "/billing/*"},
			},
			Script: ClientSideProtectionScriptConditionCriteria{
				ScriptType: "THIRD_PARTY",
				Condition:  "MATCHES",
				Values:     []string{"analytics.example.com/*"},
			},
			ActivityTypeValues: ClientSideProtectionConditionCriteria{
				Condition: "MATCHES",
				Values:    []string{"CREDIT_CARD", "PII"},
			},
		},
	}

	expectedRequestBody := `{
		"ruleName":"Block PCI Data Access",
		"description":"Prevents third-party scripts from accessing payment card data",
		"conditions":{
			"page":{"condition":"MATCHES","values":["/checkout","/billing/*"]},
			"script":{"scriptType":"THIRD_PARTY","condition":"MATCHES","values":["analytics.example.com/*"]},
			"activityTypeValues":{"condition":"MATCHES","values":["CREDIT_CARD","PII"]}
		},
		"activityType":"DATA_ACCESS",
		"action":"deny"
	}`

	tests := map[string]struct {
		params           CreateClientSideProtectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ClientSideProtectionRuleBody
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateClientSideProtectionRuleRequest{
				ConfigID: 40576, Version: 6, PolicyID: "ab_86352", Body: body,
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: CreateClientSideProtectionRuleRequest{
				ConfigID: 40576, Version: 6, PolicyID: "ab_86352", Body: body,
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"validate - missing ConfigID": {
			params: CreateClientSideProtectionRuleRequest{Version: 6, PolicyID: "ab_86352", Body: body},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank.", err.Error())
			},
		},
		"validate - missing Version": {
			params: CreateClientSideProtectionRuleRequest{ConfigID: 40576, PolicyID: "ab_86352", Body: body},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank.", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: CreateClientSideProtectionRuleRequest{ConfigID: 40576, Version: 6, Body: body},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank.", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				reqBody, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, expectedRequestBody, string(reqBody))
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateClientSideProtectionRule(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateClientSideProtectionRule(t *testing.T) {
	var result ClientSideProtectionRuleBody
	respData := compactJSON(loadFixtureBytes("testdata/TestClientSideProtectionRules/Rule.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	ruleID := int64(3000601)
	body := ClientSideProtectionRuleBody{
		RuleID:       &ruleID,
		RuleName:     "Block PCI Data Access",
		Description:  "Prevents third-party scripts from accessing payment card data",
		ActivityType: "DATA_ACCESS",
		Action:       "deny",
		Conditions: ClientSideProtectionConditions{
			Page: ClientSideProtectionConditionCriteria{
				Condition: "MATCHES",
				Values:    []string{"/checkout", "/billing/*"},
			},
			Script: ClientSideProtectionScriptConditionCriteria{
				ScriptType: "THIRD_PARTY",
				Condition:  "MATCHES",
				Values:     []string{"analytics.example.com/*"},
			},
			ActivityTypeValues: ClientSideProtectionConditionCriteria{
				Condition: "MATCHES",
				Values:    []string{"CREDIT_CARD", "PII"},
			},
		},
	}

	expectedRequestBody := `{
		"ruleId":3000601,
		"ruleName":"Block PCI Data Access",
		"description":"Prevents third-party scripts from accessing payment card data",
		"conditions":{
			"page":{"condition":"MATCHES","values":["/checkout","/billing/*"]},
			"script":{"scriptType":"THIRD_PARTY","condition":"MATCHES","values":["analytics.example.com/*"]},
			"activityTypeValues":{"condition":"MATCHES","values":["CREDIT_CARD","PII"]}
		},
		"activityType":"DATA_ACCESS",
		"action":"deny"
	}`

	tests := map[string]struct {
		params           UpdateClientSideProtectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ClientSideProtectionRuleBody
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateClientSideProtectionRuleRequest{
				ConfigID: 40576, Version: 6, PolicyID: "ab_86352", RuleID: ruleID, Body: body,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules/3000601",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: UpdateClientSideProtectionRuleRequest{
				ConfigID: 40576, Version: 6, PolicyID: "ab_86352", RuleID: ruleID, Body: body,
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules/3000601",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"validate - missing ConfigID": {
			params: UpdateClientSideProtectionRuleRequest{Version: 6, PolicyID: "ab_86352", RuleID: ruleID, Body: body},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank.", err.Error())
			},
		},
		"validate - missing Version": {
			params: UpdateClientSideProtectionRuleRequest{ConfigID: 40576, PolicyID: "ab_86352", RuleID: ruleID, Body: body},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank.", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: UpdateClientSideProtectionRuleRequest{ConfigID: 40576, Version: 6, RuleID: ruleID, Body: body},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank.", err.Error())
			},
		},
		"validate - missing RuleID": {
			params: UpdateClientSideProtectionRuleRequest{ConfigID: 40576, Version: 6, PolicyID: "ab_86352", Body: body},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: cannot be blank.", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				reqBody, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, expectedRequestBody, string(reqBody))
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateClientSideProtectionRule(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_DeleteClientSideProtectionRule(t *testing.T) {
	ruleID := int64(3000601)

	tests := map[string]struct {
		params         DeleteClientSideProtectionRuleRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			params:         DeleteClientSideProtectionRuleRequest{ConfigID: 40576, Version: 6, PolicyID: "ab_86352", RuleID: ruleID},
			responseStatus: http.StatusOK,
			responseBody:   `{}`,
			expectedPath:   "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules/3000601",
		},
		"204 No Content": {
			params:         DeleteClientSideProtectionRuleRequest{ConfigID: 40576, Version: 6, PolicyID: "ab_86352", RuleID: ruleID},
			responseStatus: http.StatusNoContent,
			responseBody:   ``,
			expectedPath:   "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules/3000601",
		},
		"400 bad request": {
			params:         DeleteClientSideProtectionRuleRequest{ConfigID: 40576, Version: 6, PolicyID: "ab_86352", RuleID: ruleID},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/40576/versions/6/security-policies/ab_86352/client-side-protection-rules/3000601",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"validate - missing ConfigID": {
			params: DeleteClientSideProtectionRuleRequest{Version: 6, PolicyID: "ab_86352", RuleID: ruleID},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank.", err.Error())
			},
		},
		"validate - missing Version": {
			params: DeleteClientSideProtectionRuleRequest{ConfigID: 40576, PolicyID: "ab_86352", RuleID: ruleID},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank.", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: DeleteClientSideProtectionRuleRequest{ConfigID: 40576, Version: 6, RuleID: ruleID},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank.", err.Error())
			},
		},
		"validate - missing RuleID": {
			params: DeleteClientSideProtectionRuleRequest{ConfigID: 40576, Version: 6, PolicyID: "ab_86352"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: cannot be blank.", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			err := client.DeleteClientSideProtectionRule(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAppSec_cpcError(t *testing.T) {
	// Build a minimal appsec client — cpcError only needs the session for p.Error() fallback.
	server := httptest.NewTLSServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	t.Cleanup(server.Close)
	p := mockAPIClient(t, server).(*appsec)

	makeResp := func(status int, body string) *http.Response {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Request:    req,
		}
	}

	t.Run("CPC error shape surfaces details field", func(t *testing.T) {
		resp := makeResp(http.StatusBadRequest, `{"title":"Bad Request","details":"Client Side Protection is disabled for this policy","statusCode":400}`)
		err := p.cpcError(resp)
		require.Error(t, err)
		var appsecErr *Error
		require.ErrorAs(t, err, &appsecErr)
		assert.Equal(t, "Bad Request", appsecErr.Title)
		assert.Equal(t, "Client Side Protection is disabled for this policy", appsecErr.Detail)
		assert.Equal(t, http.StatusBadRequest, appsecErr.StatusCode)
	})

	t.Run("standard appsec error falls back to p.Error", func(t *testing.T) {
		resp := makeResp(http.StatusBadRequest, `{"type":"https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR","title":"Invalid Input Error","detail":"configId incorrect type","status":400}`)
		err := p.cpcError(resp)
		require.Error(t, err)
		var appsecErr *Error
		require.ErrorAs(t, err, &appsecErr)
		assert.Equal(t, "Invalid Input Error", appsecErr.Title)
		assert.Equal(t, "configId incorrect type", appsecErr.Detail)
		assert.Equal(t, http.StatusBadRequest, appsecErr.StatusCode)
	})

	t.Run("empty details field falls back to p.Error", func(t *testing.T) {
		resp := makeResp(http.StatusBadRequest, `{"title":"Bad Request","details":"","statusCode":400}`)
		err := p.cpcError(resp)
		require.Error(t, err)
		var appsecErr *Error
		require.ErrorAs(t, err, &appsecErr)
		assert.Equal(t, "Bad Request", appsecErr.Title)
		assert.Empty(t, appsecErr.Detail) // p.Error parses "detail" not "details", so Detail is empty
		assert.Equal(t, http.StatusBadRequest, appsecErr.StatusCode)
	})
}
