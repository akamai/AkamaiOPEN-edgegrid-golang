package cloudaccess

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupProperties(t *testing.T) {
	tests := map[string]struct {
		request          LookupPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *LookupPropertiesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK empty": {
			request:        LookupPropertiesRequest{AccessKeyUID: 1234, Version: 1},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "properties": []
}
`,
			expectedPath: "/cam/v1/access-keys/1234/versions/1/properties",
			expectedResponse: &LookupPropertiesResponse{
				Properties: []Property{},
			},
		},
		"200 OK with properties": {
			request:        LookupPropertiesRequest{AccessKeyUID: 1234, Version: 1},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "properties": [
        {
            "accessKeyUid": 1234,
            "version": 1,
            "propertyId": "prp_5678",
            "propertyName": "test-property",
            "productionVersion": null,
            "stagingVersion": 1
        },
        {
            "accessKeyUid": 1234,
            "version": 1,
            "propertyId": "prp_6789",
            "propertyName": "test-property2",
            "productionVersion": 1,
            "stagingVersion": null
        }
    ]
}
`,
			expectedPath: "/cam/v1/access-keys/1234/versions/1/properties",
			expectedResponse: &LookupPropertiesResponse{
				Properties: []Property{
					{
						AccessKeyUID:      1234,
						Version:           1,
						PropertyID:        "prp_5678",
						PropertyName:      "test-property",
						ProductionVersion: nil,
						StagingVersion:    ptr.To(int64(1)),
					},
					{
						AccessKeyUID:      1234,
						Version:           1,
						PropertyID:        "prp_6789",
						PropertyName:      "test-property2",
						ProductionVersion: ptr.To(int64(1)),
						StagingVersion:    nil,
					},
				},
			},
		},
		"404 incorrect request": {
			request:        LookupPropertiesRequest{AccessKeyUID: 1234, Version: 10},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/cam/error-types/access-key-version-does-not-exist",
    "title": "Domain Error",
    "instance": "60126c0d-67f5-473c-bea0-16daa836dc44",
    "status": 404,
    "detail": "Version '10' for access key '1234' does not exist.",
    "problemId": "60126c0d-67f5-473c-bea0-16daa836dc44",
    "version": 10,
    "accessKeyUID": 1234
}
`,
			expectedPath: "/cam/v1/access-keys/1234/versions/10/properties",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:         "/cam/error-types/access-key-version-does-not-exist",
					Title:        "Domain Error",
					Detail:       "Version '10' for access key '1234' does not exist.",
					Instance:     "60126c0d-67f5-473c-bea0-16daa836dc44",
					Status:       404,
					ProblemID:    "60126c0d-67f5-473c-bea0-16daa836dc44",
					AccessKeyUID: 1234,
					Version:      10,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			request:        LookupPropertiesRequest{AccessKeyUID: 1234, Version: 1},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal-server-error",
    "title": "Internal Server Error",
    "detail": "Error processing request",
    "instance": "TestInstances",
    "status": 500
}
`,
			expectedPath: "/cam/v1/access-keys/1234/versions/1/properties",
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
		"validate errors": {
			request: LookupPropertiesRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "lookup properties: struct validation: AccessKeyUID: cannot be blank\nVersion: cannot be blank", err.Error())
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
			result, err := client.LookupProperties(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetAsyncPropertiesLookupID(t *testing.T) {
	tests := map[string]struct {
		request          GetAsyncPropertiesLookupIDRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAsyncPropertiesLookupIDResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request:        GetAsyncPropertiesLookupIDRequest{AccessKeyUID: 1234, Version: 1},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "lookupId": 4321,
    "retryAfter": 10
}
`,
			expectedPath:     "/cam/v1/access-keys/1234/versions/1/property-lookup-id",
			expectedResponse: &GetAsyncPropertiesLookupIDResponse{LookupID: 4321, RetryAfter: 10},
		},
		"404 incorrect request": {
			request:        GetAsyncPropertiesLookupIDRequest{AccessKeyUID: 1234, Version: 10},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/cam/error-types/access-key-version-does-not-exist",
    "title": "Domain Error",
    "instance": "60126c0d-67f5-473c-bea0-16daa836dc44",
    "status": 404,
    "detail": "Version '10' for access key '1234' does not exist.",
    "problemId": "60126c0d-67f5-473c-bea0-16daa836dc44",
    "version": 10,
    "accessKeyUID": 1234
}
`,
			expectedPath: "/cam/v1/access-keys/1234/versions/10/property-lookup-id",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:         "/cam/error-types/access-key-version-does-not-exist",
					Title:        "Domain Error",
					Detail:       "Version '10' for access key '1234' does not exist.",
					Instance:     "60126c0d-67f5-473c-bea0-16daa836dc44",
					Status:       404,
					ProblemID:    "60126c0d-67f5-473c-bea0-16daa836dc44",
					AccessKeyUID: 1234,
					Version:      10,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			request:        GetAsyncPropertiesLookupIDRequest{AccessKeyUID: 1234, Version: 1},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal-server-error",
    "title": "Internal Server Error",
    "detail": "Error processing request",
    "instance": "TestInstances",
    "status": 500
}
`,
			expectedPath: "/cam/v1/access-keys/1234/versions/1/property-lookup-id",
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
		"validate errors": {
			request: GetAsyncPropertiesLookupIDRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get lookup properties id async: struct validation: AccessKeyUID: cannot be blank\nVersion: cannot be blank", err.Error())
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
			result, err := client.GetAsyncPropertiesLookupID(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPerformAsyncPropertiesLookup(t *testing.T) {
	tests := map[string]struct {
		request          PerformAsyncPropertiesLookupRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PerformAsyncPropertiesLookupResponse
		withError        func(*testing.T, error)
	}{
		"200 OK empty": {
			request:        PerformAsyncPropertiesLookupRequest{LookupID: 4321},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "lookupId": 4321,
    "lookupStatus": "COMPLETE",
    "properties": []
}
`,
			expectedPath: "/cam/v1/property-lookups/4321",
			expectedResponse: &PerformAsyncPropertiesLookupResponse{
				LookupID:     4321,
				LookupStatus: "COMPLETE",
				Properties:   []Property{},
			},
		},
		"200 OK with properties": {
			request:        PerformAsyncPropertiesLookupRequest{LookupID: 4321},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "lookupId": 4321,
    "lookupStatus": "COMPLETE",
    "properties": [
        {
            "accessKeyUid": 1234,
            "version": 1,
            "propertyId": "prp_5678",
            "propertyName": "test-property",
            "productionVersion": null,
            "stagingVersion": 1
        },
        {
            "accessKeyUid": 1234,
            "version": 1,
            "propertyId": "prp_6789",
            "propertyName": "test-property2",
            "productionVersion": 1,
            "stagingVersion": null
        }
    ]
}
`,
			expectedPath: "/cam/v1/property-lookups/4321",
			expectedResponse: &PerformAsyncPropertiesLookupResponse{
				LookupID:     4321,
				LookupStatus: LookupComplete,
				Properties: []Property{
					{
						AccessKeyUID:      1234,
						Version:           1,
						PropertyID:        "prp_5678",
						PropertyName:      "test-property",
						ProductionVersion: nil,
						StagingVersion:    ptr.To(int64(1)),
					},
					{
						AccessKeyUID:      1234,
						Version:           1,
						PropertyID:        "prp_6789",
						PropertyName:      "test-property2",
						ProductionVersion: ptr.To(int64(1)),
						StagingVersion:    nil,
					},
				},
			},
		},
		"500 internal server error": {
			request:        PerformAsyncPropertiesLookupRequest{LookupID: 4321},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal-server-error",
		   "title": "Internal Server Error",
		   "detail": "Error processing request",
		   "instance": "TestInstances",
		   "status": 500
		}
		`,
			expectedPath: "/cam/v1/property-lookups/4321",
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
		"validate errors": {
			request: PerformAsyncPropertiesLookupRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "perform async lookup properties: struct validation: LookupID: cannot be blank", err.Error())
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
			result, err := client.PerformAsyncPropertiesLookup(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
