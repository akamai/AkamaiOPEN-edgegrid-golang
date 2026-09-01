package cloudcertificates

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	t.Parallel()
	sess, err := session.New()
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodHead,
		"/",
		nil)
	require.NoError(t, err)

	tests := map[string]struct {
		response *http.Response
		expected *Error
	}{
		"Bad request 400 - invalid field value": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`
					{
						"detail": "Invalid value '{grp_1234}' for field '{groupId}'. Failed to convert value of type 'String' to required type 'Integer'; For input string: \"grp_1234\"",
						"status": 400,
						"title": "Invalid field value.",
						"type": "/error-types/invalid-field",
						"instance": "/error-types/invalid-field?traceId=12345",
						"explanation": "Failed to convert value of type 'String' to required type 'Integer'; For input string: \"grp_1234\"",
						"parameterName": "groupId",
						"invalidParameterValue": "grp_1234"
					}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:                  "/error-types/invalid-field",
				Title:                 "Invalid field value.",
				Detail:                "Invalid value '{grp_1234}' for field '{groupId}'. Failed to convert value of type 'String' to required type 'Integer'; For input string: \"grp_1234\"",
				Status:                http.StatusBadRequest,
				Instance:              "/error-types/invalid-field?traceId=12345",
				Explanation:           "Failed to convert value of type 'String' to required type 'Integer'; For input string: \"grp_1234\"",
				ParameterName:         "groupId",
				InvalidParameterValue: "grp_1234",
			},
		},
		"Resource not found 404": {
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(strings.NewReader(`
					{
						"type": "/error-types/certificate-not-found",
						"title": "Certificate subscription is not found.",
						"instance": "/error-types/certificate-not-found?traceId=12345",
						"status": 404,
						"detail": "Certificate subscription with {certificateSubscriptionId}: {1234} is not found.",
						"certificateIdentifier": "certificateSubscriptionId",
						"certificateIdentifierValue": "1234"
					}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:                       "/error-types/certificate-not-found",
				Title:                      "Certificate subscription is not found.",
				Instance:                   "/error-types/certificate-not-found?traceId=12345",
				Detail:                     "Certificate subscription with {certificateSubscriptionId}: {1234} is not found.",
				Status:                     http.StatusNotFound,
				CertificateIdentifier:      "certificateSubscriptionId",
				CertificateIdentifierValue: "1234",
			},
		},
		"Validation error 400 - invalid parameter value": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`
					{
						"type": "/error-types/validation-failure",
						"title": "Validation failure.",
						"instance": "/error-types/validation-failure?traceId=12345",
						"status": 400,
						"detail": "Validation failed while executing the operation.",
						"errors": [
							{
								"type": "/error-types/invalid-field",
								"title": "Invalid field value.",
								"detail": "Invalid value '{[example.com]}' for field '{sans}'. SANs list cannot contain duplicates.",
								"instance": "/error-types/validation-failure?traceId=12345",
								"explanation": "SANs list cannot contain duplicates.",
								"invalidParameterValue": ["example.com"],
								"parameterName": "sans"
							}
						]
					}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "/error-types/validation-failure",
				Title:    "Validation failure.",
				Instance: "/error-types/validation-failure?traceId=12345",
				Detail:   "Validation failed while executing the operation.",
				Status:   http.StatusBadRequest,
				Errors: []SecondaryError{
					{
						Type:                  "/error-types/invalid-field",
						Title:                 "Invalid field value.",
						Detail:                "Invalid value '{[example.com]}' for field '{sans}'. SANs list cannot contain duplicates.",
						Instance:              "/error-types/validation-failure?traceId=12345",
						Explanation:           "SANs list cannot contain duplicates.",
						InvalidParameterValue: []string{"example.com"},
						ParameterName:         "sans",
					},
				},
			},
		},
		"Invalid response body, assign status code": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:  "Failed to unmarshal error body. CCM API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
		"Empty response body, assign status code": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			},
			expected: &Error{
				Title:  "Failed to unmarshal error body. CCM API failed. Check details for more information.",
				Detail: "",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*cloudcertificates).Error(tc.response)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestIs(t *testing.T) {
	t.Parallel()
	someError := Error{Type: "/some/type", Title: "some error", Status: 404, Detail: "some detail", Instance: "/some/error/instance"}

	tests := map[string]struct {
		target   Error
		expected bool
	}{
		"different error status": {
			target:   Error{Status: 401},
			expected: false,
		},
		"different error title": {
			target:   Error{Title: "other error"},
			expected: false,
		},
		"same error title": {
			target:   Error{Title: "some error"},
			expected: true,
		},
		"same error type": {
			target:   Error{Type: "/some/type"},
			expected: true,
		},
		"same error status": {
			target:   Error{Status: 404},
			expected: true,
		},
		"same error type, title and status": {
			target:   Error{Type: "/some/type", Title: "some error", Status: 404},
			expected: true,
		},
		"same error type but different title": {
			target:   Error{Type: "/some/type", Title: "other error"},
			expected: false,
		},
		"same error status and title but different type": {
			target:   Error{Type: "/other/type", Title: "some error", Status: 404},
			expected: false,
		},
		"same error status and type but different detail and instance": {
			target:   Error{Type: "/some/type", Status: 404, Detail: "other detail", Instance: "/other/error/instance"},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, someError.Is(&test.target), test.expected)
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	e := &Error{
		Type:     "/error-types/test",
		Title:    "Test Error",
		Status:   400,
		Detail:   "This is a test error",
		Instance: "/error-types/test?traceId=12345",
	}
	expected := `API error: 
{
	"type": "/error-types/test",
	"title": "Test Error",
	"status": 400,
	"detail": "This is a test error",
	"instance": "/error-types/test?traceId=12345"
}`
	assert.EqualError(t, e, expected)
}
