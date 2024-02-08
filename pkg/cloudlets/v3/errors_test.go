package v3

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
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
		"valid response, status code 500": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body: ioutil.NopCloser(strings.NewReader(
					`{"type":"a","title":"b","detail":"c"}`),
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
				Body: ioutil.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:  "Failed to unmarshal error body. Cloudlets API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*cloudlets).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestAs(t *testing.T) {
	someErrorMarshalled, _ := json.Marshal("some error")
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
			err:      Error{Status: 404, Errors: someErrorMarshalled},
			target:   Error{Status: 404, Errors: someErrorMarshalled},
			expected: true,
		},
		"same error code and different error message": {
			err:      Error{Status: 404, Errors: someErrorMarshalled},
			target:   Error{Status: 404},
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
				Body:    ioutil.NopCloser(strings.NewReader(`<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>`))},
			expected: &Error{
				Type:   "",
				Title:  "Failed to unmarshal error body. Cloudlets API failed. Check details for more information.",
				Detail: "<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>",
			},
		},
		"API failure with plain text response": {
			input: &http.Response{
				Request: req,
				Status:  "OK",
				Body:    ioutil.NopCloser(strings.NewReader("Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z"))},
			expected: &Error{
				Type:   "",
				Title:  "Failed to unmarshal error body. Cloudlets API failed. Check details for more information.",
				Detail: "Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z",
			},
		},
		"API failure with XML response": {
			input: &http.Response{
				Request: req,
				Status:  "OK",
				Body:    ioutil.NopCloser(strings.NewReader(`<Root><Item id="1" name="Example" /></Root>`))},
			expected: &Error{
				Type:   "",
				Title:  "Failed to unmarshal error body. Cloudlets API failed. Check details for more information.",
				Detail: "<Root><Item id=\"1\" name=\"Example\" /></Root>",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			as := cloudlets{
				Session: sess,
			}
			assert.Equal(t, test.expected, as.Error(test.input))
		})
	}
}
