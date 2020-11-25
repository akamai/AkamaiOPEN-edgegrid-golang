package iam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestIAM_ListGroups(t *testing.T) {
	tests := map[string]struct {
		params           ListGroupsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Group
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListGroupsRequest{
				Actions: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
			[
				{
					"groupId": 12345,
					"groupName": "Top Level group",
					"createdDate": "2012-04-28T00:00:00.000Z",
					"createdBy": "johndoe",
					"modifiedDate": "2012-04-28T00:00:00.000Z",
					"modifiedBy": "johndoe",
					"actions": {
						"edit": true,
						"delete": false
					}
				}
			]`,
			expectedPath: "/identity-management/v2/user-admin/groups?actions=true",
			expectedResponse: []Group{
				{
					GroupID:      12345,
					GroupName:    "Top Level group",
					CreatedDate:  "2012-04-28T00:00:00.000Z",
					CreatedBy:    "johndoe",
					ModifiedDate: "2012-04-28T00:00:00.000Z",
					ModifiedBy:   "johndoe",
					Actions: &GroupActions{
						Edit:   true,
						Delete: false,
					},
				},
			},
		},
		"500 internal server error": {
			params: ListGroupsRequest{
				Actions: true,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching cp codes",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/groups?actions=true",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching cp codes",
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
			result, err := client.ListGroups(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
