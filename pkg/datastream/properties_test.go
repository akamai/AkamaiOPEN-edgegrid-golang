package datastream

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
			expectedPath: "/datastream-config-api/v3/log/cdn/groups/12345/properties",
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
			expectedPath: "/datastream-config-api/v3/log/cdn/groups/12345/properties",
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
			expectedPath: "/datastream-config-api/v3/log/cdn/groups/12345/properties",
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
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
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
		expectedErr      string
		withError        error
	}{
		"200 OK": {
			request: GetDatasetFieldsRequest{
				LogType:   LogTypeCDN,
				ProductID: "",
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
			expectedPath: "/datastream-config-api/v3/log/cdn/datasets-fields",
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
			request:        GetDatasetFieldsRequest{LogType: LogTypeCDN, ProductID: "INVALID_PROD_ID"},
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
			expectedPath: "/datastream-config-api/v3/log/cdn/datasets-fields?productId=INVALID_PROD_ID",
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
		"200 OK - answerx log type": {
			request: GetDatasetFieldsRequest{
				LogType: LogTypeAnswerX,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "datasetFields": [
        {
            "datasetFieldDescription": "datasetFieldDescription_1",
            "datasetFieldGroup": "datasetFieldGroup_1",
            "datasetFieldId": 2000,
            "datasetFieldJsonKey": "datasetFieldJsonKey_1",
            "datasetFieldName": "datasetFieldName_1"
        },
        {
            "datasetFieldDescription": "datasetFieldDescription_2",
            "datasetFieldGroup": "datasetFieldGroup_2",
            "datasetFieldId": 2001,
            "datasetFieldJsonKey": "datasetFieldJsonKey_2",
            "datasetFieldName": "datasetFieldName_2"
        }
    ]
}
`,
			expectedPath: "/datastream-config-api/v3/log/answerx/datasets-fields",
			expectedResponse: &DataSets{
				DataSetFields: []DataSetField{
					{
						DatasetFieldID:          2000,
						DatasetFieldName:        "datasetFieldName_1",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_1",
						DatasetFieldGroup:       "datasetFieldGroup_1",
						DatasetFieldDescription: "datasetFieldDescription_1",
					},
					{
						DatasetFieldID:          2001,
						DatasetFieldName:        "datasetFieldName_2",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_2",
						DatasetFieldGroup:       "datasetFieldGroup_2",
						DatasetFieldDescription: "datasetFieldDescription_2",
					},
				},
			},
		},
		"validation error - answerx with product id": {
			request: GetDatasetFieldsRequest{
				LogType:   LogTypeAnswerX,
				ProductID: "Ion_Standard",
			},
			expectedErr: "ProductID",
			withError:   ErrStructValidation,
		},
		"validation error - missing log type": {
			request:     GetDatasetFieldsRequest{},
			expectedErr: "LogType",
			withError:   ErrStructValidation,
		},
		"validation error - invalid log type": {
			request:     GetDatasetFieldsRequest{LogType: "INVALID"},
			expectedErr: "LogType",
			withError:   ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetDatasetFields(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				if test.expectedErr != "" {
					assert.Contains(t, err.Error(), test.expectedErr)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_GetAppSecConfigs(t *testing.T) {
	tests := map[string]struct {
		request          GetAppSecConfigsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []AppSecConfigDetails
		withError        error
	}{
		"200 OK": {
			request: GetAppSecConfigsRequest{
				GroupID:    12345,
				ContractID: "1-ABC",
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "fileType": "LOG",
        "id": 12345,
        "latestVersion": 3,
        "name": "WAF Security File",
        "productionVersion": 2,
        "targetProduct": "KSD"
    },
    {
        "fileType": "LOG",
        "id": 67890,
        "latestVersion": 1,
        "name": "Bot Manager Config",
        "productionVersion": 1,
        "targetProduct": "BOT_MANAGER"
    }
]
`,
			expectedPath: "/datastream-config-api/v3/log/appsec/groups/12345/contracts/1-ABC/configs",
			expectedResponse: []AppSecConfigDetails{
				{
					FileType:          "LOG",
					ID:                12345,
					LatestVersion:     3,
					Name:              "WAF Security File",
					ProductionVersion: 2,
					TargetProduct:     "KSD",
				},
				{
					FileType:          "LOG",
					ID:                67890,
					LatestVersion:     1,
					Name:              "Bot Manager Config",
					ProductionVersion: 1,
					TargetProduct:     "BOT_MANAGER",
				},
			},
		},
		"validation error - missing GroupID": {
			request:   GetAppSecConfigsRequest{ContractID: "1-ABC"},
			withError: ErrStructValidation,
		},
		"validation error - missing ContractID": {
			request:   GetAppSecConfigsRequest{GroupID: 12345},
			withError: ErrStructValidation,
		},
		"403 forbidden": {
			request: GetAppSecConfigsRequest{
				GroupID:    12345,
				ContractID: "1-ABC",
			},
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
			"detail": "User does not have access to the requested group."
		}
	]
}
`,
			expectedPath: "/datastream-config-api/v3/log/appsec/groups/12345/contracts/1-ABC/configs",
			withError: &Error{
				Type:       "forbidden",
				Title:      "Forbidden",
				Instance:   "28eb43a8-97ae-4c57-98aa-258081582b92",
				StatusCode: http.StatusForbidden,
				Errors: []RequestErrors{
					{
						Type:   "forbidden",
						Title:  "Forbidden",
						Detail: "User does not have access to the requested group.",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAppSecConfigs(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_ListAnswerXServiceIDs(t *testing.T) {
	tests := map[string]struct {
		request          ListAnswerXServiceIDsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListAnswerXServiceIDsResponse
		expectedErr      string
		withError        error
	}{
		"200 OK": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "1-ABC",
				PageSize:   1000,
				Page:       1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "metadata": {"lastPage":1,"pageSize":1000,"page":1,"totalElements":2},
    "contractId": "1-ABC",
    "serviceSubletterIds": [
        {
			"ssid": 101,
            "name": "ServiceA",
            "product": "AnswerX"
        },
        {
			"ssid": 102,
            "name": "ServiceB",
            "product": "AnswerX"
        }
    ]
}
`,
			expectedPath: "/datastream-config-api/v3/log/answerx/contracts/1-ABC/answerxSSIDs?page=1&pageSize=1000",
			expectedResponse: &ListAnswerXServiceIDsResponse{
				Metadata: &PaginationMetadata{
					LastPage:      1,
					PageSize:      1000,
					Page:          1,
					TotalElements: 2,
				},
				ContractID: "1-ABC",
				AnswerXServiceIDs: []AnswerXServiceDetail{
					{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
					{SSID: 102, Name: "ServiceB", Product: "AnswerX"},
				},
			},
		},
		"200 OK - no pagination params": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "1-ABC",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "1-ABC",
    "serviceSubletterIds": [
        {
			"ssid": 101,
            "name": "ServiceA",
            "product": "AnswerX"
        },
        {
			"ssid": 102,
            "name": "ServiceB",
            "product": "AnswerX"
        }
    ]
}
`,
			expectedPath: "/datastream-config-api/v3/log/answerx/contracts/1-ABC/answerxSSIDs",
			expectedResponse: &ListAnswerXServiceIDsResponse{
				ContractID: "1-ABC",
				AnswerXServiceIDs: []AnswerXServiceDetail{
					{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
					{SSID: 102, Name: "ServiceB", Product: "AnswerX"},
				},
			},
		},
		"validation error - missing ContractID": {
			request:     ListAnswerXServiceIDsRequest{},
			expectedErr: "ContractID: cannot be blank",
			withError:   ErrStructValidation,
		},
		"validation error - page without pageSize": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "1-ABC",
				Page:       1,
			},
			expectedErr: "page and pageSize must be provided together",
			withError:   ErrStructValidation,
		},
		"validation error - pageSize without page": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "1-ABC",
				PageSize:   1000,
			},
			expectedErr: "page and pageSize must be provided together",
			withError:   ErrStructValidation,
		},
		"validation error - negative page": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "1-ABC",
				Page:       -1,
				PageSize:   1000,
			},
			expectedErr: "Page: must be no less than 0",
			withError:   ErrStructValidation,
		},
		"validation error - negative pageSize": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "1-ABC",
				Page:       1,
				PageSize:   -1000,
			},
			expectedErr: "PageSize: must be no less than 0",
			withError:   ErrStructValidation,
		},
		"403 forbidden": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "1-ABC",
				PageSize:   1000,
				Page:       1,
			},
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
			"detail": "User does not have access to the requested contract."
		}
	]
}
`,
			expectedPath: "/datastream-config-api/v3/log/answerx/contracts/1-ABC/answerxSSIDs?page=1&pageSize=1000",
			withError: &Error{
				Type:       "forbidden",
				Title:      "Forbidden",
				Instance:   "28eb43a8-97ae-4c57-98aa-258081582b92",
				StatusCode: http.StatusForbidden,
				Errors: []RequestErrors{
					{
						Type:   "forbidden",
						Title:  "Forbidden",
						Detail: "User does not have access to the requested contract.",
					},
				},
			},
		},
		"400 bad request": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "INVALID",
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
	"type": "bad-request",
	"title": "Bad Request",
	"detail": "bad request",
	"instance": "82b67b97-d98d-4bee-ac1e-ef6eaf7cac82",
	"statusCode": 400,
	"errors": [
		{
			"type": "bad-request",
			"title": "Bad Request",
			"detail": "Invalid contract ID. Please provide a valid contract."
		}
	]
}
`,
			expectedPath: "/datastream-config-api/v3/log/answerx/contracts/INVALID/answerxSSIDs",
			withError: &Error{
				Type:       "bad-request",
				Title:      "Bad Request",
				Detail:     "bad request",
				Instance:   "82b67b97-d98d-4bee-ac1e-ef6eaf7cac82",
				StatusCode: http.StatusBadRequest,
				Errors: []RequestErrors{
					{
						Type:   "bad-request",
						Title:  "Bad Request",
						Detail: "Invalid contract ID. Please provide a valid contract.",
					},
				},
			},
		},
		"request execution error": {
			request: ListAnswerXServiceIDsRequest{
				ContractID: "1-ABC",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error": "internal server error"}`,
			expectedPath:   "/datastream-config-api/v3/log/answerx/contracts/1-ABC/answerxSSIDs",
			withError:      ErrListAnswerXServiceIDs,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListAnswerXServiceIDs(context.Background(), test.request)
			if test.withError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, test.withError)
				if test.expectedErr != "" {
					assert.Contains(t, err.Error(), test.expectedErr)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
