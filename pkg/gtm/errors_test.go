package gtm

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/stretchr/testify/require"

	"github.com/tj/assert"
)

func TestJSONErrorUnmarshalling(t *testing.T) {
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
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       ioutil.NopCloser(strings.NewReader(`<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>`))},
			expected: &Error{
				Type:       "",
				Title:      "Failed to unmarshal error body. GTM API failed. Check details for more information.",
				Detail:     "<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>",
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		"API failure with plain text response": {
			input: &http.Response{
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       ioutil.NopCloser(strings.NewReader("Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z"))},
			expected: &Error{
				Type:       "",
				Title:      "Failed to unmarshal error body. GTM API failed. Check details for more information.",
				Detail:     "Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z",
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		"API failure with XML response": {
			input: &http.Response{
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       ioutil.NopCloser(strings.NewReader(`<Root><Item id="1" name="Example" /></Root>`))},
			expected: &Error{
				Type:       "",
				Title:      "Failed to unmarshal error body. GTM API failed. Check details for more information.",
				Detail:     "<Root><Item id=\"1\" name=\"Example\" /></Root>",
				StatusCode: http.StatusServiceUnavailable,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			g := gtm{
				Session: sess,
			}
			assert.Equal(t, test.expected, g.Error(test.input))
		})
	}
}
