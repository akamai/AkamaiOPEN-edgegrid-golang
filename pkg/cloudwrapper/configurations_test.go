package cloudwrapper

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfiguration(t *testing.T) {
	tests := map[string]struct {
		params           GetConfigurationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Configuration
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetConfigurationRequest{
				ConfigID: 1,
			},
			responseStatus: 200,
			responseBody: `
{
    "configId": 1,
    "configName": "TestConfigName",
    "contractId": "TestContractID",
    "propertyIds": [
        "321",
		"654"
    ],
    "comments": "TestComments",
    "status": "ACTIVE",
    "retainIdleObjects": false,
    "locations": [
        {
            "trafficTypeId": 1,
            "comments": "TestComments",
            "capacity": {
                "value": 1,
                "unit": "GB"
            },
			"mapName": "cw-s-use"
        },
		{
            "trafficTypeId": 2,
            "comments": "TestComments",
            "capacity": {
                "value": 2,
                "unit": "TB"
            },
			"mapName": "cw-s-use"
        }
    ],
    "multiCdnSettings": {
        "origins": [
            {
                "originId": "TestOriginID",
                "hostname": "TestHostname",
						"propertyId": 321
            },
			{
                "originId": "TestOriginID2",
                "hostname": "TestHostname",
                "propertyId": 654
            }
        ],
        "cdns": [
            {
                "cdnCode": "TestCDNCode",
                "enabled": true,
                "cdnAuthKeys": [
                    {
                        "authKeyName": "TestAuthKeyName"
                    }
                ],
                "ipAclCidrs": [],
                "httpsOnly": false
            },
			{
                "cdnCode": "TestCDNCode",
                "enabled": false,
                "cdnAuthKeys": [
                    {
                        "authKeyName": "TestAuthKeyName"
                    },
					{
                        "authKeyName": "TestAuthKeyName2"
                    }
                ],
                "ipAclCidrs": [
					"test1",
					"test2"
				],
                "httpsOnly": true
            }
        ],
        "dataStreams": {
            "enabled": true,
            "dataStreamIds": [
				11,
				22
			],
			"samplingRate": 999
        },
        "bocc": {
            "enabled": false,
			"conditionalSamplingFrequency": "ONE_TENTH",
			"forwardType": "ORIGIN_AND_MIDGRESS",
			"requestType": "EDGE_ONLY",
			"samplingFrequency": "ZERO"
        },
        "enableSoftAlerts": true
    },
    "capacityAlertsThreshold": 75,
    "notificationEmails": [
        "test@akamai.com"
    ],
    "lastUpdatedDate": "2023-05-10T09:55:37.000Z",
    "lastUpdatedBy": "user",
    "lastActivatedDate": "2023-05-10T10:14:49.379Z",
    "lastActivatedBy": "user"
}`,
			expectedPath: "/cloud-wrapper/v1/configurations/1",
			expectedResponse: &Configuration{
				CapacityAlertsThreshold: ptr.To(75),
				Comments:                "TestComments",
				ContractID:              "TestContractID",
				ConfigID:                1,
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestComments",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 1,
						},
						MapName: "cw-s-use",
					},
					{
						Comments:      "TestComments",
						TrafficTypeID: 2,
						Capacity: Capacity{
							Unit:  "TB",
							Value: 2,
						},
						MapName: "cw-s-use",
					},
				},
				MultiCDNSettings: &MultiCDNSettings{
					BOCC: &BOCC{
						ConditionalSamplingFrequency: SamplingFrequencyOneTenth,
						Enabled:                      false,
						ForwardType:                  ForwardTypeOriginAndMidgress,
						RequestType:                  RequestTypeEdgeOnly,
						SamplingFrequency:            SamplingFrequencyZero,
					},
					CDNs: []CDN{
						{
							CDNAuthKeys: []CDNAuthKey{
								{
									AuthKeyName: "TestAuthKeyName",
								},
							},
							CDNCode:    "TestCDNCode",
							Enabled:    true,
							HTTPSOnly:  false,
							IPACLCIDRs: []string{},
						},
						{
							CDNAuthKeys: []CDNAuthKey{
								{
									AuthKeyName: "TestAuthKeyName",
								},
								{
									AuthKeyName: "TestAuthKeyName2",
								},
							},
							CDNCode:   "TestCDNCode",
							Enabled:   false,
							HTTPSOnly: true,
							IPACLCIDRs: []string{
								"test1",
								"test2",
							},
						},
					},
					DataStreams: &DataStreams{
						DataStreamIDs: []int64{11, 22},
						Enabled:       true,
						SamplingRate:  ptr.To(999),
					},
					EnableSoftAlerts: true,
					Origins: []Origin{
						{
							Hostname:   "TestHostname",
							OriginID:   "TestOriginID",
							PropertyID: 321,
						},
						{
							Hostname:   "TestHostname",
							OriginID:   "TestOriginID2",
							PropertyID: 654,
						},
					},
				},
				Status:             "ACTIVE",
				ConfigName:         "TestConfigName",
				LastUpdatedBy:      "user",
				LastUpdatedDate:    "2023-05-10T09:55:37.000Z",
				LastActivatedBy:    ptr.To("user"),
				LastActivatedDate:  ptr.To("2023-05-10T10:14:49.379Z"),
				NotificationEmails: []string{"test@akamai.com"},
				PropertyIDs: []string{
					"321",
					"654",
				},
				RetainIdleObjects: false,
			},
		},
		"200 OK - minimal": {
			params: GetConfigurationRequest{
				ConfigID: 1,
			},
			responseStatus: 200,
			responseBody: `
{
   "configId":1,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestComments",
   "status":"ACTIVE",
   "retainIdleObjects":false,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestComments",
         "capacity":{
            "value":1,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":null,
   "capacityAlertsThreshold":null,
   "notificationEmails":[],
   "lastUpdatedDate":"2023-05-10T09:55:37.000Z",
   "lastUpdatedBy":"user",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			expectedPath: "/cloud-wrapper/v1/configurations/1",
			expectedResponse: &Configuration{
				Comments:   "TestComments",
				ContractID: "TestContractID",
				ConfigID:   1,
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestComments",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 1,
						},
						MapName: "cw-s-use",
					},
				},
				Status:             "ACTIVE",
				ConfigName:         "TestConfigName",
				LastUpdatedBy:      "user",
				LastUpdatedDate:    "2023-05-10T09:55:37.000Z",
				NotificationEmails: []string{},
				PropertyIDs: []string{
					"123",
				},
				RetainIdleObjects: false,
			},
		},
		"missing required params - validation error": {
			params: GetConfigurationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get configuration: struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: GetConfigurationRequest{
				ConfigID: 3,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/cloudwrapper/error-types/cloudwrapper-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/cloudwrapper/error-instances/abc",
    "status": 500
}`,
			expectedPath: "/cloud-wrapper/v1/configurations/3",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/cloudwrapper/error-types/cloudwrapper-server-error",
					Title:    "An unexpected error has occurred.",
					Detail:   "Error processing request",
					Instance: "/cloudwrapper/error-instances/abc",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
			result, err := client.GetConfiguration(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListConfigurations(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListConfigurationsResponse
		withError        error
	}{
		"200 OK": {
			responseStatus: 200,
			responseBody: `
{
   "configurations":[
      {
         "configId":1,
         "configName":"testcloudwrapper",
         "contractId":"testContract",
         "propertyIds":[
            "11"
         ],
         "comments":"testComments",
         "status":"ACTIVE",
         "retainIdleObjects":false,
         "locations":[
            {
               "trafficTypeId":1,
               "comments":"usageNotes",
               "capacity":{
                  "value":1,
                  "unit":"GB"
               },
			   "mapName": "cw-s-use"
            }
         ],
         "multiCdnSettings":null,
         "capacityAlertsThreshold":75,
         "notificationEmails":[
            "user@akamai.com"
         ],
         "lastUpdatedDate":"2023-05-10T09:55:37.000Z",
         "lastUpdatedBy":"user",
         "lastActivatedDate":"2023-05-10T10:14:49.379Z",
         "lastActivatedBy":"user"
      },
      {
         "configId":2,
         "configName":"testcloudwrappermcdn",
         "contractId":"testContract2",
         "propertyIds":[
            "22"
         ],
         "comments":"mcdn",
         "status":"ACTIVE",
         "retainIdleObjects":false,
         "locations":[
            {
               "trafficTypeId":2,
               "comments":"mcdn",
               "capacity":{
                  "value":2,
                  "unit":"TB"
               },
			   "mapName": "cw-s-use"
            }
         ],
         "multiCdnSettings":{
            "origins":[
               {
                  "originId":"testOrigin",
                  "hostname":"hostname.example.com",
                  "propertyId":222
               }
            ],
            "cdns":[
               {
                  "cdnCode":"testCode2",
                  "enabled":true,
                  "cdnAuthKeys":[
                     {
                        "authKeyName":"authKeyTest2"
                     }
                  ],
                  "ipAclCidrs":[
                     "2.2.2.2/22"
                  ],
                  "httpsOnly":true
               }
            ],
            "dataStreams":{
               "enabled":false,
               "dataStreamIds":[
                  2
               ]
            },
            "bocc":{
               "enabled":false
            },
            "enableSoftAlerts":true
         },
         "capacityAlertsThreshold":75,
         "notificationEmails":[
            "user@akamai.com"
         ],
         "lastUpdatedDate":"2023-05-10T09:55:37.000Z",
         "lastUpdatedBy":"user",
         "lastActivatedDate":"2023-05-10T10:14:49.379Z",
         "lastActivatedBy":"user"
      }
   ]
}`,
			expectedPath: "/cloud-wrapper/v1/configurations",
			expectedResponse: &ListConfigurationsResponse{
				Configurations: []Configuration{
					{
						CapacityAlertsThreshold: ptr.To(75),
						Comments:                "testComments",
						ContractID:              "testContract",
						ConfigID:                1,
						Locations: []ConfigLocationResp{
							{
								Comments:      "usageNotes",
								TrafficTypeID: 1,
								Capacity: Capacity{
									Unit:  "GB",
									Value: 1,
								},
								MapName: "cw-s-use",
							},
						},
						Status:             "ACTIVE",
						ConfigName:         "testcloudwrapper",
						LastUpdatedBy:      "user",
						LastUpdatedDate:    "2023-05-10T09:55:37.000Z",
						LastActivatedBy:    ptr.To("user"),
						LastActivatedDate:  ptr.To("2023-05-10T10:14:49.379Z"),
						NotificationEmails: []string{"user@akamai.com"},
						PropertyIDs: []string{
							"11",
						},
						RetainIdleObjects: false,
					},
					{
						CapacityAlertsThreshold: ptr.To(75),
						Comments:                "mcdn",
						ContractID:              "testContract2",
						ConfigID:                2,
						Locations: []ConfigLocationResp{
							{
								Comments:      "mcdn",
								TrafficTypeID: 2,
								Capacity: Capacity{
									Unit:  "TB",
									Value: 2,
								},
								MapName: "cw-s-use",
							},
						},
						MultiCDNSettings: &MultiCDNSettings{
							BOCC: &BOCC{
								Enabled: false,
							},
							CDNs: []CDN{
								{
									CDNAuthKeys: []CDNAuthKey{
										{
											AuthKeyName: "authKeyTest2",
										},
									},
									CDNCode:   "testCode2",
									Enabled:   true,
									HTTPSOnly: true,
									IPACLCIDRs: []string{
										"2.2.2.2/22",
									},
								},
							},
							DataStreams: &DataStreams{
								DataStreamIDs: []int64{2},
								Enabled:       false,
							},
							EnableSoftAlerts: true,
							Origins: []Origin{
								{
									Hostname:   "hostname.example.com",
									OriginID:   "testOrigin",
									PropertyID: 222,
								},
							},
						},
						Status:             "ACTIVE",
						ConfigName:         "testcloudwrappermcdn",
						LastUpdatedBy:      "user",
						LastUpdatedDate:    "2023-05-10T09:55:37.000Z",
						LastActivatedBy:    ptr.To("user"),
						LastActivatedDate:  ptr.To("2023-05-10T10:14:49.379Z"),
						NotificationEmails: []string{"user@akamai.com"},
						PropertyIDs:        []string{"22"},
						RetainIdleObjects:  false,
					},
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/cloudwrapper/error-types/cloudwrapper-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/cloudwrapper/error-instances/abc",
    "status": 500
}`,
			expectedPath: "/cloud-wrapper/v1/configurations",
			withError: &Error{
				Type:     "/cloudwrapper/error-types/cloudwrapper-server-error",
				Title:    "An unexpected error has occurred.",
				Detail:   "Error processing request",
				Instance: "/cloudwrapper/error-instances/abc",
				Status:   500,
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
			result, err := client.ListConfigurations(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreateConfiguration(t *testing.T) {
	tests := map[string]struct {
		params              CreateConfigurationRequest
		expectedRequestBody string
		expectedPath        string
		responseStatus      int
		responseBody        string
		expectedResponse    *Configuration
		withError           func(*testing.T, error)
	}{
		"200 OK - minimal": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  UnitGB,
								Value: 1,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"123"},
				},
			},
			expectedRequestBody: `
{
   "locations":[
      {
         "capacity":{
            "value":1,
            "unit":"GB"
         },
         "comments":"TestComments",
         "trafficTypeId":1
      }
   ],
   "propertyIds":[
      "123"
   ],
   "contractId":"TestContractID",
   "comments":"TestComments",
   "configName":"TestConfigName"
}`,
			expectedPath:   "/cloud-wrapper/v1/configurations?activate=false",
			responseStatus: 201,
			responseBody: `
{
   "configId":111,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestComments",
   "status":"IN_PROGRESS",
   "retainIdleObjects":false,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestComments",
         "capacity":{
            "value":1,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":null,
   "capacityAlertsThreshold":50,
   "notificationEmails":[
      
   ],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			expectedResponse: &Configuration{
				CapacityAlertsThreshold: ptr.To(50),
				Comments:                "TestComments",
				ContractID:              "TestContractID",
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestComments",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  UnitGB,
							Value: 1,
						},
						MapName: "cw-s-use",
					},
				},
				Status:             StatusInProgress,
				ConfigName:         "TestConfigName",
				LastUpdatedBy:      "johndoe",
				LastUpdatedDate:    "2022-06-10T13:21:14.488Z",
				NotificationEmails: []string{},
				PropertyIDs:        []string{"123"},
				ConfigID:           111,
			},
		},
		"200 OK - minimal with activate query param": {
			params: CreateConfigurationRequest{
				Activate: true,
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"123"},
				},
			},
			expectedRequestBody: `
{
   "locations":[
      {
         "capacity":{
            "value":1,
            "unit":"GB"
         },
         "comments":"TestComments",
         "trafficTypeId":1
      }
   ],
   "propertyIds":[
      "123"
   ],
   "contractId":"TestContractID",
   "comments":"TestComments",
   "configName":"TestConfigName"
}`,
			expectedPath:   "/cloud-wrapper/v1/configurations?activate=true",
			responseStatus: 201,
			responseBody: `
{
   "configId":111,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestComments",
   "status":"IN_PROGRESS",
   "retainIdleObjects":false,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestComments",
         "capacity":{
            "value":1,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":null,
   "capacityAlertsThreshold":50,
   "notificationEmails":[
      
   ],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			expectedResponse: &Configuration{
				CapacityAlertsThreshold: ptr.To(50),
				Comments:                "TestComments",
				ContractID:              "TestContractID",
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestComments",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 1,
						},
						MapName: "cw-s-use",
					},
				},
				Status:             StatusInProgress,
				ConfigName:         "TestConfigName",
				LastUpdatedBy:      "johndoe",
				LastUpdatedDate:    "2022-06-10T13:21:14.488Z",
				NotificationEmails: []string{},
				PropertyIDs:        []string{"123"},
				ConfigID:           111,
			},
		},
		"200 OK - minimal MultiCDNSettings": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							Enabled: false,
						},
						CDNs: []CDN{
							{
								CDNAuthKeys: []CDNAuthKey{
									{
										AuthKeyName: "TestAuthKeyName",
										ExpiryDate:  "TestExpiryDate",
										HeaderName:  "TestHeaderName",
										Secret:      "testtesttesttesttesttest",
									},
								},
								CDNCode: "TestCDNCode",
								Enabled: true,
							},
						},
						DataStreams: &DataStreams{
							Enabled: false,
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 123,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"123"},
				},
			},
			expectedRequestBody: `
{
   "locations":[
      {
         "capacity":{
            "value":1,
            "unit":"GB"
         },
         "comments":"TestComments",
         "trafficTypeId":1
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName",
                  "headerName":"TestHeaderName",
                  "secret":"testtesttesttesttesttest",
                  "expiryDate":"TestExpiryDate"
               }
            ]
         }
      ],
      "dataStreams":{
         "enabled":false
      },
      "bocc":{
         "enabled":false
      }
   },
   "propertyIds":[
      "123"
   ],
   "contractId":"TestContractID",
   "comments":"TestComments",
   "configName":"TestConfigName"
}`,
			expectedPath:   "/cloud-wrapper/v1/configurations?activate=false",
			responseStatus: 201,
			responseBody: `
{
   "configId":111,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestComments",
   "status":"IN_PROGRESS",
   "retainIdleObjects":false,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestComments",
         "capacity":{
            "value":1,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName",
                  "headerName":"TestHeaderName",
                  "secret":"testtesttesttesttesttest",
                  "expiryDate":"TestExpiryDate"
               }
            ],
            "ipAclCidrs":[],
            "httpsOnly":false
         }
      ],
      "dataStreams":{
         "enabled":false
      },
      "bocc":{
         "enabled":false
      },
      "enableSoftAlerts":false
   },
   "capacityAlertsThreshold":null,
   "notificationEmails":[],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			expectedResponse: &Configuration{
				CapacityAlertsThreshold: nil,
				Comments:                "TestComments",
				ContractID:              "TestContractID",
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestComments",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 1,
						},
						MapName: "cw-s-use",
					},
				},
				MultiCDNSettings: &MultiCDNSettings{
					BOCC: &BOCC{
						Enabled: false,
					},
					CDNs: []CDN{
						{
							CDNAuthKeys: []CDNAuthKey{
								{
									AuthKeyName: "TestAuthKeyName",
									ExpiryDate:  "TestExpiryDate",
									HeaderName:  "TestHeaderName",
									Secret:      "testtesttesttesttesttest",
								},
							},
							CDNCode:    "TestCDNCode",
							Enabled:    true,
							IPACLCIDRs: []string{},
						},
					},
					DataStreams: &DataStreams{
						Enabled: false,
					},
					Origins: []Origin{
						{
							Hostname:   "TestHostname",
							OriginID:   "TestOriginID",
							PropertyID: 123,
						},
					},
					EnableSoftAlerts: false,
				},
				RetainIdleObjects:  false,
				Status:             StatusInProgress,
				ConfigName:         "TestConfigName",
				LastUpdatedBy:      "johndoe",
				LastUpdatedDate:    "2022-06-10T13:21:14.488Z",
				NotificationEmails: []string{},
				PropertyIDs:        []string{"123"},
				ConfigID:           111,
			},
		},
		"200 OK - full MultiCDNSettings": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					CapacityAlertsThreshold: ptr.To(70),
					Comments:                "TestComments",
					ContractID:              "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
						{
							Comments:      "TestComments2",
							TrafficTypeID: 2,
							Capacity: Capacity{
								Unit:  "TB",
								Value: 2,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							ConditionalSamplingFrequency: SamplingFrequencyZero,
							Enabled:                      true,
							ForwardType:                  ForwardTypeOriginAndMidgress,
							RequestType:                  RequestTypeEdgeAndMidgress,
							SamplingFrequency:            SamplingFrequencyZero,
						},
						CDNs: []CDN{
							{
								CDNAuthKeys: []CDNAuthKey{
									{
										AuthKeyName: "TestAuthKeyName",
										ExpiryDate:  "TestExpiryDate",
										HeaderName:  "TestHeaderName",
										Secret:      "testtesttesttesttesttest",
									},
								},
								CDNCode:   "TestCDNCode",
								Enabled:   true,
								HTTPSOnly: true,
							},
							{
								CDNCode:   "TestCDNCode",
								Enabled:   true,
								HTTPSOnly: true,
								IPACLCIDRs: []string{
									"1.1.1.1/1",
								},
							},
						},
						DataStreams: &DataStreams{
							DataStreamIDs: []int64{1},
							Enabled:       true,
							SamplingRate:  ptr.To(10),
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 123,
							},
							{
								Hostname:   "TestHostname2",
								OriginID:   "TestOriginID2",
								PropertyID: 1234,
							},
						},
						EnableSoftAlerts: true,
					},
					ConfigName: "TestConfigName",
					NotificationEmails: []string{
						"test@test.com",
					},
					PropertyIDs:       []string{"123"},
					RetainIdleObjects: true,
				},
			},
			expectedRequestBody: `
{
   "capacityAlertsThreshold":70,
   "locations":[
      {
         "capacity":{
            "value":1,
            "unit":"GB"
         },
         "comments":"TestComments",
         "trafficTypeId":1
      },
	  {
         "capacity":{
            "value":2,
            "unit":"TB"
         },
         "comments":"TestComments2",
         "trafficTypeId":2
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         },
		 {
            "originId":"TestOriginID2",
            "hostname":"TestHostname2",
            "propertyId":1234
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName",
                  "headerName":"TestHeaderName",
                  "secret":"testtesttesttesttesttest",
                  "expiryDate":"TestExpiryDate"
               }
            ],
            "httpsOnly":true
         },
		 {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "httpsOnly":true,
			"ipAclCidrs": [
				"1.1.1.1/1"
			]
         }
      ],
      "dataStreams":{
         "enabled":true,
		 "dataStreamIds": [
			1
		 ],
		 "samplingRate": 10
      },
      "bocc":{
         "enabled":true,
		 "conditionalSamplingFrequency": "ZERO",
		 "forwardType": "ORIGIN_AND_MIDGRESS",
		 "requestType": "EDGE_AND_MIDGRESS",
		 "samplingFrequency": "ZERO"
      },
	  "enableSoftAlerts":true
   },
   "propertyIds":[
      "123"
   ],
   "notificationEmails": [
		"test@test.com"
   ],
   "retainIdleObjects":true,
   "contractId":"TestContractID",
   "comments":"TestComments",
   "configName":"TestConfigName"
}
`,
			expectedPath:   "/cloud-wrapper/v1/configurations?activate=false",
			responseStatus: 201,
			responseBody: `
{
   "configId":111,
   "capacityAlertsThreshold": 70,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestComments",
   "status":"IN_PROGRESS",
   "retainIdleObjects":true,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestComments",
         "capacity":{
            "value":1,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      },
	  {
         "trafficTypeId":2,
         "comments":"TestComments2",
         "capacity":{
            "value":2,
            "unit":"TB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         },
		 {
            "originId":"TestOriginID2",
            "hostname":"TestHostname2",
            "propertyId":1234
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName",
                  "headerName":"TestHeaderName",
                  "secret":"testtesttesttesttesttest",
                  "expiryDate":"TestExpiryDate"
               }
            ],
            "ipAclCidrs":[],
            "httpsOnly":true
         },
		 {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "httpsOnly":true,
			"ipAclCidrs": [
				"1.1.1.1/1"
			]
         }
      ],
      "dataStreams":{
         "enabled":true,
		 "dataStreamIds": [
			1
		 ],
		 "samplingRate": 10
      },
      "bocc":{
         "enabled":true,
		 "conditionalSamplingFrequency": "ZERO",
		 "forwardType": "ORIGIN_AND_MIDGRESS",
		 "requestType": "EDGE_AND_MIDGRESS",
		 "samplingFrequency": "ZERO"
      },
      "enableSoftAlerts": true
   },
   "notificationEmails":[
      "test@test.com"
   ],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			expectedResponse: &Configuration{
				CapacityAlertsThreshold: ptr.To(70),
				Comments:                "TestComments",
				ContractID:              "TestContractID",
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestComments",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 1,
						},
						MapName: "cw-s-use",
					},
					{
						Comments:      "TestComments2",
						TrafficTypeID: 2,
						Capacity: Capacity{
							Unit:  "TB",
							Value: 2,
						},
						MapName: "cw-s-use",
					},
				},
				MultiCDNSettings: &MultiCDNSettings{
					BOCC: &BOCC{
						ConditionalSamplingFrequency: SamplingFrequencyZero,
						Enabled:                      true,
						ForwardType:                  ForwardTypeOriginAndMidgress,
						RequestType:                  RequestTypeEdgeAndMidgress,
						SamplingFrequency:            SamplingFrequencyZero,
					},
					CDNs: []CDN{
						{
							CDNAuthKeys: []CDNAuthKey{
								{
									AuthKeyName: "TestAuthKeyName",
									ExpiryDate:  "TestExpiryDate",
									HeaderName:  "TestHeaderName",
									Secret:      "testtesttesttesttesttest",
								},
							},
							CDNCode:    "TestCDNCode",
							Enabled:    true,
							IPACLCIDRs: []string{},
							HTTPSOnly:  true,
						},
						{
							CDNCode:    "TestCDNCode",
							Enabled:    true,
							IPACLCIDRs: []string{"1.1.1.1/1"},
							HTTPSOnly:  true,
						},
					},
					DataStreams: &DataStreams{
						DataStreamIDs: []int64{1},
						Enabled:       true,
						SamplingRate:  ptr.To(10),
					},
					Origins: []Origin{
						{
							Hostname:   "TestHostname",
							OriginID:   "TestOriginID",
							PropertyID: 123,
						},
						{
							Hostname:   "TestHostname2",
							OriginID:   "TestOriginID2",
							PropertyID: 1234,
						},
					},
					EnableSoftAlerts: true,
				},
				RetainIdleObjects:  true,
				Status:             StatusInProgress,
				ConfigName:         "TestConfigName",
				LastUpdatedBy:      "johndoe",
				LastUpdatedDate:    "2022-06-10T13:21:14.488Z",
				NotificationEmails: []string{"test@test.com"},
				PropertyIDs:        []string{"123"},
				ConfigID:           111,
			},
		},
		"200 OK - BOCC struct fields default values": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 10,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{},
						CDNs: []CDN{
							{
								CDNAuthKeys: []CDNAuthKey{
									{
										AuthKeyName: "TestAuthKeyName",
									},
								},
								CDNCode: "TestCDNCode",
								Enabled: true,
							},
						},
						DataStreams: &DataStreams{
							Enabled: true,
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 123,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"123"},
				},
			},
			expectedRequestBody: `
{
   "locations":[
      {
         "capacity":{
            "value":10,
            "unit":"GB"
         },
         "comments":"TestComments",
         "trafficTypeId":1
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName"
               }
            ]
         }
      ],
      "dataStreams":{
         "enabled":true
      },
      "bocc":{
         "enabled":false
      }
   },
   "propertyIds":[
      "123"
   ],
   "contractId":"TestContractID",
   "comments":"TestComments",
   "configName":"TestConfigName"
}`,
			expectedPath: "/cloud-wrapper/v1/configurations?activate=false",
			responseBody: `
{
   "configId":111,
   "capacityAlertsThreshold":null,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestComments",
   "status":"IN_PROGRESS",
   "retainIdleObjects":false,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestComments",
         "capacity":{
            "value":10,
            "unit":"GB"
         }
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName"
               }
            ],
            "ipAclCidrs":[]
         }
      ],
      "dataStreams":{
         "enabled":true
      },
      "bocc":{
         "enabled":false
      },
      "enableSoftAlerts":false
   },
   "notificationEmails":[
      
   ],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			responseStatus: 201,
			expectedResponse: &Configuration{
				Comments:   "TestComments",
				ContractID: "TestContractID",
				Status:     StatusInProgress,
				ConfigID:   111,
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestComments",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 10,
						},
					},
				},
				MultiCDNSettings: &MultiCDNSettings{
					BOCC: &BOCC{
						Enabled: false,
					},
					CDNs: []CDN{
						{
							CDNAuthKeys: []CDNAuthKey{
								{
									AuthKeyName: "TestAuthKeyName",
								},
							},
							CDNCode:    "TestCDNCode",
							Enabled:    true,
							IPACLCIDRs: []string{},
						},
					},
					DataStreams: &DataStreams{
						Enabled: true,
					},
					Origins: []Origin{
						{
							Hostname:   "TestHostname",
							OriginID:   "TestOriginID",
							PropertyID: 123,
						},
					},
				},
				ConfigName:         "TestConfigName",
				PropertyIDs:        []string{"123"},
				NotificationEmails: []string{},
				LastUpdatedBy:      "johndoe",
				LastUpdatedDate:    "2022-06-10T13:21:14.488Z",
			},
		},
		"200 OK - DataStreams struct fields default values": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 10,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{},
						CDNs: []CDN{
							{
								CDNAuthKeys: []CDNAuthKey{
									{
										AuthKeyName: "TestAuthKeyName",
									},
								},
								CDNCode: "TestCDNCode",
								Enabled: true,
							},
						},
						DataStreams: &DataStreams{},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 123,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"123"},
				},
			},
			expectedRequestBody: `
{
   "locations":[
      {
         "capacity":{
            "value":10,
            "unit":"GB"
         },
         "comments":"TestComments",
         "trafficTypeId":1
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName"
               }
            ]
         }
      ],
      "dataStreams":{
         "enabled":false
      },
      "bocc":{
         "enabled":false
      }
   },
   "propertyIds":[
      "123"
   ],
   "contractId":"TestContractID",
   "comments":"TestComments",
   "configName":"TestConfigName"
}`,
			expectedPath: "/cloud-wrapper/v1/configurations?activate=false",
			responseBody: `
{
   "configId":111,
   "capacityAlertsThreshold":null,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestComments",
   "status":"IN_PROGRESS",
   "retainIdleObjects":false,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestComments",
         "capacity":{
            "value":10,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName"
               }
            ],
            "ipAclCidrs":[]
         }
      ],
      "dataStreams":{
         "enabled":false,
		 "dataStreamsIds": []
      },
      "bocc":{
         "enabled":false
      },
      "enableSoftAlerts":false
   },
   "notificationEmails":[],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			responseStatus: 201,
			expectedResponse: &Configuration{
				Comments:   "TestComments",
				ContractID: "TestContractID",
				Status:     StatusInProgress,
				ConfigID:   111,
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestComments",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 10,
						},
						MapName: "cw-s-use",
					},
				},
				MultiCDNSettings: &MultiCDNSettings{
					BOCC: &BOCC{
						Enabled: false,
					},
					CDNs: []CDN{
						{
							CDNAuthKeys: []CDNAuthKey{
								{
									AuthKeyName: "TestAuthKeyName",
								},
							},
							CDNCode:    "TestCDNCode",
							Enabled:    true,
							IPACLCIDRs: []string{},
						},
					},
					DataStreams: &DataStreams{
						Enabled: false,
					},
					Origins: []Origin{
						{
							Hostname:   "TestHostname",
							OriginID:   "TestOriginID",
							PropertyID: 123,
						},
					},
				},
				ConfigName:         "TestConfigName",
				PropertyIDs:        []string{"123"},
				NotificationEmails: []string{},
				LastUpdatedBy:      "johndoe",
				LastUpdatedDate:    "2022-06-10T13:21:14.488Z",
			},
		},
		"missing required params: comments, configName, contractID, locations and propertyIDs - validation error": {
			params: CreateConfigurationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create configuration: struct validation: Body: {\n\tComments: cannot be blank\n\tConfigName: cannot be blank\n\tContractID: cannot be blank\n\tLocations: cannot be blank\n\tPropertyIDs: cannot be blank\n}", err.Error())
			},
		},
		"missing required params - location fields": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "",
							TrafficTypeID: 0,
							Capacity:      Capacity{},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"1"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create configuration: struct validation: Body: {\n\tLocations[0]: {\n\t\tCapacity: {\n\t\t\tUnit: cannot be blank\n\t\t\tValue: cannot be blank\n\t\t}\n\t\tComments: cannot be blank\n\t\tTrafficTypeID: cannot be blank\n\t}\n}", err.Error())
			},
		},
		"missing required params - multiCDN fields": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 5,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 10,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{},
					ConfigName:       "TestConfigName",
					PropertyIDs:      []string{"1"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create configuration: struct validation: Body: {\n\tMultiCDNSettings: {\n\t\tBOCC: cannot be blank\n\t\tCDNs: cannot be blank\n\t\tDataStreams: cannot be blank\n\t\tOrigins: cannot be blank\n\t}\n}", err.Error())
			},
		},
		"missing required params - BOCC struct fields when enabled": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 5,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 10,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							Enabled: true,
						},
						CDNs: []CDN{
							{
								CDNAuthKeys: []CDNAuthKey{
									{AuthKeyName: "TestAuthKeyName"},
								},
								CDNCode: "TestCDNCode",
								Enabled: true,
							},
						},
						DataStreams: &DataStreams{
							Enabled: true,
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 1,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"1"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create configuration: struct validation: Body: {\n\tMultiCDNSettings: {\n\t\tBOCC: {\n\t\t\tConditionalSamplingFrequency: cannot be blank\n\t\t\tForwardType: cannot be blank\n\t\t\tRequestType: cannot be blank\n\t\t\tSamplingFrequency: cannot be blank\n\t\t}\n\t}\n}", err.Error())
			},
		},
		"missing required params - Origin struct fields": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{Comments: "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 5,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 10,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							Enabled: false,
						},
						CDNs: []CDN{
							{
								CDNCode:    "TestCDNCode",
								Enabled:    true,
								IPACLCIDRs: []string{"1.1.1.1/1"},
							},
						},
						DataStreams: &DataStreams{
							Enabled: true,
						},
						Origins: []Origin{
							{},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"1"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create configuration: struct validation: Body: {\n\tMultiCDNSettings: {\n\t\tOrigins[0]: {\n\t\t\tHostname: cannot be blank\n\t\t\tOriginID: cannot be blank\n\t\t\tPropertyID: cannot be blank\n\t\t}\n\t}\n}", err.Error())
			},
		},
		"validation error - at least one CDN must be enabled": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 5,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 10,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							Enabled: false,
						},
						CDNs: []CDN{
							{
								CDNCode:    "TestCDNCode",
								Enabled:    false,
								IPACLCIDRs: []string{"1.1.1.1/1"},
							},
							{
								CDNCode:    "TestCDNCode",
								Enabled:    false,
								IPACLCIDRs: []string{"1.1.1.1/1"},
							},
						},
						DataStreams: &DataStreams{
							Enabled: false,
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 1,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"1"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create configuration: struct validation: Body: {\n\tMultiCDNSettings: {\n\t\tCDNs: at least one of CDNs must be enabled\n\t}\n}", err.Error())
			},
		},
		"validation error - authKeys nor IPACLCIDRs specified": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 5,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 10,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							Enabled: false,
						},
						CDNs: []CDN{
							{
								CDNCode: "TestCDNCode",
								Enabled: false,
							},
						},
						DataStreams: &DataStreams{
							Enabled: false,
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 1,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"1"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create configuration: struct validation: Body: {\n\tMultiCDNSettings: {\n\t\tCDNs: at least one authentication method is required for CDN. Either IP ACL or header authentication must be enabled\n\t}\n}", err.Error())
			},
		},
		"struct fields validations": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					CapacityAlertsThreshold: ptr.To(20),
					Comments:                "TestComments",
					ContractID:              "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 5,
							Capacity: Capacity{
								Unit:  "MB",
								Value: 10,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							ConditionalSamplingFrequency: "a",
							Enabled:                      false,
							ForwardType:                  "a",
							RequestType:                  "a",
							SamplingFrequency:            "a",
						},
						CDNs: []CDN{
							{
								CDNCode:    "TestCDNCode",
								Enabled:    true,
								IPACLCIDRs: []string{"1.1.1.1/1"},
							},
							{
								CDNAuthKeys: []CDNAuthKey{
									{},
								},
								CDNCode: "TestCDNCode",
								Enabled: true,
							},
						},
						DataStreams: &DataStreams{
							DataStreamIDs: []int64{1},
							Enabled:       true,
							SamplingRate:  ptr.To(-10),
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 1,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"1"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create configuration: struct validation: Body: {\n\tCapacityAlertsThreshold: must be no less than 50\n\tLocations[0]: {\n\t\tCapacity: {\n\t\t\tUnit: value 'MB' is invalid. Must be one of: 'GB', 'TB'\n\t\t}\n\t}\n\tMultiCDNSettings: {\n\t\tBOCC: {\n\t\t\tConditionalSamplingFrequency: value 'a' is invalid. Must be one of: 'ZERO', 'ONE_TENTH'\n\t\t\tForwardType: value 'a' is invalid. Must be one of: 'ORIGIN_ONLY', 'MIDGRESS_ONLY', 'ORIGIN_AND_MIDGRESS'\n\t\t\tRequestType: value 'a' is invalid. Must be one of: 'EDGE_ONLY', 'EDGE_AND_MIDGRESS'\n\t\t\tSamplingFrequency: value 'a' is invalid. Must be one of: 'ZERO', 'ONE_TENTH'\n\t\t}\n\t\tCDNs[1]: {\n\t\t\tCDNAuthKeys[0]: {\n\t\t\t\tAuthKeyName: cannot be blank\n\t\t\t}\n\t\t}\n\t\tDataStreams: {\n\t\t\tSamplingRate: must be no less than 1\n\t\t}\n\t}\n}", err.Error())
			},
		},
		"500 internal server error": {
			params: CreateConfigurationRequest{
				Body: CreateConfigurationRequestBody{
					Comments:   "TestComments",
					ContractID: "TestContractID",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
					},
					ConfigName:  "TestConfigName",
					PropertyIDs: []string{"123"},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/cloudwrapper/error-types/cloudwrapper-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/cloudwrapper/error-instances/abc",
    "status": 500
}`,
			expectedPath: "/cloud-wrapper/v1/configurations?activate=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/cloudwrapper/error-types/cloudwrapper-server-error",
					Title:    "An unexpected error has occurred.",
					Detail:   "Error processing request",
					Instance: "/cloudwrapper/error-instances/abc",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateConfiguration(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateConfiguration(t *testing.T) {
	tests := map[string]struct {
		params              UpdateConfigurationRequest
		expectedRequestBody string
		expectedPath        string
		responseStatus      int
		responseBody        string
		expectedResponse    *Configuration
		withError           func(*testing.T, error)
	}{
		"200 OK - minimal": {
			params: UpdateConfigurationRequest{
				ConfigID: 111,
				Body: UpdateConfigurationRequestBody{
					Comments: "TestCommentsUpdated",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestCommentsUpdated",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
					},
					PropertyIDs: []string{
						"123",
					},
				},
			},
			expectedRequestBody: `
{
   "locations":[
      {
         "capacity":{
            "value":1,
            "unit":"GB"
         },
         "comments":"TestCommentsUpdated",
         "trafficTypeId":1
      }
   ],
   "propertyIds":[
      "123"
   ],
   "comments":"TestCommentsUpdated"
}`,
			expectedPath:   "/cloud-wrapper/v1/configurations/111?activate=false",
			responseStatus: 200,
			responseBody: `
{
   "configId":111,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestCommentsUpdated",
   "status":"IN_PROGRESS",
   "retainIdleObjects":false,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestCommentsUpdated",
         "capacity":{
            "value":1,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":null,
   "capacityAlertsThreshold":50,
   "notificationEmails":[
      
   ],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			expectedResponse: &Configuration{
				ConfigID:                111,
				CapacityAlertsThreshold: ptr.To(50),
				Comments:                "TestCommentsUpdated",
				ContractID:              "TestContractID",
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestCommentsUpdated",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 1,
						},
						MapName: "cw-s-use",
					},
				},
				Status:             StatusInProgress,
				ConfigName:         "TestConfigName",
				LastUpdatedBy:      "johndoe",
				LastUpdatedDate:    "2022-06-10T13:21:14.488Z",
				NotificationEmails: []string{},
				PropertyIDs:        []string{"123"},
				RetainIdleObjects:  false,
			},
		},
		"200 OK - minimal MultiCDNSettings": {
			params: UpdateConfigurationRequest{
				ConfigID: 111,
				Body: UpdateConfigurationRequestBody{
					Comments: "TestCommentsUpdated",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestCommentsUpdated",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
					},
					PropertyIDs: []string{
						"123",
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							Enabled: false,
						},
						CDNs: []CDN{
							{
								CDNCode: "TestCDNCode",
								Enabled: true,
								IPACLCIDRs: []string{
									"1.1.1.1/1",
								},
							},
						},
						DataStreams: &DataStreams{
							Enabled: false,
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 123,
							},
						},
					},
				},
			},
			expectedRequestBody: `
{
   "locations":[
      {
         "capacity":{
            "value":1,
            "unit":"GB"
         },
         "comments":"TestCommentsUpdated",
         "trafficTypeId":1
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "ipAclCidrs":[
               "1.1.1.1/1"
            ]
         }
      ],
      "dataStreams":{
         "enabled":false
      },
      "bocc":{
         "enabled":false
      }
   },
   "propertyIds":[
      "123"
   ],
   "comments":"TestCommentsUpdated"
}`,
			expectedPath:   "/cloud-wrapper/v1/configurations/111?activate=false",
			responseStatus: 200,
			responseBody: `
{
   "configId":111,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestCommentsUpdated",
   "status":"IN_PROGRESS",
   "retainIdleObjects":false,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestCommentsUpdated",
         "capacity":{
            "value":1,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[],
            "ipAclCidrs":[
               "1.1.1.1/1"
            ],
            "httpsOnly":false
         }
      ],
      "dataStreams":{
         "enabled":false
      },
      "bocc":{
         "enabled":false
      },
      "enableSoftAlerts":false
   },
   "capacityAlertsThreshold":null,
   "notificationEmails":[],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			expectedResponse: &Configuration{
				ConfigID:   111,
				Comments:   "TestCommentsUpdated",
				ContractID: "TestContractID",
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestCommentsUpdated",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 1,
						},
						MapName: "cw-s-use",
					},
				},
				MultiCDNSettings: &MultiCDNSettings{
					BOCC: &BOCC{
						Enabled: false,
					},
					CDNs: []CDN{
						{
							CDNAuthKeys: []CDNAuthKey{},
							CDNCode:     "TestCDNCode",
							Enabled:     true,
							HTTPSOnly:   false,
							IPACLCIDRs: []string{
								"1.1.1.1/1",
							},
						},
					},
					DataStreams: &DataStreams{
						Enabled: false,
					},
					EnableSoftAlerts: false,
					Origins: []Origin{
						{
							Hostname:   "TestHostname",
							OriginID:   "TestOriginID",
							PropertyID: 123,
						},
					},
				},
				Status:             StatusInProgress,
				ConfigName:         "TestConfigName",
				LastUpdatedBy:      "johndoe",
				LastUpdatedDate:    "2022-06-10T13:21:14.488Z",
				NotificationEmails: []string{},
				PropertyIDs:        []string{"123"},
				RetainIdleObjects:  false,
			},
		},
		"200 OK - all fields": {
			params: UpdateConfigurationRequest{
				ConfigID: 111,
				Body: UpdateConfigurationRequestBody{
					CapacityAlertsThreshold: ptr.To(80),
					Comments:                "TestCommentsUpdated",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestCommentsUpdated",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
					},
					MultiCDNSettings: &MultiCDNSettings{
						BOCC: &BOCC{
							ConditionalSamplingFrequency: SamplingFrequencyZero,
							Enabled:                      true,
							ForwardType:                  ForwardTypeOriginAndMidgress,
							RequestType:                  RequestTypeEdgeAndMidgress,
							SamplingFrequency:            SamplingFrequencyZero,
						},
						CDNs: []CDN{
							{
								CDNAuthKeys: []CDNAuthKey{
									{
										AuthKeyName: "TestAuthKeyName",
										ExpiryDate:  "TestExpiryDate",
										HeaderName:  "TestHeaderName",
										Secret:      "TestSecretTestSecret1234",
									},
								},
								CDNCode: "TestCDNCode",
								Enabled: true,
								IPACLCIDRs: []string{
									"1.1.1.1/1",
								},
								HTTPSOnly: true,
							},
						},
						DataStreams: &DataStreams{
							DataStreamIDs: []int64{1},
							Enabled:       true,
							SamplingRate:  ptr.To(10),
						},
						Origins: []Origin{
							{
								Hostname:   "TestHostname",
								OriginID:   "TestOriginID",
								PropertyID: 123,
							},
						},
					},
					NotificationEmails: []string{
						"test@test.com",
					},
					PropertyIDs: []string{
						"123",
					},
					RetainIdleObjects: true,
				},
			},
			expectedRequestBody: `
{
   "capacityAlertsThreshold":80,
   "locations":[
      {
         "capacity":{
            "value":1,
            "unit":"GB"
         },
         "comments":"TestCommentsUpdated",
         "trafficTypeId":1
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName",
                  "expiryDate":"TestExpiryDate",
                  "headerName":"TestHeaderName",
                  "secret":"TestSecretTestSecret1234"
               }
            ],
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "ipAclCidrs":[
               "1.1.1.1/1"
            ],
            "httpsOnly":true
         }
      ],
      "dataStreams":{
         "enabled":true,
         "dataStreamIds":[
            1
         ],
         "samplingRate":10
      },
      "bocc":{
         "enabled":true,
         "conditionalSamplingFrequency":"ZERO",
         "forwardType":"ORIGIN_AND_MIDGRESS",
         "requestType":"EDGE_AND_MIDGRESS",
         "samplingFrequency":"ZERO"
      }
   },
   "propertyIds":[
      "123"
   ],
   "notificationEmails":[
      "test@test.com"
   ],
   "retainIdleObjects":true,
   "comments":"TestCommentsUpdated"
}
`,
			expectedPath:   "/cloud-wrapper/v1/configurations/111?activate=false",
			responseStatus: 200,
			responseBody: `
{
   "configId":111,
   "configName":"TestConfigName",
   "contractId":"TestContractID",
   "propertyIds":[
      "123"
   ],
   "comments":"TestCommentsUpdated",
   "status":"IN_PROGRESS",
   "retainIdleObjects":true,
   "locations":[
      {
         "trafficTypeId":1,
         "comments":"TestCommentsUpdated",
         "capacity":{
            "value":1,
            "unit":"GB"
         },
		 "mapName": "cw-s-use"
      }
   ],
   "multiCdnSettings":{
      "origins":[
         {
            "originId":"TestOriginID",
            "hostname":"TestHostname",
            "propertyId":123
         }
      ],
      "cdns":[
         {
            "cdnCode":"TestCDNCode",
            "enabled":true,
            "cdnAuthKeys":[
               {
                  "authKeyName":"TestAuthKeyName",
                  "expiryDate":"TestExpiryDate",
                  "headerName":"TestHeaderName",
                  "secret":"TestSecretTestSecret1234"
               }
            ],
            "ipAclCidrs":[
               "1.1.1.1/1"
            ],
            "httpsOnly":true
         }
      ],
      "dataStreams":{
         "enabled":true,
         "dataStreamIds":[
            1
         ],
         "samplingRate":10
      },
      "bocc":{
         "enabled":true,
         "conditionalSamplingFrequency":"ZERO",
         "forwardType":"ORIGIN_AND_MIDGRESS",
         "requestType":"EDGE_AND_MIDGRESS",
         "samplingFrequency":"ZERO"
      },
      "enableSoftAlerts":true
   },
   "capacityAlertsThreshold":80,
   "notificationEmails":[
      "test@test.com"
   ],
   "lastUpdatedDate":"2022-06-10T13:21:14.488Z",
   "lastUpdatedBy":"johndoe",
   "lastActivatedDate":null,
   "lastActivatedBy":null
}`,
			expectedResponse: &Configuration{
				CapacityAlertsThreshold: ptr.To(80),
				Comments:                "TestCommentsUpdated",
				ContractID:              "TestContractID",
				ConfigID:                111,
				Locations: []ConfigLocationResp{
					{
						Comments:      "TestCommentsUpdated",
						TrafficTypeID: 1,
						Capacity: Capacity{
							Unit:  "GB",
							Value: 1,
						},
						MapName: "cw-s-use",
					},
				},
				MultiCDNSettings: &MultiCDNSettings{
					BOCC: &BOCC{
						ConditionalSamplingFrequency: SamplingFrequencyZero,
						Enabled:                      true,
						ForwardType:                  ForwardTypeOriginAndMidgress,
						RequestType:                  RequestTypeEdgeAndMidgress,
						SamplingFrequency:            SamplingFrequencyZero,
					},
					CDNs: []CDN{
						{
							CDNAuthKeys: []CDNAuthKey{
								{
									AuthKeyName: "TestAuthKeyName",
									ExpiryDate:  "TestExpiryDate",
									HeaderName:  "TestHeaderName",
									Secret:      "TestSecretTestSecret1234",
								},
							},
							CDNCode: "TestCDNCode",
							Enabled: true,
							IPACLCIDRs: []string{
								"1.1.1.1/1",
							},
							HTTPSOnly: true,
						},
					},
					DataStreams: &DataStreams{
						DataStreamIDs: []int64{1},
						Enabled:       true,
						SamplingRate:  ptr.To(10),
					},
					EnableSoftAlerts: true,
					Origins: []Origin{
						{
							Hostname:   "TestHostname",
							OriginID:   "TestOriginID",
							PropertyID: 123,
						},
					},
				},
				Status:          StatusInProgress,
				ConfigName:      "TestConfigName",
				LastUpdatedBy:   "johndoe",
				LastUpdatedDate: "2022-06-10T13:21:14.488Z",
				NotificationEmails: []string{
					"test@test.com",
				},
				PropertyIDs:       []string{"123"},
				RetainIdleObjects: true,
			},
		},
		"missing required params - validation error": {
			params: UpdateConfigurationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update configuration: struct validation: Body: {\n\tComments: cannot be blank\n\tLocations: cannot be blank\n\tPropertyIDs: cannot be blank\n}\nConfigID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: UpdateConfigurationRequest{
				ConfigID: 1,
				Body: UpdateConfigurationRequestBody{
					Comments: "TestCommentsUpdated",
					Locations: []ConfigLocationReq{
						{
							Comments:      "TestCommentsUpdated",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
					},
					PropertyIDs: []string{"1"},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/cloudwrapper/error-types/cloudwrapper-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/cloudwrapper/error-instances/abc",
    "status": 500
}`,
			expectedPath: "/cloud-wrapper/v1/configurations/1?activate=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/cloudwrapper/error-types/cloudwrapper-server-error",
					Title:    "An unexpected error has occurred.",
					Detail:   "Error processing request",
					Instance: "/cloudwrapper/error-instances/abc",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateConfiguration(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteConfiguration(t *testing.T) {
	tests := map[string]struct {
		params         DeleteConfigurationRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"202 - Accepted": {
			params: DeleteConfigurationRequest{
				ConfigID: 1,
			},
			responseStatus: 202,
			expectedPath:   "/cloud-wrapper/v1/configurations/1",
		},
		"missing required params - validation error": {
			params: DeleteConfigurationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete configuration: struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: DeleteConfigurationRequest{
				ConfigID: 1,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/cloudwrapper/error-types/cloudwrapper-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/cloudwrapper/error-instances/abc",
    "status": 500
}`,
			expectedPath: "/cloud-wrapper/v1/configurations/1",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/cloudwrapper/error-types/cloudwrapper-server-error",
					Title:    "An unexpected error has occurred.",
					Detail:   "Error processing request",
					Instance: "/cloudwrapper/error-instances/abc",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
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
			err := client.DeleteConfiguration(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestActivateConfiguration(t *testing.T) {
	tests := map[string]struct {
		params              ActivateConfigurationRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		withError           func(*testing.T, error)
	}{
		"204 - single configID": {
			params: ActivateConfigurationRequest{
				ConfigurationIDs: []int{1},
			},
			expectedRequestBody: `
{
   "configurationIds":[
      1
   ]
}`,
			responseStatus: 204,
			expectedPath:   "/cloud-wrapper/v1/configurations/activate",
		},
		"204 - multiple configIDs": {
			params: ActivateConfigurationRequest{
				ConfigurationIDs: []int{1, 2, 3},
			},
			expectedRequestBody: `
{
	"configurationIds": [
		1,
		2,
		3
	]
}`,
			responseStatus: 204,
			expectedPath:   "/cloud-wrapper/v1/configurations/activate",
		},
		"missing required params - validation error": {
			params: ActivateConfigurationRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "activate configuration: struct validation: ConfigurationIDs: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: ActivateConfigurationRequest{
				ConfigurationIDs: []int{1},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "/cloudwrapper/error-types/cloudwrapper-server-error",
    "title": "An unexpected error has occurred.",
    "detail": "Error processing request",
    "instance": "/cloudwrapper/error-instances/abc",
    "status": 500
}`,
			expectedPath: "/cloud-wrapper/v1/configurations/activate",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/cloudwrapper/error-types/cloudwrapper-server-error",
					Title:    "An unexpected error has occurred.",
					Detail:   "Error processing request",
					Instance: "/cloudwrapper/error-instances/abc",
					Status:   500,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.ActivateConfiguration(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
