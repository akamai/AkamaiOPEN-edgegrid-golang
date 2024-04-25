package gtm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGTM_ListCIDRMaps(t *testing.T) {
	var result CIDRMapList

	respData, err := loadTestData("TestGTM_ListCIDRMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           ListCIDRMapsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []CIDRMap
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: ListCIDRMapsRequest{
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/cidr-maps",
			expectedResponse: result.CIDRMapItems,
		},
		"500 internal server error": {
			params: ListCIDRMapsRequest{
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching asmap",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/cidr-maps",
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
			result, err := client.ListCIDRMaps(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_GetCIDRMap(t *testing.T) {
	var result GetCIDRMapResponse

	respData, err := loadTestData("TestGTM_GetCIDRMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           GetCIDRMapRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCIDRMapResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetCIDRMapRequest{
				MapName:    "Software-rollout",
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/cidr-maps/Software-rollout",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetCIDRMapRequest{
				MapName:    "Software-rollout",
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching asmap",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/cidr-maps/Software-rollout",
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
			result, err := client.GetCIDRMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_CreateCIDRMap(t *testing.T) {
	var result CreateCIDRMapResponse
	var req CIDRMap

	respData, err := loadTestData("TestGTM_CreateCIDRMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateCIDRMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params              CreateCIDRMapRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateCIDRMapResponse
		expectedRequestBody string
		withError           error
		headers             http.Header
	}{
		"200 OK": {
			params: CreateCIDRMapRequest{
				CIDR:       &req,
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/cidr-maps/The%20North",
			expectedResponse: &result,
			expectedRequestBody: `{
    "name": "The North",
    "defaultDatacenter": {
        "datacenterId": 5400,
        "nickname": "All Other CIDR Blocks"
    },
    "assignments": [
        {
            "datacenterId": 3134,
            "nickname": "Frostfangs and the Fist of First Men",
            "blocks": [
                "1.3.5.9",
                "1.2.3.0/24"
            ]
        },
        {
            "datacenterId": 3133,
            "nickname": "Winterfell",
            "blocks": [
                "1.2.4.0/24"
            ]
        }
    ]
}`,
		},
		"200 test": {
			params: CreateCIDRMapRequest{
				CIDR: &CIDRMap{
					DefaultDatacenter: &DatacenterBase{
						Nickname:     "test_datacenter",
						DatacenterID: 200,
					},
					Assignments: nil,
					Name:        "test_name",
					Links:       nil,
				},
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/cidr-maps/test_name",
			expectedResponse: &result,
			expectedRequestBody: `{
    "name": "test_name",
    "defaultDatacenter": {
        "datacenterId": 200,
        "nickname": "test_datacenter"
    }
}`,
		},
		"500 internal server error": {
			params: CreateCIDRMapRequest{
				CIDR:       &req,
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating asmap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/cidr-maps/The%20North",
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateCIDRMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_UpdateCIDRMap(t *testing.T) {
	var result UpdateCIDRMapResponse
	var req CIDRMap

	respData, err := loadTestData("TestGTM_CreateCIDRMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateCIDRMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           UpdateCIDRMapRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateCIDRMapResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: UpdateCIDRMapRequest{
				CIDR:       &req,
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/cidr-maps/The%20North",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: UpdateCIDRMapRequest{
				CIDR:       &req,
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating asmap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/cidr-maps/The%20North",
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
			result, err := client.UpdateCIDRMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_DeleteCIDRMap(t *testing.T) {
	var result DeleteCIDRMapResponse
	var req CIDRMap

	respData, err := loadTestData("TestGTM_CreateCIDRMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateCIDRMap.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           DeleteCIDRMapRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DeleteCIDRMapResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: DeleteCIDRMapRequest{
				MapName:    "The%20North",
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/cidr-maps/The%20North",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: DeleteCIDRMapRequest{
				MapName:    "The%20North",
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating asmap"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/cidr-maps/The%20North",
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
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteCIDRMap(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
