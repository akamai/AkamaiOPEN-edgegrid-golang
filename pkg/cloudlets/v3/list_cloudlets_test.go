package v3

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCloudlets(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []ListCloudletsItem
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
[
  {
    "cloudletName": "API_PRIORITIZATION",
    "cloudletType": "AP"
  },
  {
    "cloudletName": "AUDIENCE_SEGMENTATION",
    "cloudletType": "AS"
  },
  {
    "cloudletName": "EDGE_REDIRECTOR",
    "cloudletType": "ER"
  },
  {
    "cloudletName": "FORWARD_REWRITE",
    "cloudletType": "FR"
  },
  {
    "cloudletName": "PHASED_RELEASE",
    "cloudletType": "CD"
  },
  {
    "cloudletName": "REQUEST_CONTROL",
    "cloudletType": "IG"
  }
]`,
			expectedPath: "/cloudlets/v3/cloudlet-info",
			expectedResponse: []ListCloudletsItem{
				{
					CloudletName: "API_PRIORITIZATION",
					CloudletType: "AP",
				},
				{
					CloudletName: "AUDIENCE_SEGMENTATION",
					CloudletType: "AS",
				},
				{
					CloudletName: "EDGE_REDIRECTOR",
					CloudletType: "ER",
				},
				{
					CloudletName: "FORWARD_REWRITE",
					CloudletType: "FR",
				},
				{
					CloudletName: "PHASED_RELEASE",
					CloudletType: "CD",
				},
				{
					CloudletName: "REQUEST_CONTROL",
					CloudletType: "IG",
				},
			},
		},
		"500 Internal Server Error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"status": 500,
	"requestId": "1",
	"requestTime": "12:00",
	"clientIp": "1.1.1.1",
	"serverIp": "2.2.2.2",
	"method": "GET"
}`,
			expectedPath: "/cloudlets/v3/cloudlet-info",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "internal_error",
					Title:       "Internal Server Error",
					Status:      http.StatusInternalServerError,
					RequestID:   "1",
					RequestTime: "12:00",
					ClientIP:    "1.1.1.1",
					ServerIP:    "2.2.2.2",
					Method:      "GET",
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
			result, err := client.ListCloudlets(context.Background())
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
