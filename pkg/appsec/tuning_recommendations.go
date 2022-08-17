package appsec

import (
	"context"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The TuningRecommendations interface supports retrieving tuning recommendations for a policy or a specific attack group.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#tunning_recommendations
	TuningRecommendations interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#gettunningrecommendations
		GetTuningRecommendations(ctx context.Context, params GetTuningRecommendationsRequest) (*GetTuningRecommendationsResponse, error)
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getattackgrouprecommendations
		GetAttackGroupRecommendations(ctx context.Context, params GetAttackGroupRecommendationsRequest) (*GetAttackGroupRecommendationsResponse, error)
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

	// GetTuningRecommendationsResponse is returned from a call to GetTuningRecommendations.
	GetTuningRecommendationsResponse struct {
		AttackGroupRecommendations []AttackGroupRecommendation `json:"attackGroupRecommendations,omitempty"`
		EvaluationPeriodStart      time.Time                   `json:"evaluationPeriodStart,omitempty"`
		EvaluationPeriodEnd        time.Time                   `json:"evaluationPeriodEnd,omitempty"`
	}

	// GetAttackGroupRecommendationsResponse is returned from a call to GetAttackGroupRecommendations.
	GetAttackGroupRecommendationsResponse AttackGroupRecommendation

	// AttackGroupRecommendation is used to describe a recommendation.
	AttackGroupRecommendation struct {
		Description string                `json:"description,omitempty"`
		Evidence    *Evidences            `json:"evidences,omitempty"`
		Exception   *AttackGroupException `json:"exception,omitempty"`
		Group       string                `json:"group,omitempty"`
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

func (p *appsec) GetTuningRecommendations(ctx context.Context, params GetTuningRecommendationsRequest) (*GetTuningRecommendationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetTuningRecommendations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetTuningRecommendationsResponse

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

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetTuningRecommendations request failed: %w", err)
	}

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

	var result GetAttackGroupRecommendationsResponse

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

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAttackGroupRecommendations request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}
