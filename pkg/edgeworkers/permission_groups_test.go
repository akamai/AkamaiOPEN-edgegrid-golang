package edgeworkers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPermissionGroup(t *testing.T) {
	tests := map[string]struct {
		params           GetPermissionGroupRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PermissionGroup
		withError        error
	}{
		"200 OK - get permission group": {
			params:         GetPermissionGroupRequest{GroupID: "grp_123"},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "groupId": 123,
    "groupName": "Permission Group",
    "capabilities": [
        "VIEW",
        "EDIT",
        "DELETE",
        "VIEW_VERSION",
        "CREATE_VERSION",
        "DELETE_VERSION",
        "VIEW_ACTIVATION",
        "ACTIVATE"
    ]
}
`,
			expectedPath: "/edgeworkers/v1/groups/grp_123",
			expectedResponse: &PermissionGroup{
				ID:   123,
				Name: "Permission Group",
				Capabilities: []string{
					"VIEW",
					"EDIT",
					"DELETE",
					"VIEW_VERSION",
					"CREATE_VERSION",
					"DELETE_VERSION",
					"VIEW_ACTIVATION",
					"ACTIVATE",
				},
			},
		},
		"500 internal server error - get group which does not exist": {
			params:         GetPermissionGroupRequest{GroupID: "grp_1"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
    "title": "Server Error",
    "status": 500,
    "instance": "host_name/edgeworkers/v1/groups/grp_1",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2021-12-06T10:27:11Z"
}`,
			expectedPath: "/edgeworkers/v1/groups/grp_1",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/groups/grp_1",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T10:27:11Z",
			},
		},
		"missing group ID": {
			params:    GetPermissionGroupRequest{},
			withError: ErrStructValidation,
		},
		"403 Forbidden - incorrect credentials": {
			params:         GetPermissionGroupRequest{GroupID: "grp_123"},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/groups/123",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2021-12-06T12:45:09Z"
}`,
			expectedPath: "/edgeworkers/v1/groups/grp_123",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/groups/123",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T12:45:09Z",
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
			result, err := client.GetPermissionGroup(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListPermissionGroups(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListPermissionGroupsResponse
		withError        error
	}{
		"200 OK - list permission groups": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "groups": [
        {
            "groupId": 11111,
            "groupName": "First test group",
            "capabilities": [
                "VIEW",
                "EDIT",
                "DELETE",
                "VIEW_VERSION",
                "CREATE_VERSION",
                "DELETE_VERSION",
                "VIEW_ACTIVATION",
                "ACTIVATE"
            ]
        },
        {
            "groupId": 22222,
            "groupName": "Second test group",
            "capabilities": [
                "VIEW",
                "EDIT",
                "DELETE",
                "VIEW_VERSION",
                "CREATE_VERSION",
                "DELETE_VERSION",
                "VIEW_ACTIVATION",
                "ACTIVATE"
            ]
        },
        {
            "groupId": 33333,
            "groupName": "Third test group",
            "capabilities": [
                "VIEW",
                "EDIT",
                "DELETE",
                "VIEW_VERSION",
                "CREATE_VERSION",
                "DELETE_VERSION",
                "VIEW_ACTIVATION",
                "ACTIVATE"
            ]
        }
    ]
}`,
			expectedPath: "/edgeworkers/v1/groups",
			expectedResponse: &ListPermissionGroupsResponse{[]PermissionGroup{
				{
					ID:   11111,
					Name: "First test group",
					Capabilities: []string{
						"VIEW",
						"EDIT",
						"DELETE",
						"VIEW_VERSION",
						"CREATE_VERSION",
						"DELETE_VERSION",
						"VIEW_ACTIVATION",
						"ACTIVATE",
					},
				},
				{
					ID:   22222,
					Name: "Second test group",
					Capabilities: []string{
						"VIEW",
						"EDIT",
						"DELETE",
						"VIEW_VERSION",
						"CREATE_VERSION",
						"DELETE_VERSION",
						"VIEW_ACTIVATION",
						"ACTIVATE",
					},
				},
				{
					ID:   33333,
					Name: "Third test group",
					Capabilities: []string{
						"VIEW",
						"EDIT",
						"DELETE",
						"VIEW_VERSION",
						"CREATE_VERSION",
						"DELETE_VERSION",
						"VIEW_ACTIVATION",
						"ACTIVATE",
					},
				},
			}},
		},
		"403 Forbidden - incorrect credentials": {
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/groups",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2021-12-06T11:10:42Z"
}`,
			expectedPath: "/edgeworkers/v1/groups",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/groups",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T11:10:42Z",
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
    "title": "Server Error",
    "status": 500,
    "instance": "host_name/edgeworkers/v1/groups",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2021-12-06T10:27:11Z"
}`,
			expectedPath: "/edgeworkers/v1/groups",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/groups",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T10:27:11Z",
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
			result, err := client.ListPermissionGroups(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
