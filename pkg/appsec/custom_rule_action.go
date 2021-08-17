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
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#customruleactions
	CustomRuleAction interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomruleactions
		GetCustomRuleActions(ctx context.Context, params GetCustomRuleActionsRequest) (*GetCustomRuleActionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomruleactions
		GetCustomRuleAction(ctx context.Context, params GetCustomRuleActionRequest) (*GetCustomRuleActionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putactionruleid
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
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetCustomRuleAction")

	var rval GetCustomRuleActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomRuleAction request: %w", err)
	}

	var rvals GetCustomRuleActionsResponse

	resp, err := p.Exec(req, &rvals)
	if err != nil {
		return nil, fmt.Errorf("GetCustomRuleAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	for _, val := range rvals {
		if val.RuleID == params.RuleID {
			rval = val
			return &rval, nil
		}
	}

	return &rval, nil

}

func (p *appsec) GetCustomRuleActions(ctx context.Context, params GetCustomRuleActionsRequest) (*GetCustomRuleActionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetCustomRuleActions")

	var rval GetCustomRuleActionsResponse
	var rvalfiltered GetCustomRuleActionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomRuleActions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetCustomRuleActions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RuleID != 0 {
		for _, val := range rval {
			if val.RuleID == params.RuleID {
				rvalfiltered = append(rvalfiltered, val)
			}
		}
	} else {
		rvalfiltered = rval
	}
	return &rvalfiltered, nil

}

func (p *appsec) UpdateCustomRuleAction(ctx context.Context, params UpdateCustomRuleActionRequest) (*UpdateCustomRuleActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateCustomRuleAction")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/custom-rules/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomRuleAction request: %w", err)
	}

	var rval UpdateCustomRuleActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomRuleAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
