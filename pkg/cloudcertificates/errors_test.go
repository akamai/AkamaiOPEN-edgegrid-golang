package cloudcertificates

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	t.Parallel()
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
		is       error // if set, asserts errors.Is(res, is) matches a known sentinel error
	}{
		"403 - account not allowed": {
			response: &http.Response{
				StatusCode: http.StatusForbidden,
				Body: io.NopCloser(strings.NewReader(`
					{
						"type": "/error-types/lineage-account-not-allowed",
						"title": "Account is not allowed to use Certificate Lineage.",
						"status": 403,
						"detail": "This account is not permitted to access Certificate Lineage APIs.",
						"instance": "/error-types/lineage-account-not-allowed?traceId=1234567891011"
					}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "/error-types/lineage-account-not-allowed",
				Title:    "Account is not allowed to use Certificate Lineage.",
				Status:   http.StatusForbidden,
				Detail:   "This account is not permitted to access Certificate Lineage APIs.",
				Instance: "/error-types/lineage-account-not-allowed?traceId=1234567891011",
			},
			is: ErrLineageAccountNotAllowed,
		},
		"409 - cert already uploaded with context": {
			response: &http.Response{
				StatusCode: http.StatusConflict,
				Body: io.NopCloser(strings.NewReader(`
					{
						"type": "/error-types/cert-already-uploaded",
						"title": "A signed certificate has already been uploaded for this algorithm instance.",
						"status": 409,
						"detail": "Algorithm instance {RSA} on generation {2912} already has an accepted signed certificate. Re-upload is not permitted.",
						"instance": "/error-types/cert-already-uploaded?traceId=1234567891012",
						"context": {"generationId": 2912, "keyType": "RSA"}
					}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "/error-types/cert-already-uploaded",
				Title:    "A signed certificate has already been uploaded for this algorithm instance.",
				Status:   http.StatusConflict,
				Detail:   "Algorithm instance {RSA} on generation {2912} already has an accepted signed certificate. Re-upload is not permitted.",
				Instance: "/error-types/cert-already-uploaded?traceId=1234567891012",
				Context:  map[string]any{"generationId": float64(2912), "keyType": "RSA"},
			},
			is: ErrCertAlreadyUploaded,
		},
		"415 - media type not supported, context contains array": {
			response: &http.Response{
				StatusCode: http.StatusUnsupportedMediaType,
				Body: io.NopCloser(strings.NewReader(`
					{
						"type": "/error-types/media-type-not-supported",
						"title": "Media type not supported.",
						"status": 415,
						"detail": "Media type {application/json-patch+json} is not supported. Supported media type(s) are {[application/json]}.",
						"instance": "/error-types/media-type-not-supported?traceId=1234567891013",
						"context": {
							"allowedMediaTypes": ["application/json"],
							"unsupportedMediaType": "application/json-patch+json"
						}
					}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:     "/error-types/media-type-not-supported",
				Title:    "Media type not supported.",
				Status:   http.StatusUnsupportedMediaType,
				Detail:   "Media type {application/json-patch+json} is not supported. Supported media type(s) are {[application/json]}.",
				Instance: "/error-types/media-type-not-supported?traceId=1234567891013",
				Context: map[string]any{
					"allowedMediaTypes":    []any{"application/json"},
					"unsupportedMediaType": "application/json-patch+json",
				},
			},
			is: ErrMediaTypeNotSupported,
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
				Title:  "Failed to unmarshal error body. CCM API failed. Check details for more information.",
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
				Title:  "Failed to unmarshal error body. CCM API failed. Check details for more information.",
				Detail: "",
				Status: http.StatusInternalServerError,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*cloudcertificates).Error(tc.response)
			assert.Equal(t, tc.expected, res)
			if tc.is != nil {
				assert.True(t, errors.Is(res, tc.is))
			}
		})
	}
}

func TestIs(t *testing.T) {
	t.Parallel()
	someError := Error{Type: "/some/type", Title: "some error", Status: 404, Detail: "some detail", Instance: "/some/error/instance"}

	tests := map[string]struct {
		target   Error
		expected bool
	}{
		"different error status": {
			target:   Error{Status: 401},
			expected: false,
		},
		"different error title": {
			target:   Error{Title: "other error"},
			expected: false,
		},
		"same error title": {
			target:   Error{Title: "some error"},
			expected: true,
		},
		"same error type": {
			target:   Error{Type: "/some/type"},
			expected: true,
		},
		"same error status": {
			target:   Error{Status: 404},
			expected: true,
		},
		"same error type, title and status": {
			target:   Error{Type: "/some/type", Title: "some error", Status: 404},
			expected: true,
		},
		"same error type but different title": {
			target:   Error{Type: "/some/type", Title: "other error"},
			expected: false,
		},
		"same error status and title but different type": {
			target:   Error{Type: "/other/type", Title: "some error", Status: 404},
			expected: false,
		},
		"same error status and type but different detail and instance": {
			target:   Error{Type: "/some/type", Status: 404, Detail: "other detail", Instance: "/other/error/instance"},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, someError.Is(&test.target), test.expected)
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	e := &Error{
		Type:     "/error-types/test",
		Title:    "Test Error",
		Status:   400,
		Detail:   "This is a test error",
		Instance: "/error-types/test?traceId=12345",
	}
	expected := `API error: 
{
	"type": "/error-types/test",
	"title": "Test Error",
	"status": 400,
	"detail": "This is a test error",
	"instance": "/error-types/test?traceId=12345"
}`
	assert.EqualError(t, e, expected)
}
