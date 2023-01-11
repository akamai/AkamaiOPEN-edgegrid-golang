package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomRuleAction interface supports retrieving and updating the actions for the custom
	// rules of a configuration, or for a specific custom rule.
	CustomRuleAction interface {
		// GetCustomRuleActions returns a list of all configured custom rules for the specified configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-custom-rules
		GetCustomRuleActions(ctx context.Context, params GetCustomRuleActionsRequest) (*GetCustomRuleActionsResponse, error)

		// GetCustomRuleAction returns a specified action of a custom rule.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-custom-rules
		GetCustomRuleAction(ctx context.Context, params GetCustomRuleActionRequest) (*GetCustomRuleActionResponse, error)

		// UpdateCustomRuleAction updates the action of a custom rule.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-custom-rule
		UpdateCustomRuleAction(ctx context.Context, params UpdateCustomRuleActionRequest) (*UpdateCustomRuleActionResponse, error)
	}

	// GetCustomRuleActionsRequest is used to retrieve the custom rule actions for a configuration.
	GetCustomRuleActionsRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
		RuleID   int    `json:"ruleId"`
	}

	// GetCustomRuleActionsResponse is returned from a call to GetCustomRuleActions.
	GetCustomRuleActionsResponse []struct {
		Action                string `json:"action,omitempty"`
		CanUseAdvancedActions bool   `json:"canUseAdvancedActions,omitempty"`
		Link                  string `json:"link,omitempty"`
		Name                  string `json:"name,omitempty"`
		RuleID                int    `json:"ruleId,omitempty"`
	}

	// GetCustomRuleActionRequest is used to retrieve the action for a custom rule.
	GetCustomRuleActionRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
		RuleID   int    `json:"ruleId"`
	}

	// GetCustomRuleActionResponse is returned from a call to GetCustomRuleAction.
	GetCustomRuleActionResponse struct {
		Action                string `json:"action,omitempty"`
		CanUseAdvancedActions bool   `json:"canUseAdvancedActions,omitempty"`
		Link                  string `json:"link,omitempty"`
		Name                  string `json:"name,omitempty"`
		RuleID                int    `json:"ruleId,omitempty"`
	}

	// UpdateCustomRuleActionRequest is used to modify an existing custom rule.
	UpdateCustomRuleActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
		Action   string `json:"action"`
	}

	// UpdateCustomRuleActionResponse is returned from a call to UpdateCustomRuleAction.
	UpdateCustomRuleActionResponse struct {
		Action                string `json:"action"`
		CanUseAdvancedActions bool   `json:"canUseAdvancedActions"`
		Link                  string `json:"link"`
		Name                  string `json:"name"`
		RuleID                int    `json:"ruleId"`
	}
)

// Validate validates a GetCustomRuleActionRequest.
func (v GetCustomRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetCustomRuleActionsRequest.
func (v GetCustomRuleActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomRuleActionRequest.
func (v UpdateCustomRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"ID":       validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetCustomRuleAction(ctx context.Context, params GetCustomRuleActionRequest) (*GetCustomRuleActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetCustomRuleAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetCustomRuleActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomRuleAction request: %w", err)
	}

	var results GetCustomRuleActionsResponse

	resp, err := p.Exec(req, &results)
	if err != nil {
		return nil, fmt.Errorf("get custom rule action request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	for _, val := range results {
		if val.RuleID == params.RuleID {
			result = val
			return &result, nil
		}
	}

	return &result, nil
}

func (p *appsec) GetCustomRuleActions(ctx context.Context, params GetCustomRuleActionsRequest) (*GetCustomRuleActionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetCustomRuleActions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomRuleActions request: %w", err)
	}

	var result GetCustomRuleActionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get custom rule actions request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RuleID != 0 {
		var filteredResult GetCustomRuleActionsResponse
		for _, val := range result {
			if val.RuleID == params.RuleID {
				filteredResult = append(filteredResult, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateCustomRuleAction(ctx context.Context, params UpdateCustomRuleActionRequest) (*UpdateCustomRuleActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateCustomRuleAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-rules/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomRuleAction request: %w", err)
	}

	var result UpdateCustomRuleActionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update custom rule action request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &result, nil
}
