package iam

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListCIDRBlocksRequest contains the request parameters for the ListCIDRBlocks endpoint.
	ListCIDRBlocksRequest struct {
		Actions bool
	}

	// ListCIDRBlocksResponse describes the response from the ListCIDRBlocks endpoint.
	ListCIDRBlocksResponse []CIDRBlock

	// CIDRBlock represents a CIDR block.
	CIDRBlock struct {
		Actions      *CIDRActions `json:"actions"`
		CIDRBlock    string       `json:"cidrBlock"`
		CIDRBlockID  int64        `json:"cidrBlockId"`
		Comments     *string      `json:"comments"`
		CreatedBy    string       `json:"createdBy"`
		CreatedDate  time.Time    `json:"createdDate"`
		Enabled      bool         `json:"enabled"`
		ModifiedBy   string       `json:"modifiedBy"`
		ModifiedDate time.Time    `json:"modifiedDate"`
	}

	// CIDRActions specifies activities available for the CIDR block.
	CIDRActions struct {
		Delete bool `json:"delete"`
		Edit   bool `json:"edit"`
	}

	// CreateCIDRBlockRequest contains the request parameters for the CreateCIDRBlock endpoint.
	CreateCIDRBlockRequest struct {
		CIDRBlock string  `json:"cidrBlock"`
		Comments  *string `json:"comments,omitempty"`
		Enabled   bool    `json:"enabled"`
	}

	// CreateCIDRBlockResponse describes the response from the CreateCIDRBlock endpoint.
	CreateCIDRBlockResponse CIDRBlock

	// GetCIDRBlockRequest contains the request parameters for the GetCIDRBlock endpoint.
	GetCIDRBlockRequest struct {
		CIDRBlockID int64
		Actions     bool
	}

	// GetCIDRBlockResponse describes the response from the GetCIDRBlock endpoint.
	GetCIDRBlockResponse CIDRBlock

	// UpdateCIDRBlockRequest contains the request parameters for the UpdateCIDRBlock endpoint.
	UpdateCIDRBlockRequest struct {
		CIDRBlockID int64
		Body        UpdateCIDRBlockRequestBody
	}

	// UpdateCIDRBlockRequestBody contains the request body to be used in the UpdateCIDRBlock endpoint.
	UpdateCIDRBlockRequestBody struct {
		CIDRBlock string  `json:"cidrBlock"`
		Comments  *string `json:"comments,omitempty"`
		Enabled   bool    `json:"enabled"`
	}

	// UpdateCIDRBlockResponse describes the response of the UpdateCIDRBlock endpoint.
	UpdateCIDRBlockResponse CIDRBlock

	// DeleteCIDRBlockRequest contains the request parameters for the DeleteCIDRBlock endpoint.
	DeleteCIDRBlockRequest struct {
		CIDRBlockID int64
	}

	// ValidateCIDRBlockRequest contains the request parameters for the ValidateCIDRBlock endpoint.
	ValidateCIDRBlockRequest struct {
		CIDRBlock string
	}
)

// validateCIDR validates the format of CIDRBlock.
func validateCIDR() validation.Rule {
	return validation.By(func(value interface{}) error {
		stringVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected type 'string', got: %T", value)
		}

		_, _, err := net.ParseCIDR(stringVal)
		if err != nil {
			return err
		}

		return nil
	})
}

// Validate validates validation on CreateCIDRBlockRequest.
func (r CreateCIDRBlockRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CIDRBlock": validation.Validate(r.CIDRBlock, validation.Required, validateCIDR()),
	})
}

// Validate validates validation on GetCIDRBlockRequest.
func (r GetCIDRBlockRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CIDRBlockID": validation.Validate(r.CIDRBlockID, validation.Required, validation.Min(1)),
	})
}

// Validate validates validation on UpdateCIDRBlockRequest.
func (r UpdateCIDRBlockRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CIDRBlockID": validation.Validate(r.CIDRBlockID, validation.Required, validation.Min(1)),
		"Body":        validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates validation on UpdateCIDRBlockRequestBody.
func (r UpdateCIDRBlockRequestBody) Validate() error {
	return validation.Errors{
		"CIDRBlock": validation.Validate(r.CIDRBlock, validation.Required, validateCIDR()),
	}.Filter()
}

// Validate validates validation on DeleteCIDRBlockRequest.
func (r DeleteCIDRBlockRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CIDRBlockID": validation.Validate(r.CIDRBlockID, validation.Required, validation.Min(1)),
	})
}

