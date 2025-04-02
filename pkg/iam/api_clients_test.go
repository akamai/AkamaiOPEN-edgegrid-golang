package iam

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAM_LockAPIClient(t *testing.T) {
	tests := map[string]struct {
		params           LockAPIClientRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *LockAPIClientResponse
		withError        func(*testing.T, error)
	}{
		"200 OK with specified client": {
			params: LockAPIClientRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/lock",
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessToken": "test_token1234",
  "activeCredentialCount": 1,
  "allowAccountSwitch": false,
  "authorizedUsers": [
    "jdoe"
  ],
  "clientDescription": "Test",
  "clientId": "abcd1234",
  "clientName": "test",
  "clientType": "CLIENT",
  "createdBy": "jdoe",
  "createdDate": "2022-05-13T20:04:35.000Z",
  "isLocked": true,
  "notificationEmails": [
    "jdoe@example.com"
  ],
  "serviceConsumerToken": "test_token12345"
}`,
			expectedResponse: &LockAPIClientResponse{
				AccessToken:             "test_token1234",
				ActiveCredentialCount:   1,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"jdoe"},
				CanAutoCreateCredential: false,
				ClientDescription:       "Test",
				ClientID:                "abcd1234",
				ClientName:              "test",
				ClientType:              "CLIENT",
				CreatedBy:               "jdoe",
				CreatedDate:             test.NewTimeFromString(t, "2022-05-13T20:04:35.000Z"),
				IsLocked:                true,
				NotificationEmails:      []string{"jdoe@example.com"},
				ServiceConsumerToken:    "test_token12345",
			},
		},
		"200 OK - self": {
			params:         LockAPIClientRequest{},
			expectedPath:   "/identity-management/v3/api-clients/self/lock",
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessToken": "test_token1234",
  "activeCredentialCount": 1,
  "allowAccountSwitch": false,
  "authorizedUsers": [
    "jdoe"
  ],
  "clientDescription": "Test",
  "clientId": "abcd1234",
  "clientName": "test",
  "clientType": "CLIENT",
  "createdBy": "jdoe",
  "createdDate": "2022-05-13T20:04:35.000Z",
  "isLocked": true,
  "notificationEmails": [
    "jdoe@example.com"
  ],
  "serviceConsumerToken": "test_token12345"
}`,
			expectedResponse: &LockAPIClientResponse{
				AccessToken:             "test_token1234",
				ActiveCredentialCount:   1,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"jdoe"},
				CanAutoCreateCredential: false,
				ClientDescription:       "Test",
				ClientID:                "abcd1234",
				ClientName:              "test",
				ClientType:              "CLIENT",
				CreatedBy:               "jdoe",
				CreatedDate:             test.NewTimeFromString(t, "2022-05-13T20:04:35.000Z"),
				IsLocked:                true,
				NotificationEmails:      []string{"jdoe@example.com"},
				ServiceConsumerToken:    "test_token12345",
			},
		},
		"404 Not Found": {
			params: LockAPIClientRequest{
				ClientID: "test12344",
			},
			expectedPath:   "/identity-management/v3/api-clients/test12344/lock",
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"instance": "",
				"httpStatus": 404,
				"detail": "",
				"title": "invalid open identity",
				"type": "/identity-management/error-types/2"
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusNotFound,
					Detail:     "",
					Title:      "invalid open identity",
					Type:       "/identity-management/error-types/2",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: LockAPIClientRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/lock",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
			}`,
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

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			response, err := client.LockAPIClient(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestIAM_UnlockAPIClient(t *testing.T) {
	tests := map[string]struct {
		params           UnlockAPIClientRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *UnlockAPIClientResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UnlockAPIClientRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/unlock",
			responseStatus: http.StatusOK,
			responseBody: `
{
  "accessToken": "test_token1234",
  "activeCredentialCount": 1,
  "allowAccountSwitch": false,
  "authorizedUsers": [
    "jdoe"
  ],
  "clientDescription": "Test",
  "clientId": "abcd1234",
  "clientName": "test",
  "clientType": "CLIENT",
  "createdBy": "jdoe",
  "createdDate": "2022-05-13T20:04:35.000Z",
  "isLocked": true,
  "notificationEmails": [
    "jdoe@example.com"
  ],
  "serviceConsumerToken": "test_token12345"
}`,
			expectedResponse: &UnlockAPIClientResponse{
				AccessToken:             "test_token1234",
				ActiveCredentialCount:   1,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"jdoe"},
				CanAutoCreateCredential: false,
				ClientDescription:       "Test",
				ClientID:                "abcd1234",
				ClientName:              "test",
				ClientType:              "CLIENT",
				CreatedBy:               "jdoe",
				CreatedDate:             test.NewTimeFromString(t, "2022-05-13T20:04:35.000Z"),
				IsLocked:                true,
				NotificationEmails:      []string{"jdoe@example.com"},
				ServiceConsumerToken:    "test_token12345",
			},
		},
		"validation errors": {
			params: UnlockAPIClientRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "unlock api client: struct validation:\nClientID: cannot be blank", err.Error())
			},
		},
		"404 Not Found": {
			params: UnlockAPIClientRequest{
				ClientID: "test12344",
			},
			expectedPath:   "/identity-management/v3/api-clients/test12344/unlock",
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"instance": "",
				"httpStatus": 404,
				"detail": "",
				"title": "invalid open identity",
				"type": "/identity-management/error-types/2"
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Instance:   "",
					HTTPStatus: http.StatusNotFound,
					Detail:     "",
					Title:      "invalid open identity",
					Type:       "/identity-management/error-types/2",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: UnlockAPIClientRequest{
				ClientID: "test1234",
			},
			expectedPath:   "/identity-management/v3/api-clients/test1234/unlock",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
			}`,
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

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			response, err := client.UnlockAPIClient(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestIAM_ListAPIClients(t *testing.T) {
	tests := map[string]struct {
		params           ListAPIClientsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse ListAPIClientsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         ListAPIClientsRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "clientId": "abcdefgh12345678",
        "clientName": "test_user_1",
        "clientDescription": "test_user_1 description",
        "clientType": "CLIENT",
        "authorizedUsers": [
            "user1"
        ],
        "canAutoCreateCredential": false,
        "notificationEmails": [
            "user1@example.com"
        ],
        "activeCredentialCount": 0,
        "allowAccountSwitch": false,
        "createdDate": "2024-07-16T23:01:50.000Z",
        "createdBy": "admin",
        "isLocked": false,
        "accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
        "serviceConsumerToken": "akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1"
    },
    {
        "clientId": "hgfedcba87654321",
        "clientName": "test_user_2",
        "clientDescription": "test_user_2 description",
        "clientType": "SERVICE_ACCOUNT",
        "authorizedUsers": [
            "user2"
        ],
        "canAutoCreateCredential": true,
        "notificationEmails": [
            "user2@example.com"
        ],
        "activeCredentialCount": 1,
        "allowAccountSwitch": false,
        "createdDate": "2023-07-03T15:04:01.000Z",
        "createdBy": "admin",
        "isLocked": false,
        "accessToken": "akaa-8h7g6f5e8h7g6f5e-8h7g6f5e8h7g6f5e",
        "serviceConsumerToken": "akaa-e5f6g7h8e5f6g7h8-e5f6g7h8e5f6g7h8"
    }
]`,
			expectedPath: "/identity-management/v3/api-clients?actions=false",
			expectedResponse: ListAPIClientsResponse{
				{
					AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
					ActiveCredentialCount:   0,
					AllowAccountSwitch:      false,
					AuthorizedUsers:         []string{"user1"},
					CanAutoCreateCredential: false,
					ClientDescription:       "test_user_1 description",
					ClientID:                "abcdefgh12345678",
					ClientName:              "test_user_1",
					ClientType:              ClientClientType,
					CreatedBy:               "admin",
					CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
					IsLocked:                false,
					NotificationEmails:      []string{"user1@example.com"},
					ServiceConsumerToken:    "akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1",
				},
				{
					AccessToken:             "akaa-8h7g6f5e8h7g6f5e-8h7g6f5e8h7g6f5e",
					ActiveCredentialCount:   1,
					AllowAccountSwitch:      false,
					AuthorizedUsers:         []string{"user2"},
					CanAutoCreateCredential: true,
					ClientDescription:       "test_user_2 description",
					ClientID:                "hgfedcba87654321",
					ClientName:              "test_user_2",
					ClientType:              ServiceAccountClientType,
					CreatedBy:               "admin",
					CreatedDate:             test.NewTimeFromString(t, "2023-07-03T15:04:01.000Z"),
					IsLocked:                false,
					NotificationEmails:      []string{"user2@example.com"},
					ServiceConsumerToken:    "akaa-e5f6g7h8e5f6g7h8-e5f6g7h8e5f6g7h8",
				},
			},
		},
		"200 with actions": {
			params:         ListAPIClientsRequest{Actions: true},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "clientId": "abcdefgh12345678",
        "clientName": "test_user_1",
        "clientDescription": "test_user_1 description",
        "clientType": "CLIENT",
        "authorizedUsers": [
            "user1"
        ],
        "canAutoCreateCredential": false,
        "notificationEmails": [
            "user1@example.com"
        ],
        "activeCredentialCount": 0,
        "allowAccountSwitch": false,
        "createdDate": "2024-07-16T23:01:50.000Z",
        "createdBy": "admin",
        "isLocked": false,
        "accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
        "serviceConsumerToken": "akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1",
        "actions": {
            "lock": false,
            "unlock": false,
            "edit": false,
            "transfer": false,
            "delete": false,
            "deactivateAll": false
        }
    },
    {
        "clientId": "hgfedcba87654321",
        "clientName": "test_user_2",
        "clientDescription": "test_user_2 description",
        "clientType": "SERVICE_ACCOUNT",
        "authorizedUsers": [
            "user2"
        ],
        "canAutoCreateCredential": true,
        "notificationEmails": [
            "user2@example.com"
        ],
        "activeCredentialCount": 1,
        "allowAccountSwitch": false,
        "createdDate": "2023-07-03T15:04:01.000Z",
        "createdBy": "admin",
        "isLocked": false,
        "accessToken": "akaa-8h7g6f5e8h7g6f5e-8h7g6f5e8h7g6f5e",
        "serviceConsumerToken": "akaa-e5f6g7h8e5f6g7h8-e5f6g7h8e5f6g7h8",
        "actions": {
            "lock": true,
            "unlock": true,
            "edit": true,
            "transfer": true,
            "delete": true,
            "deactivateAll": true
        }
    }
]`,
			expectedPath: "/identity-management/v3/api-clients?actions=true",
			expectedResponse: ListAPIClientsResponse{
				{
					AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
					ActiveCredentialCount:   0,
					AllowAccountSwitch:      false,
					AuthorizedUsers:         []string{"user1"},
					CanAutoCreateCredential: false,
					ClientDescription:       "test_user_1 description",
					ClientID:                "abcdefgh12345678",
					ClientName:              "test_user_1",
					ClientType:              ClientClientType,
					CreatedBy:               "admin",
					CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
					IsLocked:                false,
					NotificationEmails:      []string{"user1@example.com"},
					ServiceConsumerToken:    "akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1",
					Actions: &ListAPIClientsActions{
						Delete:        false,
						DeactivateAll: false,
						Edit:          false,
						Lock:          false,
						Transfer:      false,
						Unlock:        false,
					},
				},
				{
					AccessToken:             "akaa-8h7g6f5e8h7g6f5e-8h7g6f5e8h7g6f5e",
					ActiveCredentialCount:   1,
					AllowAccountSwitch:      false,
					AuthorizedUsers:         []string{"user2"},
					CanAutoCreateCredential: true,
					ClientDescription:       "test_user_2 description",
					ClientID:                "hgfedcba87654321",
					ClientName:              "test_user_2",
					ClientType:              ServiceAccountClientType,
					CreatedBy:               "admin",
					CreatedDate:             test.NewTimeFromString(t, "2023-07-03T15:04:01.000Z"),
					IsLocked:                false,
					NotificationEmails:      []string{"user2@example.com"},
					ServiceConsumerToken:    "akaa-e5f6g7h8e5f6g7h8-e5f6g7h8e5f6g7h8",
					Actions: &ListAPIClientsActions{
						Delete:        true,
						DeactivateAll: true,
						Edit:          true,
						Lock:          true,
						Transfer:      true,
						Unlock:        true,
					},
				},
			},
		},
		"500 internal server error": {
			params:         ListAPIClientsRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/identity-management/v3/api-clients?actions=false",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListAPIClients(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_CreateAPIClient(t *testing.T) {
	tests := map[string]struct {
		params              CreateAPIClientRequest
		expectedPath        string
		responseStatus      int
		responseBody        string
		expectedResponse    *CreateAPIClientResponse
		expectedRequestBody string
		withError           func(*testing.T, error)
	}{
		"201 Created with allAPI, cpCodes and clone group": {
			params: CreateAPIClientRequest{
				APIAccess: APIAccess{
					AllAccessibleAPIs: true,
				},
				AuthorizedUsers:   []string{"user1"},
				ClientDescription: "test_user_1 description",
				ClientName:        "test_user_1",
				ClientType:        ClientClientType,
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: true,
				},
				NotificationEmails: []string{"user1@example.com"},
				PurgeOptions: &PurgeOptions{
					CPCodeAccess: CPCodeAccess{
						AllCurrentAndNewCPCodes: true,
					},
				},
			},
			expectedPath:   "/identity-management/v3/api-clients",
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "clientId": "abcdefgh12345678",
    "clientName": "test_user_1",
    "clientDescription": "test_user_1 description",
    "clientType": "CLIENT",
    "authorizedUsers": [
        "user1"
    ],
    "canAutoCreateCredential": false,
    "notificationEmails": [
        "user1@example.com"
    ],
    "activeCredentialCount": 0,
    "allowAccountSwitch": false,
    "createdDate": "2024-07-16T23:01:50.000Z",
    "createdBy": "admin",
    "isLocked": false,
    "groupAccess": {
        "cloneAuthorizedUserGroups": true,
        "groups": [
            {
                "groupId": 123,
                "groupName": "GroupName-G-R0UP",
                "roleId": 1,
                "roleName": "Admin",
                "roleDescription": "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
                "isBlocked": false,
                "subGroups": []
            }
        ]
    },
    "apiAccess": {
        "allAccessibleApis": true,
        "apis": [
            {
                "apiId": 1,
                "apiName": "API Client Administration",
                "description": "API Client Administration",
                "endPoint": "/identity-management",
                "documentationUrl": "https://developer.akamai.com",
                "accessLevel": "READ-ONLY"
            },
            {
                "apiId": 2,
                "apiName": "CCU APIs",
                "description": "Content control utility APIs",
                "endPoint": "/ccu",
                "documentationUrl": "https://developer.akamai.com",
                "accessLevel": "READ-WRITE"
            }
        ]
    },
    "purgeOptions": {
        "canPurgeByCpcode": false,
        "canPurgeByCacheTag": false,
        "cpcodeAccess": {
            "allCurrentAndNewCpcodes": true,
            "cpcodes": []
        }
    },
    "baseURL": "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
    "accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
    "actions": {
        "editGroups": true,
        "editApis": true,
        "lock": true,
        "unlock": false,
        "editAuth": true,
        "edit": true,
        "editSwitchAccount": false,
        "transfer": true,
        "editIpAcl": true,
        "delete": true,
        "deactivateAll": false
    }
}`,
			expectedRequestBody: `
{
  "allowAccountSwitch": false,
  "apiAccess": {
    "allAccessibleApis": true,
    "apis": null
  },
  "authorizedUsers": [
    "user1"
  ],
  "canAutoCreateCredential": false,
  "clientDescription": "test_user_1 description",
  "clientName": "test_user_1",
  "clientType": "CLIENT",
  "createCredential": false,
  "groupAccess": {
    "cloneAuthorizedUserGroups": true,
    "groups": null
  },
  "notificationEmails": [
    "user1@example.com"
  ],
  "purgeOptions": {
    "canPurgeByCacheTag": false,
    "canPurgeByCpcode": false,
    "cpcodeAccess": {
      "allCurrentAndNewCpcodes": true,
      "cpcodes": null
    }
  }
}
`,
			expectedResponse: &CreateAPIClientResponse{
				AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: false,
				ActiveCredentialCount:   0,
				ClientDescription:       "test_user_1 description",
				ClientID:                "abcdefgh12345678",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreatedBy:               "admin",
				CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
				IsLocked:                false,
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: true,
					Groups: []ClientGroup{
						{
							GroupID:         123,
							GroupName:       "GroupName-G-R0UP",
							RoleID:          1,
							RoleName:        "Admin",
							RoleDescription: "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
							Subgroups:       []ClientGroup{},
						},
					},
				},
				NotificationEmails: []string{"user1@example.com"},
				APIAccess: APIAccess{
					AllAccessibleAPIs: true,
					APIs: []API{
						{
							APIID:            1,
							APIName:          "API Client Administration",
							Description:      "API Client Administration",
							Endpoint:         "/identity-management",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadOnlyLevel,
						},
						{
							APIID:            2,
							APIName:          "CCU APIs",
							Description:      "Content control utility APIs",
							Endpoint:         "/ccu",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadWriteLevel,
						},
					},
				},
				PurgeOptions: &PurgeOptions{
					CanPurgeByCPCode:   false,
					CanPurgeByCacheTag: false,
					CPCodeAccess: CPCodeAccess{
						AllCurrentAndNewCPCodes: true,
						CPCodes:                 []int64{},
					},
				},
				BaseURL: "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
				Actions: &APIClientActions{
					EditGroups:        true,
					EditAPIs:          true,
					Lock:              true,
					Unlock:            false,
					EditAuth:          true,
					Edit:              true,
					EditSwitchAccount: false,
					Transfer:          true,
					EditIPACL:         true,
					Delete:            true,
					DeactivateAll:     false,
				},
			},
		},
		"201 Created with all fields and custom API and group": {
			params: CreateAPIClientRequest{
				AllowAccountSwitch: true,
				APIAccess: APIAccess{
					AllAccessibleAPIs: false,
					APIs: []API{
						{
							AccessLevel: ReadOnlyLevel,
							APIID:       1,
						},
						{
							AccessLevel: ReadWriteLevel,
							APIID:       2,
						},
					},
				},
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: true,
				ClientDescription:       "test_user_1 description",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreateCredential:        true,
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: false,
					Groups: []ClientGroup{
						{
							GroupID: 123,
							RoleID:  1,
						},
					},
				},
				IPACL: &IPACL{
					CIDR:   []string{"1.2.3.4/32"},
					Enable: true,
				},
				NotificationEmails: []string{"user1@example.com"},
				PurgeOptions: &PurgeOptions{
					CanPurgeByCacheTag: true,
					CanPurgeByCPCode:   true,
					CPCodeAccess: CPCodeAccess{
						AllCurrentAndNewCPCodes: false,
						CPCodes:                 []int64{321},
					},
				},
			},
			expectedPath:   "/identity-management/v3/api-clients",
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "clientId": "abcdefgh12345678",
    "clientName": "test_user_1",
    "clientDescription": "test_user_1 description",
    "clientType": "CLIENT",
    "authorizedUsers": [
        "user1"
    ],
    "canAutoCreateCredential": true,
    "notificationEmails": [
        "user1@example.com"
    ],
    "activeCredentialCount": 1,
    "allowAccountSwitch": true,
    "createdDate": "2024-07-16T23:01:50.000Z",
    "createdBy": "admin",
    "isLocked": false,
    "groupAccess": {
        "cloneAuthorizedUserGroups": false,
        "groups": [
            {
                "groupId": 123,
                "groupName": "GroupName-G-R0UP",
                "roleId": 1,
                "roleName": "Admin",
                "roleDescription": "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
                "isBlocked": false,
                "subGroups": []
            }
        ]
    },
    "apiAccess": {
        "allAccessibleApis": false,
        "apis": [
            {
                "apiId": 1,
                "apiName": "API Client Administration",
                "description": "API Client Administration",
                "endPoint": "/identity-management",
                "documentationUrl": "https://developer.akamai.com",
                "accessLevel": "READ-ONLY"
            },
            {
                "apiId": 2,
                "apiName": "CCU APIs",
                "description": "Content control utility APIs",
                "endPoint": "/ccu",
                "documentationUrl": "https://developer.akamai.com",
                "accessLevel": "READ-WRITE"
            }
        ]
    },
    "purgeOptions": {
        "canPurgeByCpcode": true,
        "canPurgeByCacheTag": true,
        "cpcodeAccess": {
            "allCurrentAndNewCpcodes": false,
            "cpcodes": [321]
        }
    },
    "baseURL": "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
    "accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
    "credentials": [
        {
            "credentialId": 456,
            "clientToken": "akaa-bc78bc78bc78bc78-bc78bc78bc78bc78",
            "clientSecret": "verysecretsecret",
            "status": "ACTIVE",
            "createdOn": "2023-01-03T07:44:08.000Z",
            "description": "desc",
            "expiresOn": "2025-01-03T07:44:08.000Z",
            "actions": {
                "deactivate": true,
                "delete": false,
                "activate": false,
                "editDescription": true,
                "editExpiration": true
            }
        }
    ],
    "actions": {
        "editGroups": true,
        "editApis": true,
        "lock": true,
        "unlock": false,
        "editAuth": true,
        "edit": true,
        "editSwitchAccount": false,
        "transfer": true,
        "editIpAcl": true,
        "delete": true,
        "deactivateAll": false
    }
}`,
			expectedRequestBody: `
{
  "allowAccountSwitch": true,
  "apiAccess": {
    "allAccessibleApis": false,
    "apis": [
      {
        "accessLevel": "READ-ONLY",
        "apiId": 1,
        "apiName": "",
        "description": "",
        "documentationUrl": "",
        "endPoint": ""
      },
      {
        "accessLevel": "READ-WRITE",
        "apiId": 2,
        "apiName": "",
        "description": "",
        "documentationUrl": "",
        "endPoint": ""
      }
    ]
  },
  "authorizedUsers": [
    "user1"
  ],
  "canAutoCreateCredential": true,
  "clientDescription": "test_user_1 description",
  "clientName": "test_user_1",
  "clientType": "CLIENT",
  "createCredential": true,
  "groupAccess": {
    "cloneAuthorizedUserGroups": false,
    "groups": [
      {
        "groupId": 123,
        "groupName": "",
        "isBlocked": false,
        "parentGroupId": 0,
        "roleDescription": "",
        "roleId": 1,
        "roleName": "",
        "subgroups": null
      }
    ]
  },
  "ipAcl": {
    "cidr": [
      "1.2.3.4/32"
    ],
    "enable": true
  },
  "notificationEmails": [
    "user1@example.com"
  ],
  "purgeOptions": {
    "canPurgeByCacheTag": true,
    "canPurgeByCpcode": true,
    "cpcodeAccess": {
      "allCurrentAndNewCpcodes": false,
      "cpcodes": [
        321
      ]
    }
  }
}
`,
			expectedResponse: &CreateAPIClientResponse{
				AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
				AllowAccountSwitch:      true,
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: true,
				ActiveCredentialCount:   1,
				ClientDescription:       "test_user_1 description",
				ClientID:                "abcdefgh12345678",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreatedBy:               "admin",
				CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
				IsLocked:                false,
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: false,
					Groups: []ClientGroup{
						{
							GroupID:         123,
							GroupName:       "GroupName-G-R0UP",
							RoleID:          1,
							RoleName:        "Admin",
							RoleDescription: "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
							Subgroups:       []ClientGroup{},
						},
					},
				},
				NotificationEmails: []string{"user1@example.com"},
				APIAccess: APIAccess{
					AllAccessibleAPIs: false,
					APIs: []API{
						{
							APIID:            1,
							APIName:          "API Client Administration",
							Description:      "API Client Administration",
							Endpoint:         "/identity-management",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadOnlyLevel,
						},
						{
							APIID:            2,
							APIName:          "CCU APIs",
							Description:      "Content control utility APIs",
							Endpoint:         "/ccu",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadWriteLevel,
						},
					},
				},
				PurgeOptions: &PurgeOptions{
					CanPurgeByCacheTag: true,
					CanPurgeByCPCode:   true,
					CPCodeAccess: CPCodeAccess{
						AllCurrentAndNewCPCodes: false,
						CPCodes:                 []int64{321},
					},
				},
				BaseURL: "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
				Credentials: []CreateAPIClientCredential{
					{
						Actions: CredentialActions{
							Deactivate:      true,
							Delete:          false,
							Activate:        false,
							EditDescription: true,
							EditExpiration:  true,
						},
						ClientToken:  "akaa-bc78bc78bc78bc78-bc78bc78bc78bc78",
						ClientSecret: "verysecretsecret",
						CreatedOn:    test.NewTimeFromString(t, "2023-01-03T07:44:08.000Z"),
						CredentialID: 456,
						Description:  "desc",
						ExpiresOn:    test.NewTimeFromString(t, "2025-01-03T07:44:08.000Z"),
						Status:       CredentialActive,
					},
				},
				Actions: &APIClientActions{
					EditGroups:        true,
					EditAPIs:          true,
					Lock:              true,
					Unlock:            false,
					EditAuth:          true,
					Edit:              true,
					EditSwitchAccount: false,
					Transfer:          true,
					EditIPACL:         true,
					Delete:            true,
					DeactivateAll:     false,
				},
			},
		},
		"validation errors": {
			params: CreateAPIClientRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create api client: struct validation:\nAPIAccess: {\n\tAPIs: cannot be blank\n}\nAuthorizedUsers: cannot be blank\nClientType: cannot be blank\nGroupAccess: {\n\tGroups: cannot be blank\n}", err.Error())
			},
		},
		"validation errors - internal validations": {
			params: CreateAPIClientRequest{APIAccess: APIAccess{APIs: []API{{}}}, AuthorizedUsers: []string{"user1"}, ClientType: "abc", GroupAccess: GroupAccess{Groups: []ClientGroup{{}}}, PurgeOptions: &PurgeOptions{CPCodeAccess: CPCodeAccess{AllCurrentAndNewCPCodes: false, CPCodes: nil}}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create api client: struct validation:\nAPIAccess: {\n\tAPIs[0]: {\n\t\tAPIID: cannot be blank\n\t\tAccessLevel: cannot be blank\n\t}\n}\nClientType: value 'abc' is invalid. Must be one of: 'CLIENT', 'SERVICE_ACCOUNT' or 'USER_CLIENT'\nGroupAccess: {\n\tGroups[0]: {\n\t\tGroupID: cannot be blank\n\t\tRoleID: cannot be blank\n\t}\n}\nPurgeOptions: {\n\tCPCodeAccess: {\n\t\tCPCodes: is required\n\t}\n}", err.Error())
			},
		},
		"500 internal server error": {
			params: CreateAPIClientRequest{
				APIAccess: APIAccess{
					AllAccessibleAPIs: true,
				},
				AuthorizedUsers:   []string{"user1"},
				ClientDescription: "test_user_1 description",
				ClientName:        "test_user_1",
				ClientType:        ClientClientType,
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: true,
				},
				NotificationEmails: []string{"user1@example.com"},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}
