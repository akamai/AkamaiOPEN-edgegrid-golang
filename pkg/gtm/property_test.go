package gtm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGTM_ListProperties(t *testing.T) {
	var result PropertyList

	respData, err := loadTestData("TestGTM_ListProperties.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           ListPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Property
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: ListPropertiesRequest{
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/properties",
			expectedResponse: result.PropertyItems,
		},
		"500 internal server error": {
			params: ListPropertiesRequest{
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching propertys",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching propertys",
				StatusCode: http.StatusInternalServerError,
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
			result, err := client.ListProperties(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_GetProperty(t *testing.T) {
	var result GetPropertyResponse

	respData, err := loadTestData("TestGTM_GetProperty.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           GetPropertyRequest
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *GetPropertyResponse
		withError        error
	}{
		"200 OK": {
			params: GetPropertyRequest{
				PropertyName: "www",
				DomainName:   "example.akadns.net",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/properties/www",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetPropertyRequest{
				PropertyName: "www",
				DomainName:   "example.akadns.net",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching property"
}`),
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties/www",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching property",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write(test.responseBody)
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetProperty(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_CreateProperty(t *testing.T) {
	tests := map[string]struct {
		params           CreatePropertyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreatePropertyResponse
		withError        bool
		assertError      func(*testing.T, error)
		headers          http.Header
	}{
		"201 Created": {
			params: CreatePropertyRequest{
				Property: &Property{
					BalanceByDownloadScore: false,
					HandoutMode:            "normal",
					IPv6:                   false,
					Name:                   "origin",
					ScoreAggregationType:   "mean",
					StaticTTL:              600,
					Type:                   "weighted-round-robin",
					UseComputedTargets:     false,
					LivenessTests: []LivenessTest{
						{
							DisableNonstandardPortWarning: false,
							HTTPError3xx:                  true,
							HTTPError4xx:                  true,
							HTTPError5xx:                  true,
							Name:                          "health-check",
							TestInterval:                  60,
							TestObject:                    "/status",
							TestObjectPort:                80,
							TestObjectProtocol:            "HTTP",
							TestTimeout:                   25.0,
						},
					},
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 3134,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.5"},
						},
						{
							DatacenterID: 3133,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.4"},
							Precedence:   nil,
						},
					},
				},
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "resource": {
        "backupCName": null,
        "backupIp": null,
        "balanceByDownloadScore": false,
        "cname": null,
        "comments": null,
        "dynamicTTL": 300,
        "failbackDelay": 0,
        "failoverDelay": 0,
        "handoutMode": "normal",
        "healthMax": null,
        "healthMultiplier": null,
        "healthThreshold": null,
        "ipv6": false,
        "lastModified": null,
        "loadImbalancePercentage": null,
        "mapName": null,
        "maxUnreachablePenalty": null,
        "name": "origin",
        "scoreAggregationType": "mean",
        "staticTTL": 600,
        "stickinessBonusConstant": 0,
        "stickinessBonusPercentage": 0,
        "type": "weighted-round-robin",
        "unreachableThreshold": null,
        "useComputedTargets": false,
        "mxRecords": [],
        "links": [
            {
                "href": "/config-gtm/v1/domains/example.akadns.net/properties/origin",
                "rel": "self"
            }
        ],
        "livenessTests": [
            {
                "disableNonstandardPortWarning": false,
                "hostHeader": "foo.example.com",
                "httpError3xx": true,
                "httpError4xx": true,
                "httpError5xx": true,
                "name": "health-check",
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
        "trafficTargets": [
            {
                "datacenterId": 3134,
                "enabled": true,
                "handoutCName": null,
                "name": null,
                "weight": 50.0,
                "servers": [
                    "1.2.3.5"
                ],
                "precedence": null
            },
            {
                "datacenterId": 3133,
                "enabled": true,
                "handoutCName": null,
                "name": null,
                "weight": 50.0,
                "servers": [
                    "1.2.3.4"
                ],
                "precedence": null
            }
        ]
    },
    "status": {
        "changeId": "eee0c3b4-0e45-4f4b-822c-7dbc60764d18",
        "message": "Change Pending",
        "passingValidation": true,
        "propagationStatus": "PENDING",
        "propagationStatusDate": "2014-04-15T11:30:27.000+0000",
        "links": [
            {
                "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                "rel": "self"
            }
        ]
    }
}
`,
			expectedResponse: &CreatePropertyResponse{
				Resource: &Property{
					BalanceByDownloadScore: false,
					HandoutMode:            "normal",
					IPv6:                   false,
					Name:                   "origin",
					ScoreAggregationType:   "mean",
					StaticTTL:              600,
					DynamicTTL:             300,
					Type:                   "weighted-round-robin",
					UseComputedTargets:     false,
					LivenessTests: []LivenessTest{
						{
							DisableNonstandardPortWarning: false,
							HTTPError3xx:                  true,
							HTTPError4xx:                  true,
							HTTPError5xx:                  true,
							Name:                          "health-check",
							TestInterval:                  60,
							TestObject:                    "/status",
							TestObjectPort:                80,
							TestObjectProtocol:            "HTTP",
							TestTimeout:                   25.0,
						},
					},
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 3134,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.5"},
						},
						{
							DatacenterID: 3133,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.4"},
						},
					},
					Links: []Link{
						{
							Href: "/config-gtm/v1/domains/example.akadns.net/properties/origin",
							Rel:  "self",
						},
					},
				},
				Status: &ResponseStatus{
					ChangeID:              "eee0c3b4-0e45-4f4b-822c-7dbc60764d18",
					Message:               "Change Pending",
					PassingValidation:     true,
					PropagationStatus:     "PENDING",
					PropagationStatusDate: "2014-04-15T11:30:27.000+0000",
					Links: []Link{
						{
							Href: "/config-gtm/v1/domains/example.akadns.net/status/current",
							Rel:  "self",
						},
					},
				},
			},
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties/origin",
		},
		"201 Created - ranked-failover": {
			params: CreatePropertyRequest{
				Property: &Property{
					BalanceByDownloadScore: false,
					HandoutMode:            "normal",
					IPv6:                   false,
					Name:                   "origin",
					ScoreAggregationType:   "mean",
					StaticTTL:              600,
					Type:                   "ranked-failover",
					UseComputedTargets:     false,
					LivenessTests: []LivenessTest{
						{
							DisableNonstandardPortWarning: false,
							HTTPError3xx:                  true,
							HTTPError4xx:                  true,
							HTTPError5xx:                  true,
							HTTPMethod:                    ptr.To("GET"),
							HTTPRequestBody:               ptr.To("TestBody"),
							Name:                          "health-check",
							TestInterval:                  60,
							TestObject:                    "/status",
							TestObjectPort:                80,
							TestObjectProtocol:            "HTTP",
							TestTimeout:                   25.0,
							Pre2023SecurityPosture:        true,
							AlternateCACertificates:       []string{"test1"},
						},
					},
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 3134,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.5"},
							Precedence:   ptr.To(255),
						},
						{
							DatacenterID: 3133,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.4"},
							Precedence:   nil,
						},
					},
				},
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "resource": {
        "backupCName": null,
        "backupIp": null,
        "balanceByDownloadScore": false,
        "cname": null,
        "comments": null,
        "dynamicTTL": 300,
        "failbackDelay": 0,
        "failoverDelay": 0,
        "handoutMode": "normal",
        "healthMax": null,
        "healthMultiplier": null,
        "healthThreshold": null,
        "ipv6": false,
        "lastModified": null,
        "loadImbalancePercentage": null,
        "mapName": null,
        "maxUnreachablePenalty": null,
        "name": "origin",
        "scoreAggregationType": "mean",
        "staticTTL": 600,
        "stickinessBonusConstant": 0,
        "stickinessBonusPercentage": 0,
        "type": "weighted-round-robin",
        "unreachableThreshold": null,
        "useComputedTargets": false,
        "mxRecords": [],
        "links": [
            {
                "href": "/config-gtm/v1/domains/example.akadns.net/properties/origin",
                "rel": "self"
            }
        ],
        "livenessTests": [
            {
                "disableNonstandardPortWarning": false,
                "hostHeader": "foo.example.com",
                "httpError3xx": true,
                "httpError4xx": true,
                "httpError5xx": true,
				"httpMethod": "GET",
				"httpRequestBody": "TestBody",
				"pre2023SecurityPosture": true,
				"alternateCACertificates": ["test1"],
                "name": "health-check",
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
        "trafficTargets": [
            {
                "datacenterId": 3134,
                "enabled": true,
                "handoutCName": null,
                "name": null,
                "weight": 50.0,
                "servers": [
                    "1.2.3.5"
                ],
                "precedence": 255
            },
            {
                "datacenterId": 3133,
                "enabled": true,
                "handoutCName": null,
                "name": null,
                "weight": 50.0,
                "servers": [
                    "1.2.3.4"
                ],
                "precedence": null
            }
        ]
    },
    "status": {
        "changeId": "eee0c3b4-0e45-4f4b-822c-7dbc60764d18",
        "message": "Change Pending",
        "passingValidation": true,
        "propagationStatus": "PENDING",
        "propagationStatusDate": "2014-04-15T11:30:27.000+0000",
        "links": [
            {
                "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                "rel": "self"
            }
        ]
    }
}
`,
			expectedResponse: &CreatePropertyResponse{
				Resource: &Property{
					BalanceByDownloadScore: false,
					HandoutMode:            "normal",
					IPv6:                   false,
					Name:                   "origin",
					ScoreAggregationType:   "mean",
					StaticTTL:              600,
					DynamicTTL:             300,
					Type:                   "weighted-round-robin",
					UseComputedTargets:     false,
					LivenessTests: []LivenessTest{
						{
							DisableNonstandardPortWarning: false,
							HTTPError3xx:                  true,
							HTTPError4xx:                  true,
							HTTPError5xx:                  true,
							HTTPMethod:                    ptr.To("GET"),
							HTTPRequestBody:               ptr.To("TestBody"),
							Pre2023SecurityPosture:        true,
							AlternateCACertificates:       []string{"test1"},
							Name:                          "health-check",
							TestInterval:                  60,
							TestObject:                    "/status",
							TestObjectPort:                80,
							TestObjectProtocol:            "HTTP",
							TestTimeout:                   25.0,
						},
					},
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 3134,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.5"},
							Precedence:   ptr.To(255),
						},
						{
							DatacenterID: 3133,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.4"},
							Precedence:   nil,
						},
					},
					Links: []Link{
						{
							Href: "/config-gtm/v1/domains/example.akadns.net/properties/origin",
							Rel:  "self",
						},
					},
				},
				Status: &ResponseStatus{
					ChangeID:              "eee0c3b4-0e45-4f4b-822c-7dbc60764d18",
					Message:               "Change Pending",
					PassingValidation:     true,
					PropagationStatus:     "PENDING",
					PropagationStatusDate: "2014-04-15T11:30:27.000+0000",
					Links: []Link{
						{
							Href: "/config-gtm/v1/domains/example.akadns.net/status/current",
							Rel:  "self",
						},
					},
				},
			},
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties/origin",
		},
		"validation error - missing precedence for ranked-failover property type": {
			params: CreatePropertyRequest{
				Property: &Property{
					Type:                 "ranked-failover",
					Name:                 "property",
					HandoutMode:          "normal",
					ScoreAggregationType: "mean",
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 1,
							Enabled:      false,
							Precedence:   nil,
						},
						{
							DatacenterID: 2,
							Enabled:      false,
							Precedence:   nil,
						},
					},
				},
			},
			withError: true,
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "TrafficTargets: property cannot have multiple primary traffic targets (targets with lowest precedence)")
			},
		},
		"validation error - precedence value over the limit": {
			params: CreatePropertyRequest{
				Property: &Property{
					Type:                 "ranked-failover",
					Name:                 "property",
					HandoutMode:          "normal",
					ScoreAggregationType: "mean",
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 1,
							Enabled:      false,
							Precedence:   ptr.To(256),
						},
					},
				},
			},
			withError: true,
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "TrafficTargets: 'Precedence' value has to be between 0 and 255")
			},
		},
		"500 internal server error": {
			params: CreatePropertyRequest{
				Property: &Property{
					Name:                 "testName",
					HandoutMode:          "normal",
					ScoreAggregationType: "mean",
					Type:                 "failover",
				},
				DomainName: "example.akadns.net",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error creating domain"
		}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties/testName",
			withError:    true,
			assertError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating domain",
					StatusCode: http.StatusInternalServerError,
				}
				assert.ErrorIs(t, err, want)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateProperty(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError {
				test.assertError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_UpdateProperty(t *testing.T) {
	tests := map[string]struct {
		params           UpdatePropertyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdatePropertyResponse
		withError        bool
		assertError      func(*testing.T, error)
		headers          http.Header
	}{
		"200 Success": {
			params: UpdatePropertyRequest{
				Property: &Property{
					BalanceByDownloadScore: false,
					HandoutMode:            "normal",
					IPv6:                   false,
					Name:                   "origin",
					ScoreAggregationType:   "mean",
					StaticTTL:              600,
					Type:                   "weighted-round-robin",
					UseComputedTargets:     false,
					LivenessTests: []LivenessTest{
						{
							DisableNonstandardPortWarning: false,
							HTTPError3xx:                  true,
							HTTPError4xx:                  true,
							HTTPError5xx:                  true,
							Name:                          "health-check",
							TestInterval:                  60,
							TestObject:                    "/status",
							TestObjectPort:                80,
							TestObjectProtocol:            "HTTP",
							TestTimeout:                   25.0,
						},
					},
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 3134,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.5"},
							Precedence:   ptr.To(255),
						},
						{
							DatacenterID: 3133,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.4"},
							Precedence:   nil,
						},
					},
				},
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "resource": {
        "backupCName": null,
        "backupIp": null,
        "balanceByDownloadScore": false,
        "cname": null,
        "comments": null,
        "dynamicTTL": 300,
        "failbackDelay": 0,
        "failoverDelay": 0,
        "handoutMode": "normal",
        "healthMax": null,
        "healthMultiplier": null,
        "healthThreshold": null,
        "ipv6": false,
        "lastModified": null,
        "loadImbalancePercentage": null,
        "mapName": null,
        "maxUnreachablePenalty": null,
        "name": "origin",
        "scoreAggregationType": "mean",
        "staticTTL": 600,
        "stickinessBonusConstant": 0,
        "stickinessBonusPercentage": 0,
        "type": "weighted-round-robin",
        "unreachableThreshold": null,
        "useComputedTargets": false,
        "mxRecords": [],
        "links": [
            {
                "href": "/config-gtm/v1/domains/example.akadns.net/properties/origin",
                "rel": "self"
            }
        ],
        "livenessTests": [
            {
                "disableNonstandardPortWarning": false,
                "hostHeader": "foo.example.com",
                "httpError3xx": true,
                "httpError4xx": true,
                "httpError5xx": true,
                "name": "health-check",
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
        "trafficTargets": [
            {
                "datacenterId": 3134,
                "enabled": true,
                "handoutCName": null,
                "name": null,
                "weight": 50.0,
                "servers": [
                    "1.2.3.5"
                ],
                "precedence": 255
            },
            {
                "datacenterId": 3133,
                "enabled": true,
                "handoutCName": null,
                "name": null,
                "weight": 50.0,
                "servers": [
                    "1.2.3.4"
                ],
                "precedence": null
            }
        ]
    },
    "status": {
        "changeId": "eee0c3b4-0e45-4f4b-822c-7dbc60764d18",
        "message": "Change Pending",
        "passingValidation": true,
        "propagationStatus": "PENDING",
        "propagationStatusDate": "2014-04-15T11:30:27.000+0000",
        "links": [
            {
                "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                "rel": "self"
            }
        ]
    }
}
`,
			expectedResponse: &UpdatePropertyResponse{
				Resource: &Property{
					BalanceByDownloadScore: false,
					HandoutMode:            "normal",
					IPv6:                   false,
					Name:                   "origin",
					ScoreAggregationType:   "mean",
					StaticTTL:              600,
					DynamicTTL:             300,
					Type:                   "weighted-round-robin",
					UseComputedTargets:     false,
					LivenessTests: []LivenessTest{
						{
							DisableNonstandardPortWarning: false,
							HTTPError3xx:                  true,
							HTTPError4xx:                  true,
							HTTPError5xx:                  true,
							Name:                          "health-check",
							TestInterval:                  60,
							TestObject:                    "/status",
							TestObjectPort:                80,
							TestObjectProtocol:            "HTTP",
							TestTimeout:                   25.0,
						},
					},
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 3134,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.5"},
							Precedence:   ptr.To(255),
						},
						{
							DatacenterID: 3133,
							Enabled:      true,
							Weight:       50.0,
							Servers:      []string{"1.2.3.4"},
						},
					},
					Links: []Link{
						{
							Href: "/config-gtm/v1/domains/example.akadns.net/properties/origin",
							Rel:  "self",
						},
					},
				},
				Status: &ResponseStatus{
					ChangeID:              "eee0c3b4-0e45-4f4b-822c-7dbc60764d18",
					Message:               "Change Pending",
					PassingValidation:     true,
					PropagationStatus:     "PENDING",
					PropagationStatusDate: "2014-04-15T11:30:27.000+0000",
					Links: []Link{
						{
							Href: "/config-gtm/v1/domains/example.akadns.net/status/current",
							Rel:  "self",
						},
					},
				},
			},
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties/origin",
		},
		"validation error - missing precedence for ranked-failover property type": {
			params: UpdatePropertyRequest{
				Property: &Property{
					Type:                 "ranked-failover",
					Name:                 "property",
					HandoutMode:          "normal",
					ScoreAggregationType: "mean",
					TrafficTargets: []TrafficTarget{
						{
							DatacenterID: 1,
							Enabled:      false,
							Precedence:   nil,
						},
						{
							DatacenterID: 2,
							Enabled:      false,
							Precedence:   nil,
						},
					},
				},
				DomainName: "example.akadns.net",
			},
			withError: true,
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "TrafficTargets: property cannot have multiple primary traffic targets (targets with lowest precedence)")
			},
		},
		"validation error - no traffic targets": {
			params: UpdatePropertyRequest{
				Property: &Property{
					Type:                 "ranked-failover",
					Name:                 "property",
					HandoutMode:          "normal",
					ScoreAggregationType: "mean",
				},
				DomainName: "example.akadns.net",
			},
			withError: true,
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "TrafficTargets: no traffic targets are enabled")
			},
		},
		"500 internal server error": {
			params: UpdatePropertyRequest{
				Property: &Property{
					Name:                 "testName",
					HandoutMode:          "normal",
					ScoreAggregationType: "mean",
					Type:                 "failover",
				},
				DomainName: "example.akadns.net",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties/testName",
			withError:    true,
			assertError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating zone",
					StatusCode: http.StatusInternalServerError,
				}
				assert.ErrorIs(t, err, want)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateProperty(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError {
				test.assertError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_DeleteProperty(t *testing.T) {
	var result DeletePropertyResponse
	var req Property

	respData, err := loadTestData("TestGTM_CreateProperty.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateProperty.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           DeletePropertyRequest
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *DeletePropertyResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: DeletePropertyRequest{
				PropertyName: "www",
				DomainName:   "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/properties/www",
		},
		"500 internal server error": {
			params: DeletePropertyRequest{
				PropertyName: "www",
				DomainName:   "example.akadns.net",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties/www",
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
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write(test.responseBody)
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteProperty(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
