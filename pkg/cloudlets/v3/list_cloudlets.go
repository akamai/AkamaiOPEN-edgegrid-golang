package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

type (
	// ListCloudletsItem contains the response data from ListCloudlets operation
	ListCloudletsItem struct {
		CloudletName string       `json:"cloudletName"`
		CloudletType CloudletType `json:"cloudletType"`
	}
)

var (
	// ErrListCloudlets is returned when ListCloudlets fails
	ErrListCloudlets = errors.New("list cloudlets")
)

func (c *cloudlets) ListCloudlets(ctx context.Context) ([]ListCloudletsItem, error) {
	logger := c.Log(ctx)
	logger.Debug("ListCloudlets")

	req, err := request.NewGet(ctx, "/cloudlets/v3/cloudlet-info").Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCloudlets, err)
	}

	var result []ListCloudletsItem
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCloudlets, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListCloudlets, c.Error(resp))
	}

	return result, nil
}
