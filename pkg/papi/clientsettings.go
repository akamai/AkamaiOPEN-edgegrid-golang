package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

type (
	// ClientSettingsBody represents both the request and response bodies for operating on client settings resource
	ClientSettingsBody struct {
		RuleFormat  string `json:"ruleFormat"`
		UsePrefixes bool   `json:"usePrefixes"`
	}
)

var (
	// ErrGetClientSettings represents error when fetching client setting fails
	ErrGetClientSettings = errors.New("fetching client settings")
	// ErrUpdateClientSettings represents error when updating client setting fails
	ErrUpdateClientSettings = errors.New("updating client settings")
)

// GetClientSettings is used to list the client settings
func (p *papi) GetClientSettings(ctx context.Context) (*ClientSettingsBody, error) {
	logger := p.Log(ctx)
	logger.Debug("GetClientSettings")

	req, err := request.NewGet(ctx, "/papi/v1/client-settings").Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetClientSettings, err)
	}

	var clientSettings ClientSettingsBody
	resp, err := p.Exec(req, &clientSettings)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetClientSettings, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetClientSettings, p.Error(resp))
	}

	return &clientSettings, nil
}

// UpdateClientSettings is used to update the client settings
// fixme body structure
func (p *papi) UpdateClientSettings(ctx context.Context, params ClientSettingsBody) (*ClientSettingsBody, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateClientSettings")

	req, err := request.NewPut(ctx, "/papi/v1/client-settings").
		WithBody(params).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateClientSettings, err)
	}

	var clientSettings ClientSettingsBody
	resp, err := p.Exec(req, &clientSettings)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateClientSettings, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetClientSettings, p.Error(resp))
	}

	return &clientSettings, nil
}
