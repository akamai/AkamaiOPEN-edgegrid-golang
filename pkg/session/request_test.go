package session

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegrid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	A string `json:"a"`
	B int    `json:"b"`
}

type testInvalid struct {
	Invalid func()
}

func TestSession_Exec(t *testing.T) {
	tests := map[string]struct {
		request             *http.Request
		out                 testStruct
		in                  []interface{}
		responseBody        string
		responseStatus      int
		expectedContentType string
		expectedUserAgent   string
		expectedMethod      string
		expectedPath        string
		expected            interface{}
		withError           error
	}{
		"GET request, use default values for request": {
			request: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/test/path", nil)
				require.NoError(t, err)
				return req
			}(),
			out:            testStruct{},
			responseBody:   `{"a":"text","b":1}`,
			responseStatus: http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/test/path",
			expected: testStruct{
				A: "text",
				B: 1,
			},
		},
		"GET request, escape query": {
			request: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/test/path?param1=some param", nil)
				require.NoError(t, err)
				return req
			}(),
			out:            testStruct{},
			responseBody:   `{"a":"text","b":1}`,
			responseStatus: http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/test/path?param1=some+param",
			expected: testStruct{
				A: "text",
				B: 1,
			},
		},
		"GET request, custom content type and user agent": {
			request: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/test/path", nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "text/plain")
				req.Header.Set("User-Agent", "other user agent")
				return req
			}(),
			out:                 testStruct{},
			responseBody:        `{"a":"text","b":1}`,
			responseStatus:      http.StatusOK,
			expectedMethod:      http.MethodGet,
			expectedPath:        "/test/path",
			expectedContentType: "text/plain",
			expectedUserAgent:   "other user agent",
			expected: testStruct{
				A: "text",
				B: 1,
			},
		},
		"POST request, custom content type and user agent": {
			request: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/test/path", nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "text/plain")
				req.Header.Set("User-Agent", "other user agent")
				return req
			}(),
			in: []interface{}{&testStruct{
				A: "text",
				B: 1,
			}},
			out:                 testStruct{},
			responseBody:        `{"a":"text","b":1}`,
			responseStatus:      http.StatusCreated,
			expectedMethod:      http.MethodPost,
			expectedPath:        "/test/path",
			expectedContentType: "text/plain",
			expectedUserAgent:   "other user agent",
			expected: testStruct{
				A: "text",
				B: 1,
			},
		},
		"POST request, invalid body": {
			request: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/test/path", nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "text/plain")
				req.Header.Set("User-Agent", "other user agent")
				return req
			}(),
			in:        []interface{}{&testInvalid{func() {}}},
			out:       testStruct{},
			withError: ErrMarshaling,
		},
		"POST request, unmarshaling error": {
			request: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/test/path", nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "text/plain")
				req.Header.Set("User-Agent", "other user agent")
				return req
			}(),
			in: []interface{}{&testStruct{
				A: "text",
				B: 1,
			}},
			out:                 testStruct{},
			responseBody:        `{"a":1,"b":1}`,
			responseStatus:      http.StatusCreated,
			expectedMethod:      http.MethodPost,
			expectedPath:        "/test/path",
			expectedContentType: "text/plain",
			expectedUserAgent:   "other user agent",
			withError:           ErrUnmarshaling,
		},
		"invalid number of input parameters": {
			in:        []interface{}{testStruct{}, testStruct{}},
			withError: ErrInvalidArgument,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, test.expectedMethod, r.Method)
				if test.expectedContentType == "" {
					assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				} else {
					assert.Equal(t, test.expectedContentType, r.Header.Get("Content-Type"))
				}
				if test.expectedUserAgent == "" {
					assert.Equal(t, "test user agent", r.Header.Get("User-Agent"))
				} else {
					assert.Equal(t, test.expectedUserAgent, r.Header.Get("User-Agent"))
				}
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))

			certPool := x509.NewCertPool()
			certPool.AddCert(mockServer.Certificate())
			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: certPool,
					},
				},
			}
			serverURL, err := url.Parse(mockServer.URL)
			require.NoError(t, err)
			s, err := New(WithSigner(&edgegrid.Config{
				Host: serverURL.Host,
			}), WithClient(httpClient), WithUserAgent("test user agent"), WithHTTPTracing(true))
			require.NoError(t, err)

			_, err = s.Exec(test.request, &test.out, test.in...)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, test.out)
		})
	}
}
