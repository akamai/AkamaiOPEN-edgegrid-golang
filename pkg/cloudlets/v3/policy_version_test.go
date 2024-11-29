package v3

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/internal/test"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPolicyVersions(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          ListPolicyVersionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListPolicyVersions
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: ListPolicyVersionsRequest{
				PolicyID: 670790,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"content": [
		{
			"createdBy": "jsmith",
			"createdDate": "2023-10-19T08:50:47.350Z",
			"description": null,
			"id": 6551184,
			"immutable": false,
			"links": [],
			"modifiedBy": "jsmith",
			"modifiedDate": "2023-10-19T08:50:47.350Z",
			"policyId": 670790,
			"version": 3
		},
		{
			"createdBy": "jsmith",
			"createdDate": "2023-10-19T08:50:46.225Z",
			"description": null,
			"id": 6551183,
			"immutable": false,
			"links": [],
			"modifiedBy": "jsmith",
			"modifiedDate": "2023-10-19T08:50:46.225Z",
			"policyId": 670790,
			"version": 2
		},
		{
			"createdBy": "jsmith",
			"createdDate": "2023-10-19T08:44:34.398Z",
			"description": null,
			"id": 6551182,
			"immutable": false,
			"links": [],
			"modifiedBy": "jsmith",
			"modifiedDate": "2023-10-19T08:44:34.398Z",
			"policyId": 670790,
			"version": 1
		}
	],
	"links": [
		{
			"href": "/cloudlets/v3/policies/670790/versions?page=0&size=1000",
			"rel": "self"
		}
	],
	"page": {
		"number": 0,
		"size": 1000,
		"totalElements": 3,
		"totalPages": 1
	}
}`,
			expectedPath: "/cloudlets/v3/policies/670790/versions?page=0",
			expectedResponse: &ListPolicyVersions{
				PolicyVersions: []ListPolicyVersionsItem{
					{
						CreatedBy:     "jsmith",
						CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
						Description:   nil,
						ID:            6551184,
						Immutable:     false,
						Links:         []Link{},
						ModifiedBy:    "jsmith",
						ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
						PolicyID:      670790,
						PolicyVersion: 3,
					},
					{
						CreatedBy:     "jsmith",
						CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:46.225Z"),
						Description:   nil,
						ID:            6551183,
						Immutable:     false,
						Links:         []Link{},
						ModifiedBy:    "jsmith",
						ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:46.225Z")),
						PolicyID:      670790,
						PolicyVersion: 2,
					},
					{
						CreatedBy:     "jsmith",
						CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:44:34.398Z"),
						Description:   nil,
						ID:            6551182,
						Immutable:     false,
						Links:         []Link{},
						ModifiedBy:    "jsmith",
						ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:44:34.398Z")),
						PolicyID:      670790,
						PolicyVersion: 1,
					},
				},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/670790/versions?page=0&size=1000",
						Rel:  "self",
					},
				},
				Page: Page{
					Number:        0,
					Size:          1000,
					TotalElements: 3,
					TotalPages:    1,
				},
			},
		},
		"200 OK with params": {
			request: ListPolicyVersionsRequest{
				PolicyID: 670790,
				Page:     4,
				Size:     10,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"content": [],
	"links": [
		{
			"href": "/cloudlets/v3/policies/670790/versions?page=0&size=10",
			"rel": "first"
		},
		{
			"href": "/cloudlets/v3/policies/670790/versions?page=3&size=10",
			"rel": "prev"
		},
		{
			"href": "/cloudlets/v3/policies/670790/versions?page=4&size=10",
			"rel": "self"
		},
		{
			"href": "/cloudlets/v3/policies/670790/versions?page=0&size=10",
			"rel": "last"
		}
	],
	"page": {
		"number": 4,
		"size": 10,
		"totalElements": 3,
		"totalPages": 1
	}
}`,
			expectedPath: "/cloudlets/v3/policies/670790/versions?page=4&size=10",
			expectedResponse: &ListPolicyVersions{
				PolicyVersions: []ListPolicyVersionsItem{},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/670790/versions?page=0&size=10",
						Rel:  "first",
					},
					{
						Href: "/cloudlets/v3/policies/670790/versions?page=3&size=10",
						Rel:  "prev",
					},
					{
						Href: "/cloudlets/v3/policies/670790/versions?page=4&size=10",
						Rel:  "self",
					},
					{
						Href: "/cloudlets/v3/policies/670790/versions?page=0&size=10",
						Rel:  "last",
					},
				},
				Page: Page{
					Number:        4,
					Size:          10,
					TotalElements: 3,
					TotalPages:    1,
				},
			},
		},
		"500 internal server error": {
			request: ListPolicyVersionsRequest{
				PolicyID: 284823,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
		   "title": "Internal Server Error",
		   "status": 500
		}`,
			expectedPath: "/cloudlets/v3/policies/284823/versions?page=0",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Status: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListPolicyVersions(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetPolicyVersion(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          GetPolicyVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PolicyVersion
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: GetPolicyVersionRequest{
				PolicyID:      670798,
				PolicyVersion: 2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
		{
			"createdBy": "jsmith",
			"createdDate": "2023-10-19T09:46:57.395Z",
			"description": "test description",
			"id": 6551191,
			"immutable": false,
			"matchRules": [],
			"matchRulesWarnings": [],
			"modifiedBy": "jsmith",
			"modifiedDate": "2023-10-19T09:46:57.395Z",
			"policyId": 670798,
			"version": 2
		}`,
			expectedPath: "/cloudlets/v3/policies/670798/versions/2",
			expectedResponse: &PolicyVersion{
				CreatedBy:          "jsmith",
				CreatedDate:        test.NewTimeFromString(t, "2023-10-19T09:46:57.395Z"),
				Description:        ptr.To("test description"),
				ID:                 6551191,
				MatchRules:         nil,
				MatchRulesWarnings: []MatchRulesWarning{},
				ModifiedBy:         "jsmith",
				ModifiedDate:       ptr.To(test.NewTimeFromString(t, "2023-10-19T09:46:57.395Z")),
				PolicyID:           670798,
				PolicyVersion:      2,
			},
		},
		"200 OK, ER with disabled rule": {
			request: GetPolicyVersionRequest{
				PolicyID:      276858,
				PolicyVersion: 1,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
	"createdBy": "jsmith",
	"createdDate": "2023-10-19T10:45:30.619Z",
	"description": null,
	"id": 6551242,
	"immutable": false,
	"matchRules": [
		{
			"type": "erMatchRule",
			"id": 0,
			"name": "er_rule",
			"start": 0,
			"end": 0,
			"matchURL": null,
			"disabled": true,
			"matches": [],
			"akaRuleId": "6d3bbc891fc0d8ce",
			"statusCode": 301,
			"redirectURL": "/path",
			"useIncomingQueryString": false
		}
	],
	"matchRulesWarnings": [],
	"modifiedBy": "jsmith",
	"modifiedDate": "2023-10-19T10:45:30.619Z",
	"policyId": 276858,
	"version": 1
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions/1",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T10:45:30.619Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				PolicyID:      276858,
				PolicyVersion: 1,
				ID:            6551242,
				MatchRules: MatchRules{
					&MatchRuleER{
						Type:                   "erMatchRule",
						End:                    0,
						ID:                     0,
						MatchURL:               "",
						Name:                   "er_rule",
						Matches:                []MatchCriteriaER{},
						RedirectURL:            "/path",
						Start:                  0,
						StatusCode:             301,
						UseIncomingQueryString: false,
						Disabled:               true,
					},
				},
				MatchRulesWarnings: []MatchRulesWarning{},
				ModifiedBy:         "jsmith",
				ModifiedDate:       ptr.To(test.NewTimeFromString(t, "2023-10-19T10:45:30.619Z")),
			},
		},
		"200 OK, AS rule with disabled=false": {
			request: GetPolicyVersionRequest{
				PolicyID:      355557,
				PolicyVersion: 2,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
		    "createdDate": "2023-10-19T09:46:57.395Z",
		    "createdBy": "jsmith",
		    "description": "Initial version",
		    "modifiedBy": "jsmith",
		    "modifiedDate": "2023-10-19T10:45:30.619Z",
		    "matchRules": [
		        {
		            "type": "asMatchRule",
		            "akaRuleId": "f58014ee0cc17ce",
		            "end": 0,
		            "forwardSettings": {
		                "originId": "originremote2",
		                "pathAndQS": "/sales/Q1/",
		                "useIncomingQueryString": true
		            },
		            "id": 0,
		            "matches": [
		                {
		                    "caseSensitive": false,
		                    "matchOperator": "equals",
		                    "matchType": "range",
		                    "negate": false,
		                    "objectMatchValue": {
		                        "type": "range",
		                        "value": [
		                            1,
		                            25
		                        ]
		                    }
		                }
		            ],
		            "name": "Q1Sales",
		            "start": 0
		        }
		    ],
		    "policyId": 355557,
		    "version": 2
		}`,
			expectedPath: "/cloudlets/v3/policies/355557/versions/2",
			expectedResponse: &PolicyVersion{
				CreatedBy:   "jsmith",
				CreatedDate: test.NewTimeFromString(t, "2023-10-19T09:46:57.395Z"),
				Description: ptr.To("Initial version"),
				MatchRules: MatchRules{
					&MatchRuleAS{
						Type: "asMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsAS{
							OriginID:               "originremote2",
							PathAndQS:              "/sales/Q1/",
							UseIncomingQueryString: true,
						},
						ID: 0,
						Matches: []MatchCriteriaAS{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "range",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueRange{
									Type: "range",
									Value: []int64{
										1, 25},
								},
							},
						},
						Name:  "Q1Sales",
						Start: 0,
					},
				},
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T10:45:30.619Z")),
				PolicyID:      355557,
				PolicyVersion: 2,
			},
		},
		"200 OK, PR rule": {
			request: GetPolicyVersionRequest{
				PolicyID:      276858,
				PolicyVersion: 6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
		    "createdDate": "2023-10-19T09:46:57.395Z",
		    "createdBy": "jsmith",
		    "description": null,
		    "modifiedBy": "jsmith",
		    "modifiedDate": "2023-10-19T09:46:57.395Z",
		    "matchRules": [
		        {
		            "type": "cdMatchRule",
		            "akaRuleId": "b151ca68e51f5a61",
		            "end": 0,
		            "forwardSettings": {
		                "originId": "fr_test_krk_dc2",
		                "percent": 11
		            },
		            "id": 0,
		            "matches": [
		                {
		                    "caseSensitive": false,
		                    "matchOperator": "equals",
		                    "matchType": "method",
		                    "negate": false,
		                    "objectMatchValue": {
		                        "type": "simple",
		                        "value": [
		                            "GET"
		                        ]
		                    }
		                }
		            ],
		            "name": "rule 1",
		            "start": 0
		        }
		    ],
		    "policyId": 325401,
		    "version": 3
		}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions/6",
			expectedResponse: &PolicyVersion{
				CreatedBy:    "jsmith",
				CreatedDate:  test.NewTimeFromString(t, "2023-10-19T09:46:57.395Z"),
				Description:  nil,
				ModifiedBy:   "jsmith",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-19T09:46:57.395Z")),
				MatchRules: MatchRules{
					&MatchRulePR{
						Type: "cdMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsPR{
							OriginID: "fr_test_krk_dc2",
							Percent:  11,
						},
						ID: 0,
						Matches: []MatchCriteriaPR{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "method",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueSimple{
									Type: "simple",
									Value: []string{
										"GET"},
								},
							},
						},
						Name:  "rule 1",
						Start: 0,
					},
				},
				PolicyID:      325401,
				PolicyVersion: 3,
			},
		},
		"200 OK, ER rule with disabled=false": {
			request: GetPolicyVersionRequest{
				PolicyID:      276858,
				PolicyVersion: 6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
		    "createdDate": "2023-10-19T09:46:57.395Z",
		    "createdBy": "jsmith",
		    "description": null,
		    "modifiedBy": "jsmith",
		    "modifiedDate": "2023-10-19T09:46:57.395Z",
		    "matchRules": [
		        {
		            "type": "erMatchRule",
		            "end": 0,
		            "id": 0,
		            "matchURL": null,
		            "name": "rul3",
		            "redirectURL": "/abc/sss",
		            "start": 0,
		            "statusCode": 307,
		            "useIncomingQueryString": false,
		            "useIncomingSchemeAndHost": true,
		            "useRelativeUrl": "copy_scheme_hostname"
		        }
		    ],
		    "policyId": 276858,
		    "version": 6
		}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions/6",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T09:46:57.395Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T09:46:57.395Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleER{
						Type:                     "erMatchRule",
						End:                      0,
						ID:                       0,
						MatchURL:                 "",
						Name:                     "rul3",
						RedirectURL:              "/abc/sss",
						Start:                    0,
						StatusCode:               307,
						UseIncomingQueryString:   false,
						UseIncomingSchemeAndHost: true,
						UseRelativeURL:           "copy_scheme_hostname",
						Disabled:                 false,
					},
				},
			},
		},
		"200 OK, FR with disabled rule": {
			request: GetPolicyVersionRequest{
				PolicyID:      276858,
				PolicyVersion: 6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
		    "createdDate": "2023-10-19T08:50:47.350Z",
		    "createdBy": "jsmith",
		    "description": null,
		    "modifiedBy": "jsmith",
		    "modifiedDate": "2023-10-19T08:50:47.350Z",
		    "matchRules": [{
		        "type": "frMatchRule",
		        "disabled": true,
		        "end": 0,
		        "id": 0,
		        "matchURL": null,
		        "forwardSettings": {
		            "pathAndQS": "/test_images/simpleimg.jpg",
		            "useIncomingQueryString": true,
					"originId": "1234"
				},
		        "name": "rule 1",
		        "start": 0
		    }],
		    "policyId": 276858,
		    "version": 6
		}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions/6",
			expectedResponse: &PolicyVersion{
				CreatedBy:     "jsmith",
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				Description:   nil,
				PolicyID:      276858,
				PolicyVersion: 6,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				MatchRules: MatchRules{
					&MatchRuleFR{
						Name:     "rule 1",
						Type:     "frMatchRule",
						Start:    0,
						End:      0,
						ID:       0,
						MatchURL: "",
						ForwardSettings: ForwardSettingsFR{
							PathAndQS:              "/test_images/simpleimg.jpg",
							UseIncomingQueryString: true,
							OriginID:               "1234",
						},
						Disabled: true,
					},
				},
			},
		},
		"500 internal server error": {
			request: GetPolicyVersionRequest{
				PolicyID: 1,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
		   "title": "Internal Server Error",
		   "status": 500
		}`,
			expectedPath: "/cloudlets/v3/policies/1/versions/0",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:   "internal_error",
					Title:  "Internal Server Error",
					Status: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetPolicyVersion(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreatePolicyVersion(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          CreatePolicyVersionRequest
		requestBody      string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PolicyVersion
		withError        error
	}{
		"201 created, simple ER": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					Description: ptr.To("Description for the policy"),
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": "Description for the policy",
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": null,
    "policyId": 276858,
    "version": 2
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   ptr.To("Description for the policy"),
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				MatchRules:    nil,
				PolicyID:      276858,
				PolicyVersion: 2,
			},
		},
		"201 created, complex AS": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleAS{
							Start: 0,
							End:   0,
							Type:  "asMatchRule",
							Name:  "Q1Sales",
							ID:    0,
							ForwardSettings: ForwardSettingsAS{
								OriginID:               "originremote2",
								PathAndQS:              "/sales/Q1/",
								UseIncomingQueryString: true,
							},
							Matches: []MatchCriteriaAS{
								{
									CaseSensitive: false,
									MatchOperator: "equals",
									MatchType:     "range",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueRange{
										Type:  "range",
										Value: []int64{1, 25},
									},
								},
								{
									CaseSensitive: false,
									MatchOperator: "equals",
									MatchType:     "method",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueSimple{
										Type:  "simple",
										Value: []string{"GET"},
									},
								},
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueObject{
										Type: "object",
										Name: "AS",
										Options: &Options{
											Value: []string{
												"text/html*",
												"text/css*",
												"application/x-javascript*",
											},
											ValueHasWildcard: true,
										},
									},
								},
							},
						},
					},
				},
				PolicyID: 355557,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": "Initial version",
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "asMatchRule",
            "akaRuleId": "f58014ee0cc17ce",
            "end": 0,
            "forwardSettings": {
                "originId": "originremote2",
                "pathAndQS": "/sales/Q1/",
                "useIncomingQueryString": true
            },
            "id": 0,
            "matches": [
				{
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "range",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "range",
                        "value": [
                            1,
                            25
                        ]
                    }
                },
				{
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "simple",
                        "value": [
                            "GET"
                        ]
                    }
                },
				{
                    "matchOperator": "equals",
                    "matchType": "header",
                    "negate": false,
					"objectMatchValue": {
						"type": "object",
						"name": "AS",
						"options": {
							"value": [
								"text/html*",
                                "text/css*",
                                "application/x-javascript*"
							],
							"valueHasWildcard": true
						}
					}
                }
            ],
            "name": "Q1Sales",
            "start": 0
        }
    ],
    "policyId": 355557,
    "version": 2
}`,
			expectedPath: "/cloudlets/v3/policies/355557/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:  test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:    "jsmith",
				Description:  ptr.To("Initial version"),
				ModifiedBy:   "jsmith",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				MatchRules: MatchRules{
					&MatchRuleAS{
						Type: "asMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsAS{
							OriginID:               "originremote2",
							PathAndQS:              "/sales/Q1/",
							UseIncomingQueryString: true,
						},
						ID: 0,
						Matches: []MatchCriteriaAS{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "range",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueRange{
									Type: "range",
									Value: []int64{
										1, 25},
								},
							},
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "method",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueSimple{
									Type:  "simple",
									Value: []string{"GET"},
								},
							},
							{
								MatchOperator: "equals",
								MatchType:     "header",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObject{
									Type: "object",
									Name: "AS",
									Options: &Options{
										Value: []string{
											"text/html*",
											"text/css*",
											"application/x-javascript*",
										},
										ValueHasWildcard: true,
									},
								},
							},
						},
						Name:  "Q1Sales",
						Start: 0,
					},
				},
				PolicyID:      355557,
				PolicyVersion: 2,
			},
		},
		"201 created, complex PR": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRulePR{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsPR{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaPR{
								{
									MatchType:     "hostname",
									MatchValue:    "3333.dom",
									MatchOperator: "equals",
									CaseSensitive: true,
									Negate:        false,
								},
								{
									MatchType:     "cookie",
									MatchValue:    "cookie=cookievalue",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
								{
									MatchType:     "extension",
									MatchValue:    "txt",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
							},
						},
						&MatchRulePR{
							Start:    0,
							End:      0,
							Type:     "cdMatchRule",
							Name:     "rule 2",
							MatchURL: "ddd.aaa",
							ID:       0,
							ForwardSettings: ForwardSettingsPR{
								OriginID: "some_origin",
								Percent:  10,
							},
						},
						&MatchRulePR{
							Type:     "cdMatchRule",
							ID:       0,
							Name:     "r1",
							Start:    0,
							End:      0,
							MatchURL: "abc.com",
							ForwardSettings: ForwardSettingsPR{
								OriginID: "some_origin",
								Percent:  10,
							},
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "cdMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": true,
                    "matchOperator": "equals",
                    "matchType": "hostname",
                    "matchValue": "3333.dom",
                    "negate": false
                },
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "cookie",
                    "matchValue": "cookie=cookievalue",
                    "negate": false
                },
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "extension",
                    "matchValue": "txt",
                    "negate": false
                }
            ],
            "name": "rul3",
            "redirectURL": "/abc/sss",
            "start": 0,
            "forwardSettings": {
                "originId": "some_origin",
                "percent": 10
            }
        },
        {
            "type": "cdMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": "ddd.aaa",
            "name": "rule 2",
            "redirectURL": "sss.com",
            "start": 0,
            "statusCode": 301,
            "useIncomingQueryString": true,
            "useRelativeUrl": "none"
        },
        {
            "type": "cdMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": "abc.com",
            "name": "r1",
            "redirectURL": "/ddd",
            "start": 0,
            "statusCode": 301,
            "useIncomingQueryString": false,
            "useIncomingSchemeAndHost": true,
            "useRelativeUrl": "copy_scheme_hostname"
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRulePR{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						ForwardSettings: ForwardSettingsPR{
							OriginID: "some_origin",
							Percent:  10,
						},
						Matches: []MatchCriteriaPR{
							{
								MatchType:     "hostname",
								MatchValue:    "3333.dom",
								MatchOperator: "equals",
								CaseSensitive: true,
								Negate:        false,
							},
							{
								MatchType:     "cookie",
								MatchValue:    "cookie=cookievalue",
								MatchOperator: "equals",
								Negate:        false,
								CaseSensitive: false,
							},
							{
								MatchType:     "extension",
								MatchValue:    "txt",
								MatchOperator: "equals",
								Negate:        false,
								CaseSensitive: false,
							},
						},
					},
					&MatchRulePR{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "ddd.aaa",
						Name:     "rule 2",
						Start:    0,
					},
					&MatchRulePR{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "abc.com",
						Name:     "r1",
						Start:    0,
					},
				},
			},
		},
		"201 created, complex PR with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRulePR{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsPR{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaPR{
								{
									CaseSensitive: true,
									MatchOperator: "equals",
									MatchType:     "method",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueSimple{
										Type:  "simple",
										Value: []string{"GET"},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			requestBody:    `{"matchRules":[{"name":"rul3","type":"cdMatchRule","matches":[{"matchType":"method","matchOperator":"equals","caseSensitive":true,"negate":false,"objectMatchValue":{"type":"simple","value":["GET"]}}],"forwardSettings":{"originId":"some_origin","percent":10}}]}`,
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "cdMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": true,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "simple",
                        "value": [
                            "GET"
                        ]
                    }
                }
            ],
            "name": "rul3",
            "redirectURL": "/abc/sss",
            "start": 0,
            "forwardSettings": {
                "originId": "some_origin",
                "percent": 10
            }
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRulePR{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						ForwardSettings: ForwardSettingsPR{
							OriginID: "some_origin",
							Percent:  10,
						},
						Matches: []MatchCriteriaPR{
							{
								CaseSensitive: true,
								MatchOperator: "equals",
								MatchType:     "method",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueSimple{
									Type:  "simple",
									Value: []string{"GET"},
								},
							},
						},
					},
				},
			},
		},
		"201 created, complex PR with objectMatchValue - object": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRulePR{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsPR{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaPR{
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueObject{
										Type: "object",
										Name: "PR",
										Options: &Options{
											Value: []string{
												"text/html*",
												"text/css*",
												"application/x-javascript*",
											},
											ValueHasWildcard: true,
										},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "cdMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "matchOperator": "equals",
                    "matchType": "hostname",
                    "negate": false,
					"objectMatchValue": {
						"type": "object",
						"name": "PR",
						"options": {
							"value": [
								"text/html*",
                                "text/css*",
                                "application/x-javascript*"
							],
							"valueHasWildcard": true
						}
					}
                }
            ],
            "name": "rul3",
            "redirectURL": "/abc/sss",
            "start": 0,
            "forwardSettings": {
                "originId": "some_origin",
                "percent": 10
            }
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRulePR{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						ForwardSettings: ForwardSettingsPR{
							OriginID: "some_origin",
							Percent:  10,
						},
						Matches: []MatchCriteriaPR{
							{
								MatchOperator: "equals",
								MatchType:     "hostname",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObject{
									Type: "object",
									Name: "PR",
									Options: &Options{
										Value: []string{
											"text/html*",
											"text/css*",
											"application/x-javascript*",
										},
										ValueHasWildcard: true,
									},
								},
							},
						},
					},
				},
			},
		},
		"validation error, complex PR with unavailable objectMatchValue type - range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRulePR{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsPR{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaPR{
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueRange{
										Type:  "range",
										Value: []int64{1, 50},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			withError: ErrStructValidation,
		},
		"validation error, complex PR missing forwardSettings": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRulePR{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsPR{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaPR{
								{
									MatchType:     "hostname",
									MatchValue:    "3333.dom",
									MatchOperator: "equals",
									CaseSensitive: true,
									Negate:        false,
								},
								{
									MatchType:     "cookie",
									MatchValue:    "cookie=cookievalue",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
								{
									MatchType:     "extension",
									MatchValue:    "txt",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
							},
						},
						&MatchRulePR{
							Start:    0,
							End:      0,
							Type:     "cdMatchRule",
							Name:     "rule 2",
							MatchURL: "ddd.aaa",
							ID:       0,
						},
						&MatchRulePR{
							Type:     "cdMatchRule",
							ID:       0,
							Name:     "r1",
							Start:    0,
							End:      0,
							MatchURL: "abc.com",
							ForwardSettings: ForwardSettingsPR{
								OriginID: "some_origin",
								Percent:  10,
							},
						},
					},
				},
				PolicyID: 276858,
			},
			withError: ErrStructValidation,
		},

		"201 created, complex ER": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleER{
							Start:          0,
							End:            0,
							Type:           "erMatchRule",
							UseRelativeURL: "copy_scheme_hostname",
							Name:           "rul3",
							StatusCode:     307,
							RedirectURL:    "/abc/sss",
							ID:             0,
							Matches: []MatchCriteriaER{
								{
									MatchType:     "hostname",
									MatchValue:    "3333.dom",
									MatchOperator: "equals",
									CaseSensitive: true,
									Negate:        false,
								},
								{
									MatchType:     "cookie",
									MatchValue:    "cookie=cookievalue",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
								{
									MatchType:     "extension",
									MatchValue:    "txt",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
							},
						},
						&MatchRuleER{
							Start:                  0,
							End:                    0,
							Type:                   "erMatchRule",
							UseRelativeURL:         "none",
							Name:                   "rule 2",
							MatchURL:               "ddd.aaa",
							RedirectURL:            "sss.com",
							StatusCode:             301,
							UseIncomingQueryString: true,
							ID:                     0,
						},
						&MatchRuleER{
							Type:                     "erMatchRule",
							ID:                       0,
							Name:                     "r1",
							Start:                    0,
							End:                      0,
							MatchURL:                 "abc.com",
							StatusCode:               301,
							RedirectURL:              "/ddd",
							UseIncomingQueryString:   false,
							UseIncomingSchemeAndHost: true,
							UseRelativeURL:           "copy_scheme_hostname",
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "erMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": true,
                    "matchOperator": "equals",
                    "matchType": "hostname",
                    "matchValue": "3333.dom",
                    "negate": false
                },
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "cookie",
                    "matchValue": "cookie=cookievalue",
                    "negate": false
                },
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "extension",
                    "matchValue": "txt",
                    "negate": false
                }
            ],
            "name": "rul3",
            "redirectURL": "/abc/sss",
            "start": 0,
            "statusCode": 307,
            "useIncomingQueryString": false,
            "useIncomingSchemeAndHost": true,
            "useRelativeUrl": "copy_scheme_hostname"
        },
        {
            "type": "erMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": "ddd.aaa",
            "name": "rule 2",
            "redirectURL": "sss.com",
            "start": 0,
            "statusCode": 301,
            "useIncomingQueryString": true,
            "useRelativeUrl": "none"
        },
        {
            "type": "erMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": "abc.com",
            "name": "r1",
            "redirectURL": "/ddd",
            "start": 0,
            "statusCode": 301,
            "useIncomingQueryString": false,
            "useIncomingSchemeAndHost": true,
            "useRelativeUrl": "copy_scheme_hostname"
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleER{
						Type:                     "erMatchRule",
						End:                      0,
						ID:                       0,
						MatchURL:                 "",
						Name:                     "rul3",
						RedirectURL:              "/abc/sss",
						Start:                    0,
						StatusCode:               307,
						UseIncomingQueryString:   false,
						UseIncomingSchemeAndHost: true,
						UseRelativeURL:           "copy_scheme_hostname",
						Matches: []MatchCriteriaER{
							{
								MatchType:     "hostname",
								MatchValue:    "3333.dom",
								MatchOperator: "equals",
								CaseSensitive: true,
								Negate:        false,
							},
							{
								MatchType:     "cookie",
								MatchValue:    "cookie=cookievalue",
								MatchOperator: "equals",
								Negate:        false,
								CaseSensitive: false,
							},
							{
								MatchType:     "extension",
								MatchValue:    "txt",
								MatchOperator: "equals",
								Negate:        false,
								CaseSensitive: false,
							},
						},
					},
					&MatchRuleER{
						Type:                   "erMatchRule",
						End:                    0,
						ID:                     0,
						MatchURL:               "ddd.aaa",
						Name:                   "rule 2",
						RedirectURL:            "sss.com",
						Start:                  0,
						StatusCode:             301,
						UseIncomingQueryString: true,
						UseRelativeURL:         "none",
					},
					&MatchRuleER{
						Type:                     "erMatchRule",
						End:                      0,
						ID:                       0,
						MatchURL:                 "abc.com",
						Name:                     "r1",
						RedirectURL:              "/ddd",
						Start:                    0,
						StatusCode:               301,
						UseIncomingQueryString:   false,
						UseIncomingSchemeAndHost: true,
						UseRelativeURL:           "copy_scheme_hostname",
					},
				},
			},
		},
		"201 created, complex ER with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleER{
							Start:          0,
							End:            0,
							Type:           "erMatchRule",
							UseRelativeURL: "copy_scheme_hostname",
							Name:           "rul3",
							StatusCode:     307,
							RedirectURL:    "/abc/sss",
							ID:             0,
							Matches: []MatchCriteriaER{
								{
									CaseSensitive: true,
									MatchOperator: "equals",
									MatchType:     "method",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueSimple{
										Type:  "simple",
										Value: []string{"GET"},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "erMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": true,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "simple",
                        "value": [
                            "GET"
                        ]
                    }
                }
            ],
            "name": "rul3",
            "redirectURL": "/abc/sss",
            "start": 0,
            "statusCode": 307,
            "useIncomingQueryString": false,
            "useIncomingSchemeAndHost": true,
            "useRelativeUrl": "copy_scheme_hostname"
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleER{
						Type:                     "erMatchRule",
						End:                      0,
						ID:                       0,
						MatchURL:                 "",
						Name:                     "rul3",
						RedirectURL:              "/abc/sss",
						Start:                    0,
						StatusCode:               307,
						UseIncomingQueryString:   false,
						UseIncomingSchemeAndHost: true,
						UseRelativeURL:           "copy_scheme_hostname",
						Matches: []MatchCriteriaER{
							{
								CaseSensitive: true,
								MatchOperator: "equals",
								MatchType:     "method",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueSimple{
									Type:  "simple",
									Value: []string{"GET"},
								},
							},
						},
					},
				},
			},
		},
		"201 created, complex ER with objectMatchValue - object": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleER{
							Start:          0,
							End:            0,
							Type:           "erMatchRule",
							UseRelativeURL: "copy_scheme_hostname",
							Name:           "rul3",
							StatusCode:     307,
							RedirectURL:    "/abc/sss",
							ID:             0,
							Matches: []MatchCriteriaER{
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueObject{
										Type: "object",
										Name: "ER",
										Options: &Options{
											Value: []string{
												"text/html*",
												"text/css*",
												"application/x-javascript*",
											},
											ValueHasWildcard: true,
										},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "erMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "matchOperator": "equals",
                    "matchType": "hostname",
                    "negate": false,
					"objectMatchValue": {
						"type": "object",
						"name": "ER",
						"options": {
							"value": [
								"text/html*",
                                "text/css*",
                                "application/x-javascript*"
							],
							"valueHasWildcard": true
						}
					}
                }
            ],
            "name": "rul3",
            "redirectURL": "/abc/sss",
            "start": 0,
            "statusCode": 307,
            "useIncomingQueryString": false,
            "useIncomingSchemeAndHost": true,
            "useRelativeUrl": "copy_scheme_hostname"
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleER{
						Type:                     "erMatchRule",
						End:                      0,
						ID:                       0,
						MatchURL:                 "",
						Name:                     "rul3",
						RedirectURL:              "/abc/sss",
						Start:                    0,
						StatusCode:               307,
						UseIncomingQueryString:   false,
						UseIncomingSchemeAndHost: true,
						UseRelativeURL:           "copy_scheme_hostname",
						Matches: []MatchCriteriaER{
							{
								MatchOperator: "equals",
								MatchType:     "hostname",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObject{
									Type: "object",
									Name: "ER",
									Options: &Options{
										Value: []string{
											"text/html*",
											"text/css*",
											"application/x-javascript*",
										},
										ValueHasWildcard: true,
									},
								},
							},
						},
					},
				},
			},
		},
		"201 created, ER with empty/no useRelativeURL": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleER{
							Start:                  0,
							End:                    0,
							Type:                   "erMatchRule",
							Name:                   "rule 2",
							MatchURL:               "ddd.aaa",
							RedirectURL:            "sss.com",
							StatusCode:             301,
							UseIncomingQueryString: true,
							ID:                     0,
						},
						&MatchRuleER{
							Type:                     "erMatchRule",
							ID:                       0,
							Name:                     "r1",
							Start:                    0,
							End:                      0,
							MatchURL:                 "abc.com",
							StatusCode:               301,
							RedirectURL:              "/ddd",
							UseIncomingQueryString:   false,
							UseIncomingSchemeAndHost: true,
							UseRelativeURL:           "",
						},
						&MatchRuleER{
							Start:       0,
							End:         0,
							Type:        "erMatchRule",
							Name:        "rul3",
							StatusCode:  307,
							RedirectURL: "/abc/sss",
							ID:          0,
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "erMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": "ddd.aaa",
            "name": "rule 2",
            "redirectURL": "sss.com",
            "start": 0,
            "statusCode": 301,
            "useIncomingQueryString": true
        },
        {
            "type": "erMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": "abc.com",
            "name": "r1",
            "redirectURL": "/ddd",
            "start": 0,
            "statusCode": 301,
            "useIncomingQueryString": false,
            "useIncomingSchemeAndHost": true,
            "useRelativeUrl": "copy_scheme_hostname"
        },
		{
			"type": "erMatchRule",
            "end": 0,
            "id": 0,
			"name": "rul3",
			"redirectURL": "/abc/sss",
			"start": 0,
			"statusCode": 307
		}
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleER{
						Type:                   "erMatchRule",
						End:                    0,
						ID:                     0,
						MatchURL:               "ddd.aaa",
						Name:                   "rule 2",
						RedirectURL:            "sss.com",
						Start:                  0,
						StatusCode:             301,
						UseIncomingQueryString: true,
					},
					&MatchRuleER{
						Type:                     "erMatchRule",
						End:                      0,
						ID:                       0,
						MatchURL:                 "abc.com",
						Name:                     "r1",
						RedirectURL:              "/ddd",
						Start:                    0,
						StatusCode:               301,
						UseIncomingQueryString:   false,
						UseIncomingSchemeAndHost: true,
						UseRelativeURL:           "copy_scheme_hostname",
					},
					&MatchRuleER{
						Start:       0,
						End:         0,
						Type:        "erMatchRule",
						Name:        "rul3",
						StatusCode:  307,
						RedirectURL: "/abc/sss",
						ID:          0,
					},
				},
			},
		},
		"validation error, complex ER with unavailable objectMatchValue type - range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleER{
							Start:          0,
							End:            0,
							Type:           "erMatchRule",
							UseRelativeURL: "copy_scheme_hostname",
							Name:           "rul3",
							StatusCode:     307,
							RedirectURL:    "/abc/sss",
							ID:             0,
							Matches: []MatchCriteriaER{
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueRange{
										Type:  "range",
										Value: []int64{1, 50},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			withError: ErrStructValidation,
		},

		"201 created, complex FR": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleFR{
							Start: 0,
							End:   0,
							Type:  "frMatchRule",
							Name:  "rul3",
							ID:    0,
							Matches: []MatchCriteriaFR{
								{
									MatchType:     "hostname",
									MatchValue:    "3333.dom",
									MatchOperator: "equals",
									CaseSensitive: true,
									Negate:        false,
								},
								{
									MatchType:     "cookie",
									MatchValue:    "cookie=cookievalue",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
								{
									MatchType:     "extension",
									MatchValue:    "txt",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
							},
							ForwardSettings: ForwardSettingsFR{
								PathAndQS:              "/test_images/simpleimg.jpg",
								UseIncomingQueryString: true,
								OriginID:               "1234",
							},
						},
						&MatchRuleFR{
							Name:     "rule 1",
							Type:     "frMatchRule",
							Start:    0,
							End:      0,
							ID:       0,
							MatchURL: "ddd.aaa",
							ForwardSettings: ForwardSettingsFR{
								PathAndQS:              "/test_images/simpleimg.jpg",
								UseIncomingQueryString: true,
								OriginID:               "1234",
							},
						},
						&MatchRuleFR{
							Name:     "rule 2",
							Type:     "frMatchRule",
							Start:    0,
							End:      0,
							ID:       0,
							MatchURL: "abc.com",
							ForwardSettings: ForwardSettingsFR{
								PathAndQS:              "/test_images/otherimage.jpg",
								UseIncomingQueryString: true,
								OriginID:               "1234",
							},
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "frMatchRule",
            "akaRuleId": "893947a3d5a85c1b",
            "end": 0,
            "forwardSettings": {
                "pathAndQS": "/test_images/otherimage.jpg",
                "useIncomingQueryString": true,
				"originId": "1234"
            },
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": true,
                    "matchOperator": "equals",
                    "matchType": "hostname",
                    "matchValue": "3333.dom",
                    "negate": false
                },
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "cookie",
                    "matchValue": "cookie=cookievalue",
                    "negate": false
                },
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "extension",
                    "matchValue": "txt",
                    "negate": false
                }
            ],
            "name": "rul3",
            "start": 0
        },
        {
            "type": "frMatchRule",
            "akaRuleId": "aa379d230efcded0",
            "end": 0,
            "forwardSettings": {
                "pathAndQS": "/test_images/simpleimg.jpg",
                "useIncomingQueryString": true,
				"originId": "1234"
            },
            "id": 0,
            "matchURL": "ddd.aaa",
            "name": "rule 1",
            "start": 0
        },
        {
            "type": "frMatchRule",
            "akaRuleId": "1afe03d843996766",
            "end": 0,
            "forwardSettings": {
                "pathAndQS": "/test_images/otherimage.jpg",
                "useIncomingQueryString": true,
				"originId": "1234"
            },
            "id": 0,
            "matchURL": "abc.com",
            "name": "rule 2",
            "start": 0
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleFR{
						Type:     "frMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						Matches: []MatchCriteriaFR{
							{
								MatchType:     "hostname",
								MatchValue:    "3333.dom",
								MatchOperator: "equals",
								CaseSensitive: true,
								Negate:        false,
							},
							{
								MatchType:     "cookie",
								MatchValue:    "cookie=cookievalue",
								MatchOperator: "equals",
								Negate:        false,
								CaseSensitive: false,
							},
							{
								MatchType:     "extension",
								MatchValue:    "txt",
								MatchOperator: "equals",
								Negate:        false,
								CaseSensitive: false,
							},
						},
						ForwardSettings: ForwardSettingsFR{
							PathAndQS:              "/test_images/otherimage.jpg",
							UseIncomingQueryString: true,
							OriginID:               "1234",
						},
					},
					&MatchRuleFR{
						Name:     "rule 1",
						Type:     "frMatchRule",
						Start:    0,
						End:      0,
						ID:       0,
						MatchURL: "ddd.aaa",
						ForwardSettings: ForwardSettingsFR{
							PathAndQS:              "/test_images/simpleimg.jpg",
							UseIncomingQueryString: true,
							OriginID:               "1234",
						},
					},
					&MatchRuleFR{
						Name:     "rule 2",
						Type:     "frMatchRule",
						Start:    0,
						End:      0,
						ID:       0,
						MatchURL: "abc.com",
						ForwardSettings: ForwardSettingsFR{
							PathAndQS:              "/test_images/otherimage.jpg",
							UseIncomingQueryString: true,
							OriginID:               "1234",
						},
					},
				},
			},
		},
		"201 created, complex FR with objectMatchValue - object": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					Description: ptr.To("New version 1630480693371"),
					MatchRules: MatchRules{
						&MatchRuleFR{
							ForwardSettings: ForwardSettingsFR{},
							Matches: []MatchCriteriaFR{
								{
									CaseSensitive: false,
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueObject{
										Type:              "object",
										Name:              "Accept",
										NameCaseSensitive: false,
										NameHasWildcard:   false,
										Options: &Options{
											Value:              []string{"asd", "qwe"},
											ValueHasWildcard:   false,
											ValueCaseSensitive: true,
											ValueEscaped:       false,
										},
									},
								},
							},
							Start: 0,
							End:   0,
							Type:  "frMatchRule",
							Name:  "rul3",
							ID:    0,
						},
					},
				},
				PolicyID: 139743,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": "New version 1630480693371",
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "frMatchRule",
            "akaRuleId": "f2168e71692e6d9f",
            "end": 0,
            "forwardSettings": {},
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "header",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "object",
                        "name": "Accept",
                        "options": {
                            "value": [
                                "asd",
                                "qwe"
                            ],
                            "valueCaseSensitive": true
                        }
                    }
                }
            ],
            "name": "rul3",
            "start": 0
        }
    ],
	"policyId": 139743,
    "version": 798
}`,
			expectedPath: "/cloudlets/v3/policies/139743/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:  test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:    "jsmith",
				Description:  ptr.To("New version 1630480693371"),
				ModifiedBy:   "jsmith",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				MatchRules: MatchRules{
					&MatchRuleFR{
						ForwardSettings: ForwardSettingsFR{},
						Matches: []MatchCriteriaFR{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "header",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObject{
									Name:              "Accept",
									Type:              "object",
									NameCaseSensitive: false,
									NameHasWildcard:   false,
									Options: &Options{
										Value:              []string{"asd", "qwe"},
										ValueHasWildcard:   false,
										ValueCaseSensitive: true,
										ValueEscaped:       false,
									},
								},
							},
						},
						Start: 0,
						End:   0,
						Type:  "frMatchRule",
						Name:  "rul3",
						ID:    0,
					},
				},
				PolicyID:      139743,
				PolicyVersion: 798,
			},
		},
		"201 created, complex FR with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					Description: ptr.To("New version 1630480693371"),
					MatchRules: MatchRules{
						&MatchRuleFR{
							ForwardSettings: ForwardSettingsFR{
								PathAndQS:              "/test_images/otherimage.jpg",
								UseIncomingQueryString: true,
							},
							Matches: []MatchCriteriaFR{
								{
									CaseSensitive: false,
									MatchOperator: "equals",
									MatchType:     "method",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueSimple{
										Type:  "simple",
										Value: []string{"GET"},
									},
								},
							},
							Start: 0,
							End:   0,
							Type:  "frMatchRule",
							Name:  "rul3",
							ID:    0,
						},
					},
				},
				PolicyID: 139743,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": "New version 1630480693371",
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "frMatchRule",
            "akaRuleId": "f2168e71692e6d9f",
            "end": 0,
            "forwardSettings": {
                "pathAndQS": "/test_images/otherimage.jpg",
                "useIncomingQueryString": true
            },
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "simple",
                        "value": [
                            "GET"
                        ]
                    }
                }
			],
            "name": "rul3",
            "start": 0
        }
    ],
	"policyId": 139743,
    "version": 798
}`,
			expectedPath: "/cloudlets/v3/policies/139743/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:  test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:    "jsmith",
				Description:  ptr.To("New version 1630480693371"),
				ModifiedBy:   "jsmith",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				MatchRules: MatchRules{
					&MatchRuleFR{
						ForwardSettings: ForwardSettingsFR{
							PathAndQS:              "/test_images/otherimage.jpg",
							UseIncomingQueryString: true,
						},
						Matches: []MatchCriteriaFR{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "method",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueSimple{
									Type:  "simple",
									Value: []string{"GET"},
								},
							},
						},
						Start: 0,
						End:   0,
						Type:  "frMatchRule",
						Name:  "rul3",
						ID:    0,
					},
				},
				PolicyID:      139743,
				PolicyVersion: 798,
			},
		},
		"201 created, complex AP with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleAP{
							Start:              0,
							End:                0,
							Type:               "apMatchRule",
							Name:               "rul3",
							PassThroughPercent: ptr.To(float64(0)),
							ID:                 0,
							Matches: []MatchCriteriaAP{
								{
									CaseSensitive: true,
									MatchOperator: "equals",
									MatchType:     "method",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueSimple{
										Type:  "simple",
										Value: []string{"GET"},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "apMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": true,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "simple",
                        "value": [
                            "GET"
                        ]
                    }
                }
            ],
            "name": "rul3",
            "start": 0,
            "useIncomingQueryString": false,
			"passThroughPercent": -1
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleAP{
						Type:               "apMatchRule",
						End:                0,
						ID:                 0,
						MatchURL:           "",
						Name:               "rul3",
						PassThroughPercent: ptr.To(float64(-1)),
						Start:              0,
						Matches: []MatchCriteriaAP{
							{
								CaseSensitive: true,
								MatchOperator: "equals",
								MatchType:     "method",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueSimple{
									Type:  "simple",
									Value: []string{"GET"},
								},
							},
						},
					},
				},
			},
		},
		"201 created, complex AP with objectMatchValue - object": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleAP{
							Start:              0,
							End:                0,
							Type:               "apMatchRule",
							Name:               "rul3",
							PassThroughPercent: ptr.To(float64(-1)),
							ID:                 0,
							Matches: []MatchCriteriaAP{
								{
									CaseSensitive: false,
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueObject{
										Type: "object",
										Name: "AP",
										Options: &Options{
											Value:            []string{"y"},
											ValueHasWildcard: true,
										},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": null,
    "matchRules": [
        {
            "type": "apMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "header",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "object",
                        "name": "AP",
						"options": {
                            "value": [
                                "y"
                            ],
                            "valueHasWildcard": true
                        }
                    }
                }
            ],
            "name": "rul3",
            "start": 0,
			"passThroughPercent": -1
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  nil,
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleAP{
						Type:               "apMatchRule",
						End:                0,
						ID:                 0,
						MatchURL:           "",
						Name:               "rul3",
						PassThroughPercent: ptr.To(float64(-1)),
						Start:              0,
						Matches: []MatchCriteriaAP{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "header",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObject{
									Type: "object",
									Name: "AP",
									Options: &Options{
										Value:            []string{"y"},
										ValueHasWildcard: true,
									},
								},
							},
						},
					},
				},
			},
		},
		"201 created, complex RC": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleRC{
							Start:     0,
							End:       0,
							Type:      "igMatchRule",
							Name:      "rul3",
							AllowDeny: DenyBranded,
							ID:        0,
							Matches: []MatchCriteriaRC{
								{
									CaseSensitive: false,
									MatchOperator: "equals",
									MatchType:     "protocol",
									MatchValue:    "https",
									Negate:        false,
								},
								{
									CaseSensitive: true,
									MatchOperator: "equals",
									MatchType:     "method",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueSimple{
										Type:  "simple",
										Value: []string{"GET"},
									},
								},
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueObject{
										Type: "object",
										Name: "RC",
										Options: &Options{
											Value: []string{
												"text/html*",
												"text/css*",
												"application/x-javascript*",
											},
											ValueHasWildcard: true,
										},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": null,
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": [
        {
            "type": "igMatchRule",
            "end": 0,
            "id": 0,
            "matchesAlways": false,
            "matches": [
				{
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "protocol",
                    "negate": false,
					"matchValue": "https"
				},
                {
                    "caseSensitive": true,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "simple",
                        "value": [
                            "GET"
                        ]
                    }
                },
				{
                    "matchOperator": "equals",
                    "matchType": "header",
                    "negate": false,
					"objectMatchValue": {
                        "type": "object",
                        "name": "RC",
						"options": {
							"value": [
								"text/html*",
								"text/css*",
								"application/x-javascript*"
                            ],
                            "valueHasWildcard": true
						}
                    }
				}
            ],
            "name": "rul3",
            "start": 0,
			"allowDeny": "denybranded"
        }
    ],
    "policyId": 276858,
    "version": 6
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   nil,
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				PolicyID:      276858,
				PolicyVersion: 6,
				MatchRules: MatchRules{
					&MatchRuleRC{
						Type:          "igMatchRule",
						End:           0,
						ID:            0,
						MatchesAlways: false,
						Name:          "rul3",
						AllowDeny:     DenyBranded,
						Start:         0,
						Matches: []MatchCriteriaRC{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "protocol",
								MatchValue:    "https",
								Negate:        false,
							},
							{
								CaseSensitive: true,
								MatchOperator: "equals",
								MatchType:     "method",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueSimple{
									Type:  "simple",
									Value: []string{"GET"},
								},
							},
							{
								MatchOperator: "equals",
								MatchType:     "header",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObject{
									Type: "object",
									Name: "RC",
									Options: &Options{
										Value: []string{
											"text/html*",
											"text/css*",
											"application/x-javascript*",
										},
										ValueHasWildcard: true,
									},
								},
							},
						},
					},
				},
			},
		},
		"validation error, complex RC with unavailable objectMatchValue type - range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleRC{
							Start:     0,
							End:       0,
							Type:      "igMatchRule",
							AllowDeny: Allow,
							Name:      "rul3",
							ID:        0,
							Matches: []MatchCriteriaRC{
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueRange{
										Type:  "range",
										Value: []int64{1, 50},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			withError: ErrStructValidation,
		},
		"validation error, complex AP with unavailable objectMatchValue type - range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleAP{
							Start:              0,
							End:                0,
							Type:               "apMatchRule",
							PassThroughPercent: ptr.To(50.50),
							Name:               "rul3",
							ID:                 0,
							Matches: []MatchCriteriaAP{
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueRange{
										Type:  "range",
										Value: []int64{1, 50},
									},
								},
							},
						},
					},
				},
				PolicyID: 276858,
			},
			withError: ErrStructValidation,
		},

		"validation error, simple RC missing allowDeny": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleRC{
							Start: 0,
							End:   0,
							Type:  "igMatchRule",
							Name:  "rul3",
							ID:    0,
						},
					},
				},
				PolicyID: 276858,
			},
			withError: ErrStructValidation,
		},

		"validation error, simple AP missing passThrughPercent": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleAP{
							Start: 0,
							End:   0,
							Type:  "apMatchRule",
							Name:  "rul3",
							ID:    0,
						},
					},
				},
				PolicyID: 276858,
			},
			withError: ErrStructValidation,
		},

		"validation error, simple AP passThroughPercent out of range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleAP{
							Start:              0,
							End:                0,
							Type:               "apMatchRule",
							PassThroughPercent: ptr.To(float64(101)),
							Name:               "rul3",
							ID:                 0,
						},
					},
				},
				PolicyID: 276858,
			},
			withError: ErrStructValidation,
		},

		"500 internal server error": {
			request: CreatePolicyVersionRequest{
				PolicyID: 1,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "status": 500
}`,
			expectedPath: "/cloudlets/v3/policies/1/versions",
			withError: &Error{
				Type:   "internal_error",
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				if test.requestBody != "" {
					buf := new(bytes.Buffer)
					_, err := buf.ReadFrom(r.Body)
					assert.NoError(t, err)
					req := buf.String()
					assert.Equal(t, test.requestBody, req)
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreatePolicyVersion(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeletePolicyVersion(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request        DeletePolicyVersionRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 no content": {
			request: DeletePolicyVersionRequest{
				PolicyID:      276858,
				PolicyVersion: 5,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/cloudlets/v3/policies/276858/versions/5",
		},

		"missing required fields": {
			request:   DeletePolicyVersionRequest{},
			withError: ErrStructValidation,
		},

		"500 internal server error": {
			request: DeletePolicyVersionRequest{
				PolicyID:      1,
				PolicyVersion: 2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "status": 500
}`,
			expectedPath: "/cloudlets/v3/policies/1/versions/2",
			withError: &Error{
				Type:   "internal_error",
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeletePolicyVersion(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdatePolicyVersion(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		request          UpdatePolicyVersionRequest
		requestBody      string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PolicyVersion
		withError        error
	}{
		"201 updated simple ER": {
			request: UpdatePolicyVersionRequest{
				UpdatePolicyVersion: UpdatePolicyVersion{
					Description: ptr.To("Updated description"),
				},
				PolicyID:      276858,
				PolicyVersion: 5,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "createdDate": "2023-10-19T08:50:47.350Z",
    "createdBy": "jsmith",
    "description": "Updated description",
    "modifiedBy": "jsmith",
    "modifiedDate": "2023-10-19T08:50:47.350Z",
    "matchRules": null,
    "policyId": 276858,
    "version": 5
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions/5",
			expectedResponse: &PolicyVersion{
				CreatedDate:   test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z"),
				CreatedBy:     "jsmith",
				Description:   ptr.To("Updated description"),
				ModifiedBy:    "jsmith",
				ModifiedDate:  ptr.To(test.NewTimeFromString(t, "2023-10-19T08:50:47.350Z")),
				MatchRules:    nil,
				PolicyID:      276858,
				PolicyVersion: 5,
			},
		},
		"201 updated simple ER with warnings": {
			request: UpdatePolicyVersionRequest{
				UpdatePolicyVersion: UpdatePolicyVersion{
					Description: ptr.To("Updated description"),
					MatchRules: MatchRules{
						&MatchRuleER{
							Name:          "er_rule",
							Type:          "erMatchRule",
							Matches:       []MatchCriteriaER{},
							MatchesAlways: false,
							StatusCode:    301,
							RedirectURL:   "/path",
						},
					},
				},
				PolicyID:      276858,
				PolicyVersion: 5,
			},
			requestBody: `
{
	"description": "Updated description",
	"id": 276858,
	"matchRules": [
		{
			"type": "erMatchRule",
			"name": "er_rule",
			"matchURL": null,
			"statusCode": 301,
			"redirectURL": "/path",
			"useIncomingQueryString": false
		}
	]
}`,
			responseStatus: http.StatusOK,
			responseBody: `
{
	"createdBy": "jsmith",
	"createdDate": "2023-10-20T09:21:04.180Z",
	"description": "Updated description",
	"id": 6552305,
	"immutable": false,
	"matchRules": [
		{
			"type": "erMatchRule",
			"id": 0,
			"name": "er_rule",
			"start": 0,
			"end": 0,
			"matchURL": null,
			"statusCode": 301,
			"redirectURL": "/path",
			"useIncomingQueryString": false
		}
	],
	"matchRulesWarnings": [
		{
			"detail": "No match match conditions.",
			"title": "Missing Match Criteria",
			"type": "/cloudlets/error-types/missing-match-criteria",
			"jsonPointer": "/matchRules/0"
		}
	],
	"modifiedBy": "jsmith",
	"modifiedDate": "2023-10-20T10:32:56.316Z",
	"policyId": 670831,
	"version": 3
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions/5",
			expectedResponse: &PolicyVersion{
				CreatedDate:  test.NewTimeFromString(t, "2023-10-20T09:21:04.180Z"),
				CreatedBy:    "jsmith",
				Description:  ptr.To("Updated description"),
				ModifiedBy:   "jsmith",
				ModifiedDate: ptr.To(test.NewTimeFromString(t, "2023-10-20T10:32:56.316Z")),
				ID:           6552305,
				MatchRules: MatchRules{
					&MatchRuleER{
						Name:          "er_rule",
						Type:          "erMatchRule",
						MatchesAlways: false,
						StatusCode:    301,
						RedirectURL:   "/path",
					},
				},
				PolicyID:      670831,
				PolicyVersion: 3,
				MatchRulesWarnings: []MatchRulesWarning{
					{
						Detail:      "No match match conditions.",
						Title:       "Missing Match Criteria",
						Type:        "/cloudlets/error-types/missing-match-criteria",
						JSONPointer: "/matchRules/0",
					},
				},
			},
		},
		"500 internal server error": {
			request: UpdatePolicyVersionRequest{
				PolicyID:      276858,
				PolicyVersion: 3,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "status": 500
}`,
			expectedPath: "/cloudlets/v3/policies/276858/versions/3",
			withError: &Error{
				Type:   "internal_error",
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			},
		},
		"validation error": {
			expectedPath: "/cloudlets/v3/policies/0",
			request: UpdatePolicyVersionRequest{
				UpdatePolicyVersion: UpdatePolicyVersion{
					Description: ptr.To(strings.Repeat("A", 256)),
				},
			},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdatePolicyVersion(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
