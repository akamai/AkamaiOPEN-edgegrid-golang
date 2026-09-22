package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListExportConfiguration(t *testing.T) {

	result := GetExportConfigurationResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestExportConfiguration/ExportConfiguration.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	aiRulesResult := GetExportConfigurationResponse{}
	aiRulesRespData := compactJSON(loadFixtureBytes("testdata/TestAIRule/ExportAIRules.json"))
	err = json.Unmarshal([]byte(aiRulesRespData), &aiRulesResult)
	require.NoError(t, err)

	aiRulesDisabledResult := GetExportConfigurationResponse{}
	aiRulesDisabledRespData := compactJSON(loadFixtureBytes("testdata/TestExportConfiguration/ExportAIRulesDisabled.json"))
	err = json.Unmarshal([]byte(aiRulesDisabledRespData), &aiRulesDisabledResult)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetExportConfigurationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetExportConfigurationResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetExportConfigurationRequest{
				ConfigID: 43253,
				Version:  15,
				Source:   "TF",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/export/configs/43253/versions/15?source=TF",
			expectedResponse: &result,
		},
		"200 OK - with AI rules data": {
			params: GetExportConfigurationRequest{
				ConfigID: 77653,
				Version:  25,
				Source:   "TF",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(aiRulesRespData),
			expectedPath:     "/appsec/v1/export/configs/77653/versions/25?source=TF",
			expectedResponse: &aiRulesResult,
		},
		"200 OK - DISABLED AI rules with empty list": {
			params: GetExportConfigurationRequest{
				ConfigID: 77653,
				Version:  25,
				Source:   "TF",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(aiRulesDisabledRespData),
			expectedPath:     "/appsec/v1/export/configs/77653/versions/25?source=TF",
			expectedResponse: &aiRulesDisabledResult,
		},
		"500 internal server error": {
			params: GetExportConfigurationRequest{
				ConfigID: 43253,
				Version:  15,
				Source:   "TF",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching propertys",
    "status": 500
}`,
			expectedPath: "/appsec/v1/export/configs/43253/versions/15?source=TF",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching propertys",
				StatusCode: http.StatusInternalServerError,
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetExportConfiguration(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// TestAppSec_ListExportConfigurationRapidRules asserts that the rapidRules block of the export
// configuration fixture actually round-trips into the SecurityPolicies.RapidRules struct.
// TestAppSec_ListExportConfiguration compares the api response against a struct unmarshaled from
// the same fixture, so it would still pass if the RapidRules field did not exist at all.
func TestAppSec_ListExportConfigurationRapidRules(t *testing.T) {
	result := GetExportConfigurationResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestExportConfiguration/ExportConfiguration.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	require.NotEmpty(t, result.SecurityPolicies)
	policy := result.SecurityPolicies[0]
	require.Equal(t, "AAAA_81230", policy.ID)

	rapidRules := policy.RapidRules
	require.NotNil(t, rapidRules, "rapidRules must be deserialized into the security policy")
	assert.Equal(t, "alert", rapidRules.DefaultAction)
	assert.True(t, rapidRules.Enabled)
	assert.Equal(t, "on", rapidRules.ThreatIntel)

	require.Len(t, rapidRules.PolicyRules, 2)

	assert.Equal(t, "deny", rapidRules.PolicyRules[0].Action)
	assert.False(t, rapidRules.PolicyRules[0].Lock)
	assert.Equal(t, 1001, rapidRules.PolicyRules[0].RuleID)
	assert.Equal(t, 1, rapidRules.PolicyRules[0].RuleVersion)
	assert.Equal(t, "SQL", rapidRules.PolicyRules[0].Group)
	assert.Equal(t, 7392, rapidRules.PolicyRules[0].RulesetVersionID)

	conditionException := rapidRules.PolicyRules[0].ConditionException
	require.NotNil(t, conditionException, "conditionException must be deserialized into the policy rule")

	require.NotNil(t, conditionException.Conditions)
	require.Len(t, *conditionException.Conditions, 1)
	assert.Equal(t, "extensionMatch", (*conditionException.Conditions)[0].Type)
	assert.Equal(t, []string{"php"}, (*conditionException.Conditions)[0].Extensions)
	assert.True(t, (*conditionException.Conditions)[0].PositiveMatch)

	require.NotNil(t, conditionException.Exception)
	require.NotNil(t, conditionException.Exception.SpecificHeaderCookieParamXMLOrJSONNames)
	exceptionNames := *conditionException.Exception.SpecificHeaderCookieParamXMLOrJSONNames
	require.Len(t, exceptionNames, 1)
	assert.Equal(t, []string{"Cookie1001"}, exceptionNames[0].Names)
	assert.Equal(t, "REQUEST_COOKIES_NAMES", exceptionNames[0].Selector)
	assert.True(t, exceptionNames[0].Wildcard)

	require.NotNil(t, conditionException.AdvancedExceptionsList)
	assert.Equal(t, "AND", conditionException.AdvancedExceptionsList.ConditionOperator)
	require.NotNil(t, conditionException.AdvancedExceptionsList.SpecificHeaderCookieParamXMLOrJSONNames)
	advancedNames := *conditionException.AdvancedExceptionsList.SpecificHeaderCookieParamXMLOrJSONNames
	require.Len(t, advancedNames, 1)
	assert.Equal(t, []string{"adv1001"}, advancedNames[0].Names)
	assert.Equal(t, "REQUEST_COOKIES", advancedNames[0].Selector)

	assert.Equal(t, "alert", rapidRules.PolicyRules[1].Action)
	assert.True(t, rapidRules.PolicyRules[1].Lock)
	assert.Equal(t, 1002, rapidRules.PolicyRules[1].RuleID)
	assert.Equal(t, 2, rapidRules.PolicyRules[1].RuleVersion)
	assert.Equal(t, "XSS", rapidRules.PolicyRules[1].Group)
	assert.Equal(t, 7392, rapidRules.PolicyRules[1].RulesetVersionID)
	assert.Nil(t, rapidRules.PolicyRules[1].ConditionException, "conditionException is optional and must stay nil when absent")
}

// TestAppSec_ListExportConfigurationRapidRulesAbsent asserts that RapidRules stays nil for
// policies which have no rapidRules block, since the cli-terraform export uses a nil RapidRules
// as the signal that rapid rules are not configured for the policy.
func TestAppSec_ListExportConfigurationRapidRulesAbsent(t *testing.T) {
	result := GetExportConfigurationResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAIRule/ExportAIRules.json"))
	require.NoError(t, json.Unmarshal([]byte(respData), &result))

	require.NotEmpty(t, result.SecurityPolicies)
	for _, policy := range result.SecurityPolicies {
		assert.Nil(t, policy.RapidRules, "policy %s should have no rapid rules", policy.ID)
	}
}
