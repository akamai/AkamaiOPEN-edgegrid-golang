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

func TestDNS_GetRecordSets(t *testing.T) {
	tests := map[string]struct {
		params           GetRecordSetsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRecordSetsResponse
		withError        error
	}{
		"200 OK": {
			params: GetRecordSetsRequest{
				Zone:      "example.com",
				QueryArgs: &RecordSetQueryArgs{},
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
			expectedPath: "/config-dns/v2/zones/example.com/recordsets?showAll=false",
			expectedResponse: &GetRecordSetsResponse{
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
			params: GetRecordSetsRequest{
				Zone:      "example.com",
				QueryArgs: &RecordSetQueryArgs{},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/example.com/recordsets?showAll=false",
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
			result, err := client.GetRecordSets(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_CreateRecordSets(t *testing.T) {
	tests := map[string]struct {
		params         CreateRecordSetsRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"200 OK": {
			params: CreateRecordSetsRequest{
				Zone: "example.com",
				RecordSets: &RecordSets{
					[]RecordSet{
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
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/recordsets",
		},
		"500 internal server error": {
			params: CreateRecordSetsRequest{
				Zone: "example.com",
				RecordSets: &RecordSets{
					[]RecordSet{
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
			err := client.CreateRecordSets(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDNS_UpdateRecordSets(t *testing.T) {
	tests := map[string]struct {
		params         UpdateRecordSetsRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"200 OK": {
			params: UpdateRecordSetsRequest{
				Zone: "example.com",
				RecordSets: &RecordSets{
					[]RecordSet{
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
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/config-dns/v2/zones/example.com/recordsets",
		},
		"500 internal server error": {
			params: UpdateRecordSetsRequest{
				Zone: "example.com",
				RecordSets: &RecordSets{
					[]RecordSet{
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
			err := client.UpdateRecordSets(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
