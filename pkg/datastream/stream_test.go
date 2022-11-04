package datastream

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDs_GetStream(t *testing.T) {
	tests := map[string]struct {
		request          GetStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DetailedStreamVersion
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: GetStreamRequest{
				StreamID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "streamId":1,
    "streamVersionId":2,
    "streamName":"ds2-sample-name",
    "datasets":[
        {
            "datasetGroupName":"group_name_1",
            "datasetGroupDescription":"group_desc_1",
            "datasetFields":[
                {
                    "datasetFieldId":1000,
                    "datasetFieldName":"dataset_field_name_1",
                    "datasetFieldDescription":"dataset_field_desc_1",
                    "order":0
                },
                {
                    "datasetFieldId":1002,
                    "datasetFieldName":"dataset_field_name_2",
                    "datasetFieldDescription":"dataset_field_desc_2",
                    "order":1
                }
            ]
        },
        {
            "datasetGroupName":"group_name_2",
            "datasetFields":[
                {
                    "datasetFieldId":1082,
                    "datasetFieldName":"dataset_field_name_3",
                    "datasetFieldDescription":"dataset_field_desc_3",
                    "order":32
                }
            ]
        }
    ],
    "connectors":[
        {
            "connectorType":"S3",
            "connectorId":13174,
            "bucket":"amzdemods2",
            "path":"/sample_path",
            "compressLogs":true,
            "connectorName":"aws_ds2_amz_demo",
            "region":"us-east-1"
        }
    ],
    "productName":"Adaptive Media Delivery",
    "productId":"Adaptive_Media_Delivery",
    "templateName":"EDGE_LOGS",
    "config":{
        "delimiter":"SPACE",
        "uploadFilePrefix":"ak",
        "uploadFileSuffix":"ds",
        "frequency":{
            "timeInSec":30
        },
        "useStaticPublicIP":false,
        "format":"STRUCTURED"
    },
    "groupId":171647,
    "groupName":"Akamai Data Delivery-P-132NZF456",
    "contractId":"P-132NZF456",
    "properties":[
        {
            "propertyId":678154,
            "propertyName":"amz.demo.com"
        }
    ],
    "streamType":"RAW_LOGS",
    "activationStatus":"ACTIVATED",
    "createdBy":"sample_username",
    "createdDate":"08-07-2021 06:00:27 GMT",
    "modifiedBy":"sample_username2",
    "modifiedDate":"08-07-2021 16:00:27 GMT",
    "emailIds":"sample_username@akamai.com"
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams/1",
			expectedResponse: &DetailedStreamVersion{
				ActivationStatus: ActivationStatusActivated,
				Config: Config{
					Delimiter: DelimiterTypePtr(DelimiterTypeSpace),
					Format:    FormatTypeStructured,
					Frequency: Frequency{
						TimeInSec: TimeInSec30,
					},
					UploadFilePrefix: "ak",
					UploadFileSuffix: "ds",
				},
				Connectors: []ConnectorDetails{
					{
						ConnectorID:   13174,
						CompressLogs:  true,
						ConnectorName: "aws_ds2_amz_demo",
						ConnectorType: ConnectorTypeS3,
						Path:          "/sample_path",
						Bucket:        "amzdemods2",
						Region:        "us-east-1",
					},
				},
				ContractID:  "P-132NZF456",
				CreatedBy:   "sample_username",
				CreatedDate: "08-07-2021 06:00:27 GMT",
				Datasets: []DataSets{
					{
						DatasetGroupName:        "group_name_1",
						DatasetGroupDescription: "group_desc_1",
						DatasetFields: []DatasetFields{
							{
								DatasetFieldID:          1000,
								DatasetFieldName:        "dataset_field_name_1",
								DatasetFieldDescription: "dataset_field_desc_1",
								Order:                   0,
							},
							{
								DatasetFieldID:          1002,
								DatasetFieldName:        "dataset_field_name_2",
								DatasetFieldDescription: "dataset_field_desc_2",
								Order:                   1,
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
								Order:                   32,
							},
						},
					},
				},
				EmailIDs:     "sample_username@akamai.com",
				GroupID:      171647,
				GroupName:    "Akamai Data Delivery-P-132NZF456",
				ModifiedBy:   "sample_username2",
				ModifiedDate: "08-07-2021 16:00:27 GMT",
				ProductID:    "Adaptive_Media_Delivery",
				ProductName:  "Adaptive Media Delivery",
				Properties: []Property{
					{
						PropertyID:   678154,
						PropertyName: "amz.demo.com",
					},
				},
				StreamID:        1,
				StreamName:      "ds2-sample-name",
				StreamType:      StreamTypeRawLogs,
				StreamVersionID: 2,
				TemplateName:    TemplateNameEdgeLogs,
			},
		},
		"validation error": {
			request: GetStreamRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"400 bad request": {
			request:        GetStreamRequest{StreamID: 12},
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/datastream-config-api/v1/log/streams/12",
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
			"detail": "Stream does not exist. Please provide valid stream."
		}
	]
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "bad-request",
					Title:      "Bad Request",
					Detail:     "bad request",
					Instance:   "82b67b97-d98d-4bee-ac1e-ef6eaf7cac82",
					StatusCode: http.StatusBadRequest,
					Errors: []RequestErrors{
						{
							Type:   "bad-request",
							Title:  "Bad Request",
							Detail: "Stream does not exist. Please provide valid stream.",
						},
					},
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
			result, err := client.GetStream(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_CreateStream(t *testing.T) {
	createStreamRequest := CreateStreamRequest{
		StreamConfiguration: StreamConfiguration{
			ActivateNow: true,
			Config: Config{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{TimeInSec: TimeInSec30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Connectors: []AbstractConnector{
				&S3Connector{
					Path:            "log/edgelogs/{ %Y/%m/%d }",
					ConnectorName:   "S3Destination",
					Bucket:          "datastream.akamai.com",
					Region:          "ap-south-1",
					AccessKey:       "AKIA6DK7TDQLVGZ3TYP1",
					SecretAccessKey: "1T2ll1H4dXWx5itGhpc7FlSbvvOvky1098nTtEMg",
				},
			},
			ContractID: "2-FGHIJ",
			DatasetFieldIDs: []int{
				1002, 1005, 1006, 1008, 1009, 1011, 1012,
				1013, 1014, 1015, 1016, 1017, 1101,
			},
			EmailIDs:     "useremail@akamai.com",
			GroupID:      tools.IntPtr(21484),
			PropertyIDs:  []int{123123, 123123},
			StreamName:   "TestStream",
			StreamType:   StreamTypeRawLogs,
			TemplateName: TemplateNameEdgeLogs,
		},
	}

	modifyRequest := func(r CreateStreamRequest, opt func(r *CreateStreamRequest)) CreateStreamRequest {
		opt(&r)
		return r
	}

	tests := map[string]struct {
		request          CreateStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedBody     string
		expectedResponse *StreamUpdate
		withError        error
	}{
		"202 Accepted": {
			request:        createStreamRequest,
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "streamVersionKey": {
        "streamId": 7050,
        "streamVersionId": 1
    }
}`,
			expectedPath: "/datastream-config-api/v1/log/streams",
			expectedResponse: &StreamUpdate{
				StreamVersionKey: StreamVersionKey{
					StreamID:        7050,
					StreamVersionID: 1,
				},
			},
			expectedBody: `
{
    "streamName": "TestStream",
    "activateNow": true,
    "streamType": "RAW_LOGS",
    "templateName": "EDGE_LOGS",
    "groupId": 21484,
    "contractId": "2-FGHIJ",
    "emailIds": "useremail@akamai.com",
    "propertyIds": [
        123123,
		123123
    ],
    "datasetFieldIds": [
        1002,
        1005,
        1006,
        1008,
        1009,
        1011,
        1012,
        1013,
        1014,
        1015,
        1016,
        1017,
        1101
    ],
    "config": {
        "uploadFilePrefix": "logs",
        "uploadFileSuffix": "ak",
        "delimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": {
            "timeInSec": 30
        }
    },
    "connectors": [
        {
            "path": "log/edgelogs/{ %Y/%m/%d }",
            "connectorName": "S3Destination",
            "bucket": "datastream.akamai.com",
            "region": "ap-south-1",
            "accessKey": "AKIA6DK7TDQLVGZ3TYP1",
            "secretAccessKey": "1T2ll1H4dXWx5itGhpc7FlSbvvOvky1098nTtEMg",
            "connectorType": "S3"
        }
    ]
}`,
		},
		"validation error - empty request": {
			request:   CreateStreamRequest{},
			withError: ErrStructValidation,
		},
		"validation error - empty connectors list": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.Connectors = []AbstractConnector{}
			}),
			withError: ErrStructValidation,
		},
		"validation error - delimiter with JSON format": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.Config = Config{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeJson,
					Frequency:        Frequency{TimeInSec: TimeInSec30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				}
			}),
			withError: ErrStructValidation,
		},
		"validation error - no delimiter with STRUCTURED format": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.Config = Config{
					Format:           FormatTypeStructured,
					Frequency:        Frequency{TimeInSec: TimeInSec30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				}
			}),
			withError: ErrStructValidation,
		},
		"validation error - missing connector configuration fields": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.Connectors = []AbstractConnector{
					&S3Connector{
						Path:          "log/edgelogs/{ %Y/%m/%d }",
						ConnectorName: "S3Destination",
						Bucket:        "datastream.akamai.com",
						Region:        "ap-south-1",
					},
				}
			}),
			withError: ErrStructValidation,
		},
		"403 forbidden": {
			request:        createStreamRequest,
			responseStatus: http.StatusForbidden,
			responseBody: `
{
	"type": "forbidden",
	"title": "Forbidden",
	"detail": "forbidden",
	"instance": "72a7654e-3f95-454f-a174-104bc946be52",
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
			expectedPath: "/datastream-config-api/v1/log/streams",
			withError: &Error{
				Type:       "forbidden",
				Title:      "Forbidden",
				Detail:     "forbidden",
				Instance:   "72a7654e-3f95-454f-a174-104bc946be52",
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
		"400 bad request": {
			request:        createStreamRequest,
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
	"type": "bad-request",
	"title": "Bad Request",
	"detail": "bad-request",
	"instance": "d0d2497e-ed93-4685-b44c-93a8eb8f3dea",
	"statusCode": 400,
	"errors": [
		{
			"type": "bad-request",
			"title": "Bad Request",
			"detail": "The credentials provided don’t give you write access to the bucket. Check your AWS credentials or bucket permissions in the S3 account and try again."
		}
	]
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams",
			withError: &Error{
				Type:       "bad-request",
				Title:      "Bad Request",
				Detail:     "bad-request",
				Instance:   "d0d2497e-ed93-4685-b44c-93a8eb8f3dea",
				StatusCode: http.StatusBadRequest,
				Errors: []RequestErrors{
					{
						Type:   "bad-request",
						Title:  "Bad Request",
						Detail: "The credentials provided don’t give you write access to the bucket. Check your AWS credentials or bucket permissions in the S3 account and try again.",
					},
				},
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

				//check request body only if we aren't testing errors
				if test.withError == nil {
					var reqBody interface{}
					err = json.NewDecoder(r.Body).Decode(&reqBody)
					require.NoError(t, err, "Error while decoding request body")

					var expectedBody interface{}
					err = json.Unmarshal([]byte(test.expectedBody), &expectedBody)
					require.NoError(t, err, "Error while parsing expected body to JSON")

					assert.Equal(t, expectedBody, reqBody)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateStream(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_UpdateStream(t *testing.T) {
	updateRequest := UpdateStreamRequest{
		StreamID: 7050,
		StreamConfiguration: StreamConfiguration{
			ActivateNow: true,
			Config: Config{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           "STRUCTURED",
				Frequency:        Frequency{TimeInSec: TimeInSec30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Connectors:      []AbstractConnector{},
			ContractID:      "P-132NZF456",
			DatasetFieldIDs: []int{1, 2, 3},
			EmailIDs:        "test@aka.mai",
			PropertyIDs:     []int{123123, 123123},
			StreamName:      "TestStream",
			StreamType:      "RAW_LOGS",
			TemplateName:    "EDGE_LOGS",
		},
	}

	modifyRequest := func(r UpdateStreamRequest, opt func(r *UpdateStreamRequest)) UpdateStreamRequest {
		opt(&r)
		return r
	}

	tests := map[string]struct {
		request          UpdateStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *StreamUpdate
		withError        error
	}{
		"202 Accepted": {
			request:        updateRequest,
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "streamVersionKey": {
        "streamId": 7050,
        "streamVersionId": 2
    }
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams/7050",
			expectedResponse: &StreamUpdate{
				StreamVersionKey: StreamVersionKey{
					StreamID:        7050,
					StreamVersionID: 2,
				},
			},
		},
		"validation error - empty request": {
			request:   UpdateStreamRequest{},
			withError: ErrStructValidation,
		},
		"validation error - delimiter with JSON format": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.Config = Config{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeJson,
					Frequency:        Frequency{TimeInSec: TimeInSec30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				}
			}),
			withError: ErrStructValidation,
		},
		"validation error - no delimiter with STRUCTURED format": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.Config = Config{
					Format:           FormatTypeStructured,
					Frequency:        Frequency{TimeInSec: TimeInSec60},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				}
			}),
			withError: ErrStructValidation,
		},
		"validation error - groupId modification": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.GroupID = tools.IntPtr(1337)
			}),
			withError: ErrStructValidation,
		},
		"validation error - missing contractId": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
			}),
			withError: ErrStructValidation,
		},
		"400 bad request": {
			request:        updateRequest,
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
	"type": "bad-request",
	"title": "Bad Request",
	"detail": "bad request",
	"instance": "a42cc1e6-fea4-4e3a-91ce-9da9819e089a",
	"statusCode": 400,
	"errors": [
		{
			"type": "bad-request",
			"title": "Bad Request",
			"detail": "Stream does not exist. Please provide valid stream."
		}
	]
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams/7050",
			withError: &Error{
				Type:       "bad-request",
				Title:      "Bad Request",
				Detail:     "bad request",
				Instance:   "a42cc1e6-fea4-4e3a-91ce-9da9819e089a",
				StatusCode: http.StatusBadRequest,
				Errors: []RequestErrors{
					{
						Type:   "bad-request",
						Title:  "Bad Request",
						Detail: "Stream does not exist. Please provide valid stream.",
					},
				},
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
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateStream(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_DeleteStream(t *testing.T) {
	tests := map[string]struct {
		request          DeleteStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DeleteStreamResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: DeleteStreamRequest{
				StreamID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "message": "Success"
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams/1",
			expectedResponse: &DeleteStreamResponse{
				Message: "Success",
			},
		},
		"validation error": {
			request: DeleteStreamRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"400 bad request": {
			request:        DeleteStreamRequest{StreamID: 12},
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/datastream-config-api/v1/log/streams/12",
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
			"detail": "Stream does not exist. Please provide valid stream."
		}
	]
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "bad-request",
					Title:      "Bad Request",
					Detail:     "bad request",
					Instance:   "82b67b97-d98d-4bee-ac1e-ef6eaf7cac82",
					StatusCode: 400,
					Errors: []RequestErrors{
						{
							Type:   "bad-request",
							Title:  "Bad Request",
							Detail: "Stream does not exist. Please provide valid stream.",
						},
					},
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteStream(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_Connectors(t *testing.T) {
	tests := map[string]struct {
		connector    AbstractConnector
		expectedJSON string
	}{
		"S3Connector": {
			connector: &S3Connector{
				Path:            "testPath",
				ConnectorName:   "testConnectorName",
				Bucket:          "testBucket",
				Region:          "testRegion",
				AccessKey:       "testAccessKey",
				SecretAccessKey: "testSecretKey",
			},
			expectedJSON: `
[{
	"path": "testPath",
	"connectorName": "testConnectorName",
	"bucket": "testBucket",
	"region": "testRegion",
	"accessKey": "testAccessKey",
	"secretAccessKey": "testSecretKey",
	"connectorType": "S3"
}]
`,
		},
		"AzureConnector": {
			connector: &AzureConnector{
				AccountName:   "testAccountName",
				AccessKey:     "testAccessKey",
				ConnectorName: "testConnectorName",
				ContainerName: "testContainerName",
				Path:          "testPath",
			},
			expectedJSON: `
[{
    "accountName": "testAccountName",
    "accessKey": "testAccessKey",
    "connectorName": "testConnectorName",
    "containerName": "testContainerName",
    "path": "testPath",
    "connectorType": "AZURE"
}]
`,
		},
		"DatadogConnector": {
			connector: &DatadogConnector{
				Service:       "testService",
				AuthToken:     "testAuthToken",
				ConnectorName: "testConnectorName",
				URL:           "testURL",
				Source:        "testSource",
				Tags:          "testTags",
				CompressLogs:  false,
			},
			expectedJSON: `
[{
    "service": "testService",
    "authToken": "testAuthToken",
    "connectorName": "testConnectorName",
    "url": "testURL",
    "source": "testSource",
    "tags": "testTags",
    "connectorType": "DATADOG",
    "compressLogs": false
}]
`,
		},
		"SplunkConnector": {
			connector: &SplunkConnector{
				ConnectorName:       "testConnectorName",
				URL:                 "testURL",
				EventCollectorToken: "testEventCollector",
				CompressLogs:        true,
				CustomHeaderName:    "custom-header",
				CustomHeaderValue:   "custom-header-value",
			},
			expectedJSON: `
[{
    "connectorName": "testConnectorName",
    "url": "testURL",
    "eventCollectorToken": "testEventCollector",
    "connectorType": "SPLUNK",
    "compressLogs": true,
	"customHeaderName": "custom-header",
	"customHeaderValue": "custom-header-value"
}]
`,
		},
		"GCSConnector": {
			connector: &GCSConnector{
				ConnectorName:      "testConnectorName",
				Bucket:             "testBucket",
				Path:               "testPath",
				ProjectID:          "testProjectID",
				ServiceAccountName: "testServiceAccountName",
				PrivateKey:         "testPrivateKey",
			},
			expectedJSON: `
[{
    "connectorType": "GCS",
    "connectorName": "testConnectorName",
    "bucket": "testBucket",
    "path": "testPath",
    "projectId": "testProjectID",
    "serviceAccountName": "testServiceAccountName",
	"privateKey": "testPrivateKey"
}]
`,
		},
		"CustomHTTPSConnector": {
			connector: &CustomHTTPSConnector{
				AuthenticationType: AuthenticationTypeBasic,
				ConnectorName:      "testConnectorName",
				URL:                "testURL",
				UserName:           "testUserName",
				Password:           "testPassword",
				CompressLogs:       true,
				CustomHeaderName:   "custom-header",
				CustomHeaderValue:  "custom-header-value",
				ContentType:        "application/json",
			},
			expectedJSON: `
[{
    "authenticationType": "BASIC",
    "connectorName": "testConnectorName",
    "url": "testURL",
    "userName": "testUserName",
    "password": "testPassword",
    "connectorType": "HTTPS",
    "compressLogs": true,
	"customHeaderName": "custom-header",
	"customHeaderValue": "custom-header-value",
	"contentType": "application/json"
}]
`,
		},
		"SumoLogicConnector": {
			connector: &SumoLogicConnector{
				ConnectorName:     "testConnectorName",
				Endpoint:          "testEndpoint",
				CollectorCode:     "testCollectorCode",
				CompressLogs:      true,
				CustomHeaderName:  "custom-header",
				CustomHeaderValue: "custom-header-value",
				ContentType:       "application/json",
			},
			expectedJSON: `
[{
    "connectorType": "SUMO_LOGIC",
    "connectorName": "testConnectorName",
    "endpoint": "testEndpoint",
    "collectorCode": "testCollectorCode",
    "compressLogs": true,
	"customHeaderName": "custom-header",
	"customHeaderValue": "custom-header-value",
	"contentType": "application/json"
}]
`,
		},
		"OracleCloudStorageConnector": {
			connector: &OracleCloudStorageConnector{
				AccessKey:       "testAccessKey",
				ConnectorName:   "testConnectorName",
				Path:            "testPath",
				Bucket:          "testBucket",
				Region:          "testRegion",
				SecretAccessKey: "testSecretAccessKey",
				Namespace:       "testNamespace",
			},
			expectedJSON: `
[{
    "accessKey": "testAccessKey",
    "connectorName": "testConnectorName",
    "path": "testPath",
    "bucket": "testBucket",
    "region": "testRegion",
    "secretAccessKey": "testSecretAccessKey",
    "connectorType": "Oracle_Cloud_Storage",
    "namespace": "testNamespace"
}]
`,
		},
		"LogglyConnector": {
			connector: &LogglyConnector{
				ConnectorName:     "testConnectorName",
				Endpoint:          "testEndpoint",
				AuthToken:         "testAuthToken",
				Tags:              "testTags",
				ContentType:       "testContentType",
				CustomHeaderName:  "testCustomHeaderName",
				CustomHeaderValue: "testCustomHeaderValue",
			},
			expectedJSON: `
[{
	"connectorType": "LOGGLY",
	"connectorName": "testConnectorName",
	"endpoint": "testEndpoint",
	"authToken": "testAuthToken",
	"tags": "testTags",
	"contentType": "testContentType",
	"customHeaderName": "testCustomHeaderName",
	"customHeaderValue": "testCustomHeaderValue"
}]
    `,
		},
		"NewRelicConnector": {
			connector: &NewRelicConnector{
				ConnectorName:     "testConnectorName",
				Endpoint:          "testEndpoint",
				AuthToken:         "testAuthToken",
				ContentType:       "testContentType",
				CustomHeaderName:  "testCustomHeaderName",
				CustomHeaderValue: "testCustomHeaderValue",
			},
			expectedJSON: `
[{
	"connectorType": "NEWRELIC",
	"connectorName": "testConnectorName",
	"endpoint": "testEndpoint",
	"authToken": "testAuthToken",
	"contentType": "testContentType",
	"customHeaderName": "testCustomHeaderName",
	"customHeaderValue": "testCustomHeaderValue"
}]
    `,
		},
		"ElasticsearchConnector": {
			connector: &ElasticsearchConnector{
				ConnectorName:     "testConnectorName",
				Endpoint:          "testEndpoint",
				IndexName:         "testIndexName",
				UserName:          "testUserName",
				Password:          "testPassword",
				ContentType:       "testContentType",
				CustomHeaderName:  "testCustomHeaderName",
				CustomHeaderValue: "testCustomHeaderValue",
				TLSHostname:       "testTLSHostname",
				CACert:            "testCACert",
				ClientCert:        "testClientCert",
				ClientKey:         "testClientKey",
			},
			expectedJSON: `
[{
	"connectorType": "ELASTICSEARCH",
	"connectorName": "testConnectorName",
	"endpoint": "testEndpoint",
	"indexName": "testIndexName",
	"userName": "testUserName",
	"password": "testPassword",
	"contentType": "testContentType",
	"customHeaderName": "testCustomHeaderName",
	"customHeaderValue": "testCustomHeaderValue",
	"tlsHostname": "testTLSHostname",
	"caCert": "testCACert",
	"clientCert": "testClientCert",
	"clientKey": "testClientKey"
}]
    `,
		},
	}

	request := CreateStreamRequest{
		StreamConfiguration: StreamConfiguration{
			ActivateNow: true,
			Config: Config{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{TimeInSec: TimeInSec30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Connectors:      nil,
			ContractID:      "P-132NZF456",
			DatasetFieldIDs: []int{1, 2, 3},
			EmailIDs:        "test@aka.mai",
			GroupID:         tools.IntPtr(123231),
			PropertyIDs:     []int{123123, 123123},
			StreamName:      "TestStream",
			StreamType:      StreamTypeRawLogs,
			TemplateName:    TemplateNameEdgeLogs,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			request.StreamConfiguration.Connectors = []AbstractConnector{test.connector}

			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var connectorMap map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&connectorMap)
				require.NoError(t, err)

				var expectedMap interface{}
				err = json.Unmarshal([]byte(test.expectedJSON), &expectedMap)
				require.NoError(t, err)

				res := reflect.DeepEqual(expectedMap, connectorMap["connectors"])
				assert.True(t, res)
			}))

			client := mockAPIClient(t, mockServer)
			_, _ = client.CreateStream(context.Background(), request)
		})
	}
}

type mockConnector struct {
	Called bool
}

func (c *mockConnector) SetConnectorType() {
	c.Called = true
}

func (c *mockConnector) Validate() error {
	return nil
}

func TestDs_setConnectorTypes(t *testing.T) {
	mockConnector := mockConnector{Called: false}

	request := CreateStreamRequest{
		StreamConfiguration: StreamConfiguration{
			ActivateNow: true,
			Config: Config{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{TimeInSec: TimeInSec30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Connectors: []AbstractConnector{
				&mockConnector,
			},
			ContractID:      "P-132NZF456",
			DatasetFieldIDs: []int{1, 2, 3},
			EmailIDs:        "test@aka.mai",
			GroupID:         tools.IntPtr(123231),
			PropertyIDs:     []int{123123, 123123},
			StreamName:      "TestStream",
			StreamType:      StreamTypeRawLogs,
			TemplateName:    TemplateNameEdgeLogs,
		},
	}

	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	}))
	client := mockAPIClient(t, mockServer)
	_, err := client.CreateStream(context.Background(), request)
	require.NoError(t, err)

	assert.True(t, mockConnector.Called)
}

func TestDs_ListStreams(t *testing.T) {
	tests := map[string]struct {
		request          ListStreamsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []StreamDetails
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request:        ListStreamsRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
[
  {
    "streamId": 1,
    "streamName": "Stream1",
    "streamVersionId": 2,
    "createdBy": "user1",
    "createdDate": "14-07-2020 07:07:40 GMT",
    "currentVersionId": 2,
    "archived": false,
    "activationStatus": "DEACTIVATED",
    "groupId": 1234,
    "groupName": "Default Group",
    "contractId": "1-ABCDE",
    "connectors": "S3-S1",
    "streamTypeName": "Logs - Raw",
    "properties": [
      {
        "propertyId": 13371337,
        "propertyName": "property_name_1"
      }
    ],
	"errors": [
      {
        "type": "ACTIVATION_ERROR",
        "title": "Activation/Deactivation Error",
        "detail": "Contact technical support."
      }
	]
  },
  {
    "streamId": 2,
    "streamName": "Stream2",
    "streamVersionId": 3,
    "createdBy": "user2",
    "createdDate": "24-07-2020 07:07:40 GMT",
    "currentVersionId": 3,
    "archived": true,
    "activationStatus": "ACTIVATED",
    "groupId": 4321,
    "groupName": "Default Group",
    "contractId": "2-ABCDE",
    "connectors": "S3-S2",
    "streamTypeName": "Logs - Raw",
    "properties": [
      {
        "propertyId": 23372337,
        "propertyName": "property_name_2"
      },
      {
        "propertyId": 33373337,
        "propertyName": "property_name_3"
      }
    ]
  }
]
`,
			expectedPath: "/datastream-config-api/v1/log/streams",
			expectedResponse: []StreamDetails{
				{
					ActivationStatus: ActivationStatusDeactivated,
					Archived:         false,
					Connectors:       "S3-S1",
					ContractID:       "1-ABCDE",
					CreatedBy:        "user1",
					CreatedDate:      "14-07-2020 07:07:40 GMT",
					CurrentVersionID: 2,
					Errors: []Errors{
						{
							Detail: "Contact technical support.",
							Title:  "Activation/Deactivation Error",
							Type:   "ACTIVATION_ERROR",
						},
					},
					GroupID:   1234,
					GroupName: "Default Group",
					Properties: []Property{
						{
							PropertyID:   13371337,
							PropertyName: "property_name_1",
						},
					},
					StreamID:        1,
					StreamName:      "Stream1",
					StreamTypeName:  "Logs - Raw",
					StreamVersionID: 2,
				},
				{
					ActivationStatus: ActivationStatusActivated,
					Archived:         true,
					Connectors:       "S3-S2",
					ContractID:       "2-ABCDE",
					CreatedBy:        "user2",
					CreatedDate:      "24-07-2020 07:07:40 GMT",
					CurrentVersionID: 3,
					Errors:           nil,
					GroupID:          4321,
					GroupName:        "Default Group",
					Properties: []Property{
						{
							PropertyID:   23372337,
							PropertyName: "property_name_2",
						},
						{
							PropertyID:   33373337,
							PropertyName: "property_name_3",
						},
					},
					StreamID:        2,
					StreamName:      "Stream2",
					StreamTypeName:  "Logs - Raw",
					StreamVersionID: 3,
				},
			},
		},
		"200 OK - with groupId": {
			request: ListStreamsRequest{
				GroupID: tools.IntPtr(1234),
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
  {
    "streamId": 2,
    "streamName": "Stream2",
    "streamVersionId": 3,
    "createdBy": "user2",
    "createdDate": "24-07-2020 07:07:40 GMT",
    "currentVersionId": 3,
    "archived": true,
    "activationStatus": "ACTIVATED",
    "groupId": 1234,
    "groupName": "Default Group",
    "contractId": "2-ABCDE",
    "connectors": "S3-S2",
    "streamTypeName": "Logs - Raw",
    "properties": [
      {
        "propertyId": 23372337,
        "propertyName": "property_name_2"
      },
      {
        "propertyId": 33373337,
        "propertyName": "property_name_3"
      }
    ]
  }
]
`,
			expectedPath: "/datastream-config-api/v1/log/streams?groupId=1234",
			expectedResponse: []StreamDetails{
				{
					ActivationStatus: ActivationStatusActivated,
					Archived:         true,
					Connectors:       "S3-S2",
					ContractID:       "2-ABCDE",
					CreatedBy:        "user2",
					CreatedDate:      "24-07-2020 07:07:40 GMT",
					CurrentVersionID: 3,
					Errors:           nil,
					GroupID:          1234,
					GroupName:        "Default Group",
					Properties: []Property{
						{
							PropertyID:   23372337,
							PropertyName: "property_name_2",
						},
						{
							PropertyID:   33373337,
							PropertyName: "property_name_3",
						},
					},
					StreamID:        2,
					StreamName:      "Stream2",
					StreamTypeName:  "Logs - Raw",
					StreamVersionID: 3,
				},
			},
		},
		"400 bad request": {
			request:        ListStreamsRequest{},
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/datastream-config-api/v1/log/streams",
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
			"detail": "Stream does not exist. Please provide valid stream."
		}
	]
}
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "bad-request",
					Title:      "Bad Request",
					Detail:     "bad request",
					Instance:   "82b67b97-d98d-4bee-ac1e-ef6eaf7cac82",
					StatusCode: http.StatusBadRequest,
					Errors: []RequestErrors{
						{
							Type:   "bad-request",
							Title:  "Bad Request",
							Detail: "Stream does not exist. Please provide valid stream.",
						},
					},
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
			result, err := client.ListStreams(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
