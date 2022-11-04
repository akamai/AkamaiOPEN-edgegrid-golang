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

// Test get EvalPenaltyBox
func TestAppSec_GetEvalPenaltyBox(t *testing.T) {

	result := GetPenaltyBoxResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestPenaltyBoxes/PenaltyBox.json"))
	json.Unmarshal([]byte(respData), &result)

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
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-penalty-box",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetPenaltyBoxRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-penalty-box",
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
			result, err := client.GetEvalPenaltyBox(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update EvalPenaltyBox.
func TestAppSec_UpdateEvalPenaltyBox(t *testing.T) {
	result := UpdatePenaltyBoxResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestPenaltyBoxes/PenaltyBox.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdatePenaltyBoxRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestPenaltyBoxes/PenaltyBox.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdatePenaltyBoxRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdatePenaltyBoxResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdatePenaltyBoxRequest{
				ConfigID:             43253,
				Version:              15,
				PolicyID:             "AAAA_81230",
				PenaltyBoxProtection: false,
				Action:               string(ActionTypeDeny),
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-penalty-box",
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
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/eval-penalty-box",
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
			result, err := client.UpdateEvalPenaltyBox(
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
