package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApsec_ListAdvancedSettingsAsePenaltyBox(t *testing.T) {

	result := GetAdvancedSettingsAsePenaltyBoxResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsAsePenaltyBoxRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsAsePenaltyBoxResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAdvancedSettingsAsePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/ase-penalty-box",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsAsePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching AdvancedSettingsAsePenaltyBox",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/ase-penalty-box",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsAsePenaltyBox",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAdvancedSettingsAsePenaltyBox(
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

// Test AdvancedSettingsAsePenaltyBox
func TestAppSec_GetAdvancedSettingsAsePenaltyBox(t *testing.T) {

	result := GetAdvancedSettingsAsePenaltyBoxResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsAsePenaltyBoxRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsAsePenaltyBoxResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsAsePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/ase-penalty-box",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsAsePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching AdvancedSettingsAsePenaltyBox"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/ase-penalty-box",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsAsePenaltyBox",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAdvancedSettingsAsePenaltyBox(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateAdvancedSettingsAsePenaltyBox(t *testing.T) {
	result := UpdateAdvancedSettingsAsePenaltyBoxResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsAsePenaltyBoxRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsAsePenaltyBoxRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsAsePenaltyBoxResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsAsePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/ase-penalty-box",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsAsePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsAsePenaltyBox"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/ase-penalty-box",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsAsePenaltyBox",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAdvancedSettingsAsePenaltyBox(
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

func TestAppSec_RemoveAdvancedSettingsAsePenaltyBox(t *testing.T) {
	result := RemoveAdvancedSettingsAsePenaltyBoxResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := RemoveAdvancedSettingsAsePenaltyBoxRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           RemoveAdvancedSettingsAsePenaltyBoxRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveAdvancedSettingsAsePenaltyBoxResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveAdvancedSettingsAsePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/ase-penalty-box",
		},
		"500 internal server error": {
			params: RemoveAdvancedSettingsAsePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsAsePenaltyBox"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/ase-penalty-box",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsAsePenaltyBox",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.RemoveAdvancedSettingsAsePenaltyBox(
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
