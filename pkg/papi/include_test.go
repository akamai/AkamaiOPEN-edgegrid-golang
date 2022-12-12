package papi

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListIncludes(t *testing.T) {
	tests := map[string]struct {
		params           ListIncludesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListIncludesResponse
		withError        error
	}{
		"200 OK - list includes given contractId and groupId": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "includes": {
        "items": [
            {
                "accountId": "test_account",
                "contractId": "test_contract",
                "groupId": "test_group",
                "latestVersion": 1,
                "stagingVersion": null,
                "productionVersion": null,
                "assetId": "test_asset",
                "includeId": "inc_123456",
                "includeName": "test_include",
                "includeType": "MICROSERVICES"
            },
            {
                "accountId": "test_account_1",
                "contractId": "test_contract",
                "groupId": "test_group",
                "latestVersion": 1,
                "stagingVersion": 1,
                "productionVersion": null,
                "assetId": "test_asset_1",
                "includeId": "inc_456789",
                "includeName": "test_include_1",
                "includeType": "COMMON_SETTINGS"
            }
		]
	}
}`,
			params: ListIncludesRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			expectedPath: "/papi/v1/includes?contractId=test_contract&groupId=test_group",
			expectedResponse: &ListIncludesResponse{
				Includes: IncludeItems{
					Items: []Include{
						{
							AccountID:     "test_account",
							AssetID:       "test_asset",
							ContractID:    "test_contract",
							GroupID:       "test_group",
							IncludeID:     "inc_123456",
							IncludeName:   "test_include",
							IncludeType:   IncludeTypeMicroServices,
							LatestVersion: 1,
						},
						{
							AccountID:      "test_account_1",
							AssetID:        "test_asset_1",
							ContractID:     "test_contract",
							GroupID:        "test_group",
							IncludeID:      "inc_456789",
							IncludeName:    "test_include_1",
							IncludeType:    IncludeTypeCommonSettings,
							LatestVersion:  1,
							StagingVersion: tools.IntPtr(1),
						},
					},
				},
			},
		},
		"200 OK - list includes given only contractId": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "includes": {
        "items": [
            {
                "accountId": "test_account",
                "contractId": "test_contract",
                "groupId": "test_group",
                "latestVersion": 1,
                "stagingVersion": null,
                "productionVersion": null,
                "assetId": "test_asset",
                "includeId": "inc_123456",
                "includeName": "test_include",
                "includeType": "MICROSERVICES"
            },
            {
                "accountId": "test_account_1",
                "contractId": "test_contract",
                "groupId": "test_group_1",
                "latestVersion": 1,
                "stagingVersion": 1,
                "productionVersion": null,
                "assetId": "test_asset_1",
                "includeId": "inc_456789",
                "includeName": "test_include_1",
                "includeType": "COMMON_SETTINGS"
            }
		]
	}
}`,
			params: ListIncludesRequest{
				ContractID: "test_contract",
			},
			expectedPath: "/papi/v1/includes?contractId=test_contract",
			expectedResponse: &ListIncludesResponse{
				Includes: IncludeItems{
					Items: []Include{
						{
							AccountID:     "test_account",
							AssetID:       "test_asset",
							ContractID:    "test_contract",
							GroupID:       "test_group",
							IncludeID:     "inc_123456",
							IncludeName:   "test_include",
							IncludeType:   "MICROSERVICES",
							LatestVersion: 1,
						},
						{
							AccountID:      "test_account_1",
							AssetID:        "test_asset_1",
							ContractID:     "test_contract",
							GroupID:        "test_group_1",
							IncludeID:      "inc_456789",
							IncludeName:    "test_include_1",
							IncludeType:    IncludeTypeCommonSettings,
							LatestVersion:  1,
							StagingVersion: tools.IntPtr(1),
						},
					},
				},
			},
		},
		"200 OK - no includes under given contractId and groupId": {
			params: ListIncludesRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "includes": {
		"items": []
	}
}`,
			expectedPath:     "/papi/v1/includes?contractId=test_contract&groupId=test_group",
			expectedResponse: &ListIncludesResponse{Includes: IncludeItems{Items: []Include{}}},
		},
		"500 internal server error": {
			params: ListIncludesRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching includes",
    "status": 500
}`,
			expectedPath: "/papi/v1/includes?contractId=test_contract&groupId=test_group",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching includes",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing contractId": {
			params: ListIncludesRequest{
				GroupID: "test_group",
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
			result, err := client.ListIncludes(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListIncludeParents(t *testing.T) {
	tests := map[string]struct {
		params           ListIncludeParentsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListIncludeParentsResponse
		withError        error
	}{
		"200 OK - list include parents given includeId, contractId and groupId": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "properties": {
        "items": [
            {
                "accountId": "test_account",
                "contractId": "test_contract",
                "groupId": "test_group",
                "propertyId": "prp_123456",
                "propertyName": "test_property",
                "stagingVersion": 1,
                "productionVersion": null,
                "assetId": "test_asset"
            }
        ]
    }
}`,
			params: ListIncludeParentsRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_456789",
			},
			expectedPath: "/papi/v1/includes/inc_456789/parents?contractId=test_contract&groupId=test_group",
			expectedResponse: &ListIncludeParentsResponse{
				Properties: ParentPropertyItems{
					Items: []ParentProperty{
						{
							AccountID:      "test_account",
							AssetID:        "test_asset",
							ContractID:     "test_contract",
							GroupID:        "test_group",
							PropertyID:     "prp_123456",
							PropertyName:   "test_property",
							StagingVersion: tools.IntPtr(1),
						},
					},
				},
			},
		},
		"200 OK - list includes given only includeId": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "properties": {
        "items": [
            {
                "accountId": "test_account",
                "contractId": "test_contract",
                "groupId": "test_group",
                "propertyId": "prp_123456",
                "propertyName": "test_property",
                "stagingVersion": 1,
                "productionVersion": null,
                "assetId": "test_asset"
            }
        ]
    }
}`,
			params: ListIncludeParentsRequest{
				IncludeID: "inc_456789",
			},
			expectedPath: "/papi/v1/includes/inc_456789/parents",
			expectedResponse: &ListIncludeParentsResponse{
				Properties: ParentPropertyItems{
					Items: []ParentProperty{
						{
							AccountID:      "test_account",
							AssetID:        "test_asset",
							ContractID:     "test_contract",
							GroupID:        "test_group",
							PropertyID:     "prp_123456",
							PropertyName:   "test_property",
							StagingVersion: tools.IntPtr(1),
						},
					},
				},
			},
		},
		"200 OK - no parents for given include": {
			params: ListIncludeParentsRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_456789",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "properties": {
		"items": []
	}
}`,
			expectedPath:     "/papi/v1/includes/inc_456789/parents?contractId=test_contract&groupId=test_group",
			expectedResponse: &ListIncludeParentsResponse{Properties: ParentPropertyItems{Items: []ParentProperty{}}},
		},
		"500 internal server error": {
			params: ListIncludeParentsRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_456789",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching properties",
    "status": 500
}`,
			expectedPath: "/papi/v1/includes/inc_456789/parents?contractId=test_contract&groupId=test_group",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching properties",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: ListIncludeParentsRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
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
			result, err := client.ListIncludeParents(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetInclude(t *testing.T) {
	tests := map[string]struct {
		params           GetIncludeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetIncludeResponse
		withError        error
	}{
		"200 OK - get include given includeId, contractId and groupId": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "includes": {
        "items": [
            {
                "accountId": "test_account",
                "contractId": "test_contract",
                "groupId": "test_group",
                "latestVersion": 1,
                "stagingVersion": null,
                "productionVersion": null,
				"propertyType": "INCLUDE",
                "assetId": "test_asset",
                "includeId": "inc_123456",
                "includeName": "test_include",
                "includeType": "MICROSERVICES"
            }
		]
	}
}`,
			params: GetIncludeRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_123456",
			},
			expectedPath: "/papi/v1/includes/inc_123456?contractId=test_contract&groupId=test_group",
			expectedResponse: &GetIncludeResponse{
				Includes: IncludeItems{
					Items: []Include{
						{
							AccountID:     "test_account",
							AssetID:       "test_asset",
							ContractID:    "test_contract",
							GroupID:       "test_group",
							IncludeID:     "inc_123456",
							IncludeName:   "test_include",
							IncludeType:   "MICROSERVICES",
							LatestVersion: 1,
							PropertyType:  tools.StringPtr("INCLUDE"),
						},
					},
				},
				Include: Include{
					AccountID:     "test_account",
					AssetID:       "test_asset",
					ContractID:    "test_contract",
					GroupID:       "test_group",
					IncludeID:     "inc_123456",
					IncludeName:   "test_include",
					IncludeType:   "MICROSERVICES",
					LatestVersion: 1,
					PropertyType:  tools.StringPtr("INCLUDE"),
				},
			},
		},
		"500 internal server error": {
			params: GetIncludeRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_123456",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error getting include",
    "status": 500
}`,
			expectedPath: "/papi/v1/includes/inc_123456?contractId=test_contract&groupId=test_group",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: GetIncludeRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing contractId": {
			params: GetIncludeRequest{
				GroupID:   "test_group",
				IncludeID: "inc_123456",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing groupId": {
			params: GetIncludeRequest{
				ContractID: "test_contract",
				IncludeID:  "inc_123456",
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
			result, err := client.GetInclude(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreateInclude(t *testing.T) {
	tests := map[string]struct {
		params              CreateIncludeRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		responseHeaders     map[string]string
		expectedPath        string
		expectedResponse    *CreateIncludeResponse
		withError           error
	}{
		"200 OK - create include": {
			expectedRequestBody: `{"includeName":"test_include","includeType":"MICROSERVICES","productId":"test_product","ruleFormat":"test_rule_format"}`,
			responseStatus:      http.StatusCreated,
			responseBody: `
{
	"includeLink": "/papi/v1/includes/inc_123456?contractId=test_contract&groupId=test_group"
}`,
			responseHeaders: map[string]string{
				"x-limit-includes-per-contract-limit":     "500",
				"x-limit-includes-per-contract-remaining": "499",
			},
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeName: "test_include",
				IncludeType: IncludeTypeMicroServices,
				ProductID:   "test_product",
				RuleFormat:  "test_rule_format",
			},
			expectedPath: "/papi/v1/includes?contractId=test_contract&groupId=test_group",
			expectedResponse: &CreateIncludeResponse{
				IncludeID:   "inc_123456",
				IncludeLink: "/papi/v1/includes/inc_123456?contractId=test_contract&groupId=test_group",
				ResponseHeaders: CreateIncludeResponseHeaders{
					IncludesLimitTotal:     "500",
					IncludesLimitRemaining: "499",
				},
			},
		},
		"200 OK - create include with clone": {
			expectedRequestBody: `{"includeName":"test_include","includeType":"MICROSERVICES","productId":"test_product","cloneFrom":{"includeId":"inc_456789","version":1}}`,
			responseStatus:      http.StatusCreated,
			responseBody: `
{
	"includeLink": "/papi/v1/includes/inc_123456?contractId=test_contract&groupId=test_group"
}`,
			responseHeaders: map[string]string{
				"x-limit-includes-per-contract-limit":     "700",
				"x-limit-includes-per-contract-remaining": "654",
			},
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeName: "test_include",
				IncludeType: IncludeTypeMicroServices,
				ProductID:   "test_product",
				CloneIncludeFrom: &CloneIncludeFrom{
					IncludeID: "inc_456789",
					Version:   1,
				},
			},
			expectedPath: "/papi/v1/includes?contractId=test_contract&groupId=test_group",
			expectedResponse: &CreateIncludeResponse{
				IncludeID:   "inc_123456",
				IncludeLink: "/papi/v1/includes/inc_123456?contractId=test_contract&groupId=test_group",
				ResponseHeaders: CreateIncludeResponseHeaders{
					IncludesLimitTotal:     "700",
					IncludesLimitRemaining: "654",
				},
			},
		},
		"500 internal server error": {
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeName: "test_include",
				IncludeType: IncludeTypeMicroServices,
				ProductID:   "test_product",
				RuleFormat:  "test_rule_format",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating include",
    "status": 500
}`,
			expectedPath: "/papi/v1/includes?contractId=test_contract&groupId=test_group",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing productId": {
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeName: "test_include",
				IncludeType: IncludeTypeMicroServices,
				RuleFormat:  "test_rule_format",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing contractId": {
			params: CreateIncludeRequest{
				GroupID:     "test_group",
				IncludeName: "test_include",
				IncludeType: IncludeTypeMicroServices,
				ProductID:   "test_product",
				RuleFormat:  "test_rule_format",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing groupId": {
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				IncludeName: "test_include",
				IncludeType: IncludeTypeMicroServices,
				ProductID:   "test_product",
				RuleFormat:  "test_rule_format",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing includeName": {
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeType: IncludeTypeMicroServices,
				ProductID:   "test_product",
				RuleFormat:  "test_rule_format",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing includeType": {
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeName: "test_include",
				ProductID:   "test_product",
				RuleFormat:  "test_rule_format",
			},
			withError: ErrStructValidation,
		},
		"validation error - incorrect includeType": {
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeName: "test_include",
				IncludeType: "test",
				ProductID:   "test_product",
				RuleFormat:  "test_rule_format",
			},
			withError: ErrStructValidation,
		},
		"validation error - cloneFrom - missing includeId": {
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeName: "test_include",
				IncludeType: IncludeTypeMicroServices,
				ProductID:   "test_product",
				CloneIncludeFrom: &CloneIncludeFrom{
					Version: 1,
				},
			},
			withError: ErrStructValidation,
		},
		"validation error - cloneFrom - missing version": {
			params: CreateIncludeRequest{
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeName: "test_include",
				IncludeType: IncludeTypeMicroServices,
				ProductID:   "test_product",
				CloneIncludeFrom: &CloneIncludeFrom{
					IncludeID: "inc_123456",
				},
			},
			withError: ErrStructValidation,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)

				if len(test.responseHeaders) > 0 {
					for header, value := range test.responseHeaders {
						w.Header().Set(header, value)
					}
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateInclude(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteInclude(t *testing.T) {
	tests := map[string]struct {
		params           DeleteIncludeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DeleteIncludeResponse
		withError        error
	}{
		"200 OK - delete given includeId, contractId and groupId": {
			responseStatus: http.StatusOK,
			responseBody: `
{
	"message": "Deletion Successful."
}`,
			params: DeleteIncludeRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_123456",
			},
			expectedPath: "/papi/v1/includes/inc_123456?contractId=test_contract&groupId=test_group",
			expectedResponse: &DeleteIncludeResponse{
				Message: "Deletion Successful.",
			},
		},
		"200 OK - delete given only includeId": {
			responseStatus: http.StatusOK,
			responseBody: `
{
	"message": "Deletion Successful."
}`,
			params: DeleteIncludeRequest{
				IncludeID: "inc_123456",
			},
			expectedPath: "/papi/v1/includes/inc_123456",
			expectedResponse: &DeleteIncludeResponse{
				Message: "Deletion Successful.",
			},
		},
		"500 internal server error": {
			params: DeleteIncludeRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_123456",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"detail": "Error deleting include",
	"status": 500
}`,
			expectedPath: "/papi/v1/includes/inc_123456?contractId=test_contract&groupId=test_group",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: DeleteIncludeRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
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
			result, err := client.DeleteInclude(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
