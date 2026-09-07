package accountprotection

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListProtectedOperations(t *testing.T) {

	tests := map[string]struct {
		params           ListProtectedOperationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListProtectedOperationsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListProtectedOperationsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusOK,
			responseBody: `
					{
						"metadata": {
							"configId": 43253,
							"configVersion": 15,
							"securityPolicyId": "AAAA_81230"
						},
						"operations": [
							{"operationId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
							{"operationId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
							{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
							{"operationId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
							{"operationId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
						]
					}
				`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/account-protection",
			expectedResponse: &ListProtectedOperationsResponse{
				Metadata: Metadata{
					ConfigID:         43253,
					ConfigVersion:    15,
					SecurityPolicyID: "AAAA_81230",
				},
				Operations: []map[string]interface{}{
					{"operationId": "b85e3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
					{"operationId": "69acad64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
					{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
					{"operationId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
					{"operationId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
				},
			},
		},
		"500 internal server error": {
			params: ListProtectedOperationsRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/account-protection",
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
			params: ListProtectedOperationsRequest{
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
			params: ListProtectedOperationsRequest{
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
			params: ListProtectedOperationsRequest{
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListProtectedOperations(
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

func Test_GetProtectedOperationByID(t *testing.T) {

	tests := map[string]struct {
		params           GetProtectedOperationByIDRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListProtectedOperationsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetProtectedOperationByIDRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusOK,
			responseBody: `
					{
						"metadata": {
							"configId": 43253,
							"configVersion": 15,
							"securityPolicyId": "AAAA_81230"
						},
						"operations": [
							{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}
						]
					}
				`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/account-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			expectedResponse: &ListProtectedOperationsResponse{
				Metadata: Metadata{
					ConfigID:         43253,
					ConfigVersion:    15,
					SecurityPolicyID: "AAAA_81230",
				},
				Operations: []map[string]interface{}{
					{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				},
			},
		},
		"500 internal server error": {
			params: GetProtectedOperationByIDRequest{
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
					"detail": "Error fetching data",
					"status": 500
				}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/account-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: GetProtectedOperationByIDRequest{
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
			params: GetProtectedOperationByIDRequest{
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
			params: GetProtectedOperationByIDRequest{
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetProtectedOperationByID(
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

func Test_CreateProtectedOperations(t *testing.T) {

	tests := map[string]struct {
		params              CreateProtectedOperationsRequest
		expectedRequestBody string
		prop                *CreateProtectedOperationsRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *ListProtectedOperationsResponse
		withError           func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateProtectedOperationsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"operations": [{"testKey":"testValue3"}]}`),
			},
			expectedRequestBody: `{"operations": [{"testKey":"testValue3"}]}`,
			responseStatus:      http.StatusCreated,
			responseBody: `
				{
			        "metadata": {
						"configId": 43253,
						"configVersion": 15,
						"securityPolicyId": "AAAA_81230"
					},
					"operations":[
						{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}
					]
				}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/account-protection",
			expectedResponse: &ListProtectedOperationsResponse{
				Metadata: Metadata{
					ConfigID:         43253,
					ConfigVersion:    15,
					SecurityPolicyID: "AAAA_81230",
				},
				Operations: []map[string]interface{}{
					{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				},
			},
		},
		"500 internal server error": {
			params: CreateProtectedOperationsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"operations": [{"testKey":"testValue3"}]}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/transactional-endpoints/account-protection",
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
			params: CreateProtectedOperationsRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"operations": [{"testKey":"testValue3"}]}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: CreateProtectedOperationsRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(`{"operations": [{"testKey":"testValue3"}]}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing JsonPayload": {
			params: CreateProtectedOperationsRequest{
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateProtectedOperations(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func Test_UpdateProtectedOperation(t *testing.T) {
	tests := map[string]struct {
		params              UpdateProtectedOperationRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    map[string]interface{}
		withError           func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateProtectedOperationRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(`{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus:      http.StatusOK,
			responseBody:        `{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedRequestBody: `{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse:    map[string]interface{}{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:        "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/transactional-endpoints/account-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: UpdateProtectedOperationRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/transactional-endpoints/account-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: UpdateProtectedOperationRequest{
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
			params: UpdateProtectedOperationRequest{
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
			params: UpdateProtectedOperationRequest{
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
			params: UpdateProtectedOperationRequest{
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
			params: UpdateProtectedOperationRequest{
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateProtectedOperation(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func Test_RemoveProtectedOperation(t *testing.T) {
	tests := map[string]struct {
		params           RemoveProtectedOperationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: RemoveProtectedOperationRequest{
				ConfigID:         43253,
				Version:          10,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/transactional-endpoints/account-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: RemoveProtectedOperationRequest{
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
			expectedPath: "/appsec/v1/configs/43253/versions/10/security-policies/AAAA_81230/transactional-endpoints/account-protection/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: RemoveProtectedOperationRequest{
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
			params: RemoveProtectedOperationRequest{
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
			params: RemoveProtectedOperationRequest{
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
			params: RemoveProtectedOperationRequest{
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			err := client.RemoveProtectedOperation(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
