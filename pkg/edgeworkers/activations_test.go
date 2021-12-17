package edgeworkers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListActivations(t *testing.T) {
	tests := map[string]struct {
		params           ListActivationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListActivationsResponse
		withError        error
	}{
		"200 OK": {
			params: ListActivationsRequest{
				EdgeWorkerID: 42,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "activations": [
        {
            "edgeWorkerId": 42,
            "version": "2",
            "activationId": 3,
            "accountId": "B-M-1KQK3WU",
            "status": "PENDING",
            "network": "PRODUCTION",
            "createdBy": "jdoe",
            "createdTime": "2018-07-09T09:03:28Z",
            "lastModifiedTime": "2018-07-09T09:04:42Z"
        },
        {
            "edgeWorkerId": 42,
            "version": "1",
            "activationId": 1,
            "accountId": "B-M-1KQK3WU",
            "status": "IN_PROGRESS",
            "network": "STAGING",
            "createdBy": "jsmith",
            "createdTime": "2018-07-09T08:13:54Z",
            "lastModifiedTime": "2018-07-09T08:35:02Z"
        }
    ]
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations",
			expectedResponse: &ListActivationsResponse{
				Activations: []Activation{
					{
						AccountID:        "B-M-1KQK3WU",
						ActivationID:     3,
						CreatedBy:        "jdoe",
						CreatedTime:      "2018-07-09T09:03:28Z",
						EdgeWorkerID:     42,
						LastModifiedTime: "2018-07-09T09:04:42Z",
						Network:          "PRODUCTION",
						Status:           "PENDING",
						Version:          "2",
					},
					{
						AccountID:        "B-M-1KQK3WU",
						ActivationID:     1,
						CreatedBy:        "jsmith",
						CreatedTime:      "2018-07-09T08:13:54Z",
						EdgeWorkerID:     42,
						LastModifiedTime: "2018-07-09T08:35:02Z",
						Network:          "STAGING",
						Status:           "IN_PROGRESS",
						Version:          "1",
					},
				},
			},
		},
		"200 OK with version query": {
			params: ListActivationsRequest{
				EdgeWorkerID: 42,
				Version:      "1",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "activations": [
        {
            "edgeWorkerId": 42,
            "version": "1",
            "activationId": 1,
            "accountId": "B-M-1KQK3WU",
            "status": "IN_PROGRESS",
            "network": "STAGING",
            "createdBy": "jsmith",
            "createdTime": "2018-07-09T08:13:54Z",
            "lastModifiedTime": "2018-07-09T08:35:02Z"
        }
    ]
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations?version=1",
			expectedResponse: &ListActivationsResponse{
				Activations: []Activation{
					{
						AccountID:        "B-M-1KQK3WU",
						ActivationID:     1,
						CreatedBy:        "jsmith",
						CreatedTime:      "2018-07-09T08:13:54Z",
						EdgeWorkerID:     42,
						LastModifiedTime: "2018-07-09T08:35:02Z",
						Network:          "STAGING",
						Status:           "IN_PROGRESS",
						Version:          "1",
					},
				},
			},
		},
		"500 internal server error": {
			params: ListActivationsRequest{
				EdgeWorkerID: 42,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/edgeworkers/error-instances/abc",
    "status": 500,
    "errorCode": "EW4303"
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
			},
		},
		"missing edge worker id": {
			params:    ListActivationsRequest{},
			withError: ErrStructValidation,
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
			result, err := client.ListActivations(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetActivation(t *testing.T) {
	tests := map[string]struct {
		params           GetActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Activation
		withError        error
	}{
		"200 OK": {
			params: GetActivationRequest{
				EdgeWorkerID: 42,
				ActivationID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"edgeWorkerId": 42,
	"version": "1",
	"activationId": 1,
	"accountId": "B-M-1KQK3WU",
	"status": "IN_PROGRESS",
	"network": "STAGING",
	"createdBy": "jsmith",
	"createdTime": "2018-07-09T08:13:54Z",
	"lastModifiedTime": "2018-07-09T08:35:02Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations/1",
			expectedResponse: &Activation{
				AccountID:        "B-M-1KQK3WU",
				ActivationID:     1,
				CreatedBy:        "jsmith",
				CreatedTime:      "2018-07-09T08:13:54Z",
				EdgeWorkerID:     42,
				LastModifiedTime: "2018-07-09T08:35:02Z",
				Network:          "STAGING",
				Status:           "IN_PROGRESS",
				Version:          "1",
			},
		},
		"500 internal server error": {
			params: GetActivationRequest{
				EdgeWorkerID: 42,
				ActivationID: 1,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/edgeworkers/error-instances/abc",
    "status": 500,
    "errorCode": "EW4303"
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations/1",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
			},
		},
		"missing activation id": {
			params: GetActivationRequest{
				EdgeWorkerID: 42,
			},
			withError: ErrStructValidation,
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
			result, err := client.GetActivation(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreateActivation(t *testing.T) {
	tests := map[string]struct {
		params           CreateActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Activation
		withError        error
	}{
		"200 OK": {
			params: CreateActivationRequest{
				EdgeWorkerID: 42,
				CreateActivation: CreateActivation{
					Network: "STAGING",
					Version: "1",
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"edgeWorkerId": 42,
	"version": "1",
	"activationId": 1,
	"accountId": "B-M-1KQK3WU",
	"status": "PRESUBMIT",
	"network": "STAGING",
	"createdBy": "jsmith",
	"createdTime": "2018-07-09T08:13:54Z",
	"lastModifiedTime": "2018-07-09T08:35:02Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations",
			expectedResponse: &Activation{
				AccountID:        "B-M-1KQK3WU",
				ActivationID:     1,
				CreatedBy:        "jsmith",
				CreatedTime:      "2018-07-09T08:13:54Z",
				EdgeWorkerID:     42,
				LastModifiedTime: "2018-07-09T08:35:02Z",
				Network:          "STAGING",
				Status:           "PRESUBMIT",
				Version:          "1",
			},
		},
		"500 internal server error": {
			params: CreateActivationRequest{
				EdgeWorkerID: 42,
				CreateActivation: CreateActivation{
					Network: "STAGING",
					Version: "1",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/edgeworkers/error-instances/abc",
    "status": 500,
    "errorCode": "EW4303"
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
			},
		},
		"missing edge worker id": {
			params: CreateActivationRequest{
				CreateActivation: CreateActivation{
					Network: ActivationNetworkStaging,
					Version: "1",
				},
			},
			withError: ErrStructValidation,
		},
		"invalid network": {
			params: CreateActivationRequest{
				CreateActivation: CreateActivation{
					Network: "invalid",
					Version: "1",
				},
			},
			withError: ErrStructValidation,
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
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateActivation(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCancelActivation(t *testing.T) {
	tests := map[string]struct {
		params           CancelActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Activation
		withError        error
	}{
		"200 OK": {
			params: CancelActivationRequest{
				EdgeWorkerID: 42,
				ActivationID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"edgeWorkerId": 42,
	"version": "1",
	"activationId": 1,
	"accountId": "B-M-1KQK3WU",
	"status": "CANCELED",
	"network": "STAGING",
	"createdBy": "jsmith",
	"createdTime": "2018-07-09T08:13:54Z",
	"lastModifiedTime": "2018-07-09T08:35:02Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations/1",
			expectedResponse: &Activation{
				AccountID:        "B-M-1KQK3WU",
				ActivationID:     1,
				CreatedBy:        "jsmith",
				CreatedTime:      "2018-07-09T08:13:54Z",
				EdgeWorkerID:     42,
				LastModifiedTime: "2018-07-09T08:35:02Z",
				Network:          "STAGING",
				Status:           "CANCELED",
				Version:          "1",
			},
		},
		"500 internal server error": {
			params: CancelActivationRequest{
				EdgeWorkerID: 42,
				ActivationID: 1,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/edgeworkers/error-instances/abc",
    "status": 500,
    "errorCode": "EW4303"
}`,
			expectedPath: "/edgeworkers/v1/ids/42/activations/1",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
			},
		},
		"missing activation id": {
			params: CancelActivationRequest{
				EdgeWorkerID: 42,
			},
			withError: ErrStructValidation,
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
			client := mockAPIClient(t, mockServer)
			result, err := client.CancelActivation(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
