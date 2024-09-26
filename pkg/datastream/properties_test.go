package datastream

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDs_GetProperties(t *testing.T) {
	tests := map[string]struct {
		request          GetPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PropertiesDetails
		withError        error
	}{
		"200 OK": {
			request: GetPropertiesRequest{
				GroupId: 12345,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "groupId": 12345,
    "properties": [
    {
        "contractId": "1-7KLGU",
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
        "contractId": "1-7KLGU",
        "propertyId": 347459,
        "propertyName": "example.com",
        "productId": "Dynamic_Site_Accelerator",
        "productName": "Dynamic Site Accelerator",
        "hostnames": [
            "example.edgekey.net"
        ]
    }
]
}
`,
			expectedPath: "/datastream-config-api/v2/log/groups/12345/properties",
			expectedResponse: &PropertiesDetails{
				GroupID: 12345,
				Properties: []PropertyDetails{
					{
						ContractID:   "1-7KLGU",
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
						ContractID:   "1-7KLGU",
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
		},
		"validation error": {
			request:   GetPropertiesRequest{},
			withError: ErrStructValidation,
		},
		"400 bad request": {
			request:        GetPropertiesRequest{GroupId: 12345},
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
			expectedPath: "/datastream-config-api/v2/log/groups/12345/properties",
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
			request:        GetPropertiesRequest{GroupId: 12345},
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
			expectedPath: "/datastream-config-api/v2/log/groups/12345/properties",
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

func TestDs_GetDatasetFields(t *testing.T) {
	tests := map[string]struct {
		request          GetDatasetFieldsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DataSets
		withError        error
	}{
		"200 OK": {
			request: GetDatasetFieldsRequest{
				ProductID: nil,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "datasetFields": [
        {
            "datasetFieldDescription": "datasetFieldDescription_1",
            "datasetFieldGroup": "datasetFieldGroup_1",
            "datasetFieldId": 1000,
            "datasetFieldJsonKey": "datasetFieldJsonKey_1",
            "datasetFieldName": "datasetFieldName_1"
        },
        {
            "datasetFieldDescription": "datasetFieldDescription_2",
            "datasetFieldGroup": "datasetFieldGroup_2",
            "datasetFieldId": 1001,
            "datasetFieldJsonKey": "datasetFieldJsonKey_2",
            "datasetFieldName": "datasetFieldName_2"
        },
        {
            "datasetFieldDescription": "datasetFieldDescription_3",
            "datasetFieldGroup": "datasetFieldGroup_3",
            "datasetFieldId": 1002,
            "datasetFieldJsonKey": "datasetFieldJsonKey_3",
            "datasetFieldName": "datasetFieldName_3"
        }
    ]
}
`,
			expectedPath: "/datastream-config-api/v2/log/datasets-fields",
			expectedResponse: &DataSets{
				DataSetFields: []DataSetField{
					{
						DatasetFieldID:          1000,
						DatasetFieldName:        "datasetFieldName_1",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_1",
						DatasetFieldGroup:       "datasetFieldGroup_1",
						DatasetFieldDescription: "datasetFieldDescription_1",
					},
					{
						DatasetFieldID:          1001,
						DatasetFieldName:        "datasetFieldName_2",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_2",
						DatasetFieldGroup:       "datasetFieldGroup_2",
						DatasetFieldDescription: "datasetFieldDescription_2",
					},
					{
						DatasetFieldID:          1002,
						DatasetFieldName:        "datasetFieldName_3",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_3",
						DatasetFieldGroup:       "datasetFieldGroup_3",
						DatasetFieldDescription: "datasetFieldDescription_3",
					},
				},
			},
		},
		"validation error - invalid product id": {
			request:        GetDatasetFieldsRequest{ProductID: ptr.To("INVALID_PROD_ID")},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
    "errors": [
        {
            "detail": "Invalid product ID. Provide the correct product ID and try again.", 
            "problemId": "800a7291-c694-434a-99b7-8940d788239a", 
            "title": "Bad Request", 
            "type": "bad-request"
        }
    ], 
    "instance": "6e067164-4a61-429a-abaf-87452fd47036", 
    "problemId": "6e067164-4a61-429a-abaf-87452fd47036", 
    "status": 400, 
    "title": "Bad Request", 
    "type": "bad-request"
}
`,
			expectedPath: "/datastream-config-api/v2/log/datasets-fields?productId=INVALID_PROD_ID",
			withError: &Error{
				Type:       "bad-request",
				Title:      "Bad Request",
				Instance:   "6e067164-4a61-429a-abaf-87452fd47036",
				StatusCode: http.StatusBadRequest,
				Errors: []RequestErrors{
					{
						Type:   "bad-request",
						Title:  "Bad Request",
						Detail: "Invalid product ID. Provide the correct product ID and try again.",
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
