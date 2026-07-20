package cloudcertificates

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

func (c *cloudcertificates) ListLineageBindings(ctx context.Context, params ListLineageBindingsRequest) (*ListLineageBindingsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListLineageBindings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListLineageBindings, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v2/lineages/%d/bindings", params.LineageID).
		AddQueryParamIf("network", string(params.Network), params.Network != "").
		AddQueryParamIf("pageSize", strconv.Itoa(params.PageSize), params.PageSize > 0).
		AddQueryParamIf("after", params.After, params.After != "").
		AddQueryParamIf("sort", string(params.Sort), params.Sort != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListLineageBindings, err)
	}

	var result ListLineageBindingsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListLineageBindings, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListLineageBindings, c.Error(resp))
	}

	return &result, nil
}
