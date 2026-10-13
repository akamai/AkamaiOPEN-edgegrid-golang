package hapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	t.Parallel()
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
				Type:   "a",
				Title:  "b",
				Detail: "c",
				Status: http.StatusInternalServerError,
			},
		},
		"valid response, status code 404": {
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(strings.NewReader(`{
					"type":"/hapi/problems/change-id-not-found",
					"title":"Invalid Change ID",
					"detail":"Change request not found for change ID 5555555",
					"changeId":"5555555",
					"status":404,
					"errors":[],
					"method": "GET",
					"requestTime": "2026-10-02T08:43:36.499655720Z",
					"instance": "/hapi/error-instances/123",
					"extensionFields":[{"key":"changeId","value":"5555555"}],
					"domainPrefix": "example.com",
					"domainSuffix": "example.net",
					"requestInstance": "http://example-cloud.com/123"
				}`)),
				Request: req,
			},
			expected: &Error{
				Type:            "/hapi/problems/change-id-not-found",
				Title:           "Invalid Change ID",
				Detail:          "Change request not found for change ID 5555555",
				ChangeID:        "5555555",
				Method:          "GET",
				RequestTime:     time.Date(2026, time.October, 2, 8, 43, 36, 499655720, time.UTC),
				Status:          http.StatusNotFound,
				Errors:          []ErrorItem{},
				ExtensionFields: []ErrorItem{{Key: "changeId", Value: "5555555"}},
				Instance:        "/hapi/error-instances/123",
				DomainPrefix:    "example.com",
				DomainSuffix:    "example.net",
				RequestInstance: "http://example-cloud.com/123",
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
				Title:  "Failed to unmarshal error body. HAPI API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
		"response body read failure": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(iotest.ErrReader(errors.New("read failed"))),
				Request:    req,
			},
			expected: &Error{Status: http.StatusBadRequest, Title: "Failed to read error body", Detail: "read failed"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := Client(sess).(*hapi).Error(test.response)
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
				Request: req,
				Status:  "OK",
				Body:    io.NopCloser(strings.NewReader(`<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>`))},
			expected: &Error{
				Type:   "",
				Title:  "Failed to unmarshal error body. HAPI API failed. Check details for more information.",
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
				Title:  "Failed to unmarshal error body. HAPI API failed. Check details for more information.",
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
				Title:  "Failed to unmarshal error body. HAPI API failed. Check details for more information.",
				Detail: "<Root><Item id=\"1\" name=\"Example\" /></Root>",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sess, _ := session.New()
			g := hapi{
				Session: sess,
			}
			assert.Equal(t, test.expected, g.Error(test.input))
		})
	}
}
