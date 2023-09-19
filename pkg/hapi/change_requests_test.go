package hapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChangeRequest(t *testing.T) {
	tests := map[string]struct {
		request          GetChangeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ChangeRequest
		withError        error
	}{
		"200 OK": {
			request:        GetChangeRequest{123},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"action": "EDIT",
	"changeId": 123,
	"edgeHostnames": [
		{
			"chinaCdn": {
				"isChinaCdn": false
			},
			"dnsZone": "edgekey.net",
			"edgeHostnameId": 112233,
			"ipVersionBehavior": "IPV6_IPV4_DUALSTACK",
			"map": "a;bcd.akamaiedge.net",
			"productId": "DSA",
			"recordName": "test123",
			"securityType": "ENHANCED-TLS",
			"serialNumber": 0,
			"slotNumber": 1234,
			"ttl": 21600,
			"useDefaultMap": true,
			"useDefaultTtl": true
		}
	],
	"status": "PENDING",
	"statusMessage": "File uploaded and awaiting validation",
	"statusUpdateDate": "2023-09-04T09:21:38.000+00:00",
	"submitDate": "2023-09-04T09:21:38.000+00:00",
	"submitter": "nobody",
	"submitterEmail": "nobody@nomail-akamai.com"
}
`,
			expectedPath: "/hapi/v1/change-requests/123",
			expectedResponse: &ChangeRequest{
				Action:   "EDIT",
				ChangeID: 123,
				EdgeHostnames: []EdgeHostname{
					{
						EdgeHostnameID:    112233,
						RecordName:        "test123",
						DNSZone:           "edgekey.net",
						SecurityType:      "ENHANCED-TLS",
						UseDefaultTTL:     true,
						UseDefaultMap:     true,
						TTL:               21600,
						Map:               "a;bcd.akamaiedge.net",
						SlotNumber:        1234,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						Comments:          "",
						ChinaCDN: ChinaCDN{
							IsChinaCDN: false,
						},
						ProductId:    "DSA",
						SerialNumber: 0,
					},
				},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: "2023-09-04T09:21:38.000+00:00",
				SubmitDate:       "2023-09-04T09:21:38.000+00:00",
				Submitter:        "nobody",
				SubmitterEmail:   "nobody@nomail-akamai.com",
			},
		},
		"403 Access Denied to Edge Hostname": {
			request:        GetChangeRequest{234},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
	"type": "/hapi/problems/access-denied-to-edge-hostname",
	"title": "Access Denied to Edge Hostname",
	"status": 403,
	"detail": "You do not have access to this edge hostname",
	"instance": "/hapi/error-instances/7a2c1b84-fe90-40f2-8391-aaaaaaaaaa",
	"requestInstance": "http://cloud.akamaiapis.net/hapi/open/v1/change-requests/234#aaaaaaa",
	"method": "GET",
	"requestTime": "2023-09-04T09:26:16.561878149Z",
	"errors": []
}`,
			expectedPath: "/hapi/v1/change-requests/234",
			withError: &Error{
				Type:            "/hapi/problems/access-denied-to-edge-hostname",
				Title:           "Access Denied to Edge Hostname",
				Status:          403,
				Detail:          "You do not have access to this edge hostname",
				Instance:        "/hapi/error-instances/7a2c1b84-fe90-40f2-8391-aaaaaaaaaa",
				RequestInstance: "http://cloud.akamaiapis.net/hapi/open/v1/change-requests/234#aaaaaaa",
				Method:          "GET",
				RequestTime:     "2023-09-04T09:26:16.561878149Z",
			},
		},
		"500 internal server error": {
			request:        GetChangeRequest{123},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting activation",
    "status": 500
}`,
			expectedPath: "/hapi/v1/change-requests/123",
			withError: &Error{
				Type:   "internal_error",
				Title:  "Internal Server Error",
				Detail: "Error deleting activation",
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
			result, err := client.GetChangeRequest(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
