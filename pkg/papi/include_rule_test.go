package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIncludeRuleTree(t *testing.T) {
	tests := map[string]struct {
		params           GetIncludeRuleTreeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetIncludeRuleTreeResponse
		withError        error
	}{
		"200 OK - get include": {
			params: GetIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				ValidateMode:   "fast",
				ValidateRules:  false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "test_account",
    "contractId": "test_contract",
    "groupId": "test_group",
    "includeId": "inc_123456",
    "includeVersion": 2,
	"includeName": "test_include",
	"includeType": "MICROSERVICES",
    "etag": "etag",
    "ruleFormat": "v2020-11-02",
    "rules": {
        "name": "default",
        "criteria": [],
        "children": [
            {
                "name": "Compress Text Content",
                "criteria": [
                    {
                        "name": "contentType",
                        "options": {
                            "matchOperator": "IS_ONE_OF",
                            "matchWildcard": true,
                            "matchCaseSensitive": false,
                            "values": [
                                "text/html*",
                                "text/css*",
                                "application/x-javascript*"
                            ]
                        }
                    }
                ],
                "behaviors": [
                    {
                        "name": "gzipResponse",
                        "options": { "behavior": "ALWAYS" }
                    }
                ]
            }
        ],
        "options": {
            "is_secure": false
        },
        "behaviors": [
            {
                "name": "origin",
                "options": {
                    "httpPort": 80,
                    "enableTrueClientIp": false,
                    "compress": true,
                    "cacheKeyHostname": "ORIGIN_HOSTNAME",
                    "forwardHostHeader": "REQUEST_HOST_HEADER",
                    "hostname": "origin.test.com",
                    "originType": "CUSTOMER"
                }
            },
            {
                "name": "cpCode",
                "options": {
                    "value": {
                        "id": 12345,
                        "name": "my CP code"
                    }
                }
            }
        ],
 		"customOverride": {
        	"overrideId": "cbo_12345",
        	"name": "mdc"
    	},
		"variables": [
            {
                "name": "VAR_NAME",
                "value": "default value",
                "description": "This is a sample Property Manager variable.",
                "hidden": false,
                "sensitive": false
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/includes/inc_123456/versions/2/rules?contractId=test_contract&groupId=test_group&validateMode=fast&validateRules=false",
			expectedResponse: &GetIncludeRuleTreeResponse{
				Response: Response{
					AccountID:  "test_account",
					ContractID: "test_contract",
					GroupID:    "test_group",
				},
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				IncludeName:    "test_include",
				IncludeType:    IncludeTypeMicroServices,
				Etag:           "etag",
				RuleFormat:     "v2020-11-02",
				Rules: Rules{
					Behaviors: []RuleBehavior{
						{
							Name: "origin",
							Options: RuleOptionsMap{
								"httpPort":           float64(80),
								"enableTrueClientIp": false,
								"compress":           true,
								"cacheKeyHostname":   "ORIGIN_HOSTNAME",
								"forwardHostHeader":  "REQUEST_HOST_HEADER",
								"hostname":           "origin.test.com",
								"originType":         "CUSTOMER",
							},
						},
						{
							Name: "cpCode",
							Options: RuleOptionsMap{
								"value": map[string]interface{}{
									"id":   float64(12345),
									"name": "my CP code",
								},
							},
						},
					},
					Children: []Rules{
						{
							Behaviors: []RuleBehavior{
								{
									Name: "gzipResponse",
									Options: RuleOptionsMap{
										"behavior": "ALWAYS",
									},
								},
							},
							Criteria: []RuleBehavior{
								{
									Locked: false,
									Name:   "contentType",
									Options: RuleOptionsMap{
										"matchOperator":      "IS_ONE_OF",
										"matchWildcard":      true,
										"matchCaseSensitive": false,
										"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
									},
								},
							},
							Name: "Compress Text Content",
						},
					},
					Criteria: []RuleBehavior{},
					Name:     "default",
					Options:  RuleOptions{IsSecure: false},
					CustomOverride: &RuleCustomOverride{
						OverrideID: "cbo_12345",
						Name:       "mdc",
					},
					Variables: []RuleVariable{
						{
							Description: "This is a sample Property Manager variable.",
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       "default value",
						},
					},
				},
			},
		},
		"200 OK - get include with ruleFormat": {
			params: GetIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				RuleFormat:     "v2020-11-02",
				ValidateMode:   "fast",
				ValidateRules:  false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "test_account",
    "contractId": "test_contract",
    "groupId": "test_group",
    "includeId": "inc_123456",
    "includeVersion": 2,
	"includeName": "test_include",
	"includeType": "MICROSERVICES",
    "etag": "etag",
    "ruleFormat": "v2020-11-02",
    "rules": {
        "name": "default",
        "criteria": [],
        "children": [
            {
                "name": "Compress Text Content",
                "criteria": [
                    {
                        "name": "contentType",
                        "options": {
                            "matchOperator": "IS_ONE_OF",
                            "matchWildcard": true,
                            "matchCaseSensitive": false,
                            "values": [
                                "text/html*",
                                "text/css*",
                                "application/x-javascript*"
                            ]
                        }
                    }
                ],
                "behaviors": [
                    {
                        "name": "gzipResponse",
                        "options": { "behavior": "ALWAYS" }
                    }
                ]
            }
        ],
        "options": {
            "is_secure": false
        },
        "behaviors": [
            {
                "name": "origin",
                "options": {
                    "httpPort": 80,
                    "enableTrueClientIp": false,
                    "compress": true,
                    "cacheKeyHostname": "ORIGIN_HOSTNAME",
                    "forwardHostHeader": "REQUEST_HOST_HEADER",
                    "hostname": "origin.test.com",
                    "originType": "CUSTOMER"
                }
            },
            {
                "name": "cpCode",
                "options": {
                    "value": {
                        "id": 12345,
                        "name": "my CP code"
                    }
                }
            }
        ],
 		"customOverride": {
        	"overrideId": "cbo_12345",
        	"name": "mdc"
    	},
		"variables": [
            {
                "name": "VAR_NAME",
                "value": "default value",
                "description": "This is a sample Property Manager variable.",
                "hidden": false,
                "sensitive": false
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/includes/inc_123456/versions/2/rules?contractId=test_contract&groupId=test_group&validateMode=fast&validateRules=false",
			expectedResponse: &GetIncludeRuleTreeResponse{
				Response: Response{
					AccountID:  "test_account",
					ContractID: "test_contract",
					GroupID:    "test_group",
				},
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				IncludeName:    "test_include",
				IncludeType:    IncludeTypeMicroServices,
				Etag:           "etag",
				RuleFormat:     "v2020-11-02",
				Rules: Rules{
					Behaviors: []RuleBehavior{
						{
							Name: "origin",
							Options: RuleOptionsMap{
								"httpPort":           float64(80),
								"enableTrueClientIp": false,
								"compress":           true,
								"cacheKeyHostname":   "ORIGIN_HOSTNAME",
								"forwardHostHeader":  "REQUEST_HOST_HEADER",
								"hostname":           "origin.test.com",
								"originType":         "CUSTOMER",
							},
						},
						{
							Name: "cpCode",
							Options: RuleOptionsMap{
								"value": map[string]interface{}{
									"id":   float64(12345),
									"name": "my CP code",
								},
							},
						},
					},
					Children: []Rules{
						{
							Behaviors: []RuleBehavior{
								{
									Name: "gzipResponse",
									Options: RuleOptionsMap{
										"behavior": "ALWAYS",
									},
								},
							},
							Criteria: []RuleBehavior{
								{
									Locked: false,
									Name:   "contentType",
									Options: RuleOptionsMap{
										"matchOperator":      "IS_ONE_OF",
										"matchWildcard":      true,
										"matchCaseSensitive": false,
										"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
									},
								},
							},
							Name: "Compress Text Content",
						},
					},
					Criteria: []RuleBehavior{},
					Name:     "default",
					Options:  RuleOptions{IsSecure: false},
					CustomOverride: &RuleCustomOverride{
						OverrideID: "cbo_12345",
						Name:       "mdc",
					},
					Variables: []RuleVariable{
						{
							Description: "This is a sample Property Manager variable.",
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       "default value",
						},
					},
				},
			},
		},
		"500 Internal Server Error": {
			params: GetIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				ValidateMode:   "fast",
				ValidateRules:  false,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching rule tree",
    "status": 500
}`,
			expectedPath: "/papi/v1/includes/inc_123456/versions/2/rules?contractId=test_contract&groupId=test_group&validateMode=fast&validateRules=false",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching rule tree",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: GetIncludeRuleTreeRequest{
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				ValidateMode:   "fast",
				ValidateRules:  false,
			},
			withError: ErrStructValidation,
		},
		"validation error - missing includeVersion": {
			params: GetIncludeRuleTreeRequest{
				IncludeID:     "inc_123456",
				ContractID:    "test_contract",
				GroupID:       "test_group",
				ValidateMode:  "fast",
				ValidateRules: false,
			},
			withError: ErrStructValidation,
		},
		"validation error - missing contractId": {
			params: GetIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				GroupID:        "test_group",
				ValidateMode:   "fast",
				ValidateRules:  false,
			},
			withError: ErrStructValidation,
		},
		"validation error - missing groupId": {
			params: GetIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				ValidateMode:   "fast",
				ValidateRules:  false,
			},
			withError: ErrStructValidation,
		},
		"validation error - invalid validation mode": {
			params: GetIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				ValidateMode:   "test",
				ValidateRules:  false,
			},
			withError: ErrStructValidation,
		},
		"validation error - invalid ruleFormat": {
			params: GetIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				RuleFormat:     "invalid",
				ValidateMode:   "fast",
				ValidateRules:  false,
			},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				if test.params.RuleFormat != "" {
					assert.Equal(t, r.Header.Get("Accept"), fmt.Sprintf("application/vnd.akamai.papirules.%s+json", test.params.RuleFormat))
				}
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetIncludeRuleTree(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateIncludeRuleTree(t *testing.T) {
	tests := map[string]struct {
		params           UpdateIncludeRuleTreeRequest
		responseStatus   int
		responseBody     string
		responseHeaders  map[string]string
		expectedPath     string
		expectedResponse *UpdateIncludeRuleTreeResponse
		withError        error
	}{
		"200 OK - update rule tree": {
			params: UpdateIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				DryRun:         true,
				ValidateMode:   "fast",
				ValidateRules:  false,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Comments: "default comment",
						Behaviors: []RuleBehavior{
							{
								Name: "origin",
								Options: RuleOptionsMap{
									"httpPort":           float64(80),
									"enableTrueClientIp": false,
									"compress":           true,
									"cacheKeyHostname":   "ORIGIN_HOSTNAME",
									"forwardHostHeader":  "REQUEST_HOST_HEADER",
									"hostname":           "origin.test.com",
									"originType":         "CUSTOMER",
								},
							},
							{
								Name: "cpCode",
								Options: RuleOptionsMap{
									"value": map[string]interface{}{
										"id":   float64(12345),
										"name": "my CP code",
									},
								},
							},
						},
						Children: []Rules{
							{
								Behaviors: []RuleBehavior{
									{
										Name: "gzipResponse",
										Options: RuleOptionsMap{
											"behavior": "ALWAYS",
										},
									},
								},
								Criteria: []RuleBehavior{
									{
										Locked: false,
										Name:   "contentType",
										Options: RuleOptionsMap{
											"matchOperator":      "IS_ONE_OF",
											"matchWildcard":      true,
											"matchCaseSensitive": false,
											"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
										},
									},
								},
								Name: "Compress Text Content",
							},
						},
						Criteria: []RuleBehavior{},
						Name:     "default",
						Options:  RuleOptions{IsSecure: false},
						CustomOverride: &RuleCustomOverride{
							OverrideID: "cbo_12345",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: "This is a sample Property Manager variable.",
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       "default value",
							},
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseHeaders: map[string]string{
				"x-limit-elements-per-property-limit":            "3000",
				"x-limit-elements-per-property-remaining":        "3000",
				"x-limit-max-nested-rules-per-include-limit":     "4",
				"x-limit-max-nested-rules-per-include-remaining": "4",
			},
			responseBody: `
{
    "accountId": "test_account",
    "contractId": "test_contract",
    "groupId": "test_group",
    "includeId": "inc_123456",
    "includeVersion": 2,
   	"includeName": "test_include",
    "includeType": "MICROSERVICES",
    "etag": "etag",
    "ruleFormat": "v2020-11-02",
	"comments": "version comment",
    "rules": {
        "name": "default",
        "comments": "default comment",
        "criteria": [],
        "children": [
            {
                "name": "Compress Text Content",
                "criteria": [
                    {
                        "name": "contentType",
                        "options": {
                            "matchOperator": "IS_ONE_OF",
                            "matchWildcard": true,
                            "matchCaseSensitive": false,
                            "values": [
                                "text/html*",
                                "text/css*",
                                "application/x-javascript*"
                            ]
                        }
                    }
                ],
                "behaviors": [
                    {
                        "name": "gzipResponse",
                        "options": { "behavior": "ALWAYS" }
                    }
                ]
            }
        ],
        "options": {
            "is_secure": false
        },
        "behaviors": [
            {
                "name": "origin",
                "options": {
                    "httpPort": 80,
                    "enableTrueClientIp": false,
                    "compress": true,
                    "cacheKeyHostname": "ORIGIN_HOSTNAME",
                    "forwardHostHeader": "REQUEST_HOST_HEADER",
                    "hostname": "origin.test.com",
                    "originType": "CUSTOMER"
                }
            },
            {
                "name": "cpCode",
                "options": {
                    "value": {
                        "id": 12345,
                        "name": "my CP code"
                    }
                }
            },
 			{
                "name": "matchAdvanced",
                "options": {
                    "value": {
                        "id": 12345,
                        "name": "my advanced match"
                    }
                },
				"locked": true,
				"uuid": "fd6a63bc-120a-4891-a5f2-c479765d5553",
				"templateUuid": "bedbac99-4ce1-43a3-96cc-b84c8cd30176"
            }
        ],
 		"customOverride": {
        	"overrideId": "cbo_12345",
        	"name": "mdc"
    	},
		"templateUuid": "bedbac99-4ce1-43a3-96cc-b84c8cd30176",
		"templateLink": "/platformtoolkit/service/ruletemplate/30582260/1?accountId=1-1TJZFB&gid=61726&ck=16.3.1.1",
		"variables": [
            {
                "name": "VAR_NAME",
                "value": "default value",
                "description": "This is a sample Property Manager variable.",
                "hidden": false,
                "sensitive": false
            }
        ]
    },
	"errors": [
        {
            "type": "https://problems.example.net/papi/v0/validation/attribute_required",
            "errorLocation": "#/rules/behaviors/0/options/reason",
            "detail": "The Reason ID option on the Control Access behavior is required."
        }
    ]
}`,
			expectedPath: "/papi/v1/includes/inc_123456/versions/2/rules?contractId=test_contract&dryRun=true&groupId=test_group&validateMode=fast&validateRules=false",
			expectedResponse: &UpdateIncludeRuleTreeResponse{
				Response: Response{
					AccountID:  "test_account",
					ContractID: "test_contract",
					GroupID:    "test_group",
					Errors: []*Error{
						{
							Type:          "https://problems.example.net/papi/v0/validation/attribute_required",
							ErrorLocation: "#/rules/behaviors/0/options/reason",
							Detail:        "The Reason ID option on the Control Access behavior is required.",
						},
					},
				},
				ResponseHeaders: UpdateIncludeResponseHeaders{
					ElementsPerPropertyRemaining:      "3000",
					ElementsPerPropertyTotal:          "3000",
					MaxNestedRulesPerIncludeRemaining: "4",
					MaxNestedRulesPerIncludeTotal:     "4",
				},
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				IncludeName:    "test_include",
				IncludeType:    IncludeTypeMicroServices,
				Etag:           "etag",
				RuleFormat:     "v2020-11-02",
				Comments:       "version comment",
				Rules: Rules{
					Comments: "default comment",
					Behaviors: []RuleBehavior{
						{
							Name: "origin",
							Options: RuleOptionsMap{
								"httpPort":           float64(80),
								"enableTrueClientIp": false,
								"compress":           true,
								"cacheKeyHostname":   "ORIGIN_HOSTNAME",
								"forwardHostHeader":  "REQUEST_HOST_HEADER",
								"hostname":           "origin.test.com",
								"originType":         "CUSTOMER",
							},
						},
						{
							Name: "cpCode",
							Options: RuleOptionsMap{
								"value": map[string]interface{}{
									"id":   float64(12345),
									"name": "my CP code",
								},
							},
						},
						{
							Name: "matchAdvanced",
							Options: RuleOptionsMap{
								"value": map[string]interface{}{
									"id":   float64(12345),
									"name": "my advanced match",
								},
							},
							Locked:       true,
							UUID:         "fd6a63bc-120a-4891-a5f2-c479765d5553",
							TemplateUuid: "bedbac99-4ce1-43a3-96cc-b84c8cd30176",
						},
					},
					Children: []Rules{
						{
							Behaviors: []RuleBehavior{
								{
									Name: "gzipResponse",
									Options: RuleOptionsMap{
										"behavior": "ALWAYS",
									},
								},
							},
							Criteria: []RuleBehavior{
								{
									Locked: false,
									Name:   "contentType",
									Options: RuleOptionsMap{
										"matchOperator":      "IS_ONE_OF",
										"matchWildcard":      true,
										"matchCaseSensitive": false,
										"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
									},
								},
							},
							Name: "Compress Text Content",
						},
					},
					Criteria: []RuleBehavior{},
					Name:     "default",
					Options:  RuleOptions{IsSecure: false},
					CustomOverride: &RuleCustomOverride{
						OverrideID: "cbo_12345",
						Name:       "mdc",
					},
					TemplateUuid: "bedbac99-4ce1-43a3-96cc-b84c8cd30176",
					TemplateLink: "/platformtoolkit/service/ruletemplate/30582260/1?accountId=1-1TJZFB&gid=61726&ck=16.3.1.1",
					Variables: []RuleVariable{
						{
							Description: "This is a sample Property Manager variable.",
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       "default value",
						},
					},
				},
			},
		},
		"500 Internal Server Error": {
			params: UpdateIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				DryRun:         true,
				ValidateMode:   "fast",
				ValidateRules:  false,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Comments: "default comment",
						Behaviors: []RuleBehavior{
							{
								Name: "origin",
								Options: RuleOptionsMap{
									"httpPort":           float64(80),
									"enableTrueClientIp": false,
									"compress":           true,
									"cacheKeyHostname":   "ORIGIN_HOSTNAME",
									"forwardHostHeader":  "REQUEST_HOST_HEADER",
									"hostname":           "origin.test.com",
									"originType":         "CUSTOMER",
								},
							},
							{
								Name: "cpCode",
								Options: RuleOptionsMap{
									"value": map[string]interface{}{
										"id":   float64(12345),
										"name": "my CP code",
									},
								},
							},
						},
						Children: []Rules{
							{
								Behaviors: []RuleBehavior{
									{
										Name: "gzipResponse",
										Options: RuleOptionsMap{
											"behavior": "ALWAYS",
										},
									},
								},
								Criteria: []RuleBehavior{
									{
										Locked: false,
										Name:   "contentType",
										Options: RuleOptionsMap{
											"matchOperator":      "IS_ONE_OF",
											"matchWildcard":      true,
											"matchCaseSensitive": false,
											"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
										},
									},
								},
								Name: "Compress Text Content",
							},
						},
						Criteria: []RuleBehavior{},
						Name:     "default",
						Options:  RuleOptions{IsSecure: false},
						CustomOverride: &RuleCustomOverride{
							OverrideID: "cbo_12345",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: "This is a sample Property Manager variable.",
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       "default value",
							},
						},
					},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating rule tree",
    "status": 500
}`,
			expectedPath: "/papi/v1/includes/inc_123456/versions/2/rules?contractId=test_contract&dryRun=true&groupId=test_group&validateMode=fast&validateRules=false",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating rule tree",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error - missing includeId": {
			params: UpdateIncludeRuleTreeRequest{
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				DryRun:         true,
				ValidateMode:   "fast",
				ValidateRules:  false,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Comments: "default comment",
						Behaviors: []RuleBehavior{
							{
								Name: "origin",
								Options: RuleOptionsMap{
									"httpPort":           float64(80),
									"enableTrueClientIp": false,
									"compress":           true,
									"cacheKeyHostname":   "ORIGIN_HOSTNAME",
									"forwardHostHeader":  "REQUEST_HOST_HEADER",
									"hostname":           "origin.test.com",
									"originType":         "CUSTOMER",
								},
							},
							{
								Name: "cpCode",
								Options: RuleOptionsMap{
									"value": map[string]interface{}{
										"id":   float64(12345),
										"name": "my CP code",
									},
								},
							},
						},
						Children: []Rules{
							{
								Behaviors: []RuleBehavior{
									{
										Name: "gzipResponse",
										Options: RuleOptionsMap{
											"behavior": "ALWAYS",
										},
									},
								},
								Criteria: []RuleBehavior{
									{
										Locked: false,
										Name:   "contentType",
										Options: RuleOptionsMap{
											"matchOperator":      "IS_ONE_OF",
											"matchWildcard":      true,
											"matchCaseSensitive": false,
											"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
										},
									},
								},
								Name: "Compress Text Content",
							},
						},
						Criteria: []RuleBehavior{},
						Name:     "default",
						Options:  RuleOptions{IsSecure: false},
						CustomOverride: &RuleCustomOverride{
							OverrideID: "cbo_12345",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: "This is a sample Property Manager variable.",
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       "default value",
							},
						},
					},
				},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing includeVersion": {
			params: UpdateIncludeRuleTreeRequest{
				IncludeID:     "inc_123456",
				ContractID:    "test_contract",
				GroupID:       "test_group",
				DryRun:        true,
				ValidateMode:  "fast",
				ValidateRules: false,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Comments: "default comment",
						Behaviors: []RuleBehavior{
							{
								Name: "origin",
								Options: RuleOptionsMap{
									"httpPort":           float64(80),
									"enableTrueClientIp": false,
									"compress":           true,
									"cacheKeyHostname":   "ORIGIN_HOSTNAME",
									"forwardHostHeader":  "REQUEST_HOST_HEADER",
									"hostname":           "origin.test.com",
									"originType":         "CUSTOMER",
								},
							},
							{
								Name: "cpCode",
								Options: RuleOptionsMap{
									"value": map[string]interface{}{
										"id":   float64(12345),
										"name": "my CP code",
									},
								},
							},
						},
						Children: []Rules{
							{
								Behaviors: []RuleBehavior{
									{
										Name: "gzipResponse",
										Options: RuleOptionsMap{
											"behavior": "ALWAYS",
										},
									},
								},
								Criteria: []RuleBehavior{
									{
										Locked: false,
										Name:   "contentType",
										Options: RuleOptionsMap{
											"matchOperator":      "IS_ONE_OF",
											"matchWildcard":      true,
											"matchCaseSensitive": false,
											"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
										},
									},
								},
								Name: "Compress Text Content",
							},
						},
						Criteria: []RuleBehavior{},
						Name:     "default",
						Options:  RuleOptions{IsSecure: false},
						CustomOverride: &RuleCustomOverride{
							OverrideID: "cbo_12345",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: "This is a sample Property Manager variable.",
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       "default value",
							},
						},
					},
				},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing contractId": {
			params: UpdateIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				GroupID:        "test_group",
				DryRun:         true,
				ValidateMode:   "fast",
				ValidateRules:  false,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Comments: "default comment",
						Behaviors: []RuleBehavior{
							{
								Name: "origin",
								Options: RuleOptionsMap{
									"httpPort":           float64(80),
									"enableTrueClientIp": false,
									"compress":           true,
									"cacheKeyHostname":   "ORIGIN_HOSTNAME",
									"forwardHostHeader":  "REQUEST_HOST_HEADER",
									"hostname":           "origin.test.com",
									"originType":         "CUSTOMER",
								},
							},
							{
								Name: "cpCode",
								Options: RuleOptionsMap{
									"value": map[string]interface{}{
										"id":   float64(12345),
										"name": "my CP code",
									},
								},
							},
						},
						Children: []Rules{
							{
								Behaviors: []RuleBehavior{
									{
										Name: "gzipResponse",
										Options: RuleOptionsMap{
											"behavior": "ALWAYS",
										},
									},
								},
								Criteria: []RuleBehavior{
									{
										Locked: false,
										Name:   "contentType",
										Options: RuleOptionsMap{
											"matchOperator":      "IS_ONE_OF",
											"matchWildcard":      true,
											"matchCaseSensitive": false,
											"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
										},
									},
								},
								Name: "Compress Text Content",
							},
						},
						Criteria: []RuleBehavior{},
						Name:     "default",
						Options:  RuleOptions{IsSecure: false},
						CustomOverride: &RuleCustomOverride{
							OverrideID: "cbo_12345",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: "This is a sample Property Manager variable.",
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       "default value",
							},
						},
					},
				},
			},
			withError: ErrStructValidation,
		},
		"validation error - missing groupId": {
			params: UpdateIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				DryRun:         true,
				ValidateMode:   "fast",
				ValidateRules:  false,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Comments: "default comment",
						Behaviors: []RuleBehavior{
							{
								Name: "origin",
								Options: RuleOptionsMap{
									"httpPort":           float64(80),
									"enableTrueClientIp": false,
									"compress":           true,
									"cacheKeyHostname":   "ORIGIN_HOSTNAME",
									"forwardHostHeader":  "REQUEST_HOST_HEADER",
									"hostname":           "origin.test.com",
									"originType":         "CUSTOMER",
								},
							},
							{
								Name: "cpCode",
								Options: RuleOptionsMap{
									"value": map[string]interface{}{
										"id":   float64(12345),
										"name": "my CP code",
									},
								},
							},
						},
						Children: []Rules{
							{
								Behaviors: []RuleBehavior{
									{
										Name: "gzipResponse",
										Options: RuleOptionsMap{
											"behavior": "ALWAYS",
										},
									},
								},
								Criteria: []RuleBehavior{
									{
										Locked: false,
										Name:   "contentType",
										Options: RuleOptionsMap{
											"matchOperator":      "IS_ONE_OF",
											"matchWildcard":      true,
											"matchCaseSensitive": false,
											"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
										},
									},
								},
								Name: "Compress Text Content",
							},
						},
						Criteria: []RuleBehavior{},
						Name:     "default",
						Options:  RuleOptions{IsSecure: false},
						CustomOverride: &RuleCustomOverride{
							OverrideID: "cbo_12345",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: "This is a sample Property Manager variable.",
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       "default value",
							},
						},
					},
				},
			},
			withError: ErrStructValidation,
		},
		"validation error - invalid validation mode": {
			params: UpdateIncludeRuleTreeRequest{
				IncludeID:      "inc_123456",
				IncludeVersion: 2,
				ContractID:     "test_contract",
				GroupID:        "test_group",
				DryRun:         true,
				ValidateMode:   "test",
				ValidateRules:  false,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Comments: "default comment",
						Behaviors: []RuleBehavior{
							{
								Name: "origin",
								Options: RuleOptionsMap{
									"httpPort":           float64(80),
									"enableTrueClientIp": false,
									"compress":           true,
									"cacheKeyHostname":   "ORIGIN_HOSTNAME",
									"forwardHostHeader":  "REQUEST_HOST_HEADER",
									"hostname":           "origin.test.com",
									"originType":         "CUSTOMER",
								},
							},
							{
								Name: "cpCode",
								Options: RuleOptionsMap{
									"value": map[string]interface{}{
										"id":   float64(12345),
										"name": "my CP code",
									},
								},
							},
						},
						Children: []Rules{
							{
								Behaviors: []RuleBehavior{
									{
										Name: "gzipResponse",
										Options: RuleOptionsMap{
											"behavior": "ALWAYS",
										},
									},
								},
								Criteria: []RuleBehavior{
									{
										Locked: false,
										Name:   "contentType",
										Options: RuleOptionsMap{
											"matchOperator":      "IS_ONE_OF",
											"matchWildcard":      true,
											"matchCaseSensitive": false,
											"values":             []interface{}{"text/html*", "text/css*", "application/x-javascript*"},
										},
									},
								},
								Name: "Compress Text Content",
							},
						},
						Criteria: []RuleBehavior{},
						Name:     "default",
						Options:  RuleOptions{IsSecure: false},
						CustomOverride: &RuleCustomOverride{
							OverrideID: "cbo_12345",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: "This is a sample Property Manager variable.",
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       "default value",
							},
						},
					},
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

				if len(test.responseHeaders) > 0 {
					for header, value := range test.responseHeaders {
						w.Header().Set(header, value)
					}
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateIncludeRuleTree(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
