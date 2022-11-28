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

func TestAppSec_ListCustomDeny(t *testing.T) {

	result := GetCustomDenyListResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomDeny/CustomDenyList.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetCustomDenyListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCustomDenyListResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetCustomDenyListRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/custom-deny",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetCustomDenyListRequest{
				ConfigID: 43253,
				Version:  15,
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/custom-deny",
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
			result, err := client.GetCustomDenyList(
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

// Test CustomDeny
func TestAppSec_GetCustomDeny(t *testing.T) {

	result := GetCustomDenyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomDeny/CustomDeny.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetCustomDenyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCustomDenyResponse
		withError        error
	}{
		"200 OK": {
			params: GetCustomDenyRequest{
				ConfigID: 43253,
				Version:  15,
				ID:       "622919",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/custom-deny/622919",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetCustomDenyRequest{
				ConfigID: 43253,
				Version:  15,
				ID:       "622919",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching CustomDeny"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/custom-deny/622919",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching CustomDeny",
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
			result, err := client.GetCustomDeny(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create CustomDeny
func TestAppSec_CreateCustomDeny(t *testing.T) {

	result := CreateCustomDenyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomDeny/CustomDeny.json"))
	json.Unmarshal([]byte(respData), &result)

	req := CreateCustomDenyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomDeny/CustomDeny.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           CreateCustomDenyRequest
		prop             *CreateCustomDenyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateCustomDenyResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateCustomDenyRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/custom-deny",
		},
		"500 internal server error": {
			params: CreateCustomDenyRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating CustomDeny"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/custom-deny",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating CustomDeny",
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
			result, err := client.CreateCustomDeny(
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

// Test Update CustomDeny
func TestAppSec_UpdateCustomDeny(t *testing.T) {
	result := UpdateCustomDenyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomDeny/CustomDeny.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateCustomDenyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomDeny/CustomDeny.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateCustomDenyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateCustomDenyResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateCustomDenyRequest{
				ConfigID: 43253,
				Version:  15,
				ID:       "deny_custom_622918",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/custom-deny/deny_custom_622918",
		},
		"500 internal server error": {
			params: UpdateCustomDenyRequest{
				ConfigID: 43253,
				Version:  15,
				ID:       "deny_custom_622918",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating CustomDeny"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/custom-deny/deny_custom_622918",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating CustomDeny",
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
			result, err := client.UpdateCustomDeny(
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

// Test Remove CustomDeny
func TestAppSec_RemoveCustomDeny(t *testing.T) {

	result := RemoveCustomDenyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestCustomDeny/CustomDenyEmpty.json"))
	json.Unmarshal([]byte(respData), &result)

	req := RemoveCustomDenyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestCustomDeny/CustomDenyEmpty.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           RemoveCustomDenyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveCustomDenyResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveCustomDenyRequest{
				ConfigID: 43253,
				Version:  15,
				ID:       "deny_custom_622918",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/custom-deny/deny_custom_622918",
		},
		"500 internal server error": {
			params: RemoveCustomDenyRequest{
				ConfigID: 43253,
				Version:  15,
				ID:       "deny_custom_622918",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting CustomDeny"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/custom-deny/deny_custom_622918",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting CustomDeny",
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
			result, err := client.RemoveCustomDeny(
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
