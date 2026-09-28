package cloudcertificates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLineage(t *testing.T) {
	t.Parallel()

	baseRequest := CreateLineageRequest{
		Body: CreateLineageRequestBody{
			ContractID:    "C-0N7RAC7",
			GroupID:       12345,
			GeoClass:      GeoClassStandardWorldwide,
			SecureNetwork: SecureNetworkEnhancedTLS,
			KeySpecs: []KeySpec{
				{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
				{KeyType: CryptographicAlgorithmECDSA, KeySize: KeySizeP256},
			},
			SANs: []string{"www.example.com", "example.com"},
			Subject: ptr.To(Subject{
				CommonName:         "example.com",
				Organization:       "Example Corp.",
				OrganizationalUnit: "IT",
				Country:            "US",
				State:              "Massachusetts",
				Locality:           "Cambridge",
			}),
		},
	}

	baseRequestBody := `{
		"contractId": "C-0N7RAC7",
		"groupId": 12345,
		"geoClass": "STANDARD_WORLDWIDE",
		"secureNetwork": "ENHANCED_TLS",
		"keySpecs": [
			{"keyType": "RSA", "keySize": "2048"},
			{"keyType": "ECDSA", "keySize": "P-256"}
		],
		"sans": ["www.example.com", "example.com"],
		"subject": {
			"commonName": "example.com",
			"organization": "Example Corp.",
			"organizationalUnit": "IT",
			"country": "US",
			"state": "Massachusetts",
			"locality": "Cambridge"
		}
	}`

	tests := map[string]struct {
		params              CreateLineageRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedResponse    *CreateLineageResponse
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"201 Created - create lineage": {
			params:              baseRequest,
			expectedPath:        "/ccm/v2/lineages",
			expectedRequestBody: baseRequestBody,
			responseStatus:      http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-07T10:18:46Z",
							"algorithmInstanceId": 6122,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-07-07T10:18:46Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-07T10:18:46Z",
							"algorithmInstanceId": 6123,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-07-07T10:18:46Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-07T10:18:46Z",
					"generationModifiedBy": "terraform-dev",
					"headGenerationId": 2912,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-07T10:18:46Z",
				"lineageId": 500001,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-07T10:18:46Z",
				"lineageName": "test non-validated domain",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "MULTIPLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				},
				"validationResults": {
					"warnings": [
						{
							"context": {"domValidatorLink": "/domain-validation/v1/domains"},
							"detail": "Domain validation must be completed before uploading the signed certificate.",
							"instance": "/error-types/domain-not-validated?traceId=1234567891014",
							"status": 400,
							"title": "Some SANs in the request are not Domain Validated.",
							"type": "/error-types/domain-not-validated"
						}
					]
				}
			}`,
			expectedResponse: &CreateLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-07T10:18:46Z")),
								AlgorithmInstanceID:          6122,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-07T10:18:46Z")),
								AlgorithmInstanceID:          6123,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-07T10:18:46Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					HeadGenerationID:     2912,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-07T10:18:46Z"),
				LineageID:           500001,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-07T10:18:46Z"),
				LineageName:         "test non-validated domain",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeMultipleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
				ValidationResults: &ValidationResults{
					Warnings: []ValidationResultItem{
						{
							Context:  map[string]any{"domValidatorLink": "/domain-validation/v1/domains"},
							Detail:   "Domain validation must be completed before uploading the signed certificate.",
							Instance: "/error-types/domain-not-validated?traceId=1234567891014",
							Status:   400,
							Title:    "Some SANs in the request are not Domain Validated.",
							Type:     "/error-types/domain-not-validated",
						},
					},
				},
			},
		},
		"201 Created - create lineage with single key spec infers SINGLE_STACK": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					GeoClass:      GeoClassStandardWorldwide,
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs: []string{"www.example.com", "example.com"},
					Subject: ptr.To(Subject{
						CommonName:         "example.com",
						Organization:       "Example Corp.",
						OrganizationalUnit: "IT",
						Country:            "US",
						State:              "Massachusetts",
						Locality:           "Cambridge",
					}),
				},
			},
			expectedPath: "/ccm/v2/lineages",
			expectedRequestBody: `{
				"contractId": "C-0N7RAC7",
				"groupId": 12345,
				"geoClass": "STANDARD_WORLDWIDE",
				"secureNetwork": "ENHANCED_TLS",
				"keySpecs": [
					{"keyType": "RSA", "keySize": "2048"}
				],
				"sans": ["www.example.com", "example.com"],
				"subject": {
					"commonName": "example.com",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"country": "US",
					"state": "Massachusetts",
					"locality": "Cambridge"
				}
			}`,
			responseStatus: http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-17T12:28:59Z",
							"algorithmInstanceId": 1714,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-07-17T12:28:59Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-17T12:28:59Z",
					"generationModifiedBy": "terraform-dev",
					"headGenerationId": 929,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-17T12:28:59Z",
				"lineageId": 500005,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-17T12:28:59Z",
				"lineageName": "www.example.com20260717122858498674",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "SINGLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &CreateLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-17T12:28:59Z")),
								AlgorithmInstanceID:          1714,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-07-17T12:28:59Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-17T12:28:59Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					HeadGenerationID:     929,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-17T12:28:59Z"),
				LineageID:           500005,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-17T12:28:59Z"),
				LineageName:         "www.example.com20260717122858498674",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeSingleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"201 Created - create lineage without GeoClass omits it from request body": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs: []string{"www.example.com", "example.com"},
					Subject: ptr.To(Subject{
						CommonName:         "example.com",
						Organization:       "Example Corp.",
						OrganizationalUnit: "IT",
						Country:            "US",
						State:              "Massachusetts",
						Locality:           "Cambridge",
					}),
				},
			},
			expectedPath: "/ccm/v2/lineages",
			expectedRequestBody: `{
				"contractId": "C-0N7RAC7",
				"groupId": 12345,
				"secureNetwork": "ENHANCED_TLS",
				"keySpecs": [
					{"keyType": "RSA", "keySize": "2048"}
				],
				"sans": ["www.example.com", "example.com"],
				"subject": {
					"commonName": "example.com",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"country": "US",
					"state": "Massachusetts",
					"locality": "Cambridge"
				}
			}`,
			responseStatus: http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-17T12:28:59Z",
							"algorithmInstanceId": 1714,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-07-17T12:28:59Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-17T12:28:59Z",
					"generationModifiedBy": "terraform-dev",
					"headGenerationId": 929,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-17T12:28:59Z",
				"lineageId": 500005,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-17T12:28:59Z",
				"lineageName": "www.example.com20260717122858498674",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "SINGLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &CreateLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-17T12:28:59Z")),
								AlgorithmInstanceID:          1714,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-07-17T12:28:59Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-17T12:28:59Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					HeadGenerationID:     929,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-17T12:28:59Z"),
				LineageID:           500005,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-17T12:28:59Z"),
				LineageName:         "www.example.com20260717122858498674",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeSingleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"201 Created - create lineage without Subject omits it from request body": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs: []string{"www.example.com", "example.com"},
				},
			},
			expectedPath: "/ccm/v2/lineages",
			expectedRequestBody: `{
				"contractId": "C-0N7RAC7",
				"groupId": 12345,
				"secureNetwork": "ENHANCED_TLS",
				"keySpecs": [
					{"keyType": "RSA", "keySize": "2048"}
				],
				"sans": ["www.example.com", "example.com"]
			}`,
			responseStatus: http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-20T10:32:09Z",
							"algorithmInstanceId": 2244,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-07-20T10:32:09Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-20T10:32:09Z",
					"generationModifiedBy": "terraform-dev",
					"headGenerationId": 1323,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-20T10:32:09Z",
				"lineageId": 500006,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-20T10:32:09Z",
				"lineageName": "www.example.com20260720103154092771",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "SINGLE_STACK",
				"subject": {}
			}`,
			expectedResponse: &CreateLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-20T10:32:09Z")),
								AlgorithmInstanceID:          2244,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-07-20T10:32:09Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-20T10:32:09Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					HeadGenerationID:     1323,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-20T10:32:09Z"),
				LineageID:           500006,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-20T10:32:09Z"),
				LineageName:         "www.example.com20260720103154092771",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeSingleStack),
				Subject:             Subject{},
			},
		},
		"201 Created - create lineage with explicit LineageName": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					LineageName:   "test 1",
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
						{KeyType: CryptographicAlgorithmECDSA, KeySize: KeySizeP256},
					},
					SANs: []string{"www.example.com", "example.com"},
					Subject: ptr.To(Subject{
						CommonName:         "example.com",
						Organization:       "Example Corp.",
						OrganizationalUnit: "IT",
						Country:            "US",
						State:              "Massachusetts",
						Locality:           "Cambridge",
					}),
				},
			},
			expectedPath: "/ccm/v2/lineages",
			expectedRequestBody: `{
				"contractId": "C-0N7RAC7",
				"groupId": 12345,
				"lineageName": "test 1",
				"secureNetwork": "ENHANCED_TLS",
				"keySpecs": [
					{"keyType": "RSA", "keySize": "2048"},
					{"keyType": "ECDSA", "keySize": "P-256"}
				],
				"sans": ["www.example.com", "example.com"],
				"subject": {
					"commonName": "example.com",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"country": "US",
					"state": "Massachusetts",
					"locality": "Cambridge"
				}
			}`,
			responseStatus: http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T08:44:31Z",
							"algorithmInstanceId": 4546,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T08:44:31Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T08:44:31Z",
							"algorithmInstanceId": 4547,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T08:44:31Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-27T08:44:31Z",
					"generationModifiedBy": "terraform-dev",
					"headGenerationId": 3171,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-27T08:44:31Z",
				"lineageId": 500008,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-27T08:44:31Z",
				"lineageName": "test 1",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "MULTIPLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &CreateLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T08:44:31Z")),
								AlgorithmInstanceID:          4546,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T08:44:31Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T08:44:31Z")),
								AlgorithmInstanceID:          4547,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T08:44:31Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T08:44:31Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					HeadGenerationID:     3171,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-27T08:44:31Z"),
				LineageID:           500008,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-27T08:44:31Z"),
				LineageName:         "test 1",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeMultipleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"201 Created - create lineage with single key spec and explicit LineageName": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					LineageName:   "single cipher",
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs: []string{"www.example.com", "example.com"},
					Subject: ptr.To(Subject{
						CommonName:         "example.com",
						Organization:       "Example Corp.",
						OrganizationalUnit: "IT",
						Country:            "US",
						State:              "Massachusetts",
						Locality:           "Cambridge",
					}),
				},
			},
			expectedPath: "/ccm/v2/lineages",
			expectedRequestBody: `{
				"contractId": "C-0N7RAC7",
				"groupId": 12345,
				"lineageName": "single cipher",
				"secureNetwork": "ENHANCED_TLS",
				"keySpecs": [
					{"keyType": "RSA", "keySize": "2048"}
				],
				"sans": ["www.example.com", "example.com"],
				"subject": {
					"commonName": "example.com",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"country": "US",
					"state": "Massachusetts",
					"locality": "Cambridge"
				}
			}`,
			responseStatus: http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T12:49:11Z",
							"algorithmInstanceId": 4673,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T12:49:11Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-27T12:49:11Z",
					"generationModifiedBy": "terraform-dev",
					"headGenerationId": 3247,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-27T12:49:11Z",
				"lineageId": 500016,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-27T12:49:11Z",
				"lineageName": "single cipher",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "SINGLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &CreateLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T12:49:11Z")),
								AlgorithmInstanceID:          4673,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T12:49:11Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T12:49:11Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					HeadGenerationID:     3247,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-27T12:49:11Z"),
				LineageID:           500016,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-27T12:49:11Z"),
				LineageName:         "single cipher",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeSingleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"403 Forbidden - account not allowed": {
			params:              baseRequest,
			expectedPath:        "/ccm/v2/lineages",
			expectedRequestBody: baseRequestBody,
			responseStatus:      http.StatusForbidden,
			responseBody: `{
				"type": "/error-types/lineage-account-not-allowed",
				"title": "Account is not allowed to use Certificate Lineage.",
				"status": 403,
				"detail": "This account is not permitted to access Certificate Lineage APIs.",
				"instance": "/error-types/lineage-account-not-allowed?traceId=1234567891011"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCreateLineage, &Error{
					Type:     "/error-types/lineage-account-not-allowed",
					Title:    "Account is not allowed to use Certificate Lineage.",
					Status:   http.StatusForbidden,
					Detail:   "This account is not permitted to access Certificate Lineage APIs.",
					Instance: "/error-types/lineage-account-not-allowed?traceId=1234567891011",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageAccountNotAllowed)
				assert.ErrorIs(t, err, ErrCreateLineage)
			},
		},
		"400 Bad Request - schema validation failure, invalid SAN": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					LineageName:   "test-lineage",
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
						{KeyType: CryptographicAlgorithmECDSA, KeySize: KeySizeP256},
					},
					SANs: []string{"www.example.com", "example com"},
					Subject: ptr.To(Subject{
						CommonName:         "example.com",
						Organization:       "Example Corp.",
						OrganizationalUnit: "IT",
						Country:            "US",
						State:              "Massachusetts",
						Locality:           "Cambridge",
					}),
				},
			},
			expectedPath: "/ccm/v2/lineages",
			expectedRequestBody: `{
				"contractId": "C-0N7RAC7",
				"groupId": 12345,
				"lineageName": "test-lineage",
				"secureNetwork": "ENHANCED_TLS",
				"keySpecs": [
					{"keyType": "RSA", "keySize": "2048"},
					{"keyType": "ECDSA", "keySize": "P-256"}
				],
				"sans": ["www.example.com", "example com"],
				"subject": {
					"commonName": "example.com",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"country": "US",
					"state": "Massachusetts",
					"locality": "Cambridge"
				}
			}`,
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/schema-validation-failure",
				"title": "Request body failed JSON schema validation.",
				"status": 400,
				"detail": "ECMA 262 regex \"^(\\*\\.)?([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,}$\" does not match input string \"example com\"",
				"instance": "/error-types/schema-validation-failure?traceId=6556339828170181119",
				"context": {
					"pointer": "/sans/1",
					"regex": "^(\\*\\.)?([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,}$",
					"schemaUri": "open-api-v1-san-domain.json",
					"string": "example com"
				}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCreateLineage, &Error{
					Type:     "/error-types/schema-validation-failure",
					Title:    "Request body failed JSON schema validation.",
					Status:   http.StatusBadRequest,
					Detail:   `ECMA 262 regex "^(\*\.)?([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$" does not match input string "example com"`,
					Instance: "/error-types/schema-validation-failure?traceId=6556339828170181119",
					Context: map[string]any{
						"pointer":   "/sans/1",
						"regex":     `^(\*\.)?([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`,
						"schemaUri": "open-api-v1-san-domain.json",
						"string":    "example com",
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrSchemaValidationFailure)
				assert.ErrorIs(t, err, ErrCreateLineage)
			},
		},
		"validation error - missing ContractID": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					GroupID:       12345,
					GeoClass:      GeoClassStandardWorldwide,
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs:    []string{"example.com"},
					Subject: ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tContractID: cannot be blank\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing GroupID": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GeoClass:      GeoClassStandardWorldwide,
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs:    []string{"example.com"},
					Subject: ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tGroupID: cannot be blank\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing KeySpecs": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					GeoClass:      GeoClassStandardWorldwide,
					SecureNetwork: SecureNetworkEnhancedTLS,
					SANs:          []string{"example.com"},
					Subject:       ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tKeySpecs: cannot be blank\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing SANs": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					GeoClass:      GeoClassStandardWorldwide,
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					Subject: ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tSANs: cannot be blank\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing SecureNetwork": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID: "C-0N7RAC7",
					GroupID:    12345,
					GeoClass:   GeoClassStandardWorldwide,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs:    []string{"example.com"},
					Subject: ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tSecureNetwork: cannot be blank\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid GeoClass": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					GeoClass:      "INVALID_CLASS",
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs:    []string{"example.com"},
					Subject: ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tGeoClass: value "+
					"'INVALID_CLASS' is invalid. Must be 'STANDARD_WORLDWIDE'\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid LineageType": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					GeoClass:      GeoClassStandardWorldwide,
					LineageType:   "INVALID_TYPE",
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs:    []string{"example.com"},
					Subject: ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tLineageType: value "+
					"'INVALID_TYPE' is invalid. Must be either 'MULTIPLE_GENERATION' or 'SINGLE_GENERATION'\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid SecureNetwork": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					GeoClass:      GeoClassStandardWorldwide,
					SecureNetwork: "NOT_A_REAL_NETWORK",
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs:    []string{"example.com"},
					Subject: ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tSecureNetwork: value "+
					"'NOT_A_REAL_NETWORK' is invalid. Must be either 'ENHANCED_TLS' or 'STANDARD_TLS'\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - LineageName too long": {
			params: CreateLineageRequest{
				Body: CreateLineageRequestBody{
					ContractID:    "C-0N7RAC7",
					GroupID:       12345,
					GeoClass:      GeoClassStandardWorldwide,
					LineageName:   strings.Repeat("a", 271),
					SecureNetwork: SecureNetworkEnhancedTLS,
					KeySpecs: []KeySpec{
						{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
					},
					SANs:    []string{"example.com"},
					Subject: ptr.To(Subject{CommonName: "example.com"}),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "creating lineage: struct validation: Body: {\n\tLineageName: the length must be no more than 270\n}")
				assert.ErrorIs(t, err, ErrCreateLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if tc.expectedRequestBody != "" {
					requestBody, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(requestBody))
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.CreateLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestLineageType_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		lineageType LineageType
		withError   bool
	}{
		"empty value is valid - server assigns default": {
			lineageType: "",
			withError:   false,
		},
		"MULTIPLE_GENERATION is valid": {
			lineageType: LineageTypeMultipleGeneration,
			withError:   false,
		},
		"SINGLE_GENERATION is valid": {
			lineageType: LineageTypeSingleGeneration,
			withError:   false,
		},
		"invalid value": {
			lineageType: "NOT_A_REAL_TYPE",
			withError:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.lineageType.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, "value 'NOT_A_REAL_TYPE' is invalid. Must be either "+
					"'MULTIPLE_GENERATION' or 'SINGLE_GENERATION'")
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGeoClass_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		geoClass  GeoClass
		withError bool
	}{
		"empty value is valid": {
			geoClass:  "",
			withError: false,
		},
		"STANDARD_WORLDWIDE is valid": {
			geoClass:  GeoClassStandardWorldwide,
			withError: false,
		},
		"invalid value": {
			geoClass:  "NOT_A_REAL_CLASS",
			withError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.geoClass.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, "value 'NOT_A_REAL_CLASS' is invalid. Must be 'STANDARD_WORLDWIDE'")
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestExpandGenerations_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expandGenerations ExpandGenerations
		withError         bool
	}{
		"empty value is valid": {
			expandGenerations: "",
			withError:         false,
		},
		"HEAD is valid": {
			expandGenerations: ExpandGenerationsHead,
			withError:         false,
		},
		"CURRENT_PRODUCTION is valid": {
			expandGenerations: ExpandGenerationsCurrentProduction,
			withError:         false,
		},
		"CURRENT_STAGING is valid": {
			expandGenerations: ExpandGenerationsCurrentStaging,
			withError:         false,
		},
		"PREVIOUS_PRODUCTION is valid": {
			expandGenerations: ExpandGenerationsPreviousProduction,
			withError:         false,
		},
		"invalid value": {
			expandGenerations: "ALL",
			withError:         true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.expandGenerations.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, "value 'ALL' is invalid. Must be one of: "+
					"'HEAD', 'CURRENT_PRODUCTION', 'CURRENT_STAGING', or 'PREVIOUS_PRODUCTION'")
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestKeySpec_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		keySpec   KeySpec
		withError bool
		errorMsg  string
	}{
		"RSA 2048 is valid": {
			keySpec:   KeySpec{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize2048},
			withError: false,
		},
		"RSA 4096 is valid": {
			keySpec:   KeySpec{KeyType: CryptographicAlgorithmRSA, KeySize: KeySize4096},
			withError: false,
		},
		"ECDSA P-256 is valid": {
			keySpec:   KeySpec{KeyType: CryptographicAlgorithmECDSA, KeySize: KeySizeP256},
			withError: false,
		},
		"ECDSA P-384 is valid": {
			keySpec:   KeySpec{KeyType: CryptographicAlgorithmECDSA, KeySize: KeySizeP384},
			withError: false,
		},
		"missing KeyType": {
			keySpec:   KeySpec{KeySize: KeySize2048},
			withError: true,
			errorMsg:  "KeyType: cannot be blank.",
		},
		"missing KeySize": {
			keySpec:   KeySpec{KeyType: CryptographicAlgorithmRSA},
			withError: true,
			errorMsg:  "KeySize: cannot be blank.",
		},
		"invalid KeyType": {
			keySpec:   KeySpec{KeyType: "INVALID", KeySize: KeySize2048},
			withError: true,
			errorMsg:  "KeyType: value 'INVALID' is invalid. Must be either 'RSA' or 'ECDSA'.",
		},
		"invalid KeySize": {
			keySpec:   KeySpec{KeyType: CryptographicAlgorithmRSA, KeySize: "1024"},
			withError: true,
			errorMsg:  "KeySize: value '1024' is invalid. Must be one of: '2048', '4096', 'P-256', or 'P-384'.",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.keySpec.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, tc.errorMsg)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetLineage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           GetLineageRequest
		responseStatus   int
		responseBody     string
		expectedResponse *GetLineageResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"200 OK - get lineage with expanded head": {
			params:         GetLineageRequest{LineageID: 500001, ExpandGenerations: []ExpandGenerations{ExpandGenerationsHead}},
			expectedPath:   "/ccm/v2/lineages/500001?expandGenerations=HEAD",
			responseStatus: http.StatusOK,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-07T10:18:46Z",
							"algorithmInstanceId": 6123,
							"algorithmInstanceModifiedBy": "terraform-dev",
							"algorithmInstanceModifiedTime": "2026-07-08T13:23:37Z",
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-07-07T10:18:46Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-07-08T13:23:37Z",
							"signedCertificateNotValidBeforeDate": "2026-07-08T13:23:37Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:02",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:02"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-07T10:18:46Z",
							"algorithmInstanceId": 6122,
							"algorithmInstanceModifiedBy": "terraform-dev",
							"algorithmInstanceModifiedTime": "2026-07-08T13:06:53Z",
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-07-07T10:18:46Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-07-08T13:06:53Z",
							"signedCertificateNotValidBeforeDate": "2026-07-08T13:06:53Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:01"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-07T10:18:46Z",
					"generationModifiedBy": "terraform-dev",
					"generationModifiedTime": "2026-07-08T13:24:17Z",
					"headGenerationId": 2912,
					"headGenerationStatus": "READY_FOR_USE"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-07T10:18:46Z",
				"lineageId": 500001,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-08T13:24:17Z",
				"lineageName": "www.example.com20260707101743398222",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "MULTIPLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &GetLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-07T10:18:46Z")),
								AlgorithmInstanceID:                 6123,
								AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
								AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-08T13:23:37Z")),
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
								CSRPEM:                              ecdsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmECDSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-08T13:23:37Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-08T13:23:37Z")),
								SignedCertificatePEM:                ptr.To(ecdsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:02"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:02"),
							},
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-07T10:18:46Z")),
								AlgorithmInstanceID:                 6122,
								AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
								AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-08T13:06:53Z")),
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
								CSRPEM:                              rsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmRSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-08T13:06:53Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-08T13:06:53Z")),
								SignedCertificatePEM:                ptr.To(rsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:01"),
							},
						},
						GenerationCreatedBy:    ptr.To("terraform-dev"),
						GenerationCreatedTime:  ptr.To(test.NewTimeFromString(t, "2026-07-07T10:18:46Z")),
						GenerationModifiedBy:   ptr.To("terraform-dev"),
						GenerationModifiedTime: ptr.To(test.NewTimeFromString(t, "2026-07-08T13:24:17Z")),
					},
					HeadGenerationID:     2912,
					HeadGenerationStatus: string(GenerationStatusReadyForUse),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-07T10:18:46Z"),
				LineageID:           500001,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-08T13:24:17Z"),
				LineageName:         "www.example.com20260707101743398222",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeMultipleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"200 OK - get lineage with all four generation pointers expanded": {
			params: GetLineageRequest{
				LineageID: 500022,
				ExpandGenerations: []ExpandGenerations{
					ExpandGenerationsHead, ExpandGenerationsCurrentProduction,
					ExpandGenerationsCurrentStaging, ExpandGenerationsPreviousProduction,
				},
			},
			expectedPath:   "/ccm/v2/lineages/500022?expandGenerations=HEAD%2CCURRENT_PRODUCTION%2CCURRENT_STAGING%2CPREVIOUS_PRODUCTION",
			responseStatus: http.StatusOK,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"currentProduction": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-23T11:20:11Z",
							"algorithmInstanceId": 4343,
							"algorithmInstanceModifiedBy": "terraform-dev",
							"algorithmInstanceModifiedTime": "2026-07-23T11:21:41Z",
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-09-24T11:20:11Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-07-23T11:21:00Z",
							"signedCertificateNotValidBeforeDate": "2026-07-23T11:21:00Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:04",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:04"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-23T11:20:11Z",
							"algorithmInstanceId": 4344,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-24T11:20:11Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						}
					],
					"firstPromotedToProductionTime": "2026-07-23T11:25:05Z",
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-23T11:20:11Z",
					"generationModifiedBy": "terraform-dev",
					"generationModifiedTime": "2026-07-23T11:21:42Z",
					"productionGenerationId": 3028,
					"productionGenerationStatus": "ACTIVE"
				},
				"currentStaging": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-23T11:20:11Z",
							"algorithmInstanceId": 4343,
							"algorithmInstanceModifiedBy": "terraform-dev",
							"algorithmInstanceModifiedTime": "2026-07-23T11:21:41Z",
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-09-24T11:20:11Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-07-23T11:21:00Z",
							"signedCertificateNotValidBeforeDate": "2026-07-23T11:21:00Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:04",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:04"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-23T11:20:11Z",
							"algorithmInstanceId": 4344,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-24T11:20:11Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						}
					],
					"firstPromotedToProductionTime": "2026-07-23T11:25:05Z",
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-23T11:20:11Z",
					"generationModifiedBy": "terraform-dev",
					"generationModifiedTime": "2026-07-23T11:21:42Z",
					"stagingGenerationId": 3028,
					"stagingGenerationStatus": "ACTIVE"
				},
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-23T11:26:01Z",
							"algorithmInstanceId": 4345,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-24T11:26:01Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-23T11:26:01Z",
							"algorithmInstanceId": 4346,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-24T11:26:01Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-23T11:26:01Z",
					"generationModifiedBy": "terraform-dev",
					"headGenerationId": 3029,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-22T09:21:09Z",
				"lineageId": 500022,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-23T11:26:01Z",
				"lineageName": "alpha 2",
				"lineageType": "MULTIPLE_GENERATION",
				"previousProduction": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-23T11:13:41Z",
							"algorithmInstanceId": 4341,
							"algorithmInstanceModifiedBy": "terraform-dev",
							"algorithmInstanceModifiedTime": "2026-07-23T11:16:00Z",
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-09-24T11:13:41Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-07-23T11:14:45Z",
							"signedCertificateNotValidBeforeDate": "2026-07-23T11:14:45Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:05",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:05"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-23T11:13:41Z",
							"algorithmInstanceId": 4342,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-24T11:13:41Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						}
					],
					"firstPromotedToProductionTime": "2026-07-23T11:18:53Z",
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-23T11:13:41Z",
					"generationModifiedBy": "terraform-dev",
					"generationModifiedTime": "2026-07-23T11:16:01Z",
					"previousProductionGenerationId": 3027,
					"previousProductionGenerationStatus": "READY_FOR_USE"
				},
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "MULTIPLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &GetLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:26:01Z")),
								AlgorithmInstanceID:          4345,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-24T11:26:01Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:26:01Z")),
								AlgorithmInstanceID:          4346,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-24T11:26:01Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:26:01Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					HeadGenerationID:     3029,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-22T09:21:09Z"),
				LineageID:           500022,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-23T11:26:01Z"),
				LineageName:         "alpha 2",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeMultipleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
				CurrentProduction: &ProductionGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:20:11Z")),
								AlgorithmInstanceID:                 4343,
								AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
								AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-23T11:21:41Z")),
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-24T11:20:11Z")),
								CSRPEM:                              rsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmRSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-23T11:21:00Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:21:00Z")),
								SignedCertificatePEM:                ptr.To(rsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:04"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:04"),
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:20:11Z")),
								AlgorithmInstanceID:          4344,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-24T11:20:11Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
						},
						FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:25:05Z")),
						GenerationCreatedBy:           ptr.To("terraform-dev"),
						GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-23T11:20:11Z")),
						GenerationModifiedBy:          ptr.To("terraform-dev"),
						GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:21:42Z")),
					},
					ProductionGenerationID:     3028,
					ProductionGenerationStatus: string(GenerationStatusActive),
				},
				CurrentStaging: &StagingGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:20:11Z")),
								AlgorithmInstanceID:                 4343,
								AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
								AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-23T11:21:41Z")),
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-24T11:20:11Z")),
								CSRPEM:                              rsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmRSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-23T11:21:00Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:21:00Z")),
								SignedCertificatePEM:                ptr.To(rsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:04"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:04"),
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:20:11Z")),
								AlgorithmInstanceID:          4344,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-24T11:20:11Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
						},
						FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:25:05Z")),
						GenerationCreatedBy:           ptr.To("terraform-dev"),
						GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-23T11:20:11Z")),
						GenerationModifiedBy:          ptr.To("terraform-dev"),
						GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:21:42Z")),
					},
					StagingGenerationID:     3028,
					StagingGenerationStatus: string(GenerationStatusActive),
				},
				PreviousProduction: &PreviousProductionGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:13:41Z")),
								AlgorithmInstanceID:                 4341,
								AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
								AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-23T11:16:00Z")),
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-24T11:13:41Z")),
								CSRPEM:                              rsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmRSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-23T11:14:45Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:14:45Z")),
								SignedCertificatePEM:                ptr.To(rsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:05"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:05"),
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:13:41Z")),
								AlgorithmInstanceID:          4342,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-24T11:13:41Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
						},
						FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:18:53Z")),
						GenerationCreatedBy:           ptr.To("terraform-dev"),
						GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-23T11:13:41Z")),
						GenerationModifiedBy:          ptr.To("terraform-dev"),
						GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:16:01Z")),
					},
					PreviousProductionGenerationID:     3027,
					PreviousProductionGenerationStatus: string(GenerationStatusReadyForUse),
				},
			},
		},
		"200 OK - get lineage, fully promoted to production and staging, no head": {
			params: GetLineageRequest{
				LineageID:         500017,
				ExpandGenerations: []ExpandGenerations{ExpandGenerationsCurrentProduction, ExpandGenerationsCurrentStaging},
			},
			expectedPath:   "/ccm/v2/lineages/500017?expandGenerations=CURRENT_PRODUCTION%2CCURRENT_STAGING",
			responseStatus: http.StatusOK,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"currentProduction": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T13:28:47Z",
							"algorithmInstanceId": 4679,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T13:28:45Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T13:28:47Z",
							"algorithmInstanceId": 4678,
							"algorithmInstanceModifiedBy": "terraform-dev",
							"algorithmInstanceModifiedTime": "2026-07-27T13:33:40Z",
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-09-28T13:28:47Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-07-27T13:29:28Z",
							"signedCertificateNotValidBeforeDate": "2026-07-27T13:29:28Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:11",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:11"
						}
					],
					"firstPromotedToProductionTime": "2026-07-27T13:33:59Z",
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-27T13:28:47Z",
					"generationModifiedBy": "terraform-dev",
					"generationModifiedTime": "2026-07-27T13:33:41Z",
					"productionGenerationId": 3250,
					"productionGenerationStatus": "ACTIVE"
				},
				"currentStaging": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T13:28:47Z",
							"algorithmInstanceId": 4679,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T13:28:45Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T13:28:47Z",
							"algorithmInstanceId": 4678,
							"algorithmInstanceModifiedBy": "terraform-dev",
							"algorithmInstanceModifiedTime": "2026-07-27T13:33:40Z",
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-09-28T13:28:47Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-07-27T13:29:28Z",
							"signedCertificateNotValidBeforeDate": "2026-07-27T13:29:28Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:11",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:11"
						}
					],
					"firstPromotedToProductionTime": "2026-07-27T13:33:59Z",
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-27T13:28:47Z",
					"generationModifiedBy": "terraform-dev",
					"generationModifiedTime": "2026-07-27T13:33:41Z",
					"stagingGenerationId": 3250,
					"stagingGenerationStatus": "ACTIVE"
				},
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-27T13:28:47Z",
				"lineageId": 500017,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-27T13:33:59Z",
				"lineageName": "auto activate",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "MULTIPLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &GetLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head:       nil,
				CurrentProduction: &ProductionGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T13:28:47Z")),
								AlgorithmInstanceID:          4679,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T13:28:45Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-27T13:28:47Z")),
								AlgorithmInstanceID:                 4678,
								AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
								AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-27T13:33:40Z")),
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-28T13:28:47Z")),
								CSRPEM:                              rsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmRSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-27T13:29:28Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-27T13:29:28Z")),
								SignedCertificatePEM:                ptr.To(rsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:11"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:11"),
							},
						},
						FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T13:33:59Z")),
						GenerationCreatedBy:           ptr.To("terraform-dev"),
						GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-27T13:28:47Z")),
						GenerationModifiedBy:          ptr.To("terraform-dev"),
						GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-27T13:33:41Z")),
					},
					ProductionGenerationID:     3250,
					ProductionGenerationStatus: string(GenerationStatusActive),
				},
				CurrentStaging: &StagingGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T13:28:47Z")),
								AlgorithmInstanceID:          4679,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T13:28:45Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-27T13:28:47Z")),
								AlgorithmInstanceID:                 4678,
								AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
								AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-27T13:33:40Z")),
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-28T13:28:47Z")),
								CSRPEM:                              rsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmRSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-27T13:29:28Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-27T13:29:28Z")),
								SignedCertificatePEM:                ptr.To(rsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:11"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:11"),
							},
						},
						FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T13:33:59Z")),
						GenerationCreatedBy:           ptr.To("terraform-dev"),
						GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-27T13:28:47Z")),
						GenerationModifiedBy:          ptr.To("terraform-dev"),
						GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-27T13:33:41Z")),
					},
					StagingGenerationID:     3250,
					StagingGenerationStatus: string(GenerationStatusActive),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-27T13:28:47Z"),
				LineageID:           500017,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-27T13:33:59Z"),
				LineageName:         "auto activate",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeMultipleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"200 OK - get lineage without expandGenerations": {
			params:         GetLineageRequest{LineageID: 500002},
			expectedPath:   "/ccm/v2/lineages/500002",
			responseStatus: http.StatusOK,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"headGenerationId": 43,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-14T09:52:09Z",
				"lineageId": 500002,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-14T09:52:09Z",
				"lineageName": "www.example.com20260714095208534842",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "MULTIPLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &GetLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					HeadGenerationID:     43,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-14T09:52:09Z"),
				LineageID:           500002,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-14T09:52:09Z"),
				LineageName:         "www.example.com20260714095208534842",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeMultipleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"200 OK - get lineage without expandGenerations, created without a subject": {
			params:         GetLineageRequest{LineageID: 500035},
			expectedPath:   "/ccm/v2/lineages/500035",
			responseStatus: http.StatusOK,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"headGenerationId": 7270,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-08-13T08:52:53Z",
				"lineageId": 500035,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-08-13T08:52:53Z",
				"lineageName": "no subject",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "MULTIPLE_STACK",
				"subject": {}
			}`,
			expectedResponse: &GetLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					HeadGenerationID:     7270,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-08-13T08:52:53Z"),
				LineageID:           500035,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-08-13T08:52:53Z"),
				LineageName:         "no subject",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeMultipleStack),
				Subject:             Subject{},
			},
		},
		"404 lineage not found": {
			params:         GetLineageRequest{LineageID: 999999},
			expectedPath:   "/ccm/v2/lineages/999999",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891023",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetLineage, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891023",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrGetLineage)
			},
		},
		"500 internal server error": {
			params:         GetLineageRequest{LineageID: 500001},
			expectedPath:   "/ccm/v2/lineages/500001",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891016"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetLineage, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891016",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrGetLineage)
			},
		},
		"validation error - missing LineageID": {
			params: GetLineageRequest{},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "getting lineage: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrGetLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid ExpandGenerations": {
			params: GetLineageRequest{LineageID: 500001, ExpandGenerations: []ExpandGenerations{"ALL"}},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "getting lineage: struct validation: 0: value 'ALL' is invalid. "+
					"Must be one of: 'HEAD', 'CURRENT_PRODUCTION', 'CURRENT_STAGING', or 'PREVIOUS_PRODUCTION'")
				assert.ErrorIs(t, err, ErrGetLineage)
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
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.GetLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListLineages(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           ListLineagesRequest
		responseStatus   int
		responseBody     string
		expectedResponse *ListLineagesResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"200 OK - list lineages, no query parameters": {
			params:         ListLineagesRequest{},
			expectedPath:   "/ccm/v2/lineages",
			responseStatus: http.StatusOK,
			responseBody: `{
				"lineages": [
					{
						"accountId": "A-CCT1234",
						"contractId": "C-0N7RAC7",
						"geoClass": "STANDARD_WORLDWIDE",
						"groupId": 12345,
						"head": {
							"headGenerationId": 903,
							"headGenerationStatus": "CSR_READY"
						},
						"keySpecs": [
							{"keySize": "2048", "keyType": "RSA"},
							{"keySize": "P-256", "keyType": "ECDSA"}
						],
						"lineageCreatedBy": "terraform-dev",
						"lineageCreatedTime": "2026-07-14T09:52:09Z",
						"lineageId": 500002,
						"lineageModifiedBy": "terraform-dev",
						"lineageModifiedTime": "2026-07-17T07:58:13Z",
						"lineageName": "www.example.com20260714095208534842",
						"lineageType": "MULTIPLE_GENERATION",
						"sans": ["www.example.com", "example.com"],
						"secureNetwork": "ENHANCED_TLS",
						"stackMode": "MULTIPLE_STACK",
						"subject": {
							"commonName": "example.com",
							"country": "US",
							"locality": "Cambridge",
							"organization": "Example Corp.",
							"organizationalUnit": "IT",
							"state": "Massachusetts"
						}
					},
					{
						"accountId": "A-CCT1234",
						"contractId": "C-0N7RAC7",
						"geoClass": "STANDARD_WORLDWIDE",
						"groupId": 12345,
						"head": {
							"headGenerationId": 830,
							"headGenerationStatus": "CSR_READY"
						},
						"keySpecs": [
							{"keySize": "2048", "keyType": "RSA"},
							{"keySize": "P-256", "keyType": "ECDSA"}
						],
						"lineageCreatedBy": "terraform-dev",
						"lineageCreatedTime": "2026-07-16T07:37:16Z",
						"lineageId": 500003,
						"lineageModifiedBy": "terraform-dev",
						"lineageModifiedTime": "2026-07-16T07:37:16Z",
						"lineageName": "www.example.com20260716073714851272",
						"lineageType": "MULTIPLE_GENERATION",
						"sans": ["www.example.com", "example.com"],
						"secureNetwork": "ENHANCED_TLS",
						"stackMode": "MULTIPLE_STACK",
						"subject": {
							"commonName": "example.com",
							"country": "US",
							"locality": "Cambridge",
							"organization": "Example Corp.",
							"organizationalUnit": "IT",
							"state": "Massachusetts"
						}
					}
				],
				"totalCount": 2
			}`,
			expectedResponse: &ListLineagesResponse{
				Lineages: []Lineage{
					{
						AccountID:  "A-CCT1234",
						ContractID: "C-0N7RAC7",
						GeoClass:   string(GeoClassStandardWorldwide),
						GroupID:    12345,
						Head: &HeadGeneration{
							HeadGenerationID:     903,
							HeadGenerationStatus: string(GenerationStatusCSRReady),
						},
						KeySpecs: []KeySpecResponse{
							{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
							{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
						},
						LineageCreatedBy:    "terraform-dev",
						LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-14T09:52:09Z"),
						LineageID:           500002,
						LineageModifiedBy:   "terraform-dev",
						LineageModifiedTime: test.NewTimeFromString(t, "2026-07-17T07:58:13Z"),
						LineageName:         "www.example.com20260714095208534842",
						LineageType:         string(LineageTypeMultipleGeneration),
						SANs:                []string{"www.example.com", "example.com"},
						SecureNetwork:       string(SecureNetworkEnhancedTLS),
						StackMode:           string(StackModeMultipleStack),
						Subject: Subject{
							CommonName:         "example.com",
							Country:            "US",
							Locality:           "Cambridge",
							Organization:       "Example Corp.",
							OrganizationalUnit: "IT",
							State:              "Massachusetts",
						},
					},
					{
						AccountID:  "A-CCT1234",
						ContractID: "C-0N7RAC7",
						GeoClass:   string(GeoClassStandardWorldwide),
						GroupID:    12345,
						Head: &HeadGeneration{
							HeadGenerationID:     830,
							HeadGenerationStatus: string(GenerationStatusCSRReady),
						},
						KeySpecs: []KeySpecResponse{
							{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
							{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
						},
						LineageCreatedBy:    "terraform-dev",
						LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-16T07:37:16Z"),
						LineageID:           500003,
						LineageModifiedBy:   "terraform-dev",
						LineageModifiedTime: test.NewTimeFromString(t, "2026-07-16T07:37:16Z"),
						LineageName:         "www.example.com20260716073714851272",
						LineageType:         string(LineageTypeMultipleGeneration),
						SANs:                []string{"www.example.com", "example.com"},
						SecureNetwork:       string(SecureNetworkEnhancedTLS),
						StackMode:           string(StackModeMultipleStack),
						Subject: Subject{
							CommonName:         "example.com",
							Country:            "US",
							Locality:           "Cambridge",
							Organization:       "Example Corp.",
							OrganizationalUnit: "IT",
							State:              "Massachusetts",
						},
					},
				},
				TotalCount: 2,
			},
		},
		"200 OK - list lineages with filters, pagination and sort": {
			params: ListLineagesRequest{
				ContractID:        "C-0N7RAC7",
				LineageName:       "example",
				LineageIDs:        []int64{500001, 500002},
				SecureNetwork:     SecureNetworkEnhancedTLS,
				StackMode:         StackModeMultipleStack,
				LineageType:       LineageTypeMultipleGeneration,
				Domain:            "example.com",
				GenerationStatus:  []GenerationStatus{GenerationStatusReadyForUse, GenerationStatusActive},
				ExpiringInDays:    ptr.To(30),
				KeyType:           CryptographicAlgorithmRSA,
				Issuer:            "Test Certificate Authority",
				ExpandGenerations: []ExpandGenerations{ExpandGenerationsHead, ExpandGenerationsCurrentProduction},
				PageSize:          25,
				After:             "cursor-abc",
				Sort:              "-modifiedDate",
			},
			expectedPath:   "/ccm/v2/lineages?after=cursor-abc&contractId=C-0N7RAC7&domain=example.com&expandGenerations=HEAD%2CCURRENT_PRODUCTION&expiringInDays=30&generationStatus=READY_FOR_USE%2CACTIVE&issuer=Test+Certificate+Authority&keyType=RSA&lineageId=500001%2C500002&lineageName=example&lineageType=MULTIPLE_GENERATION&pageSize=25&secureNetwork=ENHANCED_TLS&sort=-modifiedDate&stackMode=MULTIPLE_STACK",
			responseStatus: http.StatusOK,
			responseBody: `{
				"lineages": [
					{
						"accountId": "A-CCT1234",
						"contractId": "C-0N7RAC7",
						"geoClass": "STANDARD_WORLDWIDE",
						"groupId": 12345,
						"head": {
							"headGenerationId": 903,
							"headGenerationStatus": "CSR_READY"
						},
						"keySpecs": [
							{"keySize": "2048", "keyType": "RSA"}
						],
						"lineageCreatedBy": "terraform-dev",
						"lineageCreatedTime": "2026-07-14T09:52:09Z",
						"lineageId": 500002,
						"lineageModifiedBy": "terraform-dev",
						"lineageModifiedTime": "2026-07-17T07:58:13Z",
						"lineageName": "www.example.com20260714095208534842",
						"lineageType": "MULTIPLE_GENERATION",
						"sans": ["www.example.com", "example.com"],
						"secureNetwork": "ENHANCED_TLS",
						"stackMode": "MULTIPLE_STACK",
						"subject": {
							"commonName": "example.com",
							"country": "US",
							"locality": "Cambridge",
							"organization": "Example Corp.",
							"organizationalUnit": "IT",
							"state": "Massachusetts"
						}
					}
				],
				"nextCursor": "eyJzb3J0IjoiLW1vZGlmaWVkRGF0ZSJ9",
				"totalCount": 1
			}`,
			expectedResponse: &ListLineagesResponse{
				Lineages: []Lineage{
					{
						AccountID:  "A-CCT1234",
						ContractID: "C-0N7RAC7",
						GeoClass:   string(GeoClassStandardWorldwide),
						GroupID:    12345,
						Head: &HeadGeneration{
							HeadGenerationID:     903,
							HeadGenerationStatus: string(GenerationStatusCSRReady),
						},
						KeySpecs: []KeySpecResponse{
							{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
						},
						LineageCreatedBy:    "terraform-dev",
						LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-14T09:52:09Z"),
						LineageID:           500002,
						LineageModifiedBy:   "terraform-dev",
						LineageModifiedTime: test.NewTimeFromString(t, "2026-07-17T07:58:13Z"),
						LineageName:         "www.example.com20260714095208534842",
						LineageType:         string(LineageTypeMultipleGeneration),
						SANs:                []string{"www.example.com", "example.com"},
						SecureNetwork:       string(SecureNetworkEnhancedTLS),
						StackMode:           string(StackModeMultipleStack),
						Subject: Subject{
							CommonName:         "example.com",
							Country:            "US",
							Locality:           "Cambridge",
							Organization:       "Example Corp.",
							OrganizationalUnit: "IT",
							State:              "Massachusetts",
						},
					},
				},
				NextCursor: ptr.To("eyJzb3J0IjoiLW1vZGlmaWVkRGF0ZSJ9"),
				TotalCount: 1,
			},
		},
		"200 OK - list lineages, no matches for lineageName filter": {
			params:         ListLineagesRequest{LineageName: "no-such-lineage"},
			expectedPath:   "/ccm/v2/lineages?lineageName=no-such-lineage",
			responseStatus: http.StatusOK,
			responseBody: `{
				"lineages": [],
				"totalCount": 0
			}`,
			expectedResponse: &ListLineagesResponse{
				Lineages:   []Lineage{},
				TotalCount: 0,
			},
		},
		"200 OK - list lineages, current production/staging/previous production populated": {
			params:         ListLineagesRequest{},
			expectedPath:   "/ccm/v2/lineages",
			responseStatus: http.StatusOK,
			responseBody: `{
				"lineages": [
					{
						"accountId": "A-CCT1234",
						"contractId": "C-0N7RAC7",
						"currentProduction": {
							"productionGenerationId": 3032,
							"productionGenerationStatus": "ACTIVE"
						},
						"currentStaging": {
							"stagingGenerationId": 3032,
							"stagingGenerationStatus": "ACTIVE"
						},
						"geoClass": "STANDARD_WORLDWIDE",
						"groupId": 12345,
						"keySpecs": [
							{"keySize": "2048", "keyType": "RSA"},
							{"keySize": "P-256", "keyType": "ECDSA"}
						],
						"lineageCreatedBy": "terraform-dev",
						"lineageCreatedTime": "2026-07-22T09:21:09Z",
						"lineageId": 500022,
						"lineageModifiedBy": "terraform-dev",
						"lineageModifiedTime": "2026-07-23T13:49:13Z",
						"lineageName": "alpha 2",
						"lineageType": "MULTIPLE_GENERATION",
						"previousProduction": {
							"previousProductionGenerationId": 3029,
							"previousProductionGenerationStatus": "READY_FOR_USE"
						},
						"sans": ["www.example.com", "example.com"],
						"secureNetwork": "ENHANCED_TLS",
						"stackMode": "MULTIPLE_STACK",
						"subject": {
							"commonName": "example.com",
							"country": "US",
							"locality": "Cambridge",
							"organization": "Example Corp.",
							"organizationalUnit": "IT",
							"state": "Massachusetts"
						}
					}
				],
				"totalCount": 1
			}`,
			expectedResponse: &ListLineagesResponse{
				Lineages: []Lineage{
					{
						AccountID:  "A-CCT1234",
						ContractID: "C-0N7RAC7",
						CurrentProduction: &ProductionGeneration{
							ProductionGenerationID:     3032,
							ProductionGenerationStatus: string(GenerationStatusActive),
						},
						CurrentStaging: &StagingGeneration{
							StagingGenerationID:     3032,
							StagingGenerationStatus: string(GenerationStatusActive),
						},
						GeoClass: string(GeoClassStandardWorldwide),
						GroupID:  12345,
						KeySpecs: []KeySpecResponse{
							{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
							{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
						},
						LineageCreatedBy:    "terraform-dev",
						LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-22T09:21:09Z"),
						LineageID:           500022,
						LineageModifiedBy:   "terraform-dev",
						LineageModifiedTime: test.NewTimeFromString(t, "2026-07-23T13:49:13Z"),
						LineageName:         "alpha 2",
						LineageType:         string(LineageTypeMultipleGeneration),
						PreviousProduction: &PreviousProductionGeneration{
							PreviousProductionGenerationID:     3029,
							PreviousProductionGenerationStatus: string(GenerationStatusReadyForUse),
						},
						SANs:          []string{"www.example.com", "example.com"},
						SecureNetwork: string(SecureNetworkEnhancedTLS),
						StackMode:     string(StackModeMultipleStack),
						Subject: Subject{
							CommonName:         "example.com",
							Country:            "US",
							Locality:           "Cambridge",
							Organization:       "Example Corp.",
							OrganizationalUnit: "IT",
							State:              "Massachusetts",
						},
					},
				},
				TotalCount: 1,
			},
		},
		"200 OK - list lineages, expand all generation pointers, all four pointers populated": {
			params: ListLineagesRequest{
				ExpandGenerations: []ExpandGenerations{
					ExpandGenerationsHead, ExpandGenerationsCurrentProduction, ExpandGenerationsCurrentStaging, ExpandGenerationsPreviousProduction,
				},
			},
			expectedPath:   "/ccm/v2/lineages?expandGenerations=HEAD%2CCURRENT_PRODUCTION%2CCURRENT_STAGING%2CPREVIOUS_PRODUCTION",
			responseStatus: http.StatusOK,
			responseBody: `{
				"lineages": [
					{
						"accountId": "A-CCT1234",
						"contractId": "C-0N7RAC7",
						"currentProduction": {
							"algorithms": [
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-23T13:44:38Z",
									"algorithmInstanceId": 4351,
									"certificateStatus": "CSR_READY",
									"csrExpirationDate": "2027-09-24T11:26:01Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "RSA"
								},
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-23T13:44:38Z",
									"algorithmInstanceId": 4352,
									"certificateStatus": "READY_FOR_USE",
									"csrExpirationDate": "2027-09-24T11:26:01Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "ECDSA",
									"signedCertificateIssuer": "CN=Test Certificate Authority",
									"signedCertificateNotValidAfterDate": "2027-07-23T13:44:09Z",
									"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n",
									"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:06",
									"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:06"
								}
							],
							"firstPromotedToProductionTime": "2026-07-23T13:49:13Z",
							"generationCreatedBy": "terraform-dev",
							"generationCreatedTime": "2026-07-23T13:44:38Z",
							"generationModifiedBy": "terraform-dev",
							"productionGenerationId": 3032,
							"productionGenerationStatus": "ACTIVE"
						},
						"currentStaging": {
							"algorithms": [
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-23T13:44:38Z",
									"algorithmInstanceId": 4351,
									"certificateStatus": "CSR_READY",
									"csrExpirationDate": "2027-09-24T11:26:01Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "RSA"
								},
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-23T13:44:38Z",
									"algorithmInstanceId": 4352,
									"certificateStatus": "READY_FOR_USE",
									"csrExpirationDate": "2027-09-24T11:26:01Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "ECDSA",
									"signedCertificateIssuer": "CN=Test Certificate Authority",
									"signedCertificateNotValidAfterDate": "2027-07-23T13:44:09Z",
									"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n",
									"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:06",
									"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:06"
								}
							],
							"firstPromotedToProductionTime": "2026-07-23T13:49:13Z",
							"generationCreatedBy": "terraform-dev",
							"generationCreatedTime": "2026-07-23T13:44:38Z",
							"generationModifiedBy": "terraform-dev",
							"stagingGenerationId": 3032,
							"stagingGenerationStatus": "ACTIVE"
						},
						"geoClass": "STANDARD_WORLDWIDE",
						"groupId": 12345,
						"head": {
							"algorithms": [
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-27T08:38:52Z",
									"algorithmInstanceId": 4544,
									"certificateStatus": "CSR_READY",
									"csrExpirationDate": "2027-09-28T08:38:52Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "RSA"
								},
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-27T08:38:52Z",
									"algorithmInstanceId": 4545,
									"certificateStatus": "CSR_READY",
									"csrExpirationDate": "2027-09-28T08:38:52Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "ECDSA"
								}
							],
							"generationCreatedBy": "terraform-dev",
							"generationCreatedTime": "2026-07-27T08:38:52Z",
							"generationModifiedBy": "terraform-dev",
							"headGenerationId": 3170,
							"headGenerationStatus": "CSR_READY"
						},
						"keySpecs": [
							{"keySize": "2048", "keyType": "RSA"},
							{"keySize": "P-256", "keyType": "ECDSA"}
						],
						"lineageCreatedBy": "terraform-dev",
						"lineageCreatedTime": "2026-07-22T09:21:09Z",
						"lineageId": 500022,
						"lineageModifiedBy": "terraform-dev",
						"lineageModifiedTime": "2026-07-27T08:38:52Z",
						"lineageName": "alpha 2",
						"lineageType": "MULTIPLE_GENERATION",
						"previousProduction": {
							"algorithms": [
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-23T11:26:01Z",
									"algorithmInstanceId": 4345,
									"algorithmInstanceModifiedBy": "terraform-dev",
									"algorithmInstanceModifiedTime": "2026-07-23T12:20:51Z",
									"certificateStatus": "READY_FOR_USE",
									"csrExpirationDate": "2027-09-24T11:26:01Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "RSA",
									"signedCertificateIssuer": "CN=Test Certificate Authority",
									"signedCertificateNotValidAfterDate": "2027-07-23T12:20:09Z",
									"signedCertificateNotValidBeforeDate": "2026-07-23T12:20:09Z",
									"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
									"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:07",
									"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:07"
								},
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-23T11:26:01Z",
									"algorithmInstanceId": 4346,
									"certificateStatus": "CSR_READY",
									"csrExpirationDate": "2027-09-24T11:26:01Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "ECDSA"
								}
							],
							"firstPromotedToProductionTime": "2026-07-23T12:21:36Z",
							"generationCreatedBy": "terraform-dev",
							"generationCreatedTime": "2026-07-23T11:26:01Z",
							"generationModifiedBy": "terraform-dev",
							"generationModifiedTime": "2026-07-23T12:20:52Z",
							"previousProductionGenerationId": 3029,
							"previousProductionGenerationStatus": "READY_FOR_USE"
						},
						"sans": ["www.example.com", "example.com"],
						"secureNetwork": "ENHANCED_TLS",
						"stackMode": "MULTIPLE_STACK",
						"subject": {
							"commonName": "example.com",
							"country": "US",
							"locality": "Cambridge",
							"organization": "Example Corp.",
							"organizationalUnit": "IT",
							"state": "Massachusetts"
						}
					}
				],
				"totalCount": 1
			}`,
			expectedResponse: &ListLineagesResponse{
				Lineages: []Lineage{
					{
						AccountID:  "A-CCT1234",
						ContractID: "C-0N7RAC7",
						CurrentProduction: &ProductionGeneration{
							Generation: Generation{
								Algorithms: []Algorithm{
									{
										AlgorithmInstanceCreatedBy:   "terraform-dev",
										AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T13:44:38Z")),
										AlgorithmInstanceID:          4351,
										CertificateStatus:            CertificateStatusCSRReady,
										CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-24T11:26:01Z")),
										CSRPEM:                       rsaCSRPEM,
										KeyType:                      string(CryptographicAlgorithmRSA),
									},
									{
										AlgorithmInstanceCreatedBy:         "terraform-dev",
										AlgorithmInstanceCreatedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-23T13:44:38Z")),
										AlgorithmInstanceID:                4352,
										CertificateStatus:                  CertificateStatusReadyForUse,
										CSRExpirationDate:                  ptr.To(test.NewTimeFromString(t, "2027-09-24T11:26:01Z")),
										CSRPEM:                             ecdsaCSRPEM,
										KeyType:                            string(CryptographicAlgorithmECDSA),
										SignedCertificateIssuer:            ptr.To("CN=Test Certificate Authority"),
										SignedCertificateNotValidAfterDate: ptr.To(test.NewTimeFromString(t, "2027-07-23T13:44:09Z")),
										SignedCertificatePEM:               ptr.To(ecdsaCertPEM),
										SignedCertificateSerialNumber:      ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:06"),
										SignedCertificateSHA256Fingerprint: ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:06"),
									},
								},
								FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T13:49:13Z")),
								GenerationCreatedBy:           ptr.To("terraform-dev"),
								GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-23T13:44:38Z")),
								GenerationModifiedBy:          ptr.To("terraform-dev"),
							},
							ProductionGenerationID:     3032,
							ProductionGenerationStatus: string(GenerationStatusActive),
						},
						CurrentStaging: &StagingGeneration{
							Generation: Generation{
								Algorithms: []Algorithm{
									{
										AlgorithmInstanceCreatedBy:   "terraform-dev",
										AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T13:44:38Z")),
										AlgorithmInstanceID:          4351,
										CertificateStatus:            CertificateStatusCSRReady,
										CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-24T11:26:01Z")),
										CSRPEM:                       rsaCSRPEM,
										KeyType:                      string(CryptographicAlgorithmRSA),
									},
									{
										AlgorithmInstanceCreatedBy:         "terraform-dev",
										AlgorithmInstanceCreatedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-23T13:44:38Z")),
										AlgorithmInstanceID:                4352,
										CertificateStatus:                  CertificateStatusReadyForUse,
										CSRExpirationDate:                  ptr.To(test.NewTimeFromString(t, "2027-09-24T11:26:01Z")),
										CSRPEM:                             ecdsaCSRPEM,
										KeyType:                            string(CryptographicAlgorithmECDSA),
										SignedCertificateIssuer:            ptr.To("CN=Test Certificate Authority"),
										SignedCertificateNotValidAfterDate: ptr.To(test.NewTimeFromString(t, "2027-07-23T13:44:09Z")),
										SignedCertificatePEM:               ptr.To(ecdsaCertPEM),
										SignedCertificateSerialNumber:      ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:06"),
										SignedCertificateSHA256Fingerprint: ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:06"),
									},
								},
								FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T13:49:13Z")),
								GenerationCreatedBy:           ptr.To("terraform-dev"),
								GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-23T13:44:38Z")),
								GenerationModifiedBy:          ptr.To("terraform-dev"),
							},
							StagingGenerationID:     3032,
							StagingGenerationStatus: string(GenerationStatusActive),
						},
						GeoClass: string(GeoClassStandardWorldwide),
						GroupID:  12345,
						Head: &HeadGeneration{
							Generation: Generation{
								Algorithms: []Algorithm{
									{
										AlgorithmInstanceCreatedBy:   "terraform-dev",
										AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T08:38:52Z")),
										AlgorithmInstanceID:          4544,
										CertificateStatus:            CertificateStatusCSRReady,
										CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T08:38:52Z")),
										CSRPEM:                       rsaCSRPEM,
										KeyType:                      string(CryptographicAlgorithmRSA),
									},
									{
										AlgorithmInstanceCreatedBy:   "terraform-dev",
										AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T08:38:52Z")),
										AlgorithmInstanceID:          4545,
										CertificateStatus:            CertificateStatusCSRReady,
										CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T08:38:52Z")),
										CSRPEM:                       ecdsaCSRPEM,
										KeyType:                      string(CryptographicAlgorithmECDSA),
									},
								},
								GenerationCreatedBy:   ptr.To("terraform-dev"),
								GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T08:38:52Z")),
								GenerationModifiedBy:  ptr.To("terraform-dev"),
							},
							HeadGenerationID:     3170,
							HeadGenerationStatus: string(GenerationStatusCSRReady),
						},
						KeySpecs: []KeySpecResponse{
							{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
							{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
						},
						LineageCreatedBy:    "terraform-dev",
						LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-22T09:21:09Z"),
						LineageID:           500022,
						LineageModifiedBy:   "terraform-dev",
						LineageModifiedTime: test.NewTimeFromString(t, "2026-07-27T08:38:52Z"),
						LineageName:         "alpha 2",
						LineageType:         string(LineageTypeMultipleGeneration),
						PreviousProduction: &PreviousProductionGeneration{
							Generation: Generation{
								Algorithms: []Algorithm{
									{
										AlgorithmInstanceCreatedBy:          "terraform-dev",
										AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:26:01Z")),
										AlgorithmInstanceID:                 4345,
										AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
										AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-23T12:20:51Z")),
										CertificateStatus:                   CertificateStatusReadyForUse,
										CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-24T11:26:01Z")),
										CSRPEM:                              rsaCSRPEM,
										KeyType:                             string(CryptographicAlgorithmRSA),
										SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
										SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-23T12:20:09Z")),
										SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-23T12:20:09Z")),
										SignedCertificatePEM:                ptr.To(rsaCertPEM),
										SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:07"),
										SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:07"),
									},
									{
										AlgorithmInstanceCreatedBy:   "terraform-dev",
										AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:26:01Z")),
										AlgorithmInstanceID:          4346,
										CertificateStatus:            CertificateStatusCSRReady,
										CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-24T11:26:01Z")),
										CSRPEM:                       ecdsaCSRPEM,
										KeyType:                      string(CryptographicAlgorithmECDSA),
									},
								},
								FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T12:21:36Z")),
								GenerationCreatedBy:           ptr.To("terraform-dev"),
								GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-23T11:26:01Z")),
								GenerationModifiedBy:          ptr.To("terraform-dev"),
								GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T12:20:52Z")),
							},
							PreviousProductionGenerationID:     3029,
							PreviousProductionGenerationStatus: string(GenerationStatusReadyForUse),
						},
						SANs:          []string{"www.example.com", "example.com"},
						SecureNetwork: string(SecureNetworkEnhancedTLS),
						StackMode:     string(StackModeMultipleStack),
						Subject: Subject{
							CommonName:         "example.com",
							Country:            "US",
							Locality:           "Cambridge",
							Organization:       "Example Corp.",
							OrganizationalUnit: "IT",
							State:              "Massachusetts",
						},
					},
				},
				TotalCount: 1,
			},
		},
		"200 OK - list lineages, expand all generation pointers, single match": {
			params: ListLineagesRequest{
				LineageName: "example20260717084521257580",
				ExpandGenerations: []ExpandGenerations{
					ExpandGenerationsHead, ExpandGenerationsCurrentProduction, ExpandGenerationsCurrentStaging, ExpandGenerationsPreviousProduction,
				},
			},
			expectedPath:   "/ccm/v2/lineages?expandGenerations=HEAD%2CCURRENT_PRODUCTION%2CCURRENT_STAGING%2CPREVIOUS_PRODUCTION&lineageName=example20260717084521257580",
			responseStatus: http.StatusOK,
			responseBody: `{
				"lineages": [
					{
						"accountId": "A-CCT1234",
						"contractId": "C-0N7RAC7",
						"geoClass": "STANDARD_WORLDWIDE",
						"groupId": 12345,
						"head": {
							"algorithms": [
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-17T08:45:22Z",
									"algorithmInstanceId": 1674,
									"certificateStatus": "CSR_READY",
									"csrExpirationDate": "2027-07-17T08:45:21Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "RSA"
								},
								{
									"algorithmInstanceCreatedBy": "terraform-dev",
									"algorithmInstanceCreatedTime": "2026-07-17T08:45:22Z",
									"algorithmInstanceId": 1675,
									"certificateStatus": "CSR_READY",
									"csrExpirationDate": "2027-07-17T08:45:21Z",
									"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
									"keyType": "ECDSA"
								}
							],
							"generationCreatedBy": "terraform-dev",
							"generationCreatedTime": "2026-07-17T08:45:22Z",
							"generationModifiedBy": "terraform-dev",
							"headGenerationId": 909,
							"headGenerationStatus": "CSR_READY"
						},
						"keySpecs": [
							{"keySize": "2048", "keyType": "RSA"},
							{"keySize": "P-256", "keyType": "ECDSA"}
						],
						"lineageCreatedBy": "terraform-dev",
						"lineageCreatedTime": "2026-07-17T08:45:22Z",
						"lineageId": 500004,
						"lineageModifiedBy": "terraform-dev",
						"lineageModifiedTime": "2026-07-17T08:45:22Z",
						"lineageName": "www.example.com20260717084521257580",
						"lineageType": "MULTIPLE_GENERATION",
						"sans": ["www.example.com", "example.com"],
						"secureNetwork": "ENHANCED_TLS",
						"stackMode": "MULTIPLE_STACK",
						"subject": {
							"commonName": "example.com",
							"country": "US",
							"locality": "Cambridge",
							"organization": "Example Corp.",
							"organizationalUnit": "IT",
							"state": "Massachusetts"
						}
					}
				],
				"totalCount": 1
			}`,
			expectedResponse: &ListLineagesResponse{
				Lineages: []Lineage{
					{
						AccountID:  "A-CCT1234",
						ContractID: "C-0N7RAC7",
						GeoClass:   string(GeoClassStandardWorldwide),
						GroupID:    12345,
						Head: &HeadGeneration{
							Generation: Generation{
								Algorithms: []Algorithm{
									{
										AlgorithmInstanceCreatedBy:   "terraform-dev",
										AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-17T08:45:22Z")),
										AlgorithmInstanceID:          1674,
										CertificateStatus:            CertificateStatusCSRReady,
										CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-07-17T08:45:21Z")),
										CSRPEM:                       rsaCSRPEM,
										KeyType:                      string(CryptographicAlgorithmRSA),
									},
									{
										AlgorithmInstanceCreatedBy:   "terraform-dev",
										AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-17T08:45:22Z")),
										AlgorithmInstanceID:          1675,
										CertificateStatus:            CertificateStatusCSRReady,
										CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-07-17T08:45:21Z")),
										CSRPEM:                       ecdsaCSRPEM,
										KeyType:                      string(CryptographicAlgorithmECDSA),
									},
								},
								GenerationCreatedBy:   ptr.To("terraform-dev"),
								GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-17T08:45:22Z")),
								GenerationModifiedBy:  ptr.To("terraform-dev"),
							},
							HeadGenerationID:     909,
							HeadGenerationStatus: string(GenerationStatusCSRReady),
						},
						KeySpecs: []KeySpecResponse{
							{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
							{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
						},
						LineageCreatedBy:    "terraform-dev",
						LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-17T08:45:22Z"),
						LineageID:           500004,
						LineageModifiedBy:   "terraform-dev",
						LineageModifiedTime: test.NewTimeFromString(t, "2026-07-17T08:45:22Z"),
						LineageName:         "www.example.com20260717084521257580",
						LineageType:         string(LineageTypeMultipleGeneration),
						SANs:                []string{"www.example.com", "example.com"},
						SecureNetwork:       string(SecureNetworkEnhancedTLS),
						StackMode:           string(StackModeMultipleStack),
						Subject: Subject{
							CommonName:         "example.com",
							Country:            "US",
							Locality:           "Cambridge",
							Organization:       "Example Corp.",
							OrganizationalUnit: "IT",
							State:              "Massachusetts",
						},
					},
				},
				TotalCount: 1,
			},
		},
		"500 internal server error": {
			params:         ListLineagesRequest{},
			expectedPath:   "/ccm/v2/lineages",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891017"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineages, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891017",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrListLineages)
			},
		},
		"400 invalid or expired after cursor": {
			params:         ListLineagesRequest{After: "invalid-or-expired-cursor"},
			expectedPath:   "/ccm/v2/lineages?after=invalid-or-expired-cursor",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/invalid-cursor",
				"title": "Invalid or expired pagination cursor.",
				"status": 400,
				"detail": "The 'after' cursor value is invalid or has expired. Please restart pagination without a cursor.",
				"instance": "/error-types/invalid-cursor?traceId=1234567891018"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineages, &Error{
					Type:     "/error-types/invalid-cursor",
					Title:    "Invalid or expired pagination cursor.",
					Status:   http.StatusBadRequest,
					Detail:   "The 'after' cursor value is invalid or has expired. Please restart pagination without a cursor.",
					Instance: "/error-types/invalid-cursor?traceId=1234567891018",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInvalidCursor)
				assert.ErrorIs(t, err, ErrListLineages)
			},
		},
		"400 invalid sort parameter": {
			params:         ListLineagesRequest{Sort: "foo"},
			expectedPath:   "/ccm/v2/lineages?sort=foo",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/invalid-sort-parameter",
				"title": "Invalid sort parameter.",
				"status": 400,
				"detail": "Sort parameter '{foo}' is not valid. {Valid values: -modifiedDate, +modifiedDate, -createdDate, +createdDate, +lineageName, -lineageName, +expirationDate, -expirationDate}",
				"instance": "/error-types/invalid-sort-parameter?traceId=1234567891019",
				"context": {
					"explanation": "Valid values: -modifiedDate, +modifiedDate, -createdDate, +createdDate, +lineageName, -lineageName, +expirationDate, -expirationDate",
					"invalidParameterValue": "foo",
					"parameterName": "sort"
				}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineages, &Error{
					Type:     "/error-types/invalid-sort-parameter",
					Title:    "Invalid sort parameter.",
					Status:   http.StatusBadRequest,
					Detail:   "Sort parameter '{foo}' is not valid. {Valid values: -modifiedDate, +modifiedDate, -createdDate, +createdDate, +lineageName, -lineageName, +expirationDate, -expirationDate}",
					Instance: "/error-types/invalid-sort-parameter?traceId=1234567891019",
					Context: map[string]any{
						"explanation":           "Valid values: -modifiedDate, +modifiedDate, -createdDate, +createdDate, +lineageName, -lineageName, +expirationDate, -expirationDate",
						"invalidParameterValue": "foo",
						"parameterName":         "sort",
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInvalidSortParameter)
				assert.ErrorIs(t, err, ErrListLineages)
			},
		},
		"validation error - invalid SecureNetwork": {
			params: ListLineagesRequest{SecureNetwork: "NOT_A_REAL_NETWORK"},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineages: struct validation: SecureNetwork: value "+
					"'NOT_A_REAL_NETWORK' is invalid. Must be either 'ENHANCED_TLS' or 'STANDARD_TLS'")
				assert.ErrorIs(t, err, ErrListLineages)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid StackMode": {
			params: ListLineagesRequest{StackMode: "NOT_A_REAL_MODE"},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineages: struct validation: StackMode: value "+
					"'NOT_A_REAL_MODE' is invalid. Must be either 'SINGLE_STACK' or 'MULTIPLE_STACK'")
				assert.ErrorIs(t, err, ErrListLineages)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid LineageType": {
			params: ListLineagesRequest{LineageType: "NOT_A_REAL_TYPE"},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineages: struct validation: LineageType: value "+
					"'NOT_A_REAL_TYPE' is invalid. Must be either 'MULTIPLE_GENERATION' or 'SINGLE_GENERATION'")
				assert.ErrorIs(t, err, ErrListLineages)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid GenerationStatus": {
			params: ListLineagesRequest{GenerationStatus: []GenerationStatus{"NOT_A_REAL_STATUS"}},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineages: struct validation: 0: value "+
					"'NOT_A_REAL_STATUS' is invalid. Must be one of: 'CSR_READY', 'READY_FOR_USE', 'ACTIVE', "+
					"'ARCHIVED', or 'ABANDONED'")
				assert.ErrorIs(t, err, ErrListLineages)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - empty LineageIDs": {
			params: ListLineagesRequest{LineageIDs: []int64{}},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineages: struct validation: LineageIDs: cannot be blank")
				assert.ErrorIs(t, err, ErrListLineages)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid KeyType": {
			params: ListLineagesRequest{KeyType: "NOT_A_REAL_ALGORITHM"},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineages: struct validation: KeyType: value "+
					"'NOT_A_REAL_ALGORITHM' is invalid. Must be either 'RSA' or 'ECDSA'")
				assert.ErrorIs(t, err, ErrListLineages)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid ExpandGenerations": {
			params: ListLineagesRequest{ExpandGenerations: []ExpandGenerations{"ALL"}},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineages: struct validation: 0: value 'ALL' is invalid. "+
					"Must be one of: 'HEAD', 'CURRENT_PRODUCTION', 'CURRENT_STAGING', or 'PREVIOUS_PRODUCTION'")
				assert.ErrorIs(t, err, ErrListLineages)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - PageSize exceeds maximum": {
			params: ListLineagesRequest{PageSize: 101},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, fmt.Sprintf("listing lineages: struct validation: PageSize: must be no greater than %d", MaxListLineagesPageSize))
				assert.ErrorIs(t, err, ErrListLineages)
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
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.ListLineages(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestRenameLineage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params              RenameLineageRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedResponse    *RenameLineageResponse
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"200 OK - rename lineage": {
			params:              RenameLineageRequest{LineageID: 500001, LineageName: "renamed-lineage"},
			expectedPath:        "/ccm/v2/lineages/500001",
			expectedRequestBody: `[{"op":"replace","path":"/lineageName","value":"renamed-lineage"}]`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"headGenerationId": 2912,
					"headGenerationStatus": "READY_FOR_USE"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-07T10:18:46Z",
				"lineageId": 500001,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-08T13:24:17Z",
				"lineageName": "renamed-lineage",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "SINGLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Corp.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &RenameLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					HeadGenerationID:     2912,
					HeadGenerationStatus: string(GenerationStatusReadyForUse),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-07T10:18:46Z"),
				LineageID:           500001,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-08T13:24:17Z"),
				LineageName:         "renamed-lineage",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeSingleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Corp.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"500 internal server error": {
			params:              RenameLineageRequest{LineageID: 500001, LineageName: "renamed-lineage"},
			expectedPath:        "/ccm/v2/lineages/500001",
			expectedRequestBody: `[{"op":"replace","path":"/lineageName","value":"renamed-lineage"}]`,
			responseStatus:      http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891020"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRenameLineage, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891020",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrRenameLineage)
			},
		},
		"404 lineage not found": {
			params:              RenameLineageRequest{LineageID: 999999, LineageName: "renamed-lineage"},
			expectedPath:        "/ccm/v2/lineages/999999",
			expectedRequestBody: `[{"op":"replace","path":"/lineageName","value":"renamed-lineage"}]`,
			responseStatus:      http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891024",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRenameLineage, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891024",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrRenameLineage)
			},
		},
		"409 lineage name conflict": {
			params:              RenameLineageRequest{LineageID: 500001, LineageName: "renamed-lineage"},
			expectedPath:        "/ccm/v2/lineages/500001",
			expectedRequestBody: `[{"op":"replace","path":"/lineageName","value":"renamed-lineage"}]`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/lineage-name-conflict",
				"title": "A lineage with this name already exists for this account.",
				"status": 409,
				"detail": "A lineage named '{renamed-lineage}' already exists for account {A-CCT1234}.",
				"instance": "/error-types/lineage-name-conflict?traceId=1234567891025",
				"context": {"accountId": "A-CCT1234", "lineageName": "renamed-lineage"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRenameLineage, &Error{
					Type:     "/error-types/lineage-name-conflict",
					Title:    "A lineage with this name already exists for this account.",
					Status:   http.StatusConflict,
					Detail:   "A lineage named '{renamed-lineage}' already exists for account {A-CCT1234}.",
					Instance: "/error-types/lineage-name-conflict?traceId=1234567891025",
					Context:  map[string]any{"accountId": "A-CCT1234", "lineageName": "renamed-lineage"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNameConflict)
				assert.ErrorIs(t, err, ErrRenameLineage)
			},
		},
		"validation error - missing LineageID": {
			params: RenameLineageRequest{LineageName: "renamed-lineage"},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "renaming lineage: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrRenameLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing LineageName": {
			params: RenameLineageRequest{LineageID: 500001},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "renaming lineage: struct validation: LineageName: cannot be blank")
				assert.ErrorIs(t, err, ErrRenameLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - LineageName too long": {
			params: RenameLineageRequest{LineageID: 500001, LineageName: strings.Repeat("a", 271)},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "renaming lineage: struct validation: LineageName: the length must be between 1 and 270")
				assert.ErrorIs(t, err, ErrRenameLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPatch, r.Method)
				if tc.expectedRequestBody != "" {
					requestBody, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(requestBody))
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.RenameLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestDeleteLineage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params         DeleteLineageRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No Content - delete lineage": {
			params:         DeleteLineageRequest{LineageID: 500003},
			expectedPath:   "/ccm/v2/lineages/500003",
			responseStatus: http.StatusNoContent,
		},
		"404 lineage not found": {
			params:         DeleteLineageRequest{LineageID: 999999},
			expectedPath:   "/ccm/v2/lineages/999999",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891037",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteLineage, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891037",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrDeleteLineage)
			},
		},
		"500 internal server error": {
			params:         DeleteLineageRequest{LineageID: 500003},
			expectedPath:   "/ccm/v2/lineages/500003",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891038"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteLineage, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891038",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrDeleteLineage)
			},
		},
		"409 lineage has active production": {
			params:         DeleteLineageRequest{LineageID: 500022},
			expectedPath:   "/ccm/v2/lineages/500022",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/lineage-has-active-production",
				"title": "Cannot delete lineage with active production certificate.",
				"status": 409,
				"detail": "Lineage {500022} has an active production certificate. Deactivate before deleting.",
				"instance": "/error-types/lineage-has-active-production?traceId=1234567891039",
				"context": {"lineageId": 500022}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteLineage, &Error{
					Type:     "/error-types/lineage-has-active-production",
					Title:    "Cannot delete lineage with active production certificate.",
					Status:   http.StatusConflict,
					Detail:   "Lineage {500022} has an active production certificate. Deactivate before deleting.",
					Instance: "/error-types/lineage-has-active-production?traceId=1234567891039",
					Context:  map[string]any{"lineageId": float64(500022)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageHasActiveProduction)
				assert.ErrorIs(t, err, ErrDeleteLineage)
			},
		},
		"409 lineage has active staging": {
			params:         DeleteLineageRequest{LineageID: 500009},
			expectedPath:   "/ccm/v2/lineages/500009",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/lineage-has-active-staging",
				"title": "Cannot delete lineage with active staging certificate.",
				"status": 409,
				"detail": "Lineage {500009} has an active staging certificate. Deactivate before deleting.",
				"instance": "/error-types/lineage-has-active-staging?traceId=1234567891040",
				"context": {"lineageId": 500009}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteLineage, &Error{
					Type:     "/error-types/lineage-has-active-staging",
					Title:    "Cannot delete lineage with active staging certificate.",
					Status:   http.StatusConflict,
					Detail:   "Lineage {500009} has an active staging certificate. Deactivate before deleting.",
					Instance: "/error-types/lineage-has-active-staging?traceId=1234567891040",
					Context:  map[string]any{"lineageId": float64(500009)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageHasActiveStaging)
				assert.ErrorIs(t, err, ErrDeleteLineage)
			},
		},
		"validation error - missing LineageID": {
			params: DeleteLineageRequest{},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "deleting lineage: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrDeleteLineage)
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
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			err := client.DeleteLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestListLineageActivity(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           ListLineageActivityRequest
		responseStatus   int
		responseBody     string
		expectedResponse *ListLineageActivityResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"200 OK - list lineage activity": {
			params:         ListLineageActivityRequest{LineageID: 500005},
			expectedPath:   "/ccm/v2/lineages/500005/activity",
			responseStatus: http.StatusOK,
			responseBody: `{
				"events": [
					{
						"activityId": 3242,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-21T12:18:03Z",
						"eventType": "LINEAGE_RENAMED",
						"lineageId": 500005
					},
					{
						"activityId": 3194,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-21T11:52:21Z",
						"eventType": "LINEAGE_RENAMED",
						"lineageId": 500005
					},
					{
						"activityId": 3193,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-21T11:52:07Z",
						"eventType": "LINEAGE_RENAMED",
						"lineageId": 500005
					},
					{
						"activityId": 3192,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-21T11:50:30Z",
						"eventType": "LINEAGE_RENAMED",
						"lineageId": 500005
					},
					{
						"activityId": 2616,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-20T15:19:16Z",
						"eventType": "LINEAGE_RENAMED",
						"lineageId": 500005
					},
					{
						"activityId": 1805,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-17T12:50:44Z",
						"eventType": "CERT_UPLOADED",
						"generationId": 929,
						"lineageId": 500005
					},
					{
						"activityId": 1783,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-17T12:28:59Z",
						"eventType": "LINEAGE_CREATED",
						"generationId": 929,
						"lineageId": 500005
					}
				],
				"totalCount": 7
			}`,
			expectedResponse: &ListLineageActivityResponse{
				Events: []LineageActivityEvent{
					{
						ActivityID:  3242,
						CreatedBy:   "terraform-dev",
						CreatedTime: test.NewTimeFromString(t, "2026-07-21T12:18:03Z"),
						EventType:   ActivityEventTypeLineageRenamed,
						LineageID:   500005,
					},
					{
						ActivityID:  3194,
						CreatedBy:   "terraform-dev",
						CreatedTime: test.NewTimeFromString(t, "2026-07-21T11:52:21Z"),
						EventType:   ActivityEventTypeLineageRenamed,
						LineageID:   500005,
					},
					{
						ActivityID:  3193,
						CreatedBy:   "terraform-dev",
						CreatedTime: test.NewTimeFromString(t, "2026-07-21T11:52:07Z"),
						EventType:   ActivityEventTypeLineageRenamed,
						LineageID:   500005,
					},
					{
						ActivityID:  3192,
						CreatedBy:   "terraform-dev",
						CreatedTime: test.NewTimeFromString(t, "2026-07-21T11:50:30Z"),
						EventType:   ActivityEventTypeLineageRenamed,
						LineageID:   500005,
					},
					{
						ActivityID:  2616,
						CreatedBy:   "terraform-dev",
						CreatedTime: test.NewTimeFromString(t, "2026-07-20T15:19:16Z"),
						EventType:   ActivityEventTypeLineageRenamed,
						LineageID:   500005,
					},
					{
						ActivityID:   1805,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-17T12:50:44Z"),
						EventType:    ActivityEventTypeCertUploaded,
						GenerationID: ptr.To(int64(929)),
						LineageID:    500005,
					},
					{
						ActivityID:   1783,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-17T12:28:59Z"),
						EventType:    ActivityEventTypeLineageCreated,
						GenerationID: ptr.To(int64(929)),
						LineageID:    500005,
					},
				},
				TotalCount: 7,
			},
		},
		"200 OK - list lineage activity with pagination": {
			params:         ListLineageActivityRequest{LineageID: 500005, PageSize: 2},
			expectedPath:   "/ccm/v2/lineages/500005/activity?pageSize=2",
			responseStatus: http.StatusOK,
			responseBody: `{
				"events": [
					{
						"activityId": 3242,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-21T12:18:03Z",
						"eventType": "LINEAGE_RENAMED",
						"lineageId": 500005
					},
					{
						"activityId": 3194,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-21T11:52:21Z",
						"eventType": "LINEAGE_RENAMED",
						"lineageId": 500005
					}
				],
				"nextCursor": "MzE5NA==",
				"totalCount": 7
			}`,
			expectedResponse: &ListLineageActivityResponse{
				Events: []LineageActivityEvent{
					{
						ActivityID:  3242,
						CreatedBy:   "terraform-dev",
						CreatedTime: test.NewTimeFromString(t, "2026-07-21T12:18:03Z"),
						EventType:   ActivityEventTypeLineageRenamed,
						LineageID:   500005,
					},
					{
						ActivityID:  3194,
						CreatedBy:   "terraform-dev",
						CreatedTime: test.NewTimeFromString(t, "2026-07-21T11:52:21Z"),
						EventType:   ActivityEventTypeLineageRenamed,
						LineageID:   500005,
					},
				},
				NextCursor: ptr.To("MzE5NA=="),
				TotalCount: 7,
			},
		},
		"200 OK - list lineage activity, full lifecycle with network/outcome and generation event types": {
			params:         ListLineageActivityRequest{LineageID: 500022},
			expectedPath:   "/ccm/v2/lineages/500022/activity",
			responseStatus: http.StatusOK,
			responseBody: `{
				"events": [
					{
						"activityId": 5245,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T13:42:30Z",
						"eventType": "GENERATION_ROLLED_BACK",
						"generationId": 3029,
						"lineageId": 500022,
						"network": "PRODUCTION",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 5242,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T13:22:55Z",
						"eventType": "GENERATION_DELETED",
						"generationId": 3031,
						"lineageId": 500022
					},
					{
						"activityId": 5241,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T13:03:24Z",
						"eventType": "STAGING_REPLACED",
						"generationId": 3028,
						"lineageId": 500022,
						"network": "STAGING",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 5240,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T12:52:26Z",
						"eventType": "GENERATION_PROMOTED_TO_STAGING",
						"generationId": 3031,
						"lineageId": 500022,
						"network": "STAGING",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 5239,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T12:51:53Z",
						"eventType": "CERT_UPLOADED",
						"generationId": 3031,
						"lineageId": 500022
					},
					{
						"activityId": 5238,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T12:50:34Z",
						"eventType": "GENERATION_RENEWED",
						"generationId": 3031,
						"lineageId": 500022
					},
					{
						"activityId": 5237,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T12:46:26Z",
						"eventType": "GENERATION_ROLLED_BACK",
						"generationId": 3028,
						"lineageId": 500022,
						"network": "PRODUCTION",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 5236,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T12:21:36Z",
						"eventType": "GENERATION_PROMOTED_TO_PRODUCTION",
						"generationId": 3029,
						"lineageId": 500022,
						"network": "PRODUCTION",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 5235,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T12:20:52Z",
						"eventType": "CERT_UPLOADED",
						"generationId": 3029,
						"lineageId": 500022
					},
					{
						"activityId": 5231,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:26:01Z",
						"eventType": "GENERATION_RENEWED",
						"generationId": 3029,
						"lineageId": 500022
					},
					{
						"activityId": 5230,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:25:05Z",
						"eventType": "GENERATION_PROMOTED_TO_PRODUCTION",
						"generationId": 3028,
						"lineageId": 500022,
						"network": "PRODUCTION",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 5229,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:25:05Z",
						"eventType": "GENERATION_PROMOTED_TO_STAGING",
						"generationId": 3028,
						"lineageId": 500022,
						"network": "STAGING",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 5226,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:21:42Z",
						"eventType": "CERT_UPLOADED",
						"generationId": 3028,
						"lineageId": 500022
					},
					{
						"activityId": 5225,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:20:11Z",
						"eventType": "GENERATION_RENEWED",
						"generationId": 3028,
						"lineageId": 500022
					},
					{
						"activityId": 5224,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:18:53Z",
						"eventType": "GENERATION_PROMOTED_TO_PRODUCTION",
						"generationId": 3027,
						"lineageId": 500022,
						"network": "PRODUCTION",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 5223,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:16:01Z",
						"eventType": "CERT_UPLOADED",
						"generationId": 3027,
						"lineageId": 500022
					},
					{
						"activityId": 5222,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:13:41Z",
						"eventType": "GENERATION_RENEWED",
						"generationId": 3027,
						"lineageId": 500022
					},
					{
						"activityId": 5221,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-23T11:08:24Z",
						"eventType": "GENERATION_PROMOTED_TO_PRODUCTION",
						"generationId": 2448,
						"lineageId": 500022,
						"network": "PRODUCTION",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 4540,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-22T12:47:56Z",
						"eventType": "GENERATION_PROMOTED_TO_STAGING",
						"generationId": 2448,
						"lineageId": 500022,
						"network": "STAGING",
						"outcome": "ALL_SUCCESS"
					},
					{
						"activityId": 4296,
						"createdBy": "terraform-dev",
						"createdTime": "2026-07-22T12:31:00Z",
						"eventType": "CERT_UPLOADED",
						"generationId": 2448,
						"lineageId": 500022
					}
				],
				"nextCursor": "NDI5Ng==",
				"totalCount": 22
			}`,
			expectedResponse: &ListLineageActivityResponse{
				Events: []LineageActivityEvent{
					{
						ActivityID:   5245,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T13:42:30Z"),
						EventType:    ActivityEventTypeGenerationRolledBack,
						GenerationID: ptr.To(int64(3029)),
						LineageID:    500022,
						Network:      ptr.To("PRODUCTION"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   5242,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T13:22:55Z"),
						EventType:    ActivityEventTypeGenerationDeleted,
						GenerationID: ptr.To(int64(3031)),
						LineageID:    500022,
					},
					{
						ActivityID:   5241,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T13:03:24Z"),
						EventType:    ActivityEventTypeStagingReplaced,
						GenerationID: ptr.To(int64(3028)),
						LineageID:    500022,
						Network:      ptr.To("STAGING"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   5240,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T12:52:26Z"),
						EventType:    ActivityEventTypeGenerationPromotedToStaging,
						GenerationID: ptr.To(int64(3031)),
						LineageID:    500022,
						Network:      ptr.To("STAGING"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   5239,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T12:51:53Z"),
						EventType:    ActivityEventTypeCertUploaded,
						GenerationID: ptr.To(int64(3031)),
						LineageID:    500022,
					},
					{
						ActivityID:   5238,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T12:50:34Z"),
						EventType:    ActivityEventTypeGenerationRenewed,
						GenerationID: ptr.To(int64(3031)),
						LineageID:    500022,
					},
					{
						ActivityID:   5237,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T12:46:26Z"),
						EventType:    ActivityEventTypeGenerationRolledBack,
						GenerationID: ptr.To(int64(3028)),
						LineageID:    500022,
						Network:      ptr.To("PRODUCTION"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   5236,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T12:21:36Z"),
						EventType:    ActivityEventTypeGenerationPromotedToProduction,
						GenerationID: ptr.To(int64(3029)),
						LineageID:    500022,
						Network:      ptr.To("PRODUCTION"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   5235,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T12:20:52Z"),
						EventType:    ActivityEventTypeCertUploaded,
						GenerationID: ptr.To(int64(3029)),
						LineageID:    500022,
					},
					{
						ActivityID:   5231,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:26:01Z"),
						EventType:    ActivityEventTypeGenerationRenewed,
						GenerationID: ptr.To(int64(3029)),
						LineageID:    500022,
					},
					{
						ActivityID:   5230,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:25:05Z"),
						EventType:    ActivityEventTypeGenerationPromotedToProduction,
						GenerationID: ptr.To(int64(3028)),
						LineageID:    500022,
						Network:      ptr.To("PRODUCTION"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   5229,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:25:05Z"),
						EventType:    ActivityEventTypeGenerationPromotedToStaging,
						GenerationID: ptr.To(int64(3028)),
						LineageID:    500022,
						Network:      ptr.To("STAGING"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   5226,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:21:42Z"),
						EventType:    ActivityEventTypeCertUploaded,
						GenerationID: ptr.To(int64(3028)),
						LineageID:    500022,
					},
					{
						ActivityID:   5225,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:20:11Z"),
						EventType:    ActivityEventTypeGenerationRenewed,
						GenerationID: ptr.To(int64(3028)),
						LineageID:    500022,
					},
					{
						ActivityID:   5224,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:18:53Z"),
						EventType:    ActivityEventTypeGenerationPromotedToProduction,
						GenerationID: ptr.To(int64(3027)),
						LineageID:    500022,
						Network:      ptr.To("PRODUCTION"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   5223,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:16:01Z"),
						EventType:    ActivityEventTypeCertUploaded,
						GenerationID: ptr.To(int64(3027)),
						LineageID:    500022,
					},
					{
						ActivityID:   5222,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:13:41Z"),
						EventType:    ActivityEventTypeGenerationRenewed,
						GenerationID: ptr.To(int64(3027)),
						LineageID:    500022,
					},
					{
						ActivityID:   5221,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:08:24Z"),
						EventType:    ActivityEventTypeGenerationPromotedToProduction,
						GenerationID: ptr.To(int64(2448)),
						LineageID:    500022,
						Network:      ptr.To("PRODUCTION"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   4540,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-22T12:47:56Z"),
						EventType:    ActivityEventTypeGenerationPromotedToStaging,
						GenerationID: ptr.To(int64(2448)),
						LineageID:    500022,
						Network:      ptr.To("STAGING"),
						Outcome:      ptr.To("ALL_SUCCESS"),
					},
					{
						ActivityID:   4296,
						CreatedBy:    "terraform-dev",
						CreatedTime:  test.NewTimeFromString(t, "2026-07-22T12:31:00Z"),
						EventType:    ActivityEventTypeCertUploaded,
						GenerationID: ptr.To(int64(2448)),
						LineageID:    500022,
					},
				},
				NextCursor: ptr.To("NDI5Ng=="),
				TotalCount: 22,
			},
		},
		"404 lineage not found": {
			params:         ListLineageActivityRequest{LineageID: 999999},
			expectedPath:   "/ccm/v2/lineages/999999/activity",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891026",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineageActivity, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891026",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrListLineageActivity)
			},
		},
		"400 invalid or expired cursor": {
			params:         ListLineageActivityRequest{LineageID: 500005, Cursor: "invalid-or-expired-cursor"},
			expectedPath:   "/ccm/v2/lineages/500005/activity?after=invalid-or-expired-cursor",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/invalid-cursor",
				"title": "Invalid or expired pagination cursor.",
				"status": 400,
				"detail": "The 'after' cursor value is invalid or has expired. Please restart pagination without a cursor.",
				"instance": "/error-types/invalid-cursor?traceId=1234567891027"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineageActivity, &Error{
					Type:     "/error-types/invalid-cursor",
					Title:    "Invalid or expired pagination cursor.",
					Status:   http.StatusBadRequest,
					Detail:   "The 'after' cursor value is invalid or has expired. Please restart pagination without a cursor.",
					Instance: "/error-types/invalid-cursor?traceId=1234567891027",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInvalidCursor)
				assert.ErrorIs(t, err, ErrListLineageActivity)
			},
		},
		"500 internal server error": {
			params:         ListLineageActivityRequest{LineageID: 500005},
			expectedPath:   "/ccm/v2/lineages/500005/activity",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891022"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListLineageActivity, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891022",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrListLineageActivity)
			},
		},
		"validation error - missing LineageID": {
			params: ListLineageActivityRequest{},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing lineage activity: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrListLineageActivity)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - PageSize exceeds maximum": {
			params: ListLineageActivityRequest{LineageID: 500005, PageSize: 101},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, fmt.Sprintf("listing lineage activity: struct validation: PageSize: must be no greater than %d", MaxListLineageActivityPageSize))
				assert.ErrorIs(t, err, ErrListLineageActivity)
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
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.ListLineageActivity(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestPromoteLineage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params              PromoteLineageRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedResponse    *PromoteLineageResponse
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"202 OK - promote lineage to production, PENDING": {
			params: PromoteLineageRequest{
				LineageID:    500035,
				GenerationID: 7274,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/500035/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":7274,"networks":["PRODUCTION"]}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-08-13T09:56:01Z",
						"createdBy": "terraform-dev",
						"generationId": 7274,
						"activationId": 6102,
						"lineageId": 500035,
						"modifiedBy": "terraform-dev",
						"activationStatus": "PENDING",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-08-13T09:56:01Z"
					}
				],
				"totalCount": 0
			}`,
			expectedResponse: &PromoteLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "PROMOTE",
						ActivationCreatedTime:  time.Date(2026, 8, 13, 9, 56, 1, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           7274,
						ActivationID:           6102,
						LineageID:              500035,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "PENDING",
						TargetEnvironment:      "PRODUCTION",
						ActivationModifiedTime: time.Date(2026, 8, 13, 9, 56, 1, 0, time.UTC),
					},
				},
			},
		},
		"202 OK - promote lineage to staging, IN_PROGRESS": {
			params: PromoteLineageRequest{
				LineageID:    500036,
				GenerationID: 7273,
				Networks:     []TargetNetwork{TargetNetworkStaging},
			},
			expectedPath:        "/ccm/v2/lineages/500036/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":7273,"networks":["STAGING"]}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-08-13T09:48:41Z",
						"createdBy": "terraform-dev",
						"generationId": 7273,
						"activationId": 6100,
						"lineageId": 500036,
						"modifiedBy": "terraform-dev",
						"activationStatus": "IN_PROGRESS",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-08-13T09:49:15Z"
					}
				],
				"totalCount": 0
			}`,
			expectedResponse: &PromoteLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "PROMOTE",
						ActivationCreatedTime:  time.Date(2026, 8, 13, 9, 48, 41, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           7273,
						ActivationID:           6100,
						LineageID:              500036,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "IN_PROGRESS",
						TargetEnvironment:      "STAGING",
						ActivationModifiedTime: time.Date(2026, 8, 13, 9, 49, 15, 0, time.UTC),
					},
				},
			},
		},
		"202 OK - promote lineage, empty networks omitted from request body (defaults to production)": {
			params: PromoteLineageRequest{
				LineageID:    500027,
				GenerationID: 4075,
			},
			expectedPath:        "/ccm/v2/lineages/500027/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":4075}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-30T09:56:17Z",
						"createdBy": "terraform-dev",
						"generationId": 4075,
						"activationId": 3011,
						"lineageId": 500027,
						"modifiedBy": "terraform-dev",
						"activationStatus": "IN_PROGRESS",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-30T09:56:17Z"
					}
				],
				"totalCount": 0
			}`,
			expectedResponse: &PromoteLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "PROMOTE",
						ActivationCreatedTime:  time.Date(2026, 7, 30, 9, 56, 17, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           4075,
						ActivationID:           3011,
						LineageID:              500027,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "IN_PROGRESS",
						TargetEnvironment:      "PRODUCTION",
						ActivationModifiedTime: time.Date(2026, 7, 30, 9, 56, 17, 0, time.UTC),
					},
				},
			},
		},
		"204 No Content - generation already active on requested network (idempotent no-op)": {
			params: PromoteLineageRequest{
				LineageID:    500025,
				GenerationID: 4090,
				Networks:     []TargetNetwork{TargetNetworkStaging},
			},
			expectedPath:        "/ccm/v2/lineages/500025/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":4090,"networks":["STAGING"]}`,
			responseStatus:      http.StatusNoContent,
			expectedResponse:    nil,
		},
		"202 OK - promote lineage to both production and staging": {
			params: PromoteLineageRequest{
				LineageID:    500019,
				GenerationID: 3032,
				Networks:     []TargetNetwork{TargetNetworkStaging, TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/500019/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":3032,"networks":["STAGING","PRODUCTION"]}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T13:49:13Z",
						"createdBy": "terraform-dev",
						"generationId": 3032,
						"activationId": 1258,
						"lineageId": 500019,
						"modifiedBy": "terraform-dev",
						"activationStatus": "IN_PROGRESS",
						"targetEnvironment": "PRODUCTION/STAGING",
						"activationModifiedTime": "2026-07-23T13:49:13Z"
					}
				],
				"totalCount": 0
			}`,
			expectedResponse: &PromoteLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "PROMOTE",
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           3032,
						ActivationID:           1258,
						LineageID:              500019,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "IN_PROGRESS",
						TargetEnvironment:      "PRODUCTION/STAGING",
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
					},
				},
			},
		},
		"202 OK - first promote to both networks, partial MULTIPLE_STACK activation warnings": {
			params: PromoteLineageRequest{
				LineageID:    185485,
				GenerationID: 57611,
				Networks:     []TargetNetwork{TargetNetworkProduction, TargetNetworkStaging},
			},
			expectedPath:        "/ccm/v2/lineages/185485/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":57611,"networks":["PRODUCTION","STAGING"]}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationCreatedTime": "2026-09-17T12:20:29Z",
						"activationId": 79503,
						"activationModifiedTime": "2026-09-17T12:20:29Z",
						"activationStatus": "PENDING",
						"activationType": "PROMOTE",
						"createdBy": "terraform-dev",
						"generationId": 57611,
						"lineageId": 185485,
						"modifiedBy": "terraform-dev",
						"targetEnvironment": "STAGING"
					},
					{
						"activationCreatedTime": "2026-09-17T12:20:29Z",
						"activationId": 79504,
						"activationModifiedTime": "2026-09-17T12:20:29Z",
						"activationStatus": "PENDING",
						"activationType": "PROMOTE",
						"createdBy": "terraform-dev",
						"generationId": 57611,
						"lineageId": 185485,
						"modifiedBy": "terraform-dev",
						"targetEnvironment": "PRODUCTION"
					}
				],
				"totalCount": 2,
				"validationResults": {
					"warnings": [
						{
							"detail": "Only algorithms with certificateStatus READY_FOR_USE are deployed. The following algorithms were excluded because they have not completed certificate upload: ECDSA=CSR_READY. To include them later without renewing all certificates, call COMPLETE once the certificate is obtained from your CA; if the live certificate is approaching expiry, use RENEW instead.",
							"instance": "/error-types/partial-multiple-stack-activation?traceId=b551ba8494084468",
							"status": 202,
							"title": "Only a subset of key algorithms are being activated.",
							"type": "/error-types/partial-multiple-stack-activation"
						},
						{
							"detail": "Promotion only provisions certificate materials for lineage and makes the lineage be available for binding to hostnames on the configuration side and does NOT perform any hostname binding changes.",
							"instance": "/error-types/first-promote-notice?traceId=b551ba8494084468",
							"status": 202,
							"title": "First promotion only provisions certificate materials; it does not perform hostname binding changes.",
							"type": "/error-types/first-promote-notice"
						},
						{
							"detail": "Only algorithms with certificateStatus READY_FOR_USE are deployed. The following algorithms were excluded because they have not completed certificate upload: ECDSA=CSR_READY. To include them later without renewing all certificates, call COMPLETE once the certificate is obtained from your CA; if the live certificate is approaching expiry, use RENEW instead.",
							"instance": "/error-types/partial-multiple-stack-activation?traceId=b551ba8494084468",
							"status": 202,
							"title": "Only a subset of key algorithms are being activated.",
							"type": "/error-types/partial-multiple-stack-activation"
						},
						{
							"detail": "Promotion only provisions certificate materials for lineage and makes the lineage be available for binding to hostnames on the configuration side and does NOT perform any hostname binding changes.",
							"instance": "/error-types/first-promote-notice?traceId=b551ba8494084468",
							"status": 202,
							"title": "First promotion only provisions certificate materials; it does not perform hostname binding changes.",
							"type": "/error-types/first-promote-notice"
						}
					]
				}
			}`,
			expectedResponse: &PromoteLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "PROMOTE",
						ActivationCreatedTime:  time.Date(2026, 9, 17, 12, 20, 29, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           57611,
						ActivationID:           79503,
						LineageID:              185485,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "PENDING",
						TargetEnvironment:      "STAGING",
						ActivationModifiedTime: time.Date(2026, 9, 17, 12, 20, 29, 0, time.UTC),
					},
					{
						ActivationType:         "PROMOTE",
						ActivationCreatedTime:  time.Date(2026, 9, 17, 12, 20, 29, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           57611,
						ActivationID:           79504,
						LineageID:              185485,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "PENDING",
						TargetEnvironment:      "PRODUCTION",
						ActivationModifiedTime: time.Date(2026, 9, 17, 12, 20, 29, 0, time.UTC),
					},
				},
				ValidationResults: &ValidationResults{
					Warnings: []ValidationResultItem{
						{
							Detail:   "Only algorithms with certificateStatus READY_FOR_USE are deployed. The following algorithms were excluded because they have not completed certificate upload: ECDSA=CSR_READY. To include them later without renewing all certificates, call COMPLETE once the certificate is obtained from your CA; if the live certificate is approaching expiry, use RENEW instead.",
							Instance: "/error-types/partial-multiple-stack-activation?traceId=b551ba8494084468",
							Status:   http.StatusAccepted,
							Title:    "Only a subset of key algorithms are being activated.",
							Type:     "/error-types/partial-multiple-stack-activation",
						},
						{
							Detail:   "Promotion only provisions certificate materials for lineage and makes the lineage be available for binding to hostnames on the configuration side and does NOT perform any hostname binding changes.",
							Instance: "/error-types/first-promote-notice?traceId=b551ba8494084468",
							Status:   http.StatusAccepted,
							Title:    "First promotion only provisions certificate materials; it does not perform hostname binding changes.",
							Type:     "/error-types/first-promote-notice",
						},
						{
							Detail:   "Only algorithms with certificateStatus READY_FOR_USE are deployed. The following algorithms were excluded because they have not completed certificate upload: ECDSA=CSR_READY. To include them later without renewing all certificates, call COMPLETE once the certificate is obtained from your CA; if the live certificate is approaching expiry, use RENEW instead.",
							Instance: "/error-types/partial-multiple-stack-activation?traceId=b551ba8494084468",
							Status:   http.StatusAccepted,
							Title:    "Only a subset of key algorithms are being activated.",
							Type:     "/error-types/partial-multiple-stack-activation",
						},
						{
							Detail:   "Promotion only provisions certificate materials for lineage and makes the lineage be available for binding to hostnames on the configuration side and does NOT perform any hostname binding changes.",
							Instance: "/error-types/first-promote-notice?traceId=b551ba8494084468",
							Status:   http.StatusAccepted,
							Title:    "First promotion only provisions certificate materials; it does not perform hostname binding changes.",
							Type:     "/error-types/first-promote-notice",
						},
					},
				},
			},
		},
		"400 lineage bad request, generationId is not head generation": {
			params: PromoteLineageRequest{
				LineageID:    500038,
				GenerationID: 12345,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/500038/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":12345,"networks":["PRODUCTION"]}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/lineage-bad-request",
				"title": "Lineage bad request.",
				"status": 400,
				"detail": "Lineage operation failed due to a bad request: {PROMOTE requires the generation (generationId=7276) to be head generation. Use ROLLBACK or REPLACE_STAGING for other generation-specific activations.}",
				"instance": "/error-types/lineage-bad-request?traceId=-4315570061984641290",
				"context": {
					"reason": "PROMOTE requires the generation (generationId=7276) to be head generation. Use ROLLBACK or REPLACE_STAGING for other generation-specific activations."
				}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPromoteLineage, &Error{
					Type:     "/error-types/lineage-bad-request",
					Title:    "Lineage bad request.",
					Status:   http.StatusBadRequest,
					Detail:   "Lineage operation failed due to a bad request: {PROMOTE requires the generation (generationId=7276) to be head generation. Use ROLLBACK or REPLACE_STAGING for other generation-specific activations.}",
					Instance: "/error-types/lineage-bad-request?traceId=-4315570061984641290",
					Context:  map[string]any{"reason": "PROMOTE requires the generation (generationId=7276) to be head generation. Use ROLLBACK or REPLACE_STAGING for other generation-specific activations."},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageBadRequest)
				assert.ErrorIs(t, err, ErrPromoteLineage)
			},
		},
		"409 incomplete cert material": {
			params: PromoteLineageRequest{
				LineageID:    500026,
				GenerationID: 4068,
				Networks:     []TargetNetwork{TargetNetworkStaging},
			},
			expectedPath:        "/ccm/v2/lineages/500026/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":4068,"networks":["STAGING"]}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/incomplete-cert-material",
				"title": "Required certificate material is missing.",
				"status": 409,
				"detail": "At least one signed certificate is required before promoting generation {4068}.",
				"instance": "/error-types/incomplete-cert-material?traceId=1234567891068",
				"context": {
					"generationId": 4068,
					"reason": "At least one algorithm instance must be in READY_FOR_USE status to promote. Please upload a signed certificate before promoting."
				}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPromoteLineage, &Error{
					Type:     "/error-types/incomplete-cert-material",
					Title:    "Required certificate material is missing.",
					Status:   http.StatusConflict,
					Detail:   "At least one signed certificate is required before promoting generation {4068}.",
					Instance: "/error-types/incomplete-cert-material?traceId=1234567891068",
					Context: map[string]any{
						"generationId": float64(4068),
						"reason":       "At least one algorithm instance must be in READY_FOR_USE status to promote. Please upload a signed certificate before promoting.",
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrIncompleteCertMaterial)
				assert.ErrorIs(t, err, ErrPromoteLineage)
			},
		},
		"409 lineage no head generation": {
			params: PromoteLineageRequest{
				LineageID:    500020,
				GenerationID: 4200,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/500020/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":4200,"networks":["PRODUCTION"]}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/lineage-no-head-generation",
				"title": "No head generation exists for this lineage.",
				"status": 409,
				"detail": "No head generation exists for lineage {500020}. Cannot perform activation operation.",
				"instance": "/error-types/lineage-no-head-generation?traceId=1234567891063",
				"context": {"lineageId": 500020}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPromoteLineage, &Error{
					Type:     "/error-types/lineage-no-head-generation",
					Title:    "No head generation exists for this lineage.",
					Status:   http.StatusConflict,
					Detail:   "No head generation exists for lineage {500020}. Cannot perform activation operation.",
					Instance: "/error-types/lineage-no-head-generation?traceId=1234567891063",
					Context:  map[string]any{"lineageId": float64(500020)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNoHeadGeneration)
				assert.ErrorIs(t, err, ErrPromoteLineage)
			},
		},
		"400 first promote requires both networks": {
			params: PromoteLineageRequest{
				LineageID:    185485,
				GenerationID: 9000,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/185485/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":9000,"networks":["PRODUCTION"]}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/first-promote-requires-both-networks",
				"title": "The first promotion of a lineage must activate both STAGING and PRODUCTION.",
				"status": 400,
				"detail": "Lineage {185485} has not completed its first promotion. The first PROMOTE must target both STAGING and PRODUCTION networks.",
				"instance": "/error-types/first-promote-requires-both-networks?traceId=c73d40a204a34afa",
				"context": {"lineageId": 185485}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPromoteLineage, &Error{
					Type:     "/error-types/first-promote-requires-both-networks",
					Title:    "The first promotion of a lineage must activate both STAGING and PRODUCTION.",
					Status:   http.StatusBadRequest,
					Detail:   "Lineage {185485} has not completed its first promotion. The first PROMOTE must target both STAGING and PRODUCTION networks.",
					Instance: "/error-types/first-promote-requires-both-networks?traceId=c73d40a204a34afa",
					Context:  map[string]any{"lineageId": float64(185485)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrFirstPromoteRequiresBothNetworks)
				assert.ErrorIs(t, err, ErrPromoteLineage)
			},
		},
		"409 single generation activation not supported": {
			params: PromoteLineageRequest{
				LineageID:    500021,
				GenerationID: 4201,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/500021/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":4201,"networks":["PRODUCTION"]}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/single-generation-activation-not-supported",
				"title": "Activation operations are not supported for SINGLE_GENERATION lineages.",
				"status": 409,
				"detail": "Lineage {500021} is SINGLE_GENERATION type. PROMOTE, ROLLBACK, and REPLACE_STAGING are not applicable.",
				"instance": "/error-types/single-generation-activation-not-supported?traceId=1234567891064",
				"context": {"lineageId": 500021}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPromoteLineage, &Error{
					Type:     "/error-types/single-generation-activation-not-supported",
					Title:    "Activation operations are not supported for SINGLE_GENERATION lineages.",
					Status:   http.StatusConflict,
					Detail:   "Lineage {500021} is SINGLE_GENERATION type. PROMOTE, ROLLBACK, and REPLACE_STAGING are not applicable.",
					Instance: "/error-types/single-generation-activation-not-supported?traceId=1234567891064",
					Context:  map[string]any{"lineageId": float64(500021)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrSingleGenerationActivationNotSupported)
				assert.ErrorIs(t, err, ErrPromoteLineage)
			},
		},
		"404 lineage not found": {
			params: PromoteLineageRequest{
				LineageID:    999999,
				GenerationID: 4202,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/999999/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":4202,"networks":["PRODUCTION"]}`,
			responseStatus:      http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891065",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPromoteLineage, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891065",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrPromoteLineage)
			},
		},
		"409 activation cooldown in effect": {
			params: PromoteLineageRequest{
				LineageID:    500022,
				GenerationID: 3031,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/500022/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":3031,"networks":["PRODUCTION"]}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/activation-cooldown-in-effect",
				"title": "A recent activation is still propagating to the edge network. Wait before starting another.",
				"status": 409,
				"detail": "Lineage {500022} had a recent activation that is still propagating. Retry after {2026-07-23T13:47:30.389694Z}.",
				"instance": "/error-types/activation-cooldown-in-effect?traceId=1234567891066",
				"context": {"cooldownEndsAt": "2026-07-23T13:47:30Z", "lineageId": 500022}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPromoteLineage, &Error{
					Type:     "/error-types/activation-cooldown-in-effect",
					Title:    "A recent activation is still propagating to the edge network. Wait before starting another.",
					Status:   http.StatusConflict,
					Detail:   "Lineage {500022} had a recent activation that is still propagating. Retry after {2026-07-23T13:47:30.389694Z}.",
					Instance: "/error-types/activation-cooldown-in-effect?traceId=1234567891066",
					Context:  map[string]any{"cooldownEndsAt": "2026-07-23T13:47:30Z", "lineageId": float64(500022)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrActivationCooldownInEffect)
				assert.ErrorIs(t, err, ErrPromoteLineage)
			},
		},
		"502 upstream activation error": {
			params: PromoteLineageRequest{
				LineageID:    500029,
				GenerationID: 4203,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			expectedPath:        "/ccm/v2/lineages/500029/activations",
			expectedRequestBody: `{"operationType":"PROMOTE","generationId":4203,"networks":["PRODUCTION"]}`,
			responseStatus:      http.StatusBadGateway,
			responseBody: `{
				"type": "/error-types/upstream-activation-error",
				"title": "An upstream service returned an error while processing the activation.",
				"status": 502,
				"detail": "Activation failed due to an upstream service error.",
				"instance": "/error-types/upstream-activation-error?traceId=1234567891070",
				"context": {"reason": "activation failed due to an upstream service error"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrPromoteLineage, &Error{
					Type:     "/error-types/upstream-activation-error",
					Title:    "An upstream service returned an error while processing the activation.",
					Status:   http.StatusBadGateway,
					Detail:   "Activation failed due to an upstream service error.",
					Instance: "/error-types/upstream-activation-error?traceId=1234567891070",
					Context:  map[string]any{"reason": "activation failed due to an upstream service error"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrUpstreamActivationError)
				assert.ErrorIs(t, err, ErrPromoteLineage)
			},
		},
		"validation error - missing LineageID": {
			params: PromoteLineageRequest{
				GenerationID: 4200,
				Networks:     []TargetNetwork{TargetNetworkProduction},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "promoting lineage: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrPromoteLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing GenerationID": {
			params: PromoteLineageRequest{
				LineageID: 500019,
				Networks:  []TargetNetwork{TargetNetworkProduction},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "promoting lineage: struct validation: GenerationID: cannot be blank")
				assert.ErrorIs(t, err, ErrPromoteLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid network value": {
			params: PromoteLineageRequest{
				LineageID:    500019,
				GenerationID: 3203,
				Networks:     []TargetNetwork{"NOT_A_REAL_NETWORK"},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "promoting lineage: struct validation: 0: value 'NOT_A_REAL_NETWORK' is invalid. Must be either 'STAGING' or 'PRODUCTION'")
				assert.ErrorIs(t, err, ErrPromoteLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if tc.expectedRequestBody != "" {
					requestBody, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(requestBody))
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.PromoteLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestRollbackLineage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params              RollbackLineageRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedResponse    *RollbackLineageResponse
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"202 OK - rollback lineage": {
			params: RollbackLineageRequest{
				LineageID:    500019,
				GenerationID: 3028,
			},
			expectedPath:        "/ccm/v2/lineages/500019/activations",
			expectedRequestBody: `{"operationType":"ROLLBACK","generationId":3028}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationType": "ROLLBACK",
						"activationCreatedTime": "2026-07-23T12:46:26Z",
						"createdBy": "terraform-dev",
						"generationId": 3028,
						"activationId": 1252,
						"lineageId": 500019,
						"modifiedBy": "terraform-dev",
						"activationStatus": "IN_PROGRESS",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T12:46:26Z"
					}
				],
				"totalCount": 0
			}`,
			expectedResponse: &RollbackLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "ROLLBACK",
						ActivationCreatedTime:  time.Date(2026, 7, 23, 12, 46, 26, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           3028,
						ActivationID:           1252,
						LineageID:              500019,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "IN_PROGRESS",
						TargetEnvironment:      "PRODUCTION",
						ActivationModifiedTime: time.Date(2026, 7, 23, 12, 46, 26, 0, time.UTC),
					},
				},
			},
		},
		"202 OK - rollback lineage, PENDING": {
			params: RollbackLineageRequest{
				LineageID:    500035,
				GenerationID: 7274,
			},
			expectedPath:        "/ccm/v2/lineages/500035/activations",
			expectedRequestBody: `{"operationType":"ROLLBACK","generationId":7274}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationType": "ROLLBACK",
						"activationCreatedTime": "2026-08-13T11:05:07Z",
						"createdBy": "terraform-dev",
						"generationId": 7274,
						"activationId": 6106,
						"lineageId": 500035,
						"modifiedBy": "terraform-dev",
						"activationStatus": "PENDING",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-08-13T11:05:07Z"
					}
				],
				"totalCount": 0
			}`,
			expectedResponse: &RollbackLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "ROLLBACK",
						ActivationCreatedTime:  time.Date(2026, 8, 13, 11, 5, 7, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           7274,
						ActivationID:           6106,
						LineageID:              500035,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "PENDING",
						TargetEnvironment:      "PRODUCTION",
						ActivationModifiedTime: time.Date(2026, 8, 13, 11, 5, 7, 0, time.UTC),
					},
				},
			},
		},
		"204 No Content - previous production generation already active (idempotent no-op)": {
			params: RollbackLineageRequest{
				LineageID:    500037,
				GenerationID: 7275,
			},
			expectedPath:        "/ccm/v2/lineages/500037/activations",
			expectedRequestBody: `{"operationType":"ROLLBACK","generationId":7275}`,
			responseStatus:      http.StatusNoContent,
			expectedResponse:    nil,
		},
		"400 lineage bad request, generationId is not previous production generation": {
			params: RollbackLineageRequest{
				LineageID:    500039,
				GenerationID: 12345,
			},
			expectedPath:        "/ccm/v2/lineages/500039/activations",
			expectedRequestBody: `{"operationType":"ROLLBACK","generationId":12345}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/lineage-bad-request",
				"title": "Lineage bad request.",
				"status": 400,
				"detail": "Lineage operation failed due to a bad request: {ROLLBACK is disallowed for generation 12345: it is not the previousProduction generation.}",
				"instance": "/error-types/lineage-bad-request?traceId=8871564688816097839",
				"context": {
					"reason": "ROLLBACK is disallowed for generation 12345: it is not the previousProduction generation."
				}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRollbackLineage, &Error{
					Type:     "/error-types/lineage-bad-request",
					Title:    "Lineage bad request.",
					Status:   http.StatusBadRequest,
					Detail:   "Lineage operation failed due to a bad request: {ROLLBACK is disallowed for generation 12345: it is not the previousProduction generation.}",
					Instance: "/error-types/lineage-bad-request?traceId=8871564688816097839",
					Context:  map[string]any{"reason": "ROLLBACK is disallowed for generation 12345: it is not the previousProduction generation."},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageBadRequest)
				assert.ErrorIs(t, err, ErrRollbackLineage)
			},
		},
		"400 rollback unavailable, no previous production generation": {
			params: RollbackLineageRequest{
				LineageID:    500028,
				GenerationID: 4079,
			},
			expectedPath:        "/ccm/v2/lineages/500028/activations",
			expectedRequestBody: `{"operationType":"ROLLBACK","generationId":4079}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/rollback-unavailable",
				"title": "Rollback is unavailable.",
				"status": 400,
				"detail": "No previous production generation exists for lineage {500028} to roll back to.",
				"instance": "/error-types/rollback-unavailable?traceId=1234567891069",
				"context": {"lineageId": 500028}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRollbackLineage, &Error{
					Type:     "/error-types/rollback-unavailable",
					Title:    "Rollback is unavailable.",
					Status:   http.StatusBadRequest,
					Detail:   "No previous production generation exists for lineage {500028} to roll back to.",
					Instance: "/error-types/rollback-unavailable?traceId=1234567891069",
					Context:  map[string]any{"lineageId": float64(500028)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrRollbackUnavailable)
				assert.ErrorIs(t, err, ErrRollbackLineage)
			},
		},
		"502 upstream activation error": {
			params: RollbackLineageRequest{
				LineageID:    500030,
				GenerationID: 4204,
			},
			expectedPath:        "/ccm/v2/lineages/500030/activations",
			expectedRequestBody: `{"operationType":"ROLLBACK","generationId":4204}`,
			responseStatus:      http.StatusBadGateway,
			responseBody: `{
				"type": "/error-types/upstream-activation-error",
				"title": "An upstream service returned an error while processing the activation.",
				"status": 502,
				"detail": "Activation failed due to an upstream service error.",
				"instance": "/error-types/upstream-activation-error?traceId=1234567891071",
				"context": {"reason": "activation failed due to an upstream service error"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRollbackLineage, &Error{
					Type:     "/error-types/upstream-activation-error",
					Title:    "An upstream service returned an error while processing the activation.",
					Status:   http.StatusBadGateway,
					Detail:   "Activation failed due to an upstream service error.",
					Instance: "/error-types/upstream-activation-error?traceId=1234567891071",
					Context:  map[string]any{"reason": "activation failed due to an upstream service error"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrUpstreamActivationError)
				assert.ErrorIs(t, err, ErrRollbackLineage)
			},
		},
		"404 lineage not found": {
			params: RollbackLineageRequest{
				LineageID:    999999,
				GenerationID: 4205,
			},
			expectedPath:        "/ccm/v2/lineages/999999/activations",
			expectedRequestBody: `{"operationType":"ROLLBACK","generationId":4205}`,
			responseStatus:      http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891072",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRollbackLineage, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891072",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrRollbackLineage)
			},
		},
		"validation error - missing LineageID": {
			params: RollbackLineageRequest{
				GenerationID: 4200,
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "rolling back lineage: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrRollbackLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing GenerationID": {
			params: RollbackLineageRequest{
				LineageID: 500019,
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "rolling back lineage: struct validation: GenerationID: cannot be blank")
				assert.ErrorIs(t, err, ErrRollbackLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if tc.expectedRequestBody != "" {
					requestBody, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(requestBody))
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.RollbackLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestReplaceStagingLineage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params              ReplaceStagingLineageRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedResponse    *ReplaceStagingLineageResponse
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"202 OK - replace staging lineage": {
			params: ReplaceStagingLineageRequest{
				LineageID:    500019,
				GenerationID: 3028,
			},
			expectedPath:        "/ccm/v2/lineages/500019/activations",
			expectedRequestBody: `{"operationType":"REPLACE_STAGING","generationId":3028}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationType": "REPLACE_STAGING",
						"activationCreatedTime": "2026-07-23T13:03:24Z",
						"createdBy": "terraform-dev",
						"generationId": 3028,
						"activationId": 1254,
						"lineageId": 500019,
						"modifiedBy": "terraform-dev",
						"activationStatus": "IN_PROGRESS",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-07-23T13:03:24Z"
					}
				],
				"totalCount": 0
			}`,
			expectedResponse: &ReplaceStagingLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "REPLACE_STAGING",
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 3, 24, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           3028,
						ActivationID:           1254,
						LineageID:              500019,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "IN_PROGRESS",
						TargetEnvironment:      "STAGING",
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 3, 24, 0, time.UTC),
					},
				},
			},
		},
		"202 OK - replace staging lineage, PENDING": {
			params: ReplaceStagingLineageRequest{
				LineageID:    500035,
				GenerationID: 7274,
			},
			expectedPath:        "/ccm/v2/lineages/500035/activations",
			expectedRequestBody: `{"operationType":"REPLACE_STAGING","generationId":7274}`,
			responseStatus:      http.StatusAccepted,
			responseBody: `{
				"items": [
					{
						"activationType": "REPLACE_STAGING",
						"activationCreatedTime": "2026-08-13T11:14:16Z",
						"createdBy": "terraform-dev",
						"generationId": 7274,
						"activationId": 6108,
						"lineageId": 500035,
						"modifiedBy": "terraform-dev",
						"activationStatus": "PENDING",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-08-13T11:14:16Z"
					}
				],
				"totalCount": 0
			}`,
			expectedResponse: &ReplaceStagingLineageResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationType:         "REPLACE_STAGING",
						ActivationCreatedTime:  time.Date(2026, 8, 13, 11, 14, 16, 0, time.UTC),
						CreatedBy:              "terraform-dev",
						GenerationID:           7274,
						ActivationID:           6108,
						LineageID:              500035,
						ModifiedBy:             "terraform-dev",
						ActivationStatus:       "PENDING",
						TargetEnvironment:      "STAGING",
						ActivationModifiedTime: time.Date(2026, 8, 13, 11, 14, 16, 0, time.UTC),
					},
				},
			},
		},
		"204 No Content - GenerationID already current staging generation (idempotent no-op)": {
			params: ReplaceStagingLineageRequest{
				LineageID:    500035,
				GenerationID: 7274,
			},
			expectedPath:        "/ccm/v2/lineages/500035/activations",
			expectedRequestBody: `{"operationType":"REPLACE_STAGING","generationId":7274}`,
			responseStatus:      http.StatusNoContent,
			expectedResponse:    nil,
		},
		"409 replace staging head not on staging": {
			params: ReplaceStagingLineageRequest{
				LineageID:    500022,
				GenerationID: 3031,
			},
			expectedPath:        "/ccm/v2/lineages/500022/activations",
			expectedRequestBody: `{"operationType":"REPLACE_STAGING","generationId":3031}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/replace-staging-head-not-on-staging",
				"title": "Head generation is not currently deployed on the staging network.",
				"status": 409,
				"detail": "Cannot perform REPLACE_STAGING on lineage {500022}: head generation is not the current staging generation. REPLACE_STAGING is only valid when the head is deployed on staging.",
				"instance": "/error-types/replace-staging-head-not-on-staging?traceId=1234567891067",
				"context": {"lineageId": 500022}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrReplaceStagingLineage, &Error{
					Type:     "/error-types/replace-staging-head-not-on-staging",
					Title:    "Head generation is not currently deployed on the staging network.",
					Status:   http.StatusConflict,
					Detail:   "Cannot perform REPLACE_STAGING on lineage {500022}: head generation is not the current staging generation. REPLACE_STAGING is only valid when the head is deployed on staging.",
					Instance: "/error-types/replace-staging-head-not-on-staging?traceId=1234567891067",
					Context:  map[string]any{"lineageId": float64(500022)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrReplaceStagingHeadNotOnStaging)
				assert.ErrorIs(t, err, ErrReplaceStagingLineage)
			},
		},
		"409 no current staging generation": {
			params: ReplaceStagingLineageRequest{
				LineageID:    500027,
				GenerationID: 4075,
			},
			expectedPath:        "/ccm/v2/lineages/500027/activations",
			expectedRequestBody: `{"operationType":"REPLACE_STAGING","generationId":4075}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/no-current-staging",
				"title": "No current staging generation exists.",
				"status": 409,
				"detail": "No current staging generation exists for lineage {500027}. Cannot replace staging without an existing staging generation.",
				"instance": "/error-types/no-current-staging?traceId=1234567891075",
				"context": {"lineageId": 500027}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrReplaceStagingLineage, &Error{
					Type:     "/error-types/no-current-staging",
					Title:    "No current staging generation exists.",
					Status:   http.StatusConflict,
					Detail:   "No current staging generation exists for lineage {500027}. Cannot replace staging without an existing staging generation.",
					Instance: "/error-types/no-current-staging?traceId=1234567891075",
					Context:  map[string]any{"lineageId": float64(500027)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrNoCurrentStaging)
				assert.ErrorIs(t, err, ErrReplaceStagingLineage)
			},
		},
		"400 lineage bad request, invalid explicit GenerationID": {
			params: ReplaceStagingLineageRequest{
				LineageID:    500032,
				GenerationID: 12345,
			},
			expectedPath:        "/ccm/v2/lineages/500032/activations",
			expectedRequestBody: `{"operationType":"REPLACE_STAGING","generationId":12345}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/lineage-bad-request",
				"title": "Lineage bad request.",
				"status": 400,
				"detail": "Lineage operation failed due to a bad request: {REPLACE_STAGING generationId must be currentProduction or previousProduction.}",
				"instance": "/error-types/lineage-bad-request?traceId=-7216309907389071349",
				"context": {"reason": "REPLACE_STAGING generationId must be currentProduction or previousProduction."}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrReplaceStagingLineage, &Error{
					Type:     "/error-types/lineage-bad-request",
					Title:    "Lineage bad request.",
					Status:   http.StatusBadRequest,
					Detail:   "Lineage operation failed due to a bad request: {REPLACE_STAGING generationId must be currentProduction or previousProduction.}",
					Instance: "/error-types/lineage-bad-request?traceId=-7216309907389071349",
					Context:  map[string]any{"reason": "REPLACE_STAGING generationId must be currentProduction or previousProduction."},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageBadRequest)
				assert.ErrorIs(t, err, ErrReplaceStagingLineage)
			},
		},
		"404 lineage not found": {
			params: ReplaceStagingLineageRequest{
				LineageID:    999999,
				GenerationID: 4206,
			},
			expectedPath:        "/ccm/v2/lineages/999999/activations",
			expectedRequestBody: `{"operationType":"REPLACE_STAGING","generationId":4206}`,
			responseStatus:      http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891073",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrReplaceStagingLineage, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891073",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrReplaceStagingLineage)
			},
		},
		"validation error - missing LineageID": {
			params: ReplaceStagingLineageRequest{
				GenerationID: 4200,
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "replacing lineage staging generation: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrReplaceStagingLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing GenerationID": {
			params: ReplaceStagingLineageRequest{
				LineageID: 500019,
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "replacing lineage staging generation: struct validation: GenerationID: cannot be blank")
				assert.ErrorIs(t, err, ErrReplaceStagingLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if tc.expectedRequestBody != "" {
					requestBody, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(requestBody))
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.ReplaceStagingLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetActivationStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           GetActivationStatusRequest
		responseStatus   int
		responseBody     string
		expectedResponse *GetActivationStatusResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"200 OK - get activation status": {
			params:         GetActivationStatusRequest{LineageID: 500019, ActivationID: 1075},
			expectedPath:   "/ccm/v2/lineages/500019/activations/1075",
			responseStatus: http.StatusOK,
			responseBody: `{
				"activationType": "PROMOTE",
				"activationCreatedTime": "2026-07-22T12:47:55Z",
				"createdBy": "terraform-dev",
				"generationId": 2448,
				"activationId": 1075,
				"lineageId": 500019,
				"modifiedBy": "terraform-dev",
				"activationStatus": "INIT",
				"targetEnvironment": "STAGING",
				"activationModifiedTime": "2026-07-22T12:47:56Z"
			}`,
			expectedResponse: &GetActivationStatusResponse{
				ActivationID:           1075,
				LineageID:              500019,
				ActivationType:         "PROMOTE",
				ActivationStatus:       "INIT",
				GenerationID:           2448,
				ActivationCreatedTime:  time.Date(2026, 7, 22, 12, 47, 55, 0, time.UTC),
				ActivationModifiedTime: time.Date(2026, 7, 22, 12, 47, 56, 0, time.UTC),
				TargetEnvironment:      "STAGING",
				CreatedBy:              "terraform-dev",
				ModifiedBy:             "terraform-dev",
			},
		},
		"validation error - missing LineageID": {
			params: GetActivationStatusRequest{ActivationID: 1075},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "getting activation status: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrGetActivationStatus)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing ActivationID": {
			params: GetActivationStatusRequest{LineageID: 500019},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "getting activation status: struct validation: ActivationID: cannot be blank")
				assert.ErrorIs(t, err, ErrGetActivationStatus)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"404 activation not found": {
			params:         GetActivationStatusRequest{LineageID: 500022, ActivationID: 1255},
			expectedPath:   "/ccm/v2/lineages/500022/activations/1255",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/activation-not-found",
				"title": "Activation not found.",
				"status": 404,
				"detail": "Activation {1255} not found.",
				"instance": "/error-types/activation-not-found?traceId=1234567891077",
				"context": {"activationId": "1255"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetActivationStatus, &Error{
					Type:     "/error-types/activation-not-found",
					Title:    "Activation not found.",
					Status:   http.StatusNotFound,
					Detail:   "Activation {1255} not found.",
					Instance: "/error-types/activation-not-found?traceId=1234567891077",
					Context:  map[string]any{"activationId": "1255"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrActivationNotFound)
				assert.ErrorIs(t, err, ErrGetActivationStatus)
			},
		},
		"404 lineage not found": {
			params:         GetActivationStatusRequest{LineageID: 999999, ActivationID: 1075},
			expectedPath:   "/ccm/v2/lineages/999999/activations/1075",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891080",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetActivationStatus, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891080",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrGetActivationStatus)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.GetActivationStatus(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListActivations(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           ListActivationsRequest
		responseStatus   int
		responseBody     string
		expectedResponse *ListActivationsResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"200 OK - list activations": {
			params:         ListActivationsRequest{LineageID: 500022},
			expectedPath:   "/ccm/v2/lineages/500022/activations",
			responseStatus: http.StatusOK,
			responseBody: `{
				"items": [
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T13:49:12Z",
						"createdBy": "terraform-dev",
						"generationId": 3032,
						"activationId": 1258,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T13:49:13Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T13:49:12Z",
						"createdBy": "terraform-dev",
						"generationId": 3032,
						"activationId": 1257,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-07-23T13:49:13Z"
					},
					{
						"activationType": "ROLLBACK",
						"activationCreatedTime": "2026-07-23T13:42:29Z",
						"createdBy": "terraform-dev",
						"generationId": 3029,
						"activationId": 1256,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T13:42:30Z"
					},
					{
						"activationType": "REPLACE_STAGING",
						"activationCreatedTime": "2026-07-23T13:03:24Z",
						"createdBy": "terraform-dev",
						"generationId": 3028,
						"activationId": 1254,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-07-23T13:03:24Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T12:52:25Z",
						"createdBy": "terraform-dev",
						"generationId": 3031,
						"activationId": 1253,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-07-23T12:52:26Z"
					},
					{
						"activationType": "ROLLBACK",
						"activationCreatedTime": "2026-07-23T12:46:25Z",
						"createdBy": "terraform-dev",
						"generationId": 3028,
						"activationId": 1252,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T12:46:25Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T12:21:36Z",
						"createdBy": "terraform-dev",
						"generationId": 3029,
						"activationId": 1251,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T12:21:36Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T11:25:05Z",
						"createdBy": "terraform-dev",
						"generationId": 3028,
						"activationId": 1248,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T11:25:05Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T11:25:05Z",
						"createdBy": "terraform-dev",
						"generationId": 3028,
						"activationId": 1247,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-07-23T11:25:05Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T11:18:53Z",
						"createdBy": "terraform-dev",
						"generationId": 3027,
						"activationId": 1244,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T11:18:53Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T11:08:23Z",
						"createdBy": "terraform-dev",
						"generationId": 2448,
						"activationId": 1243,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T11:08:24Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-22T12:47:55Z",
						"createdBy": "terraform-dev",
						"generationId": 2448,
						"activationId": 1075,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-07-22T12:47:56Z"
					}
				],
				"totalCount": 12
			}`,
			expectedResponse: &ListActivationsResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationID: 1258, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3032,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 49, 12, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1257, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3032,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 49, 12, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
						TargetEnvironment:      "STAGING", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1256, LineageID: 500022, ActivationType: "ROLLBACK", ActivationStatus: "INIT", GenerationID: 3029,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 42, 29, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 42, 30, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1254, LineageID: 500022, ActivationType: "REPLACE_STAGING", ActivationStatus: "INIT", GenerationID: 3028,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 3, 24, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 3, 24, 0, time.UTC),
						TargetEnvironment:      "STAGING", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1253, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3031,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 12, 52, 25, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 12, 52, 26, 0, time.UTC),
						TargetEnvironment:      "STAGING", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1252, LineageID: 500022, ActivationType: "ROLLBACK", ActivationStatus: "INIT", GenerationID: 3028,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 12, 46, 25, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 12, 46, 25, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1251, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3029,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 12, 21, 36, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 12, 21, 36, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1248, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3028,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 11, 25, 5, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 11, 25, 5, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1247, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3028,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 11, 25, 5, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 11, 25, 5, 0, time.UTC),
						TargetEnvironment:      "STAGING", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1244, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3027,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 11, 18, 53, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 11, 18, 53, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1243, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 2448,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 11, 8, 23, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 11, 8, 24, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1075, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 2448,
						ActivationCreatedTime:  time.Date(2026, 7, 22, 12, 47, 55, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 22, 12, 47, 56, 0, time.UTC),
						TargetEnvironment:      "STAGING", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
				},
				TotalCount: 12,
			},
		},
		"200 OK - list activations, empty": {
			params:         ListActivationsRequest{LineageID: 500033},
			expectedPath:   "/ccm/v2/lineages/500033/activations",
			responseStatus: http.StatusOK,
			responseBody: `{
				"items": [],
				"totalCount": 0
			}`,
			expectedResponse: &ListActivationsResponse{
				Items:      []GetActivationStatusResponse{},
				TotalCount: 0,
			},
		},
		"200 OK - list activations, full page with populated nextCursor": {
			params:         ListActivationsRequest{LineageID: 500022, PageSize: 3},
			expectedPath:   "/ccm/v2/lineages/500022/activations?pageSize=3",
			responseStatus: http.StatusOK,
			responseBody: `{
				"items": [
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T13:49:12Z",
						"createdBy": "terraform-dev",
						"generationId": 3032,
						"activationId": 1258,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T13:49:13Z"
					},
					{
						"activationType": "PROMOTE",
						"activationCreatedTime": "2026-07-23T13:49:12Z",
						"createdBy": "terraform-dev",
						"generationId": 3032,
						"activationId": 1257,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "STAGING",
						"activationModifiedTime": "2026-07-23T13:49:13Z"
					},
					{
						"activationType": "ROLLBACK",
						"activationCreatedTime": "2026-07-23T13:42:29Z",
						"createdBy": "terraform-dev",
						"generationId": 3029,
						"activationId": 1256,
						"lineageId": 500022,
						"modifiedBy": "terraform-dev",
						"activationStatus": "INIT",
						"targetEnvironment": "PRODUCTION",
						"activationModifiedTime": "2026-07-23T13:42:30Z"
					}
				],
				"nextCursor": "MTI1Ng",
				"totalCount": 12
			}`,
			expectedResponse: &ListActivationsResponse{
				Items: []GetActivationStatusResponse{
					{
						ActivationID: 1258, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3032,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 49, 12, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1257, LineageID: 500022, ActivationType: "PROMOTE", ActivationStatus: "INIT", GenerationID: 3032,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 49, 12, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
						TargetEnvironment:      "STAGING", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
					{
						ActivationID: 1256, LineageID: 500022, ActivationType: "ROLLBACK", ActivationStatus: "INIT", GenerationID: 3029,
						ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 42, 29, 0, time.UTC),
						ActivationModifiedTime: time.Date(2026, 7, 23, 13, 42, 30, 0, time.UTC),
						TargetEnvironment:      "PRODUCTION", CreatedBy: "terraform-dev", ModifiedBy: "terraform-dev",
					},
				},
				NextCursor: ptr.To("MTI1Ng"),
				TotalCount: 12,
			},
		},
		"404 lineage not found": {
			params:         ListActivationsRequest{LineageID: 999999},
			expectedPath:   "/ccm/v2/lineages/999999/activations",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891078",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListActivations, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891078",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrListActivations)
			},
		},
		"400 invalid cursor value": {
			params:         ListActivationsRequest{LineageID: 500034, Cursor: "12423542352342343"},
			expectedPath:   "/ccm/v2/lineages/500034/activations?cursor=12423542352342343",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/invalid-field",
				"title": "Invalid field value.",
				"status": 400,
				"detail": "Invalid value '{12423542352342343}' for field '{cursor}'. Invalid cursor value.",
				"instance": "/error-types/invalid-field?traceId=1234567891079",
				"context": {
					"explanation": "Invalid cursor value.",
					"invalidParameterValue": "12423542352342343",
					"parameterName": "cursor"
				}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListActivations, &Error{
					Type:     "/error-types/invalid-field",
					Title:    "Invalid field value.",
					Status:   http.StatusBadRequest,
					Detail:   "Invalid value '{12423542352342343}' for field '{cursor}'. Invalid cursor value.",
					Instance: "/error-types/invalid-field?traceId=1234567891079",
					Context: map[string]any{
						"explanation":           "Invalid cursor value.",
						"invalidParameterValue": "12423542352342343",
						"parameterName":         "cursor",
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInvalidField)
				assert.ErrorIs(t, err, ErrListActivations)
			},
		},
		"validation error - missing LineageID": {
			params: ListActivationsRequest{},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing activations: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrListActivations)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - PageSize exceeds maximum": {
			params: ListActivationsRequest{LineageID: 500022, PageSize: 101},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, fmt.Sprintf("listing activations: struct validation: PageSize: must be no greater than %d", MaxListActivationsPageSize))
				assert.ErrorIs(t, err, ErrListActivations)
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
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.ListActivations(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestSortOrder_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		sortOrder SortOrder
		withError bool
	}{
		"empty value is valid": {
			sortOrder: "",
			withError: false,
		},
		"ASC is valid": {
			sortOrder: SortOrderAscending,
			withError: false,
		},
		"DESC is valid": {
			sortOrder: SortOrderDescending,
			withError: false,
		},
		"invalid value": {
			sortOrder: "NOT_A_REAL_SORT",
			withError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.sortOrder.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, "value 'NOT_A_REAL_SORT' is invalid. Must be either 'ASC' or 'DESC'")
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestStackMode_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stackMode StackMode
		withError bool
	}{
		"empty value is valid": {
			stackMode: "",
			withError: false,
		},
		"SINGLE_STACK is valid": {
			stackMode: StackModeSingleStack,
			withError: false,
		},
		"MULTIPLE_STACK is valid": {
			stackMode: StackModeMultipleStack,
			withError: false,
		},
		"invalid value": {
			stackMode: "NOT_A_REAL_MODE",
			withError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.stackMode.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, "value 'NOT_A_REAL_MODE' is invalid. Must be either "+
					"'SINGLE_STACK' or 'MULTIPLE_STACK'")
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGenerationStatus_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		generationStatus GenerationStatus
		withError        bool
	}{
		"empty value is valid": {
			generationStatus: "",
			withError:        false,
		},
		"CSR_READY is valid": {
			generationStatus: GenerationStatusCSRReady,
			withError:        false,
		},
		"READY_FOR_USE is valid": {
			generationStatus: GenerationStatusReadyForUse,
			withError:        false,
		},
		"ACTIVE is valid": {
			generationStatus: GenerationStatusActive,
			withError:        false,
		},
		"ARCHIVED is valid": {
			generationStatus: GenerationStatusArchived,
			withError:        false,
		},
		"ABANDONED is valid": {
			generationStatus: GenerationStatusAbandoned,
			withError:        false,
		},
		"invalid value": {
			generationStatus: "NOT_A_REAL_STATUS",
			withError:        true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.generationStatus.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, "value 'NOT_A_REAL_STATUS' is invalid. Must be one of: "+
					"'CSR_READY', 'READY_FOR_USE', 'ACTIVE', 'ARCHIVED', or 'ABANDONED'")
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestSubject_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		subject   Subject
		withError bool
		errorMsg  string
	}{
		"valid - all fields": {
			subject: Subject{
				CommonName:         "example.com",
				Organization:       "Example Corp",
				OrganizationalUnit: "Engineering",
				Country:            "US",
				State:              "Massachusetts",
				Locality:           "Cambridge",
			},
		},
		"valid - no fields specified": {
			subject: Subject{},
		},
		"valid - any non-empty subset of fields (e.g. common_name only)": {
			subject: Subject{CommonName: "example.com"},
		},
		"validation error - common_name too long": {
			subject:   Subject{CommonName: string(make([]byte, 65))},
			withError: true,
			errorMsg:  "CommonName: the length must be between 1 and 64.",
		},
		"validation error - common_name only whitespace": {
			subject:   Subject{CommonName: "   "},
			withError: true,
			errorMsg:  "CommonName: must be in a valid format.",
		},
		"validation error - organization too long": {
			subject:   Subject{Organization: string(make([]byte, 65))},
			withError: true,
			errorMsg:  "Organization: the length must be between 1 and 64.",
		},
		"validation error - organization only whitespace": {
			subject:   Subject{Organization: "   "},
			withError: true,
			errorMsg:  "Organization: must be in a valid format.",
		},
		"validation error - country too short": {
			subject:   Subject{Country: "U"},
			withError: true,
			errorMsg:  "Country: the length must be exactly 2.",
		},
		"validation error - country too long": {
			subject:   Subject{Country: "USA"},
			withError: true,
			errorMsg:  "Country: the length must be exactly 2.",
		},
		"validation error - country only whitespace": {
			subject:   Subject{Country: "  "},
			withError: true,
			errorMsg:  "Country: must be in a valid format.",
		},
		"validation error - state too long": {
			subject:   Subject{State: string(make([]byte, 129))},
			withError: true,
			errorMsg:  "State: the length must be between 1 and 128.",
		},
		"validation error - state only whitespace": {
			subject:   Subject{State: "   "},
			withError: true,
			errorMsg:  "State: must be in a valid format.",
		},
		"validation error - locality too long": {
			subject:   Subject{Locality: string(make([]byte, 129))},
			withError: true,
			errorMsg:  "Locality: the length must be between 1 and 128.",
		},
		"validation error - locality only whitespace": {
			subject:   Subject{Locality: "   "},
			withError: true,
			errorMsg:  "Locality: must be in a valid format.",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.subject.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, tc.errorMsg)
				return
			}
			assert.NoError(t, err)
		})
	}
}
