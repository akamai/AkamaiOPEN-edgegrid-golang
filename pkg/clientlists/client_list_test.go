package clientlists

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestClientList_GetClientLists(t *testing.T) {
	uri := "/client-list/v1/lists"

	tests := map[string]struct {
		params           GetClientListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetClientListsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetClientListsRequest{},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"content": [
					{
						"createDate": "2023-06-06T15:58:39.225+00:00",
						"createdBy": "ccare2",
						"deprecated": false,
						"filePrefix": "CL",
						"itemsCount": 1,
						"listId": "91596_AUDITLOGSTESTLIST",
						"listType": "CL",
						"name": "AUDIT LOGS - TEST LIST",
						"productionActivationStatus": "INACTIVE",
						"readOnly": false,
						"shared": false,
						"stagingActivationStatus": "INACTIVE",
						"tags": ["green"],
						"type": "IP",
						"updateDate": "2023-06-06T15:58:39.225+00:00",
						"updatedBy": "ccare2",
						"version": 1
					},
					{
						"createDate": "2022-11-10T14:42:04.857+00:00",
						"createdBy": "ccare2",
						"deprecated": false,
						"filePrefix": "CL",
						"itemsCount": 2,
						"listId": "85988_ANTHONYGEOLISTOPEN",
						"listType": "CL",
						"name": "AnthonyGeoListOPEN",
						"notes": "This is another Geo client list for Nov 11",
						"productionActivationStatus": "INACTIVE",
						"readOnly": false,
						"shared": false,
						"stagingActivationStatus": "INACTIVE",
						"tags": [],
						"type": "GEO",
						"updateDate": "2023-05-11T15:30:10.224+00:00",
						"updatedBy": "ccare2",
						"version": 66
					},
					{
						"createDate": "2022-10-17T13:39:25.319+00:00",
						"createdBy": "ccare2",
						"deprecated": false,
						"filePrefix": "CL",
						"itemsCount": 0,
						"listId": "85552_ANTHONYFILEHASHLIST",
						"listType": "CL",
						"name": "File Hash List",
						"notes": "This is another File hash client list for Oct 17",
						"productionActivationStatus": "PENDING_ACTIVATION",
						"readOnly": false,
						"shared": false,
						"stagingActivationStatus": "INACTIVE",
						"tags": ["blue"],
						"type": "TLS_FINGERPRINT",
						"updateDate": "2023-06-05T06:56:19.004+00:00",
						"updatedBy": "ccare2",
						"version": 343
					}
				]
			}
			`,
			expectedPath: uri,
			expectedResponse: &GetClientListsResponse{
				Content: []ListContent{
					{
						CreateDate:                 "2023-06-06T15:58:39.225+00:00",
						CreatedBy:                  "ccare2",
						Deprecated:                 false,
						ItemsCount:                 1,
						ListID:                     "91596_AUDITLOGSTESTLIST",
						ListType:                   "CL",
						Name:                       "AUDIT LOGS - TEST LIST",
						ProductionActivationStatus: "INACTIVE",
						ReadOnly:                   false,
						Shared:                     false,
						StagingActivationStatus:    "INACTIVE",
						Tags:                       []string{"green"},
						Type:                       "IP",
						UpdateDate:                 "2023-06-06T15:58:39.225+00:00",
						UpdatedBy:                  "ccare2",
						Version:                    1,
					},
					{
						CreateDate:                 "2022-11-10T14:42:04.857+00:00",
						CreatedBy:                  "ccare2",
						Deprecated:                 false,
						ItemsCount:                 2,
						ListID:                     "85988_ANTHONYGEOLISTOPEN",
						ListType:                   "CL",
						Name:                       "AnthonyGeoListOPEN",
						Notes:                      "This is another Geo client list for Nov 11",
						ProductionActivationStatus: "INACTIVE",
						ReadOnly:                   false,
						Shared:                     false,
						StagingActivationStatus:    "INACTIVE",
						Tags:                       []string{},
						Type:                       "GEO",
						UpdateDate:                 "2023-05-11T15:30:10.224+00:00",
						UpdatedBy:                  "ccare2",
						Version:                    66,
					},
					{
						CreateDate:                 "2022-10-17T13:39:25.319+00:00",
						CreatedBy:                  "ccare2",
						Deprecated:                 false,
						ItemsCount:                 0,
						ListID:                     "85552_ANTHONYFILEHASHLIST",
						ListType:                   "CL",
						Name:                       "File Hash List",
						Notes:                      "This is another File hash client list for Oct 17",
						ProductionActivationStatus: "PENDING_ACTIVATION",
						ReadOnly:                   false,
						Shared:                     false,
						StagingActivationStatus:    "INACTIVE",
						Tags:                       []string{"blue"},
						Type:                       "TLS_FINGERPRINT",
						UpdateDate:                 "2023-06-05T06:56:19.004+00:00",
						UpdatedBy:                  "ccare2",
						Version:                    343,
					},
				},
			},
		},
		"500 internal server error": {
			params:         GetClientListsRequest{},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
						"type": "internal_error",
						"title": "Internal Server Error",
						"detail": "Error fetching client lists",
						"status": 500
				}`,
			expectedPath: uri,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching client lists",
				StatusCode: http.StatusInternalServerError,
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
			result, err := client.GetClientLists(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
