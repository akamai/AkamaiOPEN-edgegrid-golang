package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// AdvancedSettingsPragma represents a collection of AdvancedSettingsPragma
//
// See: AdvancedSettingsPragma.GetAdvancedSettingsPragma()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// AdvancedSettingsPragma  contains operations available on AdvancedSettingsPragma  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpragmaheaderconfiguration
	AdvancedSettingsPragma interface {
		GetAdvancedSettingsPragma(ctx context.Context, params GetAdvancedSettingsPragmaRequest) (*GetAdvancedSettingsPragmaResponse, error)
		UpdateAdvancedSettingsPragma(ctx context.Context, params UpdateAdvancedSettingsPragmaRequest) (*UpdateAdvancedSettingsPragmaResponse, error)
	}

	GetAdvancedSettingsPragmaRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	GetAdvancedSettingsPragmaResponse struct {
		Action            string             `json:"action,,omitempty"`
		ConditionOperator string             `json:"conditionOperator,omitempty"`
		ExcludeCondition  []ExcludeCondition `json:"excludeCondition,omitempty"`
	}

	ExcludeCondition struct {
		Type          string   `json:"type"`
		PositiveMatch bool     `json:"positiveMatch"`
		Header        string   `json:"header"`
		Value         []string `json:"value"`
		Name          string   `json:"name"`
		ValueCase     bool     `json:"valueCase"`
		ValueWildcard bool     `json:"valueWildcard"`
		useHeaders    bool     `json:"useHeaders"`
	}

	UpdateAdvancedSettingsPragmaRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}
	UpdateAdvancedSettingsPragmaResponse struct {
		Action            string             `json:"action"`
		ConditionOperator string             `json:"conditionOperator"`
		ExcludeCondition  []ExcludeCondition `json:"excludeCondition"`
	}
)

// Validate validates GetAdvancedSettingsPragmaRequest
func (v GetAdvancedSettingsPragmaRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateAdvancedSettingsPragmaRequest
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
		return nil, fmt.Errorf("failed to get getadvancedsettingpragma request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getadvancedsettingspragma request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a AdvancedSettingsPragma.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putadvancedsettingsprefetch

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
		return nil, fmt.Errorf("failed to create create AdvancedSettingsPragmarequest: %w", err)
	}

	var rval UpdateAdvancedSettingsPragmaResponse
	resp, err := p.Exec(req, &rval, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create AdvancedSettingsPragma request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
