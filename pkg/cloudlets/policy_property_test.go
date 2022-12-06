package cloudlets

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPolicyProperties(t *testing.T) {
	tests := map[string]struct {
		policyID         int64
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse map[string]PolicyProperty
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			policyID:       11754,
			responseStatus: http.StatusOK,
			responseBody: `
				{
					"www.myproperty.com": {
						"groupId": 40498,
						"id": 179120478,
						"name": "www.myproperty.com",
						"newestVersion": {
							"activatedBy": "sjones",
							"activationDate": "2015-08-25",
							"cloudletsOrigins": {
								"clorigin2": {
									"id": "clorigin2",
									"hostname": "origin2.myproperty.com",
									"checksum": "0edc0bb1be7439248a77f48e806d2531",
									"type": "CUSTOMER"
								},
								"clorigin1": {
									"id": "clorigin1",
									"hostname": "origin1.myproperty.com",
									"checksum": "eefa90e680a1183725cfe2a1f00307c4",
									"type": "CUSTOMER"
								}
							},
							"version": 5,
							"referencedPolicies": [
								"fr_policy_1"
							]
						},
						"production": {
							"activatedBy": "jsmith",
							"activationDate": "2015-08-26",
							"cloudletsOrigins": {
								"clorigin2": {
									"id": "clorigin2",
									"hostname": "origin2.myproperty.com",
									"checksum": "0edc0bb1be7439248a77f48e806d2531",
									"type": "CUSTOMER"
								},
								"clorigin1": {
									"id": "clorigin1",
									"hostname": "origin1.myproperty.com",
									"checksum": "eefa90e680a1183725cfe2a1f00307c4",
									"type": "CUSTOMER"
								}
							},
							"version": 5,
							"referencedPolicies": [
								"fr_policy_1"
							]
						},
						"staging": {
							"activatedBy": "jsmith",
							"activationDate": "2015-08-25",
							"cloudletsOrigins": {
								"clorigin2": {
									"id": "clorigin2",
									"hostname": "origin2.myproperty.com",
									"checksum": "0edc0bb1be7439248a77f48e806d2531",
									"type": "CUSTOMER"
								},
								"clorigin1": {
									"id": "clorigin1",
									"hostname": "origin1.myproperty.com",
									"checksum": "eefa90e680a1183725cfe2a1f00307c4",
									"type": "CUSTOMER"
								}
							},
							"version": 5,
							"referencedPolicies": [
								"fr_policy_1"
							]
						}
					}
				}
			`,
			expectedPath: "/cloudlets/api/v2/policies/11754/properties",
			expectedResponse: map[string]PolicyProperty{
				"www.myproperty.com": {
					GroupID: 40498,
					ID:      179120478,
					Name:    "www.myproperty.com",
					NewestVersion: NetworkStatus{
						ActivatedBy:    "sjones",
						ActivationDate: "2015-08-25",
						CloudletsOrigins: map[string]CloudletsOrigin{
							"clorigin2": {
								OriginID: "clorigin2",
								Hostname: "origin2.myproperty.com",
								Checksum: "0edc0bb1be7439248a77f48e806d2531",
								Type:     "CUSTOMER",
							},
							"clorigin1": {
								OriginID: "clorigin1",
								Hostname: "origin1.myproperty.com",
								Checksum: "eefa90e680a1183725cfe2a1f00307c4",
								Type:     "CUSTOMER",
							},
						},
						Version: 5,
						ReferencedPolicies: []string{
							"fr_policy_1",
						},
					},
					Production: NetworkStatus{
						ActivatedBy:    "jsmith",
						ActivationDate: "2015-08-26",
						CloudletsOrigins: map[string]CloudletsOrigin{
							"clorigin2": {
								OriginID: "clorigin2",
								Hostname: "origin2.myproperty.com",
								Checksum: "0edc0bb1be7439248a77f48e806d2531",
								Type:     "CUSTOMER",
							},
							"clorigin1": {
								OriginID: "clorigin1",
								Hostname: "origin1.myproperty.com",
								Checksum: "eefa90e680a1183725cfe2a1f00307c4",
								Type:     "CUSTOMER",
							},
						},
						Version: 5,
						ReferencedPolicies: []string{
							"fr_policy_1",
						},
					},
					Staging: NetworkStatus{
						ActivatedBy:    "jsmith",
						ActivationDate: "2015-08-25",
						CloudletsOrigins: map[string]CloudletsOrigin{
							"clorigin2": {
								OriginID: "clorigin2",
								Hostname: "origin2.myproperty.com",
								Checksum: "0edc0bb1be7439248a77f48e806d2531",
								Type:     "CUSTOMER",
							},
							"clorigin1": {
								OriginID: "clorigin1",
								Hostname: "origin1.myproperty.com",
								Checksum: "eefa90e680a1183725cfe2a1f00307c4",
								Type:     "CUSTOMER",
							},
						},
						Version: 5,
						ReferencedPolicies: []string{
							"fr_policy_1",
						},
					},
				},
			},
		},
		"500 Internal Server Error": {
			policyID:       11754,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error making request",
					"status": 500
				}
			`,
			expectedPath: "/cloudlets/api/v2/policies/11754/properties",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
			result, err := client.GetPolicyProperties(context.Background(), GetPolicyPropertiesRequest{PolicyID: test.policyID})
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCloudlets_DeletePolicyProperty(t *testing.T) {
	tests := map[string]struct {
		params         DeletePolicyPropertyRequest
		responseStatus int
		withError      error
		expectedURL    string
	}{
		"ok deletion prod": {
			params: DeletePolicyPropertyRequest{
				PolicyID:   1234,
				PropertyID: 5678,
				Network:    PolicyActivationNetworkProduction,
			},
			responseStatus: http.StatusNoContent,
			expectedURL:    "/cloudlets/api/v2/policies/1234/properties/5678?async=true&network=prod",
		},
		"ok deletion staging": {
			params: DeletePolicyPropertyRequest{
				PolicyID:   1234,
				PropertyID: 5678,
				Network:    PolicyActivationNetworkStaging,
			},
			responseStatus: http.StatusNoContent,
			expectedURL:    "/cloudlets/api/v2/policies/1234/properties/5678?async=true&network=staging",
		},
		"ok deletion no network": {
			params: DeletePolicyPropertyRequest{
				PolicyID:   1234,
				PropertyID: 5678,
			},
			responseStatus: http.StatusNoContent,
			expectedURL:    "/cloudlets/api/v2/policies/1234/properties/5678?async=true",
		},
		"nok validation property": {
			params: DeletePolicyPropertyRequest{
				PolicyID: 1234,
				Network:  PolicyActivationNetworkProduction,
			},
			withError: ErrStructValidation,
		},
		"nok validation policy": {
			params: DeletePolicyPropertyRequest{
				PropertyID: 1234,
				Network:    PolicyActivationNetworkProduction,
			},
			withError: ErrStructValidation,
		},
		"internal server error": {
			params: DeletePolicyPropertyRequest{
				PolicyID:   1234,
				PropertyID: 5678,
				Network:    PolicyActivationNetworkProduction,
			},
			responseStatus: http.StatusInternalServerError,
			withError:      ErrDeletePolicyProperty,
			expectedURL:    "/cloudlets/api/v2/policies/1234/properties/5678?async=true&network=prod",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedURL, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
			}))
			client := mockAPIClient(t, mockServer)

			err := client.DeletePolicyProperty(context.Background(), test.params)

			if test.withError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
