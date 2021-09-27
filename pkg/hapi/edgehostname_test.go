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

func TestDeleteEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		request          DeleteEdgeHostnameRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DeleteEdgeHostnameResponse
		withError        error
	}{
		"202 Accepted": {
			request: DeleteEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-001",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "action": "DELETE",
    "changeId": 66025603,
    "edgeHostnames": [
        {
            "chinaCdn": {
                "isChinaCdn": false
            },
            "dnsZone": "edgesuite.net",
            "edgeHostnameId": 4558392,
            "recordName": "mgw-test-001",
            "securityType": "STANDARD-TLS",
            "useDefaultMap": false,
            "useDefaultTtl": false
        }
    ],
    "status": "PENDING",
    "statusMessage": "File uploaded and awaiting validation",
    "statusUpdateDate": "2021-09-23T15:07:10.000+00:00",
    "submitDate": "2021-09-23T15:07:10.000+00:00",
    "submitter": "ftzgvvigljhoq5ib",
    "submitterEmail": "ftzgvvigljhoq5ib@nomail-akamai.com"
}`,
			expectedPath: "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/mgw-test-001?comments=some+comment&statusUpdateEmail=some%40example.com",
			expectedResponse: &DeleteEdgeHostnameResponse{
				Action:   "DELETE",
				ChangeID: 66025603,
				EdgeHostnames: []EdgeHostname{{
					ChinaCDN: ChinaCDN{
						IsChinaCDN: false,
					},
					DNSZone:        "edgesuite.net",
					EdgeHostnameID: 4558392,
					RecordName:     "mgw-test-001",
					SecurityType:   "STANDARD-TLS",
					UseDefaultMap:  false,
					UseDefaultTTL:  false,
				},
				},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: "2021-09-23T15:07:10.000+00:00",
				SubmitDate:       "2021-09-23T15:07:10.000+00:00",
				Submitter:        "ftzgvvigljhoq5ib",
				SubmitterEmail:   "ftzgvvigljhoq5ib@nomail-akamai.com",
			},
		},
		"404 could not find edge hostname": {
			request: DeleteEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-003",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/hapi/problems/record-name-dns-zone-not-found",
    "title": "Invalid Record Name/DNS Zone",
    "status": 404,
    "detail": "Could not find edge hostname with record name mgw-test-003 and DNS Zone edgesuite.net",
    "instance": "/hapi/error-instances/47f08d26-00b4-4c05-a8c0-bcbc542b9bce",
    "requestInstance": "http://cloud-qa-resource-impl.luna-dev.akamaiapis.net/hapi/open/v1/dns-zones/edgesuite.net/edge-hostnames/mgw-test-003#9ea9060c",
    "method": "DELETE",
    "requestTime": "2021-09-23T15:37:28.383173Z",
    "errors": [],
    "domainPrefix": "mgw-test-003",
    "domainSuffix": "edgesuite.net"
}`,
			expectedPath: "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/mgw-test-003?comments=some+comment&statusUpdateEmail=some%40example.com",
			withError: &Error{
				Type:            "/hapi/problems/record-name-dns-zone-not-found",
				Title:           "Invalid Record Name/DNS Zone",
				Status:          404,
				Detail:          "Could not find edge hostname with record name mgw-test-003 and DNS Zone edgesuite.net",
				Instance:        "/hapi/error-instances/47f08d26-00b4-4c05-a8c0-bcbc542b9bce",
				RequestInstance: "http://cloud-qa-resource-impl.luna-dev.akamaiapis.net/hapi/open/v1/dns-zones/edgesuite.net/edge-hostnames/mgw-test-003#9ea9060c",
				Method:          "DELETE",
				RequestTime:     "2021-09-23T15:37:28.383173Z",
				DomainPrefix:    "mgw-test-003",
				DomainSuffix:    "edgesuite.net",
			},
		},
		"500 internal server error": {
			request: DeleteEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-002",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting activation",
    "status": 500
}`,
			expectedPath: "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/mgw-test-002?comments=some+comment&statusUpdateEmail=some%40example.com",
			withError: &Error{
				Type:   "internal_error",
				Title:  "Internal Server Error",
				Detail: "Error deleting activation",
				Status: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: DeleteEdgeHostnameRequest{
				RecordName:        "atv_1696855",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
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
			result, err := client.DeleteEdgeHostname(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