`,
			expectedPath: "/identity-management/v3/api-clients",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)

				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			response, err := client.CreateAPIClient(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestIAM_UpdateAPIClient(t *testing.T) {
	tests := map[string]struct {
		params              UpdateAPIClientRequest
		expectedPath        string
		responseStatus      int
		responseBody        string
		expectedResponse    *UpdateAPIClientResponse
		expectedRequestBody string
		withError           func(*testing.T, error)
	}{
		"200 Updated self": {
			params: UpdateAPIClientRequest{
				Body: UpdateAPIClientRequestBody{
					APIAccess: APIAccess{
						AllAccessibleAPIs: true,
					},
					AuthorizedUsers:   []string{"user1"},
					ClientDescription: "test_user_1 description",
					ClientName:        "test_user_1",
					ClientType:        ClientClientType,
					GroupAccess: GroupAccess{
						CloneAuthorizedUserGroups: true,
					},
					NotificationEmails: []string{"user1@example.com"},
					PurgeOptions: &PurgeOptions{
						CPCodeAccess: CPCodeAccess{
							AllCurrentAndNewCPCodes: true,
						},
					},
				},
			},
			expectedPath:   "/identity-management/v3/api-clients/self",
			responseStatus: http.StatusOK,
			responseBody: `
{
   "clientId": "abcdefgh12345678",
   "clientName": "test_user_1",
   "clientDescription": "test_user_1 description",
   "clientType": "CLIENT",
   "authorizedUsers": [
       "user1"
   ],
   "canAutoCreateCredential": false,
   "notificationEmails": [
       "user1@example.com"
   ],
   "activeCredentialCount": 0,
   "allowAccountSwitch": false,
   "createdDate": "2024-07-16T23:01:50.000Z",
   "createdBy": "admin",
   "isLocked": false,
   "groupAccess": {
       "cloneAuthorizedUserGroups": true,
       "groups": [
           {
               "groupId": 123,
               "groupName": "GroupName-G-R0UP",
               "roleId": 1,
               "roleName": "Admin",
               "roleDescription": "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
               "isBlocked": false,
               "subGroups": []
           }
       ]
   },
   "apiAccess": {
       "allAccessibleApis": true,
       "apis": [
           {
               "apiId": 1,
               "apiName": "API Client Administration",
               "description": "API Client Administration",
               "endPoint": "/identity-management",
               "documentationUrl": "https://developer.akamai.com",
               "accessLevel": "READ-ONLY"
           },
           {
               "apiId": 2,
               "apiName": "CCU APIs",
               "description": "Content control utility APIs",
               "endPoint": "/ccu",
               "documentationUrl": "https://developer.akamai.com",
               "accessLevel": "READ-WRITE"
           }
       ]
   },
   "purgeOptions": {
       "canPurgeByCpcode": false,
       "canPurgeByCacheTag": false,
       "cpcodeAccess": {
           "allCurrentAndNewCpcodes": true,
           "cpcodes": []
       }
   },
   "baseURL": "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
   "accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
   "actions": {
       "editGroups": true,
       "editApis": true,
       "lock": true,
       "unlock": false,
       "editAuth": true,
       "edit": true,
       "editSwitchAccount": false,
       "transfer": true,
       "editIpAcl": true,
       "delete": true,
       "deactivateAll": false
   }
}`,
			expectedResponse: &UpdateAPIClientResponse{
				AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: false,
				ActiveCredentialCount:   0,
				ClientDescription:       "test_user_1 description",
				ClientID:                "abcdefgh12345678",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreatedBy:               "admin",
				CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
				IsLocked:                false,
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: true,
					Groups: []ClientGroup{
						{
							GroupID:         123,
							GroupName:       "GroupName-G-R0UP",
							RoleID:          1,
							RoleName:        "Admin",
							RoleDescription: "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
							Subgroups:       []ClientGroup{},
						},
					},
				},
				NotificationEmails: []string{"user1@example.com"},
				APIAccess: APIAccess{
					AllAccessibleAPIs: true,
					APIs: []API{
						{
							APIID:            1,
							APIName:          "API Client Administration",
							Description:      "API Client Administration",
							Endpoint:         "/identity-management",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadOnlyLevel,
						},
						{
							APIID:            2,
							APIName:          "CCU APIs",
							Description:      "Content control utility APIs",
							Endpoint:         "/ccu",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadWriteLevel,
						},
					},
				},
				PurgeOptions: &PurgeOptions{
					CanPurgeByCPCode:   false,
					CanPurgeByCacheTag: false,
					CPCodeAccess: CPCodeAccess{
						AllCurrentAndNewCPCodes: true,
						CPCodes:                 []int64{},
					},
				},
				BaseURL: "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
				Actions: &APIClientActions{
					EditGroups:        true,
					EditAPIs:          true,
					Lock:              true,
					Unlock:            false,
					EditAuth:          true,
					Edit:              true,
					EditSwitchAccount: false,
					Transfer:          true,
					EditIPACL:         true,
					Delete:            true,
					DeactivateAll:     false,
				},
			},
		},
		"200 Updated with allAPI, cpCodes and clone group": {
			params: UpdateAPIClientRequest{
				ClientID: "abcdefgh12345678",
				Body: UpdateAPIClientRequestBody{
					APIAccess: APIAccess{
						AllAccessibleAPIs: true,
					},
					AuthorizedUsers:   []string{"user1"},
					ClientDescription: "test_user_1 description",
					ClientName:        "test_user_1",
					ClientType:        ClientClientType,
					GroupAccess: GroupAccess{
						CloneAuthorizedUserGroups: true,
					},
					NotificationEmails: []string{"user1@example.com"},
					PurgeOptions: &PurgeOptions{
						CPCodeAccess: CPCodeAccess{
							AllCurrentAndNewCPCodes: true,
						},
					},
				},
			},
			expectedPath:   "/identity-management/v3/api-clients/abcdefgh12345678",
			responseStatus: http.StatusOK,
			responseBody: `
{
   "clientId": "abcdefgh12345678",
   "clientName": "test_user_1",
   "clientDescription": "test_user_1 description",
   "clientType": "CLIENT",
   "authorizedUsers": [
       "user1"
   ],
   "canAutoCreateCredential": false,
   "notificationEmails": [
       "user1@example.com"
   ],
   "activeCredentialCount": 0,
   "allowAccountSwitch": false,
   "createdDate": "2024-07-16T23:01:50.000Z",
   "createdBy": "admin",
   "isLocked": false,
   "groupAccess": {
       "cloneAuthorizedUserGroups": true,
       "groups": [
           {
               "groupId": 123,
               "groupName": "GroupName-G-R0UP",
               "roleId": 1,
               "roleName": "Admin",
               "roleDescription": "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
               "isBlocked": false,
               "subGroups": []
           }
       ]
   },
   "apiAccess": {
       "allAccessibleApis": true,
       "apis": [
           {
               "apiId": 1,
               "apiName": "API Client Administration",
               "description": "API Client Administration",
               "endPoint": "/identity-management",
               "documentationUrl": "https://developer.akamai.com",
               "accessLevel": "READ-ONLY"
           },
           {
               "apiId": 2,
               "apiName": "CCU APIs",
               "description": "Content control utility APIs",
               "endPoint": "/ccu",
               "documentationUrl": "https://developer.akamai.com",
               "accessLevel": "READ-WRITE"
           }
       ]
   },
   "purgeOptions": {
       "canPurgeByCpcode": false,
       "canPurgeByCacheTag": false,
       "cpcodeAccess": {
           "allCurrentAndNewCpcodes": true,
           "cpcodes": []
       }
   },
   "baseURL": "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
   "accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
   "actions": {
       "editGroups": true,
       "editApis": true,
       "lock": true,
       "unlock": false,
       "editAuth": true,
       "edit": true,
       "editSwitchAccount": false,
       "transfer": true,
       "editIpAcl": true,
       "delete": true,
       "deactivateAll": false
   }
}`,
			expectedRequestBody: `
{
  "allowAccountSwitch": false,
  "apiAccess": {
    "allAccessibleApis": true,
    "apis": null
  },
  "authorizedUsers": [
    "user1"
  ],
  "canAutoCreateCredential": false,
  "clientDescription": "test_user_1 description",
  "clientName": "test_user_1",
  "clientType": "CLIENT",
  "groupAccess": {
    "cloneAuthorizedUserGroups": true,
    "groups": null
  },
  "notificationEmails": [
    "user1@example.com"
  ],
  "purgeOptions": {
    "canPurgeByCacheTag": false,
    "canPurgeByCpcode": false,
    "cpcodeAccess": {
      "allCurrentAndNewCpcodes": true,
      "cpcodes": null
    }
  }
}
`,
			expectedResponse: &UpdateAPIClientResponse{
				AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: false,
				ActiveCredentialCount:   0,
				ClientDescription:       "test_user_1 description",
				ClientID:                "abcdefgh12345678",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreatedBy:               "admin",
				CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
				IsLocked:                false,
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: true,
					Groups: []ClientGroup{
						{
							GroupID:         123,
							GroupName:       "GroupName-G-R0UP",
							RoleID:          1,
							RoleName:        "Admin",
							RoleDescription: "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
							Subgroups:       []ClientGroup{},
						},
					},
				},
				NotificationEmails: []string{"user1@example.com"},
				APIAccess: APIAccess{
					AllAccessibleAPIs: true,
					APIs: []API{
						{
							APIID:            1,
							APIName:          "API Client Administration",
							Description:      "API Client Administration",
							Endpoint:         "/identity-management",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadOnlyLevel,
						},
						{
							APIID:            2,
							APIName:          "CCU APIs",
							Description:      "Content control utility APIs",
							Endpoint:         "/ccu",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadWriteLevel,
						},
					},
				},
				PurgeOptions: &PurgeOptions{
					CanPurgeByCPCode:   false,
					CanPurgeByCacheTag: false,
					CPCodeAccess: CPCodeAccess{
						AllCurrentAndNewCPCodes: true,
						CPCodes:                 []int64{},
					},
				},
				BaseURL: "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
				Actions: &APIClientActions{
					EditGroups:        true,
					EditAPIs:          true,
					Lock:              true,
					Unlock:            false,
					EditAuth:          true,
					Edit:              true,
					EditSwitchAccount: false,
					Transfer:          true,
					EditIPACL:         true,
					Delete:            true,
					DeactivateAll:     false,
				},
			},
		},
		"200 Updated with all fields and custom API and group": {
			params: UpdateAPIClientRequest{
				ClientID: "abcdefgh12345678",
				Body: UpdateAPIClientRequestBody{
					AllowAccountSwitch: true,
					APIAccess: APIAccess{
						AllAccessibleAPIs: false,
						APIs: []API{
							{
								AccessLevel: ReadOnlyLevel,
								APIID:       1,
							},
							{
								AccessLevel: ReadWriteLevel,
								APIID:       2,
							},
						},
					},
					AuthorizedUsers:         []string{"user1"},
					CanAutoCreateCredential: true,
					ClientDescription:       "test_user_1 description",
					ClientName:              "test_user_1",
					ClientType:              ClientClientType,
					GroupAccess: GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []ClientGroup{
							{
								GroupID: 123,
								RoleID:  1,
							},
						},
					},
					IPACL: &IPACL{
						CIDR:   []string{"1.2.3.4/32"},
						Enable: true,
					},
					NotificationEmails: []string{"user1@example.com"},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCacheTag: true,
						CanPurgeByCPCode:   true,
						CPCodeAccess: CPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{321},
						},
					},
				},
			},
			expectedPath:   "/identity-management/v3/api-clients/abcdefgh12345678",
			responseStatus: http.StatusOK,
			responseBody: `
{
   "clientId": "abcdefgh12345678",
   "clientName": "test_user_1",
   "clientDescription": "test_user_1 description",
   "clientType": "CLIENT",
   "authorizedUsers": [
	   "user1"
   ],
   "canAutoCreateCredential": true,
   "notificationEmails": [
	   "user1@example.com"
   ],
   "activeCredentialCount": 1,
   "allowAccountSwitch": true,
   "createdDate": "2024-07-16T23:01:50.000Z",
   "createdBy": "admin",
   "isLocked": false,
   "groupAccess": {
	   "cloneAuthorizedUserGroups": false,
	   "groups": [
		   {
			   "groupId": 123,
			   "groupName": "GroupName-G-R0UP",
			   "roleId": 1,
			   "roleName": "Admin",
			   "roleDescription": "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
			   "isBlocked": false,
			   "subGroups": []
		   }
	   ]
   },
   "apiAccess": {
	   "allAccessibleApis": false,
	   "apis": [
		   {
			   "apiId": 1,
			   "apiName": "API Client Administration",
			   "description": "API Client Administration",
			   "endPoint": "/identity-management",
			   "documentationUrl": "https://developer.akamai.com",
			   "accessLevel": "READ-ONLY"
		   },
		   {
			   "apiId": 2,
			   "apiName": "CCU APIs",
			   "description": "Content control utility APIs",
			   "endPoint": "/ccu",
			   "documentationUrl": "https://developer.akamai.com",
			   "accessLevel": "READ-WRITE"
		   }
	   ]
   },
   "purgeOptions": {
	   "canPurgeByCpcode": true,
	   "canPurgeByCacheTag": true,
	   "cpcodeAccess": {
		   "allCurrentAndNewCpcodes": false,
		   "cpcodes": [321]
	   }
   },
   "baseURL": "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
   "accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
   "credentials": [
        {
            "credentialId": 456,
            "clientToken": "akaa-bc78bc78bc78bc78-bc78bc78bc78bc78",
            "status": "ACTIVE",
            "createdOn": "2023-01-03T07:44:08.000Z",
            "description": "desc",
            "expiresOn": "2025-01-03T07:44:08.000Z",
            "actions": {
                "deactivate": true,
                "delete": false,
                "activate": false,
                "editDescription": true,
                "editExpiration": true
            }
        }
    ],
   "actions": {
	   "editGroups": true,
	   "editApis": true,
	   "lock": true,
	   "unlock": false,
	   "editAuth": true,
	   "edit": true,
	   "editSwitchAccount": false,
	   "transfer": true,
	   "editIpAcl": true,
	   "delete": true,
	   "deactivateAll": false
   }
}`,
			expectedRequestBody: `
{
  "allowAccountSwitch": true,
  "apiAccess": {
    "allAccessibleApis": false,
    "apis": [
      {
        "accessLevel": "READ-ONLY",
        "apiId": 1,
        "apiName": "",
        "description": "",
        "documentationUrl": "",
        "endPoint": ""
      },
      {
        "accessLevel": "READ-WRITE",
        "apiId": 2,
        "apiName": "",
        "description": "",
        "documentationUrl": "",
        "endPoint": ""
      }
    ]
  },
  "authorizedUsers": [
    "user1"
  ],
  "canAutoCreateCredential": true,
  "clientDescription": "test_user_1 description",
  "clientName": "test_user_1",
  "clientType": "CLIENT",
  "groupAccess": {
    "cloneAuthorizedUserGroups": false,
    "groups": [
      {
        "groupId": 123,
        "groupName": "",
        "isBlocked": false,
        "parentGroupId": 0,
        "roleDescription": "",
        "roleId": 1,
        "roleName": "",
        "subgroups": null
      }
    ]
  },
  "ipAcl": {
    "cidr": [
      "1.2.3.4/32"
    ],
    "enable": true
  },
  "notificationEmails": [
    "user1@example.com"
  ],
  "purgeOptions": {
    "canPurgeByCacheTag": true,
    "canPurgeByCpcode": true,
    "cpcodeAccess": {
      "allCurrentAndNewCpcodes": false,
      "cpcodes": [
        321
      ]
    }
  }
}
`,
			expectedResponse: &UpdateAPIClientResponse{
				AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
				AllowAccountSwitch:      true,
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: true,
				ActiveCredentialCount:   1,
				ClientDescription:       "test_user_1 description",
				ClientID:                "abcdefgh12345678",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreatedBy:               "admin",
				CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
				IsLocked:                false,
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: false,
					Groups: []ClientGroup{
						{
							GroupID:         123,
							GroupName:       "GroupName-G-R0UP",
							RoleID:          1,
							RoleName:        "Admin",
							RoleDescription: "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
							Subgroups:       []ClientGroup{},
						},
					},
				},
				NotificationEmails: []string{"user1@example.com"},
				APIAccess: APIAccess{
					AllAccessibleAPIs: false,
					APIs: []API{
						{
							APIID:            1,
							APIName:          "API Client Administration",
							Description:      "API Client Administration",
							Endpoint:         "/identity-management",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadOnlyLevel,
						},
						{
							APIID:            2,
							APIName:          "CCU APIs",
							Description:      "Content control utility APIs",
							Endpoint:         "/ccu",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadWriteLevel,
						},
					},
				},
				PurgeOptions: &PurgeOptions{
					CanPurgeByCacheTag: true,
					CanPurgeByCPCode:   true,
					CPCodeAccess: CPCodeAccess{
						AllCurrentAndNewCPCodes: false,
						CPCodes:                 []int64{321},
					},
				},
				BaseURL: "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
				Credentials: []APIClientCredential{
					{
						Actions: CredentialActions{
							Deactivate:      true,
							Delete:          false,
							Activate:        false,
							EditDescription: true,
							EditExpiration:  true,
						},
						ClientToken:  "akaa-bc78bc78bc78bc78-bc78bc78bc78bc78",
						CreatedOn:    test.NewTimeFromString(t, "2023-01-03T07:44:08.000Z"),
						CredentialID: 456,
						Description:  "desc",
						ExpiresOn:    test.NewTimeFromString(t, "2025-01-03T07:44:08.000Z"),
						Status:       CredentialActive,
					},
				},
				Actions: &APIClientActions{
					EditGroups:        true,
					EditAPIs:          true,
					Lock:              true,
					Unlock:            false,
					EditAuth:          true,
					Edit:              true,
					EditSwitchAccount: false,
					Transfer:          true,
					EditIPACL:         true,
					Delete:            true,
					DeactivateAll:     false,
				},
			},
		},
		"validation errors": {
			params: UpdateAPIClientRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update api client: struct validation:\nBody: {\n\tAPIAccess: {\n\t\tAPIs: cannot be blank\n\t}\n\tAuthorizedUsers: cannot be blank\n\tClientName: cannot be blank\n\tClientType: cannot be blank\n\tGroupAccess: {\n\t\tGroups: cannot be blank\n\t}\n}", err.Error())
			},
		},
		"validation errors - internal validations": {
			params: UpdateAPIClientRequest{Body: UpdateAPIClientRequestBody{APIAccess: APIAccess{APIs: []API{{}}}, AuthorizedUsers: []string{"user1"}, ClientType: "abc", GroupAccess: GroupAccess{Groups: []ClientGroup{{}}}, PurgeOptions: &PurgeOptions{CPCodeAccess: CPCodeAccess{AllCurrentAndNewCPCodes: false, CPCodes: nil}}}},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update api client: struct validation:\nBody: {\n\tAPIAccess: {\n\t\tAPIs[0]: {\n\t\t\tAPIID: cannot be blank\n\t\t\tAccessLevel: cannot be blank\n\t\t}\n\t}\n\tClientName: cannot be blank\n\tClientType: value 'abc' is invalid. Must be one of: 'CLIENT', 'SERVICE_ACCOUNT' or 'USER_CLIENT'\n\tGroupAccess: {\n\t\tGroups[0]: {\n\t\t\tGroupID: cannot be blank\n\t\t\tRoleID: cannot be blank\n\t\t}\n\t}\n\tPurgeOptions: {\n\t\tCPCodeAccess: {\n\t\t\tCPCodes: is required\n\t\t}\n\t}\n}", err.Error())
			},
		},
		"500 internal server error": {
			params: UpdateAPIClientRequest{
				Body: UpdateAPIClientRequestBody{
					APIAccess: APIAccess{
						AllAccessibleAPIs: true,
					},
					AuthorizedUsers:   []string{"user1"},
					ClientDescription: "test_user_1 description",
					ClientName:        "test_user_1",
					ClientType:        ClientClientType,
					GroupAccess: GroupAccess{
						CloneAuthorizedUserGroups: true,
					},
					NotificationEmails: []string{"user1@example.com"},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/identity-management/v3/api-clients/self",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)

				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(body))
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			response, err := client.UpdateAPIClient(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestIAM_GetAPIClient(t *testing.T) {
	tests := map[string]struct {
		params           GetAPIClientRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAPIClientResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - Self": {
			params:         GetAPIClientRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"clientId": "abcdefgh12345678",
	"clientName": "test_user_1",
	"clientDescription": "test_user_1 description",
	"clientType": "CLIENT",
	"authorizedUsers": [
		"user1"
	],
	"canAutoCreateCredential": false,
	"notificationEmails": [
		"user1@example.com"
	],
	"activeCredentialCount": 0,
	"allowAccountSwitch": false,
	"createdDate": "2024-07-16T23:01:50.000Z",
	"createdBy": "admin",
	"isLocked": false,
	"baseURL": "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
	"accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d"
}
`,
			expectedPath: "/identity-management/v3/api-clients/self?actions=false&apiAccess=false&credentials=false&groupAccess=false&ipAcl=false",
			expectedResponse: &GetAPIClientResponse{
				AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
				BaseURL:                 "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
				ActiveCredentialCount:   0,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: false,
				ClientDescription:       "test_user_1 description",
				ClientID:                "abcdefgh12345678",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreatedBy:               "admin",
				CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
				IsLocked:                false,
				NotificationEmails:      []string{"user1@example.com"},
			},
		},
		"200 OK - with clientID": {
			params:         GetAPIClientRequest{ClientID: "abcdefgh12345678"},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"clientId": "abcdefgh12345678",
	"clientName": "test_user_1",
	"clientDescription": "test_user_1 description",
	"clientType": "CLIENT",
	"authorizedUsers": [
		"user1"
	],
	"canAutoCreateCredential": false,
	"notificationEmails": [
		"user1@example.com"
	],
	"activeCredentialCount": 0,
	"allowAccountSwitch": false,
	"createdDate": "2024-07-16T23:01:50.000Z",
	"createdBy": "admin",
	"isLocked": false,
	"baseURL": "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
	"accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d"
}
`,
			expectedPath: "/identity-management/v3/api-clients/abcdefgh12345678?actions=false&apiAccess=false&credentials=false&groupAccess=false&ipAcl=false",
			expectedResponse: &GetAPIClientResponse{
				AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
				BaseURL:                 "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
				ActiveCredentialCount:   0,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: false,
				ClientDescription:       "test_user_1 description",
				ClientID:                "abcdefgh12345678",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreatedBy:               "admin",
				CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
				IsLocked:                false,
				NotificationEmails:      []string{"user1@example.com"},
			},
		},
		"200 OK - with all query params and all fields": {
			params: GetAPIClientRequest{
				ClientID:    "abcdefgh12345678",
				Actions:     true,
				GroupAccess: true,
				APIAccess:   true,
				Credentials: true,
				IPACL:       true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"clientId": "abcdefgh12345678",
	"clientName": "test_user_1",
	"clientDescription": "test_user_1 description",
	"clientType": "CLIENT",
	"authorizedUsers": [
		"user1"
	],
	"canAutoCreateCredential": false,
	"notificationEmails": [
		"user1@example.com"
	],
	"activeCredentialCount": 0,
	"allowAccountSwitch": false,
	"createdDate": "2024-07-16T23:01:50.000Z",
	"createdBy": "admin",
	"isLocked": false,
	"baseURL": "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
	"accessToken": "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
	"groupAccess": {
        "cloneAuthorizedUserGroups": true,
        "groups": [
            {
                "groupId": 123,
                "groupName": "GroupName-G-R0UP",
                "roleId": 1,
                "roleName": "Admin",
                "roleDescription": "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
                "isBlocked": false,
                "subGroups": []
            }
        ]
    },
    "apiAccess": {
        "allAccessibleApis": true,
        "apis": [
            {
                "apiId": 1,
                "apiName": "API Client Administration",
                "description": "API Client Administration",
                "endPoint": "/identity-management",
                "documentationUrl": "https://developer.akamai.com",
                "accessLevel": "READ-ONLY"
            },
            {
                "apiId": 2,
                "apiName": "CCU APIs",
                "description": "Content control utility APIs",
                "endPoint": "/ccu",
                "documentationUrl": "https://developer.akamai.com",
                "accessLevel": "READ-WRITE"
            }
        ]
    },
    "purgeOptions": {
        "canPurgeByCpcode": false,
        "canPurgeByCacheTag": false,
        "cpcodeAccess": {
            "allCurrentAndNewCpcodes": true,
            "cpcodes": []
        }
    },
    "credentials": [
        {
            "credentialId": 456,
            "clientToken": "akaa-bc78bc78bc78bc78-bc78bc78bc78bc78",
            "status": "ACTIVE",
            "createdOn": "2023-01-03T07:44:08.000Z",
            "description": "desc",
            "expiresOn": "2025-01-03T07:44:08.000Z",
            "actions": {
                "deactivate": true,
                "delete": false,
                "activate": false,
                "editDescription": true,
                "editExpiration": true
            }
        },
        {
            "credentialId": 678,
            "clientToken": "akaa-de90de90de90de90-de90de90de90de90",
            "status": "INACTIVE",
            "createdOn": "2023-01-03T07:44:08.000Z",
            "description": "",
            "expiresOn": "2025-01-03T07:44:08.000Z",
            "actions": {
                "deactivate": false,
                "delete": false,
                "activate": false,
                "editDescription": false,
                "editExpiration": false
            }
        },
        {
            "credentialId": 789,
            "clientToken": "akaa-gh34gh34gh34gh34-gh34gh34gh34gh34",
            "status": "DELETED",
            "createdOn": "2023-01-03T07:44:08.000Z",
            "description": "",
            "expiresOn": "2025-01-03T07:44:08.000Z",
            "actions": {
                "deactivate": false,
                "delete": false,
                "activate": false,
                "editDescription": false,
                "editExpiration": false
            }
        }
    ],
	"actions": {
        "editGroups": true,
        "editApis": true,
        "lock": false,
        "unlock": false,
        "editAuth": true,
        "edit": true,
        "editSwitchAccount": false,
        "transfer": true,
        "editIpAcl": true,
        "delete": false,
        "deactivateAll": true
    }
}
`,
			expectedPath: "/identity-management/v3/api-clients/abcdefgh12345678?actions=true&apiAccess=true&credentials=true&groupAccess=true&ipAcl=true",
			expectedResponse: &GetAPIClientResponse{
				AccessToken:             "akaa-1a2b3c4d1a2b3c4d-1a2b3c4d1a2b3c4d",
				BaseURL:                 "https://akaa-d4c3b2a1d4c3b2a1-d4c3b2a1d4c3b2a1.luna-dev.akamaiapis.net",
				ActiveCredentialCount:   0,
				AllowAccountSwitch:      false,
				AuthorizedUsers:         []string{"user1"},
				CanAutoCreateCredential: false,
				ClientDescription:       "test_user_1 description",
				ClientID:                "abcdefgh12345678",
				ClientName:              "test_user_1",
				ClientType:              ClientClientType,
				CreatedBy:               "admin",
				CreatedDate:             test.NewTimeFromString(t, "2024-07-16T23:01:50.000Z"),
				IsLocked:                false,
				NotificationEmails:      []string{"user1@example.com"},
				GroupAccess: GroupAccess{
					CloneAuthorizedUserGroups: true,
					Groups: []ClientGroup{
						{
							GroupID:         123,
							GroupName:       "GroupName-G-R0UP",
							RoleID:          1,
							RoleName:        "Admin",
							RoleDescription: "This role provides the maximum access to users. An Administrator can perform admin tasks such as creating users and groups; configuration-related tasks such as creating and editing configurations; publishing tasks",
							Subgroups:       []ClientGroup{},
						},
					},
				},
				APIAccess: APIAccess{
					AllAccessibleAPIs: true,
					APIs: []API{
						{
							APIID:            1,
							APIName:          "API Client Administration",
							Description:      "API Client Administration",
							Endpoint:         "/identity-management",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadOnlyLevel,
						},
						{
							APIID:            2,
							APIName:          "CCU APIs",
							Description:      "Content control utility APIs",
							Endpoint:         "/ccu",
							DocumentationURL: "https://developer.akamai.com",
							AccessLevel:      ReadWriteLevel,
						},
					},
				},
				PurgeOptions: &PurgeOptions{
					CanPurgeByCPCode:   false,
					CanPurgeByCacheTag: false,
					CPCodeAccess: CPCodeAccess{
						AllCurrentAndNewCPCodes: true,
						CPCodes:                 []int64{},
					},
				},
				Credentials: []APIClientCredential{
					{
						Actions: CredentialActions{
							Deactivate:      true,
							Delete:          false,
							Activate:        false,
							EditDescription: true,
							EditExpiration:  true,
						},
						ClientToken:  "akaa-bc78bc78bc78bc78-bc78bc78bc78bc78",
						CreatedOn:    test.NewTimeFromString(t, "2023-01-03T07:44:08.000Z"),
						CredentialID: 456,
						Description:  "desc",
						ExpiresOn:    test.NewTimeFromString(t, "2025-01-03T07:44:08.000Z"),
						Status:       CredentialActive,
					},
					{
						Actions: CredentialActions{
							Deactivate:      false,
							Delete:          false,
							Activate:        false,
							EditDescription: false,
							EditExpiration:  false,
						},
						ClientToken:  "akaa-de90de90de90de90-de90de90de90de90",
						CreatedOn:    test.NewTimeFromString(t, "2023-01-03T07:44:08.000Z"),
						CredentialID: 678,
						Description:  "",
						ExpiresOn:    test.NewTimeFromString(t, "2025-01-03T07:44:08.000Z"),
						Status:       CredentialInactive,
					},
					{
						Actions: CredentialActions{
							Deactivate:      false,
							Delete:          false,
							Activate:        false,
							EditDescription: false,
							EditExpiration:  false,
						},
						ClientToken:  "akaa-gh34gh34gh34gh34-gh34gh34gh34gh34",
						CreatedOn:    test.NewTimeFromString(t, "2023-01-03T07:44:08.000Z"),
						CredentialID: 789,
						Description:  "",
						ExpiresOn:    test.NewTimeFromString(t, "2025-01-03T07:44:08.000Z"),
						Status:       CredentialDeleted,
					},
				},
				Actions: &APIClientActions{
					EditGroups:        true,
					EditAPIs:          true,
					Lock:              false,
					Unlock:            false,
					EditAuth:          true,
					Edit:              true,
					EditSwitchAccount: false,
					Transfer:          true,
					EditIPACL:         true,
					Delete:            false,
					DeactivateAll:     true,
				},
			},
		},
		"500 internal server error": {
			params:         GetAPIClientRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		  "type": "internal_error",
		  "title": "Internal Server Error",
		  "detail": "Error making request",
		  "status": 500
		}
		`,
			expectedPath: "/identity-management/v3/api-clients/self?actions=false&apiAccess=false&credentials=false&groupAccess=false&ipAcl=false",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAPIClient(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_DeleteAPIClient(t *testing.T) {
	tests := map[string]struct {
		params         DeleteAPIClientRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No content - self": {
			params:         DeleteAPIClientRequest{},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v3/api-clients/self",
		},
		"204 No content - by ID": {
			params:         DeleteAPIClientRequest{ClientID: "abcdefgh12345678"},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v3/api-clients/abcdefgh12345678",
		},
		"500 internal server error": {
			params:         DeleteAPIClientRequest{ClientID: "abcdefgh12345678"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/identity-management/v3/api-clients/abcdefgh12345678",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeleteAPIClient(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
