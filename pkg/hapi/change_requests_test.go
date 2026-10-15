package hapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChangeRequest(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          GetChangeRequestRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetChangeRequestResponse
		withError        error
		cancelContext    bool
	}{
		"200 OK": {
			request:        GetChangeRequestRequest{ChangeID: 123},
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
			expectedResponse: &GetChangeRequestResponse{
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
						SlotNumber:        ptr.To(int64(1234)),
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						Comments:          "",
						ChinaCDN: ChinaCDN{
							IsChinaCDN: false,
						},
						ProductID:    "DSA",
						SerialNumber: 0,
					},
				},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: test.NewTimeFromString(t, "2023-09-04T09:21:38.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2023-09-04T09:21:38.000+00:00"),
				Submitter:        "nobody",
				SubmitterEmail:   "nobody@nomail-akamai.com",
			},
		},
		"200 OK with HTTPSServiceBinding": {
			request:        GetChangeRequestRequest{ChangeID: 456},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"action": "EDIT",
	"changeId": 456,
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
			"useDefaultTtl": true,
			"httpsServiceBinding": "H2_AND_H3"
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
			expectedPath: "/hapi/v1/change-requests/456",
			expectedResponse: &GetChangeRequestResponse{
				Action:   "EDIT",
				ChangeID: 456,
				EdgeHostnames: []EdgeHostname{
					{
						EdgeHostnameID:      112233,
						RecordName:          "test123",
						DNSZone:             "edgekey.net",
						SecurityType:        "ENHANCED-TLS",
						UseDefaultTTL:       true,
						UseDefaultMap:       true,
						TTL:                 21600,
						Map:                 "a;bcd.akamaiedge.net",
						SlotNumber:          ptr.To(int64(1234)),
						IPVersionBehavior:   "IPV6_IPV4_DUALSTACK",
						ChinaCDN:            ChinaCDN{IsChinaCDN: false},
						ProductID:           "DSA",
						SerialNumber:        0,
						HTTPSServiceBinding: ptr.To("H2_AND_H3"),
					},
				},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: test.NewTimeFromString(t, "2023-09-04T09:21:38.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2023-09-04T09:21:38.000+00:00"),
				Submitter:        "nobody",
				SubmitterEmail:   "nobody@nomail-akamai.com",
			},
		},
		"200 OK with query options and optimization status": {
			request: GetChangeRequestRequest{
				ChangeID:              123,
				IncludeDeliveryPolicy: ptr.To(false),
			},
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeId":123,"optimizationStatus":"PROCESSED"}`,
			expectedPath:     "/hapi/v1/change-requests/123?includeDeliveryPolicy=false",
			expectedResponse: &GetChangeRequestResponse{ChangeID: 123, OptimizationStatus: ptr.To("PROCESSED")},
		},
		"200 OK with full metadata and delivery policy option": {
			request:        GetChangeRequestRequest{ChangeID: 123, IncludeDeliveryPolicy: ptr.To(true)},
			responseStatus: http.StatusOK,
			responseBody: `{
				"action":"CREATE","changeId":123,"comments":"Created for example.com",
				"edgeHostnames":[{
					"recordName":"www.example.com","dnsZone":"edgesuite.net","securityType":"STANDARD-TLS",
					"productId":"DSA","useDefaultMap":true,"useDefaultTtl":true,
					"dnsSerial":1234567890,
					"chinaCdn":{"isChinaCdn":true,"chinaPartner":"example-partner","customChinaCdnMap":"partner.example.com","gtmDomain":"example.akadns.net"},
					"deliveryPolicy":{"policyType":"GEO_OPT_IN","policyAttributes":[{"attributeType":"GEO","code":"US","name":"USA"}]}
				}],
				"status":"SUCCEEDED","statusMessage":"Deployed",
				"statusUpdateDate":"2018-10-17T20:21:17.0Z","statusUpdateEmail":"status@example.com",
				"submitDate":"2018-10-17T20:21:12.0Z","submitter":"example","submitterEmail":"user@example.com"
			}`,
			expectedPath: "/hapi/v1/change-requests/123?includeDeliveryPolicy=true",
			expectedResponse: &GetChangeRequestResponse{
				Action: "CREATE", ChangeID: 123, Comments: "Created for example.com",
				EdgeHostnames: []EdgeHostname{{
					RecordName: "www.example.com", DNSZone: "edgesuite.net", SecurityType: "STANDARD-TLS",
					ProductID: "DSA", UseDefaultMap: true, UseDefaultTTL: true,
					DNSSerial: ptr.To(int64(1234567890)),
					ChinaCDN: ChinaCDN{
						IsChinaCDN: true, CustomChinaCDNMap: ptr.To("partner.example.com"),
						ChinaPartner: ptr.To("example-partner"), GTMDomain: ptr.To("example.akadns.net"),
					},
					DeliveryPolicy: &DeliveryPolicy{
						PolicyType:       "GEO_OPT_IN",
						PolicyAttributes: []DeliveryPolicyAttribute{{AttributeType: ptr.To("GEO"), Code: "US", Name: ptr.To("USA")}},
					},
				}},
				Status: "SUCCEEDED", StatusMessage: "Deployed",
				StatusUpdateDate: test.NewTimeFromString(t, "2018-10-17T20:21:17.0Z"), StatusUpdateEmail: ptr.To("status@example.com"),
				SubmitDate: test.NewTimeFromString(t, "2018-10-17T20:21:12.0Z"), Submitter: "example", SubmitterEmail: "user@example.com",
			},
		},
		"200 OK with null policy and optional policy attributes": {
			request:        GetChangeRequestRequest{ChangeID: 123, IncludeDeliveryPolicy: ptr.To(true)},
			responseStatus: http.StatusOK,
			responseBody: `{
				"changeId":123,"edgeHostnames":[
					{"deliveryPolicy":null},
					{"deliveryPolicy":{"policyType":"GEO_OPT_IN","policyAttributes":[
						{"code":"US"},{"code":"IN","attributeType":null,"name":null},{"code":"EU_CH_GB","attributeType":"","name":""}
					]}}
				]
			}`,
			expectedPath: "/hapi/v1/change-requests/123?includeDeliveryPolicy=true",
			expectedResponse: &GetChangeRequestResponse{
				ChangeID: 123,
				EdgeHostnames: []EdgeHostname{
					{},
					{DeliveryPolicy: &DeliveryPolicy{
						PolicyType: "GEO_OPT_IN",
						PolicyAttributes: []DeliveryPolicyAttribute{
							{Code: "US"},
							{Code: "IN"},
							{Code: "EU_CH_GB", AttributeType: ptr.To(""), Name: ptr.To("")},
						},
					}},
				},
			},
		},
		"200 OK with null optimization status": {
			request:          GetChangeRequestRequest{ChangeID: 123},
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeId":123,"optimizationStatus":null}`,
			expectedPath:     "/hapi/v1/change-requests/123",
			expectedResponse: &GetChangeRequestResponse{ChangeID: 123},
		},
		"200 OK with empty optimization status": {
			request:          GetChangeRequestRequest{ChangeID: 123},
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeId":123,"optimizationStatus":""}`,
			expectedPath:     "/hapi/v1/change-requests/123",
			expectedResponse: &GetChangeRequestResponse{ChangeID: 123, OptimizationStatus: ptr.To("")},
		},
		"200 OK with 64-bit IDs and timestamp offsets": {
			request:        GetChangeRequestRequest{ChangeID: 2147483648},
			responseStatus: http.StatusOK,
			responseBody:   `{"changeId":2147483648,"edgeHostnames":[{"edgeHostnameId":2147483648,"ttl":2147483648,"slotNumber":2147483648,"serialNumber":2147483648}],"submitDate":"2023-09-04T14:51:38.123456789+05:30","statusUpdateDate":"2023-09-04T09:21:38.123456789Z"}`,
			expectedPath:   "/hapi/v1/change-requests/2147483648",
			expectedResponse: &GetChangeRequestResponse{
				ChangeID:         2147483648,
				EdgeHostnames:    []EdgeHostname{{EdgeHostnameID: 2147483648, TTL: 2147483648, SlotNumber: ptr.To(int64(2147483648)), SerialNumber: 2147483648}},
				SubmitDate:       test.NewTimeFromString(t, "2023-09-04T14:51:38.123456789+05:30"),
				StatusUpdateDate: test.NewTimeFromString(t, "2023-09-04T09:21:38.123456789Z"),
			},
		},
		"200 OK with null timestamps": {
			request:          GetChangeRequestRequest{ChangeID: 123},
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeId":123,"submitDate":null,"statusUpdateDate":null}`,
			expectedPath:     "/hapi/v1/change-requests/123",
			expectedResponse: &GetChangeRequestResponse{ChangeID: 123},
		},
		"200 OK with standard TLS and omitted optional fields": {
			request:        GetChangeRequestRequest{ChangeID: 123},
			responseStatus: http.StatusOK,
			responseBody: `{
				"action":"CREATE","changeId":123,"comments":"Created by Property-Manager/PAPI",
				"edgeHostnames":[{
					"chinaCdn":{"isChinaCdn":false},"comments":"Created by Property-Manager/PAPI",
					"dnsZone":"edgesuite.net","edgeHostnameId":1234,"ipVersionBehavior":"IPV4",
					"map":"a;r.akamai.net","productId":"IONNA","recordName":"example.com",
					"securityType":"STANDARD-TLS","dnsSerial":774,"serialNumber":774,"ttl":21600,
					"useDefaultMap":true,"useDefaultTtl":true
				}],
				"status":"SUCCEEDED","statusMessage":"File successfully deployed to Akamai's network",
				"statusUpdateDate":"2026-09-22T06:48:21.000+00:00","submitDate":"2026-09-22T06:41:52.000+00:00",
				"submitter":"example","submitterEmail":"user@example.com"
			}`,
			expectedPath: "/hapi/v1/change-requests/123",
			expectedResponse: &GetChangeRequestResponse{
				Action: "CREATE", ChangeID: 123, Comments: "Created by Property-Manager/PAPI",
				EdgeHostnames: []EdgeHostname{{
					Comments: "Created by Property-Manager/PAPI", DNSZone: "edgesuite.net", EdgeHostnameID: 1234,
					IPVersionBehavior: "IPV4", Map: "a;r.akamai.net", ProductID: "IONNA", RecordName: "example.com",
					SecurityType: "STANDARD-TLS", DNSSerial: ptr.To(int64(774)), SerialNumber: 774, TTL: 21600,
					UseDefaultMap: true, UseDefaultTTL: true,
				}},
				Status: "SUCCEEDED", StatusMessage: "File successfully deployed to Akamai's network",
				StatusUpdateDate: test.NewTimeFromString(t, "2026-09-22T06:48:21.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2026-09-22T06:41:52.000+00:00"),
				Submitter:        "example", SubmitterEmail: "user@example.com",
			},
		},
		"missing change ID": {
			withError: ErrStructValidation,
		},
		"invalid timestamp": {
			request:        GetChangeRequestRequest{ChangeID: 123},
			responseStatus: http.StatusOK,
			responseBody:   `{"changeId":123,"submitDate":"invalid"}`,
			expectedPath:   "/hapi/v1/change-requests/123",
			withError:      ErrGetChangeRequest,
		},
		"404 change request not found": {
			request:        GetChangeRequestRequest{ChangeID: 123},
			responseStatus: http.StatusNotFound,
			responseBody:   `{"type":"/hapi/problems/change-id-not-found","title":"Invalid Change ID","changeId":"123","status":404}`,
			expectedPath:   "/hapi/v1/change-requests/123",
			withError:      &Error{Type: "/hapi/problems/change-id-not-found", Title: "Invalid Change ID", ChangeID: "123", Status: http.StatusNotFound},
		},
		"403 Access Denied to Edge Hostname": {
			request:        GetChangeRequestRequest{ChangeID: 234},
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
				RequestTime:     test.NewTimeFromString(t, "2023-09-04T09:26:16.561878149Z"),
			},
		},
		"500 internal server error": {
			request:        GetChangeRequestRequest{ChangeID: 123},
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
		"canceled context": {
			request:       GetChangeRequestRequest{ChangeID: 123},
			cancelContext: true,
			withError:     context.Canceled,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !assert.NotEmpty(t, test.expectedPath, "unexpected HTTP request") {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			ctx := context.Background()
			if test.cancelContext {
				canceledCtx, cancel := context.WithCancel(ctx)
				cancel()
				ctx = canceledCtx
			}
			result, err := client.GetChangeRequest(ctx, test.request)
			if test.withError != nil {
				assert.ErrorIs(t, err, ErrGetChangeRequest)
				assert.Nil(t, result)
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				if errors.Is(test.withError, ErrStructValidation) {
					assert.ErrorContains(t, err, "ChangeID: cannot be blank")
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListChangeRequests(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          ListChangeRequestsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListChangeRequestsResponse
		withError        error
		errorContains    string
		cancelContext    bool
	}{
		"200 OK with defaults": {
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeRequests":[{"changeId":123,"status":"PENDING"}]}`,
			expectedPath:     "/hapi/v1/change-requests",
			expectedResponse: &ListChangeRequestsResponse{ChangeRequests: []ChangeRequest{{ChangeID: 123, Status: "PENDING"}}},
		},
		"200 OK with all query options": {
			request:        ListChangeRequestsRequest{Status: ChangeRequestStatusPending, IncludeDeliveryPolicy: ptr.To(true)},
			responseStatus: http.StatusOK,
			responseBody: `{"changeRequests":[{"changeId":123,"optimizationStatus":"PROCESSING","edgeHostnames":[{
				"recordName":"www.example.com","dnsZone":"edgesuite.net",
				"deliveryPolicy":{"policyType":"GEO_OPT_IN","policyAttributes":[{"attributeType":"GEO","code":"US","name":"USA"}]}
			}]}]}`,
			expectedPath: "/hapi/v1/change-requests?includeDeliveryPolicy=true&status=PENDING",
			expectedResponse: &ListChangeRequestsResponse{ChangeRequests: []ChangeRequest{{
				ChangeID: 123, OptimizationStatus: ptr.To("PROCESSING"),
				EdgeHostnames: []EdgeHostname{{
					RecordName: "www.example.com", DNSZone: "edgesuite.net",
					DeliveryPolicy: &DeliveryPolicy{
						PolicyType:       "GEO_OPT_IN",
						PolicyAttributes: []DeliveryPolicyAttribute{{AttributeType: ptr.To("GEO"), Code: "US", Name: ptr.To("USA")}},
					},
				}},
			}}},
		},
		"200 OK with explicit false delivery policy option and empty list": {
			request:          ListChangeRequestsRequest{IncludeDeliveryPolicy: ptr.To(false)},
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeRequests":[]}`,
			expectedPath:     "/hapi/v1/change-requests?includeDeliveryPolicy=false",
			expectedResponse: &ListChangeRequestsResponse{ChangeRequests: []ChangeRequest{}},
		},
		"invalid status": {
			request:       ListChangeRequestsRequest{Status: "INVALID"},
			withError:     ErrStructValidation,
			errorContains: "Status: must be a valid value",
		},
		"400 API error": {
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"type":"/hapi/problems/invalid-status-name","title":"Invalid Status Name","status":400}`,
			expectedPath:   "/hapi/v1/change-requests",
			withError:      &Error{Type: "/hapi/problems/invalid-status-name", Title: "Invalid Status Name", Status: http.StatusBadRequest},
		},
		"500 API error": {
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"title":"Internal Server Error","status":500}`,
			expectedPath:   "/hapi/v1/change-requests",
			withError:      &Error{Title: "Internal Server Error", Status: http.StatusInternalServerError},
		},
		"malformed success response": {
			responseStatus: http.StatusOK,
			responseBody:   `{"changeRequests":`,
			expectedPath:   "/hapi/v1/change-requests",
			withError:      session.ErrUnmarshaling,
		},
		"canceled context": {
			cancelContext: true,
			withError:     context.Canceled,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !assert.NotEmpty(t, test.expectedPath, "unexpected HTTP request") {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			ctx := context.Background()
			if test.cancelContext {
				canceledCtx, cancel := context.WithCancel(ctx)
				cancel()
				ctx = canceledCtx
			}
			result, err := client.ListChangeRequests(ctx, test.request)
			if test.withError != nil {
				assert.Nil(t, result)
				assert.ErrorIs(t, err, ErrListChangeRequests)
				assert.ErrorIs(t, err, test.withError)
				if test.errorContains != "" {
					assert.ErrorContains(t, err, test.errorContains)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListChangeRequestsPerEdgeHostname(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          ListChangeRequestsPerEdgeHostnameRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListChangeRequestsPerEdgeHostnameResponse
		withError        error
		errorContains    string
		cancelContext    bool
	}{
		"200 OK with defaults": {
			request:          ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "edgesuite.net", RecordName: "www.example.com"},
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeRequests":[{"changeId":123,"status":"SUCCEEDED"},{"changeId":124,"status":"PENDING"}]}`,
			expectedPath:     "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com/change-requests",
			expectedResponse: &ListChangeRequestsPerEdgeHostnameResponse{ChangeRequests: []ChangeRequest{{ChangeID: 123, Status: "SUCCEEDED"}, {ChangeID: 124, Status: "PENDING"}}},
		},
		"200 OK with all query options": {
			request:        ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "edgesuite.net", RecordName: "www.example.com", IncludeDeliveryPolicy: ptr.To(true)},
			responseStatus: http.StatusOK,
			responseBody: `{"changeRequests":[{"changeId":123,"optimizationStatus":"PROCESSED","edgeHostnames":[{
				"recordName":"www.example.com","dnsZone":"edgesuite.net",
				"deliveryPolicy":{"policyType":"GEO_OPT_IN","policyAttributes":[{"attributeType":"GEO","code":"US","name":"USA"}]}
			}]}]}`,
			expectedPath: "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com/change-requests?includeDeliveryPolicy=true",
			expectedResponse: &ListChangeRequestsPerEdgeHostnameResponse{ChangeRequests: []ChangeRequest{{
				ChangeID: 123, OptimizationStatus: ptr.To("PROCESSED"),
				EdgeHostnames: []EdgeHostname{{
					RecordName: "www.example.com", DNSZone: "edgesuite.net",
					DeliveryPolicy: &DeliveryPolicy{
						PolicyType:       "GEO_OPT_IN",
						PolicyAttributes: []DeliveryPolicyAttribute{{AttributeType: ptr.To("GEO"), Code: "US", Name: ptr.To("USA")}},
					},
				}},
			}}},
		},
		"200 OK with explicit false delivery policy option and empty list": {
			request:          ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "edgesuite.net", RecordName: "www.example.com", IncludeDeliveryPolicy: ptr.To(false)},
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeRequests":[]}`,
			expectedPath:     "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com/change-requests?includeDeliveryPolicy=false",
			expectedResponse: &ListChangeRequestsPerEdgeHostnameResponse{ChangeRequests: []ChangeRequest{}},
		},
		"200 OK with escaped path segments": {
			request:          ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "zone/name", RecordName: "record?name"},
			responseStatus:   http.StatusOK,
			responseBody:     `{"changeRequests":[]}`,
			expectedPath:     "/hapi/v1/dns-zones/zone%2Fname/edge-hostnames/record%3Fname/change-requests",
			expectedResponse: &ListChangeRequestsPerEdgeHostnameResponse{ChangeRequests: []ChangeRequest{}},
		},
		"missing path fields": {
			withError:     ErrStructValidation,
			errorContains: "DNSZone: cannot be blank\nRecordName: cannot be blank",
		},
		"missing DNS zone": {
			request:       ListChangeRequestsPerEdgeHostnameRequest{RecordName: "www.example.com"},
			withError:     ErrStructValidation,
			errorContains: "DNSZone: cannot be blank",
		},
		"missing record name": {
			request:       ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "edgesuite.net"},
			withError:     ErrStructValidation,
			errorContains: "RecordName: cannot be blank",
		},
		"400 API error": {
			request:        ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "edgesuite.net", RecordName: "www.example.com"},
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"type":"INVALID_DOMAIN_PREFIX_SUFFIX","title":"Invalid Record Name/DNS Zone","status":0}`,
			expectedPath:   "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com/change-requests",
			withError:      &Error{Type: "INVALID_DOMAIN_PREFIX_SUFFIX", Title: "Invalid Record Name/DNS Zone", Status: http.StatusBadRequest},
		},
		"500 API error": {
			request:        ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "edgesuite.net", RecordName: "www.example.com"},
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"title":"Internal Server Error","status":500}`,
			expectedPath:   "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com/change-requests",
			withError:      &Error{Title: "Internal Server Error", Status: http.StatusInternalServerError},
		},
		"malformed success response": {
			request:        ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "edgesuite.net", RecordName: "www.example.com"},
			responseStatus: http.StatusOK,
			responseBody:   `{"changeRequests":`,
			expectedPath:   "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com/change-requests",
			withError:      session.ErrUnmarshaling,
		},
		"canceled context": {
			request:       ListChangeRequestsPerEdgeHostnameRequest{DNSZone: "edgesuite.net", RecordName: "www.example.com"},
			cancelContext: true,
			withError:     context.Canceled,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !assert.NotEmpty(t, test.expectedPath, "unexpected HTTP request") {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			ctx := context.Background()
			if test.cancelContext {
				canceledCtx, cancel := context.WithCancel(ctx)
				cancel()
				ctx = canceledCtx
			}
			result, err := client.ListChangeRequestsPerEdgeHostname(ctx, test.request)
			if test.withError != nil {
				assert.Nil(t, result)
				assert.ErrorIs(t, err, ErrListChangeRequestsPerEdgeHostname)
				assert.ErrorIs(t, err, test.withError)
				if test.errorContains != "" {
					assert.ErrorContains(t, err, test.errorContains)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
