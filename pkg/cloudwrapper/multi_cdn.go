package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListAuthKeysRequest is a request struct
	ListAuthKeysRequest struct {
		ContractID string
		CDNCode    string
	}

	// ListAuthKeysResponse contains response list of CDN auth keys
	ListAuthKeysResponse struct {
		CDNAuthKeys []MultiCDNAuthKey `json:"cdnAuthKeys"`
	}

	// MultiCDNAuthKey contains CDN auth key information
	MultiCDNAuthKey struct {
		AuthKeyName string `json:"authKeyName"`
		ExpiryDate  string `json:"expiryDate"`
		HeaderName  string `json:"headerName"`
	}

	// ListCDNProvidersResponse contains response list of CDN providers
	ListCDNProvidersResponse struct {
		CDNProviders []CDNProvider `json:"cdnProviders"`
	}

	// CDNProvider contains CDN provider information
	CDNProvider struct {
		CDNCode string `json:"cdnCode"`
		CDNName string `json:"cdnName"`
	}
)

// Validate validates ListAuthKeysRequest
func (r ListAuthKeysRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ContractID": validation.Validate(r.ContractID, validation.Required),
		"CDNCode":    validation.Validate(r.CDNCode, validation.Required),
	})
}

var (
	// ErrListAuthKeys is returned in case an error occurs on ListAuthKeys operation
	ErrListAuthKeys = errors.New("list auth keys")
	// ErrListCDNProviders is returned in case an error occurs on ListCDNProviders operation
	ErrListCDNProviders = errors.New("list CDN providers")
)

func (c *cloudwrapper) ListAuthKeys(ctx context.Context, params ListAuthKeysRequest) (*ListAuthKeysResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListAuthKeys")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListAuthKeys, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cloud-wrapper/v1/multi-cdn/auth-keys").
		AddQueryParam("contractId", params.ContractID).
		AddQueryParam("cdnCode", params.CDNCode).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListAuthKeys, err)
	}

	var result ListAuthKeysResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListAuthKeys, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListAuthKeys, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudwrapper) ListCDNProviders(ctx context.Context) (*ListCDNProvidersResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListCDNProviders")

	req, err := request.NewGet(ctx, "/cloud-wrapper/v1/multi-cdn/providers").Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListCDNProviders, err)
	}

	var result ListCDNProvidersResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListCDNProviders, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListCDNProviders, c.Error(resp))
	}

	return &result, nil
}
