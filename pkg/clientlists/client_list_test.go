package clientlists

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetClientLists(t *testing.T) {
	uri := "/client-list/v1/lists"

	tests := map[string]struct {
		params           GetClientListsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetClientListsResponse
		withError        error
	}{
		"200 OK": {
			params:         GetClientListsRequest{},
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
				Content: []ClientList{
					{
						ListContent: ListContent{
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
					},
					{
						ListContent: ListContent{
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
					},
					{
						ListContent: ListContent{
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
		},
		"200 OK - Lists filtered by name and type": {
			params: GetClientListsRequest{
				Name: "list name",
				Type: []ClientListType{IP, GEO},
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
					}
				]
			}
			`,
			expectedPath: fmt.Sprintf(uri+"?name=%s&type=%s&type=%s", "list+name", "IP", "GEO"),
			expectedResponse: &GetClientListsResponse{
				Content: []ClientList{
					{
						ListContent: ListContent{
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
					},
				},
			},
		},
		"200 OK - Lists filtered by search and query params: includeItems, includeDeprecated, includeNetworkList, page, pageSize, sort": {
			params: GetClientListsRequest{
				Search:             "search term",
				IncludeItems:       true,
				IncludeDeprecated:  true,
				IncludeNetworkList: true,
				Page:               ptr.To(0),
				PageSize:           ptr.To(2),
				Sort:               []string{"updatedBy:desc", "value:desc"},
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
						"version": 1,
						"items": []
					}
				]
			}`,
			expectedPath: fmt.Sprintf(
				uri+"?includeDeprecated=%s&includeItems=%s&includeNetworkList=%s&page=%d&pageSize=%d&search=%s&sort=%s&sort=%s",
				"true", "true", "true", 0, 2, "search+term", "updatedBy%3Adesc", "value%3Adesc",
			),
			expectedResponse: &GetClientListsResponse{
				Content: []ClientList{
					{
						ListContent: ListContent{
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
						Items: []ListItemContent{},
					},
				},
			},
		},
		"500 internal server error": {
			params:         GetClientListsRequest{},
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
				context.Background(),
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

func TestGetClientList(t *testing.T) {
	uri := "/client-list/v1/lists/12_AB?includeItems=true"

	tests := map[string]struct {
		params           GetClientListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetClientListResponse
		withError        error
	}{
		"200 OK": {
			params: GetClientListRequest{
				ListID:       "12_AB",
				IncludeItems: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"createDate": "2023-06-06T15:58:39.225+00:00",
				"createdBy": "ccare2",
				"deprecated": false,
				"filePrefix": "CL",
				"itemsCount": 1,
				"listId": "12_AB",
				"listType": "CL",
				"name": "AUDIT LOGS - TEST LIST",
				"productionActivationStatus": "INACTIVE",
				"readOnly": false,
				"shared": false,
				"stagingActivationStatus": "INACTIVE",
				"productionActiveVersion": 2,
				"stagingActiveVersion": 2,
				"tags": ["green"],
				"type": "IP",
				"updateDate": "2023-06-06T15:58:39.225+00:00",
				"updatedBy": "ccare2",
				"version": 1,
				"groupId": 12,
				"groupName": "123_ABC",
				"contractId" :"12_CO",
				"items": [
					{
						"createDate": "2022-07-12T20:14:29.189+00:00",
						"createdBy": "ccare2",
						"createdVersion": 9,
						"productionStatus": "INACTIVE",
						"stagingStatus": "PENDING_ACTIVATION",
						"tags": [],
						"type": "IP",
						"updateDate": "2022-07-12T20:14:29.189+00:00",
						"updatedBy": "ccare2",
						"value": "7d0:1:0::0/64"
					},
					{
            "createDate": "2022-07-12T20:14:29.189+00:00",
            "createdBy": "ccare2",
            "createdVersion": 9,
            "description": "Item with description, tags, expiration date",
            "expirationDate": "2030-12-31T12:40:00.000+00:00",
            "productionStatus": "INACTIVE",
            "stagingStatus": "PENDING_ACTIVATION",
            "tags": [
                "red",
                "green",
                "blue"
            ],
            "type": "IP",
            "updateDate": "2022-07-12T20:14:29.189+00:00",
            "updatedBy": "ccare2",
            "value": "7d0:1:1::0/64"
        	}
				]
			}`,
			expectedPath: uri,
			expectedResponse: &GetClientListResponse{
				ListContent: ListContent{
					CreateDate:                 "2023-06-06T15:58:39.225+00:00",
					CreatedBy:                  "ccare2",
					Deprecated:                 false,
					ItemsCount:                 1,
					ListID:                     "12_AB",
					ListType:                   "CL",
					Name:                       "AUDIT LOGS - TEST LIST",
					ProductionActivationStatus: "INACTIVE",
					ReadOnly:                   false,
					Shared:                     false,
					StagingActivationStatus:    "INACTIVE",
					ProductionActiveVersion:    2,
					StagingActiveVersion:       2,
					Tags:                       []string{"green"},
					Type:                       "IP",
					UpdateDate:                 "2023-06-06T15:58:39.225+00:00",
					UpdatedBy:                  "ccare2",
					Version:                    1,
				},
				GroupID:    12,
				GroupName:  "123_ABC",
				ContractID: "12_CO",
				Items: []ListItemContent{
					{
						CreateDate:       "2022-07-12T20:14:29.189+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   9,
						ProductionStatus: "INACTIVE",
						StagingStatus:    "PENDING_ACTIVATION",
						Tags:             []string{},
						Type:             "IP",
						UpdateDate:       "2022-07-12T20:14:29.189+00:00",
						UpdatedBy:        "ccare2",
						Value:            "7d0:1:0::0/64",
					},
					{
						CreateDate:       "2022-07-12T20:14:29.189+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   9,
						ProductionStatus: "INACTIVE",
						StagingStatus:    "PENDING_ACTIVATION",
						Tags:             []string{"red", "green", "blue"},
						Description:      "Item with description, tags, expiration date",
						ExpirationDate:   "2030-12-31T12:40:00.000+00:00",
						Type:             "IP",
						UpdateDate:       "2022-07-12T20:14:29.189+00:00",
						UpdatedBy:        "ccare2",
						Value:            "7d0:1:1::0/64",
					},
				},
			},
		},
		"500 internal server error": {
			params: GetClientListRequest{
				ListID:       "12_AB",
				IncludeItems: true,
			},
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
		"validation error": {
			params:    GetClientListRequest{},
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
			result, err := client.GetClientList(
				session.ContextWithOptions(
					context.Background(),
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

func TestUpdateClientList(t *testing.T) {
	uri := "/client-list/v1/lists/12_12"
	request := UpdateClientListRequest{
		UpdateClientList: UpdateClientList{
			Name:  "Some New Name",
			Tags:  []string{"red"},
			Notes: "Updating list notes",
		},
		ListID: "12_12",
	}
	result := UpdateClientListResponse{
		ContractID: "M-2CF0QRI",
		GroupName:  "Kona QA16-M-2CF0QRI",
		GroupID:    12,
		ListContent: ListContent{
			CreateDate:                 "2023-04-03T15:50:34.074+00:00",
			CreatedBy:                  "ccare2",
			Deprecated:                 false,
			ItemsCount:                 51,
			ListID:                     "12_12",
			ListType:                   "CL",
			Name:                       "Some New Name",
			Tags:                       []string{"red"},
			Notes:                      "Updating list notes",
			ProductionActivationStatus: "INACTIVE",
			ReadOnly:                   false,
			Shared:                     false,
			StagingActivationStatus:    "INACTIVE",
			ProductionActiveVersion:    2,
			StagingActiveVersion:       2,
			Type:                       "IP",
			UpdateDate:                 "2023-06-15T20:28:09.047+00:00",
			UpdatedBy:                  "ccare2",
			Version:                    75,
		},
	}

	tests := map[string]struct {
		params              UpdateClientListRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *UpdateClientListResponse
		withError           error
	}{
		"200 OK": {
			params:              request,
			expectedRequestBody: `{"name":"Some New Name","notes":"Updating list notes","tags":["red"]}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"contractId": "M-2CF0QRI",
				"createDate": "2023-04-03T15:50:34.074+00:00",
				"createdBy": "ccare2",
				"deprecated": false,
				"filePrefix": "CL",
				"groupName": "Kona QA16-M-2CF0QRI",
				"groupId": 12,
				"itemsCount": 51,
				"listId": "12_12",
				"listType": "CL",
				"name": "Some New Name",
				"tags": [ "red"],
				"notes": "Updating list notes",
				"productionActivationStatus": "INACTIVE",
				"readOnly": false,
				"shared": false,
				"stagingActivationStatus": "INACTIVE",
				"productionActiveVersion":    2,
				"stagingActiveVersion":       2,
				"type": "IP",
				"updateDate": "2023-06-15T20:28:09.047+00:00",
				"updatedBy": "ccare2",
				"version": 75
			}`,
			expectedPath:     uri,
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         request,
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
		"validation error": {
			params:    UpdateClientListRequest{},
			withError: ErrStructValidation,
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

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateClientList(
				context.Background(),
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
func TestUpdateClientListItems(t *testing.T) {
	uri := "/client-list/v1/lists/12_12/items"
	request := UpdateClientListItemsRequest{
		ListID: "12_12",
		UpdateClientListItems: UpdateClientListItems{
			Append: []ListItemPayload{
				{
					Description:    "Lorem Ipsum has been the industry's standard dummy text ever since the 1500s...",
					ExpirationDate: "2026-12-26T01:32:08.375+00:00",
					Value:          "1.1.1.72",
				},
			},
			Update: []ListItemPayload{
				{
					Description:    "remove exp date and tags",
					ExpirationDate: "",
					Tags:           []string{"t"},
					Value:          "1.1.1.45",
				},
				{
					ExpirationDate: "2028-11-26T17:32:08.375+00:00",
					Value:          "1.1.1.33",
				},
			},
			Delete: []ListItemPayload{
				{
					Value: "1.1.1.38",
				},
			},
		},
	}
	result := UpdateClientListItemsResponse{
		Appended: []ListItemContent{
			{
				Description:      "Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley",
				ExpirationDate:   "2026-12-26T01:32:08.375+00:00",
				Tags:             []string{"new tag"},
				Value:            "1.1.1.75",
				CreateDate:       "2023-06-15T20:46:30.780+00:00",
				CreatedBy:        "ccare2",
				CreatedVersion:   76,
				ProductionStatus: "INACTIVE",
				StagingStatus:    "INACTIVE",
				Type:             "IP",
				UpdateDate:       "2023-06-15T20:46:30.780+00:00",
				UpdatedBy:        "ccare2",
			},
		},
		Deleted: []ListItemContent{
			{
				Value: "1.1.1.39",
			},
		},
		Updated: []ListItemContent{
			{
				Description:      "remove exp date and tags",
				Tags:             []string{"t1"},
				Value:            "1.1.1.45",
				CreateDate:       "2023-04-28T19:34:00.906+00:00",
				CreatedBy:        "ccare2",
				CreatedVersion:   54,
				ProductionStatus: "INACTIVE",
				StagingStatus:    "INACTIVE",
				Type:             "IP",
				UpdateDate:       "2023-06-15T20:46:30.765+00:00",
				UpdatedBy:        "ccare2",
			},
		},
	}

	tests := map[string]struct {
		params              UpdateClientListItemsRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *UpdateClientListItemsResponse
		withError           error
	}{
		"200 OK": {
			params:              request,
			expectedRequestBody: `{"append":[{"value":"1.1.1.72","tags":null,"description":"Lorem Ipsum has been the industry's standard dummy text ever since the 1500s...","expirationDate":"2026-12-26T01:32:08.375+00:00"}],"update":[{"value":"1.1.1.45","tags":["t"],"description":"remove exp date and tags","expirationDate":""},{"value":"1.1.1.33","tags":null,"description":"","expirationDate":"2028-11-26T17:32:08.375+00:00"}],"delete":[{"value":"1.1.1.38","tags":null,"description":"","expirationDate":""}]}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"appended": [
					{
						"createDate": "2023-06-15T20:46:30.780+00:00",
						"createdBy": "ccare2",
						"createdVersion": 76,
						"description": "Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley",
						"expirationDate": "2026-12-26T01:32:08.375+00:00",
						"productionStatus": "INACTIVE",
						"stagingStatus": "INACTIVE",
						"tags": [
							"new tag"
						],
						"type": "IP",
						"updateDate": "2023-06-15T20:46:30.780+00:00",
						"updatedBy": "ccare2",
						"value": "1.1.1.75"
					}
				],
				"deleted": [
					{
						"value": "1.1.1.39"
					}
				],
				"updated": [
					{
						"createDate": "2023-04-28T19:34:00.906+00:00",
						"createdBy": "ccare2",
						"createdVersion": 54,
						"description": "remove exp date and tags",
						"productionStatus": "INACTIVE",
						"stagingStatus": "INACTIVE",
						"tags": [
							"t1"
						],
						"type": "IP",
						"updateDate": "2023-06-15T20:46:30.765+00:00",
						"updatedBy": "ccare2",
						"value": "1.1.1.45"
					}
				]
			}`,
			expectedPath:     uri,
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         request,
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
		"validation error": {
			params:    UpdateClientListItemsRequest{},
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

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateClientListItems(
				context.Background(),
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

func TestCreateClientLists(t *testing.T) {
	uri := "/client-list/v1/lists"
	request := CreateClientListRequest{
		Name:       "TEST LIST",
		Type:       "IP",
		Notes:      "Some notes",
		Tags:       []string{"red", "green"},
		ContractID: "M-2CF0QRI",
		GroupID:    112524,
		Items: []ListItemPayload{
			{
				Value:          "1.1.1.1",
				Description:    "some description",
				Tags:           []string{},
				ExpirationDate: "2026-12-26T01:32:08.375+00:00",
			},
		},
	}
	result := CreateClientListResponse{
		ListContent: ListContent{
			ListID: "123_ABC",
			Name:   "TEST LIST",
			Type:   "IP",
			Notes:  "Some notes",
			Tags:   []string{"red", "green"},
		},
		ContractID: "M-2CF0QRI",
		GroupName:  "Group A",
		GroupID:    12,
		Items: []ListItemContent{
			{
				Value:          "1.1.1.1",
				Description:    "",
				Tags:           []string{},
				ExpirationDate: "2026-12-26T01:32:08.375+00:00",
			},
		},
	}

	tests := map[string]struct {
		params              CreateClientListRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateClientListResponse
		withError           error
	}{
		"201 Created": {
			params:              request,
			expectedRequestBody: `{"contractId":"M-2CF0QRI","groupId":112524,"name":"TEST LIST","type":"IP","notes":"Some notes","tags":["red","green"],"items":[{"value":"1.1.1.1","tags":[],"description":"some description","expirationDate":"2026-12-26T01:32:08.375+00:00"}]}`,
			responseStatus:      http.StatusCreated,
			responseBody: `{
				"listId": "123_ABC",
				"name": "TEST LIST",
				"type": "IP",
				"notes": "Some notes",
				"tags": [
					"red",
					"green"
				],
				"contractId": "M-2CF0QRI",
				"groupName": "Group A",
				"groupId": 12,
				"items": [
					{
						"value": "1.1.1.1",
						"description": "",
						"tags": [],
						"expirationDate": "2026-12-26T01:32:08.375+00:00"
					}
				]
			}
			`,
			expectedPath:     uri,
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         request,
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
		"validation error": {
			params:    CreateClientListRequest{},
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

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateClientList(
				context.Background(),
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

func TestDeleteClientLists(t *testing.T) {
	uri := "/client-list/v1/lists/12_AB"
	request := DeleteClientListRequest{
		ListID: "12_AB",
	}

	tests := map[string]struct {
		params           DeleteClientListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Error
		withError        error
	}{
		"204 NoContent": {
			params:           request,
			responseBody:     "",
			responseStatus:   http.StatusNoContent,
			expectedPath:     uri,
			expectedResponse: nil,
		},
		"500 internal server error": {
			params:         request,
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
		"validation error": {
			params:    DeleteClientListRequest{},
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
			err := client.DeleteClientList(
				context.Background(),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
