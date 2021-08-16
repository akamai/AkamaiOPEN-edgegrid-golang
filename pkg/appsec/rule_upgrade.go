package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RuleUpgrade interface supports verifying changes in Kona rule sets, and upgrading to the
	// latest rules.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#upgrade
	RuleUpgrade interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getupgradedetails
		GetRuleUpgrade(ctx context.Context, params GetRuleUpgradeRequest) (*GetRuleUpgradeResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putrules
		UpdateRuleUpgrade(ctx context.Context, params UpdateRuleUpgradeRequest) (*UpdateRuleUpgradeResponse, error)
	}

	// GetRuleUpgradeRequest is used to verify changes in the KRS rule sets.
	GetRuleUpgradeRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	// GetRuleUpgradeResponse is returned from a call to GetRuleUpgrade.
	GetRuleUpgradeResponse struct {
		Current          string `json:"current,omitempty"`
		Evaluating       string `json:"evaluating,omitempty"`
		Latest           string `json:"latest,omitempty"`
		KRSToEvalUpdates struct {
			DeletedRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"deletedRules,omitempty"`
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
			DeletedRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"deletedRules,omitempty"`
			UpdatedRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"updatedRules,omitempty"`
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
			UpdatedRules []struct {
				ID    int    `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"updatedRules,omitempty"`
		} `json:"KRSToLatestUpdates,omitempty"`
	}

	// UpdateRuleUpgradeRequest is used to upgrade to the most recent version of the KRS rule set.
	UpdateRuleUpgradeRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Upgrade  bool   `json:"upgrade"`
		Mode     string `json:"mode,omitempty"`
	}

	// UpdateRuleUpgradeResponse is returned from a call to UpdateRuleUpgrade.
	UpdateRuleUpgradeResponse struct {
		Current string `json:"current"`
		Mode    string `json:"mode"`
		Eval    string `json:"eval"`
	}
)

// Validate validates a GetRuleUpgradeRequest.
func (v GetRuleUpgradeRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateRuleUpgradeRequest.
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
		return nil, fmt.Errorf("failed to create GetRuleUpgrade request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetRuleUpgrade request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create UpdateRuleUpgrade request: %w", err)
	}

	var rval UpdateRuleUpgradeResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateRuleUpgrade request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
