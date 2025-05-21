package appsec

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The TuningRecommendations interface supports retrieving tuning recommendations for a security policy, a specific attack group or a rule
	TuningRecommendations interface {
		// GetTuningRecommendations lists available tuning recommendations for a policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-recommendations
		GetTuningRecommendations(ctx context.Context, params GetTuningRecommendationsRequest) (*GetTuningRecommendationsResponse, error)

		// GetAttackGroupRecommendations returns available tuning recommendations for an attack group.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-attack-group
		GetAttackGroupRecommendations(ctx context.Context, params GetAttackGroupRecommendationsRequest) (*GetAttackGroupRecommendationsResponse, error)

		// GetRuleRecommendations returns available tuning recommendations for a rule.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-recommendations-rule
		GetRuleRecommendations(ctx context.Context, params GetRuleRecommendationsRequest) (*GetRuleRecommendationsResponse, error)
	}

	// GetTuningRecommendationsRequest is used to retrieve tuning recommendations for a security policy.
	GetTuningRecommendationsRequest struct {
		ConfigID    int
		Version     int
		PolicyID    string
		RulesetType RulesetType
	}

	// GetAttackGroupRecommendationsRequest is used to retrieve tuning recommendations for a specific attack group.
	GetAttackGroupRecommendationsRequest struct {
		ConfigID    int
		Version     int
		PolicyID    string
		Group       string
		RulesetType RulesetType
	}

	// GetRuleRecommendationsRequest is used to retrieve tuning recommendations for a specific rule.
	GetRuleRecommendationsRequest struct {
		ConfigID    int
		Version     int
		PolicyID    string
		RuleID      int
		RulesetType RulesetType
	}

	// GetTuningRecommendationsResponse is returned from a call to GetTuningRecommendations.
	GetTuningRecommendationsResponse struct {
		AttackGroupRecommendations []AttackGroupRecommendation `json:"attackGroupRecommendations,omitempty"`
		RuleRecommendations        []RuleRecommendation        `json:"ruleRecommendations,omitempty"`
		EvaluationPeriodStart      time.Time                   `json:"evaluationPeriodStart,omitempty"`
		EvaluationPeriodEnd        time.Time                   `json:"evaluationPeriodEnd,omitempty"`
	}

	// GetAttackGroupRecommendationsResponse is returned from a call to GetAttackGroupRecommendations.
	GetAttackGroupRecommendationsResponse AttackGroupRecommendation

	// GetRuleRecommendationsResponse is returned from a call to GetRuleRecommendations.
	GetRuleRecommendationsResponse RuleRecommendation

	// AttackGroupRecommendation is used to describe a recommendation for an attack group.
	AttackGroupRecommendation struct {
		Description string                `json:"description,omitempty"`
		Evidence    *Evidences            `json:"evidences,omitempty"`
		Exception   *AttackGroupException `json:"exception,omitempty"`
		Group       string                `json:"group,omitempty"`
	}
	// RuleRecommendation is used to describe a recommendation for a rule.
	RuleRecommendation struct {
		Description string                `json:"description,omitempty"`
		Evidence    *Evidences            `json:"evidences,omitempty"`
		Exception   *AttackGroupException `json:"exception,omitempty"`
		RuleId      int                   `json:"ruleId,omitempty"`
	}

	// Evidences is used to describe evidences for a recommendation.
	Evidences []struct {
		HostEvidences     []string `json:"hostEvidences,omitempty"`
		PathEvidences     []string `json:"pathEvidences,omitempty"`
		UserDataEvidences []string `json:"userDataEvidences,omitempty"`
	}
)

// Validate validates a GetTuningRecommendationsRequest.
func (v GetTuningRecommendationsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RulesetType": validation.Validate(v.RulesetType, validation.In(RulesetTypeActive, RulesetTypeEvaluation).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'active', 'evaluation' or '' (empty)", v.RulesetType))),
	}.Filter()
}

// Validate validates a GetAttackGroupRecommendationsRequest.
func (v GetAttackGroupRecommendationsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Group":    validation.Validate(v.Group, validation.Required),
		"RulesetType": validation.Validate(v.RulesetType, validation.In(RulesetTypeActive, RulesetTypeEvaluation).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'active', 'evaluation' or '' (empty)", v.RulesetType))),
	}.Filter()
}

// Validate validates a GetAttackGroupRecommendationsRequest.
func (v GetRuleRecommendationsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
		"RulesetType": validation.Validate(v.RulesetType, validation.In(RulesetTypeActive, RulesetTypeEvaluation).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'active', 'evaluation' or '' (empty)", v.RulesetType))),
	}.Filter()
}

func (p *appsec) GetTuningRecommendations(ctx context.Context, params GetTuningRecommendationsRequest) (*GetTuningRecommendationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetTuningRecommendations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/recommendations?standardException=true&type=%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RulesetType)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTuningRecommendations request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var result GetTuningRecommendationsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get tuning recommendations request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetAttackGroupRecommendations(ctx context.Context, params GetAttackGroupRecommendationsRequest) (*GetAttackGroupRecommendationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAttackGroupRecommendations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/recommendations/attack-groups/%s?standardException=true&type=%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
		params.RulesetType,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAttackGroupRecommendations request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var result GetAttackGroupRecommendationsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get attack group recommendations request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetRuleRecommendations(ctx context.Context, params GetRuleRecommendationsRequest) (*GetRuleRecommendationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRuleRecommendations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/recommendations/rules/%d?standardException=true&type=%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
		params.RulesetType,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRuleRecommendations request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var result GetRuleRecommendationsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rule recommendations request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
