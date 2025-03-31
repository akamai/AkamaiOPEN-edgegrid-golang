package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiGetPropertyVersions(t *testing.T) {
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
		"200 OK - response with propertyType": {
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
	"propertyType": "HOSTNAME_BUCKET",
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
				PropertyType: ptr.To("HOSTNAME_BUCKET"),
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

func TestPapiGetPropertyVersion(t *testing.T) {
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
		"200 OK - response with propertyType": {
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
	"propertyType": "HOSTNAME_BUCKET",
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
				PropertyType: ptr.To("HOSTNAME_BUCKET"),
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

func TestPapiCreatePropertyVersion(t *testing.T) {
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

func TestPapiGetLatestVersion(t *testing.T) {
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
		"200 OK - response with propertyType": {
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
	"propertyType": "HOSTNAME_BUCKET",
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
				PropertyType: ptr.To("HOSTNAME_BUCKET"),
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

func TestPapiGetAvailableBehaviors(t *testing.T) {
	tests := map[string]struct {
		params           GetAvailableBehaviorsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetBehaviorsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetAvailableBehaviorsRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/2/available-behaviors",
			responseBody: `
{
    "contractId": "contract",
    "groupId": "group",
    "productId": "productID",
    "ruleFormat": "v2020-09-15",
    "behaviors": {
        "items": [
            {
                "name": "allHttpInCacheHierarchy",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallHttpInCacheHierarchy"
            },
            {
                "name": "allowDelete",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowDelete"
            },
            {
                "name": "allowOptions",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowOptions"
            }
        ]
    }
}`,
			expectedResponse: &GetBehaviorsResponse{
				ContractID: "contract",
				GroupID:    "group",
				ProductID:  "productID",
				RuleFormat: "v2020-09-15",
				AvailableBehaviors: AvailableBehaviors{Items: []Behavior{
					{
						Name:       "allHttpInCacheHierarchy",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallHttpInCacheHierarchy",
					},
					{
						Name:       "allowDelete",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowDelete",
					},
					{
						Name:       "allowOptions",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowOptions",
					},
				}},
			},
		},
		"200 OK with query parameters": {
			params: GetAvailableBehaviorsRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				GroupID:         "groupID",
				ContractID:      "contractID",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/2/available-behaviors?contractId=contractID&groupId=groupID",
			responseBody: `
{
    "contractId": "contract",
    "groupId": "group",
    "productId": "productID",
    "ruleFormat": "v2020-09-15",
    "behaviors": {
        "items": [
            {
                "name": "allHttpInCacheHierarchy",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallHttpInCacheHierarchy"
            },
            {
                "name": "allowDelete",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowDelete"
            },
            {
                "name": "allowOptions",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowOptions"
            }
        ]
    }
}`,
			expectedResponse: &GetBehaviorsResponse{
				ContractID: "contract",
				GroupID:    "group",
				ProductID:  "productID",
				RuleFormat: "v2020-09-15",
				AvailableBehaviors: AvailableBehaviors{Items: []Behavior{
					{
						Name:       "allHttpInCacheHierarchy",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallHttpInCacheHierarchy",
					},
					{
						Name:       "allowDelete",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowDelete",
					},
					{
						Name:       "allowOptions",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fbehaviors%2FallowOptions",
					},
				}},
			},
		},
		"500 Internal Server Error": {
			params: GetAvailableBehaviorsRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error fetching available behaviors",
  "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2/available-behaviors",
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
			params: GetAvailableBehaviorsRequest{
				PropertyID:      "",
				PropertyVersion: 2,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"empty property version": {
			params: GetAvailableBehaviorsRequest{
				PropertyID: "propertyID",
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

func TestPapiGetAvailableCriteria(t *testing.T) {
	tests := map[string]struct {
		params           GetAvailableCriteriaRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCriteriaResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetAvailableCriteriaRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/2/available-criteria",
			responseBody: `
{
    "contractId": "contract",
    "groupId": "group",
    "productId": "productID",
    "ruleFormat": "v2020-09-15",
    "criteria": {
        "items": [
            {
                "name": "bucket",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fbucket"
            },
            {
                "name": "cacheability",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fcacheability"
            },
            {
                "name": "chinaCdnRegion",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2FchinaCdnRegion"
            }
        ]
    }
}`,
			expectedResponse: &GetCriteriaResponse{
				ContractID: "contract",
				GroupID:    "group",
				ProductID:  "productID",
				RuleFormat: "v2020-09-15",
				AvailableCriteria: AvailableCriteria{Items: []Criteria{
					{
						Name:       "bucket",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fbucket",
					},
					{
						Name:       "cacheability",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fcacheability",
					},
					{
						Name:       "chinaCdnRegion",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2FchinaCdnRegion",
					},
				}},
			},
		},
		"200 OK with query parameters": {
			params: GetAvailableCriteriaRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				GroupID:         "groupID",
				ContractID:      "contractID",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/2/available-criteria?contractId=contractID&groupId=groupID",
			responseBody: `
{
    "contractId": "contract",
    "groupId": "group",
    "productId": "productID",
    "ruleFormat": "v2020-09-15",
    "criteria": {
        "items": [
            {
                "name": "bucket",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fbucket"
            },
            {
                "name": "cacheability",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fcacheability"
            },
            {
                "name": "chinaCdnRegion",
                "schemaLink": "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2FchinaCdnRegion"
            }
        ]
    }
}`,
			expectedResponse: &GetCriteriaResponse{
				ContractID: "contract",
				GroupID:    "group",
				ProductID:  "productID",
				RuleFormat: "v2020-09-15",
				AvailableCriteria: AvailableCriteria{Items: []Criteria{
					{
						Name:       "bucket",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fbucket",
					},
					{
						Name:       "cacheability",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2Fcacheability",
					},
					{
						Name:       "chinaCdnRegion",
						SchemaLink: "/papi/v0/schemas/products/prd_Rich_Media_Accel/latest#%2Fdefinitions%2Fcatalog%2Fcriteria%2FchinaCdnRegion",
					},
				}},
			},
		},
		"500 Internal Server Error": {
			params: GetAvailableCriteriaRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error fetching available behaviors",
  "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2/available-criteria",
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
			params: GetAvailableCriteriaRequest{
				PropertyID:      "",
				PropertyVersion: 2,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"empty property version": {
			params: GetAvailableCriteriaRequest{
				PropertyID: "propertyID",
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

func TestPapiListAvailableIncludes(t *testing.T) {
	tests := map[string]struct {
		params           ListAvailableIncludesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListAvailableIncludesResponse
		withError        error
	}{
		"200 OK - available includes given ContractID and GroupID": {
			params: ListAvailableIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
				ContractID:      "ctr_1",
				GroupID:         "grp_2",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/1/external-resources?contractId=ctr_1&groupId=grp_2",
			responseBody: `
{
    "externalResources": {
        "include": {
            "test_include_id1": {
                "id": "test_include_id1",
                "name": "test_include1",
                "includeType": "MICROSERVICES",
                "fileName": "test_include1.xml",
                "productName": "Example_Name",
                "ruleFormat": "v2020-11-02"
            },
			"test_include_id2": {
                "id": "test_include_id2",
                "name": "test_include2",
                "includeType": "MICROSERVICES",
                "fileName": "test_include2.xml",
                "productName": "Example_Name",
                "ruleFormat": "v2020-11-02"
            }
		},
		"cloudletSharedPolicy": {
			"123456": {
                "cloudletType": "TESTCLOUDLETTYPE",
                "id": 123456,
                "name": "TestName123456",
                "policyType": "SHARED"
            }
		},
		"availableCnames": [
			{
                "id": 123456,
                "name": "www.example.com",
                "domain": "example.net",
                "serialNumber": 123,
                "slot": null,
                "status": "Created",
                "ipv6": false,
                "useCases": [],
                "cname": "www.example.example.net",
                "isSecure": false,
                "isEdgeIPBindingEnabled": null
            }
		],
		"customOverrides": {},
		"customOverrides": {},
		"blacklistedCertDomains": [
			"s3.example.com"
		],
		"availableNetStorageGroups": [
			{
                "id": 123456,
                "name": "aa-example",
                "uploadDomainName": "spm.example.example.com",
                "downloadDomainName": "spm.example.example.com",
                "cpCodeList": [
                    {
                        "cpCode": 123456,
                        "g2oToken": null
                    }
                ]
            }
		],
    	"availableCpCodes": [
			{
                "id": 123123,
                "description": "example-test-subgroup",
                "products": [
                    "EXAMPLE"
                ],
                "createdDate": 1521566901000,
                "cpCodeLimits": null,
                "name": "example-test-subgroup"
            }
		],
    	"availablePolicies": {
			"applicationLoadBalancer":[
				{
                    "id": 123456,
                    "name": "0000000000000000_EXAMPLE_clone"
                }
			],
			"firstPartyMarketingPlus":[
				{
                    "id": 123456,
                    "name": "Example_Name_123456"
                }
			],
			"firstPartyMarketing":[
				{
                    "id": 123456,
                    "name": "Example_first_party"
                }	
			],
			"forwardRewrite":[
				{
                    "id": 123456,
                    "name": "ExampleName"
                }
			],
			"continuousDeployment":[
				{
                    "id": 123456,
                    "name": "Example_Name_123456"
                }
			],
			"requestControl":[
				{
                    "id": 123456,
                    "name": "EXAMPLE_MATCH_RULE_SIZE_RC"
                }
			],
			"inputValidation":[
				{
                    "id": 123456,
                    "name": "example_name"
                }
			],
			"visitorPrioritization":[
				{
                    "id": 123456,
                    "name": "0000000000000000000000_EXAMPLE"
                }
			],
			"audienceSegmentation":[
				{
                    "id": 123456,
                    "name": "ExampleTest"
                }
			],
			"apiPrioritization":[
				{
                    "id": 123456,
                    "name": "APIExampleTest"
                }
			],
			"edgeRedirector":[
				{
                    "id": 123456,
                    "name": "00000000_EXAMPLE"
                }
			]
		},
    	"cloudletSharedPolicyVirtualWaitingRoom": {
			"123456": {
                "cloudletType": "EXAMPLETYPE",
                "id": 123456,
                "name": "example_name",
                "policyType": "SHARED"
            }
		}
    }
}`,
			expectedResponse: &ListAvailableIncludesResponse{
				AvailableIncludes: []ExternalIncludeData{

					{
						IncludeID:   "test_include_id1",
						IncludeName: "test_include1",
						IncludeType: IncludeTypeMicroServices,
						FileName:    "test_include1.xml",
						ProductName: "Example_Name",
						RuleFormat:  "v2020-11-02",
					},
					{
						IncludeID:   "test_include_id2",
						IncludeName: "test_include2",
						IncludeType: IncludeTypeMicroServices,
						FileName:    "test_include2.xml",
						ProductName: "Example_Name",
						RuleFormat:  "v2020-11-02",
					},
				},
			},
		},
		"200 OK - available includes given only ContractID": {
			params: ListAvailableIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
				ContractID:      "ctr_1",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/1/external-resources?contractId=ctr_1",
			responseBody: `
{
    "externalResources": {
        "include": {
            "test_include_id1": {
                "id": "test_include_id1",
                "name": "test_include1",
                "includeType": "MICROSERVICES",
                "fileName": "test_include1.xml",
                "productName": "Example_Name",
                "ruleFormat": "v2020-11-02"
            },
			"test_include_id2": {
                "id": "test_include_id2",
                "name": "test_include2",
                "includeType": "MICROSERVICES",
                "fileName": "test_include2.xml",
                "productName": "Example_Name",
                "ruleFormat": "v2020-11-02"
            }
		},
		"cloudletSharedPolicy": {
			"123456": {
                "cloudletType": "TESTCLOUDLETTYPE",
                "id": 123456,
                "name": "TestName123456",
                "policyType": "SHARED"
            }
		},
		"availableCnames": [
			{
                "id": 123456,
                "name": "www.example.com",
                "domain": "example.net",
                "serialNumber": 123,
                "slot": null,
                "status": "Created",
                "ipv6": false,
                "useCases": [],
                "cname": "www.example.example.net",
                "isSecure": false,
                "isEdgeIPBindingEnabled": null
            }
		],
		"customOverrides": {},
		"customOverrides": {},
		"blacklistedCertDomains": [
			"s3.example.com"
		],
		"availableNetStorageGroups": [
			{
                "id": 123456,
                "name": "aa-example",
                "uploadDomainName": "spm.example.example.com",
                "downloadDomainName": "spm.example.example.com",
                "cpCodeList": [
                    {
                        "cpCode": 123456,
                        "g2oToken": null
                    }
                ]
            }
		],
    	"availableCpCodes": [
			{
                "id": 123123,
                "description": "example-test-subgroup",
                "products": [
                    "EXAMPLE"
                ],
                "createdDate": 1521566901000,
                "cpCodeLimits": null,
                "name": "example-test-subgroup"
            }
		],
    	"availablePolicies": {
			"applicationLoadBalancer":[
				{
                    "id": 123456,
                    "name": "0000000000000000_EXAMPLE_clone"
                }
			],
			"firstPartyMarketingPlus":[
				{
                    "id": 123456,
                    "name": "Example_Name_123456"
                }
			],
			"firstPartyMarketing":[
				{
                    "id": 123456,
                    "name": "Example_first_party"
                }	
			],
			"forwardRewrite":[
				{
                    "id": 123456,
                    "name": "ExampleName"
                }
			],
			"continuousDeployment":[
				{
                    "id": 123456,
                    "name": "Example_Name_123456"
                }
			],
			"requestControl":[
				{
                    "id": 123456,
                    "name": "EXAMPLE_MATCH_RULE_SIZE_RC"
                }
			],
			"inputValidation":[
				{
                    "id": 123456,
                    "name": "example_name"
                }
			],
			"visitorPrioritization":[
				{
                    "id": 123456,
                    "name": "0000000000000000000000_EXAMPLE"
                }
			],
			"audienceSegmentation":[
				{
                    "id": 123456,
                    "name": "ExampleTest"
                }
			],
			"apiPrioritization":[
				{
                    "id": 123456,
                    "name": "APIExampleTest"
                }
			],
			"edgeRedirector":[
				{
                    "id": 123456,
                    "name": "00000000_EXAMPLE"
                }
			]
		},
    	"cloudletSharedPolicyVirtualWaitingRoom": {
			"123456": {
                "cloudletType": "EXAMPLETYPE",
                "id": 123456,
                "name": "example_name",
                "policyType": "SHARED"
            }
		}
    }
}`,
			expectedResponse: &ListAvailableIncludesResponse{
				AvailableIncludes: []ExternalIncludeData{

					{
						IncludeID:   "test_include_id1",
						IncludeName: "test_include1",
						IncludeType: IncludeTypeMicroServices,
						FileName:    "test_include1.xml",
						ProductName: "Example_Name",
						RuleFormat:  "v2020-11-02",
					},
					{
						IncludeID:   "test_include_id2",
						IncludeName: "test_include2",
						IncludeType: IncludeTypeMicroServices,
						FileName:    "test_include2.xml",
						ProductName: "Example_Name",
						RuleFormat:  "v2020-11-02",
					},
				},
			},
		},
		"200 OK - available includes given only GroupID": {
			params: ListAvailableIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
				GroupID:         "grp_2",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/1/external-resources?groupId=grp_2",
			responseBody: `
{
    "externalResources": {
        "include": {
            "test_include_id1": {
                "id": "test_include_id1",
                "name": "test_include1",
                "includeType": "MICROSERVICES",
                "fileName": "test_include1.xml",
                "productName": "Example_Name",
                "ruleFormat": "v2020-11-02"
            },
			"test_include_id2": {
                "id": "test_include_id2",
                "name": "test_include2",
                "includeType": "MICROSERVICES",
                "fileName": "test_include2.xml",
                "productName": "Example_Name",
                "ruleFormat": "v2020-11-02"
            }
		},
		"cloudletSharedPolicy": {
			"123456": {
                "cloudletType": "TESTCLOUDLETTYPE",
                "id": 123456,
                "name": "TestName123456",
                "policyType": "SHARED"
            }
		},
		"availableCnames": [
			{
                "id": 123456,
                "name": "www.example.com",
                "domain": "example.net",
                "serialNumber": 123,
                "slot": null,
                "status": "Created",
                "ipv6": false,
                "useCases": [],
                "cname": "www.example.example.net",
                "isSecure": false,
                "isEdgeIPBindingEnabled": null
            }
		],
		"customOverrides": {},
		"customOverrides": {},
		"blacklistedCertDomains": [
			"s3.example.com"
		],
		"availableNetStorageGroups": [
			{
                "id": 123456,
                "name": "aa-example",
                "uploadDomainName": "spm.example.example.com",
                "downloadDomainName": "spm.example.example.com",
                "cpCodeList": [
                    {
                        "cpCode": 123456,
                        "g2oToken": null
                    }
                ]
            }
		],
    	"availableCpCodes": [
			{
                "id": 123123,
                "description": "example-test-subgroup",
                "products": [
                    "EXAMPLE"
                ],
                "createdDate": 1521566901000,
                "cpCodeLimits": null,
                "name": "example-test-subgroup"
            }
		],
    	"availablePolicies": {
			"applicationLoadBalancer":[
				{
                    "id": 123456,
                    "name": "0000000000000000_EXAMPLE_clone"
                }
			],
			"firstPartyMarketingPlus":[
				{
                    "id": 123456,
                    "name": "Example_Name_123456"
                }
			],
			"firstPartyMarketing":[
				{
                    "id": 123456,
                    "name": "Example_first_party"
                }	
			],
			"forwardRewrite":[
				{
                    "id": 123456,
                    "name": "ExampleName"
                }
			],
			"continuousDeployment":[
				{
                    "id": 123456,
                    "name": "Example_Name_123456"
                }
			],
			"requestControl":[
				{
                    "id": 123456,
                    "name": "EXAMPLE_MATCH_RULE_SIZE_RC"
                }
			],
			"inputValidation":[
				{
                    "id": 123456,
                    "name": "example_name"
                }
			],
			"visitorPrioritization":[
				{
                    "id": 123456,
                    "name": "0000000000000000000000_EXAMPLE"
                }
			],
			"audienceSegmentation":[
				{
                    "id": 123456,
                    "name": "ExampleTest"
                }
			],
			"apiPrioritization":[
				{
                    "id": 123456,
                    "name": "APIExampleTest"
                }
			],
			"edgeRedirector":[
				{
                    "id": 123456,
                    "name": "00000000_EXAMPLE"
                }
			]
		},
    	"cloudletSharedPolicyVirtualWaitingRoom": {
			"123456": {
                "cloudletType": "EXAMPLETYPE",
                "id": 123456,
                "name": "example_name",
                "policyType": "SHARED"
            }
		}
    }
}`,
			expectedResponse: &ListAvailableIncludesResponse{
				AvailableIncludes: []ExternalIncludeData{

					{
						IncludeID:   "test_include_id1",
						IncludeName: "test_include1",
						IncludeType: IncludeTypeMicroServices,
						FileName:    "test_include1.xml",
						ProductName: "Example_Name",
						RuleFormat:  "v2020-11-02",
					},
					{
						IncludeID:   "test_include_id2",
						IncludeName: "test_include2",
						IncludeType: IncludeTypeMicroServices,
						FileName:    "test_include2.xml",
						ProductName: "Example_Name",
						RuleFormat:  "v2020-11-02",
					},
				},
			},
		},
		"200 OK - available includes ContractID and GroupID not provided": {
			params: ListAvailableIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/1/external-resources",
			responseBody: `
{
    "externalResources": {
        "include": {
            "test_include_id1": {
                "id": "test_include_id1",
                "name": "test_include1",
                "includeType": "MICROSERVICES",
                "fileName": "test_include1.xml",
                "productName": "Example_Name",
                "ruleFormat": "v2020-11-02"
            },
			"test_include_id2": {
                "id": "test_include_id2",
                "name": "test_include2",
                "includeType": "MICROSERVICES",
                "fileName": "test_include2.xml",
                "productName": "Example_Name",
                "ruleFormat": "v2020-11-02"
            }
		},
		"cloudletSharedPolicy": {
			"123456": {
                "cloudletType": "TESTCLOUDLETTYPE",
                "id": 123456,
                "name": "TestName123456",
                "policyType": "SHARED"
            }
		},
		"availableCnames": [
			{
                "id": 123456,
                "name": "www.example.com",
                "domain": "example.net",
                "serialNumber": 123,
                "slot": null,
                "status": "Created",
                "ipv6": false,
                "useCases": [],
                "cname": "www.example.example.net",
                "isSecure": false,
                "isEdgeIPBindingEnabled": null
            }
		],
		"customOverrides": {},
		"customOverrides": {},
		"blacklistedCertDomains": [
			"s3.example.com"
		],
		"availableNetStorageGroups": [
			{
                "id": 123456,
                "name": "aa-example",
                "uploadDomainName": "spm.example.example.com",
                "downloadDomainName": "spm.example.example.com",
                "cpCodeList": [
                    {
                        "cpCode": 123456,
                        "g2oToken": null
                    }
                ]
            }
		],
    	"availableCpCodes": [
			{
                "id": 123123,
                "description": "example-test-subgroup",
                "products": [
                    "EXAMPLE"
                ],
                "createdDate": 1521566901000,
                "cpCodeLimits": null,
                "name": "example-test-subgroup"
            }
		],
    	"availablePolicies": {
			"applicationLoadBalancer":[
				{
                    "id": 123456,
                    "name": "0000000000000000_EXAMPLE_clone"
                }
			],
			"firstPartyMarketingPlus":[
				{
                    "id": 123456,
                    "name": "Example_Name_123456"
                }
			],
			"firstPartyMarketing":[
				{
                    "id": 123456,
                    "name": "Example_first_party"
                }	
			],
			"forwardRewrite":[
				{
                    "id": 123456,
                    "name": "ExampleName"
                }
			],
			"continuousDeployment":[
				{
                    "id": 123456,
                    "name": "Example_Name_123456"
                }
			],
			"requestControl":[
				{
                    "id": 123456,
                    "name": "EXAMPLE_MATCH_RULE_SIZE_RC"
                }
			],
			"inputValidation":[
				{
                    "id": 123456,
                    "name": "example_name"
                }
			],
			"visitorPrioritization":[
				{
                    "id": 123456,
                    "name": "0000000000000000000000_EXAMPLE"
                }
			],
			"audienceSegmentation":[
				{
                    "id": 123456,
                    "name": "ExampleTest"
                }
			],
			"apiPrioritization":[
				{
                    "id": 123456,
                    "name": "APIExampleTest"
                }
			],
			"edgeRedirector":[
				{
                    "id": 123456,
                    "name": "00000000_EXAMPLE"
                }
			]
		},
    	"cloudletSharedPolicyVirtualWaitingRoom": {
			"123456": {
                "cloudletType": "EXAMPLETYPE",
                "id": 123456,
                "name": "example_name",
                "policyType": "SHARED"
            }
		}
    }
}`,
			expectedResponse: &ListAvailableIncludesResponse{
				AvailableIncludes: []ExternalIncludeData{

					{
						IncludeID:   "test_include_id1",
						IncludeName: "test_include1",
						IncludeType: IncludeTypeMicroServices,
						FileName:    "test_include1.xml",
						ProductName: "Example_Name",
						RuleFormat:  "v2020-11-02",
					},
					{
						IncludeID:   "test_include_id2",
						IncludeName: "test_include2",
						IncludeType: IncludeTypeMicroServices,
						FileName:    "test_include2.xml",
						ProductName: "Example_Name",
						RuleFormat:  "v2020-11-02",
					},
				},
			},
		},
		"500 Internal Server Error": {
			params: ListAvailableIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
				ContractID:      "ctr_1",
				GroupID:         "grp_2",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error fetching available includes",
  "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/1/external-resources?contractId=ctr_1&groupId=grp_2",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching available includes",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing property ID": {
			params: ListAvailableIncludesRequest{
				PropertyVersion: 2,
				ContractID:      "ctr_1",
				GroupID:         "grp_2",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing property version": {
			params: ListAvailableIncludesRequest{
				PropertyID: "propertyID",
				ContractID: "ctr_1",
				GroupID:    "grp_2",
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
			result, err := client.ListAvailableIncludes(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, len(test.expectedResponse.AvailableIncludes), len(result.AvailableIncludes))
			for _, element := range test.expectedResponse.AvailableIncludes {
				assert.Contains(t, result.AvailableIncludes, element)
			}
		})
	}
}

func TestPapiListReferencedIncludes(t *testing.T) {
	tests := map[string]struct {
		params           ListReferencedIncludesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListReferencedIncludesResponse
		withError        error
	}{
		"200 OK": {
			params: ListReferencedIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
				ContractID:      "ctr_1",
				GroupID:         "grp_2",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/papi/v1/properties/propertyID/versions/1/includes?contractId=ctr_1&groupId=grp_2",
			responseBody: `
{
    "includes": {
        "items": [
            {
                "accountId": "test_account",
                "contractId": "test_contract",
                "groupId": "test_group",
                "latestVersion": 1,
                "stagingVersion": 1,
                "productionVersion": null,
                "assetId": "test_asset",
                "includeId": "inc_123456",
                "includeName": "test_include",
                "includeType": "MICROSERVICES"
            }
        ]
    }
}`,
			expectedResponse: &ListReferencedIncludesResponse{
				Includes: IncludeItems{
					Items: []Include{
						{
							AccountID:         "test_account",
							AssetID:           "test_asset",
							ContractID:        "test_contract",
							GroupID:           "test_group",
							IncludeID:         "inc_123456",
							IncludeName:       "test_include",
							IncludeType:       IncludeTypeMicroServices,
							LatestVersion:     1,
							ProductionVersion: nil,
							StagingVersion:    ptr.To(1),
						},
					},
				},
			},
		},
		"500 Internal Server Error": {
			params: ListReferencedIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
				ContractID:      "ctr_1",
				GroupID:         "grp_2",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error fetching referenced includes",
  "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/1/includes?contractId=ctr_1&groupId=grp_2",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching referenced includes",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing property ID": {
			params: ListReferencedIncludesRequest{
				PropertyVersion: 1,
				ContractID:      "ctr_1",
				GroupID:         "grp_2",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing property version": {
			params: ListReferencedIncludesRequest{
				PropertyID: "propertyID",
				ContractID: "ctr_1",
				GroupID:    "grp_2",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing contractID": {
			params: ListReferencedIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
				GroupID:         "grp_2",
			},
			withError: ErrStructValidation,
		},
		"validation error - missing groupID": {
			params: ListReferencedIncludesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 1,
				ContractID:      "ctr_1",
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
			result, err := client.ListReferencedIncludes(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
