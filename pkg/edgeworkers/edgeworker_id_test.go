package edgeworkers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEdgeWorkerID(t *testing.T) {
	tests := map[string]struct {
		params           GetEdgeWorkerIDRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *EdgeWorkerID
		withError        error
	}{
		"200 OK - get EdgeWorkerID": {
			params:         GetEdgeWorkerIDRequest{EdgeWorkerID: 12345},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "edgeWorkerId": 12345,
    "name": "EdgeWorkerID",
    "accountId": "B-123-WNKA6P",
    "groupId": 12345,
    "resourceTierId": 123,
    "createdBy": "jbond",
    "createdTime": "2021-04-19T07:08:37Z",
    "lastModifiedBy": "jbond",
    "lastModifiedTime": "2021-04-19T07:08:37Z"
}
`,
			expectedPath: "/edgeworkers/v1/ids/12345",
			expectedResponse: &EdgeWorkerID{
				EdgeWorkerID:     12345,
				Name:             "EdgeWorkerID",
				AccountID:        "B-123-WNKA6P",
				GroupID:          12345,
				ResourceTierID:   123,
				CreatedBy:        "jbond",
				CreatedTime:      "2021-04-19T07:08:37Z",
				LastModifiedBy:   "jbond",
				LastModifiedTime: "2021-04-19T07:08:37Z",
			},
		},
		"500 internal server error": {
			params:         GetEdgeWorkerIDRequest{EdgeWorkerID: 12345},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
  "title": "Server Error",
  "status": 500,
  "instance": "host_name/edgeworkers/v1/ids/12345",
  "method": "GET",
  "serverIp": "104.81.220.111",
  "clientIp": "89.64.55.111",
  "requestId": "a73affa111",
  "requestTime": "2021-12-06T10:27:11Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/12345",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/12345",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T10:27:11Z",
			},
		},
		"missing group EdgeWorkerID": {
			params:    GetEdgeWorkerIDRequest{},
			withError: ErrStructValidation,
		},
		"403 Forbidden - incorrect credentials": {
			params:         GetEdgeWorkerIDRequest{EdgeWorkerID: 12345},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/ids/12345",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2021-12-06T12:45:09Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/12345",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids/12345",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T12:45:09Z",
			},
		},
		"404 Not Found - EdgeWorkerID doesn't exist": {
			params:         GetEdgeWorkerIDRequest{EdgeWorkerID: 12345},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-not-found",
    "title": "The given resource could not be found.",
    "detail": "Unable to find the requested EdgeWorkerID ID",
    "instance": "/edgeworkers/error-instances/86d1cc10-4baf-49e1-b81a-075b72a2f6a4",
    "status": 404,
    "errorCode": "EW2002"
}`,
			expectedPath: "/edgeworkers/v1/ids/12345",
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-not-found",
				Title:     "The given resource could not be found.",
				Detail:    "Unable to find the requested EdgeWorkerID ID",
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
			result, err := client.GetEdgeWorkerID(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListEdgeWorkersID(t *testing.T) {
	tests := map[string]struct {
		params           ListEdgeWorkersIDRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListEdgeWorkersIDResponse
		withError        error
	}{
		"200 OK - list permission groups without parameters": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "edgeWorkerIds": [
        {
            "edgeWorkerId": 12345,
            "name": "edgeworker",
            "accountId": "B-3-WNK123",
            "groupId": 54321,
            "resourceTierId": 123,
            "createdBy": "jbond",
            "createdTime": "2020-05-05T21:55:36Z",
            "lastModifiedBy": "jbond",
            "lastModifiedTime": "2020-05-05T21:55:36Z"
        },
		{
            "edgeWorkerId": 12346,
            "name": "edgeworker-first",
            "accountId": "B-3-WNK123",
            "groupId": 54321,
            "resourceTierId": 234,
            "createdBy": "jbond",
            "createdTime": "2021-05-05T21:55:36Z",
            "lastModifiedBy": "jbond",
            "lastModifiedTime": "2021-05-05T21:55:36Z"
        }
    ]
}`,
			params:       ListEdgeWorkersIDRequest{},
			expectedPath: "/edgeworkers/v1/ids",
			expectedResponse: &ListEdgeWorkersIDResponse{[]EdgeWorkerID{
				{
					EdgeWorkerID:     12345,
					Name:             "edgeworker",
					AccountID:        "B-3-WNK123",
					GroupID:          54321,
					ResourceTierID:   123,
					CreatedBy:        "jbond",
					CreatedTime:      "2020-05-05T21:55:36Z",
					LastModifiedBy:   "jbond",
					LastModifiedTime: "2020-05-05T21:55:36Z",
				},
				{
					EdgeWorkerID:     12346,
					Name:             "edgeworker-first",
					AccountID:        "B-3-WNK123",
					GroupID:          54321,
					ResourceTierID:   234,
					CreatedBy:        "jbond",
					CreatedTime:      "2021-05-05T21:55:36Z",
					LastModifiedBy:   "jbond",
					LastModifiedTime: "2021-05-05T21:55:36Z",
				},
			}},
		},
		"200 OK - list permission groups with parameters": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "edgeWorkerIds": [
        {
            "edgeWorkerId": 12345,
            "name": "edgeworker",
            "accountId": "B-3-WNK123",
            "groupId": 54321,
            "resourceTierId": 234,
            "createdBy": "jbond",
            "createdTime": "2020-05-05T21:55:36Z",
            "lastModifiedBy": "jbond",
            "lastModifiedTime": "2020-05-05T21:55:36Z"
        },
		{
            "edgeWorkerId": 12346,
            "name": "edgeworker-first",
            "accountId": "B-3-WNK123",
            "groupId": 54321,
            "resourceTierId": 234,
            "createdBy": "jbond",
            "createdTime": "2021-05-05T21:55:36Z",
            "lastModifiedBy": "jbond",
            "lastModifiedTime": "2021-05-05T21:55:36Z"
        }
    ]
}`,
			params: ListEdgeWorkersIDRequest{
				GroupID:        54321,
				ResourceTierID: 234,
			},
			expectedPath: "/edgeworkers/v1/ids?groupId=54321&resourceTierId=234",
			expectedResponse: &ListEdgeWorkersIDResponse{[]EdgeWorkerID{
				{
					EdgeWorkerID:     12345,
					Name:             "edgeworker",
					AccountID:        "B-3-WNK123",
					GroupID:          54321,
					ResourceTierID:   234,
					CreatedBy:        "jbond",
					CreatedTime:      "2020-05-05T21:55:36Z",
					LastModifiedBy:   "jbond",
					LastModifiedTime: "2020-05-05T21:55:36Z",
				},
				{
					EdgeWorkerID:     12346,
					Name:             "edgeworker-first",
					AccountID:        "B-3-WNK123",
					GroupID:          54321,
					ResourceTierID:   234,
					CreatedBy:        "jbond",
					CreatedTime:      "2021-05-05T21:55:36Z",
					LastModifiedBy:   "jbond",
					LastModifiedTime: "2021-05-05T21:55:36Z",
				},
			}},
		},
		"200 OK - no EdgeWorkerID under resourceTierId": {
			params: ListEdgeWorkersIDRequest{
				GroupID:        54321,
				ResourceTierID: 123,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "edgeWorkerIds": []
}`,
			expectedPath:     "/edgeworkers/v1/ids?groupId=54321&resourceTierId=123",
			expectedResponse: &ListEdgeWorkersIDResponse{[]EdgeWorkerID{}},
		},
		"200 OK - no EdgeWorkerID under groupID": {
			params: ListEdgeWorkersIDRequest{
				GroupID:        46778,
				ResourceTierID: 123,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "edgeWorkerIds": []
}`,
			expectedPath:     "/edgeworkers/v1/ids?groupId=46778&resourceTierId=123",
			expectedResponse: &ListEdgeWorkersIDResponse{[]EdgeWorkerID{}},
		},
		"500 internal server error": {
			params:         ListEdgeWorkersIDRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
  "title": "Server Error",
  "status": 500,
  "instance": "host_name/edgeworkers/v1/ids",
  "method": "GET",
  "serverIp": "104.81.220.111",
  "clientIp": "89.64.55.111",
  "requestId": "a73affa111",
  "requestTime": "2021-12-06T10:27:11Z"
}`,
			expectedPath: "/edgeworkers/v1/ids",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-06T10:27:11Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params:         ListEdgeWorkersIDRequest{},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/ids",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2021-12-13T12:01:17Z"
}`,
			expectedPath: "/edgeworkers/v1/ids",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
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
			result, err := client.ListEdgeWorkersID(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreateEdgeWorkerID(t *testing.T) {
	tests := map[string]struct {
		params              CreateEdgeWorkerIDRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *EdgeWorkerID
		withError           error
	}{
		"201 Created - create EdgeWorkerID": {
			params: CreateEdgeWorkerIDRequest{
				GroupID:        12345,
				Name:           "New EdgeWorkerID",
				ResourceTierID: 123,
			},
			expectedRequestBody: `{"name":"New EdgeWorkerID","groupId":12345,"resourceTierId":123}`,
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "edgeWorkerId": 83969,
    "name": "New EdgeWorkerID",
    "accountId": "B-123-WNKA6P",
    "groupId": 12345,
    "resourceTierId": 123,
    "createdBy": "jbond",
    "createdTime": "2021-12-13T13:32:37Z",
    "lastModifiedBy": "jbond",
    "lastModifiedTime": "2021-12-13T13:32:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids",
			expectedResponse: &EdgeWorkerID{
				EdgeWorkerID:     83969,
				Name:             "New EdgeWorkerID",
				AccountID:        "B-123-WNKA6P",
				GroupID:          12345,
				ResourceTierID:   123,
				CreatedBy:        "jbond",
				CreatedTime:      "2021-12-13T13:32:37Z",
				LastModifiedBy:   "jbond",
				LastModifiedTime: "2021-12-13T13:32:37Z",
			},
		},
		"validation error": {
			params:    CreateEdgeWorkerIDRequest{},
			withError: ErrStructValidation,
		},
		"500 internal server error": {
			params: CreateEdgeWorkerIDRequest{
				GroupID:        12345,
				Name:           "New EdgeWorkerID",
				ResourceTierID: 123,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
  "title": "Server Error",
  "status": 500,
  "instance": "host_name/edgeworkers/v1/ids",
  "method": "POST",
  "serverIp": "104.81.220.111",
  "clientIp": "89.64.55.111",
  "requestId": "a73affa111",
  "requestTime": "2021-12-13T13:32:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids",
				Method:      "POST",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-13T13:32:37Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params: CreateEdgeWorkerIDRequest{
				GroupID:        12345,
				Name:           "New EdgeWorkerID",
				ResourceTierID: 123,
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/ids",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "POST",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2021-12-13T13:32:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "POST",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateEdgeWorkerID(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateEdgeWorkerID(t *testing.T) {
	tests := map[string]struct {
		params              UpdateEdgeWorkerIDRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *EdgeWorkerID
		withError           error
	}{
		"200 OK - update EdgeWorkerID": {
			params: UpdateEdgeWorkerIDRequest{
				Body: EdgeWorkerIDRequestBody{
					GroupID:        12345,
					Name:           "Update EdgeWorkerID",
					ResourceTierID: 123,
				},
				EdgeWorkerID: 54321,
			},
			expectedRequestBody: `{"name":"Update EdgeWorkerID","groupId":12345,"resourceTierId":123}`,
			responseStatus:      http.StatusOK,
			responseBody: `
{
    "edgeWorkerId": 54321,
    "name": "Update EdgeWorkerID",
    "accountId": "B-123-WNKA6P",
    "groupId": 12345,
    "resourceTierId": 123,
    "createdBy": "jbond",
    "createdTime": "2021-12-14T09:51:54Z",
    "lastModifiedBy": "jbond",
    "lastModifiedTime": "2021-12-14T09:51:54Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/54321",
			expectedResponse: &EdgeWorkerID{
				EdgeWorkerID:     54321,
				Name:             "Update EdgeWorkerID",
				AccountID:        "B-123-WNKA6P",
				GroupID:          12345,
				ResourceTierID:   123,
				CreatedBy:        "jbond",
				CreatedTime:      "2021-12-14T09:51:54Z",
				LastModifiedBy:   "jbond",
				LastModifiedTime: "2021-12-14T09:51:54Z",
			},
		},
		"validation error - empty body parameters": {
			params: UpdateEdgeWorkerIDRequest{
				Body:         EdgeWorkerIDRequestBody{},
				EdgeWorkerID: 54321,
			},
			withError: ErrStructValidation,
		},
		"validation error - empty edgeworker id": {
			params: UpdateEdgeWorkerIDRequest{
				Body: EdgeWorkerIDRequestBody{
					GroupID:        12345,
					Name:           "Update EdgeWorkerID",
					ResourceTierID: 123,
				},
			},
			withError: ErrStructValidation,
		},
		"500 internal server error": {
			params: UpdateEdgeWorkerIDRequest{
				Body: EdgeWorkerIDRequestBody{
					GroupID:        12345,
					Name:           "Update EdgeWorkerID",
					ResourceTierID: 123,
				},
				EdgeWorkerID: 54321,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
 "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
 "title": "Server Error",
 "status": 500,
 "instance": "host_name/edgeworkers/v1/ids/54321",
 "method": "PUT",
 "serverIp": "104.81.220.111",
 "clientIp": "89.64.55.111",
 "requestId": "a73affa111",
 "requestTime": "2021-12-13T13:32:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/54321",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/54321",
				Method:      "PUT",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-13T13:32:37Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params: UpdateEdgeWorkerIDRequest{
				Body: EdgeWorkerIDRequestBody{
					GroupID:        12345,
					Name:           "Update EdgeWorkerID",
					ResourceTierID: 123,
				},
				EdgeWorkerID: 54321,
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/ids/54321",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "PUT",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2021-12-14T10:26:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/54321",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids/54321",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "PUT",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-14T10:26:37Z",
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateEdgeWorkerID(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCloneEdgeWorkerID(t *testing.T) {
	tests := map[string]struct {
		params              CloneEdgeWorkerIDRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *EdgeWorkerID
		withError           error
	}{
		"200 OK - clone EdgeWorkerID with different resourceTierId": {
			params: CloneEdgeWorkerIDRequest{
				Body: EdgeWorkerIDRequestBody{
					GroupID:        12345,
					Name:           "Clone EdgeWorkerID",
					ResourceTierID: 123,
				},
				EdgeWorkerID: 54321,
			},
			expectedRequestBody: `{"name":"Clone EdgeWorkerID","groupId":12345,"resourceTierId":123}`,
			responseStatus:      http.StatusOK,
			responseBody: `
{
    "edgeWorkerId": 54322,
    "name": "Clone EdgeWorkerID",
    "accountId": "B-123-WNKA6P",
    "groupId": 12345,
    "resourceTierId": 123,
    "sourceEdgeWorkerId": 54321,
    "createdBy": "jbond",
    "createdTime": "2021-12-14T11:23:02Z",
    "lastModifiedBy": "jbond",
    "lastModifiedTime": "2021-12-14T11:23:02Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/54321/clone",
			expectedResponse: &EdgeWorkerID{
				EdgeWorkerID:       54322,
				Name:               "Clone EdgeWorkerID",
				AccountID:          "B-123-WNKA6P",
				GroupID:            12345,
				ResourceTierID:     123,
				SourceEdgeWorkerID: 54321,
				CreatedBy:          "jbond",
				CreatedTime:        "2021-12-14T11:23:02Z",
				LastModifiedBy:     "jbond",
				LastModifiedTime:   "2021-12-14T11:23:02Z",
			},
		},
		"validation error - empty body parameters": {
			params:    CloneEdgeWorkerIDRequest{EdgeWorkerID: 54321},
			withError: ErrStructValidation,
		},
		"validation error - empty edgeworker id": {
			params: CloneEdgeWorkerIDRequest{
				Body: EdgeWorkerIDRequestBody{
					GroupID:        12345,
					Name:           "Update EdgeWorkerID",
					ResourceTierID: 123,
				},
			},
			withError: ErrStructValidation,
		},
		"500 internal server error": {
			params: CloneEdgeWorkerIDRequest{
				Body: EdgeWorkerIDRequestBody{
					GroupID:        12345,
					Name:           "Clone EdgeWorkerID",
					ResourceTierID: 123,
				},
				EdgeWorkerID: 54321,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
"type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
"title": "Server Error",
"status": 500,
"instance": "host_name/edgeworkers/v1/ids/54321/clone",
"method": "POST",
"serverIp": "104.81.220.111",
"clientIp": "89.64.55.111",
"requestId": "a73affa111",
"requestTime": "2021-12-13T13:32:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/54321/clone",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/54321/clone",
				Method:      "POST",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-13T13:32:37Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params: CloneEdgeWorkerIDRequest{
				Body: EdgeWorkerIDRequestBody{
					GroupID:        12345,
					Name:           "Update EdgeWorkerID",
					ResourceTierID: 123,
				},
				EdgeWorkerID: 54321,
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
   "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
   "title": "Forbidden",
   "status": 403,
   "detail": "The client does not have the grant needed for the request",
   "instance": "host_name/edgeworkers/v1/ids/54321/clone",
   "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
   "method": "PUT",
   "serverIp": "104.81.220.111",
   "clientIp": "89.64.55.111",
   "requestId": "a73affa111",
   "requestTime": "2021-12-14T10:26:37Z"
}`,
			expectedPath: "/edgeworkers/v1/ids/54321/clone",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/ids/54321/clone",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "PUT",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-14T10:26:37Z",
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CloneEdgeWorkerID(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteEdgeWorkerIDRequest_Validate(t *testing.T) {
	tests := map[string]struct {
		request   DeleteEdgeWorkerIDRequest
		withError bool
	}{
		"empty id": {
			request:   DeleteEdgeWorkerIDRequest{},
			withError: true,
		},
		"ok": {
			request: DeleteEdgeWorkerIDRequest{EdgeWorkerID: 1},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.request.Validate()
			if test.withError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestEdgeworkers_DeleteEdgeWorkerID(t *testing.T) {
	tests := map[string]struct {
		params         DeleteEdgeWorkerIDRequest
		withError      error
		expectedPath   string
		responseStatus int
		responseBody   string
	}{
		"204 Deleted": {
			params:         DeleteEdgeWorkerIDRequest{EdgeWorkerID: 1},
			expectedPath:   "/edgeworkers/v1/ids/1",
			responseStatus: http.StatusNoContent,
		},
		"404 Not Found": {
			params:         DeleteEdgeWorkerIDRequest{EdgeWorkerID: 1},
			expectedPath:   "/edgeworkers/v1/ids/1",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/edgeworkers/error-types/edgeworkers-not-found",
    "title": "The given resource could not be found.",
    "detail": "Unable to delete this EdgeWorker ID",
    "instance": "/edgeworkers/error-instances/da246328-ed4a-4e5f-bed3-44e57f9ba7ef",
    "status": 404,
    "errorCode": "EW2002"
}`,
			withError: &Error{
				Type:      "/edgeworkers/error-types/edgeworkers-not-found",
				Title:     "The given resource could not be found.",
				Detail:    "Unable to delete this EdgeWorker ID",
				Instance:  "/edgeworkers/error-instances/da246328-ed4a-4e5f-bed3-44e57f9ba7ef",
				Status:    404,
				ErrorCode: "EW2002",
			},
		},
		"500 internal server error": {
			params:         DeleteEdgeWorkerIDRequest{EdgeWorkerID: 1},
			expectedPath:   "/edgeworkers/v1/ids/1",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
  "title": "Server Error",
  "status": 500,
  "instance": "host_name/edgeworkers/v1/ids/1",
  "method": "DELETE",
  "serverIp": "104.81.220.111",
  "clientIp": "89.64.55.111",
  "requestId": "a73affa111",
  "requestTime": "2021-12-17T16:32:37Z"
}`,
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/ids/1",
				Method:      "DELETE",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
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
			err := client.DeleteEdgeWorkerID(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				return
			}
			require.NoError(t, err)
		})
	}
}
