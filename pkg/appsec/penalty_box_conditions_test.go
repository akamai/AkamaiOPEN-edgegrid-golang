package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListPenaltyBoxConditions(t *testing.T) {

	result := GetPenaltyBoxConditionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestPenaltyBoxConditions/PenaltyBoxConditions.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	tests := map[string]struct {
		params           GetPenaltyBoxConditionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPenaltyBoxConditionsResponse
		withError        func(*testing.T, error)
		headers          http.Header
	}{
		"validation errors - PolicyID cannot be blank": {
			params: GetPenaltyBoxConditionsRequest{
				ConfigID: 43253,
				Version:  15,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validation errors - Version cannot be blank": {
			params: GetPenaltyBoxConditionsRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validation errors - ConfigID cannot be blank": {
			params: GetPenaltyBoxConditionsRequest{
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"200 OK": {
			params: GetPenaltyBoxConditionsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box/conditions",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetPenaltyBoxConditionsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error fetching PenaltyBoxConditions"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box/conditions",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching PenaltyBoxConditions",
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
			result, err := client.GetPenaltyBoxConditions(
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

func TestAppSec_UpdatePenaltyBoxConditions(t *testing.T) {
	result := UpdatePenaltyBoxConditionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestPenaltyBoxConditions/PenaltyBoxConditions.json"))
	err := json.Unmarshal([]byte(respData), &result)
	require.NoError(t, err)

	// eval Penalty Box conditions update test
	reqData := PenaltyBoxConditionsPayload{}
	err = json.Unmarshal(loadFixtureBytes("testdata/TestPenaltyBoxConditions/PenaltyBoxConditions.json"), &reqData)
	require.NoError(t, err)

	reqDataWithNoConditionOperator := PenaltyBoxConditionsPayload{
		ConditionOperator: "",
		Conditions:        &RuleConditions{},
	}

	reqDataWithNoConditions := PenaltyBoxConditionsPayload{
		ConditionOperator: "AND",
		Conditions:        nil,
	}

	tests := map[string]struct {
		params           UpdatePenaltyBoxConditionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdatePenaltyBoxConditionsResponse
		withError        func(*testing.T, error)
		headers          http.Header
	}{
		"validation errors - PolicyID cannot be empty string": {
			params: UpdatePenaltyBoxConditionsRequest{
				ConfigID:          43253,
				Version:           15,
				ConditionsPayload: reqData,
				PolicyID:          "",
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
			headers: nil,
		},
		"validation errors - PolicyID cannot be blank": {
			params: UpdatePenaltyBoxConditionsRequest{
				ConfigID:          43253,
				Version:           15,
				ConditionsPayload: reqData,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validation errors - ConfigID cannot be blank": {
			params: UpdatePenaltyBoxConditionsRequest{
				PolicyID:          "AAAA_81230",
				Version:           15,
				ConditionsPayload: reqData,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConfigID: cannot be blank", err.Error())
			},
		},
		"validation errors - Version cannot be blank": {
			params: UpdatePenaltyBoxConditionsRequest{
				ConfigID:          43253,
				PolicyID:          "AAAA_81230",
				ConditionsPayload: reqData,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Version: cannot be blank", err.Error())
			},
		},
		"validation errors - Request body payload cannot be blank": {
			params: UpdatePenaltyBoxConditionsRequest{
				ConfigID: 43253,
				PolicyID: "AAAA_81230",
				Version:  15,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating PenaltyBoxConditions"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box/conditions",
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConditionOperator: cannot be blank\nConditions: is required", err.Error())

			},
		},
		"validation errors - ConditionOperator cannot be blank": {
			params: UpdatePenaltyBoxConditionsRequest{
				ConfigID:          43253,
				Version:           15,
				PolicyID:          "AAAA_81230",
				ConditionsPayload: reqDataWithNoConditionOperator,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: ConditionOperator: cannot be blank", err.Error())
			},
		},
		"validation errors - Conditions cannot be blank": {
			params: UpdatePenaltyBoxConditionsRequest{
				ConfigID:          43253,
				Version:           15,
				PolicyID:          "AAAA_81230",
				ConditionsPayload: reqDataWithNoConditions,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "struct validation: Conditions: is required", err.Error())
			},
		},
		"200 Success": {
			params: UpdatePenaltyBoxConditionsRequest{
				ConfigID:          43253,
				Version:           15,
				PolicyID:          "AAAA_81230",
				ConditionsPayload: reqData,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box/conditions",
		},
		"500 internal server error": {
			params: UpdatePenaltyBoxConditionsRequest{
				ConfigID:          43253,
				Version:           15,
				PolicyID:          "AAAA_81230",
				ConditionsPayload: reqData,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
			{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error creating PenaltyBoxConditions"
			}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/penalty-box/conditions",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error creating PenaltyBoxConditions",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
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
			result, err := client.UpdatePenaltyBoxConditions(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
