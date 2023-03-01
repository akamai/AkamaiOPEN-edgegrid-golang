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

func TestListGroupsWithinNamespace(t *testing.T) {
	tests := map[string]struct {
		params           ListGroupsWithinNamespaceRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListGroupsWithinNamespaceRequest{
				Network:     NamespaceStagingNetwork,
				NamespaceID: "test_namespace",
			},
			responseStatus: http.StatusOK,
			responseBody: `[
			"test_group_name"
		]`,
			expectedPath:     "/edgekv/v1/networks/staging/namespaces/test_namespace/groups",
			expectedResponse: []string{"test_group_name"},
		},
		"500 internal server error": {
			params: ListGroupsWithinNamespaceRequest{
				Network:     NamespaceStagingNetwork,
				NamespaceID: "test_namespace",
			},
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
			expectedPath: "/edgekv/v1/networks/staging/namespaces/test_namespace/groups",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
					Title:       "Server Error",
					Status:      500,
					Instance:    "host_name/edgeworkers/v1/groups",
					Method:      "GET",
					ServerIP:    "104.81.220.111",
					ClientIP:    "89.64.55.111",
					RequestID:   "a73affa111",
					RequestTime: "2021-12-06T10:27:11Z",
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"NamespaceID - required param not provided": {
			params: ListGroupsWithinNamespaceRequest{
				Network: NamespaceStagingNetwork,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list groups within namespace: struct validation:\nNamespaceID: cannot be blank", err.Error())
			},
		},
		"Network - required param not provided": {
			params: ListGroupsWithinNamespaceRequest{
				NamespaceID: "test_namespace",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list groups within namespace: struct validation:\nNetwork: cannot be blank", err.Error())
			},
		},
		"Network, NamespaceID  - required params not provided": {
			params: ListGroupsWithinNamespaceRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list groups within namespace: struct validation:\nNamespaceID: cannot be blank\nNetwork: cannot be blank", err.Error())
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
			result, err := client.ListGroupsWithinNamespace(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
