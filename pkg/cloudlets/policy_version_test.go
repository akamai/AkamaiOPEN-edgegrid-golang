package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestGetPolicyVersion(t *testing.T) {
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
				PolicyID:  276858,
				Version:   1,
				OmitRules: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "activations": [],
    "createDate": 1629299944291,
    "createdBy": "agrebenk",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "agrebenk",
    "lastModifiedDate": 1629299944291,
    "location": "/cloudlets/api/v2/policies/276858/versions/1",
    "matchRuleFormat": "1.0",
    "matchRules": null,
    "policyId": 276858,
    "revisionId": 4811534,
    "rulesLocked": false,
    "version": 1
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/1?omitRules=true",
			expectedResponse: &PolicyVersion{
				Activations:      []*Activation{},
				CreateDate:       1629299944291,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1629299944291,
				Location:         "/cloudlets/api/v2/policies/276858/versions/1",
				MatchRuleFormat:  "1.0",
				MatchRules:       nil,
				PolicyID:         276858,
				RevisionID:       4811534,
				RulesLocked:      false,
				Version:          1,
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
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/1/versions/0?omitRules=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
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
	tests := map[string]struct {
		request          CreatePolicyVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PolicyVersion
		withError        error
	}{
		"201 created, simple ER": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					Description: "Description for the policy",
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "activations": [],
    "createDate": 1629812554924,
    "createdBy": "agrebenk",
    "deleted": false,
    "description": "Description for the policy",
    "lastModifiedBy": "agrebenk",
    "lastModifiedDate": 1629812554924,
    "location": "/cloudlets/api/v2/policies/276858/versions/2",
    "matchRuleFormat": "1.0",
    "matchRules": null,
    "policyId": 276858,
    "revisionId": 4814868,
    "rulesLocked": false,
    "version": 2
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []*Activation{},
				CreateDate:       1629812554924,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "Description for the policy",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1629812554924,
				Location:         "/cloudlets/api/v2/policies/276858/versions/2",
				MatchRuleFormat:  "1.0",
				MatchRules:       nil,
				PolicyID:         276858,
				RevisionID:       4814868,
				RulesLocked:      false,
				Version:          2,
			},
		},

		"201 created, complex ALB": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleALB{
							Matches: []MatchCriteriaALB{
								{
									MatchType:     "protocol",
									MatchValue:    "https",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
								{
									MatchType:     "path",
									MatchValue:    "/nonalb",
									MatchOperator: "equals",
									Negate:        true,
									CaseSensitive: false,
								},
							},
							Start: 0,
							End:   0,
							Type:  "albMatchRule",
							Name:  "Rule3",
							ForwardSettings: ForwardSettings{
								OriginID: "alb_test_krk_dc1_only",
							},
							ID: 0,
						},
						&MatchRuleALB{
							Matches: []MatchCriteriaALB{
								{
									MatchType:     "protocol",
									MatchValue:    "http",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
							},
							Start: 0,
							End:   0,
							Type:  "albMatchRule",
							Name:  "Rule1",
							ForwardSettings: ForwardSettings{
								OriginID: "alb_test_krk_0_100",
							},
							ID: 0,
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "activations": [],
    "createDate": 1629981546401,
    "createdBy": "agrebenk",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "agrebenk",
    "lastModifiedDate": 1629981546401,
    "location": "/cloudlets/api/v2/policies/279628/versions/2",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "akaRuleId": "8a57dcbd5565cc9a",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_dc1_only"
            },
            "id": 0,
            "location": "/cloudlets/api/v2/policies/279628/versions/2/rules/8a57dcbd5565cc9a",
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "protocol",
                    "matchValue": "https",
                    "negate": false
                },
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "path",
                    "matchValue": "/nonalb",
                    "negate": true
                }
            ],
            "name": "Rule3",
            "start": 0
        },
        {
            "type": "albMatchRule",
            "akaRuleId": "c018ee7c534b568c",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_0_100"
            },
            "id": 0,
            "location": "/cloudlets/api/v2/policies/279628/versions/2/rules/c018ee7c534b568c",
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "protocol",
                    "matchValue": "http",
                    "negate": false
                }
            ],
            "name": "Rule1",
            "start": 0
        }
    ],
    "policyId": 279628,
    "revisionId": 4815971,
    "rulesLocked": false,
    "version": 2
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []*Activation{},
				CreateDate:       1629981546401,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1629981546401,
				Location:         "/cloudlets/api/v2/policies/279628/versions/2",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type:      "albMatchRule",
						AkaRuleID: "8a57dcbd5565cc9a",
						End:       0,
						ForwardSettings: ForwardSettings{
							OriginID: "alb_test_krk_dc1_only",
						},
						ID:       0,
						Location: "/cloudlets/api/v2/policies/279628/versions/2/rules/8a57dcbd5565cc9a",
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "protocol",
								MatchValue:    "https",
								Negate:        false,
							},
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "path",
								MatchValue:    "/nonalb",
								Negate:        true,
							},
						},
						Name:  "Rule3",
						Start: 0,
					},
					&MatchRuleALB{
						Type:      "albMatchRule",
						AkaRuleID: "c018ee7c534b568c",
						End:       0,
						ForwardSettings: ForwardSettings{
							OriginID: "alb_test_krk_0_100",
						},
						ID:       0,
						Location: "/cloudlets/api/v2/policies/279628/versions/2/rules/c018ee7c534b568c",
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "protocol",
								MatchValue:    "http",
								Negate:        false,
							},
						},
						Name:  "Rule1",
						Start: 0,
					},
				},
				PolicyID:    279628,
				RevisionID:  4815971,
				RulesLocked: false,
				Version:     2,
			},
		},

		"201 created, complex ALB with objectMatchValue - object": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					Description:     "New version 1630480693371",
					MatchRuleFormat: "1.0",
					MatchRules: MatchRules{
						&MatchRuleALB{
							Matches: []MatchCriteriaALB{
								{
									MatchOperator: "equals",
									MatchType:     "header",
									ObjectMatchValue: ObjectMatchValueObjectSubtype{
										Type: "object",
										Name: "ALB",
										Options: &Options{
											Value: []interface{}{
												"y",
											},
											ValueHasWildcard: true,
										},
									},
									Negate: false,
								},
							},
							Start: 0,
							End:   0,
							Type:  "albMatchRule",
							Name:  "alb rule",
							ForwardSettings: ForwardSettings{
								OriginID: "alb_test_krk_mutable",
							},
							ID:            0,
							MatchesAlways: true,
						},
					},
				},
				PolicyID: 139743,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
		{
		    "activations": [],
		    "createDate": 1630489884742,
		    "createdBy": "agrebenk",
		    "deleted": false,
		    "description": "New version 1630480693371",
		    "lastModifiedBy": "agrebenk",
		    "lastModifiedDate": 1630489884742,
		    "location": "/cloudlets/api/v2/policies/139743/versions/796",
		    "matchRuleFormat": "1.0",
		    "matchRules": [
		        {
		            "type": "albMatchRule",
		            "akaRuleId": "93869097f18722c9",
		            "end": 0,
		            "forwardSettings": {
		                "originId": "alb_test_krk_mutable"
		            },
		            "id": 0,
		            "location": "/cloudlets/api/v2/policies/139743/versions/796/rules/93869097f18722c9",
		            "matchURL": null,
		            "matches": [
		                {
		                    "caseSensitive": false,
		                    "matchOperator": "equals",
		                    "matchType": "header",
		                    "negate": false,
		                    "objectMatchValue": {
		                        "type": "object",
		                        "name": "ALB",
		                        "options": {
		                            "value": [
		                                "y"
		                            ],
		                            "valueHasWildcard": true
		                        }
		                    }
		                }
		            ],
		            "matchesAlways": true,
		            "name": "alb rule",
		            "start": 0
		        }
		    ],
		    "policyId": 139743,
		    "revisionId": 4819432,
		    "rulesLocked": false,
		    "version": 796
		}`,
			expectedPath: "/cloudlets/api/v2/policies/139743/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []*Activation{},
				CreateDate:       1630489884742,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "New version 1630480693371",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1630489884742,
				Location:         "/cloudlets/api/v2/policies/139743/versions/796",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type:      "albMatchRule",
						AkaRuleID: "93869097f18722c9",
						End:       0,
						ForwardSettings: ForwardSettings{
							OriginID: "alb_test_krk_mutable",
						},
						ID:       0,
						Location: "/cloudlets/api/v2/policies/139743/versions/796/rules/93869097f18722c9",
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "header",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObjectSubtype{
									Type: "object",
									Name: "ALB",
									Options: &Options{
										Value: []interface{}{
											"y",
										},
										ValueHasWildcard: true,
									},
								},
							},
						},
						MatchesAlways: true,
						Name:          "alb rule",
						Start:         0,
					},
				},
				PolicyID:    139743,
				RevisionID:  4819432,
				RulesLocked: false,
				Version:     796,
			},
		},

		"201 created, complex ALB with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{

					Description:     "New version 1630480693371",
					MatchRuleFormat: "1.0",
					MatchRules: MatchRules{
						&MatchRuleALB{
							Matches: []MatchCriteriaALB{
								{
									CaseSensitive: false,
									MatchOperator: "equals",
									MatchType:     "method",
									Negate:        false,
									ObjectMatchValue: ObjectMatchValueRangeOrSimpleSubtype{
										Type:  "simple",
										Value: []interface{}{"GET"},
									},
								},
							},
							Start: 0,
							End:   0,
							Type:  "albMatchRule",
							Name:  "alb rule",
							ForwardSettings: ForwardSettings{
								OriginID: "alb_test_krk_mutable",
							},
							ID:            0,
							MatchesAlways: true,
						},
					},
				},
				PolicyID: 139743,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "activations": [],
    "createDate": 1630506044027,
    "createdBy": "agrebenk",
    "deleted": false,
    "description": "New version 1630480693371",
    "lastModifiedBy": "agrebenk",
    "lastModifiedDate": 1630506044027,
    "location": "/cloudlets/api/v2/policies/139743/versions/797",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "akaRuleId": "43776ac25e98e869",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_mutable"
            },
            "id": 0,
            "location": "/cloudlets/api/v2/policies/139743/versions/797/rules/43776ac25e98e869",
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
            "matchesAlways": true,
            "name": "alb rule",
            "start": 0
        }
    ],
    "policyId": 139743,
    "revisionId": 4819449,
    "rulesLocked": false,
    "version": 797
}`,
			expectedPath: "/cloudlets/api/v2/policies/139743/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []*Activation{},
				CreateDate:       1630506044027,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "New version 1630480693371",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1630506044027,
				Location:         "/cloudlets/api/v2/policies/139743/versions/797",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type:      "albMatchRule",
						AkaRuleID: "43776ac25e98e869",
						End:       0,
						ForwardSettings: ForwardSettings{
							OriginID: "alb_test_krk_mutable",
						},
						ID:       0,
						Location: "/cloudlets/api/v2/policies/139743/versions/797/rules/43776ac25e98e869",
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "method",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueRangeOrSimpleSubtype{
									Type: "simple",
									Value: []interface{}{
										"GET",
									},
								},
							},
						},
						MatchesAlways: true,
						Name:          "alb rule",
						Start:         0,
					},
				},
				PolicyID:    139743,
				RevisionID:  4819449,
				RulesLocked: false,
				Version:     797,
			},
		},

		"201 created, complex ALB with objectMatchValue - range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{

					Description:     "New version 1630480693371",
					MatchRuleFormat: "1.0",
					MatchRules: MatchRules{
						&MatchRuleALB{
							Matches: []MatchCriteriaALB{
								{
									CaseSensitive: false,
									MatchOperator: "equals",
									MatchType:     "range",
									Negate:        false,
									ObjectMatchValue: ObjectMatchValueRangeOrSimpleSubtype{
										Type:  "range",
										Value: []interface{}{1, 50},
									},
								},
							},
							Start: 0,
							End:   0,
							Type:  "albMatchRule",
							Name:  "alb rule",
							ForwardSettings: ForwardSettings{
								OriginID: "alb_test_krk_mutable",
							},
							ID:            0,
							MatchesAlways: true,
						},
					},
				},
				PolicyID: 139743,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "activations": [],
    "createDate": 1630507099511,
    "createdBy": "agrebenk",
    "deleted": false,
    "description": "New version 1630480693371",
    "lastModifiedBy": "agrebenk",
    "lastModifiedDate": 1630507099511,
    "location": "/cloudlets/api/v2/policies/139743/versions/798",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "akaRuleId": "69ace82d9db2ce48",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_mutable"
            },
            "id": 0,
            "location": "/cloudlets/api/v2/policies/139743/versions/798/rules/69ace82d9db2ce48",
            "matchURL": null,
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
                            50
                        ]
                    }
                }
            ],
            "matchesAlways": true,
            "name": "alb rule",
            "start": 0
        }
    ],
    "policyId": 139743,
    "revisionId": 4819450,
    "rulesLocked": false,
    "version": 798
}`,
			expectedPath: "/cloudlets/api/v2/policies/139743/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []*Activation{},
				CreateDate:       1630507099511,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "New version 1630480693371",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1630507099511,
				Location:         "/cloudlets/api/v2/policies/139743/versions/798",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type:      "albMatchRule",
						AkaRuleID: "69ace82d9db2ce48",
						End:       0,
						ForwardSettings: ForwardSettings{
							OriginID: "alb_test_krk_mutable",
						},
						ID:       0,
						Location: "/cloudlets/api/v2/policies/139743/versions/798/rules/69ace82d9db2ce48",
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "range",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueRangeOrSimpleSubtype{
									Type:  "range",
									Value: []interface{}{float64(1), float64(50)},
								},
							},
						},
						MatchesAlways: true,
						Name:          "alb rule",
						Start:         0,
					},
				},
				PolicyID:    139743,
				RevisionID:  4819450,
				RulesLocked: false,
				Version:     798,
			},
		},

		"201 created, complex ER": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleER{
							Start:          0,
							End:            0,
							Type:           "erMatchRule",
							UseRelativeUrl: "copy_scheme_hostname",
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
							UseRelativeUrl:         "none",
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
							AkaRuleID:                "e1969ed65202167f",
							StatusCode:               301,
							RedirectURL:              "/ddd",
							UseIncomingQueryString:   false,
							UseIncomingSchemeAndHost: true,
							UseRelativeUrl:           "copy_scheme_hostname",
						},
					},
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "agrebenk",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "agrebenk",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "erMatchRule",
            "akaRuleId": "a58392a7a43f19a3",
            "end": 0,
            "id": 0,
            "location": "/cloudlets/api/v2/policies/276858/versions/6/rules/a58392a7a43f19a3",
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
            "akaRuleId": "e38515c6542d2ed8",
            "end": 0,
            "id": 0,
            "location": "/cloudlets/api/v2/policies/276858/versions/6/rules/e38515c6542d2ed8",
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
            "akaRuleId": "e1969ed65202167f",
            "end": 0,
            "id": 0,
            "location": "/cloudlets/api/v2/policies/276858/versions/6/rules/e1969ed65202167f",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []*Activation{},
				CreateDate:       1629981355165,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleER{
						Type:                     "erMatchRule",
						AkaRuleID:                "a58392a7a43f19a3",
						End:                      0,
						ID:                       0,
						Location:                 "/cloudlets/api/v2/policies/276858/versions/6/rules/a58392a7a43f19a3",
						MatchURL:                 "",
						Name:                     "rul3",
						RedirectURL:              "/abc/sss",
						Start:                    0,
						StatusCode:               307,
						UseIncomingQueryString:   false,
						UseIncomingSchemeAndHost: true,
						UseRelativeUrl:           "copy_scheme_hostname",
						Matches: []MatchCriteriaER{
							{
								CaseSensitive: true,
								MatchOperator: "equals",
								MatchType:     "hostname",
								MatchValue:    "3333.dom",
								Negate:        false,
							},
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "cookie",
								MatchValue:    "cookie=cookievalue",
								Negate:        false,
							},
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "extension",
								MatchValue:    "txt",
								Negate:        false,
							},
						},
					},
					&MatchRuleER{
						Type:                   "erMatchRule",
						AkaRuleID:              "e38515c6542d2ed8",
						End:                    0,
						ID:                     0,
						Location:               "/cloudlets/api/v2/policies/276858/versions/6/rules/e38515c6542d2ed8",
						MatchURL:               "ddd.aaa",
						Name:                   "rule 2",
						RedirectURL:            "sss.com",
						Start:                  0,
						StatusCode:             301,
						UseIncomingQueryString: true,
						UseRelativeUrl:         "none",
					},
					&MatchRuleER{
						Type:                     "erMatchRule",
						AkaRuleID:                "e1969ed65202167f",
						End:                      0,
						ID:                       0,
						Location:                 "/cloudlets/api/v2/policies/276858/versions/6/rules/e1969ed65202167f",
						MatchURL:                 "abc.com",
						Name:                     "r1",
						RedirectURL:              "/ddd",
						Start:                    0,
						StatusCode:               301,
						UseIncomingQueryString:   false,
						UseIncomingSchemeAndHost: true,
						UseRelativeUrl:           "copy_scheme_hostname",
					},
				},
			},
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
  "detail": "Error creating enrollment",
  "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/1/versions",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRuleFormat: "2.0",
				},
			},
			expectedPath: "/cloudlets/api/v2/policies/0/versions",
			withError:    ErrStructValidation,
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
	tests := map[string]struct {
		request        DeletePolicyVersionRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 no content": {
			request: DeletePolicyVersionRequest{
				PolicyID: 276858,
				Version:  5,
			},
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/cloudlets/api/v2/policies/276858/versions/5",
		},
		"500 internal server error": {
			request: DeletePolicyVersionRequest{
				PolicyID: 1,
				Version:  2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error creating enrollment",
  "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/1/versions/2",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
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
	tests := map[string]struct {
		request          UpdatePolicyVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PolicyVersion
		withError        error
	}{
		"201 updated simple ER": {
			request: UpdatePolicyVersionRequest{
				UpdatePolicyVersion: UpdatePolicyVersion{
					Description: "Updated description",
				},
				PolicyID: 276858,
				Version:  5,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "activations": [],
    "createDate": 1629817335218,
    "createdBy": "agrebenk",
    "deleted": false,
    "description": "Updated description",
    "lastModifiedBy": "agrebenk",
    "lastModifiedDate": 1629821693867,
    "location": "/cloudlets/api/v2/policies/276858/versions/5",
    "matchRuleFormat": "1.0",
    "matchRules": null,
    "policyId": 276858,
    "revisionId": 4814876,
    "rulesLocked": false,
    "version": 5
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/5",
			expectedResponse: &PolicyVersion{
				CreateDate:       1629817335218,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "Updated description",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1629821693867,
				Location:         "/cloudlets/api/v2/policies/276858/versions/5",
				Activations:      []*Activation{},
				MatchRules:       nil,
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4814876,
				RulesLocked:      false,
				Version:          5,
			},
		},

		"201 updated complex ALB with warnings": {
			request: UpdatePolicyVersionRequest{
				UpdatePolicyVersion: UpdatePolicyVersion{
					Description: "Updated description",
					MatchRules: MatchRules{
						&MatchRuleALB{
							Matches: []MatchCriteriaALB{
								{
									MatchType:     "protocol",
									MatchValue:    "https",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
							},
							Start: 1,
							End:   2,
							Type:  "albMatchRule",
							Name:  "Rule3",
							ForwardSettings: ForwardSettings{
								OriginID: "alb_test_krk_dc1_only",
							},
							ID: 0,
						},
						&MatchRuleALB{
							Matches: []MatchCriteriaALB{
								{
									MatchType:     "protocol",
									MatchValue:    "http",
									MatchOperator: "equals",
									Negate:        false,
									CaseSensitive: false,
								},
							},
							Start: 0,
							End:   0,
							Type:  "albMatchRule",
							Name:  "Rule1",
							ForwardSettings: ForwardSettings{
								OriginID: "alb_test_krk_0_100",
							},
							ID: 0,
						},
					},
				},
				PolicyID: 276858,
				Version:  5,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "activations": [],
    "createDate": 1629981546401,
    "createdBy": "agrebenk",
    "deleted": false,
    "description": "Updated description",
    "lastModifiedBy": "agrebenk",
    "lastModifiedDate": 1630414029735,
    "location": "/cloudlets/api/v2/policies/279628/versions/2",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "akaRuleId": "80842d2cb25ccede",
            "end": 2,
            "forwardSettings": {
                "originId": "alb_test_krk_dc1_only"
            },
            "id": 10,
            "location": "/cloudlets/api/v2/policies/279628/versions/2/rules/80842d2cb25ccede",
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "protocol",
                    "matchValue": "https",
                    "negate": false
                }
            ],
            "name": "Rule3",
            "start": 1
        },
        {
            "type": "albMatchRule",
            "akaRuleId": "c018ee7c534b568c",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_0_100"
            },
            "id": 0,
            "location": "/cloudlets/api/v2/policies/279628/versions/2/rules/c018ee7c534b568c",
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "protocol",
                    "matchValue": "http",
                    "negate": false
                }
            ],
            "name": "Rule1",
            "start": 0
        }
    ],
    "policyId": 279628,
    "revisionId": 4815971,
    "rulesLocked": false,
    "version": 2,
    "warnings": [
        {
            "detail": "Start time is very old, possibly invalid: 1 (1970-01-01T00:00:01Z)",
            "title": "Invalid Result Value",
            "type": "/cloudlets/error-types/invalid-result-value",
            "jsonPointer": "/matchRules/0"
        },
        {
            "detail": "End time is very old, possibly invalid: 2 (1970-01-01T00:00:02Z)",
            "title": "Invalid Result Value",
            "type": "/cloudlets/error-types/invalid-result-value",
            "jsonPointer": "/matchRules/0"
        }
    ]
}
`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/5",
			expectedResponse: &PolicyVersion{
				Activations:      []*Activation{},
				CreateDate:       1629981546401,
				CreatedBy:        "agrebenk",
				Deleted:          false,
				Description:      "Updated description",
				LastModifiedBy:   "agrebenk",
				LastModifiedDate: 1630414029735,
				Location:         "/cloudlets/api/v2/policies/279628/versions/2",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type:      "albMatchRule",
						AkaRuleID: "80842d2cb25ccede",
						End:       2,
						ForwardSettings: ForwardSettings{
							OriginID: "alb_test_krk_dc1_only",
						},
						ID:       10,
						Location: "/cloudlets/api/v2/policies/279628/versions/2/rules/80842d2cb25ccede",
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "protocol",
								MatchValue:    "https",
								Negate:        false,
							},
						},
						Name:  "Rule3",
						Start: 1,
					},
					&MatchRuleALB{
						Type:      "albMatchRule",
						AkaRuleID: "c018ee7c534b568c",
						End:       0,
						ForwardSettings: ForwardSettings{
							OriginID: "alb_test_krk_0_100",
						},
						ID:       0,
						Location: "/cloudlets/api/v2/policies/279628/versions/2/rules/c018ee7c534b568c",
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "protocol",
								MatchValue:    "http",
								Negate:        false,
							},
						},
						Name:  "Rule1",
						Start: 0,
					},
				},
				PolicyID:    279628,
				RevisionID:  4815971,
				RulesLocked: false,
				Version:     2,
				Warnings: []Warning{
					{
						Detail:      "Start time is very old, possibly invalid: 1 (1970-01-01T00:00:01Z)",
						Title:       "Invalid Result Value",
						Type:        "/cloudlets/error-types/invalid-result-value",
						JSONPointer: "/matchRules/0",
					},
					{
						Detail:      "End time is very old, possibly invalid: 2 (1970-01-01T00:00:02Z)",
						Title:       "Invalid Result Value",
						Type:        "/cloudlets/error-types/invalid-result-value",
						JSONPointer: "/matchRules/0",
					},
				},
			},
		},
		"500 internal server error": {
			request: UpdatePolicyVersionRequest{
				PolicyID: 276858,
				Version:  3,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error creating enrollment",
  "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/3",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			expectedPath: "/cloudlets/api/v2/policies/0",
			request: UpdatePolicyVersionRequest{
				UpdatePolicyVersion: UpdatePolicyVersion{
					Description: strings.Repeat("A", 256),
				},
			},
			withError: ErrStructValidation,
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

func TestUnmarshalJSONMatchRules(t *testing.T) {
	tests := map[string]struct {
		withError      error
		responseBody   string
		expectedObject MatchRules
	}{
		"valid MarchRuleALB": {
			responseBody: `
	[
        {
            "type": "albMatchRule",
            "akaRuleId": "8a57dcbd5565cc9a",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_dc1_only"
            },
            "id": 0,
            "location": "/cloudlets/api/v2/policies/279628/versions/2/rules/8a57dcbd5565cc9a",
            "matchURL": null,
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "protocol",
                    "matchValue": "https",
                    "negate": false
                },
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "range",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "range",
                        "value": [
                            1,
                            50
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
                }
            ],
            "name": "Rule3",
            "start": 0
        }
    ]
`,
			expectedObject: MatchRules{
				&MatchRuleALB{
					Type:      "albMatchRule",
					AkaRuleID: "8a57dcbd5565cc9a",
					End:       0,
					ForwardSettings: ForwardSettings{
						OriginID: "alb_test_krk_dc1_only",
					},
					ID:       0,
					Location: "/cloudlets/api/v2/policies/279628/versions/2/rules/8a57dcbd5565cc9a",
					MatchURL: "",
					Matches: []MatchCriteriaALB{
						{
							CaseSensitive: false,
							MatchOperator: "equals",
							MatchType:     "protocol",
							MatchValue:    "https",
							Negate:        false,
						},
						{
							CaseSensitive: false,
							MatchOperator: "equals",
							MatchType:     "range",
							Negate:        false,
							ObjectMatchValue: &ObjectMatchValueRangeOrSimpleSubtype{
								Type: "range",
								Value: []interface{}{
									float64(1),
									float64(50),
								},
							},
						},
						{
							CaseSensitive: false,
							MatchOperator: "equals",
							MatchType:     "method",
							Negate:        false,
							ObjectMatchValue: &ObjectMatchValueRangeOrSimpleSubtype{
								Type: "simple",
								Value: []interface{}{
									"GET",
								},
							},
						},
					},
					Name:  "Rule3",
					Start: 0,
				},
			},
		},
		"invalid object value type": {
			withError: errors.New("ObjectMatchValue. UnmarshalJSON: unexpected type: foo"),
			responseBody: `
	[
        {
            "type": "albMatchRule",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "type": "foo",
                        "value": [
                            "GET"
                        ]
                    }
                }
            ],
            "name": "Rule3",
            "start": 0
        }
    ]
`,
		},
		"wrong type for object value type": {
			withError: errors.New("'type' should be a string"),
			responseBody: `
	[
        {
            "type": "albMatchRule",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "type": 1,
                        "value": [
                            "GET"
                        ]
                    }
                }
            ],
            "name": "Rule3",
            "start": 0
        }
    ]
`,
		},
		"missing object value type": {
			withError: errors.New("object should contain 'type' field"),
			responseBody: `
	[
        {
            "type": "albMatchRule",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": {
                        "value": [
                            "GET"
                        ]
                    }
                }
            ],
            "name": "Rule3",
            "start": 0
        }
    ]
`,
		},
		"invalid object value": {
			withError: errors.New("structure of ObjectMatchValue should be the map, but was string"),
			responseBody: `
	[
        {
            "type": "albMatchRule",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
                    "negate": false,
                    "objectMatchValue": ""
                }
            ],
            "name": "Rule3",
            "start": 0
        }
    ]
`,
		},
		"invalid MarchRuleAP": {
			responseBody: `
	[
        {
            "type": "apMatchRule"
        }
    ]
`,
			withError: errors.New("unsupported match rule type: apMatchRule"),
		},
		"invalid type": {
			withError: errors.New("'type' field on match rule entry should be a string"),
			responseBody: `
	[
        {
            "type": 1
        }
    ]
`,
		},
		"invalid JSON": {
			withError: errors.New("unexpected end of JSON input"),
			responseBody: `
	[
        {
            "type": "albMatchRule"
        }
    
`,
		},
		"missing type": {
			withError: errors.New("match rule entry should contain 'type' field"),
			responseBody: `
	[
        {
        }
    ]
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var matchRules MatchRules
			err := json.Unmarshal([]byte(test.responseBody), &matchRules)

			if test.withError != nil {
				assert.Equal(t, test.withError.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedObject, matchRules)
		})
	}
}
