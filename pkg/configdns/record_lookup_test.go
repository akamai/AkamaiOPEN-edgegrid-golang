package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDns_GetRecord(t *testing.T) {
	tests := map[string]struct {
		zone             string
		name             string
		recordType       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RecordBody
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
			name:           "www.example.com",
			recordType:     "A",
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"name": "www.example.com",
				"type": "A",
				"ttl": 300,
				"rdata": [
					"10.0.0.2",
					"10.0.0.3"
				]
			}`,
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			expectedResponse: &RecordBody{
				Name:       "www.example.com",
				RecordType: "A",
				TTL:        300,
				Active:     false,
				Target:     []string{"10.0.0.2", "10.0.0.3"},
			},
		},
		"500 internal server error": {
			zone:           "example.com",
			name:           "www.example.com",
			recordType:     "A",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
			withError: session.APIError{
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
			result, err := client.GetRecord(context.Background(), test.zone, test.name, test.recordType)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetRecordList(t *testing.T) {
	tests := map[string]struct {
		zone             string
		name             string
		recordType       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RecordSetResponse
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
			recordType:     "A",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"metadata": {
        "zone": "example.com",
        "page": 1,
        "pageSize": 25,
        "totalElements": 2,
        "types": [
            "A"
        ]
    },
    "recordsets": [
        {
            "name": "www.example.com",
            "type": "A",
            "ttl": 300,
            "rdata": [
                "10.0.0.2",
                "10.0.0.3"
            ]
        }
    ]
}`,
			expectedPath: "/config-dns/v2/zones/example.com/recordsets?types=A&showAll=true",
			expectedResponse: &RecordSetResponse{
				Metadata: MetadataH{
					LastPage:      0,
					Page:          1,
					PageSize:      25,
					ShowAll:       false,
					TotalElements: 2,
				},
				Recordsets: []Recordset{
					{
						Name:  "www.example.com",
						Type:  "A",
						TTL:   300,
						Rdata: []string{"10.0.0.2", "10.0.0.3"},
					},
				},
			},
		},
		"500 internal server error": {
			zone:           "example.com",
			recordType:     "A",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/recordsets?types=A&showAll=true",
			withError: session.APIError{
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
			result, err := client.GetRecordList(context.Background(), test.zone, test.name, test.recordType)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetRdata(t *testing.T) {
	tests := map[string]struct {
		zone             string
		name             string
		recordType       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        error
	}{
		"ipv6 test": {
			zone:           "example.com",
			name:           "www.example.com",
			recordType:     "AAAA",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"metadata": {
        "zone": "example.com",
        "page": 1,
        "pageSize": 25,
        "totalElements": 1,
        "types": [
            "AAAA"
        ]
    },
    "recordsets": [
        {
            "name": "www.example.com",
            "type": "AAAA",
            "ttl": 300,
            "rdata": [
                "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
            ]
		}
    ]
}`,
			expectedPath:     "/config-dns/v2/zones/example.com/recordsets?types=AAAA&showAll=true",
			expectedResponse: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		},
		"loc test": {
			zone:           "example.com",
			name:           "www.example.com",
			recordType:     "LOC",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"metadata": {
        "zone": "example.com",
        "page": 1,
        "pageSize": 25,
        "totalElements": 1,
        "types": [
            "LOC"
        ]
    },
    "recordsets": [
		{
            "name": "www.example.com",
            "type": "LOC",
            "ttl": 300,
            "rdata": [
                "52 22 23.000 N 4 53 32.000 E -2.00m 0.00m 10000m 10m"
            ]
        }
    ]
}`,
			expectedPath:     "/config-dns/v2/zones/example.com/recordsets?types=LOC&showAll=true",
			expectedResponse: []string{"52 22 23.000 N 4 53 32.000 E -2.00m 0.00m 10000.00m 10.00m"},
		},
		"500 internal server error": {
			zone:           "example.com",
			recordType:     "A",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/recordsets?types=A&showAll=true",
			withError: session.APIError{
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
			result, err := client.GetRdata(context.Background(), test.zone, test.name, test.recordType)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
