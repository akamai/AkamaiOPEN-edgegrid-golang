package cloudlets

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListOrigins(t *testing.T) {
	tests := map[string]struct {
		originType       ListOriginsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse Origins
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			originType:     ListOriginsRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
[
	{
		"description": "ALB1",
		"hostname": "",
		"originId": "alb1",
		"type": "APPLICATION_LOAD_BALANCER",
		"akamaized": false
	},
	{
		"description": "",
		"hostname": "",
		"originId": "alb2",
		"type": "APPLICATION_LOAD_BALANCER",
		"akamaized": false
	},
	{
		"description": "",
		"hostname": "dc1.foo.com",
		"originId": "dc1",
		"type": "CUSTOMER",
		"akamaized": false
	},
	{
		"description": "",
		"hostname": "dc2.foo.com",
		"originId": "dc2",
		"type": "CUSTOMER",
		"akamaized": true
	},
	{
		"description": "",
		"hostname": "download.akamai.com/12345",
		"originId": "ns1",
		"type": "NETSTORAGE",
		"akamaized": true
	},
	{
		"description": "",
		"hostname": "download.akamai.com/12345",
		"originId": "ns2",
		"type": "NETSTORAGE",
		"akamaized": true
	}
]`,
			expectedPath: "/cloudlets/api/v2/origins",
			expectedResponse: Origins{
				{
					Description: "ALB1",
					Hostname:    "",
					OriginID:    "alb1",
					Type:        "APPLICATION_LOAD_BALANCER",
					Akamaized:   false,
				},
				{
					Description: "",
					Hostname:    "",
					OriginID:    "alb2",
					Type:        "APPLICATION_LOAD_BALANCER",
					Akamaized:   false,
				},
				{
					Description: "",
					Hostname:    "dc1.foo.com",
					OriginID:    "dc1",
					Type:        "CUSTOMER",
					Akamaized:   false,
				},
				{
					Description: "",
					Hostname:    "dc2.foo.com",
					OriginID:    "dc2",
					Type:        "CUSTOMER",
					Akamaized:   true,
				},
				{
					Description: "",
					Hostname:    "download.akamai.com/12345",
					OriginID:    "ns1",
					Type:        "NETSTORAGE",
					Akamaized:   true,
				},
				{
					Description: "",
					Hostname:    "download.akamai.com/12345",
					OriginID:    "ns2",
					Type:        "NETSTORAGE",
					Akamaized:   true,
				},
			},
		},
		"200 ok with param": {
			originType:     ListOriginsRequest{Type: OriginTypeCustomer},
			responseStatus: http.StatusOK,
			responseBody: `
[
	{
		"description": "",
		"hostname": "dc1.foo.com",
		"originId": "dc1",
		"type": "CUSTOMER",
		"akamaized": false
	},
	{
		"description": "",
		"hostname": "dc2.foo.com",
		"originId": "dc2",
		"type": "CUSTOMER",
		"akamaized": true
	}
]`,
			expectedPath: "/cloudlets/api/v2/origins?type=CUSTOMER",
			expectedResponse: Origins{
				{
					Description: "",
					Hostname:    "dc1.foo.com",
					OriginID:    "dc1",
					Type:        "CUSTOMER",
					Akamaized:   false,
				},
				{
					Description: "",
					Hostname:    "dc2.foo.com",
					OriginID:    "dc2",
					Type:        "CUSTOMER",
					Akamaized:   true,
				},
			},
		},
		"500 internal server error": {
			originType:     ListOriginsRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
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
			responseBody: `
{
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
			responseBody: `
{
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
