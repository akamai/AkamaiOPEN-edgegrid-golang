package clientlists

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	internalServerErrBody = `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error fetching client lists",
					"status": 500
				}`

	internalServerErr = &Error{
		Type:       "internal_error",
		Title:      "Internal Server Error",
		Detail:     "Error fetching client lists",
		StatusCode: http.StatusInternalServerError,
	}
)

func getMockTestServer(t *testing.T, method, expectedPath string, responseStatus int, responseBody, expectedRequestBody string) *httptest.Server {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedPath, r.URL.String())
		assert.Equal(t, method, r.Method)
		w.WriteHeader(responseStatus)
		_, err := w.Write([]byte(responseBody))
		assert.NoError(t, err)
		if len(expectedRequestBody) > 0 {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Equal(t, expectedRequestBody, string(body))
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func checkResponse[T any](t *testing.T, result *T, err error, expectedResponse *T, withError error) {
	t.Helper()
	if withError != nil {
		assert.True(t, errors.Is(err, withError), "want: %s; got: %s", withError, err)
		return
	}
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

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
		"200 OK - Lists filtered by RequestHeaderNameValue type": {
			params: GetClientListsRequest{
				Type: []ClientListType{RequestHeaderNameValue},
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
						"listId": "91596_REQHEADERLIST",
						"listType": "CL",
						"name": "Request Header List",
						"productionActivationStatus": "INACTIVE",
						"readOnly": false,
						"shared": false,
						"stagingActivationStatus": "INACTIVE",
						"tags": [],
						"type": "REQUEST_HEADER_NAME_VALUE",
						"updateDate": "2023-06-06T15:58:39.225+00:00",
						"updatedBy": "ccare2",
						"version": 1
					}
				]
			}
			`,
			expectedPath: fmt.Sprintf(uri+"?type=%s", "REQUEST_HEADER_NAME_VALUE"),
			expectedResponse: &GetClientListsResponse{
				Content: []ClientList{
					{
						ListContent: ListContent{
							CreateDate:                 "2023-06-06T15:58:39.225+00:00",
							CreatedBy:                  "ccare2",
							Deprecated:                 false,
							ItemsCount:                 1,
							ListID:                     "91596_REQHEADERLIST",
							ListType:                   "CL",
							Name:                       "Request Header List",
							ProductionActivationStatus: "INACTIVE",
							ReadOnly:                   false,
							Shared:                     false,
							StagingActivationStatus:    "INACTIVE",
							Tags:                       []string{},
							Type:                       "REQUEST_HEADER_NAME_VALUE",
							UpdateDate:                 "2023-06-06T15:58:39.225+00:00",
							UpdatedBy:                  "ccare2",
							Version:                    1,
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
			responseBody:   internalServerErrBody,
			expectedPath:   uri,
			withError:      internalServerErr,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := getMockTestServer(t, http.MethodGet, test.expectedPath, test.responseStatus, test.responseBody, "")
			client := mockAPIClient(t, mockServer)
			result, err := client.GetClientLists(context.Background(), test.params)
			checkResponse(t, result, err, test.expectedResponse, test.withError)
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
		"200 OK - RequestHeaderNameValue items with key and values": {
			params: GetClientListRequest{
				ListID:       "12_RH",
				IncludeItems: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"createDate": "2023-06-06T15:58:39.225+00:00",
				"createdBy": "ccare2",
				"deprecated": false,
				"filePrefix": "CL",
				"itemsCount": 2,
				"listId": "12_RH",
				"listType": "CL",
				"name": "Request Header List",
				"productionActivationStatus": "INACTIVE",
				"readOnly": false,
				"shared": false,
				"stagingActivationStatus": "INACTIVE",
				"tags": [],
				"type": "REQUEST_HEADER_NAME_VALUE",
				"updateDate": "2023-06-06T15:58:39.225+00:00",
				"updatedBy": "ccare2",
				"version": 3,
				"groupId": 12,
				"groupName": "Group A",
				"contractId": "12_CO",
				"items": [
					{
						"key": "X-Custom-Header",
						"values": ["value1", "value2"],
						"createDate": "2023-01-01T00:00:00.000+00:00",
						"createdBy": "ccare2",
						"createdVersion": 1,
						"productionStatus": "INACTIVE",
						"stagingStatus": "INACTIVE",
						"tags": [],
						"type": "REQUEST_HEADER_NAME_VALUE",
						"updateDate": "2023-01-01T00:00:00.000+00:00",
						"updatedBy": "ccare2"
					},
					{
						"key": "Accept-Language",
						"values": ["en-US"],
						"description": "Language header match",
						"createDate": "2023-02-01T00:00:00.000+00:00",
						"createdBy": "ccare2",
						"createdVersion": 2,
						"productionStatus": "INACTIVE",
						"stagingStatus": "INACTIVE",
						"tags": ["lang"],
						"type": "REQUEST_HEADER_NAME_VALUE",
						"updateDate": "2023-02-01T00:00:00.000+00:00",
						"updatedBy": "ccare2"
					}
				]
			}`,
			expectedPath: "/client-list/v1/lists/12_RH?includeItems=true",
			expectedResponse: &GetClientListResponse{
				ListContent: ListContent{
					CreateDate:                 "2023-06-06T15:58:39.225+00:00",
					CreatedBy:                  "ccare2",
					Deprecated:                 false,
					ItemsCount:                 2,
					ListID:                     "12_RH",
					ListType:                   "CL",
					Name:                       "Request Header List",
					ProductionActivationStatus: "INACTIVE",
					ReadOnly:                   false,
					Shared:                     false,
					StagingActivationStatus:    "INACTIVE",
					Tags:                       []string{},
					Type:                       "REQUEST_HEADER_NAME_VALUE",
					UpdateDate:                 "2023-06-06T15:58:39.225+00:00",
					UpdatedBy:                  "ccare2",
					Version:                    3,
				},
				GroupID:    12,
				GroupName:  "Group A",
				ContractID: "12_CO",
				Items: []ListItemContent{
					{
						Key:              "X-Custom-Header",
						Values:           []string{"value1", "value2"},
						CreateDate:       "2023-01-01T00:00:00.000+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   1,
						ProductionStatus: "INACTIVE",
						StagingStatus:    "INACTIVE",
						Tags:             []string{},
						Type:             "REQUEST_HEADER_NAME_VALUE",
						UpdateDate:       "2023-01-01T00:00:00.000+00:00",
						UpdatedBy:        "ccare2",
					},
					{
						Key:              "Accept-Language",
						Values:           []string{"en-US"},
						Description:      "Language header match",
						CreateDate:       "2023-02-01T00:00:00.000+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   2,
						ProductionStatus: "INACTIVE",
						StagingStatus:    "INACTIVE",
						Tags:             []string{"lang"},
						Type:             "REQUEST_HEADER_NAME_VALUE",
						UpdateDate:       "2023-02-01T00:00:00.000+00:00",
						UpdatedBy:        "ccare2",
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
			responseBody:   internalServerErrBody,
			expectedPath:   uri,
			withError:      internalServerErr,
		},
		"validation error": {
			params:    GetClientListRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := getMockTestServer(t, http.MethodGet, test.expectedPath, test.responseStatus, test.responseBody, "")
			client := mockAPIClient(t, mockServer)
			result, err := client.GetClientList(session.ContextWithOptions(context.Background()), test.params)
			checkResponse(t, result, err, test.expectedResponse, test.withError)
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
			responseBody:   internalServerErrBody,
			expectedPath:   uri,
			withError:      internalServerErr,
		},
		"validation error": {
			params:    UpdateClientListRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := getMockTestServer(t, http.MethodPut, test.expectedPath, test.responseStatus, test.responseBody, test.expectedRequestBody)
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateClientList(context.Background(), test.params)
			checkResponse(t, result, err, test.expectedResponse, test.withError)
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
		"200 OK - Update RequestHeaderNameValue items": {
			params: UpdateClientListItemsRequest{
				ListID: "12_RH",
				UpdateClientListItems: UpdateClientListItems{
					Append: []ListItemPayload{
						{
							Key:    "X-New-Header",
							Values: []string{"newval"},
							Tags:   []string{},
						},
					},
					Update: []ListItemPayload{
						{
							Key:    "X-Custom-Header",
							Values: []string{"updated1", "updated2"},
							Tags:   []string{"t"},
						},
					},
					Delete: []ListItemPayload{
						{
							Key: "X-Old-Header",
						},
					},
				},
			},
			expectedRequestBody: `{"append":[{"key":"X-New-Header","values":["newval"],"tags":[],"description":"","expirationDate":""}],"update":[{"key":"X-Custom-Header","values":["updated1","updated2"],"tags":["t"],"description":"","expirationDate":""}],"delete":[{"key":"X-Old-Header","tags":null,"description":"","expirationDate":""}]}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"appended": [
					{
						"key": "X-New-Header",
						"values": ["newval"],
						"tags": [],
						"type": "REQUEST_HEADER_NAME_VALUE",
						"productionStatus": "INACTIVE",
						"stagingStatus": "INACTIVE",
						"createDate": "2023-06-15T20:46:30.780+00:00",
						"createdBy": "ccare2",
						"createdVersion": 5,
						"updateDate": "2023-06-15T20:46:30.780+00:00",
						"updatedBy": "ccare2"
					}
				],
				"deleted": [
					{
						"key": "X-Old-Header"
					}
				],
				"updated": [
					{
						"key": "X-Custom-Header",
						"values": ["updated1", "updated2"],
						"tags": ["t"],
						"type": "REQUEST_HEADER_NAME_VALUE",
						"productionStatus": "INACTIVE",
						"stagingStatus": "INACTIVE",
						"createDate": "2023-04-28T19:34:00.906+00:00",
						"createdBy": "ccare2",
						"createdVersion": 3,
						"updateDate": "2023-06-15T20:46:30.765+00:00",
						"updatedBy": "ccare2"
					}
				]
			}`,
			expectedPath: "/client-list/v1/lists/12_RH/items",
			expectedResponse: &UpdateClientListItemsResponse{
				Appended: []ListItemContent{
					{
						Key:              "X-New-Header",
						Values:           []string{"newval"},
						Tags:             []string{},
						Type:             "REQUEST_HEADER_NAME_VALUE",
						ProductionStatus: "INACTIVE",
						StagingStatus:    "INACTIVE",
						CreateDate:       "2023-06-15T20:46:30.780+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   5,
						UpdateDate:       "2023-06-15T20:46:30.780+00:00",
						UpdatedBy:        "ccare2",
					},
				},
				Deleted: []ListItemContent{
					{Key: "X-Old-Header"},
				},
				Updated: []ListItemContent{
					{
						Key:              "X-Custom-Header",
						Values:           []string{"updated1", "updated2"},
						Tags:             []string{"t"},
						Type:             "REQUEST_HEADER_NAME_VALUE",
						ProductionStatus: "INACTIVE",
						StagingStatus:    "INACTIVE",
						CreateDate:       "2023-04-28T19:34:00.906+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   3,
						UpdateDate:       "2023-06-15T20:46:30.765+00:00",
						UpdatedBy:        "ccare2",
					},
				},
			},
		},
		"500 internal server error": {
			params:         request,
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerErrBody,
			expectedPath:   uri,
			withError:      internalServerErr,
		},
		"validation error": {
			params:    UpdateClientListItemsRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := getMockTestServer(t, http.MethodPost, test.expectedPath, test.responseStatus, test.responseBody, test.expectedRequestBody)
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateClientListItems(context.Background(), test.params)
			checkResponse(t, result, err, test.expectedResponse, test.withError)
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
		"201 Created - RequestHeaderNameValue type": {
			params: CreateClientListRequest{
				Name:       "Request Header List",
				Type:       RequestHeaderNameValue,
				Notes:      "Some notes",
				Tags:       []string{"tag1"},
				ContractID: "M-2CF0QRI",
				GroupID:    112524,
				Items: []ListItemPayload{
					{
						Key:    "X-Custom-Header",
						Values: []string{"value1", "value2"},
						Tags:   []string{},
					},
				},
			},
			expectedRequestBody: `{"contractId":"M-2CF0QRI","groupId":112524,"name":"Request Header List","type":"REQUEST_HEADER_NAME_VALUE","notes":"Some notes","tags":["tag1"],"items":[{"key":"X-Custom-Header","values":["value1","value2"],"tags":[],"description":"","expirationDate":""}]}`,
			responseStatus:      http.StatusCreated,
			responseBody: `{
				"listId": "456_RH",
				"name": "Request Header List",
				"type": "REQUEST_HEADER_NAME_VALUE",
				"notes": "Some notes",
				"tags": ["tag1"],
				"contractId": "M-2CF0QRI",
				"groupName": "Group A",
				"groupId": 12,
				"items": [
					{
						"key": "X-Custom-Header",
						"values": ["value1", "value2"],
						"tags": [],
						"description": "",
						"expirationDate": ""
					}
				]
			}`,
			expectedPath: uri,
			expectedResponse: &CreateClientListResponse{
				ListContent: ListContent{
					ListID: "456_RH",
					Name:   "Request Header List",
					Type:   "REQUEST_HEADER_NAME_VALUE",
					Notes:  "Some notes",
					Tags:   []string{"tag1"},
				},
				ContractID: "M-2CF0QRI",
				GroupName:  "Group A",
				GroupID:    12,
				Items: []ListItemContent{
					{
						Key:    "X-Custom-Header",
						Values: []string{"value1", "value2"},
						Tags:   []string{},
					},
				},
			},
		},
		"500 internal server error": {
			params:         request,
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerErrBody,
			expectedPath:   uri,
			withError:      internalServerErr,
		},
		"validation error": {
			params:    CreateClientListRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := getMockTestServer(t, http.MethodPost, test.expectedPath, test.responseStatus, test.responseBody, test.expectedRequestBody)
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateClientList(context.Background(), test.params)
			checkResponse(t, result, err, test.expectedResponse, test.withError)
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
			responseBody:   internalServerErrBody,
			expectedPath:   uri,
			withError:      internalServerErr,
		},
		"validation error": {
			params:    DeleteClientListRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := getMockTestServer(t, http.MethodDelete, test.expectedPath, test.responseStatus, test.responseBody, "")
			client := mockAPIClient(t, mockServer)
			err := client.DeleteClientList(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestTranslateUsernames(t *testing.T) {
	uri := "/appsec/v1/search/user/external-uuid"
	request := TranslateUsernamesRequest{
		"user1",
		"user2",
		"user3",
	}
	result := TranslateUsernamesResponse{
		"user1": "3a453537-faa8-4525-b5db-022447bbbf2a",
		"user2": "07e29045-7739-4bd9-8cfb-9f118e000337",
		"user3": "e164394a-5ae1-4208-8487-1ac0f368ecf3",
	}

	tests := map[string]struct {
		params              TranslateUsernamesRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *TranslateUsernamesResponse
		withError           error
	}{
		"200 - translate usernames": {
			params:              request,
			expectedRequestBody: `["user1","user2","user3"]`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"user1": "3a453537-faa8-4525-b5db-022447bbbf2a",
				"user2": "07e29045-7739-4bd9-8cfb-9f118e000337",
				"user3": "e164394a-5ae1-4208-8487-1ac0f368ecf3"
			}`,
			expectedPath:     uri,
			expectedResponse: &result,
		},
		"500 internal server error": {
			params:         request,
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerErrBody,
			expectedPath:   uri,
			withError:      internalServerErr,
		},
		"validation error": {
			params:    TranslateUsernamesRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := getMockTestServer(t, http.MethodPost, test.expectedPath, test.responseStatus, test.responseBody, test.expectedRequestBody)
			client := mockAPIClient(t, mockServer)
			result, err := client.TranslateUsernames(context.Background(), test.params)
			checkResponse(t, result, err, test.expectedResponse, test.withError)
		})
	}
}

func TestGetClientListItems(t *testing.T) {
	uri := "/client-list/v1/lists/12_AB/items?showUsernames=true"

	tests := map[string]struct {
		params           GetClientListItemsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetClientListItemsResponse
		withError        error
	}{
		"200 OK - non user type list": {
			params: GetClientListItemsRequest{
				ListID: "12_AB",
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"content": [
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
			expectedResponse: &GetClientListItemsResponse{
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
		"200 OK - user type list - show username enabled": {
			params: GetClientListItemsRequest{
				ListID: "12_AB",
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"content": [
					{
						"createDate": "2022-07-12T20:14:29.189+00:00",
						"createdBy": "ccare2",
						"createdVersion": 9,
						"productionStatus": "INACTIVE",
						"stagingStatus": "PENDING_ACTIVATION",
						"tags": [],
						"type": "USER_ID",
						"updateDate": "2022-07-12T20:14:29.189+00:00",
						"updatedBy": "ccare2",
						"value": "3a453537-faa8-4525-b5db-022447bbbf2a",
						"username": "user1"
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
						"type": "USER_ID",
						"updateDate": "2022-07-12T20:14:29.189+00:00",
						"updatedBy": "ccare2",
						"value": "07e29045-7739-4bd9-8cfb-9f118e000337",
						"username": "user2"
					}
				]
			}`,
			expectedPath: uri,
			expectedResponse: &GetClientListItemsResponse{
				Items: []ListItemContent{
					{
						CreateDate:       "2022-07-12T20:14:29.189+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   9,
						ProductionStatus: "INACTIVE",
						StagingStatus:    "PENDING_ACTIVATION",
						Tags:             []string{},
						Type:             "USER_ID",
						UpdateDate:       "2022-07-12T20:14:29.189+00:00",
						UpdatedBy:        "ccare2",
						Value:            "3a453537-faa8-4525-b5db-022447bbbf2a",
						Username:         "user1",
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
						Type:             "USER_ID",
						UpdateDate:       "2022-07-12T20:14:29.189+00:00",
						UpdatedBy:        "ccare2",
						Value:            "07e29045-7739-4bd9-8cfb-9f118e000337",
						Username:         "user2",
					},
				},
			},
		},
		"200 OK - RequestHeaderNameValue items with key and values": {
			params: GetClientListItemsRequest{
				ListID: "12_RH",
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"content": [
					{
						"key": "X-Custom-Header",
						"values": ["value1", "value2"],
						"createDate": "2023-01-01T00:00:00.000+00:00",
						"createdBy": "ccare2",
						"createdVersion": 1,
						"productionStatus": "INACTIVE",
						"stagingStatus": "INACTIVE",
						"tags": [],
						"type": "REQUEST_HEADER_NAME_VALUE",
						"updateDate": "2023-01-01T00:00:00.000+00:00",
						"updatedBy": "ccare2"
					},
					{
						"key": "Accept-Language",
						"values": ["en-US"],
						"description": "Language header match",
						"createDate": "2023-02-01T00:00:00.000+00:00",
						"createdBy": "ccare2",
						"createdVersion": 2,
						"productionStatus": "ACTIVE",
						"stagingStatus": "ACTIVE",
						"tags": ["lang"],
						"type": "REQUEST_HEADER_NAME_VALUE",
						"updateDate": "2023-02-01T00:00:00.000+00:00",
						"updatedBy": "ccare2"
					}
				]
			}`,
			expectedPath: "/client-list/v1/lists/12_RH/items?showUsernames=true",
			expectedResponse: &GetClientListItemsResponse{
				Items: []ListItemContent{
					{
						Key:              "X-Custom-Header",
						Values:           []string{"value1", "value2"},
						CreateDate:       "2023-01-01T00:00:00.000+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   1,
						ProductionStatus: "INACTIVE",
						StagingStatus:    "INACTIVE",
						Tags:             []string{},
						Type:             "REQUEST_HEADER_NAME_VALUE",
						UpdateDate:       "2023-01-01T00:00:00.000+00:00",
						UpdatedBy:        "ccare2",
					},
					{
						Key:              "Accept-Language",
						Values:           []string{"en-US"},
						Description:      "Language header match",
						CreateDate:       "2023-02-01T00:00:00.000+00:00",
						CreatedBy:        "ccare2",
						CreatedVersion:   2,
						ProductionStatus: "ACTIVE",
						StagingStatus:    "ACTIVE",
						Tags:             []string{"lang"},
						Type:             "REQUEST_HEADER_NAME_VALUE",
						UpdateDate:       "2023-02-01T00:00:00.000+00:00",
						UpdatedBy:        "ccare2",
					},
				},
			},
		},
		"500 internal server error": {
			params: GetClientListItemsRequest{
				ListID: "12_AB",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerErrBody,
			expectedPath:   uri,
			withError:      internalServerErr,
		},
		"validation error": {
			params:    GetClientListItemsRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := getMockTestServer(t, http.MethodGet, test.expectedPath, test.responseStatus, test.responseBody, "")
			client := mockAPIClient(t, mockServer)
			result, err := client.GetClientListItems(session.ContextWithOptions(context.Background()), test.params)
			checkResponse(t, result, err, test.expectedResponse, test.withError)
		})
	}
}
