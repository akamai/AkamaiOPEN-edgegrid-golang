package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

type (
	// The EvalGroup interface supports creating, modifying and retrieving attack groups for evaluation.
	EvalGroup interface {
		// GetEvalGroups retrieves all attack groups currently under evaluation.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-eval-groups
		GetEvalGroups(ctx context.Context, params GetAttackGroupsRequest) (*GetAttackGroupsResponse, error)

		// GetEvalGroup retrieves a specific attack group currently under evaluation.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-eval-group
		GetEvalGroup(ctx context.Context, params GetAttackGroupRequest) (*GetAttackGroupResponse, error)

		// UpdateEvalGroup supports updating the condition and exception information for an attack group under evaluation.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-eval-group
		UpdateEvalGroup(ctx context.Context, params UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error)
	}
)

func (p *appsec) GetEvalGroup(ctx context.Context, params GetAttackGroupRequest) (*GetAttackGroupResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetEvalGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	var result GetAttackGroupResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get eval group request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetEvalGroups(ctx context.Context, params GetAttackGroupsRequest) (*GetAttackGroupsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetEvalGroups")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-groups?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalGroups request: %w", err)
	}

	var result GetAttackGroupsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get eval groups request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.Group != "" {
		var filteredResult GetAttackGroupsResponse
		for k, val := range result.AttackGroups {
			if val.Group == params.Group {
				filteredResult.AttackGroups = append(filteredResult.AttackGroups, result.AttackGroups[k])
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateEvalGroup(ctx context.Context, params UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateEvalGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-groups/%s/action-condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateEvalGroup request: %w", err)
	}

	var result UpdateAttackGroupResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update eval group request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
