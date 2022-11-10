package gtm

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGtm_NewDomain(t *testing.T) {
	client := Client(session.Must(session.New()))

	dom := client.NewDomain(context.Background(), "example.com", "primary")

	assert.Equal(t, "example.com", dom.Name)
	assert.Equal(t, "primary", dom.Type)
}

// Verify GetListDomains. Sould pass, e.g. no API errors and non nil list.
func TestGtm_ListDomains(t *testing.T) {
	var result DomainsList

	respData, err := loadTestData("TestGtm_ListDomains.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []*DomainItem
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/config-gtm/v1/domains",
			expectedResponse: result.DomainItems,
		},
		"500 internal server error": {
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching domains",
    "status": 500
}`,
			expectedPath: "/config-gtm/v1/domains",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching domains",
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
			result, err := client.ListDomains(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)))
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGtm_NullFieldMap(t *testing.T) {
	var result NullFieldMapStruct

	gob.Register(map[string]NullPerObjectAttributeStruct{})

	respData, err := loadTestData("TestGtm_NullFieldMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	resultData, err := loadTestData("TestGtm_NullFieldMap.result.gob")
	if err != nil {
		t.Fatal(err)
	}

	if err := gob.NewDecoder(bytes.NewBuffer(resultData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		arg              *Domain
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *NullFieldMapStruct
		withError        error
	}{
		"200 OK": {
			arg: &Domain{
				Name: "example.akadns.net",
				Type: "primary",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net",
			expectedResponse: &result,
		},
		"500 internal server error": {
			arg: &Domain{
				Name: "example.akadns.net",
				Type: "primary",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching null field map",
    "status": 500
}`),
			expectedPath: "/config-gtm/v1/domains/example.akadns.net",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching null field map",
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
			result, err := client.NullFieldMap(context.Background(), test.arg)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test GetDomain
// GetDomain(context.Context, string) (*Domain, error)
func TestGtm_GetDomain(t *testing.T) {
	var result Domain

	respData, err := loadTestData("TestGtm_GetDomain.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		domain           string
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *Domain
		withError        error
	}{
		"200 OK": {
			domain:           "example.akadns.net",
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net",
			expectedResponse: &result,
		},
		"500 internal server error": {
			domain:         "example.akadns.net",
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching domain"
}`),
			expectedPath: "/config-gtm/v1/domains/example.akadns.net",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching domain",
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
			result, err := client.GetDomain(context.Background(), test.domain)
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
// CreateDomain(context.Context, *Domain, map[string]string) (*DomainResponse, error)
func TestGtm_CreateDomain(t *testing.T) {
	var result DomainResponse

	respData, err := loadTestData("TestGtm_GetDomain.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		domain           Domain
		query            map[string]string
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *DomainResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			domain: Domain{
				Name: "gtmdomtest.akadns.net",
				Type: "basic",
			},
			query: map[string]string{"contractId": "1-2ABCDE"},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/config-gtm/v1/domains?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			domain: Domain{
				Name: "gtmdomtest.akadns.net",
				Type: "basic",
			},
			query:          map[string]string{"contractId": "1-2ABCDE"},
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating domain"
}`),
			expectedPath: "/config-gtm/v1/domains?contractId=1-2ABCDE",
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
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write(test.responseBody)
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateDomain(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), &test.domain, test.query)
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
// UpdateDomain(context.Context, *Domain, map[string]string) (*DomainResponse, error)
func TestGtm_UpdateDomain(t *testing.T) {
	var result DomainResponse

	respData, err := loadTestData("TestGtm_UpdateDomain.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		domain           Domain
		query            map[string]string
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *ResponseStatus
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			domain: Domain{
				EndUserMappingEnabled: false,
				Name:                  "gtmdomtest.akadns.net",
				Type:                  "basic",
			},
			query: map[string]string{"contractId": "1-2ABCDE"},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: result.Status,
			expectedPath:     "/config-gtm/v1/domains?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			domain: Domain{
				Name: "gtmdomtest.akadns.net",
				Type: "basic",
			},
			query:          map[string]string{"contractId": "1-2ABCDE"},
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/config-gtm/v1/domains?contractId=1-2ABCDE",
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
			result, err := client.UpdateDomain(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), &test.domain, test.query)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

/* Future. Presently no domain Delete endpoint.
func TestGtm_DeleteDomain(t *testing.T) {

        defer gock.Off()

        mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/"+gtmTestDomain)
        mock.
                Delete("/config-gtm/v1/domains/"+gtmTestDomain).
                HeaderPresent("Authorization").
                Reply(200).
                SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
                BodyString(`{
                        "resource" : null,
                        "status" : {
                               "changeId": "40e36abd-bfb2-4635-9fca-62175cf17007",
                               "links": [
                                     {
                                          "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
                                          "rel": "self"
                                     }
                               ],
                               "message": "Change Pending",
                               "passingValidation": true,
                               "propagationStatus": "PENDING",
                               "propagationStatusDate": "2019-04-25T14:54:00.000+00:00"
                          },
                }`)

        Init(config)

        getDomain := instantiateDomain()

        _, err := getDomain.Delete()
        assert.NoError(t, err)

}
*/
