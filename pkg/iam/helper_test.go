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

func TestIAMListAllowedCPCodes(t *testing.T) {
	tests := map[string]struct {
		params           ListAllowedCPCodesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse ListAllowedCPCodesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListAllowedCPCodesRequest{
				UserName: "jsmith",
				ListAllowedCPCodesRequestBody: ListAllowedCPCodesRequestBody{
					ClientType: ClientClientType,
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `[
  {
    "name": "Stream Analyzer (36915)",
    "value": 36915
  },
  {
    "name": "plopessa-uvod-ns (373118)",
    "value": 373118
  },
  {
    "name": "ArunNS (866797)",
    "value": 866797
  },
  {
    "name": "1234 (933076)",
    "value": 933076
  }
]`,
			expectedPath: "/identity-management/v3/users/jsmith/allowed-cpcodes",
			expectedResponse: ListAllowedCPCodesResponse{
				{
					Name:  "Stream Analyzer (36915)",
					Value: 36915,
				},
				{
					Name:  "plopessa-uvod-ns (373118)",
					Value: 373118,
				},
				{
					Name:  "ArunNS (866797)",
					Value: 866797,
				},
				{
					Name:  "1234 (933076)",
					Value: 933076,
				},
			},
		},
		"200 OK with groups": {
			params: ListAllowedCPCodesRequest{
				UserName: "jsmith",
				ListAllowedCPCodesRequestBody: ListAllowedCPCodesRequestBody{
					ClientType: ServiceAccountClientType,
					Groups: []AllowedCPCodesGroup{
						{
							GroupID: 1,
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `[
  {
    "name": "Stream Analyzer (36915)",
    "value": 36915
  },
  {
    "name": "plopessa-uvod-ns (373118)",
    "value": 373118
  },
  {
    "name": "ArunNS (866797)",
    "value": 866797
  },
  {
    "name": "1234 (933076)",
    "value": 933076
  }
]`,
			expectedPath: "/identity-management/v3/users/jsmith/allowed-cpcodes",
			expectedResponse: ListAllowedCPCodesResponse{
				{
					Name:  "Stream Analyzer (36915)",
					Value: 36915,
				},
				{
					Name:  "plopessa-uvod-ns (373118)",
					Value: 373118,
				},
				{
					Name:  "ArunNS (866797)",
					Value: 866797,
				},
				{
					Name:  "1234 (933076)",
					Value: 933076,
				},
			},
		},
		"500 internal server error": {
			params: ListAllowedCPCodesRequest{
				UserName: "jsmith",
				ListAllowedCPCodesRequestBody: ListAllowedCPCodesRequestBody{
					ClientType: ClientClientType,
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
			expectedPath: "/identity-management/v3/users/jsmith/allowed-cpcodes",
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
		"missing user name and client type": {
			params: ListAllowedCPCodesRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "list allowed CP codes: struct validation:\nClientType: cannot be blank\nUserName: cannot be blank")
			},
		},
		"group is required for client type SERVICE_ACCOUNT": {
			params: ListAllowedCPCodesRequest{
				UserName: "jsmith",
				ListAllowedCPCodesRequestBody: ListAllowedCPCodesRequestBody{
					ClientType: ServiceAccountClientType,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "list allowed CP codes: struct validation:\nGroups: cannot be blank")
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
			result, err := client.ListAllowedCPCodes(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAMListAuthorizedUsers(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse ListAuthorizedUsersResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `[
	{
      	"username": "test.example.user",
        "firstName": "Edd",
        "lastName": "Example",
        "email": "test_example@akamai.com",
        "uiIdentityId": "X-YZ-1111111"
    },
    {
        "username": "test.example.user2",
        "firstName": "Fred",
        "lastName": "Example2",
        "email": "test_example2@akamai.com",
        "uiIdentityId": "X-YZ-2222222"
    },
    {
        "username": "test.example.user3",
        "firstName": "Ted",
        "lastName": "Example3",
        "email": "test_example3@akamai.com",
        "uiIdentityId": "X-YZ-3333333"
    }
]`,
			expectedPath: "/identity-management/v3/users",
			expectedResponse: ListAuthorizedUsersResponse{
				{
					Username:     "test.example.user",
					FirstName:    "Edd",
					LastName:     "Example",
					Email:        "test_example@akamai.com",
					UIIdentityID: "X-YZ-1111111",
				},
				{
					Username:     "test.example.user2",
					FirstName:    "Fred",
					LastName:     "Example2",
					Email:        "test_example2@akamai.com",
					UIIdentityID: "X-YZ-2222222",
				},
				{
					Username:     "test.example.user3",
					FirstName:    "Ted",
					LastName:     "Example3",
					Email:        "test_example3@akamai.com",
					UIIdentityID: "X-YZ-3333333",
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
			expectedPath: "/identity-management/v3/users",
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
			result, err := client.ListAuthorizedUsers(context.Background())
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAMListAllowedAPIs(t *testing.T) {
	tests := map[string]struct {
		params           ListAllowedAPIsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse ListAllowedAPIsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListAllowedAPIsRequest{
				UserName: "jsmith",
			},
			responseStatus: http.StatusOK,
			responseBody: `[
  {
    	"apiId": 1111,
        "serviceProviderId": 1,
        "apiName": "Test API Name",
        "description": "Test API Name",
        "endPoint": "/test-api-name/",
        "documentationUrl": "https://example.akamai.com/",
        "accessLevels": [
            "READ-WRITE",
            "READ-ONLY"
        ],
        "hasAccess": false
  },
  {
    	"apiId": 2222,
        "serviceProviderId": 1,
        "apiName": "Example API Name",
        "description": "Example API Name",
        "endPoint": "/example-api-name/",
        "documentationUrl": "https://example2.akamai.com/",
        "accessLevels": [
            "READ-WRITE",
            "READ-ONLY"
        ],
        "hasAccess": false
  },
  {
    	"apiId": 3333,
        "serviceProviderId": 1,
        "apiName": "Best API Name",
        "description": "Best API Name",
        "endPoint": "/best-api-name/",
        "documentationUrl": "https://example3.akamai.com/",
        "accessLevels": [
            "READ-WRITE",
            "READ-ONLY"
        ],
        "hasAccess": false
  }
]`,
			expectedPath: "/identity-management/v3/users/jsmith/allowed-apis?allowAccountSwitch=false",
			expectedResponse: ListAllowedAPIsResponse{
				{
					APIID:             1111,
					ServiceProviderID: 1,
					APIName:           "Test API Name",
					Description:       "Test API Name",
					Endpoint:          "/test-api-name/",
					DocumentationURL:  "https://example.akamai.com/",
					AccessLevels:      []AccessLevel{ReadWriteLevel, ReadOnlyLevel},
					HasAccess:         false,
				},
				{
					APIID:             2222,
					ServiceProviderID: 1,
					APIName:           "Example API Name",
					Description:       "Example API Name",
					Endpoint:          "/example-api-name/",
					DocumentationURL:  "https://example2.akamai.com/",
					AccessLevels:      []AccessLevel{ReadWriteLevel, ReadOnlyLevel},
					HasAccess:         false,
				},
				{
					APIID:             3333,
					ServiceProviderID: 1,
					APIName:           "Best API Name",
					Description:       "Best API Name",
					Endpoint:          "/best-api-name/",
					DocumentationURL:  "https://example3.akamai.com/",
					AccessLevels:      []AccessLevel{ReadWriteLevel, ReadOnlyLevel},
					HasAccess:         false,
				},
			},
		},
		"200 OK with query params": {
			params: ListAllowedAPIsRequest{
				UserName:           "jsmith",
				ClientType:         UserClientType,
				AllowAccountSwitch: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `[
  {
    	"apiId": 1111,
        "serviceProviderId": 1,
        "apiName": "Test API Name",
        "description": "Test API Name",
        "endPoint": "/test-api-name/",
        "documentationUrl": "https://example.akamai.com/",
        "accessLevels": [
            "READ-WRITE",
            "READ-ONLY"
        ],
        "hasAccess": false
  }
]`,
			expectedPath: "/identity-management/v3/users/jsmith/allowed-apis?allowAccountSwitch=true&clientType=USER_CLIENT",
			expectedResponse: ListAllowedAPIsResponse{
				{
					APIID:             1111,
					ServiceProviderID: 1,
					APIName:           "Test API Name",
					Description:       "Test API Name",
					Endpoint:          "/test-api-name/",
					DocumentationURL:  "https://example.akamai.com/",
					AccessLevels:      []AccessLevel{ReadWriteLevel, ReadOnlyLevel},
					HasAccess:         false,
				},
			},
		},
		"500 internal server error": {
			params: ListAllowedAPIsRequest{
				UserName: "jsmith",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
				}`,
			expectedPath: "/identity-management/v3/users/jsmith/allowed-apis?allowAccountSwitch=false",
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
		"missing user name": {
			params: ListAllowedAPIsRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "list allowed APIs: struct validation:\nUserName: cannot be blank")
			},
		},
		"wrong client type": {
			params: ListAllowedAPIsRequest{
				UserName:   "jsmith",
				ClientType: "Test",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "list allowed APIs: struct validation:\nClientType: value 'Test' is invalid. Must be one of: 'CLIENT' or 'USER_CLIENT' or 'SERVICE_ACCOUNT'")
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
			result, err := client.ListAllowedAPIs(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestIAMAccessibleGroups(t *testing.T) {
	tests := map[string]struct {
		params           ListAccessibleGroupsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse ListAccessibleGroupsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListAccessibleGroupsRequest{
				UserName: "jsmith",
			},
			responseStatus: http.StatusOK,
			responseBody: `[
  	{
        "groupId": 1111,
        "groupName": "TestGroupName",
        "roleId": 123123,
        "roleName": "Test Role Name",
        "roleDescription": "Test Role Description",
        "isBlocked": false,
        "subGroups": [
            {
                "groupId": 3333,
                "groupName": "TestSubGroupName",
                "parentGroupId": 1111,
                "subGroups": []
            }
		]
	},
	{
        "groupId": 2222,
        "groupName": "TestGroupName2",
        "roleId": 321321,
        "roleName": "Test Role Name 2",
        "roleDescription": "Test Role Description 2",
        "isBlocked": false,
		"subGroups": []
	}
]`,
			expectedPath: "/identity-management/v3/users/jsmith/group-access",
			expectedResponse: ListAccessibleGroupsResponse{
				{
					GroupID:         1111,
					RoleID:          123123,
					GroupName:       "TestGroupName",
					RoleName:        "Test Role Name",
					IsBlocked:       false,
					RoleDescription: "Test Role Description",
					SubGroups: []AccessibleSubGroup{
						{
							GroupID:       3333,
							GroupName:     "TestSubGroupName",
							ParentGroupID: 1111,
							SubGroups:     []AccessibleSubGroup{},
						},
					},
				},
				{
					GroupID:         2222,
					RoleID:          321321,
					GroupName:       "TestGroupName2",
					RoleName:        "Test Role Name 2",
					IsBlocked:       false,
					RoleDescription: "Test Role Description 2",
					SubGroups:       []AccessibleSubGroup{},
				},
			},
		},
		"500 internal server error": {
			params: ListAccessibleGroupsRequest{
				UserName: "jsmith",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
				}`,
			expectedPath: "/identity-management/v3/users/jsmith/group-access",
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
		"missing user name": {
			params: ListAccessibleGroupsRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "list accessible groups: struct validation:\nUserName: cannot be blank")
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
			result, err := client.ListAccessibleGroups(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
