package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

	uri := "/cloudlets/v3/cloudlet-info"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCloudlets, err)
	}

	var result []ListCloudletsItem
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCloudlets, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListCloudlets, c.Error(resp))
	}

	return result, nil
}
