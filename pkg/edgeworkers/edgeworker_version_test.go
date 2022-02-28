package edgeworkers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEdgeWorkerVersion(t *testing.T) {
	tests := map[string]struct {
		params           GetEdgeWorkerVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *EdgeWorkerVersion
		withError        error
	}{
		"200 OK - get EdgeWorkerVersion": {
			params: GetEdgeWorkerVersionRequest{
				EdgeWorkerID: 12345,
				Version:      "1.2.3",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "edgeWorkerId": 12345,
	"version": "1.2.3",
    "accountId": "B-123-WNKA6P",
	"checksum": "868f28f16c26f46d418d83e24973520534d9ea4e4dbfd8a69ab00c1c37f28ca4",
	"sequenceNumber": 3,
    "createdBy": "jbond",
    "createdTime": "2021-04-19T07:08:37Z"
}
`,
			expectedPath: "/edgeworkers/v1/ids/12345/versions/1.2.3",
			expectedResponse: &EdgeWorkerVersion{
				EdgeWorkerID:   12345,
				Version:        "1.2.3",
				AccountID:      "B-123-WNKA6P",
				Checksum:       "868f28f16c26f46d418d83e24973520534d9ea4e4dbfd8a69ab00c1c37f28ca4",
				SequenceNumber: 3,
				CreatedBy:      "jbond",
				CreatedTime:    "2021-04-19T07:08:37Z",
			},
		},
		"500 internal server error": {
			params: GetEdgeWorkerVersionRequest{
				EdgeWorkerID: 12345,
				Version:      "1.2.3",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
  "title": "Server Error",
  "status": 500,
  "instance": "host_name/edgeworkers/v1/ids/12345/versions/1.2.3",
  "method": "GET",
  "serverIp": "10.0.0.1",
  "clientIp": "192.168.0.1",
  "requestId": "a73affa111",
  "requestTime": "2021-12-06T10:27:11Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/12345/versions/1.2.3",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/12345/versions/1.2.3",
				Method:      "GET",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T10:27:11Z",
			},
		},
		"missing EdgeWorkerID": {
			params: GetEdgeWorkerVersionRequest{
				Version: "1.2.3",
			},
			withError: ErrStructValidation,
		},
		"missing Version": {
			params: GetEdgeWorkerVersionRequest{
				EdgeWorkerID: 12345,
			},
			withError: ErrStructValidation,
		},
		"403 Forbidden - incorrect credentials": {
			params: GetEdgeWorkerVersionRequest{
				EdgeWorkerID: 12345,
				Version:      "1.2.3",
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/ids/12345/versions/1.2.3",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "GET",
    "serverIp": "10.0.0.1",
    "clientIp": "192.168.0.1",
    "requestId": "a73affa111",
    "requestTime": "2021-12-06T12:45:09Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/12345/versions/1.2.3",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids/12345/versions/1.2.3",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T12:45:09Z",
			},
		},
		"404 Not Found - EdgeWorkerID doesn't exist": {
			params: GetEdgeWorkerVersionRequest{
				EdgeWorkerID: 12345,
				Version:      "1.2.3",
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-not-found",
    "title": "The given resource could not be found.",
    "detail": "Unable to find the requested version",
    "instance": "/edgeworkers/error-instances/86d1cc10-4baf-49e1-b81a-075b72a2f6a4",
    "status": 404,
    "errorCode": "EW2002"
}`,
			expectedPath: "/edgeworkers/v1/ids/12345/versions/1.2.3",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-not-found",
				Title:     "The given resource could not be found.",
				Detail:    "Unable to find the requested version",
				Status:    404,
				Instance:  "/edgeworkers/error-instances/86d1cc10-4baf-49e1-b81a-075b72a2f6a4",
				ErrorCode: "EW2002",
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
			result, err := client.GetEdgeWorkerVersion(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListEdgeWorkerVersions(t *testing.T) {
	tests := map[string]struct {
		params           ListEdgeWorkerVersionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListEdgeWorkerVersionsResponse
		withError        error
	}{
		"200 OK - list Edgeworker versions": {
			params: ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: 88334,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "versions": [
        {
            "edgeWorkerId": 88334,
            "version": "1.23",
            "accountId": "B-3-WNKA6P",
            "checksum": "868f28f16c26f46d418d83e24973520534d9ea4e4dbfd8a69ab00c1c37f28ca4",
            "sequenceNumber": 3,
            "createdBy": "jsmith",
            "createdTime": "2021-12-20T08:28:37Z"
        },
        {
            "edgeWorkerId": 88334,
            "version": "1.24.5",
            "accountId": "B-3-WNKA6P",
            "checksum": "ad9c18a7f2ed5d7bbcd31c55b94a0a00ae1771c6a15fd9265aeae08f5ef41e1f",
            "sequenceNumber": 4,
            "createdBy": "jsmith",
            "createdTime": "2021-12-20T09:39:48Z"
        }
    ]
}`,
			expectedPath: "/edgeworkers/v1/ids/88334/versions",
			expectedResponse: &ListEdgeWorkerVersionsResponse{[]EdgeWorkerVersion{
				{
					EdgeWorkerID:   88334,
					Version:        "1.23",
					AccountID:      "B-3-WNKA6P",
					Checksum:       "868f28f16c26f46d418d83e24973520534d9ea4e4dbfd8a69ab00c1c37f28ca4",
					SequenceNumber: 3,
					CreatedBy:      "jsmith",
					CreatedTime:    "2021-12-20T08:28:37Z",
				},
				{
					EdgeWorkerID:   88334,
					Version:        "1.24.5",
					AccountID:      "B-3-WNKA6P",
					Checksum:       "ad9c18a7f2ed5d7bbcd31c55b94a0a00ae1771c6a15fd9265aeae08f5ef41e1f",
					SequenceNumber: 4,
					CreatedBy:      "jsmith",
					CreatedTime:    "2021-12-20T09:39:48Z",
				},
			}},
		},
		"200 OK - no Versions under EdgeworkerID": {
			params: ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: 88334,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
   "versions": []
}`,
			expectedPath:     "/edgeworkers/v1/ids/88334/versions",
			expectedResponse: &ListEdgeWorkerVersionsResponse{[]EdgeWorkerVersion{}},
		},
		"missing EdgeWorkerID": {
			params:    ListEdgeWorkerVersionsRequest{},
			withError: ErrStructValidation,
		},
		"500 internal server error": {
			params: ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: 88334,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
 "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
 "title": "Server Error",
 "status": 500,
 "instance": "host_name/edgeworkers/v1/ids/88334/versions",
 "method": "GET",
 "serverIp": "10.0.0.1",
 "clientIp": "192.168.0.1",
 "requestId": "a73affa111",
 "requestTime": "2021-12-06T10:27:11Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/88334/versions",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/88334/versions",
				Method:      "GET",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T10:27:11Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params: ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: 88334,
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
   "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
   "title": "Forbidden",
   "status": 403,
   "detail": "The client does not have the grant needed for the request",
   "instance": "host_name/edgeworkers/v1/ids/88334/versions",
   "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
   "method": "GET",
   "serverIp": "10.0.0.1",
   "clientIp": "192.168.0.1",
   "requestId": "a73affa111",
   "requestTime": "2021-12-13T12:01:17Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/88334/versions",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids/88334/versions",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-13T12:01:17Z",
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
			result, err := client.ListEdgeWorkerVersions(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetEdgeWorkerVersionContent(t *testing.T) {
	tests := map[string]struct {
		params         GetEdgeWorkerVersionContentRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"200 OK - get EdgeWorkerVersion": {
			params: GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: 88334,
				Version:      "1.23",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/edgeworkers/v1/ids/88334/versions/1.23/content",
		},
		"500 internal server error": {
			params: GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: 88334,
				Version:      "1.23",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
  "title": "Server Error",
  "status": 500,
  "instance": "host_name/edgeworkers/v1/ids/88334/versions/1.23/content",
  "method": "GET",
  "serverIp": "10.0.0.1",
  "clientIp": "192.168.0.1",
  "requestId": "a73affa111",
  "requestTime": "2021-12-06T10:27:11Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/88334/versions/1.23/content",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/88334/versions/1.23/content",
				Method:      "GET",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T10:27:11Z",
			},
		},
		"missing EdgeWorkerID": {
			params: GetEdgeWorkerVersionContentRequest{
				Version: "1.23",
			},
			withError: ErrStructValidation,
		},
		"missing Version": {
			params: GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: 88334,
			},
			withError: ErrStructValidation,
		},
		"403 Forbidden - incorrect credentials": {
			params: GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: 88334,
				Version:      "1.23",
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/ids/88334/versions/1.23/content",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "GET",
    "serverIp": "10.0.0.1",
    "clientIp": "192.168.0.1",
    "requestId": "a73affa111",
    "requestTime": "2021-12-06T12:45:09Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/88334/versions/1.23/content",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids/88334/versions/1.23/content",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T12:45:09Z",
			},
		},
		"404 Not Found": {
			params: GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: 88334,
				Version:      "1.23",
			},
			expectedPath:   "/edgeworkers/v1/ids/88334/versions/1.23/content",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-not-found",
    "title": "The given resource could not be found.",
    "detail": "Unable to delete the version provided",
    "instance": "/edgeworkers/error-instances/514139f4-1608-4afc-88ac-67da91696af3",
    "status": 404,
    "errorCode": "EW2002"
}`,
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-not-found",
				Title:     "The given resource could not be found.",
				Detail:    "Unable to delete the version provided",
				Instance:  "/edgeworkers/error-instances/514139f4-1608-4afc-88ac-67da91696af3",
				Status:    404,
				ErrorCode: "EW2002",
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
			_, err := client.GetEdgeWorkerVersionContent(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				if test.responseStatus != 0 {
					assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				}

				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCreateEdgeWorkerVersion(t *testing.T) {
	tests := map[string]struct {
		params           CreateEdgeWorkerVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *EdgeWorkerVersion
		withError        error
	}{
		"201 Created - create EdgeWorker Version": {
			params: CreateEdgeWorkerVersionRequest{
				EdgeWorkerID:  88334,
				ContentBundle: Bundle{bytes.NewBuffer([]byte("testing create"))},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "edgeWorkerId": 88334,
    "version": "1.23",
    "accountId": "B-3-WNKA6P",
    "checksum": "868f28f16c26f46d418d83e24973520534d9ea4e4dbfd8a69ab00c1c37f28ca4",
    "sequenceNumber": 24,
    "createdBy": "jsmith",
    "createdTime": "2021-12-21T12:57:54Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/88334/versions",
			expectedResponse: &EdgeWorkerVersion{
				EdgeWorkerID:   88334,
				Version:        "1.23",
				AccountID:      "B-3-WNKA6P",
				Checksum:       "868f28f16c26f46d418d83e24973520534d9ea4e4dbfd8a69ab00c1c37f28ca4",
				SequenceNumber: 24,
				CreatedBy:      "jsmith",
				CreatedTime:    "2021-12-21T12:57:54Z",
			},
		},
		"missing EdgeWorkerID": {
			params: CreateEdgeWorkerVersionRequest{
				ContentBundle: Bundle{bytes.NewBuffer([]byte("testing create"))},
			},
			withError: ErrStructValidation,
		},
		"missing ContentBundle": {
			params: CreateEdgeWorkerVersionRequest{
				EdgeWorkerID: 88334,
			},
			withError: ErrStructValidation,
		},
		"500 internal server error": {
			params: CreateEdgeWorkerVersionRequest{
				EdgeWorkerID:  88334,
				ContentBundle: Bundle{bytes.NewBuffer([]byte("testing create"))},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
"type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
"title": "Server Error",
"status": 500,
"instance": "host_name/edgeworkers/v1/ids/88334/versions",
"method": "POST",
"serverIp": "10.0.0.1",
"clientIp": "192.168.0.1",
"requestId": "a73affa111",
"requestTime": "2021-12-13T13:32:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/88334/versions",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/88334/versions",
				Method:      "POST",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-13T13:32:37Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params: CreateEdgeWorkerVersionRequest{
				EdgeWorkerID:  88334,
				ContentBundle: Bundle{bytes.NewBuffer([]byte("testing create"))},
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
  "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
  "title": "Forbidden",
  "status": 403,
  "detail": "The client does not have the grant needed for the request",
  "instance": "host_name/edgeworkers/v1/ids/88334/versions",
  "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
  "method": "POST",
  "serverIp": "10.0.0.1",
  "clientIp": "192.168.0.1",
  "requestId": "a73affa111",
  "requestTime": "2021-12-13T13:32:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/88334/versions",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids/88334/versions",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "POST",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-13T13:32:37Z",
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
			result, err := client.CreateEdgeWorkerVersion(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteEdgeWorkerVersion(t *testing.T) {
	tests := map[string]struct {
		params         DeleteEdgeWorkerVersionRequest
		withError      error
		expectedPath   string
		responseStatus int
		responseBody   string
	}{
		"204 Deleted": {
			params: DeleteEdgeWorkerVersionRequest{
				EdgeWorkerID: 88334,
				Version:      "1.23",
			},
			expectedPath:   "/edgeworkers/v1/ids/88334/versions/1.23",
			responseStatus: http.StatusNoContent,
		},
		"missing EdgeWorkerID": {
			params: DeleteEdgeWorkerVersionRequest{
				Version: "1.23",
			},
			withError: ErrStructValidation,
		},
		"missing Version": {
			params: DeleteEdgeWorkerVersionRequest{
				EdgeWorkerID: 88334,
			},
			withError: ErrStructValidation,
		},
		"404 Not Found": {
			params: DeleteEdgeWorkerVersionRequest{
				EdgeWorkerID: 88334,
				Version:      "1.23",
			},
			expectedPath:   "/edgeworkers/v1/ids/88334/versions/1.23",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-not-found",
    "title": "The given resource could not be found.",
    "detail": "Unable to delete the version provided",
    "instance": "/edgeworkers/error-instances/514139f4-1608-4afc-88ac-67da91696af3",
    "status": 404,
    "errorCode": "EW2002"
}`,
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-not-found",
				Title:     "The given resource could not be found.",
				Detail:    "Unable to delete the version provided",
				Instance:  "/edgeworkers/error-instances/514139f4-1608-4afc-88ac-67da91696af3",
				Status:    404,
				ErrorCode: "EW2002",
			},
		},
		"500 internal server error": {
			params: DeleteEdgeWorkerVersionRequest{
				EdgeWorkerID: 88334,
				Version:      "1.23",
			},
			expectedPath:   "/edgeworkers/v1/ids/88334/versions/1.23",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
 "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
 "title": "Server Error",
 "status": 500,
 "instance": "host_name/edgeworkers/v1/ids/88334/versions/1.23",
 "method": "DELETE",
 "serverIp": "10.0.0.1",
 "clientIp": "192.168.0.1",
 "requestId": "a73affa111",
 "requestTime": "2021-12-17T16:32:37Z"
}`,
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/88334/versions/1.23",
				Method:      "DELETE",
				ServerIP:    "10.0.0.1",
				ClientIP:    "192.168.0.1",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-17T16:32:37Z",
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
			err := client.DeleteEdgeWorkerVersion(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				if test.responseStatus != 0 {
					assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				}

				return
			}
			require.NoError(t, err)
		})
	}
}
