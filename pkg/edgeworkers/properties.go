package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListPropertiesRequest contains parameters used to list properties
	ListPropertiesRequest struct {
		EdgeWorkerID int
		ActiveOnly   bool
	}

	// ListPropertiesResponse represents a response object returned by ListPropertiesRequest
	ListPropertiesResponse struct {
		Properties                []Property `json:"properties"`
		LimitedAccessToProperties bool       `json:"limitedAccessToProperties"`
	}

	// Property represents a single property object
	Property struct {
		ID                int64  `json:"propertyId"`
		Name              string `json:"propertyName"`
		StagingVersion    int    `json:"stagingVersion"`
		ProductionVersion int    `json:"productionVersion"`
		LatestVersion     int    `json:"latestVersion"`
	}
)

// Validate validates ListPropertiesRequest
func (r ListPropertiesRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(r.EdgeWorkerID, validation.Required),
	}.Filter()
}

var (
	// ErrListProperties is returned in case an error occurs on ListProperties operation
	ErrListProperties = errors.New("list properties")
)

func (e *edgeworkers) ListProperties(ctx context.Context, params ListPropertiesRequest) (*ListPropertiesResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListProperies")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/properties", params.EdgeWorkerID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListProperties, err)
	}

	q := uri.Query()
	q.Add("activeOnly", strconv.FormatBool(params.ActiveOnly))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListProperties, err)
	}

	var result ListPropertiesResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListProperties, e.Error(resp))
	}

	return &result, nil
}
