package cps

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
