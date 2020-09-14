package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi/tools"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapi_GetProperties(t *testing.T) {
	tests := map[string]struct {
		request          GetPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPropertiesResponse
		withError        error
	}{
		"200 OK": {
			request: GetPropertiesRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"properties": {
		"items": [
			{
				"accountId": "act_1-1TJZFB",
				"contractId": "ctr_1-1TJZH5",
				"groupId": "grp_15166",
				"propertyId": "prp_175780",
				"propertyName": "example.com",
				"latestVersion": 2,
				"stagingVersion": 1,
				"productId": "prp_175780",
				"productionVersion": null,
				"assetId": "aid_101",
				"note": "Notes about example.com"
			}
		]
	}
}`,
			expectedPath: "/papi/v1/properties?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &GetPropertiesResponse{
				Properties: PropertiesItems{Items: []*Property{
					{
						AccountID:         "act_1-1TJZFB",
						ContactID:         "ctr_1-1TJZH5",
						GroupID:           "grp_15166",
						PropertyID:        "prp_175780",
						ProductID:         "prp_175780",
						PropertyName:      "example.com",
						LatestVersion:     2,
						StagingVersion:    tools.IntPtr(1),
						ProductionVersion: nil,
						AssetID:           "aid_101",
						Note:              "Notes about example.com",
					},
				}},
			},
		},
		"500 internal server error": {
			request: GetPropertiesRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching properties",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: session.APIError{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching properties",
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
			result, err := client.GetProperties(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
