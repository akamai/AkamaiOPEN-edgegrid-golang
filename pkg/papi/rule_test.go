package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapi_GetRuleTree(t *testing.T) {
	tests := map[string]struct {
		params           GetRuleTreeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRuleTreeResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetRuleTreeRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "fast",
				ValidateRules:   false,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "accountID",
    "contractId": "contract",
    "groupId": "group",
    "propertyId": "propertyID",
    "propertyVersion": 2,
    "etag": "etag",
    "ruleFormat": "v2020-09-16",
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
			expectedPath: "/papi/v1/properties/propertyID/versions/2/rules?contractId=contract&groupId=group&validateMode=fast&validateRules=false",
			expectedResponse: &GetRuleTreeResponse{
				Response: Response{
					AccountID:  "accountID",
					ContractID: "contract",
					GroupID:    "group",
				},
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				Etag:            "etag",
				RuleFormat:      "v2020-09-16",
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			params: GetRuleTreeRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "fast",
				ValidateRules:   false,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching rule tree",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2/rules?contractId=contract&groupId=group&validateMode=fast&validateRules=false",
			withError: func(t *testing.T, err error) {
				want := session.APIError{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching rule tree",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty property ID": {
			params: GetRuleTreeRequest{
				PropertyID:      "",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "fast",
				ValidateRules:   false,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"empty property version": {
			params: GetRuleTreeRequest{
				PropertyID:    "propertyID",
				ContractID:    "contract",
				GroupID:       "group",
				ValidateMode:  "fast",
				ValidateRules: false,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyVersion")
			},
		},
		"invalid validation mode": {
			params: GetRuleTreeRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "test",
				ValidateRules:   false,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ValidateMode")
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
			result, err := client.GetRuleTree(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_UpdateRuleTree(t *testing.T) {
	tests := map[string]struct {
		params           UpdateRulesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateRulesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "fast",
				ValidateRules:   false,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "accountID",
    "contractId": "contract",
    "groupId": "group",
    "propertyId": "propertyID",
    "propertyVersion": 2,
    "etag": "etag",
    "ruleFormat": "v2020-09-16",
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
    },
	"errors": [
        {
            "type": "/papi/v1/errors/validation.required_behavior",
            "title": "Missing required behavior in default rule",
            "detail": "In order for this property to work correctly behavior Content Provider Code needs to be present in the default section",
            "instance": "/papi/v1/properties/prp_173136/versions/3/rules#err_100",
            "behaviorName": "cpCode"
        }
    ]
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2/rules?contractId=contract&groupId=group&validateMode=fast&validateRules=false&dryRun=true",
			expectedResponse: &UpdateRulesResponse{
				AccountID:       "accountID",
				ContractID:      "contract",
				GroupID:         "group",
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				Etag:            "etag",
				RuleFormat:      "v2020-09-16",
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
				Errors: []RuleError{
					{
						Type:         "/papi/v1/errors/validation.required_behavior",
						Title:        "Missing required behavior in default rule",
						Detail:       "In order for this property to work correctly behavior Content Provider Code needs to be present in the default section",
						Instance:     "/papi/v1/properties/prp_173136/versions/3/rules#err_100",
						BehaviorName: "cpCode",
					},
				},
			},
		},
		"500 Internal Server Error": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "fast",
				ValidateRules:   false,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating rule tree",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2/rules?contractId=contract&groupId=group&validateMode=fast&validateRules=false&dryRun=true",
			withError: func(t *testing.T, err error) {
				want := session.APIError{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error updating rule tree",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"empty property ID": {
			params: UpdateRulesRequest{
				PropertyID:      "",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "fast",
				ValidateRules:   false,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyID")
			},
		},
		"empty property version": {
			params: UpdateRulesRequest{
				PropertyID:    "propertyID",
				ContractID:    "contract",
				GroupID:       "group",
				DryRun:        true,
				ValidateMode:  "fast",
				ValidateRules: false,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "PropertyVersion")
			},
		},
		"invalid validation mode": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "test",
				ValidateRules:   false,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ValidateMode")
			},
		},
		"empty behaviors": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "fast",
				ValidateRules:   false,
				Rules: Rules{
					Behaviors: nil,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Behaviors")
			},
		},
		"empty name": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "fast",
				ValidateRules:   false,
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
									Locked: "",
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
					Name:     "",
					Options:  &RuleOptions{IsSecure: false},
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Name")
			},
		},
		"empty name in behavior": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "test",
				ValidateRules:   false,
				Rules: Rules{
					Behaviors: []RuleBehavior{
						{
							Name: "",
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Name")
			},
		},
		"empty options in behavior": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "test",
				ValidateRules:   false,
				Rules: Rules{
					Behaviors: []RuleBehavior{
						{
							Name:    "origin",
							Options: nil,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Options")
			},
		},
		"empty name in custom override": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "test",
				ValidateRules:   false,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
					CustomOverride: &RuleCustomOverride{
						OverrideID: "cbo_12345",
						Name:       "",
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Name")
			},
		},
		"invalid overrideID in customOverride": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "test",
				ValidateRules:   false,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
					CustomOverride: &RuleCustomOverride{
						OverrideID: "",
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
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "OverrideID")
			},
		},
		"empty name in variable": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "test",
				ValidateRules:   false,
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
									Locked: "",
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
					Options:  &RuleOptions{IsSecure: false},
					CustomOverride: &RuleCustomOverride{
						OverrideID: "cbo_12345",
						Name:       "mdc",
					},
					Variables: []RuleVariable{
						{
							Description: "This is a sample Property Manager variable.",
							Hidden:      false,
							Name:        "",
							Sensitive:   false,
							Value:       "default value",
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Name")
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
			result, err := client.UpdateRuleTree(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
