package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiGetPropertyHostnameActivation(t *testing.T) {
	tests := map[string]struct {
		params         GetPropertyHostnameActivationRequest
		expected       *GetPropertyHostnameActivationResponse
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			params: GetPropertyHostnameActivationRequest{
				PropertyID:           "property_id",
				HostnameActivationID: "hostname_activation_id",
				ContractID:           "contract_id",
				GroupID:              "group_id",
			},
			expected: &GetPropertyHostnameActivationResponse{
				AccountID:  "account_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
				HostnameActivation: HostnameActivationGetItem{
					ActivationType:       "ACTIVATE",
					HostnameActivationID: "hostname_activation_id",
					PropertyName:         "property_name",
					PropertyID:           "property_id",
					Network:              "STAGING",
					Status:               "ACTIVE",
					SubmitDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
					UpdateDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
					Note:                 "",
					NotifyEmails:         []string{"nomail@akamai.com"},
				}},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId": "account_id",
				"contractId": "contract_id",
				"groupId": "group_id",
				"hostnameActivations": {
					"items": [
						{
							"propertyName": "property_name",
							"propertyId": "property_id",
							"network": "STAGING",
							"activationType": "ACTIVATE",
							"status": "ACTIVE",
							"submitDate": "2025-01-13T11:22:33Z",
							"updateDate": "2025-01-13T11:22:33Z",
							"note": "",
							"notifyEmails": [
								"nomail@akamai.com"
							],
							
							"hostnameActivationId": "hostname_activation_id"
							}
						]
					}
				}`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations/hostname_activation_id?contractId=contract_id&groupId=group_id",
			withError:    nil,
		},
		"200 OK - include hostnames": {
			params: GetPropertyHostnameActivationRequest{
				PropertyID:           "property_id",
				HostnameActivationID: "hostname_activation_id",
				ContractID:           "contract_id",
				GroupID:              "group_id",
				IncludeHostnames:     true,
			},
			expected: &GetPropertyHostnameActivationResponse{
				AccountID:  "account_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
				HostnameActivation: HostnameActivationGetItem{
					ActivationType:       "ACTIVATE",
					HostnameActivationID: "hostname_activation_id",
					PropertyName:         "property_name",
					PropertyID:           "property_id",
					Network:              "STAGING",
					Status:               "ACTIVE",
					SubmitDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
					UpdateDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
					Note:                 "",
					NotifyEmails:         []string{"nomail@akamai.com"},
					Hostnames: []PropertyHostnameItem{{
						Action:               "ADD",
						EdgeHostnameID:       "edge_hostname_id",
						CnameFrom:            "hostname.com",
						CnameTo:              "hostname.com.net",
						CertProvisioningType: "CPS_MANAGED",
						CertStatus: CertStatusItem{
							ValidationCname: ValidationCname{
								Hostname: "hostname.com.net",
								Target:   "hostname.com.net",
							},
							Staging: []StatusItem{{
								Status: "PENDING",
							}},
							Production: []StatusItem{{
								Status: "PENDING",
							}},
						},
					}},
				}},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId": "account_id",
				"contractId": "contract_id",
				"groupId": "group_id",
				"hostnameActivations": {
					"items": [
						{
							"propertyName": "property_name",
							"propertyId": "property_id",
							"network": "STAGING",
							"activationType": "ACTIVATE",
							"status": "ACTIVE",
							"submitDate": "2025-01-13T11:22:33Z",
							"updateDate": "2025-01-13T11:22:33Z",
							"note": "",
							"notifyEmails": [
								"nomail@akamai.com"
							],
							
							"hostnameActivationId": "hostname_activation_id",
							"hostnames": {
								"items": [
									{
										"edgeHostnameId": "edge_hostname_id",
										"cnameFrom": "hostname.com",
										"cnameTo": "hostname.com.net",
										"certProvisioningType": "CPS_MANAGED",
										"certStatus": {
											"production": [
												{
													"status": "PENDING"
												}
											],
											"staging": [
												{
													"status": "PENDING"
												}
											],
											"validationCname": {
												"hostname": "hostname.com.net",
												"target": "hostname.com.net"
											}
										},
										"action": "ADD"
									}
									]
								}
							}
						]
					}
				}`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations/hostname_activation_id?contractId=contract_id&groupId=group_id&includeHostnames=true",
			withError:    nil,
		},
		"validation error": {
			params: GetPropertyHostnameActivationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "GroupID: cannot be blank")
				assert.Contains(t, err.Error(), "ContractID: cannot be blank")
				assert.Contains(t, err.Error(), "PropertyID: cannot be blank")
				assert.Contains(t, err.Error(), "HostnameActivationID: cannot be blank")
			},
		},
		"500 Internal error": {
			params: GetPropertyHostnameActivationRequest{
				PropertyID:           "property_id",
				HostnameActivationID: "hostname_activation_id",
				ContractID:           "contract_id",
				GroupID:              "group_id",
				IncludeHostnames:     true,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching property hostname activation",
				"status": 500
			}`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations/hostname_activation_id?contractId=contract_id&groupId=group_id&includeHostnames=true",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching property hostname activation",
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
			result, err := client.GetPropertyHostnameActivation(context.Background(), test.params)

			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)

		})
	}
}

func TestPapiListPropertyHostnameActivations(t *testing.T) {
	tests := map[string]struct {
		params         ListPropertyHostnameActivationsRequest
		expected       *ListPropertyHostnameActivationsResponse
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			params: ListPropertyHostnameActivationsRequest{
				PropertyID: "property_id",
			},
			expected: &ListPropertyHostnameActivationsResponse{
				AccountID:  "account_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
				HostnameActivations: HostnameActivationsList{
					Items: []HostnameActivationListItem{
						{
							ActivationType:       "ACTIVATE",
							HostnameActivationID: "hostname_activation_id",
							PropertyName:         "property_name",
							PropertyID:           "property_id",
							Network:              "PRODUCTION",
							Status:               "ACTIVE",
							SubmitDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							UpdateDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							Note:                 "note",
							NotifyEmails:         []string{"nomail@akamai.com"},
						},
						{
							ActivationType:       "ACTIVATE",
							HostnameActivationID: "hostname_activation_id",
							PropertyName:         "property_name",
							PropertyID:           "property_id",
							Network:              "STAGING",
							Status:               "ACTIVE",
							SubmitDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							UpdateDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							Note:                 "note",
							NotifyEmails:         []string{"nomail@akamai.com"},
						},
					},
					TotalItems:       ptr.To(2),
					CurrentItemCount: ptr.To(2),
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
		{
						"accountId": "account_id",
						"contractId": "contract_id",
						"groupId": "group_id",
						"hostnameActivations": {
							"items": [
								{
									"propertyName": "property_name",
									"propertyId": "property_id",
									"network": "PRODUCTION",
									"activationType": "ACTIVATE",
									"status": "ACTIVE",
									"submitDate": "2025-01-13T11:22:33Z",
									"updateDate": "2025-01-13T11:22:33Z",
									"note": "note",
									"notifyEmails": [
										"nomail@akamai.com"
									],
									"hostnameActivationId": "hostname_activation_id"
								},
								{
									"propertyName": "property_name",
									"propertyId": "property_id",
									"network": "STAGING",
									"activationType": "ACTIVATE",
									"status": "ACTIVE",
									"submitDate": "2025-01-13T11:22:33Z",
									"updateDate": "2025-01-13T11:22:33Z",
									"note": "note",
									"notifyEmails": [
										"nomail@akamai.com"
									],
									"hostnameActivationId": "hostname_activation_id"
								}
							],
							"totalItems": 2,
							"currentItemCount": 2
						}
					}`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations",
			withError:    nil,
		},
		"200 OK - optional fields ": {
			params: ListPropertyHostnameActivationsRequest{
				PropertyID: "property_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
			},
			expected: &ListPropertyHostnameActivationsResponse{
				AccountID:  "account_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
				HostnameActivations: HostnameActivationsList{
					Items: []HostnameActivationListItem{
						{
							ActivationType:       "ACTIVATE",
							HostnameActivationID: "hostname_activation_id",
							PropertyName:         "property_name",
							PropertyID:           "property_id",
							Network:              "PRODUCTION",
							Status:               "ACTIVE",
							SubmitDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							UpdateDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							Note:                 "note",
							NotifyEmails:         []string{"nomail@akamai.com"},
						},
						{
							ActivationType:       "ACTIVATE",
							HostnameActivationID: "hostname_activation_id",
							PropertyName:         "property_name",
							PropertyID:           "property_id",
							Network:              "STAGING",
							Status:               "ACTIVE",
							SubmitDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							UpdateDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							Note:                 "note",
							NotifyEmails:         []string{"nomail@akamai.com"},
						},
					},
					TotalItems:       ptr.To(2),
					CurrentItemCount: ptr.To(2),
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
		{
						"accountId": "account_id",
						"contractId": "contract_id",
						"groupId": "group_id",
						"hostnameActivations": {
							"items": [
								{
									"propertyName": "property_name",
									"propertyId": "property_id",
									"network": "PRODUCTION",
									"activationType": "ACTIVATE",
									"status": "ACTIVE",
									"submitDate": "2025-01-13T11:22:33Z",
									"updateDate": "2025-01-13T11:22:33Z",
									"note": "note",
									"notifyEmails": [
										"nomail@akamai.com"
									],
									"hostnameActivationId": "hostname_activation_id"
								},
								{
									"propertyName": "property_name",
									"propertyId": "property_id",
									"network": "STAGING",
									"activationType": "ACTIVATE",
									"status": "ACTIVE",
									"submitDate": "2025-01-13T11:22:33Z",
									"updateDate": "2025-01-13T11:22:33Z",
									"note": "note",
									"notifyEmails": [
										"nomail@akamai.com"
									],
									"hostnameActivationId": "hostname_activation_id"
								}
							],
							"totalItems": 2,
							"currentItemCount": 2
						}
					}`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations?contractId=contract_id&groupId=group_id",
			withError:    nil,
		},
		"200 OK - pagination": {
			params: ListPropertyHostnameActivationsRequest{
				PropertyID: "property_id",
				Limit:      1,
				Offset:     1,
				ContractID: "contract_id",
				GroupID:    "group_id",
			},
			expected: &ListPropertyHostnameActivationsResponse{
				AccountID:  "account_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
				HostnameActivations: HostnameActivationsList{
					Items: []HostnameActivationListItem{
						{
							ActivationType:       "ACTIVATE",
							HostnameActivationID: "hostname_activation_id",
							PropertyName:         "property_name",
							PropertyID:           "property_id",
							Network:              "PRODUCTION",
							Status:               "ACTIVE",
							SubmitDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							UpdateDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
							Note:                 "note",
							NotifyEmails:         []string{"nomail@akamai.com"},
						},
					},
					TotalItems:       ptr.To(3),
					CurrentItemCount: ptr.To(1),
					NextLink:         ptr.To("/papi/v1/properties/property_id/hostname-activations?contractId=contract_id&groupId=group_id&limit=1&offset=1"),
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
		{
						"accountId": "account_id",
						"contractId": "contract_id",
						"groupId": "group_id",
						"hostnameActivations": {
							"items": [
								{
									"propertyName": "property_name",
									"propertyId": "property_id",
									"network": "PRODUCTION",
									"activationType": "ACTIVATE",
									"status": "ACTIVE",
									"submitDate": "2025-01-13T11:22:33Z",
									"updateDate": "2025-01-13T11:22:33Z",
									"note": "note",
									"notifyEmails": [
										"nomail@akamai.com"
									],
									"hostnameActivationId": "hostname_activation_id"
								}
							],
							"totalItems": 3,
							"currentItemCount": 1,
							"prevLink": "",
							"nextLink": "/papi/v1/properties/property_id/hostname-activations?contractId=contract_id&groupId=group_id&limit=1&offset=1"
						}
					}`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations?contractId=contract_id&groupId=group_id&limit=1&offset=1",
			withError:    nil,
		},
		"500 internal server error": {
			params: ListPropertyHostnameActivationsRequest{
				PropertyID: "property_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching property hostname activations list",
				"status": 500
			}`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations?contractId=contract_id&groupId=group_id",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching property hostname activations list",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error": {
			params: ListPropertyHostnameActivationsRequest{
				Offset: -1,
				Limit:  -1,
			},
			expected: &ListPropertyHostnameActivationsResponse{},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "PropertyID: cannot be blank")
				assert.Contains(t, err.Error(), "Limit: must be no less than 1")
				assert.Contains(t, err.Error(), "Offset: must be no less than 0")
			},
		},
		"validation error - GroupID provided, ContractID missing": {
			params: ListPropertyHostnameActivationsRequest{
				PropertyID: "property_id",
				GroupID:    "group_id",
			},
			expected: &ListPropertyHostnameActivationsResponse{},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "ContractID: cannot be blank when GroupID is provided")
			},
		},
		"validation error - ContractID provided, GroupID missing": {
			params: ListPropertyHostnameActivationsRequest{
				PropertyID: "property_id",
				ContractID: "contract_id",
			},
			expected: &ListPropertyHostnameActivationsResponse{},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "GroupID: cannot be blank when ContractID is provided")
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
			result, err := client.ListPropertyHostnameActivations(context.Background(), test.params)

			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)

		})
	}

}

func TestPapiCancelPropertyHostnameActivation(t *testing.T) {
	tests := map[string]struct {
		params         CancelPropertyHostnameActivationRequest
		expected       *CancelPropertyHostnameActivationResponse
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK": {
			params: CancelPropertyHostnameActivationRequest{
				PropertyID:           "property_id",
				HostnameActivationID: "hostname_activation_id",
				ContractID:           "contract_id",
				GroupID:              "group_id",
			},
			expected: &CancelPropertyHostnameActivationResponse{
				AccountID:  "account_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
				HostnameActivation: HostnameActivationCancelItem{
					ActivationType:       "ACTIVATE",
					HostnameActivationID: "hostname_activation_id",
					PropertyName:         "property_name",
					PropertyID:           "property_id",
					Network:              "STAGING",
					Status:               "PENDING_CANCELLATION",
					PropertyVersion:      1,
					SubmitDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
					UpdateDate:           test.NewTimeFromString(t, "2025-01-13T11:22:33Z"),
					Note:                 "",
					NotifyEmails:         []string{"nomail@akamai.com"},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"accountId":"account_id",
				"contractId":"contract_id",
				"groupId":"group_id",
				"hostnameActivations":{
				   "items":[
					  {
						 "propertyName":"property_name",
						 "propertyId":"property_id",
						 "network":"STAGING",
						 "activationType":"ACTIVATE",
						 "status":"PENDING_CANCELLATION",
						 "propertyVersion":1,
						 "submitDate":"2025-01-13T11:22:33Z",
						 "updateDate":"2025-01-13T11:22:33Z",
						 "note":"",
						 "notifyEmails":[
							"nomail@akamai.com"
						 ],
						 "hostnameActivationId":"hostname_activation_id"
					  }
				   ]
				}
			 }`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations/hostname_activation_id?contractId=contract_id&groupId=group_id",
			withError:    nil,
		},
		"204 No Content - activation already aborted": {
			params: CancelPropertyHostnameActivationRequest{
				PropertyID:           "property_id",
				HostnameActivationID: "hostname_activation_id",
				ContractID:           "contract_id",
				GroupID:              "group_id",
			},
			expected:       &CancelPropertyHostnameActivationResponse{},
			responseStatus: 204,
			responseBody:   "",
			expectedPath:   "/papi/v1/properties/property_id/hostname-activations/hostname_activation_id?contractId=contract_id&groupId=group_id",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "canceling hostname activation: activation already aborted", err.Error())
			},
		},
		"validation error": {
			params: CancelPropertyHostnameActivationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "GroupID: cannot be blank")
				assert.Contains(t, err.Error(), "ContractID: cannot be blank")
				assert.Contains(t, err.Error(), "PropertyID: cannot be blank")
				assert.Contains(t, err.Error(), "HostnameActivationID: cannot be blank")
			},
		},
		"500 Internal error": {
			params: CancelPropertyHostnameActivationRequest{
				PropertyID:           "property_id",
				HostnameActivationID: "hostname_activation_id",
				ContractID:           "contract_id",
				GroupID:              "group_id",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching property hostname activation",
				"status": 500
			}`,
			expectedPath: "/papi/v1/properties/property_id/hostname-activations/hostname_activation_id?contractId=contract_id&groupId=group_id",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching property hostname activation",
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
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CancelPropertyHostnameActivation(context.Background(), test.params)

			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)

		})
	}
}
