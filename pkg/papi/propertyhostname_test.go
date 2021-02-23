package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestPapi_GetPropertyVersionHostnames(t *testing.T) {
	tests := map[string]struct {
		params           GetPropertyVersionHostnamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPropertyVersionHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetPropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupID:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
				IncludeCertStatus:false,
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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&includeCertStatus=false&validateHostnames=false",
			expectedResponse: &GetPropertyVersionHostnamesResponse{
				AccountID:       "act_1-1TJZFB",
				ContractID:      "ctr_1-1TJZH5",
				GroupID:         "grp_15225",
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Etag:            "6aed418629b4e5c0",
				Hostnames: HostnameResponseItems{
					Items: []Hostname{
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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=&groupId=&includeCertStatus=false&validateHostnames=false",
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

func TestPapi_UpdatePropertyVersionHostnames(t *testing.T) {
	tests := map[string]struct {
		params           UpdatePropertyVersionHostnamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdatePropertyVersionHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupID:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
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
 				"certProvisioningType": "CPS_MANAGED"
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
				Etag:            "6aed418629b4e5c0",
				Hostnames: HostnameResponseItems{
					Items: []Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895822",
							CnameFrom:            "m.example.com",
							CnameTo:              "example.com.edgekey.net",
							CertProvisioningType: "DEFAULT",
							CertStatus:CertStatusItem{
								ValidationCname: ValidationCname{
									Hostname: "_acme-challenge.www.example.com",
									Target:   "{token}.www.example.com.akamai-domain.com",
								},
								Staging: []StatusItem{{Status:"NEEDS_VALIDATION"},

								},
								Production: []StatusItem{{Status:"NEEDS_VALIDATION"},

								},
							},

					},
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_895833",
							CnameFrom:            "example3.com",
							CnameTo:              "m.example.com.edgesuite.net",
							CertProvisioningType: "CPS_MANAGED",
						},
					},
				},
			},
		},
		"200 empty hostnames": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupID:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
				IncludeCertStatus: true,
				Hostnames:       []Hostname{{}},
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
				Hostnames: HostnameResponseItems{
					Items: []Hostname{},
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
				IncludeCertStatus:true,
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
				Hostnames: HostnameResponseItems{
					Items: []Hostname{},
				},
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
		"200 Hostnames missing": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupID:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
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
				Hostnames: HostnameResponseItems{
					Items: []Hostname{},
				},
			},
		},
		"200 Hostnames items missing": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupID:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
				Hostnames:       nil,
				IncludeCertStatus:true,
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
				Hostnames: HostnameResponseItems{
					Items: []Hostname{},
				},
			},
		},
		"200 Hostnames items empty": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupID:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
				IncludeCertStatus: true,
				Hostnames:       []Hostname{},
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
				Hostnames: HostnameResponseItems{
					Items: []Hostname{},
				},
			},
		},
		"400 Hostnames cert type is invalid": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				GroupID:         "grp_15225",
				ContractID:      "ctr_1-1TJZH5",
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
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Hostnames:       []Hostname{{}},
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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=&groupId=&includeCertStatus=true&validateHostnames=false",
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
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
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
