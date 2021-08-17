package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsPragma interface supports retrieving or modifying the pragma header
	// excluded conditions for a configuration or policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#pragmaheader
	AdvancedSettingsPragma interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpragmaheaderconfiguration
		GetAdvancedSettingsPragma(ctx context.Context, params GetAdvancedSettingsPragmaRequest) (*GetAdvancedSettingsPragmaResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpragmaheaderconfiguration
		UpdateAdvancedSettingsPragma(ctx context.Context, params UpdateAdvancedSettingsPragmaRequest) (*UpdateAdvancedSettingsPragmaResponse, error)
	}

	// GetAdvancedSettingsPragmaRequest is used to retrieve the pragma settings for a security policy.
	GetAdvancedSettingsPragmaRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	// GetAdvancedSettingsPragmaResponse is returned from a call to GetAdvancedSettingsPragma.
	GetAdvancedSettingsPragmaResponse struct {
		Action            string             `json:"action,,omitempty"`
		ConditionOperator string             `json:"conditionOperator,omitempty"`
		ExcludeCondition  []ExcludeCondition `json:"excludeCondition,omitempty"`
	}

	// ExcludeCondition describes the pragma header's excluded conditions.
	ExcludeCondition struct {
		Type          string   `json:"type"`
		PositiveMatch bool     `json:"positiveMatch"`
		Header        string   `json:"header"`
		Value         []string `json:"value"`
		Name          string   `json:"name"`
		ValueCase     bool     `json:"valueCase"`
		ValueWildcard bool     `json:"valueWildcard"`
		UseHeaders    bool     `json:"useHeaders"`
	}

	// UpdateAdvancedSettingsPragmaRequest is used to modify the pragma settings for a security policy.
	UpdateAdvancedSettingsPragmaRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// UpdateAdvancedSettingsPragmaResponse is returned from a call to UpdateAdvancedSettingsPragma.
	UpdateAdvancedSettingsPragmaResponse struct {
		Action            string             `json:"action"`
		ConditionOperator string             `json:"conditionOperator"`
		ExcludeCondition  []ExcludeCondition `json:"excludeCondition"`
	}
)

// Validate validates a GetAdvancedSettingsPragmaRequest.
func (v GetAdvancedSettingsPragmaRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateAdvancedSettingsPragmaRequest.
func (v UpdateAdvancedSettingsPragmaRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetAdvancedSettingsPragma(ctx context.Context, params GetAdvancedSettingsPragmaRequest) (*GetAdvancedSettingsPragmaResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsPragma")

	var rval GetAdvancedSettingsPragmaResponse
	var uri string

	if params.PolicyID != "" {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/pragma-header",
			params.ConfigID,
			params.Version,
			params.PolicyID)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/advanced-settings/pragma-header",
			params.ConfigID,
			params.Version)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAdvancedSettingsPragma request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetAdvancedSettingsPragma request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateAdvancedSettingsPragma(ctx context.Context, params UpdateAdvancedSettingsPragmaRequest) (*UpdateAdvancedSettingsPragmaResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsPragma")

	var uri string
	if params.PolicyID != "" {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/pragma-header",
			params.ConfigID,
			params.Version,
			params.PolicyID)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/advanced-settings/pragma-header",
			params.ConfigID,
			params.Version)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsPragma request: %w", err)
	}

	var rval UpdateAdvancedSettingsPragmaResponse
	resp, err := p.Exec(req, &rval, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("UpdateAdvancedSettingsPragma request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
