package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListActivePolicyPropertiesRequest contains request parameters for ListActivePolicyProperties
	ListActivePolicyPropertiesRequest struct {
		PolicyID int64
		Page     int
		Size     int
	}

	// ListActivePolicyPropertiesResponse contains the response data from ListActivePolicyProperties operation.
	ListActivePolicyPropertiesResponse struct {
		Page             Page                       `json:"page"`
		PolicyProperties []ListPolicyPropertiesItem `json:"content"`
		Links            []Link                     `json:"links"`
	}

	// ListPolicyPropertiesItem represents associated active properties information.
	ListPolicyPropertiesItem struct {
		GroupID       int64   `json:"groupId"`
		ID            int64   `json:"id"`
		Name          string  `json:"name"`
		Network       Network `json:"network"`
		PolicyVersion int64   `json:"version"`
	}

	// Link represents hypermedia link to help navigate through the result set.
	Link struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	}

	// Page contains informational data about pagination.
	Page struct {
		Number        int `json:"number"`
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
	}
)

var (
	// ErrListActivePolicyProperties is returned when ListActivePolicyProperties fails.
	ErrListActivePolicyProperties = errors.New("list active policy properties")
)

// Validate validates ListActivePolicyPropertiesRequest.
func (r ListActivePolicyPropertiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
		"Page":     validation.Validate(r.Page, validation.Min(0)),
		"Size":     validation.Validate(r.Size, validation.Min(10)),
	})
}

func (c *cloudlets) ListActivePolicyProperties(ctx context.Context, params ListActivePolicyPropertiesRequest) (*ListActivePolicyPropertiesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListActivePolicyProperties")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListActivePolicyProperties, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cloudlets/v3/policies/%d/properties", params.PolicyID).
		AddQueryParamIf("page", strconv.Itoa(params.Page), params.Page != 0).
		AddQueryParamIf("size", strconv.Itoa(params.Size), params.Size != 0).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListActivePolicyProperties, err)
	}

	var result ListActivePolicyPropertiesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListActivePolicyProperties, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListActivePolicyProperties, c.Error(resp))
	}

	return &result, nil
}
