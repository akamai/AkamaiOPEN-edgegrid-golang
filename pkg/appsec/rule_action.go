package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RuleAction represents a collection of RuleAction
//
// See: RuleAction.GetRuleAction()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// RuleAction  contains operations available on RuleAction  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getruleaction
	RuleAction interface {
		GetRuleActions(ctx context.Context, params GetRuleActionsRequest) (*GetRuleActionsResponse, error)
		GetRuleAction(ctx context.Context, params GetRuleActionRequest) (*GetRuleActionResponse, error)
		UpdateRuleAction(ctx context.Context, params UpdateRuleActionRequest) (*UpdateRuleActionResponse, error)
	}

	GetRuleActionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetRuleActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetRuleActionsResponse struct {
		RuleActions []struct {
			Action string `json:"action,omitempty"`
			ID     int    `json:"id,omitempty"`
		} `json:"ruleActions,omitempty"`
	}

	GetRuleActionResponse struct {
		Action string `json:"action,omitempty"`
		ID     int    `json:"id,omitempty"`
	}

	UpdateRuleActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
		Action   string `json:"action"`
	}

	UpdateRuleActionResponse struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}
)

// Validate validates GetRuleActionRequest
func (v GetRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates GetRuleActionsRequest
func (v GetRuleActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRuleActionRequest
func (v UpdateRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetRuleAction(ctx context.Context, params GetRuleActionRequest) (*GetRuleActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRuleAction")

	var rval GetRuleActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getruleaction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getruleaction  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetRuleActions(ctx context.Context, params GetRuleActionsRequest) (*GetRuleActionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRuleActions")

	var rval GetRuleActionsResponse
	var rvalfiltered GetRuleActionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getruleactions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getruleactions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RuleID != 0 {
		for _, val := range rval.RuleActions {
			if val.ID == params.RuleID {
				rvalfiltered.RuleActions = append(rvalfiltered.RuleActions, val)
			}
		}

	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

// Update will update a RuleAction.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleaction

func (p *appsec) UpdateRuleAction(ctx context.Context, params UpdateRuleActionRequest) (*UpdateRuleActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRuleAction")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create RuleActionrequest: %w", err)
	}

	var rval UpdateRuleActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create RuleAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
