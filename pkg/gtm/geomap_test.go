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

func TestGTM_ListGeoMap(t *testing.T) {
	var result GeoMapList

	respData, err := loadTestData("TestGTM_ListGeoMap.resp.json")
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
		expectedResponse []*GeoMap
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
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/geographic-maps",
			expectedResponse: result.GeoMapItems,
		},
		"500 internal server error": {
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching GeoMap",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/geographic-maps",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching GeoMap",
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
			result, err := client.ListGeoMaps(
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

func TestGTM_GetGeoMap(t *testing.T) {
	var result GeoMap

	respData, err := loadTestData("TestGTM_GetGeoMap.resp.json")
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
		expectedResponse *GeoMap
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
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/geographic-maps/Software-rollout",
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
    "detail": "Error fetching GeoMap",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/geographic-maps/Software-rollout",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching GeoMap",
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
			result, err := client.GetGeoMap(
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

func TestGTM_CreateGeoMap(t *testing.T) {
	var result GeoMapResponse
	var req GeoMap

	respData, err := loadTestData("TestGTM_CreateGeoMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateGeoMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		GeoMap           *GeoMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GeoMapResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			GeoMap:     &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
			expectedResponse: &result,
		},
		"500 internal server error": {
			GeoMap:         &req,
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating GeoMap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating GeoMap",
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
			result, err := client.CreateGeoMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.GeoMap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_UpdateGeoMap(t *testing.T) {
	var result GeoMapResponse
	var req GeoMap

	respData, err := loadTestData("TestGTM_CreateGeoMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateGeoMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		GeoMap           *GeoMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			GeoMap:     &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
			expectedResponse: result.Status,
		},
		"500 internal server error": {
			GeoMap:         &req,
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating GeoMap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating GeoMap",
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
			result, err := client.UpdateGeoMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.GeoMap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_DeleteGeoMap(t *testing.T) {
	var result GeoMapResponse
	var req GeoMap

	respData, err := loadTestData("TestGTM_CreateGeoMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateGeoMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		GeoMap           *GeoMap
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			GeoMap:     &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
			expectedResponse: result.Status,
		},
		"500 internal server error": {
			GeoMap:         &req,
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating GeoMap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating GeoMap",
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
			result, err := client.DeleteGeoMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.GeoMap, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
