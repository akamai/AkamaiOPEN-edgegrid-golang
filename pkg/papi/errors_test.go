package papi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	sess, err := session.New()
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		context.TODO(),
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
					`{"type":"a","title":"b","detail":"c"}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:       "a",
				Title:      "b",
				Detail:     "c",
				StatusCode: http.StatusInternalServerError,
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
				Title:      "Failed to unmarshal error body. PAPI API failed. Check details for more information.",
				Detail:     "test",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*papi).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestErrorIs(t *testing.T) {
	tests := map[string]struct {
		err      Error
		given    error
		expected bool
	}{
		"is ErrSBDNotEnabled": {
			err: Error{
				StatusCode: http.StatusForbidden,
				Type:       "https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/default-cert-provisioning-unavailable",
			},
			given:    ErrSBDNotEnabled,
			expected: true,
		},
		"is wrapped ErrSBDNotEnabled": {
			err: Error{
				StatusCode: http.StatusForbidden,
				Type:       "https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/default-cert-provisioning-unavailable",
			},
			given:    fmt.Errorf("oops: %w", ErrSBDNotEnabled),
			expected: true,
		},
		"is ErrDefaultCertLimitReached": {
			err: Error{
				StatusCode: http.StatusTooManyRequests,
				LimitKey:   "DEFAULT_CERTS_PER_CONTRACT",
				Remaining:  ptr.To(0),
			},
			given:    ErrDefaultCertLimitReached,
			expected: true,
		},
		"is not ErrSBDNotEnabled": {
			err: Error{
				StatusCode: http.StatusTooManyRequests,
				LimitKey:   "DEFAULT_CERTS_PER_CONTRACT",
				Remaining:  ptr.To(0),
			},
			given:    ErrSBDNotEnabled,
			expected: false,
		},
		"is not ErrDefaultCertLimitReached": {
			err: Error{
				StatusCode: http.StatusForbidden,
				Type:       "https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/default-cert-provisioning-unavailable",
			},
			given:    ErrDefaultCertLimitReached,
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := test.err.Is(test.given)
			assert.Equal(t, test.expected, res)
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
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader(`<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>`))},
			expected: &Error{
				Type:       "",
				Title:      "Failed to unmarshal error body. PAPI API failed. Check details for more information.",
				Detail:     "<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>",
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		"API failure with plain text response": {
			input: &http.Response{
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader("Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z"))},
			expected: &Error{
				Type:       "",
				StatusCode: http.StatusServiceUnavailable,
				Title:      "Failed to unmarshal error body. PAPI API failed. Check details for more information.",
				Detail:     "Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z",
			},
		},
		"API failure with XML response": {
			input: &http.Response{
				Request:    req,
				Status:     "OK",
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader(`<Root><Item id="1" name="Example" /></Root>`))},
			expected: &Error{
				Type:       "",
				StatusCode: http.StatusServiceUnavailable,
				Title:      "Failed to unmarshal error body. PAPI API failed. Check details for more information.",
				Detail:     "<Root><Item id=\"1\" name=\"Example\" /></Root>",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			p := papi{
				Session: sess,
			}
			assert.Equal(t, test.expected, p.Error(test.input))
		})
	}
}
