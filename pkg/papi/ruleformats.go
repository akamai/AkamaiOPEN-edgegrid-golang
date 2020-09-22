package papi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

type (
	RuleFormats interface {
		GetRuleFormats(context.Context) (*GetRuleFormatsResponse, error)
	}

	GetRuleFormatsResponse struct {
		RuleFormats RuleFormatItems `json:"ruleFormats"`
	}

	RuleFormatItems struct {
		Items []string `json:"items"`
	}
)

func (p *papi) GetRuleFormats(ctx context.Context) (*GetRuleFormatsResponse, error) {
	var ruleFormats GetRuleFormatsResponse

	logger := p.Log(ctx)
	logger.Debug("GetRuleFormats")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/papi/v1/rule-formats", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRuleFormats request: %w", err)
	}

	resp, err := p.Exec(req, &ruleFormats)
	if err != nil {
		return nil, fmt.Errorf("GetRuleFormats request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &ruleFormats, nil
}
