package v3

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSONMatchRules(t *testing.T) {
	tests := map[string]struct {
		withError      error
		responseBody   string
		expectedObject MatchRules
	}{
		"invalid MatchRuleXX": {
			responseBody: `
	[
        {
            "type": "xxMatchRule"
        }
    ]
`,
			withError: errors.New("unmarshalling MatchRules: unsupported match rule type: xxMatchRule"),
		},

		"invalid type": {
			withError: errors.New("unmarshalling MatchRules: 'type' field on match rule entry should be a string"),
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
            "type": "erMatchRule"
        }
    
`,
		},

		"missing type": {
			withError: errors.New("unmarshalling MatchRules: match rule entry should contain 'type' field"),
			responseBody: `
	[
        {
        }
    ]
`,
		},

		"invalid objectMatchValue type for PR - range": {
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaPR: objectMatchValue has unexpected type: 'range'"),
			responseBody: `
	[
        {
            "type": "cdMatchRule",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
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
            "name": "Rule3",
            "start": 0
        }
    ]
`,
		},

		"invalid objectMatchValue type for ER - range": {
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaER: objectMatchValue has unexpected type: 'range'"),
			responseBody: `
	[
        {
            "type": "erMatchRule",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
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
            "name": "Rule3",
            "start": 0
        }
    ]
`,
		},

		"invalid objectMatchValue type for FR - range": {
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaFR: objectMatchValue has unexpected type: 'range'"),
			responseBody: `
	[
        {
            "type": "frMatchRule",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
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
            "name": "Rule3",
            "start": 0
        }
    ]
`,
		},

		"valid MatchRulePR": {
			responseBody: `
	[
        {
            "type": "cdMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
            "forwardSettings": {
                "originId": "fr_test_krk_dc2",
                "percent": 62
            },
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
				&MatchRulePR{
					Type:     "cdMatchRule",
					End:      0,
					ID:       0,
					MatchURL: "",
					Matches: []MatchCriteriaPR{
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
							MatchType:     "method",
							Negate:        false,
							ObjectMatchValue: &ObjectMatchValueSimple{
								Type:  "simple",
								Value: []string{"GET"},
							},
						},
					},
					Name:  "Rule3",
					Start: 0,
					ForwardSettings: ForwardSettingsPR{
						OriginID: "fr_test_krk_dc2",
						Percent:  62,
					},
				},
			},
		},

		"valid MatchRuleFR": {
			responseBody: `
	[
        {
            "type": "frMatchRule",
            "end": 0,
            "id": 0,
            "matchURL": null,
			"forwardSettings": {},
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
				&MatchRuleFR{
					Type:     "frMatchRule",
					End:      0,
					ID:       0,
					MatchURL: "",
					Matches: []MatchCriteriaFR{
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
							MatchType:     "method",
							Negate:        false,
							ObjectMatchValue: &ObjectMatchValueSimple{
								Type:  "simple",
								Value: []string{"GET"},
							},
						},
					},
					Name:  "Rule3",
					Start: 0,
				},
			},
		},

		"invalid objectMatchValue type for AP - range": {
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaAP: objectMatchValue has unexpected type: 'range'"),
			responseBody: `
	[
        {
            "type": "apMatchRule",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
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
            "name": "Rule3",
            "start": 0
        }
    ]
`,
		},

		"valid MatchRuleAP": {
			responseBody: `
	[
        {
            "type": "apMatchRule",
            "end": 0,
            "passThroughPercent": 50.50,
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
				&MatchRuleAP{
					Type:               "apMatchRule",
					End:                0,
					PassThroughPercent: ptr.To(50.50),
					ID:                 0,
					MatchURL:           "",
					Matches: []MatchCriteriaAP{
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
							MatchType:     "method",
							Negate:        false,
							ObjectMatchValue: &ObjectMatchValueSimple{
								Type:  "simple",
								Value: []string{"GET"},
							},
						},
					},
					Name:  "Rule3",
					Start: 0,
				},
			},
		},
		"valid MatchRuleAS": {
			responseBody: `
	[
		{
            "name": "rule 10",
            "type": "asMatchRule",
            "matchURL": "http://source.com/test1",

            "forwardSettings": {
                "originId": "origin_remote_1",
                "pathAndQS": "/cpaths/test1.html"
            },

            "matches": [
                {
                    "matchType": "range",
                    "objectMatchValue": {
                        "type": "range",
                        "value": [  1, 100 ]
                    },
                    "matchOperator": "equals",
                    "negate": false,
                    "caseSensitive": false
                },
                {
                    "matchType": "header",
                    "objectMatchValue": {
                        "options": {
                            "value": [  "en" ]
                        },
                        "type": "object",
                        "name": "Accept-Charset"
                    },
                    "matchOperator": "equals",
                    "negate": false,
                    "caseSensitive": false
                }
            ]
        }
	]`,
			expectedObject: MatchRules{
				&MatchRuleAS{
					Name:     "rule 10",
					Type:     "asMatchRule",
					MatchURL: "http://source.com/test1",
					ForwardSettings: ForwardSettingsAS{
						OriginID:  "origin_remote_1",
						PathAndQS: "/cpaths/test1.html",
					},
					Matches: []MatchCriteriaAS{
						{
							MatchType: "range",
							ObjectMatchValue: &ObjectMatchValueRange{
								Type:  "range",
								Value: []int64{1, 100},
							},
							MatchOperator: "equals",
							CaseSensitive: false,
							Negate:        false,
						},
						{
							MatchType: "header",
							ObjectMatchValue: &ObjectMatchValueObject{
								Name: "Accept-Charset",
								Type: "object",
								Options: &Options{
									Value: []string{"en"},
								},
							},
							MatchOperator: "equals",
							Negate:        false,
							CaseSensitive: false,
						},
					},
				},
			},
		},

		"valid MatchRuleRC": {
			responseBody: `
	[
		{
			"type": "igMatchRule",
			"end": 0,
			"allowDeny": "allow",
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
	]`,
			expectedObject: MatchRules{
				&MatchRuleRC{
					Name:      "Rule3",
					Type:      "igMatchRule",
					AllowDeny: Allow,
					Matches: []MatchCriteriaRC{
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

		"invalid objectMatchValue type for RC - range": {
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaRC: objectMatchValue has unexpected type: 'range'"),
			responseBody: `
	[
        {
            "type": "igMatchRule",
			"allowDeny": "allow",
            "matches": [
                {
                    "caseSensitive": false,
                    "matchOperator": "equals",
                    "matchType": "method",
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
            "name": "Rule3",
            "start": 0
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

func TestGetObjectMatchValueType(t *testing.T) {
	tests := map[string]struct {
		withError error
		input     interface{}
		expected  string
	}{
		"success getting objectMatchValue type": {
			input: map[string]interface{}{
				"type":  "range",
				"value": []int{1, 50},
			},
			expected: "range",
		},
		"error getting objectMatchValue type - invalid type": {
			withError: errors.New("structure of objectMatchValue should be 'map', but was 'string'"),
			input:     "stringType",
		},
		"error getting objectMatchValue type - missing type": {
			withError: errors.New("objectMatchValue should contain 'type' field"),
			input: map[string]interface{}{
				"value": []int{1, 50},
			},
		},
		"error getting objectMatchValue type - type not string": {
			withError: errors.New("'type' should be a string"),
			input: map[string]interface{}{
				"type":  50,
				"value": []int{1, 50},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			objectMatchValueType, err := getObjectMatchValueType(test.input)

			if test.withError != nil {
				assert.Equal(t, test.withError.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, objectMatchValueType)
		})
	}
}

func TestConvertObjectMatchValue(t *testing.T) {
	tests := map[string]struct {
		withError bool
		input     map[string]interface{}
		output    interface{}
		expected  interface{}
	}{
		"success converting objectMatchValueRange": {
			input: map[string]interface{}{
				"type":  "range",
				"value": []int{1, 50},
			},
			output: &ObjectMatchValueRange{},
			expected: &ObjectMatchValueRange{
				Type:  "range",
				Value: []int64{1, 50},
			},
		},
		"success converting objectMatchValueSimple": {
			input: map[string]interface{}{
				"type":  "simple",
				"value": []string{"GET"},
			},
			output: &ObjectMatchValueSimple{},
			expected: &ObjectMatchValueSimple{
				Type:  "simple",
				Value: []string{"GET"},
			},
		},
		"success converting objectMatchValueObject": {
			input: map[string]interface{}{
				"type": "object",
				"name": "ER",
				"options": map[string]interface{}{
					"value": []string{
						"text/html*",
						"text/css*",
						"application/x-javascript*",
					},
					"valueHasWildcard": true,
				},
			},
			output: &ObjectMatchValueObject{},
			expected: &ObjectMatchValueObject{
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
		"error converting objectMatchValue": {
			withError: true,
			input: map[string]interface{}{
				"type":  "range",
				"value": []int{1, 50},
			},
			output: &ObjectMatchValueSimple{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			convertedObjectMatchValue, err := convertObjectMatchValue(test.input, test.output)

			if test.withError == true {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, convertedObjectMatchValue)
		})
	}
}

func TestValidateMatchRules(t *testing.T) {
	tests := map[string]struct {
		input     MatchRules
		withError string
	}{
		"valid match rules AP": {
			input: MatchRules{
				MatchRuleAP{
					Type:               "apMatchRule",
					PassThroughPercent: ptr.To(float64(-1)),
				},
				MatchRuleAP{
					Type:               "apMatchRule",
					PassThroughPercent: ptr.To(50.5),
				},
				MatchRuleAP{
					Type:               "apMatchRule",
					PassThroughPercent: ptr.To(float64(0)),
				},
				MatchRuleAP{
					Type:               "apMatchRule",
					PassThroughPercent: ptr.To(float64(100)),
				},
			},
		},
		"invalid match rules AP": {
			input: MatchRules{
				MatchRuleAP{
					Type: "matchRule",
				},
				MatchRuleAP{
					Type:               "apMatchRule",
					PassThroughPercent: ptr.To(100.1),
				},
				MatchRuleAP{
					Type:               "apMatchRule",
					PassThroughPercent: ptr.To(-1.1),
				},
			},
			withError: `
MatchRules[0]: {
	PassThroughPercent: cannot be blank
	Type: value 'matchRule' is invalid. Must be: 'apMatchRule'
}
MatchRules[1]: {
	PassThroughPercent: must be no greater than 100
}
MatchRules[2]: {
	PassThroughPercent: must be no less than -1
}`,
		},
		"valid match rules AS": {
			input: MatchRules{
				MatchRuleAS{
					Type:  "asMatchRule",
					Start: 0,
					End:   1,
				},
				MatchRuleAS{
					Type: "asMatchRule",
					ForwardSettings: ForwardSettingsAS{
						PathAndQS: "something",
						OriginID:  "something_else",
					},
				},
			},
		},
		"invalid match rules AS": {
			input: MatchRules{
				MatchRuleAS{
					Type: "matchRule",
				},
				MatchRuleAS{
					Type:  "asMatchRule",
					Start: -2,
					End:   -1,
					ForwardSettings: ForwardSettingsAS{
						OriginID: "some_id",
					},
				},
			},

			withError: `
MatchRules[0]: {
	Type: value 'matchRule' is invalid. Must be: 'asMatchRule'
}
MatchRules[1]: {
	End: must be no less than 0
	Start: must be no less than 0
}`,
		},
		"valid match rules CD": {
			input: MatchRules{
				MatchRulePR{
					Type: "cdMatchRule",
					ForwardSettings: ForwardSettingsPR{
						OriginID: "testOriginID",
						Percent:  100,
					},
				},
				MatchRulePR{
					Type: "cdMatchRule",
					ForwardSettings: ForwardSettingsPR{
						OriginID: "testOriginID",
						Percent:  1,
					},
				},
			},
		},
		"invalid match rules CD": {
			input: MatchRules{
				MatchRulePR{
					Type: "matchRule",
				},
				MatchRulePR{
					Type:            "cdMatchRule",
					ForwardSettings: ForwardSettingsPR{},
				},
				MatchRulePR{
					Type: "cdMatchRule",
					ForwardSettings: ForwardSettingsPR{
						OriginID: "testOriginID",
						Percent:  101,
					},
				},
				MatchRulePR{
					Type: "cdMatchRule",
					ForwardSettings: ForwardSettingsPR{
						OriginID: "testOriginID",
						Percent:  -1,
					},
				},
				MatchRulePR{
					Type: "cdMatchRule",
					ForwardSettings: ForwardSettingsPR{
						OriginID: "testOriginID",
						Percent:  0,
					},
				},
			},
			withError: `
MatchRules[0]: {
	ForwardSettings.OriginID: cannot be blank
	ForwardSettings.Percent: cannot be blank
	Type: value 'matchRule' is invalid. Must be: 'cdMatchRule'
}
MatchRules[1]: {
	ForwardSettings.OriginID: cannot be blank
	ForwardSettings.Percent: cannot be blank
}
MatchRules[2]: {
	ForwardSettings.Percent: must be no greater than 100
}
MatchRules[3]: {
	ForwardSettings.Percent: must be no less than 1
}
MatchRules[4]: {
	ForwardSettings.Percent: cannot be blank
}`,
		},
		"valid match rules ER": {
			input: MatchRules{
				MatchRuleER{
					Type:           "erMatchRule",
					RedirectURL:    "abc.com",
					UseRelativeURL: "none",
					StatusCode:     301,
				},
				MatchRuleER{
					Type:        "erMatchRule",
					RedirectURL: "abc.com",
					StatusCode:  301,
				},
				MatchRuleER{
					Type:          "erMatchRule",
					RedirectURL:   "abc.com",
					MatchesAlways: true,
					StatusCode:    301,
				},
				MatchRuleER{
					Type:        "erMatchRule",
					RedirectURL: "abc.com",
					Matches: []MatchCriteriaER{
						{
							MatchValue: "asd",
						},
					},
					StatusCode: 301,
				},
			},
		},
		"invalid match rules ER": {
			input: MatchRules{
				MatchRuleER{
					Type: "matchRule",
				},
				MatchRuleER{
					Type:           "erMatchRule",
					RedirectURL:    "abc.com",
					UseRelativeURL: "test",
					StatusCode:     404,
				},
				MatchRuleER{
					Type:           "erMatchRule",
					RedirectURL:    "abc.com",
					UseRelativeURL: "none",
					StatusCode:     301,
					MatchesAlways:  true,
					Matches: []MatchCriteriaER{
						{
							MatchValue: "asd",
						},
					},
				},
			},
			withError: `
MatchRules[0]: {
	RedirectURL: cannot be blank
	StatusCode: cannot be blank
	Type: value 'matchRule' is invalid. Must be: 'erMatchRule'
}
MatchRules[1]: {
	StatusCode: value '404' is invalid. Must be one of: 301, 302, 303, 307 or 308
	UseRelativeURL: value 'test' is invalid. Must be one of: 'none', 'copy_scheme_hostname', 'relative_url' or '' (empty)
}
MatchRules[2]: {
	Matches/MatchesAlways: only one of [ "Matches", "MatchesAlways" ] can be specified
}`,
		},
		"valid match rules FR": {
			input: MatchRules{
				MatchRuleFR{
					Type: "frMatchRule",
					ForwardSettings: ForwardSettingsFR{
						PathAndQS: "test",
						OriginID:  "testOriginID",
					},
				},
				MatchRuleFR{
					Type: "frMatchRule",
					ForwardSettings: ForwardSettingsFR{
						PathAndQS: "test",
						OriginID:  "testOriginID",
					},
				},
			},
		},
		"invalid match rules FR": {
			input: MatchRules{
				MatchRuleFR{
					Type: "matchRule",
				},
				MatchRuleFR{
					Type: "frMatchRule",
					ForwardSettings: ForwardSettingsFR{
						OriginID:  "testOriginID",
						PathAndQS: "",
					},
				},
			},
			withError: `
MatchRules[0]: {
	Type: value 'matchRule' is invalid. Must be: 'frMatchRule'
}`,
		},
		"valid match rules RC": {
			input: MatchRules{
				MatchRuleRC{
					Type:      "igMatchRule",
					AllowDeny: Allow,
				},
				MatchRuleRC{
					Type:      "igMatchRule",
					AllowDeny: Deny,
				},
				MatchRuleRC{
					Type:      "igMatchRule",
					AllowDeny: DenyBranded,
				},
			},
		},
		"invalid match rules RC": {
			input: MatchRules{
				MatchRuleRC{
					Type: "invalidMatchRule",
				},
				MatchRuleRC{
					Type:      "igMatchRule",
					AllowDeny: "allowBranded",
				},
				MatchRuleRC{
					Type:          "igMatchRule",
					AllowDeny:     Allow,
					MatchesAlways: true,
					Matches: []MatchCriteriaRC{
						{
							CaseSensitive: false,
							CheckIPs:      "CONNECTING_IP",
							MatchOperator: "equals",
							MatchType:     "clientip",
							MatchValue:    "1.2.3.4",
							Negate:        false,
						},
					},
				},
			},
			withError: `
MatchRules[0]: {
	AllowDeny: cannot be blank
	Type: value 'invalidMatchRule' is invalid. Must be: 'igMatchRule'
}
MatchRules[1]: {
	AllowDeny: value 'allowBranded' is invalid. Must be one of: 'allow', 'deny' or 'denybranded'
}
MatchRules[2]: {
	Matches/MatchesAlways: only one of [ "Matches", "MatchesAlways" ] can be specified
}`,
		},
		"valid match criteria - matchValue": {
			input: MatchRules{
				MatchRuleER{
					Type:        "erMatchRule",
					RedirectURL: "abc.com",
					StatusCode:  301,
					Matches: []MatchCriteriaER{
						{
							MatchType:     "method",
							MatchOperator: "equals",
							CheckIPs:      "CONNECTING_IP",
							MatchValue:    "https",
						},
					},
				},
			},
		},
		"valid match criteria - object match value": {
			input: MatchRules{
				MatchRuleER{
					Type:        "erMatchRule",
					RedirectURL: "abc.com",
					StatusCode:  301,
					Matches: []MatchCriteriaER{
						{
							MatchType:     "header",
							MatchOperator: "equals",
							CheckIPs:      "CONNECTING_IP",
							ObjectMatchValue: &ObjectMatchValueSimple{
								Type: "simple",
								Value: []string{
									"GET",
								},
							},
						},
					},
				},
			},
		},
		"invalid match criteria - matchValue and omv combinations": {
			input: MatchRules{
				MatchRuleER{
					Type:        "erMatchRule",
					RedirectURL: "abc.com",
					StatusCode:  301,
					Matches: []MatchCriteriaER{
						{
							MatchType:     "header",
							MatchOperator: "equals",
							CheckIPs:      "CONNECTING_IP",
							ObjectMatchValue: &ObjectMatchValueSimple{
								Type: "simple",
								Value: []string{
									"GET",
								},
							},
							MatchValue: "GET",
						},
						{
							MatchType:     "header",
							MatchOperator: "equals",
							CheckIPs:      "CONNECTING_IP",
						},
					},
				},
			},
			withError: `
MatchRules[0]: {
	Matches[0]: {
		MatchValue: must be blank when ObjectMatchValue is set
		ObjectMatchValue: must be blank when MatchValue is set
	}
	Matches[1]: {
		MatchValue: cannot be blank when ObjectMatchValue is blank
		ObjectMatchValue: cannot be blank when MatchValue is blank
	}
}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.input.Validate()
			if test.withError != "" {
				require.Error(t, err)
				assert.Equal(t, strings.TrimPrefix(test.withError, "\n"), err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
