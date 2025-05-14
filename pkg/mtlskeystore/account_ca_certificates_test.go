package mtlskeystore

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMTLS_Keystore_ListAccountCACertificates(t *testing.T) {
	tests := map[string]struct {
		request          ListAccountCACertificatesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListAccountCACertificatesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - multiple entries and no query param": {
			request:        ListAccountCACertificatesRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "certificates": [
        {
            "id": 1,
            "version": 1,
            "accountId": "test_account",
            "subject": "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
            "commonName": "Test Common Name",
            "keyAlgorithm": "RSA",
            "keySizeInBytes": 4096,
            "signatureAlgorithm": "SHA256_WITH_RSA",
            "certificate": "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
            "status": "CURRENT",
            "issuedDate": "2025-03-24T15:46:06Z",
            "expiryDate": "2028-03-24T15:46:06Z",
            "createdBy": "Test User",
            "createdDate": "2025-03-24T15:43:50Z"
        },
		{
            "id": 2,
            "version": 2,
            "accountId": "test_account",
            "subject": "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
            "commonName": "Test Common Name",
            "keyAlgorithm": "RSA",
            "keySizeInBytes": 4096,
            "signatureAlgorithm": "SHA256_WITH_RSA",
            "certificate": "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
            "status": "EXPIRED",
            "issuedDate": "2025-03-24T15:46:06Z",
            "expiryDate": "2028-03-24T15:46:06Z",
            "createdBy": "Test User",
            "createdDate": "2025-03-24T15:43:50Z"
        },
		{
            "id": 3,
            "version": 3,
            "accountId": "test_account",
            "subject": "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
            "commonName": "Test Common Name",
            "keyAlgorithm": "RSA",
            "keySizeInBytes": 4096,
			"qualificationDate": "2025-03-25T15:46:06Z",
            "signatureAlgorithm": "SHA256_WITH_RSA",
            "certificate": "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
            "status": "CURRENT",
            "issuedDate": "2025-03-24T15:46:06Z",
            "expiryDate": "2028-03-24T15:46:06Z",
            "createdBy": "Test User",
            "createdDate": "2025-03-24T15:43:50Z"
        }
    ]
}`,
			expectedPath: "/mtls-origin-keystore/v1/ca-certificates",
			expectedResponse: &ListAccountCACertificatesResponse{
				Certificates: []AccountCACertificate{
					{
						ID:                 1,
						Version:            1,
						AccountID:          "test_account",
						Subject:            "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
						CommonName:         "Test Common Name",
						KeyAlgorithm:       "RSA",
						KeySizeInBytes:     4096,
						SignatureAlgorithm: "SHA256_WITH_RSA",
						Certificate:        "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
						Status:             "CURRENT",
						IssuedDate:         test.NewTimeFromString(t, "2025-03-24T15:46:06Z"),
						ExpiryDate:         test.NewTimeFromString(t, "2028-03-24T15:46:06Z"),
						CreatedBy:          "Test User",
						CreatedDate:        test.NewTimeFromString(t, "2025-03-24T15:43:50Z"),
					},
					{
						ID:                 2,
						Version:            2,
						AccountID:          "test_account",
						Subject:            "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
						CommonName:         "Test Common Name",
						KeyAlgorithm:       "RSA",
						KeySizeInBytes:     4096,
						SignatureAlgorithm: "SHA256_WITH_RSA",
						Certificate:        "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
						Status:             "EXPIRED",
						IssuedDate:         test.NewTimeFromString(t, "2025-03-24T15:46:06Z"),
						ExpiryDate:         test.NewTimeFromString(t, "2028-03-24T15:46:06Z"),
						CreatedBy:          "Test User",
						CreatedDate:        test.NewTimeFromString(t, "2025-03-24T15:43:50Z"),
					},
					{
						ID:                 3,
						Version:            3,
						AccountID:          "test_account",
						Subject:            "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
						CommonName:         "Test Common Name",
						KeyAlgorithm:       "RSA",
						KeySizeInBytes:     4096,
						QualificationDate:  ptr.To(test.NewTimeFromString(t, "2025-03-25T15:46:06Z")),
						SignatureAlgorithm: "SHA256_WITH_RSA",
						Certificate:        "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
						Status:             "CURRENT",
						IssuedDate:         test.NewTimeFromString(t, "2025-03-24T15:46:06Z"),
						ExpiryDate:         test.NewTimeFromString(t, "2028-03-24T15:46:06Z"),
						CreatedBy:          "Test User",
						CreatedDate:        test.NewTimeFromString(t, "2025-03-24T15:43:50Z"),
					},
				},
			},
		},
		"200 OK - single entry with query param": {
			request: ListAccountCACertificatesRequest{
				Status: []CertificateStatus{CertificateStatusCurrent},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "certificates": [
        {
            "id": 1,
            "version": 1,
            "accountId": "test_account",
            "subject": "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
            "commonName": "Test Common Name",
            "keyAlgorithm": "RSA",
            "keySizeInBytes": 4096,
            "signatureAlgorithm": "SHA256_WITH_RSA",
            "certificate": "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
            "status": "CURRENT",
            "issuedDate": "2025-03-24T15:46:06Z",
            "expiryDate": "2028-03-24T15:46:06Z",
            "createdBy": "Test User",
            "createdDate": "2025-03-24T15:43:50Z"
        },
		{
            "id": 3,
            "version": 3,
            "accountId": "test_account",
            "subject": "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
            "commonName": "Test Common Name",
            "keyAlgorithm": "RSA",
            "keySizeInBytes": 4096,
			"qualificationDate": "2025-03-25T15:46:06Z",
            "signatureAlgorithm": "SHA256_WITH_RSA",
            "certificate": "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
            "status": "CURRENT",
            "issuedDate": "2025-03-24T15:46:06Z",
            "expiryDate": "2028-03-24T15:46:06Z",
            "createdBy": "Test User",
            "createdDate": "2025-03-24T15:43:50Z"
        }
    ]
}`,
			expectedPath: "/mtls-origin-keystore/v1/ca-certificates?status=CURRENT",
			expectedResponse: &ListAccountCACertificatesResponse{
				Certificates: []AccountCACertificate{
					{
						ID:                 1,
						Version:            1,
						AccountID:          "test_account",
						Subject:            "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
						CommonName:         "Test Common Name",
						KeyAlgorithm:       "RSA",
						KeySizeInBytes:     4096,
						SignatureAlgorithm: "SHA256_WITH_RSA",
						Certificate:        "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
						Status:             "CURRENT",
						IssuedDate:         test.NewTimeFromString(t, "2025-03-24T15:46:06Z"),
						ExpiryDate:         test.NewTimeFromString(t, "2028-03-24T15:46:06Z"),
						CreatedBy:          "Test User",
						CreatedDate:        test.NewTimeFromString(t, "2025-03-24T15:43:50Z"),
					},
					{
						ID:                 3,
						Version:            3,
						AccountID:          "test_account",
						Subject:            "/C=Test Country/O=Test Organization, Inc./OU=Test Organization Unit/CN=Test Common Name/",
						CommonName:         "Test Common Name",
						KeyAlgorithm:       "RSA",
						KeySizeInBytes:     4096,
						QualificationDate:  ptr.To(test.NewTimeFromString(t, "2025-03-25T15:46:06Z")),
						SignatureAlgorithm: "SHA256_WITH_RSA",
						Certificate:        "-----BEGIN CERTIFICATE-----\nCERTIFICATE CONTENT\n-----END CERTIFICATE-----\n",
						Status:             "CURRENT",
						IssuedDate:         test.NewTimeFromString(t, "2025-03-24T15:46:06Z"),
						ExpiryDate:         test.NewTimeFromString(t, "2028-03-24T15:46:06Z"),
						CreatedBy:          "Test User",
						CreatedDate:        test.NewTimeFromString(t, "2025-03-24T15:43:50Z"),
					},
				},
			},
		},
		"200 OK - no entries with combined query params": {
			request: ListAccountCACertificatesRequest{
				Status: []CertificateStatus{CertificateStatusExpired, CertificateStatusCurrent, CertificateStatusPrevious},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "certificates": []
}`,
			expectedPath: "/mtls-origin-keystore/v1/ca-certificates?status=EXPIRED%2CCURRENT%2CPREVIOUS",
			expectedResponse: &ListAccountCACertificatesResponse{
				Certificates: []AccountCACertificate{},
			},
		},
		"validation error - wrong status": {
			request: ListAccountCACertificatesRequest{
				Status: []CertificateStatus{"SOME_WRONG_STATUS", CertificateStatusCurrent, CertificateStatusPrevious},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list account ca certificates: struct validation: Status: list '[SOME_WRONG_STATUS CURRENT PREVIOUS]' contains invalid element 'SOME_WRONG_STATUS'. Each element must be one of: 'CURRENT', 'EXPIRED', 'PREVIOUS', or 'QUALIFYING'", err.Error())
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
   "type": "internal-server-error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "instance": "TestInstances",
   "status": 500
}`,
			expectedPath: "/mtls-origin-keystore/v1/ca-certificates",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error making request",
					Instance: "TestInstances",
					Status:   http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.ListAccountCACertificates(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestStatusesToQueryString(t *testing.T) {
	tests := map[string]struct {
		statuses []CertificateStatus
		expected string
	}{
		"Single status": {
			statuses: []CertificateStatus{CertificateStatusCurrent},
			expected: "status=CURRENT",
		},
		"Multiple statuses": {
			statuses: []CertificateStatus{CertificateStatusCurrent, CertificateStatusExpired},
			expected: "status=CURRENT%2CEXPIRED",
		},
		"All statuses": {
			statuses: []CertificateStatus{CertificateStatusCurrent, CertificateStatusExpired, CertificateStatusPrevious, CertificateStatusQualifying},
			expected: "status=CURRENT%2CEXPIRED%2CPREVIOUS%2CQUALIFYING",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := statusesToQueryString(tc.statuses)
			assert.Equal(t, tc.expected, result)
		})
	}
}
