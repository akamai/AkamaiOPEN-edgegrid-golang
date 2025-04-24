package mtlskeystore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/internal/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMTLS_Keystore_CreateClientCertificate(t *testing.T) {
	tests := map[string]struct {
		request             CreateClientCertificateRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *CreateClientCertificateResponse
		withError           func(*testing.T, error)
	}{
		"201 Created": {
			request: CreateClientCertificateRequest{
				CertificateName: "test-certificate1",
				ContractID:      "test-contract",
				Geography:       GeographyCore,
				GroupID:         12345,
				KeyAlgorithm:    ptr.To(KeyAlgorithmRSA),
				NotificationEmails: []string{
					"jsmith@akamai.com",
					"jkowalski@akamai.com",
				},
				PreferredCA:   ptr.To("AKAMAI"),
				SecureNetwork: SecureNetworkStandardTLS,
				Signer:        SignerAkamai,
				Subject:       ptr.To("CN=test-certificate1"),
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"certificateId": 1234,
	"certificateName": "test-certificate1",
	"createdBy": "jsmith",
	"createdDate": "2023-01-01T00:00:00Z",
	"geography": "CORE",
	"keyAlgorithm": "RSA",
	"notificationEmails": [
		"jsmith@akamai.com",
		"jkowalski@akamai.com"
  	],
	"secureNetwork": "STANDARD_TLS",
	"signer": "AKAMAI",
	"subject": "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate1"
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates",
			expectedRequestBody: `
{
	"certificateName": "test-certificate1",
	"contractId": "test-contract",
	"geography": "CORE",
	"groupId": 12345,
	"keyAlgorithm": "RSA",
	"notificationEmails": [
		"jsmith@akamai.com",
		"jkowalski@akamai.com"
  	],
	"preferredCa": "AKAMAI",
	"secureNetwork": "STANDARD_TLS",
	"signer":"AKAMAI", 
	"subject":"CN=test-certificate1"
}`,
			expectedResponse: &CreateClientCertificateResponse{
				CertificateID:   1234,
				CertificateName: "test-certificate1",
				CreatedBy:       "jsmith",
				CreatedDate:     test.NewTimeFromString(t, "2023-01-01T00:00:00Z"),
				Geography:       GeographyCore,
				KeyAlgorithm:    KeyAlgorithmRSA,
				NotificationEmails: []string{
					"jsmith@akamai.com",
					"jkowalski@akamai.com",
				},
				SecureNetwork: SecureNetworkStandardTLS,
				Signer:        SignerAkamai,
				Subject:       "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate1",
			},
		},
		"201 Created - only with required params": {
			request: CreateClientCertificateRequest{
				CertificateName: "test-certificate1",
				ContractID:      "test-contract",
				Geography:       GeographyChinaAndCore,
				GroupID:         12345,
				NotificationEmails: []string{
					"jsmith@akamai.com",
					"jkowalski@akamai.com",
				},
				SecureNetwork: SecureNetworkEnhancedTLS,
				Signer:        SignerThirdParty,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"certificateId": 1234,
	"certificateName": "test-certificate1",
	"createdBy": "jsmith",
	"createdDate": "2023-01-01T00:00:00Z",
	"geography": "CHINA_AND_CORE",
	"keyAlgorithm": "RSA",
	"notificationEmails": [
		"jsmith@akamai.com",
		"jkowalski@akamai.com"
  	],
	"secureNetwork": "ENHANCED_TLS",
	"signer": "THIRD_PARTY",
	"subject": "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate1"
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates",
			expectedRequestBody: `
{
	"certificateName": "test-certificate1",
	"contractId": "test-contract",
	"geography": "CHINA_AND_CORE",
	"groupId": 12345,
	"notificationEmails": [
		"jsmith@akamai.com",
		"jkowalski@akamai.com"
  	],
	"secureNetwork": "ENHANCED_TLS",
	"signer":"THIRD_PARTY"
}`,
			expectedResponse: &CreateClientCertificateResponse{
				CertificateID:   1234,
				CertificateName: "test-certificate1",
				CreatedBy:       "jsmith",
				CreatedDate:     test.NewTimeFromString(t, "2023-01-01T00:00:00Z"),
				Geography:       GeographyChinaAndCore,
				KeyAlgorithm:    KeyAlgorithmRSA,
				NotificationEmails: []string{
					"jsmith@akamai.com",
					"jkowalski@akamai.com",
				},
				SecureNetwork: SecureNetworkEnhancedTLS,
				Signer:        SignerThirdParty,
				Subject:       "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate1",
			},
		},
		"validation errors": {
			request: CreateClientCertificateRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create client certificate: struct validation: CertificateName: cannot be blank\n"+
					"ContractID: cannot be blank\nGeography: cannot be blank\nGroupID: cannot be blank\nNotificationEmails: cannot be blank\n"+
					"SecureNetwork: cannot be blank\nSigner: cannot be blank", err.Error())
			},
		},
		"validation error - incorrect geography": {
			request: CreateClientCertificateRequest{
				CertificateName: "test-certificate1",
				ContractID:      "test-contract",
				Geography:       "test-geography",
				GroupID:         12345,
				NotificationEmails: []string{
					"jsmith@akamai.com",
				},
				SecureNetwork: SecureNetworkStandardTLS,
				Signer:        SignerAkamai,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create client certificate: struct validation: Geography: value 'test-geography' is invalid. Must be one of: 'CORE', 'RUSSIA_AND_CORE', 'CHINA_AND_CORE'", err.Error())
			},
		},
		"validation error - incorrect secureNetwork": {
			request: CreateClientCertificateRequest{
				CertificateName: "test-certificate1",
				ContractID:      "test-contract",
				Geography:       GeographyRussiaAndCore,
				GroupID:         12345,
				NotificationEmails: []string{
					"jsmith@akamai.com",
				},
				SecureNetwork: "test-secure-network",
				Signer:        SignerThirdParty,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create client certificate: struct validation: SecureNetwork: value 'test-secure-network' is invalid. Must be one of: 'STANDARD_TLS', 'ENHANCED_TLS'", err.Error())
			},
		},
		"validation error - incorrect keyAlgorithm": {
			request: CreateClientCertificateRequest{
				CertificateName: "test-certificate1",
				ContractID:      "test-contract",
				Geography:       GeographyCore,
				GroupID:         12345,
				KeyAlgorithm:    (*CryptographicAlgorithm)(ptr.To("test-key-algorithm")),
				NotificationEmails: []string{
					"jsmith@akamai.com",
				},
				SecureNetwork: SecureNetworkEnhancedTLS,
				Signer:        SignerAkamai,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create client certificate: struct validation: KeyAlgorithm: value 'test-key-algorithm' is invalid. Must be one of: 'RSA', 'ECDSA'", err.Error())
			},
		},
		"validation error - incorrect signer": {
			request: CreateClientCertificateRequest{
				CertificateName: "test-certificate1",
				ContractID:      "test-contract",
				Geography:       GeographyChinaAndCore,
				GroupID:         12345,
				NotificationEmails: []string{
					"jsmith@akamai.com",
				},
				SecureNetwork: SecureNetworkStandardTLS,
				Signer:        "test-signer",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create client certificate: struct validation: Signer: value 'test-signer' is invalid. Must be one of: 'AKAMAI', 'THIRD_PARTY'", err.Error())
			},
		},
		"500 internal server error": {
			request: CreateClientCertificateRequest{
				CertificateName: "test-certificate1",
				ContractID:      "test-contract",
				Geography:       GeographyRussiaAndCore,
				GroupID:         12345,
				NotificationEmails: []string{
					"jsmith@akamai.com",
				},
				SecureNetwork: SecureNetworkEnhancedTLS,
				Signer:        SignerAkamai,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	   "type": "internal-server-error",
	   "title": "Internal Server Error",
	   "detail": "Error making request",
	   "instance": "TestInstances",
	   "status": 500
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

				if tc.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.CreateClientCertificate(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestMTLS_Keystore_PatchClientCertificate(t *testing.T) {
	tests := map[string]struct {
		request             PatchClientCertificateRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			request: PatchClientCertificateRequest{
				CertificateID: 1234,
				Body: PatchClientCertificateRequestBody{
					CertificateName: ptr.To("test-certificate1"),
					NotificationEmails: []string{
						"jsmith@akamai.com",
					},
				},
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/1234",
			expectedRequestBody: `
{
	"certificateName": "test-certificate1",
	"notificationEmails": [
			"jsmith@akamai.com"
	]
}`,
		},
		"200 OK - only certificateId provided": {
			request: PatchClientCertificateRequest{
				CertificateID: 1234,
				Body: PatchClientCertificateRequestBody{
					CertificateName:    ptr.To("test-certificate1"),
					NotificationEmails: nil,
				},
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/1234",
			expectedRequestBody: `
{
	"certificateName": "test-certificate1",
	"notificationEmails": null
}`,
		},
		"200 OK - only notificationEmails provided": {
			request: PatchClientCertificateRequest{
				CertificateID: 1234,
				Body: PatchClientCertificateRequestBody{
					CertificateName: nil,
					NotificationEmails: []string{
						"jsmith@akamai.com",
					},
				},
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/1234",
			expectedRequestBody: `
{
	"certificateName": null,
	"notificationEmails": [
			"jsmith@akamai.com"
	]
}`,
		},
		"validation error - missing certificateId": {
			request: PatchClientCertificateRequest{
				Body: PatchClientCertificateRequestBody{
					CertificateName: ptr.To("test-certificate1"),
					NotificationEmails: []string{
						"jsmith@akamai.com",
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patch client certificate: struct validation: CertificateID: cannot be blank", err.Error())
			},
		},
		"validation error - empty request body": {
			request: PatchClientCertificateRequest{
				CertificateID: 1234,
				Body:          PatchClientCertificateRequestBody{},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patch client certificate: struct validation: Body: CertificateName or NotificationEmails must be provided", err.Error())
			},
		},
		"validation error - nil value notificationEmails and certificateName": {
			request: PatchClientCertificateRequest{
				CertificateID: 1234,
				Body: PatchClientCertificateRequestBody{
					CertificateName:    nil,
					NotificationEmails: nil,
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patch client certificate: struct validation: Body: CertificateName or NotificationEmails must be provided", err.Error())
			},
		},
		"validation error - empty certificateName": {
			request: PatchClientCertificateRequest{
				CertificateID: 1234,
				Body: PatchClientCertificateRequestBody{
					CertificateName:    ptr.To(""),
					NotificationEmails: nil,
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patch client certificate: struct validation: Body: {\n\tCertificateName: value is invalid\n}", err.Error())
			},
		},
		"validation error - certificateName longer than 64": {
			request: PatchClientCertificateRequest{
				CertificateID: 1234,
				Body: PatchClientCertificateRequestBody{
					CertificateName:    ptr.To("test-certificate1-test-certificate1-test-certificate1-test-certificate1-test-certificate1"),
					NotificationEmails: nil,
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patch client certificate: struct validation: Body: {\n\tCertificateName: value 'test-certificate1-test-certificate1-test-certificate1-test-certificate1-test-certificate1' is invalid. Must be between 1 and 64 characters\n}", err.Error())
			},
		},
		"404 Not Found": {
			request: PatchClientCertificateRequest{
				CertificateID: 1,
				Body: PatchClientCertificateRequestBody{
					CertificateName:    ptr.To("test-certificate1"),
					NotificationEmails: nil,
				},
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
	"type": "not-found",
	"title": "Not Found",
	"detail": "Client certificate not found",
	"instance": "TestInstances",
	"status": 404
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates/1",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "not-found",
					Title:    "Not Found",
					Detail:   "Client certificate not found",
					Instance: "TestInstances",
					Status:   http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			request: PatchClientCertificateRequest{
				CertificateID: 1234,
				Body: PatchClientCertificateRequestBody{
					CertificateName: ptr.To("test-certificate1"),
					NotificationEmails: []string{
						"jsmith@akamai.com",
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	   "type": "internal-server-error",
	   "title": "Internal Server Error",
	   "detail": "Error making request",
	   "instance": "TestInstances",
	   "status": 500
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates/1234",
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
				assert.Equal(t, http.MethodPatch, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

				if tc.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			err := client.PatchClientCertificate(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMTLS_Keystore_GetClientCertificate(t *testing.T) {
	tests := map[string]struct {
		request          GetClientCertificateRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetClientCertificateResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: GetClientCertificateRequest{
				CertificateID: 1234,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"certificateId": 1234,
	"certificateName": "test-certificate1",
	"createdBy": "jsmith",
	"createdDate": "2023-01-01T00:00:00Z",
	"geography": "CORE",
	"keyAlgorithm": "RSA",
	"notificationEmails": [
		"jsmith@akamai.com",
		"jkowalski@akamai.com"
  	],
	"secureNetwork": "STANDARD_TLS",
	"signer": "AKAMAI",
	"subject": "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate1"
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates/1234",
			expectedResponse: &GetClientCertificateResponse{
				CertificateID:   1234,
				CertificateName: "test-certificate1",
				CreatedBy:       "jsmith",
				CreatedDate:     test.NewTimeFromString(t, "2023-01-01T00:00:00Z"),
				Geography:       "CORE",
				KeyAlgorithm:    "RSA",
				NotificationEmails: []string{
					"jsmith@akamai.com",
					"jkowalski@akamai.com",
				},
				SecureNetwork: "STANDARD_TLS",
				Signer:        "AKAMAI",
				Subject:       "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate1",
			},
		},
		"Validation error": {
			request: GetClientCertificateRequest{
				CertificateID: 0,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get client certificate: struct validation: CertificateID: cannot be blank", err.Error())
			},
		},
		"404 Not Found": {
			request: GetClientCertificateRequest{
				CertificateID: 1,
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
	"type": "not-found",
	"title": "Not Found",
	"detail": "Client certificate not found",
	"instance": "TestInstances",
	"status": 404
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates/1",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "not-found",
					Title:    "Not Found",
					Detail:   "Client certificate not found",
					Instance: "TestInstances",
					Status:   http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			request: GetClientCertificateRequest{
				CertificateID: 1234,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	   "type": "internal-server-error",
	   "title": "Internal Server Error",
	   "detail": "Error making request",
	   "instance": "TestInstances",
	   "status": 500
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates/1234",
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
			result, err := client.GetClientCertificate(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestMTLS_Keystore_ListClientCertificates(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListClientCertificatesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
{
	"certificates": [
		{
			"certificateId": 1234,
			"certificateName": "test-certificate1",
			"createdBy": "jsmith",
			"createdDate": "2023-01-01T00:00:00Z",
			"geography": "CORE",
			"keyAlgorithm": "RSA",
			"notificationEmails": [
				"jsmith@akamai.com",
				"jkowalski@akamai.com"
			  ],
			"secureNetwork": "STANDARD_TLS",
			"signer": "AKAMAI",
			"subject": "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate1"
		},
		{
			"certificateId": 12345,
			"certificateName": "test-certificate2",
			"createdBy": "jsmith",
			"createdDate": "2023-01-02T00:00:00Z",
			"geography": "CORE",
			"keyAlgorithm": "RSA",
			"notificationEmails": [
				"jsmith@akamai.com",
				"jkowalski@akamai.com"
			  ],
			"secureNetwork": "STANDARD_TLS",
			"signer": "AKAMAI",
			"subject": "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate2"
		}
	]
}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates",
			expectedResponse: &ListClientCertificatesResponse{
				Certificates: []Certificate{
					{
						CertificateID:   1234,
						CertificateName: "test-certificate1",
						CreatedBy:       "jsmith",
						CreatedDate:     test.NewTimeFromString(t, "2023-01-01T00:00:00Z"),
						Geography:       "CORE",
						KeyAlgorithm:    "RSA",
						NotificationEmails: []string{
							"jsmith@akamai.com",
							"jkowalski@akamai.com",
						},
						SecureNetwork: "STANDARD_TLS",
						Signer:        "AKAMAI",
						Subject:       "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate1",
					},
					{
						CertificateID:   12345,
						CertificateName: "test-certificate2",
						CreatedBy:       "jsmith",
						CreatedDate:     test.NewTimeFromString(t, "2023-01-02T00:00:00Z"),
						Geography:       "CORE",
						KeyAlgorithm:    "RSA",
						NotificationEmails: []string{
							"jsmith@akamai.com",
							"jkowalski@akamai.com",
						},
						SecureNetwork: "STANDARD_TLS",
						Signer:        "AKAMAI",
						Subject:       "/C=US/O=Akamai Technologies, Inc./OU=123 test-contract 12345/CN=test-certificate2",
					},
				},
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
			expectedPath: "/mtls-origin-keystore/v1/client-certificates",
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
			result, err := client.ListClientCertificates(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
