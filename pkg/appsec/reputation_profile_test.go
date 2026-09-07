package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListReputationProfile(t *testing.T) {

	result := GetReputationProfilesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetReputationProfilesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetReputationProfilesResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetReputationProfilesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/reputation-profiles",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetReputationProfilesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/reputation-profiles",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetReputationProfiles(
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

// Test ReputationProfile
func TestAppSec_GetReputationProfile(t *testing.T) {

	result := GetReputationProfileResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetReputationProfileRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetReputationProfileResponse
		withError        error
	}{
		"200 OK": {
			params: GetReputationProfileRequest{
				ConfigID:            43253,
				ConfigVersion:       15,
				ReputationProfileId: 134644,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/reputation-profiles/134644",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetReputationProfileRequest{
				ConfigID:            43253,
				ConfigVersion:       15,
				ReputationProfileId: 134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/reputation-profiles/134644",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetReputationProfile(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create ReputationProfile
func TestAppSec_CreateReputationProfile(t *testing.T) {

	result := CreateReputationProfileResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	req := CreateReputationProfileRequest{
		ConfigID:       43253,
		ConfigVersion:  15,
		JsonPayloadRaw: json.RawMessage(reqData),
	}

	tests := map[string]struct {
		params              CreateReputationProfileRequest
		prop                *CreateReputationProfileRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateReputationProfileResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{
		"201 Created": {
			params:              req,
			expectedRequestBody: reqData,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/reputation-profiles",
		},
		"500 internal server error": {
			params:         req,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating domain"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/reputation-profiles",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating domain",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateReputationProfile(
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

// Test Update ReputationProfile
func TestAppSec_UpdateReputationProfile(t *testing.T) {
	result := UpdateReputationProfileResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	req := UpdateReputationProfileRequest{
		ConfigID:            43253,
		ConfigVersion:       15,
		ReputationProfileId: 134644,
		JsonPayloadRaw:      json.RawMessage(reqData),
	}

	tests := map[string]struct {
		params              UpdateReputationProfileRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *UpdateReputationProfileResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{
		"200 Success": {
			params: req,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			expectedRequestBody: reqData,
			responseStatus:      http.StatusCreated,
			responseBody:        respData,
			expectedResponse:    &result,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/reputation-profiles/134644",
		},
		"500 internal server error": {
			params:         req,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/reputation-profiles/134644",
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateReputationProfile(
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

// Test Remove ReputationProfile
func TestAppSec_RemoveReputationProfile(t *testing.T) {

	result := RemoveReputationProfileResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfileEmpty.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := RemoveReputationProfileRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfileEmpty.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           RemoveReputationProfileRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveReputationProfileResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveReputationProfileRequest{
				ConfigID:            43253,
				ConfigVersion:       15,
				ReputationProfileId: 134644,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/rate-policies/134644",
		},
		"500 internal server error": {
			params: RemoveReputationProfileRequest{
				ConfigID:            43253,
				ConfigVersion:       15,
				ReputationProfileId: 134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/rate-policies/134644",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting match target",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.RemoveReputationProfile(
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
