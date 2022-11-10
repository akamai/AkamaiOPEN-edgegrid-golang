package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDns_GetBulkZoneCreateStatus(t *testing.T) {
	tests := map[string]struct {
		requestid        string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *BulkStatusResponse
		withError        error
	}{
		"200 OK": {
			requestid:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
			responseStatus: http.StatusOK,
			responseBody: `
  {
    "requestId": "15bc138f-8d82-451b-80b7-a56b88ffc474",
    "zonesSubmitted": 2,
    "successCount": 0,
    "failureCount": 2,
    "isComplete": true,
    "expirationDate": "2020-10-28T17:10:04.515792Z"
  }`,
			expectedPath: "/config-dns/v2/zones/create-requests/15bc138f-8d82-451b-80b7-a56b88ffc474",
			expectedResponse: &BulkStatusResponse{
				RequestId:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
				ZonesSubmitted: 2,
				SuccessCount:   0,
				FailureCount:   2,
				IsComplete:     true,
				ExpirationDate: "2020-10-28T17:10:04.515792Z",
			},
		},
		"500 internal server error": {
			requestid:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/create-requests/15bc138f-8d82-451b-80b7-a56b88ffc474",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetBulkZoneCreateStatus(context.Background(), test.requestid)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetBulkZoneCreateResult(t *testing.T) {
	tests := map[string]struct {
		requestid        string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *BulkCreateResultResponse
		withError        error
	}{
		"200 OK": {
			requestid:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
			responseStatus: http.StatusOK,
			responseBody: `
  {
    "requestId": "15bc138f-8d82-451b-80b7-a56b88ffc474",
    "successfullyCreatedZones": [],
    "failedZones": [
      {
        "zone": "one.testbulk.net",
        "failureReason": "ZONE_ALREADY_EXISTS"
      }
    ]
  }`,
			expectedPath: "/config-dns/v2/zones/create-requests/15bc138f-8d82-451b-80b7-a56b88ffc474/result",
			expectedResponse: &BulkCreateResultResponse{
				RequestId:                "15bc138f-8d82-451b-80b7-a56b88ffc474",
				SuccessfullyCreatedZones: make([]string, 0),
				FailedZones: []*BulkFailedZone{
					{
						Zone:          "one.testbulk.net",
						FailureReason: "ZONE_ALREADY_EXISTS",
					},
				},
			},
		},
		"500 internal server error": {
			requestid:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
        "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/create-requests/15bc138f-8d82-451b-80b7-a56b88ffc474/result",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetBulkZoneCreateResult(context.Background(), test.requestid)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_CreateBulkZones(t *testing.T) {
	tests := map[string]struct {
		zones            BulkZonesCreate
		query            ZoneQueryString
		responseStatus   int
		responseBody     string
		expectedResponse *BulkZonesResponse
		expectedPath     string
		withError        error
	}{
		"200 Created": {
			zones: BulkZonesCreate{
				Zones: []*ZoneCreate{
					{
						Zone:    "one.testbulk.net",
						Type:    "secondary",
						Comment: "testing bulk operations",
						Masters: []string{"1.2.3.4", "1.2.3.10"},
					},
					{
						Zone:    "two.testbulk.net",
						Type:    "secondary",
						Comment: "testing bulk operations",
						Masters: []string{"1.2.3.6", "1.2.3.70"},
					},
				},
			},
			query:          ZoneQueryString{Contract: "1-2ABCDE", Group: "testgroup"},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "requestId": "93e97a28-4e05-45f4-8b9a-cebd71155949",
    "expirationDate": "2020-10-28T19:50:36.272668Z"
}`,
			expectedResponse: &BulkZonesResponse{
				RequestId:      "93e97a28-4e05-45f4-8b9a-cebd71155949",
				ExpirationDate: "2020-10-28T19:50:36.272668Z",
			},
			expectedPath: "/config-dns/v2/zones/create-requests?contractId=1-2ABCDE&gid=testgroup",
		},
		"500 internal server error": {
			zones: BulkZonesCreate{
				Zones: []*ZoneCreate{
					{
						Zone:    "one.testbulk.net",
						Type:    "secondary",
						Comment: "testing bulk operations",
						Masters: []string{"1.2.3.4", "1.2.3.10"},
					},
					{
						Zone:    "two.testbulk.net",
						Type:    "secondary",
						Comment: "testing bulk operations",
						Masters: []string{"1.2.3.6", "1.2.3.70"},
					},
				},
			},
			query:          ZoneQueryString{Contract: "1-2ABCDE", Group: "testgroup"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/create-requests?contractId=1-2ABCDE&gid=testgroup",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateBulkZones(context.Background(), &test.zones, test.query)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Bulk Delete tests
func TestDns_GetBulkZoneDeleteStatus(t *testing.T) {
	tests := map[string]struct {
		requestid        string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *BulkStatusResponse
		withError        error
	}{
		"200 OK": {
			requestid:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
			responseStatus: http.StatusOK,
			responseBody: `
  {
    "requestId": "15bc138f-8d82-451b-80b7-a56b88ffc474",
    "zonesSubmitted": 2,
    "successCount": 0,
    "failureCount": 2,
    "isComplete": true,
    "expirationDate": "2020-10-28T17:10:04.515792Z"
  }`,
			expectedPath: "/config-dns/v2/zones/delete-requests/15bc138f-8d82-451b-80b7-a56b88ffc474",
			expectedResponse: &BulkStatusResponse{
				RequestId:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
				ZonesSubmitted: 2,
				SuccessCount:   0,
				FailureCount:   2,
				IsComplete:     true,
				ExpirationDate: "2020-10-28T17:10:04.515792Z",
			},
		},
		"500 internal server error": {
			requestid:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
        "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/delete-requests/15bc138f-8d82-451b-80b7-a56b88ffc474",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetBulkZoneDeleteStatus(context.Background(), test.requestid)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetBulkZoneDeleteResult(t *testing.T) {
	tests := map[string]struct {
		requestid        string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *BulkDeleteResultResponse
		withError        error
	}{
		"200 OK": {
			requestid:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
			responseStatus: http.StatusOK,
			responseBody: `
  {
    "requestId": "15bc138f-8d82-451b-80b7-a56b88ffc474",
    "successfullyDeletedZones": [],
    "failedZones": [
      {
        "zone": "one.testbulk.net",
        "failureReason": "ZONE_ALREADY_EXISTS"
      }
    ]
  }`,
			expectedPath: "/config-dns/v2/zones/delete-requests/15bc138f-8d82-451b-80b7-a56b88ffc474/result",
			expectedResponse: &BulkDeleteResultResponse{
				RequestId:                "15bc138f-8d82-451b-80b7-a56b88ffc474",
				SuccessfullyDeletedZones: make([]string, 0),
				FailedZones: []*BulkFailedZone{
					{
						Zone:          "one.testbulk.net",
						FailureReason: "ZONE_ALREADY_EXISTS",
					},
				},
			},
		},
		"500 internal server error": {
			requestid:      "15bc138f-8d82-451b-80b7-a56b88ffc474",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
        "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/create-requests/15bc138f-8d82-451b-80b7-a56b88ffc474/result",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetBulkZoneDeleteResult(context.Background(), test.requestid)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_DeleteBulkZones(t *testing.T) {
	tests := map[string]struct {
		zoneslist        ZoneNameListResponse
		bypassSafety     bool
		responseStatus   int
		responseBody     string
		expectedResponse *BulkZonesResponse
		expectedPath     string
		withError        error
	}{
		"200 Created": {
			zoneslist: ZoneNameListResponse{
				Zones: []string{"one.testbulk.net", "two.testbulk.net"},
			},
			bypassSafety:   true,
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "requestId": "93e97a28-4e05-45f4-8b9a-cebd71155949",
    "expirationDate": "2020-10-28T19:50:36.272668Z"
}`,
			expectedResponse: &BulkZonesResponse{
				RequestId:      "93e97a28-4e05-45f4-8b9a-cebd71155949",
				ExpirationDate: "2020-10-28T19:50:36.272668Z",
			},
			expectedPath: "/config-dns/v2/zones/delete-requests?bypassSafetyChecks=true",
		},
		"500 internal server error": {
			zoneslist: ZoneNameListResponse{
				Zones: []string{"one.testbulk.net", "two.testbulk.net"},
			},
			bypassSafety:   true,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
        "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/zones/delete-requests?bypassSafetyChecks=true",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteBulkZones(context.Background(), &test.zoneslist, test.bypassSafety)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