// Validate validates validation on ValidateCIDRBlockRequest.
func (r ValidateCIDRBlockRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CIDRBlock": validation.Validate(r.CIDRBlock, validation.Required, validateCIDR()),
	})
}

var (
	// ErrListCIDRBlocks is returned when ListCIDRBlocks fails.
	ErrListCIDRBlocks = errors.New("list CIDR blocks")
	// ErrCreateCIDRBlock is returned when CreateCIDRBlock fails.
	ErrCreateCIDRBlock = errors.New("create CIDR block")
	// ErrGetCIDRBlock is returned when GetCIDRBlock fails.
	ErrGetCIDRBlock = errors.New("get CIDR block")
	// ErrUpdateCIDRBlock is returned when UpdateCIDRBlock fails.
	ErrUpdateCIDRBlock = errors.New("update CIDR block")
	// ErrDeleteCIDRBlock is returned when DeleteCIDRBlock fails.
	ErrDeleteCIDRBlock = errors.New("delete CIDR block")
	// ErrValidateCIDRBlock is returned when ValidateCIDRBlock fails.
	ErrValidateCIDRBlock = errors.New("validate CIDR block")
)

func (i *iam) ListCIDRBlocks(ctx context.Context, params ListCIDRBlocksRequest) (ListCIDRBlocksResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListCIDRBlocks")

	uri, err := url.Parse("/identity-management/v3/user-admin/ip-acl/allowlist")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCIDRBlocks, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCIDRBlocks, err)
	}

	var result ListCIDRBlocksResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCIDRBlocks, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListCIDRBlocks, i.Error(resp))
	}

	return result, nil
}

func (i *iam) CreateCIDRBlock(ctx context.Context, params CreateCIDRBlockRequest) (*CreateCIDRBlockResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("CreateCIDRBlock")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateCIDRBlock, ErrStructValidation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/identity-management/v3/user-admin/ip-acl/allowlist", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateCIDRBlock, err)
	}

	var result CreateCIDRBlockResponse
	resp, err := i.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateCIDRBlock, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateCIDRBlock, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) GetCIDRBlock(ctx context.Context, params GetCIDRBlockRequest) (*GetCIDRBlockResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("GetCIDRBlock")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCIDRBlock, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/user-admin/ip-acl/allowlist/%d", params.CIDRBlockID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetCIDRBlock, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCIDRBlock, err)
	}

	var result GetCIDRBlockResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCIDRBlock, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCIDRBlock, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) UpdateCIDRBlock(ctx context.Context, params UpdateCIDRBlockRequest) (*UpdateCIDRBlockResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("UpdateCIDRBlock")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateCIDRBlock, ErrStructValidation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("/identity-management/v3/user-admin/ip-acl/allowlist/%d", params.CIDRBlockID), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateCIDRBlock, err)
	}

	var result UpdateCIDRBlockResponse
	resp, err := i.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateCIDRBlock, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateCIDRBlock, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) DeleteCIDRBlock(ctx context.Context, params DeleteCIDRBlockRequest) error {
	logger := i.Log(ctx)
	logger.Debug("DeleteCIDRBlock")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteCIDRBlock, ErrStructValidation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("/identity-management/v3/user-admin/ip-acl/allowlist/%d", params.CIDRBlockID), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteCIDRBlock, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteCIDRBlock, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteCIDRBlock, i.Error(resp))
	}

	return nil
}

func (i *iam) ValidateCIDRBlock(ctx context.Context, params ValidateCIDRBlockRequest) error {
	logger := i.Log(ctx)
	logger.Debug("ValidateCIDRBlock")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrValidateCIDRBlock, ErrStructValidation, err)
	}

	uri, err := url.Parse("/identity-management/v3/user-admin/ip-acl/allowlist/validate")
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrValidateCIDRBlock, err)
	}

	q := uri.Query()
	q.Add("cidrblock", params.CIDRBlock)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrValidateCIDRBlock, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrValidateCIDRBlock, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrValidateCIDRBlock, i.Error(resp))
	}

	return nil
}
