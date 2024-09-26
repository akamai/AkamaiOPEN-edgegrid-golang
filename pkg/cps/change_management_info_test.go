package cps

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChangeManagementInfo(t *testing.T) {
	tests := map[string]struct {
		params           GetChangeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ChangeManagementInfoResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
  "acknowledgementDeadline": null,
  "validationResultHash": "da39a3ee5e6b4b0d3255bfef95601890afd80709",
  "pendingState": {
    "pendingNetworkConfiguration": {
      "dnsNameSettings": null,
      "mustHaveCiphers": "ak-akamai-default2016q3",
      "networkType": null,
      "ocspStapling": "not-set",
      "preferredCiphers": "ak-akamai-default",
      "quicEnabled": "false",
      "sniOnly": "false",
      "disallowedTlsVersions": [
        "TLSv1_2"
      ]
    },
    "pendingCertificates": [
      {
        "certificateType": "third-party",
        "fullCertificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
        "keyAlgorithm": "RSA",
        "ocspStapled": "false",
        "ocspUris": null,
        "signatureAlgorithm": "SHA-256"
      }
    ]
  },
  "validationResult": {
    "errors": null,
    "warnings": [
      {
        "message": "[SAN name [san9.example.com] removed from certificate is still live on the network., SAN name [san8.example.com] removed from certificate is still live on the network.]",
        "messageCode": "no-code"
      }
    ]
  }
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/change-management-info",
			expectedResponse: &ChangeManagementInfoResponse{
				AcknowledgementDeadline: nil,
				ValidationResultHash:    "da39a3ee5e6b4b0d3255bfef95601890afd80709",
				PendingState: PendingState{
					PendingCertificates: []PendingCertificate{
						{
							CertificateType:    "third-party",
							FullCertificate:    "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							OCSPStapled:        "false",
							OCSPURIs:           nil,
							SignatureAlgorithm: "SHA-256",
							KeyAlgorithm:       "RSA",
						},
					},
					PendingNetworkConfiguration: PendingNetworkConfiguration{
						DNSNameSettings:  nil,
						MustHaveCiphers:  "ak-akamai-default2016q3",
						NetworkType:      "",
						OCSPStapling:     "not-set",
						PreferredCiphers: "ak-akamai-default",
						QUICEnabled:      "false",
						SNIOnly:          "false",
						DisallowedTLSVersions: []string{
							"TLSv1_2",
						},
					},
				},
				ValidationResult: &ValidationResult{
					Errors: nil,
					Warnings: []ValidationMessage{
						{
							Message:     "[SAN name [san9.example.com] removed from certificate is still live on the network., SAN name [san8.example.com] removed from certificate is still live on the network.]",
							MessageCode: "no-code",
						},
					},
				},
			},
		},
		"500 internal server error": {
			params: GetChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error making request",
  "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/change-management-info",
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
			params: GetChangeRequest{},
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
				assert.Equal(t, "application/vnd.akamai.cps.change-management-info.v5+json", r.Header.Get("Accept"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetChangeManagementInfo(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetChangeDeploymentInfo(t *testing.T) {
	tests := map[string]struct {
		params           GetChangeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ChangeDeploymentInfoResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "networkConfiguration": {
        "geography": "core",
        "secureNetwork": "enhanced-tls",
        "mustHaveCiphers": "ak-akamai-2020q1",
        "preferredCiphers": "ak-akamai-2020q1",
        "disallowedTlsVersions": [
            "TLSv1_1",
            "TLSv1"
        ],
        "ocspStapling": "not-set",
        "sniOnly": false,
        "quicEnabled": false,
        "dnsNames": null
    },
    "primaryCertificate": {
        "signatureAlgorithm": "SHA-1",
        "certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
        "trustChain": "",
        "expiry": "2023-08-25T13:02:15Z",
        "keyAlgorithm": "RSA"
    },
    "multiStackedCertificates": [],
    "ocspUris": [],
    "ocspStapled": false
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/change-management-info",
			expectedResponse: &ChangeDeploymentInfoResponse{
				NetworkConfiguration: DeploymentNetworkConfiguration{
					Geography:        "core",
					SecureNetwork:    "enhanced-tls",
					MustHaveCiphers:  "ak-akamai-2020q1",
					PreferredCiphers: "ak-akamai-2020q1",
					DisallowedTLSVersions: []string{
						"TLSv1_1",
						"TLSv1",
					},
					OCSPStapling: "not-set",
					SNIOnly:      false,
					QUICEnabled:  false,
					DNSNames:     nil,
				},
				PrimaryCertificate: DeploymentCertificate{
					SignatureAlgorithm: "SHA-1",
					Certificate:        "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					TrustChain:         "",
					Expiry:             "2023-08-25T13:02:15Z",
					KeyAlgorithm:       "RSA",
				},
				MultiStackedCertificates: []DeploymentCertificate{},
				OCSPURIs:                 []string{},
				OCSPStapled:              ptr.To(false),
			},
		},
		"500 internal server error": {
			params: GetChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error making request",
  "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/change-management-info",
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
				assert.Equal(t, "application/vnd.akamai.cps.deployment.v8+json", r.Header.Get("Accept"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetChangeDeploymentInfo(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAcknowledgeChangeManagement(t *testing.T) {
	tests := map[string]struct {
		params         AcknowledgementRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			params: AcknowledgementRequest{
				EnrollmentID: 1,
				ChangeID:     2,
				Acknowledgement: Acknowledgement{
					Acknowledgement: AcknowledgementAcknowledge,
				},
			},
			responseStatus: http.StatusOK,
			responseBody:   "",
			expectedPath:   "/cps/v2/enrollments/1/changes/2/input/update/change-management-ack",
		},
		"500 internal server error": {
			params: AcknowledgementRequest{
				EnrollmentID: 1,
				ChangeID:     2,
				Acknowledgement: Acknowledgement{
					Acknowledgement: AcknowledgementAcknowledge,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error making request",
  "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/update/change-management-ack",
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
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.change-id.v1+json", r.Header.Get("Accept"))
				assert.Equal(t, "application/vnd.akamai.cps.acknowledgement.v1+json; charset=utf-8", r.Header.Get("Content-Type"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.AcknowledgeChangeManagement(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
