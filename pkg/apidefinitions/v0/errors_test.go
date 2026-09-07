package v0

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
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
				Title:  "Failed to unmarshal error body",
				Detail: "invalid character 'e' in literal true (expecting 'r')",
				Status: http.StatusInternalServerError,
			},
		},
		"Invalid request, status code 400 with rejectedValue as stringified array": {
			response: &http.Response{
				Status:     "Invalid input error",
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bodyFromFile("testdata/400_bad_create_api_request_1.json")),
				Request:    req,
			},
			expected: &Error{
				Type:     "/api-definitions/error-types/invalid-input-error",
				Title:    "Invalid input error",
				Detail:   "The request you submitted is invalid. Modify the request and try again.",
				Instance: "f1d30806-544e-44db-b152-c95398d2bd43",
				Status:   http.StatusBadRequest,
				Errors: []Error{
					{
						Type:          "/api-definitions/error-types/endpoint-invalid-host",
						Title:         "Invalid host",
						Detail:        "The system couldn't recognize the hostnames: '[dummy-apr-msg.konaqa.com]'. Ensure the hostnames exist in the selected access control group.",
						Severity:      ptrToString("ERROR"),
						RejectedValue: strToPtrRejectedVal("[dummy-apr-msg.konaqa.com, dummy-bmp-msg.konaqa.com]"),
					},
				},
				Severity: ptrToString("ERROR"),
			},
		},
		"Invalid request, status code 400 with rejectedValue as plain string": {
			response: &http.Response{
				Status:     "Invalid input error",
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bodyFromFile("testdata/400_bad_create_api_request_2.json")),
				Request:    req,
			},
			expected: &Error{
				Type:     "/api-definitions/error-types/invalid-input-error",
				Title:    "Invalid input error",
				Detail:   "The request you submitted is invalid. Modify the request and try again.",
				Instance: "f1d30806-544e-44db-b152-c95398d2bd43",
				Status:   http.StatusBadRequest,
				Errors: []Error{
					{
						Type:          "/api-definitions/error-types/endpoint-invalid-host",
						Title:         "Invalid host",
						Detail:        "The system couldn't recognize the hostnames: dummy-apr-msg.konaqa.com. Ensure the hostnames exist in the selected access control group.",
						Severity:      ptrToString("ERROR"),
						RejectedValue: strToPtrRejectedVal("dummy-apr-msg.konaqa.com"),
					},
				},
				Severity: ptrToString("ERROR"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*apidefinitions).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func bodyFromFile(filePath string) *strings.Reader {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to read file: %s", err))
	}
	return strings.NewReader(string(content))
}

func ptrToString(s string) *string {
	return &s
}

func strToPtrRejectedVal(s string) *interface{} {
	var i interface{} = s
	return &i
}
