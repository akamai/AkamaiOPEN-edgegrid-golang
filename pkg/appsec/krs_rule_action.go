package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// KRSRuleAction represents a collection of KRSRuleAction
//
// See: KRSRuleAction.GetKRSRuleAction()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// KRSRuleAction  contains operations available on KRSRuleAction  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getkrsruleaction
	KRSRuleAction interface {
		GetKRSRuleActions(ctx context.Context, params GetKRSRuleActionsRequest) (*GetKRSRuleActionsResponse, error)
		GetKRSRuleAction(ctx context.Context, params GetKRSRuleActionRequest) (*GetKRSRuleActionResponse, error)
		UpdateKRSRuleAction(ctx context.Context, params UpdateKRSRuleActionRequest) (*UpdateKRSRuleActionResponse, error)
	}

	GetKRSRuleActionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetKRSRuleActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetKRSRuleActionsResponse struct {
		RuleActions []struct {
			Action string `json:"action"`
			ID     int    `json:"id"`
		} `json:"ruleActions"`
	}

	GetKRSRuleActionResponse struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}

	UpdateKRSRuleActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
		Action   string `json:"action"`
	}

	UpdateKRSRuleActionResponse struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}
)

// Validate validates GetKRSRuleActionRequest
func (v GetKRSRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates GetKRSRuleActionsRequest
func (v GetKRSRuleActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateKRSRuleActionRequest
func (v UpdateKRSRuleActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetKRSRuleAction(ctx context.Context, params GetKRSRuleActionRequest) (*GetKRSRuleActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetKRSRuleAction")

	var rval GetKRSRuleActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getkrsruleaction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getkrsruleaction  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetKRSRuleActions(ctx context.Context, params GetKRSRuleActionsRequest) (*GetKRSRuleActionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetKRSRuleActions")

	var rval GetKRSRuleActionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getkrsruleactions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getkrsruleactions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a KRSRuleAction.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putkrsruleaction

func (p *appsec) UpdateKRSRuleAction(ctx context.Context, params UpdateKRSRuleActionRequest) (*UpdateKRSRuleActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateKRSRuleAction")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create KRSRuleActionrequest: %w", err)
	}

	var rval UpdateKRSRuleActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create KRSRuleAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
