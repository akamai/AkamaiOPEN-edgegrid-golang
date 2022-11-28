package imaging

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
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
				Body: ioutil.NopCloser(strings.NewReader(
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
				Body: ioutil.NopCloser(strings.NewReader(
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
				Body: ioutil.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:  "test",
				Detail: "",
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
