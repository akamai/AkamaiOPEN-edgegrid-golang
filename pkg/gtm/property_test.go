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

func TestGtm_NewTrafficTarget(t *testing.T) {
	client := Client(session.Must(session.New()))

	tgt := client.NewTrafficTarget(context.Background())

	assert.NotNil(t, tgt)
}

func TestGtm_NewStaticRRSet(t *testing.T) {
	client := Client(session.Must(session.New()))

	set := client.NewStaticRRSet(context.Background())

	assert.NotNil(t, set)
}

func TestGtm_NewLivenessTest(t *testing.T) {
	client := Client(session.Must(session.New()))

	test := client.NewLivenessTest(context.Background(), "foo", "bar", 1, 1000)

	assert.NotNil(t, test)
	assert.Equal(t, "foo", test.Name)
	assert.Equal(t, "bar", test.TestObjectProtocol)
	assert.Equal(t, 1, test.TestInterval)
	assert.Equal(t, float32(1000), test.TestTimeout)
}

func TestGtm_NewProperty(t *testing.T) {
	client := Client(session.Must(session.New()))

	prop := client.NewProperty(context.Background(), "foo")

	assert.NotNil(t, prop)
	assert.Equal(t, prop.Name, "foo")
}

func TestGtm_ListProperties(t *testing.T) {
	var result PropertyList

	respData, err := loadTestData("TestGtm_ListProperties.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		domain           string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []*Property
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			domain: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/properties",
			expectedResponse: result.PropertyItems,
		},
		"500 internal server error": {
			domain:         "example.akadns.net",
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching propertys",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching propertys",
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
			result, err := client.ListProperties(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.domain)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test GetProperty
// GetProperty(context.Context, string) (*Property, error)
func TestGtm_GetProperty(t *testing.T) {
	var result Property

	respData, err := loadTestData("TestGtm_GetProperty.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		name             string
		domain           string
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *Property
		withError        error
	}{
		"200 OK": {
			name:             "www",
			domain:           "example.akadns.net",
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net/properties/www",
			expectedResponse: &result,
		},
		"500 internal server error": {
			name:           "www",
			domain:         "example.akadns.net",
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching property"
}`),
			expectedPath: "/config-gtm/v1/domains/example.akadns.net/properties/www",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching property",
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
				_, err := w.Write(test.responseBody)
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetProperty(context.Background(), test.name, test.domain)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Create domain.
// CreateProperty(context.Context, *Property, map[string]string) (*PropertyResponse, error)
func TestGtm_CreateProperty(t *testing.T) {
	var result PropertyResponse
	var req Property

	respData, err := loadTestData("TestGtm_CreateProperty.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateProperty.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		domain           string
		prop             *Property
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *PropertyResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			prop:   &req,
			domain: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			prop:           &req,
			domain:         "example.akadns.net",
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating domain"
}`),
			expectedPath: "/config-gtm/v1/domains/example.akadns.net?contractId=1-2ABCDE",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating domain",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write(test.responseBody)
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateProperty(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.prop, test.domain)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update domain.
// UpdateProperty(context.Context, *Property, map[string]string) (*PropertyResponse, error)
func TestGtm_UpdateProperty(t *testing.T) {
	var result PropertyResponse
	var req Property

	respData, err := loadTestData("TestGtm_CreateProperty.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateProperty.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		prop             *Property
		domain           string
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			prop:   &req,
			domain: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: result.Status,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			prop:           &req,
			domain:         "example.akadns.net",
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/config-gtm/v1/domains/example.akadns.net?contractId=1-2ABCDE",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write(test.responseBody)
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateProperty(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.prop, test.domain)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGtm_DeleteProperty(t *testing.T) {
	var result PropertyResponse
	var req Property

	respData, err := loadTestData("TestGtm_CreateProperty.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	reqData, err := loadTestData("TestGtm_CreateProperty.req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(reqData)).Decode(&req); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		prop             *Property
		domain           string
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			prop:   &req,
			domain: "example.akadns.net",
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: result.Status,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			prop:           &req,
			domain:         "example.akadns.net",
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/config-gtm/v1/domains/example.akadns.net?contractId=1-2ABCDE",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write(test.responseBody)
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteProperty(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.prop, test.domain)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
