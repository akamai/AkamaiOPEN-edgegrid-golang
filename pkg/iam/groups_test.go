package iam

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/internal/test"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestIAM_CreateGroup(t *testing.T) {
	tests := map[string]struct {
		params              GroupRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *Group
		withError           error
	}{
		"201 OK": {
			params: GroupRequest{
				GroupID:   12345,
				GroupName: "Test Group",
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"groupId": 98765,
	"groupName": "Test Group",
	"parentGroupId": 12345,
	"createdDate": "2012-04-28T00:00:00.000Z",
	"createdBy": "johndoe",
	"modifiedDate": "2012-04-28T00:00:00.000Z",
	"modifiedBy": "johndoe"
}`,
			expectedPath:        "/identity-management/v3/user-admin/groups/12345",
			expectedRequestBody: `{"groupName":"Test Group"}`,
			expectedResponse: &Group{
				GroupID:       98765,
				GroupName:     "Test Group",
				ParentGroupID: 12345,
				CreatedDate:   test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
				CreatedBy:     "johndoe",
				ModifiedDate:  test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
				ModifiedBy:    "johndoe",
			},
		},
		"500 internal server error": {
			params: GroupRequest{
				GroupID:   12345,
				GroupName: "Test Group",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/groups/12345",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"missing group id": {
			params: GroupRequest{
				GroupName: "Test Group",
			},
			withError: ErrStructValidation,
		},
		"missing group name": {
			params: GroupRequest{
				GroupID: 12345,
			},
			withError: ErrStructValidation,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)

				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateGroup(context.Background(), tc.params)
			if tc.withError != nil {
				assert.True(t, errors.Is(err, tc.withError), "want: %s; got: %s", tc.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_MoveGroup(t *testing.T) {
	tests := map[string]struct {
		params              MoveGroupRequest
		expectedRequestBody string
		responseStatus      int
		withError           error
		responseBody        string
	}{
		"204 ok": {
			responseStatus:      http.StatusNoContent,
			expectedRequestBody: `{"sourceGroupId":1,"destinationGroupId":1}`,
			params:              MoveGroupRequest{DestinationGroupID: 1, SourceGroupID: 1},
		},
		"500 internal server error": {
			responseStatus:      http.StatusInternalServerError,
			params:              MoveGroupRequest{DestinationGroupID: 1, SourceGroupID: 1},
			expectedRequestBody: `{"sourceGroupId":1,"destinationGroupId":1}`,
			withError:           ErrMoveGroup,
		},
		"validation error": {
			params:    MoveGroupRequest{},
			withError: ErrStructValidation,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/identity-management/v3/user-admin/groups/move", r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedRequestBody, string(body))
				w.WriteHeader(tc.responseStatus)
				_, err = w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.MoveGroup(context.Background(), tc.params)
			if tc.withError != nil {
				assert.True(t, errors.Is(err, tc.withError), "want: %s; got: %s", tc.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_GetGroup(t *testing.T) {
	tests := map[string]struct {
		params           GetGroupRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Group
		withError        error
	}{
		"200 OK with actions": {
			params: GetGroupRequest{
				GroupID: 12345,
				Actions: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"groupId": 12345,
	"groupName": "Top Level group",
	"createdDate": "2012-04-28T00:00:00.000Z",
	"createdBy": "johndoe",
	"modifiedDate": "2012-04-28T00:00:00.000Z",
	"modifiedBy": "johndoe",
	"actions": {
		"edit": true,
		"delete": true
	}
}`,
			expectedPath: "/identity-management/v3/user-admin/groups/12345?actions=true",
			expectedResponse: &Group{
				GroupID:      12345,
				GroupName:    "Top Level group",
				CreatedDate:  test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
				CreatedBy:    "johndoe",
				ModifiedDate: test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
				ModifiedBy:   "johndoe",
				Actions: &GroupActions{
					Edit:   true,
					Delete: true,
				},
			},
		},
		"200 OK with no actions": {
			params: GetGroupRequest{
				GroupID: 12345,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"groupId": 12345,
	"groupName": "Top Level group",
	"createdDate": "2012-04-28T00:00:00.000Z",
	"createdBy": "johndoe",
	"modifiedDate": "2012-04-28T00:00:00.000Z",
	"modifiedBy": "johndoe"
}`,
			expectedPath: "/identity-management/v3/user-admin/groups/12345?actions=false",
			expectedResponse: &Group{
				GroupID:      12345,
				GroupName:    "Top Level group",
				CreatedDate:  test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
				CreatedBy:    "johndoe",
				ModifiedDate: test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
				ModifiedBy:   "johndoe",
			},
		},
		"500 internal server error": {
			params: GetGroupRequest{
				GroupID: 12345,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/groups/12345?actions=false",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"missing group id": {
			params: GetGroupRequest{
				Actions: true,
			},
			withError: ErrStructValidation,
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
			result, err := client.GetGroup(context.Background(), tc.params)
			if tc.withError != nil {
				assert.True(t, errors.Is(err, tc.withError), "want: %s; got: %s", tc.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_ListAffectedUsers(t *testing.T) {
	tests := map[string]struct {
		params           ListAffectedUsersRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []GroupUser
		withError        error
	}{
		"200 OK": {
			params: ListAffectedUsersRequest{
				DestinationGroupID: 12345,
				SourceGroupID:      12344,
				UserType:           GainAccessUsers,
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "uiIdentityId": "test-identity",
        "firstName": "jhon",
        "lastName": "doe",
        "accountId": "test-account",
        "email": "john.doe@mycompany.com",
        "uiUserName": "john.doe@mycompany.com",
        "lastLoginDate": "2022-02-22T17:06:50.000Z"
    }
]`,
			expectedPath: "/identity-management/v3/user-admin/groups/move/12344/12345/affected-users?userType=gainAccess",
			expectedResponse: []GroupUser{
				{
					IdentityID:    "test-identity",
					FirstName:     "jhon",
					LastName:      "doe",
					AccountID:     "test-account",
					Email:         "john.doe@mycompany.com",
					UserName:      "john.doe@mycompany.com",
					LastLoginDate: test.NewTimeFromString(t, "2022-02-22T17:06:50.000Z"),
				},
			},
		},
		"500 internal server error": {
			params: ListAffectedUsersRequest{
				DestinationGroupID: 12345,
				SourceGroupID:      12344,
				UserType:           GainAccessUsers,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/groups/move/12344/12345/affected-users?userType=gainAccess",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"missing destination group id": {
			params: ListAffectedUsersRequest{
				SourceGroupID: 12344,
				UserType:      GainAccessUsers,
			},
			withError: ErrStructValidation,
		},
		"missing source group id": {
			params: ListAffectedUsersRequest{
				DestinationGroupID: 12345,
				UserType:           GainAccessUsers,
			},
			withError: ErrStructValidation,
		},
		"invalid user type": {
			params: ListAffectedUsersRequest{
				DestinationGroupID: 12345,
				SourceGroupID:      12344,
				UserType:           "different",
			},
			withError: ErrStructValidation,
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
			result, err := client.ListAffectedUsers(context.Background(), tc.params)
			if tc.withError != nil {
				assert.True(t, errors.Is(err, tc.withError), "want: %s; got: %s", tc.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_ListGroups(t *testing.T) {
	tests := map[string]struct {
		params           ListGroupsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Group
		withError        func(*testing.T, error)
	}{
		"200 OK with actions": {
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
						"delete": true
					}
				}
			]`,
			expectedPath: "/identity-management/v3/user-admin/groups?actions=true",
			expectedResponse: []Group{
				{
					GroupID:      12345,
					GroupName:    "Top Level group",
					CreatedDate:  test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
					CreatedBy:    "johndoe",
					ModifiedDate: test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
					ModifiedBy:   "johndoe",
					Actions: &GroupActions{
						Edit:   true,
						Delete: true,
					},
				},
			},
		},
		"200 OK with no actions": {
			params: ListGroupsRequest{
				Actions: false,
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
					"modifiedBy": "johndoe"
				}
			]`,
			expectedPath: "/identity-management/v3/user-admin/groups?actions=false",
			expectedResponse: []Group{
				{
					GroupID:      12345,
					GroupName:    "Top Level group",
					CreatedDate:  test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
					CreatedBy:    "johndoe",
					ModifiedDate: test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
					ModifiedBy:   "johndoe",
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
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/groups?actions=true",
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
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListGroups(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_RemoveGroup(t *testing.T) {
	tests := map[string]struct {
		params         RemoveGroupRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 deleted": {
			params: RemoveGroupRequest{
				GroupID: 12345,
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v3/user-admin/groups/12345",
		},
		"500 internal server error": {
			params: RemoveGroupRequest{
				GroupID: 12345,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/groups/12345",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"403 not authorised": {
			params: RemoveGroupRequest{
				GroupID: 12345,
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "instance": "",
    "httpStatus": 403,
    "detail": "Not Authorized to perform this action",
    "title": "Forbidden",
    "type": "/useradmin-api/error-types/1001"
}`,
			expectedPath: "/identity-management/v3/user-admin/groups/12345",
			withError: &Error{
				Instance:   "",
				StatusCode: http.StatusForbidden,
				Detail:     "Not Authorized to perform this action",
				Title:      "Forbidden",
				Type:       "/useradmin-api/error-types/1001",
				HTTPStatus: http.StatusForbidden,
			},
		},
		"missing group id": {
			params:    RemoveGroupRequest{},
			withError: ErrStructValidation,
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
			err := client.RemoveGroup(context.Background(), tc.params)
			if tc.withError != nil {
				assert.True(t, errors.Is(err, tc.withError), "want: %s; got: %s", tc.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestIAM_UpdateGroupName(t *testing.T) {
	tests := map[string]struct {
		params              GroupRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *Group
		withError           error
	}{
		"200 updated": {
			params: GroupRequest{
				GroupID:   12345,
				GroupName: "New Group Name",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"groupId": 12345,
	"groupName": "New Group Name",
	"parentGroupId": 12344,
	"createdDate": "2012-04-28T00:00:00.000Z",
	"createdBy": "johndoe",
	"modifiedDate": "2012-04-28T00:00:00.000Z",
	"modifiedBy": "johndoe"
}`,
			expectedPath:        "/identity-management/v3/user-admin/groups/12345",
			expectedRequestBody: `{"groupName":"New Group Name"}`,
			expectedResponse: &Group{
				GroupID:       12345,
				GroupName:     "New Group Name",
				ParentGroupID: 12344,
				CreatedDate:   test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
				CreatedBy:     "johndoe",
				ModifiedDate:  test.NewTimeFromString(t, "2012-04-28T00:00:00.000Z"),
				ModifiedBy:    "johndoe",
			},
		},
		"500 internal server error": {
			params: GroupRequest{
				GroupID:   12345,
				GroupName: "New Group Name",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/groups/12345",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error making request",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"missing group id": {
			params: GroupRequest{
				GroupName: "New Group Name",
			},
			withError: ErrStructValidation,
		},
		"missing group name": {
			params: GroupRequest{
				GroupID: 12345,
			},
			withError: ErrStructValidation,
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

				if len(tc.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateGroupName(context.Background(), tc.params)
			if tc.withError != nil {
				assert.True(t, errors.Is(err, tc.withError), "want: %s; got: %s", tc.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
