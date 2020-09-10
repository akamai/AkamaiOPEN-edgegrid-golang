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
		params           CPCodeParams
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCPCodesResponse
		withError        error
	}{
		"200 OK": {
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
			},
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
				GroupID:    "group",
				CPCodes: CPCodeItems{Items: []CPCode{
					{
						ID:          "cpcode_id",
						Name:        "cpcode_name",
						CreatedDate: "2020-09-10T15:06:13Z",
						ProductIDs:  []string{"prd_Web_App_Accel"},
					},
				}},
			},
		},
		"404 Not found": {
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
			},
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
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
			},
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
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "",
			},
			withError: ErrGroupEmpty,
		},
		"empty contract ID": {
			params: CPCodeParams{
				ContractID: "",
				GroupID:    "group",
			},
			withError: ErrContractEmpty,
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
			result, err := client.GetCPCodes(context.Background(), test.params)
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
		params           CPCodeParams
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCPCodesResponse
		withError        error
	}{
		"200 OK": {
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
				ID:         "cpcodeID",
			},
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
				GroupID:    "group",
				CPCodes: CPCodeItems{Items: []CPCode{
					{
						ID:          "cpcodeID",
						Name:        "cpcode_name",
						CreatedDate: "2020-09-10T15:06:13Z",
						ProductIDs:  []string{"prd_Web_App_Accel"},
					},
				}},
			},
		},
		"404 Not found": {
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
				ID:         "not_found",
			},
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
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
				ID:         "cpcodeID",
			},
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
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
				ID:         "",
			},
			withError: ErrIDEmpty,
		},
		"empty group ID": {
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "",
				ID:         "cpcodeID",
			},
			withError: ErrGroupEmpty,
		},
		"empty contract ID": {
			params: CPCodeParams{
				ContractID: "",
				GroupID:    "group",
				ID:         "cpcodeID",
			},
			withError: ErrContractEmpty,
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
			result, err := client.GetCPCode(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_CreateCPCode(t *testing.T) {
	tests := map[string]struct {
		params         CPCodeParams
		request        CreateCPCodeRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		expected       *CreateCPCodeResponse
		withError      error
	}{
		"201 Created": {
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
			},
			request: CreateCPCodeRequest{
				ProductID:  "productID",
				CPCodeName: "cpcodeName",
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "cpcodeLink": "/papi/v1/cpcodes/123?contractId=contract-1TJZFW&groupId=group"
}`,
			expectedPath: "/papi/v1/cpcodes?contractId=contract&groupId=group",
			expected: &CreateCPCodeResponse{
				CPCodeLink: "/papi/v1/cpcodes/123?contractId=contract-1TJZFW&groupId=group",
				CPCodeID:   "123",
			},
		},
		"500 Internal Server Error": {
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
			},
			request: CreateCPCodeRequest{
				ProductID:  "productID",
				CPCodeName: "cpcodeName",
			},
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
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "",
			},
			withError: ErrGroupEmpty,
		},
		"empty contract ID": {
			params: CPCodeParams{
				ContractID: "",
				GroupID:    "group",
			},
			withError: ErrContractEmpty,
		},
		"invalid response location": {
			params: CPCodeParams{
				ContractID: "contract",
				GroupID:    "group",
			},
			request: CreateCPCodeRequest{
				ProductID:  "productID",
				CPCodeName: "cpcodeName",
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "cpcodeLink": ":"
}`,
			expectedPath: "/papi/v1/cpcodes?contractId=contract&groupId=group",
			withError:    ErrInvalidLocation,
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
			result, err := client.CreateCPCode(context.Background(), CreateCPCode{test.params, test.request})
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
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
	return New(s)
}
