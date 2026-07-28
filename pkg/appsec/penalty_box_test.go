package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test PenaltyBox
func TestAppSec_GetPenaltyBox(t *testing.T) {

	result := GetPenaltyBoxResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestPenaltyBoxes/PenaltyBox.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetPenaltyBoxRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPenaltyBoxResponse
		withError        error
	}{
		"200 OK": {
			params: GetPenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetPenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching match target"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box",
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
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetPenaltyBox(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update PenaltyBox.
func TestAppSec_UpdatePenaltyBox(t *testing.T) {
	result := UpdatePenaltyBoxResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestPenaltyBoxes/PenaltyBox.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	req := UpdatePenaltyBoxRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestPenaltyBoxes/PenaltyBox.json"))
	err = json.Unmarshal([]byte(reqData), &req)
	require.NoError(t, err)

	tests := map[string]struct {
		params              UpdatePenaltyBoxRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *UpdatePenaltyBoxResponse
		withError           error
		headers             http.Header
		expectedRequestBody string
	}{
		"200 Success": {
			params: UpdatePenaltyBoxRequest{
				ConfigID:             43253,
				Version:              15,
				PolicyID:             "AAAA_81230",
				PenaltyBoxProtection: true,
				Action:               string(ActionTypeDeny),
			},
			expectedRequestBody: `{"action":"deny","penaltyBoxProtection":true}`,
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box",
		},
		"500 internal server error": {
			params: UpdatePenaltyBoxRequest{
				ConfigID:             43253,
				Version:              15,
				PolicyID:             "AAAA_81230",
				PenaltyBoxProtection: true,
				Action:               string(ActionTypeDeny),
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating zone"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box",
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
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
			result, err := client.UpdatePenaltyBox(
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

func TestGetPenaltyBoxRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req       GetPenaltyBoxRequest
		withError string
	}{
		"valid request - all fields populated": {
			req: GetPenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
		},
		"missing all required fields": {
			req:       GetPenaltyBoxRequest{},
			withError: "ConfigID: cannot be blank; PolicyID: cannot be blank; Version: cannot be blank.",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != "" {
				assert.EqualError(t, err, test.withError)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdatePenaltyBoxRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		req       UpdatePenaltyBoxRequest
		withError string
	}{
		"valid request - all fields populated": {
			req: UpdatePenaltyBoxRequest{
				ConfigID:             43253,
				Version:              15,
				PolicyID:             "AAAA_81230",
				Action:               string(ActionTypeDeny),
				PenaltyBoxProtection: true,
			},
		},
		"valid request - action alert": {
			req: UpdatePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Action:   string(ActionTypeAlert),
			},
		},
		"valid request - action none": {
			req: UpdatePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Action:   string(ActionTypeNone),
			},
		},
		"valid request - deny_custom_ ": {
			req: UpdatePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Action:   "deny_custom_1012254",
			},
		},
		"invalid action": {
			req: UpdatePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Action:   "invalid_action",
			},
			withError: "Action: value 'invalid_action' is invalid. Must be one of: alert, deny, deny_custom_{custom_deny_id}, none.",
		},
		"invalid action - deny_custom_ with blank id": {
			req: UpdatePenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
				Action:   "deny_custom_",
			},
			withError: "Action: value 'deny_custom_' is invalid. Must be one of: alert, deny, deny_custom_{custom_deny_id}, none.",
		},
		"missing all required fields": {
			req:       UpdatePenaltyBoxRequest{},
			withError: "Action: cannot be blank; ConfigID: cannot be blank; PolicyID: cannot be blank; Version: cannot be blank.",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != "" {
				assert.EqualError(t, err, test.withError)
				return
			}
			require.NoError(t, err)
		})
	}
}
