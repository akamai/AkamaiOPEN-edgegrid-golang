package cps

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDeploymentSchedule(t *testing.T) {
	tests := map[string]struct {
		params           GetDeploymentScheduleRequest
		expectedPath     string
		expectedHeaders  map[string]string
		expectedResponse *DeploymentSchedule
		responseStatus   int
		responseBody     string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetDeploymentScheduleRequest{
				ChangeID:     1,
				EnrollmentID: 10,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
	"notAfter": "2021-11-03T08:02:46.655484Z",
	"notBefore": "2021-10-03T08:02:46.655484Z"
}`,
			expectedPath: "/cps/v2/enrollments/10/changes/1/deployment-schedule",
			expectedHeaders: map[string]string{
				"Accept": "application/vnd.akamai.cps.deployment-schedule.v1+json",
			},
			expectedResponse: &DeploymentSchedule{
				NotAfter:  tools.StringPtr("2021-11-03T08:02:46.655484Z"),
				NotBefore: tools.StringPtr("2021-10-03T08:02:46.655484Z"),
			},
		},
		"500 internal server error": {
			params: GetDeploymentScheduleRequest{
				ChangeID:     1,
				EnrollmentID: 10,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/10/changes/1/deployment-schedule",
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
		"validation error missing change_id": {
			params: GetDeploymentScheduleRequest{
				EnrollmentID: 10,
			},
			expectedPath: "/cps/v2/enrollments/10/changes/1/deployment-schedule",
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "ChangeID: cannot be blank.", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"validation error missing enorollment_id": {
			params: GetDeploymentScheduleRequest{
				ChangeID: 1,
			},
			expectedPath: "/cps/v2/enrollments/10/changes/1/deployment-schedule",
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "EnrollmentID: cannot be blank.", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				for k, v := range test.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetDeploymentSchedule(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateDeploymentSchedule(t *testing.T) {
	tests := map[string]struct {
		params              UpdateDeploymentScheduleRequest
		expectedPath        string
		expectedHeaders     map[string]string
		expectedRequestBody string
		expectedResponse    *UpdateDeploymentScheduleResponse
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
	}{
		"200 OK - update deployment schedule": {
			params: UpdateDeploymentScheduleRequest{
				ChangeID:     1,
				EnrollmentID: 10,
				DeploymentSchedule: DeploymentSchedule{
					NotAfter:  tools.StringPtr("2021-11-03T08:02:46.655484Z"),
					NotBefore: tools.StringPtr("2021-10-03T08:02:46.655484Z"),
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `{
	"change": "test_change"
}`,
			expectedRequestBody: `{"notAfter":"2021-11-03T08:02:46.655484Z","notBefore":"2021-10-03T08:02:46.655484Z"}`,
			expectedPath:        "/cps/v2/enrollments/10/changes/1/deployment-schedule",
			expectedHeaders: map[string]string{
				"Accept":       "application/vnd.akamai.cps.change-id.v1+json",
				"Content-type": "application/vnd.akamai.cps.deployment-schedule.v1+json; charset=utf-8",
			},
			expectedResponse: &UpdateDeploymentScheduleResponse{
				Change: "test_change",
			},
		},
		"500 internal server error": {
			params: UpdateDeploymentScheduleRequest{
				ChangeID:     1,
				EnrollmentID: 10,
				DeploymentSchedule: DeploymentSchedule{
					NotAfter:  tools.StringPtr("2021-11-03T08:02:46.655484Z"),
					NotBefore: tools.StringPtr("2021-10-03T08:02:46.655484Z"),
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/10/changes/1/deployment-schedule",
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
		"validation error missing change_id": {
			params: UpdateDeploymentScheduleRequest{
				EnrollmentID: 10,
				DeploymentSchedule: DeploymentSchedule{
					NotAfter:  tools.StringPtr("2021-11-03T08:02:46.655484Z"),
					NotBefore: tools.StringPtr("2021-10-03T08:02:46.655484Z"),
				},
			},
			expectedPath: "/cps/v2/enrollments/10/changes/1/deployment-schedule",
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "ChangeID: cannot be blank.", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
		"validation error missing enorollment_id": {
			params: UpdateDeploymentScheduleRequest{
				ChangeID: 1,
				DeploymentSchedule: DeploymentSchedule{
					NotAfter:  tools.StringPtr("2021-11-03T08:02:46.655484Z"),
					NotBefore: tools.StringPtr("2021-10-03T08:02:46.655484Z"),
				},
			},
			expectedPath: "/cps/v2/enrollments/10/changes/1/deployment-schedule",
			withError: func(t *testing.T, err error) {
				assert.Containsf(t, err.Error(), "EnrollmentID: cannot be blank.", "want: %s; got: %s", ErrStructValidation, err)
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				for k, v := range test.expectedHeaders {
					assert.Equal(t, v, r.Header.Get(k))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateDeploymentSchedule(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
