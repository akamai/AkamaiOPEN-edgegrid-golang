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

func TestAppSec_ListEvalGroup(t *testing.T) {

	result := GetAttackGroupsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAttackGroup/AttackGroups.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAttackGroupsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAttackGroupsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAttackGroupsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-groups?includeConditionException=true",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAttackGroupsRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-groups?includeConditionException=true",
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
			result, err := client.GetEvalGroups(
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

// Test AttackGroupConditionException
func TestAppSec_GetEvalGroup(t *testing.T) {

	result := GetAttackGroupResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAttackGroup/AttackGroup.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetAttackGroupRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAttackGroupResponse
		withError        error
	}{
		"200 OK": {
			params: GetAttackGroupRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Group:    "SQL",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-groups/SQL?includeConditionException=true",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAttackGroupRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Group:    "SQL",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-groups/SQL?includeConditionException=true",
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
			result, err := client.GetEvalGroup(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update AttackGroupConditionException.
func TestAppSec_UpdateEvalGroup(t *testing.T) {
	result := UpdateAttackGroupResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAttackGroup/AttackGroup.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateAttackGroupRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAttackGroup/AttackGroup.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateAttackGroupRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAttackGroupResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAttackGroupRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Group:    "SQL",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-groups/SQL/action-condition-exception",
		},
		"500 internal server error": {
			params: UpdateAttackGroupRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Group:    "SQL",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-groups/SQL/action-condition-exception",
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
			result, err := client.UpdateEvalGroup(
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
