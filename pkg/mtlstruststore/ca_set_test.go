package mtlstruststore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCASet(t *testing.T) {
	tests := map[string]struct {
		request          CreateCASetRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *CreateCASetResponse
		responseHeaders  map[string]string
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			request: CreateCASetRequest{
				CASetName:   "test",
				Description: ptr.To("description"),
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets",
			responseStatus: http.StatusCreated,
			responseBody: `
				{
					"accountId": "A-CCOUNT",
					"caSetId": "199",
					"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
					"caSetName": "test",
					"caSetStatus": "NOT_DELETED",
					"createdBy": "jdoe",
					"createdDate": "2025-04-01T15:33:48.464941Z",
					"deletedBy": null,
					"deletedDate": null,
					"description": "",
					"latestVersion": null,
					"latestVersionLink": null,
					"productionVersion": null,
					"productionVersionLink": null,
					"stagingVersion": null,
					"stagingVersionLink": null,
					"versionsLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/",
					"removalDate": null
				}`,
			expectedResponse: &CreateCASetResponse{
				AccountID:             "A-CCOUNT",
				CASetID:               "199",
				CASetLink:             "/mtls-edge-truststore/v2/ca-sets/199",
				CASetName:             "test",
				CASetStatus:           "NOT_DELETED",
				CreatedBy:             "jdoe",
				CreatedDate:           test.NewTimeFromString(t, "2025-04-01T15:33:48.464941Z"),
				DeletedBy:             nil,
				DeletedDate:           nil,
				Description:           ptr.To(""),
				LatestVersion:         nil,
				LatestVersionLink:     nil,
				ProductionVersion:     nil,
				ProductionVersionLink: nil,
				StagingVersion:        nil,
				StagingVersionLink:    nil,
				VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/199/versions/",
				RemovalDate:           nil,
			},
			responseHeaders: map[string]string{
				"Location": "/mtls-edge-truststore/v2/ca-sets/199",
			},
		},
		"missing required request param - validation error": {
			request: CreateCASetRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create ca set failed: struct validation: CASetName: cannot be blank", err.Error())
			},
		},
		"name too short - validation error": {
			request: CreateCASetRequest{
				CASetName: "a",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create ca set failed: struct validation: CASetName: the length must be between 3 and 64", err.Error())
			},
		},
		"invalid name - validation error": {
			request: CreateCASetRequest{
				CASetName: "###A",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create ca set failed: struct validation: CASetName: allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-), percent (%) and period (.)", err.Error())
			},
		},
		"invalid name with ... - validation error": {
			request: CreateCASetRequest{
				CASetName: "AAA...A",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create ca set failed: struct validation: CASetName: cannot contain three consecutive periods (...)", err.Error())
			},
		},
		"500 internal server error": {
			request: CreateCASetRequest{
				CASetName:   "test",
				Description: ptr.To("description"),
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
			expectedPath: "/mtls-edge-truststore/v2/ca-sets",
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
		"409 error response - duplicate CA set": {
			request: CreateCASetRequest{
				CASetName:   "test",
				Description: ptr.To("description"),
			},
			responseStatus: http.StatusConflict,
			responseBody: `
{
  "contextInfo": {
    "accountId": "A-CCOUNT",
    "caSetName": "testCAsetType"
  },
  "errors": [
    {
      "detail": "CA set with caSetName testCAsetType cannot be created as another CA set with the same name exists in the account with accountId A-CCOUNT.",
      "pointer": "/caSetName",
      "title": "CA set with the same name exists in the account.",
      "type": "/mtls-edge-truststore/error-types/ca-set-name-is-not-unique"
    }
  ],
  "instance": "/mtls-edge-truststore/error-types/ca-set-name-is-not-unique/dadce1c2d3e34567",
  "status": 409,
  "title": "CA set with the same name exists in the account.",
  "type": "/mtls-edge-truststore/error-types/ca-set-name-is-not-unique"
}`,
			expectedPath: "/mtls-edge-truststore/v2/ca-sets",
			withError: func(t *testing.T, err error) {
				want := &Error{
					ContextInfo: map[string]interface{}{
						"accountId": "A-CCOUNT",
						"caSetName": "testCAsetType",
					},
					Errors: []ErrorItem{
						{
							Detail:  "CA set with caSetName testCAsetType cannot be created as another CA set with the same name exists in the account with accountId A-CCOUNT.",
							Pointer: "/caSetName",
							Title:   "CA set with the same name exists in the account.",
							Type:    "/mtls-edge-truststore/error-types/ca-set-name-is-not-unique",
						},
					},
					Instance: "/mtls-edge-truststore/error-types/ca-set-name-is-not-unique/dadce1c2d3e34567",
					Status:   409,
					Title:    "CA set with the same name exists in the account.",
					Type:     "/mtls-edge-truststore/error-types/ca-set-name-is-not-unique",
				}
				assert.True(t, reflect.DeepEqual(err, want), "want: %s; got: %s", want, err)
			},
		},
		"415 error response - wrong content type header": {
			request: CreateCASetRequest{
				CASetName:   "test",
				Description: ptr.To("description"),
			},
			responseStatus: http.StatusUnsupportedMediaType,
			responseBody: `
{
    "contextInfo": {
        "allowedMediaTypes": [
            "application/json"
        ],
        "unsupportedMediaType": "application/json1"
    },
    "detail": "Media type application/json1 is not supported. Supported media type(s) are [application/json].",
    "instance": "/mtls-edge-truststore/error-types/media-type-not-supported/5de24ab7be0613b4",
    "status": 415,
    "title": "Media type not supported.",
    "type": "/mtls-edge-truststore/error-types/media-type-not-supported"
}`,
			expectedPath: "/mtls-edge-truststore/v2/ca-sets",
			withError: func(t *testing.T, err error) {
				want := &Error{
					ContextInfo: map[string]interface{}{
						"allowedMediaTypes":    []interface{}{"application/json"},
						"unsupportedMediaType": "application/json1",
					},
					Detail:   "Media type application/json1 is not supported. Supported media type(s) are [application/json].",
					Instance: "/mtls-edge-truststore/error-types/media-type-not-supported/5de24ab7be0613b4",
					Status:   415,
					Title:    "Media type not supported.",
					Type:     "/mtls-edge-truststore/error-types/media-type-not-supported",
				}
				assert.True(t, reflect.DeepEqual(err, want), "want: %s; got: %s", want, err)
			},
		},
		"400 error response - wrong content type header": {
			request: CreateCASetRequest{
				CASetName:   "test",
				Description: ptr.To("description"),
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
    "contextInfo": {
        "message": "instance type (integer) does not match any allowed primitive type (allowed: [\"string\"])"
    },
    "instance": "/mtls-edge-truststore/error-types/json-schema-validation-error/12345a6c78caa9bb",
    "pointer": "/caSetName",
    "status": 400,
    "title": "Body failed JSON schema validation.",
    "type": "/mtls-edge-truststore/error-types/json-schema-validation-error"
}`,
			expectedPath: "/mtls-edge-truststore/v2/ca-sets",
			withError: func(t *testing.T, err error) {
				want := &Error{
					ContextInfo: map[string]interface{}{
						"message": "instance type (integer) does not match any allowed primitive type (allowed: [\"string\"])",
					},
					Instance: "/mtls-edge-truststore/error-types/json-schema-validation-error/12345a6c78caa9bb",
					Pointer:  "/caSetName",
					Status:   400,
					Title:    "Body failed JSON schema validation.",
					Type:     "/mtls-edge-truststore/error-types/json-schema-validation-error",
				}
				assert.True(t, reflect.DeepEqual(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(tc.responseHeaders) > 0 {
					for header, value := range tc.responseHeaders {
						w.Header().Set(header, value)
					}
				}

				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateCASet(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetCASet(t *testing.T) {
	tests := map[string]struct {
		params           GetCASetRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *GetCASetResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusOK,
			responseBody: `
				{
					"accountId": "A-CCOUNT",
					"caSetId": "199",
					"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
					"caSetName": "test",
					"caSetStatus": "NOT_DELETED",
					"createdBy": "jdoe",
					"createdDate": "2025-04-01T15:33:48.464941Z",
					"deletedBy": null,
					"deletedDate": null,
					"description": "",
					"latestVersion": null,
					"latestVersionLink": null,
					"productionVersion": null,
					"productionVersionLink": null,
					"stagingVersion": null,
					"stagingVersionLink": null,
					"versionsLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/",
					"removalDate": null
				}`,
			expectedResponse: &GetCASetResponse{
				AccountID:             "A-CCOUNT",
				CASetID:               "199",
				CASetLink:             "/mtls-edge-truststore/v2/ca-sets/199",
				CASetName:             "test",
				CASetStatus:           "NOT_DELETED",
				CreatedBy:             "jdoe",
				CreatedDate:           test.NewTimeFromString(t, "2025-04-01T15:33:48.464941Z"),
				DeletedBy:             nil,
				DeletedDate:           nil,
				Description:           ptr.To(""),
				LatestVersion:         nil,
				LatestVersionLink:     nil,
				ProductionVersion:     nil,
				ProductionVersionLink: nil,
				StagingVersion:        nil,
				StagingVersionLink:    nil,
				VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/199/versions/",
				RemovalDate:           nil,
			},
		},
		"missing required params - validation error": {
			params: GetCASetRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get ca set failed: struct validation: CASetID: cannot be blank", err.Error())
			},
		},
		"404 ca set key not found - custom error check": {
			params: GetCASetRequest{
				CASetID: "10",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/10",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"contextInfo": {
					"caSetId": "10"
				},
				"detail": "Cannot get CA set as the CA set with caSetId 10 is not found.",
				"status": 404,
				"title": "CA set is not found.",
				"type": "/mtls-edge-truststore/error-types/ca-set-not-found"
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetNotFound))
			},
		},
		"500 internal server error": {
			params: GetCASetRequest{
				CASetID: "199",
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
			expectedPath: "/mtls-edge-truststore/v2/ca-sets/199",
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
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCASet(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListCASets(t *testing.T) {
	tests := map[string]struct {
		request          ListCASetsRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *ListCASetsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request:        ListCASetsRequest{},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets",
			responseStatus: http.StatusOK,
			responseBody: `{
				"caSets": [
					{
						"accountId": "A-CCOUNT",
						"caSetId": "199",
						"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
						"caSetName": "test",
						"caSetStatus": "NOT_DELETED",
						"createdBy": "jdoe",
						"createdDate": "2025-04-01T15:33:48.464941Z",
						"deletedBy": null,
						"deletedDate": null,
						"description": "",
						"latestVersion": null,
						"latestVersionLink": null,
						"productionVersion": null,
						"productionVersionLink": null,
						"stagingVersion": null,
						"stagingVersionLink": null,
						"versionsLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/",
						"removalDate": null
					},
					{
						"accountId": "A-CCOUNT",
						"caSetId": "80431",
						"caSetLink": "/mtls-edge-truststore/v2/ca-sets/80431",
						"caSetName": "sktcm2-051623",
						"caSetStatus": "NOT_DELETED",
						"createdBy": "migration_run",
						"createdDate": "2023-10-17T23:04:52.491822Z",
						"deletedBy": null,
						"deletedDate": null,
						"description": "Imported from Techpreview TCM",
						"latestVersion": 2,
						"latestVersionLink": "/mtls-edge-truststore/v2/ca-sets/80431/versions/2",
						"productionVersion": 2,
						"productionVersionLink": "/mtls-edge-truststore/v2/ca-sets/80431/versions/2",
						"stagingVersion": 2,
						"stagingVersionLink": "/mtls-edge-truststore/v2/ca-sets/80431/versions/2",
						"versionsLink": "/mtls-edge-truststore/v2/ca-sets/80431/versions/",
						"removalDate": null
					},
					{
						"accountId": "A-CCOUNT",
						"caSetId": "75201",
						"caSetLink": "/mtls-edge-truststore/v2/ca-sets/75201",
						"caSetName": "CertSet-4-docs",
						"caSetStatus": "DELETED",
						"createdBy": "migration_run",
						"createdDate": "2023-10-17T23:04:52.884782Z",
						"deletedBy": "migration_run",
						"deletedDate": "2025-06-04T12:19:33.095023Z",
						"description": "Imported from Techpreview TCM",
						"latestVersion": 3,
						"latestVersionLink": "/mtls-edge-truststore/v2/ca-sets/75201/versions/3",
						"productionVersion": null,
						"productionVersionLink": null,
						"stagingVersion": null,
						"stagingVersionLink": null,
						"versionsLink": "/mtls-edge-truststore/v2/ca-sets/75201/versions/",
						"removalDate": "2025-07-01T00:00:00Z"
					}
				]
			}`,
			expectedResponse: &ListCASetsResponse{
				CASets: []CASetResponse{
					{
						AccountID:             "A-CCOUNT",
						CASetID:               "199",
						CASetLink:             "/mtls-edge-truststore/v2/ca-sets/199",
						CASetName:             "test",
						CASetStatus:           "NOT_DELETED",
						CreatedBy:             "jdoe",
						CreatedDate:           test.NewTimeFromString(t, "2025-04-01T15:33:48.464941Z"),
						DeletedBy:             nil,
						DeletedDate:           nil,
						Description:           ptr.To(""),
						LatestVersion:         nil,
						LatestVersionLink:     nil,
						ProductionVersion:     nil,
						ProductionVersionLink: nil,
						StagingVersion:        nil,
						StagingVersionLink:    nil,
						VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/199/versions/",
						RemovalDate:           nil,
					},
					{
						AccountID:             "A-CCOUNT",
						CASetID:               "80431",
						CASetLink:             "/mtls-edge-truststore/v2/ca-sets/80431",
						CASetName:             "sktcm2-051623",
						CASetStatus:           "NOT_DELETED",
						CreatedBy:             "migration_run",
						CreatedDate:           test.NewTimeFromString(t, "2023-10-17T23:04:52.491822Z"),
						DeletedBy:             nil,
						DeletedDate:           nil,
						Description:           ptr.To("Imported from Techpreview TCM"),
						LatestVersion:         ptr.To(int64(2)),
						LatestVersionLink:     ptr.To("/mtls-edge-truststore/v2/ca-sets/80431/versions/2"),
						ProductionVersion:     ptr.To(int64(2)),
						ProductionVersionLink: ptr.To("/mtls-edge-truststore/v2/ca-sets/80431/versions/2"),
						StagingVersion:        ptr.To(int64(2)),
						StagingVersionLink:    ptr.To("/mtls-edge-truststore/v2/ca-sets/80431/versions/2"),
						VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/80431/versions/",
						RemovalDate:           nil,
					},
					{
						AccountID:             "A-CCOUNT",
						CASetID:               "75201",
						CASetLink:             "/mtls-edge-truststore/v2/ca-sets/75201",
						CASetName:             "CertSet-4-docs",
						CASetStatus:           "DELETED",
						CreatedBy:             "migration_run",
						CreatedDate:           test.NewTimeFromString(t, "2023-10-17T23:04:52.884782Z"),
						DeletedBy:             ptr.To("migration_run"),
						DeletedDate:           ptr.To(test.NewTimeFromString(t, "2025-06-04T12:19:33.095023Z")),
						Description:           ptr.To("Imported from Techpreview TCM"),
						LatestVersion:         ptr.To(int64(3)),
						LatestVersionLink:     ptr.To("/mtls-edge-truststore/v2/ca-sets/75201/versions/3"),
						ProductionVersion:     nil,
						ProductionVersionLink: nil,
						StagingVersion:        nil,
						StagingVersionLink:    nil,
						VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/75201/versions/",
						RemovalDate:           ptr.To(test.NewTimeFromString(t, "2025-07-01T00:00:00Z")),
					},
				},
			},
		},
		"200 OK with filtering and empty response": {
			request: ListCASetsRequest{
				CASetNamePrefix: "foo",
				ActivatedOn:     "PRODUCTION",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets?activatedOn=PRODUCTION&caSetNamePrefix=foo",
			responseStatus: http.StatusOK,
			responseBody: `{
				"caSets": []
			}`,
			expectedResponse: &ListCASetsResponse{
				CASets: []CASetResponse{},
			},
		},
		"200 OK with filtering with non lower case value and an empty response": {
			request: ListCASetsRequest{
				CASetNamePrefix: "foo",
				ActivatedOn:     "PRODUCTION",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets?activatedOn=PRODUCTION&caSetNamePrefix=foo",
			responseStatus: http.StatusOK,
			responseBody: `{
				"caSets": []
			}`,
			expectedResponse: &ListCASetsResponse{
				CASets: []CASetResponse{},
			},
		},
		"200 OK on production and staging with an empty response": {
			request: ListCASetsRequest{
				ActivatedOn: "STAGING+PRODUCTION",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets?activatedOn=STAGING%2BPRODUCTION",
			responseStatus: http.StatusOK,
			responseBody: `{
				"caSets": []
			}`,
			expectedResponse: &ListCASetsResponse{
				CASets: []CASetResponse{},
			},
		},
		"200 OK with empty CASetStatuses - defaults to NOT_DELETED": {
			request: ListCASetsRequest{
				CASetStatuses: []string{},
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets",
			responseStatus: http.StatusOK,
			responseBody: `{
				"caSets": [
					{
						"accountId": "A-CCOUNT",
						"caSetId": "199",
						"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
						"caSetName": "test",
						"caSetStatus": "NOT_DELETED",
						"createdBy": "jdoe",
						"createdDate": "2025-04-01T15:33:48.464941Z",
						"deletedBy": null,
						"deletedDate": null,
						"description": "",
						"latestVersion": null,
						"latestVersionLink": null,
						"productionVersion": null,
						"productionVersionLink": null,
						"stagingVersion": null,
						"stagingVersionLink": null,
						"versionsLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/"
					}
				]
			}`,
			expectedResponse: &ListCASetsResponse{
				CASets: []CASetResponse{
					{
						AccountID:             "A-CCOUNT",
						CASetID:               "199",
						CASetLink:             "/mtls-edge-truststore/v2/ca-sets/199",
						CASetName:             "test",
						CASetStatus:           "NOT_DELETED",
						CreatedBy:             "jdoe",
						CreatedDate:           test.NewTimeFromString(t, "2025-04-01T15:33:48.464941Z"),
						DeletedBy:             nil,
						DeletedDate:           nil,
						Description:           ptr.To(""),
						LatestVersion:         nil,
						LatestVersionLink:     nil,
						ProductionVersion:     nil,
						ProductionVersionLink: nil,
						StagingVersion:        nil,
						StagingVersionLink:    nil,
						VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/199/versions/",
					},
				},
			},
		},
		"200 OK with single CASetStatuses filter": {
			request: ListCASetsRequest{
				CASetStatuses: []string{CASetStatusDeleted},
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets?caSetStatus=DELETED",
			responseStatus: http.StatusOK,
			responseBody: `{
				"caSets": [
					{
						"accountId": "A-CCOUNT",
						"caSetId": "75201",
						"caSetLink": "/mtls-edge-truststore/v2/ca-sets/75201",
						"caSetName": "test",
						"caSetStatus": "DELETED",
						"createdBy": "migration_run",
						"createdDate": "2023-10-17T23:04:52.884782Z",
						"deletedBy": "migration_run",
						"deletedDate": "2025-06-04T12:19:33.095023Z",
						"description": "Imported from Techpreview TCM",
						"latestVersion": 3,
						"latestVersionLink": "/mtls-edge-truststore/v2/ca-sets/75201/versions/3",
						"productionVersion": null,
						"productionVersionLink": null,
						"stagingVersion": null,
						"stagingVersionLink": null,
						"versionsLink": "/mtls-edge-truststore/v2/ca-sets/75201/versions/"
					}
				]
			}`,
			expectedResponse: &ListCASetsResponse{
				CASets: []CASetResponse{
					{
						AccountID:             "A-CCOUNT",
						CASetID:               "75201",
						CASetLink:             "/mtls-edge-truststore/v2/ca-sets/75201",
						CASetName:             "test",
						CASetStatus:           "DELETED",
						CreatedBy:             "migration_run",
						CreatedDate:           test.NewTimeFromString(t, "2023-10-17T23:04:52.884782Z"),
						DeletedBy:             ptr.To("migration_run"),
						DeletedDate:           ptr.To(test.NewTimeFromString(t, "2025-06-04T12:19:33.095023Z")),
						Description:           ptr.To("Imported from Techpreview TCM"),
						LatestVersion:         ptr.To(int64(3)),
						LatestVersionLink:     ptr.To("/mtls-edge-truststore/v2/ca-sets/75201/versions/3"),
						ProductionVersion:     nil,
						ProductionVersionLink: nil,
						StagingVersion:        nil,
						StagingVersionLink:    nil,
						VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/75201/versions/",
					},
				},
			},
		},
		"200 OK with multiple CASetStatuses filter": {
			request: ListCASetsRequest{
				CASetStatuses: []string{CASetStatusDeleted, CASetStatusNotDeleted},
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets?caSetStatus=DELETED%2CNOT_DELETED",
			responseStatus: http.StatusOK,
			responseBody: `{
				"caSets": [
					{
						"accountId": "A-CCOUNT",
						"caSetId": "199",
						"caSetLink": "/mtls-edge-truststore/v2/ca-sets/199",
						"caSetName": "test",
						"caSetStatus": "NOT_DELETED",
						"createdBy": "jdoe",
						"createdDate": "2025-04-01T15:33:48.464941Z",
						"deletedBy": null,
						"deletedDate": null,
						"description": "",
						"latestVersion": null,
						"latestVersionLink": null,
						"productionVersion": null,
						"productionVersionLink": null,
						"stagingVersion": null,
						"stagingVersionLink": null,
						"versionsLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/"
					},
					{
						"accountId": "A-CCOUNT",
						"caSetId": "75201",
						"caSetLink": "/mtls-edge-truststore/v2/ca-sets/75201",
						"caSetName": "test2",
						"caSetStatus": "DELETED",
						"createdBy": "migration_run",
						"createdDate": "2023-10-17T23:04:52.884782Z",
						"deletedBy": "migration_run",
						"deletedDate": "2025-06-04T12:19:33.095023Z",
						"description": "Imported from Techpreview TCM",
						"latestVersion": 3,
						"latestVersionLink": "/mtls-edge-truststore/v2/ca-sets/75201/versions/3",
						"productionVersion": null,
						"productionVersionLink": null,
						"stagingVersion": null,
						"stagingVersionLink": null,
						"versionsLink": "/mtls-edge-truststore/v2/ca-sets/75201/versions/"
					}
				]
			}`,
			expectedResponse: &ListCASetsResponse{
				CASets: []CASetResponse{
					{
						AccountID:             "A-CCOUNT",
						CASetID:               "199",
						CASetLink:             "/mtls-edge-truststore/v2/ca-sets/199",
						CASetName:             "test",
						CASetStatus:           "NOT_DELETED",
						CreatedBy:             "jdoe",
						CreatedDate:           test.NewTimeFromString(t, "2025-04-01T15:33:48.464941Z"),
						DeletedBy:             nil,
						DeletedDate:           nil,
						Description:           ptr.To(""),
						LatestVersion:         nil,
						LatestVersionLink:     nil,
						ProductionVersion:     nil,
						ProductionVersionLink: nil,
						StagingVersion:        nil,
						StagingVersionLink:    nil,
						VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/199/versions/",
					},
					{
						AccountID:             "A-CCOUNT",
						CASetID:               "75201",
						CASetLink:             "/mtls-edge-truststore/v2/ca-sets/75201",
						CASetName:             "test2",
						CASetStatus:           "DELETED",
						CreatedBy:             "migration_run",
						CreatedDate:           test.NewTimeFromString(t, "2023-10-17T23:04:52.884782Z"),
						DeletedBy:             ptr.To("migration_run"),
						DeletedDate:           ptr.To(test.NewTimeFromString(t, "2025-06-04T12:19:33.095023Z")),
						Description:           ptr.To("Imported from Techpreview TCM"),
						LatestVersion:         ptr.To(int64(3)),
						LatestVersionLink:     ptr.To("/mtls-edge-truststore/v2/ca-sets/75201/versions/3"),
						ProductionVersion:     nil,
						ProductionVersionLink: nil,
						StagingVersion:        nil,
						StagingVersionLink:    nil,
						VersionsLink:          "/mtls-edge-truststore/v2/ca-sets/75201/versions/",
					},
				},
			},
		},
		"invalid CASetStatuses - validation error": {
			request: ListCASetsRequest{
				CASetStatuses: []string{"INVALID"},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca sets failed: struct validation: CASetStatuses: list element 'INVALID' is invalid. "+
					"Each element must be one of: 'NOT_DELETED', 'DELETED', 'DELETING'", err.Error())
			},
		},
		"invalid network - validation error": {
			request: ListCASetsRequest{
				ActivatedOn: "PROD",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca sets failed: struct validation: ActivatedOn: value 'PROD' is invalid. "+
					"Must be one of: 'INACTIVE', 'STAGING', 'PRODUCTION', 'STAGING+PRODUCTION', 'PRODUCTION+STAGING', 'STAGING,PRODUCTION' or 'PRODUCTION,STAGING'.", err.Error())
			},
		},
		"name prefix too long - validation error": {
			request: ListCASetsRequest{
				CASetNamePrefix: strings.Repeat("A", 65),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca sets failed: struct validation: CASetNamePrefix: the length must be no more than 64", err.Error())
			},
		},
		"invalid name prefix - validation error": {
			request: ListCASetsRequest{
				CASetNamePrefix: "###A",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca sets failed: struct validation: CASetNamePrefix: allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-), percent (%) and period (.)", err.Error())
			},
		},
		"invalid name prefix with ... - validation error": {
			request: ListCASetsRequest{
				CASetNamePrefix: "AAA...A",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca sets failed: struct validation: CASetNamePrefix: cannot contain three consecutive periods (...)", err.Error())
			},
		},
		"500 internal server error": {
			request:        ListCASetsRequest{},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets",
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
			client := mockAPIClient(t, mockServer)
			result, err := client.ListCASets(context.Background(), tc.request)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestDeleteCASet(t *testing.T) {
	tests := map[string]struct {
		params         DeleteCASetRequest
		expectedPath   string
		responseStatus int
		responseBody   string
		withError      func(*testing.T, error)
	}{
		"202 Accepted": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusAccepted,
		},
		"missing required params - validation error": {
			params: DeleteCASetRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete ca set failed: struct validation: CASetID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
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
		"409 error response - delete CA set - Defense Edge": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusConflict,
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
    "caSetId" : "1",
    "caSetName" : "caSetName-30bb6f49",
    "network" : "PRODUCTION"
  },
  "detail" : "CA set cannot be deleted as CA set with caSetId 1 links to several Property Manager hostnames. You need to unlink the CA set from the hostnames to proceed. See accompanying response data for hostname details.",
  "status" : 409,
  "instance": "/mtls-edge-truststore/error-types/ca-set-bound-to-hostname/4e0069deb5f40f63",
  "title" : "CA set is linked to hostnames.",
  "type" : "/mtls-edge-truststore/error-types/ca-set-bound-to-hostname"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetBoundToHostname), "want: %s; got: %s", ErrCASetBoundToHostname, err)
			},
		},
		"404 error response - delete CA set - not found": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
  "contextInfo": {
    "caSetId": "199"
  },
  "detail": "Cannot get CA set as the CA set with caSetId 199 is not found.",
  "instance": "/mtls-edge-truststore/error-types/ca-set-not-found/982c5bdc0abffe8d",
  "status": 404,
  "title": "CA set is not found.",
  "type": "/mtls-edge-truststore/error-types/ca-set-not-found"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrDeleteCASetNotFound), "want: %s; got: %s", ErrGetCASetNotFound, err)
			},
		},
		"409 error response - delete CA set - in progress version activations": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusConflict,
			responseBody: `
{
    "contextInfo": {
        "caSetId": "199",
        "caSetName": "v2-api-create-ca-set-2-2",
        "productionLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/1/activations/12345",
        "productionStatus": "IN_PROGRESS"
    },
    "detail": "CA set with caSetId 199 cannot be deleted as it has activations in progress on one or more networks.",
    "instance": "/mtls-edge-truststore/error-types/ca-set-cannot-be-deleted-in-progress-version-activations/123dc1d57813ec84",
    "status": 409,
    "title": "CA set cannot be deleted due to in progress activations.",
    "type": "/mtls-edge-truststore/error-types/ca-set-cannot-be-deleted-in-progress-version-activations"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrDeleteActivationDeactivationInProgress), "want: %s; got: %s", ErrDeleteActivationDeactivationInProgress, err)
			},
		},
		"409 error response - delete CA set - in progress version activations staging": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusConflict,
			responseBody: `
{
    "contextInfo": {
        "caSetId": "199",
        "caSetName": "clone-caset-api-test",
        "stagingLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/4/activations/12345",
        "stagingStatus": "IN_PROGRESS"
    },
    "detail": "CA set with caSetId 199 cannot be deleted as it has activations in progress on one or more networks.",
    "instance": "/mtls-edge-truststore/error-types/ca-set-cannot-be-deleted-in-progress-version-activations/0b1a57113cf28f42",
    "status": 409,
    "title": "CA set cannot be deleted due to in progress activations.",
    "type": "/mtls-edge-truststore/error-types/ca-set-cannot-be-deleted-in-progress-version-activations"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrDeleteActivationDeactivationInProgress), "want: %s; got: %s", ErrDeleteActivationDeactivationInProgress, err)
			},
		},
		"409 error response - delete CA set - in progress version activations both networks": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusConflict,
			responseBody: `
{
    "contextInfo": {
        "caSetId": "199",
        "caSetName": "sup-m2-bugjam6",
        "productionLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/2/activations/95458",
        "productionStatus": "IN_PROGRESS",
        "stagingLink": "/mtls-edge-truststore/v2/ca-sets/199/versions/2/activations/95458", 
        "stagingStatus": "IN_PROGRESS"     
},
    "detail": "CA set with caSetId 199 cannot be deleted as it has activations in progress on one or more networks.",
    "instance": "/mtls-edge-truststore/error-types/ca-set-cannot-be-deleted-in-progress-version-activations/3f21d188e84a5566",
    "status": 409,
    "title": "CA set cannot be deleted due to in progress activations.",
    "type": "/mtls-edge-truststore/error-types/ca-set-cannot-be-deleted-in-progress-version-activations"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrDeleteActivationDeactivationInProgress), "want: %s; got: %s", ErrDeleteActivationDeactivationInProgress, err)
			},
		},
		"409 error response - CA set deletion is in progress": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusConflict,
			responseBody: `
{
    "contextInfo": {
        "caSetId": "199",
        "caSetName": "akamai-test",
        "deletionLink": "/mtls-edge-truststore/v2/ca-sets/199/status/delete",
        "productionStatus": "IN_PROGRESS",
        "stagingStatus": "IN_PROGRESS"
    },
    "detail": "Cannot delete CA set with caSetId 199 as the CA set is being deleted on one or more networks.",
    "instance": "/mtls-edge-truststore/error-types/delete-ca-set-request-in-progress/96818406fe73e90f",
    "status": 409,
    "title": "DELETE request is in progress for the CA set on the network.",
    "type": "/mtls-edge-truststore/error-types/delete-ca-set-request-in-progress"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetDeleteRequestInProgress), "want: %s; got: %s", ErrCASetDeleteRequestInProgress, err)
			},
		},
		"409 error response - CA set deletion in progress on one network and completed on other network": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusConflict,
			responseBody: `
{
    "contextInfo": {
        "caSetId": "199",
        "caSetName": "akamai-test",
        "deletionLink": "/mtls-edge-truststore/v2/ca-sets/199/status/delete",
        "productionStatus": "IN_PROGRESS",
        "stagingStatus": "COMPLETE"
    },
    "detail": "Cannot delete CA set with caSetId 133254 as the CA set is being deleted on one or more networks.",
    "instance": "/mtls-edge-truststore/error-types/delete-ca-set-request-in-progress/96818406fe73e90f",
    "status": 409,
    "title": "DELETE request is in progress for the CA set on the network.",
    "type": "/mtls-edge-truststore/error-types/delete-ca-set-request-in-progress"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetDeleteRequestInProgress), "want: %s; got: %s", ErrCASetDeleteRequestInProgress, err)
			},
		},
		"400 error response - unknown query parameters passed": {
			params: DeleteCASetRequest{
				CASetID: "199",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/199",
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
    "contextInfo": {
        "allowedQueryParameters": [
            "accountSwitchKey"
        ],
        "parameters": [
            "test"
        ]
    },
    "detail": "The query parameter 'test' is not allowed for this endpoint. The following query parameter is allowed: 'accountSwitchKey'.",
    "instance": "/mtls-edge-truststore/error-types/unknown-query-parameters/63a890d25d672e99",
    "status": 400,
    "title": "Unknown query parameters provided.",
    "type": "/mtls-edge-truststore/error-types/unknown-query-parameters"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrUnknownQueryParameters), "want: %s; got: %s", ErrUnknownQueryParameters, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeleteCASet(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestListCASetAssociations(t *testing.T) {
	tests := map[string]struct {
		params           ListCASetAssociationsRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *ListCASetAssociationsResponse
		withError        func(*testing.T, error)
	}{
		"200 - No associations": {
			params: ListCASetAssociationsRequest{
				CASetID: "1",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/associations",
			responseStatus: http.StatusOK,
			responseBody: `{
  "associations": {
    "properties": [],
    "enrollments": null
  }
}`,
			expectedResponse: &ListCASetAssociationsResponse{Associations: Associations{Properties: make([]AssociationProperty, 0)}},
		},
		"200 - Navigable property with association type value set to properties": {
			params: ListCASetAssociationsRequest{
				CASetID:         "1",
				AssociationType: AssociationTypeProperties,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/associations?associationType=properties",
			responseStatus: http.StatusOK,
			responseBody: `{
  "associations": {
    "properties": [
      {
        "propertyId": "123",
        "propertyName": "test-prp-name",
		"propertyLink": "/papi/v1/properties/123",
        "assetId": 123456,
        "groupId": 345,
        "hostnames": [
          {
            "hostName": "example.com",
            "network": "STAGING",
            "status": "ATTACHED"
          }
        ]
      }
    ],
    "enrollments": null
  }
}`,
			expectedResponse: &ListCASetAssociationsResponse{Associations: Associations{Properties: []AssociationProperty{
				{
					PropertyID:   "123",
					PropertyName: ptr.To("test-prp-name"),
					PropertyLink: "/papi/v1/properties/123",
					AssetID:      ptr.To(int64(123456)),
					GroupID:      ptr.To(int64(345)),
					Hostnames: []AssociationHostname{
						{
							Hostname: "example.com",
							Network:  "STAGING",
							Status:   "ATTACHED",
						},
					},
				},
			}}},
		},
		"200 - Non-navigable property": {
			params: ListCASetAssociationsRequest{
				CASetID: "1",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/associations",
			responseStatus: http.StatusOK,
			responseBody: `{
  "associations": {
    "properties": [
      {
        "propertyId": "123",
		"propertyLink": "/papi/v1/properties/123",
        "hostnames": [
          {
            "hostName": "example.com",
            "network": "STAGING",
            "status": "ATTACHED"
          }
        ]
      }
    ],
    "enrollments": null
  }
}`,
			expectedResponse: &ListCASetAssociationsResponse{Associations: Associations{Properties: []AssociationProperty{
				{
					PropertyID:   "123",
					PropertyLink: "/papi/v1/properties/123",
					Hostnames: []AssociationHostname{
						{
							Hostname: "example.com",
							Network:  "STAGING",
							Status:   "ATTACHED",
						},
					},
				},
			}}},
		},
		"200 - association type with value set to enrollment": {
			params: ListCASetAssociationsRequest{
				CASetID:         "1",
				AssociationType: AssociationTypeEnrollments,
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/associations?associationType=enrollments",
			responseStatus: http.StatusOK,
			responseBody: `{
  "associations": {
    "properties": null,
    "enrollments": [
      {
        "enrollmentId": 123456,
        "enrollmentLink": "/cps/v2/enrollments/123456",
        "stagingSlots": [
          78956
        ],
        "productionSlots": [
          78956,
          56478
        ],
        "cn": "example1.com"
      },
      {
        "enrollmentId": 123457,
        "enrollmentLink": "/cps/v2/enrollments/123457",
        "stagingSlots": [
          23456
        ],
        "productionSlots": [
          23456
        ],
        "cn": "example2.com"
      }
    ]
  }
}
`,
			expectedResponse: &ListCASetAssociationsResponse{Associations: Associations{Enrollments: []AssociationEnrollment{
				{
					EnrollmentID:    123456,
					EnrollmentLink:  "/cps/v2/enrollments/123456",
					StagingSlots:    []int64{78956},
					ProductionSlots: []int64{78956, 56478},
					CN:              "example1.com",
				},
				{
					EnrollmentID:    123457,
					EnrollmentLink:  "/cps/v2/enrollments/123457",
					StagingSlots:    []int64{23456},
					ProductionSlots: []int64{23456},
					CN:              "example2.com",
				},
			}}},
		},
		"missing required params - validation error": {
			params: ListCASetAssociationsRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca sets associations failed: struct validation: CASetID: cannot be blank", err.Error())
			},
		},
		"validation error - invalid association type": {
			params: ListCASetAssociationsRequest{
				CASetID:         "1",
				AssociationType: "invalid",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca sets associations failed: struct validation: AssociationType: value 'invalid' is invalid. "+
					"Must be one of: 'enrollments' or 'properties'.", err.Error())
			},
		},
		"404 ca set not found": {
			params: ListCASetAssociationsRequest{
				CASetID: "2",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/2/associations",
			responseStatus: http.StatusNotFound,
			responseBody: `{
  "contextInfo" : {
    "caSetId" : "2"
  },
  "detail" : "Cannot get CA set associations as the CA set with caSetId 2 is not found.",
  "status" : 404,
  "title" : "CA set is not found.",
  "type" : "/mtls-edge-truststore/error-types/ca-set-not-found"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetNotFound), "want: %s; got: %s", ErrGetCASetNotFound, err)
			},
		},
		"issue fetching": {
			params: ListCASetAssociationsRequest{
				CASetID: "1",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/associations",
			responseStatus: http.StatusGatewayTimeout,
			responseBody: `{
  "contextInfo" : {
    "caSetId" : 1,
    "caSetName" : "caSetName-1"
  },
  "detail" : "There are issues fetching the information about associations at this time. The request timed out. Try again later.",
  "status" : 504,
  "title" : "Couldn't fetch associations details. Timeout occurred.",
  "type" : "/mtls-edge-truststore/error-types/cannot-get-ca-set-associations-timeout"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrFetchAssociationsTimeout), "want: %s; got: %s", ErrFetchAssociationsTimeout, err)
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
			client := mockAPIClient(t, mockServer)
			result, err := client.ListCASetAssociations(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestCloneCASet(t *testing.T) {
	tests := map[string]struct {
		params              CloneCASetRequest
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedResponse    *CloneCASetResponse
		withError           func(*testing.T, error)
	}{
		"200 - No version provided": {
			params: CloneCASetRequest{
				CloneFromSetID: "1",
				NewCASetName:   "new-set",
				NewDescription: ptr.To("New CA Set"),
			},
			expectedPath: "/mtls-edge-truststore/v2/ca-sets/1/clone",
			expectedRequestBody: `{
  "caSetName":"new-set",
  "description":"New CA Set"
}`,
			responseStatus: http.StatusCreated,
			responseBody: `{
    "accountId": "1-ACC",
    "caSetId": "2",
    "caSetLink": "/mtls-edge-truststore/v2/ca-sets/2",
    "caSetName": "new-set",
    "caSetStatus": "NOT_DELETED",
    "createdBy": "user1",
    "createdDate": "2025-04-10T07:03:32.987904Z",
    "deletedBy": null,
    "deletedDate": null,
    "description": "New CA Set",
    "latestVersion": 1,
    "latestVersionLink": "/mtls-edge-truststore/v2/ca-sets/2/versions/1",
    "productionVersion": null,
    "productionVersionLink": null,
    "stagingVersion": null,
    "stagingVersionLink": null,
    "versionsLink": "/mtls-edge-truststore/v2/ca-sets/2/versions/",
	"removalDate": null
}`,
			expectedResponse: &CloneCASetResponse{
				AccountID:         "1-ACC",
				CASetID:           "2",
				CASetLink:         "/mtls-edge-truststore/v2/ca-sets/2",
				CASetName:         "new-set",
				CASetStatus:       "NOT_DELETED",
				CreatedBy:         "user1",
				CreatedDate:       test.NewTimeFromString(t, "2025-04-10T07:03:32.987904Z"),
				Description:       ptr.To("New CA Set"),
				LatestVersion:     ptr.To(int64(1)),
				LatestVersionLink: ptr.To("/mtls-edge-truststore/v2/ca-sets/2/versions/1"),
				VersionsLink:      "/mtls-edge-truststore/v2/ca-sets/2/versions/",
				RemovalDate:       nil,
			},
		},
		"200 - Version provided, no description": {
			params: CloneCASetRequest{
				CloneFromSetID:   "1",
				CloneFromVersion: 2,
				NewCASetName:     "new-set",
				NewDescription:   nil,
			},
			expectedPath: "/mtls-edge-truststore/v2/ca-sets/1/clone?cloneFromVersion=2",
			expectedRequestBody: `{
  "caSetName":"new-set",
  "description":null
}`,
			responseStatus: http.StatusCreated,
			responseBody: `{
    "accountId": "1-ACC",
    "caSetId": "2",
    "caSetLink": "/mtls-edge-truststore/v2/ca-sets/2",
    "caSetName": "new-set",
    "caSetStatus": "NOT_DELETED",
    "createdBy": "user1",
    "createdDate": "2025-04-10T07:03:32.987904Z",
    "deletedBy": null,
    "deletedDate": null,
    "description": null,
    "latestVersion": 1,
    "latestVersionLink": "/mtls-edge-truststore/v2/ca-sets/2/versions/1",
    "productionVersion": null,
    "productionVersionLink": null,
    "stagingVersion": null,
    "stagingVersionLink": null,
    "versionsLink": "/mtls-edge-truststore/v2/ca-sets/2/versions/",
	"removalDate": null
}`,
			expectedResponse: &CloneCASetResponse{
				AccountID:         "1-ACC",
				CASetID:           "2",
				CASetLink:         "/mtls-edge-truststore/v2/ca-sets/2",
				CASetName:         "new-set",
				CASetStatus:       "NOT_DELETED",
				CreatedBy:         "user1",
				CreatedDate:       test.NewTimeFromString(t, "2025-04-10T07:03:32.987904Z"),
				Description:       nil,
				LatestVersion:     ptr.To(int64(1)),
				LatestVersionLink: ptr.To("/mtls-edge-truststore/v2/ca-sets/2/versions/1"),
				VersionsLink:      "/mtls-edge-truststore/v2/ca-sets/2/versions/",
				RemovalDate:       nil,
			},
		},
		"missing required params - validation error": {
			params: CloneCASetRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "clone ca set failed: struct validation: CloneFromSetID: cannot be blank\nNewCASetName: cannot be blank", err.Error())
			},
		},
		"too short required params - validation error": {
			params: CloneCASetRequest{
				CloneFromSetID: "1",
				NewCASetName:   "a",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "clone ca set failed: struct validation: NewCASetName: the length must be between 3 and 64", err.Error())
			},
		},
		"incorrect characters in required parameters - validation error": {
			params: CloneCASetRequest{
				CloneFromSetID: "1",
				NewCASetName:   "#edgegrid",
				NewDescription: ptr.To("New CA Set"),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "clone ca set failed: struct validation: NewCASetName: allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-), percent (%) and period (.)", err.Error())
			},
		},
		"empty description - validation error": {
			params: CloneCASetRequest{
				CloneFromSetID: "1",
				NewCASetName:   "edgegrid",
				NewDescription: ptr.To(""),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "clone ca set failed: struct validation: NewDescription: cannot be blank", err.Error())
			},
		},
		"special case in required parameters - validation error": {
			params: CloneCASetRequest{
				CloneFromSetID: "1",
				NewCASetName:   "abc...cba",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "clone ca set failed: struct validation: NewCASetName: cannot contain three consecutive periods (...)", err.Error())
			},
		},
		"400 ca set does not have any version": {
			params: CloneCASetRequest{
				CloneFromSetID: "1",
				NewCASetName:   "new-set",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/clone",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
    "contextInfo": {
        "caSetId": "1",
        "caSetName": "test1"
    },
    "detail": "CA set with caSetId 1 does not contain any versions. At least one version must be present to clone the CA set.",
    "status": 400,
    "title": "CA set does not contain any versions.",
    "type": "/mtls-edge-truststore/error-types/missing-caset-version"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrMissingCASetVersion), "want: %s; got: %s", ErrMissingCASetVersion, err)
			},
		},
		"404 ca set not found": {
			params: CloneCASetRequest{
				CloneFromSetID: "2",
				NewCASetName:   "new-set",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/2/clone",
			responseStatus: http.StatusNotFound,
			responseBody: `{
		 "contextInfo" : {
		   "caSetId" : 2
		 },
		 "detail" : "Cannot clone CA set as the CA set with caSetId 2 is not found.",
		 "status" : 404,
		 "title" : "CA set is not found.",
		 "type" : "/mtls-edge-truststore/error-types/ca-set-not-found"
		}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetNotFound), "want: %s; got: %s", ErrGetCASetNotFound, err)
			},
		},
		"404 ca set version not found": {
			params: CloneCASetRequest{
				CloneFromSetID:   "1",
				CloneFromVersion: 2,
				NewCASetName:     "new-set",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/clone?cloneFromVersion=2",
			responseStatus: http.StatusNotFound,
			responseBody: `{
  "type": "/mtls-edge-truststore/error-types/ca-set-version-not-found",
  "title": "CA set version not found",
  "status": 404,
  "detail": "Cannot clone CA set as the CA set version with version 2 is not found in the CA set under caSetName test1.",
  "contextInfo": {
    "caSetName": "test1",
    "caSetId": "1",
    "version": 2
  }
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetVersionNotFound), "want: %s; got: %s", ErrGetCASetVersionNotFound, err)
			},
		},
		"409 duplicate": {
			params: CloneCASetRequest{
				CloneFromSetID: "1",
				NewCASetName:   "new-set",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/clone",
			responseStatus: http.StatusConflict,
			responseBody: `{
  "type": "/mtls-edge-truststore/error-types/ca-set-name-is-not-unique",
  "title": "CA set name already exists",
  "status": 409,
  "detail": "CA set with caSetName new-set cannot be created as another CA set with the same name exists in the account with accountId 1-ACC.",
  "contextInfo": {
    "caSetName": "new-set",
    "accountId": "1-ACC"
  }
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetNameNotUnique), "want: %s; got: %s", ErrCASetNameNotUnique, err)
			},
		},
		"422 reached limit": {
			params: CloneCASetRequest{
				CloneFromSetID: "1",
				NewCASetName:   "new-set",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/clone",
			responseStatus: http.StatusUnprocessableEntity,
			responseBody: `{
  "type": "/mtls-edge-truststore/error-types/ca-set-limit-reached",
  "title": "Cannot create a new CA set. Maximum allowed CA set limit has been reached.",
  "status": 422,
  "detail": "Cannot create CA set as you have already reached or exceeded the maximum allowed CA set limit of 2 for your account. Please delete any unused or unwanted CA sets before you attempt to create a new CA set.",
  "contextInfo": {
    "accountId": "1-ACC",
    "currentSetCount": 2,
    "limit": 2
  }
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrCASetLimitReached), "want: %s; got: %s", ErrCASetLimitReached, err)
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
			client := mockAPIClient(t, mockServer)
			result, err := client.CloneCASet(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetCASetDeletionStatus(t *testing.T) {
	tests := map[string]struct {
		params           GetCASetDeletionStatusRequest
		expectedPath     string
		responseStatus   int
		responseHeaders  map[string]string
		responseBody     string
		expectedResponse *GetCASetDeletionStatusResponse
		withError        func(*testing.T, error)
	}{
		"202 - in progress": {
			params: GetCASetDeletionStatusRequest{
				CASetID: "1",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
			responseStatus: http.StatusAccepted,
			responseHeaders: map[string]string{
				"Retry-After": "Tue, 15 Apr 2025 12:15:02 GMT",
			},
			responseBody: `{
    "caSetId": "1",
    "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1",
    "caSetName": "test1",
    "deletions": [
        {
            "failureReason": null,
            "network": "PRODUCTION",
            "percentComplete": 0,
            "status": "IN_PROGRESS"
        },
        {
            "failureReason": null,
            "network": "STAGING",
            "percentComplete": 0,
            "status": "IN_PROGRESS"
        }
    ],
    "endTime": null,
    "estimatedEndTime": "2025-04-15T12:25:02.183294Z",
    "failureReason": null,
    "resourceMethod": "delete",
    "startTime": "2025-04-15T12:10:02.039140Z",
    "status": "IN_PROGRESS",
    "statusLink": "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
	"removalDate": "2025-07-01T00:00:00Z"
}`,
			expectedResponse: &GetCASetDeletionStatusResponse{
				CASetID:   "1",
				CASetLink: "/mtls-edge-truststore/v2/ca-sets/1",
				CASetName: "test1",
				Deletions: []CASetNetworkDeleteStatus{
					{
						Network:         "PRODUCTION",
						PercentComplete: 0,
						Status:          "IN_PROGRESS",
					},
					{
						Network:         "STAGING",
						PercentComplete: 0,
						Status:          "IN_PROGRESS",
					},
				},
				EstimatedEndTime: ptr.To(test.NewTimeFromString(t, "2025-04-15T12:25:02.183294Z")),
				ResourceMethod:   ptr.To("delete"),
				StartTime:        test.NewTimeFromString(t, "2025-04-15T12:10:02.039140Z"),
				Status:           "IN_PROGRESS",
				StatusLink:       "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
				RetryAfter:       test.NewGMTTimeFromString(t, "Tue, 15 Apr 2025 12:15:02 GMT"),
				RemovalDate:      ptr.To(test.NewTimeFromString(t, "2025-07-01T00:00:00Z")),
			},
		},
		"200 - completed": {
			params: GetCASetDeletionStatusRequest{
				CASetID: "1",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
			responseStatus: http.StatusOK,
			responseBody: `{
    "caSetId": "1",
    "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1",
    "caSetName": "test1",
    "deletions": [
        {
            "failureReason": null,
            "network": "PRODUCTION",
            "percentComplete": 100,
            "status": "COMPLETE"
        },
        {
            "failureReason": null,
            "network": "STAGING",
            "percentComplete": 100,
            "status": "COMPLETE"
        }
    ],
    "endTime": "2025-04-15T12:13:30.082193Z",
    "estimatedEndTime": null,
    "failureReason": null,
    "resourceMethod": "delete",
    "startTime": "2025-04-15T12:10:02.039140Z",
    "status": "COMPLETE",
    "statusLink": "/mtls-edge-truststore/v2/ca-sets/1/status/delete"
}`,
			expectedResponse: &GetCASetDeletionStatusResponse{
				CASetID:   "1",
				CASetLink: "/mtls-edge-truststore/v2/ca-sets/1",
				CASetName: "test1",
				Deletions: []CASetNetworkDeleteStatus{
					{
						Network:         "PRODUCTION",
						PercentComplete: 100,
						Status:          "COMPLETE",
					},
					{
						Network:         "STAGING",
						PercentComplete: 100,
						Status:          "COMPLETE",
					},
				},
				EndTime:        ptr.To(test.NewTimeFromString(t, "2025-04-15T12:13:30.082193Z")),
				ResourceMethod: ptr.To("delete"),
				StartTime:      test.NewTimeFromString(t, "2025-04-15T12:10:02.039140Z"),
				Status:         "COMPLETE",
				StatusLink:     "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
			},
		},
		"207 - partial": {
			params: GetCASetDeletionStatusRequest{
				CASetID: "1",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
			responseStatus: http.StatusAccepted,
			responseHeaders: map[string]string{
				"Retry-After": "Tue, 15 Apr 2025 12:15:02 GMT",
			},
			responseBody: `{
    "caSetId": "1",
    "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1",
    "caSetName": "test1",
    "deletions": [
        {
            "failureReason": "Reason for failure",
            "network": "STAGING",
            "percentComplete": 100,
            "status": "FAILED"
        },
        {
            "network": "PRODUCTION",
            "percentComplete": 100,
            "status": "COMPLETE"
        }
    ],
    "endTime": "2025-04-15T12:13:30.082193Z",
	"failureReason": "Indication of which network had a failure in deletion.",
	"resourceMethod": "delete",
    "startTime": "2025-04-15T12:10:02.039140Z",
	"status": "FAILED",
    "statusLink": "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
	"removalDate": null
}`,
			expectedResponse: &GetCASetDeletionStatusResponse{
				CASetID:   "1",
				CASetLink: "/mtls-edge-truststore/v2/ca-sets/1",
				CASetName: "test1",
				Deletions: []CASetNetworkDeleteStatus{
					{
						Network:         "STAGING",
						PercentComplete: 100,
						Status:          "FAILED",
						FailureReason:   ptr.To("Reason for failure"),
					},
					{
						Network:         "PRODUCTION",
						PercentComplete: 100,
						Status:          "COMPLETE",
					},
				},
				EndTime:        ptr.To(test.NewTimeFromString(t, "2025-04-15T12:13:30.082193Z")),
				ResourceMethod: ptr.To("delete"),
				StartTime:      test.NewTimeFromString(t, "2025-04-15T12:10:02.039140Z"),
				Status:         "FAILED",
				StatusLink:     "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
				FailureReason:  ptr.To("Indication of which network had a failure in deletion."),
				RetryAfter:     test.NewGMTTimeFromString(t, "Tue, 15 Apr 2025 12:15:02 GMT"),
				RemovalDate:    nil,
			},
		},
		"missing required params - validation error": {
			params: GetCASetDeletionStatusRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list ca set deletion status failed: struct validation: CASetID: cannot be blank", err.Error())
			},
		},
		"404 ca set not found": {
			params: GetCASetDeletionStatusRequest{
				CASetID: "2",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/2/status/delete",
			responseStatus: http.StatusNotFound,
			responseBody: `{
 "contextInfo" : {
   "caSetId" : "2"
 },
 "detail" : "Cannot get CA set deletions as the CA set with caSetId 2 is not found.",
 "status" : 404,
 "title" : "CA set is not found.",
 "type" : "/mtls-edge-truststore/error-types/ca-set-not-found"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetNotFound), "want: %s; got: %s", ErrGetCASetNotFound, err)
			},
		},
		"400 - ca set is not during delete": {
			params: GetCASetDeletionStatusRequest{
				CASetID: "1",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/status/delete",
			responseStatus: http.StatusBadRequest,
			responseBody: `{
    "contextInfo": {
        "caSetId": "1"
    },
    "detail": "No active deletions were found for CA Certificate Set with set ID 1.",
    "status": 400,
    "title": "No active cert deletions found",
    "type": "/mtls-edge-truststore/error-types/no-active-cert-deletions"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrNoActiveCertDeletions), "want: %s; got: %s", ErrNoActiveCertDeletions, err)
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
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCASetDeletionStatus(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListCASetActivities(t *testing.T) {
	tests := map[string]struct {
		params           ListCASetActivitiesRequest
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *ListCASetActivitiesResponse
		withError        func(*testing.T, error)
	}{
		"200 - no query params": {
			params: ListCASetActivitiesRequest{
				CASetID: "1",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/activities",
			responseStatus: http.StatusOK,
			responseBody: `{
    "activities": [
        {
            "activityBy": "user1",
            "activityDate": "2025-04-18T11:30:35.866481Z",
            "network": "PRODUCTION",
            "type": "DELETE_CA_SET",
            "version": null
        },
        {
            "activityBy": "user1",
            "activityDate": "2025-04-18T11:30:35.852930Z",
            "network": "STAGING",
            "type": "DELETE_CA_SET",
            "version": null
        },
        {
            "activityBy": "user1",
            "activityDate": "2025-04-16T13:14:46.183261Z",
            "network": null,
            "type": "CREATE_CA_SET_VERSION",
            "version": 2
        },
        {
            "activityBy": "user1",
            "activityDate": "2025-04-15T13:35:48.292127Z",
            "network": null,
            "type": "CREATE_CA_SET_VERSION",
            "version": 1
        },
        {
            "activityBy": "user1",
            "activityDate": "2025-04-15T13:35:48.257574Z",
            "network": null,
            "type": "CREATE_CA_SET",
            "version": null
        }
    ],
    "caSetId": "1",
    "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1",
    "caSetName": "test1",
    "caSetStatus": "DELETED",
    "createdBy": "user1",
    "createdDate": "2025-04-15T13:35:48.211999Z",
    "deletedBy": "user1",
    "deletedDate": "2025-04-18T11:31:40.225213Z"
}`,
			expectedResponse: &ListCASetActivitiesResponse{
				Activities: []CASetActivity{
					{
						ActivityBy:   "user1",
						ActivityDate: test.NewTimeFromString(t, "2025-04-18T11:30:35.866481Z"),
						Network:      ptr.To("PRODUCTION"),
						Type:         "DELETE_CA_SET",
					},
					{
						ActivityBy:   "user1",
						ActivityDate: test.NewTimeFromString(t, "2025-04-18T11:30:35.852930Z"),
						Network:      ptr.To("STAGING"),
						Type:         "DELETE_CA_SET",
					},
					{
						ActivityBy:   "user1",
						ActivityDate: test.NewTimeFromString(t, "2025-04-16T13:14:46.183261Z"),
						Type:         "CREATE_CA_SET_VERSION",
						Version:      ptr.To(int64(2)),
					},
					{
						ActivityBy:   "user1",
						ActivityDate: test.NewTimeFromString(t, "2025-04-15T13:35:48.292127Z"),
						Type:         "CREATE_CA_SET_VERSION",
						Version:      ptr.To(int64(1)),
					},
					{
						ActivityBy:   "user1",
						ActivityDate: test.NewTimeFromString(t, "2025-04-15T13:35:48.257574Z"),
						Type:         "CREATE_CA_SET",
					},
				},
				CASetID:     "1",
				CASetLink:   "/mtls-edge-truststore/v2/ca-sets/1",
				CASetName:   "test1",
				CASetStatus: "DELETED",
				CreatedBy:   "user1",
				CreatedDate: test.NewTimeFromString(t, "2025-04-15T13:35:48.211999Z"),
				DeletedBy:   ptr.To("user1"),
				DeletedDate: ptr.To(test.NewTimeFromString(t, "2025-04-18T11:31:40.225213Z")),
			},
		},
		"200 - all query params": {
			params: ListCASetActivitiesRequest{
				CASetID: "1",
				Start:   test.NewTimeFromString(t, "2025-04-15T14:00:00Z"),
				End:     test.NewTimeFromString(t, "2025-04-17T14:00:00.00000Z"),
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/1/activities?end=2025-04-17T14%3A00%3A00Z&start=2025-04-15T14%3A00%3A00Z",
			responseStatus: http.StatusOK,
			responseBody: `{
    "activities": [
        {
            "activityBy": "user1",
            "activityDate": "2025-04-16T13:14:46.183261Z",
            "network": null,
            "type": "CREATE_CA_SET_VERSION",
            "version": 2
        }
    ],
    "caSetId": "1",
    "caSetLink": "/mtls-edge-truststore/v2/ca-sets/1",
    "caSetName": "test1",
    "caSetStatus": "DELETED",
    "createdBy": "user1",
    "createdDate": "2025-04-15T13:35:48.211999Z",
    "deletedBy": "user1",
    "deletedDate": "2025-04-18T11:31:40.225213Z"
}`,
			expectedResponse: &ListCASetActivitiesResponse{
				Activities: []CASetActivity{
					{
						ActivityBy:   "user1",
						ActivityDate: test.NewTimeFromString(t, "2025-04-16T13:14:46.183261Z"),
						Type:         "CREATE_CA_SET_VERSION",
						Version:      ptr.To(int64(2)),
					},
				},
				CASetID:     "1",
				CASetLink:   "/mtls-edge-truststore/v2/ca-sets/1",
				CASetName:   "test1",
				CASetStatus: "DELETED",
				CreatedBy:   "user1",
				CreatedDate: test.NewTimeFromString(t, "2025-04-15T13:35:48.211999Z"),
				DeletedBy:   ptr.To("user1"),
				DeletedDate: ptr.To(test.NewTimeFromString(t, "2025-04-18T11:31:40.225213Z")),
			},
		},
		"missing required params - validation error": {
			params: ListCASetActivitiesRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get ca set activities failed: struct validation: CASetID: cannot be blank", err.Error())
			},
		},
		"404 ca set not found": {
			params: ListCASetActivitiesRequest{
				CASetID: "2",
			},
			expectedPath:   "/mtls-edge-truststore/v2/ca-sets/2/activities",
			responseStatus: http.StatusNotFound,
			responseBody: `{
    "contextInfo": {
        "caSetId": "2"
    },
    "detail": "Cannot get CA set activities as the CA set with caSetId 2 is not found.",
    "status": 404,
    "title": "CA set is not found.",
    "type": "/mtls-edge-truststore/error-types/ca-set-not-found"
}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrGetCASetNotFound), "want: %s; got: %s", ErrGetCASetNotFound, err)
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
			client := mockAPIClient(t, mockServer)
			result, err := client.ListCASetActivities(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
