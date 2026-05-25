package apidefinitions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchResourceOperations(t *testing.T) {
	tests := map[string]struct {
		expectedPath     string
		expectedResponse *SearchResourceOperationsResponse
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			expectedPath:     "/api-definitions/v2/search-operations",
			expectedResponse: loadJson[SearchResourceOperationsResponse]("testdata/search_resource_operations.json"),
			responseStatus:   http.StatusOK,
			responseBody: `{
   "apiEndPoints":[
      {
         "apiEndPointHosts":[
            "www.bot-staging.net",
            "www.botman-test1.com",
            "www.botman-test.com",
            "www.test-login.com"
         ],
         "apiEndPointId":503130,
         "apiEndPointName":"E2ETest",
         "basePath":"/login",
         "caseSensitive":true,
         "link":"/api-definitions/v2/endpoints/503130/versions/4",
         "productionVersion":{
            "status":"ACTIVE",
            "timestamp":"2019-10-02T16:04:14+0000",
            "versionNumber":4
         },
         "stagingVersion":{
            "status":"ACTIVE",
            "timestamp":"2019-10-02T15:59:07+0000",
            "versionNumber":4
         }
      }
   ],
   "operations":[
      {
         "apiEndPointId":503130,
         "apiResourceId":3143802,
         "apiResourceLogicId":122202,
         "conditions":[
            {
               "apiParameterId":1636906,
               "value":"cors/form-login"
            }
         ],
         "link":"/api-definitions/v2/endpoints/503130/versions/4/resources/3143802/operations/123a567b-fe5c-4d52-9cc9-a9cd870f0829",
         "metadata":{
            "isActive":true
         },
         "method":"POST",
         "operationId":"123a567b-fe5c-4d52-9cc9-a9cd870f0829",
         "operationName":"Form Data",
         "operationPurpose":"ACCOUNT_VERIFICATION"
      }
   ],
   "resources":[
      {
         "apiEndPointId":503130,
         "apiResourceId":3143802,
         "apiResourceLogicId":122202,
         "apiResourceMethods":[
            {
               "apiParameters":[
                  {
                     "apiChildParameters":[
                        
                     ],
                     "apiParamLogicId":778229,
                     "apiParameterId":1636906,
                     "apiParameterLocation":"query",
                     "apiParameterName":"route",
                     "apiParameterNotes":null,
                     "apiParameterRequired":true,
                     "apiParameterRestriction":null,
                     "apiParameterType":"string",
                     "apiResourceMethParamId":572570,
                     "array":false,
                     "pathParamLocationId":null,
                     "response":false
                  }
               ],
               "apiResourceMethod":"POST",
               "apiResourceMethodId":365993,
               "apiResourceMethodLogicId":170106,
               "methodRestrictions":null
            }
         ],
         "apiResourceName":"Login Form Data",
         "createDate":"2019-10-02T15:39:03+0000",
         "createdBy":"jsmith",
         "link":"/api-definitions/v2/endpoints/503130/versions/4/resources/3143802",
         "lockVersion":0,
         "metadata":{
            "methodsEnabled":1,
            "methodsWithOperations":1,
            "operationCount":3
         },
         "resourcePath":"/index.php",
         "updateDate":"2019-10-02T15:39:03+0000",
         "updatedBy":"jsmith"
      }
   ]
}`,
		},
		"403 Unauthorized Access/Action": {
			expectedPath:     "/api-definitions/v2/search-operations",
			expectedResponse: loadJson[SearchResourceOperationsResponse]("testdata/search_resource_operations_403.json"),
			responseStatus:   http.StatusForbidden,
			responseBody: `{
						  "detail": "You do not have the necessary access to perform this operation.",
						  "instance": "https://problems.luna.akamaiapis.net/appsec/error-instances/f1b952dbf3a940a9",
						  "status": 403,
						  "title": "Unauthorized Access/Action",
						  "type": "https://problems.luna.akamaiapis.net/appsec/error-types/UNAUTHORIZED"
						}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "https://problems.luna.akamaiapis.net/appsec/error-types/UNAUTHORIZED",
					Title:    "Unauthorized Access/Action",
					Detail:   "You do not have the necessary access to perform this operation.",
					Instance: "https://problems.luna.akamaiapis.net/appsec/error-instances/f1b952dbf3a940a9",
					Status:   403,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 Internal Server Error": {
			expectedPath:     "/api-definitions/v2/search-operations",
			expectedResponse: loadJson[SearchResourceOperationsResponse]("testdata/search_resource_operations_500.json"),
			responseStatus:   http.StatusInternalServerError,
			responseBody: `{
							  "detail": "An unexpected error occurred on the server.",
							  "instance": "https://problems.luna.akamaiapis.net/appsec/error-instances/abc123def456",
							  "status": 500,
							  "title": "Internal Server Error",
							  "type": "https://problems.luna.akamaiapis.net/appsec/error-types/INTERNAL_SERVER_ERROR"
							}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "https://problems.luna.akamaiapis.net/appsec/error-types/INTERNAL_SERVER_ERROR",
					Title:    "Internal Server Error",
					Detail:   "An unexpected error occurred on the server.",
					Instance: "https://problems.luna.akamaiapis.net/appsec/error-instances/abc123def456",
					Status:   500,
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.SearchResourceOperations(context.Background())
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func loadJson[T any](path string) *T {
	contents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var result T
	if err := json.Unmarshal(contents, &result); err != nil {
		panic(fmt.Sprintf("failed to unmarshal JSON: %s", err))
	}
	return &result
}
