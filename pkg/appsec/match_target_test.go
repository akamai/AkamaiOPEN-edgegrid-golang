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

func TestAppSec_ListMatchTargets(t *testing.T) {

	result := GetMatchTargetsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestMatchTargets/MatchTarget.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetMatchTargetsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetMatchTargetsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetMatchTargetsRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/match-targets",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetMatchTargetsRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/match-targets",
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
			result, err := client.GetMatchTargets(
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

// Test MatchTarget
func TestAppSec_GetMatchTarget(t *testing.T) {

	result := GetMatchTargetResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestMatchTargets/MatchTargets.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetMatchTargetRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetMatchTargetResponse
		withError        error
	}{
		"200 OK": {
			params: GetMatchTargetRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				TargetID:      3008967,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/match-targets/3008967?includeChildObjectName=true",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetMatchTargetRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				TargetID:      3008967,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/match-targets/3008967?includeChildObjectName=true",
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
			result, err := client.GetMatchTarget(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create MatchTarget
func TestAppSec_CreateMatchTarget(t *testing.T) {

	result := CreateMatchTargetResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestMatchTargets/MatchTargets.json"))
	json.Unmarshal([]byte(respData), &result)

	req := CreateMatchTargetRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestMatchTargets/MatchTargets.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           CreateMatchTargetRequest
		prop             *CreateMatchTargetRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateMatchTargetResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateMatchTargetRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/match-targets",
		},
		"500 internal server error": {
			params: CreateMatchTargetRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/match-targets",
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
			result, err := client.CreateMatchTarget(
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

// Test Update MatchTarget
func TestAppSec_UpdateMatchTarget(t *testing.T) {
	result := UpdateMatchTargetResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestMatchTargets/MatchTargets.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateMatchTargetRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestMatchTargets/MatchTargets.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateMatchTargetRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateMatchTargetResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateMatchTargetRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				TargetID:      3008967,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/match-targets/3008967",
		},
		"500 internal server error": {
			params: UpdateMatchTargetRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				TargetID:      3008967,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/match-targets/3008967",
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
			result, err := client.UpdateMatchTarget(
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

// Test Remove MatchTarget
func TestAppSec_RemoveMatchTarget(t *testing.T) {

	result := RemoveMatchTargetResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestMatchTargets/MatchTargetsEmpty.json"))
	json.Unmarshal([]byte(respData), &result)

	req := RemoveMatchTargetRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestMatchTargets/MatchTargetsEmpty.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           RemoveMatchTargetRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveMatchTargetResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveMatchTargetRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				TargetID:      3008967,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/match-targets/3008967",
		},
		"500 internal server error": {
			params: RemoveMatchTargetRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				TargetID:      3008967,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/match-targets/3008967",
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
			result, err := client.RemoveMatchTarget(
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
