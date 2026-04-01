package reportinggroups

import (
	"context"
	"fmt"
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
				want := fmt.Errorf("%w: %w", ErrGetCPCodesWaterMarkLimits, &Error{
					Code:       "bad.request",
					Title:      "Bad Request",
					IncidentID: "123456",
					Details: []SecondaryError{
						{
							Code:    "invalid.data",
							Message: "Invalid contract id:INVALID_CONTRACT_ID",
						},
					},
				})
				assert.EqualError(t, err, want.Error())
				assert.ErrorIs(t, err, ErrGetCPCodesWaterMarkLimits)
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
				want := fmt.Errorf("%w: %w", ErrGetCPCodesWaterMarkLimits, &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
				})
				assert.EqualError(t, err, want.Error())
				assert.ErrorIs(t, err, ErrGetCPCodesWaterMarkLimits)
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
						DefaultTimeZone: nil,
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
						AccessGroup: AccessGroupModel{
							ContractID: "C-0N7RAC7",
							GroupID:    nil,
						},
					},
					{
						CPCodeID:        98765,
						CPCodeName:      "my-second-cp-code",
						Purgeable:       true,
						AccountID:       "A-CCT1234",
						DefaultTimeZone: ptr.To("GMT 0 (Greenwich Mean Time)"),
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
						AccessGroup: AccessGroupModel{
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
						DefaultTimeZone: ptr.To("GMT 0 (Greenwich Mean Time)"),
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
						AccessGroup: AccessGroupModel{
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
				want := fmt.Errorf("%w: %w", ErrListCPCodes, &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
				})
				assert.EqualError(t, err, want.Error())
				assert.ErrorIs(t, err, ErrListCPCodes)
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
