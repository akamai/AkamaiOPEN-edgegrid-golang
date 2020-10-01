package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapi_GetGroups(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetGroupsResponse
		withError        error
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
{
	"accountId": "act_1-1TJZFB",
	"accountName": "Example.com",
	"groups": {
		"items": [
			{
				"groupName": "Example.com-1-1TJZH5",
				"groupId": "grp_15225",
				"contractIds": [
					"ctr_1-1TJZH5"
				]
			}
		]
	}
}`,
			expectedPath: "/papi/v1/groups",
			expectedResponse: &GetGroupsResponse{
				AccountID:   "act_1-1TJZFB",
				AccountName: "Example.com",
				Groups: GroupItems{Items: []*Group{
					{
						GroupName:   "Example.com-1-1TJZH5",
						GroupID:     "grp_15225",
						ContractIDs: []string{"ctr_1-1TJZH5"},
					},
				}},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching groups",
    "status": 500
}`,
			expectedPath: "/papi/v1/groups",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching groups",
				StatusCode: http.StatusInternalServerError,
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
			result, err := client.GetGroups(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
