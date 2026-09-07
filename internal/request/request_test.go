package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestBuilder_GET(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		expectedPath string
		expectErr    error
	}{
		"GET request with basic path": {
			path:         "/test/path",
			expectedPath: "/test/path",
		},
		"GET request with string formatted path": {
			path:         "/api/objects/%s/status",
			pathArgs:     []any{"123"},
			expectedPath: "/api/objects/123/status",
		},
		"GET request with int formatted path": {
			path:         "/api/objects/%d/status",
			pathArgs:     []any{123},
			expectedPath: "/api/objects/123/status",
		},
		"GET request with multiple arguments of different types - true": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, true},
			expectedPath: "/api/objects/abc/42/active/true",
		},
		"GET request with multiple arguments of different types - false": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, false},
			expectedPath: "/api/objects/abc/42/active/false",
		},
		"expect error - invalid URL": {
			path:      "://t\tinvalid-url",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - empty URL": {
			path:      "/",
			expectErr: fmt.Errorf("path cannot be empty"),
		},
		"expect error - GET request with basic path without opening path separator": {
			path:      "test/path",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - query params provided in path": {
			path:      "/test/path?param1=value1&param2=value2",
			expectErr: fmt.Errorf("query parameters should not be provided in the path. Use one of the AddQueryParam* or AddQueryParams* methods instead"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := NewGet(ctx, tt.path, tt.pathArgs...)

			httpReq, err := req.Build()
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPath, httpReq.URL.String())
				assert.Equal(t, http.MethodGet, httpReq.Method)
				assert.Equal(t, ctx, httpReq.Context())
			}
		})
	}
}

func TestRequestBuilder_POST(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		expectedPath string
		expectErr    error
	}{
		"POST request with basic path": {
			path:         "/test/path",
			expectedPath: "/test/path",
		},
		"POST request with string formatted path": {
			path:         "/api/objects/%s/status",
			pathArgs:     []any{"123"},
			expectedPath: "/api/objects/123/status",
		},
		"POST request with int formatted path": {
			path:         "/api/objects/%d/status",
			pathArgs:     []any{123},
			expectedPath: "/api/objects/123/status",
		},
		"POST request with multiple arguments of different types - true": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, true},
			expectedPath: "/api/objects/abc/42/active/true",
		},
		"POST request with multiple arguments of different types - false": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, false},
			expectedPath: "/api/objects/abc/42/active/false",
		},
		"expect error - invalid URL": {
			path:      "://t\tinvalid-url",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - empty URL": {
			path:      "/",
			expectErr: fmt.Errorf("path cannot be empty"),
		},
		"expect error - POST request with basic path without opening path separator": {
			path:      "test/path",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - query params provided in path": {
			path:      "/test/path?param1=value1&param2=value2",
			expectErr: fmt.Errorf("query parameters should not be provided in the path. Use one of the AddQueryParam* or AddQueryParams* methods instead"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := NewPost(ctx, tt.path, tt.pathArgs...)

			httpReq, err := req.Build()
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPath, httpReq.URL.String())
				assert.Equal(t, http.MethodPost, httpReq.Method)
				assert.Equal(t, ctx, httpReq.Context())
			}
		})
	}
}

func TestRequestBuilder_PUT(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		expectedPath string
		expectErr    error
	}{
		"PUT request with basic path": {
			path:         "/test/path",
			expectedPath: "/test/path",
		},
		"PUT request with string formatted path": {
			path:         "/api/objects/%s/status",
			pathArgs:     []any{"123"},
			expectedPath: "/api/objects/123/status",
		},
		"PUT request with int formatted path": {
			path:         "/api/objects/%d/status",
			pathArgs:     []any{123},
			expectedPath: "/api/objects/123/status",
		},
		"PUT request with multiple arguments of different types - true": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, true},
			expectedPath: "/api/objects/abc/42/active/true",
		},
		"PUT request with multiple arguments of different types - false": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, false},
			expectedPath: "/api/objects/abc/42/active/false",
		},
		"expect error - invalid URL": {
			path:      "://t\tinvalid-url",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - empty URL": {
			path:      "/",
			expectErr: fmt.Errorf("path cannot be empty"),
		},
		"expect error - PUT request with basic path without opening path separator": {
			path:      "test/path",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - query params provided in path": {
			path:      "/test/path?param1=value1&param2=value2",
			expectErr: fmt.Errorf("query parameters should not be provided in the path. Use one of the AddQueryParam* or AddQueryParams* methods instead"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := NewPut(ctx, tt.path, tt.pathArgs...)

			httpReq, err := req.Build()
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPath, httpReq.URL.String())
				assert.Equal(t, http.MethodPut, httpReq.Method)
				assert.Equal(t, ctx, httpReq.Context())
			}
		})
	}
}

