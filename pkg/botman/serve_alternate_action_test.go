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

// Test Get ServeAlternateAction List
func TestBotman_GetServeAlternateActionList(t *testing.T) {

	tests := map[string]struct {
		params           GetServeAlternateActionListRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetServeAlternateActionListResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetServeAlternateActionListRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"serveAlternateActions": [
		{"actionId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"actionId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"actionId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"actionId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/response-actions/serve-alternate-actions",
			expectedResponse: &GetServeAlternateActionListResponse{
				ServeAlternateActions: []map[string]interface{}{
					{"actionId": "b85e3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
					{"actionId": "69acad64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
					{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
					{"actionId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
					{"actionId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
				},
			},
		},
		"200 OK One Record": {
			params: GetServeAlternateActionListRequest{
				ConfigID: 43253,
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"serveAlternateActions":[
		{"actionId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"actionId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"actionId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"actionId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/response-actions/serve-alternate-actions",
			expectedResponse: &GetServeAlternateActionListResponse{
				ServeAlternateActions: []map[string]interface{}{
					{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				},
			},
		},
		"500 internal server error": {
			params: GetServeAlternateActionListRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching data",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/response-actions/serve-alternate-actions",
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
			params: GetServeAlternateActionListRequest{
				Version: 15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetServeAlternateActionListRequest{
				ConfigID: 43253,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
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
			result, err := client.GetServeAlternateActionList(
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

// Test Get ServeAlternateAction
func TestBotman_GetServeAlternateAction(t *testing.T) {
	tests := map[string]struct {
		params           GetServeAlternateActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/response-actions/serve-alternate-actions/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			expectedResponse: map[string]interface{}{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
		},
		"500 internal server error": {
			params: GetServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/response-actions/serve-alternate-actions/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: GetServeAlternateActionRequest{
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetServeAlternateActionRequest{
				ConfigID: 43253,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing ActionID": {
			params: GetServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ActionID")
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
			result, err := client.GetServeAlternateAction(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create ServeAlternateAction
func TestBotman_CreateServeAlternateAction(t *testing.T) {

	tests := map[string]struct {
		params           CreateServeAlternateActionRequest
		prop             *CreateServeAlternateActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateServeAlternateActionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusCreated,
			responseBody:     `{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/15/response-actions/serve-alternate-actions",
		},
		"500 internal server error": {
			params: CreateServeAlternateActionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/response-actions/serve-alternate-actions",
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
			params: CreateServeAlternateActionRequest{
				Version:     15,
				JsonPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: CreateServeAlternateActionRequest{
				ConfigID:    43253,
				JsonPayload: json.RawMessage(`{"testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing JsonPayload": {
			params: CreateServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
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
			result, err := client.CreateServeAlternateAction(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update ServeAlternateAction
func TestBotman_UpdateServeAlternateAction(t *testing.T) {
	tests := map[string]struct {
		params           UpdateServeAlternateActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateServeAlternateActionRequest{
				ConfigID:    43253,
				Version:     10,
				ActionID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`,
			expectedResponse: map[string]interface{}{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			expectedPath:     "/appsec/v1/configs/43253/versions/10/response-actions/serve-alternate-actions/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: UpdateServeAlternateActionRequest{
				ConfigID:    43253,
				Version:     10,
				ActionID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/response-actions/serve-alternate-actions/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: UpdateServeAlternateActionRequest{
				Version:     15,
				ActionID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateServeAlternateActionRequest{
				ConfigID:    43253,
				ActionID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(`{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing JsonPayload": {
			params: UpdateServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "JsonPayload")
			},
		},
		"Missing ActionID": {
			params: UpdateServeAlternateActionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: json.RawMessage(`{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`),
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ActionID")
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
			result, err := client.UpdateServeAlternateAction(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Remove ServeAlternateAction
func TestBotman_RemoveServeAlternateAction(t *testing.T) {
	tests := map[string]struct {
		params           RemoveServeAlternateActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]interface{}
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: RemoveServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  10,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/appsec/v1/configs/43253/versions/10/response-actions/serve-alternate-actions/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
		"500 internal server error": {
			params: RemoveServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  10,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/10/response-actions/serve-alternate-actions/cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
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
			params: RemoveServeAlternateActionRequest{
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: RemoveServeAlternateActionRequest{
				ConfigID: 43253,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing ActionID": {
			params: RemoveServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ActionID")
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
			err := client.RemoveServeAlternateAction(session.ContextWithOptions(context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
