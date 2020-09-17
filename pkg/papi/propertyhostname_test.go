package papi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestPapi_GetHostnames(t *testing.T) {
	tests := map[string]struct {
		params           GetHostnamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupId:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_1-1TJZFB",
    "contractId": "ctr_1-1TJZH5",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
    "propertyVersion": 3,
    "etag": "6aed418629b4e5c0",
    "hostnames": {
        "items": [
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895822",
                "cnameFrom": "example.com",
                "cnameTo": "example.com.edgesuite.net"
            },
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895833",
                "cnameFrom": "m.example.com",
                "cnameTo": "m.example.com.edgesuite.net"
            }
        ]
    }
}

`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&validateHostnames=false",
			expectedResponse: &GetHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				Hostnames: HostnameItems{
					Items: []HostnameItem{
						{
							CnameType:      "EDGE_HOSTNAME",
							EdgeHostnameID: "ehn_895822",
							CnameFrom:      "example.com",
							CnameTo:        "example.com.edgesuite.net",
						},
						{
							CnameType:      "EDGE_HOSTNAME",
							EdgeHostnameID: "ehn_895833",
							CnameFrom:      "m.example.com",
							CnameTo:        "m.example.com.edgesuite.net",
						},
					},
				},
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
			result, err := client.GetHostnames(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_CreateHostnames(t *testing.T) {
	tests := map[string]struct {
		params           CreateHostnamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: CreateHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupID:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_1-1TJZFB",
    "contractId": "ctr_1-1TJZH5",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
    "propertyVersion": 3,
    "etag": "6aed418629b4e5c0",
    "hostnames": {
        "items": [
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895822",
                "cnameFrom": "example.com",
                "cnameTo": "example.com.edgesuite.net"
            },
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895833",
                "cnameFrom": "m.example.com",
                "cnameTo": "m.example.com.edgesuite.net"
            }
        ]
    }
}

`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&validateHostnames=false",
			expectedResponse: &CreateHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				Hostnames: HostnameItems{
					Items: []HostnameItem{
						{
							CnameType:      "EDGE_HOSTNAME",
							EdgeHostnameID: "ehn_895822",
							CnameFrom:      "example.com",
							CnameTo:        "example.com.edgesuite.net",
						},
						{
							CnameType:      "EDGE_HOSTNAME",
							EdgeHostnameID: "ehn_895833",
							CnameFrom:      "m.example.com",
							CnameTo:        "m.example.com.edgesuite.net",
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateHostnames(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
