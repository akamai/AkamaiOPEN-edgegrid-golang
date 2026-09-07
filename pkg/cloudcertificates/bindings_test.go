package cloudcertificates

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCertificateBindings(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params           ListCertificateBindingsRequest
		responseStatus   int
		responseBody     string
		expectedResponse *ListCertificateBindingsResponse
		expectedPath     string
		returnedHeaders  map[string]string
		withError        func(*testing.T, error)
	}{
		"200 - fetch of certificate bindings successful": {
			params: ListCertificateBindingsRequest{
				CertificateID: "123",
			},
			returnedHeaders: map[string]string{
				"Akamai-RateLimit-Limit":     "60",
				"Akamai-RateLimit-Remaining": "59",
			},
			expectedResponse: &ListCertificateBindingsResponse{
				Bindings: []CertificateBinding{
					{
						CertificateID: "123456",
						Hostname:      "www.example.com",
						Network:       "PRODUCTION",
						ResourceType:  "CDN_HOSTNAME",
					},
					{
						CertificateID: "654321",
						Hostname:      "secure.example.com",
						Network:       "STAGING",
						ResourceType:  "CDN_HOSTNAME",
					},
					{
						CertificateID: "789012",
						Hostname:      "api.example.com",
						Network:       "PRODUCTION",
						ResourceType:  "CDN_HOSTNAME",
					},
					// There should be 10 bindings in the response to match links correctly, but only 3 are shown here for brevity.
				},
				Links: Links{
					Next:     ptr.To("https://api.example.com/v1/certificates/123/certificate-bindings?page=2&pageSize=10"),
					Previous: nil,
					Self:     "https://api.example.com/v1/certificates/123/certificate-bindings?page=1&pageSize=10",
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			responseStatus: 200,
			expectedPath:   "/ccm/v1/certificates/123/certificate-bindings",
			responseBody: `
{
  "bindings": [
    {
      "certificateId": 123456,
      "hostname": "www.example.com",
      "network": "PRODUCTION",
      "resourceType": "CDN_HOSTNAME"
    },
    {
      "certificateId": "654321",
      "hostname": "secure.example.com",
      "network": "STAGING",
      "resourceType": "CDN_HOSTNAME"
    },
    {
      "certificateId": "789012",
      "hostname": "api.example.com",
      "network": "PRODUCTION",
      "resourceType": "CDN_HOSTNAME"
    }
  ],
  "links": {
    "next": "https://api.example.com/v1/certificates/123/certificate-bindings?page=2&pageSize=10",
    "previous": null,
    "self": "https://api.example.com/v1/certificates/123/certificate-bindings?page=1&pageSize=10"
  }
}`,
		},
		"200 - fetch of certificate bindings with paging successful": {
			params: ListCertificateBindingsRequest{
				CertificateID: "123",
				Page:          3,
				PageSize:      1,
			},
			expectedResponse: &ListCertificateBindingsResponse{
				Bindings: []CertificateBinding{
					{
						CertificateID: "789012",
						Hostname:      "api.example.com",
						Network:       "PRODUCTION",
						ResourceType:  "CDN_HOSTNAME",
					},
				},
				Links: Links{
					Next:     nil,
					Previous: ptr.To("https://api.example.com/v1/certificates/123/certificate-bindings?page=2&pageSize=1"),
					Self:     "https://api.example.com/v1/certificates/123/certificate-bindings?page=3&pageSize=1",
				},
			},
			responseStatus: 200,
			expectedPath:   "/ccm/v1/certificates/123/certificate-bindings?page=3&pageSize=1",
			responseBody: `
{
  "bindings": [
    {
      "certificateId": "789012",
      "hostname": "api.example.com",
      "network": "PRODUCTION",
      "resourceType": "CDN_HOSTNAME"
    }
  ],
  "links": {
    "next": null,
    "previous": "https://api.example.com/v1/certificates/123/certificate-bindings?page=2&pageSize=1",
    "self": "https://api.example.com/v1/certificates/123/certificate-bindings?page=3&pageSize=1"
  }
}`,
		},
		"404 resource not found - certificate not found": {
			params: ListCertificateBindingsRequest{
				CertificateID: "1234",
			},
			responseStatus: 404,
			expectedPath:   "/ccm/v1/certificates/1234/certificate-bindings",
			responseBody: `{
				"certificateIdentifier": "certificateId",
				"certificateIdentifierValue": "1234",
				"detail": "Certificate with {certificateId}: {1234} is not found.",
				"instance": "/error-types/certificate-not-found?traceId=-11111",
				"status": 404,
				"title": "Certificate is not found.",
				"type": "/error-types/certificate-not-found"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListCertificateBindings, &Error{
					Type:                       "/error-types/certificate-not-found",
					Title:                      "Certificate is not found.",
					Detail:                     "Certificate with {certificateId}: {1234} is not found.",
					Status:                     http.StatusNotFound,
					Instance:                   "/error-types/certificate-not-found?traceId=-11111",
					CertificateIdentifier:      "certificateId",
					CertificateIdentifierValue: "1234",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertificateNotFound)
				assert.ErrorIs(t, err, ErrListCertificateBindings)
			},
		},
		"500 internal server error - assert that error is ErrListCertificateBindings": {
			params: ListCertificateBindingsRequest{
				CertificateID: "123",
			},
			responseStatus: 500,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error removing certificate",
				"status": 500
			}`,
			expectedPath: "/ccm/v1/certificates/123/certificate-bindings",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListCertificateBindings, &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error removing certificate",
					Status: http.StatusInternalServerError,
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrListCertificateBindings)
			},
		},
		"validation error - missing CertificateID": {
			params:       ListCertificateBindingsRequest{},
			expectedPath: "/ccm/v1/certificates/123/certificate-bindings",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificate bindings: struct validation: CertificateID: cannot be blank",
					err.Error())
				assert.ErrorIs(t, err, ErrListCertificateBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page size less than 1": {
			params: ListCertificateBindingsRequest{
				CertificateID: "123",
				PageSize:      -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificate bindings: struct validation: PageSize: must be 1 or greater", err.Error())
				assert.ErrorIs(t, err, ErrListCertificateBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page size greater than 100": {
			params: ListCertificateBindingsRequest{
				CertificateID: "123",
				PageSize:      101,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificate bindings: struct validation: PageSize: cannot be greater than 100", err.Error())
				assert.ErrorIs(t, err, ErrListCertificateBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page value less than 1": {
			params: ListCertificateBindingsRequest{
				CertificateID: "123",
				Page:          -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificate bindings: struct validation: Page: must be 1 or greater", err.Error())
				assert.ErrorIs(t, err, ErrListCertificateBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				if len(tc.returnedHeaders) > 0 {
					for header, value := range tc.returnedHeaders {
						w.Header().Set(header, value)
					}
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.ListCertificateBindings(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListBindings(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params           ListBindingsRequest
		responseStatus   int
		responseBody     string
		expectedResponse *ListBindingsResponse
		expectedPath     string
		returnedHeaders  map[string]string
		withError        func(*testing.T, error)
	}{
		"200 - fetch of bindings successful": {
			params: ListBindingsRequest{},
			returnedHeaders: map[string]string{
				"Akamai-RateLimit-Limit":     "60",
				"Akamai-RateLimit-Remaining": "59",
			},
			expectedResponse: &ListBindingsResponse{
				Bindings: []CertificateBinding{
					{
						CertificateID: "123456",
						Hostname:      "www.example.com",
						Network:       "PRODUCTION",
						ResourceType:  "CDN_HOSTNAME",
					},
					{
						CertificateID: "654321",
						Hostname:      "secure.example.com",
						Network:       "STAGING",
						ResourceType:  "CDN_HOSTNAME",
					},
					{
						CertificateID: "789012",
						Hostname:      "api.example.com",
						Network:       "PRODUCTION",
						ResourceType:  "CDN_HOSTNAME",
					},
					// There should be 10 bindings in the response to match links correctly, but only 3 are shown here for brevity.
				},
				Links: Links{
					Next:     ptr.To("https://api.example.com/v1/certificate-bindings?page=2&pageSize=10"),
					Previous: nil,
					Self:     "https://api.example.com/v1/certificate-bindings?page=1&pageSize=10",
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			responseStatus: 200,
			expectedPath:   "/ccm/v1/certificate-bindings",
			responseBody: `
				{
				  "bindings": [
					{
					  "certificateId": 123456,
					  "hostname": "www.example.com",
					  "network": "PRODUCTION",
					  "resourceType": "CDN_HOSTNAME"
					},
					{
					  "certificateId": "654321",
					  "hostname": "secure.example.com",
					  "network": "STAGING",
					  "resourceType": "CDN_HOSTNAME"
					},
					{
					  "certificateId": "789012",
					  "hostname": "api.example.com",
					  "network": "PRODUCTION",
					  "resourceType": "CDN_HOSTNAME"
					}
				  ],
				  "links": {
					"next": "https://api.example.com/v1/certificate-bindings?page=2&pageSize=10",
					"previous": null,
					"self": "https://api.example.com/v1/certificate-bindings?page=1&pageSize=10"
				  }
				}`,
		},
		"200 - fetch of bindings with paging successful": {
			params: ListBindingsRequest{
				Page:     3,
				PageSize: 1,
			},
			expectedResponse: &ListBindingsResponse{
				Bindings: []CertificateBinding{
					{
						CertificateID: "789012",
						Hostname:      "api.example.com",
						Network:       "PRODUCTION",
						ResourceType:  "CDN_HOSTNAME",
					},
				},
				Links: Links{
					Next:     nil,
					Previous: ptr.To("https://api.example.com/v1/certificate-bindings?page=2&pageSize=1"),
					Self:     "https://api.example.com/v1/certificate-bindings?page=3&pageSize=1",
				},
			},
			responseStatus: 200,
			expectedPath:   "/ccm/v1/certificate-bindings?page=3&pageSize=1",
			responseBody: `
				{
				  "bindings": [
					{
					  "certificateId": "789012",
					  "hostname": "api.example.com",
					  "network": "PRODUCTION",
					  "resourceType": "CDN_HOSTNAME"
					}
				  ],
				  "links": {
					"next": null,
					"previous": "https://api.example.com/v1/certificate-bindings?page=2&pageSize=1",
					"self": "https://api.example.com/v1/certificate-bindings?page=3&pageSize=1"
				  }
				}`,
		},
		"200 - fetch of bindings with all filters": {
			params: ListBindingsRequest{
				ContractID:     "12345",
				GroupID:        "999",
				Domain:         "api.example.com",
				ExpiringInDays: ptr.To(int64(30)),
				Network:        "PRODUCTION",
				Page:           3,
				PageSize:       1,
			},
			expectedResponse: &ListBindingsResponse{
				Bindings: []CertificateBinding{
					{
						CertificateID: "789012",
						Hostname:      "api.example.com",
						Network:       "PRODUCTION",
						ResourceType:  "CDN_HOSTNAME",
					},
				},
				Links: Links{
					Next:     nil,
					Previous: ptr.To("https://api.example.com/v1/certificate-bindings?page=2&pageSize=1"),
					Self:     "https://api.example.com/v1/certificate-bindings?page=3&pageSize=1",
				},
			},
			responseStatus: 200,
			expectedPath:   "/ccm/v1/certificate-bindings?contractId=12345&domain=api.example.com&expiringInDays=30&groupId=999&network=PRODUCTION&page=3&pageSize=1",
			responseBody: `
				{
				  "bindings": [
					{
					  "certificateId": "789012",
					  "hostname": "api.example.com",
					  "network": "PRODUCTION",
					  "resourceType": "CDN_HOSTNAME"
					}
				  ],
				  "links": {
					"next": null,
					"previous": "https://api.example.com/v1/certificate-bindings?page=2&pageSize=1",
					"self": "https://api.example.com/v1/certificate-bindings?page=3&pageSize=1"
				  }
				}`,
		},
		"200 - empty response": {
			params: ListBindingsRequest{
				ContractID:     "12345",
				GroupID:        "999",
				Domain:         "foo.example.com",
				ExpiringInDays: ptr.To(int64(30)),
				Network:        "PRODUCTION",
			},
			expectedResponse: &ListBindingsResponse{
				Bindings: []CertificateBinding{},
				Links: Links{
					Next:     nil,
					Previous: nil,
					Self:     "",
				},
			},
			responseStatus: 200,
			expectedPath:   "/ccm/v1/certificate-bindings?contractId=12345&domain=foo.example.com&expiringInDays=30&groupId=999&network=PRODUCTION",
			responseBody: `
				{
				  "bindings": [],
				  "links": {
					"next": null,
					"previous": null,
					"self": null
				  }
				}`,
		},
		"500 internal server error - assert that error is ErrListBindings": {
			params:         ListBindingsRequest{},
			responseStatus: 500,
			responseBody: `
			{
				"instance": "/error-types/internal-error?traceId=-11111",
				"status": 500,
				"title": "An unexpected error occurred.",
				"type": "/error-types/internal-error"
			}`,
			expectedPath: "/ccm/v1/certificate-bindings",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListBindings, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=-11111",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrListBindings)
			},
		},
		"validation error - invalid network": {
			params: ListBindingsRequest{
				PageSize: 1,
				Network:  "foo",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing bindings: struct validation: Network: must be either 'STAGING' or 'PRODUCTION'", err.Error())
				assert.ErrorIs(t, err, ErrListBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page size less than 1": {
			params: ListBindingsRequest{
				PageSize: -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing bindings: struct validation: PageSize: must be 1 or greater", err.Error())
				assert.ErrorIs(t, err, ErrListBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page size greater than 100": {
			params: ListBindingsRequest{
				PageSize: 101,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing bindings: struct validation: PageSize: cannot be greater than 100", err.Error())
				assert.ErrorIs(t, err, ErrListBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page value less than 1": {
			params: ListBindingsRequest{
				Page: -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing bindings: struct validation: Page: must be 1 or greater", err.Error())
				assert.ErrorIs(t, err, ErrListBindings)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				if len(tc.returnedHeaders) > 0 {
					for header, value := range tc.returnedHeaders {
						w.Header().Set(header, value)
					}
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.ListBindings(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
