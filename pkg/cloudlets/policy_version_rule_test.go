package cloudlets

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPolicyVersionRule(t *testing.T) {
	tests := map[string]struct {
		request          GetPolicyVersionRuleRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse MatchRule
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			request: GetPolicyVersionRuleRequest{
				AkaRuleID: "abc123",
				PolicyID:  12345,
				Version:   2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "type": "erMatchRule",
    "akaRuleId": "abc123",
    "end": 0,
    "id": 0,
    "location": "/cloudlets/api/v2/policies/12345/versions/2/rules/abc123",
    "matchURL": "/redirect-to-example",
    "matches": [
        {
            "caseSensitive": false,
            "matchOperator": "contains",
            "matchType": "header",
            "negate": false,
            "objectMatchValue": {
                "type": "object",
                "name": "contentType",
                "options": {
                    "value": [
                        "text/html*",
                        "text/css*",
                        "application/x-javascript*"
                    ],
                    "valueHasWildcard": true
                }
            }
        },
        {
            "caseSensitive": false,
            "matchOperator": "contains",
            "matchType": "header",
            "negate": false,
            "objectMatchValue": {
                "type": "object",
                "name": "contentType",
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
    "name": "ER RANGE",
    "redirectURL": "https://www.example.com",
    "start": 0,
    "statusCode": 302,
    "useIncomingQueryString": true,
    "useRelativeUrl": "none"
}`,
			expectedPath: "/cloudlets/api/v2/policies/12345/versions/2/rules/abc123",
			expectedResponse: &MatchRuleER{
				Name:  "ER RANGE",
				Type:  "erMatchRule",
				Start: 0,
				End:   0,
				ID:    0,
				Matches: []MatchCriteriaER{
					{
						MatchType:     "header",
						MatchOperator: "contains",
						ObjectMatchValue: &ObjectMatchValueObject{
							Name: "contentType",
							Type: Object,
							Options: &Options{
								Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
								ValueHasWildcard: true,
							},
						},
					},
					{
						MatchType:     "header",
						MatchOperator: "contains",
						ObjectMatchValue: &ObjectMatchValueObject{
							Name: "contentType",
							Type: Object,
							Options: &Options{
								Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
								ValueHasWildcard: true,
							},
						},
					},
				},
				MatchesAlways:            false,
				UseRelativeURL:           "none",
				StatusCode:               http.StatusFound,
				RedirectURL:              "https://www.example.com",
				MatchURL:                 "/redirect-to-example",
				UseIncomingQueryString:   true,
				UseIncomingSchemeAndHost: false,
				Disabled:                 false,
				AkaRuleID:                "abc123",
				Location:                 "/cloudlets/api/v2/policies/12345/versions/2/rules/abc123",
			},
		},

		"struct validation error - missed AkaRuleID": {
			request: GetPolicyVersionRuleRequest{
				PolicyID: 12345,
				Version:  2,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get policy version rule: struct validation:\nAkaRuleID: cannot be blank", err.Error())
			},
		},
		"struct validation error - missed PolicyID": {
			request: GetPolicyVersionRuleRequest{
				AkaRuleID: "abc123",
				Version:   2,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get policy version rule: struct validation:\nPolicyID: cannot be blank", err.Error())
			},
		},
		"struct validation error - missed Version": {
			request: GetPolicyVersionRuleRequest{
				AkaRuleID: "abc123",
				PolicyID:  12345,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get policy version rule: struct validation:\nVersion: cannot be blank", err.Error())
			},
		},
		"struct validation error - missed Version, PolicyID and AkaRuleID": {
			request: GetPolicyVersionRuleRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get policy version rule: struct validation:\nAkaRuleID: cannot be blank\nPolicyID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"struct validation error - invalid Version": {
			request: GetPolicyVersionRuleRequest{
				AkaRuleID: "abc123",
				Version:   -1,
				PolicyID:  12345,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get policy version rule: struct validation:\nVersion: must be no less than 1", err.Error())
			},
		},
		"500 internal server error": {
			request: GetPolicyVersionRuleRequest{
				AkaRuleID: "abc123",
				PolicyID:  12345,
				Version:   2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/12345/versions/2/rules/abc123",
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
			result, err := client.GetPolicyVersionRule(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreatePolicyVersionRule(t *testing.T) {
	tests := map[string]struct {
		request             CreatePolicyVersionRuleRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    MatchRule
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			request: CreatePolicyVersionRuleRequest{
				Version:  2,
				PolicyID: 12345,
				Index:    2,
				MatchRule: &MatchRuleER{
					Type:     MatchRuleTypeER,
					MatchURL: "/redirect-to-example",
					Matches: []MatchCriteriaER{
						{
							CaseSensitive: false,
							MatchOperator: MatchOperatorContains,
							MatchType:     "header",
							ObjectMatchValue: ObjectMatchValueObject{
								Type: Object,
								Name: "contentType",
								Options: &Options{
									Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
									ValueHasWildcard: true,
								},
							},
						},
						{
							CaseSensitive: false,
							MatchOperator: MatchOperatorContains,
							MatchType:     "header",
							ObjectMatchValue: ObjectMatchValueObject{
								Type: Object,
								Name: "contentType",
								Options: &Options{
									Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
									ValueHasWildcard: true,
								},
							},
						},
					},
				},
			},
			expectedRequestBody: `{"type":"erMatchRule","matches":[{"matchType":"header","matchOperator":"contains","caseSensitive":false,"negate":false,"objectMatchValue":{"name":"contentType","type":"object","nameCaseSensitive":false,"nameHasWildcard":false,"options":{"value":["text/html*","text/css*","application/x-javascript*"],"valueHasWildcard":true}}},{"matchType":"header","matchOperator":"contains","caseSensitive":false,"negate":false,"objectMatchValue":{"name":"contentType","type":"object","nameCaseSensitive":false,"nameHasWildcard":false,"options":{"value":["text/html*","text/css*","application/x-javascript*"],"valueHasWildcard":true}}}],"statusCode":0,"redirectURL":"","matchURL":"/redirect-to-example","useIncomingQueryString":false,"useIncomingSchemeAndHost":false}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
    "type": "erMatchRule",
    "akaRuleId": "abc123",
    "end": 0,
    "id": 0,
    "matchURL": "/redirect-to-example",
    "matches": [
        {
            "caseSensitive": false,
            "matchOperator": "contains",
            "matchType": "header",
            "negate": false,
            "objectMatchValue": {
                "type": "object",
                "name": "contentType",
                "options": {
                    "value": [
                        "text/html*",
                        "text/css*",
                        "application/x-javascript*"
                    ],
                    "valueHasWildcard": true
                }
            }
        },
        {
            "caseSensitive": false,
            "matchOperator": "contains",
            "matchType": "header",
            "negate": false,
            "objectMatchValue": {
                "type": "object",
                "name": "contentType",
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
    "name": "ER RANGE",
    "redirectURL": "https://www.example.com",
    "start": 0,
    "statusCode": 302,
    "useIncomingQueryString": true,
    "useRelativeUrl": "none"
}`,
			expectedPath: "/cloudlets/api/v2/policies/12345/versions/2/rules?index=2",
			expectedResponse: &MatchRuleER{
				Name: "ER RANGE",
				Type: "erMatchRule",
				End:  0,
				ID:   0,
				Matches: []MatchCriteriaER{
					{
						CaseSensitive: false,
						MatchType:     "header",
						MatchOperator: "contains",
						Negate:        false,
						ObjectMatchValue: &ObjectMatchValueObject{
							Name: "contentType",
							Type: Object,
							Options: &Options{
								Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
								ValueHasWildcard: true,
							},
						},
					},
					{
						CaseSensitive: false,
						MatchType:     "header",
						MatchOperator: "contains",
						Negate:        false,
						ObjectMatchValue: &ObjectMatchValueObject{
							Name: "contentType",
							Type: Object,
							Options: &Options{
								Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
								ValueHasWildcard: true,
							},
						},
					},
				},
				MatchesAlways:          false,
				UseRelativeURL:         "none",
				StatusCode:             http.StatusFound,
				RedirectURL:            "https://www.example.com",
				MatchURL:               "/redirect-to-example",
				UseIncomingQueryString: true,
				AkaRuleID:              "abc123",
			},
		},
		"struct validation error - missed PolicyID": {
			request: CreatePolicyVersionRuleRequest{
				Version: 2,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create policy version rule: struct validation:\nPolicyID: cannot be blank", err.Error())
			},
		},
		"struct validation error - missed Version": {
			request: CreatePolicyVersionRuleRequest{
				PolicyID: 12345,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create policy version rule: struct validation:\nVersion: cannot be blank", err.Error())
			},
		},
		"struct validation error - missed Version, PolicyID": {
			request: CreatePolicyVersionRuleRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create policy version rule: struct validation:\nPolicyID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"struct validation error - invalid Version": {
			request: CreatePolicyVersionRuleRequest{
				Version:  -2,
				PolicyID: 12345,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "create policy version rule: struct validation:\nVersion: must be no less than 1", err.Error())
			},
		},
		"500 internal server error": {
			request: CreatePolicyVersionRuleRequest{
				PolicyID: 12345,
				Version:  2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/12345/versions/2/rules",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreatePolicyVersionRule(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdatePolicyVersionRule(t *testing.T) {
	tests := map[string]struct {
		request             UpdatePolicyVersionRuleRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    MatchRule
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			request: UpdatePolicyVersionRuleRequest{
				Version:   2,
				PolicyID:  12345,
				AkaRuleID: "abc123",
				MatchRule: &MatchRuleER{
					Type:     MatchRuleTypeER,
					MatchURL: "/redirect-to-example",
					Matches: []MatchCriteriaER{
						{
							CaseSensitive: false,
							MatchOperator: MatchOperatorContains,
							MatchType:     "header",
							ObjectMatchValue: ObjectMatchValueObject{
								Type: Object,
								Name: "contentType",
								Options: &Options{
									Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
									ValueHasWildcard: true,
								},
							},
						},
						{
							CaseSensitive: false,
							MatchOperator: MatchOperatorContains,
							MatchType:     "header",
							ObjectMatchValue: ObjectMatchValueObject{
								Type: Object,
								Name: "contentType",
								Options: &Options{
									Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
									ValueHasWildcard: true,
								},
							},
						},
					},
				},
			},
			expectedRequestBody: `{"type":"erMatchRule","matches":[{"matchType":"header","matchOperator":"contains","caseSensitive":false,"negate":false,"objectMatchValue":{"name":"contentType","type":"object","nameCaseSensitive":false,"nameHasWildcard":false,"options":{"value":["text/html*","text/css*","application/x-javascript*"],"valueHasWildcard":true}}},{"matchType":"header","matchOperator":"contains","caseSensitive":false,"negate":false,"objectMatchValue":{"name":"contentType","type":"object","nameCaseSensitive":false,"nameHasWildcard":false,"options":{"value":["text/html*","text/css*","application/x-javascript*"],"valueHasWildcard":true}}}],"statusCode":0,"redirectURL":"","matchURL":"/redirect-to-example","useIncomingQueryString":false,"useIncomingSchemeAndHost":false}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
    "type": "erMatchRule",
    "akaRuleId": "abc123",
    "end": 0,
    "id": 0,
    "matchURL": "/redirect-to-example",
    "matches": [
        {
            "caseSensitive": true,
            "matchOperator": "contains",
            "matchType": "header",
            "negate": false,
            "objectMatchValue": {
                "type": "object",
                "name": "contentType",
                "options": {
                    "value": [
                        "text/html*",
                        "text/css*",
                        "application/x-javascript*"
                    ],
                    "valueHasWildcard": true
                }
            }
        },
        {
            "caseSensitive": false,
            "matchOperator": "contains",
            "matchType": "header",
            "negate": false,
            "objectMatchValue": {
                "type": "object",
                "name": "contentType",
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
    "name": "ER RANGE",
    "redirectURL": "https://www.example.com",
    "start": 0,
    "statusCode": 302,
    "useIncomingQueryString": true,
    "useRelativeUrl": "none"
}`,
			expectedPath: "/cloudlets/api/v2/policies/12345/versions/2/rules/abc123",
			expectedResponse: &MatchRuleER{
				Name:  "ER RANGE",
				Type:  "erMatchRule",
				Start: 0,
				End:   0,
				ID:    0,
				Matches: []MatchCriteriaER{
					{
						MatchType:     "header",
						MatchOperator: "contains",
						CaseSensitive: true,
						ObjectMatchValue: &ObjectMatchValueObject{
							Name: "contentType",
							Type: Object,
							Options: &Options{
								Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
								ValueHasWildcard: true,
							},
						},
					},
					{
						MatchType:     "header",
						MatchOperator: "contains",
						ObjectMatchValue: &ObjectMatchValueObject{
							Name: "contentType",
							Type: Object,
							Options: &Options{
								Value:            []string{"text/html*", "text/css*", "application/x-javascript*"},
								ValueHasWildcard: true,
							},
						},
					},
				},
				MatchesAlways:            false,
				UseRelativeURL:           "none",
				StatusCode:               http.StatusFound,
				RedirectURL:              "https://www.example.com",
				MatchURL:                 "/redirect-to-example",
				UseIncomingQueryString:   true,
				UseIncomingSchemeAndHost: false,
				Disabled:                 false,
				AkaRuleID:                "abc123",
			},
		},
		"struct validation error - missed PolicyID": {
			request: UpdatePolicyVersionRuleRequest{
				Version:   2,
				AkaRuleID: "abc123",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update policy version rule: struct validation:\nPolicyID: cannot be blank", err.Error())
			},
		},
		"struct validation error - missed Version": {
			request: UpdatePolicyVersionRuleRequest{
				PolicyID:  12345,
				AkaRuleID: "abc123",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update policy version rule: struct validation:\nVersion: cannot be blank", err.Error())
			},
		},
		"struct validation error - missed AkaRuleID": {
			request: UpdatePolicyVersionRuleRequest{
				PolicyID: 12345,
				Version:  2,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update policy version rule: struct validation:\nAkaRuleID: cannot be blank", err.Error())
			},
		},
		"struct validation error - missed PolicyID, Version, AkaRuleID": {
			request: UpdatePolicyVersionRuleRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update policy version rule: struct validation:\nAkaRuleID: cannot be blank\nPolicyID: cannot be blank\nVersion: cannot be blank", err.Error())
			},
		},
		"struct validation error - invalid Version": {
			request: UpdatePolicyVersionRuleRequest{
				PolicyID:  12345,
				Version:   -2,
				AkaRuleID: "abc123",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "update policy version rule: struct validation:\nVersion: must be no less than 1", err.Error())
			},
		},
		"500 internal server error": {
			request: UpdatePolicyVersionRuleRequest{
				PolicyID:  12345,
				Version:   2,
				AkaRuleID: "abc123",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/12345/versions/2/rules/abc123",
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
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdatePolicyVersionRule(context.Background(), test.request)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
