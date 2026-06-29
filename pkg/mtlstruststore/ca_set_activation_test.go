package mtlstruststore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivateCASetVersion(t *testing.T) {
	tests := map[string]struct {
		request             ActivateCASetVersionRequest
		expectedRequestBody string
		expectedPath        string
		responseStatus      int
		responseBody        string
		expectedResponse    *ActivateCASetVersionResponse
		withError           func(*testing.T, error)
	}{
		"202 Accepted": {
			request: ActivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activate",
			responseStatus:      http.StatusAccepted,
			responseBody: `
				{
					"activationId": 84707,
					"activationLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84707",
					"caSetId": "199",
					"caSetName": "test1",
					"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
					"createdBy": "jsmith@example.com",
					"createdDate": "2023-06-01T23:02:29.876589Z",
					"failureReason": null,
					"modifiedBy": null,
					"modifiedDate": null,
					"network": "STAGING",
					"activationStatus": "IN_PROGRESS",
					"activationType": "ACTIVATE",
					"percentComplete": 0,
					"validation": {
						"warnings": []
					},
					"version": 1,
					"versionLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1"
				}`,
			expectedResponse: &ActivateCASetVersionResponse{
				ActivationID:     84707,
				ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84707",
				CASetID:          "199",
				CASetName:        "test1",
				CASetLink:        "/mtls-edge-truststore/v2/ca-sets/199",
				CreatedBy:        "jsmith@example.com",
				CreatedDate:      test.NewTimeFromString(t, "2023-06-01T23:02:29.876589Z"),
				FailureReason:    nil,
				ModifiedBy:       nil,
				ModifiedDate:     nil,
				Network:          "STAGING",
				ActivationStatus: "IN_PROGRESS",
				ActivationType:   "ACTIVATE",
				PercentComplete:  0,
				Version:          1,
				VersionLink:      "/mtls-edge-truststore/v2/ca-sets/199/versions/1",
				Validation:       &Validation{Warnings: []Warning{}},
			},
		},
		"202 Accepted with warning": {
			request: ActivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activate",
			responseStatus:      http.StatusAccepted,
			responseBody: `
				{
					"activationId": 84707,
					"activationLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84707",
					"caSetId": "199",
					"caSetName": "test1",
					"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
					"createdBy": "jsmith@example.com",
					"createdDate": "2023-06-01T23:02:29.876589Z",
					"failureReason": null,
					"modifiedBy": null,
					"modifiedDate": null,
					"network": "STAGING",
					"activationStatus": "IN_PROGRESS",
					"activationType": "ACTIVATE",
					"percentComplete": 0,
					"validation": {
						"warnings": [
							{
								"detail": "A request was submitted to activate the CA set version on production without activating on staging first. While this is allowed, it is highly recommended to first deploy the CA set version on staging network and ensure that there are no issues in serving the traffic.",
								"title": "Request was submitted to activate the CA set version on production without activating on staging first.",
								"type": "/mtls-edge-truststore/error-types/activate-version-on-production-without-activating-on-staging"
							}
						]
					},
					"version": 1,
					"versionLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1"
				}`,
			expectedResponse: &ActivateCASetVersionResponse{
				ActivationID:     84707,
				ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84707",
				CASetID:          "199",
				CASetName:        "test1",
				CASetLink:        "/mtls-edge-truststore/v2/ca-sets/199",
				CreatedBy:        "jsmith@example.com",
				CreatedDate:      test.NewTimeFromString(t, "2023-06-01T23:02:29.876589Z"),
				FailureReason:    nil,
				ModifiedBy:       nil,
				ModifiedDate:     nil,
				Network:          "STAGING",
				ActivationStatus: "IN_PROGRESS",
				ActivationType:   "ACTIVATE",
				PercentComplete:  0,
				Version:          1,
				VersionLink:      "/mtls-edge-truststore/v2/ca-sets/199/versions/1",
				Validation: &Validation{
					Warnings: []Warning{
						{
							Detail: "A request was submitted to activate the CA set version on production without activating on staging first. While this is allowed, it is highly recommended to first deploy the CA set version on staging network and ensure that there are no issues in serving the traffic.",
							Title:  "Request was submitted to activate the CA set version on production without activating on staging first.",
							Type:   "/mtls-edge-truststore/error-types/activate-version-on-production-without-activating-on-staging",
						},
					},
				},
			},
		},
		"missing required request param - validation error": {
			request: ActivateCASetVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "activate ca set version failed: struct validation: CASetID: cannot be blank\nNetwork: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"invalid network - validation error": {
			request: ActivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: "foo",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "activate ca set version failed: struct validation: Network: value 'foo' is invalid. Must be one of: 'STAGING' or 'PRODUCTION'", err.Error())
			},
		},
		"409 cannot be activated": {
			request: ActivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkStaging,
			},
			expectedRequestBody: `{"network":"STAGING"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activate",
			responseStatus:      http.StatusConflict,
			responseBody: `
				{
    "contextInfo": {
        "caSetId": "199",
        "caSetName": "DXE-5498",
        "network": "STAGING",
        "version": 1
    },
    "detail": "CA set version with version 1 cannot be activated as it is already active on the STAGING network.",
    "instance": "/mtls-edge-truststore/error-types/ca-set-version-already-active-on-network/1916882c5af2b20b",
    "status": 409,
    "title": "CA set version cannot be activated as it is already active on the network.",
    "type": "/mtls-edge-truststore/error-types/ca-set-version-already-active-on-network"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetVersionIsActiveOnNetworkCannotBeActivated))
			},
		},
		"500 internal server error": {
			request: ActivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal-server-error",
    "title": "Internal Server Error",
    "detail": "Error processing request",
    "instance": "TestInstances",
    "status": 500
}`,
			expectedPath: "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activate",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
				if tc.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ActivateCASetVersion(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestDeactivateCASetVersion(t *testing.T) {
	tests := map[string]struct {
		request             DeactivateCASetVersionRequest
		expectedRequestBody string
		expectedPath        string
		responseStatus      int
		responseBody        string
		expectedResponse    *DeactivateCASetVersionResponse
		withError           func(*testing.T, error)
	}{
		"202 Accepted": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			responseStatus:      http.StatusAccepted,
			responseBody: `
				{
					"activationId": 84707,
					"activationLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84707",
					"caSetId": "199",
					"caSetName": "test1",
					"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
					"createdBy": "jsmith@example.com",
					"createdDate": "2023-06-01T23:02:29.876589Z",
					"failureReason": null,
					"modifiedBy": null,
					"modifiedDate": null,
					"network": "STAGING",
					"activationStatus": "IN_PROGRESS",
					"activationType": "DEACTIVATE",
					"percentComplete": 0,
					"validation": {
						"warnings": []
					},
					"version": 1,
					"versionLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1"
				}`,
			expectedResponse: &DeactivateCASetVersionResponse{
				ActivationID:     84707,
				ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84707",
				CASetID:          "199",
				CASetName:        "test1",
				CASetLink:        "/mtls-edge-truststore/v2/ca-sets/199",
				CreatedBy:        "jsmith@example.com",
				CreatedDate:      test.NewTimeFromString(t, "2023-06-01T23:02:29.876589Z"),
				FailureReason:    nil,
				ModifiedBy:       nil,
				ModifiedDate:     nil,
				Network:          "STAGING",
				ActivationStatus: "IN_PROGRESS",
				ActivationType:   "DEACTIVATE",
				PercentComplete:  0,
				Version:          1,
				VersionLink:      "/mtls-edge-truststore/v2/ca-sets/199/versions/1",
				Validation:       &Validation{Warnings: []Warning{}},
			},
		},
		"202 Accepted with warning": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			responseStatus:      http.StatusAccepted,
			responseBody: `
				{
					"activationId": 84707,
					"activationLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84707",
					"caSetId": "199",
					"caSetName": "test1",
					"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
					"createdBy": "jsmith@example.com",
					"createdDate": "2023-06-01T23:02:29.876589Z",
					"failureReason": null,
					"modifiedBy": null,
					"modifiedDate": null,
					"network": "STAGING",
					"activationStatus": "IN_PROGRESS",
					"activationType": "DEACTIVATE",
					"percentComplete": 0,
					"validation": {
						"warnings": [
							{
								"detail": "Example warning detail",
								"title": "Example warning title",
								"type": "/mtls-edge-truststore/error-types/example-type"
							}
						]
					},
					"version": 1,
					"versionLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1"
				}`,
			expectedResponse: &DeactivateCASetVersionResponse{
				ActivationID:     84707,
				ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84707",
				CASetID:          "199",
				CASetName:        "test1",
				CASetLink:        "/mtls-edge-truststore/v2/ca-sets/199",
				CreatedBy:        "jsmith@example.com",
				CreatedDate:      test.NewTimeFromString(t, "2023-06-01T23:02:29.876589Z"),
				FailureReason:    nil,
				ModifiedBy:       nil,
				ModifiedDate:     nil,
				Network:          "STAGING",
				ActivationStatus: "IN_PROGRESS",
				ActivationType:   "DEACTIVATE",
				PercentComplete:  0,
				Version:          1,
				VersionLink:      "/mtls-edge-truststore/v2/ca-sets/199/versions/1",
				Validation: &Validation{
					Warnings: []Warning{
						{
							Detail: "Example warning detail",
							Title:  "Example warning title",
							Type:   "/mtls-edge-truststore/error-types/example-type",
						},
					},
				},
			},
		},
		"Error Response - CA Set is associated associated with enrollments/hostnames (Commercial)": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			responseStatus:      http.StatusConflict,
			responseBody: `
				{
    "contextInfo": {
        "associations": {
            "enrollments": [
                {
                    "cn": "some.domain.com",
                    "enrollmentId": 12345,
                    "enrollmentLink": "/cps/v2/enrollments/12345",
                    "productionSlots": [],
                    "stagingSlots": [
                        123466
                    ]
                },
                {
                    "cn": "some.domain.com",
                    "enrollmentId": 12345,
                    "enrollmentLink": "/cps/v2/enrollments/12345",
                    "productionSlots": [],
                    "stagingSlots": [
                        123456
                    ]
                }
            ],
            "properties": [
                {
                    "hostnames": [
                        {
                            "hostname": "some.domain.com",
                            "network": "PRODUCTION"
                        },
                        {
                            "hostname": "some.domain.com",
                            "network": "STAGING"
                        }
                    ],
                    "propertyId": "1234567"
                }
            ]
        },
        "caSetId": "132965",
        "caSetName": "jsmith-mets-tpc-delete-test-3"
    },
    "detail": "CA set with caSetId 132965 is linked to several Certificate Provisioning System enrollments and Property Manager hostnames. You need to unlink the CA set from the enrollments and hostnames to proceed. See accompanying data for enrollment and hostname details.",
    "instance": "/mtls-edge-truststore/error-types/ca-set-in-use-by-both-enrollments-and-hostnames/f2869f68451b8408",
    "status": 409,
    "title": "CA set is linked to both enrollments and hostnames.",
    "type": "/mtls-edge-truststore/error-types/ca-set-in-use-by-both-enrollments-and-hostnames"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetInUseByBothEnrollmentsAndHostnames))
			},
		},
		"Error Response - CA Set is associated with a hostname in Property Manager (Defense Edge)": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			responseStatus:      http.StatusConflict,
			responseBody: `
				{
					"contextInfo" : {
						"associations" : {
							"properties" : [ {
								"hostnames" : [ {
									"hostname" : "example-3.com"
								} ],
								"propertyId" : "2"
							} ]
						},
						"caSetId" : 1,
						"caSetName" : "foo",
						"network" : "PRODUCTION"
					},
					"detail" : "CA set cannot be deactivated CA set with caSetId 1 links to several Property Manager hostnames. You need to unlink the CA set from the hostnames to proceed. See accompanying response data for hostname details.",
					"status" : 409,
					"title" : "CA set is linked to hostnames.",
					"type" : "/mtls-edge-truststore/error-types/ca-set-bound-to-hostname"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetBoundToHostname))
			},
		},
		"Error Response - Another activation  is associated with the CA set version": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			responseStatus:      http.StatusConflict,
			responseBody: `
				{
					"contextInfo" : {
						"activationLink" : "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/1",
						"activationType" : "DEACTIVATE",
						"caSetId" : 199,
						"caSetName" : "foo",
						"network" : "PRODUCTION",
						"version" : 1
					},
					"detail" : "CA set with caSetId 1 and version 1 cannot be DEACTIVATED as another activation request is in progress for the CA set on the PRODUCTION network. Hypermedia link to the activation is attached.",
					"status" : 409,
					"title" : "Another activation request is in progress in the CA set.",
					"type" : "/mtls-edge-truststore/error-types/another-activation-request-in-progress-in-the-ca-set"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrAnotherActivationInProgress))
			},
		},
		"Error Response - Another deactivation is associated with the CA set version": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			responseStatus:      http.StatusConflict,
			responseBody: `
				{
					"contextInfo" : {
						"activationLink" : "/mtls-edge-truststore/v2/ca-sets/1/versions/1/activations/1",
						"activationType" : "DEACTIVATE",
						"caSetId" : 1,
						"caSetName" : "caSetName-1",
						"network" : "STAGING",
						"version" : 1
					},
					"detail" : "CA set version with version 1 cannot be DEACTIVATED as another deactivation request is in progress for the CA set with caSetId 1 on the STAGING network. Hypermedia link to the deactivation is attached.",
					"status" : 409,
					"title" : "Another deactivation request is in progress in the CA set.",
					"type" : "/mtls-edge-truststore/error-types/another-deactivation-request-in-progress-in-the-ca-set"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrAnotherDeactivationInProgress))
			},
		},
		"Error Response - Version is not activated on the network": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			expectedRequestBody: `{"network":"PRODUCTION"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			responseStatus:      http.StatusConflict,
			responseBody: `
				{
					"contextInfo" : {
						"caSetId" : 199,
						"caSetName" : "foo",
						"network" : "STAGING",
						"version" : 1
					},
					"detail" : "CA set version with version 1 cannot be deactivated as it is not active on the STAGING network.",
					"status" : 409,
					"title" : "CA set version cannot be deactivated as it is not active on the network.",
					"type" : "/mtls-edge-truststore/error-types/ca-set-version-not-active-on-network"
				}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetVersionNotActiveOnNetwork))
			},
		},
		"409 cannot be deactivated": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkStaging,
			},
			expectedRequestBody: `{"network":"STAGING"}`,
			expectedPath:        "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			responseStatus:      http.StatusConflict,
			responseBody: `
				{
    "contextInfo": {
        "caSetId": "199",
        "caSetName": "DXE-5498",
        "network": "STAGING",
        "version": 1
    },
    "detail": "CA set version with version 1 cannot be deactivated as it is not active on the STAGING network.",
    "instance": "/mtls-edge-truststore/error-types/ca-set-version-not-active-on-network/0d0a2fb6c1c58594",
    "status": 409,
    "title": "CA set version cannot be deactivated as it is not active on the network.",
    "type": "/mtls-edge-truststore/error-types/ca-set-version-not-active-on-network"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetVersionNotActiveOnNetworkCannotBeDeactivated))
			},
		},
		"missing required request param - validation error": {
			request: DeactivateCASetVersionRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "deactivate ca set version failed: struct validation: CASetID: cannot be blank\nNetwork: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"invalid network - validation error": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: "foo",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "deactivate ca set version failed: struct validation: Network: value 'foo' is invalid. Must be one of: 'STAGING' or 'PRODUCTION'", err.Error())
			},
		},
		"500 internal server error": {
			request: DeactivateCASetVersionRequest{
				CASetID: "199",
				Version: 1,
				Network: ActivationNetworkProduction,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal-server-error",
    "title": "Internal Server Error",
    "detail": "Error processing request",
    "instance": "TestInstances",
    "status": 500
}`,
			expectedPath: "/mtls-edge-truststore/v2/ca-sets/199/versions/1/deactivate",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
				if tc.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, tc.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.DeactivateCASetVersion(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetCASetVersionActivation(t *testing.T) {
	tests := map[string]struct {
		params           GetCASetVersionActivationRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		responseHeaders  map[string]string
		expectedResponse *GetCASetVersionActivationResponse
		withError        func(*testing.T, error)
	}{
		"202 Accepted": {
			params: GetCASetVersionActivationRequest{
				CASetID:      "1000",
				Version:      1,
				ActivationID: 84572,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/1/activations/84572",
			responseStatus: http.StatusAccepted,
			responseHeaders: map[string]string{
				"Retry-After": "Tue, 10 Jan 2023 11:05:00 GMT",
			},
			responseBody: `
				{
					 "activationId": 84572,
					 "activationLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/1000/activations/84572",
					 "caSetId": "1000",
					 "caSetName": "test1",
					 "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1000",  
					 "percentComplete": 0,
					 "validation": null,
					 "version": 1,
					 "versionLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/1",
					 "network": "STAGING",
					 "activationType": "ACTIVATE",
					 "activationStatus": "IN_PROGRESS",
					 "createdDate": "2023-01-10T11:00:00.494383Z",
					 "createdBy": "someone",
					 "failureReason": null,
					 "modifiedDate": "2023-01-10T12:00:00.771298Z",
					 "modifiedBy": "someone"
				}`,
			expectedResponse: &GetCASetVersionActivationResponse{
				ActivationID:     84572,
				ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/1000/activations/84572",
				CASetID:          "1000",
				CASetName:        "test1",
				CASetLink:        "/mtls-edge-truststore/v2/ca-sets/1000",
				PercentComplete:  0,
				Version:          1,
				VersionLink:      "/mtls-edge-truststore/v2/ca-sets/1000/versions/1",
				Network:          "STAGING",
				ActivationType:   "ACTIVATE",
				ActivationStatus: "IN_PROGRESS",
				CreatedDate:      test.NewTimeFromString(t, "2023-01-10T11:00:00.494383Z"),
				CreatedBy:        "someone",
				FailureReason:    nil,
				ModifiedDate:     ptr.To(test.NewTimeFromString(t, "2023-01-10T12:00:00.771298Z")),
				ModifiedBy:       ptr.To("someone"),
				RetryAfter:       test.NewGMTTimeFromString(t, "Tue, 10 Jan 2023 11:05:00 GMT"),
			},
		},
		"200 OK": {
			params: GetCASetVersionActivationRequest{
				CASetID:      "1000",
				Version:      1,
				ActivationID: 84572,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/1/activations/84572",
			responseStatus: http.StatusOK,
			responseBody: `
				{
					 "activationId": 84572,
					 "activationLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/1000/activations/84572",
					 "caSetId": "1000",
					 "caSetName": "test1",
					 "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1000",
					 "percentComplete": 100,
					 "version": 1,
					 "versionLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/1",
					 "network": "STAGING",
					 "activationType": "ACTIVATE",
					 "activationStatus": "COMPLETE",
					 "createdDate": "2023-01-10T11:00:00.494383Z",
					 "createdBy": "someone",
					 "failureReason": null,
					 "modifiedDate": "2023-01-10T12:00:00.771298Z",
					 "modifiedBy": "someone"
				}`,
			expectedResponse: &GetCASetVersionActivationResponse{
				ActivationID:     84572,
				ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/1000/activations/84572",
				CASetID:          "1000",
				CASetName:        "test1",
				CASetLink:        "/mtls-edge-truststore/v2/ca-sets/1000",
				PercentComplete:  100,
				Version:          1,
				VersionLink:      "/mtls-edge-truststore/v2/ca-sets/1000/versions/1",
				Network:          "STAGING",
				ActivationType:   "ACTIVATE",
				ActivationStatus: "COMPLETE",
				CreatedDate:      test.NewTimeFromString(t, "2023-01-10T11:00:00.494383Z"),
				CreatedBy:        "someone",
				FailureReason:    nil,
				ModifiedDate:     ptr.To(test.NewTimeFromString(t, "2023-01-10T12:00:00.771298Z")),
				ModifiedBy:       ptr.To("someone"),
			},
		},
		"missing required params - validation error": {
			params: GetCASetVersionActivationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get ca set version activation failed: struct validation: ActivationID: cannot be blank\nCASetID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"404 ca set not found - custom error check": {
			params: GetCASetVersionActivationRequest{
				CASetID:      "10",
				Version:      1,
				ActivationID: 84572,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/10/versions/1/activations/84572",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/mtls-edge-truststore/error-types/ca-set-not-found",
				"title": "CA set not found.",
				"status": 404,
				"detail": "Cannot get activation or deactivation status as the CA set with caSetId 10 is not found.",
				"contextInfo": {
					"caSetId": "10"
				}
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetNotFound))
			},
		},
		"404 ca set version not found - custom error check": {
			params: GetCASetVersionActivationRequest{
				CASetID:      "2",
				Version:      12,
				ActivationID: 84572,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/2/versions/12/activations/84572",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/mtls-edge-truststore/error-types/ca-set-version-not-found",
				"title": "CA set version not found",
				"status": 404,
				"detail": "Cannot get activation or deactivation status as the version 12 is not found in the CA set under caSetName test1.",
				"contextInfo": {
					"caSetName": "test1",
					"caSetId": "2",
					"version": 12
				}
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetVersionNotFound))
			},
		},
		"404 ca set activation not found - custom error check": {
			params: GetCASetVersionActivationRequest{
				CASetID:      "2",
				Version:      1,
				ActivationID: 2,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/2/versions/1/activations/2",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"contextInfo" : {
					"activationId" : 2,
					"caSetId" : "2",
					"caSetName" : "caSetName-1",
					"version" : 1
				},
				"detail" : "Cannot get activation or deactivation status as the activation or deactivation request with activationId 2 is not found in CA set version with version 1 under CA set with caSetId 2.",
				"status" : 404,
				"title" : "Activation or Deactivation Request is not found.",
				"type" : "/mtls-edge-truststore/error-types/activation-or-deactivation-request-not-found"
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetActivationNotFound))
			},
		},
		"500 internal server error": {
			params: GetCASetVersionActivationRequest{
				CASetID:      "199",
				Version:      1,
				ActivationID: 84572,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal-server-error",
  "title": "Internal Server Error",
  "detail": "Error processing request",
  "instance": "TestInstances",
  "status": 500
}`,
			expectedPath: "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/84572",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				if len(tc.responseHeaders) > 0 {
					for header, value := range tc.responseHeaders {
						w.Header().Set(header, value)
					}
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCASetVersionActivation(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListCASetVersionActivations(t *testing.T) {
	tests := map[string]struct {
		request          ListCASetVersionActivationsRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *ListCASetVersionActivationsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: ListCASetVersionActivationsRequest{
				CASetID: "1000",
				Version: 2,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/2/activations",
			responseStatus: http.StatusOK,
			responseBody: `{
			   "activations": [
				  {
					 "activationId": 84571,
					 "activationLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/2/activations/84571",
					 "caSetId": "1000",
					 "caSetName": "test1",
					 "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1000",
					 "percentComplete": 100,
					 "validation": null,
					 "version": 2,
					 "versionLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/3",
					 "network": "STAGING",
					 "activationType": "DEACTIVATE",
					 "activationStatus": "COMPLETE",
					 "createdDate": "2023-01-11T11:00:00.385674Z",
					 "createdBy": "someone",
					 "failureReason": null,
					 "modifiedDate": "2023-01-11T12:00:00.358292Z",
					 "modifiedBy": "someone"
				  }
			   ]
			}`,
			expectedResponse: &ListCASetVersionActivationsResponse{
				Activations: []ActivateCASetVersionResponse{
					{
						ActivationID:     84571,
						ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/2/activations/84571",
						CASetID:          "1000",
						CASetName:        "test1",
						CASetLink:        "/mtls-edge-truststore/v2/ca-sets/1000",
						PercentComplete:  100,
						Version:          2,
						VersionLink:      "/mtls-edge-truststore/v2/ca-sets/1000/versions/3",
						Network:          "STAGING",
						ActivationType:   "DEACTIVATE",
						ActivationStatus: "COMPLETE",
						CreatedDate:      test.NewTimeFromString(t, "2023-01-11T11:00:00.385674Z"),
						CreatedBy:        "someone",
						ModifiedDate:     ptr.To(test.NewTimeFromString(t, "2023-01-11T12:00:00.358292Z")),
						ModifiedBy:       ptr.To("someone"),
					},
				},
			},
		},
		"404 ca set not found - custom error check": {
			request: ListCASetVersionActivationsRequest{
				CASetID: "1000",
				Version: 2,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/2/activations",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/mtls-edge-truststore/error-types/ca-set-not-found",
				"title": "CA set not found.",
				"status": 404,
				"detail": "Cannot get CA set activations as the CA set with caSetId 1000 is not found.",
				"contextInfo": {
					"caSetId": "1000"
				}
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetNotFound))
			},
		},
		"404 ca set version not found - custom error check": {
			request: ListCASetVersionActivationsRequest{
				CASetID: "1000",
				Version: 2,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/2/activations",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/mtls-edge-truststore/error-types/ca-set-version-not-found",
				"title": "CA set version not found",
				"status": 404,
				"detail": "Cannot get CA set activations as the CA set version with version 2 is not found in the CA set under caSetName foo.",
				"contextInfo": {
					"caSetName": "foo",
					"caSetId": "1000",
					"version": 2
				}
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetVersionNotFound))
			},
		},

		"missing id - validation error": {
			request: ListCASetVersionActivationsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca set version activations failed: struct validation: CASetID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			request: ListCASetVersionActivationsRequest{
				CASetID: "199",
				Version: 2,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199/versions/2/activations",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal-server-error",
  "title": "Internal Server Error",
  "detail": "Error processing request",
  "instance": "TestInstances",
  "status": 500
}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListCASetVersionActivations(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListCASetActivations(t *testing.T) {
	tests := map[string]struct {
		request          ListCASetActivationsRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *ListCASetActivationsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: ListCASetActivationsRequest{
				CASetID: "1000",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1000/activations",
			responseStatus: http.StatusOK,
			responseBody: `{
			   "activations": [
				  {
					 "activationId": 84572,
					 "activationLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/1/activations/84572",
					 "caSetId": "1000",
					 "caSetName": "test1",
					 "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1000",
					 "percentComplete": 100,
					 "validation": null,
					 "version": 1,
					 "versionLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/2",
					 "network": "STAGING",
					 "activationType": "ACTIVATE",
					 "activationStatus": "FAILED",
					 "failureReason": null,
					 "createdDate": "2023-01-10T11:00:00.494383Z",
					 "createdBy": "someone",
					 "modifiedDate": "2023-01-10T12:00:00.771298Z",
					 "modifiedBy": "someone"
				  },
				  {
					 "activationId": 84571,
					 "activationLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/2/activations/84571",
					 "caSetId": "1000",
					 "caSetName": "test1",
					 "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1000",
					 "percentComplete": 0,
					 "validation": null,
					 "version": 2,
					 "versionLink": "/mtls-edge-truststore/v2/ca-sets/1000/versions/3",
					 "network": "STAGING",
					 "activationType": "DEACTIVATE",
					 "activationStatus": "COMPLETE",
					 "failureReason": "test failure reason",
					 "createdDate": "2023-01-11T11:00:00.385674Z",
					 "createdBy": "someone",
					 "modifiedDate": "2023-01-11T12:00:00.358292Z",
					 "modifiedBy": "someone"
				  }
			   ]
			}`,
			expectedResponse: &ListCASetActivationsResponse{
				Activations: []ActivateCASetVersionResponse{
					{
						ActivationID:     84572,
						ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/1/activations/84572",
						CASetID:          "1000",
						CASetName:        "test1",
						CASetLink:        "/mtls-edge-truststore/v2/ca-sets/1000",
						PercentComplete:  100,
						Version:          1,
						VersionLink:      "/mtls-edge-truststore/v2/ca-sets/1000/versions/2",
						Network:          "STAGING",
						ActivationType:   "ACTIVATE",
						ActivationStatus: "FAILED",
						FailureReason:    nil,
						CreatedDate:      test.NewTimeFromString(t, "2023-01-10T11:00:00.494383Z"),
						CreatedBy:        "someone",
						ModifiedDate:     ptr.To(test.NewTimeFromString(t, "2023-01-10T12:00:00.771298Z")),
						ModifiedBy:       ptr.To("someone"),
					},
					{
						ActivationID:     84571,
						ActivationLink:   "/mtls-edge-truststore/v2/ca-sets/1000/versions/2/activations/84571",
						CASetID:          "1000",
						CASetName:        "test1",
						CASetLink:        "/mtls-edge-truststore/v2/ca-sets/1000",
						PercentComplete:  0,
						Version:          2,
						VersionLink:      "/mtls-edge-truststore/v2/ca-sets/1000/versions/3",
						Network:          "STAGING",
						ActivationType:   "DEACTIVATE",
						ActivationStatus: "COMPLETE",
						FailureReason:    ptr.To("test failure reason"),
						CreatedDate:      test.NewTimeFromString(t, "2023-01-11T11:00:00.385674Z"),
						CreatedBy:        "someone",
						ModifiedDate:     ptr.To(test.NewTimeFromString(t, "2023-01-11T12:00:00.358292Z")),
						ModifiedBy:       ptr.To("someone"),
					},
				},
			},
		},
		"404 ca set not found - custom error check": {
			request: ListCASetActivationsRequest{
				CASetID: "1000",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1000/activations",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "/mtls-edge-truststore/error-types/ca-set-not-found",
				"title": "CA set not found.",
				"status": 404,
				"detail": "Cannot get CA set activations as the CA set with caSetId 1000 is not found.",
				"contextInfo": {
					"caSetId": "1000"
				}
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetNotFound))
			},
		},
		"missing id - validation error": {
			request: ListCASetActivationsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca set activations failed: struct validation: CASetID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			request: ListCASetActivationsRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199/activations",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal-server-error",
  "title": "Internal Server Error",
  "detail": "Error processing request",
  "instance": "TestInstances",
  "status": 500
}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ListCASetActivations(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
