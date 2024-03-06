package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDNS_ListGroups(t *testing.T) {
	tests := map[string]struct {
		request          ListGroupRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListGroupResponse
		withError        error
	}{
		"200 OK, when optional query parameter provided": {
			request: ListGroupRequest{
				GroupID: "9012",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
  				"groups": [
    				{
      					"groupId": 9012,
      					"groupName": "example-name",
      					"contractIds": [
        					"1-2ABCDE"
      					],
      					"permissions": [
        					"READ",
        					"WRITE",
        					"ADD",
        					"DELETE"
      					]
    				}
  				]
			}`,
			expectedPath: "/config-dns/v2/data/groups/?gid=9012",
			expectedResponse: &ListGroupResponse{
				Groups: []Group{
					{
						GroupID:   9012,
						GroupName: "example-name",
						ContractIDs: []string{
							"1-2ABCDE",
						},
						Permissions: []string{
							"READ",
							"WRITE",
							"ADD",
							"DELETE",
						},
					},
				},
			},
		},
		"200 OK, when optional query parameter not provided": {
			responseStatus: http.StatusOK,
			responseBody: `
			{
  				"groups": [
    				{
      					"groupId": 9012,
      					"groupName": "example-name1",
      					"contractIds": [
        					"1-2ABCDE"
      					],
      					"permissions": [
        					"READ",
        					"WRITE",
        					"ADD",
        					"DELETE"
      					]
    				},
{
      					"groupId": 9013,
      					"groupName": "example-name2",
      					"contractIds": [
        					"1-2ABCDE"
      					],
      					"permissions": [
        					"READ",
        					"WRITE",
        					"ADD",
        					"DELETE"
      					]
    				}
  				]
			}`,
			expectedPath: "/config-dns/v2/data/groups/",
			expectedResponse: &ListGroupResponse{
				Groups: []Group{
					{
						GroupID:   9012,
						GroupName: "example-name1",
						ContractIDs: []string{
							"1-2ABCDE",
						},
						Permissions: []string{
							"READ",
							"WRITE",
							"ADD",
							"DELETE",
						},
					},
					{
						GroupID:   9013,
						GroupName: "example-name2",
						ContractIDs: []string{
							"1-2ABCDE",
						},
						Permissions: []string{
							"READ",
							"WRITE",
							"ADD",
							"DELETE",
						},
					},
				},
			},
		},
		"500 internal server error, when optional query parameter not provided ": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
    				"detail": "Error fetching authorities",
    				"status": 500
				}`,
			expectedPath: "/config-dns/v2/data/groups/",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
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
			result, err := client.ListGroups(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
