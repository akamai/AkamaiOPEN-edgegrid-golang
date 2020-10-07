package gtm

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verify GetListDomains. Sould pass, e.g. no API errors and non nil list.
func TestGtm_ListDomains(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DomainList
		withError        error
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
			{
                            "items" : [ {
                                "name" : "gtmdomtest.akadns.net",
                                "status" : "Change Pending",
                                "acgId" : "1-3CV382",
                                "lastModified" : "2019-06-06T19:07:20.000+00:00",
                                "lastModifiedBy" : "operator",
                                "changeId" : "c3e1b771-2500-40c9-a7da-6c3cdbce1936",
                                "activationState" : "PENDING",
                                "modificationComments" : "mock test",
                                "links" : [ {
                                    "rel" : "self",
                                    "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net"
                                } ]
                            } ]
			}`,
			expectedPath: "/config-gtm/v1/domains",
			expectedResponse: &DomainList{
				DomainItems: []*DomainItem{
					{
						AcgId:        "1-2345",
						LastModified: "2014-03-03T16:02:45.000+0000",
						Name:         "example.akadns.net",
						Status:       "2014-02-20 22:56 GMT: Current configuration has been propagated to all GTM name servers",
						Links: []*link{
							{
								"href": "/config-gtm/v1/domains/example.akadns.net",
								"rel":  "self",
							},
						},
					},
					{
						AcgIdx:       "1-2345",
						LastModified: "2013-11-09T12:04:45.000+0000",
						Name:         "demo.akadns.net",
						Status:       "2014-02-20 22:56 GMT: Current configuration has been propagated to all GTM name servers",
						Links: []*link{
							{
								"href": "/config-gtm/v1/domains/demo.akadns.net",
								"rel":  "self",
							},
						},
					},
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	// ** TODO **
	//                 SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8")
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
			result, err := client.ListZones(context.Background(), test.args...)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test NewDomain
// NewDomain(context.Context, string, string) *Domain
func TestDns_NewDomain(t *testing.T) {
	client := Client(session.Must(session.New()))

	inp := Domain{
		Name: "testdomain.akadns.net",
		Type: "weighted",
	}

	out := client.NewDomain(context.Background(), inp.Name, inp.Type)

	assert.ObjectsAreEqual(&inp, out)
}

/*

TODO:: NullFieldMap does Get Domain .... need to set up full mock ...

// Retrieve map of null fields
// NullFieldMap(context.Context, *Domain) (*NullFieldMapStruct, error)
func TestDns_NullFieldMap(t *testing.T) {
	client := Client(session.Must(session.New()))

        inp := Domain{
                Name:   "testdomain.akadns.net",
                Type:   "weighted",
        }
	nullfieldout := NullFieldMapStruct{}

	out, err := client.NewZoneResponse(context.Background(), &inp)

	assert.Equal(t, err, nil)
	assert.Equal(t, out., "example.com")
}
*/

// Test GetDomain
// GetDomain(context.Context, string) (*Domain, error)
func TestDns_GetDomain(t *testing.T) {
	tests := map[string]struct {
		domain           string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Domain
		withError        error
	}{
		"200 OK": {
			domain:         "gtmdomtest.akadns.net",
			responseStatus: http.StatusOK,
			responseBody: `
			{ 
                          "cidrMaps": [], 
                          "datacenters": [
                              {
                                  "city": "Snæfellsjökull", 
                                  "cloneOf": null, 
                                  "cloudServerTargeting": false, 
                                  "continent": "EU", 
                                  "country": "IS", 
                                  "datacenterId": 3132, 
                                  "defaultLoadObject": {
                                       "loadObject": null, 
                                       "loadObjectPort": 0, 
                                       "loadServers": null
                                   }, 
                                   "latitude": 64.808, 
                                   "links": [
                                       {
                                            "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3132", 
                                            "rel": "self"
                                       }
                                    ], 
                                    "longitude": -23.776, 
                                    "nickname": "property_test_dc2", 
                                    "stateOrProvince": null, 
                                    "virtual": true
                              }
                          ], 
                          "defaultErrorPenalty": 75, 
                          "defaultSslClientCertificate": null, 
                          "defaultSslClientPrivateKey": null, 
                          "defaultTimeoutPenalty": 25, 
                          "emailNotificationList": [], 
                          "geographicMaps": [], 
                          "lastModified": "2019-04-25T14:53:12.000+00:00", 
                          "lastModifiedBy": "operator", 
                          "links": [
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net", 
                                   "rel": "self"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters", 
                                   "rel": "datacenters"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties", 
                                   "rel": "properties"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps", 
                                   "rel": "geographic-maps"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps", 
                                   "rel": "cidr-maps"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources", 
                                   "rel": "resources"
                              }
                          ], 
                          "loadFeedback": false, 
                          "loadImbalancePercentage": 10.0, 
                          "modificationComments": "Edit Property test_property", 
                          "name": "gtmdomtest.akadns.net", 
                          "properties": [
                               {
                                    "backupCName": null, 
                                    "backupIp": null, 
                                    "balanceByDownloadScore": false, 
                                    "cname": "www.boo.wow", 
                                    "comments": null, 
                                    "dynamicTTL": 300, 
                                    "failbackDelay": 0, 
                                    "failoverDelay": 0, 
                                    "handoutMode": "normal", 
                                    "healthMax": null, 
                                    "healthMultiplier": null, 
                                    "healthThreshold": null, 
                                    "ipv6": false, 
                                    "lastModified": "2019-04-25T14:53:12.000+00:00", 
                                    "links": [
                                         {
                                              "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/test_property", 
                                              "rel": "self"
                                         }
                                    ], 
                                    "livenessTests": [
                                         {
                                               "disableNonstandardPortWarning": false, 
                                               "hostHeader": null, 
                                               "httpError3xx": true, 
                                               "httpError4xx": true, 
                                               "httpError5xx": true, 
                                               "name": "health check", 
                                               "requestString": null, 
                                               "responseString": null, 
                                               "sslClientCertificate": null, 
                                               "sslClientPrivateKey": null, 
                                               "testInterval": 60, 
                                               "testObject": "/status", 
                                               "testObjectPassword": null, 
                                               "testObjectPort": 80, 
                                               "testObjectProtocol": "HTTP", 
                                               "testObjectUsername": null, 
                                               "testTimeout": 25.0
                                         }
                                    ], 
                                    "loadImbalancePercentage": 10.0, 
                                    "mapName": null, 
                                    "maxUnreachablePenalty": null, 
                                    "mxRecords": [], 
                                    "name": "test_property", 
                                    "scoreAggregationType": "mean", 
                                    "staticTTL": 600, 
                                    "stickinessBonusConstant": null, 
                                    "stickinessBonusPercentage": 50, 
                                    "trafficTargets": [
                                         {
                                              "datacenterId": 3131, 
                                              "enabled": true, 
                                              "handoutCName": null, 
                                              "name": null, 
                                              "servers": [
                                                   "1.2.3.4", 
                                                   "1.2.3.5"
                                              ], 
                                              "weight": 50.0
                                         }, 
                                         {
                                              "datacenterId": 3132, 
                                              "enabled": true, 
                                              "handoutCName": "www.google.com", 
                                              "name": null, 
                                              "servers": [], 
                                              "weight": 25.0
                                         }, 
                                         {
                                              "datacenterId": 3133, 
                                              "enabled": true, 
                                              "handoutCName": "www.comcast.com", 
                                              "name": null, 
                                              "servers": [
                                                    "www.comcast.com"
                                              ], 
                                              "weight": 25.0
                                         }
                                    ], 
                                    "type": "weighted-round-robin", 
                                    "unreachableThreshold": null, 
                                    "useComputedTargets": false
                               }
                          ], 
                          "resources": [], 
                          "status": {
                               "changeId": "40e36abd-bfb2-4635-9fca-62175cf17007", 
                               "links": [
                                     {
                                          "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current", 
                                          "rel": "self"
                                     }
                               ], 
                               "message": "Current configuration has been propagated to all GTM nameservers", 
                               "passingValidation": true, 
                               "propagationStatus": "COMPLETE", 
                               "propagationStatusDate": "2019-04-25T14:54:00.000+00:00"
                          }, 
                          "type": "weighted"
			}`,
			expectedPath: "/config-gtm/v1/domains/gtmdomtest.akadns.net",
			expectedResponse: &Domain{
				CidrMaps: make([]interface{}),
				Datacenters: []*Datacenter{
					{
						City:                 "Snæfellsjökull",
						CloneOf:              null,
						CloudServerTargeting: false,
						Continent:            "EU",
						Country:              "IS",
						DatacenterId:         3132,
						DefaultLoadObject: *LoadObject{
							LoadObject:     "",
							LoadObjectPort: 0,
							LoadServers:    make([]string),
						},
						Latitude: 64.808,
						Links: []*Link{
							{
								Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3132",
								Rel:  "self",
							},
						},
						Longitude:       -23.776,
						Nickname:        "property_test_dc2",
						StateOrProvince: "",
						Virtual:         true,
					},
				},
				DefaultErrorPenalty:         75,
				DefaultSslClientCertificate: "",
				DefaultSslClientPrivateKey:  "",
				DefaultTimeoutPenalty:       25,
				EmailNotificationList:       make([]string),
				GeographicMaps:              make([]interface{}),
				LastModified:                "2019-04-25T14:53:12.000+00:00",
				LastModifiedBy:              "operator",
				Links: []*Link{
					{
						Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net",
						Rel:  "self",
					},
					{
						Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters",
						Rel:  "datacenters",
					},
					{
						Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties",
						Rel:  "properties",
					},
					{
						Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps",
						Rel:  "geographic-maps",
					},
					{
						Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps",
						Rel:  "cidr-maps",
					},
					{
						Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources",
						Rel:  "resources",
					},
				},
				LoadFeedback:            false,
				LoadImbalancePercentage: 10.0,
				ModificationComments:    "Edit Property test_property",
				Name:                    "gtmdomtest.akadns.net",
				Properties: []*Property{
					{
						BackupCName:            "",
						BackupIp:               null,
						BalanceByDownloadScore: false,
						CName:                  "www.boo.wow",
						Comments:               "",
						DynamicTTL:             300,
						FailbackDelay:          0,
						FailoverDelay:          0,
						HandoutMode:            "normal",
						HealthMax:              0,
						HealthMultiplier:       0,
						HealthThreshold:        0,
						Ipv6:                   false,
						LastModified:           "2019-04-25T14:53:12.000+00:00",
						Links: []*Link{
							{
								Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/test_property",
								Rel:  "self",
							},
						},
						LivenessTests: []*LivenessTest{
							{
								DisableNonstandardPortWarning: false,
								HostHeader:                    "",
								HttpError3xx:                  true,
								HttpError4xx:                  true,
								HttpError5xx:                  true,
								Name:                          "health check",
								RequestString:                 "",
								ResponseString:                "",
								SslClientCertificate:          "",
								SslClientPrivateKey:           "",
								TestInterval:                  60,
								TestObject:                    "/status",
								TestObjectPassword:            "",
								TestObjectPort:                80,
								TestObjectProtocol:            "HTTP",
								TestObjectUsername:            "",
								TestTimeout:                   25.0,
							},
						},
						LoadImbalancePercentage:   10.0,
						MapName:                   "",
						MaxUnreachablePenalty:     0,
						MxRecords:                 make([]interface{}),
						Name:                      "test_property",
						ScoreAggregationType:      "mean",
						StaticTTL:                 600,
						StickinessBonusConstant:   "",
						StickinessBonusPercentage: 50,
						TrafficTargets: []*TrafficTarget{
							{
								DatacenterId: 3131,
								Enabled:      true,
								HandoutCName: "",
								Name:         null,
								Servers: []string{
									"1.2.3.4",
									"1.2.3.5",
								},
								Weight: 50.0,
							},
						},
						Type:                 "weighted-round-robin",
						UnreachableThreshold: 0,
						UseComputedTargets:   false,
					},
				},
				Resources: make([]interface{}),
				Status: {
					ChangeId: "40e36abd-bfb2-4635-9fca-62175cf17007",
					Links: []*Link{
						{
							Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
							Rel:  "self",
						},
					},
					Message:               "Current configuration has been propagated to all GTM nameservers",
					PassingValidation:     true,
					PropagationStatus:     "COMPLETE",
					PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
				},
				Type: "weighted",
			},
		},
		"500 internal server error": {
			domain:         "gtmdomtest.akadns.net",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/gtmdomtest.akadns.net",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"404 not found error": {
			domain:         "baddomain.akadns.net",
			responseStatus: http.StatusNotFoundError,
			responseBody: `
{
    "type": "not_found",
    "title": "Not Found Error",
    "detail": "Domain not found",
    "status": 404
}`,
			expectedPath: "/config-gtm/v1/domains/baddomain.akadns.net",
			withError: &Error{
				Type:       "not_found",
				Title:      "Not Found Error",
				Detail:     "Domain not found",
				StatusCode: http.StatusNotFoundError,
			},
		},
	}

	// ** TODO **
	//                 SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8")
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
			result, err := client.GetDomain(context.Background(), test.domain)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create domain.
// CreateDomain(context.Context, *Domain, map[string]string) (*DomainResponse, error)
func TestDns_CreateDomain(t *testing.T) {
	tests := map[string]struct {
		domain         Domain
		query          QueryStringMap
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"201 Created": {
			domain: Domain{
				Name: "gtmdomtest.akadns.net",
				Type: "basic",
			},
			query:          map[string]string{contractId: "1-2ABCDE"},
			responseStatus: http.StatusCreated,
			responseBody: `
		        {
                          "resource" : {
                                "cnameCoalescingEnabled" : false,
                                "defaultErrorPenalty" : 75,
                                "defaultHealthMax" : null,
                                "defaultHealthMultiplier" : null,
                                "defaultHealthThreshold" : null,
                                "defaultMaxUnreachablePenalty" : null,
                                "defaultSslClientCertificate" : null,
                                "defaultSslClientPrivateKey" : null,
                                "defaultTimeoutPenalty" : 25,
                                "defaultUnreachableThreshold" : null,
                                "emailNotificationList" : [ ],
                                "endUserMappingEnabled" : false,
                                "lastModified" : "2019-06-24T18:48:57.787+00:00",
                                "lastModifiedBy" : "operator",
                                "loadFeedback" : false,
                                "mapUpdateInterval" : 0,
                                "maxProperties" : 0,
                                "maxResources" : 512,
                                "maxTestTimeout" : 0.0,
                                "maxTTL" : 0,
                                "minPingableRegionFraction" : null,
                                "minTestInterval" : 0,
                                "minTTL" : 0,
                                "modificationComments" : null,
                                "name" : "gtmdomtest.akadns.net",
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "roundRobinPrefix" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "servermonitorPool" : null,
                                "type" : "basic",
                                "status" : {
                                        "message" : "Change Pending",
                                        "changeId" : "539872cc-6ba6-4429-acd5-90bab7fb5e9d",
                                        "propagationStatus" : "PENDING",
                                        "propagationStatusDate" : "2019-06-24T18:48:57.787+00:00",
                                        "passingValidation" : true,
                                        "links" : [ {
                                                "rel" : "self",
                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                        } ]
                                },
                                "loadImbalancePercentage" : null,
                                "domainVersionId" : null,
                                "resources" : [ ],
                                "properties" : [ ],
                                "datacenters" : [ ],
                                "geographicMaps" : [ ],
                                "cidrMaps" : [ ],
                                "asMaps" : [ ],
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net"
                                    }, {
                                        "rel" : "datacenters",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters"
                                    }, {
                                        "rel" : "properties",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties"
                                    }, {
                                        "rel" : "geographic-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps"
                                    }, {
                                        "rel" : "cidr-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps"
                                    }, {
                                        "rel" : "resources",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources"
                                    }, {
                                        "rel" : "as-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/as-maps"
                                    } ]
                                },
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "539872cc-6ba6-4429-acd5-90bab7fb5e9d",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-06-24T18:48:57.787+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                } ]
			}`,
			expectedResponse: &DomainResponse{
				Resource: {
					CnameCoalescingEnabled: false,
					DefaultErrorPenalty:    75,
					DefaultTimeoutPenalty:  25,
					EndUserMappingEnabled:  false,
					LastModified:           "2019-06-24T18:48:57.787+00:00",
					LastModifiedBy:         "operator",
					LoadFeedback:           false,
					MapUpdateInterval:      0,
					MaxProperties:          0,
					MaxResources:           512,
					MaxTestTimeout:         0.0,
					MaxTTL:                 0,
					MinTestInterval:        0,
					MinTTL:                 0,
					Name:                   "gtmdomtest.akadns.net",
					Type:                   "basic",
					Status: *ResponseStatus{
						Message:               "Change Pending",
						ChangeId:              "539872cc-6ba6-4429-acd5-90bab7fb5e9d",
						PropagationStatus:     "PENDING",
						PropagationStatusDate: "2019-06-24T18:48:57.787+00:00",
						PassingValidation:     true,
						Links: []*Link{
							{
								Rel:  "self",
								Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
							},
						},
					},
					Links: []*Link{
						{
							Rel:  "self",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net",
						},
						{
							Rel:  "datacenters",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters",
						},
						{
							Rel:  "properties",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties",
						},
						{
							Rel:  "geographic-maps",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps",
						},
						{
							Rel:  "cidr-maps",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps",
						},
						{
							Rel:  "resources",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources",
						},
						{
							Rel:  "as-maps",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/as-maps",
						},
					},
				},
				Status: &ResponseStatus{
					Message:           "Change Pending",
					ChangeId:          "40e36abd-bfb2-4635-9fca-62175cf17007",
					PropagationStatus: "PENDING",
					PassingValidation: true,
					Links: []*Link{
						{
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
							Rel:  "self",
						},
					},
				},
			},
			expectedPath: "/config-gtm/v1/domains?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			domain: Domain{
				Name: "gtmdomtest.akadns.net",
				Type: "basic",
			},
			query:          map[string]string{contractId: "1-2ABCDE"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains?contractId=1-2ABCDE",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	// ** TODO **
	//                 SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8")
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
			err := client.CreateDomain(context.Background(), &test.domain, test.query)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// Test Update domain.
// UpdateDomain(context.Context, *Domain, map[string]string) (*DomainResponse, error)
func TestDns_UpdateDomain(t *testing.T) {
	tests := map[string]struct {
		domain         Domain
		query          map[string]string
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"200 Success": {
			domain: Domain{
				EndUserMappingEnabled: false,
				Name:                  "gtmdomtest.akadns.net",
				Type:                  "basic",
			},
			query:          map[string]string{contractId: "1-2ABCDE"},
			responseStatus: http.StatusCreated,
			responseBody: `
                        {
                          "resource" : {
                                "cnameCoalescingEnabled" : false,
                                "defaultErrorPenalty" : 75,
                                "defaultHealthMax" : null,
                                "defaultHealthMultiplier" : null,
                                "defaultHealthThreshold" : null,
                                "defaultMaxUnreachablePenalty" : null,
                                "defaultSslClientCertificate" : null,
                                "defaultSslClientPrivateKey" : null,
                                "defaultTimeoutPenalty" : 25,
                                "defaultUnreachableThreshold" : null,
                                "emailNotificationList" : [ ],
                                "endUserMappingEnabled" : false,
                                "lastModified" : "2019-06-24T18:48:57.787+00:00",
                                "lastModifiedBy" : "operator",
                                "loadFeedback" : false,
                                "mapUpdateInterval" : 0,
                                "maxProperties" : 0,
                                "maxResources" : 512,
                                "maxTestTimeout" : 0.0,
                                "maxTTL" : 0,
                                "minPingableRegionFraction" : null,
                                "minTestInterval" : 0,
                                "minTTL" : 0,
                                "modificationComments" : null,
                                "name" : "gtmdomtest.akadns.net",
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "roundRobinPrefix" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "servermonitorPool" : null,
                                "type" : "basic",
                                "status" : {
                                        "message" : "Change Pending",
                                        "changeId" : "539872cc-6ba6-4429-acd5-90bab7fb5e9d",
                                        "propagationStatus" : "PENDING",
                                        "propagationStatusDate" : "2019-06-24T18:48:57.787+00:00",
                                        "passingValidation" : true,
                                        "links" : [ {
                                                "rel" : "self",
                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                        } ]
                                },
                                "loadImbalancePercentage" : null,
                                "domainVersionId" : null,
                                "resources" : [ ],
                                "properties" : [ ],
                                "datacenters" : [ ],
                                "geographicMaps" : [ ],
                                "cidrMaps" : [ ],
                                "asMaps" : [ ],
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net"
                                    }, {
                                        "rel" : "datacenters",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters"
                                    }, {
                                        "rel" : "properties",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties"
                                    }, {
                                        "rel" : "geographic-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps"
                                    }, {
                                        "rel" : "cidr-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps"
                                    }, {
                                        "rel" : "resources",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources"
                                    }, {
                                        "rel" : "as-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/as-maps"
                                    } ]
                                },
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "539872cc-6ba6-4429-acd5-90bab7fb5e9d",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-06-24T18:48:57.787+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                } ]
                        }`,
			expectedResponse: &DomainResponse{
				Resource: {
					CnameCoalescingEnabled: false,
					DefaultErrorPenalty:    75,
					DefaultTimeoutPenalty:  25,
					EndUserMappingEnabled:  false,
					LastModified:           "2019-06-24T18:48:57.787+00:00",
					LastModifiedBy:         "operator",
					LoadFeedback:           false,
					MapUpdateInterval:      0,
					MaxProperties:          0,
					MaxResources:           512,
					MaxTestTimeout:         0.0,
					MaxTTL:                 0,
					MinTestInterval:        0,
					MinTTL:                 0,
					Name:                   "gtmdomtest.akadns.net",
					Type:                   "basic",
					Status: *ResponseStatus{
						Message:               "Change Pending",
						ChangeId:              "539872cc-6ba6-4429-acd5-90bab7fb5e9d",
						PropagationStatus:     "PENDING",
						PropagationStatusDate: "2019-06-24T18:48:57.787+00:00",
						PassingValidation:     true,
						Links: []*Link{
							{
								Rel:  "self",
								Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
							},
						},
					},
					Links: []*Link{
						{
							Rel:  "self",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net",
						},
						{
							Rel:  "datacenters",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters",
						},
						{
							Rel:  "properties",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties",
						},
						{
							Rel:  "geographic-maps",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps",
						},
						{
							Rel:  "cidr-maps",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps",
						},
						{
							Rel:  "resources",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources",
						},
						{
							Rel:  "as-maps",
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/as-maps",
						},
					},
				},
				Status: &ResponseStatus{
					Message:           "Change Pending",
					ChangeId:          "40e36abd-bfb2-4635-9fca-62175cf17007",
					PropagationStatus: "PENDING",
					PassingValidation: true,
					Links: []*Link{
						{
							Href: "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
							Rel:  "self",
						},
					},
				},
			},
			expectedPath: "/config-gtm/v1/domains?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			domain: Domain{
				Name: "gtmdomtest.akadns.net",
				Type: "basic",
			},
			query:          map[string]string{contractId: "1-2ABCDE"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains?contractId=1-2ABCDE",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	// ** TODO **
	//                 SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8")
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.CreateDomain(context.Background(), &test.domain, test.query)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

/* Future. Presently no domain Delete endpoint.
func TestDeleteDomain(t *testing.T) {

        defer gock.Off()

        mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/"+gtmTestDomain)
        mock.
                Delete("/config-gtm/v1/domains/"+gtmTestDomain).
                HeaderPresent("Authorization").
                Reply(200).
                SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
                BodyString(`{
                        "resource" : null,
                        "status" : {
                               "changeId": "40e36abd-bfb2-4635-9fca-62175cf17007",
                               "links": [
                                     {
                                          "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
                                          "rel": "self"
                                     }
                               ],
                               "message": "Change Pending",
                               "passingValidation": true,
                               "propagationStatus": "PENDING",
                               "propagationStatusDate": "2019-04-25T14:54:00.000+00:00"
                          },
                }`)

        Init(config)

        getDomain := instantiateDomain()

        _, err := getDomain.Delete()
        assert.NoError(t, err)

}
*/
