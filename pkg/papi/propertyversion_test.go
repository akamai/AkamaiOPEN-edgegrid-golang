package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapi_GetPropertyVersions(t *testing.T) {
	tests := map[string]struct {
		params           GetPropertyVersionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPropertyVersionsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetPropertyVersionsRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
				Limit:      5,
				Offset:     6,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "propertyId": "propertyID",
    "propertyName": "property_name",
    "accountId": "accountID",
    "contractId": "contract",
    "groupId": "group",
    "assetId": "assetID",
    "versions": {
        "items": [
            {
                "propertyVersion": 2,
                "updatedByUser": "user",
                "updatedDate": "2020-09-14T19:06:13Z",
                "productionStatus": "INACTIVE",
                "stagingStatus": "ACTIVE",
                "etag": "etag",
                "productId": "productID",
                "note": "version note"
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions?contractId=contract&groupId=group&limit=5&offset=6",
			expectedResponse: &GetPropertyVersionsResponse{
				PropertyID:   "propertyID",
				PropertyName: "property_name",
				AccountID:    "accountID",
				ContractID:   "contract",
				GroupID:      "group",
				AssetID:      "assetID",
				Versions: PropertyVersionItems{
					Items: []PropertyVersionGetItem{
						{
							Etag:             "etag",
							Note:             "version note",
							ProductID:        "productID",
							ProductionStatus: "INACTIVE",
							PropertyVersion:  2,
							StagingStatus:    "ACTIVE",
							UpdatedByUser:    "user",
							UpdatedDate:      "2020-09-14T19:06:13Z",
						}}},
			},
		},
		"500 Internal Server Error": {
			params: GetPropertyVersionsRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
				Limit:      5,
				Offset:     6,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching property versions",
    "status": 505
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions?contractId=contract&groupId=group&limit=5&offset=6",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching property versions",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty property ID": {
			params: GetPropertyVersionsRequest{
				PropertyID: "",
				ContractID: "contract",
				GroupID:    "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
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
			result, err := client.GetPropertyVersions(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_GetPropertyVersion(t *testing.T) {
	tests := map[string]struct {
		params           GetPropertyVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPropertyVersionsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetPropertyVersionRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "propertyId": "propertyID",
    "propertyName": "property_name",
    "accountId": "accountID",
    "contractId": "contract",
    "groupId": "group",
    "assetId": "assetID",
    "versions": {
        "items": [
            {
                "propertyVersion": 2,
                "updatedByUser": "user",
                "updatedDate": "2020-09-14T19:06:13Z",
                "productionStatus": "INACTIVE",
                "stagingStatus": "ACTIVE",
                "etag": "etag",
                "productId": "productID",
                "note": "version note"
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2?contractId=contract&groupId=group",
			expectedResponse: &GetPropertyVersionsResponse{
				PropertyID:   "propertyID",
				PropertyName: "property_name",
				AccountID:    "accountID",
				ContractID:   "contract",
				GroupID:      "group",
				AssetID:      "assetID",
				Versions: PropertyVersionItems{
					Items: []PropertyVersionGetItem{
						{
							Etag:             "etag",
							Note:             "version note",
							ProductID:        "productID",
							ProductionStatus: "INACTIVE",
							PropertyVersion:  2,
							StagingStatus:    "ACTIVE",
							UpdatedByUser:    "user",
							UpdatedDate:      "2020-09-14T19:06:13Z",
						}}},
				Version: PropertyVersionGetItem{

					Etag:             "etag",
					Note:             "version note",
					ProductID:        "productID",
					ProductionStatus: "INACTIVE",
					PropertyVersion:  2,
					StagingStatus:    "ACTIVE",
					UpdatedByUser:    "user",
					UpdatedDate:      "2020-09-14T19:06:13Z",
				},
			},
		},
		"version not found": {
			params: GetPropertyVersionRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "propertyId": "propertyID",
    "propertyName": "property_name",
    "accountId": "accountID",
    "contractId": "contract",
    "groupId": "group",
    "assetId": "assetID",
    "versions": {
        "items": [
        ]
    }
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrNotFound), "want: %s; got: %s", ErrNotFound, err)
			},
		},
		"500 Internal Server Error": {
			params: GetPropertyVersionRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching property version",
    "status": 505
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching property version",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty property ID": {
			params: GetPropertyVersionRequest{
				PropertyID:      "",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"empty property version": {
			params: GetPropertyVersionRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyVersion")
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
			result, err := client.GetPropertyVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_CreatePropertyVersion(t *testing.T) {
	tests := map[string]struct {
		params           CreatePropertyVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreatePropertyVersionResponse
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			params: CreatePropertyVersionRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
				Version: PropertyVersionCreate{
					CreateFromVersion: 1,
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
		{
		   "versionLink": "/papi/v1/properties/propertyID/versions/2?contractId=contract&groupId=group"
		}`,
			expectedPath: "/papi/v1/properties/propertyID/versions?contractId=contract&groupId=group",
			expectedResponse: &CreatePropertyVersionResponse{
				VersionLink:     "/papi/v1/properties/propertyID/versions/2?contractId=contract&groupId=group",
				PropertyVersion: 2,
			},
		},
		"500 Internal Server Error": {
			params: CreatePropertyVersionRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
				Version: PropertyVersionCreate{
					CreateFromVersion: 1,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error creating property version",
		   "status": 500
		}`,
			expectedPath: "/papi/v1/properties/propertyID/versions?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating property version",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty property ID": {
			params: CreatePropertyVersionRequest{
				PropertyID: "",
				ContractID: "contract",
				GroupID:    "group",
				Version: PropertyVersionCreate{
					CreateFromVersion: 1,
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"empty CreateFromVersion": {
			params: CreatePropertyVersionRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
				Version:    PropertyVersionCreate{},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "CreateFromVersion")
			},
		},
		"invalid location": {
			params: CreatePropertyVersionRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
				Version: PropertyVersionCreate{
					CreateFromVersion: 1,
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "versionLink": ":"
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := ErrInvalidResponseLink
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"invalid version format": {
			params: CreatePropertyVersionRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
				Version: PropertyVersionCreate{
					CreateFromVersion: 1,
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "versionLink": "/papi/v1/properties/propertyID/versions/abc?contractId=contract&groupId=group"
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := ErrInvalidResponseLink
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
			result, err := client.CreatePropertyVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_GetLatestVersion(t *testing.T) {
	tests := map[string]struct {
		params           GetLatestVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPropertyVersionsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetLatestVersionRequest{
				PropertyID:  "propertyID",
				ActivatedOn: "STAGING",
				ContractID:  "contract",
				GroupID:     "group",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/latest?activatedOn=STAGING&contractId=contract&groupId=group",
			responseBody: `
{
    "propertyId": "propertyID",
    "propertyName": "property_name",
    "accountId": "accountID",
    "contractId": "contract",
    "groupId": "group",
    "assetId": "assetID",
    "versions": {
        "items": [
            {
                "propertyVersion": 2,
                "updatedByUser": "user",
                "updatedDate": "2020-09-14T19:06:13Z",
                "productionStatus": "INACTIVE",
                "stagingStatus": "ACTIVE",
                "etag": "etag",
                "productId": "productID",
                "note": "version note"
            }
        ]
    }
}`,
			expectedResponse: &GetPropertyVersionsResponse{
				PropertyID:   "propertyID",
				PropertyName: "property_name",
				AccountID:    "accountID",
				ContractID:   "contract",
				GroupID:      "group",
				AssetID:      "assetID",
				Versions: PropertyVersionItems{
					Items: []PropertyVersionGetItem{
						{
							Etag:             "etag",
							Note:             "version note",
							ProductID:        "productID",
							ProductionStatus: "INACTIVE",
							PropertyVersion:  2,
							StagingStatus:    "ACTIVE",
							UpdatedByUser:    "user",
							UpdatedDate:      "2020-09-14T19:06:13Z",
						},
					},
				},
				Version: PropertyVersionGetItem{
					Etag:             "etag",
					Note:             "version note",
					ProductID:        "productID",
					ProductionStatus: "INACTIVE",
					PropertyVersion:  2,
					StagingStatus:    "ACTIVE",
					UpdatedByUser:    "user",
					UpdatedDate:      "2020-09-14T19:06:13Z",
				},
			},
		},
		"Version not found": {
			params: GetLatestVersionRequest{
				PropertyID:  "propertyID",
				ActivatedOn: "STAGING",
				ContractID:  "contract",
				GroupID:     "group",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/latest?activatedOn=STAGING&contractId=contract&groupId=group",
			responseBody: `
{
    "propertyId": "propertyID",
    "propertyName": "property_name",
    "accountId": "accountID",
    "contractId": "contract",
    "groupId": "group",
    "assetId": "assetID",
    "versions": {
        "items": [
        ]
    }
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrNotFound), "want: %v; got: %v", ErrNotFound, err)
			},
		},
		"500 Internal Server Error": {
			params: GetLatestVersionRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error fetching latest version",
  "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/latest?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching latest version",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty property ID": {
			params: GetLatestVersionRequest{
				PropertyID:  "",
				ActivatedOn: "STAGING",
				ContractID:  "contract",
				GroupID:     "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"invalid ActivatedOn": {
			params: GetLatestVersionRequest{
				PropertyID:  "propertyID",
				ActivatedOn: "test",
				ContractID:  "contract",
				GroupID:     "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ActivatedOn")
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
			result, err := client.GetLatestVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_GetAvailableBehaviors(t *testing.T) {
	tests := map[string]struct {
		params           GetFeaturesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetFeaturesCriteriaResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetFeaturesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/2/available-behaviors?contractId=contract&groupId=group",
			responseBody: `
{
    "contractId": "contract",
    "groupId": "group",
    "productId": "productID",
    "ruleFormat": "v2020-09-15",
    "availableBehaviors": {
        "items": [
            {
                "name": "cpCode",
                "schemaLink": "/papi/v1/schemas/products/prd_Alta/latest#/definitions/catalog/behaviors/cpCode"
            }
        ]
    }
}`,
			expectedResponse: &GetFeaturesCriteriaResponse{
				ContractID: "contract",
				GroupID:    "group",
				ProductID:  "productID",
				RuleFormat: "v2020-09-15",
				AvailableBehaviors: AvailableFeatureItems{Items: []AvailableFeature{
					{
						Name:       "cpCode",
						SchemaLink: "/papi/v1/schemas/products/prd_Alta/latest#/definitions/catalog/behaviors/cpCode",
					},
				}},
			},
		},
		"500 Internal Server Error": {
			params: GetFeaturesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error fetching available behaviors",
  "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2/available-behaviors?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching available behaviors",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty property ID": {
			params: GetFeaturesRequest{
				PropertyID:      "",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"empty property version": {
			params: GetFeaturesRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyVersion")
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
			result, err := client.GetAvailableBehaviors(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_GetAvailableCriteria(t *testing.T) {
	tests := map[string]struct {
		params           GetFeaturesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetFeaturesCriteriaResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetFeaturesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/2/available-criteria?contractId=contract&groupId=group",
			responseBody: `
{
    "contractId": "contract",
    "groupId": "group",
    "productId": "productID",
    "ruleFormat": "v2020-09-15",
    "availableBehaviors": {
        "items": [
            {
                "name": "cpCode",
                "schemaLink": "/papi/v1/schemas/products/prd_Alta/latest#/definitions/catalog/behaviors/cpCode"
            }
        ]
    }
}`,
			expectedResponse: &GetFeaturesCriteriaResponse{
				ContractID: "contract",
				GroupID:    "group",
				ProductID:  "productID",
				RuleFormat: "v2020-09-15",
				AvailableBehaviors: AvailableFeatureItems{Items: []AvailableFeature{
					{
						Name:       "cpCode",
						SchemaLink: "/papi/v1/schemas/products/prd_Alta/latest#/definitions/catalog/behaviors/cpCode",
					},
				}},
			},
		},
		"500 Internal Server Error": {
			params: GetFeaturesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error fetching available behaviors",
  "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2/available-criteria?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching available behaviors",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty property ID": {
			params: GetFeaturesRequest{
				PropertyID:      "",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"empty property version": {
			params: GetFeaturesRequest{
				PropertyID: "propertyID",
				ContractID: "contract",
				GroupID:    "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyVersion")
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
			result, err := client.GetAvailableCriteria(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
