package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppsecGetAdvancedSettingsRequestBody(t *testing.T) {

	result := GetAdvancedSettingsRequestBodyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsRequestBodyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsRequestBodyResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/request-body",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
					"type": "internal_error",
					"title": "Internal Server Error",
					"detail": "Error fetching AdvancedSettingsRequestBody",
					"status": 500
				}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/request-body",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsRequestBody",
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
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAdvancedSettingsRequestBody(
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

func TestAppsecGetAdvancedSettingsRequestBodyValidation(t *testing.T) {

	result := GetAdvancedSettingsRequestBodyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsRequestBodyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsRequestBodyResponse
		withError        func(*testing.T, error)
		headers          http.Header
	}{
		"validation error configID missing": {
			params: GetAdvancedSettingsRequestBodyRequest{
				Version: 15,
			},
			headers:      http.Header{},
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/request-body",
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "ConfigID: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"validation error version missing": {
			params: GetAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
			},
			headers:      http.Header{},
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/request-body",
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "Version: cannot be blank", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
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
			result, err := client.GetAdvancedSettingsRequestBody(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppsecGetAdvancedSettingsRequestBodyPolicy(t *testing.T) {

	result := GetAdvancedSettingsRequestBodyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicy.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetAdvancedSettingsRequestBodyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetAdvancedSettingsRequestBodyResponse
		withError        error
	}{
		"200 OK": {
			params: GetAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "test_policy",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/request-body",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "test_policy",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
							{
								"type": "internal_error",
								"title": "Internal Server Error",
								"detail": "Error fetching AdvancedSettingsRequestBody"
							}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/request-body",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching AdvancedSettingsRequestBody",
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
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAdvancedSettingsRequestBody(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update AdvancedSettingsRequestBody config level.
func TestAppsecUpdateAdvancedSettingsRequestBody(t *testing.T) {
	result := UpdateAdvancedSettingsRequestBodyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsRequestBodyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsRequestBodyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsRequestBodyResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/advanced-settings/request-body",
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsRequestBody"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/request-body",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsRequestBody",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAdvancedSettingsRequestBody(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update UpdateAdvancedSettingsRequestBody policy level.
func TestAppsecUpdateAdvancedSettingsRequestBodyPolicy(t *testing.T) {
	result := UpdateAdvancedSettingsRequestBodyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicy.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsRequestBodyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicy.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsRequestBodyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsRequestBodyResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "test_policy",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/request-body",
		},
		"400 invalid input error": {
			params: UpdateAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
			{
    			"detail": "The value of the request body size parameter must be one of [default, 8, 16, 32]",
    			"title": "Invalid Input Error",
				"type": "internal_error"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/request-body",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Invalid Input Error",
				Detail:     "The value of the request body size parameter must be one of [default, 8, 16, 32]",
				StatusCode: http.StatusBadRequest,
			},
		},
		"500 internal server error": {
			params: UpdateAdvancedSettingsRequestBodyRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "test_policy",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating AdvancedSettingsRequestBody"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/request-body",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating AdvancedSettingsRequestBody",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAdvancedSettingsRequestBody(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppsecUpdateAdvancedSettingsRequestBodyPolicyWithInvalidValue(t *testing.T) {
	result := UpdateAdvancedSettingsRequestBodyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicyWithInvalidValue.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdateAdvancedSettingsRequestBodyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicyWithInvalidValue.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           UpdateAdvancedSettingsRequestBodyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateAdvancedSettingsRequestBodyResponse
		withError        error
		headers          http.Header
	}{
		"400 invalid input error": {
			params: UpdateAdvancedSettingsRequestBodyRequest{
				ConfigID:                           43253,
				Version:                            15,
				RequestBodyInspectionLimitInKB:     req.RequestBodyInspectionLimitInKB,
				RequestBodyInspectionLimitOverride: req.RequestBodyInspectionLimitOverride,
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
			{
    			"detail": "The value of the request body size parameter must be one of [default, 8, 16, 32]",
    			"title": "Invalid Input Error",
				"type": "internal_error"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/advanced-settings/request-body",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Invalid Input Error",
				Detail:     "The value of the request body size parameter must be one of [default, 8, 16, 32]",
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAdvancedSettingsRequestBody(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestAppsecUpdateAdvancedSettingsRequestBodyPolicyWithOverrideUnset(t *testing.T) {
	result := RemoveAdvancedSettingsRequestBodyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicyWithOverrideUnsetResponse.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := RemoveAdvancedSettingsRequestBodyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicyWithOverrideUnsetRequest.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params           RemoveAdvancedSettingsRequestBodyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveAdvancedSettingsRequestBodyResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveAdvancedSettingsRequestBodyRequest{
				ConfigID:                           43253,
				Version:                            15,
				PolicyID:                           "test_policy",
				RequestBodyInspectionLimitInKB:     req.RequestBodyInspectionLimitInKB,
				RequestBodyInspectionLimitOverride: req.RequestBodyInspectionLimitOverride,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/test_policy/advanced-settings/request-body",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.RemoveAdvancedSettingsRequestBody(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
