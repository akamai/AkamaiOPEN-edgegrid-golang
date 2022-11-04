package networklists

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

func TestApsec_ListNetworkListDescription(t *testing.T) {

	result := GetNetworkListDescriptionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkListDescription/NetworkListDescription.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetNetworkListDescriptionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetNetworkListDescriptionResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetNetworkListDescriptionRequest{UniqueID: "Test"},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/network-lists/Test",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetNetworkListDescriptionRequest{UniqueID: "Test"},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching NetworkListDescription",
    "status": 500
}`,
			expectedPath: "/network-list/v2/network-lists/Test",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching NetworkListDescription",
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
			result, err := client.GetNetworkListDescription(
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

// Test NetworkListDescription
func TestAppSec_GetNetworkListDescription(t *testing.T) {

	result := GetNetworkListDescriptionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkListDescription/NetworkListDescription.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetNetworkListDescriptionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetNetworkListDescriptionResponse
		withError        error
	}{
		"200 OK": {
			params:           GetNetworkListDescriptionRequest{UniqueID: "Test"},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/network-lists/Test",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetNetworkListDescriptionRequest{UniqueID: "Test"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching NetworkListDescription"
}`,
			expectedPath: "/network-list/v2/network-lists/Test",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching NetworkListDescription",
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
			result, err := client.GetNetworkListDescription(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update NetworkListDescription.
func TestAppSec_UpdateNetworkListDescription(t *testing.T) {
	result := UpdateNetworkListDescriptionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkListDescription/NetworkListDescription.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateNetworkListDescriptionRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestNetworkListDescription/NetworkListDescription.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateNetworkListDescriptionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateNetworkListDescriptionResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateNetworkListDescriptionRequest{UniqueID: "Test"},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/network-list/v2/network-lists/Test/details",
		},
		"500 internal server error": {
			params:         UpdateNetworkListDescriptionRequest{UniqueID: "Test"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating NetworkListDescription"
}`,
			expectedPath: "/network-list/v2/network-lists/Test/details",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating NetworkListDescription",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, test.expectedPath, r.URL.String())
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateNetworkListDescription(
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
