package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// AdvancedSettingsPolicyLogging represents a collection of AdvancedSettingsPolicyLogging
//
// See: AdvancedSettingsPolicyLogging.GetAdvancedSettingsPolicyLogging()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// AdvancedSettingsPolicyLogging  contains operations available on AdvancedSettingsPolicyLogging  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getadvancedsettingspolicylogging
	AdvancedSettingsPolicyLogging interface {
		//GetAdvancedSettingsPolicyLoggings(ctx context.Context, params GetAdvancedSettingsPolicyLoggingsRequest) (*GetAdvancedSettingsPolicyLoggingsResponse, error)
		GetAdvancedSettingsPolicyLogging(ctx context.Context, params GetAdvancedSettingsPolicyLoggingRequest) (*GetAdvancedSettingsPolicyLoggingResponse, error)
		UpdateAdvancedSettingsPolicyLogging(ctx context.Context, params UpdateAdvancedSettingsPolicyLoggingRequest) (*UpdateAdvancedSettingsPolicyLoggingResponse, error)
	}

	GetAdvancedSettingsPolicyLoggingRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	GetAdvancedSettingsPolicyLoggingResponse struct {
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

	UpdateAdvancedSettingsPolicyLoggingRequest struct {
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
	UpdateAdvancedSettingsPolicyLoggingResponse struct {
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
	RemoveAdvancedSettingsPolicyLoggingRequest struct {
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

	RemoveAdvancedSettingsPolicyLoggingResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates GetAdvancedSettingsPolicyLoggingRequest
func (v GetAdvancedSettingsPolicyLoggingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}



// Validate validates UpdateAdvancedSettingsPolicyLoggingRequest
func (v UpdateAdvancedSettingsPolicyLoggingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		
	}.Filter()
}

func (p *appsec) GetAdvancedSettingsPolicyLogging(ctx context.Context, params GetAdvancedSettingsPolicyLoggingRequest) (*GetAdvancedSettingsPolicyLoggingResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}
logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsPolicyLogging")

	var rval GetAdvancedSettingsPolicyLoggingResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/logging",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getadvancedsettingspolicylogging request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getadvancedsettingspolicylogging  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetAdvancedSettingsPolicyLoggings(ctx context.Context, params GetAdvancedSettingsPolicyLoggingRequest) (*GetAdvancedSettingsPolicyLoggingResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsPolicyLoggings")

	var rval GetAdvancedSettingsPolicyLoggingResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/logging",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getadvancedsettingspolicyloggings request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getadvancedsettingspolicyloggings request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a AdvancedSettingsPolicyLogging.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putadvancedsettingspolicylogging

func (p *appsec) UpdateAdvancedSettingsPolicyLogging(ctx context.Context, params UpdateAdvancedSettingsPolicyLoggingRequest) (*UpdateAdvancedSettingsPolicyLoggingResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsPolicyLogging")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/logging",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create AdvancedSettingsPolicyLoggingrequest: %w", err)
	}

	var rval UpdateAdvancedSettingsPolicyLoggingResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create AdvancedSettingsPolicyLogging request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
