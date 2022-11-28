package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListRatePolicies(t *testing.T) {

	result := GetRatePoliciesResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePolicies.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetRatePoliciesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRatePoliciesResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetRatePoliciesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/rate-policies",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetRatePoliciesRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
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
			expectedPath: "/appsec/v1/configs/43253/versions/15/rate-policies",
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
			client := mockAPIClient(t, mockServer)
			result, err := client.GetRatePolicies(
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

// Test RatePolicy
func TestAppSec_GetRatePolicy(t *testing.T) {

	result := GetRatePolicyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePolicies.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetRatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetRatePolicyResponse
		withError        error
	}{
		"200 OK": {
			params: GetRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				RatePolicyID:  134644,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/rate-policies/134644",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				RatePolicyID:  134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/rate-policies/134644",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching match target",
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
			result, err := client.GetRatePolicy(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create RatePolicy
func TestAppSec_CreateRatePolicy(t *testing.T) {

	result := CreateRatePolicyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePolicies.json"))
	json.Unmarshal([]byte(respData), &result)

	req := CreateRatePolicyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePolicies.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           CreateRatePolicyRequest
		prop             *CreateRatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateRatePolicyResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/rate-policies",
		},
		"500 internal server error": {
			params: CreateRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating domain"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/rate-policies",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating domain",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateRatePolicy(
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

// Test Create RatePolicy with negative hostnames match (using RatePoliciesHosts field)
func TestAppSec_CreateRatePolicy_NegativeMatch(t *testing.T) {

	result := CreateRatePolicyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePoliciesHosts.json"))
	json.Unmarshal([]byte(respData), &result)

	req := CreateRatePolicyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePoliciesHosts.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           CreateRatePolicyRequest
		prop             *CreateRatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateRatePolicyResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/rate-policies",
		},
		"500 internal server error": {
			params: CreateRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating domain"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/rate-policies",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating domain",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateRatePolicy(
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

// Test Update RatePolicy
func TestAppSec_UpdateRatePolicy(t *testing.T) {
	result := UpdateRatePolicyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePolicies.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateRatePolicyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePolicies.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateRatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateRatePolicyResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				RatePolicyID:  134644,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/rate-policies/134644",
		},
		"500 internal server error": {
			params: UpdateRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				RatePolicyID:  134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/rate-policies/134644",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
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
			result, err := client.UpdateRatePolicy(
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

// Test Remove RatePolicy
func TestAppSec_RemoveRatePolicy(t *testing.T) {

	result := RemoveRatePolicyResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePoliciesEmpty.json"))
	json.Unmarshal([]byte(respData), &result)

	req := RemoveRatePolicyRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestRatePolicies/RatePoliciesEmpty.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           RemoveRatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemoveRatePolicyResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: RemoveRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				RatePolicyID:  134644,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/rate-policies/134644",
		},
		"500 internal server error": {
			params: RemoveRatePolicyRequest{
				ConfigID:      43253,
				ConfigVersion: 15,
				RatePolicyID:  134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error deleting match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/rate-policies/134644",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error deleting match target",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.RemoveRatePolicy(
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