func TestRequestBuilder_PATCH(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		expectedPath string
		expectErr    error
	}{
		"PATCH request with basic path": {
			path:         "/test/path",
			expectedPath: "/test/path",
		},
		"PATCH request with string formatted path": {
			path:         "/api/objects/%s/status",
			pathArgs:     []any{"123"},
			expectedPath: "/api/objects/123/status",
		},
		"PATCH request with int formatted path": {
			path:         "/api/objects/%d/status",
			pathArgs:     []any{123},
			expectedPath: "/api/objects/123/status",
		},
		"PATCH request with multiple arguments of different types - true": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, true},
			expectedPath: "/api/objects/abc/42/active/true",
		},
		"PATCH request with multiple arguments of different types - false": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, false},
			expectedPath: "/api/objects/abc/42/active/false",
		},
		"expect error - invalid URL": {
			path:      "://t\tinvalid-url",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - empty URL": {
			path:      "/",
			expectErr: fmt.Errorf("path cannot be empty"),
		},
		"expect error - PATCH request with basic path without opening path separator": {
			path:      "test/path",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - query params provided in path": {
			path:      "/test/path?param1=value1&param2=value2",
			expectErr: fmt.Errorf("query parameters should not be provided in the path. Use one of the AddQueryParam* or AddQueryParams* methods instead"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := NewPatch(ctx, tt.path, tt.pathArgs...)

			httpReq, err := req.Build()
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPath, httpReq.URL.String())
				assert.Equal(t, http.MethodPatch, httpReq.Method)
				assert.Equal(t, ctx, httpReq.Context())
			}
		})
	}
}

func TestRequestBuilder_DELETE(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		expectedPath string
		expectErr    error
	}{
		"DELETE request with basic path": {
			path:         "/test/path",
			expectedPath: "/test/path",
		},
		"DELETE request with string formatted path": {
			path:         "/api/objects/%s/status",
			pathArgs:     []any{"123"},
			expectedPath: "/api/objects/123/status",
		},
		"DELETE request with int formatted path": {
			path:         "/api/objects/%d/status",
			pathArgs:     []any{123},
			expectedPath: "/api/objects/123/status",
		},
		"DELETE request with multiple arguments of different types - true": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, true},
			expectedPath: "/api/objects/abc/42/active/true",
		},
		"DELETE request with multiple arguments of different types - false": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, false},
			expectedPath: "/api/objects/abc/42/active/false",
		},
		"expect error - invalid URL": {
			path:      "://t\tinvalid-url",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - empty URL": {
			path:      "/",
			expectErr: fmt.Errorf("path cannot be empty"),
		},
		"expect error - DELETE request with basic path without opening path separator": {
			path:      "test/path",
			expectErr: fmt.Errorf("path must start with '/'"),
		},
		"expect error - query params provided in path": {
			path:      "/test/path?param1=value1&param2=value2",
			expectErr: fmt.Errorf("query parameters should not be provided in the path. Use one of the AddQueryParam* or AddQueryParams* methods instead"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := NewDelete(ctx, tt.path, tt.pathArgs...)

			httpReq, err := req.Build()
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPath, httpReq.URL.String())
				assert.Equal(t, http.MethodDelete, httpReq.Method)
				assert.Equal(t, ctx, httpReq.Context())
			}
		})
	}
}

