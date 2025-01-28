// Package session provides the base secure http client and request management for akamai apis
package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
)

type (
	// Session is the interface that is used by the pa
	// This allows the client itself to be more extensible and readily testable, ets.
	Session interface {
		// Exec will sign and execute a request returning the response
		// The response body will be unmarshaled in to out
		// Optionally the in value will be marshaled into the body
		Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error)

		// Sign will only sign a request, this is useful for circumstances
		// when the caller wishes to manage the http client
		Sign(r *http.Request) error

		// Log returns the logging interface for the session
		// If provided all debugging will output to this log interface
		Log(ctx context.Context) log.Interface

		// Client return the session http client
		Client() *http.Client
	}

	// session is the base akamai http client
	session struct {
		client       *http.Client
		signer       edgegrid.Signer
		log          log.Interface
		trace        bool
		userAgent    string
		requestLimit int
	}

	contextOptions struct {
		log    log.Interface
		header http.Header
	}

	// Option defines a client option
	Option func(*session) error

	contextKey string

	// ContextOption are options on the context
	ContextOption func(*contextOptions)
)

var (
	contextOptionKey = contextKey("sessionContext")
)

const (
	// Version is the client version
	Version = "10.0.0"
)

// New returns a new session
func New(opts ...Option) (Session, error) {
	var (
		defaultUserAgent = "Akamai-Open-Edgegrid-golang/" + Version + " golang/" + strings.TrimPrefix(runtime.Version(), "go")
	)

	s := &session{
		client:    http.DefaultClient,
		log:       log.Default(),
		userAgent: defaultUserAgent,
		trace:     false,
	}

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	if s.signer == nil {
		config, err := edgegrid.New()
		if err != nil {
			return nil, err
		}
		s.signer = config
	}

	return s, nil
}

// Must is a helper that will result in a panic if an error is returned
// ex. sess := Must(New())
func Must(sess Session, err error) Session {
	if err != nil {
		panic(err)
	}

	return sess
}

// WithClient creates a client using the specified http.Client
func WithClient(client *http.Client) Option {
	return func(s *session) error {
		if client == nil {
			return errors.New("client should not be nil")
		}
		s.client = client
		return nil
	}
}

// WithRetries configures the HTTP client to automatically retry failed GET requests
func WithRetries(conf RetryConfig) Option {
	return func(s *session) error {
		retryClient, err := configureRetryClient(conf, s.Sign, s.log)
		if err != nil {
			return fmt.Errorf("retry configuration failed: %w", err)
		}
		s.client = retryClient.StandardClient()
		return nil
	}
}

// WithLog sets the log interface for the client
func WithLog(l log.Interface) Option {
	return func(s *session) error {
		if l == nil {
			return errors.New("logger should not be nil")
		}
		s.log = l
		return nil
	}
}

// WithUserAgent sets the user agent string for the client
func WithUserAgent(u string) Option {
	return func(s *session) error {
		if u == "" {
			return errors.New("user agent should not be empty")
		}
		s.userAgent = u
		return nil
	}
}

// WithSigner sets the request signer for the session
func WithSigner(signer edgegrid.Signer) Option {
	return func(s *session) error {
		if signer == nil {
			return errors.New("signer should not be nil")
		}
		s.signer = signer
		return nil
	}
}

// WithRequestLimit sets the maximum number of API calls that the provider will make per second.
func WithRequestLimit(requestLimit int) Option {
	return func(s *session) error {
		s.requestLimit = requestLimit
		return nil
	}
}

// WithHTTPTracing sets the request and response dump for debugging
func WithHTTPTracing(trace bool) Option {
	return func(s *session) error {
		s.trace = trace
		return nil
	}
}

// Log returns the context logger, or the session log
func (s *session) Log(ctx context.Context) log.Interface {
	if o := ctx.Value(contextOptionKey); o != nil {
		if ops, ok := o.(*contextOptions); ok && ops.log != nil {
			return ops.log
		}
	}
	if s.log != nil {
		return s.log
	}

	// if context/session logs were not set, it will return default logger
	return log.Default()
}

// Client returns the http client interface
func (s *session) Client() *http.Client {
	return s.client
}

// ContextWithOptions adds request-specific options to the context
// This log debugs the request using only the provided log
func ContextWithOptions(ctx context.Context, opts ...ContextOption) context.Context {
	o := new(contextOptions)
	for _, opt := range opts {
		opt(o)
	}

	return context.WithValue(ctx, contextOptionKey, o)
}

// WithContextLog provides a context specific logger
func WithContextLog(l log.Interface) ContextOption {
	return func(o *contextOptions) {
		o.log = l
	}
}

// WithContextHeaders sets the context headers
func WithContextHeaders(h http.Header) ContextOption {
	return func(o *contextOptions) {
		o.header = h
	}
}

// CloseResponseBody closes response body
func CloseResponseBody(resp *http.Response) {
	_ = resp.Body.Close()
}
