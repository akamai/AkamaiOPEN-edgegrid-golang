package papi

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateIncludeVersion(t *testing.T) {
	tests := map[string]struct {
		params              CreateIncludeVersionRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateIncludeVersionResponse
		withError           error
	}{
		"201 Created": {
			params: CreateIncludeVersionRequest{
				IncludeID: "inc_12345",
				IncludeVersionRequest: IncludeVersionRequest{
					CreateFromVersion: 2,
				},
			},
			expectedRequestBody: `{"createFromVersion":2}`,
			expectedPath:        "/papi/v1/includes/inc_12345/versions",
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "versionLink": "/papi/v1/includes/inc_12345/versions/5"
}`,
			expectedResponse: &CreateIncludeVersionResponse{
				VersionLink: "/papi/v1/includes/inc_12345/versions/5",
				Version:     5,
			},
		},
		"500 internal server error": {
			params: CreateIncludeVersionRequest{
				IncludeID: "inc_12345",
				IncludeVersionRequest: IncludeVersionRequest{
					CreateFromVersion: 2,
				},
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error getting include",
    "status": 500
}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: CreateIncludeVersionRequest{
				IncludeVersionRequest: IncludeVersionRequest{
					CreateFromVersion: 2,
				},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing createFromVersion": {
			params: CreateIncludeVersionRequest{
				IncludeID: "inc_12345",
			},
			withError: ErrStructValidation,
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateIncludeVersion(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetIncludeVersion(t *testing.T) {
	tests := map[string]struct {
		params           GetIncludeVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetIncludeVersionResponse
		withError        error
	}{
		"200 OK": {
			params: GetIncludeVersionRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				Version:    2,
				IncludeID:  "inc_12345",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions/2?contractId=test_contract&groupId=test_group",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "includeId": "inc_12345",
    "includeName": "tfp_test1",
    "accountId": "act_B-3-WNKA123",
    "contractId": "test_contract",
    "groupId": "test_group",
    "assetId": "aid_11069123",
    "includeType": "MICROSERVICES",
    "versions": {
        "items": [
            {
                "updatedByUser": "test_user",
                "updatedDate": "2022-08-22T07:17:48Z",
                "productionStatus": "INACTIVE",
                "stagingStatus": "ACTIVE",
                "etag": "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
                "productId": "prd_Site_Defender",
                "ruleFormat": "v2020-11-02",
                "includeVersion": 2
            }
        ]
    }
}`,
			expectedResponse: &GetIncludeVersionResponse{
				AccountID:   "act_B-3-WNKA123",
				AssetID:     "aid_11069123",
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeID:   "inc_12345",
				IncludeName: "tfp_test1",
				IncludeType: IncludeTypeMicroServices,
				IncludeVersions: Versions{
					Items: []IncludeVersion{
						{
							UpdatedByUser:    "test_user",
							UpdatedDate:      "2022-08-22T07:17:48Z",
							ProductionStatus: VersionStatusInactive,
							Etag:             "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
							ProductID:        "prd_Site_Defender",
							RuleFormat:       "v2020-11-02",
							IncludeVersion:   2,
							StagingStatus:    VersionStatusActive,
						},
					},
				},
				IncludeVersion: IncludeVersion{
					UpdatedByUser:    "test_user",
					UpdatedDate:      "2022-08-22T07:17:48Z",
					ProductionStatus: VersionStatusInactive,
					Etag:             "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
					ProductID:        "prd_Site_Defender",
					RuleFormat:       "v2020-11-02",
					IncludeVersion:   2,
					StagingStatus:    VersionStatusActive,
				},
			},
		},
		"500 internal server error": {
			params: GetIncludeVersionRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				Version:    2,
				IncludeID:  "inc_12345",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions/2?contractId=test_contract&groupId=test_group",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error getting include",
    "status": 500
}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: GetIncludeVersionRequest{
				Version:    1,
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing contractId": {
			params: GetIncludeVersionRequest{
				IncludeID: "inc_12345",
				Version:   1,
				GroupID:   "test_group",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing groupId": {
			params: GetIncludeVersionRequest{
				IncludeID:  "inc_12345",
				ContractID: "test_contract",
				Version:    1,
			},
			withError: ErrStructValidation,
		},
		"validation error - missing version": {
			params: GetIncludeVersionRequest{
				IncludeID:  "inc_12345",
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
			result, err := client.GetIncludeVersion(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListIncludeVersions(t *testing.T) {
	tests := map[string]struct {
		params           ListIncludeVersionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListIncludeVersionsResponse
		withError        error
	}{
		"200 OK": {
			params: ListIncludeVersionsRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_12345",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions?contractId=test_contract&groupId=test_group",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "includeId": "inc_12345",
    "includeName": "tfp_test1",
    "accountId": "act_B-3-WNKA123",
    "contractId": "test_contract",
    "groupId": "test_group",
    "assetId": "aid_11069123",
    "includeType": "MICROSERVICES",
    "versions": {
        "items": [
			{
                "updatedByUser": "test_user",
                "updatedDate": "2022-10-14T08:41:00Z",
                "productionStatus": "INACTIVE",
                "stagingStatus": "INACTIVE",
                "etag": "c925d2b5fa4cc002774c752186d8faafeac7f28a",
                "productId": "prd_Site_Defender",
                "ruleFormat": "v2020-11-02",
                "includeVersion": 4
            },
            {
                "updatedByUser": "test_user",
                "updatedDate": "2022-08-23T12:39:33Z",
                "productionStatus": "INACTIVE",
                "stagingStatus": "INACTIVE",
                "etag": "f5230dfe9d50e0a4a8b643388226b36db494d7c4",
                "productId": "prd_Site_Defender",
                "ruleFormat": "v2020-11-02",
                "includeVersion": 3
            },
            {
                "updatedByUser": "test_user",
                "updatedDate": "2022-08-22T07:17:48Z",
                "productionStatus": "INACTIVE",
                "stagingStatus": "ACTIVE",
                "etag": "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
                "productId": "prd_Site_Defender",
                "ruleFormat": "v2020-11-02",
                "includeVersion": 2
            },
            {
                "updatedByUser": "test_user",
                "updatedDate": "2022-08-16T10:29:43Z",
                "productionStatus": "INACTIVE",
                "stagingStatus": "DEACTIVATED",
                "etag": "d2be894768ae4e587eae91f93f15d2217ef517d8",
                "productId": "prd_Site_Defender",
                "ruleFormat": "v2020-11-02",
                "includeVersion": 1
            }
        ]
    }
}`,
			expectedResponse: &ListIncludeVersionsResponse{
				AccountID:   "act_B-3-WNKA123",
				AssetID:     "aid_11069123",
				ContractID:  "test_contract",
				GroupID:     "test_group",
				IncludeID:   "inc_12345",
				IncludeName: "tfp_test1",
				IncludeType: IncludeTypeMicroServices,
				IncludeVersions: Versions{
					Items: []IncludeVersion{
						{
							UpdatedByUser:    "test_user",
							UpdatedDate:      "2022-10-14T08:41:00Z",
							ProductionStatus: VersionStatusInactive,
							Etag:             "c925d2b5fa4cc002774c752186d8faafeac7f28a",
							ProductID:        "prd_Site_Defender",
							RuleFormat:       "v2020-11-02",
							IncludeVersion:   4,
							StagingStatus:    VersionStatusInactive,
						},
						{
							UpdatedByUser:    "test_user",
							UpdatedDate:      "2022-08-23T12:39:33Z",
							ProductionStatus: VersionStatusInactive,
							Etag:             "f5230dfe9d50e0a4a8b643388226b36db494d7c4",
							ProductID:        "prd_Site_Defender",
							RuleFormat:       "v2020-11-02",
							IncludeVersion:   3,
							StagingStatus:    VersionStatusInactive,
						},
						{
							UpdatedByUser:    "test_user",
							UpdatedDate:      "2022-08-22T07:17:48Z",
							ProductionStatus: VersionStatusInactive,
							Etag:             "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
							ProductID:        "prd_Site_Defender",
							RuleFormat:       "v2020-11-02",
							IncludeVersion:   2,
							StagingStatus:    VersionStatusActive,
						},
						{
							UpdatedByUser:    "test_user",
							UpdatedDate:      "2022-08-16T10:29:43Z",
							ProductionStatus: VersionStatusInactive,
							Etag:             "d2be894768ae4e587eae91f93f15d2217ef517d8",
							ProductID:        "prd_Site_Defender",
							RuleFormat:       "v2020-11-02",
							IncludeVersion:   1,
							StagingStatus:    VersionStatusDeactivated,
						},
					},
				},
			},
		},
		"500 internal server error": {
			params: ListIncludeVersionsRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
				IncludeID:  "inc_12345",
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions?contractId=test_contract&groupId=test_group",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error getting include",
    "status": 500
}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: ListIncludeVersionsRequest{
				ContractID: "test_contract",
				GroupID:    "test_group",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing contractId": {
			params: ListIncludeVersionsRequest{
				GroupID:   "test_group",
				IncludeID: "inc_12345",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing groupId": {
			params: ListIncludeVersionsRequest{
				ContractID: "test_contract",
				IncludeID:  "inc_12345",
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
			result, err := client.ListIncludeVersions(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetIncludeVersionAvailableCriteria(t *testing.T) {
	tests := map[string]struct {
		params           ListAvailableCriteriaRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *AvailableCriteriaResponse
		withError        error
	}{
		"200 OK": {
			params: ListAvailableCriteriaRequest{
				IncludeID: "inc_12345",
				Version:   2,
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions/2/available-criteria",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "ctr_3-WNK123",
    "groupId": "grp_115123",
    "productId": "prd_Site_Defender",
    "ruleFormat": "v2020-11-02",
    "criteria": {
        "items": [
            {
                "name": "bucket",
                "schemaLink": "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fbucket"
            },
            {
                "name": "cacheability",
                "schemaLink": "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fcacheability"
            },
            {
                "name": "chinaCdnRegion",
                "schemaLink": "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fcriteria%2FchinaCdnRegion"
            }
        ]
    }
}`,
			expectedResponse: &AvailableCriteriaResponse{
				ContractID: "ctr_3-WNK123",
				GroupID:    "grp_115123",
				ProductID:  "prd_Site_Defender",
				RuleFormat: "v2020-11-02",
				AvailableCriteria: AvailableCriteria{
					Items: []Criteria{
						{
							Name:       "bucket",
							SchemaLink: "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fbucket",
						},
						{
							Name:       "cacheability",
							SchemaLink: "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fcacheability",
						},
						{
							Name:       "chinaCdnRegion",
							SchemaLink: "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fcriteria%2FchinaCdnRegion",
						},
					},
				},
			},
		},
		"500 internal server error": {
			params: ListAvailableCriteriaRequest{
				IncludeID: "inc_12345",
				Version:   2,
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions/2/available-criteria",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error getting include",
    "status": 500
}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: ListAvailableCriteriaRequest{
				Version: 2,
			},
			withError: ErrStructValidation,
		},
		"validation error - missing version": {
			params: ListAvailableCriteriaRequest{
				IncludeID: "inc_12345",
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
			result, err := client.ListIncludeVersionAvailableCriteria(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListIncludeVersionAvailableBehaviors(t *testing.T) {
	tests := map[string]struct {
		params           ListAvailableBehaviorsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *AvailableBehaviorsResponse
		withError        error
	}{
		"200 OK": {
			params: ListAvailableBehaviorsRequest{
				IncludeID: "inc_12345",
				Version:   2,
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions/2/available-behaviors",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "ctr_3-WNK123",
    "groupId": "grp_115123",
    "productId": "prd_Site_Defender",
    "ruleFormat": "v2020-11-02",
    "behaviors": {
        "items": [
            {
                "name": "akamaizer",
                "schemaLink": "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fbehaviors%2Fakamaizer"
            },
            {
                "name": "akamaizerTag",
                "schemaLink": "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FakamaizerTag"
            },
            {
                "name": "allHttpInCacheHierarchy",
                "schemaLink": "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallHttpInCacheHierarchy"
            },
            {
                "name": "allowDelete",
                "schemaLink": "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowDelete"
            }
        ]
    }
}`,
			expectedResponse: &AvailableBehaviorsResponse{
				ContractID: "ctr_3-WNK123",
				GroupID:    "grp_115123",
				ProductID:  "prd_Site_Defender",
				RuleFormat: "v2020-11-02",
				AvailableBehaviors: AvailableBehaviors{
					Items: []Behavior{
						{
							Name:       "akamaizer",
							SchemaLink: "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fbehaviors%2Fakamaizer",
						},
						{
							Name:       "akamaizerTag",
							SchemaLink: "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FakamaizerTag",
						},
						{
							Name:       "allHttpInCacheHierarchy",
							SchemaLink: "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallHttpInCacheHierarchy",
						},
						{
							Name:       "allowDelete",
							SchemaLink: "/papi/v0/schemas/products/prd_Site_Defender/v2020-11-02#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowDelete",
						},
					},
				},
			},
		},
		"500 internal server error": {
			params: ListAvailableBehaviorsRequest{
				IncludeID: "inc_12345",
				Version:   2,
			},
			expectedPath:   "/papi/v1/includes/inc_12345/versions/2/available-behaviors",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error getting include",
    "status": 500
}`,
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error getting include",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: ListAvailableBehaviorsRequest{
				Version: 2,
			},
			withError: ErrStructValidation,
		},
		"validation error - missing version": {
			params: ListAvailableBehaviorsRequest{
				IncludeID: "inc_12345",
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
			result, err := client.ListIncludeVersionAvailableBehaviors(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
