package papi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiPatchPropertyHostnameActivation(t *testing.T) {
	tests := map[string]struct {
		params         PatchPropertyHostnameBucketRequest
		expected       *PatchPropertyHostnameBucketResponse
		responseStatus int
		responseBody   string
		requestBody    string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"200 OK - Add": {
			params: PatchPropertyHostnameBucketRequest{
				PropertyID: "property_id",
				Body: PatchPropertyHostnameBucketBody{
					Add: []PatchPropertyHostnameBucketAdd{
						{
							EdgeHostnameID:       "edge_hostname_id",
							CertProvisioningType: "DEFAULT",
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "cname.from",
						},
					},
					Remove:       []string{},
					Network:      "STAGING",
					NotifyEmails: []string{"noemail@akamai.com"},
					Note:         "note",
				},
			},
			expected: &PatchPropertyHostnameBucketResponse{
				ActivationLink: "/papi/v1/properties/property_id/hostname-activations/activation_id?groupId=group_id&contractId=contract_id",
				ActivationID:   "activation_id",
				Hostnames: []PatchHostnameItem{{
					CertProvisioningType: "DEFAULT",
					CnameFrom:            "cname.from",
					CnameTo:              "cname.to",
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "edge_hostname_id",
					CertStatus: CertStatusItem{
						ValidationCname: ValidationCname{
							Hostname: "validation_cname_hostname",
							Target:   "validation_cname_target",
						},
						Staging: []StatusItem{{
							Status: "PENDING",
						}},
						Production: []StatusItem{{
							Status: "PENDING",
						}},
					},
					Action: "ADD",
				}},
			},
			requestBody:    `{"add":[{"edgeHostnameId":"edge_hostname_id","certProvisioningType":"DEFAULT","cnameType":"EDGE_HOSTNAME","cnameFrom":"cname.from"}],"network":"STAGING","notifyEmails":["noemail@akamai.com"],"note":"note"}`,
			responseStatus: 201,
			responseBody: `
			{
				"activationLink": "/papi/v1/properties/property_id/hostname-activations/activation_id?groupId=group_id&contractId=contract_id",
				"activationId": "activation_id",
				"hostnames": [
					{
						"cnameType": "EDGE_HOSTNAME",
						"edgeHostnameId": "edge_hostname_id",
						"cnameFrom": "cname.from",
						"cnameTo": "cname.to",
						"certProvisioningType": "DEFAULT",
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
								"hostname": "validation_cname_hostname",
								"target": "validation_cname_target"
							}
						},
						"action": "ADD"
					}
				]
			}`,
			expectedPath: "/papi/v1/properties/property_id/hostnames",
			withError:    nil,
		},
		"200 OK - Add - optional fields": {
			params: PatchPropertyHostnameBucketRequest{
				PropertyID: "property_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
				Body: PatchPropertyHostnameBucketBody{
					Add: []PatchPropertyHostnameBucketAdd{
						{
							EdgeHostnameID:       "edge_hostname_id",
							CertProvisioningType: "DEFAULT",
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "cname.from",
						},
					},
					Remove:       []string{},
					Network:      "STAGING",
					NotifyEmails: []string{"noemail@akamai.com"},
					Note:         "note",
				},
			},
			expected: &PatchPropertyHostnameBucketResponse{
				ActivationLink: "/papi/v1/properties/property_id/hostname-activations/activation_id?groupId=group_id&contractId=contract_id",
				ActivationID:   "activation_id",
				Hostnames: []PatchHostnameItem{{
					CertProvisioningType: "DEFAULT",
					CnameFrom:            "cname.from",
					CnameTo:              "cname.to",
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "edge_hostname_id",
					CertStatus: CertStatusItem{
						ValidationCname: ValidationCname{
							Hostname: "validation_cname_hostname",
							Target:   "validation_cname_target",
						},
						Staging: []StatusItem{{
							Status: "PENDING",
						}},
						Production: []StatusItem{{
							Status: "PENDING",
						}},
					},
					Action: "ADD",
				}},
			},
			requestBody:    `{"add":[{"edgeHostnameId":"edge_hostname_id","certProvisioningType":"DEFAULT","cnameType":"EDGE_HOSTNAME","cnameFrom":"cname.from"}],"network":"STAGING","notifyEmails":["noemail@akamai.com"],"note":"note"}`,
			responseStatus: 201,
			responseBody: `
			{
				"activationLink": "/papi/v1/properties/property_id/hostname-activations/activation_id?groupId=group_id&contractId=contract_id",
				"activationId": "activation_id",
				"hostnames": [
					{
						"cnameType": "EDGE_HOSTNAME",
						"edgeHostnameId": "edge_hostname_id",
						"cnameFrom": "cname.from",
						"cnameTo": "cname.to",
						"certProvisioningType": "DEFAULT",
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
								"hostname": "validation_cname_hostname",
								"target": "validation_cname_target"
							}
						},
						"action": "ADD"
					}
				]
			}`,
			expectedPath: "/papi/v1/properties/property_id/hostnames?contractId=contract_id&groupId=group_id",
			withError:    nil,
		},
		"200 OK - Remove": {
			params: PatchPropertyHostnameBucketRequest{
				PropertyID: "property_id",
				Body: PatchPropertyHostnameBucketBody{
					Remove:  []string{"www.example.com"},
					Network: "STAGING",
				},
			},
			expected: &PatchPropertyHostnameBucketResponse{
				ActivationLink: "/papi/v1/properties/property_id/hostname-activations/activation_id?groupId=group_id&contractId=contract_id",
				ActivationID:   "activation_id",
				Hostnames: []PatchHostnameItem{{
					CertProvisioningType: "DEFAULT",
					CnameFrom:            "cname.from",
					CnameTo:              "cname.to",
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "edge_hostname_id",
					CertStatus: CertStatusItem{
						ValidationCname: ValidationCname{
							Hostname: "validation_cname_hostname",
							Target:   "validation_cname_target",
						},
						Staging: []StatusItem{{
							Status: "PENDING",
						}},
						Production: []StatusItem{{
							Status: "PENDING",
						}},
					},
					Action: "REMOVE",
				}},
			},
			requestBody:    `{"remove":["www.example.com"],"network":"STAGING"}`,
			responseStatus: 201,
			responseBody: `
			{
				"activationLink": "/papi/v1/properties/property_id/hostname-activations/activation_id?groupId=group_id&contractId=contract_id",
				"activationId": "activation_id",
				"hostnames": [
					{
						"cnameType": "EDGE_HOSTNAME",
						"edgeHostnameId": "edge_hostname_id",
						"cnameFrom": "cname.from",
						"cnameTo": "cname.to",
						"certProvisioningType": "DEFAULT",
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
								"hostname": "validation_cname_hostname",
								"target": "validation_cname_target"
							}
						},
						"action": "REMOVE"
					}
				]
			}`,
			expectedPath: "/papi/v1/properties/property_id/hostnames",
			withError:    nil,
		},
		"500 internal server error": {
			params: PatchPropertyHostnameBucketRequest{
				PropertyID: "property_id",
				ContractID: "contract_id",
				GroupID:    "group_id",
				Body: PatchPropertyHostnameBucketBody{
					Remove:  []string{"www.example.com"},
					Network: "STAGING",
				},
			},
			requestBody:    `{"remove":["www.example.com"],"network":"STAGING"}`,
			responseStatus: 500,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error removing property hostname",
				"status": 500
			}`,
			expectedPath: "/papi/v1/properties/property_id/hostnames?contractId=contract_id&groupId=group_id",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error removing property hostname",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error": {
			params: PatchPropertyHostnameBucketRequest{
				Body: PatchPropertyHostnameBucketBody{
					Add:     []PatchPropertyHostnameBucketAdd{{}, {}},
					Network: "",
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patching property hostname bucket: struct validation: Body: {\n\tAdd[0]: {\n\t\tCertProvisioningType: cannot be blank\n\t\tCnameFrom: cannot be blank\n\t\tCnameType: cannot be blank\n\t\tEdgeHostnameID: cannot be blank\n\t}\n\tAdd[1]: {\n\t\tCertProvisioningType: cannot be blank\n\t\tCnameFrom: cannot be blank\n\t\tCnameType: cannot be blank\n\t\tEdgeHostnameID: cannot be blank\n\t}\n\tNetwork: cannot be blank\n}\nPropertyID: cannot be blank",
					err.Error())
			},
		},
		"validation error - empty body": {
			params: PatchPropertyHostnameBucketRequest{
				PropertyID: "property_id",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "patching property hostname bucket: struct validation: Body: {\n\t: at least one hostname is required in add or remove list\n\tNetwork: cannot be blank\n}",
					err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPatch, r.Method)
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, string(requestBody), test.requestBody)
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.PatchPropertyHostnameBucket(context.Background(), test.params)

			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)

		})
	}
}
