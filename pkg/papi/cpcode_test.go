package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"
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

func TestGetCPCodeDetail(t *testing.T) {
	tests := map[string]struct {
		id               int
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CPCodeDetailResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			id:             123,
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cpcodeId": 123,
    "cpcodeName": "test-cp-code",
    "purgeable": true,
    "accountId": "test-account-id",
    "defaultTimezone": "GMT 0 (Greenwich Mean Time)",
    "overrideTimezone": {
        "timezoneId": "0",
        "timezoneValue": "GMT 0 (Greenwich Mean Time)"
    },
    "type": "Regular",
    "contracts": [
        {
            "contractId": "test-contract-id",
            "status": "ongoing"
        }
    ],
    "products": [
        {
            "productId": "test-product-id",
            "productName": "test-product-name"
        }
    ],
    "accessGroup": {
        "groupId": null,
        "contractId": "test-contract-id"
    }
}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &CPCodeDetailResponse{
				ID:              123,
				Name:            "test-cp-code",
				Purgeable:       true,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "0",
					TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
		},
		"500 internal server error": {
			id:             123,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Server Error",
    "detail": "Error fetching cp codes",
    "status": 500
}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Server Error",
					Detail:     "Error fetching cp codes",
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
			result, err := client.GetCPCodeDetail(context.Background(), test.id)
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

func TestUpdateCPCode(t *testing.T) {
	tests := map[string]struct {
		params           UpdateCPCodeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CPCodeDetailResponse
		withError        func(*testing.T, error)
	}{
		"200 OK Update name": {
			params: UpdateCPCodeRequest{
				ID:   123,
				Name: "test-cp-code-updated",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cpcodeId": 123,
    "cpcodeName": "test-cp-code-updated",
    "purgeable": true,
    "accountId": "test-account-id",
    "defaultTimezone": "GMT 0 (Greenwich Mean Time)",
    "overrideTimezone": {
        "timezoneId": "0",
        "timezoneValue": "GMT 0 (Greenwich Mean Time)"
    },
    "type": "Regular",
    "contracts": [
        {
            "contractId": "test-contract-id",
            "status": "ongoing"
        }
    ],
    "products": [
        {
            "productId": "test-product-id",
            "productName": "test-product-name"
        }
    ],
    "accessGroup": {
        "groupId": null,
        "contractId": "test-contract-id"
    }
}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &CPCodeDetailResponse{
				ID:              123,
				Name:            "test-cp-code-updated",
				Purgeable:       true,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "0",
					TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
		},
		"200 OK Update time zone": {
			params: UpdateCPCodeRequest{
				ID:   123,
				Name: "test-cp-code",
				OverrideTimeZone: &CPCodeTimeZone{
					TimeZoneID:    "1",
					TimeZoneValue: "GMT + 1",
				},
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cpcodeId": 123,
    "cpcodeName": "test-cp-code-updated",
    "purgeable": true,
    "accountId": "test-account-id",
    "defaultTimezone": "GMT 0 (Greenwich Mean Time)",
    "overrideTimezone": {
        "timezoneId": "1",
        "timezoneValue": "GMT + 1"
    },
    "type": "Regular",
    "contracts": [
        {
            "contractId": "test-contract-id",
            "status": "ongoing"
        }
    ],
    "products": [
        {
            "productId": "test-product-id",
            "productName": "test-product-name"
        }
    ],
    "accessGroup": {
        "groupId": null,
        "contractId": "test-contract-id"
    }
}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &CPCodeDetailResponse{
				ID:              123,
				Name:            "test-cp-code-updated",
				Purgeable:       true,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "1",
					TimeZoneValue: "GMT + 1",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
		},
		"200 OK Update purgeable": {
			params: UpdateCPCodeRequest{
				ID:        123,
				Name:      "test-cp-code",
				Purgeable: tools.BoolPtr(false),
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cpcodeId": 123,
    "cpcodeName": "test-cp-code-updated",
    "purgeable": false,
    "accountId": "test-account-id",
    "defaultTimezone": "GMT 0 (Greenwich Mean Time)",
    "overrideTimezone": {
        "timezoneId": "0",
        "timezoneValue": "GMT 0 (Greenwich Mean Time)"
    },
    "type": "Regular",
    "contracts": [
        {
            "contractId": "test-contract-id",
            "status": "ongoing"
        }
    ],
    "products": [
        {
            "productId": "test-product-id",
            "productName": "test-product-name"
        }
    ],
    "accessGroup": {
        "groupId": null,
        "contractId": "test-contract-id"
    }
}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &CPCodeDetailResponse{
				ID:              123,
				Name:            "test-cp-code-updated",
				Purgeable:       false,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "0",
					TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
		},
		"500 internal server error": {
			params: UpdateCPCodeRequest{
				ID:   123,
				Name: "test-cp-code-updated",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Server Error",
    "detail": "Error updating cp code",
    "status": 500
}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Server Error",
					Detail:     "Error updating cp code",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation - id is required": {
			params: UpdateCPCodeRequest{
				Name: "test-cp-code-updated",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation - name is required": {
			params: UpdateCPCodeRequest{
				ID: 123,
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation - contracts is required": {
			params: UpdateCPCodeRequest{
				ID:   123,
				Name: "test-cp-code-updated",
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation - contract id is required": {
			params: UpdateCPCodeRequest{
				ID:   123,
				Name: "test-cp-code-updated",
				Contracts: []CPCodeContract{
					{
						Status: "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation - products is required": {
			params: UpdateCPCodeRequest{
				ID:   123,
				Name: "test-cp-code-updated",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
			},
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation - product id is required": {
			params: UpdateCPCodeRequest{
				ID:   123,
				Name: "test-cp-code-updated",
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductName: "test-product-name",
					},
				},
			},
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation - time zone id is required": {
			params: UpdateCPCodeRequest{
				ID:   123,
				Name: "test-cp-code-updated",
				OverrideTimeZone: &CPCodeTimeZone{
					TimeZoneValue: "GMT + 1",
				},
				Contracts: []CPCodeContract{
					{
						ContractID: "test-contract-id",
						Status:     "ongoing",
					},
				},
				Products: []CPCodeProduct{
					{
						ProductID:   "test-product-id",
						ProductName: "test-product-name",
					},
				},
			},
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
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
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateCPCode(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
