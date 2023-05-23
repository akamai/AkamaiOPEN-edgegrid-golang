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

func TestGetEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		edgeHostnameID   int
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetEdgeHostnameResponse
		withError        error
	}{
		"200 OK": {
			edgeHostnameID: 1234,
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"chinaCdn": {
					"isChinaCdn": false
				},
				"comments": "Created by Property-Manager/PAPI on Thu Mar 03 15:58:17 GMT 2022",
				"dnsZone": "edgekey.net",
				"edgeHostnameId": 4617960,
				"ipVersionBehavior": "IPV6_IPV4_DUALSTACK",
				"map": "e;dscx.akamaiedge.net",
				"recordName": "aws_ci_pearltest-asorigin-na-as-eu-ionp.cumulus-essl.webexp-ipqa-ion.com-v2",
				"securityType": "ENHANCED-TLS",
				"slotNumber": 47463,
				"ttl": 21600,
				"useDefaultMap": true,
				"useDefaultTtl": true
			}`,
			expectedPath: "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{
				ChinaCdn: ChinaCDN{
					IsChinaCDN: false,
				},
				Comments:          "Created by Property-Manager/PAPI on Thu Mar 03 15:58:17 GMT 2022",
				DNSZone:           "edgekey.net",
				EdgeHostnameID:    4617960,
				IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
				Map:               "e;dscx.akamaiedge.net",
				RecordName:        "aws_ci_pearltest-asorigin-na-as-eu-ionp.cumulus-essl.webexp-ipqa-ion.com-v2",
				SecurityType:      "ENHANCED-TLS",
				SlotNumber:        47463,
				TTL:               21600,
				UseDefaultMap:     true,
				UseDefaultTTL:     true,
			},
		},
		"404 could not find edge hostname": {
			edgeHostnameID: 9999,
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"type": "/hapi/problems/edge-hostname-not-found",
				"title": "Edge Hostname Not Found",
				"status": 404,
				"detail": "Edge hostname not found",
				"instance": "/hapi/error-instances/cdc47ffa-46f2-410d-8059-3f454c435e93",
				"requestInstance": "http://cloud-qa-resource-impl.luna-dev.akamaiapis.net/hapi/open/v1/edge-hostnames/9999#8a702528",
				"method": "GET",
				"requestTime": "2022-03-03T16:43:19.876613Z",
				"errors": []
			}`,
			expectedPath: "/hapi/v1/edge-hostnames/9999",
			withError: &Error{
				Type:            "/hapi/problems/edge-hostname-not-found",
				Title:           "Edge Hostname Not Found",
				Status:          404,
				Detail:          "Edge hostname not found",
				Instance:        "/hapi/error-instances/cdc47ffa-46f2-410d-8059-3f454c435e93",
				RequestInstance: "http://cloud-qa-resource-impl.luna-dev.akamaiapis.net/hapi/open/v1/edge-hostnames/9999#8a702528",
				Method:          "GET",
				RequestTime:     "2022-03-03T16:43:19.876613Z",
			},
		},
		"500 internal server error": {
			edgeHostnameID: 9999,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting activation",
				"status": 500
			}`,
			expectedPath: "/hapi/v1/edge-hostnames/9999",
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
			result, err := client.GetEdgeHostname(context.Background(), test.edgeHostnameID)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPatchEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		request          UpdateEdgeHostnameRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateEdgeHostnameResponse
		withError        error
	}{
		"202 Accepted": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-001",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "replace",
						Path:  "/ttl",
						Value: "10000",
					},
					{
						Op:    "replace",
						Path:  "/ipVersionBehavior",
						Value: "IPV4",
					},
				},
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "action": "EDIT",
    "changeId": 66025603,
    "edgeHostnames": [
        {
            "chinaCdn": {
                "isChinaCdn": false
            },
            "dnsZone": "edgesuite.net",
            "edgeHostnameId": 4558392,
			"ipVersionBehavior": "IPV4",
            "recordName": "mgw-test-001",
            "securityType": "STANDARD-TLS",
			"ttl": 10000,
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
			expectedResponse: &UpdateEdgeHostnameResponse{
				Action:   "EDIT",
				ChangeID: 66025603,
				EdgeHostnames: []EdgeHostname{{
					ChinaCDN: ChinaCDN{
						IsChinaCDN: false,
					},
					DNSZone:           "edgesuite.net",
					EdgeHostnameID:    4558392,
					RecordName:        "mgw-test-001",
					SecurityType:      "STANDARD-TLS",
					UseDefaultMap:     false,
					UseDefaultTTL:     false,
					TTL:               10000,
					IPVersionBehavior: "IPV4",
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
		"400 Incorrect body": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-001",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Path:  "/incorrect",
						Value: "some Value",
					},
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
    "type": "/hapi/problems/invalid-patch-request",
    "title": "Invalid Patch Request",
    "status": 400,
    "detail": "Invalid 'patch' request: patch replacement is only supported for 'TTL',and 'IpVersionBehavior'",
    "instance": "/hapi/error-instances/02702ac2-38a8-42a8-a482-07e1e4a93a44",
    "requestInstance": "http://cloud-qa-resource-impl.luna-dev.akamaiapis.net/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/mgw-test-001?comments=some+comment&statusUpdateEmail=some%40example.com#0e423b67",
    "method": "PATCH",
    "requestTime": "2022-05-23T13:50:06.221019Z",
    "errors": []
}`,
			expectedPath: "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/mgw-test-001?comments=some+comment&statusUpdateEmail=some%40example.com",
			withError:    ErrUpdateEdgeHostname,
		},
		"500 internal server error": {
			request: UpdateEdgeHostnameRequest{
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
			withError:    ErrUpdateEdgeHostname,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPatch, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateEdgeHostname(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
