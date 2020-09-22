package papi

import (
	"context"
	"errors"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi/tools"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPapi_GetEdgeHostnames(t *testing.T) {
	tests := map[string]struct {
		params           GetEdgeHostnamesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetEdgeHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetEdgeHostnamesRequest{
				ContractID: "contract",
				GroupID:    "group",
				Options:    []string{"opt1", "opt2"},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "acc",
    "contractId": "contract",
    "groupId": "group",
    "edgeHostnames": {
        "items": [
            {
                "edgeHostnameId": "ehID",
                "edgeHostnameDomain": "example.com.edgekey.net",
                "productId": "prdID",
                "domainPrefix": "example.com",
                "domainSuffix": "edgekey.net",
                "status": "PENDING",
                "secure": true,
                "ipVersionBehavior": "IPV4"
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1,opt2",
			expectedResponse: &GetEdgeHostnamesResponse{
				AccountID:  "acc",
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostnames: EdgeHostnameItems{Items: []EdgeHostnameGetItem{
					{
						ID:                "ehID",
						Domain:            "example.com.edgekey.net",
						ProductID:         "prdID",
						DomainPrefix:      "example.com",
						DomainSuffix:      "edgekey.net",
						Status:            "PENDING",
						Secure:            true,
						IPVersionBehavior: "IPV4",
						UseCases:          nil,
					},
				}},
			},
		},
		"500 internal server error": {
			params: GetEdgeHostnamesRequest{
				ContractID: "contract",
				GroupID:    "group",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching edge hostnames",
    "status": 500
}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := session.APIError{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching edge hostnames",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty group ID": {
			params: GetEdgeHostnamesRequest{
				ContractID: "contract",
				GroupID:    "",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "GroupID")
			},
		},
		"empty contract ID": {
			params: GetEdgeHostnamesRequest{
				ContractID: "",
				GroupID:    "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContractID")
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
			result, err := client.GetEdgeHostnames(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_GetEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		params           GetEdgeHostnameRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetEdgeHostnamesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetEdgeHostnameRequest{
				EdgeHostnameID: "ehID",
				ContractID:     "contract",
				GroupID:        "group",
				Options:        []string{"opt1", "opt2"},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "acc",
    "contractId": "contract",
    "groupId": "group",
    "edgeHostnames": {
        "items": [
            {
                "edgeHostnameId": "ehID",
                "edgeHostnameDomain": "example.com.edgekey.net",
                "productId": "prdID",
                "domainPrefix": "example.com",
                "domainSuffix": "edgekey.net",
                "status": "PENDING",
                "secure": true,
                "ipVersionBehavior": "IPV4"
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/edgehostnames/ehID?contractId=contract&groupId=group&options=opt1,opt2",
			expectedResponse: &GetEdgeHostnamesResponse{
				AccountID:  "acc",
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostnames: EdgeHostnameItems{Items: []EdgeHostnameGetItem{
					{
						ID:                "ehID",
						Domain:            "example.com.edgekey.net",
						ProductID:         "prdID",
						DomainPrefix:      "example.com",
						DomainSuffix:      "edgekey.net",
						Status:            "PENDING",
						Secure:            true,
						IPVersionBehavior: "IPV4",
						UseCases:          nil,
					},
				}},
			},
		},
		"500 internal server error": {
			params: GetEdgeHostnameRequest{
				EdgeHostnameID: "ehID",
				ContractID:     "contract",
				GroupID:        "group",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching edge hostnames",
    "status": 500
}`,
			expectedPath: "/papi/v1/edgehostnames/ehID?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := session.APIError{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching edge hostnames",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty group ID": {
			params: GetEdgeHostnameRequest{
				EdgeHostnameID: "ehID",
				ContractID:     "contract",
				GroupID:        "",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "GroupID")
			},
		},
		"empty contract ID": {
			params: GetEdgeHostnameRequest{
				EdgeHostnameID: "ehID",
				ContractID:     "",
				GroupID:        "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContractID")
			},
		},
		"empty edge hostname ID": {
			params: GetEdgeHostnameRequest{
				EdgeHostnameID: "",
				ContractID:     "contract",
				GroupID:        "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "EdgeHostnameID")
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
			result, err := client.GetEdgeHostname(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_CreateEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		params           CreateEdgeHostnameRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateEdgeHostnameResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				Options:    []string{"opt1", "opt2"},
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases: []UseCase{{
						Option:  "option",
						Type:    "GLOBAL",
						UseCase: "UseCase",
					}},
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "edgeHostnameLink": "/papi/v1/edgehostnames/ehID?contractId=contract&group=group"
}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1,opt2",
			expectedResponse: &CreateEdgeHostnameResponse{
				EdgeHostnameLink: "/papi/v1/edgehostnames/ehID?contractId=contract&group=group",
				EdgeHostnameID:   "ehID",
			},
		},
		"200 OK - STANDARD_TLS": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				Options:    []string{"opt1", "opt2"},
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgesuite.net",
					Secure:            true,
					SecureNetwork:     "STANDARD_TLS",
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "edgeHostnameLink": "/papi/v1/edgehostnames/ehID?contractId=contract&group=group"
}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1,opt2",
			expectedResponse: &CreateEdgeHostnameResponse{
				EdgeHostnameLink: "/papi/v1/edgehostnames/ehID?contractId=contract&group=group",
				EdgeHostnameID:   "ehID",
			},
		},
		"200 OK - SHARED_CERT": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				Options:    []string{"opt1", "opt2"},
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "akamaized.net",
					Secure:            true,
					SecureNetwork:     "SHARED_CERT",
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "edgeHostnameLink": "/papi/v1/edgehostnames/ehID?contractId=contract&group=group"
}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1,opt2",
			expectedResponse: &CreateEdgeHostnameResponse{
				EdgeHostnameLink: "/papi/v1/edgehostnames/ehID?contractId=contract&group=group",
				EdgeHostnameID:   "ehID",
			},
		},
		"200 OK - ENHANCED_TLS": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				Options:    []string{"opt1", "opt2"},
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					CertEnrollmentID:  5,
					Secure:            true,
					SecureNetwork:     "ENHANCED_TLS",
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "edgeHostnameLink": "/papi/v1/edgehostnames/ehID?contractId=contract&group=group"
}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1,opt2",
			expectedResponse: &CreateEdgeHostnameResponse{
				EdgeHostnameLink: "/papi/v1/edgehostnames/ehID?contractId=contract&group=group",
				EdgeHostnameID:   "ehID",
			},
		},
		"500 Internal Server Error": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				Options:    []string{"opt1", "opt2"},
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating edge hostname",
    "status": 500
}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1,opt2",
			withError: func(t *testing.T, err error) {
				want := session.APIError{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating edge hostname",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty group ID": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "GroupID")
			},
		},
		"empty contract ID": {
			params: CreateEdgeHostnameRequest{
				ContractID: "",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContractID")
			},
		},
		"empty domain prefix": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "DomainPrefix")
			},
		},
		"empty domain suffix": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "DomainSuffix")
			},
		},
		"empty product id": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ProductID")
			},
		},
		"CertEnrollmentID is required for SecureNetwork == ENHANCED_TLS": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					SecureNetwork:     "ENHANCED_TLS",
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "CertEnrollmentID")
			},
		},
		"SecureNetwork has invalid value": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					SecureNetwork:     "test",
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecureNetwork")
			},
		},
		"DomainSuffix has invalid value for SecureNetwork == STANDARD_TLS": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					SecureNetwork:     "STANDARD_TLS",
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "DomainSuffix")
			},
		},
		"DomainSuffix has invalid value for SecureNetwork == SHARED_CERT": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					SecureNetwork:     "SHARED_CERT",
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "DomainSuffix")
			},
		},
		"DomainSuffix has invalid value for SecureNetwork == ENHANCED_TLS": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "akamized.net",
					Secure:            true,
					SecureNetwork:     "STANDARD_TLS",
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "DomainSuffix")
			},
		},
		"IPVersionBehavior has invalid value": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "akamized.net",
					Secure:            true,
					IPVersionBehavior: "test",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "IPVersionBehavior")
			},
		},
		"UseCase has empty UseCase value": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "akamized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases: []UseCase{{
						Option:  "option",
						Type:    "GLOBAL",
						UseCase: "",
					}},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "UseCase")
			},
		},
		"UseCase has empty Option value": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "akamized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases: []UseCase{{
						Option:  "",
						Type:    "GLOBAL",
						UseCase: "useCase",
					}},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Option")
			},
		},
		"UseCase has invalid Type value": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "akamized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases: []UseCase{{
						Option:  "option",
						Type:    "test",
						UseCase: "UseCase",
					}},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Type")
			},
		},
		"invalid location": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				Options:    []string{"opt1", "opt2"},
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "example.com",
					DomainSuffix:      "edgekey.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases: []UseCase{{
						Option:  "option",
						Type:    "GLOBAL",
						UseCase: "UseCase",
					}},
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "edgeHostnameLink": ":"
}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1,opt2",
			withError: func(t *testing.T, err error) {
				want := tools.ErrInvalidLocation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateEdgeHostname(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
