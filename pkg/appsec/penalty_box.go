package appsec

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The PenaltyBox interface supports retrieving or modifying the penalty box settings for
	// a specified security policy
	PenaltyBox interface {
		// GetPenaltyBox returns the penalty box settings for the security policy you specify.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-penalty-box
		GetPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error)

		// UpdatePenaltyBox modifies the penalty box settings for a security policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-policy-penalty-box
		UpdatePenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error)
	}

	// GetPenaltyBoxRequest is used to retrieve the penalty box settings.
	GetPenaltyBoxRequest struct {
		ConfigID             int    `json:"-"`
		Version              int    `json:"-"`
		PolicyID             string `json:"-"`
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	// GetPenaltyBoxResponse is returned from a call to GetPenaltyBox.
	GetPenaltyBoxResponse struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection,omitempty"`
	}

	// UpdatePenaltyBoxRequest is used to modify the penalty box settings.
	UpdatePenaltyBoxRequest struct {
		ConfigID             int    `json:"-"`
		Version              int    `json:"-"`
		PolicyID             string `json:"-"`
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	// UpdatePenaltyBoxResponse is returned from a call to UpdatePenaltyBox.
	UpdatePenaltyBoxResponse struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}
)

// Validate validates a GetPenaltyBoxRequest.
func (v GetPenaltyBoxRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdatePenaltyBoxRequest.
func (v UpdatePenaltyBoxRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Action":   validation.Validate(v.Action, validation.Required, validateActions),
	}.Filter()
}

func (p *appsec) GetPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetPenaltyBox")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/penalty-box",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetPenaltyBox request: %w", err)
	}

	var result GetPenaltyBoxResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get penalty box request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdatePenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdatePenaltyBox")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/penalty-box",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdatePenaltyBox request: %w", err)
	}

	var result UpdatePenaltyBoxResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update penalty box request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

// validateActions ensures actions are correct for API call.
// Allowed values: alert, deny, deny_custom_{custom_deny_id}, none.
var validateActions = validation.By(func(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	allowed := map[string]bool{"alert": true, "deny": true, "none": true}
	if allowed[v] {
		return nil
	}
	if !strings.HasPrefix(v, "deny_custom_") || strings.TrimPrefix(v, "deny_custom_") == "" {
		return fmt.Errorf("value '%s' is invalid. Must be one of: alert, deny, deny_custom_{custom_deny_id}, none", v)
	}
	return nil
})
