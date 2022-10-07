package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsLogging interface supports retrieving, updating or removing settings
	// related to HTTP header logging.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#headerlog
	AdvancedSettingsLogging interface {
		// AdvancedSettingsLogging lists the HTTP header logging settings for a configuration or policy. If
		// the request specifies a policy, then the settings for that policy will be returned, otherwise the
		// settings for the configuration will be returned.
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#gethttpheaderloggingforaconfiguration
		GetAdvancedSettingsLogging(ctx context.Context, params GetAdvancedSettingsLoggingRequest) (*GetAdvancedSettingsLoggingResponse, error)

		// UpdateAdvancedSettingsLogging enables, disables, or updates the HTTP header logging settings for a
		// configuration or policy. If the request specifies a policy, then the settings for that policy will
		// updated, otherwise the settings for the configuration will be updated.
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#puthttpheaderloggingforaconfiguration
		UpdateAdvancedSettingsLogging(ctx context.Context, params UpdateAdvancedSettingsLoggingRequest) (*UpdateAdvancedSettingsLoggingResponse, error)

		// RemoveAdvancedSettingsLogging disables HTTP header logging for a confguration or policy. If the request
		// specifies a policy, then header logging will be disabled for that policy, otherwise logging will be
		// disabled for the configuration.
		RemoveAdvancedSettingsLogging(ctx context.Context, params RemoveAdvancedSettingsLoggingRequest) (*RemoveAdvancedSettingsLoggingResponse, error)
	}

	// GetAdvancedSettingsLoggingRequest is used to retrieve the HTTP header logging settings for a configuration or policy.
	GetAdvancedSettingsLoggingRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	// GetAdvancedSettingsLoggingResponse is returned from a call to GetAdvancedSettingsLogging.
	GetAdvancedSettingsLoggingResponse struct {
		Override        json.RawMessage                  `json:"override,omitempty"`
		AllowSampling   bool                             `json:"allowSampling"`
		Cookies         *AdvancedSettingsCookies         `json:"cookies,omitempty"`
		CustomHeaders   *AdvancedSettingsCustomHeaders   `json:"customHeaders,omitempty"`
		StandardHeaders *AdvancedSettingsStandardHeaders `json:"standardHeaders,omitempty"`
	}

	// AdvancedSettingsCookies contains a list of cookie headers to be logged or not logged depending on the Type field.
	AdvancedSettingsCookies struct {
		Type   string   `json:"type"`
		Values []string `json:"values,omitempty"`
	}

	// AdvancedSettingsCustomHeaders contains a list of custom headers to be logged or not logged depending on the Type field.
	AdvancedSettingsCustomHeaders struct {
		Type   string   `json:"type,omitempty"`
		Values []string `json:"values,omitempty"`
	}

	// AdvancedSettingsStandardHeaders contains a list of standard headers to be logged or not logged depending on the Type field.
	AdvancedSettingsStandardHeaders struct {
		Type   string   `json:"type,omitempty"`
		Values []string `json:"values,omitempty"`
	}

	// UpdateAdvancedSettingsLoggingRequest is used to update the HTTP header logging settings for a configuration or policy.
	UpdateAdvancedSettingsLoggingRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// UpdateAdvancedSettingsLoggingResponse is returned from a call to UpdateAdvancedSettingsLogging.
	UpdateAdvancedSettingsLoggingResponse struct {
		Override      bool `json:"override"`
		AllowSampling bool `json:"allowSampling"`
		Cookies       struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"cookies"`
		CustomHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"customHeaders"`
		StandardHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"standardHeaders"`
	}

	// RemoveAdvancedSettingsLoggingRequest is used to disable HTTP header logging for a configuration or policy.
	RemoveAdvancedSettingsLoggingRequest struct {
		ConfigID      int    `json:"-"`
		Version       int    `json:"-"`
		PolicyID      string `json:"-"`
		Override      bool   `json:"override"`
		AllowSampling bool   `json:"allowSampling"`
	}

	// RemoveAdvancedSettingsLoggingResponse is returned from a call to RemoveAdvancedSettingsLogging.
	RemoveAdvancedSettingsLoggingResponse struct {
		Override      bool `json:"override"`
		AllowSampling bool `json:"allowSampling"`
		Cookies       struct {
			Type string `json:"type"`
		} `json:"cookies"`
		CustomHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"customHeaders"`
		StandardHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"standardHeaders"`
	}
)

// Validate validates a GetAdvancedSettingsLoggingsRequest.
func (v GetAdvancedSettingsLoggingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateAdvancedSettingsLoggingRequest.
func (v UpdateAdvancedSettingsLoggingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a RemoveAdvancedSettingsLoggingRequest.
func (v RemoveAdvancedSettingsLoggingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetAdvancedSettingsLogging(ctx context.Context, params GetAdvancedSettingsLoggingRequest) (*GetAdvancedSettingsLoggingResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsLoggings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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
		return nil, fmt.Errorf("failed to create GetAdvancedSettingsLogging request: %w", err)
	}

	var result GetAdvancedSettingsLoggingResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get advanced settings logging request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateAdvancedSettingsLogging(ctx context.Context, params UpdateAdvancedSettingsLoggingRequest) (*UpdateAdvancedSettingsLoggingResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsLogging")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsLogging request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var result UpdateAdvancedSettingsLoggingResponse
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("update advanced settings logging request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveAdvancedSettingsLogging(ctx context.Context, params RemoveAdvancedSettingsLoggingRequest) (*RemoveAdvancedSettingsLoggingResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveAdvancedSettingsLogging")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveAdvancedSettingsLogging request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var result RemoveAdvancedSettingsLoggingResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove advanced settings logging request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
