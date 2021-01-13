package iam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi/tools"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestIAM_CreateUser(t *testing.T) {
	tests := map[string]struct {
		params           CreateUserRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *User
		withError        func(*testing.T, error)
	}{
		"201 OK": {
			params: CreateUserRequest{
				User: UserBasicInfo{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@mycompany.com",
					Phone:     "(123) 321-1234",
					Country:   "USA",
					State:     "CA",
				},
				AuthGrants:    []AuthGrant{{GroupID: 1, RoleID: tools.IntPtr(1)}},
				Notifications: UserNotifications{},
			},
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
				User: UserBasicInfo{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@mycompany.com",
					Phone:     "(123) 321-1234",
					Country:   "USA",
					State:     "CA",
				},
				AuthGrants:    []AuthGrant{{GroupID: 1, RoleID: tools.IntPtr(1)}},
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

func TestIAM_UpdateUserInfo(t *testing.T) {
	tests := map[string]struct {
		params           UpdateUserInfoRequest
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
				AuthGrants: []AuthGrant{
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
				AuthGrants: []AuthGrant{
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
