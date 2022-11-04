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

func TestAppSec_ListVersionNotes(t *testing.T) {

	result := GetVersionNotesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestVersionNotes/VersionNotes.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetVersionNotesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetVersionNotesResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetVersionNotesRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/version-notes",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetVersionNotesRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching VersionNotes",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/version-notes",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching VersionNotes",
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
			result, err := client.GetVersionNotes(
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

// Test VersionNotes
func TestAppSec_GetVersionNotes(t *testing.T) {

	result := GetVersionNotesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestVersionNotes/VersionNotes.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetVersionNotesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetVersionNotesResponse
		withError        error
	}{
		"200 OK": {
			params: GetVersionNotesRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/version-notes",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetVersionNotesRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching VersionNotes"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/version-notes",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching VersionNotes",
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
			result, err := client.GetVersionNotes(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update VersionNotes.
func TestAppSec_UpdateVersionNotes(t *testing.T) {
	result := UpdateVersionNotesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestVersionNotes/VersionNotes.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateVersionNotesRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestVersionNotes/VersionNotes.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateVersionNotesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateVersionNotesResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateVersionNotesRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/version-notes",
		},
		"500 internal server error": {
			params: UpdateVersionNotesRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating VersionNotes"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/version-notes",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating VersionNotes",
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
			result, err := client.UpdateVersionNotes(
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
