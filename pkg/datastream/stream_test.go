package datastream

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDs_GetStream(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          GetStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DetailedStreamVersion
		withError        func(*testing.T, error)
	}{
		"200 OK APPSEC Stream": {
			request: GetStreamRequest{
				LogType:  LogTypeAppSec,
				StreamID: 2,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "contractId": "P-1324", 
    "createdBy": "sample_username", 
    "createdDate": "2022-11-04T00:49:45Z", 
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": { "intervalInSeconds": 30 },
        "uploadFilePrefix": "ak",
        "uploadFileSuffix": "ds"
    },
    "destination": {
        "bucket": "sample_bucket",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample_display_name",
        "path": "/sample_path",
        "region": "us-east-1"
    },
    "groupId": 1234,
    "latestVersion": 2,
    "modifiedBy": "sample_username2",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": [
        "sample_username@akamai.com"
    ],
    "productId": "Adaptive_Media_Delivery",
	"datasetFields": [
        {
            "datasetFieldId":1000,
            "datasetFieldName":"dataset_field_name_1",
            "datasetFieldJsonKey":"dataset_field_json_key_1"
        },
        {
            "datasetFieldId":1002,
            "datasetFieldName":"dataset_field_name_2",
            "datasetFieldJsonKey":"dataset_field_json_key_2"
        },
        {
            "datasetFieldId":1082,
            "datasetFieldName":"dataset_field_name_3",
            "datasetFieldJsonKey":"dataset_field_json_key_3"
        }
    ],
    "appSecConfigs": [
        {
            "appSecId": 12345,
            "appSecName": "example_config"
        }
    ],
    "streamId": 2,
    "streamName": "ds2-appsec-stream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}`,
			expectedPath: "/datastream-config-api/v3/log/appsec/streams/2",
			expectedResponse: &DetailedStreamVersion{
				LogType:      LogTypeAppSec,
				StreamStatus: StreamStatusActivated,
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "ak",
					UploadFileSuffix: "ds",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "sample_display_name",
					DestinationType: DestinationTypeS3,
					Path:            "/sample_path",
					Bucket:          "sample_bucket",
					Region:          "us-east-1",
				},
				ContractID:  "P-1324",
				CreatedBy:   "sample_username",
				CreatedDate: "2022-11-04T00:49:45Z",
				DatasetFields: []DataSetField{
					{DatasetFieldID: 1000, DatasetFieldName: "dataset_field_name_1", DatasetFieldJsonKey: "dataset_field_json_key_1"},
					{DatasetFieldID: 1002, DatasetFieldName: "dataset_field_name_2", DatasetFieldJsonKey: "dataset_field_json_key_2"},
					{DatasetFieldID: 1082, DatasetFieldName: "dataset_field_name_3", DatasetFieldJsonKey: "dataset_field_json_key_3"},
				},
				NotificationEmails: []string{"sample_username@akamai.com"},
				GroupID:            1234,
				ModifiedBy:         "sample_username2",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				ProductID:          "Adaptive_Media_Delivery",
				AppSecConfigs: []AppSecConfig{
					{AppSecID: 12345, AppSecName: "example_config"},
				},
				StreamID:      2,
				StreamName:    "ds2-appsec-stream",
				StreamVersion: 2,
				LatestVersion: 2,
			},
		},
		"200 OK Without midgress, integrationType, and samplingPercentage": {
			request: GetStreamRequest{
				LogType:  LogTypeCDN,
				StreamID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "P-1324", 
    "createdBy": "sample_username", 
    "createdDate": "2022-11-04T00:49:45Z", 
    "datasetFields": [
        {
            "datasetFieldId":1000,
            "datasetFieldName":"dataset_field_name_1",
            "datasetFieldJsonKey":"dataset_field_json_key_1"
        },
        {
            "datasetFieldId":1002,
            "datasetFieldName":"dataset_field_name_2",
            "datasetFieldJsonKey":"dataset_field_json_key_2"
        },
        {
            "datasetFieldId":1082,
            "datasetFieldName":"dataset_field_name_3",
            "datasetFieldJsonKey":"dataset_field_json_key_3"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": { "intervalInSeconds": 30 },
        "uploadFilePrefix": "ak",
        "uploadFileSuffix": "ds"
    },
    "destination": {
        "bucket": "sample_bucket",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample_display_name",
        "path": "/sample_path",
        "region": "us-east-1"
    },
    "groupId": 1234,
    "latestVersion": 2,
    "modifiedBy": "sample_username2",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": [
        "sample_username@akamai.com"
    ],
    "productId": "Adaptive_Media_Delivery",
    "properties": [
        {
            "propertyId": 12345,
            "propertyName": "example.com"
        }
    ],
    "streamId": 1,
    "streamName": "ds2-sample-name",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams/1",
			expectedResponse: &DetailedStreamVersion{
				LogType:      LogTypeCDN,
				StreamStatus: StreamStatusActivated,
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "ak",
					UploadFileSuffix: "ds",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "sample_display_name",
					DestinationType: DestinationTypeS3,
					Path:            "/sample_path",
					Bucket:          "sample_bucket",
					Region:          "us-east-1",
				},
				ContractID:  "P-1324",
				CreatedBy:   "sample_username",
				CreatedDate: "2022-11-04T00:49:45Z",
				DatasetFields: []DataSetField{
					{DatasetFieldID: 1000, DatasetFieldName: "dataset_field_name_1", DatasetFieldJsonKey: "dataset_field_json_key_1"},
					{DatasetFieldID: 1002, DatasetFieldName: "dataset_field_name_2", DatasetFieldJsonKey: "dataset_field_json_key_2"},
					{DatasetFieldID: 1082, DatasetFieldName: "dataset_field_name_3", DatasetFieldJsonKey: "dataset_field_json_key_3"},
				},
				NotificationEmails: []string{"sample_username@akamai.com"},
				GroupID:            1234,
				ModifiedBy:         "sample_username2",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				ProductID:          "Adaptive_Media_Delivery",
				Properties: []Property{
					{PropertyID: 12345, PropertyName: "example.com"},
				},
				StreamID:      1,
				StreamName:    "ds2-sample-name",
				StreamVersion: 2,
				LatestVersion: 2,
			},
		},
		"200 OK With midgress field, integrationType DS_MANAGED and samplingPercentage 33": {
			request: GetStreamRequest{
				LogType:  LogTypeCDN,
				StreamID: 2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "X-DS-123",
    "createdBy": "admin",
    "createdDate": "2023-12-01T11:22:33Z",
    "collectMidgress": true,
    "integrationType": "DS_MANAGED",
    "samplingPercentage": 33,
    "datasetFields": [],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": { "intervalInSeconds": 22 },
        "uploadFilePrefix": "prefix",
        "uploadFileSuffix": "suffix"
    },
    "destination": {
        "bucket": "dsbucket",
        "compressLogs": false,
        "destinationType": "S3",
        "displayName": "adminBucket",
        "path": "/archive",
        "region": "us-east-2"
    },
    "groupId": 4567,
    "latestVersion": 5,
    "modifiedBy": "admin2",
    "modifiedDate": "2023-12-01T14:00:00Z",
    "notificationEmails": [
        "admin@akamai.com"
    ],
    "productId": "DSProd",
    "properties": [],
    "streamId": 2,
    "streamName": "stream-ds",
    "streamStatus": "ACTIVATED",
    "streamVersion": 5
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams/2",
			expectedResponse: &DetailedStreamVersion{
				LogType:         LogTypeCDN,
				CollectMidgress: true,
				StreamStatus:    StreamStatusActivated,
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: 22},
					UploadFilePrefix: "prefix",
					UploadFileSuffix: "suffix",
				},
				Destination: Destination{
					CompressLogs:    false,
					DisplayName:     "adminBucket",
					DestinationType: DestinationTypeS3,
					Path:            "/archive",
					Bucket:          "dsbucket",
					Region:          "us-east-2",
				},
				ContractID:         "X-DS-123",
				CreatedBy:          "admin",
				CreatedDate:        "2023-12-01T11:22:33Z",
				DatasetFields:      []DataSetField{},
				NotificationEmails: []string{"admin@akamai.com"},
				GroupID:            4567,
				ModifiedBy:         "admin2",
				ModifiedDate:       "2023-12-01T14:00:00Z",
				ProductID:          "DSProd",
				Properties:         []Property{},
				StreamID:           2,
				StreamName:         "stream-ds",
				StreamVersion:      5,
				LatestVersion:      5,
				IntegrationType:    "DS_MANAGED",
				SamplingPercentage: 33,
			},
		},
		"200 OK With integrationType PM_DEPENDENT and samplingPercentage 50": {
			request: GetStreamRequest{
				StreamID: 3,
				LogType:  LogTypeCDN,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "PM-DEP-001",
    "createdBy": "userpm",
    "createdDate": "2024-02-01T07:00:00Z",
    "integrationType": "PM_DEPENDENT",
    "samplingPercentage": 50,
    "datasetFields": [],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": { "intervalInSeconds": 90 },
        "uploadFilePrefix": "pm",
        "uploadFileSuffix": "pmx"
    },
    "destination": {
        "bucket": "pmbucket",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "pmBucket",
        "path": "/pmdata",
        "region": "eu-west-2"
    },
    "groupId": 2222,
    "latestVersion": 3,
    "modifiedBy": "pmadmin",
    "modifiedDate": "2024-02-01T09:00:00Z",
    "notificationEmails": ["userpm@akamai.com"],
    "productId": "PMProdX",
    "properties": [],
    "streamId": 3,
    "streamName": "pm_stream",
    "streamStatus": "INACTIVE",
    "streamVersion": 3
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams/3",
			expectedResponse: &DetailedStreamVersion{
				LogType:      LogTypeCDN,
				StreamStatus: StreamStatusInactive,
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: 90},
					UploadFilePrefix: "pm",
					UploadFileSuffix: "pmx",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "pmBucket",
					DestinationType: DestinationTypeS3,
					Path:            "/pmdata",
					Bucket:          "pmbucket",
					Region:          "eu-west-2",
				},
				ContractID:         "PM-DEP-001",
				CreatedBy:          "userpm",
				CreatedDate:        "2024-02-01T07:00:00Z",
				DatasetFields:      []DataSetField{},
				NotificationEmails: []string{"userpm@akamai.com"},
				GroupID:            2222,
				ModifiedBy:         "pmadmin",
				ModifiedDate:       "2024-02-01T09:00:00Z",
				ProductID:          "PMProdX",
				Properties:         []Property{},
				StreamID:           3,
				StreamName:         "pm_stream",
				StreamVersion:      3,
				LatestVersion:      3,
				IntegrationType:    "PM_DEPENDENT",
				SamplingPercentage: 50,
			},
		},
		"200 OK With integrationType HYBRID and samplingPercentage 88": {
			request: GetStreamRequest{
				StreamID: 4,
				LogType:  LogTypeCDN,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "HYB-777",
    "createdBy": "hybriduser",
    "createdDate": "2025-01-01T12:00:00Z",
    "integrationType": "HYBRID",
    "samplingPercentage": 88,
    "datasetFields": [],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": { "intervalInSeconds": 60 },
        "uploadFilePrefix": "hybrid",
        "uploadFileSuffix": "hyx"
    },
    "destination": {
        "bucket": "hybucket",
        "compressLogs": false,
        "destinationType": "S3",
        "displayName": "hyBucket",
        "path": "/hy_archive",
        "region": "ap-south-2"
    },
    "groupId": 9876,
    "latestVersion": 6,
    "modifiedBy": "hyadmin",
    "modifiedDate": "2025-01-01T16:00:00Z",
    "notificationEmails": ["hybriduser@akamai.com"],
    "productId": "HYProd",
    "properties": [],
    "streamId": 4,
    "streamName": "hy_stream",
    "streamStatus": "INACTIVE",
    "streamVersion": 6
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams/4",
			expectedResponse: &DetailedStreamVersion{
				LogType:      LogTypeCDN,
				StreamStatus: StreamStatusInactive,
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: 60},
					UploadFilePrefix: "hybrid",
					UploadFileSuffix: "hyx",
				},
				Destination: Destination{
					CompressLogs:    false,
					DisplayName:     "hyBucket",
					DestinationType: DestinationTypeS3,
					Path:            "/hy_archive",
					Bucket:          "hybucket",
					Region:          "ap-south-2",
				},
				ContractID:         "HYB-777",
				CreatedBy:          "hybriduser",
				CreatedDate:        "2025-01-01T12:00:00Z",
				DatasetFields:      []DataSetField{},
				NotificationEmails: []string{"hybriduser@akamai.com"},
				GroupID:            9876,
				ModifiedBy:         "hyadmin",
				ModifiedDate:       "2025-01-01T16:00:00Z",
				ProductID:          "HYProd",
				Properties:         []Property{},
				StreamID:           4,
				StreamName:         "hy_stream",
				StreamVersion:      6,
				LatestVersion:      6,
				IntegrationType:    "HYBRID",
				SamplingPercentage: 88,
			},
		},
		"200 OK ANSWERX Stream": {
			request: GetStreamRequest{
				LogType:  LogTypeAnswerX,
				StreamID: 5,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "contractId": "ANS-001",
    "createdBy": "answerx_user",
    "createdDate": "2025-03-01T08:00:00Z",
    "datasetFields": [
        {
            "datasetFieldId": 2000,
            "datasetFieldName": "answerx_field_1",
            "datasetFieldJsonKey": "answerx_key_1"
        },
		{
			"datasetFieldId": 2001,
			"datasetFieldName": "answerx_field_2",
			"datasetFieldJsonKey": "answerx_key_2"
		}
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": { "intervalInSeconds": 30 },
        "uploadFilePrefix": "ax",
        "uploadFileSuffix": "ds"
    },
    "destination": {
        "bucket": "answerx_bucket",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "answerx_display_name",
        "path": "/answerx_path",
        "region": "us-east-1"
    },
    "groupId": 5678,
    "latestVersion": 1,
    "modifiedBy": "answerx_user2",
    "modifiedDate": "2025-03-01T09:00:00Z",
    "notificationEmails": ["answerx@akamai.com"],
    "productId": "AnswerX_Product",
    "serviceSubletterIds": [
        {"ssid": 101, "name": "ServiceA", "product": "AnswerX"}
    ],
    "streamId": 5,
    "streamName": "answerx-stream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 1
}`,
			expectedPath: "/datastream-config-api/v3/log/answerx/streams/5",
			expectedResponse: &DetailedStreamVersion{
				LogType:      LogTypeAnswerX,
				StreamStatus: StreamStatusActivated,
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "ax",
					UploadFileSuffix: "ds",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "answerx_display_name",
					DestinationType: DestinationTypeS3,
					Path:            "/answerx_path",
					Bucket:          "answerx_bucket",
					Region:          "us-east-1",
				},
				ContractID:  "ANS-001",
				CreatedBy:   "answerx_user",
				CreatedDate: "2025-03-01T08:00:00Z",
				DatasetFields: []DataSetField{
					{DatasetFieldID: 2000, DatasetFieldName: "answerx_field_1", DatasetFieldJsonKey: "answerx_key_1"},
					{DatasetFieldID: 2001, DatasetFieldName: "answerx_field_2", DatasetFieldJsonKey: "answerx_key_2"},
				},
				NotificationEmails: []string{"answerx@akamai.com"},
				GroupID:            5678,
				ModifiedBy:         "answerx_user2",
				ModifiedDate:       "2025-03-01T09:00:00Z",
				ProductID:          "AnswerX_Product",
				AnswerXServiceIDs: []AnswerXServiceDetail{
					{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
				},
				StreamID:      5,
				StreamName:    "answerx-stream",
				StreamVersion: 1,
				LatestVersion: 1,
			},
		},
		"validation error": {
			request: GetStreamRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error - invalid log type": {
			request: GetStreamRequest{StreamID: 12, LogType: "INVALID"},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"400 bad request": {
			request:        GetStreamRequest{StreamID: 12, LogType: LogTypeCDN},
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/datastream-config-api/v3/log/cdn/streams/12",
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

		"200 OK without contractId and groupId": {
			request: GetStreamRequest{
				LogType:  LogTypeCDN,
				StreamID: 7050,
			},
			responseStatus:   http.StatusOK,
			responseBody:     createStreamResponseJSONWithoutContractAndGroupID,
			expectedPath:     "/datastream-config-api/v3/log/cdn/streams/7050",
			expectedResponse: createStreamResponseWithoutContractAndGroupID(),
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

func createStreamResponseWithoutContractAndGroupID() *DetailedStreamVersion {
	return &DetailedStreamVersion{
		CollectMidgress: true,
		CreatedBy:       "sample_username",
		CreatedDate:     "2022-11-04T00:49:45Z",
		DatasetFields: []DataSetField{
			{
				DatasetFieldName:    "field_name_1",
				DatasetFieldID:      2020,
				DatasetFieldJsonKey: "field_json_key_1",
			},
		},
		DeliveryConfiguration: DeliveryConfiguration{
			Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
			Format:           FormatTypeStructured,
			Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
			UploadFilePrefix: "logs",
			UploadFileSuffix: "ak",
		},
		Destination: Destination{
			CompressLogs:    true,
			DisplayName:     "sample-display-name",
			DestinationType: DestinationTypeS3,
			Path:            "sample-path/{%Y/%m/%d}",
			Bucket:          "datastream.com",
			Region:          "ap-south-1",
		},
		LatestVersion:      1,
		StreamID:           7050,
		StreamVersion:      1,
		StreamName:         "TestStream",
		StreamStatus:       StreamStatusActivated,
		ModifiedBy:         "sample_username2",
		ModifiedDate:       "2022-11-04T02:14:29Z",
		NotificationEmails: []string{"useremail1@akamai.com", "useremail2@akamai.com"},
		ProductID:          "Adaptive_Media_Delivery",
		Properties:         []Property{{PropertyID: 1234, PropertyName: "abcd"}, {PropertyID: 1234, PropertyName: "abcd"}},
		LogType:            LogTypeCDN,
	}
}

const createStreamResponseJSONWithoutContractAndGroupID = `
				{
					"createdBy": "sample_username",
					"createdDate": "2022-11-04T00:49:45Z",
					"collectMidgress": true,
					"datasetFields": [
						{
							"datasetFieldId":2020,
							"datasetFieldName":"field_name_1",
							"datasetFieldJsonKey":"field_json_key_1"
						}
					],
					"deliveryConfiguration": {
						"fieldDelimiter": "SPACE",
						"format": "STRUCTURED",
						"frequency": {
							"intervalInSeconds": 30
						},
						"uploadFilePrefix": "logs",
						"uploadFileSuffix": "ak"
					},
					"destination": {
						"bucket": "datastream.com",
						"compressLogs": true,
						"destinationType": "S3",
						"displayName": "sample-display-name",
						"path": "sample-path/{%Y/%m/%d}",
						"region": "ap-south-1"
					},
					"latestVersion": 1,
					"modifiedBy": "sample_username2",
					"modifiedDate": "2022-11-04T02:14:29Z",
					"notificationEmails": [
						"useremail1@akamai.com", "useremail2@akamai.com"
					],
					"productId": "Adaptive_Media_Delivery",
					"properties": [
						{
							"propertyId": 1234,
							"propertyName": "abcd"
						},
						{
							"propertyId": 1234,
							"propertyName": "abcd"
						}
					],
					"streamId": 7050,
					"streamName": "TestStream",
					"streamStatus": "ACTIVATED",
					"streamVersion": 1
				}
				`

const createStreamRequestBodyWithoutContractAndGroupID = `
				   "streamName":"TestStream",
				   "collectMidgress":true,
				   "notificationEmails":[
					  "useremail1@akamai.com",
					  "useremail2@akamai.com"
				   ],
				   "properties":[
					  {"propertyId":1234},{"propertyId":1234}
				   ],
				   "datasetFields":[{"datasetFieldId":2020}],
				   "deliveryConfiguration":{
					  "uploadFilePrefix":"logs",
					  "uploadFileSuffix":"ak",
					  "fieldDelimiter":"SPACE",
					  "format":"STRUCTURED",
					  "frequency": {"intervalInSeconds":30}
				   },
				   "destination":{
					  "path":"sample-path/{%Y/%m/%d}",
					  "displayName":"sample-display-name",
					  "bucket":"datastream.com",
					  "region":"ap-south-1",
					  "accessKey":"1234ABCD",
					  "secretAccessKey":"1234ABCD",
					  "destinationType":"S3"
				   }`

func updateStreamResponseWithoutContractAndGroupID() *DetailedStreamVersion {
	return &DetailedStreamVersion{
		CollectMidgress: true,
		CreatedBy:       "sample_username",
		CreatedDate:     "2022-11-04T00:49:45Z",
		DatasetFields:   []DataSetField{{DatasetFieldName: "field_name_1", DatasetFieldID: 2020, DatasetFieldJsonKey: "field_json_key_1"}},
		DeliveryConfiguration: DeliveryConfiguration{
			Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
			Format:           FormatTypeStructured,
			Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
			UploadFilePrefix: "logs",
			UploadFileSuffix: "ak",
		},
		Destination: Destination{
			CompressLogs:    true,
			DisplayName:     "sample-display-name",
			DestinationType: DestinationTypeS3,
			Path:            "sample-path/{%Y/%m/%d}",
			Bucket:          "datastream.com",
			Region:          "ap-south-1",
		},
		LatestVersion:      2,
		StreamID:           7050,
		StreamVersion:      2,
		StreamName:         "TestStream",
		StreamStatus:       StreamStatusActivated,
		ModifiedBy:         "modified_by_user",
		ModifiedDate:       "2022-11-04T02:14:29Z",
		NotificationEmails: []string{"useremail1@akamai.com", "useremail2@akamai.com"},
		ProductID:          "Adaptive_Media_Delivery",
		Properties:         []Property{{PropertyID: 1234, PropertyName: "sample1.com"}, {PropertyID: 1234, PropertyName: "sample2.com"}},
		LogType:            LogTypeCDN,
	}
}

const updateStreamResponseJSONWithoutContractAndGroupID = `
{
    "createdBy": "sample_username",
    "createdDate": "2022-11-04T00:49:45Z",
    "collectMidgress": true,
    "datasetFields": [
        {
            "datasetFieldId":2020,
            "datasetFieldName":"field_name_1",
            "datasetFieldJsonKey":"field_json_key_1"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": {
            "intervalInSeconds": 30
        },
        "uploadFilePrefix": "logs",
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample-display-name",
        "path": "sample-path/{%Y/%m/%d}",
        "region": "ap-south-1"
    },
    "latestVersion": 2,
    "modifiedBy": "modified_by_user",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": [
        "useremail1@akamai.com", "useremail2@akamai.com"
    ],
    "productId": "Adaptive_Media_Delivery",
    "properties": [
        {
            "propertyId": 1234,
            "propertyName": "sample1.com"
        },
        {
            "propertyId": 1234,
            "propertyName": "sample2.com"
        }
    ],
    "streamId": 7050,
    "streamName": "TestStream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}
`

const updateStreamRequestBodyWithoutContractAndGroupID = `
   "streamName":"TestStream",
   "notificationEmails":[
      "test@aka.mai",
      "useremail2@akamai.com"
   ],
   "properties":[
      {"propertyId":123123},
      {"propertyId":123123}
   ],
   "datasetFields":[
      {"datasetFieldId":1},
      {"datasetFieldId":2},
      {"datasetFieldId":3}
   ],
   "deliveryConfiguration":{
      "uploadFilePrefix":"logs",
      "uploadFileSuffix":"ak",
      "fieldDelimiter":"SPACE",
      "format":"STRUCTURED",
      "frequency": {"intervalInSeconds":30}
   },
   "destination":{
      "path":"sample-path/{%Y/%m/%d}",
      "displayName":"sample-display-name",
      "bucket":"datastream.com",
      "region":"ap-south-1",
      "accessKey":"ABC",
      "secretAccessKey":"XYZ",
      "destinationType":"S3"
   }`

func TestDs_CreateStream(t *testing.T) {
	t.Parallel()
	createStreamRequest := CreateStreamRequest{
		LogType:  LogTypeCDN,
		Activate: true,
		StreamConfiguration: StreamConfiguration{
			DeliveryConfiguration: DeliveryConfiguration{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Destination: AbstractConnector(
				&S3Connector{
					Path:            "sample-path/{%Y/%m/%d}",
					DisplayName:     "sample-display-name",
					Bucket:          "datastream.com",
					Region:          "ap-south-1",
					AccessKey:       "1234ABCD",
					SecretAccessKey: "1234ABCD",
				},
			),
			ContractID: "2-AB1234",
			DatasetFields: []DatasetFieldID{
				{DatasetFieldID: 2020},
			},
			NotificationEmails: []string{"useremail1@akamai.com", "useremail2@akamai.com"},
			GroupID:            1234,
			Properties: []PropertyID{
				{PropertyID: 1234},
				{PropertyID: 1234},
			},
			StreamName:         "TestStream",
			CollectMidgress:    true,
			SamplingPercentage: 0,
			// DO NOT set IntegrationType in request
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
		expectedResponse *DetailedStreamVersion
		expectedErr      string
		withError        error
	}{
		"201 Created APPSEC Stream": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.LogType = LogTypeAppSec
				r.StreamConfiguration.AppSecConfigs = []AppSecConfigID{
					{AppSecID: 12345},
				}
				r.StreamConfiguration.Properties = []PropertyID{}
				r.StreamConfiguration.DatasetFields = []DatasetFieldID{}
			}),
			responseStatus: http.StatusCreated,
			responseBody: `{
    "contractId": "2-AB1234", 
    "createdBy": "sample_username", 
    "createdDate": "2022-11-04T00:49:45Z", 
    "collectMidgress": true,
    "datasetFields": [
        {
            "datasetFieldId":2020,
            "datasetFieldName":"field_name_1",
            "datasetFieldJsonKey":"field_json_key_1"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE", 
        "format": "STRUCTURED", 
        "frequency": {
            "intervalInSeconds": 30
        }, 
        "uploadFilePrefix": "logs", 
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com", 
        "compressLogs": true, 
        "destinationType": "S3", 
        "displayName": "sample-display-name", 
        "path": "sample-path/{%Y/%m/%d}", 
        "region": "ap-south-1"
    },
    "groupId": 1234, 
    "latestVersion": 1, 
    "modifiedBy": "sample_username2", 
    "modifiedDate": "2022-11-04T02:14:29Z", 
    "notificationEmails": [
        "useremail1@akamai.com", "useremail2@akamai.com"
    ], 
    "productId": "Adaptive_Media_Delivery", 
	"appSecConfigs": [
        {
			"appSecId": 12345,
			"appSecName": "example_config"
		}
    ], 
    "streamId": 7050, 
    "streamName": "TestStream", 
    "streamStatus": "ACTIVATED", 
    "streamVersion": 1
}`,
			expectedPath: "/datastream-config-api/v3/log/appsec/streams?activate=true",
			expectedBody: `{
   "streamName":"TestStream",
   "groupId":1234,
   "contractId":"2-AB1234",
   "collectMidgress":true,
   "notificationEmails":[
      "useremail1@akamai.com",
      "useremail2@akamai.com"
   ],
   "appSecConfigs":[
      {"appSecId":12345}
   ],
   "deliveryConfiguration":{
      "uploadFilePrefix":"logs",
      "uploadFileSuffix":"ak",
      "fieldDelimiter":"SPACE",
      "format":"STRUCTURED",
      "frequency":{"intervalInSeconds":30}
   },
   "destination":{
      "path":"sample-path/{%Y/%m/%d}",
      "displayName":"sample-display-name",
      "bucket":"datastream.com",
      "region":"ap-south-1",
      "accessKey":"1234ABCD",
      "secretAccessKey":"1234ABCD",
      "destinationType":"S3"
   }
}
`,
			expectedResponse: &DetailedStreamVersion{
				LogType:         LogTypeAppSec,
				ContractID:      "2-AB1234",
				CreatedBy:       "sample_username",
				CreatedDate:     "2022-11-04T00:49:45Z",
				CollectMidgress: true,
				DatasetFields: []DataSetField{
					{
						DatasetFieldID:      2020,
						DatasetFieldName:    "field_name_1",
						DatasetFieldJsonKey: "field_json_key_1",
					},
				},
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: 30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				},
				Destination: Destination{
					Bucket:          "datastream.com",
					CompressLogs:    true,
					DestinationType: DestinationTypeS3,
					DisplayName:     "sample-display-name",
					Path:            "sample-path/{%Y/%m/%d}",
					Region:          "ap-south-1",
				},
				GroupID:       1234,
				LatestVersion: 1,
				ModifiedBy:    "sample_username2",
				ModifiedDate:  "2022-11-04T02:14:29Z",
				NotificationEmails: []string{
					"useremail1@akamai.com",
					"useremail2@akamai.com",
				},
				ProductID: "Adaptive_Media_Delivery",
				AppSecConfigs: []AppSecConfig{
					{
						AppSecID:   12345,
						AppSecName: "example_config",
					},
				},
				StreamID:      7050,
				StreamName:    "TestStream",
				StreamStatus:  StreamStatusActivated,
				StreamVersion: 1,
			},
		},
		"201 Created with IntegrationType PM_DEPENDENT and SamplingPercentage 75 (output only)": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.SamplingPercentage = 75
			}),
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "contractId": "2-AB1234",
    "createdBy": "pmuser",
    "createdDate": "2023-01-20T10:00:00Z",
    "integrationType": "PM_DEPENDENT",
    "samplingPercentage": 75,
    "collectMidgress": true,
    "datasetFields": [
        {"datasetFieldId":2020, "datasetFieldName":"pm_field", "datasetFieldJsonKey":"pm_key"}
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": {"intervalInSeconds": 15},
        "uploadFilePrefix": "pm_logs",
        "uploadFileSuffix": "pmx"
    },
    "destination": {
        "bucket": "pmbucket.com",
        "compressLogs": false,
        "destinationType": "S3",
        "displayName": "pm display-name",
        "path": "pm-path/{%Y/%m/%d}",
        "region": "eu-central-1"
    },
    "groupId": 1234,
    "latestVersion": 2,
    "modifiedBy": "pmuser2",
    "modifiedDate": "2023-01-21T12:00:00Z",
    "notificationEmails": ["pmuser@akamai.com"],
    "productId": "PMDelivery",
    "properties": [{"propertyId": 4321, "propertyName": "pmprop"}],
    "streamId": 9500,
    "streamName": "PMStream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams?activate=true",
			expectedBody: `
{
   "streamName":"TestStream",
   "groupId":1234,
   "contractId":"2-AB1234",
   "samplingPercentage":75,
   "collectMidgress":true,
   "notificationEmails":[
      "useremail1@akamai.com",
      "useremail2@akamai.com"
   ],
   "properties":[
      {"propertyId":1234},
      {"propertyId":1234}
   ],
   "datasetFields":[
      {"datasetFieldId":2020}
   ],
   "deliveryConfiguration":{
      "uploadFilePrefix":"logs",
      "uploadFileSuffix":"ak",
      "fieldDelimiter":"SPACE",
      "format":"STRUCTURED",
      "frequency":{"intervalInSeconds":30}
   },
   "destination":{
      "path":"sample-path/{%Y/%m/%d}",
      "displayName":"sample-display-name",
      "bucket":"datastream.com",
      "region":"ap-south-1",
      "accessKey":"1234ABCD",
      "secretAccessKey":"1234ABCD",
      "destinationType":"S3"
   }
}
`,
			expectedResponse: &DetailedStreamVersion{
				LogType:            LogTypeCDN,
				ContractID:         "2-AB1234",
				CreatedBy:          "pmuser",
				CreatedDate:        "2023-01-20T10:00:00Z",
				IntegrationType:    "PM_DEPENDENT",
				SamplingPercentage: 75,
				CollectMidgress:    true,
				DatasetFields: []DataSetField{
					{DatasetFieldID: 2020, DatasetFieldName: "pm_field", DatasetFieldJsonKey: "pm_key"},
				},
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: 15},
					UploadFilePrefix: "pm_logs",
					UploadFileSuffix: "pmx",
				},
				Destination: Destination{
					CompressLogs:    false,
					DisplayName:     "pm display-name",
					DestinationType: DestinationTypeS3,
					Path:            "pm-path/{%Y/%m/%d}",
					Bucket:          "pmbucket.com",
					Region:          "eu-central-1",
				},
				GroupID:            1234,
				LatestVersion:      2,
				StreamID:           9500,
				StreamVersion:      2,
				StreamName:         "PMStream",
				StreamStatus:       StreamStatusActivated,
				ModifiedBy:         "pmuser2",
				ModifiedDate:       "2023-01-21T12:00:00Z",
				NotificationEmails: []string{"pmuser@akamai.com"},
				ProductID:          "PMDelivery",
				Properties:         []Property{{PropertyID: 4321, PropertyName: "pmprop"}},
			},
		},
		"201 Created ActivateNow:true": {
			request:        createStreamRequest,
			responseStatus: http.StatusCreated,
			responseBody: `

{
    "contractId": "2-AB1234", 
    "createdBy": "sample_username", 
    "createdDate": "2022-11-04T00:49:45Z", 
    "collectMidgress": true,
    "datasetFields": [
        {
            "datasetFieldId":2020,
            "datasetFieldName":"field_name_1",
            "datasetFieldJsonKey":"field_json_key_1"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE", 
        "format": "STRUCTURED", 
        "frequency": {
            "intervalInSeconds": 30
        }, 
        "uploadFilePrefix": "logs", 
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com", 
        "compressLogs": true, 
        "destinationType": "S3", 
        "displayName": "sample-display-name", 
        "path": "sample-path/{%Y/%m/%d}", 
        "region": "ap-south-1"
    },
    "groupId": 1234, 
    "latestVersion": 1, 
    "modifiedBy": "sample_username2", 
    "modifiedDate": "2022-11-04T02:14:29Z", 
    "notificationEmails": [
        "useremail1@akamai.com", "useremail2@akamai.com"
    ], 
    "productId": "Adaptive_Media_Delivery", 
    "properties": [
        {
            "propertyId": 1234, 
            "propertyName": "abcd"
        },
        {
            "propertyId": 1234, 
            "propertyName": "abcd"
        }
    ], 
    "streamId": 7050, 
    "streamName": "TestStream", 
    "streamStatus": "ACTIVATED", 
    "streamVersion": 1
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams?activate=true",
			expectedResponse: &DetailedStreamVersion{
				LogType:         LogTypeCDN,
				CollectMidgress: true,
				ContractID:      "2-AB1234",
				CreatedBy:       "sample_username",
				CreatedDate:     "2022-11-04T00:49:45Z",
				DatasetFields: []DataSetField{
					{
						DatasetFieldName:    "field_name_1",
						DatasetFieldID:      2020,
						DatasetFieldJsonKey: "field_json_key_1",
					},
				},
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter: DelimiterTypePtr(DelimiterTypeSpace),
					Format:    FormatTypeStructured,
					Frequency: Frequency{
						IntervalInSeconds: IntervalInSeconds30,
					},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "sample-display-name",
					DestinationType: DestinationTypeS3,
					Path:            "sample-path/{%Y/%m/%d}",
					Bucket:          "datastream.com",
					Region:          "ap-south-1",
				},
				GroupID:            1234,
				LatestVersion:      1,
				StreamID:           7050,
				StreamVersion:      1,
				StreamName:         "TestStream",
				StreamStatus:       StreamStatusActivated,
				ModifiedBy:         "sample_username2",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				NotificationEmails: []string{"useremail1@akamai.com", "useremail2@akamai.com"},
				ProductID:          "Adaptive_Media_Delivery",
				Properties: []Property{
					{
						PropertyID:   1234,
						PropertyName: "abcd",
					},
					{
						PropertyID:   1234,
						PropertyName: "abcd",
					},
				},
			},
			expectedBody: `
{
   "streamName":"TestStream",
   "groupId":1234,
   "contractId":"2-AB1234",
   "collectMidgress":true,
   "notificationEmails":[
      "useremail1@akamai.com",
      "useremail2@akamai.com"
   ],
   "properties":[
      {"propertyId":1234},{"propertyId":1234}
   ],
   "datasetFields":[{"datasetFieldId":2020}],
   "deliveryConfiguration":{
      "uploadFilePrefix":"logs",
      "uploadFileSuffix":"ak",
      "fieldDelimiter":"SPACE",
      "format":"STRUCTURED",
      "frequency": {"intervalInSeconds":30}
   },
   "destination":{
      "path":"sample-path/{%Y/%m/%d}",
      "displayName":"sample-display-name",
      "bucket":"datastream.com",
      "region":"ap-south-1",
      "accessKey":"1234ABCD",
      "secretAccessKey":"1234ABCD",
      "destinationType":"S3"
   }
}
`,
		},
		"201 Created ActivateNow:true with omitted contractId and groupId": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			}),
			responseStatus:   http.StatusCreated,
			responseBody:     createStreamResponseJSONWithoutContractAndGroupID,
			expectedPath:     "/datastream-config-api/v3/log/cdn/streams?activate=true",
			expectedResponse: createStreamResponseWithoutContractAndGroupID(),
			expectedBody: `
				{
				   "streamName":"TestStream",
				   "collectMidgress":true,
				   "notificationEmails":[
					  "useremail1@akamai.com",
					  "useremail2@akamai.com"
				   ],
				   "properties":[
					  {"propertyId":1234},{"propertyId":1234}
				   ],
				   "datasetFields":[{"datasetFieldId":2020}],
				   "deliveryConfiguration":{
					  "uploadFilePrefix":"logs",
					  "uploadFileSuffix":"ak",
					  "fieldDelimiter":"SPACE",
					  "format":"STRUCTURED",
					  "frequency": {"intervalInSeconds":30}
				   },
				   "destination":{
					  "path":"sample-path/{%Y/%m/%d}",
					  "displayName":"sample-display-name",
					  "bucket":"datastream.com",
					  "region":"ap-south-1",
					  "accessKey":"1234ABCD",
					  "secretAccessKey":"1234ABCD",
				   "destinationType":"S3"
				   }
				}
				`,
		},
		"201 Created ActivateNow:true with only contractId": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = "1-ONLY-CONTRACT"
				r.StreamConfiguration.GroupID = 0
			}),
			responseStatus:   http.StatusCreated,
			responseBody:     createStreamResponseJSONWithoutContractAndGroupID,
			expectedPath:     "/datastream-config-api/v3/log/cdn/streams?activate=true",
			expectedResponse: createStreamResponseWithoutContractAndGroupID(),
			expectedBody: `{
				` + createStreamRequestBodyWithoutContractAndGroupID + `,
				"contractId":"1-ONLY-CONTRACT"
				}`,
		},
		"201 Created ActivateNow:true with only groupId": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 9999
			}),
			responseStatus:   http.StatusCreated,
			responseBody:     createStreamResponseJSONWithoutContractAndGroupID,
			expectedPath:     "/datastream-config-api/v3/log/cdn/streams?activate=true",
			expectedResponse: createStreamResponseWithoutContractAndGroupID(),
			expectedBody: `{
				` + createStreamRequestBodyWithoutContractAndGroupID + `,
				"groupId":9999
				}`,
		},
		"201 Created ActivateNow:true with whitespace contractId stripped": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = "   "
				r.StreamConfiguration.GroupID = 0
			}),
			responseStatus:   http.StatusCreated,
			responseBody:     createStreamResponseJSONWithoutContractAndGroupID,
			expectedPath:     "/datastream-config-api/v3/log/cdn/streams?activate=true",
			expectedResponse: createStreamResponseWithoutContractAndGroupID(),
			expectedBody: `{
				` + createStreamRequestBodyWithoutContractAndGroupID + `
				}`,
		},
		"validation error - negative groupId with omitted contractId": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = -1
			}),
			withError: ErrStructValidation,
		},
		"201 Created ActivateNow:false with omitted contractId and groupId": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.Activate = false
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			}),
			responseStatus:   http.StatusCreated,
			responseBody:     createStreamResponseJSONWithoutContractAndGroupID,
			expectedPath:     "/datastream-config-api/v3/log/cdn/streams?activate=false",
			expectedResponse: createStreamResponseWithoutContractAndGroupID(),
			expectedBody: `{
				` + createStreamRequestBodyWithoutContractAndGroupID + `
				}`,
		},
		"validation error - appsec omitted contractId and groupId": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.LogType = LogTypeAppSec
				r.StreamConfiguration.AppSecConfigs = []AppSecConfigID{{AppSecID: 12345}}
				r.StreamConfiguration.Properties = []PropertyID{}
				r.StreamConfiguration.DatasetFields = []DatasetFieldID{}
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			}),
			withError: ErrStructValidation,
		},
		"201 Created ANSWERX Stream": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.LogType = LogTypeAnswerX
				r.StreamConfiguration.Properties = nil
				r.StreamConfiguration.AnswerXServiceIDs = []AnswerXServiceID{
					{SSID: 101},
				}
				r.StreamConfiguration.CollectMidgress = false
			}),
			responseStatus: http.StatusCreated,
			responseBody: `{
    "contractId": "2-AB1234",
    "createdBy": "sample_username",
    "createdDate": "2022-11-04T00:49:45Z",
    "datasetFields": [
        {
            "datasetFieldId": 2020,
            "datasetFieldName": "field_name_1",
            "datasetFieldJsonKey": "field_json_key_1"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": { "intervalInSeconds": 30 },
        "uploadFilePrefix": "logs",
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample-display-name",
        "path": "sample-path/{%Y/%m/%d}",
        "region": "ap-south-1"
    },
    "groupId": 1234,
    "latestVersion": 1,
    "modifiedBy": "sample_username2",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": ["useremail1@akamai.com", "useremail2@akamai.com"],
    "productId": "Adaptive_Media_Delivery",
    "serviceSubletterIds": [
        {"ssid": 101, "name": "ServiceA", "product": "AnswerX"}
    ],
    "streamId": 8000,
    "streamName": "TestStream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 1
}`,
			expectedPath: "/datastream-config-api/v3/log/answerx/streams?activate=true",
			expectedBody: `{
   "streamName":"TestStream",
   "groupId":1234,
   "contractId":"2-AB1234",
   "notificationEmails":[
      "useremail1@akamai.com",
      "useremail2@akamai.com"
   ],
   "datasetFields":[{"datasetFieldId":2020}],
   "serviceSubletterIds":[{"ssid":101}],
   "deliveryConfiguration":{
      "uploadFilePrefix":"logs",
      "uploadFileSuffix":"ak",
      "fieldDelimiter":"SPACE",
      "format":"STRUCTURED",
      "frequency":{"intervalInSeconds":30}
   },
   "destination":{
      "path":"sample-path/{%Y/%m/%d}",
      "displayName":"sample-display-name",
      "bucket":"datastream.com",
      "region":"ap-south-1",
      "accessKey":"1234ABCD",
      "secretAccessKey":"1234ABCD",
      "destinationType":"S3"
   }
}
`,
			expectedResponse: &DetailedStreamVersion{
				LogType:     LogTypeAnswerX,
				ContractID:  "2-AB1234",
				CreatedBy:   "sample_username",
				CreatedDate: "2022-11-04T00:49:45Z",
				DatasetFields: []DataSetField{
					{DatasetFieldID: 2020, DatasetFieldName: "field_name_1", DatasetFieldJsonKey: "field_json_key_1"},
				},
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				},
				Destination: Destination{
					Bucket:          "datastream.com",
					CompressLogs:    true,
					DestinationType: DestinationTypeS3,
					DisplayName:     "sample-display-name",
					Path:            "sample-path/{%Y/%m/%d}",
					Region:          "ap-south-1",
				},
				GroupID:       1234,
				LatestVersion: 1,
				ModifiedBy:    "sample_username2",
				ModifiedDate:  "2022-11-04T02:14:29Z",
				NotificationEmails: []string{
					"useremail1@akamai.com",
					"useremail2@akamai.com",
				},
				ProductID: "Adaptive_Media_Delivery",
				AnswerXServiceIDs: []AnswerXServiceDetail{
					{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
				},
				StreamID:      8000,
				StreamName:    "TestStream",
				StreamStatus:  StreamStatusActivated,
				StreamVersion: 1,
			},
		},
		"validation error - empty destination": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.Destination = AbstractConnector(&S3Connector{})
			}),
			withError: ErrStructValidation,
		},
		"validation error - delimiter with JSON format": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.DeliveryConfiguration = DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeJson,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				}
			}),
			withError: ErrStructValidation,
		},
		"validation error - no delimiter with STRUCTURED format": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.DeliveryConfiguration = DeliveryConfiguration{
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				}
			}),
			withError: ErrStructValidation,
		},
		"validation error - missing destination configuration fields": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.Destination = AbstractConnector(
					&S3Connector{
						Path:        "log/edgelogs/{ %Y/%m/%d }",
						DisplayName: "S3Destination",
						Bucket:      "datastream.akamai.com",
						Region:      "ap-south-1",
					},
				)
			}),
			withError: ErrStructValidation,
		},
		"validation error - SamplingPercentage less than 1": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.SamplingPercentage = -1
			}),
			withError: ErrStructValidation,
		},
		"validation error - SamplingPercentage greater than 100": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.StreamConfiguration.SamplingPercentage = 101
			}),
			withError: ErrStructValidation,
		},
		"validation error - invalid log type": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.LogType = "INVALID"
			}),
			withError: ErrStructValidation,
		},
		"validation error - ANSWERX stream missing AnswerXServiceIDs": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.LogType = LogTypeAnswerX
				r.StreamConfiguration.AnswerXServiceIDs = nil
			}),
			expectedErr: "StreamConfiguration.AnswerXServiceIDs",
			withError:   ErrStructValidation,
		},
		"validation error - ANSWERX stream missing dataset fields": {
			request: modifyRequest(createStreamRequest, func(r *CreateStreamRequest) {
				r.LogType = LogTypeAnswerX
				r.StreamConfiguration.AnswerXServiceIDs = []AnswerXServiceID{{SSID: 101}}
				r.StreamConfiguration.DatasetFields = nil
			}),
			expectedErr: "StreamConfiguration.DatasetFields",
			withError:   ErrStructValidation,
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
			expectedPath: "/datastream-config-api/v3/log/cdn/streams?activate=true",
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
			expectedPath: "/datastream-config-api/v3/log/cdn/streams?activate=true",
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
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if test.withError == nil && test.expectedBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateStream(context.Background(), test.request)
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

func TestDs_UpdateStream(t *testing.T) {
	t.Parallel()
	updateRequest := UpdateStreamRequest{
		LogType:  LogTypeCDN,
		StreamID: 7050,
		Activate: true,
		StreamConfiguration: StreamConfiguration{
			DeliveryConfiguration: DeliveryConfiguration{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Destination: AbstractConnector(&S3Connector{
				DisplayName:     "sample-display-name",
				DestinationType: DestinationTypeS3,
				Path:            "sample-path/{%Y/%m/%d}",
				Bucket:          "datastream.com",
				Region:          "ap-south-1",
				AccessKey:       "ABC",
				SecretAccessKey: "XYZ",
			}),
			ContractID: "P-1324",
			DatasetFields: []DatasetFieldID{
				{DatasetFieldID: 1},
				{DatasetFieldID: 2},
				{DatasetFieldID: 3},
			},
			NotificationEmails: []string{"test@aka.mai", "useremail2@akamai.com"},
			Properties: []PropertyID{
				{PropertyID: 123123},
				{PropertyID: 123123},
			},
			StreamName: "TestStream",
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
		expectedBody     string
		expectedResponse *DetailedStreamVersion
		expectedErr      string
		withError        error
	}{
		"200 OK AppSec Streams": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.LogType = LogTypeAppSec
				r.StreamConfiguration.AppSecConfigs = []AppSecConfigID{
					{AppSecID: 12345},
				}
				r.StreamConfiguration.Properties = []PropertyID{}
				r.StreamConfiguration.DatasetFields = []DatasetFieldID{}
			}),
			responseStatus: http.StatusOK,
			expectedPath:   "/datastream-config-api/v3/log/appsec/streams/7050?activate=true",
			responseBody: `
{
    "contractId": "2-AB1234",
    "createdBy": "sample_username",
    "createdDate": "2022-11-04T00:49:45Z",
    "collectMidgress": true,
    "integrationType": "PM_DEPENDENT",
    "samplingPercentage": 55,
    "datasetFields": [
        {
            "datasetFieldId":2020,
            "datasetFieldName":"field_name_1",
            "datasetFieldJsonKey":"field_json_key_1"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": {
            "intervalInSeconds": 30
        },
        "uploadFilePrefix": "logs",
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample-display-name",
        "path": "sample-path/{%Y/%m/%d}",
        "region": "ap-south-1"
    },
    "groupId": 1234,
    "latestVersion": 2,
    "modifiedBy": "modified_by_user",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": [
        "useremail1@akamai.com", "useremail2@akamai.com"
    ],
    "productId": "Adaptive_Media_Delivery",
    "appSecConfigs": [
        {
            "appSecId": 12345,
			"appSecName": "example_config"
        }
    ],
    "streamId": 7050,
    "streamName": "TestStream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}
`,
			expectedResponse: &DetailedStreamVersion{
				LogType:            LogTypeAppSec,
				CollectMidgress:    true,
				ContractID:         "2-AB1234",
				CreatedBy:          "sample_username",
				CreatedDate:        "2022-11-04T00:49:45Z",
				IntegrationType:    "PM_DEPENDENT",
				SamplingPercentage: 55,
				DatasetFields: []DataSetField{
					{
						DatasetFieldName:    "field_name_1",
						DatasetFieldID:      2020,
						DatasetFieldJsonKey: "field_json_key_1",
					},
				},
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter: DelimiterTypePtr(DelimiterTypeSpace),
					Format:    FormatTypeStructured,
					Frequency: Frequency{
						IntervalInSeconds: IntervalInSeconds30,
					},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "sample-display-name",
					DestinationType: DestinationTypeS3,
					Path:            "sample-path/{%Y/%m/%d}",
					Bucket:          "datastream.com",
					Region:          "ap-south-1",
				},
				GroupID:            1234,
				LatestVersion:      2,
				StreamID:           7050,
				StreamVersion:      2,
				StreamName:         "TestStream",
				StreamStatus:       StreamStatusActivated,
				ModifiedBy:         "modified_by_user",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				NotificationEmails: []string{"useremail1@akamai.com", "useremail2@akamai.com"},
				ProductID:          "Adaptive_Media_Delivery",
				AppSecConfigs: []AppSecConfig{
					{
						AppSecID:   12345,
						AppSecName: "example_config",
					},
				},
			},
		},
		"200 OK activate:true with IntegrationType and SamplingPercentage": {
			request:        updateRequest,
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "2-AB1234",
    "createdBy": "sample_username",
    "createdDate": "2022-11-04T00:49:45Z",
    "collectMidgress": true,
    "integrationType": "PM_DEPENDENT",
    "samplingPercentage": 55,
    "datasetFields": [
        {
            "datasetFieldId":2020,
            "datasetFieldName":"field_name_1",
            "datasetFieldJsonKey":"field_json_key_1"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": {
            "intervalInSeconds": 30
        },
        "uploadFilePrefix": "logs",
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample-display-name",
        "path": "sample-path/{%Y/%m/%d}",
        "region": "ap-south-1"
    },
    "groupId": 1234,
    "latestVersion": 2,
    "modifiedBy": "modified_by_user",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": [
        "useremail1@akamai.com", "useremail2@akamai.com"
    ],
    "productId": "Adaptive_Media_Delivery",
    "properties": [
        {
            "propertyId": 1234,
            "propertyName": "sample1.com"
        },
        {
            "propertyId": 1234,
            "propertyName": "sample2.com"
        }
    ],
    "streamId": 7050,
    "streamName": "TestStream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams/7050?activate=true",
			expectedResponse: &DetailedStreamVersion{
				LogType:            LogTypeCDN,
				CollectMidgress:    true,
				ContractID:         "2-AB1234",
				CreatedBy:          "sample_username",
				CreatedDate:        "2022-11-04T00:49:45Z",
				IntegrationType:    "PM_DEPENDENT",
				SamplingPercentage: 55,
				DatasetFields: []DataSetField{
					{
						DatasetFieldName:    "field_name_1",
						DatasetFieldID:      2020,
						DatasetFieldJsonKey: "field_json_key_1",
					},
				},
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter: DelimiterTypePtr(DelimiterTypeSpace),
					Format:    FormatTypeStructured,
					Frequency: Frequency{
						IntervalInSeconds: IntervalInSeconds30,
					},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "sample-display-name",
					DestinationType: DestinationTypeS3,
					Path:            "sample-path/{%Y/%m/%d}",
					Bucket:          "datastream.com",
					Region:          "ap-south-1",
				},
				GroupID:            1234,
				LatestVersion:      2,
				StreamID:           7050,
				StreamVersion:      2,
				StreamName:         "TestStream",
				StreamStatus:       StreamStatusActivated,
				ModifiedBy:         "modified_by_user",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				NotificationEmails: []string{"useremail1@akamai.com", "useremail2@akamai.com"},
				ProductID:          "Adaptive_Media_Delivery",
				Properties: []Property{
					{
						PropertyID:   1234,
						PropertyName: "sample1.com",
					},
					{
						PropertyID:   1234,
						PropertyName: "sample2.com",
					},
				},
			},
		},
		"200 OK ANSWERX Stream": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.LogType = LogTypeAnswerX
				r.StreamConfiguration.Properties = nil
				r.StreamConfiguration.AnswerXServiceIDs = []AnswerXServiceID{
					{SSID: 101},
				}
			}),
			responseStatus: http.StatusOK,
			expectedPath:   "/datastream-config-api/v3/log/answerx/streams/7050?activate=true",
			expectedBody: `{
			   "streamName":"TestStream",
			   "contractId":"P-1324",
			   "notificationEmails":[
			      "test@aka.mai",
			      "useremail2@akamai.com"
			   ],
			   "datasetFields":[
			      {"datasetFieldId":1},
			      {"datasetFieldId":2},
			      {"datasetFieldId":3}
			   ],
			   "serviceSubletterIds":[
			      {"ssid":101}
			   ],
			   "deliveryConfiguration":{
			      "uploadFilePrefix":"logs",
			      "uploadFileSuffix":"ak",
			      "fieldDelimiter":"SPACE",
			      "format":"STRUCTURED",
			      "frequency":{"intervalInSeconds":30}
			   },
			   "destination":{
			      "path":"sample-path/{%Y/%m/%d}",
			      "displayName":"sample-display-name",
			      "bucket":"datastream.com",
			      "region":"ap-south-1",
			      "accessKey":"ABC",
			      "secretAccessKey":"XYZ",
			      "destinationType":"S3"
			   }
			}`,
			responseBody: `{
    "contractId": "P-1324",
    "createdBy": "sample_username",
    "createdDate": "2022-11-04T00:49:45Z",
    "datasetFields": [
        {"datasetFieldId": 1, "datasetFieldName": "field_1", "datasetFieldJsonKey": "key_1"},
        {"datasetFieldId": 2, "datasetFieldName": "field_2", "datasetFieldJsonKey": "key_2"},
        {"datasetFieldId": 3, "datasetFieldName": "field_3", "datasetFieldJsonKey": "key_3"}
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": { "intervalInSeconds": 30 },
        "uploadFilePrefix": "logs",
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample-display-name",
        "path": "sample-path/{%Y/%m/%d}",
        "region": "ap-south-1"
    },
    "groupId": 1234,
    "latestVersion": 2,
    "modifiedBy": "modified_by_user",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": ["test@aka.mai", "useremail2@akamai.com"],
    "productId": "AnswerX_Product",
    "serviceSubletterIds": [
        {"ssid": 101, "name": "ServiceA", "product": "AnswerX"}
    ],
    "streamId": 7050,
    "streamName": "TestStream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}`,
			expectedResponse: &DetailedStreamVersion{
				LogType:     LogTypeAnswerX,
				ContractID:  "P-1324",
				CreatedBy:   "sample_username",
				CreatedDate: "2022-11-04T00:49:45Z",
				DatasetFields: []DataSetField{
					{DatasetFieldID: 1, DatasetFieldName: "field_1", DatasetFieldJsonKey: "key_1"},
					{DatasetFieldID: 2, DatasetFieldName: "field_2", DatasetFieldJsonKey: "key_2"},
					{DatasetFieldID: 3, DatasetFieldName: "field_3", DatasetFieldJsonKey: "key_3"},
				},
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "sample-display-name",
					DestinationType: DestinationTypeS3,
					Path:            "sample-path/{%Y/%m/%d}",
					Bucket:          "datastream.com",
					Region:          "ap-south-1",
				},
				GroupID:            1234,
				LatestVersion:      2,
				StreamID:           7050,
				StreamVersion:      2,
				StreamName:         "TestStream",
				StreamStatus:       StreamStatusActivated,
				ModifiedBy:         "modified_by_user",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				NotificationEmails: []string{"test@aka.mai", "useremail2@akamai.com"},
				ProductID:          "AnswerX_Product",
				AnswerXServiceIDs: []AnswerXServiceDetail{
					{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
				},
			},
		},

		"validation error - delimiter with JSON format": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.DeliveryConfiguration = DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeJson,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				}
			}),
			withError: ErrStructValidation,
		},
		"validation error - no delimiter with STRUCTURED format": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.DeliveryConfiguration = DeliveryConfiguration{
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				}
			}),
			withError: ErrStructValidation,
		},
		"200 OK activate:true with omitted contractId and groupId": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			}),
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdBy": "sample_username",
    "createdDate": "2022-11-04T00:49:45Z",
    "collectMidgress": true,
    "datasetFields": [
        {
            "datasetFieldId":2020,
            "datasetFieldName":"field_name_1",
            "datasetFieldJsonKey":"field_json_key_1"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": {
            "intervalInSeconds": 30
        },
        "uploadFilePrefix": "logs",
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample-display-name",
        "path": "sample-path/{%Y/%m/%d}",
        "region": "ap-south-1"
    },
    "latestVersion": 2,
    "modifiedBy": "modified_by_user",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": [
        "useremail1@akamai.com", "useremail2@akamai.com"
    ],
    "productId": "Adaptive_Media_Delivery",
    "properties": [
        {
            "propertyId": 1234,
            "propertyName": "sample1.com"
        },
        {
            "propertyId": 1234,
            "propertyName": "sample2.com"
        }
    ],
    "streamId": 7050,
    "streamName": "TestStream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams/7050?activate=true",
			expectedResponse: &DetailedStreamVersion{
				CollectMidgress: true,
				CreatedBy:       "sample_username",
				CreatedDate:     "2022-11-04T00:49:45Z",
				DatasetFields:   []DataSetField{{DatasetFieldName: "field_name_1", DatasetFieldID: 2020, DatasetFieldJsonKey: "field_json_key_1"}},
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
					Format:           FormatTypeStructured,
					Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
					UploadFilePrefix: "logs",
					UploadFileSuffix: "ak",
				},
				Destination: Destination{
					CompressLogs:    true,
					DisplayName:     "sample-display-name",
					DestinationType: DestinationTypeS3,
					Path:            "sample-path/{%Y/%m/%d}",
					Bucket:          "datastream.com",
					Region:          "ap-south-1",
				},
				LatestVersion:      2,
				StreamID:           7050,
				StreamVersion:      2,
				StreamName:         "TestStream",
				StreamStatus:       StreamStatusActivated,
				ModifiedBy:         "modified_by_user",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				NotificationEmails: []string{"useremail1@akamai.com", "useremail2@akamai.com"},
				ProductID:          "Adaptive_Media_Delivery",
				Properties:         []Property{{PropertyID: 1234, PropertyName: "sample1.com"}, {PropertyID: 1234, PropertyName: "sample2.com"}},
				LogType:            LogTypeCDN,
			},
		},
		"200 OK activate:true with only contractId": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = "P-ONLY"
				r.StreamConfiguration.GroupID = 0
			}),
			responseStatus:   http.StatusOK,
			responseBody:     updateStreamResponseJSONWithoutContractAndGroupID,
			expectedPath:     "/datastream-config-api/v3/log/cdn/streams/7050?activate=true",
			expectedResponse: updateStreamResponseWithoutContractAndGroupID(),
			expectedBody: `{
				` + updateStreamRequestBodyWithoutContractAndGroupID + `,
				"contractId":"P-ONLY"
				}`,
		},
		"200 OK activate:true with only groupId": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 5678
			}),
			responseStatus:   http.StatusOK,
			responseBody:     updateStreamResponseJSONWithoutContractAndGroupID,
			expectedPath:     "/datastream-config-api/v3/log/cdn/streams/7050?activate=true",
			expectedResponse: updateStreamResponseWithoutContractAndGroupID(),
			expectedBody: `{
				` + updateStreamRequestBodyWithoutContractAndGroupID + `,
				"groupId":5678
				}`,
		},
		"200 OK activate:true with contractId and groupId": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = "2-AB1234"
				r.StreamConfiguration.GroupID = 5678
			}),
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "2-AB1234",
    "groupId": 5678,
    "createdBy": "sample_username",
    "createdDate": "2022-11-04T00:49:45Z",
    "collectMidgress": true,
    "datasetFields": [
        {
            "datasetFieldId":2020,
            "datasetFieldName":"field_name_1",
            "datasetFieldJsonKey":"field_json_key_1"
        }
    ],
    "deliveryConfiguration": {
        "fieldDelimiter": "SPACE",
        "format": "STRUCTURED",
        "frequency": {
            "intervalInSeconds": 30
        },
        "uploadFilePrefix": "logs",
        "uploadFileSuffix": "ak"
    },
    "destination": {
        "bucket": "datastream.com",
        "compressLogs": true,
        "destinationType": "S3",
        "displayName": "sample-display-name",
        "path": "sample-path/{%Y/%m/%d}",
        "region": "ap-south-1"
    },
    "latestVersion": 2,
    "modifiedBy": "modified_by_user",
    "modifiedDate": "2022-11-04T02:14:29Z",
    "notificationEmails": [
        "useremail1@akamai.com", "useremail2@akamai.com"
    ],
    "productId": "Adaptive_Media_Delivery",
    "properties": [
        {
            "propertyId": 1234,
            "propertyName": "sample1.com"
        },
        {
            "propertyId": 1234,
            "propertyName": "sample2.com"
        }
    ],
    "streamId": 7050,
    "streamName": "TestStream",
    "streamStatus": "ACTIVATED",
    "streamVersion": 2
}
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams/7050?activate=true",
			expectedResponse: func() *DetailedStreamVersion {
				resp := updateStreamResponseWithoutContractAndGroupID()
				resp.ContractID = "2-AB1234"
				resp.GroupID = 5678
				return resp
			}(),
			expectedBody: `{
				` + updateStreamRequestBodyWithoutContractAndGroupID + `,
				"contractId":"2-AB1234",
				"groupId":5678
				}`,
		},
		"validation error - negative groupId with omitted contractId": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = -5
			}),
			withError: ErrStructValidation,
		},
		"validation error - appsec omitted contractId": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.LogType = LogTypeAppSec
				r.StreamConfiguration.AppSecConfigs = []AppSecConfigID{{AppSecID: 12345}}
				r.StreamConfiguration.Properties = []PropertyID{}
				r.StreamConfiguration.DatasetFields = []DatasetFieldID{}
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			}),
			withError: ErrStructValidation,
		},
		"validation error - SamplingPercentage less than 1": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.SamplingPercentage = -1
			}),
			withError: ErrStructValidation,
		},
		"validation error - SamplingPercentage greater than 100": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.SamplingPercentage = 101
			}),
			withError: ErrStructValidation,
		},
		"validation error - invalid log type": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.LogType = "INVALID"
			}),
			expectedErr: "LogType",
			withError:   ErrStructValidation,
		},
		"validation error - ANSWERX stream missing AnswerXServiceIDs": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.LogType = LogTypeAnswerX
				r.StreamConfiguration.Properties = nil
			}),
			expectedErr: "StreamConfiguration.AnswerXServiceIDs",
			withError:   ErrStructValidation,
		},
		"validation error - ANSWERX stream missing dataset fields": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.LogType = LogTypeAnswerX
				r.StreamConfiguration.Properties = nil
				r.StreamConfiguration.AnswerXServiceIDs = []AnswerXServiceID{
					{SSID: 101},
				}
				r.StreamConfiguration.DatasetFields = nil
			}),
			expectedErr: "StreamConfiguration.DatasetFields",
			withError:   ErrStructValidation,
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
			expectedPath: "/datastream-config-api/v3/log/cdn/streams/7050?activate=true",
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
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if test.withError == nil && test.expectedBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateStream(context.Background(), test.request)
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

