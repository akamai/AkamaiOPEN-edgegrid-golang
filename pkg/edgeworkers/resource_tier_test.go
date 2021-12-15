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

func TestListResourceTiers(t *testing.T) {
	tests := map[string]struct {
		params           ListResourceTiersRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListResourceTiersResponse
		withError        error
	}{
		"200 OK": {
			params:         ListResourceTiersRequest{ContractID: "123"},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "resourceTiers": [
        {
            "resourceTierId": 100,
            "resourceTierName": "Basic Compute",
            "edgeWorkerLimits": [
                {
                    "limitName": "Maximum CPU time during initialization",
                    "limitValue": 30,
                    "limitUnit": "MILLISECOND"
                },
                {
                    "limitName": "Maximum response size for HTTP sub-requests during the responseProvider event handler",
                    "limitValue": 1048576,
                    "limitUnit": "BYTE"
                }
            ]
        },
        {
            "resourceTierId": 200,
            "resourceTierName": "Dynamic Compute",
            "edgeWorkerLimits": [
                {
                    "limitName": "Maximum number of HTTP sub-requests allowed in parallel for responseProvider",
                    "limitValue": 5,
                    "limitUnit": "COUNT"
                },
                {
                    "limitName": "Maximum wall time for HTTP sub-requests during the execution of the responseProvider event handler",
                    "limitValue": 1000,
                    "limitUnit": "MILLISECOND"
                }
            ]
        }
    ]
}`,
			expectedPath: "/edgeworkers/v1/resource-tiers?contractId=123",
			expectedResponse: &ListResourceTiersResponse{[]ResourceTier{
				{
					ID:   100,
					Name: "Basic Compute",
					EdgeWorkerLimits: []EdgeWorkerLimit{
						{
							LimitName:  "Maximum CPU time during initialization",
							LimitValue: 30,
							LimitUnit:  "MILLISECOND",
						},
						{
							LimitName:  "Maximum response size for HTTP sub-requests during the responseProvider event handler",
							LimitValue: 1048576,
							LimitUnit:  "BYTE",
						},
					},
				},
				{
					ID:   200,
					Name: "Dynamic Compute",
					EdgeWorkerLimits: []EdgeWorkerLimit{
						{
							LimitName:  "Maximum number of HTTP sub-requests allowed in parallel for responseProvider",
							LimitValue: 5,
							LimitUnit:  "COUNT",
						},
						{
							LimitName:  "Maximum wall time for HTTP sub-requests during the execution of the responseProvider event handler",
							LimitValue: 1000,
							LimitUnit:  "MILLISECOND",
						},
					},
				},
			}},
		},
		"500 internal server error": {
			params:         ListResourceTiersRequest{ContractID: "123"},
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
			expectedPath: "/edgeworkers/v1/resource-tiers?contractId=123",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
			},
		},
		"missing contract ID": {
			params:    ListResourceTiersRequest{},
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
			result, err := client.ListResourceTiers(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetResourceTier(t *testing.T) {
	tests := map[string]struct {
		params           GetResourceTierRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResourceTier
		withError        error
	}{
		"200 OK": {
			params:         GetResourceTierRequest{EdgeWorkerID: 123},
			responseStatus: http.StatusOK,
			responseBody: `
{
        
	"resourceTierId": 100,
	"resourceTierName": "Basic Compute",
	"edgeWorkerLimits": [
		{
			"limitName": "Maximum CPU time during initialization",
			"limitValue": 30,
			"limitUnit": "MILLISECOND"
		},
		{
			"limitName": "Maximum response size for HTTP sub-requests during the responseProvider event handler",
			"limitValue": 1048576,
			"limitUnit": "BYTE"
		}
	]
}
`,
			expectedPath: "/edgeworkers/v1/ids/123/resource-tier",
			expectedResponse: &ResourceTier{
				ID:   100,
				Name: "Basic Compute",
				EdgeWorkerLimits: []EdgeWorkerLimit{
					{
						LimitName:  "Maximum CPU time during initialization",
						LimitValue: 30,
						LimitUnit:  "MILLISECOND",
					},
					{
						LimitName:  "Maximum response size for HTTP sub-requests during the responseProvider event handler",
						LimitValue: 1048576,
						LimitUnit:  "BYTE",
					},
				},
			},
		},
		"500 internal server error": {
			params:         GetResourceTierRequest{EdgeWorkerID: 123},
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
			expectedPath: "/edgeworkers/v1/ids/123/resource-tier",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
			},
		},
		"missing edgeworker ID": {
			params:    GetResourceTierRequest{},
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
			result, err := client.GetResourceTier(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
