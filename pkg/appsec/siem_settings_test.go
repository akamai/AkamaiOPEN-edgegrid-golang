package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListSiemSettings(t *testing.T) {

	result := GetSiemSettingsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSiemSettings/SiemSettings.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetSiemSettingsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSiemSettingsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetSiemSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/siem",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetSiemSettingsRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/siem",
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
			result, err := client.GetSiemSettings(
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

// Test SiemSettings
func TestAppSec_GetSiemSettings(t *testing.T) {

	result := GetSiemSettingsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSiemSettings/SiemSettings.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetSiemSettingsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSiemSettingsResponse
		withError        error
	}{
		"200 OK": {
			params: GetSiemSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/siem",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetSiemSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/siem",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching match target",
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
			result, err := client.GetSiemSettings(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update SiemSettings.
func TestAppSec_UpdateSiemSettings(t *testing.T) {
	result := UpdateSiemSettingsResponse{}
	resultWithoutEnabledBotman := UpdateSiemSettingsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSiemSettings/SiemSettings.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	respDataWithoutEnableBotman := compactJSON(loadFixtureBytes("testdata/TestSiemSettings/SiemSettingsWithoutEnabledBotmanSiem.json"))
	err = json.Unmarshal([]byte(respDataWithoutEnableBotman), &resultWithoutEnabledBotman)
	require.NoError(t, err)

	req := UpdateSiemSettingsRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestSiemSettings/SiemSettings.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateSiemSettingsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateSiemSettingsResponse
		withError        error
		headers          http.Header
		errors           *regexp.Regexp
	}{
		"200 Success": {
			params: UpdateSiemSettingsRequest{
				ConfigID:                43253,
				Version:                 15,
				EnableSiem:              true,
				EnabledBotmanSiemEvents: ptr.To(false),
				Exceptions: []Exception{
					{
						ActionTypes: []string{"*"},
						Protection:  "botmanagement",
					},
					{
						ActionTypes: []string{"deny"},
						Protection:  "ipgeo",
					},
					{
						ActionTypes: []string{"alert"},
						Protection:  "rate",
					},
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/siem",
		},
		"200 Success without EnabledBotmanSiemEvents": {
			params: UpdateSiemSettingsRequest{
				ConfigID:   43253,
				Version:    15,
				EnableSiem: true,
				Exceptions: []Exception{
					{
						ActionTypes: []string{"*"},
						Protection:  "botmanagement",
					},
					{
						ActionTypes: []string{"deny"},
						Protection:  "ipgeo",
					},
					{
						ActionTypes: []string{"alert"},
						Protection:  "rate",
					},
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respDataWithoutEnableBotman,
			expectedResponse: &resultWithoutEnabledBotman,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/siem",
		},
		"400 Bad Request action types": {
			params: UpdateSiemSettingsRequest{
				ConfigID:   43253,
				Version:    15,
				EnableSiem: true,
				Exceptions: []Exception{
					{
						ActionTypes: []string{"reject"},
						Protection:  "botmanagement",
					},
					{
						ActionTypes: []string{"deny"},
						Protection:  "ipgeo",
					},
					{
						ActionTypes: []string{"alert"},
						Protection:  "rate",
					},
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusBadRequest,
			errors:         regexp.MustCompile(`struct validation: ActionTypes:.+`),
			expectedPath:   "/appsec/v1/configs/43253/versions/15/siem",
		},
		"400 Bad Request protection": {
			params: UpdateSiemSettingsRequest{
				ConfigID:   43253,
				Version:    15,
				EnableSiem: true,
				Exceptions: []Exception{
					{
						ActionTypes: []string{"tarpit"},
						Protection:  "bot",
					},
				},
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusBadRequest,
			errors:         regexp.MustCompile(`struct validation: Protection:.+`),
			expectedPath:   "/appsec/v1/configs/43253/versions/15/siem",
		},
		"500 internal server error": {
			params: UpdateSiemSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/siem",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
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
			result, err := client.UpdateSiemSettings(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)

			if test.errors != nil {
				require.Error(t, err)
				assert.Regexp(t, test.errors, err.Error())
				return
			}

			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
