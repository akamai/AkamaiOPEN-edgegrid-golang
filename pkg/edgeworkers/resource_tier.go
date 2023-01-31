package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ResourceTiers is an edgeworkers resource tiers API interface
	ResourceTiers interface {
		// ListResourceTiers lists all resource tiers for a given contract
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-resource-tiers
		ListResourceTiers(context.Context, ListResourceTiersRequest) (*ListResourceTiersResponse, error)

		// GetResourceTier returns resource tier for a given edgeworker ID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-id-resource-tier
		GetResourceTier(context.Context, GetResourceTierRequest) (*ResourceTier, error)
	}

	// ListResourceTiersRequest contains parameters used to list resource tiers
	ListResourceTiersRequest struct {
		ContractID string
	}

	// GetResourceTierRequest contains parameters used to get a resource tier
	GetResourceTierRequest struct {
		EdgeWorkerID int
	}

	// ListResourceTiersResponse represents a response object returned by ListResourceTiers
	ListResourceTiersResponse struct {
		ResourceTiers []ResourceTier `json:"resourceTiers"`
	}

	// ResourceTier represents a single resource tier object
	ResourceTier struct {
		ID               int               `json:"resourceTierId"`
		Name             string            `json:"resourceTierName"`
		EdgeWorkerLimits []EdgeWorkerLimit `json:"edgeWorkerLimits"`
	}

	// EdgeWorkerLimit represents a single edgeworker limit object
	EdgeWorkerLimit struct {
		LimitName  string `json:"limitName"`
		LimitValue int64  `json:"limitValue"`
		LimitUnit  string `json:"limitUnit"`
	}
)

// Validate validates ListResourceTiersRequest
func (r ListResourceTiersRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(r.ContractID, validation.Required),
	}.Filter()
}

// Validate validates GetResourceTierRequest
func (r GetResourceTierRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(r.EdgeWorkerID, validation.Required),
	}.Filter()
}

var (
	// ErrListResourceTiers is returned in case an error occurs on ListResourceTiers operation
	ErrListResourceTiers = errors.New("list resource tiers")
	// ErrGetResourceTier is returned in case an error occurs on GetResourceTier operation
	ErrGetResourceTier = errors.New("get a resource tier")
)

func (e *edgeworkers) ListResourceTiers(ctx context.Context, params ListResourceTiersRequest) (*ListResourceTiersResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListResourceTiers")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListResourceTiers, ErrStructValidation, err)
	}

	uri, err := url.Parse("/edgeworkers/v1/resource-tiers")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListResourceTiers, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListResourceTiers, err)
	}

	var result ListResourceTiersResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListResourceTiers, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListResourceTiers, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) GetResourceTier(ctx context.Context, params GetResourceTierRequest) (*ResourceTier, error) {
	logger := e.Log(ctx)
	logger.Debug("GetResourceTier")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetResourceTier, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/resource-tier", params.EdgeWorkerID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetResourceTier, err)
	}

	var result ResourceTier
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetResourceTier, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetResourceTier, e.Error(resp))
	}

	return &result, nil
}
