package gtm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGTM_ListASMap(t *testing.T) {
	var result ASMapList

	respData, err := loadTestData("TestGTM_ListASMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []*ASMap
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/as-maps",
			expectedResponse: result.ASMapItems,
		},
		"500 internal server error": {
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching asmap",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/as-maps",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching asmap",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListASMaps(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_GetASMap(t *testing.T) {
	var result ASMap

	respData, err := loadTestData("TestGTM_GetASMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		name             string
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ASMap
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			name:       "Software-rollout",
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/as-maps/Software-rollout",
			expectedResponse: &result,
		},
		"500 internal server error": {
			name:           "Software-rollout",
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching asmap",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/as-maps/Software-rollout",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching asmap",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetASMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.name, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_CreateASMap(t *testing.T) {
	var result ASMapResponse
	var req ASMap

	respData, err := loadTestData("TestGTM_CreateASMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateASMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		asMap            *ASMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ASMapResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			asMap:      &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
			expectedResponse: &result,
		},
		"500 internal server error": {
			asMap:          &req,
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating asmap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating asmap",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateASMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.asMap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_UpdateASMap(t *testing.T) {
	var result ASMapResponse
	var req ASMap

	respData, err := loadTestData("TestGTM_CreateASMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateASMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		asMap            *ASMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			asMap:      &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
			expectedResponse: result.Status,
		},
		"500 internal server error": {
			asMap:          &req,
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating asmap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating asmap",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateASMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.asMap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_DeleteASMap(t *testing.T) {
	var result ASMapResponse
	var req ASMap

	respData, err := loadTestData("TestGTM_CreateASMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateASMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		asMap            *ASMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			asMap:      &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
			expectedResponse: result.Status,
		},
		"500 internal server error": {
			asMap:          &req,
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating asmap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating asmap",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteASMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.asMap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
