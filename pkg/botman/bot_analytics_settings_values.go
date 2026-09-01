package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

type (
	// The BotAnalyticsSettingsValues interface supports retrieving bot analytics settings values for an account.
	BotAnalyticsSettingsValues interface {
		// GetBotAnalyticsSettingsValues https://techdocs.akamai.com/bot-manager/reference/get-bot-analytics-settings-values
		GetBotAnalyticsSettingsValues(ctx context.Context) (map[string]any, error)
	}
)

func (b *botman) GetBotAnalyticsSettingsValues(ctx context.Context) (map[string]any, error) {
	logger := b.Log(ctx)
	logger.Debug("GetBotAnalyticsSettingsValues")

	req, err := request.NewGet(ctx, "/appsec/v1/bot-analytics-settings/values").Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotAnalyticsSettingsValues request: %w", err)
	}

	var result map[string]any
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotAnalyticsSettingsValues request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
