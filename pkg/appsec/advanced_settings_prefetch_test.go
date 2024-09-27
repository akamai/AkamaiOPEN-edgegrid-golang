package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListAdvancedSettingsPrefetch(t *testing.T) {

	result := GetAdvancedSettingsPrefetchResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsPrefetchRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsPrefetchResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAdvancedSettingsPrefetchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/prefetch",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsPrefetchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching AdvancedSettingsPrefetch",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/prefetch",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsPrefetch",
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
			result, err := client.GetAdvancedSettingsPrefetch(
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

// Test AdvancedSettingsPrefetch
func TestAppSec_GetAdvancedSettingsPrefetch(t *testing.T) {

	result := GetAdvancedSettingsPrefetchResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsPrefetchRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsPrefetchResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsPrefetchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/prefetch",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsPrefetchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching AdvancedSettingsPrefetch"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/prefetch",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsPrefetch",
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
			result, err := client.GetAdvancedSettingsPrefetch(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update AdvancedSettingsPrefetch.
func TestAppSec_UpdateAdvancedSettingsPrefetch(t *testing.T) {
	result := UpdateAdvancedSettingsPrefetchResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsPrefetchRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsPrefetchRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsPrefetchResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsPrefetchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/prefetch",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsPrefetchRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsPrefetch"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/prefetch",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsPrefetch",
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
			result, err := client.UpdateAdvancedSettingsPrefetch(
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
