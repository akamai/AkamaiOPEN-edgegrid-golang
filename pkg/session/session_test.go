package session

import (
	"context"
	"net/http"
	"runtime"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegrid"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		client      *http.Client
		log         log.Interface
		userAgent   string
		httpTracing bool
		expected    *session
	}{
		"no options provided, return default session": {
			expected: &session{
				client:    http.DefaultClient,
				signer:    &edgegrid.Config{},
				log:       log.Log,
				trace:     false,
				userAgent: "Akamai-Open-Edgegrid-golang/2.0.0 golang/" + strings.TrimPrefix(runtime.Version(), "go"),
			},
		},
		"with options provided": {
			client: &http.Client{
				Timeout: 500,
			},
			log:         log.Log,
			userAgent:   "test user agent",
			httpTracing: true,
			expected: &session{
				client: &http.Client{
					Timeout: 500,
				},
				signer:    &edgegrid.Config{},
				log:       log.Log,
				trace:     true,
				userAgent: "test user agent",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cfg := &edgegrid.Config{}
			options := []Option{WithSigner(cfg)}
			if test.client != nil {
				options = append(options, WithClient(test.client))
			}
			if test.log != nil {
				options = append(options, WithLog(test.log))
			}
			if test.userAgent != "" {
				options = append(options, WithUserAgent(test.userAgent))
			}
			if test.httpTracing {
				options = append(options, WithHTTPTracing(test.httpTracing))
			}
			res, err := New(options...)
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestSession_Log(t *testing.T) {
	tests := map[string]struct {
		ctx           context.Context
		sessionLogger log.Interface
		expected      *log.Logger
	}{
		"logger found in context, omit logger from session": {
			ctx: ContextWithOptions(context.Background(), WithContextLog(&log.Logger{
				Handler: discard.New(),
				Level:   1,
			})),
			sessionLogger: &log.Logger{
				Handler: discard.New(),
				Level:   2,
			},
			expected: &log.Logger{
				Handler: discard.New(),
				Level:   1,
			},
		},
		"logger not found in context, pick logger from session": {
			ctx: context.Background(),
			sessionLogger: &log.Logger{
				Handler: discard.New(),
				Level:   2,
			},
			expected: &log.Logger{
				Handler: discard.New(),
				Level:   2,
			},
		},
		"logger not found in context or session": {
			ctx:           context.Background(),
			sessionLogger: nil,
			expected: &log.Logger{
				Handler: discard.New(),
			},
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