func TestRequestBuilder_AddQueryParam(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		queryParams  map[string][]string
		expectedPath string
		expectErr    error
	}{
		"GET request with query parameters": {
			path:         "/test/path",
			queryParams:  map[string][]string{"param1": {"value1"}, "param2": {"value2"}},
			expectedPath: "/test/path?param1=value1&param2=value2",
		},
		"GET request with multi-value query parameters": {
			path:         "/test/path",
			queryParams:  map[string][]string{"param1": {"value1", "value2"}, "param2": {"value22", "value33", "value44"}},
			expectedPath: "/test/path?param1=value1&param1=value2&param2=value22&param2=value33&param2=value44",
		},
		"GET request with multiple arguments of different types and query parameters": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, true},
			queryParams:  map[string][]string{"param1": {"value1"}, "param2": {"value2"}},
			expectedPath: "/api/objects/abc/42/active/true?param1=value1&param2=value2",
		},
		"expect error - GET request with empty query parameter key": {
			path:        "/test/path",
			queryParams: map[string][]string{"": {"value1"}},
			expectErr:   fmt.Errorf("invalid query parameters: query parameter key cannot be empty"),
		},
		"expect error - GET request with empty query parameter value": {
			path:        "/test/path",
			queryParams: map[string][]string{"param1": {""}},
			expectErr:   fmt.Errorf("invalid query parameters: 'param1' query parameter value cannot contain empty value"),
		},
		"expect error - GET request with empty query parameter value among non-empty values": {
			path:        "/test/path",
			queryParams: map[string][]string{"param1": {"value1", ""}},
			expectErr:   fmt.Errorf("invalid query parameters: 'param1' query parameter value cannot contain empty value"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			req := NewGet(ctx, tt.path, tt.pathArgs...)

			for key, value := range tt.queryParams {
				for _, v := range value {
					req = req.AddQueryParam(key, v)
				}
			}

			if tt.queryParams != nil {
				for key, value := range tt.queryParams {
					assert.Equal(t, value, req.query[key])
				}
			}

			httpReq, err := req.Build()
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPath, httpReq.URL.String())
				assert.Equal(t, http.MethodGet, httpReq.Method)
				assert.Equal(t, ctx, httpReq.Context())
			}
		})
	}
}

func TestRequestBuilder_AddQueryParamIf(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		queryParams  map[string][]string
		expectedPath string
	}{
		"PUT request with query parameters - condition met": {
			path:         "/test/path",
			queryParams:  map[string][]string{"param1": {"value1"}, "param2": {"value2"}},
			expectedPath: "/test/path?param1=value1&param2=value2",
		},
		"PUT request with query parameters - condition not met - empty string": {
			path:         "/test/path",
			queryParams:  map[string][]string{"param1": {""}, "param2": {"value2"}},
			expectedPath: "/test/path?param2=value2",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			req := NewPut(ctx, tt.path, tt.pathArgs...)

			for key, value := range tt.queryParams {
				for _, v := range value {
					req = req.AddQueryParamIf(key, v, v != "")
				}
			}

			if tt.queryParams != nil {
				for key, value := range tt.queryParams {
					if !slices.Contains(value, "") {
						assert.Equal(t, value, req.query[key])
					}
				}
			}

			httpReq, err := req.Build()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedPath, httpReq.URL.String())
			assert.Equal(t, http.MethodPut, httpReq.Method)
			assert.Equal(t, ctx, httpReq.Context())
		})
	}
}

func TestRequestBuilder_AddQueryParamFunc(t *testing.T) {
	t.Parallel()

	var int64NilPtr *int64
	var intNilPtr *int
	var stringNilPtr *string
	var boolNilPtr *bool

	tests := map[string]struct {
		path         string
		pathArgs     []any
		queryParams  map[string]any
		expectedPath string
	}{
		"DELETE request with not nil int64 query": {
			path:         "/test/path",
			queryParams:  map[string]any{"param1": ptr.To(int64(1))},
			expectedPath: "/test/path?param1=1",
		},
		"DELETE request with not nil int query": {
			path:         "/test/path",
			queryParams:  map[string]any{"param1": ptr.To(1)},
			expectedPath: "/test/path?param1=1",
		},
		"DELETE request with not nil string query": {
			path:         "/test/path",
			queryParams:  map[string]any{"param1": ptr.To("value1")},
			expectedPath: "/test/path?param1=value1",
		},
		"DELETE request with not nil bool query": {
			path:         "/test/path",
			queryParams:  map[string]any{"param1": ptr.To(true)},
			expectedPath: "/test/path?param1=true",
		},
		"DELETE request with nil int64 query": {
			path:         "/test/path",
			queryParams:  map[string]any{"param1": int64NilPtr},
			expectedPath: "/test/path",
		},
		"DELETE request with nil int query": {
			path:         "/test/path",
			queryParams:  map[string]any{"param1": intNilPtr},
			expectedPath: "/test/path",
		},
		"DELETE request with nil string query": {
			path:         "/test/path",
			queryParams:  map[string]any{"param1": stringNilPtr},
			expectedPath: "/test/path",
		},
		"DELETE request with nil bool query": {
			path:         "/test/path",
			queryParams:  map[string]any{"param1": boolNilPtr},
			expectedPath: "/test/path",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			req := NewDelete(ctx, tt.path, tt.pathArgs...)

			for key, value := range tt.queryParams {
				rv := reflect.ValueOf(value)
				req = req.AddQueryParamFunc(key, func() string {
					switch v := value.(type) {
					case *int64:
						if v != nil {
							return fmt.Sprintf("%d", *v)
						}
					case *int:
						if v != nil {
							return fmt.Sprintf("%d", *v)
						}
					case *string:
						if v != nil {
							return *v
						}
					case *bool:
						if v != nil {
							return fmt.Sprintf("%t", *v)
						}
					}
					return ""
				}, !rv.IsNil())
			}

			httpReq, err := req.Build()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedPath, httpReq.URL.String())
			assert.Equal(t, http.MethodDelete, httpReq.Method)
			assert.Equal(t, ctx, httpReq.Context())
		})
	}
}

