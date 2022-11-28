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

func TestAppSec_ListCustomRuleActions(t *testing.T) {

	result := GetCustomRuleActionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRuleAction/CustomRuleActions.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetCustomRuleActionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCustomRuleActionsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetCustomRuleActionsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/custom-rules",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetCustomRuleActionsRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/custom-rules",
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
			result, err := client.GetCustomRuleActions(
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

// Test CustomRuleAction
func TestAppSec_GetCustomRuleAction(t *testing.T) {

	result := GetCustomRuleActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRuleAction/CustomRuleActions.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetCustomRuleActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCustomRuleActionResponse
		withError        error
	}{
		"200 OK": {
			params: GetCustomRuleActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				RuleID:   12345,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/custom-rules",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetCustomRuleActionRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/custom-rules",
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
			result, err := client.GetCustomRuleAction(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update CustomRuleAction.
func TestAppSec_UpdateCustomRuleAction(t *testing.T) {
	result := UpdateCustomRuleActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomRuleAction/CustomRuleActionUpdate.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateCustomRuleActionRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomRuleAction/CustomRuleAction.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateCustomRuleActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateCustomRuleActionResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateCustomRuleActionRequest{
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
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/custom-rules/12345",
		},
		"500 internal server error": {
			params: UpdateCustomRuleActionRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/custom-rules/12345",
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
			result, err := client.UpdateCustomRuleAction(
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
