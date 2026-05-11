package reportinggroups

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCPCodesWaterMarkLimits(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           GetCPCodesWaterMarkLimitsRequest
		expectedResponse *GetCPCodesWaterMarkLimitsResponse
		expectedPath     string
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetCPCodesWaterMarkLimitsRequest{
				ContractID: "C-0N7RAC7",
			},
			expectedResponse: &GetCPCodesWaterMarkLimitsResponse{
				CurrentCapacity: 5,
				Limit:           100,
				LimitType:       "account",
			},
			responseStatus: http.StatusOK,
			responseBody: `
			{
				"currentCapacity": 5,
				"limit": 100,
				"limitType": "account"
			}`,
			expectedPath: "/cprg/v1/cpcodes/contracts/C-0N7RAC7/watermark-limits",
		},
		"validation error - empty contract ID": {
			params: GetCPCodesWaterMarkLimitsRequest{},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "ContractID")
			},
		},
		"400 bad request - invalid contract ID": {
			params: GetCPCodesWaterMarkLimitsRequest{
				ContractID: "INVALID_CONTRACT_ID",
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
			{
				"code": "bad.request",
				"title": "Bad Request",
				"incidentId": "123456",
				"details": [
					{
						"code": "invalid.data",
						"message": "Invalid contract id:INVALID_CONTRACT_ID"
					}
				]
			}`,
			expectedPath: "/cprg/v1/cpcodes/contracts/INVALID_CONTRACT_ID/watermark-limits",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "bad.request",
					Title:      "Bad Request",
					IncidentID: "123456",
					Details: []SecondaryError{
						{
							Code:    "invalid.data",
							Message: "Invalid contract id:INVALID_CONTRACT_ID",
						},
					},
					HTTPStatus: 400,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrGetCPCodesWaterMarkLimits)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrGetCPCodesWaterMarkLimits, want).Error())
			},
		},
		"500 internal server error": {
			params: GetCPCodesWaterMarkLimitsRequest{
				ContractID: "C-0N7RAC7",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"code": "internal.server.error",
				"title": "Internal Server Error",
				"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
			}`,
			expectedPath: "/cprg/v1/cpcodes/contracts/C-0N7RAC7/watermark-limits",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
					HTTPStatus: 500,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrGetCPCodesWaterMarkLimits)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrGetCPCodesWaterMarkLimits, want).Error())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCPCodesWaterMarkLimits(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestListCPCodes(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           ListCPCodesRequest
		expectedResponse *ListCPCodesResponse
		expectedPath     string
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK - no params": {
			params:         ListCPCodesRequest{},
			responseStatus: http.StatusOK,
			responseBody: `{
				"cpcodes": [
					{
						"cpcodeId": 12345,
						"cpcodeName": "my-cp-code",
						"purgeable": true,
						"accountId": "A-CCT1234",
						"defaultTimezone": null,
						"overrideTimezone": {
							"timezoneId": "0",
							"timezoneValue": "GMT 0 (Greenwich Mean Time)"
						},
						"type": "Regular",
						"contracts": [
							{
								"contractId": "C-0N7RAC7",
								"status": "ongoing"
							}
						],
						"products": [
							{
								"productId": "prd_Object_Delivery",
								"productName": "Object Delivery"
							}
						],
						"accessGroup": {
							"groupId": null,
							"contractId": "C-0N7RAC7"
						}
					},
					{
						"cpcodeId": 98765,
						"cpcodeName": "my-second-cp-code",
						"purgeable": true,
						"accountId": "A-CCT1234",
						"defaultTimezone": "GMT 0 (Greenwich Mean Time)",
						"overrideTimezone": {
							"timezoneId": "0",
							"timezoneValue": "GMT 0 (Greenwich Mean Time)"
						},
						"type": "Regular",
						"contracts": [
							{
								"contractId": "C-0N7RAC7",
								"status": "ongoing"
							}
						],
						"products": [
							{
								"productId": "prd_Site_Del",
								"productName": "Site Delivery"
							}
						],
						"accessGroup": {
							"groupId": null,
							"contractId": "C-0N7RAC7"
						}
					}
				]
			}`,
			expectedPath: "/cprg/v1/cpcodes",
			expectedResponse: &ListCPCodesResponse{
				CPCodes: []CPCodeDetails{
					{
						CPCodeID:        12345,
						CPCodeName:      "my-cp-code",
						Purgeable:       true,
						AccountID:       "A-CCT1234",
						DefaultTimeZone: "",
						OverrideTimeZone: CPCodeTimeZone{
							TimeZoneID:    "0",
							TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
						},
						Type: "Regular",
						Contracts: []CPCodeContract{
							{ContractID: "C-0N7RAC7", Status: "ongoing"},
						},
						Products: []Product{
							{ProductID: "prd_Object_Delivery", ProductName: "Object Delivery"},
						},
						AccessGroup: AccessGroup{
							ContractID: "C-0N7RAC7",
							GroupID:    nil,
						},
					},
					{
						CPCodeID:        98765,
						CPCodeName:      "my-second-cp-code",
						Purgeable:       true,
						AccountID:       "A-CCT1234",
						DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
						OverrideTimeZone: CPCodeTimeZone{
							TimeZoneID:    "0",
							TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
						},
						Type: "Regular",
						Contracts: []CPCodeContract{
							{ContractID: "C-0N7RAC7", Status: "ongoing"},
						},
						Products: []Product{
							{ProductID: "prd_Site_Del", ProductName: "Site Delivery"},
						},
						AccessGroup: AccessGroup{
							ContractID: "C-0N7RAC7",
							GroupID:    nil,
						},
					},
				},
			},
		},
		"200 OK - all params": {
			params: ListCPCodesRequest{
				ContractID: "C-0N7RAC7",
				GroupID:    "12345",
				ProductID:  "prd_Object_Delivery",
				CPCodeName: "my-cp-code",
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"cpcodes": [
					{
						"cpcodeId": 12345,
						"cpcodeName": "my-cp-code",
						"purgeable": true,
						"accountId": "A-CCT1234",
						"defaultTimezone": "GMT 0 (Greenwich Mean Time)",
						"overrideTimezone": {
							"timezoneId": "0",
							"timezoneValue": "GMT 0 (Greenwich Mean Time)"
						},
						"type": "Regular",
						"contracts": [
							{
								"contractId": "C-0N7RAC7",
								"status": "ongoing"
							}
						],
						"products": [
							{
								"productId": "prd_Object_Delivery",
								"productName": "Object Delivery"
							}
						],
						"accessGroup": {
							"groupId": null,
							"contractId": "C-0N7RAC7"
						}
					}
				]
			}`,
			expectedPath: "/cprg/v1/cpcodes?contractId=C-0N7RAC7&cpcodeName=my-cp-code&groupId=12345&productId=prd_Object_Delivery",
			expectedResponse: &ListCPCodesResponse{
				CPCodes: []CPCodeDetails{
					{
						CPCodeID:        12345,
						CPCodeName:      "my-cp-code",
						Purgeable:       true,
						AccountID:       "A-CCT1234",
						DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
						OverrideTimeZone: CPCodeTimeZone{
							TimeZoneID:    "0",
							TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
						},
						Type: "Regular",
						Contracts: []CPCodeContract{
							{ContractID: "C-0N7RAC7", Status: "ongoing"},
						},
						Products: []Product{
							{ProductID: "prd_Object_Delivery", ProductName: "Object Delivery"},
						},
						AccessGroup: AccessGroup{
							ContractID: "C-0N7RAC7",
							GroupID:    nil,
						},
					},
				},
			},
		},
		"500 internal server error": {
			params:         ListCPCodesRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"code": "internal.server.error",
				"title": "Internal Server Error",
				"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
			}`,
			expectedPath: "/cprg/v1/cpcodes",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
					HTTPStatus: 500,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrListCPCodes)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrListCPCodes, want).Error())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListCPCodes(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetCPCode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           GetCPCodeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCPCodeResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         GetCPCodeRequest{CPCodeID: 123},
			responseStatus: http.StatusOK,
			responseBody: `{
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
					"contractId": "test-contract-id",
					"groupId": 12345
				}
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &GetCPCodeResponse{
				CPCodeID:        123,
				CPCodeName:      "test-cp-code",
				Purgeable:       true,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "0",
					TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
				AccessGroup: AccessGroup{
					ContractID: "test-contract-id",
					GroupID:    ptr.To(int64(12345)),
				},
			},
		},
		"200 OK - null defaultTimezone": {
			params:         GetCPCodeRequest{CPCodeID: 123},
			responseStatus: http.StatusOK,
			responseBody: `{
				"cpcodeId": 123,
				"cpcodeName": "test-cp-code",
				"purgeable": true,
				"accountId": "test-account-id",
				"defaultTimezone": null,
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
					"contractId": "test-contract-id",
					"groupId": 12345
				}
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &GetCPCodeResponse{
				CPCodeID:        123,
				CPCodeName:      "test-cp-code",
				Purgeable:       true,
				AccountID:       "test-account-id",
				DefaultTimeZone: "",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "0",
					TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
				AccessGroup: AccessGroup{
					ContractID: "test-contract-id",
					GroupID:    ptr.To(int64(12345)),
				},
			},
		},
		"200 OK - null groupID in AccessGroup": {
			params:         GetCPCodeRequest{CPCodeID: 123},
			responseStatus: http.StatusOK,
			responseBody: `{
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
					"contractId": "test-contract-id",
					"groupId": null
				}
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &GetCPCodeResponse{
				CPCodeID:        123,
				CPCodeName:      "test-cp-code",
				Purgeable:       true,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "0",
					TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
				AccessGroup: AccessGroup{
					ContractID: "test-contract-id",
					GroupID:    nil,
				},
			},
		},
		"validation error - missing CPCodeID": {
			params: GetCPCodeRequest{},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.EqualError(t, err, "get CP code detail: struct validation: CPCodeID: cannot be blank")
			},
		},
		"403 forbidden": {
			params:         GetCPCodeRequest{CPCodeID: 123},
			responseStatus: http.StatusForbidden,
			responseBody: `{
				"code": "forbidden",
				"title": "Forbidden",
				"incidentId": "19f75985-a7e3-4eb5-9ab6-d0d92f149dfe",
				"details": [
					{
						"code": "invalid.role",
						"message": " User is not authorized to access cpcode"
					}
				]
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "forbidden",
					Title:      "Forbidden",
					IncidentID: "19f75985-a7e3-4eb5-9ab6-d0d92f149dfe",
					Details: []SecondaryError{
						{
							Code:    "invalid.role",
							Message: " User is not authorized to access cpcode",
						},
					},
					HTTPStatus: 403,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrGetCPCode)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrGetCPCode, want).Error())
			},
		},
		"500 internal server error": {
			params:         GetCPCodeRequest{CPCodeID: 123},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"code": "internal.server.error",
				"title": "Internal Server Error",
				"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
					HTTPStatus: 500,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrGetCPCode)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrGetCPCode, want).Error())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetCPCode(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestUpdateCPCode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           UpdateCPCodeRequest
		expectedResponse *UpdateCPCodeResponse
		expectedPath     string
		expectedBody     string
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK - update name": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
			},
			expectedBody:   `{"cpcodeId":123,"cpcodeName":"updated-cp-code","contracts":[{"contractId":"test-contract-id","status":"ongoing"}],"products":[{"productId":"test-product-id","productName":"test-product-name"}]}`,
			responseStatus: http.StatusOK,
			responseBody: `{
				"cpcodeId": 123,
				"cpcodeName": "updated-cp-code",
				"purgeable": true,
				"accountId": "test-account-id",
				"defaultTimezone": "GMT 0 (Greenwich Mean Time)",
				"overrideTimezone": {
					"timezoneId": "0",
					"timezoneValue": "GMT 0 (Greenwich Mean Time)"
				},
				"type": "Regular",
				"contracts": [
					{"contractId": "test-contract-id", "status": "ongoing"}
				],
				"products": [
					{"productId": "test-product-id", "productName": "test-product-name"}
				],
				"accessGroup": {
					"contractId": "test-contract-id",
					"groupId": 12345
				}
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &UpdateCPCodeResponse{
				CPCodeID:        123,
				CPCodeName:      "updated-cp-code",
				Purgeable:       true,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "0",
					TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
				AccessGroup: AccessGroup{
					ContractID: "test-contract-id",
					GroupID:    ptr.To(int64(12345)),
				},
			},
		},
		"200 OK - update timezone": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "test-cp-code",
				OverrideTimeZone: &CPCodeTimeZone{
					TimeZoneID:    "1",
					TimeZoneValue: "GMT + 1",
				},
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
			},
			expectedBody:   `{"cpcodeId":123,"cpcodeName":"test-cp-code","overrideTimezone":{"timezoneId":"1","timezoneValue":"GMT + 1"},"contracts":[{"contractId":"test-contract-id","status":"ongoing"}],"products":[{"productId":"test-product-id","productName":"test-product-name"}]}`,
			responseStatus: http.StatusOK,
			responseBody: `{
				"cpcodeId": 123,
				"cpcodeName": "test-cp-code",
				"purgeable": true,
				"accountId": "test-account-id",
				"defaultTimezone": "GMT 0 (Greenwich Mean Time)",
				"overrideTimezone": {
					"timezoneId": "1",
					"timezoneValue": "GMT + 1"
				},
				"type": "Regular",
				"contracts": [
					{"contractId": "test-contract-id", "status": "ongoing"}
				],
				"products": [
					{"productId": "test-product-id", "productName": "test-product-name"}
				],
				"accessGroup": {
					"contractId": "test-contract-id",
					"groupId": 12345
				}
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &UpdateCPCodeResponse{
				CPCodeID:        123,
				CPCodeName:      "test-cp-code",
				Purgeable:       true,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "1",
					TimeZoneValue: "GMT + 1",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
				AccessGroup: AccessGroup{
					ContractID: "test-contract-id",
					GroupID:    ptr.To(int64(12345)),
				},
			},
		},
		"200 OK - update purgeable": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "test-cp-code",
				Purgeable:  ptr.To(false),
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
			},
			expectedBody:   `{"cpcodeId":123,"cpcodeName":"test-cp-code","purgeable":false,"contracts":[{"contractId":"test-contract-id","status":"ongoing"}],"products":[{"productId":"test-product-id","productName":"test-product-name"}]}`,
			responseStatus: http.StatusOK,
			responseBody: `{
				"cpcodeId": 123,
				"cpcodeName": "test-cp-code",
				"purgeable": false,
				"accountId": "test-account-id",
				"defaultTimezone": "GMT 0 (Greenwich Mean Time)",
				"overrideTimezone": {
					"timezoneId": "0",
					"timezoneValue": "GMT 0 (Greenwich Mean Time)"
				},
				"type": "Regular",
				"contracts": [
					{"contractId": "test-contract-id", "status": "ongoing"}
				],
				"products": [
					{"productId": "test-product-id", "productName": "test-product-name"}
				],
				"accessGroup": {
					"contractId": "test-contract-id",
					"groupId": 12345
				}
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			expectedResponse: &UpdateCPCodeResponse{
				CPCodeID:        123,
				CPCodeName:      "test-cp-code",
				Purgeable:       false,
				AccountID:       "test-account-id",
				DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
				OverrideTimeZone: CPCodeTimeZone{
					TimeZoneID:    "0",
					TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
				},
				Type: "Regular",
				Contracts: []CPCodeContract{
					{ContractID: "test-contract-id", Status: "ongoing"},
				},
				Products: []Product{
					{ProductID: "test-product-id", ProductName: "test-product-name"},
				},
				AccessGroup: AccessGroup{
					ContractID: "test-contract-id",
					GroupID:    ptr.To(int64(12345)),
				},
			},
		},
		"validation error - missing CPCodeID": {
			params: UpdateCPCodeRequest{
				CPCodeName: "updated-cp-code",
				Contracts:  []CPCodeContract{{ContractID: "test-contract-id"}},
				Products:   []Product{{ProductID: "test-product-id"}},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.EqualError(t, err, "update CP code: struct validation: CPCodeID: cannot be blank")
			},
		},
		"validation error - missing CPCodeName": {
			params: UpdateCPCodeRequest{
				CPCodeID:  123,
				Contracts: []CPCodeContract{{ContractID: "test-contract-id"}},
				Products:  []Product{{ProductID: "test-product-id"}},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.EqualError(t, err, "update CP code: struct validation: CPCodeName: cannot be blank")
			},
		},
		"validation error - contracts is required": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Products:   []Product{{ProductID: "test-product-id"}},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.EqualError(t, err, "update CP code: struct validation: Contracts: cannot be blank")
			},
		},
		"validation error - contract ID is required": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Contracts:  []CPCodeContract{{Status: "ongoing"}},
				Products:   []Product{{ProductID: "test-product-id"}},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Regexp(t, `(?s)update CP code: struct validation: Contracts\[0\]:.*ContractID: cannot be blank`, err.Error())
			},
		},
		"validation error - products is required": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Contracts:  []CPCodeContract{{ContractID: "test-contract-id"}},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.EqualError(t, err, "update CP code: struct validation: Products: cannot be blank")
			},
		},
		"validation error - product ID is required": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Contracts:  []CPCodeContract{{ContractID: "test-contract-id"}},
				Products:   []Product{{ProductName: "test-product-name"}},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Regexp(t, `(?s)update CP code: struct validation: Products\[0\]:.*ProductID: cannot be blank`, err.Error())
			},
		},
		"validation error - timezone ID is required when override timezone is provided": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				OverrideTimeZone: &CPCodeTimeZone{
					TimeZoneValue: "GMT + 1",
				},
				Contracts: []CPCodeContract{{ContractID: "test-contract-id"}},
				Products:  []Product{{ProductID: "test-product-id"}},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Regexp(t, `(?s)update CP code: struct validation: OverrideTimeZone:.*TimeZoneID: cannot be blank`, err.Error())
			},
		},
		"400 bad request - contract ID not found on server": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Contracts:  []CPCodeContract{{ContractID: "C-INVALID", Status: "ongoing"}},
				Products:   []Product{{ProductID: "test-product-id", ProductName: "test-product-name"}},
			},
			expectedBody:   `{"cpcodeId":123,"cpcodeName":"updated-cp-code","contracts":[{"contractId":"C-INVALID","status":"ongoing"}],"products":[{"productId":"test-product-id","productName":"test-product-name"}]}`,
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"code": "bad.request",
				"title": "Bad Request",
				"incidentId": "8b953d5b-ccc6-455e-baa8-b3d4287f5ded",
				"details": [
					{
						"code": "invalid.data",
						"message": "At least one active contract needs to be present"
					}
				]
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "bad.request",
					Title:      "Bad Request",
					IncidentID: "8b953d5b-ccc6-455e-baa8-b3d4287f5ded",
					Details: []SecondaryError{
						{
							Code:    "invalid.data",
							Message: "At least one active contract needs to be present",
						},
					},
					HTTPStatus: 400,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrUpdateCPCode)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrUpdateCPCode, want).Error())
			},
		},
		"400 bad request - product ID not found on server": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Contracts:  []CPCodeContract{{ContractID: "test-contract-id", Status: "ongoing"}},
				Products:   []Product{{ProductID: "P-INVALID", ProductName: "test-product-name"}},
			},
			expectedBody:   `{"cpcodeId":123,"cpcodeName":"updated-cp-code","contracts":[{"contractId":"test-contract-id","status":"ongoing"}],"products":[{"productId":"P-INVALID","productName":"test-product-name"}]}`,
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"code": "bad.request",
				"title": "Bad Request",
				"incidentId": "756ef52c-1049-4c91-9121-be20dc6125f9",
				"details": [
					{
						"code": "invalid.data",
						"message": "At least one service should be present for each contract"
					}
				]
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "bad.request",
					Title:      "Bad Request",
					IncidentID: "756ef52c-1049-4c91-9121-be20dc6125f9",
					Details: []SecondaryError{
						{
							Code:    "invalid.data",
							Message: "At least one service should be present for each contract",
						},
					},
					HTTPStatus: 400,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrUpdateCPCode)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrUpdateCPCode, want).Error())
			},
		},
		"403 forbidden": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Contracts:  []CPCodeContract{{ContractID: "test-contract-id", Status: "ongoing"}},
				Products:   []Product{{ProductID: "test-product-id", ProductName: "test-product-name"}},
			},
			responseStatus: http.StatusForbidden,
			responseBody: `{
				"code": "forbidden",
				"title": "Forbidden",
				"incidentId": "19f75985-a7e3-4eb5-9ab6-d0d92f149dfe",
				"details": [
					{
						"code": "invalid.role",
						"message": " User is not authorized to access cpcode"
					}
				]
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "forbidden",
					Title:      "Forbidden",
					IncidentID: "19f75985-a7e3-4eb5-9ab6-d0d92f149dfe",
					Details: []SecondaryError{
						{
							Code:    "invalid.role",
							Message: " User is not authorized to access cpcode",
						},
					},
					HTTPStatus: 403,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrUpdateCPCode)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrUpdateCPCode, want).Error())
			},
		},
		"500 internal server error": {
			params: UpdateCPCodeRequest{
				CPCodeID:   123,
				CPCodeName: "updated-cp-code",
				Contracts:  []CPCodeContract{{ContractID: "test-contract-id", Status: "ongoing"}},
				Products:   []Product{{ProductID: "test-product-id", ProductName: "test-product-name"}},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"code": "internal.server.error",
				"title": "Internal Server Error",
				"incidentId": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
			}`,
			expectedPath: "/cprg/v1/cpcodes/123",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
					HTTPStatus: 500,
				}
				assert.ErrorIs(t, err, want)
				assert.ErrorIs(t, err, ErrUpdateCPCode)
				assert.EqualError(t, err, fmt.Errorf("%w: %w", ErrUpdateCPCode, want).Error())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if tc.expectedBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, tc.expectedBody, string(body))
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateCPCode(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}
