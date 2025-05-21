package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RuleUpgrade interface supports verifying changes in Kona rule sets, and upgrading to the
	// latest rules.
	RuleUpgrade interface {
		// GetRuleUpgrade only applies to Kona rule sets. The KRS rule sets are maintained by Akamai's security research team.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rules-upgrade-details-1
		GetRuleUpgrade(ctx context.Context, params GetRuleUpgradeRequest) (*GetRuleUpgradeResponse, error)

		// UpdateRuleUpgrade upgrades to the most recent version of the KRS rule set.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-policy-rules-1
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
	logger := p.Log(ctx)
	logger.Debug("GetRuleUpgrade")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	var result GetRuleUpgradeResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rule upgrade request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateRuleUpgrade(ctx context.Context, params UpdateRuleUpgradeRequest) (*UpdateRuleUpgradeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRuleUpgrade")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRuleUpgrade request: %w", err)
	}

	var result UpdateRuleUpgradeResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update rule upgrade request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
