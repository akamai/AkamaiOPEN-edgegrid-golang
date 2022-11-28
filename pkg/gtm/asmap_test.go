package gtm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGtm_NewAsMap(t *testing.T) {
	client := Client(session.Must(session.New()))

	asmap := client.NewAsMap(context.Background(), "foo")

	assert.Equal(t, "foo", asmap.Name)
}

func TestGtm_ListAsMap(t *testing.T) {
	var result AsMapList

	respData, err := loadTestData("TestGtm_ListAsMap.resp.json")
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
		expectedResponse []*AsMap
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
			expectedResponse: result.AsMapItems,
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
			result, err := client.ListAsMaps(
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

func TestGtm_GetAsMap(t *testing.T) {
	var result AsMap

	respData, err := loadTestData("TestGtm_GetAsMap.resp.json")
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
		expectedResponse *AsMap
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
			result, err := client.GetAsMap(
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

func TestGtm_NewASAssignment(t *testing.T) {
	client := Client(session.Must(session.New()))

	asn := client.NewASAssignment(context.Background(), nil, 100, "foo")

	assert.Equal(t, 100, asn.DatacenterId)
	assert.Equal(t, "foo", asn.Nickname)
}

func TestGtm_CreateAsMap(t *testing.T) {
	var result AsMapResponse
	var req AsMap

	respData, err := loadTestData("TestGtm_CreateAsMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateAsMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		asmap            *AsMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *AsMapResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			asmap:      &req,
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
			asmap:          &req,
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
			result, err := client.CreateAsMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.asmap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGtm_UpdateAsMap(t *testing.T) {
	var result AsMapResponse
	var req AsMap

	respData, err := loadTestData("TestGtm_CreateAsMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateAsMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		asmap            *AsMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			asmap:      &req,
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
			asmap:          &req,
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
			result, err := client.UpdateAsMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.asmap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGtm_DeleteAsMap(t *testing.T) {
	var result AsMapResponse
	var req AsMap

	respData, err := loadTestData("TestGtm_CreateAsMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateAsMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		asmap            *AsMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			asmap:      &req,
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
			asmap:          &req,
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
			result, err := client.DeleteAsMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.asmap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
