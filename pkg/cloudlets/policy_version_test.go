package cloudlets

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets/tools"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestListPolicyVersions(t *testing.T) {
	tests := map[string]struct {
		request          ListPolicyVersionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []PolicyVersion
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: ListPolicyVersionsRequest{
				PolicyID: 284823,
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "activations": [],
        "createDate": 1631191583350,
        "createdBy": "jsmith",
        "deleted": false,
        "description": "some string",
        "lastModifiedBy": "jsmith",
        "lastModifiedDate": 1631191583350,
        "location": "/cloudlets/api/v2/policies/284823/versions/2",
        "matchRuleFormat": "1.0",
        "matchRules": null,
        "policyId": 284823,
        "revisionId": 4824132,
        "rulesLocked": false,
        "version": 2
    },
    {
        "activations": [],
        "createDate": 1631190136935,
        "createdBy": "jsmith",
        "deleted": false,
        "description": null,
        "lastModifiedBy": "jsmith",
        "lastModifiedDate": 1631190136935,
        "location": "/cloudlets/api/v2/policies/284823/versions/1",
        "matchRuleFormat": "1.0",
        "matchRules": null,
        "policyId": 284823,
        "revisionId": 4824129,
        "rulesLocked": false,
        "version": 1
    }
]`,
			expectedPath: "/cloudlets/api/v2/policies/284823/versions?includeActivations=false&includeDeleted=false&includeRules=false&offset=0",
			expectedResponse: []PolicyVersion{
				{
					Activations:      []PolicyActivation{},
					CreateDate:       1631191583350,
					CreatedBy:        "jsmith",
					Deleted:          false,
					Description:      "some string",
					LastModifiedBy:   "jsmith",
					LastModifiedDate: 1631191583350,
					Location:         "/cloudlets/api/v2/policies/284823/versions/2",
					MatchRuleFormat:  "1.0",
					MatchRules:       nil,
					PolicyID:         284823,
					RevisionID:       4824132,
					RulesLocked:      false,
					Version:          2,
				},
				{
					Activations:      []PolicyActivation{},
					CreateDate:       1631190136935,
					CreatedBy:        "jsmith",
					Deleted:          false,
					Description:      "",
					LastModifiedBy:   "jsmith",
					LastModifiedDate: 1631190136935,
					Location:         "/cloudlets/api/v2/policies/284823/versions/1",
					MatchRuleFormat:  "1.0",
					MatchRules:       nil,
					PolicyID:         284823,
					RevisionID:       4824129,
					RulesLocked:      false,
					Version:          1,
				},
			},
		},

		"200 OK with params": {
			request: ListPolicyVersionsRequest{
				PolicyID:           284823,
				IncludeRules:       true,
				IncludeDeleted:     true,
				IncludeActivations: true,
				Offset:             2,
				PageSize:           tools.IntPtr(3),
			},
			responseStatus: http.StatusOK,
			responseBody: `
