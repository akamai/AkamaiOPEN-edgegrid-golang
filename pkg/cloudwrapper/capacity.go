package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

type (
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

	req, err := request.NewGet(ctx, "/cloud-wrapper/v1/capacity").
		AddQueryParams("contractIds", params.ContractIDs).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListCapacities, err)
	}

	var result ListCapacitiesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListCapacities, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListCapacities, c.Error(resp))
	}

	return &result, nil
}
