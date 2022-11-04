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

func TestAppSec_ListMatchTargetSequence(t *testing.T) {

	result := GetMatchTargetSequenceResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestMatchTargetSequence/MatchTargetSequence.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetMatchTargetSequenceRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetMatchTargetSequenceResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetMatchTargetSequenceRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Type:          "website",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/match-targets/sequence?type=website",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetMatchTargetSequenceRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Type:          "website",
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/match-targets/sequence?type=website",
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
			result, err := client.GetMatchTargetSequence(
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

// Test MatchTargetSequence
func TestAppSec_GetMatchTargetSequence(t *testing.T) {

	result := GetMatchTargetSequenceResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestMatchTargetSequence/MatchTargetSequence.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetMatchTargetSequenceRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetMatchTargetSequenceResponse
		withError        error
	}{
		"200 OK": {
			params: GetMatchTargetSequenceRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Type:          "website",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/match-targets/sequence?type=website",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetMatchTargetSequenceRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Type:          "website",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/match-targets/sequence?type=website",
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
			result, err := client.GetMatchTargetSequence(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update MatchTargetSequence.
func TestAppSec_UpdateMatchTargetSequence(t *testing.T) {
	result := UpdateMatchTargetSequenceResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestMatchTargetSequence/MatchTargetSequence.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateMatchTargetSequenceRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestMatchTargetSequence/MatchTargetSequence.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateMatchTargetSequenceRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateMatchTargetSequenceResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateMatchTargetSequenceRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Type:          "website",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/match-targets/%d",
		},
		"500 internal server error": {
			params: UpdateMatchTargetSequenceRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				Type:          "website",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/match-targets/%d",
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
			result, err := client.UpdateMatchTargetSequence(
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
