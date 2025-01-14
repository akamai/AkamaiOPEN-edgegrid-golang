package datastream

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/ptr"
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
		"200 OK Without midgress field": {
			request: GetStreamRequest{
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
        "frequency": {
            "intervalInSeconds": 30
        }, 
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
			expectedPath: "/datastream-config-api/v2/log/streams/1",
			expectedResponse: &DetailedStreamVersion{
				StreamStatus: StreamStatusActivated,
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter: DelimiterTypePtr(DelimiterTypeSpace),
					Format:    FormatTypeStructured,
					Frequency: Frequency{
						IntervalInSeconds: IntervalInSeconds30,
					},
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
					{
						DatasetFieldID:      1000,
						DatasetFieldName:    "dataset_field_name_1",
						DatasetFieldJsonKey: "dataset_field_json_key_1",
					},
					{
						DatasetFieldID:      1002,
						DatasetFieldName:    "dataset_field_name_2",
						DatasetFieldJsonKey: "dataset_field_json_key_2",
					},
					{
						DatasetFieldID:      1082,
						DatasetFieldName:    "dataset_field_name_3",
						DatasetFieldJsonKey: "dataset_field_json_key_3",
					},
				},
				NotificationEmails: []string{"sample_username@akamai.com"},
				GroupID:            1234,
				ModifiedBy:         "sample_username2",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				ProductID:          "Adaptive_Media_Delivery",
				Properties: []Property{
					{
						PropertyID:   12345,
						PropertyName: "example.com",
					},
				},
				StreamID:      1,
				StreamName:    "ds2-sample-name",
				StreamVersion: 2,
				LatestVersion: 2,
			},
		},

		"200 OK With midgress field": {
			request: GetStreamRequest{
				StreamID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "P-1324", 
    "createdBy": "sample_username", 
    "createdDate": "2022-11-04T00:49:45Z", 
    "collectMidgress": true,
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
        "frequency": {
            "intervalInSeconds": 30
        }, 
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
            "propertyId": 1234, 
            "propertyName": "sample.com"
        }
    ], 
    "streamId": 1, 
    "streamName": "ds2-sample-name", 
    "streamStatus": "ACTIVATED", 
    "streamVersion": 2
}
`,
			expectedPath: "/datastream-config-api/v2/log/streams/1",
			expectedResponse: &DetailedStreamVersion{
				CollectMidgress: true,
				StreamStatus:    StreamStatusActivated,
				DeliveryConfiguration: DeliveryConfiguration{
					Delimiter: DelimiterTypePtr(DelimiterTypeSpace),
					Format:    FormatTypeStructured,
					Frequency: Frequency{
						IntervalInSeconds: IntervalInSeconds30,
					},
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
					{
						DatasetFieldID:      1000,
						DatasetFieldName:    "dataset_field_name_1",
						DatasetFieldJsonKey: "dataset_field_json_key_1",
					},
					{
						DatasetFieldID:      1002,
						DatasetFieldName:    "dataset_field_name_2",
						DatasetFieldJsonKey: "dataset_field_json_key_2",
					},
					{
						DatasetFieldID:      1082,
						DatasetFieldName:    "dataset_field_name_3",
						DatasetFieldJsonKey: "dataset_field_json_key_3",
					},
				},
				NotificationEmails: []string{"sample_username@akamai.com"},
				GroupID:            1234,
				ModifiedBy:         "sample_username2",
				ModifiedDate:       "2022-11-04T02:14:29Z",
				ProductID:          "Adaptive_Media_Delivery",
				Properties: []Property{
					{
						PropertyID:   1234,
						PropertyName: "sample.com",
					},
				},
				StreamID:      1,
				StreamName:    "ds2-sample-name",
				StreamVersion: 2,
				LatestVersion: 2,
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
			expectedPath:   "/datastream-config-api/v2/log/streams/12",
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
				{
					DatasetFieldID: 2020,
				},
			},
			NotificationEmails: []string{"useremail1@akamai.com", "useremail2@akamai.com"},
			GroupID:            1234,
			Properties: []PropertyID{
				{
					PropertyID: 1234,
				},
				{
					PropertyID: 1234,
				},
			},
			StreamName:      "TestStream",
			CollectMidgress: true,
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
		withError        error
	}{
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
			expectedPath: "/datastream-config-api/v2/log/streams?activate=true",
			expectedResponse: &DetailedStreamVersion{
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
      {
         "propertyId":1234
      },
      {
         "propertyId":1234
      }
   ],
   "datasetFields":[
      {
         "datasetFieldId":2020
      }
   ],
   "deliveryConfiguration":{
      "uploadFilePrefix":"logs",
      "uploadFileSuffix":"ak",
      "fieldDelimiter":"SPACE",
      "format":"STRUCTURED",
      "frequency":{
         "intervalInSeconds":30
      }
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
			expectedPath: "/datastream-config-api/v2/log/streams?activate=true",
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
			expectedPath: "/datastream-config-api/v2/log/streams?activate=true",
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
		Activate: true,
		StreamConfiguration: StreamConfiguration{
			DeliveryConfiguration: DeliveryConfiguration{
				Delimiter:        DelimiterTypePtr(DelimiterTypeSpace),
				Format:           "STRUCTURED",
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
			NotificationEmails: []string{"test@aka.mai", "useremail2@akamai.com"},

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

	modifyRequest := func(r UpdateStreamRequest, opt func(r *UpdateStreamRequest)) UpdateStreamRequest {
		opt(&r)
		return r
	}

	tests := map[string]struct {
		request          UpdateStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DetailedStreamVersion
		withError        error
	}{
		"200 OK activate:true": {
			request:        updateRequest,
			responseStatus: http.StatusOK,
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
			expectedPath: "/datastream-config-api/v2/log/streams/7050?activate=true",
			expectedResponse: &DetailedStreamVersion{
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
		"validation error - groupId modification": {
			request: modifyRequest(updateRequest, func(r *UpdateStreamRequest) {
				r.StreamConfiguration.GroupID = 1337
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
			expectedPath: "/datastream-config-api/v2/log/streams/7050?activate=true",
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
		request        DeleteStreamRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			request: DeleteStreamRequest{
				StreamID: 1,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   ``,
			expectedPath:   "/datastream-config-api/v2/log/streams/1",
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
			expectedPath:   "/datastream-config-api/v2/log/streams/12",
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
			expectedPath: "/datastream-config-api/v2/log/streams",
			expectedResponse: []StreamDetails{
				{
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
			expectedPath: "/datastream-config-api/v2/log/streams?groupId=1234",
			expectedResponse: []StreamDetails{
				{
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
		"400 bad request": {
			request:        ListStreamsRequest{},
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/datastream-config-api/v2/log/streams",
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
