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

func TestAppSec_ListRatePolicyActions(t *testing.T) {

	result := GetRatePolicyActionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicyAction/RatePolicyActions.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetRatePolicyActionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRatePolicyActionsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetRatePolicyActionsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rate-policies",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetRatePolicyActionsRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rate-policies",
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
			result, err := client.GetRatePolicyActions(
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

// Test RatePolicyAction
func TestAppSec_GetRatePolicyAction(t *testing.T) {

	result := GetRatePolicyActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicyAction/RatePolicyActions.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetRatePolicyActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRatePolicyActionResponse
		withError        error
	}{
		"200 OK": {
			params: GetRatePolicyActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rate-policies",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetRatePolicyActionRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rate-policies",
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
			result, err := client.GetRatePolicyAction(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update RatePolicyAction.
func TestAppSec_UpdateRatePolicyAction(t *testing.T) {
	result := UpdateRatePolicyActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicyAction/RatePolicyActions.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateRatePolicyActionRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestRatePolicyAction/RatePolicyActions.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateRatePolicyActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateRatePolicyActionResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateRatePolicyActionRequest{
				ConfigID:     43253,
				Version:      15,
				PolicyID:     "AAAA_81230",
				RatePolicyID: 134644,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rate-policies/134644",
		},
		"500 internal server error": {
			params: UpdateRatePolicyActionRequest{
				ConfigID:     43253,
				Version:      15,
				PolicyID:     "AAAA_81230",
				RatePolicyID: 134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/rate-policies/134644",
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
			result, err := client.UpdateRatePolicyAction(
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
