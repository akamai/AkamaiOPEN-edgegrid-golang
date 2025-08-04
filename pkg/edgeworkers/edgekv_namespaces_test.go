package edgeworkers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListNamespaces(t *testing.T) {
	tests := map[string]struct {
		params         ListEdgeKVNamespacesRequest
		withError      func(*testing.T, error)
		expectedPath   string
		responseStatus int
		responseBody   string
		expectedResult *ListEdgeKVNamespacesResponse
	}{
		"200 OK - production network": {
			params: ListEdgeKVNamespacesRequest{
				Network: NamespaceProductionNetwork,
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespaces": [
					{
						"namespace": "testNs_1"
					},
					{
						"namespace": "testNs_2"
					},
					{
						"namespace": "testNs_3"
					}
				]
			}`,
			expectedResult: &ListEdgeKVNamespacesResponse{
				Namespaces: []Namespace{
					{
						Name: "testNs_1",
					},
					{
						Name: "testNs_2",
					},
					{
						Name: "testNs_3",
					},
				},
			},
		},
		"200 OK - staging network": {
			params: ListEdgeKVNamespacesRequest{
				Network: NamespaceStagingNetwork,
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespaces": [
					{
						"namespace": "testNs_1"
					},
					{
						"namespace": "testNs_2"
					},
					{
						"namespace": "testNs_3"
					}
				]
			}`,
			expectedResult: &ListEdgeKVNamespacesResponse{
				Namespaces: []Namespace{
					{
						Name: "testNs_1",
					},
					{
						Name: "testNs_2",
					},
					{
						Name: "testNs_3",
					},
				},
			},
		},
		"200 OK - details on": {
			params: ListEdgeKVNamespacesRequest{
				Network: NamespaceProductionNetwork,
				Details: true,
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces?details=on",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespaces": [
					{
						"namespace": "testNs_1",
						"retentionInSeconds": 0,
						"geoLocation": "EU",
						"groupId": 0
					},
					{
						"namespace": "testNs_2",
						"retentionInSeconds": 86400,
						"geoLocation": "JP",
						"groupId": 123
					},
					{
						"namespace": "testNs_3",
						"retentionInSeconds": 315360000,
						"geoLocation": "US",
						"groupId": 234
					}
				]
			}`,
			expectedResult: &ListEdgeKVNamespacesResponse{
				Namespaces: []Namespace{
					{
						Name:        "testNs_1",
						Retention:   ptr.To(0),
						GeoLocation: "EU",
						GroupID:     ptr.To(0),
					},
					{
						Name:        "testNs_2",
						Retention:   ptr.To(86400),
						GeoLocation: "JP",
						GroupID:     ptr.To(123),
					},
					{
						Name:        "testNs_3",
						Retention:   ptr.To(315360000),
						GeoLocation: "US",
						GroupID:     ptr.To(234),
					},
				},
			},
		},
		"missing network": {
			params: ListEdgeKVNamespacesRequest{},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"invalid network": {
			params: ListEdgeKVNamespacesRequest{
				Network: "invalidNetwork",
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: value 'invalidNetwork' is invalid. Must be one of: 'staging' or 'production'", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
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
			result, err := client.ListEdgeKVNamespaces(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetNamespace(t *testing.T) {
	tests := map[string]struct {
		params         GetEdgeKVNamespaceRequest
		withError      func(*testing.T, error)
		expectedPath   string
		responseStatus int
		responseBody   string
		expectedResult *Namespace
	}{
		"200 OK - production": {
			params: GetEdgeKVNamespaceRequest{
				Network: NamespaceProductionNetwork,
				Name:    "testNs",
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces/testNs",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespace": "testNs",
				"retentionInSeconds": 0,
				"geoLocation": "EU",
				"groupId": 0
			}`,
			expectedResult: &Namespace{
				Name:        "testNs",
				Retention:   ptr.To(0),
				GeoLocation: "EU",
				GroupID:     ptr.To(0),
			},
		},
		"200 OK - staging": {
			params: GetEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Name:    "testNs",
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespace": "testNs",
				"retentionInSeconds": 86400,
				"geoLocation": "US",
				"groupId": 0
			}`,
			expectedResult: &Namespace{
				Name:        "testNs",
				Retention:   ptr.To(86400),
				GeoLocation: "US",
				GroupID:     ptr.To(0),
			},
		},
		"400 bad request - namespace does not exist": {
			params: GetEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Name:    "testNs",
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "The requested namespace does not exist.",
					Instance:  "/edgeKV/error-instances/f65424a8-dbea-4799-a2f1-44acc45a121b",
					Status:    http.StatusBadRequest,
					Title:     "Bad Request",
					Type:      "https://learn.akamai.com",
					ErrorCode: "EKV_9000",
					AdditionalDetail: Additional{
						RequestID: "a46f61d2c9539c77",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "The requested namespace does not exist.",
				"errorCode": "EKV_9000",
				"instance": "/edgeKV/error-instances/f65424a8-dbea-4799-a2f1-44acc45a121b",
				"status": 400,
				"title": "Bad Request",
				"type": "https://learn.akamai.com",
				"additionalDetail": {
					"requestId": "a46f61d2c9539c77"
				}
			}`,
		},
		"missing required parameters": {
			params: GetEdgeKVNamespaceRequest{},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Name: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"namespace name too long": {
			params: GetEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Name:    "namespaceNameThatHasMoreThan32Letters",
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Name: the length must be between 1 and 32", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
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
			result, err := client.GetEdgeKVNamespace(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCreateNamespace(t *testing.T) {
	tests := map[string]struct {
		params         CreateEdgeKVNamespaceRequest
		withError      func(*testing.T, error)
		expectedPath   string
		responseStatus int
		responseBody   string
		expectedResult *Namespace
	}{
		"200 OK - production": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceProductionNetwork,
				Namespace: Namespace{
					Name:        "testNs",
					Retention:   ptr.To(0),
					GeoLocation: "EU",
					GroupID:     ptr.To(0),
				},
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespace": "testNs",
				"retentionInSeconds": 0,
				"geoLocation": "EU",
				"groupId": 0
			}`,
			expectedResult: &Namespace{
				Name:        "testNs",
				Retention:   ptr.To(0),
				GeoLocation: "EU",
				GroupID:     ptr.To(0),
			},
		},
		"200 OK - staging": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Namespace: Namespace{
					Name:      "testNs",
					Retention: ptr.To(86400),
					GroupID:   ptr.To(123),
				},
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespace": "testNs",
				"retentionInSeconds": 86400,
				"geoLocation": "US",
				"groupId": 123
			}`,
			expectedResult: &Namespace{
				Name:        "testNs",
				Retention:   ptr.To(86400),
				GeoLocation: "US",
				GroupID:     ptr.To(123),
			},
		},
		"400 bad request - invalid geoLocation for staging network": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Namespace: Namespace{
					Name:        "testNs",
					Retention:   ptr.To(0),
					GeoLocation: "JP",
					GroupID:     ptr.To(0),
				},
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "The staging network only supports US location.",
					Instance:  "/edgeKV/error-instances/afe8f030-30cc-4c8e-9e33-df6b86a9f947",
					Status:    http.StatusBadRequest,
					Title:     "Bad Request",
					Type:      "https://learn.akamai.com",
					ErrorCode: "EKV_2000",
					AdditionalDetail: Additional{
						RequestID: "a46f61d2c9539c77",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "The staging network only supports US location.",
    			"errorCode": "EKV_2000",
    			"instance": "/edgeKV/error-instances/afe8f030-30cc-4c8e-9e33-df6b86a9f947",
    			"status": 400,
    			"title": "Bad Request",
    			"type": "https://learn.akamai.com",
				"additionalDetail": {
					"requestId": "a46f61d2c9539c77"
				}
			}`,
		},
		"400 bad request - geoLocation for production network": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceProductionNetwork,
				Namespace: Namespace{
					Name:        "testNs",
					Retention:   ptr.To(0),
					GeoLocation: "INVALID",
					GroupID:     ptr.To(0),
				},
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "Specified geoLocation not supported. Please specify one of US, EU, JP, GLOBAL",
					Instance:  "/edgeKV/error-instances/d4ae9ce3-7068-4ff8-aef6-3477a0dadbf0",
					Status:    http.StatusBadRequest,
					Title:     "Bad Request",
					Type:      "https://learn.akamai.com",
					ErrorCode: "EKV_2000",
					AdditionalDetail: Additional{
						RequestID: "a46f61d2c9539c77",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "Specified geoLocation not supported. Please specify one of US, EU, JP, GLOBAL",
				"errorCode": "EKV_2000",
				"instance": "/edgeKV/error-instances/d4ae9ce3-7068-4ff8-aef6-3477a0dadbf0",
				"status": 400,
				"title": "Bad Request",
				"type": "https://learn.akamai.com",
				"additionalDetail": {
					"requestId": "a46f61d2c9539c77"
				}
			}`,
		},
		"413 payload too large": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Namespace: Namespace{
					Name:      "testNs",
					Retention: ptr.To(0),
					GroupID:   ptr.To(0),
				},
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "Each account is allowed 20 namespaces. This limit has already been reached.",
					Instance:  "/edgeKV/error-instances/f65424a8-dbea-4799-a2f1-44acc45a121b",
					Status:    http.StatusRequestEntityTooLarge,
					Title:     "Payload Too Large",
					Type:      "https://learn.akamai.com",
					ErrorCode: "EKV_9000",
					AdditionalDetail: Additional{
						RequestID: "a46f61d2c9539c77",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "Each account is allowed 20 namespaces. This limit has already been reached.",
				"errorCode": "EKV_9000",
				"instance": "/edgeKV/error-instances/f65424a8-dbea-4799-a2f1-44acc45a121b",
				"status": 413,
				"title": "Payload Too Large",
				"type": "https://learn.akamai.com",
				"additionalDetail": {
					"requestId": "a46f61d2c9539c77"
				}
			}`,
		},
		"missing required parameters": {
			params: CreateEdgeKVNamespaceRequest{},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Name: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Retention: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "GroupID: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"retention less than 86400": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Namespace: Namespace{
					Name:      "testNs",
					Retention: ptr.To(86399),
					GroupID:   ptr.To(0),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Retention: a non zero value specified for retention period cannot be less than 86400 or more than 315360000", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"retention more than 315360000": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Namespace: Namespace{
					Name:      "testNs",
					Retention: ptr.To(315360001),
					GroupID:   ptr.To(0),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Retention: a non zero value specified for retention period cannot be less than 86400 or more than 315360000", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"namespace name too long": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Namespace: Namespace{
					Name:      "namespaceNameThatHasMoreThan32Letters",
					Retention: ptr.To(0),
					GroupID:   ptr.To(0),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Name: the length must be between 1 and 32", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"groupID less than 0": {
			params: CreateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Namespace: Namespace{
					Name:      "groupIDLessThan0",
					Retention: ptr.To(0),
					GroupID:   ptr.To(-1),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "GroupID: cannot be less than 0", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
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
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateEdgeKVNamespace(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestUpdateNamespace(t *testing.T) {
	tests := map[string]struct {
		params         UpdateEdgeKVNamespaceRequest
		withError      func(*testing.T, error)
		expectedPath   string
		responseStatus int
		responseBody   string
		expectedResult *Namespace
	}{
		"200 OK - production": {
			params: UpdateEdgeKVNamespaceRequest{
				Network: NamespaceProductionNetwork,
				UpdateNamespace: UpdateNamespace{
					Name:      "testNs",
					Retention: ptr.To(86400),
					GroupID:   ptr.To(0),
				},
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces/testNs",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespace": "testNs",
				"retentionInSeconds": 86400,
				"geoLocation": "EU",
				"groupId": 0
			}`,
			expectedResult: &Namespace{
				Name:        "testNs",
				Retention:   ptr.To(86400),
				GeoLocation: "EU",
				GroupID:     ptr.To(0),
			},
		},
		"200 OK - staging": {
			params: UpdateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				UpdateNamespace: UpdateNamespace{
					Name:      "testNs",
					Retention: ptr.To(86400),
					GroupID:   ptr.To(123),
				},
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs",
			responseStatus: http.StatusOK,
			responseBody: `{
				"namespace": "testNs",
				"retentionInSeconds": 86400,
				"geoLocation": "US",
				"groupId": 123
			}`,
			expectedResult: &Namespace{
				Name:        "testNs",
				Retention:   ptr.To(86400),
				GeoLocation: "US",
				GroupID:     ptr.To(123),
			},
		},
		"409 conflict": {
			params: UpdateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				UpdateNamespace: UpdateNamespace{
					Name:      "testNs_2",
					Retention: ptr.To(0),
					GroupID:   ptr.To(0),
				},
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "Cannot update a namespace that does not exist on the network",
					Instance:  "/edgeKV/error-instances/792cc619-2e89-4474-b5e5-e302d2f59e05",
					Status:    http.StatusConflict,
					Title:     "Conflict",
					Type:      "https://learn.akamai.com",
					ErrorCode: "EKV_3000",
					AdditionalDetail: Additional{
						RequestID: "a46f61d2c9539c77",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs_2",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "Cannot update a namespace that does not exist on the network",
				"errorCode": "EKV_3000",
				"instance": "/edgeKV/error-instances/792cc619-2e89-4474-b5e5-e302d2f59e05",
				"status": 409,
				"title": "Conflict",
				"type": "https://learn.akamai.com",
				"additionalDetail": {
					"requestId": "a46f61d2c9539c77"
				}
			}`,
		},
		"missing required parameters": {
			params: UpdateEdgeKVNamespaceRequest{},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Name: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Retention: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "GroupID: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"namespace name too long": {
			params: UpdateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				UpdateNamespace: UpdateNamespace{
					Name:      "namespaceNameThatHasMoreThan32Letters",
					Retention: ptr.To(0),
					GroupID:   ptr.To(0),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Name: the length must be between 1 and 32", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"groupID less than 0": {
			params: UpdateEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				UpdateNamespace: UpdateNamespace{
					Name:      "groupIDLessThan0",
					Retention: ptr.To(0),
					GroupID:   ptr.To(-1),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "GroupID: cannot be less than 0", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
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
			result, err := client.UpdateEdgeKVNamespace(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestDeleteNamespace(t *testing.T) {
	tests := map[string]struct {
		params         DeleteEdgeKVNamespaceRequest
		withError      func(*testing.T, error)
		expectedPath   string
		responseStatus int
		responseBody   string
		expectedResult *DeleteEdgeKVNamespacesResponse
	}{
		"202 Accepted for async deletion": {
			params: DeleteEdgeKVNamespaceRequest{
				Network: NamespaceProductionNetwork,
				Name:    "testNs",
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces/testNs",
			responseStatus: http.StatusAccepted,
			responseBody: `{
				"scheduledDeleteTime": "2025-05-08T14:16:05.350Z"
			}`,
			expectedResult: &DeleteEdgeKVNamespacesResponse{
				ScheduledDeleteTime: ptr.To(test.NewTimeFromString(t, "2025-05-08T14:16:05.350Z")),
			},
		},
		"200 OK for sync deletion": {
			params: DeleteEdgeKVNamespaceRequest{
				Network: NamespaceProductionNetwork,
				Name:    "testNs",
				Sync:    true,
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces/testNs?sync=true",
			responseStatus: http.StatusOK,
			// Yes, the API returns different response format for sync deletion
			responseBody: `{
				"operationPerformed": "DELETED",
				"description": "Namespace 'testNs' was successfully deleted.",
				"id": "1234567"
			}`,
			expectedResult: &DeleteEdgeKVNamespacesResponse{
				ScheduledDeleteTime: nil,
			},
		},
		"400 namespace does not exist": {
			params: DeleteEdgeKVNamespaceRequest{
				Network: NamespaceProductionNetwork,
				Name:    "testNs",
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "That resource does not exist or you are missing the 'DELETE_NAMESPACE' permission for it.",
					Instance:  "/edgekv/error-instances/12345678-abcd-abcd-abcd-1234567890ab",
					Status:    http.StatusBadRequest,
					Title:     "Bad Request",
					Type:      "/edgekv/error-types/edgekv-bad-request",
					ErrorCode: "EKV_1002",
					AdditionalDetail: Additional{
						RequestID: "1234567890abcdef",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/production/namespaces/testNs",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "That resource does not exist or you are missing the 'DELETE_NAMESPACE' permission for it.",
				"errorCode": "EKV_1002",
				"instance": "/edgekv/error-instances/12345678-abcd-abcd-abcd-1234567890ab",
				"status": 400,
				"title": "Bad Request",
				"type": "/edgekv/error-types/edgekv-bad-request",
				"additionalDetail": {
					"requestId": "1234567890abcdef"
				}
			}`,
		},
		"missing required parameters": {
			params: DeleteEdgeKVNamespaceRequest{},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Name: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"namespace name too long": {
			params: DeleteEdgeKVNamespaceRequest{
				Network: NamespaceStagingNetwork,
				Name:    "namespaceNameThatHasMoreThan32Letters",
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Name: the length must be between 1 and 32", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"bad network name": {
			params: DeleteEdgeKVNamespaceRequest{
				Network: NamespaceNetwork("foo"),
				Name:    "testNs",
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: value 'foo' is invalid. Must be one of: "+
					"'staging' or 'production'", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
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
			result, err := client.DeleteEdgeKVNamespace(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetScheduledDeleteTime(t *testing.T) {
	tests := map[string]struct {
		params           GetScheduledDeleteTimeRequest
		withError        func(*testing.T, error)
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResult   *ScheduledDeleteTimeResponse
		retryAfterHeader string
	}{
		"200 OK": {
			params: GetScheduledDeleteTimeRequest{
				Network: NamespaceStagingNetwork,
				Name:    "testNs",
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs/status/scheduled-delete",
			responseStatus: http.StatusOK,
			responseBody: `{
				"scheduledDeleteTime": "2025-06-05T09:58:37.565Z"
			}`,
			retryAfterHeader: "209278",
			expectedResult: &ScheduledDeleteTimeResponse{
				ScheduledDeleteTime: time.Date(2025, time.June, 5, 9, 58, 37, 565000000, time.UTC),
				RetryAfterHeader:    "209278",
			},
		},
		"400 bad request - namespace does not exist": {
			params: GetScheduledDeleteTimeRequest{
				Network: NamespaceStagingNetwork,
				Name:    "testNs_not_exist",
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "That resource does not exist or you are missing the 'DELETE_STAGING_NAMESPACE' permission for it.",
					Instance:  "/edgekv/error-instances/b36d638f-41d4-4582-8e10-ae8d43bc3722",
					Status:    http.StatusBadRequest,
					Title:     "Bad Request",
					Type:      "/edgekv/error-types/edgekv-bad-request",
					ErrorCode: "EKV_1002",
					AdditionalDetail: Additional{
						RequestID: "776a0d334ac1e2d9",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs_not_exist/status/scheduled-delete",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "That resource does not exist or you are missing the 'DELETE_STAGING_NAMESPACE' permission for it.",
				"errorCode": "EKV_1002",
				"instance": "/edgekv/error-instances/b36d638f-41d4-4582-8e10-ae8d43bc3722",
				"status": 400,
				"title": "Bad Request",
				"type": "/edgekv/error-types/edgekv-bad-request",
				"additionalDetail": {
					"requestId": "776a0d334ac1e2d9"
				}
			}`,
		},
		"missing required parameters": {
			params: GetScheduledDeleteTimeRequest{},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Name: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"namespace name too long": {
			params: GetScheduledDeleteTimeRequest{
				Network: NamespaceStagingNetwork,
				Name:    "namespaceNameThatHasMoreThan32Letters",
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Name: the length must be between 1 and 32", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				if len(tc.retryAfterHeader) > 0 {
					w.Header().Set("Retry-After", tc.retryAfterHeader)
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetNamespaceScheduledDeleteTime(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestRescheduleNamespaceDelete(t *testing.T) {
	tests := map[string]struct {
		params           RescheduleNamespaceDeleteRequest
		withError        func(*testing.T, error)
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResult   *RescheduleNamespaceDeleteResponse
		retryAfterHeader string
	}{
		"200 reschedule namespace delete": {
			params: RescheduleNamespaceDeleteRequest{
				Network: NamespaceStagingNetwork,
				Name:    "testNs",
				Body: &ScheduledDeleteTimeRequest{
					test.NewTimeFromString(t, "2025-06-05T19:58:37.565Z"),
				},
			},
			expectedPath:     "/edgekv/v1/networks/staging/namespaces/testNs/status/scheduled-delete",
			responseStatus:   http.StatusOK,
			retryAfterHeader: "209278",
			responseBody: `{
				"scheduledDeleteTime": "2025-06-05T19:58:37.565Z"
			}`,
			expectedResult: &RescheduleNamespaceDeleteResponse{
				ScheduledDeleteTime: ScheduledDeleteTimeResponse{
					ScheduledDeleteTime: test.NewTimeFromString(t, "2025-06-05T19:58:37.565Z"),
				},
				RetryAfterHeader: "209278",
			},
		},
		"400 bad request - namespace does not exist": {
			params: RescheduleNamespaceDeleteRequest{
				Network: NamespaceStagingNetwork,
				Name:    "testNs_not_exist",
				Body: &ScheduledDeleteTimeRequest{
					test.NewTimeFromString(t, "2025-06-05T19:58:37.565Z"),
				},
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "That resource does not exist or you are missing the 'DELETE_STAGING_NAMESPACE' permission for it.",
					Instance:  "/edgekv/error-instances/f5756eb4-b670-44bc-88d7-d7bc7ab45402",
					Status:    http.StatusBadRequest,
					Title:     "Bad Request",
					Type:      "/edgekv/error-types/edgekv-bad-request",
					ErrorCode: "EKV_1002",
					AdditionalDetail: Additional{
						RequestID: "830d86d914fe7e0b",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs_not_exist/status/scheduled-delete",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "That resource does not exist or you are missing the 'DELETE_STAGING_NAMESPACE' permission for it.",
				"errorCode": "EKV_1002",
				"instance": "/edgekv/error-instances/f5756eb4-b670-44bc-88d7-d7bc7ab45402",
				"status": 400,
				"title": "Bad Request",
				"type": "/edgekv/error-types/edgekv-bad-request",
				"additionalDetail": {
					"requestId": "830d86d914fe7e0b"
				}
			}`,
		},
		"missing required parameters": {
			params: RescheduleNamespaceDeleteRequest{},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Name: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Body: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"namespace name too long": {
			params: RescheduleNamespaceDeleteRequest{
				Network: NamespaceStagingNetwork,
				Name:    "namespaceNameThatHasMoreThan32Letters",
				Body: &ScheduledDeleteTimeRequest{
					test.NewTimeFromString(t, "2025-06-05T19:58:37.565Z"),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Name: the length must be between 1 and 32", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if len(tc.retryAfterHeader) > 0 {
					w.Header().Set("Retry-After", tc.retryAfterHeader)
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.RescheduleNamespaceDelete(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCancelScheduledNamespaceDelete(t *testing.T) {
	tests := map[string]struct {
		params         CancelScheduledNamespaceDeleteRequest
		withError      func(*testing.T, error)
		expectedPath   string
		responseStatus int
		responseBody   string
	}{
		"204 cancel scheduled namespace delete": {
			params: CancelScheduledNamespaceDeleteRequest{
				Network: NamespaceStagingNetwork,
				Name:    "testNs",
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs/status/scheduled-delete",
			responseStatus: http.StatusNoContent,
		},
		"400 namespace does not exist": {
			params: CancelScheduledNamespaceDeleteRequest{
				Network: NamespaceStagingNetwork,
				Name:    "testNs_not_exist",
			},
			withError: func(t *testing.T, err error) {
				expected := &Error{
					Detail:    "That resource does not exist or you are missing the 'DELETE_STAGING_NAMESPACE' permission for it.",
					Instance:  "/edgekv/error-instances/cce5131b-a821-406a-b864-7aa680ba55bc",
					Status:    http.StatusBadRequest,
					Title:     "Bad Request",
					Type:      "/edgekv/error-types/edgekv-bad-request",
					ErrorCode: "EKV_1002",
					AdditionalDetail: Additional{
						RequestID: "f0da6846b5f903e6",
					},
				}
				assert.True(t, errors.Is(err, expected), "want: %s; got: %s", expected, err)
			},
			expectedPath:   "/edgekv/v1/networks/staging/namespaces/testNs_not_exist/status/scheduled-delete",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"detail": "That resource does not exist or you are missing the 'DELETE_STAGING_NAMESPACE' permission for it.",
				"errorCode": "EKV_1002",
				"instance": "/edgekv/error-instances/cce5131b-a821-406a-b864-7aa680ba55bc",
				"status": 400,
				"title": "Bad Request",
				"type": "/edgekv/error-types/edgekv-bad-request",
				"additionalDetail": {
					"requestId": "f0da6846b5f903e6"
				}
			}`,
		},
		"missing required parameters": {
			params: CancelScheduledNamespaceDeleteRequest{},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Name: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"namespace name too long": {
			params: CancelScheduledNamespaceDeleteRequest{
				Network: NamespaceStagingNetwork,
				Name:    "namespaceNameThatHasMoreThan32Letters",
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Name: the length must be between 1 and 32", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
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
			err := client.CancelScheduledNamespaceDelete(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
