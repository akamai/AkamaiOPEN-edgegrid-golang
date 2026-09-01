package papi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiGetPropertyVersionHostnames(t *testing.T) {
	tests := map[string]struct {
		params           GetPropertyVersionHostnamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPropertyVersionHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK with validation in progress": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_1-1TJZFB",
    "contractId": "ctr_1-1TJZH5",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
    "propertyVersion": 3,
	"propertyName": "mytestproperty.com",
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
                "cnameTo": "m.example.com.edgesuite.net",
                "domainOwnershipVerification": {
                    "challengeTokenExpiryDate": "2024-05-14T05:25:56Z",
                    "status": "VALIDATION_IN_PROGRESS",
                    "validationCname": {
                        "hostname": "validation.hostname.example.com",
                        "target": "validation.target.example.com"
                    },
                    "validationHttp": {
                        "redirectMethod": {
                            "httpRedirectFrom": "http.validation.redirect.from.example.com",
                            "httpRedirectTo": "http.validation.redirect.to.example.com"
                        },
                        "fileContentMethod": {
                            "url": "http://validation.file.content.example.com/validation.txt",
                            "body": "HTTP validation body"
                        }
                    },
                    "validationTxt": {
                        "hostname": "txt.validation.hostname.example.com",
                        "challengeToken": "token"
                    }
                }
            }
        ]
    }
}

`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=false&validateHostnames=false",
			expectedResponse: &GetPropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
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
							DomainOwnershipVerification: &DomainOwnershipVerification{
								ChallengeTokenExpiryDate: ptr.To(time.Date(2024, 5, 14, 5, 25, 56, 0, time.UTC)),
								Status:                   "VALIDATION_IN_PROGRESS",
								ValidationCname: &ValidationCname{
									Hostname: "validation.hostname.example.com",
									Target:   "validation.target.example.com",
								},
								ValidationHTTP: &ValidationHTTP{
									RedirectMethod: RedirectMethod{
										HTTPRedirectFrom: "http.validation.redirect.from.example.com",
										HTTPRedirectTo:   "http.validation.redirect.to.example.com",
									},
									FileContentMethod: FileContentMethod{
										URL:  "http://validation.file.content.example.com/validation.txt",
										Body: "HTTP validation body",
									},
								},
								ValidationTXT: &ValidationTXT{
									Hostname:       "txt.validation.hostname.example.com",
									ChallengeToken: "token",
								},
							},
						},
					},
				},
			},
		},
		"200 OK with authorization": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_12345",
				ContractID:        "ctr_C-0N7RAC7",
				IncludeCertStatus: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_A-CCT5678",
    "contractId": "ctr_C-0N7RAC7",
    "groupId": "grp_12345",
    "propertyId": "prp_175780",
    "propertyVersion": 3,
	"propertyName": "mytestproperty.com",
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
                "cnameTo": "m.example.com.edgesuite.net",
                "certStatus": {
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
                        "hostname": "_acme-challenge.m.example.com",
                        "target": "ac.secret.m.example.com.validate-akdv.net"
                    },
                    "authorization": {
                        "certDomain": "m.example.com",
                        "status": "ATTEMPTING_VALIDATION",
                        "validUntil": "2026-01-22T03:21:43Z",
                        "http01": {
                            "url": "http://m.example.com/.well-known/acme-challenge/anothersecret",
                            "body": "anothersecret.secretsecret",
                            "result": {
                                "message": "Got NXDomain error while getting A records for m.example.com/.well-known/acme-challenge/anothersecret",
                                "timestamp": "2026-01-21T13:42:12Z",
                                "source": "DNS"
                            }
                        },
                        "dns01": {
                            "value": "secret2",
                            "result": {
                                "message": "Got NXDomain error while getting TXT for _acme-challenge.m.example.com",
                                "timestamp": "2026-01-21T13:42:12Z",
                                "source": "DNS"
                            }
                        }
                    }
                }
            }
        ]
    }
}

`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_C-0N7RAC7&groupId=grp_12345&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &GetPropertyVersionHostnamesResponse{
				AccountID:       "act_A-CCT5678",
				ContractID:      "ctr_C-0N7RAC7",
				GroupID:         "grp_12345",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
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
							CertStatus: CertStatusItem{
								Production: []StatusItem{{Status: "PENDING"}},
								Staging:    []StatusItem{{Status: "PENDING"}},
								ValidationCname: ValidationCname{
									Hostname: "_acme-challenge.m.example.com",
									Target:   "ac.secret.m.example.com.validate-akdv.net",
								},
								Authorization: &Authorization{
									DNS01: &DNSAuthorization{
										Result: AuthorizationResult{
											Message:   "Got NXDomain error while getting TXT for _acme-challenge.m.example.com",
											Timestamp: test.NewTimeFromString(t, "2026-01-21T13:42:12Z"),
											Source:    "DNS",
										},
										Value: "secret2",
									},
									HTTP01: &HTTPAuthorization{
										Body: "anothersecret.secretsecret",
										Result: AuthorizationResult{
											Message:   "Got NXDomain error while getting A records for m.example.com/.well-known/acme-challenge/anothersecret",
											Timestamp: test.NewTimeFromString(t, "2026-01-21T13:42:12Z"),
											Source:    "DNS",
										},
										URL: "http://m.example.com/.well-known/acme-challenge/anothersecret",
									},
									Status:     "ATTEMPTING_VALIDATION",
									ValidUntil: ptr.To(test.NewTimeFromString(t, "2026-01-22T03:21:43Z")),
								},
							},
						},
					},
				},
			},
		},
		"200 OK validated": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_1-1TJZFB",
    "contractId": "ctr_1-1TJZH5",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
    "propertyVersion": 3,
	"propertyName": "mytestproperty.com",
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
                "cnameTo": "m.example.com.edgesuite.net",
                "domainOwnershipVerification": {
                    "status": "VALIDATED"
                }
            }
        ]
    }
}

