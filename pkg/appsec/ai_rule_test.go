package appsec

import (
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

func TestAppSec_ListAIRules(t *testing.T) {
	result := ListAIRulesResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestAIRule/AIRules.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	tests := map[string]struct {
		params           ListAIRulesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListAIRulesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:           ListAIRulesRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288"},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules",
			expectedResponse: &result,
		},
		"400 bad request": {
			params:         ListAIRulesRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288"},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"500 internal server error": {
			params:         ListAIRulesRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288"},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error500), "want: %s; got: %s", error500, err)
			},
		},
		"validate - missing ConfigID": {
			params: ListAIRulesRequest{Version: 25, PolicyID: "boBF_19288"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing Version": {
			params: ListAIRulesRequest{ConfigID: 77653, PolicyID: "boBF_19288"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: ListAIRulesRequest{ConfigID: 77653, Version: 25},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - all missing": {
			params: ListAIRulesRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank\nPolicyID: cannot be blank\nVersion: cannot be blank", err.Error())
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
			result, err := client.ListAIRules(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_GetAIRulesStatus(t *testing.T) {
	result := GetAIRulesStatusResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestAIRule/AIRulesStatus.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	tests := map[string]struct {
		params           GetAIRulesStatusRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAIRulesStatusResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:           GetAIRulesStatusRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288"},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/status",
			expectedResponse: &result,
		},
		"400 bad request": {
			params:         GetAIRulesStatusRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288"},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/status",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"500 internal server error": {
			params:         GetAIRulesStatusRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288"},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/status",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error500), "want: %s; got: %s", error500, err)
			},
		},
		"validate - missing ConfigID": {
			params: GetAIRulesStatusRequest{Version: 25, PolicyID: "boBF_19288"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing Version": {
			params: GetAIRulesStatusRequest{ConfigID: 77653, PolicyID: "boBF_19288"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: GetAIRulesStatusRequest{ConfigID: 77653, Version: 25},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - all missing": {
			params: GetAIRulesStatusRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank\nPolicyID: cannot be blank\nVersion: cannot be blank", err.Error())
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
			result, err := client.GetAIRulesStatus(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateAIRulesStatus(t *testing.T) {
	tests := map[string]struct {
		params              UpdateAIRulesStatusRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *UpdateAIRulesStatusResponse
		withError           func(*testing.T, error)
	}{
		"200 OK - enable": {
			params: UpdateAIRulesStatusRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288",
				Body: UpdateAIRulesStatusRequestBody{AIRuleStatus: "ENABLED"},
			},
			responseStatus:      http.StatusOK,
			responseBody:        `{"aiRuleStatus":"ENABLED"}`,
			expectedPath:        "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/status",
			expectedRequestBody: `{"aiRuleStatus":"ENABLED"}`,
			expectedResponse:    &UpdateAIRulesStatusResponse{AIRuleStatus: "ENABLED"},
		},
		"200 OK - disable": {
			params: UpdateAIRulesStatusRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288",
				Body: UpdateAIRulesStatusRequestBody{AIRuleStatus: "DISABLED"},
			},
			responseStatus:      http.StatusOK,
			responseBody:        `{"aiRuleStatus":"DISABLED"}`,
			expectedPath:        "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/status",
			expectedRequestBody: `{"aiRuleStatus":"DISABLED"}`,
			expectedResponse:    &UpdateAIRulesStatusResponse{AIRuleStatus: "DISABLED"},
		},
		"400 bad request": {
			params: UpdateAIRulesStatusRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288",
				Body: UpdateAIRulesStatusRequestBody{AIRuleStatus: "ENABLED"},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/status",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"500 internal server error": {
			params: UpdateAIRulesStatusRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288",
				Body: UpdateAIRulesStatusRequestBody{AIRuleStatus: "ENABLED"},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/status",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error500), "want: %s; got: %s", error500, err)
			},
		},
		"validate - missing ConfigID": {
			params: UpdateAIRulesStatusRequest{
				Version: 25, PolicyID: "boBF_19288",
				Body: UpdateAIRulesStatusRequestBody{AIRuleStatus: "ENABLED"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing Version": {
			params: UpdateAIRulesStatusRequest{
				ConfigID: 77653, PolicyID: "boBF_19288",
				Body: UpdateAIRulesStatusRequestBody{AIRuleStatus: "ENABLED"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: UpdateAIRulesStatusRequest{
				ConfigID: 77653, Version: 25,
				Body: UpdateAIRulesStatusRequestBody{AIRuleStatus: "ENABLED"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - invalid status": {
			params: UpdateAIRulesStatusRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288",
				Body: UpdateAIRulesStatusRequestBody{AIRuleStatus: "INVALID"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAIRuleStatus: must be a valid value\n}", err.Error())
			},
		},
		"validate - missing status": {
			params: UpdateAIRulesStatusRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAIRuleStatus: cannot be blank\n}", err.Error())
			},
		},
		"validate - all missing": {
			params: UpdateAIRulesStatusRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAIRuleStatus: cannot be blank\n}\nConfigID: cannot be blank\nPolicyID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAIRulesStatus(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_GetAIRuleAction(t *testing.T) {
	result := GetAIRuleActionResponse{}
	respData := compactJSON(loadFixtureBytes("testdata/TestAIRule/AIRuleAction.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	tests := map[string]struct {
		params           GetAIRuleActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAIRuleActionResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:           GetAIRuleActionRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/3001000/versions/1/action",
			expectedResponse: &result,
		},
		"404 not found": {
			params:         GetAIRuleActionRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1},
			responseStatus: http.StatusNotFound,
			responseBody:   `{"type":"not_found","title":"Not Found","status":404,"detail":"AI rule not found for ruleId=3001000 ruleVersionId=1"}`,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/3001000/versions/1/action",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, &Error{
					Type:       "not_found",
					Title:      "Not Found",
					Detail:     "AI rule not found for ruleId=3001000 ruleVersionId=1",
					StatusCode: http.StatusNotFound,
				}), "want a not-found API error; got: %v", err)
			},
		},
		"500 internal server error": {
			params:         GetAIRuleActionRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/3001000/versions/1/action",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error500), "want: %s; got: %s", error500, err)
			},
		},
		"validate - missing ConfigID": {
			params: GetAIRuleActionRequest{Version: 25, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing Version": {
			params: GetAIRuleActionRequest{ConfigID: 77653, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: GetAIRuleActionRequest{ConfigID: 77653, Version: 25, RuleID: 3001000, RuleVersionID: 1},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - missing RuleID": {
			params: GetAIRuleActionRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleVersionID: 1},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: cannot be blank", err.Error())
			},
		},
		"validate - missing RuleVersionID": {
			params: GetAIRuleActionRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 3001000},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleVersionID: cannot be blank", err.Error())
			},
		},
		"validate - invalid RuleID": {
			params: GetAIRuleActionRequest{ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 9999999, RuleVersionID: 1},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: must be a valid value", err.Error())
			},
		},
		"validate - all missing": {
			params: GetAIRuleActionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank\nPolicyID: cannot be blank\nRuleID: cannot be blank\nRuleVersionID: cannot be blank\nVersion: cannot be blank", err.Error())
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
			result, err := client.GetAIRuleAction(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateAIRuleAction(t *testing.T) {
	tests := map[string]struct {
		params              UpdateAIRuleActionRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *UpdateAIRuleActionResponse
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288",
				RuleID: 3001000, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			responseStatus:      http.StatusOK,
			responseBody:        `{"action":"deny"}`,
			expectedPath:        "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/3001000/versions/1/action",
			expectedRequestBody: `{"action":"deny"}`,
			expectedResponse:    &UpdateAIRuleActionResponse{Action: "deny"},
		},
		"400 bad request": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288",
				RuleID: 3001000, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/3001000/versions/1/action",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"500 internal server error": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288",
				RuleID: 3001000, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/3001000/versions/1/action",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error500), "want: %s; got: %s", error500, err)
			},
		},
		"validate - missing ConfigID": {
			params: UpdateAIRuleActionRequest{
				Version: 25, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing Version": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing PolicyID": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, RuleID: 3001000, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - missing RuleID": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: cannot be blank", err.Error())
			},
		},
		"validate - missing RuleVersionID": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 3001000,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleVersionID: cannot be blank", err.Error())
			},
		},
		"validate - invalid RuleID": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 9999999, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: must be a valid value", err.Error())
			},
		},
		"validate - deny_custom action accepted": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "deny_custom_12345"},
			},
			responseStatus:      http.StatusOK,
			responseBody:        `{"action":"deny_custom_12345"}`,
			expectedPath:        "/appsec/v1/configs/77653/versions/25/security-policies/boBF_19288/ai-rules/3001000/versions/1/action",
			expectedRequestBody: `{"action":"deny_custom_12345"}`,
			expectedResponse:    &UpdateAIRuleActionResponse{Action: "deny_custom_12345"},
		},
		"validate - invalid action": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1,
				Body: UpdateAIRuleActionRequestBody{Action: "block"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: must be in a valid format\n}", err.Error())
			},
		},
		"validate - missing action": {
			params: UpdateAIRuleActionRequest{
				ConfigID: 77653, Version: 25, PolicyID: "boBF_19288", RuleID: 3001000, RuleVersionID: 1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: cannot be blank\n}", err.Error())
			},
		},
		"validate - all missing": {
			params: UpdateAIRuleActionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: cannot be blank\n}\nConfigID: cannot be blank\nPolicyID: cannot be blank\nRuleID: cannot be blank\nRuleVersionID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAIRuleAction(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
