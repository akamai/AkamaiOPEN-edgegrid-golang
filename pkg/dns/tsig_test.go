package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDNS_ListTSIGKeys(t *testing.T) {
	tests := map[string]struct {
		params           ListTSIGKeysRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListTSIGKeysResponse
		withError        error
	}{
		"200 OK": {
			params: ListTSIGKeysRequest{
				TsigQuery: &TSIGQueryString{ContractIDs: []string{"1-1ABCD"}},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"metadata": {
					"totalElements": 1
				},
				"keys": [
					{
						"name": "a.test.key.",
						"algorithm": "hmac-sha256",
						"secret": "DjY16JfIi3JnSDosQWE7Xkx60MbCLo1K7hUCqng8ccg=",
						"zonesCount": 3
					}
				]
			}`,
			expectedPath: "/config-dns/v2/keys?contractIds=1-1ABCD",
			expectedResponse: &ListTSIGKeysResponse{
				Metadata: &TSIGReportMeta{
					TotalElements: 1,
				},
				Keys: []TSIGKeyResponse{
					{
						TSIGKey: TSIGKey{
							Name:      "a.test.key.",
							Algorithm: "hmac-sha256",
							Secret:    "DjY16JfIi3JnSDosQWE7Xkx60MbCLo1K7hUCqng8ccg=",
						},
						ZoneCount: 3,
					},
				},
			},
		},
		"500 internal server error": {
			params: ListTSIGKeysRequest{
				TsigQuery: &TSIGQueryString{ContractIDs: []string{"1-1ABCD"}},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/keys?contractIds=1-1ABCD",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
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
			result, err := client.ListTSIGKeys(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetTSIGKeyZones(t *testing.T) {
	tests := map[string]struct {
		params           GetTSIGKeyZonesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetTSIGKeyZonesResponse
		withError        error
	}{
		"200 OK": {
			params: GetTSIGKeyZonesRequest{
				TsigKey: &TSIGKey{
					Name:      "example.com.akamai.com.",
					Algorithm: "hmac-sha512",
					Secret:    "fakeR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"zones": [
					"river.com",
					"stream.com"
				]
			}`,
			expectedPath: "/config-dns/v2/keys/used-by",
			expectedResponse: &GetTSIGKeyZonesResponse{
				Zones: []string{"river.com", "stream.com"},
			},
		},
		"500 internal server error": {
			params: GetTSIGKeyZonesRequest{
				TsigKey: &TSIGKey{
					Name:      "example.com.akamai.com.",
					Algorithm: "hmac-sha512",
					Secret:    "fakeR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/keys/used-by",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
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
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetTSIGKeyZones(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetTSIGKeyAliases(t *testing.T) {
	tests := map[string]struct {
		params           GetTSIGKeyAliasesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetTSIGKeyAliasesResponse
		withError        error
	}{
		"200 OK": {
			params: GetTSIGKeyAliasesRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"aliases": [
					"exmaple.com",
					"river.com",
					"brook.com",
					"ocean.com"
				]
			}`,
			expectedPath: "/config-dns/v2/zones/example.com/key/used-by",
			expectedResponse: &GetTSIGKeyAliasesResponse{
				Aliases: []string{"exmaple.com", "river.com", "brook.com", "ocean.com"},
			},
		},
		"500 internal server error": {
			params: GetTSIGKeyAliasesRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/key/used-by",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
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
			result, err := client.GetTSIGKeyAliases(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_TSIGKeyBulkUpdate(t *testing.T) {
	tests := map[string]struct {
		params           UpdateTSIGKeyBulkRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneNameListResponse
		withError        error
	}{
		"200 OK": {
			params: UpdateTSIGKeyBulkRequest{
				TSIGKeyBulk: &TSIGKeyBulkPost{
					Key: &TSIGKey{
						Name:      "brook.com.akamai.com.",
						Algorithm: "hmac-sha512",
						Secret:    "fakeR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
					},
					Zones: []string{"brook.com", "river.com"},
				},
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/keys/bulk-update",
			expectedResponse: &ZoneNameListResponse{
				Aliases: []string{"exmaple.com", "river.com", "brook.com", "ocean.com"},
			},
		},
		"500 internal server error": {
			params: UpdateTSIGKeyBulkRequest{
				TSIGKeyBulk: &TSIGKeyBulkPost{
					Key: &TSIGKey{
						Name:      "brook.com.akamai.com.",
						Algorithm: "hmac-sha512",
						Secret:    "fakeR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
					},
					Zones: []string{"brook.com", "river.com"},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/keys/bulk-update",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
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
			err := client.UpdateTSIGKeyBulk(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDNS_GetTSIGKey(t *testing.T) {
	tests := map[string]struct {
		params           GetTSIGKeyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetTSIGKeyResponse
		withError        error
	}{
		"200 OK": {
			params: GetTSIGKeyRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"name": "example.com.akamai.com.",
				"algorithm": "hmac-sha512",
				"secret": "fakeR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				"zonesCount": 7
			}`,
			expectedPath: "/config-dns/v2/zones/example.com/key",
			expectedResponse: &GetTSIGKeyResponse{
				TSIGKey: TSIGKey{
					Name:      "example.com.akamai.com.",
					Algorithm: "hmac-sha512",
					Secret:    "fakeR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				},
				ZoneCount: 7,
			},
		},
		"500 internal server error": {
			params: GetTSIGKeyRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/key",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
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
			result, err := client.GetTSIGKey(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_DeleteTSIGKey(t *testing.T) {
	tests := map[string]struct {
		params           DeleteTSIGKeyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *TSIGKeyResponse
		withError        error
	}{
		"204 No Content": {
			params: DeleteTSIGKeyRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/key",
		},
		"500 internal server error": {
			params: DeleteTSIGKeyRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error Deleting TSig Key",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/key",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error Deleting TSig Key",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeleteTSIGKey(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDNS_UpdateTSIGKey(t *testing.T) {
	tests := map[string]struct {
		params           UpdateTSIGKeyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *TSIGKeyResponse
		withError        error
	}{
		"200 OK": {
			params: UpdateTSIGKeyRequest{
				TsigKey: &TSIGKey{
					Name:      "example.com.akamai.com.",
					Algorithm: "hmac-sha512",
					Secret:    "fakeR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				},
				Zone: "example.com",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/key",
		},
		"500 internal server error": {
			params: UpdateTSIGKeyRequest{
				TsigKey: &TSIGKey{
					Name:      "example.com.akamai.com.",
					Algorithm: "hmac-sha512",
					Secret:    "fakeR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				},
				Zone: "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/key",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
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
			err := client.UpdateTSIGKey(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
