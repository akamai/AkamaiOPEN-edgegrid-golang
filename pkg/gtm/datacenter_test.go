package gtm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGTM_ListDatacenters(t *testing.T) {
	var result DatacenterList

	respData, err := loadTestData("TestGTM_ListDatacenters.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           ListDatacentersRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Datacenter
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: ListDatacentersRequest{
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters",
			expectedResponse: result.DatacenterItems,
		},
		"500 internal server error": {
			params: ListDatacentersRequest{
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching datacenters"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/datacenters",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching datacenters",
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
			result, err := client.ListDatacenters(
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

func TestGTM_GetDatacenter(t *testing.T) {
	var result Datacenter

	respData, err := loadTestData("TestGTM_GetDatacenter.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           GetDatacenterRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Datacenter
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetDatacenterRequest{
				DatacenterID: 1,
				DomainName:   "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/1",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetDatacenterRequest{
				DatacenterID: 1,
				DomainName:   "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching datacenter"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/datacenters/1",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching datacenter",
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
			result, err := client.GetDatacenter(
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

func TestGTM_CreateDatacenter(t *testing.T) {
	var result CreateDatacenterResponse
	var req Datacenter

	respData, err := loadTestData("TestGTM_CreateDatacenter.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateDatacenter.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           CreateDatacenterRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreateDatacenterResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: CreateDatacenterRequest{
				Datacenter: &req,
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: CreateDatacenterRequest{
				Datacenter: &req,
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating dc"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/datacenters",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating dc",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateDatacenter(
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

func TestGTM_CreateMapsDefaultDatacenter(t *testing.T) {
	var result DatacenterResponse

	respData, err := loadTestData("TestGTM_CreateMapsDefaultDatacenter.resp.json")
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
		expectedResponse *Datacenter
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/default-datacenter-for-maps",
			expectedResponse: result.Resource,
		},
		"500 internal server error": {
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating dc"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/datacenters/default-datacenter-for-maps",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating dc",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusNotFound)
					_, err = w.Write([]byte(`
                                        {
                                            "type": "Datacenter",
                                            "title": "not found"
                                        }`))
					require.NoError(t, err)
					return
				}
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateMapsDefaultDatacenter(
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

func TestGTM_CreateIPv4DefaultDatacenter(t *testing.T) {
	var result DatacenterResponse

	respData, err := loadTestData("TestGTM_CreateIPv4DefaultDatacenter.resp.json")
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
		expectedResponse *Datacenter
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/datacenter-for-ip-version-selector-ipv4",
			expectedResponse: result.Resource,
		},
		"500 internal server error": {
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating dc"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/datacenters/datacenter-for-ip-version-selector-ipv4",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating dc",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusNotFound)
					_, err = w.Write([]byte(`
                                        {
                                            "type": "Datacenter",
                                            "title": "not found"
                                        }`))
					require.NoError(t, err)
					return
				}
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateIPv4DefaultDatacenter(
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

func TestGTM_CreateIPv6DefaultDatacenter(t *testing.T) {
	var result DatacenterResponse

	respData, err := loadTestData("TestGTM_CreateIPv6DefaultDatacenter.resp.json")
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
		expectedResponse *Datacenter
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/datacenter-for-ip-version-selector-ipv6",
			expectedResponse: result.Resource,
		},
		"500 internal server error": {
			domainName:     "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating dc"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/datacenters/datacenter-for-ip-version-selector-ipv6",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating dc",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusNotFound)
					_, err = w.Write([]byte(`
                                        {
                                            "type": "Datacenter",
                                            "title": "not found"
                                        }`))
					require.NoError(t, err)
					return
				}
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateIPv6DefaultDatacenter(
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

func TestGTM_UpdateDatacenter(t *testing.T) {
	var result UpdateDatacenterResponse
	var req Datacenter

	respData, err := loadTestData("TestGTM_CreateDatacenter.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateDatacenter.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           UpdateDatacenterRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateDatacenterResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: UpdateDatacenterRequest{
				Datacenter: &req,
				DomainName: "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/0",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: UpdateDatacenterRequest{
				Datacenter: &req,
				DomainName: "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating dc"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/datacenters/0",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating dc",
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
			result, err := client.UpdateDatacenter(
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

func TestGTM_DeleteDatacenter(t *testing.T) {
	var result DeleteDatacenterResponse
	var req Datacenter

	respData, err := loadTestData("TestGTM_CreateDatacenter.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGTM_CreateDatacenter.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           DeleteDatacenterRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DeleteDatacenterResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: DeleteDatacenterRequest{
				DatacenterID: 1,
				DomainName:   "example.akadns.net",
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/1",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: DeleteDatacenterRequest{
				DatacenterID: 1,
				DomainName:   "example.akadns.net",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error updating dc"
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/datacenters/1",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error updating dc",
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
			result, err := client.DeleteDatacenter(
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
