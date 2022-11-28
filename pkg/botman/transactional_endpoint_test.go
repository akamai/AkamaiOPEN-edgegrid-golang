package botman

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Get TransactionalEndpoint List
func TestBotman_GetTransactionalEndpointList(t *testing.T) {

	tests := map[string]struct {
		params           GetTransactionalEndpointListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetTransactionalEndpointListResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetTransactionalEndpointListRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"operations": [
		{"operationId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"operationId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"operationId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"operationId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/bot-protection",
			expectedResponse: &GetTransactionalEndpointListResponse{
				Operations: []map[string]interface{}{
					{"operationId": "b85e3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
					{"operationId": "69acad64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
					{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
					{"operationId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
					{"operationId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
				},
			},
		},
		"200 OK One Record": {
			params: GetTransactionalEndpointListRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"operations":[
		{"operationId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"operationId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"operationId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"operationId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/bot-protection",
			expectedResponse: &GetTransactionalEndpointListResponse{
				Operations: []map[string]interface{}{
					{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				},
			},
		},
		"500 internal server error": {
			params: GetTransactionalEndpointListRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching data",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/bot-protection",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: GetTransactionalEndpointListRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetTransactionalEndpointListRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: GetTransactionalEndpointListRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
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
			result, err := client.GetTransactionalEndpointList(
				session.ContextWithOptions(
					context.Background(),
				),
				test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Get TransactionalEndpoint
func TestBotman_GetTransactionalEndpoint(t *testing.T) {
	tests := map[string]struct {
		params           GetTransactionalEndpointRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/bot-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			expectedResponse: map[string]interface{}{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
		},
		"500 internal server error": {
			params: GetTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/bot-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: GetTransactionalEndpointRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetTransactionalEndpointRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: GetTransactionalEndpointRequest{
				ConfigID:    43253,
				Version:     15,
				OperationID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing OperationID": {
			params: GetTransactionalEndpointRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				Version:          15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "OperationID")
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
			result, err := client.GetTransactionalEndpoint(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create TransactionalEndpoint
func TestBotman_CreateTransactionalEndpoint(t *testing.T) {

	tests := map[string]struct {
		params           CreateTransactionalEndpointRequest
		prop             *CreateTransactionalEndpointRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusCreated,
			responseBody:     `{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/bot-protection",
		},
		"500 internal server error": {
			params: CreateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/bot-protection",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: CreateTransactionalEndpointRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: CreateTransactionalEndpointRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing JsonPayload": {
			params: CreateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "JsonPayload")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateTransactionalEndpoint(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update TransactionalEndpoint
func TestBotman_UpdateTransactionalEndpoint(t *testing.T) {
	tests := map[string]struct {
		params           UpdateTransactionalEndpointRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(`{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/transactional-endpoints/bot-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: UpdateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(`{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/transactional-endpoints/bot-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating zone",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: UpdateTransactionalEndpointRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(`{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateTransactionalEndpointRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(`{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: UpdateTransactionalEndpointRequest{
				ConfigID:    43253,
				Version:     15,
				OperationID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing JsonPayload": {
			params: UpdateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "JsonPayload")
			},
		},
		"Missing OperationID": {
			params: UpdateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "OperationID")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateTransactionalEndpoint(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Remove TransactionalEndpoint
func TestBotman_RemoveTransactionalEndpoint(t *testing.T) {
	tests := map[string]struct {
		params           RemoveTransactionalEndpointRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: RemoveTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/transactional-endpoints/bot-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: RemoveTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/transactional-endpoints/bot-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error deleting match target",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: RemoveTransactionalEndpointRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: RemoveTransactionalEndpointRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: RemoveTransactionalEndpointRequest{
				ConfigID:    43253,
				Version:     15,
				OperationID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing OperationID": {
			params: RemoveTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "OperationID")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.RemoveTransactionalEndpoint(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
