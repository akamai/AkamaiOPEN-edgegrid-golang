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

func TestAppSec_ListCustomRules(t *testing.T) {

	result := GetCustomRulesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRules.json"))
	json.Unmarshal([]byte(respData), &result)

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
	json.Unmarshal([]byte(respData), &result)

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
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
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
	json.Unmarshal([]byte(respData), &result)

	req := CreateCustomRuleRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRule.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           CreateCustomRuleRequest
		prop             *CreateCustomRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateCustomRuleResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateCustomRuleRequest{
				ConfigID: 43253,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/custom-rules",
		},
		"500 internal server error": {
			params: CreateCustomRuleRequest{
				ConfigID: 43253,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating domain"
}`),
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
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
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
	json.Unmarshal([]byte(respData), &result)

	req := UpdateCustomRuleRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRule.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateCustomRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateCustomRuleResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateCustomRuleRequest{
				ConfigID: 43253,
				ID:       60039625,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/custom-rules/%d",
		},
		"500 internal server error": {
			params: UpdateCustomRuleRequest{
				ConfigID: 43253,
				ID:       60039625,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
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
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
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
	json.Unmarshal([]byte(respData), &result)

	req := RemoveCustomRuleRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomRules/CustomRulesEmpty.json"))
	json.Unmarshal([]byte(reqData), &req)

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
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting match target"
}`),
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
