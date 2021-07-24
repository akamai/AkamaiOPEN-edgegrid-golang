package siteshield

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSiteShield_ListSiteShieldMaps(t *testing.T) {

	result := GetSiteShieldMapsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSiteShield/SiteShieldMaps.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSiteShieldMapsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/siteshield/v1/maps",
			expectedResponse: &result,
		},
		"500 internal server error": {
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching siteshieldmap",
    "status": 500
}`,
			expectedPath: "/siteshield/v1/maps",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching siteshieldmap",
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
			result, err := client.GetSiteShieldMaps(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
			)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test SiteShieldMap
func TestSiteShield_GetSiteShieldMap(t *testing.T) {

	result := SiteShieldMapResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSiteShield/SiteShieldMap.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           SiteShieldMapRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *SiteShieldMapResponse
		withError        error
	}{
		"200 OK": {
			params:           SiteShieldMapRequest{UniqueID: 1234},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/siteshield/v1/maps/1234",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         SiteShieldMapRequest{UniqueID: 1234},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching siteshieldmap"
}`,
			expectedPath: "/siteshield/v1/maps/1234",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching siteshieldmap",
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
			result, err := client.GetSiteShieldMap(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Acknowledgement SiteShieldMap
func TestSiteShield_Acknowledgement(t *testing.T) {

	result := SiteShieldMapResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestSiteShield/SiteShieldMap.json"))
	json.Unmarshal([]byte(respData), &result)

	req := SiteShieldMapRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestSiteShield/SiteShieldMap.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           SiteShieldMapRequest
		prop             *SiteShieldMapRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *SiteShieldMapResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params:           SiteShieldMapRequest{UniqueID: 1234},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/siteshield/v1/maps/1234/acknowledge",
		},
		"500 internal server error": {
			params:         SiteShieldMapRequest{UniqueID: 1234},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating siteshieldmap"
}`,
			expectedPath: "/siteshield/v1/maps/1234/acknowledge",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating siteshieldmap",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, test.expectedPath, r.URL.String())
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.AckSiteShieldMap(
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
