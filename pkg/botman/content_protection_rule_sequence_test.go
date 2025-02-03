package botman

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Get ContentProtectionRuleSequence
func TestBotman_GetContentProtectionRuleSequence(t *testing.T) {
	tests := map[string]struct {
		params           GetContentProtectionRuleSequenceRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetContentProtectionRuleSequenceResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetContentProtectionRuleSequenceRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"contentProtectionRuleSequence":["fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"]}`,
			expectedPath:   "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rule-sequence",
			expectedResponse: &GetContentProtectionRuleSequenceResponse{
				ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"},
			},
		},
		"500 internal server error": {
			params: GetContentProtectionRuleSequenceRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rule-sequence",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: GetContentProtectionRuleSequenceRequest{
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: GetContentProtectionRuleSequenceRequest{
				ConfigID:         43253,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: GetContentProtectionRuleSequenceRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
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
			result, err := client.GetContentProtectionRuleSequence(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update ContentProtectionRuleSequence.
func TestBotman_UpdateContentProtectionRuleSequence(t *testing.T) {
	tests := map[string]struct {
		params           UpdateContentProtectionRuleSequenceRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateContentProtectionRuleSequenceResponse
		withError        func(*testing.T, error)
	}{
		"200 Success": {
			params: UpdateContentProtectionRuleSequenceRequest{
				ConfigID:                      43253,
				Version:                       15,
				SecurityPolicyID:              "AAAA_81230",
				ContentProtectionRuleSequence: ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"}},
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"contentProtectionRuleSequence":["fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"]}`,
			expectedResponse: &UpdateContentProtectionRuleSequenceResponse{
				ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"},
			},
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rule-sequence",
		},
		"500 internal server error": {
			params: UpdateContentProtectionRuleSequenceRequest{
				ConfigID:                      43253,
				Version:                       15,
				SecurityPolicyID:              "AAAA_81230",
				ContentProtectionRuleSequence: ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"}},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error updating data"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/content-protection-rule-sequence",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error updating data",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"Missing ConfigID": {
			params: UpdateContentProtectionRuleSequenceRequest{
				Version:                       15,
				SecurityPolicyID:              "AAAA_81230",
				ContentProtectionRuleSequence: ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"}},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ConfigID")
			},
		},
		"Missing Version": {
			params: UpdateContentProtectionRuleSequenceRequest{
				ConfigID:                      43253,
				SecurityPolicyID:              "AAAA_81230",
				ContentProtectionRuleSequence: ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"}},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "Version")
			},
		},
		"Missing SecurityPolicyID": {
			params: UpdateContentProtectionRuleSequenceRequest{
				ConfigID:                      43253,
				Version:                       15,
				ContentProtectionRuleSequence: ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake9b8-4fd5-430e-a061-1c61df1d2ac2"}},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "SecurityPolicyID")
			},
		},
		"Missing ContentProtectionRuleSequence": {
			params: UpdateContentProtectionRuleSequenceRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContentProtectionRuleSequence")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateContentProtectionRuleSequence(
				session.ContextWithOptions(
					context.Background()), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
