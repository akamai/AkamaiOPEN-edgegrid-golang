package domainownership

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
		"Bad request 400 - missing parameter": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
    "detail": "Required parameter 'validationScope' is missing.",
    "status": 400,
    "title": "Bad Request",
    "type": "bad-request"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:   "bad-request",
				Title:  "Bad Request",
				Detail: "Required parameter 'validationScope' is missing.",
				Status: http.StatusBadRequest,
			},
		},
		"Bad request 400 - invalid check": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
    "type": "bad-request",
    "title": "Bad Request",
    "instance": "d87ca377-f864-4054-bfa8-f567c44567c5",
    "status": 400,
    "detail": "Oops, something wasn't right. Please correct the errors.",
    "errors": [
        {
            "type": "error-types/invalid",
            "title": "Invalid Check",
            "detail": "Domain cannot be invalidated for the current state.",
            "field": "domains[0].domainName",
			"problemId": "d87ca377-f864-4054-bfa8-f567c44567c5"
        }
    ],
	"problemId": "d87ca377-f864-4054-bfa8-f567c44567c5"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "bad-request",
				Title:    "Bad Request",
				Instance: "d87ca377-f864-4054-bfa8-f567c44567c5",
				Detail:   "Oops, something wasn't right. Please correct the errors.",
				Status:   http.StatusBadRequest,
				Errors: []ErrorDetail{
					{
						Type:      "error-types/invalid",
						Title:     "Invalid Check",
						Detail:    "Domain cannot be invalidated for the current state.",
						Field:     "domains[0].domainName",
						ProblemID: "d87ca377-f864-4054-bfa8-f567c44567c5",
					},
				},
				ProblemID: "d87ca377-f864-4054-bfa8-f567c44567c5",
			},
		},
		"Bad request 400 - invalid value": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
    "type": "bad-request",
    "title": "Bad Request",
    "instance": "f5334872-80ae-437c-89ed-fee729f3a8de",
    "status": 400,
    "detail": "Invalid value 'a' for query parameter validationScope.",
    "errors": [
        {
            "type": "error-types/invalid",
            "title": "Invalid Value",
            "detail": "Invalid value 'a' for query parameter validationScope.",
            "field": "validationScope",
            "problemId": "f5334872-80ae-437c-89ed-fee729f3a8de"
        }
    ],
    "problemId": "f5334872-80ae-437c-89ed-fee729f3a8de"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "bad-request",
				Title:    "Bad Request",
				Instance: "f5334872-80ae-437c-89ed-fee729f3a8de",
				Detail:   "Invalid value 'a' for query parameter validationScope.",
				Status:   http.StatusBadRequest,
				Errors: []ErrorDetail{
					{
						Type:      "error-types/invalid",
						Title:     "Invalid Value",
						Detail:    "Invalid value 'a' for query parameter validationScope.",
						Field:     "validationScope",
						ProblemID: "f5334872-80ae-437c-89ed-fee729f3a8de",
					},
				},
				ProblemID: "f5334872-80ae-437c-89ed-fee729f3a8de",
			},
		},
		"Resource not found 404": {
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(strings.NewReader(
					`{
	"type": "not-found",
	"title": "Not Found",
	"instance": "fe111e63-225d-45ea-8e0a-dd182496092d",
	"status": 404,
	"detail": "The requested resource could not be found on the server.",
	"field": "domainName",
	"value": "example.com"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Title:    "Not Found",
				Type:     "not-found",
				Detail:   "The requested resource could not be found on the server.",
				Status:   http.StatusNotFound,
				Instance: "fe111e63-225d-45ea-8e0a-dd182496092d",
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
				Title:  "Failed to unmarshal error body. Domain Ownership Manager API failed. Check details for more information.",
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
				Title:  "Failed to unmarshal error body. Domain Ownership Manager API failed. Check details for more information.",
				Detail: "",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*domainownership).Error(tc.response)
			assert.Equal(t, tc.expected, res)
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
