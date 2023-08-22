package cloudwrapper

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudwrapper_ListAuthKeys(t *testing.T) {
	tests := map[string]struct {
		params           ListAuthKeysRequest
		expectedPath     string
		responseBody     string
		expectedResponse *ListAuthKeysResponse
		responseStatus   int
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListAuthKeysRequest{
				ContractID: "test_contract",
				CDNCode:    "dn123",
			},
			expectedPath:   "/cloud-wrapper/v1/multi-cdn/auth-keys?cdnCode=dn123&contractId=test_contract",
			responseStatus: http.StatusOK,
			responseBody: `{
    "cdnAuthKeys": [
        {
            "authKeyName": "test7",
            "expiryDate": "2023-08-08",
            "headerName": "key"
        }
    ]
}`,
			expectedResponse: &ListAuthKeysResponse{
				CDNAuthKeys: []MultiCDNAuthKey{
					{
						AuthKeyName: "test7",
						ExpiryDate:  "2023-08-08",
						HeaderName:  "key",
					},
				},
			},
		},
		"missing CDNCode": {
			params: ListAuthKeysRequest{ContractID: "test_contract"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, err.Error(), "list auth keys: struct validation: CDNCode: cannot be blank")
			},
		},
		"missing ContractID": {
			params: ListAuthKeysRequest{CDNCode: "dn123"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, err.Error(), "list auth keys: struct validation: ContractID: cannot be blank")
			},
		},
		"500 internal server error": {
			params: ListAuthKeysRequest{
				ContractID: "test_contract",
				CDNCode:    "dn123",
			},
			expectedPath:   "/cloud-wrapper/v1/multi-cdn/auth-keys?cdnCode=dn123&contractId=test_contract",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error processing request",
					"status": 500
				}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error processing request",
					Status: http.StatusInternalServerError,
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
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			users, err := client.ListAuthKeys(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}

func TestCloudwrapper_ListCDNProviders(t *testing.T) {
	tests := map[string]struct {
		expectedPath     string
		responseBody     string
		expectedResponse *ListCDNProvidersResponse
		responseStatus   int
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			expectedPath:   "/cloud-wrapper/v1/multi-cdn/providers",
			responseStatus: http.StatusOK,
			responseBody: `{
    "cdnProviders": [
        {
            "cdnCode": "dn002",
            "cdnName": "Level 3 (Centurylink)"
        },
        {
            "cdnCode": "dn003",
            "cdnName": "Limelight"
        },
        {
            "cdnCode": "dn004",
            "cdnName": "CloudFront"
        }
    ]
}`,
			expectedResponse: &ListCDNProvidersResponse{
				CDNProviders: []CDNProvider{
					{
						CDNCode: "dn002",
						CDNName: "Level 3 (Centurylink)",
					},
					{
						CDNCode: "dn003",
						CDNName: "Limelight",
					},
					{
						CDNCode: "dn004",
						CDNName: "CloudFront",
					},
				},
			},
		},
		"500 internal server error": {
			expectedPath:   "/cloud-wrapper/v1/multi-cdn/providers",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error processing request",
					"status": 500
				}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error processing request",
					Status: http.StatusInternalServerError,
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
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			users, err := client.ListCDNProviders(context.Background())
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}
