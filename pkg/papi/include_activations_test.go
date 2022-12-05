package papi

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivateInclude(t *testing.T) {
	tests := map[string]struct {
		params              ActivateIncludeRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *ActivationIncludeResponse
		withError           error
		assertError         func(*testing.T, error)
	}{
		"201 Activate include acknowledging all the warnings": {
			params: ActivateIncludeRequest{
				IncludeID:              "inc_12345",
				Version:                4,
				Network:                ActivationNetworkStaging,
				Note:                   "test activation",
				NotifyEmails:           []string{"jbond@example.com"},
				AcknowledgeAllWarnings: true,
			},
			expectedRequestBody: `{"acknowledgeAllWarnings":true,"activationType":"ACTIVATE","ignoreHttpErrors":true,"includeVersion":4,"network":"STAGING","note":"test activation","notifyEmails":["jbond@example.com"]}`,
			expectedPath:        "/papi/v1/includes/inc_12345/activations",
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "activationLink": "/papi/v1/includes/inc_12345/activations/temporary-activation-id"
}`,
			expectedResponse: &ActivationIncludeResponse{
				ActivationID:   "temporary-activation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
			},
		},
		"201 Activate include": {
			params: ActivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				Network:      ActivationNetworkStaging,
				Note:         "test activation",
				NotifyEmails: []string{"jbond@example.com"},
			},
			expectedRequestBody: `{"acknowledgeAllWarnings":false,"activationType":"ACTIVATE","ignoreHttpErrors":true,"includeVersion":4,"network":"STAGING","note":"test activation","notifyEmails":["jbond@example.com"]}`,
			expectedPath:        "/papi/v1/includes/inc_12345/activations",
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "activationLink": "/papi/v1/includes/inc_12345/activations/temporary-activation-id"
}`,
			expectedResponse: &ActivationIncludeResponse{
				ActivationID:   "temporary-activation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
			},
		},
		"201 Activate include with ComplianceRecord None": {
			params: ActivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				Network:      ActivationNetworkProduction,
				Note:         "test activation",
				NotifyEmails: []string{"jbond@example.com"},
				ComplianceRecord: &ComplianceRecordNone{
					CustomerEmail:  "sb@akamai.com",
					PeerReviewedBy: "sb@akamai.com",
					UnitTested:     true,
					TicketID:       "123",
				},
			},
			expectedRequestBody: `{"acknowledgeAllWarnings":false,"activationType":"ACTIVATE","ignoreHttpErrors":true,"includeVersion":4,"network":"PRODUCTION","note":"test activation","notifyEmails":["jbond@example.com"], "complianceRecord":{"customerEmail":"sb@akamai.com", "noncomplianceReason":"NONE", "peerReviewedBy":"sb@akamai.com", "unitTested":true, "ticketId":"123"}}`,
			expectedPath:        "/papi/v1/includes/inc_12345/activations",
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "activationLink": "/papi/v1/includes/inc_12345/activations/temporary-activation-id"
}`,
			expectedResponse: &ActivationIncludeResponse{
				ActivationID:   "temporary-activation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
			},
		},
		"201 Activate include with ComplianceRecord Other": {
			params: ActivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				Network:      ActivationNetworkProduction,
				Note:         "test activation",
				NotifyEmails: []string{"jbond@example.com"},
				ComplianceRecord: &ComplianceRecordOther{
					OtherNoncomplianceReason: "some other reason",
					TicketID:                 "123",
				},
			},
			expectedRequestBody: `{"acknowledgeAllWarnings":false,"activationType":"ACTIVATE","ignoreHttpErrors":true,"includeVersion":4,"network":"PRODUCTION","note":"test activation","notifyEmails":["jbond@example.com"], "complianceRecord":{"otherNoncomplianceReason":"some other reason", "noncomplianceReason":"OTHER", "ticketId":"123"}}`,
			expectedPath:        "/papi/v1/includes/inc_12345/activations",
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "activationLink": "/papi/v1/includes/inc_12345/activations/temporary-activation-id"
}`,
			expectedResponse: &ActivationIncludeResponse{
				ActivationID:   "temporary-activation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
			},
		},
		"201 Activate include with ComplianceRecord No_Production_Traffic": {
			params: ActivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				Network:      ActivationNetworkProduction,
				Note:         "test activation",
				NotifyEmails: []string{"jbond@example.com"},
				ComplianceRecord: &ComplianceRecordNoProductionTraffic{
					TicketID: "123",
				},
			},
			expectedRequestBody: `{"acknowledgeAllWarnings":false,"activationType":"ACTIVATE","ignoreHttpErrors":true,"includeVersion":4,"network":"PRODUCTION","note":"test activation","notifyEmails":["jbond@example.com"], "complianceRecord":{"noncomplianceReason":"NO_PRODUCTION_TRAFFIC", "ticketId":"123"}}`,
			expectedPath:        "/papi/v1/includes/inc_12345/activations",
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "activationLink": "/papi/v1/includes/inc_12345/activations/temporary-activation-id"
}`,
			expectedResponse: &ActivationIncludeResponse{
				ActivationID:   "temporary-activation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
			},
		},
		"201 Activate include with ComplianceRecord Emergency": {
			params: ActivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				Network:      ActivationNetworkProduction,
				Note:         "test activation",
				NotifyEmails: []string{"jbond@example.com"},
				ComplianceRecord: &ComplianceRecordEmergency{
					TicketID: "123",
				},
			},
			expectedRequestBody: `{"acknowledgeAllWarnings":false,"activationType":"ACTIVATE","ignoreHttpErrors":true,"includeVersion":4,"network":"PRODUCTION","note":"test activation","notifyEmails":["jbond@example.com"], "complianceRecord":{"noncomplianceReason":"EMERGENCY", "ticketId":"123"}}`,
			expectedPath:        "/papi/v1/includes/inc_12345/activations",
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "activationLink": "/papi/v1/includes/inc_12345/activations/temporary-activation-id"
}`,
			expectedResponse: &ActivationIncludeResponse{
				ActivationID:   "temporary-activation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
			},
		},
		"500 internal server error": {
			params: ActivateIncludeRequest{
				IncludeID:              "inc_12345",
				Version:                4,
				Network:                ActivationNetworkStaging,
				Note:                   "test activation",
				NotifyEmails:           []string{"jbond@example.com"},
				AcknowledgeAllWarnings: true,
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error getting include",
		   "status": 500
		}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing include id": {
			params: ActivateIncludeRequest{
				Version:      4,
				Network:      ActivationNetworkStaging,
				NotifyEmails: []string{"jbond@example.com"},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing version": {
			params: ActivateIncludeRequest{
				IncludeID:    "inc_12345",
				Network:      ActivationNetworkStaging,
				NotifyEmails: []string{"jbond@example.com"},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing network": {
			params: ActivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				NotifyEmails: []string{"jbond@example.com"},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing notify emails": {
			params: ActivateIncludeRequest{
				IncludeID: "inc_12345",
				Version:   4,
				Network:   ActivationNetworkStaging,
			},
			withError: ErrStructValidation,
		},
		"validation error - missing compliance record for production network": {
			params: ActivateIncludeRequest{
				IncludeID:              "inc_12345",
				Version:                4,
				Network:                ActivationNetworkProduction,
				Note:                   "test activation",
				NotifyEmails:           []string{"jbond@example.com"},
				AcknowledgeAllWarnings: true,
			},
			withError: ErrStructValidation,
		},
		"validation error - not valid ComplianceRecordNone": {
			params: ActivateIncludeRequest{
				IncludeID:              "inc_12345",
				Version:                4,
				Network:                ActivationNetworkProduction,
				Note:                   "test activation",
				NotifyEmails:           []string{"jbond@example.com"},
				AcknowledgeAllWarnings: true,
				ComplianceRecord: &ComplianceRecordNone{
					UnitTested: true,
					TicketID:   "123",
				},
			},
			withError: ErrStructValidation,
			assertError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "CustomerEmail: cannot be blank")
				assert.Contains(t, err.Error(), "PeerReviewedBy: cannot be blank")
			},
		},
		"validation error - not valid UnitTested field for PRODUCTION activation network and ComplianceRecordNone": {
			params: ActivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				Network:      ActivationNetworkProduction,
				Note:         "test activation",
				NotifyEmails: []string{"jbond@example.com"},
				ComplianceRecord: &ComplianceRecordNone{
					CustomerEmail:  "sb@akamai.com",
					PeerReviewedBy: "sb@akamai.com",
					UnitTested:     false,
					TicketID:       "123",
				},
			},
			withError: ErrStructValidation,
			assertError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "for PRODUCTION activation network and nonComplianceRecord, UnitTested value has to be set to true, otherwise API will not work correctly")
			},
		},
		"validation error - not valid ComplianceRecordOther": {
			params: ActivateIncludeRequest{
				IncludeID:              "inc_12345",
				Version:                4,
				Network:                ActivationNetworkProduction,
				Note:                   "test activation",
				NotifyEmails:           []string{"jbond@example.com"},
				AcknowledgeAllWarnings: true,
				ComplianceRecord:       &ComplianceRecordOther{},
			},
			withError: ErrStructValidation,
			assertError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "OtherNoncomplianceReason: cannot be blank")
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ActivateInclude(context.Background(), test.params)

			if test.withError != nil || test.assertError != nil {
				if test.withError != nil {
					assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				}
				if test.assertError != nil {
					test.assertError(t, err)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeactivateInclude(t *testing.T) {
	tests := map[string]struct {
		params              DeactivateIncludeRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *DeactivationIncludeResponse
		withError           error
	}{
		"201 Activate include acknowledging all the warnings": {
			params: DeactivateIncludeRequest{
				IncludeID:              "inc_12345",
				Version:                4,
				Network:                ActivationNetworkStaging,
				Note:                   "test activation",
				NotifyEmails:           []string{"jbond@example.com"},
				AcknowledgeAllWarnings: true,
			},
			expectedRequestBody: `{"acknowledgeAllWarnings":true,"activationType":"DEACTIVATE","ignoreHttpErrors":true,"includeVersion":4,"network":"STAGING","note":"test activation","notifyEmails":["jbond@example.com"]}`,
			expectedPath:        "/papi/v1/includes/inc_12345/activations",
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "activationLink": "/papi/v1/includes/inc_12345/activations/temporary-activation-id"
}`,
			expectedResponse: &DeactivationIncludeResponse{
				ActivationID:   "temporary-activation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
			},
		},
		"201 Activate include": {
			params: DeactivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				Network:      ActivationNetworkStaging,
				Note:         "test activation",
				NotifyEmails: []string{"jbond@example.com"},
			},
			expectedRequestBody: `{"acknowledgeAllWarnings":false,"activationType":"DEACTIVATE","ignoreHttpErrors":true,"includeVersion":4,"network":"STAGING","note":"test activation","notifyEmails":["jbond@example.com"]}`,
			expectedPath:        "/papi/v1/includes/inc_12345/activations",
			responseStatus:      http.StatusCreated,
			responseBody: `
		{
		   "activationLink": "/papi/v1/includes/inc_12345/activations/temporary-activation-id"
		}`,
			expectedResponse: &DeactivationIncludeResponse{
				ActivationID:   "temporary-activation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
			},
		},
		"422 Unprocessable entity - deactivate version which is not active on some network": {
			params: DeactivateIncludeRequest{
				IncludeID:              "inc_12345",
				Version:                4,
				Network:                ActivationNetworkProduction,
				Note:                   "test activation",
				NotifyEmails:           []string{"jbond@example.com"},
				AcknowledgeAllWarnings: true,
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations",
			responseStatus: http.StatusUnprocessableEntity,
			responseBody: `
{
    "type": "https://problems.luna.akamaiapis.net/papi/v0/deactivation/include-not-active-in-production",
    "title": "Include not active in PRODUCTION",
    "detail": "The include cannot be deactivated because it is not active in PRODUCTION.",
    "instance": "https://akaa-gcplhccxrheyl6kw-bcfnozqkbaydivqp.luna-dev.akamaiapis.net/papi/v1/includes/inc_12345/activations#12345",
    "status": 422
}`,
			withError: &Error{
				Type:       "https://problems.luna.akamaiapis.net/papi/v0/deactivation/include-not-active-in-production",
				Title:      "Include not active in PRODUCTION",
				Detail:     "The include cannot be deactivated because it is not active in PRODUCTION.",
				Instance:   "https://akaa-gcplhccxrheyl6kw-bcfnozqkbaydivqp.luna-dev.akamaiapis.net/papi/v1/includes/inc_12345/activations#12345",
				StatusCode: http.StatusUnprocessableEntity,
			},
		},
		"500 internal server error": {
			params: DeactivateIncludeRequest{
				IncludeID:              "inc_12345",
				Version:                4,
				Network:                ActivationNetworkStaging,
				Note:                   "test activation",
				NotifyEmails:           []string{"jbond@example.com"},
				AcknowledgeAllWarnings: true,
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
				   "title": "Internal Server Error",
				   "detail": "Error getting include",
				   "status": 500
				}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing include id": {
			params: DeactivateIncludeRequest{
				Version:      4,
				Network:      ActivationNetworkStaging,
				NotifyEmails: []string{"jbond@example.com"},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing version": {
			params: DeactivateIncludeRequest{
				IncludeID:    "inc_12345",
				Network:      ActivationNetworkStaging,
				NotifyEmails: []string{"jbond@example.com"},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing network": {
			params: DeactivateIncludeRequest{
				IncludeID:    "inc_12345",
				Version:      4,
				NotifyEmails: []string{"jbond@example.com"},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing notify emails": {
			params: DeactivateIncludeRequest{
				IncludeID: "inc_12345",
				Version:   4,
				Network:   ActivationNetworkStaging,
			},
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
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeactivateInclude(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCancelIncludeActivation(t *testing.T) {
	tests := map[string]struct {
		params           CancelIncludeActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CancelIncludeActivationResponse
		withError        error
	}{
		"200 cancel include activation": {
			params: CancelIncludeActivationRequest{
				IncludeID:    "inc_12345",
				ContractID:   "test_contract",
				GroupID:      "test_group",
				ActivationID: "test_activation_123",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations/test_activation_123?contractId=test_contract&groupId=test_group",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "test_account",
    "contractId": "test_contract",
    "groupId": "test_group",
    "activations": {
        "items": [
            {
                "network": "STAGING",
                "activationType": "ACTIVATE",
                "status": "PENDING_CANCELLATION",
                "submitDate": "2022-12-01T13:18:57Z",
                "updateDate": "2022-12-01T13:19:04Z",
                "note": "test_note_1",
                "notifyEmails": [
                    "nomail@nomail.com"
                ],
                "fmaActivationState": "received",
                "includeId": "inc_12345",
                "includeName": "test_include_name",
                "includeVersion": 1,
                "includeActivationId": "test_activation_123"
            }
        ]
    }
}`,
			expectedResponse: &CancelIncludeActivationResponse{
				AccountID:  "test_account",
				ContractID: "test_contract",
				GroupID:    "test_group",
				Activations: IncludeActivationsRes{
					Items: []IncludeActivation{
						{
							Network:             "STAGING",
							ActivationType:      ActivationTypeActivate,
							Status:              ActivationStatusCancelling,
							SubmitDate:          "2022-12-01T13:18:57Z",
							UpdateDate:          "2022-12-01T13:19:04Z",
							Note:                "test_note_1",
							NotifyEmails:        []string{"nomail@nomail.com"},
							FMAActivationState:  "received",
							IncludeID:           "inc_12345",
							IncludeName:         "test_include_name",
							IncludeVersion:      1,
							IncludeActivationID: "test_activation_123",
						},
					},
				},
			},
		},
		"500 internal server error": {
			params: CancelIncludeActivationRequest{
				IncludeID:    "inc_12345",
				ContractID:   "test_contract",
				GroupID:      "test_group",
				ActivationID: "test_activation_123",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations/test_activation_123?contractId=test_contract&groupId=test_group",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
				   "title": "Internal Server Error",
				   "detail": "Error cancelling include activation",
				   "status": 500
				}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error cancelling include activation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing include id": {
			params: CancelIncludeActivationRequest{
				ContractID:   "test_contract",
				GroupID:      "test_group",
				ActivationID: "test_activation_123",
			},
			withError: ErrStructValidation,
		},
		"validation error - contract id": {
			params: CancelIncludeActivationRequest{
				IncludeID:    "inc_12345",
				GroupID:      "test_group",
				ActivationID: "test_activation_123",
			},
			withError: ErrStructValidation,
		},
		"validation error - group id": {
			params: CancelIncludeActivationRequest{
				IncludeID:    "inc_12345",
				ContractID:   "test_contract",
				ActivationID: "test_activation_123",
			},
			withError: ErrStructValidation,
		},
		"validation error - activation id": {
			params: CancelIncludeActivationRequest{
				IncludeID:  "inc_12345",
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			withError: ErrStructValidation,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CancelIncludeActivation(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetIncludeActivation(t *testing.T) {
	tests := map[string]struct {
		params           GetIncludeActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetIncludeActivationResponse
		withError        error
	}{
		"200 Get include activation": {
			params: GetIncludeActivationRequest{
				IncludeID:    "inc_12345",
				ActivationID: "atv_12345",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations/atv_12345",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "test_account",
    "contractId": "test_contract",
    "groupId": "test_group",
    "activations": {
        "items": [
            {
                "activationId": "atv_12345",
                "network": "STAGING",
                "activationType": "ACTIVATE",
                "status": "ACTIVE",
                "submitDate": "2022-10-27T12:27:54Z",
                "updateDate": "2022-10-27T12:28:54Z",
                "note": "DXE test activation",
                "notifyEmails": [
                    "test@example.com"
                ],
                "fmaActivationState": "steady",
                "fallbackInfo": {
                    "fastFallbackAttempted": false,
                    "fallbackVersion": 3,
                    "canFastFallback": false,
                    "steadyStateTime": 1666873734,
                    "fastFallbackExpirationTime": 1666877334,
                    "fastFallbackRecoveryState": null
                },
                "includeId": "inc_12345",
                "includeName": "tfp_test1",
                "includeType": "MICROSERVICES",
                "includeVersion": 4
            }
        ]
    }
}`,
			expectedResponse: &GetIncludeActivationResponse{
				AccountID:  "test_account",
				ContractID: "test_contract",
				GroupID:    "test_group",
				Activations: IncludeActivationsRes{
					Items: []IncludeActivation{
						{
							ActivationID:       "atv_12345",
							Network:            "STAGING",
							ActivationType:     ActivationTypeActivate,
							Status:             ActivationStatusActive,
							SubmitDate:         "2022-10-27T12:27:54Z",
							UpdateDate:         "2022-10-27T12:28:54Z",
							Note:               "DXE test activation",
							NotifyEmails:       []string{"test@example.com"},
							FMAActivationState: "steady",
							FallbackInfo: &ActivationFallbackInfo{
								FastFallbackAttempted:      false,
								FallbackVersion:            3,
								CanFastFallback:            false,
								SteadyStateTime:            1666873734,
								FastFallbackExpirationTime: 1666877334,
							},
							IncludeID:      "inc_12345",
							IncludeName:    "tfp_test1",
							IncludeType:    "MICROSERVICES",
							IncludeVersion: 4,
						},
					},
				},
				Activation: IncludeActivation{
					ActivationID:       "atv_12345",
					Network:            "STAGING",
					ActivationType:     ActivationTypeActivate,
					Status:             ActivationStatusActive,
					SubmitDate:         "2022-10-27T12:27:54Z",
					UpdateDate:         "2022-10-27T12:28:54Z",
					Note:               "DXE test activation",
					NotifyEmails:       []string{"test@example.com"},
					FMAActivationState: "steady",
					FallbackInfo: &ActivationFallbackInfo{
						FastFallbackAttempted:      false,
						FallbackVersion:            3,
						CanFastFallback:            false,
						SteadyStateTime:            1666873734,
						FastFallbackExpirationTime: 1666877334,
					},
					IncludeID:      "inc_12345",
					IncludeName:    "tfp_test1",
					IncludeType:    "MICROSERVICES",
					IncludeVersion: 4,
				},
			},
		},
		"200 Get include activation with includeActivationId": {
			params: GetIncludeActivationRequest{
				IncludeID:    "inc_12345",
				ActivationID: "5e597860-1107-461e-8dbe-4e7526e8dd02",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations/5e597860-1107-461e-8dbe-4e7526e8dd02",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "test_account",
    "contractId": "test_contract",
    "groupId": "test_group",
    "activations": {
        "items": [
            {
                "includeActivationId": "5e597860-1107-461e-8dbe-4e7526e8dd02",
                "network": "STAGING",
                "activationType": "ACTIVATE",
                "status": "ACTIVE",
                "submitDate": "2022-10-27T12:27:54Z",
                "updateDate": "2022-10-27T12:28:54Z",
                "note": "DXE test activation",
                "notifyEmails": [
                    "test@example.com"
                ],
                "fmaActivationState": "steady",
                "fallbackInfo": {
                    "fastFallbackAttempted": false,
                    "fallbackVersion": 3,
                    "canFastFallback": false,
                    "steadyStateTime": 1666873734,
                    "fastFallbackExpirationTime": 1666877334,
                    "fastFallbackRecoveryState": null
                },
                "includeId": "inc_12345",
                "includeName": "tfp_test1",
                "includeType": "MICROSERVICES",
                "includeVersion": 4
            }
        ]
    }
}`,
			expectedResponse: &GetIncludeActivationResponse{
				AccountID:  "test_account",
				ContractID: "test_contract",
				GroupID:    "test_group",
				Activations: IncludeActivationsRes{
					Items: []IncludeActivation{
						{
							IncludeActivationID: "5e597860-1107-461e-8dbe-4e7526e8dd02",
							Network:             "STAGING",
							ActivationType:      ActivationTypeActivate,
							Status:              ActivationStatusActive,
							SubmitDate:          "2022-10-27T12:27:54Z",
							UpdateDate:          "2022-10-27T12:28:54Z",
							Note:                "DXE test activation",
							NotifyEmails:        []string{"test@example.com"},
							FMAActivationState:  "steady",
							FallbackInfo: &ActivationFallbackInfo{
								FastFallbackAttempted:      false,
								FallbackVersion:            3,
								CanFastFallback:            false,
								SteadyStateTime:            1666873734,
								FastFallbackExpirationTime: 1666877334,
							},
							IncludeID:      "inc_12345",
							IncludeName:    "tfp_test1",
							IncludeType:    "MICROSERVICES",
							IncludeVersion: 4,
						},
					},
				},
				Activation: IncludeActivation{
					IncludeActivationID: "5e597860-1107-461e-8dbe-4e7526e8dd02",
					Network:             "STAGING",
					ActivationType:      ActivationTypeActivate,
					Status:              ActivationStatusActive,
					SubmitDate:          "2022-10-27T12:27:54Z",
					UpdateDate:          "2022-10-27T12:28:54Z",
					Note:                "DXE test activation",
					NotifyEmails:        []string{"test@example.com"},
					FMAActivationState:  "steady",
					FallbackInfo: &ActivationFallbackInfo{
						FastFallbackAttempted:      false,
						FallbackVersion:            3,
						CanFastFallback:            false,
						SteadyStateTime:            1666873734,
						FastFallbackExpirationTime: 1666877334,
					},
					IncludeID:      "inc_12345",
					IncludeName:    "tfp_test1",
					IncludeType:    "MICROSERVICES",
					IncludeVersion: 4,
				},
			},
		},
		"500 internal server error": {
			params: GetIncludeActivationRequest{
				IncludeID:    "inc_12345",
				ActivationID: "atv_12345",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations/atv_12345",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
				   "title": "Internal Server Error",
				   "detail": "Error getting include",
				   "status": 500
				}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing include id": {
			params: GetIncludeActivationRequest{
				ActivationID: "atv_12345",
			},
			withError: ErrStructValidation,
		},
		"validation error - activation id": {
			params: GetIncludeActivationRequest{
				IncludeID: "inc_12345",
			},
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
			result, err := client.GetIncludeActivation(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListIncludeActivations(t *testing.T) {
	tests := map[string]struct {
		params           ListIncludeActivationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListIncludeActivationsResponse
		withError        error
	}{
		"200 List include activations": {
			params: ListIncludeActivationsRequest{
				IncludeID:  "inc_12345",
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations?contractId=test_contract&groupId=test_group",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "test_account",
    "contractId": "test_contract",
    "groupId": "test_group",
    "activations": {
        "items": [
            {
                "activationId": "atv_12344",
                "network": "STAGING",
                "activationType": "ACTIVATE",
                "status": "ACTIVE",
                "submitDate": "2022-10-27T12:27:54Z",
                "updateDate": "2022-10-27T12:28:54Z",
                "note": "test activation",
                "notifyEmails": [
                    "test@example.com"
                ],
                "fmaActivationState": "steady",
                "fallbackInfo": {
                    "fastFallbackAttempted": false,
                    "fallbackVersion": 3,
                    "canFastFallback": false,
                    "steadyStateTime": 1666873734,
                    "fastFallbackExpirationTime": 1666877334,
                    "fastFallbackRecoveryState": null
                },
                "includeId": "inc_12345",
                "includeName": "tfp_test1",
                "includeType": "MICROSERVICES",
                "includeVersion": 4
            },
            {
                "activationId": "atv_12343",
                "network": "STAGING",
                "activationType": "ACTIVATE",
                "status": "ACTIVE",
                "submitDate": "2022-10-27T11:21:40Z",
                "updateDate": "2022-10-27T11:22:54Z",
                "note": "test activation",
                "notifyEmails": [
                    "test@example.com"
                ],
                "fmaActivationState": "steady",
                "fallbackInfo": {
                    "fastFallbackAttempted": false,
                    "fallbackVersion": 4,
                    "canFastFallback": false,
                    "steadyStateTime": 1666869774,
                    "fastFallbackExpirationTime": 1666873374,
                    "fastFallbackRecoveryState": null
                },
                "includeId": "inc_12345",
                "includeName": "tfp_test1",
                "includeType": "MICROSERVICES",
                "includeVersion": 3
            },
            {
                "activationId": "atv_12343",
                "network": "STAGING",
                "activationType": "DEACTIVATE",
                "status": "ACTIVE",
                "submitDate": "2022-10-26T12:41:58Z",
                "updateDate": "2022-10-26T13:03:04Z",
                "note": "test activation",
                "notifyEmails": [
                    "test@example.com"
                ],
                "includeId": "inc_12345",
                "includeName": "tfp_test1",
                "includeType": "MICROSERVICES",
                "includeVersion": 3
            },
            {
                "activationId": "atv_12342",
                "network": "STAGING",
                "activationType": "ACTIVATE",
                "status": "ACTIVE",
                "submitDate": "2022-10-26T12:37:49Z",
                "updateDate": "2022-10-26T12:38:59Z",
                "note": "test activation",
                "notifyEmails": [
                    "test@example.com"
                ],
                "fmaActivationState": "steady",
                "fallbackInfo": {
                    "fastFallbackAttempted": false,
                    "fallbackVersion": 4,
                    "canFastFallback": false,
                    "steadyStateTime": 1666787939,
                    "fastFallbackExpirationTime": 1666791539,
                    "fastFallbackRecoveryState": null
                },
                "includeId": "inc_12345",
                "includeName": "tfp_test1",
                "includeType": "MICROSERVICES",
                "includeVersion": 2
            },
            {
                "activationId": "atv_12341",
                "network": "STAGING",
                "activationType": "ACTIVATE",
                "status": "ACTIVE",
                "submitDate": "2022-08-17T09:13:18Z",
                "updateDate": "2022-08-17T09:15:35Z",
                "note": "test activation",
                "notifyEmails": [
                    "test@example.com"
                ],
                "fmaActivationState": "steady",
                "fallbackInfo": {
                    "fastFallbackAttempted": false,
                    "fallbackVersion": 4,
                    "canFastFallback": false,
                    "steadyStateTime": 1660727735,
                    "fastFallbackExpirationTime": 1660731335,
                    "fastFallbackRecoveryState": null
                },
                "includeId": "inc_12345",
                "includeName": "tfp_test1",
                "includeType": "MICROSERVICES",
                "includeVersion": 1
            }
        ]
    }
}`,
			expectedResponse: &ListIncludeActivationsResponse{
				AccountID:  "test_account",
				ContractID: "test_contract",
				GroupID:    "test_group",
				Activations: IncludeActivationsRes{
					Items: []IncludeActivation{
						{
							ActivationID:       "atv_12344",
							Network:            "STAGING",
							ActivationType:     ActivationTypeActivate,
							Status:             ActivationStatusActive,
							SubmitDate:         "2022-10-27T12:27:54Z",
							UpdateDate:         "2022-10-27T12:28:54Z",
							Note:               "test activation",
							NotifyEmails:       []string{"test@example.com"},
							FMAActivationState: "steady",
							FallbackInfo: &ActivationFallbackInfo{
								FastFallbackAttempted:      false,
								FallbackVersion:            3,
								CanFastFallback:            false,
								SteadyStateTime:            1666873734,
								FastFallbackExpirationTime: 1666877334,
							},
							IncludeID:      "inc_12345",
							IncludeName:    "tfp_test1",
							IncludeType:    "MICROSERVICES",
							IncludeVersion: 4,
						},
						{
							ActivationID:       "atv_12343",
							Network:            "STAGING",
							ActivationType:     ActivationTypeActivate,
							Status:             ActivationStatusActive,
							SubmitDate:         "2022-10-27T11:21:40Z",
							UpdateDate:         "2022-10-27T11:22:54Z",
							Note:               "test activation",
							NotifyEmails:       []string{"test@example.com"},
							FMAActivationState: "steady",
							FallbackInfo: &ActivationFallbackInfo{
								FastFallbackAttempted:      false,
								FallbackVersion:            4,
								CanFastFallback:            false,
								SteadyStateTime:            1666869774,
								FastFallbackExpirationTime: 1666873374,
							},
							IncludeID:      "inc_12345",
							IncludeName:    "tfp_test1",
							IncludeType:    "MICROSERVICES",
							IncludeVersion: 3,
						},
						{
							ActivationID:   "atv_12343",
							Network:        "STAGING",
							ActivationType: ActivationTypeDeactivate,
							Status:         ActivationStatusActive,
							SubmitDate:     "2022-10-26T12:41:58Z",
							UpdateDate:     "2022-10-26T13:03:04Z",
							Note:           "test activation",
							NotifyEmails:   []string{"test@example.com"},
							IncludeID:      "inc_12345",
							IncludeName:    "tfp_test1",
							IncludeType:    "MICROSERVICES",
							IncludeVersion: 3,
						},
						{
							ActivationID:       "atv_12342",
							Network:            "STAGING",
							ActivationType:     ActivationTypeActivate,
							Status:             ActivationStatusActive,
							SubmitDate:         "2022-10-26T12:37:49Z",
							UpdateDate:         "2022-10-26T12:38:59Z",
							Note:               "test activation",
							NotifyEmails:       []string{"test@example.com"},
							FMAActivationState: "steady",
							FallbackInfo: &ActivationFallbackInfo{
								FastFallbackAttempted:      false,
								FallbackVersion:            4,
								CanFastFallback:            false,
								SteadyStateTime:            1666787939,
								FastFallbackExpirationTime: 1666791539,
							},
							IncludeID:      "inc_12345",
							IncludeName:    "tfp_test1",
							IncludeType:    "MICROSERVICES",
							IncludeVersion: 2,
						},
						{
							ActivationID:       "atv_12341",
							Network:            "STAGING",
							ActivationType:     ActivationTypeActivate,
							Status:             ActivationStatusActive,
							SubmitDate:         "2022-08-17T09:13:18Z",
							UpdateDate:         "2022-08-17T09:15:35Z",
							Note:               "test activation",
							NotifyEmails:       []string{"test@example.com"},
							FMAActivationState: "steady",
							FallbackInfo: &ActivationFallbackInfo{
								FastFallbackAttempted:      false,
								FallbackVersion:            4,
								CanFastFallback:            false,
								SteadyStateTime:            1660727735,
								FastFallbackExpirationTime: 1660731335,
							},
							IncludeID:      "inc_12345",
							IncludeName:    "tfp_test1",
							IncludeType:    "MICROSERVICES",
							IncludeVersion: 1,
						},
					},
				},
			},
		},
		"500 internal server error": {
			params: ListIncludeActivationsRequest{
				IncludeID:  "inc_12345",
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/activations?contractId=test_contract&groupId=test_group",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
				   "title": "Internal Server Error",
				   "detail": "Error getting include",
				   "status": 500
				}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing include id": {
			params: ListIncludeActivationsRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			withError: ErrStructValidation,
		},
		"validation error - contract id": {
			params: ListIncludeActivationsRequest{
				IncludeID: "inc_12345",
				GroupID:   "test_group",
			},
			withError: ErrStructValidation,
		},
		"validation error - group id": {
			params: ListIncludeActivationsRequest{
				IncludeID:  "inc_12345",
				ContractID: "test_contract",
			},
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
			result, err := client.ListIncludeActivations(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
