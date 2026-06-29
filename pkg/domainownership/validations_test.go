package domainownership

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDomains(t *testing.T) {
	tests := map[string]struct {
		request             ValidateDomainsRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *ValidateDomainsResponse
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			request: ValidateDomainsRequest{
				Domains: []ValidateDomain{
					{
						DomainName:       "sample1.com",
						ValidationScope:  ValidationScopeHost,
						ValidationMethod: ValidationMethodHTTP,
					},
					{
						DomainName:       "sample2.com",
						ValidationScope:  ValidationScopeWildcard,
						ValidationMethod: ValidationMethodDNSCNAME,
					},
					{
						DomainName:       "sample3.com",
						ValidationScope:  ValidationScopeDomain,
						ValidationMethod: ValidationMethodDNSTXT,
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "domains": [
    {
      "domainName": "sample1.com",
      "domainStatus": "VALIDATED",
      "validationScope": "HOST"
    },
    {
      "domainName": "sample2.com",
      "domainStatus": "REQUEST_ACCEPTED",
      "validationScope": "WILDCARD"
    },
    {
      "domainName": "sample3.com",
      "domainStatus": "REQUEST_ACCEPTED",
      "validationScope": "DOMAIN"
    }
  ]
}
`,
			expectedPath:        "/domain-validation/v1/domains/validate-now",
			expectedRequestBody: `{"domains":[{"domainName":"sample1.com","validationMethod":"HTTP","validationScope":"HOST"},{"domainName":"sample2.com","validationMethod":"DNS_CNAME","validationScope":"WILDCARD"},{"domainName":"sample3.com","validationMethod":"DNS_TXT","validationScope":"DOMAIN"}]}`,
			expectedResponse: &ValidateDomainsResponse{
				Domains: []ValidateDomainResponse{
					{
						DomainName:      "sample1.com",
						DomainStatus:    "VALIDATED",
						ValidationScope: "HOST",
					},
					{
						DomainName:      "sample2.com",
						DomainStatus:    "REQUEST_ACCEPTED",
						ValidationScope: "WILDCARD",
					},
					{
						DomainName:      "sample3.com",
						DomainStatus:    "REQUEST_ACCEPTED",
						ValidationScope: "DOMAIN",
					},
				},
			},
		},
		"validation - empty domain": {
			request: ValidateDomainsRequest{
				Domains: []ValidateDomain{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "validate domains: struct validation:\nDomains: cannot be blank", err.Error())
			},
		},
		"validation - domain Name not supplied": {
			request: ValidateDomainsRequest{
				Domains: []ValidateDomain{
					{
						DomainName:       "sample1.com",
						ValidationScope:  ValidationScopeHost,
						ValidationMethod: ValidationMethodHTTP,
					},
					{
						ValidationScope:  ValidationScopeHost,
						ValidationMethod: ValidationMethodHTTP,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "validate domains: struct validation:\nDomains[1]: {\n\tDomainName: cannot be blank\n}", err.Error())
			},
		},
		"validation - validation scope not supplied": {
			request: ValidateDomainsRequest{
				Domains: []ValidateDomain{
					{
						DomainName:       "sample1.com",
						ValidationMethod: ValidationMethodHTTP,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "validate domains: struct validation:\nDomains[0]: {\n\tValidationScope: cannot be blank\n}", err.Error())
			},
		},
		"validation - validation method not supplied": {
			request: ValidateDomainsRequest{
				Domains: []ValidateDomain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "validate domains: struct validation:\nDomains[0]: {\n\tValidationMethod: cannot be blank\n}", err.Error())
			},
		},
		"validation - incorrect ValidationScope": {
			request: ValidateDomainsRequest{
				Domains: []ValidateDomain{
					{
						DomainName:       "sample1.com",
						ValidationScope:  ValidationScope("incorrect"),
						ValidationMethod: ValidationMethodHTTP,
					},
					{
						DomainName:       "sample2.com",
						ValidationScope:  ValidationScopeHost,
						ValidationMethod: ValidationMethodHTTP,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "validate domains: struct validation:\nDomains[0]: {\n\tValidationScope: value 'incorrect' is invalid. Must be one of: 'HOST', 'DOMAIN' or 'WILDCARD'\n}", err.Error())
			},
		},
		"validation - incorrect ValidationMethod": {
			request: ValidateDomainsRequest{
				Domains: []ValidateDomain{
					{
						DomainName:       "sample1.com",
						ValidationMethod: ValidationMethod("incorrect"),
						ValidationScope:  ValidationScopeHost,
					},
					{
						DomainName:       "sample2.com",
						ValidationScope:  ValidationScopeHost,
						ValidationMethod: ValidationMethodHTTP,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "validate domains: struct validation:\nDomains[0]: {\n\tValidationMethod: value must be one of: 'DNS_CNAME', 'DNS_TXT' or 'HTTP'\n}", err.Error())
			},
		},
		"500 internal server error": {
			request: ValidateDomainsRequest{
				Domains: []ValidateDomain{
					{
						DomainName:       "sample1.com",
						ValidationScope:  ValidationScopeDomain,
						ValidationMethod: ValidationMethodHTTP,
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/domain-validation/v1/domains/validate-now",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ValidateDomains(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestInvalidateDomain(t *testing.T) {
	tests := map[string]struct {
		params              InvalidateDomainRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *InvalidateDomainResponse
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			params: InvalidateDomainRequest{
				DomainName:      "sample1.com",
				ValidationScope: ValidationScopeHost,
			},
			responseStatus: http.StatusOK,
			responseBody: `
    {
      "domainName": "sample1.com",
      "domainStatus": "INVALIDATED",
      "validationScope": "HOST"
    }
`,
			expectedPath:        "/domain-validation/v1/domains/invalidate/sample1.com?validationScope=HOST",
			expectedRequestBody: "",
			expectedResponse: &InvalidateDomainResponse{
				DomainName:      "sample1.com",
				ValidationScope: "HOST",
				DomainStatus:    "INVALIDATED",
			},
		},
		"validation - domain Name not supplied": {
			params: InvalidateDomainRequest{
				ValidationScope: ValidationScopeHost,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "invalidate domain: struct validation:\nDomainName: cannot be blank", err.Error())
			},
		},
		"validation - ValidationScope not supplied": {
			params: InvalidateDomainRequest{
				DomainName: "sample1.com",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "invalidate domain: struct validation:\nValidationScope: cannot be blank", err.Error())
			},
		},
		"validation - incorrect ValidationScope": {
			params: InvalidateDomainRequest{
				DomainName:      "sample1.com",
				ValidationScope: ValidationScope("incorrect"),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "invalidate domain: struct validation:\nValidationScope: value 'incorrect' is invalid. Must be one of: 'HOST', 'DOMAIN' or 'WILDCARD'", err.Error())
			},
		},
		"500 internal server error": {
			params: InvalidateDomainRequest{

				DomainName:      "sample1.com",
				ValidationScope: ValidationScopeHost,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/domain-validation/v1/domains/invalidate/sample1.com?validationScope=HOST",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.InvalidateDomain(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestInvalidateDomains(t *testing.T) {
	tests := map[string]struct {
		request             InvalidateDomainsRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *InvalidateDomainsResponse
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			request: InvalidateDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
					{
						DomainName:      "sample2.com",
						ValidationScope: ValidationScopeWildcard,
					},
					{
						DomainName:      "sample3.com",
						ValidationScope: ValidationScopeDomain,
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "domains": [
    {
      "domainName": "sample1.com",
      "domainStatus": "INVALIDATED",
      "validationScope": "HOST"
    },
    {
      "domainName": "sample2.com",
      "domainStatus": "INVALIDATED",
      "validationScope": "WILDCARD"
    },
    {
      "domainName": "sample3.com",
      "domainStatus": "INVALIDATED",
      "validationScope": "DOMAIN"
    }
  ]
}
`,
			expectedPath:        "/domain-validation/v1/domains/invalidate",
			expectedRequestBody: `{"domains":[{"domainName":"sample1.com","validationScope":"HOST"},{"domainName":"sample2.com","validationScope":"WILDCARD"},{"domainName":"sample3.com","validationScope":"DOMAIN"}]}`,
			expectedResponse: &InvalidateDomainsResponse{
				Domains: []InvalidateDomainResponse{
					{
						DomainName:      "sample1.com",
						DomainStatus:    "INVALIDATED",
						ValidationScope: "HOST",
					},
					{
						DomainName:      "sample2.com",
						DomainStatus:    "INVALIDATED",
						ValidationScope: "WILDCARD",
					},
					{
						DomainName:      "sample3.com",
						DomainStatus:    "INVALIDATED",
						ValidationScope: "DOMAIN",
					},
				},
			},
		},
		"validation - empty domain": {
			request: InvalidateDomainsRequest{
				Domains: []Domain{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "invalidate domains: struct validation:\nDomains: cannot be blank", err.Error())
			},
		},
		"validation - domain Name not supplied": {
			request: InvalidateDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeHost,
					},
					{
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "invalidate domains: struct validation:\nDomains[1]: {\n\tDomainName: cannot be blank\n}", err.Error())
			},
		},
		"validation - validation scope not supplied": {
			request: InvalidateDomainsRequest{
				Domains: []Domain{
					{
						DomainName: "sample1.com",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "invalidate domains: struct validation:\nDomains[0]: {\n\tValidationScope: cannot be blank\n}", err.Error())
			},
		},
		"validation - incorrect ValidationScope": {
			request: InvalidateDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScope("incorrect"),
					},
					{
						DomainName:      "sample2.com",
						ValidationScope: ValidationScopeHost,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "invalidate domains: struct validation:\nDomains[0]: {\n\tValidationScope: value 'incorrect' is invalid. Must be one of: 'HOST', 'DOMAIN' or 'WILDCARD'\n}", err.Error())
			},
		},
		"500 internal server error": {
			request: InvalidateDomainsRequest{
				Domains: []Domain{
					{
						DomainName:      "sample1.com",
						ValidationScope: ValidationScopeDomain,
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/domain-validation/v1/domains/invalidate",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error making request",
					Status: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.InvalidateDomains(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
