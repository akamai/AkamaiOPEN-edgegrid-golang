package imaging

import (
	"context"
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
		"valid response, status code 400, Bad Request": {
			response: &http.Response{
				Status:     "Bad Request",
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{"type":"testType","title":"Bad Request","detail":"error","status":400,
					"extensionFields":{"requestId":"123"},"problemId":"abc123","requestId":"123"}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:   "testType",
				Title:  "Bad Request",
				Detail: "error",
				Status: http.StatusBadRequest,
				ExtensionFields: map[string]string{
					"requestId": "123",
				},
				ProblemID: "abc123",
				RequestID: "123",
			},
		},
		"valid response, status code 400, Illegal parameter value": {
			response: &http.Response{
				Status:     "Bad Request",
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{"type":"testType","title":"Illegal parameter value","detail":"error","status":400,
					"extensionFields":{"illegalValue":"abc","parameterName":"param1"},"problemId":"abc123","illegalValue":"abc","parameterName":"param1"}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:   "testType",
				Title:  "Illegal parameter value",
				Detail: "error",
				Status: http.StatusBadRequest,
				ExtensionFields: map[string]string{
					"illegalValue":  "abc",
					"parameterName": "param1",
				},
				ProblemID:     "abc123",
				IllegalValue:  "abc",
				ParameterName: "param1",
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
				Title:  "Failed to unmarshal error body. Image & Video Manager API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*imaging).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestAs(t *testing.T) {
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
		"same error code and error message": {
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

func TestJsonErrorUnmarshalling(t *testing.T) {
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
				Title:  "Failed to unmarshal error body. Image & Video Manager API failed. Check details for more information.",
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
				Title:  "Failed to unmarshal error body. Image & Video Manager API failed. Check details for more information.",
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
				Title:  "Failed to unmarshal error body. Image & Video Manager API failed. Check details for more information.",
				Detail: "<Root><Item id=\"1\" name=\"Example\" /></Root>",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			g := imaging{
				Session: sess,
			}
			assert.Equal(t, test.expected, g.Error(test.input))
		})
	}
}
