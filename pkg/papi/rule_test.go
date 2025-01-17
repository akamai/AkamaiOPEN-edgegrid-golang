package papi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiGetRuleTree(t *testing.T) {
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
							Description: ptr.To("This is a sample Property Manager variable."),
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       ptr.To("default value"),
						},
					},
				},
			},
		},
		"200 OK with originalInput set to true": {
			params: GetRuleTreeRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "fast",
				ValidateRules:   false,
				OriginalInput:   ptr.To(true),
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
							Description: ptr.To("This is a sample Property Manager variable."),
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       ptr.To("default value"),
						},
					},
				},
			},
		},
		"200 OK with originalInput set to false": {
			params: GetRuleTreeRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "fast",
				ValidateRules:   false,
				OriginalInput:   ptr.To(false),
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
			expectedPath: "/papi/v1/properties/propertyID/versions/2/rules?contractId=contract&groupId=group&originalInput=false&validateMode=fast&validateRules=false",
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
							Description: ptr.To("This is a sample Property Manager variable."),
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       ptr.To("default value"),
						},
					},
				},
			},
		},
		"200 OK nested with empty behaviour": {
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
		"children": [{
			"name": "Augment insights",
			"criteria": [],
			"children": [{
				"name": "Compress Text Content",
				"criteria": [{
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
				}],
				"behaviors": []
			}],
			"options": {
				"is_secure": false
			}
		}],
		"behaviors": [{
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
		"variables": [{
			"name": "VAR_NAME",
			"value": "default value",
			"description": "This is a sample Property Manager variable.",
			"hidden": false,
			"sensitive": false
		}]
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
							Criteria: []RuleBehavior{},
							Name:     "Augment insights",

							Children: []Rules{
								{
									Behaviors: []RuleBehavior{},
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
									Name: "Compress Text Content"},
							}}},
					Name:    "default",
					Options: RuleOptions{IsSecure: false},
					CustomOverride: &RuleCustomOverride{
						OverrideID: "cbo_12345",
						Name:       "mdc",
					},
					Variables: []RuleVariable{
						{
							Description: ptr.To("This is a sample Property Manager variable."),
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       ptr.To("default value"),
						},
					},
				},
			},
		},
		"200 OK nested with empty name": {
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
		"children": [{
			"criteria": [],
			"children": [{
				"name": "Compress Text Content",
				"criteria": [{
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
				}],
				"behaviors": []
			}],
			"options": {
				"is_secure": false
			}
		}],
		"behaviors": [{
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
		"variables": [{
			"name": "VAR_NAME",
			"value": "default value",
			"description": "This is a sample Property Manager variable.",
			"hidden": false,
			"sensitive": false
		}]
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
							Criteria: []RuleBehavior{},
							Children: []Rules{
								{
									Behaviors: []RuleBehavior{},
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
									Name: "Compress Text Content"},
							}}},
					Name:    "default",
					Options: RuleOptions{IsSecure: false},
					CustomOverride: &RuleCustomOverride{
						OverrideID: "cbo_12345",
						Name:       "mdc",
					},
					Variables: []RuleVariable{
						{
							Description: ptr.To("This is a sample Property Manager variable."),
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       ptr.To("default value"),
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
				want := &Error{
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
		"Accept header with latest ruleFormat": {
			params: GetRuleTreeRequest{
				PropertyID:      "1",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "fast",
				ValidateRules:   false,
				RuleFormat:      "latest",
			},
			responseStatus:   http.StatusOK,
			responseBody:     "{}",
			expectedPath:     "/papi/v1/properties/1/versions/2/rules?contractId=contract&groupId=group&validateMode=fast&validateRules=false",
			expectedResponse: &GetRuleTreeResponse{},
		},
		"Accept header with versioned ruleFormat": {
			params: GetRuleTreeRequest{
				PropertyID:      "1",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "fast",
				ValidateRules:   false,
				RuleFormat:      "v2021-01-01",
			},
			responseStatus:   http.StatusOK,
			responseBody:     "{}",
			expectedPath:     "/papi/v1/properties/1/versions/2/rules?contractId=contract&groupId=group&validateMode=fast&validateRules=false",
			expectedResponse: &GetRuleTreeResponse{},
		},
		"Accept header with invalid ruleFormat": {
			params: GetRuleTreeRequest{
				PropertyID:      "1",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				ValidateMode:    "fast",
				ValidateRules:   false,
				RuleFormat:      "invalid",
			},
			responseStatus: http.StatusUnsupportedMediaType,
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "RuleFormat")
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
				if test.params.RuleFormat != "" {
					assert.Equal(t, r.Header.Get("Accept"), fmt.Sprintf("application/vnd.akamai.papirules.%s+json", test.params.RuleFormat))
				}
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

func TestPapiUpdateRuleTree(t *testing.T) {
	tests := map[string]struct {
		params           UpdateRulesRequest
		requestBody      string
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
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       ptr.To("default value"),
							},
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
            "type": "/papi/v1/errors/validation.required_behavior",
            "title": "Missing required behavior in default rule",
            "detail": "In order for this property to work correctly behavior Content Provider Code needs to be present in the default section",
            "instance": "/papi/v1/properties/prp_173136/versions/3/rules#err_100",
            "behaviorName": "cpCode"
        }
    ]
}`,
			expectedPath: "/papi/v1/properties/propertyID/versions/2/rules?contractId=contract&dryRun=true&groupId=group&validateMode=fast&validateRules=false",
			expectedResponse: &UpdateRulesResponse{
				AccountID:       "accountID",
				ContractID:      "contract",
				GroupID:         "group",
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				Etag:            "etag",
				RuleFormat:      "v2020-09-16",
				Comments:        "version comment",
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
							Description: ptr.To("This is a sample Property Manager variable."),
							Hidden:      false,
							Name:        "VAR_NAME",
							Sensitive:   false,
							Value:       ptr.To("default value"),
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
		"200 OK - empty and null values for description": {
			params: UpdateRulesRequest{
				ContractID:      "ctr_id",
				GroupID:         "grp_id",
				PropertyID:      "prp_id",
				PropertyVersion: 1,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Name: "default",
						Children: []Rules{
							{
								Name:     "change fwd path",
								Children: []Rules{},
								Behaviors: []RuleBehavior{
									{
										Name: "baseDirectory",
										Options: RuleOptionsMap{
											"value": "/smth/",
										},
									},
								},
								Criteria: []RuleBehavior{
									{
										Name:   "requestHeader",
										Locked: false,
										Options: RuleOptionsMap{
											"headerName":              "Accept-Encoding",
											"matchCaseSensitiveValue": true,
											"matchOperator":           "IS_ONE_OF",
											"matchWildcardName":       false,
											"matchWildcardValue":      false,
										},
									},
								},
								CriteriaMustSatisfy: RuleCriteriaMustSatisfyAll,
							},
							{
								Name:     "caching",
								Children: []Rules{},
								Behaviors: []RuleBehavior{
									{
										Name: "caching",
										Options: RuleOptionsMap{
											"behavior":       "MAX_AGE",
											"mustRevalidate": false,
											"ttl":            "1m",
										},
									},
								},
								Criteria:            []RuleBehavior{},
								CriteriaMustSatisfy: RuleCriteriaMustSatisfyAny,
							},
						},
						Behaviors: []RuleBehavior{
							{
								Name: "origin",
								Options: RuleOptionsMap{
									"cacheKeyHostname":          "REQUEST_HOST_HEADER",
									"compress":                  true,
									"enableTrueClientIp":        true,
									"forwardHostHeader":         "REQUEST_HOST_HEADER",
									"hostname":                  "httpbin.smth.online",
									"httpPort":                  float64(80),
									"httpsPort":                 float64(443),
									"originCertificate":         "",
									"originSni":                 true,
									"originType":                "CUSTOMER",
									"ports":                     "",
									"trueClientIpClientSetting": false,
									"trueClientIpHeader":        "True-Client-IP",
									"verificationMode":          "PLATFORM_SETTINGS",
								},
							},
						},
						Options:  RuleOptions{},
						Criteria: []RuleBehavior{},
						Variables: []RuleVariable{
							{
								Name:        "TEST_EMPTY_FIELDS",
								Value:       ptr.To(""),
								Description: ptr.To(""),
								Hidden:      true,
								Sensitive:   false,
							},
							{
								Name:        "TEST_NIL_DESCRIPTION",
								Description: nil,
								Value:       ptr.To(""),
								Hidden:      true,
								Sensitive:   false,
							},
						},
						Comments: "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.",
					},
				},
			},
			requestBody:    `{"comments":"version comment","rules":{"behaviors":[{"name":"origin","options":{"cacheKeyHostname":"REQUEST_HOST_HEADER","compress":true,"enableTrueClientIp":true,"forwardHostHeader":"REQUEST_HOST_HEADER","hostname":"httpbin.smth.online","httpPort":80,"httpsPort":443,"originCertificate":"","originSni":true,"originType":"CUSTOMER","ports":"","trueClientIpClientSetting":false,"trueClientIpHeader":"True-Client-IP","verificationMode":"PLATFORM_SETTINGS"}}],"children":[{"behaviors":[{"name":"baseDirectory","options":{"value":"/smth/"}}],"criteria":[{"name":"requestHeader","options":{"headerName":"Accept-Encoding","matchCaseSensitiveValue":true,"matchOperator":"IS_ONE_OF","matchWildcardName":false,"matchWildcardValue":false}}],"name":"change fwd path","options":{},"criteriaMustSatisfy":"all"},{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"1m"}}],"name":"caching","options":{},"criteriaMustSatisfy":"any"}],"comments":"The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.","name":"default","options":{},"variables":[{"description":"","hidden":true,"name":"TEST_EMPTY_FIELDS","sensitive":false,"value":""},{"description":null,"hidden":true,"name":"TEST_NIL_DESCRIPTION","sensitive":false,"value":""}]}}`,
			responseStatus: http.StatusOK,
			responseBody: `{
    "accountId": "act_id",
    "contractId": "ctr_id",
    "groupId": "grp_id",
    "propertyId": "grp_id",
    "propertyName": "dxe-2407-reproducing",
    "propertyVersion": 1,
    "etag": "123ge123",
    "rules": {
        "name": "default",
        "children": [
            {
                "name": "change fwd path",
                "children": [],
                "behaviors": [
                    {
                        "name": "baseDirectory",
                        "options": {
                            "value": "/smth/"
                        }
                    }
                ],
                "criteria": [
                    {
                        "name": "requestHeader",
                        "options": {
                            "headerName": "Accept-Encoding",
                            "matchCaseSensitiveValue": true,
                            "matchOperator": "IS_ONE_OF",
                            "matchWildcardName": false,
                            "matchWildcardValue": false
                        }
                    }
                ],
                "criteriaMustSatisfy": "all"
            },
            {
                "name": "caching",
                "children": [],
                "behaviors": [
                    {
                        "name": "caching",
                        "options": {
                            "behavior": "MAX_AGE",
                            "mustRevalidate": false,
                            "ttl": "1m"
                        }
                    }
                ],
                "criteria": [],
                "criteriaMustSatisfy": "any"
            }
        ],
        "behaviors": [
            {
                "name": "origin",
                "options": {
                    "cacheKeyHostname": "REQUEST_HOST_HEADER",
                    "compress": true,
                    "enableTrueClientIp": true,
                    "forwardHostHeader": "REQUEST_HOST_HEADER",
                    "hostname": "httpbin.smth.online",
                    "httpPort": 80,
                    "httpsPort": 443,
                    "originCertificate": "",
                    "originSni": true,
                    "originType": "CUSTOMER",
                    "ports": "",
                    "trueClientIpClientSetting": false,
                    "trueClientIpHeader": "True-Client-IP",
                    "verificationMode": "PLATFORM_SETTINGS"
                }
            }
        ],
        "options": {},
        "variables": [
            {
                "name": "TEST_EMPTY_FIELDS",
                "value": "",
                "description": "",
                "hidden": true,
                "sensitive": false
            },
            {
                "name": "TEST_NIL_FIELDS",
                "value": "",
                "description": null,
                "hidden": true,
                "sensitive": false
            }
        ],
        "comments": "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings."
    },
    "errors": [
        {
            "type": "https://problems.luna.akamaiapis.net/papi/v0/validation/required_feature_any",
            "errorLocation": "#/rules",
            "detail": "Add a Content Provider Code behavior in any rule. See <a href=\"/dl/property-manager/property-manager-help/csh_lookup.html?id=PM_0030\" target=\"_blank\">Content Provider Code</a>."
        }
    ],
    "warnings": [
        {
            "type": "https://problems.luna.akamaiapis.net/papi/v0/validation/incompatible_condition",
            "errorLocation": "#/rules/children/0/behaviors/0"
        },
        {
            "title": "Unstable rule format",
            "type": "https://problems.luna.akamaiapis.net/papi/v0/unstable_rule_format",
            "currentRuleFormat": "latest",
            "suggestedRuleFormat": "v2023-01-05"
        }
    ],
    "ruleFormat": "latest"
}`,
			expectedPath: "/papi/v1/properties/prp_id/versions/1/rules?contractId=ctr_id&groupId=grp_id&validateRules=false",
			expectedResponse: &UpdateRulesResponse{
				AccountID:       "act_id",
				ContractID:      "ctr_id",
				GroupID:         "grp_id",
				PropertyID:      "grp_id",
				PropertyVersion: 1,
				Etag:            "123ge123",
				RuleFormat:      "latest",
				Rules: Rules{
					Name: "default",
					Children: []Rules{
						{
							Name:     "change fwd path",
							Children: []Rules{},
							Behaviors: []RuleBehavior{
								{
									Name: "baseDirectory",
									Options: RuleOptionsMap{
										"value": "/smth/",
									},
								},
							},
							Criteria: []RuleBehavior{
								{
									Name:   "requestHeader",
									Locked: false,
									Options: RuleOptionsMap{
										"headerName":              "Accept-Encoding",
										"matchCaseSensitiveValue": true,
										"matchOperator":           "IS_ONE_OF",
										"matchWildcardName":       false,
										"matchWildcardValue":      false,
									},
								},
							},
							CriteriaMustSatisfy: "all",
						},
						{
							Name:     "caching",
							Children: []Rules{},
							Behaviors: []RuleBehavior{
								{
									Name: "caching",
									Options: RuleOptionsMap{
										"behavior":       "MAX_AGE",
										"mustRevalidate": false,
										"ttl":            "1m",
									},
								},
							},
							Criteria:            []RuleBehavior{},
							CriteriaMustSatisfy: "any",
						},
					},
					Behaviors: []RuleBehavior{
						{
							Name: "origin",
							Options: RuleOptionsMap{
								"cacheKeyHostname":          "REQUEST_HOST_HEADER",
								"compress":                  true,
								"enableTrueClientIp":        true,
								"forwardHostHeader":         "REQUEST_HOST_HEADER",
								"hostname":                  "httpbin.smth.online",
								"httpPort":                  float64(80),
								"httpsPort":                 float64(443),
								"originCertificate":         "",
								"originSni":                 true,
								"originType":                "CUSTOMER",
								"ports":                     "",
								"trueClientIpClientSetting": false,
								"trueClientIpHeader":        "True-Client-IP",
								"verificationMode":          "PLATFORM_SETTINGS",
							},
						},
					},
					Variables: []RuleVariable{
						{
							Name:        "TEST_EMPTY_FIELDS",
							Value:       ptr.To(""),
							Description: ptr.To(""),
							Hidden:      true,
							Sensitive:   false,
						},
						{
							Name:        "TEST_NIL_FIELDS",
							Description: nil,
							Value:       ptr.To(""),
							Hidden:      true,
							Sensitive:   false,
						},
					},
					Comments: "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.",
				},
				Errors: []RuleError{
					{
						Type:          "https://problems.luna.akamaiapis.net/papi/v0/validation/required_feature_any",
						ErrorLocation: "#/rules",
						Detail:        "Add a Content Provider Code behavior in any rule. See <a href=\"/dl/property-manager/property-manager-help/csh_lookup.html?id=PM_0030\" target=\"_blank\">Content Provider Code</a>.",
					},
				},
				Warnings: []RuleWarnings{
					{
						Type:          "https://problems.luna.akamaiapis.net/papi/v0/validation/incompatible_condition",
						ErrorLocation: "#/rules/children/0/behaviors/0",
					},
					{
						Title:               "Unstable rule format",
						Type:                "https://problems.luna.akamaiapis.net/papi/v0/unstable_rule_format",
						CurrentRuleFormat:   "latest",
						SuggestedRuleFormat: "v2023-01-05",
					},
				},
			},
		},
		"validation error - value is null": {
			params: UpdateRulesRequest{
				ContractID:      "ctr_id",
				GroupID:         "grp_id",
				PropertyID:      "prp_id",
				PropertyVersion: 1,
				Rules: RulesUpdate{
					Comments: "version comment",
					Rules: Rules{
						Name: "default",
						Variables: []RuleVariable{
							{
								Name:        "TEST_EMPTY_FIELDS",
								Value:       ptr.To(""),
								Description: ptr.To(""),
								Hidden:      true,
								Sensitive:   false,
							},
							{
								Name:        "TEST_NIL_FIELDS",
								Description: nil,
								Value:       nil,
								Hidden:      true,
								Sensitive:   false,
							},
						},
						Comments: "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.",
					},
				},
			},
			expectedPath: "/papi/v1/properties/prp_id/versions/1/rules?contractId=ctr_id&groupId=grp_id&validateRules=false",
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "updating rule tree: struct validation:\nRules: {\n\tRules: {\n\t\tVariables[1]: {\n\t\t\tValue: is required\n\t\t}\n\t}\n}")
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
				Rules: RulesUpdate{
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
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       ptr.To("default value"),
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
			expectedPath: "/papi/v1/properties/propertyID/versions/2/rules?contractId=contract&dryRun=true&groupId=group&validateMode=fast&validateRules=false",
			withError: func(t *testing.T, err error) {
				want := &Error{
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
				Rules: RulesUpdate{
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
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       ptr.To("default value"),
							},
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
				Rules: RulesUpdate{
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
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       ptr.To("default value"),
							},
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
				Rules: RulesUpdate{
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
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       ptr.To("default value"),
							},
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
		"empty name": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "fast",
				ValidateRules:   false,
				Rules: RulesUpdate{
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
						Name:     "",
						Options:  RuleOptions{IsSecure: false},
						CustomOverride: &RuleCustomOverride{
							OverrideID: "cbo_12345",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       ptr.To("default value"),
							},
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
		"empty name in custom override": {
			params: UpdateRulesRequest{
				PropertyID:      "propertyID",
				PropertyVersion: 2,
				ContractID:      "contract",
				GroupID:         "group",
				DryRun:          true,
				ValidateMode:    "test",
				ValidateRules:   false,
				Rules: RulesUpdate{
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
							Name:       "",
						},
						Variables: []RuleVariable{
							{
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       ptr.To("default value"),
							},
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
				Rules: RulesUpdate{
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
							OverrideID: "",
							Name:       "mdc",
						},
						Variables: []RuleVariable{
							{
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "VAR_NAME",
								Sensitive:   false,
								Value:       ptr.To("default value"),
							},
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
				Rules: RulesUpdate{
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
								Description: ptr.To("This is a sample Property Manager variable."),
								Hidden:      false,
								Name:        "",
								Sensitive:   false,
								Value:       ptr.To("default value"),
							},
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
