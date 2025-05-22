package mtlskeystore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRotateClientCertificateVersion(t *testing.T) {
	tests := map[string]struct {
		request          RotateClientCertificateVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RotateClientCertificateVersionResponse
		withError        func(*testing.T, error)
	}{
		"201- Successful Rotate client certificate version": {
			request: RotateClientCertificateVersionRequest{
				CertificateID: 123,
			},
			responseStatus: http.StatusCreated,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions",
			responseBody: `{
				  "certificateBlock": {
					"certificate": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					"keyAlgorithm": "RSA",
					"trustChain": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"
				  },
				  "createdBy": "jperez",
				  "createdDate": "2024-05-21T04:35:20Z",
				  "expiryDate": "2024-08-21T04:35:21Z",
				  "issuedDate": "2024-05-21T04:35:21Z",
				  "issuer": "1360 Account CA G366",
				  "keyAlgorithm": "RSA",
				  "keySizeInBytes": "2048",
				  "signatureAlgorithm": "SHA256_WITH_RSA",
				  "status": "DEPLOYMENT_PENDING",
				  "subject": "/C=US/O=Akamai Technologies/OU=KMI/CN=/",
				  "validation": {
					"errors": [],
					"warnings": []
				  },
				  "version": 4,
				  "versionGuid": "13d16e57-22fa-4475-af0a-b2b745115128"
			}`,
			expectedResponse: &RotateClientCertificateVersionResponse{
				ClientCertificateVersion: ClientCertificateVersion{
					Version:     4,
					VersionGUID: "13d16e57-22fa-4475-af0a-b2b745115128",
					CertificateBlock: &CertificateBlock{
						Certificate:  "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
						KeyAlgorithm: "RSA",
						TrustChain:   "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					},
					CreatedBy:          "jperez",
					CreatedDate:        "2024-05-21T04:35:20Z",
					ExpiryDate:         "2024-08-21T04:35:21Z",
					IssuedDate:         "2024-05-21T04:35:21Z",
					Issuer:             "1360 Account CA G366",
					KeyAlgorithm:       "RSA",
					KeySizeInBytes:     "2048",
					SignatureAlgorithm: "SHA256_WITH_RSA",
					Status:             "DEPLOYMENT_PENDING",
					Subject:            "/C=US/O=Akamai Technologies/OU=KMI/CN=/",
					Validation: ValidationResult{
						Errors:   []ValidationDetail{},
						Warnings: []ValidationDetail{},
					},
				},
			},
		},
		"Validation error - missing CertificateID": {
			request: RotateClientCertificateVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "rotating client certificate version: validation failed: CertificateID: cannot be blank", err.Error())
			},
		},
		"Error Response 404 - Client Certificate not found": {
			request: RotateClientCertificateVersionRequest{
				CertificateID: 123,
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions",
			responseBody: `
				{
				  "detail": "The requested resource could not be found on the server.",
				  "field": "certificateId",
				  "instance": "32be8a77-cc41-495e-94d8-cb536afea149",
				  "problemId": "32be8a77-cc41-495e-94d8-cb536afea149",
				  "status": 404,
				  "title": "Resource Not Found",
				  "type": "resource-not-found",
				  "value": "180131"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrClientCertificateNotFound))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.RotateClientCertificateVersion(context.Background(), test.request)

			if test.withError != nil {
				test.withError(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetClientCertificateVersions(t *testing.T) {
	tests := map[string]struct {
		request          GetClientCertificateVersionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetClientCertificateVersionsResponse
		withError        func(*testing.T, error)
	}{
		"200 - Successful get client certificate versions": {
			request: GetClientCertificateVersionsRequest{
				CertificateID: 123,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions",
			responseBody: `{
			  "versions": [
				{
				  "certificateBlock": {
					"certificate": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					"keyAlgorithm": "RSA",
					"trustChain": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"
				  },
				  "createdBy": "jperez",
				  "createdDate": "2024-05-21T04:35:20Z",
				  "expiryDate": "2024-08-21T04:35:21Z",
				  "issuedDate": "2024-05-21T04:35:21Z",
				  "issuer": "1360 Account CA G366",
				  "keyAlgorithm": "RSA",
				  "keySizeInBytes": "2048",
				  "signatureAlgorithm": "SHA256_WITH_RSA",
				  "status": "DEPLOYMENT_PENDING",
				  "subject": "/C=US/O=Akamai Technologies/OU=KMI/CN=/",
				  "validation": {
					"errors": [],
					"warnings": []
				  },
				  "version": 4,
				  "versionGuid": "13d16e57-22fa-4475-af0a-b2b745115128"
				}
			  ]
			}`,
			expectedResponse: &GetClientCertificateVersionsResponse{
				Versions: []ClientCertificateVersion{
					{
						Version:     4,
						VersionGUID: "13d16e57-22fa-4475-af0a-b2b745115128",
						CertificateBlock: &CertificateBlock{
							Certificate:  "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
							KeyAlgorithm: "RSA",
							TrustChain:   "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
						},
						CreatedBy:          "jperez",
						CreatedDate:        "2024-05-21T04:35:20Z",
						ExpiryDate:         "2024-08-21T04:35:21Z",
						IssuedDate:         "2024-05-21T04:35:21Z",
						Issuer:             "1360 Account CA G366",
						KeyAlgorithm:       "RSA",
						KeySizeInBytes:     "2048",
						SignatureAlgorithm: "SHA256_WITH_RSA",
						Status:             "DEPLOYMENT_PENDING",
						Subject:            "/C=US/O=Akamai Technologies/OU=KMI/CN=/",
						Validation: ValidationResult{
							Errors:   []ValidationDetail{},
							Warnings: []ValidationDetail{},
						},
					},
				},
			},
		},
		"200 - Successful get client certificate versions with associated properties": {
			request: GetClientCertificateVersionsRequest{
				CertificateID:               123,
				IncludeAssociatedProperties: true,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions?includeAssociatedProperties=true",
			responseBody: `{
			  "versions": [
				{
				  "certificateBlock": {
					"certificate": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					"keyAlgorithm": "RSA",
					"trustChain": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"
				  },
				  "createdBy": "jperez",
				  "createdDate": "2024-05-21T04:35:20Z",
				  "expiryDate": "2024-08-21T04:35:21Z",
				  "issuedDate": "2024-05-21T04:35:21Z",
				  "issuer": "1360 Account CA G366",
				  "keyAlgorithm": "RSA",
				  "keySizeInBytes": "2048",
				  "signatureAlgorithm": "SHA256_WITH_RSA",
				  "status": "DEPLOYMENT_PENDING",
				  "subject": "/C=US/O=Akamai Technologies/OU=KMI/CN=/",
				  "validation": {
					"errors": [],
					"warnings": []
				  },
				  "version": 4,
				  "versionGuid": "13d16e57-22fa-4475-af0a-b2b745115128",
				  "properties": [
					{
					  "assetId": 111111,
					  "groupId": 222222,
					  "propertyName": "propertyName",
					  "propertyVersion": 3
					}
				  ]
				}
			  ]
			}`,
			expectedResponse: &GetClientCertificateVersionsResponse{
				Versions: []ClientCertificateVersion{
					{
						Version:     4,
						VersionGUID: "13d16e57-22fa-4475-af0a-b2b745115128",
						CertificateBlock: &CertificateBlock{
							Certificate:  "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
							KeyAlgorithm: "RSA",
							TrustChain:   "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
						},
						CreatedBy:          "jperez",
						CreatedDate:        "2024-05-21T04:35:20Z",
						ExpiryDate:         "2024-08-21T04:35:21Z",
						IssuedDate:         "2024-05-21T04:35:21Z",
						Issuer:             "1360 Account CA G366",
						KeyAlgorithm:       "RSA",
						KeySizeInBytes:     "2048",
						SignatureAlgorithm: "SHA256_WITH_RSA",
						Status:             "DEPLOYMENT_PENDING",
						Subject:            "/C=US/O=Akamai Technologies/OU=KMI/CN=/",
						Validation: ValidationResult{
							Errors:   []ValidationDetail{},
							Warnings: []ValidationDetail{},
						},
						AssociatedProperties: []AssociatedProperty{
							{
								AssetID:         111111,
								GroupID:         222222,
								PropertyName:    "propertyName",
								PropertyVersion: 3,
							},
						},
					},
				},
			},
		},
		"Validation error - missing CertificateID": {
			request: GetClientCertificateVersionsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "fetching client certificate versions: validation failed: CertificateID: cannot be blank", err.Error())
			},
		},
		"Error Response 404 - Client Certificate not found": {
			request: GetClientCertificateVersionsRequest{
				CertificateID: 123,
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions",
			responseBody: `
				{
				  "detail": "The requested resource could not be found on the server.",
				  "field": "certificateId",
				  "instance": "32be8a77-cc41-495e-94d8-cb536afea149",
				  "problemId": "32be8a77-cc41-495e-94d8-cb536afea149",
				  "status": 404,
				  "title": "Resource Not Found",
				  "type": "resource-not-found",
				  "value": "180131"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrClientCertificateNotFound))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.GetClientCertificateVersions(context.Background(), test.request)

			if test.withError != nil {
				test.withError(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteClientCertificateVersion(t *testing.T) {
	tests := map[string]struct {
		request          DeleteClientCertificateVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DeleteClientCertificateVersionResponse
		withError        func(*testing.T, error)
	}{
		"202- Successful submitted deletion request for client certificate version": {
			request: DeleteClientCertificateVersionRequest{
				CertificateID: 123,
				Version:       1,
			},
			responseStatus: http.StatusAccepted,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions/1",
			responseBody: `{
			  "message": "It's being scheduled to delete on 2024-05-10T00:00:00Z. The delete request will be cancelled automatically if it is used again in any delivery configuration."
			}`,
			expectedResponse: &DeleteClientCertificateVersionResponse{
				Message: "It's being scheduled to delete on 2024-05-10T00:00:00Z. The delete request will be cancelled automatically if it is used again in any delivery configuration.",
			},
		},
		"204- Successful submitted deletion request for client certificate version": {
			request: DeleteClientCertificateVersionRequest{
				CertificateID: 123,
				Version:       1,
			},
			responseStatus:   http.StatusNoContent,
			expectedPath:     "/mtls-origin-keystore/v1/client-certificates/123/versions/1",
			expectedResponse: nil,
		},
		"Validation error - missing CertificateID": {
			request: DeleteClientCertificateVersionRequest{
				Version: 1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "deleting client certificate version: validation failed: CertificateID: cannot be blank", err.Error())
			},
		},
		"Validation error - missing Version": {
			request: DeleteClientCertificateVersionRequest{
				CertificateID: 123,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "deleting client certificate version: validation failed: Version: cannot be blank", err.Error())
			},
		},
		"Error Response 404 - Client Certificate not found": {
			request: DeleteClientCertificateVersionRequest{
				CertificateID: 123,
				Version:       1,
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions/1",
			responseBody: `
				{
				  "detail": "The requested resource could not be found on the server.",
				  "field": "certificateId",
				  "instance": "32be8a77-cc41-495e-94d8-cb536afea149",
				  "problemId": "32be8a77-cc41-495e-94d8-cb536afea149",
				  "status": 404,
				  "title": "Resource Not Found",
				  "type": "resource-not-found",
				  "value": "180131"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrClientCertificateNotFound))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteClientCertificateVersion(context.Background(), test.request)

			if test.withError != nil {
				test.withError(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUploadClientCertificateVersion(t *testing.T) {
	tests := map[string]struct {
		request             UploadSignedClientCertificateRequest
		responseStatus      int
		responseBody        string
		expectedRequestBody string
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"200- Successful upload of a signed client certificate": {
			request: UploadSignedClientCertificateRequest{
				CertificateID: 123,
				Version:       1,
				Body: UploadSignedClientCertificateRequestBody{
					Certificate: "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
				},
			},
			responseStatus: http.StatusOK,
			expectedRequestBody: `{
					"certificate": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"
 			}`,
			expectedPath: "/mtls-origin-keystore/v1/client-certificates/123/versions/1/certificate-block",
		},
		"Validation error - missing CertificateID": {
			request: UploadSignedClientCertificateRequest{
				Version: 1,
				Body: UploadSignedClientCertificateRequestBody{
					Certificate: "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "uploading client certificate version: validation failed: CertificateID: cannot be blank", err.Error())
			},
		},
		"Validation error - missing Version and Certificate": {
			request: UploadSignedClientCertificateRequest{
				CertificateID: 123,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "uploading client certificate version: validation failed: Certificate: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"Validation error - missing Certificate": {
			request: UploadSignedClientCertificateRequest{
				CertificateID: 123,
				Version:       1,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "uploading client certificate version: validation failed: Certificate: cannot be blank", err.Error())
			},
		},
		"Validation error - Certificate is empty": {
			request: UploadSignedClientCertificateRequest{
				CertificateID: 123,
				Version:       1,
				Body: UploadSignedClientCertificateRequestBody{
					Certificate: "",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "uploading client certificate version: validation failed: Certificate: cannot be blank", err.Error())
			},
		},
		"Validation error - TrustChain provided is empty": {
			request: UploadSignedClientCertificateRequest{
				CertificateID: 123,
				Version:       1,
				Body: UploadSignedClientCertificateRequestBody{
					Certificate: "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					TrustChain:  ptr.To(""),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "uploading client certificate version: validation failed: TrustChain: cannot be blank", err.Error())
			},
		},
		"Error Response 404 - Client Certificate not found": {
			request: UploadSignedClientCertificateRequest{
				CertificateID: 123,
				Version:       1,
				Body: UploadSignedClientCertificateRequestBody{
					Certificate: "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					TrustChain:  ptr.To("-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"),
				},
			},
			expectedRequestBody: `{
					"certificate": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					"trustChain": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"
			}`,
			responseStatus: http.StatusNotFound,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions/1/certificate-block",
			responseBody: `
				{
				  "detail": "The requested resource could not be found on the server.",
				  "field": "certificateId",
				  "instance": "32be8a77-cc41-495e-94d8-cb536afea149",
				  "problemId": "32be8a77-cc41-495e-94d8-cb536afea149",
				  "status": 404,
				  "title": "Resource Not Found",
				  "type": "resource-not-found",
				  "value": "180131"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrClientCertificateNotFound))
			},
		},
		"Error Response 400- Duplicate Certificate name": {
			request: UploadSignedClientCertificateRequest{
				CertificateID: 123,
				Version:       1,
				Body: UploadSignedClientCertificateRequestBody{
					Certificate: "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					TrustChain:  ptr.To("-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"),
				},
			},
			expectedRequestBody: `{
					"certificate": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					"trustChain": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"
			}`,
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions/1/certificate-block",
			responseBody: `{
				  "detail": "Bad Request",
				  "errors": [
					{
					  "detail": "Certificate with same name already exists.",
					  "field": "certificateName",
					  "problemId": "00c4e7b5-dc7b-43f9-9f18-de8de3c82527",
					  "title": "Invalid Input",
					  "type": "error-types/invalid"
					}
				  ],
				  "instance": "/f311c60f-9914-4e23-be2c-db8dbb711a8a",
				  "status": 400,
				  "title": "Bad Request",
				  "type": "bad-request"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrInvalidClientCertificate))
			},
		},
		"Error Response 400- Client Certificate is invalid": {
			request: UploadSignedClientCertificateRequest{
				CertificateID: 123,
				Version:       1,
				Body: UploadSignedClientCertificateRequestBody{
					Certificate: "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					TrustChain:  ptr.To("-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"),
				},
			},
			expectedRequestBody: `{
					"certificate": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----",
					"trustChain": "-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----"
			}`,
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/mtls-origin-keystore/v1/client-certificates/123/versions/1/certificate-block",
			responseBody: `
				{
				  "detail": "Bad Request",
				  "errors": [
					{
					  "detail": "Certificate is either invalid or cannot not be accepted.",
					  "field": "certificate",
					  "problemId": "f1658521-6e7d-4051-89e6-5dce765029af",
					  "title": "Invalid Input",
					  "type": "error-types/invalid"
					}
				  ],
				  "instance": "7cc378a1-ba9e-41a6-8255-27ee0a5a3897",
				  "problemId": "7cc378a1-ba9e-41a6-8255-27ee0a5a3897",
				  "status": 400,
				  "title": "Bad Request",
				  "type": "bad-request"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrInvalidClientCertificate))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			err := client.UploadSignedClientCertificate(context.Background(), test.request)

			if test.withError != nil {
				test.withError(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
