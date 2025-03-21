package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiListActivePropertyHostnames(t *testing.T) {
	tests := map[string]struct {
		params           ListActivePropertyHostnamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListActivePropertyHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - required params": {
			params: ListActivePropertyHostnamesRequest{
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_123",
    "contractId": "ctr_123",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
	"defaultSort": "hostname:a",
    "currentSort": "hostname:d", 
    "hostnames": {
		"currentItemCount": 2,
        "items": [
            {
                "cnameFrom": "example.com",
        		"cnameType": "EDGE_HOSTNAME",
        		"productionCertType": "DEFAULT",
        		"productionCnameTo": "example.com.edgekey.net",
        		"productionEdgeHostnameID": "ehn_895822"
            },
            {
                "cnameFrom": "m-example.com",
        		"cnameType": "EDGE_HOSTNAME",
        		"stagingCertType": "DEFAULT",
        		"stagingEdgeHostnameID": "ehn_293412",
				"stagingCnameTo": "m-example.com.edgekey.net"
            }
        ]
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/hostnames",
			expectedResponse: &ListActivePropertyHostnamesResponse{
				AccountID:   "act_123",
				ContractID:  "ctr_123",
				GroupID:     "grp_15225",
				PropertyID:  "prp_175780",
				DefaultSort: SortAscending,
				CurrentSort: SortDescending,
				Hostnames: HostnamesResponseItems{
					CurrentItemCount: 2,
					Items: []HostnameItem{
						{
							CnameFrom:                "example.com",
							CnameType:                HostnameCnameTypeEdgeHostname,
							ProductionCertType:       CertTypeDefault,
							ProductionCnameTo:        "example.com.edgekey.net",
							ProductionEdgeHostnameID: "ehn_895822",
						},
						{
							CnameFrom:             "m-example.com",
							CnameType:             HostnameCnameTypeEdgeHostname,
							StagingCertType:       CertTypeDefault,
							StagingEdgeHostnameID: "ehn_293412",
							StagingCnameTo:        "m-example.com.edgekey.net",
						},
					},
				},
			},
		},
		"200 OK - all params": {
			params: ListActivePropertyHostnamesRequest{
				PropertyID:        "prp_175780",
				GroupID:           "grp_15225",
				ContractID:        "ctr_123",
				IncludeCertStatus: false,
				Offset:            0,
				Limit:             1,
				Sort:              "hostname:a",
				Hostname:          "example.com",
				CnameTo:           "example.com",
				Network:           ActivationNetworkProduction,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_123",
    "contractId": "ctr_123",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
	"defaultSort": "hostname:a",
    "currentSort": "hostname:d", 
    "hostnames": {
		"currentItemCount": 2,
        "items": [
            {
                "cnameFrom": "example.com",
        		"cnameType": "EDGE_HOSTNAME",
        		"productionCertType": "DEFAULT",
        		"productionCnameTo": "example.com.edgekey.net",
        		"productionEdgeHostnameID": "ehn_895822"
            },
            {
                "cnameFrom": "m-example.com",
        		"cnameType": "EDGE_HOSTNAME",
        		"stagingCertType": "DEFAULT",
        		"stagingEdgeHostnameID": "ehn_293412",
				"stagingCnameTo": "m-example.com.edgekey.net"
            }
        ],
		"previousLink": "previous link",
		"nextLink": "next link"
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/hostnames?cnameTo=example.com&contractId=ctr_123&groupId=grp_15225&hostname=example.com&limit=1&network=PRODUCTION&sort=hostname%3Aa",
			expectedResponse: &ListActivePropertyHostnamesResponse{
				AccountID:   "act_123",
				ContractID:  "ctr_123",
				GroupID:     "grp_15225",
				PropertyID:  "prp_175780",
				DefaultSort: SortAscending,
				CurrentSort: SortDescending,
				Hostnames: HostnamesResponseItems{
					CurrentItemCount: 2,
					Items: []HostnameItem{
						{
							CnameFrom:                "example.com",
							CnameType:                HostnameCnameTypeEdgeHostname,
							ProductionCertType:       CertTypeDefault,
							ProductionCnameTo:        "example.com.edgekey.net",
							ProductionEdgeHostnameID: "ehn_895822",
						},
						{
							CnameFrom:             "m-example.com",
							CnameType:             HostnameCnameTypeEdgeHostname,
							StagingCertType:       CertTypeDefault,
							StagingEdgeHostnameID: "ehn_293412",
							StagingCnameTo:        "m-example.com.edgekey.net",
						},
					},
					PreviousLink: ptr.To("previous link"),
					NextLink:     ptr.To("next link"),
				},
			},
		},
		"validation error PropertyID missing": {
			params: ListActivePropertyHostnamesRequest{
				Offset: 3,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"validation error Offset negative": {
			params: ListActivePropertyHostnamesRequest{
				PropertyID: "prp_175780",
				Offset:     -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Offset")
			},
		},
		"validation error Limit negative": {
			params: ListActivePropertyHostnamesRequest{
				PropertyID: "prp_175780",
				Limit:      -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Limit")
			},
		},
		"validation error Network invalid": {
			params: ListActivePropertyHostnamesRequest{
				PropertyID: "prp_175780",
				Network:    "invalid_network",
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Network")
			},
		},
		"validation error network missing": {
			params: ListActivePropertyHostnamesRequest{
				Offset: 3,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"validation error Sort method invalid": {
			params: ListActivePropertyHostnamesRequest{
				PropertyID: "prp_175780",
				Sort:       "asc",
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Sort")
			},
		},
		"500 internal server status error": {
			params: ListActivePropertyHostnamesRequest{
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching hostnames",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/hostnames",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching hostnames",
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
			result, err := client.ListActivePropertyHostnames(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapiGetActivePropertyHostnamesDiff(t *testing.T) {
	tests := map[string]struct {
		params           GetActivePropertyHostnamesDiffRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetActivePropertyHostnamesDiffResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - required params": {
			params: GetActivePropertyHostnamesDiffRequest{
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_123",
    "contractId": "ctr_123",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
    "hostnames": {
		"currentItemCount": 2,
        "items": [
            {
                "cnameFrom": "example.com",
        	    "ProductionCnameType": "EDGE_HOSTNAME",
        		"productionCnameTo": "example.com.edgekey.net",
        		"productionEdgeHostnameID": "ehn_895822",
				"productionCertProvisioningType": "CPS_MANAGED"
            },
            {
                "cnameFrom": "m-example.com",
        		"stagingCnameType":	"EDGE_HOSTNAME",
				"stagingCnameTo": "m-example.com.edgekey.net",
        		"stagingEdgeHostnameID": "ehn_293412",
				"stagingCertProvisioningType": "CPS_MANAGED"
            }
        ]
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/hostnames/diff",
			expectedResponse: &GetActivePropertyHostnamesDiffResponse{
				AccountID:  "act_123",
				ContractID: "ctr_123",
				GroupID:    "grp_15225",
				PropertyID: "prp_175780",
				Hostnames: HostnamesDiffResponseItems{
					CurrentItemCount: 2,
					Items: []HostnameDiffItem{
						{
							CnameFrom:                      "example.com",
							ProductionCnameTo:              "example.com.edgekey.net",
							ProductionCnameType:            HostnameCnameTypeEdgeHostname,
							ProductionEdgeHostnameID:       "ehn_895822",
							ProductionCertProvisioningType: CertTypeCPSManaged,
						},
						{
							CnameFrom:                   "m-example.com",
							StagingCnameTo:              "m-example.com.edgekey.net",
							StagingCnameType:            HostnameCnameTypeEdgeHostname,
							StagingEdgeHostnameID:       "ehn_293412",
							StagingCertProvisioningType: CertTypeCPSManaged,
						},
					},
				},
			},
		},
		"200 OK - all params": {
			params: GetActivePropertyHostnamesDiffRequest{
				PropertyID: "prp_175780",
				GroupID:    "grp_15225",
				ContractID: "ctr_123",
				Offset:     1,
				Limit:      1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_123",
    "contractId": "ctr_123",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
    "hostnames": {
		"currentItemCount": 2,
        "items": [
            {
                "cnameFrom": "example.com",
        	    "ProductionCnameType": "EDGE_HOSTNAME",
        		"productionCnameTo": "example.com.edgekey.net",
        		"productionEdgeHostnameID": "ehn_895822",
				"productionCertProvisioningType": "CPS_MANAGED"
            },
            {
                "cnameFrom": "m-example.com",
        		"stagingCnameType":	"EDGE_HOSTNAME",
				"stagingCnameTo": "m-example.com.edgekey.net",
        		"stagingEdgeHostnameID": "ehn_293412",
				"stagingCertProvisioningType": "CPS_MANAGED"
            }
        ],
		"previousLink": "previous link",
		"nextLink": "next link"
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/hostnames/diff?contractId=ctr_123&groupId=grp_15225&limit=1&offset=1",
			expectedResponse: &GetActivePropertyHostnamesDiffResponse{
				AccountID:  "act_123",
				ContractID: "ctr_123",
				GroupID:    "grp_15225",
				PropertyID: "prp_175780",
				Hostnames: HostnamesDiffResponseItems{
					CurrentItemCount: 2,
					Items: []HostnameDiffItem{
						{
							CnameFrom:                      "example.com",
							ProductionCnameTo:              "example.com.edgekey.net",
							ProductionCnameType:            HostnameCnameTypeEdgeHostname,
							ProductionEdgeHostnameID:       "ehn_895822",
							ProductionCertProvisioningType: CertTypeCPSManaged,
						},
						{
							CnameFrom:                   "m-example.com",
							StagingCnameTo:              "m-example.com.edgekey.net",
							StagingCnameType:            HostnameCnameTypeEdgeHostname,
							StagingEdgeHostnameID:       "ehn_293412",
							StagingCertProvisioningType: CertTypeCPSManaged,
						},
					},
					PreviousLink: ptr.To("previous link"),
					NextLink:     ptr.To("next link"),
				},
			},
		},
		"validation error PropertyID missing": {
			params: GetActivePropertyHostnamesDiffRequest{
				Offset: 3,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"validation error Offset negative": {
			params: GetActivePropertyHostnamesDiffRequest{
				PropertyID: "prp_175780",
				Offset:     -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Offset")
			},
		},
		"validation error Limit negative": {
			params: GetActivePropertyHostnamesDiffRequest{
				PropertyID: "prp_175780",
				Limit:      -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Limit")
			},
		},
		"500 internal server status error": {
			params: GetActivePropertyHostnamesDiffRequest{
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching hostnames",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/hostnames/diff",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching hostnames",
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
			result, err := client.GetActivePropertyHostnamesDiff(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
