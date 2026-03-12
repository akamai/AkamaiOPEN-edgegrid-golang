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

func TestCreateReportingGroup(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           CreateReportingGroupRequest
		expectedResponse *CreateReportingGroupResponse
		expectedPath     string
		expectedBody     string
		responseStatus   int
		responseBody     string
		returnedHeaders  map[string]string
		withError        func(*testing.T, error)
	}{
		"201 Created - minimal request": {
			params: CreateReportingGroupRequest{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractCreateModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeCreateModel{
							{
								CpCodeID: 12345,
							},
						},
					},
				},
				ReportingGroupName: "Test Reporting Group",
			},
			returnedHeaders: map[string]string{
				"X-Limit-Max-Reporting-Groups-Limit":     "100",
				"X-Limit-Max-Reporting-Groups-Remaining": "99",
			},
			expectedResponse: &CreateReportingGroupResponse{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "Test CP Code",
							},
						},
					},
				},
				ReportingGroupID:   789,
				ReportingGroupName: "Test Reporting Group",
				ResourceLimits: ResourceLimitsMetadata{
					ReportingGroupsLimitTotal:     ptr.To(int64(100)),
					ReportingGroupsLimitRemaining: ptr.To(int64(99)),
				},
			},
			responseStatus: 201,
			responseBody: `
			{
				"accessGroup": {
					"contractId": "C-0N7RAC7",
					"groupId": 456
				},
				"contracts": [
					{
						"contractId": "C-0N7RAC7",
						"cpcodes": [
							{
								"cpcodeId": 12345,
								"cpcodeName": "Test CP Code"
							}
						]
					}
				],
				"reportingGroupId": 789,
				"reportingGroupName": "Test Reporting Group"
			}`,
			expectedPath: "/cprg/v1/reporting-groups",
		},
		"201 Created - multiple CP codes": {
			params: CreateReportingGroupRequest{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractCreateModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeCreateModel{
							{
								CpCodeID: 12345,
							},
							{
								CpCodeID: 654321,
							},
						},
					},
				},
				ReportingGroupName: "Multi Contract Group",
			},
			returnedHeaders: map[string]string{
				"X-Limit-Max-Reporting-Groups-Limit":     "100",
				"X-Limit-Max-Reporting-Groups-Remaining": "99",
			},
			expectedResponse: &CreateReportingGroupResponse{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "CP Code 1",
							},
							{
								CpCodeID:   654321,
								CpCodeName: "CP Code 2",
							},
						},
					},
				},
				ReportingGroupID:   790,
				ReportingGroupName: "Multi Contract Group",
				ResourceLimits: ResourceLimitsMetadata{
					ReportingGroupsLimitTotal:     ptr.To(int64(100)),
					ReportingGroupsLimitRemaining: ptr.To(int64(99)),
				},
			},
			responseStatus: 201,
			responseBody: `
			{
				"accessGroup": {
					"contractId": "C-0N7RAC7",
					"groupId": 456
				},
				"contracts": [
					{
						"contractId": "C-0N7RAC7",
						"cpcodes": [
							{
								"cpcodeId": 12345,
								"cpcodeName": "CP Code 1"
							},
							{
								"cpcodeId": 654321,
								"cpcodeName": "CP Code 2"
							}
						]
					}
				],
				"reportingGroupId": 790,
				"reportingGroupName": "Multi Contract Group"
			}`,
			expectedPath: "/cprg/v1/reporting-groups",
		},
		"validation error - missing reporting group name": {
			params: CreateReportingGroupRequest{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractCreateModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeCreateModel{
							{
								CpCodeID: 12345,
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "ReportingGroupName")
			},
		},
		"validation error - more than one contract provided in create request": {
			params: CreateReportingGroupRequest{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractCreateModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeCreateModel{
							{
								CpCodeID: 12345,
							},
						},
					},
					{
						ContractID: "C-0N7RAC71",
						CpCodes: []CpCodeCreateModel{
							{
								CpCodeID: 67890,
							},
						},
					},
				}, ReportingGroupName: "Test Reporting Group",
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "Contracts")
			},
		},
		"validation error - missing access group contract ID": {
			params: CreateReportingGroupRequest{
				AccessGroup: AccessGroupModel{
					GroupID: ptr.To(int64(456)),
				},
				Contracts: []ContractCreateModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeCreateModel{
							{
								CpCodeID: 12345,
							},
						},
					},
				},
				ReportingGroupName: "Test Reporting Group",
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "AccessGroup")
			},
		},
		"validation error - missing access group": {
			params: CreateReportingGroupRequest{
				Contracts: []ContractCreateModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeCreateModel{
							{
								CpCodeID: 12345,
							},
						},
					},
				},
				ReportingGroupName: "Test Reporting Group",
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "AccessGroup")
			},
		},
		"validation error - missing contracts": {
			params: CreateReportingGroupRequest{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				ReportingGroupName: "Test Reporting Group",
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "Contracts")
			},
		},
		"500 internal server error": {
			params: CreateReportingGroupRequest{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractCreateModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeCreateModel{
							{
								CpCodeID: 12345,
							},
						},
					},
				},
				ReportingGroupName: "Test Reporting Group",
			},
			responseStatus: 500,
			responseBody: `
			{
				"code": "internal.server.error",
				"title": "Internal Server Error",
				"incidentID": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
			}`,
			expectedPath: "/cprg/v1/reporting-groups",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrCreateReportingGroup, &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
				})
				assert.EqualError(t, err, want.Error())
				assert.ErrorIs(t, err, ErrCreateReportingGroup)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if len(tc.returnedHeaders) > 0 {
					for header, value := range tc.returnedHeaders {
						w.Header().Set(header, value)
					}
				}
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateReportingGroup(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestGetReportingGroup(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           GetReportingGroupsRequest
		expectedResponse *GetReportingGroupResponse
		expectedPath     string
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetReportingGroupsRequest{
				ReportingGroupID: 789,
			},
			expectedResponse: &GetReportingGroupResponse{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "Test CP Code",
							},
						},
					},
				},
				ReportingGroupID:   789,
				ReportingGroupName: "Test Reporting Group",
			},
			responseStatus: 200,
			responseBody: `
			{
				"accessGroup": {
					"contractId": "C-0N7RAC7",
					"groupId": 456
				},
				"contracts": [
					{
						"contractId": "C-0N7RAC7",
						"cpcodes": [
							{
								"cpcodeId": 12345,
								"cpcodeName": "Test CP Code"
							}
						]
					}
				],
				"reportingGroupId": 789,
				"reportingGroupName": "Test Reporting Group"
			}`,
			expectedPath: "/cprg/v1/reporting-groups/789",
		},
		"validation error - missing reporting group ID": {
			params: GetReportingGroupsRequest{},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "ReportingGroupID")
			},
		},
		"500 internal server error": {
			params: GetReportingGroupsRequest{
				ReportingGroupID: 789,
			},
			responseStatus: 500,
			responseBody: `
			{
				"code": "internal.server.error",
				"title": "Internal Server Error",
				"incidentID": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
			}`,
			expectedPath: "/cprg/v1/reporting-groups/789",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrGetReportingGroup, &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
				})
				assert.EqualError(t, err, want.Error())
				assert.ErrorIs(t, err, ErrGetReportingGroup)
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
			result, err := client.GetReportingGroup(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestUpdateReportingGroup(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           UpdateReportingGroupRequest
		expectedResponse *UpdateReportingGroupResponse
		expectedPath     string
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateReportingGroupRequest{
				ReportingGroupID:   789,
				ReportingGroupName: "Updated Reporting Group",
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "Test CP Code",
							},
							{
								CpCodeID:   67890,
								CpCodeName: "Additional CP Code",
							},
						},
					},
				},
			},
			expectedResponse: &UpdateReportingGroupResponse{
				AccessGroup: AccessGroupModel{
					ContractID: "C-0N7RAC7",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "Test CP Code",
							},
							{
								CpCodeID:   67890,
								CpCodeName: "Additional CP Code",
							},
						},
					},
				},
				ReportingGroupID:   789,
				ReportingGroupName: "Updated Reporting Group",
			},
			responseStatus: 200,
			responseBody: `
			{
				"accessGroup": {
					"contractId": "C-0N7RAC7",
					"groupId": 456
				},
				"contracts": [
					{
						"contractId": "C-0N7RAC7",
						"cpcodes": [
							{
								"cpcodeId": 12345,
								"cpcodeName": "Test CP Code"
							},
							{
								"cpcodeId": 67890,
								"cpcodeName": "Additional CP Code"
							}
						]
					}
				],
				"reportingGroupId": 789,
				"reportingGroupName": "Updated Reporting Group"
			}`,
			expectedPath: "/cprg/v1/reporting-groups/789",
		},
		"validation error - missing reporting group ID": {
			params: UpdateReportingGroupRequest{
				ReportingGroupName: "Updated Reporting Group",
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "Test CP Code",
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "ReportingGroupID")
			},
		},
		"validation error - missing contracts": {
			params: UpdateReportingGroupRequest{
				ReportingGroupID:   789,
				ReportingGroupName: "Test Reporting Group",
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "Contracts")
			},
		},
		"validation error - more than one contract provided in update request": {
			params: UpdateReportingGroupRequest{
				ReportingGroupID:   789,
				ReportingGroupName: "Updated Reporting Group",
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "Test CP Code",
							},
						},
					},
					{
						ContractID: "C-0N7RAC71",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   67890,
								CpCodeName: "Test CP Code 2",
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "Contracts")
			},
		},
		"validation error - missing reporting group name": {
			params: UpdateReportingGroupRequest{
				ReportingGroupID: 789,
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "Test CP Code",
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "ReportingGroupName")
			},
		},
		"500 internal server error": {
			params: UpdateReportingGroupRequest{
				ReportingGroupID:   789,
				ReportingGroupName: "Updated Reporting Group",
				Contracts: []ContractModel{
					{
						ContractID: "C-0N7RAC7",
						CpCodes: []CpCodeModel{
							{
								CpCodeID:   12345,
								CpCodeName: "Test CP Code",
							},
						},
					},
				},
			},
			responseStatus: 500,
			responseBody: `
			{
				"code": "internal.server.error",
				"title": "Internal Server Error",
				"incidentID": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
			}`,
			expectedPath: "/cprg/v1/reporting-groups/789",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrUpdateReportingGroup, &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
				})
				assert.EqualError(t, err, want.Error())
				assert.ErrorIs(t, err, ErrUpdateReportingGroup)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateReportingGroup(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestDeleteReportingGroup(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params         DeleteReportingGroupRequest
		expectedPath   string
		responseStatus int
		responseBody   string
		withError      func(*testing.T, error)
	}{
		"204 No Content": {
			params: DeleteReportingGroupRequest{
				ReportingGroupID: 789,
			},
			responseStatus: 204,
			expectedPath:   "/cprg/v1/reporting-groups/789",
		},
		"validation error - missing reporting group ID": {
			params: DeleteReportingGroupRequest{},
			withError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrStructValidation)
				assert.Contains(t, err.Error(), "ReportingGroupID")
			},
		},
		"500 internal server error": {
			params: DeleteReportingGroupRequest{
				ReportingGroupID: 789,
			},
			responseStatus: 500,
			responseBody: `
			{
				"code": "internal.server.error",
				"title": "Internal Server Error",
				"incidentID": "aaaa111-bb22-cc33-dd44-dfeeggda5566"
			}`,
			expectedPath: "/cprg/v1/reporting-groups/789",
			withError: func(t *testing.T, err error) {
				want := fmt.Errorf("%w: %w", ErrDeleteReportingGroup, &Error{
					Code:       "internal.server.error",
					Title:      "Internal Server Error",
					IncidentID: "aaaa111-bb22-cc33-dd44-dfeeggda5566",
					Details:    nil,
				})
				assert.EqualError(t, err, want.Error())
				assert.ErrorIs(t, err, ErrDeleteReportingGroup)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeleteReportingGroup(context.Background(), tc.params)

			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
