package session

import (
	"github.com/apex/log"
	"github.com/tj/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewAPIError(t *testing.T) {
	tests := map[string]struct {
		response *http.Response
		expected APIError
	}{
		"valid response, status code 500": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body: ioutil.NopCloser(strings.NewReader(
					`{"type":"a","title":"b","detail":"c"}`),
				),
			},
			expected: APIError{
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
			},
			expected: APIError{
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := NewAPIError(test.response, log.Log)
			assert.Equal(t, test.expected, res)
		})
	}
}
