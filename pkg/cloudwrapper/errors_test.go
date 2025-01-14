package cloudwrapper

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
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
  "instance": "30109837-7ea6-4b14-a41d-50cfb12a4b03",
  "status": 400,
  "detail": "Erroneous data input",
  "errors": [
    {
      "type": "bad-request",
      "title": "Bad Request",
      "detail": "Configuration with name UpdateConfiguration already exists in account 1234-3KNWKV.",
      "illegalValue": "UpdateConfiguration",
      "illegalParameter": "configurationName"
    },
    {
      "type": "bad-request",
      "title": "Bad Request",
      "detail": "One or more ARL Property is already used in another configuration.",
      "illegalValue": [
        {
          "propertyId": "123010"
        }
      ],
      "illegalParameter": "properties"
    }
  ]
}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "bad-request",
				Title:    "Bad Request",
				Instance: "30109837-7ea6-4b14-a41d-50cfb12a4b03",
				Status:   http.StatusBadRequest,
				Detail:   "Erroneous data input",
				Errors: []ErrorItem{
					{
						Type:             "bad-request",
						Title:            "Bad Request",
						Detail:           "Configuration with name UpdateConfiguration already exists in account 1234-3KNWKV.",
						IllegalValue:     "UpdateConfiguration",
						IllegalParameter: "configurationName",
					},
					{
						Type:             "bad-request",
						Title:            "Bad Request",
						Detail:           "One or more ARL Property is already used in another configuration.",
						IllegalValue:     []any{map[string]any{"propertyId": "123010"}},
						IllegalParameter: "properties",
					},
				},
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
				Title:  "Failed to unmarshal error body. Cloud Wrapper API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*cloudwrapper).Error(test.response)
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
				Title:  "Failed to unmarshal error body. Cloud Wrapper API failed. Check details for more information.",
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
				Title:  "Failed to unmarshal error body. Cloud Wrapper API failed. Check details for more information.",
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
				Title:  "Failed to unmarshal error body. Cloud Wrapper API failed. Check details for more information.",
				Detail: "<Root><Item id=\"1\" name=\"Example\" /></Root>",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			c := cloudwrapper{
				Session: sess,
			}
			assert.Equal(t, test.expected, c.Error(test.input))
		})
	}
}
