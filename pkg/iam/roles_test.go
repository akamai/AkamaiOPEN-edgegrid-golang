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

func TestIAM_ListRoles(t *testing.T) {
	tests := map[string]struct {
		params           ListRolesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Role
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListRolesRequest{
				Actions: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "roleId": 123456,
        "roleName": "View Only",
        "roleDescription": "This role will allow you to view",
        "type": "custom",
        "createdDate": "2017-07-27T18:11:25.000Z",
        "createdBy": "john.doe@mycompany.com",
        "modifiedDate": "2017-07-27T18:11:25.000Z",
        "modifiedBy": "john.doe@mycompany.com",
        "actions": {
            "edit": true,
            "delete": true
        }
	}
]`,
			expectedPath: "/identity-management/v2/user-admin/roles?actions=true&ignoreContext=false&users=false",
			expectedResponse: []Role{
				{
					RoleID:          123456,
					RoleName:        "View Only",
					RoleDescription: "This role will allow you to view",
					RoleType:        RoleTypeCustom,
					CreatedDate:     "2017-07-27T18:11:25.000Z",
					CreatedBy:       "john.doe@mycompany.com",
					ModifiedDate:    "2017-07-27T18:11:25.000Z",
					ModifiedBy:      "john.doe@mycompany.com",
					Actions: &RoleAction{
						Edit:   true,
						Delete: true,
					},
				},
			},
		},
		"500 internal server error": {
			params: ListRolesRequest{
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
			expectedPath: "/identity-management/v2/user-admin/roles?actions=true&ignoreContext=false&users=false",
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
			result, err := client.ListRoles(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
