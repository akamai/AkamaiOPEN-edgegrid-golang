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

func TestDNS_GetAuthorities(t *testing.T) {
	tests := map[string]struct {
		params           GetAuthoritiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAuthoritiesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         GetAuthoritiesRequest{ContractIDs: "9-9XXXXX"},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"contracts": [
        {
            "contractId": "9-9XXXXX",
            "authorities": [
                "a1-118.akam.net.",
                "a2-64.akam.net.",
                "a6-66.akam.net.",
                "a18-67.akam.net.",
                "a7-64.akam.net.",
                "a11-64.akam.net."
            ]
        }
    ]
}`,
			expectedPath: "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
			expectedResponse: &GetAuthoritiesResponse{
				Contracts: []Contract{
					{
						ContractID: "9-9XXXXX",
						Authorities: []string{
							"a1-118.akam.net.",
							"a2-64.akam.net.",
							"a6-66.akam.net.",
							"a18-67.akam.net.",
							"a7-64.akam.net.",
							"a11-64.akam.net.",
						},
					},
				},
			},
		},
		"Missing arguments": {
			responseStatus: http.StatusOK,
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get authorities: struct validation: ContractIDs: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params:         GetAuthoritiesRequest{ContractIDs: "9-9XXXXX"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching authorities",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
			result, err := client.GetAuthorities(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDNS_GetNameServerRecordList(t *testing.T) {
	tests := map[string]struct {
		params           GetNameServerRecordListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        func(*testing.T, error)
	}{
		"test with valid arguments": {
			params:         GetNameServerRecordListRequest{ContractIDs: "9-9XXXXX"},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"contracts": [
        {
            "contractId": "9-9XXXXX",
            "authorities": [
                "a1-118.akam.net.",
                "a2-64.akam.net.",
                "a6-66.akam.net.",
                "a18-67.akam.net.",
                "a7-64.akam.net.",
                "a11-64.akam.net."
            ]
        }
    ]
}`,
			expectedResponse: []string{"a1-118.akam.net.", "a2-64.akam.net.", "a6-66.akam.net.", "a18-67.akam.net.", "a7-64.akam.net.", "a11-64.akam.net."},
			expectedPath:     "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
		},
		"test with missing arguments": {
			expectedPath: "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get name server record list: struct validation: ContractIDs: cannot be blank", err.Error())
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
			result, err := client.GetNameServerRecordList(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
