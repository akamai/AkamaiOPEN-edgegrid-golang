package hapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
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
		"202 Accepted with multiple notification emails": {
			request: DeleteEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "www.example.com",
				StatusUpdateEmail: []string{"first@example.com", "second@example.com"},
				Comments:          "remove & notify",
			},
			responseStatus:   http.StatusAccepted,
			responseBody:     `{"changeId":123}`,
			expectedPath:     "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com?comments=remove+%26+notify&statusUpdateEmail=first%40example.com%2Csecond%40example.com",
			expectedResponse: &DeleteEdgeHostnameResponse{ChangeID: 123},
		},
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
				StatusUpdateDate: test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				Submitter:        "ftzgvvigljhoq5ib",
				SubmitterEmail:   "ftzgvvigljhoq5ib@nomail-akamai.com",
			},
		},
		"202 Accepted with HTTPSServiceBinding in response": {
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
    "changeId": 66025604,
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
            "useDefaultTtl": false,
            "httpsServiceBinding": "H2"
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
				ChangeID: 66025604,
				EdgeHostnames: []EdgeHostname{{
					ChinaCDN: ChinaCDN{
						IsChinaCDN: false,
					},
					DNSZone:             "edgesuite.net",
					EdgeHostnameID:      4558392,
					RecordName:          "mgw-test-001",
					SecurityType:        "STANDARD-TLS",
					UseDefaultMap:       false,
					UseDefaultTTL:       false,
					HTTPSServiceBinding: ptr.To("H2"),
				}},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
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
				RequestTime:     test.NewTimeFromString(t, "2021-09-23T15:37:28.383173Z"),
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
			defer mockServer.Close()
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
	t.Parallel()
	tests := map[string]struct {
		edgeHostnameID        int64
		includeDeliveryPolicy *bool
		responseStatus        int
		responseBody          string
		expectedPath          string
		expectedResponse      *GetEdgeHostnameResponse
		withError             error
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
				"productId": "DSA",
				"map": "e;dscx.akamaiedge.net",
				"recordName": "aws_ci_pearltest-asorigin-na-as-eu-ionp.cumulus-essl.webexp-ipqa-ion.com-v2",
				"securityType": "ENHANCED-TLS",
				"slotNumber": 47463,
				"ttl": 21600,
				"useDefaultMap": true,
				"useDefaultTtl": true,
				"mapAlias": "al",
				"useCases": [
					{
						"type": "GLOBAL",
						"option": "LIVE",
						"useCase": "Segmented_Media_Mode"
					}
				]
			}`,
			expectedPath: "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{
				ChinaCDN: ChinaCDN{
					IsChinaCDN: false,
				},
				Comments:          "Created by Property-Manager/PAPI on Thu Mar 03 15:58:17 GMT 2022",
				DNSZone:           "edgekey.net",
				EdgeHostnameID:    4617960,
				IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
				ProductID:         "DSA",
				Map:               "e;dscx.akamaiedge.net",
				RecordName:        "aws_ci_pearltest-asorigin-na-as-eu-ionp.cumulus-essl.webexp-ipqa-ion.com-v2",
				SecurityType:      "ENHANCED-TLS",
				SlotNumber:        ptr.To(int64(47463)),
				TTL:               21600,
				UseDefaultMap:     true,
				UseDefaultTTL:     true,
				MapAlias:          ptr.To("al"),
				UseCases: []UseCase{
					{
						Type:    "GLOBAL",
						Option:  "LIVE",
						UseCase: "Segmented_Media_Mode",
					},
				},
			},
		},
		"200 OK with HTTPSServiceBinding": {
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
				"productId": "DSA",
				"map": "e;dscx.akamaiedge.net",
				"recordName": "aws_ci_pearltest-asorigin-na-as-eu-ionp.cumulus-essl.webexp-ipqa-ion.com-v2",
				"securityType": "ENHANCED-TLS",
				"slotNumber": 47463,
				"ttl": 21600,
				"useDefaultMap": true,
				"useDefaultTtl": true,
				"mapAlias": "al",
				"httpsServiceBinding": "H2",
				"useCases": []
			}`,
			expectedPath: "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{
				ChinaCDN: ChinaCDN{
					IsChinaCDN: false,
				},
				Comments:            "Created by Property-Manager/PAPI on Thu Mar 03 15:58:17 GMT 2022",
				DNSZone:             "edgekey.net",
				EdgeHostnameID:      4617960,
				IPVersionBehavior:   "IPV6_IPV4_DUALSTACK",
				ProductID:           "DSA",
				Map:                 "e;dscx.akamaiedge.net",
				RecordName:          "aws_ci_pearltest-asorigin-na-as-eu-ionp.cumulus-essl.webexp-ipqa-ion.com-v2",
				SecurityType:        "ENHANCED-TLS",
				SlotNumber:          ptr.To(int64(47463)),
				TTL:                 21600,
				UseDefaultMap:       true,
				UseDefaultTTL:       true,
				MapAlias:            ptr.To("al"),
				HTTPSServiceBinding: ptr.To("H2"),
				UseCases:            []UseCase{},
			},
		},
		"200 OK with hostname metadata and delivery policy": {
			edgeHostnameID: 1234,
			responseStatus: http.StatusOK,
			responseBody: `{
				"recordName":"www.example.com","dnsZone":"edgesuite.net","securityType":"STANDARD-TLS",
				"productId":"DSA","useDefaultMap":true,"useDefaultTtl":true,
				"dnsSerial":1234567890,
				"chinaCdn":{"isChinaCdn":true,"chinaPartner":"example-partner","customChinaCdnMap":"partner.example.com","gtmDomain":"example.akadns.net"},
				"deliveryPolicy":{"policyType":"GEO_OPT_IN","policyAttributes":[{"attributeType":"GEO","code":"US","name":"USA"}]}
			}`,
			expectedPath: "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{
				RecordName: "www.example.com", DNSZone: "edgesuite.net", SecurityType: "STANDARD-TLS",
				ProductID: "DSA", UseDefaultMap: true, UseDefaultTTL: true, DNSSerial: ptr.To(int64(1234567890)),
				ChinaCDN: ChinaCDN{
					IsChinaCDN: true, CustomChinaCDNMap: ptr.To("partner.example.com"),
					ChinaPartner: ptr.To("example-partner"), GTMDomain: ptr.To("example.akadns.net"),
				},
				DeliveryPolicy: &DeliveryPolicy{
					PolicyType:       "GEO_OPT_IN",
					PolicyAttributes: []DeliveryPolicyAttribute{{AttributeType: ptr.To("GEO"), Code: "US", Name: ptr.To("USA")}},
				},
			},
		},
		"200 OK with null delivery policy": {
			edgeHostnameID:   1234,
			responseStatus:   http.StatusOK,
			responseBody:     `{"deliveryPolicy":null}`,
			expectedPath:     "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{},
		},
		"200 OK with empty and null hostname metadata": {
			edgeHostnameID:   1234,
			responseStatus:   http.StatusOK,
			responseBody:     `{"dnsSerial":0,"chinaCdn":{"chinaPartner":"","gtmDomain":null}}`,
			expectedPath:     "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{DNSSerial: ptr.To(int64(0)), ChinaCDN: ChinaCDN{ChinaPartner: ptr.To("")}},
		},
		"200 OK with delivery policy explicitly included": {
			edgeHostnameID:        1234,
			includeDeliveryPolicy: ptr.To(true),
			responseStatus:        http.StatusOK,
			responseBody:          `{"deliveryPolicy":{"policyType":"GEO_OPT_IN","policyAttributes":[{"code":"US"}]}}`,
			expectedPath:          "/hapi/v1/edge-hostnames/1234?includeDeliveryPolicy=true",
			expectedResponse: &GetEdgeHostnameResponse{
				DeliveryPolicy: &DeliveryPolicy{PolicyType: "GEO_OPT_IN", PolicyAttributes: []DeliveryPolicyAttribute{{Code: "US"}}},
			},
		},
		"200 OK with delivery policy explicitly excluded": {
			edgeHostnameID:        1234,
			includeDeliveryPolicy: ptr.To(false),
			responseStatus:        http.StatusOK,
			responseBody:          `{}`,
			expectedPath:          "/hapi/v1/edge-hostnames/1234?includeDeliveryPolicy=false",
			expectedResponse:      &GetEdgeHostnameResponse{},
		},
		"200 OK with 64-bit hostname values": {
			edgeHostnameID:   2147483648,
			responseStatus:   http.StatusOK,
			responseBody:     `{"edgeHostnameId":2147483648,"ttl":2147483648,"slotNumber":2147483648,"serialNumber":2147483648}`,
			expectedPath:     "/hapi/v1/edge-hostnames/2147483648",
			expectedResponse: &GetEdgeHostnameResponse{EdgeHostnameID: 2147483648, TTL: 2147483648, SlotNumber: ptr.To(int64(2147483648)), SerialNumber: 2147483648},
		},
		"200 OK with standard TLS and omitted optional fields": {
			edgeHostnameID: 1234,
			responseStatus: http.StatusOK,
			responseBody: `{
				"chinaCdn":{"isChinaCdn":false},
				"comments":"Created by Property-Manager/PAPI",
				"dnsZone":"edgesuite.net","edgeHostnameId":1234,
				"ipVersionBehavior":"IPV4","map":"a;r.akamai.net","productId":"IONNA",
				"recordName":"example.com","securityType":"STANDARD-TLS",
				"dnsSerial":1153,"serialNumber":1153,"ttl":21600,
				"useDefaultMap":true,"useDefaultTtl":true
			}`,
			expectedPath: "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{
				Comments: "Created by Property-Manager/PAPI", DNSZone: "edgesuite.net", EdgeHostnameID: 1234,
				IPVersionBehavior: "IPV4", Map: "a;r.akamai.net", ProductID: "IONNA",
				RecordName: "example.com", SecurityType: "STANDARD-TLS",
				DNSSerial: ptr.To(int64(1153)), SerialNumber: 1153, TTL: 21600,
				UseDefaultMap: true, UseDefaultTTL: true,
			},
		},
		"200 OK with explicitly zero-valued optional fields": {
			edgeHostnameID: 1234,
			responseStatus: http.StatusOK,
			responseBody: `{
				"slotNumber":0,"serialNumber":0,"dnsSerial":0,
				"customTarget":"","isEdgeIPBindingEnabled":false,"mapAlias":"",
				"chinaCdn":{"isChinaCdn":false,"customChinaCdnMap":"","chinaPartner":"","gtmDomain":""}
			}`,
			expectedPath: "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{
				SlotNumber: ptr.To(int64(0)), SerialNumber: 0, DNSSerial: ptr.To(int64(0)),
				CustomTarget: ptr.To(""), IsEdgeIPBindingEnabled: ptr.To(false), MapAlias: ptr.To(""),
				ChinaCDN: ChinaCDN{CustomChinaCDNMap: ptr.To(""), ChinaPartner: ptr.To(""), GTMDomain: ptr.To("")},
			},
		},
		"200 OK with null optional fields": {
			edgeHostnameID: 1234,
			responseStatus: http.StatusOK,
			responseBody: `{
				"slotNumber":null,"dnsSerial":null,
				"customTarget":null,"isEdgeIPBindingEnabled":null,"mapAlias":null,
				"chinaCdn":{"customChinaCdnMap":null,"chinaPartner":null,"gtmDomain":null}
			}`,
			expectedPath:     "/hapi/v1/edge-hostnames/1234",
			expectedResponse: &GetEdgeHostnameResponse{},
		},
		"missing edge hostname ID": {
			withError: ErrStructValidation,
		},
		"invalid DNS serial type": {
			edgeHostnameID: 1234,
			responseStatus: http.StatusOK,
			responseBody:   `{"dnsSerial":"1153"}`,
			expectedPath:   "/hapi/v1/edge-hostnames/1234",
			withError:      ErrGetEdgeHostname,
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
				RequestTime:     test.NewTimeFromString(t, "2022-03-03T16:43:19.876613Z"),
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
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetEdgeHostname(context.Background(), GetEdgeHostnameRequest{
				EdgeHostnameID: test.edgeHostnameID, IncludeDeliveryPolicy: test.includeDeliveryPolicy,
			})
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				if errors.Is(test.withError, ErrStructValidation) {
					assert.ErrorContains(t, err, "EdgeHostnameID: cannot be blank")
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateEdgeHostname(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          UpdateEdgeHostnameRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedBody     string
		expectedResponse *UpdateEdgeHostnameResponse
		withError        error
	}{
		"202 Accepted with multiple notification emails": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "www.example.com",
				StatusUpdateEmail: []string{"first@example.com", "second@example.com"},
				Comments:          "update & notify",
				Body:              []UpdateEdgeHostnameRequestBody{{Op: "replace", Path: "/ttl", Value: "10000"}},
			},
			responseStatus:   http.StatusAccepted,
			responseBody:     `{"changeId":123}`,
			expectedPath:     "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com?comments=update+%26+notify&statusUpdateEmail=first%40example.com%2Csecond%40example.com",
			expectedBody:     `[{"op":"replace","path":"/ttl","value":"10000"}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{ChangeID: 123},
		},
		"202 Accepted with optimization status": {
			request: UpdateEdgeHostnameRequest{
				DNSZone: "edgesuite.net", RecordName: "www.example.com",
				Body: []UpdateEdgeHostnameRequestBody{{Op: "replace", Path: "/ttl", Value: "10000"}},
			},
			responseStatus:   http.StatusAccepted,
			responseBody:     `{"changeId":123,"optimizationStatus":"PROCESSING"}`,
			expectedPath:     "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com",
			expectedBody:     `[{"op":"replace","path":"/ttl","value":"10000"}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{ChangeID: 123, OptimizationStatus: ptr.To("PROCESSING")},
		},
		"202 Accepted with empty optimization status": {
			request: UpdateEdgeHostnameRequest{
				DNSZone: "edgesuite.net", RecordName: "www.example.com",
				Body: []UpdateEdgeHostnameRequestBody{{Op: "replace", Path: "/ttl", Value: "10000"}},
			},
			responseStatus:   http.StatusAccepted,
			responseBody:     `{"changeId":123,"optimizationStatus":""}`,
			expectedPath:     "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com",
			expectedBody:     `[{"op":"replace","path":"/ttl","value":"10000"}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{ChangeID: 123, OptimizationStatus: ptr.To("")},
		},
		"202 Accepted with null optimization status": {
			request: UpdateEdgeHostnameRequest{
				DNSZone: "edgesuite.net", RecordName: "www.example.com",
				Body: []UpdateEdgeHostnameRequestBody{{Op: "replace", Path: "/ttl", Value: "10000"}},
			},
			responseStatus:   http.StatusAccepted,
			responseBody:     `{"changeId":123,"optimizationStatus":null}`,
			expectedPath:     "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com",
			expectedBody:     `[{"op":"replace","path":"/ttl","value":"10000"}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{ChangeID: 123},
		},
		"202 Accepted with 64-bit change ID and null timestamps": {
			request: UpdateEdgeHostnameRequest{
				DNSZone: "edgesuite.net", RecordName: "www.example.com",
				Body: []UpdateEdgeHostnameRequestBody{{Op: "replace", Path: "/ttl", Value: "10000"}},
			},
			responseStatus:   http.StatusAccepted,
			responseBody:     `{"changeId":2147483648,"submitDate":null,"statusUpdateDate":null}`,
			expectedPath:     "/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.example.com",
			expectedBody:     `[{"op":"replace","path":"/ttl","value":"10000"}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{ChangeID: 2147483648},
		},
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
			expectedBody: `[
				{"op":"replace","path":"/ttl","value":"10000"},
				{"op":"replace","path":"/ipVersionBehavior","value":"IPV4"}
			]`,
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
				StatusUpdateDate: test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				Submitter:        "ftzgvvigljhoq5ib",
				SubmitterEmail:   "ftzgvvigljhoq5ib@nomail-akamai.com",
			},
		},
		"202 Accepted with HTTPSServiceBinding in response": {
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
				},
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "action": "EDIT",
    "changeId": 66025604,
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
            "useDefaultTtl": false,
            "httpsServiceBinding": "H2"
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
			expectedBody: `[{"op":"replace","path":"/ttl","value":"10000"}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{
				Action:   "EDIT",
				ChangeID: 66025604,
				EdgeHostnames: []EdgeHostname{{
					ChinaCDN: ChinaCDN{
						IsChinaCDN: false,
					},
					DNSZone:             "edgesuite.net",
					EdgeHostnameID:      4558392,
					RecordName:          "mgw-test-001",
					SecurityType:        "STANDARD-TLS",
					UseDefaultMap:       false,
					UseDefaultTTL:       false,
					TTL:                 10000,
					IPVersionBehavior:   "IPV4",
					HTTPSServiceBinding: ptr.To("H2"),
				}},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				Submitter:        "ftzgvvigljhoq5ib",
				SubmitterEmail:   "ftzgvvigljhoq5ib@nomail-akamai.com",
			},
		},
		"202 Accepted - updating httpsServiceBinding": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-001",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "replace",
						Path:  "/httpsServiceBinding",
						Value: "H2_AND_H3",
					},
				},
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "action": "EDIT",
    "changeId": 66025605,
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
            "useDefaultTtl": false,
            "httpsServiceBinding": "H2_AND_H3"
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
			expectedBody: `[{"op":"replace","path":"/httpsServiceBinding","value":"H2_AND_H3"}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{
				Action:   "EDIT",
				ChangeID: 66025605,
				EdgeHostnames: []EdgeHostname{{
					ChinaCDN: ChinaCDN{
						IsChinaCDN: false,
					},
					DNSZone:             "edgesuite.net",
					EdgeHostnameID:      4558392,
					RecordName:          "mgw-test-001",
					SecurityType:        "STANDARD-TLS",
					UseDefaultMap:       false,
					UseDefaultTTL:       false,
					HTTPSServiceBinding: ptr.To("H2_AND_H3"),
				}},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				Submitter:        "ftzgvvigljhoq5ib",
				SubmitterEmail:   "ftzgvvigljhoq5ib@nomail-akamai.com",
			},
		},
		"202 Accepted - adding httpsServiceBinding with op add": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-001",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "add",
						Path:  "/httpsServiceBinding",
						Value: "H3",
					},
				},
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "action": "EDIT",
    "changeId": 66025606,
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
            "useDefaultTtl": false,
            "httpsServiceBinding": "H3"
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
			expectedBody: `[{"op":"add","path":"/httpsServiceBinding","value":"H3"}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{
				Action:   "EDIT",
				ChangeID: 66025606,
				EdgeHostnames: []EdgeHostname{{
					ChinaCDN: ChinaCDN{
						IsChinaCDN: false,
					},
					DNSZone:             "edgesuite.net",
					EdgeHostnameID:      4558392,
					RecordName:          "mgw-test-001",
					SecurityType:        "STANDARD-TLS",
					UseDefaultMap:       false,
					UseDefaultTTL:       false,
					HTTPSServiceBinding: ptr.To("H3"),
				}},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				Submitter:        "ftzgvvigljhoq5ib",
				SubmitterEmail:   "ftzgvvigljhoq5ib@nomail-akamai.com",
			},
		},
		"202 Accepted - removing httpsServiceBinding with op remove": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-001",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:   "remove",
						Path: "/httpsServiceBinding",
					},
				},
			},
			responseStatus: http.StatusAccepted,
			responseBody: `
{
    "action": "EDIT",
    "changeId": 66025607,
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
			expectedBody: `[{"op":"remove","path":"/httpsServiceBinding","value":""}]`,
			expectedResponse: &UpdateEdgeHostnameResponse{
				Action:   "EDIT",
				ChangeID: 66025607,
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
				}},
				Status:           "PENDING",
				StatusMessage:    "File uploaded and awaiting validation",
				StatusUpdateDate: test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				SubmitDate:       test.NewTimeFromString(t, "2021-09-23T15:07:10.000+00:00"),
				Submitter:        "ftzgvvigljhoq5ib",
				SubmitterEmail:   "ftzgvvigljhoq5ib@nomail-akamai.com",
			},
		},
		"request builder error - malformed path": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:    "edgesuite.net",
				RecordName: "%",
				Body:       []UpdateEdgeHostnameRequestBody{{Op: "replace", Path: "/ttl", Value: "10000"}},
			},
			withError: ErrUpdateEdgeHostname,
		},
		"validation error - invalid httpsServiceBinding value": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:    "edgesuite.net",
				RecordName: "mgw-test-001",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "replace",
						Path:  "/httpsServiceBinding",
						Value: "INVALID",
					},
				},
			},
			withError: ErrUpdateEdgeHostname,
		},
		"validation error - invalid op add for /ttl path": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:    "edgesuite.net",
				RecordName: "mgw-test-001",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "add",
						Path:  "/ttl",
						Value: "10000",
					},
				},
			},
			withError: ErrUpdateEdgeHostname,
		},
		"validation error - invalid op remove for /ttl path": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:    "edgesuite.net",
				RecordName: "mgw-test-001",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "remove",
						Path:  "/ttl",
						Value: "10000",
					},
				},
			},
			withError: ErrUpdateEdgeHostname,
		},
		"validation error - invalid op add for /ipVersionBehavior path": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:    "edgesuite.net",
				RecordName: "mgw-test-001",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "add",
						Path:  "/ipVersionBehavior",
						Value: "IPV4",
					},
				},
			},
			withError: ErrUpdateEdgeHostname,
		},
		"validation error - invalid op remove for /ipVersionBehavior path": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:    "edgesuite.net",
				RecordName: "mgw-test-001",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "remove",
						Path:  "/ipVersionBehavior",
						Value: "IPV4",
					},
				},
			},
			withError: ErrUpdateEdgeHostname,
		},
		"validation error - invalid op for /httpsServiceBinding path": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:    "edgesuite.net",
				RecordName: "mgw-test-001",
				Body: []UpdateEdgeHostnameRequestBody{
					{
						Op:    "test",
						Path:  "/httpsServiceBinding",
						Value: "H2",
					},
				},
			},
			withError: ErrUpdateEdgeHostname,
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
			expectedBody: `[{"op":"","path":"/incorrect","value":"some Value"}]`,
			withError:    ErrUpdateEdgeHostname,
		},
		"500 internal server error": {
			request: UpdateEdgeHostnameRequest{
				DNSZone:           "edgesuite.net",
				RecordName:        "mgw-test-002",
				StatusUpdateEmail: []string{"some@example.com"},
				Comments:          "some comment",
				Body:              []UpdateEdgeHostnameRequestBody{{Op: "replace", Path: "/ttl", Value: "10000"}},
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
			expectedBody: `[{"op":"replace","path":"/ttl","value":"10000"}]`,
			withError:    ErrUpdateEdgeHostname,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPatch, r.Method)
				assert.Equal(t, "application/json-patch+json", r.Header.Get("Content-Type"))
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, test.expectedBody, string(body))
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
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

func TestGetCertificate(t *testing.T) {
	tests := map[string]struct {
		request          GetCertificateRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCertificateResponse
		withError        error
	}{
		"200 OK": {
			request: GetCertificateRequest{
				DNSZone:    "edgekey.net",
				RecordName: "mgw-test-002",
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"certificateId": "1234",
				"commonName": "example.com",
				"serialNumber": "12:34:56:78:90:AB:CD:EF",
				"slotNumber": 8927,
				"expirationDate": "2019-10-31T23:59:59Z",
				"certificateType": "SAN",
				"validationType": "DOMAIN_VALIDATION",
				"status": "PENDING",
				"availableDomains": [
				"live.example.com",
				"secure.example.com",
				"www.example.com"
				]
			}`,
			expectedPath: "/hapi/v1/dns-zones/edgekey.net/edge-hostnames/mgw-test-002/certificate",
			expectedResponse: &GetCertificateResponse{
				CertificateID:   "1234",
				CommonName:      "example.com",
				SerialNumber:    "12:34:56:78:90:AB:CD:EF",
				SlotNumber:      8927,
				ExpirationDate:  test.NewTimeFromString(t, "2019-10-31T23:59:59Z"),
				CertificateType: "SAN",
				ValidationType:  "DOMAIN_VALIDATION",
				Status:          "PENDING",
				AvailableDomains: []string{
					"live.example.com",
					"secure.example.com",
					"www.example.com",
				},
			},
		},
		"404 certificate not found": {
			request: GetCertificateRequest{
				DNSZone:    "edgekey.net",
				RecordName: "unknown",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
			{
				"type": "CERTIFICATE_NOT_FOUND",
				"title": "Certificate Not Found",
				"status": 404,
				"detail": "Details are not available for this certificate; the certificate is missing or access is denied",
				"instance": "/hapi/error-instances/a30f67cc-df20-4e02-bbc3-cf7c204a4aab",
				"requestInstance": "http://origin.pulsar.akamai.com/hapi/open/v1/dns-zones/edgekey.net/edge-hostnames/example.com/certificate?depth=ALL&accountSwitchKey=F-AC-1937217#d7aa7348",
				"method": "GET",
				"requestTime": "2022-11-30T18:51:43.482982Z",
				"errors": [],
				"extensionFields": []
			}`,
			expectedPath: "/hapi/v1/dns-zones/edgekey.net/edge-hostnames/unknown/certificate",
			withError: fmt.Errorf("%s: %s: %w", ErrGetCertificate, ErrNotFound, &Error{
				Type:            "CERTIFICATE_NOT_FOUND",
				Title:           "Certificate Not Found",
				Status:          404,
				Detail:          "Details are not available for this certificate; the certificate is missing or access is denied",
				Instance:        "/hapi/error-instances/a30f67cc-df20-4e02-bbc3-cf7c204a4aab",
				RequestInstance: "http://origin.pulsar.akamai.com/hapi/open/v1/dns-zones/edgekey.net/edge-hostnames/example.com/certificate?depth=ALL&accountSwitchKey=F-AC-1937217#d7aa7348",
				Method:          "GET",
				RequestTime:     test.NewTimeFromString(t, "2022-11-30T18:51:43.482982Z"),
			}),
		},
		"500 internal server error": {
			request: GetCertificateRequest{
				DNSZone:    "edgekey.net",
				RecordName: "mgw-test-002",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error deleting activation",
				"status": 500
			}`,
			expectedPath: "/hapi/v1/dns-zones/edgekey.net/edge-hostnames/mgw-test-002/certificate",
			withError: &Error{
				Type:   "internal_error",
				Title:  "Internal Server Error",
				Detail: "Error deleting activation",
				Status: http.StatusInternalServerError,
			},
		},
		"missing required values": {
			request:   GetCertificateRequest{},
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCertificate(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
