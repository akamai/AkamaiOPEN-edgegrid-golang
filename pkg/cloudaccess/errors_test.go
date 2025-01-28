package cloudaccess

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
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
				Status:     "Internal Server Error",
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
	"type": "bad-request",
	"title": "Bad Request",
	"instance": "test-instance-123",
	"status": 400
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "bad-request",
				Title:    "Bad Request",
				Instance: "test-instance-123",
				Status:   http.StatusBadRequest,
			},
		},
		"Invalid request 400": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
	"type": "invalid-request",
	"title": "Invalid Request",
	"instance": "test-instance-123",
	"status": 400,
	"errors": [
		{
		  "detail": "Constraint violation: accessKeyName must not be blank.",
		  "title": "Constraint Violation",
		  "type": "/cam/error-types/constraint-violation"
		},
		{
		  "detail": "Constraint violation: accessKeyName length must be between 1 and 50.",
		  "title": "Constraint Violation",
		  "type": "/cam/error-types/constraint-violation"
		}
	]
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "invalid-request",
				Title:    "Invalid Request",
				Instance: "test-instance-123",
				Status:   http.StatusBadRequest,
				Errors: []ErrorItem{
					{
						Detail: "Constraint violation: accessKeyName must not be blank.",
						Title:  "Constraint Violation",
						Type:   "/cam/error-types/constraint-violation",
					},
					{
						Detail: "Constraint violation: accessKeyName length must be between 1 and 50.",
						Title:  "Constraint Violation",
						Type:   "/cam/error-types/constraint-violation",
					},
				},
			},
		},
		"access key does not exists 404": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(strings.NewReader(
					`{
	"type": "/cam/error-types/access-key-does-not-exist",
	"title": "Domain Error",
	"detail": "Access key with accessKeyUID '1' does not exist.",
	"instance": "test-instance-123",
	"status": 404,
	"accessKeyUid": 1
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:         accessKeyNotFoundType,
				Title:        "Domain Error",
				Detail:       "Access key with accessKeyUID '1' does not exist.",
				Instance:     "test-instance-123",
				Status:       http.StatusNotFound,
				AccessKeyUID: 1,
			},
		},
		"invalid response body, assign status code": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:  "Failed to unmarshal error body. Cloud Access Manager API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*cloudaccess).Error(test.response)
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
