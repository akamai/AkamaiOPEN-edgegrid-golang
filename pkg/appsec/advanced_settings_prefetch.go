package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// AdvancedSettingsPrefetch represents a collection of AdvancedSettingsPrefetch
//
// See: AdvancedSettingsPrefetch.GetAdvancedSettingsPrefetch()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// AdvancedSettingsPrefetch  contains operations available on AdvancedSettingsPrefetch  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getadvancedsettingsprefetch
	AdvancedSettingsPrefetch interface {
		GetAdvancedSettingsPrefetch(ctx context.Context, params GetAdvancedSettingsPrefetchRequest) (*GetAdvancedSettingsPrefetchResponse, error)
		UpdateAdvancedSettingsPrefetch(ctx context.Context, params UpdateAdvancedSettingsPrefetchRequest) (*UpdateAdvancedSettingsPrefetchResponse, error)
	}

	GetAdvancedSettingsPrefetchRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	GetAdvancedSettingsPrefetchResponse struct {
		AllExtensions      bool     `json:"allExtensions"`
		EnableAppLayer     bool     `json:"enableAppLayer"`
		EnableRateControls bool     `json:"enableRateControls"`
		Extensions         []string `json:"extensions,omitempty"`
	}

	UpdateAdvancedSettingsPrefetchRequest struct {
		ConfigID           int      `json:"-"`
		Version            int      `json:"-"`
		AllExtensions      bool     `json:"allExtensions"`
		EnableAppLayer     bool     `json:"enableAppLayer"`
		EnableRateControls bool     `json:"enableRateControls"`
		Extensions         []string `json:"extensions,omitempty"`
	}
	UpdateAdvancedSettingsPrefetchResponse struct {
		AllExtensions      bool     `json:"allExtensions"`
		EnableAppLayer     bool     `json:"enableAppLayer"`
		EnableRateControls bool     `json:"enableRateControls"`
		Extensions         []string `json:"extensions,omitempty"`
	}

	RemoveAdvancedSettingsPrefetchRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Action   string `json:"action"`
	}

	RemoveAdvancedSettingsPrefetchResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates GetAdvancedSettingsPrefetchRequest
func (v GetAdvancedSettingsPrefetchRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateAdvancedSettingsPrefetchRequest
func (v UpdateAdvancedSettingsPrefetchRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetAdvancedSettingsPrefetch(ctx context.Context, params GetAdvancedSettingsPrefetchRequest) (*GetAdvancedSettingsPrefetchResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsPrefetch")

	var rval GetAdvancedSettingsPrefetchResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/prefetch",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getadvancedsettingsprefetch request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getadvancedsettingsprefetch  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a AdvancedSettingsPrefetch.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putadvancedsettingsprefetch

func (p *appsec) UpdateAdvancedSettingsPrefetch(ctx context.Context, params UpdateAdvancedSettingsPrefetchRequest) (*UpdateAdvancedSettingsPrefetchResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsPrefetch")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/prefetch",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create AdvancedSettingsPrefetchrequest: %w", err)
	}

	var rval UpdateAdvancedSettingsPrefetchResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create AdvancedSettingsPrefetch request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
