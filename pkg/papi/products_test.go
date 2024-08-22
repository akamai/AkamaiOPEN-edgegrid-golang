package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiGetProducts(t *testing.T) {
	tests := map[string]struct {
		params           GetProductsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetProductsResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetProductsRequest{
				ContractID: "ctr_1-1TJZFW",
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "accountId": "act_1-1TJZFB",
    "contractId": "ctr_1-1TJZFW",
    "products": {
        "items": [
            {
                "productName": "Alta",
                "productId": "prd_Alta"
            }
        ]
    }
}`,
			expectedPath: "/papi/v1/products?contractId=ctr_1-1TJZFW",
			expectedResponse: &GetProductsResponse{
				AccountID:  "act_1-1TJZFB",
				ContractID: "ctr_1-1TJZFW",
				Products: ProductsItems{
					Items: []ProductItem{
						{
							ProductName: "Alta",
							ProductID:   "prd_Alta",
						},
					},
				},
			},
		},
		"500 internal server error": {
			params: GetProductsRequest{
				ContractID: "ctr_1-1TJZFW",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching products",
    "status": 500
}
`,
			expectedPath: "/papi/v1/products?contractId=ctr_1-1TJZFW",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error fetching products",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error empty contract ID": {
			params: GetProductsRequest{
				ContractID: "",
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "ContractID")
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
			result, err := client.GetProducts(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
