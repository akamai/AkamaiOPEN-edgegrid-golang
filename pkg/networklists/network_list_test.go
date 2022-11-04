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

func TestNetworkList_ListNetworkList(t *testing.T) {

	result := GetNetworkListsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkLists.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetNetworkListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetNetworkListsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetNetworkListsRequest{},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/network-lists",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetNetworkListsRequest{},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching networklist",
    "status": 500
}`,
			expectedPath: "/network-list/v2/network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching networklist",
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
			result, err := client.GetNetworkLists(
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

func TestNetworkList_FilterNetworkLists(t *testing.T) {

	result := GetNetworkListsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkLists.json"))
	json.Unmarshal([]byte(respData), &result)

	expectedResult := GetNetworkListsResponse{}
	expectedResponseData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkLists_GEO.json"))
	json.Unmarshal([]byte(expectedResponseData), &expectedResult)

	tests := map[string]struct {
		params           GetNetworkListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetNetworkListsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetNetworkListsRequest{Type: "GEO"},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/network-lists",
			expectedResponse: &expectedResult,
		},
		"500 internal server error": {
			params:         GetNetworkListsRequest{},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching networklist",
    "status": 500
}`,
			expectedPath: "/network-list/v2/network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching networklist",
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
			result, err := client.GetNetworkLists(
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

// Test NetworkList
func TestNetworkList_GetNetworkList(t *testing.T) {

	result := GetNetworkListResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkList.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetNetworkListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetNetworkListResponse
		withError        error
	}{
		"200 OK": {
			params:           GetNetworkListRequest{UniqueID: "Test"},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/network-list/v2/network-lists/Test",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         GetNetworkListRequest{UniqueID: "Test"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching networklist"
}`,
			expectedPath: "/network-list/v2/network-lists/Test",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching networklist",
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
			result, err := client.GetNetworkList(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create NetworkList
func TestNetworkList_CreateNetworkList(t *testing.T) {

	result := CreateNetworkListResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkList.json"))
	json.Unmarshal([]byte(respData), &result)

	req := CreateNetworkListRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkList.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           CreateNetworkListRequest
		prop             *CreateNetworkListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateNetworkListResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateNetworkListRequest{Name: "Test"},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/network-list/v2/network-lists",
		},
		"500 internal server error": {
			params:         CreateNetworkListRequest{Name: "Test"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating networklist"
}`,
			expectedPath: "/network-list/v2/network-lists",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating networklist",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, test.expectedPath, r.URL.String())
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateNetworkList(
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

// Test Update NetworkList
func TestNetworkList_UpdateNetworkList(t *testing.T) {
	result := UpdateNetworkListResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkList.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateNetworkListRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkList.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateNetworkListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateNetworkListResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateNetworkListRequest{Name: "TEST", UniqueID: "Test"},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/network-list/v2/network-lists/Test",
		},
		"500 internal server error": {
			params:         UpdateNetworkListRequest{Name: "TEST", UniqueID: "Test"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating networklist"
}`,
			expectedPath: "/network-list/v2/network-lists/Test",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating networklist",
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
			result, err := client.UpdateNetworkList(
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

//Test delete NetworkList
func TestNetworkList_DeleteNetworkList(t *testing.T) {

	result := RemoveNetworkListResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkListEmpty.json"))
	json.Unmarshal([]byte(respData), &result)

	req := RemoveNetworkListRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestNetworkList/NetworkListEmpty.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           RemoveNetworkListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveNetworkListResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveNetworkListRequest{UniqueID: "Test"},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/network-list/v2/network-lists/Test",
		},
		"500 internal server error": {
			params:         RemoveNetworkListRequest{UniqueID: "Test"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
        {
         "type": "internal_error",
         "title": "Internal Server Error",
         "detail": "Error deleting networklist"
         }`,
			expectedPath: "/network-list/v2/network-lists/Test",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting networklist",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, test.expectedPath, r.URL.String())
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.RemoveNetworkList(
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
