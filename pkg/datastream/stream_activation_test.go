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
		expectedResponse *ActivateStreamResponse
		withError        error
	}{
		"202 accepted": {
			request:        ActivateStreamRequest{StreamID: 3},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "streamVersionKey": {
        "streamId": 1,
        "streamVersion": 3
    }
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams/3/activate",
			expectedResponse: &ActivateStreamResponse{
				StreamVersionKey: StreamUpdate{
					StreamID:      1,
					StreamVersion: 3,
				},
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
			expectedPath: "/datastream-config-api/v1/log/streams/123/activate",
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
				assert.Equal(t, http.MethodPut, r.Method)
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
		expectedResponse *DeactivateStreamResponse
		withError        error
	}{
		"202 accepted": {
			request:        DeactivateStreamRequest{StreamID: 3},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "streamVersionKey": {
        "streamId": 1,
        "streamVersion": 3
    }
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams/3/deactivate",
			expectedResponse: &DeactivateStreamResponse{
				StreamVersionKey: StreamUpdate{
					StreamID:      1,
					StreamVersion: 3,
				},
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
			expectedPath: "/datastream-config-api/v1/log/streams/123/deactivate",
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
				assert.Equal(t, http.MethodPut, r.Method)
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
        "streamId": 7050,
        "streamVersionId": 2,
        "createdBy": "user1",
        "createdDate": "16-01-2020 11:07:12 GMT",
        "isActive": false
    },
    {
        "streamId": 7050,
        "streamVersionId": 2,
        "createdBy": "user2",
        "createdDate": "16-01-2020 09:31:02 GMT",
        "isActive": true
    }
]
`,
			expectedPath: "/datastream-config-api/v1/log/streams/3/activationHistory",
			expectedResponse: []ActivationHistoryEntry{
				{
					CreatedBy:       "user1",
					CreatedDate:     "16-01-2020 11:07:12 GMT",
					IsActive:        false,
					StreamID:        7050,
					StreamVersionID: 2,
				},
				{
					CreatedBy:       "user2",
					CreatedDate:     "16-01-2020 09:31:02 GMT",
					IsActive:        true,
					StreamID:        7050,
					StreamVersionID: 2,
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
			expectedPath: "/datastream-config-api/v1/log/streams/123/activationHistory",
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
