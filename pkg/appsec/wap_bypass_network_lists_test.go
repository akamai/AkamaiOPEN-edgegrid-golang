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

func TestAppSec_ListWAPBypassNetworkLists(t *testing.T) {

	result := GetWAPBypassNetworkListsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestBypassNetworkLists/BypassNetworkLists.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetWAPBypassNetworkListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetWAPBypassNetworkListsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetWAPBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/bypass-network-lists",
			expectedResponse: &result,
		},
		"validation error - missing PolicyID": {
			params: GetWAPBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: ErrStructValidation,
		},
		"401 Not authorized - incorrect credentials": {
			params: GetWAPBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusUnauthorized,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
    "title": "Not authorized",
	"status": 401,
    "detail": "Inactive client token",
    "instance": "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens"
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/bypass-network-lists",
			withError: &Error{
				Type:       "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:      "Not authorized",
				Detail:     "Inactive client token",
				Instance:   "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens",
				StatusCode: 401,
			},
		},
		"500 internal server error": {
			params: GetWAPBypassNetworkListsRequest{
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
    "detail": "Error fetching WAPBypassNetworkLists",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/bypass-network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching WAPBypassNetworkLists",
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
			result, err := client.GetWAPBypassNetworkLists(
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

func TestAppSec_GetWAPBypassNetworkLists(t *testing.T) {

	result := GetWAPBypassNetworkListsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestBypassNetworkLists/BypassNetworkLists.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetWAPBypassNetworkListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetWAPBypassNetworkListsResponse
		withError        error
	}{
		"200 OK": {
			params: GetWAPBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/bypass-network-lists",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetWAPBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching WAPBypassNetworkLists"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/bypass-network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching WAPBypassNetworkLists",
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
			result, err := client.GetWAPBypassNetworkLists(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateWAPBypassNetworkLists(t *testing.T) {
	result := UpdateWAPBypassNetworkListsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestBypassNetworkLists/BypassNetworkLists.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateWAPBypassNetworkListsRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestBypassNetworkLists/BypassNetworkLists.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateWAPBypassNetworkListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateWAPBypassNetworkListsResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateWAPBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/bypass-network-lists",
		},
		"500 internal server error": {
			params: UpdateWAPBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating WAPBypassNetworkLists"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/bypass-network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating WAPBypassNetworkLists",
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
			result, err := client.UpdateWAPBypassNetworkLists(
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
