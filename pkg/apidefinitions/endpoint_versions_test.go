package apidefinitions

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEndpointVersions(t *testing.T) {
	tests := map[string]struct {
		params           ListEndpointVersionsRequest
		expectedPath     string
		expectedResponse *ListEndpointVersionsResponse
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListEndpointVersionsRequest{
				APIEndpointID: 2,
			},
			expectedPath: "/api-definitions/v2/endpoints/2/versions",
			expectedResponse: &ListEndpointVersionsResponse{
				TotalSize:       1,
				Page:            2,
				PageSize:        3,
				APIEndpointID:   2,
				APIEndpointName: "test_cache2",
				APIVersions: []APIVersion{
					{
						CreateDate:           "2022-08-11T06:33:29+0000",
						CreatedBy:            "tester",
						UpdateDate:           "2022-08-16T10:08:41+0000",
						UpdatedBy:            "tester",
						APIEndpointVersionID: 10101010,
						BasePath:             "/test",
						VersionNumber:        1,
						Description:          nil,
						BasedOn:              nil,
						StagingStatus:        nil,
						ProductionStatus:     nil,
						StagingDate:          nil,
						ProductionDate:       nil,
						IsVersionLocked:      false,
						Hidden:               false,
						AvailableActions: []string{
							"ACTIVATE_ON_STAGING",
							"COMPARE_ENDPOINT",
							"DELETE",
							"CLONE_VERSION",
							"EDIT_AAG_SETTINGS",
							"ACTIVATE_ON_PRODUCTION",
							"COMPARE_RESOURCE_PURPOSES",
							"HIDE_VERSION",
							"RESOURCES",
							"EDIT_ENDPOINT_DEFINITION",
							"VIEW_AAG_SETTINGS",
							"COMPARE_AAG_SETTINGS",
							"VIEW_TAPIOCA"},
						CloningStatus: nil,
						LockVersion:   1,
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"totalSize": 1,
	"page": 2,
	"pageSize": 3,
    "apiEndPointId": 2,
    "apiEndPointName": "test_cache2",
    "apiVersions": [
        {
            "createdBy": "tester",
            "createDate": "2022-08-11T06:33:29+0000",
            "updateDate": "2022-08-16T10:08:41+0000",
            "updatedBy": "tester",
            "apiEndPointVersionId": 10101010,
            "basePath": "/test",
            "versionNumber": 1,
            "description": null,
            "basedOn": null,
            "stagingStatus": null,
            "productionStatus": null,
            "stagingDate": null,
            "productionDate": null,
            "isVersionLocked": false,
            "hidden": false,
            "availableActions": [
                "ACTIVATE_ON_STAGING",
                "COMPARE_ENDPOINT",
                "DELETE",
                "CLONE_VERSION",
                "EDIT_AAG_SETTINGS",
                "ACTIVATE_ON_PRODUCTION",
                "COMPARE_RESOURCE_PURPOSES",
                "HIDE_VERSION",
                "RESOURCES",
                "EDIT_ENDPOINT_DEFINITION",
                "VIEW_AAG_SETTINGS",
                "COMPARE_AAG_SETTINGS",
                "VIEW_TAPIOCA"
            ],
            "cloningStatus": null,
            "lockVersion": 1
        }
    ]
}
`,
		}, "200 OK fields populated": {
			params: ListEndpointVersionsRequest{
				APIEndpointID: 10,
			},
			expectedPath: "/api-definitions/v2/endpoints/10/versions",
			expectedResponse: &ListEndpointVersionsResponse{
				TotalSize:       1,
				Page:            2,
				PageSize:        3,
				APIEndpointID:   10,
				APIEndpointName: "testingEndpoint",
				APIVersions: []APIVersion{
					{
						CreateDate:           "2020-02-27T15:22:46+0000",
						CreatedBy:            "user",
						UpdateDate:           "2020-02-27T15:22:46+0000",
						UpdatedBy:            "user",
						APIEndpointVersionID: 111222,
						BasePath:             "/production/test",
						VersionNumber:        3,
						Description:          ptr.To("Test description"),
						BasedOn:              ptr.To(int64(2)),
						StagingStatus:        ptr.To(ActivationStatusPending),
						ProductionStatus:     ptr.To(ActivationStatusActive),
						StagingDate:          ptr.To("2022-02-27T15:22:46+0000"),
						ProductionDate:       ptr.To("2022-02-27T15:22:46+0000"),
						IsVersionLocked:      true,
						Hidden:               true,
						AvailableActions: []string{
							"COMPARE_AAG_SETTINGS",
							"EDIT_ENDPOINT_DEFINITION",
							"VIEW_AAG_SETTINGS",
							"ACTIVATE_ON_PRODUCTION",
							"ACTIVATE_ON_STAGING",
							"CLONE_VERSION",
							"EDIT_AAG_SETTINGS",
							"COMPARE_RESOURCE_PURPOSES",
							"HIDE_VERSION",
							"COMPARE_ENDPOINT",
							"DELETE",
							"RESOURCES",
							"VIEW_TAPIOCA",
						},
						CloningStatus: ptr.To("FakeStatus"),
						LockVersion:   0,
					},
					{
						CreateDate:           "2020-01-22T20:06:07+0000",
						CreatedBy:            "user2",
						UpdateDate:           "2020-01-22T20:06:19+0000",
						UpdatedBy:            "user3",
						APIEndpointVersionID: 999000,
						BasePath:             "/production/test",
						VersionNumber:        2,
						Description:          nil,
						BasedOn:              ptr.To(int64(1)),
						StagingStatus:        nil,
						ProductionStatus:     ptr.To(ActivationStatusDeactivated),
						StagingDate:          nil,
						ProductionDate:       ptr.To("2020-01-22T21:26:45+0000"),
						IsVersionLocked:      true,
						Hidden:               false,
						AvailableActions: []string{
							"COMPARE_AAG_SETTINGS",
							"VIEW_AAG_SETTINGS",
							"ACTIVATE_ON_PRODUCTION",
							"ACTIVATE_ON_STAGING",
							"CLONE_VERSION",
							"COMPARE_RESOURCE_PURPOSES",
							"HIDE_VERSION",
							"COMPARE_ENDPOINT",
							"VIEW_TAPIOCA",
						},
						CloningStatus: nil,
						LockVersion:   2,
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"totalSize": 1,
	"page": 2,
	"pageSize": 3,
    "apiEndPointId": 10,
    "apiEndPointName": "testingEndpoint",
    "apiVersions": [
        {
            "createdBy": "user",
            "createDate": "2020-02-27T15:22:46+0000",
            "updateDate": "2020-02-27T15:22:46+0000",
            "updatedBy": "user",
            "apiEndPointVersionId": 111222,
            "basePath": "/production/test",
            "versionNumber": 3,
            "description": "Test description",
            "basedOn": 2,
            "stagingStatus": "PENDING",
            "productionStatus": "ACTIVE",
            "stagingDate": "2022-02-27T15:22:46+0000",
            "productionDate": "2022-02-27T15:22:46+0000",
            "isVersionLocked": true,
            "hidden": true,
            "availableActions": [
                "COMPARE_AAG_SETTINGS",
                "EDIT_ENDPOINT_DEFINITION",
                "VIEW_AAG_SETTINGS",
                "ACTIVATE_ON_PRODUCTION",
                "ACTIVATE_ON_STAGING",
                "CLONE_VERSION",
                "EDIT_AAG_SETTINGS",
                "COMPARE_RESOURCE_PURPOSES",
                "HIDE_VERSION",
                "COMPARE_ENDPOINT",
                "DELETE",
                "RESOURCES",
                "VIEW_TAPIOCA"
            ],
            "cloningStatus": "FakeStatus",
            "lockVersion": 0
        },
        {
            "createdBy": "user2",
            "createDate": "2020-01-22T20:06:07+0000",
            "updateDate": "2020-01-22T20:06:19+0000",
            "updatedBy": "user3",
            "apiEndPointVersionId": 999000,
            "basePath": "/production/test",
            "versionNumber": 2,
            "description": null,
            "basedOn": 1,
            "stagingStatus": null,
            "productionStatus": "DEACTIVATED",
            "stagingDate": null,
            "productionDate": "2020-01-22T21:26:45+0000",
            "isVersionLocked": true,
            "hidden": false,
            "availableActions": [
                "COMPARE_AAG_SETTINGS",
                "VIEW_AAG_SETTINGS",
                "ACTIVATE_ON_PRODUCTION",
                "ACTIVATE_ON_STAGING",
                "CLONE_VERSION",
                "COMPARE_RESOURCE_PURPOSES",
                "HIDE_VERSION",
                "COMPARE_ENDPOINT",
                "VIEW_TAPIOCA"
            ],
            "cloningStatus": null,
            "lockVersion": 2
        }
    ]
}
`,
		},
		"200 OK with query params": {
			params: ListEndpointVersionsRequest{
				APIEndpointID: 2,
				Page:          1,
				PageSize:      1,
				SortBy:        DescriptionSort,
				SortOrder:     AscSortOrder,
				Show:          AllVisibility,
			},
			expectedPath: "/api-definitions/v2/endpoints/2/versions?page=1&pageSize=1&show=ALL&sortBy=description&sortOrder=asc",
			expectedResponse: &ListEndpointVersionsResponse{
				TotalSize:       3,
				Page:            1,
				PageSize:        1,
				APIEndpointID:   2,
				APIEndpointName: "test_cache2",
				APIVersions: []APIVersion{
					{
						CreateDate:           "2022-08-11T06:33:29+0000",
						CreatedBy:            "tester",
						UpdateDate:           "2022-08-16T10:08:41+0000",
						UpdatedBy:            "tester",
						APIEndpointVersionID: 10101010,
						BasePath:             "/test",
						VersionNumber:        1,
						Description:          nil,
						BasedOn:              nil,
						StagingStatus:        nil,
						ProductionStatus:     nil,
						StagingDate:          nil,
						ProductionDate:       nil,
						IsVersionLocked:      false,
						Hidden:               false,
						AvailableActions: []string{
							"ACTIVATE_ON_STAGING",
							"COMPARE_ENDPOINT",
							"DELETE",
							"CLONE_VERSION",
							"EDIT_AAG_SETTINGS",
							"ACTIVATE_ON_PRODUCTION",
							"COMPARE_RESOURCE_PURPOSES",
							"HIDE_VERSION",
							"RESOURCES",
							"EDIT_ENDPOINT_DEFINITION",
							"VIEW_AAG_SETTINGS",
							"COMPARE_AAG_SETTINGS",
							"VIEW_TAPIOCA"},
						CloningStatus: nil,
						LockVersion:   1,
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"totalSize": 3,
	"page": 1,
	"pageSize": 1,
    "apiEndPointId": 2,
    "apiEndPointName": "test_cache2",
    "apiVersions": [
        {
            "createdBy": "tester",
            "createDate": "2022-08-11T06:33:29+0000",
            "updateDate": "2022-08-16T10:08:41+0000",
            "updatedBy": "tester",
            "apiEndPointVersionId": 10101010,
            "basePath": "/test",
            "versionNumber": 1,
            "description": null,
            "basedOn": null,
            "stagingStatus": null,
            "productionStatus": null,
            "stagingDate": null,
            "productionDate": null,
            "isVersionLocked": false,
            "hidden": false,
            "availableActions": [
                "ACTIVATE_ON_STAGING",
                "COMPARE_ENDPOINT",
                "DELETE",
                "CLONE_VERSION",
                "EDIT_AAG_SETTINGS",
                "ACTIVATE_ON_PRODUCTION",
                "COMPARE_RESOURCE_PURPOSES",
                "HIDE_VERSION",
                "RESOURCES",
                "EDIT_ENDPOINT_DEFINITION",
                "VIEW_AAG_SETTINGS",
                "COMPARE_AAG_SETTINGS",
                "VIEW_TAPIOCA"
            ],
            "cloningStatus": null,
            "lockVersion": 1
        }
    ]
}
`,
		},
		"403 forbidden": {
			params: ListEndpointVersionsRequest{
				APIEndpointID: 1,
			},
			expectedPath:   "/api-definitions/v2/endpoints/1/versions",
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "/api-definitions/error-types/forbidden",
    "status": 403,
    "title": "Forbidden",
    "instance": "TestInstance123",
    "detail": "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
    "severity": "ERROR",
    "stackTrace": "StackTraceTest"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/api-definitions/error-types/forbidden",
					Title:    "Forbidden",
					Instance: "TestInstance123",
					Status:   http.StatusForbidden,
					Detail:   "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"404 not found": {
			params: ListEndpointVersionsRequest{
				APIEndpointID: 10101,
			},
			expectedPath:   "/api-definitions/v2/endpoints/10101/versions",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "test.com/resource-impl/forward-origin-error",
    "title": "Not Found",
    "status": 404,
    "instance": "TestInstance123",
    "method": "GET",
    "serverIp": "1.1.1.1",
    "clientIp": "2.2.2.2",
    "requestId": "3222db8",
    "requestTime": "2022-08-19T08:13:39Z"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "test.com/resource-impl/forward-origin-error",
					Title:       "Not Found",
					Instance:    "TestInstance123",
					Status:      http.StatusNotFound,
					Method:      ptr.To("GET"),
					ServerIP:    ptr.To("1.1.1.1"),
					ClientIP:    ptr.To("2.2.2.2"),
					RequestID:   ptr.To("3222db8"),
					RequestTime: ptr.To("2022-08-19T08:13:39Z"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"required param not provided": {
			params: ListEndpointVersionsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list endpoint versions: struct validation: APIEndpointID: cannot be blank", err.Error())
			},
		},
		"query params invalid": {
			params: ListEndpointVersionsRequest{
				APIEndpointID: 1,
				Page:          1,
				PageSize:      1,
				SortBy:        "Wrong",
				SortOrder:     "Wrong",
				Show:          "Wrong",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list endpoint versions: struct validation: Show: value 'Wrong' is invalid. Must be one of: 'ALL', 'ONLY_HIDDEN', 'ONLY_VISIBLE'.\nSortBy: value 'Wrong' is invalid. Must be one of: 'description', 'versionNumber', 'updateDate', 'updatedBy', 'basedOn', 'stagingStatus', 'productionStatus'.\nSortOrder: value 'Wrong' is invalid. Must be one of: 'asc', 'desc'.", err.Error())
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListEndpointVersions(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateEndpointVersion(t *testing.T) {
	tests := map[string]struct {
		params              UpdateEndpointVersionRequest
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *UpdateEndpointVersionResponse
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateEndpointVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 123,
				Body: UpdateEndpointVersionRequestBody{
					SecurityScheme: &SecurityScheme{
						SecuritySchemeType: "apikey",
						SecuritySchemeDetail: SecuritySchemeDetail{
							APIKeyLocation: "cookie",
							APIKeyName:     "name",
						},
					},
					AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
						AllowUndefinedParams:          ptr.To(restrictionsBool(true)),
						PositiveSecurityVersion:       ptr.To(int64(2)),
						PositiveSecurityEnabled:       ptr.To(restrictionsBool(true)),
						AllowUndefinedResources:       ptr.To(restrictionsBool(true)),
						AllowOnlySpecUndefinedMethods: ptr.To(restrictionsBool(false)),
					},
					ContractID:            "TestContract",
					GroupID:               123,
					APIEndpointID:         123,
					APIEndpointVersion:    ptr.To(int64(321)),
					VersionNumber:         1,
					APIEndpointName:       "TestName",
					Description:           ptr.To("Test description"),
					BasePath:              "/test",
					APIEndpointScheme:     ptr.To(APIEndpointSchemeHTTP),
					ConsumeType:           ptr.To(ConsumeTypeAny),
					APIEndpointHosts:      []string{"test.com"},
					APICategoryIDs:        nil,
					LockVersion:           1,
					CaseSensitive:         ptr.To(false),
					MatchPathSegmentParam: true,
					IsGraphQL:             false,
					APIGatewayEnabled:     ptr.To(false),
					GraphQL:               false,
					APIVersionInfo: &APIVersionInfo{
						Location:      "HEADER",
						ParameterName: "ds",
						Value:         "qw",
					},
					APIResources: []APIResource{
						{
							APIResourceClonedFromID:    ptr.To(int64(1010)),
							APIResourceID:              ptr.To(int64(1)),
							APIResourceLogicID:         ptr.To(int64(2)),
							APIResourceMethodNameLists: nil,
							APIResourceMethods: []APIResourceMethod{
								{
									APIResourceMethodID:      ptr.To(int64(34)),
									APIResourceMethodLogicID: ptr.To(int64(234)),
									APIResourceMethod:        APIResourceMethodsPUT,
									APIParameters:            nil,
								},
							},
							APIResourceName: "testResource",
							Description:     "TestDesc",
							Link:            ptr.To("/Test1"),
							LockVersion:     ptr.To(int64(2)),
							Private:         ptr.To(false),
							ResourcePath:    "/res1",
						},
					},
				},
			},
			expectedPath: "/api-definitions/v2/endpoints/123/versions/1",
			expectedRequestBody: `
{
    "apiEndPointId": 123,
    "apiEndPointName": "TestName",
    "description": "Test description",
    "basePath": "/test",
    "consumeType": "any",
    "apiEndPointScheme": "http",
    "apiEndPointVersion": 321,
    "contractId": "TestContract",
    "groupId": 123,
    "versionNumber": 1,
    "apiGatewayEnabled": false,
    "apiEndPointHosts": [
        "test.com"
    ],
    "apiCategoryIds": null,
	"apiVersionInfo": {
		"location": "HEADER",
		"parameterName": "ds",
		"value": "qw"
	},
    "matchPathSegmentParam": true,
    "securityScheme": {
		"securitySchemeType": "apikey",
		"securitySchemeDetail": {
			"apiKeyLocation": "cookie",
			"apiKeyName": "name"
		}
	},
    "akamaiSecurityRestrictions": {
        "ALLOW_UNDEFINED_PARAMS": 1,
        "POSITIVE_SECURITY_VERSION": 2,
        "POSITIVE_SECURITY_ENABLED": 1,
        "ALLOW_UNDEFINED_RESOURCES": 1,
        "ALLOW_ONLY_SPEC_UNDEFINED_METHODS": 0
    },
    "graphQL": false,
    "caseSensitive": false,
    "isGraphQL": false,
	"apiResources": [
		{
            "apiResourceId": 1,
            "apiResourceName": "testResource",
            "resourcePath": "/res1",
            "description": "TestDesc",
            "link": "/Test1",
            "apiResourceClonedFromId": 1010,
            "apiResourceLogicId": 2,
            "private": false,
            "apiResourceMethods": [
                {
                    "apiResourceMethodId": 34,
                    "apiResourceMethod": "PUT",
                    "apiResourceMethodLogicId": 234
                }
            ],
            "lockVersion": 2
        }
	],
    "lockVersion": 1
}
`,
			expectedResponse: &UpdateEndpointVersionResponse{
				SecurityScheme: nil,
				AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
					AllowUndefinedParams:          ptr.To(restrictionsBool(true)),
					PositiveSecurityVersion:       ptr.To(int64(2)),
					PositiveSecurityEnabled:       ptr.To(restrictionsBool(true)),
					AllowUndefinedResources:       ptr.To(restrictionsBool(true)),
					AllowOnlySpecUndefinedMethods: ptr.To(restrictionsBool(false)),
				},
				ContractID:                "TestContract",
				GroupID:                   123,
				APIEndpointID:             123,
				APIEndpointVersion:        ptr.To(int64(321)),
				VersionNumber:             1,
				APIEndpointName:           "TestName",
				Description:               ptr.To("Test description"),
				BasePath:                  "/test",
				ClonedFromVersion:         nil,
				APIEndpointLocked:         false,
				APIEndpointScheme:         nil,
				ConsumeType:               ptr.To("any"),
				APIEndpointHosts:          []string{"test.com"},
				APICategoryIDs:            nil,
				LockVersion:               2,
				UpdatedBy:                 "user2",
				CreatedBy:                 "user",
				CreateDate:                "2022-08-18T08:51:22+0000",
				UpdateDate:                "2022-08-18T09:45:55+0000",
				PositiveConstrainsEnabled: false,
				CaseSensitive:             ptr.To(false),
				MatchPathSegmentParam:     true,
				Source:                    nil,
				StagingVersion:            &VersionState{},
				ProductionVersion:         &VersionState{},
				ProductionStatus:          nil,
				StagingStatus:             nil,
				ProtectedByAPIKey:         false,
				IsGraphQL:                 false,
				AvailableActions:          nil,
				VersionHidden:             false,
				EndpointHidden:            false,
				APISource:                 nil,
				APIGatewayEnabled:         ptr.To(false),
				GraphQL:                   false,
				DiscoveredPIIIDs:          []int64{},
				APIVersionInfo:            nil,
				APIResources: []APIResource{
					{
						CreatedBy:                  "testU",
						CreateDate:                 "2022-08-22T09:58:54+0000",
						UpdatedBy:                  "testU",
						UpdateDate:                 "2022-08-22T09:58:54+0000",
						APIResourceClonedFromID:    ptr.To(int64(1010)),
						APIResourceID:              ptr.To(int64(1)),
						APIResourceLogicID:         ptr.To(int64(2)),
						APIResourceMethodNameLists: nil,
						APIResourceMethods: []APIResourceMethod{
							{
								APIResourceMethodID:      ptr.To(int64(34)),
								APIResourceMethodLogicID: ptr.To(int64(234)),
								APIResourceMethod:        APIResourceMethodsPUT,
								APIParameters:            nil,
							},
						},
						APIResourceName: "testResource",
						Description:     "TestDesc",
						Link:            ptr.To("/Test1"),
						LockVersion:     ptr.To(int64(3)),
						Private:         ptr.To(false),
						ResourcePath:    "/res1",
					},
				},
				Locked: false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdBy": "user",
    "createDate": "2022-08-18T08:51:22+0000",
    "updateDate": "2022-08-18T09:45:55+0000",
    "updatedBy": "user2",
    "apiEndPointId": 123,
    "apiEndPointName": "TestName",
    "description": "Test description",
    "basePath": "/test",
    "consumeType": "any",
    "apiEndPointScheme": null,
    "apiEndPointVersion": 321,
    "contractId": "TestContract",
    "groupId": 123,
    "versionNumber": 1,
    "clonedFromVersion": null,
    "apiEndPointLocked": false,
    "stagingVersion": {
        "versionNumber": null,
        "status": null,
        "timestamp": null,
        "lastError": null
    },
    "productionVersion": {
        "versionNumber": null,
        "status": null,
        "timestamp": null,
        "lastError": null
    },
    "protectedByApiKey": false,
    "apiGatewayEnabled": false,
    "apiEndPointHosts": [
        "test.com"
    ],
    "apiCategoryIds": null,
    "source": null,
    "positiveConstrainsEnabled": false,
    "versionHidden": false,
    "endpointHidden": false,
    "matchPathSegmentParam": true,
    "availableActions": null,
    "apiSource": null,
    "apiSourceDetails": null,
    "cloningStatus": null,
    "securityScheme": null,
    "akamaiSecurityRestrictions": {
        "ALLOW_UNDEFINED_PARAMS": 1,
        "POSITIVE_SECURITY_VERSION": 2,
        "POSITIVE_SECURITY_ENABLED": 1,
        "ALLOW_UNDEFINED_RESOURCES": 1,
        "ALLOW_ONLY_SPEC_UNDEFINED_METHODS": 0
    },
    "discoveredPiiIds": [],
    "stagingStatus": null,
    "productionStatus": null,
    "locked": false,
    "graphQL": false,
    "caseSensitive": false,
    "isGraphQL": false,
	"apiResources": [
		{
			"createdBy": "testU",
            "createDate": "2022-08-22T09:58:54+0000",
            "updateDate": "2022-08-22T09:58:54+0000",
            "updatedBy": "testU",
            "apiResourceId": 1,
            "apiResourceName": "testResource",
            "resourcePath": "/res1",
            "description": "TestDesc",
            "link": "/Test1",
            "apiResourceClonedFromId": 1010,
            "apiResourceLogicId": 2,
            "private": false,
            "apiResourceMethods": [
                {
                    "apiResourceMethodId": 34,
                    "apiResourceMethod": "PUT",
                    "apiResourceMethodLogicId": 234
                }
            ],
            "lockVersion": 3
        }
	],
    "lockVersion": 2
}
`,
		},
		"404 not found": {
			params: UpdateEndpointVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 1,
				Body: UpdateEndpointVersionRequestBody{
					SecurityScheme: &SecurityScheme{
						SecuritySchemeType: "apikey",
						SecuritySchemeDetail: SecuritySchemeDetail{
							APIKeyLocation: "cookie",
							APIKeyName:     "name",
						},
					},
					AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
						AllowUndefinedParams:          ptr.To(restrictionsBool(true)),
						PositiveSecurityVersion:       ptr.To(int64(2)),
						PositiveSecurityEnabled:       ptr.To(restrictionsBool(true)),
						AllowUndefinedResources:       ptr.To(restrictionsBool(true)),
						AllowOnlySpecUndefinedMethods: ptr.To(restrictionsBool(false)),
					},
					ContractID:            "TestContract",
					GroupID:               123,
					APIEndpointID:         123,
					APIEndpointVersion:    ptr.To(int64(321)),
					VersionNumber:         1,
					APIEndpointName:       "TestName",
					Description:           ptr.To("Test description"),
					BasePath:              "/test",
					APIEndpointScheme:     ptr.To(APIEndpointSchemeHTTP),
					ConsumeType:           nil,
					APIEndpointHosts:      []string{"test.com"},
					APICategoryIDs:        []int64{12},
					LockVersion:           1,
					CaseSensitive:         ptr.To(true),
					MatchPathSegmentParam: true,
					IsGraphQL:             true,
					GraphQL:               false,
					APIVersionInfo: &APIVersionInfo{
						Location:      "HEADER",
						ParameterName: "ds",
						Value:         "qw",
					},
					APIResources: []APIResource{
						{
							APIResourceClonedFromID:    ptr.To(int64(1010)),
							APIResourceID:              ptr.To(int64(1)),
							APIResourceLogicID:         ptr.To(int64(2)),
							APIResourceMethodNameLists: nil,
							APIResourceMethods: []APIResourceMethod{
								{
									APIResourceMethodID:      ptr.To(int64(34)),
									APIResourceMethodLogicID: ptr.To(int64(234)),
									APIResourceMethod:        APIResourceMethodsPUT,
									APIParameters: []APIParameter{
										{
											APIParameterName:        "Test",
											APIParameterRequired:    false,
											APIParameterLocation:    "header",
											PathParamLocationID:     nil,
											APIParameterType:        "string",
											Array:                   nil,
											APIParameterNotes:       nil,
											APIParameterRestriction: nil,
											APIChildParameters:      nil,
											APIParameterID:          nil,
											APIParamLogicID:         nil,
											APIResourceMethParamID:  nil,
										},
									},
								},
							},
							APIResourceName: "testResource",
							Description:     "TestDesc",
							Link:            ptr.To("/Test1"),
							LockVersion:     ptr.To(int64(2)),
							Private:         ptr.To(false),
							ResourcePath:    "/res1",
						},
					},
				},
			},
			expectedRequestBody: `
{
    "apiEndPointId": 123,
    "apiEndPointName": "TestName",
    "description": "Test description",
    "basePath": "/test",
    "apiEndPointScheme": "http",
    "apiEndPointVersion": 321,
    "contractId": "TestContract",
    "groupId": 123,
    "versionNumber": 1,
    "apiEndPointHosts": [
        "test.com"
    ],
    "apiCategoryIds": [12],
	"apiVersionInfo": {
		"location": "HEADER",
		"parameterName": "ds",
		"value": "qw"
	},
    "matchPathSegmentParam": true,
    "securityScheme": {
		"securitySchemeType": "apikey",
		"securitySchemeDetail": {
			"apiKeyLocation": "cookie",
			"apiKeyName": "name"
		}
	},
    "akamaiSecurityRestrictions": {
        "ALLOW_UNDEFINED_PARAMS": 1,
        "POSITIVE_SECURITY_VERSION": 2,
        "POSITIVE_SECURITY_ENABLED": 1,
        "ALLOW_UNDEFINED_RESOURCES": 1,
        "ALLOW_ONLY_SPEC_UNDEFINED_METHODS": 0
    },
    "graphQL": false,
    "caseSensitive": true,
    "isGraphQL": true,
	"apiResources": [
		{
            "apiResourceId": 1,
            "apiResourceName": "testResource",
            "resourcePath": "/res1",
            "description": "TestDesc",
            "link": "/Test1",
            "apiResourceClonedFromId": 1010,
            "apiResourceLogicId": 2,
            "private": false,
            "apiResourceMethods": [
                {
                    "apiResourceMethodId": 34,
                    "apiResourceMethod": "PUT",
                    "apiResourceMethodLogicId": 234,
					"apiParameters": [
						{
							"apiParameterName": "Test",
							"apiParameterLocation": "header",
							"apiParameterType": "string",
							"apiParameterRequired": false,
							"apiChildParameters": null
						}
					]
                }
            ],
            "lockVersion": 2
        }
	],
    "lockVersion": 1
}
`,
			expectedPath:   "/api-definitions/v2/endpoints/1/versions/1",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "test.com/resource-impl/forward-origin-error",
    "title": "Not Found",
    "status": 404,
    "instance": "TestInstance123",
    "method": "GET",
    "serverIp": "1.1.1.1",
    "clientIp": "2.2.2.2",
    "requestId": "3222db8",
    "requestTime": "2022-08-19T08:13:39Z"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "test.com/resource-impl/forward-origin-error",
					Title:       "Not Found",
					Instance:    "TestInstance123",
					Status:      http.StatusNotFound,
					Method:      ptr.To("GET"),
					ServerIP:    ptr.To("1.1.1.1"),
					ClientIP:    ptr.To("2.2.2.2"),
					RequestID:   ptr.To("3222db8"),
					RequestTime: ptr.To("2022-08-19T08:13:39Z"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"409 conflict": {
			params: UpdateEndpointVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 1,
				Body: UpdateEndpointVersionRequestBody{
					SecurityScheme: &SecurityScheme{
						SecuritySchemeType: "apikey",
						SecuritySchemeDetail: SecuritySchemeDetail{
							APIKeyLocation: "cookie",
							APIKeyName:     "name",
						},
					},
					AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
						PositiveSecurityVersion:       ptr.To(int64(2)),
						AllowUndefinedParams:          ptr.To(restrictionsBool(true)),
						PositiveSecurityEnabled:       ptr.To(restrictionsBool(true)),
						AllowUndefinedResources:       ptr.To(restrictionsBool(true)),
						AllowOnlySpecUndefinedMethods: ptr.To(restrictionsBool(false)),
					},
					ContractID:            "TestContract",
					GroupID:               123,
					APIEndpointID:         123,
					APIEndpointVersion:    ptr.To(int64(321)),
					VersionNumber:         1,
					APIEndpointName:       "TestName",
					Description:           ptr.To("Test description"),
					BasePath:              "/test",
					APIEndpointScheme:     ptr.To(APIEndpointSchemeHTTP),
					ConsumeType:           nil,
					APIEndpointHosts:      []string{"test.com"},
					APICategoryIDs:        []int64{12},
					LockVersion:           1,
					MatchPathSegmentParam: true,
					IsGraphQL:             true,
					APIGatewayEnabled:     ptr.To(false),
					GraphQL:               false,
					APIVersionInfo: &APIVersionInfo{
						Location:      "HEADER",
						ParameterName: "ds",
						Value:         "qw",
					},
					APIResources: []APIResource{
						{
							APIResourceClonedFromID:    ptr.To(int64(1010)),
							APIResourceID:              ptr.To(int64(1)),
							APIResourceLogicID:         ptr.To(int64(2)),
							APIResourceMethodNameLists: nil,
							APIResourceMethods: []APIResourceMethod{
								{
									APIResourceMethodID:      ptr.To(int64(34)),
									APIResourceMethodLogicID: ptr.To(int64(234)),
									APIResourceMethod:        APIResourceMethodsPUT,
									APIParameters: []APIParameter{
										{
											APIParameterName:        "Test",
											APIParameterRequired:    false,
											APIParameterLocation:    "header",
											PathParamLocationID:     nil,
											APIParameterType:        "string",
											Array:                   nil,
											APIParameterNotes:       nil,
											APIParameterRestriction: nil,
											APIChildParameters:      nil,
											APIParameterID:          nil,
											APIParamLogicID:         nil,
											APIResourceMethParamID:  nil,
										},
									},
								},
							},
							APIResourceName: "testResource",
							Description:     "TestDesc",
							Link:            ptr.To("/Test1"),
							LockVersion:     ptr.To(int64(2)),
							Private:         ptr.To(false),
							ResourcePath:    "/res1",
						},
					},
				},
			},
			expectedPath: "/api-definitions/v2/endpoints/1/versions/1",
			expectedRequestBody: `
{
    "apiEndPointId": 123,
    "apiEndPointName": "TestName",
    "description": "Test description",
    "basePath": "/test",
    "apiEndPointScheme": "http",
    "apiEndPointVersion": 321,
    "contractId": "TestContract",
    "groupId": 123,
    "versionNumber": 1,
    "apiGatewayEnabled": false,
    "apiEndPointHosts": [
        "test.com"
    ],
    "apiCategoryIds": [12],
	"apiVersionInfo": {
		"location": "HEADER",
		"parameterName": "ds",
		"value": "qw"
	},
    "matchPathSegmentParam": true,
    "securityScheme": {
		"securitySchemeType": "apikey",
		"securitySchemeDetail": {
			"apiKeyLocation": "cookie",
			"apiKeyName": "name"
		}
	},
    "akamaiSecurityRestrictions": {
        "ALLOW_UNDEFINED_PARAMS": 1,
        "POSITIVE_SECURITY_VERSION": 2,
        "POSITIVE_SECURITY_ENABLED": 1,
        "ALLOW_UNDEFINED_RESOURCES": 1,
        "ALLOW_ONLY_SPEC_UNDEFINED_METHODS": 0
    },
    "graphQL": false,
    "isGraphQL": true,
	"apiResources": [
		{
            "apiResourceId": 1,
            "apiResourceName": "testResource",
            "resourcePath": "/res1",
            "description": "TestDesc",
            "link": "/Test1",
            "apiResourceClonedFromId": 1010,
            "apiResourceLogicId": 2,
            "private": false,
            "apiResourceMethods": [
                {
                    "apiResourceMethodId": 34,
                    "apiResourceMethod": "PUT",
                    "apiResourceMethodLogicId": 234,
					"apiParameters": [
						{
							"apiParameterName": "Test",
							"apiParameterLocation": "header",
							"apiParameterType": "string",
							"apiParameterRequired": false,
							"apiChildParameters": null
						}
					]
                }
            ],
            "lockVersion": 2
        }
	],
    "lockVersion": 1
}
`,
			responseStatus: http.StatusConflict,
			responseBody: `
{
    "type": "/api-definitions/error-types/CONCURRENT-MODIFICATION-ERROR",
    "status": 409,
    "title": "Concurrent Modification Error",
    "instance": "TestInstance123",
    "detail": "API Endpoint does not allow concurrent modification. Please get the latest API Definition and try again.",
    "severity": "ERROR",
    "stackTrace": "com.akamai.portal.apidefinitions.rest.ProblemOccurredException: Concurrent Modification Error\n\tat com.akamai.portal.apidefinitions.config.jersey.mapper.RestErrorExceptionMapper.toProblemOccurredException(RestErrorExceptionMapper.java:39)\n\tat com.akamai.portal.apidefinitions.config.jersey.mapper.RestErrorExceptionMapper.toProblemOccurredException(RestErrorExceptionMapper.java:13)\n\tat com.akamai.portal.apidefinitions.config.jersey.mapper.AbstractExceptionMapper.toResponse(AbstractExceptionMapper.java:54)\n\tat com.akamai.portal.apidefinitions.config.jersey.mapper.RestErrorExceptionMapper.toResponse(RestErrorExceptionMapper.java:48)\n\tat com.akamai.portal.apidefinitions.config.jersey.mapper.RestErrorExceptionMapper.toResponse(RestErrorExceptionMapper.java:13)\n\tat org.glassfish.jersey.server.ServerRuntime$Responder.mapException(ServerRuntime.java:528)\n\tat org.glassfish.jersey.server.ServerRuntime$Responder.process(ServerRuntime.java:405)\n\tat org.glassfish.jersey.server.ServerRuntime$1.run(ServerRuntime.java:263)\n\tat org.glassfish.jersey.internal.Errors$1.call(Errors.java:248)\n\tat org.glassfish.jersey.internal.Errors$1.call(Errors.java:244)\n\tat org.glassfish.jersey.internal.Errors.process(Errors.java:292)\n\tat org.glassfish.jersey.internal.Errors.process(Errors.java:274)\n\tat org.glassfish.jersey.internal.Errors.process(Errors.java:244)\n\tat org.glassfish.jersey.process.internal.RequestScope.runInScope(RequestScope.java:265)\n\tat org.glassfish.jersey.server.ServerRuntime.process(ServerRuntime.java:234)\n\tat org.glassfish.jersey.server.ApplicationHandler.handle(ApplicationHandler.java:684)\n\tat org.glassfish.jersey.servlet.WebComponent.serviceImpl(WebComponent.java:394)\n\tat org.glassfish.jersey.servlet.WebComponent.service(WebComponent.java:346)\n\tat org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:366)\n\tat org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:319)\n\tat org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:205)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:227)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat org.apache.tomcat.websocket.server.WsFilter.doFilter(WsFilter.java:53)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat com.akamai.portal.apidefinitions.config.MetricsFilter.doFilterInternal(MetricsFilter.java:33)\n\tat org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat com.akamai.se.logging.filters.RequestLoggingFilter.doFilter(RequestLoggingFilter.java:90)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat org.springframework.web.filter.RequestContextFilter.doFilterInternal(RequestContextFilter.java:100)\n\tat org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat com.akamai.ids.sdk.innergrid.servlet.filter.InnerGridFilter.continueProcessingRequest(InnerGridFilter.java:111)\n\tat com.akamai.ids.sdk.innergrid.servlet.filter.InnerGridFilter.performValidation(InnerGridFilter.java:96)\n\tat com.akamai.ids.sdk.innergrid.servlet.filter.InnerGridFilter.doFilter(InnerGridFilter.java:86)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat com.akamai.security.servicesgateway.identity.PulsarSecurityFilter.doFilter(PulsarSecurityFilter.java:57)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat org.springframework.web.filter.FormContentFilter.doFilterInternal(FormContentFilter.java:93)\n\tat org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat org.springframework.boot.actuate.metrics.web.servlet.WebMvcMetricsFilter.doFilterInternal(WebMvcMetricsFilter.java:96)\n\tat org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat org.springframework.web.filter.CharacterEncodingFilter.doFilterInternal(CharacterEncodingFilter.java:201)\n\tat org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat com.akamai.portal.apidefinitions.config.DebuggingHeadersFilter.doFilterInternal(DebuggingHeadersFilter.java:43)\n\tat org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117)\n\tat org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189)\n\tat org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162)\n\tat org.apache.catalina.core.StandardWrapperValve.invoke(StandardWrapperValve.java:197)\n\tat org.apache.catalina.core.StandardContextValve.invoke(StandardContextValve.java:97)\n\tat org.apache.catalina.authenticator.AuthenticatorBase.invoke(AuthenticatorBase.java:541)\n\tat org.apache.catalina.core.StandardHostValve.invoke(StandardHostValve.java:135)\n\tat org.apache.catalina.valves.ErrorReportValve.invoke(ErrorReportValve.java:92)\n\tat org.apache.catalina.core.StandardEngineValve.invoke(StandardEngineValve.java:78)\n\tat org.apache.catalina.connector.CoyoteAdapter.service(CoyoteAdapter.java:360)\n\tat org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:399)\n\tat org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:65)\n\tat org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:890)\n\tat org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1789)\n\tat org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:49)\n\tat org.apache.tomcat.util.threads.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1191)\n\tat org.apache.tomcat.util.threads.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:659)\n\tat org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:61)\n\tat java.base/java.lang.Thread.run(Thread.java:829)\nCaused by: com.akamai.security.servicesgateway.restapi.exceptions.ConcurrentModificationException: API Endpoint does not allow concurrent modification. Please get the latest API Definition and try again.\n\tat com.akamai.portal.apidefinitions.service.impl.ApiEndpointServiceImpl.updateEndpointVersion(ApiEndpointServiceImpl.java:2070)\n\tat com.akamai.portal.apidefinitions.service.impl.ApiEndpointServiceImpl.updateEndpointVersion(ApiEndpointServiceImpl.java:976)\n\tat com.akamai.portal.apidefinitions.service.impl.ApiEndpointServiceImpl$$FastClassBySpringCGLIB$$d229a872.invoke(<generated>)\n\tat org.springframework.cglib.proxy.MethodProxy.invoke(MethodProxy.java:218)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.invokeJoinpoint(CglibAopProxy.java:793)\n\tat org.springframework.aop.framework.ReflectiveMethodInvocation.proceed(ReflectiveMethodInvocation.java:163)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.proceed(CglibAopProxy.java:763)\n\tat org.springframework.aop.aspectj.MethodInvocationProceedingJoinPoint.proceed(MethodInvocationProceedingJoinPoint.java:89)\n\tat com.akamai.portal.apidefinitions.config.TransactionalMetricsAspect.timedMethod(TransactionalMetricsAspect.java:47)\n\tat jdk.internal.reflect.GeneratedMethodAccessor299.invoke(Unknown Source)\n\tat java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)\n\tat java.base/java.lang.reflect.Method.invoke(Method.java:566)\n\tat org.springframework.aop.aspectj.AbstractAspectJAdvice.invokeAdviceMethodWithGivenArgs(AbstractAspectJAdvice.java:634)\n\tat org.springframework.aop.aspectj.AbstractAspectJAdvice.invokeAdviceMethod(AbstractAspectJAdvice.java:624)\n\tat org.springframework.aop.aspectj.AspectJAroundAdvice.invoke(AspectJAroundAdvice.java:72)\n\tat org.springframework.aop.framework.ReflectiveMethodInvocation.proceed(ReflectiveMethodInvocation.java:186)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.proceed(CglibAopProxy.java:763)\n\tat org.springframework.transaction.interceptor.TransactionInterceptor$1.proceedWithInvocation(TransactionInterceptor.java:123)\n\tat org.springframework.transaction.interceptor.TransactionAspectSupport.invokeWithinTransaction(TransactionAspectSupport.java:388)\n\tat org.springframework.transaction.interceptor.TransactionInterceptor.invoke(TransactionInterceptor.java:119)\n\tat org.springframework.aop.framework.ReflectiveMethodInvocation.proceed(ReflectiveMethodInvocation.java:186)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.proceed(CglibAopProxy.java:763)\n\tat org.springframework.aop.interceptor.ExposeInvocationInterceptor.invoke(ExposeInvocationInterceptor.java:97)\n\tat org.springframework.aop.framework.ReflectiveMethodInvocation.proceed(ReflectiveMethodInvocation.java:186)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.proceed(CglibAopProxy.java:763)\n\tat org.springframework.aop.framework.CglibAopProxy$DynamicAdvisedInterceptor.intercept(CglibAopProxy.java:708)\n\tat com.akamai.portal.apidefinitions.service.impl.ApiEndpointServiceImpl$$EnhancerBySpringCGLIB$$909262de.updateEndpointVersion(<generated>)\n\tat com.akamai.portal.apidefinitions.rest.resource.ApiEndpointV2Controller.updateEndpoint(ApiEndpointV2Controller.java:433)\n\tat com.akamai.portal.apidefinitions.rest.resource.ApiEndpointV2Controller$$FastClassBySpringCGLIB$$54b2cd18.invoke(<generated>)\n\tat org.springframework.cglib.proxy.MethodProxy.invoke(MethodProxy.java:218)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.invokeJoinpoint(CglibAopProxy.java:793)\n\tat org.springframework.aop.framework.ReflectiveMethodInvocation.proceed(ReflectiveMethodInvocation.java:163)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.proceed(CglibAopProxy.java:763)\n\tat org.springframework.aop.aspectj.AspectJAfterThrowingAdvice.invoke(AspectJAfterThrowingAdvice.java:64)\n\tat org.springframework.aop.framework.ReflectiveMethodInvocation.proceed(ReflectiveMethodInvocation.java:186)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.proceed(CglibAopProxy.java:763)\n\tat org.springframework.aop.framework.adapter.MethodBeforeAdviceInterceptor.invoke(MethodBeforeAdviceInterceptor.java:58)\n\tat org.springframework.aop.framework.ReflectiveMethodInvocation.proceed(ReflectiveMethodInvocation.java:175)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.proceed(CglibAopProxy.java:763)\n\tat org.springframework.aop.interceptor.ExposeInvocationInterceptor.invoke(ExposeInvocationInterceptor.java:97)\n\tat org.springframework.aop.framework.ReflectiveMethodInvocation.proceed(ReflectiveMethodInvocation.java:186)\n\tat org.springframework.aop.framework.CglibAopProxy$CglibMethodInvocation.proceed(CglibAopProxy.java:763)\n\tat org.springframework.aop.framework.CglibAopProxy$DynamicAdvisedInterceptor.intercept(CglibAopProxy.java:708)\n\tat com.akamai.portal.apidefinitions.rest.resource.ApiEndpointV2Controller$$EnhancerBySpringCGLIB$$f7ae2e3d.updateEndpoint(<generated>)\n\tat jdk.internal.reflect.GeneratedMethodAccessor880.invoke(Unknown Source)\n\tat java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)\n\tat java.base/java.lang.reflect.Method.invoke(Method.java:566)\n\tat org.glassfish.jersey.server.model.internal.ResourceMethodInvocationHandlerFactory.lambda$static$0(ResourceMethodInvocationHandlerFactory.java:52)\n\tat org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher$1.run(AbstractJavaResourceMethodDispatcher.java:124)\n\tat org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher.invoke(AbstractJavaResourceMethodDispatcher.java:167)\n\tat org.glassfish.jersey.server.model.internal.JavaResourceMethodDispatcherProvider$TypeOutInvoker.doDispatch(JavaResourceMethodDispatcherProvider.java:219)\n\tat org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher.dispatch(AbstractJavaResourceMethodDispatcher.java:79)\n\tat org.glassfish.jersey.server.model.ResourceMethodInvoker.invoke(ResourceMethodInvoker.java:475)\n\tat org.glassfish.jersey.server.model.ResourceMethodInvoker.apply(ResourceMethodInvoker.java:397)\n\tat org.glassfish.jersey.server.model.ResourceMethodInvoker.apply(ResourceMethodInvoker.java:81)\n\tat org.glassfish.jersey.server.ServerRuntime$1.run(ServerRuntime.java:255)\n\t... 69 more\n"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/api-definitions/error-types/CONCURRENT-MODIFICATION-ERROR",
					Title:    "Concurrent Modification Error",
					Instance: "TestInstance123",
					Detail:   "API Endpoint does not allow concurrent modification. Please get the latest API Definition and try again.",
					Status:   http.StatusConflict,
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation errors": {
			params: UpdateEndpointVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update endpoint version: struct validation: APIEndpointID: cannot be blank\nBody: {\n\tAPIEndpointHosts: cannot be blank\n\tAPIEndpointID: cannot be blank\n\tAPIEndpointName: cannot be blank\n\tContractID: cannot be blank\n\tGroupID: cannot be blank\n}\nVersionNumber: cannot be blank", err.Error())
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)

				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, test.expectedRequestBody, string(requestBody))

				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateEndpointVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetEndpointVersion(t *testing.T) {
	tests := map[string]struct {
		params           GetEndpointVersionRequest
		expectedPath     string
		expectedResponse *GetEndpointVersionResponse
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetEndpointVersionRequest{
				VersionNumber: 10,
				APIEndpointID: 3,
			},
			expectedPath: "/api-definitions/v2/endpoints/3/versions/10/resources-detail",
			expectedResponse: &GetEndpointVersionResponse{
				SecurityScheme: nil,
				AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
					PositiveSecurityEnabled:       ptr.To(restrictionsBool(true)),
					PositiveSecurityVersion:       ptr.To(int64(2)),
					AllowUndefinedResources:       ptr.To(restrictionsBool(true)),
					AllowOnlySpecUndefinedMethods: ptr.To(restrictionsBool(false)),
					AllowUndefinedParams:          ptr.To(restrictionsBool(true)),
				},
				ContractID:                "TestContract",
				GroupID:                   111,
				APIEndpointID:             3,
				APIEndpointVersion:        ptr.To(int64(999)),
				VersionNumber:             10,
				APIEndpointName:           "Test",
				Description:               ptr.To("Test desc"),
				BasePath:                  "/test",
				ClonedFromVersion:         ptr.To(int64(222)),
				APIEndpointLocked:         false,
				APIEndpointScheme:         nil,
				ConsumeType:               ptr.To("any"),
				APIEndpointHosts:          []string{"test.com"},
				APICategoryIDs:            []int64{456},
				LockVersion:               10,
				UpdatedBy:                 "user2",
				CreatedBy:                 "user",
				CreateDate:                "2022-08-18T08:51:22+0000",
				UpdateDate:                "2022-08-18T09:45:55+0000",
				PositiveConstrainsEnabled: true,
				CaseSensitive:             ptr.To(false),
				MatchPathSegmentParam:     false,
				Source:                    nil,
				StagingVersion: &VersionState{
					VersionNumber: ptr.To(int64(10)),
					Status:        ptr.To(ActivationStatusDeactivated),
					Timestamp:     ptr.To("2022-07-06T09:12:04+0000"),
					LastError:     nil,
				},
				ProductionVersion: &VersionState{},
				ProductionStatus:  nil,
				StagingStatus:     nil,
				ProtectedByAPIKey: false,
				IsGraphQL:         false,
				AvailableActions: []string{
					"ACTIVATE_ON_PRODUCTION",
					"CLONE_ENDPOINT",
					"DELETE",
					"HIDE_ENDPOINT",
					"ACTIVATE_ON_STAGING",
					"EDIT_ENDPOINT_DEFINITION",
				},
				VersionHidden:     false,
				EndpointHidden:    false,
				APISource:         ptr.To("USER"),
				APISourceDetails:  nil,
				CloningStatus:     nil,
				APIGatewayEnabled: ptr.To(true),
				GraphQL:           false,
				DiscoveredPIIIDs:  []int64{},
				APIVersionInfo: &APIVersionInfo{
					Location: "BASE_PATH",
				},
				APIResources: []APIResource{},
				Locked:       false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdBy": "user",
    "createDate": "2022-08-18T08:51:22+0000",
    "updateDate": "2022-08-18T09:45:55+0000",
    "updatedBy": "user2",
    "apiEndPointId": 3,
    "apiEndPointName": "Test",
    "description": "Test desc",
    "basePath": "/test",
    "consumeType": "any",
    "apiEndPointScheme": null,
    "apiEndPointVersion": 999,
    "contractId": "TestContract",
    "groupId": 111,
    "versionNumber": 10,
    "clonedFromVersion": 222,
    "apiEndPointLocked": false,
    "stagingVersion": {
        "versionNumber": 10,
        "status": "DEACTIVATED",
        "timestamp": "2022-07-06T09:12:04+0000",
        "lastError": null
    },
    "productionVersion": {
        "versionNumber": null,
        "status": null,
        "timestamp": null,
        "lastError": null
    },
    "protectedByApiKey": false,
    "apiGatewayEnabled": true,
    "apiEndPointHosts": [
        "test.com"
    ],
    "apiCategoryIds": [
        456
    ],
    "source": null,
    "apiVersionInfo": {
        "location": "BASE_PATH",
        "parameterName": null,
        "value": null
    },
    "positiveConstrainsEnabled": true,
    "versionHidden": false,
    "endpointHidden": false,
    "matchPathSegmentParam": false,
    "availableActions": [
        "ACTIVATE_ON_PRODUCTION",
        "CLONE_ENDPOINT",
        "DELETE",
        "HIDE_ENDPOINT",
        "ACTIVATE_ON_STAGING",
        "EDIT_ENDPOINT_DEFINITION"
    ],
    "apiSource": "USER",
    "apiSourceDetails": null,
    "cloningStatus": null,
    "securityScheme": null,
    "akamaiSecurityRestrictions": {
        "ALLOW_UNDEFINED_PARAMS": 1,
        "POSITIVE_SECURITY_VERSION": 2,
        "POSITIVE_SECURITY_ENABLED": 1,
        "ALLOW_UNDEFINED_RESOURCES": 1,
        "ALLOW_ONLY_SPEC_UNDEFINED_METHODS": 0
    },
    "discoveredPiiIds": [],
    "stagingStatus": null,
    "productionStatus": null,
    "locked": false,
    "graphQL": false,
    "caseSensitive": false,
    "isGraphQL": false,
    "apiResources": [],
    "lockVersion": 10
}
`,
			withError: nil,
		},
		"200 OK all fields populated": {
			params: GetEndpointVersionRequest{
				VersionNumber: 123,
				APIEndpointID: 456,
			},
			expectedPath: "/api-definitions/v2/endpoints/456/versions/123/resources-detail",
			expectedResponse: &GetEndpointVersionResponse{
				SecurityScheme: &SecurityScheme{
					SecuritySchemeType: "apikey",
					SecuritySchemeDetail: SecuritySchemeDetail{
						APIKeyName:     "test",
						APIKeyLocation: "cookie",
					},
				},
				AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
					MaxJSONXMLElement:       ptr.To(int64(5)),
					MaxElementNameLength:    ptr.To(int64(100)),
					MaxStringLength:         ptr.To(int64(40)),
					MaxIntegerValue:         ptr.To(int64(90)),
					MaxDocDepth:             ptr.To(int64(20)),
					MaxBodySize:             ptr.To(int64(400)),
					AllowUndefinedMethodGet: ptr.To(restrictionsBool(true)),
				},
				ContractID:                "Test-Contract123",
				GroupID:                   88888,
				APIEndpointID:             456,
				APIEndpointVersion:        ptr.To(int64(1010)),
				VersionNumber:             123,
				APIEndpointName:           "test",
				Description:               ptr.To("test description"),
				BasePath:                  "/test",
				ClonedFromVersion:         ptr.To(int64(111)),
				APIEndpointLocked:         true,
				APIEndpointScheme:         ptr.To("http"),
				ConsumeType:               ptr.To("any"),
				APIEndpointHosts:          []string{"testing.com"},
				APICategoryIDs:            []int64{},
				LockVersion:               1,
				UpdatedBy:                 "tester",
				CreatedBy:                 "tester",
				CreateDate:                "2022-08-11T06:33:29+0000",
				UpdateDate:                "2022-08-16T10:08:41+0000",
				PositiveConstrainsEnabled: true,
				CaseSensitive:             ptr.To(true),
				MatchPathSegmentParam:     true,
				Source: &Source{
					Type:                 "TestType",
					APIVersion:           "TestVersion",
					SpecificationVersion: "Test",
				},
				StagingVersion: &VersionState{
					VersionNumber: ptr.To(int64(10)),
					Status:        ptr.To(ActivationStatusActive),
					Timestamp:     ptr.To("2022-08-16T10:08:41+0000"),
					LastError:     nil,
				},
				ProductionVersion: &VersionState{
					VersionNumber: ptr.To(int64(10)),
					Status:        ptr.To(ActivationStatusActive),
					Timestamp:     ptr.To("2022-08-16T10:08:41+0000"),
					LastError:     nil,
				},
				ProductionStatus:  ptr.To("PENDING"),
				StagingStatus:     ptr.To("WAITING"),
				ProtectedByAPIKey: true,
				IsGraphQL:         true,
				AvailableActions: []string{
					"ACTIVATE_ON_PRODUCTION",
					"CLONE_ENDPOINT",
					"DELETE",
					"HIDE_ENDPOINT",
					"ACTIVATE_ON_STAGING",
					"EDIT_ENDPOINT_DEFINITION",
				},
				DiscoveredPIIIDs: []int64{},
				VersionHidden:    true,
				EndpointHidden:   true,
				APISource:        ptr.To("USER"),
				APISourceDetails: []APISourceDiff{
					{
						Name:        "TestName",
						SourceValue: "TestSource",
						SavedValue:  "TestSaved",
					},
				},
				CloningStatus:     ptr.To("PENDING"),
				APIGatewayEnabled: ptr.To(true),
				GraphQL:           true,
				APIVersionInfo: &APIVersionInfo{
					Location:      "QUERY",
					ParameterName: "QUERY",
					Value:         "test",
				},
				APIResources: []APIResource{
					{
						APIResourceClonedFromID:    ptr.To(int64(10)),
						APIResourceID:              ptr.To(int64(111)),
						APIResourceLogicID:         ptr.To(int64(112)),
						APIResourceMethodNameLists: nil,
						APIResourceMethods: []APIResourceMethod{
							{
								APIResourceMethodID:      ptr.To(int64(44)),
								APIResourceMethodLogicID: ptr.To(int64(55)),
								APIResourceMethod:        APIResourceMethodsGET,
								APIParameters: []APIParameter{
									{
										APIParameterName:     "TestName",
										APIParameterRequired: false,
										APIParameterLocation: "TestLocation",
										PathParamLocationID:  ptr.To(int64(999)),
										APIParameterType:     "TestType",
										Array:                ptr.To(true),
										APIParameterNotes:    ptr.To("TestNotes"),
										APIParameterRestriction: &APIParameterRestriction{
											LengthRestriction: &LengthRestriction{
												LengthMax: 10,
												LengthMin: 1,
											},
											RangeRestriction: &RangeRestriction{
												RangeMin: 1,
												RangeMax: 15,
											},
											NumberRangeRestriction: &NumberRangeRestriction{
												NumberRangeMin: 2,
												NumberRangeMax: 30,
											},
											ArrayRestriction: &ArrayRestriction{
												MaxItems: 100,
												MinItems: 20,
											},
											XMLConversionRule: &XMLConversionRule{
												Attribute: false,
												Wrapped:   true,
												Name:      "Testing",
												Namespace: "Test",
												Prefix:    "T",
											},
										},
										APIChildParameters:     nil,
										APIParameterID:         ptr.To(int64(10)),
										APIParamLogicID:        ptr.To(int64(11)),
										APIResourceMethParamID: ptr.To(int64(12)),
									},
								},
							},
						},
						APIResourceName: "test_api",
						CreateDate:      "2022-08-11T06:33:29+0000",
						CreatedBy:       "test",
						Description:     "Test desc",
						Link:            ptr.To("TestLink"),
						LockVersion:     ptr.To(int64(0)),
						Private:         ptr.To(false),
						ResourcePath:    "/test_api",
						UpdateDate:      "2022-08-11T06:33:29+0000",
						UpdatedBy:       "test",
					},
					{
						APIResourceClonedFromID:    nil,
						APIResourceID:              ptr.To(int64(321)),
						APIResourceLogicID:         ptr.To(int64(4321)),
						APIResourceMethodNameLists: nil,
						APIResourceMethods: []APIResourceMethod{
							{
								APIResourceMethodID:      ptr.To(int64(789)),
								APIResourceMethodLogicID: ptr.To(int64(6789)),
								APIResourceMethod:        APIResourceMethodsGET,
								APIParameters:            []APIParameter{},
							},
						},
						APIResourceName: "test_api3",
						CreateDate:      "2022-08-11T06:33:29+0000",
						CreatedBy:       "Test user",
						Description:     "",
						Link:            nil,
						LockVersion:     ptr.To(int64(0)),
						Private:         ptr.To(false),
						ResourcePath:    "/test_api3",
						UpdateDate:      "2022-08-11T06:33:29+0000",
						UpdatedBy:       "Test user",
					},
				},
				Locked: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdBy": "tester",
    "createDate": "2022-08-11T06:33:29+0000",
    "updateDate": "2022-08-16T10:08:41+0000",
    "updatedBy": "tester",
    "apiEndPointId": 456,
    "apiEndPointName": "test",
    "description": "test description",
    "basePath": "/test",
    "consumeType": "any",
    "apiEndPointScheme": "http",
    "apiEndPointVersion": 1010,
    "contractId": "Test-Contract123",
    "groupId": 88888,
    "versionNumber": 123,
    "clonedFromVersion": 111,
    "apiEndPointLocked": true,
    "stagingVersion": {
        "versionNumber": 10,
        "status": "ACTIVE",
        "timestamp": "2022-08-16T10:08:41+0000",
        "lastError": null
    },
    "productionVersion": {
        "versionNumber": 10,
        "status": "ACTIVE",
        "timestamp": "2022-08-16T10:08:41+0000",
        "lastError": null
    },
    "protectedByApiKey": true,
    "apiGatewayEnabled": true,
    "apiEndPointHosts": [
        "testing.com"
    ],
    "apiCategoryIds": [],
    "source": {
		"type": "TestType",
		"apiVersion": "TestVersion",
		"specificationVersion": "Test"
	},
    "apiVersionInfo": {
        "location": "QUERY",
        "parameterName": "QUERY",
        "value": "test"
    },
    "positiveConstrainsEnabled": true,
    "versionHidden": true,
    "endpointHidden": true,
    "matchPathSegmentParam": true,
    "availableActions": [
        "ACTIVATE_ON_PRODUCTION",
        "CLONE_ENDPOINT",
        "DELETE",
        "HIDE_ENDPOINT",
        "ACTIVATE_ON_STAGING",
        "EDIT_ENDPOINT_DEFINITION"
    ],
    "apiSource": "USER",
    "apiSourceDetails": [
		{
			"name": "TestName",
			"sourceValue": "TestSource",
			"savedValue": "TestSaved"
		}
	],
    "cloningStatus": "PENDING",
    "securityScheme": {
        "securitySchemeType": "apikey",
        "securitySchemeDetail": {
            "apiKeyName": "test",
            "apiKeyLocation": "cookie"
        }
    },
    "akamaiSecurityRestrictions": {
        "MAX_JSONXML_ELEMENT": 5,
        "MAX_ELEMENT_NAME_LENGTH": 100,
        "MAX_STRING_LENGTH": 40,
        "MAX_INTEGER_VALUE": 90,
        "MAX_DOC_DEPTH": 20,
		"MAX_BODY_SIZE": 400,
		"ALLOW_UNDEFINED_METHOD_GET": 1
    },
    "discoveredPiiIds": [],
    "stagingStatus": "WAITING",
    "productionStatus": "PENDING",
    "locked": true,
    "graphQL": true,
    "caseSensitive": true,
    "isGraphQL": true,
    "apiResources": [
        {
            "createdBy": "test",
            "createDate": "2022-08-11T06:33:29+0000",
            "updateDate": "2022-08-11T06:33:29+0000",
            "updatedBy": "test",
            "apiResourceId": 111,
            "apiResourceName": "test_api",
            "resourcePath": "/test_api",
            "description": "Test desc",
            "link": "TestLink",
            "apiResourceClonedFromId": 10,
            "apiResourceLogicId": 112,
            "private": false,
            "apiResourceMethods": [
                {
                    "apiResourceMethodId": 44,
                    "apiResourceMethod": "GET",
                    "apiParameters": [
						{
							"apiParameterName": "TestName",
							"apiParameterRequired": false,
							"apiParameterLocation": "TestLocation",
							"pathParamLocationId": 999,
							"apiParameterType": "TestType",
							"array": true,
							"apiParameterNotes": "TestNotes",
							"apiParameterRestriction": {
								"lengthRestriction": {
									"lengthMax": 10,
									"lengthMin": 1
								},
								"rangeRestriction": {
									"rangeMin": 1,
									"rangeMax": 15
								},
								"numberRangeRestriction": {
									"numberRangeMin": 2,
									"numberRangeMax": 30
								},
								"arrayRestriction": {
									"collectionFormat": "Test",
									"maxItems": 100,
									"minItems": 20,
									"uniqueItems": false
								},
								"xmlConversionRule": {
									"attribute": false,
									"wrapped": true,
									"name": "Testing",
									"namespace": "Test",
									"prefix": "T"
								}
							},
							"apiParameterId": 10,
							"apiParamLogicId": 11,
							"apiResourceMethParamID": 12
						}
					],
                    "apiResourceMethodLogicId": 55,
                    "methodRestrictions": null
                }
            ],
            "lockVersion": 0
        },
        {
            "createdBy": "Test user",
            "createDate": "2022-08-11T06:33:29+0000",
            "updateDate": "2022-08-11T06:33:29+0000",
            "updatedBy": "Test user",
            "apiResourceId": 321,
            "apiResourceName": "test_api3",
            "resourcePath": "/test_api3",
            "description": null,
            "link": null,
            "apiResourceClonedFromId": null,
            "apiResourceLogicId": 4321,
            "private": false,
            "apiResourceMethods": [
                {
                    "apiResourceMethodId": 789,
                    "apiResourceMethod": "GET",
                    "apiParameters": [],
                    "apiResourceMethodLogicId": 6789,
                    "methodRestrictions": null
                }
            ],
            "lockVersion": 0
        }
    ],
    "lockVersion": 1
}
`,
			withError: nil,
		},
		"403 forbidden": {
			params: GetEndpointVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 1,
			},
			expectedPath:   "/api-definitions/v2/endpoints/1/versions/1/resources-detail",
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "/api-definitions/error-types/forbidden",
    "status": 403,
    "title": "Forbidden",
    "instance": "TestInstance123",
    "detail": "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
    "severity": "ERROR",
    "stackTrace": "StackTraceTest"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/api-definitions/error-types/forbidden",
					Title:    "Forbidden",
					Instance: "TestInstance123",
					Status:   http.StatusForbidden,
					Detail:   "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"404 not found": {
			params: GetEndpointVersionRequest{
				VersionNumber: 1010,
				APIEndpointID: 1,
			},
			expectedPath:   "/api-definitions/v2/endpoints/1/versions/1010/resources-detail",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "test.com/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "TestInstance123",
    "detail": "No Api Endpoint/Version found for endpoint ID 1 and version 10",
    "severity": "ERROR",
    "stackTrace": "StackTraceTest"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "test.com/api-definitions/error-types/NOT-FOUND",
					Title:    "Not Found",
					Detail:   "No Api Endpoint/Version found for endpoint ID 1 and version 10",
					Instance: "TestInstance123",
					Status:   http.StatusNotFound,
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"required param not provided": {
			params: GetEndpointVersionRequest{
				VersionNumber: 10,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get endpoint version: struct validation: APIEndpointID: cannot be blank", err.Error())
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetEndpointVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCloneEndpointVersion(t *testing.T) {
	tests := map[string]struct {
		params           CloneEndpointVersionRequest
		expectedPath     string
		expectedResponse *CloneEndpointVersionResponse
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: CloneEndpointVersionRequest{
				VersionNumber: 2,
				APIEndpointID: 123,
			},
			expectedPath: "/api-definitions/v2/endpoints/123/versions/2/cloneVersion",
			expectedResponse: &CloneEndpointVersionResponse{
				SecurityScheme: nil,
				AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
					AllowUndefinedParams:          ptr.To(restrictionsBool(true)),
					AllowUndefinedResources:       ptr.To(restrictionsBool(true)),
					AllowOnlySpecUndefinedMethods: ptr.To(restrictionsBool(false)),
					PositiveSecurityEnabled:       ptr.To(restrictionsBool(true)),
					PositiveSecurityVersion:       ptr.To(int64(2)),
				},
				ContractID:         "Contract123",
				GroupID:            222,
				APIEndpointID:      123,
				APIEndpointVersion: ptr.To(int64(111000)),
				VersionNumber:      3,
				APIEndpointName:    "IPQA_DXE_TEST_2",
				Description:        ptr.To("Test"),
				BasePath:           "/test",
				ClonedFromVersion:  ptr.To(int64(110000)),
				APIEndpointLocked:  false,
				APIEndpointScheme:  nil,
				ConsumeType:        ptr.To("any"),
				APIEndpointHosts: []string{
					"test.com",
					"test2.com",
				},
				APICategoryIDs:            []int64{12},
				LockVersion:               1,
				UpdatedBy:                 "user1",
				CreatedBy:                 "user1",
				CreateDate:                "2022-08-19T06:43:33+0000",
				UpdateDate:                "2022-08-19T06:43:33+0000",
				PositiveConstrainsEnabled: true,
				CaseSensitive:             ptr.To(false),
				MatchPathSegmentParam:     false,
				Source:                    nil,
				StagingVersion: &VersionState{
					VersionNumber: ptr.To(int64(2)),
					Status:        ptr.To(ActivationStatusDeactivated),
					Timestamp:     ptr.To("2022-07-06T09:12:04+0000"),
					LastError:     nil,
				},
				ProductionVersion: &VersionState{},
				ProductionStatus:  nil,
				StagingStatus:     nil,
				ProtectedByAPIKey: true,
				IsGraphQL:         false,
				AvailableActions: []string{
					"ACTIVATE_ON_STAGING",
					"CLONE_ENDPOINT",
					"DELETE",
					"ACTIVATE_ON_PRODUCTION",
					"HIDE_ENDPOINT",
					"EDIT_ENDPOINT_DEFINITION",
				},
				DiscoveredPIIIDs:  []int64{},
				VersionHidden:     false,
				EndpointHidden:    false,
				APISource:         ptr.To("USER"),
				APISourceDetails:  nil,
				CloningStatus:     nil,
				APIGatewayEnabled: ptr.To(true),
				GraphQL:           false,
				APIVersionInfo:    nil,
				APIResources:      []APIResource{},
				Locked:            false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdBy": "user1",
    "createDate": "2022-08-19T06:43:33+0000",
    "updateDate": "2022-08-19T06:43:33+0000",
    "updatedBy": "user1",
    "apiEndPointId": 123,
    "apiEndPointName": "IPQA_DXE_TEST_2",
    "description": "Test",
    "basePath": "/test",
    "consumeType": "any",
    "apiEndPointScheme": null,
    "apiEndPointVersion": 111000,
    "contractId": "Contract123",
    "groupId": 222,
    "versionNumber": 3,
    "clonedFromVersion": 110000,
    "apiEndPointLocked": false,
    "stagingVersion": {
        "versionNumber": 2,
        "status": "DEACTIVATED",
        "timestamp": "2022-07-06T09:12:04+0000",
        "lastError": null
    },
    "productionVersion": {
        "versionNumber": null,
        "status": null,
        "timestamp": null,
        "lastError": null
    },
    "protectedByApiKey": true,
    "apiGatewayEnabled": true,
    "apiEndPointHosts": [
        "test.com",
		"test2.com"
    ],
    "apiCategoryIds": [
        12
    ],
    "source": null,
    "positiveConstrainsEnabled": true,
    "versionHidden": false,
    "endpointHidden": false,
    "matchPathSegmentParam": false,
    "availableActions": [
        "ACTIVATE_ON_STAGING",
        "CLONE_ENDPOINT",
        "DELETE",
        "ACTIVATE_ON_PRODUCTION",
        "HIDE_ENDPOINT",
        "EDIT_ENDPOINT_DEFINITION"
    ],
    "apiSource": "USER",
    "apiSourceDetails": null,
    "cloningStatus": null,
    "securityScheme": null,
    "akamaiSecurityRestrictions": {
        "POSITIVE_SECURITY_VERSION": 2,
        "ALLOW_UNDEFINED_PARAMS": 1,
        "ALLOW_UNDEFINED_RESOURCES": 1,
        "POSITIVE_SECURITY_ENABLED": 1,
        "ALLOW_ONLY_SPEC_UNDEFINED_METHODS": 0
    },
    "discoveredPiiIds": [],
    "stagingStatus": null,
    "productionStatus": null,
    "locked": false,
    "graphQL": false,
    "caseSensitive": false,
    "isGraphQL": false,
    "apiResources": [],
    "lockVersion": 1
}
`,
		},
		"200 OK with resources": {
			params: CloneEndpointVersionRequest{
				VersionNumber: 10,
				APIEndpointID: 1000,
			},
			expectedPath: "/api-definitions/v2/endpoints/1000/versions/10/cloneVersion",
			expectedResponse: &CloneEndpointVersionResponse{
				SecurityScheme: &SecurityScheme{
					SecuritySchemeType: "apikey",
					SecuritySchemeDetail: SecuritySchemeDetail{
						APIKeyName:     "key-header",
						APIKeyLocation: "header",
					},
				},
				AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
					PositiveSecurityVersion: ptr.To(int64(2)),
				},
				ContractID:                "Contract2",
				GroupID:                   33,
				APIEndpointID:             1000,
				APIEndpointVersion:        ptr.To(int64(100)),
				VersionNumber:             10,
				APIEndpointName:           "TestEndpoint",
				Description:               nil,
				BasePath:                  "/test2",
				ClonedFromVersion:         ptr.To(int64(200)),
				APIEndpointLocked:         false,
				APIEndpointScheme:         nil,
				ConsumeType:               ptr.To("any"),
				APIEndpointHosts:          []string{"host.test.com"},
				APICategoryIDs:            []int64{11},
				LockVersion:               1,
				UpdatedBy:                 "test",
				CreatedBy:                 "test",
				CreateDate:                "2022-08-19T07:35:43+0000",
				UpdateDate:                "2022-08-19T07:35:43+0000",
				PositiveConstrainsEnabled: false,
				CaseSensitive:             ptr.To(true),
				MatchPathSegmentParam:     true,
				Source:                    nil,
				StagingVersion:            &VersionState{},
				ProductionVersion: &VersionState{
					VersionNumber: ptr.To(int64(2)),
					Status:        ptr.To(ActivationStatusActive),
					Timestamp:     ptr.To("2019-10-02T14:49:09+0000"),
					LastError:     nil,
				},
				ProductionStatus:  nil,
				StagingStatus:     nil,
				ProtectedByAPIKey: true,
				IsGraphQL:         false,
				AvailableActions: []string{
					"ACTIVATE_ON_PRODUCTION",
					"CLONE_ENDPOINT",
					"ACTIVATE_ON_STAGING",
					"DEACTIVATE_ON_PRODUCTION",
					"EDIT_ENDPOINT_DEFINITION",
				},
				DiscoveredPIIIDs:  []int64{},
				VersionHidden:     false,
				EndpointHidden:    false,
				APISource:         ptr.To("USER"),
				APISourceDetails:  nil,
				CloningStatus:     nil,
				APIGatewayEnabled: ptr.To(true),
				GraphQL:           false,
				APIVersionInfo:    nil,
				APIResources: []APIResource{
					{
						APIResourceClonedFromID:    nil,
						APIResourceID:              ptr.To(int64(123)),
						APIResourceLogicID:         ptr.To(int64(456)),
						APIResourceMethodNameLists: nil,
						APIResourceMethods: []APIResourceMethod{
							{
								APIResourceMethodID:      ptr.To(int64(1122)),
								APIResourceMethodLogicID: ptr.To(int64(3344)),
								APIResourceMethod:        APIResourceMethodsGET,
								APIParameters:            []APIParameter{},
							},
							{
								APIResourceMethodID:      ptr.To(int64(5566)),
								APIResourceMethodLogicID: ptr.To(int64(7788)),
								APIResourceMethod:        APIResourceMethodsPOST,
								APIParameters:            []APIParameter{},
							},
						},
						APIResourceName: "testRes",
						CreateDate:      "2022-08-19T07:35:43+0000",
						CreatedBy:       "user1",
						Description:     "",
						Link:            nil,
						LockVersion:     ptr.To(int64(0)),
						Private:         ptr.To(false),
						ResourcePath:    "test/200",
						UpdateDate:      "2022-08-19T07:35:43+0000",
						UpdatedBy:       "user1",
					},
				},
				Locked: false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdBy": "test",
    "createDate": "2022-08-19T07:35:43+0000",
    "updateDate": "2022-08-19T07:35:43+0000",
    "updatedBy": "test",
    "apiEndPointId": 1000,
    "apiEndPointName": "TestEndpoint",
    "description": null,
    "basePath": "/test2",
    "consumeType": "any",
    "apiEndPointScheme": null,
    "apiEndPointVersion": 100,
    "contractId": "Contract2",
    "groupId": 33,
    "versionNumber": 10,
    "clonedFromVersion": 200,
    "apiEndPointLocked": false,
    "stagingVersion": {
        "versionNumber": null,
        "status": null,
        "timestamp": null,
        "lastError": null
    },
    "productionVersion": {
        "versionNumber": 2,
        "status": "ACTIVE",
        "timestamp": "2019-10-02T14:49:09+0000",
        "lastError": null
    },
    "protectedByApiKey": true,
    "apiGatewayEnabled": true,
    "apiEndPointHosts": [
        "host.test.com"
    ],
    "apiCategoryIds": [
        11
    ],
    "source": null,
    "positiveConstrainsEnabled": false,
    "versionHidden": false,
    "endpointHidden": false,
    "matchPathSegmentParam": true,
    "availableActions": [
        "ACTIVATE_ON_PRODUCTION",
        "CLONE_ENDPOINT",
        "ACTIVATE_ON_STAGING",
        "DEACTIVATE_ON_PRODUCTION",
        "EDIT_ENDPOINT_DEFINITION"
    ],
    "apiSource": "USER",
    "apiSourceDetails": null,
    "cloningStatus": null,
    "securityScheme": {
        "securitySchemeType": "apikey",
        "securitySchemeDetail": {
            "apiKeyLocation": "header",
            "apiKeyName": "key-header"
        }
    },
    "akamaiSecurityRestrictions": {
        "POSITIVE_SECURITY_VERSION": 2
    },
    "discoveredPiiIds": [],
    "stagingStatus": null,
    "productionStatus": null,
    "locked": false,
    "graphQL": false,
    "caseSensitive": true,
    "isGraphQL": false,
    "apiResources": [
        {
            "createdBy": "user1",
            "createDate": "2022-08-19T07:35:43+0000",
            "updateDate": "2022-08-19T07:35:43+0000",
            "updatedBy": "user1",
            "apiResourceId": 123,
            "apiResourceName": "testRes",
            "resourcePath": "test/200",
            "description": null,
            "link": null,
            "apiResourceClonedFromId": null,
            "apiResourceLogicId": 456,
            "private": false,
            "apiResourceMethods": [
                {
                    "apiResourceMethodId": 1122,
                    "apiResourceMethod": "GET",
                    "apiParameters": [],
                    "apiResourceMethodLogicId": 3344,
                    "methodRestrictions": null
                },
                {
                    "apiResourceMethodId": 5566,
                    "apiResourceMethod": "POST",
                    "apiParameters": [],
                    "apiResourceMethodLogicId": 7788,
                    "methodRestrictions": null
                }
            ],
            "lockVersion": 0
        }
    ],
    "lockVersion": 1
}
`,
			withError: nil,
		},
		"403 forbidden": {
			params: CloneEndpointVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 1,
			},
			expectedPath:   "/api-definitions/v2/endpoints/1/versions/1/cloneVersion",
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "/api-definitions/error-types/forbidden",
    "status": 403,
    "title": "Forbidden",
    "instance": "TestInstance123",
    "detail": "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
    "severity": "ERROR",
    "stackTrace": "StackTraceTest"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/api-definitions/error-types/forbidden",
					Title:    "Forbidden",
					Instance: "TestInstance123",
					Status:   http.StatusForbidden,
					Detail:   "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"404 not found": {
			params: CloneEndpointVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 10101,
			},
			expectedPath:   "/api-definitions/v2/endpoints/10101/versions/1/cloneVersion",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "test.com/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "TestInstance123",
    "detail": "Cannot delete Api Endpoint Version. No Version found for API ID 10101",
    "severity": "ERROR",
    "stackTrace": "StackTrace123.com"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "test.com/api-definitions/error-types/NOT-FOUND",
					Title:    "Not Found",
					Detail:   "Cannot delete Api Endpoint Version. No Version found for API ID 10101",
					Instance: "TestInstance123",
					Status:   http.StatusNotFound,
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"required params not provided": {
			params: CloneEndpointVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "clone endpoint version: struct validation: APIEndpointID: cannot be blank\nVersionNumber: cannot be blank", err.Error())
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.CloneEndpointVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteEndpointVersion(t *testing.T) {
	tests := map[string]struct {
		params         DeleteEndpointVersionRequest
		expectedPath   string
		responseStatus int
		responseBody   string
		withError      func(*testing.T, error)
	}{
		"204 no content": {
			params: DeleteEndpointVersionRequest{
				VersionNumber: 10,
				APIEndpointID: 333000444,
			},
			expectedPath:   "/api-definitions/v2/endpoints/333000444/versions/10",
			responseStatus: http.StatusNoContent,
			responseBody:   ``,
			withError:      nil,
		},
		"403 forbidden": {
			params: DeleteEndpointVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 1,
			},
			expectedPath:   "/api-definitions/v2/endpoints/1/versions/1",
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "/api-definitions/error-types/forbidden",
    "status": 403,
    "title": "Forbidden",
    "instance": "TestInstance123",
    "detail": "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
    "severity": "ERROR",
    "stackTrace": "StackTraceTest"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/api-definitions/error-types/forbidden",
					Title:    "Forbidden",
					Instance: "TestInstance123",
					Status:   http.StatusForbidden,
					Detail:   "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"404 not found": {
			params: DeleteEndpointVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 10101,
			},
			expectedPath:   "/api-definitions/v2/endpoints/10101/versions/1",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "test.com/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "TestInstance123",
    "detail": "Cannot delete Api Endpoint Version. No Version found for API ID 10101",
    "severity": "ERROR",
    "stackTrace": "StackTrace123.com"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "test.com/api-definitions/error-types/NOT-FOUND",
					Title:    "Not Found",
					Detail:   "Cannot delete Api Endpoint Version. No Version found for API ID 10101",
					Instance: "TestInstance123",
					Status:   http.StatusNotFound,
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"required param not provided": {
			params: DeleteEndpointVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete endpoint version: struct validation: APIEndpointID: cannot be blank\nVersionNumber: cannot be blank", err.Error())
			},
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			err := client.DeleteEndpointVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
