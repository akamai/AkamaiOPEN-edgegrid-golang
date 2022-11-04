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

func TestAppSec_ListApiRequestConstraints(t *testing.T) {

	result := GetApiRequestConstraintsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiRequestConstraints/ApiRequestConstraints.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetApiRequestConstraintsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetApiRequestConstraintsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetApiRequestConstraintsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-request-constraints",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetApiRequestConstraintsRequest{
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
    "detail": "Error fetching ApiRequestConstraints",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-request-constraints",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ApiRequestConstraints",
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
			result, err := client.GetApiRequestConstraints(
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

// Test ApiRequestConstraints
func TestAppSec_GetApiRequestConstraints(t *testing.T) {

	result := GetApiRequestConstraintsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiRequestConstraints/ApiRequestConstraints.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetApiRequestConstraintsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetApiRequestConstraintsResponse
		withError        error
	}{
		"200 OK": {
			params: GetApiRequestConstraintsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-request-constraints",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetApiRequestConstraintsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching ApiRequestConstraints"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-request-constraints",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ApiRequestConstraints",
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
			result, err := client.GetApiRequestConstraints(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update ApiRequestConstraints.
func TestAppSec_UpdateApiRequestConstraints(t *testing.T) {
	result := UpdateApiRequestConstraintsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiRequestConstraints/ApiRequestConstraints.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateApiRequestConstraintsRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestApiRequestConstraints/ApiRequestConstraints.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateApiRequestConstraintsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateApiRequestConstraintsResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateApiRequestConstraintsRequest{
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
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-request-constraints",
		},
		"500 internal server error": {
			params: UpdateApiRequestConstraintsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating ApiRequestConstraints"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-request-constraints",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating ApiRequestConstraints",
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
			result, err := client.UpdateApiRequestConstraints(
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
