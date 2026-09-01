package apidefinitions

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	registerEndpointFileContent = "test-file-content"
)

func TestGetEndpoint(t *testing.T) {
	tests := map[string]struct {
		params           GetEndpointRequest
		expectedPath     string
		expectedResponse *GetEndpointResponse
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetEndpointRequest{
				APIEndpointID: 3,
			},
			expectedPath: "/api-definitions/v2/endpoints/3",
			expectedResponse: &GetEndpointResponse{
				AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
					AllowUndefinedParams: ptr.To(restrictionsBool(true)),
				},
				Endpoint: Endpoint{
					APIEndpointID:      3,
					APIEndpointName:    "Test",
					Description:        ptr.To("Test desc"),
					BasePath:           "/test",
					ConsumeType:        ptr.To(ConsumeTypeAny),
					APIEndpointScheme:  nil,
					APIEndpointVersion: int64(999),
					ContractID:         "TestContract",
					GroupID:            111,
					VersionNumber:      10,
					ClonedFromVersion:  ptr.To(int64(222)),
					Locked:             false,
					StagingVersion: VersionState{
						VersionNumber: ptr.To(int64(10)),
						Status:        ptr.To(ActivationStatusDeactivated),
						Timestamp:     ptr.To("2022-07-06T09:12:04+0000"),
						LastError:     nil,
					},
					ProductionVersion:   VersionState{},
					ProtectedByAPIKey:   false,
					APIGatewayEnabled:   true,
					CaseSensitive:       false,
					APIEndpointHosts:    []string{"test.com"},
					APICategoryIDs:      []int64{456},
					APIResourceBaseInfo: nil,
					Source:              nil,
					APIVersionInfo: &APIVersionInfo{
						Location: "BASE_PATH",
					},
					PositiveConstrainsEnabled: true,
					VersionHidden:             false,
					EndpointHidden:            false,
					IsGraphQL:                 false,
					MatchPathSegmentParam:     false,
					AvailableActions: []string{
						"ACTIVATE_ON_PRODUCTION",
						"CLONE_ENDPOINT",
						"DELETE",
						"HIDE_ENDPOINT",
						"ACTIVATE_ON_STAGING",
						"EDIT_ENDPOINT_DEFINITION",
					},
					APISource:   ptr.To("USER"),
					LockVersion: 10,
					UpdatedBy:   "user2",
					CreatedBy:   "user",
					CreateDate:  "2022-08-18T08:51:22+0000",
					UpdateDate:  "2022-08-18T09:45:55+0000",
				},
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
        "ALLOW_UNDEFINED_PARAMS": 1
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
		"403 forbidden": {
			params: GetEndpointRequest{
				APIEndpointID: 1,
			},
			expectedPath:   "/api-definitions/v2/endpoints/1",
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
			params: GetEndpointRequest{
				APIEndpointID: 1,
			},
			expectedPath:   "/api-definitions/v2/endpoints/1",
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
			params: GetEndpointRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get endpoint: struct validation: APIEndpointID: cannot be blank", err.Error())
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
			result, err := client.GetEndpoint(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListEndpointsRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req           ListEndpointsRequest
		errorExpected bool
	}{
		"min ok": {
			req: ListEndpointsRequest{SortBy: "name"},
		},
		"all ok": {
			req: ListEndpointsRequest{
				SortBy:            "updateDate",
				SortOrder:         "asc",
				VersionPreference: "LAST_UPDATED",
				Show:              "ALL",
			},
		},
		"bad sortBy": {
			req:           ListEndpointsRequest{SortBy: "random"},
			errorExpected: true,
		},
		"bad sortOrder": {
			req: ListEndpointsRequest{
				SortBy:    "name",
				SortOrder: "random",
			},
			errorExpected: true,
		},
		"bad version preference": {
			req: ListEndpointsRequest{
				SortBy:            "name",
				VersionPreference: "random",
			},
			errorExpected: true,
		},
		"bad show": {
			req: ListEndpointsRequest{
				SortBy: "name",
				Show:   "NONE",
			},
			errorExpected: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.errorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestSecurityScheme_Validate(t *testing.T) {
	tests := map[string]struct {
		detail    SecurityScheme
		withError func(t *testing.T, err error)
	}{
		"no securitySchemeType": {
			detail: SecurityScheme{},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "SecuritySchemeType: cannot be blank.")
			},
		},
		"bad securitySchemeType": {
			detail: SecurityScheme{SecuritySchemeType: "parameter"},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "SecuritySchemeType: value 'parameter' is not valid. Must be: 'apikey'.")
			},
		},
		"no apiKeyLocation": {
			detail: SecurityScheme{},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "APIKeyName: cannot be blank.")
			},
		},
		"bad apiKeyLocation": {
			detail: SecurityScheme{
				SecuritySchemeType: "apikey",
				SecuritySchemeDetail: SecuritySchemeDetail{
					APIKeyLocation: "wrong"},
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "SecuritySchemeDetail: (APIKeyLocation: value 'wrong' is not valid. Must be one of: 'cookie', 'header', 'query'")
			},
		},
		"no apiKeyName": {
			detail: SecurityScheme{
				SecuritySchemeType: "apikey",
				SecuritySchemeDetail: SecuritySchemeDetail{
					APIKeyLocation: "header",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "SecuritySchemeDetail: (APIKeyName: cannot be blank.).", err.Error())
			},
		},
		"ok": {
			detail: SecurityScheme{
				SecuritySchemeType: "apikey",
				SecuritySchemeDetail: SecuritySchemeDetail{
					APIKeyLocation: "cookie",
					APIKeyName:     "name",
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.detail.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAPIParameterRestriction_Validate(t *testing.T) {
	tests := map[string]struct {
		detail    APIParameterRestriction
		withError func(t *testing.T, err error)
	}{
		"min OK": {
			detail: APIParameterRestriction{ArrayRestriction: &ArrayRestriction{}},
		},
		"bad LengthRestriction": {
			detail: APIParameterRestriction{LengthRestriction: &LengthRestriction{
				LengthMax: -1,
				LengthMin: -1,
			}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "LengthRestriction: (LengthMax: must be no less than 0; LengthMin: must be no less than 0.).", err.Error())
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.detail.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAPIResourceMethod_Validate(t *testing.T) {
	apiParameterRestriction := APIParameterRestriction{
		LengthRestriction: &LengthRestriction{},
		ArrayRestriction:  &ArrayRestriction{},
	}
	tests := map[string]struct {
		detail    APIResourceMethod
		withError func(t *testing.T, err error)
	}{
		"no APIResourceMethod": {
			detail: APIResourceMethod{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "APIResourceMethod: cannot be blank.", err.Error())
			},
		},
		"no APIParameterName": {
			detail: APIResourceMethod{APIResourceMethod: APIResourceMethodsGET,
				APIParameters: []APIParameter{{
					APIParameterLocation:    "query",
					APIParameterType:        "string",
					APIParameterRestriction: &apiParameterRestriction,
				}}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "APIParameters: (0: (APIParameterName: cannot be blank.).).", err.Error())
			},
		},
		"no APIParameterLocation": {
			detail: APIResourceMethod{APIResourceMethod: APIResourceMethodsGET,
				APIParameters: []APIParameter{{
					APIParameterName:        "a name",
					APIParameterType:        "string",
					APIParameterRestriction: &apiParameterRestriction,
				}}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "APIParameters: (0: (APIParameterLocation: cannot be blank.).).", err.Error())
			},
		},
		"bad APIParameterLocation": {
			detail: APIResourceMethod{APIResourceMethod: APIResourceMethodsGET,
				APIParameters: []APIParameter{{
					APIParameterName:        "a name",
					APIParameterLocation:    "bad",
					APIParameterType:        "string",
					APIParameterRestriction: &apiParameterRestriction,
				}}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "APIParameters: (0: (APIParameterLocation: value 'bad' is not valid. Must be one of: 'query', 'header', 'path', 'cookie', 'body'.).).", err.Error())
			},
		},
		"no APIParameterType": {
			detail: APIResourceMethod{APIResourceMethod: APIResourceMethodsGET,
				APIParameters: []APIParameter{{
					APIParameterName:        "a name",
					APIParameterLocation:    "query",
					APIParameterRestriction: &apiParameterRestriction,
				}}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "APIParameters: (0: (APIParameterType: cannot be blank.).).", err.Error())
			},
		},
		"bad APIParameterType": {
			detail: APIResourceMethod{APIResourceMethod: APIResourceMethodsGET,
				APIParameters: []APIParameter{{
					APIParameterName:        "a name",
					APIParameterLocation:    "query",
					APIParameterType:        "float",
					APIParameterRestriction: &apiParameterRestriction,
				}}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "APIParameters: (0: (APIParameterType: value 'float' is not valid. Must be one of 'string', 'integer', 'number', 'boolean', or 'json/xml'.).).", err.Error())
			},
		},
		"min ok": {
			detail: APIResourceMethod{APIResourceMethod: APIResourceMethodsGET,
				APIParameters: []APIParameter{{
					APIParameterName:        "a name",
					APIParameterLocation:    "query",
					APIParameterType:        "string",
					APIParameterRestriction: &apiParameterRestriction,
				}}},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.detail.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAkamaiSecurityRestrictions_Validate(t *testing.T) {
	negativeTests := map[string]AkamaiSecurityRestrictions{
		"POSITIVE_SECURITY_VERSION": {PositiveSecurityVersion: ptr.To(int64(3))},
	}
	for name, restriction := range negativeTests {
		t.Run(name, func(t *testing.T) {
			require.Error(t, restriction.Validate())
		})
	}
	t.Run("happy path", func(t *testing.T) {
		require.NoError(t, AkamaiSecurityRestrictions{}.Validate())
	})
}

func TestRegisterEndpointRequest_Validate(t *testing.T) {
	apiResourceMethod := APIResourceMethod{
		APIResourceMethod: APIResourceMethodsGET,
		APIParameters: []APIParameter{{
			APIParameterName:     "parameter",
			APIParameterLocation: "body",
			APIParameterType:     "json/xml",
			APIParameterRestriction: &APIParameterRestriction{
				LengthRestriction: &LengthRestriction{},
				ArrayRestriction:  &ArrayRestriction{},
			},
		}}}
	apiResource := APIResource{
		APIResourceName:    "api resource",
		ResourcePath:       "/a/path",
		APIResourceMethods: []APIResourceMethod{apiResourceMethod}}
	endpointRequest := RegisterEndpointRequest{
		APIEndpointName:            "name",
		APIEndpointScheme:          "http",
		APIEndpointHosts:           []string{"akamai.com"},
		APIResources:               []APIResource{apiResource},
		APIVersionInfo:             &APIVersionInfo{Location: "BASE_PATH"},
		AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{},
		ContractID:                 "3-XXXXXXX",
		GroupID:                    33333,
	}

	noName := endpointRequest
	noName.APIEndpointName = ""

	basePathLonger := endpointRequest
	basePathLonger.BasePath = "/longer/basepath"

	basePathWithSlash := endpointRequest
	basePathWithSlash.BasePath = "/slash/"

	noHosts := endpointRequest
	noHosts.APIEndpointHosts = []string{}

	noAPIResources := endpointRequest
	noAPIResources.APIResources = nil

	noAPIResourceName := endpointRequest
	apiResourceNoName := apiResource
	apiResourceNoName.APIResourceName = ""
	noAPIResourceName.APIResources = append(noAPIResourceName.APIResources, apiResourceNoName)

	nilAPIVersionInfo := endpointRequest
	nilAPIVersionInfo.APIVersionInfo = nil

	noAPIVersionInfo := endpointRequest
	noAPIVersionInfo.APIVersionInfo = &APIVersionInfo{}

	noSecurityRestrictions := endpointRequest
	noSecurityRestrictions.AkamaiSecurityRestrictions = nil

	badSecurityRestrictions := endpointRequest
	badSecurityRestrictions.AkamaiSecurityRestrictions = &AkamaiSecurityRestrictions{PositiveSecurityVersion: ptr.To(int64(3))}

	noContractRequest := endpointRequest
	noContractRequest.ContractID = ""

	noGroupRequest := endpointRequest
	noGroupRequest.GroupID = 0

	tests := map[string]struct {
		req         RegisterEndpointRequest
		expectError bool
	}{
		"ok (empty basePath)": {
			req: endpointRequest,
		},
		"basePath longer": {
			req: basePathLonger,
		},
		"basePath slash": {
			req:         basePathWithSlash,
			expectError: true,
		},
		"fail no name": {
			req:         noName,
			expectError: true,
		},
		"fail no hosts": {
			req:         noHosts,
			expectError: true,
		},
		"fail no API resource name": {
			req:         noAPIResourceName,
			expectError: true,
		},
		"no API resources": {
			req: noAPIResources,
		},
		"nil apiVersionInfo": {
			req: nilAPIVersionInfo,
		},
		"no apiVersionInfo.Location": {
			req:         noAPIVersionInfo,
			expectError: true,
		},
		"nil akamaiSecurityRestrictions": {
			req: noSecurityRestrictions,
		},
		"bad akamaiSecurityRestrictions": {
			req:         badSecurityRestrictions,
			expectError: true,
		},
		"no contract": {
			req:         noContractRequest,
			expectError: true,
		},
		"no group": {
			req:         noGroupRequest,
			expectError: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAPIResourceMethods_Validate(t *testing.T) {
	tests := map[string]struct {
		met       []APIResourceMethod
		withError func(t *testing.T, err error)
	}{
		"ok": {
			met: []APIResourceMethod{
				{APIResourceMethod: APIResourceMethodsHEAD, APIParameters: []APIParameter{{
					APIParameterName:        "parameter",
					APIParameterLocation:    "path",
					APIParameterType:        "number",
					APIParameterRestriction: &APIParameterRestriction{ArrayRestriction: &ArrayRestriction{}},
				}}},
			},
		},
		"nok": {
			met: []APIResourceMethod{
				{
					APIResourceMethod: "random",
					APIParameters: []APIParameter{{
						APIParameterName:        "random",
						APIParameterLocation:    "path",
						APIParameterType:        "number",
						APIParameterRestriction: &APIParameterRestriction{ArrayRestriction: &ArrayRestriction{}},
					}}},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "0: (APIResourceMethod: value 'random' is not valid. Must be one of: 'GET', 'PUT', 'POST', 'DELETE', 'HEAD', 'PATCH', 'OPTIONS'.).", err.Error())
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := validation.Validate(test.met)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestListEndpoints(t *testing.T) {
	tests := map[string]struct {
		request        ListEndpointsRequest
		withError      func(*testing.T, error)
		responseBody   string
		expectedPath   string
		expectedResult *ListEndpointsResponse
		responseStatus int
	}{
		"OK": {
			responseBody: `{
    "totalSize": 466,
    "page": 1,
    "pageSize": 1,
    "apiEndPoints": [
        {
            "createdBy": "example",
            "createDate": "2020-01-22T20:24:34+0000",
            "updateDate": "2020-01-22T20:24:51+0000",
            "updatedBy": "example",
            "apiEndPointId": 400130,
            "apiEndPointName": "kk-pearlcheckc",
            "description": "The API config",
            "basePath": "/production",
            "consumeType": "any",
            "apiEndPointScheme": null,
            "apiEndPointVersion": 748268,
            "contractId": "3-XXXXXX",
            "groupId": 33333,
            "versionNumber": 5,
            "clonedFromVersion": 664137,
            "apiEndPointLocked": false,
            "stagingVersion": {
                "versionNumber": 1,
                "status": "ACTIVE",
                "timestamp": "2019-01-15T19:38:03+0000",
                "lastError": null
            },
            "productionVersion": {
                "versionNumber": 3,
                "status": "ACTIVE",
                "timestamp": "2019-10-23T18:20:34+0000",
                "lastError": null
            },
            "protectedByApiKey": true,
            "apiGatewayEnabled": true,
            "apiEndPointHosts": [
                "kk-apigateway.wildcard.webexp-ipqa-ion.com"
            ],
            "apiCategoryIds": [
                8944,
                7462
            ],
            "apiResourceBaseInfo": [
                {
                    "createdBy": "example",
                    "createDate": "2020-01-22T20:24:34+0000",
                    "updateDate": "2020-01-22T20:24:51+0000",
                    "updatedBy": "example",
                    "apiResourceId": 3437990,
                    "apiResourceName": "/all/200",
                    "resourcePath": "/all/200",
                    "description": null,
                    "link": null,
                    "apiResourceClonedFromId": null,
                    "apiResourceLogicId": 71334,
                    "private": false,
                    "lockVersion": 0
                },
                {
                    "createdBy": "example",
                    "createDate": "2020-01-22T20:24:34+0000",
                    "updateDate": "2020-01-22T20:24:51+0000",
                    "updatedBy": "example",
                    "apiResourceId": 3437991,
                    "apiResourceName": "/get/200",
                    "resourcePath": "/get/200",
                    "description": null,
                    "link": null,
                    "apiResourceClonedFromId": null,
                    "apiResourceLogicId": 71335,
                    "private": false,
                    "lockVersion": 0
                }
            ],
            "source": null,
            "positiveConstrainsEnabled": false,
            "versionHidden": false,
            "endpointHidden": false,
            "matchPathSegmentParam": true,
            "availableActions": [
                "ACTIVATE_ON_PRODUCTION",
                "ACTIVATE_ON_STAGING",
                "EDIT_ENDPOINT_DEFINITION",
                "DEACTIVATE_ON_PRODUCTION",
                "DEACTIVATE_ON_STAGING",
                "CLONE_ENDPOINT"
            ],
            "apiSource": "USER",
            "locked": false,
            "caseSensitive": true,
            "isGraphQL": false,
            "lockVersion": 1
        }
    ]
}`,
			expectedResult: &ListEndpointsResponse{
				TotalSize: 466,
				Page:      1,
				PageSize:  1,
				APIEndpoints: []Endpoint{
					{
						CreatedBy:         "example",
						CreateDate:        "2020-01-22T20:24:34+0000",
						UpdateDate:        "2020-01-22T20:24:51+0000",
						UpdatedBy:         "example",
						APIEndpointID:     400130,
						APIEndpointName:   "kk-pearlcheckc",
						Description:       ptr.To("The API config"),
						BasePath:          "/production",
						ConsumeType:       ptr.To(ConsumeTypeAny),
						ContractID:        "3-XXXXXX",
						GroupID:           33333,
						VersionNumber:     5,
						ClonedFromVersion: ptr.To(int64(664137)),
						StagingVersion: VersionState{
							VersionNumber: ptr.To(int64(1)),
							Status:        ptr.To(ActivationStatusActive),
							Timestamp:     ptr.To("2019-01-15T19:38:03+0000"),
						},
						ProductionVersion: VersionState{
							VersionNumber: ptr.To(int64(3)),
							Status:        ptr.To(ActivationStatusActive),
							Timestamp:     ptr.To("2019-10-23T18:20:34+0000"),
						},
						ProtectedByAPIKey:     true,
						APIGatewayEnabled:     true,
						APIEndpointHosts:      []string{"kk-apigateway.wildcard.webexp-ipqa-ion.com"},
						APICategoryIDs:        []int64{8944, 7462},
						MatchPathSegmentParam: true,
						AvailableActions: []string{
							"ACTIVATE_ON_PRODUCTION",
							"ACTIVATE_ON_STAGING",
							"EDIT_ENDPOINT_DEFINITION",
							"DEACTIVATE_ON_PRODUCTION",
							"DEACTIVATE_ON_STAGING",
							"CLONE_ENDPOINT",
						},
						APISource:          ptr.To("USER"),
						CaseSensitive:      true,
						LockVersion:        1,
						APIEndpointVersion: 748268,

						APIResourceBaseInfo: []APIResourceBaseInfo{
							{
								CreatedBy:          ptr.To("example"),
								CreateDate:         ptr.To("2020-01-22T20:24:34+0000"),
								UpdateDate:         ptr.To("2020-01-22T20:24:51+0000"),
								UpdatedBy:          ptr.To("example"),
								APIResourceID:      3437990,
								APIResourceName:    "/all/200",
								ResourcePath:       "/all/200",
								APIResourceLogicID: 71334,
							},
							{
								CreatedBy:          ptr.To("example"),
								CreateDate:         ptr.To("2020-01-22T20:24:34+0000"),
								UpdateDate:         ptr.To("2020-01-22T20:24:51+0000"),
								UpdatedBy:          ptr.To("example"),
								APIResourceID:      3437991,
								APIResourceName:    "/get/200",
								ResourcePath:       "/get/200",
								APIResourceLogicID: 71335,
							},
						},
					},
				},
			},
			expectedPath:   "/api-definitions/v2/endpoints?category=__UNCATEGORIZED__&contains=pearl&contractId=3-XXXXXX&groupId=33333&page=1&pageSize=25&show=ALL&sortBy=name&sortOrder=asc&versionPreference=LAST_UPDATED",
			responseStatus: http.StatusOK,
			request: ListEndpointsRequest{
				Category:          "__UNCATEGORIZED__",
				Contains:          "pearl",
				SortBy:            "name",
				SortOrder:         "asc",
				ContractID:        "3-XXXXXX",
				GroupID:           33333,
				VersionPreference: "LAST_UPDATED",
				Show:              "ALL",
			},
		},
		"500 HTTP code": {
			responseBody: `{
    "type": "a",
    "title": "b",
    "detail": "c",
    "status": 500
}`,
			expectedPath:   "/api-definitions/v2/endpoints?category=__UNCATEGORIZED__&contains=pearl&contractId=3-XXXXXX&groupId=33333&page=1&pageSize=25&show=ALL&sortBy=name&sortOrder=asc&versionPreference=LAST_UPDATED",
			responseStatus: http.StatusInternalServerError,
			request: ListEndpointsRequest{
				Category:          "__UNCATEGORIZED__",
				Contains:          "pearl",
				SortBy:            "name",
				SortOrder:         "asc",
				ContractID:        "3-XXXXXX",
				GroupID:           33333,
				VersionPreference: "LAST_UPDATED",
				Show:              "ALL",
			},
			withError: func(t *testing.T, err error) {
				expectedError := Error{
					Type:   "a",
					Title:  "b",
					Detail: "c",
					Status: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, &expectedError), "want: %s; got: %s", expectedError, err)
			},
		},
		"fail validation": {
			request: ListEndpointsRequest{SortBy: "creationDate"},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list endpoints: struct validation: SortBy: value 'creationDate' is not valid. Must be one of: 'name' or 'updateDate'", err.Error())
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
			result, err := client.ListEndpoints(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func TestRegisterEndpoint(t *testing.T) {
	tests := map[string]struct {
		body                RegisterEndpointRequest
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
		expectedResult      *RegisterEndpointResponse
	}{
		"ok created less parameters": {
			responseStatus: http.StatusCreated,
			body: RegisterEndpointRequest{
				APIEndpointName:  "bookstore API",
				APIEndpointHosts: []string{"akamai.com"},
				ContractID:       "3-XXXXXX",
				GroupID:          33333,
			},
			expectedRequestBody: `{
    "apiEndPointName": "bookstore API",
    "contractId": "3-XXXXXX",
    "groupId": 33333,
    "apiEndPointHosts": [
        "akamai.com"
    ]
}`,
			expectedPath: "/api-definitions/v2/endpoints",
			responseBody: `{
    "createdBy": "yyyyyyy",
    "createDate": "2022-06-23T13:08:52+0000",
    "updateDate": "2022-06-23T13:08:53+0000",
    "updatedBy": "yyyyyyy",
    "apiEndPointId": 1231231,
    "apiEndPointName": "bookstore API",
    "description": null,
    "basePath": "",
    "consumeType": null,
    "apiEndPointScheme": null,
    "apiEndPointVersion": 1231231,
    "contractId": "3-XXXXXX",
    "groupId": 33333,
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
    "apiGatewayEnabled": true,
    "apiEndPointHosts": [
        "akamai.com"
    ],
    "apiCategoryIds": [],
    "source": null,
    "positiveConstrainsEnabled": false,
    "versionHidden": false,
    "endpointHidden": false,
    "matchPathSegmentParam": true,
    "availableActions": [
        "HIDE_ENDPOINT",
        "ACTIVATE_ON_STAGING",
        "ACTIVATE_ON_PRODUCTION",
        "CLONE_ENDPOINT",
        "EDIT_ENDPOINT_DEFINITION",
        "DELETE"
    ],
    "apiSource": "USER",
    "securityScheme": null,
    "akamaiSecurityRestrictions": {
        "POSITIVE_SECURITY_VERSION": 2
    },
    "discoveredPiiIds": [],
    "stagingStatus": null,
    "productionStatus": null,
    "locked": false,
    "caseSensitive": true,
    "isGraphQL": false,
    "apiResources": [],
    "lockVersion": 1
}`,
			expectedResult: &RegisterEndpointResponse{
				EndpointWithResources: EndpointWithResources{
					DiscoveredPIIIDs: []int64{},
					APIResources:     []APIResource{},
					EndpointDetail: EndpointDetail{
						AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{PositiveSecurityVersion: ptr.To(int64(2))},
						Endpoint: Endpoint{
							CreatedBy:             "yyyyyyy",
							CreateDate:            "2022-06-23T13:08:52+0000",
							UpdatedBy:             "yyyyyyy",
							UpdateDate:            "2022-06-23T13:08:53+0000",
							APIEndpointID:         1231231,
							APIEndpointVersion:    1231231,
							APIEndpointName:       "bookstore API",
							ContractID:            "3-XXXXXX",
							GroupID:               33333,
							VersionNumber:         1,
							LockVersion:           1,
							CaseSensitive:         true,
							APIEndpointHosts:      []string{"akamai.com"},
							AvailableActions:      []string{"HIDE_ENDPOINT", "ACTIVATE_ON_STAGING", "ACTIVATE_ON_PRODUCTION", "CLONE_ENDPOINT", "EDIT_ENDPOINT_DEFINITION", "DELETE"},
							APICategoryIDs:        []int64{},
							APISource:             ptr.To("USER"),
							MatchPathSegmentParam: true,
							APIGatewayEnabled:     true,
						},
					},
				},
			},
		},
		"ok created more parameters": {
			responseStatus: http.StatusCreated,
			body: RegisterEndpointRequest{
				APIEndpointName:       "bookstore API",
				APIEndpointScheme:     "http",
				APIEndpointHosts:      []string{"akamai.com"},
				APIGatewayEnabled:     ptr.To(false),
				ContractID:            "3-XXXXXX",
				GroupID:               33333,
				DiscoveredPIIIDs:      []int64{1234},
				BasePath:              "/test",
				Description:           "testDescription",
				ConsumeType:           "any",
				CaseSensitive:         ptr.To(false),
				MatchPathSegmentParam: ptr.To(false),
				IsGraphQL:             ptr.To(false),
				APICategoryIDs:        []int64{1234},
				SecurityScheme: &SecurityScheme{
					SecuritySchemeType: "apikey",
					SecuritySchemeDetail: SecuritySchemeDetail{
						APIKeyName:     "test",
						APIKeyLocation: "cookie",
					},
				},
				AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
					AllowUndefinedResources: ptr.To(restrictionsBool(true)),
				},
				APISource: "USER",
				APIResources: []APIResource{
					{
						APIResourceName: "test_api",
						ResourcePath:    "/test_api",
						APIResourceMethods: []APIResourceMethod{
							{
								APIResourceMethod:        APIResourceMethodsGET,
								APIResourceMethodID:      ptr.To(int64(1231234)),
								APIResourceMethodLogicID: ptr.To(int64(1231234)),
								APIParameters: []APIParameter{
									{
										APIParameterName:       "param",
										APIParameterRequired:   false,
										APIParameterLocation:   "cookie",
										APIParameterType:       "integer",
										APIParamLogicID:        ptr.To(int64(1231234)),
										APIResourceMethParamID: ptr.To(int64(1231234)),
										APIChildParameters:     []APIParameter{},
									},
								},
							},
						},
					},
				},
				APIVersionInfo: &APIVersionInfo{
					Location:      "QUERY",
					ParameterName: "QUERY",
					Value:         "test",
				},
			},
			expectedRequestBody: `{
    "apiEndPointName": "bookstore API",
    "apiEndPointScheme": "http",
    "contractId": "3-XXXXXX",
    "groupId": 33333,
    "apiEndPointHosts": [
        "akamai.com"
    ],
	"apiGatewayEnabled": false,
	"basePath": "/test",
	"description": "testDescription",
	"consumeType": "any",
	"caseSensitive": false,
	"matchPathSegmentParam": false,
	"isGraphQL": false,
	"discoveredPiiIds": [1234],
	"apiCategoryIds": [1234],
	"securityScheme": {
        "securitySchemeType": "apikey",
        "securitySchemeDetail": {
            "apiKeyName": "test",
            "apiKeyLocation": "cookie"
        }
    },
	"akamaiSecurityRestrictions": {
        "ALLOW_UNDEFINED_RESOURCES": 1
	},
	"apiSource": "USER",
	"apiResources": [
        {
            "apiResourceName": "test_api",
            "resourcePath": "/test_api",
            "apiResourceMethods": [
                {
                    "apiResourceMethod": "GET",
					"apiResourceMethodLogicId": 1231234,
					"apiResourceMethodId": 1231234,
                    "apiParameters": [
                        {
                            "apiParameterName": "param",
                            "apiParameterRequired": false,
                            "apiParameterLocation": "cookie",
                            "apiParameterType": "integer",
							"apiParamLogicId": 1231234,
                            "apiResourceMethParamId": 1231234,
							"apiChildParameters": []
                        }
                    ]
                }
            ]
        }
    ],
	"apiVersionInfo": {
        "location": "QUERY",
        "parameterName": "QUERY",
        "value": "test"
    }
}`,
			expectedPath: "/api-definitions/v2/endpoints",
			responseBody: `{
    "createdBy": "yyyyyyy",
    "createDate": "2022-06-23T13:08:52+0000",
    "updateDate": "2022-06-23T13:08:53+0000",
    "updatedBy": "yyyyyyy",
    "apiEndPointId": 1297665,
    "apiEndPointName": "bookstore API",
    "description": "testDescription",
    "basePath": "/test",
    "consumeType": "any",
    "apiEndPointScheme": "http",
    "apiEndPointVersion": 1518543,
    "contractId": "3-XXXXXX",
    "groupId": 33333,
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
    "protectedByApiKey": true,
    "apiGatewayEnabled": false,
    "apiEndPointHosts": [
        "akamai.com"
    ],
    "apiCategoryIds": [
        12345
    ],
    "source": null,
    "apiVersionInfo": {
        "location": "QUERY",
        "parameterName": "QUERY",
        "value": "test"
    },
    "positiveConstrainsEnabled": false,
    "versionHidden": false,
    "endpointHidden": false,
    "matchPathSegmentParam": false,
    "availableActions": [
        "HIDE_ENDPOINT",
        "ACTIVATE_ON_STAGING",
        "ACTIVATE_ON_PRODUCTION",
        "CLONE_ENDPOINT",
        "EDIT_ENDPOINT_DEFINITION",
        "DELETE"
    ],
    "apiSource": "USER",
    "securityScheme": {
        "securitySchemeType": "apikey",
        "securitySchemeDetail": {
            "apiKeyLocation": "cookie",
            "apiKeyName": "test"
        }
    },
    "akamaiSecurityRestrictions": {
        "POSITIVE_SECURITY_VERSION": 2,
        "ALLOW_UNDEFINED_RESOURCES": 1
    },
    "discoveredPiiIds": [],
    "stagingStatus": null,
    "productionStatus": null,
    "locked": false,
    "caseSensitive": false,
    "isGraphQL": false,
    "apiResources": [
        {
            "createdBy": "yyyyyyy",
            "createDate": "2022-06-23T13:08:52+0000",
            "updateDate": "2022-06-23T13:08:52+0000",
            "updatedBy": "yyyyyyy",
            "apiResourceId": 123456,
            "apiResourceName": "test_api",
            "resourcePath": "/test_api",
            "description": null,
            "link": null,
            "apiResourceClonedFromId": null,
            "apiResourceLogicId": 1231234,
            "private": false,
            "apiResourceMethods": [
                {
                    "apiResourceMethodId": 1231234,
                    "apiResourceMethod": "GET",
                    "apiParameters": [
                        {
                            "apiParameterId": 1231234,
                            "apiParameterRequired": false,
                            "apiParameterName": "param",
                            "apiParameterLocation": "cookie",
                            "apiParameterType": "integer",
                            "array": false,
                            "pathParamLocationId": null,
                            "apiChildParameters": [],
                            "apiParameterRestriction": null,
                            "apiParameterNotes": null,
							"apiParamLogicId": 1231234,
                            "apiResourceMethParamId": 1231234
                        }
                    ],
                    "apiResourceMethodLogicId": 1231234
                }
            ],
            "lockVersion": 0
        }
    ],
    "lockVersion": 1
}`,
			expectedResult: &RegisterEndpointResponse{
				EndpointWithResources: EndpointWithResources{
					DiscoveredPIIIDs: []int64{},
					APIResources: []APIResource{
						{
							UpdateDate:         "2022-06-23T13:08:52+0000",
							UpdatedBy:          "yyyyyyy",
							CreateDate:         "2022-06-23T13:08:52+0000",
							CreatedBy:          "yyyyyyy",
							LockVersion:        ptr.To(int64(0)),
							Private:            ptr.To(false),
							APIResourceLogicID: ptr.To(int64(1231234)),
							APIResourceID:      ptr.To(int64(123456)),
							APIResourceName:    "test_api",
							ResourcePath:       "/test_api",
							APIResourceMethods: []APIResourceMethod{
								{
									APIResourceMethodID: ptr.To(int64(1231234)),
									APIResourceMethod:   APIResourceMethodsGET,
									APIParameters: []APIParameter{
										{
											APIParameterID:         ptr.To(int64(1231234)),
											APIParameterName:       "param",
											APIParameterRequired:   false,
											APIParameterLocation:   "cookie",
											APIParameterType:       "integer",
											Array:                  ptr.To(false),
											APIParamLogicID:        ptr.To(int64(1231234)),
											APIResourceMethParamID: ptr.To(int64(1231234)),
											APIChildParameters:     []APIParameter{},
										},
									},
									APIResourceMethodLogicID: ptr.To(int64(1231234)),
								},
							},
						},
					},
					EndpointDetail: EndpointDetail{
						SecurityScheme: &SecurityScheme{
							SecuritySchemeType: "apikey",
							SecuritySchemeDetail: SecuritySchemeDetail{
								APIKeyName:     "test",
								APIKeyLocation: "cookie",
							},
						},
						AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
							AllowUndefinedResources: ptr.To(restrictionsBool(true)),
							PositiveSecurityVersion: ptr.To(int64(2)),
						},
						Endpoint: Endpoint{
							CreatedBy:             "yyyyyyy",
							CreateDate:            "2022-06-23T13:08:52+0000",
							UpdatedBy:             "yyyyyyy",
							UpdateDate:            "2022-06-23T13:08:53+0000",
							APIEndpointID:         1297665,
							APIEndpointVersion:    1518543,
							APIEndpointName:       "bookstore API",
							Description:           ptr.To("testDescription"),
							ContractID:            "3-XXXXXX",
							BasePath:              "/test",
							APIEndpointScheme:     ptr.To(APIEndpointSchemeHTTP),
							ConsumeType:           ptr.To(ConsumeTypeAny),
							GroupID:               33333,
							VersionNumber:         1,
							ProtectedByAPIKey:     true,
							LockVersion:           1,
							CaseSensitive:         false,
							APIEndpointHosts:      []string{"akamai.com"},
							AvailableActions:      []string{"HIDE_ENDPOINT", "ACTIVATE_ON_STAGING", "ACTIVATE_ON_PRODUCTION", "CLONE_ENDPOINT", "EDIT_ENDPOINT_DEFINITION", "DELETE"},
							APICategoryIDs:        []int64{12345},
							APISource:             ptr.To("USER"),
							MatchPathSegmentParam: false,
							APIGatewayEnabled:     false,
							APIVersionInfo: &APIVersionInfo{
								Location:      "QUERY",
								ParameterName: "QUERY",
								Value:         "test",
							},
						},
					},
				},
			},
		},
		"fail validation": {
			body: RegisterEndpointRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "register endpoint: struct validation: APIEndpointHosts: cannot be blank\nAPIEndpointName: cannot be blank\nContractID: cannot be blank\nGroupID: cannot be blank", err.Error())
			},
		},
		"403 forbidden error": {
			expectedRequestBody: `{
    "apiEndPointName": "bookstore API",
    "contractId": "3-XXXXXX",
    "groupId": 33333,
    "apiEndPointHosts": [
        "akamai.com"
    ]
}`,
			responseBody: `{
  "detail": "string",
  "errors": [
    {}
  ],
  "instance": "https://problems.luna.akamaiapis.net/api-definitions/error-instances/d54686b5-21cb-4ab7-a8d6-a92282cf1749",
  "status": 403,
  "title": "Forbidden",
  "type": "https://problems.luna.akamaiapis.net/api-definitions/error-types/FORBIDDEN"
}`,
			responseStatus: http.StatusForbidden,
			body: RegisterEndpointRequest{
				APIEndpointName:  "bookstore API",
				APIEndpointHosts: []string{"akamai.com"},
				ContractID:       "3-XXXXXX",
				GroupID:          33333,
			},
			expectedPath: "/api-definitions/v2/endpoints",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrRegisterEndpoint), "want: %s; got: %s", ErrRegisterEndpoint, err)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, test.expectedRequestBody, string(requestBody))

				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.RegisterEndpoint(context.Background(), test.body)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func TestListUserEntitlements(t *testing.T) {
	tests := map[string]struct {
		withError      error
		responseBody   string
		expectedPath   string
		expectedResult ListUserEntitlementsResponse
		responseStatus int
	}{
		"OK": {
			responseBody: `[
  "API_READ",
  "API_WRITE",
  "API_VERSIONING",
  "API_FEATURES"
]`,
			expectedResult: ListUserEntitlementsResponse{
				"API_READ",
				"API_WRITE",
				"API_VERSIONING",
				"API_FEATURES",
			},
			expectedPath:   "/api-definitions/v2/endpoints/user-entitlements",
			responseStatus: http.StatusOK,
		},
		"500 HTTP code": {
			expectedPath:   "/api-definitions/v2/endpoints/user-entitlements",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
  "detail": "string",
  "errors": [
    {}
  ],
  "instance": "https://problems.luna.akamaiapis.net/api-definitions/error-instances/d54686b5-21cb-4ab7-a8d6-a92282cf1749",
  "status": 500,
  "title": "Internal server error",
  "type": "https://problems.luna.akamaiapis.net/api-definitions/error-types/INTERNAL-SERVER-ERROR"
}`,
			withError: &Error{
				Errors:   []Error{{}},
				Detail:   "string",
				Instance: "https://problems.luna.akamaiapis.net/api-definitions/error-instances/d54686b5-21cb-4ab7-a8d6-a92282cf1749",
				Status:   500,
				Title:    "Internal server error",
				Type:     "https://problems.luna.akamaiapis.net/api-definitions/error-types/INTERNAL-SERVER-ERROR",
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
			result, err := client.ListUserEntitlements(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func TestShowEndpoint(t *testing.T) {
	tests := map[string]struct {
		params           ShowEndpointRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ShowEndpointResponse
		withError        error
	}{
		"200 OK": {
			params:         ShowEndpointRequest{APIEndpointID: 12345},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdBy": "dl-terraform-dev (54321)",
    "createDate": "2022-08-02T10:07:42+0000",
    "updateDate": "2022-08-02T10:07:42+0000",
    "updatedBy": "dl-terraform-dev (54321)",
    "apiEndPointId": 1328556,
    "apiEndPointName": "TEST_DXE-1385",
    "description": "Test for DXE-1385",
    "basePath": "",
    "consumeType": null,
    "apiEndPointScheme": null,
    "apiEndPointVersion": 1550374,
    "contractId": "3-WNKA7W1",
    "groupId": 34567,
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
    "apiGatewayEnabled": true,
    "apiEndPointHosts": [
        "tls1.test"
    ],
    "apiCategoryIds": [],
    "source": null,
    "positiveConstrainsEnabled": false,
    "versionHidden": false,
    "endpointHidden": false,
    "matchPathSegmentParam": false,
    "availableActions": [
        "CLONE_ENDPOINT",
        "ACTIVATE_ON_STAGING",
        "HIDE_ENDPOINT",
        "ACTIVATE_ON_PRODUCTION",
        "EDIT_ENDPOINT_DEFINITION",
        "DELETE"
    ],
    "apiSource": "USER",
    "securityScheme": {
        "securitySchemeType": "apikey",
        "securitySchemeDetail": {
            "apiKeyName": "keyname",
            "apiKeyLocation": "cookie"
        }
    },
    "akamaiSecurityRestrictions": {
        "POSITIVE_SECURITY_VERSION": 2
    },
    "discoveredPiiIds": [],
    "stagingStatus": null,
    "productionStatus": null,
    "locked": false,
    "caseSensitive": true,
    "isGraphQL": false,
    "apiResources": [],
    "lockVersion": 1
}
`,
			expectedPath: "/api-definitions/v2/endpoints/12345/show",
			expectedResponse: &ShowEndpointResponse{
				EndpointWithResources: EndpointWithResources{
					APIResources:     []APIResource{},
					DiscoveredPIIIDs: []int64{},
					EndpointDetail: EndpointDetail{
						SecurityScheme: &SecurityScheme{
							SecuritySchemeType: "apikey",
							SecuritySchemeDetail: SecuritySchemeDetail{
								APIKeyName:     "keyname",
								APIKeyLocation: "cookie",
							},
						},
						AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
							PositiveSecurityVersion: ptr.To(int64(2)),
						},
						Endpoint: Endpoint{
							CreatedBy:                 "dl-terraform-dev (54321)",
							CreateDate:                "2022-08-02T10:07:42+0000",
							UpdateDate:                "2022-08-02T10:07:42+0000",
							UpdatedBy:                 "dl-terraform-dev (54321)",
							APIEndpointID:             1328556,
							APIEndpointName:           "TEST_DXE-1385",
							Description:               ptr.To("Test for DXE-1385"),
							APIEndpointVersion:        1550374,
							ContractID:                "3-WNKA7W1",
							GroupID:                   34567,
							VersionNumber:             1,
							Locked:                    false,
							StagingVersion:            VersionState{},
							ProductionVersion:         VersionState{},
							ProtectedByAPIKey:         false,
							APIGatewayEnabled:         true,
							APIEndpointHosts:          []string{"tls1.test"},
							APICategoryIDs:            []int64{},
							PositiveConstrainsEnabled: false,
							VersionHidden:             false,
							EndpointHidden:            false,
							MatchPathSegmentParam:     false,
							AvailableActions:          []string{"CLONE_ENDPOINT", "ACTIVATE_ON_STAGING", "HIDE_ENDPOINT", "ACTIVATE_ON_PRODUCTION", "EDIT_ENDPOINT_DEFINITION", "DELETE"},
							APISource:                 ptr.To("USER"),
							CaseSensitive:             true,
							IsGraphQL:                 false,
							LockVersion:               1,
						},
					},
				},
			},
		},
		"500 internal server error": {
			params:         ShowEndpointRequest{APIEndpointID: 12345},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		 "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
		 "title": "Server Error",
		 "status": 500,
		 "instance": "/api-definitions/v2/endpoints/12345/show",
		 "method": "POST",
		 "requestTime": "2025-12-06T10:27:11Z"
		}`,
			expectedPath: "/api-definitions/v2/endpoints/12345/show",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "/api-definitions/v2/endpoints/12345/show",
				Method:      ptr.To("POST"),
				RequestTime: ptr.To("2025-12-06T10:27:11Z"),
			},
		},
		"404 Not Found - EdgeWorkerID doesn't exist": {
			params:         ShowEndpointRequest{APIEndpointID: 12345},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
	"type": "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "aae04e5f-8cd6-4c06-9ed0-ec5d526aaa61",
    "detail": "Invalid endpoint provided.",
    "severity": "ERROR"
}`,
			expectedPath: "/api-definitions/v2/endpoints/12345/show",
			withError: &Error{
				Type:     "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
				Status:   404,
				Title:    "Not Found",
				Instance: "aae04e5f-8cd6-4c06-9ed0-ec5d526aaa61",
				Detail:   "Invalid endpoint provided.",
				Severity: ptr.To("ERROR"),
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
			result, err := client.ShowEndpoint(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestHideEndpoint(t *testing.T) {
	tests := map[string]struct {
		params           HideEndpointRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *HideEndpointResponse
		withError        error
	}{
		"200 OK": {
			params:         HideEndpointRequest{APIEndpointID: 12345},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdBy": "dl-terraform-dev (54321)",
    "createDate": "2022-08-02T10:07:42+0000",
    "updateDate": "2022-08-02T10:07:42+0000",
    "updatedBy": "dl-terraform-dev (54321)",
    "apiEndPointId": 1328556,
    "apiEndPointName": "TEST_DXE-1385",
    "description": "Test for DXE-1385",
    "basePath": "",
    "consumeType": null,
    "apiEndPointScheme": null,
    "apiEndPointVersion": 1550374,
    "contractId": "3-WNKA7W1",
    "groupId": 34567,
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
    "apiGatewayEnabled": true,
    "apiEndPointHosts": [
        "tls1.test"
    ],
    "apiCategoryIds": [],
    "source": null,
    "positiveConstrainsEnabled": false,
    "versionHidden": true,
    "endpointHidden": true,
    "matchPathSegmentParam": false,
    "availableActions": [
        "CLONE_ENDPOINT",
        "SHOW_ENDPOINT"
    ],
    "apiSource": "USER",
    "securityScheme": null,
    "akamaiSecurityRestrictions": {
        "POSITIVE_SECURITY_VERSION": 2
    },
    "discoveredPiiIds": [],
    "stagingStatus": null,
    "productionStatus": null,
    "locked": false,
    "caseSensitive": true,
    "isGraphQL": false,
    "apiResources": [],
    "lockVersion": 1
}
`,
			expectedPath: "/api-definitions/v2/endpoints/12345/hide",
			expectedResponse: &HideEndpointResponse{
				EndpointWithResources: EndpointWithResources{
					APIResources:     []APIResource{},
					DiscoveredPIIIDs: []int64{},
					EndpointDetail: EndpointDetail{
						AkamaiSecurityRestrictions: &AkamaiSecurityRestrictions{
							PositiveSecurityVersion: ptr.To(int64(2)),
						},
						Endpoint: Endpoint{
							CreatedBy:       "dl-terraform-dev (54321)",
							CreateDate:      "2022-08-02T10:07:42+0000",
							UpdateDate:      "2022-08-02T10:07:42+0000",
							UpdatedBy:       "dl-terraform-dev (54321)",
							APIEndpointID:   1328556,
							APIEndpointName: "TEST_DXE-1385",
							Description:     ptr.To("Test for DXE-1385"),

							APIEndpointVersion:        1550374,
							ContractID:                "3-WNKA7W1",
							GroupID:                   34567,
							VersionNumber:             1,
							Locked:                    false,
							StagingVersion:            VersionState{},
							ProductionVersion:         VersionState{},
							ProtectedByAPIKey:         false,
							APIGatewayEnabled:         true,
							APIEndpointHosts:          []string{"tls1.test"},
							APICategoryIDs:            []int64{},
							PositiveConstrainsEnabled: false,
							VersionHidden:             true,
							EndpointHidden:            true,
							MatchPathSegmentParam:     false,
							AvailableActions:          []string{"CLONE_ENDPOINT", "SHOW_ENDPOINT"},
							APISource:                 ptr.To("USER"),
							CaseSensitive:             true,
							IsGraphQL:                 false,
							LockVersion:               1,
						},
					},
				},
			},
		},
		"500 internal server error": {
			params:         HideEndpointRequest{APIEndpointID: 12345},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		 "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
		 "title": "Server Error",
		 "status": 500,
		 "instance": "/api-definitions/v2/endpoints/12345/hide",
		 "method": "POST",
		 "requestTime": "2025-12-06T10:27:11Z"
		}`,
			expectedPath: "/api-definitions/v2/endpoints/12345/hide",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "/api-definitions/v2/endpoints/12345/hide",
				Method:      ptr.To("POST"),
				RequestTime: ptr.To("2025-12-06T10:27:11Z"),
			},
		},
		"404 Not Found - EdgeWorkerID doesn't exist": {
			params:         HideEndpointRequest{APIEndpointID: 12345},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
	"type": "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "35838562-df64-4726-8fd3-89782aaaa3bb",
    "detail": "Invalid endpoint provided.",
    "severity": "ERROR"
}`,
			expectedPath: "/api-definitions/v2/endpoints/12345/hide",
			withError: &Error{
				Type:     "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
				Status:   404,
				Title:    "Not Found",
				Instance: "35838562-df64-4726-8fd3-89782aaaa3bb",
				Detail:   "Invalid endpoint provided.",
				Severity: ptr.To("ERROR"),
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
			result, err := client.HideEndpoint(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteEndpoint(t *testing.T) {
	tests := map[string]struct {
		params         DeleteEndpointRequest
		withError      error
		expectedPath   string
		responseStatus int
		responseBody   string
	}{
		"204 Deleted": {
			params:         DeleteEndpointRequest{APIEndpointID: 123},
			expectedPath:   "/api-definitions/v2/endpoints/123",
			responseStatus: http.StatusNoContent,
		},
		"403 Forbidden when endpoint does not exist": {
			params:         DeleteEndpointRequest{APIEndpointID: 123},
			expectedPath:   "/api-definitions/v2/endpoints/123",
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "/api-definitions/error-types/forbidden",
    "status": 403,
    "title": "Forbidden",
    "instance": "e368eeb7-d66a-4f9d-8476-479904555547",
    "detail": "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
    "severity": "ERROR"
}`,
			withError: &Error{
				Type:     "/api-definitions/error-types/forbidden",
				Status:   403,
				Title:    "Forbidden",
				Instance: "e368eeb7-d66a-4f9d-8476-479904555547",
				Detail:   "You have insufficient permissions to perform this action. Ensure that you have the correct permissions set in Identity and Access Management.",
				Severity: ptr.To("ERROR"),
			},
		},
		"500 internal server error": {
			params:         DeleteEndpointRequest{APIEndpointID: 123},
			expectedPath:   "/api-definitions/v2/endpoints/123",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "/api-definitions/error-types/smth",
  "title": "Server Error",
  "status": 500,
  "instance": "/api-definitions/v2/endpoints/123",
  "method": "DELETE",
  "requestTime": "2021-12-17T16:32:37Z"
}`,
			withError: &Error{
				Type:        "/api-definitions/error-types/smth",
				Title:       "Server Error",
				Status:      500,
				Instance:    "/api-definitions/v2/endpoints/123",
				Method:      ptr.To("DELETE"),
				RequestTime: ptr.To("2021-12-17T16:32:37Z"),
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
			err := client.DeleteEndpoint(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestRegisterEndpointFromFileRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		request   RegisterEndpointFromFileRequest
		withError func(t *testing.T, err error)
	}{
		"OK": {
			request: RegisterEndpointFromFileRequest{
				ContractID:        "ctr_1-1TJZH5",
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "raml",
				ImportFileSource:  "BODY_BASE64",
			},
		},
		"NOK: import file source URL + import file content ": {
			request: RegisterEndpointFromFileRequest{
				ContractID:        "ctr_1-1TJZH5",
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "raml",
				ImportFileSource:  "URL",
				ImportURL:         ptr.To("https://example.com/import"),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ImportFileContent: must not be set when ImportFileSource=='URL'", err.Error())
			},
		},
		"NOK: import file source BODY_BASE64 + no import file content": {
			request: RegisterEndpointFromFileRequest{
				ContractID:       "ctr_1-1TJZH5",
				GroupID:          44681,
				ImportFileFormat: "raml",
				ImportFileSource: "BODY_BASE64",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ImportFileContent: must be set when ImportFileSource=='BODY_BASE64'", err.Error())
			},
		},
		"NOK: no contract": {
			request: RegisterEndpointFromFileRequest{
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "raml",
				ImportFileSource:  "BODY_BASE64",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ContractID: cannot be blank", err.Error())
			},
		},
		"NOK: no group": {
			request: RegisterEndpointFromFileRequest{
				ContractID:        "ctr_1-1TJZH5",
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "raml",
				ImportFileSource:  "BODY_BASE64",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "GroupID: cannot be blank", err.Error())
			},
		},
		"NOK: no import file format": {
			request: RegisterEndpointFromFileRequest{
				ContractID:        "ctr_1-1TJZH5",
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileSource:  "BODY_BASE64",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ImportFileFormat: cannot be blank", err.Error())
			},
		},
		"NOK: no import file source": {
			request: RegisterEndpointFromFileRequest{
				ContractID:        "ctr_1-1TJZH5",
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "raml",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ImportFileSource: cannot be blank", err.Error())
			},
		},
		"NOK: bad import file source": {
			request: RegisterEndpointFromFileRequest{
				ContractID:        "ctr_1-1TJZH5",
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "raml",
				ImportFileSource:  "bad",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ImportFileSource: value 'bad' is not valid. Must be one of: 'URL', 'BODY_BASE64'", err.Error())
			},
		},
		"NOK: bad import URL": {
			request: RegisterEndpointFromFileRequest{
				ContractID:        "ctr_1-1TJZH5",
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "raml",
				ImportFileSource:  "URL",
				ImportURL:         ptr.To("BAD"),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ImportFileContent: must not be set when ImportFileSource=='URL'\nImportURL: must be a valid URL when ImportFileSource=='URL'", err.Error())
			},
		},
		"OK: import URL": {
			request: RegisterEndpointFromFileRequest{
				ContractID:       "ctr_1-1TJZH5",
				GroupID:          44681,
				ImportFileFormat: "raml",
				ImportFileSource: "URL",
				ImportURL:        ptr.To("https://www.akamai.com"),
			},
		},
		"NOK: import URL with file source BODY_BASE64": {
			request: RegisterEndpointFromFileRequest{
				ContractID:       "ctr_1-1TJZH5",
				GroupID:          44681,
				ImportFileFormat: "raml",
				ImportFileSource: "BODY_BASE64",
				ImportURL:        ptr.To("https://www.akamai.com"),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ImportFileContent: must be set when ImportFileSource=='BODY_BASE64'\nImportURL: should not be set when ImportFileSource=='BODY_BASE64'", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.request.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestRegisterEndpointFromFile(t *testing.T) {
	tests := map[string]struct {
		params              RegisterEndpointFromFileRequest
		withError           func(*testing.T, error)
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedResult      *RegisterEndpointFromFileResponse
	}{
		"does not validate": {
			params: RegisterEndpointFromFileRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "register endpoint from file: struct validation: ContractID: cannot be blank\nGroupID: cannot be blank\nImportFileFormat: cannot be blank\nImportFileSource: cannot be blank", err.Error())
			},
		},
		"OK": {
			params: RegisterEndpointFromFileRequest{
				ContractID:        "3-13H55B5",
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "swagger",
				ImportFileSource:  "BODY_BASE64",
			},
			expectedPath:        "/api-definitions/v2/endpoints/files",
			expectedRequestBody: `{"contractId":"3-13H55B5","groupId":44681,"importFileContent":"test-file-content","importFileFormat":"swagger","importFileSource":"BODY_BASE64"}`,
			responseStatus:      http.StatusCreated,
			expectedResult: &RegisterEndpointFromFileResponse{
				EndpointWithResources: EndpointWithResources{
					APIResources: []APIResource{
						{
							CreatedBy:          "bookstore_admin",
							CreateDate:         "2019-06-12T13:06:52+0000",
							UpdateDate:         "2019-06-12T13:06:52+0000",
							UpdatedBy:          "bookstore_admin",
							APIResourceID:      ptr.To(int64(2926712)),
							APIResourceName:    "books",
							ResourcePath:       "/books/{bookId}",
							Description:        "A book item within the bookstore API.",
							APIResourceLogicID: ptr.To(int64(118435)),
							LockVersion:        ptr.To(int64(2)),
							Private:            ptr.To(false),

							APIResourceMethods: []APIResourceMethod{
								{
									APIResourceMethodID:      ptr.To(int64(341591)),
									APIResourceMethod:        APIResourceMethodsGET,
									APIResourceMethodLogicID: ptr.To(int64(184404)),
									APIParameters: []APIParameter{
										{
											APIParameterID:         ptr.To(int64(1212945)),
											APIParameterRequired:   true,
											APIParameterName:       "bookId",
											APIParameterLocation:   "path",
											APIParameterType:       "string",
											APIParamLogicID:        ptr.To(int64(578116)),
											APIResourceMethParamID: ptr.To(int64(494448)),
											APIChildParameters:     []APIParameter{},
											Array:                  ptr.To(false),
										},
									},
								},
							},
						},
					},
					DiscoveredPIIIDs: []int64{},
					EndpointDetail: EndpointDetail{
						SecurityScheme: &SecurityScheme{
							SecuritySchemeType: "apikey",
							SecuritySchemeDetail: SecuritySchemeDetail{
								APIKeyName:     "apikey",
								APIKeyLocation: "header",
							},
						},
						Endpoint: Endpoint{
							CreatedBy:          "bookstore_admin",
							CreateDate:         "2019-06-12T13:06:52+0000",
							UpdateDate:         "2019-06-12T13:06:52+0000",
							UpdatedBy:          "bookstore_admin",
							APIEndpointID:      492375,
							APIEndpointName:    "Bookstore API",
							APIEndpointVersion: 1726570,
							APISource:          ptr.To("USER"),
							APIGatewayEnabled:  true,
							Description:        ptr.To("An API for bookstore users allowing them to retrieve book items, add new items (admin users), and modify existing items."),
							BasePath:           "/bookstore",
							ConsumeType:        ptr.To(ConsumeTypeAny),
							APIEndpointScheme:  ptr.To(APIEndpointSchemeHTTPHTTPS),
							ContractID:         "3-13H55B5",
							GroupID:            int64(44681),
							VersionNumber:      1,
							ClonedFromVersion:  ptr.To(int64(1)),
							ProtectedByAPIKey:  true,
							CaseSensitive:      true,
							LockVersion:        int64(2),
							APIEndpointHosts:   []string{"bookstore.api.akamai.com"},
							APICategoryIDs:     []int64{2, 7},
							AvailableActions: []string{"DELETE",
								"CLONE_ENDPOINT",
								"ACTIVATE_ON_PRODUCTION",
								"HIDE_ENDPOINT",
								"EDIT_ENDPOINT_DEFINITION",
								"ACTIVATE_ON_STAGING"},
							StagingVersion:    VersionState{},
							ProductionVersion: VersionState{},
						},
					},
				},
			},
			responseBody: `{
  "createdBy": "bookstore_admin",
  "createDate": "2019-06-12T13:06:52+0000",
  "updateDate": "2019-06-12T13:06:52+0000",
  "updatedBy": "bookstore_admin",
  "apiEndPointId": 492375,
  "apiEndPointName": "Bookstore API",
  "description": "An API for bookstore users allowing them to retrieve book items, add new items (admin users), and modify existing items.",
  "basePath": "/bookstore",
  "consumeType": "any",
  "apiEndPointScheme": "http/https",
  "apiEndPointVersion": 1726570,
  "contractId": "3-13H55B5",
  "groupId": 44681,
  "versionNumber": 1,
  "clonedFromVersion": 1,
  "apiEndPointLocked": false,
  "protectedByApiKey": true,
  "source": null,
  "positiveConstrainsEnabled": null,
  "versionHidden": false,
  "endpointHidden": false,
  "akamaiSecurityRestrictions": null,
  "discoveredPiiIds": [],
  "stagingStatus": null,
  "productionStatus": null,
  "caseSensitive": true,
  "matchPathSegmentParam": false,
  "apiSource": "USER",
  "apiGatewayEnabled": true,
  "locked": false,
  "isGraphQL": false,
  "lockVersion": 2,
  "apiEndPointHosts": [
    "bookstore.api.akamai.com"
  ],
  "apiCategoryIds": [
    2,
    7
  ],
  "availableActions": [
    "DELETE",
    "CLONE_ENDPOINT",
    "ACTIVATE_ON_PRODUCTION",
    "HIDE_ENDPOINT",
    "EDIT_ENDPOINT_DEFINITION",
    "ACTIVATE_ON_STAGING"
  ],
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
  "securityScheme": {
    "securitySchemeType": "apikey",
    "securitySchemeDetail": {
      "apiKeyLocation": "header",
      "apiKeyName": "apikey"
    }
  },
  "apiResources": [
    {
      "createdBy": "bookstore_admin",
      "createDate": "2019-06-12T13:06:52+0000",
      "updateDate": "2019-06-12T13:06:52+0000",
      "updatedBy": "bookstore_admin",
      "apiResourceId": 2926712,
      "apiResourceName": "books",
      "resourcePath": "/books/{bookId}",
      "description": "A book item within the bookstore API.",
      "link": null,
      "apiResourceClonedFromId": null,
      "apiResourceLogicId": 118435,
      "private": false,
      "lockVersion": 2,
      "apiResourceMethods": [
        {
          "apiResourceMethodId": 341591,
          "apiResourceMethod": "GET",
          "apiResourceMethodLogicId": 184404,
          "apiParameters": [
            {
              "apiParameterId": 1212945,
              "apiParameterRequired": true,
              "apiParameterName": "bookId",
              "apiParameterLocation": "path",
              "apiParameterType": "string",
              "array": false,
              "pathParamLocationId": null,
              "apiParamLogicId": 578116,
              "apiResourceMethParamId": 494448,
              "apiParameterRestriction": null,
              "apiParameterNotes": null,
              "apiChildParameters": []
            }
          ]
        }
      ]
    }
  ]
}`,
		},
		"NOK: an error": {
			params: RegisterEndpointFromFileRequest{
				ContractID:        "3-13H55B5",
				GroupID:           44681,
				ImportFileContent: ptr.To(registerEndpointFileContent),
				ImportFileFormat:  "raml",
				ImportFileSource:  "BODY_BASE64",
			},
			withError: func(t *testing.T, err error) {
				assert.Regexp(t, ErrRegisterEndpointFromFile, err.Error())
			},
			responseStatus: http.StatusServiceUnavailable,
			expectedPath:   "/api-definitions/v2/endpoints/files",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			response, err := client.RegisterEndpointFromFile(context.Background(), test.params)
			if test.withError != nil {
				require.Error(t, err)
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResult, response)
		})
	}
}
