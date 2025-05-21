package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
)

type (
	// EvalPenaltyBoxConditions interface supports retrieving or modifying the evaluation penalty box conditions for
	// a specified security policy in evaluation mode.
	EvalPenaltyBoxConditions interface {
		// GetEvalPenaltyBoxConditions returns the eval penalty box conditions for a security policy in evaluation mode.
		GetEvalPenaltyBoxConditions(ctx context.Context, params GetPenaltyBoxConditionsRequest) (*GetPenaltyBoxConditionsResponse, error)

		// UpdateEvalPenaltyBoxConditions modifies the eval penalty box conditions for a security policy.
		UpdateEvalPenaltyBoxConditions(ctx context.Context, params UpdatePenaltyBoxConditionsRequest) (*UpdatePenaltyBoxConditionsResponse, error)
	}
)

func (p *appsec) GetEvalPenaltyBoxConditions(ctx context.Context, params GetPenaltyBoxConditionsRequest) (*GetPenaltyBoxConditionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetEvalPenaltyBoxConditions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-penalty-box/conditions",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalPenaltyBoxConditions request: %w", err)
	}

	var result GetPenaltyBoxConditionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get eval penalty box conditions request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateEvalPenaltyBoxConditions(ctx context.Context, params UpdatePenaltyBoxConditionsRequest) (*UpdatePenaltyBoxConditionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateEvalPenaltyBoxConditions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-penalty-box/conditions",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateEvalPenaltyBoxConditions request: %w", err)
	}

	var result UpdatePenaltyBoxConditionsResponse
	resp, err := p.Exec(req, &result, params.ConditionsPayload)
	if err != nil {
		return nil, fmt.Errorf("update eval penalty box conditions request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
