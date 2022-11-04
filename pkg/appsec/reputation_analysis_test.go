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

func TestAppSec_ListReputationAnalysis(t *testing.T) {

	result := GetReputationAnalysisResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationAnalysis/ReputationAnalysis.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetReputationAnalysisRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetReputationAnalysisResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetReputationAnalysisRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-analysis",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetReputationAnalysisRequest{
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
    "detail": "Error fetching ReputationAnalysis",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-analysis",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ReputationAnalysis",
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
			result, err := client.GetReputationAnalysis(
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

// Test ReputationAnalysis
func TestAppSec_GetReputationAnalysis(t *testing.T) {

	result := GetReputationAnalysisResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationAnalysis/ReputationAnalysis.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetReputationAnalysisRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetReputationAnalysisResponse
		withError        error
	}{
		"200 OK": {
			params: GetReputationAnalysisRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-analysis",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetReputationAnalysisRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching ReputationAnalysis"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-analysis",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching ReputationAnalysis",
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
			result, err := client.GetReputationAnalysis(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update ReputationAnalysis
func TestAppSec_UpdateReputationAnalysis(t *testing.T) {
	result := UpdateReputationAnalysisResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationAnalysis/ReputationAnalysis.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateReputationAnalysisRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationAnalysis/ReputationAnalysis.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateReputationAnalysisRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateReputationAnalysisResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateReputationAnalysisRequest{
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
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-analysis",
		},
		"500 internal server error": {
			params: UpdateReputationAnalysisRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating ReputationAnalysis"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-analysis",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating ReputationAnalysis",
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
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateReputationAnalysis(
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

// Test Remove ReputationAnalysis
func TestAppSec_RemoveReputationAnalysis(t *testing.T) {

	result := RemoveReputationAnalysisResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationAnalysis/ReputationAnalysisEmpty.json"))
	json.Unmarshal([]byte(respData), &result)

	req := RemoveReputationAnalysisRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationAnalysis/ReputationAnalysisEmpty.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           RemoveReputationAnalysisRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveReputationAnalysisResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveReputationAnalysisRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-analysis",
		},
		"500 internal server error": {
			params: RemoveReputationAnalysisRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting ReputationAnalysis"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-analysis",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting ReputationAnalysis",
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
			client := mockAPIClient(t, mockServer)
			result, err := client.RemoveReputationAnalysis(
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
