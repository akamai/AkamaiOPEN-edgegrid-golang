package v3

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPolicies(t *testing.T) {
	tests := map[string]struct {
		params           ListPoliciesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListPoliciesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - two policies, one minimal and one with activation data": {
			params:         ListPoliciesRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "content": [
        {
            "cloudletType": "CD",
            "createdBy": "User1",
            "createdDate": "2023-10-23T11:21:19.896Z",
            "currentActivations": {
                "production": {
                    "effective": null,
                    "latest": null
                },
                "staging": {
                    "effective": null,
                    "latest": null
                }
            },
            "description": null,
            "groupId": 1,
            "id": 11,
            "links": [
                {
                    "href": "Link1",
                    "rel": "self"
                }
            ],
            "modifiedBy": "User2",
            "modifiedDate": "2023-10-23T11:21:19.896Z",
            "name": "Name1",
            "policyType": "SHARED"
        },
        {
            "cloudletType": "CD",
            "createdBy": "User1",
            "createdDate": "2023-10-23T11:21:19.896Z",
            "currentActivations": {
                "production": {
                    "effective": {
                        "createdBy": "User1",
						"createdDate": "2023-10-23T11:21:19.896Z",
						"finishDate": "2023-10-23T11:22:57.589Z",
                        "id": 123,
                        "links": [
                            {
                                "href": "Link1",
                                "rel": "self"
                            }
                        ],
                        "network": "PRODUCTION",
                        "operation": "ACTIVATION",
                        "policyId": 1234,
                        "policyVersion": 1,
                        "policyVersionDeleted": false,
                        "status": "SUCCESS"
                    },
                    "latest": {
                        "createdBy": "User1",
						"createdDate": "2023-10-23T11:21:19.896Z",
						"finishDate": "2023-10-23T11:22:57.589Z",
                        "id": 321,
                        "links": [
                            {
                                "href": "Link2",
                                "rel": "self"
                            }
                        ],
                        "network": "PRODUCTION",
                        "operation": "ACTIVATION",
                        "policyId": 4321,
                        "policyVersion": 1,
                        "policyVersionDeleted": false,
                        "status": "SUCCESS"
                    }
                },
                "staging": {
                    "effective": {
                        "createdBy": "User3",
						"createdDate": "2023-10-23T11:21:19.896Z",
						"finishDate": "2023-10-23T11:22:57.589Z",
                        "id": 789,
                        "links": [
                            {
                                "href": "Link3",
                                "rel": "self"
                            }
                        ],
                        "network": "STAGING",
                        "operation": "ACTIVATION",
                        "policyId": 6789,
                        "policyVersion": 1,
                        "policyVersionDeleted": false,
                        "status": "SUCCESS"
                    },
                    "latest": {
                        "createdBy": "User3",
						"createdDate": "2023-10-23T11:21:19.896Z",
						"finishDate": "2023-10-23T11:22:57.589Z",
                        "id": 987,
                        "links": [
                            {
                                "href": "Link4",
                                "rel": "self"
                            }
                        ],
                        "network": "STAGING",
                        "operation": "ACTIVATION",
                        "policyId": 9876,
                        "policyVersion": 1,
                        "policyVersionDeleted": false,
                        "status": "SUCCESS"
                    }
                }
            },
            "description": "Test",
            "groupId": 1,
            "id": 22,
            "links": [
                {
                    "href": "Link5",
                    "rel": "self"
                }
            ],
            "modifiedBy": "User1",
            "modifiedDate": "2023-10-23T11:21:19.896Z",
            "name": "TestName",
            "policyType": "SHARED"
        }
	],
    "links": [
        {
            "href": "/cloudlets/v3/policies?page=0&size=1000",
            "rel": "self"
        }
    ],
    "page": {
        "number": 0,
        "size": 1000,
        "totalElements": 54,
        "totalPages": 1
    }
}
`,
			expectedPath: "/cloudlets/v3/policies",
			expectedResponse: &ListPoliciesResponse{
				Content: []Policy{
					{
						CloudletType:       CloudletTypeCD,
						CreatedBy:          "User1",
						CreatedDate:        test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
						CurrentActivations: CurrentActivations{},
						Description:        nil,
						GroupID:            1,
						ID:                 11,
						Links: []Link{
							{
								Href: "Link1",
								Rel:  "self",
							},
						},
						ModifiedBy:   "User2",
						ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
						Name:         "Name1",
						PolicyType:   PolicyTypeShared,
					},
					{
						CloudletType: CloudletTypeCD,
						CreatedBy:    "User1",
						CreatedDate:  test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
						CurrentActivations: CurrentActivations{
							Production: ActivationInfo{
								Effective: &PolicyActivation{
									CreatedBy:   "User1",
									CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
									FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
									ID:          123,
									Links: []Link{
										{
											Href: "Link1",
											Rel:  "self",
										},
									},
									Network:              ProductionNetwork,
									Operation:            OperationActivation,
									PolicyID:             1234,
									PolicyVersion:        1,
									PolicyVersionDeleted: false,
									Status:               ActivationStatusSuccess,
								},
								Latest: &PolicyActivation{
									CreatedBy:   "User1",
									CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
									FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
									ID:          321,
									Links: []Link{
										{
											Href: "Link2",
											Rel:  "self",
										},
									},
									Network:              ProductionNetwork,
									Operation:            OperationActivation,
									PolicyID:             4321,
									PolicyVersion:        1,
									PolicyVersionDeleted: false,
									Status:               ActivationStatusSuccess,
								},
							},
							Staging: ActivationInfo{
								Effective: &PolicyActivation{
									CreatedBy:   "User3",
									CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
									FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
									ID:          789,
									Links: []Link{
										{
											Href: "Link3",
											Rel:  "self",
										},
									},
									Network:              StagingNetwork,
									Operation:            OperationActivation,
									PolicyID:             6789,
									PolicyVersion:        1,
									PolicyVersionDeleted: false,
									Status:               ActivationStatusSuccess,
								},
								Latest: &PolicyActivation{
									CreatedBy:   "User3",
									CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
									FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
									ID:          987,
									Links: []Link{
										{
											Href: "Link4",
											Rel:  "self",
										},
									},
									Network:              StagingNetwork,
									Operation:            OperationActivation,
									PolicyID:             9876,
									PolicyVersion:        1,
									PolicyVersionDeleted: false,
									Status:               ActivationStatusSuccess,
								},
							},
						},
						Description: ptr.To("Test"),
						GroupID:     1,
						ID:          22,
						Links: []Link{
							{
								Href: "Link5",
								Rel:  "self",
							},
						},
						ModifiedBy:   "User1",
						ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
						Name:         "TestName",
						PolicyType:   PolicyTypeShared,
					},
				},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies?page=0&size=1000",
						Rel:  "self",
					},
				},
				Page: Page{
					Number:        0,
					Size:          1000,
					TotalElements: 54,
					TotalPages:    1,
				},
			},
		},
		"200 OK - with query params": {
			params: ListPoliciesRequest{
				Page: 2,
				Size: 12,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "content": [
        {
            "cloudletType": "CD",
            "createdBy": "User1",
            "createdDate": "2023-10-23T11:21:19.896Z",
            "currentActivations": {
                "production": {
                    "effective": null,
                    "latest": null
                },
                "staging": {
                    "effective": null,
                    "latest": null
                }
            },
            "description": null,
            "groupId": 1,
            "id": 11,
            "links": [
                {
                    "href": "Link1",
                    "rel": "self"
                }
            ],
            "modifiedBy": "User2",
            "modifiedDate": "2023-10-23T11:21:19.896Z",
            "name": "Name1",
            "policyType": "SHARED"
        }
	],
    "links": [
        {
            "href": "/cloudlets/v3/policies?page=0&size=1000",
            "rel": "self"
        }
    ],
    "page": {
        "number": 0,
        "size": 1000,
        "totalElements": 54,
        "totalPages": 1
    }
}
`,
			expectedPath: "/cloudlets/v3/policies?page=2&size=12",
			expectedResponse: &ListPoliciesResponse{
				Content: []Policy{
					{
						CloudletType:       CloudletTypeCD,
						CreatedBy:          "User1",
						CreatedDate:        test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
						CurrentActivations: CurrentActivations{},
						Description:        nil,
						GroupID:            1,
						ID:                 11,
						Links: []Link{
							{
								Href: "Link1",
								Rel:  "self",
							},
						},
						ModifiedBy:   "User2",
						ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
						Name:         "Name1",
						PolicyType:   PolicyTypeShared,
					},
				},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies?page=0&size=1000",
						Rel:  "self",
					},
				},
				Page: Page{
					Number:        0,
					Size:          1000,
					TotalElements: 54,
					TotalPages:    1,
				},
			},
		},
		"200 OK - empty content": {
			params:         ListPoliciesRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "content": [],
    "links": [
        {
            "href": "/cloudlets/v3/policies?page=0&size=1000",
            "rel": "self"
        }
    ],
    "page": {
        "number": 0,
        "size": 1000,
        "totalElements": 0,
        "totalPages": 1
    }
}
`,
			expectedPath: "/cloudlets/v3/policies",
			expectedResponse: &ListPoliciesResponse{
				Content: []Policy{},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies?page=0&size=1000",
						Rel:  "self",
					},
				},
				Page: Page{
					Number:        0,
					Size:          1000,
					TotalElements: 0,
					TotalPages:    1,
				},
			},
		},
		"validation errors - size lower than 10, negative page number": {
			params: ListPoliciesRequest{
				Page: -2,
				Size: 5,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "list shared policies: struct validation: Page: must be no less than 0\nSize: must be no less than 10", err.Error())
			},
		},
		"500 Internal Server Error": {
			params:         ListPoliciesRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
			"title": "Internal Server Error",
			"status": 500,
			"requestId": "1",
			"requestTime": "12:00",
			"clientIp": "1.1.1.1",
			"serverIp": "2.2.2.2",
			"method": "GET"
		}`,
			expectedPath: "/cloudlets/v3/policies",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "internal_error",
					Title:       "Internal Server Error",
					Status:      http.StatusInternalServerError,
					RequestID:   "1",
					RequestTime: "12:00",
					ClientIP:    "1.1.1.1",
					ServerIP:    "2.2.2.2",
					Method:      "GET",
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
			result, err := client.ListPolicies(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreatePolicy(t *testing.T) {
	tests := map[string]struct {
		params              CreatePolicyRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *Policy
		withError           func(*testing.T, error)
	}{
		"200 OK - minimal data": {
			params: CreatePolicyRequest{
				CloudletType: CloudletTypeFR,
				GroupID:      1,
				Name:         "TestName",
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "cloudletType": "FR",
    "createdBy": "User1",
    "createdDate": "2023-10-23T11:21:19.896Z",
    "currentActivations": {
        "production": {
            "effective": null,
            "latest": null
        },
        "staging": {
            "effective": null,
            "latest": null
        }
    },
    "description": null,
    "groupId": 1,
    "id": 11,
    "links": [
        {
            "href": "Link1",
            "rel": "self"
        }
    ],
    "modifiedBy": "User1",
    "modifiedDate": "2023-10-23T11:21:19.896Z",
    "name": "TestName",
    "policyType": "SHARED"
}
`,
			expectedPath: "/cloudlets/v3/policies",
			expectedRequestBody: `
{
  "cloudletType": "FR",
  "groupId": 1,
  "name": "TestName"
}`,
			expectedResponse: &Policy{
				CloudletType:       CloudletTypeFR,
				CreatedBy:          "User1",
				CreatedDate:        test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{},
				Description:        nil,
				GroupID:            1,
				ID:                 11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "TestName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"200 OK - all data": {
			params: CreatePolicyRequest{
				CloudletType: CloudletTypeFR,
				Description:  ptr.To("Description"),
				GroupID:      1,
				Name:         "TestName",
				PolicyType:   PolicyTypeShared,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "cloudletType": "FR",
    "createdBy": "User1",
    "createdDate": "2023-10-23T11:21:19.896Z",
    "currentActivations": {
        "production": {
            "effective": null,
            "latest": null
        },
        "staging": {
            "effective": null,
            "latest": null
        }
    },
    "description": "Description",
    "groupId": 1,
    "id": 11,
    "links": [
        {
            "href": "Link1",
            "rel": "self"
        }
    ],
    "modifiedBy": "User1",
    "modifiedDate": "2023-10-23T11:21:19.896Z",
    "name": "TestName",
    "policyType": "SHARED"
}
`,
			expectedPath: "/cloudlets/v3/policies",
			expectedRequestBody: `
{
  "cloudletType": "FR",
  "description": "Description",
  "groupId": 1,
  "name": "TestName",
  "policyType": "SHARED"
}
`,
			expectedResponse: &Policy{
				CloudletType:       CloudletTypeFR,
				CreatedBy:          "User1",
				CreatedDate:        test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{},
				Description:        ptr.To("Description"),
				GroupID:            1,
				ID:                 11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "TestName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"validation errors": {
			params: CreatePolicyRequest{
				CloudletType: "Wrong Cloudlet Type",
				Description:  ptr.To(strings.Repeat("Too long description", 20)),
				GroupID:      1,
				Name:         "TestName not match",
				PolicyType:   "Wrong Policy Type",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create shared policy: struct validation: CloudletType: value 'Wrong Cloudlet Type' is invalid. Must be one of: 'AP', 'AS', 'CD', 'ER', 'FR', 'IG'\nDescription: the length must be no more than 255\nName: value 'TestName not match' is invalid. Must be of format: ^[a-z_A-Z0-9]+$\nPolicyType: value 'Wrong Policy Type' is invalid. Must be 'SHARED'", err.Error())
			},
		},
		"validation errors - missing required params": {
			params: CreatePolicyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create shared policy: struct validation: CloudletType: cannot be blank\nGroupID: cannot be blank\nName: cannot be blank", err.Error())
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreatePolicy(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeletePolicy(t *testing.T) {
	tests := map[string]struct {
		params         DeletePolicyRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204": {
			params: DeletePolicyRequest{
				PolicyID: 1,
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/cloudlets/v3/policies/1",
		},
		"validation errors - missing required param": {
			params: DeletePolicyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete shared policy: struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeletePolicy(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetPolicy(t *testing.T) {
	tests := map[string]struct {
		params           GetPolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Policy
		withError        func(*testing.T, error)
	}{
		"200 OK - minimal data": {
			params: GetPolicyRequest{
				PolicyID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `

	{
	    "cloudletType": "FR",
	    "createdBy": "User1",
	    "createdDate": "2023-10-23T11:21:19.896Z",
	    "currentActivations": {
	        "production": {
	            "effective": null,
	            "latest": null
	        },
	        "staging": {
	            "effective": null,
	            "latest": null
	        }
	    },
	    "description": null,
	    "groupId": 1,
	    "id": 11,
	    "links": [
	        {
	            "href": "Link1",
	            "rel": "self"
	        }
	    ],
	    "modifiedBy": "User1",
	    "modifiedDate": "2023-10-23T11:21:19.896Z",
	    "name": "TestName",
	    "policyType": "SHARED"
	}

`,

			expectedPath: "/cloudlets/v3/policies/1",
			expectedResponse: &Policy{
				CloudletType:       CloudletTypeFR,
				CreatedBy:          "User1",
				CreatedDate:        test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{},
				Description:        nil,
				GroupID:            1,
				ID:                 11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "TestName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"200 OK - with activation information": {
			params: GetPolicyRequest{
				PolicyID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `

	{
	    "cloudletType": "FR",
	    "createdBy": "User1",
	    "createdDate": "2023-10-23T11:21:19.896Z",
	    "currentActivations": {
			"production": {
				"effective": {
					"createdBy": "User1",
					"createdDate": "2023-10-23T11:21:19.896Z",
					"finishDate": "2023-10-23T11:22:57.589Z",
					"id": 123,
					"links": [
						{
							"href": "Link1",
							"rel": "self"
						}
					],
					"network": "PRODUCTION",
					"operation": "ACTIVATION",
					"policyId": 1234,
					"policyVersion": 1,
					"policyVersionDeleted": false,
					"status": "SUCCESS"
				},
				"latest": {
					"createdBy": "User1",
					"createdDate": "2023-10-23T11:21:19.896Z",
					"finishDate": "2023-10-23T11:22:57.589Z",
					"id": 321,
					"links": [
						{
							"href": "Link2",
							"rel": "self"
						}
					],
					"network": "PRODUCTION",
					"operation": "ACTIVATION",
					"policyId": 4321,
					"policyVersion": 1,
					"policyVersionDeleted": false,
					"status": "SUCCESS"
				}
			},
			"staging": {
				"effective": {
					"createdBy": "User3",
					"createdDate": "2023-10-23T11:21:19.896Z",
					"finishDate": "2023-10-23T11:22:57.589Z",
					"id": 789,
					"links": [
						{
							"href": "Link3",
							"rel": "self"
						}
					],
					"network": "STAGING",
					"operation": "ACTIVATION",
					"policyId": 6789,
					"policyVersion": 1,
					"policyVersionDeleted": false,
					"status": "SUCCESS"
				},
				"latest": {
					"createdBy": "User3",
					"createdDate": "2023-10-23T11:21:19.896Z",
					"finishDate": "2023-10-23T11:22:57.589Z",
					"id": 987,
					"links": [
						{
							"href": "Link4",
							"rel": "self"
						}
					],
					"network": "STAGING",
					"operation": "ACTIVATION",
					"policyId": 9876,
					"policyVersion": 1,
					"policyVersionDeleted": false,
					"status": "SUCCESS"
				}
			}
		},
	    "description": "Description",
	    "groupId": 1,
	    "id": 11,
	    "links": [
	        {
	            "href": "Link1",
	            "rel": "self"
	        }
	    ],
	    "modifiedBy": "User1",
	    "modifiedDate": "2023-10-23T11:21:19.896Z",
	    "name": "TestName",
	    "policyType": "SHARED"
	}

`,

			expectedPath: "/cloudlets/v3/policies/1",
			expectedResponse: &Policy{
				CloudletType: CloudletTypeFR,
				CreatedBy:    "User1",
				CreatedDate:  test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{
					Production: ActivationInfo{
						Effective: &PolicyActivation{
							CreatedBy:   "User1",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          123,
							Links: []Link{
								{
									Href: "Link1",
									Rel:  "self",
								},
							},
							Network:              ProductionNetwork,
							Operation:            OperationActivation,
							PolicyID:             1234,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
						Latest: &PolicyActivation{
							CreatedBy:   "User1",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          321,
							Links: []Link{
								{
									Href: "Link2",
									Rel:  "self",
								},
							},
							Network:              ProductionNetwork,
							Operation:            OperationActivation,
							PolicyID:             4321,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
					},
					Staging: ActivationInfo{
						Effective: &PolicyActivation{
							CreatedBy:   "User3",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          789,
							Links: []Link{
								{
									Href: "Link3",
									Rel:  "self",
								},
							},
							Network:              StagingNetwork,
							Operation:            OperationActivation,
							PolicyID:             6789,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
						Latest: &PolicyActivation{
							CreatedBy:   "User3",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          987,
							Links: []Link{
								{
									Href: "Link4",
									Rel:  "self",
								},
							},
							Network:              StagingNetwork,
							Operation:            OperationActivation,
							PolicyID:             9876,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
					},
				},
				Description: ptr.To("Description"),
				GroupID:     1,
				ID:          11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "TestName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"200 OK - one network is active": {
			params: GetPolicyRequest{
				PolicyID: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `

	{
	    "cloudletType": "FR",
	    "createdBy": "User1",
	    "createdDate": "2023-10-23T11:21:19.896Z",
	    "currentActivations": {
			"production": {
				"effective": null,
				"latest": null
			},
			"staging": {
				"effective": {
					"createdBy": "User3",
					"createdDate": "2023-10-23T11:21:19.896Z",
					"finishDate": "2023-10-23T11:22:57.589Z",
					"id": 789,
					"links": [
						{
							"href": "Link3",
							"rel": "self"
						}
					],
					"network": "STAGING",
					"operation": "ACTIVATION",
					"policyId": 6789,
					"policyVersion": 1,
					"policyVersionDeleted": false,
					"status": "SUCCESS"
				},
				"latest": {
					"createdBy": "User3",
					"createdDate": "2023-10-23T11:21:19.896Z",
					"finishDate": "2023-10-23T11:22:57.589Z",
					"id": 987,
					"links": [
						{
							"href": "Link4",
							"rel": "self"
						}
					],
					"network": "STAGING",
					"operation": "ACTIVATION",
					"policyId": 9876,
					"policyVersion": 1,
					"policyVersionDeleted": false,
					"status": "SUCCESS"
				}
			}
		},
	    "description": "Description",
	    "groupId": 1,
	    "id": 11,
	    "links": [
	        {
	            "href": "Link1",
	            "rel": "self"
	        }
	    ],
	    "modifiedBy": "User1",
	    "modifiedDate": "2023-10-23T11:21:19.896Z",
	    "name": "TestName",
	    "policyType": "SHARED"
	}

`,

			expectedPath: "/cloudlets/v3/policies/1",
			expectedResponse: &Policy{
				CloudletType: CloudletTypeFR,
				CreatedBy:    "User1",
				CreatedDate:  test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{
					Production: ActivationInfo{
						Effective: nil,
						Latest:    nil,
					},
					Staging: ActivationInfo{
						Effective: &PolicyActivation{
							CreatedBy:   "User3",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          789,
							Links: []Link{
								{
									Href: "Link3",
									Rel:  "self",
								},
							},
							Network:              StagingNetwork,
							Operation:            OperationActivation,
							PolicyID:             6789,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
						Latest: &PolicyActivation{
							CreatedBy:   "User3",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          987,
							Links: []Link{
								{
									Href: "Link4",
									Rel:  "self",
								},
							},
							Network:              StagingNetwork,
							Operation:            OperationActivation,
							PolicyID:             9876,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
					},
				},
				Description: ptr.To("Description"),
				GroupID:     1,
				ID:          11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "TestName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"validation errors - missing required params": {
			params: GetPolicyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get shared policy: struct validation: PolicyID: cannot be blank", err.Error())
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
			result, err := client.GetPolicy(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdatePolicy(t *testing.T) {
	tests := map[string]struct {
		params              UpdatePolicyRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *Policy
		withError           func(*testing.T, error)
	}{
		"200 OK - minimal data": {
			params: UpdatePolicyRequest{
				PolicyID: 1,
				Body: UpdatePolicyRequestBody{
					GroupID: 11,
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cloudletType": "FR",
    "createdBy": "User1",
    "createdDate": "2023-10-23T11:21:19.896Z",
    "currentActivations": {
        "production": {
            "effective": null,
            "latest": null
        },
        "staging": {
            "effective": null,
            "latest": null
        }
    },
    "description": null,
    "groupId": 1,
    "id": 11,
    "links": [
        {
            "href": "Link1",
            "rel": "self"
        }
    ],
    "modifiedBy": "User1",
    "modifiedDate": "2023-10-23T11:21:19.896Z",
    "name": "TestName",
    "policyType": "SHARED"
}
`,
			expectedPath: "/cloudlets/v3/policies/1",
			expectedRequestBody: `
{
  "groupId": 11
}
`,
			expectedResponse: &Policy{
				CloudletType:       CloudletTypeFR,
				CreatedBy:          "User1",
				CreatedDate:        test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{},
				Description:        nil,
				GroupID:            1,
				ID:                 11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "TestName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"200 OK - with description and activations": {
			params: UpdatePolicyRequest{
				PolicyID: 1,
				Body: UpdatePolicyRequestBody{
					GroupID:     11,
					Description: ptr.To("Description"),
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cloudletType": "FR",
    "createdBy": "User1",
    "createdDate": "2023-10-23T11:21:19.896Z",
    "currentActivations": {
		"production": {
			"effective": {
				"createdBy": "User1",
				"createdDate": "2023-10-23T11:21:19.896Z",
				"finishDate": "2023-10-23T11:22:57.589Z",
				"id": 123,
				"links": [
					{
						"href": "Link1",
						"rel": "self"
					}
				],
				"network": "PRODUCTION",
				"operation": "ACTIVATION",
				"policyId": 1234,
				"policyVersion": 1,
				"policyVersionDeleted": false,
				"status": "SUCCESS"
			},
			"latest": {
				"createdBy": "User1",
				"createdDate": "2023-10-23T11:21:19.896Z",
				"finishDate": "2023-10-23T11:22:57.589Z",
				"id": 321,
				"links": [
					{
						"href": "Link2",
						"rel": "self"
					}
				],
				"network": "PRODUCTION",
				"operation": "ACTIVATION",
				"policyId": 4321,
				"policyVersion": 1,
				"policyVersionDeleted": false,
				"status": "SUCCESS"
			}
		},
		"staging": {
			"effective": {
				"createdBy": "User3",
				"createdDate": "2023-10-23T11:21:19.896Z",
				"finishDate": "2023-10-23T11:22:57.589Z",
				"id": 789,
				"links": [
					{
						"href": "Link3",
						"rel": "self"
					}
				],
				"network": "STAGING",
				"operation": "ACTIVATION",
				"policyId": 6789,
				"policyVersion": 1,
				"policyVersionDeleted": false,
				"status": "SUCCESS"
			},
			"latest": {
				"createdBy": "User3",
				"createdDate": "2023-10-23T11:21:19.896Z",
				"finishDate": "2023-10-23T11:22:57.589Z",
				"id": 987,
				"links": [
					{
						"href": "Link4",
						"rel": "self"
					}
				],
				"network": "STAGING",
				"operation": "ACTIVATION",
				"policyId": 9876,
				"policyVersion": 1,
				"policyVersionDeleted": false,
				"status": "SUCCESS"
			}
		}
	},
    "description": "Description",
    "groupId": 1,
    "id": 11,
    "links": [
        {
            "href": "Link1",
            "rel": "self"
        }
    ],
    "modifiedBy": "User1",
    "modifiedDate": "2023-10-23T11:21:19.896Z",
    "name": "TestName",
    "policyType": "SHARED"
}
`,
			expectedPath: "/cloudlets/v3/policies/1",
			expectedRequestBody: `
{
  "groupId": 11,
  "description": "Description"
}
`,
			expectedResponse: &Policy{
				CloudletType: CloudletTypeFR,
				CreatedBy:    "User1",
				CreatedDate:  test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{
					Production: ActivationInfo{
						Effective: &PolicyActivation{
							CreatedBy:   "User1",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          123,
							Links: []Link{
								{
									Href: "Link1",
									Rel:  "self",
								},
							},
							Network:              ProductionNetwork,
							Operation:            OperationActivation,
							PolicyID:             1234,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
						Latest: &PolicyActivation{
							CreatedBy:   "User1",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          321,
							Links: []Link{
								{
									Href: "Link2",
									Rel:  "self",
								},
							},
							Network:              ProductionNetwork,
							Operation:            OperationActivation,
							PolicyID:             4321,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
					},
					Staging: ActivationInfo{
						Effective: &PolicyActivation{
							CreatedBy:   "User3",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          789,
							Links: []Link{
								{
									Href: "Link3",
									Rel:  "self",
								},
							},
							Network:              StagingNetwork,
							Operation:            OperationActivation,
							PolicyID:             6789,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
						Latest: &PolicyActivation{
							CreatedBy:   "User3",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          987,
							Links: []Link{
								{
									Href: "Link4",
									Rel:  "self",
								},
							},
							Network:              StagingNetwork,
							Operation:            OperationActivation,
							PolicyID:             9876,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
					},
				},
				Description: ptr.To("Description"),
				GroupID:     1,
				ID:          11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "TestName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"validation errors - missing required params": {
			params: UpdatePolicyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update shared policy: struct validation: GroupID: cannot be blank\nPolicyID: cannot be blank", err.Error())
			},
		},
		"validation errors - description too long": {
			params: UpdatePolicyRequest{
				PolicyID: 1,
				Body: UpdatePolicyRequestBody{
					GroupID:     11,
					Description: ptr.To(strings.Repeat("TestDescription", 30)),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update shared policy: struct validation: Description: the length must be no more than 255", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdatePolicy(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestClonePolicy(t *testing.T) {
	tests := map[string]struct {
		params              ClonePolicyRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *Policy
		withError           func(*testing.T, error)
	}{
		"200 OK - minimal data": {
			params: ClonePolicyRequest{
				PolicyID: 1,
				Body: ClonePolicyRequestBody{
					NewName: "NewName",
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cloudletType": "FR",
    "createdBy": "User1",
    "createdDate": "2023-10-23T11:21:19.896Z",
    "currentActivations": {
        "production": {
            "effective": null,
            "latest": null
        },
        "staging": {
            "effective": null,
            "latest": null
        }
    },
    "description": null,
    "groupId": 1,
    "id": 11,
    "links": [
        {
            "href": "Link1",
            "rel": "self"
        }
    ],
    "modifiedBy": "User1",
    "modifiedDate": "2023-10-23T11:21:19.896Z",
    "name": "NewName",
    "policyType": "SHARED"
}
`,
			expectedPath: "/cloudlets/v3/policies/1/clone",
			expectedRequestBody: `
{
  "newName": "NewName"
}
`,
			expectedResponse: &Policy{
				CloudletType:       "FR",
				CreatedBy:          "User1",
				CreatedDate:        test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{},
				Description:        nil,
				GroupID:            1,
				ID:                 11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "NewName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"200 OK - all data": {
			params: ClonePolicyRequest{
				PolicyID: 1,
				Body: ClonePolicyRequestBody{
					AdditionalVersions: []int64{1, 2},
					GroupID:            11,
					NewName:            "NewName",
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cloudletType": "FR",
    "createdBy": "User1",
    "createdDate": "2023-10-23T11:21:19.896Z",
    "currentActivations": {
		"production": {
			"effective": {
				"createdBy": "User1",
				"createdDate": "2023-10-23T11:21:19.896Z",
				"finishDate": "2023-10-23T11:22:57.589Z",
				"id": 123,
				"links": [
					{
						"href": "Link1",
						"rel": "self"
					}
				],
				"network": "PRODUCTION",
				"operation": "ACTIVATION",
				"policyId": 1234,
				"policyVersion": 1,
				"policyVersionDeleted": false,
				"status": "SUCCESS"
			},
			"latest": {
				"createdBy": "User1",
				"createdDate": "2023-10-23T11:21:19.896Z",
				"finishDate": "2023-10-23T11:22:57.589Z",
				"id": 321,
				"links": [
					{
						"href": "Link2",
						"rel": "self"
					}
				],
				"network": "PRODUCTION",
				"operation": "ACTIVATION",
				"policyId": 4321,
				"policyVersion": 1,
				"policyVersionDeleted": false,
				"status": "SUCCESS"
			}
		},
		"staging": {
			"effective": {
				"createdBy": "User3",
				"createdDate": "2023-10-23T11:21:19.896Z",
				"finishDate": "2023-10-23T11:22:57.589Z",
				"id": 789,
				"links": [
					{
						"href": "Link3",
						"rel": "self"
					}
				],
				"network": "STAGING",
				"operation": "ACTIVATION",
				"policyId": 6789,
				"policyVersion": 1,
				"policyVersionDeleted": false,
				"status": "SUCCESS"
			},
			"latest": {
				"createdBy": "User3",
				"createdDate": "2023-10-23T11:21:19.896Z",
				"finishDate": "2023-10-23T11:22:57.589Z",
				"id": 987,
				"links": [
					{
						"href": "Link4",
						"rel": "self"
					}
				],
				"network": "STAGING",
				"operation": "ACTIVATION",
				"policyId": 9876,
				"policyVersion": 1,
				"policyVersionDeleted": false,
				"status": "SUCCESS"
			}
		}
	},
    "description": "Description",
    "groupId": 1,
    "id": 11,
    "links": [
        {
            "href": "Link1",
            "rel": "self"
        }
    ],
    "modifiedBy": "User1",
    "modifiedDate": "2023-10-23T11:21:19.896Z",
    "name": "NewName",
    "policyType": "SHARED"
}
`,
			expectedPath: "/cloudlets/v3/policies/1/clone",
			expectedRequestBody: `
{
   "additionalVersions": [1, 2],
   "groupId": 11,
   "newName": "NewName"
}
`,
			expectedResponse: &Policy{
				CloudletType: "FR",
				CreatedBy:    "User1",
				CreatedDate:  test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
				CurrentActivations: CurrentActivations{
					Production: ActivationInfo{
						Effective: &PolicyActivation{
							CreatedBy:   "User1",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          123,
							Links: []Link{
								{
									Href: "Link1",
									Rel:  "self",
								},
							},
							Network:              ProductionNetwork,
							Operation:            OperationActivation,
							PolicyID:             1234,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
						Latest: &PolicyActivation{
							CreatedBy:   "User1",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          321,
							Links: []Link{
								{
									Href: "Link2",
									Rel:  "self",
								},
							},
							Network:              ProductionNetwork,
							Operation:            OperationActivation,
							PolicyID:             4321,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
					},
					Staging: ActivationInfo{
						Effective: &PolicyActivation{
							CreatedBy:   "User3",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          789,
							Links: []Link{
								{
									Href: "Link3",
									Rel:  "self",
								},
							},
							Network:              StagingNetwork,
							Operation:            OperationActivation,
							PolicyID:             6789,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
						Latest: &PolicyActivation{
							CreatedBy:   "User3",
							CreatedDate: test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z"),
							FinishDate:  ptr.To(test.NewTimeFromString(t, "2023-10-23T11:22:57.589Z")),
							ID:          987,
							Links: []Link{
								{
									Href: "Link4",
									Rel:  "self",
								},
							},
							Network:              StagingNetwork,
							Operation:            OperationActivation,
							PolicyID:             9876,
							PolicyVersion:        1,
							PolicyVersionDeleted: false,
							Status:               ActivationStatusSuccess,
						},
					},
				},
				Description: ptr.To("Description"),
				GroupID:     1,
				ID:          11,
				Links: []Link{
					{
						Href: "Link1",
						Rel:  "self",
					},
				},
				ModifiedBy:   "User1",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-23T11:21:19.896Z")),
				Name:         "NewName",
				PolicyType:   PolicyTypeShared,
			},
		},
		"validation errors - missing required params": {
			params: ClonePolicyRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "clone policy: struct validation: NewName: cannot be blank\nPolicyID: cannot be blank", err.Error())
			},
		},
		"validation errors - newName too long": {
			params: ClonePolicyRequest{
				PolicyID: 1,
				Body: ClonePolicyRequestBody{
					GroupID: 11,
					NewName: strings.Repeat("TestNameTooLong", 10),
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "clone policy: struct validation: NewName: the length must be no more than 64", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)

				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ClonePolicy(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