func TestRequestBuilder_AddQueryParams(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		queryParams  map[string][]string
		expectedPath string
		expectErr    error
	}{
		"POST request with multiple values for a single query parameter": {
			path:         "/test/path",
			queryParams:  map[string][]string{"param1": {"value1", "value2"}},
			expectedPath: "/test/path?param1=value1&param1=value2",
		},
		"POST request with multiple values for multiple query parameters": {
			path:         "/test/path",
			queryParams:  map[string][]string{"param1": {"value1", "value2"}, "param2": {"value22", "value33", "value44"}},
			expectedPath: "/test/path?param1=value1&param1=value2&param2=value22&param2=value33&param2=value44",
		},
		"POST request with multiple arguments of different types and query parameters": {
			path:         "/api/objects/%s/%d/active/%t",
			pathArgs:     []any{"abc", 42, true},
			queryParams:  map[string][]string{"param1": {"value1"}, "param2": {"value2"}},
			expectedPath: "/api/objects/abc/42/active/true?param1=value1&param2=value2",
		},
		"expect error - POST request with empty query parameter key": {
			path:        "/test/path",
			queryParams: map[string][]string{"": {"value1", "value2"}},
			expectErr:   fmt.Errorf("invalid query parameters: query parameter key cannot be empty"),
		},
		"expect error - POST request with empty query parameter value": {
			path:        "/test/path",
			queryParams: map[string][]string{"param1": {"", ""}},
			expectErr:   fmt.Errorf("invalid query parameters: 'param1' query parameter value cannot contain empty value"),
		},
		"expect error - POST request with empty query parameter value among non-empty values": {
			path:        "/test/path",
			queryParams: map[string][]string{"param1": {"value1", ""}},
			expectErr:   fmt.Errorf("invalid query parameters: 'param1' query parameter value cannot contain empty value"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := NewPost(ctx, tt.path, tt.pathArgs...)

			for key, value := range tt.queryParams {
				req = req.AddQueryParams(key, value)
			}

			if tt.queryParams != nil {
				for key, value := range tt.queryParams {
					assert.Equal(t, value, req.query[key])
				}
			}

			httpReq, err := req.Build()
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPath, httpReq.URL.String())
				assert.Equal(t, http.MethodPost, httpReq.Method)
				assert.Equal(t, ctx, httpReq.Context())
			}
		})
	}
}

func TestRequestBuilder_AddQueryParamsFunc(t *testing.T) {
	t.Parallel()

	type customString string

	t.Run("convert custom type to string", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		req := NewPost(ctx, "/test/path")

		customParams := []customString{"value1", "value2"}

		req = req.AddQueryParamsFunc("param1", func() []string {
			parsedValues := make([]string, 0, len(customParams))
			for _, v := range customParams {
				parsedValues = append(parsedValues, string(v))
			}
			return parsedValues
		}, len(customParams) > 0)

		assert.Equal(t, []string{"value1", "value2"}, req.query["param1"])
		httpReq, err := req.Build()
		require.NoError(t, err)
		assert.Equal(t, "/test/path?param1=value1&param1=value2", httpReq.URL.String())
		assert.Equal(t, http.MethodPost, httpReq.Method)
		assert.Equal(t, ctx, httpReq.Context())
	})
}

