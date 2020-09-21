package papi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPapi_SearchProperties(t *testing.T) {
	tests := map[string]struct {
		params           SearchRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedRequest  string
		expectedResponse *SearchResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: SearchRequest{
				key:   "edgeHostname",
				value: "edgesuite.net",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "versions": {
        "items": [
            {
                "accountId": "accountID_1",
                "assetId": "assetID_1",
                "contractId": "contractID_1",
                "groupId": "groupID_1",
                "productionStatus": "INACTIVE",
                "propertyId": "propertyID_1",
                "propertyName": "propertyName_1",
                "propertyVersion": 1,
                "stagingStatus": "INACTIVE",
                "updatedByUser": "user_1",
                "updatedDate": "2017-08-07T15:39:49Z"
            },
            {
                "accountId": "accountID_2",
                "assetId": "assetID_2",
                "contractId": "contractID_2",
                "groupId": "groupID_2",
                "productionStatus": "INACTIVE",
                "propertyId": "propertyID_2",
                "propertyName": "propertyName_2",
                "propertyVersion": 2,
                "stagingStatus": "INACTIVE",
                "updatedByUser": "user_2",
                "updatedDate": "2017-08-07T15:39:49Z"
            }
        ]
    }
}`,
			expectedRequest: `
{
	"edgeHostname": "edgesuite.net"
}`,
			expectedPath: "/papi/v1/search/find-by-value",
			expectedResponse: &SearchResponse{
				Versions: SearchItems{
					Items: []SearchItem{
						{
							AccountID:        "accountID_1",
							AssetID:          "assetID_1",
							ContractID:       "contractID_1",
							GroupID:          "groupID_1",
							ProductionStatus: "INACTIVE",
							PropertyID:       "propertyID_1",
							PropertyName:     "propertyName_1",
							PropertyVersion:  1,
							StagingStatus:    "INACTIVE",
							UpdatedByUser:    "user_1",
							UpdatedDate:      "2017-08-07T15:39:49Z",
						},
						{
							AccountID:        "accountID_2",
							AssetID:          "assetID_2",
							ContractID:       "contractID_2",
							GroupID:          "groupID_2",
							ProductionStatus: "INACTIVE",
							PropertyID:       "propertyID_2",
							PropertyName:     "propertyName_2",
							PropertyVersion:  2,
							StagingStatus:    "INACTIVE",
							UpdatedByUser:    "user_2",
							UpdatedDate:      "2017-08-07T15:39:49Z",
						},
					},
				},
			},
		},
		"500 Internal Server Error": {
			params: SearchRequest{
				key:   "edgeHostname",
				value: "edgesuite.net",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error searching for property",
    "status": 505
}`,
			expectedRequest: `
{
	"edgeHostname": "edgesuite.net"
}`,
			expectedPath: "/papi/v1/search/find-by-value",
			withError: func(t *testing.T, err error) {
				want := session.APIError{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error searching for property",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"invalid key": {
			params: SearchRequest{
				key:   "test",
				value: "edgesuite.net",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SearchKey")
			},
		},
		"empty key": {
			params: SearchRequest{
				key:   "",
				value: "edgesuite.net",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SearchKey")
			},
		},
		"empty value": {
			params: SearchRequest{
				key:   "edgeHostname",
				value: "",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SearchValue")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				body, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err)
				var compact bytes.Buffer
				err = json.Compact(&compact, []byte(test.expectedRequest))
				require.NoError(t, err)
				assert.Equal(t, compact.String(), string(body))
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.SearchProperties(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
