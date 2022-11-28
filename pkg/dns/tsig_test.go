package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewTsigQueryString(t *testing.T) {
	client := Client(session.Must(session.New()))

	str := client.NewTsigQueryString(context.Background())

	assert.NotNil(t, str)
}

func TestDns_ListTsigKeys(t *testing.T) {
	tests := map[string]struct {
		query            TSIGQueryString
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *TSIGReportResponse
		withError        error
	}{
		"200 OK": {
			query: TSIGQueryString{
				ContractIds: []string{"1-1ABCD"},
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
			expectedResponse: &TSIGReportResponse{
				Metadata: &TSIGReportMeta{
					TotalElements: 1,
				},
				Keys: []*TSIGKeyResponse{
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
			query: TSIGQueryString{
				ContractIds: []string{"1-1ABCD"},
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
			result, err := client.ListTsigKeys(context.Background(), &test.query)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetTsigKeyZones(t *testing.T) {
	tests := map[string]struct {
		key              TSIGKey
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneNameListResponse
		withError        error
	}{
		"200 OK": {
			key: TSIGKey{
				Name:      "example.com.akamai.com.",
				Algorithm: "hmac-sha512",
				Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
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
			expectedResponse: &ZoneNameListResponse{
				Zones: []string{"river.com", "stream.com"},
			},
		},
		"500 internal server error": {
			key: TSIGKey{
				Name:      "example.com.akamai.com.",
				Algorithm: "hmac-sha512",
				Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
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
			result, err := client.GetTsigKeyZones(context.Background(), &test.key)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetTsigKeyAliases(t *testing.T) {
	tests := map[string]struct {
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneNameListResponse
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
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
			expectedResponse: &ZoneNameListResponse{
				Aliases: []string{"exmaple.com", "river.com", "brook.com", "ocean.com"},
			},
		},
		"500 internal server error": {
			zone:           "example.com",
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
			result, err := client.GetTsigKeyAliases(context.Background(), test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_TsigKeyBulkUpdate(t *testing.T) {
	tests := map[string]struct {
		bulk             TSIGKeyBulkPost
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneNameListResponse
		withError        error
	}{
		"200 OK": {
			bulk: TSIGKeyBulkPost{
				Key: &TSIGKey{
					Name:      "brook.com.akamai.com.",
					Algorithm: "hmac-sha512",
					Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				},
				Zones: []string{"brook.com", "river.com"},
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/keys/bulk-update",
			expectedResponse: &ZoneNameListResponse{
				Aliases: []string{"exmaple.com", "river.com", "brook.com", "ocean.com"},
			},
		},
		"500 internal server error": {
			bulk: TSIGKeyBulkPost{
				Key: &TSIGKey{
					Name:      "brook.com.akamai.com.",
					Algorithm: "hmac-sha512",
					Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				},
				Zones: []string{"brook.com", "river.com"},
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
			err := client.TsigKeyBulkUpdate(context.Background(), &test.bulk)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_GetTsigKey(t *testing.T) {
	tests := map[string]struct {
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *TSIGKeyResponse
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"name": "example.com.akamai.com.",
				"algorithm": "hmac-sha512",
				"secret": "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				"zonesCount": 7
			}`,
			expectedPath: "/config-dns/v2/zones/example.com/key",
			expectedResponse: &TSIGKeyResponse{
				TSIGKey: TSIGKey{
					Name:      "example.com.akamai.com.",
					Algorithm: "hmac-sha512",
					Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
				},
				ZoneCount: 7,
			},
		},
		"500 internal server error": {
			zone:           "example.com",
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
			result, err := client.GetTsigKey(context.Background(), test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_DeleteTsigKey(t *testing.T) {
	tests := map[string]struct {
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *TSIGKeyResponse
		withError        error
	}{
		"204 No Content": {
			zone:           "example.com",
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/key",
		},
		"500 internal server error": {
			zone:           "example.com",
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
			err := client.DeleteTsigKey(context.Background(), test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_UpdateTsigKey(t *testing.T) {
	tests := map[string]struct {
		key              TSIGKey
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *TSIGKeyResponse
		withError        error
	}{
		"200 OK": {
			key: TSIGKey{
				Name:      "example.com.akamai.com.",
				Algorithm: "hmac-sha512",
				Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
			},
			zone:           "example.com",
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/key",
		},
		"500 internal server error": {
			key: TSIGKey{
				Name:      "example.com.akamai.com.",
				Algorithm: "hmac-sha512",
				Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
			},
			zone:           "example.com",
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
			err := client.UpdateTsigKey(context.Background(), &test.key, test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
