package netstorage

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetStorageListStorageGroups(t *testing.T) {
	tests := map[string]struct {
		request          ListStorageGroupsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListStorageGroupsResponse
		withError        error
	}{
		"200 OK": {
			request: ListStorageGroupsRequest{
				CPCodeID:            "12345",
				StorageGroupPurpose: "NETSTORAGE",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "items": [
        {
            "contractId": "P-12A3B4C",
            "storageGroupId": 12345,
            "storageGroupName": "My Website Failover",
            "storageGroupType": "OBJECTSTORE",
            "storageGroupPurpose": "NETSTORAGE",
            "domainPrefix": "webfailover",
            "asperaEnabled": false,
            "pciEnabled": false,
            "estimatedUsageGB": 0.01,
            "allowEdit": true,
            "provisionStatus": "PROVISIONED",
            "cpcodes": [
                {
                    "cpcodeId": 123456,
                    "downloadSecurity": "ALL_EDGE_SERVERS",
                    "useSsl": false,
                    "serveFromZip": false,
                    "sendHash": false,
                    "quickDelete": false,
                    "numberOfFiles": 154,
                    "numberOfBytes": 1,
                    "lastChangesPropagated": true,
                    "requestUriCaseConversion": "NO_CONVERSION",
                    "queryStringConversion": {
                        "queryStringConversionMode": "STRIP_ALL_INCOMING_QUERY"
                    },
                    "encodingConfig": {
                        "enforceEncoding": false
                    },
                    "dirListing": {
                        "indexFileName": "index.html",
                        "maxListSize": 0,
                        "searchOn404": "DO_NOT_SEARCH"
                    },
                    "lastModifiedBy": "some.body@myemail.net",
                    "lastModifiedDate": "2024-11-11T11:50:59Z"
                },
                {
                    "cpcodeId": 678901,
                    "downloadSecurity": "ALL_EDGE_SERVERS",
                    "useSsl": false,
                    "serveFromZip": false,
                    "sendHash": false,
                    "quickDelete": false,
                    "numberOfFiles": 67,
                    "numberOfBytes": 1,
                    "lastChangesPropagated": true,
                    "requestUriCaseConversion": "NO_CONVERSION",
                    "queryStringConversion": {
                        "queryStringConversionMode": "STRIP_ALL_INCOMING_QUERY"
                    },
                    "encodingConfig": {
                        "enforceEncoding": false
                    },
                    "dirListing": {
                        "indexFileName": "index.html",
                        "maxListSize": 0,
                        "searchOn404": "DO_NOT_SEARCH"
                    },
                    "lastModifiedBy": "some.body@myemail.net",
                    "lastModifiedDate": "2024-11-11T11:50:59Z"
                },
                {
                    "cpcodeId": 1304188,
                    "downloadSecurity": "ALL_EDGE_SERVERS",
                    "useSsl": false,
                    "serveFromZip": false,
                    "sendHash": false,
                    "quickDelete": false,
                    "numberOfFiles": 0,
                    "numberOfBytes": 0,
                    "lastChangesPropagated": true,
                    "requestUriCaseConversion": "NO_CONVERSION",
                    "queryStringConversion": {
                        "queryStringConversionMode": "STRIP_ALL_INCOMING_QUERY"
                    },
                    "pathCheckAndConversion": "DO_NOT_CHECK_PATHS",
                    "encodingConfig": {
                        "enforceEncoding": false,
                        "encoding": "UTF_8"
                    },
                    "dirListing": {
                        "maxListSize": 0,
                        "searchOn404": "DO_NOT_SEARCH"
                    },
                    "lastModifiedBy": "some.body@myemail.net",
                    "lastModifiedDate": "2024-11-11T11:50:59Z"
                },
                {
                    "cpcodeId": 1726184,
                    "downloadSecurity": "ALL_EDGE_SERVERS",
                    "useSsl": true,
                    "serveFromZip": false,
                    "sendHash": false,
                    "quickDelete": false,
                    "numberOfFiles": 42,
                    "numberOfBytes": 1,
                    "lastChangesPropagated": true,
                    "requestUriCaseConversion": "NO_CONVERSION",
                    "queryStringConversion": {
                        "queryStringConversionMode": "STRIP_ALL_INCOMING_QUERY"
                    },
                    "pathCheckAndConversion": "DO_NOT_CHECK_PATHS",
                    "encodingConfig": {
                        "enforceEncoding": false,
                        "encoding": "UTF_8"
                    },
                    "dirListing": {
                        "maxListSize": 0,
                        "searchOn404": "DO_NOT_SEARCH"
                    },
                    "lastModifiedBy": "some.body@myemail.net",
                    "lastModifiedDate": "2024-11-11T11:50:59Z"
                },
                {
                    "cpcodeId": 1726185,
                    "downloadSecurity": "ALL_EDGE_SERVERS",
                    "useSsl": true,
                    "serveFromZip": false,
                    "sendHash": false,
                    "quickDelete": false,
                    "numberOfFiles": 14,
                    "numberOfBytes": 1,
                    "lastChangesPropagated": true,
                    "requestUriCaseConversion": "NO_CONVERSION",
                    "queryStringConversion": {
                        "queryStringConversionMode": "STRIP_ALL_INCOMING_QUERY"
                    },
                    "pathCheckAndConversion": "DO_NOT_CHECK_PATHS",
                    "encodingConfig": {
                        "enforceEncoding": false,
                        "encoding": "UTF_8"
                    },
                    "dirListing": {
                        "maxListSize": 0,
                        "searchOn404": "DO_NOT_SEARCH"
                    },
                    "lastModifiedBy": "some.body@myemail.net",
                    "lastModifiedDate": "2024-11-11T11:50:59Z"
                }
            ],
            "zones": [
                {
                    "zoneName": "europe",
                    "noCapacityAction": "SPILL_OUTSIDE",
                    "allowUpload": "YES",
                    "allowDownload": "YES",
                    "lastModifiedBy": "somebody",
                    "lastModifiedDate": "2020-05-07T22:01:10Z"
                },
                {
                    "zoneName": "europe",
                    "noCapacityAction": "SPILL_OUTSIDE",
                    "allowUpload": "YES",
                    "allowDownload": "YES",
                    "lastModifiedBy": "somebody",
                    "lastModifiedDate": "2020-05-07T22:01:10Z"
                }
            ],
            "propagationStatus": {
                "status": "ACTIVE"
            },
            "lastModifiedBy": "some.body@myemail.net",
            "lastModifiedDate": "2024-11-19T15:20:20Z",
            "links": [
                {
                    "rel": "self",
                    "href": "/storage-services//storage-groups/12345"
                },
                {
                    "rel": "uploadAccounts",
                    "href": "/storage-services//upload-accounts?storageGroupId=12345"
                },
                {
                    "rel": "zones",
                    "href": "/storage-services//zones"
                }
            ]
        }
    ]
}`,
			expectedPath: "/storage/v1/storage-groups?cpcodeId=12345&storageGroupPurpose=NETSTORAGE",
			expectedResponse: &ListStorageGroupsResponse{
				Items: []*StorageGroup{{
					ContractID:          "P-12A3B4C",
					StorageGroupID:      12345,
					StorageGroupName:    "My Website Failover",
					StorageGroupType:    StorageGroupTypeObjectStore,
					StorageGroupPurpose: StorageGroupPurposeNetStorage,
					DomainPrefix:        "webfailover",
					AsperaEnabled:       false,
					PciEnabled:          false,
					EstimatedUsageGB:    0.01,
					AllowEdit:           true,
					ProvisionStatus:     StorageGroupProvisionStatusProvisioned,
					PropagationStatus:   StorageGroupPropagationStatus{Status: PropagationStatusActive},
					LastModifiedBy:      "some.body@myemail.net",
					LastModifiedDate:    "2024-11-19T15:20:20Z",
					CPCodes: []*StorageGroupCPCode{
						{
							CPCodeID:                 123456,
							DownloadSecurity:         "ALL_EDGE_SERVERS",
							UseSSL:                   false,
							ServeFromZip:             false,
							SendHash:                 false,
							QuickDelete:              false,
							NumberOfFiles:            154,
							NumberOfBytes:            1,
							LastChangesPropagated:    true,
							RequestURICaseConversion: RequestURICaseConversionNoConversion,
							QueryStringConversion: QueryStringConversion{
								QueryStringConversionMode: QueryStringConversionModeStripAllIncomingQuery,
							},
							EncodingConfig: EncodingConfig{
								EnforceEncoding: false,
							},
							DirListing: DirListing{
								IndexFileName: "index.html",
								MaxListSize:   0,
								SearchOn404:   DirListingSearchOn404DoNotSearch,
							},
							LastModifiedBy:   "some.body@myemail.net",
							LastModifiedDate: "2024-11-11T11:50:59Z",
						},
						{
							CPCodeID:                 678901,
							DownloadSecurity:         "ALL_EDGE_SERVERS",
							UseSSL:                   false,
							ServeFromZip:             false,
							SendHash:                 false,
							QuickDelete:              false,
							NumberOfFiles:            67,
							NumberOfBytes:            1,
							LastChangesPropagated:    true,
							RequestURICaseConversion: RequestURICaseConversionNoConversion,
							QueryStringConversion: QueryStringConversion{
								QueryStringConversionMode: QueryStringConversionModeStripAllIncomingQuery,
							},
							EncodingConfig: EncodingConfig{
								EnforceEncoding: false,
							},
							DirListing: DirListing{
								IndexFileName: "index.html",
								MaxListSize:   0,
								SearchOn404:   DirListingSearchOn404DoNotSearch,
							},
							LastModifiedBy:   "some.body@myemail.net",
							LastModifiedDate: "2024-11-11T11:50:59Z",
						},
						{
							CPCodeID:                 1304188,
							DownloadSecurity:         "ALL_EDGE_SERVERS",
							UseSSL:                   false,
							ServeFromZip:             false,
							SendHash:                 false,
							QuickDelete:              false,
							NumberOfFiles:            0,
							NumberOfBytes:            0,
							LastChangesPropagated:    true,
							RequestURICaseConversion: RequestURICaseConversionNoConversion,
							QueryStringConversion: QueryStringConversion{
								QueryStringConversionMode: QueryStringConversionModeStripAllIncomingQuery,
							},
							PathCheckAndConversion: PathCheckAndConversionDoNotCheckPaths,
							EncodingConfig: EncodingConfig{
								EnforceEncoding: false,
								Encoding:        EncodingUTF8,
							},
							DirListing: DirListing{
								MaxListSize: 0,
								SearchOn404: DirListingSearchOn404DoNotSearch,
							},
							LastModifiedBy:   "some.body@myemail.net",
							LastModifiedDate: "2024-11-11T11:50:59Z",
						},
						{
							CPCodeID:                 1726184,
							DownloadSecurity:         "ALL_EDGE_SERVERS",
							UseSSL:                   true,
							ServeFromZip:             false,
							SendHash:                 false,
							QuickDelete:              false,
							NumberOfFiles:            42,
							NumberOfBytes:            1,
							LastChangesPropagated:    true,
							RequestURICaseConversion: RequestURICaseConversionNoConversion,
							QueryStringConversion: QueryStringConversion{
								QueryStringConversionMode: QueryStringConversionModeStripAllIncomingQuery,
							},
							PathCheckAndConversion: PathCheckAndConversionDoNotCheckPaths,
							EncodingConfig: EncodingConfig{
								EnforceEncoding: false,
								Encoding:        EncodingUTF8,
							},
							DirListing: DirListing{
								MaxListSize: 0,
								SearchOn404: DirListingSearchOn404DoNotSearch,
							},
							LastModifiedBy:   "some.body@myemail.net",
							LastModifiedDate: "2024-11-11T11:50:59Z",
						},
						{
							CPCodeID:                 1726185,
							DownloadSecurity:         "ALL_EDGE_SERVERS",
							UseSSL:                   true,
							ServeFromZip:             false,
							SendHash:                 false,
							QuickDelete:              false,
							NumberOfFiles:            14,
							NumberOfBytes:            1,
							LastChangesPropagated:    true,
							RequestURICaseConversion: RequestURICaseConversionNoConversion,
							QueryStringConversion: QueryStringConversion{
								QueryStringConversionMode: QueryStringConversionModeStripAllIncomingQuery,
							},
							PathCheckAndConversion: PathCheckAndConversionDoNotCheckPaths,
							EncodingConfig: EncodingConfig{
								EnforceEncoding: false,
								Encoding:        EncodingUTF8,
							},
							DirListing: DirListing{
								MaxListSize: 0,
								SearchOn404: DirListingSearchOn404DoNotSearch,
							},
							LastModifiedBy:   "some.body@myemail.net",
							LastModifiedDate: "2024-11-11T11:50:59Z",
						},
					},
				},
				},
			},
		},
		"500 internal server error": {
			request:        ListStorageGroupsRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching activation",
    "status": 500
}`,
			expectedPath: "/storage/v1/storage-groups",
			withError: &Error{
				Type:   "internal_error",
				Title:  "Internal Server Error",
				Detail: "Error fetching activation",
				Status: http.StatusInternalServerError,
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
			result, err := client.ListStorageGroups(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestNetStorageGetStorageGroup(t *testing.T) {
	tests := map[string]struct {
		request          GetStorageGroupRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetStorageGroupResponse
		withError        error
	}{
		"200 OK": {
			request: GetStorageGroupRequest{
				StorageGroupID: 12345,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "contractId": "P-12A3B4C",
    "storageGroupId": 12345,
    "storageGroupName": "My Website Failover",
    "storageGroupType": "OBJECTSTORE",
    "storageGroupPurpose": "NETSTORAGE",
    "domainPrefix": "webfailover",
    "asperaEnabled": false,
    "pciEnabled": false,
    "estimatedUsageGB": 0.01,
    "allowEdit": true,
    "provisionStatus": "PROVISIONED",
    "cpcodes": [
        {
            "cpcodeId": 123456,
            "downloadSecurity": "ALL_EDGE_SERVERS",
            "useSsl": false,
            "serveFromZip": false,
            "sendHash": false,
            "quickDelete": false,
            "numberOfFiles": 154,
            "numberOfBytes": 1,
            "lastChangesPropagated": true,
            "requestUriCaseConversion": "NO_CONVERSION",
            "queryStringConversion": {
                "queryStringConversionMode": "STRIP_ALL_INCOMING_QUERY"
            },
            "encodingConfig": {
                "enforceEncoding": false
            },
            "dirListing": {
                "indexFileName": "index.html",
                "maxListSize": 0,
                "searchOn404": "DO_NOT_SEARCH"
            },
            "lastModifiedBy": "some.body@myemail.net",
            "lastModifiedDate": "2024-11-11T11:50:59Z"
        }
    ],
    "zones": [
        {
            "zoneName": "europe",
            "noCapacityAction": "SPILL_OUTSIDE",
            "allowUpload": "YES",
            "allowDownload": "YES",
            "lastModifiedBy": "somebody",
            "lastModifiedDate": "2020-05-07T22:01:10Z"
        },
        {
            "zoneName": "europe",
            "noCapacityAction": "SPILL_OUTSIDE",
            "allowUpload": "YES",
            "allowDownload": "YES",
            "lastModifiedBy": "somebody",
            "lastModifiedDate": "2020-05-07T22:01:10Z"
        }
    ],
    "propagationStatus": {
        "status": "ACTIVE"
    },
    "lastModifiedBy": "some.body@myemail.net",
    "lastModifiedDate": "2024-11-19T15:20:20Z",
    "links": [
        {
            "rel": "self",
            "href": "/storage-services//storage-groups/12345"
        },
        {
            "rel": "uploadAccounts",
            "href": "/storage-services//upload-accounts?storageGroupId=12345"
        },
        {
            "rel": "zones",
            "href": "/storage-services//zones"
        }
    ]
}`,
			expectedPath: "/storage/v1/storage-groups/12345",
			expectedResponse: &GetStorageGroupResponse{
				ContractID:          "P-12A3B4C",
				StorageGroupID:      12345,
				StorageGroupName:    "My Website Failover",
				StorageGroupType:    StorageGroupTypeObjectStore,
				StorageGroupPurpose: StorageGroupPurposeNetStorage,
				DomainPrefix:        "webfailover",
				AsperaEnabled:       false,
				PciEnabled:          false,
				EstimatedUsageGB:    0.01,
				AllowEdit:           true,
				ProvisionStatus:     StorageGroupProvisionStatusProvisioned,
				PropagationStatus:   StorageGroupPropagationStatus{Status: PropagationStatusActive},
				LastModifiedBy:      "some.body@myemail.net",
				LastModifiedDate:    "2024-11-19T15:20:20Z",
				CPCodes: []*StorageGroupCPCode{
					{
						CPCodeID:                 123456,
						DownloadSecurity:         "ALL_EDGE_SERVERS",
						UseSSL:                   false,
						ServeFromZip:             false,
						SendHash:                 false,
						QuickDelete:              false,
						NumberOfFiles:            154,
						NumberOfBytes:            1,
						LastChangesPropagated:    true,
						RequestURICaseConversion: RequestURICaseConversionNoConversion,
						QueryStringConversion: QueryStringConversion{
							QueryStringConversionMode: QueryStringConversionModeStripAllIncomingQuery,
						},
						EncodingConfig: EncodingConfig{
							EnforceEncoding: false,
						},
						DirListing: DirListing{
							IndexFileName: "index.html",
							MaxListSize:   0,
							SearchOn404:   DirListingSearchOn404DoNotSearch,
						},
						LastModifiedBy:   "some.body@myemail.net",
						LastModifiedDate: "2024-11-11T11:50:59Z",
					},
				},
			},
		},
		"storage group not found": {
			request: GetStorageGroupRequest{
				StorageGroupID: 123456789,
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
    "type": "validation-error",
    "title": "Validation failure",
    "instance": "c162061b-01bf-49f8-b8d8-eb97e9ee476b",
    "status": 400,
    "detail": "Validation failed. Please review the errors.",
    "errors": [
        {
            "type": "error-types/invalid-value",
            "title": "Invalid value",
            "detail": "Unable to find the given storage group.",
            "problemId": "0c433e81-f64f-465e-aeb4-d834b5bc717f",
            "field": "storageGroupId"
        }
    ],
    "problemId": "c162061b-01bf-49f8-b8d8-eb97e9ee476b"
}`,
			expectedPath: "/storage/v1/storage-groups/123456789",
			withError:    ErrNotFound,
		},
		"500 internal server error": {
			request: GetStorageGroupRequest{
				StorageGroupID: 123456,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching activation",
    "status": 500
}`,
			expectedPath: "/storage/v1/storage-groups/123456",
			withError: &Error{
				Type:   "internal_error",
				Title:  "Internal Server Error",
				Detail: "Error fetching activation",
				Status: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request:   GetStorageGroupRequest{},
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
			result, err := client.GetStorageGroup(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
