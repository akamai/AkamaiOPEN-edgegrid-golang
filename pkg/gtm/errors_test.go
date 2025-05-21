package gtm

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader(`<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>`))},
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
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader("Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z"))},
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
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader(`<Root><Item id="1" name="Example" /></Root>`))},
			expected: &Error{
				Type:       "",
				Title:      "Failed to unmarshal error body. GTM API failed. Check details for more information.",
				Detail:     "<Root><Item id=\"1\" name=\"Example\" /></Root>",
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		"API failure nested error": {
			input: &http.Response{
				Request:    req,
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`
{
 	"type": "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationFailed",
 	"title": "Property Validation Failure",
 	"detail": "",
 	"instance": "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm-api/v1/domains/ddzh-test-1.akadns.net/properties/property_test#d290ddf7-53da-4509-be5a-ba582614f883",
 	"errors": [
 		{
 			"type": "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError",
 			"title": "Property Validation Error",
 			"detail": "In Property \"property_test\", there are no enabled traffic targets that have any traffic allowed to go to them",
 			"errors": null
 		}
 	]
}`))},
			expected: &Error{
				Type:       "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationFailed",
				Title:      "Property Validation Failure",
				Detail:     "",
				StatusCode: http.StatusBadRequest,
				Instance:   "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm-api/v1/domains/ddzh-test-1.akadns.net/properties/property_test#d290ddf7-53da-4509-be5a-ba582614f883",
				Errors: []Error{
					{
						Type:   "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError",
						Title:  "Property Validation Error",
						Detail: "In Property \"property_test\", there are no enabled traffic targets that have any traffic allowed to go to them",
						Errors: nil,
					},
				},
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

func TestIs(t *testing.T) {
	tests := map[string]struct {
		err      Error
		target   Error
		expected bool
	}{
		"no datacenter is assigned to map target (all others)": {
			err: Error{StatusCode: 400, Type: "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError", Title: "Property Validation Error", Detail: "Invalid configuration for property \"publishprod\": no datacenter is assigned to map target (all others)",
				Errors: nil},
			target: Error{StatusCode: 400, Type: "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError", Title: "Property Validation Error", Detail: "Invalid configuration for property \"publishprod\": no datacenter is assigned to map target (all others)",
				Errors: nil},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.err.Is(&test.target), test.expected)
		})
	}
}
