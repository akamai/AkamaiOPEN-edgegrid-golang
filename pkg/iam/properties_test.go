package iam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/internal/test"
	"github.com/stretchr/testify/assert"
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
