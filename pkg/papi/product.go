package papi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
)

type (
	Product interface {
		GetProducts(context.Context, GetProductsRequest) (*GetProductsResponse, error)
	}

	Products struct {
		Items []ProductItem `json:"items"`
	}

	ProductItem struct {
		ProductName string `json:"productName"`
		ProductID   string `json:"productId"`
	}

	GetProductsRequest struct {
		ContractID string
	}

	GetProductsResponse struct {
		AccountID  string     `json:"accountId"`
		ContractID string     `json:"contractId"`
		Products   []Products `json:"products"`
	}
)

func (pr GetProductsRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(pr.ContractID, validation.Required),
	}
}

func (p *papi) GetProducts(ctx context.Context, params GetProductsRequest) (*GetProductsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetProducts")

	getURL := fmt.Sprintf("/papi/v1/products?contractId=%s", params.ContractID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getproducts request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	var products GetProductsResponse
	resp, err := p.Exec(req, &products)
	if err != nil {
		return nil, fmt.Errorf("getproducts request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &products, nil
}
