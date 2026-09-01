package accountprotection

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				Title:      "Failed to unmarshal error body. Bot Manager API failed. Check details for more information.",
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
				Title:      "Failed to unmarshal error body. Bot Manager API failed. Check details for more information.",
				Detail:     "Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z",
				StatusCode: http.StatusServiceUnavailable,
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
				Title:      "Failed to unmarshal error body. Bot Manager API failed. Check details for more information.",
				Detail:     "<Root><Item id=\"1\" name=\"Example\" /></Root>",
				StatusCode: http.StatusServiceUnavailable,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			b := accountProtection{
				Session: sess,
			}
			assert.Equal(t, test.expected, b.Error(test.input))
		})
	}
}

func TestErrorError(t *testing.T) {
	tests := map[string]struct {
		err      *Error
		expected string
	}{
		"without errors": {
			err: &Error{
				Type:       "a",
				Title:      "b",
				Detail:     "c",
				StatusCode: http.StatusBadRequest,
			},
			expected: "Title: b; Type: a; Detail: c",
		},
		"with errors": {
			err: &Error{
				Type:       "parent-type",
				Title:      "parent-title",
				Detail:     "parent-detail",
				StatusCode: http.StatusBadRequest,
				Errors: []Error{
					{
						Type:       "child1-type",
						Title:      "child1-title",
						Detail:     "child1-detail",
						StatusCode: http.StatusUnauthorized,
					},
					{
						Type:       "child2-type",
						Title:      "child2-title",
						Detail:     "child2-detail",
						StatusCode: http.StatusPaymentRequired,
					},
				},
			},
			expected: "Title: parent-title; Type: parent-type; Detail: parent-detail: [child1-detail, child2-detail ]",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.err.Error())
		})
	}
}

func TestErrorIs(t *testing.T) {
	testError := Error{
		Type:       "a",
		Title:      "b",
		Detail:     "c",
		StatusCode: http.StatusBadRequest,
	}
	tests := map[string]struct {
		err      *Error
		target   error
		expected bool
	}{
		"unwrap fail": {
			err:      &Error{},
			target:   errors.New("NOPE"),
			expected: false,
		},
		"unwrap ok": {
			err:      &testError,
			target:   fmt.Errorf("wrapped error: %w", &testError),
			expected: true,
		},
		"empty": {
			err:      &Error{},
			target:   &Error{},
			expected: true,
		},
		"same pointer": {
			err:      &testError,
			target:   &testError,
			expected: true,
		},
		"all fields equal": {
			err: &testError,
			target: &Error{
				Type:       "a",
				Title:      "b",
				Detail:     "c",
				StatusCode: http.StatusBadRequest,
			},
			expected: true,
		},
		"same child errors list": {
			err: &Error{
				Errors: []Error{{Type: "a"}},
			},
			target: &Error{
				Errors: []Error{{Type: "a"}},
			},
			expected: true,
		},
		"child errors list with different type": {
			err: &Error{
				Errors: []Error{{Type: "a"}},
			},
			target: &Error{
				Errors: []Error{{Type: "notA"}},
			},
			expected: true, // should we fix this?
		},
		"child errors list with different detail": {
			err: &Error{
				Errors: []Error{{Detail: "a"}},
			},
			target: &Error{
				Errors: []Error{{Detail: "notA"}},
			},
			expected: false,
		},
		"different type": {
			err: &Error{
				Type: "a",
			},
			target: &Error{
				Type: "notA",
			},
			expected: false,
		},
		"different title": {
			err: &Error{
				Title: "b",
			},
			target: &Error{
				Title: "notB",
			},
			expected: false,
		},
		"different detail": {
			err: &Error{
				Detail: "c",
			},
			target: &Error{
				Detail: "notC",
			},
			expected: false,
		},
		"different status code": {
			err: &Error{
				StatusCode: http.StatusBadRequest,
			},
			target: &Error{
				StatusCode: http.StatusUnauthorized,
			},
			expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, errors.Is(test.err, test.target))
		})
	}
}
