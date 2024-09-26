package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDNS_GetRecord(t *testing.T) {
	tests := map[string]struct {
		params           GetRecordRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRecordResponse
		withError        error
	}{
		"200 OK": {
			params: GetRecordRequest{
				Zone:       "example.com",
				Name:       "www.example.com",
				RecordType: "A",
			},
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
			expectedResponse: &GetRecordResponse{
				Name:       "www.example.com",
				RecordType: "A",
				TTL:        300,
				Active:     false,
				Target:     []string{"10.0.0.2", "10.0.0.3"},
			},
		},
		"500 internal server error": {
			params: GetRecordRequest{
				Zone:       "example.com",
				Name:       "www.example.com",
				RecordType: "A",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/names/www.example.com/types/A",
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
			result, err := client.GetRecord(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetRecordList(t *testing.T) {
	tests := map[string]struct {
		params           GetRecordListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRecordListResponse
		withError        error
	}{
		"200 OK": {
			params: GetRecordListRequest{
				Zone:       "example.com",
				RecordType: "A",
			},
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
			expectedPath: "/config-dns/v2/zones/example.com/recordsets?showAll=true&types=A",
			expectedResponse: &GetRecordListResponse{
				Metadata: Metadata{
					LastPage:      0,
					Page:          1,
					PageSize:      25,
					ShowAll:       false,
					TotalElements: 2,
				},
				RecordSets: []RecordSet{
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
			params: GetRecordListRequest{
				Zone:       "example.com",
				RecordType: "A",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/recordsets?showAll=true&types=A",
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
			result, err := client.GetRecordList(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetRdata(t *testing.T) {
	tests := map[string]struct {
		params           GetRdataRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        error
	}{
		"ipv6 test": {
			params: GetRdataRequest{
				Zone:       "example.com",
				RecordType: "AAAA",
				Name:       "www.example.com",
			},
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
			expectedPath:     "/config-dns/v2/zones/example.com/recordsets?showAll=true&types=AAAA",
			expectedResponse: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		},
		"loc test": {
			params: GetRdataRequest{
				Zone:       "example.com",
				RecordType: "LOC",
				Name:       "www.example.com",
			},
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
			expectedPath:     "/config-dns/v2/zones/example.com/recordsets?showAll=true&types=LOC",
			expectedResponse: []string{"52 22 23.000 N 4 53 32.000 E -2.00m 0.00m 10000.00m 10.00m"},
		},
		"500 internal server error": {
			params: GetRdataRequest{
				Zone:       "example.com",
				RecordType: "A",
				Name:       "www.example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/recordsets?showAll=true&types=A",
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
			result, err := client.GetRdata(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_TestRdata(t *testing.T) {
	client := Client(session.Must(session.New()))

	out := client.ProcessRdata(context.Background(), []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"}, "AAAA")

	assert.Equal(t, []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"}, out)

	out = client.ProcessRdata(context.Background(), []string{"52 22 23.000 N 4 53 32.000 E -2.00m 0.00m 10000.00m 10.00m"}, "LOC")

	assert.Equal(t, []string{"52 22 23.000 N 4 53 32.000 E -2.00m 0.00m 10000.00m 10.00m"}, out)
}

func TestDNS_ParseRData(t *testing.T) {
	client := Client(session.Must(session.New()))

	tests := map[string]struct {
		rType  string
		rdata  []string
		expect map[string]interface{}
	}{
		"AFSDB": {
			rType: "AFSDB",
			rdata: []string{"1 bar.com"},
			expect: map[string]interface{}{
				"subtype": 1,
				"target":  []string{"bar.com"},
			},
		},
		"SVCB": {
			rType: "SVCB",
			rdata: []string{"0 svc4.example.com."},
			expect: map[string]interface{}{
				"target":       []string{},
				"svc_priority": 0,
				"target_name":  "svc4.example.com.",
			},
		},
		"HTTPS": {
			rType: "HTTPS",
			rdata: []string{"3 https.example.com. alpn=bar port=8080"},
			expect: map[string]interface{}{
				"target":       []string{},
				"svc_priority": 3,
				"target_name":  "https.example.com.",
				"svc_params":   "alpn=bar port=8080",
			},
		},
		"SRV with default values": {
			rType: "SRV",
			rdata: []string{"10 60 5060 big.example.com.", "10 60 5060 small.example.com."},
			expect: map[string]interface{}{
				"port":     5060,
				"priority": 10,
				"weight":   60,
				"target":   []string{"big.example.com.", "small.example.com."},
			},
		},
		"SRV without default values": {
			rType: "SRV",
			rdata: []string{"10 60 5060 big.example.com.", "20 50 5060 small.example.com."},
			expect: map[string]interface{}{
				"target": []string{"10 60 5060 big.example.com.", "20 50 5060 small.example.com."},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out := client.ParseRData(context.Background(), test.rType, test.rdata)

			assert.Equal(t, test.expect, out)
		})
	}
}
