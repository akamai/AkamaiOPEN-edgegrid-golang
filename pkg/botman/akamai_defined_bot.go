package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
)

type (
	// The AkamaiDefinedBot interface supports retrieving akamai defined bots
	AkamaiDefinedBot interface {
		// GetAkamaiDefinedBotList https://techdocs.akamai.com/bot-manager/reference/get-akamai-defined-bots
		GetAkamaiDefinedBotList(ctx context.Context, params GetAkamaiDefinedBotListRequest) (*GetAkamaiDefinedBotListResponse, error)
	}

	// GetAkamaiDefinedBotListRequest is used to retrieve the akamai bot category actions for a policy.
	GetAkamaiDefinedBotListRequest struct {
		BotName string
	}

	// GetAkamaiDefinedBotListResponse is returned from a call to GetAkamaiDefinedBotList.
	GetAkamaiDefinedBotListResponse struct {
		Bots []map[string]interface{} `json:"bots"`
	}
)

func (b *botman) GetAkamaiDefinedBotList(ctx context.Context, params GetAkamaiDefinedBotListRequest) (*GetAkamaiDefinedBotListResponse, error) {

	logger := b.Log(ctx)
	logger.Debug("GetAkamaiDefinedBotList")

	uri := "/appsec/v1/akamai-defined-bots"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAkamaiDefinedBotList request: %w", err)
	}

	var result GetAkamaiDefinedBotListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAkamaiDefinedBotList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetAkamaiDefinedBotListResponse
	if params.BotName != "" {
		for _, val := range result.Bots {
			if val["botName"].(string) == params.BotName {
				filteredResult.Bots = append(filteredResult.Bots, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}
