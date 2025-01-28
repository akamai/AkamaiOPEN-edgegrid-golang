package cloudaccess

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// LookupPropertiesRequest holds parameters for LookupProperties
	LookupPropertiesRequest struct {
		AccessKeyUID int64
		Version      int64
	}

	// LookupPropertiesResponse contains response for LookupProperties
	LookupPropertiesResponse struct {
		Properties []Property `json:"properties"`
	}

	// Property holds information about property related to given access key
	Property struct {
		AccessKeyUID      int64  `json:"accessKeyUid"`
		Version           int64  `json:"version"`
		PropertyID        string `json:"propertyId"`
		PropertyName      string `json:"propertyName"`
		ProductionVersion *int64 `json:"productionVersion"`
		StagingVersion    *int64 `json:"stagingVersion"`
	}

	// GetAsyncPropertiesLookupIDRequest holds parameters for GetAsyncPropertiesLookupID
	GetAsyncPropertiesLookupIDRequest struct {
		AccessKeyUID int64
		Version      int64
	}

	// GetAsyncPropertiesLookupIDResponse contains response for GetAsyncPropertiesLookupID
	GetAsyncPropertiesLookupIDResponse struct {
		LookupID   int64 `json:"lookupId"`
		RetryAfter int64 `json:"retryAfter"`
	}

	// PerformAsyncPropertiesLookupRequest holds parameters for PerformAsyncPropertiesLookup
	PerformAsyncPropertiesLookupRequest struct {
		LookupID int64
	}

	// PerformAsyncPropertiesLookupResponse contains response for PerformAsyncPropertiesLookup
	PerformAsyncPropertiesLookupResponse struct {
		LookupID     int64        `json:"lookupId"`
		LookupStatus LookupStatus `json:"lookupStatus"`
		Properties   []Property   `json:"properties"`
	}

	// LookupStatus represents a lookup status
	LookupStatus string
)

const (
	// LookupComplete represents complete asynchronous property lookup status
	LookupComplete LookupStatus = "COMPLETE"
	// LookupError represents error asynchronous property lookup status
	LookupError LookupStatus = "ERROR"
	// LookupInProgress represents in progress asynchronous property lookup status
	LookupInProgress LookupStatus = "IN_PROGRESS"
	// LookupPending represents pending asynchronous property lookup status
	LookupPending LookupStatus = "PENDING"
	// LookupSubmitted represents submitted asynchronous property lookup status
	LookupSubmitted LookupStatus = "SUBMITTED"
)

// Validate validates LookupPropertiesRequest
func (r LookupPropertiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Version":      validation.Validate(r.Version, validation.Required),
		"AccessKeyUID": validation.Validate(r.AccessKeyUID, validation.Required),
	})
}

// Validate validates GetAsyncPropertiesLookupIDRequest
func (r GetAsyncPropertiesLookupIDRequest) Validate() interface{} {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Version":      validation.Validate(r.Version, validation.Required),
		"AccessKeyUID": validation.Validate(r.AccessKeyUID, validation.Required),
	})
}

// Validate validates PerformAsyncPropertiesLookupRequest
func (r PerformAsyncPropertiesLookupRequest) Validate() interface{} {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LookupID": validation.Validate(r.LookupID, validation.Required),
	})
}

var (
	// ErrLookupProperties is returned when LookupProperties fails
	ErrLookupProperties = errors.New("lookup properties")
	// ErrGetAsyncLookupIDProperties is returned when GetAsyncPropertiesLookupID fails
	ErrGetAsyncLookupIDProperties = errors.New("get lookup properties id async")
	// ErrPerformAsyncLookupProperties is returned when PerformAsyncPropertiesLookup fails
	ErrPerformAsyncLookupProperties = errors.New("perform async lookup properties")
)

func (c *cloudaccess) LookupProperties(ctx context.Context, params LookupPropertiesRequest) (*LookupPropertiesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("LookupProperties")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrLookupProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cam/v1/access-keys/%d/versions/%d/properties", params.AccessKeyUID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrLookupProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrLookupProperties, err)
	}

	var result LookupPropertiesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrLookupProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrLookupProperties, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudaccess) GetAsyncPropertiesLookupID(ctx context.Context, params GetAsyncPropertiesLookupIDRequest) (*GetAsyncPropertiesLookupIDResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetAsyncPropertiesLookupID")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAsyncLookupIDProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cam/v1/access-keys/%d/versions/%d/property-lookup-id", params.AccessKeyUID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetAsyncLookupIDProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetAsyncLookupIDProperties, err)
	}

	var result GetAsyncPropertiesLookupIDResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetAsyncLookupIDProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrGetAsyncLookupIDProperties, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudaccess) PerformAsyncPropertiesLookup(ctx context.Context, params PerformAsyncPropertiesLookupRequest) (*PerformAsyncPropertiesLookupResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("PerformAsyncPropertiesLookup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrPerformAsyncLookupProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cam/v1/property-lookups/%d", params.LookupID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrPerformAsyncLookupProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrPerformAsyncLookupProperties, err)
	}

	var result PerformAsyncPropertiesLookupResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrPerformAsyncLookupProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrPerformAsyncLookupProperties, c.Error(resp))
	}

	return &result, nil
}
