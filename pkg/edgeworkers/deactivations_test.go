package edgeworkers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDeactivations(t *testing.T) {
	tests := map[string]struct {
		params         ListDeactivationsRequest
		withError      error
		expectedPath   string
		responseStatus int
		responseBody   string
		expectedResult []Deactivation
	}{
		"400 bad request": {
			params:    ListDeactivationsRequest{},
			withError: ErrStructValidation,
		},
		"500 internal server error": {
			params:         ListDeactivationsRequest{EdgeWorkerID: 5},
			expectedPath:   "/edgeworkers/v1/ids/5/deactivations",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
			  "detail": "An error occurred while fetching the activation",
			  "instance": "/edgeworkers/error-instances/8b95c971-f6b3-4479-a393-202be75e43e1",
			  "status": 500,
			  "title": "An unexpected error has occurred.",
			  "type": "/edgeworkers/error-types/edgeworkers-server-error",
			  "errorCode": "EW4303"
			}`,
			withError: &Error{
				Detail:    "An error occurred while fetching the activation",
				Instance:  "/edgeworkers/error-instances/8b95c971-f6b3-4479-a393-202be75e43e1",
				Status:    http.StatusInternalServerError,
				Title:     "An unexpected error has occurred.",
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				ErrorCode: "EW4303",
			},
		},
		"200 OK": {
			responseBody: `{
					"deactivations":[
						{
							"edgeWorkerId":42,
							"version":"2",
							"deactivationId":3,
							"accountId":"1-AB2CD34",
							"status":"PENDING",
							"network":"PRODUCTION",
							"note":"EdgeWorker ID 42 is no longer used in production.",
							"createdBy":"jdoe",
							"createdTime":"2020-07-09T09:03:28Z",
							"lastModifiedTime":"2020-07-09T09:04:42Z"
						},
						{
							"edgeWorkerId":42,
							"version":"1",
							"deactivationId":1,
							"accountId":"1-AB2CD34",
							"status":"IN_PROGRESS",
							"network":"STAGING",
							"createdBy":"jsmith",
							"createdTime":"2020-07-09T08:13:54Z",
							"lastModifiedTime":"2020-07-09T08:35:02Z"
						},
						{
							"edgeWorkerId":42,
							"version":"2",
							"deactivationId":2,
							"accountId":"1-AB2CD34",
							"status":"COMPLETE",
							"network":"PRODUCTION",
							"createdBy":"asmith",
							"createdTime":"2020-07-10T14:23:42Z",
							"lastModifiedTime":"2020-07-10T14:53:25Z"
						}
					]
				}`,
			expectedPath:   "/edgeworkers/v1/ids/42/deactivations",
			responseStatus: http.StatusOK,
			expectedResult: []Deactivation{
				{
					EdgeWorkerID:     42,
					Version:          "2",
					DeactivationID:   3,
					AccountID:        "1-AB2CD34",
					Status:           "PENDING",
					Network:          ActivationNetworkProduction,
					Note:             "EdgeWorker ID 42 is no longer used in production.",
					CreatedBy:        "jdoe",
					CreatedTime:      "2020-07-09T09:03:28Z",
					LastModifiedTime: "2020-07-09T09:04:42Z",
				},
				{
					EdgeWorkerID:     42,
					Version:          "1",
					DeactivationID:   1,
					AccountID:        "1-AB2CD34",
					Status:           "IN_PROGRESS",
					Network:          ActivationNetworkStaging,
					CreatedBy:        "jsmith",
					CreatedTime:      "2020-07-09T08:13:54Z",
					LastModifiedTime: "2020-07-09T08:35:02Z",
				},
				{
					EdgeWorkerID:     42,
					Version:          "2",
					DeactivationID:   2,
					AccountID:        "1-AB2CD34",
					Status:           "COMPLETE",
					Network:          ActivationNetworkProduction,
					CreatedBy:        "asmith",
					CreatedTime:      "2020-07-10T14:23:42Z",
					LastModifiedTime: "2020-07-10T14:53:25Z",
				},
			},
			params: ListDeactivationsRequest{EdgeWorkerID: 42},
		},
		"200 OK with version": {
			responseBody: `{
					"deactivations":[
						{
							"edgeWorkerId":41,
							"version":"2",
							"deactivationId":3,
							"accountId":"1-AB2CD34",
							"status":"PENDING",
							"network":"PRODUCTION",
							"note":"EdgeWorker ID 41 is no longer used in production.",
							"createdBy":"jdoe",
							"createdTime":"2020-07-09T09:03:28Z",
							"lastModifiedTime":"2020-07-09T09:04:42Z"
						},
						{
							"edgeWorkerId":41,
							"version":"2",
							"deactivationId":2,
							"accountId":"1-AB2CD34",
							"status":"COMPLETE",
							"network":"PRODUCTION",
							"createdBy":"asmith",
							"createdTime":"2020-07-10T14:23:42Z",
							"lastModifiedTime":"2020-07-10T14:53:25Z"
						}
					]
				}`,
			expectedPath:   "/edgeworkers/v1/ids/42/deactivations?version=2",
			responseStatus: http.StatusOK,
			expectedResult: []Deactivation{
				{
					EdgeWorkerID:     41,
					Version:          "2",
					DeactivationID:   3,
					AccountID:        "1-AB2CD34",
					Status:           "PENDING",
					Network:          ActivationNetworkProduction,
					Note:             "EdgeWorker ID 41 is no longer used in production.",
					CreatedBy:        "jdoe",
					CreatedTime:      "2020-07-09T09:03:28Z",
					LastModifiedTime: "2020-07-09T09:04:42Z",
				},
				{
					EdgeWorkerID:     41,
					Version:          "2",
					DeactivationID:   2,
					AccountID:        "1-AB2CD34",
					Status:           "COMPLETE",
					Network:          ActivationNetworkProduction,
					CreatedBy:        "asmith",
					CreatedTime:      "2020-07-10T14:23:42Z",
					LastModifiedTime: "2020-07-10T14:53:25Z",
				},
			},
			params: ListDeactivationsRequest{
				EdgeWorkerID: 42,
				Version:      "2",
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

			result, err := client.ListDeactivations(context.Background(), test.params)
			if test.withError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, test.withError))
				return
			}

			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(result.Deactivations, test.expectedResult))
		})
	}
}

func TestListDeactivationsRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		params ListDeactivationsRequest
		errors validation.Errors
	}{
		"no EW ID": {
			params: ListDeactivationsRequest{},
			errors: validation.Errors{
				"EdgeWorkerID": validation.ErrorObject{}.SetCode("validation_required").SetMessage("cannot be blank"),
			},
		},
		"EW ID": {
			params: ListDeactivationsRequest{EdgeWorkerID: 1},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.params.Validate()
			if len(test.errors) != 0 {
				require.Error(t, err)
				assert.Equal(t, test.errors, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestEdgeWorkerDeactivateVersionRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		params DeactivateVersionRequest
		errors *regexp.Regexp
	}{
		"no EW ID": {
			params: DeactivateVersionRequest{
				DeactivateVersion: DeactivateVersion{
					Version: "--",
					Network: ActivationNetworkProduction,
				},
			},
			errors: regexp.MustCompile(`EdgeWorkerID.+cannot be blank.+`),
		},
		"no version": {
			params: DeactivateVersionRequest{
				EdgeWorkerID: 1,
				DeactivateVersion: DeactivateVersion{
					Network: ActivationNetworkProduction,
				},
			},
			errors: regexp.MustCompile(`DeactivateVersion:.+Version:.+cannot be blank.+`),
		},
		"bad network": {
			params: DeactivateVersionRequest{
				EdgeWorkerID: 1,
				DeactivateVersion: DeactivateVersion{
					Network: "-asdfa",
					Version: "a",
				},
			},
			errors: regexp.MustCompile(`DeactivateVersion:.+Network:.+value '-asdfa' is invalid. Must be one of: 'STAGING' or 'PRODUCTION'.+`),
		},
		"no network": {
			params: DeactivateVersionRequest{
				EdgeWorkerID: 1,
				DeactivateVersion: DeactivateVersion{
					Version: "a",
				},
			},
			errors: regexp.MustCompile(`DeactivateVersion:.+Network:.+cannot be blank.+`),
		},
		"ok": {
			params: DeactivateVersionRequest{
				EdgeWorkerID: 1,
				DeactivateVersion: DeactivateVersion{
					Version: "asdf",
					Network: ActivationNetworkStaging,
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.params.Validate()
			if test.errors != nil {
				require.Error(t, err)
				assert.Regexp(t, test.errors, err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestEdgeworkers_DeactivateVersion(t *testing.T) {
	tests := map[string]struct {
		params           DeactivateVersionRequest
		withError        error
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse Deactivation
	}{
		"400 bad request": {
			params:    DeactivateVersionRequest{},
			withError: ErrStructValidation,
		},
		"201 created": {
			params: DeactivateVersionRequest{
				EdgeWorkerID: 1,
				DeactivateVersion: DeactivateVersion{
					Version: "123",
					Network: ActivationNetworkProduction,
					Note:    "not used",
				},
			},
			expectedPath: "/edgeworkers/v1/ids/1/deactivations",
			responseBody: `{
				"edgeWorkerId": 1,
				"version": "123",
				"deactivationId": 1,
				"accountId": "B-3-WNKA6P",
				"status": "PRESUBMIT",
				"network": "PRODUCTION",
				"note": "not used",
				"createdBy": "agrebenk",
				"createdTime": "2021-12-17T10:07:35Z",
				"lastModifiedTime": "2021-12-17T10:07:35Z"
			}`,
			responseStatus: http.StatusCreated,
			expectedResponse: Deactivation{
				EdgeWorkerID:     1,
				Version:          "123",
				DeactivationID:   1,
				AccountID:        "B-3-WNKA6P",
				Status:           "PRESUBMIT",
				Network:          ActivationNetworkProduction,
				Note:             "not used",
				CreatedBy:        "agrebenk",
				CreatedTime:      "2021-12-17T10:07:35Z",
				LastModifiedTime: "2021-12-17T10:07:35Z",
			},
		},
		"500 server error": {
			params: DeactivateVersionRequest{
				EdgeWorkerID: 1,
				DeactivateVersion: DeactivateVersion{
					Version: "123",
					Network: ActivationNetworkProduction,
					Note:    "not used",
				},
			},
			expectedPath:   "/edgeworkers/v1/ids/1/deactivations",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
			  "detail": "An error occurred while fetching the activation",
			  "instance": "/edgeworkers/error-instances/8b95c971-f6b3-4479-a393-202be75e43e1",
			  "status": 500,
			  "title": "An unexpected error has occurred.",
			  "type": "/edgeworkers/error-types/edgeworkers-server-error",
			  "errorCode": "EW4303"
			}`,
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "An error occurred while fetching the activation",
				Instance:  "/edgeworkers/error-instances/8b95c971-f6b3-4479-a393-202be75e43e1",
				Status:    http.StatusInternalServerError,
				ErrorCode: "EW4303",
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
			client := mockAPIClient(t, mockServer)

			response, err := client.DeactivateVersion(context.Background(), test.params)
			if test.withError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, test.withError))
				return
			}
			require.NoError(t, err)

			assert.True(t, reflect.DeepEqual(test.expectedResponse, *response))
		})
	}
}

func TestEdgeWorkerGetDeactivationRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		request GetDeactivationRequest
		errors  validation.Errors
	}{
		"no EW ID": {
			request: GetDeactivationRequest{DeactivationID: 1},
			errors: map[string]error{
				"EdgeWorkerID": validation.ErrorObject{}.SetCode("validation_required").SetMessage("cannot be blank"),
			},
		},
		"no Deactivation ID": {
			request: GetDeactivationRequest{EdgeWorkerID: 1},
			errors: map[string]error{
				"DeactivationID": validation.ErrorObject{}.SetCode("validation_required").SetMessage("cannot be blank"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.request.Validate()
			if len(test.errors) != 0 {
				require.Error(t, err)
				assert.Equal(t, test.errors, err)
				return
			}
			assert.True(t, false)
		})
	}
}

func TestEdgeworkers_GetDeactivation(t *testing.T) {
	tests := map[string]struct {
		request              GetDeactivationRequest
		expectedDeactivation Deactivation
		withError            error
		expectedPath         string
		responseStatus       int
		responseBody         string
	}{
		"request validation error": {
			request:   GetDeactivationRequest{},
			withError: ErrStructValidation,
		},
		"404 deactivation not found": {
			request: GetDeactivationRequest{
				EdgeWorkerID:   1,
				DeactivationID: 2,
			},
			responseStatus: http.StatusNotFound,
			responseBody: `{
			  "detail": "Unable to find the requested EdgeWorker ID",
			  "errorCode": "EW2002",
			  "instance": "/edgeworkers/error-instances/76b1595d-08e5-46a8-8bc6-72d01e621303",
			  "status": 404,
			  "title": "The given resource could not be found.",
			  "type": "/edgeworkers/error-types/edgeworkers-bad-request"
			}`,
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-bad-request",
				Title:     "The given resource could not be found.",
				Detail:    "Unable to find the requested EdgeWorker ID",
				Instance:  "/edgeworkers/error-instances/76b1595d-08e5-46a8-8bc6-72d01e621303",
				Status:    http.StatusNotFound,
				ErrorCode: "EW2002",
			},
			expectedPath: "/edgeworkers/v1/ids/1/deactivations/2",
		},
		"200 ok": {
			request: GetDeactivationRequest{
				EdgeWorkerID:   1,
				DeactivationID: 2,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/edgeworkers/v1/ids/1/deactivations/2",
			responseBody: `{
				"edgeWorkerId": 1,
				"version": "1.0aSDA",
				"deactivationId": 2,
				"accountId": "B-3-WNKA6P",
				"status": "COMPLETE",
				"network": "PRODUCTION",
				"note": "not used",
				"createdBy": "agrebenk",
				"createdTime": "2021-12-17T10:07:35Z",
				"lastModifiedTime": "2021-12-17T10:24:02Z"
			}`,
			expectedDeactivation: Deactivation{
				EdgeWorkerID:     1,
				Version:          "1.0aSDA",
				DeactivationID:   2,
				AccountID:        "B-3-WNKA6P",
				Status:           "COMPLETE",
				Network:          ActivationNetworkProduction,
				Note:             "not used",
				CreatedBy:        "agrebenk",
				CreatedTime:      "2021-12-17T10:07:35Z",
				LastModifiedTime: "2021-12-17T10:24:02Z",
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

			response, err := client.GetDeactivation(context.Background(), test.request)
			if test.withError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, test.withError))
				return
			}
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(test.expectedDeactivation, *response))
		})
	}
}
