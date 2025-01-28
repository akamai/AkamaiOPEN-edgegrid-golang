package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListRapidRule(t *testing.T) {
	result := GetRapidRulesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRapidRule/RapidRules.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetRapidRulesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRapidRulesResponse
		withError        func(*testing.T, error)
		headers          http.Header
	}{
		"200 OK": {
			params: GetRapidRulesRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: GetRapidRulesRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules",
			withError: func(t *testing.T, err error) {
				want := error400
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: GetRapidRulesRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules",
			withError: func(t *testing.T, err error) {
				want := error500
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"rapid rule with ID not found": {
			params: GetRapidRulesRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   ptr.To(int64(1234567890)),
			},
			headers:        http.Header{},
			responseStatus: http.StatusOK,
			responseBody:   respData,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get rapid rule failure. rapid rule with ID: 1234567890 not found", err.Error())
			},
		},
		"validate - missing config ID": {
			params: GetRapidRulesRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing version": {
			params: GetRapidRulesRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing security policy ID": {
			params: GetRapidRulesRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - required param not provided": {
			params: GetRapidRulesRequest{},
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
			client := mockAPIClient(t, mockServer)
			result, err := client.GetRapidRules(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_GetRapidRulesDefaultAction(t *testing.T) {
	result := GetRapidRulesDefaultActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRapidRule/RapidRulesDefaultAction.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetRapidRulesDefaultActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRapidRulesDefaultActionResponse
		withError        func(*testing.T, error)
		headers          http.Header
	}{
		"200 OK": {
			params: GetRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/action",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: GetRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/action",
			withError: func(t *testing.T, err error) {
				want := error400
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: GetRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/action",
			withError: func(t *testing.T, err error) {
				want := error500
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validate - missing config ID": {
			params: GetRapidRulesDefaultActionRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing version": {
			params: GetRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing security policy ID": {
			params: GetRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - required param not provided": {
			params: GetRapidRulesDefaultActionRequest{},
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
			client := mockAPIClient(t, mockServer)
			result, err := client.GetRapidRulesDefaultAction(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_GetRapidRulesStatus(t *testing.T) {
	result := GetRapidRulesStatusResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRapidRule/RapidRulesStatus.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetRapidRulesStatusRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRapidRulesStatusResponse
		withError        func(*testing.T, error)
		headers          http.Header
	}{
		"200 OK": {
			params: GetRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/status",
			expectedResponse: &result,
		},
		"400 bad request": {
			params: GetRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/status",
			withError: func(t *testing.T, err error) {
				want := error400
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: GetRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/status",
			withError: func(t *testing.T, err error) {
				want := error500
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validate - missing config ID": {
			params: GetRapidRulesStatusRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing version": {
			params: GetRapidRulesStatusRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing security policy ID": {
			params: GetRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - required param not provided": {
			params: GetRapidRulesStatusRequest{},
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
			client := mockAPIClient(t, mockServer)
			result, err := client.GetRapidRulesStatus(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateRapidRulesStatus(t *testing.T) {
	tests := map[string]struct {
		params              UpdateRapidRulesStatusRequest
		responseStatus      int
		responseBody        string
		expectedResponse    *UpdateRapidRulesStatusResponse
		expectedRequestBody string
		expectedPath        string
		withError           func(*testing.T, error)
		headers             http.Header
	}{
		"200 Success": {
			params: UpdateRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesStatusRequestBody{
					Enabled: ptr.To(true),
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"enabled": true}`,
			expectedResponse: ptr.To(UpdateRapidRulesStatusResponse{
				Enabled: true,
			}),
			expectedRequestBody: `{"enabled":true}`,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/status"},
		"400 bad request": {
			params: UpdateRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesStatusRequestBody{
					Enabled: ptr.To(true),
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/status",
			withError: func(t *testing.T, err error) {
				want := error400
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UpdateRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesStatusRequestBody{
					Enabled: ptr.To(true),
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/status",
			withError: func(t *testing.T, err error) {
				want := error500
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validate - missing config ID": {
			params: UpdateRapidRulesStatusRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesStatusRequestBody{
					Enabled: ptr.To(true),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing version": {
			params: UpdateRapidRulesStatusRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesStatusRequestBody{
					Enabled: ptr.To(true),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing security policy ID": {
			params: UpdateRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				Body: UpdateRapidRulesStatusRequestBody{
					Enabled: ptr.To(true),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - missing status": {
			params: UpdateRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tEnabled: is required\n}", err.Error())
			},
		},
		"validate - missing status value": {
			params: UpdateRapidRulesStatusRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body:     UpdateRapidRulesStatusRequestBody{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tEnabled: is required\n}", err.Error())
			},
		},
		"validate - required param not provided": {
			params: UpdateRapidRulesStatusRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tEnabled: is required\n}\nConfigID: cannot be blank\nPolicyID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateRapidRulesStatus(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateRapidRulesDefaultAction(t *testing.T) {
	tests := map[string]struct {
		params              UpdateRapidRulesDefaultActionRequest
		responseStatus      int
		responseBody        string
		expectedResponse    *UpdateRapidRulesDefaultActionResponse
		expectedRequestBody string
		expectedPath        string
		withError           func(*testing.T, error)
		headers             http.Header
	}{
		"200 Success": {
			params: UpdateRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesDefaultActionRequestBody{
					Action: "akamai_managed",
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"action": "akamai_managed"}`,
			expectedResponse: &UpdateRapidRulesDefaultActionResponse{
				Action: "akamai_managed",
			},
			expectedRequestBody: `{"action":"akamai_managed"}`,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/action"},
		"400 bad request": {
			params: UpdateRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesDefaultActionRequestBody{
					Action: "akamai_managed",
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/action",
			withError: func(t *testing.T, err error) {
				want := error400
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UpdateRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesDefaultActionRequestBody{
					Action: "akamai_managed",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/action",
			withError: func(t *testing.T, err error) {
				want := error500
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validate - missing config ID": {
			params: UpdateRapidRulesDefaultActionRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesDefaultActionRequestBody{
					Action: "akamai_managed",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing version": {
			params: UpdateRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRulesDefaultActionRequestBody{
					Action: "akamai_managed",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing security policy ID": {
			params: UpdateRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				Body: UpdateRapidRulesDefaultActionRequestBody{
					Action: "akamai_managed",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - missing default action": {
			params: UpdateRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: cannot be blank\n}", err.Error())
			},
		},
		"validate - missing default action value": {
			params: UpdateRapidRulesDefaultActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body:     UpdateRapidRulesDefaultActionRequestBody{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: cannot be blank\n}", err.Error())
			},
		},
		"validate - required param not provided": {
			params: UpdateRapidRulesDefaultActionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: cannot be blank\n}\nConfigID: cannot be blank\nPolicyID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateRapidRulesDefaultAction(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateRapidRuleActionLock(t *testing.T) {
	tests := map[string]struct {
		params              UpdateRapidRuleActionLockRequest
		responseStatus      int
		responseBody        string
		expectedResponse    *UpdateRapidRuleActionLockResponse
		expectedRequestBody string
		expectedPath        string
		withError           func(*testing.T, error)
		headers             http.Header
	}{
		"200 Success": {
			params: UpdateRapidRuleActionLockRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body: UpdateRapidRuleActionLockRequestBody{
					Enabled: ptr.To(true),
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"enabled": true}`,
			expectedResponse: &UpdateRapidRuleActionLockResponse{
				Enabled: true,
			},
			expectedRequestBody: `{"enabled":true}`,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/lock"},
		"400 bad request": {
			params: UpdateRapidRuleActionLockRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body: UpdateRapidRuleActionLockRequestBody{
					Enabled: ptr.To(true),
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/lock",
			withError: func(t *testing.T, err error) {
				want := error400
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UpdateRapidRuleActionLockRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body: UpdateRapidRuleActionLockRequestBody{
					Enabled: ptr.To(true),
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/lock",
			withError: func(t *testing.T, err error) {
				want := error500
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validate - missing config ID": {
			params: UpdateRapidRuleActionLockRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body: UpdateRapidRuleActionLockRequestBody{
					Enabled: ptr.To(true),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing version": {
			params: UpdateRapidRuleActionLockRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body: UpdateRapidRuleActionLockRequestBody{
					Enabled: ptr.To(true),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing security policy ID": {
			params: UpdateRapidRuleActionLockRequest{
				ConfigID: 43253,
				Version:  15,
				RuleID:   999997,
				Body: UpdateRapidRuleActionLockRequestBody{
					Enabled: ptr.To(true),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - missing ruleID": {
			params: UpdateRapidRuleActionLockRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body: UpdateRapidRuleActionLockRequestBody{
					Enabled: ptr.To(true),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: cannot be blank", err.Error())
			},
		},
		"validate - missing lock": {
			params: UpdateRapidRuleActionLockRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tEnabled: is required\n}", err.Error())
			},
		},
		"validate - missing lock value": {
			params: UpdateRapidRuleActionLockRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body:     UpdateRapidRuleActionLockRequestBody{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tEnabled: is required\n}", err.Error())
			},
		},
		"validate - required param not provided": {
			params: UpdateRapidRuleActionLockRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tEnabled: is required\n}\nConfigID: cannot be blank\nPolicyID: cannot be blank\nRuleID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateRapidRuleActionLock(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateRapidRuleAction(t *testing.T) {
	tests := map[string]struct {
		params              UpdateRapidRuleActionRequest
		responseStatus      int
		responseBody        string
		expectedResponse    *UpdateRapidRuleActionResponse
		expectedRequestBody string
		expectedPath        string
		withError           func(*testing.T, error)
		headers             http.Header
	}{
		"200 Success": {
			params: UpdateRapidRuleActionRequest{
				ConfigID:    43253,
				Version:     15,
				PolicyID:    "AAAA_81230",
				RuleID:      999997,
				RuleVersion: 2,
				Body: UpdateRapidRuleActionRequestBody{
					Action: "alert",
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"action": "alert","lock": true}`,
			expectedResponse: &UpdateRapidRuleActionResponse{
				Action: "alert",
				Lock:   true,
			},
			expectedRequestBody: `{"action":"alert"}`,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/versions/2/action"},
		"400 bad request": {
			params: UpdateRapidRuleActionRequest{
				ConfigID:    43253,
				Version:     15,
				PolicyID:    "AAAA_81230",
				RuleID:      999997,
				RuleVersion: 2,
				Body: UpdateRapidRuleActionRequestBody{
					Action: "alert",
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/versions/2/action",
			withError: func(t *testing.T, err error) {
				want := error400
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UpdateRapidRuleActionRequest{
				ConfigID:    43253,
				Version:     15,
				PolicyID:    "AAAA_81230",
				RuleID:      999997,
				RuleVersion: 2,
				Body: UpdateRapidRuleActionRequestBody{
					Action: "alert",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/versions/2/action",
			withError: func(t *testing.T, err error) {
				want := error500
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validate - missing config ID": {
			params: UpdateRapidRuleActionRequest{
				Version:     15,
				PolicyID:    "AAAA_81230",
				RuleID:      999997,
				RuleVersion: 2,
				Body: UpdateRapidRuleActionRequestBody{
					Action: "alert",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing version": {
			params: UpdateRapidRuleActionRequest{
				ConfigID:    43253,
				PolicyID:    "AAAA_81230",
				RuleID:      999997,
				RuleVersion: 2,
				Body: UpdateRapidRuleActionRequestBody{
					Action: "alert",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing security policy ID": {
			params: UpdateRapidRuleActionRequest{
				ConfigID:    43253,
				Version:     15,
				RuleID:      999997,
				RuleVersion: 2,
				Body: UpdateRapidRuleActionRequestBody{
					Action: "alert",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - missing ruleID": {
			params: UpdateRapidRuleActionRequest{
				ConfigID:    43253,
				Version:     15,
				PolicyID:    "AAAA_81230",
				RuleVersion: 2,
				Body: UpdateRapidRuleActionRequestBody{
					Action: "alert",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: cannot be blank", err.Error())
			},
		},
		"validate - missing rule version": {
			params: UpdateRapidRuleActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body: UpdateRapidRuleActionRequestBody{
					Action: "alert",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleVersion: cannot be blank", err.Error())
			},
		},
		"validate - missing rule action": {
			params: UpdateRapidRuleActionRequest{
				ConfigID:    43253,
				Version:     15,
				PolicyID:    "AAAA_81230",
				RuleID:      999997,
				RuleVersion: 2,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: cannot be blank\n}", err.Error())
			},
		},
		"validate - missing rule action value": {
			params: UpdateRapidRuleActionRequest{
				ConfigID:    43253,
				Version:     15,
				PolicyID:    "AAAA_81230",
				RuleID:      999997,
				RuleVersion: 2,
				Body:        UpdateRapidRuleActionRequestBody{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: cannot be blank\n}", err.Error())
			},
		},
		"validate - required param not provided": {
			params: UpdateRapidRuleActionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body: {\n\tAction: cannot be blank\n}\nConfigID: cannot be blank\nPolicyID: cannot be blank\nRuleID: cannot be blank\nRuleVersion: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateRapidRuleAction(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateRapidRuleException(t *testing.T) {
	body := RuleConditionException{
		Exception: &RuleException{
			SpecificHeaderCookieParamXMLOrJSONNames: &SpecificHeaderCookieParamXMLOrJSONNames{
				{
					Names:    []string{"akaToken"},
					Selector: "REQUEST_HEADERS",
					Wildcard: false,
				},
			},
		},
	}

	expectedResponse := &UpdateRapidRuleExceptionResponse{
		Exception: &RuleException{
			SpecificHeaderCookieParamXMLOrJSONNames: &SpecificHeaderCookieParamXMLOrJSONNames{
				{
					Names:    []string{"akaToken"},
					Selector: "REQUEST_HEADERS",
					Wildcard: false,
				},
			},
		},
	}

	tests := map[string]struct {
		params              UpdateRapidRuleExceptionRequest
		responseStatus      int
		responseBody        string
		expectedResponse    *UpdateRapidRuleExceptionResponse
		expectedRequestBody string
		expectedPath        string
		withError           func(*testing.T, error)
		headers             http.Header
	}{
		"200 Success": {
			params: UpdateRapidRuleExceptionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body:     body,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
			  "exception": {
			    "specificHeaderCookieParamXmlOrJsonNames": [
				  {
				    "names": [
					  "akaToken"
				    ],
				    "selector": "REQUEST_HEADERS",
				    "wildcard": false
				  }
			    ]
			  }
			}`,
			expectedRequestBody: `{"exception":{"specificHeaderCookieParamXmlOrJsonNames":[{"names":["akaToken"],"selector":"REQUEST_HEADERS"}]}}`,
			expectedResponse:    expectedResponse,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/condition-exception"},
		"400 bad request": {
			params: UpdateRapidRuleExceptionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body:     body,
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/condition-exception",
			withError: func(t *testing.T, err error) {
				want := error400
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UpdateRapidRuleExceptionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body:     body,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rapid-rules/999997/condition-exception",
			withError: func(t *testing.T, err error) {
				want := error500
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validate - missing config ID": {
			params: UpdateRapidRuleExceptionRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body:     body,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validate - missing version": {
			params: UpdateRapidRuleExceptionRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
				RuleID:   999997,
				Body:     body,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validate - missing security policy ID": {
			params: UpdateRapidRuleExceptionRequest{
				ConfigID: 43253,
				Version:  15,
				RuleID:   999997,
				Body:     body,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validate - missing ruleID": {
			params: UpdateRapidRuleExceptionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Body:     body,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: RuleID: cannot be blank", err.Error())
			},
		},
		"validate - required param not provided": {
			params: UpdateRapidRuleExceptionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank\nPolicyID: cannot be blank\nRuleID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateRapidRuleException(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

var badRequest = `
{
    "type": "https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR",
    "title": "Invalid Input Error",
    "detail": "configId incorrect type",
	"status": 400
}`

var error400 = &Error{
	Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR",
	Title:      "Invalid Input Error",
	Detail:     "configId incorrect type",
	StatusCode: http.StatusBadRequest,
}

var internalServerError = `
{
    "type": "https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR",
    "title": "Internal Server Error",
    "detail": "The server was unable to complete your request. Please try again later.",
	"status": 400
}`

var error500 = &Error{
	Type:       "https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR",
	Title:      "Internal Server Error",
	Detail:     "The server was unable to complete your request. Please try again later.",
	StatusCode: http.StatusInternalServerError,
}
