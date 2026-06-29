// Package request provides utilities for constructing HTTP requests with context, method, path, and query parameters.
// It supports building requests with both standard and comma-separated query parameter formats.
package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// Builder builds an HTTP request with a context, method, path, headers and optional query parameters.
type Builder struct {
	ctx                    context.Context
	method                 string
	path                   string
	query                  url.Values
	headers                http.Header
	useCommaSeparatedQuery bool
	body                   any
}

func newBuilder(ctx context.Context, method, path string, pathArgs ...any) *Builder {
	return &Builder{
		ctx:     ctx,
		method:  method,
		path:    fmt.Sprintf(path, pathArgs...),
		query:   url.Values{},
		headers: http.Header{},
	}
}

// NewGet creates a new GET request with the specified context and path.
// Path can include format specifiers for dynamic parts, followed by
// corresponding arguments to fill those specifiers.
// Example:
// NewGet(ctx, "/items/%s/details/%d", itemID, detailID).
func NewGet(ctx context.Context, path string, pathArgs ...any) *Builder {
	return newBuilder(ctx, http.MethodGet, path, pathArgs...)
}

// NewPost creates a new POST request with the specified context and path.
// Path can include format specifiers for dynamic segments, followed by
// corresponding arguments to fill those specifiers.
// Example:
// NewPost(ctx, "/items/%s/details/%d", itemID, detailID).
func NewPost(ctx context.Context, path string, pathArgs ...any) *Builder {
	return newBuilder(ctx, http.MethodPost, path, pathArgs...)
}

// NewPut creates a new PUT request with the specified context and path.
// Path can include format specifiers for dynamic segments, followed by
// corresponding arguments to fill those specifiers.
// Example:
// NewPut(ctx, "/items/%s/details/%d", itemID, detailID).
func NewPut(ctx context.Context, path string, pathArgs ...any) *Builder {
	return newBuilder(ctx, http.MethodPut, path, pathArgs...)
}

// NewPatch creates a new PATCH request with the specified context and path.
// Path can include format specifiers for dynamic segments, followed by
// corresponding arguments to fill those specifiers.
// Example:
// NewPatch(ctx, "/items/%s/details/%d", itemID, detailID).
func NewPatch(ctx context.Context, path string, pathArgs ...any) *Builder {
	return newBuilder(ctx, http.MethodPatch, path, pathArgs...)
}

// NewDelete creates a new DELETE request with the specified context and path.
// Path can include format specifiers for dynamic segments, followed by
// corresponding arguments to fill those specifiers.
// Example:
// NewDelete(ctx, "/items/%s/details/%d", itemID, detailID).
func NewDelete(ctx context.Context, path string, pathArgs ...any) *Builder {
	return newBuilder(ctx, http.MethodDelete, path, pathArgs...)
}

// AddQueryParam adds a query parameter to the request.
func (b *Builder) AddQueryParam(key, value string) *Builder {
	b.query.Add(key, value)
	return b
}

// AddQueryParamIf adds a query parameter to the request if the condition is met.
func (b *Builder) AddQueryParamIf(key, value string, condition bool) *Builder {
	if condition {
		b.query.Add(key, value)
	}
	return b
}

// AddQueryParamFunc adds a query parameter to the request using a function to generate the value,
// if the condition is met.
func (b *Builder) AddQueryParamFunc(key string, valueFunc func() string, condition bool) *Builder {
	if condition {
		v := valueFunc()
		b.query.Add(key, v)
	}
	return b
}

// AddQueryParams adds a query parameter with multiple values to the request.
func (b *Builder) AddQueryParams(key string, value []string) *Builder {
	for _, v := range value {
		b.query.Add(key, v)
	}
	return b
}

// AddQueryParamsFunc adds a query parameter to the request using a function to generate the values,
// if the condition is met.
func (b *Builder) AddQueryParamsFunc(key string, valuesFunc func() []string, condition bool) *Builder {
	if condition {
		vv := valuesFunc()
		for _, v := range vv {
			b.query.Add(key, v)
		}
	}
	return b
}

// AddHeader adds a header to the request.
func (b *Builder) AddHeader(key, value string) *Builder {
	b.headers.Add(key, value)
	return b
}

// WithBody sets the request body. The body will be JSON-marshaled when Build is called,
// unless it is already a []byte or an io.Reader, in which case it is used as-is.
func (b *Builder) WithBody(body any) *Builder {
	b.body = body
	return b
}

// UseCommaSeparatedQuery enables the use of comma-separated query parameters.
// When enabled, multiple values for the same query parameter key will be combined
// into a single key-value pair with values separated by commas.
// For example, if you add "param" with values "value1" and "value2":
//   - By default, the query string will be: "param=value1&param=value2"
//   - With UseCommaSeparatedQuery enabled, the query string will be: "param=value1,value2"
func (b *Builder) UseCommaSeparatedQuery() *Builder {
	b.useCommaSeparatedQuery = true
	return b
}

// Build constructs the HTTP request based on the provided context, method, path, and query parameters.
// By default, standard multi-value query parameters are used.
func (b *Builder) Build() (*http.Request, error) {
	if !strings.HasPrefix(b.path, "/") {
		return nil, fmt.Errorf("path must start with '/'")
	}
	uri, err := url.Parse(b.path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path: %w", err)
	}
	if uri.RawQuery != "" {
		return nil, fmt.Errorf("query parameters should not be provided in the path. Use one of the AddQueryParam* or AddQueryParams* methods instead")
	}
	if uri.Path == "/" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	if err = validateQuery(b.query); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	if err := validateHeaders(b.headers); err != nil {
		return nil, fmt.Errorf("invalid headers: %w", err)
	}

	if len(b.query) != 0 && b.useCommaSeparatedQuery {
		commaSeparatedQuery := url.Values{}
		for key, value := range b.query {
			commaSeparatedQuery.Set(key, strings.Join(value, ","))
		}
		uri.RawQuery = commaSeparatedQuery.Encode()
	} else {
		uri.RawQuery = b.query.Encode()
	}

	var bodyReader io.Reader
	if b.body != nil {
		switch v := b.body.(type) {
		case io.Reader:
			bodyReader = v
		case []byte:
			bodyReader = bytes.NewReader(v)
		default:
			data, err := json.Marshal(b.body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(data)
		}
	}

	req, err := http.NewRequestWithContext(b.ctx, b.method, uri.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header = b.headers

	return req, nil
}

func validateQuery(query url.Values) error {
	for key, values := range query {
		if key == "" {
			return fmt.Errorf("query parameter key cannot be empty")
		}
		if slices.Contains(values, "") {
			return fmt.Errorf("'%s' query parameter value cannot contain empty value", key)
		}
	}
	return nil
}

func validateHeaders(headers http.Header) error {
	for key, values := range headers {
		if key == "" {
			return fmt.Errorf("header key cannot be empty")
		}
		if slices.Contains(values, "") {
			return fmt.Errorf("'%s' header parameter value cannot contain empty value", key)
		}
	}
	return nil
}
