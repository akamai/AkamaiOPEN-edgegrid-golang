package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

type (
	// The BotAnalyticsCookieValues interface supports retrieving bot analytics cookie values for an account
	BotAnalyticsCookieValues interface {
		// GetBotAnalyticsCookieValues https://techdocs.akamai.com/bot-manager/reference/get-akamai-defined-bots
		GetBotAnalyticsCookieValues(ctx context.Context) (map[string]interface{}, error)
	}
)

func (b *botman) GetBotAnalyticsCookieValues(ctx context.Context) (map[string]interface{}, error) {

	logger := b.Log(ctx)
	logger.Debug("GetBotAnalyticsCookieValues")

	uri := "/appsec/v1/bot-analytics-cookie/values"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotAnalyticsCookieValues request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotAnalyticsCookieValues request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
