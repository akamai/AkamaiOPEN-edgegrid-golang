package appsec

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// The EvalGroup interface supports creating, modifying and retrieving attack groups for evaluation.
	EvalGroup interface {
		// GetEvalGroups retrieves all attack groups currently under evaluation.
		GetEvalGroups(ctx context.Context, params GetAttackGroupsRequest) (*GetAttackGroupsResponse, error)

		// GetEvalGroups retrieves a specific attack group currently under evaluation.
		GetEvalGroup(ctx context.Context, params GetAttackGroupRequest) (*GetAttackGroupResponse, error)

		// UpdateEvalGroup supports updating the condition and exception information for an attack group under evaluation.
		UpdateEvalGroup(ctx context.Context, params UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error)
	}
)

func (p *appsec) GetEvalGroup(ctx context.Context, params GetAttackGroupRequest) (*GetAttackGroupResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalGroup")

	var rval GetAttackGroupResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-groups/%s?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalGroup request: %w", err)
	}
	logger.Debugf("BEFORE GetEvalGroup %v", rval)
	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetEvalGroup request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}
	logger.Debugf("GetEvalGroup %v", rval)
	return &rval, nil

}

func (p *appsec) GetEvalGroups(ctx context.Context, params GetAttackGroupsRequest) (*GetAttackGroupsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalGroupConditionExceptions")

	var rval GetAttackGroupsResponse
	var rvalfiltered GetAttackGroupsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-groups?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalGroups request: %w", err)
	}
	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetEvalGroups request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.Group != "" {
		for k, val := range rval.AttackGroups {
			if val.Group == params.Group {
				rvalfiltered.AttackGroups = append(rvalfiltered.AttackGroups, rval.AttackGroups[k])
			}
		}
	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

func (p *appsec) UpdateEvalGroup(ctx context.Context, params UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEvalGroup")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-groups/%s/action-condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateEvalGroup request: %w", err)
	}

	var rval UpdateAttackGroupResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateEvalGroup request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
