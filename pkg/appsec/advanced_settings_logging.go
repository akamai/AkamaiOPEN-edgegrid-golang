package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// AdvancedSettingsLogging represents a collection of AdvancedSettingsLogging
//
// See: AdvancedSettingsLogging.GetAdvancedSettingsLogging()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// AdvancedSettingsLogging  contains operations available on AdvancedSettingsLogging  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getadvancedsettingslogging
	AdvancedSettingsLogging interface {
		GetAdvancedSettingsLogging(ctx context.Context, params GetAdvancedSettingsLoggingRequest) (*GetAdvancedSettingsLoggingResponse, error)
		UpdateAdvancedSettingsLogging(ctx context.Context, params UpdateAdvancedSettingsLoggingRequest) (*UpdateAdvancedSettingsLoggingResponse, error)
	}

	GetAdvancedSettingsLoggingRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	GetAdvancedSettingsLoggingResponse struct {
		AllowSampling bool `json:"allowSampling"`
		Cookies       struct {
			Type string `json:"type"`
		} `json:"cookies"`
		CustomHeaders struct {
			Type   string   `json:"type,omitempty"`
			Values []string `json:"values,omitempty"`
		} `json:"customHeaders"`
		StandardHeaders struct {
			Type   string   `json:"type,omitempty"`
			Values []string `json:"values,omitempty"`
		} `json:"standardHeaders"`
	}

	UpdateAdvancedSettingsLoggingRequest struct {
		ConfigID      int    `json:"-"`
		Version       int    `json:"-"`
		PolicyID      string `json:"-"`
		AllowSampling bool   `json:"allowSampling"`
		Cookies       struct {
			Type string `json:"type"`
		} `json:"cookies"`
		CustomHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values"`
		} `json:"customHeaders"`
		StandardHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values"`
		} `json:"standardHeaders"`
	}
	UpdateAdvancedSettingsLoggingResponse struct {
		AllowSampling bool `json:"allowSampling"`
		Cookies       struct {
			Type string `json:"type"`
		} `json:"cookies"`
		CustomHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values"`
		} `json:"customHeaders"`
		StandardHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values"`
		} `json:"standardHeaders"`
	}
	RemoveAdvancedSettingsLoggingRequest struct {
		ConfigID      int    `json:"-"`
		Version       int    `json:"-"`
		PolicyID      string `json:"-"`
		AllowSampling bool   `json:"allowSampling"`
		Cookies       struct {
			Type string `json:"type"`
		} `json:"cookies"`
		CustomHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values"`
		} `json:"customHeaders"`
		StandardHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values"`
		} `json:"standardHeaders"`
	}

	RemoveAdvancedSettingsLoggingResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates GetAdvancedSettingsLoggingsRequest
func (v GetAdvancedSettingsLoggingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateAdvancedSettingsLoggingRequest
func (v UpdateAdvancedSettingsLoggingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetAdvancedSettingsLogging(ctx context.Context, params GetAdvancedSettingsLoggingRequest) (*GetAdvancedSettingsLoggingResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsLoggings")

	var rval GetAdvancedSettingsLoggingResponse
	var uri string

	if params.PolicyID != "" {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/logging",
			params.ConfigID,
			params.Version,
			params.PolicyID)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/advanced-settings/logging",
			params.ConfigID,
			params.Version)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getadvancedsettingsloggings request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getadvancedsettingsloggings request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a AdvancedSettingsLogging.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putadvancedsettingslogging

func (p *appsec) UpdateAdvancedSettingsLogging(ctx context.Context, params UpdateAdvancedSettingsLoggingRequest) (*UpdateAdvancedSettingsLoggingResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsLogging")

	var putURL string
	if params.PolicyID != "" {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/logging",
			params.ConfigID,
			params.Version,
			params.PolicyID)
	} else {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/advanced-settings/logging",
			params.ConfigID,
			params.Version)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create AdvancedSettingsLoggingrequest: %w", err)
	}

	var rval UpdateAdvancedSettingsLoggingResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create AdvancedSettingsLogging request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
