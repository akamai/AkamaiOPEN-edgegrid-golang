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

func TestAppSec_ListReputationProfile(t *testing.T) {

	result := GetReputationProfilesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	json.Unmarshal([]byte(respData), &result)

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
	json.Unmarshal([]byte(respData), &result)

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
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
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
	json.Unmarshal([]byte(respData), &result)

	req := CreateReputationProfileRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           CreateReputationProfileRequest
		prop             *CreateReputationProfileRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateReputationProfileResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateReputationProfileRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/reputation-profiles",
		},
		"500 internal server error": {
			params: CreateReputationProfileRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating domain"
}`),
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
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
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
	json.Unmarshal([]byte(respData), &result)

	req := UpdateReputationProfileRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfile.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateReputationProfileRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateReputationProfileResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateReputationProfileRequest{
				ConfigID:            43253,
				ConfigVersion:       15,
				ReputationProfileId: 134644,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/reputation-profiles/134644",
		},
		"500 internal server error": {
			params: UpdateReputationProfileRequest{
				ConfigID:            43253,
				ConfigVersion:       15,
				ReputationProfileId: 134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
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
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
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
	json.Unmarshal([]byte(respData), &result)

	req := RemoveReputationProfileRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationProfile/ReputationProfileEmpty.json"))
	json.Unmarshal([]byte(reqData), &req)

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
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting match target"
}`),
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
