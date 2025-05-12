package gtm

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGTM_ListDomains(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []DomainItem
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus: http.StatusOK,
			responseBody: `{
			"items":[{
				"acgId": "1-2345",
				"lastModified": "2014-03-03T16:02:45.000+0000",
				"name": "example.akadns.net",
				"status": "2014-02-20 22:56 GMT: Current configuration has been propagated to all GTM name servers",
				"lastModifiedBy": "test-user",
				"changeId": "abf5b76f-f9de-4404-bb2c-9d15e7b9ff5d",
				"activationState": "COMPLETE",
            	"modificationComments": "terraform test gtm domain",
            	"signAndServe": false,
            	"signAndServeAlgorithm": null,
            	"deleteRequestId": null,
				"links": [{
					"href": "/config-gtm/v1/domains/example.akadns.net",
					"rel": "self"
				}]
			},
			{
				"acgId": "1-2345",
				"lastModified": "2013-11-09T12:04:45.000+0000",
				"name": "demo.akadns.net",
				"status": "2014-02-20 22:56 GMT: Current configuration has been propagated to all GTM name servers",
 				"lastModifiedBy": "test-user",
				"changeId": "abf5b76f-f9de-4404-bb2c-9d15e7b9ff5d",
            	"activationState": "COMPLETE",
            	"modificationComments": "terraform test gtm domain",
           		"signAndServe": false,
            	"signAndServeAlgorithm": null,
            	"deleteRequestId": null,
				"links": [{
					"href": "/config-gtm/v1/domains/example.akadns.net",
					"rel": "self"
				}]
			}]}`,
			expectedPath: "/config-gtm/v1/domains",
			expectedResponse: []DomainItem{{
				AcgID:                 "1-2345",
				LastModified:          "2014-03-03T16:02:45.000+0000",
				Name:                  "example.akadns.net",
				Status:                "2014-02-20 22:56 GMT: Current configuration has been propagated to all GTM name servers",
				LastModifiedBy:        "test-user",
				ChangeID:              "abf5b76f-f9de-4404-bb2c-9d15e7b9ff5d",
				ActivationState:       "COMPLETE",
				ModificationComments:  "terraform test gtm domain",
				SignAndServe:          false,
				SignAndServeAlgorithm: "",
				DeleteRequestID:       "",
				Links: []Link{{
					Rel:  "self",
					Href: "/config-gtm/v1/domains/example.akadns.net",
				}},
			},
				{
					AcgID:                 "1-2345",
					LastModified:          "2013-11-09T12:04:45.000+0000",
					Name:                  "demo.akadns.net",
					Status:                "2014-02-20 22:56 GMT: Current configuration has been propagated to all GTM name servers",
					LastModifiedBy:        "test-user",
					ChangeID:              "abf5b76f-f9de-4404-bb2c-9d15e7b9ff5d",
					ActivationState:       "COMPLETE",
					ModificationComments:  "terraform test gtm domain",
					SignAndServe:          false,
					SignAndServeAlgorithm: "",
					DeleteRequestID:       "",
					Links: []Link{{
						Rel:  "self",
						Href: "/config-gtm/v1/domains/example.akadns.net",
					}},
				}},
		},
		"500 internal server error": {
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
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
		"Service Unavailable plain text response": {
			headers:        http.Header{},
			responseStatus: http.StatusServiceUnavailable,
			responseBody:   `Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z`,
			expectedPath:   "/config-gtm/v1/domains",
			withError: &Error{
				Type:       "",
				Title:      "Failed to unmarshal error body. GTM API failed. Check details for more information.",
				Detail:     "Your request did not succeed as this operation has reached  the limit for your account. Please try after 2024-01-16T15:20:55.945Z",
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		"Service Unavailable html response": {
			headers:        http.Header{},
			responseStatus: http.StatusServiceUnavailable,
			responseBody:   `<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>`,
			expectedPath:   "/config-gtm/v1/domains",
			withError: &Error{
				Type:       "",
				Title:      "Failed to unmarshal error body. GTM API failed. Check details for more information.",
				Detail:     "<HTML><HEAD>...</HEAD><BODY>...</BODY></HTML>",
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		"Service Unavailable xml response": {
			headers:        http.Header{},
			responseStatus: http.StatusServiceUnavailable,
			responseBody:   `<?xml version="1.0" encoding="UTF-8"?><root><item><id>1</id><name>Item 1</name></item><item><id>2</id><name>Item 2</name></item></root>`,
			expectedPath:   "/config-gtm/v1/domains",
			withError: &Error{
				Type:       "",
				Title:      "Failed to unmarshal error body. GTM API failed. Check details for more information.",
				Detail:     "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root><item><id>1</id><name>Item 1</name></item><item><id>2</id><name>Item 2</name></item></root>",
				StatusCode: http.StatusServiceUnavailable,
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

func TestGTM_NullFieldMap(t *testing.T) {
	var result NullFieldMapStruct

	gob.Register(map[string]NullPerObjectAttributeStruct{})

	respData, err := loadTestData("TestGTM_NullFieldMap.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	resultData, err := loadTestData("TestGTM_NullFieldMap.result.gob")
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

func TestGTM_GetDomain(t *testing.T) {
	var result GetDomainResponse

	respData, err := loadTestData("TestGTM_GetDomain.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           GetDomainRequest
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *GetDomainResponse
		withError        error
	}{
		"200 OK": {
			params: GetDomainRequest{
				DomainName: "example.akadns.net",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/config-gtm/v1/domains/example.akadns.net",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetDomainRequest{
				DomainName: "example.akadns.net",
			},
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
			result, err := client.GetDomain(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_CreateDomain(t *testing.T) {
	var result CreateDomainResponse

	respData, err := loadTestData("TestGTM_GetDomain.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           CreateDomainRequest
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *CreateDomainResponse
		withError        error
		headers          http.Header
	}{
		"201 Created": {
			params: CreateDomainRequest{
				Domain: &Domain{
					Name: "gtmdomtest.akadns.net",
					Type: "basic",
				},
				QueryArgs: &DomainQueryArgs{ContractID: "1-2ABCDE"},
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/config-gtm/v1/domains?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			params: CreateDomainRequest{
				Domain: &Domain{
					Name: "gtmdomtest.akadns.net",
					Type: "basic",
				},
				QueryArgs: &DomainQueryArgs{ContractID: "1-2ABCDE"},
			},
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
				assert.Equal(t, test.expectedPath, r.URL.String())
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

func TestGTM_UpdateDomain(t *testing.T) {
	var result UpdateDomainResponse

	respData, err := loadTestData("TestGTM_UpdateDomain.resp.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(respData)).Decode(&result); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params           UpdateDomainRequest
		responseStatus   int
		responseBody     []byte
		expectedPath     string
		expectedResponse *UpdateDomainResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateDomainRequest{
				Domain: &Domain{
					EndUserMappingEnabled: false,
					Name:                  "gtmdomtest.akadns.net",
					Type:                  "basic",
				},
				QueryArgs: &DomainQueryArgs{ContractID: "1-2ABCDE"},
			},
			headers: http.Header{
				"Content-Type": []string{"application/vnd.config-gtm.v1.4+json;charset=UTF-8"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/config-gtm/v1/domains/gtmdomtest.akadns.net?contractId=1-2ABCDE",
		},
		"500 internal server error": {
			params: UpdateDomainRequest{
				Domain: &Domain{
					EndUserMappingEnabled: false,
					Name:                  "gtmdomtest.akadns.net",
					Type:                  "basic",
				},
				QueryArgs: &DomainQueryArgs{ContractID: "1-2ABCDE"},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: []byte(`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/config-gtm/v1/domains/gtmdomtest.akadns.net?contractId=1-2ABCDE",
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
				assert.Equal(t, test.expectedPath, r.URL.String())
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

func TestGTM_DeleteDomains(t *testing.T) {
	tests := map[string]struct {
		request             DeleteDomainsRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *DeleteDomainsResponse
		withError           func(*testing.T, error)
	}{
		"200 Successfully submitted the delete request for domains": {
			request: DeleteDomainsRequest{
				Body: DeleteDomainsRequestBody{
					DomainNames: []string{
						"example.akadns.net",
						"demo.akadns.net",
					},
				},
			},
			expectedRequestBody: `{"domains":["example.akadns.net","demo.akadns.net"]}`,
			responseStatus:      http.StatusOK,
			expectedPath:        "/config-gtm/v1/domains/delete-requests",
			responseBody: `{
						"requestId": "e585a640-0849-4b87-8dd9-91afdaf8851c",
						"expirationDate": "2021-01-03T12:00:00Z"
			}`,
			expectedResponse: &DeleteDomainsResponse{
				RequestID:      "e585a640-0849-4b87-8dd9-91afdaf8851c",
				ExpirationDate: "2021-01-03T12:00:00Z",
			},
		},
		"200 Successfully submitted the delete request for domains - with BypassSafetyChecks query param": {
			request: DeleteDomainsRequest{
				BypassSafetyChecks: ptr.To(true),
				Body: DeleteDomainsRequestBody{
					DomainNames: []string{
						"example.akadns.net",
						"demo.akadns.net",
					},
				},
			},
			expectedRequestBody: `{"domains":["example.akadns.net","demo.akadns.net"]}`,
			responseStatus:      http.StatusOK,
			expectedPath:        "/config-gtm/v1/domains/delete-requests?bypassSafetyChecks=true",
			responseBody: `{
						"requestId": "e585a640-0849-4b87-8dd9-91afdaf8851c",
						"expirationDate": "2021-01-03T12:00:00Z"
			}`,
			expectedResponse: &DeleteDomainsResponse{
				RequestID:      "e585a640-0849-4b87-8dd9-91afdaf8851c",
				ExpirationDate: "2021-01-03T12:00:00Z",
			},
		},
		"Validation error - missing domainNames in the request": {
			request: DeleteDomainsRequest{
				Body: DeleteDomainsRequestBody{
					DomainNames: []string{},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domains: struct validation: DomainNames: cannot be blank", err.Error())
			},
		},
		"Validation error - empty domainNames in the request": {
			request: DeleteDomainsRequest{
				Body: DeleteDomainsRequestBody{
					DomainNames: []string{"example.akadns.net", ""},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "delete domains: struct validation: 1: cannot be blank", err.Error())
			},
		},
		"400 Bad Request": {
			request: DeleteDomainsRequest{
				Body: DeleteDomainsRequestBody{
					DomainNames: []string{
						"example.akadns.net",
						"demo.akadns.net",
						"devexpautomatedtest_0muwrj",
					},
				},
			},
			expectedRequestBody: `{"domains":["example.akadns.net","demo.akadns.net","devexpautomatedtest_0muwrj"]}`,
			responseStatus:      http.StatusBadRequest,
			expectedPath:        "/config-gtm/v1/domains/delete-requests",
			responseBody: `{
						"type": "https://problems.luna.akamaiapis.net/config-gtm/v1/badRequest",
    					"title": "Bad Request",
    					"instance": "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm-api/v1/domains/delete-requests?bypassSafetyChecks=false#724a2c56-a67a-4ea0-81df-50148d335a86",
    					"status": 400,
    					"detail": "some of the listed domains could not be found",
    					"problemId": "724a2c56-a67a-4ea0-81df-50148d335a86",
						"missingZones": [
        					"devexpautomatedtest_0muwrj"
    					]
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrDomainNotFound))
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.DeleteDomains(context.Background(), test.request)

			if test.withError != nil {
				test.withError(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGTM_GetDeleteDomainsStatus(t *testing.T) {
	tests := map[string]struct {
		request          DeleteDomainsStatusRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *DeleteDomainsStatusResponse
		withError        func(*testing.T, error)
	}{
		"200 Successfully gets the delete domain request status": {
			request: DeleteDomainsStatusRequest{
				RequestID: "e585a640-0849-4b87-8dd9-91afdaf8851c",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/config-gtm/v1/domains/delete-requests/e585a640-0849-4b87-8dd9-91afdaf8851c",
			responseBody: `{
							"requestId": "e585a640-0849-4b87-8dd9-91afdaf8851c",
							"domainsSubmitted": 3,
							"successCount": 2,
							"failureCount": 1,
							"isComplete": true,
							"expirationDate": "2021-01-03T12:00:00Z"
			}`,
			expectedResponse: &DeleteDomainsStatusResponse{
				RequestID:        "e585a640-0849-4b87-8dd9-91afdaf8851c",
				DomainsSubmitted: 3,
				SuccessCount:     2,
				FailureCount:     1,
				IsComplete:       true,
				ExpirationDate:   "2021-01-03T12:00:00Z",
			},
		},
		"Validation error - missing or empty RequestID in the request": {
			request: DeleteDomainsStatusRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get delete domains status: validation failed: RequestID: cannot be blank", err.Error())
			},
		},
		"404 Request ID not found": {
			request: DeleteDomainsStatusRequest{
				RequestID: "e585a640-0849-4b87-8dd9-91afdaf8851c",
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/config-gtm/v1/domains/delete-requests/e585a640-0849-4b87-8dd9-91afdaf8851c",
			responseBody: `{
    					"type": "https://problems.luna.akamaiapis.net/config-gtm/v1/notfound",
    					"title": "Resource not found.",
    					"instance": "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm-api/v1/domains/delete-requests/e585a640-0849-4b87-8dd9-91afdaf8851c#9037df98-5e58-49a9-a74f-ec247ec920f7",
    					"status": 404,
    					"detail": "Resource not found.",
    					"problemId": "9037df98-5e58-49a9-a74f-ec247ec920f7"
			}`,
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrNotFound))
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
			defer mockServer.Close()

			client := mockAPIClient(t, mockServer)
			result, err := client.GetDeleteDomainsStatus(context.Background(), test.request)

			if test.withError != nil {
				test.withError(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
