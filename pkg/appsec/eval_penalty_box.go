package appsec

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// The EvalPenaltyBox interface supports retrieving or modifying the evaluation penalty box settings for
	// a specified security policy in evaluation mode.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#evalpenaltybox
	EvalPenaltyBox interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpenaltybox
		GetEvalPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpenaltybox
		UpdateEvalPenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error)
	}
)

func (p *appsec) GetEvalPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalPenaltyBox")

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
		return nil, fmt.Errorf("GetEvalPenaltyBox request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) UpdateEvalPenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error) {

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEvalPenaltyBox")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-penalty-box",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateEvalPenaltyBox request: %w", err)
	}

	var rval UpdatePenaltyBoxResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateEvalPenaltyBox request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
