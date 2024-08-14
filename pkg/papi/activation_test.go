package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiCreateActivation(t *testing.T) {
	tests := map[string]struct {
		request          CreateActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateActivationResponse
		withError        error
		assertError      func(*testing.T, error)
	}{
		"200 OK": {
			request: CreateActivationRequest{
				PropertyID: "prp_175780",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkStaging,
					UseFastFallback: false,
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"activationLink": "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225"
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &CreateActivationResponse{
				ActivationID:   "atv_67037",
				ActivationLink: "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225",
			},
		},
		"200 Activate property with ComplianceRecord None": {
			request: CreateActivationRequest{
				PropertyID: "prp_175780",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkStaging,
					UseFastFallback: false,
					ComplianceRecord: &ComplianceRecordNone{
						CustomerEmail:  "sb@akamai.com",
						PeerReviewedBy: "sb@akamai.com",
						UnitTested:     true,
						TicketID:       "123",
					},
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"activationLink": "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225"
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &CreateActivationResponse{
				ActivationID:   "atv_67037",
				ActivationLink: "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225",
			},
		},
		"200 Activate property with ComplianceRecord Other": {
			request: CreateActivationRequest{
				PropertyID: "prp_175780",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkStaging,
					UseFastFallback: false,
					ComplianceRecord: &ComplianceRecordOther{
						OtherNoncomplianceReason: "some other reason",
						TicketID:                 "123",
					},
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"activationLink": "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225"
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &CreateActivationResponse{
				ActivationID:   "atv_67037",
				ActivationLink: "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225",
			},
		},
		"200 Activate property with ComplianceRecord No_Production_Traffic": {
			request: CreateActivationRequest{
				PropertyID: "prp_175780",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkStaging,
					UseFastFallback: false,
					ComplianceRecord: &ComplianceRecordNoProductionTraffic{
						TicketID: "123",
					},
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"activationLink": "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225"
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &CreateActivationResponse{
				ActivationID:   "atv_67037",
				ActivationLink: "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225",
			},
		},
		"200 Activate property with ComplianceRecord Emergency": {
			request: CreateActivationRequest{
				PropertyID: "prp_175780",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkStaging,
					UseFastFallback: false,
					ComplianceRecord: &ComplianceRecordEmergency{
						TicketID: "123",
					},
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"activationLink": "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225"
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &CreateActivationResponse{
				ActivationID:   "atv_67037",
				ActivationLink: "/papi/v1/properties/prp_173136/activations/atv_67037?contractId=ctr_1-1TJZFB&groupId=grp_15225",
			},
		},
		"500 internal server error": {
			request: CreateActivationRequest{
				PropertyID: "prp_175780",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkStaging,
					UseFastFallback: false,
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating activation",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating activation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing property ID": {
			request: CreateActivationRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkStaging,
					UseFastFallback: false,
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			withError: ErrStructValidation,
		},
		"validation error - not valid ComplianceRecordNone": {
			request: CreateActivationRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkStaging,
					UseFastFallback: false,
					ComplianceRecord: &ComplianceRecordNone{
						UnitTested: true,
						TicketID:   "123",
					},
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			withError: ErrStructValidation,
			assertError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "CustomerEmail: cannot be blank")
				assert.Contains(t, err.Error(), "PeerReviewedBy: cannot be blank")
			},
		},
		"validation error - not valid UnitTested field for PRODUCTION activation network and ComplianceRecordNone": {
			request: CreateActivationRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion: 1,
					Network:         ActivationNetworkProduction,
					UseFastFallback: false,
					ComplianceRecord: &ComplianceRecordNone{
						CustomerEmail:  "sb@akamai.com",
						PeerReviewedBy: "sb@akamai.com",
						UnitTested:     false,
						TicketID:       "123",
					},
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
			},
			withError: ErrStructValidation,
			assertError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "for PRODUCTION activation network and nonComplianceRecord, UnitTested value has to be set to true, otherwise API will not work correctly")
			},
		},
		"validation error - not valid ComplianceRecordOther": {
			request: CreateActivationRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Activation: Activation{
					PropertyVersion:  1,
					Network:          ActivationNetworkProduction,
					UseFastFallback:  false,
					ComplianceRecord: &ComplianceRecordOther{},
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					AcknowledgeWarnings: []string{"msg_baa4560881774a45b5fd25f5b1eab021d7c40b4f"},
				},
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
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateActivation(context.Background(), test.request)
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

func TestPapiGetActivations(t *testing.T) {
	tests := map[string]struct {
		request          GetActivationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetActivationsResponse
		withError        error
	}{
		"200 OK": {
			request: GetActivationsRequest{
				PropertyID: "prp_175780",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"accountId": "act_1-1TJZFB",
	"contractId": "ctr_1-1TJZFW",
	"groupId": "grp_15166",
	"activations": {
		"items": [
			{
				"activationId": "atv_1696985",
				"propertyName": "example.com",
				"propertyId": "prp_173136",
				"propertyVersion": 1,
				"network": "STAGING",
				"activationType": "ACTIVATE",
				"status": "PENDING",
				"submitDate": "2014-03-02T02:22:12Z",
				"updateDate": "2014-03-01T21:12:57Z",
				"note": "Sample activation",
				"fmaActivationState": "steady",
				"notifyEmails": [
					"you@example.com",
					"them@example.com"
				],
				"fallbackInfo": {
					"fastFallbackAttempted": false,
					"fallbackVersion": 10,
					"canFastFallback": true,
					"steadyStateTime": 1506448172,
					"fastFallbackExpirationTime": 1506451772,
					"fastFallbackRecoveryState": null
				}
			}
		]
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &GetActivationsResponse{
				Response: Response{
					AccountID:  "act_1-1TJZFB",
					ContractID: "ctr_1-1TJZFW",
					GroupID:    "grp_15166",
				},
				Activations: ActivationsItems{Items: []*Activation{{
					ActivationID:       "atv_1696985",
					PropertyName:       "example.com",
					PropertyID:         "prp_173136",
					PropertyVersion:    1,
					Network:            ActivationNetworkStaging,
					ActivationType:     ActivationTypeActivate,
					Status:             ActivationStatusPending,
					SubmitDate:         "2014-03-02T02:22:12Z",
					UpdateDate:         "2014-03-01T21:12:57Z",
					Note:               "Sample activation",
					FMAActivationState: "steady",
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					FallbackInfo: &ActivationFallbackInfo{
						FastFallbackAttempted:      false,
						FallbackVersion:            10,
						CanFastFallback:            true,
						SteadyStateTime:            1506448172,
						FastFallbackExpirationTime: 1506451772,
						FastFallbackRecoveryState:  nil,
					},
				}},
				},
			},
		},
		"500 internal server error": {
			request: GetActivationsRequest{
				PropertyID: "prp_175780",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching activation",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching activation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: GetActivationsRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
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
			result, err := client.GetActivations(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapiGetActivation(t *testing.T) {
	tests := map[string]struct {
		request          GetActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetActivationResponse
		withError        error
	}{
		"200 OK": {
			request: GetActivationRequest{
				PropertyID:   "prp_175780",
				ActivationID: "atv_1696855",
				ContractID:   "ctr_1-1TJZFW",
				GroupID:      "grp_15166",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"accountId": "act_1-1TJZFB",
	"contractId": "ctr_1-1TJZFW",
	"groupId": "grp_15166",
	"activations": {
		"items": [
			{
				"activationId": "atv_1696985",
				"propertyName": "example.com",
				"propertyId": "prp_173136",
				"propertyVersion": 1,
				"network": "STAGING",
				"activationType": "ACTIVATE",
				"status": "PENDING",
				"submitDate": "2014-03-02T02:22:12Z",
				"updateDate": "2014-03-01T21:12:57Z",
				"note": "Sample activation",
				"fmaActivationState": "steady",
				"notifyEmails": [
					"you@example.com",
					"them@example.com"
				],
				"fallbackInfo": {
					"fastFallbackAttempted": false,
					"fallbackVersion": 10,
					"canFastFallback": true,
					"steadyStateTime": 1506448172,
					"fastFallbackExpirationTime": 1506451772,
					"fastFallbackRecoveryState": null
				}
			}
		]
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations/atv_1696855?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &GetActivationResponse{
				GetActivationsResponse: GetActivationsResponse{
					Response: Response{
						AccountID:  "act_1-1TJZFB",
						ContractID: "ctr_1-1TJZFW",
						GroupID:    "grp_15166",
					},
					Activations: ActivationsItems{Items: []*Activation{{
						ActivationID:       "atv_1696985",
						PropertyName:       "example.com",
						PropertyID:         "prp_173136",
						PropertyVersion:    1,
						Network:            ActivationNetworkStaging,
						ActivationType:     ActivationTypeActivate,
						Status:             ActivationStatusPending,
						SubmitDate:         "2014-03-02T02:22:12Z",
						UpdateDate:         "2014-03-01T21:12:57Z",
						Note:               "Sample activation",
						FMAActivationState: "steady",
						NotifyEmails: []string{
							"you@example.com",
							"them@example.com",
						},
						FallbackInfo: &ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            10,
							CanFastFallback:            true,
							SteadyStateTime:            1506448172,
							FastFallbackExpirationTime: 1506451772,
							FastFallbackRecoveryState:  nil,
						},
					}},
					},
				},
				Activation: &Activation{
					ActivationID:       "atv_1696985",
					PropertyName:       "example.com",
					PropertyID:         "prp_173136",
					PropertyVersion:    1,
					Network:            ActivationNetworkStaging,
					ActivationType:     ActivationTypeActivate,
					Status:             ActivationStatusPending,
					SubmitDate:         "2014-03-02T02:22:12Z",
					UpdateDate:         "2014-03-01T21:12:57Z",
					Note:               "Sample activation",
					FMAActivationState: "steady",
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					FallbackInfo: &ActivationFallbackInfo{
						FastFallbackAttempted:      false,
						FallbackVersion:            10,
						CanFastFallback:            true,
						SteadyStateTime:            1506448172,
						FastFallbackExpirationTime: 1506451772,
						FastFallbackRecoveryState:  nil,
					},
				},
			},
		},
		"activation not found": {
			request: GetActivationRequest{
				PropertyID:   "prp_175780",
				ActivationID: "atv_1696855",
				ContractID:   "ctr_1-1TJZFW",
				GroupID:      "grp_15166",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"accountId": "act_1-1TJZFB",
	"contractId": "ctr_1-1TJZFW",
	"groupId": "grp_15166",
	"activations": {
		"items": [
		]
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations/atv_1696855?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError:    ErrNotFound,
		},
		"500 internal server error": {
			request: GetActivationRequest{
				PropertyID:   "prp_175780",
				ActivationID: "atv_1696855",
				ContractID:   "ctr_1-1TJZFW",
				GroupID:      "grp_15166",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching activation",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations/atv_1696855?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching activation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: GetActivationRequest{
				ActivationID: "atv_1696855",
				ContractID:   "ctr_1-1TJZFW",
				GroupID:      "grp_15166",
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
			result, err := client.GetActivation(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapiCancelActivation(t *testing.T) {
	tests := map[string]struct {
		request          CancelActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CancelActivationResponse
		withError        error
	}{
		"200 OK": {
			request: CancelActivationRequest{
				PropertyID:   "prp_175780",
				ActivationID: "atv_1696855",
				ContractID:   "ctr_1-1TJZFW",
				GroupID:      "grp_15166",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"activations": {
		"items": [
			{
				"activationId": "atv_1696985",
				"propertyName": "example.com",
				"propertyId": "prp_173136",
				"propertyVersion": 1,
				"network": "STAGING",
				"activationType": "ACTIVATE",
				"status": "ABORTED",
				"submitDate": "2014-03-02T02:22:12Z",
				"updateDate": "2014-03-01T21:12:57Z",
				"note": "Sample activation",
				"fmaActivationState": "steady",
				"notifyEmails": [
					"you@example.com",
					"them@example.com"
				],
				"fallbackInfo": {
					"fastFallbackAttempted": false,
					"fallbackVersion": 10,
					"canFastFallback": true,
					"steadyStateTime": 1506448172,
					"fastFallbackExpirationTime": 1506451772,
					"fastFallbackRecoveryState": null
				}
			}
		]
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations/atv_1696855?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &CancelActivationResponse{
				Activations: ActivationsItems{Items: []*Activation{{
					ActivationID:       "atv_1696985",
					PropertyName:       "example.com",
					PropertyID:         "prp_173136",
					PropertyVersion:    1,
					Network:            ActivationNetworkStaging,
					ActivationType:     ActivationTypeActivate,
					Status:             ActivationStatusAborted,
					SubmitDate:         "2014-03-02T02:22:12Z",
					UpdateDate:         "2014-03-01T21:12:57Z",
					Note:               "Sample activation",
					FMAActivationState: "steady",
					NotifyEmails: []string{
						"you@example.com",
						"them@example.com",
					},
					FallbackInfo: &ActivationFallbackInfo{
						FastFallbackAttempted:      false,
						FallbackVersion:            10,
						CanFastFallback:            true,
						SteadyStateTime:            1506448172,
						FastFallbackExpirationTime: 1506451772,
						FastFallbackRecoveryState:  nil,
					},
				}},
				},
			},
		},
		"500 internal server error": {
			request: CancelActivationRequest{
				PropertyID:   "prp_175780",
				ActivationID: "atv_1696855",
				ContractID:   "ctr_1-1TJZFW",
				GroupID:      "grp_15166",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting activation",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780/activations/atv_1696855?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting activation",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: CancelActivationRequest{
				ActivationID: "atv_1696855",
				ContractID:   "ctr_1-1TJZFW",
				GroupID:      "grp_15166",
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
			result, err := client.CancelActivation(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
