package cloudwrapper

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestListCapacity(t *testing.T) {
	tests := map[string]struct {
		request          ListCapacitiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListCapacitiesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: ListCapacitiesRequest{
				ContractIDs: nil,
			},
			responseStatus: 200,
			expectedPath:   "/cloud-wrapper/v1/capacity",
			responseBody: `
			{
				"capacities": [
					{
						"locationId": 1,
						"locationName": "US East",
						"contractId": "A-BCDEFG",
						"type": "MEDIA",
						"approvedCapacity": {
							"value": 2000,
							"unit": "GB"
						},
						"assignedCapacity": {
							"value": 2,
							"unit": "GB"
						},
						"unassignedCapacity": {
							"value": 1998,
							"unit": "GB"
						}
					}
				]
			}`,
			expectedResponse: &ListCapacitiesResponse{
				Capacities: []LocationCapacity{
					{
						LocationID:   1,
						LocationName: "US East",
						ContractID:   "A-BCDEFG",
						Type:         Media,
						ApprovedCapacity: Capacity{
							Value: 2000,
							Unit:  "GB",
						},
						AssignedCapacity: Capacity{
							Value: 2,
							Unit:  "GB",
						},
						UnassignedCapacity: Capacity{
							Value: 1998,
							Unit:  "GB",
						},
					},
				},
			},
		},
		"200 OK with contracts": {
			request: ListCapacitiesRequest{
				ContractIDs: []string{"A-BCDEF", "B-CDEFG"},
			},
			responseStatus: 200,
			expectedPath:   "/cloud-wrapper/v1/capacity?contractIds=A-BCDEF&contractIds=B-CDEFG",
			responseBody: `
			{
				"capacities": [
					{
						"locationId": 1,
						"locationName": "US East",
						"contractId": "A-BCDEFG",
						"type": "WEB_ENHANCED_TLS",
						"approvedCapacity": {
							"value": 10,
							"unit": "TB"
						},
						"assignedCapacity": {
							"value": 1,
							"unit": "TB"
						},
						"unassignedCapacity": {
							"value": 9,
							"unit": "TB"
						}
					}
				]
			}`,
			expectedResponse: &ListCapacitiesResponse{
				Capacities: []LocationCapacity{
					{
						LocationID:   1,
						LocationName: "US East",
						ContractID:   "A-BCDEFG",
						Type:         WebEnhancedTLS,
						ApprovedCapacity: Capacity{
							Value: 10,
							Unit:  UnitTB,
						},
						AssignedCapacity: Capacity{
							Value: 1,
							Unit:  UnitTB,
						},
						UnassignedCapacity: Capacity{
							Value: 9,
							Unit:  UnitTB,
						},
					},
				},
			},
		},
		"401 not authorized": {
			request:        ListCapacitiesRequest{},
			responseStatus: 401,
			responseBody: `{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
    "title": "Not authorized",
    "status": 401,
    "detail": "The signature does not match",
    "instance": "https://instance.luna-dev.akamaiapis.net/cloud-wrapper/v1/capacity",
    "method": "GET",
    "serverIp": "2.2.2.2",
    "clientIp": "3.3.3.3",
    "requestId": "a7a7a7a7a7a",
    "requestTime": "2023-05-22T10:05:22Z"
}`,
			expectedPath:     "/cloud-wrapper/v1/capacity",
			expectedResponse: nil,
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
					Title:       "Not authorized",
					Instance:    "https://instance.luna-dev.akamaiapis.net/cloud-wrapper/v1/capacity",
					Status:      401,
					Detail:      "The signature does not match",
					Method:      "GET",
					ServerIP:    "2.2.2.2",
					ClientIP:    "3.3.3.3",
					RequestID:   "a7a7a7a7a7a",
					RequestTime: "2023-05-22T10:05:22Z",
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
		"500": {
			request:          ListCapacitiesRequest{},
			responseStatus:   500,
			expectedPath:     "/cloud-wrapper/v1/capacity",
			expectedResponse: nil,
			responseBody: `{
				"type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				"title": "Server Error",
				"status": 500,
				"instance": "https://instance.luna-dev.akamaiapis.net/cloud-wrapper/v1/capacity",
				"method": "GET",
				"serverIp": "2.2.2.2",
				"clientIp": "3.3.3.3",
				"requestId": "a7a7a7a7a7a",
				"requestTime": "2021-12-06T10:27:11Z"
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
					Title:       "Server Error",
					Status:      500,
					Instance:    "https://instance.luna-dev.akamaiapis.net/cloud-wrapper/v1/capacity",
					Method:      "GET",
					ServerIP:    "2.2.2.2",
					ClientIP:    "3.3.3.3",
					RequestID:   "a7a7a7a7a7a",
					RequestTime: "2021-12-06T10:27:11Z",
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
			result, err := client.ListCapacities(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
