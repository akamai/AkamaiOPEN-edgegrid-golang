package clientlists

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// Lists interface to support creating, retrieving, updating and removing client lists.
	Lists interface {
		// GetClientLists lists all client lists accessible for an authenticated user
		//
		// See: https://techdocs.akamai.com/client-lists/reference/get-lists
		GetClientLists(ctx context.Context, params GetClientListsRequest) (*GetClientListsResponse, error)
	}

	// GetClientListsRequest contains request parameters for GetClientLists method
	GetClientListsRequest struct {
	}

	// GetClientListsResponse contains response parameters from GetClientLists method
	GetClientListsResponse struct {
		Content []ListContent
	}

	// ListContent contains list content
	ListContent struct {
		Name                       string   `json:"name"`
		Type                       string   `json:"type"`
		Notes                      string   `json:"notes"`
		Tags                       []string `json:"tags"`
		ListID                     string   `json:"listId"`
		Version                    int      `json:"version"`
		ItemsCount                 int      `json:"itemsCount"`
		CreateDate                 string   `json:"createDate"`
		CreatedBy                  string   `json:"createdBy"`
		UpdateDate                 string   `json:"updateDate"`
		UpdatedBy                  string   `json:"updatedBy"`
		ProductionActivationStatus string   `json:"productionActivationStatus"`
		StagingActivationStatus    string   `json:"stagingActivationStatus"`
		ListType                   string   `json:"listType"`
		Shared                     bool     `json:"shared"`
		ReadOnly                   bool     `json:"readOnly"`
		Deprecated                 bool     `json:"deprecated"`
	}
)

func (p *clientlists) GetClientLists(ctx context.Context, _ GetClientListsRequest) (*GetClientListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetClientLists")

	uri := "/client-list/v1/lists"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getClientLists request: %s", err.Error())
	}

	var rval GetClientListsResponse

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getClientLists request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
