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

func TestListPolicyActivations(t *testing.T) {
	tests := map[string]struct {
		parameters       ListPolicyActivationsRequest
		uri              string
		responseStatus   int
		responseBody     string
		expectedResponse []PolicyActivation
		withError        error
	}{
		"200 staging ok": {
			parameters: ListPolicyActivationsRequest{
				PropertyName: "www.rc-cloudlet.com",
				Network:      "staging",
				PolicyID:     1234,
			},
			uri:            "/cloudlets/api/v2/policies/1234/activations?network=staging&propertyName=www.rc-cloudlet.com",
			responseStatus: http.StatusOK,
			responseBody: `[
				{
					"serviceVersion": null,
					"apiVersion": "2.0",
					"network": "staging",
					"policyInfo": {
						"policyId": 2962,
						"name": "RequestControlPolicy",
						"version": 1,
						"status": "active",
						"statusDetail": "File successfully deployed to Akamai's network",
						"activationDate": 1427428800000,
						"activatedBy": "jsmith"
					},
					"propertyInfo": {
						"name": "www.rc-cloudlet.com",
						"version": 0,
						"groupId": 40498,
						"status": "inactive",
						"activatedBy": null,
						"activationDate": 0
					}
				}
			]`,
			expectedResponse: []PolicyActivation{{
				APIVersion: "2.0",
				Network:    PolicyActivationNetworkStaging,
				PropertyInfo: PropertyInfo{
					Name:           "www.rc-cloudlet.com",
					Version:        0,
					GroupID:        40498,
					Status:         PolicyActivationStatusInactive,
					ActivatedBy:    "",
					ActivationDate: 0,
				},
				PolicyInfo: PolicyInfo{
					PolicyID:       2962,
					Name:           "RequestControlPolicy",
					Version:        1,
					Status:         PolicyActivationStatusActive,
					StatusDetail:   "File successfully deployed to Akamai's network",
					ActivatedBy:    "jsmith",
					ActivationDate: 1427428800000,
				},
			}},
		},
		"empty Network should not appear in uri query": {
			parameters: ListPolicyActivationsRequest{
				PropertyName: "www.rc-cloudlet.com",
				Network:      "",
				PolicyID:     1234,
			},
			responseBody:     `[]`,
			expectedResponse: []PolicyActivation{},
			responseStatus:   http.StatusOK,
			uri:              "/cloudlets/api/v2/policies/1234/activations?propertyName=www.rc-cloudlet.com",
		},
		"empty PropertyName should not appear in uri query": {
			parameters: ListPolicyActivationsRequest{
				PropertyName: "",
				Network:      "staging",
				PolicyID:     1234,
			},
			responseBody:     `[]`,
			expectedResponse: []PolicyActivation{},
			responseStatus:   http.StatusOK,
			uri:              "/cloudlets/api/v2/policies/1234/activations?network=staging",
		},
		"not valid network": {
			parameters: ListPolicyActivationsRequest{
				PropertyName: "www.rc-cloudlet.com",
				Network:      "not valid",
				PolicyID:     1234,
			},
			withError: ErrStructValidation,
		},
		"404 not found": {
			parameters: ListPolicyActivationsRequest{
				PropertyName: "www.rc-cloudlet.com",
				Network:      "staging",
				PolicyID:     1234,
			},
			responseStatus: http.StatusNotFound,
			uri:            "/cloudlets/api/v2/policies/1234/activations?network=staging&propertyName=www.rc-cloudlet.com",
			withError:      &Error{StatusCode: 404, Title: "Failed to unmarshal error body. Cloudlets API failed. Check details for more information."},
		},
		"500 server error": {
			parameters: ListPolicyActivationsRequest{
				PropertyName: "www.rc-cloudlet.com",
				Network:      "staging",
				PolicyID:     1234,
			},
			responseStatus: http.StatusInternalServerError,
			uri:            "/cloudlets/api/v2/policies/1234/activations?network=staging&propertyName=www.rc-cloudlet.com",
			withError:      &Error{StatusCode: 500, Title: "Failed to unmarshal error body. Cloudlets API failed. Check details for more information."},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.uri, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListPolicyActivations(context.Background(), test.parameters)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestActivatePolicyVersion(t *testing.T) {
	tests := map[string]struct {
		parameters         ActivatePolicyVersionRequest
		responseStatus     int
		uri                string
		responseBody       string
		expectedResponse   []PolicyActivation
		expectedActivation PolicyVersionActivation
		withError          func(*testing.T, error)
	}{
		"200 Policy version activation": {
			responseStatus: http.StatusOK,
			parameters: ActivatePolicyVersionRequest{
				PolicyID: 1234,
				Version:  1,
				PolicyVersionActivation: PolicyVersionActivation{
					Network:                 PolicyActivationNetworkStaging,
					AdditionalPropertyNames: []string{"www.rc-cloudlet.com"},
				},
			},
			responseBody: `[
				{
					"serviceVersion": null,
					"apiVersion": "2.0",
					"network": "staging",
					"policyInfo": {
						"policyId": 2962,
						"name": "RequestControlPolicy",
						"version": 1,
						"status": "pending",
						"statusDetail": "initial",
						"activationDate": 1427428800000,
						"activatedBy": "jsmith"
					},
					"propertyInfo": {
						"name": "www.rc-cloudlet.com",
						"version": 0,
						"groupId": 40498,
						"status": "inactive",
						"activatedBy": null,
						"activationDate": 0
					}
				}
			]`,
			expectedResponse: []PolicyActivation{
				{
					APIVersion: "2.0",
					Network:    PolicyActivationNetworkStaging,
					PropertyInfo: PropertyInfo{
						Name:           "www.rc-cloudlet.com",
						Version:        0,
						GroupID:        40498,
						Status:         PolicyActivationStatusInactive,
						ActivatedBy:    "",
						ActivationDate: 0,
					},
					PolicyInfo: PolicyInfo{
						PolicyID:       2962,
						Name:           "RequestControlPolicy",
						Version:        1,
						Status:         PolicyActivationStatusPending,
						StatusDetail:   "initial",
						ActivatedBy:    "jsmith",
						ActivationDate: 1427428800000,
					},
				},
			},
		},
		"any request validation error": {
			parameters: ActivatePolicyVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"any kind of server error": {
			responseStatus: http.StatusInternalServerError,
			parameters: ActivatePolicyVersionRequest{
				PolicyID: 1234,
				Version:  1,
				PolicyVersionActivation: PolicyVersionActivation{
					Network:                 PolicyActivationNetworkStaging,
					AdditionalPropertyNames: []string{"www.rc-cloudlet.com"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrActivatePolicyVersion), "want: %s; got: %s", ErrActivatePolicyVersion, err)
			},
		},
		"property name not existing": {
			responseStatus: http.StatusBadRequest,
			parameters: ActivatePolicyVersionRequest{
				PolicyID: 1234,
				Version:  1,
				PolicyVersionActivation: PolicyVersionActivation{
					Network:                 PolicyActivationNetworkStaging,
					AdditionalPropertyNames: []string{"www.rc-cloudlet.com"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Requested propertyName \\\"XYZ\\\" does not exist", "want: %s; got: %s", ErrActivatePolicyVersion, err)
				assert.True(t, errors.Is(err, ErrActivatePolicyVersion), "want: %s; got: %s", ErrActivatePolicyVersion, err)
			},
			responseBody: `
				{
					"detail": "Requested propertyName \"XYZ\" does not exist",
					"errorCode": -1,
					"errorMessage": "Requested property Name \"XYZ\" does not exist",
					"instance": "s8dsf8sf8df8",
					"stackTrace": "java.lang.IllegalArgumentException: Requested property Name \"XYZ\" does not exist\n\tat com.akamai..."
				}
			`,
		},
		"validation errors": {
			responseStatus: http.StatusBadRequest,
			parameters: ActivatePolicyVersionRequest{
				PolicyID: 1234,
				Version:  1,
				PolicyVersionActivation: PolicyVersionActivation{
					Network:                 "",
					AdditionalPropertyNames: []string{},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "RequestBody.Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "RequestBody.AdditionalPropertyNames: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ActivatePolicyVersion(context.Background(), test.parameters)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
