package edgegrid

import (
	"encoding/base64"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestConfig_createAuthHeader(t *testing.T) {
	tests := map[string]struct {
		config    Config
		request   *http.Request
		expected  authHeader
		withError error
	}{
		"method is GET": {
			config: Config{
				ClientToken: "12345",
				AccessToken: "54321",
				MaxBody:     MaxBodySize,
			},
			request: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://akamai.com/test/path?query=test", nil)
				require.NoError(t, err)
				return req
			}(),
			expected: authHeader{
				authType:    authType,
				clientToken: "12345",
				accessToken: "54321",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := test.config.createAuthHeader(test.request)
			assert.Equal(t, test.expected.authType, res.authType)
			assert.Equal(t, test.expected.accessToken, res.accessToken)
			assert.Equal(t, test.expected.clientToken, res.clientToken)
			assert.NotEmpty(t, res.signature)
			_, err := uuid.Parse(res.nonce)
			assert.NoError(t, err)
			_, err = base64.StdEncoding.DecodeString(res.signature)
			require.NoError(t, err)
			_, err = time.Parse("20060102T15:04:05-0700", res.timestamp)
			assert.NoError(t, err)
		})
	}
}

func TestCanonicalizeHeaders(t *testing.T) {
	tests := map[string]struct {
		requestHeaders http.Header
		headersToSign  []string
		expected       string
	}{
		"found matching request headers": {
			requestHeaders: map[string][]string{
				"A": {"val1"},
				"B": {"  VAL   2   "},
				"C": {"V A L 3"},
			},
			headersToSign: []string{"B", "C"},
			expected:      "b:val 2\tc:v a l 3",
		},
		"no matching headers found": {
			requestHeaders: map[string][]string{
				"A": {"val1"},
				"B": {"  VAL   2   "},
				"C": {"V A L 3"},
			},
			headersToSign: []string{"D", "E"},
			expected:      "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := canonicalizeHeaders(test.requestHeaders, test.headersToSign)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestCreateContentHash(t *testing.T) {
	tests := map[string]struct {
		httpMethod  string
		body        string
		resultEmpty bool
	}{
		"PUT request": {
			httpMethod:  http.MethodPut,
			body:        `{"key":"value"}`,
			resultEmpty: true,
		},
		"POST request, empty body": {
			httpMethod:  http.MethodPost,
			body:        "",
			resultEmpty: true,
		},
		"POST request, body is not empty": {
			httpMethod:  http.MethodPost,
			body:        `{"key":"value"}`,
			resultEmpty: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(test.httpMethod, "", strings.NewReader(test.body))
			require.NoError(t, err)
			res := createContentHash(req, MaxBodySize)
			if test.resultEmpty {
				assert.Empty(t, res)
				return
			}
			require.NotEmpty(t, res)
			_, err = base64.StdEncoding.DecodeString(res)
			assert.NoError(t, err)
		})
	}
}

func TestAuthHeader_String(t *testing.T) {
	tests := map[string]struct {
		given    authHeader
		expected string
	}{
		"signature is empty": {
			given: authHeader{
				authType:    "A",
				clientToken: "B",
				accessToken: "C",
				timestamp:   "D",
				nonce:       "E",
			},
			expected: "A client_token=B;access_token=C;timestamp=D;nonce=E;",
		},
		"signature is not empty": {
			given: authHeader{
				authType:    "A",
				clientToken: "B",
				accessToken: "C",
				timestamp:   "D",
				nonce:       "E",
				signature:   "F",
			},
			expected: "A client_token=B;access_token=C;timestamp=D;nonce=E;signature=F",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := test.given.String()
			assert.Equal(t, test.expected, res)
		})
	}
}
