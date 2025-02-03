package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// The AkamaiBotCategory interface supports retrieving akamai bot categories
	AkamaiBotCategory interface {
		// GetAkamaiBotCategoryList https://techdocs.akamai.com/bot-manager/reference/get-akamai-bot-categories
		GetAkamaiBotCategoryList(ctx context.Context, params GetAkamaiBotCategoryListRequest) (*GetAkamaiBotCategoryListResponse, error)
	}

	// GetAkamaiBotCategoryListRequest is used to retrieve the akamai bot category actions for a policy.
	GetAkamaiBotCategoryListRequest struct {
		CategoryName string
	}

	// GetAkamaiBotCategoryListResponse is returned from a call to GetAkamaiBotCategoryList.
	GetAkamaiBotCategoryListResponse struct {
		Categories []map[string]interface{} `json:"categories"`
	}
)

func (b *botman) GetAkamaiBotCategoryList(ctx context.Context, params GetAkamaiBotCategoryListRequest) (*GetAkamaiBotCategoryListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetAkamaiBotCategoryList")

	uri := "/appsec/v1/akamai-bot-categories"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAkamaiBotCategoryList request: %w", err)
	}

	var result GetAkamaiBotCategoryListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAkamaiBotCategoryList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetAkamaiBotCategoryListResponse
	if params.CategoryName != "" {
		for _, val := range result.Categories {
			if val["categoryName"].(string) == params.CategoryName {
				filteredResult.Categories = append(filteredResult.Categories, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}
