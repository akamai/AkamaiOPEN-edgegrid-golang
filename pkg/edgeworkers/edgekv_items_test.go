package edgeworkers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestListItems(t *testing.T) {
	tests := map[string]struct {
		params           ListItemsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListItemsResponse
		withError        error
	}{
		"200 OK - list edgekv items": {
			params: ListItemsRequest{
				ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
				[
					"US",
					"DE"
				]
			`,
			expectedPath: "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries",
			expectedResponse: &ListItemsResponse{
				"US",
				"DE",
			},
		},
		"validation - incorrect network": {
			params: ListItemsRequest{
				ItemsRequestParams{
					Network:     "stag",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing network": {
			params: ListItemsRequest{
				ItemsRequestParams{
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing namespace_id": {
			params: ListItemsRequest{
				ItemsRequestParams{
					Network: "staging",
					GroupID: "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing group_id": {
			params: ListItemsRequest{
				ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
				},
			},
			withError: ErrStructValidation,
		},
		"500 - internal server error": {
			params: ListItemsRequest{
				ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "https://learn.akamai.com",
				"title": "Internal Server Error",
				"detail": "An internal error occurred.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 500,
				"errorCode": "EKV_0000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Internal Server Error",
				Detail:    "An internal error occurred.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    500,
				ErrorCode: "EKV_0000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
			},
		},
		"404 - empty group": {
			params: ListItemsRequest{
				ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "https://learn.akamai.com",
				"title": "Not Found",
				"detail": "The requested group is empty or not found.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 404,
				"errorCode": "EKV_9000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Not Found",
				Detail:    "The requested group is empty or not found.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    404,
				ErrorCode: "EKV_9000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
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
			result, err := client.ListItems(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetItem(t *testing.T) {
	itemText := Item("English")
	itemPairs := Item(`
{
	"currency": "$",
	"flag": "/us.png",
	"name": "United States"
}
`)

	tests := map[string]struct {
		params           GetItemRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Item
		withError        error
	}{
		"200 OK - get edgekv item text": {
			params: GetItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus:   http.StatusOK,
			responseBody:     "English",
			expectedPath:     "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/GB",
			expectedResponse: &itemText,
		},
		"200 OK - get edgekv item key value pairs": {
			params: GetItemRequest{
				ItemID: "US",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"currency": "$",
	"flag": "/us.png",
	"name": "United States"
}
`,
			expectedPath:     "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/US",
			expectedResponse: &itemPairs,
		},
		"validation - incorrect network": {
			params: GetItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "stagng",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing network": {
			params: GetItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing namespace_id": {
			params: GetItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network: "staging",
					GroupID: "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing group_id": {
			params: GetItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing item_id": {
			params: GetItemRequest{
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"500 - internal server error": {
			params: GetItemRequest{
				ItemID: "US",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "https://learn.akamai.com",
				"title": "Internal Server Error",
				"detail": "An internal error occurred.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 500,
				"errorCode": "EKV_0000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/US",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Internal Server Error",
				Detail:    "An internal error occurred.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    500,
				ErrorCode: "EKV_0000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
			},
		},
		"404 - item doesn't exists": {
			params: GetItemRequest{
				ItemID: "US",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"type": "https://learn.akamai.com",
				"title": "Not Found",
				"detail": "Item 'GB' or group 'countries' was not found in the database.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 404,
				"errorCode": "EKV_9000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/US",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Not Found",
				Detail:    "Item 'GB' or group 'countries' was not found in the database.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    404,
				ErrorCode: "EKV_9000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
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
			result, err := client.GetItem(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpsertItem(t *testing.T) {
	itemText := Item("English")
	itemTextResponse := "Item was upserted in KV store with database 123456, namespace marketing, group countries, and key GB."
	itemPairs := Item(`
{
	"currency": "$",
	"flag": "/us.png",
	"name": "United States"
}
`)
	itemPairsResponse := "Item was upserted in KV store with database 123456, namespace marketing, group countries, and key US."

	tests := map[string]struct {
		params           UpsertItemRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *string
		withError        error
	}{
		"200 OK - upsert edgekv item text": {
			params: UpsertItemRequest{
				ItemID:   "GB",
				ItemData: itemText,
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus:   http.StatusOK,
			responseBody:     "Item was upserted in KV store with database 123456, namespace marketing, group countries, and key GB.",
			expectedPath:     "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/GB",
			expectedResponse: &itemTextResponse,
		},
		"200 OK - upsert edgekv item key value pairs": {
			params: UpsertItemRequest{
				ItemID:   "US",
				ItemData: itemPairs,
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus:   http.StatusOK,
			responseBody:     "Item was upserted in KV store with database 123456, namespace marketing, group countries, and key US.",
			expectedPath:     "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/US",
			expectedResponse: &itemPairsResponse,
		},
		"validation - incorrect network": {
			params: UpsertItemRequest{
				ItemID:   "US",
				ItemData: itemPairs,
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staing",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing network": {
			params: UpsertItemRequest{
				ItemID:   "US",
				ItemData: itemPairs,
				ItemsRequestParams: ItemsRequestParams{
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing namespace_id": {
			params: UpsertItemRequest{
				ItemID:   "US",
				ItemData: itemPairs,
				ItemsRequestParams: ItemsRequestParams{
					Network: "staging",
					GroupID: "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing group_id": {
			params: UpsertItemRequest{
				ItemID:   "US",
				ItemData: itemPairs,
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing item_id": {
			params: UpsertItemRequest{
				ItemData: itemPairs,
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing item_data": {
			params: UpsertItemRequest{
				ItemID: "US",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"500 - internal server error": {
			params: UpsertItemRequest{
				ItemID:   "US",
				ItemData: itemPairs,
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "https://learn.akamai.com",
				"title": "Internal Server Error",
				"detail": "An internal error occurred.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 500,
				"errorCode": "EKV_0000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/US",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Internal Server Error",
				Detail:    "An internal error occurred.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    500,
				ErrorCode: "EKV_0000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
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
			result, err := client.UpsertItem(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}

}

func TestDeleteItem(t *testing.T) {
	itemDeleteResponse := "Item was marked for deletion from database, namespace marketing, group countries, and key GB."

	tests := map[string]struct {
		params           DeleteItemRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *string
		withError        error
	}{
		"200 OK - delete edgekv item": {
			params: DeleteItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus:   http.StatusOK,
			responseBody:     "Item was marked for deletion from database, namespace marketing, group countries, and key GB.",
			expectedPath:     "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/GB",
			expectedResponse: &itemDeleteResponse,
		},
		"validation - incorrect network": {
			params: DeleteItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "sting",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing network": {
			params: DeleteItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing namespace_id": {
			params: DeleteItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network: "staging",
					GroupID: "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing group_id": {
			params: DeleteItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
				},
			},
			withError: ErrStructValidation,
		},
		"validation - missing item_id": {
			params: DeleteItemRequest{
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			withError: ErrStructValidation,
		},
		"500 - internal server error": {
			params: DeleteItemRequest{
				ItemID: "GB",
				ItemsRequestParams: ItemsRequestParams{
					Network:     "staging",
					NamespaceID: "marketing",
					GroupID:     "countries",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"type": "https://learn.akamai.com",
				"title": "Internal Server Error",
				"detail": "An internal error occurred.",
				"instance": "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				"status": 500,
				"errorCode": "EKV_0000",
				"additionalDetail": {
					"requestId": "db6e61d461c20395"
				}
			}`,
			expectedPath: "/edgekv/v1/networks/staging/namespaces/marketing/groups/countries/items/GB",
			withError: &Error{
				Type:      "https://learn.akamai.com",
				Title:     "Internal Server Error",
				Detail:    "An internal error occurred.",
				Instance:  "/edgeKV/error-instances/1386a423-377c-4dba-b746-abe84738f5c5",
				Status:    500,
				ErrorCode: "EKV_0000",
				AdditionalDetail: Additional{
					RequestID: "db6e61d461c20395",
				},
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
			result, err := client.DeleteItem(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}

}
