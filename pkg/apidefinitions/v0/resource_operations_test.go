package v0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TestUpdateResourceOperation(t *testing.T) {
	tests := map[string]struct {
		params              UpdateResourceOperationRequest
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *UpdateResourceOperationResponse
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
	}{
		"201 OK Doc Sample": {
			params: UpdateResourceOperationRequest{
				VersionNumber: 1,
				APIID:         123,
				Body:          ResourceOperationsRequestBody(*updateResourceOperationResponse()),
			},
			expectedPath:        "/api-definitions/v0/endpoints/123/versions/1/operations",
			expectedRequestBody: loadJson("testdata/resource_operations_01.json"),
			expectedResponse:    updateResourceOperationResponse(),
			responseStatus:      http.StatusOK,
			responseBody:        loadJson("testdata/resource_operations_01.json"),
		},
		"201 OK Empty": {
			params: UpdateResourceOperationRequest{
				VersionNumber: 1,
				APIID:         123,
				Body:          ResourceOperationsRequestBody{},
			},
			expectedPath:        "/api-definitions/v0/endpoints/123/versions/1/operations",
			expectedRequestBody: loadJson("testdata/resource_operations_empty.json"),
			expectedResponse:    &UpdateResourceOperationResponse{},
			responseStatus:      http.StatusOK,
			responseBody:        loadJson("testdata/resource_operations_empty.json"),
		},
		"400 Bad Request": {
			params: UpdateResourceOperationRequest{
				VersionNumber: 1,
				APIID:         123,
				Body:          ResourceOperationsRequestBody(*badRequest()),
			},
			expectedPath:        "/api-definitions/v0/endpoints/123/versions/1/operations",
			expectedRequestBody: loadJson("testdata/resource_operations_bad_request.json"),
			responseStatus:      http.StatusBadRequest,
			responseBody:        loadJson("testdata/400_bad_request_api_operations.json"),
			withError: func(t *testing.T, err error) {
				want := loadErrorResponse()
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"404 not found": {
			params: UpdateResourceOperationRequest{
				APIID:         1,
				VersionNumber: 2,
				Body:          ResourceOperationsRequestBody{},
			},
			expectedRequestBody: `
									{"operations":null}
								`,
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2/operations",
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
		"validation errors": {
			params: UpdateResourceOperationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update resource operations: struct validation: APIID: cannot be blank\nVersionNumber: cannot be blank", err.Error())
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
			result, err := client.UpdateResourceOperation(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetResourceOperation(t *testing.T) {
	tests := map[string]struct {
		params           GetResourceOperationRequest
		expectedPath     string
		expectedResponse *GetResourceOperationResponse
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetResourceOperationRequest{
				APIID:         11,
				VersionNumber: 22,
			},
			expectedPath:     "/api-definitions/v0/endpoints/11/versions/22/operations",
			expectedResponse: &GetResourceOperationResponse{},
			responseStatus:   http.StatusOK,
			responseBody: `
			{
			  "operations": null
			}
			`,
			withError: nil,
		},
		"200 OK all fields populated": {
			params: GetResourceOperationRequest{
				APIID:         456,
				VersionNumber: 123,
			},
			expectedPath:     "/api-definitions/v0/endpoints/456/versions/123/operations",
			expectedResponse: getResourceOperationResponse(),
			responseStatus:   http.StatusOK,
			responseBody:     loadJson("testdata/resource_operations_02.json"),
			withError:        nil,
		},
		"403 forbidden": {
			params: GetResourceOperationRequest{
				APIID:         1,
				VersionNumber: 2,
			},
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2/operations",
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
			params: GetResourceOperationRequest{
				APIID:         1,
				VersionNumber: 1010,
			},
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/1010/operations",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "test.com/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "TestInstance123",
    "detail": "No Api Endpoint/Version found for endpoint ID 1 and version 1010",
    "severity": "ERROR",
    "stackTrace": "StackTraceTest"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "test.com/api-definitions/error-types/NOT-FOUND",
					Title:    "Not Found",
					Detail:   "No Api Endpoint/Version found for endpoint ID 1 and version 1010",
					Instance: "TestInstance123",
					Status:   http.StatusNotFound,
					Severity: ptr.To("ERROR"),
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"required param not provided": {
			params: GetResourceOperationRequest{
				VersionNumber: 11,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get resource operations: struct validation: APIID: cannot be blank", err.Error())
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
			result, err := client.GetResourceOperation(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}

}

func TestDeleteResourceOperation(t *testing.T) {
	tests := map[string]struct {
		params              DeleteResourceOperationRequest
		expectedPath        string
		expectedResponse    *DeleteResourceOperationResponse
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
		expectedRequestBody string
	}{
		"200 OK": {
			params: DeleteResourceOperationRequest{
				APIID:         11,
				VersionNumber: 11,
			},
			expectedPath:        "/api-definitions/v0/endpoints/11/versions/11/operations",
			expectedRequestBody: loadJson("testdata/resource_operations_empty.json"),
			expectedResponse:    &deleteResourceOperationResponse,
			responseStatus:      http.StatusOK,
			withError:           nil,
			responseBody:        loadJson("testdata/delete_Resource_Operations_Response.json"),
		},
		"required params not present": {
			params: DeleteResourceOperationRequest{
				VersionNumber: 11,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete resource operations: struct validation: APIID: cannot be blank", err.Error())
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteResourceOperation(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func getResourceOperationResponse() *GetResourceOperationResponse {
	resourceOperations := orderedmap.New[string, *orderedmap.OrderedMap[string, Operation]]()

	operationsForIndex := orderedmap.New[string, Operation]()

	operationsForIndex.Set("testPurpose", Operation{
		Method:  ptr.To("POST"),
		Purpose: ptr.To("login"),
		Parameters: func() *orderedmap.OrderedMap[string, OperationParameter] {
			omap := orderedmap.New[string, OperationParameter]()
			omap.Set("us979898ername", OperationParameter{
				Path:         []string{"username"},
				Location:     ptr.To("request_body"),
				UsedForLogin: ptr.To(false),
			})
			return omap
		}(),
		FailureConditions: []OperationCondition{
			{
				HeaderName:                 ptr.To("Content-Length"),
				PositiveMatch:              ptr.To(true),
				SuppressFromClientResponse: ptr.To(false),
				Type:                       ptr.To("header_value"),
				ValueCase:                  ptr.To(false),
				ValueWildcard:              ptr.To(true),
				Path:                       ptr.To(""),
				Values:                     []string{"28?"},
			},
		},
		SuccessConditions: []OperationCondition{
			{
				HeaderName:                 ptr.To("Content-Type"),
				PositiveMatch:              ptr.To(true),
				SuppressFromClientResponse: ptr.To(false),
				Type:                       ptr.To(""),
				ValueCase:                  ptr.To(false),
				ValueWildcard:              ptr.To(true),
				Path:                       ptr.To(""),
				Values:                     []string{},
			},
		},
		Conditions: []ParameterPathCondition{
			{
				Path:          []string{"root", "username"},
				Location:      ptr.To("request_body"),
				PositiveMatch: ptr.To(true),
			},
		},
	})

	resourceOperations.Set("/index.php*", operationsForIndex)

	operationsForLogin := orderedmap.New[string, Operation]()

	operationsForLogin.Set("test-login", Operation{
		Method:             ptr.To("POST"),
		Purpose:            ptr.To("login"),
		MultistepGroupName: ptr.To("msg-group-test"),
		Parameters: func() *orderedmap.OrderedMap[string, OperationParameter] {
			omap := orderedmap.New[string, OperationParameter]()
			omap.Set("username", OperationParameter{
				ReferenceParameterLocation: ptr.To("USER_EMAIL"),
			})
			return omap
		}(),
		SuccessConditions: []OperationCondition{
			{
				PositiveMatch: ptr.To(true),
				Type:          ptr.To("HTTP_STATUS"),
				Values:        []string{"200"},
			},
		},
		StepSuccessConditions: []OperationCondition{
			{
				PositiveMatch: ptr.To(true),
				Type:          ptr.To("HTTP_STATUS"),
				Values:        []string{"200"},
			},
		},
	})

	resourceOperations.Set("/login", operationsForLogin)

	var getResourceOperationResponse = GetResourceOperationResponse{
		ResourceOperations: resourceOperations,
	}

	return &getResourceOperationResponse
}

func updateResourceOperationResponse() *UpdateResourceOperationResponse {
	resourceOperations := orderedmap.New[string, *orderedmap.OrderedMap[string, Operation]]()

	operationsForResource := orderedmap.New[string, Operation]()

	operationsForResource.Set("UpdateOperation", Operation{
		Method:  ptr.To("POST"),
		Purpose: ptr.To("login"),
		Parameters: func() *orderedmap.OrderedMap[string, OperationParameter] {
			omap := orderedmap.New[string, OperationParameter]()
			omap.Set("param1", OperationParameter{ // Fixed: Using OrderedMap
				Path:         []string{"path", "to", "param1"},
				Location:     ptr.To("body"),
				UsedForLogin: ptr.To(true),
			})
			return omap
		}(),
		FailureConditions: []OperationCondition{
			{
				HeaderName:    ptr.To("X-Error"),
				PositiveMatch: ptr.To(false),
				Type:          ptr.To("ErrorCondition"),
				Values:        []string{"400", "500"},
			},
		},
		SuccessConditions: []OperationCondition{
			{
				HeaderName:    ptr.To("X-Success"),
				PositiveMatch: ptr.To(true),
				Type:          ptr.To("SuccessCondition"),
				Values:        []string{"201"},
			},
		},
	})

	resourceOperations.Set("/api/resource", operationsForResource)

	updateResourceOperationResponse := UpdateResourceOperationResponse{
		ResourceOperations: resourceOperations,
	}

	return &updateResourceOperationResponse
}

func badRequest() *UpdateResourceOperationResponse {
	resourceOps := orderedmap.New[string, *orderedmap.OrderedMap[string, Operation]]()
	operationsForBase := orderedmap.New[string, Operation]()

	operationsForBase.Set("test login", Operation{
		Method:  ptr.To("POST"),
		Purpose: ptr.To("login"),
	})

	resourceOps.Set("/base", operationsForBase)

	updateResourceOperationResponse := UpdateResourceOperationResponse{
		ResourceOperations: resourceOps,
	}

	return &updateResourceOperationResponse
}

var deleteResourceOperationResponse = DeleteResourceOperationResponse{
	APIID:         11,
	VersionNumber: 11,
	Status:        200,
	Detail:        "Api resource operations for Endpoint is Deleted",
}

func loadErrorResponse() *Error {
	jsonStr := `
	{
		"type": "/api-definitions/error-types/invalid-input-error",
		"title": "Invalid input error",
		"detail": "The request you submitted is invalid. Modify the request and try again.",
		"instance": "id_001",
		"status": 400,
		"severity": "ERROR",
		"errors": [
			{
				"type": "/api-definitions/error-types/resource-path-operation-check",
				"title": "resource-path-operation-check.title",
				"detail": "resource-path-operation-check.detail",
				"severity": "ERROR",
				"field": "put.dto.resourceOperations[/base].\u003cmap value\u003e[test login].operationParameter",
				"rejectedValue": {
        	            			"method": "POST",
        	            			"operationPurpose": "login"
					}
			},
			{
				"type": "/api-definitions/error-types/resource-path-operation-check",
				"title": "resource-path-operation-check.title",
				"detail": "resource-path-operation-check.detail",
				"severity": "ERROR",
				"field": "put.dto.resourceOperations[/base].\u003cmap value\u003e[test login].operationParameter.username",
				"rejectedValue": {
        	            			"method": "POST",
        	            			"operationPurpose": "login"
					}
			}
		]
	}`

	var apiError Error
	err := json.Unmarshal([]byte(jsonStr), &apiError)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return &apiError
}
