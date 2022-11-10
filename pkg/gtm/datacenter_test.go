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

func TestGtm_NewDatacenterResponse(t *testing.T) {
	client := Client(session.Must(session.New()))

	dcr := client.NewDatacenterResponse(context.Background())

	assert.NotNil(t, dcr)
}

func TestGtm_NewDatacenter(t *testing.T) {
	client := Client(session.Must(session.New()))

	dc := client.NewDatacenter(context.Background())

	assert.NotNil(t, dc)
}

func TestGtm_ListDatacenters(t *testing.T) {
	var result DatacenterList

	respData, err := loadTestData("TestGtm_ListDatacenters.resp.json")
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
		expectedResponse []*Datacenter
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
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters",
			expectedResponse: result.DatacenterItems,
		},
		"500 internal server error": {
			domainName:     "example.akadns.net",
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

func TestGtm_GetDatacenter(t *testing.T) {
	var result Datacenter

	respData, err := loadTestData("TestGtm_GetDatacenter.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		id               int
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Datacenter
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			id:         1,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/1",
			expectedResponse: &result,
		},
		"500 internal server error": {
			id:             1,
			domainName:     "example.akadns.net",
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
					session.WithContextHeaders(test.headers)), test.id, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGtm_CreateDatacenter(t *testing.T) {
	var result DatacenterResponse
	var req Datacenter

	respData, err := loadTestData("TestGtm_CreateDatacenter.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateDatacenter.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		dc               *Datacenter
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DatacenterResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			dc:         &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters",
			expectedResponse: &result,
		},
		"500 internal server error": {
			dc:             &req,
			domainName:     "example.akadns.net",
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
					session.WithContextHeaders(test.headers)), test.dc, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGtm_CreateMapsDefaultDatacenter(t *testing.T) {
	var result DatacenterResponse

	respData, err := loadTestData("TestGtm_CreateMapsDefaultDatacenter.resp.json")
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
					w.Write([]byte(`
                                        {
                                            "type": "Datacenter",
                                            "title": "not found"
                                        }`))
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

func TestGtm_CreateIPv4DefaultDatacenter(t *testing.T) {
	var result DatacenterResponse

	respData, err := loadTestData("TestGtm_CreateIPv4DefaultDatacenter.resp.json")
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
					w.Write([]byte(`
                                        {
                                            "type": "Datacenter",
                                            "title": "not found"
                                        }`))
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

func TestGtm_CreateIPv6DefaultDatacenter(t *testing.T) {
	var result DatacenterResponse

	respData, err := loadTestData("TestGtm_CreateIPv6DefaultDatacenter.resp.json")
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
					w.Write([]byte(`
                                        {
                                            "type": "Datacenter",
                                            "title": "not found"
                                        }`))
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

func TestGtm_UpdateDatacenter(t *testing.T) {
	var result DatacenterResponse
	var req Datacenter

	respData, err := loadTestData("TestGtm_CreateDatacenter.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateDatacenter.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		dc               *Datacenter
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			dc:         &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/0",
			expectedResponse: result.Status,
		},
		"500 internal server error": {
			dc:             &req,
			domainName:     "example.akadns.net",
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
					session.WithContextHeaders(test.headers)), test.dc, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGtm_DeleteDatacenter(t *testing.T) {
	var result DatacenterResponse
	var req Datacenter

	respData, err := loadTestData("TestGtm_CreateDatacenter.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateDatacenter.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		dc               *Datacenter
		domainName       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			dc:         &req,
			domainName: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/datacenters/0",
			expectedResponse: result.Status,
		},
		"500 internal server error": {
			dc:             &req,
			domainName:     "example.akadns.net",
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
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteDatacenter(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.dc, test.domainName)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
