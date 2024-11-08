package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListGroupsWithinNamespaceRequest contains parameters used to get groups within a namespace
	ListGroupsWithinNamespaceRequest struct {
		Network     NamespaceNetwork
		NamespaceID string
	}
)

// Validate validates ListGroupsWithinNamespaceRequest
func (g ListGroupsWithinNamespaceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network":     validation.Validate(g.Network, validation.Required),
		"NamespaceID": validation.Validate(g.NamespaceID, validation.Required),
	})
}

// ErrListGroupsWithinNamespace is returned in case an error occurs on ListGroupsWithinNamespace operation
var ErrListGroupsWithinNamespace = errors.New("list groups within namespace")

func (e *edgeworkers) ListGroupsWithinNamespace(ctx context.Context, params ListGroupsWithinNamespaceRequest) ([]string, error) {
	logger := e.Log(ctx)
	logger.Debug("ListGroupsWithinNamespace")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListGroupsWithinNamespace, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgekv/v1/networks/%s/namespaces/%s/groups", params.Network, params.NamespaceID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListGroupsWithinNamespace, err)
	}

	var result []string
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListGroupsWithinNamespace, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListGroupsWithinNamespace, e.Error(resp))
	}

	return result, nil
}
