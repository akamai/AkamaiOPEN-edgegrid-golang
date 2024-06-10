package v3

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPolicyActivations(t *testing.T) {
	t.Parallel()

	var policyID int64 = 1234

	tests := map[string]struct {
		params           ListPolicyActivationsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PolicyActivations
		withError        bool
		assertError      func(*testing.T, error)
	}{
		"200 OK": {
			params: ListPolicyActivationsRequest{
				PolicyID: policyID,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/cloudlets/v3/policies/1234/activations",
			expectedResponse: &PolicyActivations{
				PolicyActivations: []PolicyActivation{
					{
						CreatedBy:   "testUser",
						CreatedDate: test.NewTimeFromString(t, "2023-10-25T10:33:47.982Z"),
						FinishDate:  nil,
						ID:          234,
						Links: []Link{
							{
								Href: "/cloudlets/v3/policies/1234/activations/234",
								Rel:  "self",
							},
						},
						Network:              StagingNetwork,
						Operation:            OperationDeactivation,
						PolicyID:             policyID,
						PolicyVersion:        1,
						PolicyVersionDeleted: false,
						Status:               ActivationStatusInProgress,
					},
					{
						CreatedBy:   "testUser",
						CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
						FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
						ID:          123,
						Links: []Link{
							{
								Href: "/cloudlets/v3/policies/1234/activations/123",
								Rel:  "self",
							},
						},
						Network:              StagingNetwork,
						Operation:            OperationActivation,
						PolicyID:             policyID,
						PolicyVersion:        1,
						PolicyVersionDeleted: false,
						Status:               ActivationStatusSuccess,
					},
				},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/1234/activations?page=0&size=1000",
						Rel:  "self",
					},
				},
				Page: Page{
					Number:        0,
					Size:          1000,
					TotalElements: 2,
					TotalPages:    1,
				},
			},
			responseBody: `
{
	"content": [
		{
			"createdBy": "testUser",
			"createdDate": "2023-10-25T10:33:47.982Z",
			"finishDate": null,
			"id": 234,
			"links": [
				{
					"href": "/cloudlets/v3/policies/1234/activations/234",
					"rel": "self"
				}
			],
			"network": "STAGING",
			"operation": "DEACTIVATION",
			"policyId": 1234,
			"policyVersion": 1,
			"policyVersionDeleted": false,
			"status": "IN_PROGRESS"
		},
		{
			"createdBy": "testUser",
			"createdDate": "2023-10-23T11:21:19.896Z",
			"finishDate": "2023-10-23T11:22:57.589Z",
			"id": 123,
			"links": [
				{
					"href": "/cloudlets/v3/policies/1234/activations/123",
					"rel": "self"
				}
			],
			"network": "STAGING",
			"operation": "ACTIVATION",
			"policyId": 1234,
			"policyVersion": 1,
			"policyVersionDeleted": false,
			"status": "SUCCESS"
		}
	],
	"links": [
		{
			"href": "/cloudlets/v3/policies/1234/activations?page=0&size=1000",
			"rel": "self"
		}
	],
	"page": {
		"number": 0,
		"size": 1000,
		"totalElements": 2,
		"totalPages": 1
	}
}`,
		},
		"200 OK with query": {
			params: ListPolicyActivationsRequest{
				PolicyID: policyID,
				Page:     1,
				Size:     10,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/cloudlets/v3/policies/1234/activations?page=1&size=10",
			expectedResponse: &PolicyActivations{
				PolicyActivations: []PolicyActivation{},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/1234/activations?page=0&size=10",
						Rel:  "first",
					},
					{
						Href: "/cloudlets/v3/policies/1234/activations?page=0&size=10",
						Rel:  "prev",
					},
					{
						Href: "/cloudlets/v3/policies/1234/activations?page=1&size=10",
						Rel:  "self",
					},
					{
						Href: "/cloudlets/v3/policies/1234/activations?page=0&size=10",
						Rel:  "last",
					},
				},
				Page: Page{
					Number:        1,
					Size:          10,
					TotalElements: 2,
					TotalPages:    1,
				},
			},
			responseBody: `
{
	"content": [],
	"links": [
		{
			"href": "/cloudlets/v3/policies/1234/activations?page=0&size=10",
			"rel": "first"
		},
		{
			"href": "/cloudlets/v3/policies/1234/activations?page=0&size=10",
			"rel": "prev"
		},
		{
			"href": "/cloudlets/v3/policies/1234/activations?page=1&size=10",
			"rel": "self"
		},
		{
			"href": "/cloudlets/v3/policies/1234/activations?page=0&size=10",
			"rel": "last"
		}
	],
	"page": {
		"number": 1,
		"size": 10,
		"totalElements": 2,
		"totalPages": 1
	}
}`,
		},
		"500 Internal Server Error": {
			params: ListPolicyActivationsRequest{
				PolicyID: policyID,
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/cloudlets/v3/policies/1234/activations",
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"status": 500,
	"requestId": "1",
	"requestTime": "12:00",
	"clientIp": "1.1.1.1",
	"serverIp": "2.2.2.2",
	"method": "GET"
}`,
			withError: true,
			assertError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "internal_error",
					Title:       "Internal Server Error",
					Status:      http.StatusInternalServerError,
					RequestID:   "1",
					RequestTime: "12:00",
					ClientIP:    "1.1.1.1",
					ServerIP:    "2.2.2.2",
					Method:      "GET",
				}
				assert.ErrorIs(t, err, want)
			},
		},
		"request validation failed": {
			params: ListPolicyActivationsRequest{
				Page: -1,
				Size: 5,
			},
			withError: true,
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "Page: must be no less than 0")
				assert.ErrorContains(t, err, "PolicyID: cannot be blank")
				assert.ErrorContains(t, err, "Size: must be no less than 10")
			},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListPolicyActivations(context.Background(), tc.params)

			if tc.withError {
				tc.assertError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestActivatePolicy(t *testing.T) {
	t.Parallel()

	var policyID int64 = 1234

	tests := map[string]struct {
		params           ActivatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedReqBody  string
		expectedResponse *PolicyActivation
		withError        bool
		assertError      func(*testing.T, error)
	}{
		"202 Accepted": {
			params: ActivatePolicyRequest{
				PolicyID:      policyID,
				Network:       StagingNetwork,
				PolicyVersion: 1,
			},
			responseStatus: http.StatusAccepted,
			expectedPath:   "/cloudlets/v3/policies/1234/activations",
			expectedResponse: &PolicyActivation{
				CreatedBy:   "testUser",
				CreatedDate: test.NewTimeFromString(t, "2023-10-25T10:33:47.982Z"),
				FinishDate:  nil,
				ID:          123,
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/1234/activations/123",
						Rel:  "self",
					},
				},
				Network:              StagingNetwork,
				Operation:            OperationActivation,
				PolicyID:             policyID,
				PolicyVersion:        1,
				PolicyVersionDeleted: false,
				Status:               ActivationStatusInProgress,
			},
			expectedReqBody: `
{
	"network": "STAGING",
	"policyVersion": 1,
	"operation": "ACTIVATION"
}`,
			responseBody: `
{
	"createdBy": "testUser",
	"createdDate": "2023-10-25T10:33:47.982Z",
	"finishDate": null,
	"id": 123,
	"links": [
		{
			"href": "/cloudlets/v3/policies/1234/activations/123",
			"rel": "self"
		}
	],
	"network": "STAGING",
	"operation": "ACTIVATION",
	"policyId": 1234,
	"policyVersion": 1,
	"policyVersionDeleted": false,
	"status": "IN_PROGRESS"
}`,
		},
		"500 Internal Server Error": {
			params: ActivatePolicyRequest{
				PolicyID:      policyID,
				Network:       StagingNetwork,
				PolicyVersion: 1,
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/cloudlets/v3/policies/1234/activations",
			expectedReqBody: `
{
	"network": "STAGING",
	"policyVersion": 1,
	"operation": "ACTIVATION"
}`,
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"status": 500,
	"requestId": "1",
	"requestTime": "12:00",
	"clientIp": "1.1.1.1",
	"serverIp": "2.2.2.2",
	"method": "GET"
}`,
			withError: true,
			assertError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "internal_error",
					Title:       "Internal Server Error",
					Status:      http.StatusInternalServerError,
					RequestID:   "1",
					RequestTime: "12:00",
					ClientIP:    "1.1.1.1",
					ServerIP:    "2.2.2.2",
					Method:      "GET",
				}
				assert.ErrorIs(t, err, want)
			},
		},
		"request validation failed": {
			params:    ActivatePolicyRequest{},
			withError: true,
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "PolicyID: cannot be blank")
				assert.ErrorContains(t, err, "PolicyVersion: cannot be blank")
				assert.ErrorContains(t, err, "Network: cannot be blank")
			},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
				if tc.expectedReqBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, tc.expectedReqBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ActivatePolicy(context.Background(), tc.params)

			if tc.withError {
				tc.assertError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestDeactivatePolicy(t *testing.T) {
	t.Parallel()

	var policyID int64 = 1234

	tests := map[string]struct {
		params           DeactivatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedReqBody  string
		expectedResponse *PolicyActivation
		withError        bool
		assertError      func(*testing.T, error)
	}{
		"202 Accepted": {
			params: DeactivatePolicyRequest{
				PolicyID:      policyID,
				Network:       StagingNetwork,
				PolicyVersion: 1,
			},
			responseStatus: http.StatusAccepted,
			expectedPath:   "/cloudlets/v3/policies/1234/activations",
			expectedResponse: &PolicyActivation{
				CreatedBy:   "testUser",
				CreatedDate: test.NewTimeFromString(t, "2023-10-25T10:33:47.982Z"),
				FinishDate:  nil,
				ID:          123,
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/1234/activations/123",
						Rel:  "self",
					},
				},
				Network:              StagingNetwork,
				Operation:            OperationDeactivation,
				PolicyID:             policyID,
				PolicyVersion:        1,
				PolicyVersionDeleted: false,
				Status:               ActivationStatusInProgress,
			},
			expectedReqBody: `
{
	"network": "STAGING",
	"policyVersion": 1,
	"operation": "DEACTIVATION"
}`,
			responseBody: `
{
	"createdBy": "testUser",
	"createdDate": "2023-10-25T10:33:47.982Z",
	"finishDate": null,
	"id": 123,
	"links": [
		{
			"href": "/cloudlets/v3/policies/1234/activations/123",
			"rel": "self"
		}
	],
	"network": "STAGING",
	"operation": "DEACTIVATION",
	"policyId": 1234,
	"policyVersion": 1,
	"policyVersionDeleted": false,
	"status": "IN_PROGRESS"
}`,
		},
		"500 Internal Server Error": {
			params: DeactivatePolicyRequest{
				PolicyID:      policyID,
				Network:       ProductionNetwork,
				PolicyVersion: 1,
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/cloudlets/v3/policies/1234/activations",
			expectedReqBody: `
{
	"network": "PRODUCTION",
	"policyVersion": 1,
	"operation": "DEACTIVATION"
}`,
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"status": 500,
	"requestId": "1",
	"requestTime": "12:00",
	"clientIp": "1.1.1.1",
	"serverIp": "2.2.2.2",
	"method": "GET"
}`,
			withError: true,
			assertError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "internal_error",
					Title:       "Internal Server Error",
					Status:      http.StatusInternalServerError,
					RequestID:   "1",
					RequestTime: "12:00",
					ClientIP:    "1.1.1.1",
					ServerIP:    "2.2.2.2",
					Method:      "GET",
				}
				assert.ErrorIs(t, err, want)
			},
		},
		"request validation failed": {
			params: DeactivatePolicyRequest{
				Network: "OTHER",
			},
			withError: true,
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "PolicyID: cannot be blank")
				assert.ErrorContains(t, err, "PolicyVersion: cannot be blank")
				assert.ErrorContains(t, err, "Network: value 'OTHER' is invalid. Must be one of: 'STAGING' or 'PRODUCTION'")
			},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
				if tc.expectedReqBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, tc.expectedReqBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeactivatePolicy(context.Background(), tc.params)

			if tc.withError {
				tc.assertError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetPolicyActivation(t *testing.T) {
	t.Parallel()

	var policyID int64 = 1234
	var activationID int64 = 123

	tests := map[string]struct {
		params           GetPolicyActivationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PolicyActivation
		withError        bool
		assertError      func(*testing.T, error)
	}{
		"200 OK": {
			params: GetPolicyActivationRequest{
				PolicyID:     policyID,
				ActivationID: activationID,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/cloudlets/v3/policies/1234/activations/123",
			expectedResponse: &PolicyActivation{
				CreatedBy:   "testUser",
				CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
				ID:          activationID,
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/1234/activations/123",
						Rel:  "self",
					},
				},
				Network:              StagingNetwork,
				Operation:            OperationActivation,
				PolicyID:             policyID,
				PolicyVersion:        1,
				PolicyVersionDeleted: false,
				Status:               ActivationStatusSuccess,
			},
			responseBody: `
{
	"createdBy": "testUser",
	"createdDate": "2023-10-23T11:21:19.896Z",
	"finishDate": "2023-10-23T11:22:57.589Z",
	"id": 123,
	"links": [
		{
			"href": "/cloudlets/v3/policies/1234/activations/123",
			"rel": "self"
		}
	],
	"network": "STAGING",
	"operation": "ACTIVATION",
	"policyId": 1234,
	"policyVersion": 1,
	"policyVersionDeleted": false,
	"status": "SUCCESS"
}`,
		},
		"404 Not Found": {
			params: GetPolicyActivationRequest{
				PolicyID:     policyID,
				ActivationID: 1,
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/cloudlets/v3/policies/1234/activations/1",
			responseBody: `
{
	"instance": "testInstance",
	"status": 404,
	"title": "Not found",
	"type": "/cloudlets/v3/error-types/not-found",
	"errors": [
		{
			"detail": "Activation with id '1' not found.",
			"title": "Not found"
		}
	]
}`,
			withError: true,
			assertError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/cloudlets/v3/error-types/not-found",
					Title:    "Not found",
					Status:   http.StatusNotFound,
					Instance: "testInstance",
					Errors:   json.RawMessage(`[{"detail": "Activation with id '1' not found.", "title": "Not found"}]`),
				}
				assert.ErrorIs(t, err, want)
			},
		},
		"500 Internal Server Error": {
			params: GetPolicyActivationRequest{
				PolicyID:     policyID,
				ActivationID: activationID,
			},
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/cloudlets/v3/policies/1234/activations/123",
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"status": 500,
	"requestId": "1",
	"requestTime": "12:00",
	"clientIp": "1.1.1.1",
	"serverIp": "2.2.2.2",
	"method": "GET"
}`,
			withError: true,
			assertError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "internal_error",
					Title:       "Internal Server Error",
					Status:      http.StatusInternalServerError,
					RequestID:   "1",
					RequestTime: "12:00",
					ClientIP:    "1.1.1.1",
					ServerIP:    "2.2.2.2",
					Method:      "GET",
				}
				assert.ErrorIs(t, err, want)
			},
		},
		"request validation failed": {
			params:    GetPolicyActivationRequest{},
			withError: true,
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "PolicyID: cannot be blank")
				assert.ErrorContains(t, err, "ActivationID: cannot be blank")
			},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetPolicyActivation(context.Background(), tc.params)

			if tc.withError {
				tc.assertError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
