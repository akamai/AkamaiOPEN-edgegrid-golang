package papi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

type (
	ClientSettings interface {
		GetClientSettings(context.Context) (*ClientSettingsBody, error)
		UpdateClientSettings(context.Context, ClientSettingsBody) (*ClientSettingsBody, error)
	}

	ClientSettingsBody struct {
		RuleFormat  string `json:"ruleFormat"`
		UsePrefixes bool   `json:"usePrefixes"`
	}
)

// GetClientSettings is used to list the client settings
func (p *papi) GetClientSettings(ctx context.Context) (*ClientSettingsBody, error) {
	logger := p.Log(ctx)
	logger.Debug("GetClientSettings")

	getURL := "/papi/v1/client-settings"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getclientsettings request: %w", err)
	}

	var clientSettings ClientSettingsBody
	resp, err := p.Exec(req, &clientSettings)
	if err != nil {
		return nil, fmt.Errorf("getclientsettings request failed: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &clientSettings, nil
}

// UpdateClientSettings is used to update the client settings
// fixme body structure
func (p *papi) UpdateClientSettings(ctx context.Context, params ClientSettingsBody) (*ClientSettingsBody, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateClientSettings")

	putURL := "/papi/v1/client-settings"
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create updateclientsettings request: %w", err)
	}

	var clientSettings ClientSettingsBody
	resp, err := p.Exec(req, &clientSettings, params)
	if err != nil {
		return nil, fmt.Errorf("updateclientsettings request failed: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &clientSettings, nil
}
