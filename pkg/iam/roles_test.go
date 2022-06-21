package iam

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestIAM_CreateRole(t *testing.T) {
	tests := map[string]struct {
		params              CreateRoleRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *Role
		withError           error
	}{
		"201 Created": {
			params: CreateRoleRequest{
				Name:         "Terraform admin",
				Description:  "Admin granted role for tests",
				GrantedRoles: []GrantedRoleID{{ID: 12345}},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "roleId": 123456,
    "roleName": "Terraform admin",
    "roleDescription": "Admin granted role for tests",
    "type": "custom",
    "createdDate": "2022-04-11T10:52:03.811Z",
    "createdBy": "jBond",
    "modifiedDate": "2022-04-11T10:52:03.811Z",
    "modifiedBy": "jBond",
    "actions": {
        "edit": true,
        "delete": true
    },
    "grantedRoles": [
        {
            "grantedRoleId": 12345,
            "grantedRoleName": "WebAP User",
            "grantedRoleDescription": "Web Application Protector User Role"
        }
    ]
}`,
			expectedPath:        "/identity-management/v2/user-admin/roles",
			expectedRequestBody: `{"roleName":"Terraform admin","roleDescription":"Admin granted role for tests","grantedRoles":[{"grantedRoleId":12345}]}`,
			expectedResponse: &Role{
				RoleID:          123456,
				RoleName:        "Terraform admin",
				RoleDescription: "Admin granted role for tests",
				RoleType:        RoleTypeCustom,
				CreatedDate:     "2022-04-11T10:52:03.811Z",
				CreatedBy:       "jBond",
				ModifiedDate:    "2022-04-11T10:52:03.811Z",
				ModifiedBy:      "jBond",
				Actions: &RoleAction{
					Edit:   true,
					Delete: true,
				},
				GrantedRoles: []RoleGrantedRole{
					{
						RoleID:      12345,
						RoleName:    "WebAP User",
						Description: "Web Application Protector User Role",
					},
				},
			},
		},
		"500 Internal server error": {
			params: CreateRoleRequest{
				Name:         "Terraform admin",
				Description:  "Admin granted role for tests",
				GrantedRoles: []GrantedRoleID{{ID: 12345}},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/roles",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
				StatusCode: http.StatusInternalServerError,
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateRole(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAM_GetRole(t *testing.T) {
	tests := map[string]struct {
		params           GetRoleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Role
		withError        error
	}{
		"200 OK with query params": {
			params: GetRoleRequest{
				ID:           123456,
				Actions:      true,
				GrantedRoles: true,
				Users:        true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "roleId": 123456,
    "roleName": "Terraform admin updated",
    "roleDescription": "Admin granted role for tests",
    "type": "custom",
    "createdDate": "2022-04-11T10:52:03.000Z",
    "createdBy": "jBond",
    "modifiedDate": "2022-04-11T10:59:30.000Z",
    "modifiedBy": "jBond",
    "actions": {
        "edit": true,
        "delete": true
    },
    "grantedRoles": [
        {
            "grantedRoleId": 12345,
            "grantedRoleName": "View Audience Analytics Reports",
            "grantedRoleDescription": "Publisher Self-Provisioning"
        },
        {
            "grantedRoleId": 54321,
            "grantedRoleName": "WebAP User",
            "grantedRoleDescription": "Web Application Protector User Role"
        }
    ],
	"users": [
        {
            "uiIdentityId": "USER1",
            "firstName": "John",
            "lastName": "Smith",
            "accountId": "ACCOUNT1",
            "email": "example@akamai.com",
            "lastLoginDate": "2016-02-17T18:46:42.000Z"
        },
        {
            "uiIdentityId": "USER2",
            "firstName": "Steve",
            "lastName": "Smith",
            "accountId": "ACCOUNT2",
            "email": "example1@akamai.com",
            "lastLoginDate": "2016-02-17T18:46:42.000Z"
        }
	]
}`,
			expectedPath: "/identity-management/v2/user-admin/roles/123456?actions=true&grantedRoles=true&users=true",
			expectedResponse: &Role{
				RoleID:          123456,
				RoleName:        "Terraform admin updated",
				RoleDescription: "Admin granted role for tests",
				RoleType:        RoleTypeCustom,
				CreatedDate:     "2022-04-11T10:52:03.000Z",
				CreatedBy:       "jBond",
				ModifiedDate:    "2022-04-11T10:59:30.000Z",
				ModifiedBy:      "jBond",
				Actions: &RoleAction{
					Edit:   true,
					Delete: true,
				},
				GrantedRoles: []RoleGrantedRole{
					{
						RoleID:      12345,
						RoleName:    "View Audience Analytics Reports",
						Description: "Publisher Self-Provisioning",
					},
					{
						RoleID:      54321,
						RoleName:    "WebAP User",
						Description: "Web Application Protector User Role",
					},
				},
				Users: []RoleUser{
					{
						UIIdentityID:  "USER1",
						FirstName:     "John",
						LastName:      "Smith",
						AccountID:     "ACCOUNT1",
						Email:         "example@akamai.com",
						LastLoginDate: "2016-02-17T18:46:42.000Z",
					},
					{
						UIIdentityID:  "USER2",
						FirstName:     "Steve",
						LastName:      "Smith",
						AccountID:     "ACCOUNT2",
						Email:         "example1@akamai.com",
						LastLoginDate: "2016-02-17T18:46:42.000Z",
					},
				},
			},
		},
		"200 OK without query params": {
			params:         GetRoleRequest{ID: 123456},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "roleId": 123456,
    "roleName": "Terraform admin updated",
    "roleDescription": "Admin granted role for tests",
    "type": "custom",
    "createdDate": "2022-04-11T10:52:03.000Z",
    "createdBy": "jBond",
    "modifiedDate": "2022-04-11T10:59:30.000Z",
    "modifiedBy": "jBond"
}`,
			expectedPath: "/identity-management/v2/user-admin/roles/123456?actions=false&grantedRoles=false&users=false",
			expectedResponse: &Role{
				RoleID:          123456,
				RoleName:        "Terraform admin updated",
				RoleDescription: "Admin granted role for tests",
				RoleType:        RoleTypeCustom,
				CreatedDate:     "2022-04-11T10:52:03.000Z",
				CreatedBy:       "jBond",
				ModifiedDate:    "2022-04-11T10:59:30.000Z",
				ModifiedBy:      "jBond",
			},
		},
		"404 Not found": {
			params:         GetRoleRequest{ID: 123456},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/identity-management/v2/user-admin/roles/123456?actions=false&grantedRoles=false&users=false",
			responseBody: `
{
    "instance": "",
    "httpStatus": 404,
    "detail": "Role ID not found",
    "title": "Role ID not found",
    "type": "/useradmin-api/error-types/1311"
}`,
			withError: &Error{
				Instance:   "",
				HTTPStatus: http.StatusNotFound,
				Detail:     "Role ID not found",
				Title:      "Role ID not found",
				Type:       "/useradmin-api/error-types/1311",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 Internal server error": {
			params:         GetRoleRequest{ID: 123456},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"detail": "Error making request",
	"status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/roles/123456?actions=false&grantedRoles=false&users=false",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
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
			result, err := client.GetRole(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAM_UpdateRole(t *testing.T) {
	tests := map[string]struct {
		params              UpdateRoleRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *Role
		withError           error
	}{
		"200 OK - update only granted roles an name": {
			params: UpdateRoleRequest{
				ID: 123456,
				RoleRequest: RoleRequest{
					Name: "Terraform admin updated",
					GrantedRoles: []GrantedRoleID{
						{
							ID: 54321,
						},
						{
							ID: 12345,
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "roleId": 123456,
    "roleName": "Terraform admin updated",
    "roleDescription": "Admin granted role for tests",
    "type": "custom",
    "createdDate": "2022-04-11T10:52:03.000Z",
    "createdBy": "jBond",
    "modifiedDate": "2022-04-11T10:59:30.000Z",
    "modifiedBy": "jBond",
    "actions": {
        "edit": true,
        "delete": true
    },
    "grantedRoles": [
        {
            "grantedRoleId": 54321,
            "grantedRoleName": "View Audience Analytics Reports",
            "grantedRoleDescription": "Publisher Self-Provisioning"
        },
        {
            "grantedRoleId": 12345,
            "grantedRoleName": "WebAP User",
            "grantedRoleDescription": "Web Application Protector User Role"
        }
    ]
}`,
			expectedPath:        "/identity-management/v2/user-admin/roles/123456",
			expectedRequestBody: `{"roleName":"Terraform admin updated","grantedRoles":[{"grantedRoleId":54321},{"grantedRoleId":12345}]}`,
			expectedResponse: &Role{
				RoleID:          123456,
				RoleName:        "Terraform admin updated",
				RoleDescription: "Admin granted role for tests",
				RoleType:        RoleTypeCustom,
				CreatedDate:     "2022-04-11T10:52:03.000Z",
				CreatedBy:       "jBond",
				ModifiedDate:    "2022-04-11T10:59:30.000Z",
				ModifiedBy:      "jBond",
				Actions: &RoleAction{
					Edit:   true,
					Delete: true,
				},
				GrantedRoles: []RoleGrantedRole{
					{
						RoleID:      54321,
						RoleName:    "View Audience Analytics Reports",
						Description: "Publisher Self-Provisioning",
					},
					{
						RoleID:      12345,
						RoleName:    "WebAP User",
						Description: "Web Application Protector User Role",
					},
				},
			},
		},
		"500 Internal server error": {
			params: UpdateRoleRequest{
				ID: 123456,
				RoleRequest: RoleRequest{
					Name: "Terraform admin updated",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"detail": "Error making request",
	"status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/roles/123456",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
				StatusCode: http.StatusInternalServerError,
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateRole(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAM_DeleteRole(t *testing.T) {
	tests := map[string]struct {
		params         DeleteRoleRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 Deleted": {
			params:         DeleteRoleRequest{ID: 123456},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v2/user-admin/roles/123456",
		},
		"404 Not found": {
			params:         DeleteRoleRequest{ID: 123456},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/identity-management/v2/user-admin/roles/123456",
			responseBody: `
{
    "instance": "",
    "httpStatus": 404,
    "detail": "",
    "title": "Role not found",
    "type": "/useradmin-api/error-types/1311"
}`,
			withError: &Error{
				Instance:   "",
				HTTPStatus: http.StatusNotFound,
				Detail:     "",
				Title:      "Role not found",
				Type:       "/useradmin-api/error-types/1311",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 Internal server error": {
			params:         DeleteRoleRequest{ID: 123456},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/roles/123456",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
				StatusCode: http.StatusInternalServerError,
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
			err := client.DeleteRole(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_ListRoles(t *testing.T) {
	tests := map[string]struct {
		params           ListRolesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Role
		withError        error
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
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/roles?actions=true&ignoreContext=false&users=false",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
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
			result, err := client.ListRoles(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAM_ListGrantableRoles(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []RoleGrantedRole
		withError        error
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "grantedRoleId": 123456,
        "grantedRoleName": "first role name",
        "grantedRoleDescription": "first role description"
    },
    {
        "grantedRoleId": 654321,
        "grantedRoleName": "second role name",
        "grantedRoleDescription": "second role description"
    }
]`,
			expectedPath: "/identity-management/v2/user-admin/roles/grantable-roles",
			expectedResponse: []RoleGrantedRole{
				{
					RoleID:      123456,
					RoleName:    "first role name",
					Description: "first role description",
				},
				{
					RoleID:      654321,
					RoleName:    "second role name",
					Description: "second role description",
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/roles/grantable-roles",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
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
			result, err := client.ListGrantableRoles(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
