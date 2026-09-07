package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAdvancedSettingsURLEvasionDefenseRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req       GetAdvancedSettingsURLEvasionDefenseRequest
		withError func(*testing.T, error)
	}{
		"valid": {
			req: GetAdvancedSettingsURLEvasionDefenseRequest{
				ConfigID: 43253,
				Version:  15,
			},
		},
		"missing config id": {
			req: GetAdvancedSettingsURLEvasionDefenseRequest{
				Version: 15,
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "ConfigID: cannot be blank")
			},
		},
		"missing version": {
			req: GetAdvancedSettingsURLEvasionDefenseRequest{
				ConfigID: 43253,
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Version: cannot be blank")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != nil {
				require.Error(t, err)
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdateAdvancedSettingsURLEvasionDefenseRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req       UpdateAdvancedSettingsURLEvasionDefenseRequest
		withError func(*testing.T, error)
	}{
		"valid": {
			req: UpdateAdvancedSettingsURLEvasionDefenseRequest{
				ConfigID: 43253,
				Version:  15,
			},
		},
		"missing config id": {
			req: UpdateAdvancedSettingsURLEvasionDefenseRequest{
				Version: 15,
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "ConfigID: cannot be blank")
			},
		},
		"missing version": {
			req: UpdateAdvancedSettingsURLEvasionDefenseRequest{
				ConfigID: 43253,
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Version: cannot be blank")
			},
		},
		"invalid condition operator in rules": {
			req: UpdateAdvancedSettingsURLEvasionDefenseRequest{
				ConfigID: 43253,
				Version:  15,
				Body: UpdateAdvancedSettingsURLEvasionDefenseRequestBody{
					Rules: []AdvancedSettingsURLEvasionDefenseRuleRequest{
						{
							RuleID:            3002500,
							ConditionOperator: (*AdvancedSettingsURLEvasionDefenseConditionOperator)(ptr.To("INVALID")),
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "ConditionOperator: value 'INVALID' is invalid. Must be one of: 'AND', 'OR'")
			},
		},
		"invalid condition field for type": {
			req: UpdateAdvancedSettingsURLEvasionDefenseRequest{
				ConfigID: 43253,
				Version:  15,
				Body: UpdateAdvancedSettingsURLEvasionDefenseRequestBody{
					Rules: []AdvancedSettingsURLEvasionDefenseRuleRequest{
						{
							RuleID: 3002500,
							Action: "alert",
							Conditions: []AdvancedSettingsURLEvasionDefenseRuleCondition{
								{
									Type:    "hostMatch",
									Methods: []string{"GET"},
								},
							},
						},
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Regexp(t, `(?s)Conditions\[0\]:\s*\{.*Methods: field not supported for type 'hostMatch'; allowed fields: \[Hosts Type PositiveMatch\]`, err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != nil {
				require.Error(t, err)
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAdvancedSettingsURLEvasionDefenseRuleConditions_Validate(t *testing.T) {
	tests := map[string]struct {
		conditions AdvancedSettingsURLEvasionDefenseRuleCondition
		withError  func(*testing.T, error)
	}{
		"valid request header match": {
			conditions: AdvancedSettingsURLEvasionDefenseRuleCondition{
				Type:          "requestHeaderMatch",
				Header:        "X-Test",
				Value:         "abc",
				ValueCase:     true,
				ValueWildcard: true,
			},
		},
		"unsupported field for type": {
			conditions: AdvancedSettingsURLEvasionDefenseRuleCondition{
				Type:    "extensionMatch",
				Methods: []string{"POST"},
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Methods: field not supported for type 'extensionMatch'")
			},
		},
		"missing type": {
			conditions: AdvancedSettingsURLEvasionDefenseRuleCondition{},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Type: cannot be blank")
			},
		},
		"unknown type ignores allowed fields validation": {
			conditions: AdvancedSettingsURLEvasionDefenseRuleCondition{
				Type:    "futureType",
				Methods: []string{"POST"},
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Type: value 'futureType' is invalid. Must be one of: 'clientListMatch', 'extensionMatch', 'filenameMatch', 'hostMatch', 'ipMatch', 'pathMatch', 'requestHeaderMatch', 'requestMethodMatch', 'uriQueryMatch'")
				assert.NotContains(t, err.Error(), "Methods: field not supported")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.conditions.Validate()
			if test.withError != nil {
				require.Error(t, err)
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAppSec_GetAdvancedSettingsURLEvasionDefense(t *testing.T) {
	result := GetAdvancedSettingsURLEvasionDefenseResponse{}
	data := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsUrlEvasionDefense/AdvancedSettingsUrlEvasionDefense.json"))
	err := json.Unmarshal([]byte(data), &result)
	require.NoError(t, err)
	tests := map[string]struct {
		params           GetAdvancedSettingsURLEvasionDefenseRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsURLEvasionDefenseResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetAdvancedSettingsURLEvasionDefenseRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus:   http.StatusOK,
			responseBody:     data,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/url-evasion-defense",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsURLEvasionDefenseRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching URL evasion defense",
				"status": 500
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/url-evasion-defense",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching URL evasion defense",
					StatusCode: http.StatusInternalServerError,
				}), "want matching API error, got: %v", err)
			},
		},
		"validation error": {
			params: GetAdvancedSettingsURLEvasionDefenseRequest{
				Version: 15,
			},
			withError: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "ConfigID: cannot be blank")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.expectedPath == "" {
					t.Fatalf("unexpected request for validation-only test: %s", r.URL.String())
				}
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAdvancedSettingsURLEvasionDefense(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateAdvancedSettingsURLEvasionDefense(t *testing.T) {
	result := UpdateAdvancedSettingsURLEvasionDefenseResponse{}
	data := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsUrlEvasionDefense/AdvancedSettingsUrlEvasionDefense.json"))
	err := json.Unmarshal([]byte(data), &result)
	require.NoError(t, err)

	requestBody := UpdateAdvancedSettingsURLEvasionDefenseRequestBody{}
	err = json.Unmarshal([]byte(data), &requestBody)
	require.NoError(t, err)

	request := UpdateAdvancedSettingsURLEvasionDefenseRequest{
		ConfigID: 43253,
		Version:  15,
		Body:     requestBody,
	}
	marshaledRequestBody, err := json.Marshal(request.Body)
	require.NoError(t, err)

	tests := map[string]struct {
		params              UpdateAdvancedSettingsURLEvasionDefenseRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *UpdateAdvancedSettingsURLEvasionDefenseResponse
		expectedRequestBody string
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			params:              request,
			responseStatus:      http.StatusOK,
			responseBody:        data,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/advanced-settings/url-evasion-defense?disableLegacyEvasivePathMatch=false",
			expectedResponse:    &result,
			expectedRequestBody: string(marshaledRequestBody),
		},
		"500 internal server error": {
			params:         request,
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error updating URL evasion defense",
				"status": 500
			}`,
			expectedPath:        "/appsec/v1/configs/43253/versions/15/advanced-settings/url-evasion-defense?disableLegacyEvasivePathMatch=false",
			expectedRequestBody: string(marshaledRequestBody),
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error updating URL evasion defense",
					StatusCode: http.StatusInternalServerError,
				}), "want matching API error, got: %v", err)
			},
		},
		"validation error": {
			params: UpdateAdvancedSettingsURLEvasionDefenseRequest{
				Version: 15,
			},
			withError: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "ConfigID: cannot be blank")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.expectedPath == "" {
					t.Fatalf("unexpected request for validation-only test: %s", r.URL.String())
				}
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAdvancedSettingsURLEvasionDefense(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
