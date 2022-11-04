package cloudlets

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLoadBalancerActivations(t *testing.T) {
	tests := map[string]struct {
		params           ListLoadBalancerActivationsRequest
		originID         string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []LoadBalancerActivation
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListLoadBalancerActivationsRequest{
				OriginID: "clorigin1",
			},
			responseStatus: http.StatusOK,
			responseBody: `
				[
					{
						"activatedBy": "jjones",
						"activatedDate": "2016-05-03T18:41:34.251Z",
						"network": "PRODUCTION",
						"originId": "clorigin1",
						"status": "active",
						"version": 1
					},
					{
						"activatedBy": "ejnovak",
						"activatedDate": "2016-04-07T18:41:34.461Z",
						"network": "STAGING",
						"originId": "clorigin1",
						"status": "active",
						"version": 2
					}
				]
			`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin1/activations",
			expectedResponse: []LoadBalancerActivation{
				{
					ActivatedBy:   "jjones",
					ActivatedDate: "2016-05-03T18:41:34.251Z",
					Network:       LoadBalancerActivationNetworkProduction,
					OriginID:      "clorigin1",
					Status:        LoadBalancerActivationStatusActive,
					Version:       1,
				},
				{
					ActivatedBy:   "ejnovak",
					ActivatedDate: "2016-04-07T18:41:34.461Z",
					Network:       LoadBalancerActivationNetworkStaging,
					OriginID:      "clorigin1",
					Status:        LoadBalancerActivationStatusActive,
					Version:       2,
				},
			},
		},
		"200 OK with optional params": {
			params: ListLoadBalancerActivationsRequest{
				OriginID:   "clorigin1",
				Network:    "prod",
				LatestOnly: true,
				PageSize:   tools.Int64Ptr(3),
				Page:       tools.Int64Ptr(1),
			},
			responseStatus: http.StatusOK,
			responseBody: `
				[
					{
						"activatedBy": "jjones",
						"activatedDate": "2016-05-03T18:41:34.251Z",
						"network": "PRODUCTION",
						"originId": "clorigin1",
						"status": "active",
						"version": 1
					},
					{
						"activatedBy": "ejnovak",
						"activatedDate": "2016-05-07T18:41:34.461Z",
						"network": "PRODUCTION",
						"originId": "clorigin1",
						"status": "deactivated",
						"version": 2
					}
				]
			`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin1/activations?latestOnly=true&network=prod&page=1&pageSize=3",
			expectedResponse: []LoadBalancerActivation{
				{
					ActivatedBy:   "jjones",
					ActivatedDate: "2016-05-03T18:41:34.251Z",
					Network:       LoadBalancerActivationNetworkProduction,
					OriginID:      "clorigin1",
					Status:        LoadBalancerActivationStatusActive,
					Version:       1,
				},
				{
					ActivatedBy:   "ejnovak",
					ActivatedDate: "2016-05-07T18:41:34.461Z",
					Network:       LoadBalancerActivationNetworkProduction,
					OriginID:      "clorigin1",
					Status:        LoadBalancerActivationStatusDeactivated,
					Version:       2,
				},
			},
		},
		"500 Internal Server Error": {
			params: ListLoadBalancerActivationsRequest{
				OriginID: "clorigin1",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error making request",
					"status": 500
				}
			`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin1/activations",
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
		"Validation Errors": {
			params: ListLoadBalancerActivationsRequest{
				Network: "PRODUCTION",
			},
			responseStatus: http.StatusInternalServerError,
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: value 'PRODUCTION' is invalid. Must be one of: 'staging', 'prod' or '' (empty)", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "OriginID: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
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
			result, err := client.ListLoadBalancerActivations(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestActivateLoadBalancerVersion(t *testing.T) {
	tests := map[string]struct {
		params           ActivateLoadBalancerVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *LoadBalancerActivation
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ActivateLoadBalancerVersionRequest{
				OriginID: "clorigin1",
				Async:    false,
				LoadBalancerVersionActivation: LoadBalancerVersionActivation{
					Network: LoadBalancerActivationNetworkProduction,
					DryRun:  false,
					Version: 1,
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
				{
					"activatedBy": "jjones",
					"activatedDate": "2016-04-07T18:41:34.251Z",
					"network": "PRODUCTION",
					"originId": "clorigin1",
					"status": "active",
					"dryrun": false,
					"version": 1
				}
			`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin1/activations?async=false",
			expectedResponse: &LoadBalancerActivation{
				ActivatedBy:   "jjones",
				ActivatedDate: "2016-04-07T18:41:34.251Z",
				Network:       LoadBalancerActivationNetworkProduction,
				OriginID:      "clorigin1",
				Status:        LoadBalancerActivationStatusActive,
				Version:       1,
				DryRun:        false,
			},
		},
		"500 Internal Server Error": {
			params: ActivateLoadBalancerVersionRequest{
				OriginID: "clorigin1",
				Async:    false,
				LoadBalancerVersionActivation: LoadBalancerVersionActivation{
					Network: LoadBalancerActivationNetworkStaging,
					DryRun:  false,
					Version: 1,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error making request",
					"status": 500
				}
			`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin1/activations?async=false",
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
		"Validation Errors": {
			params: ActivateLoadBalancerVersionRequest{
				OriginID: "",
				Async:    false,
				LoadBalancerVersionActivation: LoadBalancerVersionActivation{
					Network: "",
					DryRun:  false,
					Version: -1,
				},
			},
			responseStatus: http.StatusInternalServerError,
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Network: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "OriginID: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.Containsf(t, err.Error(), "Version: must be no less than 0", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
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
			result, err := client.ActivateLoadBalancerVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
