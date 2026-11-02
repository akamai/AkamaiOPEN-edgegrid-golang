package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_GetClientsideProtectionConfiguration(t *testing.T) {
	result := GetClientsideProtectionConfigurationResponse{}
	data := compactJSON(loadFixtureBytes("testdata/TestClientsideProtectionConfiguration/GetClientsideProtectionConfiguration.json"))
	err := json.Unmarshal([]byte(data), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetClientsideProtectionConfigurationRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetClientsideProtectionConfigurationResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetClientsideProtectionConfigurationRequest{
				ConfigID:         43253,
				Version:          7,
				SecurityPolicyID: "test_policy_id",
			},
			responseStatus:   http.StatusOK,
			responseBody:     data,
			expectedPath:     "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			expectedResponse: &result,
		},
		"200 OK with scriptTagAttributes": {
			params: GetClientsideProtectionConfigurationRequest{
				ConfigID:         43253,
				Version:          7,
				SecurityPolicyID: "test_policy_id",
			},
			responseStatus: http.StatusOK,
			responseBody:   compactJSON(loadFixtureBytes("testdata/TestClientsideProtectionConfiguration/GetClientsideProtectionConfigurationWithScriptTagAttributes.json")),
			expectedPath:   "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			expectedResponse: func() *GetClientsideProtectionConfigurationResponse {
				var r GetClientsideProtectionConfigurationResponse
				require.NoError(t, json.Unmarshal(loadFixtureBytes("testdata/TestClientsideProtectionConfiguration/GetClientsideProtectionConfigurationWithScriptTagAttributes.json"), &r))
				return &r
			}(),
		},
		"500 internal server error": {
			params: GetClientsideProtectionConfigurationRequest{
				ConfigID:         43253,
				Version:          7,
				SecurityPolicyID: "test_policy_id",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error500), "want: %s; got: %s", error500, err)
			},
		},
		"400 bad request": {
			params: GetClientsideProtectionConfigurationRequest{
				ConfigID:         43253,
				Version:          7,
				SecurityPolicyID: "test_policy_id",
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"validate - missing ConfigID": {
			params: GetClientsideProtectionConfigurationRequest{
				Version:          7,
				SecurityPolicyID: "test_policy_id",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank.", err.Error())
			},
		},
		"validate - missing Version": {
			params: GetClientsideProtectionConfigurationRequest{
				ConfigID:         43253,
				SecurityPolicyID: "test_policy_id",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank.", err.Error())
			},
		},
		"validate - missing SecurityPolicyID": {
			params: GetClientsideProtectionConfigurationRequest{
				ConfigID: 43253,
				Version:  7,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: SecurityPolicyID: cannot be blank.", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := newMockAPIClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			})
			result, err := client.GetClientsideProtectionConfiguration(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppSec_UpdateClientsideProtectionConfiguration(t *testing.T) {
	result := UpdateClientsideProtectionConfigurationResponse{}
	data := compactJSON(loadFixtureBytes("testdata/TestClientsideProtectionConfiguration/UpdateClientsideProtectionConfiguration.json"))
	err := json.Unmarshal([]byte(data), &result)
	require.NoError(t, err)

	params := UpdateClientsideProtectionConfigurationRequest{
		ConfigID:         43253,
		Version:          7,
		SecurityPolicyID: "test_policy_id",
		Body: UpdateClientsideProtectionConfigurationRequestBody{
			EdgeInjection: ClientsideProtectionEdgeInjectionRequest{
				LoadScriptAsync: true,
			},
			EdgeTestParameters: ClientsideProtectionEdgeTestParameters{
				DisableInjectionKey:   "disable-key",
				DisableInjectionValue: "disable-value",
				ForceInjectionKey:     "force-key",
				ForceInjectionValue:   "force-value",
			},
			InjectionCriteria: ClientsideProtectionInjectionCriteria{
				InjectionPathIncludes: ptr.To("/include-path"),
				InjectionPathExcludes: ptr.To("/exclude-path"),
				InjectionPercent:      100,
				InjectionPolicy:       ptr.To("head"),
			},
		},
	}

	baseRequestBody := `{
		"edgeInjection":{"loadScriptAsync":true},
		"edgeTestParameters":{
			"disableInjectionKey":"disable-key",
			"disableInjectionValue":"disable-value",
			"forceInjectionKey":"force-key",
			"forceInjectionValue":"force-value"
		},
		"injectionCriteria":{
			"injectionPathIncludes":"/include-path",
			"injectionPathExcludes":"/exclude-path",
			"injectionPercent":100,
			"injectionPolicy":"head"
		}
	}`

	tests := map[string]struct {
		params              UpdateClientsideProtectionConfigurationRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *UpdateClientsideProtectionConfigurationResponse
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			params:              params,
			responseStatus:      http.StatusOK,
			responseBody:        data,
			expectedPath:        "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			expectedRequestBody: baseRequestBody,
			expectedResponse:    &result,
		},
		"201 Created": {
			params:              params,
			responseStatus:      http.StatusCreated,
			responseBody:        data,
			expectedPath:        "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			expectedRequestBody: baseRequestBody,
			expectedResponse:    &result,
		},
		"200 OK with scriptTagAttributes": {
			params: UpdateClientsideProtectionConfigurationRequest{
				ConfigID:         43253,
				Version:          7,
				SecurityPolicyID: "test_policy_id",
				Body: UpdateClientsideProtectionConfigurationRequestBody{
					EdgeInjection: ClientsideProtectionEdgeInjectionRequest{
						LoadScriptAsync: true,
					},
					EdgeTestParameters: ClientsideProtectionEdgeTestParameters{
						DisableInjectionKey:   "disable-key",
						DisableInjectionValue: "disable-value",
						ForceInjectionKey:     "force-key",
						ForceInjectionValue:   "force-value",
					},
					InjectionCriteria: ClientsideProtectionInjectionCriteria{
						InjectionPathIncludes: ptr.To("/include-path"),
						InjectionPathExcludes: ptr.To("/exclude-path"),
						InjectionPercent:      100,
						InjectionPolicy:       ptr.To("head"),
						ScriptTagAttributes:   []string{`nonce="abc123"`, `data-integrity="sha256"`},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody:   compactJSON(loadFixtureBytes("testdata/TestClientsideProtectionConfiguration/GetClientsideProtectionConfigurationWithScriptTagAttributes.json")),
			expectedPath:   "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			expectedRequestBody: `{
				"edgeInjection":{"loadScriptAsync":true},
				"edgeTestParameters":{
					"disableInjectionKey":"disable-key",
					"disableInjectionValue":"disable-value",
					"forceInjectionKey":"force-key",
					"forceInjectionValue":"force-value"
				},
				"injectionCriteria":{
					"injectionPathIncludes":"/include-path",
					"injectionPathExcludes":"/exclude-path",
					"injectionPercent":100,
					"injectionPolicy":"head",
					"scriptTagAttributes":["nonce=\"abc123\"","data-integrity=\"sha256\""]
				}
			}`,
			expectedResponse: func() *UpdateClientsideProtectionConfigurationResponse {
				var r UpdateClientsideProtectionConfigurationResponse
				require.NoError(t, json.Unmarshal(loadFixtureBytes("testdata/TestClientsideProtectionConfiguration/GetClientsideProtectionConfigurationWithScriptTagAttributes.json"), &r))
				return &r
			}(),
		},
		"500 internal server error": {
			params:         params,
			responseStatus: http.StatusInternalServerError,
			responseBody:   internalServerError,
			expectedPath:   "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error500), "want: %s; got: %s", error500, err)
			},
		},
		"400 bad request": {
			params:         params,
			responseStatus: http.StatusBadRequest,
			responseBody:   badRequest,
			expectedPath:   "/appsec/v1/configs/43253/versions/7/security-policies/test_policy_id/client-side-protection",
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, error400), "want: %s; got: %s", error400, err)
			},
		},
		"validate - missing ConfigID": {
			params: UpdateClientsideProtectionConfigurationRequest{
				Version:          7,
				SecurityPolicyID: "test_policy_id",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank.", err.Error())
			},
		},
		"validate - missing Version": {
			params: UpdateClientsideProtectionConfigurationRequest{
				ConfigID:         43253,
				SecurityPolicyID: "test_policy_id",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank.", err.Error())
			},
		},
		"validate - missing SecurityPolicyID": {
			params: UpdateClientsideProtectionConfigurationRequest{
				ConfigID: 43253,
				Version:  7,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: SecurityPolicyID: cannot be blank.", err.Error())
			},
		},
		"validate - InjectionPercent out of range": {
			params: UpdateClientsideProtectionConfigurationRequest{
				ConfigID:         43253,
				Version:          7,
				SecurityPolicyID: "test_policy_id",
				Body: UpdateClientsideProtectionConfigurationRequestBody{
					InjectionCriteria: ClientsideProtectionInjectionCriteria{
						InjectionPercent: 101,
					},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Body.InjectionCriteria.InjectionPercent: must be no greater than 100.", err.Error())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := newMockAPIClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			})
			result, err := client.UpdateClientsideProtectionConfiguration(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
