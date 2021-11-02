package cloudlets

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPropertyActivations(t *testing.T) {
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
				Network:    VersionActivationNetworkStaging,
				PropertyInfo: PropertyInfo{
					Name:           "www.rc-cloudlet.com",
					Version:        0,
					GroupID:        40498,
					Status:         StatusInactive,
					ActivatedBy:    "",
					ActivationDate: 0,
				},
				PolicyInfo: PolicyInfo{
					PolicyID:       2962,
					Name:           "RequestControlPolicy",
					Version:        1,
					Status:         StatusActive,
					StatusDetail:   "File successfully deployed to Akamai's network",
					ActivatedBy:    "jsmith",
					ActivationDate: 1427428800000,
				},
			}},
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
			withError:      &Error{StatusCode: 404},
		},
		"500 server error": {
			parameters: ListPolicyActivationsRequest{
				PropertyName: "www.rc-cloudlet.com",
				Network:      "staging",
				PolicyID:     1234,
			},
			responseStatus: http.StatusInternalServerError,
			uri:            "/cloudlets/api/v2/policies/1234/activations?network=staging&propertyName=www.rc-cloudlet.com",
			withError:      &Error{StatusCode: 500},
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
		expectedActivation PolicyActivation
		withError          *regexp.Regexp
	}{
		"200 Policy version activation": {
			responseStatus: http.StatusOK,
			parameters: ActivatePolicyVersionRequest{
				PolicyID: 1234,
				Version:  1,
				RequestBody: ActivatePolicyVersionRequestBody{
					Network:                 VersionActivationNetworkStaging,
					AdditionalPropertyNames: []string{"www.rc-cloudlet.com"},
				},
			},
			responseBody: `{
    "revisionId": 11870,
    "policyId": 1001,
    "version": 2,
    "description": "v2",
    "createdBy": "jsmith",
    "createDate": 1427133784903,
    "lastModifiedBy": "sjones",
    "lastModifiedDate": 1427133805767,
    "activations": [],
    "matchRules": [],
    "rulesLocked": false
}`,
		},
		"any request validation error": {
			parameters: ActivatePolicyVersionRequest{},
			withError:  regexp.MustCompile(ErrStructValidation.Error()),
		},
		"any kind of server error": {
			responseStatus: http.StatusInternalServerError,
			parameters: ActivatePolicyVersionRequest{
				PolicyID: 1234,
				Version:  1,
				RequestBody: ActivatePolicyVersionRequestBody{
					Network:                 VersionActivationNetworkStaging,
					AdditionalPropertyNames: []string{"www.rc-cloudlet.com"},
				},
			},
			withError: regexp.MustCompile(ErrActivatePolicyVersion.Error()),
		},
		"property name not existing": {
			responseStatus: http.StatusBadRequest,
			parameters: ActivatePolicyVersionRequest{
				PolicyID: 1234,
				Version:  1,
				RequestBody: ActivatePolicyVersionRequestBody{
					Network:                 VersionActivationNetworkStaging,
					AdditionalPropertyNames: []string{"www.rc-cloudlet.com"},
				},
			},
			withError: regexp.MustCompile(`"Requested propertyName \\"XYZ\\" does not exist"`),
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
		"empty property names": {
			responseStatus: http.StatusBadRequest,
			parameters: ActivatePolicyVersionRequest{
				PolicyID: 1234,
				Version:  1,
				RequestBody: ActivatePolicyVersionRequestBody{
					Network:                 VersionActivationNetworkStaging,
					AdditionalPropertyNames: []string{},
				},
			},
			withError: regexp.MustCompile(`struct validation: RequestBody.AdditionalPropertyNames: cannot be blank`),
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
			err := client.ActivatePolicyVersion(context.Background(), test.parameters)
			if test.withError != nil {
				assert.True(t, test.withError.MatchString(err.Error()), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
