package appsec

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test IPGeo
func TestAppSec_GetIPGeo(t *testing.T) {
	tests := map[string]struct {
		params           GetIPGeoRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetIPGeoResponse
		withError        error
	}{
		"200 OK": {
			params: GetIPGeoRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"block": "blockSpecificIPGeo",
				"asnControls": {
					"blockedIPNetworkLists": {
						"networkList": [
							"12345_ASNTEST"
						]
					}
				},
				"geoControls": {
					"blockedIPNetworkLists": {
						"networkList": [
							"72138_TEST1"
						]
					}
				},
				"ipControls": {
					"allowedIPNetworkLists": {
						"networkList": [
							"56921_TEST"
						]
					},
					"blockedIPNetworkLists": {
						"networkList": [
							"53712_TESTLIST123"
						]
					}
				},
				"ukraineGeoControl": {
					"action": "alert"
				}
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/ip-geo-firewall",
			expectedResponse: &GetIPGeoResponse{
				Block: "blockSpecificIPGeo",
				GeoControls: &IPGeoGeoControls{
					BlockedIPNetworkLists: &IPGeoNetworkLists{
						NetworkList: []string{"72138_TEST1"},
					},
				},
				IPControls: &IPGeoIPControls{
					AllowedIPNetworkLists: &IPGeoNetworkLists{
						NetworkList: []string{"56921_TEST"},
					},
					BlockedIPNetworkLists: &IPGeoNetworkLists{
						NetworkList: []string{"53712_TESTLIST123"},
					},
				},
				ASNControls: &IPGeoASNControls{
					BlockedIPNetworkLists: &IPGeoNetworkLists{
						NetworkList: []string{"12345_ASNTEST"},
					},
				},
				UkraineGeoControls: &UkraineGeoControl{
					Action: "alert",
				},
			},
		},
		"500 internal server error": {
			params: GetIPGeoRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/ip-geo-firewall",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching match target",
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
			result, err := client.GetIPGeo(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update IPGeo.
func TestAppSec_UpdateIPGeo(t *testing.T) {
	tests := map[string]struct {
		params           UpdateIPGeoRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateIPGeoResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateIPGeoRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
				"block": "blockSpecificIPGeo",
				"asnControls": {
					"blockedIPNetworkLists": {
						"networkList": [
							"12345_ASNTEST"
						]
					}
				},
				"geoControls": {
					"blockedIPNetworkLists": {
						"networkList": [
							"72138_TEST1"
						]
					}
				},
				"ipControls": {
					"allowedIPNetworkLists": {
						"networkList": [
							"56921_TEST"
						]
					},
					"blockedIPNetworkLists": {
						"networkList": [
							"53712_TESTLIST123"
						]
					}
				},
				"ukraineGeoControl": {
					"action": "alert"
				}
			}`,
			expectedResponse: &UpdateIPGeoResponse{
				Block: "blockSpecificIPGeo",
				GeoControls: &IPGeoGeoControls{
					BlockedIPNetworkLists: &IPGeoNetworkLists{
						NetworkList: []string{"72138_TEST1"},
					},
				},
				IPControls: &IPGeoIPControls{
					AllowedIPNetworkLists: &IPGeoNetworkLists{
						NetworkList: []string{"56921_TEST"},
					},
					BlockedIPNetworkLists: &IPGeoNetworkLists{
						NetworkList: []string{"53712_TESTLIST123"},
					},
				},
				ASNControls: &IPGeoASNControls{
					BlockedIPNetworkLists: &IPGeoNetworkLists{
						NetworkList: []string{"12345_ASNTEST"},
					},
				},
				UkraineGeoControls: &UkraineGeoControl{
					Action: "alert",
				},
			},
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/ip-geo-firewall",
		},
		"500 internal server error": {
			params: UpdateIPGeoRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/ip-geo-firewall",
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
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateIPGeo(
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
