package dns

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDNS_ListZones(t *testing.T) {

	tests := map[string]struct {
		params           ListZonesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneListResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: ListZonesRequest{
				ContractIDs: "1-1ACYUM",
				Search:      "org",
				SortBy:      "-contractId,zone",
				Types:       "primary,alias",
				Page:        1,
				PageSize:    25,
			},
			headers: http.Header{
				"Accept": []string{"application/json"},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"metadata": {
					"page": 1,
					"pageSize": 3,
					"showAll": false,
					"totalElements": 17,
					"contractIds": [
						"1-2ABCDE"
					]
				},
				"zones": [
					{
						"contractId": "1-2ABCDE",
						"zone": "example.com",
						"type": "primary",
						"aliasCount": 1,
						"signAndServe": false,
						"versionId": "ae02357c-693d-4ac4-b33d-8352d9b7c786",
						"lastModifiedDate": "2017-01-03T12:00:00Z",
						"lastModifiedBy": "user28",
						"lastActivationDate": "2017-01-03T12:00:00Z",
						"activationState": "ACTIVE"
					}
				]
			}`,
			expectedPath: "/config-dns/v2/zones?contractIds=1-1ACYUM&search=org&sortBy=-contractId%2Czone&types=primary%2Calias&page=1&pageSize=25&showAll=false",
			expectedResponse: &ZoneListResponse{
				Metadata: &ListMetadata{
					Page:          1,
					PageSize:      3,
					ShowAll:       false,
					TotalElements: 17,
					ContractIDs:   []string{"1-2ABCDE"},
				},
				Zones: []ZoneResponse{
					{
						ContractID:         "1-2ABCDE",
						Zone:               "example.com",
						Type:               "primary",
						AliasCount:         1,
						SignAndServe:       false,
						VersionID:          "ae02357c-693d-4ac4-b33d-8352d9b7c786",
						LastModifiedDate:   "2017-01-03T12:00:00Z",
						LastModifiedBy:     "user28",
						LastActivationDate: "2017-01-03T12:00:00Z",
						ActivationState:    "ACTIVE",
					},
				},
			},
		},
		"500 internal server error": {
			params: ListZonesRequest{
				ContractIDs: "1-1ACYUM",
				Search:      "org",
				SortBy:      "-contractId,zone",
				Types:       "primary,alias",
				Page:        1,
				PageSize:    25,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones?contractIds=1-1ACYUM&search=org&sortBy=-contractId%2Czone&types=primary%2Calias&page=1&pageSize=25&showAll=false",
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
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListZones(
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

func TestDNS_GetZonesDNSSecStatus(t *testing.T) {

	tests := map[string]struct {
		zones               []string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *GetZonesDNSSecStatusResponse
		withError           error
		headers             http.Header
	}{
		"200 OK current records only": {
			zones: []string{"foo.test.net"},
			headers: http.Header{
				"Accept": []string{"application/json"},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"dnsSecStatuses": [
					{
						"zone": "foo.test.net",
						"alerts": [
							"PARENT_DS_MISSING"
						],
						"currentRecords": {
							"dsRecord": "foo.test.net. 86400 IN DS 42061 7 2 ( DUMMY_HASH_1 ) ",
							"dnskeyRecord": "foo.test.net. 7200 IN DNSKEY 257 3 7 (DUMMY_HASH_2 ) ",
							"lastModifiedDate": "2024-05-28T06:58:26Z",
							"expectedTtl": 0
						}
					}
				]
			}`,
			expectedPath:        "/config-dns/v2/zones/dns-sec-status",
			expectedRequestBody: `{"zones":["foo.test.net"]}`,
			expectedResponse: &GetZonesDNSSecStatusResponse{
				DNSSecStatuses: []SecStatus{{
					Zone:   "foo.test.net",
					Alerts: []string{"PARENT_DS_MISSING"},
					CurrentRecords: SecRecords{
						DNSKeyRecord:     "foo.test.net. 7200 IN DNSKEY 257 3 7 (DUMMY_HASH_2 ) ",
						DSRecord:         "foo.test.net. 86400 IN DS 42061 7 2 ( DUMMY_HASH_1 ) ",
						ExpectedTTL:      0,
						LastModifiedDate: test.NewTimeFromString(t, "2024-05-28T06:58:26Z"),
					},
				}},
			},
		},
		"200 OK new records returned": {
			zones: []string{"foo.test.net"},
			headers: http.Header{
				"Accept": []string{"application/json"},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"dnsSecStatuses": [
					{
						"alerts": [
							"PARENT_DS_MISSING"
						],
						"currentRecords": {
							"dnskeyRecord": "foo.test.net. 7200 IN DNSKEY 257 3 13 (DUMMY_HASH_1 ) ",
							"dsRecord": "foo.test.net. 86400 IN DS 3622 13 2 ( DUMMY_HASH_2 ) ",
							"expectedTtl": 3600,
							"lastModifiedDate": "2022-06-19T10:14:35Z"
						},
						"newRecords": {
							"dnskeyRecord": "foo.test.net. 7200 IN DNSKEY 257 3 13 (DUMMY_HASH_3 ) ",
							"dsRecord": "foo.test.net. 86400 IN DS 39035 13 2 ( DUMMY_HASH_4 ) ",
							"expectedTtl": 3600,
							"lastModifiedDate": "2023-06-19T10:14:35Z"
						},
						"zone": "foo.test.net"
					}
				]
			}`,
			expectedPath:        "/config-dns/v2/zones/dns-sec-status",
			expectedRequestBody: `{"zones":["foo.test.net"]}`,
			expectedResponse: &GetZonesDNSSecStatusResponse{
				DNSSecStatuses: []SecStatus{{
					Zone:   "foo.test.net",
					Alerts: []string{"PARENT_DS_MISSING"},
					CurrentRecords: SecRecords{
						DNSKeyRecord:     "foo.test.net. 7200 IN DNSKEY 257 3 13 (DUMMY_HASH_1 ) ",
						DSRecord:         "foo.test.net. 86400 IN DS 3622 13 2 ( DUMMY_HASH_2 ) ",
						ExpectedTTL:      3600,
						LastModifiedDate: test.NewTimeFromString(t, "2022-06-19T10:14:35Z"),
					},
					NewRecords: &SecRecords{
						DNSKeyRecord:     "foo.test.net. 7200 IN DNSKEY 257 3 13 (DUMMY_HASH_3 ) ",
						DSRecord:         "foo.test.net. 86400 IN DS 39035 13 2 ( DUMMY_HASH_4 ) ",
						ExpectedTTL:      3600,
						LastModifiedDate: test.NewTimeFromString(t, "2023-06-19T10:14:35Z"),
					},
				}},
			},
		},
		"500 internal server error": {
			zones:          []string{"foo.test.net"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
			   "type": "https://problems.luna.akamaiapis.net/authoritative-dns/serverError",
			   "title": "Server error",
			   "instance": "29aa48de-ec7d-4214-ad6c-649163889be7",
			   "status": 500,
			   "detail": "An internal error occurred.",
			   "problemId": "29aa48de-ec7d-4214-ad6c-649163889be7"
			}`,
			expectedPath:        "/config-dns/v2/zones/dns-sec-status",
			expectedRequestBody: `{"zones":["foo.test.net"]}`,
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/authoritative-dns/serverError",
				Title:      "Server error",
				Detail:     "An internal error occurred.",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error: empty zone list": {
			zones:     []string{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetZonesDNSSecStatus(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)),
				GetZonesDNSSecStatusRequest{
					Zones: test.zones})
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetZone(t *testing.T) {
	tests := map[string]struct {
		params           GetZoneRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetZoneResponse
		withError        error
	}{
		"200 OK": {
			params: GetZoneRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"contractId": "1-2ABCDE",
				"zone": "example.com",
				"type": "primary",
				"aliasCount": 1,
				"signAndServe": true,
				"signAndServeAlgorithm": "RSA_SHA256",
				"versionId": "ae02357c-693d-4ac4-b33d-8352d9b7c786",
				"lastModifiedDate": "2017-01-03T12:00:00Z",
				"lastModifiedBy": "user28",
				"lastActivationDate": "2017-01-03T12:00:00Z",
				"activationState": "ACTIVE"
			}`,
			expectedPath: "/config-dns/v2/zones/example.com",
			expectedResponse: &GetZoneResponse{
				ContractID:            "1-2ABCDE",
				Zone:                  "example.com",
				Type:                  "primary",
				AliasCount:            1,
				SignAndServe:          true,
				SignAndServeAlgorithm: "RSA_SHA256",
				VersionID:             "ae02357c-693d-4ac4-b33d-8352d9b7c786",
				LastModifiedDate:      "2017-01-03T12:00:00Z",
				LastModifiedBy:        "user28",
				LastActivationDate:    "2017-01-03T12:00:00Z",
				ActivationState:       "ACTIVE",
			},
		},
		"500 internal server error": {
			params: GetZoneRequest{
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
			expectedPath: "/config-dns/v2/zones/example.com",
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
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetZone(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetZoneMasterFile(t *testing.T) {
	tests := map[string]struct {
		params           GetMasterZoneFileRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse string
		withError        error
	}{
		"200 OK": {
			params: GetMasterZoneFileRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `"example.com.        10000    IN SOA ns1.akamaidns.com. webmaster.example.com. 1 28800 14400 2419200 86400
example.com.        10000    IN NS  ns1.akamaidns.com.
example.com.        10000    IN NS  ns2.akamaidns.com.
example.com.            300 IN  A   10.0.0.1
example.com.            300 IN  A   10.0.0.2
www.example.com.        300 IN  A   10.0.0.1
www.example.com.        300 IN  A   10.0.0.2"`,
			expectedPath: "/config-dns/v2/zones/example.com/zone-file",
			expectedResponse: `"example.com.        10000    IN SOA ns1.akamaidns.com. webmaster.example.com. 1 28800 14400 2419200 86400
example.com.        10000    IN NS  ns1.akamaidns.com.
example.com.        10000    IN NS  ns2.akamaidns.com.
example.com.            300 IN  A   10.0.0.1
example.com.            300 IN  A   10.0.0.2
www.example.com.        300 IN  A   10.0.0.1
www.example.com.        300 IN  A   10.0.0.2"`,
		},
		"500 internal server error": {
			params: GetMasterZoneFileRequest{
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
			expectedPath: "/config-dns/v2/zones/example.com/zone-file",
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
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetMasterZoneFile(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_UpdateZoneMasterFile(t *testing.T) {
	tests := map[string]struct {
		params         PostMasterZoneFileRequest
		responseStatus int
		expectedPath   string
		responseBody   string
		withError      error
	}{
		"204 Updated": {
			params: PostMasterZoneFileRequest{
				Zone: "example.com",
				FileData: `"example.com.        10000    IN SOA ns1.akamaidns.com. webmaster.example.com. 1 28800 14400 2419200 86400
example.com.        10000    IN NS  ns1.akamaidns.com.
example.com.        10000    IN NS  ns2.akamaidns.com.
example.com.            300 IN  A   10.0.0.1
example.com.            300 IN  A   10.0.0.2
www.example.com.        300 IN  A   10.0.0.1
www.example.com.        300 IN  A   10.0.0.2"`,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/config-dns/v2/zones/example.com/zone-file",
		},
		"500 internal server error": {
			params: PostMasterZoneFileRequest{
				Zone: "example.com",
				FileData: `"example.com.        10000    IN SOA ns1.akamaidns.com. webmaster.example.com. 1 28800 14400 2419200 86400
example.com.        10000    IN NS  ns1.akamaidns.com.
example.com.        10000    IN NS  ns2.akamaidns.com.
example.com.            300 IN  A   10.0.0.1
example.com.            300 IN  A   10.0.0.2
www.example.com.        300 IN  A   10.0.0.1
www.example.com.        300 IN  A   10.0.0.2"`,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/zone-file",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.PostMasterZoneFile(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDNS_GetChangeList(t *testing.T) {
	tests := map[string]struct {
		params           GetChangeListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetChangeListResponse
		withError        error
	}{
		"200 OK": {
			params: GetChangeListRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"zone": "example.com",
				"changeTag": "476754f4-d605-479f-853b-db854d7254fa",
				"zoneVersionId": "1d9c887c-49bb-4382-87a6-d1bf690aa58f",
				"lastModifiedDate": "2017-02-01T12:00:12.524Z",
				"stale": false
			}`,
			expectedPath: "/config-dns/v2/zones/example.com",
			expectedResponse: &GetChangeListResponse{
				Zone:             "example.com",
				ChangeTag:        "476754f4-d605-479f-853b-db854d7254fa",
				ZoneVersionID:    "1d9c887c-49bb-4382-87a6-d1bf690aa58f",
				LastModifiedDate: "2017-02-01T12:00:12.524Z",
				Stale:            false,
			},
		},
		"500 internal server error": {
			params: GetChangeListRequest{
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
			expectedPath: "/config-dns/v2/zones/example.com",
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
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetChangeList(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetMasterZoneFile(t *testing.T) {
	tests := map[string]struct {
		params           GetMasterZoneFileRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse string
		withError        error
	}{
		"200 OK": {
			params: GetMasterZoneFileRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			example.com.        10000    IN SOA ns1.akamaidns.com. webmaster.example.com. 1 28800 14400 2419200 86400
			example.com.        10000    IN NS  ns1.akamaidns.com.
			example.com.        10000    IN NS  ns2.akamaidns.com.
			example.com.            300 IN  A   10.0.0.1
			example.com.            300 IN  A   10.0.0.2
			www.example.com.        300 IN  A   10.0.0.1
			www.example.com.        300 IN  A   10.0.0.2`,
			expectedPath: "/config-dns/v2/zones/example.com/zone-file",
			expectedResponse: `
			example.com.        10000    IN SOA ns1.akamaidns.com. webmaster.example.com. 1 28800 14400 2419200 86400
			example.com.        10000    IN NS  ns1.akamaidns.com.
			example.com.        10000    IN NS  ns2.akamaidns.com.
			example.com.            300 IN  A   10.0.0.1
			example.com.            300 IN  A   10.0.0.2
			www.example.com.        300 IN  A   10.0.0.1
			www.example.com.        300 IN  A   10.0.0.2`,
		},
		"500 internal server error": {
			params: GetMasterZoneFileRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching master zone file",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/zone-file",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching master zone file",
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
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetMasterZoneFile(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_CreateZone(t *testing.T) {
	tests := map[string]struct {
		params         CreateZoneRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"201 Created": {
			params: CreateZoneRequest{
				CreateZone: &ZoneCreate{
					Zone:       "example.com",
					ContractID: "1-2ABCDE",
					Type:       "primary",
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
			{
				"contractId": "1-2ABCDE",
				"zone": "other.com",
				"type": "primary",
				"aliasCount": 1,
				"signAndServe": false,
				"comment": "Initial add",
				"versionId": "7949b2db-ac43-4773-a3ec-dc93202142fd",
				"lastModifiedDate": "2016-12-11T03:21:00Z",
				"lastModifiedBy": "user31",
				"lastActivationDate": "2017-01-03T12:00:00Z",
				"activationState": "ERROR",
				"masters": [
					"1.2.3.4",
					"1.2.3.5"
				],
				"tsigKey": {
					"name": "other.com.akamai.com.",
					"algorithm": "hmac-sha512",
					"secret": "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="
				}
			}`,
			expectedPath: "/config-dns/v2/zones?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			params: CreateZoneRequest{
				CreateZone: &ZoneCreate{
					Zone:       "example.com",
					ContractID: "1-2ABCDE",
					Type:       "primary",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones?contractId=1-2ABCDE",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.CreateZone(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDNS_SaveChangelist(t *testing.T) {
	tests := map[string]struct {
		params         SaveChangeListRequest
		zone           ZoneCreate
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"201 Created": {
			params: SaveChangeListRequest{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "primary",
			},
			responseStatus: http.StatusCreated,
			expectedPath:   "/config-dns/v2/changelists?zone=example.com",
		},
		"500 internal server error": {
			params: SaveChangeListRequest{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "primary",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/changelists?zone=example.com",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating zone",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
			err := client.SaveChangeList(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDNS_SubmitChangelist(t *testing.T) {
	tests := map[string]struct {
		params         SubmitChangeListRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 No Content": {
			params: SubmitChangeListRequest{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "primary",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/changelists?zone=example.com",
		},
		"500 internal server error": {
			params: SubmitChangeListRequest{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "secondary",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/changelists?zone=example.com",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.SubmitChangeList(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDNS_UpdateZone(t *testing.T) {
	tests := map[string]struct {
		params         UpdateZoneRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"200 OK": {
			params: UpdateZoneRequest{
				CreateZone: &ZoneCreate{
					Zone:       "example.com",
					ContractID: "1-2ABCDE",
					Type:       "primary",
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"contractId": "1-2ABCDE",
				"zone": "other.com",
				"type": "primary",
				"aliasCount": 1,
				"signAndServe": false,
				"comment": "Initial add",
				"versionId": "7949b2db-ac43-4773-a3ec-dc93202142fd",
				"lastModifiedDate": "2016-12-11T03:21:00Z",
				"lastModifiedBy": "user31",
				"lastActivationDate": "2017-01-03T12:00:00Z",
				"activationState": "ERROR",
				"masters": [
					"1.2.3.4",
					"1.2.3.5"
				],
				"tsigKey": {
					"name": "other.com.akamai.com.",
					"algorithm": "hmac-sha512",
					"secret": "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="
				}
			}`,
			expectedPath: "/config-dns/v2/zones?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			params: UpdateZoneRequest{
				CreateZone: &ZoneCreate{
					Zone:       "example.com",
					ContractID: "1-2ABCDE",
					Type:       "secondary",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones?contractId=1-2ABCDE",
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
			err := client.UpdateZone(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDNS_GetZoneNames(t *testing.T) {
	tests := map[string]struct {
		params           GetZoneNamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetZoneNamesResponse
		withError        error
	}{
		"200 OK": {
			params: GetZoneNamesRequest{
				Zone: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"names": [
					"example.com",
					"www.example.com",
					"ftp.example.com",
					"space.example.com",
					"bar.example.com"
				]
			}`,
			expectedPath: "/config-dns/v2/zones/example.com/names",
			expectedResponse: &GetZoneNamesResponse{
				Names: []string{"example.com", "www.example.com", "ftp.example.com", "space.example.com", "bar.example.com"},
			},
		},
		"500 internal server error": {
			params: GetZoneNamesRequest{
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
			expectedPath: "/config-dns/v2/zones/example.com/names",
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
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetZoneNames(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetZoneNameTypes(t *testing.T) {
	tests := map[string]struct {
		params           GetZoneNameTypesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetZoneNameTypesResponse
		withError        error
	}{
		"200 OK": {
			params: GetZoneNameTypesRequest{
				Zone:     "example.com",
				ZoneName: "names",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"types": [
					"A",
					"AAAA",
					"MX"
				]
			}`,
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types",
			expectedResponse: &GetZoneNameTypesResponse{
				Types: []string{"A", "AAAA", "MX"},
			},
		},
		"500 internal server error": {
			params: GetZoneNameTypesRequest{
				Zone:     "example.com",
				ZoneName: "names",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types",
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
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetZoneNameTypes(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func Test_ValidateZoneErrors(t *testing.T) {
	tests := map[string]ZoneCreate{
		"empty zone": {},
		"bad type": {
			Zone: "example.com",
			Type: "BAD",
		},
		"secondary tsig": {
			Zone: "example.com",
			Type: "PRIMARY",
			TSIGKey: &TSIGKey{
				Name: "example.com",
			},
		},
		"alias empty target": {
			Zone:   "example.com",
			Type:   "ALIAS",
			Target: "",
		},
		"alias masters": {
			Zone:    "example.com",
			Type:    "ALIAS",
			Target:  "10.0.0.1",
			Masters: []string{"master"},
		},
		"alias sign": {
			Zone:         "example.com",
			Type:         "ALIAS",
			Target:       "10.0.0.1",
			SignAndServe: true,
		},
		"alias sign algo": {
			Zone:                  "example.com",
			Type:                  "ALIAS",
			Target:                "10.0.0.1",
			SignAndServe:          false,
			SignAndServeAlgorithm: "foo",
		},
		"primary bad target": {
			Zone:   "example.com",
			Type:   "PRIMARY",
			Target: "10.0.0.1",
		},
		"primary bad masters": {
			Zone:    "example.com",
			Type:    "PRIMARY",
			Masters: []string{"foo"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateZone(&test)
			assert.NotNil(t, err)
		})
	}
}
