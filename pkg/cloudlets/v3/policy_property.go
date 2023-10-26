package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListActivePolicyPropertiesRequest contains request parameters for ListActivePolicyProperties
	ListActivePolicyPropertiesRequest struct {
		PolicyID int64
		Page     int
		Size     int
	}

	// Network represents network on which policy version or property can be activated on
	Network string

	// PolicyProperty contains the response data from GetPolicyProperties operation
	PolicyProperty struct {
		Page    Page      `json:"page"`
		Content []Content `json:"content"`
		Links   []Link    `json:"links"`
	}

	// Content represents associated active properties information
	Content struct {
		GroupID int64   `json:"groupId"`
		ID      int64   `json:"id"`
		Name    string  `json:"name"`
		Network Network `json:"network"`
		Version int64   `json:"version"`
		Links   []Link  `json:"links"`
	}

	// Link represents hypermedia link to help navigate through the result set
	Link struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	}

	// Page contains informational data about pagination
	Page struct {
		Number        int `json:"number"`
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
	}
)

var (
	// ErrListActivePolicyProperties is returned when ListActivePolicyProperties fails
	ErrListActivePolicyProperties = errors.New("list active policy properties")
)

// Validate validates GetPolicyPropertiesRequest
func (r ListActivePolicyPropertiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
		"Page":     validation.Validate(r.Page, validation.Min(0)),
		"Size":     validation.Validate(r.Size, validation.Min(10)),
	})
}

func (c *cloudlets) ListActivePolicyProperties(ctx context.Context, params ListActivePolicyPropertiesRequest) (*PolicyProperty, error) {
	logger := c.Log(ctx)
	logger.Debug("ListActivePolicyProperties")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListActivePolicyProperties, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/v3/policies/%d/properties", params.PolicyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListActivePolicyProperties, err)
	}

	q := uri.Query()
	if params.Page != 0 {
		q.Add("page", strconv.Itoa(params.Page))
	}
	if params.Size != 0 {
		q.Add("size", strconv.Itoa(params.Size))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListActivePolicyProperties, err)
	}

	var result PolicyProperty
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListActivePolicyProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListActivePolicyProperties, c.Error(resp))
	}

	return &result, nil
}
