package datastream

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDs_GetProperties(t *testing.T) {
	tests := map[string]struct {
		request          GetPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Property
		withError        error
	}{
		"200 OK": {
			request: GetPropertiesRequest{
				GroupId:   12345,
				ProductId: "Download_Delivery",
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "propertyId": 382631,
        "propertyName": "customp.akamai.com",
        "productId": "Ion_Standard",
        "productName": "Ion Standard",
        "hostnames": [
            "customp.akamaize.net",
            "customp.akamaized-staging.net"
        ]
    },
    {
        "propertyId": 347459,
        "propertyName": "example.com",
        "productId": "Dynamic_Site_Accelerator",
        "productName": "Dynamic Site Accelerator",
        "hostnames": [
            "example.edgekey.net"
        ]
    }
]
`,
			expectedPath: "/datastream-config-api/v1/log/properties/product/Download_Delivery/group/12345",
			expectedResponse: []Property{
				{
					PropertyID:   382631,
					PropertyName: "customp.akamai.com",
					ProductID:    "Ion_Standard",
					ProductName:  "Ion Standard",
					Hostnames: []string{
						"customp.akamaize.net",
						"customp.akamaized-staging.net",
					},
				},
				{
					PropertyID:   347459,
					PropertyName: "example.com",
					ProductID:    "Dynamic_Site_Accelerator",
					ProductName:  "Dynamic Site Accelerator",
					Hostnames: []string{
						"example.edgekey.net",
					},
				},
			},
		},
		"validation error": {
			request:   GetPropertiesRequest{},
			withError: ErrStructValidation,
		},
		"400 bad request": {
			request:        GetPropertiesRequest{GroupId: 12345, ProductId: "testProductName"},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
	"type": "bad-request",
	"title": "Bad Request",
	"detail": "",
	"instance": "baf2671f-7b3a-406d-9dd8-63ef20a01296",
	"statusCode": 400,
	"errors": [
		{
			"type": "bad-request",
			"title": "Bad Request",
			"detail": "Invalid Product Name"
		}
	]
}
`,
			expectedPath: "/datastream-config-api/v1/log/properties/product/testProductName/group/12345",
			withError: &Error{
				Type:       "bad-request",
				Title:      "Bad Request",
				Instance:   "baf2671f-7b3a-406d-9dd8-63ef20a01296",
				StatusCode: http.StatusBadRequest,
				Errors: []RequestErrors{
					{
						Type:   "bad-request",
						Title:  "Bad Request",
						Detail: "Invalid Product Name",
					},
				},
			},
		},
		"403 forbidden": {
			request:        GetPropertiesRequest{GroupId: 12345, ProductId: "testProductName"},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
	"type": "forbidden",
	"title": "Forbidden",
	"detail": "",
	"instance": "28eb43a8-97ae-4c57-98aa-258081582b92",
	"statusCode": 403,
	"errors": [
		{
			"type": "forbidden",
			"title": "Forbidden",
			"detail": "User is not having access for the group. Access denied, please contact support."
		}
	]
}
`,
			expectedPath: "/datastream-config-api/v1/log/properties/product/testProductName/group/12345",
			withError: &Error{
				Type:       "forbidden",
				Title:      "Forbidden",
				Instance:   "28eb43a8-97ae-4c57-98aa-258081582b92",
				StatusCode: http.StatusForbidden,
				Errors: []RequestErrors{
					{
						Type:   "forbidden",
						Title:  "Forbidden",
						Detail: "User is not having access for the group. Access denied, please contact support.",
					},
				},
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
			result, err := client.GetProperties(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_GetPropertiesByGroup(t *testing.T) {
	tests := map[string]struct {
		request          GetPropertiesByGroupRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Property
		withError        error
	}{
		"200 OK": {
			request: GetPropertiesByGroupRequest{
				GroupId: 12345,
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "propertyId": 382631,
        "propertyName": "customp.akamai.com",
        "productId": "Ion_Standard",
        "productName": "Ion Standard",
        "hostnames": [
            "customp.akamaize.net",
            "customp.akamaized-staging.net"
        ]
    },
    {
        "propertyId": 347459,
        "propertyName": "example.com",
        "productId": "Dynamic_Site_Accelerator",
        "productName": "Dynamic Site Accelerator",
        "hostnames": [
            "example.edgekey.net"
        ]
    }
]
`,
			expectedPath: "/datastream-config-api/v1/log/properties/group/12345",
			expectedResponse: []Property{
				{
					PropertyID:   382631,
					PropertyName: "customp.akamai.com",
					ProductID:    "Ion_Standard",
					ProductName:  "Ion Standard",
					Hostnames: []string{
						"customp.akamaize.net",
						"customp.akamaized-staging.net",
					},
				},
				{
					PropertyID:   347459,
					PropertyName: "example.com",
					ProductID:    "Dynamic_Site_Accelerator",
					ProductName:  "Dynamic Site Accelerator",
					Hostnames: []string{
						"example.edgekey.net",
					},
				},
			},
		},
		"validation error": {
			request:   GetPropertiesByGroupRequest{},
			withError: ErrStructValidation,
		},
		"403 access forbidden": {
			request:        GetPropertiesByGroupRequest{GroupId: 12345},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
	"type": "forbidden",
	"title": "Forbidden",
	"detail": "",
	"instance": "04fde003-428b-4c2c-94fe-6109af9d231c",
	"statusCode": 403,
	"errors": [
		{
			"type": "forbidden",
			"title": "Forbidden",
			"detail": "User is not having access for the group. Access denied, please contact support."
		}
	]
}
`,
			expectedPath: "/datastream-config-api/v1/log/properties/group/12345",
			withError: &Error{
				Type:       "forbidden",
				Title:      "Forbidden",
				Instance:   "04fde003-428b-4c2c-94fe-6109af9d231c",
				StatusCode: http.StatusForbidden,
				Errors: []RequestErrors{
					{
						Type:   "forbidden",
						Title:  "Forbidden",
						Detail: "User is not having access for the group. Access denied, please contact support.",
					},
				},
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
			result, err := client.GetPropertiesByGroup(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_GetDatasetFields(t *testing.T) {
	tests := map[string]struct {
		request          GetDatasetFieldsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []DataSets
		withError        error
	}{
		"200 OK": {
			request: GetDatasetFieldsRequest{
				TemplateName: TemplateNameEdgeLogs,
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
	{
	    "datasetGroupName":"group_name_1",
	    "datasetGroupDescription":"group_desc_1",
	    "datasetFields":[
	        {
	            "datasetFieldId":1000,
	            "datasetFieldName":"dataset_field_name_1",
	            "datasetFieldDescription":"dataset_field_desc_1"
	        },
	        {
	            "datasetFieldId":1002,
	            "datasetFieldName":"dataset_field_name_2",
	            "datasetFieldDescription":"dataset_field_desc_2"
	        }
	    ]
	},
	{
	    "datasetGroupName":"group_name_2",
	    "datasetFields":[
	        {
	            "datasetFieldId":1082,
	            "datasetFieldName":"dataset_field_name_3",
	            "datasetFieldDescription":"dataset_field_desc_3"
	        }
	    ]
	}
]
`,
			expectedPath: "/datastream-config-api/v1/log/datasets/template/EDGE_LOGS",
			expectedResponse: []DataSets{
				{
					DatasetGroupName:        "group_name_1",
					DatasetGroupDescription: "group_desc_1",
					DatasetFields: []DatasetFields{
						{
							DatasetFieldID:          1000,
							DatasetFieldName:        "dataset_field_name_1",
							DatasetFieldDescription: "dataset_field_desc_1",
						},
						{
							DatasetFieldID:          1002,
							DatasetFieldName:        "dataset_field_name_2",
							DatasetFieldDescription: "dataset_field_desc_2",
						},
					},
				},
				{
					DatasetGroupName: "group_name_2",
					DatasetFields: []DatasetFields{
						{
							DatasetFieldID:          1082,
							DatasetFieldName:        "dataset_field_name_3",
							DatasetFieldDescription: "dataset_field_desc_3",
						},
					},
				},
			},
		},
		"validation error - empty request": {
			request:   GetDatasetFieldsRequest{},
			withError: ErrStructValidation,
		},
		"validation error - invalid enum value": {
			request:   GetDatasetFieldsRequest{TemplateName: "invalidTemplateName"},
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
			result, err := client.GetDatasetFields(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
