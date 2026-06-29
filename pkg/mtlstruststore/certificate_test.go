package mtlstruststore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCertificates(t *testing.T) {
	tests := map[string]struct {
		params              ValidateCertificatesRequest
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedResponse    *ValidateCertificatesResponse
		withError           func(*testing.T, error)
	}{
		"200 - valid certs provided": {
			params: ValidateCertificatesRequest{
				Certificates: []ValidateCertificate{
					{
						CertificatePEM: "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
					},
					{
						CertificatePEM: "-----BEGIN CERTIFICATE-----\nCERT2\n-----END CERTIFICATE-----",
						Description:    ptr.To("desc"),
					},
				},
			},
			expectedPath: "/mtls-edge-truststore/v2/certificates/validate",
			expectedRequestBody: `{
	"allowInsecureSha1":false,
	"certificates": [
		{
		  "certificatePem": "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----"
		},
		{
		  "certificatePem": "-----BEGIN CERTIFICATE-----\nCERT2\n-----END CERTIFICATE-----",
		  "description": "desc"
		}
	]
}`,
			responseStatus: http.StatusOK,
			responseBody: `{
    "allowInsecureSha1": false,
    "certificates": [
        {
            "certificatePem": "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
            "endDate": "2033-04-22T22:49:13Z",
            "fingerprint": "aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae",
            "issuer": "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
            "serialNumber": "dfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdf",
            "signatureAlgorithm": "SHA256WITHRSA",
            "startDate": "2023-04-25T22:49:13Z",
            "subject": "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US"
        },
        {
            "certificatePem": "-----BEGIN CERTIFICATE-----\nCERT2\n-----END CERTIFICATE-----",
            "description": "desc",
            "endDate": "2033-04-22T23:22:06Z",
            "fingerprint": "bcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbc",
            "issuer": "EMAILADDRESS=test2@example.com, CN=test-example1.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
            "serialNumber": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
            "signatureAlgorithm": "SHA256WITHRSA",
            "startDate": "2023-04-25T23:22:06Z",
            "subject": "EMAILADDRESS=test2@example.com, CN=test-example1.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US"
        }
    ],
    "validation": {
        "warnings": []
    }
}`,
			expectedResponse: &ValidateCertificatesResponse{
				AllowInsecureSHA1: false,
				Certificates: []ValidateCertificateResponse{
					{
						CertificatePEM:     "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
						EndDate:            test.NewTimeFromString(t, "2033-04-22T22:49:13Z"),
						Fingerprint:        "aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae",
						Issuer:             "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
						SerialNumber:       "dfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdf",
						SignatureAlgorithm: "SHA256WITHRSA",
						StartDate:          test.NewTimeFromString(t, "2023-04-25T22:49:13Z"),
						Subject:            "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
					},
					{
						CertificatePEM:     "-----BEGIN CERTIFICATE-----\nCERT2\n-----END CERTIFICATE-----",
						Description:        ptr.To("desc"),
						EndDate:            test.NewTimeFromString(t, "2033-04-22T23:22:06Z"),
						Fingerprint:        "bcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbc",
						Issuer:             "EMAILADDRESS=test2@example.com, CN=test-example1.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
						SerialNumber:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						SignatureAlgorithm: "SHA256WITHRSA",
						StartDate:          test.NewTimeFromString(t, "2023-04-25T23:22:06Z"),
						Subject:            "EMAILADDRESS=test2@example.com, CN=test-example1.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
					},
				},
				Validation: Validation{Warnings: []Warning{}},
			},
		},
		"200 - valid certs but there is a duplication (warning)": {
			params: ValidateCertificatesRequest{
				Certificates: []ValidateCertificate{
					{
						CertificatePEM: "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
					},
					{
						CertificatePEM: "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
					},
				},
			},
			expectedPath: "/mtls-edge-truststore/v2/certificates/validate",
			expectedRequestBody: `{
	"allowInsecureSha1":false,
	"certificates": [
		{
			"certificatePem": "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----"
		},
		{
			"certificatePem": "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----"
		}
	]
}`,
			responseStatus: http.StatusOK,
			responseBody: `{
    "allowInsecureSha1": false,
    "certificates": [
        {
            "certificatePem": "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
            "endDate": "2033-04-22T22:49:13Z",
            "fingerprint": "aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae",
            "issuer": "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
            "serialNumber": "dfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdf",
            "signatureAlgorithm": "SHA256WITHRSA",
            "startDate": "2023-04-25T22:49:13Z",
            "subject": "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US"
        },
        {
            "certificatePem": "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
            "endDate": "2033-04-22T22:49:13Z",
            "fingerprint": "aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae",
            "issuer": "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
            "serialNumber": "dfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdf",
            "signatureAlgorithm": "SHA256WITHRSA",
            "startDate": "2023-04-25T22:49:13Z",
            "subject": "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US"
        }
    ],
    "validation": {
        "warnings": [
            {
                "contextInfo": {
                    "fingerPrint": "aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae"
                },
                "detail": "The certificate with the fingerprint aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae has been submitted more than once. Duplicate certificates are not allowed.",
                "pointer": "/certificates/1",
                "title": "Duplicate certificate has been submitted in the certificates.",
                "type": "/mtls-edge-truststore/v2/error-types/duplicate-certificate"
            }
        ]
    }
}`,
			expectedResponse: &ValidateCertificatesResponse{
				AllowInsecureSHA1: false,
				Certificates: []ValidateCertificateResponse{
					{
						CertificatePEM:     "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
						EndDate:            test.NewTimeFromString(t, "2033-04-22T22:49:13Z"),
						Fingerprint:        "aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae",
						Issuer:             "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
						SerialNumber:       "dfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdf",
						SignatureAlgorithm: "SHA256WITHRSA",
						StartDate:          test.NewTimeFromString(t, "2023-04-25T22:49:13Z"),
						Subject:            "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
					},
					{
						CertificatePEM:     "-----BEGIN CERTIFICATE-----\nCERT1\n-----END CERTIFICATE-----",
						EndDate:            test.NewTimeFromString(t, "2033-04-22T22:49:13Z"),
						Fingerprint:        "aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae",
						Issuer:             "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
						SerialNumber:       "dfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdf",
						SignatureAlgorithm: "SHA256WITHRSA",
						StartDate:          test.NewTimeFromString(t, "2023-04-25T22:49:13Z"),
						Subject:            "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US",
					},
				},
				Validation: Validation{
					Warnings: []Warning{
						{
							ContextInfo: map[string]interface{}{
								"fingerPrint": "aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae",
							},
							Detail:  "The certificate with the fingerprint aeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeaeae has been submitted more than once. Duplicate certificates are not allowed.",
							Pointer: "/certificates/1",
							Title:   "Duplicate certificate has been submitted in the certificates.",
							Type:    "/mtls-edge-truststore/v2/error-types/duplicate-certificate",
						},
					},
				},
			},
		},
		"missing required params - validation error": {
			params: ValidateCertificatesRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "validate certificates failed: struct validation: Certificates: cannot be blank", err.Error())
			},
		},
		"certificate is empty - validation error": {
			params: ValidateCertificatesRequest{Certificates: []ValidateCertificate{{CertificatePEM: ""}}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "validate certificates failed: struct validation: Certificates[0]: {\n\tCertificatePEM: cannot be blank\n}", err.Error())
			},
		},
		"400 invalid certificate": {
			params: ValidateCertificatesRequest{
				Certificates: []ValidateCertificate{
					{
						CertificatePEM: "malformed",
						Description:    ptr.To("desc"),
					},
				},
			},
			expectedPath:   "/mtls-edge-truststore/v2/certificates/validate",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
    "errors": [
        {
            "contextInfo": {
                "certificatePem": "malformed",
                "description": "desc"
            },
            "detail": "Certificate PEM string is missing -----BEGIN CERTIFICATE----- header.",
            "pointer": "/certificates/0"
        }
    ],
    "status": 400,
    "title": "Certificate(s) has failed validation.",
    "type": "/mtls-edge-truststore/error-types/certificate-validation-failure"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCertValidationFailure), "want: %s; got: %s", ErrCertValidationFailure, err)
			},
		},
		"400 expired certificate as duplicate": {
			params: ValidateCertificatesRequest{
				Certificates: []ValidateCertificate{
					{
						CertificatePEM: "malformed",
						Description:    ptr.To("desc"),
					},
				},
			},
			expectedPath:   "/mtls-edge-truststore/v2/certificates/validate",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
    "errors" : [ {
        "contextInfo" : {
            "checkDate" : "2025-03-06T10:00:21Z",
            "expiryDate" : "2023-04-23T17:00:34Z",
            "fingerPrint" : "bcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbc",
            "subject" : "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US"
        },
        "detail" : "The certificate with subject EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US and fingerprint bcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbc has expired. Expiry date is 2023-04-23T17:00:34Z. The check was performed on 2025-03-06T10:00:21Z.",
        "pointer" : "/certificates/0"
    }, {
        "contextInfo" : {
            "checkDate" : "2025-03-06T10:00:21Z",
            "expiryDate" : "2023-04-23T17:00:34Z",
            "fingerPrint" : "bcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbc",
            "subject" : "EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US"
        },
        "detail" : "The certificate with subject EMAILADDRESS=test1@example.com, CN=duplicate-test-cn.com, OU=Media BU, O=Akamai, L=SF, ST=CA, C=US and fingerprint bcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbcbc has expired. Expiry date is 2023-04-23T17:00:34Z. The check was performed on 2025-03-06T10:00:21Z.",
        "pointer" : "/certificates/1"
    } ],
    "status" : 400,
    "title" : "Certificate(s) has failed validation.",
    "type" : "/mtls-edge-truststore/error-types/certificate-validation-failure"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCertValidationFailure), "want: %s; got: %s", ErrCertValidationFailure, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if tc.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ValidateCertificates(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
