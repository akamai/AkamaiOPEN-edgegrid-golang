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

func TestInitializeEdgeKV(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *EdgeKVInitializationStatus
		withError        error
	}{
		"201 Created": {
			responseStatus: http.StatusCreated,
			responseBody: `
			{
				"accountStatus": "INITIALIZED",
				"cpcode": "123456",
				"productionStatus": "INITIALIZED",
				"stagingStatus": "INITIALIZED"
			}`,
			expectedPath: "/edgekv/v1/initialize",
			expectedResponse: &EdgeKVInitializationStatus{
				AccountStatus:    "INITIALIZED",
				CPCode:           "123456",
				ProductionStatus: "INITIALIZED",
				StagingStatus:    "INITIALIZED",
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "https://learn.akamai.com",
				"title": "Internal Server Error",
				"detail": "An internal error occurred.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 500,
				"errorCode": "EKV_0000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,

			expectedPath: "/edgekv/v1/initialize",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Internal Server Error",
				Detail:    "An internal error occurred.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    500,
				ErrorCode: "EKV_0000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
			},
		},
		"503 service unavailable error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "https://learn.akamai.com",
				"title": "Service Unavailable Error",
				"detail": "An internal error occurred.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 503,
				"errorCode": "EKV_0000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/initialize",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Service Unavailable Error",
				Detail:    "An internal error occurred.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    503,
				ErrorCode: "EKV_0000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.InitializeEdgeKV(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetEdgeKVInitializeStatus(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *EdgeKVInitializationStatus
		withError        error
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountStatus": "INITIALIZED",
				"cpcode": "123456",
				"productionStatus": "INITIALIZED",
				"stagingStatus": "INITIALIZED"
			}`,
			expectedPath: "/edgekv/v1/initialize",
			expectedResponse: &EdgeKVInitializationStatus{
				AccountStatus:    "INITIALIZED",
				CPCode:           "123456",
				ProductionStatus: "INITIALIZED",
				StagingStatus:    "INITIALIZED",
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "https://learn.akamai.com",
				"title": "Internal Server Error",
				"detail": "An internal error occurred.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 500,
				"errorCode": "EKV_0000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/initialize",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Internal Server Error",
				Detail:    "An internal error occurred.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    500,
				ErrorCode: "EKV_0000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
			},
		},
		"503 service unavailable error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "https://learn.akamai.com",
				"title": "Service Unavailable Error",
				"detail": "An internal error occurred.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 503,
				"errorCode": "EKV_0000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/initialize",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Service Unavailable Error",
				Detail:    "An internal error occurred.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    503,
				ErrorCode: "EKV_0000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
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
			result, err := client.GetEdgeKVInitializationStatus(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
