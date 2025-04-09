package mtlskeystore

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
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
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
	"type": "bad-request",
	"title": "Bad Request",
	"instance": "7a51100f-77ac-48dc-bd6b-ca6fcf7e820c",
	"status": 400,
	"detail": "Invalid value for field: groupId",
	"problemId": "7a51100f-77ac-48dc-bd6b-ca6fcf7e820c"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Title:     "Bad Request",
				Type:      "bad-request",
				Detail:    "Invalid value for field: groupId",
				Status:    http.StatusBadRequest,
				ProblemID: "7a51100f-77ac-48dc-bd6b-ca6fcf7e820c",
				Instance:  "7a51100f-77ac-48dc-bd6b-ca6fcf7e820c",
			},
		},
		"Bad request 400 - nested error": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(
					`{
	"type": "bad-request",
	"title": "Bad Request",
	"instance": "6225f560-c2d0-4291-974e-cd8037547529",
	"status": 400,
	"detail": "Bad Request",
	"errors": [
		{
			"type": "error-types/invalid",
			"title": "Invalid Input",
			"detail": "Certificate with same name already exists.",
			"problemId": "d5621f91-e6dc-4b3b-ab37-f085826abc44",
			"field": "certificateName"
		}
	],
	"problemId": "514c5964-dd71-4f90-98fd-7972db559273"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Title:     "Bad Request",
				Type:      "bad-request",
				Detail:    "Bad Request",
				Status:    http.StatusBadRequest,
				ProblemID: "514c5964-dd71-4f90-98fd-7972db559273",
				Instance:  "6225f560-c2d0-4291-974e-cd8037547529",
				Errors: []Error{{
					Title:     "Invalid Input",
					Type:      "error-types/invalid",
					Detail:    "Certificate with same name already exists.",
					ProblemID: "d5621f91-e6dc-4b3b-ab37-f085826abc44",
					Field:     "certificateName",
				}},
			},
		},
		"Resource not found 404": {
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(strings.NewReader(
					`{
	"type": "resource-not-found",
	"title": "Resource Not Found",
	"instance": "a2aa0865-219e-4345-9b4d-44e035f7c246",
	"status": 404,
	"detail": "The requested resource could not be found on the server.",
	"problemId": "8346a5d6-a339-4a5b-9dcf-33d482989c78",
	"field": "certificateId"
}`),
				),
				Request: req,
			},
			expected: &Error{
				Title:     "Resource Not Found",
				Type:      "resource-not-found",
				Detail:    "The requested resource could not be found on the server.",
				Status:    http.StatusNotFound,
				ProblemID: "8346a5d6-a339-4a5b-9dcf-33d482989c78",
				Instance:  "a2aa0865-219e-4345-9b4d-44e035f7c246",
				Field:     "certificateId",
			},
		},
		"Invalid response body, assign status code": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:  "Failed to unmarshal error body. mTLS Keystore API failed. Check details for more information.",
				Detail: "test",
				Status: http.StatusInternalServerError,
			},
		},
		"Empty response body, assign status code": {
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			},
			expected: &Error{
				Title:  "Failed to unmarshal error body. mTLS Keystore API failed. Check details for more information.",
				Detail: "",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*mtlskeystore).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestErrorIs(t *testing.T) {
	tests := map[string]struct {
		err      *Error
		target   error
		expected bool
	}{
		"nil target": {
			err:      &Error{Status: 404},
			target:   nil,
			expected: false,
		},
		"target not of type *Error": {
			err:      &Error{Status: 404},
			target:   errors.New("some error"),
			expected: false,
		},
		"both errors are the same instance": {
			err:      &Error{Status: 404},
			target:   &Error{Status: 404},
			expected: true,
		},
		"same status, different details": {
			err:      &Error{Status: 404, Detail: "Not Found"},
			target:   &Error{Status: 404, Detail: "Different Detail"},
			expected: false,
		},
		"completely different errors": {
			err:      &Error{Status: 404, Detail: "Not Found"},
			target:   &Error{Status: 500, Detail: "Internal Server Error"},
			expected: false,
		},
		"same status and detail": {
			err:      &Error{Status: 404, Detail: "Not Found"},
			target:   &Error{Status: 404, Detail: "Not Found"},
			expected: true,
		},
		"nested error comparison": {
			err: &Error{
				Status: 404,
				Errors: []Error{
					{Status: 400, Detail: "Nested Error"},
				},
			},
			target: &Error{
				Status: 404,
				Errors: []Error{
					{Status: 400, Detail: "Nested Error"},
				},
			},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.err.Is(test.target), test.expected)
		})
	}
}
