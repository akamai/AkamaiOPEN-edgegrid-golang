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
	// GetRuleFormatsResponse contains the response body of GET /rule-formats request
	GetRuleFormatsResponse struct {
		RuleFormats RuleFormatItems `json:"ruleFormats"`
	}

	// RuleFormatItems contains a list of rule formats
	RuleFormatItems struct {
		Items []string `json:"items"`
	}
)

var (
	// ErrGetRuleFormats represents error when fetching rule formats fails
	ErrGetRuleFormats = errors.New("fetching rule formats")
)

func (p *papi) GetRuleFormats(ctx context.Context) (*GetRuleFormatsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRuleFormats")

	var ruleFormats GetRuleFormatsResponse

	req, err := request.NewGet(ctx, "/papi/v1/rule-formats").Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetRuleFormats, err)
	}

	resp, err := p.Exec(req, &ruleFormats)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetRuleFormats, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetRuleFormats, p.Error(resp))
	}

	return &ruleFormats, nil
}
