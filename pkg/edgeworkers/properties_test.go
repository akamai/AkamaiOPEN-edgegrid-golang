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

func TestListProperties(t *testing.T) {
	tests := map[string]struct {
		params           ListPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListPropertiesResponse
		withError        error
	}{
		"200 OK - no query": {
			params:         ListPropertiesRequest{EdgeWorkerID: 123},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"properties": [
		{
			"propertyId": 100,
			"propertyName": "property1",
			"stagingVersion": null,
			"productionVersion": 2,
			"latestVersion": 2
		},
		{
			"propertyId": 101,
			"propertyName": "property2",
			"stagingVersion": 1,
			"productionVersion": null,
			"latestVersion": 1
		},
		{
			"propertyId": 102,
			"propertyName": "property3",
			"stagingVersion": null,
			"productionVersion": null,
			"latestVersion": 3
		}
	],
	"limitedAccessToProperties": true
}`,
			expectedPath: "/edgeworkers/v1/ids/123/properties?activeOnly=false",
			expectedResponse: &ListPropertiesResponse{
				Properties: []Property{
					{
						ID:                100,
						Name:              "property1",
						StagingVersion:    0,
						ProductionVersion: 2,
						LatestVersion:     2,
					},
					{
						ID:                101,
						Name:              "property2",
						StagingVersion:    1,
						ProductionVersion: 0,
						LatestVersion:     1,
					},
					{
						ID:                102,
						Name:              "property3",
						StagingVersion:    0,
						ProductionVersion: 0,
						LatestVersion:     3,
					},
				},
				LimitedAccessToProperties: true,
			},
		},
		"200 OK - with query": {
			params:         ListPropertiesRequest{EdgeWorkerID: 123, ActiveOnly: true},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"properties": [
		{
			"propertyId": 100,
			"propertyName": "property1",
			"stagingVersion": null,
			"productionVersion": 2,
			"latestVersion": 2
		},
		{
			"propertyId": 101,
			"propertyName": "property2",
			"stagingVersion": 1,
			"productionVersion": null,
			"latestVersion": 1
		}
	],
	"limitedAccessToProperties": false
}`,
			expectedPath: "/edgeworkers/v1/ids/123/properties?activeOnly=true",
			expectedResponse: &ListPropertiesResponse{
				Properties: []Property{
					{
						ID:                100,
						Name:              "property1",
						StagingVersion:    0,
						ProductionVersion: 2,
						LatestVersion:     2,
					},
					{
						ID:                101,
						Name:              "property2",
						StagingVersion:    1,
						ProductionVersion: 0,
						LatestVersion:     1,
					},
				},
				LimitedAccessToProperties: false,
			},
		},
		"500 internal server error": {
			params:         ListPropertiesRequest{EdgeWorkerID: 123},
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
			expectedPath: "/edgeworkers/v1/ids/123/properties?activeOnly=false",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-server-error",
				Title:     "An unexpected error has occurred.",
				Detail:    "Error processing request",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    500,
				ErrorCode: "EW4303",
			},
		},
		"404 resource not found": {
			params:         ListPropertiesRequest{EdgeWorkerID: 123},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
	"type": "/edgeworkers/error-types/edgeworkers-not-found",
	"title": "The given resource could not be found.",
	"detail": "Unable to find the requested EdgeWorker ID",
	"instance": "/edgeworkers/error-instances/abc",
	"status": 404,
	"errorCode": "EW2002"
}`,
			expectedPath: "/edgeworkers/v1/ids/123/properties?activeOnly=false",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-not-found",
				Title:     "The given resource could not be found.",
				Detail:    "Unable to find the requested EdgeWorker ID",
				Instance:  "/edgeworkers/error-instances/abc",
				Status:    404,
				ErrorCode: "EW2002",
			},
		},
		"missing edgeworker ID": {
			params:    ListPropertiesRequest{},
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
			result, err := client.ListProperties(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
