package session

import (
	"context"
	"net/http"
	"runtime"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		options  []Option
		expected *session
		err      string
	}{
		"no options provided, return default session": {
			options: []Option{WithSigner(&edgegrid.Config{})},
			expected: &session{
				client:    http.DefaultClient,
				signer:    &edgegrid.Config{},
				log:       log.Default(),
				trace:     false,
				userAgent: "Akamai-Open-Edgegrid-golang/11.0.0 golang/" + strings.TrimPrefix(runtime.Version(), "go"),
			},
		},
		"nil client provided, return error": {
			options: []Option{WithClient(nil)},
			err:     "client should not be nil",
		},
		"nil log provided, return error": {
			options: []Option{WithLog(nil)},
			err:     "logger should not be nil",
		},
		"empty user agent provided, return error": {
			options: []Option{WithUserAgent("")},
			err:     "user agent should not be empty",
		},
		"nil signer provided, return error": {
			options: []Option{WithSigner(nil)},
			err:     "signer should not be nil",
		},
		"invalid retries provided, return error": {
			options: []Option{WithRetries(RetryConfig{
				RetryMax:          -1,
				RetryWaitMin:      -1,
				RetryWaitMax:      -2,
				ExcludedEndpoints: []string{"f:o#[]o"},
			})},
			err: "retry configuration failed: maximum number of retries cannot be negative\n" +
				"minimum retry wait time cannot be negative\n" +
				"maximum retry wait time cannot be negative\n" +
				"maximum retry wait time cannot be shorter than minimum retry wait time\n" +
				"malformed exclude endpoint pattern: syntax error in pattern: f:o#[]o",
		},
		"with options provided": {
			options: []Option{
				WithSigner(&edgegrid.Config{}),
				WithClient(&http.Client{Timeout: 500}),
				WithLog(log.NOPLogger()),
				WithUserAgent("test user agent"),
				WithHTTPTracing(true)},
			expected: &session{
				client: &http.Client{
					Timeout: 500,
				},
				signer:    &edgegrid.Config{},
				log:       log.NOPLogger(),
				trace:     true,
				userAgent: "test user agent",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := New(test.options...)
			if test.err != "" {
				assert.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, res)
			}
		})
	}
}

func TestSession_Log(t *testing.T) {
	tests := map[string]struct {
		ctx           context.Context
		sessionLogger log.Interface
		expected      log.Interface
	}{
		"logger found in context, omit logger from session": {
			ctx:           ContextWithOptions(context.Background(), WithContextLog(log.NOPLogger())),
			sessionLogger: log.Default(),
			expected:      log.NOPLogger(),
		},
		"logger not found in context, pick logger from session": {
			ctx:           context.Background(),
			sessionLogger: log.NOPLogger(),
			expected:      log.NOPLogger(),
		},
		"logger not found in context or session": {
			ctx:           context.Background(),
			sessionLogger: nil,
			expected:      log.Default(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := session{log: test.sessionLogger}
			res := s.Log(test.ctx)
			assert.Equal(t, test.expected, res)
		})
	}
}
