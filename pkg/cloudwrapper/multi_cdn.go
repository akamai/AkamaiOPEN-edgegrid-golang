package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// MultiCDN is the cloudwrapper Multi-CDN API interface
	MultiCDN interface {
		// ListAuthKeys lists the cdnAuthKeys for a specified contractId and cdnCode
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-auth-keys
		ListAuthKeys(context.Context, ListAuthKeysRequest) (*ListAuthKeysResponse, error)

		// ListCDNProviders lists CDN providers
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-providers
		ListCDNProviders(context.Context) (*ListCDNProvidersResponse, error)
	}

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
		return nil, fmt.Errorf("%s: %w: %s", ErrListAuthKeys, ErrStructValidation, err)
	}

	uri, err := url.Parse("/cloud-wrapper/v1/multi-cdn/auth-keys")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListAuthKeys, err)
	}
	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("cdnCode", params.CDNCode)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAuthKeys, err)
	}

	var result ListAuthKeysResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAuthKeys, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAuthKeys, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudwrapper) ListCDNProviders(ctx context.Context) (*ListCDNProvidersResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListCDNProviders")

	uri := "/cloud-wrapper/v1/multi-cdn/providers"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCDNProviders, err)
	}

	var result ListCDNProvidersResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCDNProviders, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListCDNProviders, c.Error(resp))
	}

	return &result, nil
}
