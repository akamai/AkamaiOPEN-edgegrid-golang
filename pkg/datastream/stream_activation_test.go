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

func TestDs_ActivateStream(t *testing.T) {
	tests := map[string]struct {
		request          ActivateStreamRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ActivateStreamResponse
		withError        func(*testing.T, error)
	}{
		"202 accepted": {
			request:        ActivateStreamRequest{StreamID: 3},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "streamVersionKey": {
        "streamId": 1,
        "streamVersionId": 3
    }
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams/3/activate",
			expectedResponse: &ActivateStreamResponse{
				StreamVersionKey: StreamVersionKey{
					StreamID:        1,
					StreamVersionID: 3,
				},
			},
		},
		"validation error": {
			request: ActivateStreamRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
				test.withError(t, err)
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
		withError        func(*testing.T, error)
	}{
		"202 accepted": {
			request:        DeactivateStreamRequest{StreamID: 3},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "streamVersionKey": {
        "streamId": 1,
        "streamVersionId": 3
    }
}
`,
			expectedPath: "/datastream-config-api/v1/log/streams/3/deactivate",
			expectedResponse: &DeactivateStreamResponse{
				StreamVersionKey: StreamVersionKey{
					StreamID:        1,
					StreamVersionID: 3,
				},
			},
		},
		"validation error": {
			request: DeactivateStreamRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
				test.withError(t, err)
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
		withError        func(*testing.T, error)
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
			request: GetActivationHistoryRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
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
			result, err := client.GetActivationHistory(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
