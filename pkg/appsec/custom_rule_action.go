package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CustomRuleAction represents a collection of CustomRuleAction
//
// See: CustomRuleAction.GetCustomRuleAction()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// CustomRuleAction  contains operations available on CustomRuleAction  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomruleaction
	CustomRuleAction interface {
		GetCustomRuleActions(ctx context.Context, params GetCustomRuleActionsRequest) (*GetCustomRuleActionsResponse, error)
		GetCustomRuleAction(ctx context.Context, params GetCustomRuleActionRequest) (*GetCustomRuleActionResponse, error)
		UpdateCustomRuleAction(ctx context.Context, params UpdateCustomRuleActionRequest) (*UpdateCustomRuleActionResponse, error)
	}

	GetCustomRuleActionsRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
		RuleID   int    `json:"ruleId"`
	}

	GetCustomRuleActionsResponse []struct {
		Action                string `json:"action"`
		CanUseAdvancedActions bool   `json:"canUseAdvancedActions"`
		Link                  string `json:"link"`
		Name                  string `json:"name"`
		RuleID                int    `json:"ruleId"`
	}

	GetCustomRuleActionRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		PolicyID string `json:"policyId"`
		RuleID   int    `json:"ruleId"`
	}

	GetCustomRuleActionResponse struct {
		Action                string `json:"action"`
		CanUseAdvancedActions bool   `json:"canUseAdvancedActions"`
		Link                  string `json:"link"`
		Name                  string `json:"name"`
		RuleID                int    `json:"ruleId"`
	}

	UpdateCustomRuleActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
		Action   string `json:"action"`
	}

	UpdateCustomRuleActionResponse struct {
		Action                string `json:"action"`
		CanUseAdvancedActions bool   `json:"canUseAdvancedActions"`
		Link                  string `json:"link"`
		Name                  string `json:"name"`
		RuleID                int    `json:"ruleId"`
	}
)

// Validate validates GetCustomRuleActionRequest
func (v GetCustomRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetCustomRuleActionsRequest
func (v GetCustomRuleActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateCustomRuleActionRequest
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
		return nil, fmt.Errorf("failed to create getcustomruleaction request: %w", err)
	}

	var rvals GetCustomRuleActionsResponse

	resp, err := p.Exec(req, &rvals)
	if err != nil {
		return nil, fmt.Errorf("getcustomruleaction request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create getcustomruleactions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getcustomruleactions request failed: %w", err)
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

// Update will update a CustomRuleAction.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putcustomruleaction

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
		return nil, fmt.Errorf("failed to create create CustomRuleActionrequest: %w", err)
	}

	var rval UpdateCustomRuleActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create CustomRuleAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
