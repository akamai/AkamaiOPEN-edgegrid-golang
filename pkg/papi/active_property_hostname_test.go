package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
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
		"200 OK - all types": {
			params: ListActivePropertyHostnamesRequest{
				PropertyID:        "prp_175780",
				IncludeCertStatus: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
				{
				  "accountId": "act_123",
				  "availableSort": [
					"hostname:a",
					"hostname:d"
				  ],
				  "contractId": "ctr_123",
				  "currentSort": "hostname:a",
				  "defaultSort": "hostname:a",
				  "groupId": "grp_15225",
				  "hostnames": {
					"currentItemCount": 5,
					"items": [
					  {
						"cnameFrom": "example.com",
						"cnameType": "EDGE_HOSTNAME",
						"productionCertType": "DEFAULT",
						"productionCnameTo": "example.com.edgekey.net",
						"productionEdgeHostnameId": "ehn_895822"
					  },
					  {
						"cnameFrom": "m-example.com",
						"cnameType": "EDGE_HOSTNAME",
						"stagingCertType": "DEFAULT",
						"stagingCnameTo": "m-example.com.edgekey.net",
						"stagingEdgeHostnameId": "ehn_293412"
					  },
					  {
						"certStatus": {
						  "authorization": {
							"dns01": {
							  "result": {
								"message": "dns01 cps dry run cname/TXT incomplete",
								"source": "CPS",
								"timestamp": "2024-07-25T16:17:37Z"
							  },
							  "value": "dummy-unique-value-for-DNS-TXT-record"
							},
							"http01": {
							  "body": "unique http body content",
							  "result": {
								"message": "http01 cps dry run fail reason",
								"source": "CPS",
								"timestamp": "2024-07-25T16:17:37Z"
							  },
							  "url": "/.well-known/acme-challenge/"
							},
							"status": "ATTEMPTING_VALIDATION",
							"validUntil": "2024-07-25T16:17:37Z"
						  },
						  "certExpirationDate": "2024-07-25T16:17:37Z",
						  "production": [
							{
							  "status": "PENDING"
							}
						  ],
						  "staging": [
							{
							  "status": "PENDING"
							}
						  ],
						  "validationCname": {
							"hostname": "_acme-challenge.www.example.com",
							"target": "{token}.www.example.com.akamai-domain.com"
						  }
						},
						"cnameFrom": "example2.com",
						"cnameType": "EDGE_HOSTNAME",
						"stagingCertType": "DEFAULT",
						"stagingCnameTo": "example2.com.edgekey.net",
						"stagingEdgeHostnameId": "ehn_895822"
					  },
					  {
						"ccmCertStatus": {
						  "ecdsaStagingStatus": "DEPLOYED",
						  "rsaStagingStatus": "DEPLOYED"
						},
						"ccmCertificates": {
						  "ecdsaCertId": "98765",
						  "ecdsaCertLink": "/ccm/v1/certificates/98765",
						  "rsaCertId": "12345",
						  "rsaCertLink": "/ccm/v1/certificates/12345"
						},
						"cnameFrom": "www.example-ccm.com",
						"cnameType": "EDGE_HOSTNAME",
						"mtls": {
						  "caSetId": "524125",
						  "caSetLink": "/mtls-edge-truststore/v2/ca-sets/524125",
						  "checkClientOcsp": false,
						  "sendCaSetClient": false
						},
						"stagingCertType": "CCM",
						"stagingEdgeHostnameId": "ehn_7123",
						"tlsConfiguration": {
						  "cipherProfile": "ak-akamai-2020q1",
						  "disallowedTlsVersions": [
							"TLSv1_1",
							"TLSv1"
						  ],
						  "fipsMode": false,
						  "stapleServerOcspResponse": true
						}
					  },
					  {
						"ccmCertStatus": {
						  "ecdsaProductionStatus": "DEPLOYED",
						  "rsaProductionStatus": "DEPLOYED"
						},
						"ccmCertificates": {
						  "ecdsaCertId": "98765",
						  "ecdsaCertLink": "/ccm/v1/certificates/98765",
						  "rsaCertId": "12345",
						  "rsaCertLink": "/ccm/v1/certificates/12345"
						},
						"cnameFrom": "www.example-ccm.com",
						"cnameType": "EDGE_HOSTNAME",
						"mtls": {
						  "caSetId": "524125",
						  "caSetLink": "/mtls-edge-truststore/v2/ca-sets/524125",
						  "checkClientOcsp": false,
						  "sendCaSetClient": false
						},
						"productionCertType": "CCM",
						"productionEdgeHostnameId": "ehn_7123",
						"tlsConfiguration": {
						  "cipherProfile": "ak-akamai-2020q1",
						  "disallowedTlsVersions": [
							"TLSv1_1",
							"TLSv1"
						  ],
						  "fipsMode": false,
						  "stapleServerOcspResponse": true
						}
					  }
					],
					"nextLink": "/papi/v1/properties/prp_175780/hostnames?offset=1&groupId=grp_15225&contractId=ctr_K-0N7RAK71&limit=3",
					"totalItems": 5
				  },
				  "propertyId": "prp_175780",
				  "propertyName": "mytestproperty.com"
				}`,
			expectedPath: "/papi/v1/properties/prp_175780/hostnames?includeCertStatus=true",
			expectedResponse: &ListActivePropertyHostnamesResponse{
				AccountID:    "act_123",
				ContractID:   "ctr_123",
				GroupID:      "grp_15225",
				PropertyID:   "prp_175780",
				PropertyName: "mytestproperty.com",
				DefaultSort:  SortAscending,
				CurrentSort:  SortAscending,
				AvailableSort: []SortOrder{
					SortAscending,
					SortDescending,
				},
				Hostnames: HostnamesResponseItems{
					CurrentItemCount: 5,
					TotalItems:       5,
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
						{
							CnameFrom:             "example2.com",
							CnameType:             HostnameCnameTypeEdgeHostname,
							StagingCertType:       CertTypeDefault,
							StagingEdgeHostnameID: "ehn_895822",
							StagingCnameTo:        "example2.com.edgekey.net",
							CertStatus: &CertStatusItem{
								Production: []StatusItem{
									{
										Status: "PENDING",
									},
								},
								Staging: []StatusItem{
									{
										Status: "PENDING",
									},
								},
								ValidationCname: ValidationCname{
									Hostname: "_acme-challenge.www.example.com",
									Target:   "{token}.www.example.com.akamai-domain.com",
								},
								Authorization: &Authorization{
									DNS01: &DNSAuthorization{
										Result: AuthorizationResult{
											Message:   "dns01 cps dry run cname/TXT incomplete",
											Source:    "CPS",
											Timestamp: test.NewTimeFromString(t, "2024-07-25T16:17:37Z"),
										},
										Value: "dummy-unique-value-for-DNS-TXT-record",
									},
									HTTP01: &HTTPAuthorization{
										Body: "unique http body content",
										Result: AuthorizationResult{
											Message:   "http01 cps dry run fail reason",
											Source:    "CPS",
											Timestamp: test.NewTimeFromString(t, "2024-07-25T16:17:37Z"),
										},
										URL: "/.well-known/acme-challenge/",
									},
									Status:     "ATTEMPTING_VALIDATION",
									ValidUntil: ptr.To(test.NewTimeFromString(t, "2024-07-25T16:17:37Z")),
								},
							},
						},
						{
							CnameFrom:             "www.example-ccm.com",
							CnameType:             HostnameCnameTypeEdgeHostname,
							StagingCertType:       CertTypeCCM,
							StagingEdgeHostnameID: "ehn_7123",
							MTLS: &MTLSResp{
								CASetLink: "/mtls-edge-truststore/v2/ca-sets/524125",
								MTLS: MTLS{
									CASetID:         "524125",
									CheckClientOCSP: false,
									SendCASetClient: false,
								},
							},
							TLSConfiguration: &TLSConfiguration{
								CipherProfile:            "ak-akamai-2020q1",
								DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
								FIPSMode:                 false,
								StapleServerOcspResponse: true, //TODO OCSP
							},
							CCMCertificates: &CCMCertificatesResp{
								CCMCertificates: CCMCertificates{
									ECDSACertID: "98765",
									RSACertID:   "12345",
								},
								ECDSACertLink: "/ccm/v1/certificates/98765",
								RSACertLink:   "/ccm/v1/certificates/12345",
							},
							CCMCertStatus: &CCMCertStatus{
								ECDSAStagingStatus: "DEPLOYED",
								RSAStagingStatus:   "DEPLOYED",
							},
						},
						{
							CnameFrom:                "www.example-ccm.com",
							CnameType:                HostnameCnameTypeEdgeHostname,
							ProductionCertType:       CertTypeCCM,
							ProductionEdgeHostnameID: "ehn_7123",
							MTLS: &MTLSResp{
								CASetLink: "/mtls-edge-truststore/v2/ca-sets/524125",
								MTLS: MTLS{
									CASetID:         "524125",
									CheckClientOCSP: false,
									SendCASetClient: false,
								},
							},
							TLSConfiguration: &TLSConfiguration{
								CipherProfile:            "ak-akamai-2020q1",
								DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
								FIPSMode:                 false,
								StapleServerOcspResponse: true,
							},
							CCMCertificates: &CCMCertificatesResp{
								CCMCertificates: CCMCertificates{
									ECDSACertID: "98765",
									RSACertID:   "12345",
								},
								ECDSACertLink: "/ccm/v1/certificates/98765",
								RSACertLink:   "/ccm/v1/certificates/12345",
							},
							CCMCertStatus: &CCMCertStatus{
								ECDSAProductionStatus: "DEPLOYED",
								RSAProductionStatus:   "DEPLOYED",
							},
						},
					},
					NextLink: ptr.To("/papi/v1/properties/prp_175780/hostnames?offset=1&groupId=grp_15225&contractId=ctr_K-0N7RAK71&limit=3"),
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
			defer mockServer.Close()
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
			defer mockServer.Close()
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

func TestPapiListActivePropertyHostnamesForAccount(t *testing.T) {
	tests := map[string]struct {
		params           ListActiveAccountHostnamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListActiveAccountHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - no params": {
			params:         ListActiveAccountHostnamesRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_123",
    "currentSort": "hostname:a",
    "defaultSort": "hostname:a",
	"availableSort": [
			"hostname:a",
			"hostname:d"
		],
    "hostnames": {
		"currentItemCount": 2,
		"nextLink": "/papi/v1/hostnames?offset=2",
		"totalItems": 23321,
        "items": [
            {
                "cnameFrom": "example-test-prod.com",
                "productionEdgeHostnameId": "ehn_1",
                "productionCertType": "DEFAULT",
                "productionCnameTo": "example-test-prod.com.edgesuite.net",
                "productionCnameType": "EDGE_HOSTNAME",
                "contractId": "ctr_G-1",
                "groupId": "grp_1",
                "propertyId": "prp_1",
                "propertyName": "test_1",
                "propertyType": "TRADITIONAL",
                "productionProductId": "prd_SPM",
                "latestVersion": 1
            },
            {
                "cnameFrom": "example-test-staging.com",
                "stagingEdgeHostnameId": "ehn_2",
                "stagingCertType": "CPS_MANAGED",
                "stagingCnameTo": "example-test-staging.com.edgesuite.net",
                "stagingCnameType": "EDGE_HOSTNAME",
                "contractId": "ctr_G-2",
                "groupId": "grp_2",
                "propertyId": "prp_2",
                "propertyName": "test-2",
                "propertyType": "HOSTNAME_BUCKET",
                "stagingProductId": "prd_Site_Accel",
                "latestVersion": 1
            }
        ]
    }
}
`,
			expectedPath: "/papi/v1/hostnames",
			expectedResponse: &ListActiveAccountHostnamesResponse{
				AccountID:   "act_123",
				CurrentSort: string(SortAscending),
				DefaultSort: string(SortAscending),
				AvailableSort: []string{
					string(SortAscending),
					string(SortDescending),
				},
				Hostnames: ActiveAccountHostnames{
					CurrentItemCount: 2,
					TotalItems:       23321,
					NextLink:         ptr.To("/papi/v1/hostnames?offset=2"),
					Items: []ActiveAccountHostnameItem{
						{
							CnameFrom:                "example-test-prod.com",
							ContractID:               "ctr_G-1",
							GroupID:                  "grp_1",
							LatestVersion:            1,
							ProductionCertType:       ptr.To(string(CertTypeDefault)),
							ProductionCnameTo:        ptr.To("example-test-prod.com.edgesuite.net"),
							ProductionCnameType:      ptr.To("EDGE_HOSTNAME"),
							ProductionEdgeHostnameID: ptr.To("ehn_1"),
							ProductionProductID:      ptr.To("prd_SPM"),
							PropertyID:               "prp_1",
							PropertyName:             "test_1",
							PropertyType:             "TRADITIONAL",
						},
						{
							CnameFrom:             "example-test-staging.com",
							ContractID:            "ctr_G-2",
							GroupID:               "grp_2",
							LatestVersion:         1,
							StagingCertType:       ptr.To(string(CertTypeCPSManaged)),
							StagingCnameTo:        ptr.To("example-test-staging.com.edgesuite.net"),
							StagingCnameType:      ptr.To("EDGE_HOSTNAME"),
							StagingEdgeHostnameID: ptr.To("ehn_2"),
							StagingProductID:      ptr.To("prd_Site_Accel"),
							PropertyID:            "prp_2",
							PropertyName:          "test-2",
							PropertyType:          "HOSTNAME_BUCKET",
						},
					},
				},
			},
		},
		"200 OK - empty list": {
			params:         ListActiveAccountHostnamesRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_123",
    "currentSort": "hostname:a",
    "defaultSort": "hostname:a",
    "hostnames": {
		"currentItemCount": 0,
		"totalItems": 0,
        "items": []
    }
}
`,
			expectedPath: "/papi/v1/hostnames",
			expectedResponse: &ListActiveAccountHostnamesResponse{
				AccountID:   "act_123",
				CurrentSort: string(SortAscending),
				DefaultSort: string(SortAscending),
				Hostnames: ActiveAccountHostnames{
					CurrentItemCount: 0,
					TotalItems:       0,
					Items:            []ActiveAccountHostnameItem{},
				},
			},
		},
		"200 OK - all params": {
			params: ListActiveAccountHostnamesRequest{
				Offset:     0,
				Limit:      2,
				Sort:       SortAscending,
				Hostname:   "example.com",
				CnameTo:    "example.com.edgekey.net",
				Network:    ActivationNetworkProduction,
				ContractID: "ctr_3",
				GroupID:    "grp_3",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_123",
    "currentSort": "hostname:a",
    "defaultSort": "hostname:a",
    "hostnames": {
		"currentItemCount": 2,
        "items": [
            {
                "cnameFrom": "example.com",
                "productionEdgeHostnameId": "ehn_1",
                "productionCertType": "DEFAULT",
                "productionCnameTo": "example.com.edgekey.net",
                "productionCnameType": "EDGE_HOSTNAME",
                "contractId": "ctr_3",
                "groupId": "grp_3",
                "propertyId": "prp_3",
                "propertyName": "test-3",
                "propertyType": "TRADITIONAL",
                "productionProductId": "prd_SPM",
                "latestVersion": 2
            },
            {
                "cnameFrom": "example.com",
                "stagingEdgeHostnameId": "ehn_2",
                "stagingCertType": "CPS_MANAGED",
                "stagingCnameTo": "example.com.edgekey.net",
                "stagingCnameType": "EDGE_HOSTNAME",
                "contractId": "ctr_3",
                "groupId": "grp_3",
                "propertyId": "prp_3",
                "propertyName": "test-3",
                "propertyType": "HOSTNAME_BUCKET",
                "stagingProductId": "prd_Site_Accel",
                "latestVersion": 2
            }
        ],
		"previousLink": "/papi/v1/hostnames?offset=2&limit=2",
		"totalItems": 2323
    }
}
`,
			expectedPath: "/papi/v1/hostnames?cnameTo=example.com.edgekey.net&contractId=ctr_3&groupId=grp_3&hostname=example.com&limit=2&network=PRODUCTION&sort=hostname%3Aa",
			expectedResponse: &ListActiveAccountHostnamesResponse{
				AccountID:   "act_123",
				CurrentSort: string(SortAscending),
				DefaultSort: string(SortAscending),
				Hostnames: ActiveAccountHostnames{
					CurrentItemCount: 2,
					TotalItems:       2323,
					PreviousLink:     ptr.To("/papi/v1/hostnames?offset=2&limit=2"),
					Items: []ActiveAccountHostnameItem{
						{
							CnameFrom:                "example.com",
							ContractID:               "ctr_3",
							GroupID:                  "grp_3",
							LatestVersion:            2,
							ProductionCertType:       ptr.To(string(CertTypeDefault)),
							ProductionCnameTo:        ptr.To("example.com.edgekey.net"),
							ProductionCnameType:      ptr.To("EDGE_HOSTNAME"),
							ProductionEdgeHostnameID: ptr.To("ehn_1"),
							ProductionProductID:      ptr.To("prd_SPM"),
							PropertyID:               "prp_3",
							PropertyName:             "test-3",
							PropertyType:             "TRADITIONAL",
						},
						{
							CnameFrom:             "example.com",
							ContractID:            "ctr_3",
							GroupID:               "grp_3",
							LatestVersion:         2,
							StagingCertType:       ptr.To(string(CertTypeCPSManaged)),
							StagingCnameTo:        ptr.To("example.com.edgekey.net"),
							StagingCnameType:      ptr.To("EDGE_HOSTNAME"),
							StagingEdgeHostnameID: ptr.To("ehn_2"),
							StagingProductID:      ptr.To("prd_Site_Accel"),
							PropertyID:            "prp_3",
							PropertyName:          "test-3",
							PropertyType:          "HOSTNAME_BUCKET",
						},
					},
				},
			},
		},
		"validation error - Offset negative": {
			params: ListActiveAccountHostnamesRequest{
				Offset: -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Offset")
			},
		},
		"validation error - Limit negative": {
			params: ListActiveAccountHostnamesRequest{
				Limit: -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Limit")
			},
		},
		"validation error - Limit exceeds max": {
			params: ListActiveAccountHostnamesRequest{
				Limit: 1000,
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Limit")
			},
		},
		"validation error - invalid Network": {
			params: ListActiveAccountHostnamesRequest{
				Network: "INVALID",
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Network")
			},
		},
		"validation error - invalid Sort": {
			params: ListActiveAccountHostnamesRequest{
				Sort: "invalid",
			},
			withError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Sort")
			},
		},
		"500 internal server error": {
			params:         ListActiveAccountHostnamesRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching hostnames for account",
    "status": 500
}`,
			expectedPath: "/papi/v1/hostnames",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching hostnames for account",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListActiveAccountHostnames(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