[]`,
			expectedPath:     "/cloudlets/api/v2/policies/284823/versions?includeActivations=true&includeDeleted=true&includeRules=true&offset=2&pageSize=3",
			expectedResponse: []PolicyVersion{},
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
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/284823/versions?includeActivations=false&includeDeleted=false&includeRules=false&offset=0",
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
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
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
				Activations:      []PolicyActivation{},
				CreateDate:       1629299944291,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
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
		"200 OK, ER with disabled rule": {
			request: GetPolicyVersionRequest{
				PolicyID: 276858,
				Version:  6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "erMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "disabled": true,
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/6?omitRules=false",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
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
						Disabled:                 true,
					},
				},
			},
		},

		"200 OK, ALB with disabled rule": {
			request: GetPolicyVersionRequest{
				PolicyID: 276858,
				Version:  6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "disabled": true,
            "name": "rul3",
            "start": 0
        }
    ],
    "policyId": 276858,
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/6?omitRules=false",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type:     "albMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						Disabled: true,
					},
				},
			},
		},
		"200 OK, CD rule": {
			request: GetPolicyVersionRequest{
				PolicyID: 276858,
				Version:  6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "activations": [],
    "createDate": 1638547203265,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1638547203265,
    "location": "/cloudlets/api/v2/policies/325401/versions/3",
    "matchRuleFormat": "1.0",
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
            "location": "/cloudlets/api/v2/policies/325401/versions/3/rules/b151ca68e51f5a61",
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
    "revisionId": 4888857,
    "rulesLocked": false,
    "version": 3
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/6?omitRules=false",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1638547203265,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1638547203265,
				Location:         "/cloudlets/api/v2/policies/325401/versions/3",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleCD{
						Type: "cdMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsCD{
							OriginID: "fr_test_krk_dc2",
							Percent:  11,
						},
						ID: 0,
						Matches: []MatchCriteriaCD{
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
				PolicyID:    325401,
				RevisionID:  4888857,
				RulesLocked: false,
				Version:     3,
			},
		},
		"200 OK, ER rule with disabled=false": {
			request: GetPolicyVersionRequest{
				PolicyID: 276858,
				Version:  6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/6?omitRules=false",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
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
		"200 OK, ALB rule with disabled=false": {
			request: GetPolicyVersionRequest{
				PolicyID: 276858,
				Version:  6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "name": "rul3",
            "start": 0
        }
    ],
    "policyId": 276858,
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/6?omitRules=false",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type:     "albMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						Disabled: false,
					},
				},
			},
		},
		"200 OK, FR with disabled rule": {
			request: GetPolicyVersionRequest{
				PolicyID: 276858,
				Version:  6,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions/6?omitRules=false",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
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
					Description: "Description for the policy",
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "activations": [],
    "createDate": 1629812554924,
    "createdBy": "jsmith",
    "deleted": false,
    "description": "Description for the policy",
    "lastModifiedBy": "jsmith",
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
				Activations:      []PolicyActivation{},
				CreateDate:       1629812554924,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "Description for the policy",
				LastModifiedBy:   "jsmith",
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
							ForwardSettings: ForwardSettingsALB{
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
							ForwardSettings: ForwardSettingsALB{
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
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981546401,
    "location": "/cloudlets/api/v2/policies/279628/versions/2",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_dc1_only"
            },
            "id": 0,
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
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_0_100"
            },
            "id": 0,
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
				Activations:      []PolicyActivation{},
				CreateDate:       1629981546401,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981546401,
				Location:         "/cloudlets/api/v2/policies/279628/versions/2",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type: "albMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsALB{
							OriginID: "alb_test_krk_dc1_only",
						},
						ID:       0,
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
						Type: "albMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsALB{
							OriginID: "alb_test_krk_0_100",
						},
						ID:       0,
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
									ObjectMatchValue: &ObjectMatchValueObject{
										Type: "object",
										Name: "ALB",
										Options: &Options{
											Value:            []string{"y"},
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
							ForwardSettings: ForwardSettingsALB{
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
		    "createdBy": "jsmith",
		    "deleted": false,
		    "description": "New version 1630480693371",
		    "lastModifiedBy": "jsmith",
		    "lastModifiedDate": 1630489884742,
		    "location": "/cloudlets/api/v2/policies/139743/versions/796",
		    "matchRuleFormat": "1.0",
		    "matchRules": [
		        {
		            "type": "albMatchRule",
		            "end": 0,
		            "forwardSettings": {
		                "originId": "alb_test_krk_mutable"
		            },
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
				Activations:      []PolicyActivation{},
				CreateDate:       1630489884742,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "New version 1630480693371",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1630489884742,
				Location:         "/cloudlets/api/v2/policies/139743/versions/796",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type: "albMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsALB{
							OriginID: "alb_test_krk_mutable",
						},
						ID:       0,
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "header",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObject{
									Type: "object",
									Name: "ALB",
									Options: &Options{
										Value:            []string{"y"},
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
									ObjectMatchValue: &ObjectMatchValueSimple{
										Type:  "simple",
										Value: []string{"GET"},
									},
								},
							},
							Start: 0,
							End:   0,
							Type:  "albMatchRule",
							Name:  "alb rule",
							ForwardSettings: ForwardSettingsALB{
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
    "createdBy": "jsmith",
    "deleted": false,
    "description": "New version 1630480693371",
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1630506044027,
    "location": "/cloudlets/api/v2/policies/139743/versions/797",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_mutable"
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
				Activations:      []PolicyActivation{},
				CreateDate:       1630506044027,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "New version 1630480693371",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1630506044027,
				Location:         "/cloudlets/api/v2/policies/139743/versions/797",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type: "albMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsALB{
							OriginID: "alb_test_krk_mutable",
						},
						ID:       0,
						MatchURL: "",
						Matches: []MatchCriteriaALB{
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
									ObjectMatchValue: &ObjectMatchValueRange{
										Type:  "range",
										Value: []int64{1, 50},
									},
								},
							},
							Start: 0,
							End:   0,
							Type:  "albMatchRule",
							Name:  "alb rule",
							ForwardSettings: ForwardSettingsALB{
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
    "createdBy": "jsmith",
    "deleted": false,
    "description": "New version 1630480693371",
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1630507099511,
    "location": "/cloudlets/api/v2/policies/139743/versions/798",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_mutable"
            },
            "id": 0,
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
				Activations:      []PolicyActivation{},
				CreateDate:       1630507099511,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "New version 1630480693371",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1630507099511,
				Location:         "/cloudlets/api/v2/policies/139743/versions/798",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type: "albMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsALB{
							OriginID: "alb_test_krk_mutable",
						},
						ID:       0,
						MatchURL: "",
						Matches: []MatchCriteriaALB{
							{
								CaseSensitive: false,
								MatchOperator: "equals",
								MatchType:     "range",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueRange{
									Type:  "range",
									Value: []int64{1, 50},
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

		"201 created, complex CD": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleCD{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsCD{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaCD{
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
						&MatchRuleCD{
							Start:    0,
							End:      0,
							Type:     "cdMatchRule",
							Name:     "rule 2",
							MatchURL: "ddd.aaa",
							ID:       0,
							ForwardSettings: ForwardSettingsCD{
								OriginID: "some_origin",
								Percent:  10,
							},
						},
						&MatchRuleCD{
							Type:     "cdMatchRule",
							ID:       0,
							Name:     "r1",
							Start:    0,
							End:      0,
							MatchURL: "abc.com",
							ForwardSettings: ForwardSettingsCD{
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleCD{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						ForwardSettings: ForwardSettingsCD{
							OriginID: "some_origin",
							Percent:  10,
						},
						Matches: []MatchCriteriaCD{
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
					&MatchRuleCD{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "ddd.aaa",
						Name:     "rule 2",
						Start:    0,
					},
					&MatchRuleCD{
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
		"201 created, complex CD with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleCD{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsCD{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaCD{
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleCD{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						ForwardSettings: ForwardSettingsCD{
							OriginID: "some_origin",
							Percent:  10,
						},
						Matches: []MatchCriteriaCD{
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
		"201 created, complex CD with objectMatchValue - object": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleCD{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsCD{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaCD{
								{
									MatchOperator: "equals",
									MatchType:     "header",
									Negate:        false,
									ObjectMatchValue: &ObjectMatchValueObject{
										Type: "object",
										Name: "CD",
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
						"name": "CD",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleCD{
						Type:     "cdMatchRule",
						End:      0,
						ID:       0,
						MatchURL: "",
						Name:     "rul3",
						Start:    0,
						ForwardSettings: ForwardSettingsCD{
							OriginID: "some_origin",
							Percent:  10,
						},
						Matches: []MatchCriteriaCD{
							{
								MatchOperator: "equals",
								MatchType:     "hostname",
								Negate:        false,
								ObjectMatchValue: &ObjectMatchValueObject{
									Type: "object",
									Name: "CD",
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
		"validation error, complex CD with unavailable objectMatchValue type - range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleCD{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsCD{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaCD{
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
		"validation error, complex CD missing forwardSettings": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleCD{
							Start: 0,
							End:   0,
							Type:  "cdMatchRule",
							Name:  "rul3",
							ID:    0,
							ForwardSettings: ForwardSettingsCD{
								OriginID: "some_origin",
								Percent:  10,
							},
							Matches: []MatchCriteriaCD{
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
						&MatchRuleCD{
							Start:    0,
							End:      0,
							Type:     "cdMatchRule",
							Name:     "rule 2",
							MatchURL: "ddd.aaa",
							ID:       0,
						},
						&MatchRuleCD{
							Type:     "cdMatchRule",
							ID:       0,
							Name:     "r1",
							Start:    0,
							End:      0,
							MatchURL: "abc.com",
							ForwardSettings: ForwardSettingsCD{
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
            "location": "/cloudlets/api/v2/policies/276858/versions/1/rules/893947a3d5a85c1b",
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
            "location": "/cloudlets/api/v2/policies/276858/versions/1/rules/aa379d230efcded0",
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
            "location": "/cloudlets/api/v2/policies/276858/versions/1/rules/1afe03d843996766",
            "matchURL": "abc.com",
            "name": "rule 2",
            "start": 0
        }
    ],
    "policyId": 276858,
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
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
					Description:     "New version 1630480693371",
					MatchRuleFormat: "1.0",
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
    "activations": [],
    "createDate": 1630507099511,
    "createdBy": "jsmith",
    "deleted": false,
    "description": "New version 1630480693371",
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1630507099511,
    "location": "/cloudlets/api/v2/policies/139743/versions/798",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4819450,
    "rulesLocked": false,
    "version": 798
}`,
			expectedPath: "/cloudlets/api/v2/policies/139743/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1630507099511,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "New version 1630480693371",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1630507099511,
				Location:         "/cloudlets/api/v2/policies/139743/versions/798",
				MatchRuleFormat:  "1.0",
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
				PolicyID:    139743,
				RevisionID:  4819450,
				RulesLocked: false,
				Version:     798,
			},
		},
		"201 created, complex FR with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					Description:     "New version 1630480693371",
					MatchRuleFormat: "1.0",
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
    "activations": [],
    "createDate": 1630507099511,
    "createdBy": "jsmith",
    "deleted": false,
    "description": "New version 1630480693371",
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1630507099511,
    "location": "/cloudlets/api/v2/policies/139743/versions/798",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4819450,
    "rulesLocked": false,
    "version": 798
}`,
			expectedPath: "/cloudlets/api/v2/policies/139743/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1630507099511,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "New version 1630480693371",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1630507099511,
				Location:         "/cloudlets/api/v2/policies/139743/versions/798",
				MatchRuleFormat:  "1.0",
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
				PolicyID:    139743,
				RevisionID:  4819450,
				RulesLocked: false,
				Version:     798,
			},
		},

		"201 created, complex VP with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleVP{
							Start:              0,
							End:                0,
							Type:               "vpMatchRule",
							Name:               "rul3",
							PassThroughPercent: tools.Float64Ptr(-1),
							ID:                 0,
							Matches: []MatchCriteriaVP{
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "vpMatchRule",
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
			"passThroughPercent": -1
        }
    ],
    "policyId": 276858,
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleVP{
						Type:               "vpMatchRule",
						End:                0,
						ID:                 0,
						MatchURL:           "",
						Name:               "rul3",
						PassThroughPercent: tools.Float64Ptr(-1),
						Start:              0,
						Matches: []MatchCriteriaVP{
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

		"201 created, complex AP with objectMatchValue - simple": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleAP{
							Start:              0,
							End:                0,
							Type:               "apMatchRule",
							Name:               "rul3",
							PassThroughPercent: tools.Float64Ptr(0),
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleAP{
						Type:               "apMatchRule",
						End:                0,
						ID:                 0,
						MatchURL:           "",
						Name:               "rul3",
						PassThroughPercent: tools.Float64Ptr(-1),
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
							PassThroughPercent: tools.Float64Ptr(-1),
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
    "activations": [],
    "createDate": 1629981355165,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629981355165,
    "location": "/cloudlets/api/v2/policies/276858/versions/6",
    "matchRuleFormat": "1.0",
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
    "revisionId": 4815968,
    "rulesLocked": false,
    "version": 6
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858/versions",
			expectedResponse: &PolicyVersion{
				Activations:      []PolicyActivation{},
				CreateDate:       1629981355165,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629981355165,
				Location:         "/cloudlets/api/v2/policies/276858/versions/6",
				MatchRuleFormat:  "1.0",
				PolicyID:         276858,
				RevisionID:       4815968,
				RulesLocked:      false,
				Version:          6,
				MatchRules: MatchRules{
					&MatchRuleAP{
						Type:               "apMatchRule",
						End:                0,
						ID:                 0,
						MatchURL:           "",
						Name:               "rul3",
						PassThroughPercent: tools.Float64Ptr(-1),
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

		"validation error, complex VP with unavailable objectMatchValue type - range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleVP{
							Start:              0,
							End:                0,
							Type:               "vpMatchRule",
							PassThroughPercent: tools.Float64Ptr(50.50),
							Name:               "rul3",
							ID:                 0,
							Matches: []MatchCriteriaVP{
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
							PassThroughPercent: tools.Float64Ptr(50.50),
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

		"validation error, simple VP missing passThroughPercent": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleVP{
							Start: 0,
							End:   0,
							Type:  "vpMatchRule",
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

		"validation error, simple VP passThroughPercent out of range": {
			request: CreatePolicyVersionRequest{
				CreatePolicyVersion: CreatePolicyVersion{
					MatchRules: MatchRules{
						&MatchRuleVP{
							Start:              0,
							End:                0,
							Type:               "vpMatchRule",
							PassThroughPercent: tools.Float64Ptr(101),
							Name:               "rul3",
							ID:                 0,
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
							PassThroughPercent: tools.Float64Ptr(101),
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
    "createdBy": "jsmith",
    "deleted": false,
    "description": "Updated description",
    "lastModifiedBy": "jsmith",
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
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "Updated description",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629821693867,
				Location:         "/cloudlets/api/v2/policies/276858/versions/5",
				Activations:      []PolicyActivation{},
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
							ForwardSettings: ForwardSettingsALB{
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
							ForwardSettings: ForwardSettingsALB{
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
    "createdBy": "jsmith",
    "deleted": false,
    "description": "Updated description",
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1630414029735,
    "location": "/cloudlets/api/v2/policies/279628/versions/2",
    "matchRuleFormat": "1.0",
    "matchRules": [
        {
            "type": "albMatchRule",
            "end": 2,
            "forwardSettings": {
                "originId": "alb_test_krk_dc1_only"
            },
            "id": 10,
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
            "end": 0,
            "forwardSettings": {
                "originId": "alb_test_krk_0_100"
            },
            "id": 0,
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
				Activations:      []PolicyActivation{},
				CreateDate:       1629981546401,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "Updated description",
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1630414029735,
				Location:         "/cloudlets/api/v2/policies/279628/versions/2",
				MatchRuleFormat:  "1.0",
				MatchRules: MatchRules{
					&MatchRuleALB{
						Type: "albMatchRule",
						End:  2,
						ForwardSettings: ForwardSettingsALB{
							OriginID: "alb_test_krk_dc1_only",
						},
						ID:       10,
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
						Type: "albMatchRule",
						End:  0,
						ForwardSettings: ForwardSettingsALB{
							OriginID: "alb_test_krk_0_100",
						},
						ID:       0,
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
