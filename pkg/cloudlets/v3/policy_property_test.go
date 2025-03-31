package v3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ListActivePolicyProperties(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params           ListActivePolicyPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListActivePolicyPropertiesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK - no query params": {
			params: ListActivePolicyPropertiesRequest{
				PolicyID: 5,
			},
			responseStatus: http.StatusOK,
			responseBody: `
				{
  "page": {
    "number": 0,
    "size": 1000,
    "totalElements": 2,
    "totalPages": 1
  },
  "content": [
    {
      "groupId": 5,
      "id": 1234,
      "name": "property",
      "network": "PRODUCTION",
      "version": 1
    },
    {
      "groupId": 5,
      "id": 1233,
      "name": "property",
      "network": "STAGING",
      "version": 1
    }
  ],
  "links": [
    {
      "href": "/cloudlets/v3/policies/101/properties?page=0&size=1000",
      "rel": "self"
    }
  ]
}`,
			expectedPath: "/cloudlets/v3/policies/5/properties",
			expectedResponse: &ListActivePolicyPropertiesResponse{
				Page: Page{
					Number:        0,
					Size:          1000,
					TotalElements: 2,
					TotalPages:    1,
				},
				PolicyProperties: []ListPolicyPropertiesItem{
					{
						GroupID:       5,
						ID:            1234,
						Name:          "property",
						Network:       "PRODUCTION",
						PolicyVersion: 1,
					},
					{
						GroupID:       5,
						ID:            1233,
						Name:          "property",
						Network:       "STAGING",
						PolicyVersion: 1,
					},
				},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/101/properties?page=0&size=1000",
						Rel:  "self",
					},
				},
			},
		},
		"200 OK - with query params": {
			params: ListActivePolicyPropertiesRequest{
				PolicyID: 5,
				Page:     50,
				Size:     1000,
			},
			responseStatus: http.StatusOK,
			responseBody: `
				{
  "page": {
    "number": 50,
    "size": 1000,
    "totalElements": 2,
    "totalPages": 1
  },
  "content": [
    {
      "groupId": 5,
      "id": 1234,
      "name": "property",
      "network": "PRODUCTION",
      "version": 1
    },
    {
      "groupId": 5,
      "id": 1233,
      "name": "property",
      "network": "STAGING",
      "version": 1
    }
  ],
  "links": [
    {
      "href": "/cloudlets/v3/policies/101/properties?page=50&size=1000",
      "rel": "self"
    }
  ]
}`,
			expectedPath: "/cloudlets/v3/policies/5/properties?page=50&size=1000",
			expectedResponse: &ListActivePolicyPropertiesResponse{
				Page: Page{
					Number:        50,
					Size:          1000,
					TotalElements: 2,
					TotalPages:    1,
				},
				PolicyProperties: []ListPolicyPropertiesItem{
					{
						GroupID:       5,
						ID:            1234,
						Name:          "property",
						Network:       "PRODUCTION",
						PolicyVersion: 1,
					},
					{
						GroupID:       5,
						ID:            1233,
						Name:          "property",
						Network:       "STAGING",
						PolicyVersion: 1,
					},
				},
				Links: []Link{
					{
						Href: "/cloudlets/v3/policies/101/properties?page=50&size=1000",
						Rel:  "self",
					},
				},
			},
		},
		"200 OK - empty": {
			params: ListActivePolicyPropertiesRequest{
				PolicyID: 5,
				Page:     0,
				Size:     1000,
			},
			responseStatus: http.StatusOK,
			responseBody: `
				{
  "page": {
    "number": 0,
    "size": 1000,
    "totalElements": 2,
    "totalPages": 1
  },
  "content": [],
  "links": []
}`,
			expectedPath: "/cloudlets/v3/policies/5/properties?size=1000",
			expectedResponse: &ListActivePolicyPropertiesResponse{
				Page: Page{
					Number:        0,
					Size:          1000,
					TotalElements: 2,
					TotalPages:    1,
				},
				PolicyProperties: []ListPolicyPropertiesItem{},
				Links:            []Link{},
			},
		},
		"validation errors - missing required params": {
			params: ListActivePolicyPropertiesRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get policy properties: struct validation: PolicyID: cannot be blank", err.Error())
			},
		},
		"validation errors - size lower than 10, negative page number": {
			params: ListActivePolicyPropertiesRequest{
				PolicyID: 1,
				Page:     -2,
				Size:     5,
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "get policy properties: struct validation: Page: must be no less than 0\nSize: must be no less than 10", err.Error())
			},
		},
		"500 Internal Server Error": {
			params: ListActivePolicyPropertiesRequest{
				PolicyID: 1,
				Page:     0,
				Size:     1000,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
	"title": "Internal Server Error",
	"status": 500,
	"requestId": "1",
	"requestTime": "12:00",
	"clientIp": "1.1.1.1",
	"serverIp": "2.2.2.2",
	"method": "GET"
}`,
			expectedPath: "/cloudlets/v3/policies/1/properties?size=1000",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:        "internal_error",
					Title:       "Internal Server Error",
					Status:      http.StatusInternalServerError,
					RequestID:   "1",
					RequestTime: "12:00",
					ClientIP:    "1.1.1.1",
					ServerIP:    "2.2.2.2",
					Method:      "GET",
				}
				assert.ErrorIs(t, err, want)
			},
		},
		"404 Not found": {
			params: ListActivePolicyPropertiesRequest{
				PolicyID: 1,
				Page:     0,
				Size:     1000,
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "instance": "TestInstance",
    "status": 404,
    "title": "Not found",
    "type": "/cloudlets/v3/error-types/not-found",
    "errors": [
        {
            "detail": "Policy with id 1 not found.",
            "title": "Not found"
        }
    ]
}`,
			expectedPath: "/cloudlets/v3/policies/1/properties?size=1000",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:     "/cloudlets/v3/error-types/not-found",
					Title:    "Not found",
					Status:   http.StatusNotFound,
					Instance: "TestInstance",
					Errors: json.RawMessage(`
[
	{
		"detail": "Policy with id 1 not found.",
		"title": "Not found"
	}
]`)}
				assert.ErrorIs(t, err, want)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListActivePolicyProperties(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
