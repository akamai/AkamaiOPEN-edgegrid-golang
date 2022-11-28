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

func TestDns_ListZones(t *testing.T) {

	tests := map[string]struct {
		args             []ZoneListQueryArgs
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneListResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			args: []ZoneListQueryArgs{
				{
					ContractIDs: "1-1ACYUM",
					Search:      "org",
					SortBy:      "-contractId,zone",
					Types:       "primary,alias",
					Page:        1,
					PageSize:    25,
				},
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
				Zones: []*ZoneResponse{
					{
						ContractID:         "1-2ABCDE",
						Zone:               "example.com",
						Type:               "primary",
						AliasCount:         1,
						SignAndServe:       false,
						VersionId:          "ae02357c-693d-4ac4-b33d-8352d9b7c786",
						LastModifiedDate:   "2017-01-03T12:00:00Z",
						LastModifiedBy:     "user28",
						LastActivationDate: "2017-01-03T12:00:00Z",
						ActivationState:    "ACTIVE",
					},
				},
			},
		},
		"500 internal server error": {
			args: []ZoneListQueryArgs{
				{
					ContractIDs: "1-1ACYUM",
					Search:      "org",
					SortBy:      "-contractId,zone",
					Types:       "primary,alias",
					Page:        1,
					PageSize:    25,
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
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListZones(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.args...)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_NewZone(t *testing.T) {
	client := Client(session.Must(session.New()))

	inp := ZoneCreate{
		Type:   "A",
		Target: "10.0.0.1",
	}

	out := client.NewZone(context.Background(), inp)

	assert.Equal(t, &inp, out)
}

func TestDns_NewZoneResponse(t *testing.T) {
	client := Client(session.Must(session.New()))

	out := client.NewZoneResponse(context.Background(), "example.com")

	assert.Equal(t, out.Zone, "example.com")
}

func TestDns_NewChangeListResponse(t *testing.T) {
	client := Client(session.Must(session.New()))

	out := client.NewChangeListResponse(context.Background(), "example.com")

	assert.Equal(t, out.Zone, "example.com")
}

func TestDns_NewZoneQueryString(t *testing.T) {
	client := Client(session.Must(session.New()))

	out := client.NewZoneQueryString(context.Background(), "foo", "bar")

	assert.Equal(t, out.Contract, "foo")
	assert.Equal(t, out.Group, "bar")
}

func TestDns_GetZone(t *testing.T) {
	tests := map[string]struct {
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneResponse
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
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
			expectedResponse: &ZoneResponse{
				ContractID:            "1-2ABCDE",
				Zone:                  "example.com",
				Type:                  "primary",
				AliasCount:            1,
				SignAndServe:          true,
				SignAndServeAlgorithm: "RSA_SHA256",
				VersionId:             "ae02357c-693d-4ac4-b33d-8352d9b7c786",
				LastModifiedDate:      "2017-01-03T12:00:00Z",
				LastModifiedBy:        "user28",
				LastActivationDate:    "2017-01-03T12:00:00Z",
				ActivationState:       "ACTIVE",
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
			result, err := client.GetZone(context.Background(), test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetZoneMasterFile(t *testing.T) {
	tests := map[string]struct {
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse string
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
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
			zone:           "example.com",
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
			result, err := client.GetMasterZoneFile(context.Background(), test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_UpdateZoneMasterFile(t *testing.T) {
	tests := map[string]struct {
		zone           string
		masterfile     string
		responseStatus int
		expectedPath   string
		responseBody   string
		withError      error
	}{
		"204 Updated": {
			zone: "example.com",
			masterfile: `"example.com.        10000    IN SOA ns1.akamaidns.com. webmaster.example.com. 1 28800 14400 2419200 86400
example.com.        10000    IN NS  ns1.akamaidns.com.
example.com.        10000    IN NS  ns2.akamaidns.com.
example.com.            300 IN  A   10.0.0.1
example.com.            300 IN  A   10.0.0.2
www.example.com.        300 IN  A   10.0.0.1
www.example.com.        300 IN  A   10.0.0.2"`,
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/config-dns/v2/zones/example.com/zone-file",
		},
		"500 internal server error": {
			zone: "example.com",
			masterfile: `"example.com.        10000    IN SOA ns1.akamaidns.com. webmaster.example.com. 1 28800 14400 2419200 86400
example.com.        10000    IN NS  ns1.akamaidns.com.
example.com.        10000    IN NS  ns2.akamaidns.com.
example.com.            300 IN  A   10.0.0.1
example.com.            300 IN  A   10.0.0.2
www.example.com.        300 IN  A   10.0.0.1
www.example.com.        300 IN  A   10.0.0.2"`,
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
			err := client.PostMasterZoneFile(context.Background(), test.zone, test.masterfile)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_GetChangeList(t *testing.T) {
	tests := map[string]struct {
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ChangeListResponse
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
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
			expectedResponse: &ChangeListResponse{
				Zone:             "example.com",
				ChangeTag:        "476754f4-d605-479f-853b-db854d7254fa",
				ZoneVersionID:    "1d9c887c-49bb-4382-87a6-d1bf690aa58f",
				LastModifiedDate: "2017-02-01T12:00:12.524Z",
				Stale:            false,
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
			result, err := client.GetChangeList(context.Background(), test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetMasterZoneFile(t *testing.T) {
	tests := map[string]struct {
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse string
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
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
			zone:           "example.com",
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
			result, err := client.GetMasterZoneFile(context.Background(), test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_CreateZone(t *testing.T) {
	tests := map[string]struct {
		zone           ZoneCreate
		query          ZoneQueryString
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"201 Created": {
			zone: ZoneCreate{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "primary",
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
			zone: ZoneCreate{
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
			err := client.CreateZone(context.Background(), &test.zone, test.query, true)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_SaveChangelist(t *testing.T) {
	tests := map[string]struct {
		zone           ZoneCreate
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"201 Created": {
			zone: ZoneCreate{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "primary",
			},
			responseStatus: http.StatusCreated,
			expectedPath:   "/config-dns/v2/changelists?zone=example.com",
		},
		"500 internal server error": {
			zone: ZoneCreate{
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
			err := client.SaveChangelist(context.Background(), &test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_SubmitChangelist(t *testing.T) {
	tests := map[string]struct {
		zone           ZoneCreate
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 No Content": {
			zone: ZoneCreate{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "primary",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/changelists?zone=example.com",
		},
		"500 internal server error": {
			zone: ZoneCreate{
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
			err := client.SubmitChangelist(context.Background(), &test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_UpdateZone(t *testing.T) {
	tests := map[string]struct {
		zone           ZoneCreate
		query          ZoneQueryString
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"200 OK": {
			zone: ZoneCreate{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "primary",
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
			zone: ZoneCreate{
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
			err := client.UpdateZone(context.Background(), &test.zone, test.query)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_DeleteZone(t *testing.T) {
	tests := map[string]struct {
		zone           ZoneCreate
		query          ZoneQueryString
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 No Content": {
			zone: ZoneCreate{
				Zone:       "example.com",
				ContractID: "1-2ABCDE",
				Type:       "primary",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			zone: ZoneCreate{
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
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeleteZone(context.Background(), &test.zone, test.query)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_GetZoneNames(t *testing.T) {
	tests := map[string]struct {
		zone             string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneNamesResponse
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
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
			expectedResponse: &ZoneNamesResponse{
				Names: []string{"example.com", "www.example.com", "ftp.example.com", "space.example.com", "bar.example.com"},
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
			result, err := client.GetZoneNames(context.Background(), test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetZoneNameTypes(t *testing.T) {
	tests := map[string]struct {
		zone             string
		zname            string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ZoneNameTypesResponse
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
			zname:          "names",
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
			expectedResponse: &ZoneNameTypesResponse{
				Types: []string{"A", "AAAA", "MX"},
			},
		},
		"500 internal server error": {
			zone:           "example.com",
			zname:          "names",
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
			result, err := client.GetZoneNameTypes(context.Background(), test.zname, test.zone)
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
	client := Client(session.Must(session.New()))

	tests := map[string]ZoneCreate{
		"empty zone": {},
		"bad type": {
			Zone: "example.com",
			Type: "BAD",
		},
		"secondary tsig": {
			Zone: "example.com",
			Type: "PRIMARY",
			TsigKey: &TSIGKey{
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
			err := client.ValidateZone(context.Background(), &test)
			assert.NotNil(t, err)
		})
	}
}
