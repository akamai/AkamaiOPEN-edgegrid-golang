package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListConfigurationClone(t *testing.T) {

	result := GetConfigurationCloneResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestConfigurationClone/ConfigurationClone.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetConfigurationCloneRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetConfigurationCloneResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetConfigurationCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetConfigurationCloneRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15",
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
			result, err := client.GetConfigurationClone(
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

// Test ConfigurationClone
func TestAppSec_GetConfigurationClone(t *testing.T) {

	result := GetConfigurationCloneResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestConfigurationClone/ConfigurationClone.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetConfigurationCloneRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetConfigurationCloneResponse
		withError        error
	}{
		"200 OK": {
			params: GetConfigurationCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetConfigurationCloneRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching ConfigurationClone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ConfigurationClone",
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
			result, err := client.GetConfigurationClone(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create ConfigurationClone
func TestAppSec_CreateConfigurationClone(t *testing.T) {

	result := CreateConfigurationCloneResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestConfigurationClone/ConfigurationClone.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := CreateConfigurationCloneRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestConfigurationClone/ConfigurationClone.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           CreateConfigurationCloneRequest
		prop             *CreateConfigurationCloneRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateConfigurationCloneResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateConfigurationCloneRequest{Name: "Test", CreateFrom: struct {
				ConfigID int `json:"configId"`
				Version  int `json:"version"`
			}{ConfigID: 42345,
				Version: 7}},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/",
		},
		"500 internal server error": {
			params: CreateConfigurationCloneRequest{Name: "Test", CreateFrom: struct {
				ConfigID int `json:"configId"`
				Version  int `json:"version"`
			}{ConfigID: 42345,
				Version: 7}},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating ConfigurationClone"
			}`,
			expectedPath: "/appsec/v1/configs/",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating ConfigurationClone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateConfigurationClone(
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
