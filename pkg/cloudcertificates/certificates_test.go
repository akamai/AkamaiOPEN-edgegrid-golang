package cloudcertificates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCertificate(t *testing.T) {
	t.Parallel()

	baseRequest := CreateCertificateRequest{
		ContractID: "111",
		GroupID:    "222",
		Body: CreateCertificateRequestBody{
			CertificateName: "test-cert",
			SANs:            []string{"example.com", "www.example.com"},
			SecureNetwork:   "ENHANCED_TLS",
			KeyType:         "RSA",
			KeySize:         "2048",
			Subject: &Subject{
				CommonName:   "example.com",
				Country:      "US",
				State:        "Massachusetts",
				Locality:     "Cambridge",
				Organization: "ExampleOrg",
			},
		},
	}

	baseResponseBody := `{
		"accountId": "A-CCT7890",
		"certificateId": "123",
		"certificateName": "test-cert",
		"certificateStatus": "CSR_READY",
		"certificateType": "THIRD_PARTY",
		"contractId": "C-0N7RAC7",
		"createdBy": "jsmith",
		"createdDate": "2025-09-01T06:16:05.952613Z",
		"csrExpirationDate": "2026-11-03T06:16:07Z",
		"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
		"keySize": "2048",
		"keyType": "RSA",
		"modifiedBy": "jsmith",
		"modifiedDate": "2025-09-02T06:16:05.952613Z",
		"sans": [
			"example.com",
			"www.example.com"
		],
		"secureNetwork": "ENHANCED_TLS",
		"signedCertificateIssuer": null,
		"signedCertificateNotValidAfterDate": null,
		"signedCertificateNotValidBeforeDate": null,
		"signedCertificatePem": null,
		"signedCertificateSHA256Fingerprint": null,
		"signedCertificateSerialNumber": null,
		"subject": {
			"commonName": "example.com",
			"country": "US",
			"locality": "Cambridge",
			"organization": "ExampleOrg",
			"state": "Massachusetts"
		},
		"trustChainPem": null
	}`

	expectedResponseWithoutRateLimits := &CreateCertificateResponse{
		Certificate: Certificate{
			AccountID:         "A-CCT7890",
			CertificateID:     "123",
			CertificateName:   "test-cert",
			CertificateStatus: "CSR_READY",
			CertificateType:   "THIRD_PARTY",
			ContractID:        "C-0N7RAC7",
			CreatedBy:         "jsmith",
			CreatedDate:       test.NewTimeFromString(t, "2025-09-01T06:16:05.952613Z"),
			CSRExpirationDate: test.NewTimeFromString(t, "2026-11-03T06:16:07Z"),
			CSRPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
			KeySize:           "2048",
			KeyType:           "RSA",
			ModifiedBy:        "jsmith",
			ModifiedDate:      test.NewTimeFromString(t, "2025-09-02T06:16:05.952613Z"),
			SANs:              []string{"example.com", "www.example.com"},
			SecureNetwork:     "ENHANCED_TLS",
			Subject: &Subject{
				Country:      "US",
				Organization: "ExampleOrg",
				State:        "Massachusetts",
				Locality:     "Cambridge",
				CommonName:   "example.com",
			},
		},
		ResourceLimits: ResourceLimitsMetadata{
			CertificateLimitTotal:     ptr.To(int64(50)),
			CertificateLimitRemaining: ptr.To(int64(27)),
		},
		RateLimits: RateLimitsMetadata{
			Limit:     nil,
			Remaining: nil,
		},
	}

	tests := map[string]struct {
		params              CreateCertificateRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedResponse    *CreateCertificateResponse
		expectedPath        string
		returnedHeaders     map[string]string
		withError           func(*testing.T, error)
	}{
		"201 Created - create certificate with all possible fields": {
			params:              baseRequest,
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert","keyType":"RSA","keySize":"2048","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			returnedHeaders: map[string]string{
				"Akamai-Limit-Certificates":           "50",
				"Akamai-Limit-Certificates-Remaining": "27",
				"Akamai-RateLimit-Limit":              "60",
				"Akamai-RateLimit-Remaining":          "59",
			},
			expectedResponse: &CreateCertificateResponse{
				Certificate: Certificate{
					AccountID:         "A-CCT7890",
					CertificateID:     "123",
					CertificateName:   "test-cert",
					CertificateStatus: "CSR_READY",
					CertificateType:   "THIRD_PARTY",
					ContractID:        "C-0N7RAC7",
					CreatedBy:         "jsmith",
					CreatedDate:       test.NewTimeFromString(t, "2025-09-01T06:16:05.952613Z"),
					CSRExpirationDate: test.NewTimeFromString(t, "2026-11-03T06:16:07Z"),
					CSRPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:           "2048",
					KeyType:           "RSA",
					ModifiedBy:        "jsmith",
					ModifiedDate:      test.NewTimeFromString(t, "2025-09-02T06:16:05.952613Z"),
					SANs:              []string{"example.com", "www.example.com"},
					SecureNetwork:     "ENHANCED_TLS",
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
				ResourceLimits: ResourceLimitsMetadata{
					CertificateLimitTotal:     ptr.To(int64(50)),
					CertificateLimitRemaining: ptr.To(int64(27)),
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			responseStatus: 201,
			responseBody:   baseResponseBody,
		},
		"201 Created - no rate limit headers": {
			params:              baseRequest,
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert","keyType":"RSA","keySize":"2048","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			returnedHeaders: map[string]string{
				"Akamai-Limit-Certificates":           "50",
				"Akamai-Limit-Certificates-Remaining": "27",
			},
			expectedResponse: expectedResponseWithoutRateLimits,
			responseStatus:   201,
			responseBody:     baseResponseBody,
		},
		"201 Created - empty rate limit headers": {
			params:              baseRequest,
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert","keyType":"RSA","keySize":"2048","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			returnedHeaders: map[string]string{
				"Akamai-Limit-Certificates":           "50",
				"Akamai-Limit-Certificates-Remaining": "27",
				"Akamai-RateLimit-Limit":              "",
				"Akamai-RateLimit-Remaining":          "",
			},
			expectedResponse: expectedResponseWithoutRateLimits,
			responseStatus:   201,
			responseBody:     baseResponseBody,
		},
		"201 Created - create certificate without name": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					SANs:          []string{"example.com", "www.example.com"},
					SecureNetwork: "ENHANCED_TLS",
					KeyType:       "RSA",
					KeySize:       "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"keyType":"RSA","keySize":"2048","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			returnedHeaders: map[string]string{
				"Akamai-Limit-Certificates":           "50",
				"Akamai-Limit-Certificates-Remaining": "27",
				"Akamai-RateLimit-Limit":              "60",
				"Akamai-RateLimit-Remaining":          "59",
			},
			expectedResponse: &CreateCertificateResponse{
				Certificate: Certificate{
					AccountID:         "A-CCT7890",
					CertificateID:     "123",
					CertificateName:   "example.com20250111105236090681",
					CertificateStatus: "CSR_READY",
					CertificateType:   "THIRD_PARTY",
					ContractID:        "C-0N7RAC7",
					CreatedBy:         "jsmith",
					CreatedDate:       test.NewTimeFromString(t, "2025-09-01T06:16:05.952613Z"),
					CSRExpirationDate: test.NewTimeFromString(t, "2026-11-03T06:16:07Z"),
					CSRPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:           "2048",
					KeyType:           "RSA",
					ModifiedBy:        "jsmith",
					ModifiedDate:      test.NewTimeFromString(t, "2025-09-02T06:16:05.952613Z"),
					SANs:              []string{"example.com", "www.example.com"},
					SecureNetwork:     "ENHANCED_TLS",
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
				ResourceLimits: ResourceLimitsMetadata{
					CertificateLimitTotal:     ptr.To(int64(50)),
					CertificateLimitRemaining: ptr.To(int64(27)),
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			responseStatus: 201,
			responseBody: `{
				"accountId": "A-CCT7890",
				"certificateId": "123",
				"certificateName": "example.com20250111105236090681",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "C-0N7RAC7",
				"createdBy": "jsmith",
				"createdDate": "2025-09-01T06:16:05.952613Z",
				"csrExpirationDate": "2026-11-03T06:16:07Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "jsmith",
				"modifiedDate": "2025-09-02T06:16:05.952613Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
		},
		"201 Created - create certificate but cannot parse limit headers": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert","keyType":"RSA","keySize":"2048","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			returnedHeaders: map[string]string{
				"Akamai-Limit-Certificates":           "not-parsable",
				"Akamai-Limit-Certificates-Remaining": "not-parsable",
			},
			expectedResponse: &CreateCertificateResponse{
				Certificate: Certificate{
					AccountID:         "A-CCT7890",
					CertificateID:     "123",
					CertificateName:   "test-cert",
					CertificateStatus: "CSR_READY",
					CertificateType:   "THIRD_PARTY",
					ContractID:        "C-0N7RAC7",
					CreatedBy:         "jsmith",
					CreatedDate:       test.NewTimeFromString(t, "2025-09-01T06:16:05.952613Z"),
					CSRExpirationDate: test.NewTimeFromString(t, "2026-11-03T06:16:07Z"),
					CSRPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:           "2048",
					KeyType:           "RSA",
					ModifiedBy:        "jsmith",
					ModifiedDate:      test.NewTimeFromString(t, "2025-09-02T06:16:05.952613Z"),
					SANs:              []string{"example.com", "www.example.com"},
					SecureNetwork:     "ENHANCED_TLS",
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
				ResourceLimits: ResourceLimitsMetadata{
					CertificateLimitTotal:     nil,
					CertificateLimitRemaining: nil,
				},
			},
			responseStatus: 201,
			responseBody: `{
				"accountId": "A-CCT7890",
				"certificateId": "123",
				"certificateName": "test-cert",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "C-0N7RAC7",
				"createdBy": "jsmith",
				"createdDate": "2025-09-01T06:16:05.952613Z",
				"csrExpirationDate": "2026-11-03T06:16:07Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "jsmith",
				"modifiedDate": "2025-09-02T06:16:05.952613Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
		},
		"201 Created - create certificate with no subject": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
				},
			},
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert","keyType":"RSA","keySize":"2048","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"]}`,
			expectedResponse: &CreateCertificateResponse{
				Certificate: Certificate{
					AccountID:         "A-CCT7890",
					CertificateID:     "123",
					CertificateName:   "test-cert",
					CertificateStatus: "CSR_READY",
					CertificateType:   "THIRD_PARTY",
					ContractID:        "C-0N7RAC7",
					CreatedBy:         "jsmith",
					CreatedDate:       test.NewTimeFromString(t, "2025-09-01T06:16:05.952613Z"),
					CSRExpirationDate: test.NewTimeFromString(t, "2026-11-03T06:16:07Z"),
					CSRPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:           "2048",
					KeyType:           "RSA",
					ModifiedBy:        "jsmith",
					ModifiedDate:      test.NewTimeFromString(t, "2025-09-02T06:16:05.952613Z"),
					SANs:              []string{"example.com", "www.example.com"},
					SecureNetwork:     "ENHANCED_TLS",
				},
			},
			responseStatus: 201,
			responseBody: `{
				"accountId": "A-CCT7890",
				"certificateId": "123",
				"certificateName": "test-cert",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "C-0N7RAC7",
				"createdBy": "jsmith",
				"createdDate": "2025-09-01T06:16:05.952613Z",
				"csrExpirationDate": "2026-11-03T06:16:07Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "jsmith",
				"modifiedDate": "2025-09-02T06:16:05.952613Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"trustChainPem": null
			}`,
		},
		"201 Created - create certificate with ECDSA P-384 key size": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert-ecdsa-p384",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "ECDSA",
					KeySize:         "P-384",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert-ecdsa-p384","keyType":"ECDSA","keySize":"P-384","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			returnedHeaders: map[string]string{
				"Akamai-Limit-Certificates":           "50",
				"Akamai-Limit-Certificates-Remaining": "27",
				"Akamai-RateLimit-Limit":              "60",
				"Akamai-RateLimit-Remaining":          "59",
			},
			expectedResponse: &CreateCertificateResponse{
				Certificate: Certificate{
					AccountID:         "A-CCT7890",
					CertificateID:     "456",
					CertificateName:   "test-cert-ecdsa-p384",
					CertificateStatus: "CSR_READY",
					CertificateType:   "THIRD_PARTY",
					ContractID:        "C-0N7RAC7",
					CreatedBy:         "jsmith",
					CreatedDate:       test.NewTimeFromString(t, "2025-09-01T06:16:05.952613Z"),
					CSRExpirationDate: test.NewTimeFromString(t, "2026-11-03T06:16:07Z"),
					CSRPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nexample-ECDSA-P384-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:           "P-384",
					KeyType:           "ECDSA",
					ModifiedBy:        "jsmith",
					ModifiedDate:      test.NewTimeFromString(t, "2025-09-02T06:16:05.952613Z"),
					SANs:              []string{"example.com", "www.example.com"},
					SecureNetwork:     "ENHANCED_TLS",
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
				ResourceLimits: ResourceLimitsMetadata{
					CertificateLimitTotal:     ptr.To(int64(50)),
					CertificateLimitRemaining: ptr.To(int64(27)),
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			responseStatus: 201,
			responseBody: `{
				"accountId": "A-CCT7890",
				"certificateId": "456",
				"certificateName": "test-cert-ecdsa-p384",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "C-0N7RAC7",
				"createdBy": "jsmith",
				"createdDate": "2025-09-01T06:16:05.952613Z",
				"csrExpirationDate": "2026-11-03T06:16:07Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-ECDSA-P384-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "P-384",
				"keyType": "ECDSA",
				"modifiedBy": "jsmith",
				"modifiedDate": "2025-09-02T06:16:05.952613Z",
				"sans": ["example.com", "www.example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
		},
		"201 Created - create certificate with STANDARD_TLS secure network": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "STANDARD_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert","keyType":"RSA","keySize":"2048","secureNetwork":"STANDARD_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			returnedHeaders: map[string]string{
				"Akamai-Limit-Certificates":           "50",
				"Akamai-Limit-Certificates-Remaining": "27",
				"Akamai-RateLimit-Limit":              "60",
				"Akamai-RateLimit-Remaining":          "59",
			},
			expectedResponse: &CreateCertificateResponse{
				Certificate: Certificate{
					AccountID:         "A-CCT7890",
					CertificateID:     "123",
					CertificateName:   "test-cert",
					CertificateStatus: "CSR_READY",
					CertificateType:   "THIRD_PARTY",
					ContractID:        "C-0N7RAC7",
					CreatedBy:         "jsmith",
					CreatedDate:       test.NewTimeFromString(t, "2025-09-01T06:16:05.952613Z"),
					CSRExpirationDate: test.NewTimeFromString(t, "2026-11-03T06:16:07Z"),
					CSRPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:           "2048",
					KeyType:           "RSA",
					ModifiedBy:        "jsmith",
					ModifiedDate:      test.NewTimeFromString(t, "2025-09-02T06:16:05.952613Z"),
					SANs:              []string{"example.com", "www.example.com"},
					SecureNetwork:     "STANDARD_TLS",
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
				ResourceLimits: ResourceLimitsMetadata{
					CertificateLimitTotal:     ptr.To(int64(50)),
					CertificateLimitRemaining: ptr.To(int64(27)),
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			responseStatus: 201,
			responseBody: `{
				"accountId": "A-CCT7890",
				"certificateId": "123",
				"certificateName": "test-cert",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "C-0N7RAC7",
				"createdBy": "jsmith",
				"createdDate": "2025-09-01T06:16:05.952613Z",
				"csrExpirationDate": "2026-11-03T06:16:07Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "jsmith",
				"modifiedDate": "2025-09-02T06:16:05.952613Z",
				"sans": ["example.com", "www.example.com"],
				"secureNetwork": "STANDARD_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
		},
		"validation error - missing required ContractID": {
			params: CreateCertificateRequest{
				GroupID: "123",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: ContractID: cannot be blank",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing required GroupID": {
			params: CreateCertificateRequest{
				ContractID: "111",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: GroupID: cannot be blank",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid certificate name": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert##",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: CertificateName: the input can only contain digits (1-9), letters (a-z, A-Z), spaces, hyphens, periods, and underscores.",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid key type": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "INVALID",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: KeyType: value 'INVALID' is invalid. Must be either 'RSA' or 'ECDSA'",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid key size": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "1111",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: KeySize: value '1111' is invalid. Must be one of: '2048', 'P-256', or 'P-384'",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid secure network": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "WRONG_NETWORK",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: SecureNetwork: value 'WRONG_NETWORK' is invalid. Must be either 'ENHANCED_TLS' or 'STANDARD_TLS'",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid subject country code length": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "TESTETSTES",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: Subject: Country: the length must be exactly 2",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid subject country code format": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "  ",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: Subject: Country: must be in a valid format",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid subject locality length": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     strings.Repeat("A", 129),
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: Subject: Locality: the length must be between 1 and 128",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid subject locality format": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "		",
						Organization: "ExampleOrg",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: Subject: Locality: must be in a valid format",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid subject state length": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        strings.Repeat("A", 129),
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			expectedPath: "/ccm/v1/certificates?contractId=111&groupId=222",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: Subject: State: the length must be between 1 and 128",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			}},
		"validation error - invalid subject state format": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "		",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			expectedPath: "/ccm/v1/certificates?contractId=111&groupId=222",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: Subject: State: must be in a valid format",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			}},
		"validation error - invalid subject organization length": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: strings.Repeat("A", 65),
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: Subject: Organization: the length must be between 1 and 64",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			}},
		"validation error - invalid subject organization format": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "		",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "creating certificate: struct validation: Subject: Organization: must be in a valid format",
					err.Error())
				assert.ErrorIs(t, err, ErrCreateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"409 - certificate name already in use": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			responseStatus:      409,
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert","keyType":"RSA","keySize":"2048","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			responseBody: `
			{
				"certificateIdentifier": "certificateName",
				"certificateIdentifierValue": "test-cert",
				"detail": "Certificate with {certificateName}: {test-cert} already exists with the current account Id!",
				"instance": "/error-types/certificate-name-already-in-use?traceId=-123",
				"status": 409,
				"title": "Certificate name already in use.",
				"type": "/error-types/certificate-name-already-in-use"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCreateCertificate, &Error{
					Type:                       "/error-types/certificate-name-already-in-use",
					Title:                      "Certificate name already in use.",
					Detail:                     "Certificate with {certificateName}: {test-cert} already exists with the current account Id!",
					Status:                     http.StatusConflict,
					Instance:                   "/error-types/certificate-name-already-in-use?traceId=-123",
					CertificateIdentifier:      "certificateName",
					CertificateIdentifierValue: "test-cert",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertificateNameInUse)
				assert.ErrorIs(t, err, ErrCreateCertificate)
			},
		},
		"400 - account not permitted to use P-384 key size": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert-ecdsa-p384",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "ECDSA",
					KeySize:         "P-384",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: `{"certificateName":"test-cert-ecdsa-p384","keyType":"ECDSA","keySize":"P-384","secureNetwork":"ENHANCED_TLS","sans":["example.com","www.example.com"],"subject":{"commonName":"example.com","organization":"ExampleOrg","country":"US","state":"Massachusetts","locality":"Cambridge"}}`,
			responseStatus:      400,
			responseBody: `{
				"detail": "Validation failed while executing the operation.",
				"errors": [
					{
						"accountId": "A-CCT1234",
						"detail": "Account {A-CCT1234} is not permitted to create certificate for the key size {P-384}.",
						"instance": "/error-types/account-key-size-not-allowed?traceId=123456789",
						"keySize": "P-384",
						"title": "Forbidden.",
						"type": "/error-types/account-key-size-not-allowed"
					}
				],
				"instance": "/error-types/validation-failure?traceId=123456789",
				"status": 400,
				"title": "Validation Failure.",
				"type": "/error-types/validation-failure"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCreateCertificate, &Error{
					Type:     "/error-types/validation-failure",
					Title:    "Validation Failure.",
					Detail:   "Validation failed while executing the operation.",
					Status:   http.StatusBadRequest,
					Instance: "/error-types/validation-failure?traceId=123456789",
					Errors: []SecondaryError{
						{
							Type:     "/error-types/account-key-size-not-allowed",
							Title:    "Forbidden.",
							Detail:   "Account {A-CCT1234} is not permitted to create certificate for the key size {P-384}.",
							Instance: "/error-types/account-key-size-not-allowed?traceId=123456789",
						},
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCreateCertificate)
			},
		},
		"500 internal server error - assert that error is ErrCreateCertificate": {
			params: CreateCertificateRequest{
				ContractID: "111",
				GroupID:    "222",
				Body: CreateCertificateRequestBody{
					CertificateName: "test-cert",
					SANs:            []string{"example.com", "www.example.com"},
					SecureNetwork:   "ENHANCED_TLS",
					KeyType:         "RSA",
					KeySize:         "2048",
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
					},
				},
			},
			expectedPath:        "/ccm/v1/certificates?contractId=111&groupId=222",
			expectedRequestBody: "{\"certificateName\":\"test-cert\",\"keyType\":\"RSA\",\"keySize\":\"2048\",\"secureNetwork\":\"ENHANCED_TLS\",\"sans\":[\"example.com\",\"www.example.com\"],\"subject\":{\"commonName\":\"example.com\",\"organization\":\"ExampleOrg\",\"country\":\"US\",\"state\":\"Massachusetts\",\"locality\":\"Cambridge\"}}",
			responseStatus:      500,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error removing certificate",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCreateCertificate, &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error removing certificate",
					Status: http.StatusInternalServerError,
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCreateCertificate)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(tc.returnedHeaders) > 0 {
					for header, value := range tc.returnedHeaders {
						w.Header().Set(header, value)
					}
				}
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRequestBody, string(requestBody))
				w.WriteHeader(tc.responseStatus)
				_, err = w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateCertificate(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetCertificate(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params           GetCertificateRequest
		responseStatus   int
		responseBody     string
		expectedResponse *GetCertificateResponse
		expectedPath     string
		expectedHeaders  map[string]string
		returnedHeaders  map[string]string
		withError        func(*testing.T, error)
	}{
		"200 -fetch of certificate successful": {
			params: GetCertificateRequest{
				CertificateID: "123",
			},
			returnedHeaders: map[string]string{
				"Akamai-RateLimit-Limit":     "60",
				"Akamai-RateLimit-Remaining": "59",
			},
			expectedResponse: &GetCertificateResponse{
				Certificate: Certificate{
					AccountID:                          "A-CCT7890",
					CertificateID:                      "123",
					CertificateName:                    "test-cert",
					CertificateStatus:                  "CSR_READY",
					CertificateType:                    "THIRD_PARTY",
					ContractID:                         "C-0N7RAC7",
					CreatedBy:                          "jsmith",
					CreatedDate:                        test.NewTimeFromString(t, "2025-09-01T06:16:05.952613Z"),
					CSRExpirationDate:                  test.NewTimeFromString(t, "2026-11-03T06:16:07Z"),
					CSRPEM:                             "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                            "2048",
					KeyType:                            "RSA",
					ModifiedBy:                         "jsmith",
					ModifiedDate:                       test.NewTimeFromString(t, "2025-09-02T06:16:05.952613Z"),
					SANs:                               []string{"example.com", "www.example.com"},
					SecureNetwork:                      "ENHANCED_TLS",
					SignedCertificateIssuer:            nil,
					SignedCertificateSerialNumber:      nil,
					SignedCertificateSHA256Fingerprint: nil,
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			responseStatus: 200,
			expectedPath:   "/ccm/v1/certificates/123",
			responseBody: `{
				"accountId": "A-CCT7890",
				"certificateId": "123",
				"certificateName": "test-cert",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "C-0N7RAC7",
				"createdBy": "jsmith",
				"createdDate": "2025-09-01T06:16:05.952613Z",
				"csrExpirationDate": "2026-11-03T06:16:07Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "jsmith",
				"modifiedDate": "2025-09-02T06:16:05.952613Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
		},
		"404 resource not found -certificate not found": {
			params: GetCertificateRequest{
				CertificateID: "1234",
			},
			responseStatus: 404,
			expectedPath:   "/ccm/v1/certificates/1234",
			responseBody: `{
				"certificateIdentifier": "certificateSubscriptionId",
				"certificateIdentifierValue": "1234",
				"detail": "Certificate subscription with {certificateSubscriptionId}: {1234} is not found.",
				"instance": "/error-types/certificate-not-found?traceId=-1234",
				"status": 404,
				"title": "Certificate subscription is not found.",
				"type": "/error-types/certificate-not-found"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetCertificate, &Error{
					Type:                       "/error-types/certificate-not-found",
					Title:                      "Certificate subscription is not found.",
					Detail:                     "Certificate subscription with {certificateSubscriptionId}: {1234} is not found.",
					Status:                     http.StatusNotFound,
					Instance:                   "/error-types/certificate-not-found?traceId=-1234",
					CertificateIdentifier:      "certificateSubscriptionId",
					CertificateIdentifierValue: "1234",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertificateNotFound)
				assert.ErrorIs(t, err, ErrGetCertificate)
			},
		},
		"500 internal server error - assert that error is ErrGetCertificate": {
			params: GetCertificateRequest{
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
			expectedPath: "/ccm/v1/certificates/123",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetCertificate, &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error removing certificate",
					Status: http.StatusInternalServerError,
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrGetCertificate)
			},
		},
		"validation error -missing CertificateID": {
			params:       GetCertificateRequest{},
			expectedPath: "/ccm/v1/certificates/123",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "getting certificate: struct validation: CertificateID: cannot be blank",
					err.Error())
				assert.ErrorIs(t, err, ErrGetCertificate)
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
				for k, v := range tc.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.GetCertificate(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestDeleteCertificate(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params           DeleteCertificateRequest
		responseStatus   int
		responseBody     string
		expectedResponse *DeleteCertificateResponse
		expectedPath     string
		returnedHeaders  map[string]string
		expectedHeaders  map[string]string
		withError        func(*testing.T, error)
	}{
		"204 certificate deleted successfully": {
			params: DeleteCertificateRequest{
				CertificateID: "certificate_123",
			},
			returnedHeaders: map[string]string{
				"Akamai-RateLimit-Limit":     "60",
				"Akamai-RateLimit-Remaining": "59",
			},
			expectedResponse: &DeleteCertificateResponse{
				Limit:     ptr.To(int64(60)),
				Remaining: ptr.To(int64(59)),
			},
			responseStatus: 204,
			expectedPath:   "/ccm/v1/certificates/certificate_123",
		},
		"404 resource not found -certificate not found": {
			params: DeleteCertificateRequest{
				CertificateID: "1234",
			},
			responseStatus: 404,
			expectedPath:   "/ccm/v1/certificates/1234",
			responseBody: `{
				"certificateIdentifier": "certificateSubscriptionId",
				"certificateIdentifierValue": "1234",
				"detail": "Certificate subscription with {certificateSubscriptionId}: {1234} is not found.",
				"instance": "/error-types/certificate-not-found?traceId=-1234",
				"status": 404,
				"title": "Certificate subscription is not found.",
				"type": "/error-types/certificate-not-found"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteCertificate, &Error{
					Type:                       "/error-types/certificate-not-found",
					Title:                      "Certificate subscription is not found.",
					Detail:                     "Certificate subscription with {certificateSubscriptionId}: {1234} is not found.",
					Status:                     http.StatusNotFound,
					Instance:                   "/error-types/certificate-not-found?traceId=-1234",
					CertificateIdentifier:      "certificateSubscriptionId",
					CertificateIdentifierValue: "1234",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertificateNotFound)
				assert.ErrorIs(t, err, ErrDeleteCertificate)
			},
		},
		"500 internal server error - assert that error is ErrDeleteCertificate": {
			params: DeleteCertificateRequest{
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
			expectedPath: "/ccm/v1/certificates/123",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteCertificate, &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error removing certificate",
					Status: http.StatusInternalServerError,
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrDeleteCertificate)
			},
		},
		"validate - missing CertificateID": {
			params:       DeleteCertificateRequest{},
			expectedPath: "/ccm/v1/certificates/certificate_123",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "deleting certificate: struct validation: CertificateID: cannot be blank",
					err.Error())
				assert.ErrorIs(t, err, ErrDeleteCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				if len(tc.returnedHeaders) > 0 {
					for header, value := range tc.returnedHeaders {
						w.Header().Set(header, value)
					}
				}
				w.WriteHeader(tc.responseStatus)
				if tc.responseBody != "" {
					_, err := w.Write([]byte(tc.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteCertificate(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestPatchCertificate(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params              PatchCertificateRequest
		responseStatus      int
		responseBody        string
		expectedResponse    *PatchCertificateResponse
		expectedRequestBody string
		expectedPath        string
		expectedHeaders     map[string]string
		returnedHeaders     map[string]string
		withError           func(*testing.T, error)
	}{
		"200 OK - only rename with all allowed characters": {
			params: PatchCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To("test 0123456789.-_"),
			},
			returnedHeaders: map[string]string{
				"Akamai-RateLimit-Limit":     "60",
				"Akamai-RateLimit-Remaining": "59",
			},
			expectedResponse: &PatchCertificateResponse{
				Certificate: Certificate{
					AccountID:               "acc_123",
					CertificateID:           "123",
					CertificateName:         "test 0123456789.-_",
					CertificateStatus:       "CSR_READY",
					CertificateType:         "THIRD_PARTY",
					ContractID:              "A-123",
					CreatedBy:               "user",
					CreatedDate:             test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:       test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                  "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                 "2048",
					KeyType:                 "RSA",
					ModifiedBy:              "user",
					ModifiedDate:            test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                    []string{"example.com", "www.example.com"},
					SecureNetwork:           "ENHANCED_TLS",
					SignedCertificateIssuer: nil,
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			expectedRequestBody: "[{\"op\":\"replace\",\"path\":\"/certificateName\",\"value\":\"test 0123456789.-_\"}]",
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "test 0123456789.-_",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
		},
		"200 OK - reset name by providing empty value": {
			params: PatchCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To(""),
			},
			expectedResponse: &PatchCertificateResponse{
				Certificate: Certificate{
					AccountID:               "acc_123",
					CertificateID:           "123",
					CertificateName:         "example.com20250822092651008941",
					CertificateStatus:       "CSR_READY",
					CertificateType:         "THIRD_PARTY",
					ContractID:              "A-123",
					CreatedBy:               "user",
					CreatedDate:             test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:       test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                  "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                 "2048",
					KeyType:                 "RSA",
					ModifiedBy:              "user",
					ModifiedDate:            test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                    []string{"example.com", "www.example.com"},
					SecureNetwork:           "ENHANCED_TLS",
					SignedCertificateIssuer: nil,
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
			},
			expectedRequestBody: "[{\"op\":\"replace\",\"path\":\"/certificateName\",\"value\":\"\"}]",
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "example.com20250822092651008941",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
		},
		"200 OK - upload signed certificate PEM only": {
			params: PatchCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
			},
			expectedResponse: &PatchCertificateResponse{
				Certificate: Certificate{
					AccountID:                           "acc_123",
					CertificateID:                       "123",
					CertificateName:                     "Certificate-name-rename",
					CertificateStatus:                   "READY_FOR_USE",
					CertificateType:                     "THIRD_PARTY",
					ContractID:                          "A-123",
					CreatedBy:                           "user",
					CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                             "2048",
					KeyType:                             "RSA",
					ModifiedBy:                          "user",
					ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                                []string{"example.com", "www.example.com"},
					SecureNetwork:                       "ENHANCED_TLS",
					SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
					SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
					SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
					SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
					},
					TrustChainPEM: nil,
				},
			},
			expectedRequestBody: `[{"op":"add","path":"/signedCertificatePem","value":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"}]`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "Certificate-name-rename",
				"certificateStatus": "READY_FOR_USE",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
				"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
				"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
				"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				"signedCertificateSHA256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
				"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
		},
		"200 OK - upload signed certificate PEM with trust chain": {
			params: PatchCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				TrustChainPEM:        "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n",
			},
			expectedResponse: &PatchCertificateResponse{
				Certificate: Certificate{
					AccountID:                           "acc_123",
					CertificateID:                       "123",
					CertificateName:                     "Certificate-name-rename",
					CertificateStatus:                   "READY_FOR_USE",
					CertificateType:                     "THIRD_PARTY",
					ContractID:                          "A-123",
					CreatedBy:                           "user",
					CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                             "2048",
					KeyType:                             "RSA",
					ModifiedBy:                          "user",
					ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                                []string{"example.com", "www.example.com"},
					SecureNetwork:                       "ENHANCED_TLS",
					SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
					SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
					SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
					SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
					},
					TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
				},
			},
			expectedRequestBody: `[{"op":"add","path":"/signedCertificatePem","value":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"},{"op":"add","path":"/trustChainPem","value":"-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"}]`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "Certificate-name-rename",
				"certificateStatus": "READY_FOR_USE",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
				"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
				"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
				"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				"signedCertificateSHA256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
				"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
		},
		"200 OK - upload signed certificate with AcknowledgeWarnings query param": {
			params: PatchCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				AcknowledgeWarnings:  true,
			},
			expectedResponse: &PatchCertificateResponse{
				Certificate: Certificate{
					AccountID:                           "acc_123",
					CertificateID:                       "123",
					CertificateName:                     "Certificate-name-rename",
					CertificateStatus:                   "READY_FOR_USE",
					CertificateType:                     "THIRD_PARTY",
					ContractID:                          "A-123",
					CreatedBy:                           "user",
					CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                             "2048",
					KeyType:                             "RSA",
					ModifiedBy:                          "user",
					ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                                []string{"example.com", "www.example.com"},
					SecureNetwork:                       "ENHANCED_TLS",
					SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
					SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
					SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
					SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
					},
					TrustChainPEM: nil,
				},
			},
			expectedRequestBody: `[{"op":"add","path":"/signedCertificatePem","value":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"}]`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "Certificate-name-rename",
				"certificateStatus": "READY_FOR_USE",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
				"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
				"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
				"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				"signedCertificateSHA256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
				"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
			expectedPath: "/ccm/v1/certificates/123?acknowledgeWarnings=true",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
		},
		"200 OK - upload signed certificate and trust chain PEM with AcknowledgeWarnings query param and rename certificate": {
			params: PatchCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				TrustChainPEM:        "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n",
				CertificateName:      ptr.To("Certificate-name-rename"),
				AcknowledgeWarnings:  true,
			},
			expectedResponse: &PatchCertificateResponse{
				Certificate: Certificate{
					AccountID:                           "acc_123",
					CertificateID:                       "123",
					CertificateName:                     "Certificate-name-rename",
					CertificateStatus:                   "READY_FOR_USE",
					CertificateType:                     "THIRD_PARTY",
					ContractID:                          "A-123",
					CreatedBy:                           "user",
					CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                             "2048",
					KeyType:                             "RSA",
					ModifiedBy:                          "user",
					ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                                []string{"example.com", "www.example.com"},
					SecureNetwork:                       "ENHANCED_TLS",
					SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
					SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
					SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
					SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
					},
					TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
				},
			},
			expectedRequestBody: `[{"op":"add","path":"/signedCertificatePem","value":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"},{"op":"add","path":"/trustChainPem","value":"-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"},{"op":"replace","path":"/certificateName","value":"Certificate-name-rename"}]`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "Certificate-name-rename",
				"certificateStatus": "READY_FOR_USE",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
				"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
				"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
				"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				"signedCertificateSHA256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
				"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"
			}`,
			expectedPath: "/ccm/v1/certificates/123?acknowledgeWarnings=true",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
		},
		"409 OK - warnings in the response for certificate pem": {
			params: PatchCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
			},
			expectedRequestBody: `[{"op":"add","path":"/signedCertificatePem","value":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"}]`,
			responseStatus:      409,
			responseBody: `
			{
				"data": {
					"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
					"signedCertificates": [
						{
							"certificatePem": "-----BEGIN CERTIFICATE-----\nexample-PEM\n-----END CERTIFICATE-----",
							"createdBy": null,
							"createdDate": null,
							"displayName": null,
							"endDate": "2027-11-22T12:45:19Z",
							"fingerprint": null,
							"issuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
							"sans": [
								"example.com",
								"www.example.com"
							],
							"serialNumber": "1234567890",
							"signatureAlgorithm": "SHA256WITHRSA",
							"startDate": "2025-08-22T11:45:19Z",
							"subject": {
								"commonName": "example.com",
								"country": "US",
								"locality": "Cambridge",
								"state": "Massachusetts"
							},
							"validation": {
								"errors": [],
								"notices": [],
								"warnings": [
									{
										"detail": "Message: Certificate validity period is above the maximum 398 days.. Name: LEAF_CERTIFICATE",
										"instance": "/error-types/certificate-validation-warning?traceId=123456789",
										"message": "Certificate validity period is above the maximum 398 days.",
										"name": "LEAF_CERTIFICATE",
										"title": "Certificate validation warning.",
										"type": "/error-types/certificate-validation-warning"
									}
								]
							}
						}
					],
					"trustChain": [],
					"trustChainPem": null,
					"validation": {
						"errors": [],
						"notices": [],
						"warnings": [
							{
								"detail": "Message: Name: MaxMinExpirationDateValidator Message: RSA certificate expiration is longer than allowed. Must expire within 398 days. Name: UNKNOWN",
								"instance": "/error-types/certificate-validation-warning?traceId=123456789",
								"message": "Name: MaxMinExpirationDateValidator Message: RSA certificate expiration is longer than allowed. Must expire within 398 days",
								"name": "UNKNOWN",
								"status": 400,
								"title": "Certificate validation warning.",
								"type": "/error-types/certificate-validation-warning"
							},
							{
								"detail": "Message: Name: TrustChainRequiredValidator Message: RSA certificate does not come with a trust chain and this is a non-standard practice. Name: UNKNOWN",
								"instance": "/error-types/certificate-validation-warning?traceId=123456789",
								"message": "Name: TrustChainRequiredValidator Message: RSA certificate does not come with a trust chain and this is a non-standard practice",
								"name": "UNKNOWN",
								"status": 400,
								"title": "Certificate validation warning.",
								"type": "/error-types/certificate-validation-warning"
							}
						]
					}
				},
				"detail": "Warnings detected in one or more of the uploaded certificates.",
				"instance": "/error-types/upload-certificate-validation-warnings?traceId=123456789",
				"status": 409,
				"title": "Validation warnings for uploaded signed certificate(s) and trust chain(s) for one or more key types.",
				"type": "/error-types/upload-certificate-validation-warnings"
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPatchCertificate, &Error{
					Type:     "/error-types/upload-certificate-validation-warnings",
					Title:    "Validation warnings for uploaded signed certificate(s) and trust chain(s) for one or more key types.",
					Detail:   "Warnings detected in one or more of the uploaded certificates.",
					Instance: "/error-types/upload-certificate-validation-warnings?traceId=123456789",
					Status:   409,
					Data: &ValidationData{
						SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
						SignedCertificates: []PEMValidation{
							{
								CertificatePEM:     "-----BEGIN CERTIFICATE-----\nexample-PEM\n-----END CERTIFICATE-----",
								CreatedBy:          nil,
								CreatedDate:        nil,
								DisplayName:        nil,
								EndDate:            ptr.To("2027-11-22T12:45:19Z"),
								Fingerprint:        nil,
								Issuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
								SerialNumber:       ptr.To("1234567890"),
								SignatureAlgorithm: ptr.To("SHA256WITHRSA"),
								StartDate:          ptr.To("2025-08-22T11:45:19Z"),
								Subject: &Subject{
									CommonName: "example.com",
									Country:    "US",
									Locality:   "Cambridge",
									State:      "Massachusetts",
								},
								Validation: &ValidationResult{
									Errors:  []ValidationDetail{},
									Notices: []ValidationDetail{},
									Warnings: []ValidationDetail{
										{
											Detail:   "Message: Certificate validity period is above the maximum 398 days.. Name: LEAF_CERTIFICATE",
											Instance: "/error-types/certificate-validation-warning?traceId=123456789",
											Message:  "Certificate validity period is above the maximum 398 days.",
											Name:     "LEAF_CERTIFICATE",
											Title:    "Certificate validation warning.",
											Type:     "/error-types/certificate-validation-warning",
										},
									},
								},
							},
						},
						TrustChain:    []PEMValidation{},
						TrustChainPEM: nil,
						Validation: &ValidationResult{
							Errors:  []ValidationDetail{},
							Notices: []ValidationDetail{},
							Warnings: []ValidationDetail{
								{
									Detail:   "Message: Name: MaxMinExpirationDateValidator Message: RSA certificate expiration is longer than allowed. Must expire within 398 days. Name: UNKNOWN",
									Instance: "/error-types/certificate-validation-warning?traceId=123456789",
									Message:  "Name: MaxMinExpirationDateValidator Message: RSA certificate expiration is longer than allowed. Must expire within 398 days",
									Name:     "UNKNOWN",
									Status:   ptr.To(400),
									Title:    "Certificate validation warning.",
									Type:     "/error-types/certificate-validation-warning",
								},
								{
									Detail:   "Message: Name: TrustChainRequiredValidator Message: RSA certificate does not come with a trust chain and this is a non-standard practice. Name: UNKNOWN",
									Instance: "/error-types/certificate-validation-warning?traceId=123456789",
									Message:  "Name: TrustChainRequiredValidator Message: RSA certificate does not come with a trust chain and this is a non-standard practice",
									Name:     "UNKNOWN",
									Status:   ptr.To(400),
									Title:    "Certificate validation warning.",
									Type:     "/error-types/certificate-validation-warning",
								},
							},
						},
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
			},
		},
		"409 OK - warnings in the response for certificate pem and trust chain": {
			params: PatchCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
			},
			expectedRequestBody: `[{"op":"add","path":"/signedCertificatePem","value":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"}]`,
			responseStatus:      409,
			responseBody: `
			{
				"data": {
					"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
					"signedCertificates": [
						{
							"certificatePem": "-----BEGIN CERTIFICATE-----\nexample-PEM\n-----END CERTIFICATE-----",
							"createdBy": null,
							"createdDate": null,
							"displayName": null,
							"endDate": "2027-11-22T12:45:19Z",
							"fingerprint": null,
							"issuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
							"sans": [
								"example.com",
								"www.example.com"
							],
							"serialNumber": "1234567890",
							"signatureAlgorithm": "SHA256WITHRSA",
							"startDate": "2025-08-22T11:45:19Z",
							"subject": {
								"commonName": "example.com",
								"country": "US",
								"locality": "Cambridge",
								"organization": null,
								"state": "Massachusetts"
							},
							"validation": {
								"errors": [],
								"notices": [],
								"warnings": [
									{
										"detail": "Message: Certificate validity period is above the maximum 398 days.. Name: LEAF_CERTIFICATE",
										"instance": "/error-types/certificate-validation-warning?traceId=123456789",
										"message": "Certificate validity period is above the maximum 398 days.",
										"name": "LEAF_CERTIFICATE",
										"title": "Certificate validation warning.",
										"type": "/error-types/certificate-validation-warning"
									}
								]
							}
						}
					],
					"trustChain": [
						{
							"certificatePem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----",
							"createdBy": null,
							"createdDate": null,
							"displayName": null,
							"endDate": "2027-11-22T12:45:19Z",
							"fingerprint": null,
							"issuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
							"sans": null,
							"serialNumber": "1234567890",
							"signatureAlgorithm": "SHA256WITHRSA",
							"startDate": "2025-08-22T11:45:19Z",
							"subject": null,
							"validation": {
								"errors": [
									{
										"detail": "Message: Certificate is not an intermediate trust chain certificate. Name: TRUST_CHAIN_CERTIFICATE",
										"instance": "/error-types/certificate-validation-failed?traceId=123456789",
										"message": "Certificate is not an intermediate trust chain certificate.",
										"name": "TRUST_CHAIN_CERTIFICATE",
										"title": "Certificate validation error.",
										"type": "/error-types/certificate-validation-failed"
									}
								],
								"notices": [],
								"warnings": []
							}
						}
					],
					"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n",
					"validation": {
						"errors": [],
						"notices": [],
						"warnings": []
					}
				},
				"detail": "Errors detected in one or more of the uploaded certificates.",
				"instance": "/error-types/upload-certificate-validation-failed?traceId=123456789",
				"status": 400,
				"title": "Validation failed for uploaded signed certificate(s) and trust chain(s) for one or more key types.",
				"type": "/error-types/upload-certificate-validation-failed"
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPatchCertificate, &Error{
					Type:     "/error-types/upload-certificate-validation-failed",
					Title:    "Validation failed for uploaded signed certificate(s) and trust chain(s) for one or more key types.",
					Detail:   "Errors detected in one or more of the uploaded certificates.",
					Instance: "/error-types/upload-certificate-validation-failed?traceId=123456789",
					Status:   409,
					Data: &ValidationData{
						SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
						SignedCertificates: []PEMValidation{
							{
								CertificatePEM:     "-----BEGIN CERTIFICATE-----\nexample-PEM\n-----END CERTIFICATE-----",
								CreatedBy:          nil,
								CreatedDate:        nil,
								DisplayName:        nil,
								EndDate:            ptr.To("2027-11-22T12:45:19Z"),
								Fingerprint:        nil,
								Issuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
								SerialNumber:       ptr.To("1234567890"),
								SignatureAlgorithm: ptr.To("SHA256WITHRSA"),
								StartDate:          ptr.To("2025-08-22T11:45:19Z"),
								Subject: &Subject{
									CommonName: "example.com",
									Country:    "US",
									Locality:   "Cambridge",
									State:      "Massachusetts",
								},
								Validation: &ValidationResult{
									Errors:  []ValidationDetail{},
									Notices: []ValidationDetail{},
									Warnings: []ValidationDetail{
										{
											Detail:   "Message: Certificate validity period is above the maximum 398 days.. Name: LEAF_CERTIFICATE",
											Instance: "/error-types/certificate-validation-warning?traceId=123456789",
											Message:  "Certificate validity period is above the maximum 398 days.",
											Name:     "LEAF_CERTIFICATE",
											Title:    "Certificate validation warning.",
											Type:     "/error-types/certificate-validation-warning",
										},
									},
								},
							},
						},
						TrustChain: []PEMValidation{
							{
								CertificatePEM:     "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----",
								CreatedBy:          nil,
								CreatedDate:        nil,
								DisplayName:        nil,
								EndDate:            ptr.To("2027-11-22T12:45:19Z"),
								Fingerprint:        nil,
								Issuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
								SerialNumber:       ptr.To("1234567890"),
								SignatureAlgorithm: ptr.To("SHA256WITHRSA"),
								StartDate:          ptr.To("2025-08-22T11:45:19Z"),
								Subject:            nil,
								Validation: &ValidationResult{
									Errors: []ValidationDetail{
										{
											Detail:   "Message: Certificate is not an intermediate trust chain certificate. Name: TRUST_CHAIN_CERTIFICATE",
											Instance: "/error-types/certificate-validation-failed?traceId=123456789",
											Message:  "Certificate is not an intermediate trust chain certificate.",
											Name:     "TRUST_CHAIN_CERTIFICATE",
											Title:    "Certificate validation error.",
											Type:     "/error-types/certificate-validation-failed",
										},
									},
									Notices:  []ValidationDetail{},
									Warnings: []ValidationDetail{},
								},
							},
						},
						TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
						Validation: &ValidationResult{
							Errors:   []ValidationDetail{},
							Notices:  []ValidationDetail{},
							Warnings: []ValidationDetail{},
						},
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
			},
		},
		"409 Conflict - name already in use": {
			params: PatchCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To("duplicate-name"),
			},
			expectedRequestBody: "[{\"op\":\"replace\",\"path\":\"/certificateName\",\"value\":\"duplicate-name\"}]",
			responseStatus:      409,
			responseBody: `
			{
				"certificateIdentifier": "certificateName",
				"certificateIdentifierValue": "duplicate-name.com1234567890",
				"detail": "Certificate with {certificateName}: {duplicate-name.com1234567890} already exists with the current account Id!",
				"instance": "/error-types/certificate-name-already-in-use?traceId=-123",
				"status": 409,
				"title": "Certificate name already in use.",
				"type": "/error-types/certificate-name-already-in-use"
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json-patch+json",
			},
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPatchCertificate, &Error{
					Type:                       "/error-types/certificate-name-already-in-use",
					Title:                      "Certificate name already in use.",
					Detail:                     "Certificate with {certificateName}: {duplicate-name.com1234567890} already exists with the current account Id!",
					Status:                     http.StatusConflict,
					Instance:                   "/error-types/certificate-name-already-in-use?traceId=-123",
					CertificateIdentifier:      "certificateName",
					CertificateIdentifierValue: "duplicate-name.com1234567890",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertificateNameInUse)
			},
		},
		"validation error - name too long - assert that error is ErrPatchCertificate and ErrStructValidation": {
			params: PatchCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To(strings.Repeat("A", 271)),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patching certificate: struct validation: CertificateName: the length must be no more than 270",
					err.Error())
				assert.ErrorIs(t, err, ErrPatchCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - name contains invalid characters": {
			params: PatchCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To("Invalid@Name"),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patching certificate: struct validation: CertificateName: the input can only contain digits (1-9), letters (a-z, A-Z), spaces, hyphens, periods, and underscores.",
					err.Error())
				assert.ErrorIs(t, err, ErrPatchCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing required parameters": {
			params: PatchCertificateRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patching certificate: struct validation: CertificateID: cannot be blank\nrequired parameters: at least one of SignedCertificatePEM or CertificateName must be provided",
					err.Error())
				assert.ErrorIs(t, err, ErrPatchCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"URL parsing error - assert that error is ErrPatchCertificate when parsing URL": {
			params: PatchCertificateRequest{
				CertificateID: "123 wrong url",
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "patching certificate: struct validation: required parameters: at least one of SignedCertificatePEM or CertificateName must be provided")
				assert.ErrorIs(t, err, ErrPatchCertificate, "want: %s; got: %s", ErrPatchCertificate, err)
				assert.ErrorIs(t, err, ErrStructValidation, "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"500 internal server error - assert that error is ErrPatchCertificate": {
			params: PatchCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To("New-Certificate-Name"),
			},
			expectedRequestBody: "[{\"op\":\"replace\",\"path\":\"/certificateName\",\"value\":\"New-Certificate-Name\"}]",
			responseStatus:      500,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error removing property hostname",
				"status": 500
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPatchCertificate, &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error removing property hostname",
					Status: http.StatusInternalServerError,
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrPatchCertificate)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPatch, r.Method)
				if len(tc.returnedHeaders) > 0 {
					for header, value := range tc.returnedHeaders {
						w.Header().Set(header, value)
					}
				}
				for k, v := range tc.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRequestBody, string(requestBody))
				w.WriteHeader(tc.responseStatus)
				_, err = w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.PatchCertificate(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestUpdateCertificate(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params              UpdateCertificateRequest
		responseStatus      int
		responseBody        string
		expectedResponse    *UpdateCertificateResponse
		expectedRequestBody string
		expectedPath        string
		expectedHeaders     map[string]string
		returnedHeaders     map[string]string
		withError           func(*testing.T, error)
	}{
		"200 OK - only rename with all allowed characters": {
			params: UpdateCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To("test 0123456789.-_"),
			},
			returnedHeaders: map[string]string{
				"Akamai-RateLimit-Limit":     "60",
				"Akamai-RateLimit-Remaining": "59",
			},
			expectedResponse: &UpdateCertificateResponse{
				Certificate: Certificate{
					AccountID:               "acc_123",
					CertificateID:           "123",
					CertificateName:         "test 0123456789.-_",
					CertificateStatus:       "CSR_READY",
					CertificateType:         "THIRD_PARTY",
					ContractID:              "A-123",
					CreatedBy:               "user",
					CreatedDate:             test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:       test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                  "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                 "2048",
					KeyType:                 "RSA",
					ModifiedBy:              "user",
					ModifiedDate:            test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                    []string{"example.com", "www.example.com"},
					SecureNetwork:           "ENHANCED_TLS",
					SignedCertificateIssuer: nil,
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			expectedRequestBody: `{"certificateName":"test 0123456789.-_"}`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "test 0123456789.-_",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		"200 OK - reset name by providing empty value": {
			params: UpdateCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To(""),
			},
			expectedResponse: &UpdateCertificateResponse{
				Certificate: Certificate{
					AccountID:               "acc_123",
					CertificateID:           "123",
					CertificateName:         "example.com20250822092651008941",
					CertificateStatus:       "CSR_READY",
					CertificateType:         "THIRD_PARTY",
					ContractID:              "A-123",
					CreatedBy:               "user",
					CreatedDate:             test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:       test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                  "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                 "2048",
					KeyType:                 "RSA",
					ModifiedBy:              "user",
					ModifiedDate:            test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                    []string{"example.com", "www.example.com"},
					SecureNetwork:           "ENHANCED_TLS",
					SignedCertificateIssuer: nil,
					Subject: &Subject{
						Country:      "US",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
						Locality:     "Cambridge",
						CommonName:   "example.com",
					},
				},
			},
			expectedRequestBody: "{\"certificateName\":\"\"}",
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "example.com20250822092651008941",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		"200 OK - upload signed certificate PEM only": {
			params: UpdateCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
			},
			expectedResponse: &UpdateCertificateResponse{
				Certificate: Certificate{
					AccountID:                           "acc_123",
					CertificateID:                       "123",
					CertificateName:                     "Certificate-name-rename",
					CertificateStatus:                   "READY_FOR_USE",
					CertificateType:                     "THIRD_PARTY",
					ContractID:                          "A-123",
					CreatedBy:                           "user",
					CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                             "2048",
					KeyType:                             "RSA",
					ModifiedBy:                          "user",
					ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                                []string{"example.com", "www.example.com"},
					SecureNetwork:                       "ENHANCED_TLS",
					SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
					SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
					SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
					SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
					},
					TrustChainPEM: nil,
				},
			},
			expectedRequestBody: `{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"}`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "Certificate-name-rename",
				"certificateStatus": "READY_FOR_USE",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
				"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
				"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
				"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				"signedCertificateSHA256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
				"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		"200 OK - upload signed certificate PEM with trust chain": {
			params: UpdateCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				TrustChainPEM:        "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n",
			},
			expectedResponse: &UpdateCertificateResponse{
				Certificate: Certificate{
					AccountID:                           "acc_123",
					CertificateID:                       "123",
					CertificateName:                     "Certificate-name-rename",
					CertificateStatus:                   "READY_FOR_USE",
					CertificateType:                     "THIRD_PARTY",
					ContractID:                          "A-123",
					CreatedBy:                           "user",
					CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                             "2048",
					KeyType:                             "RSA",
					ModifiedBy:                          "user",
					ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                                []string{"example.com", "www.example.com"},
					SecureNetwork:                       "ENHANCED_TLS",
					SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
					SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
					SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
					SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
					},
					TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
				},
			},
			expectedRequestBody: `{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n","trustChainPem":"-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"}`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "Certificate-name-rename",
				"certificateStatus": "READY_FOR_USE",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
				"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
				"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
				"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				"signedCertificateSHA256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
				"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		"200 OK - upload signed certificate with AcknowledgeWarnings query param": {
			params: UpdateCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				AcknowledgeWarnings:  true,
			},
			expectedResponse: &UpdateCertificateResponse{
				Certificate: Certificate{
					AccountID:                           "acc_123",
					CertificateID:                       "123",
					CertificateName:                     "Certificate-name-rename",
					CertificateStatus:                   "READY_FOR_USE",
					CertificateType:                     "THIRD_PARTY",
					ContractID:                          "A-123",
					CreatedBy:                           "user",
					CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                             "2048",
					KeyType:                             "RSA",
					ModifiedBy:                          "user",
					ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                                []string{"example.com", "www.example.com"},
					SecureNetwork:                       "ENHANCED_TLS",
					SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
					SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
					SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
					SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
					},
					TrustChainPEM: nil,
				},
			},
			expectedRequestBody: `{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"}`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "Certificate-name-rename",
				"certificateStatus": "READY_FOR_USE",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
				"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
				"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
				"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				"signedCertificateSHA256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
				"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
			expectedPath: "/ccm/v1/certificates/123?acknowledgeWarnings=true",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		"200 OK - upload signed certificate and trust chain PEM with AcknowledgeWarnings query param and rename certificate": {
			params: UpdateCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				TrustChainPEM:        "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n",
				CertificateName:      ptr.To("Certificate-name-rename"),
				AcknowledgeWarnings:  true,
			},
			expectedResponse: &UpdateCertificateResponse{
				Certificate: Certificate{
					AccountID:                           "acc_123",
					CertificateID:                       "123",
					CertificateName:                     "Certificate-name-rename",
					CertificateStatus:                   "READY_FOR_USE",
					CertificateType:                     "THIRD_PARTY",
					ContractID:                          "A-123",
					CreatedBy:                           "user",
					CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
					CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
					KeySize:                             "2048",
					KeyType:                             "RSA",
					ModifiedBy:                          "user",
					ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:32.607358Z"),
					SANs:                                []string{"example.com", "www.example.com"},
					SecureNetwork:                       "ENHANCED_TLS",
					SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
					SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
					SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
					SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
					Subject: &Subject{
						CommonName:   "example.com",
						Country:      "US",
						Locality:     "Cambridge",
						Organization: "ExampleOrg",
						State:        "Massachusetts",
					},
					TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
				},
			},
			expectedRequestBody: `{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n","trustChainPem":"-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n","certificateName":"Certificate-name-rename"}`,
			responseStatus:      200,
			responseBody: `
			{
				"accountId": "acc_123",
				"certificateId": "123",
				"certificateName": "Certificate-name-rename",
				"certificateStatus": "READY_FOR_USE",
				"certificateType": "THIRD_PARTY",
				"contractId": "A-123",
				"createdBy": "user",
				"createdDate": "2025-08-22T09:01:32.607357Z",
				"csrExpirationDate": "2026-10-24T09:01:34Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "user",
				"modifiedDate": "2025-08-22T09:01:32.607358Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
				"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
				"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
				"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
				"signedCertificateSHA256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
				"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"
			}`,
			expectedPath: "/ccm/v1/certificates/123?acknowledgeWarnings=true",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		"409 OK - warnings in the response for certificate pem": {
			params: UpdateCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
			},
			expectedRequestBody: `{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"}`,
			responseStatus:      409,
			responseBody: `
			{
				"data": {
					"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
					"signedCertificates": [
						{
							"certificatePem": "-----BEGIN CERTIFICATE-----\nexample-PEM\n-----END CERTIFICATE-----",
							"createdBy": null,
							"createdDate": null,
							"displayName": null,
							"endDate": "2027-11-22T12:45:19Z",
							"fingerprint": null,
							"issuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
							"sans": [
								"example.com",
								"www.example.com"
							],
							"serialNumber": "1234567890",
							"signatureAlgorithm": "SHA256WITHRSA",
							"startDate": "2025-08-22T11:45:19Z",
							"subject": {
								"commonName": "example.com",
								"country": "US",
								"locality": "Cambridge",
								"state": "Massachusetts"
							},
							"validation": {
								"errors": [],
								"notices": [],
								"warnings": [
									{
										"detail": "Message: Certificate validity period is above the maximum 398 days.. Name: LEAF_CERTIFICATE",
										"instance": "/error-types/certificate-validation-warning?traceId=123456789",
										"message": "Certificate validity period is above the maximum 398 days.",
										"name": "LEAF_CERTIFICATE",
										"title": "Certificate validation warning.",
										"type": "/error-types/certificate-validation-warning"
									}
								]
							}
						}
					],
					"trustChain": [],
					"trustChainPem": null,
					"validation": {
						"errors": [],
						"notices": [],
						"warnings": [
							{
								"detail": "Message: Name: MaxMinExpirationDateValidator Message: RSA certificate expiration is longer than allowed. Must expire within 398 days. Name: UNKNOWN",
								"instance": "/error-types/certificate-validation-warning?traceId=123456789",
								"message": "Name: MaxMinExpirationDateValidator Message: RSA certificate expiration is longer than allowed. Must expire within 398 days",
								"name": "UNKNOWN",
								"status": 400,
								"title": "Certificate validation warning.",
								"type": "/error-types/certificate-validation-warning"
							},
							{
								"detail": "Message: Name: TrustChainRequiredValidator Message: RSA certificate does not come with a trust chain and this is a non-standard practice. Name: UNKNOWN",
								"instance": "/error-types/certificate-validation-warning?traceId=123456789",
								"message": "Name: TrustChainRequiredValidator Message: RSA certificate does not come with a trust chain and this is a non-standard practice",
								"name": "UNKNOWN",
								"status": 400,
								"title": "Certificate validation warning.",
								"type": "/error-types/certificate-validation-warning"
							}
						]
					}
				},
				"detail": "Warnings detected in one or more of the uploaded certificates.",
				"instance": "/error-types/upload-certificate-validation-warnings?traceId=123456789",
				"status": 409,
				"title": "Validation warnings for uploaded signed certificate(s) and trust chain(s) for one or more key types.",
				"type": "/error-types/upload-certificate-validation-warnings"
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUpdateCertificate, &Error{
					Type:     "/error-types/upload-certificate-validation-warnings",
					Title:    "Validation warnings for uploaded signed certificate(s) and trust chain(s) for one or more key types.",
					Detail:   "Warnings detected in one or more of the uploaded certificates.",
					Instance: "/error-types/upload-certificate-validation-warnings?traceId=123456789",
					Status:   409,
					Data: &ValidationData{
						SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
						SignedCertificates: []PEMValidation{
							{
								CertificatePEM:     "-----BEGIN CERTIFICATE-----\nexample-PEM\n-----END CERTIFICATE-----",
								CreatedBy:          nil,
								CreatedDate:        nil,
								DisplayName:        nil,
								EndDate:            ptr.To("2027-11-22T12:45:19Z"),
								Fingerprint:        nil,
								Issuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
								SerialNumber:       ptr.To("1234567890"),
								SignatureAlgorithm: ptr.To("SHA256WITHRSA"),
								StartDate:          ptr.To("2025-08-22T11:45:19Z"),
								Subject: &Subject{
									CommonName: "example.com",
									Country:    "US",
									Locality:   "Cambridge",
									State:      "Massachusetts",
								},
								Validation: &ValidationResult{
									Errors:  []ValidationDetail{},
									Notices: []ValidationDetail{},
									Warnings: []ValidationDetail{
										{
											Detail:   "Message: Certificate validity period is above the maximum 398 days.. Name: LEAF_CERTIFICATE",
											Instance: "/error-types/certificate-validation-warning?traceId=123456789",
											Message:  "Certificate validity period is above the maximum 398 days.",
											Name:     "LEAF_CERTIFICATE",
											Title:    "Certificate validation warning.",
											Type:     "/error-types/certificate-validation-warning",
										},
									},
								},
							},
						},
						TrustChain:    []PEMValidation{},
						TrustChainPEM: nil,
						Validation: &ValidationResult{
							Errors:  []ValidationDetail{},
							Notices: []ValidationDetail{},
							Warnings: []ValidationDetail{
								{
									Detail:   "Message: Name: MaxMinExpirationDateValidator Message: RSA certificate expiration is longer than allowed. Must expire within 398 days. Name: UNKNOWN",
									Instance: "/error-types/certificate-validation-warning?traceId=123456789",
									Message:  "Name: MaxMinExpirationDateValidator Message: RSA certificate expiration is longer than allowed. Must expire within 398 days",
									Name:     "UNKNOWN",
									Status:   ptr.To(400),
									Title:    "Certificate validation warning.",
									Type:     "/error-types/certificate-validation-warning",
								},
								{
									Detail:   "Message: Name: TrustChainRequiredValidator Message: RSA certificate does not come with a trust chain and this is a non-standard practice. Name: UNKNOWN",
									Instance: "/error-types/certificate-validation-warning?traceId=123456789",
									Message:  "Name: TrustChainRequiredValidator Message: RSA certificate does not come with a trust chain and this is a non-standard practice",
									Name:     "UNKNOWN",
									Status:   ptr.To(400),
									Title:    "Certificate validation warning.",
									Type:     "/error-types/certificate-validation-warning",
								},
							},
						},
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
			},
		},
		"409 OK - warnings in the response for certificate pem and trust chain": {
			params: UpdateCertificateRequest{
				CertificateID:        "123",
				SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
			},
			expectedRequestBody: `{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"}`,
			responseStatus:      409,
			responseBody: `
			{
				"data": {
					"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
					"signedCertificates": [
						{
							"certificatePem": "-----BEGIN CERTIFICATE-----\nexample-PEM\n-----END CERTIFICATE-----",
							"createdBy": null,
							"createdDate": null,
							"displayName": null,
							"endDate": "2027-11-22T12:45:19Z",
							"fingerprint": null,
							"issuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
							"sans": [
								"example.com",
								"www.example.com"
							],
							"serialNumber": "1234567890",
							"signatureAlgorithm": "SHA256WITHRSA",
							"startDate": "2025-08-22T11:45:19Z",
							"subject": {
								"commonName": "example.com",
								"country": "US",
								"locality": "Cambridge",
								"organization": null,
								"state": "Massachusetts"
							},
							"validation": {
								"errors": [],
								"notices": [],
								"warnings": [
									{
										"detail": "Message: Certificate validity period is above the maximum 398 days.. Name: LEAF_CERTIFICATE",
										"instance": "/error-types/certificate-validation-warning?traceId=123456789",
										"message": "Certificate validity period is above the maximum 398 days.",
										"name": "LEAF_CERTIFICATE",
										"title": "Certificate validation warning.",
										"type": "/error-types/certificate-validation-warning"
									}
								]
							}
						}
					],
					"trustChain": [
						{
							"certificatePem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----",
							"createdBy": null,
							"createdDate": null,
							"displayName": null,
							"endDate": "2027-11-22T12:45:19Z",
							"fingerprint": null,
							"issuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
							"sans": null,
							"serialNumber": "1234567890",
							"signatureAlgorithm": "SHA256WITHRSA",
							"startDate": "2025-08-22T11:45:19Z",
							"subject": null,
							"validation": {
								"errors": [
									{
										"detail": "Message: Certificate is not an intermediate trust chain certificate. Name: TRUST_CHAIN_CERTIFICATE",
										"instance": "/error-types/certificate-validation-failed?traceId=123456789",
										"message": "Certificate is not an intermediate trust chain certificate.",
										"name": "TRUST_CHAIN_CERTIFICATE",
										"title": "Certificate validation error.",
										"type": "/error-types/certificate-validation-failed"
									}
								],
								"notices": [],
								"warnings": []
							}
						}
					],
					"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n",
					"validation": {
						"errors": [],
						"notices": [],
						"warnings": []
					}
				},
				"detail": "Errors detected in one or more of the uploaded certificates.",
				"instance": "/error-types/upload-certificate-validation-failed?traceId=123456789",
				"status": 400,
				"title": "Validation failed for uploaded signed certificate(s) and trust chain(s) for one or more key types.",
				"type": "/error-types/upload-certificate-validation-failed"
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUpdateCertificate, &Error{
					Type:     "/error-types/upload-certificate-validation-failed",
					Title:    "Validation failed for uploaded signed certificate(s) and trust chain(s) for one or more key types.",
					Detail:   "Errors detected in one or more of the uploaded certificates.",
					Instance: "/error-types/upload-certificate-validation-failed?traceId=123456789",
					Status:   409,
					Data: &ValidationData{
						SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
						SignedCertificates: []PEMValidation{
							{
								CertificatePEM:     "-----BEGIN CERTIFICATE-----\nexample-PEM\n-----END CERTIFICATE-----",
								CreatedBy:          nil,
								CreatedDate:        nil,
								DisplayName:        nil,
								EndDate:            ptr.To("2027-11-22T12:45:19Z"),
								Fingerprint:        nil,
								Issuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
								SerialNumber:       ptr.To("1234567890"),
								SignatureAlgorithm: ptr.To("SHA256WITHRSA"),
								StartDate:          ptr.To("2025-08-22T11:45:19Z"),
								Subject: &Subject{
									CommonName: "example.com",
									Country:    "US",
									Locality:   "Cambridge",
									State:      "Massachusetts",
								},
								Validation: &ValidationResult{
									Errors:  []ValidationDetail{},
									Notices: []ValidationDetail{},
									Warnings: []ValidationDetail{
										{
											Detail:   "Message: Certificate validity period is above the maximum 398 days.. Name: LEAF_CERTIFICATE",
											Instance: "/error-types/certificate-validation-warning?traceId=123456789",
											Message:  "Certificate validity period is above the maximum 398 days.",
											Name:     "LEAF_CERTIFICATE",
											Title:    "Certificate validation warning.",
											Type:     "/error-types/certificate-validation-warning",
										},
									},
								},
							},
						},
						TrustChain: []PEMValidation{
							{
								CertificatePEM:     "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----",
								CreatedBy:          nil,
								CreatedDate:        nil,
								DisplayName:        nil,
								EndDate:            ptr.To("2027-11-22T12:45:19Z"),
								Fingerprint:        nil,
								Issuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
								SerialNumber:       ptr.To("1234567890"),
								SignatureAlgorithm: ptr.To("SHA256WITHRSA"),
								StartDate:          ptr.To("2025-08-22T11:45:19Z"),
								Subject:            nil,
								Validation: &ValidationResult{
									Errors: []ValidationDetail{
										{
											Detail:   "Message: Certificate is not an intermediate trust chain certificate. Name: TRUST_CHAIN_CERTIFICATE",
											Instance: "/error-types/certificate-validation-failed?traceId=123456789",
											Message:  "Certificate is not an intermediate trust chain certificate.",
											Name:     "TRUST_CHAIN_CERTIFICATE",
											Title:    "Certificate validation error.",
											Type:     "/error-types/certificate-validation-failed",
										},
									},
									Notices:  []ValidationDetail{},
									Warnings: []ValidationDetail{},
								},
							},
						},
						TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
						Validation: &ValidationResult{
							Errors:   []ValidationDetail{},
							Notices:  []ValidationDetail{},
							Warnings: []ValidationDetail{},
						},
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
			},
		},
		"409 Conflict - name already in use": {
			params: UpdateCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To("duplicate-name"),
			},
			expectedRequestBody: `{"certificateName":"duplicate-name"}`,
			responseStatus:      409,
			responseBody: `
			{
				"certificateIdentifier": "certificateName",
				"certificateIdentifierValue": "duplicate-name.com1234567890",
				"detail": "Certificate with {certificateName}: {duplicate-name.com1234567890} already exists with the current account Id!",
				"instance": "/error-types/certificate-name-already-in-use?traceId=-123",
				"status": 409,
				"title": "Certificate name already in use.",
				"type": "/error-types/certificate-name-already-in-use"
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUpdateCertificate, &Error{
					Type:                       "/error-types/certificate-name-already-in-use",
					Title:                      "Certificate name already in use.",
					Detail:                     "Certificate with {certificateName}: {duplicate-name.com1234567890} already exists with the current account Id!",
					Status:                     http.StatusConflict,
					Instance:                   "/error-types/certificate-name-already-in-use?traceId=-123",
					CertificateIdentifier:      "certificateName",
					CertificateIdentifierValue: "duplicate-name.com1234567890",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertificateNameInUse)
			},
		},
		"validation error - name too long - assert that error is ErrUpdateCertificate and ErrStructValidation": {
			params: UpdateCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To(strings.Repeat("A", 271)),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "updating certificate: struct validation: CertificateName: the length must be no more than 270",
					err.Error())
				assert.ErrorIs(t, err, ErrUpdateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - name contains invalid characters": {
			params: UpdateCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To("Invalid@Name"),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "updating certificate: struct validation: CertificateName: the input can only contain digits (1-9), letters (a-z, A-Z), spaces, hyphens, periods, and underscores.",
					err.Error())
				assert.ErrorIs(t, err, ErrUpdateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing required parameters": {
			params: UpdateCertificateRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "updating certificate: struct validation: CertificateID: cannot be blank", err.Error())
				assert.ErrorIs(t, err, ErrUpdateCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"500 internal server error - assert that error is ErrUpdateCertificate": {
			params: UpdateCertificateRequest{
				CertificateID:   "123",
				CertificateName: ptr.To("New-Certificate-Name"),
			},
			expectedRequestBody: `{"certificateName":"New-Certificate-Name"}`,
			responseStatus:      500,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error removing property hostname",
				"status": 500
			}`,
			expectedPath: "/ccm/v1/certificates/123",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUpdateCertificate, &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error removing property hostname",
					Status: http.StatusInternalServerError,
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrUpdateCertificate)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if len(tc.returnedHeaders) > 0 {
					for header, value := range tc.returnedHeaders {
						w.Header().Set(header, value)
					}
				}
				for k, v := range tc.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRequestBody, string(requestBody))
				w.WriteHeader(tc.responseStatus)
				_, err = w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateCertificate(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListCertificates(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params           ListCertificatesRequest
		expectedResponse *ListCertificatesResponse
		responseStatus   int
		responseBody     string
		expectedPath     string
		returnedHeaders  map[string]string
		withError        func(t *testing.T, err error)
	}{
		"200 OK - list certificates with no params": {
			params: ListCertificatesRequest{},
			returnedHeaders: map[string]string{
				"Akamai-RateLimit-Limit":     "60",
				"Akamai-RateLimit-Remaining": "59",
			},
			expectedResponse: &ListCertificatesResponse{
				Certificates: []Certificate{
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert1_1234",
						CertificateName:                     "Test Certificate1",
						CertificateStatus:                   "ACTIVE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "2048",
						KeyType:                             "RSA",
						SANs:                                []string{"test-example.com", "www.test-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName: "test-example.com",
							Country:    "US",
							Locality:   "Cambridge",
							State:      "Massachusetts",
						},
						TrustChainPEM: nil,
					},
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert2_1234",
						CertificateName:                     "Test Certificate2",
						CertificateStatus:                   "READY_FOR_USE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-09-15T10:20:30.123456Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "P-256",
						KeyType:                             "ECDSA",
						SANs:                                []string{"test2-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("O=test organization NY,L=New York,ST=NY,C=US"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName:   "test2-example.com",
							Country:      "US",
							Locality:     "New York",
							Organization: "test organization NY",
							State:        "New York",
						},
						TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
					},
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert3_1234",
						CertificateName:                     "Test Certificate3",
						CertificateStatus:                   "READY_FOR_USE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-10-05T11:30:45.654321Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "2048",
						KeyType:                             "RSA",
						SANs:                                []string{"test3-example.com", "www.test3-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName: "test3-example.com",
							Country:    "US",
							Locality:   "San Francisco",
							State:      "California",
						},
						TrustChainPEM: nil,
					},
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert4_1234",
						CertificateName:                     "Test Certificate4",
						CertificateStatus:                   "READY_FOR_USE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2023-10-05T11:30:45.654321Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2024-10-24T09:01:34Z"),
						KeySize:                             "P-256",
						KeyType:                             "ECDSA",
						SANs:                                []string{"test4-example.com", "www.test4-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("O=test organization LA,L=Los Angeles,ST=California,C=US"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName:   "test4-example.com",
							Country:      "US",
							Locality:     "Los Angeles",
							Organization: "test organization LA",
							State:        "California",
						},
						TrustChainPEM: nil,
					},
				},
				Links: Links{
					Self:     "/ccm/v1/certificates?page=1&pageSize=10",
					Next:     nil,
					Previous: nil,
				},
				Metadata: ListMetadata{
					TotalItems: 4,
					TotalPages: 1,
				},
				RateLimits: RateLimitsMetadata{
					Limit:     ptr.To(int64(60)),
					Remaining: ptr.To(int64(59)),
				},
			},
			responseStatus: 200,
			responseBody: `
			{
				"certificates": [
					{
						"accountId": "test-account-id",
						"certificateId": "cert1_1234",
						"certificateName": "Test Certificate1",
						"certificateStatus": "ACTIVE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-08-22T09:01:32.607357Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "2048",
						"keyType": "RSA",
						"sans": [
							"test-example.com",
							"www.test-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test-example.com",
							"country": "US",
							"locality": "Cambridge",
							"state": "Massachusetts"
						},
						"trustChainPem": null
					},
					{
						"accountId": "test-account-id",
						"certificateId": "cert2_1234",
						"certificateName": "Test Certificate2",
						"certificateStatus": "READY_FOR_USE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-09-15T10:20:30.123456Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "P-256",
						"keyType": "ECDSA",
						"sans": [
							"test2-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "O=test organization NY,L=New York,ST=NY,C=US",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test2-example.com",
							"country": "US",
							"locality": "New York",
							"organization": "test organization NY",
							"state": "New York"
						},
						"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"
					},
					{
						"accountId": "test-account-id",
						"certificateId": "cert3_1234",
						"certificateName": "Test Certificate3",
						"certificateStatus": "READY_FOR_USE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-10-05T11:30:45.654321Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "2048",
						"keyType": "RSA",
						"sans": [
							"test3-example.com",
							"www.test3-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test3-example.com",
							"country": "US",
							"locality": "San Francisco",
							"state": "California"
						},
						"trustChainPem": null
					},
					{
						"accountId": "test-account-id",
						"certificateId": "cert4_1234",
						"certificateName": "Test Certificate4",
						"certificateStatus": "READY_FOR_USE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2023-10-05T11:30:45.654321Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2024-10-24T09:01:34Z",
						"keySize": "P-256",
						"keyType": "ECDSA",
						"sans": [
							"test4-example.com",
							"www.test4-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "O=test organization LA,L=Los Angeles,ST=California,C=US",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test4-example.com",
							"country": "US",
							"locality": "Los Angeles",
							"organization": "test organization LA",
							"state": "California"
						},
						"trustChainPem": null
					}
				],
				"links": {
					"self": "/ccm/v1/certificates?page=1&pageSize=10",
					"next": null,
					"previous": null
				},
				"metadata": {
					"totalItems": 4,
					"totalPages": 1
				}
			}`,
			expectedPath: "/ccm/v1/certificates",
		},
		"200 OK - with domain filtering": {
			params: ListCertificatesRequest{
				ContractID: "A-123",
				GroupID:    "1234",
				Domain:     "test3-example.com",
			},
			expectedResponse: &ListCertificatesResponse{
				Certificates: []Certificate{
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert3_1234",
						CertificateName:                     "Test Certificate3",
						CertificateStatus:                   "CSR_READY",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-10-05T11:30:45.654321Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "2048",
						KeyType:                             "RSA",
						SANs:                                []string{"test3-example.com", "www.test3-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             nil,
						SignedCertificateNotValidAfterDate:  nil,
						SignedCertificateNotValidBeforeDate: nil,
						SignedCertificatePEM:                nil,
						SignedCertificateSHA256Fingerprint:  nil,
						SignedCertificateSerialNumber:       nil,
						Subject: &Subject{
							CommonName: "test3-example.com",
							Country:    "US",
							Locality:   "San Francisco",
							State:      "California",
						},
						TrustChainPEM: nil,
					},
				},
				Links: Links{
					Self:     "/ccm/v1/certificates?page=1&pageSize=10",
					Next:     nil,
					Previous: nil,
				},
				Metadata: ListMetadata{
					TotalItems: 1,
					TotalPages: 1,
				},
			},
			responseStatus: 200,
			responseBody: `
			{
				"certificates": [
					{
						"accountId": "test-account-id",
						"certificateId": "cert3_1234",
						"certificateName": "Test Certificate3",
						"certificateStatus": "CSR_READY",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-10-05T11:30:45.654321Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "2048",
						"keyType": "RSA",
						"sans": [
							"test3-example.com",
							"www.test3-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"subject": {
							"commonName": "test3-example.com",
							"country": "US",
							"locality": "San Francisco",
							"state": "California"
						},
						"trustChainPem": null
					}
				],
				"links": {
					"self": "/ccm/v1/certificates?page=1&pageSize=10",
					"next": null,
					"previous": null
				},
				"metadata": {
					"totalItems": 1,
					"totalPages": 1
				}
			}`,
			expectedPath: "/ccm/v1/certificates?contractId=A-123&domain=test3-example.com&groupId=1234",
		},
		"200 OK - with status filtering": {
			params: ListCertificatesRequest{
				CertificateStatus: []CertificateStatus{StatusActive, StatusReadyForUse},
			},
			expectedResponse: &ListCertificatesResponse{
				Certificates: []Certificate{
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert1_1234",
						CertificateName:                     "Test Certificate1",
						CertificateStatus:                   "ACTIVE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "2048",
						KeyType:                             "RSA",
						SANs:                                []string{"test-example.com", "www.test-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName: "test-example.com",
							Country:    "US",
							Locality:   "Cambridge",
							State:      "Massachusetts",
						},
						TrustChainPEM: nil,
					},
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert3_1234",
						CertificateName:                     "Test Certificate3",
						CertificateStatus:                   "READY_FOR_USE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-10-05T11:30:45.654321Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-10-05T11:33:45.654321Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "2048",
						KeyType:                             "RSA",
						SANs:                                []string{"test3-example.com", "www.test3-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName: "test3-example.com",
							Country:    "US",
							Locality:   "San Francisco",
							State:      "California",
						},
						TrustChainPEM: nil,
					},
				},
				Links: Links{
					Self:     "/ccm/v1/certificates?page=1&pageSize=10",
					Next:     nil,
					Previous: nil,
				},
				Metadata: ListMetadata{
					TotalItems: 3,
					TotalPages: 1,
				},
			},
			responseStatus: 200,
			responseBody: `
			{
				"certificates": [
					{
						"accountId": "test-account-id",
						"certificateId": "cert1_1234",
						"certificateName": "Test Certificate1",
						"certificateStatus": "ACTIVE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-08-22T09:01:32.607357Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "2048",
						"keyType": "RSA",
						"sans": [
							"test-example.com",
							"www.test-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test-example.com",
							"country": "US",
							"locality": "Cambridge",
							"state": "Massachusetts"
						},
						"trustChainPem": null
					},
					{
						"accountId": "test-account-id",
						"certificateId": "cert3_1234",
						"certificateName": "Test Certificate3",
						"certificateStatus": "READY_FOR_USE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-10-05T11:30:45.654321Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-10-05T11:33:45.654321Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "2048",
						"keyType": "RSA",
						"sans": [
							"test3-example.com",
							"www.test3-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test3-example.com",
							"country": "US",
							"locality": "San Francisco",
							"state": "California"
						},
						"trustChainPem": null
					}
				],
				"links": {
					"self": "/ccm/v1/certificates?page=1&pageSize=10",
					"next": null,
					"previous": null
				},
				"metadata": {
					"totalItems": 3,
					"totalPages": 1
				}
			}`,
			expectedPath: "/ccm/v1/certificates?certificateStatus=ACTIVE%2CREADY_FOR_USE",
		},
		"200 OK - with certificate name filtering": {
			params: ListCertificatesRequest{
				CertificateName: "Test Certificate1",
			},
			expectedResponse: &ListCertificatesResponse{
				Certificates: []Certificate{
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert1_1234",
						CertificateName:                     "Test Certificate1",
						CertificateStatus:                   "ACTIVE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "2048",
						KeyType:                             "RSA",
						SANs:                                []string{"test-example.com", "www.test-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName: "test-example.com",
							Country:    "US",
							Locality:   "Cambridge",
							State:      "Massachusetts",
						},
						TrustChainPEM: nil,
					},
				},
				Links: Links{
					Self:     "/ccm/v1/certificates?page=1&pageSize=10",
					Next:     nil,
					Previous: nil,
				},
				Metadata: ListMetadata{
					TotalItems: 1,
					TotalPages: 1,
				},
			},
			responseStatus: 200,
			responseBody: `
			{
				"certificates": [
					{
						"accountId": "test-account-id",
						"certificateId": "cert1_1234",
						"certificateName": "Test Certificate1",
						"certificateStatus": "ACTIVE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-08-22T09:01:32.607357Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "2048",
						"keyType": "RSA",
						"sans": [
							"test-example.com",
							"www.test-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test-example.com",
							"country": "US",
							"locality": "Cambridge",
							"state": "Massachusetts"
						},
						"trustChainPem": null
					}
				],
				"links": {
					"self": "/ccm/v1/certificates?page=1&pageSize=10",
					"next": null,
					"previous": null
				},
				"metadata": {
					"totalItems": 1,
					"totalPages": 1
				}
			}`,
			expectedPath: "/ccm/v1/certificates?certificateName=Test+Certificate1",
		},
		"200 OK - with includeCertificateMaterials set to true and issuer filtering": {
			params: ListCertificatesRequest{
				IncludeCertificateMaterials: true,
				Issuer:                      "test organization",
			},
			expectedResponse: &ListCertificatesResponse{
				Certificates: []Certificate{
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert2_1234",
						CertificateName:                     "Test Certificate2",
						CertificateStatus:                   "READY_FOR_USE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-09-15T10:20:30.123456Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "P-256",
						KeyType:                             "ECDSA",
						SANs:                                []string{"test2-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("O=test organization NY,L=New York,ST=NY,C=US"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
						SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
						SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
						Subject: &Subject{
							CommonName:   "test2-example.com",
							Country:      "US",
							Locality:     "New York",
							Organization: "test organization NY",
							State:        "New York",
						},
						TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
					},
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert4_1234",
						CertificateName:                     "Test Certificate4",
						CertificateStatus:                   "READY_FOR_USE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2023-10-05T11:30:45.654321Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2024-10-24T09:01:34Z"),
						CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
						KeySize:                             "P-256",
						KeyType:                             "ECDSA",
						SANs:                                []string{"test4-example.com", "www.test4-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("O=test organization LA,L=Los Angeles,ST=California,C=US"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n"),
						SignedCertificateSHA256Fingerprint:  ptr.To("4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98"),
						SignedCertificateSerialNumber:       ptr.To("a2:84:7d:dc:97:f1"),
						Subject: &Subject{
							CommonName:   "test4-example.com",
							Country:      "US",
							Locality:     "Los Angeles",
							Organization: "test organization LA",
							State:        "California",
						},
						TrustChainPEM: nil,
					},
				},
				Links: Links{
					Self:     "/ccm/v1/certificates?page=1&pageSize=10",
					Next:     nil,
					Previous: nil,
				},
				Metadata: ListMetadata{
					TotalItems: 3,
					TotalPages: 1,
				},
			},
			responseStatus: 200,
			responseBody: `
			{
				"certificates": [
					{
						"accountId": "test-account-id",	
						"certificateId": "cert2_1234",
						"certificateName": "Test Certificate2",
						"certificateStatus": "READY_FOR_USE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-09-15T10:20:30.123456Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "P-256",
						"keyType": "ECDSA",
						"sans": [
							"test2-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "O=test organization NY,L=New York,ST=NY,C=US",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
						"signedCertificateSha256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
						"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
						"subject": {
							"commonName": "test2-example.com",
							"country": "US",
							"locality": "New York",
							"organization": "test organization NY",
							"state": "New York"
						},
						"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"
					},
					{
						"accountId": "test-account-id",
						"certificateId": "cert4_1234",
						"certificateName": "Test Certificate4",
						"certificateStatus": "READY_FOR_USE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2023-10-05T11:30:45.654321Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2024-10-24T09:01:34Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
						"keySize": "P-256",
						"keyType": "ECDSA",
						"sans": [
							"test4-example.com",
							"www.test4-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "O=test organization LA,L=Los Angeles,ST=California,C=US",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nexample-signed-PEM\n-----END CERTIFICATE-----\n",
						"signedCertificateSha256Fingerprint": "4E:69:28:A1:CE:F1:E4:97:CE:39:FE:12:98",
						"signedCertificateSerialNumber": "a2:84:7d:dc:97:f1",
						"subject": {
							"commonName": "test4-example.com",
							"country": "US",
							"locality": "Los Angeles",
							"organization": "test organization LA",
							"state": "California"
						},
						"trustChainPem": null
					}
				],
				"links": {
					"self": "/ccm/v1/certificates?page=1&pageSize=10",
					"next": null,
					"previous": null
				},
				"metadata": {
					"totalItems": 3,
					"totalPages": 1
				}
			}`,
			expectedPath: "/ccm/v1/certificates?includeCertificateMaterials=true&issuer=test+organization",
		},
		"200 OK - with expiringInDays set to less than 0 to find and return expired certificates": {
			params: ListCertificatesRequest{
				ExpiringInDays: ptr.To(int64(-1)),
			},
			expectedResponse: &ListCertificatesResponse{
				Certificates: []Certificate{
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert4_1234",
						CertificateName:                     "Test Certificate4",
						CertificateStatus:                   "READY_FOR_USE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2023-10-05T11:30:45.654321Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2024-10-24T09:01:34Z"),
						CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
						KeySize:                             "P-256",
						KeyType:                             "ECDSA",
						SANs:                                []string{"test4-example.com", "www.test4-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("O=test organization LA,L=Los Angeles,ST=California,C=US"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName:   "test4-example.com",
							Country:      "US",
							Locality:     "Los Angeles",
							Organization: "test organization LA",
							State:        "California",
						},
						TrustChainPEM: nil,
					},
				},
				Links: Links{
					Self:     "/ccm/v1/certificates?page=1&pageSize=10",
					Next:     nil,
					Previous: nil,
				},
				Metadata: ListMetadata{
					TotalItems: 1,
					TotalPages: 1,
				},
			},
			responseStatus: 200,
			responseBody: `
			{
				"certificates": [
					{
						"accountId": "test-account-id",
						"certificateId": "cert4_1234",
						"certificateName": "Test Certificate4",
						"certificateStatus": "READY_FOR_USE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2023-10-05T11:30:45.654321Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2024-10-24T09:01:34Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
						"keySize": "P-256",
						"keyType": "ECDSA",
						"sans": [
							"test4-example.com",
							"www.test4-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "O=test organization LA,L=Los Angeles,ST=California,C=US",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test4-example.com",
							"country": "US",
							"locality": "Los Angeles",
							"organization": "test organization LA",
							"state": "California"
						},
						"trustChainPem": null
					}
				],
				"links": {
					"self": "/ccm/v1/certificates?page=1&pageSize=10",
					"next": null,
					"previous": null
				},
				"metadata": {
					"totalItems": 1,
					"totalPages": 1
				}
			}`,
			expectedPath: "/ccm/v1/certificates?expiringInDays=-1",
		},
		"200 OK - with pagination and sorting by certificate name": {
			params: ListCertificatesRequest{
				PageSize: 2,
				Page:     2,
				Sort:     "-certificateName",
			},
			expectedResponse: &ListCertificatesResponse{
				Certificates: []Certificate{
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert2_1234",
						CertificateName:                     "Test Certificate2",
						CertificateStatus:                   "READY_FOR_USE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-09-15T10:20:30.123456Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-09-15T10:20:33.123456Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "P-256",
						KeyType:                             "ECDSA",
						SANs:                                []string{"test2-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("O=test organization,L=New York,ST=NY,C=US"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName:   "test2-example.com",
							Country:      "US",
							Locality:     "New York",
							Organization: "test organization",
							State:        "New York",
						},
						TrustChainPEM: ptr.To("-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"),
					},
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert1_1234",
						CertificateName:                     "Test Certificate1",
						CertificateStatus:                   "ACTIVE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "2048",
						KeyType:                             "RSA",
						SANs:                                []string{"test-example.com", "www.test-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName: "test-example.com",
							Country:    "US",
							Locality:   "Cambridge",
							State:      "Massachusetts",
						},
						TrustChainPEM: nil,
					},
				},
				Links: Links{
					Self:     "/ccm/v1/certificates?sort=-certificateName&page=2&pageSize=2",
					Next:     nil,
					Previous: ptr.To("/ccm/v1/certificates?sort=-certificateName&page=1&pageSize=2"),
				},
				Metadata: ListMetadata{
					TotalItems: 4,
					TotalPages: 2,
				},
			},
			responseStatus: 200,
			responseBody: `
			{
				"certificates": [
					{
						"accountId": "test-account-id",
						"certificateId": "cert2_1234",
						"certificateName": "Test Certificate2",
						"certificateStatus": "READY_FOR_USE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-09-15T10:20:30.123456Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-09-15T10:20:33.123456Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "P-256",
						"keyType": "ECDSA",
						"sans": [
							"test2-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "O=test organization,L=New York,ST=NY,C=US",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test2-example.com",
							"country": "US",
							"locality": "New York",
							"organization": "test organization",
							"state": "New York"
						},
						"trustChainPem": "-----BEGIN CERTIFICATE-----\nexample-trust-chain-PEM\n-----END CERTIFICATE-----\n"
					},
					{
						"accountId": "test-account-id",
						"certificateId": "cert1_1234",
						"certificateName": "Test Certificate1",
						"certificateStatus": "ACTIVE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-08-22T09:01:32.607357Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "2048",
						"keyType": "RSA",
						"sans": [
							"test-example.com",
							"www.test-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test-example.com",
							"country": "US",
							"locality": "Cambridge",
							"state": "Massachusetts"
						},
						"trustChainPem": null
					}
				],
				"links": {
					"self": "/ccm/v1/certificates?sort=-certificateName&page=2&pageSize=2",
					"next": null,
					"previous": "/ccm/v1/certificates?sort=-certificateName&page=1&pageSize=2"
				},
				"metadata": {
					"totalItems": 4,
					"totalPages": 2
				}
			}`,
			expectedPath: "/ccm/v1/certificates?page=2&pageSize=2&sort=-certificateName",
		},
		"200 OK - with pagination and keyType filtering": {
			params: ListCertificatesRequest{
				KeyType:  CryptographicAlgorithmRSA,
				PageSize: 1,
				Page:     1,
			},
			expectedResponse: &ListCertificatesResponse{
				Certificates: []Certificate{
					{
						AccountID:                           "test-account-id",
						CertificateID:                       "cert1_1234",
						CertificateName:                     "Test Certificate1",
						CertificateStatus:                   "ACTIVE",
						CertificateType:                     "THIRD_PARTY",
						CreatedBy:                           "jkowalski",
						CreatedDate:                         test.NewTimeFromString(t, "2025-08-22T09:01:32.607357Z"),
						ModifiedBy:                          "jkowalski",
						ModifiedDate:                        test.NewTimeFromString(t, "2025-08-22T09:01:33.607357Z"),
						CSRExpirationDate:                   test.NewTimeFromString(t, "2026-10-24T09:01:34Z"),
						KeySize:                             "2048",
						KeyType:                             "RSA",
						SANs:                                []string{"test-example.com", "www.test-example.com"},
						SecureNetwork:                       "ENHANCED_TLS",
						SignedCertificateIssuer:             ptr.To("CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-11-22T12:11:31Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2025-08-22T11:11:31Z")),
						Subject: &Subject{
							CommonName: "test-example.com",
							Country:    "US",
							Locality:   "Cambridge",
							State:      "Massachusetts",
						},
						TrustChainPEM: nil,
					},
				},
				Links: Links{
					Self:     "/ccm/v1/certificates?keyType=RSA&page=1&pageSize=1",
					Next:     ptr.To("/ccm/v1/certificates?keyType=RSA&page=2&pageSize=1"),
					Previous: nil,
				},
				Metadata: ListMetadata{
					TotalItems: 2,
					TotalPages: 2,
				},
			},
			responseStatus: 200,
			responseBody: `
			{
				"certificates": [
					{
						"accountId": "test-account-id",
						"certificateId": "cert1_1234",
						"certificateName": "Test Certificate1",
						"certificateStatus": "ACTIVE",
						"certificateType": "THIRD_PARTY",
						"createdBy": "jkowalski",
						"createdDate": "2025-08-22T09:01:32.607357Z",
						"modifiedBy": "jkowalski",
						"modifiedDate": "2025-08-22T09:01:33.607357Z",
						"csrExpirationDate": "2026-10-24T09:01:34Z",
						"keySize": "2048",
						"keyType": "RSA",
						"sans": [
							"test-example.com",
							"www.test-example.com"
						],
						"secureNetwork": "ENHANCED_TLS",
						"signedCertificateIssuer": "CN=mkcert user (name surname),OU=organization (name surname),O=mkcert development CA",
						"signedCertificateNotValidAfterDate": "2027-11-22T12:11:31Z",
						"signedCertificateNotValidBeforeDate": "2025-08-22T11:11:31Z",
						"subject": {
							"commonName": "test-example.com",
							"country": "US",
							"locality": "Cambridge",
							"state": "Massachusetts"
						},
						"trustChainPem": null
					}
				],
				"links": {
					"self": "/ccm/v1/certificates?keyType=RSA&page=1&pageSize=1",
					"next": "/ccm/v1/certificates?keyType=RSA&page=2&pageSize=1",
					"previous": null
				},
				"metadata": {
					"totalItems": 2,
					"totalPages": 2
				}
			}`,
			expectedPath: "/ccm/v1/certificates?keyType=RSA&page=1&pageSize=1",
		},
		"validation error - invalid certificate status": {
			params: ListCertificatesRequest{
				CertificateStatus: []CertificateStatus{"INVALID_STATUS"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificates: struct validation: CertificateStatus: list '[INVALID_STATUS]' contains invalid element 'INVALID_STATUS'. "+
					"Each element must be one of: 'ACTIVE', 'READY_FOR_USE', 'CSR_READY'", err.Error())
				assert.ErrorIs(t, err, ErrListCertificates)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid key type": {
			params: ListCertificatesRequest{
				KeyType: "INVALID_KEY",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificates: struct validation: KeyType: value 'INVALID_KEY' is invalid. Must be either 'RSA' or 'ECDSA'", err.Error())
				assert.ErrorIs(t, err, ErrListCertificates)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page size less than 1": {
			params: ListCertificatesRequest{
				PageSize: -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificates: struct validation: PageSize: must be 1 or greater", err.Error())
				assert.ErrorIs(t, err, ErrListCertificates)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page size greater than 100": {
			params: ListCertificatesRequest{
				PageSize: 101,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificates: struct validation: PageSize: cannot be greater than 100", err.Error())
				assert.ErrorIs(t, err, ErrListCertificates)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - page value less than 1": {
			params: ListCertificatesRequest{
				Page: -1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificates: struct validation: Page: must be 1 or greater", err.Error())
				assert.ErrorIs(t, err, ErrListCertificates)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid field for sort": {
			params: ListCertificatesRequest{
				Sort: "invalid_sort",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificates: struct validation: Sort: must be a comma-separated list of fields, optionally prefixed by + or - (e.g. +createdDate,-certificateName)", err.Error())
				assert.ErrorIs(t, err, ErrListCertificates)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid prefix for sort": {
			params: ListCertificatesRequest{
				Sort: "*createdDate",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "listing certificates: struct validation: Sort: must be a comma-separated list of fields, optionally prefixed by + or - (e.g. +createdDate,-certificateName)", err.Error())
				assert.ErrorIs(t, err, ErrListCertificates)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"500 internal server error - assert that error is ErrListCertificates": {
			params:         ListCertificatesRequest{},
			responseStatus: 500,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error retrieving certificates",
				"status": 500
			}`,
			expectedPath: "/ccm/v1/certificates",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListCertificates, &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error retrieving certificates",
					Status: http.StatusInternalServerError,
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrListCertificates)
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
			client := mockAPIClient(t, mockServer)
			result, err := client.ListCertificates(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestSortValidationRule(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		sortValue string
		hasError  bool
	}{
		"valid sort value - single field, ascending": {
			sortValue: "+modifiedDate",
			hasError:  false,
		},
		"valid sort value - single field, descending": {
			sortValue: "-certificateName",
			hasError:  false,
		},
		"valid sort value - single field, without prefix": {
			sortValue: "modifiedDate",
			hasError:  false,
		},
		"valid sort value - multiple fields, mixed order": {
			sortValue: "+certificateName,-expirationDate,-createdDate",
			hasError:  false,
		},
		"valid sort value - empty string": {
			sortValue: "",
			hasError:  true,
		},
		"invalid sort value - invalid field": {
			sortValue: "invalidField",
			hasError:  true,
		},
		"invalid sort value - invalid prefix": {
			sortValue: "*createdDate",
			hasError:  true,
		},
		"invalid sort value - invalid separator": {
			sortValue: "+createdDate|-certificateName",
			hasError:  true,
		},
		"invalid sort value - multiple fields, mixed order": {
			sortValue: "+certificateName,-modifiedBy,-createdDate",
			hasError:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := sortValidationRule(tc.sortValue)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
