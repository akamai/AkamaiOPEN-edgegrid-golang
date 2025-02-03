package cps

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

func TestListDeployments(t *testing.T) {
	tests := map[string]struct {
		params           ListDeploymentsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedHeaders  map[string]string
		expectedResponse *ListDeploymentsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         ListDeploymentsRequest{EnrollmentID: 10},
			responseStatus: http.StatusOK,
			responseBody: `{
  "production": {
    "ocspStapled": false,
    "ocspUris": [],
    "networkConfiguration": {
      "geography": "core",
      "mustHaveCiphers": "ak-akamai-2020q1",
      "ocspStapling": "on",
      "preferredCiphers": "ak-akamai-2020q1",
      "quicEnabled": false,
      "secureNetwork": "standard-tls",
      "sniOnly": true,
      "disallowedTlsVersions": [
        "TLSv1",
        "TLSv1_1"
      ],
      "dnsNames": [
        "san2.example.com",
        "san1.example.com"
      ]
    },
    "primaryCertificate": {
      "certificate": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 93Nvw==\n-----END CERTIFICATE-----",
      "expiry": "2022-02-05T13:21:21Z",
      "signatureAlgorithm": "SHA-256",
      "trustChain": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... Qs/v0=\n-----END CERTIFICATE-----"
    },
    "multiStackedCertificates": [
      {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... nMweq/\n-----END CERTIFICATE-----",
        "expiry": "2022-02-05T13:21:20Z",
        "signatureAlgorithm": "SHA-256",
        "trustChain": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... KEUp0=\n-----END CERTIFICATE-----"
      }
    ]
  },
  "staging": {
    "ocspStapled": false,
    "ocspUris": [],
    "networkConfiguration": {
      "geography": "core",
      "mustHaveCiphers": "ak-akamai-2020q1",
      "ocspStapling": "on",
      "preferredCiphers": "ak-akamai-2020q1",
      "quicEnabled": false,
      "secureNetwork": "standard-tls",
      "sniOnly": true,
      "disallowedTlsVersions": [
        "TLSv1",
        "TLSv1_1"
      ],
      "dnsNames": [
        "san2.example.com",
        "san1.example.com"
      ]
    },
    "primaryCertificate": {
      "certificate": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 93Nvw==\n-----END CERTIFICATE-----",
      "expiry": "2022-02-05T13:21:21Z",
      "signatureAlgorithm": "SHA-256",
      "trustChain": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 9JQs/v0=\n-----END CERTIFICATE-----"
    },
    "multiStackedCertificates": [
      {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... nMweq/\n-----END CERTIFICATE-----",
        "expiry": "2022-02-05T13:21:20Z",
        "signatureAlgorithm": "SHA-256",
        "trustChain": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... KEUp0=\n-----END CERTIFICATE-----"
      }
    ]
  }
}`,
			expectedPath: "/cps/v2/enrollments/10/deployments",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.deployments.v8+json",
			},
			expectedResponse: &ListDeploymentsResponse{
				Production: &Deployment{
					OCSPStapled: ptr.To(false),
					OCSPURIs:    []string{},
					NetworkConfiguration: DeploymentNetworkConfiguration{
						Geography:        "core",
						MustHaveCiphers:  "ak-akamai-2020q1",
						OCSPStapling:     "on",
						PreferredCiphers: "ak-akamai-2020q1",
						QUICEnabled:      false,
						SecureNetwork:    "standard-tls",
						SNIOnly:          true,
						DisallowedTLSVersions: []string{
							"TLSv1",
							"TLSv1_1",
						},
						DNSNames: []string{
							"san2.example.com",
							"san1.example.com",
						},
					},
					PrimaryCertificate: DeploymentCertificate{
						Certificate:        "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 93Nvw==\n-----END CERTIFICATE-----",
						Expiry:             "2022-02-05T13:21:21Z",
						SignatureAlgorithm: "SHA-256",
						TrustChain:         "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... Qs/v0=\n-----END CERTIFICATE-----",
					},
					MultiStackedCertificates: []DeploymentCertificate{
						{
							Certificate:        "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... nMweq/\n-----END CERTIFICATE-----",
							Expiry:             "2022-02-05T13:21:20Z",
							SignatureAlgorithm: "SHA-256",
							TrustChain:         "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... KEUp0=\n-----END CERTIFICATE-----",
						},
					},
				},
				Staging: &Deployment{
					OCSPStapled: ptr.To(false),
					OCSPURIs:    []string{},
					NetworkConfiguration: DeploymentNetworkConfiguration{
						Geography:        "core",
						MustHaveCiphers:  "ak-akamai-2020q1",
						OCSPStapling:     "on",
						PreferredCiphers: "ak-akamai-2020q1",
						QUICEnabled:      false,
						SecureNetwork:    "standard-tls",
						SNIOnly:          true,
						DisallowedTLSVersions: []string{
							"TLSv1",
							"TLSv1_1",
						},
						DNSNames: []string{
							"san2.example.com",
							"san1.example.com",
						},
					},
					PrimaryCertificate: DeploymentCertificate{
						Certificate:        "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 93Nvw==\n-----END CERTIFICATE-----",
						Expiry:             "2022-02-05T13:21:21Z",
						SignatureAlgorithm: "SHA-256",
						TrustChain:         "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 9JQs/v0=\n-----END CERTIFICATE-----",
					},
					MultiStackedCertificates: []DeploymentCertificate{
						{
							Certificate:        "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... nMweq/\n-----END CERTIFICATE-----",
							Expiry:             "2022-02-05T13:21:20Z",
							SignatureAlgorithm: "SHA-256",
							TrustChain:         "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... KEUp0=\n-----END CERTIFICATE-----",
						},
					},
				},
			},
		},
		"500 internal server error": {
			params:         ListDeploymentsRequest{EnrollmentID: 500},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/500/deployments",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.deployments.v8+json",
			},
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
				for k, v := range test.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListDeployments(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetProductionDeployment(t *testing.T) {
	tests := map[string]struct {
		params           GetDeploymentRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedHeaders  map[string]string
		expectedResponse *GetProductionDeploymentResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         GetDeploymentRequest{EnrollmentID: 10},
			responseStatus: http.StatusOK,
			responseBody: `{
    "ocspStapled": false,
    "ocspUris": [],
    "networkConfiguration": {
      "geography": "core",
      "mustHaveCiphers": "ak-akamai-2020q1",
      "ocspStapling": "on",
      "preferredCiphers": "ak-akamai-2020q1",
      "quicEnabled": false,
      "secureNetwork": "standard-tls",
      "sniOnly": true,
      "disallowedTlsVersions": [
        "TLSv1",
        "TLSv1_1"
      ],
      "dnsNames": [
        "san2.example.com",
        "san1.example.com"
      ]
    },
    "primaryCertificate": {
      "certificate": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 93Nvw==\n-----END CERTIFICATE-----",
      "expiry": "2022-02-05T13:21:21Z",
      "signatureAlgorithm": "SHA-256",
      "trustChain": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... Qs/v0=\n-----END CERTIFICATE-----"
    },
    "multiStackedCertificates": [
      {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... nMweq/\n-----END CERTIFICATE-----",
        "expiry": "2022-02-05T13:21:20Z",
        "signatureAlgorithm": "SHA-256",
        "trustChain": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... KEUp0=\n-----END CERTIFICATE-----"
      }
    ]
}`,
			expectedPath: "/cps/v2/enrollments/10/deployments/production",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.deployment.v8+json",
			},
			expectedResponse: &GetProductionDeploymentResponse{
				OCSPStapled: ptr.To(false),
				OCSPURIs:    []string{},
				NetworkConfiguration: DeploymentNetworkConfiguration{
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QUICEnabled:      false,
					SecureNetwork:    "standard-tls",
					SNIOnly:          true,
					DisallowedTLSVersions: []string{
						"TLSv1",
						"TLSv1_1",
					},
					DNSNames: []string{
						"san2.example.com",
						"san1.example.com",
					},
				},
				PrimaryCertificate: DeploymentCertificate{
					Certificate:        "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 93Nvw==\n-----END CERTIFICATE-----",
					Expiry:             "2022-02-05T13:21:21Z",
					SignatureAlgorithm: "SHA-256",
					TrustChain:         "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... Qs/v0=\n-----END CERTIFICATE-----",
				},
				MultiStackedCertificates: []DeploymentCertificate{
					{
						Certificate:        "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... nMweq/\n-----END CERTIFICATE-----",
						Expiry:             "2022-02-05T13:21:20Z",
						SignatureAlgorithm: "SHA-256",
						TrustChain:         "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... KEUp0=\n-----END CERTIFICATE-----",
					},
				},
			},
		},
		"500 internal server error": {
			params:         GetDeploymentRequest{EnrollmentID: 500},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/500/deployments/production",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.deployment.v8+json",
			},
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error": {
			params: GetDeploymentRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				for k, v := range test.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetProductionDeployment(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetStagingDeployment(t *testing.T) {
	tests := map[string]struct {
		params           GetDeploymentRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedHeaders  map[string]string
		expectedResponse *GetStagingDeploymentResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         GetDeploymentRequest{EnrollmentID: 10},
			responseStatus: http.StatusOK,
			responseBody: `{
    "ocspStapled": false,
    "ocspUris": [],
    "networkConfiguration": {
      "geography": "core",
      "mustHaveCiphers": "ak-akamai-2020q1",
      "ocspStapling": "on",
      "preferredCiphers": "ak-akamai-2020q1",
      "quicEnabled": false,
      "secureNetwork": "standard-tls",
      "sniOnly": true,
      "disallowedTlsVersions": [
        "TLSv1",
        "TLSv1_1"
      ],
      "dnsNames": [
        "san2.example.com",
        "san1.example.com"
      ]
    },
    "primaryCertificate": {
      "certificate": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 93Nvw==\n-----END CERTIFICATE-----",
      "expiry": "2022-02-05T13:21:21Z",
      "signatureAlgorithm": "SHA-256",
      "trustChain": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... Qs/v0=\n-----END CERTIFICATE-----"
    },
    "multiStackedCertificates": [
      {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... nMweq/\n-----END CERTIFICATE-----",
        "expiry": "2022-02-05T13:21:20Z",
        "signatureAlgorithm": "SHA-256",
        "trustChain": "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... KEUp0=\n-----END CERTIFICATE-----"
      }
    ]
}`,
			expectedPath: "/cps/v2/enrollments/10/deployments/staging",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.deployment.v8+json",
			},
			expectedResponse: &GetStagingDeploymentResponse{

				OCSPStapled: ptr.To(false),
				OCSPURIs:    []string{},
				NetworkConfiguration: DeploymentNetworkConfiguration{
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QUICEnabled:      false,
					SecureNetwork:    "standard-tls",
					SNIOnly:          true,
					DisallowedTLSVersions: []string{
						"TLSv1",
						"TLSv1_1",
					},
					DNSNames: []string{
						"san2.example.com",
						"san1.example.com",
					},
				},
				PrimaryCertificate: DeploymentCertificate{
					Certificate:        "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... 93Nvw==\n-----END CERTIFICATE-----",
					Expiry:             "2022-02-05T13:21:21Z",
					SignatureAlgorithm: "SHA-256",
					TrustChain:         "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... Qs/v0=\n-----END CERTIFICATE-----",
				},
				MultiStackedCertificates: []DeploymentCertificate{
					{
						Certificate:        "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... nMweq/\n-----END CERTIFICATE-----",
						Expiry:             "2022-02-05T13:21:20Z",
						SignatureAlgorithm: "SHA-256",
						TrustChain:         "-----BEGIN CERTIFICATE-----\nMIID <sample - removed for readability> .... KEUp0=\n-----END CERTIFICATE-----",
					},
				},
			},
		},
		"500 internal server error": {
			params:         GetDeploymentRequest{EnrollmentID: 500},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/500/deployments/staging",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.deployment.v8+json",
			},
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error": {
			params: GetDeploymentRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				for k, v := range test.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetStagingDeployment(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
