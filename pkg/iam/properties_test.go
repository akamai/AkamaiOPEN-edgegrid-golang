package iam

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProperties(t *testing.T) {
	tests := map[string]struct {
		params           ListPropertiesRequest
		responseStatus   int
		expectedPath     string
		responseBody     string
		expectedResponse *ListPropertiesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - no query params": {
			params:         ListPropertiesRequest{},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/user-admin/properties?actions=false",
			responseBody: `
[
    {
        "propertyId": 1,
        "propertyName": "property1",
        "propertyTypeDescription": "Site",
        "groupId": 11,
        "groupName": "group1"
    },
    {
        "propertyId": 2,
        "propertyName": "property2",
        "propertyTypeDescription": "Site",
        "groupId": 22,
        "groupName": "group2"
    }
]
`,
			expectedResponse: &ListPropertiesResponse{
				{
					PropertyID:              1,
					PropertyName:            "property1",
					PropertyTypeDescription: "Site",
					GroupID:                 11,
					GroupName:               "group1",
				},
				{
					PropertyID:              2,
					PropertyName:            "property2",
					PropertyTypeDescription: "Site",
					GroupID:                 22,
					GroupName:               "group2",
				},
			},
		},
		"200 OK - with query params": {
			params: ListPropertiesRequest{
				Actions: true,
				GroupID: 12345,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/user-admin/properties?actions=true&groupId=12345",
			responseBody: `
[
    {
        "propertyId": 1,
        "propertyName": "property1",
        "propertyTypeDescription": "Site",
        "groupId": 12345,
        "groupName": "group1",
		"actions": {
			"move": false
		}
    }
]
`,
			expectedResponse: &ListPropertiesResponse{
				{
					PropertyID:              1,
					PropertyName:            "property1",
					PropertyTypeDescription: "Site",
					GroupID:                 12345,
					GroupName:               "group1",
					Actions: PropertyActions{
						Move: false,
					},
				},
			},
		},
		"200 OK - no properties": {
			params:           ListPropertiesRequest{},
			responseStatus:   http.StatusOK,
			expectedPath:     "/identity-management/v3/user-admin/properties?actions=false",
			responseBody:     `[]`,
			expectedResponse: &ListPropertiesResponse{},
		},
		"500 internal server error": {
			params:         ListPropertiesRequest{},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/identity-management/v3/user-admin/properties?actions=false",
			responseBody: `
					{
						"type": "internal_error",
						"title": "Internal Server Error",
						"detail": "Error processing request",
						"status": 500
					}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error processing request",
					StatusCode: http.StatusInternalServerError,
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
			users, err := client.ListProperties(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}

func TestGetProperty(t *testing.T) {
	tests := map[string]struct {
		params           GetPropertyRequest
		responseStatus   int
		expectedPath     string
		responseBody     string
		expectedResponse *GetPropertyResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetPropertyRequest{
				PropertyID: 1,
				GroupID:    11,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/user-admin/properties/1?groupId=11",
			responseBody: `
{
    "createdDate": "2023-08-18T09:10:37.000Z",
    "createdBy": "user1",
    "modifiedDate": "2023-08-18T09:10:37.000Z",
    "modifiedBy": "user2",
    "groupName": "group1",
    "groupId": 11,
    "arlConfigFile": "test.xml",
    "propertyId": 1,
    "propertyName": "name1"
}
`,
			expectedResponse: &GetPropertyResponse{
				ARLConfigFile: "test.xml",
				CreatedBy:     "user1",
				CreatedDate:   test.NewTimeFromString(t, "2023-08-18T09:10:37.000Z"),
				GroupID:       11,
				GroupName:     "group1",
				ModifiedBy:    "user2",
				ModifiedDate:  test.NewTimeFromString(t, "2023-08-18T09:10:37.000Z"),
				PropertyID:    1,
				PropertyName:  "name1",
			},
		},
		"validation errors": {
			params: GetPropertyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get property: struct validation:\nGroupID: cannot be blank\nPropertyID: cannot be blank", err.Error())
			},
		},
		"404 not found": {
			params: GetPropertyRequest{
				PropertyID: 1,
				GroupID:    11,
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/identity-management/v3/user-admin/properties/1?groupId=11",
			responseBody: `
{
	"instance": "",
	"httpStatus": 404,
	"detail": "",
	"title": "Property not found",
	"type": "/useradmin-api/error-types/1806"
}					
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "/useradmin-api/error-types/1806",
					Title:      "Property not found",
					StatusCode: http.StatusNotFound,
					HTTPStatus: http.StatusNotFound,
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
			users, err := client.GetProperty(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}

func TestMoveProperty(t *testing.T) {
	tests := map[string]struct {
		params              MovePropertyRequest
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
	}{
		"204 OK": {
			params: MovePropertyRequest{
				PropertyID: 1,
				BodyParams: MovePropertyReqBody{
					DestinationGroupID: 22,
					SourceGroupID:      11,
				},
			},
			expectedRequestBody: `
{
	"destinationGroupId": 22,
	"sourceGroupId": 11
}`,
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v3/user-admin/properties/1",
		},
		"validation errors": {
			params: MovePropertyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "move property: struct validation: BodyParams: DestinationGroupID: cannot be blank\nSourceGroupID: cannot be blank\nPropertyID: cannot be blank", err.Error())
			},
		},
		"400 not allowed": {
			params: MovePropertyRequest{
				PropertyID: 1,
				BodyParams: MovePropertyReqBody{
					DestinationGroupID: 22,
					SourceGroupID:      11,
				},
			},
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/identity-management/v3/user-admin/properties/1",
			responseBody: `
{
    "instance": "",
    "httpStatus": 400,
    "detail": "Property move is not allowed from the group 11",
    "title": "Validation Exception",
    "type": "/useradmin-api/error-types/1806"
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "/useradmin-api/error-types/1806",
					Title:      "Validation Exception",
					Detail:     "Property move is not allowed from the group 11",
					HTTPStatus: http.StatusBadRequest,
					StatusCode: http.StatusBadRequest,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if test.responseBody != "" {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.MoveProperty(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestMapPropertyIDToName(t *testing.T) {
	tests := map[string]struct {
		params           MapPropertyIDToNameRequest
		responseStatus   int
		responseBody     string
		expectedResponse *string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: MapPropertyIDToNameRequest{
				PropertyID: 1,
				GroupID:    11,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdDate": "2023-08-18T09:10:37.000Z",
    "createdBy": "user1",
    "modifiedDate": "2023-08-18T09:10:37.000Z",
    "modifiedBy": "user2",
    "groupName": "group1",
    "groupId": 11,
    "arlConfigFile": "test.xml",
    "propertyId": 1,
    "propertyName": "name1"
}
`,
			expectedResponse: ptr.To("name1"),
		},
		"validation errors": {
			params: MapPropertyIDToNameRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "map property by id: struct validation:\nGroupID: cannot be blank\nPropertyID: cannot be blank", err.Error())
			},
		},
		"404 not found": {
			params: MapPropertyIDToNameRequest{
				PropertyID: 1,
				GroupID:    11,
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
		{
			"instance": "",
			"httpStatus": 404,
			"detail": "",
			"title": "Property not found",
			"type": "/useradmin-api/error-types/1806"
		}
		`,
			withError: func(t *testing.T, err error) {
				assert.Equal(t, err.Error(), `map property by id: request failed: get property: API error: 
{
	"type": "/useradmin-api/error-types/1806",
	"title": "Property not found",
	"detail": "",
	"statusCode": 404,
	"httpStatus": 404
}`)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			users, err := client.MapPropertyIDToName(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}

func TestMapPropertyNameToID(t *testing.T) {
	listPropertiesResponse := `
[
    {
        "propertyId": 1,
        "propertyName": "name1",
        "propertyTypeDescription": "Site",
        "groupId": 11,
        "groupName": "group1"
    },
    {
        "propertyId": 2,
        "propertyName": "name2",
        "propertyTypeDescription": "Site",
        "groupId": 22,
        "groupName": "group2"
    }
]
`
	tests := map[string]struct {
		name             MapPropertyNameToIDRequest
		responseStatus   int
		responseBody     string
		expectedResponse *int64
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			name:             "name2",
			responseStatus:   http.StatusOK,
			responseBody:     listPropertiesResponse,
			expectedResponse: ptr.To(int64(2)),
		},
		"200 but not found": {
			name:           "name3",
			responseStatus: http.StatusOK,
			responseBody:   listPropertiesResponse,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrNoProperty))
				assert.Equal(t, "no such property: name3", err.Error())
			},
		},
		"validation errors": {
			name: "",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "map property by name: struct validation:\n name cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			name:           "name2",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
					{
						"type": "internal_error",
						"title": "Internal Server Error",
						"detail": "Error processing request",
						"status": 500
					}`,
			withError: func(t *testing.T, err error) {
				assert.Equal(t, err.Error(), `map property by name: request failed: list properties: API error: 
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"detail": "Error processing request",
	"statusCode": 500
}`)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			users, err := client.MapPropertyNameToID(context.Background(), test.name)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}
