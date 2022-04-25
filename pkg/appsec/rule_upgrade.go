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
		ConfigID int
		Version  int
		PolicyID string
	}

	// GetRuleUpgradeResponse is returned from a call to GetRuleUpgrade.
	GetRuleUpgradeResponse struct {
		Current            string             `json:"current,omitempty"`
		Evaluating         string             `json:"evaluating,omitempty"`
		Latest             string             `json:"latest,omitempty"`
		KRSToEvalUpdates   *RulesetUpdateData `json:"KRSToEvalUpdates,omitempty"`
		EvalToEvalUpdates  *RulesetUpdateData `json:"EvalToEvalUpdates,omitempty"`
		KRSToLatestUpdates *RulesetUpdateData `json:"KRSToLatestUpdates,omitempty"`
	}

	// RulesetUpdateData is used to report all updates to rules and attack groups in the ruleset.
	RulesetUpdateData struct {
		DeletedRules        *RuleData  `json:"deletedRules,omitempty"`
		NewRules            *RuleData  `json:"newRules,omitempty"`
		UpdatedRules        *RuleData  `json:"updatedRules,omitempty"`
		DeletedAttackGroups *GroupData `json:"deletedAttackGroups,omitempty"`
		UpdatedAttackGroups *GroupData `json:"updatedAttackGroups,omitempty"`
		NewAttackGroups     *GroupData `json:"newAttackGroups,omitempty"`
	}

	// RuleData contains updates to rules
	RuleData []struct {
		ID    int    `json:"id,omitempty"`
		Title string `json:"title,omitempty"`
	}

	// GroupData contains updates to attack groups
	GroupData []struct {
		Group     int    `json:"group,omitempty"`
		GroupName string `json:"groupName,omitempty"`
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
