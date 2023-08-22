package cloudwrapper

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tj/assert"
)

func TestCloudwrapper_ListLocations(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		expectedPath     string
		responseBody     string
		expectedResponse *ListLocationResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			expectedPath:   "/cloud-wrapper/v1/locations",
			responseBody: `{
    "locations": [
        {
            "locationId": 1,
            "locationName": "US East",
            "trafficTypes": [
                {
                    "trafficTypeId": 1,
                    "trafficType": "TEST_TT1",
                    "mapName": "cw-essl-use"
                },
                {
                    "trafficTypeId": 2,
                    "trafficType": "TEST_TT2",
                    "mapName": "cw-s-use-live"
                }
            ],
            "multiCdnLocationId": "0123"
        },
        {
            "locationId": 2,
            "locationName": "US West",
            "trafficTypes": [
                {
                    "trafficTypeId": 3,
                    "trafficType": "TEST_TT1",
                    "mapName": "cw-essl-use"
                },
                {
                    "trafficTypeId": 4,
                    "trafficType": "TEST_TT2",
                    "mapName": "cw-s-use-live"
                }
            ],
            "multiCdnLocationId": "4567"
        }
	]}`,
			expectedResponse: &ListLocationResponse{
				Locations: []Location{
					{
						LocationID:   1,
						LocationName: "US East",
						TrafficTypes: []TrafficTypeItem{
							{
								TrafficTypeID: 1,
								TrafficType:   "TEST_TT1",
								MapName:       "cw-essl-use",
							},
							{
								TrafficTypeID: 2,
								TrafficType:   "TEST_TT2",
								MapName:       "cw-s-use-live",
							},
						},
						MultiCDNLocationID: "0123",
					},
					{
						LocationID:   2,
						LocationName: "US West",
						TrafficTypes: []TrafficTypeItem{
							{
								TrafficTypeID: 3,
								TrafficType:   "TEST_TT1",
								MapName:       "cw-essl-use",
							},
							{
								TrafficTypeID: 4,
								TrafficType:   "TEST_TT2",
								MapName:       "cw-s-use-live",
							},
						},
						MultiCDNLocationID: "4567",
					},
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/cloud-wrapper/v1/locations",
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error processing request",
				"status": 500
			}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Detail: "Error processing request",
					Status: http.StatusInternalServerError,
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
			users, err := client.ListLocations(context.Background())
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, users)
		})
	}
}
