package cloudcertificates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadSignedCertificate(t *testing.T) {
	t.Parallel()

	baseRequest := UploadSignedCertificateRequest{
		LineageID:           500001,
		GenerationID:        2912,
		AcknowledgeWarnings: true,
		Body: UploadSignedCertificateRequestBody{
			Algorithms: map[CryptographicAlgorithm]SignedCertificate{
				CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
			},
		},
	}

	tests := map[string]struct {
		params              UploadSignedCertificateRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedResponse    *UploadSignedCertificateResponse
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"200 OK - upload RSA signed certificate": {
			params:              baseRequest,
			expectedPath:        "/ccm/v2/lineages/500001/generations/2912?acknowledgeWarnings=true",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"certificateStatus": "CSR_READY",
						"csrExpirationDate": "2027-07-07T10:18:46Z",
						"keyType": "ECDSA"
					},
					{
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-07-07T10:18:46Z",
						"keyType": "RSA",
						"signedCertificateIssuer": "CN=Test Certificate Authority",
						"signedCertificateNotValidAfterDate": "2027-07-08T13:06:53Z",
						"signedCertificateNotValidBeforeDate": "2026-07-08T13:06:53Z",
						"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01",
						"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:01"
					}
				],
				"generationId": 2912,
				"generationModifiedBy": "terraform-dev",
				"generationModifiedTime": "2026-07-08T13:14:25Z",
				"generationStatus": "READY_FOR_USE",
				"lineageId": 500001
			}`,
			expectedResponse: &UploadSignedCertificateResponse{
				Algorithms: []Algorithm{
					{
						CertificateStatus: CertificateStatusCSRReady,
						CSRExpirationDate: ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
						KeyType:           string(CryptographicAlgorithmECDSA),
					},
					{
						CertificateStatus:                   CertificateStatusReadyForUse,
						CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
						KeyType:                             string(CryptographicAlgorithmRSA),
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-08T13:06:53Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-08T13:06:53Z")),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:01"),
					},
				},
				GenerationID:           2912,
				GenerationModifiedBy:   "terraform-dev",
				GenerationModifiedTime: test.NewTimeFromString(t, "2026-07-08T13:14:25Z"),
				GenerationStatus:       string(GenerationStatusReadyForUse),
				LineageID:              500001,
			},
		},
		"200 OK - upload RSA signed certificate, empty (non-nil) AutoActivate is omitted from request body": {
			params: UploadSignedCertificateRequest{
				LineageID:           500001,
				GenerationID:        2912,
				AcknowledgeWarnings: true,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
					AutoActivate: []TargetNetwork{},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500001/generations/2912?acknowledgeWarnings=true",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"certificateStatus": "CSR_READY",
						"csrExpirationDate": "2027-07-07T10:18:46Z",
						"keyType": "ECDSA"
					},
					{
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-07-07T10:18:46Z",
						"keyType": "RSA",
						"signedCertificateIssuer": "CN=Test Certificate Authority",
						"signedCertificateNotValidAfterDate": "2027-07-08T13:06:53Z",
						"signedCertificateNotValidBeforeDate": "2026-07-08T13:06:53Z",
						"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01",
						"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:01"
					}
				],
				"generationId": 2912,
				"generationModifiedBy": "terraform-dev",
				"generationModifiedTime": "2026-07-08T13:14:25Z",
				"generationStatus": "READY_FOR_USE",
				"lineageId": 500001
			}`,
			expectedResponse: &UploadSignedCertificateResponse{
				Algorithms: []Algorithm{
					{
						CertificateStatus: CertificateStatusCSRReady,
						CSRExpirationDate: ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
						KeyType:           string(CryptographicAlgorithmECDSA),
					},
					{
						CertificateStatus:                   CertificateStatusReadyForUse,
						CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
						KeyType:                             string(CryptographicAlgorithmRSA),
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-08T13:06:53Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-08T13:06:53Z")),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:01"),
					},
				},
				GenerationID:           2912,
				GenerationModifiedBy:   "terraform-dev",
				GenerationModifiedTime: test.NewTimeFromString(t, "2026-07-08T13:14:25Z"),
				GenerationStatus:       string(GenerationStatusReadyForUse),
				LineageID:              500001,
			},
		},
		"200 OK - upload RSA signed certificate, single-cipher (SINGLE_STACK) lineage": {
			params: UploadSignedCertificateRequest{
				LineageID:    500016,
				GenerationID: 3247,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500016/generations/3247",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-09-28T12:49:11Z",
						"keyType": "RSA",
						"signedCertificateIssuer": "CN=Test Certificate Authority",
						"signedCertificateNotValidAfterDate": "2027-07-27T12:50:55Z",
						"signedCertificateNotValidBeforeDate": "2026-07-27T12:50:55Z",
						"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:10",
						"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:10"
					}
				],
				"generationId": 3247,
				"generationModifiedBy": "terraform-dev",
				"generationModifiedTime": "2026-07-27T12:51:26Z",
				"generationStatus": "READY_FOR_USE",
				"lineageId": 500016
			}`,
			expectedResponse: &UploadSignedCertificateResponse{
				Algorithms: []Algorithm{
					{
						CertificateStatus:                   CertificateStatusReadyForUse,
						CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-28T12:49:11Z")),
						KeyType:                             string(CryptographicAlgorithmRSA),
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-27T12:50:55Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-27T12:50:55Z")),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:10"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:10"),
					},
				},
				GenerationID:           3247,
				GenerationModifiedBy:   "terraform-dev",
				GenerationModifiedTime: test.NewTimeFromString(t, "2026-07-27T12:51:26Z"),
				GenerationStatus:       string(GenerationStatusReadyForUse),
				LineageID:              500016,
			},
		},
		"200 OK - upload RSA cert with autoActivate STAGING and PRODUCTION": {
			params: UploadSignedCertificateRequest{
				LineageID:           500017,
				GenerationID:        3250,
				AcknowledgeWarnings: true,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
					AutoActivate: []TargetNetwork{TargetNetworkStaging, TargetNetworkProduction},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500017/generations/3250?acknowledgeWarnings=true",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}},"autoActivate":["STAGING","PRODUCTION"]}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"certificateStatus": "CSR_READY",
						"csrExpirationDate": "2027-09-28T13:28:45Z",
						"keyType": "ECDSA"
					},
					{
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-09-28T13:28:47Z",
						"keyType": "RSA",
						"signedCertificateIssuer": "CN=Test Certificate Authority",
						"signedCertificateNotValidAfterDate": "2027-07-27T13:29:28Z",
						"signedCertificateNotValidBeforeDate": "2026-07-27T13:29:28Z",
						"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:11",
						"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:11"
					}
				],
				"firstPromotedToProductionAt": "2026-07-27T13:33:59Z",
				"generationId": 3250,
				"generationModifiedBy": "terraform-dev",
				"generationModifiedTime": "2026-07-27T13:33:41Z",
				"generationStatus": "ACTIVE",
				"lineageId": 500017
			}`,
			expectedResponse: &UploadSignedCertificateResponse{
				Algorithms: []Algorithm{
					{
						CertificateStatus: CertificateStatusCSRReady,
						CSRExpirationDate: ptr.To(test.NewTimeFromString(t, "2027-09-28T13:28:45Z")),
						KeyType:           string(CryptographicAlgorithmECDSA),
					},
					{
						CertificateStatus:                   CertificateStatusReadyForUse,
						CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-28T13:28:47Z")),
						KeyType:                             string(CryptographicAlgorithmRSA),
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-27T13:29:28Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-27T13:29:28Z")),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:11"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:11"),
					},
				},
				FirstPromotedToProductionAt: ptr.To(test.NewTimeFromString(t, "2026-07-27T13:33:59Z")),
				GenerationID:                3250,
				GenerationModifiedBy:        "terraform-dev",
				GenerationModifiedTime:      test.NewTimeFromString(t, "2026-07-27T13:33:41Z"),
				GenerationStatus:            string(GenerationStatusActive),
				LineageID:                   500017,
			},
		},
		"200 OK - upload ECDSA cert, RSA already uploaded": {
			params: UploadSignedCertificateRequest{
				LineageID:    500001,
				GenerationID: 2912,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500001/generations/2912",
			expectedRequestBody: `{"algorithms":{"ECDSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-07-07T10:18:46Z",
						"keyType": "RSA",
						"signedCertificateIssuer": "CN=Test Certificate Authority",
						"signedCertificateNotValidAfterDate": "2027-07-08T13:06:53Z",
						"signedCertificateNotValidBeforeDate": "2026-07-08T13:06:53Z",
						"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01",
						"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:01"
					},
					{
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-07-07T10:18:46Z",
						"keyType": "ECDSA",
						"signedCertificateIssuer": "CN=Test Certificate Authority",
						"signedCertificateNotValidAfterDate": "2027-07-08T13:23:37Z",
						"signedCertificateNotValidBeforeDate": "2026-07-08T13:23:37Z",
						"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:02",
						"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:02"
					}
				],
				"generationId": 2912,
				"generationModifiedBy": "terraform-dev",
				"generationModifiedTime": "2026-07-08T13:24:17Z",
				"generationStatus": "READY_FOR_USE",
				"lineageId": 500001
			}`,
			expectedResponse: &UploadSignedCertificateResponse{
				Algorithms: []Algorithm{
					{
						CertificateStatus:                   CertificateStatusReadyForUse,
						CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
						KeyType:                             string(CryptographicAlgorithmRSA),
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-08T13:06:53Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-08T13:06:53Z")),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:01"),
					},
					{
						CertificateStatus:                   CertificateStatusReadyForUse,
						CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-07-07T10:18:46Z")),
						KeyType:                             string(CryptographicAlgorithmECDSA),
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-08T13:23:37Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-08T13:23:37Z")),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:02"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:02"),
					},
				},
				GenerationID:           2912,
				GenerationModifiedBy:   "terraform-dev",
				GenerationModifiedTime: test.NewTimeFromString(t, "2026-07-08T13:24:17Z"),
				GenerationStatus:       string(GenerationStatusReadyForUse),
				LineageID:              500001,
			},
		},
		"409 cert already uploaded": {
			params:              baseRequest,
			expectedPath:        "/ccm/v2/lineages/500001/generations/2912?acknowledgeWarnings=true",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/cert-already-uploaded",
				"title": "A signed certificate has already been uploaded for this algorithm instance.",
				"status": 409,
				"detail": "Algorithm instance {RSA} on generation {2912} already has an accepted signed certificate. Re-upload is not permitted.",
				"instance": "/error-types/cert-already-uploaded?traceId=1234567891012",
				"context": {"generationId": 2912, "keyType": "RSA"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/cert-already-uploaded",
					Title:    "A signed certificate has already been uploaded for this algorithm instance.",
					Status:   http.StatusConflict,
					Detail:   "Algorithm instance {RSA} on generation {2912} already has an accepted signed certificate. Re-upload is not permitted.",
					Instance: "/error-types/cert-already-uploaded?traceId=1234567891012",
					Context:  map[string]any{"generationId": float64(2912), "keyType": "RSA"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertAlreadyUploaded)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"409 generation immutable": {
			params: UploadSignedCertificateRequest{
				LineageID:    500022,
				GenerationID: 2448,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500022/generations/2448",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/generation-immutable",
				"title": "Generation is immutable and cannot be modified.",
				"status": 409,
				"detail": "Generation {2448} is immutable and cannot be modified.",
				"instance": "/error-types/generation-immutable?traceId=1234567891056",
				"context": {"generationId": 2448}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/generation-immutable",
					Title:    "Generation is immutable and cannot be modified.",
					Status:   http.StatusConflict,
					Detail:   "Generation {2448} is immutable and cannot be modified.",
					Instance: "/error-types/generation-immutable?traceId=1234567891056",
					Context:  map[string]any{"generationId": float64(2448)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrGenerationImmutable)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"400 cert expiry invalid": {
			params:              baseRequest,
			expectedPath:        "/ccm/v2/lineages/500001/generations/2912?acknowledgeWarnings=true",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/cert-expiry-invalid",
				"title": "Uploaded certificate validity dates are invalid.",
				"status": 400,
				"detail": "Cert for key type {RSA} has invalid validity dates: {Certificate validity period of 823 days exceeds the maximum allowed 398 days.}",
				"instance": "/error-types/cert-expiry-invalid?traceId=1234567891028",
				"context": {"keyType": "RSA", "reason": "Certificate validity period of 823 days exceeds the maximum allowed 398 days."}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/cert-expiry-invalid",
					Title:    "Uploaded certificate validity dates are invalid.",
					Status:   http.StatusBadRequest,
					Detail:   "Cert for key type {RSA} has invalid validity dates: {Certificate validity period of 823 days exceeds the maximum allowed 398 days.}",
					Instance: "/error-types/cert-expiry-invalid?traceId=1234567891028",
					Context:  map[string]any{"keyType": "RSA", "reason": "Certificate validity period of 823 days exceeds the maximum allowed 398 days."},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertExpiryInvalid)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"404 lineage not found": {
			params: UploadSignedCertificateRequest{
				LineageID:    999999,
				GenerationID: 2912,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/999999/generations/2912",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891029",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891029",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"404 generation not found": {
			params: UploadSignedCertificateRequest{
				LineageID:    500001,
				GenerationID: 109,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500001/generations/109",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/certificate-not-found",
				"title": "Certificate is not found.",
				"status": 404,
				"detail": "Certificate with {generationId}: {109} is not found.",
				"instance": "/error-types/certificate-not-found?traceId=1234567891030",
				"context": {"certificateIdentifier": "generationId", "certificateIdentifierValue": "109"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/certificate-not-found",
					Title:    "Certificate is not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate with {generationId}: {109} is not found.",
					Instance: "/error-types/certificate-not-found?traceId=1234567891030",
					Context:  map[string]any{"certificateIdentifier": "generationId", "certificateIdentifierValue": "109"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertificateNotFound)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"400 cert parse error": {
			params:              baseRequest,
			expectedPath:        "/ccm/v2/lineages/500001/generations/2912?acknowledgeWarnings=true",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/cert-parse-error",
				"title": "Uploaded certificate is malformed or invalid.",
				"status": 400,
				"detail": "Cert for key type {RSA} is invalid and cannot be parsed: {Failed to parse signed certificate: Failed to parse PEM data}",
				"instance": "/error-types/cert-parse-error?traceId=1234567891031",
				"context": {"keyType": "RSA", "reason": "Failed to parse signed certificate: Failed to parse PEM data"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/cert-parse-error",
					Title:    "Uploaded certificate is malformed or invalid.",
					Status:   http.StatusBadRequest,
					Detail:   "Cert for key type {RSA} is invalid and cannot be parsed: {Failed to parse signed certificate: Failed to parse PEM data}",
					Instance: "/error-types/cert-parse-error?traceId=1234567891031",
					Context:  map[string]any{"keyType": "RSA", "reason": "Failed to parse signed certificate: Failed to parse PEM data"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertParseError)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"400 domain not validated upload failed": {
			params:              baseRequest,
			expectedPath:        "/ccm/v2/lineages/500001/generations/2912?acknowledgeWarnings=true",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/domain-not-validated-upload-failed",
				"title": "Certificate upload failed: Some SANs are not Domain Validated.",
				"status": 400,
				"detail": "Domain validation must be completed for all SANs before uploading the signed certificate. Sample unvalidated domains: {example.com, www.example.com}",
				"instance": "/error-types/domain-not-validated-upload-failed?traceId=1234567891032",
				"context": {"domValidatorLink": "/domain-validation/v1/domains", "domains": "example.com, www.example.com"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/domain-not-validated-upload-failed",
					Title:    "Certificate upload failed: Some SANs are not Domain Validated.",
					Status:   http.StatusBadRequest,
					Detail:   "Domain validation must be completed for all SANs before uploading the signed certificate. Sample unvalidated domains: {example.com, www.example.com}",
					Instance: "/error-types/domain-not-validated-upload-failed?traceId=1234567891032",
					Context:  map[string]any{"domValidatorLink": "/domain-validation/v1/domains", "domains": "example.com, www.example.com"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrDomainNotValidatedUploadFailed)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"400 cert csr mismatch": {
			params: UploadSignedCertificateRequest{
				LineageID:    500005,
				GenerationID: 929,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500005/generations/929",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/cert-csr-mismatch",
				"title": "Certificate does not match stored CSR.",
				"status": 400,
				"detail": "Cert for key type {RSA} does not match stored CSR: {Signed certificate public key does not match the CSR public key on record.}",
				"instance": "/error-types/cert-csr-mismatch?traceId=1234567891033",
				"context": {"generationId": 929, "keyType": "RSA", "lineageId": 500005, "reason": "Signed certificate public key does not match the CSR public key on record."}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/cert-csr-mismatch",
					Title:    "Certificate does not match stored CSR.",
					Status:   http.StatusBadRequest,
					Detail:   "Cert for key type {RSA} does not match stored CSR: {Signed certificate public key does not match the CSR public key on record.}",
					Instance: "/error-types/cert-csr-mismatch?traceId=1234567891033",
					Context:  map[string]any{"generationId": float64(929), "keyType": "RSA", "lineageId": float64(500005), "reason": "Signed certificate public key does not match the CSR public key on record."},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertCSRMismatch)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"400 unknown key type": {
			params: UploadSignedCertificateRequest{
				LineageID:    500005,
				GenerationID: 929,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500005/generations/929",
			expectedRequestBody: `{"algorithms":{"ECDSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/unknown-key-type",
				"title": "Key type is not part of this lineage.",
				"status": 400,
				"detail": "Key type {ECDSA} is not part of this lineage.",
				"instance": "/error-types/unknown-key-type?traceId=1234567891034",
				"context": {"generationId": 929, "keyType": "ECDSA", "lineageId": 500005}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUploadSignedCertificate, &Error{
					Type:     "/error-types/unknown-key-type",
					Title:    "Key type is not part of this lineage.",
					Status:   http.StatusBadRequest,
					Detail:   "Key type {ECDSA} is not part of this lineage.",
					Instance: "/error-types/unknown-key-type?traceId=1234567891034",
					Context:  map[string]any{"generationId": float64(929), "keyType": "ECDSA", "lineageId": float64(500005)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrUnknownKeyType)
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
			},
		},
		"validation error - missing LineageID": {
			params: UploadSignedCertificateRequest{
				GenerationID: 2912,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "uploading signed certificate: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing GenerationID": {
			params: UploadSignedCertificateRequest{
				LineageID: 500001,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "uploading signed certificate: struct validation: GenerationID: cannot be blank")
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing algorithms": {
			params: UploadSignedCertificateRequest{
				LineageID:    500001,
				GenerationID: 2912,
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "uploading signed certificate: struct validation: Body: {\n\tAlgorithms: cannot be blank\n}")
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - trust chain without signed certificate": {
			params: UploadSignedCertificateRequest{
				LineageID:    500001,
				GenerationID: 2912,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {TrustChainPEM: rsaCertPEM},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "uploading signed certificate: struct validation: Body: {\n\tAlgorithms: {\n\t\tRSA: {\n\t\t\tSignedCertificatePEM: cannot be blank\n\t\t}\n\t}\n}")
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid autoActivate value": {
			params: UploadSignedCertificateRequest{
				LineageID:    500001,
				GenerationID: 2912,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
					AutoActivate: []TargetNetwork{"NOT_A_REAL_NETWORK"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "uploading signed certificate: struct validation: Body: {\n\t0: value 'NOT_A_REAL_NETWORK' is invalid. Must be either 'STAGING' or 'PRODUCTION'\n}")
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - empty algorithm key": {
			params: UploadSignedCertificateRequest{
				LineageID:    500001,
				GenerationID: 2912,
				Body: UploadSignedCertificateRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						"": {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "uploading signed certificate: struct validation: Body: {\n\tAlgorithms: must not contain an empty key\n}")
				assert.ErrorIs(t, err, ErrUploadSignedCertificate)
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
				assert.Equal(t, mediaTypeCertificateLineageUploadV3, r.Header.Get("Accept"))
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
			result, err := client.UploadSignedCertificate(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestRenewLineage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           RenewLineageRequest
		responseStatus   int
		responseBody     string
		expectedResponse *RenewLineageResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"201 Created - renew lineage": {
			params:       RenewLineageRequest{LineageID: 500022, ConfirmAbandonHead: true},
			expectedPath: "/ccm/v2/lineages/500022/generations?confirmAbandonHead=true",
			responseBody: `{
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
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T08:38:52Z",
							"algorithmInstanceId": 4545,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T08:38:52Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T08:38:52Z",
							"algorithmInstanceId": 4544,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T08:38:52Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
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
			}`,
			responseStatus: http.StatusCreated,
			expectedResponse: &RenewLineageResponse{
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
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T08:38:52Z")),
								AlgorithmInstanceID:          4545,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T08:38:52Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T08:38:52Z")),
								AlgorithmInstanceID:          4544,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T08:38:52Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
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
		"201 Created - renew lineage without confirmAbandonHead": {
			params:       RenewLineageRequest{LineageID: 500010},
			expectedPath: "/ccm/v2/lineages/500010/generations",
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T10:15:33Z",
							"algorithmInstanceId": 4571,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T10:15:32Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-07-27T10:15:33Z",
							"algorithmInstanceId": 4572,
							"certificateStatus": "CSR_READY",
							"csrExpirationDate": "2027-09-28T10:15:32Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-07-27T10:15:33Z",
					"generationModifiedBy": "terraform-dev",
					"headGenerationId": 3185,
					"headGenerationStatus": "CSR_READY"
				},
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-27T10:15:06Z",
				"lineageId": 500010,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-27T10:15:33Z",
				"lineageName": "renew 1",
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
			responseStatus: http.StatusCreated,
			expectedResponse: &RenewLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				GeoClass:   string(GeoClassStandardWorldwide),
				GroupID:    12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T10:15:33Z")),
								AlgorithmInstanceID:          4571,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T10:15:32Z")),
								CSRPEM:                       rsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmRSA),
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T10:15:33Z")),
								AlgorithmInstanceID:          4572,
								CertificateStatus:            CertificateStatusCSRReady,
								CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T10:15:32Z")),
								CSRPEM:                       ecdsaCSRPEM,
								KeyType:                      string(CryptographicAlgorithmECDSA),
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T10:15:33Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					HeadGenerationID:     3185,
					HeadGenerationStatus: string(GenerationStatusCSRReady),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-27T10:15:06Z"),
				LineageID:           500010,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-27T10:15:33Z"),
				LineageName:         "renew 1",
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
		"400 renew blocked head on staging": {
			params:         RenewLineageRequest{LineageID: 500022},
			expectedPath:   "/ccm/v2/lineages/500022/generations",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/renew-blocked-head-on-staging",
				"title": "Renew is not allowed in the current lineage state.",
				"status": 400,
				"detail": "Cannot renew lineage {500022} while a head generation and a staging generation are both present. Promote or abandon the head generation first.",
				"instance": "/error-types/renew-blocked-head-on-staging?traceId=1234567891041",
				"context": {"lineageId": 500022}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRenewLineage, &Error{
					Type:     "/error-types/renew-blocked-head-on-staging",
					Title:    "Renew is not allowed in the current lineage state.",
					Status:   http.StatusBadRequest,
					Detail:   "Cannot renew lineage {500022} while a head generation and a staging generation are both present. Promote or abandon the head generation first.",
					Instance: "/error-types/renew-blocked-head-on-staging?traceId=1234567891041",
					Context:  map[string]any{"lineageId": float64(500022)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrRenewBlockedHeadOnStaging)
				assert.ErrorIs(t, err, ErrRenewLineage)
			},
		},
		"409 confirmation required to abandon head": {
			params:         RenewLineageRequest{LineageID: 500010},
			expectedPath:   "/ccm/v2/lineages/500010/generations",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/confirmation-required",
				"title": "Confirmation required to abandon head generation.",
				"status": 409,
				"detail": "Head generation for lineage {500010} has a valid certificate. Set confirmAbandonHead=true to proceed.",
				"instance": "/error-types/confirmation-required?traceId=1234567891043",
				"context": {"lineageId": 500010}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRenewLineage, &Error{
					Type:     "/error-types/confirmation-required",
					Title:    "Confirmation required to abandon head generation.",
					Status:   http.StatusConflict,
					Detail:   "Head generation for lineage {500010} has a valid certificate. Set confirmAbandonHead=true to proceed.",
					Instance: "/error-types/confirmation-required?traceId=1234567891043",
					Context:  map[string]any{"lineageId": float64(500010)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrConfirmationRequired)
				assert.ErrorIs(t, err, ErrRenewLineage)
			},
		},
		"409 single-generation lineage cannot be renewed": {
			params:         RenewLineageRequest{LineageID: 500011},
			expectedPath:   "/ccm/v2/lineages/500011/generations",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/single-generation-no-renew",
				"title": "Single-generation lineages cannot be renewed without a production certificate.",
				"status": 409,
				"detail": "Lineage {500011} is SINGLE_GENERATION type with no production certificate. Deploy to production before renewing.",
				"instance": "/error-types/single-generation-no-renew?traceId=1234567891044",
				"context": {"lineageId": 500011}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRenewLineage, &Error{
					Type:     "/error-types/single-generation-no-renew",
					Title:    "Single-generation lineages cannot be renewed without a production certificate.",
					Status:   http.StatusConflict,
					Detail:   "Lineage {500011} is SINGLE_GENERATION type with no production certificate. Deploy to production before renewing.",
					Instance: "/error-types/single-generation-no-renew?traceId=1234567891044",
					Context:  map[string]any{"lineageId": float64(500011)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrSingleGenerationNoRenew)
				assert.ErrorIs(t, err, ErrRenewLineage)
			},
		},
		"404 lineage not found": {
			params:         RenewLineageRequest{LineageID: 999999},
			expectedPath:   "/ccm/v2/lineages/999999/generations",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891045",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRenewLineage, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891045",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrRenewLineage)
			},
		},
		"500 internal server error": {
			params:         RenewLineageRequest{LineageID: 500022},
			expectedPath:   "/ccm/v2/lineages/500022/generations",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891042"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrRenewLineage, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891042",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrRenewLineage)
			},
		},
		"validation error - missing LineageID": {
			params: RenewLineageRequest{},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "renewing lineage: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrRenewLineage)
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
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.RenewLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestCompleteLineage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params              CompleteLineageRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedResponse    *CompleteLineageResponse
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"201 Created - complete lineage": {
			params: CompleteLineageRequest{
				LineageID: 500022,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500022/generations/complete",
			expectedRequestBody: `{"algorithms":{"ECDSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"currentProduction": {
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
					"productionGenerationId": 3029,
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
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:08",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:08"
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
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-22T09:21:09Z",
				"lineageId": 500022,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-23T13:42:30Z",
				"lineageName": "alpha 2",
				"lineageType": "MULTIPLE_GENERATION",
				"previousProduction": {
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
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:08",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:08"
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
					"previousProductionGenerationId": 3028,
					"previousProductionGenerationStatus": "ACTIVE"
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
			expectedResponse: &CompleteLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				CurrentProduction: &ProductionGeneration{
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
					ProductionGenerationID:     3029,
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
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:08"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:08"),
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
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-23T13:42:30Z"),
				LineageName:         "alpha 2",
				LineageType:         string(LineageTypeMultipleGeneration),
				PreviousProduction: &PreviousProductionGeneration{
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
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:08"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:08"),
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
					PreviousProductionGenerationID:     3028,
					PreviousProductionGenerationStatus: string(GenerationStatusActive),
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
		"201 Created - complete lineage with autoActivate STAGING and PRODUCTION": {
			params: CompleteLineageRequest{
				LineageID: 500023,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
					AutoActivate: []TargetNetwork{TargetNetworkStaging, TargetNetworkProduction},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500023/generations/complete",
			expectedRequestBody: `{"algorithms":{"ECDSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"}},"autoActivate":["STAGING","PRODUCTION"]}`,
			responseStatus:      http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"currentProduction": {
					"productionGenerationId": 3040,
					"productionGenerationStatus": "ACTIVE"
				},
				"currentStaging": {
					"stagingGenerationId": 3040,
					"stagingGenerationStatus": "ACTIVE"
				},
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"keySpecs": [
					{"keySize": "2048", "keyType": "RSA"},
					{"keySize": "P-256", "keyType": "ECDSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-07-23T11:26:01Z",
				"lineageId": 500023,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-07-23T13:42:30Z",
				"lineageName": "beta 1",
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
			expectedResponse: &CompleteLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				CurrentProduction: &ProductionGeneration{
					ProductionGenerationID:     3040,
					ProductionGenerationStatus: string(GenerationStatusActive),
				},
				CurrentStaging: &StagingGeneration{
					StagingGenerationID:     3040,
					StagingGenerationStatus: string(GenerationStatusActive),
				},
				GeoClass: string(GeoClassStandardWorldwide),
				GroupID:  12345,
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-07-23T11:26:01Z"),
				LineageID:           500023,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-07-23T13:42:30Z"),
				LineageName:         "beta 1",
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
		"201 Created - complete lineage with acknowledgeWarnings=true after warnings": {
			params: CompleteLineageRequest{
				LineageID:           500012,
				AcknowledgeWarnings: true,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500012/generations/complete?acknowledgeWarnings=true",
			expectedRequestBody: `{"algorithms":{"ECDSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusCreated,
			responseBody: `{
				"accountId": "A-CCT1234",
				"contractId": "C-0N7RAC7",
				"currentProduction": {
					"productionGenerationId": 3208,
					"productionGenerationStatus": "READY_FOR_USE"
				},
				"currentStaging": {
					"stagingGenerationId": 3208,
					"stagingGenerationStatus": "READY_FOR_USE"
				},
				"geoClass": "STANDARD_WORLDWIDE",
				"groupId": 12345,
				"head": {
					"algorithms": [
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-09-17T12:47:52Z",
							"algorithmInstanceId": 4401,
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-11-19T12:12:15Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "ECDSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-09-17T12:41:46Z",
							"signedCertificateNotValidBeforeDate": "2026-09-17T12:41:46Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:09",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:09"
						},
						{
							"algorithmInstanceCreatedBy": "terraform-dev",
							"algorithmInstanceCreatedTime": "2026-09-17T12:47:52Z",
							"algorithmInstanceId": 4402,
							"certificateStatus": "READY_FOR_USE",
							"csrExpirationDate": "2027-11-19T12:12:15Z",
							"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
							"keyType": "RSA",
							"signedCertificateIssuer": "CN=Test Certificate Authority",
							"signedCertificateNotValidAfterDate": "2027-09-17T12:14:26Z",
							"signedCertificateNotValidBeforeDate": "2026-09-17T12:14:26Z",
							"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
							"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:10",
							"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:10"
						}
					],
					"generationCreatedBy": "terraform-dev",
					"generationCreatedTime": "2026-09-17T12:47:52Z",
					"generationModifiedBy": "terraform-dev",
					"generationModifiedTime": "2026-09-17T12:47:52Z",
					"headGenerationId": 3210,
					"headGenerationStatus": "READY_FOR_USE"
				},
				"keySpecs": [
					{"keySize": "P-256", "keyType": "ECDSA"},
					{"keySize": "2048", "keyType": "RSA"}
				],
				"lineageCreatedBy": "terraform-dev",
				"lineageCreatedTime": "2026-09-17T12:12:16Z",
				"lineageId": 500012,
				"lineageModifiedBy": "terraform-dev",
				"lineageModifiedTime": "2026-09-17T12:20:30Z",
				"lineageName": "gamma 1",
				"lineageType": "MULTIPLE_GENERATION",
				"sans": ["www.example.com", "example.com"],
				"secureNetwork": "ENHANCED_TLS",
				"stackMode": "MULTIPLE_STACK",
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "Example Inc.",
					"organizationalUnit": "IT",
					"state": "Massachusetts"
				}
			}`,
			expectedResponse: &CompleteLineageResponse{
				AccountID:  "A-CCT1234",
				ContractID: "C-0N7RAC7",
				CurrentProduction: &ProductionGeneration{
					ProductionGenerationID:     3208,
					ProductionGenerationStatus: string(GenerationStatusReadyForUse),
				},
				CurrentStaging: &StagingGeneration{
					StagingGenerationID:     3208,
					StagingGenerationStatus: string(GenerationStatusReadyForUse),
				},
				GeoClass: string(GeoClassStandardWorldwide),
				GroupID:  12345,
				Head: &HeadGeneration{
					Generation: Generation{
						Algorithms: []Algorithm{
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-09-17T12:47:52Z")),
								AlgorithmInstanceID:                 4401,
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-11-19T12:12:15Z")),
								CSRPEM:                              ecdsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmECDSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-09-17T12:41:46Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-09-17T12:41:46Z")),
								SignedCertificatePEM:                ptr.To(ecdsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:09"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:09"),
							},
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-09-17T12:47:52Z")),
								AlgorithmInstanceID:                 4402,
								CertificateStatus:                   CertificateStatusReadyForUse,
								CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-11-19T12:12:15Z")),
								CSRPEM:                              rsaCSRPEM,
								KeyType:                             string(CryptographicAlgorithmRSA),
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-09-17T12:14:26Z")),
								SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-09-17T12:14:26Z")),
								SignedCertificatePEM:                ptr.To(rsaCertPEM),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:10"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:10"),
							},
						},
						GenerationCreatedBy:    ptr.To("terraform-dev"),
						GenerationCreatedTime:  ptr.To(test.NewTimeFromString(t, "2026-09-17T12:47:52Z")),
						GenerationModifiedBy:   ptr.To("terraform-dev"),
						GenerationModifiedTime: ptr.To(test.NewTimeFromString(t, "2026-09-17T12:47:52Z")),
					},
					HeadGenerationID:     3210,
					HeadGenerationStatus: string(GenerationStatusReadyForUse),
				},
				KeySpecs: []KeySpecResponse{
					{KeySize: string(KeySizeP256), KeyType: string(CryptographicAlgorithmECDSA)},
					{KeySize: string(KeySize2048), KeyType: string(CryptographicAlgorithmRSA)},
				},
				LineageCreatedBy:    "terraform-dev",
				LineageCreatedTime:  test.NewTimeFromString(t, "2026-09-17T12:12:16Z"),
				LineageID:           500012,
				LineageModifiedBy:   "terraform-dev",
				LineageModifiedTime: test.NewTimeFromString(t, "2026-09-17T12:20:30Z"),
				LineageName:         "gamma 1",
				LineageType:         string(LineageTypeMultipleGeneration),
				SANs:                []string{"www.example.com", "example.com"},
				SecureNetwork:       string(SecureNetworkEnhancedTLS),
				StackMode:           string(StackModeMultipleStack),
				Subject: Subject{
					CommonName:         "example.com",
					Country:            "US",
					Locality:           "Cambridge",
					Organization:       "Example Inc.",
					OrganizationalUnit: "IT",
					State:              "Massachusetts",
				},
			},
		},
		"500 internal server error": {
			params: CompleteLineageRequest{
				LineageID: 500022,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:   "/ccm/v2/lineages/500022/generations/complete",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891046"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCompleteLineage, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891046",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrCompleteLineage)
			},
		},
		"404 lineage not found": {
			params: CompleteLineageRequest{
				LineageID: 999999,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:   "/ccm/v2/lineages/999999/generations/complete",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891052",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCompleteLineage, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891052",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrCompleteLineage)
			},
		},
		"409 no current production": {
			params: CompleteLineageRequest{
				LineageID: 500013,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:   "/ccm/v2/lineages/500013/generations/complete",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/no-current-production",
				"title": "No current production generation exists.",
				"status": 409,
				"detail": "Lineage {500013} has no current production generation to complete.",
				"instance": "/error-types/no-current-production?traceId=1234567891051",
				"context": {"lineageId": 500013}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCompleteLineage, &Error{
					Type:     "/error-types/no-current-production",
					Title:    "No current production generation exists.",
					Status:   http.StatusConflict,
					Detail:   "Lineage {500013} has no current production generation to complete.",
					Instance: "/error-types/no-current-production?traceId=1234567891051",
					Context:  map[string]any{"lineageId": float64(500013)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrNoCurrentProduction)
				assert.ErrorIs(t, err, ErrCompleteLineage)
			},
		},
		"409 complete precondition failed - head already exists": {
			params: CompleteLineageRequest{
				LineageID: 500022,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:   "/ccm/v2/lineages/500022/generations/complete",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/complete-precondition-failed",
				"title": "Complete operation precondition failed.",
				"status": 409,
				"detail": "Lineage {500022} failed complete precondition: {HEAD_ALREADY_EXISTS}",
				"instance": "/error-types/complete-precondition-failed?traceId=1234567891047",
				"context": {"lineageId": 500022, "reason": "HEAD_ALREADY_EXISTS"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCompleteLineage, &Error{
					Type:     "/error-types/complete-precondition-failed",
					Title:    "Complete operation precondition failed.",
					Status:   http.StatusConflict,
					Detail:   "Lineage {500022} failed complete precondition: {HEAD_ALREADY_EXISTS}",
					Instance: "/error-types/complete-precondition-failed?traceId=1234567891047",
					Context:  map[string]any{"lineageId": float64(500022), "reason": "HEAD_ALREADY_EXISTS"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCompletePreconditionFailed)
				assert.ErrorIs(t, err, ErrCompleteLineage)
			},
		},
		"409 complete precondition failed - key type mismatch": {
			params: CompleteLineageRequest{
				LineageID: 500012,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA: {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500012/generations/complete",
			expectedRequestBody: `{"algorithms":{"RSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/complete-precondition-failed",
				"title": "Complete operation precondition failed.",
				"status": 409,
				"detail": "Lineage {500012} failed complete precondition: {KEY_TYPE_MISMATCH}",
				"instance": "/error-types/complete-precondition-failed?traceId=1234567891049",
				"context": {"lineageId": 500012, "reason": "KEY_TYPE_MISMATCH"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCompleteLineage, &Error{
					Type:     "/error-types/complete-precondition-failed",
					Title:    "Complete operation precondition failed.",
					Status:   http.StatusConflict,
					Detail:   "Lineage {500012} failed complete precondition: {KEY_TYPE_MISMATCH}",
					Instance: "/error-types/complete-precondition-failed?traceId=1234567891049",
					Context:  map[string]any{"lineageId": float64(500012), "reason": "KEY_TYPE_MISMATCH"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCompletePreconditionFailed)
				assert.ErrorIs(t, err, ErrCompleteLineage)
			},
		},
		"400 cert parse error": {
			params: CompleteLineageRequest{
				LineageID: 500012,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500012/generations/complete",
			expectedRequestBody: `{"algorithms":{"ECDSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/cert-parse-error",
				"title": "Uploaded certificate is malformed or invalid.",
				"status": 400,
				"detail": "Cert for key type {ECDSA} is invalid and cannot be parsed: {Failed to parse signed certificate: Failed to parse PEM data}",
				"instance": "/error-types/cert-parse-error?traceId=1234567891050",
				"context": {"keyType": "ECDSA", "reason": "Failed to parse signed certificate: Failed to parse PEM data"}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCompleteLineage, &Error{
					Type:     "/error-types/cert-parse-error",
					Title:    "Uploaded certificate is malformed or invalid.",
					Status:   http.StatusBadRequest,
					Detail:   "Cert for key type {ECDSA} is invalid and cannot be parsed: {Failed to parse signed certificate: Failed to parse PEM data}",
					Instance: "/error-types/cert-parse-error?traceId=1234567891050",
					Context:  map[string]any{"keyType": "ECDSA", "reason": "Failed to parse signed certificate: Failed to parse PEM data"},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertParseError)
				assert.ErrorIs(t, err, ErrCompleteLineage)
			},
		},
		"400 cert csr mismatch": {
			params: CompleteLineageRequest{
				LineageID: 500012,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500012/generations/complete",
			expectedRequestBody: `{"algorithms":{"ECDSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusBadRequest,
			responseBody: `{
				"type": "/error-types/cert-csr-mismatch",
				"title": "Certificate does not match stored CSR.",
				"status": 400,
				"detail": "Cert for key type {ECDSA} does not match stored CSR: {Signed certificate public key does not match the CSR public key on record.}",
				"instance": "/error-types/cert-csr-mismatch?traceId=1234567891048",
				"context": {"generationId": 3203, "keyType": "ECDSA", "lineageId": 500012, "reason": "Signed certificate public key does not match the CSR public key on record."}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCompleteLineage, &Error{
					Type:     "/error-types/cert-csr-mismatch",
					Title:    "Certificate does not match stored CSR.",
					Status:   http.StatusBadRequest,
					Detail:   "Cert for key type {ECDSA} does not match stored CSR: {Signed certificate public key does not match the CSR public key on record.}",
					Instance: "/error-types/cert-csr-mismatch?traceId=1234567891048",
					Context:  map[string]any{"generationId": float64(3203), "keyType": "ECDSA", "lineageId": float64(500012), "reason": "Signed certificate public key does not match the CSR public key on record."},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrCertCSRMismatch)
				assert.ErrorIs(t, err, ErrCompleteLineage)
			},
		},
		"409 lineage upload validation warnings": {
			params: CompleteLineageRequest{
				LineageID: 500012,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			expectedPath:        "/ccm/v2/lineages/500012/generations/complete",
			expectedRequestBody: `{"algorithms":{"ECDSA":{"signedCertificatePem":"-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"}}}`,
			responseStatus:      http.StatusConflict,
			responseBody: `{
				"type": "/error-types/lineage-upload-validation-warnings",
				"title": "Certificate upload was rejected due to warnings that require acknowledgment.",
				"status": 409,
				"detail": "Certificate upload was rejected due to warnings that require acknowledgment.",
				"instance": "/error-types/lineage-upload-validation-warnings?traceId=b2ff3575e95a4d58",
				"context": {
					"validationWarnings": ["ECDSA certificate does not come with a trust chain; this is a non-standard practice."]
				}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCompleteLineage, &Error{
					Type:     "/error-types/lineage-upload-validation-warnings",
					Title:    "Certificate upload was rejected due to warnings that require acknowledgment.",
					Status:   http.StatusConflict,
					Detail:   "Certificate upload was rejected due to warnings that require acknowledgment.",
					Instance: "/error-types/lineage-upload-validation-warnings?traceId=b2ff3575e95a4d58",
					Context: map[string]any{
						"validationWarnings": []any{"ECDSA certificate does not come with a trust chain; this is a non-standard practice."},
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageUploadValidationWarnings)
				assert.ErrorIs(t, err, ErrCompleteLineage)
			},
		},
		"validation error - missing LineageID": {
			params: CompleteLineageRequest{
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "completing lineage: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrCompleteLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing Body": {
			params: CompleteLineageRequest{LineageID: 500022},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "completing lineage: struct validation: Body: {\n\tAlgorithms: cannot be blank\n}")
				assert.ErrorIs(t, err, ErrCompleteLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - more than one algorithm": {
			params: CompleteLineageRequest{
				LineageID: 500022,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmRSA:   {SignedCertificatePEM: rsaCertPEM},
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "completing lineage: struct validation: Body: {\n\tAlgorithms: the length must be exactly 1\n}")
				assert.ErrorIs(t, err, ErrCompleteLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - empty algorithm key": {
			params: CompleteLineageRequest{
				LineageID: 500022,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						"": {SignedCertificatePEM: rsaCertPEM},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "completing lineage: struct validation: Body: {\n\tAlgorithms: must not contain an empty key\n}")
				assert.ErrorIs(t, err, ErrCompleteLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - trust chain without signed certificate": {
			params: CompleteLineageRequest{
				LineageID: 500022,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {TrustChainPEM: ecdsaCertPEM},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "completing lineage: struct validation: Body: {\n\tAlgorithms: {\n\t\tECDSA: {\n\t\t\tSignedCertificatePEM: cannot be blank\n\t\t}\n\t}\n}")
				assert.ErrorIs(t, err, ErrCompleteLineage)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - invalid autoActivate value": {
			params: CompleteLineageRequest{
				LineageID: 500022,
				Body: CompleteLineageRequestBody{
					Algorithms: map[CryptographicAlgorithm]SignedCertificate{
						CryptographicAlgorithmECDSA: {SignedCertificatePEM: ecdsaCertPEM},
					},
					AutoActivate: []TargetNetwork{"NOT_A_REAL_NETWORK"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "completing lineage: struct validation: Body: {\n\t0: value 'NOT_A_REAL_NETWORK' is invalid. Must be either 'STAGING' or 'PRODUCTION'\n}")
				assert.ErrorIs(t, err, ErrCompleteLineage)
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
			result, err := client.CompleteLineage(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetGeneration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           GetGenerationRequest
		responseStatus   int
		responseBody     string
		expectedResponse *GetGenerationResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"200 OK - get generation": {
			params:         GetGenerationRequest{LineageID: 500005, GenerationID: 929},
			expectedPath:   "/ccm/v2/lineages/500005/generations/929",
			responseStatus: http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"algorithmInstanceCreatedBy": "terraform-dev",
						"algorithmInstanceCreatedTime": "2026-07-17T12:28:59Z",
						"algorithmInstanceId": 1714,
						"algorithmInstanceModifiedBy": "terraform-dev",
						"algorithmInstanceModifiedTime": "2026-07-17T12:50:43Z",
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-07-17T12:28:59Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						"keyType": "RSA",
						"signedCertificateIssuer": "CN=Test Certificate Authority",
						"signedCertificateNotValidAfterDate": "2027-07-17T12:50:04Z",
						"signedCertificateNotValidBeforeDate": "2026-07-17T12:50:04Z",
						"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
						"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:03",
						"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:03"
					}
				],
				"generationCreatedBy": "terraform-dev",
				"generationCreatedTime": "2026-07-17T12:28:59Z",
				"generationId": 929,
				"generationModifiedBy": "terraform-dev",
				"generationModifiedTime": "2026-07-17T12:50:44Z",
				"generationStatus": "READY_FOR_USE"
			}`,
			expectedResponse: &GetGenerationResponse{
				Generation: Generation{
					Algorithms: []Algorithm{
						{
							AlgorithmInstanceCreatedBy:          "terraform-dev",
							AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-17T12:28:59Z")),
							AlgorithmInstanceID:                 1714,
							AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
							AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-17T12:50:43Z")),
							CertificateStatus:                   CertificateStatusReadyForUse,
							CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-07-17T12:28:59Z")),
							CSRPEM:                              rsaCSRPEM,
							KeyType:                             string(CryptographicAlgorithmRSA),
							SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
							SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-17T12:50:04Z")),
							SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-17T12:50:04Z")),
							SignedCertificatePEM:                ptr.To(rsaCertPEM),
							SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:03"),
							SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:03"),
						},
					},
					GenerationCreatedBy:    ptr.To("terraform-dev"),
					GenerationCreatedTime:  ptr.To(test.NewTimeFromString(t, "2026-07-17T12:28:59Z")),
					GenerationModifiedBy:   ptr.To("terraform-dev"),
					GenerationModifiedTime: ptr.To(test.NewTimeFromString(t, "2026-07-17T12:50:44Z")),
				},
				GenerationID:     929,
				GenerationStatus: "READY_FOR_USE",
			},
		},
		"200 OK - get generation, head with two pending CSR algorithms": {
			params:         GetGenerationRequest{LineageID: 500014, GenerationID: 3218},
			expectedPath:   "/ccm/v2/lineages/500014/generations/3218",
			responseStatus: http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"algorithmInstanceCreatedBy": "terraform-dev",
						"algorithmInstanceCreatedTime": "2026-07-27T11:26:38Z",
						"algorithmInstanceId": 4632,
						"certificateStatus": "CSR_READY",
						"csrExpirationDate": "2027-09-28T11:26:38Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						"keyType": "ECDSA"
					},
					{
						"algorithmInstanceCreatedBy": "terraform-dev",
						"algorithmInstanceCreatedTime": "2026-07-27T11:26:38Z",
						"algorithmInstanceId": 4631,
						"certificateStatus": "CSR_READY",
						"csrExpirationDate": "2027-09-28T11:26:38Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						"keyType": "RSA"
					}
				],
				"generationCreatedBy": "terraform-dev",
				"generationCreatedTime": "2026-07-27T11:26:38Z",
				"generationId": 3218,
				"generationModifiedBy": "terraform-dev",
				"generationStatus": "CSR_READY"
			}`,
			expectedResponse: &GetGenerationResponse{
				Generation: Generation{
					Algorithms: []Algorithm{
						{
							AlgorithmInstanceCreatedBy:   "terraform-dev",
							AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T11:26:38Z")),
							AlgorithmInstanceID:          4632,
							CertificateStatus:            CertificateStatusCSRReady,
							CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T11:26:38Z")),
							CSRPEM:                       ecdsaCSRPEM,
							KeyType:                      string(CryptographicAlgorithmECDSA),
						},
						{
							AlgorithmInstanceCreatedBy:   "terraform-dev",
							AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T11:26:38Z")),
							AlgorithmInstanceID:          4631,
							CertificateStatus:            CertificateStatusCSRReady,
							CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T11:26:38Z")),
							CSRPEM:                       rsaCSRPEM,
							KeyType:                      string(CryptographicAlgorithmRSA),
						},
					},
					GenerationCreatedBy:   ptr.To("terraform-dev"),
					GenerationCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T11:26:38Z")),
					GenerationModifiedBy:  ptr.To("terraform-dev"),
				},
				GenerationID:     3218,
				GenerationStatus: "CSR_READY",
			},
		},
		"200 OK - get generation, mixed CSR_READY and READY_FOR_USE algorithms": {
			params:         GetGenerationRequest{LineageID: 500015, GenerationID: 3191},
			expectedPath:   "/ccm/v2/lineages/500015/generations/3191",
			responseStatus: http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"algorithmInstanceCreatedBy": "terraform-dev",
						"algorithmInstanceCreatedTime": "2026-07-27T10:24:20Z",
						"algorithmInstanceId": 4583,
						"certificateStatus": "CSR_READY",
						"csrExpirationDate": "2027-09-28T10:24:20Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						"keyType": "ECDSA"
					},
					{
						"algorithmInstanceCreatedBy": "terraform-dev",
						"algorithmInstanceCreatedTime": "2026-07-27T10:24:20Z",
						"algorithmInstanceId": 4582,
						"algorithmInstanceModifiedBy": "terraform-dev",
						"algorithmInstanceModifiedTime": "2026-07-27T11:21:51Z",
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-09-28T10:24:20Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						"keyType": "RSA",
						"signedCertificateIssuer": "CN=Test Certificate Authority",
						"signedCertificateNotValidAfterDate": "2027-07-27T11:21:21Z",
						"signedCertificateNotValidBeforeDate": "2026-07-27T11:21:21Z",
						"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
						"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:09",
						"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:09"
					}
				],
				"generationCreatedBy": "terraform-dev",
				"generationCreatedTime": "2026-07-27T10:24:20Z",
				"generationId": 3191,
				"generationModifiedBy": "terraform-dev",
				"generationModifiedTime": "2026-07-27T11:21:52Z",
				"generationStatus": "READY_FOR_USE"
			}`,
			expectedResponse: &GetGenerationResponse{
				Generation: Generation{
					Algorithms: []Algorithm{
						{
							AlgorithmInstanceCreatedBy:   "terraform-dev",
							AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T10:24:20Z")),
							AlgorithmInstanceID:          4583,
							CertificateStatus:            CertificateStatusCSRReady,
							CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-28T10:24:20Z")),
							CSRPEM:                       ecdsaCSRPEM,
							KeyType:                      string(CryptographicAlgorithmECDSA),
						},
						{
							AlgorithmInstanceCreatedBy:          "terraform-dev",
							AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-27T10:24:20Z")),
							AlgorithmInstanceID:                 4582,
							AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
							AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-27T11:21:51Z")),
							CertificateStatus:                   CertificateStatusReadyForUse,
							CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-28T10:24:20Z")),
							CSRPEM:                              rsaCSRPEM,
							KeyType:                             string(CryptographicAlgorithmRSA),
							SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
							SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-27T11:21:21Z")),
							SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-27T11:21:21Z")),
							SignedCertificatePEM:                ptr.To(rsaCertPEM),
							SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:09"),
							SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:09"),
						},
					},
					GenerationCreatedBy:    ptr.To("terraform-dev"),
					GenerationCreatedTime:  ptr.To(test.NewTimeFromString(t, "2026-07-27T10:24:20Z")),
					GenerationModifiedBy:   ptr.To("terraform-dev"),
					GenerationModifiedTime: ptr.To(test.NewTimeFromString(t, "2026-07-27T11:21:52Z")),
				},
				GenerationID:     3191,
				GenerationStatus: "READY_FOR_USE",
			},
		},
		"200 OK - get generation, RSA algorithm with a trust chain": {
			params:         GetGenerationRequest{LineageID: 172153, GenerationID: 33490},
			expectedPath:   "/ccm/v2/lineages/172153/generations/33490",
			responseStatus: http.StatusOK,
			responseBody: `{
				"algorithms": [
					{
						"algorithmInstanceCreatedBy": "terraform-dev",
						"algorithmInstanceCreatedTime": "2026-09-11T12:28:48Z",
						"algorithmInstanceId": 61758,
						"algorithmInstanceModifiedBy": "terraform-dev",
						"algorithmInstanceModifiedTime": "2026-09-11T12:29:03Z",
						"certificateStatus": "READY_FOR_USE",
						"csrExpirationDate": "2027-11-13T12:28:48Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						"keyType": "RSA",
						"signedCertificateIssuer": "CN=test.com,O=Akamai",
						"signedCertificateNotValidAfterDate": "2027-09-11T12:28:48Z",
						"signedCertificateNotValidBeforeDate": "2026-09-11T12:28:48Z",
						"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
						"signedCertificateSerialNumber": "da:db:82:98:89:5c:d6:fd:a8:bc:46:3f:1b:29:45:1a",
						"signedCertificateSha256Fingerprint": "43:D9:E5:1D:59:01:22:D9:E7:E8:94:E6:99:F2:09:89:24:17:13:8D:6B:1F:48:AA:95:B4:F2:21:B7:CF:5F:6D",
						"trustChainPem": "-----BEGIN CERTIFICATE-----\nRSA-TRUST-CHAIN\n-----END CERTIFICATE-----\n"
					},
					{
						"algorithmInstanceCreatedBy": "terraform-dev",
						"algorithmInstanceCreatedTime": "2026-09-11T12:28:48Z",
						"algorithmInstanceId": 61757,
						"certificateStatus": "CSR_READY",
						"csrExpirationDate": "2027-11-13T12:28:48Z",
						"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						"keyType": "ECDSA"
					}
				],
				"generationCreatedBy": "terraform-dev",
				"generationCreatedTime": "2026-09-11T12:28:48Z",
				"generationId": 33490,
				"generationModifiedBy": "terraform-dev",
				"generationModifiedTime": "2026-09-11T12:29:03Z",
				"generationStatus": "READY_FOR_USE"
			}`,
			expectedResponse: &GetGenerationResponse{
				Generation: Generation{
					Algorithms: []Algorithm{
						{
							AlgorithmInstanceCreatedBy:          "terraform-dev",
							AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-09-11T12:28:48Z")),
							AlgorithmInstanceID:                 61758,
							AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
							AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-09-11T12:29:03Z")),
							CertificateStatus:                   CertificateStatusReadyForUse,
							CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-11-13T12:28:48Z")),
							CSRPEM:                              rsaCSRPEM,
							KeyType:                             string(CryptographicAlgorithmRSA),
							SignedCertificateIssuer:             ptr.To("CN=test.com,O=Akamai"),
							SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-09-11T12:28:48Z")),
							SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-09-11T12:28:48Z")),
							SignedCertificatePEM:                ptr.To(rsaCertPEM),
							SignedCertificateSerialNumber:       ptr.To("da:db:82:98:89:5c:d6:fd:a8:bc:46:3f:1b:29:45:1a"),
							SignedCertificateSHA256Fingerprint:  ptr.To("43:D9:E5:1D:59:01:22:D9:E7:E8:94:E6:99:F2:09:89:24:17:13:8D:6B:1F:48:AA:95:B4:F2:21:B7:CF:5F:6D"),
							TrustChainPEM:                       ptr.To(rsaTrustChainPEM),
						},
						{
							AlgorithmInstanceCreatedBy:   "terraform-dev",
							AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-09-11T12:28:48Z")),
							AlgorithmInstanceID:          61757,
							CertificateStatus:            CertificateStatusCSRReady,
							CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-11-13T12:28:48Z")),
							CSRPEM:                       ecdsaCSRPEM,
							KeyType:                      string(CryptographicAlgorithmECDSA),
						},
					},
					GenerationCreatedBy:    ptr.To("terraform-dev"),
					GenerationCreatedTime:  ptr.To(test.NewTimeFromString(t, "2026-09-11T12:28:48Z")),
					GenerationModifiedBy:   ptr.To("terraform-dev"),
					GenerationModifiedTime: ptr.To(test.NewTimeFromString(t, "2026-09-11T12:29:03Z")),
				},
				GenerationID:     33490,
				GenerationStatus: "READY_FOR_USE",
			},
		},
		"500 internal server error": {
			params:         GetGenerationRequest{LineageID: 500005, GenerationID: 929},
			expectedPath:   "/ccm/v2/lineages/500005/generations/929",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891020"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetGeneration, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891020",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrGetGeneration)
			},
		},
		"404 generation not found": {
			params:         GetGenerationRequest{LineageID: 500005, GenerationID: 930},
			expectedPath:   "/ccm/v2/lineages/500005/generations/930",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/generation-not-found",
				"title": "Generation not found.",
				"status": 404,
				"detail": "Generation {930} in Lineage {500005} not found.",
				"instance": "/error-types/generation-not-found?traceId=1234567891035",
				"context": {"generationId": 930, "lineageId": 500005}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetGeneration, &Error{
					Type:     "/error-types/generation-not-found",
					Title:    "Generation not found.",
					Status:   http.StatusNotFound,
					Detail:   "Generation {930} in Lineage {500005} not found.",
					Instance: "/error-types/generation-not-found?traceId=1234567891035",
					Context:  map[string]any{"generationId": float64(930), "lineageId": float64(500005)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrGenerationNotFound)
				assert.ErrorIs(t, err, ErrGetGeneration)
			},
		},
		"404 lineage not found": {
			params:         GetGenerationRequest{LineageID: 999999, GenerationID: 929},
			expectedPath:   "/ccm/v2/lineages/999999/generations/929",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891036",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetGeneration, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891036",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrGetGeneration)
			},
		},
		"validation error - missing LineageID": {
			params: GetGenerationRequest{GenerationID: 929},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "getting generation: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrGetGeneration)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing GenerationID": {
			params: GetGenerationRequest{LineageID: 500005},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "getting generation: struct validation: GenerationID: cannot be blank")
				assert.ErrorIs(t, err, ErrGetGeneration)
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
			result, err := client.GetGeneration(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestDeleteGeneration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params         DeleteGenerationRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No Content - delete generation": {
			params:         DeleteGenerationRequest{LineageID: 500017, GenerationID: 3250},
			expectedPath:   "/ccm/v2/lineages/500017/generations/3250",
			responseStatus: http.StatusNoContent,
		},
		"204 No Content - delete previousProduction generation with acknowledgeRollbackCandidateRemoval": {
			params: DeleteGenerationRequest{
				LineageID:                           500022,
				GenerationID:                        3029,
				AcknowledgeRollbackCandidateRemoval: true,
			},
			expectedPath:   "/ccm/v2/lineages/500022/generations/3029?acknowledgeRollbackCandidateRemoval=true",
			responseStatus: http.StatusNoContent,
		},
		"404 lineage not found": {
			params:         DeleteGenerationRequest{LineageID: 999999, GenerationID: 3250},
			expectedPath:   "/ccm/v2/lineages/999999/generations/3250",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891060",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteGeneration, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891060",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrDeleteGeneration)
			},
		},
		"404 generation not found": {
			params:         DeleteGenerationRequest{LineageID: 500022, GenerationID: 3031},
			expectedPath:   "/ccm/v2/lineages/500022/generations/3031",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/generation-not-found",
				"title": "Generation not found.",
				"status": 404,
				"detail": "Generation {3031} in Lineage {500022} not found.",
				"instance": "/error-types/generation-not-found?traceId=1234567891057",
				"context": {"generationId": 3031, "lineageId": 500022}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteGeneration, &Error{
					Type:     "/error-types/generation-not-found",
					Title:    "Generation not found.",
					Status:   http.StatusNotFound,
					Detail:   "Generation {3031} in Lineage {500022} not found.",
					Instance: "/error-types/generation-not-found?traceId=1234567891057",
					Context:  map[string]any{"generationId": float64(3031), "lineageId": float64(500022)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrGenerationNotFound)
				assert.ErrorIs(t, err, ErrDeleteGeneration)
			},
		},
		"409 generation delete conflict - removes rollback candidate": {
			params:         DeleteGenerationRequest{LineageID: 500022, GenerationID: 3029},
			expectedPath:   "/ccm/v2/lineages/500022/generations/3029",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/generation-delete-conflict",
				"title": "Generation delete conflict.",
				"status": 409,
				"detail": "Cannot delete generation {3029}: {Deleting this generation will remove the rollback candidate (previousProduction). If you wish to proceed, retry with acknowledgeRollbackCandidateRemoval=true}",
				"instance": "/error-types/generation-delete-conflict?traceId=1234567891062",
				"context": {
					"generationId": 3029,
					"reason": "Deleting this generation will remove the rollback candidate (previousProduction). If you wish to proceed, retry with acknowledgeRollbackCandidateRemoval=true"
				}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteGeneration, &Error{
					Type:     "/error-types/generation-delete-conflict",
					Title:    "Generation delete conflict.",
					Status:   http.StatusConflict,
					Detail:   "Cannot delete generation {3029}: {Deleting this generation will remove the rollback candidate (previousProduction). If you wish to proceed, retry with acknowledgeRollbackCandidateRemoval=true}",
					Instance: "/error-types/generation-delete-conflict?traceId=1234567891062",
					Context: map[string]any{
						"generationId": float64(3029),
						"reason":       "Deleting this generation will remove the rollback candidate (previousProduction). If you wish to proceed, retry with acknowledgeRollbackCandidateRemoval=true",
					},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrGenerationDeleteConflict)
				assert.ErrorIs(t, err, ErrDeleteGeneration)
			},
		},
		"409 generation deployed": {
			params:         DeleteGenerationRequest{LineageID: 500017, GenerationID: 3177},
			expectedPath:   "/ccm/v2/lineages/500017/generations/3177",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/generation-deployed",
				"title": "Cannot delete a deployed generation.",
				"status": 409,
				"detail": "Generation {3177} is currently deployed on staging or production.",
				"instance": "/error-types/generation-deployed?traceId=1234567891058",
				"context": {"generationId": 3177}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteGeneration, &Error{
					Type:     "/error-types/generation-deployed",
					Title:    "Cannot delete a deployed generation.",
					Status:   http.StatusConflict,
					Detail:   "Generation {3177} is currently deployed on staging or production.",
					Instance: "/error-types/generation-deployed?traceId=1234567891058",
					Context:  map[string]any{"generationId": float64(3177)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrGenerationDeployed)
				assert.ErrorIs(t, err, ErrDeleteGeneration)
			},
		},
		"409 last generation": {
			params:         DeleteGenerationRequest{LineageID: 500018, GenerationID: 3251},
			expectedPath:   "/ccm/v2/lineages/500018/generations/3251",
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "/error-types/last-generation",
				"title": "Cannot delete the last generation in a lineage.",
				"status": 409,
				"detail": "Generation is the last remaining generation for lineage {500018}. Delete the lineage instead.",
				"instance": "/error-types/last-generation?traceId=1234567891061",
				"context": {"lineageId": 500018}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteGeneration, &Error{
					Type:     "/error-types/last-generation",
					Title:    "Cannot delete the last generation in a lineage.",
					Status:   http.StatusConflict,
					Detail:   "Generation is the last remaining generation for lineage {500018}. Delete the lineage instead.",
					Instance: "/error-types/last-generation?traceId=1234567891061",
					Context:  map[string]any{"lineageId": float64(500018)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLastGeneration)
				assert.ErrorIs(t, err, ErrDeleteGeneration)
			},
		},
		"500 internal server error": {
			params:         DeleteGenerationRequest{LineageID: 500017, GenerationID: 3250},
			expectedPath:   "/ccm/v2/lineages/500017/generations/3250",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891059"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteGeneration, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891059",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrDeleteGeneration)
			},
		},
		"validation error - missing LineageID": {
			params: DeleteGenerationRequest{GenerationID: 3250},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "deleting generation: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrDeleteGeneration)
				assert.ErrorIs(t, err, ErrStructValidation)
			},
		},
		"validation error - missing GenerationID": {
			params: DeleteGenerationRequest{LineageID: 500017},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "deleting generation: struct validation: GenerationID: cannot be blank")
				assert.ErrorIs(t, err, ErrDeleteGeneration)
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
			err := client.DeleteGeneration(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestListArchivedGenerations(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           ListArchivedGenerationsRequest
		responseStatus   int
		responseBody     string
		expectedResponse *ListArchivedGenerationsResponse
		expectedPath     string
		withError        func(*testing.T, error)
	}{
		"200 OK - list archived generations": {
			params:         ListArchivedGenerationsRequest{LineageID: 500022},
			expectedPath:   "/ccm/v2/lineages/500022/history",
			responseStatus: http.StatusOK,
			responseBody: `{
				"items": [
					{
						"algorithms": [
							{
								"algorithmInstanceId": 4341,
								"certificateStatus": "READY_FOR_USE",
								"keyType": "RSA",
								"signedCertificateNotValidAfterDate": "2027-07-23T11:14:45Z",
								"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:12"
							},
							{
								"algorithmInstanceId": 4342,
								"certificateStatus": "CSR_READY",
								"keyType": "ECDSA"
							}
						],
						"firstPromotedToProductionTime": "2026-07-23T11:18:53Z",
						"generationCreatedBy": "terraform-dev",
						"generationCreatedTime": "2026-07-23T11:13:41Z",
						"generationId": 3027,
						"generationModifiedBy": "terraform-dev",
						"generationModifiedTime": "2026-07-23T11:16:01Z",
						"generationStatus": "READY_FOR_USE"
					},
					{
						"algorithms": [
							{
								"algorithmInstanceId": 3677,
								"certificateStatus": "READY_FOR_USE",
								"keyType": "RSA",
								"signedCertificateNotValidAfterDate": "2027-07-22T12:29:36Z",
								"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:13"
							},
							{
								"algorithmInstanceId": 3678,
								"certificateStatus": "CSR_READY",
								"keyType": "ECDSA"
							}
						],
						"firstPromotedToProductionTime": "2026-07-23T11:08:24Z",
						"generationCreatedBy": "terraform-dev",
						"generationCreatedTime": "2026-07-22T09:21:09Z",
						"generationId": 2448,
						"generationModifiedBy": "terraform-dev",
						"generationModifiedTime": "2026-07-22T12:31:00Z",
						"generationStatus": "READY_FOR_USE"
					},
					{
						"algorithms": [
							{
								"algorithmInstanceId": 4343,
								"certificateStatus": "READY_FOR_USE",
								"keyType": "RSA",
								"signedCertificateNotValidAfterDate": "2027-07-23T11:21:00Z",
								"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:14"
							},
							{
								"algorithmInstanceId": 4344,
								"certificateStatus": "CSR_READY",
								"keyType": "ECDSA"
							}
						],
						"firstPromotedToProductionTime": "2026-07-23T11:25:05Z",
						"generationCreatedBy": "terraform-dev",
						"generationCreatedTime": "2026-07-23T11:20:11Z",
						"generationId": 3028,
						"generationModifiedBy": "terraform-dev",
						"generationModifiedTime": "2026-07-23T11:21:42Z",
						"generationStatus": "READY_FOR_USE"
					}
				]
			}`,
			expectedResponse: &ListArchivedGenerationsResponse{
				Items: []ArchivedGeneration{
					{
						Generation: Generation{
							Algorithms: []Algorithm{
								{
									AlgorithmInstanceID:                4341,
									CertificateStatus:                  CertificateStatusReadyForUse,
									KeyType:                            string(CryptographicAlgorithmRSA),
									SignedCertificateNotValidAfterDate: ptr.To(test.NewTimeFromString(t, "2027-07-23T11:14:45Z")),
									SignedCertificateSerialNumber:      ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:12"),
								},
								{
									AlgorithmInstanceID: 4342,
									CertificateStatus:   CertificateStatusCSRReady,
									KeyType:             string(CryptographicAlgorithmECDSA),
								},
							},
							FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:18:53Z")),
							GenerationCreatedBy:           ptr.To("terraform-dev"),
							GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-23T11:13:41Z")),
							GenerationModifiedBy:          ptr.To("terraform-dev"),
							GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:16:01Z")),
						},
						GenerationID:     3027,
						GenerationStatus: string(GenerationStatusReadyForUse),
					},
					{
						Generation: Generation{
							Algorithms: []Algorithm{
								{
									AlgorithmInstanceID:                3677,
									CertificateStatus:                  CertificateStatusReadyForUse,
									KeyType:                            string(CryptographicAlgorithmRSA),
									SignedCertificateNotValidAfterDate: ptr.To(test.NewTimeFromString(t, "2027-07-22T12:29:36Z")),
									SignedCertificateSerialNumber:      ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:13"),
								},
								{
									AlgorithmInstanceID: 3678,
									CertificateStatus:   CertificateStatusCSRReady,
									KeyType:             string(CryptographicAlgorithmECDSA),
								},
							},
							FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:08:24Z")),
							GenerationCreatedBy:           ptr.To("terraform-dev"),
							GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-22T09:21:09Z")),
							GenerationModifiedBy:          ptr.To("terraform-dev"),
							GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-22T12:31:00Z")),
						},
						GenerationID:     2448,
						GenerationStatus: string(GenerationStatusReadyForUse),
					},
					{
						Generation: Generation{
							Algorithms: []Algorithm{
								{
									AlgorithmInstanceID:                4343,
									CertificateStatus:                  CertificateStatusReadyForUse,
									KeyType:                            string(CryptographicAlgorithmRSA),
									SignedCertificateNotValidAfterDate: ptr.To(test.NewTimeFromString(t, "2027-07-23T11:21:00Z")),
									SignedCertificateSerialNumber:      ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:14"),
								},
								{
									AlgorithmInstanceID: 4344,
									CertificateStatus:   CertificateStatusCSRReady,
									KeyType:             string(CryptographicAlgorithmECDSA),
								},
							},
							FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:25:05Z")),
							GenerationCreatedBy:           ptr.To("terraform-dev"),
							GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-23T11:20:11Z")),
							GenerationModifiedBy:          ptr.To("terraform-dev"),
							GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-23T11:21:42Z")),
						},
						GenerationID:     3028,
						GenerationStatus: string(GenerationStatusReadyForUse),
					},
				},
			},
		},
		"200 OK - list archived generations, includeAlgorithms=true, full algorithm detail": {
			params:         ListArchivedGenerationsRequest{LineageID: 500022, IncludeAlgorithms: true},
			expectedPath:   "/ccm/v2/lineages/500022/history?includeAlgorithms=true",
			responseStatus: http.StatusOK,
			responseBody: `{
				"items": [
					{
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
								"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:12",
								"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:12"
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
						"generationId": 3027,
						"generationModifiedBy": "terraform-dev",
						"generationModifiedTime": "2026-07-23T11:16:01Z",
						"generationStatus": "READY_FOR_USE"
					},
					{
						"algorithms": [
							{
								"algorithmInstanceCreatedBy": "terraform-dev",
								"algorithmInstanceCreatedTime": "2026-07-22T09:21:09Z",
								"algorithmInstanceId": 3677,
								"algorithmInstanceModifiedBy": "terraform-dev",
								"algorithmInstanceModifiedTime": "2026-07-22T12:31:00Z",
								"certificateStatus": "READY_FOR_USE",
								"csrExpirationDate": "2027-09-23T09:21:08Z",
								"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
								"keyType": "RSA",
								"signedCertificateIssuer": "CN=Test Certificate Authority",
								"signedCertificateNotValidAfterDate": "2027-07-22T12:29:36Z",
								"signedCertificateNotValidBeforeDate": "2026-07-22T12:29:36Z",
								"signedCertificatePem": "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n",
								"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:13",
								"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:13"
							},
							{
								"algorithmInstanceCreatedBy": "terraform-dev",
								"algorithmInstanceCreatedTime": "2026-07-22T09:21:09Z",
								"algorithmInstanceId": 3678,
								"certificateStatus": "CSR_READY",
								"csrExpirationDate": "2027-09-23T09:21:08Z",
								"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
								"keyType": "ECDSA"
							}
						],
						"firstPromotedToProductionTime": "2026-07-23T11:08:24Z",
						"generationCreatedBy": "terraform-dev",
						"generationCreatedTime": "2026-07-22T09:21:09Z",
						"generationId": 2448,
						"generationModifiedBy": "terraform-dev",
						"generationModifiedTime": "2026-07-22T12:31:00Z",
						"generationStatus": "READY_FOR_USE"
					},
					{
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
								"signedCertificateSerialNumber": "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:14",
								"signedCertificateSha256Fingerprint": "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:14"
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
						"generationId": 3028,
						"generationModifiedBy": "terraform-dev",
						"generationModifiedTime": "2026-07-23T11:21:42Z",
						"generationStatus": "READY_FOR_USE"
					}
				]
			}`,
			expectedResponse: &ListArchivedGenerationsResponse{
				Items: []ArchivedGeneration{
					{
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
									SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:12"),
									SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:12"),
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
						GenerationID:     3027,
						GenerationStatus: string(GenerationStatusReadyForUse),
					},
					{
						Generation: Generation{
							Algorithms: []Algorithm{
								{
									AlgorithmInstanceCreatedBy:          "terraform-dev",
									AlgorithmInstanceCreatedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-22T09:21:09Z")),
									AlgorithmInstanceID:                 3677,
									AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
									AlgorithmInstanceModifiedTime:       ptr.To(test.NewTimeFromString(t, "2026-07-22T12:31:00Z")),
									CertificateStatus:                   CertificateStatusReadyForUse,
									CSRExpirationDate:                   ptr.To(test.NewTimeFromString(t, "2027-09-23T09:21:08Z")),
									CSRPEM:                              rsaCSRPEM,
									KeyType:                             string(CryptographicAlgorithmRSA),
									SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
									SignedCertificateNotValidAfterDate:  ptr.To(test.NewTimeFromString(t, "2027-07-22T12:29:36Z")),
									SignedCertificateNotValidBeforeDate: ptr.To(test.NewTimeFromString(t, "2026-07-22T12:29:36Z")),
									SignedCertificatePEM:                ptr.To(rsaCertPEM),
									SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:13"),
									SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:13"),
								},
								{
									AlgorithmInstanceCreatedBy:   "terraform-dev",
									AlgorithmInstanceCreatedTime: ptr.To(test.NewTimeFromString(t, "2026-07-22T09:21:09Z")),
									AlgorithmInstanceID:          3678,
									CertificateStatus:            CertificateStatusCSRReady,
									CSRExpirationDate:            ptr.To(test.NewTimeFromString(t, "2027-09-23T09:21:08Z")),
									CSRPEM:                       ecdsaCSRPEM,
									KeyType:                      string(CryptographicAlgorithmECDSA),
								},
							},
							FirstPromotedToProductionTime: ptr.To(test.NewTimeFromString(t, "2026-07-23T11:08:24Z")),
							GenerationCreatedBy:           ptr.To("terraform-dev"),
							GenerationCreatedTime:         ptr.To(test.NewTimeFromString(t, "2026-07-22T09:21:09Z")),
							GenerationModifiedBy:          ptr.To("terraform-dev"),
							GenerationModifiedTime:        ptr.To(test.NewTimeFromString(t, "2026-07-22T12:31:00Z")),
						},
						GenerationID:     2448,
						GenerationStatus: string(GenerationStatusReadyForUse),
					},
					{
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
									SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:14"),
									SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:14"),
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
						GenerationID:     3028,
						GenerationStatus: string(GenerationStatusReadyForUse),
					},
				},
			},
		},
		"200 OK - list archived generations, empty list": {
			params:         ListArchivedGenerationsRequest{LineageID: 500017},
			expectedPath:   "/ccm/v2/lineages/500017/history",
			responseStatus: http.StatusOK,
			responseBody:   `{"items": []}`,
			expectedResponse: &ListArchivedGenerationsResponse{
				Items: []ArchivedGeneration{},
			},
		},
		"404 lineage not found": {
			params:         ListArchivedGenerationsRequest{LineageID: 999999},
			expectedPath:   "/ccm/v2/lineages/999999/history",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/error-types/lineage-not-found",
				"title": "Certificate lineage not found.",
				"status": 404,
				"detail": "Certificate lineage {999999} not found.",
				"instance": "/error-types/lineage-not-found?traceId=1234567891054",
				"context": {"lineageId": 999999}
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListArchivedGenerations, &Error{
					Type:     "/error-types/lineage-not-found",
					Title:    "Certificate lineage not found.",
					Status:   http.StatusNotFound,
					Detail:   "Certificate lineage {999999} not found.",
					Instance: "/error-types/lineage-not-found?traceId=1234567891054",
					Context:  map[string]any{"lineageId": float64(999999)},
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrLineageNotFound)
				assert.ErrorIs(t, err, ErrListArchivedGenerations)
			},
		},
		"500 internal server error": {
			params:         ListArchivedGenerationsRequest{LineageID: 500022},
			expectedPath:   "/ccm/v2/lineages/500022/history",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "/error-types/internal-error",
				"title": "An unexpected error occurred.",
				"status": 500,
				"instance": "/error-types/internal-error?traceId=1234567891053"
			}`,
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrListArchivedGenerations, &Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891053",
				})
				assert.EqualError(t, err, want.Error(), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrInternalError)
				assert.ErrorIs(t, err, ErrListArchivedGenerations)
			},
		},
		"validation error - missing LineageID": {
			params: ListArchivedGenerationsRequest{},
			withError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "listing archived generations: struct validation: LineageID: cannot be blank")
				assert.ErrorIs(t, err, ErrListArchivedGenerations)
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
			result, err := client.ListArchivedGenerations(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestSignedCertificate_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cert      SignedCertificate
		withError bool
		errorMsg  string
	}{
		"valid - signed cert only": {
			cert: SignedCertificate{SignedCertificatePEM: rsaCertPEM},
		},
		"valid - signed cert and trust chain": {
			cert: SignedCertificate{SignedCertificatePEM: rsaCertPEM, TrustChainPEM: rsaCertPEM},
		},
		"validation error - empty": {
			cert:      SignedCertificate{},
			withError: true,
			errorMsg:  "SignedCertificatePEM: cannot be blank.",
		},
		"validation error - trust chain without signed certificate": {
			cert:      SignedCertificate{TrustChainPEM: rsaCertPEM},
			withError: true,
			errorMsg:  "SignedCertificatePEM: cannot be blank.",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.cert.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, tc.errorMsg)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestTargetNetwork_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		targetNetwork TargetNetwork
		withError     bool
	}{
		"empty value is valid": {
			targetNetwork: "",
			withError:     false,
		},
		"STAGING is valid": {
			targetNetwork: TargetNetworkStaging,
			withError:     false,
		},
		"PRODUCTION is valid": {
			targetNetwork: TargetNetworkProduction,
			withError:     false,
		},
		"invalid value": {
			targetNetwork: "NOT_A_REAL_NETWORK",
			withError:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.targetNetwork.Validate()
			if tc.withError {
				require.Error(t, err)
				assert.EqualError(t, err, "value 'NOT_A_REAL_NETWORK' is invalid. Must be either 'STAGING' or 'PRODUCTION'")
				return
			}
			assert.NoError(t, err)
		})
	}
}
