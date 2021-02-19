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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=&groupId=&validateHostnames=false",
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
				Hostnames: []Hostname{
					{
						CnameType:            "EDGE_HOSTNAME",
						CnameFrom:            "m.example.com",
						CnameTo:              "example.com.edgesuite.net",
						CertProvisioningType: "CPS_MANAGED",
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
                "cnameFrom": "example.com",
                "cnameTo": "example.com.edgesuite.net"
                "certProvisioningType: "CPS_MANAGED"
            },
            {
                "cnameType": "EDGE_HOSTNAME",
                "edgeHostnameId": "ehn_895833",
                "cnameFrom": "m.example.com",
                "cnameTo": "m.example.com.edgesuite.net"
 				"certProvisioningType: "CPS_MANAGED"
            }
        ]
    }
}
`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&validateHostnames=false",
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
							CnameType:      "EDGE_HOSTNAME",
							EdgeHostnameID: "ehn_895822",
							CnameFrom:      "example.com",
							CnameTo:        "example.com.edgesuite.net",
							CertProvisioningType: "CPS_MANAGED",
						},
						{
							CnameType:      "EDGE_HOSTNAME",
							EdgeHostnameID: "ehn_895833",
							CnameFrom:      "m.example.com",
							CnameTo:        "m.example.com.edgesuite.net",
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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&validateHostnames=false",
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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&validateHostnames=true",
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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&validateHostnames=false",
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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&validateHostnames=false",
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
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=ctr_1-1TJZH5&groupId=grp_15225&validateHostnames=false",
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
		"500 internal server status error": {
			params: UpdatePropertyVersionHostnamesRequest{
				PropertyID:      "prp_175780",
				PropertyVersion: 3,
				Hostnames:       []Hostname{{}},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating hostnames",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/versions/3/hostnames?contractId=&groupId=&validateHostnames=false",
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
