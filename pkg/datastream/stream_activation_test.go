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

func TestDs_ActivateStream(t *testing.T) {
	tests := map[string]struct {
		request          ActivateStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DetailedStreamVersion
		withError        error
	}{
		"200 OK": {
			request:        ActivateStreamRequest{StreamID: 3},
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
    "streamId": 3, 
    "streamName": "ds2-sample-name", 
    "streamStatus": "ACTIVATING", 
    "streamVersion": 2
}
`,
			expectedPath: "/datastream-config-api/v2/log/streams/3/activate",
			expectedResponse: &DetailedStreamVersion{
				CollectMidgress: true,
				StreamStatus:    StreamStatusActivating,
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
				StreamID:      3,
				StreamName:    "ds2-sample-name",
				StreamVersion: 2,
				LatestVersion: 2,
			},
		},
		"validation error": {
			request:   ActivateStreamRequest{},
			withError: ErrStructValidation,
		},
		"400 bad request": {
			request:        ActivateStreamRequest{StreamID: 123},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
	"type": "bad-request",
	"title": "Bad Request",
	"detail": "",
	"instance": "df22bc0f-ca8d-4bdb-afea-ffdeef819e22",
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
			expectedPath: "/datastream-config-api/v2/log/streams/123/activate",
			withError: &Error{
				Type:       "bad-request",
				Title:      "Bad Request",
				Instance:   "df22bc0f-ca8d-4bdb-afea-ffdeef819e22",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ActivateStream(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_DeactivateStream(t *testing.T) {
	tests := map[string]struct {
		request          DeactivateStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DetailedStreamVersion
		withError        error
	}{
		"200 ok": {
			request:        DeactivateStreamRequest{StreamID: 3},
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
    "streamId": 3, 
    "streamName": "ds2-sample-name", 
    "streamStatus": "DEACTIVATING", 
    "streamVersion": 2
}
`,
			expectedPath: "/datastream-config-api/v2/log/streams/3/deactivate",
			expectedResponse: &DetailedStreamVersion{
				CollectMidgress: true,
				StreamStatus:    StreamStatusDeactivating,
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
				StreamID:      3,
				StreamName:    "ds2-sample-name",
				StreamVersion: 2,
				LatestVersion: 2,
			},
		},
		"validation error": {
			request:   DeactivateStreamRequest{},
			withError: ErrStructValidation,
		},
		"400 bad request": {
			request:        DeactivateStreamRequest{StreamID: 123},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
	"type": "bad-request",
	"title": "Bad Request",
	"detail": "",
	"instance": "df22bc0f-ca8d-4bdb-afea-ffdeef819e22",
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
			expectedPath: "/datastream-config-api/v2/log/streams/123/deactivate",
			withError: &Error{
				Type:       "bad-request",
				Title:      "Bad Request",
				Instance:   "df22bc0f-ca8d-4bdb-afea-ffdeef819e22",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeactivateStream(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDs_GetActivationHistory(t *testing.T) {
	tests := map[string]struct {
		request          GetActivationHistoryRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []ActivationHistoryEntry
		withError        error
	}{
		"200 OK": {
			request:        GetActivationHistoryRequest{StreamID: 3},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "streamId": 3,
        "streamVersion": 2,
        "modifiedBy": "user1",
        "modifiedDate": "16-01-2020 11:07:12 GMT",
        "status": "DEACTIVATED"
    },
    {
        "streamId": 3,
        "streamVersion": 2,
        "modifiedBy": "user2",
        "modifiedDate": "16-01-2020 09:31:02 GMT",
        "status": "ACTIVATED"
    }
]
`,
			expectedPath: "/datastream-config-api/v2/log/streams/3/activation-history",
			expectedResponse: []ActivationHistoryEntry{
				{
					ModifiedBy:    "user1",
					ModifiedDate:  "16-01-2020 11:07:12 GMT",
					Status:        StreamStatusDeactivated,
					StreamID:      3,
					StreamVersion: 2,
				},
				{
					ModifiedBy:    "user2",
					ModifiedDate:  "16-01-2020 09:31:02 GMT",
					Status:        StreamStatusActivated,
					StreamID:      3,
					StreamVersion: 2,
				},
			},
		},
		"validation error": {
			request:   GetActivationHistoryRequest{},
			withError: ErrStructValidation,
		},
		"400 bad request": {
			request:        GetActivationHistoryRequest{StreamID: 123},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
	"type": "bad-request",
	"title": "Bad Request",
	"detail": "",
	"instance": "df22bc0f-ca8d-4bdb-afea-ffdeef819e22",
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
			expectedPath: "/datastream-config-api/v2/log/streams/123/activation-history",
			withError: &Error{
				Type:       "bad-request",
				Title:      "Bad Request",
				Instance:   "df22bc0f-ca8d-4bdb-afea-ffdeef819e22",
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
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetActivationHistory(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
