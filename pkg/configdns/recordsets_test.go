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

func TestDns_GetRecordsets(t *testing.T) {
	tests := map[string]struct {
		zone             string
		args             []RecordsetQueryArgs
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RecordSetResponse
		withError        error
	}{
		"200 OK": {
			zone:           "example.com",
			args:           []RecordsetQueryArgs{},
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
			expectedPath: "/config-dns/v2/zones/example.com/recordsets",
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
			args:           []RecordsetQueryArgs{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/recordsets",
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
			result, err := client.GetRecordsets(context.Background(), test.zone, test.args...)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_CreateRecordsets(t *testing.T) {
	tests := map[string]struct {
		zone             string
		sets             *Recordsets
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RecordSetResponse
		withError        error
	}{
		"200 OK": {
			zone: "example.com",
			sets: &Recordsets{
				[]Recordset{
					{
						Name: "www.example.com",
						Type: "A",
						TTL:  300,
						Rdata: []string{
							"10.0.0.2",
							"10.0.0.3",
						},
					},
				},
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/recordsets",
		},
		"500 internal server error": {
			zone: "example.com",
			sets: &Recordsets{
				[]Recordset{
					{
						Name: "www.example.com",
						Type: "A",
						TTL:  300,
						Rdata: []string{
							"10.0.0.2",
							"10.0.0.3",
						},
					},
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
			expectedPath: "/config-dns/v2/zones/example.com/recordsets",
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
			err := client.CreateRecordsets(context.Background(), test.sets, test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDns_UpdateRecordsets(t *testing.T) {
	tests := map[string]struct {
		zone             string
		sets             *Recordsets
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RecordSetResponse
		withError        error
	}{
		"200 OK": {
			zone: "example.com",
			sets: &Recordsets{
				[]Recordset{
					{
						Name: "www.example.com",
						Type: "A",
						TTL:  300,
						Rdata: []string{
							"10.0.0.2",
							"10.0.0.3",
						},
					},
				},
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/recordsets",
		},
		"500 internal server error": {
			zone: "example.com",
			sets: &Recordsets{
				[]Recordset{
					{
						Name: "www.example.com",
						Type: "A",
						TTL:  300,
						Rdata: []string{
							"10.0.0.2",
							"10.0.0.3",
						},
					},
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
			expectedPath: "/config-dns/v2/zones/example.com/recordsets",
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
			err := client.UpdateRecordsets(context.Background(), test.sets, test.zone)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
