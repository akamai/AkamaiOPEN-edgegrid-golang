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

func TestAppSec_ListBypassNetworkLists(t *testing.T) {

	result := GetBypassNetworkListsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestBypassNetworkLists/BypassNetworkLists.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetBypassNetworkListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetBypassNetworkListsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/bypass-network-lists",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching BypassNetworkLists",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/bypass-network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching BypassNetworkLists",
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
			result, err := client.GetBypassNetworkLists(
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

// Test BypassNetworkLists
func TestAppSec_GetBypassNetworkLists(t *testing.T) {

	result := GetBypassNetworkListsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestBypassNetworkLists/BypassNetworkLists.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetBypassNetworkListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetBypassNetworkListsResponse
		withError        error
	}{
		"200 OK": {
			params: GetBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/bypass-network-lists",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching BypassNetworkLists"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/bypass-network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching BypassNetworkLists",
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
			result, err := client.GetBypassNetworkLists(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update BypassNetworkLists.
func TestAppSec_UpdateBypassNetworkLists(t *testing.T) {
	result := UpdateBypassNetworkListsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestBypassNetworkLists/BypassNetworkLists.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateBypassNetworkListsRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestBypassNetworkLists/BypassNetworkLists.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateBypassNetworkListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateBypassNetworkListsResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/bypass-network-lists",
		},
		"500 internal server error": {
			params: UpdateBypassNetworkListsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating BypassNetworkLists"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/bypass-network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating BypassNetworkLists",
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
			result, err := client.UpdateBypassNetworkLists(
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
