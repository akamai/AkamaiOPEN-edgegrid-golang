package cloudlets

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestUnmarshalJSONMatchRules(t *testing.T) {
	tests := map[string]struct {
		withError      error
		responseBody   string
		expectedObject MatchRules
	}{
		"valid MatchRuleALB": {
			responseBody: `
	[
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
							MatchType:     "range",
							Negate:        false,
							ObjectMatchValue: &ObjectMatchValueRange{
								Type:  "range",
								Value: []int64{1, 50},
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
					},
					Name:  "Rule3",
					Start: 0,
				},
			},
		},

		"invalid objectMatchValue type for ALB - foo": {
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaALB: objectMatchValue has unexpected type: 'foo'"),
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
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaALB: 'type' should be a string"),
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
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaALB: objectMatchValue should contain 'type' field"),
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
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaALB: structure of objectMatchValue should be 'map', but was 'string'"),
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

		"invalid MatchRuleAP": {
			responseBody: `
	[
        {
            "type": "apMatchRule"
        }
    ]
`,
			withError: errors.New("unmarshalling MatchRules: unsupported match rule type: apMatchRule"),
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
            "type": "albMatchRule"
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

		"invalid objectMatchValue type for VP - range": {
			withError: errors.New("unmarshalling MatchRules: unmarshalling MatchCriteriaVP: objectMatchValue has unexpected type: 'range'"),
			responseBody: `
	[
        {
            "type": "vpMatchRule",
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

		"valid MatchRuleVP": {
			responseBody: `
	[
        {
            "type": "vpMatchRule",
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
				&MatchRuleVP{
					Type:               "vpMatchRule",
					End:                0,
					PassThroughPercent: 50.50,
					ID:                 0,
					MatchURL:           "",
					Matches: []MatchCriteriaVP{
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
