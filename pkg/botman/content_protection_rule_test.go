package botman

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

// Test Get ContentProtectionRule List
func TestBotman_GetContentProtectionRuleList(t *testing.T) {

	tests := map[string]struct {
		params           GetContentProtectionRuleListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetContentProtectionRuleListResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetContentProtectionRuleListRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"contentProtectionRules": [
		{"contentProtectionRuleId":"fake3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"contentProtectionRuleId":"fakead64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"contentProtectionRuleId":"fake4ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"contentProtectionRuleId":"faked85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rules",
			expectedResponse: &GetContentProtectionRuleListResponse{
				ContentProtectionRules: []map[string]interface{}{
					{"contentProtectionRuleId": "fake3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
					{"contentProtectionRuleId": "fakead64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
					{"contentProtectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
					{"contentProtectionRuleId": "fake4ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
					{"contentProtectionRuleId": "faked85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
				},
			},
		},
		"200 OK One Record": {
			params: GetContentProtectionRuleListRequest{
				ConfigID:                43253,
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"contentProtectionRules":[
		{"contentProtectionRuleId":"fake3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"contentProtectionRuleId":"fakead64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"contentProtectionRuleId":"fake4ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"contentProtectionRuleId":"faked85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rules",
			expectedResponse: &GetContentProtectionRuleListResponse{
				ContentProtectionRules: []map[string]interface{}{
					{"contentProtectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				},
			},
		},
		"500 internal server error": {
			params: GetContentProtectionRuleListRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rules",
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
			params: GetContentProtectionRuleListRequest{
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
			params: GetContentProtectionRuleListRequest{
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
			params: GetContentProtectionRuleListRequest{
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
			result, err := client.GetContentProtectionRuleList(
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

// Test Get ContentProtectionRule
func TestBotman_GetContentProtectionRule(t *testing.T) {
	tests := map[string]struct {
		params           GetContentProtectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
			expectedResponse: map[string]interface{}{"contentProtectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
		},
		"500 internal server error": {
			params: GetContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: GetContentProtectionRuleRequest{
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetContentProtectionRuleRequest{
				ConfigID:                43253,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: GetContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing ContentProtectionRuleID": {
			params: GetContentProtectionRuleRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				Version:          15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContentProtectionRuleID")
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
			result, err := client.GetContentProtectionRule(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create ContentProtectionRule
func TestBotman_CreateContentProtectionRule(t *testing.T) {

	tests := map[string]struct {
		params           CreateContentProtectionRuleRequest
		prop             *CreateContentProtectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateContentProtectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusCreated,
			responseBody:     `{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"contentProtectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rules",
		},
		"500 internal server error": {
			params: CreateContentProtectionRuleRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rules",
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
			params: CreateContentProtectionRuleRequest{
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
			params: CreateContentProtectionRuleRequest{
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
			params: CreateContentProtectionRuleRequest{
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
			result, err := client.CreateContentProtectionRule(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update ContentProtectionRule
func TestBotman_UpdateContentProtectionRule(t *testing.T) {
	tests := map[string]struct {
		params           UpdateContentProtectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 10,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:             json.RawMessage(`{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"contentProtectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/content-protection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: UpdateContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 10,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:             json.RawMessage(`{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/content-protection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: UpdateContentProtectionRuleRequest{
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:             json.RawMessage(`{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateContentProtectionRuleRequest{
				ConfigID:                43253,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:             json.RawMessage(`{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: UpdateContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:             json.RawMessage(`{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing JsonPayload": {
			params: UpdateContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "JsonPayload")
			},
		},
		"Missing ContentProtectionRuleID": {
			params: UpdateContentProtectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContentProtectionRuleID")
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
			result, err := client.UpdateContentProtectionRule(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Remove ContentProtectionRule
func TestBotman_RemoveContentProtectionRule(t *testing.T) {
	tests := map[string]struct {
		params           RemoveContentProtectionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: RemoveContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 10,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/content-protection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: RemoveContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 10,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/content-protection-rules/fake3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: RemoveContentProtectionRuleRequest{
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: RemoveContentProtectionRuleRequest{
				ConfigID:                43253,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: RemoveContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing ContentProtectionRuleID": {
			params: RemoveContentProtectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContentProtectionRuleID")
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
			err := client.RemoveContentProtectionRule(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