`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=false&validateHostnames=false",
			expectedResponse: &GetPropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
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
							DomainOwnershipVerification: &DomainOwnershipVerification{
								Status: "VALIDATED",
							},
						},
					},
				},
			},
		},
		"200 OK - support for CCM certs in addition to existing cert types": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_C-0N7RAC7",
				IncludeCertStatus: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accountId": "act_A-CCT5678",
  "contractId": "ctr_C-0N7RAC7",
  "etag": "6aed418629b4e5c0",
  "groupId": "grp_54321",
  "hostnames": {
    "items": [
      {
        "certProvisioningType": "CPS_MANAGED",
        "cnameFrom": "example.com",
        "cnameTo": "example.com.edgesuite.net",
        "cnameType": "EDGE_HOSTNAME",
        "edgeHostnameId": "ehn_895824"
      },
      {
        "ccmCertStatus": {
          "ecdsaProductionStatus": "DEPLOYED",
          "ecdsaStagingStatus": "DEPLOYED",
          "rsaProductionStatus": "DEPLOYED",
          "rsaStagingStatus": "DEPLOYED"
        },
        "ccmCertificates": {
          "ecdsaCertId": "98765",
          "ecdsaCertLink": "/ccm/v1/certificates/98765",
          "rsaCertId": "12345",
          "rsaCertLink": "/ccm/v1/certificates/12345"
        },
        "certProvisioningType": "CCM",
        "cnameFrom": "www.example-ccm.com",
        "cnameTo": "example.com.edgesuite.net",
        "cnameType": "EDGE_HOSTNAME",
        "edgeHostnameId": "ehn_7123",
        "mtls": {
          "caSetId": "524125",
          "caSetLink": "/mtls-edge-truststore/v2/ca-sets/524125",
          "checkClientOcsp": false,
          "sendCaSetClient": false
        },
        "tlsConfiguration": {
          "cipherProfile": "ak-akamai-2020q1",
          "disallowedTlsVersions": [
            "TLSv1_1",
            "TLSv1"
          ],
          "stapleServerOcspResponse": true,
          "fipsMode": false
        }
      }
    ]
  },
  "propertyId": "prp_123456",
  "propertyName": "mytestproperty.com",
  "propertyVersion": 1
}`,
			expectedPath: "/papi/v1/properties/prp_123456/versions/3/hostnames?contractId=ctr_C-0N7RAC7&groupId=grp_54321&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &GetPropertyVersionHostnamesResponse{
				AccountID:       "act_A-CCT5678",
				ContractID:      "ctr_C-0N7RAC7",
				GroupID:         "grp_54321",
				PropertyID:      "prp_123456",
				PropertyVersion: 1,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CertProvisioningType: "CPS_MANAGED",
							CnameFrom:            "example.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895824",
						},
						{
							CertProvisioningType: "CCM",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_7123",
							CCMCertStatus: &CCMCertStatus{
								ECDSAProductionStatus: "DEPLOYED",
								ECDSAStagingStatus:    "DEPLOYED",
								RSAProductionStatus:   "DEPLOYED",
								RSAStagingStatus:      "DEPLOYED",
							},
							CCMCertificates: &CCMCertificatesResp{
								CCMCertificates: CCMCertificates{
									ECDSACertID: "98765",
									RSACertID:   "12345",
								},
								ECDSACertLink: "/ccm/v1/certificates/98765",
								RSACertLink:   "/ccm/v1/certificates/12345",
							},
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
								StapleServerOcspResponse: true,
								FIPSMode:                 false,
							},
						},
					},
				},
			},
		},
		"302 - missed contractId and groupId": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				IncludeCertStatus: false,
			},
			responseStatus: http.StatusFound,
			responseBody: `
{
    "redirectLink": "/papi/v1/properties/prp_175780/versions/1/hostnames?groupId=12345&contractId=G-12RS3N4"
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?includeCertStatus=false&validateHostnames=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					RedirectLink: ptr.To("/papi/v1/properties/prp_175780/versions/1/hostnames?groupId=12345&contractId=G-12RS3N4"),
					StatusCode:   http.StatusFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error PropertyID missing": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyVersion: 3,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"validation error PropertyVersion missing": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID: "prp_175780",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyVersion")
			},
		},
		"404 not found error": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: false,
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
	"type": "not_found",
    "title": "Not Found",
    "detail": "The requested resource was not found",
    "status": 404
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=false&validateHostnames=false",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%s: %w: %s", ErrGetPropertyVersionHostnames, ErrNotFound, "request failed")
				assert.True(t, errors.Is(err, ErrNotFound), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server status error": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching hostnames",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?includeCertStatus=false&validateHostnames=false",
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
			result, err := client.GetPropertyVersionHostnames(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapiUpdatePropertyVersionHostnames(t *testing.T) {
	tests := map[string]struct {
		params           UpdatePropertyVersionHostnamesRequest
		requestBody      string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdatePropertyVersionHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK with validation in progress": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CnameType:            "EDGE_HOSTNAME",
						CnameFrom:            "m.example.com",
						CnameTo:              "example.com.edgekey.net",
						CertProvisioningType: "DEFAULT",
					},
					{
						CnameType:            "EDGE_HOSTNAME",
						EdgeHostnameID:       "ehn_895824",
						CnameFrom:            "example3.com",
						CertProvisioningType: "CPS_MANAGED",
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_1-1TJZFB",
    "contractId": "ctr_1-1TJZH5",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
    "propertyVersion": 3,
	"propertyName": "mytestproperty.com",
    "etag": "6aed418629b4e5c0",
    "hostnames": {
        "items": [
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895822",
                "cnameFrom": "m.example.com",
                "cnameTo": "example.com.edgekey.net",
                "certProvisioningType": "DEFAULT",
                "certStatus": {
                    "validationCname": {
                        "hostname": "_acme-challenge.www.example.com",
                        "target": "{token}.www.example.com.akamai-domain.com"
                    },
                    "staging": [
                        {
                            "status": "NEEDS_VALIDATION"
                        }
                    ],
                    "production": [
                        {
                            "status": "NEEDS_VALIDATION"
                        }
                    ]
                }
            },
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895833",
                "cnameFrom": "example3.com",
                "cnameTo": "m.example.com.edgesuite.net",
 				"certProvisioningType": "CPS_MANAGED",
				"domainOwnershipVerification": {
					"challengeTokenExpiryDate": "2024-05-14T05:25:56Z",
					"status": "VALIDATION_IN_PROGRESS",
					"validationCname": {
						"hostname": "validation.hostname.example.com",
						"target": "validation.target.example.com"
					},
					"validationHttp": {
						"redirectMethod": {
							"httpRedirectFrom": "http.validation.redirect.from.example.com",
							"httpRedirectTo": "http.validation.redirect.to.example.com"
						},
						"fileContentMethod": {
							"url": "http://validation.file.content.example.com/validation.txt",
							"body": "HTTP validation body"
						}
					},
					"validationTxt": {
						"hostname": "txt.validation.hostname.example.com",
						"challengeToken": "token"
					}
				}	
            }
        ]
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				PropertyName:    "mytestproperty.com",
				Etag:            "6aed418629b4e5c0",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895822",
							CnameFrom:            "m.example.com",
							CnameTo:              "example.com.edgekey.net",
							CertProvisioningType: "DEFAULT",
							CertStatus: CertStatusItem{
								ValidationCname: ValidationCname{
									Hostname: "_acme-challenge.www.example.com",
									Target:   "{token}.www.example.com.akamai-domain.com",
								},
								Staging:    []StatusItem{{Status: "NEEDS_VALIDATION"}},
								Production: []StatusItem{{Status: "NEEDS_VALIDATION"}},
							},
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895833",
							CnameFrom:            "example3.com",
							CnameTo:              "m.example.com.edgesuite.net",
							CertProvisioningType: "CPS_MANAGED",
							DomainOwnershipVerification: &DomainOwnershipVerification{
								ChallengeTokenExpiryDate: ptr.To(time.Date(2024, 5, 14, 5, 25, 56, 0, time.UTC)),
								Status:                   "VALIDATION_IN_PROGRESS",
								ValidationCname: &ValidationCname{
									Hostname: "validation.hostname.example.com",
									Target:   "validation.target.example.com",
								},
								ValidationHTTP: &ValidationHTTP{
									RedirectMethod: RedirectMethod{
										HTTPRedirectFrom: "http.validation.redirect.from.example.com",
										HTTPRedirectTo:   "http.validation.redirect.to.example.com",
									},
									FileContentMethod: FileContentMethod{
										URL:  "http://validation.file.content.example.com/validation.txt",
										Body: "HTTP validation body",
									},
								},
								ValidationTXT: &ValidationTXT{
									Hostname:       "txt.validation.hostname.example.com",
									ChallengeToken: "token",
								},
							},
						},
					},
				},
			},
		},
		"200 OK with authorization data": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_12345",
				ContractID:        "ctr_C-0N7RAC7",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CnameType:            "EDGE_HOSTNAME",
						CnameFrom:            "m.example.com",
						CnameTo:              "example.com.edgekey.net",
						CertProvisioningType: "DEFAULT",
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_A-CCT5678",
    "contractId": "ctr_C-0N7RAC7",
    "groupId": "grp_12345",
    "propertyId": "prp_175780",
    "propertyVersion": 3,
	"propertyName": "mytestproperty.com",
    "etag": "6aed418629b4e5c0",
    "hostnames": {
        "items": [
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895822",
                "cnameFrom": "m.example.com",
                "cnameTo": "example.com.edgekey.net",
                "certProvisioningType": "DEFAULT",
                "certStatus": {
                    "validationCname": {
                        "hostname": "_acme-challenge.www.example.com",
                        "target": "{token}.www.example.com.akamai-domain.com"
                    },
                    "staging": [
                        {
                            "status": "NEEDS_VALIDATION"
                        }
                    ],
                    "production": [
                        {
                            "status": "NEEDS_VALIDATION"
                        }
                    ],
                    "authorization": {
                        "certDomain": "m.example.com",
                        "status": "ATTEMPTING_VALIDATION",
                        "validUntil": "2026-01-22T03:21:43Z",
                        "http01": {
                            "url": "http://m.example.com/.well-known/acme-challenge/anothersecret",
                            "body": "anothersecret.secretsecret",
                            "result": {
                                "message": "Got NXDomain error while getting A records for m.example.com/.well-known/acme-challenge/anothersecret",
                                "timestamp": "2026-01-21T13:42:12Z",
                                "source": "DNS"
                            }
                        },
                        "dns01": {
                            "value": "secret2",
                            "result": {
                                "message": "Got NXDomain error while getting TXT for _acme-challenge.m.example.com",
                                "timestamp": "2026-01-21T13:42:12Z",
                                "source": "DNS"
                            }
                        }
                    }
                }
            }
        ]
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_C-0N7RAC7&groupId=grp_12345&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_A-CCT5678",
				ContractID:      "ctr_C-0N7RAC7",
				GroupID:         "grp_12345",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				PropertyName:    "mytestproperty.com",
				Etag:            "6aed418629b4e5c0",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895822",
							CnameFrom:            "m.example.com",
							CnameTo:              "example.com.edgekey.net",
							CertProvisioningType: "DEFAULT",
							CertStatus: CertStatusItem{
								ValidationCname: ValidationCname{
									Hostname: "_acme-challenge.www.example.com",
									Target:   "{token}.www.example.com.akamai-domain.com",
								},
								Staging:    []StatusItem{{Status: "NEEDS_VALIDATION"}},
								Production: []StatusItem{{Status: "NEEDS_VALIDATION"}},
								Authorization: &Authorization{
									DNS01: &DNSAuthorization{
										Result: AuthorizationResult{
											Message:   "Got NXDomain error while getting TXT for _acme-challenge.m.example.com",
											Timestamp: test.NewTimeFromString(t, "2026-01-21T13:42:12Z"),
											Source:    "DNS",
										},
										Value: "secret2",
									},
									HTTP01: &HTTPAuthorization{
										Body: "anothersecret.secretsecret",
										Result: AuthorizationResult{
											Message:   "Got NXDomain error while getting A records for m.example.com/.well-known/acme-challenge/anothersecret",
											Timestamp: test.NewTimeFromString(t, "2026-01-21T13:42:12Z"),
											Source:    "DNS",
										},
										URL: "http://m.example.com/.well-known/acme-challenge/anothersecret",
									},
									Status:     "ATTEMPTING_VALIDATION",
									ValidUntil: ptr.To(test.NewTimeFromString(t, "2026-01-22T03:21:43Z")),
								},
							},
						},
					},
				},
			},
		},
		"200 OK validated": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CnameType:            "EDGE_HOSTNAME",
						EdgeHostnameID:       "ehn_895824",
						CnameFrom:            "example3.com",
						CertProvisioningType: "CPS_MANAGED",
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "act_1-1TJZFB",
    "contractId": "ctr_1-1TJZH5",
    "groupId": "grp_15225",
    "propertyId": "prp_175780",
    "propertyVersion": 3,
	"propertyName": "mytestproperty.com",
    "etag": "6aed418629b4e5c0",
    "hostnames": {
        "items": [
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895833",
                "cnameFrom": "example3.com",
                "cnameTo": "m.example.com.edgesuite.net",
 				"certProvisioningType": "CPS_MANAGED",
				"domainOwnershipVerification": {
					"status": "VALIDATED"
				}	
            }
        ]
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				PropertyName:    "mytestproperty.com",
				Etag:            "6aed418629b4e5c0",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895833",
							CnameFrom:            "example3.com",
							CnameTo:              "m.example.com.edgesuite.net",
							CertProvisioningType: "CPS_MANAGED",
							DomainOwnershipVerification: &DomainOwnershipVerification{
								Status: "VALIDATED",
							},
						},
					},
				},
			},
		},
		"200 OK updating hostnames with CCM certs in addition to existing cert types": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CPS_MANAGED",
						CnameFrom:            "example.com",
						CnameType:            "EDGE_HOSTNAME",
						EdgeHostnameID:       "ehn_895824",
					},
					{
						CertProvisioningType: "DEFAULT",
						CnameFrom:            "example-tst.com",
						CnameTo:              "example-tst.com.edgekey.net",
						CnameType:            "EDGE_HOSTNAME",
					},
					{
						CertProvisioningType: "CCM",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						CCMCertificates: &CCMCertificates{
							RSACertID:   "12345",
							ECDSACertID: "98765",
						},
						MTLS: &MTLS{
							CASetID:         "524125",
							CheckClientOCSP: false,
							SendCASetClient: false,
						},
						TLSConfiguration: &TLSConfiguration{
							CipherProfile:            "ak-akamai-2020q1",
							DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
							StapleServerOcspResponse: true,
							FIPSMode:                 false,
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			requestBody:    `[{"cnameType":"EDGE_HOSTNAME","edgeHostnameId":"ehn_895824","cnameFrom":"example.com","certProvisioningType":"CPS_MANAGED"},{"cnameType":"EDGE_HOSTNAME","cnameFrom":"example-tst.com","cnameTo":"example-tst.com.edgekey.net","certProvisioningType":"DEFAULT"},{"cnameType":"EDGE_HOSTNAME","cnameFrom":"www.example-ccm.com","cnameTo":"example.com.edgesuite.net","certProvisioningType":"CCM","ccmCertificates":{"ecdsaCertId":"98765","rsaCertId":"12345"},"mtls":{"caSetId":"524125"},"tlsConfiguration":{"cipherProfile":"ak-akamai-2020q1","disallowedTlsVersions":["TLSv1_1","TLSv1"],"stapleServerOcspResponse":true}}]`,
			responseBody: `
{
  "accountId": "act_B-G-12RS3M4",
  "contractId": "ctr_1-2ABCD3",
  "etag": "6aed418629b4e5c0",
  "groupId": "grp_54321",
  "hostnames": {
    "items": [
      {
        "certProvisioningType": "CPS_MANAGED",
        "cnameFrom": "example.com",
        "cnameTo": "example.com.edgesuite.net",
        "cnameType": "EDGE_HOSTNAME",
        "edgeHostnameId": "ehn_895824"
      },
      {
        "ccmCertStatus": {
          "ecdsaProductionStatus": "PENDING",
          "ecdsaStagingStatus": "PENDING",
          "rsaProductionStatus": "PENDING",
          "rsaStagingStatus": "PENDING"
        },
        "ccmCertificates": {
          "ecdsaCertId": "98765",
          "ecdsaCertLink": "/ccm/v1/certificates/98765",
          "rsaCertId": "12345",
          "rsaCertLink": "/ccm/v1/certificates/12345"
        },
        "certProvisioningType": "CCM",
        "cnameFrom": "www.example-ccm.com",
        "cnameTo": "example.com.edgesuite.net",
        "cnameType": "EDGE_HOSTNAME",
        "edgeHostnameId": "ehn_7123",
        "mtls": {
          "caSetId": "524125",
          "caSetLink": "/mtls-edge-truststore/v2/ca-sets/524125",
          "checkClientOcsp": false,
          "sendCaSetClient": false
        },
        "tlsConfiguration": {
          "cipherProfile": "ak-akamai-2020q1",
          "disallowedTlsVersions": [
            "TLSv1_1",
            "TLSv1"
          ],
          "stapleServerOcspResponse": true,
          "fipsMode": false
        }
      }
    ]
  },
  "propertyId": "prp_123456",
  "propertyName": "mytestproperty.com",
  "propertyVersion": 1
}
`,
			expectedPath: "/papi/v1/properties/prp_123456/versions/3/hostnames?contractId=ctr_1-2ABCD3&groupId=grp_54321&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_B-G-12RS3M4",
				ContractID:      "ctr_1-2ABCD3",
				GroupID:         "grp_54321",
				PropertyID:      "prp_123456",
				PropertyVersion: 1,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CertProvisioningType: "CPS_MANAGED",
							CnameFrom:            "example.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895824",
						},
						{
							CertProvisioningType: "CCM",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_7123",
							CCMCertStatus: &CCMCertStatus{
								ECDSAProductionStatus: "PENDING",
								ECDSAStagingStatus:    "PENDING",
								RSAProductionStatus:   "PENDING",
								RSAStagingStatus:      "PENDING",
							},
							CCMCertificates: &CCMCertificatesResp{
								CCMCertificates: CCMCertificates{
									ECDSACertID: "98765",
									RSACertID:   "12345",
								},
								ECDSACertLink: "/ccm/v1/certificates/98765",
								RSACertLink:   "/ccm/v1/certificates/12345",
							},
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
								StapleServerOcspResponse: true,
								FIPSMode:                 false,
							},
						},
					},
				},
			},
		},
		"200 empty hostnames": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: true,
				Hostnames:         []Hostname{{}},
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
	"propertyName": "mytestproperty.com",
    "hostnames": {
        "items": []
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{},
				},
			},
		},
		"200 VerifyHostnames true empty hostnames": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				ValidateHostnames: true,
				IncludeCertStatus: true,
				Hostnames:         []Hostname{{}},
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
	"propertyName": "mytestproperty.com",
	"validateHostnames": true,
    "hostnames": {
        "items": []
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=true&validateHostnames=true",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{},
				},
			},
		},
		"302 - missed contractId and groupId": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CnameType:            "EDGE_HOSTNAME",
						CnameFrom:            "m.example.com",
						CnameTo:              "example.com.edgekey.net",
						CertProvisioningType: "DEFAULT",
					},
					{
						CnameType:            "EDGE_HOSTNAME",
						EdgeHostnameID:       "ehn_895824",
						CnameFrom:            "example3.com",
						CertProvisioningType: "CPS_MANAGED",
					},
				},
			},
			responseStatus: http.StatusFound,
			responseBody: `
{
    "redirectLink": "/papi/v1/properties/prp_175780/versions/3/hostnames?groupId=12345&contractId=G-12RS3N4"
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?includeCertStatus=true&validateHostnames=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					RedirectLink: ptr.To("/papi/v1/properties/prp_175780/versions/3/hostnames?groupId=12345&contractId=G-12RS3N4"),
					StatusCode:   http.StatusFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error PropertyID missing": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyVersion: 3,
				Hostnames:       []Hostname{{}},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"validation error PropertyVersion missing": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID: "prp_175780",
				Hostnames:  []Hostname{{}},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyVersion")
			},
		},
		"validation error - CCMCertificates RSACertID, ECDSACertID and MTLS CASetID are not a digit": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CCM",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						CCMCertificates: &CCMCertificates{
							RSACertID:   "12345a",
							ECDSACertID: "98765a",
						},
						MTLS: &MTLS{
							CASetID:         "524125a",
							CheckClientOCSP: false,
							SendCASetClient: false,
						},
						TLSConfiguration: &TLSConfiguration{
							CipherProfile:            "ak-akamai-2020q1",
							DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
							StapleServerOcspResponse: true,
							FIPSMode:                 false,
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "{\n\tCCMCertificates: {\n\t\tECDSACertID: must contain digits only\n\t\tRSACertID: must contain digits only\n\t}\n\tMTLS: {\n\t\tCASetID: must contain digits only\n\t}\n}")
			},
		},
		"validation error - one of RSACertID and ECDSACertID must be provided for CCM": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CCM",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						CCMCertificates:      &CCMCertificates{},
						MTLS: &MTLS{
							CASetID:         "524125",
							CheckClientOCSP: false,
							SendCASetClient: false,
						},
						TLSConfiguration: &TLSConfiguration{
							CipherProfile:            "ak-akamai-2020q1",
							DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
							StapleServerOcspResponse: true,
							FIPSMode:                 false,
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "either RSACertID or ECDSACertID must be provided")
			},
		},
		"validation error - CCMCertificates is required for CCM": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CCM",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						MTLS: &MTLS{
							CASetID:         "524125",
							CheckClientOCSP: false,
							SendCASetClient: false,
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "struct validation: Hostnames[0]: {\n\tValidateCCMHostname: when using `certProvisioningType` set to `CCM`, the request body must contain `ccmCertificates` with at least `rsaCertId` or `ecdsaCertId`")
			},
		},
		"validation error - MTLS is only valid for CCM": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CPS_MANAGED",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						MTLS: &MTLS{
							CASetID:         "524125",
							CheckClientOCSP: false,
							SendCASetClient: false,
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "the mTLS configuration is provided without `certProvisioningType` set to `CCM`")
			},
		},
		"validation error - TLSConfiguration is only valid for CCM": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CPS_MANAGED",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						TLSConfiguration: &TLSConfiguration{
							CipherProfile:            "ak-akamai-2020q1",
							DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
							StapleServerOcspResponse: true,
							FIPSMode:                 false,
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "the TLS configuration is provided without `certProvisioningType` set to `CCM`")
			},
		},
		"validation error - CCMCertificates is only valid for CCM": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CPS_MANAGED",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						CCMCertificates: &CCMCertificates{
							RSACertID:   "12345",
							ECDSACertID: "98765",
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "the CCM cert details are provided without `certProvisioningType` set to `CCM`")
			},
		},
		"validation error - missed CASetID when MTLS provided": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:      "prp_123456",
				PropertyVersion: 3,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CCM",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						CCMCertificates: &CCMCertificates{
							RSACertID: "12345",
						},
						MTLS: &MTLS{},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "CASetID: cannot be blank")
			},
		},
		"validation error - empty cipher profile": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CertProvisioningType: "CCM",
						CnameFrom:            "www.example-ccm.com",
						CnameTo:              "example.com.edgesuite.net",
						CnameType:            "EDGE_HOSTNAME",
						CCMCertificates: &CCMCertificates{
							RSACertID:   "12345",
							ECDSACertID: "98765",
						},
						MTLS: &MTLS{
							CASetID:         "524125",
							CheckClientOCSP: false,
							SendCASetClient: false,
						},
						TLSConfiguration: &TLSConfiguration{
							CipherProfile:            "",
							DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
							StapleServerOcspResponse: true,
							FIPSMode:                 false,
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "{\n\tValidateCCMHostname: the cipher profile is empty in the TLS configuration\n}")
			},
		},
		"200 Hostnames missing": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: true,
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
	"propertyName": "mytestproperty.com",
	"validateHostnames": false,
	"hostnames": {
		"items": []
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{},
				},
			},
		},
		"200 Hostnames items missing": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				Hostnames:         nil,
				IncludeCertStatus: true,
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
	"propertyName": "mytestproperty.com",
	"validateHostnames": false,
	"hostnames": {
		"items": []
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{},
				},
			},
		},
		"200 Hostnames items empty": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: true,
				Hostnames:         []Hostname{},
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
	"propertyName": "mytestproperty.com",
	"validateHostnames": false,
	"hostnames": {
		"items": []
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=true&validateHostnames=false",
			expectedResponse: &UpdatePropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				PropertyName:    "mytestproperty.com",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{},
				},
			},
		},
		"400 Hostnames cert type is invalid": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				GroupID:           "grp_15225",
				ContractID:        "ctr_1-1TJZH5",
				IncludeCertStatus: true,
				Hostnames: []Hostname{
					{
						CnameType:            "EDGE_HOSTNAME",
						CnameFrom:            "m.example.com",
						CnameTo:              "example.com.edgesuite.net",
						CertProvisioningType: "INVALID_TYPE",
					},
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
    "type": "https://problems.luna.akamaiapis.net/papi/v0/json-mapping-error",
    "title": "Unable to interpret JSON",
    "detail": "Your input could not be interpreted as the expected JSON format. Cannot deserialize value of type com.akamai.platformtk.entities.HostnameRelation$CertProvisioningType from String INVALID_TYPE: not one of the values accepted for Enum class: [DEFAULT, CPS_MANAGED]\n at [Source: (org.apache.catalina.connector.CoyoteInputStream); line: 6, column: 41] (through reference chain: java.util.ArrayList[0]->com.akamai.luna.papi.model.HostnameItem[certProvisioningType]).",
    "status": 400
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=true&validateHostnames=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "https://problems.luna.akamaiapis.net/papi/v0/json-mapping-error",
					Title:      "Unable to interpret JSON",
					Detail:     "Your input could not be interpreted as the expected JSON format. Cannot deserialize value of type com.akamai.platformtk.entities.HostnameRelation$CertProvisioningType from String INVALID_TYPE: not one of the values accepted for Enum class: [DEFAULT, CPS_MANAGED]\n at [Source: (org.apache.catalina.connector.CoyoteInputStream); line: 6, column: 41] (through reference chain: java.util.ArrayList[0]->com.akamai.luna.papi.model.HostnameItem[certProvisioningType]).",
					StatusCode: http.StatusBadRequest,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server status error": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:        "prp_175780",
				PropertyVersion:   3,
				Hostnames:         []Hostname{{}},
				IncludeCertStatus: true,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating hostnames",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?includeCertStatus=true&validateHostnames=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error updating hostnames",
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
				assert.Equal(t, http.MethodPut, r.Method)
				if test.requestBody != "" {
					buf := new(bytes.Buffer)
					_, err := buf.ReadFrom(r.Body)
					assert.NoError(t, err)
					req := buf.String()
					assert.Equal(t, test.requestBody, req)
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdatePropertyVersionHostnames(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapiPatchPropertyVersionHostnames(t *testing.T) {
	tests := map[string]struct {
		params              PatchPropertyVersionHostnamesRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *PatchPropertyVersionHostnamesResponse
		withError           func(*testing.T, error)
	}{
		"200 OK - only add": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CnameType: HostnameCnameTypeEdgeHostname,
							CnameFrom: "m.example.com",
							CnameTo:   "example.com.edgekey.net",
						},
						{
							CnameType:            HostnameCnameTypeEdgeHostname,
							EdgeHostnameID:       "ehn_666",
							CnameFrom:            "example3.com",
							CertProvisioningType: CertTypeDefault,
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId": "act_789",
				"contractId": "ctr_456",
				"groupId": "grp_321",
				"propertyId": "prp_123",
				"propertyName": "test-property",
				"propertyVersion": 1,
				"etag": "123abc456def7890",
				"propertyName": "test-property",
				"hostnames": {
					"items": [
						{
							"cnameType": "EDGE_HOSTNAME",
							"edgeHostnameId": "ehn_555",
							"cnameFrom": "m.example.com",
							"cnameTo": "example.com.edgekey.net",
							"certProvisioningType": "CPS_MANAGED"
						},
						{
							"cnameType": "EDGE_HOSTNAME",
							"edgeHostnameId": "ehn_666",
							"cnameFrom": "example3.com",
							"cnameTo": "m.example.com.edgesuite.net",
							"certProvisioningType": "DEFAULT"
						}
					]
				}
			}`,
			expectedPath:        "/papi/v1/properties/prp_123/versions/1/hostnames",
			expectedRequestBody: `{"add":[{"cnameFrom":"m.example.com","cnameType":"EDGE_HOSTNAME","cnameTo":"example.com.edgekey.net"},{"cnameFrom":"example3.com","cnameType":"EDGE_HOSTNAME","certProvisioningType":"DEFAULT","edgeHostnameId":"ehn_666"}]}`,
			expectedResponse: &PatchPropertyVersionHostnamesResponse{
				AccountID:       "act_789",
				ContractID:      "ctr_456",
				GroupID:         "grp_321",
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				PropertyName:    "test-property",
				Etag:            "123abc456def7890",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "m.example.com",
							CnameTo:              "example.com.edgekey.net",
							EdgeHostnameID:       "ehn_555",
							CertProvisioningType: "CPS_MANAGED",
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_666",
							CnameFrom:            "example3.com",
							CertProvisioningType: "DEFAULT",
							CnameTo:              "m.example.com.edgesuite.net",
						},
					},
				},
			},
		},
		"200 OK - support for adding hostnames with CCM certs in addition to the existing cert types": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_12345",
				PropertyVersion: 1,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: CertTypeCPSManaged,
							CnameFrom:            "example.com",
							CnameType:            HostnameCnameTypeEdgeHostname,
							EdgeHostnameID:       "ehn_895824",
						},
						{
							CertProvisioningType: CertTypeDefault,
							CnameFrom:            "example-tst.com",
							CnameTo:              "example-tst.com.edgekey.net",
							CnameType:            HostnameCnameTypeEdgeHostname,
						},
						{
							CertProvisioningType: CertTypeCCM,
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            HostnameCnameTypeEdgeHostname,
							CCMCertificates: &CCMCertificates{
								RSACertID:   "12345",
								ECDSACertID: "98765",
							},
							MTLS: &MTLS{
								CASetID:         "524125",
								CheckClientOCSP: false,
								SendCASetClient: false,
							},
							TLSConfiguration: &TLSConfiguration{
								CipherProfile:            "ak-akamai-2020q1",
								DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
								StapleServerOcspResponse: true,
								FIPSMode:                 false,
							},
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accountId": "act_A-CCT5678",
  "contractId": "ctr_K-0N1RAK23",
  "etag": "6aed418629b4e5c0",
  "groupId": "grp_12345",
  "hostnames": {
    "items": [
      {
        "certProvisioningType": "CPS_MANAGED",
        "cnameFrom": "example.com",
        "cnameTo": "example.com.edgesuite.net",
        "cnameType": "EDGE_HOSTNAME",
        "edgeHostnameId": "ehn_895824"
      },
      {
        "ccmCertStatus": {
          "ecdsaProductionStatus": "PENDING",
          "ecdsaStagingStatus": "PENDING",
          "rsaProductionStatus": "PENDING",
          "rsaStagingStatus": "PENDING"
        },
        "ccmCertificates": {
          "ecdsaCertId": "98765",
          "ecdsaCertLink": "/ccm/v1/certificates/98765",
          "rsaCertId": "12345",
          "rsaCertLink": "/ccm/v1/certificates/12345"
        },
        "certProvisioningType": "CCM",
        "cnameFrom": "www.example-ccm.com",
        "cnameTo": "example.com.edgesuite.net",
        "cnameType": "EDGE_HOSTNAME",
        "edgeHostnameId": "ehn_7123",
        "mtls": {
          "caSetId": "524125",
          "caSetLink": "/mtls-edge-truststore/v2/ca-sets/524125",
          "checkClientOcsp": false,
          "sendCaSetClient": false
        },
        "tlsConfiguration": {
          "cipherProfile": "ak-akamai-2020q1",
          "disallowedTlsVersions": [
            "TLSv1_1",
            "TLSv1"
          ],
          "stapleServerOcspResponse": true,
          "fipsMode": false
        }
      }
    ]
  },
  "propertyId": "prp_12345",
  "propertyName": "mytestproperty.com",
  "propertyVersion": 1
}`,
			expectedPath:        "/papi/v1/properties/prp_12345/versions/1/hostnames",
			expectedRequestBody: `{"add":[{"cnameFrom":"example.com","cnameType":"EDGE_HOSTNAME","certProvisioningType":"CPS_MANAGED","edgeHostnameId":"ehn_895824"},{"cnameFrom":"example-tst.com","cnameType":"EDGE_HOSTNAME","cnameTo":"example-tst.com.edgekey.net","certProvisioningType":"DEFAULT"},{"cnameFrom":"www.example-ccm.com","cnameType":"EDGE_HOSTNAME","cnameTo":"example.com.edgesuite.net","certProvisioningType":"CCM","mtls":{"caSetId":"524125"},"tlsConfiguration":{"cipherProfile":"ak-akamai-2020q1","disallowedTlsVersions":["TLSv1_1","TLSv1"],"stapleServerOcspResponse":true},"ccmCertificates":{"ecdsaCertId":"98765","rsaCertId":"12345"}}]}`,
			expectedResponse: &PatchPropertyVersionHostnamesResponse{
				AccountID:       "act_A-CCT5678",
				ContractID:      "ctr_K-0N1RAK23",
				GroupID:         "grp_12345",
				Etag:            "6aed418629b4e5c0",
				PropertyID:      "prp_12345",
				PropertyName:    "mytestproperty.com",
				PropertyVersion: 1,
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895824",
							CnameFrom:            "example.com",
							CnameTo:              "example.com.edgesuite.net",
							CertProvisioningType: "CPS_MANAGED",
							CertStatus: CertStatusItem{
								ValidationCname: ValidationCname{},
							},
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_7123",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CertProvisioningType: "CCM",
							CertStatus: CertStatusItem{
								ValidationCname: ValidationCname{},
							},
							CCMCertStatus: &CCMCertStatus{
								ECDSAProductionStatus: "PENDING",
								ECDSAStagingStatus:    "PENDING",
								RSAProductionStatus:   "PENDING",
								RSAStagingStatus:      "PENDING",
							},
							CCMCertificates: &CCMCertificatesResp{
								CCMCertificates: CCMCertificates{
									ECDSACertID: "98765",
									RSACertID:   "12345",
								},
								ECDSACertLink: "/ccm/v1/certificates/98765",
								RSACertLink:   "/ccm/v1/certificates/12345",
							},
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
								StapleServerOcspResponse: true,
								FIPSMode:                 false,
							},
						},
					},
				},
			},
		},
		"200 OK - only remove": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Remove: []string{
						"ehn_555",
						"ehn_666",
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId": "act_789",
				"contractId": "ctr_456",
				"groupId": "grp_321",
				"propertyId": "prp_123",
				"propertyName": "test-property",
				"propertyVersion": 1,
				"etag": "123abc456def7890",
				"hostnames": {
					"items": []
				}
			}`,
			expectedPath:        "/papi/v1/properties/prp_123/versions/1/hostnames",
			expectedRequestBody: `{"remove":["ehn_555","ehn_666"]}`,
			expectedResponse: &PatchPropertyVersionHostnamesResponse{
				AccountID:       "act_789",
				ContractID:      "ctr_456",
				GroupID:         "grp_321",
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				PropertyName:    "test-property",
				Etag:            "123abc456def7890",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{},
				},
			},
		},
		"200 OK - add and remove": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Remove: []string{
						"ehn_555",
						"ehn_666",
					},
					Add: []HostnameAdd{
						{
							CnameType: HostnameCnameTypeEdgeHostname,
							CnameFrom: "m.example.com",
							CnameTo:   "example.com.edgekey.net",
						},
						{
							CnameType:            HostnameCnameTypeEdgeHostname,
							EdgeHostnameID:       "ehn_888",
							CnameFrom:            "example3.com",
							CertProvisioningType: CertTypeDefault,
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId": "act_789",
				"contractId": "ctr_456",
				"groupId": "grp_321",
				"propertyId": "prp_123",
				"propertyName": "test-property",
				"propertyVersion": 1,
				"etag": "123abc456def7890",
				"hostnames": {
					"items": [
						{
							"cnameType": "EDGE_HOSTNAME",
							"edgeHostnameId": "ehn_777",
							"cnameFrom": "m.example.com",
							"cnameTo": "example.com.edgekey.net",
							"certProvisioningType": "CPS_MANAGED"
						},
						{
							"cnameType": "EDGE_HOSTNAME",
							"edgeHostnameId": "ehn_888",
							"cnameFrom": "example3.com",
							"cnameTo": "m.example.com.edgesuite.net",
							"certProvisioningType": "DEFAULT"
						}
					]
				}
			}`,
			expectedPath:        "/papi/v1/properties/prp_123/versions/1/hostnames",
			expectedRequestBody: `{"add":[{"cnameFrom":"m.example.com","cnameType":"EDGE_HOSTNAME","cnameTo":"example.com.edgekey.net"},{"cnameFrom":"example3.com","cnameType":"EDGE_HOSTNAME","certProvisioningType":"DEFAULT","edgeHostnameId":"ehn_888"}],"remove":["ehn_555","ehn_666"]}`,
			expectedResponse: &PatchPropertyVersionHostnamesResponse{
				AccountID:       "act_789",
				ContractID:      "ctr_456",
				GroupID:         "grp_321",
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				PropertyName:    "test-property",
				Etag:            "123abc456def7890",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "m.example.com",
							CnameTo:              "example.com.edgekey.net",
							EdgeHostnameID:       "ehn_777",
							CertProvisioningType: "CPS_MANAGED",
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_888",
							CnameFrom:            "example3.com",
							CertProvisioningType: "DEFAULT",
							CnameTo:              "m.example.com.edgesuite.net",
						},
					},
				},
			},
		},
		"200 OK - with optional fields and validation in progress": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:        "prp_123",
				PropertyVersion:   1,
				ContractID:        "ctr_456",
				GroupID:           "grp_321",
				IncludeCertStatus: true,
				ValidateHostnames: true,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CnameType: HostnameCnameTypeEdgeHostname,
							CnameFrom: "m.example.com",
							CnameTo:   "example.com.edgekey.net",
						},
						{
							CnameType:            HostnameCnameTypeEdgeHostname,
							EdgeHostnameID:       "ehn_666",
							CnameFrom:            "example3.com",
							CertProvisioningType: CertTypeDefault,
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId": "act_789",
				"contractId": "ctr_456",
				"groupId": "grp_321",
				"propertyId": "prp_123",
				"propertyName": "test-property",
				"propertyVersion": 1,
				"etag": "123abc456def7890",
				"hostnames": {
					"items": [
						{
							"cnameType": "EDGE_HOSTNAME",
							"edgeHostnameId": "ehn_555",
							"cnameFrom": "m.example.com",
							"cnameTo": "example.com.edgekey.net",
							"certProvisioningType": "CPS_MANAGED"
						},
						{
							"cnameType": "EDGE_HOSTNAME",
							"edgeHostnameId": "ehn_666",
							"cnameFrom": "example3.com",
							"cnameTo": "m.example.com.edgesuite.net",
							"certProvisioningType": "DEFAULT",
							"certStatus": {
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
									"hostname": "_acme-challenge.example3.com",
									"target": "ac.ee03f141752f0e52808f6669ad50ad43.example3.com.validate-akdv.net"
								}
							},
							"domainOwnershipVerification": {
								"challengeTokenExpiryDate": "2024-05-14T05:25:56Z",
								"status": "VALIDATION_IN_PROGRESS",
								"validationCname": {
									"hostname": "validation.hostname.example.com",
									"target": "validation.target.example.com"
								},
								"validationHttp": {
									"redirectMethod": {
										"httpRedirectFrom": "http.validation.redirect.from.example.com",
										"httpRedirectTo": "http.validation.redirect.to.example.com"
									},
									"fileContentMethod": {
										"url": "http://validation.file.content.example.com/validation.txt",
										"body": "HTTP validation body"
									}
								},
								"validationTxt": {
									"hostname": "txt.validation.hostname.example.com",
									"challengeToken": "token"
								}
							}
						}
					]
				}
			}`,
			expectedPath:        "/papi/v1/properties/prp_123/versions/1/hostnames?contractId=ctr_456&groupId=grp_321&includeCertStatus=true&validateHostnames=true",
			expectedRequestBody: `{"add":[{"cnameFrom":"m.example.com","cnameType":"EDGE_HOSTNAME","cnameTo":"example.com.edgekey.net"},{"cnameFrom":"example3.com","cnameType":"EDGE_HOSTNAME","certProvisioningType":"DEFAULT","edgeHostnameId":"ehn_666"}]}`,
			expectedResponse: &PatchPropertyVersionHostnamesResponse{
				AccountID:       "act_789",
				ContractID:      "ctr_456",
				GroupID:         "grp_321",
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				PropertyName:    "test-property",
				Etag:            "123abc456def7890",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "m.example.com",
							CnameTo:              "example.com.edgekey.net",
							EdgeHostnameID:       "ehn_555",
							CertProvisioningType: "CPS_MANAGED",
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_666",
							CnameFrom:            "example3.com",
							CertProvisioningType: "DEFAULT",
							CnameTo:              "m.example.com.edgesuite.net",
							CertStatus: CertStatusItem{
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
									Hostname: "_acme-challenge.example3.com",
									Target:   "ac.ee03f141752f0e52808f6669ad50ad43.example3.com.validate-akdv.net",
								},
							},
							DomainOwnershipVerification: &DomainOwnershipVerification{
								ChallengeTokenExpiryDate: ptr.To(time.Date(2024, 5, 14, 5, 25, 56, 0, time.UTC)),
								Status:                   "VALIDATION_IN_PROGRESS",
								ValidationCname: &ValidationCname{
									Hostname: "validation.hostname.example.com",
									Target:   "validation.target.example.com",
								},
								ValidationHTTP: &ValidationHTTP{
									RedirectMethod: RedirectMethod{
										HTTPRedirectFrom: "http.validation.redirect.from.example.com",
										HTTPRedirectTo:   "http.validation.redirect.to.example.com",
									},
									FileContentMethod: FileContentMethod{
										URL:  "http://validation.file.content.example.com/validation.txt",
										Body: "HTTP validation body",
									},
								},
								ValidationTXT: &ValidationTXT{
									Hostname:       "txt.validation.hostname.example.com",
									ChallengeToken: "token",
								},
							},
						},
					},
				},
			},
		},
		"200 OK - with authorization object returned": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:        "prp_123",
				PropertyVersion:   1,
				ContractID:        "ctr_456",
				GroupID:           "grp_321",
				IncludeCertStatus: true,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CnameType: HostnameCnameTypeEdgeHostname,
							CnameFrom: "m.example.com",
							CnameTo:   "example.com.edgekey.net",
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId": "act_789",
				"contractId": "ctr_456",
				"groupId": "grp_321",
				"propertyId": "prp_123",
				"propertyName": "test-property",
				"propertyVersion": 1,
				"etag": "123abc456def7890",
				"hostnames": {
					"items": [
						{
							"cnameType": "EDGE_HOSTNAME",
							"edgeHostnameId": "ehn_555",
							"cnameFrom": "m.example.com",
							"cnameTo": "example.com.edgekey.net",
							"certProvisioningType": "CPS_MANAGED",
							"certStatus": {
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
									"hostname": "_acme-challenge.m.example.com",
									"target": "ac.secret.m.example.com.validate-akdv.net"
								},
								"authorization": {
									"certDomain": "m.example.com",
									"status": "ATTEMPTING_VALIDATION",
									"validUntil": "2026-01-22T03:21:43Z",
									"http01": {
										"url": "http://m.example.com/.well-known/acme-challenge/anothersecret",
										"body": "anothersecret.secretsecret",
										"result": {
											"message": "Got NXDomain error while getting A records for m.example.com/.well-known/acme-challenge/anothersecret",
											"timestamp": "2026-01-21T13:42:12Z",
											"source": "DNS"
										}
									},
									"dns01": {
										"value": "secret2",
										"result": {
											"message": "Got NXDomain error while getting TXT for _acme-challenge.m.example.com",
											"timestamp": "2026-01-21T13:42:12Z",
											"source": "DNS"
										}
									}
								}
							}
						}
					]
				}
			}`,
			expectedPath:        "/papi/v1/properties/prp_123/versions/1/hostnames?contractId=ctr_456&groupId=grp_321&includeCertStatus=true",
			expectedRequestBody: `{"add":[{"cnameFrom":"m.example.com","cnameType":"EDGE_HOSTNAME","cnameTo":"example.com.edgekey.net"}]}`,
			expectedResponse: &PatchPropertyVersionHostnamesResponse{
				AccountID:       "act_789",
				ContractID:      "ctr_456",
				GroupID:         "grp_321",
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				PropertyName:    "test-property",
				Etag:            "123abc456def7890",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "m.example.com",
							CnameTo:              "example.com.edgekey.net",
							EdgeHostnameID:       "ehn_555",
							CertProvisioningType: "CPS_MANAGED",
							CertStatus: CertStatusItem{
								Production: []StatusItem{{Status: "PENDING"}},
								Staging:    []StatusItem{{Status: "PENDING"}},
								ValidationCname: ValidationCname{
									Hostname: "_acme-challenge.m.example.com",
									Target:   "ac.secret.m.example.com.validate-akdv.net",
								},
								Authorization: &Authorization{
									DNS01: &DNSAuthorization{
										Result: AuthorizationResult{
											Message:   "Got NXDomain error while getting TXT for _acme-challenge.m.example.com",
											Timestamp: test.NewTimeFromString(t, "2026-01-21T13:42:12Z"),
											Source:    "DNS",
										},
										Value: "secret2",
									},
									HTTP01: &HTTPAuthorization{
										Body: "anothersecret.secretsecret",
										Result: AuthorizationResult{
											Message:   "Got NXDomain error while getting A records for m.example.com/.well-known/acme-challenge/anothersecret",
											Timestamp: test.NewTimeFromString(t, "2026-01-21T13:42:12Z"),
											Source:    "DNS",
										},
										URL: "http://m.example.com/.well-known/acme-challenge/anothersecret",
									},
									Status:     "ATTEMPTING_VALIDATION",
									ValidUntil: ptr.To(test.NewTimeFromString(t, "2026-01-22T03:21:43Z")),
								},
							},
						},
					},
				},
			},
		},
		"200 OK validated": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:        "prp_123",
				PropertyVersion:   1,
				ContractID:        "ctr_456",
				GroupID:           "grp_321",
				IncludeCertStatus: true,
				ValidateHostnames: true,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CnameType:            HostnameCnameTypeEdgeHostname,
							EdgeHostnameID:       "ehn_666",
							CnameFrom:            "example3.com",
							CertProvisioningType: CertTypeDefault,
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId": "act_789",
				"contractId": "ctr_456",
				"groupId": "grp_321",
				"propertyId": "prp_123",
				"propertyName": "test-property",
				"propertyVersion": 1,
				"etag": "123abc456def7890",
				"hostnames": {
					"items": [
						{
							"cnameType": "EDGE_HOSTNAME",
							"edgeHostnameId": "ehn_666",
							"cnameFrom": "example3.com",
							"cnameTo": "m.example.com.edgesuite.net",
							"certProvisioningType": "DEFAULT",
							"certStatus": {
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
									"hostname": "_acme-challenge.example3.com",
									"target": "ac.ee03f141752f0e52808f6669ad50ad43.example3.com.validate-akdv.net"
								}
							},
							"domainOwnershipVerification": {
								"status": "VALIDATED"
							}
						}
					]
				}
			}`,
			expectedPath:        "/papi/v1/properties/prp_123/versions/1/hostnames?contractId=ctr_456&groupId=grp_321&includeCertStatus=true&validateHostnames=true",
			expectedRequestBody: `{"add":[{"cnameFrom":"example3.com","cnameType":"EDGE_HOSTNAME","certProvisioningType":"DEFAULT","edgeHostnameId":"ehn_666"}]}`,
			expectedResponse: &PatchPropertyVersionHostnamesResponse{
				AccountID:       "act_789",
				ContractID:      "ctr_456",
				GroupID:         "grp_321",
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				PropertyName:    "test-property",
				Etag:            "123abc456def7890",
				Hostnames: HostnameResponseItems{
					Items: []HostnameResponseItem{
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_666",
							CnameFrom:            "example3.com",
							CertProvisioningType: "DEFAULT",
							CnameTo:              "m.example.com.edgesuite.net",
							CertStatus: CertStatusItem{
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
									Hostname: "_acme-challenge.example3.com",
									Target:   "ac.ee03f141752f0e52808f6669ad50ad43.example3.com.validate-akdv.net",
								},
							},
							DomainOwnershipVerification: &DomainOwnershipVerification{
								Status: "VALIDATED",
							},
						},
					},
				},
			},
		},
		"validation errors missing required fields": {
			params: PatchPropertyVersionHostnamesRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patching hostnames: struct validation: PropertyID: cannot be blank\nPropertyVersion: cannot be blank", err.Error())
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.ErrorIs(t, err, ErrPatchPropertyVersionHostnames)
			},
		},
		"validation errors missing required fields in body": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{{}},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patching hostnames: struct validation: Body: {\n\tAdd[0]: {\n\t\tCnameFrom: cannot be blank\n\t\trequired parameters: either CnameTo or EdgeHostnameID must be provided\n\t}\n}", err.Error())
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.ErrorIs(t, err, ErrPatchPropertyVersionHostnames)
			},
		},
		"validation errors- invalid CertProvisioningType and CnameType": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CnameType:            HostnameCnameType("WRONG"),
							CnameFrom:            "m.example.com",
							CnameTo:              "example.com.edgekey.net",
							CertProvisioningType: CertType("WRONG"),
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patching hostnames: struct validation: Body: {\n\tAdd[0]: {\n\t\tCertProvisioningType: value 'WRONG' is invalid. Must be one of: 'CPS_MANAGED', 'DEFAULT' or 'CCM'\n\t\tCnameType: value 'WRONG' is invalid. There is only one supported value of: EDGE_HOSTNAME\n\t}\n}",
					err.Error())
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.ErrorIs(t, err, ErrPatchPropertyVersionHostnames)
			},
		},
		"validation error - CCMCertificates RSACertID, ECDSACertID and MTLS CASetID are not a digit": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123456",
				PropertyVersion: 3,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: "CCM",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							CCMCertificates: &CCMCertificates{
								RSACertID:   "12345a",
								ECDSACertID: "98765a",
							},
							MTLS: &MTLS{
								CASetID:         "524125a",
								CheckClientOCSP: false,
								SendCASetClient: false,
							},
							TLSConfiguration: &TLSConfiguration{
								CipherProfile:            "ak-akamai-2020q1",
								DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
								StapleServerOcspResponse: true,
								FIPSMode:                 false,
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "{\n\t\tCCMCertificates: {\n\t\t\tECDSACertID: must contain digits only\n\t\t\tRSACertID: must contain digits only\n\t\t}\n\t\tMTLS: {\n\t\t\tCASetID: must contain digits only")
			},
		},
		"validation error - one of RSACertID and ECDSACertID must be provided for CCM": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123456",
				PropertyVersion: 3,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: "CCM",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							CCMCertificates:      &CCMCertificates{},
							MTLS: &MTLS{
								CASetID:         "524125",
								CheckClientOCSP: false,
								SendCASetClient: false,
							},
							TLSConfiguration: &TLSConfiguration{
								CipherProfile:            "ak-akamai-2020q1",
								DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
								StapleServerOcspResponse: true,
								FIPSMode:                 false,
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "either RSACertID or ECDSACertID must be provided")
			},
		},
		"validation error - CCMCertificates is required for CCM": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123456",
				PropertyVersion: 3,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: "CCM",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							MTLS: &MTLS{
								CASetID:         "524125",
								CheckClientOCSP: false,
								SendCASetClient: false,
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "struct validation: Body: {\n\tAdd[0]: {\n\t\tValidateCCMHostname: when using `certProvisioningType` set to `CCM`, the request body must contain `ccmCertificates` with at least `rsaCertId` or `ecdsaCertId`\n\t}\n}")
			},
		},
		"validation error - MTLS is only valid for CCM": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123456",
				PropertyVersion: 3,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: "CPS_MANAGED",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							MTLS: &MTLS{
								CASetID:         "524125",
								CheckClientOCSP: false,
								SendCASetClient: false,
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "the mTLS configuration is provided without `certProvisioningType` set to `CCM`")
			},
		},
		"validation error - TLSConfiguration is only valid for CCM": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: "CPS_MANAGED",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							TLSConfiguration: &TLSConfiguration{
								CipherProfile:            "ak-akamai-2020q1",
								DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
								StapleServerOcspResponse: true,
								FIPSMode:                 false,
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "the TLS configuration is provided without `certProvisioningType` set to `CCM`")
			},
		},
		"validation error - CCMCertificates is only valid for CCM": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: "CPS_MANAGED",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							CCMCertificates: &CCMCertificates{
								RSACertID:   "12345",
								ECDSACertID: "98765",
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "the CCM cert details are provided without `certProvisioningType` set to `CCM`")
			},
		},
		"validation error - missed CASetID when MTLS provided": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123456",
				PropertyVersion: 3,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: "CCM",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							CCMCertificates: &CCMCertificates{
								RSACertID: "12345",
							},
							MTLS: &MTLS{},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "CASetID: cannot be blank")
			},
		},
		"validation error - empty cipher profile": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:        "prp_123456",
				PropertyVersion:   3,
				GroupID:           "grp_54321",
				ContractID:        "ctr_1-2ABCD3",
				IncludeCertStatus: true,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CertProvisioningType: "CCM",
							CnameFrom:            "www.example-ccm.com",
							CnameTo:              "example.com.edgesuite.net",
							CnameType:            "EDGE_HOSTNAME",
							CCMCertificates: &CCMCertificates{
								RSACertID:   "12345",
								ECDSACertID: "98765",
							},
							MTLS: &MTLS{
								CASetID:         "524125",
								CheckClientOCSP: false,
								SendCASetClient: false,
							},
							TLSConfiguration: &TLSConfiguration{
								CipherProfile:            "",
								DisallowedTLSVersions:    []string{"TLSv1_1", "TLSv1"},
								StapleServerOcspResponse: true,
								FIPSMode:                 false,
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "{\n\t\tValidateCCMHostname: the cipher profile is empty in the TLS configuration\n\t}")
			},
		},
		"500 internal server status error": {
			params: PatchPropertyVersionHostnamesRequest{
				PropertyID:      "prp_123",
				PropertyVersion: 1,
				Body: PatchPropertyVersionHostnamesRequestBody{
					Add: []HostnameAdd{
						{
							CnameType: HostnameCnameTypeEdgeHostname,
							CnameFrom: "m.example.com",
							CnameTo:   "example.com.edgekey.net",
						},
						{
							CnameType:            HostnameCnameTypeEdgeHostname,
							EdgeHostnameID:       "ehn_666",
							CnameFrom:            "example3.com",
							CertProvisioningType: CertTypeDefault,
						},
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/papi/v1/properties/prp_123/versions/1/hostnames",
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error updating hostnames",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error updating hostnames",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrPatchPropertyVersionHostnames)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPatch, r.Method)

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}

				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.PatchPropertyVersionHostnames(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapiGetAuditHistory(t *testing.T) {
	tests := map[string]struct {
		params           GetAuditHistoryRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAuditHistoryResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - one entry in audit history": {
			params: GetAuditHistoryRequest{
				Hostname: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"hostname": "example.com",
				"history": {
					"items": [
						{
							"action": "ADD",
							"certProvisioningType": "CPS_MANAGED",
							"cnameTo": "example.com.edgekey.net",
							"contractId": "C-0N7RAC7",
							"edgeHostnameId": "ehn_123",
							"groupId": "12345",
							"network": "PRODUCTION",
							"propertyId": "prp_123",
							"timestamp": "2023-10-26T12:00:00Z",
							"user": "user_123"
						}
					]
				}
			}`,
			expectedPath: "/papi/v1/hostnames/example.com/audit-history",
			expectedResponse: &GetAuditHistoryResponse{
				Hostname: "example.com",
				History: HostnameHistory{
					Items: []HostnameHistoryItem{
						{
							Action:               "ADD",
							CertProvisioningType: "CPS_MANAGED",
							CnameTo:              "example.com.edgekey.net",
							ContractID:           "C-0N7RAC7",
							EdgeHostnameID:       "ehn_123",
							GroupID:              "12345",
							Network:              "PRODUCTION",
							PropertyID:           "prp_123",
							Timestamp:            "2023-10-26T12:00:00Z",
							User:                 "user_123",
						},
					},
				},
			},
		},
		"200 OK - multiple entries in audit history": {
			params: GetAuditHistoryRequest{
				Hostname: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"hostname": "example.com",
				"history": {
					"items": [
						{
							"action": "ADD",
							"certProvisioningType": "CPS_MANAGED",
							"cnameTo": "example.com.edgekey.net",
							"contractId": "C-0N7RAC7",
							"edgeHostnameId": "ehn_123",
							"groupId": "12345",
							"network": "PRODUCTION",
							"propertyId": "prp_123",
							"timestamp": "2023-10-26T12:00:00Z",
							"user": "user_123"
						},
						{
							"action": "MODIFY",
							"certProvisioningType": "CPS_MANAGED",
							"cnameTo": "example.com.edgekey.net",
							"contractId": "C-0N7RAC7",
							"edgeHostnameId": "ehn_123",
							"groupId": "12345",
							"network": "STAGING",
							"propertyId": "prp_123",
							"timestamp": "2023-10-27T14:30:00Z",
							"user": "user_123"
						},
						{
							"action": "ACTIVATE",
							"certProvisioningType": "CPS_MANAGED",
							"cnameTo": "example.com.edgekey.net",
							"contractId": "C-0N7RAC7",
							"edgeHostnameId": "ehn_123",
							"groupId": "12345",
							"network": "PRODUCTION",
							"propertyId": "prp_123",
							"timestamp": "2023-10-28T16:45:00Z",
							"user": "user_123"
						}
					]
				}
			}`,
			expectedPath: "/papi/v1/hostnames/example.com/audit-history",
			expectedResponse: &GetAuditHistoryResponse{
				Hostname: "example.com",
				History: HostnameHistory{
					Items: []HostnameHistoryItem{
						{
							Action:               "ADD",
							CertProvisioningType: "CPS_MANAGED",
							CnameTo:              "example.com.edgekey.net",
							ContractID:           "C-0N7RAC7",
							EdgeHostnameID:       "ehn_123",
							GroupID:              "12345",
							Network:              "PRODUCTION",
							PropertyID:           "prp_123",
							Timestamp:            "2023-10-26T12:00:00Z",
							User:                 "user_123",
						},
						{
							Action:               "MODIFY",
							CertProvisioningType: "CPS_MANAGED",
							CnameTo:              "example.com.edgekey.net",
							ContractID:           "C-0N7RAC7",
							EdgeHostnameID:       "ehn_123",
							GroupID:              "12345",
							Network:              "STAGING",
							PropertyID:           "prp_123",
							Timestamp:            "2023-10-27T14:30:00Z",
							User:                 "user_123",
						},
						{
							Action:               "ACTIVATE",
							CertProvisioningType: "CPS_MANAGED",
							CnameTo:              "example.com.edgekey.net",
							ContractID:           "C-0N7RAC7",
							EdgeHostnameID:       "ehn_123",
							GroupID:              "12345",
							Network:              "PRODUCTION",
							PropertyID:           "prp_123",
							Timestamp:            "2023-10-28T16:45:00Z",
							User:                 "user_123",
						},
					},
				},
			},
		},
		"200 OK - empty audit history": {
			params: GetAuditHistoryRequest{
				Hostname: "example.com",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"hostname": "example.com",
				"history": {
					"items": []
				}
			}`,
			expectedPath: "/papi/v1/hostnames/example.com/audit-history",
			expectedResponse: &GetAuditHistoryResponse{
				Hostname: "example.com",
				History: HostnameHistory{
					Items: []HostnameHistoryItem{},
				},
			},
		},
		"500 Internal Server Error": {
			params: GetAuditHistoryRequest{
				Hostname: "example.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching audit history",
				"status": 500
			}`,
			expectedPath: "/papi/v1/hostnames/example.com/audit-history",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching audit history",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrGetAuditHistory)
			},
		},
		"validation error - missing hostname": {
			params: GetAuditHistoryRequest{},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.ErrorIs(t, err, ErrGetAuditHistory)
				assert.Contains(t, err.Error(), "Hostname: cannot be blank")
			},
		},
		"validation error - empty hostname string": {
			params: GetAuditHistoryRequest{
				Hostname: "",
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.ErrorIs(t, err, ErrGetAuditHistory)
				assert.Contains(t, err.Error(), "Hostname: cannot be blank")
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
			result, err := client.GetAuditHistory(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
