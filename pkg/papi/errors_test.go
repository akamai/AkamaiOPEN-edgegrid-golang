package papi

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
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
				Body: ioutil.NopCloser(strings.NewReader(
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
				Body: ioutil.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:      "Failed to unmarshal error body",
				Detail:     "invalid character 'e' in literal true (expecting 'r')",
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
				Remaining:  tools.IntPtr(0),
			},
			given:    ErrDefaultCertLimitReached,
			expected: true,
		},
		"is not ErrSBDNotEnabled": {
			err: Error{
				StatusCode: http.StatusTooManyRequests,
				LimitKey:   "DEFAULT_CERTS_PER_CONTRACT",
				Remaining:  tools.IntPtr(0),
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