func TestRequestBuilder_AddHeader(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		headers   map[string][]string
		expectErr error
	}{
		"PUT request with no headers": {},
		"PUT request with one header": {
			headers: map[string][]string{"testHeader": {"testValue"}},
		},
		"PUT request with two headers": {
			headers: map[string][]string{"testHeader": {"testValue"}, "testHeader2": {"testValue2"}},
		},
		"PUT request with two headers, one contains multiple values": {
			headers: map[string][]string{"testHeader": {"testValue"}, "testHeader2": {"testValue2", "testValue3"}},
		},
		"expect error - PUT request with one header with empty key": {
			headers:   map[string][]string{"": {"testValue"}},
			expectErr: fmt.Errorf("invalid headers: header key cannot be empty"),
		},
		"expect error - PUT request with one header with empty value": {
			headers:   map[string][]string{"testHeader": {""}},
			expectErr: fmt.Errorf("invalid headers: 'Testheader' header parameter value cannot contain empty value"),
		},
		"expect error - PUT request with one header and multiple values, one empty value": {
			headers:   map[string][]string{"testHeader": {"", "testValue2"}},
			expectErr: fmt.Errorf("invalid headers: 'Testheader' header parameter value cannot contain empty value"),
		},
		"expect error - PUT request with two headers, one with empty value": {
			headers:   map[string][]string{"testHeader": {""}, "testHeader2": {"testValue2"}},
			expectErr: fmt.Errorf("invalid headers: 'Testheader' header parameter value cannot contain empty value"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			req := NewPut(ctx, "/test/path")

			for key, values := range tt.headers {
				for _, value := range values {
					req = req.AddHeader(key, value)
				}
			}

			httpReq, err := req.Build()
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
				if len(tt.headers) > 0 {
					for key, values := range tt.headers {
						for _, value := range values {
							actualHeaders := httpReq.Header.Values(key)
							assert.Contains(t, actualHeaders, value)
						}
					}
				} else {
					assert.Empty(t, httpReq.Header)
				}
				assert.Equal(t, "/test/path", httpReq.URL.String())
				assert.Equal(t, http.MethodPut, httpReq.Method)
				assert.Equal(t, ctx, httpReq.Context())
			}
		})
	}
}

func TestRequestBuilder_UseCommaSeparatedQuery(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path         string
		pathArgs     []any
		queryParams  map[string][]string
		expectedPath string
	}{
		"PATCH request with multiple values for a single query parameter": {
			path:         "/test/path",
			queryParams:  map[string][]string{"param1": {"value1", "value2"}},
			expectedPath: "/test/path?param1=value1%2Cvalue2",
		},
		"PATCH request with multiple values for multiple query parameters": {
			path:         "/test/path",
			queryParams:  map[string][]string{"param1": {"value1", "value2"}, "param2": {"value22", "value33", "value44"}},
			expectedPath: "/test/path?param1=value1%2Cvalue2&param2=value22%2Cvalue33%2Cvalue44",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := NewPatch(ctx, tt.path, tt.pathArgs...)

			for key, value := range tt.queryParams {
				req = req.AddQueryParams(key, value)
			}

			if tt.queryParams != nil {
				for key, value := range tt.queryParams {
					assert.Equal(t, value, req.query[key])
				}
			}

			req = req.UseCommaSeparatedQuery()

			httpReq, err := req.Build()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedPath, httpReq.URL.String())
			assert.Equal(t, http.MethodPatch, httpReq.Method)
			assert.Equal(t, ctx, httpReq.Context())
		})
	}
}

func TestRequestBuilder_WithBody(t *testing.T) {
	t.Parallel()

	type bodyStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := map[string]struct {
		body          any
		expectedBody  string
		expectedNil   bool
		expectedError error
	}{
		"POST with struct body": {
			body:         bodyStruct{Name: "test", Value: 42},
			expectedBody: `{"name":"test","value":42}`,
		},
		"POST with map body": {
			body:         map[string]string{"key": "value"},
			expectedBody: `{"key":"value"}`,
		},
		"POST with slice body": {
			body:         []int{1, 2, 3},
			expectedBody: `[1,2,3]`,
		},
		"POST without body - nil body": {
			body:        nil,
			expectedNil: true,
		},
		"POST with body that cannot be marshaled": {
			body:          make(chan int),
			expectedError: fmt.Errorf("failed to marshal request body"),
		},
		"POST with []byte body - used as-is, no JSON marshaling": {
			body:         []byte(`{"raw":"bytes"}`),
			expectedBody: `{"raw":"bytes"}`,
		},
		"POST with io.Reader body - used as-is, no JSON marshaling": {
			body:         strings.NewReader(`{"raw":"reader"}`),
			expectedBody: `{"raw":"reader"}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := NewPost(ctx, "/test/path")
			if tt.body != nil {
				req = req.WithBody(tt.body)
			}

			httpReq, err := req.Build()

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				return
			}

			require.NoError(t, err)

			if tt.expectedNil {
				assert.Nil(t, httpReq.Body)
				assert.Equal(t, int64(0), httpReq.ContentLength)
				return
			}

			require.NotNil(t, httpReq.Body)
			data, err := io.ReadAll(httpReq.Body)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expectedBody, string(data))
			assert.Equal(t, int64(len(data)), httpReq.ContentLength)
		})
	}
}
