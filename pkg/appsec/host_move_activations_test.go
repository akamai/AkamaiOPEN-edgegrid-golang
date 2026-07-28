package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_GetHostMoveValidation(t *testing.T) {

	result := GetHostMoveValidationResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestHostMoveActivations/HostMoveValidation.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetHostMoveValidationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetHostMoveValidationResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 1,
				Network:       NetworkStaging,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/1/network/STAGING/host-move-validation",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 1,
				Network:       NetworkStaging,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching host move validation",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/1/network/STAGING/host-move-validation",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching host move validation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"missing ConfigID": {
			params: GetHostMoveValidationRequest{
				ConfigVersion: 1,
				Network:       NetworkStaging,
			},
			withError: ErrStructValidation,
		},
		"missing ConfigVersion": {
			params: GetHostMoveValidationRequest{
				ConfigID: 43253,
				Network:  NetworkStaging,
			},
			withError: ErrStructValidation,
		},
		"missing Network": {
			params: GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 1,
			},
			withError: ErrStructValidation,
		},
		"incorrect Network": {
			params: GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 1,
				Network:       "INCORRECT",
			},
			withError: ErrStructValidation,
		},
		"400 bad request - empty hosts to move": {
			params: GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 1,
				Network:       NetworkStaging,
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "bad_request",
				"title": "Bad Request",
				"detail": "No hosts require moving for this configuration",
				"status": 400
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/1/network/STAGING/host-move-validation",
			withError: &Error{
				Type:       "bad_request",
				Title:      "Bad Request",
				Detail:     "No hosts require moving for this configuration",
				StatusCode: http.StatusBadRequest,
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
			result, err := client.GetHostMoveValidation(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_CreateActivationsWithHostMove(t *testing.T) {

	result := CreateActivationsWithHostMoveResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestHostMoveActivations/CreateActivationsWithHostMove.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params              CreateActivationsWithHostMoveRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateActivationsWithHostMoveResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{
		"200 OK": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID:                         43253,
				ConfigVersion:                    1,
				Action:                           "ACTIVATE",
				Network:                          "STAGING",
				Note:                             "Test activation with host move",
				NotificationEmails:               []string{"test@example.com"},
				AcknowledgedInvalidHosts:         []string{},
				AcknowledgedInvalidHostsByConfig: []AcknowledgedInvalidHostsByConfig{},
				HostsToMove:                      []HostToMove{},
			},
			expectedRequestBody: `{"action":"ACTIVATE","network":"STAGING","note":"Test activation with host move","notificationEmails":["test@example.com"],"acknowledgedInvalidHosts":[],"acknowledgedInvalidHostsByConfig":[],"hostsToMove":[],"supportId":""}`,
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/1/activations-with-host-move",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID:                         43253,
				ConfigVersion:                    1,
				Action:                           "ACTIVATE",
				Network:                          "STAGING",
				Note:                             "Test activation with host move",
				NotificationEmails:               []string{"test@example.com"},
				AcknowledgedInvalidHosts:         []string{},
				AcknowledgedInvalidHostsByConfig: []AcknowledgedInvalidHostsByConfig{},
				HostsToMove:                      []HostToMove{},
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating activation with host move",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/1/activations-with-host-move",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating activation with host move",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"400 bad request - invalid action": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID:                         43253,
				ConfigVersion:                    1,
				Action:                           "INVALID_ACTION",
				Network:                          "STAGING",
				Note:                             "Test activation with invalid action",
				NotificationEmails:               []string{"test@example.com"},
				AcknowledgedInvalidHosts:         []string{},
				AcknowledgedInvalidHostsByConfig: []AcknowledgedInvalidHostsByConfig{},
				HostsToMove:                      []HostToMove{},
			},
			headers:        http.Header{},
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"type": "bad_request",
				"title": "Bad Request",
				"detail": "Invalid action specified. Valid actions are ACTIVATE, DEACTIVATE",
				"status": 400
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/1/activations-with-host-move",
			withError: &Error{
				Type:       "bad_request",
				Title:      "Bad Request",
				Detail:     "Invalid action specified. Valid actions are ACTIVATE, DEACTIVATE",
				StatusCode: http.StatusBadRequest,
			},
		},
		"409 conflict - activation in progress": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID:                         43253,
				ConfigVersion:                    1,
				Action:                           "ACTIVATE",
				Network:                          "STAGING",
				Note:                             "Test activation conflict",
				NotificationEmails:               []string{"test@example.com"},
				AcknowledgedInvalidHosts:         []string{},
				AcknowledgedInvalidHostsByConfig: []AcknowledgedInvalidHostsByConfig{},
				HostsToMove:                      []HostToMove{},
			},
			headers:        http.Header{},
			responseStatus: http.StatusConflict,
			responseBody: `{
				"type": "conflict",
				"title": "Conflict",
				"detail": "Another activation is already in progress for this configuration",
				"status": 409
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/1/activations-with-host-move",
			withError: &Error{
				Type:       "conflict",
				Title:      "Conflict",
				Detail:     "Another activation is already in progress for this configuration",
				StatusCode: http.StatusConflict,
			},
		},
		"422 unprocessable entity - invalid hosts": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID:                         43253,
				ConfigVersion:                    1,
				Action:                           "ACTIVATE",
				Network:                          "STAGING",
				Note:                             "Test activation with invalid hosts",
				NotificationEmails:               []string{"test@example.com"},
				AcknowledgedInvalidHosts:         []string{},
				AcknowledgedInvalidHostsByConfig: []AcknowledgedInvalidHostsByConfig{},
				HostsToMove:                      []HostToMove{},
			},
			headers:        http.Header{},
			responseStatus: http.StatusUnprocessableEntity,
			responseBody: `{
				"type": "unprocessable_entity",
				"title": "Unprocessable Entity",
				"detail": "Some hosts in the configuration are invalid and must be acknowledged",
				"status": 422
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/1/activations-with-host-move",
			withError: &Error{
				Type:       "unprocessable_entity",
				Title:      "Unprocessable Entity",
				Detail:     "Some hosts in the configuration are invalid and must be acknowledged",
				StatusCode: http.StatusUnprocessableEntity,
			},
		},
		"missing ConfigID": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigVersion: 1,
				Action:        "ACTIVATE",
				Network:       "STAGING",
			},
			withError: ErrStructValidation,
		},
		"missing ConfigVersion": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID: 43253,
				Action:   "ACTIVATE",
				Network:  "STAGING",
			},
			withError: ErrStructValidation,
		},
		"missing Action": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID:      43253,
				ConfigVersion: 1,
				Network:       "STAGING",
			},
			withError: ErrStructValidation,
		},
		"missing Network": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID:      43253,
				ConfigVersion: 1,
				Action:        "ACTIVATE",
			},
			withError: ErrStructValidation,
		},
		"incorrect Network": {
			params: CreateActivationsWithHostMoveRequest{
				ConfigID:      43253,
				ConfigVersion: 1,
				Action:        "ACTIVATE",
				Network:       "INCORRECT",
			},
			withError: ErrStructValidation,
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
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateActivationsWithHostMove(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
