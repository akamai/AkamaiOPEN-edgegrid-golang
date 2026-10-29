package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
)

type (
	// ListLocationResponse represents a response object returned by ListLocations
	ListLocationResponse struct {
		Locations []Location `json:"locations"`
	}

	// Location represents a Location object
	Location struct {
		LocationID         int               `json:"locationId"`
		LocationName       string            `json:"locationName"`
		MultiCDNLocationID string            `json:"multiCdnLocationId"`
		TrafficTypes       []TrafficTypeItem `json:"trafficTypes"`
	}

	// TrafficTypeItem represents a TrafficType object for the location
	TrafficTypeItem struct {
		TrafficTypeID int    `json:"trafficTypeId"`
		TrafficType   string `json:"trafficType"`
		MapName       string `json:"mapName"`
	}
)

var (
	// ErrListLocations is returned when ListLocations fails
	ErrListLocations = errors.New("list locations")
)

func (c *cloudwrapper) ListLocations(ctx context.Context) (*ListLocationResponse, error) {
	req, err := request.NewGet(ctx, "/cloud-wrapper/v1/locations").Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request:\n%w", ErrListLocations, err)
	}

	var locations ListLocationResponse
	resp, err := c.Exec(req, &locations)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed:\n%w", ErrListLocations, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListLocations, c.Error(resp))
	}

	return &locations, nil
}
