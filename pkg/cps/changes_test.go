package cps

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetChangeStatus(t *testing.T) {
	tests := map[string]struct {
		params           GetChangeStatusRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Change
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetChangeStatusRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "statusInfo": {
        "status": "wait-upload-third-party",
        "state": "awaiting-input",
        "description": "Waiting for you to upload and submit your third party certificate and trust chain.",
        "error": null,
        "deploymentSchedule": {
            "notBefore": null,
            "notAfter": null
        }
    },
    "allowedInput": [
        {
            "type": "third-party-certificate",
            "requiredToProceed": true,
            "info": "/cps/v2/enrollments/1/changes/2/input/info/third-party-csr",
            "update": "/cps/v2/enrollments/1/changes/2/input/update/third-party-cert-and-trust-chain"
        }
    ]
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2",
			expectedResponse: &Change{
				AllowedInput: []AllowedInput{
					{
						Info:              "/cps/v2/enrollments/1/changes/2/input/info/third-party-csr",
						RequiredToProceed: true,
						Type:              "third-party-certificate",
						Update:            "/cps/v2/enrollments/1/changes/2/input/update/third-party-cert-and-trust-chain",
					},
				},
				StatusInfo: &StatusInfo{
					DeploymentSchedule: &DeploymentSchedule{},
					Description:        "Waiting for you to upload and submit your third party certificate and trust chain.",
					State:              "awaiting-input",
					Status:             "wait-upload-third-party",
				},
			},
		},
		"500 internal server error": {
			params: GetChangeStatusRequest{
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
			expectedPath: "/cps/v2/enrollments/1/changes/2",
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
			params: GetChangeStatusRequest{},
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
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetChangeStatus(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCancelChange(t *testing.T) {
	tests := map[string]struct {
		request          CancelChangeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CancelChangeResponse
		withError        error
	}{
		"200 OK": {
			request: CancelChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"change": "/cps/v2/enrollments/1/changes/2"
}`,
			expectedPath:     "/cps/v2/enrollments/1/changes/2",
			expectedResponse: &CancelChangeResponse{Change: "/cps/v2/enrollments/1/changes/2"},
		},
		"500 internal server error": {
			request: CancelChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error canceling change",
    "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error canceling change",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request:   CancelChangeRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.change-id.v1+json", r.Header.Get("Accept"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CancelChange(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateChange(t *testing.T) {
	tests := map[string]struct {
		request             UpdateChangeRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *UpdateChangeResponse
		expectedContentType string
		withError           error
	}{
		"200 ok, dv validation": {
			request: UpdateChangeRequest{
				Certificate: Certificate{
					Certificate: "test-cert",
					TrustChain:  "test-trust-chain",
				},
				EnrollmentID:          1,
				ChangeID:              2,
				AllowedInputTypeParam: "lets-encrypt-challenges-completed",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	 "change": "/cps/v2/enrollments/1/changes/2"
}`,
			expectedPath:        "/cps/v2/enrollments/1/changes/2/input/update/lets-encrypt-challenges-completed",
			expectedResponse:    &UpdateChangeResponse{Change: "/cps/v2/enrollments/1/changes/2"},
			expectedContentType: "application/vnd.akamai.cps.acknowledgement.v1+json",
		},
		"200 ok, third party": {
			request: UpdateChangeRequest{
				Certificate: Certificate{
					Certificate: "test-cert",
					TrustChain:  "test-trust-chain",
				},
				EnrollmentID:          1,
				ChangeID:              2,
				AllowedInputTypeParam: "third-party-cert-and-trust-chain",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	 "change": "/cps/v2/enrollments/1/changes/2"
}`,
			expectedPath:        "/cps/v2/enrollments/1/changes/2/input/update/third-party-cert-and-trust-chain",
			expectedResponse:    &UpdateChangeResponse{Change: "/cps/v2/enrollments/1/changes/2"},
			expectedContentType: "application/vnd.akamai.cps.certificate-and-trust-chain.v1+json",
		},
		"500 internal server error": {
			request: UpdateChangeRequest{
				Certificate: Certificate{
					Certificate: "test-cert",
					TrustChain:  "test-trust-chain",
				},
				EnrollmentID:          1,
				ChangeID:              2,
				AllowedInputTypeParam: "third-party-cert-and-trust-chain",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
 "title": "Internal Server Error",
 "detail": "Error updating change",
 "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/update/third-party-cert-and-trust-chain",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating change",
				StatusCode: http.StatusInternalServerError,
			},
			expectedContentType: "application/vnd.akamai.cps.certificate-and-trust-chain.v1+json",
		},
		"validation error, invalid allowed input type param": {
			request: UpdateChangeRequest{
				Certificate: Certificate{
					Certificate: "test-cert",
					TrustChain:  "test-trust-chain",
				},
				EnrollmentID:          1,
				ChangeID:              2,
				AllowedInputTypeParam: "abc",
			},
			withError: ErrStructValidation,
		},
		"validation error, no enrollment id": {
			request: UpdateChangeRequest{
				Certificate: Certificate{
					Certificate: "test-cert",
					TrustChain:  "test-trust-chain",
				},
				ChangeID:              2,
				AllowedInputTypeParam: "third-party-cert-and-trust-chain",
			},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.change-id.v1+json", r.Header.Get("Accept"))
				assert.Equal(t, test.expectedContentType, r.Header.Get("Content-Type"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateChange(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetChangeLetsEncryptChallenges(t *testing.T) {
	tests := map[string]struct {
		params           GetChangeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DvChallenges
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
    "dv": [
        {
            "status": "Awaiting user",
            "error": "The domain is not ready for validation.",
            "validationStatus": "RESPONSE_ERROR",
            "requestTimestamp": "2018-09-05T15:55:49Z",
            "validatedTimestamp": "2018-09-05T17:53:22Z",
            "expires": "2018-09-06T17:55:17Z",
            "challenges": [
                {
                    "type": "dns-01",
                    "status": "pending",
                    "error": null,
                    "token": "cGBnw-3YO7rUhq61EuuHqcGrYkaQWALAgi8szTqRoHA",
                    "responseBody": "0yVISDJjpXR7BXzR5QgfA51tt-I6aKremGnPwK_lvH4",
                    "fullPath": "_acme-challenge.www.cps-example-dv.com.",
                    "redirectFullPath": "",
                    "validationRecords": []
                }
            ],
            "domain": "www.cps-example-dv.com"
        }
    ]
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/lets-encrypt-challenges",
			expectedResponse: &DvChallenges{DV: []DV{
				{
					Challenges: []Challenges{
						{
							FullPath:          "_acme-challenge.www.cps-example-dv.com.",
							ResponseBody:      "0yVISDJjpXR7BXzR5QgfA51tt-I6aKremGnPwK_lvH4",
							Status:            "pending",
							Token:             "cGBnw-3YO7rUhq61EuuHqcGrYkaQWALAgi8szTqRoHA",
							Type:              "dns-01",
							ValidationRecords: []ValidationRecords{},
						},
					},
					Domain:             "www.cps-example-dv.com",
					Error:              "The domain is not ready for validation.",
					Expires:            "2018-09-06T17:55:17Z",
					RequestTimestamp:   "2018-09-05T15:55:49Z",
					Status:             "Awaiting user",
					ValidatedTimestamp: "2018-09-05T17:53:22Z",
					ValidationStatus:   "RESPONSE_ERROR",
				},
			}},
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
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/lets-encrypt-challenges",
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
				assert.Equal(t, "application/vnd.akamai.cps.dv-challenges.v2+json", r.Header.Get("Accept"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetChangeLetsEncryptChallenges(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
