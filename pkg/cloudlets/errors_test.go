package cloudlets

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
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
				Title:      "test",
				Detail:     "",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*cloudlets).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestAs(t *testing.T) {
	someErrorMarshalled, _ := json.Marshal("some error")
	tests := map[string]struct {
		err      Error
		target   Error
		expected bool
	}{
		"different error code": {
			err:      Error{StatusCode: 404},
			target:   Error{StatusCode: 401},
			expected: false,
		},
		"same error code": {
			err:      Error{StatusCode: 404},
			target:   Error{StatusCode: 404},
			expected: true,
		},
		"same error code and error message": {
			err:      Error{StatusCode: 404, Errors: someErrorMarshalled},
			target:   Error{StatusCode: 404, Errors: someErrorMarshalled},
			expected: true,
		},
		"same error code and different error message": {
			err:      Error{StatusCode: 404, Errors: someErrorMarshalled},
			target:   Error{StatusCode: 404},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.err.Is(&test.target), test.expected)
		})
	}
}
