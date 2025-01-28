package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// The EvalPenaltyBox interface supports retrieving or modifying the evaluation penalty box settings for
	// a specified security policy in evaluation mode.
	EvalPenaltyBox interface {
		// GetEvalPenaltyBox returns the penalty box settings for a security policy in evaluation mode.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-eval-policy-penalty-box
		GetEvalPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error)

		// UpdateEvalPenaltyBox updates the penalty box settings for a security policy in evaluation mode.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-eval-policy-penalty-box
		UpdateEvalPenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error)
	}
)

func (p *appsec) GetEvalPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetEvalPenaltyBox")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-penalty-box",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalPenaltyBox request: %w", err)
	}

	var result GetPenaltyBoxResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get eval penalty box request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateEvalPenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateEvalPenaltyBox")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-penalty-box",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateEvalPenaltyBox request: %w", err)
	}

	var result UpdatePenaltyBoxResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update eval penalty box request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
