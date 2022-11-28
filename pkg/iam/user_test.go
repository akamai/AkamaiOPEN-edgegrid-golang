package iam

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestIAM_CreateUser(t *testing.T) {
	tests := map[string]struct {
		params           CreateUserRequest
		requestBody      string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *User
		withError        func(*testing.T, error)
	}{
		"201 OK": {
			params: CreateUserRequest{
				UserBasicInfo: UserBasicInfo{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@mycompany.com",
					Phone:     "(123) 321-1234",
					Country:   "USA",
					State:     "CA",
				},
				AuthGrants:    []AuthGrantRequest{{GroupID: 1, RoleID: tools.IntPtr(1)}},
				Notifications: UserNotifications{},
			},
			requestBody:    `{"firstName":"John","lastName":"Doe","email":"john.doe@mycompany.com","phone":"(123) 321-1234","jobTitle":"","tfaEnabled":false,"state":"CA","country":"USA","authGrants":[{"groupId":1,"isBlocked":false,"roleId":1}],"notifications":{"enableEmailNotifications":false,"options":{"newUserNotification":false,"passwordExpiry":false,"proactive":null,"upgrade":null}}}`,
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"uiIdentityId": "A-BC-1234567",
	"firstName": "John",
	"lastName": "Doe",
	"email": "john.doe@mycompany.com",
	"phone": "(123) 321-1234",
	"state": "CA",
	"country": "USA"
}`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities?sendEmail=false",
			expectedResponse: &User{
				IdentityID: "A-BC-1234567",
				UserBasicInfo: UserBasicInfo{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@mycompany.com",
					Phone:     "(123) 321-1234",
					Country:   "USA",
					State:     "CA",
				},
			},
		},
		"500 internal server error": {
			params: CreateUserRequest{
				UserBasicInfo: UserBasicInfo{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@mycompany.com",
					Phone:     "(123) 321-1234",
					Country:   "USA",
					State:     "CA",
				},
				AuthGrants:    []AuthGrantRequest{{GroupID: 1, RoleID: tools.IntPtr(1)}},
				Notifications: UserNotifications{},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities?sendEmail=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if test.requestBody != "" {
					buf := new(bytes.Buffer)
					_, err := buf.ReadFrom(r.Body)
					assert.NoError(t, err)
					req := buf.String()
					assert.Equal(t, test.requestBody, req)
				}
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateUser(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAM_GetUser(t *testing.T) {
	tests := map[string]struct {
		params           GetUserRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *User
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetUserRequest{
				IdentityID: "A-BC-1234567",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"uiIdentityId": "A-BC-1234567",
	"firstName": "John",
	"lastName": "Doe",
	"email": "john.doe@mycompany.com",
	"phone": "(123) 321-1234",
	"state": "CA",
	"country": "USA"
}`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities/A-BC-1234567?actions=false&authGrants=false&notifications=false",
			expectedResponse: &User{
				IdentityID: "A-BC-1234567",
				UserBasicInfo: UserBasicInfo{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@mycompany.com",
					Phone:     "(123) 321-1234",
					Country:   "USA",
					State:     "CA",
				},
			},
		},
		"500 internal server error": {
			params: GetUserRequest{
				IdentityID: "A-BC-1234567",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities/A-BC-1234567?actions=false&authGrants=false&notifications=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
			result, err := client.GetUser(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIam_ListUsers(t *testing.T) {
	tests := map[string]struct {
		params           ListUsersRequest
		responseStatus   int
		expectedPath     string
		responseBody     string
		expectedResponse []UserListItem
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListUsersRequest{
				GroupID:    tools.Int64Ptr(12345),
				AuthGrants: true,
				Actions:    true,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities?actions=true&authGrants=true&groupId=12345",
			responseBody: `[
			  {
				"uiIdentityId": "A-B-123456",
				"firstName": "John",
				"lastName": "Doe",
				"uiUserName": "johndoe",
				"email": "john.doe@mycompany.com",
				"accountId": "1-123A",
				"lastLoginDate": "2016-01-13T17:53:57Z",
				"tfaEnabled": true,
				"tfaConfigured": true,
				"isLocked": false,
				"actions": {
				  "resetPassword": true,
				  "delete": true,
				  "edit": true,
				  "apiClient": true,
				  "thirdPartyAccess": true,
				  "isCloneable": true,
				  "editProfile": true,
				  "canEditTFA": false
				},
				"authGrants": [
				  {
					"groupId": 12345,
					"roleId": 12,
					"groupName": "mygroup",
					"roleName": "admin",
					"roleDescription": "This is a new role that has been created to",
					"isBlocked": false
				  }
				]
			  }
			]`,
			expectedResponse: []UserListItem{
				{
					IdentityID:    "A-B-123456",
					FirstName:     "John",
					LastName:      "Doe",
					UserName:      "johndoe",
					Email:         "john.doe@mycompany.com",
					AccountID:     "1-123A",
					TFAEnabled:    true,
					LastLoginDate: "2016-01-13T17:53:57Z",
					TFAConfigured: true,
					IsLocked:      false,
					Actions: &UserActions{
						APIClient:        true,
						Delete:           true,
						Edit:             true,
						IsCloneable:      true,
						ResetPassword:    true,
						ThirdPartyAccess: true,
						EditProfile:      true,
					},
					AuthGrants: []AuthGrant{
						{
							GroupID:         12345,
							RoleID:          tools.IntPtr(12),
							GroupName:       "mygroup",
							RoleName:        "admin",
							RoleDescription: "This is a new role that has been created to",
						},
					},
				},
			},
		},
		"200 OK, no actions nor grants": {
			params: ListUsersRequest{
				GroupID: tools.Int64Ptr(12345),
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities?actions=false&authGrants=false&groupId=12345",
			responseBody: `[
			  {
				"uiIdentityId": "A-B-123456",
				"firstName": "John",
				"lastName": "Doe",
				"uiUserName": "johndoe",
				"email": "john.doe@mycompany.com",
				"accountId": "1-123A",
				"lastLoginDate": "2016-01-13T17:53:57Z",
				"tfaEnabled": true,
				"tfaConfigured": true,
				"isLocked": false
			  }
			]`,
			expectedResponse: []UserListItem{
				{
					IdentityID:    "A-B-123456",
					FirstName:     "John",
					LastName:      "Doe",
					UserName:      "johndoe",
					Email:         "john.doe@mycompany.com",
					AccountID:     "1-123A",
					TFAEnabled:    true,
					LastLoginDate: "2016-01-13T17:53:57Z",
					TFAConfigured: true,
					IsLocked:      false,
				},
			},
		},
		"no group id": {
			params:       ListUsersRequest{},
			expectedPath: "/identity-management/v2/user-admin/ui-identities?actions=false&authGrants=false",
			responseBody: `[
			  {
				"uiIdentityId": "A-B-123456",
				"firstName": "John",
				"lastName": "Doe",
				"uiUserName": "johndoe",
				"email": "john.doe@mycompany.com",
				"accountId": "1-123A",
				"lastLoginDate": "2016-01-13T17:53:57Z",
				"tfaEnabled": true,
				"tfaConfigured": true,
				"isLocked": false
			  }
			]`,
			expectedResponse: []UserListItem{
				{
					IdentityID:    "A-B-123456",
					FirstName:     "John",
					LastName:      "Doe",
					UserName:      "johndoe",
					Email:         "john.doe@mycompany.com",
					AccountID:     "1-123A",
					TFAEnabled:    true,
					LastLoginDate: "2016-01-13T17:53:57Z",
					TFAConfigured: true,
					IsLocked:      false,
				},
			},
			responseStatus: http.StatusOK,
		},
		"500 internal server error": {
			params: ListUsersRequest{
				GroupID:    tools.Int64Ptr(12345),
				AuthGrants: true,
				Actions:    true,
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/identity-management/v2/user-admin/ui-identities?actions=true&authGrants=true&groupId=12345",
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error processing request",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error processing request",
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
			users, err := client.ListUsers(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}

func TestIAM_UpdateUserInfo(t *testing.T) {
	tests := map[string]struct {
		params           UpdateUserInfoRequest
		requestBody      string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UserBasicInfo
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateUserInfoRequest{
				IdentityID: "1-ABCDE",
				User: UserBasicInfo{
					FirstName:         "John",
					LastName:          "Doe",
					Email:             "john.doe@mycompany.com",
					Phone:             "(123) 321-1234",
					Country:           "USA",
					State:             "CA",
					PreferredLanguage: "English",
					ContactType:       "Billing",
					SessionTimeOut:    tools.IntPtr(30),
					TimeZone:          "GMT",
				},
			},
			requestBody:    `{"firstName":"John","lastName":"Doe","email":"john.doe@mycompany.com","phone":"(123) 321-1234","timeZone":"GMT","jobTitle":"","tfaEnabled":false,"state":"CA","country":"USA","contactType":"Billing","preferredLanguage":"English","sessionTimeOut":30}`,
			responseStatus: http.StatusOK,
			responseBody: `
{
	"firstName": "John",
	"lastName": "Doe",
	"email": "john.doe@mycompany.com",
	"phone": "(123) 321-1234",
	"state": "CA",
	"country": "USA",
	"preferredLanguage": "English",
	"contactType": "Billing",
	"sessionTimeOut": 30,
	"timeZone": "GMT"
}`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities/1-ABCDE/basic-info",
			expectedResponse: &UserBasicInfo{
				FirstName:         "John",
				LastName:          "Doe",
				Email:             "john.doe@mycompany.com",
				Phone:             "(123) 321-1234",
				Country:           "USA",
				State:             "CA",
				PreferredLanguage: "English",
				ContactType:       "Billing",
				SessionTimeOut:    tools.IntPtr(30),
				TimeZone:          "GMT",
			},
		},
		"500 internal server error": {
			params: UpdateUserInfoRequest{
				IdentityID: "1-ABCDE",
				User: UserBasicInfo{
					FirstName:         "John",
					LastName:          "Doe",
					Email:             "john.doe@mycompany.com",
					Phone:             "(123) 321-1234",
					Country:           "USA",
					State:             "CA",
					PreferredLanguage: "English",
					ContactType:       "Billing",
					SessionTimeOut:    tools.IntPtr(30),
					TimeZone:          "GMT",
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
			expectedPath: "/identity-management/v2/user-admin/ui-identities/1-ABCDE/basic-info",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
				assert.Equal(t, http.MethodPut, r.Method)
				if test.requestBody != "" {
					buf := new(bytes.Buffer)
					_, err := buf.ReadFrom(r.Body)
					assert.NoError(t, err)
					req := buf.String()
					assert.Equal(t, test.requestBody, req)
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateUserInfo(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAM_UpdateUserNotifications(t *testing.T) {
	tests := map[string]struct {
		params           UpdateUserNotificationsRequest
		requestBody      string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UserNotifications
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateUserNotificationsRequest{
				IdentityID: "1-ABCDE",
				Notifications: UserNotifications{
					EnableEmail: true,
					Options: UserNotificationOptions{
						Upgrade:        []string{"NetStorage", "Other Upgrade Notifications (Planned)"},
						Proactive:      []string{"EdgeScape", "EdgeSuite (HTTP Content Delivery)"},
						PasswordExpiry: true,
						NewUser:        true,
					},
				},
			},
			requestBody:    `{"enableEmailNotifications":true,"options":{"newUserNotification":true,"passwordExpiry":true,"proactive":["EdgeScape","EdgeSuite (HTTP Content Delivery)"],"upgrade":["NetStorage","Other Upgrade Notifications (Planned)"]}}`,
			responseStatus: http.StatusOK,
			responseBody: `
{
	"enableEmailNotifications": true,
	"options": {
		"upgrade": [
			"NetStorage",
			"Other Upgrade Notifications (Planned)"
		],
		"proactive": [
			"EdgeScape",
			"EdgeSuite (HTTP Content Delivery)"
		],
		"passwordExpiry": true,
		"newUserNotification": true
	}
}`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities/1-ABCDE/notifications",
			expectedResponse: &UserNotifications{
				EnableEmail: true,
				Options: UserNotificationOptions{
					Upgrade:        []string{"NetStorage", "Other Upgrade Notifications (Planned)"},
					Proactive:      []string{"EdgeScape", "EdgeSuite (HTTP Content Delivery)"},
					PasswordExpiry: true,
					NewUser:        true,
				},
			},
		},
		"500 internal server error": {
			params: UpdateUserNotificationsRequest{
				IdentityID: "1-ABCDE",
				Notifications: UserNotifications{
					EnableEmail: true,
					Options: UserNotificationOptions{
						Upgrade:        []string{"NetStorage", "Other Upgrade Notifications (Planned)"},
						Proactive:      []string{"EdgeScape", "EdgeSuite (HTTP Content Delivery)"},
						PasswordExpiry: true,
						NewUser:        true,
					},
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
			expectedPath: "/identity-management/v2/user-admin/ui-identities/1-ABCDE/notifications",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
				assert.Equal(t, http.MethodPut, r.Method)
				if test.requestBody != "" {
					buf := new(bytes.Buffer)
					_, err := buf.ReadFrom(r.Body)
					assert.NoError(t, err)
					req := buf.String()
					assert.Equal(t, test.requestBody, req)
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateUserNotifications(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAM_UpdateUserAuthGrants(t *testing.T) {
	tests := map[string]struct {
		params           UpdateUserAuthGrantsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []AuthGrant
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateUserAuthGrantsRequest{
				IdentityID: "1-ABCDE",
				AuthGrants: []AuthGrantRequest{
					{
						GroupID: 12345,
						RoleID:  tools.IntPtr(16),
						Subgroups: []AuthGrantRequest{
							{
								GroupID: 54321,
							},
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
	{
		"groupId": 12345,
		"roleId": 16,
		"subGroups": [
			{
				"groupId": 54321
			}
		]
	}
]`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities/1-ABCDE/auth-grants",
			expectedResponse: []AuthGrant{
				{
					GroupID: 12345,
					RoleID:  tools.IntPtr(16),
					Subgroups: []AuthGrant{
						{
							GroupID: 54321,
						},
					},
				},
			},
		},
		"500 internal server error": {
			params: UpdateUserAuthGrantsRequest{
				IdentityID: "1-ABCDE",
				AuthGrants: []AuthGrantRequest{
					{
						GroupID: 12345,
						RoleID:  tools.IntPtr(16),
						Subgroups: []AuthGrantRequest{
							{
								GroupID: 54321,
							},
						},
					},
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
			expectedPath: "/identity-management/v2/user-admin/ui-identities/1-ABCDE/auth-grants",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateUserAuthGrants(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAM_RemoveUser(t *testing.T) {
	tests := map[string]struct {
		params         RemoveUserRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			params: RemoveUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusOK,
			responseBody:   "",
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE",
		},
		"204 No Content": {
			params: RemoveUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE",
		},
		"500 internal server error": {
			params: RemoveUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities/1-ABCDE",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.RemoveUser(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_UpdateTFA(t *testing.T) {
	tests := map[string]struct {
		params         UpdateTFARequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No Content": {
			params: UpdateTFARequest{
				IdentityID: "1-ABCDE",
				Action:     TFAActionEnable,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/identity-management/v2/user-admin/ui-identities/1-ABCDE/tfa?action=enable",
		},
		"500 internal server error": {
			params: UpdateTFARequest{
				IdentityID: "1-ABCDE",
				Action:     TFAActionDisable,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v2/user-admin/ui-identities/1-ABCDE/tfa?action=disable",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.UpdateTFA(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
