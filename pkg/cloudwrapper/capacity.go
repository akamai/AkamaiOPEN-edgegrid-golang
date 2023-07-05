package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type (
	// Capacities is a Cloud Wrapper API interface.
	Capacities interface {
		// ListCapacities fetches capacities available for a given contractId.
		// If no contract id is provided, lists all available capacity locations
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/getcapacityinventory
		ListCapacities(context.Context, ListCapacitiesRequest) (*ListCapacitiesResponse, error)
	}

	// ListCapacitiesRequest is a request struct
	ListCapacitiesRequest struct {
		ContractIDs []string
	}

	// ListCapacitiesResponse contains response list of location capacities
	ListCapacitiesResponse struct {
		Capacities []LocationCapacity `json:"capacities"`
	}

	// LocationCapacity contains location capacity information
	LocationCapacity struct {
		LocationID         int          `json:"locationId"`
		LocationName       string       `json:"locationName"`
		ContractID         string       `json:"contractId"`
		Type               CapacityType `json:"type"`
		ApprovedCapacity   Capacity     `json:"approvedCapacity"`
		AssignedCapacity   Capacity     `json:"assignedCapacity"`
		UnassignedCapacity Capacity     `json:"unassignedCapacity"`
	}

	// CapacityType is a type of the properties, this capacity is related to
	CapacityType string

	// Capacity struct holds capacity information
	Capacity struct {
		Value int64 `json:"value"`
		Unit  Unit  `json:"unit"`
	}

	// Unit type of capacity value. Can be either GB or TB
	Unit string
)

const (
	// UnitGB is a const value for capacity unit
	UnitGB Unit = "GB"
	// UnitTB is a const value for capacity unit
	UnitTB Unit = "TB"
)

const (
	// CapacityTypeMedia type
	CapacityTypeMedia CapacityType = "MEDIA"
	// CapacityTypeWebStandardTLS type
	CapacityTypeWebStandardTLS CapacityType = "WEB_STANDARD_TLS"
	// CapacityTypeWebEnhancedTLS type
	CapacityTypeWebEnhancedTLS CapacityType = "WEB_ENHANCED_TLS"
)

var (
	// ErrListCapacities is returned in case an error occurs on ListCapacities operation
	ErrListCapacities = errors.New("list capacities")
)

func (c *cloudwrapper) ListCapacities(ctx context.Context, params ListCapacitiesRequest) (*ListCapacitiesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListCapacities")

	uri, err := url.Parse("/cloud-wrapper/v1/capacity")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCapacities, err)
	}

	q := uri.Query()
	for _, contractID := range params.ContractIDs {
		q.Add("contractIds", contractID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCapacities, err)
	}

	var result ListCapacitiesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCapacities, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListCapacities, c.Error(resp))
	}

	return &result, nil
}
