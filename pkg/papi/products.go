package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetProductsRequest contains data required to list products associated to a contract
	GetProductsRequest struct {
		ContractID string
	}

	// GetProductsResponse contains details about all products associated to a contract
	GetProductsResponse struct {
		AccountID  string        `json:"accountId"`
		ContractID string        `json:"contractId"`
		Products   ProductsItems `json:"products"`
	}

	// ProductsItems contains a list of ProductItem
	ProductsItems struct {
		Items []ProductItem `json:"items"`
	}

	// ProductItem contains product resource data
	ProductItem struct {
		ProductName string `json:"productName"`
		ProductID   string `json:"productId"`
	}
)

// Validate validates GetProductsRequest
func (pr GetProductsRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(pr.ContractID, validation.Required),
	}.Filter()
}

var (
	// ErrGetProducts represents error when fetching products fails
	ErrGetProducts = errors.New("fetching products")
)

// GetProducts is used to list all products for a given contract
func (p *papi) GetProducts(ctx context.Context, params GetProductsRequest) (*GetProductsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetProducts, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetProducts")

	getURL := fmt.Sprintf("/papi/v1/products?contractId=%s", params.ContractID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetProducts, err)
	}

	var products GetProductsResponse
	resp, err := p.Exec(req, &products)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetProducts, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetProducts, p.Error(resp))
	}

	return &products, nil
}
