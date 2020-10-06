package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapi_GetCPCodes(t *testing.T) {
	tests := map[string]struct {
		params           GetCPCodesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCPCodesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetCPCodesRequest{
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
		"500 internal server error": {
			params: GetCPCodesRequest{
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
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching cp codes",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty group ID": {
			params: GetCPCodesRequest{
				ContractID: "contract",
				GroupID:    "",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "GroupID")
			},
		},
		"empty contract ID": {
			params: GetCPCodesRequest{
				ContractID: "",
				GroupID:    "group",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContractID")
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
			result, err := client.GetCPCodes(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_GetCPCode(t *testing.T) {
	tests := map[string]struct {
		params           GetCPCodeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCPCodesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetCPCodeRequest{
				ContractID: "contract",
				GroupID:    "group",
				CPCodeID:   "cpcodeID",
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
				CPCode: CPCode{
					ID:          "cpcodeID",
					Name:        "cpcode_name",
					CreatedDate: "2020-09-10T15:06:13Z",
					ProductIDs:  []string{"prd_Web_App_Accel"},
				},
			},
		},
		"CP Code not found": {
			params: GetCPCodeRequest{
				ContractID: "contract",
				GroupID:    "group",
				CPCodeID:   "cpcodeID",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "acc",
    "contractId": "contract",
    "groupId": "group",
    "cpcodes": {
        "items": [
        ]
    }
}`,
			expectedPath: "/papi/v1/cpcodes/cpcodeID?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrNotFound), "want: %v; got: %v", ErrNotFound, err)
			},
		},
		"500 internal server error": {
			params: GetCPCodeRequest{
				ContractID: "contract",
				GroupID:    "group",
				CPCodeID:   "cpcodeID",
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
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching cp codes",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty cpcode ID": {
			params: GetCPCodeRequest{
				ContractID: "contract",
				GroupID:    "group",
				CPCodeID:   "",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "CPCodeID")
			},
		},
		"empty group ID": {
			params: GetCPCodeRequest{
				ContractID: "contract",
				GroupID:    "",
				CPCodeID:   "cpcodeID",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "GroupID")
			},
		},
		"empty contract ID": {
			params: GetCPCodeRequest{
				ContractID: "",
				GroupID:    "group",
				CPCodeID:   "cpcodeID",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContractID")
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
			result, err := client.GetCPCode(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_CreateCPCode(t *testing.T) {
	tests := map[string]struct {
		params         CreateCPCodeRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		expected       *CreateCPCodeResponse
		withError      func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateCPCodeRequest{
				ContractID: "contract",
				GroupID:    "group",
				CPCode: CreateCPCode{
					ProductID:  "productID",
					CPCodeName: "cpcodeName",
				},
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
			params: CreateCPCodeRequest{
				ContractID: "contract",
				GroupID:    "group",
				CPCode: CreateCPCode{
					ProductID:  "productID",
					CPCodeName: "cpcodeName",
				},
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
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching cp codes",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty group ID": {
			params: CreateCPCodeRequest{
				ContractID: "contract",
				GroupID:    "",
				CPCode: CreateCPCode{
					ProductID:  "productID",
					CPCodeName: "cpCodeName",
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "GroupID")
			},
		},
		"empty contract ID": {
			params: CreateCPCodeRequest{
				ContractID: "",
				GroupID:    "group",
				CPCode: CreateCPCode{
					ProductID:  "productID",
					CPCodeName: "cpCodeName",
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContractID")
			},
		},
		"empty product ID": {
			params: CreateCPCodeRequest{
				ContractID: "contractID",
				GroupID:    "group",
				CPCode: CreateCPCode{
					ProductID:  "",
					CPCodeName: "cpCodeName",
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ProductID")
			},
		},
		"empty cp code name": {
			params: CreateCPCodeRequest{
				ContractID: "",
				GroupID:    "group",
				CPCode: CreateCPCode{
					ProductID:  "productID",
					CPCodeName: "",
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "CPCodeName")
			},
		},
		"invalid response location": {
			params: CreateCPCodeRequest{
				ContractID: "contract",
				GroupID:    "group",
				CPCode: CreateCPCode{
					ProductID:  "productID",
					CPCodeName: "cpcodeName",
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "cpcodeLink": ":"
}`,
			expectedPath: "/papi/v1/cpcodes?contractId=contract&groupId=group",
			withError: func(t *testing.T, err error) {
				want := ErrInvalidResponseLink
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
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
			result, err := client.CreateCPCode(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
