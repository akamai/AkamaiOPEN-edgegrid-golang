package iam

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/ptr"
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
					FirstName:                "John",
					LastName:                 "Doe",
					Email:                    "john.doe@mycompany.com",
					Phone:                    "(123) 321-1234",
					Country:                  "USA",
					State:                    "CA",
					AdditionalAuthentication: NoneAuthentication,
				},
				AuthGrants: []AuthGrantRequest{{GroupID: 1, RoleID: ptr.To(1)}},
			},
			requestBody:    `{"firstName":"John","lastName":"Doe","email":"john.doe@mycompany.com","phone":"(123) 321-1234","jobTitle":"","tfaEnabled":false,"state":"CA","country":"USA","additionalAuthentication":"NONE","authGrants":[{"groupId":1,"isBlocked":false,"roleId":1}]}`,
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"uiIdentityId": "A-BC-1234567",
	"firstName": "John",
	"lastName": "Doe",
	"email": "john.doe@mycompany.com",
	"phone": "(123) 321-1234",
	"state": "CA",
	"country": "USA",
	"additionalAuthenticationConfigured": false,
	"additionalAuthentication": "NONE"
}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities?sendEmail=false",
			expectedResponse: &User{
				IdentityID: "A-BC-1234567",
				UserBasicInfo: UserBasicInfo{
					FirstName:                "John",
					LastName:                 "Doe",
					Email:                    "john.doe@mycompany.com",
					Phone:                    "(123) 321-1234",
					Country:                  "USA",
					State:                    "CA",
					AdditionalAuthentication: NoneAuthentication,
				},
				AdditionalAuthenticationConfigured: false,
			},
		},
		"201 OK - all fields": {
			params: CreateUserRequest{
				UserBasicInfo: UserBasicInfo{
					FirstName:                "John",
					LastName:                 "Doe",
					UserName:                 "UserName",
					Email:                    "john.doe@mycompany.com",
					Phone:                    "(123) 321-1234",
					TimeZone:                 "GMT+2",
					JobTitle:                 "Title",
					SecondaryEmail:           "second@email.com",
					MobilePhone:              "123123123",
					Address:                  "Address",
					City:                     "City",
					State:                    "CA",
					ZipCode:                  "11-111",
					Country:                  "USA",
					ContactType:              "Dev",
					PreferredLanguage:        "EN",
					SessionTimeOut:           ptr.To(1),
					AdditionalAuthentication: MFAAuthentication,
				},
				AuthGrants: []AuthGrantRequest{{GroupID: 1, RoleID: ptr.To(1)}},
				Notifications: &UserNotifications{
					EnableEmail: false,
					Options: UserNotificationOptions{
						NewUser:                   false,
						PasswordExpiry:            false,
						Proactive:                 []string{"Test1"},
						Upgrade:                   []string{"Test2"},
						APIClientCredentialExpiry: false,
					},
				},
				SendEmail: true,
			},
			requestBody:    `{"firstName":"John","lastName":"Doe","uiUserName":"UserName","email":"john.doe@mycompany.com","phone":"(123) 321-1234","timeZone":"GMT+2","jobTitle":"Title","tfaEnabled":false,"secondaryEmail":"second@email.com","mobilePhone":"123123123","address":"Address","city":"City","state":"CA","zipCode":"11-111","country":"USA","contactType":"Dev","preferredLanguage":"EN","sessionTimeOut":1,"additionalAuthentication":"MFA","authGrants":[{"groupId":1,"isBlocked":false,"roleId":1}],"notifications":{"enableEmailNotifications":false,"options":{"newUserNotification":false,"passwordExpiry":false,"proactive":["Test1"],"upgrade":["Test2"],"apiClientCredentialExpiryNotification":false}}}`,
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"uiIdentityId": "A-BC-1234567",
	"firstName": "John",
	"lastName": "Doe",
	"email": "john.doe@mycompany.com",
	"phone": "(123) 321-1234",
	"state": "CA",
	"country": "USA",
	"additionalAuthenticationConfigured": false,
	"additionalAuthentication": "NONE"
}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities?sendEmail=true",
			expectedResponse: &User{
				IdentityID: "A-BC-1234567",
				UserBasicInfo: UserBasicInfo{
					FirstName:                "John",
					LastName:                 "Doe",
					Email:                    "john.doe@mycompany.com",
					Phone:                    "(123) 321-1234",
					Country:                  "USA",
					State:                    "CA",
					AdditionalAuthentication: NoneAuthentication,
				},
				AdditionalAuthenticationConfigured: false,
			},
		},
		"validation errors": {
			params: CreateUserRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create user: struct validation:\nAdditionalAuthentication: cannot be blank; AuthGrants: cannot be blank; Country: cannot be blank; Email: cannot be blank; FirstName: cannot be blank; LastName: cannot be blank.", err.Error())
			},
		},
		"500 internal server error": {
			params: CreateUserRequest{
				UserBasicInfo: UserBasicInfo{
					FirstName:                "John",
					LastName:                 "Doe",
					Email:                    "john.doe@mycompany.com",
					Phone:                    "(123) 321-1234",
					Country:                  "USA",
					State:                    "CA",
					AdditionalAuthentication: TFAAuthentication,
				},
				AuthGrants: []AuthGrantRequest{{GroupID: 1, RoleID: ptr.To(1)}},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities?sendEmail=false",
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
	"country": "USA",
	"accountId": "sampleID",
	"userStatus": "PENDING"
}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/A-BC-1234567?actions=false&authGrants=false&notifications=false",
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
				UserStatus: "PENDING",
				AccountID:  "sampleID",
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/A-BC-1234567?actions=false&authGrants=false&notifications=false",
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
				GroupID:    ptr.To(int64(12345)),
				AuthGrants: true,
				Actions:    true,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/user-admin/ui-identities?actions=true&authGrants=true&groupId=12345",
			responseBody: `[
			  {
				"uiIdentityId": "A-B-123456",
				"firstName": "John",
				"lastName": "Doe",
				"uiUserName": "johndoe",
				"email": "john.doe@mycompany.com",
				"accountId": "1-123A",
				"lastLoginDate": "2016-01-13T17:53:57.000Z",
				"tfaEnabled": true,
				"tfaConfigured": true,
				"isLocked": false,
				"additionalAuthentication": "TFA",
				"additionalAuthenticationConfigured": false,
				"actions": {
				  "resetPassword": true,
				  "delete": true,
				  "edit": true,
				  "apiClient": true,
				  "thirdPartyAccess": true,
				  "isCloneable": true,
				  "editProfile": true,
				  "canEditTFA": true,
				  "canEditMFA": true,
				  "canEditNone": true
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
					IdentityID:                         "A-B-123456",
					FirstName:                          "John",
					LastName:                           "Doe",
					UserName:                           "johndoe",
					Email:                              "john.doe@mycompany.com",
					AccountID:                          "1-123A",
					TFAEnabled:                         true,
					LastLoginDate:                      test.NewTimeFromString(t, "2016-01-13T17:53:57.000Z"),
					TFAConfigured:                      true,
					IsLocked:                           false,
					AdditionalAuthentication:           TFAAuthentication,
					AdditionalAuthenticationConfigured: false,
					Actions: &UserActions{
						APIClient:        true,
						Delete:           true,
						Edit:             true,
						IsCloneable:      true,
						ResetPassword:    true,
						ThirdPartyAccess: true,
						EditProfile:      true,
						CanEditMFA:       true,
						CanEditNone:      true,
						CanEditTFA:       true,
					},
					AuthGrants: []AuthGrant{
						{
							GroupID:         12345,
							RoleID:          ptr.To(12),
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
				GroupID: ptr.To(int64(12345)),
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/user-admin/ui-identities?actions=false&authGrants=false&groupId=12345",
			responseBody: `[
			  {
				"uiIdentityId": "A-B-123456",
				"firstName": "John",
				"lastName": "Doe",
				"uiUserName": "johndoe",
				"email": "john.doe@mycompany.com",
				"accountId": "1-123A",
				"lastLoginDate": "2016-01-13T17:53:57.000Z",
				"tfaEnabled": true,
				"tfaConfigured": true,
				"isLocked": false,
				"additionalAuthentication": "MFA",
				"additionalAuthenticationConfigured": true
			  }
			]`,
			expectedResponse: []UserListItem{
				{
					IdentityID:                         "A-B-123456",
					FirstName:                          "John",
					LastName:                           "Doe",
					UserName:                           "johndoe",
					Email:                              "john.doe@mycompany.com",
					AccountID:                          "1-123A",
					TFAEnabled:                         true,
					LastLoginDate:                      test.NewTimeFromString(t, "2016-01-13T17:53:57.000Z"),
					TFAConfigured:                      true,
					IsLocked:                           false,
					AdditionalAuthenticationConfigured: true,
					AdditionalAuthentication:           MFAAuthentication,
				},
			},
		},
		"200 OK, no group id": {
			params:       ListUsersRequest{},
			expectedPath: "/identity-management/v3/user-admin/ui-identities?actions=false&authGrants=false",
			responseBody: `[
			  {
				"uiIdentityId": "A-B-123456",
				"firstName": "John",
				"lastName": "Doe",
				"uiUserName": "johndoe",
				"email": "john.doe@mycompany.com",
				"accountId": "1-123A",
				"lastLoginDate": "2016-01-13T17:53:57.000Z",
				"tfaEnabled": true,
				"tfaConfigured": true,
				"isLocked": false,
				"additionalAuthentication": "TFA",
				"additionalAuthenticationConfigured": true
			  }
			]`,
			expectedResponse: []UserListItem{
				{
					IdentityID:                         "A-B-123456",
					FirstName:                          "John",
					LastName:                           "Doe",
					UserName:                           "johndoe",
					Email:                              "john.doe@mycompany.com",
					AccountID:                          "1-123A",
					TFAEnabled:                         true,
					LastLoginDate:                      test.NewTimeFromString(t, "2016-01-13T17:53:57.000Z"),
					TFAConfigured:                      true,
					IsLocked:                           false,
					AdditionalAuthentication:           TFAAuthentication,
					AdditionalAuthenticationConfigured: true,
				},
			},
			responseStatus: http.StatusOK,
		},
		"500 internal server error": {
			params: ListUsersRequest{
				GroupID:    ptr.To(int64(12345)),
				AuthGrants: true,
				Actions:    true,
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/identity-management/v3/user-admin/ui-identities?actions=true&authGrants=true&groupId=12345",
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
					FirstName:                "John",
					LastName:                 "Doe",
					Email:                    "john.doe@mycompany.com",
					Phone:                    "(123) 321-1234",
					Country:                  "USA",
					State:                    "CA",
					PreferredLanguage:        "English",
					ContactType:              "Billing",
					SessionTimeOut:           ptr.To(30),
					TimeZone:                 "GMT",
					AdditionalAuthentication: NoneAuthentication,
				},
			},
			requestBody:    `{"firstName":"John","lastName":"Doe","email":"john.doe@mycompany.com","phone":"(123) 321-1234","timeZone":"GMT","jobTitle":"","tfaEnabled":false,"state":"CA","country":"USA","contactType":"Billing","preferredLanguage":"English","sessionTimeOut":30,"additionalAuthentication":"NONE"}`,
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
	"timeZone": "GMT",
	"additionalAuthentication": "NONE"
}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/basic-info",
			expectedResponse: &UserBasicInfo{
				FirstName:                "John",
				LastName:                 "Doe",
				Email:                    "john.doe@mycompany.com",
				Phone:                    "(123) 321-1234",
				Country:                  "USA",
				State:                    "CA",
				PreferredLanguage:        "English",
				ContactType:              "Billing",
				SessionTimeOut:           ptr.To(30),
				TimeZone:                 "GMT",
				AdditionalAuthentication: NoneAuthentication,
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
					SessionTimeOut:    ptr.To(30),
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/basic-info",
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
				Notifications: &UserNotifications{
					EnableEmail: true,
					Options: UserNotificationOptions{
						Upgrade:                   []string{"NetStorage", "Other Upgrade Notifications (Planned)"},
						Proactive:                 []string{"EdgeScape", "EdgeSuite (HTTP Content Delivery)"},
						PasswordExpiry:            true,
						NewUser:                   true,
						APIClientCredentialExpiry: true,
					},
				},
			},
			requestBody:    `{"enableEmailNotifications":true,"options":{"newUserNotification":true,"passwordExpiry":true,"proactive":["EdgeScape","EdgeSuite (HTTP Content Delivery)"],"upgrade":["NetStorage","Other Upgrade Notifications (Planned)"],"apiClientCredentialExpiryNotification":true}}`,
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
		"newUserNotification": true,
		"apiClientCredentialExpiryNotification": true
	}
}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/notifications",
			expectedResponse: &UserNotifications{
				EnableEmail: true,
				Options: UserNotificationOptions{
					Upgrade:                   []string{"NetStorage", "Other Upgrade Notifications (Planned)"},
					Proactive:                 []string{"EdgeScape", "EdgeSuite (HTTP Content Delivery)"},
					PasswordExpiry:            true,
					NewUser:                   true,
					APIClientCredentialExpiry: true,
				},
			},
		},
		"validation errors": {
			params: UpdateUserNotificationsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update user notifications: struct validation:\nIdentityID: cannot be blank\nNotifications: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: UpdateUserNotificationsRequest{
				IdentityID: "1-ABCDE",
				Notifications: &UserNotifications{
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/notifications",
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
						RoleID:  ptr.To(16),
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/auth-grants",
			expectedResponse: []AuthGrant{
				{
					GroupID: 12345,
					RoleID:  ptr.To(16),
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
						RoleID:  ptr.To(16),
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/auth-grants",
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
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE",
		},
		"204 No Content": {
			params: RemoveUserRequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE",
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE",
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

func TestIAM_UpdateMFA(t *testing.T) {
	tests := map[string]struct {
		params         UpdateMFARequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No Content": {
			params: UpdateMFARequest{
				IdentityID: "1-ABCDE",
				Value:      MFAAuthentication,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE/additionalAuthentication",
		},
		"500 internal server error": {
			params: UpdateMFARequest{
				IdentityID: "1-ABCDE",
				Value:      MFAAuthentication,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/additionalAuthentication",
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
			err := client.UpdateMFA(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_ResetMFA(t *testing.T) {
	tests := map[string]struct {
		params         ResetMFARequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No Content": {
			params: ResetMFARequest{
				IdentityID: "1-ABCDE",
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/identity-management/v3/user-admin/ui-identities/1-ABCDE/additionalAuthentication/reset",
		},
		"500 internal server error": {
			params: ResetMFARequest{
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
			expectedPath: "/identity-management/v3/user-admin/ui-identities/1-ABCDE/additionalAuthentication/reset",
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
			err := client.ResetMFA(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
