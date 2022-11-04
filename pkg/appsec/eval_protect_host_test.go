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

func TestAppSec_ListEvalProtectHost(t *testing.T) {

	result := GetEvalProtectHostResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestEvalProtectHost/EvalProtectHost.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetEvalProtectHostRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetEvalProtectHostResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetEvalProtectHostRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/selected-hostnames/eval-hostnames",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetEvalProtectHostRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching EvalProtectHost",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/selected-hostnames/eval-hostnames",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching EvalProtectHost",
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
			result, err := client.GetEvalProtectHost(
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

// Test EvalProtectHost
func TestAppSec_GetEvalProtectHost(t *testing.T) {

	result := GetEvalProtectHostResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestEvalProtectHost/EvalProtectHost.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetEvalProtectHostRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetEvalProtectHostResponse
		withError        error
	}{
		"200 OK": {
			params: GetEvalProtectHostRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/selected-hostnames/eval-hostnames",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetEvalProtectHostRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching EvalProtectHost"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/selected-hostnames/eval-hostnames",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching EvalProtectHost",
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
			result, err := client.GetEvalProtectHost(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update EvalProtectHost.
func TestAppSec_UpdateEvalProtectHost(t *testing.T) {
	result := UpdateEvalProtectHostResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestEvalProtectHost/EvalProtectHost.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateEvalProtectHostRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestEvalProtectHost/EvalProtectHost.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateEvalProtectHostRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateEvalProtectHostResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateEvalProtectHostRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/protect-eval-hostnames",
		},
		"500 internal server error": {
			params: UpdateEvalProtectHostRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating EvalProtectHost"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/protect-eval-hostnames",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating EvalProtectHost",
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
			result, err := client.UpdateEvalProtectHost(
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
