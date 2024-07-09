package botman

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Get ContentProtectionJavaScriptInjectionRule List
func TestBotman_GetContentProtectionJavaScriptInjectionRuleList(t *testing.T) {

	tests := map[string]struct {
		params           GetContentProtectionJavaScriptInjectionRuleListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetContentProtectionJavaScriptInjectionRuleListResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetContentProtectionJavaScriptInjectionRuleListRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"contentProtectionJavaScriptInjectionRules": [
		{"contentProtectionJavaScriptInjectionRuleId":"fake3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"contentProtectionJavaScriptInjectionRuleId":"fakead64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"contentProtectionJavaScriptInjectionRuleId":"fake4ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"contentProtectionJavaScriptInjectionRuleId":"faked85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-javascript-injection-rules",
			expectedResponse: &GetContentProtectionJavaScriptInjectionRuleListResponse{
				ContentProtectionJavaScriptInjectionRules: []map[string]interface{}{
					{"contentProtectionJavaScriptInjectionRuleId": "fake3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
					{"contentProtectionJavaScriptInjectionRuleId": "fakead64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
					{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
					{"contentProtectionJavaScriptInjectionRuleId": "fake4ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
					{"contentProtectionJavaScriptInjectionRuleId": "faked85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
				},
			},
		},
		"200 OK One Record": {
			params: GetContentProtectionJavaScriptInjectionRuleListRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"contentProtectionJavaScriptInjectionRules":[
		{"contentProtectionJavaScriptInjectionRuleId":"fake3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"contentProtectionJavaScriptInjectionRuleId":"fakead64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"contentProtectionJavaScriptInjectionRuleId":"fake4ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"contentProtectionJavaScriptInjectionRuleId":"faked85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-javascript-injection-rules",
			expectedResponse: &GetContentProtectionJavaScriptInjectionRuleListResponse{
				ContentProtectionJavaScriptInjectionRules: []map[string]interface{}{
					{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				},
			},
		},
		"500 internal server error": {
			params: GetContentProtectionJavaScriptInjectionRuleListRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching data",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-javascript-injection-rules",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: GetContentProtectionJavaScriptInjectionRuleListRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetContentProtectionJavaScriptInjectionRuleListRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: GetContentProtectionJavaScriptInjectionRuleListRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
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
			result, err := client.GetContentProtectionJavaScriptInjectionRuleList(
				session.ContextWithOptions(
					context.Background(),
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

// Test Get ContentProtectionJavaScriptInjectionRule
func TestBotman_GetContentProtectionJavaScriptInjectionRule(t *testing.T) {
	tests := map[string]struct {
		params           GetContentProtectionJavaScriptInjectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-javascript-injection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
			expectedResponse: map[string]interface{}{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
		},
		"500 internal server error": {
			params: GetContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-javascript-injection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: GetContentProtectionJavaScriptInjectionRuleRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: GetContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID: 43253,
				Version:  15,
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing ContentProtectionJavaScriptInjectionRuleID": {
			params: GetContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				Version:          15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContentProtectionJavaScriptInjectionRuleID")
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
			result, err := client.GetContentProtectionJavaScriptInjectionRule(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create ContentProtectionJavaScriptInjectionRule
func TestBotman_CreateContentProtectionJavaScriptInjectionRule(t *testing.T) {

	tests := map[string]struct {
		params           CreateContentProtectionJavaScriptInjectionRuleRequest
		prop             *CreateContentProtectionJavaScriptInjectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusCreated,
			responseBody:     `{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-javascript-injection-rules",
		},
		"500 internal server error": {
			params: CreateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-javascript-injection-rules",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: CreateContentProtectionJavaScriptInjectionRuleRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: CreateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing JsonPayload": {
			params: CreateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "JsonPayload")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateContentProtectionJavaScriptInjectionRule(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update ContentProtectionJavaScriptInjectionRule
func TestBotman_UpdateContentProtectionJavaScriptInjectionRule(t *testing.T) {
	tests := map[string]struct {
		params           UpdateContentProtectionJavaScriptInjectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/content-protection-javascript-injection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: UpdateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/content-protection-javascript-injection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating zone",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: UpdateContentProtectionJavaScriptInjectionRuleRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: UpdateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID: 43253,
				Version:  15,
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing JsonPayload": {
			params: UpdateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "JsonPayload")
			},
		},
		"Missing ContentProtectionJavaScriptInjectionRuleID": {
			params: UpdateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContentProtectionJavaScriptInjectionRuleID")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateContentProtectionJavaScriptInjectionRule(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Remove ContentProtectionJavaScriptInjectionRule
func TestBotman_RemoveContentProtectionJavaScriptInjectionRule(t *testing.T) {
	tests := map[string]struct {
		params           RemoveContentProtectionJavaScriptInjectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: RemoveContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/content-protection-javascript-injection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: RemoveContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/content-protection-javascript-injection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error deleting match target",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: RemoveContentProtectionJavaScriptInjectionRuleRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: RemoveContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: RemoveContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID: 43253,
				Version:  15,
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing ContentProtectionJavaScriptInjectionRuleID": {
			params: RemoveContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContentProtectionJavaScriptInjectionRuleID")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.RemoveContentProtectionJavaScriptInjectionRule(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
