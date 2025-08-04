package clientlists

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateActivation(t *testing.T) {
	uri := "/client-list/v1/lists/1234_NORTHAMERICAGEOALLOWLIST/activations"

	tests := map[string]struct {
		params              CreateActivationRequest
		responseStatus      int
		expectedRequestBody string
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateActivationResponse
		withError           error
	}{
		"200 OK": {
			params: CreateActivationRequest{
				ListID: "1234_NORTHAMERICAGEOALLOWLIST",
				ActivationParams: ActivationParams{
					Action:                 Activate,
					Network:                Production,
					Comments:               "Activation of GEO allowlist list",
					SiebelTicketID:         "12_B",
					NotificationRecipients: []string{"a@a.com", "c@c.com"},
				},
			},
			expectedRequestBody: `{"action":"ACTIVATE","comments":"Activation of GEO allowlist list","network":"PRODUCTION","notificationRecipients":["a@a.com","c@c.com"],"siebelTicketId":"12_B"}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"action": "ACTIVATE",
				"activationStatus": "PENDING_ACTIVATION",
				"listId": "1234_NORTHAMERICAGEOALLOWLIST",
				"network": "PRODUCTION",
				"notificationRecipients": ["aa@dd.com"],
				"version": 1,
				"activationId": 12,
				"createDate": "2023-04-05T18:46:56.365Z",
				"createdBy": "jdoe",
				"network": "PRODUCTION",
				"comments": "Activation of GEO allowlist list",
				"siebelTicketId": "12_AB"
			}`,
			expectedPath: uri,
			expectedResponse: &CreateActivationResponse{
				Action:                 "ACTIVATE",
				ActivationID:           12,
				ActivationStatus:       PendingActivation,
				CreateDate:             time.Date(2023, 4, 5, 18, 46, 56, 365000000, time.UTC).Format(time.RFC3339Nano),
				CreatedBy:              "jdoe",
				Comments:               "Activation of GEO allowlist list",
				ListID:                 "1234_NORTHAMERICAGEOALLOWLIST",
				Network:                Production,
				NotificationRecipients: []string{"aa@dd.com"},
				SiebelTicketID:         "12_AB",
				Version:                1,
			},
		},
		"500 internal server error": {
			params: CreateActivationRequest{
				ListID: "1234_NORTHAMERICAGEOALLOWLIST",
				ActivationParams: ActivationParams{
					Network: Production,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error creating client lists activation",
					"status": 500
				}`,
			expectedPath: uri,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating client lists activation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			params:    CreateActivationRequest{},
			withError: ErrStructValidation,
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

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateActivation(
				session.ContextWithOptions(
					context.Background(),
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

func TestCreateDeactivation(t *testing.T) {
	uri := "/client-list/v1/lists/1234_NORTHAMERICAGEOALLOWLIST/activations"
	deactivationReq := CreateDeactivationRequest{
		ListID: "1234_NORTHAMERICAGEOALLOWLIST",
		ActivationParams: ActivationParams{
			Action:                 Deactivate,
			Network:                Production,
			Comments:               "Activation of GEO allowlist list",
			SiebelTicketID:         "12_B",
			NotificationRecipients: []string{"a@a.com", "c@c.com"},
		},
	}

	tests := map[string]struct {
		params              CreateDeactivationRequest
		responseStatus      int
		expectedRequestBody string
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateDeactivationResponse
		withError           error
	}{
		"200 OK": {
			params:              deactivationReq,
			expectedRequestBody: `{"action":"DEACTIVATE","comments":"Activation of GEO allowlist list","network":"PRODUCTION","notificationRecipients":["a@a.com","c@c.com"],"siebelTicketId":"12_B"}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"action": "DEACTIVATE",
				"activationStatus": "PENDING_DEACTIVATION",
				"listId": "1234_NORTHAMERICAGEOALLOWLIST",
				"network": "PRODUCTION",
				"notificationRecipients": ["aa@dd.com"],
				"version": 1,
				"activationId": 12,
				"createDate": "2023-04-05T18:46:56.365Z",
				"createdBy": "jdoe",
				"network": "PRODUCTION",
				"comments": "Activation of GEO allowlist list",
				"siebelTicketId": "12_AB"
			}`,
			expectedPath: uri,
			expectedResponse: &CreateDeactivationResponse{
				Action:                 "DEACTIVATE",
				ActivationID:           12,
				ActivationStatus:       PendingDeactivation,
				CreateDate:             time.Date(2023, 4, 5, 18, 46, 56, 365000000, time.UTC).Format(time.RFC3339Nano),
				CreatedBy:              "jdoe",
				Comments:               "Activation of GEO allowlist list",
				ListID:                 "1234_NORTHAMERICAGEOALLOWLIST",
				Network:                Production,
				NotificationRecipients: []string{"aa@dd.com"},
				SiebelTicketID:         "12_AB",
				Version:                1,
			},
		},
		"403 forbidden operation error": {
			params:         deactivationReq,
			responseStatus: http.StatusForbidden,
			responseBody: `
				{
					"type": "https://problems.luna.akamaiapis.net/client-list/error-types/FORBIDDEN-OPERATION",
					"title": "This operation is invalid",
					"detail": "Only not used client list are allowed to be deactivated",
					"status": 403
				}`,
			expectedPath: uri,
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/client-list/error-types/FORBIDDEN-OPERATION",
				Title:      "This operation is invalid",
				Detail:     "Only not used client list are allowed to be deactivated",
				StatusCode: http.StatusForbidden,
			},
		},
		"404 not found error": {
			params:         deactivationReq,
			responseStatus: http.StatusNotFound,
			responseBody: `
				{
					"type": "https://problems.luna.akamaiapis.net/client-list/error-types/NOT-FOUND",
					"title": "Not Found",
					"detail": "Client list with uid 185417_IPTYPECLIENTLIST not found",
					"status": 404
				}`,
			expectedPath: uri,
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/client-list/error-types/NOT-FOUND",
				Title:      "Not Found",
				Detail:     "Client list with uid 185417_IPTYPECLIENTLIST not found",
				StatusCode: http.StatusNotFound,
			},
		},
		"500 internal server error": {
			params:         deactivationReq,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "https://problems.luna.akamaiapis.net/client-list/error-types/INTERNAL-SERVER-ERROR",
					"title": "Internal Server Error",
					"detail": "Error creating client lists deactivation",
					"status": 500
				}`,
			expectedPath: uri,
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/client-list/error-types/INTERNAL-SERVER-ERROR",
				Title:      "Internal Server Error",
				Detail:     "Error creating client lists deactivation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			params:    CreateDeactivationRequest{},
			withError: ErrStructValidation,
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

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateDeactivation(
				session.ContextWithOptions(
					context.Background(),
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

func TestGetActivation(t *testing.T) {
	uri := "/client-list/v1/activations/12"

	tests := map[string]struct {
		params           GetActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetActivationResponse
		withError        error
	}{
		"200 OK": {
			params: GetActivationRequest{
				ActivationID: 12,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"action": "ACTIVATE",
				"activationId": 12,
				"activationStatus": "PENDING_ACTIVATION",
				"comments": "latest activation",
				"createDate": "2023-04-05T18:46:56.365Z",
				"createdBy": "jdoe",
				"fast": true,
				"listId": "1234_NORTHAMERICAGEOALLOWLIST",
				"network": "PRODUCTION",
				"notificationRecipients": [
						"qw@ff.com"
				],
				"siebelTicketId": "q",
				"version": 1
		}`,
			expectedPath: uri,
			expectedResponse: &GetActivationResponse{
				ActivationID:      12,
				ListID:            "1234_NORTHAMERICAGEOALLOWLIST",
				Version:           1,
				CreateDate:        "2023-04-05T18:46:56.365Z",
				CreatedBy:         "jdoe",
				Fast:              true,
				InitialActivation: false,
				ActivationStatus:  "PENDING_ACTIVATION",
				ActivationParams: ActivationParams{
					Action:                 Activate,
					NotificationRecipients: []string{"qw@ff.com"},
					Comments:               "latest activation",
					Network:                Production,
					SiebelTicketID:         "q",
				},
			},
		},
		"500 internal server error": {
			params: GetActivationRequest{
				ActivationID: 12,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error fetching client lists activation",
					"status": 500
				}`,
			expectedPath: uri,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching client lists activation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			params:    GetActivationRequest{},
			withError: ErrStructValidation,
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
			result, err := client.GetActivation(
				session.ContextWithOptions(
					context.Background(),
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

func TestGetActivationStatus(t *testing.T) {
	uri := "/client-list/v1/lists/1234_NORTHAMERICAGEOALLOWLIST/environments/PRODUCTION/status"

	tests := map[string]struct {
		params           GetActivationStatusRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetActivationStatusResponse
		withError        error
	}{
		"200 OK": {
			params: GetActivationStatusRequest{
				ListID:  "1234_NORTHAMERICAGEOALLOWLIST",
				Network: Production,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"action": "ACTIVATE",
				"activationStatus": "PENDING_ACTIVATION",
				"listId": "1234_NORTHAMERICAGEOALLOWLIST",
				"network": "PRODUCTION",
				"notificationRecipients": [],
				"version": 1,
				"activationId": 12,
				"createDate": "2023-04-05T18:46:56.365Z",
				"createdBy": "jdoe",
				"network": "PRODUCTION",
				"comments": "Activation of GEO allowlist list",
				"siebelTicketId": "12_AB"
			}`,
			expectedPath: uri,
			expectedResponse: &GetActivationStatusResponse{
				Action:                 "ACTIVATE",
				ActivationID:           12,
				ActivationStatus:       PendingActivation,
				CreateDate:             "2023-04-05T18:46:56.365Z",
				CreatedBy:              "jdoe",
				Comments:               "Activation of GEO allowlist list",
				ListID:                 "1234_NORTHAMERICAGEOALLOWLIST",
				Network:                Production,
				NotificationRecipients: []string{},
				SiebelTicketID:         "12_AB",
				Version:                1,
			},
		},
		"500 internal server error": {
			params: GetActivationStatusRequest{
				ListID:  "1234_NORTHAMERICAGEOALLOWLIST",
				Network: Production,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error fetching client lists activation",
					"status": 500
				}`,
			expectedPath: uri,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching client lists activation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			params:    GetActivationStatusRequest{},
			withError: ErrStructValidation,
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
			result, err := client.GetActivationStatus(
				session.ContextWithOptions(
					context.Background(),
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
