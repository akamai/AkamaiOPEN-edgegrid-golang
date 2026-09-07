package edgeworkers

import (
	"context"
	"errors"
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
		"valid response, status code 500": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(
					`{"type":"a","title":"b","detail":"c","status":500}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:   "a",
				Title:  "b",
				Detail: "c",
				Status: http.StatusInternalServerError,
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
				Title:  "Failed to unmarshal error body. Edgeworkers API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*edgeworkers).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestIs(t *testing.T) {
	tests := map[string]struct {
		err      error
		target   error
		expected bool
	}{
		"different error code": {
			err:      &Error{Status: 404},
			target:   &Error{Status: 401},
			expected: false,
		},
		"same error code": {
			err:      &Error{Status: 404},
			target:   &Error{Status: 404},
			expected: true,
		},
		"same error code and title": {
			err:      &Error{Status: 404, Title: "some error"},
			target:   &Error{Status: 404, Title: "some error"},
			expected: true,
		},
		"same error code and different error message": {
			err:      &Error{Status: 404, Title: "some error"},
			target:   &Error{Status: 404, Title: "other error"},
			expected: false,
		},
		"ErrNotFound matched by status 404 and EKV_9000": {
			err:      &Error{Status: http.StatusNotFound, ErrorCode: "EKV_9000"},
			target:   ErrNotFound,
			expected: true,
		},
		"ErrNotFound not matched when status is 400": {
			err:      &Error{Status: http.StatusBadRequest, ErrorCode: "EKV_9000"},
			target:   ErrNotFound,
			expected: false,
		},
		"ErrNamespaceNoScheduledDelete matched by status 400, EKV_9000 and matching detail": {
			err: &Error{
				Status:    http.StatusBadRequest,
				ErrorCode: "EKV_9000",
				Detail:    "No current scheduled delete for DevExpAutomatedTest-c7OIFc. Must first schedule the delete before modifying it.",
			},
			target:   ErrNamespaceNoScheduledDelete,
			expected: true,
		},
		"ErrNamespaceNoScheduledDelete not matched when status is 404": {
			err: &Error{
				Status:    http.StatusNotFound,
				ErrorCode: "EKV_9000",
				Detail:    "No current scheduled delete for some-namespace.",
			},
			target:   ErrNamespaceNoScheduledDelete,
			expected: false,
		},
		"ErrNamespaceNoScheduledDelete not matched when error code differs": {
			err: &Error{
				Status:    http.StatusBadRequest,
				ErrorCode: "EKV_0000",
				Detail:    "No current scheduled delete for some-namespace.",
			},
			target:   ErrNamespaceNoScheduledDelete,
			expected: false,
		},
		"ErrNamespaceNoScheduledDelete not matched when detail does not contain expected substring": {
			err: &Error{
				Status:    http.StatusBadRequest,
				ErrorCode: "EKV_9000",
				Detail:    "The requested namespace does not exist.",
			},
			target:   ErrNamespaceNoScheduledDelete,
			expected: false,
		},
		"ErrNamespaceNotFound matched by status 400, EKV_9000 and exact detail": {
			err: &Error{
				Status:    http.StatusBadRequest,
				ErrorCode: "EKV_9000",
				Detail:    "The requested namespace does not exist.",
			},
			target:   ErrNamespaceNotFound,
			expected: true,
		},
		"ErrNamespaceNotFound not matched when status is 404": {
			err: &Error{
				Status:    http.StatusNotFound,
				ErrorCode: "EKV_9000",
				Detail:    "The requested namespace does not exist.",
			},
			target:   ErrNamespaceNotFound,
			expected: false,
		},
		"ErrNamespaceNotFound not matched when error code differs": {
			err: &Error{
				Status:    http.StatusBadRequest,
				ErrorCode: "EKV_0000",
				Detail:    "The requested namespace does not exist.",
			},
			target:   ErrNamespaceNotFound,
			expected: false,
		},
		"ErrNamespaceNotFound not matched when detail differs": {
			err: &Error{
				Status:    http.StatusBadRequest,
				ErrorCode: "EKV_9000",
				Detail:    "No current scheduled delete for some-namespace.",
			},
			target:   ErrNamespaceNotFound,
			expected: false,
		},
		"ErrVersionBeingDeactivated matched by error code EW1031": {
			err:      &Error{ErrorCode: "EW1031"},
			target:   ErrVersionBeingDeactivated,
			expected: true,
		},
		"ErrVersionAlreadyDeactivated matched by error code EW1032": {
			err:      &Error{ErrorCode: "EW1032"},
			target:   ErrVersionAlreadyDeactivated,
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, errors.Is(test.err, test.target))
		})
	}
}

func TestValidationErrorsParsing(t *testing.T) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodHead,
		"/",
		nil)
	require.NoError(t, err)
	tests := map[string]struct {
		input    *http.Response
		expected *Error
	}{
		"API failure with HTML response": {
			input: &http.Response{
				Request: req,
				Status:  "OK",
				Body:    io.NopCloser(strings.NewReader(`<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>`))},
			expected: &Error{
				Type:   "",
				Title:  "Failed to unmarshal error body. Edgeworkers API failed. Check details for more information.",
				Detail: "<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>",
			},
		},
		"API failure with plain text response": {
			input: &http.Response{
				Request: req,
				Status:  "OK",
				Body:    io.NopCloser(strings.NewReader("Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z"))},
			expected: &Error{
				Type:   "",
				Title:  "Failed to unmarshal error body. Edgeworkers API failed. Check details for more information.",
				Detail: "Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z",
			},
		},
		"API failure with XML response": {
			input: &http.Response{
				Request: req,
				Status:  "OK",
				Body:    io.NopCloser(strings.NewReader(`<Root><Item id="1" name="Example" /></Root>`))},
			expected: &Error{
				Type:   "",
				Title:  "Failed to unmarshal error body. Edgeworkers API failed. Check details for more information.",
				Detail: "<Root><Item id=\"1\" name=\"Example\" /></Root>",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			e := edgeworkers{
				Session: sess,
			}
			assert.Equal(t, test.expected, e.Error(test.input))
		})
	}
}
