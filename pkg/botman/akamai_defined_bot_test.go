package botman

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Get AkamaiDefinedBot List
func TestBotman_GetAkamaiDefinedBotList(t *testing.T) {

	tests := map[string]struct {
		params           GetAkamaiDefinedBotListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAkamaiDefinedBotListResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
{
	"bots": [
		{"botId":"b85e3eaa-d334-466d-857e-33308ce416be", "botName":"Test Name 1", "testKey":"testValue1"},
		{"botId":"69acad64-7459-4c1d-9bad-672600150127", "botName":"Test Name 2", "testKey":"testValue2"},
		{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "botName":"Test Name 3", "testKey":"testValue3"},
		{"botId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "botName":"Test Name 4", "testKey":"testValue4"},
		{"botId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "botName":"Test Name 5", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/akamai-defined-bots",
			expectedResponse: &GetAkamaiDefinedBotListResponse{
				Bots: []map[string]interface{}{
					{"botId": "b85e3eaa-d334-466d-857e-33308ce416be", "botName": "Test Name 1", "testKey": "testValue1"},
					{"botId": "69acad64-7459-4c1d-9bad-672600150127", "botName": "Test Name 2", "testKey": "testValue2"},
					{"botId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "botName": "Test Name 3", "testKey": "testValue3"},
					{"botId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "botName": "Test Name 4", "testKey": "testValue4"},
					{"botId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "botName": "Test Name 5", "testKey": "testValue5"},
				},
			},
		},
		"200 OK One Record": {
			params: GetAkamaiDefinedBotListRequest{
				BotName: "Test Name 3",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"bots":[
		{"botId":"b85e3eaa-d334-466d-857e-33308ce416be", "botName":"Test Name 1", "testKey":"testValue1"},
		{"botId":"69acad64-7459-4c1d-9bad-672600150127", "botName":"Test Name 2", "testKey":"testValue2"},
		{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "botName":"Test Name 3", "testKey":"testValue3"},
		{"botId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "botName":"Test Name 4", "testKey":"testValue4"},
		{"botId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "botName":"Test Name 5", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/akamai-defined-bots",
			expectedResponse: &GetAkamaiDefinedBotListResponse{
				Bots: []map[string]interface{}{
					{"botId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "botName": "Test Name 3", "testKey": "testValue3"},
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching data",
    "status": 500
}`,
			expectedPath: "/appsec/v1/akamai-defined-bots",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching data",
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
			result, err := client.GetAkamaiDefinedBotList(
				session.ContextWithOptions(
					context.Background(),
				),
				test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