func TestDs_DeleteStream(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request        DeleteStreamRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No Content AppSec Stream": {
			request: DeleteStreamRequest{
				LogType:  LogTypeAppSec,
				StreamID: 1,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   ``,
			expectedPath:   "/datastream-config-api/v3/log/appsec/streams/1",
		},
		"204 No Content ANSWERX Stream": {
			request: DeleteStreamRequest{
				LogType:  LogTypeAnswerX,
				StreamID: 1,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   ``,
			expectedPath:   "/datastream-config-api/v3/log/answerx/streams/1",
		},
		"204 No Content CDN Stream": {
			request: DeleteStreamRequest{
				LogType:  LogTypeCDN,
				StreamID: 1,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   ``,
			expectedPath:   "/datastream-config-api/v3/log/cdn/streams/1",
		},
		"validation error": {
			request: DeleteStreamRequest{LogType: LogTypeCDN},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error - invalid log type": {
			request: DeleteStreamRequest{StreamID: 1, LogType: "INVALID"},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"400 bad request": {
			request: DeleteStreamRequest{
				StreamID: 12,
				LogType:  LogTypeCDN,
			},
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/datastream-config-api/v3/log/cdn/streams/12",
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
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			err := client.DeleteStream(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDs_Destinations(t *testing.T) {
	tests := map[string]struct {
		destination  AbstractConnector
		expectedJSON string
	}{
		"S3Connector": {
			destination: &S3Connector{
				Path:            "testPath",
				DisplayName:     "testDisplayName",
				Bucket:          "testBucket",
				Region:          "testRegion",
				AccessKey:       "testAccessKey",
				SecretAccessKey: "testSecretKey",
			},
			expectedJSON: `
{
	"path": "testPath",
	"displayName": "testDisplayName",
	"bucket": "testBucket",
	"region": "testRegion",
	"accessKey": "testAccessKey",
	"secretAccessKey": "testSecretKey",
	"destinationType": "S3"
}
`,
		},
		"AzureConnector": {
			destination: &AzureConnector{
				AccountName:   "testAccountName",
				AccessKey:     "testAccessKey",
				DisplayName:   "testDisplayName",
				ContainerName: "testContainerName",
				Path:          "testPath",
			},
			expectedJSON: `
{
    "accountName": "testAccountName",
    "accessKey": "testAccessKey",
    "displayName": "testDisplayName",
    "containerName": "testContainerName",
    "path": "testPath",
    "destinationType": "AZURE"
}
`,
		},
		"DatadogConnector": {
			destination: &DatadogConnector{
				Service:      "testService",
				AuthToken:    "testAuthToken",
				DisplayName:  "testDisplayName",
				Endpoint:     "testURL",
				Source:       "testSource",
				Tags:         "testTags",
				CompressLogs: false,
			},
			expectedJSON: `
{
    "service": "testService",
    "authToken": "testAuthToken",
    "displayName": "testDisplayName",
    "endpoint": "testURL",
    "source": "testSource",
    "tags": "testTags",
    "destinationType": "DATADOG",
    "compressLogs": false
}
`,
		},
		"SplunkConnector": {
			destination: &SplunkConnector{
				DisplayName:         "testDisplayName",
				Endpoint:            "testURL",
				EventCollectorToken: "testEventCollector",
				CompressLogs:        true,
				CustomHeaderName:    "custom-header",
				CustomHeaderValue:   "custom-header-value",
			},
			expectedJSON: `
{
    "displayName": "testDisplayName",
    "endpoint": "testURL",
    "eventCollectorToken": "testEventCollector",
    "destinationType": "SPLUNK",
    "compressLogs": true,
	"customHeaderName": "custom-header",
	"customHeaderValue": "custom-header-value"
}
`,
		},
		"GCSConnector": {
			destination: &GCSConnector{
				DisplayName:        "testDisplayName",
				Bucket:             "testBucket",
				Path:               "testPath",
				ProjectID:          "testProjectID",
				ServiceAccountName: "testServiceAccountName",
				PrivateKey:         "testPrivateKey",
			},
			expectedJSON: `
{
    "destinationType": "GCS",
    "displayName": "testDisplayName",
    "bucket": "testBucket",
    "path": "testPath",
    "projectId": "testProjectID",
    "serviceAccountName": "testServiceAccountName",
	"privateKey": "testPrivateKey"
}
`,
		},
		"CustomHTTPSConnector": {
			destination: &CustomHTTPSConnector{
				AuthenticationType: AuthenticationTypeBasic,
				DisplayName:        "testDisplayName",
				Endpoint:           "testURL",
				UserName:           "testUserName",
				Password:           "testPassword",
				CompressLogs:       true,
				CustomHeaderName:   "custom-header",
				CustomHeaderValue:  "custom-header-value",
				ContentType:        "application/json",
			},
			expectedJSON: `
{
    "authenticationType": "BASIC",
    "displayName": "testDisplayName",
    "endpoint": "testURL",
    "userName": "testUserName",
    "password": "testPassword",
    "destinationType": "HTTPS",
    "compressLogs": true,
	"customHeaderName": "custom-header",
	"customHeaderValue": "custom-header-value",
	"contentType": "application/json"
}
`,
		},
		"SumoLogicConnector": {
			destination: &SumoLogicConnector{
				DisplayName:       "testDisplayName",
				Endpoint:          "testEndpoint",
				CollectorCode:     "testCollectorCode",
				CompressLogs:      true,
				CustomHeaderName:  "custom-header",
				CustomHeaderValue: "custom-header-value",
				ContentType:       "application/json",
			},
			expectedJSON: `
{
    "destinationType": "SUMO_LOGIC",
    "displayName": "testDisplayName",
    "endpoint": "testEndpoint",
    "collectorCode": "testCollectorCode",
    "compressLogs": true,
	"customHeaderName": "custom-header",
	"customHeaderValue": "custom-header-value",
	"contentType": "application/json"
}
`,
		},
		"OracleCloudStorageConnector": {
			destination: &OracleCloudStorageConnector{
				AccessKey:       "testAccessKey",
				DisplayName:     "testDisplayName",
				Path:            "testPath",
				Bucket:          "testBucket",
				Region:          "testRegion",
				SecretAccessKey: "testSecretAccessKey",
				Namespace:       "testNamespace",
			},
			expectedJSON: `
{
    "accessKey": "testAccessKey",
    "displayName": "testDisplayName",
    "path": "testPath",
    "bucket": "testBucket",
    "region": "testRegion",
    "secretAccessKey": "testSecretAccessKey",
    "destinationType": "Oracle_Cloud_Storage",
    "namespace": "testNamespace"
}
`,
		},
		"LogglyConnector": {
			destination: &LogglyConnector{
				DisplayName:       "testDisplayName",
				Endpoint:          "testEndpoint",
				AuthToken:         "testAuthToken",
				Tags:              "testTags",
				ContentType:       "testContentType",
				CustomHeaderName:  "testCustomHeaderName",
				CustomHeaderValue: "testCustomHeaderValue",
			},
			expectedJSON: `
{
	"destinationType": "LOGGLY",
	"displayName": "testDisplayName",
	"endpoint": "testEndpoint",
	"authToken": "testAuthToken",
	"tags": "testTags",
	"contentType": "testContentType",
	"customHeaderName": "testCustomHeaderName",
	"customHeaderValue": "testCustomHeaderValue"
}
    `,
		},
		"NewRelicConnector": {
			destination: &NewRelicConnector{
				DisplayName:       "testDisplayName",
				Endpoint:          "testEndpoint",
				AuthToken:         "testAuthToken",
				ContentType:       "testContentType",
				CustomHeaderName:  "testCustomHeaderName",
				CustomHeaderValue: "testCustomHeaderValue",
			},
			expectedJSON: `
{
	"destinationType": "NEWRELIC",
	"displayName": "testDisplayName",
	"endpoint": "testEndpoint",
	"authToken": "testAuthToken",
	"contentType": "testContentType",
	"customHeaderName": "testCustomHeaderName",
	"customHeaderValue": "testCustomHeaderValue"
}
    `,
		},
		"ElasticsearchConnector": {
			destination: &ElasticsearchConnector{
				DisplayName:       "testDisplayName",
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
{
	"destinationType": "ELASTICSEARCH",
	"displayName": "testDisplayName",
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
}
`,
		},
		"S3CompatibleConnector": {
			destination: &S3CompatibleConnector{
				Path:            "testPath",
				DisplayName:     "testDisplayName",
				Bucket:          "testBucket",
				Region:          "testRegion",
				AccessKey:       "testAccessKey",
				SecretAccessKey: "testSecretKey",
				Endpoint:        "testEndpoint",
			},
			expectedJSON: `
{
	"path": "testPath",
	"displayName": "testDisplayName",
	"bucket": "testBucket",
	"region": "testRegion",
	"accessKey": "testAccessKey",
	"secretAccessKey": "testSecretKey",
	"destinationType": "S3_COMPATIBLE",
	"endpoint": "testEndpoint"
}
`,
		},
		"DynatraceConnector": {
			destination: &DynatraceConnector{
				DisplayName:       "testDisplayName",
				Endpoint:          "testEndpoint",
				AuthToken:         "testAuthToken",
				CustomHeaderName:  "testCustomHeaderName",
				CustomHeaderValue: "testCustomHeaderValue",
			},
			expectedJSON: `
{
	"destinationType": "DYNATRACE",
	"displayName": "testDisplayName",
	"endpoint": "testEndpoint",
	"authToken": "testAuthToken",
	"customHeaderName": "testCustomHeaderName",
	"customHeaderValue": "testCustomHeaderValue"
}
    `,
		},
		"TrafficPeakConnector": {
			destination: &TrafficPeakConnector{
				AuthenticationType: AuthenticationTypeBasic,
				DisplayName:        "testDisplayName",
				Endpoint:           "testURL",
				UserName:           "testUserName",
				Password:           "testPassword",
				CompressLogs:       true,
				CustomHeaderName:   "custom-header",
				CustomHeaderValue:  "custom-header-value",
				ContentType:        "application/json",
			},
			expectedJSON: `
{
    "authenticationType": "BASIC",
    "displayName": "testDisplayName",
    "endpoint": "testURL",
    "userName": "testUserName",
    "password": "testPassword",
    "destinationType": "TRAFFICPEAK",
    "compressLogs": true,
	"customHeaderName": "custom-header",
	"customHeaderValue": "custom-header-value",
	"contentType": "application/json"
}
`,
		},
	}

	request := CreateStreamRequest{
		Activate: true,
		StreamConfiguration: StreamConfiguration{
			DeliveryConfiguration: DeliveryConfiguration{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Destination: nil,
			ContractID:  "P-1324",
			DatasetFields: []DatasetFieldID{
				{
					DatasetFieldID: 1,
				},
				{
					DatasetFieldID: 2,
				},
				{
					DatasetFieldID: 3,
				},
			},

			NotificationEmails: []string{"test@aka.mai"},
			GroupID:            123231,
			Properties: []PropertyID{
				{
					PropertyID: 123123,
				},
				{
					PropertyID: 123123,
				},
			},
			StreamName: "TestStream",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request.StreamConfiguration.Destination = test.destination

			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				var destinationMap map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&destinationMap)
				require.NoError(t, err)

				var expectedMap interface{}
				err = json.Unmarshal([]byte(test.expectedJSON), &expectedMap)
				require.NoError(t, err)

				res := reflect.DeepEqual(expectedMap, destinationMap["destination"])
				assert.True(t, res)
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			_, _ = client.CreateStream(context.Background(), request)
		})
	}
}

type mockConnector struct {
	Called bool
}

func (c *mockConnector) SetDestinationType() {
	c.Called = true
}

func (c *mockConnector) Validate() error {
	return nil
}

func TestDs_setDestinationTypes(t *testing.T) {
	mockConnector := mockConnector{Called: false}

	request := CreateStreamRequest{
		LogType:  LogTypeCDN,
		Activate: true,
		StreamConfiguration: StreamConfiguration{
			DeliveryConfiguration: DeliveryConfiguration{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Destination: AbstractConnector(
				&mockConnector,
			),
			ContractID: "P-1324",

			DatasetFields: []DatasetFieldID{
				{
					DatasetFieldID: 1,
				},
				{
					DatasetFieldID: 2,
				},
				{
					DatasetFieldID: 3,
				},
			},

			NotificationEmails: []string{"test@aka.mai"},
			GroupID:            123231,
			Properties: []PropertyID{
				{
					PropertyID: 123123,
				},
				{
					PropertyID: 123123,
				},
			},
			StreamName: "TestStream",
		},
	}

	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	}))
	defer mockServer.Close()
	client := mockAPIClient(t, mockServer)
	_, err := client.CreateStream(context.Background(), request)
	require.NoError(t, err)

	assert.True(t, mockConnector.Called)
}

func TestDs_ListStreams(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          ListStreamsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []StreamDetails
		withError        func(*testing.T, error)
	}{
		"200 OK - AppSec streams": {
			request:        ListStreamsRequest{LogType: LogTypeAppSec, GroupID: ptr.To(123)},
			responseStatus: http.StatusOK,
			responseBody: `
[
   {
      "contractId":"1-ABC",
      "createdBy":"abc",
      "createdDate":"2022-04-21T17:02:58Z",
      "groupId":123,
      "latestVersion":15,
      "modifiedBy":"abc",
      "modifiedDate":"2022-12-26T17:00:03Z",
      "productId":"API_Acceleration",
	  "appSecConfigs" : [
	  	{
			"appSecId": 123,
			"appSecName": "test-appsec-config"
	  	}
	 ],
      "streamId":123,
      "streamName":"test-stream-1",
      "streamStatus":"ACTIVATED",
      "streamVersion":15
   },
   {
      "contractId":"1-123",
      "createdBy":"abc",
      "createdDate":"2023-01-03T12:44:15Z",
      "groupId":123,
      "latestVersion":1,
      "modifiedBy":"abc",
      "modifiedDate":"2023-01-03T12:44:15Z",
      "productId":"Download_Delivery",
      "appSecConfigs" : [
	  	{
			"appSecId": 123,
			"appSecName": "test-appsec-config"
	  	}
	  ],
      "streamId":123,
      "streamName":"test-stream-2",
      "streamStatus":"INACTIVE",
      "streamVersion":1
   }
]
`,
			expectedPath: "/datastream-config-api/v3/log/appsec/streams?groupId=123",
			expectedResponse: []StreamDetails{
				{
					LogType:       LogTypeAppSec,
					StreamStatus:  StreamStatusActivated,
					ProductID:     "API_Acceleration",
					ModifiedBy:    "abc",
					ModifiedDate:  "2022-12-26T17:00:03Z",
					ContractID:    "1-ABC",
					CreatedBy:     "abc",
					CreatedDate:   "2022-04-21T17:02:58Z",
					LatestVersion: 15,
					GroupID:       123,
					AppSecConfigs: []AppSecConfig{
						{
							AppSecID:   123,
							AppSecName: "test-appsec-config",
						},
					},
					StreamID:      123,
					StreamName:    "test-stream-1",
					StreamVersion: 15,
				},
				{
					LogType:       LogTypeAppSec,
					StreamStatus:  StreamStatusInactive,
					ProductID:     "Download_Delivery",
					ModifiedBy:    "abc",
					ModifiedDate:  "2023-01-03T12:44:15Z",
					ContractID:    "1-123",
					CreatedBy:     "abc",
					CreatedDate:   "2023-01-03T12:44:15Z",
					LatestVersion: 1,
					GroupID:       123,
					AppSecConfigs: []AppSecConfig{
						{
							AppSecID:   123,
							AppSecName: "test-appsec-config",
						},
					},
					StreamID:      123,
					StreamName:    "test-stream-2",
					StreamVersion: 1,
				},
			},
		},
		"200 OK - AnswerX streams": {
			request:        ListStreamsRequest{LogType: LogTypeAnswerX},
			responseStatus: http.StatusOK,
			responseBody: `
[
   {
      "contractId":"ANS-123",
      "createdBy":"answerx-user",
      "createdDate":"2025-03-01T08:00:00Z",
      "groupId":5678,
      "latestVersion":2,
      "modifiedBy":"answerx-admin",
      "modifiedDate":"2025-03-01T09:00:00Z",
      "productId":"AnswerX_Product",
      "serviceSubletterIds":[
         {
            "ssid":101,
            "name":"ServiceA",
            "product":"AnswerX"
         },
         {
            "ssid":202,
            "name":"ServiceB",
            "product":"AnswerX Plus"
         }
      ],
      "streamId":456,
      "streamName":"answerx-stream",
      "streamStatus":"ACTIVATED",
      "streamVersion":2
   }
]
`,
			expectedPath: "/datastream-config-api/v3/log/answerx/streams",
			expectedResponse: []StreamDetails{
				{
					LogType:       LogTypeAnswerX,
					StreamStatus:  StreamStatusActivated,
					ProductID:     "AnswerX_Product",
					ModifiedBy:    "answerx-admin",
					ModifiedDate:  "2025-03-01T09:00:00Z",
					ContractID:    "ANS-123",
					CreatedBy:     "answerx-user",
					CreatedDate:   "2025-03-01T08:00:00Z",
					LatestVersion: 2,
					GroupID:       5678,
					AnswerXServiceIDs: []AnswerXServiceDetail{
						{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
						{SSID: 202, Name: "ServiceB", Product: "AnswerX Plus"},
					},
					StreamID:      456,
					StreamName:    "answerx-stream",
					StreamVersion: 2,
				},
			},
		},
		"200 OK": {
			request:        ListStreamsRequest{LogType: LogTypeCDN},
			responseStatus: http.StatusOK,
			responseBody: `
[
   {
      "contractId":"1-ABC",
      "createdBy":"abc",
      "createdDate":"2022-04-21T17:02:58Z",
      "groupId":123,
      "latestVersion":15,
      "modifiedBy":"abc",
      "modifiedDate":"2022-12-26T17:00:03Z",
      "productId":"API_Acceleration",
      "properties":[
         {
            "propertyId":123,
            "propertyName":"example.com"
         },
         {
            "propertyId":123,
            "propertyName":"abc.media"
         }
      ],
      "streamId":123,
      "streamName":"test-stream-1",
      "streamStatus":"ACTIVATED",
      "streamVersion":15
   },
   {
      "contractId":"1-123",
      "createdBy":"abc",
      "createdDate":"2023-01-03T12:44:15Z",
      "groupId":123,
      "latestVersion":1,
      "modifiedBy":"abc",
      "modifiedDate":"2023-01-03T12:44:15Z",
      "productId":"Download_Delivery",
      "properties":[
         {
            "propertyId":123,
            "propertyName":"abc"
         }
      ],
      "streamId":123,
      "streamName":"test-stream-2",
      "streamStatus":"INACTIVE",
      "streamVersion":1
   }
]
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams",
			expectedResponse: []StreamDetails{
				{
					LogType:       LogTypeCDN,
					StreamStatus:  StreamStatusActivated,
					ProductID:     "API_Acceleration",
					ModifiedBy:    "abc",
					ModifiedDate:  "2022-12-26T17:00:03Z",
					ContractID:    "1-ABC",
					CreatedBy:     "abc",
					CreatedDate:   "2022-04-21T17:02:58Z",
					LatestVersion: 15,
					GroupID:       123,
					Properties: []Property{
						{
							PropertyID:   123,
							PropertyName: "example.com",
						},
						{
							PropertyID:   123,
							PropertyName: "abc.media",
						},
					},
					StreamID:      123,
					StreamName:    "test-stream-1",
					StreamVersion: 15,
				},
				{
					LogType:       LogTypeCDN,
					StreamStatus:  StreamStatusInactive,
					ProductID:     "Download_Delivery",
					ModifiedBy:    "abc",
					ModifiedDate:  "2023-01-03T12:44:15Z",
					ContractID:    "1-123",
					CreatedBy:     "abc",
					CreatedDate:   "2023-01-03T12:44:15Z",
					LatestVersion: 1,
					GroupID:       123,
					Properties: []Property{
						{
							PropertyID:   123,
							PropertyName: "abc",
						},
					},
					StreamID:      123,
					StreamName:    "test-stream-2",
					StreamVersion: 1,
				},
			},
		},
		"200 OK - with groupId": {
			request: ListStreamsRequest{
				GroupID: ptr.To(1234),
				LogType: LogTypeCDN,
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
  {
        "contractId": "1-123", 
        "createdBy": "abc", 
        "createdDate": "2022-07-25T08:36:32Z", 
        "groupId": 123, 
        "latestVersion": 2, 
        "modifiedBy": "abc", 
        "modifiedDate": "2022-12-26T20:00:02Z", 
        "productId": "Object_Delivery", 
        "properties": [
            {
                "propertyId": 123, 
                "propertyName": "abc.net"
            }
        ], 
        "streamId": 123, 
        "streamName": "test-stream", 
        "streamStatus": "ACTIVATED", 
        "streamVersion": 2
    }
]
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams?groupId=1234",
			expectedResponse: []StreamDetails{
				{
					LogType:       LogTypeCDN,
					StreamStatus:  StreamStatusActivated,
					ProductID:     "Object_Delivery",
					ModifiedBy:    "abc",
					ModifiedDate:  "2022-12-26T20:00:02Z",
					ContractID:    "1-123",
					CreatedBy:     "abc",
					CreatedDate:   "2022-07-25T08:36:32Z",
					LatestVersion: 2,
					GroupID:       123,
					Properties: []Property{
						{
							PropertyID:   123,
							PropertyName: "abc.net",
						},
					},
					StreamID:      123,
					StreamName:    "test-stream",
					StreamVersion: 2,
				},
			},
		},
		"200 OK - with IntegrationType and SamplingPercentage": {
			request:        ListStreamsRequest{LogType: LogTypeCDN},
			responseStatus: http.StatusOK,
			responseBody: `
[
   {
      "contractId":"PM-123",
      "createdBy":"pmuser",
      "createdDate":"2024-01-15T10:30:00Z",
      "groupId":456,
      "latestVersion":3,
      "modifiedBy":"pmuser",
      "modifiedDate":"2024-01-16T14:20:00Z",
      "productId":"Premium_Delivery",
      "properties":[
         {
            "propertyId":789,
            "propertyName":"premium.example.com"
         }
      ],
      "streamId":555,
      "streamName":"premium-stream",
      "streamStatus":"ACTIVATED",
      "streamVersion":3,
      "integrationType":"PM_DEPENDENT",
      "samplingPercentage":75
   },
   {
      "contractId":"DS-456",
      "createdBy":"dsuser",
      "createdDate":"2024-02-01T08:00:00Z",
      "groupId":789,
      "latestVersion":1,
      "modifiedBy":"dsuser",
      "modifiedDate":"2024-02-01T08:00:00Z",
      "productId":"Data_Stream",
      "properties":[
         {
            "propertyId":999,
            "propertyName":"datastream.example.com"
         }
      ],
      "streamId":777,
      "streamName":"ds-managed-stream",
      "streamStatus":"INACTIVE",
      "streamVersion":1,
      "integrationType":"DS_MANAGED",
      "samplingPercentage":100
   }
]
`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams",
			expectedResponse: []StreamDetails{
				{
					LogType:       LogTypeCDN,
					StreamStatus:  StreamStatusActivated,
					ProductID:     "Premium_Delivery",
					ModifiedBy:    "pmuser",
					ModifiedDate:  "2024-01-16T14:20:00Z",
					ContractID:    "PM-123",
					CreatedBy:     "pmuser",
					CreatedDate:   "2024-01-15T10:30:00Z",
					LatestVersion: 3,
					GroupID:       456,
					Properties: []Property{
						{
							PropertyID:   789,
							PropertyName: "premium.example.com",
						},
					},
					StreamID:           555,
					StreamName:         "premium-stream",
					StreamVersion:      3,
					IntegrationType:    "PM_DEPENDENT",
					SamplingPercentage: 75,
				},
				{
					LogType:       LogTypeCDN,
					StreamStatus:  StreamStatusInactive,
					ProductID:     "Data_Stream",
					ModifiedBy:    "dsuser",
					ModifiedDate:  "2024-02-01T08:00:00Z",
					ContractID:    "DS-456",
					CreatedBy:     "dsuser",
					CreatedDate:   "2024-02-01T08:00:00Z",
					LatestVersion: 1,
					GroupID:       789,
					Properties: []Property{
						{
							PropertyID:   999,
							PropertyName: "datastream.example.com",
						},
					},
					StreamID:           777,
					StreamName:         "ds-managed-stream",
					StreamVersion:      1,
					IntegrationType:    "DS_MANAGED",
					SamplingPercentage: 100,
				},
			},
		},
		"validation error - missing log type": {
			request: ListStreamsRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error - invalid log type": {
			request: ListStreamsRequest{LogType: "INVALID"},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"400 bad request": {
			request:        ListStreamsRequest{LogType: LogTypeCDN},
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/datastream-config-api/v3/log/cdn/streams",
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
		"200 OK without contractId and groupId": {
			request:        ListStreamsRequest{LogType: LogTypeCDN},
			responseStatus: http.StatusOK,
			responseBody: `[
   {
      "createdBy":"abc",
      "createdDate":"2022-04-21T17:02:58Z",
      "latestVersion":1,
      "modifiedBy":"abc",
      "modifiedDate":"2022-12-26T17:00:03Z",
      "productId":"API_Acceleration",
      "properties":[],
      "appSecConfigs":[],
      "streamId":99,
      "streamName":"no-contract-group-stream",
      "streamStatus":"ACTIVATED",
      "streamVersion":1
   }
]`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams",
			expectedResponse: []StreamDetails{
				{
					LogType:       LogTypeCDN,
					StreamStatus:  StreamStatusActivated,
					ProductID:     "API_Acceleration",
					ModifiedBy:    "abc",
					ModifiedDate:  "2022-12-26T17:00:03Z",
					CreatedBy:     "abc",
					CreatedDate:   "2022-04-21T17:02:58Z",
					LatestVersion: 1,
					Properties:    []Property{},
					AppSecConfigs: []AppSecConfig{},
					StreamID:      99,
					StreamName:    "no-contract-group-stream",
					StreamVersion: 1,
				},
			},
		},
		"200 OK groupId filter with response omitting contractId and groupId": {
			request:        ListStreamsRequest{LogType: LogTypeCDN, GroupID: ptr.To(1234)},
			responseStatus: http.StatusOK,
			responseBody: `[
   {
      "streamId":1,
      "streamName":"filtered",
      "streamStatus":"ACTIVATED",
      "streamVersion":1,
      "latestVersion":1,
      "createdBy":"u",
      "createdDate":"d",
      "modifiedBy":"u",
      "modifiedDate":"d",
      "productId":"p",
      "properties":[],
      "appSecConfigs":[]
   }
]`,
			expectedPath: "/datastream-config-api/v3/log/cdn/streams?groupId=1234",
			expectedResponse: []StreamDetails{
				{
					LogType:       LogTypeCDN,
					StreamStatus:  StreamStatusActivated,
					ProductID:     "p",
					ModifiedBy:    "u",
					ModifiedDate:  "d",
					CreatedBy:     "u",
					CreatedDate:   "d",
					LatestVersion: 1,
					Properties:    []Property{},
					AppSecConfigs: []AppSecConfig{},
					StreamID:      1,
					StreamName:    "filtered",
					StreamVersion: 1,
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

func validCDNCreateStreamRequest() CreateStreamRequest {
	return CreateStreamRequest{
		LogType:  LogTypeCDN,
		Activate: true,
		StreamConfiguration: StreamConfiguration{
			DeliveryConfiguration: DeliveryConfiguration{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Destination: AbstractConnector(&S3Connector{
				Path:            "sample-path/{%Y/%m/%d}",
				DisplayName:     "sample-display-name",
				Bucket:          "datastream.com",
				Region:          "ap-south-1",
				AccessKey:       "1234ABCD",
				SecretAccessKey: "1234ABCD",
			}),
			ContractID:         "2-AB1234",
			GroupID:            1234,
			DatasetFields:      []DatasetFieldID{{DatasetFieldID: 2020}},
			NotificationEmails: []string{"useremail1@akamai.com"},
			Properties:         []PropertyID{{PropertyID: 1234}},
			StreamName:         "TestStream",
			CollectMidgress:    true,
		},
	}
}

func validCDNUpdateStreamRequest() UpdateStreamRequest {
	return UpdateStreamRequest{
		LogType:  LogTypeCDN,
		StreamID: 7050,
		Activate: true,
		StreamConfiguration: StreamConfiguration{
			DeliveryConfiguration: DeliveryConfiguration{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           FormatTypeStructured,
				Frequency:        Frequency{IntervalInSeconds: IntervalInSeconds30},
				UploadFilePrefix: "logs",
				UploadFileSuffix: "ak",
			},
			Destination: AbstractConnector(&S3Connector{
				DisplayName:     "sample-display-name",
				DestinationType: DestinationTypeS3,
				Path:            "sample-path/{%Y/%m/%d}",
				Bucket:          "datastream.com",
				Region:          "ap-south-1",
				AccessKey:       "ABC",
				SecretAccessKey: "XYZ",
			}),
			ContractID:         "P-1324",
			DatasetFields:      []DatasetFieldID{{DatasetFieldID: 1}},
			NotificationEmails: []string{"test@aka.mai"},
			Properties:         []PropertyID{{PropertyID: 123123}},
			StreamName:         "TestStream",
		},
	}
}

// prepareAndValidateCreateStreamRequest sets the destination type, trims ContractID, then runs Validate.
func prepareAndValidateCreateStreamRequest(req CreateStreamRequest) error {
	setDestinationType(&req.StreamConfiguration)
	req.StreamConfiguration.ContractID = strings.TrimSpace(req.StreamConfiguration.ContractID)
	return req.Validate()
}

// prepareAndValidateUpdateStreamRequest sets the destination type, trims ContractID, then runs Validate.
func prepareAndValidateUpdateStreamRequest(req UpdateStreamRequest) error {
	setDestinationType(&req.StreamConfiguration)
	req.StreamConfiguration.ContractID = strings.TrimSpace(req.StreamConfiguration.ContractID)
	return req.Validate()
}

func validAppSecCreateStreamRequest() CreateStreamRequest {
	req := validCDNCreateStreamRequest()
	req.LogType = LogTypeAppSec
	req.StreamConfiguration.Properties = []PropertyID{}
	req.StreamConfiguration.DatasetFields = []DatasetFieldID{}
	req.StreamConfiguration.AppSecConfigs = []AppSecConfigID{{AppSecID: 12345}}
	return req
}

func TestCreateStreamRequest_Validate_OptionalContractAndGroup(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mutate  func(*CreateStreamRequest)
		wantErr string
	}{
		"both omitted": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			},
		},
		"only contract present": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = "1-ABC"
				r.StreamConfiguration.GroupID = 0
			},
		},
		"only group present": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 42
			},
		},
		"both present": {},
		"groupId one": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.GroupID = 1
			},
		},
		"groupId max int": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.GroupID = 2147483647
			},
		},
		"whitespace contract passes client validation": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = "   "
				r.StreamConfiguration.GroupID = 0
			},
		},
		"appsec both omitted": {
			mutate: func(r *CreateStreamRequest) {
				*r = validAppSecCreateStreamRequest()
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			},
			wantErr: "StreamConfiguration.ContractId",
		},
		"appsec only contract": {
			mutate: func(r *CreateStreamRequest) {
				*r = validAppSecCreateStreamRequest()
				r.StreamConfiguration.ContractID = "P-999"
				r.StreamConfiguration.GroupID = 0
			},
			wantErr: "StreamConfiguration.GroupID",
		},
		"appsec only group": {
			mutate: func(r *CreateStreamRequest) {
				*r = validAppSecCreateStreamRequest()
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 88
			},
			wantErr: "StreamConfiguration.ContractId",
		},
		"appsec whitespace contract": {
			mutate: func(r *CreateStreamRequest) {
				*r = validAppSecCreateStreamRequest()
				r.StreamConfiguration.ContractID = "   "
				r.StreamConfiguration.GroupID = 1
			},
			wantErr: "StreamConfiguration.ContractId",
		},
		"groupId negative": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.GroupID = -1
			},
			wantErr: "StreamConfiguration.GroupID",
		},
		"omitted optional still requires stream name": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
				r.StreamConfiguration.StreamName = ""
			},
			wantErr: "StreamConfiguration.StreamName",
		},
		"omitted optional still requires destination": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
				r.StreamConfiguration.Destination = AbstractConnector(&S3Connector{})
			},
			wantErr: "StreamConfiguration.Destination",
		},
		"omitted optional still requires valid log type": {
			mutate: func(r *CreateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
				r.LogType = "INVALID"
			},
			wantErr: "LogType",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := validCDNCreateStreamRequest()
			if tc.mutate != nil {
				tc.mutate(&req)
			}
			err := prepareAndValidateCreateStreamRequest(req)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestUpdateStreamRequest_Validate_OptionalContractAndGroup(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mutate  func(*UpdateStreamRequest)
		wantErr string
	}{
		"both omitted": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			},
		},
		"only contract present": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = "1-ABC"
				r.StreamConfiguration.GroupID = 0
			},
		},
		"only group present": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 42
			},
		},
		"non-zero groupId allowed on update": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.GroupID = 5
			},
		},
		"groupId negative": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.GroupID = -1
			},
			wantErr: "StreamConfiguration.GroupID",
		},
		"groupId max int": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.GroupID = 2147483647
			},
		},
		"omitted optional still requires stream name": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
				r.StreamConfiguration.StreamName = ""
			},
			wantErr: "StreamConfiguration.StreamName",
		},
		"omitted optional still requires destination": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
				r.StreamConfiguration.Destination = AbstractConnector(&S3Connector{})
			},
			wantErr: "StreamConfiguration.Destination",
		},
		"omitted optional still requires valid log type": {
			mutate: func(r *UpdateStreamRequest) {
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
				r.LogType = "INVALID"
			},
			wantErr: "LogType",
		},
		"appsec both omitted": {
			mutate: func(r *UpdateStreamRequest) {
				*r = validAppSecUpdateStreamRequest()
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 0
			},
			wantErr: "StreamConfiguration.ContractId",
		},
		"appsec both present": {
			mutate: func(r *UpdateStreamRequest) {
				*r = validAppSecUpdateStreamRequest()
			},
		},
		"appsec whitespace contract": {
			mutate: func(r *UpdateStreamRequest) {
				*r = validAppSecUpdateStreamRequest()
				r.StreamConfiguration.ContractID = "  "
				r.StreamConfiguration.GroupID = 0
			},
			wantErr: "StreamConfiguration.ContractId",
		},
		"appsec non-zero groupId rejected": {
			mutate: func(r *UpdateStreamRequest) {
				*r = validAppSecUpdateStreamRequest()
				r.StreamConfiguration.GroupID = 1
			},
			wantErr: "StreamConfiguration.GroupID",
		},
		"appsec only contract": {
			mutate: func(r *UpdateStreamRequest) {
				*r = validAppSecUpdateStreamRequest()
				r.StreamConfiguration.ContractID = "P-999"
				r.StreamConfiguration.GroupID = 0
			},
		},
		"appsec only group rejected": {
			mutate: func(r *UpdateStreamRequest) {
				*r = validAppSecUpdateStreamRequest()
				r.StreamConfiguration.ContractID = ""
				r.StreamConfiguration.GroupID = 88
			},
			wantErr: "StreamConfiguration.ContractId",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := validCDNUpdateStreamRequest()
			if tc.mutate != nil {
				tc.mutate(&req)
			}
			err := prepareAndValidateUpdateStreamRequest(req)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestStreamConfiguration_JSONOptionalFields(t *testing.T) {
	t.Parallel()

	base := func() StreamConfiguration {
		return StreamConfiguration{
			StreamName: "TestStream",
			DeliveryConfiguration: DeliveryConfiguration{
				Delimiter: DelimiterTypePtr(DelimiterTypeSpace),
				Format:    FormatTypeStructured,
				Frequency: Frequency{IntervalInSeconds: IntervalInSeconds30},
			},
			Destination: AbstractConnector(&S3Connector{
				Path:            "path",
				DisplayName:     "display",
				Bucket:          "bucket",
				Region:          "us-east-1",
				AccessKey:       "key",
				SecretAccessKey: "secret",
			}),
			DatasetFields: []DatasetFieldID{{DatasetFieldID: 1}},
			Properties:    []PropertyID{{PropertyID: 1}},
		}
	}

	tests := map[string]struct {
		mutate      func(*StreamConfiguration)
		contains    []string
		notContains []string
		checkMap    func(*testing.T, map[string]any)
		checkBody   func(*testing.T, []byte)
	}{
		"marshal omits zero contractId and groupId": {
			mutate: func(c *StreamConfiguration) {
				c.ContractID = ""
				c.GroupID = 0
			},
			notContains: []string{"contractId", "groupId"},
			checkMap: func(t *testing.T, decoded map[string]any) {
				assert.Contains(t, decoded, "streamName")
				assert.Contains(t, decoded, "destination")
				assert.Contains(t, decoded, "deliveryConfiguration")
			},
		},
		"marshal includes populated contractId and groupId": {
			mutate: func(c *StreamConfiguration) {
				c.ContractID = "2-AB1234"
				c.GroupID = 1234
			},
			contains: []string{`"contractId":"2-AB1234"`, `"groupId":1234`},
		},
		"marshal only group omitted": {
			mutate: func(c *StreamConfiguration) {
				c.ContractID = "2-AB1234"
				c.GroupID = 0
			},
			contains:    []string{`"contractId":"2-AB1234"`},
			notContains: []string{"groupId"},
		},
		"marshal only contract omitted": {
			mutate: func(c *StreamConfiguration) {
				c.ContractID = ""
				c.GroupID = 1234
			},
			contains:    []string{`"groupId":1234`},
			notContains: []string{"contractId"},
		},
		"marshal includes whitespace contractId": {
			mutate: func(c *StreamConfiguration) {
				c.ContractID = "   "
				c.GroupID = 0
			},
			contains:    []string{`"contractId":"   "`},
			notContains: []string{"groupId"},
		},
		"marshal round-trip preserves optional field keys": {
			mutate: func(c *StreamConfiguration) {
				c.ContractID = "P-999"
				c.GroupID = 77
			},
			checkMap: func(t *testing.T, decoded map[string]any) {
				assert.Equal(t, "P-999", decoded["contractId"])
				assert.Equal(t, float64(77), decoded["groupId"])
			},
		},
		"marshal omits optional keys when zero after population cleared": {
			mutate: func(c *StreamConfiguration) {
				c.ContractID = "P-999"
				c.GroupID = 77
				c.ContractID = ""
				c.GroupID = 0
			},
			checkMap: func(t *testing.T, decoded map[string]any) {
				assert.NotContains(t, decoded, "contractId")
				assert.NotContains(t, decoded, "groupId")
			},
		},
		"deterministic marshal order": {
			mutate: func(c *StreamConfiguration) {
				c.ContractID = "C-1"
				c.GroupID = 9
			},
			checkBody: func(t *testing.T, first []byte) {
				cfg := base()
				cfg.ContractID = "C-1"
				cfg.GroupID = 9
				second, err := json.Marshal(cfg)
				require.NoError(t, err)
				assert.Contains(t, string(first), "contractId")
				assert.Contains(t, string(first), "groupId")
				assert.Equal(t, string(first), string(second))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cfg := base()
			if tc.mutate != nil {
				tc.mutate(&cfg)
			}
			data, err := json.Marshal(cfg)
			require.NoError(t, err)
			body := string(data)
			for _, s := range tc.contains {
				assert.Contains(t, body, s)
			}
			for _, s := range tc.notContains {
				assert.NotContains(t, body, s)
			}
			if tc.checkMap != nil {
				var decoded map[string]any
				require.NoError(t, json.Unmarshal(data, &decoded))
				tc.checkMap(t, decoded)
			}
			if tc.checkBody != nil {
				tc.checkBody(t, data)
			}
		})
	}
}

func TestStreamResponseTypes_JSONOptionalFields(t *testing.T) {
	t.Parallel()

	const detailedBase = `{"streamId":7050,"streamName":"TestStream","streamStatus":"ACTIVATED","streamVersion":1,"latestVersion":1,"createdBy":"u","createdDate":"2022-11-04T00:49:45Z","modifiedBy":"u","modifiedDate":"2022-11-04T02:14:29Z","productId":"p","datasetFields":[],"deliveryConfiguration":{"format":"JSON","frequency":{"intervalInSeconds":30}},"destination":{"destinationType":"S3"},"notificationEmails":[],"properties":[]}`
	const detailsBase = `{"streamId":1,"streamName":"n","streamStatus":"ACTIVATED","streamVersion":1,"latestVersion":1,"createdBy":"u","createdDate":"d","modifiedBy":"u","modifiedDate":"d","productId":"p","properties":[],"appSecConfigs":[]}`

	tests := map[string]struct {
		check func(*testing.T)
	}{
		"DetailedStreamVersion unmarshals without optional fields": {
			check: func(t *testing.T) {
				var v DetailedStreamVersion
				require.NoError(t, json.Unmarshal([]byte(detailedBase), &v))
				assert.Empty(t, v.ContractID)
				assert.Zero(t, v.GroupID)
				assert.Equal(t, int64(7050), v.StreamID)
			},
		},
		"StreamDetails unmarshals without optional fields": {
			check: func(t *testing.T) {
				var s StreamDetails
				require.NoError(t, json.Unmarshal([]byte(detailsBase), &s))
				assert.Empty(t, s.ContractID)
				assert.Zero(t, s.GroupID)
			},
		},
		"DetailedStreamVersion unmarshals empty contractId and zero groupId": {
			check: func(t *testing.T) {
				raw := `{"contractId":"","groupId":0,"streamId":7050,"streamName":"TestStream","streamStatus":"ACTIVATED","streamVersion":1,"latestVersion":1,"createdBy":"u","createdDate":"2022-11-04T00:49:45Z","modifiedBy":"u","modifiedDate":"2022-11-04T02:14:29Z","productId":"p","datasetFields":[],"deliveryConfiguration":{"format":"JSON","frequency":{"intervalInSeconds":30}},"destination":{"destinationType":"S3"},"notificationEmails":[],"properties":[]}`
				var v DetailedStreamVersion
				require.NoError(t, json.Unmarshal([]byte(raw), &v))
				assert.Empty(t, v.ContractID)
				assert.Zero(t, v.GroupID)
			},
		},
		"StreamDetails unmarshals empty contractId and zero groupId": {
			check: func(t *testing.T) {
				raw := `{"contractId":"","groupId":0,"streamId":1,"streamName":"n","streamStatus":"ACTIVATED","streamVersion":1,"latestVersion":1,"createdBy":"u","createdDate":"d","modifiedBy":"u","modifiedDate":"d","productId":"p","properties":[],"appSecConfigs":[]}`
				var s StreamDetails
				require.NoError(t, json.Unmarshal([]byte(raw), &s))
				assert.Empty(t, s.ContractID)
				assert.Zero(t, s.GroupID)
			},
		},
		"DetailedStreamVersion and StreamDetails unmarshal populated optional fields": {
			check: func(t *testing.T) {
				raw := `{
					"streamId": 1,
					"streamName": "n",
					"streamStatus": "ACTIVATED",
					"streamVersion": 1,
					"latestVersion": 1,
					"contractId": "1-ABCDE",
					"groupId": 211516,
					"createdBy": "u",
					"createdDate": "d",
					"modifiedBy": "u",
					"modifiedDate": "d",
					"productId": "p",
					"properties": [],
					"appSecConfigs": []
				}`
				var listed StreamDetails
				require.NoError(t, json.Unmarshal([]byte(raw), &listed))
				assert.Equal(t, "1-ABCDE", listed.ContractID)
				assert.Equal(t, 211516, listed.GroupID)

				var detailed DetailedStreamVersion
				require.NoError(t, json.Unmarshal([]byte(raw), &detailed))
				assert.Equal(t, "1-ABCDE", detailed.ContractID)
				assert.Equal(t, 211516, detailed.GroupID)
			},
		},
		"DetailedStreamVersion marshal includes populated optional fields": {
			check: func(t *testing.T) {
				v := DetailedStreamVersion{
					StreamID:   1,
					StreamName: "s",
					ContractID: "C-1",
					GroupID:    42,
				}
				data, err := json.Marshal(v)
				require.NoError(t, err)
				body := string(data)
				assert.Contains(t, body, `"contractId":"C-1"`)
				assert.Contains(t, body, `"groupId":42`)
			},
		},
		"StreamDetails marshal includes populated optional fields": {
			check: func(t *testing.T) {
				s := StreamDetails{
					StreamID:   99,
					StreamName: "stream",
					ContractID: "C-1",
					GroupID:    12,
				}
				data, err := json.Marshal(s)
				require.NoError(t, err)
				body := string(data)
				assert.Contains(t, body, `"contractId":"C-1"`)
				assert.Contains(t, body, `"groupId":12`)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tc.check(t)
		})
	}
}

func validAppSecUpdateStreamRequest() UpdateStreamRequest {
	req := validCDNUpdateStreamRequest()
	req.LogType = LogTypeAppSec
	req.StreamConfiguration.Properties = []PropertyID{}
	req.StreamConfiguration.DatasetFields = []DatasetFieldID{}
	req.StreamConfiguration.AppSecConfigs = []AppSecConfigID{{AppSecID: 12345}}
	return req
}

func TestStreamOptionalFields_ConcurrentMarshal(t *testing.T) {
	t.Parallel()
	cfg := StreamConfiguration{
		StreamName: "TestStream",
		ContractID: "",
		GroupID:    0,
		DeliveryConfiguration: DeliveryConfiguration{
			Format:    FormatTypeJson,
			Frequency: Frequency{IntervalInSeconds: IntervalInSeconds30},
		},
		Destination: AbstractConnector(&S3Connector{
			Path:            "path",
			DisplayName:     "display",
			Bucket:          "bucket",
			Region:          "us-east-1",
			AccessKey:       "key",
			SecretAccessKey: "secret",
		}),
	}

	const workers = 32
	errCh := make(chan error, workers)
	for i := 0; i < workers; i++ {
		go func() {
			data, err := json.Marshal(cfg)
			if err != nil {
				errCh <- err
				return
			}
			body := string(data)
			if strings.Contains(body, "contractId") || strings.Contains(body, "groupId") {
				errCh <- errors.New("optional fields must be omitted")
				return
			}
			errCh <- nil
		}()
	}
	for i := 0; i < workers; i++ {
		require.NoError(t, <-errCh)
	}
}
