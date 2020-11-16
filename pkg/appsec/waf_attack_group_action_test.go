package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApsec_ListWAFAttackGroupAction(t *testing.T) {

	result := GetWAFAttackGroupActionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestWAFAttackGroupAction/WAFAttackGroupActions.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetWAFAttackGroupActionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetWAFAttackGroupActionsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetWAFAttackGroupActionsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/attack-groups/",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetWAFAttackGroupActionsRequest{
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
    "detail": "Error fetching attack group actions",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/attack-groups/",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching attack group actions",
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
			result, err := client.GetWAFAttackGroupActions(
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

// Test WAFAttackGroupAction
func TestAppSec_GetWAFAttackGroupAction(t *testing.T) {

	result := GetWAFAttackGroupActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestWAFAttackGroupAction/WAFAttackGroupAction.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetWAFAttackGroupActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetWAFAttackGroupActionResponse
		withError        error
	}{
		"200 OK": {
			params: GetWAFAttackGroupActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Group:    "SQL",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/attack-groups/SQL",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetWAFAttackGroupActionRequest{
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
    "detail": "Error fetching attack group action"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/attack-groups/SQL",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching attack group action",
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
			result, err := client.GetWAFAttackGroupAction(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update WAFAttackGroupAction.
func TestAppSec_UpdateWAFAttackGroupAction(t *testing.T) {
	result := UpdateWAFAttackGroupActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestWAFAttackGroupAction/WAFAttackGroupAction.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateWAFAttackGroupActionRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestWAFAttackGroupAction/WAFAttackGroupActionReq.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateWAFAttackGroupActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateWAFAttackGroupActionResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateWAFAttackGroupActionRequest{
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
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/attack-groups/%s",
		},
		"500 internal server error": {
			params: UpdateWAFAttackGroupActionRequest{
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
    "detail": "Error updating attack group action"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/attack-groups/%s",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating attack group action",
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
			result, err := client.UpdateWAFAttackGroupAction(
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
