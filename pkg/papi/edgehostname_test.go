package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiGetEdgeHostnames(t *testing.T) {
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
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1%2Copt2",
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
				want := &Error{
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

func TestPapiGetEdgeHostname(t *testing.T) {
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
			expectedPath: "/papi/v1/edgehostnames/ehID?contractId=contract&groupId=group&options=opt1%2Copt2",
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
				EdgeHostname: EdgeHostnameGetItem{
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
			},
		},
		"Edge hostname not found": {
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
        ]
    }
}`,
			expectedPath: "/papi/v1/edgehostnames/ehID?contractId=contract&groupId=group&options=opt1%2Copt2",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrNotFound), "want: %v; got: %v", ErrNotFound, err)
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
				want := &Error{
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

func TestPapiCreateEdgeHostname(t *testing.T) {
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
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1%2Copt2",
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
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1%2Copt2",
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
					DomainPrefix:      "example-com",
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
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1%2Copt2",
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
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1%2Copt2",
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
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1%2Copt2",
			withError: func(t *testing.T, err error) {
				want := &Error{
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
		"invalid edge hostname domain prefix for the akamaized.net domain suffix - The character '#' isn't allowed in the domain prefix.": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "tes#t",
					DomainSuffix:      "akamaized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := "A prefix for the edge hostname with the \"akamaized.net\" suffix must begin with a letter, end with a letter or digit, and contain only letters, digits, and hyphens, for example, abc-def, or abc-123"
				assert.True(t, err != nil && strings.Contains(err.Error(), want), "Expected error containing %q, got %v", want, err)
			},
		},
		"invalid edge hostname domain prefix for the `akamaized.net` domain suffix. The domain prefix contains non-UTF-8 characters ('t中esãt').": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "t中esãt",
					DomainSuffix:      "akamaized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := "A prefix for the edge hostname with the \"akamaized.net\" suffix must begin with a letter, end with a letter or digit, and contain only letters, digits, and hyphens, for example, abc-def, or abc-123"
				assert.True(t, err != nil && strings.Contains(err.Error(), want), "Expected error containing %q, got %v", want, err)
			},
		},
		"invalid edge hostname domain prefix for `akamaized.net`. The domain prefix can't end with a hyphen": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "test-",
					DomainSuffix:      "akamaized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := "A prefix for the edge hostname with the \"akamaized.net\" suffix must begin with a letter, end with a letter or digit, and contain only letters, digits, and hyphens, for example, abc-def, or abc-123"
				assert.True(t, err != nil && strings.Contains(err.Error(), want), "expected error containing %q, got %v", want, err)
			},
		},
		"invalid edge hostname domain prefix for `edgesuite.net`. The domain prefix can't end with two consecutive dots": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "test..",
					DomainSuffix:      "edgesuite.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := "A prefix for the edge hostname with the \"edgesuite.net\" suffix must begin with a letter, end with a letter, digit or dot, and contain only letters, digits, dots, and hyphens, for example, abc-def.123.456., or abc.123-def"
				assert.True(t, err != nil && strings.Contains(err.Error(), want), "Expected error containing %q, got %v", want, err)
			},
		},
		"invalid edge hostname domain prefix. The domain prefix exceeds the maximum allowed length of 63 characters": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "testABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567",
					DomainSuffix:      "akamaized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := `The edge hostname prefix must be at least 4 character(s) and no more than 63 characters for "akamaized.net" suffix; you provided 64 character(s)`
				assert.True(t, err != nil && strings.Contains(err.Error(), want), "Expected error containing %q, got %v", want, err)
			},
		},
		"invalid edge hostname domain prefix. The domain prefix less the minimum required length of 4 characters": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "nm1",
					DomainSuffix:      "akamaized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := `The edge hostname prefix must be at least 4 character(s) and no more than 63 characters for "akamaized.net" suffix; you provided 3 character(s)`
				assert.True(t, err != nil && strings.Contains(err.Error(), want), "Expected error containing %q, got %v", want, err)
			},
		},
		"invalid edge hostname domain prefix. The domain prefix less the minimum required length of 1 character": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "",
					DomainSuffix:      ".edgesuite.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			withError: func(t *testing.T, err error) {
				want := `DomainPrefix: cannot be blank`
				assert.True(t, err != nil && strings.Contains(err.Error(), want), "Expected error containing %q, got %v", want, err)
			},
		},
		"valid edge hostname with hyphen in domain prefix name for akamaized.net, create edge hostname": {
			params: CreateEdgeHostnameRequest{
				ContractID: "contract",
				GroupID:    "group",
				EdgeHostname: EdgeHostnameCreate{
					ProductID:         "product",
					DomainPrefix:      "test-1",
					DomainSuffix:      "akamaized.net",
					Secure:            true,
					IPVersionBehavior: "IPV4",
					UseCases:          nil,
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
			{
    			"edgeHostnameLink": "/papi/v1/edgehostnames/ehID?contractId=contract&group=group"
			}`,
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group",
			expectedResponse: &CreateEdgeHostnameResponse{
				EdgeHostnameLink: "/papi/v1/edgehostnames/ehID?contractId=contract&group=group",
				EdgeHostnameID:   "ehID",
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
			expectedPath: "/papi/v1/edgehostnames?contractId=contract&groupId=group&options=opt1%2Copt2",
			withError: func(t *testing.T, err error) {
				want := ErrInvalidResponseLink
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
