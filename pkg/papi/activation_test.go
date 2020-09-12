package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapi_GetActivation(t *testing.T) {
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
				AccountID:  "act_1-1TJZFB",
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
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
			withError: session.APIError{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching activation",
				StatusCode: http.StatusInternalServerError,
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
