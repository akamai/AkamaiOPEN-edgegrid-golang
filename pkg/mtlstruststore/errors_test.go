package mtlstruststore

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
		"Bad request 400": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
    "contextInfo": {
        "parameterName": "caSetId"
    },
    "detail": "The type of Parameter caSetId is invalid.",
    "status": 400,
    "title": "Type mismatch for any of path variable or query param or both.",
    "type": "/mtls-edge-truststore/error-types/path-variable-query-param-type-mismatch"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:   "/mtls-edge-truststore/error-types/path-variable-query-param-type-mismatch",
				Title:  "Type mismatch for any of path variable or query param or both.",
				Detail: "The type of Parameter caSetId is invalid.",
				Status: http.StatusBadRequest,
				ContextInfo: map[string]any{
					"parameterName": "caSetId",
				},
			},
		},
		"Invalid request 400": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
    "contextInfo": {
        "invalidParameterValue": "v2-api-create-ca set",
        "parameterName": "caSetName"
    },
    "errors": [
        {
            "contextInfo": {},
            "detail": "Provided CA set name v2-api-create-ca set does not match validation constraints. Allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-), percent (%) and period (.) with no three consecutive periods (…). Length must be between 3 and 64 characters.",
            "pointer": "/caSetName"
        }
    ],
    "status": 400,
    "title": "Invalid field value.",
    "type": "/mtls-edge-truststore/error-types/invalid-field"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:   "/mtls-edge-truststore/error-types/invalid-field",
				Title:  "Invalid field value.",
				Status: http.StatusBadRequest,
				ContextInfo: map[string]any{
					"invalidParameterValue": "v2-api-create-ca set",
					"parameterName":         "caSetName",
				},
				Errors: []ErrorItem{
					{
						Detail:      "Provided CA set name v2-api-create-ca set does not match validation constraints. Allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-), percent (%) and period (.) with no three consecutive periods (…). Length must be between 3 and 64 characters.",
						Pointer:     "/caSetName",
						ContextInfo: map[string]any{},
					},
				},
			},
		},
		"CASet does not exists 404": {
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(strings.NewReader(
					`{
    "contextInfo": {
        "caSetId": 0
    },
    "detail": "Cannot get CA set as the CA set with caSetId 0 is not found.",
    "status": 404,
    "title": "CA set is not found.",
    "type": "/mtls-edge-truststore/error-types/ca-set-not-found"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:   "/mtls-edge-truststore/error-types/ca-set-not-found",
				Title:  "CA set is not found.",
				Detail: "Cannot get CA set as the CA set with caSetId 0 is not found.",
				Status: http.StatusNotFound,
				ContextInfo: map[string]any{
					"caSetId": float64(0),
				},
			},
		},
		"invalid response body, assign status code": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:  "Failed to unmarshal error body. mTLS Truststore API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*mtlstruststore).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestIs(t *testing.T) {
	tests := map[string]struct {
		err      Error
		target   Error
		expected bool
	}{
		"different error code": {
			err:      Error{Status: 404},
			target:   Error{Status: 401},
			expected: false,
		},
		"same error code": {
			err:      Error{Status: 404},
			target:   Error{Status: 404},
			expected: true,
		},
		"same error code and title": {
			err:      Error{Status: 404, Title: "some error"},
			target:   Error{Status: 404, Title: "some error"},
			expected: true,
		},
		"same error code and different error message": {
			err:      Error{Status: 404, Title: "some error"},
			target:   Error{Status: 404, Title: "other error"},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.err.Is(&test.target), test.expected)
		})
	}
}
