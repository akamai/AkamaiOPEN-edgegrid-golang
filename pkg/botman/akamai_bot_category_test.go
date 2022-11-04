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

// Test Get AkamaiBotCategory List
func TestBotman_GetAkamaiBotCategoryList(t *testing.T) {

	tests := map[string]struct {
		params           GetAkamaiBotCategoryListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAkamaiBotCategoryListResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
{
	"categories": [
		{"categoryId":"b85e3eaa-d334-466d-857e-33308ce416be", "categoryName":"Test Name 1", "testKey":"testValue1"},
		{"categoryId":"69acad64-7459-4c1d-9bad-672600150127", "categoryName":"Test Name 2", "testKey":"testValue2"},
		{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "categoryName":"Test Name 3", "testKey":"testValue3"},
		{"categoryId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "categoryName":"Test Name 4", "testKey":"testValue4"},
		{"categoryId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "categoryName":"Test Name 5", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/akamai-bot-categories",
			expectedResponse: &GetAkamaiBotCategoryListResponse{
				Categories: []map[string]interface{}{
					{"categoryId": "b85e3eaa-d334-466d-857e-33308ce416be", "categoryName": "Test Name 1", "testKey": "testValue1"},
					{"categoryId": "69acad64-7459-4c1d-9bad-672600150127", "categoryName": "Test Name 2", "testKey": "testValue2"},
					{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "categoryName": "Test Name 3", "testKey": "testValue3"},
					{"categoryId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "categoryName": "Test Name 4", "testKey": "testValue4"},
					{"categoryId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "categoryName": "Test Name 5", "testKey": "testValue5"},
				},
			},
		},
		"200 OK One Record": {
			params: GetAkamaiBotCategoryListRequest{
				CategoryName: "Test Name 3",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"categories":[
		{"categoryId":"b85e3eaa-d334-466d-857e-33308ce416be", "categoryName":"Test Name 1", "testKey":"testValue1"},
		{"categoryId":"69acad64-7459-4c1d-9bad-672600150127", "categoryName":"Test Name 2", "testKey":"testValue2"},
		{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "categoryName":"Test Name 3", "testKey":"testValue3"},
		{"categoryId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "categoryName":"Test Name 4", "testKey":"testValue4"},
		{"categoryId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "categoryName":"Test Name 5", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/akamai-bot-categories",
			expectedResponse: &GetAkamaiBotCategoryListResponse{
				Categories: []map[string]interface{}{
					{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "categoryName": "Test Name 3", "testKey": "testValue3"},
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
			expectedPath: "/appsec/v1/akamai-bot-categories",
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
			result, err := client.GetAkamaiBotCategoryList(
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
