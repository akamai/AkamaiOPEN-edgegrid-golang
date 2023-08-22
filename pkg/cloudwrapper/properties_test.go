package cloudwrapper

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProperties(t *testing.T) {
	tests := map[string]struct {
		params           ListPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListPropertiesResponse
		withError        error
	}{
		"200 OK - multiple properties": {
			responseStatus: 200,
			responseBody: `
{
   "properties":[
      {
         "propertyId":1,
         "propertyName":"TestPropertyName1",
         "contractId":"TestContractID1",
         "groupId":11,
         "type":"MEDIA"
      },
      {
         "propertyId":2,
         "propertyName":"TestPropertyName2",
         "contractId":"TestContractID2",
         "groupId":22,
         "type":"WEB"
      },
      {
         "propertyId":3,
         "propertyName":"TestPropertyName3",
         "contractId":"TestContractID3",
         "groupId":33,
         "type":"WEB"
      }
   ]
}`,
			expectedPath: "/cloud-wrapper/v1/properties?unused=false",
			expectedResponse: &ListPropertiesResponse{
				Properties: []Property{
					{
						GroupID:      11,
						ContractID:   "TestContractID1",
						PropertyID:   1,
						PropertyName: "TestPropertyName1",
						Type:         PropertyTypeMedia,
					},
					{
						GroupID:      22,
						ContractID:   "TestContractID2",
						PropertyID:   2,
						PropertyName: "TestPropertyName2",
						Type:         PropertyTypeWeb,
					},
					{
						GroupID:      33,
						ContractID:   "TestContractID3",
						PropertyID:   3,
						PropertyName: "TestPropertyName3",
						Type:         PropertyTypeWeb,
					},
				},
			},
		},
		"200 OK - single property": {
			responseStatus: 200,
			responseBody: `
{
   "properties":[
      {
         "propertyId":1,
         "propertyName":"TestPropertyName1",
         "contractId":"TestContractID1",
         "groupId":11,
         "type":"MEDIA"
      }
   ]
}`,
			expectedPath: "/cloud-wrapper/v1/properties?unused=false",
			expectedResponse: &ListPropertiesResponse{
				Properties: []Property{
					{
						GroupID:      11,
						ContractID:   "TestContractID1",
						PropertyID:   1,
						PropertyName: "TestPropertyName1",
						Type:         PropertyTypeMedia,
					},
				},
			},
		},
		"200 OK - properties with query params": {
			params: ListPropertiesRequest{
				Unused: true,
				ContractIDs: []string{
					"TestContractID1",
					"TestContractID2",
				},
			},
			responseStatus: 200,
			responseBody: `
{
   "properties":[
      {
         "propertyId":1,
         "propertyName":"TestPropertyName1",
         "contractId":"TestContractID1",
         "groupId":11,
         "type":"MEDIA"
      },
      {
         "propertyId":2,
         "propertyName":"TestPropertyName2",
         "contractId":"TestContractID2",
         "groupId":22,
         "type":"MEDIA"
      }
   ]
}`,
			expectedPath: "/cloud-wrapper/v1/properties?contractIds=TestContractID1&contractIds=TestContractID2&unused=true",
			expectedResponse: &ListPropertiesResponse{
				Properties: []Property{
					{
						GroupID:      11,
						ContractID:   "TestContractID1",
						PropertyID:   1,
						PropertyName: "TestPropertyName1",
						Type:         PropertyTypeMedia,
					},
					{
						GroupID:      22,
						ContractID:   "TestContractID2",
						PropertyID:   2,
						PropertyName: "TestPropertyName2",
						Type:         PropertyTypeMedia,
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
			expectedPath: "/cloud-wrapper/v1/properties?unused=false",
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
			result, err := client.ListProperties(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListOrigins(t *testing.T) {
	tests := map[string]struct {
		params           ListOriginsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListOriginsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - multiple objects": {
			params: ListOriginsRequest{
				PropertyID: 1,
				ContractID: "TestContractID",
				GroupID:    11,
			},
			responseStatus: 200,
			responseBody: `
{
   "default":[
      {
         "originType":"CUSTOMER",
         "hostname":"origin-www.example.com"
      },
      {
         "originType":"NET_STORAGE",
         "hostname":"origin-www.example2.com"
      }
   ],
   "children":[
      {
         "name":"Default CORS Policy",
         "behaviors":[
            {
               "originType":"NET_STORAGE",
               "hostname":"origin-www.example3.com"
            }
         ]
      },
      {
         "name":"Cloud Wrapper",
         "behaviors":[
            {
               "originType":"CUSTOMER",
               "hostname":"origin-www.example4.com"
            },
            {
               "originType":"CUSTOMER",
               "hostname":"origin-www.example5.com"
            }
         ]
      }
   ]
}`,
			expectedPath: "/cloud-wrapper/v1/properties/1/origins?contractId=TestContractID&groupId=11",
			expectedResponse: &ListOriginsResponse{
				Children: []Child{
					{
						Name: "Default CORS Policy",
						Behaviors: []Behavior{
							{
								Hostname:   "origin-www.example3.com",
								OriginType: OriginTypeNetStorage,
							},
						},
					},
					{
						Name: "Cloud Wrapper",
						Behaviors: []Behavior{
							{
								Hostname:   "origin-www.example4.com",
								OriginType: OriginTypeCustomer,
							},
							{
								Hostname:   "origin-www.example5.com",
								OriginType: OriginTypeCustomer,
							},
						},
					},
				},
				Default: []Behavior{
					{
						Hostname:   "origin-www.example.com",
						OriginType: OriginTypeCustomer,
					},
					{
						Hostname:   "origin-www.example2.com",
						OriginType: OriginTypeNetStorage,
					},
				},
			},
		},
		"200 OK - empty behaviors": {
			params: ListOriginsRequest{
				PropertyID: 1,
				ContractID: "TestContractID",
				GroupID:    11,
			},
			responseStatus: 200,
			responseBody: `
{
   "default":[
      {
         "originType":"CUSTOMER",
         "hostname":"test.com"
      }
   ],
   "children":[
      {
         "name":"Default CORS Policy",
         "behaviors":[
            
         ]
      }
   ]
}`,
			expectedPath: "/cloud-wrapper/v1/properties/1/origins?contractId=TestContractID&groupId=11",
			expectedResponse: &ListOriginsResponse{
				Children: []Child{
					{
						Name:      "Default CORS Policy",
						Behaviors: []Behavior{},
					},
				},
				Default: []Behavior{
					{
						Hostname:   "test.com",
						OriginType: OriginTypeCustomer,
					},
				},
			},
		},
		"missing required params - validation errors": {
			params: ListOriginsRequest{
				PropertyID: 0,
				ContractID: "",
				GroupID:    0,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list origins: struct validation: ContractID: cannot be blank\nGroupID: cannot be blank\nPropertyID: cannot be blank", err.Error())
			},
		},
		"500 internal server error": {
			params: ListOriginsRequest{
				PropertyID: 1,
				ContractID: "TestContractID",
				GroupID:    11,
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
			expectedPath: "/cloud-wrapper/v1/properties/1/origins?contractId=TestContractID&groupId=11",
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
			result, err := client.ListOrigins(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
