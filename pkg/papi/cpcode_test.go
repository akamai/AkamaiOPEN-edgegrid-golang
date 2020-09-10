package papi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPapi_GetCPCodes(t *testing.T) {
	tests := map[string]struct {
		contractID       string
		groupID          string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCPCodesResponse
		withError        error
	}{
		"200 OK": {
			contractID:     "contract",
			groupID:        "group",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "acc",
    "contractId": "contract",
    "groupId": "group",
    "cpcodes": {
        "items": [
            {
                "cpcodeId": "cpcode_id",
                "cpcodeName": "cpcode_name",
                "createdDate": "2020-09-10T15:06:13Z",
                "productIds": [
                    "prd_Web_App_Accel"
                ]
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/cpcodes?contractId=contract&groupId=group",
			expectedResponse: &GetCPCodesResponse{
				AccountID:  "acc",
				ContractID: "contract",
				GroupId:    "group",
				CPCodes: CPCodeItems{Items: []CPCode{
					{
						ID:          "cpcode_id",
						Name:        "cpcode_name",
						CreatedDate: "2020-09-10T15:06:13Z",
						ProductIds:  []string{"prd_Web_App_Accel"},
					},
				}},
			},
		},
		"404 Not found": {
			contractID:     "contract",
			groupID:        "group",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
	"type": "not_found",
    "title": "Not Found",
    "detail": "Could not find cp codes",
    "status": 404
}`,
			expectedPath: "/papi/v1/cpcodes?contractId=contract&groupId=group",
			withError:    session.ErrNotFound,
		},
		"500 internal server error": {
			contractID:     "contract",
			groupID:        "group",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching cp codes",
    "status": 500
}`,
			expectedPath: "/papi/v1/cpcodes?contractId=contract&groupId=group",
			withError: session.APIError{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching cp codes",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"empty group ID": {
			contractID: "contract",
			groupID:    "",
			withError:  ErrGroupEmpty,
		},
		"empty contract ID": {
			contractID: "",
			groupID:    "group",
			withError:  ErrContractEmpty,
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
			result, err := client.GetCPCodes(context.Background(), test.contractID, test.groupID)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_GetCPCode(t *testing.T) {
	tests := map[string]struct {
		contractID       string
		groupID          string
		id               string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCPCodesResponse
		withError        error
	}{
		"200 OK": {
			contractID:     "contract",
			groupID:        "group",
			id:             "cpcodeID",
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "acc",
    "contractId": "contract",
    "groupId": "group",
    "cpcodes": {
        "items": [
            {
                "cpcodeId": "cpcodeID",
                "cpcodeName": "cpcode_name",
                "createdDate": "2020-09-10T15:06:13Z",
                "productIds": [
                    "prd_Web_App_Accel"
                ]
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/cpcodes/cpcodeID?contractId=contract&groupId=group",
			expectedResponse: &GetCPCodesResponse{
				AccountID:  "acc",
				ContractID: "contract",
				GroupId:    "group",
				CPCodes: CPCodeItems{Items: []CPCode{
					{
						ID:          "cpcodeID",
						Name:        "cpcode_name",
						CreatedDate: "2020-09-10T15:06:13Z",
						ProductIds:  []string{"prd_Web_App_Accel"},
					},
				}},
			},
		},
		"404 Not found": {
			contractID:     "contract",
			groupID:        "group",
			id:             "not_found",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
	"type": "not_found",
    "title": "Not Found",
    "detail": "Could not find cp code",
    "status": 404
}`,
			expectedPath: "/papi/v1/cpcodes/not_found?contractId=contract&groupId=group",
			withError:    session.ErrNotFound,
		},
		"500 internal server error": {
			contractID:     "contract",
			groupID:        "group",
			id:             "cpcodeID",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching cp codes",
    "status": 500
}`,
			expectedPath: "/papi/v1/cpcodes/cpcodeID?contractId=contract&groupId=group",
			withError: session.APIError{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching cp codes",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"empty cpcode ID": {
			contractID: "contract",
			groupID:    "group",
			id:         "",
			withError:  ErrIDEmpty,
		},
		"empty group ID": {
			contractID: "contract",
			groupID:    "",
			id:         "id",
			withError:  ErrGroupEmpty,
		},
		"empty contract ID": {
			contractID: "",
			groupID:    "group",
			id:         "id",
			withError:  ErrContractEmpty,
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
			result, err := client.GetCPCode(context.Background(), test.id, test.contractID, test.groupID)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func mockAPIClient(t *testing.T, mockServer *httptest.Server) PAPI {
	serverURL, err := url.Parse(mockServer.URL)
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	certPool.AddCert(mockServer.Certificate())
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	s, err := session.New(session.WithClient(httpClient), session.WithConfig(&edgegrid.Config{Host: serverURL.Host}))
	assert.NoError(t, err)
	return API(s)
}
