package cloudlets

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListOrigins(t *testing.T) {
	tests := map[string]struct {
		originType       ListOriginsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []OriginResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			originType:     ListOriginsRequest{},
			responseStatus: http.StatusOK,
			responseBody: `[
				{
					
					"hostname": "",
					"description": "ALB1",
					"originId": "alb1",
					"type": "APPLICATION_LOAD_BALANCER",
					"akamaized": false
				},
				{
					"hostname": "",
					"description": "",
					"originId": "alb2",
					"type": "APPLICATION_LOAD_BALANCER",
					"akamaized": false
				},
				{
					"hostname": "dc1.foo.com",
					"description": "",
					"originId": "dc1",
					"type": "CUSTOMER",
					"akamaized": false
				},
				{
					
					"hostname": "dc2.foo.com",
					"description": "",
					"originId": "dc2",
					"type": "CUSTOMER",
					"akamaized": true
				},
				{
					"hostname": "download.akamai.com/12345",
					"description": "",
					"originId": "ns1",
					"type": "NETSTORAGE",
					"akamaized": true
				},
				{
					
					"hostname": "download.akamai.com/12345",
					"description": "",
					"originId": "ns2",
					"type": "NETSTORAGE",
					"akamaized": true
				}
			]`,
			expectedPath: "/cloudlets/api/v2/origins",
			expectedResponse: []OriginResponse{
				{
					Hostname: "",
					Origin: Origin{
						OriginID:    "alb1",
						Description: "ALB1",
						Type:        "APPLICATION_LOAD_BALANCER",
						Akamaized:   false,
					},
				},
				{
					Hostname: "",
					Origin: Origin{
						Description: "",
						OriginID:    "alb2",
						Type:        "APPLICATION_LOAD_BALANCER",
						Akamaized:   false,
					},
				},
				{
					Hostname: "dc1.foo.com",
					Origin: Origin{
						Description: "",
						OriginID:    "dc1",
						Type:        "CUSTOMER",
						Akamaized:   false,
					},
				},
				{
					Hostname: "dc2.foo.com",
					Origin: Origin{
						Description: "",
						OriginID:    "dc2",
						Type:        "CUSTOMER",
						Akamaized:   true,
					},
				},
				{
					Hostname: "download.akamai.com/12345",
					Origin: Origin{
						Description: "",
						OriginID:    "ns1",
						Type:        "NETSTORAGE",
						Akamaized:   true,
					},
				},
				{
					Hostname: "download.akamai.com/12345",
					Origin: Origin{
						Description: "",
						OriginID:    "ns2",
						Type:        "NETSTORAGE",
						Akamaized:   true,
					},
				},
			},
		},
		"200 ok with param": {
			originType:     ListOriginsRequest{Type: OriginTypeCustomer},
			responseStatus: http.StatusOK,
			responseBody: `[
				{
					"hostname": "dc1.foo.com",
					"description": "",
					"originId": "dc1",
					"type": "CUSTOMER",
					"akamaized": false
				},
				{
					
					"hostname": "dc2.foo.com",
					"description": "",
					"originId": "dc2",
					"type": "CUSTOMER",
					"akamaized": true
				}
			]`,
			expectedPath: "/cloudlets/api/v2/origins?type=CUSTOMER",
			expectedResponse: []OriginResponse{
				{
					Hostname: "dc1.foo.com",
					Origin: Origin{
						Description: "",
						OriginID:    "dc1",
						Type:        "CUSTOMER",
						Akamaized:   false,
					},
				},
				{
					Hostname: "dc2.foo.com",
					Origin: Origin{
						Description: "",
						OriginID:    "dc2",
						Type:        "CUSTOMER",
						Akamaized:   true,
					},
				},
			},
		},
		"500 internal server error": {
			originType:     ListOriginsRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "internal_error",
			   	"title": "Internal Server Error",
			   	"detail": "Error making request",
			   	"status": 500
			}`,
			expectedPath: "/cloudlets/api/v2/origins",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
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
			result, err := client.ListOrigins(context.Background(), test.originType)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetOrigin(t *testing.T) {
	tests := map[string]struct {
		originID         string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Origin
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			originID:       "alb1",
			responseStatus: http.StatusOK,
			responseBody: `{
				"description": "ALB1",
				"originId": "alb1",
				"type": "APPLICATION_LOAD_BALANCER",
				"checksum": "abcdefg1111hijklmn22222fff76yae3"
			}`,
			expectedPath: "/cloudlets/api/v2/origins/alb1",
			expectedResponse: &Origin{
				Description: "ALB1",
				OriginID:    "alb1",
				Type:        OriginTypeApplicationLoadBalancer,
				Checksum:    "abcdefg1111hijklmn22222fff76yae3",
			},
		},
		"500 internal server error": {
			originID:       "ALB1",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "internal_error",
			   	"title": "Internal Server Error",
			   	"detail": "Error making request",
			   	"status": 500
			}`,
			expectedPath: "/cloudlets/api/v2/origins/ALB1",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
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
			result, err := client.GetOrigin(context.Background(), test.originID)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreateOrigin(t *testing.T) {
	tests := map[string]struct {
		request          LoadBalancerOriginRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Origin
		withError        error
	}{
		"201 created": {
			request: LoadBalancerOriginRequest{
				OriginID:    "first",
				Description: "create first Origin",
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
			   "originId": "first",
			   "akamaized": true,
			   "checksum": "9c0fc1f3e9ea7eb2e090f2bf53709e45",
			   "description": "create first Origin",
			   "type": "APPLICATION_LOAD_BALANCER"
			}`,
			expectedPath: "/cloudlets/api/v2/origins",
			expectedResponse: &Origin{
				OriginID:    "first",
				Description: "create first Origin",
				Akamaized:   true,
				Type:        OriginTypeApplicationLoadBalancer,
				Checksum:    "9c0fc1f3e9ea7eb2e090f2bf53709e45",
			},
		},
		"500 internal server error": {
			request: LoadBalancerOriginRequest{
				OriginID:    "second",
				Description: "create second Origin",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
				  "type": "internal_error",
				  "title": "Internal Server Error",
				  "detail": "Error creating enrollment",
				  "status": 500 
				}`,
			expectedPath: "/cloudlets/api/v2/origins",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
			expectedResponse: &Origin{
				OriginID:    "first",
				Description: "create first Origin",
				Akamaized:   false,
				Type:        OriginTypeApplicationLoadBalancer,
				Checksum:    "9c0fc1f3e9ea7eb2e090f2bf53709e45",
			},
		},
		"validation error": {
			request:   LoadBalancerOriginRequest{},
			withError: ErrStructValidation,
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
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateOrigin(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateOrigin(t *testing.T) {
	tests := map[string]struct {
		request             LoadBalancerOriginRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *Origin
		withError           error
	}{
		"200 updated": {
			request: LoadBalancerOriginRequest{
				OriginID:    "first",
				Description: "update first Origin",
			},
			expectedRequestBody: `{"description":"update first Origin"}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
			   "originId": "first",
			   "akamaized": true,
			   "checksum": "9c0fc1f3e9ea7eb2e090f2bf53709e45",
			   "description": "update first Origin",
			   "type": "APPLICATION_LOAD_BALANCER"
			}`,
			expectedPath: "/cloudlets/api/v2/origins/first",
			expectedResponse: &Origin{
				OriginID:    "first",
				Description: "update first Origin",
				Akamaized:   true,
				Type:        OriginTypeApplicationLoadBalancer,
				Checksum:    "9c0fc1f3e9ea7eb2e090f2bf53709e45",
			},
		},
		"500 internal server error": {
			request: LoadBalancerOriginRequest{
				OriginID:    "second",
				Description: "create second Origin",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
				  "type": "internal_error",
				  "title": "Internal Server Error",
				  "detail": "Error creating enrollment",
				  "status": 500 
				}`,
			expectedPath: "/cloudlets/api/v2/origins/second",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
			expectedResponse: &Origin{
				OriginID:    "second",
				Description: "update first Origin",
				Akamaized:   false,
				Type:        OriginTypeApplicationLoadBalancer,
				Checksum:    "9c0fc1f3e9ea7eb2e090f2bf53709e45",
			},
		},
		"validation error": {
			request:   LoadBalancerOriginRequest{},
			withError: ErrStructValidation,
		},
		"validation error -  OriginID exceeds max length, which is 63": {
			request: LoadBalancerOriginRequest{
				OriginID: "ExceedMaxLenghtExceedMaxLenghtExceedMaxLenghtExceedMaxLenghtExce",
			},
			withError: ErrStructValidation,
		},
		"validation error - OriginID value less than min, which is 2": {
			request: LoadBalancerOriginRequest{
				OriginID: "E",
			},
			withError: ErrStructValidation,
		},
		"validation error - Description exceeds max length, which is 255": {
			request: LoadBalancerOriginRequest{
				OriginID:    "first",
				Description: "Test for creating APPLICATION_LOAD_BALANCER origin type, Test for creating APPLICATION_LOAD_BALANCER origin type, Test for creating APPLICATION_LOAD_BALANCER origin type, Test for creating APPLICATION_LOAD_BALANCER origin type,Test for creating exceed valu",
			},
			withError: ErrStructValidation,
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateOrigin(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
