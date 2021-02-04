package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RuleUpgrade represents a collection of RuleUpgrade
//
// See: RuleUpgrade.GetRuleUpgrade()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// RuleUpgrade  contains operations available on RuleUpgrade  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getruleupgrade
	RuleUpgrade interface {
		GetRuleUpgrade(ctx context.Context, params GetRuleUpgradeRequest) (*GetRuleUpgradeResponse, error)
		UpdateRuleUpgrade(ctx context.Context, params UpdateRuleUpgradeRequest) (*UpdateRuleUpgradeResponse, error)
	}

	GetRuleUpgradeRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	GetRuleUpgradeResponse struct {
		Current          string `json:"current,omitempty"`
		Evaluating       string `json:"evaluating,omitempty"`
		Latest           string `json:"latest,omitempty"`
		KRSToEvalUpdates struct {
			UpdatedRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"updatedRules,omitempty"`
			NewRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"newRules,omitempty"`
		} `json:"KRSToEvalUpdates,omitempty"`
		EvalToEvalUpdates struct {
			NewRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"newRules,omitempty"`
		} `json:"EvalToEvalUpdates,omitempty"`
		KRSToLatestUpdates struct {
			DeletedRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"deletedRules,omitempty"`
			NewRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"newRules,omitempty"`
		} `json:"KRSToLatestUpdates,omitempty"`
	}

	UpdateRuleUpgradeRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Upgrade  bool   `json:"upgrade"`
	}

	UpdateRuleUpgradeResponse struct {
		Current string `json:"current"`
		Mode    string `json:"mode"`
		Eval    string `json:"eval"`
	}
)

// Validate validates GetRuleUpgradeRequest
func (v GetRuleUpgradeRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRuleUpgradeRequest
func (v UpdateRuleUpgradeRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetRuleUpgrade(ctx context.Context, params GetRuleUpgradeRequest) (*GetRuleUpgradeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRuleUpgrade")

	var rval GetRuleUpgradeResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/upgrade-details",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getruleupgrade request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getruleupgrade  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a RuleUpgrade.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleupgrade

func (p *appsec) UpdateRuleUpgrade(ctx context.Context, params UpdateRuleUpgradeRequest) (*UpdateRuleUpgradeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRuleUpgrade")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create RuleUpgraderequest: %w", err)
	}

	var rval UpdateRuleUpgradeResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create RuleUpgrade request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
