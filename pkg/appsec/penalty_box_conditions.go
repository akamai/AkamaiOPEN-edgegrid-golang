package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PenaltyBoxConditions interface supports retrieving or modifying the penalty box conditions for
	// a specified security policy
	PenaltyBoxConditions interface {
		// GetPenaltyBoxConditions retrieves the penalty box conditions
		GetPenaltyBoxConditions(ctx context.Context, params GetPenaltyBoxConditionsRequest) (*GetPenaltyBoxConditionsResponse, error)

		// UpdatePenaltyBoxConditions modifies the penalty box conditions for a security policy.
		UpdatePenaltyBoxConditions(ctx context.Context, params UpdatePenaltyBoxConditionsRequest) (*UpdatePenaltyBoxConditionsResponse, error)
	}

	// GetPenaltyBoxConditionsRequest describes the GET request for penalty box conditions.
	GetPenaltyBoxConditionsRequest struct {
		ConfigID int
		Version  int
		PolicyID string
	}

	// UpdatePenaltyBoxConditionsRequest describes the PUT request to modify the penalty box conditions.
	UpdatePenaltyBoxConditionsRequest struct {
		ConfigID          int
		Version           int
		PolicyID          string
		ConditionsPayload PenaltyBoxConditionsPayload
	}

	// PenaltyBoxConditionsPayload describes the penalty box conditions with operator
	PenaltyBoxConditionsPayload struct {
		ConditionOperator string          `json:"conditionOperator"`
		Conditions        *RuleConditions `json:"conditions"`
	}

	// GetPenaltyBoxConditionsResponse describes the response with the penalty box conditions.
	GetPenaltyBoxConditionsResponse PenaltyBoxConditionsPayload

	// UpdatePenaltyBoxConditionsResponse describes the response with the update of penalty box conditions.
	UpdatePenaltyBoxConditionsResponse PenaltyBoxConditionsPayload
)

// Validate validates a GetPenaltyBoxConditionsRequest.
func (v GetPenaltyBoxConditionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	})
}

// Validate validates an UpdatePenaltyBoxConditionsRequest.
func (v UpdatePenaltyBoxConditionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":          validation.Validate(v.ConfigID, validation.Required),
		"Version":           validation.Validate(v.Version, validation.Required),
		"PolicyID":          validation.Validate(v.PolicyID, validation.Required),
		"ConditionsPayload": validation.Validate(v.ConditionsPayload, validation.Required),
	})
}

// Validate validates an ConditionsPayload
func (v PenaltyBoxConditionsPayload) Validate() error {
	return validation.Errors{
		"ConditionOperator": validation.Validate(v.ConditionOperator, validation.Required),
		"Conditions":        validation.Validate(v.Conditions, validation.NotNil),
	}.Filter()
}

func (p *appsec) GetPenaltyBoxConditions(ctx context.Context, params GetPenaltyBoxConditionsRequest) (*GetPenaltyBoxConditionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetPenaltyBoxConditions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/penalty-box/conditions",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetPenaltyBoxCondition request: %w", err)
	}

	var result GetPenaltyBoxConditionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get penalty box condition request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdatePenaltyBoxConditions(ctx context.Context, params UpdatePenaltyBoxConditionsRequest) (*UpdatePenaltyBoxConditionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdatePenaltyBoxConditions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/penalty-box/conditions",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed update penalty box condition request: %w", err)
	}

	var result UpdatePenaltyBoxConditionsResponse
	resp, err := p.Exec(req, &result, params.ConditionsPayload)
	if err != nil {
		return nil, fmt.Errorf("update penalty box condition request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
